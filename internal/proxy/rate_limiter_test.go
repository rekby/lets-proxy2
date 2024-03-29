package proxy

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/maxatome/go-testdeep"
)

func TestRateLimiter_Allow(t *testing.T) {
	t.Run("should trigger an error if cache size is less than zero", func(t *testing.T) {
		_, err := NewRateLimiter(RateLimitParams{RateLimit: 1, CacheSize: -1})
		testdeep.CmpError(t, err)
	})

	t.Run("should always allow if rate limit is zero", func(t *testing.T) {
		td := testdeep.NewT(t)
		rateLimiter, err := NewRateLimiter(RateLimitParams{RateLimit: 0})
		req, _ := http.NewRequest(http.MethodGet, "http://url.com", nil)

		td.CmpNoError(err)
		td.True(rateLimiter.Allow(req))
	})
}

func TestMaxRequestsPerSec(t *testing.T) {
	req1, _ := http.NewRequest("GET", "http://url1.com", nil)
	req1.RemoteAddr = "ip1"

	req2, _ := http.NewRequest("GET", "http://url2.com", nil)
	req2.RemoteAddr = "ip2"

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	canceledReq, _ := http.NewRequestWithContext(ctx, "GET", "http://canceled.com", nil)
	canceledReq.RemoteAddr = "canceled"

	type reqSpec struct {
		req                      *http.Request
		wantAllowedRequestAround int // expected the limiter to allow +-1
	}

	tests := []struct {
		name string

		rateLimit  int
		timeWindow time.Duration
		testTime   time.Duration

		reqSpecs []reqSpec
	}{
		{
			name: "should limit the amount of requests per second",

			rateLimit:  10,
			timeWindow: time.Second,
			testTime:   time.Second,

			reqSpecs: []reqSpec{
				{
					req:                      req1,
					wantAllowedRequestAround: 10,
				},
			},
		},
		{
			name: "should restart the timer for the next time window",

			rateLimit:  10,
			timeWindow: 500 * time.Millisecond,
			testTime:   time.Second,

			reqSpecs: []reqSpec{
				{
					req:                      req1,
					wantAllowedRequestAround: 20,
				},
			},
		},
		{
			name: "requests from different IPs should NOT influence each other",

			rateLimit:  10,
			timeWindow: time.Second,
			testTime:   time.Second,

			reqSpecs: []reqSpec{
				{
					req:                      req1,
					wantAllowedRequestAround: 10,
				},
				{
					req:                      req2,
					wantAllowedRequestAround: 10,
				},
			},
		},
		{
			name: "canceled request should always fail",

			rateLimit:  10,
			timeWindow: time.Second,
			testTime:   time.Second,

			reqSpecs: []reqSpec{
				{
					req:                      canceledReq,
					wantAllowedRequestAround: 0,
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Preparations
			startTime := time.Now()
			clock := clockwork.NewFakeClockAt(startTime)
			limiter, err := NewRateLimiter(RateLimitParams{
				RateLimit:  tt.rateLimit,
				TimeWindow: tt.timeWindow,
				Burst:      1,
				CacheSize:  100,
				Clock:      clock,
			})
			testdeep.CmpNoError(t, err)

			successCounters := make([]int, len(tt.reqSpecs))
			reqCounters := make([]int, len(tt.reqSpecs))

			// The test itself
			for clock.Since(startTime) < tt.testTime {
				for idx, spec := range tt.reqSpecs {
					reqCounters[idx]++
					if limiter.Allow(spec.req) {
						successCounters[idx]++
					}
				}
				clock.Advance(time.Millisecond)
			}

			// Check the expectations
			for idx, spec := range tt.reqSpecs {
				testdeep.CmpBetween(
					t,
					successCounters[idx],
					spec.wantAllowedRequestAround-1,
					spec.wantAllowedRequestAround+1,
					testdeep.BoundsInIn,
				)
				testdeep.CmpGt(t, reqCounters[idx], successCounters[idx])
			}
		})
	}
}
