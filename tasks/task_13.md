Objective: Store the OpenAI API key in an encrypted config file and prompt the user the first time if it's missing.

Acceptance Criteria:
- On startup, the CLI loads configuration from $HOME/.config/cmd/config.yaml using viper.
- If the API key is not present, ask the user to input it on stdin.
- The key is encrypted before saving to config.yaml.
- Subsequent runs decrypt and use the stored key without prompting.
- All verification commands pass.

Implementation Checklist:
- [x] Add viper dependency and config loading logic.
- [x] Implement simple AES encryption/decryption utilities.
- [x] Prompt for API key when missing and write encrypted value.
- [x] Update unit tests for new behavior.
- [x] Run golines, golangci-lint, staticcheck, go vet, go test -race.

Status: Completed
