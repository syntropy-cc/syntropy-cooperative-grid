package constants

// Application constants
const (
	// Application info
	AppName        = "Syntropy Cooperative Grid"
	AppVersion     = "1.0.0"
	AppDescription = "A comprehensive management system for the Syntropy Cooperative Grid"

	// API constants
	APIVersion     = "v1"
	APIPrefix      = "/api/" + APIVersion
	DefaultTimeout = 30 // seconds

	// Database constants
	DefaultPageSize = 50
	MaxPageSize     = 1000

	// File constants
	MaxFileSize = 100 * 1024 * 1024 // 100MB
	AllowedFileTypes = "jpg,jpeg,png,gif,pdf,txt,json,yaml,yml"

	// Network constants
	DefaultPort     = 8080
	DefaultHost     = "localhost"
	MaxConnections  = 1000
	ReadTimeout     = 30 // seconds
	WriteTimeout    = 30 // seconds
	IdleTimeout     = 120 // seconds

	// Security constants
	JWTExpirationTime = 24 * 60 * 60 // 24 hours in seconds
	BCryptCost        = 12
	MinPasswordLength = 8
	MaxPasswordLength = 128

	// Rate limiting constants
	RateLimitRequests = 1000
	RateLimitWindow   = 3600 // 1 hour in seconds

	// Cache constants
	CacheTTL     = 300 // 5 minutes in seconds
	CacheMaxSize = 1000

	// Logging constants
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
	LogLevelFatal = "fatal"

	// Environment constants
	EnvDevelopment = "development"
	EnvStaging     = "staging"
	EnvProduction  = "production"
	EnvTest        = "test"
)

// Node constants
const (
	// Node statuses
	NodeStatusCreating   = "creating"
	NodeStatusRunning    = "running"
	NodeStatusStopped    = "stopped"
	NodeStatusError      = "error"
	NodeStatusRestarting = "restarting"
	NodeStatusUpdating   = "updating"

	// Node types
	NodeTypePhysical = "physical"
	NodeTypeVirtual  = "virtual"
	NodeTypeCloud    = "cloud"

	// Node capabilities
	NodeCapabilityContainer = "container"
	NodeCapabilityStorage   = "storage"
	NodeCapabilityNetwork   = "network"
	NodeCapabilityCompute   = "compute"
)

// Container constants
const (
	// Container statuses
	ContainerStatusCreating   = "creating"
	ContainerStatusRunning    = "running"
	ContainerStatusStopped    = "stopped"
	ContainerStatusError      = "error"
	ContainerStatusRestarting = "restarting"
	ContainerStatusUpdating   = "updating"

	// Container types
	ContainerTypeApplication = "application"
	ContainerTypeService     = "service"
	ContainerTypeDatabase    = "database"
	ContainerTypeCache       = "cache"

	// Container resources
	DefaultCPU    = "100m"
	DefaultMemory = "128Mi"
	MaxCPU        = "2000m"
	MaxMemory     = "4Gi"
)

// Network constants
const (
	// Network statuses
	NetworkStatusActive   = "active"
	NetworkStatusInactive = "inactive"
	NetworkStatusError    = "error"

	// Network protocols
	NetworkProtocolTCP  = "tcp"
	NetworkProtocolUDP  = "udp"
	NetworkProtocolHTTP = "http"
	NetworkProtocolHTTPS = "https"
	NetworkProtocolWebSocket = "websocket"

	// Network types
	NetworkTypeMesh      = "mesh"
	NetworkTypeOverlay   = "overlay"
	NetworkTypeBridge    = "bridge"
	NetworkTypeHost      = "host"
)

// Cooperative constants
const (
	// Transaction types
	TransactionTypeServiceReward = "service_reward"
	TransactionTypeResourceUsage = "resource_usage"
	TransactionTypeParticipation = "participation"
	TransactionTypeTransfer      = "transfer"
	TransactionTypePenalty       = "penalty"
	TransactionTypeBonus         = "bonus"

	// Transaction statuses
	TransactionStatusPending   = "pending"
	TransactionStatusCompleted = "completed"
	TransactionStatusFailed    = "failed"
	TransactionStatusCancelled = "cancelled"

	// Proposal statuses
	ProposalStatusActive   = "active"
	ProposalStatusVoting   = "voting"
	ProposalStatusPassed   = "passed"
	ProposalStatusRejected = "rejected"
	ProposalStatusExpired  = "expired"

	// Vote types
	VoteTypeYes     = "yes"
	VoteTypeNo      = "no"
	VoteTypeAbstain = "abstain"

	// Trust levels
	TrustLevelUnknown  = "unknown"
	TrustLevelLow      = "low"
	TrustLevelMedium   = "medium"
	TrustLevelHigh     = "high"
	TrustLevelVeryHigh = "very_high"

	// Credit limits
	MinCreditAmount = 0.01
	MaxCreditAmount = 1000000.0
	DefaultCreditBalance = 100.0
)

// USB constants
const (
	// USB device types
	USBDeviceTypeFlash = "flash"
	USBDeviceTypeHDD   = "hdd"
	USBDeviceTypeSSD   = "ssd"

	// USB device statuses
	USBDeviceStatusAvailable = "available"
	USBDeviceStatusInUse     = "in_use"
	USBDeviceStatusError     = "error"
	USBDeviceStatusFormatting = "formatting"

	// USB device sizes
	MinUSBDeviceSize = 1024 * 1024 * 1024 // 1GB
	MaxUSBDeviceSize = 1024 * 1024 * 1024 * 1024 // 1TB

	// USB device formats
	USBDeviceFormatFAT32 = "fat32"
	USBDeviceFormatNTFS  = "ntfs"
	USBDeviceFormatEXT4  = "ext4"
	USBDeviceFormatXFS   = "xfs"
)

// Platform constants
const (
	// Operating systems
	OSWindows = "windows"
	OSLinux   = "linux"
	OSMacOS   = "macos"
	OSWSL     = "wsl"

	// Architectures
	ArchAMD64 = "amd64"
	ArchARM64 = "arm64"
	ArchARM   = "arm"
	Arch386   = "386"
)

// Configuration keys
const (
	// Database configuration
	ConfigKeyDatabaseURL      = "database.url"
	ConfigKeyDatabaseMaxConns = "database.max_connections"
	ConfigKeyDatabaseTimeout  = "database.timeout"

	// Redis configuration
	ConfigKeyRedisURL     = "redis.url"
	ConfigKeyRedisTimeout = "redis.timeout"

	// API configuration
	ConfigKeyAPIPort     = "api.port"
	ConfigKeyAPIHost     = "api.host"
	ConfigKeyAPITimeout  = "api.timeout"

	// Security configuration
	ConfigKeyJWTSecret     = "security.jwt_secret"
	ConfigKeyJWTExpiration = "security.jwt_expiration"
	ConfigKeyBCryptCost    = "security.bcrypt_cost"

	// Logging configuration
	ConfigKeyLogLevel  = "logging.level"
	ConfigKeyLogFormat = "logging.format"
	ConfigKeyLogOutput = "logging.output"

	// Monitoring configuration
	ConfigKeyMetricsEnabled = "monitoring.metrics_enabled"
	ConfigKeyMetricsPort    = "monitoring.metrics_port"
	ConfigKeyHealthPort     = "monitoring.health_port"
)

// Default configuration values
var (
	// Database defaults
	DefaultDatabaseURL      = "postgres://syntropy:syntropy@localhost:5432/syntropy?sslmode=disable"
	DefaultDatabaseMaxConns = 25
	DefaultDatabaseTimeout  = 30

	// Redis defaults
	DefaultRedisURL     = "redis://localhost:6379"
	DefaultRedisTimeout = 5

	// API defaults
	DefaultAPIPort     = 8080
	DefaultAPIHost     = "localhost"
	DefaultAPITimeout  = 30

	// Security defaults
	DefaultJWTExpiration = 24 * 60 * 60 // 24 hours
	DefaultBCryptCost    = 12

	// Logging defaults
	DefaultLogLevel  = LogLevelInfo
	DefaultLogFormat = "json"
	DefaultLogOutput = "stdout"

	// Monitoring defaults
	DefaultMetricsEnabled = true
	DefaultMetricsPort    = 9090
	DefaultHealthPort     = 8081
)
