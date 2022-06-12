package th

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"

	"github.com/gojuno/minimock/v3"
	"github.com/rekby/fixenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

var (
	freeListenAddressMutex sync.Mutex
	freeListenAddressUsed  = map[string]bool{}
)

func MockController(e fixenv.Env) minimock.MockController {
	var c *minimock.Controller

	return fixenv.CacheWithCleanup(e, nil, nil, func() (minimock.MockController, fixenv.FixtureCleanupFunc, error) {
		c = minimock.NewController(e.T().(minimock.Tester))
		return c, c.Finish, nil
	})
}

func NewLocalTcpListener(e fixenv.Env) *net.TCPListener {
	//antiCache, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt32))
	//if err != nil {
	//	e.T().Fatalf("Failed create upper bound for random anticache")
	//}
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		e.T().Fatalf("Failed create tcp listener: %w", err)
	}

	e.T().Cleanup(func() {
		_ = listener.Close()
	})

	return listener.(*net.TCPListener)
}

func NewFreeLocalTcpAddress(e fixenv.Env) *net.TCPAddr {
	for {
		listener := NewLocalTcpListener(e)
		_ = listener.Close()

		addr := listener.Addr()
		addrS := addr.String()
		freeListenAddressMutex.Lock()
		used := freeListenAddressUsed[addrS]
		freeListenAddressUsed[addrS] = true
		freeListenAddressMutex.Unlock()

		if !used {
			return addr.(*net.TCPAddr)
		}
	}
}

func HttpQuery(e fixenv.Env) string {
	var server *httptest.Server
	return fixenv.CacheWithCleanup(e, nil, nil, func() (res string, cleanupFunc fixenv.FixtureCleanupFunc, err error) {
		server = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			u := *request.URL
			u.Scheme = ""
			u.Host = ""
			_, _ = writer.Write([]byte(u.String()))
		}))
		return server.URL, server.Close, nil
	})
}

func TmpDir(e fixenv.Env) string {
	var dirPath string
	return fixenv.CacheWithCleanup(e, nil, nil, func() (res string, cleanup fixenv.FixtureCleanupFunc, err error) {
		dirPath, err = ioutil.TempDir("", "lets-proxy2-test-")
		if err != nil {
			cleanup = func() {
				_ = os.RemoveAll(dirPath)
			}
		}
		return dirPath, cleanup, err
	})
}

func StdLogger(e fixenv.Env) *log.Logger {
	return fixenv.Cache(e, "", nil, func() (*log.Logger, error) {
		return log.New(testWriter{e.T().(testLogger)}, "", log.LstdFlags), nil
	})
}

func ZapLogger(e fixenv.Env) *zap.Logger {
	return fixenv.CacheWithCleanup(e, nil, nil, func() (*zap.Logger, fixenv.FixtureCleanupFunc, error) {
		logger := zaptest.NewLogger(e.T().(zaptest.TestingT), zaptest.WrapOptions(zap.Development()))
		return logger, func() { _ = logger.Sync() }, nil
	})
}

type testLogger interface {
	Log(args ...any)
}
type testWriter struct {
	t testLogger
}

func (t testWriter) Write(p []byte) (n int, err error) {
	t.t.Log(string(p))
	return len(p), nil
}
