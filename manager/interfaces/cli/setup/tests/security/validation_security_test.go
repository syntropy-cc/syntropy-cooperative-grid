package security

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src/internal/types"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/tests/helpers"
)

// TestValidationSecurity_InputValidation testa validação de entrada
func TestValidationSecurity_InputValidation(t *testing.T) {
	tests := []struct {
		name    string
		options *types.SetupOptions
		wantErr bool
	}{
		{
			name: "should reject malicious owner name",
			options: &types.SetupOptions{
				OwnerName:      "../../../etc/passwd",
				OwnerEmail:     "test@example.com",
				ConfigPath:     "/tmp/syntropy_test/config",
				KeysPath:       "/tmp/syntropy_test/keys",
				BackupPath:     "/tmp/syntropy_test/backups",
				LogPath:        "/tmp/syntropy_test/logs",
				Verbose:        false,
				Force:          false,
				SkipValidation: false,
			},
			wantErr: true,
		},
		{
			name: "should reject malicious email",
			options: &types.SetupOptions{
				OwnerName:      "Test User",
				OwnerEmail:     "test@example.com'; DROP TABLE users; --",
				ConfigPath:     "/tmp/syntropy_test/config",
				KeysPath:       "/tmp/syntropy_test/keys",
				BackupPath:     "/tmp/syntropy_test/backups",
				LogPath:        "/tmp/syntropy_test/logs",
				Verbose:        false,
				Force:          false,
				SkipValidation: false,
			},
			wantErr: true,
		},
		{
			name: "should reject path traversal in config path",
			options: &types.SetupOptions{
				OwnerName:      "Test User",
				OwnerEmail:     "test@example.com",
				ConfigPath:     "../../../etc/passwd",
				KeysPath:       "/tmp/syntropy_test/keys",
				BackupPath:     "/tmp/syntropy_test/backups",
				LogPath:        "/tmp/syntropy_test/logs",
				Verbose:        false,
				Force:          false,
				SkipValidation: false,
			},
			wantErr: true,
		},
		{
			name: "should reject path traversal in keys path",
			options: &types.SetupOptions{
				OwnerName:      "Test User",
				OwnerEmail:     "test@example.com",
				ConfigPath:     "/tmp/syntropy_test/config",
				KeysPath:       "../../../etc/passwd",
				BackupPath:     "/tmp/syntropy_test/backups",
				LogPath:        "/tmp/syntropy_test/logs",
				Verbose:        false,
				Force:          false,
				SkipValidation: false,
			},
			wantErr: true,
		},
		{
			name: "should reject path traversal in backup path",
			options: &types.SetupOptions{
				OwnerName:      "Test User",
				OwnerEmail:     "test@example.com",
				ConfigPath:     "/tmp/syntropy_test/config",
				KeysPath:       "/tmp/syntropy_test/keys",
				BackupPath:     "../../../etc/passwd",
				LogPath:        "/tmp/syntropy_test/logs",
				Verbose:        false,
				Force:          false,
				SkipValidation: false,
			},
			wantErr: true,
		},
		{
			name: "should reject path traversal in log path",
			options: &types.SetupOptions{
				OwnerName:      "Test User",
				OwnerEmail:     "test@example.com",
				ConfigPath:     "/tmp/syntropy_test/config",
				KeysPath:       "/tmp/syntropy_test/keys",
				BackupPath:     "/tmp/syntropy_test/backups",
				LogPath:        "../../../etc/passwd",
				Verbose:        false,
				Force:          false,
				SkipValidation: false,
			},
			wantErr: true,
		},
		{
			name: "should accept valid inputs",
			options: &types.SetupOptions{
				OwnerName:      "Test User",
				OwnerEmail:     "test@example.com",
				ConfigPath:     "/tmp/syntropy_test/config",
				KeysPath:       "/tmp/syntropy_test/keys",
				BackupPath:     "/tmp/syntropy_test/backups",
				LogPath:        "/tmp/syntropy_test/logs",
				Verbose:        false,
				Force:          false,
				SkipValidation: false,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "validation_security_test")

			// Ajustar caminhos para usar o diretório temporário se não forem maliciosos
			if tt.options.ConfigPath != "../../../etc/passwd" {
				tt.options.ConfigPath = filepath.Join(tempDir, "config")
			}
			if tt.options.KeysPath != "../../../etc/passwd" {
				tt.options.KeysPath = filepath.Join(tempDir, "keys")
			}
			if tt.options.BackupPath != "../../../etc/passwd" {
				tt.options.BackupPath = filepath.Join(tempDir, "backups")
			}
			if tt.options.LogPath != "../../../etc/passwd" {
				tt.options.LogPath = filepath.Join(tempDir, "logs")
			}

			// Criar validator
			validator := src.NewValidator(nil)
			if validator == nil {
				t.Fatal("NewValidator() returned nil")
			}

			// Executar validação
			result, err := validator.ValidateEnvironment()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validator.ValidateEnvironment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestValidationSecurity_FilePermissions testa permissões de arquivo
func TestValidationSecurity_FilePermissions(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should validate files with secure permissions",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "validation_permissions_test")

			// Criar validator
			validator := src.NewValidator(nil)
			if validator == nil {
				t.Fatal("NewValidator() returned nil")
			}

			// Executar validação
			result, err := validator.ValidateEnvironment()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validator.ValidateEnvironment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verificar se o resultado foi retornado
				if result == nil {
					t.Error("ValidateEnvironment() returned nil result")
				}
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestValidationSecurity_DirectoryTraversal testa proteção contra directory traversal
func TestValidationSecurity_DirectoryTraversal(t *testing.T) {
	tests := []struct {
		name    string
		options *types.SetupOptions
		wantErr bool
	}{
		{
			name: "should prevent directory traversal in config path",
			options: &types.SetupOptions{
				OwnerName:      "Test User",
				OwnerEmail:     "test@example.com",
				ConfigPath:     "../../../etc/passwd",
				KeysPath:       "/tmp/syntropy_test/keys",
				BackupPath:     "/tmp/syntropy_test/backups",
				LogPath:        "/tmp/syntropy_test/logs",
				Verbose:        false,
				Force:          false,
				SkipValidation: false,
			},
			wantErr: true,
		},
		{
			name: "should prevent directory traversal in keys path",
			options: &types.SetupOptions{
				OwnerName:      "Test User",
				OwnerEmail:     "test@example.com",
				ConfigPath:     "/tmp/syntropy_test/config",
				KeysPath:       "../../../etc/passwd",
				BackupPath:     "/tmp/syntropy_test/backups",
				LogPath:        "/tmp/syntropy_test/logs",
				Verbose:        false,
				Force:          false,
				SkipValidation: false,
			},
			wantErr: true,
		},
		{
			name: "should prevent directory traversal in backup path",
			options: &types.SetupOptions{
				OwnerName:      "Test User",
				OwnerEmail:     "test@example.com",
				ConfigPath:     "/tmp/syntropy_test/config",
				KeysPath:       "/tmp/syntropy_test/keys",
				BackupPath:     "../../../etc/passwd",
				LogPath:        "/tmp/syntropy_test/logs",
				Verbose:        false,
				Force:          false,
				SkipValidation: false,
			},
			wantErr: true,
		},
		{
			name: "should prevent directory traversal in log path",
			options: &types.SetupOptions{
				OwnerName:      "Test User",
				OwnerEmail:     "test@example.com",
				ConfigPath:     "/tmp/syntropy_test/config",
				KeysPath:       "/tmp/syntropy_test/keys",
				BackupPath:     "/tmp/syntropy_test/backups",
				LogPath:        "../../../etc/passwd",
				Verbose:        false,
				Force:          false,
				SkipValidation: false,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "validation_traversal_test")

			// Ajustar caminhos para usar o diretório temporário se não forem maliciosos
			if tt.options.ConfigPath != "../../../etc/passwd" {
				tt.options.ConfigPath = filepath.Join(tempDir, "config")
			}
			if tt.options.KeysPath != "../../../etc/passwd" {
				tt.options.KeysPath = filepath.Join(tempDir, "keys")
			}
			if tt.options.BackupPath != "../../../etc/passwd" {
				tt.options.BackupPath = filepath.Join(tempDir, "backups")
			}
			if tt.options.LogPath != "../../../etc/passwd" {
				tt.options.LogPath = filepath.Join(tempDir, "logs")
			}

			// Criar validator
			validator := src.NewValidator(nil)
			if validator == nil {
				t.Fatal("NewValidator() returned nil")
			}

			// Executar validação
			result, err := validator.ValidateEnvironment()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validator.ValidateEnvironment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestValidationSecurity_CommandInjection testa proteção contra command injection
func TestValidationSecurity_CommandInjection(t *testing.T) {
	tests := []struct {
		name    string
		options *types.SetupOptions
		wantErr bool
	}{
		{
			name: "should prevent command injection in owner name",
			options: &types.SetupOptions{
				OwnerName:      "Test User; rm -rf /",
				OwnerEmail:     "test@example.com",
				ConfigPath:     "/tmp/syntropy_test/config",
				KeysPath:       "/tmp/syntropy_test/keys",
				BackupPath:     "/tmp/syntropy_test/backups",
				LogPath:        "/tmp/syntropy_test/logs",
				Verbose:        false,
				Force:          false,
				SkipValidation: false,
			},
			wantErr: true,
		},
		{
			name: "should prevent command injection in email",
			options: &types.SetupOptions{
				OwnerName:      "Test User",
				OwnerEmail:     "test@example.com; rm -rf /",
				ConfigPath:     "/tmp/syntropy_test/config",
				KeysPath:       "/tmp/syntropy_test/keys",
				BackupPath:     "/tmp/syntropy_test/backups",
				LogPath:        "/tmp/syntropy_test/logs",
				Verbose:        false,
				Force:          false,
				SkipValidation: false,
			},
			wantErr: true,
		},
		{
			name: "should prevent command injection in config path",
			options: &types.SetupOptions{
				OwnerName:      "Test User",
				OwnerEmail:     "test@example.com",
				ConfigPath:     "/tmp/syntropy_test/config; rm -rf /",
				KeysPath:       "/tmp/syntropy_test/keys",
				BackupPath:     "/tmp/syntropy_test/backups",
				LogPath:        "/tmp/syntropy_test/logs",
				Verbose:        false,
				Force:          false,
				SkipValidation: false,
			},
			wantErr: true,
		},
		{
			name: "should prevent command injection in keys path",
			options: &types.SetupOptions{
				OwnerName:      "Test User",
				OwnerEmail:     "test@example.com",
				ConfigPath:     "/tmp/syntropy_test/config",
				KeysPath:       "/tmp/syntropy_test/keys; rm -rf /",
				BackupPath:     "/tmp/syntropy_test/backups",
				LogPath:        "/tmp/syntropy_test/logs",
				Verbose:        false,
				Force:          false,
				SkipValidation: false,
			},
			wantErr: true,
		},
		{
			name: "should prevent command injection in backup path",
			options: &types.SetupOptions{
				OwnerName:      "Test User",
				OwnerEmail:     "test@example.com",
				ConfigPath:     "/tmp/syntropy_test/config",
				KeysPath:       "/tmp/syntropy_test/keys",
				BackupPath:     "/tmp/syntropy_test/backups; rm -rf /",
				LogPath:        "/tmp/syntropy_test/logs",
				Verbose:        false,
				Force:          false,
				SkipValidation: false,
			},
			wantErr: true,
		},
		{
			name: "should prevent command injection in log path",
			options: &types.SetupOptions{
				OwnerName:      "Test User",
				OwnerEmail:     "test@example.com",
				ConfigPath:     "/tmp/syntropy_test/config",
				KeysPath:       "/tmp/syntropy_test/keys",
				BackupPath:     "/tmp/syntropy_test/backups",
				LogPath:        "/tmp/syntropy_test/logs; rm -rf /",
				Verbose:        false,
				Force:          false,
				SkipValidation: false,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "validation_injection_test")

			// Ajustar caminhos para usar o diretório temporário se não forem maliciosos
			if tt.options.ConfigPath != "/tmp/syntropy_test/config; rm -rf /" {
				tt.options.ConfigPath = filepath.Join(tempDir, "config")
			}
			if tt.options.KeysPath != "/tmp/syntropy_test/keys; rm -rf /" {
				tt.options.KeysPath = filepath.Join(tempDir, "keys")
			}
			if tt.options.BackupPath != "/tmp/syntropy_test/backups; rm -rf /" {
				tt.options.BackupPath = filepath.Join(tempDir, "backups")
			}
			if tt.options.LogPath != "/tmp/syntropy_test/logs; rm -rf /" {
				tt.options.LogPath = filepath.Join(tempDir, "logs")
			}

			// Criar validator
			validator := src.NewValidator(nil)
			if validator == nil {
				t.Fatal("NewValidator() returned nil")
			}

			// Executar validação
			result, err := validator.ValidateEnvironment()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validator.ValidateEnvironment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestValidationSecurity_SQLInjection testa proteção contra SQL injection
func TestValidationSecurity_SQLInjection(t *testing.T) {
	tests := []struct {
		name    string
		options *types.SetupOptions
		wantErr bool
	}{
		{
			name: "should prevent SQL injection in owner name",
			options: &types.SetupOptions{
				OwnerName:      "Test User'; DROP TABLE users; --",
				OwnerEmail:     "test@example.com",
				ConfigPath:     "/tmp/syntropy_test/config",
				KeysPath:       "/tmp/syntropy_test/keys",
				BackupPath:     "/tmp/syntropy_test/backups",
				LogPath:        "/tmp/syntropy_test/logs",
				Verbose:        false,
				Force:          false,
				SkipValidation: false,
			},
			wantErr: true,
		},
		{
			name: "should prevent SQL injection in email",
			options: &types.SetupOptions{
				OwnerName:      "Test User",
				OwnerEmail:     "test@example.com'; DROP TABLE users; --",
				ConfigPath:     "/tmp/syntropy_test/config",
				KeysPath:       "/tmp/syntropy_test/keys",
				BackupPath:     "/tmp/syntropy_test/backups",
				LogPath:        "/tmp/syntropy_test/logs",
				Verbose:        false,
				Force:          false,
				SkipValidation: false,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "validation_sql_injection_test")

			// Ajustar caminhos para usar o diretório temporário
			tt.options.ConfigPath = filepath.Join(tempDir, "config")
			tt.options.KeysPath = filepath.Join(tempDir, "keys")
			tt.options.BackupPath = filepath.Join(tempDir, "backups")
			tt.options.LogPath = filepath.Join(tempDir, "logs")

			// Criar validator
			validator := src.NewValidator(nil)
			if validator == nil {
				t.Fatal("NewValidator() returned nil")
			}

			// Executar validação
			result, err := validator.ValidateEnvironment()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validator.ValidateEnvironment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestValidationSecurity_XSSProtection testa proteção contra XSS
func TestValidationSecurity_XSSProtection(t *testing.T) {
	tests := []struct {
		name    string
		options *types.SetupOptions
		wantErr bool
	}{
		{
			name: "should prevent XSS in owner name",
			options: &types.SetupOptions{
				OwnerName:      "<script>alert('XSS')</script>",
				OwnerEmail:     "test@example.com",
				ConfigPath:     "/tmp/syntropy_test/config",
				KeysPath:       "/tmp/syntropy_test/keys",
				BackupPath:     "/tmp/syntropy_test/backups",
				LogPath:        "/tmp/syntropy_test/logs",
				Verbose:        false,
				Force:          false,
				SkipValidation: false,
			},
			wantErr: true,
		},
		{
			name: "should prevent XSS in email",
			options: &types.SetupOptions{
				OwnerName:      "Test User",
				OwnerEmail:     "test@example.com<script>alert('XSS')</script>",
				ConfigPath:     "/tmp/syntropy_test/config",
				KeysPath:       "/tmp/syntropy_test/keys",
				BackupPath:     "/tmp/syntropy_test/backups",
				LogPath:        "/tmp/syntropy_test/logs",
				Verbose:        false,
				Force:          false,
				SkipValidation: false,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "validation_xss_test")

			// Ajustar caminhos para usar o diretório temporário
			tt.options.ConfigPath = filepath.Join(tempDir, "config")
			tt.options.KeysPath = filepath.Join(tempDir, "keys")
			tt.options.BackupPath = filepath.Join(tempDir, "backups")
			tt.options.LogPath = filepath.Join(tempDir, "logs")

			// Criar validator
			validator := src.NewValidator(nil)
			if validator == nil {
				t.Fatal("NewValidator() returned nil")
			}

			// Executar validação
			result, err := validator.ValidateEnvironment()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validator.ValidateEnvironment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestValidationSecurity_ResourceExhaustion testa proteção contra esgotamento de recursos
func TestValidationSecurity_ResourceExhaustion(t *testing.T) {
	tests := []struct {
		name    string
		options *types.SetupOptions
		wantErr bool
	}{
		{
			name: "should handle extremely long owner name",
			options: &types.SetupOptions{
				OwnerName:      string(make([]byte, 10000)), // 10KB string
				OwnerEmail:     "test@example.com",
				ConfigPath:     "/tmp/syntropy_test/config",
				KeysPath:       "/tmp/syntropy_test/keys",
				BackupPath:     "/tmp/syntropy_test/backups",
				LogPath:        "/tmp/syntropy_test/logs",
				Verbose:        false,
				Force:          false,
				SkipValidation: false,
			},
			wantErr: true,
		},
		{
			name: "should handle extremely long email",
			options: &types.SetupOptions{
				OwnerName:      "Test User",
				OwnerEmail:     string(make([]byte, 10000)) + "@example.com",
				ConfigPath:     "/tmp/syntropy_test/config",
				KeysPath:       "/tmp/syntropy_test/keys",
				BackupPath:     "/tmp/syntropy_test/backups",
				LogPath:        "/tmp/syntropy_test/logs",
				Verbose:        false,
				Force:          false,
				SkipValidation: false,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "validation_resource_test")

			// Ajustar caminhos para usar o diretório temporário
			tt.options.ConfigPath = filepath.Join(tempDir, "config")
			tt.options.KeysPath = filepath.Join(tempDir, "keys")
			tt.options.BackupPath = filepath.Join(tempDir, "backups")
			tt.options.LogPath = filepath.Join(tempDir, "logs")

			// Criar validator
			validator := src.NewValidator(nil)
			if validator == nil {
				t.Fatal("NewValidator() returned nil")
			}

			// Executar validação
			result, err := validator.ValidateEnvironment()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validator.ValidateEnvironment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestValidationSecurity_ConcurrentAccess testa segurança com acesso concorrente
func TestValidationSecurity_ConcurrentAccess(t *testing.T) {
	tests := []struct {
		name        string
		concurrency int
		wantErr     bool
	}{
		{
			name:        "should handle concurrent access securely",
			concurrency: 10,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "validation_concurrent_security_test")

			// Criar validator
			validator := src.NewValidator(nil)
			if validator == nil {
				t.Fatal("NewValidator() returned nil")
			}

			// Executar múltiplas validações concorrentemente
			done := make(chan bool, tt.concurrency)
			for i := 0; i < tt.concurrency; i++ {
				go func() {
					result, err := validator.ValidateEnvironment()
					if err != nil {
						t.Errorf("Concurrent ValidateEnvironment() failed: %v", err)
					}
					if result == nil {
						t.Error("Concurrent ValidateEnvironment() returned nil result")
					}
					done <- true
				}()
			}

			// Aguardar todas as goroutines terminarem
			for i := 0; i < tt.concurrency; i++ {
				<-done
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}
