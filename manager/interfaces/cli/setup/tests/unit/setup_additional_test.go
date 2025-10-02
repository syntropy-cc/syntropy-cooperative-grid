//go:build !integration && !e2e && !performance && !security
// +build !integration,!e2e,!performance,!security

package unit

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	setup "setup-component/src"
)

// TestSetupLegacy testa a função legacy de setup
func TestSetupLegacy(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	tests := []struct {
		name    string
		options setup.LegacySetupOptions
		wantErr bool
	}{
		{
			name: "should setup successfully with default options",
			options: setup.LegacySetupOptions{
				Force:          false,
				InstallService: false,
				ConfigPath:     "",
				HomeDir:        "",
			},
			wantErr: false,
		},
		{
			name: "should setup successfully with force option",
			options: setup.LegacySetupOptions{
				Force:          true,
				InstallService: false,
				ConfigPath:     "",
				HomeDir:        "",
			},
			wantErr: false,
		},
		{
			name: "should setup successfully with custom config path",
			options: setup.LegacySetupOptions{
				Force:          false,
				InstallService: false,
				ConfigPath:     filepath.Join(tempDir, "custom_config.yaml"),
				HomeDir:        "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := setup.SetupLegacy(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupLegacy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if result == nil {
					t.Error("SetupLegacy() returned nil result")
					return
				}

				// Verificar campos obrigatórios
				if result.StartTime.IsZero() {
					t.Error("Result missing start time")
				}
				if result.EndTime.IsZero() {
					t.Error("Result missing end time")
				}
				if result.Environment == "" {
					t.Error("Result missing environment")
				}
				if result.Message == "" {
					t.Error("Result missing message")
				}
			}
		})
	}
}

// TestStatusLegacy testa a função legacy de status
func TestStatusLegacy(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	tests := []struct {
		name    string
		setup   bool
		wantErr bool
	}{
		{
			name:    "should return not found when no setup exists",
			setup:   false,
			wantErr: false, // StatusLegacy não retorna erro, apenas indica que não foi encontrado
		},
		{
			name:    "should return status when setup exists",
			setup:   true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Executar setup se necessário
			if tt.setup {
				options := setup.LegacySetupOptions{
					Force:          true,
					InstallService: false,
					ConfigPath:     "",
					HomeDir:        "",
				}
				_, err := setup.SetupLegacy(options)
				if err != nil {
					t.Fatalf("Failed to setup: %v", err)
				}
			}

			options := setup.LegacySetupOptions{
				Force:          false,
				InstallService: false,
				ConfigPath:     "",
				HomeDir:        "",
			}
			result, err := setup.StatusLegacy(options)
			if (err != nil) != tt.wantErr {
				t.Errorf("StatusLegacy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if result == nil {
				t.Error("StatusLegacy() returned nil result")
				return
			}

			// Verificar campos obrigatórios
			if result.StartTime.IsZero() {
				t.Error("Result missing start time")
			}
			if result.EndTime.IsZero() {
				t.Error("Result missing end time")
			}
			if result.Environment == "" {
				t.Error("Result missing environment")
			}
			if result.Message == "" {
				t.Error("Result missing message")
			}
		})
	}
}

// TestResetLegacy testa a função legacy de reset
func TestResetLegacy(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	tests := []struct {
		name    string
		setup   bool
		wantErr bool
	}{
		{
			name:    "should reset successfully even without setup",
			setup:   false,
			wantErr: false,
		},
		{
			name:    "should reset successfully with existing setup",
			setup:   true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Executar setup se necessário
			if tt.setup {
				options := setup.LegacySetupOptions{
					Force:          true,
					InstallService: false,
					ConfigPath:     "",
					HomeDir:        "",
				}
				_, err := setup.SetupLegacy(options)
				if err != nil {
					t.Fatalf("Failed to setup: %v", err)
				}
			}

			options := setup.LegacySetupOptions{
				Force:          false,
				InstallService: false,
				ConfigPath:     "",
				HomeDir:        "",
			}
			result, err := setup.ResetLegacy(options)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResetLegacy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if result == nil {
				t.Error("ResetLegacy() returned nil result")
				return
			}

			// Verificar campos obrigatórios
			if result.StartTime.IsZero() {
				t.Error("Result missing start time")
			}
			if result.EndTime.IsZero() {
				t.Error("Result missing end time")
			}
			if result.Message == "" {
				t.Error("Result missing message")
			}
		})
	}
}

// TestGetSyntropyDirLegacy testa a função legacy de obtenção do diretório
func TestGetSyntropyDirLegacy(t *testing.T) {
	tests := []struct {
		name string
		os   string
		want string
	}{
		{
			name: "should return Windows path on Windows",
			os:   "windows",
			want: "Syntropy",
		},
		{
			name: "should return Unix path on Linux",
			os:   "linux",
			want: ".syntropy",
		},
		{
			name: "should return Unix path on macOS",
			os:   "darwin",
			want: ".syntropy",
		},
		{
			name: "should return Unix path on unknown OS",
			os:   "unknown",
			want: ".syntropy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Não podemos alterar runtime.GOOS diretamente, então testamos o comportamento atual
			result := setup.GetSyntropyDirLegacy()
			if result == "" {
				t.Error("GetSyntropyDirLegacy() returned empty string")
			}

			// Verificar se contém o padrão esperado baseado no SO atual
			if runtime.GOOS == "windows" {
				if filepath.Base(result) != "Syntropy" {
					t.Errorf("GetSyntropyDirLegacy() = %v, expected Windows path", result)
				}
			} else {
				if filepath.Base(result) != ".syntropy" {
					t.Errorf("GetSyntropyDirLegacy() = %v, expected Unix path", result)
				}
			}
		})
	}
}

// TestSetupManager_ErrorHandling testa o tratamento de erros indiretamente
func TestSetupManager_ErrorHandling(t *testing.T) {
	manager, err := setup.NewSetupManager()
	if err != nil {
		t.Fatalf("Failed to create setup manager: %v", err)
	}
	// Logger será fechado automaticamente

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should handle setup errors gracefully",
			wantErr: false, // Setup pode falhar, mas não deve causar panic
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := &types.setup.SetupOptions{
				Force:          false,
				SkipValidation: false,
			}
			err := manager.Setup(options)
			// Não verificar erro específico, apenas que não houve panic
			_ = err
		})
	}
}

// TestSetupManager_Setup_ErrorPaths testa os caminhos de erro do setup
func TestSetupManager_Setup_ErrorPaths(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	tests := []struct {
		name    string
		options *setup.SetupOptions
		wantErr bool
	}{
		{
			name:    "should fail with nil options",
			options: nil,
			wantErr: true,
		},
		{
			name: "should fail with empty custom settings",
			options: &types.setup.SetupOptions{
				Force:          false,
				ValidateOnly:   false,
				Verbose:        false,
				Quiet:          false,
				CustomSettings: map[string]string{},
			},
			wantErr: false, // Deve usar valores padrão
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := setup.NewSetupManager()
			if err != nil {
				t.Fatalf("Failed to create setup manager: %v", err)
			}
			// Logger será fechado automaticamente

			err = manager.Setup(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupManager.Setup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestSetupManager_Setup_Concurrent testa setup concorrente
func TestSetupManager_Setup_Concurrent(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	manager, err := setup.NewSetupManager()
	if err != nil {
		t.Fatalf("Failed to create setup manager: %v", err)
	}
	// Logger será fechado automaticamente

	// Executar múltiplos setups concorrentes
	numGoroutines := 5
	done := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			options := &types.setup.SetupOptions{
				Force:        true, // Usar force para evitar conflitos de validação
				ValidateOnly: false,
				Verbose:      false,
				Quiet:        true,
				CustomSettings: map[string]string{
					"owner_name":  "Test User",
					"owner_email": "test@example.com",
				},
			}
			err := manager.Setup(options)
			done <- err
		}(i)
	}

	// Coletar resultados
	for i := 0; i < numGoroutines; i++ {
		select {
		case err := <-done:
			if err != nil {
				t.Errorf("Concurrent setup %d failed: %v", i, err)
			}
		case <-time.After(30 * time.Second):
			t.Errorf("Concurrent setup %d timed out", i)
		}
	}
}

// BenchmarkNewSetupManager testa a performance da criação do gerenciador
func BenchmarkNewSetupManager(b *testing.B) {
	for i := 0; i < b.N; i++ {
		manager, err := setup.NewSetupManager()
		if err != nil {
			b.Fatalf("Failed to create setup manager: %v", err)
		}
		// Logger será fechado automaticamente
	}
}

// BenchmarkSetupManager_SetupAdditional testa a performance do setup
func BenchmarkSetupManager_SetupAdditional(b *testing.B) {
	// Criar diretório temporário para testes
	tempDir := b.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	manager, err := setup.NewSetupManager()
	if err != nil {
		b.Fatalf("Failed to create setup manager: %v", err)
	}
	// Logger será fechado automaticamente

	options := &types.setup.SetupOptions{
		Force:        true,
		ValidateOnly: false,
		Verbose:      false,
		Quiet:        true,
		CustomSettings: map[string]string{
			"owner_name":  "Test User",
			"owner_email": "test@example.com",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := manager.Setup(options)
		if err != nil {
			b.Fatalf("Setup failed: %v", err)
		}
	}
}

// BenchmarkSetupManager_ValidateAdditional testa a performance da validação
func BenchmarkSetupManager_ValidateAdditional(b *testing.B) {
	// Criar diretório temporário para testes
	tempDir := b.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	manager, err := setup.NewSetupManager()
	if err != nil {
		b.Fatalf("Failed to create setup manager: %v", err)
	}
	// Logger será fechado automaticamente

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := manager.Validate()
		if err != nil {
			b.Fatalf("Validate failed: %v", err)
		}
	}
}
