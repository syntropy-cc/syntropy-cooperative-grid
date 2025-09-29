//go:build windows
// +build windows

package unit

import (
	"testing"
	"time"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src/internal/types"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/tests/helpers"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/tests/mocks"
)

// TestWindowsValidator_DetectOS testa a detecção do SO Windows
func TestWindowsValidator_DetectOS(t *testing.T) {
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
			logger := &mocks.MockSetupLogger{}
			validator := src.NewOSValidator(logger)

			// Verificar se é um WindowsValidator
			windowsValidator, ok := validator.(*src.WindowsValidator)
			if !ok {
				t.Skip("Not running on Windows, skipping Windows-specific tests")
				return
			}

			result, err := windowsValidator.DetectOS()
			if (err != nil) != tt.wantErr {
				t.Errorf("WindowsValidator.DetectOS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if result == nil {
					t.Error("WindowsValidator.DetectOS() returned nil result")
					return
				}
				helpers.AssertStringEqual(t, result.Name, "windows", "OS Name")
				helpers.AssertStringNotEmpty(t, result.Version, "OS Version")
				helpers.AssertStringNotEmpty(t, result.Architecture, "Architecture")
				helpers.AssertStringNotEmpty(t, result.Build, "Build")
				helpers.AssertStringEqual(t, result.Kernel, "nt", "Kernel")
			}
		})
	}
}

// TestWindowsValidator_ValidateResources testa a validação de recursos no Windows
func TestWindowsValidator_ValidateResources(t *testing.T) {
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
			logger := &mocks.MockSetupLogger{}
			validator := src.NewOSValidator(logger)

			// Verificar se é um WindowsValidator
			windowsValidator, ok := validator.(*src.WindowsValidator)
			if !ok {
				t.Skip("Not running on Windows, skipping Windows-specific tests")
				return
			}

			result, err := windowsValidator.ValidateResources()
			if (err != nil) != tt.wantErr {
				t.Errorf("WindowsValidator.ValidateResources() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if result == nil {
					t.Error("WindowsValidator.ValidateResources() returned nil result")
					return
				}
				helpers.AssertFloat64Equal(t, result.TotalMemoryGB, 8.0, 0.1, "TotalMemoryGB")
				helpers.AssertFloat64Equal(t, result.AvailableMemGB, 4.0, 0.1, "AvailableMemGB")
				helpers.AssertIntEqual(t, result.CPUCores, 4, "CPUCores")
				helpers.AssertFloat64Equal(t, result.DiskSpaceGB, 50.0, 0.1, "DiskSpaceGB")
			}
		})
	}
}

// TestWindowsValidator_ValidatePermissions testa a validação de permissões no Windows
func TestWindowsValidator_ValidatePermissions(t *testing.T) {
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
			logger := &mocks.MockSetupLogger{}
			validator := src.NewOSValidator(logger)

			// Verificar se é um WindowsValidator
			windowsValidator, ok := validator.(*src.WindowsValidator)
			if !ok {
				t.Skip("Not running on Windows, skipping Windows-specific tests")
				return
			}

			result, err := windowsValidator.ValidatePermissions()
			if (err != nil) != tt.wantErr {
				t.Errorf("WindowsValidator.ValidatePermissions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if result == nil {
					t.Error("WindowsValidator.ValidatePermissions() returned nil result")
					return
				}
				helpers.AssertBoolEqual(t, result.HasAdminRights, false, "HasAdminRights")
				helpers.AssertStringEqual(t, result.UserID, "user", "UserID")
				helpers.AssertStringEqual(t, result.GroupID, "users", "GroupID")
				helpers.AssertSliceLength(t, result.Capabilities, 2, "Capabilities")
			}
		})
	}
}

// TestWindowsValidator_InstallDependencies testa a instalação de dependências no Windows
func TestWindowsValidator_InstallDependencies(t *testing.T) {
	tests := []struct {
		name    string
		deps    []types.Dependency
		wantErr bool
	}{
		{
			name:    "should install Windows dependencies successfully",
			deps:    []types.Dependency{},
			wantErr: false,
		},
		{
			name: "should handle non-empty dependencies list",
			deps: []types.Dependency{
				{
					Name:      "powershell",
					Version:   "5.1+",
					Required:  true,
					Installed: false,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &mocks.MockSetupLogger{}
			validator := src.NewOSValidator(logger)

			// Verificar se é um WindowsValidator
			windowsValidator, ok := validator.(*src.WindowsValidator)
			if !ok {
				t.Skip("Not running on Windows, skipping Windows-specific tests")
				return
			}

			err := windowsValidator.InstallDependencies(tt.deps)
			if (err != nil) != tt.wantErr {
				t.Errorf("WindowsValidator.InstallDependencies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

// TestWindowsValidator_ConfigureEnvironment testa a configuração do ambiente Windows
func TestWindowsValidator_ConfigureEnvironment(t *testing.T) {
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
			logger := &mocks.MockSetupLogger{}
			validator := src.NewOSValidator(logger)

			// Verificar se é um WindowsValidator
			windowsValidator, ok := validator.(*src.WindowsValidator)
			if !ok {
				t.Skip("Not running on Windows, skipping Windows-specific tests")
				return
			}

			err := windowsValidator.ConfigureEnvironment()
			if (err != nil) != tt.wantErr {
				t.Errorf("WindowsValidator.ConfigureEnvironment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

// TestWindowsValidator_EdgeCases testa casos extremos do WindowsValidator
func TestWindowsValidator_EdgeCases(t *testing.T) {
	t.Run("should handle nil logger", func(t *testing.T) {
		validator := src.NewOSValidator(nil)

		// Verificar se é um WindowsValidator
		windowsValidator, ok := validator.(*src.WindowsValidator)
		if !ok {
			t.Skip("Not running on Windows, skipping Windows-specific tests")
			return
		}

		// Should not panic
		_, err := windowsValidator.DetectOS()
		if err != nil {
			t.Errorf("DetectOS() failed with nil logger: %v", err)
		}
	})

	t.Run("should handle nil dependencies", func(t *testing.T) {
		logger := &mocks.MockSetupLogger{}
		validator := src.NewOSValidator(logger)

		// Verificar se é um WindowsValidator
		windowsValidator, ok := validator.(*src.WindowsValidator)
		if !ok {
			t.Skip("Not running on Windows, skipping Windows-specific tests")
			return
		}

		err := windowsValidator.InstallDependencies(nil)
		if err != nil {
			t.Errorf("InstallDependencies() failed with nil dependencies: %v", err)
		}
	})
}

// TestWindowsValidator_Concurrency testa concorrência do WindowsValidator
func TestWindowsValidator_Concurrency(t *testing.T) {
	t.Run("should handle concurrent calls", func(t *testing.T) {
		logger := &mocks.MockSetupLogger{}
		validator := src.NewOSValidator(logger)

		// Verificar se é um WindowsValidator
		windowsValidator, ok := validator.(*src.WindowsValidator)
		if !ok {
			t.Skip("Not running on Windows, skipping Windows-specific tests")
			return
		}

		// Executar múltiplas chamadas concorrentemente
		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func() {
				_, err := windowsValidator.DetectOS()
				if err != nil {
					t.Errorf("Concurrent DetectOS() failed: %v", err)
				}
				done <- true
			}()
		}

		// Aguardar todas as goroutines terminarem
		for i := 0; i < 10; i++ {
			<-done
		}
	})
}

// TestWindowsValidator_Performance testa performance do WindowsValidator
func TestWindowsValidator_Performance(t *testing.T) {
	t.Run("should complete operations within reasonable time", func(t *testing.T) {
		logger := &mocks.MockSetupLogger{}
		validator := src.NewOSValidator(logger)

		// Verificar se é um WindowsValidator
		windowsValidator, ok := validator.(*src.WindowsValidator)
		if !ok {
			t.Skip("Not running on Windows, skipping Windows-specific tests")
			return
		}

		start := time.Now()
		_, err := windowsValidator.DetectOS()
		elapsed := time.Since(start)

		if err != nil {
			t.Errorf("DetectOS() failed: %v", err)
		}

		if elapsed > 100*time.Millisecond {
			t.Errorf("DetectOS() took too long: %v", elapsed)
		}
	})
}
