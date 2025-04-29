package server

import (
	"context"
	"errors"
	"github.com/ent1k1377/load_balancer/internal/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Start запускает HTTP сервер и начинает слушать на указанном адресе.
func Start(addr string, handler http.Handler) {
	server := &http.Server{Addr: addr, Handler: handler}
	go gracefulShutdown(server)

	logger.Infof("Starting server on %s", addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatalf("Server error: %v", err)
	}
}

// gracefulShutdown обрабатывает завершение работы сервера по сигналам ОС (SIGINT, SIGTERM).
// Он ожидает сигнал завершения, а затем корректно закрывает сервер.
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
