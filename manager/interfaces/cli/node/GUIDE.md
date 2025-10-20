# Node Component - Implementation Guide

## 1. Component Overview

### Purpose
The Node Component is responsible for automated provisioning and registration of physical nodes in the Syntropy Cooperative Grid. It implements a plug-and-play system that creates bootable USB devices and automatically registers nodes without user intervention.

### Key Responsibilities
- **USB Detection**: Automatically detect and validate USB devices across platforms
- **Configuration Generation**: Generate secure NodeID, SSH keys, and certificates
- **Cloud-Init Creation**: Create complete installation scripts for Ubuntu Server
- **USB Writing**: Write bootable USB devices with custom configurations
- **Node Registration**: Handle secure handshake and registration protocols
- **Heartbeat Monitoring**: Maintain continuous connection monitoring
- **State Management**: Track node lifecycle from creation to active deployment

### Architecture Pattern
**Event-Driven Architecture** with two main subcomponents:
- **Create Subcomponent**: Handles node creation workflow
- **Registration Subcomponent**: Handles node registration and monitoring

## 2. File Structure

```
node/
├── GUIDE.md                     # This implementation guide
├── README.md                    # User documentation
├── TODO.md                      # Progress tracking
├── docs/                        # Complete documentation
│   ├── DOC.md                   # Quick start guide
│   ├── DEV.md                   # Developer documentation
│   ├── API.md                   # API reference
│   ├── TEST.md                  # Testing guide
│   └── LEARN.md                 # Learning journey
├── src/                         # Source code (READ-ONLY for tests)
│   ├── node.go                  # Main orchestrator (500 lines)
│   ├── types.go                 # Public types and interfaces
│   ├── create.go                # Create subcomponent (300 lines)
│   ├── registration.go          # Registration subcomponent (300 lines)
│   ├── usb_detector.go          # USB detection interface
│   ├── usb_detector_windows.go  # Windows USB detection
│   ├── usb_detector_linux.go    # Linux USB detection
│   ├── usb_detector_macos.go    # macOS USB detection
│   ├── iso_downloader.go        # ISO download and caching
│   ├── cloud_init_generator.go  # Cloud-init script generation
│   ├── usb_writer.go            # USB writing interface
│   ├── usb_writer_windows.go    # Windows USB writing
│   ├── usb_writer_linux.go      # Linux USB writing
│   ├── usb_writer_macos.go      # macOS USB writing
│   ├── auto_config_generator.go # Configuration generation
│   ├── token_integration.go     # Setup component integration
│   ├── handshake.go             # Handshake protocol
│   ├── listener.go              # TCP listener for registration
│   ├── heartbeat.go             # Heartbeat monitoring
│   ├── node_state.go            # Node state management
│   └── internal/                # Private implementation
│       ├── types/
│       │   └── interfaces.go    # Internal interfaces
│       ├── helpers/
│       │   └── helpers.go       # Utility functions
│       └── constants/
│           └── constants.go     # System constants
└── tests/                       # Test files (WRITE-ONLY)
    ├── unit/                    # Unit tests
    ├── integration/             # Integration tests
    ├── e2e/                     # End-to-end tests
    ├── performance/             # Performance tests
    ├── security/                # Security tests
    ├── fixtures/                # Test data
    ├── mocks/                   # Test doubles
    └── helpers/                 # Test utilities
```

## 3. Core Interfaces

### NodeManager Interface
```go
// Main orchestrator for node operations
type NodeManager struct {
    createSubcomponent      *CreateSubcomponent
    registrationSubcomponent *RegistrationSubcomponent
    eventBus               *EventBus
    nodeState              *NodeStateManager
    tokenIntegration       *TokenIntegration
    logger                 *Logger
    mutex                  sync.RWMutex
}

// Public API methods
func (nm *NodeManager) CreateNode(options *CreateOptions) (*CreateResult, error)
func (nm *NodeManager) ListNodes() (*NodeList, error)
func (nm *NodeManager) GetNodeStatus(nodeID string) (*NodeStatus, error)
func (nm *NodeManager) GetNodeLogs(nodeID string, options *LogOptions) (*NodeLogs, error)
func (nm *NodeManager) DeleteNode(nodeID string) error
func (nm *NodeManager) StartRegistrationListener() error
func (nm *NodeManager) StopRegistrationListener() error
```

### CreateSubcomponent Interface
```go
// Handles node creation workflow
type CreateSubcomponent struct {
    usbDetector       USBDetector
    isoDownloader     ISODownloader
    configGenerator   ConfigGenerator
    cloudInitGenerator CloudInitGenerator
    usbWriter         USBWriter
    logger            *Logger
    mutex             sync.RWMutex
}

// Public methods
func (cs *CreateSubcomponent) OrchestrateCreation(options *CreateOptions) (*CreateResult, error)
func (cs *CreateSubcomponent) DetectUSBDevices() ([]USBDevice, error)
func (cs *CreateSubcomponent) GenerateNodeConfig() (*NodeConfig, error)
func (cs *CreateSubcomponent) CreateBootableUSB(device USBDevice, config *NodeConfig) error
func (cs *CreateSubcomponent) ValidatePrerequisites() error
```

### RegistrationSubcomponent Interface
```go
// Handles node registration and monitoring
type RegistrationSubcomponent struct {
    listener          Listener
    handshake         HandshakeManager
    heartbeat         HeartbeatManager
    nodeState         NodeStateManager
    logger            *Logger
    mutex             sync.RWMutex
}

// Public methods
func (rs *RegistrationSubcomponent) StartListener(port int) error
func (rs *RegistrationSubcomponent) StopListener() error
func (rs *RegistrationSubcomponent) ProcessHandshake(conn net.Conn) error
func (rs *RegistrationSubcomponent) StartHeartbeat(nodeID string) error
func (rs *RegistrationSubcomponent) StopHeartbeat(nodeID string) error
func (rs *RegistrationSubcomponent) GetActiveNodes() []*NodeStatus
func (rs *RegistrationSubcomponent) GetPendingNodes() []*NodeStatus
```

### USBDetector Interface
```go
// Platform-agnostic USB detection
type USBDetector interface {
    DetectAvailable() ([]USBDevice, error)
    ValidateDevice(device USBDevice) error
    GetDeviceInfo(device USBDevice) (*DeviceInfo, error)
    IsSystemDevice(device USBDevice) bool
    SelectBestDevice(devices []USBDevice) (*USBDevice, error)
}

// Platform-specific implementations
type WindowsUSBDetector struct{ /* Windows-specific fields */ }
type LinuxUSBDetector struct{ /* Linux-specific fields */ }
type MacOSUSBDetector struct{ /* macOS-specific fields */ }
```

### USBWriter Interface
```go
// Platform-agnostic USB writing
type USBWriter interface {
    WriteToUSB(device USBDevice, isoPath string, cloudInit *CloudInitConfig) error
    ValidateWriteOperation(device USBDevice, isoPath string) error
    GetWriteProgress(device USBDevice) (*WriteProgress, error)
    CancelWriteOperation(device USBDevice) error
    VerifyWriteOperation(device USBDevice, isoPath string) error
}

// Platform-specific implementations
type WindowsUSBWriter struct{ /* Windows-specific fields */ }
type LinuxUSBWriter struct{ /* Linux-specific fields */ }
type MacOSUSBWriter struct{ /* macOS-specific fields */ }
```

### HandshakeManager Interface
```go
// Secure handshake protocol
type HandshakeManager struct {
    tokenValidator    TokenValidator
    certificateManager CertificateManager
    logger            *Logger
    mutex             sync.RWMutex
}

// Public methods
func (hm *HandshakeManager) ProcessNodeAnnouncement(announcement *NodeAnnouncement) (*HandshakeResponse, error)
func (hm *HandshakeManager) ValidateGridToken(token string) error
func (hm *HandshakeManager) ValidateNodeCertificate(cert string) error
func (hm *HandshakeManager) GenerateHandshakeResponse(nodeID string) (*HandshakeResponse, error)
func (hm *HandshakeManager) LogHandshakeEvent(nodeID string, event string, details map[string]interface{})
```

## 4. Data Structures

### Core Types
```go
// Node configuration structure
type NodeConfig struct {
    NodeID            string    `json:"node_id"`
    GridToken         string    `json:"grid_token"`
    SSHPublicKey      string    `json:"ssh_public_key"`
    SSHPrivateKey     string    `json:"ssh_private_key"`
    NodeCertificate   string    `json:"node_certificate"`
    CommandStationIP  string    `json:"command_station_ip"`
    CreatedAt         time.Time `json:"created_at"`
    ExpiresAt         time.Time `json:"expires_at"`
}

// USB device structure
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

// Node status structure
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

// Hardware information
type HardwareInfo struct {
    CPUCores    int    `json:"cpu_cores"`
    MemoryGB    int    `json:"memory_gb"`
    DiskGB      int    `json:"disk_gb"`
    Hostname    string `json:"hostname"`
    IPAddress   string `json:"ip_address"`
    OSVersion   string `json:"os_version"`
    KernelVersion string `json:"kernel_version"`
}

// Cloud-init configuration
type CloudInitConfig struct {
    UserData     string `yaml:"user_data"`
    NetworkConfig string `yaml:"network_config"`
    MetaData     string `yaml:"meta_data"`
    Valid        bool   `yaml:"valid"`
}

// Node announcement (handshake)
type NodeAnnouncement struct {
    Type           string        `json:"type"`
    NodeID         string        `json:"node_id"`
    GridToken      string        `json:"grid_token"`
    NodeCertificate string       `json:"node_certificate"`
    Hardware       *HardwareInfo `json:"hardware"`
    Timestamp      time.Time     `json:"timestamp"`
}

// Handshake response
type HandshakeResponse struct {
    Status              string                 `json:"status"`
    Message             string                 `json:"message"`
    CommandStationCert  string                 `json:"command_station_cert"`
    SSHConfig           map[string]interface{} `json:"ssh_config"`
    WorkloadConfig      map[string]interface{} `json:"workload_config"`
    Timestamp           time.Time              `json:"timestamp"`
}

// Event structure
type Event struct {
    Type      string                 `json:"type"`
    NodeID    string                 `json:"node_id"`
    Data      map[string]interface{} `json:"data"`
    Timestamp time.Time              `json:"timestamp"`
}

// Error types
type NodeError struct {
    Code      string `json:"code"`
    Message   string `json:"message"`
    Details   string `json:"details"`
    Timestamp time.Time `json:"timestamp"`
}

// Create options
type CreateOptions struct {
    USBPath       string `json:"usb_path"`
    ForceDownload bool   `json:"force_download"`
    ISOPath       string `json:"iso_path"`
    NodeID        string `json:"node_id"`
    Timeout       int    `json:"timeout"`
}

// Create result
type CreateResult struct {
    NodeID    string      `json:"node_id"`
    USBPath   string      `json:"usb_path"`
    Status    string      `json:"status"`
    Message   string      `json:"message"`
    CreatedAt time.Time   `json:"created_at"`
    Config    *NodeConfig `json:"config"`
}
```

## 5. Key Algorithms

### USB Detection Algorithm
```go
func (ud *USBDetector) DetectAvailable() ([]USBDevice, error) {
    // Platform-specific detection
    devices, err := ud.platformDetect()
    if err != nil {
        return nil, err
    }
    
    // Filter out system devices
    filtered := make([]USBDevice, 0)
    for _, device := range devices {
        if !ud.IsSystemDevice(device) && ud.ValidateDevice(device) {
            filtered = append(filtered, device)
        }
    }
    
    // Sort by capacity and speed
    sort.Slice(filtered, func(i, j int) bool {
        if filtered[i].Capacity != filtered[j].Capacity {
            return filtered[i].Capacity > filtered[j].Capacity
        }
        return filtered[i].Speed > filtered[j].Speed
    })
    
    return filtered, nil
}
```

### Configuration Generation Algorithm
```go
func (cg *ConfigGenerator) GenerateNodeConfig() (*NodeConfig, error) {
    // Generate unique NodeID
    nodeID, err := cg.generateNodeID()
    if err != nil {
        return nil, err
    }
    
    // Generate SSH key pair
    sshKeys, err := cg.generateSSHKeys()
    if err != nil {
        return nil, err
    }
    
    // Generate node certificate
    cert, err := cg.generateNodeCertificate(nodeID)
    if err != nil {
        return nil, err
    }
    
    // Get Grid Token from Setup component
    gridToken, err := cg.tokenIntegration.GetGridToken()
    if err != nil {
        return nil, err
    }
    
    // Detect Command Station IP
    commandStationIP, err := cg.detectCommandStationIP()
    if err != nil {
        return nil, err
    }
    
    return &NodeConfig{
        NodeID:           nodeID,
        GridToken:        gridToken,
        SSHPublicKey:     sshKeys.PublicKey,
        SSHPrivateKey:    sshKeys.PrivateKey,
        NodeCertificate:  cert,
        CommandStationIP: commandStationIP,
        CreatedAt:        time.Now(),
        ExpiresAt:        time.Now().Add(24 * time.Hour),
    }, nil
}
```

### Handshake Validation Algorithm
```go
func (hm *HandshakeManager) ProcessNodeAnnouncement(announcement *NodeAnnouncement) (*HandshakeResponse, error) {
    // Validate Grid Token
    if err := hm.ValidateGridToken(announcement.GridToken); err != nil {
        return nil, fmt.Errorf("invalid grid token: %w", err)
    }
    
    // Validate Node Certificate
    if err := hm.ValidateNodeCertificate(announcement.NodeCertificate); err != nil {
        return nil, fmt.Errorf("invalid node certificate: %w", err)
    }
    
    // Validate NodeID format
    if err := validateNodeID(announcement.NodeID); err != nil {
        return nil, fmt.Errorf("invalid node ID: %w", err)
    }
    
    // Check if node is in pending state
    if !hm.nodeState.IsNodePending(announcement.NodeID) {
        return nil, fmt.Errorf("node %s is not in pending state", announcement.NodeID)
    }
    
    // Generate response
    response := &HandshakeResponse{
        Status:             "accepted",
        Message:            "Node registered successfully",
        CommandStationCert: hm.certificateManager.GetCommandStationCert(),
        SSHConfig: map[string]interface{}{
            "port":          22,
            "key_algorithm": "rsa",
        },
        WorkloadConfig: map[string]interface{}{
            "docker_enabled": true,
            "resource_limits": map[string]interface{}{
                "cpu_cores": announcement.Hardware.CPUCores,
                "memory_gb": announcement.Hardware.MemoryGB,
            },
        },
        Timestamp: time.Now(),
    }
    
    // Update node state
    hm.nodeState.TransitionToActive(announcement.NodeID, announcement.IPAddress)
    
    return response, nil
}
```

### Cloud-Init Generation Algorithm
```go
func (cig *CloudInitGenerator) GenerateCloudInit(config *NodeConfig) (*CloudInitConfig, error) {
    // Load template
    template, err := cig.loadTemplate("user-data.yaml")
    if err != nil {
        return nil, err
    }
    
    // Prepare template variables
    variables := map[string]interface{}{
        "NodeID":            config.NodeID,
        "GridToken":         config.GridToken,
        "CommandStationIP":  config.CommandStationIP,
        "SSHPublicKey":      config.SSHPublicKey,
        "NodeCertificate":   config.NodeCertificate,
        "CreatedAt":         config.CreatedAt.Format(time.RFC3339),
    }
    
    // Execute template
    userData, err := cig.executeTemplate(template, variables)
    if err != nil {
        return nil, err
    }
    
    // Generate network config
    networkConfig, err := cig.generateNetworkConfig()
    if err != nil {
        return nil, err
    }
    
    // Generate meta data
    metaData, err := cig.generateMetaData(config.NodeID)
    if err != nil {
        return nil, err
    }
    
    // Validate YAML
    if err := cig.validateYAML(userData); err != nil {
        return nil, err
    }
    
    return &CloudInitConfig{
        UserData:     userData,
        NetworkConfig: networkConfig,
        MetaData:     metaData,
        Valid:        true,
    }, nil
}
```

## 6. Error Handling

### Error Types
```go
// Custom error types
var (
    ErrNoUSBFound        = errors.New("no USB devices found")
    ErrUSBTooSmall       = errors.New("USB device too small")
    ErrInvalidISO        = errors.New("invalid ISO file")
    ErrTokenValidation   = errors.New("token validation failed")
    ErrNodeNotFound      = errors.New("node not found")
    ErrInvalidNodeID     = errors.New("invalid node ID")
    ErrHandshakeTimeout  = errors.New("handshake timeout")
    ErrHeartbeatFailed   = errors.New("heartbeat failed")
    ErrUSBWriteFailed    = errors.New("USB write failed")
    ErrCloudInitInvalid  = errors.New("cloud-init validation failed")
)

// Error handling with context
func (nm *NodeManager) CreateNode(options *CreateOptions) (*CreateResult, error) {
    // Validate prerequisites
    if err := nm.validatePrerequisites(); err != nil {
        return nil, fmt.Errorf("prerequisites validation failed: %w", err)
    }
    
    // Create node
    result, err := nm.createSubcomponent.OrchestrateCreation(options)
    if err != nil {
        // Log error with context
        nm.logger.Error("node creation failed",
            "error", err,
            "options", options,
            "timestamp", time.Now(),
        )
        return nil, fmt.Errorf("node creation failed: %w", err)
    }
    
    // Publish event
    nm.eventBus.Publish(Event{
        Type:      "node_created",
        NodeID:    result.NodeID,
        Data:      map[string]interface{}{"result": result},
        Timestamp: time.Now(),
    })
    
    return result, nil
}
```

### Error Recovery
```go
// Retry mechanism with exponential backoff
func (nm *NodeManager) createNodeWithRetry(options *CreateOptions) (*CreateResult, error) {
    maxRetries := 3
    baseDelay := time.Second
    
    for attempt := 1; attempt <= maxRetries; attempt++ {
        result, err := nm.createSubcomponent.OrchestrateCreation(options)
        if err == nil {
            return result, nil
        }
        
        // Check if error is retryable
        if !isRetryableError(err) {
            return nil, err
        }
        
        // Calculate delay with exponential backoff
        delay := time.Duration(attempt) * baseDelay
        nm.logger.Warn("node creation attempt failed, retrying",
            "attempt", attempt,
            "error", err,
            "delay", delay,
        )
        
        time.Sleep(delay)
    }
    
    return nil, fmt.Errorf("node creation failed after %d attempts", maxRetries)
}

// Check if error is retryable
func isRetryableError(err error) bool {
    switch {
    case errors.Is(err, ErrUSBWriteFailed):
        return true
    case errors.Is(err, ErrHandshakeTimeout):
        return true
    case errors.Is(err, ErrHeartbeatFailed):
        return true
    default:
        return false
    }
}
```

## 7. Configuration Management

### Configuration Structure
```go
// Node component configuration
type NodeConfig struct {
    // USB detection settings
    USBDetection struct {
        MinCapacity    int64         `yaml:"min_capacity"`
        MaxCapacity    int64         `yaml:"max_capacity"`
        PreferredSpeed int           `yaml:"preferred_speed"`
        ScanInterval   time.Duration `yaml:"scan_interval"`
    } `yaml:"usb_detection"`
    
    // ISO download settings
    ISODownload struct {
        URL            string        `yaml:"url"`
        CacheDir       string        `yaml:"cache_dir"`
        SHA256Checksum string        `yaml:"sha256_checksum"`
        Timeout        time.Duration `yaml:"timeout"`
    } `yaml:"iso_download"`
    
    // Registration settings
    Registration struct {
        Port           int           `yaml:"port"`
        Timeout        time.Duration `yaml:"timeout"`
        MaxConcurrent  int           `yaml:"max_concurrent"`
        RetryAttempts  int           `yaml:"retry_attempts"`
    } `yaml:"registration"`
    
    // Heartbeat settings
    Heartbeat struct {
        Interval       time.Duration `yaml:"interval"`
        Timeout        time.Duration `yaml:"timeout"`
        MaxFailures    int           `yaml:"max_failures"`
        RetryInterval  time.Duration `yaml:"retry_interval"`
    } `yaml:"heartbeat"`
    
    // Security settings
    Security struct {
        TokenExpiry    time.Duration `yaml:"token_expiry"`
        CertExpiry     time.Duration `yaml:"cert_expiry"`
        KeySize        int           `yaml:"key_size"`
        HashAlgorithm  string        `yaml:"hash_algorithm"`
    } `yaml:"security"`
    
    // Logging settings
    Logging struct {
        Level          string        `yaml:"level"`
        Format         string        `yaml:"format"`
        Output         string        `yaml:"output"`
        MaxSize        int           `yaml:"max_size"`
        MaxBackups     int           `yaml:"max_backups"`
        MaxAge         int           `yaml:"max_age"`
    } `yaml:"logging"`
}
```

### Configuration Loading
```go
// Load configuration from file
func LoadConfig(configPath string) (*NodeConfig, error) {
    // Check if config file exists
    if _, err := os.Stat(configPath); os.IsNotExist(err) {
        // Create default configuration
        config := getDefaultConfig()
        if err := SaveConfig(config, configPath); err != nil {
            return nil, fmt.Errorf("failed to create default config: %w", err)
        }
        return config, nil
    }
    
    // Read configuration file
    data, err := ioutil.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }
    
    // Parse YAML
    var config NodeConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse config file: %w", err)
    }
    
    // Validate configuration
    if err := validateConfig(&config); err != nil {
        return nil, fmt.Errorf("invalid configuration: %w", err)
    }
    
    return &config, nil
}

// Default configuration
func getDefaultConfig() *NodeConfig {
    return &NodeConfig{
        USBDetection: struct {
            MinCapacity    int64         `yaml:"min_capacity"`
            MaxCapacity    int64         `yaml:"max_capacity"`
            PreferredSpeed int           `yaml:"preferred_speed"`
            ScanInterval   time.Duration `yaml:"scan_interval"`
        }{
            MinCapacity:    8 * 1024 * 1024 * 1024, // 8GB
            MaxCapacity:    64 * 1024 * 1024 * 1024, // 64GB
            PreferredSpeed: 3, // USB 3.0
            ScanInterval:   5 * time.Second,
        },
        ISODownload: struct {
            URL            string        `yaml:"url"`
            CacheDir       string        `yaml:"cache_dir"`
            SHA256Checksum string        `yaml:"sha256_checksum"`
            Timeout        time.Duration `yaml:"timeout"`
        }{
            URL:            "https://releases.ubuntu.com/24.04/ubuntu-24.04-server-amd64.iso",
            CacheDir:       "~/.syntropy/cache/isos",
            SHA256Checksum: "a4acfda10b18da50e2ec50ccaf860d7f3b4b165bf3f4e4604dac366b93d1b25",
            Timeout:        30 * time.Minute,
        },
        Registration: struct {
            Port           int           `yaml:"port"`
            Timeout        time.Duration `yaml:"timeout"`
            MaxConcurrent  int           `yaml:"max_concurrent"`
            RetryAttempts  int           `yaml:"retry_attempts"`
        }{
            Port:           51000,
            Timeout:        30 * time.Minute,
            MaxConcurrent:  10,
            RetryAttempts:  3,
        },
        Heartbeat: struct {
            Interval       time.Duration `yaml:"interval"`
            Timeout        time.Duration `yaml:"timeout"`
            MaxFailures    int           `yaml:"max_failures"`
            RetryInterval  time.Duration `yaml:"retry_interval"`
        }{
            Interval:       30 * time.Second,
            Timeout:        10 * time.Second,
            MaxFailures:    3,
            RetryInterval:  5 * time.Second,
        },
        Security: struct {
            TokenExpiry    time.Duration `yaml:"token_expiry"`
            CertExpiry     time.Duration `yaml:"cert_expiry"`
            KeySize        int           `yaml:"key_size"`
            HashAlgorithm  string        `yaml:"hash_algorithm"`
        }{
            TokenExpiry:    24 * time.Hour,
            CertExpiry:     365 * 24 * time.Hour,
            KeySize:        2048,
            HashAlgorithm:  "SHA256",
        },
        Logging: struct {
            Level          string        `yaml:"level"`
            Format         string        `yaml:"format"`
            Output         string        `yaml:"output"`
            MaxSize        int           `yaml:"max_size"`
            MaxBackups     int           `yaml:"max_backups"`
            MaxAge         int           `yaml:"max_age"`
        }{
            Level:          "info",
            Format:         "json",
            Output:         "stdout",
            MaxSize:        100, // MB
            MaxBackups:     5,
            MaxAge:         30, // days
        },
    }
}
```

## 8. Security Considerations

### Security Architecture
The Node Component implements a 3-layer security architecture:

#### Layer 1: Grid Token (Global Level)
- **Source**: Setup Component via system Keyring
- **Usage**: Initial node authentication in the grid
- **Storage**: Command Station in system Keyring, Node in encrypted file
- **Validation**: Handshake initial validation

#### Layer 2: Node Certificate (Node Level)
- **Source**: Generated automatically during node creation
- **Usage**: Unique and permanent node identification
- **Format**: Ed25519 key pair
- **Storage**: Command Station in node directory, Node in encrypted file
- **Validation**: All communications after registration

#### Layer 3: SSH Keys (Communication Level)
- **Source**: Generated automatically during node creation
- **Usage**: SSH communication between Command Station and Node
- **Format**: RSA 2048 bits
- **Storage**: Command Station in node directory, Node in SSH directory
- **Validation**: SSH connections

### Security Implementation
```go
// Token validation
func (hm *HandshakeManager) ValidateGridToken(token string) error {
    // Get token from Setup component
    validToken, err := hm.tokenIntegration.GetGridToken()
    if err != nil {
        return fmt.Errorf("failed to get grid token: %w", err)
    }
    
    // Compare tokens
    if !hm.compareTokens(token, validToken) {
        return ErrTokenValidation
    }
    
    return nil
}

// Certificate validation
func (hm *HandshakeManager) ValidateNodeCertificate(cert string) error {
    // Parse certificate
    block, _ := pem.Decode([]byte(cert))
    if block == nil {
        return fmt.Errorf("failed to decode certificate")
    }
    
    // Validate certificate
    if err := hm.certificateManager.ValidateCertificate(block.Bytes); err != nil {
        return fmt.Errorf("certificate validation failed: %w", err)
    }
    
    return nil
}

// SSH key generation
func (cg *ConfigGenerator) generateSSHKeys() (*SSHKeys, error) {
    // Generate RSA key pair
    privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        return nil, fmt.Errorf("failed to generate private key: %w", err)
    }
    
    // Get public key
    publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
    if err != nil {
        return nil, fmt.Errorf("failed to get public key: %w", err)
    }
    
    // Encode keys
    privateKeyPEM := pem.EncodeToMemory(&pem.Block{
        Type:  "RSA PRIVATE KEY",
        Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
    })
    
    publicKeySSH := ssh.MarshalAuthorizedKey(publicKey)
    
    return &SSHKeys{
        PrivateKey: string(privateKeyPEM),
        PublicKey:  string(publicKeySSH),
    }, nil
}
```

### Security Best Practices
```go
// Secure file permissions
func setSecurePermissions(filePath string) error {
    // Set restrictive permissions
    if err := os.Chmod(filePath, 0600); err != nil {
        return fmt.Errorf("failed to set file permissions: %w", err)
    }
    
    return nil
}

// Secure data encryption
func encryptSensitiveData(data []byte, key []byte) ([]byte, error) {
    // Create AES cipher
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, fmt.Errorf("failed to create cipher: %w", err)
    }
    
    // Create GCM mode
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, fmt.Errorf("failed to create GCM: %w", err)
    }
    
    // Generate nonce
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, fmt.Errorf("failed to generate nonce: %w", err)
    }
    
    // Encrypt data
    ciphertext := gcm.Seal(nonce, nonce, data, nil)
    
    return ciphertext, nil
}

// Input validation
func validateNodeID(nodeID string) error {
    // Check format: node-XX where XX is number
    pattern := regexp.MustCompile(`^node-\d{2,3}$`)
    if !pattern.MatchString(nodeID) {
        return fmt.Errorf("invalid node ID format: %s", nodeID)
    }
    
    // Check length
    if len(nodeID) < 6 || len(nodeID) > 10 {
        return fmt.Errorf("node ID length invalid: %d", len(nodeID))
    }
    
    return nil
}
```

## 9. Testing Strategy

### Test Structure
```go
// Test file organization
// tests/unit/node_test.go
func TestNodeManager_CreateNode(t *testing.T) {
    // Setup test environment
    nodeManager := setupTestNodeManager(t)
    defer cleanupTestEnvironment(t)
    
    // Test cases
    tests := []struct {
        name    string
        options *CreateOptions
        wantErr bool
    }{
        {
            name: "successful node creation",
            options: &CreateOptions{
                USBPath: "/dev/sdb",
            },
            wantErr: false,
        },
        {
            name: "invalid USB path",
            options: &CreateOptions{
                USBPath: "/dev/invalid",
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := nodeManager.CreateNode(tt.options)
            if (err != nil) != tt.wantErr {
                t.Errorf("CreateNode() error = %v, wantErr %v", err, tt.wantErr)
            }
            if !tt.wantErr && result == nil {
                t.Error("CreateNode() returned nil result for successful case")
            }
        })
    }
}

// tests/integration/plug_and_play_test.go
func TestPlugAndPlayWorkflow(t *testing.T) {
    // Setup integration test environment
    env := setupIntegrationTest(t)
    defer env.Cleanup()
    
    // Test complete workflow
    t.Run("complete workflow", func(t *testing.T) {
        // Create node
        result, err := env.NodeManager.CreateNode(&CreateOptions{
            USBPath: env.TestUSBPath,
        })
        require.NoError(t, err)
        require.NotNil(t, result)
        
        // Start registration listener
        err = env.NodeManager.StartRegistrationListener()
        require.NoError(t, err)
        
        // Simulate node connection
        err = env.SimulateNodeConnection(result.NodeID)
        require.NoError(t, err)
        
        // Verify node registration
        status, err := env.NodeManager.GetNodeStatus(result.NodeID)
        require.NoError(t, err)
        assert.Equal(t, "active", status.Status)
    })
}

// tests/e2e/scenarios/single_node_creation_test.go
func TestSingleNodeCreation(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping E2E test in short mode")
    }
    
    // Setup E2E test environment
    env := setupE2ETest(t)
    defer env.Cleanup()
    
    // Test with real hardware
    t.Run("real hardware test", func(t *testing.T) {
        // Create node with real USB
        result, err := env.NodeManager.CreateNode(&CreateOptions{
            USBPath: env.RealUSBPath,
        })
        require.NoError(t, err)
        
        // Verify USB was created
        assert.FileExists(t, result.USBPath)
        
        // Test node registration
        err = env.TestNodeRegistration(result.NodeID)
        require.NoError(t, err)
    })
}
```

### Mock Implementations
```go
// tests/mocks/usb_detector_mock.go
type MockUSBDetector struct {
    devices []USBDevice
    err     error
}

func (m *MockUSBDetector) DetectAvailable() ([]USBDevice, error) {
    if m.err != nil {
        return nil, m.err
    }
    return m.devices, nil
}

func (m *MockUSBDetector) ValidateDevice(device USBDevice) error {
    if device.Capacity < 8*1024*1024*1024 {
        return ErrUSBTooSmall
    }
    return nil
}

func (m *MockUSBDetector) GetDeviceInfo(device USBDevice) (*DeviceInfo, error) {
    return &DeviceInfo{
        Path:     device.Path,
        Capacity: device.Capacity,
        Speed:    device.Speed,
    }, nil
}

func (m *MockUSBDetector) IsSystemDevice(device USBDevice) bool {
    return device.IsSystem
}

func (m *MockUSBDetector) SelectBestDevice(devices []USBDevice) (*USBDevice, error) {
    if len(devices) == 0 {
        return nil, ErrNoUSBFound
    }
    
    // Select device with largest capacity
    best := &devices[0]
    for i := 1; i < len(devices); i++ {
        if devices[i].Capacity > best.Capacity {
            best = &devices[i]
        }
    }
    
    return best, nil
}

// tests/mocks/token_integration_mock.go
type MockTokenIntegration struct {
    gridToken string
    err       error
}

func (m *MockTokenIntegration) GetGridToken() (string, error) {
    if m.err != nil {
        return "", m.err
    }
    return m.gridToken, nil
}

func (m *MockTokenIntegration) ValidateToken(token string) error {
    if token != m.gridToken {
        return ErrTokenValidation
    }
    return nil
}
```

### Test Coverage Requirements
```go
// Coverage requirements
// 100% line coverage for all source files
// 100% branch coverage for all conditional statements
// 100% function coverage for all public and private functions

// Example coverage test
func TestCoverage(t *testing.T) {
    // Run tests with coverage
    cmd := exec.Command("go", "test", "-coverprofile=coverage.out", "./...")
    err := cmd.Run()
    require.NoError(t, err)
    
    // Check coverage
    cmd = exec.Command("go", "tool", "cover", "-func=coverage.out")
    output, err := cmd.Output()
    require.NoError(t, err)
    
    // Parse coverage output
    lines := strings.Split(string(output), "\n")
    for _, line := range lines {
        if strings.Contains(line, "total:") {
            // Extract coverage percentage
            parts := strings.Fields(line)
            if len(parts) >= 3 {
                coverage := strings.TrimSuffix(parts[2], "%")
                coverageFloat, err := strconv.ParseFloat(coverage, 64)
                require.NoError(t, err)
                
                // Assert 100% coverage
                assert.Equal(t, 100.0, coverageFloat, "coverage must be 100%%")
            }
        }
    }
}
```

## 10. Performance Considerations

### Performance Targets
```go
// Performance benchmarks
func BenchmarkUSBDetection(b *testing.B) {
    detector := &USBDetector{}
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := detector.DetectAvailable()
        if err != nil {
            b.Fatal(err)
        }
    }
}

func BenchmarkNodeCreation(b *testing.B) {
    nodeManager := setupTestNodeManager(b)
    defer cleanupTestEnvironment(b)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := nodeManager.CreateNode(&CreateOptions{
            USBPath: "/dev/sdb",
        })
        if err != nil {
            b.Fatal(err)
        }
    }
}

func BenchmarkHandshake(b *testing.B) {
    handshakeManager := &HandshakeManager{}
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        announcement := &NodeAnnouncement{
            Type:           "node_announcement",
            NodeID:         "node-01",
            GridToken:      "valid-token",
            NodeCertificate: "valid-cert",
            Hardware:       &HardwareInfo{CPUCores: 8, MemoryGB: 16},
            Timestamp:      time.Now(),
        }
        
        _, err := handshakeManager.ProcessNodeAnnouncement(announcement)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

### Performance Optimization
```go
// Concurrent USB operations
func (cs *CreateSubcomponent) OrchestrateCreation(options *CreateOptions) (*CreateResult, error) {
    // Use goroutines for concurrent operations
    var wg sync.WaitGroup
    var devices []USBDevice
    var config *NodeConfig
    var err error
    
    // Concurrent USB detection and config generation
    wg.Add(2)
    
    go func() {
        defer wg.Done()
        devices, err = cs.usbDetector.DetectAvailable()
    }()
    
    go func() {
        defer wg.Done()
        config, err = cs.configGenerator.GenerateNodeConfig()
    }()
    
    wg.Wait()
    if err != nil {
        return nil, err
    }
    
    // Select best device
    device, err := cs.usbDetector.SelectBestDevice(devices)
    if err != nil {
        return nil, err
    }
    
    // Create bootable USB
    if err := cs.usbWriter.WriteToUSB(*device, options.ISOPath, config.CloudInit); err != nil {
        return nil, err
    }
    
    return &CreateResult{
        NodeID:    config.NodeID,
        USBPath:   device.Path,
        Status:    "created",
        Message:   "Node created successfully",
        CreatedAt: time.Now(),
        Config:    config,
    }, nil
}

// Connection pooling for heartbeat
type ConnectionPool struct {
    connections map[string]*Connection
    mutex       sync.RWMutex
}

func (cp *ConnectionPool) GetConnection(nodeID string) (*Connection, error) {
    cp.mutex.RLock()
    conn, exists := cp.connections[nodeID]
    cp.mutex.RUnlock()
    
    if exists && conn.IsHealthy() {
        return conn, nil
    }
    
    // Create new connection
    cp.mutex.Lock()
    defer cp.mutex.Unlock()
    
    conn, err := cp.createConnection(nodeID)
    if err != nil {
        return nil, err
    }
    
    cp.connections[nodeID] = conn
    return conn, nil
}
```

### Memory Management
```go
// Memory-efficient USB device handling
func (ud *USBDetector) DetectAvailable() ([]USBDevice, error) {
    // Use streaming for large device lists
    devices := make([]USBDevice, 0, 10) // Pre-allocate capacity
    
    // Process devices in batches
    batchSize := 100
    for offset := 0; ; offset += batchSize {
        batch, err := ud.detectBatch(offset, batchSize)
        if err != nil {
            return nil, err
        }
        
        if len(batch) == 0 {
            break
        }
        
        // Filter and add to result
        for _, device := range batch {
            if ud.ValidateDevice(device) {
                devices = append(devices, device)
            }
        }
    }
    
    return devices, nil
}

// Garbage collection optimization
func (nm *NodeManager) cleanupInactiveNodes() {
    nm.mutex.Lock()
    defer nm.mutex.Unlock()
    
    // Remove nodes inactive for more than 24 hours
    cutoff := time.Now().Add(-24 * time.Hour)
    for nodeID, status := range nm.nodeState.GetInactiveNodes() {
        if status.LastHeartbeat.Before(cutoff) {
            nm.nodeState.RemoveNode(nodeID)
        }
    }
}
```

## 11. Integration Points

### Setup Component Integration
```go
// Token integration with Setup component
type TokenIntegration struct {
    setupTokenManager *setup.TokenManager
    logger            *Logger
}

func (ti *TokenIntegration) GetGridToken() (string, error) {
    // Get token from Setup component
    token, err := ti.setupTokenManager.GetGridToken()
    if err != nil {
        return "", fmt.Errorf("failed to get grid token from setup: %w", err)
    }
    
    // Validate token format
    if err := ti.validateTokenFormat(token); err != nil {
        return "", fmt.Errorf("invalid token format: %w", err)
    }
    
    return token, nil
}

func (ti *TokenIntegration) ValidateToken(token string) error {
    // Validate against Setup component
    if err := ti.setupTokenManager.ValidateToken(token); err != nil {
        return fmt.Errorf("token validation failed: %w", err)
    }
    
    return nil
}

// CLI integration
func init() {
    // Add node commands to CLI
    rootCmd.AddCommand(nodeCmd)
    
    nodeCmd.AddCommand(nodeCreateCmd)
    nodeCmd.AddCommand(nodeListCmd)
    nodeCmd.AddCommand(nodeStatusCmd)
    nodeCmd.AddCommand(nodeLogsCmd)
}

// Node create command
var nodeCreateCmd = &cobra.Command{
    Use:   "create",
    Short: "Create a new node",
    Long:  "Create a new node with automatic USB detection and configuration",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Parse command options
        options := &CreateOptions{
            USBPath:       cmd.Flag("usb").Value.String(),
            ForceDownload: cmd.Flag("force-download").Changed,
            ISOPath:       cmd.Flag("iso").Value.String(),
        }
        
        // Create node
        nodeManager := node.NewNodeManager()
        result, err := nodeManager.CreateNode(options)
        if err != nil {
            return fmt.Errorf("failed to create node: %w", err)
        }
        
        // Display result
        fmt.Printf("Node created successfully: %s\n", result.NodeID)
        fmt.Printf("USB device: %s\n", result.USBPath)
        fmt.Printf("Status: %s\n", result.Status)
        
        return nil
    },
}
```

### Future Component Integration
```go
// Workload component integration (future)
type WorkloadIntegration struct {
    workloadManager *workload.WorkloadManager
    logger          *Logger
}

func (wi *WorkloadIntegration) DeployWorkload(nodeID string, workload *workload.Workload) error {
    // Check if node is active
    if !wi.isNodeActive(nodeID) {
        return fmt.Errorf("node %s is not active", nodeID)
    }
    
    // Deploy workload
    if err := wi.workloadManager.Deploy(nodeID, workload); err != nil {
        return fmt.Errorf("failed to deploy workload: %w", err)
    }
    
    return nil
}

// Management component integration (future)
type ManagementIntegration struct {
    managementManager *management.ManagementManager
    logger            *Logger
}

func (mi *ManagementIntegration) GetNodeMetrics(nodeID string) (*management.NodeMetrics, error) {
    // Get metrics from management component
    metrics, err := mi.managementManager.GetNodeMetrics(nodeID)
    if err != nil {
        return nil, fmt.Errorf("failed to get node metrics: %w", err)
    }
    
    return metrics, nil
}
```

## 12. Deployment and Operations

### Deployment Configuration
```go
// Deployment configuration
type DeploymentConfig struct {
    Environment string `yaml:"environment"`
    Version     string `yaml:"version"`
    LogLevel    string `yaml:"log_level"`
    
    // Service configuration
    Service struct {
        Name        string `yaml:"name"`
        Port        int    `yaml:"port"`
        HealthCheck string `yaml:"health_check"`
    } `yaml:"service"`
    
    // Database configuration
    Database struct {
        Host     string `yaml:"host"`
        Port     int    `yaml:"port"`
        Database string `yaml:"database"`
        Username string `yaml:"username"`
        Password string `yaml:"password"`
    } `yaml:"database"`
    
    // Monitoring configuration
    Monitoring struct {
        Enabled  bool   `yaml:"enabled"`
        Endpoint string `yaml:"endpoint"`
        Interval int    `yaml:"interval"`
    } `yaml:"monitoring"`
}

// Load deployment configuration
func LoadDeploymentConfig() (*DeploymentConfig, error) {
    configPath := os.Getenv("SYNTHROPY_CONFIG_PATH")
    if configPath == "" {
        configPath = "/etc/syntropy/config.yaml"
    }
    
    config, err := LoadConfig(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to load deployment config: %w", err)
    }
    
    return config, nil
}
```

### Health Checks
```go
// Health check implementation
type HealthChecker struct {
    nodeManager *NodeManager
    logger      *Logger
}

func (hc *HealthChecker) CheckHealth() *HealthStatus {
    status := &HealthStatus{
        Status:    "healthy",
        Timestamp: time.Now(),
        Checks:    make(map[string]CheckResult),
    }
    
    // Check USB detection
    if err := hc.checkUSBDetection(); err != nil {
        status.Checks["usb_detection"] = CheckResult{
            Status: "unhealthy",
            Error:  err.Error(),
        }
        status.Status = "unhealthy"
    } else {
        status.Checks["usb_detection"] = CheckResult{
            Status: "healthy",
        }
    }
    
    // Check registration listener
    if err := hc.checkRegistrationListener(); err != nil {
        status.Checks["registration_listener"] = CheckResult{
            Status: "unhealthy",
            Error:  err.Error(),
        }
        status.Status = "unhealthy"
    } else {
        status.Checks["registration_listener"] = CheckResult{
            Status: "healthy",
        }
    }
    
    // Check active nodes
    if err := hc.checkActiveNodes(); err != nil {
        status.Checks["active_nodes"] = CheckResult{
            Status: "unhealthy",
            Error:  err.Error(),
        }
        status.Status = "unhealthy"
    } else {
        status.Checks["active_nodes"] = CheckResult{
            Status: "healthy",
        }
    }
    
    return status
}

func (hc *HealthChecker) checkUSBDetection() error {
    devices, err := hc.nodeManager.createSubcomponent.usbDetector.DetectAvailable()
    if err != nil {
        return fmt.Errorf("USB detection failed: %w", err)
    }
    
    if len(devices) == 0 {
        return errors.New("no USB devices detected")
    }
    
    return nil
}

func (hc *HealthChecker) checkRegistrationListener() error {
    if !hc.nodeManager.registrationSubcomponent.listener.IsRunning() {
        return errors.New("registration listener is not running")
    }
    
    return nil
}

func (hc *HealthChecker) checkActiveNodes() error {
    activeNodes := hc.nodeManager.registrationSubcomponent.GetActiveNodes()
    
    // Check if any nodes are unhealthy
    for _, node := range activeNodes {
        if time.Since(node.LastHeartbeat) > 5*time.Minute {
            return fmt.Errorf("node %s has not sent heartbeat for %v", 
                node.NodeID, time.Since(node.LastHeartbeat))
        }
    }
    
    return nil
}
```

### Monitoring and Metrics
```go
// Metrics collection
type MetricsCollector struct {
    nodeManager *NodeManager
    logger      *Logger
}

func (mc *MetricsCollector) CollectMetrics() *Metrics {
    metrics := &Metrics{
        Timestamp: time.Now(),
        Counters:  make(map[string]int64),
        Gauges:    make(map[string]float64),
        Histograms: make(map[string]*Histogram),
    }
    
    // Node creation metrics
    metrics.Counters["nodes_created_total"] = mc.getNodesCreatedTotal()
    metrics.Counters["nodes_registered_total"] = mc.getNodesRegisteredTotal()
    metrics.Counters["nodes_failed_total"] = mc.getNodesFailedTotal()
    
    // Active node metrics
    activeNodes := mc.nodeManager.registrationSubcomponent.GetActiveNodes()
    metrics.Gauges["active_nodes"] = float64(len(activeNodes))
    
    // USB detection metrics
    devices, err := mc.nodeManager.createSubcomponent.usbDetector.DetectAvailable()
    if err == nil {
        metrics.Gauges["usb_devices_detected"] = float64(len(devices))
    }
    
    // Heartbeat metrics
    metrics.Histograms["heartbeat_latency"] = mc.getHeartbeatLatencyHistogram()
    
    return metrics
}

func (mc *MetricsCollector) getNodesCreatedTotal() int64 {
    // Implementation to get total nodes created
    return 0 // Placeholder
}

func (mc *MetricsCollector) getNodesRegisteredTotal() int64 {
    // Implementation to get total nodes registered
    return 0 // Placeholder
}

func (mc *MetricsCollector) getNodesFailedTotal() int64 {
    // Implementation to get total nodes failed
    return 0 // Placeholder
}

func (mc *MetricsCollector) getHeartbeatLatencyHistogram() *Histogram {
    // Implementation to get heartbeat latency histogram
    return &Histogram{
        Buckets: make(map[float64]int64),
        Count:   0,
        Sum:     0,
    }
}
```

## 13. Troubleshooting Guide

### Common Issues

#### Issue 1: USB Detection Failure
**Symptoms**: "No USB devices found" error
**Causes**: 
- Insufficient permissions
- USB device not properly connected
- Platform-specific detection issues

**Solutions**:
```bash
# Linux: Add user to disk group
sudo usermod -a -G disk $USER
newgrp disk

# Windows: Run as Administrator
# Right-click PowerShell and select "Run as Administrator"

# macOS: Grant Full Disk Access
# System Preferences > Security & Privacy > Privacy > Full Disk Access
```

#### Issue 2: Handshake Timeout
**Symptoms**: Node created but never registers
**Causes**:
- Firewall blocking port 51000
- Network connectivity issues
- Token validation failures

**Solutions**:
```bash
# Check firewall settings
sudo ufw status
sudo ufw allow 51000

# Check network connectivity
ping command-station-ip
telnet command-station-ip 51000

# Verify token
syntropy setup status
```

#### Issue 3: Cloud-Init Generation Failure
**Symptoms**: "Cloud-init validation failed" error
**Causes**:
- Invalid YAML template
- Missing template variables
- Template execution errors

**Solutions**:
```bash
# Validate YAML syntax
yamllint user-data.yaml

# Check template variables
syntropy node create --debug

# Verify template files
ls -la /opt/syntropy/templates/
```

### Debug Mode
```go
// Debug mode implementation
func (nm *NodeManager) CreateNodeWithDebug(options *CreateOptions) (*CreateResult, error) {
    // Enable debug logging
    nm.logger.SetLevel("debug")
    
    // Log detailed information
    nm.logger.Debug("Starting node creation",
        "options", options,
        "timestamp", time.Now(),
    )
    
    // Create node with detailed logging
    result, err := nm.createSubcomponent.OrchestrateCreation(options)
    if err != nil {
        nm.logger.Error("Node creation failed",
            "error", err,
            "options", options,
        )
        return nil, err
    }
    
    nm.logger.Debug("Node creation successful",
        "result", result,
        "timestamp", time.Now(),
    )
    
    return result, nil
}
```

### Log Analysis
```go
// Log analysis tools
func (nm *NodeManager) AnalyzeLogs(nodeID string) (*LogAnalysis, error) {
    analysis := &LogAnalysis{
        NodeID:    nodeID,
        Timestamp: time.Now(),
        Issues:    make([]LogIssue, 0),
        Metrics:   make(map[string]float64),
    }
    
    // Get logs for node
    logs, err := nm.GetNodeLogs(nodeID, &LogOptions{Lines: 1000})
    if err != nil {
        return nil, err
    }
    
    // Analyze logs for issues
    for _, logEntry := range logs.Logs {
        if strings.Contains(logEntry, "ERROR") {
            analysis.Issues = append(analysis.Issues, LogIssue{
                Level:   "error",
                Message: logEntry,
                Count:   1,
            })
        }
        
        if strings.Contains(logEntry, "WARN") {
            analysis.Issues = append(analysis.Issues, LogIssue{
                Level:   "warning",
                Message: logEntry,
                Count:   1,
            })
        }
    }
    
    // Calculate metrics
    analysis.Metrics["error_count"] = float64(len(analysis.Issues))
    analysis.Metrics["log_entries"] = float64(len(logs.Logs))
    
    return analysis, nil
}
```

## 14. Future Enhancements

### Planned Features

#### Feature 1: Multi-OS Support
**Description**: Support for multiple operating systems beyond Ubuntu Server
**Implementation**: Extend cloud-init generator with OS-specific templates
**Timeline**: Q2 2025

```go
// Multi-OS cloud-init generator
type MultiOSCloudInitGenerator struct {
    templates map[string]*Template
    logger    *Logger
}

func (mcg *MultiOSCloudInitGenerator) GenerateCloudInit(os string, config *NodeConfig) (*CloudInitConfig, error) {
    template, exists := mcg.templates[os]
    if !exists {
        return nil, fmt.Errorf("unsupported OS: %s", os)
    }
    
    // Generate OS-specific cloud-init
    return mcg.generateOSSpecificCloudInit(template, config)
}
```

#### Feature 2: Container-Based Nodes
**Description**: Support for container-based node deployment
**Implementation**: Extend node creation to support container images
**Timeline**: Q3 2025

```go
// Container-based node creation
type ContainerNodeCreator struct {
    containerRuntime ContainerRuntime
    imageManager     ImageManager
    logger           *Logger
}

func (cnc *ContainerNodeCreator) CreateContainerNode(options *ContainerCreateOptions) (*CreateResult, error) {
    // Pull container image
    if err := cnc.imageManager.PullImage(options.Image); err != nil {
        return nil, err
    }
    
    // Create container
    container, err := cnc.containerRuntime.CreateContainer(options)
    if err != nil {
        return nil, err
    }
    
    return &CreateResult{
        NodeID:    container.ID,
        Status:    "created",
        Message:   "Container node created successfully",
        CreatedAt: time.Now(),
    }, nil
}
```

#### Feature 3: Advanced Monitoring
**Description**: Enhanced monitoring and observability
**Implementation**: Integrate with Prometheus and Grafana
**Timeline**: Q4 2025

```go
// Advanced monitoring
type AdvancedMonitor struct {
    prometheusClient *prometheus.Client
    grafanaClient    *grafana.Client
    logger           *Logger
}

func (am *AdvancedMonitor) SetupMonitoring(nodeID string) error {
    // Setup Prometheus metrics
    if err := am.setupPrometheusMetrics(nodeID); err != nil {
        return err
    }
    
    // Setup Grafana dashboards
    if err := am.setupGrafanaDashboards(nodeID); err != nil {
        return err
    }
    
    return nil
}
```

### Research Areas
- **Edge Computing**: Optimizing node provisioning for edge environments
- **AI-Driven Automation**: Using AI to optimize node placement and configuration
- **Quantum-Safe Cryptography**: Preparing for post-quantum security threats
- **Green Computing**: Optimizing for energy efficiency

## 15. Contributing Guidelines

### Development Workflow
1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/new-feature`
3. **Make changes**: Follow coding standards and add tests
4. **Run tests**: Ensure all tests pass and coverage is 100%
5. **Submit pull request**: Include description and test results

### Coding Standards
```go
// Code formatting
// Use gofmt and goimports
// Follow Go naming conventions
// Add comprehensive comments

// Example of well-formatted code
// NodeManager handles the creation and management of nodes in the Syntropy Grid.
// It orchestrates the node creation workflow and manages the registration process.
type NodeManager struct {
    // createSubcomponent handles the node creation workflow
    createSubcomponent *CreateSubcomponent
    
    // registrationSubcomponent handles node registration and monitoring
    registrationSubcomponent *RegistrationSubcomponent
    
    // eventBus handles event communication between components
    eventBus *EventBus
    
    // logger provides structured logging
    logger *Logger
    
    // mutex ensures thread-safe operations
    mutex sync.RWMutex
}

// CreateNode creates a new node with the specified options.
// It returns a CreateResult containing the node information or an error if creation fails.
func (nm *NodeManager) CreateNode(options *CreateOptions) (*CreateResult, error) {
    // Validate input
    if err := nm.validateCreateOptions(options); err != nil {
        return nil, fmt.Errorf("invalid options: %w", err)
    }
    
    // Create node
    result, err := nm.createSubcomponent.OrchestrateCreation(options)
    if err != nil {
        return nil, fmt.Errorf("node creation failed: %w", err)
    }
    
    // Publish event
    nm.eventBus.Publish(Event{
        Type:      "node_created",
        NodeID:    result.NodeID,
        Data:      map[string]interface{}{"result": result},
        Timestamp: time.Now(),
    })
    
    return result, nil
}
```

### Testing Requirements
- **100% test coverage** for all source files
- **Unit tests** for all public functions
- **Integration tests** for component interactions
- **E2E tests** for complete workflows
- **Performance tests** for critical paths

### Documentation Requirements
- **API documentation** for all public interfaces
- **Code comments** for complex algorithms
- **README updates** for new features
- **Example code** for common use cases

## 16. Conclusion

### Summary
The Node Component provides a comprehensive solution for automated node provisioning in the Syntropy Cooperative Grid. It implements a plug-and-play system that eliminates manual configuration and ensures secure, reliable node deployment.

### Key Achievements
- **Complete Automation**: Zero-touch node creation and registration
- **Multi-Platform Support**: Windows, Linux, and macOS compatibility
- **Security**: Three-layer authentication system
- **Reliability**: Comprehensive error handling and recovery
- **Testability**: 100% test coverage with comprehensive test suite

### Next Steps
1. **Implementation**: Follow this guide to implement the component
2. **Testing**: Ensure all tests pass and coverage is 100%
3. **Documentation**: Complete all documentation files
4. **Integration**: Integrate with CLI and other components
5. **Deployment**: Deploy and monitor in production environment

### Success Criteria
- ✅ `syntropy node create` creates USB plug-and-play automatically
- ✅ Cloud-init installs Syntropy CLI on node
- ✅ Node executes handshake automatically on boot
- ✅ Registration is accepted by Command Station
- ✅ Node is ready for workloads without intervention
- ✅ Heartbeat maintains active connection
- ✅ Tests achieve 100% coverage
- ✅ Documentation is complete and accurate

---

**Implementation Status**: Ready for development
**Estimated Development Time**: 3-4 weeks
**Team Size**: 2-3 developers
**Priority**: High (Critical path component)
