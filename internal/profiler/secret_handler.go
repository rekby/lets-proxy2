package profiler

import (
	"net"
	"net/http"

	"go.uber.org/zap"
)

type secretHandler struct {
	AllowedNetworks []net.IPNet
	logger          *zap.Logger
	next            http.Handler
}

func (s secretHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	remoteIP, err := net.ResolveTCPAddr("tcp", req.RemoteAddr)
	if err != nil {
		s.logger.Error("Parse remote address", zap.String("remote_addr", req.RemoteAddr), zap.Error(err))
		http.Error(resp, "Internal error", http.StatusInternalServerError)
		return
	}

	localLogger := s.logger.With(zap.Stringer("path", req.URL), zap.String("remote_addr", req.RemoteAddr))
	for _, subnet := range s.AllowedNetworks {
		if subnet.Contains(remoteIP.IP) {
			localLogger.Info("Profiler access")
			s.next.ServeHTTP(resp, req)
			return
		}
	}

	localLogger.Error("Profiler deny")
	http.Error(resp, "Forbidden", http.StatusForbidden)
}
