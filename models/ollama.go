package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"inference-gateway/types"
	"net/http"
	"time"
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
