package migrations

import (
	"fmt"
	"sort"
	"time"

	"github.com/transaction-tracker/backend/internal/logger"

	"gorm.io/gorm"
)

// Migration represents a database migration step.
type Migration struct {
	ID          string               // Unique migration ID (e.g. 001_init_schema)
	Description string               // Human-readable description
	Up          func(*gorm.DB) error // Function to apply the migration
	Down        func(*gorm.DB) error // Function to rollback the migration (optional)
}

// MigrationRecord represents a record of an applied migration in the DB.
type MigrationRecord struct {
	ID        string `gorm:"primaryKey;size:255"`
	AppliedAt time.Time
}

func (MigrationRecord) TableName() string { return "schema_migrations" }

// Migrator manages and applies migrations.
type Migrator struct {
	db         *gorm.DB
	migrations []Migration
}

// NewMigrator creates a new migrator instance.
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{db: db, migrations: make([]Migration, 0)}
}

// NewMigratorWithDefaults creates a migrator and adds all default migrations.
func NewMigratorWithDefaults(db *gorm.DB) *Migrator {
	m := NewMigrator(db)
	for _, mig := range GetAllMigrations() {
		m.AddMigration(mig)
	}
	return m
}

// AddMigration adds a migration to the migrator (private).
func (m *Migrator) AddMigration(migration Migration) {
	m.migrations = append(m.migrations, migration)
}

// ensureMigrationTable creates the migration table if it doesn't exist.
func (m *Migrator) ensureMigrationTable() error {
	return m.db.AutoMigrate(&MigrationRecord{})
}

// getAppliedMigrations returns a set of applied migration IDs.
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

// ApplyAll runs all pending migrations.
func (m *Migrator) ApplyAll() error {
	if err := m.ensureMigrationTable(); err != nil {
		return fmt.Errorf("failed to ensure migration table: %w", err)
	}
	applied, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}
	sort.Slice(m.migrations, func(i, j int) bool { return m.migrations[i].ID < m.migrations[j].ID })
	var appliedCount int
	for _, migration := range m.migrations {
		if applied[migration.ID] {
			continue
		}
		logger.Info("Applying migration", logger.H{"id": migration.ID, "desc": migration.Description})
		tx := m.db.Begin()
		if err := migration.Up(tx); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to apply migration %s: %w", migration.ID, err)
		}
		record := MigrationRecord{ID: migration.ID, AppliedAt: time.Now().UTC()}
		if err := tx.Create(&record).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", migration.ID, err)
		}
		if err := tx.Commit().Error; err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", migration.ID, err)
		}
		appliedCount++
		logger.Info("Migration applied successfully", logger.H{"id": migration.ID})
	}
	if appliedCount == 0 {
		logger.Info("No pending migrations found", logger.H{})
	} else {
		logger.Info("Migrations applied", logger.H{"count": appliedCount})
	}
	return nil
}

// RollbackLast rolls back the last applied migration.
func (m *Migrator) RollbackLast() error {
	if err := m.ensureMigrationTable(); err != nil {
		return fmt.Errorf("failed to ensure migration table: %w", err)
	}
	var lastRecord MigrationRecord
	if err := m.db.Order("applied_at DESC").First(&lastRecord).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Info("No migrations to rollback", logger.H{})
			return nil
		}
		return fmt.Errorf("failed to get last migration: %w", err)
	}
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
	logger.Info("Rolling back migration", logger.H{"id": targetMigration.ID, "desc": targetMigration.Description})
	tx := m.db.Begin()
	if err := targetMigration.Down(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to rollback migration %s: %w", targetMigration.ID, err)
	}
	if err := tx.Delete(&lastRecord).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to remove migration record %s: %w", targetMigration.ID, err)
	}
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit rollback %s: %w", targetMigration.ID, err)
	}
	logger.Info("Migration rolled back successfully", logger.H{"id": targetMigration.ID})
	return nil
}

// Status prints the status of all migrations.
func (m *Migrator) Status() error {
	if err := m.ensureMigrationTable(); err != nil {
		return fmt.Errorf("failed to ensure migration table: %w", err)
	}
	applied, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}
	sort.Slice(m.migrations, func(i, j int) bool { return m.migrations[i].ID < m.migrations[j].ID })
	logger.Info("Migration Status", logger.H{})
	logger.Info("================", logger.H{})
	for _, migration := range m.migrations {
		status := "Pending"
		if applied[migration.ID] {
			status = "Applied"
		}
		logger.Info("Migration", logger.H{"status": status, "id": migration.ID, "desc": migration.Description})
	}
	return nil
}
