package node

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/node/src/internal/types"
)

// SSHKeyProvider provides SSH key management for Node Component
type SSHKeyProvider struct {
	logger types.Logger
}

// NewSSHKeyProvider creates a new SSH key provider
func NewSSHKeyProvider(logger types.Logger) *SSHKeyProvider {
	return &SSHKeyProvider{
		logger: logger,
	}
}

// GetSSHPublicKey returns the SSH public key from Setup Component
func (skp *SSHKeyProvider) GetSSHPublicKey() (string, error) {
	skp.logger.Debug("Loading SSH public key from Setup Component")

	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	// Try to load from owner.key.pub file
	publicKeyPath := filepath.Join(homeDir, ".syntropy", "keys", "owner.key.pub")
	if _, err := os.Stat(publicKeyPath); err != nil {
		return "", fmt.Errorf("owner key not found - please run 'syntropy setup' first: %w", err)
	}

	// Read the public key
	publicKeyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return "", fmt.Errorf("failed to read owner key: %w", err)
	}

	publicKey := string(publicKeyBytes)

	// Ensure proper format for authorized_keys
	publicKey = skp.formatSSHPublicKey(publicKey)

	skp.logger.Debug("SSH public key loaded successfully", "key_preview", publicKey[:20]+"...")
	return publicKey, nil
}

// GetSSHPrivateKeyPath returns the path to the SSH private key
func (skp *SSHKeyProvider) GetSSHPrivateKeyPath() (string, error) {
	skp.logger.Debug("Getting SSH private key path from Setup Component")

	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	// Check if owner key exists
	privateKeyPath := filepath.Join(homeDir, ".syntropy", "keys", "owner.key")
	if _, err := os.Stat(privateKeyPath); err != nil {
		return "", fmt.Errorf("owner key not found - please run 'syntropy setup' first: %w", err)
	}

	skp.logger.Debug("SSH private key path retrieved", "path", privateKeyPath)
	return privateKeyPath, nil
}

// ValidateSSHKeys validates that SSH keys exist and are properly formatted
func (skp *SSHKeyProvider) ValidateSSHKeys() error {
	skp.logger.Debug("Validating SSH keys from Setup Component")

	// Check if public key exists and is readable
	publicKey, err := skp.GetSSHPublicKey()
	if err != nil {
		return fmt.Errorf("SSH public key validation failed: %w", err)
	}

	// Check if private key exists and is readable
	privateKeyPath, err := skp.GetSSHPrivateKeyPath()
	if err != nil {
		return fmt.Errorf("SSH private key validation failed: %w", err)
	}

	// Basic format validation
	if len(publicKey) < 50 {
		return fmt.Errorf("SSH public key appears to be too short")
	}

	if !skp.isValidSSHPublicKey(publicKey) {
		return fmt.Errorf("SSH public key format is invalid")
	}

	// Check private key file permissions (should be 600)
	stat, err := os.Stat(privateKeyPath)
	if err != nil {
		return fmt.Errorf("failed to stat private key: %w", err)
	}

	mode := stat.Mode()
	if mode&077 != 0 {
		skp.logger.Warn("SSH private key has overly permissive permissions", "mode", mode.String())
	}

	skp.logger.Debug("SSH keys validation passed")
	return nil
}

// formatSSHPublicKey formats the SSH public key for authorized_keys
func (skp *SSHKeyProvider) formatSSHPublicKey(publicKey string) string {
	// Remove any trailing newlines and whitespace
	publicKey = strings.TrimSpace(publicKey)

	// Ensure it ends with a newline
	if !strings.HasSuffix(publicKey, "\n") {
		publicKey = publicKey + "\n"
	}

	return publicKey
}

// isValidSSHPublicKey validates the format of an SSH public key
func (skp *SSHKeyProvider) isValidSSHPublicKey(publicKey string) bool {
	// Basic validation - check for common SSH key types
	validTypes := []string{
		"ssh-rsa",
		"ssh-ed25519",
		"ecdsa-sha2-",
		"ssh-dss",
	}

	publicKey = strings.TrimSpace(publicKey)

	for _, keyType := range validTypes {
		if strings.HasPrefix(publicKey, keyType) {
			return true
		}
	}

	return false
}

// GetSSHKeyInfo returns information about the SSH keys
func (skp *SSHKeyProvider) GetSSHKeyInfo() (*SSHKeyInfo, error) {
	skp.logger.Debug("Getting SSH key information")

	publicKey, err := skp.GetSSHPublicKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get public key: %w", err)
	}

	privateKeyPath, err := skp.GetSSHPrivateKeyPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get private key path: %w", err)
	}

	// Get file stats
	stat, err := os.Stat(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat private key: %w", err)
	}

	// Extract key type from public key
	keyType := "unknown"
	if strings.HasPrefix(publicKey, "ssh-rsa") {
		keyType = "rsa"
	} else if strings.HasPrefix(publicKey, "ssh-ed25519") {
		keyType = "ed25519"
	} else if strings.HasPrefix(publicKey, "ecdsa-sha2-") {
		keyType = "ecdsa"
	}

	info := &SSHKeyInfo{
		PublicKey:      publicKey,
		PrivateKeyPath: privateKeyPath,
		KeyType:        keyType,
		CreatedAt:      stat.ModTime(),
		Size:           stat.Size(),
	}

	skp.logger.Debug("SSH key information retrieved", "type", keyType, "size", stat.Size())
	return info, nil
}

// SSHKeyInfo contains information about SSH keys
type SSHKeyInfo struct {
	PublicKey      string    `json:"public_key"`
	PrivateKeyPath string    `json:"private_key_path"`
	KeyType        string    `json:"key_type"`
	CreatedAt      time.Time `json:"created_at"`
	Size           int64     `json:"size"`
}
