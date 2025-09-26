// Package setup provides functionality for setting up the Syntropy CLI environment
package setup

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	apiTypes "github.com/syntropy-cc/syntropy-cooperative-grid/manager/api/types"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src/internal/types"
)

// ErrNotImplemented is returned when a functionality is not implemented for the current operating system
var ErrNotImplemented = errors.New("functionality not implemented for this operating system")

// SetupOptions defines the options for the setup process
type SetupOptions = types.SetupOptions

// SetupResult contains the result of the setup process
type SetupResult = types.SetupResult

// Setup configures the environment for the Syntropy CLI
func Setup(options types.SetupOptions) (*types.SetupResult, error) {
	fmt.Println("Starting Syntropy CLI setup...")

	// Create API integration
	apiIntegration := NewAPIIntegration()

	// Convert local types to API types
	apiOptions := convertToAPISetupOptions(options)
	apiEnvironment := getCurrentEnvironment()

	// Use API central for setup
	apiResult, err := apiIntegration.SetupWithAPI(apiOptions, apiEnvironment, "cli")
	if err != nil {
		// Fallback to local implementation if API fails
		fmt.Println("API setup failed, falling back to local implementation...")
		switch runtime.GOOS {
		case "windows":
			return setupWindows(options)
		case "linux":
			return setupLinuxImpl(options)
		case "darwin":
			return setupDarwin(options)
		default:
			return nil, fmt.Errorf("%w: %s", ErrNotImplemented, runtime.GOOS)
		}
	}

	// Convert API result back to local types
	return convertFromAPISetupResult(apiResult), nil
}

// Status checks the installation status of the Syntropy CLI
func Status(options types.SetupOptions) (*types.SetupResult, error) {
	fmt.Println("Checking Syntropy CLI status...")

	// Create API integration
	apiIntegration := NewAPIIntegration()

	// Get status using API central
	status, err := apiIntegration.GetSetupStatusWithAPI("cli")
	if err != nil {
		// Fallback to local implementation if API fails
		fmt.Println("API status check failed, falling back to local implementation...")
		switch runtime.GOOS {
		case "windows":
			return statusWindows(options)
		case "linux":
			return statusLinux(options)
		case "darwin":
			return nil, fmt.Errorf("%w: %s", ErrNotImplemented, runtime.GOOS)
		default:
			return nil, fmt.Errorf("%w: %s", ErrNotImplemented, runtime.GOOS)
		}
	}

	// Convert API status to local result
	return convertStatusToSetupResult(status), nil
}

// Reset resets the Syntropy CLI configuration
func Reset(options types.SetupOptions) (*types.SetupResult, error) {
	fmt.Println("Resetting Syntropy CLI configuration...")

	// Create API integration
	apiIntegration := NewAPIIntegration()

	// Reset using API central
	err := apiIntegration.ResetSetupWithAPI("cli")
	if err != nil {
		// Fallback to local implementation if API fails
		fmt.Println("API reset failed, falling back to local implementation...")
		switch runtime.GOOS {
		case "windows":
			return resetWindows(options)
		case "linux":
			return resetLinux(options)
		case "darwin":
			return nil, fmt.Errorf("%w: %s", ErrNotImplemented, runtime.GOOS)
		default:
			return nil, fmt.Errorf("%w: %s", ErrNotImplemented, runtime.GOOS)
		}
	}

	// Return success result
	return &types.SetupResult{
		Success:   true,
		StartTime: time.Now(),
		EndTime:   time.Now(),
		Message:   "Reset completed successfully",
	}, nil
}

// GetSyntropyDir returns the default directory for the Syntropy CLI
func GetSyntropyDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to temporary directory in case of error
		return filepath.Join(os.TempDir(), "syntropy")
	}

	switch runtime.GOOS {
	case "windows":
		return filepath.Join(homeDir, "Syntropy")
	case "linux", "darwin":
		return filepath.Join(homeDir, ".syntropy")
	default:
		return filepath.Join(homeDir, ".syntropy")
	}
}

// setupDarwin implements macOS-specific setup (placeholder)
func setupDarwin(options types.SetupOptions) (*types.SetupResult, error) {
	return nil, fmt.Errorf("%w: darwin", ErrNotImplemented)
}

// setupWindows is a stub for Windows-specific function when compiled on other systems
func setupWindows(options types.SetupOptions) (*types.SetupResult, error) {
	return nil, fmt.Errorf("%w: windows (stub)", ErrNotImplemented)
}

// statusWindows is a stub for Windows-specific function when compiled on other systems
func statusWindows(options types.SetupOptions) (*types.SetupResult, error) {
	return nil, fmt.Errorf("%w: windows (stub)", ErrNotImplemented)
}

// resetWindows is a stub for Windows-specific function when compiled on other systems
func resetWindows(options types.SetupOptions) (*types.SetupResult, error) {
	return nil, fmt.Errorf("%w: windows (stub)", ErrNotImplemented)
}

// Conversion functions between local and API types

// convertToAPISetupOptions converts local SetupOptions to API SetupOptions
func convertToAPISetupOptions(local types.SetupOptions) *apiTypes.SetupOptions {
	return &apiTypes.SetupOptions{
		Force:          local.Force,
		InstallService: local.InstallService,
		ConfigPath:     local.ConfigPath,
		HomeDir:        local.HomeDir,
		Interface:      "cli",
		CustomOptions: map[string]interface{}{
			"source": "cli_setup",
		},
	}
}

// convertFromAPISetupResult converts API SetupResult to local SetupResult
func convertFromAPISetupResult(api *apiTypes.SetupResult) *types.SetupResult {
	return &types.SetupResult{
		Success:     api.Success,
		StartTime:   api.StartTime,
		EndTime:     api.EndTime,
		ConfigPath:  api.ConfigPath,
		Environment: api.Environment,
		Options: types.SetupOptions{
			Force:          api.Options.Force,
			InstallService: api.Options.InstallService,
			ConfigPath:     api.Options.ConfigPath,
			HomeDir:        api.Options.HomeDir,
		},
		Error: api.Error,
	}
}

// getCurrentEnvironment gets current environment information
func getCurrentEnvironment() *apiTypes.EnvironmentInfo {
	homeDir, _ := os.UserHomeDir()
	return &apiTypes.EnvironmentInfo{
		OS:              runtime.GOOS,
		OSVersion:       "unknown", // Would be populated by actual detection
		Architecture:    runtime.GOARCH,
		HomeDir:         homeDir,
		HasAdminRights:  true,  // Would be detected
		AvailableDiskGB: 100.0, // Would be calculated
		HasInternet:     true,  // Would be tested
		EnvironmentVars: make(map[string]string),
		Features:        []string{},
		Capabilities:    []string{},
	}
}

// convertStatusToSetupResult converts API status to local SetupResult
func convertStatusToSetupResult(status map[string]interface{}) *types.SetupResult {
	success := true
	if status["status"] != "active" {
		success = false
	}

	return &types.SetupResult{
		Success:     success,
		StartTime:   time.Now(),
		EndTime:     time.Now(),
		ConfigPath:  status["config_path"].(string),
		Environment: status["interface"].(string),
		Message:     fmt.Sprintf("Status: %s", status["status"]),
	}
}
