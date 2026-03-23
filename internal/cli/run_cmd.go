package cli

import (
	"github.com/spf13/cobra"
	"github.com/timerzz/cc-venv/internal/env"
)

func newRunCmd() *cobra.Command {
	return &cobra.Command{
		Use:                "run <name> [claude args...]",
		Short:              "Run Claude Code in an environment",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return cobra.MinimumNArgs(1)(cmd, args)
			}

			e, err := env.Load(args[0])
			if err != nil {
				return err
			}

			extraArgs := args[1:]
			if len(extraArgs) > 0 && extraArgs[0] == "--" {
				extraArgs = extraArgs[1:]
			}

			return env.RunClaude(e, extraArgs)
		},
	}
}
