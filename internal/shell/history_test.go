package shell

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAppendHistory_Bash(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "bash_history")
	os.Setenv("HISTFILE", file)
	defer os.Unsetenv("HISTFILE")

	if err := appendHistory("/bin/bash", "echo hi"); err != nil {
		t.Fatalf("appendHistory bash: %v", err)
	}
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if string(data) != "echo hi\n" {
		t.Errorf("unexpected bash history: %q", string(data))
	}
}

func TestAppendHistory_Zsh(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "zsh_history")
	os.Setenv("HISTFILE", file)
	defer os.Unsetenv("HISTFILE")

	if err := appendHistory("/bin/zsh", "ls -l"); err != nil {
		t.Fatalf("appendHistory zsh: %v", err)
	}
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if !strings.HasSuffix(string(data), ";ls -l\n") {
		t.Errorf("unexpected zsh history: %q", string(data))
	}
}

func TestAppendHistory_Fish(t *testing.T) {
	dir := t.TempDir()
	os.Setenv("HOME", dir)
	defer os.Unsetenv("HOME")

	if err := appendHistory("/usr/bin/fish", "git status"); err != nil {
		t.Fatalf("appendHistory fish: %v", err)
	}
	path := filepath.Join(dir, ".local/share/fish/fish_history")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if !strings.Contains(string(data), "git status") {
		t.Errorf("unexpected fish history: %q", string(data))
	}
}

func TestAppendHistory_Fish_Custom(t *testing.T) {
	dir := t.TempDir()
	os.Setenv("HOME", dir)
	os.Setenv("XDG_DATA_HOME", filepath.Join(dir, "xdg"))
	os.Setenv("fish_history", "alt")
	defer func() {
		os.Unsetenv("HOME")
		os.Unsetenv("XDG_DATA_HOME")
		os.Unsetenv("fish_history")
	}()

	if err := appendHistory("/usr/bin/fish", "ls"); err != nil {
		t.Fatalf("appendHistory fish custom: %v", err)
	}
	path := filepath.Join(dir, "xdg", "fish", "alt_history")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if !strings.Contains(string(data), "ls") {
		t.Errorf("unexpected fish history: %q", string(data))
	}
}
