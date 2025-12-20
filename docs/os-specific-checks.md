# OS-Specific Tool Detection

## Overview

The `azd-ext-doctor` extension now includes intelligent, OS-aware tool detection that adapts to the operating system's conventions and common tool installations.

## Implementation Details

### Docker/Podman Detection

**Function**: `CheckDockerWithOS(goos string)`

| OS | Primary Check | Secondary Check | Notes |
|---|---|---|---|
| macOS | Docker | Podman | Docker Desktop is standard |
| Windows | Docker | Podman | Docker Desktop is standard |
| Linux | Docker | Podman | Both are common; Docker still more prevalent |

**Daemon Status**: Both tools are checked for daemon/service availability using `docker info` or `podman info`.

### Python Detection

**Function**: `CheckPythonWithOS(goos string)`

| OS | Primary Command | Secondary Command | Rationale |
|---|---|---|---|
| Windows | `python` | `python3` | Microsoft Store and standard installers use `python` |
| macOS | `python3` | `python` | Avoids Python 2.x; Python 3 is standard |
| Linux | `python3` | `python` | Avoids Python 2.x; Python 3 is standard |

### PowerShell Detection

**Function**: `CheckPwshWithOS(goos string)`

| OS | Primary Command | Secondary Command | Notes |
|---|---|---|---|
| Windows | `pwsh` | `powershell` | Checks PowerShell 7+ first, falls back to 5.1 |
| macOS | `pwsh` | None | Only PowerShell Core available |
| Linux | `pwsh` | None | Only PowerShell Core available |

### Bash Detection

**Function**: `CheckBashWithOS(goos string)`

| OS | Availability | Notes |
|---|---|---|
| macOS | Standard | Built-in shell |
| Linux | Standard | Built-in shell |
| Windows | Optional | Git Bash, WSL, or Cygwin |

## Testing

Each OS-specific function has comprehensive test coverage:

- `TestCheckDockerWithOS`: Tests Docker/Podman detection across all platforms
- `TestCheckPythonWithOS`: Tests Python command detection for each OS
- `TestCheckPwshWithOS`: Tests PowerShell detection including fallbacks
- `TestCheckBashWithOS`: Tests bash availability checks

### Test Structure

```go
func TestCheckDockerWithOS(t *testing.T) {
    t.Run("Docker on macOS", ...)
    t.Run("Podman on Linux", ...)
    t.Run("Docker daemon not running on Windows", ...)
    t.Run("Neither docker nor podman found", ...)
}
```

## Backward Compatibility

The public API remains unchanged. Existing functions like `CheckDocker()`, `CheckPython()`, etc. now internally call the OS-aware versions using `runtime.GOOS`:

```go
func CheckDocker() CheckResult {
    return CheckDockerWithOS(runtime.GOOS)
}
```

This ensures:
- Zero breaking changes for existing code
- OS detection happens automatically
- Test code can override OS via `CheckDockerWithOS("linux")`

## Benefits

1. **Correct Tool Detection**: Checks the right command for each platform
2. **Better User Experience**: More accurate results, fewer false negatives
3. **Flexible Testing**: Can test OS-specific behavior without changing the OS
4. **Future-Proof**: Easy to add new OS-specific logic as needed

## Example Output

### macOS with Docker Desktop
```
(✓) Done     docker               (Docker version 24.0.0)
(✓) Done     docker Daemon        (Running)
```

### Linux with Podman
```
(✓) Done     podman               (podman version 4.5.0)
(✓) Done     podman Daemon        (Running)
```

### Windows with Python from Microsoft Store
```
(✓) Done     python               (Python 3.11.0)
```

### macOS/Linux with Python 3
```
(✓) Done     python3              (Python 3.11.0)
```
