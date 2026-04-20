package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	gwerrors "inference-gateway/errors"
	"inference-gateway/models"
	"inference-gateway/types"
)

type ModelsHandler struct {
	Ollama *models.OllamaClient
}

func NewModelsHandler(ollama *models.OllamaClient) *ModelsHandler {
	return &ModelsHandler{Ollama: ollama}
}

func (h *ModelsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// Only allow GET
	if r.Method != http.MethodGet {
		writeError(w, gwerrors.ErrMethodNotAllowed)
		return
	}

	log.Println("fetching available models from ollama")

	ollamaResp, err := h.Ollama.GetModels()
	if err != nil {
		log.Printf("models error: %v", err)
		writeError(w, gwerrors.New(gwerrors.ErrModelUnavailable, err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.ModelsReponse{
		Models: ollamaResp.Models,
	})
}
