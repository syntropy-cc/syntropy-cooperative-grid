/**
 * File: token_manager.go
 * Purpose: Gerenciamento seguro de Grid Token usando Keyring do sistema operacional
 * Dependencies: github.com/zalando/go-keyring, crypto/rand, crypto/sha256
 * Exports: TokenManager struct, NewTokenManager function
 * Author: Syntropy Setup Component
 * Created: 2025-01-27
 * Modified: 2025-01-27
 * Version: 1.0.0
 *
 * Business Context:
 * O TokenManager é responsável pelo gerenciamento seguro de Grid Tokens,
 * implementando armazenamento criptográfico através do Keyring nativo do
 * sistema operacional. Garante que os tokens sejam gerados de forma segura,
 * armazenados com máxima proteção e gerenciados durante todo o ciclo de vida.
 */

package setup

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"setup-component/src/internal/types"

	"github.com/zalando/go-keyring"
)

// Constantes para o Keyring
const (
	KeyringService = "syntropy-grid"
	KeyringUser    = "grid-token"
	TokenVersion   = "1.0.0"
)

// TokenManager implementa a interface TokenManager para gerenciamento seguro de Grid Token
type TokenManager struct {
	logger           types.SetupLogger
	tokensDir        string
	keyringAvailable bool
}

// NewTokenManager cria um novo gerenciador de tokens
func NewTokenManager(logger types.SetupLogger) *TokenManager {
	// Detectar se o keyring está disponível
	keyringAvailable := isKeyringAvailable()

	// Criar diretório de tokens se necessário
	homeDir, _ := os.UserHomeDir()
	tokensDir := filepath.Join(homeDir, ".syntropy", "tokens")

	// Garantir que o diretório seja criado com permissões corretas
	if err := os.MkdirAll(tokensDir, 0700); err != nil {
		logger.LogWarning("failed to create tokens directory", map[string]interface{}{
			"error": err.Error(),
			"path":  tokensDir,
		})
	}

	return &TokenManager{
		logger:           logger,
		tokensDir:        tokensDir, // Agora usa ~/.syntropy/tokens/ em vez de backups
		keyringAvailable: keyringAvailable,
	}
}

// GenerateToken gera um novo Grid Token (UUID v4)
func (tm *TokenManager) GenerateToken() (string, error) {
	tm.logger.LogStep("token_generation_start", map[string]interface{}{
		"algorithm":         "UUID v4",
		"keyring_available": tm.keyringAvailable,
	})

	// Gerar 16 bytes aleatórios para UUID v4
	tokenBytes := make([]byte, 16)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Definir bits de versão e variante para UUID v4
	tokenBytes[6] = (tokenBytes[6] & 0x0f) | 0x40 // Version 4
	tokenBytes[8] = (tokenBytes[8] & 0x3f) | 0x80 // Variant bits

	// Formatar como UUID
	token := fmt.Sprintf("%x-%x-%x-%x-%x",
		tokenBytes[0:4],
		tokenBytes[4:6],
		tokenBytes[6:8],
		tokenBytes[8:10],
		tokenBytes[10:16],
	)

	tm.logger.LogStep("token_generation_completed", map[string]interface{}{
		"token_preview": token[:8] + "...[HIDDEN]",
		"token_length":  len(token),
	})

	return token, nil
}

// SaveToken salva token no keyring do sistema
func (tm *TokenManager) SaveToken(token string) error {
	tm.logger.LogStep("token_save_start", map[string]interface{}{
		"keyring_available": tm.keyringAvailable,
		"service":           KeyringService,
	})

	if !tm.keyringAvailable {
		return tm.saveTokenToFile(token)
	}

	// Validar token antes de salvar
	if err := tm.ValidateToken(token); err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	// Salvar no keyring
	if err := keyring.Set(KeyringService, KeyringUser, token); err != nil {
		tm.logger.LogWarning("keyring_save_failed", map[string]interface{}{
			"error":            err.Error(),
			"fallback_to_file": true,
		})
		return tm.saveTokenToFile(token)
	}

	tm.logger.LogStep("token_save_completed", map[string]interface{}{
		"storage_method": "keyring",
		"service":        KeyringService,
		"token_preview":  token[:8] + "...[HIDDEN]",
	})

	return nil
}

// LoadToken carrega token do keyring
func (tm *TokenManager) LoadToken() (string, error) {
	tm.logger.LogStep("token_load_start", map[string]interface{}{
		"keyring_available": tm.keyringAvailable,
	})

	if !tm.keyringAvailable {
		return tm.loadTokenFromFile()
	}

	// Tentar carregar do keyring primeiro
	token, err := keyring.Get(KeyringService, KeyringUser)
	if err != nil {
		tm.logger.LogWarning("keyring_load_failed", map[string]interface{}{
			"error":            err.Error(),
			"fallback_to_file": true,
		})
		return tm.loadTokenFromFile()
	}

	// Validar token carregado
	if err := tm.ValidateToken(token); err != nil {
		return "", fmt.Errorf("loaded token is invalid: %w", err)
	}

	tm.logger.LogStep("token_load_completed", map[string]interface{}{
		"storage_method": "keyring",
		"token_preview":  token[:8] + "...[HIDDEN]",
	})

	return token, nil
}

// DeleteToken remove token do keyring
func (tm *TokenManager) DeleteToken() error {
	tm.logger.LogStep("token_delete_start", map[string]interface{}{
		"keyring_available": tm.keyringAvailable,
	})

	if !tm.keyringAvailable {
		return tm.deleteTokenFromFile()
	}

	// Tentar remover do keyring
	if err := keyring.Delete(KeyringService, KeyringUser); err != nil {
		tm.logger.LogWarning("keyring_delete_failed", map[string]interface{}{
			"error":            err.Error(),
			"fallback_to_file": true,
		})
		return tm.deleteTokenFromFile()
	}

	tm.logger.LogStep("token_delete_completed", map[string]interface{}{
		"storage_method": "keyring",
	})

	return nil
}

// TokenExists verifica se um token existe
func (tm *TokenManager) TokenExists() (bool, error) {
	tm.logger.LogStep("token_exists_check", map[string]interface{}{
		"keyring_available": tm.keyringAvailable,
	})

	if !tm.keyringAvailable {
		return tm.tokenExistsInFile()
	}

	// Verificar no keyring
	_, err := keyring.Get(KeyringService, KeyringUser)
	if err != nil {
		if err == keyring.ErrNotFound {
			// Verificar se existe em arquivo como fallback
			return tm.tokenExistsInFile()
		}
		return false, fmt.Errorf("failed to check token existence: %w", err)
	}

	return true, nil
}

// RotateToken gera um novo token e substitui o existente
func (tm *TokenManager) RotateToken() (string, error) {
	tm.logger.LogStep("token_rotation_start", map[string]interface{}{
		"keyring_available": tm.keyringAvailable,
	})

	// Gerar novo token
	newToken, err := tm.GenerateToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate new token: %w", err)
	}

	// Salvar novo token
	if err := tm.SaveToken(newToken); err != nil {
		return "", fmt.Errorf("failed to save new token: %w", err)
	}

	tm.logger.LogStep("token_rotation_completed", map[string]interface{}{
		"new_token_preview": newToken[:8] + "...[HIDDEN]",
	})

	return newToken, nil
}

// ExportToken exporta token para backup
func (tm *TokenManager) ExportToken(outputPath string) error {
	tm.logger.LogStep("token_export_start", map[string]interface{}{
		"output_path": outputPath,
	})

	// Carregar token
	token, err := tm.LoadToken()
	if err != nil {
		return fmt.Errorf("failed to load token for export: %w", err)
	}

	// Criar backup
	backup := types.TokenBackup{
		Token:     token,
		CreatedAt: time.Now(),
		Version:   TokenVersion,
		Checksum:  tm.calculateChecksum(token),
	}

	// Serializar backup
	data, err := json.MarshalIndent(backup, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal backup: %w", err)
	}

	// Escrever arquivo
	if err := os.WriteFile(outputPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}

	tm.logger.LogStep("token_export_completed", map[string]interface{}{
		"output_path": outputPath,
		"backup_size": len(data),
	})

	return nil
}

// ImportToken importa token de backup
func (tm *TokenManager) ImportToken(inputPath string) error {
	tm.logger.LogStep("token_import_start", map[string]interface{}{
		"input_path": inputPath,
	})

	// Ler arquivo de backup
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %w", err)
	}

	// Deserializar backup
	var backup types.TokenBackup
	if err := json.Unmarshal(data, &backup); err != nil {
		return fmt.Errorf("failed to unmarshal backup: %w", err)
	}

	// Validar checksum
	expectedChecksum := tm.calculateChecksum(backup.Token)
	if backup.Checksum != expectedChecksum {
		return fmt.Errorf("backup checksum validation failed")
	}

	// Validar token
	if err := tm.ValidateToken(backup.Token); err != nil {
		return fmt.Errorf("imported token is invalid: %w", err)
	}

	// Salvar token
	if err := tm.SaveToken(backup.Token); err != nil {
		return fmt.Errorf("failed to save imported token: %w", err)
	}

	tm.logger.LogStep("token_import_completed", map[string]interface{}{
		"input_path":    inputPath,
		"token_preview": backup.Token[:8] + "...[HIDDEN]",
	})

	return nil
}

// ValidateToken valida formato e estrutura do token
func (tm *TokenManager) ValidateToken(token string) error {
	// Verificar comprimento (UUID v4 tem 36 caracteres)
	if len(token) != 36 {
		return fmt.Errorf("invalid token length: expected 36, got %d", len(token))
	}

	// Verificar formato UUID v4
	if !isValidUUIDv4(token) {
		return fmt.Errorf("invalid UUID v4 format")
	}

	return nil
}

// Métodos auxiliares privados

// isKeyringAvailable verifica se o keyring está disponível
func isKeyringAvailable() bool {
	// Tentar uma operação simples para verificar disponibilidade
	err := keyring.Set("test-service", "test-user", "test-value")
	if err != nil {
		return false
	}

	// Limpar o teste
	keyring.Delete("test-service", "test-user")
	return true
}

// saveTokenToFile salva token em arquivo como fallback
func (tm *TokenManager) saveTokenToFile(token string) error {
	// Criar diretório de tokens se não existir
	if err := os.MkdirAll(tm.tokensDir, 0700); err != nil {
		return fmt.Errorf("failed to create tokens directory: %w", err)
	}

	// Criar backup
	backup := types.TokenBackup{
		Token:     token,
		CreatedAt: time.Now(),
		Version:   TokenVersion,
		Checksum:  tm.calculateChecksum(token),
	}

	// Serializar e salvar
	data, err := json.MarshalIndent(backup, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal token backup: %w", err)
	}

	backupPath := filepath.Join(tm.tokensDir, "grid-token.json")
	if err := os.WriteFile(backupPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write token backup: %w", err)
	}

	tm.logger.LogStep("token_save_fallback", map[string]interface{}{
		"storage_method": "file",
		"tokens_path":    backupPath,
	})

	return nil
}

// loadTokenFromFile carrega token de arquivo como fallback
func (tm *TokenManager) loadTokenFromFile() (string, error) {
	backupPath := filepath.Join(tm.tokensDir, "grid-token.json")

	// Verificar se arquivo existe
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return "", fmt.Errorf("token not found in keyring or file")
	}

	// Ler arquivo
	data, err := os.ReadFile(backupPath)
	if err != nil {
		return "", fmt.Errorf("failed to read token backup: %w", err)
	}

	// Deserializar backup
	var backup types.TokenBackup
	if err := json.Unmarshal(data, &backup); err != nil {
		return "", fmt.Errorf("failed to unmarshal token backup: %w", err)
	}

	// Validar checksum
	expectedChecksum := tm.calculateChecksum(backup.Token)
	if backup.Checksum != expectedChecksum {
		return "", fmt.Errorf("token backup checksum validation failed")
	}

	tm.logger.LogStep("token_load_fallback", map[string]interface{}{
		"storage_method": "file",
		"tokens_path":    backupPath,
	})

	return backup.Token, nil
}

// deleteTokenFromFile remove token de arquivo
func (tm *TokenManager) deleteTokenFromFile() error {
	backupPath := filepath.Join(tm.tokensDir, "grid-token.json")

	if err := os.Remove(backupPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete token backup: %w", err)
	}

	tm.logger.LogStep("token_delete_fallback", map[string]interface{}{
		"storage_method": "file",
		"tokens_path":    backupPath,
	})

	return nil
}

// tokenExistsInFile verifica se token existe em arquivo
func (tm *TokenManager) tokenExistsInFile() (bool, error) {
	backupPath := filepath.Join(tm.tokensDir, "grid-token.json")
	_, err := os.Stat(backupPath)
	return !os.IsNotExist(err), nil
}

// calculateChecksum calcula checksum SHA256 do token
func (tm *TokenManager) calculateChecksum(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// isValidUUIDv4 valida formato UUID v4
func isValidUUIDv4(uuid string) bool {
	if len(uuid) != 36 {
		return false
	}

	// Verificar posições dos hífens
	if uuid[8] != '-' || uuid[13] != '-' || uuid[18] != '-' || uuid[23] != '-' {
		return false
	}

	// Verificar versão 4 (bit 12-15 do time_hi_and_version)
	if uuid[14] != '4' {
		return false
	}

	// Verificar variante (bits 6-7 do clock_seq_hi_and_reserved)
	if uuid[19] < '8' || uuid[19] > 'b' {
		return false
	}

	return true
}
