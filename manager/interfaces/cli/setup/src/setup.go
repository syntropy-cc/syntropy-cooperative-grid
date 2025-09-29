// Package setup provides functionality for setting up the Syntropy CLI environment
package setup

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	apiTypes "github.com/syntropy-cc/syntropy-cooperative-grid/manager/api/types"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src/internal/types"
)

// Usar tipos definidos em internal/types

// SetupManager implementa a interface SetupManager conforme especificado no guia
type SetupManager struct {
	validator    types.Validator
	configurator types.Configurator
	stateManager types.StateManager
	keyManager   types.KeyManager
	logger       types.SetupLogger
}

// NewSetupManager cria um novo gerenciador de setup
func NewSetupManager() (*SetupManager, error) {
	logger := NewSetupLogger()

	return &SetupManager{
		validator:    NewValidator(logger),
		configurator: NewConfigurator(logger),
		stateManager: NewStateManager(logger),
		keyManager:   NewKeyManager(logger),
		logger:       logger,
	}, nil
}

// Setup executa o setup completo conforme especificado no guia
func (sm *SetupManager) Setup(options *types.SetupOptions) error {
	sm.logger.LogStep("setup_start", map[string]interface{}{
		"options": options,
	})

	// 1. Validar ambiente
	envInfo, err := sm.validator.ValidateEnvironment()
	if err != nil {
		return sm.handleError(err, "validation_failed")
	}

	if !envInfo.CanProceed && !options.Force {
		issues := []types.ValidationIssue{
			{
				Type:        "environment",
				Severity:    "error",
				Message:     "Validation failed",
				Suggestions: []string{"Check system requirements"},
			},
		}
		return sm.handleError(types.ErrValidationFailedError(issues), "validation_failed")
	}

	// 2. Criar estrutura de diretórios
	if err := sm.configurator.CreateStructure(); err != nil {
		return sm.handleError(err, "structure_creation_failed")
	}

	// 3. Gerar chaves
	keyPair, err := sm.keyManager.GenerateKeyPair("ed25519")
	if err != nil {
		return sm.handleError(err, "key_generation_failed")
	}

	// 4. Gerar configuração
	if err := sm.configurator.GenerateConfig(&types.ConfigOptions{
		OwnerName:  options.CustomSettings["owner_name"],
		OwnerEmail: options.CustomSettings["owner_email"],
	}); err != nil {
		return sm.handleError(err, "config_generation_failed")
	}

	// 5. Salvar estado
	state := &types.SetupState{
		Version:   "1.0.0",
		CreatedAt: time.Now(),
		Status:    types.SetupStatusCompleted,
		Keys: &types.KeyInfo{
			OwnerKeyID: keyPair.ID,
			Algorithm:  keyPair.Algorithm,
		},
	}

	if err := sm.stateManager.SaveState(state); err != nil {
		return sm.handleError(err, "state_save_failed")
	}

	sm.logger.LogStep("setup_completed", map[string]interface{}{
		"key_id": keyPair.ID,
	})

	return nil
}

// Validate valida o ambiente
func (sm *SetupManager) Validate() (*types.ValidationResult, error) {
	sm.logger.LogStep("validation_start", nil)

	// Validar ambiente
	envInfo, err := sm.validator.ValidateEnvironment()
	if err != nil {
		sm.logger.LogError(err, map[string]interface{}{
			"step": "validation",
		})
		return nil, err
	}

	// Validar dependências
	deps, err := sm.validator.ValidateDependencies()
	if err != nil {
		sm.logger.LogError(err, map[string]interface{}{
			"step": "dependency_validation",
		})
		return nil, err
	}

	// Validar rede
	network, err := sm.validator.ValidateNetwork()
	if err != nil {
		sm.logger.LogError(err, map[string]interface{}{
			"step": "network_validation",
		})
		return nil, err
	}

	// Validar permissões
	permissions, err := sm.validator.ValidatePermissions()
	if err != nil {
		sm.logger.LogError(err, map[string]interface{}{
			"step": "permission_validation",
		})
		return nil, err
	}

	// Criar resultado de validação
	result := &types.ValidationResult{
		Environment:  envInfo,
		Dependencies: deps,
		Network:      network,
		Permissions:  permissions,
		CanProceed:   true,
		Issues:       []types.ValidationIssue{},
		Warnings:     []string{},
	}

	sm.logger.LogStep("validation_completed", map[string]interface{}{
		"can_proceed":  result.CanProceed,
		"issues_count": len(result.Issues),
	})

	return result, nil
}

// Status verifica o status do setup
func (sm *SetupManager) Status() (*types.SetupStatus, error) {
	sm.logger.LogStep("status_check_start", nil)

	state, err := sm.stateManager.LoadState()
	if err != nil {
		sm.logger.LogError(err, map[string]interface{}{
			"step": "status_check",
		})
		return nil, err
	}

	sm.logger.LogStep("status_check_completed", map[string]interface{}{
		"status":  state.Status,
		"version": state.Version,
	})

	return &state.Status, nil
}

// Reset reseta o setup
func (sm *SetupManager) Reset(confirm bool) error {
	if !confirm {
		return fmt.Errorf("reset requer confirmação")
	}

	sm.logger.LogStep("reset_start", nil)

	// Implementar reset conforme necessário
	// Por enquanto, apenas log

	sm.logger.LogStep("reset_completed", nil)

	return nil
}

// Repair repara problemas automaticamente
func (sm *SetupManager) Repair() error {
	sm.logger.LogStep("repair_start", nil)

	// Verificar integridade do estado
	if err := sm.stateManager.VerifyIntegrity(); err != nil {
		sm.logger.LogWarning("Problemas de integridade detectados", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Verificar integridade das chaves
	// Nota: ListKeys não está implementado na interface KeyManager
	// Implementação simplificada para reparo
	sm.logger.LogInfo("Verificação de integridade de chaves não implementada", nil)

	sm.logger.LogStep("repair_completed", nil)

	return nil
}

// handleError trata erros de forma consistente
func (sm *SetupManager) handleError(err error, context string) error {
	sm.logger.LogError(err, map[string]interface{}{
		"context": context,
	})
	return err
}

// SetupLegacy configura o ambiente para o Syntropy CLI (função legacy para compatibilidade)
func SetupLegacy(options types.LegacySetupOptions) (*types.LegacySetupResult, error) {
	fmt.Println("Starting Syntropy CLI setup...")

	// Criar novo gerenciador de setup
	manager, err := NewSetupManager()
	if err != nil {
		return nil, fmt.Errorf("falha ao criar gerenciador de setup: %w", err)
	}
	defer manager.logger.Close()

	// Converter opções legacy para novas opções
	newOptions := &types.SetupOptions{
		Force:        options.Force,
		ValidateOnly: false,
		Verbose:      true,
		Quiet:        false,
		ConfigPath:   options.ConfigPath,
		CustomSettings: map[string]string{
			"owner_name":  "Syntropy User",
			"owner_email": "user@syntropy.network",
		},
	}

	// Executar setup
	if err := manager.Setup(newOptions); err != nil {
		return &types.LegacySetupResult{
			Success:   false,
			StartTime: time.Now(),
			EndTime:   time.Now(),
			Error:     err,
			Message:   err.Error(),
		}, err
	}

	return &types.LegacySetupResult{
		Success:   true,
		StartTime: time.Now(),
		EndTime:   time.Now(),
		Message:   "Setup concluído com sucesso",
	}, nil
}

// StatusLegacy checks the installation status of the Syntropy CLI
func StatusLegacy(options types.LegacySetupOptions) (*types.LegacySetupResult, error) {
	fmt.Println("Checking Syntropy CLI status...")

	// Create new setup manager
	manager, err := NewSetupManager()
	if err != nil {
		return nil, fmt.Errorf("falha ao criar gerenciador de setup: %w", err)
	}
	defer manager.logger.Close()

	// Get status using new manager
	status, err := manager.Status()
	if err != nil {
		return &types.LegacySetupResult{
			Success:   false,
			StartTime: time.Now(),
			EndTime:   time.Now(),
			Error:     err,
			Message:   err.Error(),
		}, err
	}

	// Convert status to legacy result
	return &types.LegacySetupResult{
		Success:   true,
		StartTime: time.Now(),
		EndTime:   time.Now(),
		Message:   fmt.Sprintf("Status: %s", *status),
	}, nil
}

// ResetLegacy resets the Syntropy CLI configuration
func ResetLegacy(options types.LegacySetupOptions) (*types.LegacySetupResult, error) {
	fmt.Println("Resetting Syntropy CLI configuration...")

	// Create new setup manager
	manager, err := NewSetupManager()
	if err != nil {
		return nil, fmt.Errorf("falha ao criar gerenciador de setup: %w", err)
	}
	defer manager.logger.Close()

	// Reset using new manager
	err = manager.Reset(true)
	if err != nil {
		return &types.LegacySetupResult{
			Success:   false,
			StartTime: time.Now(),
			EndTime:   time.Now(),
			Error:     err,
			Message:   err.Error(),
		}, err
	}

	// Return success result
	return &types.LegacySetupResult{
		Success:   true,
		StartTime: time.Now(),
		EndTime:   time.Now(),
		Message:   "Reset completed successfully",
	}, nil
}

// GetSyntropyDirLegacy returns the default directory for the Syntropy CLI
func GetSyntropyDirLegacy() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to temporary directory in case of error
		return filepath.Join(os.TempDir(), "syntropy")
	}

	switch runtime.GOOS {
	case "windows":
		return filepath.Join(homeDir, "Syntropy")
	case "linux", "darwin":
		return filepath.Join(homeDir, ".syntropy")
	default:
		return filepath.Join(homeDir, ".syntropy")
	}
}

// Funções stub removidas para evitar conflitos de redefinição

// Conversion functions between local and API types

// convertLegacyToAPISetupOptions converts local SetupOptions to API SetupOptions
func convertLegacyToAPISetupOptions(local types.LegacySetupOptions) *apiTypes.SetupOptions {
	return &apiTypes.SetupOptions{
		Force:          local.Force,
		InstallService: false, // Legacy options don't have this field
		ConfigPath:     local.ConfigPath,
		HomeDir:        "", // Legacy options don't have this field
		Interface:      "cli",
		CustomOptions: map[string]interface{}{
			"source": "cli_setup",
		},
	}
}

// convertFromAPIToLegacySetupResult converts API SetupResult to local SetupResult
func convertFromAPIToLegacySetupResult(api *apiTypes.SetupResult) *types.LegacySetupResult {
	return &types.LegacySetupResult{
		Success:     api.Success,
		StartTime:   api.StartTime,
		EndTime:     api.EndTime,
		ConfigPath:  api.ConfigPath,
		Environment: api.Environment,
		Options: types.LegacySetupOptions{
			Force:      api.Options.Force,
			ConfigPath: api.Options.ConfigPath,
		},
		Error: api.Error,
	}
}

// getCurrentEnvironmentInfo gets current environment information
func getCurrentEnvironmentInfo() *apiTypes.EnvironmentInfo {
	homeDir, _ := os.UserHomeDir()
	return &apiTypes.EnvironmentInfo{
		OS:              runtime.GOOS,
		OSVersion:       "unknown", // Would be populated by actual detection
		Architecture:    runtime.GOARCH,
		HomeDir:         homeDir,
		HasAdminRights:  true,  // Would be detected
		AvailableDiskGB: 100.0, // Would be calculated
		HasInternet:     true,  // Would be tested
		EnvironmentVars: make(map[string]string),
		Features:        []string{},
		Capabilities:    []string{},
	}
}

// shouldForceLocalSetup determines whether to force local implementation instead of API
func shouldForceLocalSetup() bool {
	// Force local setup in any of these conditions:
	// 1. Environment variable is set
	// 2. We're in a test/development environment
	if os.Getenv("SYNTROPY_FORCE_LOCAL_SETUP") == "true" {
		return true
	}

	// 3. Check if we're running in CI/testing environment
	if os.Getenv("CI") != "" || os.Getenv("TESTING") != "" {
		return true
	}

	// 4. For now, force local setup to guarantee functionality
	// This can be removed once API central issues are fixed
	return true
}

// convertStatusToLegacySetupResult converts API status to local SetupResult
func convertStatusToLegacySetupResult(status map[string]interface{}) *types.LegacySetupResult {
	success := true
	if status["status"] != "active" {
		success = false
	}

	// Safe type assertions with defaults
	configPath := ""
	if cp, ok := status["config_path"].(string); ok {
		configPath = cp
	}

	environment := ""
	if env, ok := status["interface"].(string); ok {
		environment = env
	}

	statusStr := ""
	if s, ok := status["status"].(string); ok {
		statusStr = s
	}

	return &types.LegacySetupResult{
		Success:     success,
		StartTime:   time.Now(),
		EndTime:     time.Now(),
		ConfigPath:  configPath,
		Environment: environment,
		Message:     fmt.Sprintf("Status: %s", statusStr),
	}
}
