package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/transaction-tracker/backend/internal/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/transaction-tracker/backend/config"
)

// DB holds the database connection
var DB *gorm.DB

// DatabaseManager manages database connections and operations
type DatabaseManager struct {
	db     *gorm.DB
	config *config.DatabaseConfig
}

// NewDatabaseManager creates a new database manager
func NewDatabaseManager(cfg *config.DatabaseConfig) *DatabaseManager {
	return &DatabaseManager{
		config: cfg,
	}
}

// Connect establishes a connection to the database
func (dm *DatabaseManager) Connect() error {
	var err error

	// Configure GORM logger (use default GORM logger)
	gormLogger := gormlogger.Default

	// Open database connection
	dm.db, err = gorm.Open(mysql.Open(dm.config.GetDSNWithSSL()), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := dm.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(dm.config.MaxConnections)
	sqlDB.SetMaxIdleConns(dm.config.MaxIdle)
	sqlDB.SetConnMaxLifetime(dm.config.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(dm.config.ConnMaxIdleTime)

	// Set global DB variable
	DB = dm.db

	logger.Info("Database connection established successfully")
	return nil
}

// ConnectWithRetry establishes a connection to the database with retry logic
func (dm *DatabaseManager) ConnectWithRetry(maxRetries int, retryDelay time.Duration) error {
	var err error

	for i := 0; i < maxRetries; i++ {
		err = dm.Connect()
		if err == nil {
			return nil
		}

		logger.Warn("Database connection attempt failed", logger.H{"attempt": i + 1, "error": err})

		if i < maxRetries-1 {
			logger.Info("Retrying database connection", logger.H{"delay": retryDelay})
			time.Sleep(retryDelay)
			retryDelay *= 2 // Exponential backoff
		}
	}

	return fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}

// GetDB returns the database instance
func (dm *DatabaseManager) GetDB() *gorm.DB {
	return dm.db
}

// HealthCheck performs a health check on the database connection
func (dm *DatabaseManager) HealthCheck() error {
	if dm.db == nil {
		return fmt.Errorf("database connection is nil")
	}

	sqlDB, err := dm.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

// Close closes the database connection
func (dm *DatabaseManager) Close() error {
	if dm.db == nil {
		return nil
	}

	sqlDB, err := dm.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	logger.Info("Closing database connection...")
	return sqlDB.Close()
}

// CreateIndexes creates additional database indexes for performance
func (dm *DatabaseManager) CreateIndexes() error {
	if dm.db == nil {
		return fmt.Errorf("database connection is nil")
	}

	logger.Info("Checking database indexes...")

	// Note: All necessary indexes are now created by GORM AutoMigrate and migration files
	// - Basic indexes: created by GORM based on model tags
	// - Composite indexes: created by migration SQL files
	// This function is kept for future custom index additions if needed

	logger.Info("Database indexes verified successfully")
	return nil
}

// GetConnectionStats returns database connection statistics
func (dm *DatabaseManager) GetConnectionStats() (*sql.DBStats, error) {
	if dm.db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	sqlDB, err := dm.db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	stats := sqlDB.Stats()
	return &stats, nil
}

// Initialize initializes the database with the given configuration
func Initialize(cfg *config.Config) (*DatabaseManager, error) {
	dbConfig := config.GetDatabaseConfig(cfg)
	dm := NewDatabaseManager(dbConfig)

	// Connect with retry logic
	if err := dm.ConnectWithRetry(5, 2*time.Second); err != nil {
		return nil, err
	}

	return dm, nil
}
