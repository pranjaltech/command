/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"context"
	"encoding/json"
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

func NewRootCmd(client *llm.Client, collector envCollector, sel selector, run runner) *cobra.Command {
	return &cobra.Command{
		Use:   "cmd [flags] <prompt>",
		Short: "Convert natural language into shell commands",
		Long: "cmd translates English instructions into shell commands using OpenAI." +
			" Configuration is read from $HOME/.config/cmd/config.yaml or $CMD_CONFIG." +
			" Fields:\n  api_key - OpenAI token (encrypted)\n  model - model name" +
			" (default " + config.DefaultModel + ")\n  temperature - sampling temperature",
		Version:      Version,
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if *client == nil {
				return errors.New("api key not configured")
			}
			phrase := strings.Join(args, " ")
			if debug {
				fmt.Fprintf(os.Stderr, "prompt: %s\n", phrase)
			}
			env, err := collector.Collect()
			if err != nil {
				return err
			}
			if debug {
				data, _ := json.MarshalIndent(env, "", "  ")
				fmt.Fprintf(os.Stderr, "env: %s\n", data)
			}
			attempts := 0
			var cmds []string
			for {
				var err error
				cmds, err = (*client).GenerateCommands(cmd.Context(), phrase, env)
				if err == nil {
					break
				}
				var nc llm.NeedClarificationError
				if errors.As(err, &nc) && attempts < 2 {
					fmt.Fprintln(os.Stderr, nc.Question)
					fmt.Fprint(os.Stderr, "> ")
					reader := bufio.NewReader(os.Stdin)
					extra, _ := reader.ReadString('\n')
					phrase += "\n" + strings.TrimSpace(extra)
					attempts++
					continue
				}
				return err
			}
			if debug {
				fmt.Fprintf(os.Stderr, "suggestions: %v\n", cmds)
			}
			choice, err := sel.Select(cmds)
			if err != nil {
				return err
			}
			if debug {
				fmt.Fprintf(os.Stderr, "selected: %s\n", choice)
			}
			return run.Run(cmd.Context(), choice)
		},
	}
}

var rootCmd *cobra.Command
var (
	cfg         *config.Config
	apiKey      string
	model       string
	temperature float32
	client      llm.Client
	debug       bool
)

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

	var err error
	cfg, err = config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: %v\n", err)
	}
	apiKey = os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		apiKey = cfg.APIKey
	}
	model = cfg.Model
	temperature = cfg.Temperature

	rootCmd = NewRootCmd(&client, probe.NewProbe(), ui.NewSelector(), shell.NewRunner())
	rootCmd.AddCommand(configCmd)
	rootCmd.PersistentFlags().StringVar(&model, "model", model, "OpenAI model")
	rootCmd.PersistentFlags().Float32Var(&temperature, "temperature", temperature, "sampling temperature")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "verbose debug output")

	rootCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if apiKey == "" {
			if term.IsTerminal(int(os.Stdin.Fd())) {
				fmt.Fprint(os.Stderr, "Enter OpenAI API key: ")
				reader := bufio.NewReader(os.Stdin)
				key, _ := reader.ReadString('\n')
				apiKey = strings.TrimSpace(key)
			}
		}
		var err error
		oa, err := llm.NewOpenAIClient(apiKey, model, temperature)
		if err != nil {
			return err
		}
		if debug {
			oa.EnableDebug(os.Stderr)
		}
		client = oa
		cfg.APIKey = apiKey
		cfg.Model = model
		cfg.Temperature = temperature
		return config.Save(cfg)
	}
}
