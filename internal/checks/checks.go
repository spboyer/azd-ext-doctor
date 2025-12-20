package checks

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/blang/semver/v4"
)

// Runner interface allows mocking of exec.Command
type Runner interface {
	Output(name string, args ...string) ([]byte, error)
	Run(name string, args ...string) error
}

// RunnerWithContext is an optional extension to Runner that allows command execution
// to respect cancellation and timeouts.
type RunnerWithContext interface {
	OutputContext(ctx context.Context, name string, args ...string) ([]byte, error)
	RunContext(ctx context.Context, name string, args ...string) error
}

type RealRunner struct{}

func (r *RealRunner) Output(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).Output()
}

func (r *RealRunner) OutputContext(ctx context.Context, name string, args ...string) ([]byte, error) {
	return exec.CommandContext(ctx, name, args...).Output()
}

func (r *RealRunner) Run(name string, args ...string) error {
	return exec.Command(name, args...).Run()
}

func (r *RealRunner) RunContext(ctx context.Context, name string, args ...string) error {
	return exec.CommandContext(ctx, name, args...).Run()
}

var CommandRunner Runner = &RealRunner{}

func runnerOutput(ctx context.Context, name string, args ...string) ([]byte, error) {
	if r, ok := CommandRunner.(RunnerWithContext); ok {
		return r.OutputContext(ctx, name, args...)
	}

	return CommandRunner.Output(name, args...)
}

func runnerRun(ctx context.Context, name string, args ...string) error {
	if r, ok := CommandRunner.(RunnerWithContext); ok {
		return r.RunContext(ctx, name, args...)
	}

	return CommandRunner.Run(name, args...)
}

type CheckResult struct {
	Name      string
	Installed bool
	Version   string
	Running   bool
	HasDaemon bool
	Error     error
}

func CheckTool(name string, args ...string) CheckResult {
	out, err := CommandRunner.Output(name, args...)
	if err != nil {
		return CheckResult{Name: name, Installed: false, Error: err}
	}
	return CheckResult{Name: name, Installed: true, Version: strings.TrimSpace(string(out)), Running: true}
}

func CheckDocker() CheckResult {
	return CheckDockerWithOS(runtime.GOOS)
}

// CheckDockerWithOS checks for Docker or Podman with OS-specific priority
func CheckDockerWithOS(goos string) CheckResult {
	// Determine check order based on OS
	// Linux: Podman is increasingly common, check both
	// macOS/Windows: Docker Desktop is standard
	var primaryTool, secondaryTool string
	var primaryCmd, secondaryCmd string

	switch goos {
	case "linux":
		// On Linux, try Docker first (still most common), then Podman
		primaryTool = "docker"
		primaryCmd = "docker"
		secondaryTool = "podman"
		secondaryCmd = "podman"
	case "darwin", "windows":
		// On macOS and Windows, Docker Desktop is standard
		primaryTool = "docker"
		primaryCmd = "docker"
		secondaryTool = "podman"
		secondaryCmd = "podman"
	default:
		// Unknown OS, try docker first
		primaryTool = "docker"
		primaryCmd = "docker"
		secondaryTool = "podman"
		secondaryCmd = "podman"
	}

	// Try primary tool
	res := CheckTool(primaryCmd, "--version")
	if res.Installed {
		res.Name = primaryTool
		res.HasDaemon = true
		// Check if daemon is running
		if err := CommandRunner.Run(primaryCmd, "info"); err != nil {
			res.Running = false
			res.Error = fmt.Errorf("daemon not running")
		}
		return res
	}

	// Fallback to secondary tool
	res = CheckTool(secondaryCmd, "--version")
	if res.Installed {
		res.Name = secondaryTool
		res.HasDaemon = true
		// Check if daemon is running
		// Note: Podman on Linux often runs rootless/daemonless
		if err := CommandRunner.Run(secondaryCmd, "info"); err != nil {
			res.Running = false
			res.Error = fmt.Errorf("not running or not configured")
		}
		return res
	}

	return CheckResult{
		Name:      "docker/podman",
		Installed: false,
		Error:     fmt.Errorf("neither docker nor podman found"),
	}
}

func CheckNode() CheckResult {
	return CheckTool("node", "--version")
}

func CheckPython() CheckResult {
	return CheckPythonWithOS(runtime.GOOS)
}

// CheckPythonWithOS checks for Python with OS-specific command priority
func CheckPythonWithOS(goos string) CheckResult {
	var primaryCmd, secondaryCmd string

	switch goos {
	case "windows":
		// On Windows, 'python' is more common (from Microsoft Store or installer)
		primaryCmd = "python"
		secondaryCmd = "python3"
	case "darwin", "linux":
		// On macOS and Linux, 'python3' is standard to avoid Python 2.x
		primaryCmd = "python3"
		secondaryCmd = "python"
	default:
		// Unknown OS, try python3 first
		primaryCmd = "python3"
		secondaryCmd = "python"
	}

	res := CheckTool(primaryCmd, "--version")
	if !res.Installed {
		res = CheckTool(secondaryCmd, "--version")
		if res.Installed {
			// Update name to show which command worked
			res.Name = secondaryCmd
		}
	} else {
		res.Name = primaryCmd
	}

	if !res.Installed {
		res.Name = "python"
	}

	return res
}

func CheckDotNet() CheckResult {
	return CheckTool("dotnet", "--version")
}

func CheckBash() CheckResult {
	return CheckBashWithOS(runtime.GOOS)
}

// CheckBashWithOS checks for Bash with OS-specific expectations
func CheckBashWithOS(goos string) CheckResult {
	switch goos {
	case "windows":
		// On Windows, bash might be from Git Bash, WSL, or Cygwin
		// It's optional but common for development
		res := CheckTool("bash", "--version")
		res.Name = "bash"
		return res
	case "darwin", "linux":
		// On macOS and Linux, bash is standard
		res := CheckTool("bash", "--version")
		res.Name = "bash"
		return res
	default:
		res := CheckTool("bash", "--version")
		res.Name = "bash"
		return res
	}
}

func CheckPwsh() CheckResult {
	return CheckPwshWithOS(runtime.GOOS)
}

// CheckPwshWithOS checks for PowerShell with OS-specific logic
func CheckPwshWithOS(goos string) CheckResult {
	var primaryCmd, secondaryCmd string

	switch goos {
	case "windows":
		// On Windows, try pwsh (PowerShell 7+) first, then fall back to powershell (5.1)
		primaryCmd = "pwsh"
		secondaryCmd = "powershell"
	case "darwin", "linux":
		// On macOS and Linux, only pwsh is available (PowerShell Core)
		primaryCmd = "pwsh"
		secondaryCmd = "" // No fallback
	default:
		// Unknown OS
		primaryCmd = "pwsh"
		secondaryCmd = "powershell"
	}

	res := CheckTool(primaryCmd, "--version")
	if !res.Installed && secondaryCmd != "" {
		res = CheckTool(secondaryCmd, "--version")
		if res.Installed {
			res.Name = secondaryCmd
		}
	} else if res.Installed {
		res.Name = primaryCmd
	}

	if !res.Installed {
		res.Name = "pwsh/powershell"
	}

	return res
}

func CheckAzureFunctionsCoreTools() CheckResult {
	return CheckTool("func", "--version")
}

func CheckSwaCli() CheckResult {
	return CheckTool("swa", "--version")
}

func CheckGit() CheckResult {
	return CheckTool("git", "--version")
}

func CheckGh() CheckResult {
	return CheckTool("gh", "--version")
}

func CheckTerraform() CheckResult {
	return CheckTool("terraform", "--version")
}

type AzdExtension struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

func GetInstalledExtensions() ([]AzdExtension, error) {
	out, err := CommandRunner.Output("azd", "extension", "list", "--installed", "--output", "json")
	if err != nil {
		return nil, err
	}
	var extensions []AzdExtension
	if err := json.Unmarshal(out, &extensions); err != nil {
		return nil, fmt.Errorf("failed to parse azd extension list: %w", err)
	}
	return extensions, nil
}

func CheckExtension(extensions []AzdExtension, id, requiredRange string) CheckResult {
	var installedVer string
	found := false
	for _, ext := range extensions {
		if ext.Id == id {
			installedVer = ext.Version
			found = true
			break
		}
	}

	if !found {
		return CheckResult{Name: "extension " + id, Installed: false, Error: fmt.Errorf("extension not installed")}
	}

	if requiredRange != "" {
		v, err := semver.Parse(installedVer)
		if err != nil {
			return CheckResult{Name: "extension " + id, Installed: true, Version: installedVer, Error: fmt.Errorf("invalid installed version format: %w", err)}
		}
		expectedRange, err := semver.ParseRange(requiredRange)
		if err != nil {
			return CheckResult{Name: "extension " + id, Installed: true, Version: installedVer, Error: fmt.Errorf("invalid required version range: %w", err)}
		}
		if !expectedRange(v) {
			return CheckResult{Name: "extension " + id, Installed: true, Version: installedVer, Error: fmt.Errorf("version %s does not satisfy range %s", installedVer, requiredRange)}
		}
	}

	return CheckResult{Name: "extension " + id, Installed: true, Version: installedVer, Running: true}
}
