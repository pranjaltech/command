package probe

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// EnvInfo holds environment details collected for command generation.
type EnvInfo struct {
	OS           string
	Arch         string
	Kernel       string
	WorkDir      string
	GitRoot      string
	GitBranch    string
	GitDirty     bool
	Shell        string
	ShellVersion string
}

// CmdRunner executes a command and returns its combined output.
type CmdRunner interface {
	CombinedOutput(name string, arg ...string) ([]byte, error)
}

type runnerFunc func(name string, arg ...string) ([]byte, error)

func (f runnerFunc) CombinedOutput(name string, arg ...string) ([]byte, error) {
	return f(name, arg...)
}

// Probe gathers environment information.
type Probe struct {
	run CmdRunner
}

// NewProbe returns a Probe using the os/exec runner.
func NewProbe() *Probe {
	return &Probe{run: runnerFunc(defaultRun)}
}

func defaultRun(name string, arg ...string) ([]byte, error) {
	return exec.Command(name, arg...).CombinedOutput()
}

// Collect gathers information about the runtime environment.
func (p *Probe) Collect() (EnvInfo, error) {
	info := EnvInfo{
		OS:    runtime.GOOS,
		Arch:  runtime.GOARCH,
		Shell: os.Getenv("SHELL"),
	}
	if wd, err := os.Getwd(); err == nil {
		info.WorkDir = wd
	}
	if out, err := p.run.CombinedOutput("uname", "-r"); err == nil {
		info.Kernel = strings.TrimSpace(string(out))
	}
	if out, err := p.run.CombinedOutput("git", "rev-parse", "--show-toplevel"); err == nil {
		info.GitRoot = strings.TrimSpace(string(out))
	}
	if out, err := p.run.CombinedOutput("git", "rev-parse", "--abbrev-ref", "HEAD"); err == nil {
		info.GitBranch = strings.TrimSpace(string(out))
	}
	if out, err := p.run.CombinedOutput("git", "status", "--porcelain"); err == nil {
		if strings.TrimSpace(string(out)) != "" {
			info.GitDirty = true
		}
	}
	return info, nil
}
