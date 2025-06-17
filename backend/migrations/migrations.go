package migrations

import (
	"strings"
	"time"

	"github.com/transaction-tracker/backend/internal/models"
	"gorm.io/gorm"
)

// GetAllMigrations returns all available migrations
func GetAllMigrations() []Migration {
	return []Migration{
		{
			ID:          "001_create_users_table",
			Description: "Create users table with authentication and profile information",
			CreatedAt:   time.Date(2025, 6, 16, 0, 0, 0, 0, time.UTC),
			Up: func(db *gorm.DB) error {
				// Create users table using GORM AutoMigrate
				return db.AutoMigrate(&models.User{})
			},
			Down: func(db *gorm.DB) error {
				// Drop users table
				return db.Migrator().DropTable(&models.User{})
			},
		},
		{
			ID:          "002_create_transactions_table",
			Description: "Create transactions table with comprehensive financial transaction tracking",
			CreatedAt:   time.Date(2025, 6, 16, 0, 1, 0, 0, time.UTC),
			Up: func(db *gorm.DB) error {
				// Create transactions table using GORM AutoMigrate
				return db.AutoMigrate(&models.Transaction{})
			},
			Down: func(db *gorm.DB) error {
				// Drop transactions table
				return db.Migrator().DropTable(&models.Transaction{})
			},
		},
		{
			ID:          "003_create_indexes",
			Description: "Create additional indexes for performance optimization",
			CreatedAt:   time.Date(2025, 6, 16, 0, 2, 0, 0, time.UTC),
			Up: func(db *gorm.DB) error {
				// Create composite indexes for common queries
				indexes := []string{
					"CREATE INDEX idx_transactions_user_date ON transactions(user_id, transaction_date)",
					"CREATE INDEX idx_transactions_symbol_date ON transactions(symbol, transaction_date)",
					"CREATE INDEX idx_transactions_type_status ON transactions(type, status)",
					"CREATE INDEX idx_users_email_active ON users(email, is_active)",
				}
				
				for _, index := range indexes {
					if err := db.Exec(index).Error; err != nil {
						// Ignore error if index already exists
						if !strings.Contains(err.Error(), "Duplicate key name") {
							return err
						}
					}
				}
				
				return nil
			},
			Down: func(db *gorm.DB) error {
				// Drop additional indexes
				indexes := []string{
					"DROP INDEX idx_transactions_user_date ON transactions",
					"DROP INDEX idx_transactions_symbol_date ON transactions",
					"DROP INDEX idx_transactions_type_status ON transactions",
					"DROP INDEX idx_users_email_active ON users",
				}
				
				for _, index := range indexes {
					if err := db.Exec(index).Error; err != nil {
						// Ignore error if index doesn't exist
						if !strings.Contains(err.Error(), "check that it exists") {
							return err
						}
					}
				}
				
				return nil
			},
		},
	}
}
