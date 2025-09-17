package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/transaction-tracker/backend/internal/logger"

	"github.com/transaction-tracker/backend/api/routes"
	"github.com/transaction-tracker/backend/config"
	"github.com/transaction-tracker/backend/internal/database"
)

func main() {

	// Initialize structured logger
	logger.InitLogger()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load configuration", err, logger.H{})
		os.Exit(1)
	}

	// Initialize database
	logger.Info("Initializing database...")
	dm, err := database.Initialize(cfg)
	if err != nil {
		logger.Error("Failed to initialize database", err, logger.H{})
		os.Exit(1)
	}
	defer func() {
		logger.Info("Closing database connection...")
		if err := dm.Close(); err != nil {
			logger.Info("Backend server started", logger.H{"port": cfg.ServerAddress})
		}
	}()

	// Initialize router
	router := routes.SetupRouter(cfg)

	// Setup graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		logger.Info("Shutting down gracefully...")
		os.Exit(0)
	}()

	// Start server
	logger.Info("Server starting", logger.H{"address": cfg.ServerAddress})
	if err := router.Run(cfg.ServerAddress); err != nil {
		logger.Error("Failed to start server", err, logger.H{"address": cfg.ServerAddress})
		os.Exit(1)
	}
}
