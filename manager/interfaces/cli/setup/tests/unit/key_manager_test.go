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

// TestNewKeyManager testa a criação de um novo KeyManager
func TestNewKeyManager(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "should create key manager successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &mocks.MockSetupLogger{}
			keyManager := src.NewKeyManager(logger)
			if keyManager == nil {
				t.Error("NewKeyManager() returned nil key manager")
			}
		})
	}
}

// TestKeyManager_GenerateKeyPair testa o método GenerateKeyPair
func TestKeyManager_GenerateKeyPair(t *testing.T) {
	tests := []struct {
		name       string
		mockLogger *mocks.MockSetupLogger
		wantErr    bool
	}{
		{
			name:       "should generate key pair successfully",
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "key_manager_test")

			// Mock do logger
			logger := &mocks.MockSetupLogger{}

			// Criar key manager com diretório temporário
			keyManager := &src.KeyManager{
				Logger: logger,
			}

			keyPair, err := keyManager.GenerateKeyPair()
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyManager.GenerateKeyPair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if keyPair == nil {
					t.Error("KeyManager.GenerateKeyPair() returned nil keyPair")
					return
				}
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

// TestKeyManager_StoreKeyPair testa o método StoreKeyPair
func TestKeyManager_StoreKeyPair(t *testing.T) {
	tests := []struct {
		name       string
		keyPair    *types.KeyPair
		mockLogger *mocks.MockSetupLogger
		wantErr    bool
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
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    false,
		},
		{
			name:       "should fail when key pair is nil",
			keyPair:    nil,
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    true,
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
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "key_manager_store_test")
			keysDir := filepath.Join(tempDir, "keys")
			os.MkdirAll(keysDir, 0755)

			// Mock do logger
			logger := &mocks.MockSetupLogger{}

			// Criar key manager com diretório temporário
			keyManager := &src.KeyManager{
				Logger: logger,
			}

			err := keyManager.StoreKeyPair(tt.keyPair)
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyManager.StoreKeyPair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestKeyManager_LoadKeyPair testa o método LoadKeyPair
func TestKeyManager_LoadKeyPair(t *testing.T) {
	tests := []struct {
		name       string
		keyID      string
		setupFunc  func(string) error
		mockLogger *mocks.MockSetupLogger
		wantErr    bool
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
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    false,
		},
		{
			name:  "should fail when key file does not exist",
			keyID: "nonexistent-key",
			setupFunc: func(keysDir string) error {
				// Não criar arquivo de chave
				return nil
			},
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    true,
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
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "key_manager_load_test")
			keysDir := filepath.Join(tempDir, "keys")
			os.MkdirAll(keysDir, 0755)

			// Setup do teste
			if tt.setupFunc != nil {
				if err := tt.setupFunc(keysDir); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			// Mock do logger
			logger := &mocks.MockSetupLogger{}

			// Criar key manager com diretório temporário
			keyManager := &src.KeyManager{
				Logger: logger,
			}

			keyPair, err := keyManager.LoadKeyPair(tt.keyID)
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyManager.LoadKeyPair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if keyPair == nil {
					t.Error("KeyManager.LoadKeyPair() returned nil keyPair")
					return
				}
				helpers.AssertStringEqual(t, keyPair.ID, tt.keyID, "KeyPair ID")
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestKeyManager_RotateKeys testa o método RotateKeys
func TestKeyManager_RotateKeys(t *testing.T) {
	tests := []struct {
		name       string
		mockLogger *mocks.MockSetupLogger
		wantErr    bool
	}{
		{
			name:       "should rotate keys successfully",
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "key_manager_rotate_test")
			keysDir := filepath.Join(tempDir, "keys")
			os.MkdirAll(keysDir, 0755)

			// Mock do logger
			logger := &mocks.MockSetupLogger{}

			// Criar key manager com diretório temporário
			keyManager := &src.KeyManager{
				Logger: logger,
			}

			newKeyPair, err := keyManager.RotateKeys()
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyManager.RotateKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if newKeyPair == nil {
					t.Error("KeyManager.RotateKeys() returned nil keyPair")
					return
				}
				helpers.AssertStringNotEmpty(t, newKeyPair.ID, "New KeyPair ID")
				helpers.AssertStringEqual(t, newKeyPair.Algorithm, "ed25519", "New KeyPair Algorithm")
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestKeyManager_VerifyKeyIntegrity testa o método VerifyKeyIntegrity
func TestKeyManager_VerifyKeyIntegrity(t *testing.T) {
	tests := []struct {
		name       string
		keyID      string
		setupFunc  func(string) error
		mockLogger *mocks.MockSetupLogger
		wantErr    bool
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
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    false,
		},
		{
			name:  "should fail when key file does not exist",
			keyID: "nonexistent-key",
			setupFunc: func(keysDir string) error {
				// Não criar arquivo de chave
				return nil
			},
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    true,
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
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "key_manager_verify_test")
			keysDir := filepath.Join(tempDir, "keys")
			os.MkdirAll(keysDir, 0755)

			// Setup do teste
			if tt.setupFunc != nil {
				if err := tt.setupFunc(keysDir); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			// Mock do logger
			logger := &mocks.MockSetupLogger{}

			// Criar key manager com diretório temporário
			keyManager := &src.KeyManager{
				Logger: logger,
			}

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

// TestKeyManager_BackupKeys testa o método BackupKeys
func TestKeyManager_BackupKeys(t *testing.T) {
	tests := []struct {
		name       string
		backupPath string
		setupFunc  func(string) error
		mockLogger *mocks.MockSetupLogger
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
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    false,
		},
		{
			name:       "should fail when keys directory does not exist",
			backupPath: "keys_backup.tar.gz",
			setupFunc: func(keysDir string) error {
				// Não criar diretório de chaves
				return nil
			},
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    true,
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
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "key_manager_backup_test")
			keysDir := filepath.Join(tempDir, "keys")
			backupDir := filepath.Join(tempDir, "backups")
			os.MkdirAll(backupDir, 0755)

			// Setup do teste
			if tt.setupFunc != nil {
				if err := tt.setupFunc(keysDir); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			// Mock do logger
			logger := &mocks.MockSetupLogger{}

			// Criar key manager com diretório temporário
			keyManager := &src.KeyManager{
				Logger: logger,
			}

			fullBackupPath := filepath.Join(backupDir, tt.backupPath)
			err := keyManager.BackupKeys(fullBackupPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyManager.BackupKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestKeyManager_RestoreKeys testa o método RestoreKeys
func TestKeyManager_RestoreKeys(t *testing.T) {
	tests := []struct {
		name       string
		backupPath string
		setupFunc  func(string) error
		mockLogger *mocks.MockSetupLogger
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
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    false,
		},
		{
			name:       "should fail when backup file does not exist",
			backupPath: "nonexistent_backup.tar.gz",
			setupFunc: func(backupDir string) error {
				// Não criar arquivo de backup
				return nil
			},
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    true,
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
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "key_manager_restore_test")
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

			// Mock do logger
			logger := &mocks.MockSetupLogger{}

			// Criar key manager com diretório temporário
			keyManager := &src.KeyManager{
				Logger: logger,
			}

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

// TestKeyManager_ListKeys testa o método ListKeys
func TestKeyManager_ListKeys(t *testing.T) {
	tests := []struct {
		name       string
		setupFunc  func(string) error
		mockLogger *mocks.MockSetupLogger
		wantErr    bool
		wantCount  int
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
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    false,
			wantCount:  3,
		},
		{
			name: "should return empty list when no keys exist",
			setupFunc: func(keysDir string) error {
				// Não criar arquivos de chave
				return nil
			},
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    false,
			wantCount:  0,
		},
		{
			name: "should fail when keys directory does not exist",
			setupFunc: func(keysDir string) error {
				// Não criar diretório de chaves
				return nil
			},
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    true,
			wantCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "key_manager_list_test")
			keysDir := filepath.Join(tempDir, "keys")

			// Setup do teste
			if tt.setupFunc != nil {
				if err := tt.setupFunc(keysDir); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			// Mock do logger
			logger := &mocks.MockSetupLogger{}

			// Criar key manager com diretório temporário
			keyManager := &src.KeyManager{
				Logger: logger,
			}

			keys, err := keyManager.ListKeys()
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyManager.ListKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(keys) != tt.wantCount {
					t.Errorf("KeyManager.ListKeys() returned %d keys, want %d", len(keys), tt.wantCount)
				}
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestKeyManager_EdgeCases testa casos extremos do KeyManager
func TestKeyManager_EdgeCases(t *testing.T) {
	t.Run("should handle nil logger", func(t *testing.T) {
		keyManager := &src.KeyManager{
			Logger: nil,
		}

		// Should not panic
		keyPair, err := keyManager.GenerateKeyPair()
		if err != nil {
			t.Errorf("GenerateKeyPair() failed with nil logger: %v", err)
		}
		if keyPair == nil {
			t.Error("GenerateKeyPair() returned nil keyPair with nil logger")
		}
	})

	t.Run("should handle empty key ID", func(t *testing.T) {
		tempDir := helpers.CreateTempDir(t, "key_manager_edge_test")
		keysDir := filepath.Join(tempDir, "keys")
		os.MkdirAll(keysDir, 0755)

		keyManager := &src.KeyManager{
			Logger: &mocks.MockSetupLogger{},
		}

		_, err := keyManager.LoadKeyPair("")
		if err == nil {
			t.Error("LoadKeyPair() should fail with empty key ID")
		}

		os.RemoveAll(tempDir)
	})
}

// TestKeyManager_Concurrency testa concorrência do KeyManager
func TestKeyManager_Concurrency(t *testing.T) {
	t.Run("should handle concurrent key generation", func(t *testing.T) {
		tempDir := helpers.CreateTempDir(t, "key_manager_concurrent_test")

		keyManager := &src.KeyManager{
			Logger: &mocks.MockSetupLogger{},
		}

		// Executar múltiplas chamadas de GenerateKeyPair concorrentemente
		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func() {
				keyPair, err := keyManager.GenerateKeyPair()
				if err != nil {
					t.Errorf("Concurrent GenerateKeyPair() failed: %v", err)
				}
				if keyPair == nil {
					t.Error("Concurrent GenerateKeyPair() returned nil keyPair")
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

// TestKeyManager_Performance testa performance do KeyManager
func TestKeyManager_Performance(t *testing.T) {
	t.Run("should complete operations within reasonable time", func(t *testing.T) {
		tempDir := helpers.CreateTempDir(t, "key_manager_perf_test")

		keyManager := &src.KeyManager{
			Logger: &mocks.MockSetupLogger{},
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

		if elapsed > 1*time.Second {
			t.Errorf("GenerateKeyPair() took too long: %v", elapsed)
		}

		os.RemoveAll(tempDir)
	})
}

// TestKeyManager_Security testa aspectos de segurança do KeyManager
func TestKeyManager_Security(t *testing.T) {
	t.Run("should generate cryptographically secure keys", func(t *testing.T) {
		tempDir := helpers.CreateTempDir(t, "key_manager_security_test")

		keyManager := &src.KeyManager{
			Logger: &mocks.MockSetupLogger{},
		}

		// Gerar múltiplas chaves e verificar que são diferentes
		keyPairs := make([]*types.KeyPair, 10)
		for i := 0; i < 10; i++ {
			keyPair, err := keyManager.GenerateKeyPair()
			if err != nil {
				t.Errorf("GenerateKeyPair() failed: %v", err)
				return
			}
			keyPairs[i] = keyPair
		}

		// Verificar que todas as chaves são diferentes
		for i := 0; i < len(keyPairs); i++ {
			for j := i + 1; j < len(keyPairs); j++ {
				if keyPairs[i].ID == keyPairs[j].ID {
					t.Error("Generated duplicate key IDs")
				}
				if keyPairs[i].Fingerprint == keyPairs[j].Fingerprint {
					t.Error("Generated duplicate key fingerprints")
				}
				if keyPairs[i].PublicKey == keyPairs[j].PublicKey {
					t.Error("Generated duplicate public keys")
				}
				if keyPairs[i].PrivateKey == keyPairs[j].PrivateKey {
					t.Error("Generated duplicate private keys")
				}
			}
		}

		os.RemoveAll(tempDir)
	})

	t.Run("should handle key rotation securely", func(t *testing.T) {
		tempDir := helpers.CreateTempDir(t, "key_manager_rotation_test")

		keyManager := &src.KeyManager{
			Logger: &mocks.MockSetupLogger{},
		}

		// Gerar chave inicial
		initialKey, err := keyManager.GenerateKeyPair()
		if err != nil {
			t.Errorf("GenerateKeyPair() failed: %v", err)
			return
		}

		// Rotacionar chaves
		newKey, err := keyManager.RotateKeys()
		if err != nil {
			t.Errorf("RotateKeys() failed: %v", err)
			return
		}

		// Verificar que a nova chave é diferente da inicial
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

		os.RemoveAll(tempDir)
	})
}
