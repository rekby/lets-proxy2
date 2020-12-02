package proxy

import (
	"net/http"
)

var _ http.RoundTripper = TransportLogger{}
