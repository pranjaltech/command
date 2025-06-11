Objective: Add an end-to-end test that exercises the real OpenAI API using the OPENAI_API_KEY secret.

Acceptance Criteria:
- A Go test under test/e2e runs the compiled cmd binary with the real API key.
- Test asserts the command output is non-empty.
- Test is skipped when OPENAI_API_KEY is unset.
- All CI checks (golangci-lint, staticcheck, golines, go vet, go test -race) pass.

Implementation Checklist:
- [x] Create test/e2e directory and add new Go test.
- [x] Use os/exec to run `go run ./main.go <prompt>` with OPENAI_API_KEY set.
- [x] Skip test if OPENAI_API_KEY not set.
- [x] Record test result.
- [x] Run all verification commands.

Status: Completed
