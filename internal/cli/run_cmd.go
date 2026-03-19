package cli

import (
	"github.com/spf13/cobra"
	"github.com/timerzz/cc-venv/internal/env"
)

func newRunCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run <name>",
		Short: "Run Claude Code in an environment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			e, err := env.Load(args[0])
			if err != nil {
				return err
			}

			return env.RunClaude(e)
		},
	}
}
