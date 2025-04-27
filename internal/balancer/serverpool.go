package balancer

import (
	"github.com/ent1k1377/load_balancer/internal/logger"
	"net/http"
	"sync/atomic"
)

type StrategyFunc func(pool *ServerPool) *Backend

type ServerPool struct {
	backends []*Backend
	current  uint64
	strategy StrategyFunc
}

func NewServerPool(strategy StrategyFunc) *ServerPool {
	return &ServerPool{
		strategy: strategy,
	}
}

func (s *ServerPool) AddBackend(newBackend *Backend) {
	logger.Infof("Adding new backend: %s", newBackend.URL.String())
	s.backends = append(s.backends, newBackend)
}

func (s *ServerPool) LoadBalancer(w http.ResponseWriter, r *http.Request) {
	back := s.GetNextBackend()
	if back != nil {
		logger.Infof("Backend is %s", back.URL.String())
		back.ReverseProxy.ServeHTTP(w, r)
		return
	}

	logger.Error("All backends are not working")
	http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
}

func (s *ServerPool) GetNextBackend() *Backend {
	if s.strategy == nil {
		logger.Error("No backend strategy")
		return nil
	}

	if len(s.backends) == 0 {
		logger.Error("The backend pool is empty")
		return nil
	}

	return s.strategy(s)
}

func (s *ServerPool) NextIndex() int {
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(len(s.backends)))
}
