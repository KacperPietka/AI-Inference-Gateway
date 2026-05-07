package cache

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

const KeyPrefix = "generate:"

// TTL controls how long a cached response lives in Redis
type Config struct {
	TTLSeconds int
}

func GenerateKey(prompt, model string) string {

	// Normalize the prompt
	normalized := strings.ToLower(strings.TrimSpace(prompt))

	// Hash the normalized prompt
	hash := sha256.Sum256([]byte(normalized))

	return fmt.Sprintf("%s%s:%x", KeyPrefix, model, hash[:8])
}
