package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"spboyer.azd.doctor/internal/checks"
)

// TestCheckCommand_NoProjectFile tests that the check command
// runs only non-project checks when there's no azure.yaml file
func TestCheckCommand_NoProjectFile(t *testing.T) {
	// Save original runner and restore after test
	origRunner := checks.CommandRunner
	defer func() { checks.CommandRunner = origRunner }()

	// Mock runner that returns success for all tools
	mockRunner := &MockRunner{
		OutputFunc: func(name string, args ...string) ([]byte, error) {
			switch name {
			case "azd":
				return []byte("azd version 1.0.0"), nil
			case "git":
				return []byte("git version 2.0.0"), nil
			case "gh":
				return []byte("gh version 2.0.0"), nil
			case "docker":
				return []byte("Docker version 20.0.0"), nil
			case "node":
				return []byte("v18.0.0"), nil
			case "python", "python3":
				return []byte("Python 3.10.0"), nil
			case "dotnet":
				return []byte("8.0.0"), nil
			case "bash":
				return []byte("GNU bash, version 5.0.0"), nil
			case "pwsh":
				return []byte("PowerShell 7.0.0"), nil
			case "func":
				return []byte("4.0.0"), nil
			default:
				return []byte("1.0.0"), nil
			}
		},
		RunFunc: func(name string, args ...string) error {
			// docker info, podman info - success
			return nil
		},
	}
	checks.CommandRunner = mockRunner

	// Create a temporary directory without azure.yaml
	tmpDir, err := os.MkdirTemp("", "test-no-project-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Change to temporary directory
	origDir, _ := os.Getwd()
	err = os.Chdir(tmpDir)
	assert.NoError(t, err)
	defer os.Chdir(origDir)

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run check command
	cmd := NewCheckCommand()
	cmd.SetArgs([]string{"--skip-auth"})
	err = cmd.Execute()

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Assertions
	assert.NoError(t, err, "check command should succeed")

	// Should include non-project checks
	assert.Contains(t, output, "AZD Checks", "Should run AZD checks")
	assert.Contains(t, output, "azd", "Should check azd")
	assert.Contains(t, output, "git", "Should check git")
	assert.Contains(t, output, "gh", "Should check gh")

	// Should indicate project file not found
	assert.Contains(t, output, "Project Checks", "Should run project checks")
	assert.Contains(t, output, "Not found (azure.yaml/azure.yml)", "Should indicate project file not found")

	// Should run generic checks
	assert.Contains(t, output, "Generic Checks", "Should run generic checks")
	assert.Contains(t, output, "docker", "Should check docker")
	assert.Contains(t, output, "node", "Should check node")
	assert.Contains(t, output, "python", "Should check python")

	// Should NOT run project-specific checks
	assert.NotContains(t, output, "Running Infra", "Should NOT run infra checks")
	assert.NotContains(t, output, "Running Service", "Should NOT run service checks")
}

// TestCheckCommand_WithProjectFile tests that the check command
// runs all checks including project-specific ones when azure.yaml exists
func TestCheckCommand_WithProjectFile(t *testing.T) {
	// Save original runner and restore after test
	origRunner := checks.CommandRunner
	defer func() { checks.CommandRunner = origRunner }()

	// Mock runner
	mockRunner := &MockRunner{
		OutputFunc: func(name string, args ...string) ([]byte, error) {
			switch name {
			case "azd":
				return []byte("azd version 1.0.0"), nil
			case "git":
				return []byte("git version 2.0.0"), nil
			case "gh":
				return []byte("gh version 2.0.0"), nil
			case "dotnet":
				return []byte("8.0.0"), nil
			case "docker":
				return []byte("Docker version 20.0.0"), nil
			default:
				return []byte("1.0.0"), nil
			}
		},
		RunFunc: func(name string, args ...string) error {
			return nil
		},
	}
	checks.CommandRunner = mockRunner

	// Create a temporary directory with azure.yaml
	tmpDir, err := os.MkdirTemp("", "test-with-project-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create azure.yaml
	azureYaml := `name: test-project
services:
  api:
    host: containerapp
    language: dotnet
infra:
  provider: bicep
`
	err = os.WriteFile(tmpDir+"/azure.yaml", []byte(azureYaml), 0644)
	assert.NoError(t, err)

	// Change to temporary directory
	origDir, _ := os.Getwd()
	err = os.Chdir(tmpDir)
	assert.NoError(t, err)
	defer os.Chdir(origDir)

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run check command
	cmd := NewCheckCommand()
	cmd.SetArgs([]string{"--skip-auth"})
	err = cmd.Execute()

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Assertions
	assert.NoError(t, err, "check command should succeed")

	// Should include AZD checks
	assert.Contains(t, output, "AZD Checks", "Should run AZD checks")

	// Should find project file
	assert.Contains(t, output, "Project Checks", "Should run project checks")
	assert.Contains(t, output, "azure.yaml", "Should find project file")
	assert.Contains(t, output, "test-project", "Should show project name")

	// Should run infra checks with completion status
	assert.Contains(t, output, "Infra", "Should run infra checks")
	// Should show provider info for bicep
	if strings.Contains(output, "bicep") {
		assert.Contains(t, output, "Provider", "Should show provider info")
	}

	// Should run service checks
	assert.Contains(t, output, "Service", "Should run service checks")
	assert.Contains(t, output, "api", "Should check api service")
	assert.Contains(t, output, "dotnet", "Should check dotnet for service")
}
