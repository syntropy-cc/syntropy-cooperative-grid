# Node Component

## What is This?

[MACRO_VIEW]
The Node Component is the automated provisioning and registration system for physical nodes in the Syntropy Cooperative Grid network, enabling plug-and-play deployment of compute resources.
[/MACRO_VIEW]

[MESO_VIEW]
This component integrates with the Syntropy Manager CLI module to provide seamless node creation, working alongside the Setup Component for token management and preparing for the Workload Component's orchestration needs.
[/MESO_VIEW]

[MICRO_VIEW]
The specific problem this component solves is the manual, error-prone process of creating and registering physical nodes, replacing it with a fully automated system that requires zero user intervention after initial command execution.
[/MICRO_VIEW]

## Why Use This?

### Problems It Solves
- **Manual Node Provisioning**: Eliminates the need for manual USB creation, ISO downloads, and configuration
- **Registration Complexity**: Automates the complex handshake and registration process between nodes and command station
- **Multi-Platform Inconsistency**: Provides consistent node creation experience across Windows, Linux, and macOS
- **Security Configuration**: Automatically generates and manages cryptographic keys, certificates, and tokens
- **Error-Prone Manual Steps**: Removes human error from node deployment process

### Key Benefits
- **Zero-Touch Deployment**: Create nodes with a single command, everything else is automatic
- **Plug-and-Play Experience**: Insert USB into hardware, node automatically registers and becomes ready
- **Enterprise Security**: Three-layer authentication system with Grid Token, Node Certificates, and SSH keys
- **Multi-Platform Support**: Works consistently across all major operating systems
- **Production Ready**: Complete systemd integration, heartbeat monitoring, and error recovery

## Quick Start

### Prerequisites
- Setup Component configured (100% complete)
- Grid Token available via system Keyring
- USB device (minimum 8GB capacity)
- Network connectivity for ISO download
- Ubuntu Server 24.04 LTS (downloaded automatically)

### Installation
```bash
# Node Component is part of Syntropy Manager CLI
# Ensure Setup Component is configured first
syntropy setup status
```

### Basic Usage
```bash
# Create a new node (fully automated)
syntropy node create

# Input: No parameters required - everything is automatic
# Output: USB bootable device ready for hardware deployment

# List all nodes (created, pending, active)
syntropy node list

# Input: No parameters
# Output: Table showing all nodes with status

# Check status of specific node
syntropy node status node-01

# Input: Node ID
# Output: Detailed status including IP, uptime, health
```

## Features

| Feature | Description | Status |
|---------|-------------|--------|
| USB Auto-Detection | Automatically detects and validates USB devices | Stable |
| ISO Download & Cache | Downloads Ubuntu Server with SHA256 verification | Stable |
| Auto Config Generation | Generates NodeID, SSH keys, certificates automatically | Stable |
| Cloud-Init Integration | Creates complete installation scripts for nodes | Stable |
| Multi-Platform USB Writing | Cross-platform USB creation (Windows/Linux/macOS) | Stable |
| Automatic Registration | Handshake and registration without user intervention | Stable |
| Heartbeat Monitoring | Continuous health monitoring of registered nodes | Stable |
| Security Validation | Prevents writing to system disks, validates tokens | Stable |
| Multiple Node Support | Manages multiple pending and active nodes simultaneously | Stable |

## Component Structure
```
node/
├── docs/           # Complete documentation
├── src/            # Source code
├── tests/          # Test files (100% coverage)
├── examples/       # Usage examples
└── config/         # Configuration files
```

## Next Steps
- [Explore the API](./API.md) - Detailed usage instructions
- [Developer Guide](./DEV.md) - Understanding the internals
- [Testing Guide](./TEST.md) - Running and writing tests
- [Learning Path](./LEARN.md) - Deep dive into concepts and theory
- [Examples](../examples/) - More usage scenarios
- [Implementation Guide](../GUIDE.md) - For LLM implementation

## Support
- Issue Tracker: [Link to be added]
- Discussion Forum: [Link to be added]
- Contact: [Contact method to be added]

## License
MIT License - See LICENSE file for details
