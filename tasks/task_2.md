Objective: Implement environment probing to collect system context for LLM requests.

Acceptance Criteria:
- A new package `internal/probe` exposes `Probe` with a `Collect() (EnvInfo, error)` method.
- `EnvInfo` captures OS, Arch, kernel version, working directory, shell, git root, branch and dirty state.
- Unit tests cover probe logic using a stubbed command runner.
- `go test -race ./...` and static analysis pass.

Implementation Checklist:
- [x] Create `internal/probe` package with `EnvInfo` struct and `Probe` type.
- [x] Add an injectable command runner to avoid executing commands in tests.
- [x] Implement `Collect` using standard library and git commands.
- [x] Write unit tests for `Collect` using a stubbed runner.
- [x] Run golines, gofmt, staticcheck, go vet, and go test.

Environment notes:
- Tests must not depend on actual git or shell; use stubs.

Status: Completed

