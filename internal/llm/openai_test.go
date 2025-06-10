package llm

import (
	"context"
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

func TestOpenAIClient_GenerateCommand(t *testing.T) {
	stub := &stubChat{
		resp: openai.ChatCompletionResponse{
			Choices: []openai.ChatCompletionChoice{{
				Message: openai.ChatCompletionMessage{Content: "ls"},
			}},
		},
	}
	client := &OpenAIClient{api: stub}
	env := probe.EnvInfo{OS: "linux"}

	got, err := client.GenerateCommand(context.Background(), "list", env)
	if err != nil {
		t.Fatalf("GenerateCommand() error = %v", err)
	}
	if got != "ls" {
		t.Errorf("want ls, got %q", got)
	}
	if stub.req.Model != openai.GPT4oMini {
		t.Errorf("expected model %s", openai.GPT4oMini)
	}
	if len(stub.req.Messages) == 0 || stub.req.Messages[0].Role != openai.ChatMessageRoleSystem {
		t.Fatalf("system message missing")
	}
}
