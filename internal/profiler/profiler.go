package profiler

import (
	"net/http"
	"net/http/pprof"

	"github.com/rekby/lets-proxy2/internal/secrethandler"

	"go.uber.org/zap"
)

type Config struct {
	secrethandler.Config

	Enable      bool
	BindAddress string
}

type Profiler struct {
	secretHandler secrethandler.SecretHandler
}

func New(logger *zap.Logger, config Config) *Profiler {
	mux := http.NewServeMux()

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	return &Profiler{
		secretHandler: secrethandler.New(logger, config.Config, mux),
	}
}

func (p Profiler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	p.secretHandler.ServeHTTP(resp, req)
}
