package config

import "os"

type Config struct {
	OllamaURL    string
	DefaultModel string
	ServerPort   string
}

func Load() *Config {
	return &Config{
		OllamaURL:    getEnv("OLLAMA_URL", "http://localhost:11434/api/generate"),
		DefaultModel: getEnv("DEFAULT_MODEL", "tinyllama"),
		ServerPort:   getEnv("SERVER_PORT", ":8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
