package balancer

import (
	"github.com/ent1k1377/load_balancer/internal/logger"
	"github.com/ent1k1377/load_balancer/internal/utils"
	"net/http"
	"sync"
	"time"
)

func NewServerPool(strategy Strategy) *ServerPool {
	return &ServerPool{
		strategy: strategy,
	}
}

func (s *ServerPool) LoadBalancer(w http.ResponseWriter, r *http.Request) {
	back, release, err := s.strategy.NextBackend(s.backends)
	if err != nil {
		logger.Errorf("All backends are not working: %s", err)
		utils.WriteJSONError(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer release()

	logger.Infof("Backend is %s", back.URL.String())
	back.ReverseProxy.ServeHTTP(w, r)
}

func (s *ServerPool) StartHealthChecks(interval int) {
	go func() {
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			s.healthCheck()
		}
	}()
}

func (s *ServerPool) addBackend(newBackend *Backend) {
	logger.Infof("Adding new backend: %s", newBackend.URL.String())
	s.backends = append(s.backends, newBackend)
}

func (s *ServerPool) healthCheck() {
	var wg sync.WaitGroup
	defer wg.Add(len(s.backends))

	for _, b := range s.backends {
		go func(backend *Backend) {
			defer wg.Done()
			alive := backend.checkAlive()
			backend.SetAlive(alive)
		}(b)
	}

	wg.Wait()
}
