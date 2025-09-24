# Syntropy CLI Setup Component

This component provides functionality for setting up the Syntropy CLI environment across different operating systems.

## Supported Platforms

- Windows
- Linux (Ubuntu, Debian, CentOS, Fedora, RHEL)
- macOS (planned)

## Features

- Environment validation
- Configuration management
- Service installation
- Cross-platform compatibility

## Usage

### Setup

```go
import "github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup"

options := types.SetupOptions{
    Force:          false,
    InstallService: true,
    ConfigPath:     "/custom/path/config.yaml", // Optional
    HomeDir:        "/custom/home",             // Optional
}

result, err := setup.Setup(options)
if err != nil {
    log.Fatalf("Setup failed: %v", err)
}

fmt.Printf("Setup successful: %v\n", result.Success)
fmt.Printf("Config path: %s\n", result.ConfigPath)
```

### Status Check

```go
result, err := setup.Status(options)
if err != nil {
    log.Fatalf("Status check failed: %v", err)
}

fmt.Printf("Status: %v\n", result.Success)
```

### Reset

```go
result, err := setup.Reset(options)
if err != nil {
    log.Fatalf("Reset failed: %v", err)
}

fmt.Printf("Reset successful: %v\n", result.Success)
```

## Platform-Specific Details

### Windows

- Validates Windows version, administrator rights, disk space, PowerShell version
- Installs Windows service using PowerShell
- Configuration stored in `%USERPROFILE%\Syntropy`

### Linux

- Validates Linux distribution, root permissions, disk space, system dependencies
- Installs systemd service (when requested)
- Configuration stored in `$HOME/.syntropy`
- Supports common Linux distributions (Ubuntu, Debian, CentOS, Fedora, RHEL)

#### Linux Service Installation

When the `InstallService` option is set to `true`, the setup will:

1. Create a systemd service file in `$HOME/.syntropy/services`
2. Generate an installation script
3. Provide instructions for installing the service with root privileges

To install the service after setup:

```bash
sudo $HOME/.syntropy/services/install_service.sh
```

## Development

### Adding Support for New Platforms

To add support for a new platform:

1. Create platform-specific validation file (e.g., `validation_darwin.go`)
2. Create platform-specific setup file (e.g., `setup_darwin.go`)
3. Create platform-specific configuration file (e.g., `configuration_darwin.go`)
4. Update the main `setup.go` file to include the new platform
5. Add unit tests for the new platform

### Running Tests

```bash
go test -v ./...
```

## License

This component is part of the Syntropy Cooperative Grid project and is subject to the project's license terms.