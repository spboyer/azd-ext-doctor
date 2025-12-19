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
	client, err := azdext.NewAzdClient()
	if err != nil {
		return fmt.Errorf("failed to create azd client: %w", err)
	}
	defer client.Close()

	// Use the extension ID from extension.yaml
	eventManager := azdext.NewEventManager("spboyer.azd.doctor", client)
	defer eventManager.Close()

	// Register handlers
	if err := eventManager.AddProjectEventHandler(ctx, "preprovision", onPreProvision); err != nil {
		return fmt.Errorf("failed to register preprovision handler: %w", err)
	}
	if err := eventManager.AddProjectEventHandler(ctx, "predeploy", onPreDeploy); err != nil {
		return fmt.Errorf("failed to register predeploy handler: %w", err)
	}
	if err := eventManager.AddProjectEventHandler(ctx, "preup", onPreUp); err != nil {
		return fmt.Errorf("failed to register preup handler: %w", err)
	}

	return eventManager.Receive(ctx)
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
