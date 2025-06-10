package e2e

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestEndToEnd_Run uses the real OpenAI API to ensure the CLI works end to end.
func TestEndToEnd_Run(t *testing.T) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		t.Skip("OPENAI_API_KEY not set; skipping e2e test")
	}

	cmd := exec.Command("go", "run", "./main.go", "list all directories")
	cmd.Env = append(os.Environ(), "OPENAI_API_KEY="+apiKey)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("running cmd: %v\n%s", err, out)
	}
	if len(strings.TrimSpace(string(out))) == 0 {
		t.Fatalf("expected output, got empty")
	}
}
