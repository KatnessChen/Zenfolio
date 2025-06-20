package migrations

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"gorm.io/gorm"
)

// GetAllMigrations returns all available migrations
func GetAllMigrations() []Migration {
	return []Migration{
		{
			ID:          "001_create_users_table",
			Description: "Create users table with authentication and profile information (simplified schema)",
			CreatedAt:   time.Date(2025, 6, 16, 0, 0, 0, 0, time.UTC),
			Up: func(db *gorm.DB) error {
				return executeSQLFile(db, "001_create_users_table.sql")
			},
			Down: func(db *gorm.DB) error {
				// Drop users table and its indexes
				return db.Exec("DROP TABLE IF EXISTS users").Error
			},
		},
		{
			ID:          "002_create_transactions_table",
			Description: "Create transactions table with simplified financial transaction tracking",
			CreatedAt:   time.Date(2025, 6, 16, 0, 1, 0, 0, time.UTC),
			Up: func(db *gorm.DB) error {
				return executeSQLFile(db, "002_create_transactions_table.sql")
			},
			Down: func(db *gorm.DB) error {
				// Drop transactions table and its indexes
				return db.Exec("DROP TABLE IF EXISTS transactions").Error
			},
		},
		{
			ID:          "003_create_jwt_tokens_table",
			Description: "Create JWT tokens table for token lifecycle management",
			CreatedAt:   time.Date(2025, 6, 19, 0, 0, 0, 0, time.UTC),
			Up: func(db *gorm.DB) error {
				return executeSQLFile(db, "003_create_jwt_tokens_table.sql")
			},
			Down: func(db *gorm.DB) error {
				// Drop jwt_tokens table and its indexes
				return db.Exec("DROP TABLE IF EXISTS jwt_tokens").Error
			},
		},
	}
}

// executeSQLFile reads and executes a SQL migration file
func executeSQLFile(db *gorm.DB, filename string) error {
	// Get the directory where this Go file is located
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("failed to get current file path")
	}

	migrationsDir := filepath.Dir(currentFile)
	sqlFilePath := filepath.Join(migrationsDir, filename)

	// Read the SQL file
	sqlBytes, err := os.ReadFile(sqlFilePath)
	if err != nil {
		return fmt.Errorf("failed to read SQL file %s: %w", filename, err)
	}

	sqlContent := string(sqlBytes)

	// Remove comments and empty lines
	lines := strings.Split(sqlContent, "\n")
	var cleanLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines and comment lines
		if line == "" || strings.HasPrefix(line, "--") || strings.HasPrefix(line, "//") {
			continue
		}
		cleanLines = append(cleanLines, line)
	}

	if len(cleanLines) == 0 {
		return fmt.Errorf("no SQL statements found in file %s", filename)
	}

	// Join lines and split by semicolon to get individual statements
	fullSQL := strings.Join(cleanLines, " ")
	statements := strings.Split(fullSQL, ";")

	// Execute each statement
	for _, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		}

		if err := db.Exec(statement).Error; err != nil {
			// For CREATE TABLE IF NOT EXISTS and CREATE INDEX, ignore "already exists" errors
			if strings.Contains(strings.ToLower(statement), "if not exists") &&
				(strings.Contains(err.Error(), "already exists") ||
					strings.Contains(err.Error(), "Duplicate key name")) {
				continue
			}
			return fmt.Errorf("failed to execute SQL statement: %s\nError: %w", statement, err)
		}
	}

	return nil
}
