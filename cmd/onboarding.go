package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"command/internal/config"
)

var providerOptions = []struct {
	Name string
	Key  string
	URL  string
}{
	{"OpenAI", "openai", "https://api.openai.com/v1"},
	{"Anthropic", "anthropic", "https://api.anthropic.com"},
	{"Gemini (Google)", "gemini", "https://generativelanguage.googleapis.com/v1beta"},
	{"OpenRouter", "openrouter", "https://openrouter.ai/api/v1"},
	{"Ollama (local)", "ollama", "http://localhost:11434"},
}

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

	fmt.Printf("Selected Provider: %s\n", sel.Name)
	fmt.Print("Please provide an API key: ")
	key, _ := reader.ReadString('\n')
	key = strings.TrimSpace(key)
	if key == "" {
		return fmt.Errorf("api key is required")
	}

	fmt.Printf("Confirm the API URL: %s\n> ", sel.URL)
	urlInput, _ := reader.ReadString('\n')
	urlInput = strings.TrimSpace(urlInput)
	if urlInput == "" {
		urlInput = sel.URL
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
