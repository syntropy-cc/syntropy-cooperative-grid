// Package types provides shared configuration type definitions for the API
package types

import (
	"time"
)

// ConfigRequest represents a configuration request
type ConfigRequest struct {
	Type        string                 `json:"type"`        // Configuration type (setup, manager, security, network)
	Options     ConfigOptions          `json:"options"`     // Configuration options
	Environment *EnvironmentInfo       `json:"environment"` // Environment information
	Interface   string                 `json:"interface"`   // Interface type
	UserID      string                 `json:"user_id"`     // User identifier
	SessionID   string                 `json:"session_id"`  // Session identifier
	Template    string                 `json:"template"`    // Configuration template
	CustomData  map[string]interface{} `json:"custom_data"` // Custom configuration data
}

// ConfigOptions represents configuration options
type ConfigOptions struct {
	Force            bool     `json:"force"`             // Force configuration generation
	Backup           bool     `json:"backup"`            // Create backup before changes
	Validate         bool     `json:"validate"`          // Validate configuration
	Encrypt          bool     `json:"encrypt"`           // Encrypt sensitive data
	Format           string   `json:"format"`            // Output format (yaml, json, toml)
	IncludeDefaults  bool     `json:"include_defaults"`  // Include default values
	Categories       []string `json:"categories"`        // Categories to include
	ExcludeSensitive bool     `json:"exclude_sensitive"` // Exclude sensitive data
}

// ConfigResponse represents a configuration response
type ConfigResponse struct {
	Success  bool            `json:"success"`            // Success status
	Config   *SetupConfig    `json:"config,omitempty"`   // Generated configuration
	Error    *ErrorDetail    `json:"error,omitempty"`    // Error details
	Message  string          `json:"message"`            // Response message
	Code     int             `json:"code"`               // Response code
	Metadata *ConfigMetadata `json:"metadata,omitempty"` // Configuration metadata
}

// ConfigTemplate represents a configuration template
type ConfigTemplate struct {
	Name        string                 `json:"name"`        // Template name
	Description string                 `json:"description"` // Template description
	Version     string                 `json:"version"`     // Template version
	Interface   string                 `json:"interface"`   // Interface type
	Environment string                 `json:"environment"` // Target environment
	Content     string                 `json:"content"`     // Template content
	Variables   []TemplateVariable     `json:"variables"`   // Template variables
	Validation  *TemplateValidation    `json:"validation"`  // Template validation rules
	Metadata    map[string]interface{} `json:"metadata"`    // Template metadata
}

// TemplateVariable represents a template variable
type TemplateVariable struct {
	Name        string      `json:"name"`        // Variable name
	Type        string      `json:"type"`        // Variable type (string, int, bool, etc.)
	Default     interface{} `json:"default"`     // Default value
	Required    bool        `json:"required"`    // Whether it's required
	Description string      `json:"description"` // Variable description
	Validation  string      `json:"validation"`  // Validation rule
	Options     []string    `json:"options"`     // Available options
	Condition   string      `json:"condition"`   // Conditional requirement
	Sensitive   bool        `json:"sensitive"`   // Whether it's sensitive data
}

// TemplateValidation represents template validation rules
type TemplateValidation struct {
	Schema   string                 `json:"schema"`   // JSON schema
	Rules    []ValidationRule       `json:"rules"`    // Validation rules
	Required []string               `json:"required"` // Required fields
	Optional []string               `json:"optional"` // Optional fields
	Defaults map[string]interface{} `json:"defaults"` // Default values
	Metadata map[string]interface{} `json:"metadata"` // Validation metadata
}

// ValidationRule represents a validation rule
type ValidationRule struct {
	Field     string      `json:"field"`     // Field to validate
	Type      string      `json:"type"`      // Validation type
	Pattern   string      `json:"pattern"`   // Regex pattern
	Min       interface{} `json:"min"`       // Minimum value
	Max       interface{} `json:"max"`       // Maximum value
	Values    []string    `json:"values"`    // Allowed values
	Message   string      `json:"message"`   // Error message
	Condition string      `json:"condition"` // Conditional validation
}

// ConfigBackup represents a configuration backup
type ConfigBackup struct {
	ID          string                 `json:"id"`          // Backup ID
	Name        string                 `json:"name"`        // Backup name
	Description string                 `json:"description"` // Backup description
	Config      *SetupConfig           `json:"config"`      // Backup configuration
	Timestamp   time.Time              `json:"timestamp"`   // Backup timestamp
	Size        int64                  `json:"size"`        // Backup size in bytes
	Checksum    string                 `json:"checksum"`    // Backup checksum
	Encrypted   bool                   `json:"encrypted"`   // Whether backup is encrypted
	Compressed  bool                   `json:"compressed"`  // Whether backup is compressed
	Metadata    map[string]interface{} `json:"metadata"`    // Backup metadata
}

// ConfigRestoreRequest represents a configuration restore request
type ConfigRestoreRequest struct {
	BackupID   string                 `json:"backup_id"`   // Backup ID to restore
	Options    RestoreOptions         `json:"options"`     // Restore options
	UserID     string                 `json:"user_id"`     // User identifier
	SessionID  string                 `json:"session_id"`  // Session identifier
	CustomData map[string]interface{} `json:"custom_data"` // Custom restore data
}

// RestoreOptions represents restore options
type RestoreOptions struct {
	Force      bool     `json:"force"`      // Force restore even if conflicts exist
	Backup     bool     `json:"backup"`     // Create backup before restore
	Validate   bool     `json:"validate"`   // Validate restored configuration
	Selective  bool     `json:"selective"`  // Selective restore
	Categories []string `json:"categories"` // Categories to restore
	Exclude    []string `json:"exclude"`    // Categories to exclude
	DryRun     bool     `json:"dry_run"`    // Dry run mode
}

// ConfigRestoreResponse represents a configuration restore response
type ConfigRestoreResponse struct {
	Success  bool          `json:"success"`          // Success status
	Config   *SetupConfig  `json:"config,omitempty"` // Restored configuration
	Backup   *ConfigBackup `json:"backup,omitempty"` // Created backup
	Error    *ErrorDetail  `json:"error,omitempty"`  // Error details
	Message  string        `json:"message"`          // Response message
	Code     int           `json:"code"`             // Response code
	Warnings []string      `json:"warnings"`         // Restore warnings
}

// ConfigListRequest represents a configuration list request
type ConfigListRequest struct {
	Type       string                 `json:"type"`       // Configuration type filter
	Interface  string                 `json:"interface"`  // Interface type filter
	UserID     string                 `json:"user_id"`    // User identifier
	SessionID  string                 `json:"session_id"` // Session identifier
	Pagination PaginationOptions      `json:"pagination"` // Pagination options
	Sort       SortOptions            `json:"sort"`       // Sort options
	Filter     map[string]interface{} `json:"filter"`     // Filter options
}

// PaginationOptions represents pagination options
type PaginationOptions struct {
	Page     int `json:"page"`      // Page number
	PageSize int `json:"page_size"` // Page size
	Total    int `json:"total"`     // Total items
}

// SortOptions represents sort options
type SortOptions struct {
	Field string `json:"field"` // Sort field
	Order string `json:"order"` // Sort order (asc, desc)
}

// ConfigListResponse represents a configuration list response
type ConfigListResponse struct {
	Success    bool               `json:"success"`         // Success status
	Configs    []ConfigSummary    `json:"configs"`         // Configuration list
	Error      *ErrorDetail       `json:"error,omitempty"` // Error details
	Message    string             `json:"message"`         // Response message
	Code       int                `json:"code"`            // Response code
	Pagination *PaginationOptions `json:"pagination"`      // Pagination info
}

// ConfigSummary represents a configuration summary
type ConfigSummary struct {
	ID          string    `json:"id"`          // Configuration ID
	Name        string    `json:"name"`        // Configuration name
	Type        string    `json:"type"`        // Configuration type
	Interface   string    `json:"interface"`   // Interface type
	Environment string    `json:"environment"` // Environment
	Version     string    `json:"version"`     // Configuration version
	CreatedAt   time.Time `json:"created_at"`  // Creation timestamp
	UpdatedAt   time.Time `json:"updated_at"`  // Last update timestamp
	Size        int64     `json:"size"`        // Configuration size
	Checksum    string    `json:"checksum"`    // Configuration checksum
	Status      string    `json:"status"`      // Configuration status
}

// ConfigType represents configuration types
type ConfigType string

const (
	ConfigTypeSetup     ConfigType = "setup"
	ConfigTypeManager   ConfigType = "manager"
	ConfigTypeSecurity  ConfigType = "security"
	ConfigTypeNetwork   ConfigType = "network"
	ConfigTypeDatabase  ConfigType = "database"
	ConfigTypeInterface ConfigType = "interface"
)

// ConfigFormat represents configuration formats
type ConfigFormat string

const (
	ConfigFormatYAML ConfigFormat = "yaml"
	ConfigFormatJSON ConfigFormat = "json"
	ConfigFormatTOML ConfigFormat = "toml"
	ConfigFormatXML  ConfigFormat = "xml"
	ConfigFormatINI  ConfigFormat = "ini"
)

// ConfigStatus represents configuration status
type ConfigStatus string

const (
	ConfigStatusActive   ConfigStatus = "active"
	ConfigStatusInactive ConfigStatus = "inactive"
	ConfigStatusDraft    ConfigStatus = "draft"
	ConfigStatusArchived ConfigStatus = "archived"
	ConfigStatusError    ConfigStatus = "error"
)
