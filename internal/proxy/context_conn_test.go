package proxy

import (
	"net"

	"github.com/rekby/lets-proxy2/internal/manager"
)

var (
	_ net.Conn           = ContextConnextion{}
	_ manager.GetContext = ContextConnextion{}
)
