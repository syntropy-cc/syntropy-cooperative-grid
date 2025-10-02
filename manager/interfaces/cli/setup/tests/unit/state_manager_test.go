//go:build !integration && !e2e && !performance && !security
// +build !integration,!e2e,!performance,!security

package unit

import (
	"os"
	"path/filepath"
	"testing"

	setup "setup-component/src"
)

// TestNewStateManager testa a criação do gerenciador de estado
func TestNewStateManager(t *testing.T) {
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
			name:    "should create state manager successfully",
			logger:  logger,
			wantErr: false,
		},
		{
			name:    "should create state manager with nil logger",
			logger:  nil,
			wantErr: false, // Logger pode ser nil
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stateManager := setup.NewStateManager(tt.logger)
			if stateManager == nil {
				t.Error("NewStateManager() returned nil state manager")
			}
		})
	}
}

// TestStateManager_LoadState testa o carregamento de estado
func TestStateManager_LoadState(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	stateManager := setup.NewStateManager(logger)

	tests := []struct {
		name    string
		setup   bool
		wantErr bool
	}{
		{
			name:    "should fail load when state file does not exist",
			setup:   false,
			wantErr: true,
		},
		{
			name:    "should load state successfully when state file exists",
			setup:   true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar arquivo de estado se necessário
			if tt.setup {
				state := &types.setup.SetupState{
					Status:  "completed",
					Version: "1.0.0",
					Config:  map[string]interface{}{"test": "value"},
				}
				err := stateManager.SaveState(state)
				if err != nil {
					t.Fatalf("Failed to save state: %v", err)
				}
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

				// Verificar campos obrigatórios
				if state.Status == "" {
					t.Error("State missing status")
				}
				if state.Version == "" {
					t.Error("State missing version")
				}
				// Verificar se o estado foi carregado corretamente
				if state.Status == "" {
					t.Error("State missing status")
				}
			}
		})
	}
}

// TestStateManager_SaveState testa o salvamento de estado
func TestStateManager_SaveState(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	stateManager := setup.NewStateManager(logger)

	tests := []struct {
		name    string
		state   *setup.SetupState
		wantErr bool
	}{
		{
			name:    "should fail save with nil state",
			state:   nil,
			wantErr: true,
		},
		{
			name: "should save state successfully",
			state: &types.setup.SetupState{
				Status:    "completed",
				Version:   "1.0.0",
				Timestamp: "2023-01-01T00:00:00Z",
				Config:    map[string]interface{}{"test": "value"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := stateManager.SaveState(tt.state)
			if (err != nil) != tt.wantErr {
				t.Errorf("StateManager.SaveState() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verificar se o arquivo foi salvo
			if !tt.wantErr {
				statePath := filepath.Join(tempDir, ".syntropy", "state", "setup_state.json")
				if _, err := os.Stat(statePath); os.IsNotExist(err) {
					t.Errorf("State file not saved: %s", statePath)
				}
			}
		})
	}
}

// TestStateManager_UpdateState testa a atualização de estado
func TestStateManager_UpdateState(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	stateManager := setup.NewStateManager(logger)

	tests := []struct {
		name    string
		setup   bool
		wantErr bool
	}{
		{
			name:    "should fail update when no state exists",
			setup:   false,
			wantErr: true,
		},
		{
			name:    "should update state successfully when state exists",
			setup:   true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar estado se necessário
			if tt.setup {
				state := &types.setup.SetupState{
					Status:  "in_progress",
					Version: "1.0.0",
					Config:  map[string]interface{}{"test": "value"},
				}
				err := stateManager.SaveState(state)
				if err != nil {
					t.Fatalf("Failed to save state: %v", err)
				}
			}

			updates := func(state *setup.SetupState) error {
				state.Status = "completed"
				return nil
			}

			err := stateManager.UpdateState(updates)
			if (err != nil) != tt.wantErr {
				t.Errorf("StateManager.UpdateState() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verificar se o estado foi atualizado
			if !tt.wantErr {
				state, err := stateManager.LoadState()
				if err != nil {
					t.Errorf("Failed to load updated state: %v", err)
					return
				}

				if state.Status != "completed" {
					t.Errorf("State status not updated: got %s, want completed", state.Status)
				}
			}
		})
	}
}

// TestStateManager_GetState testa a obtenção de estado
func TestStateManager_GetState(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	stateManager := setup.NewStateManager(logger)

	tests := []struct {
		name    string
		setup   bool
		wantErr bool
	}{
		{
			name:    "should fail get when no state exists",
			setup:   false,
			wantErr: true,
		},
		{
			name:    "should get state successfully when state exists",
			setup:   true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar estado se necessário
			if tt.setup {
				state := &types.setup.SetupState{
					Status:  "completed",
					Version: "1.0.0",
					Config:  map[string]interface{}{"test": "value"},
				}
				err := stateManager.SaveState(state)
				if err != nil {
					t.Fatalf("Failed to save state: %v", err)
				}
			}

			state, err := stateManager.LoadState()
			if (err != nil) != tt.wantErr {
				t.Errorf("StateManager.GetState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if state == nil {
					t.Error("StateManager.GetState() returned nil state")
					return
				}

				// Verificar campos obrigatórios
				if state.Status == "" {
					t.Error("State missing status")
				}
				if state.Version == "" {
					t.Error("State missing version")
				}
				// Verificar se o estado foi carregado corretamente
				if state.Status == "" {
					t.Error("State missing status")
				}
			}
		})
	}
}

// TestStateManager_ResetState testa o reset de estado
func TestStateManager_ResetState(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	stateManager := setup.NewStateManager(logger)

	tests := []struct {
		name    string
		setup   bool
		wantErr bool
	}{
		{
			name:    "should reset state successfully when no state exists",
			setup:   false,
			wantErr: false,
		},
		{
			name:    "should reset state successfully when state exists",
			setup:   true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar estado se necessário
			if tt.setup {
				state := &types.setup.SetupState{
					Status:  "completed",
					Version: "1.0.0",
					Config:  map[string]interface{}{"test": "value"},
				}
				err := stateManager.SaveState(state)
				if err != nil {
					t.Fatalf("Failed to save state: %v", err)
				}
			}

			err := stateManager.Reset()
			if (err != nil) != tt.wantErr {
				t.Errorf("StateManager.ResetState() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verificar se o estado foi resetado
			if !tt.wantErr {
				state, err := stateManager.LoadState()
				if err != nil {
					t.Errorf("Failed to get reset state: %v", err)
					return
				}

				if state.Status != "not_started" {
					t.Errorf("State status not reset: got %s, want not_started", state.Status)
				}
			}
		})
	}
}

// TestStateManager_BackupState testa o backup de estado
func TestStateManager_BackupState(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	stateManager := setup.NewStateManager(logger)

	tests := []struct {
		name    string
		setup   bool
		wantErr bool
	}{
		{
			name:    "should fail backup when no state exists",
			setup:   false,
			wantErr: true,
		},
		{
			name:    "should backup state successfully when state exists",
			setup:   true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar estado se necessário
			if tt.setup {
				state := &types.setup.SetupState{
					Status:  "completed",
					Version: "1.0.0",
					Config:  map[string]interface{}{"test": "value"},
				}
				err := stateManager.SaveState(state)
				if err != nil {
					t.Fatalf("Failed to save state: %v", err)
				}
			}

			backupPath, err := stateManager.BackupState("test_backup")
			if (err != nil) != tt.wantErr {
				t.Errorf("StateManager.BackupState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if backupPath == "" {
					t.Error("StateManager.BackupState() returned empty backup path")
				}

				// Verificar se o backup foi criado
				if _, err := os.Stat(backupPath); os.IsNotExist(err) {
					t.Errorf("Backup file not created: %s", backupPath)
				}
			}
		})
	}
}

// TestStateManager_RestoreState testa a restauração de estado
func TestStateManager_RestoreState(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	stateManager := setup.NewStateManager(logger)

	tests := []struct {
		name    string
		setup   bool
		wantErr bool
	}{
		{
			name:    "should fail restore when backup does not exist",
			setup:   false,
			wantErr: true,
		},
		{
			name:    "should restore state successfully when backup exists",
			setup:   true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var backupPath string

			// Criar backup se necessário
			if tt.setup {
				state := &types.setup.SetupState{
					Status:  "completed",
					Version: "1.0.0",
					Config:  map[string]interface{}{"test": "value"},
				}
				err := stateManager.SaveState(state)
				if err != nil {
					t.Fatalf("Failed to save state: %v", err)
				}

				backupPath, err = stateManager.BackupState("test_backup")
				if err != nil {
					t.Fatalf("Failed to backup state: %v", err)
				}
			} else {
				backupPath = filepath.Join(tempDir, "nonexistent_backup.json")
			}

			err := stateManager.RestoreState(backupPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("StateManager.RestoreState() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestStateManager_VerifyStateIntegrity testa a verificação de integridade de estado
func TestStateManager_VerifyStateIntegrity(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	stateManager := setup.NewStateManager(logger)

	tests := []struct {
		name    string
		setup   bool
		wantErr bool
	}{
		{
			name:    "should fail verify when no state exists",
			setup:   false,
			wantErr: true,
		},
		{
			name:    "should verify state integrity successfully when state exists",
			setup:   true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar estado se necessário
			if tt.setup {
				state := &types.setup.SetupState{
					Status:  "completed",
					Version: "1.0.0",
					Config:  map[string]interface{}{"test": "value"},
				}
				err := stateManager.SaveState(state)
				if err != nil {
					t.Fatalf("Failed to save state: %v", err)
				}
			}

			valid, err := stateManager.VerifyIntegrity()
			if (err != nil) != tt.wantErr {
				t.Errorf("StateManager.VerifyStateIntegrity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !valid {
					t.Error("StateManager.VerifyStateIntegrity() returned false for valid state")
				}
			}
		})
	}
}

// TestStateManager_GetStatePath testa a obtenção do caminho de estado
func TestStateManager_GetStatePath(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	stateManager := setup.NewStateManager(logger)

	tests := []struct {
		name string
		want string
	}{
		{
			name: "should return state path",
			want: filepath.Join(tempDir, ".syntropy", "state", "setup_state.json"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stateManager.LoadStatePath()
			if result != tt.want {
				t.Errorf("StateManager.GetStatePath() = %v, want %v", result, tt.want)
			}
		})
	}
}

// TestStateManager_GetStateDir testa a obtenção do diretório de estado
func TestStateManager_GetStateDir(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	stateManager := setup.NewStateManager(logger)

	tests := []struct {
		name string
		want string
	}{
		{
			name: "should return state directory path",
			want: filepath.Join(tempDir, ".syntropy", "state"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stateManager.GetStatePath()
			if result != tt.want {
				t.Errorf("StateManager.GetStatePath() = %v, want %v", result, tt.want)
			}
		})
	}
}
