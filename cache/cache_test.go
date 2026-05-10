package cache

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
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

func TestRedisCache(t *testing.T) {

	cache, err := NewRedisCache("127.0.0.1:6379")
	if err != nil {
		t.Skip("Redis not available")
	}
	defer cache.Close()

	ctx := context.Background()
	key := "test:key"
	value := "test response"

	cache.Delete(ctx, key)
	defer cache.Delete(ctx, key)

	_, err = cache.Get(ctx, key)
	if !errors.Is(err, ErrCacheMiss) {
		t.Errorf("expected ErrCacheMiss got %v", err)
	}

	err = cache.Set(ctx, key, value, 5*time.Second)
	if err != nil {
		t.Errorf("expected no error no Set got %v", err)
	}

	got, err := cache.Get(ctx, key)
	if err != nil {
		t.Errorf("expected no error no Get got %v", err)
	}
	if got != value {
		t.Errorf("expected %s got %s", value, got)
	}
}

func TestErrCacheMiss(t *testing.T) {
	err := ErrCacheMiss
	if !errors.Is(err, ErrCacheMiss) {
		t.Errorf("expected ErrCacheMiss got %v", err)
	}
}
