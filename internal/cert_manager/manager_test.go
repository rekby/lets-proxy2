//nolint:golint
package cert_manager

//go:generate minimock -i github.com/rekby/lets-proxy2/internal/cache.Bytes -o ./cache_bytes_mock_test.go

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/rekby/fixenv"
	"github.com/rekby/lets-proxy2/internal/cache"

	"go.uber.org/zap"

	zc "github.com/rekby/zapcontext"

	"github.com/gojuno/minimock/v3"

	"github.com/maxatome/go-testdeep"

	"github.com/rekby/lets-proxy2/internal/th"

	"golang.org/x/crypto/acme"
)

const rsaKeyLength = 2048

type contextConnection struct {
	net.Conn
	context.Context
}

func (c contextConnection) GetContext() context.Context {
	return c.Context
}

//nolint:gochecknoinits
func init() {
	zc.SetDefaultLogger(zap.NewNop())
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
}

func createTestClientManager(env fixenv.Env, t *testing.T) *AcmeClientManagerMock {
	resp, err := http.Get(th.AcmeServerDirURL(env))
	if err != nil {
		t.Fatalf("Can't connect to buoulder server: %q", err)
	}
	_ = resp.Body.Close()

	client := acme.Client{}
	client.HTTPClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				//nolint:gosec
				InsecureSkipVerify: true,
			},
		},
	}

	client.DirectoryURL = th.AcmeServerDirURL(env)
	client.Key, _ = rsa.GenerateKey(rand.Reader, rsaKeyLength)
	_, err = client.Register(context.Background(), &acme.Account{}, func(tosURL string) bool {
		return true
	})

	if err != nil {
		t.Fatalf("Can't initialize acme client: %v", err)
	}

	clientManager := NewAcmeClientManagerMock(t)
	clientManager.CloseMock.Return(nil)
	clientManager.GetClientMock.Return(&client, func() {}, nil)
	return clientManager
}

func TestGetKeyType(t *testing.T) {
	td := testdeep.NewT(t)
	cert := &tls.Certificate{
		PrivateKey: &rsa.PrivateKey{},
	}
	td.CmpDeeply(getKeyType(cert), KeyRSA)

	cert = &tls.Certificate{
		PrivateKey: &ecdsa.PrivateKey{},
	}
	td.CmpDeeply(getKeyType(cert), KeyECDSA)

	cert = &tls.Certificate{
		PrivateKey: "string - no key",
	}
	td.CmpPanic(func() {
		getKeyType(cert)
	}, "unexexpected key type: .string")

	td.CmpPanic(func() {
		getKeyType(nil)
	}, "cert is nil")
}

func TestStoreLoadCertificate(t *testing.T) {
	e, ctx, flush := th.NewEnv(t)
	defer flush()

	certBytes, keyBytes := fastCreateTestCert([]string{"domain.com"}, time.Now())
	cert, err := tls.X509KeyPair(certBytes, keyBytes)
	e.CmpNoError(err)
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	e.CmpNoError(err)

	mc := minimock.NewController(t)
	c := make(map[string][]byte)
	cacheMock := NewBytesMock(mc)
	cacheMock.PutMock.Set(func(ctx context.Context, key string, data []byte) (err error) {
		c[key] = data
		return nil
	})
	cacheMock.GetMock.Set(func(ctx context.Context, key string) (ba1 []byte, err error) {
		data, ok := c[key]
		if ok {
			return data, nil
		}
		return nil, cache.ErrCacheMiss
	})

	cd := CertDescription{MainDomain: "asd", KeyType: KeyRSA}
	err = storeCertificate(ctx, cacheMock, cd, &cert)
	e.CmpNoError(err)

	resCert, err := loadCertificateFromCache(ctx, cacheMock, cd)
	e.CmpNoError(err)

	e.CmpNoError(err)
	e.Cmp(resCert, &cert)
}

func TestIsNeedRenew(t *testing.T) {
	td := testdeep.NewT(t)
	var cert = &tls.Certificate{}
	cert.Leaf = &x509.Certificate{NotAfter: time.Date(2000, 7, 31, 0, 0, 0, 0, time.UTC)}
	td.True(isNeedRenew(cert, time.Date(2000, 7, 31, 0, 0, 0, 1, time.UTC)))
	td.True(isNeedRenew(cert, time.Date(2000, 7, 1, 0, 0, 0, 1, time.UTC)))
	td.False(isNeedRenew(cert, time.Date(2000, 7, 1, 0, 0, 0, 0, time.UTC)))
	td.False(isNeedRenew(cert, time.Date(2000, 6, 30, 0, 0, 0, 0, time.UTC)))
}

type testManagerContext struct {
	ctx context.Context

	manager                *Manager
	connContext            contextConnection
	conn                   *ConnMock
	cache                  *BytesMock
	certForDomainAuthorize *ValueMock
	certState              *ValueMock
	clientManager          *AcmeClientManagerMock
	domainChecker          *DomainCheckerMock
	httpTokens             *BytesMock
}

func TestManager_CertForLockedDomain(t *testing.T) {
	td := testdeep.NewT(t)
	c, cancel := createManager(t)
	defer cancel()

	c.certState.GetMock.Return(&certState{}, nil)
	c.cache.GetMock.Set(func(ctx context.Context, key string) (ba1 []byte, err error) {
		if key == "test.ru.lock" {
			return []byte{}, nil
		}
		return nil, cache.ErrCacheMiss
	})

	res, err := c.manager.GetCertificate(&tls.ClientHelloInfo{Conn: c.connContext, ServerName: "test.ru"})
	td.Nil(res)
	td.CmpError(err)
}

func TestManager_CertForDenied(t *testing.T) {
	td := testdeep.NewT(t)
	c, cancel := createManager(t)
	defer cancel()

	c.certState.GetMock.Return(&certState{}, nil)
	c.cache.GetMock.Return(nil, cache.ErrCacheMiss)
	c.domainChecker.IsDomainAllowedMock.Return(false, nil)

	res, err := c.manager.GetCertificate(&tls.ClientHelloInfo{Conn: c.connContext, ServerName: "test.ru"})
	td.Nil(res)
	td.CmpError(err)
}

func TestManagerFilterTlsHello(t *testing.T) {
	t.Run("AllowInsecureChipers_True", func(t *testing.T) {
		e, ctx, flush := th.NewEnv(t)
		defer flush()

		m := Manager{}
		m.AllowInsecureTLSChipers = true

		hello := tls.ClientHelloInfo{
			CipherSuites: []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256, tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA, tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384},
		}
		m.filterTlsHello(ctx, &hello)

		expectedHello := tls.ClientHelloInfo{
			CipherSuites: []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256, tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA, tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384},
		}
		e.Cmp(hello, expectedHello)
	})
	t.Run("AllowInsecureChipers_False", func(t *testing.T) {
		e, ctx, flush := th.NewEnv(t)
		defer flush()

		m := Manager{}
		m.AllowInsecureTLSChipers = false

		hello := tls.ClientHelloInfo{
			CipherSuites: []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256, tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA, tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384},
		}
		m.filterTlsHello(ctx, &hello)

		expectedHello := tls.ClientHelloInfo{
			CipherSuites: []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256, tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384},
		}
		e.Cmp(hello, expectedHello)
	})
}

func TestGetCertificateDenyCertificates(t *testing.T) {
	td := testdeep.NewT(t)
	m := Manager{}

	_, err := m.getCertificate(nil, "", KeyRSA)
	td.Cmp(err, errRSADenied)

	_, err = m.getCertificate(nil, "", KeyECDSA)
	td.Cmp(err, errECDSADenied)

	_, err = m.getCertificate(nil, "", "")
	td.Cmp(err, errCertTypeUnknown)
}

func createManager(t *testing.T) (res testManagerContext, cancel func()) {
	ctx, ctxCancel := th.TestContext(t)
	mc := minimock.NewController(t)

	res.ctx = ctx
	res.conn = NewConnMock(mc)
	res.connContext = contextConnection{
		Conn:    res.conn,
		Context: zc.WithLogger(context.Background(), zap.NewNop()),
	}
	res.cache = NewBytesMock(mc)
	res.clientManager = NewAcmeClientManagerMock(mc)
	res.certForDomainAuthorize = NewValueMock(mc)
	res.certState = NewValueMock(mc)
	res.domainChecker = NewDomainCheckerMock(mc)
	res.httpTokens = NewBytesMock(mc)

	res.manager = &Manager{
		CertificateIssueTimeout: time.Second,
		Cache:                   res.cache,
		acmeClientManager:       res.clientManager,
		DomainChecker:           res.domainChecker,
		EnableHTTPValidation:    true,
		EnableTLSValidation:     true,
		AllowRSACert:            true,
		AllowECDSACert:          true,
		certForDomainAuthorize:  res.certForDomainAuthorize,
		certState:               res.certState,
		httpTokens:              res.httpTokens,
	}
	res.manager.initMetrics(nil)
	return res, func() {
		mc.Finish()
		ctxCancel()
	}
}
