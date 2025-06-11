Objective: Cache PATH binary scans to avoid repeated directory reads when building the system prompt.

Acceptance Criteria:
- uniqueBinaries caches its result so repeated calls don't rescan the PATH.
- Unit tests and verification commands continue to pass.

Implementation Checklist:
- [x] Add sync.Once and cached list in openai.go.
- [x] Update buildSystemPrompt to use cached data.
- [x] Run verification commands.

Status: Completed

