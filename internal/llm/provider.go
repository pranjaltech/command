package llm

import (
	"fmt"
)

// NewClient returns an LLM client for the given provider.
func NewClient(provider, apiKey, apiURL, model string, temperature float32) (Client, error) {
	switch provider {
	case "openai":
		return NewOpenAIClient(apiKey, apiURL, model, temperature)
	default:
		return nil, fmt.Errorf("provider %s not supported", provider)
	}
}
