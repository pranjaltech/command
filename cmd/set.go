/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"

	"command/internal/config"
	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set <field> <value>",
	Short: "Update a configuration field",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		switch args[0] {
		case "api_key":
			cfg.APIKey = args[1]
		case "model":
			cfg.Model = args[1]
		case "temperature":
			f, err := strconv.ParseFloat(args[1], 32)
			if err != nil {
				return err
			}
			cfg.Temperature = float32(f)
		default:
			return fmt.Errorf("unknown field %q", args[0])
		}
		return config.Save(cfg)
	},
}

func init() {
	configCmd.AddCommand(setCmd)
}
