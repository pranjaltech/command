Objective: Prevent empty prompt from calling OpenAI and show a helpful error.

Acceptance Criteria:
- Running `cmd` with no prompt exits with an error before hitting the API.
- README mentions that a prompt is required.
- Unit test covers the missing prompt case.
- All verification commands pass.

Implementation Checklist:
- [x] Add `Args: cobra.MinimumNArgs(1)` in root command.
- [x] Write failing test for no arguments.
- [x] Update README usage section.
- [x] Run golangci-lint, staticcheck, golines, go vet, go test -race.

Status: Completed
