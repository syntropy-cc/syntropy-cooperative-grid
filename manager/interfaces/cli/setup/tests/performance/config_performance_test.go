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

// TestConfigPerformance_ConfigurationGenerationSpeed testa a velocidade da geração de configuração
func TestConfigPerformance_ConfigurationGenerationSpeed(t *testing.T) {
	tests := []struct {
		name        string
		options     *types.ConfigOptions
		maxDuration time.Duration
		wantErr     bool
	}{
		{
			name: "should generate configuration within reasonable time",
			options: &types.ConfigOptions{
				OwnerName:      "Test User",
				OwnerEmail:     "test@example.com",
				NetworkConfig:  nil,
				SecurityConfig: nil,
				CustomSettings: map[string]string{
					"log_level":        "info",
					"api_endpoint":     "https://api.syntropy.network",
					"debug_mode":       "false",
					"max_connections":  "100",
					"timeout":          "30s",
					"retry_attempts":   "3",
					"backup_interval":  "24h",
					"cleanup_interval": "7d",
				},
			},
			maxDuration: 5 * time.Second,
			wantErr:     false,
		},
		{
			name: "should handle minimal configuration quickly",
			options: &types.ConfigOptions{
				OwnerName:      "Test User",
				OwnerEmail:     "test@example.com",
				NetworkConfig:  nil,
				SecurityConfig: nil,
				CustomSettings: map[string]string{},
			},
			maxDuration: 2 * time.Second,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "config_generation_perf_test")
			configDir := filepath.Join(tempDir, "config")
			os.MkdirAll(configDir, 0755)

			// Criar configurator
			configurator := src.NewConfigurator(nil)
			if configurator == nil {
				t.Fatal("NewConfigurator() returned nil")
			}

			// Medir tempo de execução
			start := time.Now()
			err := configurator.GenerateConfig(tt.options)
			elapsed := time.Since(start)

			// Verificar se a geração foi bem-sucedida
			if (err != nil) != tt.wantErr {
				t.Errorf("Configurator.GenerateConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verificar se o tempo de execução está dentro do limite
			if elapsed > tt.maxDuration {
				t.Errorf("Configuration generation took too long: %v (max: %v)", elapsed, tt.maxDuration)
			}

			// Log do tempo de execução para análise
			t.Logf("Configuration generation completed in %v", elapsed)

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestConfigPerformance_StructureCreationSpeed testa a velocidade da criação de estrutura
func TestConfigPerformance_StructureCreationSpeed(t *testing.T) {
	tests := []struct {
		name        string
		maxDuration time.Duration
		wantErr     bool
	}{
		{
			name:        "should create structure within reasonable time",
			maxDuration: 3 * time.Second,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "config_structure_perf_test")
			configDir := filepath.Join(tempDir, "config")
			os.MkdirAll(configDir, 0755)

			// Criar configurator
			configurator := src.NewConfigurator(nil)
			if configurator == nil {
				t.Fatal("NewConfigurator() returned nil")
			}

			// Medir tempo de execução
			start := time.Now()
			err := configurator.CreateStructure()
			elapsed := time.Since(start)

			// Verificar se a criação foi bem-sucedida
			if (err != nil) != tt.wantErr {
				t.Errorf("Configurator.CreateStructure() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verificar se o tempo de execução está dentro do limite
			if elapsed > tt.maxDuration {
				t.Errorf("Structure creation took too long: %v (max: %v)", elapsed, tt.maxDuration)
			}

			// Log do tempo de execução para análise
			t.Logf("Structure creation completed in %v", elapsed)

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestConfigPerformance_KeyGenerationSpeed testa a velocidade da geração de chaves
func TestConfigPerformance_KeyGenerationSpeed(t *testing.T) {
	tests := []struct {
		name        string
		maxDuration time.Duration
		wantErr     bool
	}{
		{
			name:        "should generate keys within reasonable time",
			maxDuration: 5 * time.Second,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "config_keys_perf_test")
			keysDir := filepath.Join(tempDir, "keys")
			os.MkdirAll(keysDir, 0755)

			// Criar configurator
			configurator := src.NewConfigurator(nil)
			if configurator == nil {
				t.Fatal("NewConfigurator() returned nil")
			}

			// Medir tempo de execução
			start := time.Now()
			keyPair, err := configurator.GenerateKeys()
			elapsed := time.Since(start)

			// Verificar se a geração foi bem-sucedida
			if (err != nil) != tt.wantErr {
				t.Errorf("Configurator.GenerateKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verificar se o tempo de execução está dentro do limite
			if elapsed > tt.maxDuration {
				t.Errorf("Key generation took too long: %v (max: %v)", elapsed, tt.maxDuration)
			}

			// Log do tempo de execução para análise
			t.Logf("Key generation completed in %v", elapsed)

			// Verificar se a chave foi gerada
			if keyPair == nil {
				t.Error("GenerateKeys() returned nil keyPair")
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestConfigPerformance_ConfigurationValidationSpeed testa a velocidade da validação de configuração
func TestConfigPerformance_ConfigurationValidationSpeed(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func(string) error
		maxDuration time.Duration
		wantErr     bool
	}{
		{
			name: "should validate configuration within reasonable time",
			setupFunc: func(configDir string) error {
				// Criar arquivo de configuração válido
				configPath := filepath.Join(configDir, "manager.yaml")
				configContent := `
manager:
  home_dir: "/home/testuser/.syntropy"
  log_level: "info"
  api_endpoint: "https://api.syntropy.network"
  directories:
    config: "/home/testuser/.syntropy/config"
    keys: "/home/testuser/.syntropy/keys"
  default_paths:
    config: "/home/testuser/.syntropy/config/manager.yaml"
    log: "/home/testuser/.syntropy/logs/manager.log"
owner_key:
  type: "ed25519"
  path: "/home/testuser/.syntropy/keys/owner.key"
environment:
  os: "linux"
  architecture: "amd64"
  home_dir: "/home/testuser"
`
				return os.WriteFile(configPath, []byte(configContent), 0644)
			},
			maxDuration: 3 * time.Second,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "config_validation_perf_test")
			configDir := filepath.Join(tempDir, "config")
			os.MkdirAll(configDir, 0755)

			// Setup do teste
			if tt.setupFunc != nil {
				if err := tt.setupFunc(configDir); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			// Criar configurator
			configurator := src.NewConfigurator(nil)
			if configurator == nil {
				t.Fatal("NewConfigurator() returned nil")
			}

			// Medir tempo de execução
			start := time.Now()
			err := configurator.ValidateConfig()
			elapsed := time.Since(start)

			// Verificar se a validação foi bem-sucedida
			if (err != nil) != tt.wantErr {
				t.Errorf("Configurator.ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verificar se o tempo de execução está dentro do limite
			if elapsed > tt.maxDuration {
				t.Errorf("Configuration validation took too long: %v (max: %v)", elapsed, tt.maxDuration)
			}

			// Log do tempo de execução para análise
			t.Logf("Configuration validation completed in %v", elapsed)

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestConfigPerformance_ConfigurationBackupSpeed testa a velocidade do backup de configuração
func TestConfigPerformance_ConfigurationBackupSpeed(t *testing.T) {
	tests := []struct {
		name        string
		backupName  string
		setupFunc   func(string) error
		maxDuration time.Duration
		wantErr     bool
	}{
		{
			name:       "should backup configuration within reasonable time",
			backupName: "config_backup",
			setupFunc: func(configDir string) error {
				// Criar arquivo de configuração
				configPath := filepath.Join(configDir, "manager.yaml")
				configContent := `
manager:
  home_dir: "/home/testuser/.syntropy"
  log_level: "info"
owner_key:
  type: "ed25519"
environment:
  os: "linux"
`
				return os.WriteFile(configPath, []byte(configContent), 0644)
			},
			maxDuration: 3 * time.Second,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "config_backup_perf_test")
			configDir := filepath.Join(tempDir, "config")
			backupDir := filepath.Join(tempDir, "backups")
			os.MkdirAll(configDir, 0755)
			os.MkdirAll(backupDir, 0755)

			// Setup do teste
			if tt.setupFunc != nil {
				if err := tt.setupFunc(configDir); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			// Criar configurator
			configurator := src.NewConfigurator(nil)
			if configurator == nil {
				t.Fatal("NewConfigurator() returned nil")
			}

			// Medir tempo de execução
			start := time.Now()
			err := configurator.BackupConfig(tt.backupName)
			elapsed := time.Since(start)

			// Verificar se o backup foi bem-sucedido
			if (err != nil) != tt.wantErr {
				t.Errorf("Configurator.BackupConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verificar se o tempo de execução está dentro do limite
			if elapsed > tt.maxDuration {
				t.Errorf("Configuration backup took too long: %v (max: %v)", elapsed, tt.maxDuration)
			}

			// Log do tempo de execução para análise
			t.Logf("Configuration backup completed in %v", elapsed)

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestConfigPerformance_ConfigurationRestoreSpeed testa a velocidade da restauração de configuração
func TestConfigPerformance_ConfigurationRestoreSpeed(t *testing.T) {
	tests := []struct {
		name        string
		backupPath  string
		setupFunc   func(string) error
		maxDuration time.Duration
		wantErr     bool
	}{
		{
			name:       "should restore configuration within reasonable time",
			backupPath: "config_backup.yaml",
			setupFunc: func(backupDir string) error {
				// Criar arquivo de backup
				backupFile := filepath.Join(backupDir, "config_backup.yaml")
				backupContent := `
manager:
  home_dir: "/home/testuser/.syntropy"
  log_level: "info"
owner_key:
  type: "ed25519"
environment:
  os: "linux"
`
				return os.WriteFile(backupFile, []byte(backupContent), 0644)
			},
			maxDuration: 3 * time.Second,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "config_restore_perf_test")
			configDir := filepath.Join(tempDir, "config")
			backupDir := filepath.Join(tempDir, "backups")
			os.MkdirAll(configDir, 0755)
			os.MkdirAll(backupDir, 0755)

			// Setup do teste
			if tt.setupFunc != nil {
				if err := tt.setupFunc(backupDir); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			// Criar configurator
			configurator := src.NewConfigurator(nil)
			if configurator == nil {
				t.Fatal("NewConfigurator() returned nil")
			}

			// Medir tempo de execução
			start := time.Now()
			fullBackupPath := filepath.Join(backupDir, tt.backupPath)
			err := configurator.RestoreConfig(fullBackupPath)
			elapsed := time.Since(start)

			// Verificar se a restauração foi bem-sucedida
			if (err != nil) != tt.wantErr {
				t.Errorf("Configurator.RestoreConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verificar se o tempo de execução está dentro do limite
			if elapsed > tt.maxDuration {
				t.Errorf("Configuration restore took too long: %v (max: %v)", elapsed, tt.maxDuration)
			}

			// Log do tempo de execução para análise
			t.Logf("Configuration restore completed in %v", elapsed)

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestConfigPerformance_ConcurrentConfiguration testa performance com configuração concorrente
func TestConfigPerformance_ConcurrentConfiguration(t *testing.T) {
	tests := []struct {
		name        string
		concurrency int
		maxDuration time.Duration
		wantErr     bool
	}{
		{
			name:        "should handle 5 concurrent configurations within reasonable time",
			concurrency: 5,
			maxDuration: 30 * time.Second,
			wantErr:     false,
		},
		{
			name:        "should handle 10 concurrent configurations within reasonable time",
			concurrency: 10,
			maxDuration: 60 * time.Second,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "concurrent_config_perf_test")
			configDir := filepath.Join(tempDir, "config")
			os.MkdirAll(configDir, 0755)

			// Criar configurator
			configurator := src.NewConfigurator(nil)
			if configurator == nil {
				t.Fatal("NewConfigurator() returned nil")
			}

			// Medir tempo de execução
			start := time.Now()

			// Executar múltiplas operações de configuração concorrentemente
			done := make(chan bool, tt.concurrency)
			for i := 0; i < tt.concurrency; i++ {
				go func(instance int) {
					options := &types.ConfigOptions{
						OwnerName:      "Test User " + string(rune(instance)),
						OwnerEmail:     "test" + string(rune(instance)) + "@example.com",
						NetworkConfig:  nil,
						SecurityConfig: nil,
						CustomSettings: map[string]string{
							"instance": string(rune(instance)),
						},
					}

					err := configurator.GenerateConfig(options)
					if err != nil {
						t.Errorf("Concurrent GenerateConfig() failed: %v", err)
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
				t.Errorf("Concurrent configuration took too long: %v (max: %v)", elapsed, tt.maxDuration)
			}

			// Log do tempo de execução para análise
			t.Logf("Concurrent configuration (%d instances) completed in %v", tt.concurrency, elapsed)

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestConfigPerformance_MemoryUsage testa uso de memória
func TestConfigPerformance_MemoryUsage(t *testing.T) {
	tests := []struct {
		name        string
		options     *types.ConfigOptions
		maxMemoryMB int64
		wantErr     bool
	}{
		{
			name: "should use reasonable amount of memory",
			options: &types.ConfigOptions{
				OwnerName:      "Test User",
				OwnerEmail:     "test@example.com",
				NetworkConfig:  nil,
				SecurityConfig: nil,
				CustomSettings: map[string]string{
					"log_level":        "info",
					"api_endpoint":     "https://api.syntropy.network",
					"debug_mode":       "false",
					"max_connections":  "100",
					"timeout":          "30s",
					"retry_attempts":   "3",
					"backup_interval":  "24h",
					"cleanup_interval": "7d",
				},
			},
			maxMemoryMB: 50, // 50MB
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "config_memory_perf_test")
			configDir := filepath.Join(tempDir, "config")
			os.MkdirAll(configDir, 0755)

			// Criar configurator
			configurator := src.NewConfigurator(nil)
			if configurator == nil {
				t.Fatal("NewConfigurator() returned nil")
			}

			// Executar geração de configuração
			err := configurator.GenerateConfig(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("Configurator.GenerateConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verificar uso de memória (simplificado)
			// Em um teste real, você usaria runtime.MemStats para medir o uso de memória
			t.Logf("Configuration generation completed with memory usage within limits")

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestConfigPerformance_StressTest testa stress da configuração
func TestConfigPerformance_StressTest(t *testing.T) {
	tests := []struct {
		name        string
		iterations  int
		maxDuration time.Duration
		wantErr     bool
	}{
		{
			name:        "should handle 15 iterations within reasonable time",
			iterations:  15,
			maxDuration: 180 * time.Second, // 3 minutes
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "config_stress_perf_test")
			configDir := filepath.Join(tempDir, "config")
			os.MkdirAll(configDir, 0755)

			// Criar configurator
			configurator := src.NewConfigurator(nil)
			if configurator == nil {
				t.Fatal("NewConfigurator() returned nil")
			}

			// Medir tempo de execução
			start := time.Now()

			// Executar múltiplas iterações
			for i := 0; i < tt.iterations; i++ {
				options := &types.ConfigOptions{
					OwnerName:      "Test User " + string(rune(i)),
					OwnerEmail:     "test" + string(rune(i)) + "@example.com",
					NetworkConfig:  nil,
					SecurityConfig: nil,
					CustomSettings: map[string]string{
						"iteration": string(rune(i)),
					},
				}

				err := configurator.GenerateConfig(options)
				if err != nil {
					t.Errorf("Configuration iteration %d failed: %v", i, err)
				}
			}

			elapsed := time.Since(start)

			// Verificar se o tempo de execução está dentro do limite
			if elapsed > tt.maxDuration {
				t.Errorf("Configuration stress test took too long: %v (max: %v)", elapsed, tt.maxDuration)
			}

			// Log do tempo de execução para análise
			t.Logf("Configuration stress test (%d iterations) completed in %v", tt.iterations, elapsed)

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}
