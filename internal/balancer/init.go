package balancer

import (
	"github.com/ent1k1377/load_balancer/internal/config"
	"github.com/ent1k1377/load_balancer/internal/logger"
)

// InitServerPool инициализирует ServerPool с выбранной стратегией
// и добавляет в него список backend'ов из конфигурации. Также запускает health check.
func InitServerPool(cfg *config.Config) *ServerPool {
	strategy := strategyByName(cfg.Strategy)
	pool := NewServerPool(strategy)

	for _, rawURL := range cfg.Backends {
		if backend, err := NewBackend(rawURL); err == nil {
			pool.addBackend(backend)
		} else {
			logger.Errorf("Failed to create backend: %v", err)
		}
	}

	go pool.StartHealthChecks(cfg.HealthCheckInterval)
	return pool
}

// strategyByName возвращает реализацию Strategy по её имени из конфигурации.
// Поддерживаемые значения: "round_robin", "least_connections".
func strategyByName(name string) Strategy {
	switch name {
	case "round_robin":
		return NewRoundRobinStrategy()
	case "least_connections":
		return NewLeastConnections()
	default:
		logger.Fatalf("Unknown strategy: %s", name)
		return nil
	}
}
