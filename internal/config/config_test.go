package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSave(t *testing.T) {
	dir := t.TempDir()
	os.Setenv("CMD_CONFIG", filepath.Join(dir, "config.yaml"))
	defer os.Unsetenv("CMD_CONFIG")

	c := &Config{APIKey: "secret", Model: "gpt-4", Temperature: 0.5}
	if err := Save(c); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.APIKey != "secret" {
		t.Errorf("expected %q, got %q", "secret", got.APIKey)
	}
	if got.Model != "gpt-4" || got.Temperature != 0.5 {
		t.Errorf("unexpected model/temperature: %#v", got)
	}
}
