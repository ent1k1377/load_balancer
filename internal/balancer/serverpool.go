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

// LoadBalancer — это обработчик HTTP-запросов, который использует стратегию для выбора подходящего backend'а.
// Он выбирает backend, направляет запрос на выбранный backend
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

// StartHealthChecks запускает периодические проверки состояния backend'ов.
// Задается интервал в секундах, через который будет выполняться проверка каждого backend'а.
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

// healthCheck выполняет проверку состояния каждого backend'а.
// Каждый backend проверяется на доступность, и его состояние обновляется.
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
