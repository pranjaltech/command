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

If telemetry is enabled, set `LANGFUSE_HOST`, `LANGFUSE_PUBLIC_KEY` and `LANGFUSE_SECRET_KEY` so events are sent to your Langfuse instance.

## For Developers

Clone the repo and install the helper tools:

```bash
./.codex/setup.sh
```

Run the verification suite before sending patches:

```bash
golangci-lint run ./...
staticcheck ./...
Golines -m 120 -w $(git ls-files '*.go')
go vet ./...
go test -race -coverprofile=coverage.out ./...
```

### Contributing
- Follow Go conventions (`go fmt` etc.).
- Keep tests passing and add new ones for your changes.
- Document behaviour in this README when it affects users.
- Manage API keys and Langfuse credentials via environment variables; never commit secrets.

