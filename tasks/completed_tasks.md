# Completed Tasks


## scaffold_cli

Objective: Scaffold the cmd CLI tool and implement minimal natural language to command conversion.

Acceptance Criteria:
- `cmd <text>` prints a command suggestion derived from the input.
- Code organised with packages following Clean Architecture principles (internal/core, internal/adapters, cmd/).
- All tests pass with `go test -race ./...` and coverage is >95% (for now we just have minimal tests, so coverage may not reach target but we set target for small code). 
- Static analysis passes via `staticcheck` and `go vet`.

Implementation Checklist:
- [x] run `go mod init command`
- [x] run `go install github.com/spf13/cobra-cli@latest`
- [x] run `cobra-cli init --pkg-name command`
- [x] implement a simple converter function in internal/core that maps `"list all directories"` to `"ls -d */"` and returns unknown for others
- [x] wire root command to call the converter and print the resulting command
- [x] add unit tests for the converter and CLI command
- [x] create .codex/setup.sh to install dependencies (cobra-cli, staticcheck, golines, golangci-lint)
- [x] run formatting and static analysis
- [x] run go test -race ./...

Environment notes:
- Nothing special yet.

Status: Completed


---

## tui_edit_suggestions

Objective: Improve the TUI to display full command suggestions and allow editing before execution.

Acceptance Criteria:
- Command list shows complete strings without truncation.
- After selecting a command, user can edit it in a text input field.
- Selector returns the edited command on confirmation.
- Unit tests cover editing workflow.
- All verification commands pass.

Implementation Checklist:
- [x] Expand list width to accommodate long commands.
- [x] Add textinput component for editing the chosen command.
- [x] Update selector model and unit tests.
- [x] Run golines, gofmt, golangci-lint, staticcheck, go vet, go test -race.

Status: Completed

---

## json_parse_objects

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

---

## error_handling_docs_install

Objective: Improve error handling and user docs, plus provide install scripts.

Acceptance Criteria:
- API errors from OpenAI are reported with a concise message without showing raw JSON details.
- CLI does not print usage when generation fails.
- README explains usage, development, installation and uninstallation.
- `scripts/install.sh` builds and installs the binary on macOS and Linux.
- `scripts/uninstall.sh` removes the installed binary.
- All verification commands pass.

Implementation Checklist:
- [x] Detect `openai.APIError` in `OpenAIClient.GenerateCommands` and wrap with friendly message.
- [x] Set `SilenceUsage` on the root command.
- [x] Add unit test for API error handling.
- [x] Write expanded README with installation instructions.
- [x] Create install/uninstall scripts.
- [x] Run golines, golangci-lint, staticcheck, go vet, go test -race.

Status: Completed

---

## encrypted_api_key_config

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

---

## mac_install_fix

Objective: Fix install script error on Mac

Acceptance Criteria:
- `scripts/install.sh` builds the binary despite existing cmd directory
- Installing results in binary copied to PREFIX/cmd
- Uninstall script remains unaffected
- All verification commands pass

Implementation Checklist:
- [x] Modify install.sh to build to a temporary file
- [x] Copy temporary binary to target location
- [x] Remove temporary file
- [x] Run verification commands

Status: Completed

---

## install_permission_error

Objective: Handle permission errors in install script

Acceptance Criteria:
- install.sh checks write permission and uses sudo when needed
- README explains running with sudo or custom PREFIX
- All verification commands pass

Implementation Checklist:
- [x] detect non-writable BIN_DIR and invoke sudo
- [x] update README with sudo note
- [x] run golangci-lint, staticcheck, golines, go vet, go test -race

Status: Completed

---

## configurable_model_temperature

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

---

## config_mgmt_uninstall_permissions

Objective: Provide user-friendly configuration management and fix uninstall script permissions.

Acceptance Criteria:
- `cmd config view` prints current settings with API key redacted.
- `cmd config set <field> <value>` updates the config file.
- README documents how to view and change configs.
- uninstall.sh uses sudo when removing from protected directories.
- All verification commands pass.

Status: Completed

---

## empty_prompt_error

Objective: Prevent empty prompt from calling OpenAI and show a helpful error.

Acceptance Criteria:
- Running `cmd` with no prompt exits with an error before hitting the API.
- README mentions that a prompt is required.
- Unit test covers the missing prompt case.
- All verification commands pass.

Implementation Checklist:
- [x] Add `Args: cobra.MinimumNArgs(1)` in root command.
- [x] Write failing test for no arguments.
- [x] Update README usage section.
- [x] Run golangci-lint, staticcheck, golines, go vet, go test -race.

Status: Completed

---

## debug_mode_verbose

Objective: Add a debug mode that outputs verbose information for troubleshooting.

Acceptance Criteria:
- A global `--debug` flag writes details about environment collection, LLM prompts, responses and selected commands to stderr.
- README documents the new flag.
- Unit test covers debug logging in the LLM client.
- All verification commands pass.

Implementation Checklist:
- [x] Add persistent `--debug` flag in root command.
- [x] Add debug printing in root command.
- [x] Extend OpenAI client with optional debug writer and implement logging.
- [x] Write unit test for debug output.
- [x] Update README.
- [x] Run golangci-lint, staticcheck, golines, go vet, go test -race.

Status: Completed

---

## env_probe

Objective: Implement environment probing to collect system context for LLM requests.

Acceptance Criteria:
- A new package `internal/probe` exposes `Probe` with a `Collect() (EnvInfo, error)` method.
- `EnvInfo` captures OS, Arch, kernel version, working directory, shell, git root, branch and dirty state.
- Unit tests cover probe logic using a stubbed command runner.
- `go test -race ./...` and static analysis pass.

Implementation Checklist:
- [x] Create `internal/probe` package with `EnvInfo` struct and `Probe` type.
- [x] Add an injectable command runner to avoid executing commands in tests.
- [x] Implement `Collect` using standard library and git commands.
- [x] Write unit tests for `Collect` using a stubbed runner.
- [x] Run golines, gofmt, staticcheck, go vet, and go test.

Environment notes:
- Tests must not depend on actual git or shell; use stubs.

Status: Completed


---

## shell_history_save

Objective: Save executed commands to the user's shell history.

Acceptance Criteria:
- After running a command via `cmd`, the command appears in the history file for bash, zsh and fish.
- History updates do not cause command execution to fail if writing fails.
- Unit tests cover history file updates for all three shells.
- All verification commands pass.

Implementation Checklist:
- [x] Add history appending logic in internal/shell.
- [x] Call history writer from the shell runner.
- [x] Write unit tests for bash, zsh and fish history entries.
- [x] Add end-to-end tests verifying history across shells.
- [x] Run verification commands.

Status: Completed

---

## llm_client

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


---

## integrate_env_llm

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

---

## dotenv_support

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

---

## missing_api_key_error

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

---

## e2e_real_openai

Objective: Add an end-to-end test that exercises the real OpenAI API using the OPENAI_API_KEY secret.

Acceptance Criteria:
- A Go test under test/e2e runs the compiled cmd binary with the real API key.
- Test asserts the command output is non-empty.
- Test is skipped when OPENAI_API_KEY is unset.
- All CI checks (golangci-lint, staticcheck, golines, go vet, go test -race) pass.

Implementation Checklist:
- [x] Create test/e2e directory and add new Go test.
- [x] Use os/exec to run `go run ./main.go <prompt>` with OPENAI_API_KEY set.
- [x] Skip test if OPENAI_API_KEY not set.
- [x] Record test result.
- [x] Run all verification commands.

Status: Completed

---

## tui_interactive_selection

Objective: Implement interactive suggestion selection using Bubbletea showing three commands returned by the LLM.

Acceptance Criteria:
- LLM client returns up to three command suggestions.
- CLI displays suggestions in a TUI list with arrow-key navigation (first selected by default).
- Pressing Enter executes the chosen command via the user's shell.
- Unit tests cover selection logic using stubs.
- All verification commands pass.

Implementation Checklist:
- [x] Update llm.Client to return []string from GenerateCommands.
- [x] Parse newline-delimited suggestions in OpenAI client implementation.
- [x] Add ui package with Bubbletea selector and interface.
- [x] Add shell runner interface for executing commands.
- [x] Wire root command to use selector and runner.
- [x] Update tests for new interfaces and behaviour.
- [x] Update e2e test to press Enter automatically.

Status: Completed

---

## structured_json_responses

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

---
## add_version_flag
Objective: Provide version information via `cmd --version`.

Acceptance Criteria:
- `cmd --version` prints the build version.
- Version is injected during build using ldflags.
- Unit test verifies `--version` output.
- Documentation explains how to display the version.
- Verification commands succeed.

Implementation Checklist:
- [x] Declare `Version` variable in the `cmd` package.
- [x] Set `rootCmd.Version` to that variable.
- [x] Add ldflags in `.goreleaser.yaml` to fill `Version`.
- [x] Add unit test for `--version` flag.
- [x] Mention the flag in README.
- [x] Run verification commands.

Status: Completed

---
## clarification_flow
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

---
## debug_mode_improvements
Objective: Improve debug mode output and UI.

Acceptance Criteria:
- Debug output logs the full system prompt including environment details.
- UI selector no longer shows the "List" title bar.
- Unit test updated for new debug logs.
- All verification commands pass.

Implementation Checklist:
- [x] Log system prompt and user prompt separately in OpenAI client when debug is enabled.
- [x] Hide list title in selector model.
- [x] Update debug output unit test.
- [x] Run verification commands.

Status: Completed

---
## homebrew_packaging
Objective: Add Homebrew formula and installation instructions so users can install via brew.

Acceptance Criteria:
- Formula file builds the CLI using Go and includes a test block.
- README documents installation using the raw formula URL.
- goreleaser config includes a brews section for generating the formula.
- Verification commands pass.

Implementation Checklist:
- [x] Add Formula/cmd.rb with go build and test.
- [x] Configure goreleaser to produce Homebrew formula.
- [x] Update README with brew install command.
- [x] Run golangci-lint, staticcheck, go vet, golines, go test -race.

Status: Completed

---
## homebrew_versioned_release
Objective: Implement versioned release workflow with Homebrew tap

Acceptance Criteria:
- goreleaser builds cross-platform binaries and generates a Homebrew formula
  pointing at release tarballs
- README explains installing the stable release from the pranjaltech/tools tap
  and installing development versions from a raw formula URL
- Tap repository is configurable so generated .rb files can be copied to
  https://github.com/pranjaltech/homebrew-tools
- Verification commands pass

Implementation Checklist:
- [x] Update goreleaser config with project name, release settings and tap info
- [x] Document stable and development installation methods in README
- [x] Run golangci-lint, staticcheck, golines, go vet, go test -race

Status: Completed

---
## llm_prompt_improvement
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

---
## path_scan_cache
Objective: Cache PATH binary scans to avoid repeated directory reads when building the system prompt.

Acceptance Criteria:
- uniqueBinaries caches its result so repeated calls don't rescan the PATH.
- Unit tests and verification commands continue to pass.

Implementation Checklist:
- [x] Add sync.Once and cached list in openai.go.
- [x] Update buildSystemPrompt to use cached data.
- [x] Run verification commands.

Status: Completed


---
## remove_formula_folder
Objective: Remove obsolete Homebrew Formula directory.

Acceptance Criteria:
- `Formula/` folder and files are deleted.
- README has no references to the formula.
- `goreleaser` config contains only cask settings.
- Verification commands succeed.

Implementation Checklist:
- [x] Delete `Formula/cmd.rb`.
- [x] Ensure README mentions only Homebrew casks.
- [x] Run verification commands.

Status: Completed

---
## switch_to_homebrew_cask
Objective: Replace deprecated brews section with homebrew_casks and update documentation.

Acceptance Criteria:
- `.goreleaser.yaml` uses `homebrew_casks` configuration.
- README shows `brew install --cask` instructions for stable and development versions.
- Development cask file exists under `Casks/`.
- Verification commands succeed.

Implementation Checklist:
- [x] Update goreleaser config.
- [x] Update README with dev cask URL.
- [x] Add `Casks/cmd-dev.rb`.
- [x] Run goreleaser check.
- [x] Run verification commands.

Status: Completed

---
## unify_dev_cask
Objective: Use a single Homebrew cask file `cmd.rb` for both stable and development installs.

Acceptance Criteria:
- `Casks/cmd.rb` exists and contains the cask definition.
- No `Casks/cmd-dev.rb` file remains.
- README refers to the raw `cmd.rb` URL for development installation.
- Verification commands succeed.

Implementation Checklist:
- [x] Remove `cmd-dev.rb`.
- [x] Add `cmd.rb` with the same content.
- [x] Update README installation section.
- [x] Run `goreleaser check`.
- [x] Run verification commands.

Status: Completed

---
## update_goreleaser_config
Objective: Update goreleaser configuration to version 2 to resolve release errors.

Acceptance Criteria:
- `goreleaser check` passes without errors.
- Release workflow works with `goreleaser release --clean`.
- Verification commands succeed.

Implementation Checklist:
- [x] Change `.goreleaser.yaml` to start with `version: 2`.
- [x] Run `goreleaser check` to verify config (initially shows deprecation warnings).
- [x] Update deprecated fields under `archives` to use `formats` keys.
- [x] Re-run `goreleaser check` until no warnings remain.
- [x] Run verification commands.

Status: Completed

---
