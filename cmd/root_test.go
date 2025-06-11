package cmd

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	"command/internal/llm"
	"command/internal/probe"
)

type stubLLM struct{ out []string }

func (s stubLLM) GenerateCommands(ctx context.Context, prompt string, env probe.EnvInfo) ([]string, error) {
	return s.out, nil
}

var _ llm.Client = stubLLM{}

type clarifyLLM struct{ calls int }

func (c *clarifyLLM) GenerateCommands(ctx context.Context, prompt string, env probe.EnvInfo) ([]string, error) {
	c.calls++
	if c.calls == 1 {
		return nil, llm.NeedClarificationError{Question: "which dir?"}
	}
	return []string{"ls"}, nil
}

type stubSelector struct{ pick string }

func (s stubSelector) Select(opts []string) (string, error) { return s.pick, nil }

type stubRunner struct{ cmd string }

func (r *stubRunner) Run(ctx context.Context, cmd string) error { r.cmd = cmd; return nil }

type stubProbe struct{}

func (stubProbe) Collect() (probe.EnvInfo, error) { return probe.EnvInfo{}, nil }

func TestRootCmd(t *testing.T) {
	r := &stubRunner{}
	c := llm.Client(stubLLM{out: []string{"ls"}})
	cmd := NewRootCmd(&c, stubProbe{}, stubSelector{pick: "ls"}, r)
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"list"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if r.cmd != "ls" {
		t.Errorf("expected runner to execute ls, got %q", r.cmd)
	}
}

func TestRootCmd_NoClient(t *testing.T) {
	var c llm.Client
	cmd := NewRootCmd(&c, stubProbe{}, stubSelector{pick: ""}, &stubRunner{})
	cmd.SetArgs([]string{"noop"})
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "api key not configured") {
		t.Fatalf("expected error when api key missing, got %v", err)
	}
}

func TestRootCmd_NoPrompt(t *testing.T) {
	r := &stubRunner{}
	c := llm.Client(stubLLM{out: []string{"ls"}})
	cmd := NewRootCmd(&c, stubProbe{}, stubSelector{pick: "ls"}, r)
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "at least 1 arg") {
		t.Fatalf("expected arg error, got %v", err)
	}
}

func TestRootCmd_ClarificationFlow(t *testing.T) {
	r := &stubRunner{}
	clar := &clarifyLLM{}
	c := llm.Client(clar)
	cmd := NewRootCmd(&c, stubProbe{}, stubSelector{pick: "ls"}, r)
	// provide clarification input
	rPipe, wPipe, _ := os.Pipe()
	oldStdin := os.Stdin
	os.Stdin = rPipe
	_, _ = wPipe.WriteString("src\n")
	wPipe.Close()
	defer func() { os.Stdin = oldStdin }()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"list files"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if clar.calls != 2 {
		t.Errorf("expected 2 calls, got %d", clar.calls)
	}
	if r.cmd != "ls" {
		t.Errorf("expected runner to execute ls, got %q", r.cmd)
	}
}

func TestRootCmd_VersionFlag(t *testing.T) {
	c := llm.Client(stubLLM{out: []string{"ls"}})
	cmd := NewRootCmd(&c, stubProbe{}, stubSelector{pick: "ls"}, &stubRunner{})
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--version"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	out := strings.TrimSpace(buf.String())
	if !strings.Contains(out, Version) {
		t.Errorf("expected output to contain %q, got %q", Version, out)
	}
}
