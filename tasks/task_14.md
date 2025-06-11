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
