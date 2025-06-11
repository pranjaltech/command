Objective: Improve debug mode output and UI.

Acceptance Criteria:
- Debug output logs the full system prompt including environment details.
- UI selector no longer shows the "List" title bar.
- Unit test updated for new debug logs.
- All verification commands pass.

Implementation Checklist:
- [x] Log system prompt and user prompt separately in OpenAI client when debug is enabled.
- [x] Hide list title in selector model.
- [x] Update debug output unit test.
- [ ] Run verification commands.

Status: In Progress
