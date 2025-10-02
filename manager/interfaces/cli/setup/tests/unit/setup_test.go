//go:build !integration && !e2e && !performance && !security
// +build !integration,!e2e,!performance,!security

package unit

import (
	"os"
	"path/filepath"
	"testing"

	setup "setup-component/src"
)

// TestNewSetupManager testa a criação do gerenciador de setup
func TestNewSetupManager(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should create setup manager successfully",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := setup.NewSetupManager()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSetupManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && manager == nil {
				t.Error("NewSetupManager() returned nil manager")
			}
		})
	}
}

// TestSetupManager_Setup testa o fluxo principal de setup
func TestSetupManager_Setup(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	tests := []struct {
		name    string
		options *setup.SetupOptions
		wantErr bool
	}{
		{
			name: "should setup successfully with valid options",
			options: &setup.SetupOptions{
				Force:        false,
				ValidateOnly: false,
				Verbose:      false,
				Quiet:        false,
				CustomSettings: map[string]string{
					"owner_name":  "Test User",
					"owner_email": "test@example.com",
				},
			},
			wantErr: false,
		},
		{
			name: "should setup successfully with force option",
			options: &setup.SetupOptions{
				Force:        true,
				ValidateOnly: false,
				Verbose:      false,
				Quiet:        false,
				CustomSettings: map[string]string{
					"owner_name":  "Test User",
					"owner_email": "test@example.com",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := setup.NewSetupManager()
			if err != nil {
				t.Fatalf("Failed to create setup manager: %v", err)
			}
			// Logger será fechado automaticamente

			err = manager.Setup(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupManager.Setup() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verificar se os arquivos foram criados
			if !tt.wantErr {
				configPath := filepath.Join(tempDir, ".syntropy", "config", "manager.yaml")
				if _, err := os.Stat(configPath); os.IsNotExist(err) {
					t.Errorf("Config file not created: %s", configPath)
				}

				statePath := filepath.Join(tempDir, ".syntropy", "state", "setup_state.json")
				if _, err := os.Stat(statePath); os.IsNotExist(err) {
					t.Errorf("State file not created: %s", statePath)
				}

				keysPath := filepath.Join(tempDir, ".syntropy", "keys", "owner.key")
				if _, err := os.Stat(keysPath); os.IsNotExist(err) {
					t.Errorf("Private key file not created: %s", keysPath)
				}
			}
		})
	}
}

// TestSetupManager_Validate testa a validação do ambiente
func TestSetupManager_Validate(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should validate environment successfully",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := setup.NewSetupManager()
			if err != nil {
				t.Fatalf("Failed to create setup manager: %v", err)
			}
			// Logger será fechado automaticamente

			result, err := manager.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupManager.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if result == nil {
					t.Error("SetupManager.Validate() returned nil result")
					return
				}

				// Verificar campos obrigatórios
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
		})
	}
}

// TestSetupManager_Status testa a verificação de status
func TestSetupManager_Status(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	tests := []struct {
		name    string
		setup   bool
		wantErr bool
	}{
		{
			name:    "should return error when no setup exists",
			setup:   false,
			wantErr: true,
		},
		{
			name:    "should return status when setup exists",
			setup:   true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := setup.NewSetupManager()
			if err != nil {
				t.Fatalf("Failed to create setup manager: %v", err)
			}
			// Logger será fechado automaticamente

			// Executar setup se necessário
			if tt.setup {
				options := &types.setup.SetupOptions{
					Force:        true,
					ValidateOnly: false,
					Verbose:      false,
					Quiet:        true,
					CustomSettings: map[string]string{
						"owner_name":  "Test User",
						"owner_email": "test@example.com",
					},
				}
				err = manager.Setup(options)
				if err != nil {
					t.Fatalf("Failed to setup: %v", err)
				}
			}

			status, err := manager.Status()
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupManager.Status() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && status == nil {
				t.Error("SetupManager.Status() returned nil status")
			}
		})
	}
}

// TestSetupManager_Reset testa o reset do setup
func TestSetupManager_Reset(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	tests := []struct {
		name    string
		confirm bool
		wantErr bool
	}{
		{
			name:    "should return error when not confirmed",
			confirm: false,
			wantErr: true,
		},
		{
			name:    "should reset successfully when confirmed",
			confirm: true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := setup.NewSetupManager()
			if err != nil {
				t.Fatalf("Failed to create setup manager: %v", err)
			}
			// Logger será fechado automaticamente

			// Executar setup primeiro
			options := &types.setup.SetupOptions{
				Force:        true,
				ValidateOnly: false,
				Verbose:      false,
				Quiet:        true,
				CustomSettings: map[string]string{
					"owner_name":  "Test User",
					"owner_email": "test@example.com",
				},
			}
			err = manager.Setup(options)
			if err != nil {
				t.Fatalf("Failed to setup: %v", err)
			}

			// Verificar se os arquivos existem antes do reset
			statePath := filepath.Join(tempDir, ".syntropy", "state", "setup_state.json")
			if _, err := os.Stat(statePath); os.IsNotExist(err) {
				t.Fatalf("State file should exist before reset: %s", statePath)
			}

			err = manager.Reset(tt.confirm)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupManager.Reset() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verificar se os arquivos foram removidos após reset bem-sucedido
			if !tt.wantErr {
				if _, err := os.Stat(statePath); !os.IsNotExist(err) {
					t.Errorf("State file should be removed after reset: %s", statePath)
				}
			}
		})
	}
}

// TestSetupManager_Repair testa o reparo do setup
func TestSetupManager_Repair(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should repair successfully",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := setup.NewSetupManager()
			if err != nil {
				t.Fatalf("Failed to create setup manager: %v", err)
			}
			// Logger será fechado automaticamente

			err = manager.Repair()
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupManager.Repair() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
