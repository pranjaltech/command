# cmd

A command-line helper that turns English instructions into shell commands. It gathers some environment details to improve accuracy and lets you review the suggested commands before running them.

## For Users

### Install
- **Homebrew**: `brew install --cask pranjaltech/tools/cmd`
- **Manual build**: run `scripts/install.sh` and it will place the binary under `/usr/local/bin`.

The tool works on macOS and Linux with Go 1.22+ installed.

### What can it do?
- Translate natural language into runnable shell commands.
- Let you pick and edit the generated command.
- Remember your configuration in `~/.config/cmd/config.yaml`.

Run `cmd "your prompt"` to start.

The first run guides you through choosing a model provider, entering an API key and deciding whether to enable anonymous telemetry. Supported providers are OpenAI, Anthropic, Gemini, OpenRouter and Ollama. You can also set credentials via environment variables such as `OPENAI_API_KEY` or `ANTHROPIC_API_KEY`.

Default models used for each provider (latest light models as of 0.2.0):

| Provider   | Default model       | Notes                                  |
|-----------|---------------------|----------------------------------------|
| OpenAI     | gpt-4.1-nano        | Fastest & cheapest with 1M context     |
| Anthropic  | claude-haiku-4-5    | 4-5x faster than Sonnet 4.5            |
| Gemini     | gemini-3-flash      | Pro-level intelligence at Flash speed  |
| OpenRouter | x-ai/grok-4.1-fast  | Best agentic tool calling with 2M ctx  |
| Ollama     | qwen3:8b            | Best local lightweight model           |

After onboarding, you can change these via `cmd config`.

If telemetry is enabled, set `LANGFUSE_HOST`, `LANGFUSE_PUBLIC_KEY` and `LANGFUSE_SECRET_KEY` so events are sent to your Langfuse instance.

## For Developers

Clone the repo and install the helper tools:

```bash
./.codex/setup.sh
```

### Building

Build the binary with version information:

```bash
make build        # Build with version info (commit hash, build date)
make install      # Install to $GOPATH/bin
make test         # Run tests
make clean        # Remove build artifacts
```

The Makefile automatically injects version info via Go's ldflags. To build manually:

```bash
go build -ldflags "-X command/cmd.Version=0.2.0 -X command/cmd.Commit=$(git rev-parse --short HEAD) -X command/cmd.Date=$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o command .
```

Check the version:

```bash
./command --version
# Output:
# cmd version 0.2.0
# commit: 065ffc7
# built at: 2026-01-26T16:22:54Z
# go: go1.25.6
```

### Testing

Run the verification suite before sending patches:

```bash
make test               # Run all tests
make test-coverage      # Generate coverage report
golangci-lint run ./...
staticcheck ./...
go vet ./...
```

### Contributing
- Follow Go conventions (`go fmt` etc.).
- Keep tests passing and add new ones for your changes.
- Document behaviour in this README when it affects users.
- Manage API keys and Langfuse credentials via environment variables; never commit secrets.

