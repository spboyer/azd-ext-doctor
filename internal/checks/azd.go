package checks

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/azure/azure-dev/cli/azd/pkg/azdext"
)

// IAzdClient defines the interface for the azd extension client
type IAzdClient interface {
	Deployment() azdext.DeploymentServiceClient
	Project() azdext.ProjectServiceClient
}

type LoginStatus string

const (
	LoginStatusSuccess         LoginStatus = "success"
	LoginStatusUnauthenticated LoginStatus = "unauthenticated"
)

type LoginResult struct {
	Status    LoginStatus `json:"status"`
	ExpiresOn *time.Time  `json:"expiresOn,omitempty"`
}

func CheckAzdVersion() CheckResult {
	return CheckTool("azd", "version")
}

func CheckAzdLogin(ctx context.Context, client IAzdClient) CheckResult {
	// Check login status using CLI command
	out, err := CommandRunner.Output("azd", "auth", "login", "--check-status")
	outputStr := strings.TrimSpace(string(out))

	if err != nil {
		msg := "Not logged in"
		if len(outputStr) > 0 {
			msg = outputStr
		}
		return CheckResult{Name: "azd auth", Installed: true, Version: msg, Running: false, Error: fmt.Errorf("not logged in: %w", err)}
	}

	if outputStr == "" {
		outputStr = "Logged in"
	}

	return CheckResult{Name: "azd auth", Installed: true, Version: outputStr, Running: true}
}

func CheckAzdInit(ctx context.Context, client IAzdClient) CheckResult {
	// Check if project is initialized using extension client
	resp, err := client.Project().Get(ctx, &azdext.EmptyRequest{})
	if err != nil {
		return CheckResult{Name: "azd project", Installed: false, Error: fmt.Errorf("project not initialized: %w", err)}
	}

	if resp.Project == nil {
		return CheckResult{Name: "azd project", Installed: false, Error: fmt.Errorf("project not initialized (no project returned)")}
	}

	return CheckResult{Name: "azd project", Installed: true, Version: fmt.Sprintf("Initialized (%s)", resp.Project.Name), Running: true}
}
