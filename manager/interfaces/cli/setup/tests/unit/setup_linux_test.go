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

func TestStatusLinux(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "syntropy-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create the .syntropy directory structure
	syntropyDir := filepath.Join(tempDir, ".syntropy")
	configDir := filepath.Join(syntropyDir, "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create syntropy config directory: %v", err)
	}

	// Create a dummy config file
	configPath := filepath.Join(configDir, "manager.yaml")
	if err := os.WriteFile(configPath, []byte("test: true"), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Set the HOME environment variable temporarily for this test
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Test status with custom home directory
	options := types.SetupOptions{
		HomeDir:    tempDir,
		ConfigPath: configPath,
	}

	result, err := setup.Status(options)
	if err != nil {
		t.Fatalf("Status check failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected status check to succeed")
	}

	if result.ConfigPath != configPath {
		t.Errorf("Expected config path to be %s, got %s", configPath, result.ConfigPath)
	}
}

func TestResetLinux(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "syntropy-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create the .syntropy directory structure
	syntropyDir := filepath.Join(tempDir, ".syntropy")
	if err := os.MkdirAll(syntropyDir, 0755); err != nil {
		t.Fatalf("Failed to create syntropy directory: %v", err)
	}

	// Test reset with force option
	options := types.SetupOptions{
		HomeDir: tempDir,
		Force:   true,
	}

	// Set the HOME environment variable temporarily for this test
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	result, err := setup.Reset(options)
	if err != nil {
		t.Fatalf("Reset failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected reset to succeed")
	}

	// Check if directory was removed
	if _, err := os.Stat(syntropyDir); !os.IsNotExist(err) {
		t.Error("Expected syntropy directory to be removed")
	}
}
