package mocks

import (
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src/internal/types"
)

// MockValidator implementa a interface Validator para testes
type MockValidator struct {
	ValidateEnvironmentFunc  func() (*types.EnvironmentInfo, error)
	ValidateDependenciesFunc func() (*types.DependencyStatus, error)
	ValidateNetworkFunc      func() (*types.NetworkInfo, error)
	ValidatePermissionsFunc  func() (*types.PermissionStatus, error)
	FixIssuesFunc            func(issues []types.ValidationIssue) error
}

// ValidateEnvironment chama a função mock
func (m *MockValidator) ValidateEnvironment() (*types.EnvironmentInfo, error) {
	if m.ValidateEnvironmentFunc != nil {
		return m.ValidateEnvironmentFunc()
	}
	return &types.EnvironmentInfo{
		OS:              "linux",
		OSVersion:       "20.04",
		Architecture:    "amd64",
		HasAdminRights:  false,
		AvailableDiskGB: 100.0,
		HasInternet:     true,
		HomeDir:         "/home/testuser",
		CanProceed:      true,
		Issues:          []string{},
	}, nil
}

// ValidateDependencies chama a função mock
func (m *MockValidator) ValidateDependencies() (*types.DependencyStatus, error) {
	if m.ValidateDependenciesFunc != nil {
		return m.ValidateDependenciesFunc()
	}
	return &types.DependencyStatus{
		Required:  []types.Dependency{},
		Installed: []types.Dependency{},
		Missing:   []types.Dependency{},
		Outdated:  []types.Dependency{},
	}, nil
}

// ValidateNetwork chama a função mock
func (m *MockValidator) ValidateNetwork() (*types.NetworkInfo, error) {
	if m.ValidateNetworkFunc != nil {
		return m.ValidateNetworkFunc()
	}
	return &types.NetworkInfo{
		HasInternet:     true,
		Connectivity:    true,
		ProxyConfigured: false,
		FirewallActive:  false,
		PortsOpen:       []int{8080, 9090},
	}, nil
}

// ValidatePermissions chama a função mock
func (m *MockValidator) ValidatePermissions() (*types.PermissionStatus, error) {
	if m.ValidatePermissionsFunc != nil {
		return m.ValidatePermissionsFunc()
	}
	return &types.PermissionStatus{
		FileSystem: true,
		Network:    true,
		Service:    false,
		Admin:      false,
		Issues:     []string{},
	}, nil
}

// FixIssues chama a função mock
func (m *MockValidator) FixIssues(issues []types.ValidationIssue) error {
	if m.FixIssuesFunc != nil {
		return m.FixIssuesFunc(issues)
	}
	return nil
}

// MockOSValidator implementa a interface OSValidator para testes
type MockOSValidator struct {
	DetectOSFunc             func() (*types.OSInfo, error)
	ValidateResourcesFunc    func() (*types.ResourceInfo, error)
	ValidatePermissionsFunc  func() (*types.PermissionInfo, error)
	InstallDependenciesFunc  func(deps []types.Dependency) error
	ConfigureEnvironmentFunc func() error
}

// DetectOS chama a função mock
func (m *MockOSValidator) DetectOS() (*types.OSInfo, error) {
	if m.DetectOSFunc != nil {
		return m.DetectOSFunc()
	}
	return &types.OSInfo{
		Name:         "linux",
		Version:      "20.04",
		Architecture: "amd64",
		Build:        "5.4.0",
		Kernel:       "5.4.0-42-generic",
	}, nil
}

// ValidateResources chama a função mock
func (m *MockOSValidator) ValidateResources() (*types.ResourceInfo, error) {
	if m.ValidateResourcesFunc != nil {
		return m.ValidateResourcesFunc()
	}
	return &types.ResourceInfo{
		TotalMemoryGB:  8.0,
		AvailableMemGB: 4.0,
		CPUCores:       4,
		DiskSpaceGB:    50.0,
	}, nil
}

// ValidatePermissions chama a função mock
func (m *MockOSValidator) ValidatePermissions() (*types.PermissionInfo, error) {
	if m.ValidatePermissionsFunc != nil {
		return m.ValidatePermissionsFunc()
	}
	return &types.PermissionInfo{
		HasAdminRights: false,
		UserID:         "1000",
		GroupID:        "1000",
		Capabilities:   []string{"file_system", "network"},
	}, nil
}

// InstallDependencies chama a função mock
func (m *MockOSValidator) InstallDependencies(deps []types.Dependency) error {
	if m.InstallDependenciesFunc != nil {
		return m.InstallDependenciesFunc(deps)
	}
	return nil
}

// ConfigureEnvironment chama a função mock
func (m *MockOSValidator) ConfigureEnvironment() error {
	if m.ConfigureEnvironmentFunc != nil {
		return m.ConfigureEnvironmentFunc()
	}
	return nil
}
