package llm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"command/internal/config"
	"command/internal/log"
	"command/internal/probe"
)

// ChatClient abstracts the OpenAI client for testability.
type ChatClient interface {
	CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
}

// OpenAIClient implements Client using the official OpenAI SDK.
type OpenAIClient struct {
	api   ChatClient
	model string
	debug bool
}

// EnableDebug turns on verbose logging to the provided writer.
func (c *OpenAIClient) EnableDebug(w io.Writer) {
	c.debug = true
	log.Enable(w)
}

// NewOpenAIClient constructs an OpenAI-based LLM client.
func NewOpenAIClient(apiKey, baseURL, model string) (*OpenAIClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY not set")
	}
	cfg := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		cfg.BaseURL = baseURL
	}
	if model == "" {
		model = config.DefaultModelOpenAI
	}
	return &OpenAIClient{
		api:   openai.NewClientWithConfig(cfg),
		model: model,
	}, nil
}

// GenerateCommands returns command suggestions from the LLM.
func (c *OpenAIClient) GenerateCommands(ctx context.Context, prompt string, env probe.EnvInfo) ([]string, error) {
	sysPrompt, err := BuildSystemPrompt(env)
	if err != nil {
		return nil, err
	}
	if c.debug {
		log.Debugf("llm system prompt: %s", sysPrompt)
		log.Debugf("llm user prompt: %s", prompt)
		if data, err := json.MarshalIndent(env, "", "  "); err == nil {
			log.Debugf("llm env: %s", data)
		}
	}
	req := openai.ChatCompletionRequest{
		Model: c.model,
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
		log.Debugf("llm error: %v", err)
	}
	if err != nil {
		var apiErr *openai.APIError
		if errors.As(err, &apiErr) {
			return nil, fmt.Errorf("openai request failed: %s (status %d)", apiErr.Message, apiErr.HTTPStatusCode)
		}
		return nil, fmt.Errorf("chat completion: %w", err)
	}
	if c.debug {
		log.Debugf("llm raw response: %s", strings.TrimSpace(resp.Choices[0].Message.Content))
	}
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned")
	}

	cmds, err := ParseCommands(resp.Choices[0].Message.Content)
	if c.debug && err == nil {
		log.Debugf("llm parsed commands: %v", cmds)
	}
	return cmds, err
}
