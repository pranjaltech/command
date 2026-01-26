package llm

import (
	"strings"
	"testing"

	"command/internal/config"
)

func TestOpenRouterClient_Constructor(t *testing.T) {
	// Test missing API key
	_, err := NewOpenRouterClient("", "", "")
	if err == nil || !strings.Contains(err.Error(), "OPENROUTER_API_KEY not set") {
		t.Errorf("expected API key error, got %v", err)
	}

	// Test with API key
	client, err := NewOpenRouterClient("test-key", "", "")
	if err != nil {
		t.Fatalf("NewOpenRouterClient() error = %v", err)
	}
	if client.model != config.DefaultModelOpenRouter {
		t.Errorf("expected default model %s, got %s", config.DefaultModelOpenRouter, client.model)
	}
}

func TestOpenRouterClient_DefaultBaseURL(t *testing.T) {
	client, err := NewOpenRouterClient("test-key", "", "")
	if err != nil {
		t.Fatalf("NewOpenRouterClient() error = %v", err)
	}
	// The base URL is set internally in the OpenAI client config
	// We can verify the wrapper was created successfully
	if client.OpenAIClient == nil {
		t.Error("expected embedded OpenAIClient to be set")
	}
}

func TestOpenRouterClient_CustomSettings(t *testing.T) {
	client, err := NewOpenRouterClient("test-key", "https://custom.openrouter.ai/v1", "anthropic/claude-3-sonnet")
	if err != nil {
		t.Fatalf("NewOpenRouterClient() error = %v", err)
	}
	if client.model != "anthropic/claude-3-sonnet" {
		t.Errorf("expected model anthropic/claude-3-sonnet, got %s", client.model)
	}
}

func TestOpenRouterClient_DebugOutput(t *testing.T) {
	client, _ := NewOpenRouterClient("test-key", "", "")
	var buf strings.Builder
	client.EnableDebug(&buf)
	if !client.debug {
		t.Error("expected debug to be enabled")
	}
}
