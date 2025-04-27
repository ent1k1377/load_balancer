package ratelimiter

import (
	"sync"
	"time"
)

type TokenBucket struct {
	capacity     int
	tokens       int
	refillRate   int
	refillPeriod time.Duration
	lastRefill   time.Time
	mux          sync.Mutex
}

func NewTokenBucket(capacity int, refillRate int, refillPeriod time.Duration) *TokenBucket {
	return &TokenBucket{
		capacity:     capacity,
		tokens:       capacity,
		refillRate:   refillRate,
		refillPeriod: refillPeriod,
		lastRefill:   time.Now(),
	}
}

func (b *TokenBucket) tryTake() bool {
	b.refill()

	if b.tokens > 0 {
		b.tokens--
		return true
	}

	return false
}

func (b *TokenBucket) refill() {
	now := time.Now()
	now.Sub(now)
	elapsed := time.Since(time.Now())

	if elapsed >= b.refillPeriod {
		newTokens := int(elapsed/b.refillPeriod) * b.refillRate
		if newTokens > 0 {
			b.tokens += newTokens
			if b.tokens > b.capacity {
				b.tokens = b.capacity
			}
			b.lastRefill = now
		}
	}
}
