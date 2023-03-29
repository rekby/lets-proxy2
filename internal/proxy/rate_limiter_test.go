package proxy

import (
	"context"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/maxatome/go-testdeep"
)

func TestRateLimiter_Wait(t *testing.T) {
	req1, _ := http.NewRequest("GET", "http://url1.com", nil)
	req1.RemoteAddr = "ip1"

	req2, _ := http.NewRequest("GET", "http://url2.com", nil)
	req1.RemoteAddr = "ip2"

	ctx, close := context.WithCancel(context.Background())
	close()
	reqWithClosedCtx, _ := http.NewRequestWithContext(ctx, "GET", "http://with-closed-ctx.com", nil)
	reqWithClosedCtx.RemoteAddr = "with-closed-ctx"

	type want struct {
		errorNumber int32
		cacheLen    int
	}
	tests := []struct {
		name string

		rateLimit  int
		timeWindow int
		burst      int
		cacheSize  int

		requests []*http.Request

		want want
	}{
		{
			name: "queries with the same IP should be handled sequentially",

			rateLimit:  1,
			timeWindow: 1,
			burst:      1,
			cacheSize:  100,

			requests: []*http.Request{req1, req1, req1},

			want: want{
				errorNumber: 0,
				cacheLen:    1,
			},
		},
		{
			name: "queries with different IPs should be handled in parallel",

			rateLimit:  1,
			timeWindow: 1,
			burst:      1,
			cacheSize:  100,

			requests: []*http.Request{req1, req1, req2},

			want: want{
				errorNumber: 0,
				cacheLen:    2,
			},
		},
		{
			name: "waiting should fail in case the context is done",

			rateLimit:  1,
			timeWindow: 1,
			burst:      1,
			cacheSize:  100,

			requests: []*http.Request{req1, reqWithClosedCtx, req2},

			want: want{
				errorNumber: 1,
				cacheLen:    3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			limiter, err := NewRateLimiter(tt.rateLimit, tt.timeWindow, tt.burst, tt.cacheSize)
			testdeep.CmpNoError(t, err)

			var errCounter int32 = 0
			var wg sync.WaitGroup

			for _, req := range tt.requests {
				req := req
				wg.Add(1)
				go func() {
					err := limiter.Wait(req)
					if err != nil {
						atomic.AddInt32(&errCounter, 1)
					}
					wg.Done()
				}()
			}

			wg.Wait()

			testdeep.Cmp(t, limiter.cache.Len(), tt.want.cacheLen, "incorrect cache length")
			testdeep.Cmp(t, errCounter, tt.want.errorNumber, "incorrect number of errors")
		})
	}
}
