package node

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/syntropy-grid/syntropy-cooperative-grid/manager/interfaces/cli/node/src/internal/types"
)

// HeartbeatManager manages heartbeat communication with nodes
type HeartbeatManager struct {
	nodeStateManager types.NodeStateManager
	logger           types.Logger
	ctx              context.Context
	cancel           context.CancelFunc
	mutex            sync.RWMutex
	heartbeats       map[string]*NodeHeartbeat
	interval         time.Duration
	timeout          time.Duration
	maxFailures      int
	running          bool
}

// NodeHeartbeat represents heartbeat state for a node
type NodeHeartbeat struct {
	NodeID           string          `json:"node_id"`
	LastSeen         time.Time       `json:"last_seen"`
	LastResponse     time.Time       `json:"last_response"`
	FailureCount     int             `json:"failure_count"`
	Status           string          `json:"status"` // active, inactive, failed
	ConnectionInfo   *ConnectionInfo `json:"connection_info"`
	Metrics          *NodeMetrics    `json:"metrics"`
	LastError        string          `json:"last_error,omitempty"`
	NextHeartbeat    time.Time       `json:"next_heartbeat"`
	RetryAttempts    int             `json:"retry_attempts"`
	MaxRetryAttempts int             `json:"max_retry_attempts"`
}

// ConnectionInfo represents connection information for a node
type ConnectionInfo struct {
	RemoteAddr  string        `json:"remote_addr"`
	Port        int           `json:"port"`
	Protocol    string        `json:"protocol"`
	ConnectedAt time.Time     `json:"connected_at"`
	LastPing    time.Time     `json:"last_ping"`
	Latency     time.Duration `json:"latency"`
}

// NodeMetrics represents metrics collected from a node
type NodeMetrics struct {
	CPUUsage      float64       `json:"cpu_usage"`
	MemoryUsage   float64       `json:"memory_usage"`
	DiskUsage     float64       `json:"disk_usage"`
	NetworkIn     int64         `json:"network_in"`
	NetworkOut    int64         `json:"network_out"`
	Uptime        time.Duration `json:"uptime"`
	WorkloadCount int           `json:"workload_count"`
	LastUpdated   time.Time     `json:"last_updated"`
}

// HeartbeatRequest represents a heartbeat request
type HeartbeatRequest struct {
	Type      string       `json:"type"`
	NodeID    string       `json:"node_id"`
	Timestamp time.Time    `json:"timestamp"`
	Metrics   *NodeMetrics `json:"metrics,omitempty"`
}

// HeartbeatResponse represents a heartbeat response
type HeartbeatResponse struct {
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	NextPing  time.Time `json:"next_ping"`
}

// NewHeartbeatManager creates a new heartbeat manager
func NewHeartbeatManager(
	nodeStateManager types.NodeStateManager,
	logger types.Logger,
) *HeartbeatManager {
	ctx, cancel := context.WithCancel(context.Background())

	return &HeartbeatManager{
		nodeStateManager: nodeStateManager,
		logger:           logger,
		ctx:              ctx,
		cancel:           cancel,
		heartbeats:       make(map[string]*NodeHeartbeat),
		interval:         30 * time.Second,
		timeout:          10 * time.Second,
		maxFailures:      3,
	}
}

// Start starts the heartbeat manager
func (hm *HeartbeatManager) Start() error {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()

	if hm.running {
		return fmt.Errorf("heartbeat manager is already running")
	}

	hm.logger.Info("Starting heartbeat manager",
		"interval", hm.interval,
		"timeout", hm.timeout,
		"max_failures", hm.maxFailures)

	hm.running = true

	// Start heartbeat processing in a goroutine
	go hm.processHeartbeats()

	hm.logger.Info("Heartbeat manager started successfully")
	return nil
}

// Stop stops the heartbeat manager
func (hm *HeartbeatManager) Stop() error {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()

	if !hm.running {
		return fmt.Errorf("heartbeat manager is not running")
	}

	hm.logger.Info("Stopping heartbeat manager")

	hm.running = false
	hm.cancel()

	hm.logger.Info("Heartbeat manager stopped")
	return nil
}

// IsRunning returns whether the heartbeat manager is running
func (hm *HeartbeatManager) IsRunning() bool {
	hm.mutex.RLock()
	defer hm.mutex.RUnlock()
	return hm.running
}

// AddNode adds a node to heartbeat monitoring
func (hm *HeartbeatManager) AddNode(nodeID string, connectionInfo *ConnectionInfo) error {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()

	if _, exists := hm.heartbeats[nodeID]; exists {
		return fmt.Errorf("node %s is already being monitored", nodeID)
	}

	heartbeat := &NodeHeartbeat{
		NodeID:           nodeID,
		LastSeen:         time.Now(),
		LastResponse:     time.Now(),
		FailureCount:     0,
		Status:           "active",
		ConnectionInfo:   connectionInfo,
		Metrics:          &NodeMetrics{},
		NextHeartbeat:    time.Now().Add(hm.interval),
		RetryAttempts:    0,
		MaxRetryAttempts: 5,
	}

	hm.heartbeats[nodeID] = heartbeat

	hm.logger.Info("Node added to heartbeat monitoring", "node_id", nodeID)
	return nil
}

// RemoveNode removes a node from heartbeat monitoring
func (hm *HeartbeatManager) RemoveNode(nodeID string) error {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()

	if _, exists := hm.heartbeats[nodeID]; !exists {
		return fmt.Errorf("node %s is not being monitored", nodeID)
	}

	delete(hm.heartbeats, nodeID)

	hm.logger.Info("Node removed from heartbeat monitoring", "node_id", nodeID)
	return nil
}

// ProcessHeartbeat processes an incoming heartbeat from a node
func (hm *HeartbeatManager) ProcessHeartbeat(ctx context.Context, conn net.Conn, request *HeartbeatRequest) (*HeartbeatResponse, error) {
	hm.logger.Debug("Processing heartbeat", "node_id", request.NodeID)

	// Update heartbeat state
	hm.mutex.Lock()
	heartbeat, exists := hm.heartbeats[request.NodeID]
	if !exists {
		hm.mutex.Unlock()
		return nil, fmt.Errorf("node %s is not being monitored", request.NodeID)
	}

	// Update heartbeat information
	heartbeat.LastSeen = time.Now()
	heartbeat.LastResponse = time.Now()
	heartbeat.FailureCount = 0
	heartbeat.Status = "active"
	heartbeat.RetryAttempts = 0
	heartbeat.LastError = ""

	// Update metrics if provided
	if request.Metrics != nil {
		heartbeat.Metrics = request.Metrics
		heartbeat.Metrics.LastUpdated = time.Now()
	}

	// Update connection info
	if heartbeat.ConnectionInfo != nil {
		heartbeat.ConnectionInfo.LastPing = time.Now()
		heartbeat.ConnectionInfo.Latency = time.Since(request.Timestamp)
	}

	// Calculate next heartbeat time
	heartbeat.NextHeartbeat = time.Now().Add(hm.interval)

	hm.mutex.Unlock()

	// Update node state in node state manager
	if err := hm.nodeStateManager.UpdateNodeStatus(request.NodeID, "active"); err != nil {
		hm.logger.Error("Failed to update node status", "node_id", request.NodeID, "error", err)
	}

	// Create response
	response := &HeartbeatResponse{
		Status:    "ack",
		Message:   "Heartbeat received",
		Timestamp: time.Now(),
		NextPing:  heartbeat.NextHeartbeat,
	}

	hm.logger.Debug("Heartbeat processed successfully", "node_id", request.NodeID)
	return response, nil
}

// GetHeartbeatStatus returns the heartbeat status for a node
func (hm *HeartbeatManager) GetHeartbeatStatus(nodeID string) (*NodeHeartbeat, error) {
	hm.mutex.RLock()
	defer hm.mutex.RUnlock()

	heartbeat, exists := hm.heartbeats[nodeID]
	if !exists {
		return nil, fmt.Errorf("node %s is not being monitored", nodeID)
	}

	return heartbeat, nil
}

// GetAllHeartbeatStatuses returns heartbeat statuses for all nodes
func (hm *HeartbeatManager) GetAllHeartbeatStatuses() map[string]*NodeHeartbeat {
	hm.mutex.RLock()
	defer hm.mutex.RUnlock()

	statuses := make(map[string]*NodeHeartbeat)
	for nodeID, heartbeat := range hm.heartbeats {
		statuses[nodeID] = heartbeat
	}

	return statuses
}

// GetActiveNodes returns list of active nodes
func (hm *HeartbeatManager) GetActiveNodes() []string {
	hm.mutex.RLock()
	defer hm.mutex.RUnlock()

	var activeNodes []string
	for nodeID, heartbeat := range hm.heartbeats {
		if heartbeat.Status == "active" {
			activeNodes = append(activeNodes, nodeID)
		}
	}

	return activeNodes
}

// GetInactiveNodes returns list of inactive nodes
func (hm *HeartbeatManager) GetInactiveNodes() []string {
	hm.mutex.RLock()
	defer hm.mutex.RUnlock()

	var inactiveNodes []string
	for nodeID, heartbeat := range hm.heartbeats {
		if heartbeat.Status == "inactive" || heartbeat.Status == "failed" {
			inactiveNodes = append(inactiveNodes, nodeID)
		}
	}

	return inactiveNodes
}

// Private methods

// processHeartbeats processes heartbeats for all monitored nodes
func (hm *HeartbeatManager) processHeartbeats() {
	hm.logger.Info("Starting heartbeat processing")

	ticker := time.NewTicker(hm.interval)
	defer ticker.Stop()

	for {
		select {
		case <-hm.ctx.Done():
			hm.logger.Debug("Heartbeat processing context cancelled")
			return
		case <-ticker.C:
			hm.checkHeartbeats()
		}
	}
}

// checkHeartbeats checks heartbeat status for all nodes
func (hm *HeartbeatManager) checkHeartbeats() {
	hm.mutex.RLock()
	nodes := make([]string, 0, len(hm.heartbeats))
	for nodeID := range hm.heartbeats {
		nodes = append(nodes, nodeID)
	}
	hm.mutex.RUnlock()

	for _, nodeID := range nodes {
		hm.checkNodeHeartbeat(nodeID)
	}
}

// checkNodeHeartbeat checks heartbeat status for a specific node
func (hm *HeartbeatManager) checkNodeHeartbeat(nodeID string) {
	hm.mutex.Lock()
	heartbeat, exists := hm.heartbeats[nodeID]
	if !exists {
		hm.mutex.Unlock()
		return
	}

	// Check if heartbeat is overdue
	now := time.Now()
	if now.After(heartbeat.NextHeartbeat) {
		// Heartbeat is overdue
		heartbeat.FailureCount++
		heartbeat.LastError = fmt.Sprintf("Heartbeat overdue by %v", now.Sub(heartbeat.NextHeartbeat))

		hm.logger.Warn("Heartbeat overdue",
			"node_id", nodeID,
			"failure_count", heartbeat.FailureCount,
			"overdue_by", now.Sub(heartbeat.NextHeartbeat))

		// Check if node should be marked as inactive
		if heartbeat.FailureCount >= hm.maxFailures {
			heartbeat.Status = "inactive"
			hm.logger.Error("Node marked as inactive due to heartbeat failures",
				"node_id", nodeID,
				"failure_count", heartbeat.FailureCount)

			// Update node state
			if err := hm.nodeStateManager.UpdateNodeStatus(nodeID, "inactive"); err != nil {
				hm.logger.Error("Failed to update node status to inactive", "node_id", nodeID, "error", err)
			}
		}

		// Schedule retry
		heartbeat.NextHeartbeat = now.Add(hm.interval)
	}

	hm.mutex.Unlock()
}

// HeartbeatServer handles incoming heartbeat connections
type HeartbeatServer struct {
	heartbeatManager *HeartbeatManager
	listener         net.Listener
	logger           types.Logger
	port             int
	running          bool
	ctx              context.Context
	cancel           context.CancelFunc
}

// NewHeartbeatServer creates a new heartbeat server
func NewHeartbeatServer(
	heartbeatManager *HeartbeatManager,
	logger types.Logger,
	port int,
) *HeartbeatServer {
	ctx, cancel := context.WithCancel(context.Background())

	return &HeartbeatServer{
		heartbeatManager: heartbeatManager,
		logger:           logger,
		port:             port,
		ctx:              ctx,
		cancel:           cancel,
	}
}

// Start starts the heartbeat server
func (hs *HeartbeatServer) Start() error {
	hs.logger.Info("Starting heartbeat server", "port", hs.port)

	// Create listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", hs.port))
	if err != nil {
		return fmt.Errorf("failed to create heartbeat listener: %w", err)
	}

	hs.listener = listener
	hs.running = true

	// Start accepting connections in a goroutine
	go hs.acceptConnections()

	hs.logger.Info("Heartbeat server started successfully", "port", hs.port)
	return nil
}

// Stop stops the heartbeat server
func (hs *HeartbeatServer) Stop() error {
	hs.logger.Info("Stopping heartbeat server")

	hs.running = false
	hs.cancel()

	if hs.listener != nil {
		if err := hs.listener.Close(); err != nil {
			hs.logger.Error("Error closing heartbeat listener", "error", err)
		}
	}

	hs.logger.Info("Heartbeat server stopped")
	return nil
}

// IsRunning returns whether the heartbeat server is running
func (hs *HeartbeatServer) IsRunning() bool {
	return hs.running
}

// acceptConnections accepts incoming heartbeat connections
func (hs *HeartbeatServer) acceptConnections() {
	for hs.running {
		select {
		case <-hs.ctx.Done():
			hs.logger.Debug("Heartbeat server context cancelled")
			return
		default:
			// Accept connection with timeout
			conn, err := hs.listener.Accept()
			if err != nil {
				if hs.running {
					hs.logger.Error("Failed to accept heartbeat connection", "error", err)
				}
				continue
			}

			// Handle connection in a goroutine
			go hs.handleConnection(conn)
		}
	}
}

// handleConnection handles a single heartbeat connection
func (hs *HeartbeatServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	remoteAddr := conn.RemoteAddr().String()
	hs.logger.Debug("Handling heartbeat connection", "remote_addr", remoteAddr)

	// Set connection timeout
	conn.SetDeadline(time.Now().Add(30 * time.Second))

	// Read heartbeat request
	var request HeartbeatRequest
	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(&request); err != nil {
		hs.logger.Error("Failed to decode heartbeat request", "remote_addr", remoteAddr, "error", err)
		return
	}

	// Process heartbeat
	response, err := hs.heartbeatManager.ProcessHeartbeat(hs.ctx, conn, &request)
	if err != nil {
		hs.logger.Error("Failed to process heartbeat", "remote_addr", remoteAddr, "error", err)
		return
	}

	// Send response
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(response); err != nil {
		hs.logger.Error("Failed to send heartbeat response", "remote_addr", remoteAddr, "error", err)
		return
	}

	hs.logger.Debug("Heartbeat connection handled successfully",
		"remote_addr", remoteAddr,
		"node_id", request.NodeID)
}

// HeartbeatClient sends heartbeats to the command station
type HeartbeatClient struct {
	commandStationAddr string
	nodeID             string
	interval           time.Duration
	timeout            time.Duration
	logger             types.Logger
	ctx                context.Context
	cancel             context.CancelFunc
	running            bool
	mutex              sync.RWMutex
}

// NewHeartbeatClient creates a new heartbeat client
func NewHeartbeatClient(
	commandStationAddr string,
	nodeID string,
	logger types.Logger,
) *HeartbeatClient {
	ctx, cancel := context.WithCancel(context.Background())

	return &HeartbeatClient{
		commandStationAddr: commandStationAddr,
		nodeID:             nodeID,
		interval:           30 * time.Second,
		timeout:            10 * time.Second,
		logger:             logger,
		ctx:                ctx,
		cancel:             cancel,
	}
}

// Start starts the heartbeat client
func (hc *HeartbeatClient) Start() error {
	hc.mutex.Lock()
	defer hc.mutex.Unlock()

	if hc.running {
		return fmt.Errorf("heartbeat client is already running")
	}

	hc.logger.Info("Starting heartbeat client",
		"node_id", hc.nodeID,
		"command_station", hc.commandStationAddr,
		"interval", hc.interval)

	hc.running = true

	// Start sending heartbeats in a goroutine
	go hc.sendHeartbeats()

	hc.logger.Info("Heartbeat client started successfully")
	return nil
}

// Stop stops the heartbeat client
func (hc *HeartbeatClient) Stop() error {
	hc.mutex.Lock()
	defer hc.mutex.Unlock()

	if !hc.running {
		return fmt.Errorf("heartbeat client is not running")
	}

	hc.logger.Info("Stopping heartbeat client")

	hc.running = false
	hc.cancel()

	hc.logger.Info("Heartbeat client stopped")
	return nil
}

// IsRunning returns whether the heartbeat client is running
func (hc *HeartbeatClient) IsRunning() bool {
	hc.mutex.RLock()
	defer hc.mutex.RUnlock()
	return hc.running
}

// sendHeartbeats sends periodic heartbeats to the command station
func (hc *HeartbeatClient) sendHeartbeats() {
	hc.logger.Info("Starting to send heartbeats")

	ticker := time.NewTicker(hc.interval)
	defer ticker.Stop()

	for {
		select {
		case <-hc.ctx.Done():
			hc.logger.Debug("Heartbeat client context cancelled")
			return
		case <-ticker.C:
			hc.sendHeartbeat()
		}
	}
}

// sendHeartbeat sends a single heartbeat to the command station
func (hc *HeartbeatClient) sendHeartbeat() {
	hc.logger.Debug("Sending heartbeat", "node_id", hc.nodeID)

	// Create heartbeat request
	request := &HeartbeatRequest{
		Type:      "heartbeat",
		NodeID:    hc.nodeID,
		Timestamp: time.Now(),
		Metrics:   hc.collectMetrics(),
	}

	// Send heartbeat
	if err := hc.sendHeartbeatRequest(request); err != nil {
		hc.logger.Error("Failed to send heartbeat", "node_id", hc.nodeID, "error", err)
		return
	}

	hc.logger.Debug("Heartbeat sent successfully", "node_id", hc.nodeID)
}

// sendHeartbeatRequest sends a heartbeat request to the command station
func (hc *HeartbeatClient) sendHeartbeatRequest(request *HeartbeatRequest) error {
	// Connect to command station
	conn, err := net.DialTimeout("tcp", hc.commandStationAddr, hc.timeout)
	if err != nil {
		return fmt.Errorf("failed to connect to command station: %w", err)
	}
	defer conn.Close()

	// Set connection timeout
	conn.SetDeadline(time.Now().Add(hc.timeout))

	// Send request
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(request); err != nil {
		return fmt.Errorf("failed to send heartbeat request: %w", err)
	}

	// Read response
	var response HeartbeatResponse
	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(&response); err != nil {
		return fmt.Errorf("failed to read heartbeat response: %w", err)
	}

	if response.Status != "ack" {
		return fmt.Errorf("heartbeat response status: %s", response.Status)
	}

	return nil
}

// collectMetrics collects basic metrics from the node
func (hc *HeartbeatClient) collectMetrics() *NodeMetrics {
	// This is a simplified implementation
	// In production, we would collect real system metrics
	return &NodeMetrics{
		CPUUsage:      0.0,                    // Placeholder
		MemoryUsage:   0.0,                    // Placeholder
		DiskUsage:     0.0,                    // Placeholder
		NetworkIn:     0,                      // Placeholder
		NetworkOut:    0,                      // Placeholder
		Uptime:        time.Since(time.Now()), // Placeholder
		WorkloadCount: 0,                      // Placeholder
		LastUpdated:   time.Now(),
	}
}
