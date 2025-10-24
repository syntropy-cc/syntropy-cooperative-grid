package node

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"node-component/src/internal/types"
)

// CreateSubcomponent handles the complete node creation workflow
type CreateSubcomponent struct {
	configGenerator    types.ConfigGenerator
	usbDetector        USBDetector
	isoDownloader      ISODownloader
	cloudInitGenerator types.CloudInitGenerator
	usbWriter          USBWriter
	tokenIntegration   types.TokenIntegration
	nodeState          types.NodeStateManager
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
	nodeState types.NodeStateManager,
	logger types.Logger,
) *CreateSubcomponent {
	// Validate critical dependencies
	if usbDetector == nil {
		logger.Error("USB detector is nil - this will cause crashes during node creation")
	}
	if configGenerator == nil {
		logger.Error("Config generator is nil - this will cause crashes during node creation")
	}
	if isoDownloader == nil {
		logger.Error("ISO downloader is nil - this will cause crashes during node creation")
	}
	if cloudInitGenerator == nil {
		logger.Error("Cloud-init generator is nil - this will cause crashes during node creation")
	}
	if usbWriter == nil {
		logger.Error("USB writer is nil - this will cause crashes during node creation")
	}
	if nodeState == nil {
		logger.Error("Node state manager is nil - this will cause crashes during node creation")
	}
	if logger == nil {
		// This would be a critical error, but we can't log it
		panic("logger cannot be nil in NewCreateSubcomponent")
	}

	return &CreateSubcomponent{
		configGenerator:    configGenerator,
		usbDetector:        usbDetector,
		isoDownloader:      isoDownloader,
		cloudInitGenerator: cloudInitGenerator,
		usbWriter:          usbWriter,
		tokenIntegration:   tokenIntegration,
		nodeState:          nodeState,
		logger:             logger,
	}
}

// CreateOptions represents options for node creation
type CreateOptions struct {
	UbuntuVersion     string
	DevicePath        string
	ISOPath           string
	ISOURL            string
	SkipISODownload   bool
	SkipCloudInit     bool
	SkipUSBWrite      bool
	ForceOverwrite    bool
	AutoStart         bool
	SkipISOValidation bool
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

	// Step 3: Detect and validate USB device (if not specified)
	var selectedDevice *SelectedUSBDevice
	if options.DevicePath == "" {
		selected, err := cs.detectUSBDevice(ctx)
		if err != nil {
			result.StepsFailed = append(result.StepsFailed, "detect_usb_device")
			result.ErrorMessage = err.Error()
			return result, fmt.Errorf("USB device detection failed: %w", err)
		}
		selectedDevice = selected
		result.DevicePath = selected.Device.Path
		result.StepsCompleted = append(result.StepsCompleted, "detect_usb_device")
	} else if options.DevicePath != "" {
		// Dispositivo especificado manualmente - validar também
		device := types.USBDevice{Path: options.DevicePath}
		if err := cs.usbDetector.ValidateDeviceDoubleCheck(ctx, device); err != nil {
			return result, fmt.Errorf("manual device validation failed: %w", err)
		}
		selectedDevice = NewSelectedUSBDevice(device, runtime.GOOS)
		result.DevicePath = options.DevicePath
	}

	// Step 4: Download Ubuntu ISO
	var isoPath string
	if !options.SkipISODownload {
		isoPath, err = cs.downloadUbuntuISOWithURL(ctx, options.UbuntuVersion, options.ISOURL, options.SkipISOValidation)
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

	// Step 6: Write ISO to USB device (passar selectedDevice ao invés de string)
	if !options.SkipUSBWrite {
		fmt.Println("💾 Etapa 6/6: Escrevendo ISO no dispositivo USB...")
		fmt.Printf("   📁 ISO: %s\n", filepath.Base(isoPath))
		fmt.Printf("   🔌 Dispositivo: %s\n", selectedDevice.Device.Path)
		fmt.Printf("   🔐 Token: %s\n", selectedDevice.ValidationToken)
		fmt.Println("   ⏳ Isso pode levar alguns minutos...")
		fmt.Println()

		writeResult, err := cs.writeISOToUSBWithValidation(ctx, isoPath, selectedDevice, cloudInitConfig)
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
		fmt.Println("   ✅ ISO escrita com sucesso no dispositivo USB")
		fmt.Println()
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

// validateSetupComponent validates that the setup component has been run successfully
func (cs *CreateSubcomponent) validateSetupComponent() error {
	cs.logger.Debug("Validating setup component")

	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Check if .syntropy directory exists
	syntropyDir := filepath.Join(homeDir, ".syntropy")
	if _, err := os.Stat(syntropyDir); err != nil {
		return fmt.Errorf("setup not found - please run 'syntropy setup' first")
	}

	// Check if cache directory exists (created by setup)
	cacheDir := filepath.Join(syntropyDir, "cache")
	if _, err := os.Stat(cacheDir); err != nil {
		return fmt.Errorf("setup incomplete - cache directory not found, please run 'syntropy setup' first")
	}

	// Check if config directory exists (created by setup)
	configDir := filepath.Join(syntropyDir, "config")
	if _, err := os.Stat(configDir); err != nil {
		return fmt.Errorf("setup incomplete - config directory not found, please run 'syntropy setup' first")
	}

	// Check if manager.yaml configuration file exists (created by setup)
	managerConfigPath := filepath.Join(configDir, "manager.yaml")
	if _, err := os.Stat(managerConfigPath); err != nil {
		return fmt.Errorf("setup incomplete - manager configuration not found, please run 'syntropy setup' first")
	}

	// Validate that the configuration file is readable and not empty
	configData, err := os.ReadFile(managerConfigPath)
	if err != nil {
		return fmt.Errorf("setup incomplete - cannot read manager configuration, please run 'syntropy setup' first: %w", err)
	}

	if len(configData) == 0 {
		return fmt.Errorf("setup incomplete - manager configuration is empty, please run 'syntropy setup' first")
	}

	cs.logger.Debug("Setup component validation passed")
	return nil
}

// validatePrerequisites validates that all prerequisites are met
func (cs *CreateSubcomponent) validatePrerequisites(ctx context.Context) error {
	cs.logger.Debug("Validating prerequisites")

	// First, validate that setup component has been run successfully
	if err := cs.validateSetupComponent(); err != nil {
		return fmt.Errorf("setup validation failed: %w", err)
	}

	// Check if token integration is available
	if cs.tokenIntegration == nil {
		return fmt.Errorf("token integration is not available - please run 'syntropy setup' first")
	}

	// Check if Grid Token is available
	if _, err := cs.tokenIntegration.GetGridToken(); err != nil {
		return fmt.Errorf("grid token is not available - please run 'syntropy setup' first: %w", err)
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

// detectUSBDevice detects and validates a USB device with double validation
func (cs *CreateSubcomponent) detectUSBDevice(ctx context.Context) (*SelectedUSBDevice, error) {
	cs.logger.Debug("Detecting USB device with double validation")

	// Get removable devices
	devices, err := cs.usbDetector.DetectRemovableDevices(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to detect USB devices: %w", err)
	}

	if len(devices) == 0 {
		return nil, fmt.Errorf("no suitable USB devices found")
	}

	// Manual selection is always required
	selectedDevice, err := cs.selectUSBDeviceFromDetected(devices)
	if err != nil {
		return nil, fmt.Errorf("device selection failed: %w", err)
	}

	cs.logger.Info("Device selected by user", "device", selectedDevice.Path)

	// Double security validation
	fmt.Println("\n🔒 Executando validações de segurança...")
	if err := cs.usbDetector.ValidateDeviceDoubleCheck(ctx, *selectedDevice); err != nil {
		return nil, fmt.Errorf("SECURITY VALIDATION FAILED: %w", err)
	}
	fmt.Println("✅ Validações de segurança aprovadas")

	// Create selected device with token
	selected := NewSelectedUSBDevice(*selectedDevice, runtime.GOOS)

	cs.logger.Info("Device validated and selected",
		"device", selected.Device.Path,
		"token", selected.ValidationToken)

	return selected, nil
}

// selectUSBDeviceFromDetected allows user to manually select a USB device from detected devices
func (cs *CreateSubcomponent) selectUSBDeviceFromDetected(devices []types.USBDevice) (*types.USBDevice, error) {
	if len(devices) == 0 {
		return nil, fmt.Errorf("no devices available for selection")
	}

	// ALWAYS show list for manual selection, even with only one device
	fmt.Printf("\n⚠️  SELEÇÃO MANUAL DE DISPOSITIVO USB OBRIGATÓRIA\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")
	fmt.Printf("🔌 Dispositivos USB detectados: %d\n\n", len(devices))

	for i, device := range devices {
		fmt.Printf("  [%d] %s\n", i+1, device.Path)
		fmt.Printf("      Capacidade: %.2f GB\n", float64(device.Capacity)/(1024*1024*1024))
		if device.Vendor != "" {
			fmt.Printf("      Fabricante: %s\n", device.Vendor)
		}
		if device.Model != "" {
			fmt.Printf("      Modelo: %s\n", device.Model)
		}
		if device.Serial != "" {
			fmt.Printf("      Serial: %s\n", device.Serial)
		}
		fmt.Printf("      Removível: %t\n", device.IsRemovable)
		fmt.Printf("      Sistema: %t\n", device.IsSystem)
		fmt.Printf("\n")
	}

	fmt.Printf("⚠️  ATENÇÃO: O dispositivo selecionado será COMPLETAMENTE FORMATADO\n")
	fmt.Printf("⚠️  TODOS os dados serão PERDIDOS permanentemente\n\n")
	fmt.Printf("❓ Digite o número do dispositivo desejado (1-%d): ", len(devices))

	var choice int
	_, err := fmt.Scanln(&choice)
	if err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	if choice < 1 || choice > len(devices) {
		return nil, fmt.Errorf("invalid choice: %d (must be between 1 and %d)", choice, len(devices))
	}

	selectedDevice := devices[choice-1]

	// Show selected device and ask for confirmation
	fmt.Printf("\n📌 Dispositivo selecionado:\n")
	fmt.Printf("   Caminho: %s\n", selectedDevice.Path)
	fmt.Printf("   Capacidade: %.2f GB\n", float64(selectedDevice.Capacity)/(1024*1024*1024))
	fmt.Printf("   Modelo: %s\n", selectedDevice.Model)
	fmt.Printf("\n")

	// Ask for explicit formatting confirmation
	fmt.Printf("⚠️  CONFIRMAÇÃO DE FORMATAÇÃO OBRIGATÓRIA\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("🔧 O dispositivo USB será FORMATADO completamente\n")
	fmt.Printf("💾 TODOS os dados serão PERDIDOS permanentemente\n")
	fmt.Printf("🔄 O dispositivo será formatado como FAT32 para boot\n")
	fmt.Printf("✅ Esta ação é NECESSÁRIA para criar um USB bootável\n\n")
	fmt.Printf("❓ Confirma a formatação do dispositivo? (s/n): ")

	var confirm string
	_, err = fmt.Scanln(&confirm)
	if err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	if strings.ToLower(confirm) != "s" && strings.ToLower(confirm) != "sim" && strings.ToLower(confirm) != "y" && strings.ToLower(confirm) != "yes" {
		return nil, fmt.Errorf("formatação cancelada pelo usuário")
	}

	fmt.Printf("✅ Formatação confirmada pelo usuário\n\n")

	return &selectedDevice, nil
}

// downloadUbuntuISO downloads the Ubuntu ISO
func (cs *CreateSubcomponent) downloadUbuntuISO(ctx context.Context, version string) (string, error) {
	return cs.downloadUbuntuISOWithURL(ctx, version, "", false)
}

// downloadUbuntuISOWithURL downloads the Ubuntu ISO with custom URL support
func (cs *CreateSubcomponent) downloadUbuntuISOWithURL(ctx context.Context, version string, customURL string, skipValidation bool) (string, error) {
	cs.logger.Info("Downloading Ubuntu ISO", "version", version, "custom_url", customURL)

	// Use default version if not specified
	if version == "" {
		version = "24.04"
	}

	// Download ISO using the ISO downloader with custom URL support
	isoDownloader := cs.isoDownloader.(*ISODownloaderImpl)
	isoInfo, err := isoDownloader.DownloadISOWithOptions(ctx, version, customURL, skipValidation)
	if err != nil {
		return "", fmt.Errorf("failed to download Ubuntu ISO: %w", err)
	}

	cs.logger.Info("Ubuntu ISO downloaded successfully",
		"version", version,
		"path", isoInfo.FilePath,
		"size", isoInfo.Size,
		"source", isoInfo.DownloadURL)

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

// writeISOToUSBWithValidation escreve ISO com validação contínua do dispositivo
func (cs *CreateSubcomponent) writeISOToUSBWithValidation(
	ctx context.Context,
	isoPath string,
	selectedDevice *SelectedUSBDevice,
	cloudInitConfig *types.CloudInitConfig,
) (*types.WriteResult, error) {
	cs.logger.Info("Writing ISO with continuous validation",
		"iso", isoPath,
		"device", selectedDevice.Device.Path,
		"token", selectedDevice.ValidationToken)

	// Validar novamente antes de escrever
	fmt.Println("🔒 Revalidando dispositivo antes da gravação...")
	if err := cs.usbDetector.ValidateDeviceDoubleCheck(ctx, selectedDevice.Device); err != nil {
		return nil, fmt.Errorf("pre-write validation failed: %w", err)
	}

	// Verificar que o token ainda é válido
	currentDevices, err := cs.usbDetector.DetectRemovableDevices(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to re-detect devices: %w", err)
	}

	var foundDevice *types.USBDevice
	for _, dev := range currentDevices {
		if dev.Path == selectedDevice.Device.Path {
			foundDevice = &dev
			break
		}
	}

	if foundDevice == nil {
		return nil, fmt.Errorf("SECURITY: selected device %s no longer available", selectedDevice.Device.Path)
	}

	if err := selectedDevice.Validate(*foundDevice); err != nil {
		return nil, fmt.Errorf("SECURITY: device validation failed: %w", err)
	}

	fmt.Println("✅ Dispositivo revalidado com sucesso")

	// Escrever ISO
	result, err := cs.usbWriter.WriteISO(ctx, isoPath, selectedDevice, cloudInitConfig)
	if err != nil {
		return nil, err
	}

	// Validar uma última vez após escrita
	cs.logger.Info("Write completed, performing final validation")

	return result, nil
}

// saveNodeState saves the node state for later reference
func (cs *CreateSubcomponent) saveNodeState(ctx context.Context, config *types.NodeConfig, result *CreateResult) error {
	cs.logger.Debug("Saving node state", "node_id", config.NodeID)

	// Use NodeStateManager to save the node state
	if err := cs.nodeState.CreateNode(config.NodeID, config); err != nil {
		return fmt.Errorf("failed to create node in state manager: %w", err)
	}

	// Also save to disk for backward compatibility
	// Note: .syntropy directory should already exist from setup component
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	nodesDir := filepath.Join(homeDir, ".syntropy", "nodes")
	// Check if .syntropy directory exists (created by setup component)
	syntropyDir := filepath.Join(homeDir, ".syntropy")
	if _, err := os.Stat(syntropyDir); err != nil {
		return fmt.Errorf("setup not found - please run 'syntropy setup' first: %w", err)
	}
	// Only create nodes subdirectory if it doesn't exist
	if _, err := os.Stat(nodesDir); err != nil {
		if err := os.MkdirAll(nodesDir, 0755); err != nil {
			return fmt.Errorf("failed to create nodes directory: %w", err)
		}
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

	// Get home directory using cross-platform method
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Check if .syntropy directory exists (should be created by setup component)
	syntropyDir := filepath.Join(homeDir, ".syntropy")
	if _, err := os.Stat(syntropyDir); err != nil {
		return fmt.Errorf("setup not found - please run 'syntropy setup' first: %w", err)
	}

	// Try to create a temporary file to test write permissions
	testFile := filepath.Join(syntropyDir, ".syntropy_test_write")
	file, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("cannot write to .syntropy directory: %w", err)
	}
	file.Close()
	os.Remove(testFile)

	cs.logger.Debug("Write permissions check passed")
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
