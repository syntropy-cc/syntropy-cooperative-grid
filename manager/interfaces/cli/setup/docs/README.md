# Syntropy CLI Setup Component

## What is This?
[MACRO_VIEW]
The Setup Component is the foundational initialization system for the Syntropy Cooperative Grid CLI, responsible for configuring the entire Syntropy environment on user systems.
[/MACRO_VIEW]

[MESO_VIEW]
This component integrates with the Manager's CLI interface and API services to provide seamless environment setup, configuration management, and system validation across multiple operating systems.
[/MESO_VIEW]

[MICRO_VIEW]
The Setup Component handles environment validation, directory structure creation, cryptographic key generation, configuration file management, and optional system service installation.
[/MICRO_VIEW]

## Why Use This?
### Problems It Solves
- **Complex Environment Setup**: Automates the intricate process of configuring Syntropy CLI across different operating systems
- **Configuration Management**: Generates and manages YAML configuration files with proper directory structures
- **Security Initialization**: Creates Ed25519 cryptographic key pairs for secure network participation
- **System Integration**: Optionally installs system services for background operation
- **Environment Validation**: Ensures system compatibility and resource availability before setup

### Key Benefits
- **Cross-Platform Support**: Works seamlessly on Windows, Linux, and macOS
- **API-First Architecture**: Integrates with centralized API services with local fallback
- **Secure by Default**: Generates cryptographic keys and follows security best practices
- **Idempotent Operations**: Safe to run multiple times without side effects
- **Comprehensive Validation**: Checks system requirements, permissions, and dependencies

## Quick Start
### Prerequisites
- Go 1.19 or higher
- Administrative privileges (recommended for service installation)
- Internet connectivity for API integration
- Minimum 1GB available disk space

### Installation
```bash
# Clone the repository
git clone https://github.com/syntropy-cc/syntropy-cooperative-grid.git
cd syntropy-cooperative-grid/manager/interfaces/cli/setup

# Build the setup component
go build -o syntropy-setup .
```

### Basic Usage
```go
package main

import (
    "fmt"
    "github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup"
    "github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/internal/types"
)

func main() {
    // Configure setup options
    options := types.SetupOptions{
        Force:          false,
        InstallService: true,
        ConfigPath:     "", // Use default path
        HomeDir:        "", // Use default home directory
    }
    
    // Execute setup
    result, err := setup.Setup(options)
    if err != nil {
        fmt.Printf("Setup failed: %v\n", err)
        return
    }
    
    fmt.Printf("Setup completed successfully!\n")
    fmt.Printf("Configuration file: %s\n", result.ConfigPath)
    fmt.Printf("Duration: %s\n", result.EndTime.Sub(result.StartTime))
}
```

**Input**: SetupOptions with configuration preferences
**Output**: SetupResult with success status, configuration path, and timing information

## Features
| Feature | Description | Status |
|---------|-------------|--------|
| Cross-Platform Setup | Windows, Linux, macOS support | Stable |
| API Integration | Centralized setup via API services | Stable |
| Environment Validation | System compatibility checks | Stable |
| Key Generation | Ed25519 cryptographic key pairs | Stable |
| Service Installation | Optional system service setup | Stable |
| Configuration Management | YAML config file generation | Stable |
| Status Checking | Installation status verification | Stable |
| Reset Functionality | Clean uninstallation | Stable |

## Component Structure
```
setup/
├── docs/                    # Documentation
│   ├── README.md           # This file
│   ├── DEV.md             # Developer documentation
│   ├── API.md             # API reference
│   └── TEST.md            # Testing guide
├── config/                 # Configuration templates
│   ├── defaults/          # Default configurations
│   ├── schemas/           # Configuration schemas
│   └── templates/         # YAML templates
├── internal/              # Internal packages
│   ├── services/          # Service implementations
│   ├── types/             # Type definitions
│   └── utils/             # Utility functions
├── tests/                 # Test files
│   ├── fixtures/          # Test data
│   ├── integration/       # Integration tests
│   └── unit/              # Unit tests
├── setup.go               # Main setup logic
├── setup_linux.go         # Linux-specific implementation
├── setup_windows.go       # Windows-specific implementation
├── api_integration.go     # API service integration
├── validation_linux.go    # Linux environment validation
├── validation_windows.go  # Windows environment validation
├── configuration_linux.go # Linux configuration
└── configuration_windows.go # Windows configuration
```

## Next Steps
- [Explore the API](./API.md) - Detailed usage instructions and function reference
- [Developer Guide](./DEV.md) - Understanding the internals and architecture
- [Testing Guide](./TEST.md) - Running and writing tests
- [Examples](../examples/) - More usage scenarios and integration patterns

## Support
- Issue Tracker: [GitHub Issues](https://github.com/syntropy-cc/syntropy-cooperative-grid/issues)
- Discussion Forum: [GitHub Discussions](https://github.com/syntropy-cc/syntropy-cooperative-grid/discussions)
- Contact: [Syntropy Community](https://syntropy.com/community)

## License
Apache 2.0 - See LICENSE file for details