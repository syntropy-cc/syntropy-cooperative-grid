package node

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"node-component/src/internal/constants"
	"node-component/src/internal/helpers"
	"node-component/src/internal/types"
)

// NodeStateManager manages the state of nodes in the system
type NodeStateManager struct {
	logger             types.Logger
	mutex              sync.RWMutex
	pendingNodes       map[string]*PendingNode
	activeNodes        map[string]*ActiveNode
	inactiveNodes      map[string]*InactiveNode
	stateDir           string
	persistenceEnabled bool
}

// PendingNode represents a node waiting for connection
type PendingNode struct {
	NodeID          string                 `json:"node_id"`
	CreatedAt       time.Time              `json:"created_at"`
	ExpectedAt      time.Time              `json:"expected_at"`
	Config          *types.NodeConfig      `json:"config"`
	SSHKeys         *types.SSHKeys         `json:"ssh_keys"`
	CloudInitConfig *types.CloudInitConfig `json:"cloud_init_config"`
	USBDevice       *types.USBDevice       `json:"usb_device"`
	ISOPath         string                 `json:"iso_path"`
	Status          string                 `json:"status"` // waiting, timeout, cancelled
	TimeoutDuration time.Duration          `json:"timeout_duration"`
	LastChecked     time.Time              `json:"last_checked"`
	RetryCount      int                    `json:"retry_count"`
	MaxRetries      int                    `json:"max_retries"`
}

// ActiveNode represents a connected and active node
type ActiveNode struct {
	NodeID               string                `json:"node_id"`
	RegisteredAt         time.Time             `json:"registered_at"`
	Config               *types.NodeConfig     `json:"config"`
	SSHKeys              *types.SSHKeys        `json:"ssh_keys"`
	ConnectionInfo       *types.ConnectionInfo `json:"connection_info"`
	HardwareInfo         *types.HardwareInfo   `json:"hardware_info"`
	LastSeen             time.Time             `json:"last_seen"`
	LastHeartbeat        time.Time             `json:"last_heartbeat"`
	Status               string                `json:"status"` // active, busy, maintenance
	WorkloadCount        int                   `json:"workload_count"`
	MaxWorkloads         int                   `json:"max_workloads"`
	Metrics              *types.NodeMetrics    `json:"metrics"`
	HeartbeatFailures    int                   `json:"heartbeat_failures"`
	MaxHeartbeatFailures int                   `json:"max_heartbeat_failures"`
}

// InactiveNode represents a disconnected or failed node
type InactiveNode struct {
	NodeID           string                `json:"node_id"`
	LastActiveAt     time.Time             `json:"last_active_at"`
	InactiveSince    time.Time             `json:"inactive_since"`
	Config           *types.NodeConfig     `json:"config"`
	SSHKeys          *types.SSHKeys        `json:"ssh_keys"`
	ConnectionInfo   *types.ConnectionInfo `json:"connection_info"`
	HardwareInfo     *types.HardwareInfo   `json:"hardware_info"`
	Status           string                `json:"status"` // disconnected, failed, timeout, maintenance
	FailureReason    string                `json:"failure_reason"`
	LastError        string                `json:"last_error"`
	RetryAttempts    int                   `json:"retry_attempts"`
	MaxRetryAttempts int                   `json:"max_retry_attempts"`
	NextRetryAt      time.Time             `json:"next_retry_at"`
	RecoveryActions  []string              `json:"recovery_actions"`
}

// NodeState represents the overall state of a node
type NodeState struct {
	NodeID         string            `json:"node_id"`
	CurrentState   string            `json:"current_state"` // pending, active, inactive
	LastTransition time.Time         `json:"last_transition"`
	History        []StateTransition `json:"history"`
}

// StateTransition represents a transition between states
type StateTransition struct {
	FromState   string                 `json:"from_state"`
	ToState     string                 `json:"to_state"`
	Timestamp   time.Time              `json:"timestamp"`
	Reason      string                 `json:"reason"`
	TriggeredBy string                 `json:"triggered_by"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// NewNodeStateManager creates a new node state manager
func NewNodeStateManager(logger types.Logger) *NodeStateManager {
	stateDir := filepath.Join(constants.DefaultConfigDir, "nodes")

	return &NodeStateManager{
		logger:             logger,
		pendingNodes:       make(map[string]*PendingNode),
		activeNodes:        make(map[string]*ActiveNode),
		inactiveNodes:      make(map[string]*InactiveNode),
		stateDir:           stateDir,
		persistenceEnabled: true,
	}
}

// Interface compliance with types.NodeStateManager

// CreateNode creates a new pending node entry with the provided config
func (nsm *NodeStateManager) CreateNode(nodeID string, config *types.NodeConfig) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}
	if config.NodeID == "" {
		config.NodeID = nodeID
	}
	// Create a minimal pending node (other fields can be filled later by orchestrators)
	return nsm.AddPendingNode(config, nil, nil, nil, "")
}

// GetNode returns a unified NodeStatus regardless of the internal state bucket
func (nsm *NodeStateManager) GetNode(nodeID string) (*types.NodeStatus, error) {
	nsm.mutex.RLock()
	defer nsm.mutex.RUnlock()

	if an, ok := nsm.activeNodes[nodeID]; ok {
		return &types.NodeStatus{
			NodeID:        an.NodeID,
			Status:        an.Status,
			IPAddress:     safeConnIP(an.ConnectionInfo),
			Uptime:        time.Since(an.RegisteredAt),
			LastHeartbeat: an.LastHeartbeat,
			Hardware:      an.HardwareInfo,
			CreatedAt:     an.RegisteredAt,
			RegisteredAt:  an.RegisteredAt,
		}, nil
	}
	if pn, ok := nsm.pendingNodes[nodeID]; ok {
		return &types.NodeStatus{
			NodeID:        pn.NodeID,
			Status:        "pending",
			IPAddress:     "",
			Uptime:        0,
			LastHeartbeat: time.Time{},
			Hardware:      nil,
			CreatedAt:     pn.CreatedAt,
			RegisteredAt:  time.Time{},
		}, nil
	}
	if in, ok := nsm.inactiveNodes[nodeID]; ok {
		return &types.NodeStatus{
			NodeID:        in.NodeID,
			Status:        in.Status,
			IPAddress:     safeConnIP(in.ConnectionInfo),
			Uptime:        0,
			LastHeartbeat: in.LastActiveAt,
			Hardware:      in.HardwareInfo,
			CreatedAt:     in.LastActiveAt,
			RegisteredAt:  time.Time{},
		}, nil
	}
	return nil, fmt.Errorf("node %s not found", nodeID)
}

// TransitionToActive moves a node to active state with provided IP
func (nsm *NodeStateManager) TransitionToActive(nodeID string, ipAddress string) error {
	conn := &types.ConnectionInfo{RemoteAddr: ipAddress, Port: 0, Protocol: "tcp", ConnectedAt: time.Now(), LastPing: time.Now(), Latency: 0}
	return nsm.ActivateNode(nodeID, conn, nil)
}

// TransitionToInactive moves a node to inactive state with a generic reason
func (nsm *NodeStateManager) TransitionToInactive(nodeID string) error {
	return nsm.DeactivateNode(nodeID, "manual_transition", "")
}

// GetActiveNodes returns active nodes as []*types.NodeStatus
func (nsm *NodeStateManager) GetActiveNodes() []*types.NodeStatus {
	nsm.mutex.RLock()
	defer nsm.mutex.RUnlock()
	out := make([]*types.NodeStatus, 0, len(nsm.activeNodes))
	for _, an := range nsm.activeNodes {
		out = append(out, &types.NodeStatus{
			NodeID:        an.NodeID,
			Status:        an.Status,
			IPAddress:     safeConnIP(an.ConnectionInfo),
			Uptime:        time.Since(an.RegisteredAt),
			LastHeartbeat: an.LastHeartbeat,
			Hardware:      an.HardwareInfo,
			CreatedAt:     an.RegisteredAt,
			RegisteredAt:  an.RegisteredAt,
		})
	}
	return out
}

// GetPendingNodes returns pending nodes as []*types.NodeStatus
func (nsm *NodeStateManager) GetPendingNodes() []*types.NodeStatus {
	nsm.mutex.RLock()
	defer nsm.mutex.RUnlock()
	out := make([]*types.NodeStatus, 0, len(nsm.pendingNodes))
	for _, pn := range nsm.pendingNodes {
		out = append(out, &types.NodeStatus{
			NodeID:        pn.NodeID,
			Status:        "pending",
			IPAddress:     "",
			Uptime:        0,
			LastHeartbeat: time.Time{},
			Hardware:      nil,
			CreatedAt:     pn.CreatedAt,
			RegisteredAt:  time.Time{},
		})
	}
	return out
}

// GetInactiveNodes returns inactive nodes as []*types.NodeStatus
func (nsm *NodeStateManager) GetInactiveNodes() []*types.NodeStatus {
	nsm.mutex.RLock()
	defer nsm.mutex.RUnlock()
	out := make([]*types.NodeStatus, 0, len(nsm.inactiveNodes))
	for _, in := range nsm.inactiveNodes {
		out = append(out, &types.NodeStatus{
			NodeID:        in.NodeID,
			Status:        in.Status,
			IPAddress:     safeConnIP(in.ConnectionInfo),
			Uptime:        0,
			LastHeartbeat: in.LastActiveAt,
			Hardware:      in.HardwareInfo,
			CreatedAt:     in.LastActiveAt,
			RegisteredAt:  time.Time{},
		})
	}
	return out
}

// RemoveNode removes a node from any state and deletes persistence
func (nsm *NodeStateManager) RemoveNode(nodeID string) error {
	nsm.mutex.Lock()
	defer nsm.mutex.Unlock()

	deleted := false
	if _, ok := nsm.pendingNodes[nodeID]; ok {
		delete(nsm.pendingNodes, nodeID)
		_ = nsm.removePendingNodeFromDisk(nodeID)
		deleted = true
	}
	if _, ok := nsm.activeNodes[nodeID]; ok {
		delete(nsm.activeNodes, nodeID)
		_ = nsm.removeActiveNodeFromDisk(nodeID)
		deleted = true
	}
	if _, ok := nsm.inactiveNodes[nodeID]; ok {
		delete(nsm.inactiveNodes, nodeID)
		_ = nsm.removeInactiveNodeFromDisk(nodeID)
		deleted = true
	}
	if !deleted {
		return fmt.Errorf("node %s not found", nodeID)
	}
	return nil
}

// SaveState persists all current nodes to disk
func (nsm *NodeStateManager) SaveState() error {
	nsm.mutex.RLock()
	defer nsm.mutex.RUnlock()
	for _, pn := range nsm.pendingNodes {
		if err := nsm.savePendingNodeToDisk(pn); err != nil {
			return err
		}
	}
	for _, an := range nsm.activeNodes {
		if err := nsm.saveActiveNodeToDisk(an); err != nil {
			return err
		}
	}
	for _, in := range nsm.inactiveNodes {
		if err := nsm.saveInactiveNodeToDisk(in); err != nil {
			return err
		}
	}
	return nil
}

// LoadState loads all nodes from disk
func (nsm *NodeStateManager) LoadState() error {
	nsm.mutex.Lock()
	defer nsm.mutex.Unlock()
	return nsm.loadStateFromDisk()
}

// Helpers
func safeConnIP(ci *types.ConnectionInfo) string {
	if ci == nil {
		return ""
	}
	return ci.RemoteAddr
}

// Initialize initializes the node state manager
func (nsm *NodeStateManager) Initialize() error {
	nsm.logger.Info("Initializing node state manager", "state_dir", nsm.stateDir)

	// Create state directory
	if err := os.MkdirAll(nsm.stateDir, 0755); err != nil {
		return fmt.Errorf("failed to create state directory: %w", err)
	}

	// Load existing state from disk
	if nsm.persistenceEnabled {
		if err := nsm.loadStateFromDisk(); err != nil {
			nsm.logger.Warn("Failed to load state from disk", "error", err)
		}
	}

	nsm.logger.Info("Node state manager initialized",
		"pending_nodes", len(nsm.pendingNodes),
		"active_nodes", len(nsm.activeNodes),
		"inactive_nodes", len(nsm.inactiveNodes))

	return nil
}

// AddPendingNode adds a node to the pending state
func (nsm *NodeStateManager) AddPendingNode(nodeConfig *types.NodeConfig, sshKeys *types.SSHKeys, cloudInitConfig *types.CloudInitConfig, usbDevice *types.USBDevice, isoPath string) error {
	nsm.mutex.Lock()
	defer nsm.mutex.Unlock()

	if _, exists := nsm.pendingNodes[nodeConfig.NodeID]; exists {
		return fmt.Errorf("node %s is already pending", nodeConfig.NodeID)
	}

	pendingNode := &PendingNode{
		NodeID:          nodeConfig.NodeID,
		CreatedAt:       time.Now(),
		ExpectedAt:      time.Now().Add(60 * time.Minute), // Expect connection within 30 minutes
		Config:          nodeConfig,
		SSHKeys:         sshKeys,
		CloudInitConfig: cloudInitConfig,
		USBDevice:       usbDevice,
		ISOPath:         isoPath,
		Status:          "waiting",
		TimeoutDuration: 60 * time.Minute,
		LastChecked:     time.Now(),
		RetryCount:      0,
		MaxRetries:      3,
	}

	nsm.pendingNodes[nodeConfig.NodeID] = pendingNode

	// Persist to disk
	if nsm.persistenceEnabled {
		if err := nsm.savePendingNodeToDisk(pendingNode); err != nil {
			nsm.logger.Error("Failed to save pending node to disk", "node_id", nodeConfig.NodeID, "error", err)
		}
	}

	nsm.logger.Info("Node added to pending state", "node_id", nodeConfig.NodeID)
	return nil
}

// ActivateNode moves a node from pending to active state
func (nsm *NodeStateManager) ActivateNode(nodeID string, connectionInfo *types.ConnectionInfo, hardwareInfo *types.HardwareInfo) error {
	nsm.mutex.Lock()
	defer nsm.mutex.Unlock()

	// Check if node is pending
	pendingNode, exists := nsm.pendingNodes[nodeID]
	if !exists {
		return fmt.Errorf("node %s is not in pending state", nodeID)
	}

	// Create active node
	activeNode := &ActiveNode{
		NodeID:               nodeID,
		RegisteredAt:         time.Now(),
		Config:               pendingNode.Config,
		SSHKeys:              pendingNode.SSHKeys,
		ConnectionInfo:       connectionInfo,
		HardwareInfo:         hardwareInfo,
		LastSeen:             time.Now(),
		LastHeartbeat:        time.Now(),
		Status:               "active",
		WorkloadCount:        0,
		MaxWorkloads:         calculateMaxWorkloads(hardwareInfo),
		Metrics:              &types.NodeMetrics{},
		HeartbeatFailures:    0,
		MaxHeartbeatFailures: 3,
	}

	// Move from pending to active
	delete(nsm.pendingNodes, nodeID)
	nsm.activeNodes[nodeID] = activeNode

	// Persist changes
	if nsm.persistenceEnabled {
		if err := nsm.saveActiveNodeToDisk(activeNode); err != nil {
			nsm.logger.Error("Failed to save active node to disk", "node_id", nodeID, "error", err)
		}
		if err := nsm.removePendingNodeFromDisk(nodeID); err != nil {
			nsm.logger.Error("Failed to remove pending node from disk", "node_id", nodeID, "error", err)
		}
	}

	nsm.logger.Info("Node activated", "node_id", nodeID)
	return nil
}

// DeactivateNode moves a node from active to inactive state
func (nsm *NodeStateManager) DeactivateNode(nodeID string, reason string, lastError string) error {
	nsm.mutex.Lock()
	defer nsm.mutex.Unlock()

	// Check if node is active
	activeNode, exists := nsm.activeNodes[nodeID]
	if !exists {
		return fmt.Errorf("node %s is not in active state", nodeID)
	}

	// Create inactive node
	inactiveNode := &InactiveNode{
		NodeID:           nodeID,
		LastActiveAt:     activeNode.LastSeen,
		InactiveSince:    time.Now(),
		Config:           activeNode.Config,
		SSHKeys:          activeNode.SSHKeys,
		ConnectionInfo:   activeNode.ConnectionInfo,
		HardwareInfo:     activeNode.HardwareInfo,
		Status:           "disconnected",
		FailureReason:    reason,
		LastError:        lastError,
		RetryAttempts:    0,
		MaxRetryAttempts: 5,
		NextRetryAt:      time.Now().Add(5 * time.Minute),
		RecoveryActions:  []string{"heartbeat_retry", "connection_retry"},
	}

	// Move from active to inactive
	delete(nsm.activeNodes, nodeID)
	nsm.inactiveNodes[nodeID] = inactiveNode

	// Persist changes
	if nsm.persistenceEnabled {
		if err := nsm.saveInactiveNodeToDisk(inactiveNode); err != nil {
			nsm.logger.Error("Failed to save inactive node to disk", "node_id", nodeID, "error", err)
		}
		if err := nsm.removeActiveNodeFromDisk(nodeID); err != nil {
			nsm.logger.Error("Failed to remove active node from disk", "node_id", nodeID, "error", err)
		}
	}

	nsm.logger.Info("Node deactivated", "node_id", nodeID, "reason", reason)
	return nil
}

// UpdateNodeStatus updates the status of a node
func (nsm *NodeStateManager) UpdateNodeStatus(nodeID string, status string) error {
	nsm.mutex.Lock()
	defer nsm.mutex.Unlock()

	// Check if node is active
	if activeNode, exists := nsm.activeNodes[nodeID]; exists {
		activeNode.Status = status
		activeNode.LastSeen = time.Now()

		if nsm.persistenceEnabled {
			if err := nsm.saveActiveNodeToDisk(activeNode); err != nil {
				nsm.logger.Error("Failed to save active node status to disk", "node_id", nodeID, "error", err)
			}
		}

		nsm.logger.Debug("Updated active node status", "node_id", nodeID, "status", status)
		return nil
	}

	// Check if node is pending
	if pendingNode, exists := nsm.pendingNodes[nodeID]; exists {
		pendingNode.Status = status
		pendingNode.LastChecked = time.Now()

		if nsm.persistenceEnabled {
			if err := nsm.savePendingNodeToDisk(pendingNode); err != nil {
				nsm.logger.Error("Failed to save pending node status to disk", "node_id", nodeID, "error", err)
			}
		}

		nsm.logger.Debug("Updated pending node status", "node_id", nodeID, "status", status)
		return nil
	}

	// Check if node is inactive
	if inactiveNode, exists := nsm.inactiveNodes[nodeID]; exists {
		inactiveNode.Status = status

		if nsm.persistenceEnabled {
			if err := nsm.saveInactiveNodeToDisk(inactiveNode); err != nil {
				nsm.logger.Error("Failed to save inactive node status to disk", "node_id", nodeID, "error", err)
			}
		}

		nsm.logger.Debug("Updated inactive node status", "node_id", nodeID, "status", status)
		return nil
	}

	return fmt.Errorf("node %s not found in any state", nodeID)
}

// GetNodeState returns the current state of a node
func (nsm *NodeStateManager) GetNodeState(nodeID string) (*NodeState, error) {
	nsm.mutex.RLock()
	defer nsm.mutex.RUnlock()

	// Check pending nodes
	if _, exists := nsm.pendingNodes[nodeID]; exists {
		return &NodeState{
			NodeID:         nodeID,
			CurrentState:   "pending",
			LastTransition: time.Now(),
		}, nil
	}

	// Check active nodes
	if _, exists := nsm.activeNodes[nodeID]; exists {
		return &NodeState{
			NodeID:         nodeID,
			CurrentState:   "active",
			LastTransition: time.Now(),
		}, nil
	}

	// Check inactive nodes
	if _, exists := nsm.inactiveNodes[nodeID]; exists {
		return &NodeState{
			NodeID:         nodeID,
			CurrentState:   "inactive",
			LastTransition: time.Now(),
		}, nil
	}

	return nil, fmt.Errorf("node %s not found", nodeID)
}

// (removed legacy getters returning internal node structs to satisfy interface)

// IsNodePending checks if a node is in pending state
func (nsm *NodeStateManager) IsNodePending(nodeID string) bool {
	nsm.mutex.RLock()
	defer nsm.mutex.RUnlock()

	_, exists := nsm.pendingNodes[nodeID]
	return exists
}

// IsNodeActive checks if a node is in active state
func (nsm *NodeStateManager) IsNodeActive(nodeID string) bool {
	nsm.mutex.RLock()
	defer nsm.mutex.RUnlock()

	_, exists := nsm.activeNodes[nodeID]
	return exists
}

// IsNodeInactive checks if a node is in inactive state
func (nsm *NodeStateManager) IsNodeInactive(nodeID string) bool {
	nsm.mutex.RLock()
	defer nsm.mutex.RUnlock()

	_, exists := nsm.inactiveNodes[nodeID]
	return exists
}

// GetNodeCounts returns the count of nodes in each state
func (nsm *NodeStateManager) GetNodeCounts() map[string]int {
	nsm.mutex.RLock()
	defer nsm.mutex.RUnlock()

	return map[string]int{
		"pending":  len(nsm.pendingNodes),
		"active":   len(nsm.activeNodes),
		"inactive": len(nsm.inactiveNodes),
		"total":    len(nsm.pendingNodes) + len(nsm.activeNodes) + len(nsm.inactiveNodes),
	}
}

// UpdateNodeHeartbeat updates the heartbeat information for an active node
func (nsm *NodeStateManager) UpdateNodeHeartbeat(nodeID string, metrics *types.NodeMetrics) error {
	nsm.mutex.Lock()
	defer nsm.mutex.Unlock()

	activeNode, exists := nsm.activeNodes[nodeID]
	if !exists {
		return fmt.Errorf("node %s is not in active state", nodeID)
	}

	activeNode.LastHeartbeat = time.Now()
	activeNode.LastSeen = time.Now()
	activeNode.HeartbeatFailures = 0

	if metrics != nil {
		activeNode.Metrics = metrics
	}

	if nsm.persistenceEnabled {
		if err := nsm.saveActiveNodeToDisk(activeNode); err != nil {
			nsm.logger.Error("Failed to save active node heartbeat to disk", "node_id", nodeID, "error", err)
		}
	}

	return nil
}

// IncrementNodeHeartbeatFailures increments the heartbeat failure count for an active node
func (nsm *NodeStateManager) IncrementNodeHeartbeatFailures(nodeID string) error {
	nsm.mutex.Lock()
	defer nsm.mutex.Unlock()

	activeNode, exists := nsm.activeNodes[nodeID]
	if !exists {
		return fmt.Errorf("node %s is not in active state", nodeID)
	}

	activeNode.HeartbeatFailures++

	if nsm.persistenceEnabled {
		if err := nsm.saveActiveNodeToDisk(activeNode); err != nil {
			nsm.logger.Error("Failed to save active node heartbeat failures to disk", "node_id", nodeID, "error", err)
		}
	}

	// Check if node should be deactivated
	if activeNode.HeartbeatFailures >= activeNode.MaxHeartbeatFailures {
		nsm.logger.Warn("Node heartbeat failures exceeded threshold",
			"node_id", nodeID,
			"failures", activeNode.HeartbeatFailures,
			"max_failures", activeNode.MaxHeartbeatFailures)

		// Move to inactive state
		if err := nsm.deactivateNodeInternal(nodeID, "heartbeat_failures", "Too many heartbeat failures"); err != nil {
			nsm.logger.Error("Failed to deactivate node due to heartbeat failures", "node_id", nodeID, "error", err)
		}
	}

	return nil
}

// Private methods

// deactivateNodeInternal deactivates a node without acquiring the lock (internal use)
func (nsm *NodeStateManager) deactivateNodeInternal(nodeID string, reason string, lastError string) error {
	activeNode, exists := nsm.activeNodes[nodeID]
	if !exists {
		return fmt.Errorf("node %s is not in active state", nodeID)
	}

	// Create inactive node
	inactiveNode := &InactiveNode{
		NodeID:           nodeID,
		LastActiveAt:     activeNode.LastSeen,
		InactiveSince:    time.Now(),
		Config:           activeNode.Config,
		SSHKeys:          activeNode.SSHKeys,
		ConnectionInfo:   activeNode.ConnectionInfo,
		HardwareInfo:     activeNode.HardwareInfo,
		Status:           "disconnected",
		FailureReason:    reason,
		LastError:        lastError,
		RetryAttempts:    0,
		MaxRetryAttempts: 5,
		NextRetryAt:      time.Now().Add(5 * time.Minute),
		RecoveryActions:  []string{"heartbeat_retry", "connection_retry"},
	}

	// Move from active to inactive
	delete(nsm.activeNodes, nodeID)
	nsm.inactiveNodes[nodeID] = inactiveNode

	// Persist changes
	if nsm.persistenceEnabled {
		if err := nsm.saveInactiveNodeToDisk(inactiveNode); err != nil {
			nsm.logger.Error("Failed to save inactive node to disk", "node_id", nodeID, "error", err)
		}
		if err := nsm.removeActiveNodeFromDisk(nodeID); err != nil {
			nsm.logger.Error("Failed to remove active node from disk", "node_id", nodeID, "error", err)
		}
	}

	return nil
}

// Persistence methods

// loadStateFromDisk loads the state from disk
func (nsm *NodeStateManager) loadStateFromDisk() error {
	// Load pending nodes
	if err := nsm.loadPendingNodesFromDisk(); err != nil {
		nsm.logger.Error("Failed to load pending nodes from disk", "error", err)
	}

	// Load active nodes
	if err := nsm.loadActiveNodesFromDisk(); err != nil {
		nsm.logger.Error("Failed to load active nodes from disk", "error", err)
	}

	// Load inactive nodes
	if err := nsm.loadInactiveNodesFromDisk(); err != nil {
		nsm.logger.Error("Failed to load inactive nodes from disk", "error", err)
	}

	return nil
}

// loadPendingNodesFromDisk loads pending nodes from disk
func (nsm *NodeStateManager) loadPendingNodesFromDisk() error {
	pendingDir := filepath.Join(nsm.stateDir, "pending")
	if _, err := os.Stat(pendingDir); os.IsNotExist(err) {
		return nil // No pending nodes
	}

	files, err := filepath.Glob(filepath.Join(pendingDir, "*.json"))
	if err != nil {
		return err
	}

	for _, file := range files {
		var pendingNode PendingNode
		if err := helpers.ReadJSONFile(file, &pendingNode); err != nil {
			nsm.logger.Error("Failed to read pending node file", "file", file, "error", err)
			continue
		}
		nsm.pendingNodes[pendingNode.NodeID] = &pendingNode
	}

	return nil
}

// loadActiveNodesFromDisk loads active nodes from disk
func (nsm *NodeStateManager) loadActiveNodesFromDisk() error {
	activeDir := filepath.Join(nsm.stateDir, "active")
	if _, err := os.Stat(activeDir); os.IsNotExist(err) {
		return nil // No active nodes
	}

	files, err := filepath.Glob(filepath.Join(activeDir, "*.json"))
	if err != nil {
		return err
	}

	for _, file := range files {
		var activeNode ActiveNode
		if err := helpers.ReadJSONFile(file, &activeNode); err != nil {
			nsm.logger.Error("Failed to read active node file", "file", file, "error", err)
			continue
		}
		nsm.activeNodes[activeNode.NodeID] = &activeNode
	}

	return nil
}

// loadInactiveNodesFromDisk loads inactive nodes from disk
func (nsm *NodeStateManager) loadInactiveNodesFromDisk() error {
	inactiveDir := filepath.Join(nsm.stateDir, "inactive")
	if _, err := os.Stat(inactiveDir); os.IsNotExist(err) {
		return nil // No inactive nodes
	}

	files, err := filepath.Glob(filepath.Join(inactiveDir, "*.json"))
	if err != nil {
		return err
	}

	for _, file := range files {
		var inactiveNode InactiveNode
		if err := helpers.ReadJSONFile(file, &inactiveNode); err != nil {
			nsm.logger.Error("Failed to read inactive node file", "file", file, "error", err)
			continue
		}
		nsm.inactiveNodes[inactiveNode.NodeID] = &inactiveNode
	}

	return nil
}

// savePendingNodeToDisk saves a pending node to disk
func (nsm *NodeStateManager) savePendingNodeToDisk(node *PendingNode) error {
	pendingDir := filepath.Join(nsm.stateDir, "pending")
	if err := os.MkdirAll(pendingDir, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(pendingDir, node.NodeID+".json")
	return helpers.WriteJSONFile(filePath, node)
}

// saveActiveNodeToDisk saves an active node to disk
func (nsm *NodeStateManager) saveActiveNodeToDisk(node *ActiveNode) error {
	activeDir := filepath.Join(nsm.stateDir, "active")
	if err := os.MkdirAll(activeDir, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(activeDir, node.NodeID+".json")
	return helpers.WriteJSONFile(filePath, node)
}

// saveInactiveNodeToDisk saves an inactive node to disk
func (nsm *NodeStateManager) saveInactiveNodeToDisk(node *InactiveNode) error {
	inactiveDir := filepath.Join(nsm.stateDir, "inactive")
	if err := os.MkdirAll(inactiveDir, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(inactiveDir, node.NodeID+".json")
	return helpers.WriteJSONFile(filePath, node)
}

// removePendingNodeFromDisk removes a pending node from disk
func (nsm *NodeStateManager) removePendingNodeFromDisk(nodeID string) error {
	filePath := filepath.Join(nsm.stateDir, "pending", nodeID+".json")
	return os.Remove(filePath)
}

// removeActiveNodeFromDisk removes an active node from disk
func (nsm *NodeStateManager) removeActiveNodeFromDisk(nodeID string) error {
	filePath := filepath.Join(nsm.stateDir, "active", nodeID+".json")
	return os.Remove(filePath)
}

// removeInactiveNodeFromDisk removes an inactive node from disk
func (nsm *NodeStateManager) removeInactiveNodeFromDisk(nodeID string) error {
	filePath := filepath.Join(nsm.stateDir, "inactive", nodeID+".json")
	return os.Remove(filePath)
}

// Helper function to calculate max workloads based on hardware
func calculateMaxWorkloads(hardwareInfo *types.HardwareInfo) int {
	if hardwareInfo == nil {
		return 1
	}

	// Simple calculation based on CPU cores and memory
	// Each workload should have at least 1 CPU core and 2GB RAM
	maxByCPU := hardwareInfo.CPUCores
	maxByMemory := hardwareInfo.MemoryGB / 2

	// Return the minimum of the two
	if maxByCPU < maxByMemory {
		return maxByCPU
	}
	return maxByMemory
}
