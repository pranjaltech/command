package llm

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"command/internal/probe"
)

const PromptTemplate = `You are **cmd**, a terminal-command assistant.

## 1. Environment (immutable for this session)
%s

## 2. Dynamic context (inject when available)
### Installed binaries (truncated PATH scan)
%s

### Top-level files in $PWD
%s

### Git status (if inside a repo)
%s

### User aliases & functions
%s

## 3. Output contract  -- obey exactly
- Respond ONLY with valid JSON—no markdown, no prose.
- Schema:
  {
    "commands": ["<cmd1>", "<cmd2>", "<cmd3>"],   // ≤3, most relevant first
    "need_clarification": "<question or null>",
    "notes": "<one-sentence rationale or null>"
  }
- Never output destructive commands (e.g. 'rm -rf /') unless the user explicitly requests them and they are clearly marked by adding "dangerous": true alongside that entry.
- If the request is ambiguous, leave "commands" empty and set "need_clarification".

## 4. Style rules
- Fish syntax by default.
- Chain commands with && only when later commands depend on earlier success.
- Prefer concise flags (ls -la > ls --all --long).

## 5. Few-shot examples
U: list go files in current dir
A: {"commands":["ls *.go"],"need_clarification":null,"notes":null}

U: initialise a git repo and push first commit
A: {"commands":["git init && git add . && git commit -m \"init\" && git branch -M main && git remote add origin <url> && git push -u origin main"],"need_clarification":null,"notes":"uses main branch by default"}

Begin!`

// NeedClarificationError is returned when the model asks a follow-up question
// instead of providing commands.
type NeedClarificationError struct{ Question string }

func (e NeedClarificationError) Error() string {
	if e.Question == "" {
		return "clarification requested"
	}
	return "clarification requested: " + e.Question
}

// BuildSystemPrompt constructs the system prompt from environment info.
func BuildSystemPrompt(env probe.EnvInfo) (string, error) {
	envJSON, err := json.Marshal(env)
	if err != nil {
		return "", fmt.Errorf("marshal env: %w", err)
	}
	files, _ := os.ReadDir(env.WorkDir)
	var names []string
	for i, f := range files {
		if i >= 20 {
			break
		}
		names = append(names, f.Name())
	}
	fileList := strings.Join(names, " ")

	bins := UniqueBinaries(20)

	gitStatus := fmt.Sprintf("{\"root\":%q,\"branch\":%q,\"dirty\":%t}", env.GitRoot, env.GitBranch, env.GitDirty)

	return fmt.Sprintf(PromptTemplate, envJSON, bins, fileList, gitStatus, "[]"), nil
}

var (
	binsOnce sync.Once
	binsList []string
)

// UniqueBinaries returns a space-separated list of unique binaries from PATH.
func UniqueBinaries(limit int) string {
	binsOnce.Do(func() {
		seen := make(map[string]struct{})
		for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
			entries, err := os.ReadDir(dir)
			if err != nil {
				continue
			}
			for _, e := range entries {
				if len(binsList) >= limit {
					return
				}
				name := e.Name()
				if _, ok := seen[name]; ok {
					continue
				}
				seen[name] = struct{}{}
				binsList = append(binsList, name)
			}
		}
	})
	if len(binsList) > limit {
		return strings.Join(binsList[:limit], " ")
	}
	return strings.Join(binsList, " ")
}

// LLMResponse is the expected JSON structure from all providers.
type LLMResponse struct {
	Commands          []json.RawMessage `json:"commands"`
	NeedClarification json.RawMessage   `json:"need_clarification"`
	Notes             json.RawMessage   `json:"notes"`
}

// ParseCommands extracts command strings from the LLM response content.
func ParseCommands(content string) ([]string, error) {
	var raw LLMResponse
	if err := json.Unmarshal([]byte(strings.TrimSpace(content)), &raw); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	if len(raw.NeedClarification) > 0 && string(raw.NeedClarification) != "null" {
		var q string
		if err := json.Unmarshal(raw.NeedClarification, &q); err == nil && q != "" {
			return nil, NeedClarificationError{Question: q}
		}
	}

	var out []string
	for _, r := range raw.Commands {
		var s string
		if err := json.Unmarshal(r, &s); err == nil {
			out = append(out, s)
			continue
		}
		var obj map[string]interface{}
		if err := json.Unmarshal(r, &obj); err == nil {
			if v, ok := obj["command"].(string); ok {
				out = append(out, v)
				continue
			}
			if v, ok := obj["cmd"].(string); ok {
				out = append(out, v)
				continue
			}
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no commands parsed")
	}
	if len(out) > 3 {
		out = out[:3]
	}
	return out, nil
}
