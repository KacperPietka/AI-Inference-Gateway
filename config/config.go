package config

import (
	"os"
	"strconv"
)

type Config struct {
	OllamaURL           string
	DefaultModel        string
	SecondaryModel      string
	ServerPort          string
	RedisURL            string
	RateLimitRequests   int
	RateLimitWindowSecs int
	CacheTTLSeconds     int
	GeminiAPIKey        string
	GeminiModel         string
}

func Load() *Config {
	return &Config{
		OllamaURL:           getEnv("OLLAMA_URL", "http://localhost:11434/api/generate"),
		DefaultModel:        getEnv("DEFAULT_MODEL", "tinyllama"),
		SecondaryModel:      getEnv("SECONDARY_MODEL", "mistral"),
		ServerPort:          getEnv("SERVER_PORT", ":8080"),
		RedisURL:            getEnv("REDIS_URL", "localhost:6379"),
		RateLimitRequests:   getEnvInt("RATE_LIMIT_REQUESTS", 10),
		RateLimitWindowSecs: getEnvInt("RATE_LIMIT_WINDOW_SECONDS", 60),
		CacheTTLSeconds:     getEnvInt("CACHE_TTL_SECONDS", 3600),
		GeminiAPIKey:        getEnv("GEMINI_API_KEY", ""),
		GeminiModel:         getEnv("GEMINI_MODEL", "gemini-3-flash-preview"),
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
