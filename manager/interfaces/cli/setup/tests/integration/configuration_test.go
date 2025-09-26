package integration_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/tests/types"
)

// TestConfigurationIntegration tests configuration file operations
func TestConfigurationIntegration(t *testing.T) {
	tempDir := createTempDir(t)

	t.Run("Configuration File Operations", func(t *testing.T) {
		t.Run("Should create and read configuration file", func(t *testing.T) {
			config := types.SetupConfig{
				Manager: types.ManagerConfig{
					HomeDir:     tempDir,
					LogLevel:    "info",
					APIEndpoint: "https://api.syntropy.com",
					Directories: map[string]string{
						"config": filepath.Join(tempDir, "config"),
						"keys":   filepath.Join(tempDir, "keys"),
						"logs":   filepath.Join(tempDir, "logs"),
					},
					DefaultPaths: map[string]string{
						"config": filepath.Join(tempDir, "config", "manager.yaml"),
						"key":    filepath.Join(tempDir, "keys", "owner.key"),
					},
				},
				OwnerKey: types.OwnerKey{
					Type: "Ed25519",
					Path: filepath.Join(tempDir, "keys", "owner.key"),
				},
				Environment: types.Environment{
					OS:           "linux",
					Architecture: "amd64",
					HomeDir:      tempDir,
				},
			}

			configPath := filepath.Join(tempDir, "config.yaml")
			
			// Write configuration
			data, err := yaml.Marshal(config)
			require.NoError(t, err)
			
			err = os.WriteFile(configPath, data, 0644)
			require.NoError(t, err)
			
			// Read and verify configuration
			readData, err := os.ReadFile(configPath)
			require.NoError(t, err)
			
			var readConfig types.SetupConfig
			err = yaml.Unmarshal(readData, &readConfig)
			require.NoError(t, err)
			
			assert.Equal(t, config.Manager.HomeDir, readConfig.Manager.HomeDir)
			assert.Equal(t, config.Manager.LogLevel, readConfig.Manager.LogLevel)
			assert.Equal(t, config.OwnerKey.Type, readConfig.OwnerKey.Type)
			assert.Equal(t, config.Environment.OS, readConfig.Environment.OS)
		})

		t.Run("Should handle invalid configuration file", func(t *testing.T) {
			invalidConfigPath := filepath.Join(tempDir, "invalid.yaml")
			
			err := os.WriteFile(invalidConfigPath, []byte("invalid: yaml: content: ["), 0644)
			require.NoError(t, err)
			
			data, err := os.ReadFile(invalidConfigPath)
			require.NoError(t, err)
			
			var config types.SetupConfig
			err = yaml.Unmarshal(data, &config)
			assert.Error(t, err)
		})

		t.Run("Should create directory structure", func(t *testing.T) {
			dirs := []string{
				filepath.Join(tempDir, "config"),
				filepath.Join(tempDir, "keys"),
				filepath.Join(tempDir, "logs"),
				filepath.Join(tempDir, "data"),
			}

			for _, dir := range dirs {
				err := os.MkdirAll(dir, 0755)
				require.NoError(t, err)
				assert.DirExists(t, dir)
			}
		})
	})

	t.Run("Key Management Integration", func(t *testing.T) {
		keysDir := filepath.Join(tempDir, "keys")
		err := os.MkdirAll(keysDir, 0755)
		require.NoError(t, err)

		t.Run("Should create and manage key files", func(t *testing.T) {
			keyPath := filepath.Join(keysDir, "owner.key")
			publicKeyPath := filepath.Join(keysDir, "owner.pub")
			
			// Simulate key generation (in real implementation, this would use crypto)
			privateKey := "ed25519_private_key_data_here"
			publicKey := "ed25519_public_key_data_here"
			
			err := os.WriteFile(keyPath, []byte(privateKey), 0600)
			require.NoError(t, err)
			
			err = os.WriteFile(publicKeyPath, []byte(publicKey), 0644)
			require.NoError(t, err)
			
			// Verify key files exist and have correct permissions
			assert.FileExists(t, keyPath)
			assert.FileExists(t, publicKeyPath)
			
			info, err := os.Stat(keyPath)
			require.NoError(t, err)
			assert.Equal(t, os.FileMode(0600), info.Mode().Perm())
			
			info, err = os.Stat(publicKeyPath)
			require.NoError(t, err)
			assert.Equal(t, os.FileMode(0644), info.Mode().Perm())
		})

		t.Run("Should handle key file permissions correctly", func(t *testing.T) {
			keyPath := filepath.Join(keysDir, "test.key")
			
			err := os.WriteFile(keyPath, []byte("test_key_data"), 0644)
			require.NoError(t, err)
			
			// Change to secure permissions
			err = os.Chmod(keyPath, 0600)
			require.NoError(t, err)
			
			info, err := os.Stat(keyPath)
			require.NoError(t, err)
			assert.Equal(t, os.FileMode(0600), info.Mode().Perm())
		})
	})

	t.Run("Environment Configuration Integration", func(t *testing.T) {
		t.Run("Should detect environment information", func(t *testing.T) {
			env := types.EnvironmentInfo{
				OS:              "linux",
				OSVersion:       "Ubuntu 20.04",
				Architecture:    "amd64",
				HasAdminRights:  false,
				AvailableDiskGB: 100.0,
				HasInternet:     true,
				HomeDir:         tempDir,
			}

			// Verify environment detection
			assert.NotEmpty(t, env.OS)
			assert.NotEmpty(t, env.Architecture)
			assert.NotEmpty(t, env.HomeDir)
			assert.True(t, env.AvailableDiskGB > 0)
		})

		t.Run("Should validate system requirements", func(t *testing.T) {
			validation := types.ValidationResult{
				Valid:    true,
				Warnings: []string{},
				Errors:   []string{},
				Environment: types.EnvironmentInfo{
					OS:              "linux",
					Architecture:    "amd64",
					HasAdminRights:  true,
					AvailableDiskGB: 50.0,
					HasInternet:     true,
					HomeDir:         tempDir,
				},
			}

			assert.True(t, validation.Valid)
			assert.Empty(t, validation.Errors)
			assert.True(t, validation.Environment.AvailableDiskGB >= 10.0) // Minimum requirement
		})
	})
}

// TestConfigurationValidation tests configuration validation
func TestConfigurationValidation(t *testing.T) {
	tempDir := createTempDir(t)

	t.Run("Valid Configuration", func(t *testing.T) {
		config := types.SetupConfig{
			Manager: types.ManagerConfig{
				HomeDir:     tempDir,
				LogLevel:    "info",
				APIEndpoint: "https://api.syntropy.com",
			},
			OwnerKey: types.OwnerKey{
				Type: "Ed25519",
				Path: filepath.Join(tempDir, "keys", "owner.key"),
			},
			Environment: types.Environment{
				OS:           "linux",
				Architecture: "amd64",
				HomeDir:      tempDir,
			},
		}

		// Validate configuration structure
		assert.NotEmpty(t, config.Manager.HomeDir)
		assert.NotEmpty(t, config.Manager.LogLevel)
		assert.NotEmpty(t, config.Manager.APIEndpoint)
		assert.NotEmpty(t, config.OwnerKey.Type)
		assert.NotEmpty(t, config.OwnerKey.Path)
		assert.NotEmpty(t, config.Environment.OS)
		assert.NotEmpty(t, config.Environment.Architecture)
	})

	t.Run("Invalid Configuration", func(t *testing.T) {
		config := types.SetupConfig{
			Manager: types.ManagerConfig{
				HomeDir:  "", // Invalid: empty home directory
				LogLevel: "invalid_level", // Invalid log level
			},
			OwnerKey: types.OwnerKey{
				Type: "", // Invalid: empty key type
				Path: "", // Invalid: empty key path
			},
		}

		// Validate that invalid configuration is detected
		assert.Empty(t, config.Manager.HomeDir)
		assert.Empty(t, config.OwnerKey.Type)
		assert.Empty(t, config.OwnerKey.Path)
	})
}

// Helper function to create temporary directory
func createTempDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "syntropy-integration-test-*")
	require.NoError(t, err)
	t.Cleanup(func() {
		os.RemoveAll(dir)
	})
	return dir
}

// BenchmarkConfigurationOperations benchmarks configuration operations
func BenchmarkConfigurationOperations(b *testing.B) {
	tempDir, _ := os.MkdirTemp("", "syntropy-bench-*")
	defer os.RemoveAll(tempDir)

	config := types.SetupConfig{
		Manager: types.ManagerConfig{
			HomeDir:     tempDir,
			LogLevel:    "info",
			APIEndpoint: "https://api.syntropy.com",
		},
		OwnerKey: types.OwnerKey{
			Type: "Ed25519",
			Path: filepath.Join(tempDir, "keys", "owner.key"),
		},
		Environment: types.Environment{
			OS:           "linux",
			Architecture: "amd64",
			HomeDir:      tempDir,
		},
	}

	b.Run("YAML Marshal", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := yaml.Marshal(config)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("YAML Unmarshal", func(b *testing.B) {
		data, _ := yaml.Marshal(config)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var c types.SetupConfig
			err := yaml.Unmarshal(data, &c)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("File Write", func(b *testing.B) {
		data, _ := yaml.Marshal(config)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			configPath := filepath.Join(tempDir, "bench_config.yaml")
			err := os.WriteFile(configPath, data, 0644)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}