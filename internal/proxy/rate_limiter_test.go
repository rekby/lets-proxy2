package proxy

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep"
)

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
			timeWindow: 1000,
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
			timeWindow: 500,
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
			timeWindow: 1000,
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
			timeWindow: 1000,
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Preparations
			limiter, err := NewRateLimiter(RateLimitParams{
				RateLimit:  tt.rateLimit,
				TimeWindow: tt.timeWindow,
				Burst:      1,
				CacheSize:  100,
			})
			testdeep.CmpNoError(t, err)

			startTime := time.Now()
			successCounters := make([]int, len(tt.reqSpecs))
			reqCounters := make([]int, len(tt.reqSpecs))

			// The test itself
			for time.Since(startTime) < tt.testTime {
				for idx, spec := range tt.reqSpecs {
					reqCounters[idx]++
					if limiter.Allow(spec.req) {
						successCounters[idx]++
					}
				}
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
