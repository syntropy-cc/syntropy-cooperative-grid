package node

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"node-component/src/internal/types"

	"github.com/zalando/go-keyring"
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

	// Log debug para identificar o ambiente (apenas em modo debug)
	if adapter.logger != nil {
		adapter.logger.Debug("Creating SetupTokenManagerAdapter",
			"home_dir", homeDir,
			"tokens_dir", tokensDir,
			"os", runtime.GOOS,
			"arch", runtime.GOARCH)
	}

	// Garantir que o diretório seja criado com permissões corretas
	if err := os.MkdirAll(tokensDir, 0700); err != nil {
		adapter.logger.Warn("Failed to create tokens directory", "error", err, "path", tokensDir)
	}

	// Detectar se o keyring está disponível (mesmo que o Setup Component)
	keyringAvailable := isKeyringAvailable()

	return &SetupTokenManagerAdapter{
		adapter:          adapter,
		keyringAvailable: keyringAvailable,
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
	// Usar a mesma lógica do Setup Component: keyring primeiro, fallback para arquivo
	if !stma.keyringAvailable {
		return stma.loadTokenFromFile()
	}

	// Tentar carregar do keyring primeiro
	token, err := keyring.Get(KeyringService, KeyringUser)
	if err != nil {
		if err == keyring.ErrNotFound {
			// Log que está usando fallback para arquivo
			if stma.adapter.logger != nil {
				stma.adapter.logger.Debug("Token not found in keyring, using file fallback")
			}
			// Verificar se existe em arquivo como fallback
			return stma.loadTokenFromFile()
		}
		return "", fmt.Errorf("failed to load token from keyring: %w", err)
	}

	// Log que carregou do keyring
	if stma.adapter.logger != nil {
		stma.adapter.logger.Debug("Token loaded from keyring")
	}

	// Validar token carregado
	if err := stma.ValidateToken(token); err != nil {
		return "", fmt.Errorf("loaded token is invalid: %w", err)
	}

	return token, nil
}

// DeleteToken deletes the token (delegates to Setup Component)
func (stma *SetupTokenManagerAdapter) DeleteToken() error {
	// This would normally call the Setup Component's TokenManager
	return fmt.Errorf("token deletion should be done via 'syntropy setup' command")
}

// TokenExists checks if token exists (delegates to Setup Component)
func (stma *SetupTokenManagerAdapter) TokenExists() (bool, error) {
	// Usar a mesma lógica do Setup Component: keyring primeiro, fallback para arquivo
	if !stma.keyringAvailable {
		return stma.tokenExistsInFile()
	}

	// Verificar no keyring primeiro
	_, err := keyring.Get(KeyringService, KeyringUser)
	if err != nil {
		if err == keyring.ErrNotFound {
			// Verificar se existe em arquivo como fallback
			return stma.tokenExistsInFile()
		}
		return false, fmt.Errorf("failed to check token existence: %w", err)
	}

	return true, nil
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
		// Log detalhado para debug
		if stma.adapter.logger != nil {
			stma.adapter.logger.Warn("Token checksum validation failed",
				"expected", expectedChecksum,
				"actual", backup.Checksum,
				"token_preview", backup.Token[:8]+"...",
				"file_path", backupPath)
		}
		return "", fmt.Errorf("token backup checksum validation failed: expected %s, got %s", expectedChecksum, backup.Checksum)
	}

	// Log sucesso para debug
	if stma.adapter.logger != nil {
		stma.adapter.logger.Debug("Token loaded successfully from file",
			"token_preview", backup.Token[:8]+"...",
			"file_path", backupPath,
			"checksum_valid", true)
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

// isKeyringAvailable verifica se o keyring está disponível (mesma lógica do SetupManager)
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
