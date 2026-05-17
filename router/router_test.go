package router

import (
	"context"
	"inference-gateway/types"
	"testing"
)

func TestIsShortPrompt(t *testing.T) {
	if !IsShortPrompt("hi") {
		t.Error("expected short prompt")
	}
	if IsShortPrompt("This is a very long prompt that exceeds fifty characters easily") {
		t.Error("expected long prompt")
	}
}

func TestContainsCode(t *testing.T) {
	if !ContainsCode("write a function to sort a list") {
		t.Error("expected code prompt")
	}
	if ContainsCode("What is Kubernetes?") {
		t.Error("expected non-code prompt")
	}
}

func TestRouterFallback(t *testing.T) {
	// Mock provider for testing
	mock := &mockProvider{}
	r := New(mock, "tinyllama")

	provider, model := r.Route("hi")
	if provider != mock {
		t.Error("expected fallback provider")
	}
	if model != "tinyllama" {
		t.Errorf("expected tinyllama got %s", model)
	}
}

func TestRouterFirstMatchWins(t *testing.T) {
	mock1 := &mockProvider{}
	mock2 := &mockProvider{}
	fallback := &mockProvider{}

	r := New(fallback, "tinyllama")
	r.AddRule(Rule{
		Name:      "short",
		Condition: IsShortPrompt,
		Provider:  mock1,
		Model:     "fast-model",
	})
	r.AddRule(Rule{
		Name:      "long",
		Condition: IsLongPrompt,
		Provider:  mock2,
		Model:     "smart-model",
	})

	// Short prompt → mock1
	provider, model := r.Route("hi")
	if provider != mock1 {
		t.Error("expected mock1 for short prompt")
	}
	if model != "fast-model" {
		t.Errorf("expected fast-model got %s", model)
	}

	// Long prompt → mock2
	provider, model = r.Route("This is a very long prompt that exceeds fifty characters easily")
	if provider != mock2 {
		t.Error("expected mock2 for long prompt")
	}
	if model != "smart-model" {
		t.Errorf("expected smart-model got %s", model)
	}
}

// mockProvider satisfies interfaces.ModelProvider for testing
type mockProvider struct{}

func (m *mockProvider) Generate(prompt, model string) (*types.OllamaResponse, error) {
	return &types.OllamaResponse{Response: "mock", Model: model}, nil
}
func (m *mockProvider) Ping(ctx context.Context) error { return nil }
func (m *mockProvider) GetModels() (*types.OllamaModelsResponse, error) {
	return &types.OllamaModelsResponse{}, nil
}
