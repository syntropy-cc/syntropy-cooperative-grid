# Syntropy Setup Component Unit Tests

This document provides instructions for running and understanding the unit tests for the Syntropy Cooperative Grid's setup component.

## Prerequisites

Before running the tests, ensure you have:

1. Go installed (version 1.16 or later recommended)
2. The Syntropy Cooperative Grid repository cloned locally
3. All dependencies installed (run `go mod download` in the project root)

## Running the Tests

### Running All Tests

To run all unit tests for the setup component:

```bash
cd /path/to/syntropy-cooperative-grid/manager/interfaces/cli/setup/tests/unit
go test -v ./...
```

### Running Specific Tests

To run a specific test file:

```bash
go test -v ./setup_test.go
```

To run tests with a specific build tag (e.g., Linux-specific tests):

```bash
go test -v -tags=linux ./...
```

### Test Coverage

To generate test coverage reports:

```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

This will create an HTML report showing which code paths are covered by tests.

## Expected Test Behaviors

The unit tests verify the following behaviors:

1. **Directory Management**
   - Creation of the Syntropy directory structure
   - Proper handling of configuration files
   - Reset functionality that cleans up directories

2. **Configuration**
   - Proper loading and validation of configuration files
   - Default configuration values are set correctly
   - Custom configuration paths are respected

3. **Environment Validation**
   - OS and architecture detection
   - Environment requirements validation
   - Force flag behavior for overriding warnings

4. **Setup Process**
   - Success and error handling during setup
   - Proper reporting of setup results

## Common Test Scenarios

### Testing Basic Setup

The tests verify that the basic setup process creates all necessary directories and configuration files:

```go
// Example from tests
result, err := setup.Status(options)
if err != nil {
    t.Fatalf("Status check failed: %v", err)
}

if !result.Success {
    t.Error("Expected status check to succeed")
}
```

### Testing Reset Functionality

Tests ensure that the reset functionality properly removes existing Syntropy directories:

```go
// Example from tests
result, err := setup.Reset(options)
if err != nil {
    t.Fatalf("Reset failed: %v", err)
}

if !result.Success {
    t.Error("Expected reset to succeed")
}
```

### Testing Environment Validation

Tests verify that the environment validation correctly identifies the OS, architecture, and other system requirements:

```go
// Example from tests
result, err := setup.ValidateLinuxEnvironment(false)
if err != nil {
    t.Fatalf("ValidateLinuxEnvironment failed: %v", err)
}
```

## Troubleshooting

If you encounter test failures:

1. Ensure your Go environment is properly set up
2. Check that you're running the tests with the correct build tags for your OS
3. Verify that you have the necessary permissions to create and delete files in the test directories
4. For Linux-specific tests, ensure you're running on a Linux system or using appropriate virtualization

## Additional Resources

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Syntropy Cooperative Grid Documentation](https://docs.syntropy.network/)