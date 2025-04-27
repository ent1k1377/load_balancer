package main

import (
	"context"
	"errors"
	"github.com/ent1k1377/load_balancer/internal/balancer"
	"github.com/ent1k1377/load_balancer/internal/config"
	"github.com/ent1k1377/load_balancer/internal/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func createServerPool(cfg *config.Config) *balancer.ServerPool {
	pool := balancer.NewServerPool(balancer.RoundRobinStrategy)

	for _, rawURL := range cfg.Backends {
		backend, err := balancer.NewBackend(rawURL)
		if err != nil {
			logger.Errorf("Failed to create backend: %v", err)
			continue
		}
		pool.AddBackend(backend)
	}

	logger.Infof("Serverpool formed")
	return pool
}

func startServer(addr string, handler http.Handler) {
	server := http.Server{
		Addr:    addr,
		Handler: handler,
	}
	go gracefulShutdown(&server)

	logger.Infof("Starting server on addrerss %s", addr)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatalf("Failed to start server: %v", err)
	}

	logger.Infof("Server exited gracefully")
}

func gracefulShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	logger.Infof("Received shutdown signal: %v", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("Failed to gracefully shutdown server: %v", err)
	}
}
