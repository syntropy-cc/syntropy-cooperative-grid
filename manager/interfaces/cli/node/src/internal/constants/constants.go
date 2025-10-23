package constants

import (
	"time"
)

// Node Component Constants

// Default Configuration Values
const (
	// USB Detection
	DefaultMinUSBCapacity    = 8 * 1024 * 1024 * 1024  // 8GB
	DefaultMaxUSBCapacity    = 64 * 1024 * 1024 * 1024 // 64GB
	DefaultPreferredUSBSpeed = 3                       // USB 3.0
	DefaultUSBScanInterval   = 5 * time.Second

	// ISO Download
	DefaultISOURL            = "https://releases.ubuntu.com/24.04/ubuntu-24.04-server-amd64.iso"
	DefaultISOCacheDir       = "~/.syntropy/cache/isos"
	DefaultISOSHA256Checksum = "a4acfda10b18da50e2ec50ccaf860d7f3b4b165bf3f4e4604dac366b93d1b25"

	// Configuration
	DefaultConfigDir          = "~/.syntropy"
	DefaultISODownloadTimeout = 60 * time.Minute

	// Registration
	DefaultRegistrationPort    = 51000
	DefaultRegistrationTimeout = 30 * time.Minute
	DefaultMaxConcurrentNodes  = 10
	DefaultRegistrationRetries = 3

	// Heartbeat
	DefaultHeartbeatInterval      = 30 * time.Second
	DefaultHeartbeatTimeout       = 10 * time.Second
	DefaultMaxHeartbeatFailures   = 3
	DefaultHeartbeatRetryInterval = 5 * time.Second

	// Security
	DefaultTokenExpiry   = 24 * time.Hour
	DefaultCertExpiry    = 365 * 24 * time.Hour
	DefaultSSHKeySize    = 2048
	DefaultHashAlgorithm = "SHA256"

	// Logging
	DefaultLogLevel      = "info"
	DefaultLogFormat     = "json"
	DefaultLogOutput     = "stdout"
	DefaultLogMaxSize    = 100 // MB
	DefaultLogMaxBackups = 5
	DefaultLogMaxAge     = 30 // days
)

// File Paths
const (
	// Configuration paths
	DefaultConfigPath         = "~/.syntropy/config.yaml"
	DefaultNodeConfigDir      = "~/.syntropy/nodes"
	DefaultCacheDir           = "~/.syntropy/cache"
	DefaultLogDir             = "~/.syntropy/logs"
	DefaultTemplatesDir       = "~/.syntropy/templates"
	DefaultCloudInitOutputDir = "~/.syntropy/cloud-init"

	// Node configuration files
	NodeConfigFile = "node.json"
	NodeStateFile  = "state.json"
	NodeLogFile    = "node.log"
	NodeSSHKeyFile = "ssh_key"
	NodeCertFile   = "node_cert.pem"

	// Cloud-init template paths
	UserDataTemplate      = "templates/user-data.yaml"
	NetworkConfigTemplate = "templates/network-config.yaml"
	MetaDataTemplate      = "templates/meta-data.yaml"

	// Infrastructure paths
	InfrastructureDir   = "infrastructure"
	CloudInitDir        = "cloud-init"
	CloudInitScriptsDir = "scripts"

	// Node filesystem paths (on the actual node)
	NodeSyntropyDir = "/opt/syntropy"
	NodeConfigDir   = "/opt/syntropy/config"
	NodeBinDir      = "/opt/syntropy/bin"
	NodeLogDir      = "/opt/syntropy/logs"
	NodeDataDir     = "/opt/syntropy/data"
	NodeSSHDir      = "/home/syntropy/.ssh"
	NodeSystemdDir  = "/etc/systemd/system"

	// Node configuration files (on the actual node)
	NodeIDFile             = "/opt/syntropy/config/node_id"
	NodeGridTokenFile      = "/opt/syntropy/config/grid_token"
	NodeCommandStationFile = "/opt/syntropy/config/command_station_ip"
	NodeCertFileOnNode     = "/opt/syntropy/config/node_certificate.pem"
	NodeSSHKeyFileOnNode   = "/home/syntropy/.ssh/authorized_keys"
)

// Network Configuration
const (
	// Ports
	DefaultSSHPort   = 22
	DefaultHTTPPort  = 80
	DefaultHTTPSPort = 443
	DefaultDNSPort   = 53

	// Network timeouts
	DefaultConnectTimeout = 10 * time.Second
	DefaultReadTimeout    = 30 * time.Second
	DefaultWriteTimeout   = 30 * time.Second

	// Network retries
	DefaultMaxRetries    = 3
	DefaultRetryInterval = 5 * time.Second
)

// Security Constants
const (
	// Key sizes
	RSAKeySize     = 2048
	Ed25519KeySize = 32
	AESKeySize     = 32
	HMACKeySize    = 32

	// Hash sizes
	SHA256HashSize = 32
	SHA512HashSize = 64

	// Certificate fields
	DefaultCertSubject     = "CN=syntropy-node,O=syntropy-grid"
	DefaultCertValidity    = 365 * 24 * time.Hour
	DefaultCertKeyUsage    = "Digital Signature, Key Encipherment"
	DefaultCertExtKeyUsage = "Client Auth, Server Auth"
)

// Cloud-init Constants
const (
	// Cloud-init versions
	CloudInitVersion = "#cloud-config"

	// User configuration
	DefaultUserName   = "syntropy"
	DefaultUserGroups = "adm,audio,cdrom,dialout,dip,floppy,lxd,netdev,plugdev,sudo,video,docker"
	DefaultUserShell  = "/bin/bash"

	// Package lists
	BasePackages   = "curl,wget,git,build-essential,net-tools,htop,jq,ca-certificates,gnupg,lsb-release"
	DockerPackages = "docker.io,docker-compose"

	// Service names
	HeartbeatServiceName = "syntropy-heartbeat.service"
	RegisterServiceName  = "syntropy-register.service"
	DockerServiceName    = "docker.service"

	// Script names
	AutoRegisterScript = "auto-register.sh"
	HeartbeatScript    = "heartbeat.sh"
)

// Platform-specific Constants
const (
	// Windows
	WindowsWMICPath       = "wmic"
	WindowsDiskpartPath   = "diskpart"
	WindowsPowerShellPath = "powershell"
	WindowsDDPath         = "dd"

	// Linux
	LinuxLSBLKPath  = "lsblk"
	LinuxFdiskPath  = "fdisk"
	LinuxDDPath     = "dd"
	LinuxSyncPath   = "sync"
	LinuxNetcatPath = "nc"
	LinuxSSHPath    = "ssh"

	// macOS
	MacOSDiskutilPath = "diskutil"
	MacOSDDPath       = "dd"
	MacOSNetcatPath   = "nc"
)

// Error Messages
const (
	ErrMsgNoUSBFound         = "no USB devices found"
	ErrMsgUSBTooSmall        = "USB device too small (minimum 8GB required)"
	ErrMsgInvalidISO         = "invalid ISO file"
	ErrMsgTokenValidation    = "token validation failed"
	ErrMsgNodeNotFound       = "node not found"
	ErrMsgInvalidNodeID      = "invalid node ID format"
	ErrMsgHandshakeTimeout   = "handshake timeout"
	ErrMsgHeartbeatFailed    = "heartbeat failed"
	ErrMsgUSBWriteFailed     = "USB write operation failed"
	ErrMsgCloudInitInvalid   = "cloud-init validation failed"
	ErrMsgRegistrationFailed = "node registration failed"
	ErrMsgListenerFailed     = "listener failed to start"
	ErrMsgPermissionDenied   = "permission denied"
	ErrMsgNetworkError       = "network error"
	ErrMsgSystemDevice       = "cannot use system device"
)

// Success Messages
const (
	SuccessMsgNodeCreated        = "node created successfully"
	SuccessMsgNodeRegistered     = "node registered successfully"
	SuccessMsgUSBDetected        = "USB device detected"
	SuccessMsgISODownloaded      = "ISO downloaded successfully"
	SuccessMsgCloudInitGenerated = "cloud-init generated successfully"
	SuccessMsgUSBWritten         = "USB written successfully"
	SuccessMsgHandshakeComplete  = "handshake completed successfully"
	SuccessMsgHeartbeatActive    = "heartbeat active"
)

// Node States
const (
	StatePending     = "pending"
	StateActive      = "active"
	StateInactive    = "inactive"
	StateFailed      = "failed"
	StateCreating    = "creating"
	StateRegistering = "registering"
)

// Event Types
const (
	EventNodeCreated        = "node_created"
	EventNodeRegistered     = "node_registered"
	EventNodeDisconnected   = "node_disconnected"
	EventHeartbeatFailed    = "heartbeat_failed"
	EventNodeFailed         = "node_failed"
	EventUSBDetected        = "usb_detected"
	EventISODownloaded      = "iso_downloaded"
	EventCloudInitGenerated = "cloud_init_generated"
	EventUSBWritten         = "usb_written"
	EventHandshakeStarted   = "handshake_started"
	EventHandshakeComplete  = "handshake_complete"
	EventHeartbeatStarted   = "heartbeat_started"
	EventHeartbeatStopped   = "heartbeat_stopped"
)

// HTTP Status Codes
const (
	HTTPStatusOK                  = 200
	HTTPStatusCreated             = 201
	HTTPStatusBadRequest          = 400
	HTTPStatusUnauthorized        = 401
	HTTPStatusForbidden           = 403
	HTTPStatusNotFound            = 404
	HTTPStatusInternalServerError = 500
	HTTPStatusServiceUnavailable  = 503
)

// Time Formats
const (
	TimeFormatRFC3339  = "2006-01-02T15:04:05Z07:00"
	TimeFormatISO8601  = "2006-01-02T15:04:05.000Z"
	TimeFormatLog      = "2006-01-02 15:04:05"
	TimeFormatFilename = "2006-01-02_15-04-05"
)

// File Permissions
const (
	FileModeReadOnly   = 0400
	FileModeWriteOnly  = 0200
	FileModeReadWrite  = 0600
	FileModeExecutable = 0700
	FileModeDirectory  = 0755
	FileModeConfigFile = 0600
	FileModeLogFile    = 0644
)

// Buffer Sizes
const (
	DefaultBufferSize = 4096
	LargeBufferSize   = 65536
	MaxBufferSize     = 1048576 // 1MB
	NetworkBufferSize = 8192
	FileBufferSize    = 32768
)

// Retry Configuration
const (
	DefaultRetryDelay      = 1 * time.Second
	MaxRetryDelay          = 30 * time.Second
	RetryBackoffMultiplier = 2.0
)

// Validation Constants
const (
	MinNodeIDLength = 6
	MaxNodeIDLength = 20
	NodeIDPattern   = `^node-\d{2,3}$`
	MinTokenLength  = 32
	MaxTokenLength  = 256
	MinCertLength   = 100
	MaxCertLength   = 8192
)

// Performance Constants
const (
	DefaultConcurrency = 4
	MaxConcurrency     = 16
	DefaultBatchSize   = 100
	MaxBatchSize       = 1000
	DefaultCacheSize   = 100
	MaxCacheSize       = 1000
)

// Monitoring Constants
const (
	DefaultMetricsInterval     = 60 * time.Second
	DefaultHealthCheckInterval = 30 * time.Second
	DefaultLogRotationInterval = 24 * time.Hour
	DefaultCleanupInterval     = 7 * 24 * time.Hour
)

// Development Constants
const (
	DebugMode       = false
	VerboseLogging  = false
	EnableProfiling = false
	EnableMetrics   = true
	EnableTracing   = false
	TestMode        = false
)

// HTTP timeout for downloads
const (
	DefaultHTTPTimeout = 30 * time.Minute
)

// Ubuntu ISO Constants
const (
	// Ubuntu 24.04 LTS Server ISO
	Ubuntu2404ServerSHA256 = "a4acfda10b18da50e2ec50ccaf860d7f20b389df8765611142305c0e911d16fd"
	Ubuntu2404ServerSize   = int64(2040110000) // ~1.9GB

	// Ubuntu 22.04 LTS Server ISO
	Ubuntu2204ServerSHA256 = "10f2300071c92004b7293d9b0a9a0c9d7b6b6b6b6b6b6b6b6b6b6b6b6b6b6b6"
	Ubuntu2204ServerSize   = int64(1932740000) // ~1.8GB

	// Ubuntu 20.04 LTS Server ISO
	Ubuntu2004ServerSHA256 = "20f2300071c92004b7293d9b0a9a0c9d7b6b6b6b6b6b6b6b6b6b6b6b6b6b6b6"
	Ubuntu2004ServerSize   = int64(1825360000) // ~1.7GB
)
