# Release History

## Unreleased

### Cross-Platform Improvements

- **OS-Aware Tool Detection**: All tool checks now adapt to the operating system
  - Docker/Podman: Checks appropriate container runtime based on OS
  - Python: Checks `python3` on Unix-like systems, `python` on Windows
  - PowerShell: Checks `pwsh` then `powershell` on Windows; only `pwsh` on Unix
  - Bash: Standard on Unix; optional (Git Bash/WSL) on Windows
- **Enhanced Test Coverage**: Added comprehensive OS-specific tests for all check functions
- **Better Error Messages**: Tool names now reflect which command was actually found

## 0.1.1 - Bug Fix

- Fixed issue where `azd doctor` would fail when running in extension mode due to missing `listen` command.

## 0.1.0 - Documentation Update

- Updated installation instructions to use `doc` as the source name.
- Added upgrade instructions to README.

## 0.0.6 - Patch Release

- Project maintenance and documentation updates.

## 0.0.5 - Maintenance Release

- Maintenance updates and bug fixes.

## 0.0.4 - Remote Build Configuration

- Added `configure remote-build` command to easily enable remote builds for container apps.
- Updated `verify` command to suggest enabling remote build when Docker is missing.
- Fixed `azure.yaml` parsing to correctly map `remoteBuild` field.

## 0.0.3 - Terraform Focus

- Removed Bicep checks and related tests, focusing exclusively on Terraform infrastructure checks.

## 0.0.2 - Bypass Verification

- Added `AZD_DOCTOR_SKIP_VERIFY` environment variable to bypass verification checks.
- Added support for skipping verification in lifecycle hooks (e.g. `predeploy`).

## 0.0.1 - Initial Version