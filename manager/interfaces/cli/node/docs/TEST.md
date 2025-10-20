# Node Component - Testing Guide

## Testing Philosophy

[MACRO_VIEW]
The Node Component's testing strategy contributes to the overall project quality by ensuring 100% reliability of the critical node provisioning pipeline, which is essential for the Syntropy Cooperative Grid's infrastructure automation and operational excellence.
[/MACRO_VIEW]

[MESO_VIEW]
Testing integrates with the Syntropy Manager CLI module's testing strategy through shared test utilities, consistent patterns, and comprehensive integration testing that validates the interaction between Node, Setup, and future components like Workload.
[/MESO_VIEW]

[MICRO_VIEW]
The specific testing approach focuses on comprehensive coverage of the two-subcomponent architecture (Create + Registration), platform-specific implementations, security validation, and end-to-end plug-and-play workflows.
[/MICRO_VIEW]

## Test Pyramid

### Distribution Strategy
- **Unit Tests**: 70% of test suite
  - Individual component testing
  - Mocked dependencies
  - Fast execution (< 1 second per test)
- **Integration Tests**: 25% of test suite
  - Component interaction testing
  - Real dependencies with controlled environment
  - Medium execution time (< 10 seconds per test)
- **End-to-End Tests**: 5% of test suite
  - Full workflow testing
  - Real hardware simulation
  - Slower execution (< 60 seconds per test)

### Test Categories
| Category | Purpose | Tools | Coverage Target |
|----------|---------|-------|-----------------|
| Unit | Component isolation | Go testing, mocks | 100% line/branch |
| Integration | Component interaction | Test containers, real USB simulation | 100% API coverage |
| E2E | Full workflow | Virtual machines, real hardware | 100% user scenarios |
| Performance | Load and stress testing | Benchmarks, profiling | 95% performance targets |
| Security | Security validation | Security scanners, penetration testing | 100% security controls |

## Test Environment Setup

### Prerequisites
- Go 1.21+ installed
- Docker for integration tests
- VirtualBox/VMware for E2E tests
- USB devices for hardware testing

### Environment Configuration
```bash
# 1. Clone repository
git clone <repository-url>
cd syntropy-cooperative-grid/manager/interfaces/cli/node

# 2. Install dependencies
go mod tidy
go mod download

# 3. Setup test environment
export SYNTHROPY_TEST_MODE=1
export SYNTHROPY_TEST_USB_PATH=/dev/sdb  # Use test USB device
export SYNTHROPY_TEST_CACHE_DIR=/tmp/syntropy-test-cache

# 4. Verify setup
go test ./tests/unit/... -v
```

### Test Data Setup
```bash
# Create test fixtures
mkdir -p tests/fixtures/{valid,invalid,edge-cases}
mkdir -p tests/fixtures/valid/node-configs
mkdir -p tests/fixtures/valid/cloud-init
mkdir -p tests/fixtures/invalid/malformed-configs
mkdir -p tests/fixtures/invalid/invalid-tokens
mkdir -p tests/fixtures/edge-cases/boundary-values
```

## Running Tests

### Quick Start
```bash
# Run all tests
go test ./...

# Run specific test categories
go test ./tests/unit/... -v
go test ./tests/integration/... -v
go test ./tests/e2e/... -v

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Platform-Specific Instructions

#### Windows (PowerShell)
```powershell
# Enable USB device access
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser

# Run tests
go test ./tests/unit/... -v
go test ./tests/integration/... -v

# Check USB devices
Get-WmiObject -Class Win32_LogicalDisk | Where-Object {$_.DriveType -eq 2}
```

#### Linux (bash)
```bash
# Add user to disk group for USB access
sudo usermod -a -G disk $USER
newgrp disk

# Run tests
go test ./tests/unit/... -v
go test ./tests/integration/... -v

# Check USB devices
lsblk
```

#### macOS (zsh)
```zsh
# Grant disk access permissions
# System Preferences > Security & Privacy > Privacy > Full Disk Access

# Run tests
go test ./tests/unit/... -v
go test ./tests/integration/... -v

# Check USB devices
diskutil list
```

### Test Execution Strategies

#### Parallel Execution
```bash
# Run unit tests in parallel (default)
go test -parallel 4 ./tests/unit/...

# Run integration tests sequentially (USB access conflicts)
go test -parallel 1 ./tests/integration/...
```

#### Test Filtering
```bash
# Run tests matching pattern
go test -run "TestUSBDetection" ./...

# Skip long-running tests
go test -short ./...

# Run only unit tests
go test ./tests/unit/... -tags=unit
```

## Test Organization

### Directory Structure
```
tests/
├── unit/                    # Unit tests (70% of suite)
│   ├── node_test.go        # NodeManager tests
│   ├── create_test.go      # CreateSubcomponent tests
│   ├── registration_test.go # RegistrationSubcomponent tests
│   ├── usb_detector_*_test.go # Platform-specific USB tests
│   ├── cloud_init_generator_test.go
│   ├── handshake_test.go
│   ├── listener_test.go
│   ├── heartbeat_test.go
│   └── node_state_test.go
├── integration/             # Integration tests (25% of suite)
│   ├── api/
│   │   └── token_integration_test.go
│   ├── create/
│   │   ├── usb_creation_flow_test.go
│   │   └── cloud_init_injection_test.go
│   ├── registration/
│   │   ├── handshake_flow_test.go
│   │   └── node_lifecycle_test.go
│   └── plug_and_play_test.go
├── e2e/                     # End-to-end tests (5% of suite)
│   └── scenarios/
│       ├── single_node_creation_test.go
│       ├── multiple_nodes_test.go
│       └── node_failure_recovery_test.go
├── performance/             # Performance tests
│   ├── usb_write_performance_test.go
│   ├── handshake_latency_test.go
│   └── concurrent_nodes_test.go
├── security/                # Security tests
│   ├── token_validation_test.go
│   ├── certificate_validation_test.go
│   └── ssh_security_test.go
├── fixtures/                # Test data
├── mocks/                   # Test doubles
└── helpers/                 # Test utilities
```

### Test File Naming Convention
- `*_test.go` for all test files
- `*_linux_test.go` for Linux-specific tests
- `*_windows_test.go` for Windows-specific tests
- `*_integration_test.go` for integration tests
- `*_e2e_test.go` for end-to-end tests

## Code Coverage Requirements

### Coverage Targets
| Component | Line Coverage | Branch Coverage | Function Coverage |
|-----------|---------------|-----------------|-------------------|
| src/node.go | 100% | 100% | 100% |
| src/create.go | 100% | 100% | 100% |
| src/registration.go | 100% | 100% | 100% |
| src/usb_detector*.go | 100% | 100% | 100% |
| src/cloud_init_generator.go | 100% | 100% | 100% |
| src/handshake.go | 100% | 100% | 100% |
| src/listener.go | 100% | 100% | 100% |
| src/heartbeat.go | 100% | 100% | 100% |
| src/node_state.go | 100% | 100% | 100% |
| **Overall** | **100%** | **100%** | **100%** |

### Coverage Verification
```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# Check coverage thresholds
go test -cover ./... | grep -E "(coverage: [0-9]+%)"
```

### Coverage Exclusions
```go
// Exclude generated code
//go:build !generate

// Exclude platform-specific code in wrong environment
// +build linux
// +build !windows

// Exclude test files from coverage
// +build test
```

## CI/CD Integration

### GitHub Actions Example
```yaml
name: Node Component Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21
    
    - name: Run unit tests
      run: go test ./tests/unit/... -v -coverprofile=unit.out
    
    - name: Run integration tests
      run: go test ./tests/integration/... -v -coverprofile=integration.out
      env:
        SYNTHROPY_TEST_MODE: 1
    
    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out
```

### Pre-commit Hooks
```bash
#!/bin/bash
# .git/hooks/pre-commit

echo "Running Node Component tests..."

# Run unit tests
go test ./tests/unit/... -short
if [ $? -ne 0 ]; then
    echo "Unit tests failed"
    exit 1
fi

# Check coverage
coverage=$(go test -cover ./tests/unit/... | grep "coverage:" | awk '{print $5}' | sed 's/%//')
if (( $(echo "$coverage < 100" | bc -l) )); then
    echo "Coverage too low: $coverage%"
    exit 1
fi

echo "All tests passed!"
```

## Performance Testing

### Benchmark Tests
```go
// tests/performance/usb_write_performance_test.go
func BenchmarkUSBWrite(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // Test USB write performance
    }
}

func BenchmarkHandshakeLatency(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // Test handshake latency
    }
}
```

### Performance Targets
| Metric | Target | Measurement |
|--------|--------|-------------|
| USB Creation Time | < 5 minutes | Time to create bootable USB |
| Handshake Latency | < 30 seconds | Time from connection to registration |
| Memory Usage | < 200MB | Peak memory during operations |
| CPU Usage | < 15% | Peak CPU during operations |

### Load Testing
```bash
# Run performance benchmarks
go test -bench=. ./tests/performance/...

# Run with profiling
go test -bench=. -cpuprofile=cpu.prof ./tests/performance/...
go tool pprof cpu.prof
```

## Security Testing

### Security Test Categories
| Category | Purpose | Tools | Coverage |
|----------|---------|-------|----------|
| Token Validation | Verify token security | Custom validators | 100% |
| Certificate Validation | Verify certificate security | OpenSSL, custom | 100% |
| SSH Security | Verify SSH key security | SSH tools, custom | 100% |
| USB Security | Verify USB write security | Custom validators | 100% |

### Security Test Examples
```go
// tests/security/token_validation_test.go
func TestTokenValidation(t *testing.T) {
    tests := []struct {
        name    string
        token   string
        wantErr bool
    }{
        {"valid_token", "valid-token", false},
        {"invalid_token", "invalid", true},
        {"empty_token", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateToken(tt.token)
            if (err != nil) != tt.wantErr {
                t.Errorf("validateToken() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Debugging Tests

### Test Debugging Techniques
| Issue | Debug Method | Tools |
|-------|--------------|-------|
| USB Detection Failure | Verbose logging, device enumeration | lsblk, WMIC, diskutil |
| Handshake Timeout | Network packet capture, timeout analysis | Wireshark, tcpdump |
| State Management Bugs | State transition logging, race detection | go race detector, custom logger |
| Performance Issues | Profiling, benchmark analysis | pprof, benchstat |

### Debug Mode
```bash
# Enable debug mode
export SYNTHROPY_DEBUG=1
export SYNTHROPY_LOG_LEVEL=debug

# Run tests with debug output
go test -v ./tests/unit/...

# Run specific test with debug
go test -run TestUSBDetection -v ./tests/unit/...
```

### Test Logging
```go
// Enable test logging
func TestUSBDetection(t *testing.T) {
    if testing.Verbose() {
        t.Log("Testing USB detection...")
    }
    
    // Test implementation
}
```

## Troubleshooting by Platform

### Windows Issues
| Issue | Symptoms | Solution |
|-------|----------|----------|
| USB Access Denied | "Access denied" errors | Run PowerShell as Administrator |
| WMIC Not Found | "WMIC not recognized" | Ensure Windows Management Instrumentation is enabled |
| PowerShell Execution Policy | Script execution blocked | Set-ExecutionPolicy RemoteSigned |

### Linux Issues
| Issue | Symptoms | Solution |
|-------|----------|----------|
| Permission Denied | "Permission denied" on USB | Add user to disk group: `sudo usermod -a -G disk $USER` |
| USB Device Not Found | "No USB devices detected" | Check udev rules, ensure USB is properly mounted |
| Missing Dependencies | "Command not found" errors | Install required packages: `sudo apt install util-linux` |

### macOS Issues
| Issue | Symptoms | Solution |
|-------|----------|----------|
| Disk Access Denied | "Operation not permitted" | Grant Full Disk Access in System Preferences |
| USB Device Not Recognized | "No USB devices found" | Check System Information > USB |
| Missing Developer Tools | "Command not found" | Install Xcode Command Line Tools |

## Test Maintenance

### Test Health Metrics
| Metric | Target | Measurement |
|--------|--------|-------------|
| Test Execution Time | < 10 minutes | Total test suite runtime |
| Test Flakiness | < 1% | Percentage of intermittently failing tests |
| Test Coverage | 100% | Line, branch, and function coverage |
| Test Maintenance | < 2 hours/week | Time spent maintaining tests |

### Test Refactoring Guidelines
1. **Extract Common Setup**: Move common test setup to helper functions
2. **Parameterize Tests**: Use table-driven tests for multiple scenarios
3. **Mock External Dependencies**: Use mocks for network, filesystem, and hardware
4. **Clean Test Data**: Ensure tests clean up after themselves

### Test Documentation
- Each test file should have a header comment explaining its purpose
- Complex test scenarios should be documented with examples
- Test data should be documented with expected outcomes

## Educational Resources
For learning how to write effective tests and understanding testing principles, see [LEARN.md](./LEARN.md).
