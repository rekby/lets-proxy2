package manager

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"net"
	"net/http"
	"testing"
	"time"

	"golang.org/x/crypto/acme"
)

const testACMEServer = "http://localhost:4000/directory"
const rsaKeyLength = 2048

func TestManager_GetCertificate(t *testing.T) {

	domain := "asdfadfwefdc.com"

	manager := Manager{
		Client:                  createTestClient(t),
		GetCertContext:          context.Background(),
		CertificateIssueTimeout: time.Minute,
	}

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
}

func createTestClient(t *testing.T) *acme.Client {
	_, err := http.Get(testACMEServer)
	if err != nil {
		t.Fatal("Can't connect to buoulder server")
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
