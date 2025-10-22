package node

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"node-component/src/internal/constants"
	"node-component/src/internal/helpers"
	"node-component/src/internal/types"
)

// Type aliases for compatibility - using types from types package

// NodeManager is the main orchestrator for node operations
type NodeManager struct {
	// Subcomponents
	createSubcomponent       *CreateSubcomponent
	registrationSubcomponent *RegistrationSubcomponent

	// Core services
	eventBus         types.EventBus
	nodeState        types.NodeStateManager
	tokenIntegration types.TokenIntegration
	logger           types.Logger

	// Component services
	configGenerator    types.ConfigGenerator
	usbDetector        USBDetector
	isoDownloader      ISODownloader
	cloudInitGenerator types.CloudInitGenerator
	usbWriter          USBWriter

	// Configuration
	config *types.Configuration

	// State management
	mutex     sync.RWMutex
	isRunning bool
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewNodeManager creates a new NodeManager instance
func NewNodeManager() *NodeManager {
	ctx, cancel := context.WithCancel(context.Background())

	return &NodeManager{
		eventBus:  NewEventBus(),
		nodeState: NewNodeStateManager(NewLogger()),
		logger:    NewLogger(),
		config:    GetDefaultConfiguration(),
		ctx:       ctx,
		cancel:    cancel,
		isRunning: false,
	}
}

// Initialize initializes the NodeManager with required dependencies
func (nm *NodeManager) Initialize() error {
	nm.mutex.Lock()
	defer nm.mutex.Unlock()

	nm.logger.Info("Initializing NodeManager...")

	// Initialize token integration with Setup Component
	setupAdapter := NewSetupAdapter(nm.logger)
	nm.tokenIntegration = setupAdapter.CreateTokenIntegration()

	// Initialize the token integration with Setup Component
	if err := setupAdapter.InitializeTokenIntegration(nm.tokenIntegration); err != nil {
		// Check if this is a "token not found" error vs other errors
		if strings.Contains(err.Error(), "grid token not found in setup component") {
			// Token doesn't exist - this is expected and should fail commands that need token
			nm.logger.Warn("Grid token not found - some commands will not be available", "error", err)
			nm.tokenIntegration = nil
		} else {
			// Other errors (permission, corruption, etc.) - fail initialization
			nm.logger.Error("Failed to initialize token integration due to system error", "error", err)
			return fmt.Errorf("token integration initialization failed: %w", err)
		}
	} else {
		// Token integration successful - no warnings needed
		nm.logger.Debug("Token integration initialized successfully")
	}

	// Initialize component services
	nm.configGenerator = NewAutoConfigGenerator(nm.tokenIntegration, nm.logger)
	factory := NewUSBDetectorFactory()
	nm.usbDetector, _ = factory.CreateUSBDetector(nm.logger)
	nm.isoDownloader = NewISODownloader(nm.logger)
	nm.cloudInitGenerator = NewCloudInitGenerator(nm.logger)
	writerFactory := NewUSBWriterFactory()
	nm.usbWriter, _ = writerFactory.CreateUSBWriter(nm.logger)

	// Initialize create subcomponent
	nm.createSubcomponent = NewCreateSubcomponent(
		nm.configGenerator,
		nm.usbDetector,
		nm.isoDownloader,
		nm.cloudInitGenerator,
		nm.usbWriter,
		nm.tokenIntegration,
		nm.nodeState,
		nm.logger,
	)

	// Initialize registration subcomponent
	nm.registrationSubcomponent = NewRegistrationSubcomponent(
		nm.nodeState,
		nm.eventBus,
		nm.logger,
		nm.config,
	)

	// Load existing node state
	if err := nm.nodeState.LoadState(); err != nil {
		nm.logger.Warn("Failed to load existing node state", "error", err)
	}

	// Subscribe to events
	nm.eventBus.Subscribe(constants.EventNodeCreated, nm.handleNodeCreated)
	nm.eventBus.Subscribe(constants.EventNodeRegistered, nm.handleNodeRegistered)
	nm.eventBus.Subscribe(constants.EventNodeDisconnected, nm.handleNodeDisconnected)

	nm.isRunning = true
	nm.logger.Info("NodeManager initialized successfully")

	return nil
}

// CreateNode creates a new node with the specified options
func (nm *NodeManager) CreateNode(options *CreateOptions) (*CreateResult, error) {
	nm.mutex.RLock()
	if !nm.isRunning {
		nm.mutex.RUnlock()
		return nil, fmt.Errorf("NodeManager is not running")
	}
	nm.mutex.RUnlock()

	nm.logger.Info("Creating new node", "options", options)

	// Validate options
	if err := nm.validateCreateOptions(options); err != nil {
		return nil, fmt.Errorf("invalid options: %w", err)
	}

	// Create node using create subcomponent
	result, err := nm.createSubcomponent.CreateNode(context.Background(), *options)
	if err != nil {
		nm.logger.Error("Node creation failed", "error", err)
		return nil, fmt.Errorf("node creation failed: %w", err)
	}

	// Node state is already saved by the CreateSubcomponent

	// Publish event
	nm.eventBus.Publish(types.Event{
		Type:      constants.EventNodeCreated,
		NodeID:    result.NodeID,
		Data:      map[string]interface{}{"result": result},
		Timestamp: time.Now(),
	})

	nm.logger.Info("Node created successfully", "nodeID", result.NodeID, "devicePath", result.DevicePath)

	return result, nil
}

// ListNodes lists all nodes with their current status
func (nm *NodeManager) ListNodes() (*NodeList, error) {
	nm.mutex.RLock()
	defer nm.mutex.RUnlock()

	if !nm.isRunning {
		return nil, fmt.Errorf("NodeManager is not running")
	}

	activeNodes := nm.nodeState.GetActiveNodes()
	pendingNodes := nm.nodeState.GetPendingNodes()
	inactiveNodes := nm.nodeState.GetInactiveNodes()

	total := len(activeNodes) + len(pendingNodes) + len(inactiveNodes)

	return &NodeList{
		Active:   convertToNodeStatusSlice(activeNodes),
		Pending:  convertToNodeStatusSlice(pendingNodes),
		Inactive: convertToNodeStatusSlice(inactiveNodes),
		Total:    total,
	}, nil
}

// GetNodeStatus gets detailed status of a specific node
func (nm *NodeManager) GetNodeStatus(nodeID string) (*NodeStatus, error) {
	nm.mutex.RLock()
	defer nm.mutex.RUnlock()

	if !nm.isRunning {
		return nil, fmt.Errorf("NodeManager is not running")
	}

	// Validate node ID
	if err := helpers.ValidateNodeID(nodeID); err != nil {
		return nil, fmt.Errorf("invalid node ID: %w", err)
	}

	status, err := nm.nodeState.GetNode(nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get node status: %w", err)
	}

	return convertToNodeStatus(status), nil
}

// GetNodeLogs retrieves logs for a specific node
func (nm *NodeManager) GetNodeLogs(nodeID string, options *LogOptions) (*NodeLogs, error) {
	nm.mutex.RLock()
	defer nm.mutex.RUnlock()

	if !nm.isRunning {
		return nil, fmt.Errorf("NodeManager is not running")
	}

	// Validate node ID
	if err := helpers.ValidateNodeID(nodeID); err != nil {
		return nil, fmt.Errorf("invalid node ID: %w", err)
	}

	// Check if node exists
	if _, err := nm.nodeState.GetNode(nodeID); err != nil {
		return nil, fmt.Errorf("node not found: %w", err)
	}

	// Read logs from file
	logPath, err := helpers.ExpandPath(fmt.Sprintf("%s/%s/%s", constants.DefaultNodeConfigDir, nodeID, constants.NodeLogFile))
	if err != nil {
		return nil, fmt.Errorf("failed to expand log path: %w", err)
	}
	logData, err := helpers.ReadFileSafely(logPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read node logs: %w", err)
	}

	// Parse logs
	logLines := strings.Split(string(logData), "\n")

	// Apply filters
	if options.Lines > 0 && len(logLines) > options.Lines {
		logLines = logLines[len(logLines)-options.Lines:]
	}

	return &NodeLogs{
		NodeID:    nodeID,
		Logs:      logLines,
		Timestamp: time.Now(),
		Service:   options.Service,
	}, nil
}

// DeleteNode deletes a node and cleans up its resources
func (nm *NodeManager) DeleteNode(nodeID string) error {
	nm.mutex.Lock()
	defer nm.mutex.Unlock()

	if !nm.isRunning {
		return fmt.Errorf("NodeManager is not running")
	}

	nm.logger.Info("Deleting node", "nodeID", nodeID)

	// Validate node ID
	if err := helpers.ValidateNodeID(nodeID); err != nil {
		return fmt.Errorf("invalid node ID: %w", err)
	}

	// Stop heartbeat if active
	if nm.registrationSubcomponent.IsHeartbeatActive(nodeID) {
		if err := nm.registrationSubcomponent.StopHeartbeat(nodeID); err != nil {
			nm.logger.Warn("Failed to stop heartbeat", "nodeID", nodeID, "error", err)
		}
	}

	// Remove from state
	if err := nm.nodeState.RemoveNode(nodeID); err != nil {
		return fmt.Errorf("failed to remove node from state: %w", err)
	}

	// Clean up files
	nodeDir, err := helpers.ExpandPath(fmt.Sprintf("%s/%s", constants.DefaultNodeConfigDir, nodeID))
	if err != nil {
		nm.logger.Warn("Failed to expand node directory path", "nodeID", nodeID, "error", err)
	} else {
		if err := os.RemoveAll(nodeDir); err != nil {
			nm.logger.Warn("Failed to clean up node directory", "nodeID", nodeID, "error", err)
		}
	}

	nm.logger.Info("Node deleted successfully", "nodeID", nodeID)

	return nil
}

// StartRegistrationListener starts the registration listener
func (nm *NodeManager) StartRegistrationListener() error {
	nm.mutex.Lock()
	defer nm.mutex.Unlock()

	if !nm.isRunning {
		return fmt.Errorf("NodeManager is not running")
	}

	nm.logger.Info("Starting registration listener")

	if err := nm.registrationSubcomponent.StartListener(constants.DefaultRegistrationPort); err != nil {
		return fmt.Errorf("failed to start registration listener: %w", err)
	}

	nm.logger.Info("Registration listener started successfully")

	return nil
}

// StopRegistrationListener stops the registration listener
func (nm *NodeManager) StopRegistrationListener() error {
	nm.mutex.Lock()
	defer nm.mutex.Unlock()

	if !nm.isRunning {
		return fmt.Errorf("NodeManager is not running")
	}

	nm.logger.Info("Stopping registration listener")

	if err := nm.registrationSubcomponent.StopListener(); err != nil {
		return fmt.Errorf("failed to stop registration listener: %w", err)
	}

	nm.logger.Info("Registration listener stopped successfully")

	return nil
}

// Shutdown gracefully shuts down the NodeManager
func (nm *NodeManager) Shutdown() error {
	nm.mutex.Lock()
	defer nm.mutex.Unlock()

	if !nm.isRunning {
		return nil
	}

	nm.logger.Info("Shutting down NodeManager...")

	// Stop registration listener
	if err := nm.registrationSubcomponent.StopListener(); err != nil {
		nm.logger.Warn("Failed to stop registration listener", "error", err)
	}

	// Stop all heartbeats
	activeNodes := nm.nodeState.GetActiveNodes()
	for _, node := range activeNodes {
		if err := nm.registrationSubcomponent.StopHeartbeat(node.NodeID); err != nil {
			nm.logger.Warn("Failed to stop heartbeat", "nodeID", node.NodeID, "error", err)
		}
	}

	// Save state
	if err := nm.nodeState.SaveState(); err != nil {
		nm.logger.Warn("Failed to save node state", "error", err)
	}

	// Close event bus
	nm.eventBus.Close()

	// Cancel context
	nm.cancel()

	nm.isRunning = false
	nm.logger.Info("NodeManager shutdown completed")

	return nil
}

// GetConfiguration returns the current configuration
func (nm *NodeManager) GetConfiguration() *types.Configuration {
	nm.mutex.RLock()
	defer nm.mutex.RUnlock()

	return nm.config
}

// UpdateConfiguration updates the configuration
func (nm *NodeManager) UpdateConfiguration(config *types.Configuration) error {
	nm.mutex.Lock()
	defer nm.mutex.Unlock()

	if !nm.isRunning {
		return fmt.Errorf("NodeManager is not running")
	}

	// TODO: Add configuration validation
	// Validate configuration
	// if err := config.Validate(); err != nil {
	//	return fmt.Errorf("invalid configuration: %w", err)
	// }

	nm.config = config
	nm.logger.Info("Configuration updated successfully")

	return nil
}

// IsRunning returns whether the NodeManager is running
func (nm *NodeManager) IsRunning() bool {
	nm.mutex.RLock()
	defer nm.mutex.RUnlock()

	return nm.isRunning
}

// Private methods

// validateCreateOptions validates the create options
func (nm *NodeManager) validateCreateOptions(options *CreateOptions) error {
	if options == nil {
		return fmt.Errorf("options cannot be nil")
	}

	// Validate device path if provided
	if options.DevicePath != "" {
		if !helpers.FileExists(options.DevicePath) {
			return fmt.Errorf("device path does not exist: %s", options.DevicePath)
		}
	}

	// Validate Ubuntu version if provided
	if options.UbuntuVersion != "" {
		if !nm.isValidUbuntuVersion(options.UbuntuVersion) {
			return fmt.Errorf("invalid Ubuntu version: %s", options.UbuntuVersion)
		}
	}

	return nil
}

// isValidUbuntuVersion validates if the Ubuntu version is supported
func (nm *NodeManager) isValidUbuntuVersion(version string) bool {
	supportedVersions := []string{"24.04", "22.04", "20.04"}
	for _, v := range supportedVersions {
		if v == version {
			return true
		}
	}
	return false
}

// Helper functions for type conversion

// convertToNodeStatusSlice converts []*types.NodeStatus to []*NodeStatus
func convertToNodeStatusSlice(typesNodes []*types.NodeStatus) []*NodeStatus {
	var nodes []*NodeStatus
	for _, tNode := range typesNodes {
		var hardware *HardwareInfo
		if tNode.Hardware != nil {
			hardware = &HardwareInfo{
				CPUCores:      tNode.Hardware.CPUCores,
				MemoryGB:      tNode.Hardware.MemoryGB,
				DiskGB:        tNode.Hardware.DiskGB,
				Hostname:      tNode.Hardware.Hostname,
				IPAddress:     tNode.Hardware.IPAddress,
				OSVersion:     tNode.Hardware.OSVersion,
				KernelVersion: tNode.Hardware.KernelVersion,
			}
		}

		nodes = append(nodes, &NodeStatus{
			NodeID:        tNode.NodeID,
			Status:        tNode.Status,
			IPAddress:     tNode.IPAddress,
			Uptime:        tNode.Uptime,
			LastHeartbeat: tNode.LastHeartbeat,
			Hardware:      hardware,
			CreatedAt:     tNode.CreatedAt,
			RegisteredAt:  tNode.RegisteredAt,
		})
	}
	return nodes
}

// convertToNodeStatus converts *types.NodeStatus to *NodeStatus
func convertToNodeStatus(typesNode *types.NodeStatus) *NodeStatus {
	var hardware *HardwareInfo
	if typesNode.Hardware != nil {
		hardware = &HardwareInfo{
			CPUCores:      typesNode.Hardware.CPUCores,
			MemoryGB:      typesNode.Hardware.MemoryGB,
			DiskGB:        typesNode.Hardware.DiskGB,
			Hostname:      typesNode.Hardware.Hostname,
			IPAddress:     typesNode.Hardware.IPAddress,
			OSVersion:     typesNode.Hardware.OSVersion,
			KernelVersion: typesNode.Hardware.KernelVersion,
		}
	}

	return &NodeStatus{
		NodeID:        typesNode.NodeID,
		Status:        typesNode.Status,
		IPAddress:     typesNode.IPAddress,
		Uptime:        typesNode.Uptime,
		LastHeartbeat: typesNode.LastHeartbeat,
		Hardware:      hardware,
		CreatedAt:     typesNode.CreatedAt,
		RegisteredAt:  typesNode.RegisteredAt,
	}
}

// convertToTypesNodeConfig converts *NodeConfig to *types.NodeConfig
func convertToTypesNodeConfig(config *NodeConfig) *types.NodeConfig {
	return &types.NodeConfig{
		NodeID:           config.NodeID,
		GridToken:        config.GridToken,
		SSHPublicKey:     config.SSHPublicKey,
		SSHPrivateKey:    config.SSHPrivateKey,
		NodeCertificate:  config.NodeCertificate,
		CommandStationIP: config.CommandStationIP,
		CreatedAt:        config.CreatedAt,
		ExpiresAt:        config.ExpiresAt,
	}
}

// Event handlers

// handleNodeCreated handles node creation events
func (nm *NodeManager) handleNodeCreated(event types.Event) {
	nm.logger.Info("Node created event received", "nodeID", event.NodeID)

	// Start registration listener if not already running
	if !nm.registrationSubcomponent.IsListenerRunning() {
		if err := nm.registrationSubcomponent.StartListener(constants.DefaultRegistrationPort); err != nil {
			nm.logger.Error("Failed to start registration listener", "error", err)
		}
	}
}

// handleNodeRegistered handles node registration events
func (nm *NodeManager) handleNodeRegistered(event types.Event) {
	nm.logger.Info("Node registered event received", "nodeID", event.NodeID)

	// Start heartbeat for the registered node
	if err := nm.registrationSubcomponent.StartHeartbeat(event.NodeID, nil); err != nil {
		nm.logger.Error("Failed to start heartbeat", "nodeID", event.NodeID, "error", err)
	}
}

// handleNodeDisconnected handles node disconnection events
func (nm *NodeManager) handleNodeDisconnected(event types.Event) {
	nm.logger.Info("Node disconnected event received", "nodeID", event.NodeID)

	// Stop heartbeat for the disconnected node
	if err := nm.registrationSubcomponent.StopHeartbeat(event.NodeID); err != nil {
		nm.logger.Warn("Failed to stop heartbeat", "nodeID", event.NodeID, "error", err)
	}

	// Update node state
	if err := nm.nodeState.TransitionToInactive(event.NodeID); err != nil {
		nm.logger.Warn("Failed to update node state", "nodeID", event.NodeID, "error", err)
	}
}

// Constructor functions for dependencies

// NewEventBus creates a new event bus
func NewEventBus() types.EventBus {
	return &eventBus{
		subscribers: make(map[string][]types.EventHandler),
		mutex:       sync.RWMutex{},
	}
}

// NewLogger creates a new logger
func NewLogger() types.Logger {
	homeDir, _ := os.UserHomeDir()
	logDir := filepath.Join(homeDir, ".syntropy", "logs")

	// Criar diretório de logs se não existir
	if err := os.MkdirAll(logDir, 0755); err != nil {
		// Se não conseguir criar diretório, usar diretório temporário
		logDir = os.TempDir()
	}

	// Criar arquivo de log com timestamp
	timestamp := time.Now().Format("20060102")
	logPath := filepath.Join(logDir, fmt.Sprintf("syntropy-%s.log", timestamp))

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		// Fallback para stdout se não conseguir criar arquivo
		logFile = os.Stdout
	}

	return &logger{
		level:   constants.DefaultLogLevel,
		logFile: logFile,
		logDir:  logDir,
	}
}

// NewNodeStateManager is now implemented in node_state.go

// NewTokenIntegration creates a new token integration
func NewTokenIntegration() types.TokenIntegration {
	return NewTokenIntegrationInstance()
}

// NewCreateSubcomponent is now implemented in create.go

// NewRegistrationSubcomponent creates a new registration subcomponent
func NewRegistrationSubcomponent(
	nodeState types.NodeStateManager,
	eventBus types.EventBus,
	logger types.Logger,
	config *types.Configuration,
) *RegistrationSubcomponent {
	return &RegistrationSubcomponent{
		nodeState: nodeState,
		eventBus:  eventBus,
		logger:    logger,
		config:    config,
	}
}

// GetDefaultConfiguration returns the default configuration
func GetDefaultConfiguration() *types.Configuration {
	return &types.Configuration{
		USBDetection: struct {
			MinCapacity    int64         `yaml:"min_capacity"`
			MaxCapacity    int64         `yaml:"max_capacity"`
			PreferredSpeed int           `yaml:"preferred_speed"`
			ScanInterval   time.Duration `yaml:"scan_interval"`
		}{
			MinCapacity:    constants.DefaultMinUSBCapacity,
			MaxCapacity:    constants.DefaultMaxUSBCapacity,
			PreferredSpeed: constants.DefaultPreferredUSBSpeed,
			ScanInterval:   constants.DefaultUSBScanInterval,
		},
		ISODownload: struct {
			URL            string        `yaml:"url"`
			CacheDir       string        `yaml:"cache_dir"`
			SHA256Checksum string        `yaml:"sha256_checksum"`
			Timeout        time.Duration `yaml:"timeout"`
		}{
			URL:            constants.DefaultISOURL,
			CacheDir:       constants.DefaultISOCacheDir,
			SHA256Checksum: constants.DefaultISOSHA256Checksum,
			Timeout:        constants.DefaultISODownloadTimeout,
		},
		Registration: struct {
			Port          int           `yaml:"port"`
			Timeout       time.Duration `yaml:"timeout"`
			MaxConcurrent int           `yaml:"max_concurrent"`
			RetryAttempts int           `yaml:"retry_attempts"`
		}{
			Port:          constants.DefaultRegistrationPort,
			Timeout:       constants.DefaultRegistrationTimeout,
			MaxConcurrent: constants.DefaultMaxConcurrentNodes,
			RetryAttempts: constants.DefaultRegistrationRetries,
		},
		Heartbeat: struct {
			Interval      time.Duration `yaml:"interval"`
			Timeout       time.Duration `yaml:"timeout"`
			MaxFailures   int           `yaml:"max_failures"`
			RetryInterval time.Duration `yaml:"retry_interval"`
		}{
			Interval:      constants.DefaultHeartbeatInterval,
			Timeout:       constants.DefaultHeartbeatTimeout,
			MaxFailures:   constants.DefaultMaxHeartbeatFailures,
			RetryInterval: constants.DefaultHeartbeatRetryInterval,
		},
		Security: struct {
			TokenExpiry   time.Duration `yaml:"token_expiry"`
			CertExpiry    time.Duration `yaml:"cert_expiry"`
			KeySize       int           `yaml:"key_size"`
			HashAlgorithm string        `yaml:"hash_algorithm"`
		}{
			TokenExpiry:   constants.DefaultTokenExpiry,
			CertExpiry:    constants.DefaultCertExpiry,
			KeySize:       constants.DefaultSSHKeySize,
			HashAlgorithm: constants.DefaultHashAlgorithm,
		},
		Logging: struct {
			Level      string `yaml:"level"`
			Format     string `yaml:"format"`
			Output     string `yaml:"output"`
			MaxSize    int    `yaml:"max_size"`
			MaxBackups int    `yaml:"max_backups"`
			MaxAge     int    `yaml:"max_age"`
		}{
			Level:      constants.DefaultLogLevel,
			Format:     constants.DefaultLogFormat,
			Output:     constants.DefaultLogOutput,
			MaxSize:    constants.DefaultLogMaxSize,
			MaxBackups: constants.DefaultLogMaxBackups,
			MaxAge:     constants.DefaultLogMaxAge,
		},
	}
}
