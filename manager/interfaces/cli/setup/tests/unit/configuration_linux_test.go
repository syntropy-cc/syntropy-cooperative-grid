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

func TestConfigureLinuxEnvironment(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "syntropy-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock validation result
	validationResult := &types.ValidationResult{
		Valid: true,
		Environment: types.EnvironmentInfo{
			OS:           "linux",
			Architecture: "amd64",
			HomeDir:      tempDir,
		},
	}

	// Create setup options with the temp directory
	options := types.SetupOptions{
		HomeDir:    tempDir,
		ConfigPath: filepath.Join(tempDir, ".syntropy", "config", "test-manager.yaml"),
	}

	// Test configuration
	err = setup.ConfigureLinuxEnvironment(validationResult, options)
	if err != nil {
		t.Fatalf("ConfigureLinuxEnvironment failed: %v", err)
	}

	// Check if directories were created
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
			t.Errorf("Expected directory %s to exist", dir)
		}
	}

	// Check if config file was created
	configPath := filepath.Join(syntropyDir, "config", "test-manager.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("Expected config file %s to exist", configPath)
	}
}
