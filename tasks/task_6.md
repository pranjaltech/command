Objective: Gracefully handle missing OPENAI_API_KEY so the CLI returns an informative error instead of failing at API request time.

Acceptance Criteria:
- When the key is absent, running the command returns an error "OPENAI_API_KEY not set".
- Unit tests cover this case.
- Static analysis and tests continue to pass.

Implementation Checklist:
- [x] Update llm.NewOpenAIClient to return an error if api key is empty.
- [x] Root command checks for nil client and fails early.
- [x] Add unit test for missing client.
- [x] Run golines, gofmt, golangci-lint, staticcheck, go vet, go test -race.

Status: Completed
