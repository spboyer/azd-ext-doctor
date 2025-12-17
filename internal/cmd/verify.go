package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/azure/azure-dev/cli/azd/pkg/azdext"
	"github.com/spf13/cobra"
	"spboyer.azd.doctor/internal/checks"
)

func NewVerifyCommand() *cobra.Command {
	var targetCommand string
	var authTimeout time.Duration

	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify environment for a specific azd command (up, provision, deploy)",
		Long: `Verifies that the environment meets all requirements for running a specific azd command.

This command performs strict checks for:
- Required tools (azd, git)
- Authentication status (must be logged in)
- Project-specific requirements based on azure.yaml (languages, Docker, Functions Core Tools)

It is automatically invoked by azd before 'up', 'provision', and 'deploy' commands, but can also be run manually for debugging.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunVerify(cmd.Context(), targetCommand, authTimeout)
		},
	}

	cmd.Flags().StringVar(&targetCommand, "command", "up", "The azd command to verify for (up, provision, deploy)")
	cmd.Flags().DurationVar(&authTimeout, "auth-timeout", 5*time.Second, "Timeout for azd auth status check")

	return cmd
}

func RunVerify(ctx context.Context, targetCommand string, authTimeout time.Duration) error {
	// Check for bypass environment variable
	// AZD_DOCTOR_SKIP_VERIFY can be:
	// - "true", "1", "all": Skip all verification
	// - "up", "provision", "deploy": Skip verification for specific command
	// - Comma separated list: "provision,deploy"
	skipVerify := os.Getenv("AZD_DOCTOR_SKIP_VERIFY")
	if skipVerify != "" {
		// Check for global skip
		if skipVerify == "true" || skipVerify == "1" || skipVerify == "all" {
			printSuccess("Verification", "Skipped by AZD_DOCTOR_SKIP_VERIFY")
			return nil
		}
	}

	// Determine the actual command context if running as a hook
	// azd sets AZD_HOOK_NAME to the name of the hook (e.g. predeploy, preprovision)
	hookName := os.Getenv("AZD_HOOK_NAME")
	if hookName != "" {
		switch hookName {
		case "predeploy":
			targetCommand = "deploy"
		case "preprovision":
			targetCommand = "provision"
		case "preup":
			targetCommand = "up"
		}
	}

	// Default to 'up' if not specified, or validate input
	if targetCommand == "" {
		targetCommand = "up"
	}

	// Check if specific command should be skipped
	if skipVerify != "" {
		// Simple contains check for now, could be more robust with splitting
		if contains(skipVerify, targetCommand) {
			printSuccess("Verification", fmt.Sprintf("Skipped for %s by AZD_DOCTOR_SKIP_VERIFY", targetCommand))
			return nil
		}
	}

	if targetCommand != "up" && targetCommand != "provision" && targetCommand != "deploy" {
		return fmt.Errorf("invalid command target: %s. Must be one of: up, provision, deploy", targetCommand)
	}

	printRunning("Verifying for", targetCommand)

	// 1. Common Checks (azd, git, gh)
	if err := requireCheck(checks.CheckAzdVersion()); err != nil {
		return err
	}
	if err := requireCheck(checks.CheckGit()); err != nil {
		return err
	}
	// gh is often optional, but let's check it as warning or skip if not critical?
	// For now, let's treat it as non-critical for 'provision'/'deploy' unless we know we need it.
	// But 'azd pipeline' needs it. 'azd up' might not strictly need it unless using repo.
	// Let's just log it but not fail? Or maybe fail if it's missing?
	// The user said "specifically the ones that would log as errors or 'stops'".
	// 'gh' is not strictly required for local provision/deploy.
	// So we skip strict requirement for gh.

	// 2. Auth Check
	authCtx, cancel := context.WithTimeout(ctx, authTimeout)
	defer cancel()
	// We need an azd client for auth check if we want to be thorough, but CheckAzdLogin uses CLI mostly.
	// However, CheckAzdLogin signature is (ctx, client).
	// Let's try to create client.
	azdClient, err := azdext.NewAzdClient()
	if err != nil {
		// If we can't create client, we might still be able to check login via CLI?
		// But CheckAzdLogin takes client.
		// Let's proceed with what we have.
	}
	if azdClient != nil {
		defer azdClient.Close()
	}

	loginRes := checks.CheckAzdLogin(authCtx, azdClient)
	if !loginRes.Installed || loginRes.Error != nil {
		printFailure(loginRes.Name, "Not logged in or error")
		return fmt.Errorf("azd auth check failed: %v", loginRes.Error)
	}
	printSuccess(loginRes.Name, loginRes.Version)

	// 3. Project Checks
	projectFile := "azure.yaml"
	if _, err := os.Stat(projectFile); os.IsNotExist(err) {
		projectFile = "azure.yml"
	}

	if _, err := os.Stat(projectFile); os.IsNotExist(err) {
		return fmt.Errorf("project file (azure.yaml/yml) not found, required for %s", targetCommand)
	}

	config, err := checks.LoadProjectConfig(projectFile)
	if err != nil {
		return fmt.Errorf("failed to load project config: %w", err)
	}

	// Required Extensions Check
	if len(config.RequiredVersions.Extensions) > 0 {
		for name, version := range config.RequiredVersions.Extensions {
			if err := requireCheck(checks.CheckExtension(name, version)); err != nil {
				return err
			}
		}
	}

	// Infra Checks (provision/up)
	if targetCommand == "provision" || targetCommand == "up" {
		provider := config.Infra.Provider
		// Default to bicep if empty
		if provider == "" || provider == "bicep" {
			if err := requireCheck(checks.CheckBicep()); err != nil {
				return err
			}
		} else if provider == "terraform" {
			if err := requireCheck(checks.CheckTerraform()); err != nil {
				return err
			}
		}
	}

	// Service Checks (deploy/up)
	if targetCommand == "deploy" || targetCommand == "up" {
		// Check Services
		checkedLangs := make(map[string]bool)
		checkedTools := make(map[string]bool)

		for _, svc := range config.Services {
			// Language Checks
			if svc.Language != "" && !checkedLangs[svc.Language] {
				var res checks.CheckResult
				switch svc.Language {
				case "js", "ts":
					res = checks.CheckNode()
				case "py", "python":
					res = checks.CheckPython()
				case "csharp", "fsharp", "dotnet":
					res = checks.CheckDotNet()
				}

				if res.Name != "" {
					if err := requireCheck(res); err != nil {
						return err
					}
				}
				checkedLangs[svc.Language] = true
			}

			// Container Checks
			isContainerHost := svc.Host == "containerapp" || svc.Host == "aks"
			needsBuild := svc.Image == ""

			if isContainerHost && !svc.Docker.Remote && needsBuild {
				if !checkedTools["docker"] {
					if err := requireCheck(checks.CheckDocker()); err != nil {
						return err
					}
					checkedTools["docker"] = true
				}
			}

			// Functions Checks
			if svc.Host == "function" {
				if !checkedTools["func"] {
					if err := requireCheck(checks.CheckAzureFunctionsCoreTools()); err != nil {
						return err
					}
					checkedTools["func"] = true
				}
			}

			// Static Web Apps Checks
			if svc.Host == "staticwebapp" {
				if !checkedTools["swa"] {
					if err := requireCheck(checks.CheckSwaCli()); err != nil {
						return err
					}
					checkedTools["swa"] = true
				}
			}
		}
	}

	printSuccess("Verification", "Passed")
	return nil
}

func requireCheck(res checks.CheckResult) error {
	if !res.Installed {
		printFailure(res.Name, "Not found")
		return fmt.Errorf("required tool not found: %s", res.Name)
	}
	if res.Error != nil {
		printFailure(res.Name, res.Error.Error())
		return fmt.Errorf("check failed for %s: %w", res.Name, res.Error)
	}
	if res.HasDaemon && !res.Running {
		printFailure(res.Name, "Daemon not running")
		return fmt.Errorf("%s daemon is not running", res.Name)
	}
	printSuccess(res.Name, res.Version)
	return nil
}

func contains(s, substr string) bool {
	parts := strings.Split(s, ",")
	for _, p := range parts {
		if strings.TrimSpace(p) == substr {
			return true
		}
	}
	return false
}
