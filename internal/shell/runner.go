package shell

import (
	"context"
	"os"
	"os/exec"
)

// Runner executes shell commands.
type Runner interface {
	Run(ctx context.Context, cmd string) error
}

type execRunner struct{}

// NewRunner returns a default shell runner.
func NewRunner() Runner { return execRunner{} }

func (execRunner) Run(ctx context.Context, cmdStr string) error {
	sh := os.Getenv("SHELL")
	if sh == "" {
		sh = "/bin/sh"
	}
	c := exec.CommandContext(ctx, sh, "-c", cmdStr)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	// Ignore history errors so command execution result is returned.
	_ = appendHistory(sh, cmdStr)
	return err
}
