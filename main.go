package main

import (
	"context"
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

func main() {
	// Load config
	cfg := config.Load()

	// Build dependencies
	ollamaClient := models.NewOllamaClient(cfg.OllamaURL)
	generateHandler := handlers.NewGenerateHandler(ollamaClient, cfg.DefaultModel)
	modelsHandler := handlers.NewModelsHandler(ollamaClient)

	// Register routes
	http.HandleFunc("/health", middleware.Logger(handlers.HealthHandler))
	http.HandleFunc("/generate", middleware.Logger(generateHandler.Handle))
	http.HandleFunc("/models", middleware.Logger(modelsHandler.Handle))

	server := &http.Server{
		Addr: cfg.ServerPort,
	}

	// go launches a goroutine, the server runs in the background (immediately executed)
	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Creates a channel that carries os.Signal values (buffered)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// blocks yntil a signal arrives (keeps running after shutting down to finish a request)
	<-quit

	log.Println("Shutting down server, finishing in-progress requests...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}

	log.Println("Server exited cleanly")
}
