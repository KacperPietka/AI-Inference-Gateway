package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"inference-gateway/types"
)

type OllamaClient struct {
	URL        string
	HTTPClient *http.Client
}

// Creates a new client with a timeout
func NewOllamaClient(url string) *OllamaClient {
	return &OllamaClient{
		URL: url,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// Ping checks if Ollama is reachable
func (c *OllamaClient) Ping(ctx context.Context) error {
	url := strings.Replace(c.URL, "/api/generate", "/", 1)

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)

	if err != nil {
		return fmt.Errorf("failed to build ping request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("ollama unreachable: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

func (c *OllamaClient) GetModels() (*types.OllamaModelsReponse, error) {
	// This is a GET request - simpler than Generate
	// Ollama's models endpoint is /api/tags
	url := strings.Replace(c.URL, "/api/generate", "/api/tags", 1)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call ollama models: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama returned statis %d", resp.StatusCode)
	}

	var modelsResp types.OllamaModelsReponse
	err = json.NewDecoder(resp.Body).Decode(&modelsResp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode models response: %w", err)
	}

	return &modelsResp, nil
}

func (c *OllamaClient) Generate(prompt, model string) (*types.OllamaResponse, error) {
	// Build request
	ollamaReq := types.OllamaRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	}

	body, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Call Ollama
	resp, err := c.HTTPClient.Post(c.URL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to call ollama: %w", err)
	}
	defer resp.Body.Close()

	// Check status code explicitly
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama returned status %d", resp.StatusCode)
	}

	// Decode response
	var ollamaResp types.OllamaResponse
	err = json.NewDecoder(resp.Body).Decode(&ollamaResp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ollama response: %w", err)
	}

	return &ollamaResp, nil
}
