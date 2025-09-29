package e2e

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src/internal/types"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/tests/helpers"
)

// TestSetupE2E_CompleteSetupFlow testa o fluxo completo do setup end-to-end
func TestSetupE2E_CompleteSetupFlow(t *testing.T) {
	tests := []struct {
		name    string
		options *types.SetupOptions
		wantErr bool
	}{
		{
			name: "should complete full setup flow successfully",
			options: &types.SetupOptions{
				OwnerName:      "Test User",
				OwnerEmail:     "test@example.com",
				ConfigPath:     "/tmp/syntropy_e2e_test/config",
				KeysPath:       "/tmp/syntropy_e2e_test/keys",
				BackupPath:     "/tmp/syntropy_e2e_test/backups",
				LogPath:        "/tmp/syntropy_e2e_test/logs",
				Verbose:        true,
				Force:          false,
				SkipValidation: false,
			},
			wantErr: false,
		},
		{
			name: "should complete setup with force flag",
			options: &types.SetupOptions{
				OwnerName:      "Test User",
				OwnerEmail:     "test@example.com",
				ConfigPath:     "/tmp/syntropy_e2e_test/config",
				KeysPath:       "/tmp/syntropy_e2e_test/keys",
				BackupPath:     "/tmp/syntropy_e2e_test/backups",
				LogPath:        "/tmp/syntropy_e2e_test/logs",
				Verbose:        false,
				Force:          true,
				SkipValidation: false,
			},
			wantErr: false,
		},
		{
			name: "should complete setup with skip validation",
			options: &types.SetupOptions{
				OwnerName:      "Test User",
				OwnerEmail:     "test@example.com",
				ConfigPath:     "/tmp/syntropy_e2e_test/config",
				KeysPath:       "/tmp/syntropy_e2e_test/keys",
				BackupPath:     "/tmp/syntropy_e2e_test/backups",
				LogPath:        "/tmp/syntropy_e2e_test/logs",
				Verbose:        false,
				Force:          false,
				SkipValidation: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "setup_e2e_test")

			// Ajustar caminhos para usar o diretório temporário
			tt.options.ConfigPath = filepath.Join(tempDir, "config")
			tt.options.KeysPath = filepath.Join(tempDir, "keys")
			tt.options.BackupPath = filepath.Join(tempDir, "backups")
			tt.options.LogPath = filepath.Join(tempDir, "logs")

			// Criar setup manager
			setupManager := src.NewSetupManager()
			if setupManager == nil {
				t.Fatal("NewSetupManager() returned nil")
			}

			// Executar setup completo
			err := setupManager.Setup(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupManager.Setup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verificar se todos os diretórios foram criados
				helpers.AssertDirExists(t, tt.options.ConfigPath, "Config directory")
				helpers.AssertDirExists(t, tt.options.KeysPath, "Keys directory")
				helpers.AssertDirExists(t, tt.options.BackupPath, "Backup directory")
				helpers.AssertDirExists(t, tt.options.LogPath, "Log directory")

				// Verificar se os arquivos de configuração foram criados
				configFile := filepath.Join(tt.options.ConfigPath, "manager.yaml")
				helpers.AssertFileExists(t, configFile, "Manager config file")

				// Verificar se as chaves foram geradas
				keys, err := setupManager.ListKeys()
				if err != nil {
					t.Errorf("ListKeys() failed: %v", err)
				}
				if len(keys) == 0 {
					t.Error("No keys were generated")
				}

				// Verificar status do setup
				status, err := setupManager.Status()
				if err != nil {
					t.Errorf("Status() failed: %v", err)
				}
				if status == nil {
					t.Error("Status() returned nil")
				}

				// Verificar se o setup está funcionando corretamente
				if status.Status != "completed" {
					t.Errorf("Setup status is %s, expected completed", status.Status)
				}
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestSetupE2E_ValidationFlow testa o fluxo de validação end-to-end
func TestSetupE2E_ValidationFlow(t *testing.T) {
	tests := []struct {
		name    string
		options *types.SetupOptions
		wantErr bool
	}{
		{
			name: "should validate environment end-to-end",
			options: &types.SetupOptions{
				OwnerName:      "Test User",
				OwnerEmail:     "test@example.com",
				ConfigPath:     "/tmp/syntropy_e2e_test/config",
				KeysPath:       "/tmp/syntropy_e2e_test/keys",
				BackupPath:     "/tmp/syntropy_e2e_test/backups",
				LogPath:        "/tmp/syntropy_e2e_test/logs",
				Verbose:        true,
				Force:          false,
				SkipValidation: false,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "validation_e2e_test")

			// Ajustar caminhos para usar o diretório temporário
			tt.options.ConfigPath = filepath.Join(tempDir, "config")
			tt.options.KeysPath = filepath.Join(tempDir, "keys")
			tt.options.BackupPath = filepath.Join(tempDir, "backups")
			tt.options.LogPath = filepath.Join(tempDir, "logs")

			// Criar setup manager
			setupManager := src.NewSetupManager()
			if setupManager == nil {
				t.Fatal("NewSetupManager() returned nil")
			}

			// Executar validação completa
			result, err := setupManager.Validate(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupManager.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if result == nil {
					t.Error("Validate() returned nil result")
					return
				}

				// Verificar se a validação foi bem-sucedida
				if !result.Success {
					t.Errorf("Validation failed: %v", result.Errors)
				}

				// Verificar se não há erros críticos
				if len(result.Errors) > 0 {
					t.Errorf("Validation returned errors: %v", result.Errors)
				}

				// Verificar se a validação incluiu todos os aspectos
				if result.Environment == nil {
					t.Error("Validation result missing environment info")
				}
				if result.Dependencies == nil {
					t.Error("Validation result missing dependencies info")
				}
				if result.Network == nil {
					t.Error("Validation result missing network info")
				}
				if result.Permissions == nil {
					t.Error("Validation result missing permissions info")
				}
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestSetupE2E_ConfigurationFlow testa o fluxo de configuração end-to-end
func TestSetupE2E_ConfigurationFlow(t *testing.T) {
	tests := []struct {
		name    string
		options *types.ConfigOptions
		wantErr bool
	}{
		{
			name: "should generate complete configuration end-to-end",
			options: &types.ConfigOptions{
				OwnerName:      "Test User",
				OwnerEmail:     "test@example.com",
				NetworkConfig:  nil,
				SecurityConfig: nil,
				CustomSettings: map[string]string{
					"log_level":       "info",
					"api_endpoint":    "https://api.syntropy.network",
					"debug_mode":      "false",
					"max_connections": "100",
					"timeout":         "30s",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "config_e2e_test")
			configDir := filepath.Join(tempDir, "config")
			os.MkdirAll(configDir, 0755)

			// Criar configurator
			configurator := src.NewConfigurator(nil)
			if configurator == nil {
				t.Fatal("NewConfigurator() returned nil")
			}

			// Executar geração de configuração completa
			err := configurator.GenerateConfig(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("Configurator.GenerateConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verificar se o arquivo de configuração foi criado
				configFile := filepath.Join(configDir, "manager.yaml")
				helpers.AssertFileExists(t, configFile, "Manager config file")

				// Verificar se a configuração é válida
				err = configurator.ValidateConfig()
				if err != nil {
					t.Errorf("ValidateConfig() failed: %v", err)
				}

				// Verificar se a estrutura de diretórios foi criada
				err = configurator.CreateStructure()
				if err != nil {
					t.Errorf("CreateStructure() failed: %v", err)
				}

				// Verificar se as chaves foram geradas
				keyPair, err := configurator.GenerateKeys()
				if err != nil {
					t.Errorf("GenerateKeys() failed: %v", err)
				}
				if keyPair == nil {
					t.Error("GenerateKeys() returned nil keyPair")
				}
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestSetupE2E_KeyManagementFlow testa o fluxo de gerenciamento de chaves end-to-end
func TestSetupE2E_KeyManagementFlow(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should manage keys end-to-end successfully",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "keys_e2e_test")
			keysDir := filepath.Join(tempDir, "keys")
			os.MkdirAll(keysDir, 0755)

			// Criar key manager
			keyManager := src.NewKeyManager(nil)
			if keyManager == nil {
				t.Fatal("NewKeyManager() returned nil")
			}

			// Executar fluxo completo de gerenciamento de chaves
			// 1. Gerar par de chaves
			keyPair, err := keyManager.GenerateKeyPair()
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyManager.GenerateKeyPair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if keyPair == nil {
					t.Error("GenerateKeyPair() returned nil keyPair")
					return
				}

				// 2. Armazenar chave
				err = keyManager.StoreKeyPair(keyPair)
				if err != nil {
					t.Errorf("StoreKeyPair() failed: %v", err)
				}

				// 3. Carregar chave
				loadedKey, err := keyManager.LoadKeyPair(keyPair.ID)
				if err != nil {
					t.Errorf("LoadKeyPair() failed: %v", err)
				}
				if loadedKey == nil {
					t.Error("LoadKeyPair() returned nil keyPair")
					return
				}

				// 4. Verificar se as chaves são iguais
				if keyPair.ID != loadedKey.ID {
					t.Error("Loaded key ID does not match original")
				}
				if keyPair.PublicKey != loadedKey.PublicKey {
					t.Error("Loaded public key does not match original")
				}
				if keyPair.PrivateKey != loadedKey.PrivateKey {
					t.Error("Loaded private key does not match original")
				}

				// 5. Verificar integridade da chave
				err = keyManager.VerifyKeyIntegrity(keyPair.ID)
				if err != nil {
					t.Errorf("VerifyKeyIntegrity() failed: %v", err)
				}

				// 6. Listar chaves
				keys, err := keyManager.ListKeys()
				if err != nil {
					t.Errorf("ListKeys() failed: %v", err)
				}
				if len(keys) == 0 {
					t.Error("ListKeys() returned empty list")
				}

				// 7. Rotacionar chaves
				newKey, err := keyManager.RotateKeys()
				if err != nil {
					t.Errorf("RotateKeys() failed: %v", err)
				}
				if newKey == nil {
					t.Error("RotateKeys() returned nil keyPair")
				}

				// 8. Verificar se a nova chave é diferente da inicial
				if keyPair.ID == newKey.ID {
					t.Error("Key rotation generated same key ID")
				}
				if keyPair.Fingerprint == newKey.Fingerprint {
					t.Error("Key rotation generated same key fingerprint")
				}
				if keyPair.PublicKey == newKey.PublicKey {
					t.Error("Key rotation generated same public key")
				}
				if keyPair.PrivateKey == newKey.PrivateKey {
					t.Error("Key rotation generated same private key")
				}
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestSetupE2E_StateManagementFlow testa o fluxo de gerenciamento de estado end-to-end
func TestSetupE2E_StateManagementFlow(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should manage state end-to-end successfully",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "state_e2e_test")
			stateDir := filepath.Join(tempDir, "state")
			os.MkdirAll(stateDir, 0755)

			// Criar state manager
			stateManager := src.NewStateManager(nil)
			if stateManager == nil {
				t.Fatal("NewStateManager() returned nil")
			}

			// Executar fluxo completo de gerenciamento de estado
			// 1. Criar estado inicial
			initialState := &types.SetupState{
				CurrentStep:    "initialization",
				CompletedSteps: []string{},
				Errors:         []string{},
				Warnings:       []string{},
				StartTime:      time.Now(),
				LastUpdated:    time.Now(),
				ConfigPath:     "/tmp/syntropy_e2e_test/config",
				KeysPath:       "/tmp/syntropy_e2e_test/keys",
				BackupPath:     "/tmp/syntropy_e2e_test/backups",
			}

			// 2. Salvar estado
			err := stateManager.SaveState(initialState)
			if (err != nil) != tt.wantErr {
				t.Errorf("StateManager.SaveState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// 3. Carregar estado
				loadedState, err := stateManager.LoadState()
				if err != nil {
					t.Errorf("LoadState() failed: %v", err)
				}
				if loadedState == nil {
					t.Error("LoadState() returned nil state")
					return
				}

				// 4. Verificar se os estados são iguais
				if initialState.CurrentStep != loadedState.CurrentStep {
					t.Error("Loaded state current step does not match original")
				}
				if len(initialState.CompletedSteps) != len(loadedState.CompletedSteps) {
					t.Error("Loaded state completed steps length does not match original")
				}

				// 5. Atualizar estado
				updates := map[string]interface{}{
					"current_step":    "validation",
					"completed_steps": []string{"initialization"},
				}
				err = stateManager.UpdateState(updates)
				if err != nil {
					t.Errorf("UpdateState() failed: %v", err)
				}

				// 6. Verificar integridade do estado
				err = stateManager.VerifyIntegrity()
				if err != nil {
					t.Errorf("VerifyIntegrity() failed: %v", err)
				}

				// 7. Fazer backup do estado
				backupPath := filepath.Join(tempDir, "state_backup.json")
				err = stateManager.BackupState(backupPath)
				if err != nil {
					t.Errorf("BackupState() failed: %v", err)
				}

				// 8. Verificar se o backup foi criado
				helpers.AssertFileExists(t, backupPath, "State backup file")

				// 9. Restaurar estado do backup
				err = stateManager.RestoreState(backupPath)
				if err != nil {
					t.Errorf("RestoreState() failed: %v", err)
				}

				// 10. Verificar se o estado foi restaurado corretamente
				restoredState, err := stateManager.LoadState()
				if err != nil {
					t.Errorf("LoadState() after restore failed: %v", err)
				}
				if restoredState == nil {
					t.Error("LoadState() after restore returned nil state")
				}
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestSetupE2E_ErrorHandling testa o tratamento de erros end-to-end
func TestSetupE2E_ErrorHandling(t *testing.T) {
	tests := []struct {
		name    string
		options *types.SetupOptions
		wantErr bool
	}{
		{
			name: "should handle invalid owner email end-to-end",
			options: &types.SetupOptions{
				OwnerName:      "Test User",
				OwnerEmail:     "invalid-email",
				ConfigPath:     "/tmp/syntropy_e2e_test/config",
				KeysPath:       "/tmp/syntropy_e2e_test/keys",
				BackupPath:     "/tmp/syntropy_e2e_test/backups",
				LogPath:        "/tmp/syntropy_e2e_test/logs",
				Verbose:        false,
				Force:          false,
				SkipValidation: false,
			},
			wantErr: true,
		},
		{
			name: "should handle empty owner name end-to-end",
			options: &types.SetupOptions{
				OwnerName:      "",
				OwnerEmail:     "test@example.com",
				ConfigPath:     "/tmp/syntropy_e2e_test/config",
				KeysPath:       "/tmp/syntropy_e2e_test/keys",
				BackupPath:     "/tmp/syntropy_e2e_test/backups",
				LogPath:        "/tmp/syntropy_e2e_test/logs",
				Verbose:        false,
				Force:          false,
				SkipValidation: false,
			},
			wantErr: false, // Should not fail, just use empty name
		},
		{
			name: "should handle invalid paths end-to-end",
			options: &types.SetupOptions{
				OwnerName:      "Test User",
				OwnerEmail:     "test@example.com",
				ConfigPath:     "/invalid/path/config",
				KeysPath:       "/invalid/path/keys",
				BackupPath:     "/invalid/path/backups",
				LogPath:        "/invalid/path/logs",
				Verbose:        false,
				Force:          false,
				SkipValidation: false,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "error_e2e_test")

			// Ajustar caminhos para usar o diretório temporário se não forem inválidos
			if tt.options.ConfigPath != "/invalid/path/config" {
				tt.options.ConfigPath = filepath.Join(tempDir, "config")
				tt.options.KeysPath = filepath.Join(tempDir, "keys")
				tt.options.BackupPath = filepath.Join(tempDir, "backups")
				tt.options.LogPath = filepath.Join(tempDir, "logs")
			}

			// Criar setup manager
			setupManager := src.NewSetupManager()
			if setupManager == nil {
				t.Fatal("NewSetupManager() returned nil")
			}

			// Executar setup que pode falhar
			err := setupManager.Setup(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupManager.Setup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestSetupE2E_Concurrency testa concorrência end-to-end
func TestSetupE2E_Concurrency(t *testing.T) {
	t.Run("should handle concurrent setup operations end-to-end", func(t *testing.T) {
		// Criar diretório temporário para testes
		tempDir := helpers.CreateTempDir(t, "concurrency_e2e_test")

		// Criar setup manager
		setupManager := src.NewSetupManager()
		if setupManager == nil {
			t.Fatal("NewSetupManager() returned nil")
		}

		// Executar múltiplas operações de setup concorrentemente
		done := make(chan bool, 5)
		for i := 0; i < 5; i++ {
			go func(instance int) {
				options := &types.SetupOptions{
					OwnerName:      "Test User " + string(rune(instance)),
					OwnerEmail:     "test" + string(rune(instance)) + "@example.com",
					ConfigPath:     filepath.Join(tempDir, "config_"+string(rune(instance))),
					KeysPath:       filepath.Join(tempDir, "keys_"+string(rune(instance))),
					BackupPath:     filepath.Join(tempDir, "backups_"+string(rune(instance))),
					LogPath:        filepath.Join(tempDir, "logs_"+string(rune(instance))),
					Verbose:        false,
					Force:          true,
					SkipValidation: true,
				}

				err := setupManager.Setup(options)
				if err != nil {
					t.Errorf("Concurrent Setup() failed: %v", err)
				}
				done <- true
			}(i)
		}

		// Aguardar todas as goroutines terminarem
		for i := 0; i < 5; i++ {
			<-done
		}

		// Limpar diretório temporário
		os.RemoveAll(tempDir)
	})
}

// TestSetupE2E_Performance testa performance end-to-end
func TestSetupE2E_Performance(t *testing.T) {
	t.Run("should complete setup within reasonable time end-to-end", func(t *testing.T) {
		// Criar diretório temporário para testes
		tempDir := helpers.CreateTempDir(t, "performance_e2e_test")

		// Criar setup manager
		setupManager := src.NewSetupManager()
		if setupManager == nil {
			t.Fatal("NewSetupManager() returned nil")
		}

		options := &types.SetupOptions{
			OwnerName:      "Test User",
			OwnerEmail:     "test@example.com",
			ConfigPath:     filepath.Join(tempDir, "config"),
			KeysPath:       filepath.Join(tempDir, "keys"),
			BackupPath:     filepath.Join(tempDir, "backups"),
			LogPath:        filepath.Join(tempDir, "logs"),
			Verbose:        false,
			Force:          false,
			SkipValidation: true,
		}

		start := time.Now()
		err := setupManager.Setup(options)
		elapsed := time.Since(start)

		if err != nil {
			t.Errorf("Setup() failed: %v", err)
		}

		if elapsed > 60*time.Second {
			t.Errorf("Setup() took too long: %v", elapsed)
		}

		// Limpar diretório temporário
		os.RemoveAll(tempDir)
	})
}
