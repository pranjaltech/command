package cmd

import (
	"bytes"
	"testing"
)

func TestRootCmd(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"list all directories"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	got := buf.String()
	want := "ls -d */\n"
	if got != want {
		t.Errorf("want %q got %q", want, got)
	}
}
