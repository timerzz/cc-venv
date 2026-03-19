package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/timerzz/cc-venv/internal/env"
)

func newCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create <name>",
		Short: "Create a named environment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			e, err := env.Create(args[0])
			if err != nil {
				return err
			}

			fmt.Printf("created environment %q at %s\n", e.Name, e.RootPath)
			return nil
		},
	}
}
