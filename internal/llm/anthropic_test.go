package llm

import (
	"errors"
	"strings"
	"testing"

	"command/internal/config"
)

// Note: The Anthropic SDK client is not easily mockable due to its structure.
// These tests focus on the parsing and error handling logic which uses the shared
// common.go utilities. Integration tests with real API calls should be done separately.

func TestAnthropicClient_Constructor(t *testing.T) {
	// Test missing API key
	_, err := NewAnthropicClient("", "", "")
	if err == nil || !strings.Contains(err.Error(), "ANTHROPIC_API_KEY not set") {
		t.Errorf("expected API key error, got %v", err)
	}

	// Test with API key
	client, err := NewAnthropicClient("test-key", "", "")
	if err != nil {
		t.Fatalf("NewAnthropicClient() error = %v", err)
	}
	if client.model != config.DefaultModelAnthropic {
		t.Errorf("expected default model %s, got %s", config.DefaultModelAnthropic, client.model)
	}
}

func TestAnthropicClient_Constructor_CustomSettings(t *testing.T) {
	client, err := NewAnthropicClient("test-key", "https://custom.api.com", "claude-3-haiku")
	if err != nil {
		t.Fatalf("NewAnthropicClient() error = %v", err)
	}
	if client.model != "claude-3-haiku" {
		t.Errorf("expected model claude-3-haiku, got %s", client.model)
	}
}

func TestAnthropicClient_DebugOutput(t *testing.T) {
	client, _ := NewAnthropicClient("test-key", "", "")
	var buf strings.Builder
	client.EnableDebug(&buf)
	if !client.debug {
		t.Error("expected debug to be enabled")
	}
}

// TestParseCommands_AnthropicFormat tests parsing responses in formats Anthropic might return
func TestParseCommands_AnthropicFormat(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		want     []string
		wantErr  bool
		errType  error
	}{
		{
			name:    "simple string array",
			content: `{"commands":["ls","pwd"],"need_clarification":null,"notes":null}`,
			want:    []string{"ls", "pwd"},
		},
		{
			name:    "object format",
			content: `{"commands":[{"command":"ls"},{"command":"pwd"}],"need_clarification":null}`,
			want:    []string{"ls", "pwd"},
		},
		{
			name:    "clarification request",
			content: `{"commands":[],"need_clarification":"which directory?"}`,
			wantErr: true,
			errType: NeedClarificationError{},
		},
		{
			name:    "with notes",
			content: `{"commands":["git status"],"need_clarification":null,"notes":"shows working tree status"}`,
			want:    []string{"git status"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCommands(tt.content)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseCommands() expected error, got nil")
					return
				}
				if tt.errType != nil {
					var nc NeedClarificationError
					if !errors.As(err, &nc) {
						t.Errorf("ParseCommands() error type = %T, want NeedClarificationError", err)
					}
				}
				return
			}
			if err != nil {
				t.Errorf("ParseCommands() error = %v", err)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("ParseCommands() got %v, want %v", got, tt.want)
				return
			}
			for i, cmd := range got {
				if cmd != tt.want[i] {
					t.Errorf("ParseCommands()[%d] = %q, want %q", i, cmd, tt.want[i])
				}
			}
		})
	}
}
