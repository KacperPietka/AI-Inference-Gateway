package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"inference-gateway/cache"
	gwerrors "inference-gateway/errors"
	"inference-gateway/interfaces"
	"inference-gateway/models"
	"inference-gateway/types"
)

type GenerateHandler struct {
	Ollama       interfaces.ModelProvider
	DefaultModel string
	Logger       *slog.Logger
	Cache        cache.Cache
	CacheTTL     time.Duration
}

func NewGenerateHandler(ollama *models.OllamaClient, defaultModel string, logger *slog.Logger, c cache.Cache, cacheTTL time.Duration) *GenerateHandler {
	return &GenerateHandler{
		Ollama:       ollama,
		DefaultModel: defaultModel,
		Logger:       logger,
		Cache:        c,
		CacheTTL:     cacheTTL,
	}
}

func writeError(w http.ResponseWriter, err *gwerrors.GatewayError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)
	json.NewEncoder(w).Encode(types.ErrorResponse{
		Error: err.Message,
		Code:  err.Code,
	})
}

func (h *GenerateHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, gwerrors.ErrMethodNotAllowed)
		return
	}

	var req types.GenerateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, gwerrors.ErrInvalidRequest)
		return
	}

	if req.Prompt == "" {
		writeError(w, gwerrors.ErrPromptRequired)
		return
	}

	if req.Model == "" {
		req.Model = h.DefaultModel
	}

	ctx := context.WithValue(r.Context(), types.ModelKey, req.Model)
	r = r.WithContext(ctx)

	key := cache.GenerateKey(req.Prompt, req.Model)

	if cached, err := h.Cache.Get(r.Context(), key); err == nil {
		h.Logger.Info("cache hit", "key", key, "model", req.Model)

		w.Header().Set("X-Cache", "HIT")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(cached))
		return
	}

	h.Logger.Info("calling ollama",
		"model", req.Model,
		"prompt", req.Prompt,
	)

	ollamaResp, err := h.Ollama.Generate(req.Prompt, req.Model)
	if err != nil {
		// Wrap the underluing error with context
		h.Logger.Error("ollama error", "error", err)
		writeError(w, gwerrors.New(gwerrors.ErrModelUnavailable, err))
		return
	}

	respBytes, err := json.Marshal(types.GenerateResponse{
		Response: ollamaResp.Response,
		Model:    ollamaResp.Model,
		Cached:   false,
	})
	if err != nil {
		writeError(w, gwerrors.ErrInternalServer)
		return
	}

	if err := h.Cache.Set(r.Context(), key, string(respBytes), h.CacheTTL); err != nil {
		h.Logger.Error("failed to cache response", "error", err)
	}

	w.Header().Set("X-Cache", "MISS")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}
