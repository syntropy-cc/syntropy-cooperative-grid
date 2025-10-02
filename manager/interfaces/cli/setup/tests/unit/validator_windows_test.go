//go:build windows && !integration && !e2e && !performance && !security
// +build windows,!integration,!e2e,!performance,!security

package unit

import (
	"os"
	"testing"

	setup "setup-component/src"
)

// TestWindowsValidator_DetectOS testa a detecção do SO Windows
func TestWindowsValidator_DetectOS(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	validator := setup.NewWindowsValidator(logger)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should detect Windows OS successfully",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			osInfo, err := validator.DetectOS()
			if (err != nil) != tt.wantErr {
				t.Errorf("WindowsValidator.DetectOS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if osInfo == nil {
					t.Error("WindowsValidator.DetectOS() returned nil OS info")
					return
				}

				// Verificar campos específicos do Windows
				if osInfo.Name != "windows" {
					t.Errorf("WindowsValidator.DetectOS() = %v, want windows", osInfo.Name)
				}
				if osInfo.Version == "" {
					t.Error("WindowsValidator.DetectOS() missing version")
				}
				if osInfo.Architecture == "" {
					t.Error("WindowsValidator.DetectOS() missing architecture")
				}
				if osInfo.Build == "" {
					t.Error("WindowsValidator.DetectOS() missing build")
				}
				if osInfo.Kernel == "" {
					t.Error("WindowsValidator.DetectOS() missing kernel")
				}
			}
		})
	}
}

// TestWindowsValidator_ValidateResources testa a validação de recursos no Windows
func TestWindowsValidator_ValidateResources(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	validator := setup.NewWindowsValidator(logger)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should validate Windows resources successfully",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resources, err := validator.ValidateResources()
			if (err != nil) != tt.wantErr {
				t.Errorf("WindowsValidator.ValidateResources() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if resources == nil {
					t.Error("WindowsValidator.ValidateResources() returned nil resources")
					return
				}

				// Verificar campos obrigatórios
				if resources.TotalMemoryGB <= 0 {
					t.Error("WindowsValidator.ValidateResources() invalid total memory")
				}
				if resources.AvailableMemGB <= 0 {
					t.Error("WindowsValidator.ValidateResources() invalid available memory")
				}
				if resources.CPUCores <= 0 {
					t.Error("WindowsValidator.ValidateResources() invalid CPU cores")
				}
				if resources.DiskSpaceGB <= 0 {
					t.Error("WindowsValidator.ValidateResources() invalid disk space")
				}
			}
		})
	}
}

// TestWindowsValidator_ValidatePermissions testa a validação de permissões no Windows
func TestWindowsValidator_ValidatePermissions(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	validator := setup.NewWindowsValidator(logger)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should validate Windows permissions successfully",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			permissions, err := validator.ValidatePermissions()
			if (err != nil) != tt.wantErr {
				t.Errorf("WindowsValidator.ValidatePermissions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if permissions == nil {
					t.Error("WindowsValidator.ValidatePermissions() returned nil permissions")
					return
				}

				// Verificar campos obrigatórios
				if permissions.UserID == "" {
					t.Error("WindowsValidator.ValidatePermissions() missing user ID")
				}
				if permissions.GroupID == "" {
					t.Error("WindowsValidator.ValidatePermissions() missing group ID")
				}
				if permissions.Capabilities == nil {
					t.Error("WindowsValidator.ValidatePermissions() missing capabilities")
				}
			}
		})
	}
}

// TestWindowsValidator_InstallDependencies testa a instalação de dependências no Windows
func TestWindowsValidator_InstallDependencies(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	validator := setup.NewWindowsValidator(logger)

	tests := []struct {
		name    string
		deps    []types.Dependency
		wantErr bool
	}{
		{
			name:    "should handle empty dependencies list",
			deps:    []types.Dependency{},
			wantErr: false,
		},
		{
			name: "should handle dependencies list",
			deps: []types.Dependency{
				{
					Name:     "powershell",
					Version:  "5.1+",
					Required: true,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.InstallDependencies(tt.deps)
			if (err != nil) != tt.wantErr {
				t.Errorf("WindowsValidator.InstallDependencies() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestWindowsValidator_ConfigureEnvironment testa a configuração do ambiente Windows
func TestWindowsValidator_ConfigureEnvironment(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	validator := setup.NewWindowsValidator(logger)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should configure Windows environment successfully",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ConfigureEnvironment()
			if (err != nil) != tt.wantErr {
				t.Errorf("WindowsValidator.ConfigureEnvironment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
