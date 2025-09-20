package config

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/transaction-tracker/backend/internal/logger"
)

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// RedisClient wraps redis client with additional functionality
type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisClient creates a new Redis client from environment variables
func NewRedisClient() (*RedisClient, error) {
	config := &RedisConfig{
		Host:     getEnvOrDefault("REDIS_HOST", "zenfolio_backend_redis"),
		Port:     getEnvOrDefaultInt("REDIS_PORT", 6379),
		Password: getEnvOrDefault("REDIS_PASSWORD", ""),
		DB:       getEnvOrDefaultInt("REDIS_DB", 0),
	}

	logger.Info("Redis Client Initialization", logger.H{
		"operation": "redis_init",
		"host":      config.Host,
		"port":      config.Port,
		"db":        config.DB,
	})

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	// Test connection
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		logger.Error("Redis Connection Failed", err, logger.H{
			"operation": "redis_connection_failed",
			"host":      config.Host,
			"port":      config.Port,
		})
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("Redis Connection Successful", logger.H{
		"operation": "redis_connection_success",
		"host":      config.Host,
		"port":      config.Port,
		"db":        config.DB,
	})

	return &RedisClient{
		client: client,
		ctx:    ctx,
	}, nil
}

// Set stores a key-value pair with optional expiration
func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	logger.Info("Redis Cache Set", logger.H{
		"operation": "cache_set",
		"key":       key,
		"ttl":       expiration.String(),
	})

	err := r.client.Set(r.ctx, key, value, expiration).Err()
	if err != nil {
		logger.Warn("Redis Cache Set Failed", logger.H{
			"operation": "cache_set",
			"key":       key,
			"error":     err.Error(),
		})
	}
	return err
}

// Get retrieves a value by key
func (r *RedisClient) Get(key string) (string, error) {
	logger.Info("Redis Cache Get", logger.H{
		"operation": "cache_get",
		"key":       key,
	})

	result, err := r.client.Get(r.ctx, key).Result()

	if err != nil {
		if err == redis.Nil {
			logger.Info("Redis Cache Miss", logger.H{
				"operation": "cache_miss",
				"key":       key,
			})
		} else {
			logger.Warn("Redis Cache Get Failed", logger.H{
				"operation": "cache_get",
				"key":       key,
				"error":     err.Error(),
			})
		}
		return result, err
	}

	logger.Info("Redis Cache Hit", logger.H{
		"operation": "cache_hit",
		"key":       key,
	})

	return result, nil
}

// Del deletes keys
func (r *RedisClient) Del(keys ...string) error {
	logger.Info("Redis Cache Delete", logger.H{
		"operation": "cache_delete",
		"keys":      keys,
		"count":     len(keys),
	})

	err := r.client.Del(r.ctx, keys...).Err()
	if err != nil {
		logger.Warn("Redis Cache Delete Failed", logger.H{
			"operation": "cache_delete",
			"keys":      keys,
			"error":     err.Error(),
		})
	}
	return err
}

// Exists checks if key exists
func (r *RedisClient) Exists(key string) (bool, error) {
	logger.Info("Redis Cache Exists Check", logger.H{
		"operation": "cache_exists",
		"key":       key,
	})

	result, err := r.client.Exists(r.ctx, key).Result()
	exists := result > 0

	if err != nil {
		logger.Warn("Redis Cache Exists Check Failed", logger.H{
			"operation": "cache_exists",
			"key":       key,
			"error":     err.Error(),
		})
	} else {
		logger.Info("Redis Cache Exists Result", logger.H{
			"operation": "cache_exists",
			"key":       key,
			"exists":    exists,
		})
	}

	return exists, err
}

// Close closes the Redis connection
func (r *RedisClient) Close() error {
	logger.Info("Redis Connection Closing", logger.H{
		"operation": "redis_close",
	})

	err := r.client.Close()
	if err != nil {
		logger.Error("Redis Connection Close Failed", err, logger.H{
			"operation": "redis_close_failed",
		})
	} else {
		logger.Info("Redis Connection Closed Successfully", logger.H{
			"operation": "redis_close_success",
		})
	}

	return err
}
