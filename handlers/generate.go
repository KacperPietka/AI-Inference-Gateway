package handlers

import (
	"encoding/json"
	"inference-gateway/models"
	"inference-gateway/types"
	"log"
	"net/http"
)

type GenerateHandler struct {
	Ollama       *models.OllamaClient
	DefaultModel string
}

func NewGenerateHandler(ollama *models.OllamaClient, defaultModel string) *GenerateHandler {
	return &GenerateHandler{
		Ollama:       ollama,
		DefaultModel: defaultModel,
	}
}

func writeError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(types.ErrorResponse{Error: message, Code: code})
}

func (h *GenerateHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req types.GenerateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Prompt == "" {
		writeError(w, "prompt is required", http.StatusBadRequest)
		return
	}

	if req.Model == "" {
		req.Model = h.DefaultModel
	}

	log.Printf("calling ollama model=%s prompt=%q", req.Model, req.Prompt)

	ollamaResp, err := h.Ollama.Generate(req.Prompt, req.Model)
	if err != nil {
		log.Printf("ollama error: %v", err)
		writeError(w, "failed to call model", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.GenerateResponse{
		Response: ollamaResp.Response,
		Model:    ollamaResp.Model,
		Cached:   false,
	})
}
