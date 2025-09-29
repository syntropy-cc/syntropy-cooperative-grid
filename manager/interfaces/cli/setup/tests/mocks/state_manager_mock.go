package mocks

import (
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src/internal/types"
)

// MockStateManager implementa a interface StateManager para testes
type MockStateManager struct {
	LoadStateFunc       func() (*types.SetupState, error)
	SaveStateFunc       func(state *types.SetupState) error
	UpdateStateFunc     func(update func(*types.SetupState) error) error
	BackupStateFunc     func(name string) error
	RestoreStateFunc    func(backupPath string) error
	VerifyIntegrityFunc func() error
}

// LoadState chama a função mock
func (m *MockStateManager) LoadState() (*types.SetupState, error) {
	if m.LoadStateFunc != nil {
		return m.LoadStateFunc()
	}
	return &types.SetupState{
		Version:       "1.0.0",
		Status:        types.SetupStatusNotStarted,
		Environment:   &types.EnvironmentInfo{},
		Configuration: &types.ConfigInfo{},
		Keys:          &types.KeyInfo{},
		Metadata:      make(map[string]string),
	}, nil
}

// SaveState chama a função mock
func (m *MockStateManager) SaveState(state *types.SetupState) error {
	if m.SaveStateFunc != nil {
		return m.SaveStateFunc(state)
	}
	return nil
}

// UpdateState chama a função mock
func (m *MockStateManager) UpdateState(update func(*types.SetupState) error) error {
	if m.UpdateStateFunc != nil {
		return m.UpdateStateFunc(update)
	}
	return nil
}

// BackupState chama a função mock
func (m *MockStateManager) BackupState(name string) error {
	if m.BackupStateFunc != nil {
		return m.BackupStateFunc(name)
	}
	return nil
}

// RestoreState chama a função mock
func (m *MockStateManager) RestoreState(backupPath string) error {
	if m.RestoreStateFunc != nil {
		return m.RestoreStateFunc(backupPath)
	}
	return nil
}

// VerifyIntegrity chama a função mock
func (m *MockStateManager) VerifyIntegrity() error {
	if m.VerifyIntegrityFunc != nil {
		return m.VerifyIntegrityFunc()
	}
	return nil
}
