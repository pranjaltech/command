Objective: Load environment variables from a `.env` file so local runs can easily set credentials.

Acceptance Criteria:
- On startup the CLI loads variables using `godotenv.Load()`.
- Missing `.env` file is not treated as an error.
- New dependency `github.com/joho/godotenv` is added and tests still pass.
- CI checks (golines, gofmt, staticcheck, golangci-lint, go vet, go test -race) succeed.

Implementation Checklist:
- [x] Add `godotenv` to go.mod and tidy.
- [x] Call `godotenv.Load()` in `cmd` package init.
- [x] Run formatting, static analysis, and tests.

Status: Completed
