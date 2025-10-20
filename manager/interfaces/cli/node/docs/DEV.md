# Node Component - Developer Documentation

## Architecture Overview

[MACRO_VIEW]
The Node Component implements an Event-Driven Architecture pattern within the Syntropy Cooperative Grid system, serving as the automated provisioning layer that bridges physical hardware deployment with network registration, enabling the larger grid orchestration system to function.
[/MACRO_VIEW]

[MESO_VIEW]
This component integrates with the Syntropy Manager CLI module through well-defined interfaces, consuming Grid Tokens from the Setup Component and preparing nodes for Workload Component orchestration, maintaining loose coupling through event-based communication patterns.
[/MESO_VIEW]

[MICRO_VIEW]
The internal architecture follows a two-subcomponent pattern (Create + Registration) with shared state management, implementing thread-safe operations for concurrent node management and providing plug-and-play automation through cloud-init integration.
[/MICRO_VIEW]

## Design Decisions

### Architectural Pattern
**Pattern Used**: Event-Driven Architecture with Command Pattern
**Justification**: Enables loose coupling between node creation and registration phases, supports concurrent operations, and provides clear separation of concerns
**Trade-offs**: 
- Pros: Scalable, maintainable, testable, supports async operations
- Cons: Increased complexity, requires careful event ordering, debugging can be challenging

### Core Abstractions

| Abstraction | Purpose | Design Principle |
|-------------|---------|------------------|
| NodeManager | Central orchestration and state management | Single Responsibility Principle |
| CreateSubcomponent | Encapsulates all node creation logic | Interface Segregation Principle |
| RegistrationSubcomponent | Handles registration and lifecycle management | Dependency Inversion Principle |
| USBDetector | Abstracts platform-specific USB detection | Open/Closed Principle |
| CloudInitGenerator | Generates platform-agnostic installation scripts | Liskov Substitution Principle |

## Component Internals

### Directory Structure Deep Dive
```
src/
├── node.go                  # Main orchestrator (500 lines)
│   ├── NodeManager struct   # Central state management
│   ├── CreateNode()         # Public API for node creation
│   └── Event handling       # Coordinates subcomponents
├── create.go                # Create subcomponent (300 lines)
│   ├── CreateOrchestrator   # Coordinates creation workflow
│   ├── USB detection        # Platform abstraction
│   └── Cloud-init generation # Installation automation
├── registration.go          # Registration subcomponent (300 lines)
│   ├── RegistrationManager  # Handles registration lifecycle
│   ├── Handshake protocol   # Secure authentication
│   └── Heartbeat management # Connection monitoring
├── usb_detector*.go         # Platform-specific USB detection
├── cloud_init_generator.go  # Installation script generation
├── handshake.go             # Security protocol implementation
├── listener.go              # TCP listener for node connections
├── heartbeat.go             # Connection health monitoring
├── node_state.go            # State management and persistence
└── internal/                # Private implementation details
    ├── types/               # Internal type definitions
    ├── helpers/             # Utility functions
    └── constants/           # System constants
```

### Core Components

#### Component: NodeManager
##### Responsibility
Central orchestration of node creation and registration lifecycle, maintaining thread-safe state for all nodes (pending, active, inactive).

##### Collaborators
- CreateSubcomponent: Initiates node creation workflow
- RegistrationSubcomponent: Handles registration and monitoring
- TokenIntegration: Validates Grid Tokens from Setup Component
- NodeState: Manages persistent state across operations

##### Key Algorithms
| Algorithm | Complexity | Use Case |
|-----------|------------|----------|
| Node Creation | Time: O(n), Space: O(1) | Creating new nodes with USB detection |
| State Transitions | Time: O(1), Space: O(1) | Moving nodes between pending/active/inactive |
| Concurrent Access | Time: O(1), Space: O(1) | Thread-safe operations with RWMutex |

##### State Management
```
Initial State (No Nodes)
    ↓ CreateNode() triggered
Pending State (USB Created, Waiting for Connection)
    ↓ Handshake successful
Active State (Connected, Heartbeat Active)
    ↓ Connection lost
Inactive State (Disconnected, Requires Intervention)
```

#### Component: CloudInitGenerator
##### Responsibility
Generates complete cloud-init configuration for Ubuntu Server installation, including Syntropy CLI installation and automatic registration scripts.

##### Collaborators
- AutoConfigGenerator: Provides NodeID, tokens, and certificates
- TemplateEngine: Processes cloud-init templates with variables
- SecurityValidator: Ensures generated configurations are secure

##### Key Algorithms
| Algorithm | Complexity | Use Case |
|-----------|------------|----------|
| Template Processing | Time: O(m), Space: O(n) | Where m=template size, n=variables |
| YAML Validation | Time: O(n), Space: O(n) | Validating generated cloud-init |
| Security Sanitization | Time: O(k), Space: O(1) | Where k=security rules |

### Data Flow Architecture
```
User Command (syntropy node create)
    ↓
NodeManager.CreateNode()
    ↓
CreateSubcomponent.OrchestrateCreation()
    ├── USBDetector.DetectAvailable()
    ├── AutoConfigGenerator.GenerateConfig()
    ├── CloudInitGenerator.GenerateScripts()
    └── USBWriter.CreateBootableUSB()
    ↓
RegistrationSubcomponent.StartListener()
    ↓
Node connects and sends NodeAnnouncement
    ↓
Handshake.ValidateAndAccept()
    ↓
NodeState.TransitionToActive()
    ↓
Heartbeat.StartMonitoring()
```

### Dependency Graph
```
NodeManager
├── depends on → CreateSubcomponent
├── depends on → RegistrationSubcomponent
└── depends on → TokenIntegration
    └── depends on → Setup Component (TokenManager)

CreateSubcomponent
├── depends on → USBDetector (platform-specific)
├── depends on → AutoConfigGenerator
├── depends on → CloudInitGenerator
└── depends on → USBWriter (platform-specific)

RegistrationSubcomponent
├── depends on → Handshake
├── depends on → Listener
└── depends on → Heartbeat
```

## Extension Points

### How to Add New Features
1. **Identify Extension Point**
   - Look for interfaces in src/internal/types/interfaces.go
   - Consider if feature belongs in Create or Registration subcomponent
   - Determine if platform-specific implementation needed

2. **Implement Interface/Contract**
   - Extend existing interfaces or create new ones
   - Follow existing patterns for error handling and logging
   - Ensure thread-safe operations for state management

3. **Register Component**
   - Add to appropriate subcomponent orchestrator
   - Update dependency injection in NodeManager
   - Add configuration options if needed

### Plugin Architecture
The component supports platform-specific plugins through build constraints:
- `// +build windows` for Windows-specific implementations
- `// +build linux` for Linux-specific implementations
- `// +build darwin` for macOS-specific implementations

## Performance Characteristics

### Resource Usage
| Resource | Typical Usage | Maximum Usage | Scaling Factor |
|----------|--------------|---------------|----------------|
| Memory | 50MB base + 10MB per active node | 200MB (20 nodes) | O(n) where n=active nodes |
| CPU | 5% during creation, 1% idle | 15% during concurrent operations | O(k) where k=concurrent operations |
| I/O | 1GB during ISO download/cache | 10GB for multiple node creation | O(m) where m=number of nodes created |

### Optimization Strategies
1. **Concurrent USB Operations**
   - Implementation: Parallel USB detection and writing
   - Impact: 3x faster for multiple node creation
   - Trade-off: Increased memory usage

2. **ISO Caching**
   - Implementation: Local cache with SHA256 verification
   - Impact: 90% faster for subsequent node creation
   - Trade-off: Disk space usage (~4GB per cached ISO)

3. **Connection Pooling**
   - Implementation: Reuse TCP connections for heartbeat
   - Impact: 50% reduction in connection overhead
   - Trade-off: Increased memory per active connection

## Security Considerations

### Threat Model
| Threat | Mitigation | Residual Risk |
|--------|------------|---------------|
| Unauthorized Node Registration | Grid Token validation, Node Certificate verification | Low - Multi-layer authentication |
| USB-based Attacks | Validation of USB devices, prevention of system disk writing | Low - Strict validation rules |
| Network Interception | TLS encryption for handshake, SSH keys for communication | Low - Encrypted communication |
| Token Theft | Keyring integration, no token storage in files | Low - System-level security |

### Security Boundaries
```
[Trusted Zone: Command Station]
    ↓ Secure Handshake (TLS + Grid Token)
[Security Boundary: Network]
    ↓ Node Certificate + SSH Keys
[Trusted Zone: Registered Node]
```

## Development Workflow

### Setting Up Development Environment
```bash
# Step 1: Ensure Setup Component is configured
cd manager/interfaces/cli/setup
go test ./...

# Step 2: Install Node Component dependencies
cd ../node
go mod tidy
go mod download

# Step 3: Verify setup
go test ./tests/unit/... -v
```

### Code Organization Principles
- **Separation of Concerns**: Create and Registration subcomponents are independent
- **Dependency Injection**: All dependencies injected through constructors
- **Error Handling**: Structured errors with context and recovery suggestions

### Debugging Techniques
| Scenario | Technique | Tools |
|----------|-----------|-------|
| USB Detection Issues | Platform-specific logging, device enumeration | lsblk (Linux), WMIC (Windows) |
| Handshake Failures | Network packet capture, token validation logging | Wireshark, tcpdump |
| State Management Bugs | State transition logging, race condition detection | go race detector, custom state logger |

## Monitoring and Observability

### Key Metrics
| Metric | Purpose | Alert Threshold |
|--------|---------|-----------------|
| Node Creation Success Rate | Monitor provisioning reliability | < 95% success rate |
| Handshake Latency | Monitor network performance | > 30 seconds |
| Active Node Count | Monitor grid capacity | Based on expected capacity |
| Heartbeat Failure Rate | Monitor node health | > 5% failure rate |

### Logging Strategy
- **Debug Level**: Detailed operation flow, USB detection details
- **Info Level**: Node state transitions, successful operations
- **Error Level**: Failures, security violations, system errors

### Debugging Hooks
```bash
# Enable verbose debugging
export SYNTHROPY_NODE_DEBUG=1
export SYNTHROPY_NODE_LOG_LEVEL=debug

# Run with debugging enabled
syntropy node create --debug
```

## Maintenance Guidelines

### Code Health Metrics
- Cyclomatic Complexity: Maximum 10 per function
- Coupling: Maximum 3 dependencies per component
- Cohesion: Minimum 80% related functionality per file

### Refactoring Triggers
1. Function exceeds 50 lines → Extract to helper function
2. File exceeds 500 lines → Split into logical subcomponents
3. Test coverage drops below 100% → Add missing test cases

## Migration Guide

### Breaking Changes Policy
Breaking changes are avoided through interface versioning. When breaking changes are necessary, they are introduced with deprecation warnings and migration paths.

### Version Compatibility Matrix
| Component Version | Compatible With | Migration Required |
|-------------------|-----------------|-------------------|
| 2.x | Setup Component 3.x | No |
| 1.x | Setup Component 2.x | Yes - Token format changed |

## Troubleshooting Development Issues

### Common Problems
| Symptom | Likely Cause | Solution |
|---------|--------------|----------|
| USB detection fails | Insufficient permissions | Run with sudo or add user to disk group |
| Handshake timeout | Firewall blocking port 51000 | Configure firewall to allow TCP:51000 |
| Cloud-init generation fails | Invalid YAML template | Validate template syntax, check variable substitution |

## Contributing

### Code Review Checklist
- [ ] Follows architectural patterns (Event-Driven, Command)
- [ ] Maintains abstraction boundaries between subcomponents
- [ ] Includes performance impact analysis for new features
- [ ] Updates relevant documentation (DEV.md, API.md)
- [ ] Maintains 100% test coverage

## Further Learning
For theoretical foundations and pedagogical insights, see [LEARN.md](./LEARN.md).
