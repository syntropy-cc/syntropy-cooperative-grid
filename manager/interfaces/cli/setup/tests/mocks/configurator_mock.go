package mocks

import (
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src/internal/types"
)

// MockConfigurator implementa a interface Configurator para testes
type MockConfigurator struct {
	GenerateConfigFunc  func(options *types.ConfigOptions) error
	CreateStructureFunc func() error
	GenerateKeysFunc    func() (*types.KeyPair, error)
	ValidateConfigFunc  func() error
	BackupConfigFunc    func(name string) error
	RestoreConfigFunc   func(backupPath string) error
}

// GenerateConfig chama a função mock
func (m *MockConfigurator) GenerateConfig(options *types.ConfigOptions) error {
	if m.GenerateConfigFunc != nil {
		return m.GenerateConfigFunc(options)
	}
	return nil
}

// CreateStructure chama a função mock
func (m *MockConfigurator) CreateStructure() error {
	if m.CreateStructureFunc != nil {
		return m.CreateStructureFunc()
	}
	return nil
}

// GenerateKeys chama a função mock
func (m *MockConfigurator) GenerateKeys() (*types.KeyPair, error) {
	if m.GenerateKeysFunc != nil {
		return m.GenerateKeysFunc()
	}
	return &types.KeyPair{
		ID:          "mock_key_12345",
		Algorithm:   "ed25519",
		PrivateKey:  []byte("mock_private_key"),
		PublicKey:   []byte("mock_public_key"),
		Fingerprint: "mock_fingerprint",
		Metadata: map[string]string{
			"generated_by": "mock",
			"version":      "1.0.0",
		},
	}, nil
}

// ValidateConfig chama a função mock
func (m *MockConfigurator) ValidateConfig() error {
	if m.ValidateConfigFunc != nil {
		return m.ValidateConfigFunc()
	}
	return nil
}

// BackupConfig chama a função mock
func (m *MockConfigurator) BackupConfig(name string) error {
	if m.BackupConfigFunc != nil {
		return m.BackupConfigFunc(name)
	}
	return nil
}

// RestoreConfig chama a função mock
func (m *MockConfigurator) RestoreConfig(backupPath string) error {
	if m.RestoreConfigFunc != nil {
		return m.RestoreConfigFunc(backupPath)
	}
	return nil
}
