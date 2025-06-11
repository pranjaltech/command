Objective: Handle LLM need_clarification responses by prompting the user and retrying.

Acceptance Criteria:
- OpenAI client returns a typed error when the response includes a non-empty `need_clarification` string.
- Root command asks the user the question and sends the extra information to the LLM (up to two times).
- Unit tests cover parsing of the field and the clarification loop.
- All verification commands pass.

Implementation Checklist:
- [x] Add NeedClarificationError type and parse logic in OpenAI client.
- [x] Update root command to prompt the user and retry when this error occurs.
- [x] Extend openai_test.go to cover the new error.
- [x] Add root command unit test simulating clarification flow.
- [x] Run gofmt, golines, golangci-lint, staticcheck, go vet, go test -race.

Status: Completed
