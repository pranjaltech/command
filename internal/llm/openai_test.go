package llm

import (
	"context"
	"errors"
	"strings"
	"testing"

	openai "github.com/sashabaranov/go-openai"

	"command/internal/config"
	"command/internal/probe"
)

type stubChat struct {
	req  openai.ChatCompletionRequest
	resp openai.ChatCompletionResponse
	err  error
}

func (s *stubChat) CreateChatCompletion(
	ctx context.Context,
	req openai.ChatCompletionRequest,
) (openai.ChatCompletionResponse, error) {
	s.req = req
	return s.resp, s.err
}

func TestOpenAIClient_GenerateCommands(t *testing.T) {
	stub := &stubChat{
		resp: openai.ChatCompletionResponse{
			Choices: []openai.ChatCompletionChoice{{
				Message: openai.ChatCompletionMessage{
					Content: "{\"commands\":[{\"command\":\"ls\"},{\"command\":\"ls -l\"}]}",
				},
			}},
		},
	}
	client := &OpenAIClient{api: stub, model: config.DefaultModel, temperature: config.DefaultTemperature}
	env := probe.EnvInfo{OS: "linux"}

	got, err := client.GenerateCommands(context.Background(), "list", env)
	if err != nil {
		t.Fatalf("GenerateCommands() error = %v", err)
	}
	if len(got) != 2 || got[0] != "ls" || got[1] != "ls -l" {
		t.Errorf("unexpected output: %#v", got)
	}
	if stub.req.Model != config.DefaultModel {
		t.Errorf("expected model %s", config.DefaultModel)
	}
	if len(stub.req.Messages) == 0 || stub.req.Messages[0].Role != openai.ChatMessageRoleSystem {
		t.Fatalf("system message missing")
	}
	if !strings.Contains(stub.req.Messages[0].Content, "Output contract") ||
		!strings.Contains(stub.req.Messages[0].Content, "Few-shot examples") {
		t.Errorf("system prompt missing sections: %q", stub.req.Messages[0].Content)
	}
	if stub.req.ResponseFormat == nil ||
		stub.req.ResponseFormat.Type != openai.ChatCompletionResponseFormatTypeJSONObject {
		t.Errorf("expected json response format")
	}
}

func TestOpenAIClient_GenerateCommands_StringArray(t *testing.T) {
	stub := &stubChat{
		resp: openai.ChatCompletionResponse{
			Choices: []openai.ChatCompletionChoice{{
				Message: openai.ChatCompletionMessage{Content: "{\"commands\":[\"ls\",\"ls -l\"]}"},
			}},
		},
	}
	client := &OpenAIClient{api: stub, model: config.DefaultModel, temperature: config.DefaultTemperature}
	env := probe.EnvInfo{OS: "linux"}

	got, err := client.GenerateCommands(context.Background(), "list", env)
	if err != nil {
		t.Fatalf("GenerateCommands() error = %v", err)
	}
	if len(got) != 2 || got[0] != "ls" || got[1] != "ls -l" {
		t.Errorf("unexpected output: %#v", got)
	}
}

func TestOpenAIClient_GenerateCommands_APIError(t *testing.T) {
	apiErr := &openai.APIError{HTTPStatusCode: 400, Message: "bad"}
	stub := &stubChat{err: apiErr}
	client := &OpenAIClient{api: stub, model: config.DefaultModel, temperature: config.DefaultTemperature}
	env := probe.EnvInfo{OS: "linux"}
	_, err := client.GenerateCommands(context.Background(), "", env)
	if err == nil || !strings.Contains(err.Error(), "openai request failed") {
		t.Fatalf("expected wrapped API error, got %v", err)
	}
}

func TestOpenAIClient_DebugOutput(t *testing.T) {
	stub := &stubChat{
		resp: openai.ChatCompletionResponse{
			Choices: []openai.ChatCompletionChoice{{
				Message: openai.ChatCompletionMessage{Content: "{\"commands\":[\"ls\"]}"},
			}},
		},
	}
	client := &OpenAIClient{api: stub, model: config.DefaultModel, temperature: config.DefaultTemperature}
	var buf strings.Builder
	client.EnableDebug(&buf)
	env := probe.EnvInfo{OS: "linux"}
	_, err := client.GenerateCommands(context.Background(), "list", env)
	if err != nil {
		t.Fatalf("GenerateCommands() error = %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "llm system prompt:") || !strings.Contains(out, "llm user prompt: list") ||
		!strings.Contains(out, "llm raw response") {
		t.Errorf("debug output missing, got: %s", out)
	}
}

func TestOpenAIClient_NeedClarification(t *testing.T) {
	stub := &stubChat{
		resp: openai.ChatCompletionResponse{
			Choices: []openai.ChatCompletionChoice{{
				Message: openai.ChatCompletionMessage{
					Content: "{\"commands\":[],\"need_clarification\":\"which dir?\"}",
				},
			}},
		},
	}
	client := &OpenAIClient{api: stub, model: config.DefaultModel, temperature: config.DefaultTemperature}
	env := probe.EnvInfo{OS: "linux"}
	_, err := client.GenerateCommands(context.Background(), "list", env)
	var nc NeedClarificationError
	if !errors.As(err, &nc) || !strings.Contains(nc.Question, "which dir") {
		t.Fatalf("expected clarification error, got %v", err)
	}
}
