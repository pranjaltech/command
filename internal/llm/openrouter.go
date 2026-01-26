package llm

import (
	"context"
	"fmt"
	"io"

	"command/internal/config"
	"command/internal/probe"
)

const OpenRouterBaseURL = "https://openrouter.ai/api/v1"

// OpenRouterClient wraps OpenAIClient for OpenRouter compatibility.
// OpenRouter provides an OpenAI-compatible API endpoint.
type OpenRouterClient struct {
	*OpenAIClient
}

// NewOpenRouterClient constructs an OpenRouter-based LLM client.
func NewOpenRouterClient(apiKey, baseURL, model string) (*OpenRouterClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("OPENROUTER_API_KEY not set")
	}
	if baseURL == "" {
		baseURL = OpenRouterBaseURL
	}
	if model == "" {
		model = config.DefaultModelOpenRouter
	}

	openaiClient, err := NewOpenAIClient(apiKey, baseURL, model)
	if err != nil {
		return nil, err
	}

	return &OpenRouterClient{OpenAIClient: openaiClient}, nil
}

// GenerateCommands delegates to the embedded OpenAIClient.
func (c *OpenRouterClient) GenerateCommands(ctx context.Context, prompt string, env probe.EnvInfo) ([]string, error) {
	return c.OpenAIClient.GenerateCommands(ctx, prompt, env)
}

// EnableDebug delegates to the embedded OpenAIClient.
func (c *OpenRouterClient) EnableDebug(w io.Writer) {
	c.OpenAIClient.EnableDebug(w)
}
