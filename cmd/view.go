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
		key := "(none)"
		if cfg.APIKey != "" {
			if len(cfg.APIKey) <= 4 {
				key = strings.Repeat("*", len(cfg.APIKey))
			} else {
				key = strings.Repeat("*", len(cfg.APIKey)-4) + cfg.APIKey[len(cfg.APIKey)-4:]
			}
		}
		fmt.Printf("model: %s\n", cfg.Model)
		fmt.Printf("temperature: %.2f\n", cfg.Temperature)
		fmt.Printf("api_key: %s\n", key)
		return nil
	},
}

func init() {
	configCmd.AddCommand(viewCmd)
}
