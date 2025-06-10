package llm

import (
	"context"
	"strings"
	"testing"

	openai "github.com/sashabaranov/go-openai"

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
				Message: openai.ChatCompletionMessage{Content: "{\"commands\":[\"ls\",\"ls -l\"]}"},
			}},
		},
	}
	client := &OpenAIClient{api: stub}
	env := probe.EnvInfo{OS: "linux"}

	got, err := client.GenerateCommands(context.Background(), "list", env)
	if err != nil {
		t.Fatalf("GenerateCommands() error = %v", err)
	}
	if len(got) != 2 || got[0] != "ls" || got[1] != "ls -l" {
		t.Errorf("unexpected output: %#v", got)
	}
	if stub.req.Model != openai.GPT4oMini {
		t.Errorf("expected model %s", openai.GPT4oMini)
	}
	if len(stub.req.Messages) == 0 || stub.req.Messages[0].Role != openai.ChatMessageRoleSystem {
		t.Fatalf("system message missing")
	}
	if !strings.Contains(stub.req.Messages[0].Content, "Respond with JSON") {
		t.Errorf("system prompt missing instruction: %q", stub.req.Messages[0].Content)
	}
	if stub.req.ResponseFormat == nil ||
		stub.req.ResponseFormat.Type != openai.ChatCompletionResponseFormatTypeJSONObject {
		t.Errorf("expected json response format")
	}
}
