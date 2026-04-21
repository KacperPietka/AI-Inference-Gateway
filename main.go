package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"inference-gateway/config"
	"inference-gateway/handlers"
	"inference-gateway/middleware"
	"inference-gateway/models"
)

func printBanner(cfg *config.Config) {
	fmt.Println()
	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║         Inference Gateway              ║")
	fmt.Println("╚════════════════════════════════════════╝")
	fmt.Printf("→ Version:       %s\n", "v0.1.0")
	fmt.Printf("→ Port:          %s\n", cfg.ServerPort)
	fmt.Printf("→ Ollama URL:    %s\n", cfg.OllamaURL)
	fmt.Printf("→ Default Model: %s\n", cfg.DefaultModel)
	fmt.Println()
	fmt.Println("Routes:")
	fmt.Println("  GET  /health")
	fmt.Println("  GET  /models")
	fmt.Println("  POST /generate")
	fmt.Println()
	fmt.Println("Server ready. Press Ctrl+C to stop.")
	fmt.Println()
}

func main() {
	// Load config
	cfg := config.Load()

	// Created the structured logger once which is shared across all middleware
	logger := middleware.NewLogger()

	// Build dependencies
	ollamaClient := models.NewOllamaClient(cfg.OllamaURL)
	generateHandler := handlers.NewGenerateHandler(ollamaClient, cfg.DefaultModel, logger)
	modelsHandler := handlers.NewModelsHandler(ollamaClient)
	healthHandler := handlers.NewHealthHandler(ollamaClient, cfg.DefaultModel)

	// Register routes and pass logger to each middleware wrapped
	http.HandleFunc("/health", middleware.Logger(logger, healthHandler.Handle))
	http.HandleFunc("/generate", middleware.Logger(logger, generateHandler.Handle))
	http.HandleFunc("/models", middleware.Logger(logger, modelsHandler.Handle))

	server := &http.Server{
		Addr: cfg.ServerPort,
	}

	printBanner(cfg)

	// go launches a goroutine, the server runs in the background (immediately executed)
	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
		}
	}()

	// Creates a channel that carries os.Signal values (buffered)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// blocks yntil a signal arrives (keeps running after shutting down to finish a request)
	<-quit

	logger.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("forced shutdown", "error", err)
	}

	logger.Info("server exited cleanly")
}
