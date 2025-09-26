// Package setup provides integration with the API central
package setup

import (
	"fmt"
	"time"

	confighandlers "github.com/syntropy-cc/syntropy-cooperative-grid/manager/api/handlers/config"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/api/middleware"
	configsvc "github.com/syntropy-cc/syntropy-cooperative-grid/manager/api/services/config"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/api/services/validation"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/api/types"
)

// APIIntegration provides integration with the API central
type APIIntegration struct {
	configHandler     *confighandlers.ConfigHandler
	setupHandler      *confighandlers.SetupHandler
	validationHandler *confighandlers.ValidationHandler
	validationService *validation.ValidationService
	configService     *configsvc.ConfigService
	logger            middleware.Logger
}

// NewAPIIntegration creates a new API integration
func NewAPIIntegration() *APIIntegration {
	logger := middleware.NewSimpleLogger()

	// services
	validationService := validation.NewValidationService(logger)
	configService := configsvc.NewConfigService(logger)

	// handlers
	configHandler := confighandlers.NewConfigHandler(configService, validationService, logger)
	setupHandler := confighandlers.NewSetupHandler(configService, validationService, configService.SetupService(), logger)
	validationHandler := confighandlers.NewValidationHandler(validationService, logger)

	return &APIIntegration{
		configHandler:     configHandler,
		setupHandler:      setupHandler,
		validationHandler: validationHandler,
		validationService: validationService,
		configService:     configService,
		logger:            logger,
	}
}

// SetupWithAPI performs setup using the API central
func (ai *APIIntegration) SetupWithAPI(options *types.SetupOptions, environment *types.EnvironmentInfo, interfaceType string) (*types.SetupResult, error) {
	ai.logger.Info("Starting setup with API integration", map[string]interface{}{
		"interface": interfaceType,
		"force":     options.Force,
	})

	// Create setup request
	req := &types.SetupRequest{
		Options:     *options,
		Environment: environment,
		Interface:   interfaceType,
		UserID:      "cli_user", // In production, get from authentication
		SessionID:   generateSessionID(),
		CustomData: map[string]interface{}{
			"source":    "cli_setup",
			"timestamp": time.Now(),
		},
	}

	// Perform setup using API central
	result, err := ai.configService.SetupService().ExecuteSetup(req)
	if err != nil {
		ai.logger.Error("Setup execution failed", map[string]interface{}{
			"error":     err.Error(),
			"interface": interfaceType,
		})
		return nil, fmt.Errorf("setup execution failed: %w", err)
	}

	ai.logger.Info("Setup completed with API integration", map[string]interface{}{
		"interface": interfaceType,
		"success":   result.Success,
		"duration":  result.Duration.String(),
	})

	return result, nil
}

// ValidateWithAPI performs validation using the API central
func (ai *APIIntegration) ValidateWithAPI(environment *types.EnvironmentInfo, interfaceType string) (*types.ValidationResult, error) {
	ai.logger.Info("Starting validation with API integration", map[string]interface{}{
		"interface": interfaceType,
	})

	// Create validation request
	req := &types.ValidationRequest{
		Type:        "comprehensive",
		Environment: environment,
		Interface:   interfaceType,
		UserID:      "cli_user",
		SessionID:   generateSessionID(),
		Options: types.ValidationOptions{
			Detailed: true,
			AutoFix:  false,
			Parallel: true,
		},
		CustomData: map[string]interface{}{
			"source":    "cli_validation",
			"timestamp": time.Now(),
		},
	}

	// Perform validation using API central
	result, err := ai.validationService.ValidateAll(req)
	if err != nil {
		ai.logger.Error("Validation failed", map[string]interface{}{
			"error":     err.Error(),
			"interface": interfaceType,
		})
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	ai.logger.Info("Validation completed with API integration", map[string]interface{}{
		"interface": interfaceType,
		"valid":     result.Valid,
		"errors":    len(result.Errors),
		"warnings":  len(result.Warnings),
		"duration":  result.Duration.String(),
	})

	return result, nil
}

// GenerateConfigWithAPI generates configuration using the API central
func (ai *APIIntegration) GenerateConfigWithAPI(environment *types.EnvironmentInfo, interfaceType string, configType string) (*types.SetupConfig, error) {
	ai.logger.Info("Generating configuration with API integration", map[string]interface{}{
		"interface":   interfaceType,
		"config_type": configType,
	})

	// Create configuration request
	req := &types.ConfigRequest{
		Type:        configType,
		Environment: environment,
		Interface:   interfaceType,
		UserID:      "cli_user",
		SessionID:   generateSessionID(),
		Options: types.ConfigOptions{
			Force:           false,
			Backup:          true,
			Validate:        true,
			Encrypt:         false,
			Format:          "yaml",
			IncludeDefaults: true,
		},
		CustomData: map[string]interface{}{
			"source":    "cli_config_generation",
			"timestamp": time.Now(),
		},
	}

	// Generate configuration using API central
	config, err := ai.configService.GenerateConfig(req)
	if err != nil {
		ai.logger.Error("Configuration generation failed", map[string]interface{}{
			"error":     err.Error(),
			"interface": interfaceType,
		})
		return nil, fmt.Errorf("configuration generation failed: %w", err)
	}

	ai.logger.Info("Configuration generated with API integration", map[string]interface{}{
		"interface":   interfaceType,
		"config_type": configType,
	})

	return config, nil
}

// ValidateConfigWithAPI validates configuration using the API central
func (ai *APIIntegration) ValidateConfigWithAPI(config *types.SetupConfig, interfaceType string) (*types.ValidationResult, error) {
	ai.logger.Info("Validating configuration with API integration", map[string]interface{}{
		"interface": interfaceType,
	})

	// Create validation request
	req := &types.ValidationRequest{
		Type:      "configuration",
		Interface: interfaceType,
		UserID:    "cli_user",
		SessionID: generateSessionID(),
		CustomData: map[string]interface{}{
			"source":    "cli_config_validation",
			"timestamp": time.Now(),
		},
	}

	// Validate configuration using API central
	result, err := ai.validationService.ValidateConfig(req, config)
	if err != nil {
		ai.logger.Error("Configuration validation failed", map[string]interface{}{
			"error":     err.Error(),
			"interface": interfaceType,
		})
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	ai.logger.Info("Configuration validation completed with API integration", map[string]interface{}{
		"interface": interfaceType,
		"valid":     result.Valid,
		"errors":    len(result.Errors),
		"warnings":  len(result.Warnings),
	})

	return result, nil
}

// GetSetupStatusWithAPI gets setup status using the API central
func (ai *APIIntegration) GetSetupStatusWithAPI(interfaceType string) (map[string]interface{}, error) {
	ai.logger.Info("Getting setup status with API integration", map[string]interface{}{
		"interface": interfaceType,
	})

	// Get setup status using API central
	status, err := ai.configService.SetupService().GetSetupStatus(interfaceType, "cli_user")
	if err != nil {
		ai.logger.Error("Failed to get setup status", map[string]interface{}{
			"error":     err.Error(),
			"interface": interfaceType,
		})
		return nil, fmt.Errorf("failed to get setup status: %w", err)
	}

	return status, nil
}

// ResetSetupWithAPI resets setup using the API central
func (ai *APIIntegration) ResetSetupWithAPI(interfaceType string) error {
	ai.logger.Info("Resetting setup with API integration", map[string]interface{}{
		"interface": interfaceType,
	})

	// Create reset request
	req := &types.SetupRequest{
		Interface: interfaceType,
		UserID:    "cli_user",
		SessionID: generateSessionID(),
		CustomData: map[string]interface{}{
			"source":    "cli_reset",
			"timestamp": time.Now(),
		},
	}

	// Reset setup using API central
	err := ai.configService.SetupService().ResetSetup(req)
	if err != nil {
		ai.logger.Error("Setup reset failed", map[string]interface{}{
			"error":     err.Error(),
			"interface": interfaceType,
		})
		return fmt.Errorf("setup reset failed: %w", err)
	}

	ai.logger.Info("Setup reset completed with API integration", map[string]interface{}{
		"interface": interfaceType,
	})

	return nil
}

// CreateBackupWithAPI creates backup using the API central
func (ai *APIIntegration) CreateBackupWithAPI(environment *types.EnvironmentInfo, interfaceType string) (*types.ConfigBackup, error) {
	ai.logger.Info("Creating backup with API integration", map[string]interface{}{
		"interface": interfaceType,
	})

	// Create backup request
	req := &types.ConfigRequest{
		Type:        "setup",
		Environment: environment,
		Interface:   interfaceType,
		UserID:      "cli_user",
		SessionID:   generateSessionID(),
		Options: types.ConfigOptions{
			Backup:   true,
			Encrypt:  true,
			Validate: true,
		},
		CustomData: map[string]interface{}{
			"source":    "cli_backup",
			"timestamp": time.Now(),
		},
	}

	// Create backup using API central
	backup, err := ai.configService.CreateBackup(req)
	if err != nil {
		ai.logger.Error("Backup creation failed", map[string]interface{}{
			"error":     err.Error(),
			"interface": interfaceType,
		})
		return nil, fmt.Errorf("backup creation failed: %w", err)
	}

	ai.logger.Info("Backup created with API integration", map[string]interface{}{
		"interface": interfaceType,
		"backup_id": backup.ID,
	})

	return backup, nil
}

// GetTemplateWithAPI gets configuration template using the API central
func (ai *APIIntegration) GetTemplateWithAPI(interfaceType, environment, templateName string) (*types.ConfigTemplate, error) {
	ai.logger.Info("Getting template with API integration", map[string]interface{}{
		"interface":   interfaceType,
		"environment": environment,
		"template":    templateName,
	})

	// Get template using API central
	template, err := ai.configService.GetTemplate(interfaceType, environment, templateName)
	if err != nil {
		ai.logger.Error("Failed to get template", map[string]interface{}{
			"error":     err.Error(),
			"interface": interfaceType,
		})
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	return template, nil
}

// Helper function to generate session ID
func generateSessionID() string {
	return fmt.Sprintf("cli_session_%d", time.Now().Unix())
}
