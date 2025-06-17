package migrations

import (
	"fmt"
	"log"
	"sort"
	"time"

	"gorm.io/gorm"
)

// Migration represents a database migration
type Migration struct {
	ID          string
	Description string
	Up          func(*gorm.DB) error
	Down        func(*gorm.DB) error
	CreatedAt   time.Time
}

// MigrationRecord represents a migration record in the database
type MigrationRecord struct {
	ID        string `gorm:"primaryKey;size:255"`
	AppliedAt time.Time
}

// TableName specifies the table name for migration records
func (MigrationRecord) TableName() string {
	return "schema_migrations"
}

// Migrator handles database migrations
type Migrator struct {
	db         *gorm.DB
	migrations []Migration
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{
		db:         db,
		migrations: make([]Migration, 0),
	}
}

// AddMigration adds a migration to the migrator
func (m *Migrator) AddMigration(migration Migration) {
	m.migrations = append(m.migrations, migration)
}

// ensureMigrationTable creates the migration table if it doesn't exist
func (m *Migrator) ensureMigrationTable() error {
	return m.db.AutoMigrate(&MigrationRecord{})
}

// getAppliedMigrations returns a list of applied migration IDs
func (m *Migrator) getAppliedMigrations() (map[string]bool, error) {
	var records []MigrationRecord
	if err := m.db.Find(&records).Error; err != nil {
		return nil, err
	}

	applied := make(map[string]bool)
	for _, record := range records {
		applied[record.ID] = true
	}

	return applied, nil
}

// Up runs all pending migrations
func (m *Migrator) Up() error {
	if err := m.ensureMigrationTable(); err != nil {
		return fmt.Errorf("failed to ensure migration table: %w", err)
	}

	applied, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Sort migrations by ID
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].ID < m.migrations[j].ID
	})

	var appliedCount int
	for _, migration := range m.migrations {
		if applied[migration.ID] {
			continue
		}

		log.Printf("Applying migration: %s - %s", migration.ID, migration.Description)

		// Start transaction
		tx := m.db.Begin()

		// Run migration
		if err := migration.Up(tx); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to apply migration %s: %w", migration.ID, err)
		}

		// Record migration
		record := MigrationRecord{
			ID:        migration.ID,
			AppliedAt: time.Now().UTC(),
		}

		if err := tx.Create(&record).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", migration.ID, err)
		}

		// Commit transaction
		if err := tx.Commit().Error; err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", migration.ID, err)
		}

		appliedCount++
		log.Printf("Migration %s applied successfully", migration.ID)
	}

	if appliedCount == 0 {
		log.Println("No pending migrations found")
	} else {
		log.Printf("Applied %d migrations successfully", appliedCount)
	}

	return nil
}

// Down rolls back the last migration
func (m *Migrator) Down() error {
	if err := m.ensureMigrationTable(); err != nil {
		return fmt.Errorf("failed to ensure migration table: %w", err)
	}

	// Get the last applied migration
	var lastRecord MigrationRecord
	if err := m.db.Order("applied_at DESC").First(&lastRecord).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("No migrations to rollback")
			return nil
		}
		return fmt.Errorf("failed to get last migration: %w", err)
	}

	// Find the migration
	var targetMigration *Migration
	for _, migration := range m.migrations {
		if migration.ID == lastRecord.ID {
			targetMigration = &migration
			break
		}
	}

	if targetMigration == nil {
		return fmt.Errorf("migration %s not found in migration list", lastRecord.ID)
	}

	if targetMigration.Down == nil {
		return fmt.Errorf("migration %s does not have a down method", lastRecord.ID)
	}

	log.Printf("Rolling back migration: %s - %s", targetMigration.ID, targetMigration.Description)

	// Start transaction
	tx := m.db.Begin()

	// Run rollback
	if err := targetMigration.Down(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to rollback migration %s: %w", targetMigration.ID, err)
	}

	// Remove migration record
	if err := tx.Delete(&lastRecord).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to remove migration record %s: %w", targetMigration.ID, err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit rollback %s: %w", targetMigration.ID, err)
	}

	log.Printf("Migration %s rolled back successfully", targetMigration.ID)
	return nil
}

// Status shows the status of all migrations
func (m *Migrator) Status() error {
	if err := m.ensureMigrationTable(); err != nil {
		return fmt.Errorf("failed to ensure migration table: %w", err)
	}

	applied, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Sort migrations by ID
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].ID < m.migrations[j].ID
	})

	log.Println("Migration Status:")
	log.Println("================")

	for _, migration := range m.migrations {
		status := "Pending"
		if applied[migration.ID] {
			status = "Applied"
		}
		log.Printf("[%s] %s - %s", status, migration.ID, migration.Description)
	}

	return nil
}
