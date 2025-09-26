package fixtures

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/tests/types"
)

// TestDataManager manages test fixtures and data
type TestDataManager struct {
	BaseDir string
}

// NewTestDataManager creates a new test data manager
func NewTestDataManager(baseDir string) *TestDataManager {
	return &TestDataManager{
		BaseDir: baseDir,
	}
}

// ValidConfigurations returns valid configuration test data
func (tdm *TestDataManager) ValidConfigurations() []ConfigurationFixture {
	return []ConfigurationFixture{
		{
			Name:        "minimal_config",
			Description: "Minimal valid configuration",
			Content: `manager:
  home_dir: /home/user/.syntropy
  log_level: info
  api_endpoint: https://api.syntropy.com
`,
			Expected: types.SetupConfig{
				Manager: types.ManagerConfig{
					HomeDir:     "/home/user/.syntropy",
					LogLevel:    "info",
					APIEndpoint: "https://api.syntropy.com",
				},
			},
		},
		{
			Name:        "complete_config",
			Description: "Complete configuration with all options",
			Content: `manager:
  home_dir: /home/user/.syntropy
  log_level: debug
  api_endpoint: https://api.syntropy.com
  timeout: 30s
  retry_attempts: 3

owner_key:
  type: ed25519
  path: /home/user/.syntropy/keys/owner.key
  public_key: "ed25519_public_key_here"

environment:
  os: linux
  architecture: amd64
  home_dir: /home/user/.syntropy
`,
			Expected: types.SetupConfig{
				Manager: types.ManagerConfig{
					HomeDir:     "/home/user/.syntropy",
					LogLevel:    "debug",
					APIEndpoint: "https://api.syntropy.com",
				},
				OwnerKey: types.OwnerKey{
					Type:      "ed25519",
					Path:      "/home/user/.syntropy/keys/owner.key",
					PublicKey: "ed25519_public_key_here",
				},
				Environment: types.Environment{
					OS:           "linux",
					Architecture: "amd64",
					HomeDir:      "/home/user/.syntropy",
				},
			},
		},
		{
			Name:        "production_config",
			Description: "Production-ready configuration",
			Content: `manager:
  home_dir: /opt/syntropy
  log_level: warn
  api_endpoint: https://prod-api.syntropy.com
  timeout: 60s
  retry_attempts: 5
  enable_metrics: true
  enable_tracing: false

owner_key:
  type: ed25519
  path: /opt/syntropy/keys/production.key

environment:
  os: ` + runtime.GOOS + `
  architecture: ` + runtime.GOARCH + `
  home_dir: /opt/syntropy
`,
			Expected: types.SetupConfig{
				Manager: types.ManagerConfig{
					HomeDir:     "/opt/syntropy",
					LogLevel:    "warn",
					APIEndpoint: "https://prod-api.syntropy.com",
				},
				OwnerKey: types.OwnerKey{
					Type: "ed25519",
					Path: "/opt/syntropy/keys/production.key",
				},
				Environment: types.Environment{
					OS:           runtime.GOOS,
					Architecture: runtime.GOARCH,
					HomeDir:      "/opt/syntropy",
				},
			},
		},
	}
}

// InvalidConfigurations returns invalid configuration test data
func (tdm *TestDataManager) InvalidConfigurations() []ConfigurationFixture {
	return []ConfigurationFixture{
		{
			Name:        "empty_config",
			Description: "Empty configuration file",
			Content:     "",
			ExpectedError: "empty configuration",
		},
		{
			Name:        "invalid_yaml",
			Description: "Invalid YAML syntax",
			Content: `manager:
  home_dir: /home/user
  log_level: info
  invalid: yaml: syntax: [
`,
			ExpectedError: "yaml syntax error",
		},
		{
			Name:        "missing_required_fields",
			Description: "Missing required configuration fields",
			Content: `manager:
  log_level: info
`,
			ExpectedError: "missing required field",
		},
		{
			Name:        "invalid_log_level",
			Description: "Invalid log level value",
			Content: `manager:
  home_dir: /home/user/.syntropy
  log_level: invalid_level
  api_endpoint: https://api.syntropy.com
`,
			ExpectedError: "invalid log level",
		},
		{
			Name:        "invalid_api_endpoint",
			Description: "Invalid API endpoint URL",
			Content: `manager:
  home_dir: /home/user/.syntropy
  log_level: info
  api_endpoint: not-a-valid-url
`,
			ExpectedError: "invalid API endpoint",
		},
		{
			Name:        "malicious_path_traversal",
			Description: "Path traversal attempt in configuration",
			Content: `manager:
  home_dir: ../../../etc/passwd
  log_level: info
  api_endpoint: https://api.syntropy.com
`,
			ExpectedError: "invalid path",
		},
	}
}

// EdgeCaseConfigurations returns edge case configuration test data
func (tdm *TestDataManager) EdgeCaseConfigurations() []ConfigurationFixture {
	return []ConfigurationFixture{
		{
			Name:        "very_long_paths",
			Description: "Configuration with very long file paths",
			Content: fmt.Sprintf(`manager:
  home_dir: %s
  log_level: info
  api_endpoint: https://api.syntropy.com
`, strings.Repeat("/very/long/path", 50)),
			ExpectedError: "path too long",
		},
		{
			Name:        "unicode_characters",
			Description: "Configuration with Unicode characters",
			Content: `manager:
  home_dir: /home/用户/.syntropy
  log_level: info
  api_endpoint: https://api.syntropy.com
  description: "Configuration with émojis 🚀 and ñoñó"
`,
		},
		{
			Name:        "large_config_file",
			Description: "Very large configuration file",
			Content:     generateLargeConfig(1000),
		},
		{
			Name:        "windows_paths",
			Description: "Windows-style paths",
			Content: `manager:
  home_dir: C:\Users\User\AppData\Roaming\Syntropy
  log_level: info
  api_endpoint: https://api.syntropy.com
`,
		},
		{
			Name:        "special_characters",
			Description: "Configuration with special characters",
			Content: `manager:
  home_dir: "/home/user with spaces/.syntropy"
  log_level: info
  api_endpoint: "https://api.syntropy.com"
  special_field: "Value with !@#$%^&*()_+-={}[]|\\:;\"'<>?,./"
`,
		},
	}
}

// SetupOptionsFixtures returns various setup options for testing
func (tdm *TestDataManager) SetupOptionsFixtures() []SetupOptionsFixture {
	return []SetupOptionsFixture{
		{
			Name:        "default_options",
			Description: "Default setup options",
			Options: types.SetupOptions{
				Force:          false,
				InstallService: false,
				ConfigPath:     "",
				HomeDir:        "",
			},
		},
		{
			Name:        "force_setup",
			Description: "Setup with force option enabled",
			Options: types.SetupOptions{
				Force:          true,
				InstallService: false,
				ConfigPath:     "",
				HomeDir:        "",
			},
		},
		{
			Name:        "service_installation",
			Description: "Setup with service installation",
			Options: types.SetupOptions{
				Force:          false,
				InstallService: true,
				ConfigPath:     "",
				HomeDir:        "",
			},
		},
		{
			Name:        "custom_paths",
			Description: "Setup with custom configuration and home paths",
			Options: types.SetupOptions{
				Force:          false,
				InstallService: false,
				ConfigPath:     "/custom/path/config.yaml",
				HomeDir:        "/custom/home/dir",
			},
		},
		{
			Name:        "all_options_enabled",
			Description: "Setup with all options enabled",
			Options: types.SetupOptions{
				Force:          true,
				InstallService: true,
				ConfigPath:     "/opt/syntropy/config.yaml",
				HomeDir:        "/opt/syntropy",
			},
		},
	}
}

// ValidationResultFixtures returns validation result test data
func (tdm *TestDataManager) ValidationResultFixtures() []ValidationResultFixture {
	return []ValidationResultFixture{
		{
			Name:        "valid_environment",
			Description: "Valid environment validation result",
			Result: types.ValidationResult{
				Valid:    true,
				Warnings: []string{},
				Errors:   []string{},
				Environment: types.EnvironmentInfo{
					OS:              runtime.GOOS,
					Architecture:    runtime.GOARCH,
					HasAdminRights:  false,
					AvailableDiskGB: 100.0,
					HasInternet:     true,
					HomeDir:         "/home/user",
					SystemResources: types.SystemResources{
						TotalMemoryGB: 16.0,
						CPUCores:      8,
					},
				},
			},
		},
		{
			Name:        "validation_with_warnings",
			Description: "Validation result with warnings",
			Result: types.ValidationResult{
				Valid: true,
				Warnings: []string{
					"Low disk space available",
					"Slow internet connection detected",
				},
				Errors: []string{},
				Environment: types.EnvironmentInfo{
					OS:              runtime.GOOS,
					Architecture:    runtime.GOARCH,
					HasAdminRights:  false,
					AvailableDiskGB: 5.0,
					HasInternet:     true,
					HomeDir:         "/home/user",
				},
			},
		},
		{
			Name:        "validation_failure",
			Description: "Failed validation result",
			Result: types.ValidationResult{
				Valid:    false,
				Warnings: []string{},
				Errors: []string{
					"Insufficient disk space",
					"Missing required dependencies",
					"No internet connection",
				},
				Environment: types.EnvironmentInfo{
					OS:              runtime.GOOS,
					Architecture:    runtime.GOARCH,
					HasAdminRights:  false,
					AvailableDiskGB: 0.5,
					HasInternet:     false,
					HomeDir:         "/home/user",
				},
			},
		},
	}
}

// SetupResultFixtures returns setup result test data
func (tdm *TestDataManager) SetupResultFixtures() []SetupResultFixture {
	now := time.Now()
	return []SetupResultFixture{
		{
			Name:        "successful_setup",
			Description: "Successful setup result",
			Result: types.SetupResult{
				Success:     true,
				StartTime:   now,
				EndTime:     now.Add(5 * time.Second),
				ConfigPath:  "/home/user/.syntropy/config.yaml",
				Environment: runtime.GOOS,
				Options: types.SetupOptions{
					Force:          false,
					InstallService: false,
					HomeDir:        "/home/user/.syntropy",
				},
				Message: "Setup completed successfully",
			},
		},
		{
			Name:        "failed_setup",
			Description: "Failed setup result",
			Result: types.SetupResult{
				Success:     false,
				StartTime:   now,
				EndTime:     now.Add(2 * time.Second),
				ConfigPath:  "",
				Environment: runtime.GOOS,
				Options: types.SetupOptions{
					Force:          false,
					InstallService: true,
					HomeDir:        "/invalid/path",
				},
				Error:   fmt.Errorf("permission denied"),
				Message: "Setup failed due to insufficient permissions",
			},
		},
		{
			Name:        "setup_with_service",
			Description: "Successful setup with service installation",
			Result: types.SetupResult{
				Success:     true,
				StartTime:   now,
				EndTime:     now.Add(10 * time.Second),
				ConfigPath:  "/opt/syntropy/config.yaml",
				Environment: runtime.GOOS,
				Options: types.SetupOptions{
					Force:          false,
					InstallService: true,
					HomeDir:        "/opt/syntropy",
				},
				Message: "Setup completed with service installation",
			},
		},
	}
}

// CreateTestFiles creates test files in the specified directory
func (tdm *TestDataManager) CreateTestFiles(dir string) error {
	testFiles := map[string]string{
		"valid-config.yaml": `manager:
  home_dir: ` + dir + `
  log_level: info
  api_endpoint: https://api.syntropy.com
`,
		"invalid-config.yaml": `invalid: yaml: content: [`,
		"empty-config.yaml":   "",
		"large-config.yaml":   generateLargeConfig(100),
		"test-key.pem": `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC7VJTUt9Us8cKB
wQNneCjmrSk+12MwZ/wjQjbDHWoI9DjdFHnTbg5AjuHSPiQyPArvQqBdlVyLDXc
-----END PRIVATE KEY-----`,
		"test-data.json": `{
  "test": true,
  "data": [1, 2, 3],
  "nested": {
    "value": "test"
  }
}`,
	}

	for filename, content := range testFiles {
		filePath := filepath.Join(dir, filename)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create test file %s: %w", filename, err)
		}
	}

	return nil
}

// CreateTestDirectories creates test directory structure
func (tdm *TestDataManager) CreateTestDirectories(baseDir string) error {
	dirs := []string{
		"config",
		"keys",
		"logs",
		"temp",
		"nested/deep/structure",
		"permissions-test",
	}

	for _, dir := range dirs {
		dirPath := filepath.Join(baseDir, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dirPath, err)
		}
	}

	// Create a directory with restricted permissions for testing
	restrictedDir := filepath.Join(baseDir, "restricted")
	if err := os.MkdirAll(restrictedDir, 0000); err != nil {
		return fmt.Errorf("failed to create restricted directory: %w", err)
	}

	return nil
}

// CleanupTestData removes all test data
func (tdm *TestDataManager) CleanupTestData() error {
	return os.RemoveAll(tdm.BaseDir)
}

// Fixture types

// ConfigurationFixture represents a configuration test fixture
type ConfigurationFixture struct {
	Name          string
	Description   string
	Content       string
	Expected      types.SetupConfig
	ExpectedError string
}

// SetupOptionsFixture represents setup options test fixture
type SetupOptionsFixture struct {
	Name        string
	Description string
	Options     types.SetupOptions
}

// ValidationResultFixture represents validation result test fixture
type ValidationResultFixture struct {
	Name        string
	Description string
	Result      types.ValidationResult
}

// SetupResultFixture represents setup result test fixture
type SetupResultFixture struct {
	Name        string
	Description string
	Result      types.SetupResult
}

// Helper functions

func generateLargeConfig(entries int) string {
	var builder strings.Builder
	builder.WriteString("manager:\n")
	builder.WriteString("  home_dir: /home/user/.syntropy\n")
	builder.WriteString("  log_level: info\n")
	builder.WriteString("  api_endpoint: https://api.syntropy.com\n")
	builder.WriteString("  large_data:\n")

	for i := 0; i < entries; i++ {
		builder.WriteString(fmt.Sprintf("    entry_%d: value_%d\n", i, i))
	}

	return builder.String()
}

// GetMockAPIResponses returns mock API responses for testing
func (tdm *TestDataManager) GetMockAPIResponses() map[string]string {
	return map[string]string{
		"setup_success": `{
  "success": true,
  "message": "Setup completed successfully",
  "config": {
    "home_dir": "/home/user/.syntropy",
    "api_endpoint": "https://api.syntropy.com"
  }
}`,
		"setup_failure": `{
  "success": false,
  "error": "Invalid configuration provided",
  "details": "Missing required field: home_dir"
}`,
		"validation_success": `{
  "valid": true,
  "warnings": [],
  "errors": [],
  "environment": {
    "os": "linux",
    "architecture": "amd64",
    "has_admin_rights": false
  }
}`,
		"validation_failure": `{
  "valid": false,
  "warnings": ["Low disk space"],
  "errors": ["Missing dependencies"],
  "environment": {
    "os": "linux",
    "architecture": "amd64",
    "has_admin_rights": false
  }
}`,
		"status_active": `{
  "status": "active",
  "uptime": "2h30m",
  "last_check": "2023-01-01T12:00:00Z"
}`,
		"status_inactive": `{
  "status": "inactive",
  "last_error": "Connection timeout",
  "last_check": "2023-01-01T11:45:00Z"
}`,
	}
}

// GetTestEnvironments returns different test environment configurations
func (tdm *TestDataManager) GetTestEnvironments() []types.EnvironmentInfo {
	return []types.EnvironmentInfo{
		{
			OS:              "linux",
			Architecture:    "amd64",
			HasAdminRights:  false,
			AvailableDiskGB: 100.0,
			HasInternet:     true,
			HomeDir:         "/home/user",
			SystemResources: types.SystemResources{
				TotalMemoryGB: 16.0,
				CPUCores:      8,
			},
		},
		{
			OS:              "windows",
			Architecture:    "amd64",
			HasAdminRights:  true,
			AvailableDiskGB: 250.0,
			HasInternet:     true,
			HomeDir:         "C:\\Users\\User",
			SystemResources: types.SystemResources{
				TotalMemoryGB: 32.0,
				CPUCores:      16,
			},
		},
		{
			OS:              "darwin",
			Architecture:    "arm64",
			HasAdminRights:  false,
			AvailableDiskGB: 500.0,
			HasInternet:     true,
			HomeDir:         "/Users/user",
			SystemResources: types.SystemResources{
				TotalMemoryGB: 64.0,
				CPUCores:      10,
			},
		},
		{
			OS:              "linux",
			Architecture:    "arm64",
			HasAdminRights:  false,
			AvailableDiskGB: 2.0, // Low disk space
			HasInternet:     false, // No internet
			HomeDir:         "/home/user",
			SystemResources: types.SystemResources{
				TotalMemoryGB: 4.0, // Low memory
				CPUCores:      2,   // Low CPU
			},
		},
	}
}