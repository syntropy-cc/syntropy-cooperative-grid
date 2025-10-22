package node

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"node-component/src/internal/types"
)

// Constantes para o Keyring (mesmas do SetupManager)
const (
	KeyringService = "syntropy-grid"
	KeyringUser    = "grid-token"
)

// SetupAdapter adapts Setup Component interfaces to Node Component
type SetupAdapter struct {
	logger types.Logger
}

// NewSetupAdapter creates a new setup adapter
func NewSetupAdapter(logger types.Logger) *SetupAdapter {
	return &SetupAdapter{
		logger: logger,
	}
}

// SetupTokenManagerAdapter adapts Setup TokenManager to Node TokenIntegration interface
type SetupTokenManagerAdapter struct {
	adapter          *SetupAdapter
	keyringAvailable bool
	tokensDir        string
}

// NewSetupTokenManagerAdapter creates a new token manager adapter
func NewSetupTokenManagerAdapter(adapter *SetupAdapter) *SetupTokenManagerAdapter {
	// Criar diretório de tokens se necessário
	homeDir, _ := os.UserHomeDir()
	tokensDir := filepath.Join(homeDir, ".syntropy", "tokens")

	// Garantir que o diretório seja criado com permissões corretas
	if err := os.MkdirAll(tokensDir, 0700); err != nil {
		adapter.logger.Warn("Failed to create tokens directory", "error", err, "path", tokensDir)
	}

	return &SetupTokenManagerAdapter{
		adapter:          adapter,
		keyringAvailable: false, // Simplificado: usar apenas arquivo
		tokensDir:        tokensDir,
	}
}

// GenerateToken generates a new token (delegates to Setup Component)
func (stma *SetupTokenManagerAdapter) GenerateToken() (string, error) {
	// This would normally call the Setup Component's TokenManager
	// For now, we'll implement a basic version that checks if setup was run
	return stma.LoadToken()
}

// SaveToken saves a token (delegates to Setup Component)
func (stma *SetupTokenManagerAdapter) SaveToken(token string) error {
	// This would normally call the Setup Component's TokenManager
	// For now, we'll return an error as this should be done via setup command
	return fmt.Errorf("token saving should be done via 'syntropy setup' command")
}

// LoadToken loads the token from Setup Component storage
func (stma *SetupTokenManagerAdapter) LoadToken() (string, error) {
	// Para simplificar e evitar dependências externas, usar apenas fallback de arquivo
	// O SetupManager já gerencia o keyring, então o NodeManager pode usar apenas o arquivo
	return stma.loadTokenFromFile()
}

// DeleteToken deletes the token (delegates to Setup Component)
func (stma *SetupTokenManagerAdapter) DeleteToken() error {
	// This would normally call the Setup Component's TokenManager
	return fmt.Errorf("token deletion should be done via 'syntropy setup' command")
}

// TokenExists checks if token exists (delegates to Setup Component)
func (stma *SetupTokenManagerAdapter) TokenExists() (bool, error) {
	// Para simplificar e evitar dependências externas, verificar apenas arquivo
	// O SetupManager já gerencia o keyring, então o NodeManager pode usar apenas o arquivo
	return stma.tokenExistsInFile()
}

// RotateToken rotates the token (delegates to Setup Component)
func (stma *SetupTokenManagerAdapter) RotateToken() (string, error) {
	// This would normally call the Setup Component's TokenManager
	return "", fmt.Errorf("token rotation should be done via 'syntropy setup' command")
}

// ExportToken exports the token (delegates to Setup Component)
func (stma *SetupTokenManagerAdapter) ExportToken(outputPath string) error {
	// This would normally call the Setup Component's TokenManager
	return fmt.Errorf("token export should be done via 'syntropy setup' command")
}

// ImportToken imports the token (delegates to Setup Component)
func (stma *SetupTokenManagerAdapter) ImportToken(inputPath string) error {
	// This would normally call the Setup Component's TokenManager
	return fmt.Errorf("token import should be done via 'syntropy setup' command")
}

// ValidateToken validates the token (delegates to Setup Component)
func (stma *SetupTokenManagerAdapter) ValidateToken(token string) error {
	// Basic validation - check if token is not empty and has reasonable length
	if token == "" {
		return fmt.Errorf("token cannot be empty")
	}

	if len(token) < 10 {
		return fmt.Errorf("token appears to be too short")
	}

	// This would normally call the Setup Component's TokenManager for full validation
	return nil
}

// loadTokenFromFile carrega token de arquivo como fallback (mesma lógica do SetupManager)
func (stma *SetupTokenManagerAdapter) loadTokenFromFile() (string, error) {
	backupPath := filepath.Join(stma.tokensDir, "grid-token.json")

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
	var backup struct {
		Token     string    `json:"token"`
		CreatedAt time.Time `json:"created_at"`
		Version   string    `json:"version"`
		Checksum  string    `json:"checksum"`
	}

	if err := json.Unmarshal(data, &backup); err != nil {
		return "", fmt.Errorf("failed to unmarshal token backup: %w", err)
	}

	// Validar checksum (mesma lógica do SetupManager)
	expectedChecksum := stma.calculateChecksum(backup.Token)
	if backup.Checksum != expectedChecksum {
		return "", fmt.Errorf("token backup checksum validation failed")
	}

	return backup.Token, nil
}

// tokenExistsInFile verifica se token existe em arquivo (mesma lógica do SetupManager)
func (stma *SetupTokenManagerAdapter) tokenExistsInFile() (bool, error) {
	backupPath := filepath.Join(stma.tokensDir, "grid-token.json")
	_, err := os.Stat(backupPath)
	return !os.IsNotExist(err), nil
}

// calculateChecksum calcula checksum SHA256 do token (mesma lógica do SetupManager)
func (stma *SetupTokenManagerAdapter) calculateChecksum(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// CreateTokenIntegration creates a TokenIntegration instance using the Setup Component
func (sa *SetupAdapter) CreateTokenIntegration() types.TokenIntegration {
	return NewTokenIntegrationInstance()
}

// InitializeTokenIntegration initializes the token integration with Setup Component
func (sa *SetupAdapter) InitializeTokenIntegration(tokenIntegration types.TokenIntegration) error {
	// Create the adapter
	tokenManagerAdapter := NewSetupTokenManagerAdapter(sa)

	// Initialize the token integration
	if initTokenIntegration, ok := tokenIntegration.(*TokenIntegration); ok {
		return initTokenIntegration.Initialize(tokenManagerAdapter, sa.logger)
	}

	return fmt.Errorf("invalid token integration type")
}
