// Package mocks provides mock implementations for testing the setup component
package mocks

import (
	"context"
	"time"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/tests/types"
)

// MockSetupService provides a mock implementation of the setup service
type MockSetupService struct {
	SetupFunc   func(ctx context.Context, options types.SetupOptions) (types.SetupResult, error)
	StatusFunc  func() (types.ValidationResult, error)
	ResetFunc   func(ctx context.Context) error
	ValidateFunc func() (types.ValidationResult, error)
	
	// Call tracking
	SetupCalls   []types.SetupOptions
	StatusCalls  int
	ResetCalls   int
	ValidateCalls int
}

// NewMockSetupService creates a new mock setup service
func NewMockSetupService() *MockSetupService {
	return &MockSetupService{
		SetupCalls: make([]types.SetupOptions, 0),
	}
}

// Setup mocks the setup operation
func (m *MockSetupService) Setup(ctx context.Context, options types.SetupOptions) (types.SetupResult, error) {
	m.SetupCalls = append(m.SetupCalls, options)
	
	if m.SetupFunc != nil {
		return m.SetupFunc(ctx, options)
	}
	
	// Default successful setup
	return types.SetupResult{
		Success:     true,
		StartTime:   time.Now().Add(-5 * time.Second),
		EndTime:     time.Now(),
		ConfigPath:  "/home/user/.syntropy/config.yaml",
		Environment: "linux",
		Options:     options,
		Message:     "Setup completed successfully",
	}, nil
}

// Status mocks the status check operation
func (m *MockSetupService) Status() (types.ValidationResult, error) {
	m.StatusCalls++
	
	if m.StatusFunc != nil {
		return m.StatusFunc()
	}
	
	// Default valid status
	return types.ValidationResult{
		Valid:    true,
		Warnings: []string{},
		Errors:   []string{},
		Environment: types.EnvironmentInfo{
			OS:              "linux",
			Architecture:    "amd64",
			HasAdminRights:  false,
			AvailableDiskGB: 100.0,
			HasInternet:     true,
			HomeDir:         "/home/user",
			SystemResources: types.SystemResources{
				TotalMemoryGB: 16.0,
				CPUCores:      8,
			},
		},
	}, nil
}

// Reset mocks the reset operation
func (m *MockSetupService) Reset(ctx context.Context) error {
	m.ResetCalls++
	
	if m.ResetFunc != nil {
		return m.ResetFunc(ctx)
	}
	
	// Default successful reset
	return nil
}

// Validate mocks the validation operation
func (m *MockSetupService) Validate() (types.ValidationResult, error) {
	m.ValidateCalls++
	
	if m.ValidateFunc != nil {
		return m.ValidateFunc()
	}
	
	// Default valid validation
	return types.ValidationResult{
		Valid:    true,
		Warnings: []string{},
		Errors:   []string{},
		Environment: types.EnvironmentInfo{
			OS:              "linux",
			Architecture:    "amd64",
			HasAdminRights:  false,
			AvailableDiskGB: 100.0,
			HasInternet:     true,
			HomeDir:         "/home/user",
			SystemResources: types.SystemResources{
				TotalMemoryGB: 16.0,
				CPUCores:      8,
			},
		},
	}, nil
}

// SetupError configures the mock to return an error on setup
func (m *MockSetupService) SetupError(err error) {
	m.SetupFunc = func(ctx context.Context, options types.SetupOptions) (types.SetupResult, error) {
		return types.SetupResult{
			Success:     false,
			StartTime:   time.Now(),
			EndTime:     time.Now(),
			Environment: "linux",
			Options:     options,
			Error:       err,
			Message:     err.Error(),
		}, err
	}
}

// StatusError configures the mock to return an error on status check
func (m *MockSetupService) StatusError(err error) {
	m.StatusFunc = func() (types.ValidationResult, error) {
		return types.ValidationResult{}, err
	}
}

// ResetError configures the mock to return an error on reset
func (m *MockSetupService) ResetError(err error) {
	m.ResetFunc = func(ctx context.Context) error {
		return err
	}
}

// ValidateError configures the mock to return an error on validation
func (m *MockSetupService) ValidateError(err error) {
	m.ValidateFunc = func() (types.ValidationResult, error) {
		return types.ValidationResult{}, err
	}
}

// InvalidEnvironment configures the mock to return invalid environment
func (m *MockSetupService) InvalidEnvironment(warnings, errors []string) {
	result := types.ValidationResult{
		Valid:    false,
		Warnings: warnings,
		Errors:   errors,
		Environment: types.EnvironmentInfo{
			OS:              "linux",
			Architecture:    "amd64",
			HasAdminRights:  false,
			AvailableDiskGB: 1.0, // Low disk space
			HasInternet:     false, // No internet
			HomeDir:         "/home/user",
			SystemResources: types.SystemResources{
				TotalMemoryGB: 2.0, // Low memory
				CPUCores:      1,   // Low CPU
			},
		},
	}
	
	m.StatusFunc = func() (types.ValidationResult, error) {
		return result, nil
	}
	
	m.ValidateFunc = func() (types.ValidationResult, error) {
		return result, nil
	}
}

// ResetMock resets all call counters and functions
func (m *MockSetupService) ResetMock() {
	m.SetupFunc = nil
	m.StatusFunc = nil
	m.ResetFunc = nil
	m.ValidateFunc = nil
	m.SetupCalls = make([]types.SetupOptions, 0)
	m.StatusCalls = 0
	m.ResetCalls = 0
	m.ValidateCalls = 0
}

// GetSetupCallCount returns the number of setup calls
func (m *MockSetupService) GetSetupCallCount() int {
	return len(m.SetupCalls)
}

// GetLastSetupOptions returns the options from the last setup call
func (m *MockSetupService) GetLastSetupOptions() *types.SetupOptions {
	if len(m.SetupCalls) == 0 {
		return nil
	}
	return &m.SetupCalls[len(m.SetupCalls)-1]
}