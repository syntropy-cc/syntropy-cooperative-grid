// Package types provides type definitions for testing the setup component
package types

import (
	"errors"
	"time"
)

// SetupOptions defines the options for the setup process
type SetupOptions struct {
	Force          bool   // Force setup even if validations fail
	InstallService bool   // Install system service
	ConfigPath     string // Custom configuration file path
	HomeDir        string // Custom home directory
}

// SetupResult contains the result of the setup process
type SetupResult struct {
	Success     bool         // Indicates if the setup was successful
	StartTime   time.Time    // Setup start time
	EndTime     time.Time    // Setup end time
	ConfigPath  string       // Configuration file path
	Environment string       // Environment (windows, linux, darwin)
	Options     SetupOptions // Options used in the setup
	Error       error        // Error, if any
	Message     string       // Human-readable message
}

// ValidationResult represents the result of environment validation
type ValidationResult struct {
	Valid       bool     // Whether the environment is valid
	Warnings    []string // Warnings encountered during validation
	Errors      []string // Errors encountered during validation
	Environment EnvironmentInfo
}

// EnvironmentInfo contains information about the environment
type EnvironmentInfo struct {
	OS              string          // Operating system name
	OSVersion       string          // Operating system version
	Architecture    string          // System architecture
	HasAdminRights  bool            // Whether the user has admin rights
	PowerShellVer   string          // PowerShell version (Windows only)
	AvailableDiskGB float64         // Available disk space in GB
	HasInternet     bool            // Whether internet connectivity is available
	HomeDir         string          // User home directory
	SystemResources SystemResources // System resource information
}

// SetupConfig represents the configuration for the setup process
type SetupConfig struct {
	Manager     ManagerConfig // Manager configuration
	OwnerKey    OwnerKey      // Owner key configuration
	Environment Environment   // Environment configuration
}

// ManagerConfig represents the configuration for the Syntropy Manager
type ManagerConfig struct {
	HomeDir      string            // Home directory for Syntropy
	LogLevel     string            // Log level
	APIEndpoint  string            // API endpoint
	Directories  map[string]string // Directory paths
	DefaultPaths map[string]string // Default file paths
}

// OwnerKey represents the owner key configuration
type OwnerKey struct {
	Type      string // Key type (e.g., Ed25519)
	Path      string // Path to the key file
	PublicKey string // Public key
}

// Environment represents the environment configuration
type Environment struct {
	OS           string // Operating system
	Architecture string // System architecture
	HomeDir      string // User home directory
}

// SystemResources contains information about system resources
type SystemResources struct {
	TotalMemoryGB  float64 // Total memory in GB
	AvailableMemGB float64 // Available memory in GB
	CPUCores       int     // Number of CPU cores
	DiskSpaceGB    float64 // Available disk space in GB
}

// ErrNotImplemented is returned when a functionality is not implemented
var ErrNotImplemented = errors.New("functionality not implemented for this operating system")