package proxy

import (
	"crypto/tls"
	"net"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	zc "github.com/rekby/zapcontext"
)

// copy from go 1.10, need for compile with go 1.10 compiler
// https://github.com/golang/go/blob/b0cb374daf646454998bac7b393f3236a2ab6aca/src/net/http/transport.go#L40
var defaultTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}).DialContext,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

type Transport struct {
	IgnoreHttpsCertificate bool
}

func (t Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.getTransport(req).RoundTrip(req)
}

func (t Transport) getTransport(req *http.Request) *http.Transport {
	logger := zc.L(req.Context())

	if req.URL.Scheme == ProtocolHTTP {
		logger.Debug("Use default http transport")
		return defaultTransport
	}

	host := req.Host
	if strings.Contains(host, ":") { //strip port
		parts := strings.SplitN(host, ":", 2)
		host = parts[0]
	}

	var transport http.Transport
	transport = *defaultTransport
	transport.TLSClientConfig = &tls.Config{ServerName: host}
	transport.TLSClientConfig.InsecureSkipVerify = t.IgnoreHttpsCertificate

	logger.Debug("Use https transport",
		zap.Bool("ignore_cert", transport.TLSClientConfig.InsecureSkipVerify),
		zap.String("tls_server_name", host),
		zap.String("header_host", req.Header.Get("HOST")),
	)

	return &transport
}
