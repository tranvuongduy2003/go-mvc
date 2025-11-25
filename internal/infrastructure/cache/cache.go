package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/logger"
)

// Service provides caching functionality using Redis
type Service struct {
	client *redis.Client
	logger *logger.Logger
}

// CacheOptions contains configuration options for cache operations
type CacheOptions struct {
	TTL time.Duration
}

// NewCacheService creates a new cache service
func NewCacheService(client *redis.Client, logger *logger.Logger) *Service {
	return &Service{
		client: client,
		logger: logger,
	}
}

// Set stores a value in cache with optional TTL
func (s *Service) Set(ctx context.Context, key string, value interface{}, options *CacheOptions) error {
	s.logger.Debugf("Setting cache key: %s", key)

	// Serialize value to JSON
	data, err := json.Marshal(value)
	if err != nil {
		s.logger.Errorf("Failed to marshal value for cache key %s: %v", key, err)
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	// Determine TTL
	var ttl time.Duration
	if options != nil && options.TTL > 0 {
		ttl = options.TTL
	} else {
		ttl = 0 // No expiration
	}

	// Set value in Redis
	if err := s.client.Set(ctx, key, data, ttl).Err(); err != nil {
		s.logger.Errorf("Failed to set cache key %s: %v", key, err)
		return fmt.Errorf("failed to set cache: %w", err)
	}

	s.logger.Debugf("Successfully set cache key: %s", key)
	return nil
}

// Get retrieves a value from cache
func (s *Service) Get(ctx context.Context, key string, dest interface{}) error {
	s.logger.Debugf("Getting cache key: %s", key)

	// Get value from Redis
	data, err := s.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			s.logger.Debugf("Cache miss for key: %s", key)
			return ErrCacheMiss
		}
		s.logger.Errorf("Failed to get cache key %s: %v", key, err)
		return fmt.Errorf("failed to get cache: %w", err)
	}

	// Deserialize JSON to destination
	if err := json.Unmarshal([]byte(data), dest); err != nil {
		s.logger.Errorf("Failed to unmarshal cache value for key %s: %v", key, err)
		return fmt.Errorf("failed to unmarshal cache value: %w", err)
	}

	s.logger.Debugf("Successfully got cache key: %s", key)
	return nil
}

// Delete removes a value from cache
func (s *Service) Delete(ctx context.Context, key string) error {
	s.logger.Debugf("Deleting cache key: %s", key)

	if err := s.client.Del(ctx, key).Err(); err != nil {
		s.logger.Errorf("Failed to delete cache key %s: %v", key, err)
		return fmt.Errorf("failed to delete cache: %w", err)
	}

	s.logger.Debugf("Successfully deleted cache key: %s", key)
	return nil
}

// DeletePattern removes all keys matching a pattern
func (s *Service) DeletePattern(ctx context.Context, pattern string) error {
	s.logger.Debugf("Deleting cache keys with pattern: %s", pattern)

	// Get all keys matching the pattern
	keys, err := s.client.Keys(ctx, pattern).Result()
	if err != nil {
		s.logger.Errorf("Failed to get keys for pattern %s: %v", pattern, err)
		return fmt.Errorf("failed to get keys: %w", err)
	}

	if len(keys) == 0 {
		s.logger.Debugf("No keys found for pattern: %s", pattern)
		return nil
	}

	// Delete all matching keys
	if err := s.client.Del(ctx, keys...).Err(); err != nil {
		s.logger.Errorf("Failed to delete keys for pattern %s: %v", pattern, err)
		return fmt.Errorf("failed to delete keys: %w", err)
	}

	s.logger.Debugf("Successfully deleted %d keys for pattern: %s", len(keys), pattern)
	return nil
}

// Exists checks if a key exists in cache
func (s *Service) Exists(ctx context.Context, key string) (bool, error) {
	s.logger.Debugf("Checking existence of cache key: %s", key)

	count, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		s.logger.Errorf("Failed to check existence of cache key %s: %v", key, err)
		return false, fmt.Errorf("failed to check cache existence: %w", err)
	}

	exists := count > 0
	s.logger.Debugf("Cache key %s exists: %v", key, exists)
	return exists, nil
}

// TTL returns the remaining time to live of a key
func (s *Service) TTL(ctx context.Context, key string) (time.Duration, error) {
	s.logger.Debugf("Getting TTL for cache key: %s", key)

	ttl, err := s.client.TTL(ctx, key).Result()
	if err != nil {
		s.logger.Errorf("Failed to get TTL for cache key %s: %v", key, err)
		return 0, fmt.Errorf("failed to get TTL: %w", err)
	}

	s.logger.Debugf("TTL for cache key %s: %v", key, ttl)
	return ttl, nil
}

// Expire sets a timeout on a key
func (s *Service) Expire(ctx context.Context, key string, expiration time.Duration) error {
	s.logger.Debugf("Setting expiration for cache key: %s, expiration: %v", key, expiration)

	if err := s.client.Expire(ctx, key, expiration).Err(); err != nil {
		s.logger.Errorf("Failed to set expiration for cache key %s: %v", key, err)
		return fmt.Errorf("failed to set expiration: %w", err)
	}

	s.logger.Debugf("Successfully set expiration for cache key: %s", key)
	return nil
}

// Increment increments a numeric value in cache
func (s *Service) Increment(ctx context.Context, key string) (int64, error) {
	s.logger.Debugf("Incrementing cache key: %s", key)

	val, err := s.client.Incr(ctx, key).Result()
	if err != nil {
		s.logger.Errorf("Failed to increment cache key %s: %v", key, err)
		return 0, fmt.Errorf("failed to increment: %w", err)
	}

	s.logger.Debugf("Successfully incremented cache key %s to %d", key, val)
	return val, nil
}

// IncrementBy increments a numeric value by a specific amount
func (s *Service) IncrementBy(ctx context.Context, key string, value int64) (int64, error) {
	s.logger.Debugf("Incrementing cache key %s by %d", key, value)

	val, err := s.client.IncrBy(ctx, key, value).Result()
	if err != nil {
		s.logger.Errorf("Failed to increment cache key %s by %d: %v", key, value, err)
		return 0, fmt.Errorf("failed to increment by: %w", err)
	}

	s.logger.Debugf("Successfully incremented cache key %s by %d to %d", key, value, val)
	return val, nil
}

// SetNX sets a key only if it doesn't exist (atomic operation)
func (s *Service) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	s.logger.Debugf("Setting cache key with NX: %s", key)

	// Serialize value to JSON
	data, err := json.Marshal(value)
	if err != nil {
		s.logger.Errorf("Failed to marshal value for cache key %s: %v", key, err)
		return false, fmt.Errorf("failed to marshal value: %w", err)
	}

	// Set value only if key doesn't exist
	success, err := s.client.SetNX(ctx, key, data, expiration).Result()
	if err != nil {
		s.logger.Errorf("Failed to set NX cache key %s: %v", key, err)
		return false, fmt.Errorf("failed to set NX cache: %w", err)
	}

	s.logger.Debugf("SetNX for cache key %s: %v", key, success)
	return success, nil
}

// GetOrSet retrieves a value from cache, or sets it if not found
func (s *Service) GetOrSet(ctx context.Context, key string, dest interface{}, setter func() (interface{}, error), options *CacheOptions) error {
	s.logger.Debugf("GetOrSet for cache key: %s", key)

	// Try to get from cache first
	err := s.Get(ctx, key, dest)
	if err == nil {
		s.logger.Debugf("Cache hit for key: %s", key)
		return nil
	}

	if err != ErrCacheMiss {
		s.logger.Errorf("Unexpected error getting cache key %s: %v", key, err)
		return err
	}

	s.logger.Debugf("Cache miss for key %s, calling setter", key)

	// Cache miss, call setter to get value
	value, err := setter()
	if err != nil {
		s.logger.Errorf("Setter failed for cache key %s: %v", key, err)
		return fmt.Errorf("setter failed: %w", err)
	}

	// Set value in cache
	if err := s.Set(ctx, key, value, options); err != nil {
		s.logger.Errorf("Failed to set cache after setter for key %s: %v", key, err)
		// Don't return error here, just log it
	}

	// Unmarshal the value to destination
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal setter result: %w", err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("failed to unmarshal setter result: %w", err)
	}

	s.logger.Debugf("Successfully executed GetOrSet for key: %s", key)
	return nil
}

// Ping checks the connection to Redis
func (s *Service) Ping(ctx context.Context) error {
	s.logger.Debug("Pinging Redis")

	if err := s.client.Ping(ctx).Err(); err != nil {
		s.logger.Errorf("Failed to ping Redis: %v", err)
		return fmt.Errorf("failed to ping Redis: %w", err)
	}

	s.logger.Debug("Successfully pinged Redis")
	return nil
}

// Close closes the Redis connection
func (s *Service) Close() error {
	s.logger.Info("Closing Redis connection")

	if err := s.client.Close(); err != nil {
		s.logger.Errorf("Failed to close Redis connection: %v", err)
		return fmt.Errorf("failed to close Redis connection: %w", err)
	}

	s.logger.Info("Successfully closed Redis connection")
	return nil
}
