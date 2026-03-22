package cli

import "github.com/spf13/cobra"

func Execute() error {
	return newRootCmd().Execute()
}

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "ccv",
		Short: "Manage named Claude Code virtual environments",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	rootCmd.AddCommand(
		newCreateCmd(),
		newListCmd(),
		newActiveCmd(),
		newRemoveCmd(),
		newRunCmd(),
		newWebCmd(),
		newExportCmd(),
		newImportCmd(),
	)

	return rootCmd
}
