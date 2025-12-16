package checks

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
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

func CheckGit() CheckResult {
	return CheckTool("git", "--version")
}

func CheckGh() CheckResult {
	return CheckTool("gh", "--version")
}
