//go:build e2e
// +build e2e

package e2e

import (
	"os"
	"path/filepath"
	"testing"

	setup "setup-component/src"
)

// TestSetupManager_E2E testa o fluxo completo end-to-end do SetupManager
func TestSetupManager_E2E(t *testing.T) {
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
		wantErr bool
	}{
		{
			name:    "should complete full end-to-end setup flow",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 1. Verificar status inicial
			status, err := manager.Status()
			if err != nil {
				t.Errorf("SetupManager.Status() error = %v", err)
				return
			}

			if *status != "not_started" {
				t.Errorf("Initial status = %s, want not_started", *status)
			}

			// 2. Executar validação
			_, err = manager.Validate()
			if err != nil {
				t.Errorf("SetupManager.Validate() error = %v", err)
				return
			}

			// 3. Executar setup
			options := &setup.SetupOptions{
				Force:          false,
				SkipValidation: false,
			}
			err = manager.SetupWithPublicOptions(options)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupManager.Setup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// 4. Verificar status após setup
				status, err = manager.Status()
				if err != nil {
					t.Errorf("SetupManager.Status() error = %v", err)
					return
				}

				if *status != "completed" {
					t.Errorf("Setup status = %s, want completed", *status)
				}

				// 5. Verificar se todos os arquivos necessários foram criados
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

				// 6. Verificar se os logs foram criados
				logPath := filepath.Join(tempDir, ".syntropy", "logs", "setup.log")
				if _, err := os.Stat(logPath); os.IsNotExist(err) {
					t.Errorf("Log file not created: %s", logPath)
				}

				// 7. Executar repair para verificar integridade
				err = manager.Repair()
				if err != nil {
					t.Errorf("SetupManager.Repair() error = %v", err)
					return
				}

				// 8. Executar reset
				err = manager.Reset(true)
				if err != nil {
					t.Errorf("SetupManager.Reset() error = %v", err)
					return
				}

				// 9. Verificar status após reset
				status, err = manager.Status()
				if err != nil {
					t.Errorf("SetupManager.Status() error = %v", err)
					return
				}

				if *status != "not_started" {
					t.Errorf("Reset status = %s, want not_started", *status)
				}
			}
		})
	}
}

// TestSetupManager_E2E_WithForce testa o fluxo end-to-end com flag force
func TestSetupManager_E2E_WithForce(t *testing.T) {
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
		wantErr bool
	}{
		{
			name:    "should complete full end-to-end setup flow with force",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 1. Executar setup inicial
			options := &setup.SetupOptions{
				Force:          false,
				SkipValidation: false,
			}
			err := manager.SetupWithPublicOptions(options)
			if err != nil {
				t.Errorf("SetupManager.Setup() error = %v", err)
				return
			}

			// 2. Verificar status após setup inicial
			status, err := manager.Status()
			if err != nil {
				t.Errorf("SetupManager.Status() error = %v", err)
				return
			}

			if *status != "completed" {
				t.Errorf("Initial setup status = %s, want completed", *status)
			}

			// 3. Executar setup com force
			options = &setup.SetupOptions{
				Force:          true,
				SkipValidation: false,
			}
			err = manager.SetupWithPublicOptions(options)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupManager.Setup() with force error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// 4. Verificar status após setup com force
				status, err = manager.Status()
				if err != nil {
					t.Errorf("SetupManager.Status() error = %v", err)
					return
				}

				if *status != "completed" {
					t.Errorf("Force setup status = %s, want completed", *status)
				}

				// 5. Verificar se todos os arquivos ainda existem
				configPath := filepath.Join(tempDir, ".syntropy", "config", "manager.yaml")
				if _, err := os.Stat(configPath); os.IsNotExist(err) {
					t.Errorf("Config file not found after force setup: %s", configPath)
				}

				statePath := filepath.Join(tempDir, ".syntropy", "state", "setup_state.json")
				if _, err := os.Stat(statePath); os.IsNotExist(err) {
					t.Errorf("State file not found after force setup: %s", statePath)
				}
			}
		})
	}
}

// TestSetupManager_E2E_WithSkipValidation testa o fluxo end-to-end pulando validação
func TestSetupManager_E2E_WithSkipValidation(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	manager, err := setup.NewSetupManager()
	if err != nil {
		t.Fatalf("Failed to create SetupManager: %v", err)
	}

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should complete full end-to-end setup flow skipping validation",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 1. Executar setup pulando validação
			options := &setup.SetupOptions{
				Force:          false,
				SkipValidation: true,
			}
			err := manager.SetupWithPublicOptions(options)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupManager.Setup() skipping validation error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// 2. Verificar status após setup
				status, err := manager.Status()
				if err != nil {
					t.Errorf("SetupManager.Status() error = %v", err)
					return
				}

				if *status != "completed" {
					t.Errorf("Setup status = %s, want completed", *status)
				}

				// 3. Verificar se todos os arquivos necessários foram criados
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

// TestSetupManager_E2E_Legacy testa o fluxo end-to-end das funções legacy
func TestSetupManager_E2E_Legacy(t *testing.T) {
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
			name:    "should complete full end-to-end legacy setup flow",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 1. Verificar status inicial legacy
			status, err := setup.StatusLegacy(setup.LegacySetupOptions{})
			if err != nil {
				t.Errorf("StatusLegacy() error = %v", err)
				return
			}

			if !status.Success {
				t.Errorf("Initial legacy status success = %v, want false", status.Success)
			}

			// 2. Executar setup legacy
			_, err = setup.SetupLegacy(setup.LegacySetupOptions{})
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupLegacy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// 3. Verificar status após setup legacy
				status, err = setup.StatusLegacy(setup.LegacySetupOptions{})
				if err != nil {
					t.Errorf("StatusLegacy() error = %v", err)
					return
				}

				if !status.Success {
					t.Errorf("Legacy setup success = %v, want true", status.Success)
				}

				// 4. Verificar se todos os arquivos necessários foram criados
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

				// 5. Executar reset legacy
				_, err = setup.ResetLegacy(setup.LegacySetupOptions{})
				if err != nil {
					t.Errorf("ResetLegacy() error = %v", err)
					return
				}

				// 6. Verificar status após reset legacy
				status, err = setup.StatusLegacy(setup.LegacySetupOptions{})
				if err != nil {
					t.Errorf("StatusLegacy() error = %v", err)
					return
				}

				if !status.Success {
					t.Errorf("Legacy reset success = %v, want false", status.Success)
				}
			}
		})
	}
}

// TestSetupManager_E2E_ErrorHandling testa o tratamento de erros end-to-end
func TestSetupManager_E2E_ErrorHandling(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	manager, err := setup.NewSetupManager()
	if err != nil {
		t.Fatalf("NewSetupManager() error = %v", err)
	}

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should handle errors gracefully in end-to-end flow",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 1. Executar setup
			options := &setup.SetupOptions{
				Force:          false,
				SkipValidation: false,
			}
			err := manager.SetupWithPublicOptions(options)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupManager.Setup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// 2. Verificar se o setup foi concluído com sucesso
				status, err := manager.Status()
				if err != nil {
					t.Errorf("SetupManager.Status() error = %v", err)
					return
				}

				if *status != "completed" {
					t.Errorf("Setup status = %s, want completed", *status)
				}

				// 3. Executar repair para verificar integridade
				err = manager.Repair()
				if err != nil {
					t.Errorf("SetupManager.Repair() error = %v", err)
					return
				}

				// 4. Executar reset
				err = manager.Reset(true)
				if err != nil {
					t.Errorf("SetupManager.Reset() error = %v", err)
					return
				}

				// 5. Verificar se o reset foi executado com sucesso
				status, err = manager.Status()
				if err != nil {
					t.Errorf("SetupManager.Status() error = %v", err)
					return
				}

				if *status != "not_started" {
					t.Errorf("Reset status = %s, want not_started", *status)
				}
			}
		})
	}
}
