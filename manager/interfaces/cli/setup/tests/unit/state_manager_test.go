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

// TestNewStateManager testa a criação de um novo StateManager
func TestNewStateManager(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "should create state manager successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &mocks.MockSetupLogger{}
			stateManager := src.NewStateManager(logger)
			if stateManager == nil {
				t.Error("NewStateManager() returned nil state manager")
			}
		})
	}
}

// TestStateManager_LoadState testa o método LoadState
func TestStateManager_LoadState(t *testing.T) {
	tests := []struct {
		name       string
		setupFunc  func(string) error
		mockLogger *mocks.MockSetupLogger
		wantErr    bool
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
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    false,
		},
		{
			name: "should fail when state file does not exist",
			setupFunc: func(stateDir string) error {
				// Não criar arquivo de estado
				return nil
			},
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    true,
		},
		{
			name: "should fail when state file is invalid",
			setupFunc: func(stateDir string) error {
				// Criar arquivo de estado inválido
				stateFile := filepath.Join(stateDir, "state.json")
				stateContent := "invalid json content"
				return os.WriteFile(stateFile, []byte(stateContent), 0644)
			},
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "state_manager_load_test")
			stateDir := filepath.Join(tempDir, "state")
			os.MkdirAll(stateDir, 0755)

			// Setup do teste
			if tt.setupFunc != nil {
				if err := tt.setupFunc(stateDir); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			// Mock do logger
			logger := &mocks.MockSetupLogger{}

			// Criar state manager com diretório temporário
			stateManager := &src.StateManager{
				Logger: logger,
			}

			state, err := stateManager.LoadState()
			if (err != nil) != tt.wantErr {
				t.Errorf("StateManager.LoadState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if state == nil {
					t.Error("StateManager.LoadState() returned nil state")
					return
				}
				helpers.AssertStringEqual(t, state.CurrentStep, "validation", "State CurrentStep")
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestStateManager_SaveState testa o método SaveState
func TestStateManager_SaveState(t *testing.T) {
	tests := []struct {
		name       string
		state      *types.SetupState
		mockLogger *mocks.MockSetupLogger
		wantErr    bool
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
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    false,
		},
		{
			name:       "should fail when state is nil",
			state:      nil,
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    true,
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
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "state_manager_save_test")
			stateDir := filepath.Join(tempDir, "state")
			os.MkdirAll(stateDir, 0755)

			// Mock do logger
			logger := &mocks.MockSetupLogger{}

			// Criar state manager com diretório temporário
			stateManager := &src.StateManager{
				Logger: logger,
			}

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

// TestStateManager_UpdateState testa o método UpdateState
func TestStateManager_UpdateState(t *testing.T) {
	tests := []struct {
		name       string
		updates    map[string]interface{}
		setupFunc  func(string) error
		mockLogger *mocks.MockSetupLogger
		wantErr    bool
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
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    false,
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
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    true,
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
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "state_manager_update_test")
			stateDir := filepath.Join(tempDir, "state")
			os.MkdirAll(stateDir, 0755)

			// Setup do teste
			if tt.setupFunc != nil {
				if err := tt.setupFunc(stateDir); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			// Mock do logger
			logger := &mocks.MockSetupLogger{}

			// Criar state manager com diretório temporário
			stateManager := &src.StateManager{
				Logger: logger,
			}

			err := stateManager.UpdateState(tt.updates)
			if (err != nil) != tt.wantErr {
				t.Errorf("StateManager.UpdateState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestStateManager_BackupState testa o método BackupState
func TestStateManager_BackupState(t *testing.T) {
	tests := []struct {
		name       string
		backupPath string
		setupFunc  func(string) error
		mockLogger *mocks.MockSetupLogger
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
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    false,
		},
		{
			name:       "should fail when state file does not exist",
			backupPath: "state_backup.json",
			setupFunc: func(stateDir string) error {
				// Não criar arquivo de estado
				return nil
			},
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    true,
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
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "state_manager_backup_test")
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

			// Mock do logger
			logger := &mocks.MockSetupLogger{}

			// Criar state manager com diretório temporário
			stateManager := &src.StateManager{
				Logger: logger,
			}

			fullBackupPath := filepath.Join(backupDir, tt.backupPath)
			err := stateManager.BackupState(fullBackupPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("StateManager.BackupState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestStateManager_RestoreState testa o método RestoreState
func TestStateManager_RestoreState(t *testing.T) {
	tests := []struct {
		name       string
		backupPath string
		setupFunc  func(string) error
		mockLogger *mocks.MockSetupLogger
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
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    false,
		},
		{
			name:       "should fail when backup file does not exist",
			backupPath: "nonexistent_backup.json",
			setupFunc: func(backupDir string) error {
				// Não criar arquivo de backup
				return nil
			},
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    true,
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
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "state_manager_restore_test")
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

			// Mock do logger
			logger := &mocks.MockSetupLogger{}

			// Criar state manager com diretório temporário
			stateManager := &src.StateManager{
				Logger: logger,
			}

			fullBackupPath := filepath.Join(backupDir, tt.backupPath)
			err := stateManager.RestoreState(fullBackupPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("StateManager.RestoreState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestStateManager_VerifyIntegrity testa o método VerifyIntegrity
func TestStateManager_VerifyIntegrity(t *testing.T) {
	tests := []struct {
		name       string
		setupFunc  func(string) error
		mockLogger *mocks.MockSetupLogger
		wantErr    bool
	}{
		{
			name: "should verify integrity successfully",
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
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    false,
		},
		{
			name: "should fail when state file does not exist",
			setupFunc: func(stateDir string) error {
				// Não criar arquivo de estado
				return nil
			},
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    true,
		},
		{
			name: "should fail when state file is corrupted",
			setupFunc: func(stateDir string) error {
				// Criar arquivo de estado corrompido
				stateFile := filepath.Join(stateDir, "state.json")
				stateContent := "corrupted content"
				return os.WriteFile(stateFile, []byte(stateContent), 0644)
			},
			mockLogger: &mocks.MockSetupLogger{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "state_manager_verify_test")
			stateDir := filepath.Join(tempDir, "state")
			os.MkdirAll(stateDir, 0755)

			// Setup do teste
			if tt.setupFunc != nil {
				if err := tt.setupFunc(stateDir); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			// Mock do logger
			logger := &mocks.MockSetupLogger{}

			// Criar state manager com diretório temporário
			stateManager := &src.StateManager{
				Logger: logger,
			}

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

// TestStateManager_EdgeCases testa casos extremos do StateManager
func TestStateManager_EdgeCases(t *testing.T) {
	t.Run("should handle nil logger", func(t *testing.T) {
		stateManager := &src.StateManager{
			Logger: nil,
		}

		// Should not panic
		state, err := stateManager.LoadState()
		if err != nil {
			t.Errorf("LoadState() failed with nil logger: %v", err)
		}
		if state == nil {
			t.Error("LoadState() returned nil state with nil logger")
		}
	})

	t.Run("should handle empty updates", func(t *testing.T) {
		tempDir := helpers.CreateTempDir(t, "state_manager_edge_test")
		stateDir := filepath.Join(tempDir, "state")
		os.MkdirAll(stateDir, 0755)

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
		os.WriteFile(stateFile, []byte(stateContent), 0644)

		stateManager := &src.StateManager{
			Logger: &mocks.MockSetupLogger{},
		}

		err := stateManager.UpdateState(map[string]interface{}{})
		if err != nil {
			t.Errorf("UpdateState() failed with empty updates: %v", err)
		}

		os.RemoveAll(tempDir)
	})
}

// TestStateManager_Concurrency testa concorrência do StateManager
func TestStateManager_Concurrency(t *testing.T) {
	t.Run("should handle concurrent state updates", func(t *testing.T) {
		tempDir := helpers.CreateTempDir(t, "state_manager_concurrent_test")
		stateDir := filepath.Join(tempDir, "state")
		os.MkdirAll(stateDir, 0755)

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
		os.WriteFile(stateFile, []byte(stateContent), 0644)

		stateManager := &src.StateManager{
			Logger: &mocks.MockSetupLogger{},
		}

		// Executar múltiplas chamadas de UpdateState concorrentemente
		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func(step int) {
				updates := map[string]interface{}{
					"current_step": "step_" + string(rune(step)),
				}
				err := stateManager.UpdateState(updates)
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

		os.RemoveAll(tempDir)
	})
}

// TestStateManager_Performance testa performance do StateManager
func TestStateManager_Performance(t *testing.T) {
	t.Run("should complete operations within reasonable time", func(t *testing.T) {
		tempDir := helpers.CreateTempDir(t, "state_manager_perf_test")
		stateDir := filepath.Join(tempDir, "state")
		os.MkdirAll(stateDir, 0755)

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
		os.WriteFile(stateFile, []byte(stateContent), 0644)

		stateManager := &src.StateManager{
			Logger: &mocks.MockSetupLogger{},
		}

		start := time.Now()
		state, err := stateManager.LoadState()
		elapsed := time.Since(start)

		if err != nil {
			t.Errorf("LoadState() failed: %v", err)
		}
		if state == nil {
			t.Error("LoadState() returned nil state")
		}

		if elapsed > 1*time.Second {
			t.Errorf("LoadState() took too long: %v", elapsed)
		}

		os.RemoveAll(tempDir)
	})
}

// TestStateManager_Atomicity testa atomicidade das operações do StateManager
func TestStateManager_Atomicity(t *testing.T) {
	t.Run("should maintain state consistency during concurrent operations", func(t *testing.T) {
		tempDir := helpers.CreateTempDir(t, "state_manager_atomic_test")
		stateDir := filepath.Join(tempDir, "state")
		os.MkdirAll(stateDir, 0755)

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
		os.WriteFile(stateFile, []byte(stateContent), 0644)

		stateManager := &src.StateManager{
			Logger: &mocks.MockSetupLogger{},
		}

		// Executar operações concorrentes
		done := make(chan bool, 3)

		// Operação 1: Atualizar estado
		go func() {
			updates := map[string]interface{}{
				"current_step": "configuration",
			}
			err := stateManager.UpdateState(updates)
			if err != nil {
				t.Errorf("UpdateState() failed: %v", err)
			}
			done <- true
		}()

		// Operação 2: Carregar estado
		go func() {
			state, err := stateManager.LoadState()
			if err != nil {
				t.Errorf("LoadState() failed: %v", err)
			}
			if state == nil {
				t.Error("LoadState() returned nil state")
			}
			done <- true
		}()

		// Operação 3: Verificar integridade
		go func() {
			err := stateManager.VerifyIntegrity()
			if err != nil {
				t.Errorf("VerifyIntegrity() failed: %v", err)
			}
			done <- true
		}()

		// Aguardar todas as operações terminarem
		for i := 0; i < 3; i++ {
			<-done
		}

		// Verificar que o estado final é consistente
		finalState, err := stateManager.LoadState()
		if err != nil {
			t.Errorf("Final LoadState() failed: %v", err)
		}
		if finalState == nil {
			t.Error("Final LoadState() returned nil state")
		}

		os.RemoveAll(tempDir)
	})
}
