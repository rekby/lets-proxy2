package manager

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"reflect"
	"sync"
	"time"

	"github.com/rekby/lets-proxy2/internal/log"

	"github.com/rekby/zapcontext"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme"
)

const (
	tlsAlpn01 = "tls-alpn-01"
)

var (
	allowedChallenges  = []string{tlsAlpn01}
	domainKeyRSALength = 2048
)

var ErrCacheMiss = errors.New("lets proxy: certificate cache miss")
var haveNoCert = errors.New("have no certificate for domain")
var notImplementedError = errors.New("not implemented yet")

type Cache interface {
	// Get returns a certificate data for the specified key.
	// If there's no such key, Get returns ErrCacheMiss.
	Get(ctx context.Context, key string) ([]byte, error)

	// Put stores the data in the cache under the specified key.
	// Underlying implementations may use any data storage format,
	// as long as the reverse operation, Get, results in the original data.
	Put(ctx context.Context, key string, data []byte) error

	// Delete removes a certificate data from the cache under the specified key.
	// If there's no such key in the cache, Delete returns nil.
	Delete(ctx context.Context, key string) error
}

type GetContext interface {
	GetContext() context.Context
}

// Interface inspired to https://godoc.org/golang.org/x/crypto/acme/autocert#Manager but not compatible gurantee
type Manager struct {
	CertificateIssueTimeout time.Duration
	Cache                   Cache

	// Client is used to perform low-level operations, such as account registration
	// and requesting new certificates.
	//
	// If Client is nil, a zero-value acme.Client is used with acme.LetsEncryptURL
	// as directory endpoint. If the Client.Key is nil, a new ECDSA P-256 key is
	// generated and, if Cache is not nil, stored in cache.
	//
	// Mutating the field after the first call of GetCertificate method will have no effect.
	Client *acme.Client

	certTokensMu sync.RWMutex
	certTokens   map[DomainName]*tls.Certificate

	certStateMu sync.Mutex
	certState   map[certNameType]*certState
}

func New(ctx context.Context, client *acme.Client) *Manager {
	res := Manager{}
	res.Client = client
	res.certTokens = make(map[DomainName]*tls.Certificate)
	res.certState = make(map[certNameType]*certState)
	res.CertificateIssueTimeout = time.Minute
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
	logger.Info("Get certificate")
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
		cert, err = validCertDer([]DomainName{needDomain}, cert.Certificate, cert.PrivateKey, time.Now())
		if cert != nil {
			logger.Debug("Got certificate from local state", log.Cert(cert))
			return cert, nil
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

	cachedCert, err := getCertificate(ctx, m.Cache, certName, "rsa")
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

	key, err := m.certKeyGet(ctx, certName)
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

var generatedKey *rsa.PrivateKey

func (m *Manager) certKeyGet(ctx context.Context, domain certNameType) (crypto.Signer, error) {
	logger := zc.L(ctx)

	//TODO: load/save with cache
	if generatedKey != nil {
		logger.Debug("Get RSA key from cache.")
		return generatedKey, nil
	}

	logger.Debug("Generate new rsa key")
	key, err := rsa.GenerateKey(rand.Reader, domainKeyRSALength)
	if err == nil {
		generatedKey = key
		return key, nil
	}

	logger.Error("Can't generate rsa key", zap.Error(err))
	return nil, err
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
		return func(ctx context.Context) { go m.deleteCertToken(ctx, domain) }, nil
	default:
		logger.Error("Unknow challenge type", zap.Reflect("challenge", challenge))
		return nil, errors.New("unknown challenge type")
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
func storeCertificate(ctx context.Context, cache Cache, certName certNameType,
	cert *tls.Certificate) {
	logger := zc.L(ctx)
	if cache == nil {
		logger.Debug("Can't save certificate to nil cache")
		return
	}

	var keyType string

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
		keyType = "rsa"
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

	certKeyName := string(certName) + "." + keyType + ".cer"
	keyKeyName := string(certName) + "." + keyType + ".key"

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

func getCertificate(ctx context.Context, cache Cache, certName certNameType, keyType string) (cert *tls.Certificate, err error) {
	logger := zc.L(ctx)
	logger.Debug("Check certificate in cache")
	defer func() {
		logger.Debug("Check certificate in cache", log.Cert(cert), zap.Error(err))
	}()

	certKeyName := string(certName) + "." + keyType + ".cer"
	keyKeyName := string(certName) + "." + keyType + ".key"

	certBytes, err := cache.Get(ctx, certKeyName)
	if err != nil {
		return nil, err
	}

	keyBytes, err := cache.Get(ctx, keyKeyName)
	if err != nil {
		return nil, err
	}

	cert2, err := tls.X509KeyPair(certBytes, keyBytes)
	if err != nil {
		return nil, err
	}
	return validCertTls(&cert2, nil, time.Now())
}
