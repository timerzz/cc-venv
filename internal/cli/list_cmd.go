package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/timerzz/cc-venv/internal/env"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List environments",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			envs, err := env.List()
			if err != nil {
				return err
			}

			for _, e := range envs {
				fmt.Printf("%s\t%s\n", e.Name, e.RootPath)
			}

			return nil
		},
	}
}
