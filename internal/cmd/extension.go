package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/azure/azure-dev/cli/azd/pkg/azdext"
	"github.com/spf13/cobra"
)

func debugLog(format string, args ...interface{}) {
	f, _ := os.OpenFile("/tmp/azd-doctor-debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if f != nil {
		defer f.Close()
		fmt.Fprintf(f, format+"\n", args...)
	}
}

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
	debugLog("RunExtensionHost started")

	// Create a new context that includes the AZD access token
	ctx = azdext.WithAccessToken(ctx)

	// Create a new AZD client
	client, err := azdext.NewAzdClient()
	if err != nil {
		debugLog("Failed to create client: %v", err)
		return fmt.Errorf("failed to create azd client: %w", err)
	}
	defer client.Close()

	debugLog("Registering event handlers")

	// Use the ExtensionHost pattern which handles Ready() signaling automatically
	host := azdext.NewExtensionHost(client).
		WithServiceEventHandler("prepackage", onPrePackage, nil).
		WithProjectEventHandler("preprovision", onPreProvision).
		WithProjectEventHandler("predeploy", onPreDeploy)

	debugLog("Event handlers registered:")
	debugLog("  - service event: prepackage (no filters)")
	debugLog("  - project event: preprovision")
	debugLog("  - project event: predeploy")
	debugLog("Starting host.Run()")

	// Start listening for events
	// This is a blocking call and will not return until the server connection is closed
	if err := host.Run(ctx); err != nil {
		debugLog("host.Run() error: %v", err)
		return fmt.Errorf("failed to run extension: %w", err)
	}

	debugLog("host.Run() completed normally")
	return nil
}

func onPrePackage(ctx context.Context, args *azdext.ServiceEventArgs) error {
	debugLog("====== onPrePackage CALLED ======")
	debugLog("Service Name: %s", args.Service.Name)
	debugLog("Service Host: %s", args.Service.Host)
	debugLog("Service Language: %s", args.Service.Language)
	debugLog("==================================")

	// Write to stderr so azd shows it
	fmt.Fprintf(os.Stderr, "\n[azd doctor] Verifying environment for packaging service: %s\n", args.Service.Name)

	err := RunVerify(ctx, "package", 5*time.Second)
	if err != nil {
		debugLog("onPrePackage RunVerify returned error: %v", err)
		fmt.Fprintf(os.Stderr, "[azd doctor] Verification failed: %v\n\n", err)
	} else {
		debugLog("onPrePackage RunVerify completed successfully")
		fmt.Fprintf(os.Stderr, "[azd doctor] Environment verified\n\n")
	}
	return err
}

func onPreProvision(ctx context.Context, args *azdext.ProjectEventArgs) error {
	debugLog("onPreProvision called for project: %s", args.Project.Name)
	fmt.Fprintf(os.Stderr, "\n[azd doctor] Verifying environment for provisioning\n")

	// Default timeout for auth check in lifecycle events
	err := RunVerify(ctx, "provision", 5*time.Second)
	if err != nil {
		debugLog("onPreProvision RunVerify returned error: %v", err)
		fmt.Fprintf(os.Stderr, "[azd doctor] Verification failed: %v\n\n", err)
	} else {
		debugLog("onPreProvision RunVerify completed successfully")
		fmt.Fprintf(os.Stderr, "[azd doctor] Environment verified\n\n")
	}
	return err
}

func onPreDeploy(ctx context.Context, args *azdext.ProjectEventArgs) error {
	debugLog("onPreDeploy called for project: %s", args.Project.Name)
	fmt.Fprintf(os.Stderr, "\n[azd doctor] Verifying environment for deployment\n")

	err := RunVerify(ctx, "deploy", 5*time.Second)
	if err != nil {
		debugLog("onPreDeploy RunVerify returned error: %v", err)
		fmt.Fprintf(os.Stderr, "[azd doctor] Verification failed: %v\n\n", err)
	} else {
		debugLog("onPreDeploy RunVerify completed successfully")
		fmt.Fprintf(os.Stderr, "[azd doctor] Environment verified\n\n")
	}
	return err
}
