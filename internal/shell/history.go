package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	file := fishHistoryPath()
	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	line := fmt.Sprintf("- cmd: %s\n  when: %d\n", cmd, time.Now().Unix())
	return appendLine(file, line)
}

func fishHistoryPath() string {
	if f := os.Getenv("fish_history"); f != "" {
		if strings.HasSuffix(f, "_history") {
			return f
		}
		return filepath.Join(dataHome(), "fish", f+"_history")
	}
	return filepath.Join(dataHome(), "fish", "fish_history")
}

func dataHome() string {
	if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
		return xdg
	}
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".local", "share")
	}
	return "/tmp"
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
