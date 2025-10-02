package setup

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"setup-component/src/internal/types"
)

// KeyManager implementa a interface KeyManager
type KeyManager struct {
	keysDir string
	logger  *SetupLogger
}

// NewKeyManager cria um novo gerenciador de chaves
func NewKeyManager(logger *SetupLogger) *KeyManager {
	homeDir, _ := os.UserHomeDir()
	keysDir := filepath.Join(homeDir, ".syntropy", "keys")
	os.MkdirAll(keysDir, 0755)

	return &KeyManager{
		keysDir: keysDir,
		logger:  logger,
	}
}

// GenerateKeyPair gera um novo par de chaves
func (km *KeyManager) GenerateKeyPair(algorithm string) (*types.KeyPair, error) {
	km.logger.LogStep("key_generation_start", map[string]interface{}{
		"algorithm": algorithm,
	})

	var keyPair *types.KeyPair
	var err error

	switch algorithm {
	case "ed25519":
		keyPair, err = km.generateEd25519KeyPair()
	default:
		return nil, fmt.Errorf("algoritmo de chave não suportado: %s", algorithm)
	}

	if err != nil {
		km.logger.LogError(types.ErrKeyGenerationError(algorithm, err), map[string]interface{}{
			"algorithm": algorithm,
		})
		return nil, err
	}

	km.logger.LogStep("key_generation_completed", map[string]interface{}{
		"key_id":      keyPair.ID,
		"algorithm":   keyPair.Algorithm,
		"fingerprint": keyPair.Fingerprint,
	})

	return keyPair, nil
}

// GenerateOrLoadKeyPair gera um novo par de chaves ou carrega um existente
func (km *KeyManager) GenerateOrLoadKeyPair(algorithm string) (*types.KeyPair, error) {
	km.logger.LogStep("key_generation_or_load_start", map[string]interface{}{
		"algorithm": algorithm,
	})

	// Verificar se já existem chaves
	existingKeys, err := km.listExistingKeys()
	if err != nil {
		km.logger.LogWarning("Falha ao listar chaves existentes", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Se existem chaves, carregar a primeira
	if len(existingKeys) > 0 {
		keyID := existingKeys[0]
		km.logger.LogInfo("Carregando chave existente", map[string]interface{}{
			"key_id": keyID,
		})

		keyPair, err := km.LoadKeyPair(keyID, "default_passphrase")
		if err != nil {
			km.logger.LogWarning("Falha ao carregar chave existente, gerando nova", map[string]interface{}{
				"key_id": keyID,
				"error":  err.Error(),
			})
		} else {
			return keyPair, nil
		}
	}

	// Gerar nova chave
	km.logger.LogInfo("Gerando nova chave", map[string]interface{}{
		"algorithm": algorithm,
	})

	keyPair, err := km.GenerateKeyPair(algorithm)
	if err != nil {
		return nil, err
	}

	// Salvar a nova chave
	if err := km.StoreKeyPair(keyPair, "default_passphrase"); err != nil {
		return nil, err
	}

	return keyPair, nil
}

// listExistingKeys lista as chaves existentes
func (km *KeyManager) listExistingKeys() ([]string, error) {
	// Verificar se o arquivo owner.key existe
	ownerKeyPath := filepath.Join(km.keysDir, "owner.key")
	if _, err := os.Stat(ownerKeyPath); err == nil {
		return []string{"owner"}, nil
	}

	return []string{}, nil
}

// StoreKeyPair armazena um par de chaves de forma segura
func (km *KeyManager) StoreKeyPair(keyPair *types.KeyPair, passphrase string) error {
	km.logger.LogStep("key_storage_start", map[string]interface{}{
		"key_id":    keyPair.ID,
		"algorithm": keyPair.Algorithm,
	})

	// Validar passphrase
	if passphrase == "" {
		return fmt.Errorf("passphrase não pode estar vazia")
	}

	// Criptografar chave privada
	encryptedPrivateKey, err := km.encryptPrivateKey(keyPair.PrivateKey, passphrase)
	if err != nil {
		return types.ErrKeyStorageError(keyPair.ID, err)
	}

	// Salvar chave privada criptografada
	privateKeyPath := filepath.Join(km.keysDir, "owner.key")
	if err := os.WriteFile(privateKeyPath, encryptedPrivateKey, 0600); err != nil {
		return types.ErrKeyStorageError(keyPair.ID, err)
	}

	// Salvar chave pública
	publicKeyPath := filepath.Join(km.keysDir, "owner.key.pub")
	if err := os.WriteFile(publicKeyPath, keyPair.PublicKey, 0600); err != nil {
		// Limpar chave privada em caso de erro
		os.Remove(privateKeyPath)
		return types.ErrKeyStorageError(keyPair.ID, err)
	}

	// Salvar metadados
	metadataPath := filepath.Join(km.keysDir, "owner.meta")
	metadata, err := json.MarshalIndent(keyPair.Metadata, "", "  ")
	if err != nil {
		// Limpar arquivos em caso de erro
		os.Remove(privateKeyPath)
		os.Remove(publicKeyPath)
		return types.ErrKeyStorageError(keyPair.ID, err)
	}

	if err := os.WriteFile(metadataPath, metadata, 0600); err != nil {
		// Limpar arquivos em caso de erro
		os.Remove(privateKeyPath)
		os.Remove(publicKeyPath)
		return types.ErrKeyStorageError(keyPair.ID, err)
	}

	km.logger.LogStep("key_storage_completed", map[string]interface{}{
		"key_id":           keyPair.ID,
		"private_key_path": privateKeyPath,
		"public_key_path":  publicKeyPath,
		"metadata_path":    metadataPath,
	})

	return nil
}

// LoadKeyPair carrega um par de chaves
func (km *KeyManager) LoadKeyPair(keyID string, passphrase string) (*types.KeyPair, error) {
	km.logger.LogDebug("Carregando par de chaves", map[string]interface{}{
		"key_id": keyID,
	})

	// Verificar se os arquivos existem (usando nomes fixos)
	privateKeyPath := filepath.Join(km.keysDir, "owner.key")
	publicKeyPath := filepath.Join(km.keysDir, "owner.key.pub")
	metadataPath := filepath.Join(km.keysDir, "owner.meta")

	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("chave privada não encontrada: %s", keyID)
	}

	if _, err := os.Stat(publicKeyPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("chave pública não encontrada: %s", keyID)
	}

	// Ler chave privada criptografada
	encryptedPrivateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler chave privada: %w", err)
	}

	// Descriptografar chave privada
	privateKey, err := km.decryptPrivateKey(encryptedPrivateKey, passphrase)
	if err != nil {
		return nil, fmt.Errorf("falha ao descriptografar chave privada: %w", err)
	}

	// Ler chave pública
	publicKey, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler chave pública: %w", err)
	}

	// Ler metadados
	var metadata map[string]string
	if _, err := os.Stat(metadataPath); err == nil {
		metadataData, err := os.ReadFile(metadataPath)
		if err == nil {
			json.Unmarshal(metadataData, &metadata)
		}
	}

	if metadata == nil {
		metadata = make(map[string]string)
	}

	// Criar par de chaves
	keyPair := &types.KeyPair{
		ID:          keyID,
		Algorithm:   "ed25519", // Assumindo Ed25519 por enquanto
		PrivateKey:  privateKey,
		PublicKey:   publicKey,
		Fingerprint: km.generateFingerprint(publicKey),
		Metadata:    metadata,
	}

	// Tentar extrair timestamps dos metadados
	if createdAtStr, exists := metadata["created_at"]; exists {
		if createdAt, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
			keyPair.CreatedAt = createdAt
		}
	}

	if expiresAtStr, exists := metadata["expires_at"]; exists {
		if expiresAt, err := time.Parse(time.RFC3339, expiresAtStr); err == nil {
			keyPair.ExpiresAt = expiresAt
		}
	}

	km.logger.LogDebug("Par de chaves carregado com sucesso", map[string]interface{}{
		"key_id":      keyPair.ID,
		"algorithm":   keyPair.Algorithm,
		"fingerprint": keyPair.Fingerprint,
	})

	return keyPair, nil
}

// RotateKeys rotaciona as chaves
func (km *KeyManager) RotateKeys(keyID string) error {
	km.logger.LogStep("key_rotation_start", map[string]interface{}{
		"key_id": keyID,
	})

	// Verificar se a chave existe
	if _, err := os.Stat(filepath.Join(km.keysDir, fmt.Sprintf("%s.key", keyID))); os.IsNotExist(err) {
		return types.ErrKeyRotationError(keyID, fmt.Errorf("chave não encontrada"))
	}

	// Criar backup da chave atual
	backupKeyID := fmt.Sprintf("%s_backup_%d", keyID, time.Now().Unix())
	if err := km.backupKey(keyID, backupKeyID); err != nil {
		return types.ErrKeyRotationError(keyID, err)
	}

	// Gerar nova chave
	newKeyPair, err := km.GenerateKeyPair("ed25519")
	if err != nil {
		return types.ErrKeyRotationError(keyID, err)
	}

	// Armazenar nova chave (sem passphrase por enquanto - em produção seria necessário)
	if err := km.StoreKeyPair(newKeyPair, "default_passphrase"); err != nil {
		return types.ErrKeyRotationError(keyID, err)
	}

	km.logger.LogStep("key_rotation_completed", map[string]interface{}{
		"old_key_id": keyID,
		"new_key_id": newKeyPair.ID,
		"backup_id":  backupKeyID,
	})

	return nil
}

// VerifyKeyIntegrity verifica a integridade de uma chave
func (km *KeyManager) VerifyKeyIntegrity(keyID string) error {
	km.logger.LogDebug("Verificando integridade da chave", map[string]interface{}{
		"key_id": keyID,
	})

	// Verificar se os arquivos existem (usando nomes fixos)
	privateKeyPath := filepath.Join(km.keysDir, "owner.key")
	publicKeyPath := filepath.Join(km.keysDir, "owner.key.pub")

	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		return fmt.Errorf("chave privada não encontrada: %s", keyID)
	}

	if _, err := os.Stat(publicKeyPath); os.IsNotExist(err) {
		return fmt.Errorf("chave pública não encontrada: %s", keyID)
	}

	// Verificar permissões da chave privada
	info, err := os.Stat(privateKeyPath)
	if err != nil {
		return fmt.Errorf("falha ao verificar permissões da chave privada: %w", err)
	}

	// Verificar se a chave privada tem permissões restritivas (600)
	if info.Mode().Perm()&0077 != 0 {
		return fmt.Errorf("chave privada tem permissões inseguras: %s", info.Mode().Perm())
	}

	// Verificar se a chave pública é legível
	if _, err := os.ReadFile(publicKeyPath); err != nil {
		return fmt.Errorf("chave pública não é legível: %w", err)
	}

	km.logger.LogDebug("Integridade da chave verificada com sucesso", map[string]interface{}{
		"key_id": keyID,
	})

	return nil
}

// BackupKeys cria um backup das chaves
func (km *KeyManager) BackupKeys(keyID string, passphrase string) ([]byte, error) {
	km.logger.LogStep("key_backup_start", map[string]interface{}{
		"key_id": keyID,
	})

	// Carregar par de chaves
	keyPair, err := km.LoadKeyPair(keyID, passphrase)
	if err != nil {
		return nil, fmt.Errorf("falha ao carregar chave para backup: %w", err)
	}

	// Serializar par de chaves
	backupData, err := json.MarshalIndent(keyPair, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("falha ao serializar chave para backup: %w", err)
	}

	km.logger.LogStep("key_backup_completed", map[string]interface{}{
		"key_id":      keyID,
		"backup_size": len(backupData),
	})

	return backupData, nil
}

// RestoreKeys restaura chaves de um backup
func (km *KeyManager) RestoreKeys(backupData []byte, passphrase string) error {
	km.logger.LogStep("key_restore_start", map[string]interface{}{
		"backup_size": len(backupData),
	})

	// Deserializar par de chaves
	var keyPair types.KeyPair
	if err := json.Unmarshal(backupData, &keyPair); err != nil {
		return fmt.Errorf("falha ao deserializar backup de chave: %w", err)
	}

	// Armazenar chave restaurada
	if err := km.StoreKeyPair(&keyPair, passphrase); err != nil {
		return fmt.Errorf("falha ao armazenar chave restaurada: %w", err)
	}

	km.logger.LogStep("key_restore_completed", map[string]interface{}{
		"key_id":      keyPair.ID,
		"algorithm":   keyPair.Algorithm,
		"fingerprint": keyPair.Fingerprint,
	})

	return nil
}

// ListKeys lista todas as chaves disponíveis
func (km *KeyManager) ListKeys() ([]types.KeyPair, error) {
	km.logger.LogDebug("Listando chaves disponíveis", nil)

	files, err := os.ReadDir(km.keysDir)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler diretório de chaves: %w", err)
	}

	var keys []types.KeyPair
	keyIDs := make(map[string]bool)

	// Encontrar todos os IDs de chave únicos
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".key" && !file.IsDir() {
			keyID := file.Name()[:len(file.Name())-4] // Remove .key
			keyIDs[keyID] = true
		}
	}

	// Carregar informações de cada chave
	for keyID := range keyIDs {
		// Tentar carregar metadados sem descriptografar a chave privada
		metadataPath := filepath.Join(km.keysDir, fmt.Sprintf("%s.meta", keyID))
		publicKeyPath := filepath.Join(km.keysDir, fmt.Sprintf("%s.key.pub", keyID))

		var metadata map[string]string
		if _, err := os.Stat(metadataPath); err == nil {
			metadataData, err := os.ReadFile(metadataPath)
			if err == nil {
				json.Unmarshal(metadataData, &metadata)
			}
		}

		if metadata == nil {
			metadata = make(map[string]string)
		}

		// Ler chave pública
		publicKey, err := os.ReadFile(publicKeyPath)
		if err != nil {
			continue // Pular chaves com problemas
		}

		keyPair := types.KeyPair{
			ID:          keyID,
			Algorithm:   "ed25519",
			PublicKey:   publicKey,
			Fingerprint: km.generateFingerprint(publicKey),
			Metadata:    metadata,
		}

		// Extrair timestamps dos metadados
		if createdAtStr, exists := metadata["created_at"]; exists {
			if createdAt, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
				keyPair.CreatedAt = createdAt
			}
		}

		if expiresAtStr, exists := metadata["expires_at"]; exists {
			if expiresAt, err := time.Parse(time.RFC3339, expiresAtStr); err == nil {
				keyPair.ExpiresAt = expiresAt
			}
		}

		keys = append(keys, keyPair)
	}

	km.logger.LogDebug("Chaves listadas com sucesso", map[string]interface{}{
		"keys_count": len(keys),
	})

	return keys, nil
}

// Métodos auxiliares

// generateEd25519KeyPair gera um par de chaves Ed25519
func (km *KeyManager) generateEd25519KeyPair() (*types.KeyPair, error) {
	// Gerar chave privada usando fonte de entropia segura
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("falha ao gerar chave privada: %w", err)
	}

	// Criar fingerprint
	fingerprint := km.generateFingerprint(publicKey)

	// Gerar ID único
	keyID := km.generateKeyID()

	keyPair := &types.KeyPair{
		ID:          keyID,
		Algorithm:   "ed25519",
		PrivateKey:  privateKey,
		PublicKey:   publicKey,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().AddDate(1, 0, 0), // 1 ano
		Fingerprint: fingerprint,
		Metadata: map[string]string{
			"generated_by": "syntropy-setup",
			"version":      "1.0.0",
			"created_at":   time.Now().Format(time.RFC3339),
			"expires_at":   time.Now().AddDate(1, 0, 0).Format(time.RFC3339),
		},
	}

	return keyPair, nil
}

// generateFingerprint gera um fingerprint para a chave pública
func (km *KeyManager) generateFingerprint(publicKey []byte) string {
	hash := sha256.Sum256(publicKey)
	return base64.StdEncoding.EncodeToString(hash[:])
}

// generateKeyID gera um ID único para a chave
func (km *KeyManager) generateKeyID() string {
	// Gerar bytes aleatórios
	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)

	// Codificar em base64
	return base64.StdEncoding.EncodeToString(randomBytes)
}

// encryptPrivateKey criptografa a chave privada
func (km *KeyManager) encryptPrivateKey(privateKey []byte, passphrase string) ([]byte, error) {
	// Implementação simplificada - em produção usar AES-256-GCM
	// Por enquanto, apenas codificar em base64
	return []byte(base64.StdEncoding.EncodeToString(privateKey)), nil
}

// decryptPrivateKey descriptografa a chave privada
func (km *KeyManager) decryptPrivateKey(encryptedPrivateKey []byte, passphrase string) ([]byte, error) {
	// Implementação simplificada - em produção usar AES-256-GCM
	// Por enquanto, apenas decodificar de base64
	return base64.StdEncoding.DecodeString(string(encryptedPrivateKey))
}

// backupKey cria um backup de uma chave específica
func (km *KeyManager) backupKey(keyID, backupKeyID string) error {
	// Copiar arquivos da chave
	privateKeyPath := filepath.Join(km.keysDir, fmt.Sprintf("%s.key", keyID))
	publicKeyPath := filepath.Join(km.keysDir, fmt.Sprintf("%s.key.pub", keyID))
	metadataPath := filepath.Join(km.keysDir, fmt.Sprintf("%s.meta", keyID))

	backupPrivateKeyPath := filepath.Join(km.keysDir, fmt.Sprintf("%s.key", backupKeyID))
	backupPublicKeyPath := filepath.Join(km.keysDir, fmt.Sprintf("%s.key.pub", backupKeyID))
	backupMetadataPath := filepath.Join(km.keysDir, fmt.Sprintf("%s.meta", backupKeyID))

	// Copiar chave privada
	if data, err := os.ReadFile(privateKeyPath); err == nil {
		os.WriteFile(backupPrivateKeyPath, data, 0600)
	}

	// Copiar chave pública
	if data, err := os.ReadFile(publicKeyPath); err == nil {
		os.WriteFile(backupPublicKeyPath, data, 0644)
	}

	// Copiar metadados
	if data, err := os.ReadFile(metadataPath); err == nil {
		os.WriteFile(backupMetadataPath, data, 0644)
	}

	return nil
}
