// Package cmd houses Cobra commands.
package cmd

import (
	"fmt"
	"os"
	"strings"

	"command/internal/config"
	"command/internal/ui"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage cmd configuration",
	Long:  "View or modify settings stored in $HOME/.config/cmd/config.yaml",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if cfg.Provider == "" {
			return runOnboarding()
		}
		return runConfigMenu(cfg)
	},
}

func runConfigMenu(cfg *config.Config) error {
	for {
		// Show current configuration
		showConfigSummary(cfg)

		// Show menu options
		telemetryStatus := "enabled"
		if cfg.TelemetryDisable {
			telemetryStatus = "disabled"
		}

		options := []string{
			"Change provider",
			"Update API key",
			"Change model",
			fmt.Sprintf("Toggle telemetry (currently: %s)", telemetryStatus),
			"Exit",
		}

		fmt.Println("\nWhat would you like to do?")
		picker := ui.NewPicker()
		idx, err := picker.Pick(options)
		if err != nil {
			return err
		}

		switch idx {
		case 0: // Change provider
			if err := handleChangeProvider(cfg); err != nil {
				return err
			}
		case 1: // Update API key
			if err := handleUpdateAPIKey(cfg); err != nil {
				return err
			}
		case 2: // Change model
			if err := handleChangeModel(cfg); err != nil {
				return err
			}
		case 3: // Toggle telemetry
			cfg.TelemetryDisable = !cfg.TelemetryDisable
			if err := config.Save(cfg); err != nil {
				return err
			}
			if cfg.TelemetryDisable {
				fmt.Println("\n" + ui.Checkmark() + " Telemetry disabled")
			} else {
				fmt.Println("\n" + ui.Checkmark() + " Telemetry enabled")
			}
		case 4: // Exit
			return nil
		}
	}
}

func showConfigSummary(cfg *config.Config) {
	p := cfg.Providers[cfg.Provider]

	// Get provider display name
	providerName := cfg.Provider
	if opt, ok := providerMap[cfg.Provider]; ok {
		providerName = opt.Name
	}

	var content strings.Builder
	content.WriteString(fmt.Sprintf("Provider: %s\n", ui.BoldStyle.Render(providerName)))
	content.WriteString(fmt.Sprintf("API Key:  %s %s\n", ui.MaskAPIKey(p.APIKey), ui.Checkmark()))
	content.WriteString(fmt.Sprintf("Model:    %s", cfg.Model))

	fmt.Println()
	fmt.Println(ui.BoxStyle.Render(content.String()))
}

func handleChangeProvider(cfg *config.Config) error {
	// Build provider list for UI
	providers := make([]ui.ProviderInfo, len(providerOptions))
	for i, p := range providerOptions {
		providers[i] = ui.ProviderInfo{Name: p.Name, Key: p.Key}
	}

	// Get current config for display
	p := cfg.Providers[cfg.Provider]
	existing := &ui.ExistingConfig{
		Provider: cfg.Provider,
		APIKey:   p.APIKey,
		Model:    cfg.Model,
	}

	fmt.Println()
	result, err := ui.RunConfigUI(providers, existing)
	if err != nil {
		return err
	}

	if result.Canceled {
		fmt.Println("\nCanceled.")
		return nil
	}

	// Get provider details
	sel := providerMap[result.Provider]

	// Update provider
	cfg.Provider = sel.Key

	// Update or create provider credentials
	if cfg.Providers == nil {
		cfg.Providers = make(map[string]config.Provider)
	}

	// Determine API URL: environment variable > existing config > default
	existingProvider := cfg.Providers[sel.Key]
	apiURL := sel.URL
	if envURL := os.Getenv(sel.URLEnv); envURL != "" {
		apiURL = envURL
	} else if existingProvider.APIURL != "" {
		apiURL = existingProvider.APIURL
	}

	cfg.Providers[sel.Key] = config.Provider{
		APIKey: result.APIKey,
		APIURL: apiURL,
	}

	// Update model to provider default
	cfg.Model = config.DefaultModelForProvider(sel.Key)

	if err := config.Save(cfg); err != nil {
		return err
	}

	fmt.Println("\n" + ui.Checkmark() + " Provider changed to " + ui.BoldStyle.Render(sel.Name))
	return nil
}

func handleUpdateAPIKey(cfg *config.Config) error {
	// Build provider list for UI (just current provider)
	providers := []ui.ProviderInfo{
		{Name: providerMap[cfg.Provider].Name, Key: cfg.Provider},
	}

	// Get current config for display
	p := cfg.Providers[cfg.Provider]
	existing := &ui.ExistingConfig{
		Provider: cfg.Provider,
		APIKey:   p.APIKey,
		Model:    cfg.Model,
	}

	// Create a simple key update flow
	fmt.Println()
	result, err := ui.RunConfigUI(providers, existing)
	if err != nil {
		return err
	}

	if result.Canceled {
		fmt.Println("\nCanceled.")
		return nil
	}

	// Update API key
	provider := cfg.Providers[cfg.Provider]
	provider.APIKey = result.APIKey
	cfg.Providers[cfg.Provider] = provider

	if err := config.Save(cfg); err != nil {
		return err
	}

	fmt.Println("\n" + ui.Checkmark() + " API key updated")
	return nil
}

func handleChangeModel(cfg *config.Config) error {
	// Show current model and provider-specific default
	defaultModel := config.DefaultModelForProvider(cfg.Provider)

	options := []string{
		fmt.Sprintf("Use default (%s)", defaultModel),
		"Enter custom model",
	}

	fmt.Println("\nSelect model option:")
	picker := ui.NewPicker()
	idx, err := picker.Pick(options)
	if err != nil {
		return err
	}

	var newModel string
	switch idx {
	case 0: // Use default
		newModel = defaultModel
	case 1: // Custom
		selector := ui.NewSelector()
		models := []string{cfg.Model} // Start with current model
		fmt.Println("\nEnter model name (edit or replace):")
		newModel, err = selector.Select(models)
		if err != nil {
			return err
		}
		newModel = strings.TrimSpace(newModel)
		if newModel == "" {
			fmt.Println("\nCanceled.")
			return nil
		}
	}

	cfg.Model = newModel
	if err := config.Save(cfg); err != nil {
		return err
	}

	fmt.Println("\n" + ui.Checkmark() + " Model changed to " + ui.BoldStyle.Render(newModel))
	return nil
}
