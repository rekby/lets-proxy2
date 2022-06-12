package proxy

import (
	"bytes"
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/maxatome/go-testdeep"
	"github.com/rekby/fixenv"
	zc "github.com/rekby/zapcontext"
	"golang.org/x/xerrors"

	"github.com/rekby/lets-proxy2/internal/th"
)

func TestHttpProxy_HandleHttpValidationDefault(t *testing.T) {
	ctx, flush := th.TestContext(t)
	defer flush()

	td := testdeep.NewT(t)

	td.FailureIsFatal()
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1)})
	td.CmpNoError(err)

	defer func() { _ = listener.Close() }()

	proxy := NewHTTPProxy(ctx, listener)
	td.False(proxy.HandleHTTPValidation(&httptest.ResponseRecorder{}, nil))
}

func TestHttpProxy_getContextDefault(t *testing.T) {
	ctx, flush := th.TestContext(t)
	defer flush()

	td := testdeep.NewT(t)

	td.FailureIsFatal()
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1)})
	td.CmpNoError(err)
	defer func() { _ = listener.Close() }()

	proxy := NewHTTPProxy(ctx, listener)
	ctx2, err := proxy.GetContext(nil)
	td.NotNil(zc.L(ctx2))
	td.CmpNoError(err)
}

// nolint:unused
// need for mock generator
type HTTPProxyTest interface {
	GetContext(req *http.Request) (context.Context, error)
	HandleHTTPValidation(w http.ResponseWriter, r *http.Request) bool
}

func TestNewHttpProxy(t *testing.T) {
	var resp *http.Response
	var res []byte
	ctx, flush := th.TestContext(t)
	defer flush()

	td := testdeep.NewT(t)
	var mc = minimock.NewController(td)
	defer mc.Finish()

	td.FailureIsFatal()
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1)})
	defer th.Close(listener)
	td.CmpNoError(err)

	prefix := "http://" + listener.Addr().String()

	td.FailureIsFatal(false)

	transport := NewRoundTripperMock(mc)
	transport.RoundTripMock.Set(func(req *http.Request) (resp *http.Response, err error) {
		mux := http.ServeMux{}
		mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte{1, 2, 3})
		})
		respRecorder := &httptest.ResponseRecorder{
			Body: &bytes.Buffer{},
		}
		mux.ServeHTTP(respRecorder, req)
		return respRecorder.Result(), nil
	})

	proxyTest := NewHTTPProxyTestMock(mc)
	proxyTest.GetContextMock.Set(func(req *http.Request) (c1 context.Context, err2 error) {
		return ctx, nil
	})
	proxyTest.HandleHTTPValidationMock.Set(func(w http.ResponseWriter, r *http.Request) (b1 bool) {
		if strings.HasPrefix(r.URL.Path, "/asdf") {
			w.WriteHeader(http.StatusAccepted)
			_, _ = w.Write([]byte{3, 4})
			return true
		}
		return false
	})

	directorMock := NewDirectorMock(mc)
	directorMock.DirectorMock.Set(func(request *http.Request) error {
		if request.URL == nil {
			request.URL = &url.URL{}
		}
		request.URL.Scheme = ProtocolHTTP
		request.URL.Host = listener.Addr().String()
		return nil
	})

	proxy := NewHTTPProxy(ctx, listener)
	defer th.Close(proxy)

	proxy.GetContext = proxyTest.GetContext
	proxy.Director = directorMock
	proxy.HandleHTTPValidation = proxyTest.HandleHTTPValidation
	proxy.HTTPTransport = transport
	go func() { _ = proxy.Start() }()

	//nolint:gosec
	resp, err = http.Get(prefix)
	td.CmpNoError(err)
	td.CmpDeeply(http.StatusOK, resp.StatusCode)
	res, err = ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	td.CmpNoError(err)
	td.CmpDeeply(res, []byte{1, 2, 3})

	resp, err = http.Get(prefix + "/asdfg")
	td.CmpNoError(err)
	td.CmpDeeply(http.StatusAccepted, resp.StatusCode)
	res, err = ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	td.CmpNoError(err)
	td.CmpDeeply(res, []byte{3, 4})
}

func TestNoDoubleSlashredirectIssue177(t *testing.T) {
	e, _, flush := th.NewEnv(t)
	defer flush()

	_, addr := httpProxy(e, th.HttpQuery(e))

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Get(addr + "/test//")
	e.CmpNoError(err)

	body, err := ioutil.ReadAll(resp.Body)
	e.CmpNoError(err)
	e.CmpNoError(resp.Body.Close())

	e.Cmp(string(body), "/test//")
}

func httpProxy(e *th.Env, dst string) (*HTTPProxy, string) {
	var proxy *HTTPProxy
	var addr string
	e.CacheWithCleanup(dst, nil, func() (res interface{}, cleanup fixenv.FixtureCleanupFunc, err error) {
		dstURL, err := url.Parse(dst)
		if err != nil {
			return nil, nil, xerrors.Errorf("failed parse dst url %q: %%w", dst, err)
		}
		listener := th.NewLocalTcpListener(e)
		addr = "http://" + listener.Addr().String()
		proxy = NewHTTPProxy(e.Ctx, listener)
		proxy.Director = NewDirectorChain(
			DirectorHost(dstURL.Host),
			DirectorSetScheme(dstURL.Scheme),
		)
		errChan := make(chan error, 1)

		go func() { errChan <- proxy.Start() }()

		select {
		case proxyStartErr := <-errChan:
			err = proxyStartErr
		case <-time.After(time.Millisecond * 10):
			// pass
		}

		cleanup = func() { _ = proxy.Close() }
		return proxy, cleanup, err
	})
	return proxy, addr
}
