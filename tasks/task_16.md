Objective: Expose model and temperature configuration and improve CLI help.

Acceptance Criteria:
- Config struct includes model and temperature fields with defaults.
- Flags --model and --temperature override and persist these settings.
- Help output describes configuration options.
- README documents model and temperature configuration.
- All verification commands pass.

Implementation Checklist:
- [x] Update config package with new fields and tests.
- [x] Extend OpenAI client constructor with model and temperature.
- [x] Modify root command to handle flags and save config.
- [x] Update README and help text.
- [x] Run golangci-lint, staticcheck, golines, go vet, go test -race.

Status: Completed
