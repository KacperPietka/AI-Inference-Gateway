package main

import (
	"log"
	"net/http"

	"inference-gateway/config"
	"inference-gateway/handlers"
	"inference-gateway/middleware"
	"inference-gateway/models"
)

func main() {
	// Load config
	cfg := config.Load()

	// Build dependencies
	ollamaClient := models.NewOllamaClient(cfg.OllamaURL)
	generateHandler := handlers.NewGenerateHandler(ollamaClient, cfg.DefaultModel)

	// Register routes
	http.HandleFunc("/health", middleware.Logger(handlers.HealthHandler))
	http.HandleFunc("/generate", middleware.Logger(generateHandler.Handle))

	// Start server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(cfg.ServerPort, nil))
}
