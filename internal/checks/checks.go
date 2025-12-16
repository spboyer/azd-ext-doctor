package checks

import (
	"fmt"
	"os/exec"
	"strings"
)

type CheckResult struct {
	Name      string
	Installed bool
	Version   string
	Running   bool
	HasDaemon bool
	Error     error
}

func CheckTool(name string, args ...string) CheckResult {
	cmd := exec.Command(name, args...)
	out, err := cmd.Output()
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
		if err := exec.Command("docker", "info").Run(); err != nil {
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
		if err := exec.Command("podman", "info").Run(); err != nil {
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
