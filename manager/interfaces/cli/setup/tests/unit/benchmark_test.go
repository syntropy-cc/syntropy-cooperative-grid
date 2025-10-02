//go:build !integration && !e2e && !performance && !security
// +build !integration,!e2e,!performance,!security

package unit

import (
	"errors"
	"os"
	"testing"

	setup "setup-component/src"
)

// BenchmarkSetupManager_Setup testa a performance do método Setup
func BenchmarkSetupManager_Setup(b *testing.B) {
	// Criar diretório temporário para testes
	tempDir := b.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	manager, err := setup.NewSetupManager()
	if err != nil {
		b.Fatalf("Failed to create setup manager: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		options := &setup.SetupOptions{
			Force:          true,
			SkipValidation: true,
			CustomSettings: map[string]string{
				"owner_name":  "Test User",
				"owner_email": "test@example.com",
			},
		}
		_ = manager.Setup(options)
	}
}

// BenchmarkSetupManager_Validate testa a performance do método Validate
func BenchmarkSetupManager_Validate(b *testing.B) {
	// Criar diretório temporário para testes
	tempDir := b.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	manager, err := setup.NewSetupManager()
	if err != nil {
		b.Fatalf("Failed to create setup manager: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = manager.Validate()
	}
}

// BenchmarkSetupManager_Status testa a performance do método Status
func BenchmarkSetupManager_Status(b *testing.B) {
	// Criar diretório temporário para testes
	tempDir := b.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	manager, err := setup.NewSetupManager()
	if err != nil {
		b.Fatalf("Failed to create setup manager: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = manager.Status()
	}
}

// BenchmarkSetupManager_Reset testa a performance do método Reset
func BenchmarkSetupManager_Reset(b *testing.B) {
	// Criar diretório temporário para testes
	tempDir := b.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	manager, err := setup.NewSetupManager()
	if err != nil {
		b.Fatalf("Failed to create setup manager: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = manager.Reset(true)
	}
}

// BenchmarkSetupManager_Repair testa a performance do método Repair
func BenchmarkSetupManager_Repair(b *testing.B) {
	// Criar diretório temporário para testes
	tempDir := b.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	manager, err := setup.NewSetupManager()
	if err != nil {
		b.Fatalf("Failed to create setup manager: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = manager.Repair()
	}
}

// BenchmarkSetupLogger_LogStep testa a performance do método LogStep
func BenchmarkSetupLogger_LogStep(b *testing.B) {
	logger := setup.NewSetupLogger()
	defer logger.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.LogStep("test_step", map[string]interface{}{
			"iteration": i,
			"message":   "Test message",
		})
	}
}

// BenchmarkSetupLogger_LogDebug testa a performance do método LogDebug
func BenchmarkSetupLogger_LogDebug(b *testing.B) {
	logger := setup.NewSetupLogger()
	defer logger.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.LogDebug("Test debug message", map[string]interface{}{
			"iteration": i,
		})
	}
}

// BenchmarkSetupLogger_LogError testa a performance do método LogError
func BenchmarkSetupLogger_LogError(b *testing.B) {
	logger := setup.NewSetupLogger()
	defer logger.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.LogError(errors.New("test error message"), map[string]interface{}{
			"iteration": i,
		})
	}
}

// BenchmarkSetupLogger_Close testa a performance do método Close
func BenchmarkSetupLogger_Close(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger := setup.NewSetupLogger()
		logger.Close()
	}
}
