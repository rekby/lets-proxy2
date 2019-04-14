package manager

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	td "github.com/maxatome/go-testdeep"
	"github.com/rekby/lets-proxy2/internal/th"

	"golang.org/x/crypto/acme"
)

const testACMEServer = "http://localhost:4000/directory"
const rsaKeyLength = 2048

func TestManager_GetCertificate(t *testing.T) {
	ctx, flush := th.TestContext()
	defer flush()

	manager := New(ctx, createTestClient(t))

	lisneter, err := net.ListenTCP("tcp", &net.TCPAddr{Port: 5001})

	if err != nil {
		t.Fatal(err)
	}
	defer lisneter.Close()

	go func() {
		for {
			conn, err := lisneter.Accept()
			if conn != nil {
				t.Log("incoming connection")

				tlsConn := tls.Server(conn, &tls.Config{
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

				conn.Close()
			}
			if err != nil {
				break
			}
		}
	}()

	t.Run("OneCert", func(t *testing.T) {
		domain := "onecert.ru"
		cert, err := manager.GetCertificate(&tls.ClientHelloInfo{ServerName: domain})
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

	t.Run("OneCertCamelCase", func(t *testing.T) {
		domain := "onecertCamelCase.ru"
		cert, err := manager.GetCertificate(&tls.ClientHelloInfo{ServerName: domain})
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

	t.Run("ParallelCert", func(t *testing.T) {
		oldManagetCtx := manager.GetCertContext
		manager.GetCertContext = th.NoLog(ctx)
		defer func() {
			manager.GetCertContext = oldManagetCtx
		}()
		domain := "ParallelCert.ru"
		const cnt = 100

		chanCerts := make(chan *tls.Certificate, cnt)

		var wg sync.WaitGroup
		wg.Add(cnt)

		for i := 0; i < cnt; i++ {
			go func() {
				cert, err := manager.GetCertificate(&tls.ClientHelloInfo{ServerName: domain})
				if err != nil {
					t.Fatal(err)
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
