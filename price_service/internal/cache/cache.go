package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/transaction-tracker/price_service/internal/config"
	"github.com/transaction-tracker/price_service/internal/logger"
)

type Service struct {
	client     *redis.Client
	defaultTTL time.Duration
}

func NewService(cfg *config.Config) (*Service, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Service{
		client:     rdb,
		defaultTTL: cfg.Cache.DefaultTTL,
	}, nil
}

// Redis wrapper methods with logging
func (s *Service) redisGet(ctx context.Context, key string) (string, error) {
	logger.Info("Redis Cache Get", logger.H{
		"operation": "cache_get",
		"key":       key,
	})

	result, err := s.client.Get(ctx, key).Result()

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

func (s *Service) redisSet(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	logger.Info("Redis Cache Set", logger.H{
		"operation": "cache_set",
		"key":       key,
		"ttl":       expiration.String(),
	})

	err := s.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		logger.Warn("Redis Cache Set Failed", logger.H{
			"operation": "cache_set",
			"key":       key,
			"error":     err.Error(),
		})
	}
	return err
}

func (s *Service) redisDel(ctx context.Context, keys ...string) error {
	logger.Info("Redis Cache Delete", logger.H{
		"operation": "cache_delete",
		"keys":      keys,
		"count":     len(keys),
	})

	err := s.client.Del(ctx, keys...).Err()
	if err != nil {
		logger.Warn("Redis Cache Delete Failed", logger.H{
			"operation": "cache_delete",
			"keys":      keys,
			"error":     err.Error(),
		})
	}
	return err
}

func (s *Service) Close() error {
	return s.client.Close()
}
