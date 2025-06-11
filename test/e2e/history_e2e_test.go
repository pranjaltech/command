package e2e

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"command/cmd"
	"command/internal/llm"
	"command/internal/probe"
	"command/internal/shell"
)

type fakeClient struct{}

func (fakeClient) GenerateCommands(ctx context.Context, prompt string, env probe.EnvInfo) ([]string, error) {
	return []string{prompt}, nil
}

type fakeCollector struct{}

func (fakeCollector) Collect() (probe.EnvInfo, error) { return probe.EnvInfo{}, nil }

type fakeSelector struct{}

func (fakeSelector) Select(opts []string) (string, error) { return opts[0], nil }

func TestHistoryUpdates(t *testing.T) {
	shells := []string{"bash", "zsh", "fish"}
	for _, sh := range shells {
		t.Run(sh, func(t *testing.T) {
			tmp := t.TempDir()
			shellPath := filepath.Join(tmp, sh)
			script := "#!/usr/bin/bash\nexec /usr/bin/bash \"$@\""
			if err := os.WriteFile(shellPath, []byte(script), 0o755); err != nil {
				t.Fatalf("write shell: %v", err)
			}
			var histFile string
			if sh == "fish" {
				os.Setenv("HOME", tmp)
				histFile = filepath.Join(tmp, ".local/share/fish/fish_history")
			} else {
				os.Setenv("HISTFILE", filepath.Join(tmp, sh+"_history"))
				os.Setenv("HOME", tmp)
				histFile = os.Getenv("HISTFILE")
			}
			os.Setenv("SHELL", shellPath)
			defer func() {
				os.Unsetenv("SHELL")
				os.Unsetenv("HISTFILE")
				os.Unsetenv("HOME")
			}()
			var fc llm.Client = fakeClient{}
			root := cmd.NewRootCmd(&fc, fakeCollector{}, fakeSelector{}, shell.NewRunner())
			root.SetArgs([]string{"echo hi"})
			if err := root.Execute(); err != nil {
				t.Fatalf("execute: %v", err)
			}
			data, err := os.ReadFile(histFile)
			if err != nil {
				t.Fatalf("read history: %v", err)
			}
			if !strings.Contains(string(data), "echo hi") {
				t.Errorf("history missing command: %q", data)
			}
		})
	}
}
