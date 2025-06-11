Objective: Enhance the LLM prompt with richer context and stricter output rules.

Acceptance Criteria:
- System prompt follows the example from the repo README with environment section, dynamic context placeholders, output contract, style rules and few-shot examples.
- OpenAI client uses this prompt template when generating commands.
- Unit tests updated to validate new prompt pieces.
- All verification commands pass.

Implementation Checklist:
- [x] Add system prompt template and helper to build it with environment info.
- [x] Update OpenAI client to use new template.
- [x] Extend unit tests in openai_test.go for the new prompt.
- [x] Run gofmt, golines, golangci-lint, staticcheck, go vet, go test -race.

Status: Completed
