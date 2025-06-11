# <a name="agents.md"></a>AGENTS.md
**Project:** cmd – AI‑assisted command‑line tool (Go)

**Scope of this file:** Whole repository. If another AGENTS.md exists deeper in the tree, its rules override these for that subtree (Codex follows the *closest* file).

-----
## <a name="why-this-file-exists"></a>0. Why this file exists
OpenAI **Codex** (cloud agent) and human contributors both read this document to understand *how* work should be done here. Anything ambiguous in a task prompt should default to the guidelines below.

*If a rule here conflicts with a newer official best practice, Codex must fetch the latest docs and propose an update alongside its code.*

-----
## <a name="highlevel-goal"></a>1. High‑Level Goal
Build a **maintainable, robust, and performant Go CLI** that converts natural‑language prompts (e.g. cmd list all directories here) into safe shell commands after giving the user a chance to review/edit. Written with **TDD** and shipped with **≥ 95 % test coverage**.

-----
## <a name="architecture-principles"></a>2. Architecture Principles

|Guideline|Enforcement|
| :- | :- |
|**Clean Architecture**|Domain logic in internal/core; adapters in internal/adapters/...; CLI glue in cmd/ generated via [cobra-cli].|
|**Small packages** (<400 LOC)|Split by responsibility; avoid cyclic imports.|
|**Dependency Injection**|Accept interfaces; use constructors returning concrete structs.|
|**Pure Functions First**|Side‑effect‑free code is easier to test and review.|
|**Configurable & Mockable I/O**|LLM, OS probes, and exec runner are behind interfaces so unit tests don’t hit network or spawn processes.|

-----
## <a name="coding-standards"></a>3. Coding Standards
- Go ≥ 1.22 — use Modules
- go fmt, go vet, staticcheck, golangci‑lint **MUST** pass.
- Wrap errors (fmt.Errorf("%w", err)) and provide context.
- No naked panic outside of main() startup.
- Follow [Effective Go] and [Uber Go Style] where not conflicting.
-----
## <a name="dependencies-tooling"></a>4. Dependencies & Tooling

|Purpose|Tool|Notes|
| :- | :- | :- |
|CLI scaffold|**cobra‑cli**|Run cobra-cli init cmd for new commands.|
|TUI|**bubbletea + bubbles**|For arrow‑key editing & confirmation UI.|
|LLM client|**go‑openai** (github.com/sashabaranov/go-openai)|Abstract behind LLMClient interface.|
|Config|**viper**|Reads & hot‑reloads YAML at $HOME/.config/cmd/config.yaml.|
|Mocks|**gomock** (go generate)|Keep tests hermetic.|
|Static analysis|**staticcheck**, **golangci‑lint**|Enforced in CI.|
|Release|**goreleaser**|Cross‑compile & create Homebrew tap.|

**Codex directive:** Whenever a *scaffolding tool* exists (e.g. cobra-cli, goreleaser init, go test -c), **call the tool**; do *not* hand‑write boilerplate.

-----
## <a name="environment-configuration-for-codex"></a>5. Environment Configuration for Codex
Codex runs each task inside the openai/codex-universal Docker image.

1. A repo‑local setup script at .codex/setup.sh installs extra deps (e.g. go install honnef.co/go/tools/cmd/staticcheck@latest).
1. The script should assume **internet is available only during setup**; the agent phase has internet *off* unless the task explicitly enables it.
1. Secrets must be provided via the workspace UI; **never commit secrets**.
1. Use $CODEX\_PROXY\_CERT for any custom tooling that needs to trust the MITM proxy cert.
-----
## <a name="environment-awareness-logic-app-runtime"></a>6. Environment Awareness Logic (App runtime)
Before calling the LLM, collect and JSON‑encode:

- OS, ARCH, kernel version
- Current directory
- Active virtual envs (Python venv/conda, Node nvm, Go workspace, Xcode, etc.)
- Git repo status (root, branch, dirty flag)
- Shell ($SHELL) and version

Cache expensive probes per invocation; do not shell‑out when std lib suffices.

-----
## <a name="configuration-management-userlevel"></a>7. Configuration Management (User‑level)
- Default location: $HOME/.config/cmd/config.yaml (override $CMD\_CONFIG).
- Managed through cmd config interactive sub‑command (Bubbletea form).
- Versioned; migrate automatically if schema changes.

llm**:**\
`  `provider**:** openai\
`  `model**:** gpt-4o-mini\
`  `temperature**:** 0.2\
ui**:**\
`  `suggestions**:** 3   *# Number of command options shown*\
telemetry**:**\
`  `disable**:** false

-----
## <a name="commandgeneration-flow-runtime"></a>8. Command‑Generation Flow (Runtime)
flowchart TD\
`    `A[User types "cmd <prompt>"] --> B(Build Context)\
`    `B --> C(LLM Request)\
`    `C --> D{LLM Responses}\
`    `D --> E[Render suggestions]\
`    `E -->|Edit/Confirm| F(Execute)\
`    `E -->|Abort| G(Exit)

- Never execute without explicit confirmation (Enter key).
- Use exec.CommandContext with user’s login shell.
- Stream stdout/stderr in a scrollable pane.
-----
## <a name="testing-strategy"></a>9. Testing Strategy

|Layer|Tooling|Notes|
| :- | :- | :- |
|Unit|go test, **gomock**|No network, file‑system, or shell.|
|Integration|go test ./test/...|Spawn compiled binary in tmp dir, fake env vars.|
|End‑to‑end|expect scripts in test/e2e|Record TTY interaction.|

- Coverage target **≥ 95 %** (go test -covermode=atomic).
- Run go test -race and staticcheck in CI.
- Golden files stored in testdata/. Regenerate via go test ./... -update.
-----
## <a name="x0701ab4f8caa69a491757eb1c163f0258f59701"></a>10. Continuous Verification Checklist (Codex must run before PR)
./codex/setup.sh         *# ensure fresh env similar to agent*\
golangci-lint run ./...\
staticcheck ./...\
golines -m 120 -w $(git ls-files '*.go')\
go vet ./...\
go test -race -coverprofile=coverage.out ./...

-----
## <a name="pullrequest-template-autofilled-by-codex"></a>11. Pull‑Request Template (Autofilled by Codex)
\### What & Why\
\- \
\### How\
\- \
\### Tests Added/Updated\
\- [ ] Unit  -details\
\- [ ] Integration\
\- [ ] E2E\
\### Verification\
\```bash\
$ golangci-lint run ./... && staticcheck ./... && go test -race ./...
### <a name="checklist"></a>Checklist

\---\
\
\## 12. Self‑Improvement Loop\
1\. \*\*Doc watch\*\* – Before using any external API or Go stdlib added after Go 1.22, fetch docs & adjust.\
2\. \*\*Performance audit\*\* – On every 10th commit, run `go test -bench=.` and attach a benchmark diff if changes touch hot paths.\
3\. \*\*Dep upgrades\*\* – Weekly: `go get -u ./...` then `go mod tidy` guarded by tests.\
\
\---\
\
\## 13. Security & Account Hygiene\
\* Codex will \*never\* be given production credentials; tests should run in disposable containers.\
\* Do not hard‑code API keys; use env vars & Vault/1Password when manual testing.\
\* All contributors must have MFA enabled on GitHub.\
\
\---\
\
\## 14. Prompting Tips for Codex Tasks (repo‑specific)\
\* \*\*Start small\*\* – Ask Codex to scaffold a single command or write a failing test first.\
\* \*\*Provide repro\*\* – Include sample env variables & expected terminal output.\
\* \*\*Verify\*\* – End each task asking Codex to run `./scripts/ci\_local.sh` and report results.\
\
\---\
\
\## 15. Task‑File Workflow (Modus Operandi)\
Codex and human contributors \*\*must\*\* manage feature or bug work through explicit \*task files\* to ensure traceability and incremental progress.\
\
1\. \*\*Create a task stub\*\*  \
`   `\*Location:\* `docs/tasks/` or `tasks/` (create if absent).  \
`   `*Filename schema:* `task_<short-description>.md` using lowercase words separated by underscores (e.g. `task_add_login_flag.md`).\
2\. \*\*Plan before code\*\*  \
`   `Inside the file list:\
`   `\* Objective / user story\
`   `\* Acceptance criteria\
`   `\* Implementation checklist (TDD‑oriented)\
`   `\* Environment or setup notes\
3\. \*\*Iterative development loop\*\*  \
`   `a. Pick the \*next unchecked item\* in the checklist.  \
`   `b. Write a failing test (TDD).  \
`   `c. Implement the minimal code to pass the test.  \
`   `d. Run `go test -race ./...` and static analysis.\
4\. \*\*Update the task file\*\*  \
`   `Record status (`✅ done` / `🚧 in‑progress`), known bugs, `TODO`, or `FIX\_LATER` items after each green test run.\
5\. \*\*Reference before next chunk\*\*  \
`   `Always open the `task_<short-description>.md` and confirm remaining work before coding again.\
6\. \*\*Close‑out\*\*  \
`   `When all acceptance criteria are met and CI passes, mark the task file as \*\*Completed\*\* and reference it in the PR description.\
\
`   `Move completed tasks into `completed_tasks.md` to keep the directory tidy.
\> \*Rationale:\* This provides a transparent audit trail for Codex actions, eases context loading, and enforces disciplined TDD.\
\
\### Remember\
\* \*\*No quick fixes\*\* – optimise for maintainability.\
\* Use official tools for scaffolding.\
\* Keep logs clear, errors actionable, builds deterministic.\
\* Unsure?  Fetch docs ➜ update tests ➜ then code.\
\
[Effective Go]: https://go.dev/doc/effective\_go\
[Uber Go Style]: https://github.com/uber-go/guide/blob/master/style.md\
[`cobra-cli`]: https://github.com/spf13/cobra
