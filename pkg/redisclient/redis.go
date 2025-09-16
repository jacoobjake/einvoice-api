// pkg/redisclient/redis.go
package redisclient

import (
	"context"
	"time"

	cfg_redis "github.com/jacoobjake/einvoice-api/config/redis"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	rdb *redis.Client
}

// New creates a new Redis client with given options.
func NewRedisClient(redisCfg *cfg_redis.RedisConfig) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})

	return &RedisClient{rdb: rdb}
}

// Set stores a value with optional expiration.
func (c *RedisClient) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return c.rdb.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value by key.
func (c *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return c.rdb.Get(ctx, key).Result()
}

// Exists checks if a key exists.
func (c *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.rdb.Exists(ctx, key).Result()
	return n > 0, err
}

// Delete removes a key.
func (c *RedisClient) Delete(ctx context.Context, key string) error {
	return c.rdb.Del(ctx, key).Err()
}

// Expire sets a TTL on a key.
func (c *RedisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.rdb.Expire(ctx, key, expiration).Err()
}
