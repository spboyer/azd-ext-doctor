package checks

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
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
	// Check for docker
	res := CheckTool("docker", "--version")
	if res.Installed {
		res.HasDaemon = true
		// Check if docker daemon is running
		if err := CommandRunner.Run("docker", "info"); err != nil {
			res.Running = false
			res.Error = fmt.Errorf("daemon not running")
		}
		return res
	}
	// Fallback to podman
	res = CheckTool("podman", "--version")
	if res.Installed {
		res.Name = "podman" // Update name to reflect what was found
		res.HasDaemon = true
		// Check if podman is working (podman info)
		if err := CommandRunner.Run("podman", "info"); err != nil {
			res.Running = false
			res.Error = fmt.Errorf("daemon not running")
		}
		return res
	}
	return CheckResult{Name: "docker/podman", Installed: false, Error: fmt.Errorf("neither docker nor podman found")}
}

func CheckNode() CheckResult {
	return CheckTool("node", "--version")
}

func CheckPython() CheckResult {
	res := CheckTool("python3", "--version")
	if !res.Installed {
		return CheckTool("python", "--version")
	}
	return res
}

func CheckDotNet() CheckResult {
	return CheckTool("dotnet", "--version")
}

func CheckBash() CheckResult {
	return CheckTool("bash", "--version")
}

func CheckPwsh() CheckResult {
	res := CheckTool("pwsh", "--version")
	if !res.Installed {
		return CheckTool("powershell", "--version")
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

func CheckBicep() CheckResult {
	return CheckTool("bicep", "--version")
}

func CheckTerraform() CheckResult {
	return CheckTool("terraform", "--version")
}

type AzdExtension struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func GetInstalledExtensions() ([]AzdExtension, error) {
	out, err := CommandRunner.Output("azd", "extension", "list", "--output", "json")
	if err != nil {
		return nil, err
	}
	var extensions []AzdExtension
	if err := json.Unmarshal(out, &extensions); err != nil {
		return nil, fmt.Errorf("failed to parse azd extension list: %w", err)
	}
	return extensions, nil
}

func CheckExtension(name, requiredRange string) CheckResult {
	extensions, err := GetInstalledExtensions()
	if err != nil {
		return CheckResult{Name: "extension " + name, Installed: false, Error: fmt.Errorf("failed to list extensions: %w", err)}
	}

	var installedVer string
	found := false
	for _, ext := range extensions {
		if ext.Name == name {
			installedVer = ext.Version
			found = true
			break
		}
	}

	if !found {
		return CheckResult{Name: "extension " + name, Installed: false, Error: fmt.Errorf("extension not installed")}
	}

	if requiredRange != "" {
		v, err := semver.Parse(installedVer)
		if err != nil {
			return CheckResult{Name: "extension " + name, Installed: true, Version: installedVer, Error: fmt.Errorf("invalid installed version format: %w", err)}
		}
		expectedRange, err := semver.ParseRange(requiredRange)
		if err != nil {
			return CheckResult{Name: "extension " + name, Installed: true, Version: installedVer, Error: fmt.Errorf("invalid required version range: %w", err)}
		}
		if !expectedRange(v) {
			return CheckResult{Name: "extension " + name, Installed: true, Version: installedVer, Error: fmt.Errorf("version %s does not satisfy range %s", installedVer, requiredRange)}
		}
	}

	return CheckResult{Name: "extension " + name, Installed: true, Version: installedVer, Running: true}
}
