package llm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"command/internal/config"
	"command/internal/probe"
)

// ChatClient abstracts the OpenAI client for testability.
type ChatClient interface {
	CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
}

// OpenAIClient implements Client using the official OpenAI SDK.
type OpenAIClient struct {
	api         ChatClient
	model       string
	temperature float32
}

// NewOpenAIClient constructs an OpenAI-based LLM client.
func NewOpenAIClient(apiKey, model string, temperature float32) (*OpenAIClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY not set")
	}
	cfg := openai.DefaultConfig(apiKey)
	if model == "" {
		model = config.DefaultModel
	}
	if temperature == 0 {
		temperature = config.DefaultTemperature
	}
	return &OpenAIClient{api: openai.NewClientWithConfig(cfg), model: model, temperature: temperature}, nil
}

// GenerateCommand returns a command suggestion from the LLM.
func (c *OpenAIClient) GenerateCommands(ctx context.Context, prompt string, env probe.EnvInfo) ([]string, error) {
	envJSON, err := json.Marshal(env)
	if err != nil {
		return nil, fmt.Errorf("marshal env: %w", err)
	}
	req := openai.ChatCompletionRequest{
		Model:       c.model,
		Temperature: c.temperature,
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		},
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				Content: "You are a CLI assistant. Environment: " + string(
					envJSON,
				) + ". Respond with JSON: {\"commands\": [<cmd>...]} limited to three items.",
			},
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	}
	resp, err := c.api.CreateChatCompletion(ctx, req)
	if err != nil {
		var apiErr *openai.APIError
		if errors.As(err, &apiErr) {
			return nil, fmt.Errorf("openai request failed: %s (status %d)", apiErr.Message, apiErr.HTTPStatusCode)
		}
		return nil, fmt.Errorf("chat completion: %w", err)
	}
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned")
	}
	var raw struct {
		Commands []json.RawMessage `json:"commands"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(resp.Choices[0].Message.Content)), &raw); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	var out []string
	for _, r := range raw.Commands {
		var s string
		if err := json.Unmarshal(r, &s); err == nil {
			out = append(out, s)
			continue
		}
		var obj map[string]interface{}
		if err := json.Unmarshal(r, &obj); err == nil {
			if v, ok := obj["command"].(string); ok {
				out = append(out, v)
				continue
			}
			if v, ok := obj["cmd"].(string); ok {
				out = append(out, v)
				continue
			}
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no commands parsed")
	}
	if len(out) > 3 {
		out = out[:3]
	}
	return out, nil
}
