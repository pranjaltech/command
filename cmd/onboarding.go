package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"command/internal/config"
)

type providerOption struct {
	Name   string
	Key    string
	URL    string
	KeyEnv string
	URLEnv string
}

var providerOptions = []providerOption{
	{"OpenAI", "openai", "https://api.openai.com/v1", "OPENAI_API_KEY", "OPENAI_BASE_URL"},
	{"Anthropic", "anthropic", "https://api.anthropic.com", "ANTHROPIC_API_KEY", "ANTHROPIC_API_URL"},
	{
		"Gemini (Google)",
		"gemini",
		"https://generativelanguage.googleapis.com/v1beta",
		"GEMINI_API_KEY",
		"GEMINI_API_URL",
	},
	{"OpenRouter", "openrouter", "https://openrouter.ai/api/v1", "OPENROUTER_API_KEY", "OPENROUTER_API_URL"},
	{"Ollama (local)", "ollama", "http://localhost:11434", "OLLAMA_API_KEY", "OLLAMA_API_URL"},
}

var providerMap = func() map[string]providerOption {
	m := make(map[string]providerOption)
	for _, p := range providerOptions {
		m[p.Key] = p
	}
	return m
}()

func runOnboarding() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("This tool requires an AI model to work. Please select your AI model provider:")
	for i, p := range providerOptions {
		fmt.Printf("%d. %s\n", i+1, p.Name)
	}
	fmt.Print("> ")
	choiceStr, _ := reader.ReadString('\n')
	choiceStr = strings.TrimSpace(choiceStr)
	idx, _ := strconv.Atoi(choiceStr)
	if idx < 1 || idx > len(providerOptions) {
		return fmt.Errorf("invalid choice")
	}
	sel := providerOptions[idx-1]

	envKey := os.Getenv(sel.KeyEnv)
	if envKey != "" {
		fmt.Printf("Using %s from %s\n", sel.Name, sel.KeyEnv)
	}

	fmt.Printf("Selected Provider: %s\n", sel.Name)
	fmt.Print("Please provide an API key: ")
	key, _ := reader.ReadString('\n')
	key = strings.TrimSpace(key)
	if key == "" {
		key = envKey
	}
	if key == "" {
		return fmt.Errorf("api key is required")
	}

	envURL := os.Getenv(sel.URLEnv)
	if envURL != "" {
		fmt.Printf("Using %s for API URL from %s\n", envURL, sel.URLEnv)
	}
	fmt.Printf("Confirm the API URL: %s\n> ", sel.URL)
	urlInput, _ := reader.ReadString('\n')
	urlInput = strings.TrimSpace(urlInput)
	if urlInput == "" {
		if envURL != "" {
			urlInput = envURL
		} else {
			urlInput = sel.URL
		}
	}

	fmt.Print("Can we collect some anonymous telemetry to improve this tool? [y/N]: ")
	teleStr, _ := reader.ReadString('\n')
	tele := strings.TrimSpace(strings.ToLower(teleStr))
	enableTelemetry := tele == "y" || tele == "yes"

	cfg := &config.Config{
		Provider:         sel.Key,
		Providers:        map[string]config.Provider{sel.Key: {APIKey: key, APIURL: urlInput}},
		Model:            config.DefaultModel,
		Temperature:      config.DefaultTemperature,
		TelemetryDisable: !enableTelemetry,
	}
	if err := config.Save(cfg); err != nil {
		return err
	}
	fmt.Println("cmd is ready! Type `cmd \"the command you want\"` to start.")
	return nil
}
