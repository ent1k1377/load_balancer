package balancer

import (
	"github.com/ent1k1377/load_balancer/internal/logger"
	"sync/atomic"
)

func RoundRobinStrategy(pool *ServerPool) *Backend {
	next := pool.NextIndex()
	l := len(pool.backends) + next

	for i := next; i < l; i++ {
		idx := i % len(pool.backends)
		if pool.backends[idx].IsAlive() {
			if i != next {
				atomic.StoreUint64(&pool.current, uint64(idx))
			}

			return pool.backends[idx]
		}
	}

	logger.Warn("RoundRobinStrategy: no alive backends available")
	return nil
}
