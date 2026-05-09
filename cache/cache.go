package cache

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"
	"time"
)

const KeyPrefix = "generate:"

// TTL controls how long a cached response lives in Redis
type Config struct {
	TTLSeconds int
}

var ErrCacheMiss = errors.New("cache miss")

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Close() error
}

func GenerateKey(prompt, model string) string {

	// Normalize the prompt
	normalized := strings.ToLower(strings.TrimSpace(prompt))

	// Hash the normalized prompt
	hash := sha256.Sum256([]byte(normalized))

	return fmt.Sprintf("%s%s:%x", KeyPrefix, model, hash[:8])
}
