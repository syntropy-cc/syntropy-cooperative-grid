# Build and Test Instructions - Syntropy CLI Manager

This document provides step-by-step instructions for building and testing the Syntropy CLI Manager code on Linux and Windows systems.

## 📋 Overview

The **Syntropy CLI Manager** is the main command-line interface for managing the Syntropy Cooperative Grid. It provides a unified interface for all management operations, allowing users to control the network through simple and intuitive commands.

The CLI is located in `manager/interfaces/cli/` and includes:
- **Main entry point**: `main.go` (Cobra-based CLI)
- **Setup component**: `setup/` (first component, with more to come)
- **Build scripts**: Automated build and test scripts
- **Cross-platform support**: Linux, Windows, and macOS

## 🔧 Prerequisites

### Common Dependencies
- **Go 1.22.5+** (as specified in `go.mod`)
- **Git** for version control
- **Make** (optional, but recommended)

### Platform-Specific Dependencies

#### Linux
- `systemd` (for services)
- `systemctl` (for service management)
- Administrator permissions (for service installation)

#### Windows
- **PowerShell 5.1+**
- Administrator permissions
- **Windows Service Control Manager**

## 🏗️ Project Structure

```
manager/interfaces/cli/
├── main.go                     # Main CLI entry point (Cobra)
├── setup/                      # Setup component (first of many)
│   ├── setup.go               # Setup orchestrator
│   ├── setup_linux.go         # Linux implementation
│   ├── setup_windows.go       # Windows implementation
│   ├── validation_linux.go    # Linux validation
│   ├── validation_windows.go  # Windows validation
│   ├── configuration_linux.go # Linux configuration
│   ├── configuration_windows.go # Windows configuration
│   ├── internal/              # Internal types and services
│   ├── tests/                 # Unit and integration tests
│   └── config/                # Configuration templates
├── build.sh                   # Linux/macOS build script
├── build.ps1                  # Windows build script
├── Makefile                   # Make-based build system
├── BUILD_AND_TEST.md          # This document
└── README.md                  # User documentation
```

## 🐧 Building on Linux

### Step 1: Prepare Environment

```bash
# Navigate to the CLI directory
cd /home/jescott/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli

# Verify Go version
go version
# Should show Go 1.22.5 or higher

# Verify we're in the correct directory
pwd
# Should show: /home/jescott/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli
```

### Step 2: Install Dependencies

```bash
# Download Go dependencies
go mod download

# Verify and organize dependencies
go mod tidy

# Verify no dependency issues
go mod verify
```

### Step 3: Build the CLI Manager

#### Option A: Automated Build (Recommended)

```bash
# Make the build script executable
chmod +x build.sh

# Build everything automatically
./build.sh

# Or build for specific platforms
./build.sh current    # Current platform only
./build.sh linux      # Linux only
./build.sh windows    # Windows only (cross-compilation)
./build.sh darwin     # macOS only (cross-compilation)
```

#### Option B: Using Make

```bash
# Build for current platform
make build

# Build for all platforms
make cross-build

# Build with all features
make all

# Run tests
make test

# Clean build artifacts
make clean
```

#### Option C: Manual Build

```bash
# Build for current platform
go build -ldflags "-X main.version=$(date +%Y%m%d-%H%M%S) -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) -X main.gitCommit=$(git rev-parse --short HEAD)" -o build/syntropy main.go

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o build/syntropy-linux-amd64 main.go

# Build for Windows (cross-compilation)
GOOS=windows GOARCH=amd64 go build -o build/syntropy-windows-amd64.exe main.go
```

## 🪟 Building on Windows

### Step 1: Prepare Environment

```powershell
# Open PowerShell as Administrator
# Navigate to the CLI directory
cd C:\Users\$env:USERNAME\syntropy-cc\syntropy-cooperative-grid\manager\interfaces\cli

# Verify Go version
go version
# Should show Go 1.22.5 or higher

# Verify we're in the correct directory
Get-Location
# Should show the CLI directory path
```

### Step 2: Install Dependencies

```powershell
# Download Go dependencies
go mod download

# Verify and organize dependencies
go mod tidy

# Verify no dependency issues
go mod verify
```

### Step 3: Build the CLI Manager

#### Option A: Automated Build (Recommended)

```powershell
# Build everything automatically
.\build.ps1

# Or build for specific platforms
.\build.ps1 current    # Current platform only
.\build.ps1 linux      # Linux only (cross-compilation)
.\build.ps1 windows    # Windows only
.\build.ps1 darwin     # macOS only (cross-compilation)
```

#### Option B: Using Make (if available)

```powershell
# Build for current platform
make build

# Build for all platforms
make cross-build

# Build with all features
make all
```

#### Option C: Manual Build

```powershell
# Build for Windows
$buildTime = Get-Date -Format "yyyy-MM-ddTHH:mm:ssZ"
$version = Get-Date -Format "yyyyMMdd-HHmmss"
$gitCommit = git rev-parse --short HEAD

go build -ldflags "-X main.version=$version -X main.buildTime=$buildTime -X main.gitCommit=$gitCommit" -o build\syntropy.exe main.go

# Build for Linux (cross-compilation)
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o build\syntropy-linux-amd64 main.go

# Restore environment
$env:GOOS = "windows"
$env:GOARCH = "amd64"
```

## 🧪 Testing

### Run Automated Tests

#### Linux
```bash
# Run tests using build script
./build.sh test

# Or using Make
make test
make test-coverage
make test-race

# Or manually
go test -v ./...
go test -v -cover ./...
go test -v -race ./...
```

#### Windows
```powershell
# Run tests using build script
.\build.ps1 test

# Or manually
go test -v .\...
go test -v -cover .\...
go test -v -race .\...
```

### Test Setup Component Specifically

```bash
# Test setup component
cd setup
go test -v ./...

# Test with coverage
go test -v -cover ./...

# Test specific functionality
go test -v -run TestSetupFlow ./setup_test.go
```

## 🚀 Execution and Manual Testing

### Test the CLI Manager

#### Linux
```bash
# Navigate to build directory
cd build

# Test basic functionality
./syntropy --help
./syntropy --version

# Test setup commands
./syntropy setup --help
./syntropy setup validate
./syntropy setup status
```

#### Windows
```powershell
# Navigate to build directory
cd build

# Test basic functionality
.\syntropy.exe --help
.\syntropy.exe --version

# Test setup commands
.\syntropy.exe setup --help
.\syntropy.exe setup validate
.\syntropy.exe setup status
```

### Test Setup Component Functionality

```bash
# Run setup process
./syntropy setup run --force

# Check setup status
./syntropy setup status

# Validate environment
./syntropy setup validate

# Reset configuration (if needed)
./syntropy setup reset --force
```

## 📊 Quality Verification

### Static Code Analysis

```bash
# Format code
go fmt ./...

# Run go vet
go vet ./...

# Run linter (if installed)
golangci-lint run

# Check for security issues (if installed)
gosec ./...
```

### Dependency Verification

```bash
# Check dependencies
go list -json -deps ./...

# Check for vulnerabilities
go list -json -deps ./... | nancy sleuth
```

## 🔍 Troubleshooting

### Common Build Issues

#### Error: "package not found"
```bash
# Solution: Download dependencies
go mod download
go mod tidy
```

#### Error: "build constraints exclude all Go files"
```bash
# Solution: Check build tags
# For Linux: go build -tags linux
# For Windows: go build -tags windows
```

#### Error: "permission denied" (Linux)
```bash
# Solution: Give execution permissions
chmod +x build/syntropy
```

#### Error: "cannot find main package" (Windows)
```powershell
# Solution: Verify you're in the correct directory
# The main.go file should be in manager/interfaces/cli/main.go
```

### Test Issues

#### Tests fail with "not implemented"
```bash
# This is expected for unimplemented features
# The stubs return ErrNotImplemented for non-implemented functionality
```

#### Integration tests fail
```bash
# Check if system dependencies are installed
# On Linux: systemd, systemctl
# On Windows: PowerShell, administrator permissions
```

### Execution Issues

#### "command not found" on Linux
```bash
# Solution: Add to PATH or use full path
export PATH=$PATH:/path/to/syntropy
# or
./build/syntropy --help
```

#### "execution policy" on Windows
```powershell
# Solution: Change execution policy
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

## 📝 Usage Examples

### Complete Build Example (Linux)

```bash
#!/bin/bash
# Complete build script for Linux

set -e

echo "=== Building Syntropy CLI Manager - Linux ==="

# Navigate to CLI directory
cd /home/jescott/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli

# Check Go
echo "Checking Go..."
go version

# Download dependencies
echo "Downloading dependencies..."
go mod download
go mod tidy

# Build CLI Manager
echo "Building CLI Manager..."
./build.sh

# Run tests
echo "Running tests..."
./build.sh test

# Verify binary
echo "Verifying binary..."
./build/syntropy --help

echo "=== Build completed successfully! ==="
```

### Complete Build Example (Windows)

```powershell
# Complete build script for Windows

Write-Host "=== Building Syntropy CLI Manager - Windows ===" -ForegroundColor Green

# Navigate to CLI directory
Set-Location "C:\Users\$env:USERNAME\syntropy-cc\syntropy-cooperative-grid\manager\interfaces\cli"

# Check Go
Write-Host "Checking Go..."
go version

# Download dependencies
Write-Host "Downloading dependencies..."
go mod download
go mod tidy

# Build CLI Manager
Write-Host "Building CLI Manager..."
.\build.ps1

# Run tests
Write-Host "Running tests..."
.\build.ps1 test

# Verify binary
Write-Host "Verifying binary..."
.\build\syntropy.exe --help

Write-Host "=== Build completed successfully! ===" -ForegroundColor Green
```

## 📚 Additional Resources

### Related Documentation
- [README.md](./README.md) - User documentation
- [GUIDE.md](./GUIDE.md) - Development guide
- [setup/README.md](./setup/README.md) - Setup component documentation

### Useful Commands
```bash
# Clean Go cache
go clean -cache

# Check unnecessary dependencies
go mod why <package>

# Update dependencies
go get -u ./...

# Check for vulnerabilities
go list -json -deps ./... | nancy sleuth
```

### Helpful Links
- [Official Go Documentation](https://golang.org/doc/)
- [Go Build Constraints](https://pkg.go.dev/go/build#hdr-Build_Constraints)
- [Go Testing](https://golang.org/pkg/testing/)
- [Cobra CLI Library](https://github.com/spf13/cobra)

---

## ✅ Build Checklist

### Linux
- [ ] Go 1.22.5+ installed
- [ ] Dependencies downloaded (`go mod download`)
- [ ] Code compiled without errors
- [ ] Tests run successfully
- [ ] Binary created and functional
- [ ] Execution permissions set

### Windows
- [ ] Go 1.22.5+ installed
- [ ] PowerShell 5.1+ available
- [ ] Dependencies downloaded (`go mod download`)
- [ ] Code compiled without errors
- [ ] Tests run successfully
- [ ] Executable created and functional
- [ ] Execution policy configured

### Both Platforms
- [ ] Linting run without errors
- [ ] Security analysis completed
- [ ] Documentation updated
- [ ] Integration tests run
- [ ] Functionality tested manually

---

**Last Updated**: $(date)
**Version**: 1.0
**Author**: Syntropy Development Team
