package cache

import (
	"context"
	"time"
)

type RedisCache struct{}

// Compile time check
var _ Cache = (*RedisCache)(nil)

func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return "", ErrCacheMiss
}

func (c *RedisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return nil
}

func (c *RedisCache) Close() error {
	return nil
}
