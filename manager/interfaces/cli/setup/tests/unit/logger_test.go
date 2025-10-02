//go:build !integration && !e2e && !performance && !security
// +build !integration,!e2e,!performance,!security

package unit

import (
	"os"
	"path/filepath"
	"testing"

	setup "setup-component/src"
)

// TestNewSetupLogger testa a criação do logger
func TestNewSetupLogger(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should create logger successfully",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := setup.NewSetupLogger()
			if logger == nil {
				t.Error("NewSetupLogger() returned nil logger")
			}
		})
	}
}

// TestSetupLogger_Info testa o log de nível INFO
func TestSetupLogger_Info(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	tests := []struct {
		name    string
		message string
		fields  map[string]interface{}
		wantErr bool
	}{
		{
			name:    "should log info message successfully",
			message: "Test info message",
			fields:  map[string]interface{}{"test": "value"},
			wantErr: false,
		},
		{
			name:    "should log info message with empty fields",
			message: "Test info message",
			fields:  map[string]interface{}{},
			wantErr: false,
		},
		{
			name:    "should log info message with nil fields",
			message: "Test info message",
			fields:  nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.LogInfo(tt.message, tt.fields)
			// Não há retorno de erro, então apenas verificamos se não houve panic
		})
	}
}

// TestSetupLogger_Error testa o log de nível ERROR
func TestSetupLogger_Error(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	tests := []struct {
		name    string
		message string
		fields  map[string]interface{}
		wantErr bool
	}{
		{
			name:    "should log error message successfully",
			message: "Test error message",
			fields:  map[string]interface{}{"error": "test error"},
			wantErr: false,
		},
		{
			name:    "should log error message with empty fields",
			message: "Test error message",
			fields:  map[string]interface{}{},
			wantErr: false,
		},
		{
			name:    "should log error message with nil fields",
			message: "Test error message",
			fields:  nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.LogError(tt.message, tt.fields)
			// Não há retorno de erro, então apenas verificamos se não houve panic
		})
	}
}

// TestSetupLogger_Warn testa o log de nível WARN
func TestSetupLogger_Warn(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	tests := []struct {
		name    string
		message string
		fields  map[string]interface{}
		wantErr bool
	}{
		{
			name:    "should log warn message successfully",
			message: "Test warn message",
			fields:  map[string]interface{}{"warning": "test warning"},
			wantErr: false,
		},
		{
			name:    "should log warn message with empty fields",
			message: "Test warn message",
			fields:  map[string]interface{}{},
			wantErr: false,
		},
		{
			name:    "should log warn message with nil fields",
			message: "Test warn message",
			fields:  nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.LogWarning(tt.message, tt.fields)
			// Não há retorno de erro, então apenas verificamos se não houve panic
		})
	}
}

// TestSetupLogger_Debug testa o log de nível DEBUG
func TestSetupLogger_Debug(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	tests := []struct {
		name    string
		message string
		fields  map[string]interface{}
		wantErr bool
	}{
		{
			name:    "should log debug message successfully",
			message: "Test debug message",
			fields:  map[string]interface{}{"debug": "test debug"},
			wantErr: false,
		},
		{
			name:    "should log debug message with empty fields",
			message: "Test debug message",
			fields:  map[string]interface{}{},
			wantErr: false,
		},
		{
			name:    "should log debug message with nil fields",
			message: "Test debug message",
			fields:  nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.LogDebug(tt.message, tt.fields)
			// Não há retorno de erro, então apenas verificamos se não houve panic
		})
	}
}

// TestSetupLogger_SetLevel testa a definição do nível de log
func TestSetupLogger_SetLevel(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	tests := []struct {
		name  string
		level types.LogLevel
	}{
		{
			name:  "should set level to INFO",
			level: types.LogLevelInfo,
		},
		{
			name:  "should set level to ERROR",
			level: types.LogLevelError,
		},
		{
			name:  "should set level to WARN",
			level: types.LogLevelWarn,
		},
		{
			name:  "should set level to DEBUG",
			level: types.LogLevelDebug,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.SetVerbose(tt.level)
			// Não há retorno de erro, então apenas verificamos se não houve panic
		})
	}
}

// TestSetupLogger_GetLevel testa a obtenção do nível de log
func TestSetupLogger_GetLevel(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	tests := []struct {
		name  string
		level types.LogLevel
	}{
		{
			name:  "should get level INFO",
			level: types.LogLevelInfo,
		},
		{
			name:  "should get level ERROR",
			level: types.LogLevelError,
		},
		{
			name:  "should get level WARN",
			level: types.LogLevelWarn,
		},
		{
			name:  "should get level DEBUG",
			level: types.LogLevelDebug,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.SetVerbose(tt.level)
			level := logger.IsVerbose()
			if level != tt.level {
				t.Errorf("SetupLogger.GetLevel() = %v, want %v", level, tt.level)
			}
		})
	}
}

// TestSetupLogger_SetCorrelationID testa a definição do ID de correlação
func TestSetupLogger_SetCorrelationID(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	tests := []struct {
		name          string
		correlationID string
	}{
		{
			name:          "should set correlation ID successfully",
			correlationID: "test-correlation-id",
		},
		{
			name:          "should set empty correlation ID",
			correlationID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.SetVerbose(tt.correlationID)
			// Não há retorno de erro, então apenas verificamos se não houve panic
		})
	}
}

// TestSetupLogger_GetCorrelationID testa a obtenção do ID de correlação
func TestSetupLogger_GetCorrelationID(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	tests := []struct {
		name          string
		correlationID string
	}{
		{
			name:          "should get correlation ID successfully",
			correlationID: "test-correlation-id",
		},
		{
			name:          "should get empty correlation ID",
			correlationID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.SetVerbose(tt.correlationID)
			correlationID := logger.IsVerbose()
			if correlationID != tt.correlationID {
				t.Errorf("SetupLogger.GetCorrelationID() = %v, want %v", correlationID, tt.correlationID)
			}
		})
	}
}

// TestSetupLogger_ExportLogs testa a exportação de logs
func TestSetupLogger_ExportLogs(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	tests := []struct {
		name    string
		format  types.LogFormat
		path    string
		wantErr bool
	}{
		{
			name:    "should export logs in JSON format",
			format:  types.LogFormatJSON,
			path:    filepath.Join(tempDir, "logs.json"),
			wantErr: false,
		},
		{
			name:    "should export logs in CSV format",
			format:  types.LogFormatCSV,
			path:    filepath.Join(tempDir, "logs.csv"),
			wantErr: false,
		},
		{
			name:    "should export logs in TXT format",
			format:  types.LogFormatTXT,
			path:    filepath.Join(tempDir, "logs.txt"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Adicionar alguns logs antes da exportação
			logger.LogInfo("Test info message", map[string]interface{}{"test": "value"})
			logger.LogError("Test error message", map[string]interface{}{"error": "test error"})

			err := logger.ExportLogs(tt.format, tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupLogger.ExportLogs() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verificar se o arquivo foi criado
			if !tt.wantErr {
				if _, err := os.Stat(tt.path); os.IsNotExist(err) {
					t.Errorf("Export file not created: %s", tt.path)
				}
			}
		})
	}
}

// TestSetupLogger_RotateLogs testa a rotação de logs
func TestSetupLogger_RotateLogs(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should rotate logs successfully",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Adicionar alguns logs antes da rotação
			logger.LogInfo("Test info message", map[string]interface{}{"test": "value"})
			logger.LogError("Test error message", map[string]interface{}{"error": "test error"})

			err := logger.RotateLogs()
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupLogger.RotateLogs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestSetupLogger_Close testa o fechamento do logger
func TestSetupLogger_Close(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()

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
			err := logger.Close()
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupLogger.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestSetupLogger_GetLogsDir testa a obtenção do diretório de logs
func TestSetupLogger_GetLogsDir(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	tests := []struct {
		name string
		want string
	}{
		{
			name: "should return logs directory path",
			want: filepath.Join(tempDir, ".syntropy", "logs"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := logger.GetLogPath()
			if result != tt.want {
				t.Errorf("SetupLogger.GetLogsDir() = %v, want %v", result, tt.want)
			}
		})
	}
}

// TestSetupLogger_GetLogPath testa a obtenção do caminho de log
func TestSetupLogger_GetLogPath(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	tests := []struct {
		name string
		want string
	}{
		{
			name: "should return log file path",
			want: filepath.Join(tempDir, ".syntropy", "logs", "setup.log"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := logger.GetLogPath()
			if result != tt.want {
				t.Errorf("SetupLogger.GetLogPath() = %v, want %v", result, tt.want)
			}
		})
	}
}
