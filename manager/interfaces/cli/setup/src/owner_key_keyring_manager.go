package setup

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"setup-component/src/internal/types"

	"github.com/zalando/go-keyring"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/pbkdf2"
)

const (
	KeyringServiceOwnerKey = "syntropy-owner-key"
	KeyringUserPrivateKey  = "owner-private-key"
)

// OwnerKeyKeyringManager manages Owner Keys in system keyring
type OwnerKeyKeyringManager struct {
	logger           types.SetupLogger
	keysDir          string
	keyringAvailable bool
}

func NewOwnerKeyKeyringManager(logger types.SetupLogger) *OwnerKeyKeyringManager {
	keyringAvailable := isKeyringAvailable()
	homeDir, _ := os.UserHomeDir()
	keysDir := filepath.Join(homeDir, ".syntropy", "keys")
	os.MkdirAll(keysDir, 0700)

	return &OwnerKeyKeyringManager{
		logger:           logger,
		keysDir:          keysDir,
		keyringAvailable: keyringAvailable,
	}
}

// StorePrivateKeyInKeyring stores encrypted private key in system keyring
func (okm *OwnerKeyKeyringManager) StorePrivateKeyInKeyring(keyPair *types.KeyPair, passphrase string) error {
	okm.logger.LogStep("owner_key_keyring_storage_start", map[string]interface{}{
		"key_id":            keyPair.ID,
		"keyring_available": okm.keyringAvailable,
	})

	if !okm.keyringAvailable {
		return okm.storePrivateKeyToFile(keyPair, passphrase)
	}

	// 1. Encrypt private key with cascade encryption
	encryptedKey, err := okm.encryptPrivateKeyCascade(keyPair.PrivateKey, passphrase)
	if err != nil {
		return fmt.Errorf("failed to encrypt private key: %w", err)
	}

	// 2. Create keyring entry with metadata
	keyringData := map[string]interface{}{
		"encrypted_key": base64.StdEncoding.EncodeToString(encryptedKey),
		"key_id":        keyPair.ID,
		"algorithm":     keyPair.Algorithm,
		"fingerprint":   keyPair.Fingerprint,
		"created_at":    keyPair.CreatedAt.Format(time.RFC3339),
		"metadata":      keyPair.Metadata,
	}

	jsonData, err := json.Marshal(keyringData)
	if err != nil {
		return fmt.Errorf("failed to marshal keyring data: %w", err)
	}

	// 3. Store in system keyring
	if err := keyring.Set(KeyringServiceOwnerKey, KeyringUserPrivateKey, string(jsonData)); err != nil {
		okm.logger.LogWarning("keyring_save_failed", map[string]interface{}{
			"error":            err.Error(),
			"fallback_to_file": true,
		})
		return okm.storePrivateKeyToFile(keyPair, passphrase)
	}

	// 4. Save only public key and metadata to disk
	if err := okm.savePublicKeyAndMetadata(keyPair); err != nil {
		return err
	}

	okm.logger.LogStep("owner_key_keyring_storage_completed", map[string]interface{}{
		"key_id":         keyPair.ID,
		"storage_method": "keyring",
		"fingerprint":    keyPair.Fingerprint,
	})

	return nil
}

// LoadPrivateKeyFromKeyring loads and decrypts private key from keyring
func (okm *OwnerKeyKeyringManager) LoadPrivateKeyFromKeyring(keyID string, passphrase string) (*types.KeyPair, error) {
	okm.logger.LogStep("owner_key_keyring_load_start", map[string]interface{}{
		"key_id":            keyID,
		"keyring_available": okm.keyringAvailable,
	})

	if !okm.keyringAvailable {
		return okm.loadPrivateKeyFromFile(keyID, passphrase)
	}

	// 1. Load from keyring
	jsonData, err := keyring.Get(KeyringServiceOwnerKey, KeyringUserPrivateKey)
	if err != nil {
		okm.logger.LogWarning("keyring_load_failed", map[string]interface{}{
			"error":            err.Error(),
			"fallback_to_file": true,
		})
		return okm.loadPrivateKeyFromFile(keyID, passphrase)
	}

	// 2. Unmarshal keyring data
	var keyringData map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &keyringData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal keyring data: %w", err)
	}

	// 3. Decode and decrypt private key
	encryptedKeyB64 := keyringData["encrypted_key"].(string)
	encryptedKey, err := base64.StdEncoding.DecodeString(encryptedKeyB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode encrypted key: %w", err)
	}

	privateKey, err := okm.decryptPrivateKeyCascade(encryptedKey, passphrase)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt private key: %w", err)
	}

	// 4. Load public key from disk
	publicKeyPath := filepath.Join(okm.keysDir, "owner.key.pub")
	publicKey, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key: %w", err)
	}

	// 5. Reconstruct KeyPair
	keyPair := &types.KeyPair{
		ID:          keyringData["key_id"].(string),
		Algorithm:   keyringData["algorithm"].(string),
		PrivateKey:  privateKey,
		PublicKey:   publicKey,
		Fingerprint: keyringData["fingerprint"].(string),
		Metadata:    keyringData["metadata"].(map[string]string),
	}

	if createdAtStr, ok := keyringData["created_at"].(string); ok {
		if createdAt, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
			keyPair.CreatedAt = createdAt
		}
	}

	okm.logger.LogStep("owner_key_keyring_load_completed", map[string]interface{}{
		"key_id":         keyPair.ID,
		"storage_method": "keyring",
	})

	return keyPair, nil
}

// DeletePrivateKeyFromKeyring deletes private key from keyring
func (okm *OwnerKeyKeyringManager) DeletePrivateKeyFromKeyring() error {
	if !okm.keyringAvailable {
		return okm.deletePrivateKeyFile()
	}

	if err := keyring.Delete(KeyringServiceOwnerKey, KeyringUserPrivateKey); err != nil {
		okm.logger.LogWarning("keyring_delete_failed", map[string]interface{}{
			"error": err.Error(),
		})
		return okm.deletePrivateKeyFile()
	}

	return nil
}

// encryptPrivateKeyCascade implements cascade encryption (same as key_manager.go)
func (okm *OwnerKeyKeyringManager) encryptPrivateKeyCascade(privateKey []byte, passphrase string) ([]byte, error) {
	// 1. Derivar chave com argon2id (resistente a GPU)
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	key1 := argon2.IDKey([]byte(passphrase), salt, 3, 64*1024, 4, 32)

	// 2. Segunda derivação com PBKDF2 (defesa em profundidade)
	key2 := pbkdf2.Key(key1, salt, 100000, 32, sha256.New)

	// 3. Criptografar com AES-256-GCM
	block, err := aes.NewCipher(key2)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, privateKey, nil)

	// 4. Computar HMAC para verificação de integridade
	mac := hmac.New(sha256.New, key2)
	mac.Write(ciphertext)
	hmacValue := mac.Sum(nil)

	// 5. Empacotar: salt + ciphertext + hmac
	result := append(salt, ciphertext...)
	result = append(result, hmacValue...)

	return result, nil
}

// decryptPrivateKeyCascade implements cascade decryption
func (okm *OwnerKeyKeyringManager) decryptPrivateKeyCascade(encryptedKey []byte, passphrase string) ([]byte, error) {
	// 1. Extrair componentes: salt + ciphertext + hmac
	if len(encryptedKey) < 96 { // 32 (salt) + 32 (hmac) + mínimo para ciphertext
		return nil, fmt.Errorf("encrypted data too short")
	}

	salt := encryptedKey[:32]
	ciphertext := encryptedKey[32 : len(encryptedKey)-32]
	hmacValue := encryptedKey[len(encryptedKey)-32:]

	// 2. Derivar chaves (mesmo processo da criptografia)
	key1 := argon2.IDKey([]byte(passphrase), salt, 3, 64*1024, 4, 32)
	key2 := pbkdf2.Key(key1, salt, 100000, 32, sha256.New)

	// 3. Verificar HMAC
	mac := hmac.New(sha256.New, key2)
	mac.Write(ciphertext)
	expectedHMAC := mac.Sum(nil)

	if !hmac.Equal(hmacValue, expectedHMAC) {
		return nil, fmt.Errorf("HMAC verification failed - data may be corrupted")
	}

	// 4. Descriptografar com AES-256-GCM
	block, err := aes.NewCipher(key2)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// savePublicKeyAndMetadata saves only public data to disk
func (okm *OwnerKeyKeyringManager) savePublicKeyAndMetadata(keyPair *types.KeyPair) error {
	publicKeyPath := filepath.Join(okm.keysDir, "owner.key.pub")
	if err := os.WriteFile(publicKeyPath, keyPair.PublicKey, 0644); err != nil {
		return fmt.Errorf("failed to write public key: %w", err)
	}

	metadataPath := filepath.Join(okm.keysDir, "owner.meta")
	metadata, err := json.MarshalIndent(keyPair.Metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	if err := os.WriteFile(metadataPath, metadata, 0600); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	return nil
}

// Fallback methods for file storage (only if keyring is not available)
func (okm *OwnerKeyKeyringManager) storePrivateKeyToFile(keyPair *types.KeyPair, passphrase string) error {
	// Fallback: use existing file-based storage from key_manager.go
	encryptedKey, err := okm.encryptPrivateKeyCascade(keyPair.PrivateKey, passphrase)
	if err != nil {
		return err
	}

	privateKeyPath := filepath.Join(okm.keysDir, "owner.key")
	if err := os.WriteFile(privateKeyPath, encryptedKey, 0600); err != nil {
		return err
	}

	return okm.savePublicKeyAndMetadata(keyPair)
}

func (okm *OwnerKeyKeyringManager) loadPrivateKeyFromFile(keyID string, passphrase string) (*types.KeyPair, error) {
	// Fallback: use existing file-based loading from key_manager.go
	privateKeyPath := filepath.Join(okm.keysDir, "owner.key")
	encryptedKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	privateKey, err := okm.decryptPrivateKeyCascade(encryptedKey, passphrase)
	if err != nil {
		return nil, err
	}

	// Load public key
	publicKeyPath := filepath.Join(okm.keysDir, "owner.key.pub")
	publicKey, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	// Load metadata
	metadataPath := filepath.Join(okm.keysDir, "owner.meta")
	var metadata map[string]string
	if metadataBytes, err := os.ReadFile(metadataPath); err == nil {
		json.Unmarshal(metadataBytes, &metadata)
	}

	return &types.KeyPair{
		ID:          keyID,
		Algorithm:   "Ed25519",
		PrivateKey:  privateKey,
		PublicKey:   publicKey,
		Fingerprint: fmt.Sprintf("%x", sha256.Sum256(publicKey)),
		Metadata:    metadata,
	}, nil
}

func (okm *OwnerKeyKeyringManager) deletePrivateKeyFile() error {
	privateKeyPath := filepath.Join(okm.keysDir, "owner.key")
	return os.Remove(privateKeyPath)
}
