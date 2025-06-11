Objective: Use a single Homebrew cask file `cmd.rb` for both stable and development installs.

Acceptance Criteria:
- `Casks/cmd.rb` exists and contains the cask definition.
- No `Casks/cmd-dev.rb` file remains.
- README refers to the raw `cmd.rb` URL for development installation.
- Verification commands succeed.

Implementation Checklist:
- [x] Remove `cmd-dev.rb`.
- [x] Add `cmd.rb` with the same content.
- [x] Update README installation section.
- [x] Run `goreleaser check`.
- [x] Run verification commands.

Status: Completed
