// Package integration provides integration tests for the API
package integration

import (
	"testing"
	"time"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/api/middleware"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/api/services/config"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/api/services/validation"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/api/types"
)

// TestSetupIntegration tests the complete setup integration
func TestSetupIntegration(t *testing.T) {
	logger := middleware.NewSimpleLogger()

	// Create services
	validationService := validation.NewValidationService(logger)
	configService := config.NewConfigService(logger)

	// Create test environment
	environment := &types.EnvironmentInfo{
		OS:              "windows",
		OSVersion:       "10.0.19042",
		Architecture:    "amd64",
		HomeDir:         "C:\\Users\\Test",
		HasAdminRights:  true,
		AvailableDiskGB: 50.0,
		HasInternet:     true,
		EnvironmentVars: map[string]string{
			"SYNTROPY_HOME": "C:\\Users\\Test\\.syntropy",
		},
		Features:     []string{"windows_support"},
		Capabilities: []string{"powershell", "windows_service"},
	}

	// Test 1: Environment Validation
	t.Run("EnvironmentValidation", func(t *testing.T) {
		req := &types.ValidationRequest{
			Type:        "environment",
			Environment: environment,
			Interface:   "cli",
			UserID:      "test_user",
			SessionID:   "test_session",
			Options: types.ValidationOptions{
				Detailed: true,
				AutoFix:  false,
				Parallel: false,
			},
		}

		result, err := validationService.ValidateEnvironment(req)
		if err != nil {
			t.Fatalf("Environment validation failed: %v", err)
		}

		if !result.Valid {
			t.Errorf("Environment validation should be valid, but got errors: %v", result.Errors)
		}

		if result.Environment.OS != "windows" {
			t.Errorf("Expected OS 'windows', got '%s'", result.Environment.OS)
		}

		t.Logf("Environment validation completed successfully")
	})

	// Test 2: Security Validation
	t.Run("SecurityValidation", func(t *testing.T) {
		req := &types.ValidationRequest{
			Type:        "security",
			Environment: environment,
			Interface:   "cli",
			UserID:      "test_user",
			SessionID:   "test_session",
			Options: types.ValidationOptions{
				Detailed: true,
				AutoFix:  false,
				Parallel: false,
			},
		}

		result, err := validationService.ValidateSecurity(req)
		if err != nil {
			t.Fatalf("Security validation failed: %v", err)
		}

		if !result.Security.EncryptionAvailable {
			t.Error("Encryption should be available")
		}

		if !result.Security.SecureRandom {
			t.Error("Secure random should be available")
		}

		if !result.Security.KeyGeneration {
			t.Error("Key generation should be available")
		}

		t.Logf("Security validation completed successfully")
	})

	// Test 3: Performance Validation
	t.Run("PerformanceValidation", func(t *testing.T) {
		req := &types.ValidationRequest{
			Type:        "performance",
			Environment: environment,
			Interface:   "cli",
			UserID:      "test_user",
			SessionID:   "test_session",
			Options: types.ValidationOptions{
				Detailed: true,
				AutoFix:  false,
				Parallel: false,
			},
		}

		result, err := validationService.ValidatePerformance(req)
		if err != nil {
			t.Fatalf("Performance validation failed: %v", err)
		}

		if result.Performance.OverallScore < 0 || result.Performance.OverallScore > 100 {
			t.Errorf("Overall score should be between 0 and 100, got %.2f", result.Performance.OverallScore)
		}

		t.Logf("Performance validation completed successfully, score: %.2f", result.Performance.OverallScore)
	})

	// Test 4: Comprehensive Validation
	t.Run("ComprehensiveValidation", func(t *testing.T) {
		req := &types.ValidationRequest{
			Type:        "comprehensive",
			Environment: environment,
			Interface:   "cli",
			UserID:      "test_user",
			SessionID:   "test_session",
			Options: types.ValidationOptions{
				Detailed: true,
				AutoFix:  false,
				Parallel: true,
			},
		}

		result, err := validationService.ValidateAll(req)
		if err != nil {
			t.Fatalf("Comprehensive validation failed: %v", err)
		}

		t.Logf("Comprehensive validation completed - Valid: %v, Errors: %d, Warnings: %d",
			result.Valid, len(result.Errors), len(result.Warnings))
	})

	// Test 5: Configuration Generation
	t.Run("ConfigurationGeneration", func(t *testing.T) {
		req := &types.ConfigRequest{
			Type:        "setup",
			Environment: environment,
			Interface:   "cli",
			UserID:      "test_user",
			SessionID:   "test_session",
			Options: types.ConfigOptions{
				Force:           false,
				Backup:          true,
				Validate:        true,
				Encrypt:         false,
				Format:          "yaml",
				IncludeDefaults: true,
			},
		}

		config, err := configService.GenerateConfig(req)
		if err != nil {
			t.Fatalf("Configuration generation failed: %v", err)
		}

		if config.Manager.HomeDir == "" {
			t.Error("Manager home directory should not be empty")
		}

		if config.OwnerKey.Path == "" {
			t.Error("Owner key path should not be empty")
		}

		if config.Interface.Type != "cli" {
			t.Errorf("Expected interface type 'cli', got '%s'", config.Interface.Type)
		}

		t.Logf("Configuration generated successfully")
	})

	// Test 6: Configuration Validation
	t.Run("ConfigurationValidation", func(t *testing.T) {
		// First generate a configuration
		configReq := &types.ConfigRequest{
			Type:        "setup",
			Environment: environment,
			Interface:   "cli",
			UserID:      "test_user",
			SessionID:   "test_session",
			Options: types.ConfigOptions{
				Validate: true,
			},
		}

		config, err := configService.GenerateConfig(configReq)
		if err != nil {
			t.Fatalf("Configuration generation failed: %v", err)
		}

		// Then validate it
		validationReq := &types.ValidationRequest{
			Type:      "configuration",
			Interface: "cli",
			UserID:    "test_user",
			SessionID: "test_session",
		}

		result, err := validationService.ValidateConfig(validationReq, config)
		if err != nil {
			t.Fatalf("Configuration validation failed: %v", err)
		}

		if !result.Valid {
			t.Errorf("Generated configuration should be valid, but got errors: %v", result.Errors)
		}

		t.Logf("Configuration validation completed successfully")
	})

	// Test 7: Setup Execution
	t.Run("SetupExecution", func(t *testing.T) {
		setupReq := &types.SetupRequest{
			Options: types.SetupOptions{
				Force:          false,
				InstallService: false,
				ConfigPath:     "",
				HomeDir:        "C:\\Users\\Test\\.syntropy",
				Interface:      "cli",
			},
			Environment: environment,
			Interface:   "cli",
			UserID:      "test_user",
			SessionID:   "test_session",
			CustomData: map[string]interface{}{
				"test": true,
			},
		}

		result, err := configService.SetupService().ExecuteSetup(setupReq)
		if err != nil {
			t.Fatalf("Setup execution failed: %v", err)
		}

		if !result.Success {
			t.Errorf("Setup should succeed, but got error: %v", result.Error)
		}

		if result.Interface != "cli" {
			t.Errorf("Expected interface 'cli', got '%s'", result.Interface)
		}

		if result.Duration <= 0 {
			t.Error("Setup duration should be positive")
		}

		t.Logf("Setup execution completed successfully in %v", result.Duration)
	})

	// Test 8: Setup Status
	t.Run("SetupStatus", func(t *testing.T) {
		status, err := configService.SetupService().GetSetupStatus("cli", "test_user")
		if err != nil {
			t.Fatalf("Setup status retrieval failed: %v", err)
		}

		if status["interface"] != "cli" {
			t.Errorf("Expected interface 'cli', got '%s'", status["interface"])
		}

		if status["user_id"] != "test_user" {
			t.Errorf("Expected user_id 'test_user', got '%s'", status["user_id"])
		}

		t.Logf("Setup status retrieved successfully")
	})

	// Test 9: Configuration Backup
	t.Run("ConfigurationBackup", func(t *testing.T) {
		backupReq := &types.ConfigRequest{
			Type:        "setup",
			Environment: environment,
			Interface:   "cli",
			UserID:      "test_user",
			SessionID:   "test_session",
			Options: types.ConfigOptions{
				Backup:  true,
				Encrypt: true,
			},
		}

		backup, err := configService.CreateBackup(backupReq)
		if err != nil {
			t.Fatalf("Backup creation failed: %v", err)
		}

		if backup.ID == "" {
			t.Error("Backup ID should not be empty")
		}

		if backup.Config == nil {
			t.Error("Backup configuration should not be nil")
		}

		if backup.Timestamp.IsZero() {
			t.Error("Backup timestamp should not be zero")
		}

		t.Logf("Configuration backup created successfully: %s", backup.ID)
	})

	// Test 10: Template Generation
	t.Run("TemplateGeneration", func(t *testing.T) {
		template, err := configService.GetTemplate("cli", "windows", "default")
		if err != nil {
			t.Fatalf("Template generation failed: %v", err)
		}

		if template.Name == "" {
			t.Error("Template name should not be empty")
		}

		if template.Interface != "cli" {
			t.Errorf("Expected interface 'cli', got '%s'", template.Interface)
		}

		if template.Environment != "windows" {
			t.Errorf("Expected environment 'windows', got '%s'", template.Environment)
		}

		if template.Content == "" {
			t.Error("Template content should not be empty")
		}

		t.Logf("Template generated successfully")
	})
}

// TestSetupIntegrationParallel tests setup integration with parallel execution
func TestSetupIntegrationParallel(t *testing.T) {
	logger := middleware.NewSimpleLogger()
	validationService := validation.NewValidationService(logger)

	environment := &types.EnvironmentInfo{
		OS:              "linux",
		OSVersion:       "Ubuntu 20.04",
		Architecture:    "amd64",
		HomeDir:         "/home/test",
		HasAdminRights:  true,
		AvailableDiskGB: 100.0,
		HasInternet:     true,
	}

	t.Run("ParallelValidation", func(t *testing.T) {
		req := &types.ValidationRequest{
			Type:        "comprehensive",
			Environment: environment,
			Interface:   "cli",
			UserID:      "test_user",
			SessionID:   "test_session",
			Options: types.ValidationOptions{
				Detailed: true,
				AutoFix:  false,
				Parallel: true,
			},
		}

		startTime := time.Now()
		result, err := validationService.ValidateAll(req)
		duration := time.Since(startTime)

		if err != nil {
			t.Fatalf("Parallel validation failed: %v", err)
		}

		t.Logf("Parallel validation completed in %v - Valid: %v", duration, result.Valid)

		// Parallel execution should be faster than sequential
		// (This is a basic test - in production you'd have more sophisticated timing tests)
		if duration > 5*time.Second {
			t.Logf("Warning: Parallel validation took longer than expected: %v", duration)
		}
	})
}

// TestSetupIntegrationErrorHandling tests error handling in setup integration
func TestSetupIntegrationErrorHandling(t *testing.T) {
	logger := middleware.NewSimpleLogger()
	validationService := validation.NewValidationService(logger)

	t.Run("InvalidEnvironment", func(t *testing.T) {
		// Test with invalid environment
		invalidEnv := &types.EnvironmentInfo{
			OS:              "unsupported",
			OSVersion:       "",
			Architecture:    "",
			HomeDir:         "",
			HasAdminRights:  false,
			AvailableDiskGB: 0,
			HasInternet:     false,
		}

		req := &types.ValidationRequest{
			Type:        "environment",
			Environment: invalidEnv,
			Interface:   "cli",
			UserID:      "test_user",
			SessionID:   "test_session",
		}

		result, err := validationService.ValidateEnvironment(req)
		if err != nil {
			t.Fatalf("Environment validation should not fail with error, got: %v", err)
		}

		// Should have validation errors
		if result.Valid {
			t.Error("Environment validation should not be valid for unsupported OS")
		}

		if len(result.Errors) == 0 {
			t.Error("Should have validation errors for invalid environment")
		}

		t.Logf("Invalid environment handled correctly with %d errors", len(result.Errors))
	})

	t.Run("NilConfiguration", func(t *testing.T) {
		req := &types.ValidationRequest{
			Type:      "configuration",
			Interface: "cli",
			UserID:    "test_user",
			SessionID: "test_session",
		}

		result, err := validationService.ValidateConfig(req, nil)
		if err != nil {
			t.Fatalf("Configuration validation should not fail with error, got: %v", err)
		}

		// Should have validation errors
		if result.Valid {
			t.Error("Configuration validation should not be valid for nil config")
		}

		if len(result.Errors) == 0 {
			t.Error("Should have validation errors for nil configuration")
		}

		t.Logf("Nil configuration handled correctly with %d errors", len(result.Errors))
	})
}

// BenchmarkSetupIntegration benchmarks the setup integration performance
func BenchmarkSetupIntegration(b *testing.B) {
	logger := middleware.NewSimpleLogger()
	validationService := validation.NewValidationService(logger)
	configService := config.NewConfigService(logger)

	environment := &types.EnvironmentInfo{
		OS:              "windows",
		OSVersion:       "10.0.19042",
		Architecture:    "amd64",
		HomeDir:         "C:\\Users\\Test",
		HasAdminRights:  true,
		AvailableDiskGB: 50.0,
		HasInternet:     true,
	}

	b.Run("EnvironmentValidation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			req := &types.ValidationRequest{
				Type:        "environment",
				Environment: environment,
				Interface:   "cli",
				UserID:      "test_user",
				SessionID:   "test_session",
			}

			_, err := validationService.ValidateEnvironment(req)
			if err != nil {
				b.Fatalf("Environment validation failed: %v", err)
			}
		}
	})

	b.Run("ConfigurationGeneration", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			req := &types.ConfigRequest{
				Type:        "setup",
				Environment: environment,
				Interface:   "cli",
				UserID:      "test_user",
				SessionID:   "test_session",
			}

			_, err := configService.GenerateConfig(req)
			if err != nil {
				b.Fatalf("Configuration generation failed: %v", err)
			}
		}
	})
}
