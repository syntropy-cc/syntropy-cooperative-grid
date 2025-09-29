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

// TestStateIntegration_StateLoading testa o carregamento de estado
func TestStateIntegration_StateLoading(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(string) error
		wantErr   bool
	}{
		{
			name: "should load state successfully",
			setupFunc: func(stateDir string) error {
				// Criar arquivo de estado
				stateFile := filepath.Join(stateDir, "state.json")
				stateContent := `{
					"current_step": "validation",
					"completed_steps": ["initialization"],
					"errors": [],
					"warnings": [],
					"start_time": "2023-01-01T00:00:00Z",
					"last_updated": "2023-01-01T00:00:00Z",
					"config_path": "/home/testuser/.syntropy/config",
					"keys_path": "/home/testuser/.syntropy/keys",
					"backup_path": "/home/testuser/.syntropy/backups"
				}`
				return os.WriteFile(stateFile, []byte(stateContent), 0644)
			},
			wantErr: false,
		},
		{
			name: "should fail when state file does not exist",
			setupFunc: func(stateDir string) error {
				// Não criar arquivo de estado
				return nil
			},
			wantErr: true,
		},
		{
			name: "should fail when state file is invalid",
			setupFunc: func(stateDir string) error {
				// Criar arquivo de estado inválido
				stateFile := filepath.Join(stateDir, "state.json")
				stateContent := "invalid json content"
				return os.WriteFile(stateFile, []byte(stateContent), 0644)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "state_loading_test")
			stateDir := filepath.Join(tempDir, "state")
			os.MkdirAll(stateDir, 0755)

			// Setup do teste
			if tt.setupFunc != nil {
				if err := tt.setupFunc(stateDir); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			// Criar state manager
			stateManager := src.NewStateManager(nil)
			if stateManager == nil {
				t.Fatal("NewStateManager() returned nil")
			}

			// Executar carregamento de estado
			state, err := stateManager.LoadState()
			if (err != nil) != tt.wantErr {
				t.Errorf("StateManager.LoadState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if state == nil {
					t.Error("LoadState() returned nil state")
					return
				}

				// Verificar se o estado foi carregado corretamente
				helpers.AssertStringEqual(t, state.CurrentStep, "validation", "State CurrentStep")
				helpers.AssertStringEqual(t, state.ConfigPath, "/home/testuser/.syntropy/config", "State ConfigPath")
				helpers.AssertStringEqual(t, state.KeysPath, "/home/testuser/.syntropy/keys", "State KeysPath")
				helpers.AssertStringEqual(t, state.BackupPath, "/home/testuser/.syntropy/backups", "State BackupPath")
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestStateIntegration_StateSaving testa o salvamento de estado
func TestStateIntegration_StateSaving(t *testing.T) {
	tests := []struct {
		name    string
		state   *types.SetupState
		wantErr bool
	}{
		{
			name: "should save state successfully",
			state: &types.SetupState{
				CurrentStep:    "validation",
				CompletedSteps: []string{"initialization"},
				Errors:         []string{},
				Warnings:       []string{},
				StartTime:      time.Now(),
				LastUpdated:    time.Now(),
				ConfigPath:     "/home/testuser/.syntropy/config",
				KeysPath:       "/home/testuser/.syntropy/keys",
				BackupPath:     "/home/testuser/.syntropy/backups",
			},
			wantErr: false,
		},
		{
			name:    "should fail when state is nil",
			state:   nil,
			wantErr: true,
		},
		{
			name: "should handle empty state",
			state: &types.SetupState{
				CurrentStep:    "",
				CompletedSteps: []string{},
				Errors:         []string{},
				Warnings:       []string{},
				StartTime:      time.Now(),
				LastUpdated:    time.Now(),
				ConfigPath:     "",
				KeysPath:       "",
				BackupPath:     "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "state_saving_test")
			stateDir := filepath.Join(tempDir, "state")
			os.MkdirAll(stateDir, 0755)

			// Criar state manager
			stateManager := src.NewStateManager(nil)
			if stateManager == nil {
				t.Fatal("NewStateManager() returned nil")
			}

			// Executar salvamento de estado
			err := stateManager.SaveState(tt.state)
			if (err != nil) != tt.wantErr {
				t.Errorf("StateManager.SaveState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verificar se o arquivo de estado foi criado
				stateFile := filepath.Join(stateDir, "state.json")
				helpers.AssertFileExists(t, stateFile, "State file")
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestStateIntegration_StateUpdating testa a atualização de estado
func TestStateIntegration_StateUpdating(t *testing.T) {
	tests := []struct {
		name      string
		updates   map[string]interface{}
		setupFunc func(string) error
		wantErr   bool
	}{
		{
			name: "should update state successfully",
			updates: map[string]interface{}{
				"current_step":    "configuration",
				"completed_steps": []string{"initialization", "validation"},
			},
			setupFunc: func(stateDir string) error {
				// Criar arquivo de estado inicial
				stateFile := filepath.Join(stateDir, "state.json")
				stateContent := `{
					"current_step": "validation",
					"completed_steps": ["initialization"],
					"errors": [],
					"warnings": [],
					"start_time": "2023-01-01T00:00:00Z",
					"last_updated": "2023-01-01T00:00:00Z",
					"config_path": "/home/testuser/.syntropy/config",
					"keys_path": "/home/testuser/.syntropy/keys",
					"backup_path": "/home/testuser/.syntropy/backups"
				}`
				return os.WriteFile(stateFile, []byte(stateContent), 0644)
			},
			wantErr: false,
		},
		{
			name: "should fail when state file does not exist",
			updates: map[string]interface{}{
				"current_step": "configuration",
			},
			setupFunc: func(stateDir string) error {
				// Não criar arquivo de estado
				return nil
			},
			wantErr: true,
		},
		{
			name:    "should handle empty updates",
			updates: map[string]interface{}{},
			setupFunc: func(stateDir string) error {
				// Criar arquivo de estado inicial
				stateFile := filepath.Join(stateDir, "state.json")
				stateContent := `{
					"current_step": "validation",
					"completed_steps": ["initialization"],
					"errors": [],
					"warnings": [],
					"start_time": "2023-01-01T00:00:00Z",
					"last_updated": "2023-01-01T00:00:00Z",
					"config_path": "/home/testuser/.syntropy/config",
					"keys_path": "/home/testuser/.syntropy/keys",
					"backup_path": "/home/testuser/.syntropy/backups"
				}`
				return os.WriteFile(stateFile, []byte(stateContent), 0644)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "state_updating_test")
			stateDir := filepath.Join(tempDir, "state")
			os.MkdirAll(stateDir, 0755)

			// Setup do teste
			if tt.setupFunc != nil {
				if err := tt.setupFunc(stateDir); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			// Criar state manager
			stateManager := src.NewStateManager(nil)
			if stateManager == nil {
				t.Fatal("NewStateManager() returned nil")
			}

			// Executar atualização de estado
			err := stateManager.UpdateState(tt.updates)
			if (err != nil) != tt.wantErr {
				t.Errorf("StateManager.UpdateState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verificar se o arquivo de estado foi atualizado
				stateFile := filepath.Join(stateDir, "state.json")
				helpers.AssertFileExists(t, stateFile, "Updated state file")
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestStateIntegration_StateBackup testa o backup de estado
func TestStateIntegration_StateBackup(t *testing.T) {
	tests := []struct {
		name       string
		backupPath string
		setupFunc  func(string) error
		wantErr    bool
	}{
		{
			name:       "should backup state successfully",
			backupPath: "state_backup.json",
			setupFunc: func(stateDir string) error {
				// Criar arquivo de estado
				stateFile := filepath.Join(stateDir, "state.json")
				stateContent := `{
					"current_step": "validation",
					"completed_steps": ["initialization"],
					"errors": [],
					"warnings": [],
					"start_time": "2023-01-01T00:00:00Z",
					"last_updated": "2023-01-01T00:00:00Z",
					"config_path": "/home/testuser/.syntropy/config",
					"keys_path": "/home/testuser/.syntropy/keys",
					"backup_path": "/home/testuser/.syntropy/backups"
				}`
				return os.WriteFile(stateFile, []byte(stateContent), 0644)
			},
			wantErr: false,
		},
		{
			name:       "should handle empty backup path",
			backupPath: "",
			setupFunc: func(stateDir string) error {
				// Criar arquivo de estado
				stateFile := filepath.Join(stateDir, "state.json")
				stateContent := `{
					"current_step": "validation",
					"completed_steps": ["initialization"],
					"errors": [],
					"warnings": [],
					"start_time": "2023-01-01T00:00:00Z",
					"last_updated": "2023-01-01T00:00:00Z",
					"config_path": "/home/testuser/.syntropy/config",
					"keys_path": "/home/testuser/.syntropy/keys",
					"backup_path": "/home/testuser/.syntropy/backups"
				}`
				return os.WriteFile(stateFile, []byte(stateContent), 0644)
			},
			wantErr: false,
		},
		{
			name:       "should fail when state file does not exist",
			backupPath: "state_backup.json",
			setupFunc: func(stateDir string) error {
				// Não criar arquivo de estado
				return nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "state_backup_test")
			stateDir := filepath.Join(tempDir, "state")
			backupDir := filepath.Join(tempDir, "backups")
			os.MkdirAll(stateDir, 0755)
			os.MkdirAll(backupDir, 0755)

			// Setup do teste
			if tt.setupFunc != nil {
				if err := tt.setupFunc(stateDir); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			// Criar state manager
			stateManager := src.NewStateManager(nil)
			if stateManager == nil {
				t.Fatal("NewStateManager() returned nil")
			}

			// Executar backup de estado
			fullBackupPath := filepath.Join(backupDir, tt.backupPath)
			err := stateManager.BackupState(fullBackupPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("StateManager.BackupState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verificar se o backup foi criado
				helpers.AssertFileExists(t, fullBackupPath, "State backup file")
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestStateIntegration_StateRestore testa a restauração de estado
func TestStateIntegration_StateRestore(t *testing.T) {
	tests := []struct {
		name       string
		backupPath string
		setupFunc  func(string) error
		wantErr    bool
	}{
		{
			name:       "should restore state successfully",
			backupPath: "state_backup.json",
			setupFunc: func(backupDir string) error {
				// Criar arquivo de backup
				backupFile := filepath.Join(backupDir, "state_backup.json")
				backupContent := `{
					"current_step": "validation",
					"completed_steps": ["initialization"],
					"errors": [],
					"warnings": [],
					"start_time": "2023-01-01T00:00:00Z",
					"last_updated": "2023-01-01T00:00:00Z",
					"config_path": "/home/testuser/.syntropy/config",
					"keys_path": "/home/testuser/.syntropy/keys",
					"backup_path": "/home/testuser/.syntropy/backups"
				}`
				return os.WriteFile(backupFile, []byte(backupContent), 0644)
			},
			wantErr: false,
		},
		{
			name:       "should fail when backup file does not exist",
			backupPath: "nonexistent_backup.json",
			setupFunc: func(backupDir string) error {
				// Não criar arquivo de backup
				return nil
			},
			wantErr: true,
		},
		{
			name:       "should fail when backup file is invalid",
			backupPath: "invalid_backup.json",
			setupFunc: func(backupDir string) error {
				// Criar arquivo de backup inválido
				backupFile := filepath.Join(backupDir, "invalid_backup.json")
				backupContent := "invalid json content"
				return os.WriteFile(backupFile, []byte(backupContent), 0644)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "state_restore_test")
			stateDir := filepath.Join(tempDir, "state")
			backupDir := filepath.Join(tempDir, "backups")
			os.MkdirAll(stateDir, 0755)
			os.MkdirAll(backupDir, 0755)

			// Setup do teste
			if tt.setupFunc != nil {
				if err := tt.setupFunc(backupDir); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			// Criar state manager
			stateManager := src.NewStateManager(nil)
			if stateManager == nil {
				t.Fatal("NewStateManager() returned nil")
			}

			// Executar restauração de estado
			fullBackupPath := filepath.Join(backupDir, tt.backupPath)
			err := stateManager.RestoreState(fullBackupPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("StateManager.RestoreState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verificar se o estado foi restaurado
				stateFile := filepath.Join(stateDir, "state.json")
				helpers.AssertFileExists(t, stateFile, "Restored state file")
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestStateIntegration_StateIntegrity testa a verificação de integridade de estado
func TestStateIntegration_StateIntegrity(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(string) error
		wantErr   bool
	}{
		{
			name: "should verify state integrity successfully",
			setupFunc: func(stateDir string) error {
				// Criar arquivo de estado válido
				stateFile := filepath.Join(stateDir, "state.json")
				stateContent := `{
					"current_step": "validation",
					"completed_steps": ["initialization"],
					"errors": [],
					"warnings": [],
					"start_time": "2023-01-01T00:00:00Z",
					"last_updated": "2023-01-01T00:00:00Z",
					"config_path": "/home/testuser/.syntropy/config",
					"keys_path": "/home/testuser/.syntropy/keys",
					"backup_path": "/home/testuser/.syntropy/backups"
				}`
				return os.WriteFile(stateFile, []byte(stateContent), 0644)
			},
			wantErr: false,
		},
		{
			name: "should fail when state file does not exist",
			setupFunc: func(stateDir string) error {
				// Não criar arquivo de estado
				return nil
			},
			wantErr: true,
		},
		{
			name: "should fail when state file is corrupted",
			setupFunc: func(stateDir string) error {
				// Criar arquivo de estado corrompido
				stateFile := filepath.Join(stateDir, "state.json")
				stateContent := "corrupted content"
				return os.WriteFile(stateFile, []byte(stateContent), 0644)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "state_integrity_test")
			stateDir := filepath.Join(tempDir, "state")
			os.MkdirAll(stateDir, 0755)

			// Setup do teste
			if tt.setupFunc != nil {
				if err := tt.setupFunc(stateDir); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			// Criar state manager
			stateManager := src.NewStateManager(nil)
			if stateManager == nil {
				t.Fatal("NewStateManager() returned nil")
			}

			// Executar verificação de integridade
			err := stateManager.VerifyIntegrity()
			if (err != nil) != tt.wantErr {
				t.Errorf("StateManager.VerifyIntegrity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestStateIntegration_ErrorHandling testa o tratamento de erros
func TestStateIntegration_ErrorHandling(t *testing.T) {
	tests := []struct {
		name    string
		state   *types.SetupState
		wantErr bool
	}{
		{
			name:    "should handle nil state",
			state:   nil,
			wantErr: true,
		},
		{
			name: "should handle invalid state data",
			state: &types.SetupState{
				CurrentStep:    "invalid-step",
				CompletedSteps: []string{"invalid-step"},
				Errors:         []string{"invalid error"},
				Warnings:       []string{"invalid warning"},
				StartTime:      time.Now(),
				LastUpdated:    time.Now(),
				ConfigPath:     "/invalid/path",
				KeysPath:       "/invalid/path",
				BackupPath:     "/invalid/path",
			},
			wantErr: false, // Should not fail, just use invalid data
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "state_error_test")
			stateDir := filepath.Join(tempDir, "state")
			os.MkdirAll(stateDir, 0755)

			// Criar state manager
			stateManager := src.NewStateManager(nil)
			if stateManager == nil {
				t.Fatal("NewStateManager() returned nil")
			}

			// Executar operação que pode falhar
			err := stateManager.SaveState(tt.state)
			if (err != nil) != tt.wantErr {
				t.Errorf("StateManager.SaveState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestStateIntegration_Concurrency testa concorrência
func TestStateIntegration_Concurrency(t *testing.T) {
	t.Run("should handle concurrent state operations", func(t *testing.T) {
		// Criar diretório temporário para testes
		tempDir := helpers.CreateTempDir(t, "state_concurrent_test")
		stateDir := filepath.Join(tempDir, "state")
		os.MkdirAll(stateDir, 0755)

		// Criar state manager
		stateManager := src.NewStateManager(nil)
		if stateManager == nil {
			t.Fatal("NewStateManager() returned nil")
		}

		// Executar múltiplas operações de estado concorrentemente
		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func(instance int) {
				// Criar estado
				state := &types.SetupState{
					CurrentStep:    "step_" + string(rune(instance)),
					CompletedSteps: []string{"initialization"},
					Errors:         []string{},
					Warnings:       []string{},
					StartTime:      time.Now(),
					LastUpdated:    time.Now(),
					ConfigPath:     "/tmp/syntropy_test/config",
					KeysPath:       "/tmp/syntropy_test/keys",
					BackupPath:     "/tmp/syntropy_test/backups",
				}

				// Salvar estado
				err := stateManager.SaveState(state)
				if err != nil {
					t.Errorf("Concurrent SaveState() failed: %v", err)
				}

				// Carregar estado
				loadedState, err := stateManager.LoadState()
				if err != nil {
					t.Errorf("Concurrent LoadState() failed: %v", err)
				}
				if loadedState == nil {
					t.Error("Concurrent LoadState() returned nil state")
				}

				// Atualizar estado
				updates := map[string]interface{}{
					"current_step": "updated_step_" + string(rune(instance)),
				}
				err = stateManager.UpdateState(updates)
				if err != nil {
					t.Errorf("Concurrent UpdateState() failed: %v", err)
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

// TestStateIntegration_Performance testa performance
func TestStateIntegration_Performance(t *testing.T) {
	t.Run("should complete state operations within reasonable time", func(t *testing.T) {
		// Criar diretório temporário para testes
		tempDir := helpers.CreateTempDir(t, "state_perf_test")
		stateDir := filepath.Join(tempDir, "state")
		os.MkdirAll(stateDir, 0755)

		// Criar state manager
		stateManager := src.NewStateManager(nil)
		if stateManager == nil {
			t.Fatal("NewStateManager() returned nil")
		}

		// Criar estado
		state := &types.SetupState{
			CurrentStep:    "validation",
			CompletedSteps: []string{"initialization"},
			Errors:         []string{},
			Warnings:       []string{},
			StartTime:      time.Now(),
			LastUpdated:    time.Now(),
			ConfigPath:     "/tmp/syntropy_test/config",
			KeysPath:       "/tmp/syntropy_test/keys",
			BackupPath:     "/tmp/syntropy_test/backups",
		}

		start := time.Now()
		err := stateManager.SaveState(state)
		elapsed := time.Since(start)

		if err != nil {
			t.Errorf("SaveState() failed: %v", err)
		}

		if elapsed > 5*time.Second {
			t.Errorf("SaveState() took too long: %v", elapsed)
		}

		// Limpar diretório temporário
		os.RemoveAll(tempDir)
	})
}
