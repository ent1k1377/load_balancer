package balancer

import (
	"github.com/ent1k1377/load_balancer/internal/config"
	"github.com/ent1k1377/load_balancer/internal/logger"
)

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
