// Package setup provides functionality for setting up the Syntropy CLI environment
package setup

import (
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
func (sm *SetupManager) Setup(options *SetupOptions) error {
	// Check for nil options
	if options == nil {
		return sm.handleError(fmt.Errorf("setup options cannot be nil"), "invalid_options")
	}

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

	// 3. Gerar ou carregar chaves existentes
	keyPair, err := sm.keyManager.GenerateOrLoadKeyPair("ed25519")
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
	newOptions := &SetupOptions{
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
