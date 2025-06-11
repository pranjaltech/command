// Package cmd houses Cobra commands.
package cmd

import (
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage cmd configuration",
	Long:  "View or modify settings stored in $HOME/.config/cmd/config.yaml",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}
