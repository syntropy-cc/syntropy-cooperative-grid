package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/syntropy-cc/cooperative-grid/interfaces/cli/internal/cli"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	// Create root command
	rootCmd := &cobra.Command{
		Use:   "syntropy-cli",
		Short: "Syntropy Cooperative Grid Management CLI",
		Long: `Syntropy Cooperative Grid Management CLI

A comprehensive command-line interface for managing nodes, containers, 
networks, and cooperative services in the Syntropy Cooperative Grid.

Examples:
  syntropy-cli node list
  syntropy-cli node create --usb /dev/sdb --name "node-01"
  syntropy-cli container deploy --image nginx --node node-01
  syntropy-cli network status`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
	}

	// Add subcommands
	rootCmd.AddCommand(cli.NewNodeCommand())
	rootCmd.AddCommand(cli.NewContainerCommand())
	rootCmd.AddCommand(cli.NewNetworkCommand())
	rootCmd.AddCommand(cli.NewCooperativeCommand())
	rootCmd.AddCommand(cli.NewConfigCommand())

	// Execute command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
