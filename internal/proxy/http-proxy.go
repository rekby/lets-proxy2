package proxy

import (
	"context"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/rekby/lets-proxy2/internal/contexthelper"

	"github.com/rekby/lets-proxy2/internal/contextlabel"

	"github.com/rekby/lets-proxy2/internal/log"

	"go.uber.org/zap"

	zc "github.com/rekby/zapcontext"
)

type Director interface {
	Director(request *http.Request) error
}

type HTTPProxy struct {
	GetContext           func(req *http.Request) (context.Context, error)
	HandleHTTPValidation func(w http.ResponseWriter, r *http.Request) bool
	Director             Director // modify requests to backend.
	HTTPTransport        http.RoundTripper
	EnableAccessLog      bool

	logger           *zap.Logger
	listener         net.Listener
	httpReverseProxy httputil.ReverseProxy
	IdleTimeout      time.Duration
	httpServer       http.Server
}

func NewHTTPProxy(ctx context.Context, listener net.Listener) *HTTPProxy {
	const defaultHTTPPort = 80
	res := &HTTPProxy{
		HandleHTTPValidation: func(_ http.ResponseWriter, _ *http.Request) bool {
			return false
		},
		Director:   NewDirectorSameIP(defaultHTTPPort),
		GetContext: getContext,
		listener:   listener,
		logger:     zc.L(ctx),
		httpServer: http.Server{},
	}
	res.httpReverseProxy.Director = res.director
	return res
}

func (p *HTTPProxy) Close() error {
	return p.httpServer.Close()
}

// Start - finish initialization of proxy and start handling request.
// It is sync method, always return with non nil error: if handle stopped by context or if error on start handling.
// Any public fields must not change after Start called
func (p *HTTPProxy) Start() error {
	if p.HTTPTransport != nil {
		p.logger.Info("Set transport to reverse proxy")
		p.httpReverseProxy.Transport = p.HTTPTransport
	}

	if p.EnableAccessLog {
		p.httpReverseProxy.Transport = NewTransportLogger(p.httpReverseProxy.Transport)
	}
	p.logger.Info("Access log", zap.Bool("enabled", p.EnableAccessLog))

	p.httpServer.Handler = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !p.HandleHTTPValidation(writer, request) {
			p.httpReverseProxy.ServeHTTP(writer, request)
		}
	})
	p.httpServer.IdleTimeout = p.IdleTimeout

	p.logger.Info("Http builtin reverse proxy start")
	err := p.httpServer.Serve(p.listener)
	return err
}

func getContext(_ *http.Request) (context.Context, error) {
	return zc.WithLogger(context.WithValue(context.Background(), contextlabel.ConnectionID, "conn-id-none"), zap.NewNop()), nil
}

func (p *HTTPProxy) director(request *http.Request) {
	ctx, err := p.GetContext(request)

	logger := zc.L(ctx)
	log.DebugDPanic(logger, err, "Get connection context for request")
	*request = *request.WithContext(contexthelper.CombineContext(ctx, request.Context()))

	if request.URL == nil {
		request.URL = &url.URL{}
	}
	err = p.Director.Director(request)
	log.DebugPanic(logger, err, "Apply directors")
}
