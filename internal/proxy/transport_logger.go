package proxy

import (
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/rekby/lets-proxy2/internal/log"
)

type TransportLogger struct {
	Transport http.RoundTripper
}

func (t TransportLogger) RoundTrip(request *http.Request) (resp *http.Response, err error) {
	start := time.Now()

	defer func() {
		log.InfoErrorCtx(request.Context(), err, "Request",
			zap.Duration("duration_without_body", time.Since(start)),
			zap.String("initiator_addr", request.RemoteAddr),
			zap.String("metod", request.Method),
			zap.String("host", request.Host),
			zap.String("path", request.URL.Path),
			zap.String("query", request.URL.RawQuery),
			zap.Int("status_code", resp.StatusCode),
			zap.Int64("request_content_length", request.ContentLength),
			zap.Int64("resp_content_length", resp.ContentLength),
		)
	}()

	return t.Transport.RoundTrip(request)
}

func NewTransportLogger(transport http.RoundTripper) TransportLogger {
	if transport == nil {
		return TransportLogger{Transport: http.DefaultTransport}
	}
	return TransportLogger{Transport: transport}
}
