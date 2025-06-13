/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"command/internal/config"
	"github.com/spf13/cobra"
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		p := cfg.Providers[cfg.Provider]
		key := "(none)"
		if p.APIKey != "" {
			if len(p.APIKey) <= 4 {
				key = strings.Repeat("*", len(p.APIKey))
			} else {
				key = strings.Repeat("*", len(p.APIKey)-4) + p.APIKey[len(p.APIKey)-4:]
			}
		}
		fmt.Printf("provider: %s\n", cfg.Provider)
		fmt.Printf("api_url: %s\n", p.APIURL)
		fmt.Printf("model: %s\n", cfg.Model)
		fmt.Printf("temperature: %.2f\n", cfg.Temperature)
		fmt.Printf("api_key: %s\n", key)
		if cfg.TelemetryDisable {
			fmt.Println("telemetry: disabled")
		} else {
			fmt.Println("telemetry: enabled")
		}
		return nil
	},
}

func init() {
	configCmd.AddCommand(viewCmd)
}
