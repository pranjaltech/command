Objective: Improve error handling and user docs, plus provide install scripts.

Acceptance Criteria:
- API errors from OpenAI are reported with a concise message without showing raw JSON details.
- CLI does not print usage when generation fails.
- README explains usage, development, installation and uninstallation.
- `scripts/install.sh` builds and installs the binary on macOS and Linux.
- `scripts/uninstall.sh` removes the installed binary.
- All verification commands pass.

Implementation Checklist:
- [x] Detect `openai.APIError` in `OpenAIClient.GenerateCommands` and wrap with friendly message.
- [x] Set `SilenceUsage` on the root command.
- [x] Add unit test for API error handling.
- [x] Write expanded README with installation instructions.
- [x] Create install/uninstall scripts.
- [x] Run golines, golangci-lint, staticcheck, go vet, go test -race.

Status: Completed
