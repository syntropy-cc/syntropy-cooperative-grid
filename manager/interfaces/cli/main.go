package main

import (
	"fmt"
	"os"
	"runtime"

	setup "setup-component/src"

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

	// Token commands
	rootCmd.AddCommand(tokenCmd)

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
and prepare the environment for managing nodes in the cooperative grid.

A secure Grid Token is automatically generated and stored during setup.
Use --generate-grid-token to explicitly enable token generation (default behavior).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")
		configPath, _ := cmd.Flags().GetString("config-path")

		// Criar SetupManager para usar o novo sistema
		manager, err := setup.NewSetupManager()
		if err != nil {
			return fmt.Errorf("failed to create setup manager: %w", err)
		}

		// Preparar opções do novo sistema
		setupOptions := &setup.SetupOptions{
			Force:          force,
			ValidateOnly:   false,
			TestMode:       false,
			Verbose:        false,
			Quiet:          false,
			ConfigPath:     configPath,
			CustomSettings: map[string]string{
				// Token é gerado automaticamente por padrão
				// Só desabilita se explicitamente solicitado via flag
			},
		}

		fmt.Println("Starting Syntropy Manager setup...")

		// Executar setup com novo sistema
		err = manager.SetupWithPublicOptions(setupOptions)
		if err != nil {
			return fmt.Errorf("setup failed: %w", err)
		}

		fmt.Printf("✅ Setup completed successfully!\n")

		// Verificar se token foi criado
		exists, err := manager.GridTokenExists()
		if err == nil && exists {
			fmt.Printf("🔐 Grid Token generated and stored securely\n")
			fmt.Printf("📁 Token location: System Keyring (%s)\n", runtime.GOOS)
			fmt.Printf("💡 Use 'syntropy token show' to view token preview\n")
		} else {
			fmt.Printf("⚠️  Grid Token not created (may already exist or was disabled)\n")
		}

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

		options := setup.LegacySetupOptions{
			ConfigPath: configPath,
		}

		fmt.Println("Checking Syntropy Manager status...")
		result, err := setup.StatusLegacy(options)
		if err != nil {
			return fmt.Errorf("status check failed: %w", err)
		}

		if result.Success {
			fmt.Printf("✅ Syntropy Manager is properly configured\n")
			fmt.Printf("📁 Configuration: %s\n", result.ConfigPath)
			fmt.Printf("🖥️  Environment: %s\n", result.Environment)
		} else {
			fmt.Printf("❌ Setup issues detected: %v\n", result.Error)
			fmt.Printf("💡 Run 'syntropy setup run' to fix issues\n")
			return result.Error
		}

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

		options := setup.LegacySetupOptions{
			Force:      force,
			ConfigPath: configPath,
		}

		fmt.Println("Resetting Syntropy Manager configuration...")
		result, err := setup.ResetLegacy(options)
		if err != nil {
			return fmt.Errorf("reset failed: %w", err)
		}

		if result.Success {
			fmt.Printf("✅ Reset completed successfully!\n")
			fmt.Printf("💡 Run 'syntropy setup run' to reconfigure\n")
		} else {
			fmt.Printf("❌ Reset failed: %v\n", result.Error)
			return result.Error
		}

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

		options := setup.LegacySetupOptions{
			ConfigPath: configPath,
		}

		fmt.Println("Validating system environment...")
		result, err := setup.StatusLegacy(options)
		if err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}

		if result.Success {
			fmt.Printf("✅ System environment is ready for Syntropy Manager\n")
			fmt.Printf("🖥️  Environment: %s\n", result.Environment)
		} else {
			fmt.Printf("❌ Environment validation failed: %v\n", result.Error)
			fmt.Printf("💡 Fix the issues above and run validation again\n")
			return result.Error
		}

		return nil
	},
}

func init() {
	// Setup run flags
	setupRunCmd.Flags().Bool("force", false, "force setup even if validation fails")
	setupRunCmd.Flags().Bool("install-service", false, "install system service")
	setupRunCmd.Flags().String("config-path", "", "custom configuration file path")
	setupRunCmd.Flags().Bool("generate-grid-token", false, "generate Grid Token during setup")

	// Setup status flags
	setupStatusCmd.Flags().String("config-path", "", "custom configuration file path")

	// Setup reset flags
	setupResetCmd.Flags().Bool("force", false, "skip confirmation prompt")
	setupResetCmd.Flags().String("config-path", "", "custom configuration file path")

	// Setup validate flags
	setupValidateCmd.Flags().String("config-path", "", "custom configuration file path")
}

// tokenCmd represents the token command
var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Manage Grid Tokens for Syntropy Cooperative Grid",
	Long: `Manage Grid Tokens for the Syntropy Cooperative Grid network.

Grid Tokens are secure authentication tokens used to identify and authenticate
with the Syntropy Cooperative Grid network. They are stored securely in the
system keyring and can be managed through these commands.

Available operations:
- Generate new Grid Tokens
- View existing Grid Tokens (with security confirmation)
- Rotate Grid Tokens for enhanced security
- Export/Import Grid Tokens for backup and recovery
- Delete Grid Tokens when no longer needed`,
}

func init() {
	// Add token subcommands
	tokenCmd.AddCommand(tokenShowCmd)
	tokenCmd.AddCommand(tokenGenerateCmd)
	tokenCmd.AddCommand(tokenRotateCmd)
	tokenCmd.AddCommand(tokenExportCmd)
	tokenCmd.AddCommand(tokenImportCmd)
	tokenCmd.AddCommand(tokenDeleteCmd)
}

// tokenShowCmd represents the token show command
var tokenShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Display the current Grid Token (with security confirmation)",
	Long: `Display the current Grid Token with security confirmation.

This command will prompt for confirmation before displaying the token,
as Grid Tokens are sensitive authentication credentials that should be
kept secure. The token will be displayed in a masked format by default.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		showFull, _ := cmd.Flags().GetBool("full")
		confirm, _ := cmd.Flags().GetBool("confirm")

		// Create SetupManager
		manager, err := setup.NewSetupManager()
		if err != nil {
			return fmt.Errorf("failed to create setup manager: %w", err)
		}

		// Check if token exists
		exists, err := manager.GridTokenExists()
		if err != nil {
			return fmt.Errorf("failed to check token existence: %w", err)
		}

		if !exists {
			fmt.Println("❌ No Grid Token found")
			fmt.Println("💡 Run 'syntropy token generate' to create a new token")
			return nil
		}

		// Get token
		token, err := manager.GetGridToken()
		if err != nil {
			return fmt.Errorf("failed to retrieve token: %w", err)
		}

		if showFull && confirm {
			fmt.Printf("🔐 Grid Token: %s\n", token)
		} else if showFull {
			fmt.Println("⚠️  To display the full token, use --confirm flag")
			fmt.Printf("🔐 Grid Token Preview: %s\n", token[:8]+"...[HIDDEN]")
		} else {
			fmt.Printf("🔐 Grid Token Preview: %s\n", token[:8]+"...[HIDDEN]")
			fmt.Println("💡 Use --full --confirm to display the complete token")
		}

		return nil
	},
}

// tokenGenerateCmd represents the token generate command
var tokenGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a new Grid Token",
	Long: `Generate a new Grid Token for the Syntropy Cooperative Grid network.

This will create a new secure UUID v4 token and store it in the system keyring.
If a token already exists, it will be replaced with the new one.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		showToken, _ := cmd.Flags().GetBool("show")

		// Create SetupManager
		manager, err := setup.NewSetupManager()
		if err != nil {
			return fmt.Errorf("failed to create setup manager: %w", err)
		}

		// Check if token already exists
		exists, err := manager.GridTokenExists()
		if err != nil {
			return fmt.Errorf("failed to check token existence: %w", err)
		}

		if exists {
			fmt.Println("⚠️  A Grid Token already exists")
			fmt.Println("💡 Use 'syntropy token rotate' to generate a new token")
			return nil
		}

		// Generate new token
		fmt.Println("🔐 Generating new Grid Token...")
		token, err := manager.GenerateGridToken()
		if err != nil {
			return fmt.Errorf("failed to generate token: %w", err)
		}

		fmt.Println("✅ Grid Token generated successfully!")
		fmt.Printf("📁 Storage: System Keyring (%s)\n", runtime.GOOS)

		if showToken {
			fmt.Printf("🔐 Grid Token: %s\n", token)
		} else {
			fmt.Printf("🔐 Grid Token Preview: %s\n", token[:8]+"...[HIDDEN]")
			fmt.Println("💡 Use 'syntropy token show --full --confirm' to view the complete token")
		}

		return nil
	},
}

// tokenRotateCmd represents the token rotate command
var tokenRotateCmd = &cobra.Command{
	Use:   "rotate",
	Short: "Rotate the current Grid Token",
	Long: `Rotate the current Grid Token by generating a new one.

This will create a new secure UUID v4 token to replace the existing one.
The old token will be invalidated and the new token will be stored in the
system keyring.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		showToken, _ := cmd.Flags().GetBool("show")

		// Create SetupManager
		manager, err := setup.NewSetupManager()
		if err != nil {
			return fmt.Errorf("failed to create setup manager: %w", err)
		}

		// Check if token exists
		exists, err := manager.GridTokenExists()
		if err != nil {
			return fmt.Errorf("failed to check token existence: %w", err)
		}

		if !exists {
			fmt.Println("❌ No Grid Token found to rotate")
			fmt.Println("💡 Run 'syntropy token generate' to create a new token")
			return nil
		}

		// Rotate token
		fmt.Println("🔄 Rotating Grid Token...")
		newToken, err := manager.RotateGridToken()
		if err != nil {
			return fmt.Errorf("failed to rotate token: %w", err)
		}

		fmt.Println("✅ Grid Token rotated successfully!")
		fmt.Printf("📁 Storage: System Keyring (%s)\n", runtime.GOOS)

		if showToken {
			fmt.Printf("🔐 New Grid Token: %s\n", newToken)
		} else {
			fmt.Printf("🔐 New Grid Token Preview: %s\n", newToken[:8]+"...[HIDDEN]")
			fmt.Println("💡 Use 'syntropy token show --full --confirm' to view the complete token")
		}

		return nil
	},
}

// tokenExportCmd represents the token export command
var tokenExportCmd = &cobra.Command{
	Use:   "export [output-file]",
	Short: "Export Grid Token to a backup file",
	Long: `Export the current Grid Token to a backup file for safekeeping.

This will create a secure backup file containing the token with checksum
validation. The backup file should be stored securely as it contains
sensitive authentication credentials.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		outputPath := args[0]

		// Create SetupManager
		manager, err := setup.NewSetupManager()
		if err != nil {
			return fmt.Errorf("failed to create setup manager: %w", err)
		}

		// Check if token exists
		exists, err := manager.GridTokenExists()
		if err != nil {
			return fmt.Errorf("failed to check token existence: %w", err)
		}

		if !exists {
			fmt.Println("❌ No Grid Token found to export")
			fmt.Println("💡 Run 'syntropy token generate' to create a new token")
			return nil
		}

		// Export token
		fmt.Printf("📤 Exporting Grid Token to: %s\n", outputPath)
		err = manager.ExportGridToken(outputPath)
		if err != nil {
			return fmt.Errorf("failed to export token: %w", err)
		}

		fmt.Println("✅ Grid Token exported successfully!")
		fmt.Println("🔒 Backup file contains encrypted token with checksum validation")
		fmt.Println("⚠️  Store the backup file securely and do not share it")

		return nil
	},
}

// tokenImportCmd represents the token import command
var tokenImportCmd = &cobra.Command{
	Use:   "import [input-file]",
	Short: "Import Grid Token from a backup file",
	Long: `Import a Grid Token from a backup file.

This will restore a previously exported Grid Token from a backup file.
The backup file will be validated for integrity before importing.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputPath := args[0]

		// Create SetupManager
		manager, err := setup.NewSetupManager()
		if err != nil {
			return fmt.Errorf("failed to create setup manager: %w", err)
		}

		// Check if token already exists
		exists, err := manager.GridTokenExists()
		if err != nil {
			return fmt.Errorf("failed to check token existence: %w", err)
		}

		if exists {
			fmt.Println("⚠️  A Grid Token already exists")
			fmt.Println("💡 Use 'syntropy token rotate' to replace the existing token")
			return nil
		}

		// Import token
		fmt.Printf("📥 Importing Grid Token from: %s\n", inputPath)
		err = manager.ImportGridToken(inputPath)
		if err != nil {
			return fmt.Errorf("failed to import token: %w", err)
		}

		fmt.Println("✅ Grid Token imported successfully!")
		fmt.Printf("📁 Storage: System Keyring (%s)\n", runtime.GOOS)
		fmt.Println("💡 Use 'syntropy token show' to verify the imported token")

		return nil
	},
}

// tokenDeleteCmd represents the token delete command
var tokenDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete the current Grid Token",
	Long: `Delete the current Grid Token from the system keyring.

⚠️  WARNING: This will permanently remove the Grid Token and you will
need to generate a new one or import from backup to authenticate with
the Syntropy Cooperative Grid network.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")

		// Create SetupManager
		manager, err := setup.NewSetupManager()
		if err != nil {
			return fmt.Errorf("failed to create setup manager: %w", err)
		}

		// Check if token exists
		exists, err := manager.GridTokenExists()
		if err != nil {
			return fmt.Errorf("failed to check token existence: %w", err)
		}

		if !exists {
			fmt.Println("❌ No Grid Token found to delete")
			return nil
		}

		// Confirmation prompt
		if !force {
			fmt.Print("⚠️  Are you sure you want to delete the Grid Token? (yes/no): ")
			var response string
			fmt.Scanln(&response)
			if response != "yes" {
				fmt.Println("❌ Token deletion cancelled")
				return nil
			}
		}

		// Delete token
		fmt.Println("🗑️  Deleting Grid Token...")
		err = manager.DeleteGridToken()
		if err != nil {
			return fmt.Errorf("failed to delete token: %w", err)
		}

		fmt.Println("✅ Grid Token deleted successfully!")
		fmt.Println("💡 Use 'syntropy token generate' to create a new token")

		return nil
	},
}

func init() {
	// Token command flags
	tokenShowCmd.Flags().Bool("full", false, "show full token (requires --confirm)")
	tokenShowCmd.Flags().Bool("confirm", false, "confirm display of full token")

	tokenGenerateCmd.Flags().Bool("show", false, "display the generated token")

	tokenRotateCmd.Flags().Bool("show", false, "display the new token")

	tokenDeleteCmd.Flags().Bool("force", false, "skip confirmation prompt")
}

func main() {
	Execute()
}
