package profiler

import (
	"net"
	"net/http"
	"net/http/pprof"

	"go.uber.org/zap"

	"github.com/rekby/lets-proxy2/internal/log"
)

type Config struct {
	Enable          bool
	AllowedNetworks []string
	BindAddress     string
}

type Profiler struct {
	secretHandler secretHandler
}

func New(logger *zap.Logger, config Config) *Profiler {
	mux := http.NewServeMux()

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	var allowedNetworks []net.IPNet
	for _, network := range config.AllowedNetworks {
		_, parsedNet, err := net.ParseCIDR(network)
		log.InfoError(logger, err, "Parse allowed CIDR", zap.Stringer("network", parsedNet))
		if err == nil {
			allowedNetworks = append(allowedNetworks, *parsedNet)
		}
	}

	handler := secretHandler{
		logger:          logger,
		next:            mux,
		AllowedNetworks: allowedNetworks,
	}
	return &Profiler{
		secretHandler: handler,
	}
}

func (p Profiler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	p.secretHandler.ServeHTTP(resp, req)
}
