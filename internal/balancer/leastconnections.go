package balancer

import (
	"fmt"
	"math"
	"sync"
)

type LeastConnectionsStrategy struct {
	activeConnections map[*Backend]int
	mux               sync.Mutex
}

func NewLeastConnections() *LeastConnectionsStrategy {
	return &LeastConnectionsStrategy{
		activeConnections: make(map[*Backend]int),
	}
}

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
