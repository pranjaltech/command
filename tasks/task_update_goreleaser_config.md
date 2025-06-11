Objective: Update goreleaser configuration to version 2 to resolve release errors.

Acceptance Criteria:
- `goreleaser check` passes without errors.
- Release workflow works with `goreleaser release --clean`.
- Verification commands succeed.

Implementation Checklist:
- [x] Change `.goreleaser.yaml` to start with `version: 2`.
- [x] Run `goreleaser check` to verify config. (shows deprecation warnings)
- [x] Run verification commands.

Status: Completed
