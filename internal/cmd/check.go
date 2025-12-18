package cmd

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/azure/azure-dev/cli/azd/pkg/azdext"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"spboyer.azd.doctor/internal/checks"
)

func NewCheckCommand() *cobra.Command {
	var skipAuth bool
	var authTimeout time.Duration

	checkCmd := &cobra.Command{
		Use:   "check",
		Short: "Run the doctor checks",
		RunE: func(cmd *cobra.Command, args []string) error {
			printRunning("Doctor Checks", "Starting...")

			// 1) Determine Project File
			projectFile := "azure.yaml"
			if _, err := os.Stat(projectFile); os.IsNotExist(err) {
				projectFile = "azure.yml"
				if _, err := os.Stat(projectFile); os.IsNotExist(err) {
					fmt.Println()
					printRunning("AZD Checks", "Checking tools")
					printResult(checks.CheckAzdVersion())
					printResult(checks.CheckGit())
					printResult(checks.CheckGh())

					fmt.Println()
					printRunning("Project Checks", "Checking azd project")
					printInfo("Project File", "Not found (azure.yaml/azure.yml)")
					printInfo("Project Name", "Unknown")

					fmt.Println()
					printRunning("Generic Checks", "Checking common dependencies")
					printResult(checks.CheckDocker())
					printResult(checks.CheckNode())
					printResult(checks.CheckPython())
					printResult(checks.CheckDotNet())
					printResult(checks.CheckBash())
					printResult(checks.CheckPwsh())
					printResult(checks.CheckAzureFunctionsCoreTools())

					fmt.Println()
					if skipAuth {
						printInfo("Azd Auth", "Skipped")
						return nil
					}
					printRunning("Azd Auth", "Checking login status")
					authCtx, cancel := context.WithTimeout(cmd.Context(), authTimeout)
					defer cancel()
					printResult(checks.CheckAzdLogin(authCtx, nil))
					return nil
				}
			}

			// 2) Load Project
			config, err := checks.LoadProjectConfig(projectFile)
			if err != nil {
				return fmt.Errorf("failed to load project config: %w", err)
			}

			// Initialize azd client only when we have a project file.
			ctx := azdext.WithAccessToken(cmd.Context())
			azdClient, err := azdext.NewAzdClient()
			if err != nil {
				return fmt.Errorf("failed to create azd client: %w", err)
			}
			defer azdClient.Close()

			// 3) AZD / Tool Checks
			fmt.Println()
			printRunning("AZD Checks", "Checking tools")
			printResult(checks.CheckAzdVersion())
			printResult(checks.CheckGit())
			printResult(checks.CheckGh())

			// 4) Project Checks
			fmt.Println()
			printRunning("Project Checks", "Checking azd project")
			printSuccess("Project File", projectFile)
			printInfo("Project Name", config.Name)
			printResult(checks.CheckAzdInit(ctx, azdClient))

			// 5) Project Hooks
			if len(config.Hooks) > 0 {
				fmt.Println()
				printRunning("Project Hooks", "Checking requirements")
				checkHooks(config.Hooks)
			}

			// 6) Required Extensions
			if len(config.RequiredVersions.Extensions) > 0 {
				fmt.Println()
				printRunning("Extensions", "Checking requirements")
				for name, version := range config.RequiredVersions.Extensions {
					printResult(checks.CheckExtension(name, version))
				}
			}

			// 7) Infra Checks
			fmt.Println()
			printRunning("Infra", "Checking requirements")
			provider := config.Infra.Provider
			if provider == "terraform" {
				printResult(checks.CheckTerraform())
			} else {
				// Default provider is bicep
				if provider == "" {
					provider = "bicep"
				}
				printInfo("Provider", provider)
			}

			// 8) Services
			checkedLangs := make(map[string]bool)
			checkedTools := make(map[string]bool)

			for name, svc := range config.Services {
				fmt.Println()
				printRunning("Service", fmt.Sprintf("%s (%s, %s)", name, svc.Host, svc.Language))

				// Check Language Requirements
				if svc.Language != "" && !checkedLangs[svc.Language] {
					switch svc.Language {
					case "js", "ts":
						printResult(checks.CheckNode())
					case "py", "python":
						printResult(checks.CheckPython())
					case "csharp", "fsharp", "dotnet":
						printResult(checks.CheckDotNet())
					default:
						// Unknown language - no checks.
					}
					checkedLangs[svc.Language] = true
				}

				// Check Hooks (Service Level)
				if len(svc.Hooks) > 0 {
					checkHooks(svc.Hooks)
				}

				// Check Container Requirements
				isContainerHost := svc.Host == "containerapp" || svc.Host == "aks"
				needsBuild := svc.Image == "" // If image is provided, assume pre-built.

				if isContainerHost && !svc.Docker.Remote && needsBuild {
					if !checkedTools["docker"] {
						printResult(checks.CheckDocker())
						checkedTools["docker"] = true
					}
				}

				// Check Azure Functions
				if svc.Host == "function" {
					if !checkedTools["func"] {
						printResult(checks.CheckAzureFunctionsCoreTools())
						checkedTools["func"] = true
					}
				}

				// Check Static Web Apps
				if svc.Host == "staticwebapp" {
					if !checkedTools["swa"] {
						printResult(checks.CheckSwaCli())
						checkedTools["swa"] = true
					}
				}
			}

			// 9) Azd Auth (separate + optional + timeout)
			fmt.Println()
			if skipAuth {
				printInfo("Azd Auth", "Skipped")
				return nil
			}
			printRunning("Azd Auth", "Checking login status")
			authCtx, cancel := context.WithTimeout(cmd.Context(), authTimeout)
			defer cancel()
			printResult(checks.CheckAzdLogin(authCtx, azdClient))

			return nil
		},
	}

	checkCmd.Flags().BoolVar(&skipAuth, "skip-auth", false, "Skip azd auth status check")
	checkCmd.Flags().DurationVar(&authTimeout, "auth-timeout", 5*time.Second, "Timeout for azd auth status check")

	return checkCmd
}

func runGenericChecks(ctx context.Context, azdClient checks.IAzdClient) {
	printResult(checks.CheckAzdVersion())
	printResult(checks.CheckAzdLogin(ctx, azdClient))
	printResult(checks.CheckGit())
	printResult(checks.CheckGh())
	printResult(checks.CheckDocker())
	printResult(checks.CheckNode())
	printResult(checks.CheckPython())
	printResult(checks.CheckDotNet())
	printResult(checks.CheckBash())
	printResult(checks.CheckPwsh())
	printResult(checks.CheckAzureFunctionsCoreTools())
}

func checkHooks(hooks checks.Hooks) {
	checkedShells := make(map[string]bool)

	for _, hookConfig := range hooks {
		shell := hookConfig.Shell
		if shell == "" {
			// Default based on OS
			if runtime.GOOS == "windows" {
				shell = "pwsh" // or powershell
			} else {
				shell = "sh" // or bash
			}
		}

		if !checkedShells[shell] {
			// fmt.Printf("  Hook '%s' requires shell: %s\n", hookName, shell)
			switch shell {
			case "sh", "bash":
				printResult(checks.CheckBash())
			case "pwsh", "powershell":
				printResult(checks.CheckPwsh())
			default:
				printInfo("Unknown Shell", shell)
			}
			checkedShells[shell] = true
		}
	}
}

func printResult(res checks.CheckResult) {
	if res.Installed {
		printSuccess(res.Name, res.Version)
		if res.HasDaemon {
			if res.Running {
				printSuccess(fmt.Sprintf("%s Daemon", res.Name), "Running")
			} else {
				printFailure(fmt.Sprintf("%s Daemon", res.Name), "Not running")
			}
		}
	} else {
		printFailure(res.Name, "Not found")
	}
}

// Styling helpers matching azd x builder
// Format: (SYMBOL) STATUS  MESSAGE  (DETAILS)

func printSuccess(message, details string) {
	fmt.Printf("%s %s  %-20s  %s\n",
		color.GreenString("(✓)"),
		color.GreenString("Done   "),
		message,
		color.HiBlackString("(%s)", details))
}

func printFailure(message, details string) {
	fmt.Printf("%s %s  %-20s  %s\n",
		color.RedString("(x)"),
		color.RedString("Error  "),
		message,
		color.HiBlackString("(%s)", details))
}

func printRunning(message, details string) {
	fmt.Printf("%s %s  %-20s  %s\n",
		color.CyanString("(-)"),
		color.CyanString("Running"),
		message,
		color.HiBlackString("(%s)", details))
}

func printInfo(message, details string) {
	fmt.Printf("%s %s  %-20s  %s\n",
		color.WhiteString("(i)"),
		color.WhiteString("Info   "),
		message,
		color.HiBlackString("(%s)", details))
}
