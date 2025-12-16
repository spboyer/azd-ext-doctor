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
	Error     error
}

func CheckTool(name string, args ...string) CheckResult {
	cmd := exec.Command(name, args...)
	out, err := cmd.Output()
	if err != nil {
		return CheckResult{Name: name, Installed: false, Error: err}
	}
	return CheckResult{Name: name, Installed: true, Version: strings.TrimSpace(string(out))}
}

func CheckDocker() CheckResult {
	// Check for docker
	res := CheckTool("docker", "--version")
	if res.Installed {
		return res
	}
	// Fallback to podman
	res = CheckTool("podman", "--version")
	if res.Installed {
		res.Name = "podman" // Update name to reflect what was found
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
