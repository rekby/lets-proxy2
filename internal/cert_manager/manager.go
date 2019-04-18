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

	"github.com/rekby/zapcontext"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme"
)

const (
	tlsAlpn01     = "tls-alpn-01"
	http01        = "http-01"
	httpWellKnown = "/.well-known/acme-challenge/"
)

var (
	globalAllowedChallenges = []string{tlsAlpn01}
	domainKeyRSALength      = 2048
)

var haveNoCert = errors.New("have no certificate for domain")
var notImplementedError = errors.New("not implemented yet")

type GetContext interface {
	GetContext() context.Context
}

type keyType string

const keyRSA keyType = "rsa"

// Interface inspired to https://godoc.org/golang.org/x/crypto/acme/autocert#Manager but not compatible gurantee
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
	EnableHttpValidation bool
	EnableTlsValidation  bool

	// will rewrite to Cache in future
	// https://github.com/rekby/lets-proxy2/issues/32
	certTokensMu sync.RWMutex
	certTokens   map[DomainName]*tls.Certificate

	certStateMu sync.Mutex
	certState   map[certNameType]*certState

	httpTokens *cache.MemoryCache
}

func New(ctx context.Context, client *acme.Client) *Manager {
	res := Manager{}
	res.Client = client
	res.certTokens = make(map[DomainName]*tls.Certificate)
	res.certState = make(map[certNameType]*certState)
	res.CertificateIssueTimeout = time.Minute
	res.httpTokens = cache.NewMemoryCache("Http validation tokens")
	res.EnableTlsValidation = true
	return &res
}

// GetCertificate implements the tls.Config.GetCertificate hook.
//
func (m *Manager) GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	// TODO: get context of connection
	var ctx context.Context
	if getContext, ok := hello.Conn.(GetContext); ok {
		ctx = getContext.GetContext()
	} else {
		defaultLogger := zc.L(nil)
		defaultLogger.DPanic("hello.Conn must implement GetContext interface")
		ctx = zc.WithLogger(context.Background(), defaultLogger)
	}

	needDomain := normalizeDomain(hello.ServerName)

	logger := zc.L(ctx).With(logDomain(needDomain))
	logger.Info("Get certificate", zap.String("original_domain", hello.ServerName))
	if isTlsAlpn01Hello(hello) {
		logger.Debug("It is tls-alpn-01 token request.")

		m.certTokensMu.RLock()
		cert := m.certTokens[needDomain]
		m.certTokensMu.RUnlock()

		if cert == nil {
			logger.Warn("Doesn't have token for request domain")
			return nil, haveNoCert
		}
		return cert, nil
	}

	// TODO: check disk cache

	certName := certNameFromDomain(needDomain)

	logger = logger.With(logCetName(certName))
	ctx = zc.WithLogger(ctx, zc.L(ctx).With(logCetName(certName)))

	certState := m.certStateGet(certName)
	cert, err := certState.Cert()
	if cert != nil {
		logger.Debug("Got certificate from local state", log.Cert(cert))

		// Disable check for locked certificates
		cert, err = validCertDer([]DomainName{needDomain}, cert.Certificate, cert.PrivateKey, time.Now())
		if err != nil {
			logger.Debug("In memory cached certificate doesn't valid. Issue new.", log.Cert(cert), zap.Error(err))
		}
	}
	if err != nil {
		logger.Debug("Can't get certificate from local state", zap.Error(err))
	}

	// TODO: check domain
	certIssueContext, cancelFunc := context.WithTimeout(ctx, m.CertificateIssueTimeout)
	defer cancelFunc()

	// TODO: receive cert for domain and subdomains same time
	res, err := m.createCertificateForDomains(certIssueContext, certName, domainNamesFromCertificateName(certName), needDomain)
	if err == nil {
		logger.Info("Certificate issued.", log.Cert(res),
			zap.Time("expire", res.Leaf.NotAfter))
		return res, nil
	} else {
		logger.Warn("Can't issue certificate", zap.Error(err))
		return nil, haveNoCert
	}
}

func (m *Manager) certStateGet(certName certNameType) *certState {
	m.certStateMu.Lock()
	defer m.certStateMu.Unlock()

	res := m.certState[certName]
	if res == nil {
		res = &certState{}
		m.certState[certName] = res
	}
	return res
}

func (m *Manager) createCertificateForDomains(ctx context.Context, certName certNameType, domainNames []DomainName, needDomain DomainName) (res *tls.Certificate, err error) {
	logger := zc.L(ctx).With(logCetName(certName))
	ctx = zc.WithLogger(ctx, logger)

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

	cachedCert, err := getCertificate(ctx, m.Cache, certName, keyRSA)
	if err == nil {
		logger.Debug("Certificate loaded from cache")
		return cachedCert, nil
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

		logger.Info("Start detached process for revoke pending authorizations", zap.Strings("uries", authUries))

		revokeContext := context.Background()
		revokeContext = zc.WithLogger(revokeContext, logger.With(zap.Bool("detached", true)))
		go m.revokePendingAuthorizations(revokeContext, authUries)
	}()

	var allowedChallenges []string
	if m.EnableTlsValidation {
		allowedChallenges = append(allowedChallenges, tlsAlpn01)
	}
	if m.EnableHttpValidation {
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
		} else {
			logger.Debug("Select challenge for authorize", zap.Reflect("challenge", chal))
		}

		cleanup, err := m.fulfill(ctx, chal, domain)
		if err != nil {
			logger.Error("Can't set domain token", zap.Reflect("chal", chal), zap.Error(err))
			continue
		}
		//noinspection GoDeferInLoop
		defer cleanup(ctx)

		receivedChallenge, err := m.Client.Accept(ctx, chal)
		if err == nil {
			logger.Debug("Receive authorize answer", zap.Reflect("challenge", receivedChallenge))
		} else {
			logger.Error("Can't authorize domain", zap.Reflect("challenge", receivedChallenge), zap.Error(err))
			continue
		}

		receivedAuth, err := m.Client.WaitAuthorization(ctx, authz.URI)
		if err == nil {
			logger.Debug("Receive domain authorization", zap.Reflect("authorization", receivedAuth))
		} else {
			logger.Error("Dont receive wait authorization", zap.Reflect("authorizarion", receivedAuth), zap.Error(err))
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
	der, _, err := m.Client.CreateCert(ctx, csr, 0, true)

	if err != nil {
		logger.Error("Can't issue certificate", zap.Error(err))
		return nil, err
	}

	cert, err := validCertDer(domains, der, key, time.Now())

	if err == nil {
		logger.Info("Certificated issued")
		storeCertificate(ctx, m.Cache, certName, cert)
		return cert, nil
	} else {
		logger.Error("Receive invalid certificate", zap.Error(err))
		return nil, err
	}
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
	} else {
		if err == cache.ErrCacheMiss {
			logger.Debug("Cert key no in cache. Create new.")
		} else {
			logger.Error("Error while check cert key in cache", zap.Error(err))
			return nil, err
		}
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
		return func(localContext context.Context) { _ = m.httpTokens.Delete(localContext, key) }, nil
	default:
		logger.Error("Unknow challenge type", zap.Reflect("challenge", challenge))
		return nil, errors.New("unknown challenge type")
	}
}

func (m *Manager) IsHttpValidationRequest(r *http.Request) bool {
	if r == nil || r.URL == nil {
		return false
	}
	if r.Method != http.MethodGet {
		return false
	}

	return strings.HasPrefix(r.URL.Path, httpWellKnown)
}

func (m *Manager) HandleHttpValidation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := zc.L(ctx)
	if !m.IsHttpValidationRequest(r) {
		logger.DPanic("Pass non validation url to handle in cert manager.", zap.Reflect("req", r))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token := strings.TrimPrefix(r.URL.Path, httpWellKnown)
	domain := normalizeDomain(r.Host)
	resp, err := m.httpTokens.Get(ctx, domain.ASCII()+"/"+token)
	if err == nil {
		_, err = w.Write(resp)
		logger.Warn("Error write http token answer to response", logDomain(domain), zap.String("token", token), zap.Error(err))
	} else {
		logger.Warn("Have no validation token", logDomain(domain), zap.String("token", token), zap.Error(err))
	}
}

func (m *Manager) putCertToken(ctx context.Context, key DomainName, certificate *tls.Certificate) {
	logger := zc.L(ctx)
	logger.Debug("Put cert token", zap.String("key", string(key)))

	var overwriteToken bool

	defer func() {
		if overwriteToken {
			logger.Warn("Unexpected cert token already exist", zap.String("key", key.String()))
		}
	}()

	m.certTokensMu.Lock()
	defer m.certTokensMu.Unlock()

	if m.certTokens == nil {
		m.certTokens = make(map[DomainName]*tls.Certificate)
	}

	_, overwriteToken = m.certTokens[key]
	m.certTokens[key] = certificate
}

func (m *Manager) deleteCertToken(ctx context.Context, key DomainName) {
	logger := zc.L(ctx)
	logger.Debug("Delete cert token", zap.String("key", key.String()))

	var exist bool
	defer func() {
		if !exist {
			logger.Warn("Cert token for delete doesn't exist", zap.String("key", string(key)))
		}
	}()

	m.certTokensMu.Lock()
	defer m.certTokensMu.Unlock()

	_, exist = m.certTokens[key]
	delete(m.certTokens, key)
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
		logger.Debug("Check certificate in cache", log.Cert(cert), zap.Error(err))
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
	return validCertTls(&cert2, nil, time.Now())
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
			return fail(fmt.Errorf("tls: failed to find PEM block with type ending in \"PRIVATE KEY\" in key input after skipping PEM blocks of the following types: %v", skippedBlockTypes))
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
