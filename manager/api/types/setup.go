// Package types provides shared type definitions for the API
package types

import (
	"errors"
	"time"
)

// SetupOptions defines the options for the setup process
type SetupOptions struct {
	Force          bool                   `json:"force" yaml:"force"`                     // Force setup even if validations fail
	InstallService bool                   `json:"install_service" yaml:"install_service"` // Install system service
	ConfigPath     string                 `json:"config_path" yaml:"config_path"`         // Custom configuration file path
	HomeDir        string                 `json:"home_dir" yaml:"home_dir"`               // Custom home directory
	Interface      string                 `json:"interface" yaml:"interface"`             // Interface type (cli, web, desktop, mobile)
	CustomOptions  map[string]interface{} `json:"custom_options" yaml:"custom_options"`   // Interface-specific options
}

// SetupResult contains the result of the setup process
type SetupResult struct {
	Success     bool          `json:"success"`         // Indicates if the setup was successful
	StartTime   time.Time     `json:"start_time"`      // Setup start time
	EndTime     time.Time     `json:"end_time"`        // Setup end time
	Duration    time.Duration `json:"duration"`        // Setup duration
	ConfigPath  string        `json:"config_path"`     // Configuration file path
	Environment string        `json:"environment"`     // Environment (windows, linux, darwin)
	Interface   string        `json:"interface"`       // Interface type
	Options     SetupOptions  `json:"options"`         // Options used in the setup
	Config      *SetupConfig  `json:"config"`          // Generated configuration
	Message     string        `json:"message"`         // Success/error message
	Error       error         `json:"error,omitempty"` // Error, if any
}

// SetupConfig represents the configuration for the setup process
type SetupConfig struct {
	Manager     ManagerConfig   `json:"manager" yaml:"manager"`         // Manager configuration
	OwnerKey    OwnerKey        `json:"owner_key" yaml:"owner_key"`     // Owner key configuration
	Environment Environment     `json:"environment" yaml:"environment"` // Environment configuration
	Interface   InterfaceConfig `json:"interface" yaml:"interface"`     // Interface-specific configuration
	Security    SecurityConfig  `json:"security" yaml:"security"`       // Security configuration
	Network     NetworkConfig   `json:"network" yaml:"network"`         // Network configuration
	Metadata    ConfigMetadata  `json:"metadata" yaml:"metadata"`       // Configuration metadata
}

// ManagerConfig represents the configuration for the Syntropy Manager
type ManagerConfig struct {
	HomeDir      string            `json:"home_dir" yaml:"home_dir"`           // Home directory for Syntropy
	LogLevel     string            `json:"log_level" yaml:"log_level"`         // Log level
	APIEndpoint  string            `json:"api_endpoint" yaml:"api_endpoint"`   // API endpoint
	Directories  map[string]string `json:"directories" yaml:"directories"`     // Directory paths
	DefaultPaths map[string]string `json:"default_paths" yaml:"default_paths"` // Default file paths
	Database     DatabaseConfig    `json:"database" yaml:"database"`           // Database configuration
}

// OwnerKey represents the owner key configuration
type OwnerKey struct {
	Type      string    `json:"type" yaml:"type"`             // Key type (e.g., Ed25519)
	Path      string    `json:"path" yaml:"path"`             // Path to the key file
	PublicKey string    `json:"public_key" yaml:"public_key"` // Public key
	CreatedAt time.Time `json:"created_at" yaml:"created_at"` // Creation timestamp
	Algorithm string    `json:"algorithm" yaml:"algorithm"`   // Key algorithm
	Size      int       `json:"size" yaml:"size"`             // Key size in bits
}

// Environment represents the environment configuration
type Environment struct {
	OS           string            `json:"os" yaml:"os"`                     // Operating system
	Architecture string            `json:"architecture" yaml:"architecture"` // System architecture
	HomeDir      string            `json:"home_dir" yaml:"home_dir"`         // User home directory
	Variables    map[string]string `json:"variables" yaml:"variables"`       // Environment variables
	Features     []string          `json:"features" yaml:"features"`         // Available features
}

// InterfaceConfig represents interface-specific configuration
type InterfaceConfig struct {
	Type        string                 `json:"type" yaml:"type"`               // Interface type
	Theme       string                 `json:"theme" yaml:"theme"`             // UI theme
	Language    string                 `json:"language" yaml:"language"`       // UI language
	Settings    map[string]interface{} `json:"settings" yaml:"settings"`       // Interface-specific settings
	Permissions []string               `json:"permissions" yaml:"permissions"` // Required permissions
}

// SecurityConfig represents security configuration
type SecurityConfig struct {
	EncryptionAlgorithm string   `json:"encryption_algorithm" yaml:"encryption_algorithm"` // Encryption algorithm
	KeyRotationDays     int      `json:"key_rotation_days" yaml:"key_rotation_days"`       // Key rotation interval
	AllowedIPs          []string `json:"allowed_ips" yaml:"allowed_ips"`                   // Allowed IP addresses
	SSLEnabled          bool     `json:"ssl_enabled" yaml:"ssl_enabled"`                   // SSL/TLS enabled
	CertPath            string   `json:"cert_path" yaml:"cert_path"`                       // Certificate path
}

// NetworkConfig represents network configuration
type NetworkConfig struct {
	Port        int      `json:"port" yaml:"port"`               // Default port
	Host        string   `json:"host" yaml:"host"`               // Default host
	Endpoints   []string `json:"endpoints" yaml:"endpoints"`     // API endpoints
	Timeout     int      `json:"timeout" yaml:"timeout"`         // Connection timeout
	Retries     int      `json:"retries" yaml:"retries"`         // Retry attempts
	Compression bool     `json:"compression" yaml:"compression"` // Enable compression
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Type     string `json:"type" yaml:"type"`         // Database type
	Host     string `json:"host" yaml:"host"`         // Database host
	Port     int    `json:"port" yaml:"port"`         // Database port
	Name     string `json:"name" yaml:"name"`         // Database name
	Username string `json:"username" yaml:"username"` // Database username
	SSLMode  string `json:"ssl_mode" yaml:"ssl_mode"` // SSL mode
}

// ConfigMetadata represents configuration metadata
type ConfigMetadata struct {
	Version     string    `json:"version" yaml:"version"`         // Configuration version
	CreatedAt   time.Time `json:"created_at" yaml:"created_at"`   // Creation timestamp
	UpdatedAt   time.Time `json:"updated_at" yaml:"updated_at"`   // Last update timestamp
	CreatedBy   string    `json:"created_by" yaml:"created_by"`   // Created by
	Interface   string    `json:"interface" yaml:"interface"`     // Interface type
	Environment string    `json:"environment" yaml:"environment"` // Environment
	Checksum    string    `json:"checksum" yaml:"checksum"`       // Configuration checksum
}

// SetupRequest represents a setup request from any interface
type SetupRequest struct {
	Options     SetupOptions           `json:"options" yaml:"options"`         // Setup options
	Environment *EnvironmentInfo       `json:"environment" yaml:"environment"` // Environment info
	Interface   string                 `json:"interface" yaml:"interface"`     // Interface type
	UserID      string                 `json:"user_id" yaml:"user_id"`         // User identifier
	SessionID   string                 `json:"session_id" yaml:"session_id"`   // Session identifier
	CustomData  map[string]interface{} `json:"custom_data" yaml:"custom_data"` // Custom data
}

// SetupResponse represents a setup response for any interface
type SetupResponse struct {
	Success bool         `json:"success"`          // Success status
	Result  *SetupResult `json:"result,omitempty"` // Setup result
	Error   *ErrorDetail `json:"error,omitempty"`  // Error details
	Message string       `json:"message"`          // Response message
	Code    int          `json:"code"`             // Response code
}

// ErrorDetail represents detailed error information
type ErrorDetail struct {
	Code    string `json:"code"`    // Error code
	Message string `json:"message"` // Error message
	Details string `json:"details"` // Error details
	Field   string `json:"field"`   // Field that caused the error
}

// ErrNotImplemented is returned when a functionality is not implemented
var ErrNotImplemented = errors.New("functionality not implemented for this operating system")

// ErrInvalidInterface is returned when an invalid interface is specified
var ErrInvalidInterface = errors.New("invalid interface type")

// ErrSetupAlreadyExists is returned when setup already exists
var ErrSetupAlreadyExists = errors.New("setup already exists")

// ErrSetupNotFound is returned when setup is not found
var ErrSetupNotFound = errors.New("setup not found")
