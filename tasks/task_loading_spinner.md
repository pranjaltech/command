Objective: Improve UX with loading animation and centralized debug logging.

Acceptance Criteria:
- A simple spinner shows while waiting for LLM responses.
- Debug messages are printed using a reusable helper and appear in dim gray.
- Existing tests updated for new debug formatting.
- Verification commands pass.

Implementation Checklist:
- [x] Add ui.Loader for spinner.
- [x] Create log package with Enable and Debugf helpers.
- [x] Refactor root and OpenAI client to use log.Debugf.
- [x] Update tests for colored debug output.
- [x] Run gofmt, golines and verification commands.

Status: Completed
