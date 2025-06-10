Objective: Wire root command to the new environment probe and LLM client so user prompts generate suggestions via OpenAI.

Acceptance Criteria:
- `cmd <text>` uses `probe.Probe` and `llm.Client` to return a command suggestion.
- API key read from $OPENAI_API_KEY for default CLI.
- Unit tests stub the probe and LLM; no network calls.
- Static analysis and tests pass.

Implementation Checklist:
- [ ] Add `NewRootCmd` constructor accepting dependencies.
- [ ] Read OPENAI_API_KEY and create default clients in package init.
- [ ] Update `rootCmd` to collect env and call `GenerateCommand`.
- [ ] Update tests with stub implementations.
- [ ] Run golines, gofmt, staticcheck, go vet, golangci-lint, go test -race.

Environment notes:
- Tests must not hit network; use stubs.

Status: Completed
