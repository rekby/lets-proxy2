package th

import (
	"crypto/tls"
	"net"
	"net/http"
	"os"

	"github.com/letsencrypt/pebble/v2/ca"
	"github.com/letsencrypt/pebble/v2/db"
	"github.com/letsencrypt/pebble/v2/va"
	"github.com/letsencrypt/pebble/v2/wfe"
	"github.com/rekby/fixenv"
	"github.com/rekby/lets-proxy2/internal/th/testcert"
)

const (
	pebbleHttpValidationPort = 5002
	pebbleTLSValidationPort  = 5001
)

type PebbleServer struct {
	HTTPSDirectoryURL string
	TLSCert           tls.Certificate
}

func Pebble(e fixenv.Env) *PebbleServer {
	return fixenv.CacheWithCleanup(e, nil, &fixenv.FixtureOptions{}, func() (*PebbleServer, fixenv.FixtureCleanupFunc, error) {
		const strictFalse = false

		db := db.NewMemoryStore()

		ca := ca.New(StdLogger(e), db, "", 0, 1, 0)
		_ = os.Setenv("PEBBLE_VA_NOSLEEP", "true")
		va := va.New(StdLogger(e), PebbleHTTPValidationPort(e), PebbleTLSValidationPort(e), strictFalse, "")

		wfe := wfe.New(StdLogger(e), db, va, ca, strictFalse, false)
		mux := wfe.Handler()

		tlsConfig := &tls.Config{Certificates: []tls.Certificate{LocalhostCert(e)}}
		tlsListener := tls.NewListener(acmeServerListener(e), tlsConfig)
		httpServer := http.Server{}
		httpServer.Handler = mux
		go func() {
			err := httpServer.Serve(tlsListener)
			if err != http.ErrServerClosed {
				panic(err)
			}
		}()

		res := &PebbleServer{
			HTTPSDirectoryURL: AcmeServerDirURL(e),
		}
		cleanup := func() {
			_ = httpServer.Close()
			_ = tlsListener.Close()
		}
		return res, cleanup, nil
	})
}

func PebbleHTTPValidationPort(e fixenv.Env) int {
	return pebbleHttpValidationPort
}

func PebbleTLSValidationPort(e fixenv.Env) int {
	return pebbleTLSValidationPort
}

func acmeServerListener(e fixenv.Env) *net.TCPListener {
	return fixenv.Cache(e, nil, nil, func() (*net.TCPListener, error) {
		return NewLocalTcpListener(e), nil
	})
}

func AcmeServerDirURL(e fixenv.Env) string {
	return "https://" + acmeServerListener(e).Addr().String() + "/dir"
}

func LocalhostCert(e fixenv.Env) tls.Certificate {
	return fixenv.Cache(e, nil, nil, func() (tls.Certificate, error) {
		return tls.X509KeyPair(testcert.LocalhostCert, testcert.LocalhostKey)
	})
}
