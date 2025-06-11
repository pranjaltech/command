Objective: Implement versioned release workflow with Homebrew tap

Acceptance Criteria:
- goreleaser builds cross-platform binaries and generates a Homebrew formula
  pointing at release tarballs
- README explains installing the stable release from the pranjaltech/tools tap
  and installing development versions from a raw formula URL
- Tap repository is configurable so generated .rb files can be copied to
  https://github.com/pranjaltech/homebrew-tools
- Verification commands pass

Implementation Checklist:
- [x] Update goreleaser config with project name, release settings and tap info
- [x] Document stable and development installation methods in README
- [x] Run golangci-lint, staticcheck, golines, go vet, go test -race

Status: Completed
