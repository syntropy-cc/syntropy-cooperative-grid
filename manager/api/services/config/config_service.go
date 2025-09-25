// Package config provides configuration services for the API
package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"manager/api/middleware"
	"manager/api/types"
)

// ConfigService provides configuration services for all interfaces
type ConfigService struct {
	setupService *SetupService
	logger       middleware.Logger
}

// NewConfigService creates a new configuration service
func NewConfigService(logger middleware.Logger) *ConfigService {
	return &ConfigService{
		setupService: NewSetupService(logger),
		logger:       logger,
	}
}

// GenerateConfig generates a complete configuration
func (cs *ConfigService) GenerateConfig(req *types.ConfigRequest) (*types.SetupConfig, error) {
	cs.logger.Info("Generating configuration", map[string]interface{}{
		"interface":   req.Interface,
		"type":        req.Type,
		"environment": req.Environment.OS,
	})

	// Create base configuration
	config := &types.SetupConfig{
		Manager:     cs.generateManagerConfig(req),
		OwnerKey:    cs.generateOwnerKey(req),
		Environment: cs.generateEnvironmentConfig(req),
		Interface:   cs.generateInterfaceConfig(req),
		Security:    cs.generateSecurityConfig(req),
		Network:     cs.generateNetworkConfig(req),
		Metadata:    cs.generateConfigMetadata(req),
	}

	// Validate generated configuration
	if err := cs.validateGeneratedConfig(config, req); err != nil {
		cs.logger.Error("Generated configuration validation failed", map[string]interface{}{
			"error":     err.Error(),
			"interface": req.Interface,
		})
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	cs.logger.Info("Configuration generated successfully", map[string]interface{}{
		"interface":   req.Interface,
		"config_type": req.Type,
	})

	return config, nil
}

// CreateBackup creates a backup of the current configuration
func (cs *ConfigService) CreateBackup(req *types.ConfigRequest) (*types.ConfigBackup, error) {
	cs.logger.Info("Creating configuration backup", map[string]interface{}{
		"interface": req.Interface,
		"user_id":   req.UserID,
	})

	// Generate current configuration
	config, err := cs.GenerateConfig(req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate configuration for backup: %w", err)
	}

	// Create backup
	backup := &types.ConfigBackup{
		ID:          fmt.Sprintf("backup_%s_%d", req.Interface, time.Now().Unix()),
		Name:        fmt.Sprintf("Configuration backup for %s", req.Interface),
		Description: fmt.Sprintf("Automatic backup created on %s", time.Now().Format(time.RFC3339)),
		Config:      config,
		Timestamp:   time.Now(),
		Size:        cs.calculateConfigSize(config),
		Checksum:    cs.calculateConfigChecksum(config),
		Encrypted:   req.Options.Encrypt,
		Compressed:  true,
		Metadata: map[string]interface{}{
			"interface":  req.Interface,
			"user_id":    req.UserID,
			"session_id": req.SessionID,
			"created_by": "config_service",
		},
	}

	cs.logger.Info("Configuration backup created", map[string]interface{}{
		"backup_id": backup.ID,
		"size":      backup.Size,
		"interface": req.Interface,
	})

	return backup, nil
}

// RestoreConfig restores a configuration from backup
func (cs *ConfigService) RestoreConfig(req *types.ConfigRestoreRequest) (*types.ConfigRestoreResponse, error) {
	cs.logger.Info("Restoring configuration from backup", map[string]interface{}{
		"backup_id": req.BackupID,
		"user_id":   req.UserID,
	})

	// Get backup (in production, this would fetch from storage)
	backup, err := cs.getBackup(req.BackupID)
	if err != nil {
		return &types.ConfigRestoreResponse{
			Success: false,
			Error: &types.ErrorDetail{
				Code:    "BACKUP_NOT_FOUND",
				Message: "Backup not found",
				Details: err.Error(),
			},
			Message: "Failed to restore configuration",
			Code:    404,
		}, nil
	}

	// Validate backup
	if err := cs.validateBackup(backup); err != nil {
		return &types.ConfigRestoreResponse{
			Success: false,
			Error: &types.ErrorDetail{
				Code:    "INVALID_BACKUP",
				Message: "Invalid backup",
				Details: err.Error(),
			},
			Message: "Failed to restore configuration",
			Code:    400,
		}, nil
	}

	// Create backup of current configuration if requested
	var currentBackup *types.ConfigBackup
	if req.Options.Backup {
		currentConfigReq := &types.ConfigRequest{
			Interface: backup.Metadata["interface"].(string),
			UserID:    req.UserID,
			SessionID: req.SessionID,
			Options: types.ConfigOptions{
				Backup: true,
			},
		}
		currentBackup, err = cs.CreateBackup(currentConfigReq)
		if err != nil {
			cs.logger.Error("Failed to create backup before restore", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}

	// Restore configuration
	if err := cs.applyConfiguration(backup.Config); err != nil {
		return &types.ConfigRestoreResponse{
			Success: false,
			Error: &types.ErrorDetail{
				Code:    "RESTORE_FAILED",
				Message: "Failed to restore configuration",
				Details: err.Error(),
			},
			Message: "Configuration restore failed",
			Code:    500,
		}, nil
	}

	cs.logger.Info("Configuration restored successfully", map[string]interface{}{
		"backup_id": req.BackupID,
		"user_id":   req.UserID,
	})

	return &types.ConfigRestoreResponse{
		Success:  true,
		Config:   backup.Config,
		Backup:   currentBackup,
		Message:  "Configuration restored successfully",
		Code:     200,
		Warnings: []string{},
	}, nil
}

// ListConfigs lists available configurations
func (cs *ConfigService) ListConfigs(req *types.ConfigListRequest) ([]types.ConfigSummary, error) {
	cs.logger.Info("Listing configurations", map[string]interface{}{
		"interface": req.Interface,
		"user_id":   req.UserID,
	})

	// In production, this would fetch from storage
	// For now, return mock data
	configs := []types.ConfigSummary{
		{
			ID:          "config_1",
			Name:        "Default Configuration",
			Type:        "setup",
			Interface:   req.Interface,
			Environment: "windows",
			Version:     "1.0.0",
			CreatedAt:   time.Now().Add(-24 * time.Hour),
			UpdatedAt:   time.Now(),
			Size:        1024,
			Checksum:    "abc123",
			Status:      string(types.ConfigStatusActive),
		},
		{
			ID:          "config_2",
			Name:        "Backup Configuration",
			Type:        "setup",
			Interface:   req.Interface,
			Environment: "windows",
			Version:     "1.0.0",
			CreatedAt:   time.Now().Add(-48 * time.Hour),
			UpdatedAt:   time.Now().Add(-24 * time.Hour),
			Size:        1024,
			Checksum:    "def456",
			Status:      string(types.ConfigStatusInactive),
		},
	}

	cs.logger.Info("Configurations listed", map[string]interface{}{
		"count":     len(configs),
		"interface": req.Interface,
	})

	return configs, nil
}

// GetTemplate gets a configuration template
func (cs *ConfigService) GetTemplate(interfaceType, environment, templateName string) (*types.ConfigTemplate, error) {
	cs.logger.Info("Getting configuration template", map[string]interface{}{
		"interface":   interfaceType,
		"environment": environment,
		"template":    templateName,
	})

	// Generate template based on interface and environment
	template := &types.ConfigTemplate{
		Name:        templateName,
		Description: fmt.Sprintf("Configuration template for %s on %s", interfaceType, environment),
		Version:     "1.0.0",
		Interface:   interfaceType,
		Environment: environment,
		Content:     cs.generateTemplateContent(interfaceType, environment),
		Variables:   cs.generateTemplateVariables(interfaceType, environment),
		Validation:  cs.generateTemplateValidation(),
		Metadata: map[string]interface{}{
			"created_by": "config_service",
			"created_at": time.Now(),
		},
	}

	return template, nil
}

// Helper methods for configuration generation

func (cs *ConfigService) generateManagerConfig(req *types.ConfigRequest) types.ManagerConfig {
	homeDir := req.Environment.HomeDir
	if homeDir == "" {
		homeDir = filepath.Join(os.Getenv("HOME"), ".syntropy")
	}

	return types.ManagerConfig{
		HomeDir:     homeDir,
		LogLevel:    "info",
		APIEndpoint: "http://localhost:8080",
		Directories: map[string]string{
			"config":  filepath.Join(homeDir, "config"),
			"keys":    filepath.Join(homeDir, "keys"),
			"logs":    filepath.Join(homeDir, "logs"),
			"cache":   filepath.Join(homeDir, "cache"),
			"backups": filepath.Join(homeDir, "backups"),
		},
		DefaultPaths: map[string]string{
			"manager_config": filepath.Join(homeDir, "config", "manager.yaml"),
			"owner_key":      filepath.Join(homeDir, "keys", "owner.key"),
			"owner_pub":      filepath.Join(homeDir, "keys", "owner.key.pub"),
		},
		Database: types.DatabaseConfig{
			Type:     "sqlite",
			Host:     "localhost",
			Port:     0,
			Name:     filepath.Join(homeDir, "syntropy.db"),
			Username: "",
			SSLMode:  "disable",
		},
	}
}

func (cs *ConfigService) generateOwnerKey(req *types.ConfigRequest) types.OwnerKey {
	// Generate RSA key pair
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		cs.logger.Error("Failed to generate owner key", map[string]interface{}{
			"error": err.Error(),
		})
		// Return empty key - this will be caught by validation
		return types.OwnerKey{}
	}

	// Extract public key
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		cs.logger.Error("Failed to marshal public key", map[string]interface{}{
			"error": err.Error(),
		})
		return types.OwnerKey{}
	}

	// Encode public key to PEM
	pubKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	})

	return types.OwnerKey{
		Type:      "RSA",
		Path:      filepath.Join(req.Environment.HomeDir, ".syntropy", "keys", "owner.key"),
		PublicKey: string(pubKeyPEM),
		CreatedAt: time.Now(),
		Algorithm: "RSA",
		Size:      2048,
	}
}

func (cs *ConfigService) generateEnvironmentConfig(req *types.ConfigRequest) types.Environment {
	return types.Environment{
		OS:           req.Environment.OS,
		Architecture: req.Environment.Architecture,
		HomeDir:      req.Environment.HomeDir,
		Variables:    req.Environment.EnvironmentVars,
		Features:     req.Environment.Features,
	}
}

func (cs *ConfigService) generateInterfaceConfig(req *types.ConfigRequest) types.InterfaceConfig {
	return types.InterfaceConfig{
		Type:     req.Interface,
		Theme:    "default",
		Language: "en",
		Settings: map[string]interface{}{
			"auto_update":   true,
			"notifications": true,
			"log_level":     "info",
		},
		Permissions: []string{
			"read_config",
			"write_config",
			"manage_keys",
			"access_network",
		},
	}
}

func (cs *ConfigService) generateSecurityConfig(req *types.ConfigRequest) types.SecurityConfig {
	return types.SecurityConfig{
		EncryptionAlgorithm: "AES-256-GCM",
		KeyRotationDays:     90,
		AllowedIPs:          []string{"127.0.0.1", "::1"},
		SSLEnabled:          true,
		CertPath:            filepath.Join(req.Environment.HomeDir, ".syntropy", "certs", "server.crt"),
	}
}

func (cs *ConfigService) generateNetworkConfig(req *types.ConfigRequest) types.NetworkConfig {
	return types.NetworkConfig{
		Port:        8080,
		Host:        "localhost",
		Endpoints:   []string{"http://localhost:8080"},
		Timeout:     30,
		Retries:     3,
		Compression: true,
	}
}

func (cs *ConfigService) generateConfigMetadata(req *types.ConfigRequest) types.ConfigMetadata {
	return types.ConfigMetadata{
		Version:     "1.0.0",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		CreatedBy:   req.UserID,
		Interface:   req.Interface,
		Environment: req.Environment.OS,
		Checksum:    "",
	}
}

func (cs *ConfigService) validateGeneratedConfig(config *types.SetupConfig, req *types.ConfigRequest) error {
	// Basic validation
	if config.Manager.HomeDir == "" {
		return fmt.Errorf("manager home directory is required")
	}

	if config.OwnerKey.Path == "" {
		return fmt.Errorf("owner key path is required")
	}

	if config.OwnerKey.PublicKey == "" {
		return fmt.Errorf("owner public key is required")
	}

	return nil
}

func (cs *ConfigService) calculateConfigSize(config *types.SetupConfig) int64 {
	// Simplified size calculation
	// In production, you would serialize and measure the actual size
	return 1024 // Mock size
}

func (cs *ConfigService) calculateConfigChecksum(config *types.SetupConfig) string {
	// Simplified checksum calculation
	// In production, you would calculate a proper hash
	return fmt.Sprintf("checksum_%d", time.Now().Unix())
}

func (cs *ConfigService) getBackup(backupID string) (*types.ConfigBackup, error) {
	// In production, this would fetch from storage
	// For now, return a mock backup
	return &types.ConfigBackup{
		ID:          backupID,
		Name:        "Mock Backup",
		Description: "Mock backup for testing",
		Timestamp:   time.Now(),
		Size:        1024,
		Checksum:    "mock_checksum",
		Encrypted:   false,
		Compressed:  false,
		Metadata:    map[string]interface{}{},
	}, nil
}

func (cs *ConfigService) validateBackup(backup *types.ConfigBackup) error {
	if backup.ID == "" {
		return fmt.Errorf("backup ID is required")
	}

	if backup.Config == nil {
		return fmt.Errorf("backup configuration is required")
	}

	return nil
}

func (cs *ConfigService) applyConfiguration(config *types.SetupConfig) error {
	// In production, this would apply the configuration to the system
	cs.logger.Info("Applying configuration", map[string]interface{}{
		"interface":   config.Interface.Type,
		"environment": config.Environment.OS,
	})

	return nil
}

func (cs *ConfigService) generateTemplateContent(interfaceType, environment string) string {
	// Generate template content based on interface and environment
	return fmt.Sprintf(`
# Configuration Template for %s on %s
interface:
  type: %s
  theme: default
  language: en

environment:
  os: %s
  home_dir: ~/.syntropy

manager:
  log_level: info
  api_endpoint: http://localhost:8080

security:
  encryption_algorithm: AES-256-GCM
  key_rotation_days: 90

network:
  port: 8080
  host: localhost
  compression: true
`, interfaceType, environment, interfaceType, environment)
}

func (cs *ConfigService) generateTemplateVariables(interfaceType, environment string) []types.TemplateVariable {
	return []types.TemplateVariable{
		{
			Name:        "home_dir",
			Type:        "string",
			Default:     "~/.syntropy",
			Required:    true,
			Description: "Home directory for Syntropy configuration",
			Validation:  "path",
			Sensitive:   false,
		},
		{
			Name:        "log_level",
			Type:        "string",
			Default:     "info",
			Required:    false,
			Description: "Logging level",
			Validation:  "enum",
			Options:     []string{"debug", "info", "warn", "error"},
			Sensitive:   false,
		},
		{
			Name:        "api_endpoint",
			Type:        "string",
			Default:     "http://localhost:8080",
			Required:    true,
			Description: "API endpoint URL",
			Validation:  "url",
			Sensitive:   false,
		},
	}
}

func (cs *ConfigService) generateTemplateValidation() *types.TemplateValidation {
	return &types.TemplateValidation{
		Schema: `{
			"type": "object",
			"properties": {
				"interface": {"type": "object"},
				"environment": {"type": "object"},
				"manager": {"type": "object"},
				"security": {"type": "object"},
				"network": {"type": "object"}
			},
			"required": ["interface", "environment", "manager"]
		}`,
		Required: []string{"interface", "environment", "manager"},
		Optional: []string{"security", "network"},
		Defaults: map[string]interface{}{
			"log_level": "info",
			"port":      8080,
		},
		Metadata: map[string]interface{}{
			"version": "1.0.0",
		},
	}
}
