package proxy

import (
	"net"

	"github.com/rekby/lets-proxy2/internal/cert_manager"
)

var (
	_ net.Conn                = ContextConnextion{}
	_ cert_manager.GetContext = ContextConnextion{}
)
