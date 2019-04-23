package proxy

import (
	"context"
	"net/http"
	"net/url"
	"testing"

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

func TestGetContext(t *testing.T) {
	ctx := getContext(nil)
	testdeep.CmpNotNil(t, zc.L(ctx))
}

//func TestNewHttpProxy(t *testing.T) {
//	ctx, flush := th.TestContext()
//	defer flush()
//
//	td := testdeep.NewT(t)
//	var mc = minimock.NewController(td)
//	defer mc.Finish()
//
//	type listenerAnswer struct {
//		c   net.Conn
//		err error
//	}
//	listenerAnswerChan := make(chan listenerAnswer, 1)
//	defer func() {
//		close(listenerAnswerChan)
//	}()
//
//	listener := NewListenerMock(mc)
//	listener.AcceptMock.Set(func() (c1 net.Conn, err error) {
//		ans := <-listenerAnswerChan
//		return ans.c, ans.err
//	})
//
//	proxy :=
//}
