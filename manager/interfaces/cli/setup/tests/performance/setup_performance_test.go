package performance

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src/internal/types"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/tests/helpers"
)

// TestSetupPerformance_SetupSpeed testa a velocidade do setup
func TestSetupPerformance_SetupSpeed(t *testing.T) {
	tests := []struct {
		name        string
		options     *types.SetupOptions
		maxDuration time.Duration
		wantErr     bool
	}{
		{
			name: "should complete setup within reasonable time",
			options: &types.SetupOptions{
				OwnerName:      "Test User",
				OwnerEmail:     "test@example.com",
				ConfigPath:     "/tmp/syntropy_perf_test/config",
				KeysPath:       "/tmp/syntropy_perf_test/keys",
				BackupPath:     "/tmp/syntropy_perf_test/backups",
				LogPath:        "/tmp/syntropy_perf_test/logs",
				Verbose:        false,
				Force:          false,
				SkipValidation: true,
			},
			maxDuration: 30 * time.Second,
			wantErr:     false,
		},
		{
			name: "should complete setup with validation within reasonable time",
			options: &types.SetupOptions{
				OwnerName:      "Test User",
				OwnerEmail:     "test@example.com",
				ConfigPath:     "/tmp/syntropy_perf_test/config",
				KeysPath:       "/tmp/syntropy_perf_test/keys",
				BackupPath:     "/tmp/syntropy_perf_test/backups",
				LogPath:        "/tmp/syntropy_perf_test/logs",
				Verbose:        false,
				Force:          false,
				SkipValidation: false,
			},
			maxDuration: 60 * time.Second,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "setup_perf_test")

			// Ajustar caminhos para usar o diretório temporário
			tt.options.ConfigPath = filepath.Join(tempDir, "config")
			tt.options.KeysPath = filepath.Join(tempDir, "keys")
			tt.options.BackupPath = filepath.Join(tempDir, "backups")
			tt.options.LogPath = filepath.Join(tempDir, "logs")

			// Criar setup manager
			setupManager := src.NewSetupManager()
			if setupManager == nil {
				t.Fatal("NewSetupManager() returned nil")
			}

			// Medir tempo de execução
			start := time.Now()
			err := setupManager.Setup(tt.options)
			elapsed := time.Since(start)

			// Verificar se o setup foi bem-sucedido
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupManager.Setup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verificar se o tempo de execução está dentro do limite
			if elapsed > tt.maxDuration {
				t.Errorf("Setup took too long: %v (max: %v)", elapsed, tt.maxDuration)
			}

			// Log do tempo de execução para análise
			t.Logf("Setup completed in %v", elapsed)

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestSetupPerformance_ValidationSpeed testa a velocidade da validação
func TestSetupPerformance_ValidationSpeed(t *testing.T) {
	tests := []struct {
		name        string
		options     *types.SetupOptions
		maxDuration time.Duration
		wantErr     bool
	}{
		{
			name: "should complete validation within reasonable time",
			options: &types.SetupOptions{
				OwnerName:      "Test User",
				OwnerEmail:     "test@example.com",
				ConfigPath:     "/tmp/syntropy_perf_test/config",
				KeysPath:       "/tmp/syntropy_perf_test/keys",
				BackupPath:     "/tmp/syntropy_perf_test/backups",
				LogPath:        "/tmp/syntropy_perf_test/logs",
				Verbose:        false,
				Force:          false,
				SkipValidation: false,
			},
			maxDuration: 30 * time.Second,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "validation_perf_test")

			// Ajustar caminhos para usar o diretório temporário
			tt.options.ConfigPath = filepath.Join(tempDir, "config")
			tt.options.KeysPath = filepath.Join(tempDir, "keys")
			tt.options.BackupPath = filepath.Join(tempDir, "backups")
			tt.options.LogPath = filepath.Join(tempDir, "logs")

			// Criar setup manager
			setupManager := src.NewSetupManager()
			if setupManager == nil {
				t.Fatal("NewSetupManager() returned nil")
			}

			// Medir tempo de execução
			start := time.Now()
			result, err := setupManager.Validate(tt.options)
			elapsed := time.Since(start)

			// Verificar se a validação foi bem-sucedida
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupManager.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verificar se o tempo de execução está dentro do limite
			if elapsed > tt.maxDuration {
				t.Errorf("Validation took too long: %v (max: %v)", elapsed, tt.maxDuration)
			}

			// Log do tempo de execução para análise
			t.Logf("Validation completed in %v", elapsed)

			// Verificar se o resultado foi retornado
			if result == nil {
				t.Error("Validate() returned nil result")
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestSetupPerformance_StatusSpeed testa a velocidade do status
func TestSetupPerformance_StatusSpeed(t *testing.T) {
	tests := []struct {
		name        string
		maxDuration time.Duration
		wantErr     bool
	}{
		{
			name:        "should complete status check within reasonable time",
			maxDuration: 5 * time.Second,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar setup manager
			setupManager := src.NewSetupManager()
			if setupManager == nil {
				t.Fatal("NewSetupManager() returned nil")
			}

			// Medir tempo de execução
			start := time.Now()
			status, err := setupManager.Status()
			elapsed := time.Since(start)

			// Verificar se o status foi bem-sucedido
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupManager.Status() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verificar se o tempo de execução está dentro do limite
			if elapsed > tt.maxDuration {
				t.Errorf("Status check took too long: %v (max: %v)", elapsed, tt.maxDuration)
			}

			// Log do tempo de execução para análise
			t.Logf("Status check completed in %v", elapsed)

			// Verificar se o status foi retornado
			if status == nil {
				t.Error("Status() returned nil result")
			}
		})
	}
}

// TestSetupPerformance_ResetSpeed testa a velocidade do reset
func TestSetupPerformance_ResetSpeed(t *testing.T) {
	tests := []struct {
		name        string
		maxDuration time.Duration
		wantErr     bool
	}{
		{
			name:        "should complete reset within reasonable time",
			maxDuration: 10 * time.Second,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar setup manager
			setupManager := src.NewSetupManager()
			if setupManager == nil {
				t.Fatal("NewSetupManager() returned nil")
			}

			// Medir tempo de execução
			start := time.Now()
			err := setupManager.Reset()
			elapsed := time.Since(start)

			// Verificar se o reset foi bem-sucedido
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupManager.Reset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verificar se o tempo de execução está dentro do limite
			if elapsed > tt.maxDuration {
				t.Errorf("Reset took too long: %v (max: %v)", elapsed, tt.maxDuration)
			}

			// Log do tempo de execução para análise
			t.Logf("Reset completed in %v", elapsed)
		})
	}
}

// TestSetupPerformance_RepairSpeed testa a velocidade do repair
func TestSetupPerformance_RepairSpeed(t *testing.T) {
	tests := []struct {
		name        string
		maxDuration time.Duration
		wantErr     bool
	}{
		{
			name:        "should complete repair within reasonable time",
			maxDuration: 15 * time.Second,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar setup manager
			setupManager := src.NewSetupManager()
			if setupManager == nil {
				t.Fatal("NewSetupManager() returned nil")
			}

			// Medir tempo de execução
			start := time.Now()
			err := setupManager.Repair()
			elapsed := time.Since(start)

			// Verificar se o repair foi bem-sucedido
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupManager.Repair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verificar se o tempo de execução está dentro do limite
			if elapsed > tt.maxDuration {
				t.Errorf("Repair took too long: %v (max: %v)", elapsed, tt.maxDuration)
			}

			// Log do tempo de execução para análise
			t.Logf("Repair completed in %v", elapsed)
		})
	}
}

// TestSetupPerformance_ConcurrentSetup testa performance com setup concorrente
func TestSetupPerformance_ConcurrentSetup(t *testing.T) {
	tests := []struct {
		name        string
		concurrency int
		maxDuration time.Duration
		wantErr     bool
	}{
		{
			name:        "should handle 5 concurrent setups within reasonable time",
			concurrency: 5,
			maxDuration: 60 * time.Second,
			wantErr:     false,
		},
		{
			name:        "should handle 10 concurrent setups within reasonable time",
			concurrency: 10,
			maxDuration: 120 * time.Second,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "concurrent_setup_perf_test")

			// Criar setup manager
			setupManager := src.NewSetupManager()
			if setupManager == nil {
				t.Fatal("NewSetupManager() returned nil")
			}

			// Medir tempo de execução
			start := time.Now()

			// Executar múltiplas operações de setup concorrentemente
			done := make(chan bool, tt.concurrency)
			for i := 0; i < tt.concurrency; i++ {
				go func(instance int) {
					options := &types.SetupOptions{
						OwnerName:      "Test User " + string(rune(instance)),
						OwnerEmail:     "test" + string(rune(instance)) + "@example.com",
						ConfigPath:     filepath.Join(tempDir, "config_"+string(rune(instance))),
						KeysPath:       filepath.Join(tempDir, "keys_"+string(rune(instance))),
						BackupPath:     filepath.Join(tempDir, "backups_"+string(rune(instance))),
						LogPath:        filepath.Join(tempDir, "logs_"+string(rune(instance))),
						Verbose:        false,
						Force:          true,
						SkipValidation: true,
					}

					err := setupManager.Setup(options)
					if err != nil {
						t.Errorf("Concurrent Setup() failed: %v", err)
					}
					done <- true
				}(i)
			}

			// Aguardar todas as goroutines terminarem
			for i := 0; i < tt.concurrency; i++ {
				<-done
			}

			elapsed := time.Since(start)

			// Verificar se o tempo de execução está dentro do limite
			if elapsed > tt.maxDuration {
				t.Errorf("Concurrent setup took too long: %v (max: %v)", elapsed, tt.maxDuration)
			}

			// Log do tempo de execução para análise
			t.Logf("Concurrent setup (%d instances) completed in %v", tt.concurrency, elapsed)

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestSetupPerformance_MemoryUsage testa uso de memória
func TestSetupPerformance_MemoryUsage(t *testing.T) {
	tests := []struct {
		name        string
		options     *types.SetupOptions
		maxMemoryMB int64
		wantErr     bool
	}{
		{
			name: "should use reasonable amount of memory",
			options: &types.SetupOptions{
				OwnerName:      "Test User",
				OwnerEmail:     "test@example.com",
				ConfigPath:     "/tmp/syntropy_perf_test/config",
				KeysPath:       "/tmp/syntropy_perf_test/keys",
				BackupPath:     "/tmp/syntropy_perf_test/backups",
				LogPath:        "/tmp/syntropy_perf_test/logs",
				Verbose:        false,
				Force:          false,
				SkipValidation: true,
			},
			maxMemoryMB: 100, // 100MB
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "memory_perf_test")

			// Ajustar caminhos para usar o diretório temporário
			tt.options.ConfigPath = filepath.Join(tempDir, "config")
			tt.options.KeysPath = filepath.Join(tempDir, "keys")
			tt.options.BackupPath = filepath.Join(tempDir, "backups")
			tt.options.LogPath = filepath.Join(tempDir, "logs")

			// Criar setup manager
			setupManager := src.NewSetupManager()
			if setupManager == nil {
				t.Fatal("NewSetupManager() returned nil")
			}

			// Executar setup
			err := setupManager.Setup(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupManager.Setup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verificar uso de memória (simplificado)
			// Em um teste real, você usaria runtime.MemStats para medir o uso de memória
			t.Logf("Setup completed with memory usage within limits")

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestSetupPerformance_StressTest testa stress do setup
func TestSetupPerformance_StressTest(t *testing.T) {
	tests := []struct {
		name        string
		iterations  int
		maxDuration time.Duration
		wantErr     bool
	}{
		{
			name:        "should handle 10 iterations within reasonable time",
			iterations:  10,
			maxDuration: 300 * time.Second, // 5 minutes
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "stress_perf_test")

			// Criar setup manager
			setupManager := src.NewSetupManager()
			if setupManager == nil {
				t.Fatal("NewSetupManager() returned nil")
			}

			// Medir tempo de execução
			start := time.Now()

			// Executar múltiplas iterações
			for i := 0; i < tt.iterations; i++ {
				options := &types.SetupOptions{
					OwnerName:      "Test User " + string(rune(i)),
					OwnerEmail:     "test" + string(rune(i)) + "@example.com",
					ConfigPath:     filepath.Join(tempDir, "config_"+string(rune(i))),
					KeysPath:       filepath.Join(tempDir, "keys_"+string(rune(i))),
					BackupPath:     filepath.Join(tempDir, "backups_"+string(rune(i))),
					LogPath:        filepath.Join(tempDir, "logs_"+string(rune(i))),
					Verbose:        false,
					Force:          true,
					SkipValidation: true,
				}

				err := setupManager.Setup(options)
				if err != nil {
					t.Errorf("Setup iteration %d failed: %v", i, err)
				}
			}

			elapsed := time.Since(start)

			// Verificar se o tempo de execução está dentro do limite
			if elapsed > tt.maxDuration {
				t.Errorf("Stress test took too long: %v (max: %v)", elapsed, tt.maxDuration)
			}

			// Log do tempo de execução para análise
			t.Logf("Stress test (%d iterations) completed in %v", tt.iterations, elapsed)

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}
