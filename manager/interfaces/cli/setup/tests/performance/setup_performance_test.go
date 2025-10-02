//go:build performance
// +build performance

package performance

import (
	"os"
	"testing"
	"time"

	setup "setup-component/src"
)

// TestSetupManager_Performance testa a performance do SetupManager
func TestSetupManager_Performance(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	manager := setup.NewSetupManager(logger)

	tests := []struct {
		name        string
		options     *setup.SetupOptions
		maxDuration time.Duration
	}{
		{
			name: "should complete setup within acceptable time",
			options: &types.setup.SetupOptions{
				Force:          false,
				SkipValidation: false,
			},
			maxDuration: 30 * time.Second,
		},
		{
			name: "should complete setup with force within acceptable time",
			options: &types.setup.SetupOptions{
				Force:          true,
				SkipValidation: false,
			},
			maxDuration: 30 * time.Second,
		},
		{
			name: "should complete setup skipping validation within acceptable time",
			options: &types.setup.SetupOptions{
				Force:          false,
				SkipValidation: true,
			},
			maxDuration: 20 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()

			err := manager.Setup(tt.options)
			if err != nil {
				t.Errorf("SetupManager.Setup() error = %v", err)
				return
			}

			duration := time.Since(start)
			if duration > tt.maxDuration {
				t.Errorf("Setup took %v, which exceeds maximum duration of %v", duration, tt.maxDuration)
			}

			t.Logf("Setup completed in %v", duration)
		})
	}
}

// TestSetupManager_Validation_Performance testa a performance da validação
func TestSetupManager_Validation_Performance(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	manager := setup.NewSetupManager(logger)

	tests := []struct {
		name        string
		maxDuration time.Duration
	}{
		{
			name:        "should complete validation within acceptable time",
			maxDuration: 10 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()

			err := manager.Validate()
			if err != nil {
				t.Errorf("SetupManager.Validate() error = %v", err)
				return
			}

			duration := time.Since(start)
			if duration > tt.maxDuration {
				t.Errorf("Validation took %v, which exceeds maximum duration of %v", duration, tt.maxDuration)
			}

			t.Logf("Validation completed in %v", duration)
		})
	}
}

// TestSetupManager_Status_Performance testa a performance do status
func TestSetupManager_Status_Performance(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	manager := setup.NewSetupManager(logger)

	tests := []struct {
		name        string
		maxDuration time.Duration
	}{
		{
			name:        "should get status within acceptable time",
			maxDuration: 1 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()

			status, err := manager.Status()
			if err != nil {
				t.Errorf("SetupManager.Status() error = %v", err)
				return
			}

			duration := time.Since(start)
			if duration > tt.maxDuration {
				t.Errorf("Status took %v, which exceeds maximum duration of %v", duration, tt.maxDuration)
			}

			t.Logf("Status completed in %v", duration)
			t.Logf("Status: %+v", status)
		})
	}
}

// TestSetupManager_Reset_Performance testa a performance do reset
func TestSetupManager_Reset_Performance(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	manager := setup.NewSetupManager(logger)

	tests := []struct {
		name        string
		setup       bool
		maxDuration time.Duration
	}{
		{
			name:        "should reset within acceptable time when no setup exists",
			setup:       false,
			maxDuration: 5 * time.Second,
		},
		{
			name:        "should reset within acceptable time when setup exists",
			setup:       true,
			maxDuration: 10 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Executar setup se necessário
			if tt.setup {
				options := &types.setup.SetupOptions{
					Force:          true,
					SkipValidation: false,
				}
				err := manager.Setup(options)
				if err != nil {
					t.Fatalf("Failed to setup: %v", err)
				}
			}

			start := time.Now()

			err := manager.Reset()
			if err != nil {
				t.Errorf("SetupManager.Reset() error = %v", err)
				return
			}

			duration := time.Since(start)
			if duration > tt.maxDuration {
				t.Errorf("Reset took %v, which exceeds maximum duration of %v", duration, tt.maxDuration)
			}

			t.Logf("Reset completed in %v", duration)
		})
	}
}

// TestSetupManager_Repair_Performance testa a performance do repair
func TestSetupManager_Repair_Performance(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	manager := setup.NewSetupManager(logger)

	tests := []struct {
		name        string
		setup       bool
		maxDuration time.Duration
	}{
		{
			name:        "should repair within acceptable time when no setup exists",
			setup:       false,
			maxDuration: 5 * time.Second,
		},
		{
			name:        "should repair within acceptable time when setup exists",
			setup:       true,
			maxDuration: 15 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Executar setup se necessário
			if tt.setup {
				options := &types.setup.SetupOptions{
					Force:          true,
					SkipValidation: false,
				}
				err := manager.Setup(options)
				if err != nil {
					t.Fatalf("Failed to setup: %v", err)
				}
			}

			start := time.Now()

			err := manager.Repair()
			if err != nil {
				t.Errorf("SetupManager.Repair() error = %v", err)
				return
			}

			duration := time.Since(start)
			if duration > tt.maxDuration {
				t.Errorf("Repair took %v, which exceeds maximum duration of %v", duration, tt.maxDuration)
			}

			t.Logf("Repair completed in %v", duration)
		})
	}
}

// TestSetupManager_Legacy_Performance testa a performance das funções legacy
func TestSetupManager_Legacy_Performance(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	tests := []struct {
		name        string
		maxDuration time.Duration
	}{
		{
			name:        "should complete legacy setup within acceptable time",
			maxDuration: 30 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()

			err := setup.SetupLegacy()
			if err != nil {
				t.Errorf("SetupLegacy() error = %v", err)
				return
			}

			duration := time.Since(start)
			if duration > tt.maxDuration {
				t.Errorf("Legacy setup took %v, which exceeds maximum duration of %v", duration, tt.maxDuration)
			}

			t.Logf("Legacy setup completed in %v", duration)
		})
	}
}

// TestSetupManager_StatusLegacy_Performance testa a performance do status legacy
func TestSetupManager_StatusLegacy_Performance(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	tests := []struct {
		name        string
		setup       bool
		maxDuration time.Duration
	}{
		{
			name:        "should get legacy status within acceptable time when no setup exists",
			setup:       false,
			maxDuration: 1 * time.Second,
		},
		{
			name:        "should get legacy status within acceptable time when setup exists",
			setup:       true,
			maxDuration: 1 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Executar setup se necessário
			if tt.setup {
				err := setup.SetupLegacy()
				if err != nil {
					t.Fatalf("Failed to setup legacy: %v", err)
				}
			}

			start := time.Now()

			status, err := setup.StatusLegacy()
			if err != nil {
				t.Errorf("StatusLegacy() error = %v", err)
				return
			}

			duration := time.Since(start)
			if duration > tt.maxDuration {
				t.Errorf("Legacy status took %v, which exceeds maximum duration of %v", duration, tt.maxDuration)
			}

			t.Logf("Legacy status completed in %v", duration)
			t.Logf("Legacy status: %+v", status)
		})
	}
}

// TestSetupManager_ResetLegacy_Performance testa a performance do reset legacy
func TestSetupManager_ResetLegacy_Performance(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	tests := []struct {
		name        string
		setup       bool
		maxDuration time.Duration
	}{
		{
			name:        "should reset legacy within acceptable time when no setup exists",
			setup:       false,
			maxDuration: 5 * time.Second,
		},
		{
			name:        "should reset legacy within acceptable time when setup exists",
			setup:       true,
			maxDuration: 10 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Executar setup se necessário
			if tt.setup {
				err := setup.SetupLegacy()
				if err != nil {
					t.Fatalf("Failed to setup legacy: %v", err)
				}
			}

			start := time.Now()

			err := setup.ResetLegacy()
			if err != nil {
				t.Errorf("ResetLegacy() error = %v", err)
				return
			}

			duration := time.Since(start)
			if duration > tt.maxDuration {
				t.Errorf("Legacy reset took %v, which exceeds maximum duration of %v", duration, tt.maxDuration)
			}

			t.Logf("Legacy reset completed in %v", duration)
		})
	}
}

// TestSetupManager_Concurrent_Performance testa a performance com operações concorrentes
func TestSetupManager_Concurrent_Performance(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	manager := setup.NewSetupManager(logger)

	tests := []struct {
		name        string
		maxDuration time.Duration
	}{
		{
			name:        "should handle concurrent operations within acceptable time",
			maxDuration: 60 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()

			// Executar operações concorrentes
			done := make(chan bool, 3)

			// Setup
			go func() {
				options := &types.setup.SetupOptions{
					Force:          true,
					SkipValidation: false,
				}
				_ = manager.Setup(options)
				done <- true
			}()

			// Status
			go func() {
				_, _ = manager.Status()
				done <- true
			}()

			// Validate
			go func() {
				_ = manager.Validate()
				done <- true
			}()

			// Aguardar todas as operações
			for i := 0; i < 3; i++ {
				<-done
			}

			duration := time.Since(start)
			if duration > tt.maxDuration {
				t.Errorf("Concurrent operations took %v, which exceeds maximum duration of %v", duration, tt.maxDuration)
			}

			t.Logf("Concurrent operations completed in %v", duration)
		})
	}
}
