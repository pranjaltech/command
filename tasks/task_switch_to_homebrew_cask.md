Objective: Replace deprecated brews section with homebrew_casks and update documentation.

Acceptance Criteria:
- `.goreleaser.yaml` uses `homebrew_casks` configuration.
- README shows `brew install --cask` instructions for stable and development versions.
- Development cask file exists under `Casks/`.
- Verification commands succeed.

Implementation Checklist:
- [x] Update goreleaser config.
- [x] Update README with dev cask URL.
- [x] Add `Casks/cmd-dev.rb`.
- [x] Run goreleaser check.
- [x] Run verification commands.

Status: Completed
