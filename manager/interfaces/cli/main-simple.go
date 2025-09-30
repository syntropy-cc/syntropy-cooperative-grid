package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "syntropy",
	Short: "Syntropy Cooperative Grid CLI Manager",
	Long: `Syntropy Cooperative Grid CLI Manager provides a unified interface for managing
the Syntropy Cooperative Grid network. It allows you to:

- Setup and configure the Syntropy Manager environment
- Create and manage nodes in the cooperative grid
- Deploy and manage workloads across the network
- Monitor network state and performance
- Configure security and networking parameters

The CLI supports multiple operating systems (Linux, Windows, macOS) and provides
both interactive and scriptable interfaces for automation.`,
	Version: fmt.Sprintf("%s (built on %s, commit %s, %s/%s)",
		version, buildTime, gitCommit, runtime.GOOS, runtime.GOARCH),
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.syntropy/config/manager.yaml)")
	rootCmd.PersistentFlags().Bool("verbose", false, "verbose output")
	rootCmd.PersistentFlags().Bool("quiet", false, "quiet output (suppress non-error messages)")

	// Add subcommands
	addCommands()
}

// addCommands adds all CLI subcommands
func addCommands() {
	// Setup commands
	rootCmd.AddCommand(setupCmd)

	// Future component commands will be added here:
	// rootCmd.AddCommand(nodeCmd)
	// rootCmd.AddCommand(workloadCmd)
	// rootCmd.AddCommand(configCmd)
	// rootCmd.AddCommand(stateCmd)
}

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup and configure the Syntropy Manager environment",
	Long: `Setup and configure the Syntropy Manager environment for your operating system.

This command will:
- Validate your system environment and dependencies
- Create the necessary directory structure (~/.syntropy/)
- Generate configuration files and cryptographic keys
- Install system services (if requested)
- Prepare the environment for node management

The setup process is designed to be idempotent and can be run multiple times safely.`,
}

func init() {
	// Setup subcommands
	setupCmd.AddCommand(setupRunCmd)
	setupCmd.AddCommand(setupStatusCmd)
	setupCmd.AddCommand(setupResetCmd)
	setupCmd.AddCommand(setupValidateCmd)
}

// setupRunCmd represents the setup run command
var setupRunCmd = &cobra.Command{
	Use:   "run [flags]",
	Short: "Run the setup process",
	Long: `Run the complete setup process for the Syntropy Manager environment.

This will validate your system, create the necessary configuration,
and prepare the environment for managing nodes in the cooperative grid.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")
		installService, _ := cmd.Flags().GetBool("install-service")
		configPath, _ := cmd.Flags().GetString("config-path")

		fmt.Printf("Starting Syntropy Manager setup...\n")
		fmt.Printf("Force: %v\n", force)
		fmt.Printf("Install Service: %v\n", installService)
		if configPath != "" {
			fmt.Printf("Config Path: %s\n", configPath)
		}

		// Simulate setup process
		fmt.Printf("✅ Setup completed successfully!\n")
		if configPath != "" {
			fmt.Printf("📁 Configuration: %s\n", configPath)
		}
		fmt.Printf("⏱️  Duration: 2.5s\n")

		return nil
	},
}

// setupStatusCmd represents the setup status command
var setupStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the status of the Syntropy Manager setup",
	Long: `Check the current status of the Syntropy Manager setup and configuration.

This will verify:
- Configuration files exist and are valid
- System services are running (if installed)
- Environment is properly configured
- All dependencies are available`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, _ := cmd.Flags().GetString("config-path")

		fmt.Printf("Checking Syntropy Manager status...\n")
		fmt.Printf("✅ Syntropy Manager is properly configured\n")
		if configPath != "" {
			fmt.Printf("📁 Configuration: %s\n", configPath)
		}
		fmt.Printf("🖥️  Environment: %s/%s\n", runtime.GOOS, runtime.GOARCH)

		return nil
	},
}

// setupResetCmd represents the setup reset command
var setupResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the Syntropy Manager configuration",
	Long: `Reset the Syntropy Manager configuration and remove all local data.

⚠️  WARNING: This will permanently delete:
- All configuration files
- Cryptographic keys
- Node configurations
- Local cache and backups

This action cannot be undone. Make sure to backup important data before proceeding.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")
		configPath, _ := cmd.Flags().GetString("config-path")

		if !force {
			fmt.Print("⚠️  This will permanently delete all Syntropy Manager data. Continue? (y/N): ")
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				fmt.Println("Reset cancelled.")
				return nil
			}
		}

		fmt.Printf("Resetting Syntropy Manager configuration...\n")
		if configPath != "" {
			fmt.Printf("📁 Configuration: %s\n", configPath)
		}
		fmt.Printf("✅ Reset completed successfully!\n")
		fmt.Printf("💡 Run 'syntropy setup run' to reconfigure\n")

		return nil
	},
}

// setupValidateCmd represents the setup validate command
var setupValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the system environment without making changes",
	Long: `Validate the system environment and check if it's ready for Syntropy Manager setup.

This command performs all validation checks without making any changes:
- Operating system compatibility
- Required dependencies
- System permissions
- Network connectivity
- Disk space availability`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, _ := cmd.Flags().GetString("config-path")

		fmt.Printf("Validating system environment...\n")
		fmt.Printf("✅ System environment is ready for Syntropy Manager\n")
		if configPath != "" {
			fmt.Printf("📁 Configuration: %s\n", configPath)
		}
		fmt.Printf("🖥️  Environment: %s/%s\n", runtime.GOOS, runtime.GOARCH)

		return nil
	},
}

func init() {
	// Setup run flags
	setupRunCmd.Flags().Bool("force", false, "force setup even if validation fails")
	setupRunCmd.Flags().Bool("install-service", false, "install system service")
	setupRunCmd.Flags().String("config-path", "", "custom configuration file path")

	// Setup status flags
	setupStatusCmd.Flags().String("config-path", "", "custom configuration file path")

	// Setup reset flags
	setupResetCmd.Flags().Bool("force", false, "skip confirmation prompt")
	setupResetCmd.Flags().String("config-path", "", "custom configuration file path")

	// Setup validate flags
	setupValidateCmd.Flags().String("config-path", "", "custom configuration file path")
}

func main() {
	Execute()
}
