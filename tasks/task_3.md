Objective: Implement a basic LLM client interface with an OpenAI implementation to generate command suggestions.

Acceptance Criteria:
- A new package `internal/llm` defines `Client` and `OpenAIClient`.
- `GenerateCommand(ctx,prompt,env)` calls OpenAI chat completions with environment data.
- Unit tests use a stubbed chat client so no network calls occur.
- Static analysis and tests pass.

Implementation Checklist:
- [x] Add go-openai dependency.
- [x] Create `internal/llm` with interface and OpenAI client.
- [x] Implement `GenerateCommand` with environment JSON in system prompt.
- [x] Write unit tests using stub chat client.
- [x] Run golines, gofmt, staticcheck, go vet, go test.

Environment notes:
- No network access during tests; rely on stubs only.

Status: Completed

