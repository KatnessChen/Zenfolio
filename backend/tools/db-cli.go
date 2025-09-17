package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/transaction-tracker/backend/internal/logger"

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

	// Initialize structured logger
	logger.InitLogger()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load configuration", err, logger.H{})
		os.Exit(1)
	}

	// Initialize database
	dm, err := database.Initialize(cfg)
	if err != nil {
		logger.Error("Failed to initialize database", err, logger.H{})
		os.Exit(1)
	}
	defer dm.Close()

	switch *action {
	case "migrate":
		if err := runMigrations(dm); err != nil {
			logger.Error("Migration failed", err, logger.H{})
			os.Exit(1)
		}
	case "rollback":
		if err := rollbackMigration(dm); err != nil {
			logger.Error("Rollback failed", err, logger.H{})
			os.Exit(1)
		}
	case "status":
		if err := showMigrationStatus(dm); err != nil {
			logger.Error("Status check failed", err, logger.H{})
			os.Exit(1)
		}
	case "seed":
		if err := seedDatabase(dm, *env); err != nil {
			logger.Error("Seeding failed", err, logger.H{})
			os.Exit(1)
		}
	case "health":
		if err := checkHealth(dm); err != nil {
			logger.Error("Health check failed", err, logger.H{})
			os.Exit(1)
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
	logger.Info("Running database migrations...", logger.H{})
	// Use the helper function that applies all default migrations
	return migrations.ApplyAllMigrations(dm.GetDB())
}

func rollbackMigration(dm *database.DatabaseManager) error {
	logger.Info("Rolling back last migration...", logger.H{})
	migrator := migrations.NewMigratorWithDefaults(dm.GetDB())
	return migrator.RollbackLast()
}

func showMigrationStatus(dm *database.DatabaseManager) error {
	logger.Info("Checking migration status...", logger.H{})
	migrator := migrations.NewMigratorWithDefaults(dm.GetDB())
	return migrator.Status()
}

func seedDatabase(dm *database.DatabaseManager, env string) error {
	logger.Info("Seeding database", logger.H{"environment": env})
	seeder := database.NewSeeder(dm.GetDB())
	switch env {
	case "development":
		return seeder.SeedDevelopmentData()
	case "staging":
		// Staging might have different seed data
		return seeder.SeedDevelopmentData()
	case "production":
		logger.Warn("Seeding production database is not recommended", logger.H{})
		fmt.Print("Are you sure you want to seed production data? (yes/no): ")
		var response string
		if _, err := fmt.Scanln(&response); err != nil {
			logger.Error("Error reading input", err, logger.H{})
			return err
		}
		if response != "yes" {
			logger.Info("Seeding cancelled", logger.H{})
			return nil
		}
		return seeder.SeedDevelopmentData()
	default:
		return fmt.Errorf("unknown environment: %s", env)
	}
}

func checkHealth(dm *database.DatabaseManager) error {
	logger.Info("Checking database health...", logger.H{})
	// Test connection
	if err := dm.HealthCheck(); err != nil {
		logger.Error("Database health check failed", err, logger.H{})
		return err
	}
	logger.Info("Database connection is healthy", logger.H{})
	// Get connection statistics
	stats, err := dm.GetConnectionStats()
	if err != nil {
		logger.Warn("Could not get connection stats", logger.H{"error": err.Error()})
	} else {
		logger.Info("Connection Stats", logger.H{
			"MaxOpenConnections": stats.MaxOpenConnections,
			"OpenConnections":    stats.OpenConnections,
			"InUse":              stats.InUse,
			"Idle":               stats.Idle,
			"WaitCount":          stats.WaitCount,
			"WaitDuration":       stats.WaitDuration,
			"MaxIdleClosed":      stats.MaxIdleClosed,
			"MaxIdleTimeClosed":  stats.MaxIdleTimeClosed,
			"MaxLifetimeClosed":  stats.MaxLifetimeClosed,
		})
	}
	return nil
}
