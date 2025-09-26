// Package types provides shared validation type definitions for the API
package types

import (
	"time"
)

// ValidationResult represents the result of environment validation
type ValidationResult struct {
	Valid         bool              `json:"valid"`         // Whether the environment is valid
	Warnings      []ValidationItem  `json:"warnings"`      // Warnings encountered during validation
	Errors        []ValidationItem  `json:"errors"`        // Errors encountered during validation
	Environment   *EnvironmentInfo  `json:"environment"`   // Environment information
	Resources     *SystemResources  `json:"resources"`     // System resources
	Compatibility *Compatibility    `json:"compatibility"` // Compatibility information
	Security      *SecurityCheck    `json:"security"`      // Security validation
	Performance   *PerformanceCheck `json:"performance"`   // Performance validation
	Timestamp     time.Time         `json:"timestamp"`     // Validation timestamp
	Duration      time.Duration     `json:"duration"`      // Validation duration
	Interface     string            `json:"interface"`     // Interface type
	Version       string            `json:"version"`       // Validation version
}

// ValidationItem represents a single validation item (warning or error)
type ValidationItem struct {
	Code       string                 `json:"code"`               // Validation code
	Message    string                 `json:"message"`            // Human-readable message
	Severity   string                 `json:"severity"`           // Severity level (info, warning, error, critical)
	Category   string                 `json:"category"`           // Category (environment, security, performance, etc.)
	Field      string                 `json:"field"`              // Field that failed validation
	Value      interface{}            `json:"value"`              // Value that failed validation
	Expected   interface{}            `json:"expected"`           // Expected value
	Actual     interface{}            `json:"actual"`             // Actual value
	Suggestion string                 `json:"suggestion"`         // Suggestion to fix
	Fixable    bool                   `json:"fixable"`            // Whether this can be auto-fixed
	AutoFix    *AutoFixInfo           `json:"auto_fix,omitempty"` // Auto-fix information
	Metadata   map[string]interface{} `json:"metadata"`           // Additional metadata
	Timestamp  time.Time              `json:"timestamp"`          // When this validation occurred
}

// AutoFixInfo represents information for automatic fixing
type AutoFixInfo struct {
	Available bool                   `json:"available"` // Whether auto-fix is available
	Command   string                 `json:"command"`   // Command to fix
	Script    string                 `json:"script"`    // Script to fix
	Manual    string                 `json:"manual"`    // Manual fix instructions
	Risk      string                 `json:"risk"`      // Risk level (low, medium, high)
	Backup    bool                   `json:"backup"`    // Whether backup is required
	Metadata  map[string]interface{} `json:"metadata"`  // Additional fix metadata
}

// EnvironmentInfo contains information about the environment
type EnvironmentInfo struct {
	OS              string            `json:"os"`                // Operating system name
	OSVersion       string            `json:"os_version"`        // Operating system version
	Architecture    string            `json:"architecture"`      // System architecture
	KernelVersion   string            `json:"kernel_version"`    // Kernel version
	HasAdminRights  bool              `json:"has_admin_rights"`  // Whether the user has admin rights
	PowerShellVer   string            `json:"powershell_ver"`    // PowerShell version (Windows only)
	AvailableDiskGB float64           `json:"available_disk_gb"` // Available disk space in GB
	HasInternet     bool              `json:"has_internet"`      // Whether internet connectivity is available
	HomeDir         string            `json:"home_dir"`          // User home directory
	TempDir         string            `json:"temp_dir"`          // Temporary directory
	PathSeparator   string            `json:"path_separator"`    // Path separator for the OS
	EnvironmentVars map[string]string `json:"environment_vars"`  // Environment variables
	Features        []string          `json:"features"`          // Available features
	Limitations     []string          `json:"limitations"`       // System limitations
	Capabilities    []string          `json:"capabilities"`      // System capabilities
}

// SystemResources contains information about system resources
type SystemResources struct {
	TotalMemoryGB    float64   `json:"total_memory_gb"`     // Total memory in GB
	AvailableMemGB   float64   `json:"available_mem_gb"`    // Available memory in GB
	UsedMemoryGB     float64   `json:"used_memory_gb"`      // Used memory in GB
	CPUCores         int       `json:"cpu_cores"`           // Number of CPU cores
	CPUModel         string    `json:"cpu_model"`           // CPU model
	CPUSpeed         float64   `json:"cpu_speed"`           // CPU speed in GHz
	DiskSpaceGB      float64   `json:"disk_space_gb"`       // Available disk space in GB
	TotalDiskSpaceGB float64   `json:"total_disk_space_gb"` // Total disk space in GB
	NetworkSpeed     float64   `json:"network_speed"`       // Network speed in Mbps
	GPUMemory        float64   `json:"gpu_memory"`          // GPU memory in GB
	GPUModel         string    `json:"gpu_model"`           // GPU model
	LoadAverage      []float64 `json:"load_average"`        // System load average
}

// Compatibility represents compatibility information
type Compatibility struct {
	SupportedOS      []string          `json:"supported_os"`      // Supported operating systems
	MinOSVersion     map[string]string `json:"min_os_version"`    // Minimum OS versions
	RecommendedOS    []string          `json:"recommended_os"`    // Recommended operating systems
	Architecture     []string          `json:"architecture"`      // Supported architectures
	Dependencies     []Dependency      `json:"dependencies"`      // Required dependencies
	OptionalFeatures []Feature         `json:"optional_features"` // Optional features
	KnownIssues      []KnownIssue      `json:"known_issues"`      // Known issues
	Workarounds      []Workaround      `json:"workarounds"`       // Available workarounds
}

// Dependency represents a system dependency
type Dependency struct {
	Name        string `json:"name"`        // Dependency name
	Version     string `json:"version"`     // Required version
	Current     string `json:"current"`     // Current version
	Status      string `json:"status"`      // Status (installed, missing, outdated)
	Required    bool   `json:"required"`    // Whether it's required
	Installable bool   `json:"installable"` // Whether it can be installed
	InstallCmd  string `json:"install_cmd"` // Installation command
	CheckCmd    string `json:"check_cmd"`   // Check command
	Path        string `json:"path"`        // Installation path
}

// Feature represents a system feature
type Feature struct {
	Name        string `json:"name"`        // Feature name
	Description string `json:"description"` // Feature description
	Available   bool   `json:"available"`   // Whether it's available
	Required    bool   `json:"required"`    // Whether it's required
	Version     string `json:"version"`     // Feature version
	Enabled     bool   `json:"enabled"`     // Whether it's enabled
}

// KnownIssue represents a known issue
type KnownIssue struct {
	ID          string `json:"id"`          // Issue ID
	Description string `json:"description"` // Issue description
	Severity    string `json:"severity"`    // Severity level
	OS          string `json:"os"`          // Affected OS
	Version     string `json:"version"`     // Affected version
	Workaround  string `json:"workaround"`  // Available workaround
	Fixed       bool   `json:"fixed"`       // Whether it's fixed
	FixedIn     string `json:"fixed_in"`    // Fixed in version
}

// Workaround represents a workaround for an issue
type Workaround struct {
	ID          string `json:"id"`          // Workaround ID
	Description string `json:"description"` // Workaround description
	Command     string `json:"command"`     // Command to apply workaround
	Script      string `json:"script"`      // Script to apply workaround
	Manual      string `json:"manual"`      // Manual instructions
	Risk        string `json:"risk"`        // Risk level
	Reversible  bool   `json:"reversible"`  // Whether it's reversible
}

// SecurityCheck represents security validation
type SecurityCheck struct {
	EncryptionAvailable bool     `json:"encryption_available"` // Whether encryption is available
	SecureRandom        bool     `json:"secure_random"`        // Whether secure random is available
	KeyGeneration       bool     `json:"key_generation"`       // Whether key generation works
	FilePermissions     bool     `json:"file_permissions"`     // Whether file permissions are correct
	NetworkSecurity     bool     `json:"network_security"`     // Whether network is secure
	Vulnerabilities     []string `json:"vulnerabilities"`      // Known vulnerabilities
	Recommendations     []string `json:"recommendations"`      // Security recommendations
	Compliance          []string `json:"compliance"`           // Compliance standards met
}

// PerformanceCheck represents performance validation
type PerformanceCheck struct {
	DiskIOPerformance  float64     `json:"disk_io_performance"` // Disk I/O performance score
	NetworkPerformance float64     `json:"network_performance"` // Network performance score
	MemoryPerformance  float64     `json:"memory_performance"`  // Memory performance score
	CPUPerformance     float64     `json:"cpu_performance"`     // CPU performance score
	OverallScore       float64     `json:"overall_score"`       // Overall performance score
	Bottlenecks        []string    `json:"bottlenecks"`         // Performance bottlenecks
	Optimizations      []string    `json:"optimizations"`       // Optimization suggestions
	Benchmarks         []Benchmark `json:"benchmarks"`          // Performance benchmarks
}

// Benchmark represents a performance benchmark
type Benchmark struct {
	Name        string        `json:"name"`        // Benchmark name
	Description string        `json:"description"` // Benchmark description
	Score       float64       `json:"score"`       // Benchmark score
	Duration    time.Duration `json:"duration"`    // Benchmark duration
	Throughput  float64       `json:"throughput"`  // Throughput measurement
	Latency     time.Duration `json:"latency"`     // Latency measurement
	Status      string        `json:"status"`      // Benchmark status
}

// ValidationRequest represents a validation request
type ValidationRequest struct {
	Type        string                 `json:"type"`        // Validation type (environment, config, security, performance)
	Options     ValidationOptions      `json:"options"`     // Validation options
	Environment *EnvironmentInfo       `json:"environment"` // Environment to validate
	Interface   string                 `json:"interface"`   // Interface type
	UserID      string                 `json:"user_id"`     // User identifier
	SessionID   string                 `json:"session_id"`  // Session identifier
	CustomData  map[string]interface{} `json:"custom_data"` // Custom validation data
}

// ValidationOptions represents validation options
type ValidationOptions struct {
	SkipOptional      bool     `json:"skip_optional"`      // Skip optional validations
	AutoFix           bool     `json:"auto_fix"`           // Attempt auto-fix
	Detailed          bool     `json:"detailed"`           // Detailed validation
	Categories        []string `json:"categories"`         // Categories to validate
	ExcludeCategories []string `json:"exclude_categories"` // Categories to exclude
	Timeout           int      `json:"timeout"`            // Validation timeout in seconds
	Parallel          bool     `json:"parallel"`           // Run validations in parallel
}

// ValidationResponse represents a validation response
type ValidationResponse struct {
	Success bool              `json:"success"`          // Success status
	Result  *ValidationResult `json:"result,omitempty"` // Validation result
	Error   *ErrorDetail      `json:"error,omitempty"`  // Error details
	Message string            `json:"message"`          // Response message
	Code    int               `json:"code"`             // Response code
}

// ValidationCategory represents validation categories
type ValidationCategory string

const (
	CategoryEnvironment   ValidationCategory = "environment"
	CategorySecurity      ValidationCategory = "security"
	CategoryPerformance   ValidationCategory = "performance"
	CategoryCompatibility ValidationCategory = "compatibility"
	CategoryDependencies  ValidationCategory = "dependencies"
	CategoryConfiguration ValidationCategory = "configuration"
	CategoryNetwork       ValidationCategory = "network"
	CategoryStorage       ValidationCategory = "storage"
)
