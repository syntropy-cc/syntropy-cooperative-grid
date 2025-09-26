// Package helpers provides common test utilities and helper functions
package helpers

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/tests/types"
)

// TestEnvironment provides a controlled test environment
type TestEnvironment struct {
	TempDir    string
	ConfigPath string
	HomeDir    string
	Cleanup    func()
}

// SetupTestEnvironment creates a temporary test environment
func SetupTestEnvironment(t *testing.T) *TestEnvironment {
	t.Helper()
	
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "syntropy-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	
	// Create subdirectories
	homeDir := filepath.Join(tempDir, "home")
	configDir := filepath.Join(homeDir, ".syntropy")
	configPath := filepath.Join(configDir, "config.yaml")
	
	if err := os.MkdirAll(configDir, 0755); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to create config directory: %v", err)
	}
	
	// Cleanup function
	cleanup := func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Warning: Failed to cleanup temp directory %s: %v", tempDir, err)
		}
	}
	
	return &TestEnvironment{
		TempDir:    tempDir,
		ConfigPath: configPath,
		HomeDir:    homeDir,
		Cleanup:    cleanup,
	}
}

// CreateTestConfig creates a test configuration file
func CreateTestConfig(t *testing.T, path string, config types.SetupConfig) {
	t.Helper()
	
	content := `manager:
  home_dir: ` + config.Manager.HomeDir + `
  log_level: ` + config.Manager.LogLevel + `
  api_endpoint: ` + config.Manager.APIEndpoint + `
environment:
  os: ` + config.Environment.OS + `
  architecture: ` + config.Environment.Architecture + `
  home_dir: ` + config.Environment.HomeDir + `
`
	
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}
}

// CreateInvalidConfig creates an invalid configuration file for testing
func CreateInvalidConfig(t *testing.T, path string) {
	t.Helper()
	
	content := `invalid_yaml:
  - this is not valid
  - yaml: [unclosed bracket
`
	
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create invalid config: %v", err)
	}
}

// GetTestSetupOptions returns test setup options
func GetTestSetupOptions() types.SetupOptions {
	return types.SetupOptions{
		Force:          false,
		InstallService: false,
		ConfigPath:     "",
		HomeDir:        "",
	}
}

// GetTestSetupConfig returns a test setup configuration
func GetTestSetupConfig(homeDir string) types.SetupConfig {
	return types.SetupConfig{
		Manager: types.ManagerConfig{
			HomeDir:     homeDir,
			LogLevel:    "info",
			APIEndpoint: "https://api.syntropy.com",
		},
		Environment: types.Environment{
			OS:           runtime.GOOS,
			Architecture: runtime.GOARCH,
			HomeDir:      homeDir,
		},
	}
}

// GetTestEnvironmentInfo returns test environment information
func GetTestEnvironmentInfo() types.EnvironmentInfo {
	return types.EnvironmentInfo{
		OS:              runtime.GOOS,
		Architecture:    runtime.GOARCH,
		HasAdminRights:  false,
		AvailableDiskGB: 100.0,
		HasInternet:     true,
		HomeDir:         "/tmp/test",
		SystemResources: types.SystemResources{
			TotalMemoryGB: 16.0,
			CPUCores:      8,
		},
	}
}

// AssertSetupResult validates a setup result
func AssertSetupResult(t *testing.T, result types.SetupResult, expectSuccess bool) {
	t.Helper()
	
	if result.Success != expectSuccess {
		t.Errorf("Expected success=%v, got %v", expectSuccess, result.Success)
	}
	
	if result.StartTime.IsZero() {
		t.Error("StartTime should not be zero")
	}
	
	if result.EndTime.IsZero() {
		t.Error("EndTime should not be zero")
	}
	
	if result.EndTime.Before(result.StartTime) {
		t.Error("EndTime should be after StartTime")
	}
	
	if expectSuccess && result.Error != nil {
		t.Errorf("Expected no error for successful setup, got: %v", result.Error)
	}
	
	if !expectSuccess && result.Error == nil {
		t.Error("Expected error for failed setup, got nil")
	}
}

// AssertValidationResult validates a validation result
func AssertValidationResult(t *testing.T, result types.ValidationResult, expectValid bool) {
	t.Helper()
	
	if result.Valid != expectValid {
		t.Errorf("Expected valid=%v, got %v", expectValid, result.Valid)
	}
	
	if !expectValid && len(result.Errors) == 0 {
		t.Error("Expected errors for invalid validation result")
	}
	
	if expectValid && len(result.Errors) > 0 {
		t.Errorf("Expected no errors for valid validation result, got: %v", result.Errors)
	}
}

// WithTimeout creates a context with timeout for testing
func WithTimeout(t *testing.T, timeout time.Duration) (context.Context, context.CancelFunc) {
	t.Helper()
	return context.WithTimeout(context.Background(), timeout)
}

// WithCancel creates a cancellable context for testing
func WithCancel(t *testing.T) (context.Context, context.CancelFunc) {
	t.Helper()
	return context.WithCancel(context.Background())
}

// SkipIfNotLinux skips the test if not running on Linux
func SkipIfNotLinux(t *testing.T) {
	t.Helper()
	if runtime.GOOS != "linux" {
		t.Skip("Skipping test: only supported on Linux")
	}
}

// SkipIfNotWindows skips the test if not running on Windows
func SkipIfNotWindows(t *testing.T) {
	t.Helper()
	if runtime.GOOS != "windows" {
		t.Skip("Skipping test: only supported on Windows")
	}
}

// SkipIfNotDarwin skips the test if not running on macOS
func SkipIfNotDarwin(t *testing.T) {
	t.Helper()
	if runtime.GOOS != "darwin" {
		t.Skip("Skipping test: only supported on macOS")
	}
}

// RequireAdmin skips the test if not running with admin privileges
func RequireAdmin(t *testing.T) {
	t.Helper()
	if !IsAdmin() {
		t.Skip("Skipping test: requires admin privileges")
	}
}

// IsAdmin checks if running with admin privileges
func IsAdmin() bool {
	switch runtime.GOOS {
	case "windows":
		// On Windows, check if we can write to system directory
		testFile := filepath.Join(os.Getenv("WINDIR"), "temp", "admin_test")
		if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
			return false
		}
		os.Remove(testFile)
		return true
	case "linux", "darwin":
		// On Unix-like systems, check if UID is 0
		return os.Getuid() == 0
	default:
		return false
	}
}

// CreateLargeFile creates a large file for testing disk space scenarios
func CreateLargeFile(t *testing.T, path string, sizeGB float64) {
	t.Helper()
	
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create large file: %v", err)
	}
	defer file.Close()
	
	// Write in chunks to avoid memory issues
	chunkSize := 1024 * 1024 // 1MB chunks
	totalBytes := int64(sizeGB * 1024 * 1024 * 1024)
	chunk := make([]byte, chunkSize)
	
	for written := int64(0); written < totalBytes; written += int64(chunkSize) {
		remaining := totalBytes - written
		if remaining < int64(chunkSize) {
			chunk = chunk[:remaining]
		}
		
		if _, err := file.Write(chunk); err != nil {
			t.Fatalf("Failed to write to large file: %v", err)
		}
	}
}

// WaitForCondition waits for a condition to be true with timeout
func WaitForCondition(t *testing.T, condition func() bool, timeout time.Duration, message string) {
	t.Helper()
	
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			t.Fatalf("Timeout waiting for condition: %s", message)
		case <-ticker.C:
			if condition() {
				return
			}
		}
	}
}

// AssertFileExists checks if a file exists
func AssertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("Expected file to exist: %s", path)
	}
}

// AssertFileNotExists checks if a file does not exist
func AssertFileNotExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err == nil {
		t.Errorf("Expected file to not exist: %s", path)
	}
}

// AssertDirectoryExists checks if a directory exists
func AssertDirectoryExists(t *testing.T, path string) {
	t.Helper()
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		t.Errorf("Expected directory to exist: %s", path)
		return
	}
	if !info.IsDir() {
		t.Errorf("Expected %s to be a directory", path)
	}
}

// AssertFileContains checks if a file contains specific content
func AssertFileContains(t *testing.T, path, content string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", path, err)
	}
	
	if !strings.Contains(string(data), content) {
		t.Errorf("File %s does not contain expected content: %s", path, content)
	}
}

// AssertFilePermissions checks file permissions
func AssertFilePermissions(t *testing.T, path string, expectedPerm os.FileMode) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Failed to stat file %s: %v", path, err)
	}
	
	actualPerm := info.Mode().Perm()
	if actualPerm != expectedPerm {
		t.Errorf("File %s has permissions %o, expected %o", path, actualPerm, expectedPerm)
	}
}

// GetAvailableDiskSpace returns available disk space in GB
func GetAvailableDiskSpace(t *testing.T, path string) float64 {
	t.Helper()
	
	// This is a simplified implementation
	// In a real scenario, you'd use syscalls to get actual disk space
	return 100.0 // Default to 100GB for testing
}

// SimulateNetworkDelay simulates network delay for testing
func SimulateNetworkDelay(delay time.Duration) {
	time.Sleep(delay)
}

// GenerateTestData generates test data of specified size
func GenerateTestData(size int) []byte {
	data := make([]byte, size)
	for i := range data {
		data[i] = byte(i % 256)
	}
	return data
}