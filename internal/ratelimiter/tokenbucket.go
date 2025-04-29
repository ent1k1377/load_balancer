package ratelimiter

import (
	"sync"
	"time"
)

// TokenBucket представляет собой структуру для ограничения количества запросов с использованием алгоритма "токенное ведро".
type TokenBucket struct {
	// capacity - максимальное количество токенов в ведре.
	capacity int
	// tokens - текущее количество токенов в ведре.
	tokens int
	// refillRate - количество токенов, которые добавляются в ведро за каждый refillPeriod.
	refillRate int
	// refillPeriod - период времени, через который обновляется количество токенов в ведре.
	refillPeriod time.Duration
	// lastRefill - время последнего пополнения токенов.
	lastRefill time.Time
	mux        sync.Mutex
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

// tryTake пытается забрать токен из ведра. Если токен доступен, уменьшает их количество на 1 и возвращает true.
// Если токенов нет, возвращает false.
func (b *TokenBucket) tryTake() bool {
	b.refill()
	b.mux.Lock()
	defer b.mux.Unlock()

	if b.tokens > 0 {
		b.tokens--
		return true
	}

	return false
}

// refill пополняет ведро токенами, если прошел достаточный период времени с последнего пополнения.
func (b *TokenBucket) refill() {
	b.mux.Lock()
	defer b.mux.Unlock()

	now := time.Now()
	elapsed := now.Sub(b.lastRefill)

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
