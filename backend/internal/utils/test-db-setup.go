package utils

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/transaction-tracker/backend/migrations"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// TestDBConfig holds DB connection info
type TestDBConfig struct {
	Host   string
	Port   string
	User   string
	Pass   string
	DBName string
}

// GetTestDBConfig loads DB config from env or defaults
func GetTestDBConfig() TestDBConfig {
	cfg := TestDBConfig{
		Host:   os.Getenv("TEST_DB_HOST"),
		Port:   os.Getenv("TEST_DB_PORT"),
		User:   os.Getenv("TEST_DB_USER"),
		Pass:   os.Getenv("TEST_DB_PASS"),
		DBName: os.Getenv("TEST_DB_NAME"),
	}
	if cfg.Host == "" {
		cfg.Host = "localhost"
	}
	if cfg.Port == "" {
		cfg.Port = "3306"
	}
	if cfg.User == "" {
		cfg.User = "root"
	}
	if cfg.Pass == "" {
		cfg.Pass = "root"
	}
	if cfg.DBName == "" {
		cfg.DBName = "transaction_tracker_test"
	}
	return cfg
}

// GetTestDB returns a *gorm.DB for the test DB
func GetTestDB(t *testing.T) *gorm.DB {
	cfg := GetTestDBConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	return db
}

// ResetTestDB drops all tables and runs all migrations for a clean test DB.
func ResetTestDB(db *gorm.DB, dbName string, t *testing.T) {
	sqlDB, err := db.DB()
	require.NoError(t, err)
	_, _ = sqlDB.Exec("SET FOREIGN_KEY_CHECKS = 0;")
	// Drop tables if they exist (including migration tracking table)
	tables := []string{"jwt_tokens", "transactions", "users", "schema_migrations"} // Order matters for foreign keys
	for _, table := range tables {
		_, _ = sqlDB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s;", table))
	}
	_, _ = sqlDB.Exec("SET FOREIGN_KEY_CHECKS = 1;")
	// Run all migrations
	require.NoError(t, migrations.ApplyAllMigrations(db))
}

// SetupTestDB creates a MySQL test database connection, cleans tables, and runs real migrations.
func SetupTestDB(t *testing.T) *gorm.DB {
	cfg := GetTestDBConfig()
	// First connect without specifying the database to create it if it doesn't exist
	rootDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local", cfg.User, cfg.Pass, cfg.Host, cfg.Port)
	rootDB, err := gorm.Open(mysql.Open(rootDSN), &gorm.Config{})
	require.NoError(t, err)

	// Create the test database if it doesn't exist
	err = rootDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", cfg.DBName)).Error
	require.NoError(t, err)

	// Now connect to the test database
	db := GetTestDB(t)

	// Reset and migrate
	ResetTestDB(db, cfg.DBName, t)

	return db
}
