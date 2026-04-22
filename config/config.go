package config

import (
	"os"
	"strconv"
)

type Config struct {
	OllamaURL           string
	DefaultModel        string
	ServerPort          string
	RedisURL            string
	RateLimitRequests   int
	RateLimitWindowSecs int
}

func Load() *Config {
	return &Config{
		OllamaURL:           getEnv("OLLAMA_URL", "http://localhost:11434/api/generate"),
		DefaultModel:        getEnv("DEFAULT_MODEL", "tinyllama"),
		ServerPort:          getEnv("SERVER_PORT", ":8080"),
		RedisURL:            getEnv("REDIS_URL", "localhost:6379"),
		RateLimitRequests:   getEnvInt("RATE_LIMIT_REQUESTS", 10),
		RateLimitWindowSecs: getEnvInt("RATE_LIMIT_WINDOW_SECONDS", 60),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}
