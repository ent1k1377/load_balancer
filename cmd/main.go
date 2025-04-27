package main

import (
	"fmt"
	"github.com/ent1k1377/load_balancer/internal/config"
	"github.com/ent1k1377/load_balancer/internal/logger"
	"net/http"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	pool := createServerPool(cfg)

	// rateLimiter := ratelimiter.NewRateLimiter(10, 1, time.Millisecond*500)

	loadBalancer := http.HandlerFunc(pool.LoadBalancer)
	// handler := rateLimiter.Middleware(loadBalancer)
	startServer(fmt.Sprintf(":%d", cfg.Port), loadBalancer)
}
