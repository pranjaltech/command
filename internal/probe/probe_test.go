package probe

import (
	"errors"
	"strings"
	"testing"
)

type stubRunner struct{ outputs map[string]string }

func (s stubRunner) CombinedOutput(name string, arg ...string) ([]byte, error) {
	key := name + " " + strings.Join(arg, " ")
	out, ok := s.outputs[key]
	if !ok {
		return nil, errors.New("unexpected command")
	}
	return []byte(out), nil
}

func TestProbeCollect(t *testing.T) {
	r := stubRunner{outputs: map[string]string{
		"uname -r":                        "5.0\n",
		"git rev-parse --show-toplevel":   "/repo\n",
		"git rev-parse --abbrev-ref HEAD": "main\n",
		"git status --porcelain":          " M file.go\n",
	}}
	p := &Probe{run: r}
	info, err := p.Collect()
	if err != nil {
		t.Fatalf("collect: %v", err)
	}
	if info.Kernel != "5.0" {
		t.Errorf("kernel: want 5.0 got %q", info.Kernel)
	}
	if info.GitRoot != "/repo" {
		t.Errorf("root: want /repo got %q", info.GitRoot)
	}
	if info.GitBranch != "main" {
		t.Errorf("branch: want main got %q", info.GitBranch)
	}
	if !info.GitDirty {
		t.Errorf("dirty: expected true")
	}
}
