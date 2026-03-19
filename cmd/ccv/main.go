package main

import (
	"fmt"
	"os"

	"github.com/timerzz/cc-venv/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "ccv:", err)
		os.Exit(1)
	}
}
