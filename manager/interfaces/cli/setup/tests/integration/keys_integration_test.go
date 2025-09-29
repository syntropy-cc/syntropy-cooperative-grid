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

// TestKeysIntegration_KeyGeneration testa a geração de chaves
func TestKeysIntegration_KeyGeneration(t *testing.T) {
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
			tempDir := helpers.CreateTempDir(t, "keys_generation_test")
			keysDir := filepath.Join(tempDir, "keys")
			os.MkdirAll(keysDir, 0755)

			// Criar key manager
			keyManager := src.NewKeyManager(nil)
			if keyManager == nil {
				t.Fatal("NewKeyManager() returned nil")
			}

			// Executar geração de chaves
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

// TestKeysIntegration_KeyStorage testa o armazenamento de chaves
func TestKeysIntegration_KeyStorage(t *testing.T) {
	tests := []struct {
		name    string
		keyPair *types.KeyPair
		wantErr bool
	}{
		{
			name: "should store key pair successfully",
			keyPair: &types.KeyPair{
				ID:          "test-key-1",
				Algorithm:   "ed25519",
				Fingerprint: "test-fingerprint",
				PublicKey:   "test-public-key",
				PrivateKey:  "test-private-key",
				CreatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name:    "should fail when key pair is nil",
			keyPair: nil,
			wantErr: true,
		},
		{
			name: "should fail when key pair has empty ID",
			keyPair: &types.KeyPair{
				ID:          "",
				Algorithm:   "ed25519",
				Fingerprint: "test-fingerprint",
				PublicKey:   "test-public-key",
				PrivateKey:  "test-private-key",
				CreatedAt:   time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "keys_storage_test")
			keysDir := filepath.Join(tempDir, "keys")
			os.MkdirAll(keysDir, 0755)

			// Criar key manager
			keyManager := src.NewKeyManager(nil)
			if keyManager == nil {
				t.Fatal("NewKeyManager() returned nil")
			}

			// Executar armazenamento de chave
			err := keyManager.StoreKeyPair(tt.keyPair)
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyManager.StoreKeyPair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verificar se a chave foi armazenada
				keyFile := filepath.Join(keysDir, tt.keyPair.ID+".json")
				helpers.AssertFileExists(t, keyFile, "Key file")
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestKeysIntegration_KeyLoading testa o carregamento de chaves
func TestKeysIntegration_KeyLoading(t *testing.T) {
	tests := []struct {
		name      string
		keyID     string
		setupFunc func(string) error
		wantErr   bool
	}{
		{
			name:  "should load key pair successfully",
			keyID: "test-key-1",
			setupFunc: func(keysDir string) error {
				// Criar arquivo de chave
				keyFile := filepath.Join(keysDir, "test-key-1.json")
				keyContent := `{
					"id": "test-key-1",
					"algorithm": "ed25519",
					"fingerprint": "test-fingerprint",
					"public_key": "test-public-key",
					"private_key": "test-private-key",
					"created_at": "2023-01-01T00:00:00Z"
				}`
				return os.WriteFile(keyFile, []byte(keyContent), 0644)
			},
			wantErr: false,
		},
		{
			name:  "should fail when key file does not exist",
			keyID: "nonexistent-key",
			setupFunc: func(keysDir string) error {
				// Não criar arquivo de chave
				return nil
			},
			wantErr: true,
		},
		{
			name:  "should fail when key file is invalid",
			keyID: "invalid-key",
			setupFunc: func(keysDir string) error {
				// Criar arquivo de chave inválido
				keyFile := filepath.Join(keysDir, "invalid-key.json")
				keyContent := "invalid json content"
				return os.WriteFile(keyFile, []byte(keyContent), 0644)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "keys_loading_test")
			keysDir := filepath.Join(tempDir, "keys")
			os.MkdirAll(keysDir, 0755)

			// Setup do teste
			if tt.setupFunc != nil {
				if err := tt.setupFunc(keysDir); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			// Criar key manager
			keyManager := src.NewKeyManager(nil)
			if keyManager == nil {
				t.Fatal("NewKeyManager() returned nil")
			}

			// Executar carregamento de chave
			keyPair, err := keyManager.LoadKeyPair(tt.keyID)
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyManager.LoadKeyPair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if keyPair == nil {
					t.Error("LoadKeyPair() returned nil keyPair")
					return
				}

				// Verificar se a chave foi carregada corretamente
				helpers.AssertStringEqual(t, keyPair.ID, tt.keyID, "KeyPair ID")
				helpers.AssertStringEqual(t, keyPair.Algorithm, "ed25519", "KeyPair Algorithm")
				helpers.AssertStringEqual(t, keyPair.Fingerprint, "test-fingerprint", "KeyPair Fingerprint")
				helpers.AssertStringEqual(t, keyPair.PublicKey, "test-public-key", "KeyPair PublicKey")
				helpers.AssertStringEqual(t, keyPair.PrivateKey, "test-private-key", "KeyPair PrivateKey")
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestKeysIntegration_KeyRotation testa a rotação de chaves
func TestKeysIntegration_KeyRotation(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should rotate keys successfully",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "keys_rotation_test")
			keysDir := filepath.Join(tempDir, "keys")
			os.MkdirAll(keysDir, 0755)

			// Criar key manager
			keyManager := src.NewKeyManager(nil)
			if keyManager == nil {
				t.Fatal("NewKeyManager() returned nil")
			}

			// Gerar chave inicial
			initialKey, err := keyManager.GenerateKeyPair()
			if err != nil {
				t.Errorf("GenerateKeyPair() failed: %v", err)
				return
			}

			// Executar rotação de chaves
			newKey, err := keyManager.RotateKeys()
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyManager.RotateKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if newKey == nil {
					t.Error("RotateKeys() returned nil keyPair")
					return
				}

				// Verificar se a nova chave é diferente da inicial
				if initialKey.ID == newKey.ID {
					t.Error("Key rotation generated same key ID")
				}
				if initialKey.Fingerprint == newKey.Fingerprint {
					t.Error("Key rotation generated same key fingerprint")
				}
				if initialKey.PublicKey == newKey.PublicKey {
					t.Error("Key rotation generated same public key")
				}
				if initialKey.PrivateKey == newKey.PrivateKey {
					t.Error("Key rotation generated same private key")
				}
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestKeysIntegration_KeyIntegrity testa a verificação de integridade de chaves
func TestKeysIntegration_KeyIntegrity(t *testing.T) {
	tests := []struct {
		name      string
		keyID     string
		setupFunc func(string) error
		wantErr   bool
	}{
		{
			name:  "should verify key integrity successfully",
			keyID: "test-key-1",
			setupFunc: func(keysDir string) error {
				// Criar arquivo de chave válido
				keyFile := filepath.Join(keysDir, "test-key-1.json")
				keyContent := `{
					"id": "test-key-1",
					"algorithm": "ed25519",
					"fingerprint": "test-fingerprint",
					"public_key": "test-public-key",
					"private_key": "test-private-key",
					"created_at": "2023-01-01T00:00:00Z"
				}`
				return os.WriteFile(keyFile, []byte(keyContent), 0644)
			},
			wantErr: false,
		},
		{
			name:  "should fail when key file does not exist",
			keyID: "nonexistent-key",
			setupFunc: func(keysDir string) error {
				// Não criar arquivo de chave
				return nil
			},
			wantErr: true,
		},
		{
			name:  "should fail when key file is corrupted",
			keyID: "corrupted-key",
			setupFunc: func(keysDir string) error {
				// Criar arquivo de chave corrompido
				keyFile := filepath.Join(keysDir, "corrupted-key.json")
				keyContent := "corrupted content"
				return os.WriteFile(keyFile, []byte(keyContent), 0644)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "keys_integrity_test")
			keysDir := filepath.Join(tempDir, "keys")
			os.MkdirAll(keysDir, 0755)

			// Setup do teste
			if tt.setupFunc != nil {
				if err := tt.setupFunc(keysDir); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			// Criar key manager
			keyManager := src.NewKeyManager(nil)
			if keyManager == nil {
				t.Fatal("NewKeyManager() returned nil")
			}

			// Executar verificação de integridade
			err := keyManager.VerifyKeyIntegrity(tt.keyID)
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyManager.VerifyKeyIntegrity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestKeysIntegration_KeyBackup testa o backup de chaves
func TestKeysIntegration_KeyBackup(t *testing.T) {
	tests := []struct {
		name       string
		backupPath string
		setupFunc  func(string) error
		wantErr    bool
	}{
		{
			name:       "should backup keys successfully",
			backupPath: "keys_backup.tar.gz",
			setupFunc: func(keysDir string) error {
				// Criar arquivo de chave
				keyFile := filepath.Join(keysDir, "test-key-1.json")
				keyContent := `{
					"id": "test-key-1",
					"algorithm": "ed25519",
					"fingerprint": "test-fingerprint",
					"public_key": "test-public-key",
					"private_key": "test-private-key",
					"created_at": "2023-01-01T00:00:00Z"
				}`
				return os.WriteFile(keyFile, []byte(keyContent), 0644)
			},
			wantErr: false,
		},
		{
			name:       "should handle empty backup path",
			backupPath: "",
			setupFunc: func(keysDir string) error {
				// Criar arquivo de chave
				keyFile := filepath.Join(keysDir, "test-key-1.json")
				keyContent := `{
					"id": "test-key-1",
					"algorithm": "ed25519",
					"fingerprint": "test-fingerprint",
					"public_key": "test-public-key",
					"private_key": "test-private-key",
					"created_at": "2023-01-01T00:00:00Z"
				}`
				return os.WriteFile(keyFile, []byte(keyContent), 0644)
			},
			wantErr: false,
		},
		{
			name:       "should fail when keys directory does not exist",
			backupPath: "keys_backup.tar.gz",
			setupFunc: func(keysDir string) error {
				// Não criar diretório de chaves
				return nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "keys_backup_test")
			keysDir := filepath.Join(tempDir, "keys")
			backupDir := filepath.Join(tempDir, "backups")
			os.MkdirAll(backupDir, 0755)

			// Setup do teste
			if tt.setupFunc != nil {
				if err := tt.setupFunc(keysDir); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			// Criar key manager
			keyManager := src.NewKeyManager(nil)
			if keyManager == nil {
				t.Fatal("NewKeyManager() returned nil")
			}

			// Executar backup de chaves
			fullBackupPath := filepath.Join(backupDir, tt.backupPath)
			err := keyManager.BackupKeys(fullBackupPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyManager.BackupKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verificar se o backup foi criado
				helpers.AssertFileExists(t, fullBackupPath, "Keys backup file")
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestKeysIntegration_KeyRestore testa a restauração de chaves
func TestKeysIntegration_KeyRestore(t *testing.T) {
	tests := []struct {
		name       string
		backupPath string
		setupFunc  func(string) error
		wantErr    bool
	}{
		{
			name:       "should restore keys successfully",
			backupPath: "keys_backup.tar.gz",
			setupFunc: func(backupDir string) error {
				// Criar arquivo de backup
				backupFile := filepath.Join(backupDir, "keys_backup.tar.gz")
				backupContent := []byte("fake backup content")
				return os.WriteFile(backupFile, backupContent, 0644)
			},
			wantErr: false,
		},
		{
			name:       "should fail when backup file does not exist",
			backupPath: "nonexistent_backup.tar.gz",
			setupFunc: func(backupDir string) error {
				// Não criar arquivo de backup
				return nil
			},
			wantErr: true,
		},
		{
			name:       "should fail when backup file is corrupted",
			backupPath: "corrupted_backup.tar.gz",
			setupFunc: func(backupDir string) error {
				// Criar arquivo de backup corrompido
				backupFile := filepath.Join(backupDir, "corrupted_backup.tar.gz")
				backupContent := []byte("corrupted backup content")
				return os.WriteFile(backupFile, backupContent, 0644)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "keys_restore_test")
			keysDir := filepath.Join(tempDir, "keys")
			backupDir := filepath.Join(tempDir, "backups")
			os.MkdirAll(keysDir, 0755)
			os.MkdirAll(backupDir, 0755)

			// Setup do teste
			if tt.setupFunc != nil {
				if err := tt.setupFunc(backupDir); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			// Criar key manager
			keyManager := src.NewKeyManager(nil)
			if keyManager == nil {
				t.Fatal("NewKeyManager() returned nil")
			}

			// Executar restauração de chaves
			fullBackupPath := filepath.Join(backupDir, tt.backupPath)
			err := keyManager.RestoreKeys(fullBackupPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyManager.RestoreKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestKeysIntegration_KeyListing testa a listagem de chaves
func TestKeysIntegration_KeyListing(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(string) error
		wantErr   bool
		wantCount int
	}{
		{
			name: "should list keys successfully",
			setupFunc: func(keysDir string) error {
				// Criar múltiplos arquivos de chave
				for i := 1; i <= 3; i++ {
					keyFile := filepath.Join(keysDir, "test-key-"+string(rune(i))+".json")
					keyContent := `{
						"id": "test-key-` + string(rune(i)) + `",
						"algorithm": "ed25519",
						"fingerprint": "test-fingerprint-` + string(rune(i)) + `",
						"public_key": "test-public-key-` + string(rune(i)) + `",
						"private_key": "test-private-key-` + string(rune(i)) + `",
						"created_at": "2023-01-01T00:00:00Z"
					}`
					if err := os.WriteFile(keyFile, []byte(keyContent), 0644); err != nil {
						return err
					}
				}
				return nil
			},
			wantErr:   false,
			wantCount: 3,
		},
		{
			name: "should return empty list when no keys exist",
			setupFunc: func(keysDir string) error {
				// Não criar arquivos de chave
				return nil
			},
			wantErr:   false,
			wantCount: 0,
		},
		{
			name: "should fail when keys directory does not exist",
			setupFunc: func(keysDir string) error {
				// Não criar diretório de chaves
				return nil
			},
			wantErr:   true,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "keys_listing_test")
			keysDir := filepath.Join(tempDir, "keys")

			// Setup do teste
			if tt.setupFunc != nil {
				if err := tt.setupFunc(keysDir); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			// Criar key manager
			keyManager := src.NewKeyManager(nil)
			if keyManager == nil {
				t.Fatal("NewKeyManager() returned nil")
			}

			// Executar listagem de chaves
			keys, err := keyManager.ListKeys()
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyManager.ListKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(keys) != tt.wantCount {
					t.Errorf("ListKeys() returned %d keys, want %d", len(keys), tt.wantCount)
				}
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestKeysIntegration_ErrorHandling testa o tratamento de erros
func TestKeysIntegration_ErrorHandling(t *testing.T) {
	tests := []struct {
		name    string
		keyID   string
		wantErr bool
	}{
		{
			name:    "should handle empty key ID",
			keyID:   "",
			wantErr: true,
		},
		{
			name:    "should handle invalid key ID",
			keyID:   "invalid-key-id",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "keys_error_test")
			keysDir := filepath.Join(tempDir, "keys")
			os.MkdirAll(keysDir, 0755)

			// Criar key manager
			keyManager := src.NewKeyManager(nil)
			if keyManager == nil {
				t.Fatal("NewKeyManager() returned nil")
			}

			// Executar operação que deve falhar
			_, err := keyManager.LoadKeyPair(tt.keyID)
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyManager.LoadKeyPair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestKeysIntegration_Concurrency testa concorrência
func TestKeysIntegration_Concurrency(t *testing.T) {
	t.Run("should handle concurrent key operations", func(t *testing.T) {
		// Criar diretório temporário para testes
		tempDir := helpers.CreateTempDir(t, "keys_concurrent_test")
		keysDir := filepath.Join(tempDir, "keys")
		os.MkdirAll(keysDir, 0755)

		// Criar key manager
		keyManager := src.NewKeyManager(nil)
		if keyManager == nil {
			t.Fatal("NewKeyManager() returned nil")
		}

		// Executar múltiplas operações de chave concorrentemente
		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func(instance int) {
				// Gerar chave
				keyPair, err := keyManager.GenerateKeyPair()
				if err != nil {
					t.Errorf("Concurrent GenerateKeyPair() failed: %v", err)
				}
				if keyPair == nil {
					t.Error("Concurrent GenerateKeyPair() returned nil keyPair")
				}

				// Armazenar chave
				if keyPair != nil {
					err = keyManager.StoreKeyPair(keyPair)
					if err != nil {
						t.Errorf("Concurrent StoreKeyPair() failed: %v", err)
					}
				}

				done <- true
			}(i)
		}

		// Aguardar todas as goroutines terminarem
		for i := 0; i < 10; i++ {
			<-done
		}

		// Limpar diretório temporário
		os.RemoveAll(tempDir)
	})
}

// TestKeysIntegration_Performance testa performance
func TestKeysIntegration_Performance(t *testing.T) {
	t.Run("should complete key operations within reasonable time", func(t *testing.T) {
		// Criar diretório temporário para testes
		tempDir := helpers.CreateTempDir(t, "keys_perf_test")
		keysDir := filepath.Join(tempDir, "keys")
		os.MkdirAll(keysDir, 0755)

		// Criar key manager
		keyManager := src.NewKeyManager(nil)
		if keyManager == nil {
			t.Fatal("NewKeyManager() returned nil")
		}

		start := time.Now()
		keyPair, err := keyManager.GenerateKeyPair()
		elapsed := time.Since(start)

		if err != nil {
			t.Errorf("GenerateKeyPair() failed: %v", err)
		}
		if keyPair == nil {
			t.Error("GenerateKeyPair() returned nil keyPair")
		}

		if elapsed > 5*time.Second {
			t.Errorf("GenerateKeyPair() took too long: %v", elapsed)
		}

		// Limpar diretório temporário
		os.RemoveAll(tempDir)
	})
}
