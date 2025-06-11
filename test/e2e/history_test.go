package e2e

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"command/internal/shell"
)

func runShell(t *testing.T, shellName string, historyPathParts []string, want string) {
	t.Helper()
	path, err := exec.LookPath(shellName)
	if err != nil {
		t.Skipf("%s not installed", shellName)
	}
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("SHELL", path)

	r := shell.NewRunner()
	if err := r.Run(context.Background(), "echo hi"); err != nil {
		t.Fatalf("run: %v", err)
	}
	histFile := filepath.Join(append([]string{tmp}, historyPathParts...)...)
	data, err := os.ReadFile(histFile)
	if err != nil {
		t.Fatalf("read history: %v", err)
	}
	if !strings.Contains(string(data), want) {
		t.Fatalf("expected %q in history, got %q", want, data)
	}
}

func TestEndToEnd_History_Bash(t *testing.T) {
	runShell(t, "bash", []string{".bash_history"}, "echo hi")
}

func TestEndToEnd_History_Zsh(t *testing.T) {
	runShell(t, "zsh", []string{".zsh_history"}, "echo hi")
}

func TestEndToEnd_History_Fish(t *testing.T) {
	runShell(t, "fish", []string{".local", "share", "fish", "fish_history"}, "cmd: echo hi")
}
