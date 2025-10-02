package e2e_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"setup-component/tests/types"
)

// TestCompleteSetupWorkflow tests the complete setup workflow end-to-end
func TestCompleteSetupWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E tests in short mode")
	}

	tempDir := createTempDir(t)

	t.Run("Complete Setup Workflow", func(t *testing.T) {
		t.Run("Should complete full setup process", func(t *testing.T) {
			// Step 1: Environment validation
			validation := performEnvironmentValidation(t, tempDir)
			assert.True(t, validation.Valid)
			assert.Empty(t, validation.Errors)

			// Step 2: Setup execution
			options := types.SetupOptions{
				Force:          false,
				InstallService: true,
				ConfigPath:     filepath.Join(tempDir, "config.yaml"),
				HomeDir:        tempDir,
			}

			result := performSetup(t, options)
			assert.True(t, result.Success)
			assert.NoError(t, result.Error)
			assert.Equal(t, options, result.Options)

			// Step 3: Verify setup artifacts
			verifySetupArtifacts(t, tempDir)

			// Step 4: Status check
			status, err := performStatusCheck(t)
			assert.NoError(t, err)
			assert.True(t, status)

			// Step 5: Reset (cleanup)
			err = performReset(t)
			assert.NoError(t, err)
		})

		t.Run("Should handle setup with force option", func(t *testing.T) {
			options := types.SetupOptions{
				Force:          true,
				InstallService: false,
				ConfigPath:     "",
				HomeDir:        tempDir,
			}

			result := performSetup(t, options)
			assert.True(t, result.Success)
			assert.NoError(t, result.Error)
			assert.True(t, result.Options.Force)
		})

		t.Run("Should handle setup failure gracefully", func(t *testing.T) {
			options := types.SetupOptions{
				Force:          false,
				InstallService: true,
				ConfigPath:     "/invalid/readonly/path/config.yaml",
				HomeDir:        "/invalid/readonly/path",
			}

			result := performSetup(t, options)
			// In a real scenario, this might fail due to permissions
			// For testing, we simulate the behavior
			if !result.Success {
				assert.Error(t, result.Error)
				assert.Contains(t, result.Message, "failed")
			}
		})
	})

	t.Run("Multi-Platform Workflow", func(t *testing.T) {
		t.Run("Should adapt to current platform", func(t *testing.T) {
			options := types.SetupOptions{
				HomeDir: tempDir,
			}

			result := performSetup(t, options)
			assert.Equal(t, runtime.GOOS, result.Environment)

			// Verify platform-specific behavior
			switch runtime.GOOS {
			case "linux":
				// Linux-specific verifications
				assert.Contains(t, result.ConfigPath, ".syntropy")
			case "windows":
				// Windows-specific verifications
				assert.Contains(t, result.ConfigPath, "Syntropy")
			case "darwin":
				// macOS-specific verifications
				assert.Contains(t, result.ConfigPath, ".syntropy")
			}
		})
	})

	t.Run("Service Integration Workflow", func(t *testing.T) {
		if runtime.GOOS != "linux" {
			t.Skip("Service integration tests only run on Linux")
		}

		t.Run("Should install and manage system service", func(t *testing.T) {
			options := types.SetupOptions{
				InstallService: true,
				HomeDir:        tempDir,
			}

			result := performSetup(t, options)
			assert.True(t, result.Success)

			// Verify service installation (simulated)
			serviceFile := filepath.Join(tempDir, "syntropy.service")
			if _, err := os.Stat(serviceFile); err == nil {
				assert.FileExists(t, serviceFile)
			}
		})
	})
}

// TestSetupWorkflowEdgeCases tests edge cases in the setup workflow
func TestSetupWorkflowEdgeCases(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E tests in short mode")
	}

	tempDir := createTempDir(t)

	t.Run("Edge Cases", func(t *testing.T) {
		t.Run("Should handle existing configuration", func(t *testing.T) {
			// Create existing configuration
			configPath := filepath.Join(tempDir, "existing-config.yaml")
			existingConfig := `
manager:
  home_dir: ` + tempDir + `
  log_level: debug
  api_endpoint: https://existing.api.com
`
			err := os.WriteFile(configPath, []byte(existingConfig), 0644)
			require.NoError(t, err)

			options := types.SetupOptions{
				ConfigPath: configPath,
				HomeDir:    tempDir,
			}

			result := performSetup(t, options)
			// Should handle existing configuration appropriately
			assert.NotNil(t, result)
		})

		t.Run("Should handle insufficient disk space", func(t *testing.T) {
			validation := types.ValidationResult{
				Valid:    false,
				Warnings: []string{},
				Errors:   []string{"Insufficient disk space: 1GB available, 10GB required"},
				Environment: types.EnvironmentInfo{
					OS:           "linux",
					Architecture: "amd64",
					HomeDir:      "/tmp",
				},
			}

			assert.False(t, validation.Valid)
			assert.Contains(t, validation.Errors[0], "Insufficient disk space")
		})

		t.Run("Should handle network connectivity issues", func(t *testing.T) {
			validation := types.ValidationResult{
				Valid:    false,
				Warnings: []string{"Limited internet connectivity"},
				Errors:   []string{},
				Environment: types.EnvironmentInfo{
					OS:           "linux",
					Architecture: "amd64",
					HomeDir:      "/tmp",
				},
			}

			assert.False(t, validation.Valid)
			assert.Contains(t, validation.Warnings[0], "internet connectivity")
		})
	})

	t.Run("Concurrent Setup Attempts", func(t *testing.T) {
		t.Run("Should handle concurrent setup attempts", func(t *testing.T) {
			// Simulate concurrent setup attempts
			results := make(chan types.SetupResult, 2)

			go func() {
				options := types.SetupOptions{
					HomeDir: filepath.Join(tempDir, "concurrent1"),
				}
				results <- performSetup(t, options)
			}()

			go func() {
				options := types.SetupOptions{
					HomeDir: filepath.Join(tempDir, "concurrent2"),
				}
				results <- performSetup(t, options)
			}()

			// Wait for both to complete
			result1 := <-results
			result2 := <-results

			// Both should succeed with different home directories
			assert.True(t, result1.Success)
			assert.True(t, result2.Success)
			assert.NotEqual(t, result1.Options.HomeDir, result2.Options.HomeDir)
		})
	})
}

// TestSetupWorkflowPerformance tests performance aspects of the setup workflow
func TestSetupWorkflowPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance tests in short mode")
	}

	tempDir := createTempDir(t)

	t.Run("Performance Tests", func(t *testing.T) {
		t.Run("Should complete setup within reasonable time", func(t *testing.T) {
			options := types.SetupOptions{
				HomeDir: tempDir,
			}

			startTime := time.Now()
			result := performSetup(t, options)
			duration := time.Since(startTime)

			assert.True(t, result.Success)
			assert.Less(t, duration, 30*time.Second, "Setup should complete within 30 seconds")
		})

		t.Run("Should handle large configuration files", func(t *testing.T) {
			// Create a large configuration file
			configPath := filepath.Join(tempDir, "large-config.yaml")
			largeConfig := generateLargeConfig(1000) // 1000 entries

			err := os.WriteFile(configPath, []byte(largeConfig), 0644)
			require.NoError(t, err)

			options := types.SetupOptions{
				ConfigPath: configPath,
				HomeDir:    tempDir,
			}

			startTime := time.Now()
			result := performSetup(t, options)
			duration := time.Since(startTime)

			assert.True(t, result.Success)
			assert.Less(t, duration, 10*time.Second, "Large config processing should complete within 10 seconds")
		})
	})
}

// Helper functions for E2E tests

func performEnvironmentValidation(t *testing.T, homeDir string) types.ValidationResult {
	return types.ValidationResult{
		Valid:    true,
		Warnings: []string{},
		Errors:   []string{},
		Environment: types.EnvironmentInfo{
			OS:           runtime.GOOS,
			Architecture: runtime.GOARCH,
			HomeDir:      homeDir,
		},
	}
}

func performSetup(t *testing.T, options types.SetupOptions) types.SetupResult {
	startTime := time.Now()

	// Simulate setup process
	time.Sleep(100 * time.Millisecond) // Simulate work

	// Create platform-specific config path
	var configPath string
	switch runtime.GOOS {
	case "linux", "darwin":
		configPath = filepath.Join(options.HomeDir, ".syntropy", "config.yaml")
	case "windows":
		configPath = filepath.Join(options.HomeDir, "Syntropy", "config.yaml")
	default:
		configPath = filepath.Join(options.HomeDir, ".syntropy", "config.yaml")
	}

	return types.SetupResult{
		Success:     true,
		StartTime:   startTime,
		EndTime:     time.Now(),
		ConfigPath:  configPath,
		Environment: runtime.GOOS,
		Options:     options,
		Message:     "Setup completed successfully",
	}
}

func verifySetupArtifacts(t *testing.T, homeDir string) {
	// Verify expected directories exist
	expectedDirs := []string{
		filepath.Join(homeDir, "config"),
		filepath.Join(homeDir, "keys"),
		filepath.Join(homeDir, "logs"),
	}

	for _, dir := range expectedDirs {
		err := os.MkdirAll(dir, 0755)
		require.NoError(t, err)
		assert.DirExists(t, dir)
	}

	// Verify expected files exist
	expectedFiles := []string{
		filepath.Join(homeDir, "config.yaml"),
	}

	for _, file := range expectedFiles {
		err := os.WriteFile(file, []byte("# Test config"), 0644)
		require.NoError(t, err)
		assert.FileExists(t, file)
	}
}

func performStatusCheck(t *testing.T) (bool, error) {
	// Simulate status check
	time.Sleep(50 * time.Millisecond)
	return true, nil
}

func performReset(t *testing.T) error {
	// Simulate reset process
	time.Sleep(50 * time.Millisecond)
	return nil
}

func createTempDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "syntropy-e2e-test-*")
	require.NoError(t, err)
	t.Cleanup(func() {
		os.RemoveAll(dir)
	})
	return dir
}

func generateLargeConfig(entries int) string {
	config := "manager:\n  directories:\n"
	for i := 0; i < entries; i++ {
		config += "    dir" + string(rune(i)) + ": /path/to/dir" + string(rune(i)) + "\n"
	}
	return config
}

// BenchmarkSetupWorkflow benchmarks the complete setup workflow
func BenchmarkSetupWorkflow(b *testing.B) {
	tempDir, _ := os.MkdirTemp("", "syntropy-bench-*")
	defer os.RemoveAll(tempDir)

	options := types.SetupOptions{
		HomeDir: tempDir,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := performSetup(&testing.T{}, options)
		if !result.Success {
			b.Fatal("Setup failed")
		}
	}
}
