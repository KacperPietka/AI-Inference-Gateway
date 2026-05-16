package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"inference-gateway/interfaces"
	"inference-gateway/types"
)

// Compile time check
var _ interfaces.ModelProvider = (*GeminiClient)(nil)

const geminiBaseURL = "https://generativelanguage.googleapis.com/v1beta/models"

type GeminiClient struct {
	APIKey     string
	Model      string
	HTTPClient *http.Client
}

// Creates a new Gemini Client
func NewGeminiClient(apiKey, model string) *GeminiClient {
	if apiKey == "" {
		return nil
	}
	return &GeminiClient{
		APIKey: apiKey,
		Model:  model,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *GeminiClient) Generate(prompt, model string) (*types.OllamaResponse, error) {
	if model == "" {
		model = c.Model
	}

	geminiReq := types.GeminiRequest{
		Contents: []types.GeminiContent{
			{
				Parts: []types.GeminiPart{
					{Text: prompt},
				},
			},
		},
	}

	body, err := json.Marshal(geminiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal gemini request: %w", err)
	}

	url := fmt.Sprintf("%s/%s:generateContent?key=%s", geminiBaseURL, model, c.APIKey)

	resp, err := c.HTTPClient.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to call gemini: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gemini returned status %d", resp.StatusCode)
	}

	var geminiResp types.GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return nil, fmt.Errorf("failed to decode gemini response : %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("gemini returned empty response")
	}

	// Normalize to standard response
	return &types.OllamaResponse{
		Response: geminiResp.Candidates[0].Content.Parts[0].Text,
		Model:    model,
	}, nil
}

func (c *GeminiClient) Ping(ctx context.Context) error {
	url := fmt.Sprintf("%s?key=%s", geminiBaseURL, c.APIKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to build gemini ping request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("gemini unreachable: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

// Returns available Gemini models
func (c *GeminiClient) GetModels() (*types.OllamaModelsResponse, error) {
	return &types.OllamaModelsResponse{
		Models: []types.OllamaModel{
			{Name: "gemini-2.0-flash", Size: 0},
		},
	}, nil
}
