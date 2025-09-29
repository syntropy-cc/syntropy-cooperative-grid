package unit

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/tests/helpers"
)

// TestNewSetupLogger testa a criação de um novo SetupLogger
func TestNewSetupLogger(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "should create setup logger successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := src.NewSetupLogger()
			if logger == nil {
				t.Error("NewSetupLogger() returned nil logger")
			}
		})
	}
}

// TestSetupLogger_SetVerbose testa o método SetVerbose
func TestSetupLogger_SetVerbose(t *testing.T) {
	tests := []struct {
		name    string
		verbose bool
		wantErr bool
	}{
		{
			name:    "should set verbose mode successfully",
			verbose: true,
			wantErr: false,
		},
		{
			name:    "should disable verbose mode successfully",
			verbose: false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "logger_verbose_test")
			logDir := filepath.Join(tempDir, "logs")
			os.MkdirAll(logDir, 0755)

			// Criar logger com diretório temporário
			logger := &src.SetupLogger{
				LogDir: logDir,
			}

			err := logger.SetVerbose(tt.verbose)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupLogger.SetVerbose() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestSetupLogger_SetQuiet testa o método SetQuiet
func TestSetupLogger_SetQuiet(t *testing.T) {
	tests := []struct {
		name    string
		quiet   bool
		wantErr bool
	}{
		{
			name:    "should set quiet mode successfully",
			quiet:   true,
			wantErr: false,
		},
		{
			name:    "should disable quiet mode successfully",
			quiet:   false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "logger_quiet_test")
			logDir := filepath.Join(tempDir, "logs")
			os.MkdirAll(logDir, 0755)

			// Criar logger com diretório temporário
			logger := &src.SetupLogger{
				LogDir: logDir,
			}

			err := logger.SetQuiet(tt.quiet)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupLogger.SetQuiet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestSetupLogger_LogStep testa o método LogStep
func TestSetupLogger_LogStep(t *testing.T) {
	tests := []struct {
		name    string
		step    string
		message string
		wantErr bool
	}{
		{
			name:    "should log step successfully",
			step:    "validation",
			message: "Starting validation process",
			wantErr: false,
		},
		{
			name:    "should handle empty step",
			step:    "",
			message: "Empty step message",
			wantErr: false,
		},
		{
			name:    "should handle empty message",
			step:    "validation",
			message: "",
			wantErr: false,
		},
		{
			name:    "should handle both empty step and message",
			step:    "",
			message: "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "logger_step_test")
			logDir := filepath.Join(tempDir, "logs")
			os.MkdirAll(logDir, 0755)

			// Criar logger com diretório temporário
			logger := &src.SetupLogger{
				LogDir: logDir,
			}

			err := logger.LogStep(tt.step, tt.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupLogger.LogStep() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestSetupLogger_LogError testa o método LogError
func TestSetupLogger_LogError(t *testing.T) {
	tests := []struct {
		name    string
		message string
		err     error
		wantErr bool
	}{
		{
			name:    "should log error successfully",
			message: "Validation failed",
			err:     errors.New("validation error"),
			wantErr: false,
		},
		{
			name:    "should handle nil error",
			message: "Error without details",
			err:     nil,
			wantErr: false,
		},
		{
			name:    "should handle empty message",
			message: "",
			err:     errors.New("error without message"),
			wantErr: false,
		},
		{
			name:    "should handle both empty message and nil error",
			message: "",
			err:     nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "logger_error_test")
			logDir := filepath.Join(tempDir, "logs")
			os.MkdirAll(logDir, 0755)

			// Criar logger com diretório temporário
			logger := &src.SetupLogger{
				LogDir: logDir,
			}

			err := logger.LogError(tt.message, tt.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupLogger.LogError() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestSetupLogger_LogWarning testa o método LogWarning
func TestSetupLogger_LogWarning(t *testing.T) {
	tests := []struct {
		name    string
		message string
		wantErr bool
	}{
		{
			name:    "should log warning successfully",
			message: "This is a warning message",
			wantErr: false,
		},
		{
			name:    "should handle empty message",
			message: "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "logger_warning_test")
			logDir := filepath.Join(tempDir, "logs")
			os.MkdirAll(logDir, 0755)

			// Criar logger com diretório temporário
			logger := &src.SetupLogger{
				LogDir: logDir,
			}

			err := logger.LogWarning(tt.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupLogger.LogWarning() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestSetupLogger_LogInfo testa o método LogInfo
func TestSetupLogger_LogInfo(t *testing.T) {
	tests := []struct {
		name    string
		message string
		wantErr bool
	}{
		{
			name:    "should log info successfully",
			message: "This is an info message",
			wantErr: false,
		},
		{
			name:    "should handle empty message",
			message: "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "logger_info_test")
			logDir := filepath.Join(tempDir, "logs")
			os.MkdirAll(logDir, 0755)

			// Criar logger com diretório temporário
			logger := &src.SetupLogger{
				LogDir: logDir,
			}

			err := logger.LogInfo(tt.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupLogger.LogInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestSetupLogger_LogDebug testa o método LogDebug
func TestSetupLogger_LogDebug(t *testing.T) {
	tests := []struct {
		name    string
		message string
		wantErr bool
	}{
		{
			name:    "should log debug successfully",
			message: "This is a debug message",
			wantErr: false,
		},
		{
			name:    "should handle empty message",
			message: "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "logger_debug_test")
			logDir := filepath.Join(tempDir, "logs")
			os.MkdirAll(logDir, 0755)

			// Criar logger com diretório temporário
			logger := &src.SetupLogger{
				LogDir: logDir,
			}

			err := logger.LogDebug(tt.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupLogger.LogDebug() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestSetupLogger_ExportLogs testa o método ExportLogs
func TestSetupLogger_ExportLogs(t *testing.T) {
	tests := []struct {
		name       string
		exportPath string
		setupFunc  func(string) error
		wantErr    bool
	}{
		{
			name:       "should export logs successfully",
			exportPath: "logs_export.json",
			setupFunc: func(logDir string) error {
				// Criar arquivo de log
				logFile := filepath.Join(logDir, "setup.log")
				logContent := `2023-01-01T00:00:00Z [INFO] Setup started
2023-01-01T00:00:01Z [ERROR] Validation failed
2023-01-01T00:00:02Z [WARN] Warning message
2023-01-01T00:00:03Z [DEBUG] Debug message`
				return os.WriteFile(logFile, []byte(logContent), 0644)
			},
			wantErr: false,
		},
		{
			name:       "should fail when log directory does not exist",
			exportPath: "logs_export.json",
			setupFunc: func(logDir string) error {
				// Não criar diretório de logs
				return nil
			},
			wantErr: true,
		},
		{
			name:       "should handle empty export path",
			exportPath: "",
			setupFunc: func(logDir string) error {
				// Criar arquivo de log
				logFile := filepath.Join(logDir, "setup.log")
				logContent := `2023-01-01T00:00:00Z [INFO] Setup started`
				return os.WriteFile(logFile, []byte(logContent), 0644)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "logger_export_test")
			logDir := filepath.Join(tempDir, "logs")
			exportDir := filepath.Join(tempDir, "exports")
			os.MkdirAll(exportDir, 0755)

			// Setup do teste
			if tt.setupFunc != nil {
				if err := tt.setupFunc(logDir); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			// Criar logger com diretório temporário
			logger := &src.SetupLogger{
				LogDir: logDir,
			}

			fullExportPath := filepath.Join(exportDir, tt.exportPath)
			err := logger.ExportLogs(fullExportPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupLogger.ExportLogs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestSetupLogger_Close testa o método Close
func TestSetupLogger_Close(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should close logger successfully",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "logger_close_test")
			logDir := filepath.Join(tempDir, "logs")
			os.MkdirAll(logDir, 0755)

			// Criar logger com diretório temporário
			logger := &src.SetupLogger{
				LogDir: logDir,
			}

			err := logger.Close()
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupLogger.Close() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestSetupLogger_RotateLogs testa o método RotateLogs
func TestSetupLogger_RotateLogs(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(string) error
		wantErr   bool
	}{
		{
			name: "should rotate logs successfully",
			setupFunc: func(logDir string) error {
				// Criar arquivo de log
				logFile := filepath.Join(logDir, "setup.log")
				logContent := `2023-01-01T00:00:00Z [INFO] Setup started
2023-01-01T00:00:01Z [ERROR] Validation failed
2023-01-01T00:00:02Z [WARN] Warning message
2023-01-01T00:00:03Z [DEBUG] Debug message`
				return os.WriteFile(logFile, []byte(logContent), 0644)
			},
			wantErr: false,
		},
		{
			name: "should handle empty log directory",
			setupFunc: func(logDir string) error {
				// Não criar arquivo de log
				return nil
			},
			wantErr: false,
		},
		{
			name: "should handle non-existent log directory",
			setupFunc: func(logDir string) error {
				// Não criar diretório de logs
				return nil
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório temporário para testes
			tempDir := helpers.CreateTempDir(t, "logger_rotate_test")
			logDir := filepath.Join(tempDir, "logs")

			// Setup do teste
			if tt.setupFunc != nil {
				if err := tt.setupFunc(logDir); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			// Criar logger com diretório temporário
			logger := &src.SetupLogger{
				LogDir: logDir,
			}

			err := logger.RotateLogs()
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupLogger.RotateLogs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Limpar diretório temporário
			os.RemoveAll(tempDir)
		})
	}
}

// TestSetupLogger_EdgeCases testa casos extremos do SetupLogger
func TestSetupLogger_EdgeCases(t *testing.T) {
	t.Run("should handle nil logger", func(t *testing.T) {
		var logger *src.SetupLogger = nil

		// Should not panic
		err := logger.SetVerbose(true)
		if err != nil {
			t.Errorf("SetVerbose() failed with nil logger: %v", err)
		}
	})

	t.Run("should handle empty log directory", func(t *testing.T) {
		tempDir := helpers.CreateTempDir(t, "logger_edge_test")

		logger := &src.SetupLogger{
			LogDir: "",
		}

		// Should not panic
		err := logger.LogInfo("Test message")
		if err != nil {
			t.Errorf("LogInfo() failed with empty log directory: %v", err)
		}

		os.RemoveAll(tempDir)
	})
}

// TestSetupLogger_Concurrency testa concorrência do SetupLogger
func TestSetupLogger_Concurrency(t *testing.T) {
	t.Run("should handle concurrent logging", func(t *testing.T) {
		tempDir := helpers.CreateTempDir(t, "logger_concurrent_test")
		logDir := filepath.Join(tempDir, "logs")
		os.MkdirAll(logDir, 0755)

		logger := &src.SetupLogger{
			LogDir: logDir,
		}

		// Executar múltiplas chamadas de logging concorrentemente
		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func(step int) {
				err := logger.LogInfo("Concurrent log message " + string(rune(step)))
				if err != nil {
					t.Errorf("Concurrent LogInfo() failed: %v", err)
				}
				done <- true
			}(i)
		}

		// Aguardar todas as goroutines terminarem
		for i := 0; i < 10; i++ {
			<-done
		}

		os.RemoveAll(tempDir)
	})
}

// TestSetupLogger_Performance testa performance do SetupLogger
func TestSetupLogger_Performance(t *testing.T) {
	t.Run("should complete operations within reasonable time", func(t *testing.T) {
		tempDir := helpers.CreateTempDir(t, "logger_perf_test")
		logDir := filepath.Join(tempDir, "logs")
		os.MkdirAll(logDir, 0755)

		logger := &src.SetupLogger{
			LogDir: logDir,
		}

		start := time.Now()
		err := logger.LogInfo("Performance test message")
		elapsed := time.Since(start)

		if err != nil {
			t.Errorf("LogInfo() failed: %v", err)
		}

		if elapsed > 1*time.Second {
			t.Errorf("LogInfo() took too long: %v", elapsed)
		}

		os.RemoveAll(tempDir)
	})
}

// TestSetupLogger_LogLevels testa diferentes níveis de log
func TestSetupLogger_LogLevels(t *testing.T) {
	t.Run("should handle all log levels correctly", func(t *testing.T) {
		tempDir := helpers.CreateTempDir(t, "logger_levels_test")
		logDir := filepath.Join(tempDir, "logs")
		os.MkdirAll(logDir, 0755)

		logger := &src.SetupLogger{
			LogDir: logDir,
		}

		// Testar todos os níveis de log
		testCases := []struct {
			name    string
			logFunc func(string) error
		}{
			{"INFO", logger.LogInfo},
			{"WARNING", logger.LogWarning},
			{"DEBUG", logger.LogDebug},
		}

		for _, tc := range testCases {
			err := tc.logFunc("Test " + tc.name + " message")
			if err != nil {
				t.Errorf("Log%s() failed: %v", tc.name, err)
			}
		}

		// Testar LogStep
		err := logger.LogStep("test_step", "Test step message")
		if err != nil {
			t.Errorf("LogStep() failed: %v", err)
		}

		// Testar LogError
		err = logger.LogError("Test error message", errors.New("test error"))
		if err != nil {
			t.Errorf("LogError() failed: %v", err)
		}

		os.RemoveAll(tempDir)
	})
}

// TestSetupLogger_LogRotation testa rotação de logs
func TestSetupLogger_LogRotation(t *testing.T) {
	t.Run("should rotate logs correctly", func(t *testing.T) {
		tempDir := helpers.CreateTempDir(t, "logger_rotation_test")
		logDir := filepath.Join(tempDir, "logs")
		os.MkdirAll(logDir, 0755)

		logger := &src.SetupLogger{
			LogDir: logDir,
		}

		// Criar arquivo de log inicial
		logFile := filepath.Join(logDir, "setup.log")
		logContent := `2023-01-01T00:00:00Z [INFO] Setup started
2023-01-01T00:00:01Z [ERROR] Validation failed
2023-01-01T00:00:02Z [WARN] Warning message
2023-01-01T00:00:03Z [DEBUG] Debug message`
		os.WriteFile(logFile, []byte(logContent), 0644)

		// Rotacionar logs
		err := logger.RotateLogs()
		if err != nil {
			t.Errorf("RotateLogs() failed: %v", err)
		}

		// Verificar que o arquivo de log foi rotacionado
		rotatedFile := filepath.Join(logDir, "setup.log.1")
		if _, err := os.Stat(rotatedFile); os.IsNotExist(err) {
			t.Error("Log rotation did not create rotated file")
		}

		os.RemoveAll(tempDir)
	})
}
