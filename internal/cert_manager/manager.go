//nolint:golint
package cert_manager

import (
	"bytes"
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/rekby/lets-proxy2/internal/contexthelper"

	"github.com/rekby/lets-proxy2/internal/metrics"

	"github.com/prometheus/client_golang/prometheus"

	"golang.org/x/xerrors"

	"github.com/rekby/lets-proxy2/internal/cache"
	"go.uber.org/zap/zapcore"

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
const renewBeforeExpire = time.Hour * 24 * 30
const revokeAuthorizationTimeout = 5 * time.Minute
const cleanupTimeout = time.Minute

var errHaveNoCert = errors.New("have no certificate for domain") // may return for any internal error
var errRSADenied = xerrors.New("RSA certificate denied by config")
var errECDSADenied = xerrors.New("ECDSA certificate denied by config")
var errCertTypeUnknown = xerrors.New("unknown cert type")

type GetContext interface {
	GetContext() context.Context
}

type KeyType string

const KeyRSA KeyType = "rsa"
const KeyECDSA KeyType = "ecdsa"

func (t KeyType) Generate() (crypto.Signer, error) {
	switch t {
	case KeyRSA:
		return rsa.GenerateKey(rand.Reader, domainKeyRSALength)
	case KeyECDSA:
		return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	default:
		panic("Unexpected key type for generate: " + t.String())
	}
}

func (t KeyType) String() string {
	return string(t)
}

// Interface inspired to https://godoc.org/golang.org/x/crypto/acme/autocert#Manager but not compatible guarantee
type Manager struct {
	CertificateIssueTimeout time.Duration
	Cache                   cache.Bytes

	// Subdomains, auto-issued with main domain.
	// Every subdomain must have suffix dot. For example: "www."
	AutoSubdomains []string

	// Client is used to perform low-level operations, such as account registration
	// and requesting new certificates.
	//
	// If Client is nil, a zero-value acme.Client is used with acme.LetsEncryptURL
	// as directory endpoint. If the Client.Key is nil, a new ECDSA P-256 key is
	// generated and, if Cache is not nil, stored in cache.
	//
	// Mutating the field after the first call of GetCertificate method will have no effect.
	Client               AcmeClient
	DomainChecker        DomainChecker
	EnableHTTPValidation bool
	EnableTLSValidation  bool
	SaveJSONMeta         bool
	AllowECDSACert       bool
	AllowRSACert         bool

	certForDomainAuthorize cache.Value

	certStateMu sync.Mutex
	certState   cache.Value

	httpTokens cache.Bytes

	// metrics
	handleCertStart, certRequestStart   metrics.ProcessStartFunc
	handleCertFinish, certRequestFinish metrics.ProcessFinishFunc
}

func New(client AcmeClient, c cache.Bytes, r prometheus.Registerer) *Manager {
	res := Manager{}
	res.Client = client
	res.certForDomainAuthorize = cache.NewMemoryValueLRU("authcert")
	res.certState = cache.NewMemoryValueLRU("certstate")
	res.CertificateIssueTimeout = time.Minute
	res.httpTokens = cache.NewMemoryCache("Http validation tokens")
	res.Cache = c
	res.EnableTLSValidation = true
	res.DomainChecker = managerDefaults{}
	res.AllowRSACert = true
	res.AllowECDSACert = true

	res.initMetrics(r)
	return &res
}

// GetCertificate implements the tls.Config.GetCertificate hook.
func (m *Manager) GetCertificate(hello *tls.ClientHelloInfo) (resultCert *tls.Certificate, err error) {
	m.handleCertStart()
	defer func() {
		m.handleCertFinish(err)
	}()
	ctx := hello.Conn.(GetContext).GetContext()

	logger := zc.L(ctx)

	needDomain, err := normalizeDomain(hello.ServerName)
	log.DebugInfo(logger, err, "Domain name normalization", zap.String("original", hello.ServerName), logDomain(needDomain))
	if err != nil {
		return nil, errHaveNoCert
	}

	logger = logger.With(logDomain(needDomain))
	ctx = zc.WithLogger(ctx, logger)
	defer log.HandlePanic(logger)

	logger.Info("Get certificate", zap.String("original_domain", hello.ServerName))
	if isTLSALPN01Hello(hello) {
		return m.handleTLSALPN(ctx, needDomain)
	}

	certType := KeyRSA
	if supportsECDSA(hello) {
		certType = KeyECDSA
	}
	cert, err := m.getCertificate(ctx, needDomain, certType)
	log.DebugInfo(logger, err, "Got certificate", log.Cert(cert))
	if err == nil || certType == KeyRSA {
		return cert, err
	}

	logger.Info("ECDSA certificate was failed, try to get RSA certificate")
	ctx = zc.WithLogger(ctx, logger.With(zap.String("retry_type", "rsa")))
	return m.getCertificate(ctx, needDomain, KeyRSA)
}

//nolint:funlen,gocognit
func (m *Manager) getCertificate(ctx context.Context, needDomain DomainName, certType KeyType) (resultCert *tls.Certificate, err error) {
	switch certType {
	case KeyRSA:
		if !m.AllowRSACert {
			return nil, errRSADenied
		}
	case KeyECDSA:
		if !m.AllowECDSACert {
			return nil, errECDSADenied
		}
	default:
		return nil, errCertTypeUnknown
	}

	certDescription := CertDescriptionFromDomain(needDomain, certType, m.AutoSubdomains)

	logger := zc.L(ctx).With(certDescription.ZapField())
	ctx = zc.WithLogger(ctx, zc.L(ctx).With(certDescription.ZapField()))

	now := time.Now()

	var locked = false
	var lockedChecked = false

	defer func() {
		if isNeedRenew(resultCert, now) {
			if !lockedChecked {
				locked, err = isCertLocked(ctx, m.Cache, certDescription)
				log.DebugError(logger, err, "Check locked before renew", zap.Bool("locked", locked))
				if err != nil {
					return
				}
			}
			if !locked {
				// handlepanic: in renewCertInBackground
				go m.renewCertInBackground(ctx, needDomain, certDescription)
			}
		}
	}()

	certState := m.certStateGet(ctx, certDescription)
	cert, err := certState.Cert()
	if cert != nil {
		logger.Debug("Got certificate from local state", log.Cert(cert))

		cert, err = validCertTLS(cert, []DomainName{needDomain}, certState.GetUseAsIs(), now)
		logger.Debug("Validate certificate from local state", zap.Error(err))
		if err == nil {
			return cert, nil
		}
	}
	if err != nil {
		logLevel := zapcore.ErrorLevel
		if err == cache.ErrCacheMiss {
			logLevel = zapcore.DebugLevel
		}
		log.LevelParam(logger, logLevel, "Can't get certificate from local state", zap.Error(err))
		if logLevel == zapcore.ErrorLevel {
			return nil, errHaveNoCert
		}
	}

	locked, err = isCertLocked(ctx, m.Cache, certDescription)
	lockedChecked = true
	log.DebugDPanic(logger, err, "Check if certificate locked", zap.Bool("locked", locked))
	if err != nil {
		return nil, errHaveNoCert
	}

	cert, err = loadCertificateFromCache(ctx, m.Cache, certDescription)
	logLevel := zapcore.ErrorLevel
	if err == nil || err == cache.ErrCacheMiss {
		logLevel = zapcore.DebugLevel
	}
	log.LevelParam(logger, logLevel, "Load certificate from cache", zap.Error(err))

	if err == nil {
		cert, err = validCertDer([]DomainName{needDomain}, cert.Certificate, cert.PrivateKey, locked, now)
		logger.Debug("Check if certificate ok", zap.Error(err))
		if err == nil {
			certState.CertSet(ctx, locked, cert)
			return cert, nil
		}
	}
	if err != cache.ErrCacheMiss && err != errCertExpired {
		return nil, errHaveNoCert
	}

	if locked {
		return nil, errHaveNoCert
	}

	return m.issueNewCert(ctx, needDomain, certDescription)
}

func (m *Manager) issueNewCert(ctx context.Context, needDomain DomainName, cd CertDescription) (cert *tls.Certificate, err error) {
	m.certRequestStart()
	defer func() {
		m.certRequestFinish(err)
	}()
	logger := zc.L(ctx)

	allowed, err := m.DomainChecker.IsDomainAllowed(ctx, needDomain.ASCII())
	log.DebugError(logger, err, "Check if domain allowed for certificate", zap.Bool("allowed", allowed))
	if err != nil {
		return nil, errHaveNoCert
	}
	if !allowed {
		logger.Info("Deny certificate issue by filter")
		return nil, errHaveNoCert
	}
	certIssueContext, cancelFunc := context.WithTimeout(ctx, m.CertificateIssueTimeout)
	defer cancelFunc()

	domains := cd.DomainNames()
	domains, err = filterDomains(ctx, m.DomainChecker, domains, needDomain)
	log.DebugError(logger, err, "Filter domains", logDomains(domains))

	res, err := m.createCertificateForDomains(certIssueContext, cd, domains)
	if err == nil {
		logger.Info("Certificate issued.", log.Cert(res),
			zap.Time("expire", res.Leaf.NotAfter))
		return res, nil
	}
	logger.Warn("Can't issue certificate", zap.Error(err))
	return nil, errHaveNoCert
}

func (m *Manager) handleTLSALPN(ctx context.Context, needDomain DomainName) (*tls.Certificate, error) {
	logger := zc.L(ctx)
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

func filterDomains(ctx context.Context, checker DomainChecker, originalDomains []DomainName, needDomain DomainName) ([]DomainName, error) {
	logger := zc.L(ctx)
	logger.Debug("filter domains from certificate list", logDomains(originalDomains))
	var allowedDomains = make(chan DomainName, len(originalDomains))
	var hasNeedDomain bool
	var wg sync.WaitGroup

	wg.Add(len(originalDomains))

	for _, domain := range originalDomains {
		domain := domain // pin var

		go func() {
			defer wg.Done()
			defer log.HandlePanic(logger)

			allow, err := checker.IsDomainAllowed(ctx, domain.ASCII())
			logger.Debug("Check domain", logDomain(domain), zap.Bool("allowed", allow), zap.Error(err))
			if !allow {
				return
			}

			if domain == needDomain {
				hasNeedDomain = true
			}

			allowedDomains <- domain
		}()
	}

	wg.Wait()
	close(allowedDomains)

	if !hasNeedDomain {
		return nil, errors.New("need domain doesn't contained to cert list domains after filter")
	}

	res := make([]DomainName, 0, len(allowedDomains))
	for domain := range allowedDomains {
		res = append(res, domain)
	}
	return res, nil
}

func (m *Manager) certStateGet(ctx context.Context, cd CertDescription) *certState {
	m.certStateMu.Lock()
	defer m.certStateMu.Unlock()

	resInterface, err := m.certState.Get(ctx, cd.String())
	if err == cache.ErrCacheMiss {
		err = nil
	}
	log.DebugFatalCtx(ctx, err, "Got cert state from cache", zap.Bool("is_empty", resInterface == nil))
	if resInterface == nil {
		resInterface = &certState{}
		err = m.certState.Put(ctx, cd.String(), resInterface)
		log.DebugFatalCtx(ctx, err, "Put empty cert state to cache")
	}
	return resInterface.(*certState)
}

func (m *Manager) createCertificateForDomains(ctx context.Context, cd CertDescription, domainNames []DomainName) (res *tls.Certificate, err error) {
	logger := zc.L(ctx).With(logDomains(domainNames))
	certState := m.certStateGet(ctx, cd)

	if !certState.StartIssue(ctx) {
		waitTimeout, waitTimeoutCancel := context.WithTimeout(ctx, m.CertificateIssueTimeout)
		defer waitTimeoutCancel()

		logger.Debug("Certificate issue in process already - wait result")
		return certState.WaitFinishIssue(waitTimeout)
	}
	// outer func need for get argument values in defer time
	defer func() {
		certState.FinishIssue(ctx, res, err)
	}()

	logger.Debug("Start issue process")

	order, err := m.createOrderForDomains(ctx, domainNames...)
	log.DebugWarning(logger, err, "Domains authorized")
	if err != nil {
		return nil, errors.New("order authorization error")
	}

	res, err = m.issueCertificate(ctx, cd, order)
	log.DebugWarning(logger, err, "Issue certificate")
	return res, err
}

func (m *Manager) supportedChallenges() []string {
	var allowedChallenges []string
	if m.EnableTLSValidation {
		allowedChallenges = append(allowedChallenges, tlsAlpn01)
	}
	if m.EnableHTTPValidation {
		allowedChallenges = append(allowedChallenges, http01)
	}
	return allowedChallenges
}

// createOrderForDomains similar to func (m *Manager) verifyRFC(ctx context.Context, client *acme.Client, domain string) (*acme.Order, error)
// from acme/autocert
//nolint:funlen,gocognit
func (m *Manager) createOrderForDomains(ctx context.Context, domains ...DomainName) (*acme.Order, error) {
	client := m.Client
	logger := zc.L(ctx)
	challengeTypes := m.supportedChallenges()
	logger.Debug("Start order authorization.")
	var order *acme.Order

	firstLoop := true
authorizeOrderLoop:
	for {
		if ctx.Err() != nil {
			return nil, xerrors.Errorf("context canceled: %w", ctx.Err())
		}

		if firstLoop {
			firstLoop = false
		} else {
			time.Sleep(time.Second)
		}
		authIDs := make([]acme.AuthzID, len(domains))
		for i := range domains {
			authIDs[i] = acme.AuthzID{Type: "dns", Value: domains[i].ASCII()}
		}
		var err error
		order, err = m.Client.AuthorizeOrder(ctx, authIDs)
		log.DebugError(logger, err, "Create authorization order.", zap.Reflect("order", order))
		if err != nil {
			return nil, err
		}

		//noinspection GoDeferInLoop
		defer func(order *acme.Order) {
			go func(order *acme.Order) {
				defer log.HandlePanic(logger)

				revokeLogger := logger.Named("background_auth_revoker")

				revokeCtx, cancel := context.WithTimeout(context.Background(), revokeAuthorizationTimeout)
				defer cancel()

				revokeCtx = zc.WithLogger(revokeCtx, revokeLogger)
				m.deactivatePendingAuthz(revokeCtx, order.AuthzURLs)
			}(order)
		}(order)

		switch order.Status {
		case acme.StatusReady:
			logger.Debug("Authorization completed")
			break authorizeOrderLoop

		case acme.StatusPending:
		// pass
		default:
			logger.Error("Invalid new order status", zap.String("status", order.Status), zap.String("uri", order.URI))
			return nil, errors.New("invalid new order status")
		}

		logger.Debug("Start authorization step")

		// Satisfy all pending authorizations.
	authDomainLoop:
		for _, zurl := range order.AuthzURLs {
			z, err := client.GetAuthorization(ctx, zurl)

			log.DebugError(logger, err, "Get authorization object.", zap.Reflect("authorization", z))
			if err != nil {
				return nil, err
			}
			// force hide outer logger - for log specific domain on each iteration
			var logger = logger.With(logDomain(DomainName(z.Identifier.Value)))

			hasCompatibleChallenge := false
		challengeTypeLoop:
			for _, challengeType := range challengeTypes {
				if z.Status != acme.StatusPending {
					// We are interested only in pending authorizations.
					logger.Debug("Skip authorize process", zap.String("status", z.Status))
					continue authDomainLoop
				}

				// Pick the next preferred challenge.
				var chal = pickChallenge(challengeType, z.Challenges)
				logger.Debug("Selected challenge", zap.Any("challenge", chal))
				if chal == nil {
					continue challengeTypeLoop
				}
				hasCompatibleChallenge = true

				// Respond to the challenge and wait for validation result.
				cleanup, err := m.fulfill(ctx, chal, DomainName(z.Identifier.Value))
				log.DebugError(logger, err, "Write respond to challenge")
				if err != nil {
					continue authorizeOrderLoop
				}
				cleanupContext, cleanupContextCancel := context.WithTimeout(contexthelper.DropCancelContext(ctx), cleanupTimeout)
				//noinspection GoDeferInLoop
				defer cleanupContextCancel()
				//noinspection GoDeferInLoop
				defer cleanup(cleanupContext)

				authorizedChallenge, err := client.Accept(ctx, chal)
				log.DebugError(logger, err, "accept authorization", zap.Reflect("authorized_challenge", authorizedChallenge))
				if err != nil {
					continue authorizeOrderLoop
				}
				authorization, err := client.WaitAuthorization(ctx, z.URI)
				log.DebugError(logger, err, "wait authorization", zap.Reflect("authorization", authorization))
				if err != nil {
					continue authorizeOrderLoop
				}
			}
			if !hasCompatibleChallenge {
				logger.Error("No compatible challenges")
				return nil, fmt.Errorf("unable to satisfy %q for domain %q: no viable challenge type found", z.URI, z.Identifier.Value)
			}
		}

		// All authorizations are satisfied.
		// Wait for the CA to update the order status.
		order, err = client.WaitOrder(ctx, order.URI)
		log.DebugWarning(logger, err, "Wait order authorization.", zap.Reflect("order", order))
		if err == nil {
			break authorizeOrderLoop
		}
	}
	return order, nil
}

func (m *Manager) issueCertificate(ctx context.Context, cd CertDescription, order *acme.Order) (*tls.Certificate, error) {
	if len(order.Identifiers) == 0 {
		return nil, errors.New("no domains for issue certificate")
	}

	domains := make([]DomainName, len(order.Identifiers))
	for i := range order.Identifiers {
		domains[i] = DomainName(order.Identifiers[i].Value)
	}
	logger := zc.L(ctx).With(logDomains(domains))

	key, err := m.certKeyGetOrCreate(ctx, cd)
	log.DebugError(logger, err, "Get cert key")
	if err != nil {
		return nil, err
	}

	csr, err := createCertRequest(key, domains[0], domains...)
	log.DebugDPanic(logger, err, "Create certificate request")
	if err != nil {
		return nil, err
	}

	der, _, err := m.Client.CreateOrderCert(ctx, order.FinalizeURL, csr, true)
	log.InfoError(logger, err, "Receive certificate from acme server")
	if err != nil {
		return nil, err
	}

	cert, err := validCertDer(domains, der, key, false, time.Now())
	log.DebugDPanic(logger, err, "Check certificate is valid")
	if err != nil {
		return nil, err
	}

	err = storeCertificate(ctx, m.Cache, cd, cert)
	log.DebugDPanic(logger, err, "Certificate stored")
	if err != nil {
		return nil, err
	}
	if m.SaveJSONMeta {
		err = storeCertificateMeta(ctx, m.Cache, cd, cert)
		if err != nil {
			return nil, err
		}
	}
	return cert, nil
}

func (m *Manager) renewCertInBackground(ctx context.Context, needDomain DomainName, cd CertDescription) {
	// detach from request lifetime, but save log context
	logger := zc.L(ctx).Named("background")
	defer log.HandlePanic(logger)
	ctx, ctxCancel := context.WithTimeout(context.Background(), m.CertificateIssueTimeout)
	defer ctxCancel()

	logger.Debug("Start reissue certificate in background")

	ctx = zc.WithLogger(ctx, logger)
	_, err := m.issueNewCert(ctx, needDomain, cd)
	log.DebugError(logger, err, "Cert reissue in background finished")
}

func (m *Manager) deactivatePendingAuthz(ctx context.Context, uries []string) {
	logger := zc.L(ctx)

	for _, uri := range uries {
		localLogger := logger.With(zap.String("uri", uri))
		authorization, err := m.Client.GetAuthorization(ctx, uri)
		log.DebugError(localLogger, err, "Get authorization", zap.Reflect("authorization", authorization))
		if err != nil {
			continue
		}
		if authorization.Status == acme.StatusPending {
			err := m.Client.RevokeAuthorization(ctx, uri)
			log.DebugInfo(localLogger, err, "Revoke authorization", zap.String("uri", uri))
		} else {
			localLogger.Debug("Authorization not in pending state. Skip revoke.", zap.String("status", authorization.Status))
		}
	}
}

func (m *Manager) certKeyGetOrCreate(ctx context.Context, cd CertDescription) (crypto.Signer, error) {
	logger := zc.L(ctx)

	key, err := getCertificateKey(ctx, m.Cache, cd)
	logger.Debug("Got certificate key from cache and reuse old key", zap.Error(err))
	if err == nil {
		return key, nil
	}
	if err != cache.ErrCacheMiss {
		return nil, err
	}

	key, err = cd.KeyType.Generate()
	log.InfoError(logger, err, "Generate new key")
	return key, err
}

func (m *Manager) fulfill(ctx context.Context, challenge *acme.Challenge, domain DomainName) (func(context.Context), error) {
	logger := zc.L(ctx)

	switch challenge.Type {
	case tlsAlpn01:
		cert, err := m.Client.TLSALPN01ChallengeCert(challenge.Token, domain.String())
		log.DebugError(logger, err, "Got TLSALPN01ChallengeCert", log.Cert(&cert))
		if err != nil {
			return nil, err
		}
		m.putCertToken(ctx, domain, &cert)
		return func(localContext context.Context) {
			// handlepanic: in deleteCertToken
			go m.deleteCertToken(localContext, domain)
		}, nil
	case http01:
		resp, err := m.Client.HTTP01ChallengeResponse(challenge.Token)
		if err != nil {
			return nil, err
		}
		key := domain.ASCII() + "/" + challenge.Token
		err = m.httpTokens.Put(ctx, key, []byte(resp))
		log.DebugError(logger, err, "Put token for http-01", zap.String("key", key))
		if err == nil {
			return func(localContext context.Context) { _ = m.httpTokens.Delete(localContext, key) }, nil
		} else {
			return nil, err
		}
	default:
		logger.Error("Unknow challenge type", zap.Reflect("challenge", challenge))
		return nil, errors.New("unknown challenge type")
	}
}

func (m *Manager) initMetrics(r prometheus.Registerer) {
	m.handleCertStart, m.handleCertFinish = metrics.ToefCounters(r, "handle_cert", "handled certificates")
	m.certRequestStart, m.certRequestFinish = metrics.ToefCounters(r, "cert_request", "request certificates from lets-encrypt")
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

func (m *Manager) HandleHTTPValidation(w http.ResponseWriter, r *http.Request) bool {
	if !m.isHTTPValidationRequest(r) {
		return false
	}

	ctx := r.Context()
	logger := zc.L(ctx)

	token := strings.TrimPrefix(r.URL.Path, httpWellKnown)
	domain, err := normalizeDomain(r.Host)
	log.DebugInfo(logger, err, "Domain normalization", zap.String("original", r.Host), logDomain(domain))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return true
	}
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
	defer log.HandlePanicCtx(ctx)

	err := m.certForDomainAuthorize.Delete(ctx, key.String())
	log.DebugDPanicCtx(ctx, err, "Delete cert token", zap.String("key", key.String()))
}

// It isn't atomic syncronized - caller must not save two certificates with same name same time
func storeCertificate(ctx context.Context, cache cache.Bytes, cd CertDescription,
	cert *tls.Certificate) error {
	logger := zc.L(ctx)

	locked, err := isCertLocked(ctx, cache, cd)
	log.DebugError(logger, err, "Check if cert locked", zap.Bool("key_locked", locked))
	if locked {
		logger.DPanic("Logical error - try to save to locked certificate")
		return xerrors.New("Try save to locked certificate")
	}

	var keyType = getKeyType(cert)

	var certBuf bytes.Buffer

	for _, block := range cert.Certificate {
		err := pem.Encode(&certBuf, &pem.Block{Type: "CERTIFICATE", Bytes: block})
		if err != nil {
			logger.DPanic("Can't encode pem block of certificate", zap.Error(err), zap.Binary("block", block))
			return err
		}
	}

	var privateKeyBytes []byte

	switch keyType {
	case KeyRSA:
		privateKey := cert.PrivateKey.(*rsa.PrivateKey)
		keyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
		privateKeyBytes = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyBytes})
	case KeyECDSA:
		privateKey := cert.PrivateKey.(*ecdsa.PrivateKey)
		keyBytes, err := x509.MarshalECPrivateKey(privateKey)
		if err != nil {
			return err
		}
		privateKeyBytes = pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes})
	default:
		logger.DPanic("Unknow private key type", zap.String("type", reflect.TypeOf(cert.PrivateKey).String()))
		return errors.New("unknow private key type")
	}

	certKeyName := cd.CertStoreName()
	keyKeyName := cd.KeyStoreName()

	err = cache.Put(ctx, certKeyName, certBuf.Bytes())
	zc.InfoError(logger, err, "Store certificate file", zap.String("cert_key", certKeyName))
	if err != nil {
		return err
	}

	err = cache.Put(ctx, keyKeyName, privateKeyBytes)
	zc.InfoError(logger, err, "Store key file", zap.String("key_key", keyKeyName))
	if err != nil {
		_ = cache.Delete(ctx, certKeyName)
		return err
	}
	return nil
}

func storeCertificateMeta(ctx context.Context, storage cache.Bytes, cd CertDescription, certificate *tls.Certificate) error {
	info := struct {
		Domains    []string
		ExpireDate time.Time
	}{
		Domains:    certificate.Leaf.DNSNames,
		ExpireDate: certificate.Leaf.NotAfter,
	}
	infoBytes, _ := json.MarshalIndent(info, "", "    ")
	err := storage.Put(ctx, cd.MetaStoreName(), infoBytes)
	log.DebugDPanicCtx(ctx, err, "Save cert metadata")
	return err
}

func getKeyType(cert *tls.Certificate) KeyType {
	if cert == nil {
		panic("cert is nil")
	}
	switch cert.PrivateKey.(type) {
	case *rsa.PrivateKey:
		return KeyRSA
	case *ecdsa.PrivateKey:
		return KeyECDSA
	default:
		panic("unexexpected key type: " + reflect.TypeOf(cert.PrivateKey).PkgPath() + "." + reflect.TypeOf(cert.PrivateKey).Name())
	}
}

func loadCertificateFromCache(ctx context.Context, c cache.Bytes, cd CertDescription) (cert *tls.Certificate, err error) {
	logger := zc.L(ctx)
	logger.Debug("Check certificate in cache")
	defer func() {
		logger.Debug("Checked certificate in cache", log.Cert(cert), zap.Error(err))
	}()

	certCerName := cd.CertStoreName()

	certBytes, err := c.Get(ctx, certCerName)
	logLevel := zapcore.ErrorLevel
	if err == nil || err == cache.ErrCacheMiss {
		logLevel = zapcore.DebugLevel
	}
	log.LevelParam(logger, logLevel, "Get certificate from cache", zap.Error(err))

	if err != nil {
		return nil, err
	}

	keyBytes, err := getCertificateKeyBytes(ctx, c, cd)
	logLevel = zapcore.ErrorLevel
	if err == nil || err == cache.ErrCacheMiss {
		logLevel = zapcore.DebugLevel
	}
	log.LevelParam(logger, logLevel, "Get certificate key from cache")
	if err != nil {
		return nil, err
	}

	cert2, err := tls.X509KeyPair(certBytes, keyBytes)
	log.DebugError(logger, err, "Combine cert and key into pair")
	if err != nil {
		// logical error, may be system failure
		return nil, err
	}
	if len(cert2.Certificate) > 0 {
		cert2.Leaf, err = x509.ParseCertificate(cert2.Certificate[0])
		if err != nil {
			// logical error, may be system failure
			return nil, err
		}
	}
	locked, err := isCertLocked(ctx, c, cd)
	log.DebugError(logger, err, "Check if certificate locked")
	if err != nil {
		// logical error, may be system failure
		return nil, err
	}
	return validCertTLS(&cert2, nil, locked, time.Now())
}

func getCertificateKeyBytes(ctx context.Context, cache cache.Bytes, cd CertDescription) ([]byte, error) {
	keyKeyName := cd.KeyStoreName()
	return cache.Get(ctx, keyKeyName)
}

func getCertificateKey(ctx context.Context, cache cache.Bytes, cd CertDescription) (crypto.Signer, error) {
	certBytes, err := getCertificateKeyBytes(ctx, cache, cd)
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
	return cert.Leaf.NotAfter.Add(-renewBeforeExpire).Before(now)
}

func isCertLocked(ctx context.Context, storage cache.Bytes, certName CertDescription) (bool, error) {
	lockName := certName.LockName()

	_, err := storage.Get(ctx, lockName)
	switch err {
	case cache.ErrCacheMiss:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, err
	}
}

// copy from golang.org/x/crypto/acme/autocert/autocert.go
// https://github.com/golang/crypto/blob/87dc89f01550277dc22b74ffcf4cd89fa2f40f4c/acme/autocert/autocert.go#L322
func supportsECDSA(hello *tls.ClientHelloInfo) bool {
	// The "signature_algorithms" extension, if present, limits the key exchange
	// algorithms allowed by the cipher suites. See RFC 5246, section 7.4.1.4.1.
	if hello.SignatureSchemes != nil {
		ecdsaOK := false
	schemeLoop:
		for _, scheme := range hello.SignatureSchemes {
			const tlsECDSAWithSHA1 tls.SignatureScheme = 0x0203 // constant added in Go 1.10
			switch scheme {
			case tlsECDSAWithSHA1, tls.ECDSAWithP256AndSHA256,
				tls.ECDSAWithP384AndSHA384, tls.ECDSAWithP521AndSHA512:
				ecdsaOK = true
				break schemeLoop
			}
		}
		if !ecdsaOK {
			return false
		}
	}
	if hello.SupportedCurves != nil {
		ecdsaOK := false
		for _, curve := range hello.SupportedCurves {
			if curve == tls.CurveP256 {
				ecdsaOK = true
				break
			}
		}
		if !ecdsaOK {
			return false
		}
	}
	for _, suite := range hello.CipherSuites {
		switch suite {
		case tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305:
			return true
		}
	}
	return false
}
