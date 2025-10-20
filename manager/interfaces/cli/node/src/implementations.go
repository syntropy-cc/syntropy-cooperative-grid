package node

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/syntropy-grid/syntropy-cooperative-grid/manager/interfaces/cli/node/src/internal/constants"
	"github.com/syntropy-grid/syntropy-cooperative-grid/manager/interfaces/cli/node/src/internal/types"
)

// CreateSubcomponent is now implemented in create.go

// RegistrationSubcomponent handles node registration and monitoring
type RegistrationSubcomponent struct {
	nodeState types.NodeStateManager
	eventBus  types.EventBus
	logger    types.Logger
	config    *types.Configuration
	mutex     sync.RWMutex
}

// EventBus implementation
type eventBus struct {
	subscribers map[string][]types.EventHandler
	mutex       sync.RWMutex
}

func (eb *eventBus) Subscribe(eventType string, handler types.EventHandler) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()
	eb.subscribers[eventType] = append(eb.subscribers[eventType], handler)
}

func (eb *eventBus) Unsubscribe(eventType string, handler types.EventHandler) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()
	handlers := eb.subscribers[eventType]
	for i, h := range handlers {
		if fmt.Sprintf("%p", h) == fmt.Sprintf("%p", handler) {
			eb.subscribers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
}

func (eb *eventBus) Publish(event types.Event) {
	eb.mutex.RLock()
	handlers := eb.subscribers[event.Type]
	eb.mutex.RUnlock()

	for _, handler := range handlers {
		go handler(event)
	}
}

func (eb *eventBus) Close() {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()
	eb.subscribers = make(map[string][]types.EventHandler)
}

// Logger implementation
type logger struct {
	level string
}

func (l *logger) Debug(msg string, fields ...interface{}) {
	if l.level == "debug" {
		fmt.Printf("[DEBUG] %s %v\n", msg, fields)
	}
}

func (l *logger) Info(msg string, fields ...interface{}) {
	fmt.Printf("[INFO] %s %v\n", msg, fields)
}

func (l *logger) Warn(msg string, fields ...interface{}) {
	fmt.Printf("[WARN] %s %v\n", msg, fields)
}

func (l *logger) Error(msg string, fields ...interface{}) {
	fmt.Printf("[ERROR] %s %v\n", msg, fields)
}

func (l *logger) Fatal(msg string, fields ...interface{}) {
	fmt.Printf("[FATAL] %s %v\n", msg, fields)
}

func (l *logger) SetLevel(level string) {
	l.level = level
}

func (l *logger) WithFields(fields map[string]interface{}) types.Logger {
	return l
}

// NodeStateManager implementation
type nodeStateManager struct {
	nodes map[string]*types.NodeStatus
	mutex sync.RWMutex
}

func (nsm *nodeStateManager) CreateNode(nodeID string, config *types.NodeConfig) error {
	nsm.mutex.Lock()
	defer nsm.mutex.Unlock()

	nsm.nodes[nodeID] = &types.NodeStatus{
		NodeID:    nodeID,
		Status:    constants.StatePending,
		CreatedAt: time.Now(),
	}

	return nil
}

func (nsm *nodeStateManager) GetNode(nodeID string) (*types.NodeStatus, error) {
	nsm.mutex.RLock()
	defer nsm.mutex.RUnlock()

	node, exists := nsm.nodes[nodeID]
	if !exists {
		return nil, fmt.Errorf("node not found: %s", nodeID)
	}

	return node, nil
}

func (nsm *nodeStateManager) UpdateNodeStatus(nodeID string, status string) error {
	nsm.mutex.Lock()
	defer nsm.mutex.Unlock()

	node, exists := nsm.nodes[nodeID]
	if !exists {
		return fmt.Errorf("node not found: %s", nodeID)
	}

	node.Status = status
	return nil
}

func (nsm *nodeStateManager) TransitionToActive(nodeID string, ipAddress string) error {
	nsm.mutex.Lock()
	defer nsm.mutex.Unlock()

	node, exists := nsm.nodes[nodeID]
	if !exists {
		return fmt.Errorf("node not found: %s", nodeID)
	}

	node.Status = constants.StateActive
	node.IPAddress = ipAddress
	node.RegisteredAt = time.Now()

	return nil
}

func (nsm *nodeStateManager) TransitionToInactive(nodeID string) error {
	nsm.mutex.Lock()
	defer nsm.mutex.Unlock()

	node, exists := nsm.nodes[nodeID]
	if !exists {
		return fmt.Errorf("node not found: %s", nodeID)
	}

	node.Status = constants.StateInactive
	return nil
}

func (nsm *nodeStateManager) GetActiveNodes() []*types.NodeStatus {
	nsm.mutex.RLock()
	defer nsm.mutex.RUnlock()

	var activeNodes []*types.NodeStatus
	for _, node := range nsm.nodes {
		if node.Status == constants.StateActive {
			activeNodes = append(activeNodes, node)
		}
	}

	return activeNodes
}

func (nsm *nodeStateManager) GetPendingNodes() []*types.NodeStatus {
	nsm.mutex.RLock()
	defer nsm.mutex.RUnlock()

	var pendingNodes []*types.NodeStatus
	for _, node := range nsm.nodes {
		if node.Status == constants.StatePending {
			pendingNodes = append(pendingNodes, node)
		}
	}

	return pendingNodes
}

func (nsm *nodeStateManager) GetInactiveNodes() []*types.NodeStatus {
	nsm.mutex.RLock()
	defer nsm.mutex.RUnlock()

	var inactiveNodes []*types.NodeStatus
	for _, node := range nsm.nodes {
		if node.Status == constants.StateInactive {
			inactiveNodes = append(inactiveNodes, node)
		}
	}

	return inactiveNodes
}

func (nsm *nodeStateManager) RemoveNode(nodeID string) error {
	nsm.mutex.Lock()
	defer nsm.mutex.Unlock()

	delete(nsm.nodes, nodeID)
	return nil
}

func (nsm *nodeStateManager) IsNodePending(nodeID string) bool {
	nsm.mutex.RLock()
	defer nsm.mutex.RUnlock()

	node, exists := nsm.nodes[nodeID]
	return exists && node.Status == constants.StatePending
}

func (nsm *nodeStateManager) IsNodeActive(nodeID string) bool {
	nsm.mutex.RLock()
	defer nsm.mutex.RUnlock()

	node, exists := nsm.nodes[nodeID]
	return exists && node.Status == constants.StateActive
}

func (nsm *nodeStateManager) SaveState() error {
	// TODO: Implement state persistence
	return nil
}

func (nsm *nodeStateManager) LoadState() error {
	// TODO: Implement state loading
	return nil
}

// TokenIntegration implementation removed - now using real integration via setup_adapter.go

// CreateSubcomponent methods are now implemented in create.go

// RegistrationSubcomponent methods
func (rs *RegistrationSubcomponent) StartListener(port int) error {
	// TODO: Implement TCP listener
	return nil
}

func (rs *RegistrationSubcomponent) StopListener() error {
	// TODO: Implement listener stop
	return nil
}

func (rs *RegistrationSubcomponent) IsListenerRunning() bool {
	// TODO: Implement listener status check
	return false
}

func (rs *RegistrationSubcomponent) StartHeartbeat(nodeID string, conn net.Conn) error {
	// TODO: Implement heartbeat start
	return nil
}

func (rs *RegistrationSubcomponent) StopHeartbeat(nodeID string) error {
	// TODO: Implement heartbeat stop
	return nil
}

func (rs *RegistrationSubcomponent) IsHeartbeatActive(nodeID string) bool {
	// TODO: Implement heartbeat status check
	return false
}

// Configuration validation - moved to a helper function
func validateConfiguration(c *types.Configuration) error {
	// TODO: Implement configuration validation
	return nil
}
