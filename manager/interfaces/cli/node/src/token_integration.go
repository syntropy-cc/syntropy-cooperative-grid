package node

import (
	"fmt"
	"time"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/node/src/internal/constants"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/node/src/internal/types"
	// Import Setup Component interfaces - TODO: Fix import path when Setup Component is available
	// setupTypes "github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/srgithub.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/node/src/internal/types"
)

// SetupTokenManager interface for integration with Setup Component
type SetupTokenManager interface {
	GenerateToken() (string, error)
	SaveToken(token string) error
	LoadToken() (string, error)
	DeleteToken() error
	TokenExists() (bool, error)
	RotateToken() (string, error)
	ExportToken(outputPath string) error
	ImportToken(inputPath string) error
	ValidateToken(token string) error
}

// TokenIntegration handles integration with the Setup Component's TokenManager
type TokenIntegration struct {
	setupTokenManager SetupTokenManager
	logger            types.Logger
	cache             *TokenCache
}

// TokenCache caches token information to avoid repeated keyring access
type TokenCache struct {
	token      string
	expiresAt  time.Time
	isValid    bool
	lastAccess time.Time
}

// NewTokenIntegration creates a new token integration instance
func NewTokenIntegrationInstance() types.TokenIntegration {
	return &TokenIntegration{
		cache: &TokenCache{
			isValid: false,
		},
	}
}

// Initialize initializes the token integration with Setup Component
func (ti *TokenIntegration) Initialize(setupTokenManager SetupTokenManager, logger types.Logger) error {
	ti.setupTokenManager = setupTokenManager
	ti.logger = logger

	// Verify token exists in Setup Component
	exists, err := ti.setupTokenManager.TokenExists()
	if err != nil {
		return fmt.Errorf("failed to check token existence: %w", err)
	}

	if !exists {
		return fmt.Errorf("grid token not found in setup component - please run 'syntropy setup' first")
	}

	ti.logger.Info("Token integration initialized successfully")
	return nil
}

// GetGridToken retrieves the Grid Token from the Setup Component's TokenManager
func (ti *TokenIntegration) GetGridToken() (string, error) {
	// Check cache first
	if ti.isTokenCacheValid() {
		ti.cache.lastAccess = time.Now()
		ti.logger.Debug("Using cached grid token")
		return ti.cache.token, nil
	}

	// Validate setup token manager is initialized
	if ti.setupTokenManager == nil {
		return "", fmt.Errorf("token integration not initialized - call Initialize() first")
	}

	// Load token from Setup Component
	token, err := ti.setupTokenManager.LoadToken()
	if err != nil {
		return "", fmt.Errorf("failed to load grid token from setup component: %w", err)
	}

	// Validate token format
	if err := ti.validateTokenFormat(token); err != nil {
		return "", fmt.Errorf("invalid token format: %w", err)
	}

	// Update cache
	ti.cache.token = token
	ti.cache.expiresAt = time.Now().Add(constants.DefaultTokenExpiry)
	ti.cache.isValid = true
	ti.cache.lastAccess = time.Now()

	ti.logger.Debug("Grid token loaded successfully from setup component")
	return token, nil
}

// ValidateToken validates a token against the Setup Component's TokenManager
func (ti *TokenIntegration) ValidateToken(token string) error {
	// Validate setup token manager is initialized
	if ti.setupTokenManager == nil {
		return fmt.Errorf("token integration not initialized - call Initialize() first")
	}

	// Validate token format first
	if err := ti.validateTokenFormat(token); err != nil {
		return fmt.Errorf("invalid token format: %w", err)
	}

	// Validate token using Setup Component
	if err := ti.setupTokenManager.ValidateToken(token); err != nil {
		return fmt.Errorf("token validation failed: %w", err)
	}

	ti.logger.Debug("Token validated successfully")
	return nil
}

// RefreshToken refreshes the token from the Setup Component
func (ti *TokenIntegration) RefreshToken() (string, error) {
	// Validate setup token manager is initialized
	if ti.setupTokenManager == nil {
		return "", fmt.Errorf("token integration not initialized - call Initialize() first")
	}

	// Check if token needs rotation
	needsRotation, err := ti.tokenNeedsRotation()
	if err != nil {
		return "", fmt.Errorf("failed to check token rotation: %w", err)
	}

	if needsRotation {
		// Rotate token using Setup Component
		newToken, err := ti.setupTokenManager.RotateToken()
		if err != nil {
			return "", fmt.Errorf("failed to rotate token: %w", err)
		}

		// Update cache
		ti.cache.token = newToken
		ti.cache.expiresAt = time.Now().Add(constants.DefaultTokenExpiry)
		ti.cache.isValid = true
		ti.cache.lastAccess = time.Now()

		ti.logger.Info("Token rotated successfully")
		return newToken, nil
	}

	// Return current token
	return ti.GetGridToken()
}

// GetTokenExpiry returns the token expiry time
func (ti *TokenIntegration) GetTokenExpiry() (time.Time, error) {
	if !ti.cache.isValid {
		// Load token to populate cache
		_, err := ti.GetGridToken()
		if err != nil {
			return time.Time{}, fmt.Errorf("failed to load token: %w", err)
		}
	}

	return ti.cache.expiresAt, nil
}

// IsTokenValid checks if the current token is valid
func (ti *TokenIntegration) IsTokenValid(token string) bool {
	// Quick format check
	if err := ti.validateTokenFormat(token); err != nil {
		return false
	}

	// Check if token matches cached token
	if ti.cache.isValid && ti.cache.token == token {
		return !ti.isTokenExpired()
	}

	// Validate against Setup Component
	if err := ti.ValidateToken(token); err != nil {
		return false
	}

	return true
}

// Private helper methods

// validateTokenFormat validates the basic format of a token
func (ti *TokenIntegration) validateTokenFormat(token string) error {
	if token == "" {
		return fmt.Errorf("token cannot be empty")
	}

	if len(token) < constants.MinTokenLength {
		return fmt.Errorf("token too short: minimum %d characters required", constants.MinTokenLength)
	}

	if len(token) > constants.MaxTokenLength {
		return fmt.Errorf("token too long: maximum %d characters allowed", constants.MaxTokenLength)
	}

	// Check if token contains only valid characters (alphanumeric + hyphens)
	for _, char := range token {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '-') {
			return fmt.Errorf("token contains invalid characters: only alphanumeric and hyphens allowed")
		}
	}

	return nil
}

// isTokenCacheValid checks if the cached token is still valid
func (ti *TokenIntegration) isTokenCacheValid() bool {
	if !ti.cache.isValid {
		return false
	}

	if ti.isTokenExpired() {
		ti.cache.isValid = false
		return false
	}

	// Check if cache is too old (refresh every hour)
	if time.Since(ti.cache.lastAccess) > time.Hour {
		return false
	}

	return true
}

// isTokenExpired checks if the token has expired
func (ti *TokenIntegration) isTokenExpired() bool {
	return time.Now().After(ti.cache.expiresAt)
}

// tokenNeedsRotation checks if the token needs to be rotated
func (ti *TokenIntegration) tokenNeedsRotation() (bool, error) {
	if !ti.cache.isValid {
		return false, nil
	}

	// Rotate if token expires within 1 hour
	rotationThreshold := time.Hour
	return time.Until(ti.cache.expiresAt) < rotationThreshold, nil
}

// ClearCache clears the token cache
func (ti *TokenIntegration) ClearCache() {
	ti.cache.isValid = false
	ti.cache.token = ""
	ti.cache.expiresAt = time.Time{}
	ti.cache.lastAccess = time.Time{}

	ti.logger.Debug("Token cache cleared")
}

// GetCacheInfo returns information about the current cache state
func (ti *TokenIntegration) GetCacheInfo() map[string]interface{} {
	return map[string]interface{}{
		"is_valid":      ti.cache.isValid,
		"expires_at":    ti.cache.expiresAt,
		"last_access":   ti.cache.lastAccess,
		"is_expired":    ti.isTokenExpired(),
		"needs_refresh": !ti.isTokenCacheValid(),
	}
}

// TokenIntegrationFactory creates a new token integration instance
type TokenIntegrationFactory struct{}

// NewTokenIntegrationFactory creates a new factory
func NewTokenIntegrationFactory() *TokenIntegrationFactory {
	return &TokenIntegrationFactory{}
}

// CreateTokenIntegration creates a token integration with Setup Component
func (tif *TokenIntegrationFactory) CreateTokenIntegration() (types.TokenIntegration, error) {
	// This would typically connect to the Setup Component
	// For now, return a mock implementation that will be replaced
	// when the Setup Component is properly integrated

	return &TokenIntegration{
		cache: &TokenCache{
			isValid: false,
		},
	}, nil
}

// MockTokenIntegration provides a mock implementation for testing
type MockTokenIntegration struct {
	token      string
	isValid    bool
	shouldFail bool
}

// NewMockTokenIntegration creates a mock token integration for testing
func NewMockTokenIntegration(token string, isValid bool) *MockTokenIntegration {
	return &MockTokenIntegration{
		token:      token,
		isValid:    isValid,
		shouldFail: false,
	}
}

// SetShouldFail configures the mock to fail operations
func (mti *MockTokenIntegration) SetShouldFail(shouldFail bool) {
	mti.shouldFail = shouldFail
}

// GetGridToken returns the mock token
func (mti *MockTokenIntegration) GetGridToken() (string, error) {
	if mti.shouldFail {
		return "", fmt.Errorf("mock token integration failure")
	}
	return mti.token, nil
}

// ValidateToken validates the mock token
func (mti *MockTokenIntegration) ValidateToken(token string) error {
	if mti.shouldFail {
		return fmt.Errorf("mock token validation failure")
	}

	if token == mti.token && mti.isValid {
		return nil
	}

	return fmt.Errorf("invalid token")
}

// RefreshToken returns the mock token
func (mti *MockTokenIntegration) RefreshToken() (string, error) {
	if mti.shouldFail {
		return "", fmt.Errorf("mock token refresh failure")
	}
	return mti.token, nil
}

// GetTokenExpiry returns a mock expiry time
func (mti *MockTokenIntegration) GetTokenExpiry() (time.Time, error) {
	if mti.shouldFail {
		return time.Time{}, fmt.Errorf("mock token expiry failure")
	}
	return time.Now().Add(constants.DefaultTokenExpiry), nil
}

// IsTokenValid checks if the mock token is valid
func (mti *MockTokenIntegration) IsTokenValid(token string) bool {
	return token == mti.token && mti.isValid && !mti.shouldFail
}
