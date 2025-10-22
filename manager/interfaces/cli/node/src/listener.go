package node

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"node-component/src/internal/types"
)

// Listener manages TCP connections for node registration
type Listener struct {
	handshakeManager *HandshakeManager
	port             int
	host             string
	listener         net.Listener
	running          bool
	ctx              context.Context
	cancel           context.CancelFunc
	logger           types.Logger
	mutex            sync.RWMutex
	connections      map[string]*Connection
	stats            *ListenerStats
}

// Connection represents an active connection
type Connection struct {
	Conn         net.Conn
	RemoteAddr   string
	StartTime    time.Time
	LastActivity time.Time
	NodeID       string
	Status       string
}

// ListenerStats represents statistics about the listener
type ListenerStats struct {
	TotalConnections     int           `json:"total_connections"`
	ActiveConnections    int           `json:"active_connections"`
	SuccessfulHandshakes int           `json:"successful_handshakes"`
	FailedHandshakes     int           `json:"failed_handshakes"`
	TotalBytesReceived   int64         `json:"total_bytes_received"`
	TotalBytesSent       int64         `json:"total_bytes_sent"`
	Uptime               time.Duration `json:"uptime"`
	StartTime            time.Time     `json:"start_time"`
}

// NewListener creates a new listener
func NewListener(
	handshakeManager *HandshakeManager,
	logger types.Logger,
	port int,
) *Listener {
	ctx, cancel := context.WithCancel(context.Background())

	return &Listener{
		handshakeManager: handshakeManager,
		port:             port,
		host:             "0.0.0.0",
		ctx:              ctx,
		cancel:           cancel,
		logger:           logger,
		connections:      make(map[string]*Connection),
		stats: &ListenerStats{
			StartTime: time.Now(),
		},
	}
}

// Start starts the listener
func (l *Listener) Start() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.running {
		return fmt.Errorf("listener is already running")
	}

	l.logger.Info("Starting listener", "host", l.host, "port", l.port)

	// Create listener
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", l.host, l.port))
	if err != nil {
		return fmt.Errorf("failed to create listener: %w", err)
	}

	l.listener = listener
	l.running = true

	// Start accepting connections in a goroutine
	go l.acceptConnections()

	l.logger.Info("Listener started successfully", "host", l.host, "port", l.port)
	return nil
}

// Stop stops the listener
func (l *Listener) Stop() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if !l.running {
		return fmt.Errorf("listener is not running")
	}

	l.logger.Info("Stopping listener")

	l.running = false
	l.cancel()

	// Close all active connections
	l.closeAllConnections()

	// Close listener
	if l.listener != nil {
		if err := l.listener.Close(); err != nil {
			l.logger.Error("Error closing listener", "error", err)
		}
	}

	l.logger.Info("Listener stopped")
	return nil
}

// IsRunning returns whether the listener is running
func (l *Listener) IsRunning() bool {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.running
}

// GetStats returns listener statistics
func (l *Listener) GetStats() *ListenerStats {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	stats := *l.stats
	stats.Uptime = time.Since(l.stats.StartTime)
	stats.ActiveConnections = len(l.connections)

	return &stats
}

// GetConnections returns information about active connections
func (l *Listener) GetConnections() []*Connection {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	connections := make([]*Connection, 0, len(l.connections))
	for _, conn := range l.connections {
		connections = append(connections, conn)
	}

	return connections
}

// Private methods

// acceptConnections accepts incoming connections
func (l *Listener) acceptConnections() {
	l.logger.Info("Accepting connections on listener")

	for l.isRunning() {
		select {
		case <-l.ctx.Done():
			l.logger.Debug("Listener context cancelled")
			return
		default:
			// Accept connection with timeout
			conn, err := l.listener.Accept()
			if err != nil {
				if l.isRunning() {
					l.logger.Error("Failed to accept connection", "error", err)
				}
				continue
			}

			// Handle connection in a goroutine
			go l.handleConnection(conn)
		}
	}
}

// handleConnection handles a single connection
func (l *Listener) handleConnection(conn net.Conn) {
	remoteAddr := conn.RemoteAddr().String()

	l.logger.Info("New connection received", "remote_addr", remoteAddr)

	// Create connection record
	connection := &Connection{
		Conn:         conn,
		RemoteAddr:   remoteAddr,
		StartTime:    time.Now(),
		LastActivity: time.Now(),
		Status:       "connecting",
	}

	// Add to connections map
	l.addConnection(remoteAddr, connection)
	defer l.removeConnection(remoteAddr)
	defer conn.Close()

	// Set connection timeout
	conn.SetDeadline(time.Now().Add(60 * time.Minute))

	// Process handshake
	if err := l.processHandshake(connection); err != nil {
		l.logger.Error("Handshake failed", "remote_addr", remoteAddr, "error", err)
		l.stats.FailedHandshakes++
		return
	}

	l.logger.Info("Handshake successful", "remote_addr", remoteAddr, "node_id", connection.NodeID)
	l.stats.SuccessfulHandshakes++
}

// processHandshake processes the handshake for a connection
func (l *Listener) processHandshake(connection *Connection) error {
	conn := connection.Conn

	// Read handshake request
	var request HandshakeRequest
	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(&request); err != nil {
		return fmt.Errorf("failed to decode handshake request: %w", err)
	}

	// Update connection info
	connection.NodeID = request.NodeID
	connection.Status = "handshaking"
	connection.LastActivity = time.Now()

	// Update stats
	l.stats.TotalConnections++

	// Process handshake using handshake manager
	response, err := l.handshakeManager.ProcessHandshake(l.ctx, conn, &request)
	if err != nil {
		return fmt.Errorf("handshake processing failed: %w", err)
	}

	// Send response
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(response); err != nil {
		return fmt.Errorf("failed to send handshake response: %w", err)
	}

	// Update connection status
	connection.Status = "completed"
	connection.LastActivity = time.Now()

	return nil
}

// addConnection adds a connection to the connections map
func (l *Listener) addConnection(remoteAddr string, connection *Connection) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.connections[remoteAddr] = connection
}

// removeConnection removes a connection from the connections map
func (l *Listener) removeConnection(remoteAddr string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	delete(l.connections, remoteAddr)
}

// closeAllConnections closes all active connections
func (l *Listener) closeAllConnections() {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	for remoteAddr, connection := range l.connections {
		l.logger.Debug("Closing connection", "remote_addr", remoteAddr)
		if err := connection.Conn.Close(); err != nil {
			l.logger.Error("Error closing connection", "remote_addr", remoteAddr, "error", err)
		}
	}

	l.connections = make(map[string]*Connection)
}

// isRunning checks if the listener is running (without lock)
func (l *Listener) isRunning() bool {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.running
}

// ListenerManager manages multiple listeners
type ListenerManager struct {
	listeners map[int]*Listener
	logger    types.Logger
	mutex     sync.RWMutex
}

// NewListenerManager creates a new listener manager
func NewListenerManager(logger types.Logger) *ListenerManager {
	return &ListenerManager{
		listeners: make(map[int]*Listener),
		logger:    logger,
	}
}

// CreateListener creates a new listener on the specified port
func (lm *ListenerManager) CreateListener(
	handshakeManager *HandshakeManager,
	port int,
) (*Listener, error) {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()

	// Check if port is already in use
	if _, exists := lm.listeners[port]; exists {
		return nil, fmt.Errorf("listener already exists on port %d", port)
	}

	// Create new listener
	listener := NewListener(handshakeManager, lm.logger, port)
	lm.listeners[port] = listener

	lm.logger.Info("Listener created", "port", port)
	return listener, nil
}

// StartListener starts a listener on the specified port
func (lm *ListenerManager) StartListener(port int) error {
	lm.mutex.RLock()
	listener, exists := lm.listeners[port]
	lm.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("listener not found on port %d", port)
	}

	return listener.Start()
}

// StopListener stops a listener on the specified port
func (lm *ListenerManager) StopListener(port int) error {
	lm.mutex.RLock()
	listener, exists := lm.listeners[port]
	lm.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("listener not found on port %d", port)
	}

	return listener.Stop()
}

// StopAllListeners stops all listeners
func (lm *ListenerManager) StopAllListeners() error {
	lm.mutex.RLock()
	ports := make([]int, 0, len(lm.listeners))
	for port := range lm.listeners {
		ports = append(ports, port)
	}
	lm.mutex.RUnlock()

	var errors []error
	for _, port := range ports {
		if err := lm.StopListener(port); err != nil {
			errors = append(errors, fmt.Errorf("failed to stop listener on port %d: %w", port, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to stop some listeners: %v", errors)
	}

	return nil
}

// GetListener returns a listener by port
func (lm *ListenerManager) GetListener(port int) (*Listener, error) {
	lm.mutex.RLock()
	defer lm.mutex.RUnlock()

	listener, exists := lm.listeners[port]
	if !exists {
		return nil, fmt.Errorf("listener not found on port %d", port)
	}

	return listener, nil
}

// GetAllListeners returns all listeners
func (lm *ListenerManager) GetAllListeners() map[int]*Listener {
	lm.mutex.RLock()
	defer lm.mutex.RUnlock()

	listeners := make(map[int]*Listener)
	for port, listener := range lm.listeners {
		listeners[port] = listener
	}

	return listeners
}

// GetListenerStats returns statistics for all listeners
func (lm *ListenerManager) GetListenerStats() map[int]*ListenerStats {
	lm.mutex.RLock()
	defer lm.mutex.RUnlock()

	stats := make(map[int]*ListenerStats)
	for port, listener := range lm.listeners {
		stats[port] = listener.GetStats()
	}

	return stats
}

// IsListenerRunning checks if a listener is running on the specified port
func (lm *ListenerManager) IsListenerRunning(port int) bool {
	lm.mutex.RLock()
	defer lm.mutex.RUnlock()

	listener, exists := lm.listeners[port]
	if !exists {
		return false
	}

	return listener.IsRunning()
}

// GetActiveConnections returns active connections for all listeners
func (lm *ListenerManager) GetActiveConnections() map[int][]*Connection {
	lm.mutex.RLock()
	defer lm.mutex.RUnlock()

	connections := make(map[int][]*Connection)
	for port, listener := range lm.listeners {
		connections[port] = listener.GetConnections()
	}

	return connections
}

// AutoListener automatically manages listeners for node registration
type AutoListener struct {
	listenerManager  *ListenerManager
	handshakeManager *HandshakeManager
	defaultPort      int
	logger           types.Logger
	running          bool
	mutex            sync.RWMutex
}

// NewAutoListener creates a new auto listener
func NewAutoListener(
	listenerManager *ListenerManager,
	handshakeManager *HandshakeManager,
	logger types.Logger,
) *AutoListener {
	return &AutoListener{
		listenerManager:  listenerManager,
		handshakeManager: handshakeManager,
		defaultPort:      51000,
		logger:           logger,
	}
}

// Start starts the auto listener
func (al *AutoListener) Start() error {
	al.mutex.Lock()
	defer al.mutex.Unlock()

	if al.running {
		return fmt.Errorf("auto listener is already running")
	}

	al.logger.Info("Starting auto listener", "default_port", al.defaultPort)

	// Create and start default listener
	listener, err := al.listenerManager.CreateListener(al.handshakeManager, al.defaultPort)
	if err != nil {
		return fmt.Errorf("failed to create default listener: %w", err)
	}

	if err := listener.Start(); err != nil {
		return fmt.Errorf("failed to start default listener: %w", err)
	}

	al.running = true
	al.logger.Info("Auto listener started successfully")
	return nil
}

// Stop stops the auto listener
func (al *AutoListener) Stop() error {
	al.mutex.Lock()
	defer al.mutex.Unlock()

	if !al.running {
		return fmt.Errorf("auto listener is not running")
	}

	al.logger.Info("Stopping auto listener")

	// Stop all listeners
	if err := al.listenerManager.StopAllListeners(); err != nil {
		al.logger.Error("Error stopping listeners", "error", err)
	}

	al.running = false
	al.logger.Info("Auto listener stopped")
	return nil
}

// IsRunning returns whether the auto listener is running
func (al *AutoListener) IsRunning() bool {
	al.mutex.RLock()
	defer al.mutex.RUnlock()
	return al.running
}

// GetStats returns statistics for the auto listener
func (al *AutoListener) GetStats() map[int]*ListenerStats {
	return al.listenerManager.GetListenerStats()
}

// GetActiveConnections returns active connections
func (al *AutoListener) GetActiveConnections() map[int][]*Connection {
	return al.listenerManager.GetActiveConnections()
}

// StartListenerForNode starts a listener for a specific node
func (al *AutoListener) StartListenerForNode(nodeID string, port int) error {
	al.logger.Info("Starting listener for node", "node_id", nodeID, "port", port)

	// Create listener for the node
	listener, err := al.listenerManager.CreateListener(al.handshakeManager, port)
	if err != nil {
		return fmt.Errorf("failed to create listener for node %s: %w", nodeID, err)
	}

	// Start the listener
	if err := listener.Start(); err != nil {
		return fmt.Errorf("failed to start listener for node %s: %w", nodeID, err)
	}

	al.logger.Info("Listener started for node", "node_id", nodeID, "port", port)
	return nil
}

// StopListenerForNode stops a listener for a specific node
func (al *AutoListener) StopListenerForNode(nodeID string, port int) error {
	al.logger.Info("Stopping listener for node", "node_id", nodeID, "port", port)

	// Stop the listener
	if err := al.listenerManager.StopListener(port); err != nil {
		return fmt.Errorf("failed to stop listener for node %s: %w", nodeID, err)
	}

	al.logger.Info("Listener stopped for node", "node_id", nodeID, "port", port)
	return nil
}
