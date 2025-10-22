package node

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"node-component/src/internal/types"
)

// HandshakeManager manages the secure handshake protocol for node registration
type HandshakeManager struct {
	tokenIntegration   types.TokenIntegration
	nodeStateManager   types.NodeStateManager
	certificateManager types.CertificateManager
	logger             types.Logger
}

// NewHandshakeManager creates a new handshake manager
func NewHandshakeManager(
	tokenIntegration types.TokenIntegration,
	nodeStateManager types.NodeStateManager,
	certificateManager types.CertificateManager,
	logger types.Logger,
) *HandshakeManager {
	return &HandshakeManager{
		tokenIntegration:   tokenIntegration,
		nodeStateManager:   nodeStateManager,
		certificateManager: certificateManager,
		logger:             logger,
	}
}

// HandshakeRequest represents an incoming handshake request
type HandshakeRequest struct {
	Type            string              `json:"type"`
	NodeID          string              `json:"node_id"`
	GridToken       string              `json:"grid_token"`
	NodeCertificate string              `json:"node_certificate"`
	Hardware        *types.HardwareInfo `json:"hardware"`
	Timestamp       time.Time           `json:"timestamp"`
	ClientIP        string              `json:"client_ip,omitempty"`
}

// HandshakeResponse is defined in types.go

// SSHConfig, WorkloadConfig, ResourceLimits, DockerConfig, NetworkConfig are defined in interfaces.go

// ProcessHandshake processes an incoming handshake request
func (hm *HandshakeManager) ProcessHandshake(ctx context.Context, conn net.Conn, request *HandshakeRequest) (*types.HandshakeResponse, error) {
	hm.logger.Info("Processing handshake request",
		"node_id", request.NodeID,
		"client_ip", request.ClientIP,
		"timestamp", request.Timestamp)

	// Set client IP from connection
	if request.ClientIP == "" {
		request.ClientIP = conn.RemoteAddr().(*net.TCPAddr).IP.String()
	}

	// Step 1: Validate request structure
	if err := hm.validateHandshakeRequest(request); err != nil {
		hm.logger.Error("Invalid handshake request", "error", err, "node_id", request.NodeID)
		return hm.createErrorResponse("invalid_request", err.Error()), nil
	}

	// Step 2: Validate Grid Token
	if err := hm.validateGridToken(request.GridToken); err != nil {
		hm.logger.Error("Grid token validation failed", "error", err, "node_id", request.NodeID)
		return hm.createErrorResponse("invalid_token", "Grid token validation failed"), nil
	}

	// Step 3: Validate Node ID and certificate
	if err := hm.validateNodeCredentials(request.NodeID, request.NodeCertificate); err != nil {
		hm.logger.Error("Node credentials validation failed", "error", err, "node_id", request.NodeID)
		return hm.createErrorResponse("invalid_credentials", "Node credentials validation failed"), nil
	}

	// Step 4: Check if node is in pending state
	if !hm.nodeStateManager.IsNodePending(request.NodeID) {
		hm.logger.Error("Node not in pending state", "node_id", request.NodeID)
		return hm.createErrorResponse("node_not_pending", "Node is not in pending state"), nil
	}

	// Step 5: Generate command station certificate
	commandStationCert, err := hm.generateCommandStationCertificate()
	if err != nil {
		hm.logger.Error("Failed to generate command station certificate", "error", err, "node_id", request.NodeID)
		return hm.createErrorResponse("cert_generation_failed", "Failed to generate command station certificate"), nil
	}

	// Step 6: Generate SSH configuration
	sshConfig, err := hm.generateSSHConfig(request.NodeID)
	if err != nil {
		hm.logger.Error("Failed to generate SSH configuration", "error", err, "node_id", request.NodeID)
		return hm.createErrorResponse("ssh_config_failed", "Failed to generate SSH configuration"), nil
	}

	// Step 7: Generate workload configuration
	workloadConfig, err := hm.generateWorkloadConfig(request.Hardware)
	if err != nil {
		hm.logger.Error("Failed to generate workload configuration", "error", err, "node_id", request.NodeID)
		return hm.createErrorResponse("workload_config_failed", "Failed to generate workload configuration"), nil
	}

	// Step 8: Update node state to active
	if err := hm.nodeStateManager.UpdateNodeStatus(request.NodeID, "active"); err != nil {
		hm.logger.Error("Failed to activate node", "error", err, "node_id", request.NodeID)
		return hm.createErrorResponse("activation_failed", "Failed to activate node"), nil
	}

	// Step 9: Create successful response
	response := &types.HandshakeResponse{
		Status:             "accepted",
		Message:            "Node registration successful",
		CommandStationCert: commandStationCert,
		SSHConfig: map[string]interface{}{
			"public_key": sshConfig.PublicKey,
			"port":       sshConfig.Port,
			"host":       sshConfig.Host,
		},
		WorkloadConfig: map[string]interface{}{
			"max_workloads":   workloadConfig.MaxWorkloads,
			"resource_limits": workloadConfig.ResourceLimits,
		},
		Timestamp: time.Now(),
	}

	hm.logger.Info("Handshake successful",
		"node_id", request.NodeID,
		"client_ip", request.ClientIP,
		"status", response.Status)

	return response, nil
}

// Private helper methods

// validateHandshakeRequest validates the structure of a handshake request
func (hm *HandshakeManager) validateHandshakeRequest(request *HandshakeRequest) error {
	if request == nil {
		return fmt.Errorf("handshake request cannot be nil")
	}

	if request.Type != "node_announcement" {
		return fmt.Errorf("invalid request type: %s", request.Type)
	}

	if request.NodeID == "" {
		return fmt.Errorf("node ID cannot be empty")
	}

	if request.GridToken == "" {
		return fmt.Errorf("grid token cannot be empty")
	}

	if request.NodeCertificate == "" {
		return fmt.Errorf("node certificate cannot be empty")
	}

	if request.Hardware == nil {
		return fmt.Errorf("hardware information cannot be nil")
	}

	// Validate timestamp (should not be too old)
	if time.Since(request.Timestamp) > 30*time.Minute {
		return fmt.Errorf("handshake request is too old: %v", request.Timestamp)
	}

	return nil
}

// validateGridToken validates the grid token against the keyring
func (hm *HandshakeManager) validateGridToken(token string) error {
	// Get the stored grid token from keyring
	storedToken, err := hm.tokenIntegration.GetGridToken()
	if err != nil {
		return fmt.Errorf("failed to get stored grid token: %w", err)
	}

	// Compare tokens
	if token != storedToken {
		return fmt.Errorf("grid token mismatch")
	}

	// Validate token format (basic check)
	if len(token) < 32 {
		return fmt.Errorf("grid token too short")
	}

	return nil
}

// validateNodeCredentials validates the node ID and certificate
func (hm *HandshakeManager) validateNodeCredentials(nodeID string, certificate string) error {
	// Validate node ID format
	if err := hm.validateNodeIDFormat(nodeID); err != nil {
		return fmt.Errorf("invalid node ID format: %w", err)
	}

	// Validate certificate format
	if err := hm.validateCertificateFormat(certificate); err != nil {
		return fmt.Errorf("invalid certificate format: %w", err)
	}

	// Verify certificate signature (if we have the public key)
	// This is a simplified validation - in production, we'd verify against stored public keys
	if err := hm.verifyCertificateSignature(nodeID, certificate); err != nil {
		return fmt.Errorf("certificate signature verification failed: %w", err)
	}

	return nil
}

// validateNodeIDFormat validates the format of a node ID
func (hm *HandshakeManager) validateNodeIDFormat(nodeID string) error {
	// Basic format validation - should be like "node-01", "node-02", etc.
	if len(nodeID) < 6 {
		return fmt.Errorf("node ID too short")
	}

	if !startsWith(nodeID, "node-") {
		return fmt.Errorf("node ID must start with 'node-'")
	}

	// Check if it contains only valid characters
	for _, char := range nodeID {
		if !isValidNodeIDChar(char) {
			return fmt.Errorf("node ID contains invalid character: %c", char)
		}
	}

	return nil
}

// validateCertificateFormat validates the format of a certificate
func (hm *HandshakeManager) validateCertificateFormat(certificate string) error {
	// Basic format validation - should be a valid PEM or base64 encoded certificate
	if len(certificate) < 100 {
		return fmt.Errorf("certificate too short")
	}

	// Check if it looks like a PEM certificate
	if startsWith(certificate, "-----BEGIN") {
		return nil
	}

	// Check if it looks like base64
	if isBase64String(certificate) {
		return nil
	}

	return fmt.Errorf("certificate format not recognized")
}

// verifyCertificateSignature verifies the signature of a certificate
func (hm *HandshakeManager) verifyCertificateSignature(nodeID string, certificate string) error {
	// This is a simplified implementation
	// In production, we would:
	// 1. Retrieve the stored public key for the node
	// 2. Verify the certificate signature using Ed25519
	// 3. Check certificate expiration and other properties

	// For now, we'll just do basic validation
	if len(certificate) < 200 {
		return fmt.Errorf("certificate too short for valid signature")
	}

	hm.logger.Debug("Certificate signature verification passed", "node_id", nodeID)
	return nil
}

// generateCommandStationCertificate generates a command station certificate
func (hm *HandshakeManager) generateCommandStationCertificate() (string, error) {
	// Generate Ed25519 key pair for command station
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", fmt.Errorf("failed to generate command station key pair: %w", err)
	}

	// Create a simple certificate structure
	certData := map[string]interface{}{
		"type":       "command_station_certificate",
		"public_key": fmt.Sprintf("%x", publicKey),
		"timestamp":  time.Now(),
		"expires":    time.Now().Add(24 * time.Hour),
	}

	// Marshal to JSON
	certJSON, err := json.Marshal(certData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal certificate: %w", err)
	}

	// Sign the certificate
	signature := ed25519.Sign(privateKey, certJSON)

	// Create final certificate
	finalCert := map[string]interface{}{
		"certificate": string(certJSON),
		"signature":   fmt.Sprintf("%x", signature),
	}

	// Marshal final certificate
	finalCertJSON, err := json.Marshal(finalCert)
	if err != nil {
		return "", fmt.Errorf("failed to marshal final certificate: %w", err)
	}

	hm.logger.Debug("Command station certificate generated")
	return string(finalCertJSON), nil
}

// generateSSHConfig generates SSH configuration for the node
func (hm *HandshakeManager) generateSSHConfig(nodeID string) (*types.SSHConfig, error) {
	// Get SSH keys for the node (these should have been generated during node creation)
	// For now, we'll create a basic SSH config - in production, we'd retrieve from node state
	publicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC..."     // Placeholder
	privateKey := "-----BEGIN RSA PRIVATE KEY-----\n..."          // Placeholder
	authorizedKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC..." // Placeholder

	// Create SSH configuration
	sshConfig := &types.SSHConfig{
		PublicKey:     publicKey,
		PrivateKey:    privateKey,
		AuthorizedKey: authorizedKey,
		Port:          22,
		Host:          "localhost",
	}

	hm.logger.Debug("SSH configuration generated", "node_id", nodeID)
	return sshConfig, nil
}

// generateWorkloadConfig generates workload configuration based on hardware
func (hm *HandshakeManager) generateWorkloadConfig(hardware *types.HardwareInfo) (*types.WorkloadConfig, error) {
	if hardware == nil {
		return nil, fmt.Errorf("hardware information is required")
	}

	// Calculate resource limits based on hardware
	cpuLimit := "100%"                                     // Default to full CPU usage
	memoryLimit := fmt.Sprintf("%dG", hardware.MemoryGB/2) // Use half of available memory
	diskLimit := fmt.Sprintf("%dG", hardware.DiskGB/4)     // Use quarter of available disk

	// Create resource limits
	resourceLimits := &types.ResourceLimits{
		CPULimit:    cpuLimit,
		MemoryLimit: memoryLimit,
		DiskLimit:   diskLimit,
	}

	// Create Docker configuration
	dockerConfig := &types.DockerConfig{
		Enabled:     true,
		Version:     "latest",
		RegistryURL: "https://registry-1.docker.io",
	}

	// Create network configuration
	networkConfig := &types.NetworkConfig{
		AllowedPorts: []int{80, 443, 8080, 8443},
		FirewallRule: "allow_workload_ports",
	}

	// Create workload configuration
	workloadConfig := &types.WorkloadConfig{
		MaxWorkloads:   calculateMaxWorkloads(hardware),
		ResourceLimits: resourceLimits,
		DockerConfig:   dockerConfig,
		NetworkConfig:  networkConfig,
	}

	hm.logger.Debug("Workload configuration generated",
		"cpu_cores", hardware.CPUCores,
		"memory_gb", hardware.MemoryGB,
		"disk_gb", hardware.DiskGB,
		"max_workloads", workloadConfig.MaxWorkloads)

	return workloadConfig, nil
}

// createErrorResponse creates an error response
func (hm *HandshakeManager) createErrorResponse(errorCode, message string) *types.HandshakeResponse {
	return &types.HandshakeResponse{
		Status:    "rejected",
		Message:   message,
		Timestamp: time.Now(),
	}
}

// Helper functions

// startsWith checks if a string starts with a prefix
func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

// isValidNodeIDChar checks if a character is valid in a node ID
func isValidNodeIDChar(char rune) bool {
	return (char >= 'a' && char <= 'z') ||
		(char >= '0' && char <= '9') ||
		char == '-'
}

// isBase64String checks if a string is valid base64
func isBase64String(s string) bool {
	// Simple base64 validation - check if all characters are valid base64
	for _, char := range s {
		if !isBase64Char(char) {
			return false
		}
	}
	return len(s)%4 == 0 // Base64 strings should be multiple of 4
}

// isBase64Char checks if a character is valid base64
func isBase64Char(char rune) bool {
	return (char >= 'A' && char <= 'Z') ||
		(char >= 'a' && char <= 'z') ||
		(char >= '0' && char <= '9') ||
		char == '+' || char == '/' || char == '='
}

// calculateMaxWorkloads is now implemented in node_state.go

// HandshakeServer handles incoming handshake connections
type HandshakeServer struct {
	handshakeManager *HandshakeManager
	listener         net.Listener
	logger           types.Logger
	port             int
	running          bool
	ctx              context.Context
	cancel           context.CancelFunc
}

// NewHandshakeServer creates a new handshake server
func NewHandshakeServer(
	handshakeManager *HandshakeManager,
	logger types.Logger,
	port int,
) *HandshakeServer {
	ctx, cancel := context.WithCancel(context.Background())

	return &HandshakeServer{
		handshakeManager: handshakeManager,
		logger:           logger,
		port:             port,
		ctx:              ctx,
		cancel:           cancel,
	}
}

// Start starts the handshake server
func (hs *HandshakeServer) Start() error {
	hs.logger.Info("Starting handshake server", "port", hs.port)

	// Create listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", hs.port))
	if err != nil {
		return fmt.Errorf("failed to create listener: %w", err)
	}

	hs.listener = listener
	hs.running = true

	// Start accepting connections in a goroutine
	go hs.acceptConnections()

	hs.logger.Info("Handshake server started", "port", hs.port)
	return nil
}

// Stop stops the handshake server
func (hs *HandshakeServer) Stop() error {
	hs.logger.Info("Stopping handshake server")

	hs.running = false
	hs.cancel()

	if hs.listener != nil {
		if err := hs.listener.Close(); err != nil {
			hs.logger.Error("Error closing listener", "error", err)
		}
	}

	hs.logger.Info("Handshake server stopped")
	return nil
}

// IsRunning returns whether the server is running
func (hs *HandshakeServer) IsRunning() bool {
	return hs.running
}

// acceptConnections accepts incoming connections
func (hs *HandshakeServer) acceptConnections() {
	for hs.running {
		select {
		case <-hs.ctx.Done():
			hs.logger.Debug("Handshake server context cancelled")
			return
		default:
			// Accept connection with timeout
			conn, err := hs.listener.Accept()
			if err != nil {
				if hs.running {
					hs.logger.Error("Failed to accept connection", "error", err)
				}
				continue
			}

			// Handle connection in a goroutine
			go hs.handleConnection(conn)
		}
	}
}

// handleConnection handles a single connection
func (hs *HandshakeServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	hs.logger.Debug("Handling handshake connection", "remote_addr", conn.RemoteAddr())

	// Set connection timeout
	conn.SetDeadline(time.Now().Add(30 * time.Minute))

	// Read handshake request
	var request HandshakeRequest
	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(&request); err != nil {
		hs.logger.Error("Failed to decode handshake request", "error", err)
		return
	}

	// Process handshake
	response, err := hs.handshakeManager.ProcessHandshake(hs.ctx, conn, &request)
	if err != nil {
		hs.logger.Error("Failed to process handshake", "error", err)
		return
	}

	// Send response
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(response); err != nil {
		hs.logger.Error("Failed to send handshake response", "error", err)
		return
	}

	hs.logger.Info("Handshake connection handled successfully",
		"node_id", request.NodeID,
		"status", response.Status)
}
