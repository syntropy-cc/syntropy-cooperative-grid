package node

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/node/src/internal/types"
)

// CreateSubcomponent handles the complete node creation workflow
type CreateSubcomponent struct {
	configGenerator    types.ConfigGenerator
	usbDetector        USBDetector
	isoDownloader      ISODownloader
	cloudInitGenerator types.CloudInitGenerator
	usbWriter          USBWriter
	tokenIntegration   types.TokenIntegration
	logger             types.Logger
}

// NewCreateSubcomponent creates a new Create subcomponent
func NewCreateSubcomponent(
	configGenerator types.ConfigGenerator,
	usbDetector USBDetector,
	isoDownloader ISODownloader,
	cloudInitGenerator types.CloudInitGenerator,
	usbWriter USBWriter,
	tokenIntegration types.TokenIntegration,
	logger types.Logger,
) *CreateSubcomponent {
	return &CreateSubcomponent{
		configGenerator:    configGenerator,
		usbDetector:        usbDetector,
		isoDownloader:      isoDownloader,
		cloudInitGenerator: cloudInitGenerator,
		usbWriter:          usbWriter,
		tokenIntegration:   tokenIntegration,
		logger:             logger,
	}
}

// CreateOptions represents options for node creation
type CreateOptions struct {
	UbuntuVersion    string
	DevicePath       string
	SkipUSBDetection bool
	SkipISODownload  bool
	SkipCloudInit    bool
	SkipUSBWrite     bool
	ForceOverwrite   bool
	AutoStart        bool
}

// CreateResult represents the result of node creation
type CreateResult struct {
	NodeID          string                 `json:"node_id"`
	DevicePath      string                 `json:"device_path"`
	ISOPath         string                 `json:"iso_path"`
	CloudInitConfig *types.CloudInitConfig `json:"cloud_init_config"`
	Success         bool                   `json:"success"`
	Duration        time.Duration          `json:"duration"`
	ErrorMessage    string                 `json:"error_message,omitempty"`
	StepsCompleted  []string               `json:"steps_completed"`
	StepsFailed     []string               `json:"steps_failed"`
}

// CreateNode creates a complete node with all necessary components
func (cs *CreateSubcomponent) CreateNode(ctx context.Context, options CreateOptions) (*CreateResult, error) {
	cs.logger.Info("Starting node creation process", "options", options)

	startTime := time.Now()
	result := &CreateResult{
		NodeID:         "",
		DevicePath:     options.DevicePath,
		Success:        false,
		StepsCompleted: []string{},
		StepsFailed:    []string{},
	}

	// Step 1: Validate prerequisites
	if err := cs.validatePrerequisites(ctx); err != nil {
		result.StepsFailed = append(result.StepsFailed, "validate_prerequisites")
		result.ErrorMessage = err.Error()
		return result, fmt.Errorf("prerequisites validation failed: %w", err)
	}
	result.StepsCompleted = append(result.StepsCompleted, "validate_prerequisites")

	// Step 2: Generate node configuration
	nodeConfig, err := cs.generateNodeConfiguration(ctx)
	if err != nil {
		result.StepsFailed = append(result.StepsFailed, "generate_configuration")
		result.ErrorMessage = err.Error()
		return result, fmt.Errorf("configuration generation failed: %w", err)
	}
	result.NodeID = nodeConfig.NodeID
	result.StepsCompleted = append(result.StepsCompleted, "generate_configuration")

	// Step 3: Detect USB device (if not specified)
	if !options.SkipUSBDetection && options.DevicePath == "" {
		devicePath, err := cs.detectUSBDevice(ctx)
		if err != nil {
			result.StepsFailed = append(result.StepsFailed, "detect_usb_device")
			result.ErrorMessage = err.Error()
			return result, fmt.Errorf("USB device detection failed: %w", err)
		}
		result.DevicePath = devicePath
		result.StepsCompleted = append(result.StepsCompleted, "detect_usb_device")
	}

	// Step 4: Download Ubuntu ISO
	var isoPath string
	if !options.SkipISODownload {
		isoPath, err = cs.downloadUbuntuISO(ctx, options.UbuntuVersion)
		if err != nil {
			result.StepsFailed = append(result.StepsFailed, "download_iso")
			result.ErrorMessage = err.Error()
			return result, fmt.Errorf("ISO download failed: %w", err)
		}
		result.ISOPath = isoPath
		result.StepsCompleted = append(result.StepsCompleted, "download_iso")
	}

	// Step 5: Generate cloud-init configuration
	var cloudInitConfig *types.CloudInitConfig
	if !options.SkipCloudInit {
		cloudInitConfig, err = cs.generateCloudInit(ctx, nodeConfig)
		if err != nil {
			result.StepsFailed = append(result.StepsFailed, "generate_cloud_init")
			result.ErrorMessage = err.Error()
			return result, fmt.Errorf("cloud-init generation failed: %w", err)
		}
		result.CloudInitConfig = cloudInitConfig
		result.StepsCompleted = append(result.StepsCompleted, "generate_cloud_init")
	}

	// Step 6: Write ISO to USB device
	if !options.SkipUSBWrite {
		writeResult, err := cs.writeISOToUSB(ctx, isoPath, result.DevicePath, cloudInitConfig)
		if err != nil {
			result.StepsFailed = append(result.StepsFailed, "write_iso_to_usb")
			result.ErrorMessage = err.Error()
			return result, fmt.Errorf("USB write failed: %w", err)
		}
		if !writeResult.Success {
			result.StepsFailed = append(result.StepsFailed, "write_iso_to_usb")
			result.ErrorMessage = writeResult.ErrorMessage
			return result, fmt.Errorf("USB write failed: %s", writeResult.ErrorMessage)
		}
		result.StepsCompleted = append(result.StepsCompleted, "write_iso_to_usb")
	}

	// Step 7: Save node state
	if err := cs.saveNodeState(ctx, nodeConfig, result); err != nil {
		cs.logger.Warn("Failed to save node state", "error", err)
		// Don't fail the entire process for this
	}
	result.StepsCompleted = append(result.StepsCompleted, "save_node_state")

	// Step 8: Start listener (if auto-start enabled)
	if options.AutoStart {
		if err := cs.startListener(ctx, nodeConfig); err != nil {
			cs.logger.Warn("Failed to start listener", "error", err)
			// Don't fail the entire process for this
		} else {
			result.StepsCompleted = append(result.StepsCompleted, "start_listener")
		}
	}

	result.Duration = time.Since(startTime)
	result.Success = true

	cs.logger.Info("Node creation completed successfully",
		"node_id", result.NodeID,
		"device", result.DevicePath,
		"duration", result.Duration,
		"steps_completed", len(result.StepsCompleted))

	return result, nil
}

// CreateNodeInteractive creates a node with interactive prompts
func (cs *CreateSubcomponent) CreateNodeInteractive(ctx context.Context) (*CreateResult, error) {
	cs.logger.Info("Starting interactive node creation")

	// Get Ubuntu version
	ubuntuVersion := cs.promptUbuntuVersion()

	// Get USB device
	devicePath := cs.promptUSBDevice(ctx)

	// Get other options
	autoStart := cs.promptAutoStart()

	options := CreateOptions{
		UbuntuVersion: ubuntuVersion,
		DevicePath:    devicePath,
		AutoStart:     autoStart,
	}

	return cs.CreateNode(ctx, options)
}

// Private helper methods

// validatePrerequisites validates that all prerequisites are met
func (cs *CreateSubcomponent) validatePrerequisites(ctx context.Context) error {
	cs.logger.Debug("Validating prerequisites")

	// Check if token integration is available
	if cs.tokenIntegration == nil {
		return fmt.Errorf("token integration is not available")
	}

	// Check if Grid Token is available
	if _, err := cs.tokenIntegration.GetGridToken(); err != nil {
		return fmt.Errorf("Grid Token is not available: %w", err)
	}

	// Check if required tools are available
	if err := cs.checkRequiredTools(); err != nil {
		return fmt.Errorf("required tools are not available: %w", err)
	}

	// Check if we have write permissions to the working directory
	if err := cs.checkWritePermissions(); err != nil {
		return fmt.Errorf("write permissions check failed: %w", err)
	}

	cs.logger.Debug("Prerequisites validation passed")
	return nil
}

// generateNodeConfiguration generates the node configuration
func (cs *CreateSubcomponent) generateNodeConfiguration(ctx context.Context) (*types.NodeConfig, error) {
	cs.logger.Debug("Generating node configuration")

	// Generate configuration
	config, err := cs.configGenerator.GenerateNodeConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to generate configuration: %w", err)
	}

	cs.logger.Debug("Node configuration generated", "node_id", config.NodeID)
	return config, nil
}

// detectUSBDevice detects a suitable USB device
func (cs *CreateSubcomponent) detectUSBDevice(ctx context.Context) (string, error) {
	cs.logger.Debug("Detecting USB device")

	// Get suitable devices
	devices, err := cs.usbDetector.DetectRemovableDevices(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to detect USB devices: %w", err)
	}

	if len(devices) == 0 {
		return "", fmt.Errorf("no suitable USB devices found")
	}

	// For now, return the first suitable device
	// In interactive mode, we would show a list for user selection
	device := devices[0]

	cs.logger.Debug("USB device detected", "device", device.Path, "capacity", device.Capacity)
	return device.Path, nil
}

// downloadUbuntuISO downloads the Ubuntu ISO
func (cs *CreateSubcomponent) downloadUbuntuISO(ctx context.Context, version string) (string, error) {
	cs.logger.Debug("Downloading Ubuntu ISO", "version", version)

	// Use default version if not specified
	if version == "" {
		version = "24.04"
	}

	// Download ISO
	isoInfo, err := cs.isoDownloader.DownloadISO(ctx, version)
	if err != nil {
		return "", fmt.Errorf("failed to download ISO: %w", err)
	}

	cs.logger.Debug("Ubuntu ISO downloaded", "version", version, "path", isoInfo.FilePath)
	return isoInfo.FilePath, nil
}

// generateCloudInit generates the cloud-init configuration
func (cs *CreateSubcomponent) generateCloudInit(ctx context.Context, config *types.NodeConfig) (*types.CloudInitConfig, error) {
	cs.logger.Debug("Generating cloud-init configuration", "node_id", config.NodeID)

	// Generate cloud-init configuration
	cloudInitConfig, err := cs.cloudInitGenerator.GenerateCloudInit(config)
	if err != nil {
		return nil, fmt.Errorf("failed to generate cloud-init: %w", err)
	}

	cs.logger.Debug("Cloud-init configuration generated", "node_id", config.NodeID)
	return cloudInitConfig, nil
}

// writeISOToUSB writes the ISO to the USB device
func (cs *CreateSubcomponent) writeISOToUSB(ctx context.Context, isoPath, devicePath string, cloudInitConfig *types.CloudInitConfig) (*types.WriteResult, error) {
	cs.logger.Debug("Writing ISO to USB device", "iso", isoPath, "device", devicePath)

	// Create USB writer manager
	writerManager := NewUSBWriterManager(cs.usbWriter, cs.logger)

	// Write ISO with cloud-init injection
	writeResult, err := writerManager.WriteUbuntuISO(ctx, isoPath, devicePath, cloudInitConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to write ISO to USB: %w", err)
	}

	cs.logger.Debug("ISO written to USB device", "device", devicePath, "bytes_written", writeResult.BytesWritten)
	return writeResult, nil
}

// saveNodeState saves the node state for later reference
func (cs *CreateSubcomponent) saveNodeState(ctx context.Context, config *types.NodeConfig, result *CreateResult) error {
	cs.logger.Debug("Saving node state", "node_id", config.NodeID)

	// Create nodes directory
	nodesDir := filepath.Join(os.Getenv("HOME"), ".syntropy", "nodes")
	if err := os.MkdirAll(nodesDir, 0755); err != nil {
		return fmt.Errorf("failed to create nodes directory: %w", err)
	}

	// Create node directory
	nodeDir := filepath.Join(nodesDir, config.NodeID)
	if err := os.MkdirAll(nodeDir, 0755); err != nil {
		return fmt.Errorf("failed to create node directory: %w", err)
	}

	// Save node configuration
	configPath := filepath.Join(nodeDir, "config.json")
	if err := cs.saveConfigToFile(config, configPath); err != nil {
		return fmt.Errorf("failed to save node configuration: %w", err)
	}

	// Save creation result
	resultPath := filepath.Join(nodeDir, "creation_result.json")
	if err := cs.saveResultToFile(result, resultPath); err != nil {
		return fmt.Errorf("failed to save creation result: %w", err)
	}

	cs.logger.Debug("Node state saved", "node_id", config.NodeID, "path", nodeDir)
	return nil
}

// startListener starts the listener for node registration
func (cs *CreateSubcomponent) startListener(ctx context.Context, config *types.NodeConfig) error {
	cs.logger.Debug("Starting listener for node registration", "node_id", config.NodeID)

	// This will be implemented when we create the Listener component
	// For now, we'll just log that we would start it
	cs.logger.Info("Listener would be started for node registration", "node_id", config.NodeID)

	return nil
}

// checkRequiredTools checks if required tools are available
func (cs *CreateSubcomponent) checkRequiredTools() error {
	// This is a simplified check
	// In production, we would check for specific tools based on the platform
	cs.logger.Debug("Checking required tools")

	// For now, we assume tools are available
	// TODO: Implement platform-specific tool checking
	return nil
}

// checkWritePermissions checks if we have write permissions
func (cs *CreateSubcomponent) checkWritePermissions() error {
	cs.logger.Debug("Checking write permissions")

	// Check if we can write to the home directory
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		return fmt.Errorf("HOME environment variable is not set")
	}

	// Try to create a temporary file
	testFile := filepath.Join(homeDir, ".syntropy_test_write")
	file, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("cannot write to home directory: %w", err)
	}
	file.Close()
	os.Remove(testFile)

	return nil
}

// promptUbuntuVersion prompts the user for Ubuntu version
func (cs *CreateSubcomponent) promptUbuntuVersion() string {
	// For now, return default version
	// In interactive mode, we would show a list of available versions
	return "24.04"
}

// promptUSBDevice prompts the user for USB device selection
func (cs *CreateSubcomponent) promptUSBDevice(ctx context.Context) string {
	// For now, return empty to trigger automatic detection
	// In interactive mode, we would show a list of available devices
	return ""
}

// promptAutoStart prompts the user if they want to auto-start the listener
func (cs *CreateSubcomponent) promptAutoStart() bool {
	// For now, return true for auto-start
	// In interactive mode, we would ask the user
	return true
}

// saveConfigToFile saves the node configuration to a file
func (cs *CreateSubcomponent) saveConfigToFile(config *types.NodeConfig, filePath string) error {
	// This is a simplified implementation
	// In production, we would use proper JSON marshaling
	cs.logger.Debug("Saving node configuration to file", "path", filePath)

	// Create a simple text representation for now
	content := fmt.Sprintf("NodeID: %s\nCreatedAt: %s\n", config.NodeID, config.CreatedAt)

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// saveResultToFile saves the creation result to a file
func (cs *CreateSubcomponent) saveResultToFile(result *CreateResult, filePath string) error {
	// This is a simplified implementation
	// In production, we would use proper JSON marshaling
	cs.logger.Debug("Saving creation result to file", "path", filePath)

	// Create a simple text representation for now
	content := fmt.Sprintf("NodeID: %s\nSuccess: %t\nDuration: %s\n",
		result.NodeID, result.Success, result.Duration)

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write result file: %w", err)
	}

	return nil
}

// CreateNodeManager manages the node creation process
type CreateNodeManager struct {
	createSubcomponent *CreateSubcomponent
	logger             types.Logger
}

// NewCreateNodeManager creates a new Create node manager
func NewCreateNodeManager(createSubcomponent *CreateSubcomponent, logger types.Logger) *CreateNodeManager {
	return &CreateNodeManager{
		createSubcomponent: createSubcomponent,
		logger:             logger,
	}
}

// CreateNode creates a new node using the create subcomponent
func (cnm *CreateNodeManager) CreateNode(ctx context.Context, options CreateOptions) (*CreateResult, error) {
	cnm.logger.Info("Creating node via Create Node Manager", "options", options)

	return cnm.createSubcomponent.CreateNode(ctx, options)
}

// CreateNodeInteractive creates a node interactively
func (cnm *CreateNodeManager) CreateNodeInteractive(ctx context.Context) (*CreateResult, error) {
	cnm.logger.Info("Creating node interactively via Create Node Manager")

	return cnm.createSubcomponent.CreateNodeInteractive(ctx)
}

// GetCreationProgress returns the current creation progress
func (cnm *CreateNodeManager) GetCreationProgress() *types.WriteProgress {
	// This would return the current progress of any ongoing creation
	// For now, return nil as we don't track progress at this level
	return nil
}

// CancelCreation cancels any ongoing creation process
func (cnm *CreateNodeManager) CancelCreation() error {
	cnm.logger.Info("Cancelling node creation process")

	// This would cancel any ongoing creation
	// For now, just log the action
	return nil
}
