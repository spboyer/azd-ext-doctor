# Copilot Instructions for azd-ext-doctor

## Architecture Overview
This is a Go-based CLI extension for the Azure Developer CLI (`azd`). It is built using the **Cobra** library.
- **Entry Point**: `main.go` initializes the root command.
- **Commands**: Located in `internal/cmd/`. Key commands:
  - `check`: General health check of tools.
  - `verify`: Project-specific validation based on `azure.yaml`.
  - `context`: Displays current environment context.
- **Core Logic**: `internal/checks/` contains the implementation of tool checks (e.g., Docker, Node, Python).
- **Manifests**:
  - `extension.yaml`: Defines extension metadata, capabilities, and lifecycle hooks.
  - `registry.json`: Used for distribution, mapping versions to platform-specific artifacts.

## Critical Workflows

### Build
- **Standard**: `go build -o azd-ext-doctor .`
- **Scripts**: Use `build.sh` (Linux/macOS) or `build.ps1` (Windows) for cross-platform builds.

### Test
- **Run Tests**: `go test -v ./...`
- **Pattern**: Tests often mock shell commands. See `internal/cmd/verify_test.go` for examples using `MockRunner`.

### Release
1. Bump version in `extension.yaml` and `internal/cmd/version.go`.
2. Update `changelog.md`.
3. Tag and push.
4. After artifacts are generated, update `registry.json` checksums using:
   ```bash
   python3 scripts/update_registry.py <checksums_file> registry.json <version>
   ```

## Project Conventions

### Coding Standards (Go)
- **Error Handling**: Always wrap errors using `fmt.Errorf("context: %w", err)` to preserve the stack/chain.
- **Context Usage**: Pass `context.Context` as the first argument to functions performing I/O or long-running tasks.
- **Configuration**: Use the strongly-typed `AzureYaml` struct in `internal/checks/project.go` when working with `azure.yaml`. Do not parse YAML manually.
- **Project Layout**: Follow standard Go project layout.
  - `internal/cmd`: CLI command implementations.
  - `internal/checks`: Core logic and domain types.

### Command Execution & Mocking
- **Do not use `exec.Command` directly.** Use the `checks.CommandRunner` interface.
- **Mocking**: In tests, replace `checks.CommandRunner` with a `MockRunner` to simulate command output and errors without running actual shell commands.
  ```go
  // Example Mock
  checks.CommandRunner = &MockRunner{
      OutputFunc: func(name string, args ...string) ([]byte, error) {
          return []byte("1.0.0"), nil
      },
  }
  ```

### Verification Logic
- **Bypass Mechanism**: The `verify` command must respect `AZD_DOCTOR_SKIP_VERIFY` (global or command-specific).
- **Hook Context**: Use `AZD_HOOK_NAME` to infer the target command (e.g., `predeploy` -> `deploy`) if not explicitly provided.
- **DRY Principles**:
  - Use `checks.CheckTool` for standard CLI version checks.
  - Use `requireCheck` helper in `verify.go` to process `CheckResult` structs consistently.

### Testing Patterns
- **Coverage**: Ensure unit tests are created and maintained for all functions, especially new logic.
- **Framework**: Use `testing` package with `github.com/stretchr/testify/assert`.
- **Structure**: Prefer table-driven tests for covering multiple scenarios (e.g., different env vars or inputs).
- **Isolation**: Tests in `internal/cmd` should not execute real external commands; always mock.

### Documentation & Maintenance
- **README**: Keep `README.md` and installation instructions up-to-date with any feature changes or new commands.
- **Changelog**: Update `changelog.md` for every release or significant change.

## Integration Points
- **`azd`**: Invokes this binary. The extension expects `AZD_SERVER` env var when running in extension host mode.
- **`azure.yaml`**: Source of truth for project requirements.
