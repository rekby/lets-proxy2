package tlslistener

import (
	"context"
	"net"
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
	} else {
		return c.CloseFunc()
	}
}
