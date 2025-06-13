package llm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	openai "github.com/sashabaranov/go-openai"

	"command/internal/config"
	"command/internal/log"
	"command/internal/probe"
)

const promptTemplate = `You are **cmd**, a terminal-command assistant.

## 1. Environment (immutable for this session)
%s

## 2. Dynamic context (inject when available)
### Installed binaries (truncated PATH scan)
%s

### Top-level files in $PWD
%s

### Git status (if inside a repo)
%s

### User aliases & functions
%s

## 3. Output contract  -- obey exactly
- Respond ONLY with valid JSON—no markdown, no prose.
- Schema:
  {
    "commands": ["<cmd1>", "<cmd2>", "<cmd3>"],   // ≤3, most relevant first
    "need_clarification": "<question or null>",
    "notes": "<one-sentence rationale or null>"
  }
- Never output destructive commands (e.g. 'rm -rf /') unless the user explicitly requests them and they are clearly marked by adding "dangerous": true alongside that entry.
- If the request is ambiguous, leave "commands" empty and set "need_clarification".

## 4. Style rules
- Fish syntax by default.
- Chain commands with && only when later commands depend on earlier success.
- Prefer concise flags (ls -la > ls --all --long).

## 5. Few-shot examples
U: list go files in current dir
A: {"commands":["ls *.go"],"need_clarification":null,"notes":null}

U: initialise a git repo and push first commit
A: {"commands":["git init && git add . && git commit -m \"init\" && git branch -M main && git remote add origin <url> && git push -u origin main"],"need_clarification":null,"notes":"uses main branch by default"}

Begin!`

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
}

// NeedClarificationError is returned when the model asks a follow-up question
// instead of providing commands.
type NeedClarificationError struct{ Question string }

func (e NeedClarificationError) Error() string {
	if e.Question == "" {
		return "clarification requested"
	}
	return "clarification requested: " + e.Question
}

func buildSystemPrompt(env probe.EnvInfo) (string, error) {
	envJSON, err := json.Marshal(env)
	if err != nil {
		return "", fmt.Errorf("marshal env: %w", err)
	}
	files, _ := os.ReadDir(env.WorkDir)
	var names []string
	for i, f := range files {
		if i >= 20 {
			break
		}
		names = append(names, f.Name())
	}
	fileList := strings.Join(names, " ")

	bins := uniqueBinaries(20)

	gitStatus := fmt.Sprintf("{\"root\":%q,\"branch\":%q,\"dirty\":%t}", env.GitRoot, env.GitBranch, env.GitDirty)

	return fmt.Sprintf(promptTemplate, envJSON, bins, fileList, gitStatus, "[]"), nil
}

var (
	binsOnce sync.Once
	binsList []string
)

func uniqueBinaries(limit int) string {
	binsOnce.Do(func() {
		seen := make(map[string]struct{})
		for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
			entries, err := os.ReadDir(dir)
			if err != nil {
				continue
			}
			for _, e := range entries {
				if len(binsList) >= limit {
					return
				}
				name := e.Name()
				if _, ok := seen[name]; ok {
					continue
				}
				seen[name] = struct{}{}
				binsList = append(binsList, name)
			}
		}
	})
	if len(binsList) > limit {
		return strings.Join(binsList[:limit], " ")
	}
	return strings.Join(binsList, " ")
}

// EnableDebug turns on verbose logging to the provided writer.
func (c *OpenAIClient) EnableDebug(w io.Writer) {
	c.debug = true
	log.Enable(w)
}

// NewOpenAIClient constructs an OpenAI-based LLM client.
func NewOpenAIClient(apiKey, baseURL, model string, temperature float32) (*OpenAIClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY not set")
	}
	cfg := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		cfg.BaseURL = baseURL
	}
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
	}, nil
}

// GenerateCommand returns a command suggestion from the LLM.
func (c *OpenAIClient) GenerateCommands(ctx context.Context, prompt string, env probe.EnvInfo) ([]string, error) {
	sysPrompt, err := buildSystemPrompt(env)
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
	var raw struct {
		Commands          []json.RawMessage `json:"commands"`
		NeedClarification json.RawMessage   `json:"need_clarification"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(resp.Choices[0].Message.Content)), &raw); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	if len(raw.NeedClarification) > 0 && string(raw.NeedClarification) != "null" {
		var q string
		if err := json.Unmarshal(raw.NeedClarification, &q); err == nil && q != "" {
			return nil, NeedClarificationError{Question: q}
		}
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
		log.Debugf("llm parsed commands: %v", out)
	}
	return out, nil
}
