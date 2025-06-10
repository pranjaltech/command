# PROJECT.md – cmd CLI

**Purpose:** High‑level product requirements & phased implementation plan for the AI‑assisted command‑line tool (cmd). This file guides product discussions, milestone planning, and progress tracking. See **AGENTS.md** for day‑to‑day engineering rules.

## 1\. Product Vision

Give every developer and power user an “AI shell companion” that turns plain English into safe, context‑aware shell commands—while keeping humans in the driver’s seat.

- **Audience:** Software engineers, data scientists, DevOps, and tech‑savvy users on macOS, Linux, and Windows.
- **Value Props**
    1. ✨ _Speed:_ Type intent once; run the right command without googling flags.
    2. 🛡️ _Safety:_ Review/edit every suggestion before execution.
    3. 🧠 _Context‑aware:_ Detect env (Python venv, git status, etc.) to craft correct commands.
    4. 🔧 _Extensible:_ Pluggable LLMs, custom presets, and scripting hooks.

## 2\. Functional Requirements

| ID  | Description | Priority |
| --- | --- | --- |
| **F‑1** | Parse natural‑language prompt (stdin/argv) and request command suggestions from selected LLM. | Must‑have |
| **F‑2** | Collect environment context (OS, cwd, VCS, virtual env) and include in LLM payload. | Must‑have |
| **F‑3** | Render 1‑N suggestions in a TUI list, editable with arrow keys & inline editing. | Must‑have |
| **F‑4** | Execute the confirmed command in the user’s shell; stream output in real time. | Must‑have |
| **F‑5** | Interactive config sub‑command (cmd config) with persistence under $HOME/.config/cmd. | Must‑have |
| **F‑6** | Offline fallback: if LLM unavailable, provide cached commands or graceful error. | Should |
| **F‑7** | Telemetry opt‑in to anonymously measure usage & failures. | Should |
| **F‑8** | Plugin hooks: allow custom command post‑processors or validators. | Could |

### Non‑functional

- **Performance:** Suggestions < 1.5 s (network permitting). Cold start < 200 ms.
- **Reliability:** 99.5 % success rate in executing command after confirmation.
- **Security:** No command runs without explicit user approval. All secrets stay local.
- **Portability:** Single statically‑linked binary per OS/Arch via goreleaser.
- **Test Coverage:** ≥ 95 % line & branch.

## 3\. High‑Level Architecture

deploymentDiagram  
user-->cmd: prompt & keystrokes  
cmd-->envProbe: gather context  
cmd-->llmClient: JSON payload  
llmClient-->LLM(API)  
LLM-->cmd: command suggestions  
cmd-->shellExec: confirmed command  
shellExec-->cmd: stdout/stderr

- **cmd (CLI layer):** Cobra‑generated command tree, Bubbletea TUI.
- **envProbe:** Pure‑Go utilities under internal/probe.
- **llmClient:** Interface + OpenAI implementation; future providers.
- **shellExec:** Wrapper using exec.Cmd with context cancellation.

## 4\. Implementation Roadmap

| Phase | Timeline | Deliverables |
| --- | --- | --- |
| **0\. Bootstrap** | Week 1 | Repo scaffold (cobra-cli init), CI pipeline, AGENTS.md, PROJECT.md. |
| **1\. Core Engine** | Weeks 2‑3 | envProbe module, llmClient interface + OpenAI impl (unit‑tested). |
| **2\. TUI MVP** | Weeks 4‑5 | Bubbletea list view, inline editor, command execution, basic config file. |
| **3\. Testing & Hardening** | Weeks 6‑7 | 95 % coverage, E2E expect tests, race‑free, golden files. |
| **4\. Packaging** | Week 8 | goreleaser config, Homebrew tap, GitHub releases, docs site stub. |
| **5\. Nice‑to‑haves** | Weeks 9‑10 | Plugins, offline cache, telemetry, additional LLM providers. |

**Milestone exit criteria**: Each phase must pass CI, meet coverage target, and have an updated changelog.

## 5\. Risks & Mitigations

| Risk | Impact | Likelihood | Mitigation |
| --- | --- | --- | --- |
| LLM latency / cost spikes | Poor UX | Medium | Add caching; allow local models. |
| Incorrect or dangerous commands | High | Medium | Strict confirmation flow; future static analysis of command. |
| Windows shell quirks | Medium | Medium | Abstract shellExec; unit test on Windows CI runner. |
| API key leakage in logs | High | Low | Redact tokens; integrate secret scanner in CI. |

## 6\. LLM & Configuration Decisions

### 6.1 LLM Providers (v1)

- **Cloud first:** Ship with **OpenAI** (GPT‑4o‑mini default). Roadmap: add **Google Gemini** and **Anthropic Claude** once stable Go SDKs are available.
- **Local mode:** Integrate with **Ollama** so users can load any GGUF model. Provide an automated download of a quantised **Code Llama‑Instruct 7B** on first run for offline usage.

### 6.2 Prompt Memory Strategy

For 1.0 we remain **stateless**—one prompt → suggestions → execute. No per‑directory or global chat memory. (Potential future enhancement: memory-v2.)

### 6.3 Config Versioning & Migration

Include a config_version field in config.yaml. If the binary detects an older version, it applies a small patch migrator bundled with the release. Full migration tooling will be addressed post‑1.0 once the schema stabilises.

## 7\. Acceptance Criteria v1.0

Acceptance Criteria v1.0

- ✅ User can run cmd "list all directories here" and get at least one correct ls -d \*/ suggestion.
- ✅ Works on macOS (zsh) and Ubuntu (bash) in Git repo and in plain dir.
- ✅ cmd config lets user switch LLM model and suggestion count interactively.
- ✅ Installation via brew install cmd and pre‑built Linux/Windows binaries.
- ✅ Entire test suite go test ./... green with -race and ≥ 95 % coverage.

### Appendix A – Key Dependencies & Versions

| Tool | Min Version | Reason |
| --- | --- | --- |
| Go  | 1.22 | Generics maturity & perf. |
| Cobra | v1.8 | CLI scaffold. |
| Bubbletea | v0.25 | TUI. |
| Viper | v1.18 | Config. |
| go‑openai | latest | LLM API. |
| gomock | v1.6 | Tests. |
| goreleaser | v2.x | Distribution. |
