package ratelimiter

import (
	"github.com/ent1k1377/load_balancer/internal/utils"
	"net/http"
	"strings"
	"sync"
	"time"
)

type RateLimiter struct {
	buckets             map[string]*TokenBucket
	defaultCapacity     int
	defaultRefillRate   int
	defaultRefillPeriod time.Duration
	mux                 sync.RWMutex
}

func NewRateLimiter(defaultCapacity, defaultRefillRate int, defaultRefillPeriod time.Duration) *RateLimiter {
	rl := &RateLimiter{
		buckets:             make(map[string]*TokenBucket),
		defaultCapacity:     defaultCapacity,
		defaultRefillRate:   defaultRefillRate,
		defaultRefillPeriod: defaultRefillPeriod,
	}

	return rl
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		addr := r.RemoteAddr
		ip := strings.Split(addr, ":")[0]
		if rl.Allow(ip) {
			next.ServeHTTP(w, r)
			return
		}

		utils.WriteJSONError(w, "Too many requests", http.StatusTooManyRequests)
	})
}

func (rl *RateLimiter) Allow(client string) bool {
	rl.mux.Lock()
	defer rl.mux.Unlock()
	bucket, exists := rl.buckets[client]

	if !exists {
		bucket = NewTokenBucket(rl.defaultCapacity, rl.defaultRefillRate, rl.defaultRefillPeriod)
		rl.buckets[client] = bucket
	}

	return bucket.tryTake()
}
