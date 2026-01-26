package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"

	"command/internal/config"
	"command/internal/log"
	"command/internal/probe"
)

// GeminiClient implements Client using the official Google Generative AI SDK.
type GeminiClient struct {
	client *genai.Client
	model  string
	debug  bool
}

// NewGeminiClient constructs a Gemini-based LLM client.
func NewGeminiClient(ctx context.Context, apiKey, baseURL, model string) (*GeminiClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY not set")
	}

	opts := []option.ClientOption{
		option.WithAPIKey(apiKey),
	}
	if baseURL != "" {
		opts = append(opts, option.WithEndpoint(baseURL))
	}

	client, err := genai.NewClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("create gemini client: %w", err)
	}

	if model == "" {
		model = config.DefaultModelGemini
	}

	return &GeminiClient{
		client: client,
		model:  model,
	}, nil
}

// EnableDebug turns on verbose logging to the provided writer.
func (c *GeminiClient) EnableDebug(w io.Writer) {
	c.debug = true
	log.Enable(w)
}

// Close releases resources associated with the client.
func (c *GeminiClient) Close() error {
	return c.client.Close()
}

// GenerateCommands returns command suggestions from Gemini.
func (c *GeminiClient) GenerateCommands(ctx context.Context, prompt string, env probe.EnvInfo) ([]string, error) {
	sysPrompt, err := BuildSystemPrompt(env)
	if err != nil {
		return nil, err
	}

	// Gemini doesn't have native JSON mode like OpenAI, so we add explicit instructions
	jsonInstruction := "\n\nIMPORTANT: Respond ONLY with valid JSON matching the schema above. No markdown code blocks, no prose, no explanations."
	fullSystemPrompt := sysPrompt + jsonInstruction

	if c.debug {
		log.Debugf("llm system prompt: %s", fullSystemPrompt)
		log.Debugf("llm user prompt: %s", prompt)
		if data, err := json.MarshalIndent(env, "", "  "); err == nil {
			log.Debugf("llm env: %s", data)
		}
	}

	model := c.client.GenerativeModel(c.model)
	model.SystemInstruction = genai.NewUserContent(genai.Text(fullSystemPrompt))

	// Configure for JSON output
	model.ResponseMIMEType = "application/json"

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if c.debug && err != nil {
		log.Debugf("llm error: %v", err)
	}
	if err != nil {
		return nil, fmt.Errorf("gemini request failed: %w", err)
	}

	// Extract text content from response
	var content string
	if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
		for _, part := range resp.Candidates[0].Content.Parts {
			if text, ok := part.(genai.Text); ok {
				content = string(text)
				break
			}
		}
	}

	if content == "" {
		return nil, fmt.Errorf("no content in response")
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
