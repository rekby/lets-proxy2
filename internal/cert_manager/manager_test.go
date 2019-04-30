//nolint:golint
package cert_manager

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/rekby/lets-proxy2/internal/cache"

	"go.uber.org/zap"

	zc "github.com/rekby/zapcontext"

	"github.com/gojuno/minimock"

	"github.com/maxatome/go-testdeep"
	td "github.com/maxatome/go-testdeep"
	"github.com/rekby/lets-proxy2/internal/th"

	"golang.org/x/crypto/acme"
)

const testACMEServer = "http://localhost:4000/directory"
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
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zc.SetDefaultLogger(logger)
}

func TestManager_GetCertificateTls(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}

	ctx, flush := th.TestContext()
	defer flush()

	mc := minimock.NewController(t)
	defer mc.Finish()

	cacheMock := NewCacheMock(mc)
	cacheMock.GetMock.Set(func(ctx context.Context, key string) (ba1 []byte, err error) {
		zc.L(ctx).Debug("Cache mock get", zap.String("key", key))
		return nil, cache.ErrCacheMiss
	})
	cacheMock.PutMock.Set(func(ctx context.Context, key string, data []byte) (err error) {
		zc.L(ctx).Debug("Cache mock put", zap.String("key", key))
		return nil
	})

	manager := New(ctx, createTestClient(t), cacheMock)

	lisneter, err := net.ListenTCP("tcp", &net.TCPAddr{Port: 5001})

	if err != nil {
		t.Fatal(err)
	}
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

	t.Run("OneCert", func(t *testing.T) {
		domain := "onecert.ru"

		cert, err := manager.GetCertificate(&tls.ClientHelloInfo{ServerName: domain, Conn: contextConnection{Context: ctx}})
		if err != nil {
			t.Fatal(err)
		}

		if cert.Leaf.NotBefore.After(time.Now()) {
			t.Error(cert.Leaf.NotBefore)
		}
		if cert.Leaf.NotAfter.Before(time.Now()) {
			t.Error(cert.Leaf.NotAfter)
		}
		if cert.Leaf.VerifyHostname(domain) != nil {
			t.Error(cert.Leaf.VerifyHostname(domain))
		}
		if cert.Leaf.VerifyHostname("www."+domain) != nil {
			t.Error(cert.Leaf.VerifyHostname(domain))
		}
	})

	t.Run("punycode-domain", func(t *testing.T) {
		domain := "xn--80adjurfhd.xn--p1ai" // проверка.рф

		cert, err := manager.GetCertificate(&tls.ClientHelloInfo{ServerName: domain, Conn: contextConnection{Context: ctx}})
		if err != nil {
			t.Fatal(err)
		}

		if cert.Leaf.NotBefore.After(time.Now()) {
			t.Error(cert.Leaf.NotBefore)
		}
		if cert.Leaf.NotAfter.Before(time.Now()) {
			t.Error(cert.Leaf.NotAfter)
		}
		if cert.Leaf.VerifyHostname(domain) != nil {
			t.Error(cert.Leaf.VerifyHostname(domain))
		}
		if cert.Leaf.VerifyHostname("www."+domain) != nil {
			t.Error(cert.Leaf.VerifyHostname(domain))
		}
	})

	t.Run("OneCertCamelCase", func(t *testing.T) {
		domain := "onecertCamelCase.ru"
		cert, err := manager.GetCertificate(&tls.ClientHelloInfo{ServerName: domain, Conn: contextConnection{Context: ctx}})
		if err != nil {
			t.Fatal(err)
		}

		if cert.Leaf.NotBefore.After(time.Now()) {
			t.Error(cert.Leaf.NotBefore)
		}
		if cert.Leaf.NotAfter.Before(time.Now()) {
			t.Error(cert.Leaf.NotAfter)
		}
		if cert.Leaf.VerifyHostname(domain) != nil {
			t.Error(cert.Leaf.VerifyHostname(domain))
		}
	})

	//nolint[:dupl]
	t.Run("ParallelCert", func(t *testing.T) {
		// change top loevel logger
		// no parallelize
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
				cert, err := manager.GetCertificate(&tls.ClientHelloInfo{ServerName: domain, Conn: contextConnection{Context: ctx}})
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
			testdeep.CmpDeeply(t, cert2, cert)
		}
	})
}

func TestManager_GetCertificateHttp01(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}

	ctx, flush := th.TestContext()
	defer flush()

	mc := minimock.NewController(t)
	defer mc.Finish()

	cacheMock := NewCacheMock(mc)
	cacheMock.GetMock.Set(func(ctx context.Context, key string) (ba1 []byte, err error) {
		zc.L(ctx).Debug("Cache mock get", zap.String("key", key))
		return nil, cache.ErrCacheMiss
	})
	cacheMock.PutMock.Set(func(ctx context.Context, key string, data []byte) (err error) {
		zc.L(ctx).Debug("Cache mock put", zap.String("key", key))
		return nil
	})

	manager := New(ctx, createTestClient(t), cacheMock)
	manager.EnableTLSValidation = false
	manager.EnableHTTPValidation = true

	lisneter, err := net.ListenTCP("tcp", &net.TCPAddr{Port: 5002})

	if err != nil {
		t.Fatal(err)
	}
	defer lisneter.Close()

	go func() {
		ctx := zc.WithLogger(context.Background(), logger)
		server := http.Server{}
		mux := http.ServeMux{}
		mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
			request = request.WithContext(ctx)
			if manager.isHTTPValidationRequest(request) {
				logger.Info("Handle validation request", zap.Reflect("request", request))
				manager.HandleHttpValidation(writer, request)
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

	t.Run("OneCert", func(t *testing.T) {
		domain := "onecert.ru"

		cert, err := manager.GetCertificate(&tls.ClientHelloInfo{ServerName: domain, Conn: contextConnection{Context: ctx}})
		if err != nil {
			t.Fatal(err)
		}

		if cert.Leaf.NotBefore.After(time.Now()) {
			t.Error(cert.Leaf.NotBefore)
		}
		if cert.Leaf.NotAfter.Before(time.Now()) {
			t.Error(cert.Leaf.NotAfter)
		}
		if cert.Leaf.VerifyHostname(domain) != nil {
			t.Error(cert.Leaf.VerifyHostname(domain))
		}
		if cert.Leaf.VerifyHostname("www."+domain) != nil {
			t.Error(cert.Leaf.VerifyHostname(domain))
		}
	})

	t.Run("punycode-domain", func(t *testing.T) {
		domain := "xn--80adjurfhd.xn--p1ai" // проверка.рф

		cert, err := manager.GetCertificate(&tls.ClientHelloInfo{ServerName: domain, Conn: contextConnection{Context: ctx}})
		if err != nil {
			t.Fatal(err)
		}

		if cert.Leaf.NotBefore.After(time.Now()) {
			t.Error(cert.Leaf.NotBefore)
		}
		if cert.Leaf.NotAfter.Before(time.Now()) {
			t.Error(cert.Leaf.NotAfter)
		}
		if cert.Leaf.VerifyHostname(domain) != nil {
			t.Error(cert.Leaf.VerifyHostname(domain))
		}
		if cert.Leaf.VerifyHostname("www."+domain) != nil {
			t.Error(cert.Leaf.VerifyHostname(domain))
		}
	})

	t.Run("OneCertCamelCase", func(t *testing.T) {
		domain := "onecertCamelCase.ru"
		cert, err := manager.GetCertificate(&tls.ClientHelloInfo{ServerName: domain, Conn: contextConnection{Context: ctx}})
		if err != nil {
			t.Fatal(err)
		}

		if cert.Leaf.NotBefore.After(time.Now()) {
			t.Error(cert.Leaf.NotBefore)
		}
		if cert.Leaf.NotAfter.Before(time.Now()) {
			t.Error(cert.Leaf.NotAfter)
		}
		if cert.Leaf.VerifyHostname(domain) != nil {
			t.Error(cert.Leaf.VerifyHostname(domain))
		}
	})

	//nolint[:dupl]
	t.Run("ParallelCert", func(t *testing.T) {
		// change top loevel logger
		// no parallelize
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
				cert, err := manager.GetCertificate(&tls.ClientHelloInfo{ServerName: domain, Conn: contextConnection{Context: ctx}})
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
			td.CmpDeeply(t, cert2, cert)
		}
	})
}

func createTestClient(t *testing.T) *acme.Client {
	_, err := http.Get(testACMEServer)
	if err != nil {
		t.Fatalf("Can't connect to buoulder server: %q", err)
	}

	client := acme.Client{}
	client.HTTPClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				//nolint:gosec
				InsecureSkipVerify: true,
			},
		},
	}

	client.DirectoryURL = testACMEServer
	client.Key, _ = rsa.GenerateKey(rand.Reader, rsaKeyLength)
	_, err = client.Register(context.Background(), &acme.Account{}, func(tosURL string) bool {
		return true
	})

	if err != nil {
		t.Fatal("Can't initialize acme client.")
	}
	return &client
}

func TestStoreCertificate(t *testing.T) {
	ctx, flush := th.TestContext()
	defer flush()

	//nolint:gosec
	key, _ := rsa.GenerateKey(rand.Reader, 512)

	cert := &tls.Certificate{Certificate: [][]byte{
		{1, 2, 3},
		{4, 5, 6},
	},
		PrivateKey: key,
	}

	mc := minimock.NewController(t)
	cacheMock := NewCacheMock(mc).PutMock.Set(func(ctx context.Context, key string, data []byte) (err error) {
		fmt.Println(key)
		fmt.Println(string(data))
		return nil
	})

	storeCertificate(ctx, cacheMock, "asd", cert)
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
