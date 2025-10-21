package node

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/node/src/internal/constants"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/node/src/internal/helpers"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/node/src/internal/types"
)

// AutoConfigGenerator generates automatic configurations for nodes
type AutoConfigGenerator struct {
	tokenIntegration types.TokenIntegration
	logger           types.Logger
	nodeConfigDir    string
}

// NewAutoConfigGenerator creates a new auto config generator
func NewAutoConfigGenerator(tokenIntegration types.TokenIntegration, logger types.Logger) types.ConfigGenerator {
	nodeConfigDir, err := helpers.ExpandPath(constants.DefaultNodeConfigDir)
	if err != nil {
		logger.Error("Failed to expand node config directory", "error", err)
		nodeConfigDir = constants.DefaultNodeConfigDir
	}

	return &AutoConfigGenerator{
		tokenIntegration: tokenIntegration,
		logger:           logger,
		nodeConfigDir:    nodeConfigDir,
	}
}

// GenerateNodeConfig generates a complete node configuration
func (acg *AutoConfigGenerator) GenerateNodeConfig() (*types.NodeConfig, error) {
	acg.logger.Info("Generating node configuration...")

	// Generate NodeID
	nodeID, err := acg.GenerateNodeID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate node ID: %w", err)
	}

	// Get SSH public key from Setup Component (owner.key)
	sshKeyProvider := NewSSHKeyProvider(acg.logger)
	sshPublicKey, err := sshKeyProvider.GetSSHPublicKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get SSH public key from Setup Component: %w", err)
	}

	// Generate node certificate
	nodeCert, err := acg.GenerateNodeCertificate(nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate node certificate: %w", err)
	}

	// Get Grid Token from Setup Component
	gridToken, err := acg.tokenIntegration.GetGridToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get grid token: %w", err)
	}

	// Detect Command Station IP
	commandStationIP, err := acg.DetectCommandStationIP()
	if err != nil {
		return nil, fmt.Errorf("failed to detect command station IP: %w", err)
	}

	// Create NodeConfig
	config := &types.NodeConfig{
		NodeID:           nodeID,
		GridToken:        gridToken,
		SSHPublicKey:     sshPublicKey,
		SSHPrivateKey:    "", // No longer used - private key stays in Setup Component
		NodeCertificate:  nodeCert,
		CommandStationIP: commandStationIP,
		CreatedAt:        time.Now(),
		ExpiresAt:        time.Now().Add(constants.DefaultTokenExpiry),
	}

	// Validate configuration
	if err := acg.ValidateConfig(config); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	// Save configuration to file
	if err := acg.saveNodeConfig(config); err != nil {
		acg.logger.Warn("Failed to save node configuration", "error", err)
		// Don't fail the entire operation, just log the warning
	}

	acg.logger.Info("Node configuration generated successfully", "nodeID", nodeID)
	return config, nil
}

// GenerateNodeID generates a sequential node ID
func (acg *AutoConfigGenerator) GenerateNodeID() (string, error) {
	// Get existing node IDs
	existingNodes, err := acg.getExistingNodeIDs()
	if err != nil {
		acg.logger.Warn("Failed to get existing node IDs", "error", err)
		// Continue with empty list if we can't read existing nodes
		existingNodes = []string{}
	}

	// Generate new node ID
	nodeID := helpers.GenerateNodeID(existingNodes)

	// Validate the generated node ID
	if err := helpers.ValidateNodeID(nodeID); err != nil {
		return "", fmt.Errorf("generated invalid node ID: %w", err)
	}

	acg.logger.Debug("Generated node ID", "nodeID", nodeID)
	return nodeID, nil
}

// GenerateSSHKeys generates a pair of SSH keys
func (acg *AutoConfigGenerator) GenerateSSHKeys() (interface{}, error) {
	acg.logger.Debug("Generating SSH keys...")

	sshKeys, err := helpers.GenerateSSHKeys()
	if err != nil {
		return nil, fmt.Errorf("failed to generate SSH keys: %w", err)
	}

	acg.logger.Debug("SSH keys generated successfully")
	return sshKeys, nil
}

// GenerateNodeCertificate generates a node certificate
func (acg *AutoConfigGenerator) GenerateNodeCertificate(nodeID string) (string, error) {
	acg.logger.Debug("Generating node certificate", "nodeID", nodeID)

	cert, err := helpers.GenerateNodeCertificate(nodeID)
	if err != nil {
		return "", fmt.Errorf("failed to generate node certificate: %w", err)
	}

	acg.logger.Debug("Node certificate generated successfully", "nodeID", nodeID)
	return cert, nil
}

// DetectCommandStationIP detects the local IP address of the command station
func (acg *AutoConfigGenerator) DetectCommandStationIP() (string, error) {
	acg.logger.Debug("Detecting command station IP...")

	ip, err := helpers.DetectCommandStationIP()
	if err != nil {
		return "", fmt.Errorf("failed to detect command station IP: %w", err)
	}

	acg.logger.Debug("Command station IP detected", "ip", ip)
	return ip, nil
}

// ValidateConfig validates a node configuration
func (acg *AutoConfigGenerator) ValidateConfig(config *types.NodeConfig) error {
	// Validate NodeID
	if err := helpers.ValidateNodeID(config.NodeID); err != nil {
		return fmt.Errorf("invalid node ID: %w", err)
	}

	// Validate Grid Token
	if err := acg.tokenIntegration.ValidateToken(config.GridToken); err != nil {
		return fmt.Errorf("invalid grid token: %w", err)
	}

	// Validate SSH public key
	if config.SSHPublicKey == "" {
		return fmt.Errorf("SSH public key cannot be empty")
	}

	// Validate node certificate
	if config.NodeCertificate == "" {
		return fmt.Errorf("node certificate cannot be empty")
	}

	// Validate Command Station IP
	if config.CommandStationIP == "" {
		return fmt.Errorf("command station IP cannot be empty")
	}

	// Validate timestamps
	if config.CreatedAt.IsZero() {
		return fmt.Errorf("created timestamp cannot be zero")
	}

	if config.ExpiresAt.IsZero() {
		return fmt.Errorf("expires timestamp cannot be zero")
	}

	if config.ExpiresAt.Before(config.CreatedAt) {
		return fmt.Errorf("expires timestamp must be after created timestamp")
	}

	acg.logger.Debug("Configuration validation passed")
	return nil
}

// Private helper methods

// getExistingNodeIDs gets a list of existing node IDs from the filesystem
func (acg *AutoConfigGenerator) getExistingNodeIDs() ([]string, error) {
	// Ensure node config directory exists
	if err := helpers.EnsureDirectory(acg.nodeConfigDir); err != nil {
		return nil, fmt.Errorf("failed to ensure node config directory: %w", err)
	}

	// Read directory contents
	entries, err := os.ReadDir(acg.nodeConfigDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read node config directory: %w", err)
	}

	var nodeIDs []string
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "node-") {
			nodeIDs = append(nodeIDs, entry.Name())
		}
	}

	// Sort node IDs for consistent ordering
	sort.Strings(nodeIDs)

	acg.logger.Debug("Found existing node IDs", "count", len(nodeIDs), "nodes", nodeIDs)
	return nodeIDs, nil
}

// saveNodeConfig saves the node configuration to a file
func (acg *AutoConfigGenerator) saveNodeConfig(config *types.NodeConfig) error {
	// Create node directory
	nodeDir := filepath.Join(acg.nodeConfigDir, config.NodeID)
	if err := helpers.EnsureDirectory(nodeDir); err != nil {
		return fmt.Errorf("failed to create node directory: %w", err)
	}

	// Save main configuration file
	configFile := filepath.Join(nodeDir, constants.NodeConfigFile)
	configData := fmt.Sprintf(`{
  "node_id": "%s",
  "grid_token": "%s",
  "ssh_public_key": "%s",
  "ssh_private_key": "%s",
  "node_certificate": "%s",
  "command_station_ip": "%s",
  "created_at": "%s",
  "expires_at": "%s"
}`,
		config.NodeID,
		config.GridToken,
		config.SSHPublicKey,
		config.SSHPrivateKey,
		config.NodeCertificate,
		config.CommandStationIP,
		config.CreatedAt.Format(time.RFC3339),
		config.ExpiresAt.Format(time.RFC3339),
	)

	if err := helpers.WriteFileWithPermissions(configFile, []byte(configData), constants.FileModeConfigFile); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// SSH private key is no longer stored per node - it stays in Setup Component
	// Only the public key is used for node configuration

	// Save certificate file separately
	certFile := filepath.Join(nodeDir, constants.NodeCertFile)
	if err := helpers.WriteFileWithPermissions(certFile, []byte(config.NodeCertificate), constants.FileModeConfigFile); err != nil {
		return fmt.Errorf("failed to write certificate file: %w", err)
	}

	acg.logger.Debug("Node configuration saved", "nodeID", config.NodeID, "directory", nodeDir)
	return nil
}

// LoadNodeConfig loads a node configuration from file
func (acg *AutoConfigGenerator) LoadNodeConfig(nodeID string) (*types.NodeConfig, error) {
	configFile := filepath.Join(acg.nodeConfigDir, nodeID, constants.NodeConfigFile)

	if !helpers.FileExists(configFile) {
		return nil, fmt.Errorf("node configuration file not found: %s", configFile)
	}

	configData, err := helpers.ReadFileSafely(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse JSON configuration
	// This is a simplified implementation - in production, use proper JSON unmarshaling
	// For now, we'll create a basic config structure
	_ = configData // Use configData to avoid unused variable warning

	config := &types.NodeConfig{
		NodeID:    nodeID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(constants.DefaultTokenExpiry),
	}

	// SSH private key is no longer stored per node - it stays in Setup Component
	// Load SSH public key from Setup Component instead
	sshKeyProvider := NewSSHKeyProvider(acg.logger)
	if sshPublicKey, err := sshKeyProvider.GetSSHPublicKey(); err == nil {
		config.SSHPublicKey = sshPublicKey
	}

	certFile := filepath.Join(acg.nodeConfigDir, nodeID, constants.NodeCertFile)
	if helpers.FileExists(certFile) {
		certData, err := helpers.ReadFileSafely(certFile)
		if err == nil {
			config.NodeCertificate = string(certData)
		}
	}

	return config, nil
}

// DeleteNodeConfig deletes a node configuration
func (acg *AutoConfigGenerator) DeleteNodeConfig(nodeID string) error {
	nodeDir := filepath.Join(acg.nodeConfigDir, nodeID)

	if !helpers.FileExists(nodeDir) {
		acg.logger.Warn("Node directory does not exist", "nodeID", nodeID, "directory", nodeDir)
		return nil
	}

	if err := os.RemoveAll(nodeDir); err != nil {
		return fmt.Errorf("failed to delete node directory: %w", err)
	}

	acg.logger.Debug("Node configuration deleted", "nodeID", nodeID, "directory", nodeDir)
	return nil
}

// ListNodeConfigs lists all available node configurations
func (acg *AutoConfigGenerator) ListNodeConfigs() ([]string, error) {
	return acg.getExistingNodeIDs()
}

// GetNextNodeID gets the next available node ID without creating a configuration
func (acg *AutoConfigGenerator) GetNextNodeID() (string, error) {
	existingNodes, err := acg.getExistingNodeIDs()
	if err != nil {
		acg.logger.Warn("Failed to get existing node IDs", "error", err)
		existingNodes = []string{}
	}

	return helpers.GenerateNodeID(existingNodes), nil
}

// ValidateNodeIDFormat validates if a node ID follows the correct format
func (acg *AutoConfigGenerator) ValidateNodeIDFormat(nodeID string) error {
	return helpers.ValidateNodeID(nodeID)
}

// GetNodeIDSequenceNumber extracts the sequence number from a node ID
func (acg *AutoConfigGenerator) GetNodeIDSequenceNumber(nodeID string) (int, error) {
	if err := helpers.ValidateNodeID(nodeID); err != nil {
		return 0, fmt.Errorf("invalid node ID format: %w", err)
	}

	// Extract number from "node-XX" format
	numStr := strings.TrimPrefix(nodeID, "node-")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse node ID number: %w", err)
	}

	return num, nil
}
