package interfaces

import (
	"context"
	"inference-gateway/types"
)

// ModelProvider defines what any model client must implement
type ModelProvider interface {
	Generate(prompt, model string) (*types.OllamaResponse, error)
	Ping(ctx context.Context) error
	GetModels() (*types.OllamaModelsResponse, error)
}
