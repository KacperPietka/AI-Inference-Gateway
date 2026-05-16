package models

import (
	"context"
	"os"
	"testing"
)

func TestGeminiGenerate(t *testing.T) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		t.Skip("GEMINI_API_KEY not set")
	}

	client := NewGeminiClient(apiKey, "gemini-3-flash-preview")

	resp, err := client.Generate("What is Kubernetes in one sentence?", "")
	if err != nil {
		t.Fatalf("expected no error got %v", err)
	}

	if resp.Response == "" {
		t.Error("expected non-empty response")
	}

	t.Logf("Gemini response: %s", resp.Response)
}

func TestGeminiPing(t *testing.T) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		t.Skip("GEMINI_API_KEY not set")
	}

	client := NewGeminiClient(apiKey, "gemini-3-flash-preview")

	if err := client.Ping(context.Background()); err != nil {
		t.Fatalf("expected ping to succeed got %v", err)
	}
}
