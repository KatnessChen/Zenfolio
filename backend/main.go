package main

import (
	"log"

	"github.com/transaction-tracker/backend/api/routes"
	"github.com/transaction-tracker/backend/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize router
	r := routes.SetupRouter(cfg)

	// Start server
	log.Printf("Server starting on %s", cfg.ServerAddress)
	if err := r.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
