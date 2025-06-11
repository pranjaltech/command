package shell

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Runner executes shell commands.
type Runner interface {
	Run(ctx context.Context, cmd string) error
}

var commandContext = exec.CommandContext

type execRunner struct{}

// NewRunner returns a default shell runner.
func NewRunner() Runner { return execRunner{} }

func historyFile(shellPath string) string {
	home := os.Getenv("HOME")
	switch filepath.Base(shellPath) {
	case "bash":
		return filepath.Join(home, ".bash_history")
	case "zsh":
		return filepath.Join(home, ".zsh_history")
	case "fish":
		return filepath.Join(home, ".local", "share", "fish", "fish_history")
	default:
		return ""
	}
}

func appendHistory(shellPath, cmdStr string) {
	hf := historyFile(shellPath)
	if hf == "" {
		return
	}
	if err := os.MkdirAll(filepath.Dir(hf), 0o755); err != nil {
		return
	}
	f, err := os.OpenFile(hf, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return
	}
	defer f.Close()
	if filepath.Base(shellPath) == "fish" {
		fmt.Fprintf(f, "- cmd: %s\n  when: %d\n", cmdStr, time.Now().Unix())
	} else {
		fmt.Fprintln(f, cmdStr)
	}
}

func (execRunner) Run(ctx context.Context, cmdStr string) error {
	sh := os.Getenv("SHELL")
	if sh == "" {
		sh = "/bin/sh"
	}
	appendHistory(sh, cmdStr)
	c := commandContext(ctx, sh, "-c", cmdStr)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
