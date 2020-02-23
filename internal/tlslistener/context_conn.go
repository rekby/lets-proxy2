package tlslistener

import (
	"context"
	"net"

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
	go func() {
		if conn.connCloseHandler != nil {
			conn.connCloseHandler(xerrors.New("Leak connection"))
			conn.connCloseHandler = nil
		}
		zc.L(conn.Context).Warn("Leaked connection")
		_ = conn.Close()
	}()
}
