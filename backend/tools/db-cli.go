package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/transaction-tracker/backend/config"
	"github.com/transaction-tracker/backend/internal/database"
	"github.com/transaction-tracker/backend/migrations"
)

func main() {
	var (
		action = flag.String("action", "", "Action to perform: migrate, rollback, status, seed, health")
		env    = flag.String("env", "development", "Environment: development, staging, production")
	)
	flag.Parse()

	if *action == "" {
		printUsage()
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	dm, err := database.Initialize(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dm.Close()

	switch *action {
	case "migrate":
		if err := runMigrations(dm); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
	case "rollback":
		if err := rollbackMigration(dm); err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
	case "status":
		if err := showMigrationStatus(dm); err != nil {
			log.Fatalf("Status check failed: %v", err)
		}
	case "seed":
		if err := seedDatabase(dm, *env); err != nil {
			log.Fatalf("Seeding failed: %v", err)
		}
	case "health":
		if err := checkHealth(dm); err != nil {
			log.Fatalf("Health check failed: %v", err)
		}
	default:
		fmt.Printf("Unknown action: %s\n", *action)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Database CLI Tool")
	fmt.Println("Usage: go run tools/db-cli.go -action=<action> [-env=<environment>]")
	fmt.Println()
	fmt.Println("Actions:")
	fmt.Println("  migrate   - Run pending database migrations")
	fmt.Println("  rollback  - Rollback the last migration")
	fmt.Println("  status    - Show migration status")
	fmt.Println("  seed      - Seed database with sample data")
	fmt.Println("  health    - Check database health")
	fmt.Println()
	fmt.Println("Environments:")
	fmt.Println("  development - Development environment (default)")
	fmt.Println("  staging     - Staging environment")
	fmt.Println("  production  - Production environment")
}

func runMigrations(dm *database.DatabaseManager) error {
	log.Println("Running database migrations...")
	
	migrator := migrations.NewMigrator(dm.GetDB())
	
	// Add all migrations
	for _, migration := range migrations.GetAllMigrations() {
		migrator.AddMigration(migration)
	}
	
	return migrator.Up()
}

func rollbackMigration(dm *database.DatabaseManager) error {
	log.Println("Rolling back last migration...")
	
	migrator := migrations.NewMigrator(dm.GetDB())
	
	// Add all migrations
	for _, migration := range migrations.GetAllMigrations() {
		migrator.AddMigration(migration)
	}
	
	return migrator.Down()
}

func showMigrationStatus(dm *database.DatabaseManager) error {
	log.Println("Checking migration status...")
	
	migrator := migrations.NewMigrator(dm.GetDB())
	
	// Add all migrations
	for _, migration := range migrations.GetAllMigrations() {
		migrator.AddMigration(migration)
	}
	
	return migrator.Status()
}

func seedDatabase(dm *database.DatabaseManager, env string) error {
	log.Printf("Seeding database for environment: %s", env)
	
	seeder := database.NewSeeder(dm.GetDB())
	
	switch env {
	case "development":
		return seeder.SeedDevelopmentData()
	case "staging":
		// Staging might have different seed data
		return seeder.SeedDevelopmentData()
	case "production":
		log.Println("Warning: Seeding production database is not recommended")
		fmt.Print("Are you sure you want to seed production data? (yes/no): ")
		var response string
		fmt.Scanln(&response)
		if response != "yes" {
			log.Println("Seeding cancelled")
			return nil
		}
		return seeder.SeedDevelopmentData()
	default:
		return fmt.Errorf("unknown environment: %s", env)
	}
}

func checkHealth(dm *database.DatabaseManager) error {
	log.Println("Checking database health...")
	
	// Test connection
	if err := dm.HealthCheck(); err != nil {
		log.Printf("❌ Database health check failed: %v", err)
		return err
	}
	
	log.Println("✅ Database connection is healthy")
	
	// Get connection statistics
	stats, err := dm.GetConnectionStats()
	if err != nil {
		log.Printf("Warning: Could not get connection stats: %v", err)
	} else {
		log.Printf("📊 Connection Stats:")
		log.Printf("   Max Open Connections: %d", stats.MaxOpenConnections)
		log.Printf("   Open Connections: %d", stats.OpenConnections)
		log.Printf("   In Use: %d", stats.InUse)
		log.Printf("   Idle: %d", stats.Idle)
		log.Printf("   Wait Count: %d", stats.WaitCount)
		log.Printf("   Wait Duration: %v", stats.WaitDuration)
		log.Printf("   Max Idle Closed: %d", stats.MaxIdleClosed)
		log.Printf("   Max Idle Time Closed: %d", stats.MaxIdleTimeClosed)
		log.Printf("   Max Lifetime Closed: %d", stats.MaxLifetimeClosed)
	}
	
	return nil
}
