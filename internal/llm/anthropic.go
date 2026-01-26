package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"

	"command/internal/config"
	"command/internal/log"
	"command/internal/probe"
)

// AnthropicClient implements Client using the official Anthropic SDK.
type AnthropicClient struct {
	client anthropic.Client
	model  string
	debug  bool
}

// NewAnthropicClient constructs an Anthropic-based LLM client.
func NewAnthropicClient(apiKey, baseURL, model string) (*AnthropicClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY not set")
	}

	opts := []option.RequestOption{
		option.WithAPIKey(apiKey),
	}
	if baseURL != "" {
		opts = append(opts, option.WithBaseURL(baseURL))
	}

	if model == "" {
		model = config.DefaultModelAnthropic
	}

	return &AnthropicClient{
		client: anthropic.NewClient(opts...),
		model:  model,
	}, nil
}

// EnableDebug turns on verbose logging to the provided writer.
func (c *AnthropicClient) EnableDebug(w io.Writer) {
	c.debug = true
	log.Enable(w)
}

// GenerateCommands returns command suggestions from Claude.
func (c *AnthropicClient) GenerateCommands(ctx context.Context, prompt string, env probe.EnvInfo) ([]string, error) {
	sysPrompt, err := BuildSystemPrompt(env)
	if err != nil {
		return nil, err
	}

	// Anthropic doesn't have native JSON mode, so we add explicit instructions
	jsonInstruction := "\n\nIMPORTANT: Respond ONLY with valid JSON matching the schema above. No markdown code blocks, no prose, no explanations."
	fullSystemPrompt := sysPrompt + jsonInstruction

	if c.debug {
		log.Debugf("llm system prompt: %s", fullSystemPrompt)
		log.Debugf("llm user prompt: %s", prompt)
		if data, err := json.MarshalIndent(env, "", "  "); err == nil {
			log.Debugf("llm env: %s", data)
		}
	}

	message, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.Model(c.model),
		MaxTokens: int64(1024),
		System: []anthropic.TextBlockParam{
			{
				Text: fullSystemPrompt,
				Type: "text",
			},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})

	if c.debug && err != nil {
		log.Debugf("llm error: %v", err)
	}
	if err != nil {
		return nil, fmt.Errorf("anthropic request failed: %w", err)
	}

	// Extract text content from response
	var content string
	for _, block := range message.Content {
		if block.Type == "text" {
			content = block.Text
			break
		}
	}

	if c.debug {
		log.Debugf("llm raw response: %s", strings.TrimSpace(content))
	}

	cmds, err := ParseCommands(content)
	if c.debug && err == nil {
		log.Debugf("llm parsed commands: %v", cmds)
	}
	return cmds, err
}
