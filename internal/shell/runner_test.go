package shell

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestExecRunner_History_Bash(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("SHELL", "/bin/bash")
	old := commandContext
	defer func() { commandContext = old }()
	var got []string
	commandContext = func(ctx context.Context, name string, args ...string) *exec.Cmd {
		got = append([]string{name}, args...)
		return exec.CommandContext(ctx, "true")
	}

	r := NewRunner()
	if err := r.Run(context.Background(), "echo hi"); err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(got) == 0 || !strings.Contains(got[0], "bash") {
		t.Fatalf("expected bash exec, got %v", got)
	}
	data, err := os.ReadFile(filepath.Join(tmp, ".bash_history"))
	if err != nil {
		t.Fatalf("read history: %v", err)
	}
	if string(data) != "echo hi\n" {
		t.Errorf("history content: %q", data)
	}
}

func TestExecRunner_History_Zsh(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("SHELL", "/bin/zsh")
	old := commandContext
	defer func() { commandContext = old }()
	var got []string
	commandContext = func(ctx context.Context, name string, args ...string) *exec.Cmd {
		got = append([]string{name}, args...)
		return exec.CommandContext(ctx, "true")
	}
	r := NewRunner()
	if err := r.Run(context.Background(), "echo hi"); err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(got) == 0 || !strings.Contains(got[0], "zsh") {
		t.Fatalf("expected zsh exec, got %v", got)
	}
	data, err := os.ReadFile(filepath.Join(tmp, ".zsh_history"))
	if err != nil {
		t.Fatalf("read history: %v", err)
	}
	if string(data) != "echo hi\n" {
		t.Errorf("history content: %q", data)
	}
}

func TestExecRunner_History_Fish(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("SHELL", "/usr/bin/fish")
	old := commandContext
	defer func() { commandContext = old }()
	var got []string
	commandContext = func(ctx context.Context, name string, args ...string) *exec.Cmd {
		got = append([]string{name}, args...)
		return exec.CommandContext(ctx, "true")
	}
	r := NewRunner()
	if err := r.Run(context.Background(), "echo hi"); err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(got) == 0 || !strings.Contains(got[0], "fish") {
		t.Fatalf("expected fish exec, got %v", got)
	}
	data, err := os.ReadFile(filepath.Join(tmp, ".local", "share", "fish", "fish_history"))
	if err != nil {
		t.Fatalf("read history: %v", err)
	}
	if !strings.Contains(string(data), "cmd: echo hi") {
		t.Errorf("history content: %q", data)
	}
}
