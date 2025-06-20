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
	"io"
	"os"
	"strings"

	"command/internal/config"
	"command/internal/llm"
	"command/internal/log"
	"command/internal/probe"
	"command/internal/shell"
	"command/internal/telemetry"
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
		Args:         cobra.ArbitraryArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}
			if *client == nil {
				return errors.New("api key not configured")
			}
			phrase := strings.Join(args, " ")
			log.Debugf("prompt: %s", phrase)
			env, err := collector.Collect()
			if err != nil {
				return err
			}
			data, _ := json.MarshalIndent(env, "", "  ")
			log.Debugf("env: %s", data)
			attempts := 0
			var cmds []string
			for {
				var err error
				l := ui.NewLoader()
				l.Start()
				cmds, err = (*client).GenerateCommands(cmd.Context(), phrase, env)
				l.Stop()
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
			log.Debugf("suggestions: %v", cmds)
			if track != nil {
				track.Generation(phrase, model, cmds)
			}
			choice, err := sel.Select(cmds)
			if err != nil {
				return err
			}
			log.Debugf("selected: %s", choice)
			return run.Run(cmd.Context(), choice)
		},
	}
}

var rootCmd *cobra.Command
var (
	cfg         *config.Config
	model       string
	temperature float32
	client      llm.Client
	track       telemetry.Tracker
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
	model = cfg.Model
	temperature = cfg.Temperature

	rootCmd = NewRootCmd(&client, probe.NewProbe(), ui.NewSelector(), shell.NewRunner())
	rootCmd.AddCommand(configCmd)
	rootCmd.PersistentFlags().StringVar(&model, "model", model, "OpenAI model")
	rootCmd.PersistentFlags().Float32Var(&temperature, "temperature", temperature, "sampling temperature")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "verbose debug output")

	rootCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if cfg.Provider == "" {
			if err := runOnboarding(); err != nil {
				return err
			}
			var err error
			cfg, err = config.Load()
			if err != nil {
				return err
			}
		}

		p, ok := cfg.Providers[cfg.Provider]
		if !ok {
			p = config.Provider{}
		}
		opt, okOpt := providerMap[cfg.Provider]
		if !okOpt {
			return fmt.Errorf("provider %s not supported", cfg.Provider)
		}
		if p.APIKey == "" {
			if v := os.Getenv(opt.KeyEnv); v != "" {
				p.APIKey = v
			}
		}
		if p.APIURL == "" {
			if v := os.Getenv(opt.URLEnv); v != "" {
				p.APIURL = v
			} else if opt.URL != "" {
				p.APIURL = opt.URL
			}
		}
		if p.APIKey == "" && term.IsTerminal(int(os.Stdin.Fd())) {
			fmt.Fprintf(os.Stderr, "Enter %s API key: ", cfg.Provider)
			reader := bufio.NewReader(os.Stdin)
			key, _ := reader.ReadString('\n')
			p.APIKey = strings.TrimSpace(key)
		}
		cfg.Providers[cfg.Provider] = p

		if debug {
			log.Enable(os.Stderr)
		}

		var err error
		client, err = llm.NewClient(cfg.Provider, p.APIKey, p.APIURL, model, temperature)
		if err != nil {
			return err
		}
		if debug {
			if dbg, ok := client.(interface{ EnableDebug(io.Writer) }); ok {
				dbg.EnableDebug(os.Stderr)
			}
		}

		if cfg.TelemetryDisable {
			track = telemetry.Disabled()
		} else {
			track = telemetry.NewFromEnv(cmd.Context(), debug)
		}

		cfg.Model = model
		cfg.Temperature = temperature
		return config.Save(cfg)
	}
}
