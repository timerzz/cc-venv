package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/timerzz/cc-venv/internal/importer"
)

func newImportCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "import <archive>",
		Short: "Import an environment archive",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := importer.ImportArchive(args[0], importer.Options{
				Force: force,
			})
			if err != nil {
				return err
			}

			fmt.Printf("Imported environment %q to %s\n", result.EnvName, result.Path)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "force overwrite existing environment")

	return cmd
}
