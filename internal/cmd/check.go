package cmd

import (
	"fmt"

	"github.com/Azure/azure-dev/cli/azd/pkg/azdext"
	"github.com/spf13/cobra"
	"spboyer.azd.doctor/internal/checks"
)

func NewCheckCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "check",
		Short: "Run the doctor checks",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := azdext.WithAccessToken(cmd.Context())
			azdClient, err := azdext.NewAzdClient()
			if err != nil {
				return fmt.Errorf("failed to create azd client: %w", err)
			}
			defer azdClient.Close()

			fmt.Println("Running doctor checks...")

			// Get Project Config
			projectResp, err := azdClient.Project().Get(ctx, &azdext.EmptyRequest{})
			if err != nil {
				// If we can't get the project, we might not be in an azd project, or azd isn't running the extension in a context where it can provide it.
				// But we can still check for generic tools.
				fmt.Printf("Warning: Could not retrieve project context: %v\n", err)
				fmt.Println("Checking generic tools...")
				printResult(checks.CheckDocker())
				printResult(checks.CheckNode())
				printResult(checks.CheckPython())
				printResult(checks.CheckDotNet())
				return nil
			}

			if projectResp.Project == nil {
				fmt.Println("No project context found. Are you in an initialized azd project?")
				return nil
			}

			fmt.Printf("Checking prerequisites for project: %s\n", projectResp.Project.Name)

			// Always check for container runtime
			printResult(checks.CheckDocker())

			// Analyze services and check requirements
			checkedLangs := make(map[string]bool)

			for name, svc := range projectResp.Project.Services {
				fmt.Printf("\nService '%s' (Host: %s, Language: %s)\n", name, svc.Host, svc.Language)

				if !checkedLangs[svc.Language] {
					switch svc.Language {
					case "js", "ts":
						printResult(checks.CheckNode())
					case "py", "python":
						printResult(checks.CheckPython())
					case "csharp", "fsharp", "dotnet":
						printResult(checks.CheckDotNet())
					default:
						fmt.Printf("  [?] Unknown language requirement for: %s\n", svc.Language)
					}
					checkedLangs[svc.Language] = true
				}
			}

			return nil
		},
	}
}

func printResult(res checks.CheckResult) {
	if res.Installed {
		fmt.Printf("  [✓] %s: %s\n", res.Name, res.Version)
	} else {
		fmt.Printf("  [✗] %s: Not found\n", res.Name)
	}
}
