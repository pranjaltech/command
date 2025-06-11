package llm

import (
	"context"

	"command/internal/probe"
)

// Client generates command suggestions given a prompt and environment.
type Client interface {
	GenerateCommands(ctx context.Context, prompt string, env probe.EnvInfo) ([]string, error)
}
