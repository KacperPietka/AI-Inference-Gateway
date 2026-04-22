package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	client      *redis.Client
	maxRequests int
	window      time.Duration
}

// Creates a RateLimited connected to Redis
func New(redisURL string, maxRequests int, windowSecs int) (*RateLimiter, error) {
	client := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})

	// Verify connection immediately -- fail fast
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &RateLimiter{
		client:      client,
		maxRequests: maxRequests,
		window:      time.Duration(windowSecs) * time.Second,
	}, nil
}

// Allow checks if a user is within their limit
// Return true if the request is allowed
func (rl *RateLimiter) Allow(ctx context.Context, userID string) (bool, error) {
	key := fmt.Sprintf("ratelimit:%s", userID)

	count, err := rl.client.Incr(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("redis incr failed: %w", err)
	}

	// First request -> set the expiry
	if count == 1 {
		if err := rl.client.Expire(ctx, key, rl.window).Err(); err != nil {
			return false, fmt.Errorf("redis expire failed: %w", err)
		}
	}

	return count <= int64(rl.maxRequests), nil
}

// Returns how many requests the user has left in the current window
func (rl *RateLimiter) Remaining(ctx context.Context, userID string) (int, error) {
	key := fmt.Sprintf("ratelimit:%s", userID)

	count, err := rl.client.Get(ctx, key).Int()
	if err == redis.Nil {
		return rl.maxRequests, nil
	}

	if err != nil {
		return 0, fmt.Errorf("redis get failed: %w", err)
	}

	remaining := rl.maxRequests - count
	if remaining < 0 {
		return 0, nil
	}
	return remaining, nil
}

func (rl *RateLimiter) Close() error {
	return rl.client.Close()
}
