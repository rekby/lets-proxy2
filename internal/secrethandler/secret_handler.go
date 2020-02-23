package secrethandler

import (
	"net"
	"net/http"
	"net/url"

	"github.com/rekby/lets-proxy2/internal/log"

	"go.uber.org/zap"
)

const (
	errAccessDeniedMess = "Access denied"
	maxURLLen           = 100
	passwordArgName     = "password"
)

type Config struct {
	AllowedNetworks    []string
	Password           string
	AllowEmptyPassword bool
}

type SecretHandler struct {
	allowedNetworks    []net.IPNet
	allowEmptyPassword bool
	password           string
	logger             *zap.Logger
	next               http.Handler
}

func (m SecretHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := m.logger.With(zap.Stringer("path", r.URL), zap.String("remote_address", r.RemoteAddr))
	if len(m.allowedNetworks) > 0 {
		remoteIP, err := net.ResolveTCPAddr("tcp", r.RemoteAddr)
		if err != nil {
			m.logger.Error("Parse remote address", zap.String("remote_addr", r.RemoteAddr), zap.Error(err))
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		networkAllow := false
		for i := range m.allowedNetworks {
			subnet := m.allowedNetworks[i]
			contains := subnet.Contains(remoteIP.IP)
			logger.Debug("Check contains", zap.Stringer("subnet", &subnet), zap.Bool("contains", contains))
			if contains {
				logger.Debug("Allow by subnet", zap.Stringer("subnet", &subnet))
				networkAllow = true
				break
			}
		}
		if !networkAllow {
			logger.Error("Deny by remote address")
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
	}

	if m.password == "" && !m.allowEmptyPassword {
		http.Error(w, errAccessDeniedMess, http.StatusForbidden)
		return
	}

	if len(r.URL.RawQuery) > maxURLLen {
		http.Error(w, "Very long url", http.StatusRequestURITooLong)
		return
	}

	values, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, "Bad query", http.StatusBadRequest)
		return
	}

	pass := values.Get(passwordArgName)

	if pass != m.password || !m.allowEmptyPassword && pass == "" {
		http.Error(w, "Bad password", http.StatusForbidden)
		return
	}

	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "Bad method", http.StatusMethodNotAllowed)
		return
	}

	logger.Info("Allow remote access")
	m.next.ServeHTTP(w, r)
}

func New(logger *zap.Logger, config Config, next http.Handler) SecretHandler {
	localLogger := logger.Named("create_secret_handler")
	var allowedNetworksIP []net.IPNet
	for _, network := range config.AllowedNetworks {
		_, parsedNet, err := net.ParseCIDR(network)
		log.InfoError(localLogger, err, "Parse allowed CIDR", zap.Stringer("network", parsedNet))
		if err == nil {
			allowedNetworksIP = append(allowedNetworksIP, *parsedNet)
		}
	}

	var secretHandler SecretHandler
	secretHandler.allowedNetworks = allowedNetworksIP
	secretHandler.password = config.Password
	secretHandler.allowEmptyPassword = config.AllowEmptyPassword
	secretHandler.next = next
	secretHandler.logger = logger
	return secretHandler
}
