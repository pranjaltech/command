Objective: Use structured JSON responses from OpenAI instead of parsing plain text.

Acceptance Criteria:
- OpenAI client requests `json_object` response format.
- System prompt instructs the model to return `{\"commands\": [<cmd>]}`.
- Response is unmarshaled and trimmed to at most three commands.
- Unit test verifies request and parsing logic.
- All verification commands pass.

Implementation Checklist:
- [ ] Update `OpenAIClient.GenerateCommands` to set `ResponseFormat` and parse JSON.
- [ ] Adjust unit tests for new format.
- [ ] Run gofmt, golines, golangci-lint, staticcheck, go vet, go test -race.

Status: Completed
