//nolint:golint
package cert_manager

import (
	"bytes"
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/rekby/lets-proxy2/internal/cache"

	"github.com/rekby/lets-proxy2/internal/log"

	zc "github.com/rekby/zapcontext"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme"
)

const (
	tlsAlpn01     = "tls-alpn-01"
	http01        = "http-01"
	httpWellKnown = "/.well-known/acme-challenge/"
)

const domainKeyRSALength = 2048

var errHaveNoCert = errors.New("have no certificate for domain")

//nolint:varcheck,deadcode,unused
var errNotImplementedError = errors.New("not implemented yet")

type GetContext interface {
	GetContext() context.Context
}

type keyType string

const keyRSA keyType = "rsa"

// Interface inspired to https://godoc.org/golang.org/x/crypto/acme/autocert#Manager but not compatible guarantee
type Manager struct {
	CertificateIssueTimeout time.Duration
	Cache                   cache.Cache

	// Client is used to perform low-level operations, such as account registration
	// and requesting new certificates.
	//
	// If Client is nil, a zero-value acme.Client is used with acme.LetsEncryptURL
	// as directory endpoint. If the Client.Key is nil, a new ECDSA P-256 key is
	// generated and, if Cache is not nil, stored in cache.
	//
	// Mutating the field after the first call of GetCertificate method will have no effect.
	Client               *acme.Client
	EnableHTTPValidation bool
	EnableTLSValidation  bool

	certForDomainAuthorize cache.Value

	certStateMu sync.Mutex
	certState   map[certNameType]*certState

	httpTokens *cache.MemoryCache
}

func New(ctx context.Context, client *acme.Client, c cache.Cache) *Manager {
	res := Manager{}
	res.Client = client
	res.certForDomainAuthorize = cache.NewMemoryValue("authcert")
	res.certState = make(map[certNameType]*certState)
	res.CertificateIssueTimeout = time.Minute
	res.httpTokens = cache.NewMemoryCache("Http validation tokens")
	res.Cache = c
	res.EnableTLSValidation = true
	return &res
}

// GetCertificate implements the tls.Config.GetCertificate hook.
//
func (m *Manager) GetCertificate(hello *tls.ClientHelloInfo) (resultCert *tls.Certificate, err error) {
	// TODO: get context of connection
	var ctx context.Context
	if getContext, ok := hello.Conn.(GetContext); ok {
		ctx = getContext.GetContext()
	} else {
		defaultLogger := zc.L(context.Background())
		defaultLogger.DPanic("hello.Conn must implement GetContext interface")
		ctx = zc.WithLogger(context.Background(), defaultLogger)
	}

	needDomain := normalizeDomain(hello.ServerName)

	logger := zc.L(ctx).With(logDomain(needDomain))
	logger.Info("Get certificate", zap.String("original_domain", hello.ServerName))
	if isTLSALPN01Hello(hello) {
		logger.Debug("It is tls-alpn-01 token request.")

		certInterface, err := m.certForDomainAuthorize.Get(ctx, needDomain.String())
		logger.Debug("Got authcert from cache", zap.Error(err))

		cert, _ := certInterface.(*tls.Certificate)

		if cert == nil {
			logger.Warn("Doesn't have token for request domain")
			return nil, errHaveNoCert
		}
		return cert, nil
	}

	certName := certNameFromDomain(needDomain)

	logger = logger.With(logCetName(certName))
	ctx = zc.WithLogger(ctx, zc.L(ctx).With(logCetName(certName)))

	now := time.Now()
	defer func() {
		if isNeedRenew(resultCert, now) {
			go m.renewCertInBackground(ctx, certName)
		}
	}()

	certState := m.certStateGet(certName)
	cert, err := certState.Cert()
	if cert != nil {
		logger.Debug("Got certificate from local state", log.Cert(cert))

		// TODO: Disable check for locked certificates https://github.com/rekby/lets-proxy2/issues/48
		cert, err = validCertDer([]DomainName{needDomain}, cert.Certificate, cert.PrivateKey, now)
		logger.Debug("Validate certificate from local state", zap.Error(err))
		if err == nil {
			return cert, nil
		}
	}
	if err != nil {
		logger.Debug("Can't get certificate from local state", zap.Error(err))
	}

	cert, err = getCertificate(ctx, m.Cache, certName, keyRSA)
	if err == nil {
		logger.Debug("Certificate loaded from cache")

		// TODO: Disable check for locked certificates https://github.com/rekby/lets-proxy2/issues/48
		cert, err = validCertDer([]DomainName{needDomain}, cert.Certificate, cert.PrivateKey, now)
		logger.Debug("Check if certificate ok", zap.Error(err))
		if err == nil {
			certState.CertSet(ctx, cert)
			return cert, nil
		}
	}

	// TODO: check domain
	certIssueContext, cancelFunc := context.WithTimeout(ctx, m.CertificateIssueTimeout)
	defer cancelFunc()

	// TODO: receive cert for domain and subdomains same time
	res, err := m.createCertificateForDomains(certIssueContext, certName, domainNamesFromCertificateName(certName),
		needDomain)
	if err == nil {
		logger.Info("Certificate issued.", log.Cert(res),
			zap.Time("expire", res.Leaf.NotAfter))
		return res, nil
	}
	logger.Warn("Can't issue certificate", zap.Error(err))
	return nil, errHaveNoCert

}

func (m *Manager) certStateGet(certName certNameType) *certState {
	m.certStateMu.Lock()

	res := m.certState[certName]
	if res == nil {
		res = &certState{}
		m.certState[certName] = res
	}
	m.certStateMu.Unlock()
	return res
}

func (m *Manager) createCertificateForDomains(ctx context.Context, certName certNameType,
	domainNames []DomainName, needDomain DomainName) (res *tls.Certificate, err error) {

	logger := zc.L(ctx)
	certState := m.certStateGet(certName)
	if certState.StartIssue(ctx) {
		// outer func need for get argument values in defer time
		defer func() {
			certState.FinishIssue(ctx, res, err)
		}()
	} else {
		waitTimeout, waitTimeoutCancel := context.WithTimeout(ctx, m.CertificateIssueTimeout)
		defer waitTimeoutCancel()

		return certState.WaitFinishIssue(waitTimeout)
	}

	var authorizeDomainsWg sync.WaitGroup
	var authorizedDomainsMu sync.Mutex
	var authorizedDomains []DomainName
	var needDomainAuthorized = false

	authorizeDomainsWg.Add(len(domainNames))
	for _, authorizeDomain := range domainNames {
		go func(domain DomainName) {
			localLogger := logger.With(logDomain(domain))
			localCtx := zc.WithLogger(ctx, localLogger)

			err = m.authorizeDomain(localCtx, domain)
			if err == nil {
				localLogger.Debug("Domain authorized.")

				authorizedDomainsMu.Lock()
				authorizedDomains = append(authorizedDomains, domain)
				if domain == needDomain {
					needDomainAuthorized = true
				}
				authorizedDomainsMu.Unlock()

			} else {
				localLogger.Warn("Domain doesn't authorized.", zap.Error(err))
			}
			authorizeDomainsWg.Done()
		}(authorizeDomain)
	}
	logger.Debug("Wait for domains authorization.", logDomains(domainNames))
	authorizeDomainsWg.Wait()

	if !needDomainAuthorized {
		logger.Warn("Need domain doesn't authorized.", logDomain(needDomain),
			logDomainsNamed("authorized_domains", authorizedDomains))
		return nil, errors.New("need domain doesn't authorized")
	}

	res, err = m.issueCertificate(ctx, certName, authorizedDomains)
	if err == nil {
		logger.Debug("Certificate created.")
	} else {
		logger.Warn("Can't issue certificate", zap.Error(err))
	}
	return res, err
}

func (m *Manager) authorizeDomain(ctx context.Context, domain DomainName) error {
	logger := zc.L(ctx).With(logDomain(domain))
	pendingAuthorizations := make(map[string]struct{})
	defer func() {
		if len(pendingAuthorizations) == 0 {
			logger.Debug("No pending authorizations to clean")
			return
		}

		var authUries []string
		for k := range pendingAuthorizations {
			authUries = append(authUries, k)
		}

		logger.Info("Start detached process for revoke pending authorizations",
			zap.Strings("uries", authUries))

		revokeContext := context.Background()
		revokeContext = zc.WithLogger(revokeContext, logger.With(zap.Bool("detached", true)))
		go m.revokePendingAuthorizations(revokeContext, authUries)
	}()

	var allowedChallenges []string
	if m.EnableTLSValidation {
		allowedChallenges = append(allowedChallenges, tlsAlpn01)
	}
	if m.EnableHTTPValidation {
		allowedChallenges = append(allowedChallenges, http01)
	}

	var nextChallengeTypeIndex int
	for {
		// Start domain authorization and get the challenge.
		authz, err := m.Client.Authorize(ctx, domain.String())
		if err == nil {
			logger.Debug("Got authorization description.", zap.Reflect("auth_object", authz))
		} else {
			logger.Error("Can't get domain authorization description", zap.Error(err))
			return err
		}

		// No point in accepting challenges if the authorization status
		// is in a final state.
		switch authz.Status {
		case acme.StatusValid:
			logger.Debug("Domain already authorized")
			return nil // already authorized
		case acme.StatusInvalid:
			logger.Warn("Domain has invalid authorization", zap.String("auth_uri", authz.URI))
			return errors.New("invalid authorization status")
		}

		pendingAuthorizations[authz.URI] = struct{}{}

		// Pick the next preferred challenge.
		var chal *acme.Challenge
		for chal == nil && nextChallengeTypeIndex < len(allowedChallenges) {
			logger.Debug("Check if accept challenge",
				zap.String("challenge", allowedChallenges[nextChallengeTypeIndex]))
			chal = pickChallenge(allowedChallenges[nextChallengeTypeIndex], authz.Challenges)
			nextChallengeTypeIndex++
		}
		if chal == nil {
			logger.Warn("Unable to authorize domain. No compatible challenges.")
			return errors.New("unable authorize domain")
		}
		logger.Debug("Select challenge for authorize", zap.Reflect("challenge", chal))

		cleanup, err := m.fulfill(ctx, chal, domain)
		if cleanup != nil {
			//noinspection GoDeferInLoop
			defer cleanup(ctx)
		}

		if err != nil {
			logger.Error("Can't set domain token", zap.Reflect("chal", chal), zap.Error(err))
			continue
		}

		receivedChallenge, err := m.Client.Accept(ctx, chal)
		if err == nil {
			logger.Debug("Receive authorize answer", zap.Reflect("challenge", receivedChallenge))
		} else {
			logger.Error("Can't authorize domain", zap.Reflect("challenge", receivedChallenge), zap.Error(err))
			continue
		}

		receivedAuth, err := m.Client.WaitAuthorization(ctx, authz.URI)
		log.DebugError(logger, err, "Receive domain authorization", zap.Reflect("authorization", receivedAuth))
		if err != nil {
			continue
		}

		logger.Info("Domain authorized")
		delete(pendingAuthorizations, authz.URI)
		return nil
	}
}

func (m *Manager) issueCertificate(ctx context.Context, certName certNameType, domains []DomainName) (*tls.Certificate, error) {
	if len(domains) == 0 {
		return nil, errors.New("no domains for issue certificate")
	}
	logger := zc.L(ctx).With(logDomains(domains))

	key, err := m.certKeyGetOrCreate(ctx, certName, keyRSA)
	if err != nil {
		logger.Error("Can't get domain key", zap.Error(err))
		return nil, err
	}
	csr, err := createCertRequest(key, domains[0], domains...)
	log.DebugDPanic(logger, err, "Create certificate request")
	if err != nil {
		return nil, err
	}

	der, _, err := m.Client.CreateCert(ctx, csr, 0, true)
	log.InfoError(logger, err, "Receive certificate from acme server")
	if err != nil {
		return nil, err
	}

	cert, err := validCertDer(domains, der, key, time.Now())
	log.DebugDPanic(logger, err, "Check certificate is valid")
	if err == nil {
		storeCertificate(ctx, m.Cache, certName, cert)
		return cert, nil
	}
	return nil, err
}

func (m *Manager) renewCertInBackground(ctx context.Context, certName certNameType) {
	// detach from request lifetime, but save log context
	logger := zc.L(ctx).Named("background")
	ctx, ctxCancel := context.WithTimeout(context.Background(), m.CertificateIssueTimeout)
	defer ctxCancel()

	ctx = zc.WithLogger(ctx, logger)
	certState := m.certStateGet(certName)

	if !certState.StartIssue(ctx) {
		// already has other cert issue process
		return
	}
	domains := domainNamesFromCertificateName(certName)
	logger.Info("Start background certificate issue")
	cert, err := m.createCertificateForDomains(ctx, certName, domains, "")
	certState.FinishIssue(ctx, cert, err)
	log.InfoError(logger, err, "Renew certificate in background finished", log.Cert(cert))
}

func (m *Manager) revokePendingAuthorizations(ctx context.Context, uries []string) {
	logger := zc.L(ctx)
	for _, uri := range uries {
		err := m.Client.RevokeAuthorization(ctx, uri)
		if err == nil {
			logger.Debug("Revoke authorization ok", zap.String("uri", uri))
		} else {
			logger.Error("Can't revoke authorization", zap.String("uri", uri), zap.Error(err))
		}
	}
}

func (m *Manager) certKeyGetOrCreate(ctx context.Context, certName certNameType, keyType keyType) (crypto.Signer, error) {
	logger := zc.L(ctx)

	if keyType != keyRSA {
		logger.DPanic("Unknown key type", zap.String("key_type", string(keyType)))
		return nil, errors.New("unknown key type for generate key")
	}

	key, err := getCertificateKey(ctx, m.Cache, certName, keyType)
	if err == nil {
		logger.Debug("Got certificate key from cache and reuse old key")
		return key, nil
	}
	if err == cache.ErrCacheMiss {
		logger.Debug("Cert key no in cache. Create new.")
	} else {
		logger.Error("Error while check cert key in cache", zap.Error(err))
		return nil, err
	}

	logger.Debug("Generate new rsa key")
	key, err = rsa.GenerateKey(rand.Reader, domainKeyRSALength)
	if err != nil {
		logger.Error("Can't generate rsa key", zap.Error(err))
		return nil, err
	}

	return key, nil
}

func (m *Manager) fulfill(ctx context.Context, challenge *acme.Challenge, domain DomainName) (func(context.Context), error) {
	logger := zc.L(ctx)
	switch challenge.Type {
	case tlsAlpn01:
		cert, err := m.Client.TLSALPN01ChallengeCert(challenge.Token, domain.String())
		if err != nil {
			return nil, err
		}
		m.putCertToken(ctx, domain, &cert)
		return func(localContext context.Context) { go m.deleteCertToken(localContext, domain) }, nil
	case http01:
		resp, err := m.Client.HTTP01ChallengeResponse(challenge.Token)
		if err != nil {
			return nil, err
		}
		key := domain.ASCII() + "/" + challenge.Token
		err = m.httpTokens.Put(ctx, key, []byte(resp))
		log.DebugError(logger, err, "Put token for http-01", zap.String("key", key))
		return func(localContext context.Context) { _ = m.httpTokens.Delete(localContext, key) }, err
	default:
		logger.Error("Unknow challenge type", zap.Reflect("challenge", challenge))
		return nil, errors.New("unknown challenge type")
	}
}

func (m *Manager) isHTTPValidationRequest(r *http.Request) bool {
	if r == nil || r.URL == nil {
		return false
	}
	if r.Method != http.MethodGet {
		return false
	}

	return strings.HasPrefix(r.URL.Path, httpWellKnown)
}

func (m *Manager) HandleHttpValidation(w http.ResponseWriter, r *http.Request) bool {
	if !m.isHTTPValidationRequest(r) {
		return false
	}

	ctx := r.Context()
	logger := zc.L(ctx)

	token := strings.TrimPrefix(r.URL.Path, httpWellKnown)
	domain := normalizeDomain(r.Host)
	resp, err := m.httpTokens.Get(ctx, domain.ASCII()+"/"+token)
	logger.Debug("Get http token", zap.Error(err))
	if err == nil {
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(resp)
		log.DebugInfo(logger, err, "Error write http token answer to response", logDomain(domain), zap.String("token", token))
	} else {
		logger.Warn("Have no validation token", logDomain(domain), zap.String("token", token), zap.Error(err))
	}
	return true
}

func (m *Manager) putCertToken(ctx context.Context, key DomainName, certificate *tls.Certificate) {
	err := m.certForDomainAuthorize.Put(ctx, key.String(), certificate)
	log.DebugDPanicCtx(ctx, err, "Put cert token", zap.String("key", string(key)))
}

func (m *Manager) deleteCertToken(ctx context.Context, key DomainName) {
	err := m.certForDomainAuthorize.Delete(ctx, key.String())
	log.DebugDPanicCtx(ctx, err, "Delete cert token", zap.String("key", key.String()))
}

// It isn't atomic syncronized - caller must not save two certificates with same name same time
func storeCertificate(ctx context.Context, cache cache.Cache, certName certNameType,
	cert *tls.Certificate) {
	logger := zc.L(ctx)
	if cache == nil {
		logger.Debug("Can't save certificate to nil cache")
		return
	}

	var keyType keyType

	var certBuf bytes.Buffer
	for _, block := range cert.Certificate {
		pemBlock := pem.Block{Type: "CERTIFICATE", Bytes: block}
		err := pem.Encode(&certBuf, &pemBlock)
		if err != nil {
			logger.DPanic("Can't encode pem block of certificate", zap.Error(err), zap.Binary("block", block))
			return
		}
	}

	var privateKeyBuf bytes.Buffer

	switch privateKey := cert.PrivateKey.(type) {
	case *rsa.PrivateKey:
		keyType = keyRSA
		keyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
		pemBlock := pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyBytes}
		err := pem.Encode(&privateKeyBuf, &pemBlock)
		if err != nil {
			logger.DPanic("Can't marshal rsa private key", zap.Error(err))
			return
		}
	default:
		logger.DPanic("Unknow private key type", zap.String("type", reflect.TypeOf(cert.PrivateKey).String()))
		return
	}

	if keyType == "" {
		logger.DPanic("store cert key type doesn't init")
	}

	certKeyName := string(certName) + "." + string(keyType) + ".cer"
	keyKeyName := string(certName) + "." + string(keyType) + ".key"

	err := cache.Put(ctx, certKeyName, certBuf.Bytes())
	if err != nil {
		logger.Error("Can't store certificate file", zap.Error(err))
		return
	}

	err = cache.Put(ctx, keyKeyName, privateKeyBuf.Bytes())
	if err != nil {
		_ = cache.Delete(ctx, certKeyName)
		logger.Error("Can't store certificate key file", zap.Error(err))
	}
}

func getCertificate(ctx context.Context, cache cache.Cache, certName certNameType, keyType keyType) (cert *tls.Certificate, err error) {
	logger := zc.L(ctx)
	logger.Debug("Check certificate in cache")
	defer func() {
		logger.Debug("Checked certificate in cache", log.Cert(cert), zap.Error(err))
	}()

	certKeyName := string(certName) + "." + string(keyType) + ".cer"

	certBytes, err := cache.Get(ctx, certKeyName)
	if err != nil {
		return nil, err
	}
	keyBytes, err := getCertificateKeyBytes(ctx, cache, certName, keyType)
	if err != nil {
		return nil, err
	}

	cert2, err := tls.X509KeyPair(certBytes, keyBytes)
	if err != nil {
		return nil, err
	}
	if len(cert2.Certificate) > 0 {
		cert2.Leaf, err = x509.ParseCertificate(cert2.Certificate[0])
		if err != nil {
			return nil, err
		}
	}
	return validCertTLS(&cert2, nil, time.Now())
}

func getCertificateKeyBytes(ctx context.Context, cache cache.Cache, certName certNameType, keyType keyType) ([]byte, error) {
	keyKeyName := string(certName) + "." + string(keyType) + ".key"
	return cache.Get(ctx, keyKeyName)
}

func getCertificateKey(ctx context.Context, cache cache.Cache, certName certNameType, keyType keyType) (crypto.Signer, error) {
	certBytes, err := getCertificateKeyBytes(ctx, cache, certName, keyType)
	if err != nil {
		return nil, err
	}
	return parsePrivateKey(certBytes)
}

func parsePrivateKey(keyPEMBlock []byte) (crypto.Signer, error) {
	// extract from tls.go, standard lib. func X509KeyPair(certPEMBlock, keyPEMBlock []byte) (Certificate, error)
	// X509KeyPair parses a public/private key pair from a pair of
	// PEM encoded data. On successful return, Certificate.Leaf will be nil because
	// the parsed form of the certificate is not retained.
	fail := func(err error) (crypto.Signer, error) { return nil, err }

	var keyDERBlock *pem.Block
	var skippedBlockTypes []string
	for {
		keyDERBlock, keyPEMBlock = pem.Decode(keyPEMBlock)
		if keyDERBlock == nil {
			return fail(fmt.Errorf("tls: failed to find PEM block with type ending in \"PRIVATE KEY\" in key "+
				"input after skipping PEM blocks of the following types: %v", skippedBlockTypes))
		}
		if keyDERBlock.Type == "PRIVATE KEY" || strings.HasSuffix(keyDERBlock.Type, " PRIVATE KEY") {
			break
		}
		skippedBlockTypes = append(skippedBlockTypes, keyDERBlock.Type)
	}

	// bedge key bytes
	der := keyDERBlock.Bytes

	// copy from tls.go, standard lib. func parsePrivateKey(der []byte) (crypto.PrivateKey, error)
	//
	// Attempt to parse the given private key DER block. OpenSSL 0.9.8 generates
	// PKCS#1 private keys by default, while OpenSSL 1.0.0 generates PKCS#8 keys.
	// OpenSSL ecparam generates SEC1 EC private keys for ECDSA. We try all three.
	// func parsePrivateKey(der []byte) (crypto.PrivateKey, error) {
	if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
		return key, nil
	}
	if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
		switch key := key.(type) {
		// change from original. separate case need for allow return signer interface
		case *rsa.PrivateKey:
			return key, nil
		case *ecdsa.PrivateKey:
			return key, nil
		default:
			return nil, errors.New("tls: found unknown private key type in PKCS#8 wrapping")
		}
	}
	if key, err := x509.ParseECPrivateKey(der); err == nil {
		return key, nil
	}

	return nil, errors.New("tls: failed to parse private key")
}

func isNeedRenew(cert *tls.Certificate, now time.Time) bool {
	if cert == nil || cert.Leaf == nil {
		return false
	}
	return cert.Leaf.NotAfter.Add(-time.Hour * 24 * 30).Before(now)
}
