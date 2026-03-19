package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/timerzz/cc-venv/internal/exporter"
)

func newExportCmd() *cobra.Command {
	var outputPath string

	cmd := &cobra.Command{
		Use:   "export <name>",
		Short: "Export an environment to an archive",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := exporter.Export(args[0], exporter.Options{
				OutputPath: outputPath,
			})
			if err != nil {
				return err
			}

			fmt.Printf("Exported environment %q to %s\n", result.EnvName, result.ArchivePath)
			return nil
		},
	}

	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "output archive path (default: <name>-<timestamp>.tar.gz)")

	return cmd
}
