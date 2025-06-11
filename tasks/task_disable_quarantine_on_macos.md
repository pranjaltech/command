Objective: Ensure macOS installs don't require command signing by removing the quarantine attribute.

Acceptance Criteria:
- `.goreleaser.yaml` defines a post-install hook to clear `com.apple.quarantine`.
- `Casks/cmd.rb` includes the same post-install steps.
- `goreleaser check` passes.
- Verification commands succeed.

Implementation Checklist:
- [x] Update `.goreleaser.yaml` with `homebrew_casks` post-install hook.
- [x] Update `Casks/cmd.rb` with matching uninstall step.
- [x] Run `goreleaser check`.
- [x] Run verification commands.

Status: Completed
