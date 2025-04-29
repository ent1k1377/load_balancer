package main

import (
	"fmt"
	"github.com/ent1k1377/load_balancer/internal/balancer"
	"github.com/ent1k1377/load_balancer/internal/config"
	"github.com/ent1k1377/load_balancer/internal/logger"
	"github.com/ent1k1377/load_balancer/internal/ratelimiter"
	"github.com/ent1k1377/load_balancer/internal/server"
	"net/http"
	"time"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}
	logger.Infof("%+v", cfg)

	pool := balancer.InitServerPool(cfg)

	rateLimiter := ratelimiter.NewRateLimiter(cfg.DefaultCapacity, cfg.DefaultRefillRate, time.Duration(cfg.DefaultRefillPeriod)*time.Millisecond)

	loadBalancer := http.HandlerFunc(pool.LoadBalancer)
	handler := rateLimiter.Middleware(loadBalancer)

	server.Start(fmt.Sprintf(":%d", cfg.Port), handler)
}
