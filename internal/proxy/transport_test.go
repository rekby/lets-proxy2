package proxy

import (
	"net/http"
	"testing"
	"time"

	"github.com/rekby/lets-proxy2/internal/th"

	"github.com/maxatome/go-testdeep"
)

func TestTransport_GetTransport(t *testing.T) {
	ctx, flush := th.TestContext(t)
	defer flush()

	td := testdeep.NewT(t)

	tr := Transport{}
	r, _ := http.NewRequest(http.MethodGet, "http://www.ru", nil)
	r = r.WithContext(ctx)
	httpTransport := tr.getTransport(r)
	td.True(httpTransport == defaultHTTPTransport) // equal pointers

	tr = Transport{IgnoreHTTPSCertificate: false}
	r, _ = http.NewRequest(http.MethodGet, "https://www.ru", nil)
	r = r.WithContext(ctx)
	httpTransport = tr.getTransport(r)
	td.True(httpTransport != defaultTransport()) // different pointers
	td.Cmp(httpTransport.TLSClientConfig.ServerName, "www.ru")
	td.Cmp(httpTransport.TLSClientConfig.InsecureSkipVerify, false)

	tr = Transport{IgnoreHTTPSCertificate: true}
	r, _ = http.NewRequest(http.MethodGet, "https://www.ru", nil)
	r = r.WithContext(ctx)
	httpTransport = tr.getTransport(r)
	td.True(httpTransport != defaultTransport()) // different pointers
	td.Cmp(httpTransport.TLSClientConfig.ServerName, "www.ru")
	td.Cmp(httpTransport.TLSClientConfig.InsecureSkipVerify, true)
}

func TestTransport_RoundTrip(t *testing.T) {
	t.Run("should return status 429 when request is not allowed by rate limiter", func(t *testing.T) {
		td := testdeep.NewT(t)
		rateLimiter, err := NewRateLimiter(RateLimitParams{
			// -1 force rate limiter to never allow a request
			RateLimit:  -1,
			TimeWindow: time.Second,
			Burst:      1,
			CacheSize:  100,
		})
		td.CmpNoError(err)

		tr := Transport{RateLimiter: rateLimiter}
		req, _ := http.NewRequest(http.MethodGet, "http://www.ru", nil)

		resp, err := tr.RoundTrip(req)

		td.CmpNoError(err)
		td.Cmp(resp.StatusCode, http.StatusTooManyRequests, "should return '429 Too Many Request'")
	})
}
