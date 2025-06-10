/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"command/internal/config"
	"command/internal/llm"
	"command/internal/probe"
	"command/internal/shell"
	"command/internal/ui"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// rootCmd represents the base command when called without any subcommands
type envCollector interface{ Collect() (probe.EnvInfo, error) }
type selector interface {
	Select([]string) (string, error)
}
type runner interface {
	Run(ctx context.Context, cmd string) error
}

func NewRootCmd(client llm.Client, collector envCollector, sel selector, run runner) *cobra.Command {
	return &cobra.Command{
		Use:          "cmd",
		Short:        "Convert natural language into shell commands",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if client == nil {
				return errors.New("api key not configured")
			}
			phrase := strings.Join(args, " ")
			env, err := collector.Collect()
			if err != nil {
				return err
			}
			cmds, err := client.GenerateCommands(cmd.Context(), phrase, env)
			if err != nil {
				return err
			}
			choice, err := sel.Select(cmds)
			if err != nil {
				return err
			}
			return run.Run(cmd.Context(), choice)
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
	// Load environment variables from .env when present. Ignore errors so
	// the CLI still works without the file.
	_ = godotenv.Load()

	apiKey := os.Getenv("OPENAI_API_KEY")
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: %v\n", err)
	}
	if apiKey == "" {
		apiKey = cfg.APIKey
	}
	if apiKey == "" {
		if term.IsTerminal(int(os.Stdin.Fd())) {
			fmt.Fprint(os.Stderr, "Enter OpenAI API key: ")
			reader := bufio.NewReader(os.Stdin)
			key, _ := reader.ReadString('\n')
			apiKey = strings.TrimSpace(key)
			cfg.APIKey = apiKey
			if err := config.Save(cfg); err != nil {
				fmt.Fprintf(os.Stderr, "warning: %v\n", err)
			}
		}
	}
	client, err := llm.NewOpenAIClient(apiKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: %v\n", err)
	}
	rootCmd = NewRootCmd(client, probe.NewProbe(), ui.NewSelector(), shell.NewRunner())
}
