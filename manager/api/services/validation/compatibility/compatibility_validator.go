// Package compatibility provides compatibility validation for the API
package compatibility

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/api/middleware"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/api/types"
)

// CompatibilityValidator validates compatibility aspects of the environment
type CompatibilityValidator struct {
	logger middleware.Logger
}

// NewCompatibilityValidator creates a new compatibility validator
func NewCompatibilityValidator(logger middleware.Logger) *CompatibilityValidator {
	return &CompatibilityValidator{
		logger: logger,
	}
}

// Validate performs compatibility validation
func (cv *CompatibilityValidator) Validate(req *types.ValidationRequest, result *types.ValidationResult) error {
	cv.logger.Info("Starting compatibility validation", map[string]interface{}{
		"interface": req.Interface,
		"os":        runtime.GOOS,
	})

	// Initialize compatibility information
	result.Compatibility = &types.Compatibility{
		SupportedOS:      []string{},
		MinOSVersion:     make(map[string]string),
		RecommendedOS:    []string{},
		Architecture:     []string{},
		Dependencies:     []types.Dependency{},
		OptionalFeatures: []types.Feature{},
		KnownIssues:      []types.KnownIssue{},
		Workarounds:      []types.Workaround{},
	}

	// Validate operating system compatibility
	if err := cv.validateOperatingSystem(req, result); err != nil {
		cv.logger.Error("Operating system compatibility validation failed", map[string]interface{}{
			"error": err.Error(),
			"os":    runtime.GOOS,
		})
	}

	// Validate architecture compatibility
	if err := cv.validateArchitecture(req, result); err != nil {
		cv.logger.Error("Architecture compatibility validation failed", map[string]interface{}{
			"error": err.Error(),
			"arch":  runtime.GOARCH,
		})
	}

	// Validate interface compatibility
	if err := cv.validateInterfaceCompatibility(req, result); err != nil {
		cv.logger.Error("Interface compatibility validation failed", map[string]interface{}{
			"error":     err.Error(),
			"interface": req.Interface,
		})
	}

	// Check known issues
	if err := cv.checkKnownIssues(req, result); err != nil {
		cv.logger.Error("Known issues check failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Check available workarounds
	if err := cv.checkWorkarounds(req, result); err != nil {
		cv.logger.Error("Workarounds check failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	cv.logger.Info("Compatibility validation completed", map[string]interface{}{
		"supported_os": result.Compatibility.SupportedOS,
		"architecture": result.Compatibility.Architecture,
		"known_issues": len(result.Compatibility.KnownIssues),
		"workarounds":  len(result.Compatibility.Workarounds),
	})

	return nil
}

// validateOperatingSystem validates operating system compatibility
func (cv *CompatibilityValidator) validateOperatingSystem(req *types.ValidationRequest, result *types.ValidationResult) error {
	currentOS := runtime.GOOS

	// Define supported operating systems
	supportedOS := []string{"windows", "linux", "darwin"}
	result.Compatibility.SupportedOS = supportedOS

	// Define minimum OS versions
	minVersions := map[string]string{
		"windows": "10.0.1903", // Windows 10 version 1903
		"linux":   "4.15.0",    // Linux kernel 4.15+
		"darwin":  "10.15.0",   // macOS Catalina
	}
	result.Compatibility.MinOSVersion = minVersions

	// Define recommended operating systems
	recommendedOS := []string{"windows", "ubuntu", "darwin"}
	result.Compatibility.RecommendedOS = recommendedOS

	// Check if current OS is supported
	isSupported := false
	for _, os := range supportedOS {
		if os == currentOS {
			isSupported = true
			break
		}
	}

	if !isSupported {
		result.Errors = append(result.Errors, types.ValidationItem{
			Code:      "UNSUPPORTED_OS",
			Message:   fmt.Sprintf("Operating system '%s' is not supported", currentOS),
			Severity:  string(types.SeverityError),
			Category:  string(types.CategoryCompatibility),
			Fixable:   false,
			Timestamp: time.Now(),
		})
		return fmt.Errorf("unsupported operating system: %s", currentOS)
	}

	// Check OS version compatibility
	if minVersion, exists := minVersions[currentOS]; exists {
		// This is a simplified version check
		// In production, you would parse and compare actual OS versions
		cv.logger.Info("OS version compatibility check", map[string]interface{}{
			"os":          currentOS,
			"min_version": minVersion,
		})

		// For now, assume version is compatible
		result.Environment.Features = append(result.Environment.Features, fmt.Sprintf("%s_support", currentOS))
	}

	return nil
}

// validateArchitecture validates architecture compatibility
func (cv *CompatibilityValidator) validateArchitecture(req *types.ValidationRequest, result *types.ValidationResult) error {
	currentArch := runtime.GOARCH

	// Define supported architectures
	supportedArchs := []string{"amd64", "386", "arm64", "arm"}
	result.Compatibility.Architecture = supportedArchs

	// Check if current architecture is supported
	isSupported := false
	for _, arch := range supportedArchs {
		if arch == currentArch {
			isSupported = true
			break
		}
	}

	if !isSupported {
		result.Warnings = append(result.Warnings, types.ValidationItem{
			Code:      "UNSUPPORTED_ARCHITECTURE",
			Message:   fmt.Sprintf("Architecture '%s' may not be fully supported", currentArch),
			Severity:  string(types.SeverityWarning),
			Category:  string(types.CategoryCompatibility),
			Fixable:   false,
			Timestamp: time.Now(),
		})
	}

	// Add architecture-specific features
	switch currentArch {
	case "amd64":
		result.Environment.Features = append(result.Environment.Features, "x86_64_support")
	case "arm64":
		result.Environment.Features = append(result.Environment.Features, "arm64_support")
	case "arm":
		result.Environment.Features = append(result.Environment.Features, "arm_support")
	case "386":
		result.Environment.Features = append(result.Environment.Features, "x86_support")
	}

	return nil
}

// validateInterfaceCompatibility validates interface-specific compatibility
func (cv *CompatibilityValidator) validateInterfaceCompatibility(req *types.ValidationRequest, result *types.ValidationResult) error {
	interfaceType := req.Interface

	// Define interface compatibility matrix
	compatibilityMatrix := map[string]map[string]bool{
		"windows": {
			"cli":     true,
			"web":     true,
			"desktop": true,
			"mobile":  false, // Windows Mobile is deprecated
		},
		"linux": {
			"cli":     true,
			"web":     true,
			"desktop": true,
			"mobile":  false, // Limited mobile support on Linux
		},
		"darwin": {
			"cli":     true,
			"web":     true,
			"desktop": true,
			"mobile":  true, // iOS support
		},
	}

	currentOS := runtime.GOOS
	supported, exists := compatibilityMatrix[currentOS][interfaceType]

	if !exists || !supported {
		result.Warnings = append(result.Warnings, types.ValidationItem{
			Code:      "INTERFACE_COMPATIBILITY_WARNING",
			Message:   fmt.Sprintf("Interface '%s' may have limited compatibility on %s", interfaceType, currentOS),
			Severity:  string(types.SeverityWarning),
			Category:  string(types.CategoryCompatibility),
			Fixable:   false,
			Timestamp: time.Now(),
		})
	}

	// Add interface-specific features
	switch interfaceType {
	case "cli":
		result.Environment.Features = append(result.Environment.Features, "cli_interface")
		result.Environment.Capabilities = append(result.Environment.Capabilities, "command_line")
	case "web":
		result.Environment.Features = append(result.Environment.Features, "web_interface")
		result.Environment.Capabilities = append(result.Environment.Capabilities, "browser_support")
	case "desktop":
		result.Environment.Features = append(result.Environment.Features, "desktop_interface")
		result.Environment.Capabilities = append(result.Environment.Capabilities, "native_app")
	case "mobile":
		result.Environment.Features = append(result.Environment.Features, "mobile_interface")
		result.Environment.Capabilities = append(result.Environment.Capabilities, "mobile_app")
	}

	return nil
}

// checkKnownIssues checks for known compatibility issues
func (cv *CompatibilityValidator) checkKnownIssues(req *types.ValidationRequest, result *types.ValidationResult) error {
	currentOS := runtime.GOOS
	interfaceType := req.Interface

	// Define known issues based on OS and interface combinations
	knownIssues := []types.KnownIssue{
		{
			ID:          "WINDOWS_POWERSHELL_VERSION",
			Description: "PowerShell version 5.1 or higher is required on Windows",
			Severity:    "medium",
			OS:          "windows",
			Version:     "all",
			Workaround:  "Install Windows Management Framework 5.1 or later",
			Fixed:       false,
		},
		{
			ID:          "LINUX_SYSTEMD_REQUIRED",
			Description: "Systemd is required for service management on Linux",
			Severity:    "high",
			OS:          "linux",
			Version:     "all",
			Workaround:  "Use a Linux distribution with systemd support",
			Fixed:       false,
		},
		{
			ID:          "DARWIN_XCODE_REQUIRED",
			Description: "Xcode Command Line Tools are required on macOS",
			Severity:    "medium",
			OS:          "darwin",
			Version:     "all",
			Workaround:  "Install Xcode Command Line Tools: xcode-select --install",
			Fixed:       false,
		},
		{
			ID:          "MOBILE_IOS_VERSION",
			Description: "iOS 12.0 or higher is required for mobile interface",
			Severity:    "medium",
			OS:          "darwin",
			Version:     "mobile",
			Workaround:  "Update iOS to version 12.0 or later",
			Fixed:       false,
		},
	}

	// Filter issues relevant to current environment
	var relevantIssues []types.KnownIssue
	for _, issue := range knownIssues {
		if issue.OS == currentOS || issue.OS == "all" {
			if issue.Version == "all" || issue.Version == interfaceType {
				relevantIssues = append(relevantIssues, issue)
			}
		}
	}

	result.Compatibility.KnownIssues = relevantIssues

	// Add validation items for relevant issues
	for _, issue := range relevantIssues {
		severity := types.SeverityWarning
		if issue.Severity == "high" {
			severity = types.SeverityError
		} else if issue.Severity == "critical" {
			severity = types.SeverityCritical
		}

		result.Warnings = append(result.Warnings, types.ValidationItem{
			Code:      issue.ID,
			Message:   issue.Description,
			Severity:  string(severity),
			Category:  string(types.CategoryCompatibility),
			Fixable:   issue.Workaround != "",
			Timestamp: time.Now(),
		})
	}

	return nil
}

// checkWorkarounds checks for available workarounds
func (cv *CompatibilityValidator) checkWorkarounds(req *types.ValidationRequest, result *types.ValidationResult) error {
	currentOS := runtime.GOOS
	interfaceType := req.Interface

	// Define workarounds based on OS and interface combinations
	workarounds := []types.Workaround{
		{
			ID:          "WINDOWS_POWERSHELL_UPDATE",
			Description: "Update PowerShell to version 5.1 or higher",
			Command:     "Install-Module -Name PowerShellGet -Force -AllowClobber",
			Script:      "powershell -Command \"Install-Module -Name PowerShellGet -Force\"",
			Manual:      "Download and install Windows Management Framework 5.1 from Microsoft",
			Risk:        "low",
			Reversible:  true,
		},
		{
			ID:          "LINUX_SYSTEMD_ENABLE",
			Description: "Enable systemd support",
			Command:     "systemctl --version",
			Script:      "systemctl status",
			Manual:      "Ensure systemd is installed and running",
			Risk:        "low",
			Reversible:  true,
		},
		{
			ID:          "DARWIN_XCODE_INSTALL",
			Description: "Install Xcode Command Line Tools",
			Command:     "xcode-select --install",
			Script:      "xcode-select --install",
			Manual:      "Run 'xcode-select --install' in Terminal",
			Risk:        "low",
			Reversible:  true,
		},
		{
			ID:          "MOBILE_IOS_UPDATE",
			Description: "Update iOS to supported version",
			Command:     "Settings > General > Software Update",
			Script:      "",
			Manual:      "Update iOS through Settings app",
			Risk:        "medium",
			Reversible:  false,
		},
	}

	// Filter workarounds relevant to current environment
	var relevantWorkarounds []types.Workaround
	for _, workaround := range workarounds {
		// This is a simplified filtering - in production, you'd have more sophisticated matching
		if strings.Contains(workaround.ID, strings.ToUpper(currentOS)) ||
			strings.Contains(workaround.ID, strings.ToUpper(interfaceType)) {
			relevantWorkarounds = append(relevantWorkarounds, workaround)
		}
	}

	result.Compatibility.Workarounds = relevantWorkarounds

	return nil
}
