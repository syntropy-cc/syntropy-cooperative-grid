package node

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/syntropy-grid/syntropy-cooperative-grid/manager/interfaces/cli/node/src/internal/types"
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

	// Fallback to file-based storage
	tokenFile := filepath.Join(homeDir, ".syntropy", "tokens", "grid-token")
	if _, err := os.Stat(tokenFile); err == nil {
		tokenBytes, err := os.ReadFile(tokenFile)
		if err != nil {
			return "", fmt.Errorf("failed to read token file: %w", err)
		}
		return string(tokenBytes), nil
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
	_, err := stma.LoadToken()
	return err == nil, nil
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
	// This would use the keyring library to load from system keyring
	// For now, we'll return an error to indicate keyring is not available
	// In a real implementation, this would use github.com/zalando/go-keyring
	return "", fmt.Errorf("keyring not available")
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
