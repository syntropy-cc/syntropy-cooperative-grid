//go:build linux
// +build linux

package unit

import (
	"testing"
	"time"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src/internal/types"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/tests/helpers"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/tests/mocks"
)

// TestLinuxValidator_DetectOS testa a detecção do SO Linux
func TestLinuxValidator_DetectOS(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should detect Linux OS successfully",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &mocks.MockSetupLogger{}
			validator := src.NewOSValidator(logger)

			// Verificar se é um LinuxValidator
			linuxValidator, ok := validator.(*src.LinuxValidator)
			if !ok {
				t.Skip("Not running on Linux, skipping Linux-specific tests")
				return
			}

			result, err := linuxValidator.DetectOS()
			if (err != nil) != tt.wantErr {
				t.Errorf("LinuxValidator.DetectOS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if result == nil {
					t.Error("LinuxValidator.DetectOS() returned nil result")
					return
				}
				helpers.AssertStringEqual(t, result.Name, "linux", "OS Name")
				helpers.AssertStringEqual(t, result.Version, "20.04", "OS Version")
				helpers.AssertStringNotEmpty(t, result.Architecture, "Architecture")
				helpers.AssertStringEqual(t, result.Build, "5.4.0", "Build")
				helpers.AssertStringEqual(t, result.Kernel, "5.4.0-42-generic", "Kernel")
			}
		})
	}
}

// TestLinuxValidator_ValidateResources testa a validação de recursos no Linux
func TestLinuxValidator_ValidateResources(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should validate Linux resources successfully",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &mocks.MockSetupLogger{}
			validator := src.NewOSValidator(logger)

			// Verificar se é um LinuxValidator
			linuxValidator, ok := validator.(*src.LinuxValidator)
			if !ok {
				t.Skip("Not running on Linux, skipping Linux-specific tests")
				return
			}

			result, err := linuxValidator.ValidateResources()
			if (err != nil) != tt.wantErr {
				t.Errorf("LinuxValidator.ValidateResources() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if result == nil {
					t.Error("LinuxValidator.ValidateResources() returned nil result")
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

// TestLinuxValidator_ValidatePermissions testa a validação de permissões no Linux
func TestLinuxValidator_ValidatePermissions(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should validate Linux permissions successfully",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &mocks.MockSetupLogger{}
			validator := src.NewOSValidator(logger)

			// Verificar se é um LinuxValidator
			linuxValidator, ok := validator.(*src.LinuxValidator)
			if !ok {
				t.Skip("Not running on Linux, skipping Linux-specific tests")
				return
			}

			result, err := linuxValidator.ValidatePermissions()
			if (err != nil) != tt.wantErr {
				t.Errorf("LinuxValidator.ValidatePermissions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if result == nil {
					t.Error("LinuxValidator.ValidatePermissions() returned nil result")
					return
				}
				helpers.AssertBoolEqual(t, result.HasAdminRights, false, "HasAdminRights")
				helpers.AssertStringEqual(t, result.UserID, "1000", "UserID")
				helpers.AssertStringEqual(t, result.GroupID, "1000", "GroupID")
				helpers.AssertSliceLength(t, result.Capabilities, 2, "Capabilities")
			}
		})
	}
}

// TestLinuxValidator_InstallDependencies testa a instalação de dependências no Linux
func TestLinuxValidator_InstallDependencies(t *testing.T) {
	tests := []struct {
		name    string
		deps    []types.Dependency
		wantErr bool
	}{
		{
			name:    "should install Linux dependencies successfully",
			deps:    []types.Dependency{},
			wantErr: false,
		},
		{
			name: "should handle non-empty dependencies list",
			deps: []types.Dependency{
				{
					Name:      "curl",
					Version:   "any",
					Required:  true,
					Installed: false,
				},
				{
					Name:      "systemctl",
					Version:   "any",
					Required:  false,
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

			// Verificar se é um LinuxValidator
			linuxValidator, ok := validator.(*src.LinuxValidator)
			if !ok {
				t.Skip("Not running on Linux, skipping Linux-specific tests")
				return
			}

			err := linuxValidator.InstallDependencies(tt.deps)
			if (err != nil) != tt.wantErr {
				t.Errorf("LinuxValidator.InstallDependencies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

// TestLinuxValidator_ConfigureEnvironment testa a configuração do ambiente Linux
func TestLinuxValidator_ConfigureEnvironment(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should configure Linux environment successfully",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &mocks.MockSetupLogger{}
			validator := src.NewOSValidator(logger)

			// Verificar se é um LinuxValidator
			linuxValidator, ok := validator.(*src.LinuxValidator)
			if !ok {
				t.Skip("Not running on Linux, skipping Linux-specific tests")
				return
			}

			err := linuxValidator.ConfigureEnvironment()
			if (err != nil) != tt.wantErr {
				t.Errorf("LinuxValidator.ConfigureEnvironment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

// TestLinuxValidator_EdgeCases testa casos extremos do LinuxValidator
func TestLinuxValidator_EdgeCases(t *testing.T) {
	t.Run("should handle nil logger", func(t *testing.T) {
		validator := src.NewOSValidator(nil)

		// Verificar se é um LinuxValidator
		linuxValidator, ok := validator.(*src.LinuxValidator)
		if !ok {
			t.Skip("Not running on Linux, skipping Linux-specific tests")
			return
		}

		// Should not panic
		_, err := linuxValidator.DetectOS()
		if err != nil {
			t.Errorf("DetectOS() failed with nil logger: %v", err)
		}
	})

	t.Run("should handle nil dependencies", func(t *testing.T) {
		logger := &mocks.MockSetupLogger{}
		validator := src.NewOSValidator(logger)

		// Verificar se é um LinuxValidator
		linuxValidator, ok := validator.(*src.LinuxValidator)
		if !ok {
			t.Skip("Not running on Linux, skipping Linux-specific tests")
			return
		}

		err := linuxValidator.InstallDependencies(nil)
		if err != nil {
			t.Errorf("InstallDependencies() failed with nil dependencies: %v", err)
		}
	})
}

// TestLinuxValidator_Concurrency testa concorrência do LinuxValidator
func TestLinuxValidator_Concurrency(t *testing.T) {
	t.Run("should handle concurrent calls", func(t *testing.T) {
		logger := &mocks.MockSetupLogger{}
		validator := src.NewOSValidator(logger)

		// Verificar se é um LinuxValidator
		linuxValidator, ok := validator.(*src.LinuxValidator)
		if !ok {
			t.Skip("Not running on Linux, skipping Linux-specific tests")
			return
		}

		// Executar múltiplas chamadas concorrentemente
		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func() {
				_, err := linuxValidator.DetectOS()
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

// TestLinuxValidator_Performance testa performance do LinuxValidator
func TestLinuxValidator_Performance(t *testing.T) {
	t.Run("should complete operations within reasonable time", func(t *testing.T) {
		logger := &mocks.MockSetupLogger{}
		validator := src.NewOSValidator(logger)

		// Verificar se é um LinuxValidator
		linuxValidator, ok := validator.(*src.LinuxValidator)
		if !ok {
			t.Skip("Not running on Linux, skipping Linux-specific tests")
			return
		}

		start := time.Now()
		_, err := linuxValidator.DetectOS()
		elapsed := time.Since(start)

		if err != nil {
			t.Errorf("DetectOS() failed: %v", err)
		}

		if elapsed > 100*time.Millisecond {
			t.Errorf("DetectOS() took too long: %v", elapsed)
		}
	})
}
