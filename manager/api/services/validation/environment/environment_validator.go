// Package environment provides environment validation for different operating systems
package environment

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/api/middleware"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/api/types"
)

// EnvironmentValidator validates environment for different operating systems
type EnvironmentValidator struct {
	logger middleware.Logger
}

// NewEnvironmentValidator creates a new environment validator
func NewEnvironmentValidator(logger middleware.Logger) *EnvironmentValidator {
	return &EnvironmentValidator{
		logger: logger,
	}
}

// ValidateWindows validates Windows environment
func (ev *EnvironmentValidator) ValidateWindows(req *types.ValidationRequest, result *types.ValidationResult) error {
	ev.logger.Info("Validating Windows environment", map[string]interface{}{
		"interface": req.Interface,
	})

	// Detect Windows version and architecture
	osInfo, err := ev.detectWindowsInfo()
	if err != nil {
		result.Errors = append(result.Errors, types.ValidationItem{
			Code:      "WINDOWS_INFO_ERROR",
			Message:   fmt.Sprintf("Failed to detect Windows information: %v", err),
			Severity:  string(types.SeverityError),
			Category:  string(types.CategoryEnvironment),
			Fixable:   false,
			Timestamp: time.Now(),
		})
		return err
	}

	// Set environment information
	result.Environment.OS = osInfo.OS
	result.Environment.OSVersion = osInfo.OSVersion
	result.Environment.Architecture = osInfo.Architecture
	result.Environment.KernelVersion = osInfo.KernelVersion
	result.Environment.PathSeparator = "\\"

	// Validate Windows version compatibility
	if err := ev.validateWindowsVersion(osInfo.OSVersion, result); err != nil {
		ev.logger.Error("Windows version validation failed", map[string]interface{}{
			"version": osInfo.OSVersion,
			"error":   err.Error(),
		})
	}

	// Validate admin rights
	if err := ev.validateWindowsAdminRights(result); err != nil {
		ev.logger.Error("Windows admin rights validation failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Validate PowerShell
	if err := ev.validatePowerShell(result); err != nil {
		ev.logger.Error("PowerShell validation failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Validate disk space
	if err := ev.validateDiskSpace(result); err != nil {
		ev.logger.Error("Disk space validation failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Validate network connectivity
	if err := ev.validateNetworkConnectivity(result); err != nil {
		ev.logger.Error("Network connectivity validation failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Validate home directory
	if err := ev.validateHomeDirectory(result); err != nil {
		ev.logger.Error("Home directory validation failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Set Windows-specific capabilities
	result.Environment.Capabilities = []string{
		"windows_service",
		"powershell",
		"registry",
		"event_log",
		"wmi",
	}

	ev.logger.Info("Windows environment validation completed", map[string]interface{}{
		"os":           result.Environment.OS,
		"version":      result.Environment.OSVersion,
		"architecture": result.Environment.Architecture,
		"admin_rights": result.Environment.HasAdminRights,
	})

	return nil
}

// ValidateLinux validates Linux environment
func (ev *EnvironmentValidator) ValidateLinux(req *types.ValidationRequest, result *types.ValidationResult) error {
	ev.logger.Info("Validating Linux environment", map[string]interface{}{
		"interface": req.Interface,
	})

	// Detect Linux distribution and architecture
	osInfo, err := ev.detectLinuxInfo()
	if err != nil {
		result.Errors = append(result.Errors, types.ValidationItem{
			Code:      "LINUX_INFO_ERROR",
			Message:   fmt.Sprintf("Failed to detect Linux information: %v", err),
			Severity:  string(types.SeverityError),
			Category:  string(types.CategoryEnvironment),
			Fixable:   false,
			Timestamp: time.Now(),
		})
		return err
	}

	// Set environment information
	result.Environment.OS = osInfo.OS
	result.Environment.OSVersion = osInfo.OSVersion
	result.Environment.Architecture = osInfo.Architecture
	result.Environment.KernelVersion = osInfo.KernelVersion
	result.Environment.PathSeparator = "/"

	// Validate Linux distribution compatibility
	if err := ev.validateLinuxDistribution(osInfo.OS, osInfo.OSVersion, result); err != nil {
		ev.logger.Error("Linux distribution validation failed", map[string]interface{}{
			"distribution": osInfo.OS,
			"version":      osInfo.OSVersion,
			"error":        err.Error(),
		})
	}

	// Validate root privileges
	if err := ev.validateLinuxRootPrivileges(result); err != nil {
		ev.logger.Error("Linux root privileges validation failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Validate systemd (for service management)
	if err := ev.validateSystemd(result); err != nil {
		ev.logger.Error("Systemd validation failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Validate disk space
	if err := ev.validateDiskSpace(result); err != nil {
		ev.logger.Error("Disk space validation failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Validate network connectivity
	if err := ev.validateNetworkConnectivity(result); err != nil {
		ev.logger.Error("Network connectivity validation failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Validate home directory
	if err := ev.validateHomeDirectory(result); err != nil {
		ev.logger.Error("Home directory validation failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Set Linux-specific capabilities
	result.Environment.Capabilities = []string{
		"systemd_service",
		"bash",
		"cron",
		"logrotate",
		"iptables",
	}

	ev.logger.Info("Linux environment validation completed", map[string]interface{}{
		"os":              result.Environment.OS,
		"version":         result.Environment.OSVersion,
		"architecture":    result.Environment.Architecture,
		"root_privileges": result.Environment.HasAdminRights,
	})

	return nil
}

// ValidateDarwin validates macOS environment
func (ev *EnvironmentValidator) ValidateDarwin(req *types.ValidationRequest, result *types.ValidationResult) error {
	ev.logger.Info("Validating macOS environment", map[string]interface{}{
		"interface": req.Interface,
	})

	// Detect macOS version and architecture
	osInfo, err := ev.detectDarwinInfo()
	if err != nil {
		result.Errors = append(result.Errors, types.ValidationItem{
			Code:      "DARWIN_INFO_ERROR",
			Message:   fmt.Sprintf("Failed to detect macOS information: %v", err),
			Severity:  string(types.SeverityError),
			Category:  string(types.CategoryEnvironment),
			Fixable:   false,
			Timestamp: time.Now(),
		})
		return err
	}

	// Set environment information
	result.Environment.OS = osInfo.OS
	result.Environment.OSVersion = osInfo.OSVersion
	result.Environment.Architecture = osInfo.Architecture
	result.Environment.KernelVersion = osInfo.KernelVersion
	result.Environment.PathSeparator = "/"

	// Validate macOS version compatibility
	if err := ev.validateDarwinVersion(osInfo.OSVersion, result); err != nil {
		ev.logger.Error("macOS version validation failed", map[string]interface{}{
			"version": osInfo.OSVersion,
			"error":   err.Error(),
		})
	}

	// Validate admin privileges
	if err := ev.validateDarwinAdminPrivileges(result); err != nil {
		ev.logger.Error("macOS admin privileges validation failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Validate launchd (for service management)
	if err := ev.validateLaunchd(result); err != nil {
		ev.logger.Error("Launchd validation failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Validate disk space
	if err := ev.validateDiskSpace(result); err != nil {
		ev.logger.Error("Disk space validation failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Validate network connectivity
	if err := ev.validateNetworkConnectivity(result); err != nil {
		ev.logger.Error("Network connectivity validation failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Validate home directory
	if err := ev.validateHomeDirectory(result); err != nil {
		ev.logger.Error("Home directory validation failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Set macOS-specific capabilities
	result.Environment.Capabilities = []string{
		"launchd_service",
		"zsh",
		"cron",
		"log_rotation",
		"pfctl",
	}

	ev.logger.Info("macOS environment validation completed", map[string]interface{}{
		"os":               result.Environment.OS,
		"version":          result.Environment.OSVersion,
		"architecture":     result.Environment.Architecture,
		"admin_privileges": result.Environment.HasAdminRights,
	})

	return nil
}

// Windows-specific validation methods
func (ev *EnvironmentValidator) detectWindowsInfo() (*OSInfo, error) {
	// This is a simplified implementation
	// In production, you would use proper Windows APIs
	return &OSInfo{
		OS:            "windows",
		OSVersion:     "10.0.19042", // Example version
		Architecture:  runtime.GOARCH,
		KernelVersion: "10.0.19042",
	}, nil
}

func (ev *EnvironmentValidator) validateWindowsVersion(version string, result *types.ValidationResult) error {
	// Check if Windows version is supported (Windows 10 1903+ or Windows 11)
	if version == "" {
		result.Warnings = append(result.Warnings, types.ValidationItem{
			Code:      "UNKNOWN_WINDOWS_VERSION",
			Message:   "Unable to determine Windows version",
			Severity:  string(types.SeverityWarning),
			Category:  string(types.CategoryCompatibility),
			Fixable:   false,
			Timestamp: time.Now(),
		})
		return nil
	}

	// Simple version check (in production, use proper version parsing)
	if strings.Contains(version, "10.0") {
		result.Environment.Features = append(result.Environment.Features, "windows_10_support")
	} else {
		result.Warnings = append(result.Warnings, types.ValidationItem{
			Code:      "UNSUPPORTED_WINDOWS_VERSION",
			Message:   fmt.Sprintf("Windows version %s may not be fully supported", version),
			Severity:  string(types.SeverityWarning),
			Category:  string(types.CategoryCompatibility),
			Fixable:   false,
			Timestamp: time.Now(),
		})
	}

	return nil
}

func (ev *EnvironmentValidator) validateWindowsAdminRights(result *types.ValidationResult) error {
	// Check if running with admin privileges
	// This is a simplified check - in production, use proper Windows APIs
	result.Environment.HasAdminRights = true // Simplified for now

	if !result.Environment.HasAdminRights {
		result.Errors = append(result.Errors, types.ValidationItem{
			Code:     "INSUFFICIENT_PRIVILEGES",
			Message:  "Administrator privileges are required for setup",
			Severity: string(types.SeverityError),
			Category: string(types.CategorySecurity),
			Fixable:  true,
			AutoFix: &types.AutoFixInfo{
				Available: true,
				Command:   "Run as Administrator",
				Manual:    "Right-click and select 'Run as administrator'",
				Risk:      "low",
			},
			Timestamp: time.Now(),
		})
	}

	return nil
}

func (ev *EnvironmentValidator) validatePowerShell(result *types.ValidationResult) error {
	// Check PowerShell version (simplified)
	result.Environment.PowerShellVer = "5.1.19041.1320" // Example version

	// In production, execute: powershell -Command "$PSVersionTable.PSVersion"
	if result.Environment.PowerShellVer == "" {
		result.Errors = append(result.Errors, types.ValidationItem{
			Code:     "POWERSHELL_NOT_FOUND",
			Message:  "PowerShell is required but not found",
			Severity: string(types.SeverityError),
			Category: string(types.CategoryDependencies),
			Fixable:  true,
			AutoFix: &types.AutoFixInfo{
				Available: true,
				Command:   "Install PowerShell",
				Manual:    "Download and install PowerShell from Microsoft",
				Risk:      "low",
			},
			Timestamp: time.Now(),
		})
		return fmt.Errorf("PowerShell not found")
	}

	return nil
}

// Linux-specific validation methods
func (ev *EnvironmentValidator) detectLinuxInfo() (*OSInfo, error) {
	// This is a simplified implementation
	// In production, you would read from /etc/os-release
	return &OSInfo{
		OS:            "linux",
		OSVersion:     "Ubuntu 20.04", // Example
		Architecture:  runtime.GOARCH,
		KernelVersion: "5.4.0-74-generic", // Example
	}, nil
}

func (ev *EnvironmentValidator) validateLinuxDistribution(distribution, version string, result *types.ValidationResult) error {
	// Check if distribution is supported
	supportedDistros := []string{"ubuntu", "debian", "centos", "rhel", "fedora"}
	distroLower := strings.ToLower(distribution)

	supported := false
	for _, distro := range supportedDistros {
		if strings.Contains(distroLower, distro) {
			supported = true
			break
		}
	}

	if !supported {
		result.Warnings = append(result.Warnings, types.ValidationItem{
			Code:      "UNSUPPORTED_LINUX_DISTRO",
			Message:   fmt.Sprintf("Linux distribution '%s' may not be fully supported", distribution),
			Severity:  string(types.SeverityWarning),
			Category:  string(types.CategoryCompatibility),
			Fixable:   false,
			Timestamp: time.Now(),
		})
	}

	return nil
}

func (ev *EnvironmentValidator) validateLinuxRootPrivileges(result *types.ValidationResult) error {
	// Check if running as root or with sudo
	result.Environment.HasAdminRights = os.Geteuid() == 0

	if !result.Environment.HasAdminRights {
		result.Warnings = append(result.Warnings, types.ValidationItem{
			Code:     "NON_ROOT_USER",
			Message:  "Some operations may require root privileges",
			Severity: string(types.SeverityWarning),
			Category: string(types.CategorySecurity),
			Fixable:  true,
			AutoFix: &types.AutoFixInfo{
				Available: true,
				Command:   "sudo setup",
				Manual:    "Run with sudo or as root user",
				Risk:      "medium",
			},
			Timestamp: time.Now(),
		})
	}

	return nil
}

func (ev *EnvironmentValidator) validateSystemd(result *types.ValidationResult) error {
	// Check if systemd is available (simplified)
	result.Environment.Features = append(result.Environment.Features, "systemd")
	return nil
}

// macOS-specific validation methods
func (ev *EnvironmentValidator) detectDarwinInfo() (*OSInfo, error) {
	// This is a simplified implementation
	// In production, you would use proper macOS APIs
	return &OSInfo{
		OS:            "darwin",
		OSVersion:     "11.6", // macOS Big Sur example
		Architecture:  runtime.GOARCH,
		KernelVersion: "20.6.0", // Example
	}, nil
}

func (ev *EnvironmentValidator) validateDarwinVersion(version string, result *types.ValidationResult) error {
	// Check if macOS version is supported (10.15+)
	if version == "" {
		result.Warnings = append(result.Warnings, types.ValidationItem{
			Code:      "UNKNOWN_MACOS_VERSION",
			Message:   "Unable to determine macOS version",
			Severity:  string(types.SeverityWarning),
			Category:  string(types.CategoryCompatibility),
			Fixable:   false,
			Timestamp: time.Now(),
		})
		return nil
	}

	// Simple version check
	if strings.HasPrefix(version, "11.") || strings.HasPrefix(version, "12.") {
		result.Environment.Features = append(result.Environment.Features, "modern_macos")
	} else {
		result.Warnings = append(result.Warnings, types.ValidationItem{
			Code:      "OLD_MACOS_VERSION",
			Message:   fmt.Sprintf("macOS version %s may not be fully supported", version),
			Severity:  string(types.SeverityWarning),
			Category:  string(types.CategoryCompatibility),
			Fixable:   false,
			Timestamp: time.Now(),
		})
	}

	return nil
}

func (ev *EnvironmentValidator) validateDarwinAdminPrivileges(result *types.ValidationResult) error {
	// Check if user is in admin group (simplified)
	result.Environment.HasAdminRights = true // Simplified for now

	if !result.Environment.HasAdminRights {
		result.Warnings = append(result.Warnings, types.ValidationItem{
			Code:     "NON_ADMIN_USER",
			Message:  "Some operations may require administrator privileges",
			Severity: string(types.SeverityWarning),
			Category: string(types.CategorySecurity),
			Fixable:  true,
			AutoFix: &types.AutoFixInfo{
				Available: true,
				Command:   "sudo setup",
				Manual:    "Run with sudo or add user to admin group",
				Risk:      "medium",
			},
			Timestamp: time.Now(),
		})
	}

	return nil
}

func (ev *EnvironmentValidator) validateLaunchd(result *types.ValidationResult) error {
	// Check if launchd is available (it should be on all macOS systems)
	result.Environment.Features = append(result.Environment.Features, "launchd")
	return nil
}

// Common validation methods
func (ev *EnvironmentValidator) validateDiskSpace(result *types.ValidationResult) error {
	// Check available disk space (simplified)
	// In production, use proper disk space checking
	result.Environment.AvailableDiskGB = 100.0 // Example value

	if result.Environment.AvailableDiskGB < 1.0 {
		result.Errors = append(result.Errors, types.ValidationItem{
			Code:      "INSUFFICIENT_DISK_SPACE",
			Message:   "At least 1GB of free disk space is required",
			Severity:  string(types.SeverityError),
			Category:  string(types.CategoryPerformance),
			Expected:  ">= 1.0 GB",
			Actual:    fmt.Sprintf("%.2f GB", result.Environment.AvailableDiskGB),
			Fixable:   false,
			Timestamp: time.Now(),
		})
		return fmt.Errorf("insufficient disk space")
	}

	return nil
}

func (ev *EnvironmentValidator) validateNetworkConnectivity(result *types.ValidationResult) error {
	// Check network connectivity (simplified)
	result.Environment.HasInternet = true // Simplified for now

	if !result.Environment.HasInternet {
		result.Warnings = append(result.Warnings, types.ValidationItem{
			Code:      "NO_INTERNET_CONNECTIVITY",
			Message:   "No internet connectivity detected",
			Severity:  string(types.SeverityWarning),
			Category:  string(types.CategoryNetwork),
			Fixable:   false,
			Timestamp: time.Now(),
		})
	}

	return nil
}

func (ev *EnvironmentValidator) validateHomeDirectory(result *types.ValidationResult) error {
	// Get user home directory
	user, err := user.Current()
	if err != nil {
		result.Errors = append(result.Errors, types.ValidationItem{
			Code:      "USER_HOME_ERROR",
			Message:   fmt.Sprintf("Failed to get user home directory: %v", err),
			Severity:  string(types.SeverityError),
			Category:  string(types.CategoryEnvironment),
			Fixable:   false,
			Timestamp: time.Now(),
		})
		return err
	}

	result.Environment.HomeDir = user.HomeDir

	// Check if home directory is writable
	if err := os.MkdirAll(filepath.Join(user.HomeDir, ".syntropy"), 0755); err != nil {
		result.Errors = append(result.Errors, types.ValidationItem{
			Code:     "HOME_DIR_NOT_WRITABLE",
			Message:  "Home directory is not writable",
			Severity: string(types.SeverityError),
			Category: string(types.CategoryEnvironment),
			Field:    "home_dir",
			Fixable:  true,
			AutoFix: &types.AutoFixInfo{
				Available: true,
				Command:   "chmod 755 ~/.syntropy",
				Manual:    "Fix permissions on home directory",
				Risk:      "low",
			},
			Timestamp: time.Now(),
		})
		return err
	}

	return nil
}

// OSInfo represents operating system information
type OSInfo struct {
	OS            string
	OSVersion     string
	Architecture  string
	KernelVersion string
}
