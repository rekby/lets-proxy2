package tlslistener

import (
	"context"
	"net"

	zc "github.com/rekby/zapcontext"
)

type ContextConnextion struct {
	net.Conn
	context.Context
	CloseFunc func() error
}

func (c ContextConnextion) GetContext() context.Context {
	return c.Context
}

func (c ContextConnextion) Close() error {
	if c.CloseFunc == nil {
		return c.Conn.Close()
	}
	return c.CloseFunc()
}

func finalizeContextConnection(conn *ContextConnextion) {
	go func() {
		zc.L(conn.Context).Warn("Leaked connection")
		_ = conn.Close()
	}()
}
