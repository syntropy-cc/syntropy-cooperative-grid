package node

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/node/src/internal/constants"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/node/src/internal/helpers"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/node/src/internal/types"
	"golang.org/x/crypto/pbkdf2"
)

// CloudInitGenerator generates cloud-init configurations for nodes
type CloudInitGenerator struct {
	logger       types.Logger
	templatesDir string
	outputDir    string
}

// CloudInitTemplateData contains data for cloud-init templates
type CloudInitTemplateData struct {
	NodeID             string
	GridTokenEncrypted string // CHANGED: was GridToken
	SSHPublicKey       string
	CommandStationIP   string
	NodeCertificate    string
	CreatedAt          string
	InstanceID         string
	Hostname           string
}

// NewCloudInitGenerator creates a new cloud-init generator
func NewCloudInitGenerator(logger types.Logger) types.CloudInitGenerator {
	templatesDir := filepath.Join(constants.DefaultTemplatesDir, "cloud-init")
	outputDir := constants.DefaultCloudInitOutputDir

	return &CloudInitGenerator{
		logger:       logger,
		templatesDir: templatesDir,
		outputDir:    outputDir,
	}
}

// GenerateCloudInit generates complete cloud-init configuration for a node
func (cig *CloudInitGenerator) GenerateCloudInit(config *types.NodeConfig) (*types.CloudInitConfig, error) {
	cig.logger.Info("Generating cloud-init configuration", "nodeID", config.NodeID)

	// Criptografar token antes de incluir no template
	encryptedToken, err := cig.encryptGridToken(config.GridToken, config.NodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt grid token: %w", err)
	}

	// Prepare template data with default options
	templateData := &CloudInitTemplateData{
		NodeID:             config.NodeID,
		GridTokenEncrypted: encryptedToken.Ciphertext,
		SSHPublicKey:       config.SSHPublicKey,
		CommandStationIP:   config.CommandStationIP,
		NodeCertificate:    config.NodeCertificate,
		CreatedAt:          config.CreatedAt.Format(time.RFC3339),
		InstanceID:         fmt.Sprintf("syntropy-node-%s", config.NodeID),
		Hostname:           config.NodeID,
	}

	// Generate user-data.yaml
	userData, err := cig.generateUserData(templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate user-data: %w", err)
	}

	// Generate network-config.yaml
	networkConfig, err := cig.generateNetworkConfig(templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate network-config: %w", err)
	}

	// Generate meta-data.yaml
	metaData, err := cig.generateMetaData(templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate meta-data: %w", err)
	}

	// Create output directory
	if err := helpers.EnsureDirectory(cig.outputDir); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	// Save cloud-init files
	nodeDir := filepath.Join(cig.outputDir, config.NodeID)
	if err := helpers.EnsureDirectory(nodeDir); err != nil {
		return nil, fmt.Errorf("failed to create node directory: %w", err)
	}

	userDataFile := filepath.Join(nodeDir, "user-data")
	if err := helpers.WriteFileWithPermissions(userDataFile, []byte(userData), constants.FileModeConfigFile); err != nil {
		return nil, fmt.Errorf("failed to write user-data: %w", err)
	}

	networkConfigFile := filepath.Join(nodeDir, "network-config")
	if err := helpers.WriteFileWithPermissions(networkConfigFile, []byte(networkConfig), constants.FileModeConfigFile); err != nil {
		return nil, fmt.Errorf("failed to write network-config: %w", err)
	}

	metaDataFile := filepath.Join(nodeDir, "meta-data")
	if err := helpers.WriteFileWithPermissions(metaDataFile, []byte(metaData), constants.FileModeConfigFile); err != nil {
		return nil, fmt.Errorf("failed to write meta-data: %w", err)
	}

	// Create cloud-init config object
	cloudInitConfig := &types.CloudInitConfig{
		NodeID:            config.NodeID,
		UserData:          userData,
		NetworkConfig:     networkConfig,
		MetaData:          metaData,
		UserDataFile:      userDataFile,
		NetworkConfigFile: networkConfigFile,
		MetaDataFile:      metaDataFile,
		CreatedAt:         time.Now(),
	}

	cig.logger.Info("Cloud-init configuration generated successfully", "nodeID", config.NodeID)
	return cloudInitConfig, nil
}

// generateUserData generates the user-data.yaml file
func (cig *CloudInitGenerator) generateUserData(data *CloudInitTemplateData) (string, error) {
	// Caminho para template
	templatePath := filepath.Join("infrastructure", "cloud-init", "user-data-template.yaml")

	// Fallback para custom template
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		homeDir, _ := os.UserHomeDir()
		templatePath = filepath.Join(homeDir, ".syntropy", "templates", "cloud-init", "user-data.yaml")
	}

	// Ler template
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read user-data template: %w", err)
	}

	// Processar template
	tmpl, err := template.New("user-data").Parse(string(templateContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse user-data template: %w", err)
	}

	var output strings.Builder
	if err := tmpl.Execute(&output, data); err != nil {
		return "", fmt.Errorf("failed to execute user-data template: %w", err)
	}

	return output.String(), nil
}

// encryptGridToken encrypts token with node-specific key
func (cig *CloudInitGenerator) encryptGridToken(token string, nodeID string) (*types.EncryptedTokenData, error) {
	// Gerar chave derivada do Node Certificate
	nodeSeed := sha256.Sum256([]byte(nodeID))
	nodeKey := pbkdf2.Key(nodeSeed[:], []byte("syntropy-grid"), 50000, 32, sha256.New)

	// Criptografar token
	block, err := aes.NewCipher(nodeKey)
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

	ciphertext := gcm.Seal(nonce, nonce, []byte(token), nil)

	// Retornar dados estruturados
	return &types.EncryptedTokenData{
		Ciphertext:  base64.StdEncoding.EncodeToString(ciphertext),
		NodeID:      nodeID,
		EncryptedAt: time.Now(),
	}, nil
}

// generateNetworkConfig generates the network-config.yaml file
func (cig *CloudInitGenerator) generateNetworkConfig(data *CloudInitTemplateData) (string, error) {
	// Caminho para template
	templatePath := filepath.Join("infrastructure", "cloud-init", "network-config-template.yaml")

	// Fallback para custom template
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		homeDir, _ := os.UserHomeDir()
		templatePath = filepath.Join(homeDir, ".syntropy", "templates", "cloud-init", "network-config.yaml")
	}

	// Ler template
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read network-config template: %w", err)
	}

	// Processar template
	tmpl, err := template.New("network-config").Parse(string(templateContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse network-config template: %w", err)
	}

	var output strings.Builder
	if err := tmpl.Execute(&output, data); err != nil {
		return "", fmt.Errorf("failed to execute network-config template: %w", err)
	}

	return output.String(), nil
}

// generateMetaData generates the meta-data.yaml file
func (cig *CloudInitGenerator) generateMetaData(data *CloudInitTemplateData) (string, error) {
	// Caminho para template
	templatePath := filepath.Join("infrastructure", "cloud-init", "meta-data-template.yaml")

	// Fallback para custom template
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		homeDir, _ := os.UserHomeDir()
		templatePath = filepath.Join(homeDir, ".syntropy", "templates", "cloud-init", "meta-data.yaml")
	}

	// Ler template
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read meta-data template: %w", err)
	}

	// Processar template
	tmpl, err := template.New("meta-data").Parse(string(templateContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse meta-data template: %w", err)
	}

	var output strings.Builder
	if err := tmpl.Execute(&output, data); err != nil {
		return "", fmt.Errorf("failed to execute meta-data template: %w", err)
	}

	return output.String(), nil
}

// ValidateCloudInit validates cloud-init configuration
func (cig *CloudInitGenerator) ValidateCloudInit(config *types.CloudInitConfig) error {
	// Validate user-data
	if config.UserData == "" {
		return fmt.Errorf("user-data cannot be empty")
	}

	// Validate network-config
	if config.NetworkConfig == "" {
		return fmt.Errorf("network-config cannot be empty")
	}

	// Validate meta-data
	if config.MetaData == "" {
		return fmt.Errorf("meta-data cannot be empty")
	}

	// Validate file paths
	if config.UserDataFile == "" || config.NetworkConfigFile == "" || config.MetaDataFile == "" {
		return fmt.Errorf("cloud-init file paths cannot be empty")
	}

	// Check if files exist
	if !helpers.FileExists(config.UserDataFile) {
		return fmt.Errorf("user-data file does not exist: %s", config.UserDataFile)
	}

	if !helpers.FileExists(config.NetworkConfigFile) {
		return fmt.Errorf("network-config file does not exist: %s", config.NetworkConfigFile)
	}

	if !helpers.FileExists(config.MetaDataFile) {
		return fmt.Errorf("meta-data file does not exist: %s", config.MetaDataFile)
	}

	cig.logger.Debug("Cloud-init configuration validation passed")
	return nil
}

// LoadCloudInit loads cloud-init configuration from files
func (cig *CloudInitGenerator) LoadCloudInit(nodeID string) (*types.CloudInitConfig, error) {
	nodeDir := filepath.Join(cig.outputDir, nodeID)

	// Check if directory exists
	if !helpers.FileExists(nodeDir) {
		return nil, fmt.Errorf("cloud-init directory not found for node: %s", nodeID)
	}

	// Load user-data
	userDataFile := filepath.Join(nodeDir, "user-data")
	userData, err := helpers.ReadFileSafely(userDataFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read user-data: %w", err)
	}

	// Load network-config
	networkConfigFile := filepath.Join(nodeDir, "network-config")
	networkConfig, err := helpers.ReadFileSafely(networkConfigFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read network-config: %w", err)
	}

	// Load meta-data
	metaDataFile := filepath.Join(nodeDir, "meta-data")
	metaData, err := helpers.ReadFileSafely(metaDataFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read meta-data: %w", err)
	}

	// Get file modification time
	stat, err := os.Stat(userDataFile)
	if err != nil {
		return nil, fmt.Errorf("failed to get file stats: %w", err)
	}

	cloudInitConfig := &types.CloudInitConfig{
		NodeID:            nodeID,
		UserData:          string(userData),
		NetworkConfig:     string(networkConfig),
		MetaData:          string(metaData),
		UserDataFile:      userDataFile,
		NetworkConfigFile: networkConfigFile,
		MetaDataFile:      metaDataFile,
		CreatedAt:         stat.ModTime(),
	}

	return cloudInitConfig, nil
}

// DeleteCloudInit deletes cloud-init configuration for a node
func (cig *CloudInitGenerator) DeleteCloudInit(nodeID string) error {
	nodeDir := filepath.Join(cig.outputDir, nodeID)

	if !helpers.FileExists(nodeDir) {
		cig.logger.Warn("Cloud-init directory does not exist", "nodeID", nodeID, "directory", nodeDir)
		return nil
	}

	if err := os.RemoveAll(nodeDir); err != nil {
		return fmt.Errorf("failed to delete cloud-init directory: %w", err)
	}

	cig.logger.Debug("Cloud-init configuration deleted", "nodeID", nodeID, "directory", nodeDir)
	return nil
}

// ListCloudInitConfigs lists all available cloud-init configurations
func (cig *CloudInitGenerator) ListCloudInitConfigs() ([]string, error) {
	if !helpers.FileExists(cig.outputDir) {
		return []string{}, nil
	}

	entries, err := os.ReadDir(cig.outputDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read cloud-init output directory: %w", err)
	}

	var nodeIDs []string
	for _, entry := range entries {
		if entry.IsDir() {
			nodeIDs = append(nodeIDs, entry.Name())
		}
	}

	return nodeIDs, nil
}

// GetCloudInitTemplate returns the template content for a specific file
func (cig *CloudInitGenerator) GetCloudInitTemplate(templateName string) (string, error) {
	switch templateName {
	case "user-data":
		return constants.UserDataTemplate, nil
	case "network-config":
		return constants.NetworkConfigTemplate, nil
	case "meta-data":
		return constants.MetaDataTemplate, nil
	default:
		return "", fmt.Errorf("unknown template: %s", templateName)
	}
}

// UpdateCloudInitTemplate updates a cloud-init template
func (cig *CloudInitGenerator) UpdateCloudInitTemplate(templateName, content string) error {
	// In a real implementation, this would update the template file
	// For now, we'll just validate the content
	if content == "" {
		return fmt.Errorf("template content cannot be empty")
	}

	// Basic validation - check if it's valid YAML-like content
	if !strings.Contains(content, ":") {
		return fmt.Errorf("template content appears to be invalid")
	}

	cig.logger.Debug("Template updated", "template", templateName)
	return nil
}

// ValidateTemplate validates a cloud-init template
func (cig *CloudInitGenerator) ValidateTemplate(templateName, content string) error {
	if content == "" {
		return fmt.Errorf("template content cannot be empty")
	}

	// Try to parse as template
	_, err := template.New(templateName).Parse(content)
	if err != nil {
		return fmt.Errorf("invalid template syntax: %w", err)
	}

	// Basic content validation
	switch templateName {
	case "user-data":
		if !strings.Contains(content, "#cloud-config") {
			return fmt.Errorf("user-data template must start with #cloud-config")
		}
	case "network-config":
		if !strings.Contains(content, "version") {
			return fmt.Errorf("network-config template must contain version field")
		}
	case "meta-data":
		if !strings.Contains(content, "instance-id") {
			return fmt.Errorf("meta-data template must contain instance-id field")
		}
	}

	return nil
}

// GetCloudInitStats returns statistics about cloud-init configurations
func (cig *CloudInitGenerator) GetCloudInitStats() (*types.CloudInitStats, error) {
	configs, err := cig.ListCloudInitConfigs()
	if err != nil {
		return nil, fmt.Errorf("failed to list cloud-init configs: %w", err)
	}

	stats := &types.CloudInitStats{
		TotalConfigs: len(configs),
		NodeIDs:      configs,
		OutputDir:    cig.outputDir,
		TemplatesDir: cig.templatesDir,
	}

	// Calculate total size
	var totalSize int64
	for _, nodeID := range configs {
		nodeDir := filepath.Join(cig.outputDir, nodeID)
		if helpers.FileExists(nodeDir) {
			// Get directory size
			size, err := helpers.GetDirectorySize(nodeDir)
			if err == nil {
				totalSize += size
			}
		}
	}

	stats.TotalSize = totalSize
	stats.LastUpdated = time.Now()

	return stats, nil
}

// GenerateNetworkConfig generates network configuration
func (cig *CloudInitGenerator) GenerateNetworkConfig() (string, error) {
	templateData := &CloudInitTemplateData{
		NodeID:             "default",
		GridTokenEncrypted: "",
		SSHPublicKey:       "",
		CommandStationIP:   "",
		NodeCertificate:    "",
		CreatedAt:          time.Now().Format(time.RFC3339),
		InstanceID:         "default",
		Hostname:           "default",
	}

	return cig.generateNetworkConfig(templateData)
}

// LoadTemplate loads a template by name
func (cig *CloudInitGenerator) LoadTemplate(templateName string) (string, error) {
	return cig.GetCloudInitTemplate(templateName)
}

// GenerateMetaData generates meta-data for a node
func (cig *CloudInitGenerator) GenerateMetaData(nodeID string) (string, error) {
	templateData := &CloudInitTemplateData{
		NodeID:     nodeID,
		InstanceID: fmt.Sprintf("syntropy-node-%s", nodeID),
		Hostname:   nodeID,
	}

	return cig.generateMetaData(templateData)
}

// ExecuteTemplate executes a template with variables
func (cig *CloudInitGenerator) ExecuteTemplate(templateContent string, variables map[string]interface{}) (string, error) {
	tmpl, err := template.New("cloud-init-template").Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var output strings.Builder
	if err := tmpl.Execute(&output, variables); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return output.String(), nil
}

// CleanupOldConfigs removes cloud-init configurations older than specified duration
func (cig *CloudInitGenerator) CleanupOldConfigs(maxAge time.Duration) error {
	configs, err := cig.ListCloudInitConfigs()
	if err != nil {
		return fmt.Errorf("failed to list cloud-init configs: %w", err)
	}

	var cleanedCount int
	cutoffTime := time.Now().Add(-maxAge)

	for _, nodeID := range configs {
		nodeDir := filepath.Join(cig.outputDir, nodeID)

		// Get directory modification time
		stat, err := os.Stat(nodeDir)
		if err != nil {
			cig.logger.Warn("Failed to get directory stats", "nodeID", nodeID, "error", err)
			continue
		}

		// Check if directory is older than cutoff
		if stat.ModTime().Before(cutoffTime) {
			if err := os.RemoveAll(nodeDir); err != nil {
				cig.logger.Warn("Failed to delete old cloud-init config", "nodeID", nodeID, "error", err)
				continue
			}

			cleanedCount++
			cig.logger.Debug("Cleaned up old cloud-init config", "nodeID", nodeID)
		}
	}

	cig.logger.Info("Cloud-init cleanup completed", "cleanedCount", cleanedCount, "maxAge", maxAge)
	return nil
}
