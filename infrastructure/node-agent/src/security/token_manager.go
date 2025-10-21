package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/crypto/pbkdf2"
)

type SecureTokenManager struct {
	tokenPath string
}

// NewSecureTokenManager creates a new secure token manager
func NewSecureTokenManager() *SecureTokenManager {
	return &SecureTokenManager{
		tokenPath: "/opt/syntropy/config/grid_token",
	}
}

// DecryptAndStoreToken decrypts token from cloud-init and stores securely
func (stm *SecureTokenManager) DecryptAndStoreToken(encryptedToken string, nodeID string) error {
	// 1. Derivar chave do NodeID (mesma lógica do cloud-init generator)
	nodeSeed := sha256.Sum256([]byte(nodeID))
	nodeKey := pbkdf2.Key(nodeSeed[:], []byte("syntropy-grid"), 50000, 32, sha256.New)

	// 2. Descriptografar
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedToken)
	if err != nil {
		return fmt.Errorf("failed to decode encrypted token: %w", err)
	}

	block, err := aes.NewCipher(nodeKey)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("failed to decrypt token: %w", err)
	}

	// 3. Armazenar com permissões restritivas (0400 - apenas root lê)
	if err := os.MkdirAll(filepath.Dir(stm.tokenPath), 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(stm.tokenPath, plaintext, 0400); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	// 4. Limpar variáveis da memória
	for i := range plaintext {
		plaintext[i] = 0
	}
	for i := range nodeKey {
		nodeKey[i] = 0
	}

	return nil
}

// LoadToken loads the stored token
func (stm *SecureTokenManager) LoadToken() (string, error) {
	tokenBytes, err := os.ReadFile(stm.tokenPath)
	if err != nil {
		return "", fmt.Errorf("failed to read token file: %w", err)
	}

	token := string(tokenBytes)

	// Clear memory after use
	for i := range tokenBytes {
		tokenBytes[i] = 0
	}

	return token, nil
}

// DeleteToken removes the stored token
func (stm *SecureTokenManager) DeleteToken() error {
	return os.Remove(stm.tokenPath)
}

// TokenExists checks if token file exists
func (stm *SecureTokenManager) TokenExists() bool {
	_, err := os.Stat(stm.tokenPath)
	return !os.IsNotExist(err)
}
