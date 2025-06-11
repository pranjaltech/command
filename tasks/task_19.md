Objective: Add a debug mode that outputs verbose information for troubleshooting.

Acceptance Criteria:
- A global `--debug` flag writes details about environment collection, LLM prompts, responses and selected commands to stderr.
- README documents the new flag.
- Unit test covers debug logging in the LLM client.
- All verification commands pass.

Implementation Checklist:
- [x] Add persistent `--debug` flag in root command.
- [x] Add debug printing in root command.
- [x] Extend OpenAI client with optional debug writer and implement logging.
- [x] Write unit test for debug output.
- [x] Update README.
- [x] Run golangci-lint, staticcheck, golines, go vet, go test -race.

Status: Completed
