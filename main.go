// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"context"
	"os"

	"github.com/fatih/color"
	"spboyer.azd.doctor/internal/cmd"
)

func init() {
	forceColorVal, has := os.LookupEnv("FORCE_COLOR")
	if has && forceColorVal == "1" {
		color.NoColor = false
	}
}

func main() {
	ctx := context.Background()

	// Check if running in extension mode (invoked by azd for lifecycle events)
	// azd sets AZD_SERVER when running extensions or custom commands.
	// If no arguments are provided and AZD_SERVER is set, we assume it's the extension host mode.
	if len(os.Args) == 1 && os.Getenv("AZD_SERVER") != "" {
		if err := cmd.RunExtensionHost(ctx); err != nil {
			color.Red("Extension Host Error: %v", err)
			os.Exit(1)
		}
		return
	}

	// Execute the root command
	rootCmd := cmd.NewRootCommand()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		color.Red("Error: %v", err)
		os.Exit(1)
	}
}
