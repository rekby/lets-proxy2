package tlslistener

import (
	"context"
	"net"

	"github.com/rekby/lets-proxy2/internal/log"
	zc "github.com/rekby/zapcontext"
	"go.uber.org/zap"
)

type Config struct {
	TLSAddresses  []string
	TCPAddresses  []string
	MinTLSVersion string
}

func (c Config) Apply(ctx context.Context, l *ListenersHandler) error {
	logger := zc.L(ctx)

	var tlsListeners = make([]net.Listener, 0, len(c.TLSAddresses))
	for _, addr := range c.TLSAddresses { //nolint:wsl
		listener, err := net.Listen("tcp", addr)
		log.DebugError(logger, err, "Start listen tls binding", zap.String("address", addr))
		if err != nil {
			return err
		}

		tlsListeners = append(tlsListeners, listener)
	}

	var tcpListeners = make([]net.Listener, 0, len(c.TCPAddresses))

	for _, addr := range c.TCPAddresses {
		listener, err := net.Listen("tcp", addr)
		log.DebugError(logger, err, "Start listen tcp binding", zap.String("address", addr))
		if err != nil {
			return err
		}

		tcpListeners = append(tcpListeners, listener)
	}
	l.ListenersForHandleTLS = tlsListeners
	l.Listeners = tcpListeners

	if tlsVersion, err := ParseTLSVersion(c.MinTLSVersion); err == nil {
		l.MinTLSVersion = tlsVersion
		logger.Info("Min tls version", zap.String("tls_version", c.MinTLSVersion))
	} else {
		return err
	}

	return nil
}
