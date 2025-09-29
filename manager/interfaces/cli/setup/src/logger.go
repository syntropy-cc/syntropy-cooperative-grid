package setup

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src/internal/types"
)

// SetupLogger implementa a interface SetupLogger
type SetupLogger struct {
	logDir        string
	logFile       *os.File
	verbose       bool
	quiet         bool
	correlationID string
}

// NewSetupLogger cria um novo logger estruturado
func NewSetupLogger() *SetupLogger {
	homeDir, _ := os.UserHomeDir()
	logDir := filepath.Join(homeDir, ".syntropy", "logs")

	// Criar diretório de logs se não existir
	os.MkdirAll(logDir, 0755)

	// Criar arquivo de log com timestamp
	timestamp := time.Now().Format("20060102_150405")
	logPath := filepath.Join(logDir, fmt.Sprintf("setup_%s.log", timestamp))

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		// Fallback para stdout se não conseguir criar arquivo
		logFile = os.Stdout
	}

	return &SetupLogger{
		logDir:        logDir,
		logFile:       logFile,
		verbose:       false,
		quiet:         false,
		correlationID: generateCorrelationID(),
	}
}

// SetVerbose define se o logger deve ser verboso
func (sl *SetupLogger) SetVerbose(verbose bool) {
	sl.verbose = verbose
}

// SetQuiet define se o logger deve ser silencioso
func (sl *SetupLogger) SetQuiet(quiet bool) {
	sl.quiet = quiet
}

// LogStep registra uma etapa do setup
func (sl *SetupLogger) LogStep(step string, data map[string]interface{}) {
	if sl.quiet {
		return
	}

	entry := &types.LogEntry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   fmt.Sprintf("Setup step: %s", step),
		Step:      step,
		Data:      data,
	}

	sl.writeLogEntry(entry)

	if sl.verbose {
		fmt.Printf("[STEP] %s\n", step)
		if data != nil {
			for key, value := range data {
				fmt.Printf("  %s: %v\n", key, value)
			}
		}
	}
}

// LogError registra um erro
func (sl *SetupLogger) LogError(err error, context map[string]interface{}) {
	entry := &types.LogEntry{
		Timestamp: time.Now(),
		Level:     "ERROR",
		Message:   err.Error(),
		Error:     err.Error(),
		Data:      context,
	}

	sl.writeLogEntry(entry)

	if !sl.quiet {
		fmt.Printf("[ERROR] %s\n", err.Error())
		if context != nil && sl.verbose {
			for key, value := range context {
				fmt.Printf("  %s: %v\n", key, value)
			}
		}
	}
}

// LogWarning registra um warning
func (sl *SetupLogger) LogWarning(message string, data map[string]interface{}) {
	if sl.quiet {
		return
	}

	entry := &types.LogEntry{
		Timestamp: time.Now(),
		Level:     "WARN",
		Message:   message,
		Data:      data,
	}

	sl.writeLogEntry(entry)

	fmt.Printf("[WARNING] %s\n", message)
	if data != nil && sl.verbose {
		for key, value := range data {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}
}

// LogInfo registra uma informação
func (sl *SetupLogger) LogInfo(message string, data map[string]interface{}) {
	if sl.quiet {
		return
	}

	entry := &types.LogEntry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   message,
		Data:      data,
	}

	sl.writeLogEntry(entry)

	if sl.verbose {
		fmt.Printf("[INFO] %s\n", message)
		if data != nil {
			for key, value := range data {
				fmt.Printf("  %s: %v\n", key, value)
			}
		}
	}
}

// LogDebug registra uma informação de debug
func (sl *SetupLogger) LogDebug(message string, data map[string]interface{}) {
	if sl.quiet || !sl.verbose {
		return
	}

	entry := &types.LogEntry{
		Timestamp: time.Now(),
		Level:     "DEBUG",
		Message:   message,
		Data:      data,
	}

	sl.writeLogEntry(entry)

	fmt.Printf("[DEBUG] %s\n", message)
	if data != nil {
		for key, value := range data {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}
}

// ExportLogs exporta logs em formato específico
func (sl *SetupLogger) ExportLogs(format string, outputPath string) error {
	switch format {
	case "json":
		return sl.exportJSONLogs(outputPath)
	case "csv":
		return sl.exportCSVLogs(outputPath)
	case "txt":
		return sl.exportTextLogs(outputPath)
	default:
		return fmt.Errorf("formato de exportação não suportado: %s", format)
	}
}

// writeLogEntry escreve uma entrada de log
func (sl *SetupLogger) writeLogEntry(entry *types.LogEntry) {
	// Adicionar metadados
	if entry.Data == nil {
		entry.Data = make(map[string]interface{})
	}
	entry.Data["correlation_id"] = sl.correlationID
	entry.Data["os"] = runtime.GOOS
	entry.Data["arch"] = runtime.GOARCH

	// Serializar para JSON
	jsonData, err := json.Marshal(entry)
	if err != nil {
		// Fallback para formato simples
		sl.logFile.WriteString(fmt.Sprintf("%s [%s] %s\n",
			entry.Timestamp.Format(time.RFC3339),
			entry.Level,
			entry.Message))
		return
	}

	// Escrever no arquivo
	sl.logFile.WriteString(string(jsonData) + "\n")
	sl.logFile.Sync()
}

// exportJSONLogs exporta logs em formato JSON
func (sl *SetupLogger) exportJSONLogs(outputPath string) error {
	// Ler arquivo de log atual
	logData, err := os.ReadFile(sl.logFile.Name())
	if err != nil {
		return fmt.Errorf("falha ao ler arquivo de log: %w", err)
	}

	// Escrever para arquivo de saída
	return os.WriteFile(outputPath, logData, 0644)
}

// exportCSVLogs exporta logs em formato CSV
func (sl *SetupLogger) exportCSVLogs(outputPath string) error {
	// Implementação simplificada - em produção seria mais robusta
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("falha ao criar arquivo CSV: %w", err)
	}
	defer file.Close()

	// Escrever cabeçalho CSV
	file.WriteString("timestamp,level,message,step,error\n")

	// Ler e processar logs
	logData, err := os.ReadFile(sl.logFile.Name())
	if err != nil {
		return fmt.Errorf("falha ao ler arquivo de log: %w", err)
	}

	lines := strings.Split(string(logData), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		var entry types.LogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue // Pular linhas inválidas
		}

		// Escrever linha CSV
		file.WriteString(fmt.Sprintf("%s,%s,%s,%s,%s\n",
			entry.Timestamp.Format(time.RFC3339),
			entry.Level,
			escapeCSV(entry.Message),
			entry.Step,
			escapeCSV(entry.Error)))
	}

	return nil
}

// exportTextLogs exporta logs em formato texto
func (sl *SetupLogger) exportTextLogs(outputPath string) error {
	// Ler arquivo de log atual
	logData, err := os.ReadFile(sl.logFile.Name())
	if err != nil {
		return fmt.Errorf("falha ao ler arquivo de log: %w", err)
	}

	// Converter JSON para texto legível
	var textOutput string
	lines := strings.Split(string(logData), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		var entry types.LogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue // Pular linhas inválidas
		}

		textOutput += fmt.Sprintf("[%s] [%s] %s\n",
			entry.Timestamp.Format("2006-01-02 15:04:05"),
			entry.Level,
			entry.Message)

		if entry.Step != "" {
			textOutput += fmt.Sprintf("  Step: %s\n", entry.Step)
		}

		if entry.Error != "" {
			textOutput += fmt.Sprintf("  Error: %s\n", entry.Error)
		}

		if entry.Data != nil && len(entry.Data) > 0 {
			textOutput += "  Data:\n"
			for key, value := range entry.Data {
				textOutput += fmt.Sprintf("    %s: %v\n", key, value)
			}
		}

		textOutput += "\n"
	}

	return os.WriteFile(outputPath, []byte(textOutput), 0644)
}

// Close fecha o logger
func (sl *SetupLogger) Close() error {
	if sl.logFile != nil && sl.logFile != os.Stdout {
		return sl.logFile.Close()
	}
	return nil
}

// RotateLogs rotaciona os logs
func (sl *SetupLogger) RotateLogs() error {
	// Fechar arquivo atual
	if err := sl.Close(); err != nil {
		return err
	}

	// Criar novo arquivo de log
	timestamp := time.Now().Format("20060102_150405")
	logPath := filepath.Join(sl.logDir, fmt.Sprintf("setup_%s.log", timestamp))

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("falha ao criar novo arquivo de log: %w", err)
	}

	sl.logFile = logFile
	sl.correlationID = generateCorrelationID()

	return nil
}

// generateCorrelationID gera um ID de correlação único
func generateCorrelationID() string {
	return fmt.Sprintf("setup_%d_%d", time.Now().Unix(), os.Getpid())
}

// escapeCSV escapa strings para CSV
func escapeCSV(s string) string {
	if s == "" {
		return ""
	}

	// Escapar aspas duplas
	s = fmt.Sprintf("\"%s\"", s)
	return s
}
