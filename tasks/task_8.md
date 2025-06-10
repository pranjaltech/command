Objective: Implement interactive suggestion selection using Bubbletea showing three commands returned by the LLM.

Acceptance Criteria:
- LLM client returns up to three command suggestions.
- CLI displays suggestions in a TUI list with arrow-key navigation (first selected by default).
- Pressing Enter executes the chosen command via the user's shell.
- Unit tests cover selection logic using stubs.
- All verification commands pass.

Implementation Checklist:
- [x] Update llm.Client to return []string from GenerateCommands.
- [x] Parse newline-delimited suggestions in OpenAI client implementation.
- [x] Add ui package with Bubbletea selector and interface.
- [x] Add shell runner interface for executing commands.
- [x] Wire root command to use selector and runner.
- [x] Update tests for new interfaces and behaviour.
- [x] Update e2e test to press Enter automatically.

Status: Completed
