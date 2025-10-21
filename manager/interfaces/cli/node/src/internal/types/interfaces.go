package types

import (
	"context"
	"net"
	"time"
)

// Local types for interfaces
type Event struct {
	Type      string                 `json:"type"`
	NodeID    string                 `json:"node_id"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}

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

type CloudInitConfig struct {
	NodeID            string    `json:"node_id"`
	UserData          string    `yaml:"user_data"`
	NetworkConfig     string    `yaml:"network_config"`
	MetaData          string    `yaml:"meta_data"`
	UserDataFile      string    `json:"user_data_file"`
	NetworkConfigFile string    `json:"network_config_file"`
	MetaDataFile      string    `json:"meta_data_file"`
	Valid             bool      `yaml:"valid"`
	CreatedAt         time.Time `json:"created_at"`
}

// EncryptedTokenData encrypted token for cloud-init
type EncryptedTokenData struct {
	Ciphertext  string    `json:"ciphertext"`
	NodeID      string    `json:"node_id"`
	EncryptedAt time.Time `json:"encrypted_at"`
}

type NodeAnnouncement struct {
	Type            string        `json:"type"`
	NodeID          string        `json:"node_id"`
	GridToken       string        `json:"grid_token"`
	NodeCertificate string        `json:"node_certificate"`
	Hardware        *HardwareInfo `json:"hardware"`
	Timestamp       time.Time     `json:"timestamp"`
}

type HardwareInfo struct {
	CPUCores      int    `json:"cpu_cores"`
	MemoryGB      int    `json:"memory_gb"`
	DiskGB        int    `json:"disk_gb"`
	Hostname      string `json:"hostname"`
	IPAddress     string `json:"ip_address"`
	OSVersion     string `json:"os_version"`
	KernelVersion string `json:"kernel_version"`
}

type HandshakeResponse struct {
	Status             string                 `json:"status"`
	Message            string                 `json:"message"`
	CommandStationCert string                 `json:"command_station_cert"`
	SSHConfig          map[string]interface{} `json:"ssh_config"`
	WorkloadConfig     map[string]interface{} `json:"workload_config"`
	Timestamp          time.Time              `json:"timestamp"`
}

type DeviceInfo struct {
	Path     string `json:"path"`
	Capacity int64  `json:"capacity"`
	Speed    int    `json:"speed"`
	Vendor   string `json:"vendor"`
	Model    string `json:"model"`
	Serial   string `json:"serial"`
}

// CloudInitOptions represents options for cloud-init generation
type CloudInitOptions struct {
	NetworkInterface string   `json:"network_interface"`
	DHCPEnabled      bool     `json:"dhcp_enabled"`
	StaticIP         string   `json:"static_ip"`
	Gateway          string   `json:"gateway"`
	DNSServers       []string `json:"dns_servers"`
	Timezone         string   `json:"timezone"`
	Locale           string   `json:"locale"`
}

// CloudInitStats represents statistics about cloud-init configurations
type CloudInitStats struct {
	TotalConfigs int       `json:"total_configs"`
	NodeIDs      []string  `json:"node_ids"`
	OutputDir    string    `json:"output_dir"`
	TemplatesDir string    `json:"templates_dir"`
	TotalSize    int64     `json:"total_size"`
	LastUpdated  time.Time `json:"last_updated"`
}

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

// EventHandler defines the interface for handling events
type EventHandler func(event Event)

// EventBus defines the interface for event management
type EventBus interface {
	Subscribe(eventType string, handler EventHandler)
	Unsubscribe(eventType string, handler EventHandler)
	Publish(event Event)
	Close()
}

// Logger defines the interface for structured logging
type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
	SetLevel(level string)
	WithFields(fields map[string]interface{}) Logger
}

// USBDetector defines the interface for USB device detection
type USBDetector interface {
	DetectAvailable() ([]USBDevice, error)
	ValidateDevice(device USBDevice) error
	GetDeviceInfo(device USBDevice) (*DeviceInfo, error)
	IsSystemDevice(device USBDevice) bool
	SelectBestDevice(devices []USBDevice) (*USBDevice, error)
}

// USBWriter defines the interface for USB writing operations
type USBWriter interface {
	WriteToUSB(device USBDevice, isoPath string, cloudInit *CloudInitConfig) error
	ValidateWriteOperation(device USBDevice, isoPath string) error
	GetWriteProgress(device USBDevice) (*WriteProgress, error)
	CancelWriteOperation(device USBDevice) error
	VerifyWriteOperation(device USBDevice, isoPath string) error
}

// ISODownloader defines the interface for ISO download operations
type ISODownloader interface {
	DownloadISO(url string, cacheDir string) (string, error)
	GetCachedISO(cacheDir string) (string, error)
	VerifyISO(isoPath string, expectedSHA256 string) error
	GetDownloadProgress(url string) (*DownloadProgress, error)
	CancelDownload(url string) error
}

// CloudInitGenerator defines the interface for cloud-init generation
type CloudInitGenerator interface {
	GenerateCloudInit(config *NodeConfig) (*CloudInitConfig, error)
	ValidateCloudInit(cloudInit *CloudInitConfig) error
	LoadTemplate(templateName string) (string, error)
	ExecuteTemplate(template string, variables map[string]interface{}) (string, error)
	GenerateNetworkConfig() (string, error)
	GenerateMetaData(nodeID string) (string, error)
}

// ConfigGenerator defines the interface for configuration generation
type ConfigGenerator interface {
	GenerateNodeConfig() (*NodeConfig, error)
	GenerateNodeID() (string, error)
	GenerateSSHKeys() (interface{}, error)
	GenerateNodeCertificate(nodeID string) (string, error)
	DetectCommandStationIP() (string, error)
	ValidateConfig(config *NodeConfig) error
}

// TokenIntegration defines the interface for token management integration
type TokenIntegration interface {
	GetGridToken() (string, error)
	ValidateToken(token string) error
	RefreshToken() (string, error)
	GetTokenExpiry() (time.Time, error)
	IsTokenValid(token string) bool
}

// HandshakeManager defines the interface for handshake operations
type HandshakeManager interface {
	ProcessNodeAnnouncement(announcement *NodeAnnouncement) (*HandshakeResponse, error)
	ValidateGridToken(token string) error
	ValidateNodeCertificate(cert string) error
	GenerateHandshakeResponse(nodeID string) (*HandshakeResponse, error)
	LogHandshakeEvent(nodeID string, event string, details map[string]interface{})
}

// Listener defines the interface for TCP listener operations
type Listener interface {
	Start(port int) error
	Stop() error
	IsRunning() bool
	GetActiveConnections() []net.Conn
	AcceptConnection() (net.Conn, error)
	CloseConnection(conn net.Conn) error
	SetTimeout(timeout time.Duration)
}

// HeartbeatManager defines the interface for heartbeat operations
type HeartbeatManager interface {
	StartHeartbeat(nodeID string, conn net.Conn) error
	StopHeartbeat(nodeID string) error
	IsHeartbeatActive(nodeID string) bool
	GetHeartbeatStatus(nodeID string) (*HeartbeatStatus, error)
	ProcessHeartbeat(nodeID string, data []byte) error
	GetFailedNodes() []string
	ResetHeartbeatFailures(nodeID string) error
}

// NodeStateManager defines the interface for node state management
type NodeStateManager interface {
	CreateNode(nodeID string, config *NodeConfig) error
	GetNode(nodeID string) (*NodeStatus, error)
	UpdateNodeStatus(nodeID string, status string) error
	TransitionToActive(nodeID string, ipAddress string) error
	TransitionToInactive(nodeID string) error
	GetActiveNodes() []*NodeStatus
	GetPendingNodes() []*NodeStatus
	GetInactiveNodes() []*NodeStatus
	RemoveNode(nodeID string) error
	IsNodePending(nodeID string) bool
	IsNodeActive(nodeID string) bool
	SaveState() error
	LoadState() error
}

// CertificateManager defines the interface for certificate operations
type CertificateManager interface {
	GenerateCertificate(nodeID string) (string, error)
	ValidateCertificate(cert []byte) error
	GetCommandStationCert() string
	SignCertificate(cert []byte) ([]byte, error)
	VerifySignature(cert []byte, signature []byte) error
}

// NetworkManager defines the interface for network operations
type NetworkManager interface {
	GetLocalIP() (string, error)
	IsPortAvailable(port int) bool
	TestConnectivity(host string, port int) error
	GetNetworkInterfaces() ([]NetworkInterface, error)
	ConfigureFirewall(port int) error
}

// FileSystemManager defines the interface for filesystem operations
type FileSystemManager interface {
	CreateDirectory(path string) error
	WriteFile(path string, data []byte) error
	ReadFile(path string) ([]byte, error)
	FileExists(path string) bool
	DeleteFile(path string) error
	GetFileSize(path string) (int64, error)
	SetPermissions(path string, mode int) error
}

// ProcessManager defines the interface for process operations
type ProcessManager interface {
	ExecuteCommand(command string, args ...string) ([]byte, error)
	ExecuteCommandWithContext(ctx context.Context, command string, args ...string) ([]byte, error)
	IsProcessRunning(pid int) bool
	KillProcess(pid int) error
	GetProcessInfo(pid int) (*ProcessInfo, error)
}

// SecurityManager defines the interface for security operations
type SecurityManager interface {
	GenerateSecureKey(size int) ([]byte, error)
	EncryptData(data []byte, key []byte) ([]byte, error)
	DecryptData(data []byte, key []byte) ([]byte, error)
	HashData(data []byte) string
	ValidateHash(data []byte, hash string) bool
	GenerateRandomBytes(size int) ([]byte, error)
}

// MetricsCollector defines the interface for metrics collection
type MetricsCollector interface {
	CollectMetrics() (*Metrics, error)
	GetNodeMetrics(nodeID string) (*NodeMetrics, error)
	GetSystemMetrics() (*SystemMetrics, error)
	UpdateCounter(name string, value int64)
	UpdateGauge(name string, value float64)
	RecordHistogram(name string, value float64)
}

// ConfigurationManager defines the interface for configuration management
type ConfigurationManager interface {
	LoadConfig(configPath string) (*Configuration, error)
	SaveConfig(config *Configuration, configPath string) error
	GetDefaultConfig() *Configuration
	ValidateConfig(config *Configuration) error
	MergeConfigs(base *Configuration, override *Configuration) *Configuration
}

// Supporting types for interfaces

// DownloadProgress represents the progress of a download operation
type DownloadProgress struct {
	URL             string        `json:"url"`
	BytesDownloaded int64         `json:"bytes_downloaded"`
	TotalBytes      int64         `json:"total_bytes"`
	Percentage      float64       `json:"percentage"`
	Speed           float64       `json:"speed"` // MB/s
	ETA             time.Duration `json:"eta"`
}

// HeartbeatStatus represents the status of a heartbeat connection
type HeartbeatStatus struct {
	NodeID           string        `json:"node_id"`
	IsActive         bool          `json:"is_active"`
	LastHeartbeat    time.Time     `json:"last_heartbeat"`
	FailureCount     int           `json:"failure_count"`
	Latency          time.Duration `json:"latency"`
	ConnectionStatus string        `json:"connection_status"`
}

// NetworkInterface represents a network interface
type NetworkInterface struct {
	Name    string   `json:"name"`
	IP      string   `json:"ip"`
	Netmask string   `json:"netmask"`
	Gateway string   `json:"gateway"`
	DNS     []string `json:"dns"`
}

// ProcessInfo represents information about a process
type ProcessInfo struct {
	PID       int       `json:"pid"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CPU       float64   `json:"cpu"`
	Memory    int64     `json:"memory"`
	StartTime time.Time `json:"start_time"`
}

// Metrics represents collected metrics
type Metrics struct {
	Timestamp  time.Time             `json:"timestamp"`
	Counters   map[string]int64      `json:"counters"`
	Gauges     map[string]float64    `json:"gauges"`
	Histograms map[string]*Histogram `json:"histograms"`
}

// NodeMetrics represents metrics for a specific node
type NodeMetrics struct {
	NodeID      string        `json:"node_id"`
	CPUUsage    float64       `json:"cpu_usage"`
	MemoryUsage float64       `json:"memory_usage"`
	DiskUsage   float64       `json:"disk_usage"`
	NetworkIn   int64         `json:"network_in"`
	NetworkOut  int64         `json:"network_out"`
	Uptime      time.Duration `json:"uptime"`
	Timestamp   time.Time     `json:"timestamp"`
}

// SystemMetrics represents system-wide metrics
type SystemMetrics struct {
	TotalCPU    float64       `json:"total_cpu"`
	TotalMemory int64         `json:"total_memory"`
	FreeMemory  int64         `json:"free_memory"`
	TotalDisk   int64         `json:"total_disk"`
	FreeDisk    int64         `json:"free_disk"`
	LoadAverage []float64     `json:"load_average"`
	Uptime      time.Duration `json:"uptime"`
	Timestamp   time.Time     `json:"timestamp"`
}

// Histogram represents a histogram metric
type Histogram struct {
	Buckets map[float64]int64 `json:"buckets"`
	Count   int64             `json:"count"`
	Sum     float64           `json:"sum"`
}

// Configuration represents the component configuration
type Configuration struct {
	USBDetection struct {
		MinCapacity    int64         `yaml:"min_capacity"`
		MaxCapacity    int64         `yaml:"max_capacity"`
		PreferredSpeed int           `yaml:"preferred_speed"`
		ScanInterval   time.Duration `yaml:"scan_interval"`
	} `yaml:"usb_detection"`

	ISODownload struct {
		URL            string        `yaml:"url"`
		CacheDir       string        `yaml:"cache_dir"`
		SHA256Checksum string        `yaml:"sha256_checksum"`
		Timeout        time.Duration `yaml:"timeout"`
	} `yaml:"iso_download"`

	Registration struct {
		Port          int           `yaml:"port"`
		Timeout       time.Duration `yaml:"timeout"`
		MaxConcurrent int           `yaml:"max_concurrent"`
		RetryAttempts int           `yaml:"retry_attempts"`
	} `yaml:"registration"`

	Heartbeat struct {
		Interval      time.Duration `yaml:"interval"`
		Timeout       time.Duration `yaml:"timeout"`
		MaxFailures   int           `yaml:"max_failures"`
		RetryInterval time.Duration `yaml:"retry_interval"`
	} `yaml:"heartbeat"`

	Security struct {
		TokenExpiry   time.Duration `yaml:"token_expiry"`
		CertExpiry    time.Duration `yaml:"cert_expiry"`
		KeySize       int           `yaml:"key_size"`
		HashAlgorithm string        `yaml:"hash_algorithm"`
	} `yaml:"security"`

	Logging struct {
		Level      string `yaml:"level"`
		Format     string `yaml:"format"`
		Output     string `yaml:"output"`
		MaxSize    int    `yaml:"max_size"`
		MaxBackups int    `yaml:"max_backups"`
		MaxAge     int    `yaml:"max_age"`
	} `yaml:"logging"`
}

// UbuntuVersion represents information about an Ubuntu version
type UbuntuVersion struct {
	Version     string    `json:"version"`
	LTS         bool      `json:"lts"`
	FileName    string    `json:"file_name"`
	DownloadURL string    `json:"download_url"`
	SHA256      string    `json:"sha256"`
	Size        int64     `json:"size"`
	ReleaseDate time.Time `json:"release_date"`
}

// ISOInfo represents information about a downloaded ISO
type ISOInfo struct {
	Version      string    `json:"version"`
	FilePath     string    `json:"file_path"`
	FileName     string    `json:"file_name"`
	Size         int64     `json:"size"`
	SHA256       string    `json:"sha256"`
	DownloadURL  string    `json:"download_url"`
	DownloadedAt time.Time `json:"downloaded_at"`
}

// ISOCacheStats represents statistics about ISO cache
type ISOCacheStats struct {
	CacheDir   string    `json:"cache_dir"`
	TotalFiles int       `json:"total_files"`
	TotalSize  int64     `json:"total_size"`
	OldestFile time.Time `json:"oldest_file"`
	NewestFile time.Time `json:"newest_file"`
}

// WriteResult represents the result of a write operation
type WriteResult struct {
	DevicePath        string        `json:"device_path"`
	ISOPath           string        `json:"iso_path"`
	BytesWritten      int64         `json:"bytes_written"`
	Duration          time.Duration `json:"duration"`
	Success           bool          `json:"success"`
	ErrorMessage      string        `json:"error_message,omitempty"`
	CloudInitInjected bool          `json:"cloud_init_injected"`
}

// WriteProgress represents the progress of a write operation
type WriteProgress struct {
	BytesWritten int64         `json:"bytes_written"`
	TotalBytes   int64         `json:"total_bytes"`
	Percentage   float64       `json:"percentage"`
	Speed        float64       `json:"speed"` // MB/s
	ETA          time.Duration `json:"eta"`
	ElapsedTime  time.Duration `json:"elapsed_time"`
}

// SSHConfig represents SSH configuration for the node
type SSHConfig struct {
	PublicKey     string `json:"public_key"`
	PrivateKey    string `json:"private_key"`
	AuthorizedKey string `json:"authorized_key"`
	Port          int    `json:"port"`
	Host          string `json:"host"`
}

// ResourceLimits represents resource limits for the node
type ResourceLimits struct {
	CPULimit    string `json:"cpu_limit"`
	MemoryLimit string `json:"memory_limit"`
	DiskLimit   string `json:"disk_limit"`
}

// DockerConfig represents Docker configuration
type DockerConfig struct {
	Enabled     bool   `json:"enabled"`
	Version     string `json:"version"`
	RegistryURL string `json:"registry_url"`
}

// NetworkConfig represents network configuration
type NetworkConfig struct {
	AllowedPorts []int  `json:"allowed_ports"`
	FirewallRule string `json:"firewall_rule"`
}

// WorkloadConfig represents workload configuration for the node
type WorkloadConfig struct {
	MaxWorkloads   int             `json:"max_workloads"`
	ResourceLimits *ResourceLimits `json:"resource_limits"`
	DockerConfig   *DockerConfig   `json:"docker_config"`
	NetworkConfig  *NetworkConfig  `json:"network_config"`
}

// SSHKeys represents SSH key pair for a node
type SSHKeys struct {
	PublicKey     string `json:"public_key"`
	PrivateKey    string `json:"private_key"`
	AuthorizedKey string `json:"authorized_key"`
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
