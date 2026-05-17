package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"inference-gateway/cache"
	gwerrors "inference-gateway/errors"
	"inference-gateway/router"
	"inference-gateway/types"
)

type GenerateHandler struct {
	Router       *router.Router
	DefaultModel string
	Logger       *slog.Logger
	Cache        cache.Cache
	CacheTTL     time.Duration
}

func NewGenerateHandler(r *router.Router, defaultModel string, logger *slog.Logger, c cache.Cache, cacheTTL time.Duration) *GenerateHandler {
	return &GenerateHandler{
		Router:       r,
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

	// Router decides provider and model
	provider, model := h.Router.Route(req.Prompt)
	if provider == nil {
		writeError(w, gwerrors.ErrModelUnavailable)
		return
	}

	// Respect explicit user model choice
	if req.Model != "" {
		model = req.Model
	}

	ctx := context.WithValue(r.Context(), types.ModelKey, model)
	r = r.WithContext(ctx)

	key := cache.GenerateKey(req.Prompt, model)

	if cached, err := h.Cache.Get(r.Context(), key); err == nil {
		h.Logger.Info("cache hit", "key", key, "model", model)

		var cachedResp types.GenerateResponse
		if err := json.Unmarshal([]byte(cached), &cachedResp); err == nil {
			cachedResp.Cached = true
			respBytes, _ := json.Marshal(cachedResp)
			w.Header().Set("X-Cache", "HIT")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(respBytes)
			return

		}
	}

	resp, err := provider.Generate(req.Prompt, model)
	if err != nil {
		h.Logger.Error("provider error", "error", err)
		writeError(w, gwerrors.New(gwerrors.ErrModelUnavailable, err))
		return
	}

	respBytes, err := json.Marshal(types.GenerateResponse{
		Response: resp.Response,
		Model:    resp.Model,
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
