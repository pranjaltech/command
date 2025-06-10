Objective: Scaffold the cmd CLI tool and implement minimal natural language to command conversion.

Acceptance Criteria:
- `cmd <text>` prints a command suggestion derived from the input.
- Code organised with packages following Clean Architecture principles (internal/core, internal/adapters, cmd/).
- All tests pass with `go test -race ./...` and coverage is >95% (for now we just have minimal tests, so coverage may not reach target but we set target for small code). 
- Static analysis passes via `staticcheck` and `go vet`.

Implementation Checklist:
- [x] run `go mod init command`
- [x] run `go install github.com/spf13/cobra-cli@latest`
- [x] run `cobra-cli init --pkg-name command`
- [x] implement a simple converter function in internal/core that maps `"list all directories"` to `"ls -d */"` and returns unknown for others
- [x] wire root command to call the converter and print the resulting command
- [x] add unit tests for the converter and CLI command
- [x] create .codex/setup.sh to install dependencies (cobra-cli, staticcheck, golines, golangci-lint)
- [x] run formatting and static analysis
- [x] run go test -race ./...

Environment notes:
- Nothing special yet.
