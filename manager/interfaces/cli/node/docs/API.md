# Node Component - API Documentation

## API Overview

[MACRO_VIEW]
The Node Component's API philosophy centers around simplicity and automation, providing a minimal interface that hides complexity while enabling powerful plug-and-play node provisioning within the Syntropy Cooperative Grid ecosystem.
[/MACRO_VIEW]

[MESO_VIEW]
This API integrates seamlessly with the Syntropy Manager CLI module's command structure, providing consistent patterns with other components while maintaining clear boundaries for node-specific operations.
[/MESO_VIEW]

[MICRO_VIEW]
The specific capabilities exposed through this API include automated node creation, lifecycle management, status monitoring, and configuration management, all designed for zero-touch operation.
[/MICRO_VIEW]

## API Principles
- **Simplicity**: Single command creates complete nodes with minimal parameters
- **Idempotency**: Safe to retry operations without side effects
- **Consistency**: All commands follow verb-noun pattern with consistent error handling
- **Automation**: Default behavior requires no user interaction
- **Security**: All operations validate tokens and certificates automatically

## Authentication & Authorization

### Authentication Methods
| Method | Use Case | Example |
|--------|----------|---------|
| Grid Token | Initial node registration | Validated against Setup Component Keyring |
| Node Certificate | Ongoing node communication | Ed25519 certificate validation |
| SSH Keys | Command execution on nodes | RSA 2048-bit key authentication |

### Required Permissions
| Endpoint/Method | Permission Level | Scope |
|-----------------|------------------|--------|
| CreateNode | User | Local system USB access, network connectivity |
| ListNodes | User | Read access to node state directory |
| GetNodeStatus | User | Read access to specific node state |
| GetNodeLogs | User | Read access to node log files |

## Core API Reference

### Node Creation Operations

#### Method: CreateNode
**Purpose**: Creates a new node with automatic USB detection and configuration

**Signature**:
```go
func (nm *NodeManager) CreateNode(options *CreateOptions) (*CreateResult, error)
```

**Parameters**:
| Parameter | Type | Required | Default | Description | Constraints |
|-----------|------|----------|---------|-------------|-------------|
| options.USBPath | string | No | "" | Specific USB device path | Must be valid USB device |
| options.ForceDownload | bool | No | false | Force ISO re-download | N/A |
| options.ISOPath | string | No | "" | Custom ISO file path | Must be valid Ubuntu Server ISO |

**Returns**:
```go
type CreateResult struct {
    NodeID     string    // Generated node ID (e.g., "node-01")
    USBPath    string    // USB device used for creation
    Status     string    // Creation status ("created", "failed")
    Message    string    // Human-readable result message
    CreatedAt  time.Time // Timestamp of creation
    Config     *NodeConfig // Generated configuration
}
```

**Errors**:
| Error Code | Condition | Resolution |
|------------|-----------|------------|
| ErrNoUSBFound | No USB devices detected | Insert USB device and retry |
| ErrUSBTooSmall | USB device smaller than 8GB | Use USB device with 8GB+ capacity |
| ErrInvalidISO | ISO file corrupted or invalid | Re-download ISO or use different file |
| ErrTokenValidation | Grid Token validation failed | Ensure Setup Component is configured |

**Example**:
```go
// Request
options := &CreateOptions{
    USBPath: "/dev/sdb",
}
result, err := nodeManager.CreateNode(options)

// Response - Success
{
    "node_id": "node-01",
    "usb_path": "/dev/sdb",
    "status": "created",
    "message": "Node created successfully. USB ready for deployment.",
    "created_at": "2025-01-27T10:30:00Z",
    "config": {
        "node_id": "node-01",
        "grid_token": "***",
        "ssh_public_key": "ssh-rsa AAAAB3N...",
        "command_station_ip": "192.168.1.100"
    }
}

// Response - Error
{
    "error": "no_usb_found",
    "message": "No USB devices detected. Please insert a USB device and try again.",
    "suggestion": "Ensure USB device is properly connected and has at least 8GB capacity"
}
```

**Notes**:
- Operation is fully automated - no user interaction required
- USB device will be completely overwritten
- Generated USB contains complete Ubuntu Server installation

#### Method: CreateNodeWithAutoDetection
**Purpose**: Creates a node using automatic USB detection (default behavior)

**Signature**:
```go
func (nm *NodeManager) CreateNodeWithAutoDetection() (*CreateResult, error)
```

**Parameters**: None (fully automatic)

**Returns**: Same as CreateNode

**Example**:
```bash
# CLI equivalent
syntropy node create
```

### Node Management Operations

#### Method: ListNodes
**Purpose**: Lists all nodes with their current status

**Signature**:
```go
func (nm *NodeManager) ListNodes() (*NodeList, error)
```

**Parameters**: None

**Returns**:
```go
type NodeList struct {
    Active   []*NodeStatus // Currently active nodes
    Pending  []*NodeStatus // Nodes waiting for connection
    Inactive []*NodeStatus // Disconnected nodes
    Total    int           // Total number of nodes
}
```

**Example**:
```go
// Response
{
    "active": [
        {
            "node_id": "node-01",
            "status": "active",
            "ip_address": "192.168.1.101",
            "uptime": "2h30m15s",
            "last_heartbeat": "2025-01-27T10:30:00Z"
        }
    ],
    "pending": [
        {
            "node_id": "node-02",
            "status": "pending",
            "time_left": "25m30s",
            "created_at": "2025-01-27T10:00:00Z"
        }
    ],
    "inactive": [],
    "total": 2
}
```

#### Method: GetNodeStatus
**Purpose**: Gets detailed status of a specific node

**Signature**:
```go
func (nm *NodeManager) GetNodeStatus(nodeID string) (*NodeStatus, error)
```

**Parameters**:
| Parameter | Type | Required | Default | Description | Constraints |
|-----------|------|----------|---------|-------------|-------------|
| nodeID | string | Yes | N/A | Node identifier | Must be valid node ID format |

**Returns**:
```go
type NodeStatus struct {
    NodeID        string    // Node identifier
    Status        string    // Current status (active, pending, inactive)
    IPAddress     string    // Node IP address (if active)
    Uptime        time.Duration // How long node has been active
    LastHeartbeat time.Time // Last successful heartbeat
    Hardware      *HardwareInfo // Hardware specifications
    CreatedAt     time.Time // When node was created
    RegisteredAt  time.Time // When node registered (if active)
}
```

**Errors**:
| Error Code | Condition | Resolution |
|------------|-----------|------------|
| ErrNodeNotFound | Node ID does not exist | Check node ID or create new node |
| ErrInvalidNodeID | Invalid node ID format | Use format "node-XX" where XX is number |

#### Method: GetNodeLogs
**Purpose**: Retrieves logs for a specific node

**Signature**:
```go
func (nm *NodeManager) GetNodeLogs(nodeID string, options *LogOptions) (*NodeLogs, error)
```

**Parameters**:
| Parameter | Type | Required | Default | Description | Constraints |
|-----------|------|----------|---------|-------------|-------------|
| nodeID | string | Yes | N/A | Node identifier | Must be valid node ID |
| options.Lines | int | No | 100 | Number of log lines | 1-1000 |
| options.Follow | bool | No | false | Follow logs in real-time | N/A |
| options.Service | string | No | "" | Specific service logs | "registration", "heartbeat", "system" |

**Returns**:
```go
type NodeLogs struct {
    NodeID    string   // Node identifier
    Logs      []string // Log entries
    Timestamp time.Time // When logs were retrieved
    Service   string   // Service name (if filtered)
}
```

## Event API / Callbacks

### Event: NodeCreated
**Triggered When**: New node is successfully created

**Payload**:
```json
{
  "eventType": "node_created",
  "timestamp": "2025-01-27T10:30:00Z",
  "data": {
    "node_id": "node-01",
    "usb_path": "/dev/sdb",
    "config": {
      "node_id": "node-01",
      "command_station_ip": "192.168.1.100"
    }
  }
}
```

### Event: NodeRegistered
**Triggered When**: Node successfully completes handshake and registers

**Payload**:
```json
{
  "eventType": "node_registered",
  "timestamp": "2025-01-27T10:35:00Z",
  "data": {
    "node_id": "node-01",
    "ip_address": "192.168.1.101",
    "hardware": {
      "cpu_cores": 8,
      "memory_gb": 16,
      "disk_gb": 500
    }
  }
}
```

### Event: NodeDisconnected
**Triggered When**: Node loses connection or heartbeat fails

**Payload**:
```json
{
  "eventType": "node_disconnected",
  "timestamp": "2025-01-27T10:40:00Z",
  "data": {
    "node_id": "node-01",
    "last_heartbeat": "2025-01-27T10:39:30Z",
    "reason": "heartbeat_timeout"
  }
}
```

**Subscribe Example**:
```go
// Subscribe to node events
eventBus.Subscribe("node_registered", func(event Event) {
    log.Printf("Node %s registered successfully", event.Data.NodeID)
})
```

## Cloud-Init API

### Template Variables
The Cloud-Init generator accepts the following template variables:

| Variable | Type | Description | Example |
|----------|------|-------------|---------|
| .NodeID | string | Unique node identifier | "node-01" |
| .GridToken | string | Grid authentication token | "abc123-def456-..." |
| .CommandStationIP | string | IP address of command station | "192.168.1.100" |
| .SSHPublicKey | string | SSH public key for node access | "ssh-rsa AAAAB3N..." |
| .SSHPrivateKey | string | SSH private key (encrypted) | "-----BEGIN RSA..." |
| .NodeCertificate | string | Node certificate for authentication | "-----BEGIN CERT..." |
| .CreatedAt | time.Time | Node creation timestamp | "2025-01-27T10:30:00Z" |

### Cloud-Init Generation
**Method**: GenerateCloudInit
**Purpose**: Generate complete cloud-init configuration for node

**Signature**:
```go
func (cig *CloudInitGenerator) GenerateCloudInit(config *NodeConfig) (*CloudInitResult, error)
```

**Returns**:
```go
type CloudInitResult struct {
    UserData     string // user-data.yaml content
    NetworkConfig string // network-config.yaml content
    MetaData     string // meta-data.yaml content
    Valid        bool   // YAML validation result
}
```

## Protocol API

### Handshake Protocol
**Endpoint**: TCP:51000
**Purpose**: Secure node registration with command station

#### NodeAnnouncement Message
```json
{
  "type": "node_announcement",
  "node_id": "node-01",
  "grid_token": "abc123-def456-...",
  "node_certificate": "-----BEGIN CERT...",
  "hardware": {
    "cpu_cores": 8,
    "memory_gb": 16,
    "disk_gb": 500,
    "hostname": "node-01",
    "ip_address": "192.168.1.101"
  },
  "timestamp": "2025-01-27T10:30:00Z"
}
```

#### HandshakeResponse Message
```json
{
  "status": "accepted",
  "message": "Node registered successfully",
  "command_station_cert": "-----BEGIN CERT...",
  "ssh_config": {
    "port": 22,
    "key_algorithm": "rsa"
  },
  "workload_config": {
    "docker_enabled": true,
    "resource_limits": {
      "cpu_cores": 8,
      "memory_gb": 16
    }
  }
}
```

## Rate Limiting

| Endpoint Category | Requests/Minute | Burst Limit | Retry-After |
|-------------------|-----------------|-------------|-------------|
| Node Creation | 10 | 3 | 60 seconds |
| Status Queries | 100 | 20 | 10 seconds |
| Log Retrieval | 50 | 10 | 30 seconds |

## API Versioning

### Version Strategy
**Current Version**: v1.0
**Supported Versions**: v1.0
**Deprecation Policy**: 6-month notice for breaking changes

### Version Differences
| Feature | v1.0 | v1.1 (planned) | v2.0 (future) |
|---------|------|----------------|---------------|
| USB Detection | Basic | Enhanced validation | Advanced filtering |
| Cloud-Init | Ubuntu 24.04 | Multi-OS support | Container-based |
| Security | 3-layer auth | Enhanced crypto | Zero-trust |

## Testing the API

### Test Environment
**Base URL**: Local CLI (no HTTP API)
**Test Credentials**: Use Setup Component configuration
**Limitations**: Requires physical USB devices for creation tests

### Example Test Flow
```bash
# 1. Setup test environment
syntropy setup status

# 2. Create test node
syntropy node create --usb /dev/sdb

# 3. Verify node creation
syntropy node list

# 4. Check node status
syntropy node status node-01

# 5. Retrieve logs
syntropy node logs node-01
```

## API Best Practices

### DO
- Always check node status before performing operations
- Use automatic USB detection unless specific device needed
- Monitor events for node lifecycle changes
- Handle errors gracefully with appropriate user feedback

### DON'T
- Don't create multiple nodes simultaneously on same USB
- Don't ignore timeout errors during node creation
- Don't attempt to use nodes before registration completes
- Don't hardcode node IDs in automation scripts

## Troubleshooting

### Common Integration Issues
| Issue | Symptoms | Solution |
|-------|----------|----------|
| USB Detection Failure | "No USB devices found" | Check permissions, ensure USB is connected |
| Handshake Timeout | Node created but never registers | Check firewall, verify network connectivity |
| Token Validation Error | "Invalid Grid Token" | Ensure Setup Component is properly configured |

## API Metrics

### SLA
| Metric | Target | Measurement |
|--------|--------|-------------|
| Availability | 99.9% | Based on CLI command success rate |
| Response Time (p95) | < 2 seconds | Time from command to completion |
| Node Creation Success | > 95% | Successful node creation rate |

## Educational Resources
For conceptual understanding and learning exercises, see [LEARN.md](./LEARN.md).
