Objective: Handle permission errors in install script

Acceptance Criteria:
- install.sh checks write permission and uses sudo when needed
- README explains running with sudo or custom PREFIX
- All verification commands pass

Implementation Checklist:
- [x] detect non-writable BIN_DIR and invoke sudo
- [x] update README with sudo note
- [x] run golangci-lint, staticcheck, golines, go vet, go test -race

Status: Completed
