package cache

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	// Same prompt + model = same key
	key1 := GenerateKey("What is Kubernetes?", "tinyllama")
	key2 := GenerateKey("What is Kubernetes?", "tinyllama")
	if key1 != key2 {
		t.Errorf("expected same key, got %s and %s", key1, key2)
	}

	// Case insensitive — same key
	key3 := GenerateKey("what is kubernetes?", "tinyllama")
	if key1 != key3 {
		t.Errorf("expected same key for different casing, got %s and %s", key1, key3)
	}

	// Different model = different key
	key4 := GenerateKey("What is Kubernetes?", "mistral")
	if key1 == key4 {
		t.Errorf("expected different key for different model")
	}

	// Different prompt = different key
	key5 := GenerateKey("What is Docker?", "tinyllama")
	if key1 == key5 {
		t.Errorf("expected different key for different prompt")
	}

	// Print example key so you can see the format
	fmt.Println("Example key:", key1)
}

func TestErrCacheMiss(t *testing.T) {
	cache := &RedisCache{}

	_, err := cache.Get(context.Background(), "any-key")
	if !errors.Is(err, ErrCacheMiss) {
		t.Errorf("expected ErrCacheMiss got %v", err)
	}
}
