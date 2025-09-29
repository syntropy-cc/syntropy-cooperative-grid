package mocks

import (
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src/internal/types"
)

// MockKeyManager implementa a interface KeyManager para testes
type MockKeyManager struct {
	GenerateKeyPairFunc    func(algorithm string) (*types.KeyPair, error)
	StoreKeyPairFunc       func(keyPair *types.KeyPair, passphrase string) error
	LoadKeyPairFunc        func(keyID string, passphrase string) (*types.KeyPair, error)
	RotateKeysFunc         func(keyID string) error
	VerifyKeyIntegrityFunc func(keyID string) error
	BackupKeysFunc         func(keyID string, passphrase string) ([]byte, error)
	RestoreKeysFunc        func(backupData []byte, passphrase string) error
}

// GenerateKeyPair chama a função mock
func (m *MockKeyManager) GenerateKeyPair(algorithm string) (*types.KeyPair, error) {
	if m.GenerateKeyPairFunc != nil {
		return m.GenerateKeyPairFunc(algorithm)
	}
	return &types.KeyPair{
		ID:          "mock_key_12345",
		Algorithm:   algorithm,
		PrivateKey:  []byte("mock_private_key"),
		PublicKey:   []byte("mock_public_key"),
		Fingerprint: "mock_fingerprint",
		Metadata: map[string]string{
			"generated_by": "mock",
			"version":      "1.0.0",
		},
	}, nil
}

// StoreKeyPair chama a função mock
func (m *MockKeyManager) StoreKeyPair(keyPair *types.KeyPair, passphrase string) error {
	if m.StoreKeyPairFunc != nil {
		return m.StoreKeyPairFunc(keyPair, passphrase)
	}
	return nil
}

// LoadKeyPair chama a função mock
func (m *MockKeyManager) LoadKeyPair(keyID string, passphrase string) (*types.KeyPair, error) {
	if m.LoadKeyPairFunc != nil {
		return m.LoadKeyPairFunc(keyID, passphrase)
	}
	return &types.KeyPair{
		ID:          keyID,
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

// RotateKeys chama a função mock
func (m *MockKeyManager) RotateKeys(keyID string) error {
	if m.RotateKeysFunc != nil {
		return m.RotateKeysFunc(keyID)
	}
	return nil
}

// VerifyKeyIntegrity chama a função mock
func (m *MockKeyManager) VerifyKeyIntegrity(keyID string) error {
	if m.VerifyKeyIntegrityFunc != nil {
		return m.VerifyKeyIntegrityFunc(keyID)
	}
	return nil
}

// BackupKeys chama a função mock
func (m *MockKeyManager) BackupKeys(keyID string, passphrase string) ([]byte, error) {
	if m.BackupKeysFunc != nil {
		return m.BackupKeysFunc(keyID, passphrase)
	}
	return []byte("mock_backup_data"), nil
}

// RestoreKeys chama a função mock
func (m *MockKeyManager) RestoreKeys(backupData []byte, passphrase string) error {
	if m.RestoreKeysFunc != nil {
		return m.RestoreKeysFunc(backupData, passphrase)
	}
	return nil
}
