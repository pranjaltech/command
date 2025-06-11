Objective: Remove obsolete Homebrew Formula directory.

Acceptance Criteria:
- `Formula/` folder and files are deleted.
- README has no references to the formula.
- `goreleaser` config contains only cask settings.
- Verification commands succeed.

Implementation Checklist:
- [x] Delete `Formula/cmd.rb`.
- [x] Ensure README mentions only Homebrew casks.
- [x] Run verification commands.

Status: Completed
