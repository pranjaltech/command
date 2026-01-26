package llm

import (
	"context"
	"io"

	openai "github.com/sashabaranov/go-openai"

	"command/internal/config"
	"command/internal/log"
	"command/internal/probe"
)

const OllamaBaseURL = "http://localhost:11434/v1"

// OllamaClient wraps OpenAIClient for Ollama compatibility.
// Ollama provides an OpenAI-compatible API endpoint locally.
type OllamaClient struct {
	*OpenAIClient
}

// NewOllamaClient constructs an Ollama-based LLM client.
// Note: Ollama running locally typically doesn't require an API key.
func NewOllamaClient(apiKey, baseURL, model string) (*OllamaClient, error) {
	if baseURL == "" {
		baseURL = OllamaBaseURL
	}
	if model == "" {
		model = config.DefaultModelOllama
	}

	// Ollama doesn't require an API key for local usage, but the OpenAI SDK
	// requires a non-empty key. Use a placeholder if none provided.
	if apiKey == "" {
		apiKey = "ollama"
	}

	cfg := openai.DefaultConfig(apiKey)
	cfg.BaseURL = baseURL

	client := &OpenAIClient{
		api:   openai.NewClientWithConfig(cfg),
		model: model,
	}

	return &OllamaClient{OpenAIClient: client}, nil
}

// GenerateCommands delegates to the embedded OpenAIClient.
func (c *OllamaClient) GenerateCommands(ctx context.Context, prompt string, env probe.EnvInfo) ([]string, error) {
	return c.OpenAIClient.GenerateCommands(ctx, prompt, env)
}

// EnableDebug delegates to the embedded OpenAIClient.
func (c *OllamaClient) EnableDebug(w io.Writer) {
	c.debug = true
	log.Enable(w)
}
