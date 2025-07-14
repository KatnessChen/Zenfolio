package migrations

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gorm.io/gorm"
)

// GetAllMigrations returns all available migrations
func GetAllMigrations() []Migration {
	return []Migration{
		{
			ID:          "000_init_schema",
			Description: "Unified schema: users, transactions, jwt_tokens with UUID primary keys",
			Up: func(db *gorm.DB) error {
				return executeSQLFile(db, "000_init_schema.sql")
			},
			Down: func(db *gorm.DB) error {
				return db.Exec("DROP TABLE IF EXISTS jwt_tokens; DROP TABLE IF EXISTS transactions; DROP TABLE IF EXISTS users;").Error
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

// ApplyAllMigrations creates a migrator with all default migrations and applies them.
func ApplyAllMigrations(db *gorm.DB) error {
	migrator := NewMigratorWithDefaults(db)
	return migrator.ApplyAll()
}
