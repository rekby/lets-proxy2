package proxy

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/rekby/lets-proxy2/internal/log"

	"go.uber.org/zap"

	zc "github.com/rekby/zapcontext"
)

const httpPort = 80

type HttpProxy struct {
	GetDestination       func(ctx context.Context, remoteAddr string) (addr string, err error)
	GetContext           func(req *http.Request) context.Context
	HandleHttpValidation func(w http.ResponseWriter, r *http.Request) bool

	listener         net.Listener
	httpReverseProxy httputil.ReverseProxy
}

func NewHttpProxy(ctx context.Context, listener net.Listener) *HttpProxy {
	res := &HttpProxy{
		GetDestination: getDestination,
		HandleHttpValidation: func(_ http.ResponseWriter, _ *http.Request) bool {
			return false
		},
		GetContext: getContext,
		listener:   listener,
	}
	res.httpReverseProxy.Director = res.director

	mux := &http.ServeMux{}
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if !res.HandleHttpValidation(writer, request) {
			res.httpReverseProxy.ServeHTTP(writer, request)
		}
	})
	httpServer := http.Server{}
	httpServer.Handler = mux

	go func() {
		<-ctx.Done()
		err := httpServer.Close()
		log.DebugErrorCtx(ctx, err, "Http builtin reverse proxy stop becouse context cancelled")
	}()

	go func() {
		zc.L(ctx).Info("Http builtin reverse proxy start")
		err := httpServer.Serve(res.listener)
		if err == http.ErrServerClosed {
			err = nil
		}
		log.DebugErrorCtx(ctx, err, "Http builtin reverse proxy stop")
	}()

	return res
}

func (p *HttpProxy) SetTransport(transport http.RoundTripper) {
	p.httpReverseProxy.Transport = transport
}

func getDestination(_ context.Context, remoteAddr string) (addr string, err error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", remoteAddr)
	if err != nil {
		return "", errors.New("default get destination accept only tcp addresses")
	}
	tcpAddr.Port = 80
	return tcpAddr.String(), nil
}

func getContext(_ *http.Request) context.Context {
	return zc.WithLogger(context.Background(), zap.NewNop())
}

func (p *HttpProxy) director(request *http.Request) {
	ctx := p.GetContext(request)
	if request.URL == nil {
		request.URL = &url.URL{}
	}
	dest, err := p.GetDestination(ctx, request.RemoteAddr)
	log.DebugErrorCtx(ctx, err, "Get destination", zap.String("dest_addr", dest))
	request.URL.Scheme = "http"
	request.URL.Host = dest // If err != nil and dest invalid - is ok, becouse it will error proxy
}
