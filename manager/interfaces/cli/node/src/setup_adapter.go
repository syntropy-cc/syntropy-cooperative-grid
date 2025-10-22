package node

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"node-component/src/internal/types"
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
	adapter *SetupAdapter
}

// NewSetupTokenManagerAdapter creates a new token manager adapter
func NewSetupTokenManagerAdapter(adapter *SetupAdapter) *SetupTokenManagerAdapter {
	return &SetupTokenManagerAdapter{
		adapter: adapter,
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
	// Check if setup was run by looking for token in keyring or file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	// Try to load from keyring first (preferred method)
	token, err := stma.loadTokenFromKeyring()
	if err == nil && token != "" {
		return token, nil
	}

	// Fallback: ler grid-token.json (formato do Setup)
	tokenFile := filepath.Join(homeDir, ".syntropy", "tokens", "grid-token.json")
	if _, err := os.Stat(tokenFile); err == nil {
		tokenBytes, err := os.ReadFile(tokenFile)
		if err != nil {
			return "", fmt.Errorf("failed to read token file: %w", err)
		}

		// Parse JSON estrutura TokenBackup
		var backup struct {
			Token     string    `json:"token"`
			CreatedAt time.Time `json:"created_at"`
			Version   string    `json:"version"`
			Checksum  string    `json:"checksum"`
		}

		if err := json.Unmarshal(tokenBytes, &backup); err != nil {
			return "", fmt.Errorf("failed to parse token file: %w", err)
		}

		return backup.Token, nil
	}

	return "", fmt.Errorf("grid token not found - please run 'syntropy setup' first")
}

// DeleteToken deletes the token (delegates to Setup Component)
func (stma *SetupTokenManagerAdapter) DeleteToken() error {
	// This would normally call the Setup Component's TokenManager
	return fmt.Errorf("token deletion should be done via 'syntropy setup' command")
}

// TokenExists checks if token exists (delegates to Setup Component)
func (stma *SetupTokenManagerAdapter) TokenExists() (bool, error) {
	// Check if token exists by looking for the token file directly
	// This is more reliable than calling LoadToken() which might fail for other reasons
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false, fmt.Errorf("failed to get home directory: %w", err)
	}

	// Check if token file exists
	tokenFile := filepath.Join(homeDir, ".syntropy", "tokens", "grid-token.json")
	if _, err := os.Stat(tokenFile); err == nil {
		return true, nil
	}

	// Also check keyring (even though it's not implemented yet)
	// This ensures consistency with the Setup Component's approach
	_, err = stma.loadTokenFromKeyring()
	if err == nil {
		return true, nil
	}

	return false, nil
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

// loadTokenFromKeyring attempts to load token from system keyring
func (stma *SetupTokenManagerAdapter) loadTokenFromKeyring() (string, error) {
	// Import the keyring library for actual keyring access
	// For now, we'll implement a basic keyring check using the same approach as the Setup Component
	// This should use the same keyring service and user as defined in the Setup Component

	// Check if keyring is available by trying to access it
	// We'll use the same constants as the Setup Component
	const (
		KeyringService = "syntropy-grid"
		KeyringUser    = "grid-token"
	)

	// Try to load from keyring using the same approach as Setup Component
	// For now, we'll return an error to fallback to file-based storage
	// In a production implementation, this would use github.com/zalando/go-keyring
	// like: return keyring.Get(KeyringService, KeyringUser)
	return "", fmt.Errorf("keyring not available - using file fallback")
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
