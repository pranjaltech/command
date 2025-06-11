Objective: Provide user-friendly configuration management and fix uninstall script permissions.

Acceptance Criteria:
- `cmd config view` prints current settings with API key redacted.
- `cmd config set <field> <value>` updates the config file.
- README documents how to view and change configs.
- uninstall.sh uses sudo when removing from protected directories.
- All verification commands pass.

Status: Completed
