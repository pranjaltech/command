Objective: Preserve executed commands in shell history

Acceptance Criteria:
- Runner appends executed commands to the user's shell history (bash, zsh, fish).
- Unit tests cover bash, zsh, and fish cases using temporary history files.
- E2E tests verify history updates when real shells are available.
- All verification commands pass.

Implementation Checklist:
- [x] Update shell runner to write commands to history files.
- [x] Add unit tests for history preservation per shell.
- [x] Add e2e tests for bash, zsh and fish shells.
- [x] Run golangci-lint, staticcheck, golines, go vet, go test -race.

Status: Completed
