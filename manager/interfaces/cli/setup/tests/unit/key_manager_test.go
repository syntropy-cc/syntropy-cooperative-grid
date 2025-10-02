package unit

import (
	"os"
	"path/filepath"
	"testing"

	setup "setup-component/src"
)

// TestNewKeyManager testa a criação de um novo gerenciador de chaves
func TestNewKeyManager(t *testing.T) {
	logger := setup.NewSetupLogger()
	defer logger.Close()

	keyManager := setup.NewKeyManager(logger)
	if keyManager == nil {
		t.Error("NewKeyManager() returned nil")
	}
}

// TestKeyManager_GenerateKeyPair testa a geração de par de chaves através do SetupManager
func TestKeyManager_GenerateKeyPair(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	manager, err := setup.NewSetupManager()
	if err != nil {
		t.Fatalf("Failed to create SetupManager: %v", err)
	}

	tests := []struct {
		name    string
		options *setup.SetupOptions
		wantErr bool
	}{
		{
			name: "should generate key pair successfully",
			options: &setup.SetupOptions{
				Force:          false,
				SkipValidation: false,
				CustomSettings: map[string]string{
					"algorithm": "ed25519",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.SetupWithPublicOptions(tt.options)
		if (err != nil) != tt.wantErr {
				t.Errorf("SetupManager.Setup() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verificar se os arquivos de chave foram criados
			if !tt.wantErr {
				keyPath := filepath.Join(tempDir, ".syntropy", "keys", "owner.key")
				if _, err := os.Stat(keyPath); os.IsNotExist(err) {
					t.Errorf("Key file not created: %s", keyPath)
				}
			}
		})
	}
}

// TestKeyManager_StoreKeyPair testa o armazenamento de par de chaves através do SetupManager
func TestKeyManager_StoreKeyPair(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	manager, err := setup.NewSetupManager()
	if err != nil {
		t.Fatalf("Failed to create SetupManager: %v", err)
	}

	tests := []struct {
		name    string
		options *setup.SetupOptions
		wantErr bool
	}{
		{
			name: "should generate and store key pair successfully",
			options: &setup.SetupOptions{
				Force:          false,
				SkipValidation: false,
				CustomSettings: map[string]string{
					"algorithm": "ed25519",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.SetupWithPublicOptions(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupManager.Setup() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verificar se os arquivos de chave foram criados
			if !tt.wantErr {
				keyPath := filepath.Join(tempDir, ".syntropy", "keys", "owner.key")
				if _, err := os.Stat(keyPath); os.IsNotExist(err) {
					t.Errorf("Key file not saved: %s", keyPath)
				}
			}
		})
	}
}

// TestKeyManager_LoadKeyPair testa o carregamento de par de chaves através do SetupManager
func TestKeyManager_LoadKeyPair(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	manager, err := setup.NewSetupManager()
	if err != nil {
		t.Fatalf("Failed to create SetupManager: %v", err)
	}

	// Primeiro, criar uma chave
	options := &setup.SetupOptions{
		Force:          false,
		SkipValidation: false,
		CustomSettings: map[string]string{
			"algorithm": "ed25519",
		},
	}

	err = manager.SetupWithPublicOptions(options)
	if err != nil {
		t.Fatalf("Failed to setup keys: %v", err)
	}

	// Verificar se a chave foi criada
	keyPath := filepath.Join(tempDir, ".syntropy", "keys", "owner.key")
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		t.Errorf("Key file not found: %s", keyPath)
	}

	// Testar status para verificar se as chaves estão carregadas
	status, err := manager.Status()
	if err != nil {
		t.Errorf("Status failed: %v", err)
	}
	if status == nil {
		t.Error("Status should not be nil")
	}
}

// TestKeyManager_Validation testa a validação através do SetupManager
func TestKeyManager_Validation(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	manager, err := setup.NewSetupManager()
	if err != nil {
		t.Fatalf("Failed to create SetupManager: %v", err)
	}

	// Primeiro, criar uma chave
	options := &setup.SetupOptions{
		Force:          false,
		SkipValidation: false,
		CustomSettings: map[string]string{
			"algorithm": "ed25519",
		},
	}

	err = manager.SetupWithPublicOptions(options)
	if err != nil {
		t.Fatalf("Failed to setup keys: %v", err)
	}

	// Testar validação
	_, err = manager.Validate()
	if err != nil {
		t.Errorf("Validation failed: %v", err)
	}
}

// TestKeyManager_Reset testa o reset através do SetupManager
func TestKeyManager_Reset(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	manager, err := setup.NewSetupManager()
	if err != nil {
		t.Fatalf("Failed to create SetupManager: %v", err)
	}

	// Primeiro, criar uma chave
	options := &setup.SetupOptions{
		Force:          false,
		SkipValidation: false,
		CustomSettings: map[string]string{
			"algorithm": "ed25519",
		},
	}

	err = manager.SetupWithPublicOptions(options)
	if err != nil {
		t.Fatalf("Failed to setup keys: %v", err)
	}

	// Testar reset
	err = manager.Reset(true)
	if err != nil {
		t.Errorf("Reset failed: %v", err)
	}
}

// TestKeyManager_Repair testa o reparo através do SetupManager
func TestKeyManager_Repair(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	manager, err := setup.NewSetupManager()
	if err != nil {
		t.Fatalf("Failed to create SetupManager: %v", err)
	}

	// Testar reparo
	err = manager.Repair()
	if err != nil {
		t.Errorf("Repair failed: %v", err)
	}
}
