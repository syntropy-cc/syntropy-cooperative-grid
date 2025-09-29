package integration

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src/internal/types"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/tests/helpers"
)

// TestConfigIntegration_ConfigurationGeneration testa a geração de configuração
func TestConfigIntegration_ConfigurationGeneration(t *testing.T) {
	tests := []struct {
		name    string
		options *types.ConfigOptions
		wantErr bool
	}{
		{
			name: "should generate configuration successfully",
			options: &types.ConfigOptions{
				OwnerName:      "Test User",
				OwnerEmail:     "test@example.com",
				NetworkConfig:  nil,
				SecurityConfig: nil,
				CustomSettings: map[string]string{
					"log_level":    "info",
					"api_endpoint": "https://api.syntropy.network",
					"debug_mode":   "false",
				},
			},
			wantErr: false,
		},
		{
			name: "should handle minimal configuration",
			options: &types.ConfigOptions{
				OwnerName:      "Test User",
				OwnerEmail:     "test@example.com",
				NetworkConfig:  nil,
				SecurityConfig: nil,
				CustomSettings: map[string]string{},
			},
			wantErr: false,
		},
		{
			name: "should handle empty owner name",
			options: &types.ConfigOptions{
				OwnerName:      "",
				OwnerEmail:     "test@example.com",
				NetworkConfig:  nil,
				SecurityConfig: nil,
				CustomSettings: map[string]string{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "config_generation_test")
			configDir := filepath.Join(tempDir, "config")
			os.MkdirAll(configDir, 0755)

			// Criar configurator
			configurator := src.NewConfigurator(nil)
			if configurator == nil {
				t.Fatal("NewConfigurator() returned nil")
			}

			// Executar geração de configuração
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
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestConfigIntegration_StructureCreation testa a criação de estrutura
func TestConfigIntegration_StructureCreation(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should create structure successfully",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "config_structure_test")
			configDir := filepath.Join(tempDir, "config")
			os.MkdirAll(configDir, 0755)

			// Criar configurator
			configurator := src.NewConfigurator(nil)
			if configurator == nil {
				t.Fatal("NewConfigurator() returned nil")
			}

			// Executar criação de estrutura
			err := configurator.CreateStructure()
			if (err != nil) != tt.wantErr {
				t.Errorf("Configurator.CreateStructure() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verificar se os diretórios foram criados
				helpers.AssertDirExists(t, configDir, "Config directory")
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestConfigIntegration_KeyGeneration testa a geração de chaves
func TestConfigIntegration_KeyGeneration(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should generate keys successfully",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "config_keys_test")
			keysDir := filepath.Join(tempDir, "keys")
			os.MkdirAll(keysDir, 0755)

			// Criar configurator
			configurator := src.NewConfigurator(nil)
			if configurator == nil {
				t.Fatal("NewConfigurator() returned nil")
			}

			// Executar geração de chaves
			keyPair, err := configurator.GenerateKeys()
			if (err != nil) != tt.wantErr {
				t.Errorf("Configurator.GenerateKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if keyPair == nil {
					t.Error("GenerateKeys() returned nil keyPair")
					return
				}

				// Verificar se a chave foi gerada corretamente
				helpers.AssertStringNotEmpty(t, keyPair.ID, "KeyPair ID")
				helpers.AssertStringEqual(t, keyPair.Algorithm, "ed25519", "KeyPair Algorithm")
				helpers.AssertStringNotEmpty(t, keyPair.Fingerprint, "KeyPair Fingerprint")
				helpers.AssertStringNotEmpty(t, keyPair.PublicKey, "KeyPair PublicKey")
				helpers.AssertStringNotEmpty(t, keyPair.PrivateKey, "KeyPair PrivateKey")
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestConfigIntegration_ConfigurationValidation testa a validação de configuração
func TestConfigIntegration_ConfigurationValidation(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(string) error
		wantErr   bool
	}{
		{
			name: "should validate valid configuration",
			setupFunc: func(configDir string) error {
				// Criar arquivo de configuração válido
				configPath := filepath.Join(configDir, "manager.yaml")
				configContent := `
manager:
  home_dir: "/home/testuser/.syntropy"
  log_level: "info"
  api_endpoint: "https://api.syntropy.network"
  directories:
    config: "/home/testuser/.syntropy/config"
    keys: "/home/testuser/.syntropy/keys"
  default_paths:
    config: "/home/testuser/.syntropy/config/manager.yaml"
    log: "/home/testuser/.syntropy/logs/manager.log"
owner_key:
  type: "ed25519"
  path: "/home/testuser/.syntropy/keys/owner.key"
environment:
  os: "linux"
  architecture: "amd64"
  home_dir: "/home/testuser"
`
				return os.WriteFile(configPath, []byte(configContent), 0644)
			},
			wantErr: false,
		},
		{
			name: "should fail validation for invalid configuration",
			setupFunc: func(configDir string) error {
				// Criar arquivo de configuração inválido
				configPath := filepath.Join(configDir, "manager.yaml")
				configContent := "invalid yaml content: ["
				return os.WriteFile(configPath, []byte(configContent), 0644)
			},
			wantErr: true,
		},
		{
			name: "should fail validation for missing configuration",
			setupFunc: func(configDir string) error {
				// Não criar arquivo de configuração
				return nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "config_validation_test")
			configDir := filepath.Join(tempDir, "config")
			os.MkdirAll(configDir, 0755)

			// Setup do teste
			if tt.setupFunc != nil {
				if err := tt.setupFunc(configDir); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			// Criar configurator
			configurator := src.NewConfigurator(nil)
			if configurator == nil {
				t.Fatal("NewConfigurator() returned nil")
			}

			// Executar validação de configuração
			err := configurator.ValidateConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("Configurator.ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestConfigIntegration_ConfigurationBackup testa o backup de configuração
func TestConfigIntegration_ConfigurationBackup(t *testing.T) {
	tests := []struct {
		name       string
		backupName string
		setupFunc  func(string) error
		wantErr    bool
	}{
		{
			name:       "should backup configuration successfully",
			backupName: "config_backup",
			setupFunc: func(configDir string) error {
				// Criar arquivo de configuração
				configPath := filepath.Join(configDir, "manager.yaml")
				configContent := `
manager:
  home_dir: "/home/testuser/.syntropy"
  log_level: "info"
owner_key:
  type: "ed25519"
environment:
  os: "linux"
`
				return os.WriteFile(configPath, []byte(configContent), 0644)
			},
			wantErr: false,
		},
		{
			name:       "should handle empty backup name",
			backupName: "",
			setupFunc: func(configDir string) error {
				// Criar arquivo de configuração
				configPath := filepath.Join(configDir, "manager.yaml")
				configContent := `
manager:
  home_dir: "/home/testuser/.syntropy"
`
				return os.WriteFile(configPath, []byte(configContent), 0644)
			},
			wantErr: false,
		},
		{
			name:       "should fail backup for missing configuration",
			backupName: "config_backup",
			setupFunc: func(configDir string) error {
				// Não criar arquivo de configuração
				return nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "config_backup_test")
			configDir := filepath.Join(tempDir, "config")
			backupDir := filepath.Join(tempDir, "backups")
			os.MkdirAll(configDir, 0755)
			os.MkdirAll(backupDir, 0755)

			// Setup do teste
			if tt.setupFunc != nil {
				if err := tt.setupFunc(configDir); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			// Criar configurator
			configurator := src.NewConfigurator(nil)
			if configurator == nil {
				t.Fatal("NewConfigurator() returned nil")
			}

			// Executar backup de configuração
			err := configurator.BackupConfig(tt.backupName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Configurator.BackupConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verificar se o backup foi criado
				backupFile := filepath.Join(backupDir, tt.backupName+".yaml")
				helpers.AssertFileExists(t, backupFile, "Config backup file")
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestConfigIntegration_ConfigurationRestore testa a restauração de configuração
func TestConfigIntegration_ConfigurationRestore(t *testing.T) {
	tests := []struct {
		name       string
		backupPath string
		setupFunc  func(string) error
		wantErr    bool
	}{
		{
			name:       "should restore configuration successfully",
			backupPath: "config_backup.yaml",
			setupFunc: func(backupDir string) error {
				// Criar arquivo de backup
				backupFile := filepath.Join(backupDir, "config_backup.yaml")
				backupContent := `
manager:
  home_dir: "/home/testuser/.syntropy"
  log_level: "info"
owner_key:
  type: "ed25519"
environment:
  os: "linux"
`
				return os.WriteFile(backupFile, []byte(backupContent), 0644)
			},
			wantErr: false,
		},
		{
			name:       "should fail restore for missing backup",
			backupPath: "nonexistent_backup.yaml",
			setupFunc: func(backupDir string) error {
				// Não criar arquivo de backup
				return nil
			},
			wantErr: true,
		},
		{
			name:       "should fail restore for invalid backup",
			backupPath: "invalid_backup.yaml",
			setupFunc: func(backupDir string) error {
				// Criar arquivo de backup inválido
				backupFile := filepath.Join(backupDir, "invalid_backup.yaml")
				backupContent := "invalid yaml content: ["
				return os.WriteFile(backupFile, []byte(backupContent), 0644)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "config_restore_test")
			configDir := filepath.Join(tempDir, "config")
			backupDir := filepath.Join(tempDir, "backups")
			os.MkdirAll(configDir, 0755)
			os.MkdirAll(backupDir, 0755)

			// Setup do teste
			if tt.setupFunc != nil {
				if err := tt.setupFunc(backupDir); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			// Criar configurator
			configurator := src.NewConfigurator(nil)
			if configurator == nil {
				t.Fatal("NewConfigurator() returned nil")
			}

			// Executar restauração de configuração
			fullBackupPath := filepath.Join(backupDir, tt.backupPath)
			err := configurator.RestoreConfig(fullBackupPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Configurator.RestoreConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verificar se a configuração foi restaurada
				configFile := filepath.Join(configDir, "manager.yaml")
				helpers.AssertFileExists(t, configFile, "Restored config file")
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestConfigIntegration_ErrorHandling testa o tratamento de erros
func TestConfigIntegration_ErrorHandling(t *testing.T) {
	tests := []struct {
		name    string
		options *types.ConfigOptions
		wantErr bool
	}{
		{
			name:    "should handle nil options",
			options: nil,
			wantErr: true,
		},
		{
			name: "should handle invalid email",
			options: &types.ConfigOptions{
				OwnerName:      "Test User",
				OwnerEmail:     "invalid-email",
				NetworkConfig:  nil,
				SecurityConfig: nil,
				CustomSettings: map[string]string{},
			},
			wantErr: false, // Should not fail, just use invalid email
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "config_error_test")
			configDir := filepath.Join(tempDir, "config")
			os.MkdirAll(configDir, 0755)

			// Criar configurator
			configurator := src.NewConfigurator(nil)
			if configurator == nil {
				t.Fatal("NewConfigurator() returned nil")
			}

			// Executar geração de configuração
			err := configurator.GenerateConfig(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("Configurator.GenerateConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestConfigIntegration_Concurrency testa concorrência
func TestConfigIntegration_Concurrency(t *testing.T) {
	t.Run("should handle concurrent configuration operations", func(t *testing.T) {
		// Criar diretório temporário para testes
		tempDir := helpers.CreateTempDir(t, "config_concurrent_test")
		configDir := filepath.Join(tempDir, "config")
		os.MkdirAll(configDir, 0755)

		// Criar configurator
		configurator := src.NewConfigurator(nil)
		if configurator == nil {
			t.Fatal("NewConfigurator() returned nil")
		}

		// Executar múltiplas operações de configuração concorrentemente
		done := make(chan bool, 5)
		for i := 0; i < 5; i++ {
			go func(instance int) {
				options := &types.ConfigOptions{
					OwnerName:      "Test User " + string(rune(instance)),
					OwnerEmail:     "test" + string(rune(instance)) + "@example.com",
					NetworkConfig:  nil,
					SecurityConfig: nil,
					CustomSettings: map[string]string{
						"instance": string(rune(instance)),
					},
				}

				err := configurator.GenerateConfig(options)
				if err != nil {
					t.Errorf("Concurrent GenerateConfig() failed: %v", err)
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

// TestConfigIntegration_Performance testa performance
func TestConfigIntegration_Performance(t *testing.T) {
	t.Run("should complete configuration operations within reasonable time", func(t *testing.T) {
		// Criar diretório temporário para testes
		tempDir := helpers.CreateTempDir(t, "config_perf_test")
		configDir := filepath.Join(tempDir, "config")
		os.MkdirAll(configDir, 0755)

		// Criar configurator
		configurator := src.NewConfigurator(nil)
		if configurator == nil {
			t.Fatal("NewConfigurator() returned nil")
		}

		options := &types.ConfigOptions{
			OwnerName:      "Test User",
			OwnerEmail:     "test@example.com",
			NetworkConfig:  nil,
			SecurityConfig: nil,
			CustomSettings: map[string]string{
				"log_level":    "info",
				"api_endpoint": "https://api.syntropy.network",
				"debug_mode":   "false",
			},
		}

		start := time.Now()
		err := configurator.GenerateConfig(options)
		elapsed := time.Since(start)

		if err != nil {
			t.Errorf("GenerateConfig() failed: %v", err)
		}

		if elapsed > 5*time.Second {
			t.Errorf("GenerateConfig() took too long: %v", elapsed)
		}

		// Limpar diretório temporário
		os.RemoveAll(tempDir)
	})
}
