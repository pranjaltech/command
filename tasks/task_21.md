Objective: Save executed commands to the user's shell history.

Acceptance Criteria:
- After running a command via `cmd`, the command appears in the history file for bash, zsh and fish.
- History updates do not cause command execution to fail if writing fails.
- Unit tests cover history file updates for all three shells.
- All verification commands pass.

Implementation Checklist:
- [x] Add history appending logic in internal/shell.
- [x] Call history writer from the shell runner.
- [x] Write unit tests for bash, zsh and fish history entries.
- [x] Add end-to-end tests verifying history across shells.
- [x] Run verification commands.

Status: Completed
