/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"strings"

	"command/internal/llm"
	"command/internal/probe"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
type envCollector interface {
	Collect() (probe.EnvInfo, error)
}

func NewRootCmd(client llm.Client, collector envCollector) *cobra.Command {
	return &cobra.Command{
		Use:   "cmd",
		Short: "Convert natural language into shell commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			phrase := strings.Join(args, " ")
			env, err := collector.Collect()
			if err != nil {
				return err
			}
			out, err := client.GenerateCommand(cmd.Context(), phrase, env)
			if err != nil {
				return err
			}
			cmd.Println(out)
			return nil
		},
	}
}

var rootCmd *cobra.Command

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	rootCmd = NewRootCmd(llm.NewOpenAIClient(apiKey), probe.NewProbe())
}
