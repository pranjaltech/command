Objective: Provide version information via `cmd --version`.

Acceptance Criteria:
- `cmd --version` prints the build version.
- Version is injected during build using ldflags.
- Unit test verifies `--version` output.
- Documentation explains how to display the version.
- Verification commands succeed.

Implementation Checklist:
- [x] Declare `Version` variable in the `cmd` package.
- [x] Set `rootCmd.Version` to that variable.
- [x] Add ldflags in `.goreleaser.yaml` to fill `Version`.
- [x] Add unit test for `--version` flag.
- [x] Mention the flag in README.
- [x] Run verification commands.

Status: Completed
