package cli

import (
	"github.com/spf13/cobra"
	"github.com/timerzz/cc-venv/internal/env"
)

func newRemoveCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove an environment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return env.Remove(args[0], force)
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "remove without confirmation")
	return cmd
}
