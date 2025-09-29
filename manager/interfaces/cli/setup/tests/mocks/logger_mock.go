package mocks

import (
	"errors"
)

// MockSetupLogger implementa a interface SetupLogger para testes
type MockSetupLogger struct {
	LogStepFunc    func(step string, data map[string]interface{})
	LogErrorFunc   func(err error, context map[string]interface{})
	LogWarningFunc func(message string, data map[string]interface{})
	LogInfoFunc    func(message string, data map[string]interface{})
	LogDebugFunc   func(message string, data map[string]interface{})
	ExportLogsFunc func(format string, outputPath string) error
	CloseFunc      func() error
}

// LogStep chama a função mock
func (m *MockSetupLogger) LogStep(step string, data map[string]interface{}) {
	if m.LogStepFunc != nil {
		m.LogStepFunc(step, data)
	}
}

// LogError chama a função mock
func (m *MockSetupLogger) LogError(err error, context map[string]interface{}) {
	if m.LogErrorFunc != nil {
		m.LogErrorFunc(err, context)
	}
}

// LogWarning chama a função mock
func (m *MockSetupLogger) LogWarning(message string, data map[string]interface{}) {
	if m.LogWarningFunc != nil {
		m.LogWarningFunc(message, data)
	}
}

// LogInfo chama a função mock
func (m *MockSetupLogger) LogInfo(message string, data map[string]interface{}) {
	if m.LogInfoFunc != nil {
		m.LogInfoFunc(message, data)
	}
}

// LogDebug chama a função mock
func (m *MockSetupLogger) LogDebug(message string, data map[string]interface{}) {
	if m.LogDebugFunc != nil {
		m.LogDebugFunc(message, data)
	}
}

// ExportLogs chama a função mock
func (m *MockSetupLogger) ExportLogs(format string, outputPath string) error {
	if m.ExportLogsFunc != nil {
		return m.ExportLogsFunc(format, outputPath)
	}
	return nil
}

// Close chama a função mock
func (m *MockSetupLogger) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

// MockSetupLoggerWithErrors implementa SetupLogger que retorna erros para testes de falha
type MockSetupLoggerWithErrors struct {
	*MockSetupLogger
}

// NewMockSetupLoggerWithErrors cria um logger mock que retorna erros
func NewMockSetupLoggerWithErrors() *MockSetupLoggerWithErrors {
	return &MockSetupLoggerWithErrors{
		MockSetupLogger: &MockSetupLogger{
			ExportLogsFunc: func(format string, outputPath string) error {
				return errors.New("mock export error")
			},
			CloseFunc: func() error {
				return errors.New("mock close error")
			},
		},
	}
}
