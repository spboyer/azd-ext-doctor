package cmd

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"spboyer.azd.doctor/internal/checks"
)

// MockRunner implements checks.Runner for testing
type MockRunner struct {
	OutputFunc func(name string, args ...string) ([]byte, error)
	RunFunc    func(name string, args ...string) error
}

func (m *MockRunner) Output(name string, args ...string) ([]byte, error) {
	if m.OutputFunc != nil {
		return m.OutputFunc(name, args...)
	}
	return nil, nil
}

func (m *MockRunner) Run(name string, args ...string) error {
	if m.RunFunc != nil {
		return m.RunFunc(name, args...)
	}
	return nil
}

func TestRunVerify_Skip(t *testing.T) {
	// Save original runner
	origRunner := checks.CommandRunner
	defer func() { checks.CommandRunner = origRunner }()

	// Set up a runner that always fails.
	// If verification is NOT skipped, RunVerify should fail because it calls checks (starting with azd version check).
	// If verification IS skipped, RunVerify should succeed (return nil) without calling checks.
	failingRunner := &MockRunner{
		OutputFunc: func(name string, args ...string) ([]byte, error) {
			return nil, fmt.Errorf("command execution failed intentionally")
		},
	}
	checks.CommandRunner = failingRunner

	tests := []struct {
		name       string
		envVar     string
		targetCmd  string
		shouldSkip bool
	}{
		{
			name:       "No Skip",
			envVar:     "",
			targetCmd:  "up",
			shouldSkip: false,
		},
		{
			name:       "Skip All (true)",
			envVar:     "true",
			targetCmd:  "up",
			shouldSkip: true,
		},
		{
			name:       "Skip All (1)",
			envVar:     "1",
			targetCmd:  "deploy",
			shouldSkip: true,
		},
		{
			name:       "Skip All (all)",
			envVar:     "all",
			targetCmd:  "provision",
			shouldSkip: true,
		},
		{
			name:       "Skip Specific (deploy)",
			envVar:     "deploy",
			targetCmd:  "deploy",
			shouldSkip: true,
		},
		{
			name:       "Skip Specific (deploy) but running provision",
			envVar:     "deploy",
			targetCmd:  "provision",
			shouldSkip: false,
		},
		{
			name:       "Skip Multiple (deploy,provision)",
			envVar:     "deploy,provision",
			targetCmd:  "provision",
			shouldSkip: true,
		},
		{
			name:       "Skip Multiple with spaces (deploy, provision)",
			envVar:     "deploy, provision",
			targetCmd:  "provision",
			shouldSkip: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("AZD_DOCTOR_SKIP_VERIFY", tt.envVar)
			defer os.Unsetenv("AZD_DOCTOR_SKIP_VERIFY")

			err := RunVerify(context.Background(), tt.targetCmd, 1*time.Second)
			if tt.shouldSkip {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				// The error comes from requireCheck -> CheckAzdVersion -> CheckTool -> CommandRunner.Output
				// CheckTool returns Installed: false, Error: err
				// requireCheck returns "required tool not found: azd"
				assert.Contains(t, err.Error(), "required tool not found")
			}
		})
	}
}

func TestRunVerify_HookContext(t *testing.T) {
	// Save original runner
	origRunner := checks.CommandRunner
	defer func() { checks.CommandRunner = origRunner }()

	failingRunner := &MockRunner{
		OutputFunc: func(name string, args ...string) ([]byte, error) {
			return nil, fmt.Errorf("command execution failed intentionally")
		},
	}
	checks.CommandRunner = failingRunner

	tests := []struct {
		name       string
		hookName   string
		skipVar    string
		shouldSkip bool
	}{
		{
			name:       "Predeploy Hook - Skip Deploy",
			hookName:   "predeploy",
			skipVar:    "deploy",
			shouldSkip: true,
		},
		{
			name:       "Preprovision Hook - Skip Provision",
			hookName:   "preprovision",
			skipVar:    "provision",
			shouldSkip: true,
		},
		{
			name:       "Preup Hook - Skip Up",
			hookName:   "preup",
			skipVar:    "up",
			shouldSkip: true,
		},
		{
			name:       "Predeploy Hook - Skip Provision (Mismatch)",
			hookName:   "predeploy",
			skipVar:    "provision",
			shouldSkip: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("AZD_HOOK_NAME", tt.hookName)
			os.Setenv("AZD_DOCTOR_SKIP_VERIFY", tt.skipVar)
			defer os.Unsetenv("AZD_HOOK_NAME")
			defer os.Unsetenv("AZD_DOCTOR_SKIP_VERIFY")

			// Pass empty targetCommand so it infers from hook
			err := RunVerify(context.Background(), "", 1*time.Second)
			if tt.shouldSkip {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
