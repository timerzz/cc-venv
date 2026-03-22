package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/timerzz/cc-venv/internal/web"
)

func newWebCmd() *cobra.Command {
	var (
		port    int
		devMode bool
		noOpen  bool
	)

	cmd := &cobra.Command{
		Use:   "web",
		Short: "Start the local web management server",
		Long: `Start a local HTTP server for managing ccv environments through a web interface.

The server provides a REST API for environment management and serves a web UI.
By default, it opens your browser automatically.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := web.Config{
				Port:    port,
				DevMode: devMode,
				NoOpen:  noOpen,
			}

			server := web.NewServer(cfg)

			// 启动浏览器
			if !noOpen {
				go func() {
					url := fmt.Sprintf("http://localhost:%d", port)
					if err := web.OpenBrowser(url); err != nil {
						fmt.Printf("warning: failed to open browser: %v\n", err)
					}
				}()
			}

			return server.Start()
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", 3000, "port to listen on")
	cmd.Flags().BoolVar(&devMode, "dev", false, "development mode (no embedded frontend)")
	cmd.Flags().BoolVar(&noOpen, "no-open", false, "don't open browser automatically")

	return cmd
}
