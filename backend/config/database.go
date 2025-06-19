package config

import (
	"fmt"
	"time"
)

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host            string
	Port            string
	Name            string
	User            string
	Password        string
	SSLMode         string
	MaxConnections  int
	MaxIdle         int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	Charset         string
	ParseTime       bool
	Loc             string
}

// GetDatabaseConfig returns database configuration based on environment
func GetDatabaseConfig(cfg *Config) *DatabaseConfig {
	return &DatabaseConfig{
		Host:            getEnvOrDefault("DB_HOST", "localhost"),
		Port:            getEnvOrDefault("DB_PORT", "3306"),
		Name:            getEnvOrDefault("DB_NAME", "transaction_tracker_dev"),
		User:            getEnvOrDefault("DB_USER", "root"),
		Password:        getEnvOrDefault("DB_PASSWORD", "root"),
		SSLMode:         getEnvOrDefault("DB_SSL_MODE", "disable"),
		MaxConnections:  getEnvOrDefaultInt("DB_MAX_CONNECTIONS", 100),
		MaxIdle:         getEnvOrDefaultInt("DB_MAX_IDLE", 10),
		ConnMaxLifetime: time.Duration(getEnvOrDefaultInt("DB_CONN_MAX_LIFETIME", 3600)) * time.Second,
		ConnMaxIdleTime: time.Duration(getEnvOrDefaultInt("DB_CONN_MAX_IDLE_TIME", 1800)) * time.Second,
		Charset:         getEnvOrDefault("DB_CHARSET", "utf8mb4"),
		ParseTime:       true,
		Loc:             getEnvOrDefault("DB_LOC", "Local"),
	}
}

// GetDSN returns the database connection string
func (db *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%t&loc=%s",
		db.User,
		db.Password,
		db.Host,
		db.Port,
		db.Name,
		db.Charset,
		db.ParseTime,
		db.Loc,
	)
}

// GetDSNWithSSL returns the database connection string with SSL configuration
func (db *DatabaseConfig) GetDSNWithSSL() string {
	dsn := db.GetDSN()
	if db.SSLMode != "disable" {
		dsn += "&tls=" + db.SSLMode
	}
	return dsn
}
