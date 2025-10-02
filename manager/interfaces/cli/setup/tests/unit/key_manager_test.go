//go:build !integration && !e2e && !performance && !security
// +build !integration,!e2e,!performance,!security

package unit

import (
	"os"
	"path/filepath"
	"testing"

	setup "setup-component/src"
)

// TestNewKeyManager testa a criação do gerenciador de chaves
func TestNewKeyManager(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	tests := []struct {
		name    string
		logger  *setup.SetupLogger
		wantErr bool
	}{
		{
			name:    "should create key manager successfully",
			logger:  logger,
			wantErr: false,
		},
		{
			name:    "should create key manager with nil logger",
			logger:  nil,
			wantErr: false, // Logger pode ser nil
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyManager := setup.NewKeyManager(tt.logger)
			if keyManager == nil {
				t.Error("NewKeyManager() returned nil key manager")
			}
		})
	}
}

// TestKeyManager_GenerateKeyPair testa a geração de par de chaves
func TestKeyManager_GenerateKeyPair(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	keyManager := setup.NewKeyManager(logger)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should generate key pair successfully",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyPair, err := keyManager.GenerateKeyPair("test-key")
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyManager.GenerateKeyPair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if keyPair == nil {
					t.Error("KeyManager.GenerateKeyPair() returned nil key pair")
					return
				}

				// Verificar campos obrigatórios
				if keyPair.ID == "" {
					t.Error("Key pair missing ID")
				}
				if keyPair.Algorithm == "" {
					t.Error("Key pair missing algorithm")
				}
				if keyPair.PrivateKey == nil {
					t.Error("Key pair missing private key")
				}
				if keyPair.PublicKey == nil {
					t.Error("Key pair missing public key")
				}
				if keyPair.Fingerprint == "" {
					t.Error("Key pair missing fingerprint")
				}
			}
		})
	}
}

// TestKeyManager_StoreKeyPair testa o armazenamento de par de chaves
func TestKeyManager_StoreKeyPair(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	keyManager := setup.NewKeyManager(logger)

	tests := []struct {
		name    string
		keyPair *setup.KeyPair
		wantErr bool
	}{
		{
			name:    "should fail store with nil key pair",
			keyPair: nil,
			wantErr: true,
		},
		{
			name: "should store key pair successfully",
			keyPair: &types.setup.KeyPair{
				ID:          "test-key",
				Algorithm:   "ed25519",
				PrivateKey:  []byte("private key"),
				PublicKey:   []byte("public key"),
				Fingerprint: "test-fingerprint",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := keyManager.StoreKeyPair(tt.keyPair)
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyManager.StoreKeyPair() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verificar se o arquivo foi salvo
			if !tt.wantErr {
				keyPath := filepath.Join(tempDir, ".syntropy", "keys", tt.keyPair.ID+".key")
				if _, err := os.Stat(keyPath); os.IsNotExist(err) {
					t.Errorf("Key file not saved: %s", keyPath)
				}
			}
		})
	}
}

// TestKeyManager_LoadKeyPair testa o carregamento de par de chaves
func TestKeyManager_LoadKeyPair(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	keyManager := setup.NewKeyManager(logger)

	tests := []struct {
		name    string
		keyID   string
		setup   bool
		wantErr bool
	}{
		{
			name:    "should fail load when key does not exist",
			keyID:   "nonexistent-key",
			setup:   false,
			wantErr: true,
		},
		{
			name:    "should load key pair successfully when key exists",
			keyID:   "test-key",
			setup:   true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar chave se necessário
			if tt.setup {
				keyPair := &types.setup.KeyPair{
					ID:          tt.keyID,
					Algorithm:   "ed25519",
					PrivateKey:  []byte("private key"),
					PublicKey:   []byte("public key"),
					Fingerprint: "test-fingerprint",
				}
				err := keyManager.StoreKeyPair(keyPair)
				if err != nil {
					t.Fatalf("Failed to save key pair: %v", err)
				}
			}

			keyPair, err := keyManager.LoadKeyPair(tt.keyID, "test-password")
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyManager.LoadKeyPair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if keyPair == nil {
					t.Error("KeyManager.LoadKeyPair() returned nil key pair")
					return
				}

				// Verificar campos obrigatórios
				if keyPair.ID != tt.keyID {
					t.Errorf("Key pair ID mismatch: got %s, want %s", keyPair.ID, tt.keyID)
				}
				if keyPair.Algorithm == "" {
					t.Error("Key pair missing algorithm")
				}
				if keyPair.PrivateKey == nil {
					t.Error("Key pair missing private key")
				}
				if keyPair.PublicKey == nil {
					t.Error("Key pair missing public key")
				}
				if keyPair.Fingerprint == "" {
					t.Error("Key pair missing fingerprint")
				}
			}
		})
	}
}

// TestKeyManager_ListKeys testa a listagem de chaves
func TestKeyManager_ListKeys(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	keyManager := setup.NewKeyManager(logger)

	tests := []struct {
		name    string
		setup   bool
		wantErr bool
	}{
		{
			name:    "should list keys successfully when keys exist",
			setup:   true,
			wantErr: false,
		},
		{
			name:    "should list keys successfully when no keys exist",
			setup:   false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar chaves se necessário
			if tt.setup {
				keyPairs := []string{"key1", "key2", "key3"}
				for _, keyID := range keyPairs {
					keyPair := &types.setup.KeyPair{
						ID:          keyID,
						Algorithm:   "ed25519",
						PrivateKey:  []byte("private key"),
						PublicKey:   []byte("public key"),
						Fingerprint: "test-fingerprint",
					}
					err := keyManager.StoreKeyPair(keyPair)
					if err != nil {
						t.Fatalf("Failed to save key pair %s: %v", keyID, err)
					}
				}
			}

			keys, err := keyManager.ListKeys()
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyManager.ListKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if keys == nil {
					t.Error("KeyManager.ListKeys() returned nil keys")
				}
			}
		})
	}
}

// TestKeyManager_RotateKeys testa a rotação de chaves
func TestKeyManager_RotateKeys(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	keyManager := setup.NewKeyManager(logger)

	tests := []struct {
		name    string
		keyID   string
		setup   bool
		wantErr bool
	}{
		{
			name:    "should fail rotate when key does not exist",
			keyID:   "nonexistent-key",
			setup:   false,
			wantErr: true,
		},
		{
			name:    "should rotate keys successfully when key exists",
			keyID:   "test-key",
			setup:   true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar chave se necessário
			if tt.setup {
				keyPair := &types.setup.KeyPair{
					ID:          tt.keyID,
					Algorithm:   "ed25519",
					PrivateKey:  []byte("private key"),
					PublicKey:   []byte("public key"),
					Fingerprint: "test-fingerprint",
				}
				err := keyManager.StoreKeyPair(keyPair)
				if err != nil {
					t.Fatalf("Failed to save key pair: %v", err)
				}
			}

			err := keyManager.RotateKeys(tt.keyID)
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyManager.RotateKeys() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestKeyManager_VerifyKeyIntegrity testa a verificação de integridade de chave
func TestKeyManager_VerifyKeyIntegrity(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	keyManager := setup.NewKeyManager(logger)

	tests := []struct {
		name    string
		keyID   string
		setup   bool
		wantErr bool
	}{
		{
			name:    "should fail verify when key does not exist",
			keyID:   "nonexistent-key",
			setup:   false,
			wantErr: true,
		},
		{
			name:    "should verify key integrity successfully when key exists",
			keyID:   "test-key",
			setup:   true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar chave se necessário
			if tt.setup {
				keyPair := &types.setup.KeyPair{
					ID:          tt.keyID,
					Algorithm:   "ed25519",
					PrivateKey:  []byte("private key"),
					PublicKey:   []byte("public key"),
					Fingerprint: "test-fingerprint",
				}
				err := keyManager.StoreKeyPair(keyPair)
				if err != nil {
					t.Fatalf("Failed to save key pair: %v", err)
				}
			}

			valid, err := keyManager.VerifyKeyIntegrity(tt.keyID)
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyManager.VerifyKeyIntegrity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !valid {
					t.Error("KeyManager.VerifyKeyIntegrity() returned false for valid key")
				}
			}
		})
	}
}

// TestKeyManager_BackupKeys testa o backup de chaves
func TestKeyManager_BackupKeys(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	keyManager := setup.NewKeyManager(logger)

	tests := []struct {
		name    string
		keyID   string
		setup   bool
		wantErr bool
	}{
		{
			name:    "should fail backup when key does not exist",
			keyID:   "nonexistent-key",
			setup:   false,
			wantErr: true,
		},
		{
			name:    "should backup keys successfully when key exists",
			keyID:   "test-key",
			setup:   true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar chave se necessário
			if tt.setup {
				keyPair := &types.setup.KeyPair{
					ID:          tt.keyID,
					Algorithm:   "ed25519",
					PrivateKey:  []byte("private key"),
					PublicKey:   []byte("public key"),
					Fingerprint: "test-fingerprint",
				}
				err := keyManager.StoreKeyPair(keyPair)
				if err != nil {
					t.Fatalf("Failed to save key pair: %v", err)
				}
			}

			backupPath, err := keyManager.BackupKeys(tt.keyID, "test_backup")
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyManager.BackupKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if backupPath == "" {
					t.Error("KeyManager.BackupKeys() returned empty backup path")
				}

				// Verificar se o backup foi criado
				if _, err := os.Stat(backupPath); os.IsNotExist(err) {
					t.Errorf("Backup file not created: %s", backupPath)
				}
			}
		})
	}
}

// TestKeyManager_RestoreKeys testa a restauração de chaves
func TestKeyManager_RestoreKeys(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	keyManager := setup.NewKeyManager(logger)

	tests := []struct {
		name    string
		keyID   string
		setup   bool
		wantErr bool
	}{
		{
			name:    "should fail restore when backup does not exist",
			keyID:   "nonexistent-key",
			setup:   false,
			wantErr: true,
		},
		{
			name:    "should restore keys successfully when backup exists",
			keyID:   "test-key",
			setup:   true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var backupPath string

			// Criar backup se necessário
			if tt.setup {
				keyPair := &types.setup.KeyPair{
					ID:          tt.keyID,
					Algorithm:   "ed25519",
					PrivateKey:  []byte("private key"),
					PublicKey:   []byte("public key"),
					Fingerprint: "test-fingerprint",
				}
				err := keyManager.StoreKeyPair(keyPair)
				if err != nil {
					t.Fatalf("Failed to save key pair: %v", err)
				}

				backupPath, err = keyManager.BackupKeys(tt.keyID, "test_backup")
				if err != nil {
					t.Fatalf("Failed to backup key pair: %v", err)
				}
			} else {
				backupPath = filepath.Join(tempDir, "nonexistent_backup.key")
			}

			err := keyManager.RestoreKeys(backupPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyManager.RestoreKeys() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
