package proxy

import (
	"net/http"
	"sync"

	"github.com/hashicorp/golang-lru/v2"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	rateLimit  int
	timeWindow int
	burst      int

	mx    sync.RWMutex
	cache *lru.Cache[string, *rate.Limiter]
}

func NewRateLimiter(rateLimit, timeWindow, burst, cacheSize int) (*RateLimiter, error) {
	if rateLimit == 0 {
		return &RateLimiter{}, nil
	}

	cache, err := lru.New[string, *rate.Limiter](cacheSize)
	if err != nil {
		return nil, err
	}

	return &RateLimiter{
		rateLimit:  rateLimit,
		timeWindow: timeWindow,
		burst:      burst,
		cache:      cache,
	}, nil
}

func (rl *RateLimiter) Allow(r *http.Request) bool {
	if rl.rateLimit == 0 {
		return true
	}
	if r.Context().Err() != nil {
		return false
	}

	return rl.getLimiter(r).Allow()
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
		limiter = rate.NewLimiter(rate.Limit(rl.rateLimit*1000/rl.timeWindow), rl.burst)
		rl.cache.Add(ip, limiter)
	}

	return limiter
}

func getIP(r *http.Request) string {
	return r.RemoteAddr
}
