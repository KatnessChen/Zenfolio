package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/transaction-tracker/backend/api/routes"
	"github.com/transaction-tracker/backend/config"
	"github.com/transaction-tracker/backend/internal/database"
)

func main() {
	// Load configuration
	config, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	log.Println("Initializing database...")
	dm, err := database.Initialize(config)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		log.Println("Closing database connection...")
		if err := dm.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	// Initialize router
	router := routes.SetupRouter(config)

	// Setup graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down gracefully...")
		os.Exit(0)
	}()

	// Start server
	log.Printf("Server starting on %s", config.ServerAddress)
	if err := router.Run(config.ServerAddress); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
