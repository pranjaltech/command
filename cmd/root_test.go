package cmd

import (
	"bytes"
	"context"
	"testing"

	"command/internal/llm"
	"command/internal/probe"
)

type stubLLM struct{ out []string }

func (s stubLLM) GenerateCommands(ctx context.Context, prompt string, env probe.EnvInfo) ([]string, error) {
	return s.out, nil
}

var _ llm.Client = stubLLM{}

type stubSelector struct{ pick string }

func (s stubSelector) Select(opts []string) (string, error) { return s.pick, nil }

type stubRunner struct{ cmd string }

func (r *stubRunner) Run(ctx context.Context, cmd string) error { r.cmd = cmd; return nil }

type stubProbe struct{}

func (stubProbe) Collect() (probe.EnvInfo, error) { return probe.EnvInfo{}, nil }

func TestRootCmd(t *testing.T) {
	r := &stubRunner{}
	cmd := NewRootCmd(stubLLM{out: []string{"ls"}}, stubProbe{}, stubSelector{pick: "ls"}, r)
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
	cmd := NewRootCmd(nil, stubProbe{}, stubSelector{pick: ""}, &stubRunner{})
	cmd.SetArgs([]string{"noop"})
	if err := cmd.Execute(); err == nil {
		t.Fatalf("expected error when api key missing")
	}
}
