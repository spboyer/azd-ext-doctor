package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/azure/azure-dev/cli/azd/pkg/azdext"
	"github.com/spf13/cobra"
)

func NewListenCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "listen",
		Short:  "Starts the extension in server mode",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunExtensionHost(cmd.Context())
		},
	}
	return cmd
}

func RunExtensionHost(ctx context.Context) error {
	// Create a new context that includes the AZD access token
	ctx = azdext.WithAccessToken(ctx)

	// Create a new AZD client
	client, err := azdext.NewAzdClient()
	if err != nil {
		return fmt.Errorf("failed to create azd client: %w", err)
	}
	defer client.Close()

	// Use the ExtensionHost pattern which handles Ready() signaling automatically
	host := azdext.NewExtensionHost(client).
		WithProjectEventHandler("preprovision", onPreProvision).
		WithProjectEventHandler("predeploy", onPreDeploy).
		WithProjectEventHandler("preup", onPreUp)

	// Start listening for events
	// This is a blocking call and will not return until the server connection is closed
	if err := host.Run(ctx); err != nil {
		return fmt.Errorf("failed to run extension: %w", err)
	}

	return nil
}

func onPreProvision(ctx context.Context, args *azdext.ProjectEventArgs) error {
	// Default timeout for auth check in lifecycle events
	return RunVerify(ctx, "provision", 5*time.Second)
}

func onPreDeploy(ctx context.Context, args *azdext.ProjectEventArgs) error {
	return RunVerify(ctx, "deploy", 5*time.Second)
}

func onPreUp(ctx context.Context, args *azdext.ProjectEventArgs) error {
	return RunVerify(ctx, "up", 5*time.Second)
}
