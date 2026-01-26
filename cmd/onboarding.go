package cmd

import (
	"fmt"
	"os"

	"command/internal/config"
	"command/internal/ui"
)

type providerOption struct {
	Name   string
	Key    string
	URL    string
	KeyEnv string
	URLEnv string // For environment variable fallback (not shown in UI)
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
	fmt.Println("Welcome to cmd! Let's set up your AI provider.")
	fmt.Println()

	// Build provider list for UI
	providers := make([]ui.ProviderInfo, len(providerOptions))
	for i, p := range providerOptions {
		providers[i] = ui.ProviderInfo{Name: p.Name, Key: p.Key}
	}

	// Run the configuration UI
	result, err := ui.RunConfigUI(providers, nil)
	if err != nil {
		return err
	}

	if result.Canceled {
		fmt.Println("\nSetup canceled.")
		return nil
	}

	// Get provider details
	sel := providerMap[result.Provider]

	// Check for environment variable fallback
	apiKey := result.APIKey
	if apiKey == "" {
		envKey := os.Getenv(sel.KeyEnv)
		if envKey != "" {
			apiKey = envKey
			fmt.Printf("\nUsing API key from %s\n", sel.KeyEnv)
		}
	}

	// Ollama doesn't require an API key (local usage)
	if apiKey == "" && sel.Key != "ollama" {
		return fmt.Errorf("api key is required")
	}

	// Check for API URL environment variable
	apiURL := sel.URL
	if envURL := os.Getenv(sel.URLEnv); envURL != "" {
		apiURL = envURL
		fmt.Printf("Using API URL from %s: %s\n", sel.URLEnv, envURL)
	}

	// Save configuration
	cfg := &config.Config{
		Provider:         sel.Key,
		Providers:        map[string]config.Provider{sel.Key: {APIKey: apiKey, APIURL: apiURL}},
		Model:            config.DefaultModelForProvider(sel.Key),
		TelemetryDisable: true, // Default to disabled
	}
	if err := config.Save(cfg); err != nil {
		return err
	}

	fmt.Println("\n" + ui.Checkmark() + " Configuration saved!")
	fmt.Println("\nRun " + ui.BoldStyle.Render("cmd \"your command\"") + " to get started.")
	fmt.Println(ui.MutedStyle.Render("Tip: Use 'cmd config' to change settings anytime."))

	return nil
}
