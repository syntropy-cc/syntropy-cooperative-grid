// Package setup provides functionality for setting up the Syntropy CLI environment
package setup

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"setup-component/src/internal/types"
)

// Public types for external use

// SetupError represents a setup error (re-export from types)
type SetupError = types.SetupError

// LegacySetupOptions defines the options for the setup process (legacy compatibility)
type LegacySetupOptions struct {
	Force          bool   // Force setup even if validations fail
	InstallService bool   // Install system service
	ConfigPath     string // Custom configuration file path
	HomeDir        string // Custom home directory
}

// LegacySetupResult contains the result of the setup process (legacy compatibility)
type LegacySetupResult struct {
	Success     bool               // Indicates if the setup was successful
	StartTime   time.Time          // Setup start time
	EndTime     time.Time          // Setup end time
	ConfigPath  string             // Configuration file path
	Environment string             // Environment (windows, linux, darwin)
	Options     LegacySetupOptions // Options used in the setup
	Error       error              // Error, if any
	Message     string             // Human-readable message
}

// Usar tipos definidos em internal/types

// SetupManager implementa a interface SetupManager conforme especificado no guia
type SetupManager struct {
	validator    types.Validator
	configurator types.Configurator
	stateManager types.StateManager
	keyManager   types.KeyManager
	tokenManager types.TokenManager
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
		tokenManager: NewTokenManager(logger),
		logger:       logger,
	}, nil
}

// Setup executa o setup completo conforme especificado no guia
func (sm *SetupManager) Setup(options *types.SetupOptions) error {
	// Check for nil options
	if options == nil {
		return sm.handleError(fmt.Errorf("setup options cannot be nil"), "invalid_options")
	}

	sm.logger.LogStep("setup_start", map[string]interface{}{
		"options": options,
	})

	// 1. Validar ambiente
	envInfo, err := sm.validator.ValidateEnvironmentWithOptions(options)
	if err != nil {
		return sm.handleError(err, "validation_failed")
	}

	// Verificar se setup já existe
	existingState, err := sm.stateManager.LoadState()
	if err == nil && existingState.Status == types.SetupStatusCompleted {
		if !options.Force {
			return types.ErrSetupAlreadyExistsError(existingState.CreatedAt, existingState.Version)
		}

		// Se --force, criar backup automático antes de sobrescrever
		sm.logger.LogInfo("Setup existente detectado, criando backup automático...", nil)
		backupName := fmt.Sprintf("pre_override_%d", time.Now().Unix())
		if err := sm.createFullBackup(backupName); err != nil {
			sm.logger.LogWarning("Falha ao criar backup automático", map[string]interface{}{
				"error": err.Error(),
			})
		}
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

	// 3. Gerar ou carregar chaves existentes
	keyPair, err := sm.keyManager.GenerateOrLoadKeyPair("ed25519")
	if err != nil {
		return sm.handleError(err, "key_generation_failed")
	}

	// 4. Gerar Grid Token automaticamente (a menos que seja desabilitado)
	var gridToken string
	shouldGenerateToken := true

	// Verificar se geração de token foi explicitamente desabilitada
	if options.CustomSettings["generate_grid_token"] == "false" {
		shouldGenerateToken = false
	}

	// Verificar se token já existe
	tokenExists, err := sm.tokenManager.TokenExists()
	if err == nil && tokenExists {
		sm.logger.LogInfo("Grid Token already exists, skipping generation", nil)
		shouldGenerateToken = false
	}

	if shouldGenerateToken {
		sm.logger.LogStep("grid_token_generation_start", map[string]interface{}{
			"keyring_available": true, // Será verificado internamente
		})

		gridToken, err = sm.tokenManager.GenerateToken()
		if err != nil {
			return sm.handleError(err, "grid_token_generation_failed")
		}

		// Salvar token no keyring do sistema
		if err := sm.tokenManager.SaveToken(gridToken); err != nil {
			return sm.handleError(err, "grid_token_save_failed")
		}

		sm.logger.LogStep("grid_token_generation_completed", map[string]interface{}{
			"token_preview":  gridToken[:8] + "...[HIDDEN]",
			"storage_method": "keyring",
		})
	}

	// 5. Gerar configuração
	if err := sm.configurator.GenerateConfig(&types.ConfigOptions{
		OwnerName:  options.CustomSettings["owner_name"],
		OwnerEmail: options.CustomSettings["owner_email"],
	}); err != nil {
		return sm.handleError(err, "config_generation_failed")
	}

	// 6. Salvar estado
	state := &types.SetupState{
		Version:   "1.0.0",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    types.SetupStatusCompleted,
		Environment: &types.EnvironmentInfo{
			OS:           envInfo.OS,
			OSVersion:    envInfo.OSVersion,
			Architecture: envInfo.Architecture,
			HomeDir:      envInfo.HomeDir,
			CanProceed:   envInfo.CanProceed,
		},
		Keys: &types.KeyInfo{
			OwnerKeyID: keyPair.ID,
			Algorithm:  keyPair.Algorithm,
			CreatedAt:  keyPair.CreatedAt,
			ExpiresAt:  keyPair.ExpiresAt,
		},
		Metadata: map[string]string{
			"setup_version": "1.0.0",
			"setup_method":  "automated",
		},
	}

	// Adicionar informações do Grid Token aos metadados se foi gerado
	if gridToken != "" {
		state.Metadata["grid_token_generated"] = "true"
		state.Metadata["grid_token_preview"] = gridToken[:8] + "...[HIDDEN]"
		state.Metadata["grid_token_storage"] = "keyring"
	}

	if err := sm.stateManager.SaveState(state); err != nil {
		return sm.handleError(err, "state_save_failed")
	}

	sm.logger.LogStep("setup_completed", map[string]interface{}{
		"state": state,
	})

	return nil
}

// SetupWithPublicOptions executa o setup usando as opções públicas
func (sm *SetupManager) SetupWithPublicOptions(options *SetupOptions) error {
	if options == nil {
		return sm.handleError(fmt.Errorf("setup options cannot be nil"), "invalid_options")
	}

	// Convert public SetupOptions to internal types.SetupOptions
	internalOptions := &types.SetupOptions{
		Force:          options.Force,
		ValidateOnly:   options.ValidateOnly,
		TestMode:       options.TestMode,
		Verbose:        options.Verbose,
		Quiet:          options.Quiet,
		ConfigPath:     options.ConfigPath,
		CustomSettings: options.CustomSettings,
	}

	return sm.Setup(internalOptions)
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

	// Criar backup automático antes de deletar
	sm.logger.LogInfo("Criando backup automático antes de deletar...", nil)
	backupName := fmt.Sprintf("pre_reset_%d", time.Now().Unix())
	if err := sm.createFullBackup(backupName); err != nil {
		sm.logger.LogWarning("Falha ao criar backup", map[string]interface{}{
			"error": err.Error(),
		})
	} else {
		homeDir, _ := os.UserHomeDir()
		backupPath := filepath.Join(homeDir, ".syntropy", "backups")
		sm.logger.LogInfo("Backup criado com sucesso", map[string]interface{}{
			"backup_location": backupPath,
			"backup_name":     backupName,
		})
	}

	// Deletar Grid Token antes de remover arquivos
	sm.logger.LogInfo("Removendo Grid Token...", nil)
	if err := sm.tokenManager.DeleteToken(); err != nil {
		sm.logger.LogWarning("Falha ao remover Grid Token", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Remover arquivo de estado
	homeDir, _ := os.UserHomeDir()
	statePath := filepath.Join(homeDir, ".syntropy", "state", "setup_state.json")
	if _, err := os.Stat(statePath); err == nil {
		if err := os.Remove(statePath); err != nil {
			sm.logger.LogWarning("Falha ao remover arquivo de estado", map[string]interface{}{
				"state_path": statePath,
				"error":      err.Error(),
			})
		} else {
			sm.logger.LogInfo("Arquivo de estado removido", map[string]interface{}{
				"state_path": statePath,
			})
		}
	}

	// Remover diretório de configuração
	syntropyDir := filepath.Join(homeDir, ".syntropy")
	configDir := filepath.Join(syntropyDir, "config")

	if _, err := os.Stat(configDir); err == nil {
		if err := os.RemoveAll(configDir); err != nil {
			sm.logger.LogWarning("Falha ao remover diretório de configuração", map[string]interface{}{
				"config_dir": configDir,
				"error":      err.Error(),
			})
		} else {
			sm.logger.LogInfo("Diretório de configuração removido", map[string]interface{}{
				"config_dir": configDir,
			})
		}
	}

	// Remover diretório de chaves
	keysDir := filepath.Join(syntropyDir, "keys")
	if _, err := os.Stat(keysDir); err == nil {
		if err := os.RemoveAll(keysDir); err != nil {
			sm.logger.LogWarning("Falha ao remover diretório de chaves", map[string]interface{}{
				"keys_dir": keysDir,
				"error":    err.Error(),
			})
		} else {
			sm.logger.LogInfo("Diretório de chaves removido", map[string]interface{}{
				"keys_dir": keysDir,
			})
		}
	}

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
func SetupLegacy(options LegacySetupOptions) (*LegacySetupResult, error) {
	fmt.Println("Starting Syntropy CLI setup...")

	// Criar novo gerenciador de setup
	manager, err := NewSetupManager()
	if err != nil {
		return nil, fmt.Errorf("falha ao criar gerenciador de setup: %w", err)
	}
	defer manager.logger.Close()

	// Verificar se já existe um setup
	existingState, err := manager.stateManager.LoadState()
	if err == nil && existingState.Status == types.SetupStatusCompleted {
		// Setup já existe, perguntar se deve substituir
		if !options.Force {
			fmt.Println("⚠️  Já existe uma configuração do Syntropy Manager.")
			fmt.Printf("   📁 Configuração atual: %s\n", filepath.Join(os.Getenv("HOME"), ".syntropy"))
			fmt.Print("   Deseja substituí-la? (y/N): ")
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				fmt.Println("Setup cancelado pelo usuário.")
				return &LegacySetupResult{
					Success:   false,
					StartTime: time.Now(),
					EndTime:   time.Now(),
					Message:   "Setup cancelado pelo usuário",
				}, nil
			}
		}

		// Criar backup do setup existente
		fmt.Println("📦 Criando backup do setup existente...")
		backupName := fmt.Sprintf("pre_setup_%d", time.Now().Unix())
		if err := manager.stateManager.BackupState(backupName); err != nil {
			fmt.Printf("⚠️  Aviso: Falha ao criar backup: %v\n", err)
		} else {
			fmt.Printf("✅ Backup criado: %s\n", backupName)
		}

		// Fazer backup de todas as pastas exceto backups
		homeDir, _ := os.UserHomeDir()
		syntropyDir := filepath.Join(homeDir, ".syntropy")
		backupDir := filepath.Join(syntropyDir, "backups", "full_backup")

		if err := os.MkdirAll(backupDir, 0755); err == nil {
			backupPath := filepath.Join(backupDir, fmt.Sprintf("backup_%d", time.Now().Unix()))
			if err := backupAllDirectories(syntropyDir, backupPath); err != nil {
				fmt.Printf("⚠️  Aviso: Falha ao fazer backup completo: %v\n", err)
			} else {
				fmt.Printf("✅ Backup completo criado: %s\n", backupPath)
				fmt.Printf("🔒 AVISO DE SEGURANÇA: Os backups contêm chaves criptográficas sensíveis!\n")
				fmt.Printf("   - Gerencie os backups com cuidado\n")
				fmt.Printf("   - Considere criptografar os backups\n")
				fmt.Printf("   - Remova backups antigos regularmente\n")
				fmt.Printf("   - Nunca compartilhe backups não criptografados\n")
			}
		}
	}

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

	// Executar setup com fallback para validação
	if err := manager.Setup(newOptions); err != nil {
		// Se a validação falhou e não estamos forçando, tentar novamente com force
		if !options.Force && strings.Contains(err.Error(), "Falha na validação do ambiente") {
			fmt.Println("⚠️  Validação do ambiente falhou, mas prosseguindo com setup básico...")
			newOptions.Force = true
			if err := manager.Setup(newOptions); err != nil {
				return &LegacySetupResult{
					Success:   false,
					StartTime: time.Now(),
					EndTime:   time.Now(),
					Error:     err,
					Message:   err.Error(),
				}, err
			}
		} else {
			return &LegacySetupResult{
				Success:   false,
				StartTime: time.Now(),
				EndTime:   time.Now(),
				Error:     err,
				Message:   err.Error(),
			}, err
		}
	}

	// Obter caminho da configuração
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, ".syntropy", "config", "manager.yaml")

	return &LegacySetupResult{
		Success:     true,
		StartTime:   time.Now(),
		EndTime:     time.Now(),
		ConfigPath:  configPath,
		Environment: fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		Message:     "Setup concluído com sucesso",
	}, nil
}

// StatusLegacy checks the installation status of the Syntropy CLI
func StatusLegacy(options LegacySetupOptions) (*LegacySetupResult, error) {
	fmt.Println("Checking Syntropy CLI status...")

	// Create new setup manager
	manager, err := NewSetupManager()
	if err != nil {
		return nil, fmt.Errorf("falha ao criar gerenciador de setup: %w", err)
	}
	defer manager.logger.Close()

	// Check if setup actually exists by trying to load state
	state, err := manager.stateManager.LoadState()
	if err != nil {
		// Check if this is specifically a "file not found" error (setup not run yet)
		if setupErr, ok := err.(*types.SetupError); ok && setupErr.Code == types.ErrStateLoad {
			// Check if the error message indicates file not found
			if setupErr.Cause != nil && strings.Contains(setupErr.Cause.Error(), "arquivo de estado não encontrado") {
				return &LegacySetupResult{
					Success:   false,
					StartTime: time.Now(),
					EndTime:   time.Now(),
					Message:   "Setup não foi executado ainda",
				}, nil
			}
		}

		// For other errors (corruption, permission issues, etc.)
		return &LegacySetupResult{
			Success:   false,
			StartTime: time.Now(),
			EndTime:   time.Now(),
			Error:     err,
			Message:   "Setup não encontrado ou corrompido",
		}, nil
	}

	// Check if setup is actually completed (not just initial state)
	if state.Status != types.SetupStatusCompleted {
		return &LegacySetupResult{
			Success:   false,
			StartTime: time.Now(),
			EndTime:   time.Now(),
			Message:   fmt.Sprintf("Setup não concluído. Status atual: %s", state.Status),
		}, nil
	}

	// Get environment info
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, ".syntropy", "config", "manager.yaml")

	// Check if config file actually exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &LegacySetupResult{
			Success:   false,
			StartTime: time.Now(),
			EndTime:   time.Now(),
			Message:   "Arquivo de configuração não encontrado",
		}, nil
	}

	// Convert status to legacy result
	return &LegacySetupResult{
		Success:     true,
		StartTime:   time.Now(),
		EndTime:     time.Now(),
		ConfigPath:  configPath,
		Environment: fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		Message:     "Syntropy Manager está configurado corretamente",
	}, nil
}

// ResetLegacy resets the Syntropy CLI configuration
func ResetLegacy(options LegacySetupOptions) (*LegacySetupResult, error) {
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
		return &LegacySetupResult{
			Success:   false,
			StartTime: time.Now(),
			EndTime:   time.Now(),
			Error:     err,
			Message:   err.Error(),
		}, err
	}

	// Return success result
	return &LegacySetupResult{
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

// Helper functions for environment detection

// getCurrentEnvironmentInfo gets current environment information
func getCurrentEnvironmentInfo() *types.EnvironmentInfo {
	homeDir, _ := os.UserHomeDir()
	return &types.EnvironmentInfo{
		OS:              runtime.GOOS,
		OSVersion:       "unknown", // Would be populated by actual detection
		Architecture:    runtime.GOARCH,
		HomeDir:         homeDir,
		HasAdminRights:  true,  // Would be detected
		AvailableDiskGB: 100.0, // Would be calculated
		HasInternet:     true,  // Would be tested
		CanProceed:      true,
		Issues:          []string{},
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
func convertStatusToLegacySetupResult(status map[string]interface{}) *LegacySetupResult {
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

	return &LegacySetupResult{
		Success:     success,
		StartTime:   time.Now(),
		EndTime:     time.Now(),
		ConfigPath:  configPath,
		Environment: environment,
		Message:     fmt.Sprintf("Status: %s", statusStr),
	}
}

// copyDirectory copia um diretório recursivamente
func copyDirectory(src, dst string) error {
	// Criar diretório de destino
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	// Ler diretório fonte
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// Copiar cada entrada
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Recursivamente copiar subdiretório
			if err := copyDirectory(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copiar arquivo
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile copia um arquivo
func copyFile(src, dst string) error {
	// Abrir arquivo fonte
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Criar arquivo de destino
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copiar conteúdo
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	// Sincronizar arquivo de destino
	return dstFile.Sync()
}

// backupAllDirectories faz backup de todas as pastas exceto a pasta backups
func backupAllDirectories(syntropyDir, backupPath string) error {
	// Criar diretório de backup
	if err := os.MkdirAll(backupPath, 0755); err != nil {
		return err
	}

	// Ler todas as entradas do diretório .syntropy
	entries, err := os.ReadDir(syntropyDir)
	if err != nil {
		return err
	}

	// Copiar cada diretório exceto 'backups'
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() != "backups" {
			srcPath := filepath.Join(syntropyDir, entry.Name())
			dstPath := filepath.Join(backupPath, entry.Name())

			if err := copyDirectory(srcPath, dstPath); err != nil {
				return fmt.Errorf("falha ao copiar diretório %s: %w", entry.Name(), err)
			}
		}
	}

	return nil
}

// createFullBackup cria backup completo de toda a estrutura .syntropy
func (sm *SetupManager) createFullBackup(name string) error {
	homeDir, _ := os.UserHomeDir()
	syntropyDir := filepath.Join(homeDir, ".syntropy")
	backupsDir := filepath.Join(syntropyDir, "backups")

	timestamp := time.Now().Format("20060102_150405")
	backupName := fmt.Sprintf("%s_%s", name, timestamp)
	backupPath := filepath.Join(backupsDir, backupName)

	// Criar diretório de backup
	if err := os.MkdirAll(backupPath, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Copiar todas as pastas exceto 'backups'
	dirsToBackup := []string{"config", "keys", "tokens", "nodes", "logs", "state"}
	for _, dir := range dirsToBackup {
		srcPath := filepath.Join(syntropyDir, dir)
		if _, err := os.Stat(srcPath); err == nil {
			dstPath := filepath.Join(backupPath, dir)
			if err := copyDirectory(srcPath, dstPath); err != nil {
				sm.logger.LogWarning("Failed to backup directory", map[string]interface{}{
					"directory": dir,
					"error":     err.Error(),
				})
			}
		}
	}

	sm.logger.LogInfo("Backup completo criado", map[string]interface{}{
		"backup_path": backupPath,
		"backup_name": backupName,
	})

	return nil
}

// Public methods for CLI interface

// GetStatus returns the current setup status
func (sm *SetupManager) GetStatus() (*types.SetupStatus, error) {
	state, err := sm.stateManager.LoadState()
	if err != nil {
		return nil, err
	}
	return &state.Status, nil
}

// ValidateEnvironmentWithOptions validates environment with options
func (sm *SetupManager) ValidateEnvironmentWithOptions(options *SetupOptions) (*types.EnvironmentInfo, error) {
	// Convert SetupOptions to types.SetupOptions
	typesOptions := &types.SetupOptions{
		Force:          options.Force,
		ValidateOnly:   options.ValidateOnly,
		TestMode:       options.TestMode,
		Verbose:        options.Verbose,
		Quiet:          options.Quiet,
		CustomSettings: options.CustomSettings,
	}
	return sm.validator.ValidateEnvironmentWithOptions(typesOptions)
}

// GetOwnerKeyInfo returns information about the Owner Keys
func (sm *SetupManager) GetOwnerKeyInfo() (*types.OwnerKeyInfo, error) {
	homeDir, _ := os.UserHomeDir()
	keysDir := filepath.Join(homeDir, ".syntropy", "keys")

	privateKeyPath := filepath.Join(keysDir, "owner.key")
	publicKeyPath := filepath.Join(keysDir, "owner.key.pub")

	// Verificar se as chaves existem
	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("owner key not found")
	}

	// Ler chave pública
	publicKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key: %w", err)
	}

	// Calcular fingerprint
	fingerprint := sm.calculateFingerprint(publicKeyData)

	// Obter informações do arquivo
	privateKeyInfo, _ := os.Stat(privateKeyPath)

	return &types.OwnerKeyInfo{
		Algorithm:   "ed25519",
		Fingerprint: fingerprint,
		PublicKey:   string(publicKeyData),
		CreatedAt:   privateKeyInfo.ModTime(),
		Path:        privateKeyPath,
	}, nil
}

// GetOwnerPublicKey returns the public key
func (sm *SetupManager) GetOwnerPublicKey() (string, error) {
	homeDir, _ := os.UserHomeDir()
	publicKeyPath := filepath.Join(homeDir, ".syntropy", "keys", "owner.key.pub")

	publicKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return "", fmt.Errorf("failed to read public key: %w", err)
	}

	return string(publicKeyData), nil
}

// GetOwnerPrivateKey returns the private key (SENSITIVE USE)
func (sm *SetupManager) GetOwnerPrivateKey(passphrase string) (string, error) {
	homeDir, _ := os.UserHomeDir()
	privateKeyPath := filepath.Join(homeDir, ".syntropy", "keys", "owner.key")

	encryptedData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return "", fmt.Errorf("failed to read private key: %w", err)
	}

	// For now, return the raw data (in real implementation, decrypt here)
	return string(encryptedData), nil
}

// ExportOwnerKeys exports Owner Keys for backup
func (sm *SetupManager) ExportOwnerKeys(outputPath string, includePrivate bool, passphrase string) error {
	homeDir, _ := os.UserHomeDir()
	keysDir := filepath.Join(homeDir, ".syntropy", "keys")

	publicKeyPath := filepath.Join(keysDir, "owner.key.pub")
	publicKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read public key: %w", err)
	}

	export := &types.OwnerKeyExport{
		PublicKey:  string(publicKeyData),
		ExportedAt: time.Now(),
		Version:    "1.0.0",
	}

	if includePrivate {
		privateKeyPath := filepath.Join(keysDir, "owner.key")
		privateKeyData, err := os.ReadFile(privateKeyPath)
		if err != nil {
			return fmt.Errorf("failed to read private key: %w", err)
		}

		// For now, store as is (in real implementation, encrypt with passphrase)
		export.PrivateKeyEncrypted = string(privateKeyData)
	}

	// Serialize export
	data, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal export: %w", err)
	}

	// Save file
	if err := os.WriteFile(outputPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write export file: %w", err)
	}

	return nil
}

// ImportOwnerKeys imports Owner Keys from backup
func (sm *SetupManager) ImportOwnerKeys(inputPath string, passphrase string) error {
	// Read export file
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read import file: %w", err)
	}

	// Deserialize export
	var export types.OwnerKeyExport
	if err := json.Unmarshal(data, &export); err != nil {
		return fmt.Errorf("failed to unmarshal export: %w", err)
	}

	homeDir, _ := os.UserHomeDir()
	keysDir := filepath.Join(homeDir, ".syntropy", "keys")

	// Save public key
	publicKeyPath := filepath.Join(keysDir, "owner.key.pub")
	if err := os.WriteFile(publicKeyPath, []byte(export.PublicKey), 0644); err != nil {
		return fmt.Errorf("failed to write public key: %w", err)
	}

	// If includes private key, import it
	if export.PrivateKeyEncrypted != "" {
		privateKeyPath := filepath.Join(keysDir, "owner.key")
		if err := os.WriteFile(privateKeyPath, []byte(export.PrivateKeyEncrypted), 0600); err != nil {
			return fmt.Errorf("failed to write private key: %w", err)
		}
	}

	return nil
}

// calculateFingerprint calculates key fingerprint
func (sm *SetupManager) calculateFingerprint(keyData []byte) string {
	hash := sha256.Sum256(keyData)
	return hex.EncodeToString(hash[:])
}
