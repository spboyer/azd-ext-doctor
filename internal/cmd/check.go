package cmd

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/azure/azure-dev/cli/azd/pkg/azdext"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"spboyer.azd.doctor/internal/checks"
)

func NewCheckCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "check",
		Short: "Run the doctor checks",
		RunE: func(cmd *cobra.Command, args []string) error {
			printRunning("Doctor Checks", "Starting...")

			// Initialize azd client
			ctx := azdext.WithAccessToken(cmd.Context())
			azdClient, err := azdext.NewAzdClient()
			if err != nil {
				return fmt.Errorf("failed to create azd client: %w", err)
			}
			defer azdClient.Close()

			// 1. Determine Project File
			projectFile := "azure.yaml"
			if _, err := os.Stat(projectFile); os.IsNotExist(err) {
				projectFile = "azure.yml"
				if _, err := os.Stat(projectFile); os.IsNotExist(err) {
					printInfo("Project File", "Not found")
					printRunning("Generic Checks", "Running...")
					runGenericChecks(ctx, azdClient)
					return nil
				}
			}

			// 2. Load Project
			printSuccess("Project File", projectFile)
			config, err := checks.LoadProjectConfig(projectFile)
			if err != nil {
				return fmt.Errorf("failed to load project config: %w", err)
			}

			printInfo("Project Name", config.Name)

			// AZD Checks
			fmt.Println()
			printRunning("AZD Checks", "Checking azd environment")
			printResult(checks.CheckAzdVersion())
			printResult(checks.CheckAzdLogin(ctx, azdClient))
			printResult(checks.CheckAzdInit(ctx, azdClient))
			printResult(checks.CheckGit())
			printResult(checks.CheckGh())

			// 3. Check Hooks (Root)
			if len(config.Hooks) > 0 {
				fmt.Println()
				printRunning("Project Hooks", "Checking requirements")
				checkHooks(config.Hooks)
			}

			// 4. Check Services
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
						// fmt.Printf("  [?] Unknown language requirement for: %s\n", svc.Language)
					}
					checkedLangs[svc.Language] = true
				}

				// Check Hooks (Service Level)
				if len(svc.Hooks) > 0 {
					checkHooks(svc.Hooks)
				}

				// Check Container Requirements
				// If host is container-based AND not a remote build (implied by lack of 'remote' flag or presence of 'image' without build)
				// Logic:
				// - If Host is containerapp or aks
				// - AND Docker.Remote is false
				// - AND Image is empty (implies we are building from source)
				isContainerHost := svc.Host == "containerapp" || svc.Host == "aks"
				needsBuild := svc.Image == "" // If image is provided, we assume pre-built

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
			}

			return nil
		},
	}
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
