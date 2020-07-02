package tlslistener

import (
	"context"
	"net"

	"github.com/rekby/lets-proxy2/internal/log"

	"golang.org/x/xerrors"

	"github.com/rekby/lets-proxy2/internal/metrics"

	zc "github.com/rekby/zapcontext"
)

type ContextConnextion struct {
	net.Conn
	context.Context
	CloseFunc        func() error
	connCloseHandler metrics.ProcessFinishFunc
}

func (c ContextConnextion) GetContext() context.Context {
	return c.Context
}

func (c ContextConnextion) Close() error {
	if c.connCloseHandler != nil {
		c.connCloseHandler(nil)
		c.connCloseHandler = nil
	}
	if c.CloseFunc == nil {
		return c.Conn.Close()
	}
	return c.CloseFunc()
}

func finalizeContextConnection(conn *ContextConnextion) {
	logger := zc.L(conn.Context)
	defer log.HandlePanic(logger)

	if conn.connCloseHandler != nil {
		conn.connCloseHandler(xerrors.New("Leak connection"))
		conn.connCloseHandler = nil
	}
	logger.Warn("Leaked connection")
	_ = conn.Close()
}
