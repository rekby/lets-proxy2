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

	"github.com/gojuno/minimock"

	zc "github.com/rekby/zapcontext"
	"go.uber.org/zap"

	"github.com/maxatome/go-testdeep"

	"github.com/rekby/lets-proxy2/internal/th"
)

func TestGetDestination(t *testing.T) {
	ctx, flush := th.TestContext()
	defer flush()

	var dest string
	var err error

	td := testdeep.NewT(t)

	dest, err = getDestination(ctx, "127.0.0.1:443")
	td.String(dest, "127.0.0.1:80")
	td.CmpNoError(err)

	dest, err = getDestination(ctx, "127.0.0.1:444")
	td.String(dest, "127.0.0.1:80")
	td.CmpNoError(err)

	dest, err = getDestination(ctx, "127.0.0.2:443")
	td.String(dest, "127.0.0.2:80")
	td.CmpNoError(err)

	dest, err = getDestination(ctx, "127.0.0.2")
	td.String(dest, "")
	td.CmpError(err)
}

func TestHttpProxy_SetTransport(t *testing.T) {
	td := testdeep.NewT(t)

	proxy := HttpProxy{}
	transport := NewRoundTripperMock(t)
	proxy.SetTransport(transport)

	td.CmpDeeply(proxy.httpReverseProxy.Transport, transport)
	transport.MinimockFinish()

}

func TestHttpProxy_HandleHttpValidationDefault(t *testing.T) {
	ctx, flush := th.TestContext()
	defer flush()

	td := testdeep.NewT(t)

	td.FailureIsFatal()
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1)})
	td.CmpNoError(err)
	defer func() { _ = listener.Close() }()

	proxy := NewHttpProxy(ctx, listener)
	td.False(proxy.HandleHttpValidation(&httptest.ResponseRecorder{}, nil))
}

func TestHttpProxy_getContextDefault(t *testing.T) {
	ctx, flush := th.TestContext()
	defer flush()

	td := testdeep.NewT(t)

	td.FailureIsFatal()
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1)})
	td.CmpNoError(err)
	defer func() { _ = listener.Close() }()

	proxy := NewHttpProxy(ctx, listener)
	ctx2 := proxy.GetContext(nil)
	td.NotNil(zc.L(ctx2))
}

func TestHttpProxy_Director(t *testing.T) {
	ctx, flush := th.TestContext()
	defer flush()

	var req *http.Request
	td := testdeep.NewT(t)
	proxy := HttpProxy{}
	proxy.GetContext = func(req *http.Request) context.Context {
		return zc.WithLogger(ctx, zap.NewNop())
	}
	proxy.GetDestination = func(ctx context.Context, remoteAddr string) (addr string, err error) {
		return "1.2.3.4:80", err
	}

	req = &http.Request{}
	proxy.director(req)
	td.CmpDeeply(req, &http.Request{URL: &url.URL{Host: "1.2.3.4:80", Scheme: "http"}})
}

type HttpProxyTest interface {
	GetDestination(ctx context.Context, remoteAddr string) (addr string, err error)
	GetContext(req *http.Request) context.Context
	HandleHttpValidation(w http.ResponseWriter, r *http.Request) bool
}

func TestNewHttpProxy(t *testing.T) {
	var resp *http.Response
	var res []byte
	ctx, flush := th.TestContext()
	defer flush()

	td := testdeep.NewT(t)
	var mc = minimock.NewController(td)
	defer mc.Finish()

	td.FailureIsFatal()
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1)})
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

	proxyTest := NewHttpProxyTestMock(mc)
	proxyTest.GetContextMock.Set(func(req *http.Request) (c1 context.Context) {
		return ctx
	})
	proxyTest.HandleHttpValidationMock.Set(func(w http.ResponseWriter, r *http.Request) (b1 bool) {
		if strings.HasPrefix(r.URL.Path, "/asdf") {
			w.WriteHeader(http.StatusAccepted)
			_, _ = w.Write([]byte{3, 4})
			return true
		} else {
			return false
		}
	})

	proxy := NewHttpProxy(ctx, listener)
	proxy.GetContext = proxyTest.GetContext
	proxy.GetDestination = proxyTest.GetDestination
	proxy.HandleHttpValidation = proxyTest.HandleHttpValidation
	proxy.SetTransport(transport)
	proxyTest.GetDestinationMock.Return("1.2.3.4:80", nil)

	resp, err = http.Get(prefix)
	td.CmpNoError(err)
	td.CmpDeeply(http.StatusOK, resp.StatusCode)
	res, err = ioutil.ReadAll(resp.Body)
	td.CmpNoError(err)
	td.CmpDeeply(res, []byte{1, 2, 3})

	resp, err = http.Get(prefix + "/asdfg")
	td.CmpNoError(err)
	td.CmpDeeply(http.StatusAccepted, resp.StatusCode)
	res, err = ioutil.ReadAll(resp.Body)
	td.CmpNoError(err)
	td.CmpDeeply(res, []byte{3, 4})
}
