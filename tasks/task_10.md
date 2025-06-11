Objective: Improve the TUI to display full command suggestions and allow editing before execution.

Acceptance Criteria:
- Command list shows complete strings without truncation.
- After selecting a command, user can edit it in a text input field.
- Selector returns the edited command on confirmation.
- Unit tests cover editing workflow.
- All verification commands pass.

Implementation Checklist:
- [x] Expand list width to accommodate long commands.
- [x] Add textinput component for editing the chosen command.
- [x] Update selector model and unit tests.
- [x] Run golines, gofmt, golangci-lint, staticcheck, go vet, go test -race.

Status: Completed
