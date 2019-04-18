package proxy

import (
	"context"
	"net"
)

type ContextConnextion struct {
	net.Conn
	context.Context
}

func (c ContextConnextion) GetContext() context.Context {
	return c.Context
}
