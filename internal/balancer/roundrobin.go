package balancer

import (
	"errors"
	"github.com/ent1k1377/load_balancer/internal/logger"
	"sync/atomic"
)

// RoundRobinStrategy реализует стратегию круговой балансировки нагрузки (Round Robin),
// при которой запросы направляются поочередно на доступные backend'ы.
type RoundRobinStrategy struct {
	current uint64 // Индекс текущего выбранного backend'а
}

func NewRoundRobinStrategy() *RoundRobinStrategy {
	return &RoundRobinStrategy{
		current: 0,
	}
}

// NextBackend выбирает следующий доступный backend из списка backends.
// Возвращает выбранный backend и пустую функцию для освобождения ресурса (не используется в данной стратегии),
// а также ошибку, если доступных backend'ов нет.
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

// choose выбирает backend, который будет обработан следующим. Алгоритм использует
// текущий индекс и проверяет доступность каждого backend'а по очереди.
// Если backend не доступен, продолжается проверка следующих до тех пор, пока не будет найден живой backend.
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

// nextIndex вычисляет индекс следующего backend'а для выбора.
// Индекс вычисляется как остаток от деления текущего индекса на количество доступных backend'ов.
// Это обеспечивает круговую балансировку (Round Robin).
func (s *RoundRobinStrategy) nextIndex(amountBackends int) int {
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(amountBackends))
}
