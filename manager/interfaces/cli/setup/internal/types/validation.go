// Package types provides type definitions for the setup component
package types

// ValidationResult represents the result of environment validation
type ValidationResult struct {
	Valid       bool     // Whether the environment is valid
	Warnings    []string // Warnings encountered during validation
	Errors      []string // Errors encountered during validation
	Environment EnvironmentInfo
}

// EnvironmentInfo contains information about the environment
type EnvironmentInfo struct {
	OS              string  // Operating system name
	OSVersion       string  // Operating system version
	Architecture    string  // System architecture
	HasAdminRights  bool    // Whether the user has admin rights
	PowerShellVer   string  // PowerShell version (Windows only)
	AvailableDiskGB float64 // Available disk space in GB
	HasInternet     bool    // Whether internet connectivity is available
	HomeDir         string  // User home directory
}

// SystemResources contains information about system resources
type SystemResources struct {
	TotalMemoryGB  float64 // Total memory in GB
	AvailableMemGB float64 // Available memory in GB
	CPUCores       int     // Number of CPU cores
	DiskSpaceGB    float64 // Available disk space in GB
}
