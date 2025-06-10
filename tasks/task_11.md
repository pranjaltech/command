Objective: Fix JSON parsing of OpenAI command suggestions when items are objects.

Acceptance Criteria:
- `OpenAIClient.GenerateCommands` handles arrays of strings or objects with a `command` field.
- Unit tests cover both response formats.
- All verification commands pass.

Implementation Checklist:
- [ ] Update parsing logic in `OpenAIClient`.
- [ ] Add unit tests for object responses.
- [ ] Run gofmt, golines, golangci-lint, staticcheck, go vet, go test -race.

Status: Completed
