// Package types provides type definitions for the setup component
package types

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