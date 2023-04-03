package proxy

import (
	"net/http"
	"sync"
	"time"

	"github.com/hashicorp/golang-lru/v2"
	"github.com/jonboulle/clockwork"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	rateLimit  int
	timeWindow time.Duration
	burst      int

	clock clockwork.Clock
	mx    sync.RWMutex
	cache *lru.Cache[string, *rate.Limiter]
}

type RateLimitParams struct {
	RateLimit  int
	TimeWindow time.Duration
	Burst      int
	CacheSize  int
	Clock      clockwork.Clock
}

func NewRateLimiter(params RateLimitParams) (*RateLimiter, error) {
	if params.RateLimit == 0 {
		return &RateLimiter{}, nil
	}

	cache, err := lru.New[string, *rate.Limiter](params.CacheSize)
	if err != nil {
		return nil, err
	}

	self := &RateLimiter{
		rateLimit:  params.RateLimit,
		timeWindow: params.TimeWindow,
		burst:      params.Burst,
		cache:      cache,
		clock:      params.Clock,
	}

	if self.clock == nil {
		self.clock = clockwork.NewRealClock()
	}

	return self, nil
}

func (rl *RateLimiter) Allow(r *http.Request) bool {
	if rl.rateLimit == 0 {
		return true
	}
	if r.Context().Err() != nil {
		return false
	}

	return rl.getLimiter(r).AllowN(rl.clock.Now(), 1)
}

func (rl *RateLimiter) getLimiter(r *http.Request) *rate.Limiter {
	rl.mx.RLock()
	ip := getIP(r)

	limiter, ok := rl.cache.Get(ip)
	if ok {
		rl.mx.RUnlock()
		return limiter
	}

	rl.mx.RUnlock()
	rl.mx.Lock()
	defer rl.mx.Unlock()

	// we need to check cache again to avoid data race
	limiter, ok = rl.cache.Get(ip)
	if !ok {
		limiter = rate.NewLimiter(rate.Limit(float64(rl.rateLimit)/rl.timeWindow.Seconds()), rl.burst)
		rl.cache.Add(ip, limiter)
	}

	return limiter
}

func getIP(r *http.Request) string {
	return r.RemoteAddr
}
