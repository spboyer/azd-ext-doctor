// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package cmd

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "azd doctor <command> [options]",
		Short:         "Checks for pre-reqs and template requirements as needed.",
		SilenceUsage:  true,
		SilenceErrors: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug mode")

	rootCmd.AddCommand(newVersionCommand())
	rootCmd.AddCommand(NewCheckCommand())
	rootCmd.AddCommand(NewVerifyCommand())
	rootCmd.AddCommand(NewConfigureCommand())
	rootCmd.AddCommand(newContextCommand())
	rootCmd.AddCommand(NewListenCommand())

	return rootCmd
}
