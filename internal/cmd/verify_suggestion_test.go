package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"spboyer.azd.doctor/internal/checks"
)

func TestRunVerify_DockerSuggestion(t *testing.T) {
	// Save original runner
	origRunner := checks.CommandRunner
	defer func() { checks.CommandRunner = origRunner }()

	// Mock runner
	mockRunner := &MockRunner{
		OutputFunc: func(name string, args ...string) ([]byte, error) {
			switch name {
			case "azd":
				return []byte("1.0.0"), nil
			case "git":
				return []byte("2.0.0"), nil
			case "gh":
				return []byte("2.0.0"), nil
			case "node":
				return []byte("18.0.0"), nil
			case "docker":
				return nil, fmt.Errorf("command not found")
			default:
				return nil, fmt.Errorf("unknown command: %s", name)
			}
		},
		RunFunc: func(name string, args ...string) error {
			return nil
		},
	}
	checks.CommandRunner = mockRunner

	// Create temp directory for project
	tmpDir, err := os.MkdirTemp("", "azd-doctor-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create azure.yaml
	content := `
name: test-project
services:
  api:
    language: js
    host: containerapp
    project: ./src/api
`
	err = os.WriteFile(filepath.Join(tmpDir, "azure.yaml"), []byte(content), 0644)
	require.NoError(t, err)

	// Change to temp dir
	cwd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(cwd)
	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Run Verify
	err = RunVerify(context.Background(), "up", 1*time.Second)
	
	// Assert error contains tip
	require.Error(t, err)
	assert.Contains(t, err.Error(), "required tool not found: docker")
	assert.Contains(t, err.Error(), "Tip: You can enable remote build in azure.yaml")
	assert.Contains(t, err.Error(), "remoteBuild: true")
	assert.Contains(t, err.Error(), "azd doctor configure remote-build")
}
