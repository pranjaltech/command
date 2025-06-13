Objective: Remove excess blank lines from the TUI so lists and loaders display compactly.

Acceptance Criteria:
- Command selection list has no blank space above the help line.
- Onboarding pickers display without empty lines between header and options.
- Verification commands pass.

Implementation Checklist:
- [x] Adjust list help style to remove top padding.
- [x] Ensure loader doesn't emit extra newlines.
 - [x] Update unit tests if needed.

Status: Completed
