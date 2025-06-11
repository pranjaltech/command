# cmd

`cmd` converts natural language instructions into shell commands using OpenAI. It collects details about your environment to craft accurate suggestions and lets you choose the command to execute.

## Installation

Use the provided script to build and install the binary to `/usr/local/bin` (override `PREFIX` to change the target directory). If the directory is not writable, run the script with `sudo` or set `PREFIX` to a path you own:

```bash
scripts/install.sh
```

To uninstall:

```bash
scripts/uninstall.sh
```

## Usage

1. Run `cmd <prompt>` and enter your OpenAI API key when prompted. It is stored
   encrypted under `$HOME/.config/cmd/config.yaml` for future use.
2. Subsequent runs reuse the saved key automatically.

```bash
cmd list all directories
```

A list of up to three commands is shown. Use the arrow keys to pick one and press `Enter` to run it. The first command is selected by default.

### Configuration

Settings are stored in `$HOME/.config/cmd/config.yaml` (override with `CMD_CONFIG`).
Fields:

- `api_key` – your OpenAI API token (encrypted)
- `model` – model name to use (default `gpt-4o-mini`)
- `temperature` – sampling temperature for the LLM

You can override `model` and `temperature` with `--model` and `--temperature` flags. The values are persisted back to the config file.

## Development

This project requires Go 1.22+. After cloning, install helper tools with `.codex/setup.sh` and run the verification suite:

```bash
./.codex/setup.sh
golangci-lint run ./...
staticcheck ./...
golines -m 120 -w $(git ls-files '*.go')
go vet ./...
go test -race -coverprofile=coverage.out ./...
```

Project goals, phases and architecture decisions are documented in [`PROJECT.md`](PROJECT.md). Task progress can be found in the `tasks/` directory.
