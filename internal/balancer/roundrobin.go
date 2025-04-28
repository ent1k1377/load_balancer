package balancer

import (
	"errors"
	"github.com/ent1k1377/load_balancer/internal/logger"
	"sync/atomic"
)

type RoundRobinStrategy struct {
	current uint64
}

func NewRoundRobinStrategy() *RoundRobinStrategy {
	return &RoundRobinStrategy{
		current: 0,
	}
}

func (s *RoundRobinStrategy) NextBackend(backends []*Backend) (*Backend, func(), error) {
	if len(backends) == 0 {
		return nil, emptyFunc, errors.New("no available backend")
	}

	backend := s.choose(backends)
	if backend == nil {
		return nil, emptyFunc, errors.New("no available backend")
	}

	return backend, emptyFunc, nil
}

func (s *RoundRobinStrategy) choose(backends []*Backend) *Backend {
	next := s.nextIndex(len(backends))
	l := len(backends) + next

	for i := next; i < l; i++ {
		idx := i % len(backends)
		if backends[idx].IsAlive() {
			if i != next {
				atomic.StoreUint64(&s.current, uint64(idx))
			}

			return backends[idx]
		}
	}

	logger.Warn("RoundRobinStrategy: no alive backends available")
	return nil
}

func (s *RoundRobinStrategy) nextIndex(amountBackends int) int {
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(amountBackends))
}
