package llm

import (
	"strings"
	"testing"

	"command/internal/config"
)

func TestOllamaClient_Constructor(t *testing.T) {
	// Test without API key (should work for Ollama)
	client, err := NewOllamaClient("", "", "")
	if err != nil {
		t.Fatalf("NewOllamaClient() error = %v", err)
	}
	if client.model != config.DefaultModelOllama {
		t.Errorf("expected default model %s, got %s", config.DefaultModelOllama, client.model)
	}
}

func TestOllamaClient_DefaultBaseURL(t *testing.T) {
	client, err := NewOllamaClient("", "", "")
	if err != nil {
		t.Fatalf("NewOllamaClient() error = %v", err)
	}
	// The base URL is set internally in the OpenAI client config
	// We can verify the wrapper was created successfully
	if client.OpenAIClient == nil {
		t.Error("expected embedded OpenAIClient to be set")
	}
}

func TestOllamaClient_CustomSettings(t *testing.T) {
	client, err := NewOllamaClient("", "http://192.168.1.100:11434/v1", "codellama")
	if err != nil {
		t.Fatalf("NewOllamaClient() error = %v", err)
	}
	if client.model != "codellama" {
		t.Errorf("expected model codellama, got %s", client.model)
	}
}

func TestOllamaClient_DebugOutput(t *testing.T) {
	client, _ := NewOllamaClient("", "", "")
	var buf strings.Builder
	client.EnableDebug(&buf)
	if !client.debug {
		t.Error("expected debug to be enabled")
	}
}

func TestOllamaClient_WithAPIKey(t *testing.T) {
	// Test with explicit API key (for remote Ollama instances)
	client, err := NewOllamaClient("custom-key", "", "")
	if err != nil {
		t.Fatalf("NewOllamaClient() error = %v", err)
	}
	if client.OpenAIClient == nil {
		t.Error("expected embedded OpenAIClient to be set")
	}
}
