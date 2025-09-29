package unit

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src/internal/types"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/tests/helpers"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/tests/mocks"
)

// TestNewConfigurator testa a criação de um novo Configurator
func TestNewConfigurator(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "should create configurator successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &mocks.MockSetupLogger{}
			configurator := src.NewConfigurator(logger)
			if configurator == nil {
				t.Error("NewConfigurator() returned nil configurator")
			}
		})
	}
}

// TestConfigurator_GenerateConfig testa o método GenerateConfig
func TestConfigurator_GenerateConfig(t *testing.T) {
	tests := []struct {
		name       string
		options    *types.ConfigOptions
		mockLogger *mocks.MockSetupLogger
		wantErr    bool
	}{
		{
			name:       "should generate config successfully",
			options:    helpers.CreateValidConfigOptions(),
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    false,
		},
		{
			name:       "should handle nil options",
			options:    nil,
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    true,
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
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    false, // Should not fail, just use empty name
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "configurator_test")

			// Mock do logger
			logger := &mocks.MockSetupLogger{}

			// Criar configurator com diretório temporário
			configurator := &src.Configurator{
				Logger: logger,
			}

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

// TestConfigurator_CreateStructure testa o método CreateStructure
func TestConfigurator_CreateStructure(t *testing.T) {
	tests := []struct {
		name       string
		mockLogger *mocks.MockSetupLogger
		wantErr    bool
	}{
		{
			name:       "should create structure successfully",
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "configurator_structure_test")

			// Mock do logger
			logger := &mocks.MockSetupLogger{}

			// Criar configurator com diretório temporário
			configurator := &src.Configurator{
				Logger: logger,
			}

			err := configurator.CreateStructure()
			if (err != nil) != tt.wantErr {
				t.Errorf("Configurator.CreateStructure() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestConfigurator_GenerateKeys testa o método GenerateKeys
func TestConfigurator_GenerateKeys(t *testing.T) {
	tests := []struct {
		name       string
		mockLogger *mocks.MockSetupLogger
		wantErr    bool
	}{
		{
			name:       "should generate keys successfully",
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "configurator_keys_test")

			// Mock do logger
			logger := &mocks.MockSetupLogger{}

			// Criar configurator com diretório temporário
			configurator := &src.Configurator{
				Logger: logger,
			}

			keyPair, err := configurator.GenerateKeys()
			if (err != nil) != tt.wantErr {
				t.Errorf("Configurator.GenerateKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if keyPair == nil {
					t.Error("Configurator.GenerateKeys() returned nil keyPair")
					return
				}
				helpers.AssertStringNotEmpty(t, keyPair.ID, "KeyPair ID")
				helpers.AssertStringEqual(t, keyPair.Algorithm, "ed25519", "KeyPair Algorithm")
				helpers.AssertStringNotEmpty(t, keyPair.Fingerprint, "KeyPair Fingerprint")
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestConfigurator_ValidateConfig testa o método ValidateConfig
func TestConfigurator_ValidateConfig(t *testing.T) {
	tests := []struct {
		name       string
		setupFunc  func(string) error
		mockLogger *mocks.MockSetupLogger
		wantErr    bool
	}{
		{
			name: "should validate config successfully when config exists",
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
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    false,
		},
		{
			name: "should fail when config file does not exist",
			setupFunc: func(configDir string) error {
				// Não criar arquivo de configuração
				return nil
			},
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    true,
		},
		{
			name: "should fail when config file is invalid",
			setupFunc: func(configDir string) error {
				// Criar arquivo de configuração inválido
				configPath := filepath.Join(configDir, "manager.yaml")
				configContent := "invalid yaml content: ["
				return os.WriteFile(configPath, []byte(configContent), 0644)
			},
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "configurator_validate_test")
			configDir := filepath.Join(tempDir, "config")
			os.MkdirAll(configDir, 0755)

			// Setup do teste
			if tt.setupFunc != nil {
				if err := tt.setupFunc(configDir); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			// Mock do logger
			logger := &mocks.MockSetupLogger{}

			// Criar configurator com diretório temporário
			configurator := &src.Configurator{
				Logger: logger,
			}

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

// TestConfigurator_BackupConfig testa o método BackupConfig
func TestConfigurator_BackupConfig(t *testing.T) {
	tests := []struct {
		name       string
		backupName string
		setupFunc  func(string) error
		mockLogger *mocks.MockSetupLogger
		wantErr    bool
	}{
		{
			name:       "should backup config successfully",
			backupName: "test_backup",
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
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    false,
		},
		{
			name:       "should fail when config file does not exist",
			backupName: "test_backup",
			setupFunc: func(configDir string) error {
				// Não criar arquivo de configuração
				return nil
			},
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    true,
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
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "configurator_backup_test")
			configDir := filepath.Join(tempDir, "config")
			os.MkdirAll(configDir, 0755)

			// Setup do teste
			if tt.setupFunc != nil {
				if err := tt.setupFunc(configDir); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			// Mock do logger
			logger := &mocks.MockSetupLogger{}

			// Criar configurator com diretório temporário
			configurator := &src.Configurator{
				Logger: logger,
			}

			err := configurator.BackupConfig(tt.backupName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Configurator.BackupConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestConfigurator_RestoreConfig testa o método RestoreConfig
func TestConfigurator_RestoreConfig(t *testing.T) {
	tests := []struct {
		name       string
		backupPath string
		setupFunc  func(string) error
		mockLogger *mocks.MockSetupLogger
		wantErr    bool
	}{
		{
			name:       "should restore config successfully",
			backupPath: "test_backup.yaml",
			setupFunc: func(backupDir string) error {
				// Criar arquivo de backup
				backupFile := filepath.Join(backupDir, "test_backup.yaml")
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
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    false,
		},
		{
			name:       "should fail when backup file does not exist",
			backupPath: "nonexistent_backup.yaml",
			setupFunc: func(backupDir string) error {
				// Não criar arquivo de backup
				return nil
			},
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    true,
		},
		{
			name:       "should fail when backup file is invalid",
			backupPath: "invalid_backup.yaml",
			setupFunc: func(backupDir string) error {
				// Criar arquivo de backup inválido
				backupFile := filepath.Join(backupDir, "invalid_backup.yaml")
				backupContent := "invalid yaml content: ["
				return os.WriteFile(backupFile, []byte(backupContent), 0644)
			},
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "configurator_restore_test")
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

			// Mock do logger
			logger := &mocks.MockSetupLogger{}

			// Criar configurator com diretório temporário
			configurator := &src.Configurator{
				Logger: logger,
			}

			fullBackupPath := filepath.Join(backupDir, tt.backupPath)
			err := configurator.RestoreConfig(fullBackupPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Configurator.RestoreConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestConfigurator_EdgeCases testa casos extremos do Configurator
func TestConfigurator_EdgeCases(t *testing.T) {
	t.Run("should handle nil logger", func(t *testing.T) {
		configurator := &src.Configurator{
			Logger: nil,
		}

		// Should not panic
		err := configurator.CreateStructure()
		if err != nil {
			t.Errorf("CreateStructure() failed with nil logger: %v", err)
		}
	})

	t.Run("should handle empty backup name", func(t *testing.T) {
		tempDir := helpers.CreateTempDir(t, "configurator_edge_test")
		configDir := filepath.Join(tempDir, "config")
		os.MkdirAll(configDir, 0755)

		// Criar arquivo de configuração
		configPath := filepath.Join(configDir, "manager.yaml")
		configContent := `manager: {}`
		os.WriteFile(configPath, []byte(configContent), 0644)

		configurator := &src.Configurator{
			Logger: &mocks.MockSetupLogger{},
		}

		err := configurator.BackupConfig("")
		if err != nil {
			t.Errorf("BackupConfig() failed with empty name: %v", err)
		}

		os.RemoveAll(tempDir)
	})
}

// TestConfigurator_Concurrency testa concorrência do Configurator
func TestConfigurator_Concurrency(t *testing.T) {
	t.Run("should handle concurrent structure creation", func(t *testing.T) {
		tempDir := helpers.CreateTempDir(t, "configurator_concurrent_test")

		configurator := &src.Configurator{
			Logger: &mocks.MockSetupLogger{},
		}

		// Executar múltiplas chamadas de CreateStructure concorrentemente
		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func() {
				err := configurator.CreateStructure()
				if err != nil {
					t.Errorf("Concurrent CreateStructure() failed: %v", err)
				}
				done <- true
			}()
		}

		// Aguardar todas as goroutines terminarem
		for i := 0; i < 10; i++ {
			<-done
		}

		os.RemoveAll(tempDir)
	})
}

// TestConfigurator_Performance testa performance do Configurator
func TestConfigurator_Performance(t *testing.T) {
	t.Run("should complete operations within reasonable time", func(t *testing.T) {
		tempDir := helpers.CreateTempDir(t, "configurator_perf_test")

		configurator := &src.Configurator{
			Logger: &mocks.MockSetupLogger{},
		}

		start := time.Now()
		err := configurator.CreateStructure()
		elapsed := time.Since(start)

		if err != nil {
			t.Errorf("CreateStructure() failed: %v", err)
		}

		if elapsed > 1*time.Second {
			t.Errorf("CreateStructure() took too long: %v", elapsed)
		}

		os.RemoveAll(tempDir)
	})
}
