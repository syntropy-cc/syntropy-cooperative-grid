//go:build integration
// +build integration

package integration

import (
	"os"
	"path/filepath"
	"testing"

	setup "setup-component/src"
)

// TestSetupManager_Integration testa a integração completa do SetupManager
func TestSetupManager_Integration(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	manager := setup.NewSetupManager(logger)

	tests := []struct {
		name    string
		options *setup.SetupOptions
		wantErr bool
	}{
		{
			name: "should complete full setup successfully",
			options: &types.setup.SetupOptions{
				Force:          false,
				SkipValidation: false,
			},
			wantErr: false,
		},
		{
			name: "should complete setup with force flag",
			options: &types.setup.SetupOptions{
				Force:          true,
				SkipValidation: false,
			},
			wantErr: false,
		},
		{
			name: "should complete setup skipping validation",
			options: &types.setup.SetupOptions{
				Force:          false,
				SkipValidation: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Executar setup
			err := manager.Setup(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupManager.Setup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verificar se o setup foi concluído com sucesso
				status, err := manager.Status()
				if err != nil {
					t.Errorf("SetupManager.Status() error = %v", err)
					return
				}

				if status.Status != "completed" {
					t.Errorf("Setup status = %s, want completed", status.Status)
				}

				// Verificar se os arquivos necessários foram criados
				configPath := filepath.Join(tempDir, ".syntropy", "config", "manager.yaml")
				if _, err := os.Stat(configPath); os.IsNotExist(err) {
					t.Errorf("Config file not created: %s", configPath)
				}

				statePath := filepath.Join(tempDir, ".syntropy", "state", "setup_state.json")
				if _, err := os.Stat(statePath); os.IsNotExist(err) {
					t.Errorf("State file not created: %s", statePath)
				}

				keysDir := filepath.Join(tempDir, ".syntropy", "keys")
				if _, err := os.Stat(keysDir); os.IsNotExist(err) {
					t.Errorf("Keys directory not created: %s", keysDir)
				}

				logsDir := filepath.Join(tempDir, ".syntropy", "logs")
				if _, err := os.Stat(logsDir); os.IsNotExist(err) {
					t.Errorf("Logs directory not created: %s", logsDir)
				}
			}
		})
	}
}

// TestSetupManager_Validation_Integration testa a integração da validação
func TestSetupManager_Validation_Integration(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	manager := setup.NewSetupManager(logger)

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
			// Executar validação
			err := manager.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupManager.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestSetupManager_Status_Integration testa a integração do status
func TestSetupManager_Status_Integration(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	manager := setup.NewSetupManager(logger)

	tests := []struct {
		name    string
		setup   bool
		wantErr bool
	}{
		{
			name:    "should get status when no setup exists",
			setup:   false,
			wantErr: false,
		},
		{
			name:    "should get status when setup exists",
			setup:   true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Executar setup se necessário
			if tt.setup {
				options := &types.setup.SetupOptions{
					Force:          true,
					SkipValidation: false,
				}
				err := manager.Setup(options)
				if err != nil {
					t.Fatalf("Failed to setup: %v", err)
				}
			}

			// Obter status
			status, err := manager.Status()
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupManager.Status() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if status == nil {
					t.Error("SetupManager.Status() returned nil status")
					return
				}

				// Verificar campos obrigatórios
				if status.Status == "" {
					t.Error("Status missing status field")
				}
				if status.Version == "" {
					t.Error("Status missing version field")
				}
				if status.Timestamp == "" {
					t.Error("Status missing timestamp field")
				}
			}
		})
	}
}

// TestSetupManager_Reset_Integration testa a integração do reset
func TestSetupManager_Reset_Integration(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	manager := setup.NewSetupManager(logger)

	tests := []struct {
		name    string
		setup   bool
		wantErr bool
	}{
		{
			name:    "should reset when no setup exists",
			setup:   false,
			wantErr: false,
		},
		{
			name:    "should reset when setup exists",
			setup:   true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Executar setup se necessário
			if tt.setup {
				options := &types.setup.SetupOptions{
					Force:          true,
					SkipValidation: false,
				}
				err := manager.Setup(options)
				if err != nil {
					t.Fatalf("Failed to setup: %v", err)
				}
			}

			// Executar reset
			err := manager.Reset()
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupManager.Reset() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verificar se o reset foi executado
			if !tt.wantErr {
				status, err := manager.Status()
				if err != nil {
					t.Errorf("SetupManager.Status() error = %v", err)
					return
				}

				if status.Status != "not_started" {
					t.Errorf("Reset status = %s, want not_started", status.Status)
				}
			}
		})
	}
}

// TestSetupManager_Repair_Integration testa a integração do repair
func TestSetupManager_Repair_Integration(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	manager := setup.NewSetupManager(logger)

	tests := []struct {
		name    string
		setup   bool
		wantErr bool
	}{
		{
			name:    "should repair when no setup exists",
			setup:   false,
			wantErr: false,
		},
		{
			name:    "should repair when setup exists",
			setup:   true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Executar setup se necessário
			if tt.setup {
				options := &types.setup.SetupOptions{
					Force:          true,
					SkipValidation: false,
				}
				err := manager.Setup(options)
				if err != nil {
					t.Fatalf("Failed to setup: %v", err)
				}
			}

			// Executar repair
			err := manager.Repair()
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupManager.Repair() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestSetupManager_Legacy_Integration testa a integração das funções legacy
func TestSetupManager_Legacy_Integration(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should execute legacy setup successfully",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Executar setup legacy
			err := setup.SetupLegacy()
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupLegacy() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verificar se o setup foi concluído
			if !tt.wantErr {
				status, err := setup.StatusLegacy()
				if err != nil {
					t.Errorf("StatusLegacy() error = %v", err)
					return
				}

				if status.Status != "completed" {
					t.Errorf("Legacy setup status = %s, want completed", status.Status)
				}
			}
		})
	}
}

// TestSetupManager_StatusLegacy_Integration testa a integração do status legacy
func TestSetupManager_StatusLegacy_Integration(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	tests := []struct {
		name    string
		setup   bool
		wantErr bool
	}{
		{
			name:    "should get legacy status when no setup exists",
			setup:   false,
			wantErr: false,
		},
		{
			name:    "should get legacy status when setup exists",
			setup:   true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Executar setup se necessário
			if tt.setup {
				err := setup.SetupLegacy()
				if err != nil {
					t.Fatalf("Failed to setup legacy: %v", err)
				}
			}

			// Obter status legacy
			status, err := setup.StatusLegacy()
			if (err != nil) != tt.wantErr {
				t.Errorf("StatusLegacy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if status == nil {
					t.Error("StatusLegacy() returned nil status")
					return
				}

				// Verificar campos obrigatórios
				if status.Status == "" {
					t.Error("Legacy status missing status field")
				}
				if status.Version == "" {
					t.Error("Legacy status missing version field")
				}
				if status.Timestamp == "" {
					t.Error("Legacy status missing timestamp field")
				}
			}
		})
	}
}

// TestSetupManager_ResetLegacy_Integration testa a integração do reset legacy
func TestSetupManager_ResetLegacy_Integration(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	tests := []struct {
		name    string
		setup   bool
		wantErr bool
	}{
		{
			name:    "should reset legacy when no setup exists",
			setup:   false,
			wantErr: false,
		},
		{
			name:    "should reset legacy when setup exists",
			setup:   true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Executar setup se necessário
			if tt.setup {
				err := setup.SetupLegacy()
				if err != nil {
					t.Fatalf("Failed to setup legacy: %v", err)
				}
			}

			// Executar reset legacy
			err := setup.ResetLegacy()
			if (err != nil) != tt.wantErr {
				t.Errorf("ResetLegacy() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verificar se o reset foi executado
			if !tt.wantErr {
				status, err := setup.StatusLegacy()
				if err != nil {
					t.Errorf("StatusLegacy() error = %v", err)
					return
				}

				if status.Status != "not_started" {
					t.Errorf("Legacy reset status = %s, want not_started", status.Status)
				}
			}
		})
	}
}
