package balancer

import (
	"fmt"
	"math"
	"sync"
)

// LeastConnectionsStrategy реализует стратегию балансировки нагрузки,
// при которой запрос направляется на backend с наименьшим числом активных подключений.
type LeastConnectionsStrategy struct {
	activeConnections map[*Backend]int // Количество активных подключений на каждый backend
	mux               sync.Mutex
}

func NewLeastConnections() *LeastConnectionsStrategy {
	return &LeastConnectionsStrategy{
		activeConnections: make(map[*Backend]int),
	}
}

// NextBackend выбирает backend с наименьшим числом активных подключений среди доступных.
// Возвращает также функцию release, которую нужно вызвать после завершения обработки запроса,
// чтобы освободить подключение, и ошибку в случае отсутствия доступных backend'ов.
func (s *LeastConnectionsStrategy) NextBackend(backends []*Backend) (*Backend, func(), error) {
	if len(backends) == 0 {
		return nil, emptyFunc, fmt.Errorf("no backends available")
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	var selected *Backend
	minConnections := math.MaxInt32

	for _, backend := range backends {
		activeConnection := s.activeConnections[backend]
		if backend.IsAlive() && activeConnection < minConnections {
			minConnections = activeConnection
			selected = backend
		}
	}

	s.activeConnections[selected]++

	return selected,
		func() {
			s.mux.Lock()
			s.activeConnections[selected]--
			s.mux.Unlock()
		},
		nil
}
