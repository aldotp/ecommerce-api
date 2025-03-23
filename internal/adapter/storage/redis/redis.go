package redis

import (
	"context"
	"errors"
	"time"

	"github.com/aldotp/ecommerce-go-api/internal/adapter/config"
	"github.com/aldotp/ecommerce-go-api/internal/core/port"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

// New creates a new instance of Redis
func New(ctx context.Context, config *config.Redis) (port.CacheInterface, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       0,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &Redis{client}, nil
}

// Set stores the value in Redis
func (r *Redis) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

// Get retrieves the value from Redis
func (r *Redis) Get(ctx context.Context, key string) ([]byte, error) {
	res, err := r.client.Get(ctx, key).Result()
	return []byte(res), err
}

// Delete removes a key from Redis
func (r *Redis) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// DeleteByPrefix removes keys with a given prefix
func (r *Redis) DeleteByPrefix(ctx context.Context, prefix string) error {
	var cursor uint64
	var keys []string

	for {
		var err error
		keys, cursor, err = r.client.Scan(ctx, cursor, prefix, 100).Result()
		if err != nil {
			return err
		}

		for _, key := range keys {
			err := r.client.Del(ctx, key).Err()
			if err != nil {
				return err
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

// Close closes the Redis connection
func (r *Redis) Close() error {
	return r.client.Close()
}

// AcquireLock tries to acquire a lock using SETNX
func (r *Redis) AcquireLock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	acquired, err := r.client.SetNX(ctx, key, "locked", ttl).Result()
	if err != nil {
		return false, err
	}
	return acquired, nil
}

// ReleaseLock deletes the lock key
func (r *Redis) ReleaseLock(ctx context.Context, key string) error {
	deleted, err := r.client.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	if deleted == 0 {
		return errors.New("lock not found or already released")
	}
	return nil
}
