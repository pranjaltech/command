package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"command/internal/probe"
)

// ChatClient abstracts the OpenAI client for testability.
type ChatClient interface {
	CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
}

// OpenAIClient implements Client using the official OpenAI SDK.
type OpenAIClient struct {
	api ChatClient
}

// NewOpenAIClient constructs an OpenAI-based LLM client.
func NewOpenAIClient(apiKey string) (*OpenAIClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY not set")
	}
	cfg := openai.DefaultConfig(apiKey)
	return &OpenAIClient{api: openai.NewClientWithConfig(cfg)}, nil
}

// GenerateCommand returns a command suggestion from the LLM.
func (c *OpenAIClient) GenerateCommands(ctx context.Context, prompt string, env probe.EnvInfo) ([]string, error) {
	envJSON, err := json.Marshal(env)
	if err != nil {
		return nil, fmt.Errorf("marshal env: %w", err)
	}
	req := openai.ChatCompletionRequest{
		Model:       openai.GPT4oMini,
		Temperature: 0.2,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a CLI assistant. Environment: " + string(envJSON),
			},
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	}
	resp, err := c.api.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("chat completion: %w", err)
	}
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned")
	}
	lines := strings.Split(resp.Choices[0].Message.Content, "\n")
	out := make([]string, 0, len(lines))
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		out = append(out, ln)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no commands parsed")
	}
	return out, nil
}
