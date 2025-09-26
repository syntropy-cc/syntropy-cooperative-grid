# Syntropy CLI Setup Component - API Documentation

## Overview
[MACRO_VIEW]
The Setup Component provides a comprehensive API for initializing and managing Syntropy Cooperative Grid environments across multiple platforms, offering both programmatic interfaces and CLI commands for setup, validation, and configuration management.
[/MACRO_VIEW]

[MESO_VIEW]
API design follows RESTful principles for remote operations and functional programming patterns for local operations, ensuring consistent behavior across different integration scenarios while maintaining backward compatibility and extensibility.
[/MESO_VIEW]

[MICRO_VIEW]
Interface contracts use strongly-typed Go structs with comprehensive validation, error handling, and result standardization, providing clear boundaries between public APIs and internal implementation details.
[/MICRO_VIEW]

## Public API Reference

### Core Functions

#### Setup Function
Initializes the Syntropy environment with comprehensive validation and configuration.

```go
func Setup(options types.SetupOptions) (*types.SetupResult, error)
```

**Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| options | types.SetupOptions | Yes | Configuration options for setup process |

**Returns:**
| Type | Description |
|------|-------------|
| *types.SetupResult | Detailed setup results including paths and configuration |
| error | Error information if setup fails |

**Example Usage:**
```go
package main

import (
    "fmt"
    "log"
    "github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup"
    "github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/internal/types"
)

func main() {
    // Basic setup with default options
    options := types.SetupOptions{
        Force:          false,
        InstallService: true,
        ConfigPath:     "", // Use default path
        HomeDir:        "", // Use default home directory
    }
    
    result, err := setup.Setup(options)
    if err != nil {
        log.Fatalf("Setup failed: %v", err)
    }
    
    fmt.Printf("Setup completed successfully!\n")
    fmt.Printf("Config path: %s\n", result.ConfigPath)
    fmt.Printf("Environment: %s/%s\n", result.Environment.OS, result.Environment.Architecture)
    fmt.Printf("Duration: %v\n", result.EndTime.Sub(result.StartTime))
}
```

**Advanced Usage with Custom Configuration:**
```go
func advancedSetup() {
    options := types.SetupOptions{
        Force:          true,  // Force reinstallation
        InstallService: true,  // Install system service
        ConfigPath:     "/custom/path/config.yaml",
        HomeDir:        "/custom/syntropy/home",
    }
    
    result, err := setup.Setup(options)
    if err != nil {
        // Handle specific error types
        switch err {
        case types.ErrNotImplemented:
            log.Printf("Platform not supported: %v", err)
        default:
            log.Printf("Setup error: %v", err)
        }
        return
    }
    
    // Process successful result
    if result.Success {
        fmt.Printf("Setup completed in %v\n", result.EndTime.Sub(result.StartTime))
        fmt.Printf("Configuration saved to: %s\n", result.ConfigPath)
    }
}
```

#### Status Function
Retrieves the current status of the Syntropy environment installation.

```go
func Status() (*types.SetupResult, error)
```

**Returns:**
| Type | Description |
|------|-------------|
| *types.SetupResult | Current installation status and configuration |
| error | Error information if status check fails |

**Example Usage:**
```go
func checkStatus() {
    status, err := setup.Status()
    if err != nil {
        log.Printf("Failed to get status: %v", err)
        return
    }
    
    if status.Success {
        fmt.Printf("Syntropy is installed and configured\n")
        fmt.Printf("Config location: %s\n", status.ConfigPath)
        fmt.Printf("Environment: %s\n", status.Environment.OS)
    } else {
        fmt.Printf("Syntropy is not properly installed\n")
        if status.Error != "" {
            fmt.Printf("Error: %s\n", status.Error)
        }
    }
}
```

#### Reset Function
Removes the Syntropy environment configuration and services.

```go
func Reset() error
```

**Returns:**
| Type | Description |
|------|-------------|
| error | Error information if reset fails |

**Example Usage:**
```go
func resetEnvironment() {
    fmt.Printf("Resetting Syntropy environment...\n")
    
    err := setup.Reset()
    if err != nil {
        log.Printf("Reset failed: %v", err)
        return
    }
    
    fmt.Printf("Environment reset successfully\n")
}
```

## Data Types Reference

### SetupOptions
Configuration options for the setup process.

```go
type SetupOptions struct {
    Force          bool   `json:"force" yaml:"force"`
    InstallService bool   `json:"install_service" yaml:"install_service"`
    ConfigPath     string `json:"config_path" yaml:"config_path"`
    HomeDir        string `json:"home_dir" yaml:"home_dir"`
}
```

**Field Descriptions:**
| Field | Type | Description | Default |
|-------|------|-------------|---------|
| Force | bool | Force reinstallation even if already configured | false |
| InstallService | bool | Install system service (systemd/Windows service) | true |
| ConfigPath | string | Custom path for configuration file | "" (auto-detect) |
| HomeDir | string | Custom home directory for Syntropy files | "" (auto-detect) |

**Validation Rules:**
- ConfigPath: Must be absolute path if specified, parent directory must exist
- HomeDir: Must be absolute path if specified, must be writable
- Force: When true, existing configurations will be overwritten
- InstallService: Requires appropriate system permissions

### SetupResult
Comprehensive result information from setup operations.

```go
type SetupResult struct {
    Success     bool                `json:"success" yaml:"success"`
    StartTime   time.Time          `json:"start_time" yaml:"start_time"`
    EndTime     time.Time          `json:"end_time" yaml:"end_time"`
    ConfigPath  string             `json:"config_path" yaml:"config_path"`
    Environment types.Environment  `json:"environment" yaml:"environment"`
    Options     types.SetupOptions `json:"options" yaml:"options"`
    Error       string             `json:"error,omitempty" yaml:"error,omitempty"`
}
```

**Field Descriptions:**
| Field | Type | Description |
|-------|------|-------------|
| Success | bool | Whether the operation completed successfully |
| StartTime | time.Time | When the operation started |
| EndTime | time.Time | When the operation completed |
| ConfigPath | string | Path to the generated configuration file |
| Environment | Environment | System environment information |
| Options | SetupOptions | Options used for the setup |
| Error | string | Error message if operation failed |

### Environment
System environment information collected during setup.

```go
type Environment struct {
    OS           string `json:"os" yaml:"os"`
    Architecture string `json:"architecture" yaml:"architecture"`
    HomeDir      string `json:"home_dir" yaml:"home_dir"`
}
```

**Field Descriptions:**
| Field | Type | Description | Possible Values |
|-------|------|-------------|-----------------|
| OS | string | Operating system | "linux", "windows", "darwin" |
| Architecture | string | System architecture | "amd64", "arm64", "386" |
| HomeDir | string | User's home directory | Platform-specific path |

### ValidationResult
Environment validation results with detailed feedback.

```go
type ValidationResult struct {
    Valid           bool                `json:"valid" yaml:"valid"`
    Warnings        []string           `json:"warnings,omitempty" yaml:"warnings,omitempty"`
    Errors          []string           `json:"errors,omitempty" yaml:"errors,omitempty"`
    EnvironmentInfo types.EnvironmentInfo `json:"environment_info" yaml:"environment_info"`
}
```

**Field Descriptions:**
| Field | Type | Description |
|-------|------|-------------|
| Valid | bool | Whether the environment passes all validation checks |
| Warnings | []string | Non-critical issues that should be addressed |
| Errors | []string | Critical issues that prevent setup |
| EnvironmentInfo | EnvironmentInfo | Detailed system information |

### EnvironmentInfo
Detailed system environment information for validation.

```go
type EnvironmentInfo struct {
    OS              string  `json:"os" yaml:"os"`
    OSVersion       string  `json:"os_version" yaml:"os_version"`
    Architecture    string  `json:"architecture" yaml:"architecture"`
    HasAdminRights  bool    `json:"has_admin_rights" yaml:"has_admin_rights"`
    PowerShellVer   string  `json:"powershell_version,omitempty" yaml:"powershell_version,omitempty"`
    AvailableDiskGB float64 `json:"available_disk_gb" yaml:"available_disk_gb"`
    HasInternet     bool    `json:"has_internet" yaml:"has_internet"`
    HomeDir         string  `json:"home_dir" yaml:"home_dir"`
}
```

## API Integration Interface

### APIIntegration
Interface for interacting with centralized Syntropy API services.

```go
type APIIntegration struct {
    configHandler   *config.Handler
    validationSvc   *validation.Service
    setupSvc        *setup.Service
    logger          *logger.Logger
}
```

#### NewAPIIntegration
Creates a new API integration instance.

```go
func NewAPIIntegration(configHandler *config.Handler, validationSvc *validation.Service, setupSvc *setup.Service, logger *logger.Logger) *APIIntegration
```

**Example Usage:**
```go
func createAPIIntegration() *APIIntegration {
    configHandler := config.NewHandler()
    validationSvc := validation.NewService()
    setupSvc := setup.NewService()
    logger := logger.NewLogger()
    
    return NewAPIIntegration(configHandler, validationSvc, setupSvc, logger)
}
```

#### SetupWithAPI
Performs setup using centralized API services.

```go
func (api *APIIntegration) SetupWithAPI(options types.SetupOptions) (*types.SetupResult, error)
```

**Example Usage:**
```go
func setupWithAPI() {
    api := createAPIIntegration()
    
    options := types.SetupOptions{
        Force:          false,
        InstallService: true,
    }
    
    result, err := api.SetupWithAPI(options)
    if err != nil {
        log.Printf("API setup failed: %v", err)
        return
    }
    
    fmt.Printf("API setup completed: %s\n", result.ConfigPath)
}
```

#### ValidateWithAPI
Validates environment using API services.

```go
func (api *APIIntegration) ValidateWithAPI() (*types.ValidationResult, error)
```

#### GetSetupStatusWithAPI
Retrieves setup status from API services.

```go
func (api *APIIntegration) GetSetupStatusWithAPI() (*types.SetupResult, error)
```

#### ResetSetupWithAPI
Resets setup using API services.

```go
func (api *APIIntegration) ResetSetupWithAPI() error
```

## Error Handling

### Error Types
The setup component defines specific error types for different failure scenarios.

```go
var ErrNotImplemented = errors.New("setup not implemented for this platform")
```

### Error Handling Patterns

#### Basic Error Handling
```go
result, err := setup.Setup(options)
if err != nil {
    // Log the error with context
    log.Printf("Setup failed: %v", err)
    
    // Check for specific error types
    if errors.Is(err, types.ErrNotImplemented) {
        fmt.Printf("This platform is not yet supported\n")
        return
    }
    
    // Handle generic errors
    fmt.Printf("Setup error: %v\n", err)
    return
}
```

#### Advanced Error Handling with Recovery
```go
func robustSetup(options types.SetupOptions) (*types.SetupResult, error) {
    // Attempt setup with retry logic
    maxRetries := 3
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        result, err := setup.Setup(options)
        if err == nil {
            return result, nil
        }
        
        lastErr = err
        
        // Check if error is retryable
        if isRetryableError(err) {
            log.Printf("Setup attempt %d failed, retrying: %v", i+1, err)
            time.Sleep(time.Duration(i+1) * time.Second)
            continue
        }
        
        // Non-retryable error, fail immediately
        break
    }
    
    return nil, fmt.Errorf("setup failed after %d attempts: %w", maxRetries, lastErr)
}

func isRetryableError(err error) bool {
    // Define which errors are worth retrying
    retryableErrors := []string{
        "network timeout",
        "connection refused",
        "temporary failure",
    }
    
    errStr := err.Error()
    for _, retryable := range retryableErrors {
        if strings.Contains(errStr, retryable) {
            return true
        }
    }
    
    return false
}
```

## Configuration Management

### Configuration Structure
The setup component generates YAML configuration files with the following structure:

```yaml
# Generated Syntropy Manager Configuration
home_dir: "/home/user/.syntropy"
log_level: "info"
api_endpoint: "https://api.syntropy.com"

directories:
  config: "/home/user/.syntropy/config"
  keys: "/home/user/.syntropy/keys"
  nodes: "/home/user/.syntropy/nodes"
  logs: "/home/user/.syntropy/logs"
  cache: "/home/user/.syntropy/cache"
  backups: "/home/user/.syntropy/backups"
  scripts: "/home/user/.syntropy/scripts"

default_paths:
  owner_key: "/home/user/.syntropy/keys/owner.key"
  config: "/home/user/.syntropy/config/manager.yaml"

owner_key:
  type: "ed25519"
  path: "/home/user/.syntropy/keys/owner.key"
  public_key: "ed25519_public_key_here"

environment:
  os: "linux"
  architecture: "amd64"
  home_dir: "/home/user"
```

### Configuration Access
```go
func loadConfiguration(configPath string) (*types.SetupConfig, error) {
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config: %w", err)
    }
    
    var config types.SetupConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse config: %w", err)
    }
    
    return &config, nil
}
```

## Platform-Specific APIs

### Linux Implementation
Linux-specific functions with build tags `//go:build linux`.

```go
func setupLinuxImpl(options types.SetupOptions) (*types.SetupResult, error)
func statusLinux() (*types.SetupResult, error)
func resetLinux() error
```

**Linux-Specific Features:**
- Systemd service installation
- Directory permissions (0755 for directories, 0600 for keys)
- Package manager integration checks
- SELinux compatibility validation

### Windows Implementation
Windows-specific functions with build tags `//go:build windows`.

```go
func setupWindows(options types.SetupOptions) (*types.SetupResult, error)
func statusWindows() (*types.SetupResult, error)
func resetWindows() error
```

**Windows-Specific Features:**
- Windows Service installation
- Registry configuration
- PowerShell script execution
- UAC permission handling

## Testing APIs

### Test Utilities
The component provides testing utilities for integration testing.

```go
func CreateTestEnvironment(t *testing.T) (string, func())
func MockAPIIntegration() *APIIntegration
func ValidateTestResult(t *testing.T, result *types.SetupResult)
```

**Example Test Usage:**
```go
func TestSetupIntegration(t *testing.T) {
    // Create isolated test environment
    testDir, cleanup := CreateTestEnvironment(t)
    defer cleanup()
    
    // Configure test options
    options := types.SetupOptions{
        Force:      true,
        HomeDir:    testDir,
        ConfigPath: filepath.Join(testDir, "config.yaml"),
    }
    
    // Execute setup
    result, err := setup.Setup(options)
    require.NoError(t, err)
    
    // Validate results
    ValidateTestResult(t, result)
    assert.True(t, result.Success)
    assert.FileExists(t, result.ConfigPath)
}
```

## Migration and Compatibility

### Version Compatibility
| API Version | Supported Versions | Breaking Changes |
|-------------|-------------------|------------------|
| 2.x | Manager 3.x, CLI 2.x | None |
| 1.x | Manager 2.x, CLI 1.x | SetupOptions structure changed |

### Migration Examples
```go
// Migrating from v1.x to v2.x
func migrateFromV1() {
    // Old v1.x usage
    // result, err := setup.Setup(force, installService, configPath)
    
    // New v2.x usage
    options := types.SetupOptions{
        Force:          true,
        InstallService: true,
        ConfigPath:     "/custom/path/config.yaml",
    }
    result, err := setup.Setup(options)
    // Handle result...
}
```

## Performance Considerations

### Optimization Guidelines
| Operation | Typical Duration | Optimization Strategy |
|-----------|------------------|----------------------|
| Setup | 5-15 seconds | Parallel validation, template caching |
| Status Check | 100-500ms | Configuration caching, lazy loading |
| Reset | 2-5 seconds | Batch file operations, service cleanup |

### Resource Management
```go
func efficientSetup(options types.SetupOptions) (*types.SetupResult, error) {
    // Use context for timeout control
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Implement setup with context
    return setupWithContext(ctx, options)
}
```

## Security Best Practices

### Secure API Usage
```go
func secureSetup() {
    options := types.SetupOptions{
        Force:          false, // Avoid force unless necessary
        InstallService: true,  // Install service for proper isolation
        ConfigPath:     "",    // Use default secure paths
        HomeDir:        "",    // Use default user directory
    }
    
    // Validate options before setup
    if err := validateSecureOptions(options); err != nil {
        log.Fatalf("Insecure options: %v", err)
    }
    
    result, err := setup.Setup(options)
    if err != nil {
        // Don't log sensitive information
        log.Printf("Setup failed: %v", sanitizeError(err))
        return
    }
    
    // Verify secure configuration
    if err := verifySecureSetup(result); err != nil {
        log.Printf("Security validation failed: %v", err)
    }
}
```

## Troubleshooting

### Common API Issues
| Issue | Symptoms | Solution |
|-------|----------|----------|
| Permission Denied | Setup fails with access errors | Run with appropriate privileges or use non-privileged paths |
| Network Timeout | API integration fails | Check network connectivity and API endpoint availability |
| Invalid Configuration | Setup succeeds but status fails | Validate configuration file syntax and required fields |
| Service Installation Failure | Setup completes but service not running | Check system service requirements and permissions |

### Debug Mode
```go
func debugSetup() {
    // Enable debug logging
    os.Setenv("SYNTROPY_DEBUG", "true")
    os.Setenv("SYNTROPY_LOG_LEVEL", "debug")
    
    options := types.SetupOptions{
        Force:          true,
        InstallService: true,
    }
    
    result, err := setup.Setup(options)
    if err != nil {
        // Debug information will be logged automatically
        log.Printf("Debug setup failed: %v", err)
    }
}
```