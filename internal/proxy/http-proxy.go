package proxy

import (
	"context"
	"net"
	"net/http"
	"net/http/httputil"

	"github.com/rekby/lets-proxy2/internal/contexthelper"

	"github.com/rekby/lets-proxy2/internal/contextlabel"

	"github.com/rekby/lets-proxy2/internal/log"

	"go.uber.org/zap"

	zc "github.com/rekby/zapcontext"
)

type Director interface {
	Director(request *http.Request)
}

type HTTPProxy struct {
	GetContext           func(req *http.Request) (context.Context, error)
	HandleHTTPValidation func(w http.ResponseWriter, r *http.Request) bool
	Director             Director // modify requests to backend.
	HTTPTransport        http.RoundTripper

	ctx              context.Context
	listener         net.Listener
	httpReverseProxy httputil.ReverseProxy
}

func NewHTTPProxy(ctx context.Context, listener net.Listener) *HTTPProxy {
	res := &HTTPProxy{
		HandleHTTPValidation: func(_ http.ResponseWriter, _ *http.Request) bool {
			return false
		},
		Director:   NewDirectorSameIP(80),
		GetContext: getContext,
		listener:   listener,
		ctx:        ctx,
	}
	res.httpReverseProxy.Director = res.director

	return res
}

// Start - finish initialization of proxy and start handling request.
// It is sync method, always return with non nil error: if handle stopped by context or if error on start handling.
// Any public fields must not change after Start called
func (p *HTTPProxy) Start() error {
	if p.HTTPTransport != nil {
		p.httpReverseProxy.Transport = p.HTTPTransport
	}

	mux := &http.ServeMux{}
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if !p.HandleHTTPValidation(writer, request) {
			p.httpReverseProxy.ServeHTTP(writer, request)
		}
	})
	httpServer := http.Server{}
	httpServer.Handler = mux

	go func() {
		<-p.ctx.Done()
		err := httpServer.Close()
		log.DebugErrorCtx(p.ctx, err, "Http builtin reverse proxy stop because context cancelled")
	}()

	zc.L(p.ctx).Info("Http builtin reverse proxy start")
	err := httpServer.Serve(p.listener)
	return err
}

func getContext(_ *http.Request) (context.Context, error) {
	return zc.WithLogger(context.WithValue(context.Background(), contextlabel.ConnectionID, "conn-id-none"), zap.NewNop()), nil
}

func (p *HTTPProxy) director(request *http.Request) {
	ctx, err := p.GetContext(request)
	log.DebugDPanicCtx(ctx, err, "Get connection context for request")
	*request = *request.WithContext(contexthelper.CombineContext(ctx, request.Context()))
	p.Director.Director(request)
}
