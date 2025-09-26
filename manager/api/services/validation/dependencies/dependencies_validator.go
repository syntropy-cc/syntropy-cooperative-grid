// Package dependencies provides dependency validation for the API
package dependencies

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/api/middleware"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/api/types"
)

// DependenciesValidator validates system dependencies
type DependenciesValidator struct {
	logger middleware.Logger
}

// NewDependenciesValidator creates a new dependencies validator
func NewDependenciesValidator(logger middleware.Logger) *DependenciesValidator {
	return &DependenciesValidator{
		logger: logger,
	}
}

// Validate performs dependency validation
func (dv *DependenciesValidator) Validate(req *types.ValidationRequest, result *types.ValidationResult) error {
	dv.logger.Info("Starting dependencies validation", map[string]interface{}{
		"interface": req.Interface,
		"os":        runtime.GOOS,
	})

	// Get dependencies based on operating system
	dependencies := dv.getDependenciesForOS(runtime.GOOS)

	// Validate each dependency
	for _, dep := range dependencies {
		if err := dv.validateDependency(dep, result); err != nil {
			dv.logger.Error("Dependency validation failed", map[string]interface{}{
				"dependency": dep.Name,
				"error":      err.Error(),
			})
		}
	}

	// Check for optional dependencies
	optionalDeps := dv.getOptionalDependencies(runtime.GOOS)
	for _, dep := range optionalDeps {
		dv.checkOptionalDependency(dep, result)
	}

	dv.logger.Info("Dependencies validation completed", map[string]interface{}{
		"total_dependencies":    len(dependencies),
		"optional_dependencies": len(optionalDeps),
	})

	return nil
}

// getDependenciesForOS returns required dependencies for the specified OS
func (dv *DependenciesValidator) getDependenciesForOS(os string) []types.Dependency {
	switch os {
	case "windows":
		return dv.getWindowsDependencies()
	case "linux":
		return dv.getLinuxDependencies()
	case "darwin":
		return dv.getDarwinDependencies()
	default:
		return []types.Dependency{}
	}
}

// getWindowsDependencies returns Windows-specific dependencies
func (dv *DependenciesValidator) getWindowsDependencies() []types.Dependency {
	return []types.Dependency{
		{
			Name:        "PowerShell",
			Version:     "5.1.0",
			Status:      "unknown",
			Required:    true,
			Installable: true,
			InstallCmd:  "Install PowerShell from Microsoft Store or download from Microsoft",
			CheckCmd:    "powershell -Command \"$PSVersionTable.PSVersion\"",
			Path:        "powershell.exe",
		},
		{
			Name:        "Windows Management Framework",
			Version:     "5.1.0",
			Status:      "unknown",
			Required:    true,
			Installable: true,
			InstallCmd:  "Download and install WMF 5.1 from Microsoft",
			CheckCmd:    "powershell -Command \"Get-Host\"",
			Path:        "powershell.exe",
		},
		{
			Name:        "Microsoft Visual C++ Redistributable",
			Version:     "2019",
			Status:      "unknown",
			Required:    true,
			Installable: true,
			InstallCmd:  "Download and install Visual C++ Redistributable from Microsoft",
			CheckCmd:    "reg query \"HKEY_LOCAL_MACHINE\\SOFTWARE\\Microsoft\\VisualStudio\\14.0\\VC\\Runtimes\\x64\"",
			Path:        "",
		},
	}
}

// getLinuxDependencies returns Linux-specific dependencies
func (dv *DependenciesValidator) getLinuxDependencies() []types.Dependency {
	return []types.Dependency{
		{
			Name:        "systemd",
			Version:     "230",
			Status:      "unknown",
			Required:    true,
			Installable: true,
			InstallCmd:  "sudo apt-get install systemd (Ubuntu/Debian) or sudo yum install systemd (RHEL/CentOS)",
			CheckCmd:    "systemctl --version",
			Path:        "/bin/systemctl",
		},
		{
			Name:        "curl",
			Version:     "7.0",
			Status:      "unknown",
			Required:    true,
			Installable: true,
			InstallCmd:  "sudo apt-get install curl (Ubuntu/Debian) or sudo yum install curl (RHEL/CentOS)",
			CheckCmd:    "curl --version",
			Path:        "/usr/bin/curl",
		},
		{
			Name:        "wget",
			Version:     "1.0",
			Status:      "unknown",
			Required:    true,
			Installable: true,
			InstallCmd:  "sudo apt-get install wget (Ubuntu/Debian) or sudo yum install wget (RHEL/CentOS)",
			CheckCmd:    "wget --version",
			Path:        "/usr/bin/wget",
		},
		{
			Name:        "openssl",
			Version:     "1.1.0",
			Status:      "unknown",
			Required:    true,
			Installable: true,
			InstallCmd:  "sudo apt-get install openssl (Ubuntu/Debian) or sudo yum install openssl (RHEL/CentOS)",
			CheckCmd:    "openssl version",
			Path:        "/usr/bin/openssl",
		},
	}
}

// getDarwinDependencies returns macOS-specific dependencies
func (dv *DependenciesValidator) getDarwinDependencies() []types.Dependency {
	return []types.Dependency{
		{
			Name:        "Xcode Command Line Tools",
			Version:     "12.0",
			Status:      "unknown",
			Required:    true,
			Installable: true,
			InstallCmd:  "xcode-select --install",
			CheckCmd:    "xcode-select -p",
			Path:        "/usr/bin/xcode-select",
		},
		{
			Name:        "Homebrew",
			Version:     "3.0",
			Status:      "unknown",
			Required:    false,
			Installable: true,
			InstallCmd:  "/bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\"",
			CheckCmd:    "brew --version",
			Path:        "/usr/local/bin/brew",
		},
		{
			Name:        "curl",
			Version:     "7.0",
			Status:      "unknown",
			Required:    true,
			Installable: true,
			InstallCmd:  "brew install curl",
			CheckCmd:    "curl --version",
			Path:        "/usr/bin/curl",
		},
	}
}

// getOptionalDependencies returns optional dependencies
func (dv *DependenciesValidator) getOptionalDependencies(os string) []types.Dependency {
	return []types.Dependency{
		{
			Name:        "Docker",
			Version:     "20.0",
			Status:      "unknown",
			Required:    false,
			Installable: true,
			InstallCmd:  "Install Docker Desktop",
			CheckCmd:    "docker --version",
			Path:        "docker",
		},
		{
			Name:        "Git",
			Version:     "2.0",
			Status:      "unknown",
			Required:    false,
			Installable: true,
			InstallCmd:  "Install Git from official website",
			CheckCmd:    "git --version",
			Path:        "git",
		},
		{
			Name:        "Node.js",
			Version:     "16.0",
			Status:      "unknown",
			Required:    false,
			Installable: true,
			InstallCmd:  "Install Node.js from official website",
			CheckCmd:    "node --version",
			Path:        "node",
		},
	}
}

// validateDependency validates a single dependency
func (dv *DependenciesValidator) validateDependency(dep types.Dependency, result *types.ValidationResult) error {
	dv.logger.Debug("Validating dependency", map[string]interface{}{
		"name":     dep.Name,
		"required": dep.Required,
	})

	// Check if dependency is installed
	installed, version, err := dv.checkDependencyInstalled(dep)
	if err != nil {
		dv.logger.Error("Failed to check dependency", map[string]interface{}{
			"dependency": dep.Name,
			"error":      err.Error(),
		})

		// Update dependency status
		dep.Status = "error"
		dep.Current = "unknown"
	} else if installed {
		dep.Status = "installed"
		dep.Current = version

		// Check if version meets requirements
		if !dv.isVersionCompatible(version, dep.Version) {
			dep.Status = "outdated"
			result.Warnings = append(result.Warnings, types.ValidationItem{
				Code:     "DEPENDENCY_OUTDATED",
				Message:  fmt.Sprintf("Dependency '%s' version %s is outdated, required: %s", dep.Name, version, dep.Version),
				Severity: string(types.SeverityWarning),
				Category: string(types.CategoryDependencies),
				Field:    dep.Name,
				Expected: dep.Version,
				Actual:   version,
				Fixable:  dep.Installable,
				AutoFix: &types.AutoFixInfo{
					Available: dep.Installable,
					Command:   dep.InstallCmd,
					Manual:    fmt.Sprintf("Update %s to version %s or later", dep.Name, dep.Version),
					Risk:      "medium",
				},
				Timestamp: time.Now(),
			})
		} else {
			result.Environment.Features = append(result.Environment.Features, fmt.Sprintf("%s_support", strings.ToLower(dep.Name)))
		}
	} else {
		dep.Status = "missing"
		dep.Current = "not installed"

		if dep.Required {
			result.Errors = append(result.Errors, types.ValidationItem{
				Code:     "REQUIRED_DEPENDENCY_MISSING",
				Message:  fmt.Sprintf("Required dependency '%s' is not installed", dep.Name),
				Severity: string(types.SeverityError),
				Category: string(types.CategoryDependencies),
				Field:    dep.Name,
				Expected: fmt.Sprintf("installed (%s)", dep.Version),
				Actual:   "not installed",
				Fixable:  dep.Installable,
				AutoFix: &types.AutoFixInfo{
					Available: dep.Installable,
					Command:   dep.InstallCmd,
					Manual:    fmt.Sprintf("Install %s version %s or later", dep.Name, dep.Version),
					Risk:      "low",
				},
				Timestamp: time.Now(),
			})
		} else {
			result.Warnings = append(result.Warnings, types.ValidationItem{
				Code:     "OPTIONAL_DEPENDENCY_MISSING",
				Message:  fmt.Sprintf("Optional dependency '%s' is not installed", dep.Name),
				Severity: string(types.SeverityInfo),
				Category: string(types.CategoryDependencies),
				Field:    dep.Name,
				Expected: fmt.Sprintf("installed (%s)", dep.Version),
				Actual:   "not installed",
				Fixable:  dep.Installable,
				AutoFix: &types.AutoFixInfo{
					Available: dep.Installable,
					Command:   dep.InstallCmd,
					Manual:    fmt.Sprintf("Install %s version %s or later for enhanced functionality", dep.Name, dep.Version),
					Risk:      "low",
				},
				Timestamp: time.Now(),
			})
		}
	}

	// Add dependency to compatibility result
	result.Compatibility.Dependencies = append(result.Compatibility.Dependencies, dep)

	return nil
}

// checkOptionalDependency checks an optional dependency
func (dv *DependenciesValidator) checkOptionalDependency(dep types.Dependency, result *types.ValidationResult) {
	installed, version, err := dv.checkDependencyInstalled(dep)
	if err != nil {
		dep.Status = "error"
		dep.Current = "unknown"
	} else if installed {
		dep.Status = "installed"
		dep.Current = version
		result.Environment.Features = append(result.Environment.Features, fmt.Sprintf("%s_support", strings.ToLower(dep.Name)))
	} else {
		dep.Status = "missing"
		dep.Current = "not installed"
	}

	// Add to optional features
	result.Compatibility.OptionalFeatures = append(result.Compatibility.OptionalFeatures, types.Feature{
		Name:        dep.Name,
		Description: fmt.Sprintf("Optional dependency: %s", dep.Name),
		Available:   installed,
		Required:    false,
		Version:     version,
		Enabled:     installed,
	})
}

// checkDependencyInstalled checks if a dependency is installed
func (dv *DependenciesValidator) checkDependencyInstalled(dep types.Dependency) (bool, string, error) {
	// First, check if the command exists in PATH
	if dep.CheckCmd != "" {
		output, err := exec.Command("sh", "-c", dep.CheckCmd).Output()
		if err != nil {
			// Try alternative check methods
			return dv.checkDependencyAlternative(dep)
		}

		// Parse version from output
		version := dv.extractVersionFromOutput(string(output), dep.Name)
		return true, version, nil
	}

	// Check by path
	if dep.Path != "" {
		if _, err := os.Stat(dep.Path); err == nil {
			return true, "unknown", nil
		}
	}

	return false, "", nil
}

// checkDependencyAlternative checks dependency using alternative methods
func (dv *DependenciesValidator) checkDependencyAlternative(dep types.Dependency) (bool, string, error) {
	// Try to execute the command directly
	if dep.Path != "" {
		cmd := exec.Command(dep.Path, "--version")
		output, err := cmd.Output()
		if err == nil {
			version := dv.extractVersionFromOutput(string(output), dep.Name)
			return true, version, nil
		}
	}

	// Try common version flags
	versionFlags := []string{"--version", "-version", "-v", "--help"}
	for _, flag := range versionFlags {
		cmd := exec.Command(dep.Name, flag)
		output, err := cmd.Output()
		if err == nil {
			version := dv.extractVersionFromOutput(string(output), dep.Name)
			return true, version, nil
		}
	}

	return false, "", nil
}

// extractVersionFromOutput extracts version information from command output
func (dv *DependenciesValidator) extractVersionFromOutput(output, dependency string) string {
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Look for version patterns
		if strings.Contains(strings.ToLower(line), "version") {
			// Extract version number
			parts := strings.Fields(line)
			for _, part := range parts {
				if dv.isVersionNumber(part) {
					return part
				}
			}
		}

		// Look for version at the beginning of line
		if dv.isVersionNumber(line) {
			return line
		}
	}

	return "unknown"
}

// isVersionNumber checks if a string looks like a version number
func (dv *DependenciesValidator) isVersionNumber(s string) bool {
	// Remove common prefixes/suffixes
	s = strings.TrimPrefix(s, "v")
	s = strings.TrimSuffix(s, ",")
	s = strings.TrimSuffix(s, ";")

	// Check if it contains version-like patterns
	if strings.Contains(s, ".") && len(s) >= 3 {
		return true
	}

	// Check if it's a simple number
	for _, char := range s {
		if char < '0' || char > '9' {
			return false
		}
	}

	return len(s) > 0
}

// isVersionCompatible checks if a version meets the minimum requirement
func (dv *DependenciesValidator) isVersionCompatible(current, required string) bool {
	// This is a simplified version comparison
	// In production, you would use proper semantic versioning

	// If version is unknown, assume it's compatible
	if current == "unknown" {
		return true
	}

	// Simple string comparison for now
	// In production, parse and compare semantic versions
	return strings.Compare(current, required) >= 0
}
