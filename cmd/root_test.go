package cmd

import (
	"bytes"
	"context"
	"testing"

	"command/internal/llm"
	"command/internal/probe"
)

type stubLLM struct{ out string }

func (s stubLLM) GenerateCommand(ctx context.Context, prompt string, env probe.EnvInfo) (string, error) {
	return s.out, nil
}

var _ llm.Client = stubLLM{}

type stubProbe struct{}

func (stubProbe) Collect() (probe.EnvInfo, error) { return probe.EnvInfo{}, nil }

func TestRootCmd(t *testing.T) {
	cmd := NewRootCmd(stubLLM{out: "ls"}, stubProbe{})
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"list"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if got := buf.String(); got != "ls\n" {
		t.Errorf("want ls got %q", got)
	}
}
