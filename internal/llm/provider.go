package llm

import (
	"context"
	"fmt"
)

// NewClient returns an LLM client for the given provider.
func NewClient(provider, apiKey, apiURL, model string) (Client, error) {
	switch provider {
	case "openai":
		return NewOpenAIClient(apiKey, apiURL, model)
	case "anthropic":
		return NewAnthropicClient(apiKey, apiURL, model)
	case "openrouter":
		return NewOpenRouterClient(apiKey, apiURL, model)
	case "gemini":
		// Gemini client creation requires context; use background context here.
		// The actual request context is passed to GenerateCommands.
		return NewGeminiClient(context.Background(), apiKey, apiURL, model)
	case "ollama":
		return NewOllamaClient(apiKey, apiURL, model)
	default:
		return nil, fmt.Errorf("provider %s not supported", provider)
	}
}
