package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/transaction-tracker/price_service/api/routes"
	"github.com/transaction-tracker/price_service/internal/config"
	"github.com/transaction-tracker/price_service/internal/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Setup router
	router := routes.SetupRouter(cfg)

	// Setup server
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		// Use structured logging for server startup
		logger.Info("Starting price service", logger.H{
			"port": cfg.Server.Port,
			"gin_mode": os.Getenv("GIN_MODE"),
		})

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server", err, logger.H{"port": cfg.Server.Port})
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server", logger.H{"signal": "received"})

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", err, logger.H{})
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited gracefully", logger.H{})
}
