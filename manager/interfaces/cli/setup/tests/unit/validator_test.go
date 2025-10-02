//go:build !integration && !e2e && !performance && !security
// +build !integration,!e2e,!performance,!security

package unit

import (
	"os"
	"testing"

	setup "setup-component/src"
	"setup-component/src/internal/types"
)

// TestNewValidator testa a criação do validador
func TestNewValidator(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	tests := []struct {
		name    string
		logger  *setup.SetupLogger
		wantErr bool
	}{
		{
			name:    "should create validator successfully",
			logger:  logger,
			wantErr: false,
		},
		{
			name:    "should create validator with nil logger",
			logger:  nil,
			wantErr: false, // Logger pode ser nil
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := setup.NewValidator(tt.logger)
			if validator == nil {
				t.Error("NewValidator() returned nil validator")
			}
		})
	}
}

// TestNewOSValidator testa a criação do validador específico por SO
func TestNewOSValidator(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	tests := []struct {
		name   string
		logger *setup.SetupLogger
		os     string
	}{
		{
			name:   "should create Windows validator on Windows",
			logger: logger,
			os:     "windows",
		},
		{
			name:   "should create Linux validator on Linux",
			logger: logger,
			os:     "linux",
		},
		{
			name:   "should create Darwin validator on macOS",
			logger: logger,
			os:     "darwin",
		},
		{
			name:   "should create Generic validator on unknown OS",
			logger: logger,
			os:     "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Não podemos alterar runtime.GOOS diretamente, então testamos o comportamento atual
			validator := setup.NewOSValidator(tt.logger)
			if validator == nil {
				t.Error("NewOSValidator() returned nil validator")
			}
		})
	}
}

// TestValidator_ValidateEnvironment testa a validação do ambiente
func TestValidator_ValidateEnvironment(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	validator := setup.NewValidator(logger)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should validate environment successfully",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			envInfo, err := validator.ValidateEnvironment()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validator.ValidateEnvironment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if envInfo == nil {
					t.Error("Validator.ValidateEnvironment() returned nil environment info")
					return
				}

				// Verificar campos obrigatórios
				if envInfo.OS == "" {
					t.Error("Environment info missing OS")
				}
				if envInfo.Architecture == "" {
					t.Error("Environment info missing Architecture")
				}
				if envInfo.HomeDir == "" {
					t.Error("Environment info missing HomeDir")
				}
			}
		})
	}
}

// TestValidator_ValidateDependencies testa a validação de dependências
func TestValidator_ValidateDependencies(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	validator := setup.NewValidator(logger)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should validate dependencies successfully",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps, err := validator.ValidateDependencies()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validator.ValidateDependencies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if deps == nil {
					t.Error("Validator.ValidateDependencies() returned nil dependencies")
					return
				}

				// Verificar campos obrigatórios
				if deps.Required == nil {
					t.Error("Dependencies missing Required list")
				}
				if deps.Installed == nil {
					t.Error("Dependencies missing Installed list")
				}
				if deps.Missing == nil {
					t.Error("Dependencies missing Missing list")
				}
				if deps.Outdated == nil {
					t.Error("Dependencies missing Outdated list")
				}
			}
		})
	}
}

// TestValidator_ValidateNetwork testa a validação de rede
func TestValidator_ValidateNetwork(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	validator := setup.NewValidator(logger)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should validate network successfully",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			network, err := validator.ValidateNetwork()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validator.ValidateNetwork() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if network == nil {
					t.Error("Validator.ValidateNetwork() returned nil network info")
					return
				}

				// Verificar campos obrigatórios
				if network.PortsOpen == nil {
					t.Error("Network info missing PortsOpen list")
				}
			}
		})
	}
}

// TestValidator_ValidatePermissions testa a validação de permissões
func TestValidator_ValidatePermissions(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	validator := setup.NewValidator(logger)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should validate permissions successfully",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			permissions, err := validator.ValidatePermissions()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validator.ValidatePermissions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if permissions == nil {
					t.Error("Validator.ValidatePermissions() returned nil permissions")
					return
				}

				// Verificar campos obrigatórios
				if permissions.Issues == nil {
					t.Error("Permissions missing Issues list")
				}
			}
		})
	}
}

// TestValidator_FixIssues testa a correção de problemas
func TestValidator_FixIssues(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	validator := setup.NewValidator(logger)

	tests := []struct {
		name    string
		issues  []types.ValidationIssue
		wantErr bool
	}{
		{
			name:    "should handle empty issues list",
			issues:  []types.ValidationIssue{},
			wantErr: false,
		},
		{
			name: "should handle issues list with problems",
			issues: []types.ValidationIssue{
				{
					Type:        "test_issue",
					Severity:    "warning",
					Message:     "Test issue",
					Suggestions: []string{"Fix it"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.FixIssues(tt.issues)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validator.FixIssues() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidator_ValidateAll testa a validação completa
func TestValidator_ValidateAll(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	validator := setup.NewValidator(logger)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should validate all successfully",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validator.ValidateAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validator.ValidateAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if result == nil {
					t.Error("Validator.ValidateAll() returned nil result")
					return
				}

				// Verificar campos obrigatórios
				if result.Environment == nil {
					t.Error("Validation result missing environment info")
				}
				if result.Dependencies == nil {
					t.Error("Validation result missing dependencies info")
				}
				if result.Network == nil {
					t.Error("Validation result missing network info")
				}
				if result.Permissions == nil {
					t.Error("Validation result missing permissions info")
				}
				if result.Issues == nil {
					t.Error("Validation result missing issues list")
				}
				if result.Warnings == nil {
					t.Error("Validation result missing warnings list")
				}
			}
		})
	}
}
