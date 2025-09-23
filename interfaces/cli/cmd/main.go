package main

import (
	"os"

	"syntropy-cc/cooperative-grid/interfaces/cli/internal/cli"
)

func main() {
	rootCmd := cli.NewRootCommand()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}