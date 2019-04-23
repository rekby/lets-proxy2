package tlslistener

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/rekby/lets-proxy2/internal/th"
)

func TestProxyListener(t *testing.T) {

	listener := listenerType{connections: make(chan net.Conn)}

	// test proxy
	conn1 := &net.TCPConn{}
	go func() {
		err := listener.Put(conn1)
		if err != nil {
			t.Error(err)
		}
	}()

	rconn1, err := listener.Accept()
	if err != nil {
		t.Error()
	}
	if rconn1.(*net.TCPConn) != conn1 {
		t.Error()
	}

	conn2 := &net.TCPConn{}
	go func() {
		err := listener.Put(conn2)
		if err != nil {
			t.Error(err)
		}
	}()

	rconn2, err := listener.Accept()
	if err != nil {
		t.Error()
	}
	if rconn2.(*net.TCPConn) != conn2 {
		t.Error()
	}

	// test reject accept on close
	go func() {
		time.Sleep(time.Millisecond)
		err := listener.Close()
		if err != nil {
			t.Error(err)
		}
	}()

	rconn3, err := listener.Accept()
	if rconn3 != nil {
		t.Error(rconn3)
	}
	if err == nil {
		t.Error()
	}

	// second close
	err = listener.Close()
	if err == nil {
		t.Error()
	}
}

func TestProxyTLS(t *testing.T) {
	ctx, flush := th.TestContext()
	defer flush()

	httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, "Hello, client")
	}))
	defer httpServer.Close()

	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0})
	if err != nil {
		t.Fatal(err)
	}

	testUrl := "https://" + listener.Addr().String()

	proxy := proxyType{
		GetCertificate: dummyGetCertificate,
		TargetAddr:     strings.TrimPrefix(httpServer.URL, "http://"),
		TLSListener:    listener,
	}

	go func() {
		_ = proxy.Start(ctx)
	}()

	httpClient := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	resp, err := httpClient.Get(testUrl)
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()

	bodyS := string(body)
	if bodyS != "Hello, client" {
		t.Error(bodyS)
	}
}

func dummyGetCertificate(info *tls.ClientHelloInfo) (certificate *tls.Certificate, e error) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)

	certTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(123),
		DNSNames:     []string{info.ServerName}, NotAfter: time.Now().Add(time.Hour),
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, certTemplate, certTemplate, key.Public(), key)
	if err != nil {
		panic(err)
	}
	leaf, err := x509.ParseCertificate(certBytes)
	certificate = &tls.Certificate{
		Leaf:        leaf,
		PrivateKey:  key,
		Certificate: [][]byte{certBytes},
	}
	return certificate, err
}
