//nolint:golint
package cert_manager

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep"

	"github.com/rekby/lets-proxy2/internal/proxy"

	"github.com/rekby/lets-proxy2/internal/domain"

	"github.com/rekby/lets-proxy2/internal/docker"

	"github.com/gojuno/minimock/v3"
	"github.com/rekby/lets-proxy2/internal/cache"
	"github.com/rekby/lets-proxy2/internal/th"
	zc "github.com/rekby/zapcontext"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme"
)

const forceRsaDomain = "force-rsa.ru"
const testCertIssueTimeout = time.Second * 30

func TestManager_GetCertificateHttp01(t *testing.T) {
	t.Parallel()

	ctx, flush := th.TestContext(t)
	defer flush()

	logger := zc.L(ctx)

	mc := minimock.NewController(t)
	defer mc.Finish()

	manager := New(createTestClientManager(t), newCacheMock(mc), nil)
	manager.CertificateIssueTimeout = testCertIssueTimeout
	manager.AutoSubdomains = []string{"www."}
	manager.EnableTLSValidation = false
	manager.EnableHTTPValidation = true

	lisneter, err := net.ListenTCP("tcp", &net.TCPAddr{Port: 5002})

	if err != nil {
		t.Fatal(err)
	}

	//noinspection GoUnhandledErrorResult
	defer lisneter.Close()

	go func() {
		ctx := zc.WithLogger(context.Background(), logger)
		server := http.Server{}
		mux := http.ServeMux{}
		mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
			request = request.WithContext(ctx)
			if manager.isHTTPValidationRequest(request) {
				requestStr := fmt.Sprintf("%v %v", request.Method, request.URL)
				logger.Info("Handle validation request", zap.String("request", requestStr))
				manager.HandleHTTPValidation(writer, request)
			} else {
				logger.Warn("Handle non validation request")
				writer.WriteHeader(http.StatusInternalServerError)
			}
		})
		server.Handler = &mux
		go func() {
			<-ctx.Done()
			_ = server.Shutdown(context.Background())
		}()
		err = server.Serve(lisneter)
		logger.Info("http server shutdown", zap.Error(err))
	}()

	getCertificatesTests(t, manager, ctx, logger)
}

func TestManager_GetCertificateTls(t *testing.T) {
	t.Parallel()

	ctx, flush := th.TestContext(t)
	defer flush()

	logger := zc.L(ctx)

	mc := minimock.NewController(t)
	defer mc.Finish()

	manager := New(createTestClientManager(t), newCacheMock(mc), nil)
	manager.CertificateIssueTimeout = testCertIssueTimeout
	manager.AutoSubdomains = []string{"www."}
	manager.EnableTLSValidation = true
	manager.EnableHTTPValidation = false

	lisneter, err := net.ListenTCP("tcp", &net.TCPAddr{Port: 5001})

	if err != nil {
		t.Fatal(err)
	}

	//noinspection GoUnhandledErrorResult
	defer lisneter.Close()

	go func() {
		counter := 0
		for {
			conn, err := lisneter.Accept()
			if conn != nil {
				t.Log("incoming connection")
				ctx := zc.WithLogger(context.Background(), logger.With(zap.Int("connection_id", counter)))

				tlsConn := tls.Server(contextConnection{conn, ctx}, &tls.Config{
					NextProtos: []string{
						"h2", "http/1.1", // enable HTTP/2
						acme.ALPNProto, // enable tls-alpn ACME challenges
					},
					GetCertificate: manager.GetCertificate,
				})

				err := tlsConn.Handshake()
				if err == nil {
					t.Log("Handshake ok")
				} else {
					t.Error(err)
				}

				err = conn.Close()
				if err != nil {
					t.Error(err)
				}
			}
			if err != nil {
				break
			}
		}
	}()
	getCertificatesTests(t, manager, ctx, logger)
}

func TestDockerFindDomain(t *testing.T) {
	ctx, flush := th.TestContext(t)
	defer flush()

	dClient, err := docker.New(docker.Config{LabelDomain: "lets-proxy.domain", DefaultHttpPort: 80})
	if err != nil {
		t.Fatal(err)
	}
	dn, err := domain.NormalizeDomain("docker-test.internal")
	if err != nil {
		t.Fatal(err)
	}
	target, err := dClient.GetTarget(ctx, dn)
	if err != nil {
		t.Error(err)
	}

	ipS, port, err := net.SplitHostPort(target.TargetAddress)
	if err != nil {
		t.Error(err)
	}
	if port != "80" {
		t.Errorf("need port '80', got port '%v'", port)
	}
	ip := net.ParseIP(ipS)
	if len(ip) == 0 {
		t.Error("can't parse ip")
	}
}

func TestProxyToDocker(t *testing.T) {
	td := testdeep.NewT(t)
	ctx, flush := th.TestContext(t)
	defer flush()

	lisneter, err := net.ListenTCP("tcp", &net.TCPAddr{Port: 5003})
	if err != nil {
		t.Fatalf("Create listener: %v", err)
	}
	defer func() { _ = lisneter.Close() }()

	dockerClient, err := docker.New(docker.Config{
		DefaultHttpPort: 80,
		LabelDomain:     "lets-proxy.domain",
	})
	proxyConfig := proxy.Config{
		DefaultTarget: "127.0.0.1",
	}

	p := proxy.NewHTTPProxy(ctx, lisneter)
	err = proxyConfig.Apply(ctx, p, dockerClient)
	td.CmpNoError(err)

	go func() {
		err := p.Start()
		if err != http.ErrServerClosed {
			td.CmpNoError(err)
		}
	}()
	defer func() {
		_ = p.Close()
		time.Sleep(time.Second / 10)
	}()

	req, err := http.NewRequest(http.MethodGet, "http://docker-test.internal", nil)
	td.CmpNoError(err)

	req.URL.Host = "localhost:5003"
	res, err := http.DefaultClient.Do(req)
	td.CmpNoError(err)

	body, err := ioutil.ReadAll(res.Body)
	td.CmpNoError(err)
	_ = res.Body.Close()
	bodyS := strings.ToLower(string(body))
	if !strings.Contains(bodyS, "nginx") {
		td.Errorf("Unexpected body:\n'''\n%v\n'''", bodyS)
	}
}

func fastCreateTestCert(domains []string, now time.Time) (certBytes, keyBytes []byte) {
	template := x509.Certificate{
		SerialNumber: big.NewInt(123),
		Subject:      pkix.Name{CommonName: domains[0]},
		NotBefore:    now.Add(-time.Hour),
		NotAfter:     now.Add(time.Hour),
		DNSNames:     domains,
	}
	priv, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		panic(err)
	}
	der, err := x509.CreateCertificate(rand.Reader, &template, &template, priv.Public(), priv)
	if err != nil {
		panic(err)
	}

	certBytes = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyBytes = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	return certBytes, keyBytes
}

func newCacheMock(t minimock.Tester) *BytesMock {
	cacheMock := NewBytesMock(t)
	forceRSACert, forceRSAKey := fastCreateTestCert([]string{"force-rsa.ru", "www.force-rsa.ru"}, time.Now())
	cacheMock.GetMock.Set(func(ctx context.Context, key string) (ba1 []byte, err error) {
		zc.L(ctx).Debug("Cache mock get", zap.String("key", key))

		switch key {
		case "locked.ru.lock":
			return []byte{}, nil

		// force-rsa locked cert
		case forceRsaDomain + ".lock":
			return []byte{}, nil
		case forceRsaDomain + ".rsa.cer":
			return forceRSACert, nil
		case forceRsaDomain + ".rsa.key":
			return forceRSAKey, nil
		}

		return nil, cache.ErrCacheMiss
	})
	cacheMock.PutMock.Set(func(ctx context.Context, key string, data []byte) (err error) {
		zc.L(ctx).Debug("Cache mock put", zap.String("key", key))
		return nil
	})
	return cacheMock
}

func getCertificatesTests(t *testing.T, manager *Manager, ctx context.Context, logger *zap.Logger) {
	t.Run("ECDSA", func(t *testing.T) {
		getCertificatesTestsKeyType(t, manager, KeyECDSA, ctx, logger)
	})
	t.Run("RSA", func(t *testing.T) {
		getCertificatesTestsKeyType(t, manager, KeyRSA, ctx, logger)
	})
}

func getCertificatesTestsKeyType(t *testing.T, manager *Manager, keyType KeyType, ctx context.Context, logger *zap.Logger) {
	t.Run("OneCert", func(t *testing.T) {
		t.Parallel()
		checkOkDomain(ctx, t, manager, keyType, keyType, "onecert.ru")
	})

	t.Run("punycode-domain", func(t *testing.T) {
		t.Parallel()
		checkOkDomain(ctx, t, manager, keyType, keyType, "xn--80adjurfhd.xn--p1ai") // проверка.рф
	})

	t.Run("OneCertCamelCase", func(t *testing.T) {
		t.Parallel()
		checkOkDomain(ctx, t, manager, keyType, keyType, "onecertCamelCase.ru")
	})

	t.Run("Locked", func(t *testing.T) {
		td := testdeep.NewT(t)
		t.Parallel()
		domain := "locked.ru"

		cert, err := manager.GetCertificate(createTLSHello(ctx, keyType, domain))
		td.CmpError(err)
		td.Nil(cert)
	})

	//nolint[:dupl]
	t.Run("ParallelCert", func(t *testing.T) {
		td := testdeep.NewT(t)
		t.Parallel()

		// change top level logger
		oldLogger := logger
		logger = zap.NewNop()
		defer func() {
			logger = oldLogger
		}()

		domain := "ParallelCert.ru"
		const cnt = 100

		chanCerts := make(chan *tls.Certificate, cnt)

		var wg sync.WaitGroup
		wg.Add(cnt)

		for i := 0; i < cnt; i++ {
			go func() {
				cert, err := manager.GetCertificate(createTLSHello(ctx, keyType, domain))
				if err != nil {
					t.Error(err)
				}
				chanCerts <- cert
				wg.Done()
			}()
		}

		wg.Wait()
		cert := <-chanCerts
		for i := 0; i < len(chanCerts)-1; i++ {
			cert2 := <-chanCerts
			td.CmpDeeply(cert, cert2)
		}
	})

	t.Run("RenewSoonExpiredCert", func(t *testing.T) {
		t.Parallel()
		const domain = "soon-expired.com"

		// issue certificate
		cert, err := manager.GetCertificate(createTLSHello(ctx, keyType, domain))
		if err != nil {
			t.Errorf("cant issue certificate: %v", err)
			return
		}
		certNumber := cert.Leaf.SerialNumber
		newExpire := time.Now().Add(time.Hour)
		cert.Leaf.NotAfter = newExpire

		// get expired soon certificate and trigger reissue new
		cert, err = manager.GetCertificate(createTLSHello(ctx, keyType, domain))
		if err != nil {
			t.Errorf("cant issue certificate: %v", err)
			return
		}
		if certNumber.Cmp(cert.Leaf.SerialNumber) != 0 {
			t.Error("Got other sertificate, need same.")
		}
		if !cert.Leaf.NotAfter.Equal(newExpire) {
			t.Errorf("Bad expire time: '%v' instead of '%v'", cert.Leaf.NotAfter, newExpire)
		}

		time.Sleep(time.Second * 10)

		// get renewed cert
		cert, err = manager.GetCertificate(createTLSHello(ctx, keyType, domain))
		if err != nil {
			t.Errorf("cant issue certificate: %v", err)
			return
		}
		if certNumber.Cmp(cert.Leaf.SerialNumber) == 0 {
			t.Error("Need new certificate")
		}
		if !cert.Leaf.NotAfter.After(newExpire) {
			t.Errorf("Bad expire time: %v", cert.Leaf.NotAfter)
		}
	})

	t.Run("Force-rsa", func(t *testing.T) {
		t.Parallel()
		checkOkDomain(ctx, t, manager, keyType, KeyRSA, forceRsaDomain)
	})
}

func createTLSHello(ctx context.Context, keyType KeyType, domain string) *tls.ClientHelloInfo {
	switch keyType {
	case KeyRSA:
		return &tls.ClientHelloInfo{
			ServerName: domain,
			Conn:       contextConnection{Context: ctx},
		}
	case KeyECDSA:
		return &tls.ClientHelloInfo{
			ServerName:       domain,
			Conn:             contextConnection{Context: ctx},
			SignatureSchemes: []tls.SignatureScheme{tls.ECDSAWithP256AndSHA256},
			SupportedCurves:  []tls.CurveID{tls.CurveP256},
			CipherSuites:     []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384},
		}
	default:
		panic("unexpected key type")
	}
}

func checkOkDomain(ctx context.Context, t *testing.T, manager *Manager, keyTypeHello KeyType, keyTypeCert KeyType, domain string) {
	cert, err := manager.GetCertificate(createTLSHello(ctx, keyTypeHello, domain))
	if err != nil {
		t.Fatal(err)
	}
	if getKeyType(cert) != keyTypeCert {
		t.Errorf("Bad certificate key type. Expected: '%v', got: '%v'", keyTypeCert, getKeyType(cert))
	}

	certDomain := strings.TrimPrefix(strings.ToLower(domain), "www.")
	if cert.Leaf.NotBefore.After(time.Now()) {
		t.Error(cert.Leaf.NotBefore)
	}
	if cert.Leaf.NotAfter.Before(time.Now()) {
		t.Error(cert.Leaf.NotAfter)
	}
	if cert.Leaf.VerifyHostname(certDomain) != nil {
		t.Error(cert.Leaf.VerifyHostname(certDomain))
	}
	if cert.Leaf.VerifyHostname("www."+certDomain) != nil {
		t.Error(cert.Leaf.VerifyHostname(certDomain))
	}
}
