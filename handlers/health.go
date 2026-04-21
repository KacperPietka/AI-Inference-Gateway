package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"inference-gateway/models"
	"inference-gateway/types"
)

var startTime = time.Now()

type HealthHandler struct {
	Ollama       *models.OllamaClient
	DefaultModel string
}

func NewHealthHandler(ollama *models.OllamaClient, defaultModel string) *HealthHandler {
	return &HealthHandler{
		Ollama:       ollama,
		DefaultModel: defaultModel,
	}
}

func (h *HealthHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	ollamaStatus := "reachable"
	overallStatus := "healthy"

	if err := h.Ollama.Ping(ctx); err != nil {
		ollamaStatus = "unreachable"
		overallStatus = "degraded"
	}

	uptime := time.Since(startTime).Round(time.Second).String()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(types.HealthResponse{
		Status: overallStatus,
		Uptime: uptime,
		Ollama: types.OllamaHealth{
			Status: ollamaStatus,
			Model:  h.DefaultModel,
		},
	})
}
