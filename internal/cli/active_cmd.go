package cli

import (
	"github.com/spf13/cobra"
	"github.com/timerzz/cc-venv/internal/env"
)

func newActiveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "active <name>",
		Short: "Enter the environment shell",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			e, err := env.Load(args[0])
			if err != nil {
				return err
			}

			return env.Activate(e)
		},
	}
}
