package node

import (
	"crypto/sha256"
	"fmt"
	"node-component/src/internal/types"
	"time"
)

// NodeConfig represents the complete configuration for a node
type NodeConfig struct {
	NodeID           string    `json:"node_id"`
	GridToken        string    `json:"grid_token"`
	SSHPublicKey     string    `json:"ssh_public_key"`
	SSHPrivateKey    string    `json:"ssh_private_key"`
	NodeCertificate  string    `json:"node_certificate"`
	CommandStationIP string    `json:"command_station_ip"`
	CreatedAt        time.Time `json:"created_at"`
	ExpiresAt        time.Time `json:"expires_at"`
}

// USBDevice represents a USB storage device
type USBDevice struct {
	Path        string `json:"path"`
	RawPath     string `json:"raw_path"`
	Capacity    int64  `json:"capacity"`
	Speed       int    `json:"speed"`
	Vendor      string `json:"vendor"`
	Model       string `json:"model"`
	Serial      string `json:"serial"`
	IsSystem    bool   `json:"is_system"`
	IsRemovable bool   `json:"is_removable"`
}

// SelectedUSBDevice represents a USB device that has been validated and selected by the user
type SelectedUSBDevice struct {
	Device          types.USBDevice `json:"device"`
	SelectedAt      time.Time       `json:"selected_at"`
	ValidationToken string          `json:"validation_token"` // Unique token to ensure it's the same device
	Platform        string          `json:"platform"`
}

// NewSelectedUSBDevice creates a selected device after validations
func NewSelectedUSBDevice(device types.USBDevice, platform string) *SelectedUSBDevice {
	// Generate unique token based on device characteristics
	token := generateDeviceToken(device)
	return &SelectedUSBDevice{
		Device:          device,
		SelectedAt:      time.Now(),
		ValidationToken: token,
		Platform:        platform,
	}
}

// Validate verifies if the device is still valid and is the same
func (s *SelectedUSBDevice) Validate(currentDevice types.USBDevice) error {
	currentToken := generateDeviceToken(currentDevice)
	if currentToken != s.ValidationToken {
		return fmt.Errorf("device mismatch: expected token %s, got %s",
			s.ValidationToken, currentToken)
	}
	return nil
}

// generateDeviceToken generates unique token based on device characteristics
func generateDeviceToken(device types.USBDevice) string {
	data := fmt.Sprintf("%s-%s-%s-%d",
		device.Path, device.Model, device.Serial, device.Capacity)
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash[:8])
}

// NodeStatus represents the current status of a node
type NodeStatus struct {
	NodeID        string        `json:"node_id"`
	Status        string        `json:"status"`
	IPAddress     string        `json:"ip_address"`
	Uptime        time.Duration `json:"uptime"`
	LastHeartbeat time.Time     `json:"last_heartbeat"`
	Hardware      *HardwareInfo `json:"hardware"`
	CreatedAt     time.Time     `json:"created_at"`
	RegisteredAt  time.Time     `json:"registered_at"`
}

// HardwareInfo represents hardware specifications of a node
type HardwareInfo struct {
	CPUCores      int    `json:"cpu_cores"`
	MemoryGB      int    `json:"memory_gb"`
	DiskGB        int    `json:"disk_gb"`
	Hostname      string `json:"hostname"`
	IPAddress     string `json:"ip_address"`
	OSVersion     string `json:"os_version"`
	KernelVersion string `json:"kernel_version"`
}

// CloudInitConfig represents cloud-init configuration files
type CloudInitConfig struct {
	UserData      string `yaml:"user_data"`
	NetworkConfig string `yaml:"network_config"`
	MetaData      string `yaml:"meta_data"`
	Valid         bool   `yaml:"valid"`
}

// NodeAnnouncement represents the message sent by a node during handshake
type NodeAnnouncement struct {
	Type            string        `json:"type"`
	NodeID          string        `json:"node_id"`
	GridToken       string        `json:"grid_token"`
	NodeCertificate string        `json:"node_certificate"`
	Hardware        *HardwareInfo `json:"hardware"`
	Timestamp       time.Time     `json:"timestamp"`
}

// HandshakeResponse represents the response from command station to node
type HandshakeResponse struct {
	Status             string                 `json:"status"`
	Message            string                 `json:"message"`
	CommandStationCert string                 `json:"command_station_cert"`
	SSHConfig          map[string]interface{} `json:"ssh_config"`
	WorkloadConfig     map[string]interface{} `json:"workload_config"`
	Timestamp          time.Time              `json:"timestamp"`
}

// Event represents an event in the node lifecycle
type Event struct {
	Type      string                 `json:"type"`
	NodeID    string                 `json:"node_id"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}

// NodeError represents a structured error with context
type NodeError struct {
	Code      string    `json:"code"`
	Message   string    `json:"message"`
	Details   string    `json:"details"`
	Timestamp time.Time `json:"timestamp"`
}

// CreateOptions and CreateResult are now defined in create.go

// NodeList represents a list of nodes with their status
type NodeList struct {
	Active   []*NodeStatus `json:"active"`
	Pending  []*NodeStatus `json:"pending"`
	Inactive []*NodeStatus `json:"inactive"`
	Total    int           `json:"total"`
}

// LogOptions represents options for retrieving logs
type LogOptions struct {
	Lines   int    `json:"lines"`
	Follow  bool   `json:"follow"`
	Service string `json:"service"`
}

// NodeLogs represents logs for a specific node
type NodeLogs struct {
	NodeID    string    `json:"node_id"`
	Logs      []string  `json:"logs"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
}

// SSHKeys represents a pair of SSH keys
type SSHKeys struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

// DeviceInfo represents detailed information about a device
type DeviceInfo struct {
	Path     string `json:"path"`
	Capacity int64  `json:"capacity"`
	Speed    int    `json:"speed"`
	Vendor   string `json:"vendor"`
	Model    string `json:"model"`
	Serial   string `json:"serial"`
}

// WriteProgress represents the progress of a USB write operation
type WriteProgress struct {
	Device       USBDevice     `json:"device"`
	BytesWritten int64         `json:"bytes_written"`
	TotalBytes   int64         `json:"total_bytes"`
	Percentage   float64       `json:"percentage"`
	Speed        float64       `json:"speed"` // MB/s
	ETA          time.Duration `json:"eta"`
}

// Node states
const (
	StatePending  = "pending"
	StateActive   = "active"
	StateInactive = "inactive"
	StateFailed   = "failed"
)

// Event types
const (
	EventNodeCreated      = "node_created"
	EventNodeRegistered   = "node_registered"
	EventNodeDisconnected = "node_disconnected"
	EventHeartbeatFailed  = "heartbeat_failed"
	EventNodeFailed       = "node_failed"
)

// Default values
const (
	DefaultRegistrationPort     = 51000
	DefaultHeartbeatInterval    = 30 * time.Second
	DefaultHeartbeatTimeout     = 10 * time.Second
	DefaultMaxHeartbeatFailures = 3
	DefaultRegistrationTimeout  = 30 * time.Minute
	DefaultMinUSBCapacity       = 8 * 1024 * 1024 * 1024 // 8GB
	DefaultSSHKeySize           = 2048
	DefaultTokenExpiry          = 24 * time.Hour
)

// Error codes
const (
	ErrCodeNoUSBFound         = "no_usb_found"
	ErrCodeUSBTooSmall        = "usb_too_small"
	ErrCodeInvalidISO         = "invalid_iso"
	ErrCodeTokenValidation    = "token_validation_failed"
	ErrCodeNodeNotFound       = "node_not_found"
	ErrCodeInvalidNodeID      = "invalid_node_id"
	ErrCodeHandshakeTimeout   = "handshake_timeout"
	ErrCodeHeartbeatFailed    = "heartbeat_failed"
	ErrCodeUSBWriteFailed     = "usb_write_failed"
	ErrCodeCloudInitInvalid   = "cloud_init_invalid"
	ErrCodeRegistrationFailed = "registration_failed"
	ErrCodeListenerFailed     = "listener_failed"
	ErrCodePermissionDenied   = "permission_denied"
	ErrCodeNetworkError       = "network_error"
	ErrCodeSystemDevice       = "system_device"
)
