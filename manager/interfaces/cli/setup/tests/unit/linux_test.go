//go:build linux
// +build linux

package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/internal/types"
)

// TestLinuxValidation tests the Linux environment validation
func TestLinuxValidation(t *testing.T) {
	result, err := setup.ValidateLinuxEnvironment(false)
	if err != nil {
		t.Fatalf("Linux validation failed with error: %v", err)
	}

	// Check that basic environment info is populated
	if result.Environment.OS == "" {
		t.Error("OS information is missing")
	}
	if result.Environment.OSVersion == "" {
		t.Error("OS version information is missing")
	}
	if result.Environment.Architecture == "" {
		t.Error("Architecture information is missing")
	}
	if result.Environment.HomeDir == "" {
		t.Error("Home directory information is missing")
	}
}

// TestLinuxSetup tests the Linux setup process
func TestLinuxSetup(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "syntropy-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Configure test options
	options := types.SetupOptions{
		Force:          true,
		InstallService: false,
		ConfigPath:     filepath.Join(tempDir, "config.yaml"),
		HomeDir:        tempDir,
	}

	// Test Setup
	setupResult, err := setup.Setup(options)
	if err != nil {
		t.Fatalf("Setup failed with error: %v", err)
	}

	if !setupResult.Success {
		t.Error("Setup was not successful")
	}

	// Verify config file was created
	configFile := setupResult.ConfigPath
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Errorf("Config file was not created: %s", configFile)
	}

	// Verify directory structure
	syntropyDir := filepath.Join(tempDir, ".syntropy")
	dirsToCheck := []string{
		filepath.Join(syntropyDir, "config"),
		filepath.Join(syntropyDir, "logs"),
		filepath.Join(syntropyDir, "data"),
		filepath.Join(syntropyDir, "bin"),
		filepath.Join(syntropyDir, "services"),
	}

	for _, dir := range dirsToCheck {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("Expected directory was not created: %s", dir)
		}
	}
}

// TestLinuxStatus tests the Linux status check
func TestLinuxStatus(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "syntropy-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Configure test options
	options := types.SetupOptions{
		Force:          true,
		InstallService: false,
		ConfigPath:     filepath.Join(tempDir, "config.yaml"),
		HomeDir:        tempDir,
	}

	// First run setup
	_, err = setup.Setup(options)
	if err != nil {
		t.Fatalf("Setup failed with error: %v", err)
	}

	// Then check status
	statusResult, err := setup.Status(options)
	if err != nil {
		t.Fatalf("Status check failed with error: %v", err)
	}

	if !statusResult.Success {
		t.Error("Status check was not successful")
	}

	if statusResult.ConfigPath != options.ConfigPath {
		t.Errorf("Expected config path %s, got %s", options.ConfigPath, statusResult.ConfigPath)
	}
}

// TestLinuxReset tests the Linux reset functionality
func TestLinuxReset(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "syntropy-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Configure test options
	options := types.SetupOptions{
		Force:          true,
		InstallService: false,
		ConfigPath:     filepath.Join(tempDir, "config.yaml"),
		HomeDir:        tempDir,
	}

	// First run setup
	_, err = setup.Setup(options)
	if err != nil {
		t.Fatalf("Setup failed with error: %v", err)
	}

	// Then reset
	resetResult, err := setup.Reset(options)
	if err != nil {
		t.Fatalf("Reset failed with error: %v", err)
	}

	if !resetResult.Success {
		t.Error("Reset was not successful")
	}

	// Verify syntropy directory was removed
	syntropyDir := filepath.Join(tempDir, ".syntropy")
	if _, err := os.Stat(syntropyDir); !os.IsNotExist(err) {
		t.Errorf("Syntropy directory was not removed after reset: %s", syntropyDir)
	}
}
