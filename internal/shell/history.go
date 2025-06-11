package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func appendHistory(shellPath, cmd string) error {
	shell := filepath.Base(shellPath)
	switch shell {
	case "bash":
		return appendBashHistory(cmd)
	case "zsh":
		return appendZshHistory(cmd)
	case "fish":
		return appendFishHistory(cmd)
	default:
		return nil
	}
}

func historyFile(defaultName string) string {
	if f := os.Getenv("HISTFILE"); f != "" {
		return f
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join("/tmp", defaultName)
	}
	return filepath.Join(home, defaultName)
}

func appendBashHistory(cmd string) error {
	file := historyFile(".bash_history")
	return appendLine(file, cmd+"\n")
}

func appendZshHistory(cmd string) error {
	file := historyFile(".zsh_history")
	line := fmt.Sprintf(": %d:0;%s\n", time.Now().Unix(), cmd)
	return appendLine(file, line)
}

func appendFishHistory(cmd string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	dir := filepath.Join(home, ".local/share/fish")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	file := filepath.Join(dir, "fish_history")
	line := fmt.Sprintf("- cmd: %s\n  when: %d\n", cmd, time.Now().Unix())
	return appendLine(file, line)
}

func appendLine(path, line string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(line)
	return err
}
