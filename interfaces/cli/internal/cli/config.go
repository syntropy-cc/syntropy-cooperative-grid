package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewConfigCommand creates the configuration command
func NewConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
		Long: `Manage CLI configuration and settings.

Configuration includes API endpoints, authentication tokens, and 
other CLI-specific settings.`,
	}

	// Add subcommands
	cmd.AddCommand(newConfigShowCommand())
	cmd.AddCommand(newConfigSetCommand())
	cmd.AddCommand(newConfigGetCommand())
	cmd.AddCommand(newConfigResetCommand())

	return cmd
}

// newConfigShowCommand creates the config show command
func newConfigShowCommand() *cobra.Command {
	var (
		format string
	)

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		Long:  `Show the current CLI configuration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement config show
			fmt.Println("Current Configuration:")
			fmt.Println("  API Endpoint: https://api.syntropy.coop")
			fmt.Println("  Auth Token: ********************")
			fmt.Println("  Default Format: table")
			fmt.Println("  Timeout: 30s")
			fmt.Println("  Debug: false")
			fmt.Printf("Format: %s\n", format)

			return nil
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "table", "Output format (table, json, yaml)")

	return cmd
}

// newConfigSetCommand creates the config set command
func newConfigSetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Long:  `Set a configuration value for the CLI.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			value := args[1]

			// TODO: Implement config set
			fmt.Printf("Setting %s = %s\n", key, value)

			return nil
		},
	}

	return cmd
}

// newConfigGetCommand creates the config get command
func newConfigGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <key>",
		Short: "Get a configuration value",
		Long:  `Get a configuration value from the CLI.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]

			// TODO: Implement config get
			fmt.Printf("Getting %s\n", key)

			return nil
		},
	}

	return cmd
}

// newConfigResetCommand creates the config reset command
func newConfigResetCommand() *cobra.Command {
	var (
		force bool
	)

	cmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset configuration to defaults",
		Long:  `Reset all configuration values to their defaults.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Confirmation prompt
			if !force {
				fmt.Print("Are you sure you want to reset all configuration? (y/N): ")
				var response string
				fmt.Scanln(&response)
				if response != "y" && response != "Y" {
					fmt.Println("Operation cancelled.")
					return nil
				}
			}

			// TODO: Implement config reset
			fmt.Println("Resetting configuration to defaults...")

			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Skip confirmation prompt")

	return cmd
}
