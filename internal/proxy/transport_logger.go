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
		var respStatusCode int
		var respContentLength int64
		if resp != nil {
			respStatusCode = resp.StatusCode
			respContentLength = resp.ContentLength
		}
		log.InfoErrorCtx(request.Context(), err, "Request",
			zap.Duration("duration_without_body", time.Since(start)),
			zap.String("initiator_addr", request.RemoteAddr),
			zap.String("metod", request.Method),
			zap.String("host", request.Host),
			zap.String("path", request.URL.Path),
			zap.String("query", request.URL.RawQuery),
			zap.Int("status_code", respStatusCode),
			zap.Int64("request_content_length", request.ContentLength),
			zap.Int64("resp_content_length", respContentLength),
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
