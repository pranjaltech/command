// Package cmd houses Cobra commands.
package cmd

import (
	"command/internal/config"
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
		return cmd.Help()
	},
}
