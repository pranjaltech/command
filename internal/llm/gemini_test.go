package llm

import (
	"context"
	"strings"
	"testing"

	"command/internal/config"
)

// Note: The Gemini SDK client requires context and is not easily mockable.
// These tests focus on constructor validation. Integration tests with real API
// calls should be done separately.

func TestGeminiClient_Constructor(t *testing.T) {
	ctx := context.Background()

	// Test missing API key
	_, err := NewGeminiClient(ctx, "", "", "")
	if err == nil || !strings.Contains(err.Error(), "GEMINI_API_KEY not set") {
		t.Errorf("expected API key error, got %v", err)
	}
}

func TestGeminiClient_DefaultModel(t *testing.T) {
	// We can't fully test without a real API key, but we can verify defaults
	// are set correctly in the config
	if config.DefaultModelGemini != "gemini-3-flash" {
		t.Errorf("expected default model gemini-3-flash, got %s", config.DefaultModelGemini)
	}
}
