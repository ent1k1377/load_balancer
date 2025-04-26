package main

import (
	"fmt"
	"github.com/ent1k1377/load_balancer/internal/balancer"
	"github.com/ent1k1377/load_balancer/internal/config"
	"github.com/ent1k1377/load_balancer/internal/logger"
	"net/http"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		logger.Fatalf("failed to load config: %v", err)
	}

	pool := balancer.NewServerPool(balancer.RoundRobinStrategy)
	for _, rawURL := range cfg.Backends {
		backend, err := balancer.NewBackend(rawURL)
		if err != nil {
			logger.Errorf("failed to create backend: %v", err)
			continue
		}
		pool.AddBackend(backend)
		logger.Infof("backend created: %v", backend)
	}

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: http.HandlerFunc(pool.LoadBalancer),
	}

	logger.Infof("starting server on port %d", cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		logger.Fatalf("failed to start server: %v", err)
	}
}
