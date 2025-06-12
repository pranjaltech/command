Objective: Add onboarding with multi-provider and telemetry support.

Acceptance Criteria:
- Running `cmd` without args shows help text.
- `cmd config` triggers onboarding when config missing, otherwise shows help.
- Config stores provider, API key, API URL and telemetry preference; multiple providers can be saved.
- Onboarding prompts user for provider, API key, URL (with defaults), and telemetry opt-in.
- Langfuse telemetry is sent only if user enables it.
- README split into user and developer sections with contribution guidelines.
- Verification suite passes.

Implementation Checklist:
- [x] Extend Config struct with provider map and telemetry flag.
- [x] Update Save/Load and tests for new fields.
- [x] Add onboarding function in cmd package.
- [x] Adjust root and config commands to invoke onboarding when needed.
- [x] Integrate langfuse-go library behind telemetry package.
- [x] Record generation events when telemetry enabled.
- [x] Revise README structure.
- [x] Run formatting and verification commands.

Status: Completed
