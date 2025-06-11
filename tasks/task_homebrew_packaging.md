Objective: Add Homebrew formula and installation instructions so users can install via brew.

Acceptance Criteria:
- Formula file builds the CLI using Go and includes a test block.
- README documents installation using the raw formula URL.
- goreleaser config includes a brews section for generating the formula.
- Verification commands pass.

Implementation Checklist:
- [x] Add Formula/cmd.rb with go build and test.
- [x] Configure goreleaser to produce Homebrew formula.
- [x] Update README with brew install command.
- [x] Run golangci-lint, staticcheck, go vet, golines, go test -race.

Status: Completed
