package llm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
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
	debug       bool
	out         io.Writer
}

// EnableDebug turns on verbose logging to the provided writer.
func (c *OpenAIClient) EnableDebug(w io.Writer) {
	c.debug = true
	if w != nil {
		c.out = w
	}
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
	return &OpenAIClient{
		api:         openai.NewClientWithConfig(cfg),
		model:       model,
		temperature: temperature,
		out:         os.Stderr,
	}, nil
}

// GenerateCommand returns a command suggestion from the LLM.
func (c *OpenAIClient) GenerateCommands(ctx context.Context, prompt string, env probe.EnvInfo) ([]string, error) {
	envJSON, err := json.Marshal(env)
	if err != nil {
		return nil, fmt.Errorf("marshal env: %w", err)
	}
	sysPrompt := "You are a CLI assistant. Environment:" + string(envJSON) +
		". Respond with JSON: {\"commands\": [<cmd>...]} limited to three items."
	if c.debug {
		fmt.Fprintf(c.out, "llm system prompt: %s\n", sysPrompt)
		fmt.Fprintf(c.out, "llm user prompt: %s\n", prompt)
		var pretty []byte
		if p, err := json.MarshalIndent(env, "", "  "); err == nil {
			pretty = p
		} else {
			pretty = envJSON
		}
		fmt.Fprintf(c.out, "llm env: %s\n", pretty)
	}
	req := openai.ChatCompletionRequest{
		Model:       c.model,
		Temperature: c.temperature,
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		},
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: sysPrompt},
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	}
	resp, err := c.api.CreateChatCompletion(ctx, req)
	if c.debug && err != nil {
		fmt.Fprintf(c.out, "llm error: %v\n", err)
	}
	if err != nil {
		var apiErr *openai.APIError
		if errors.As(err, &apiErr) {
			return nil, fmt.Errorf("openai request failed: %s (status %d)", apiErr.Message, apiErr.HTTPStatusCode)
		}
		return nil, fmt.Errorf("chat completion: %w", err)
	}
	if c.debug {
		fmt.Fprintf(c.out, "llm raw response: %s\n", strings.TrimSpace(resp.Choices[0].Message.Content))
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
	if c.debug {
		fmt.Fprintf(c.out, "llm parsed commands: %v\n", out)
	}
	return out, nil
}
