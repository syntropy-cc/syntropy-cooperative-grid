package types

import (
	"time"
)

// TestConfig representa uma configuração de teste
type TestConfig struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Options     map[string]interface{} `json:"options"`
	Expected    ExpectedResult         `json:"expected"`
}

// ExpectedResult representa o resultado esperado de um teste
type ExpectedResult struct {
	Success bool                   `json:"success"`
	Error   string                 `json:"error,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// TestCase representa um caso de teste
type TestCase struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Setup       func() error           `json:"-"`
	Test        func() error           `json:"-"`
	Teardown    func() error           `json:"-"`
	Expected    ExpectedResult         `json:"expected"`
	Timeout     time.Duration          `json:"timeout"`
	Skip        bool                   `json:"skip"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

// TestSuite representa uma suíte de testes
type TestSuite struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Cases       []TestCase    `json:"cases"`
	Setup       func() error  `json:"-"`
	Teardown    func() error  `json:"-"`
	Timeout     time.Duration `json:"timeout"`
}

// TestResult representa o resultado de um teste
type TestResult struct {
	Name      string                 `json:"name"`
	Success   bool                   `json:"success"`
	Error     string                 `json:"error,omitempty"`
	Duration  time.Duration          `json:"duration"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// PerformanceMetrics representa métricas de performance
type PerformanceMetrics struct {
	Duration    time.Duration `json:"duration"`
	MemoryUsage uint64        `json:"memory_usage"`
	CPUTime     time.Duration `json:"cpu_time"`
	Iterations  int           `json:"iterations"`
}

// SecurityTest representa um teste de segurança
type SecurityTest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Type        SecurityTestType       `json:"type"`
	Payload     string                 `json:"payload"`
	Expected    SecurityExpectedResult `json:"expected"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

// SecurityTestType representa o tipo de teste de segurança
type SecurityTestType string

const (
	SecurityTestTypeInjection       SecurityTestType = "injection"
	SecurityTestTypePathTraversal   SecurityTestType = "path_traversal"
	SecurityTestTypePermission      SecurityTestType = "permission"
	SecurityTestTypeAuthentication  SecurityTestType = "authentication"
	SecurityTestTypeAuthorization   SecurityTestType = "authorization"
	SecurityTestTypeDataValidation  SecurityTestType = "data_validation"
	SecurityTestTypeFileSecurity    SecurityTestType = "file_security"
	SecurityTestTypeNetworkSecurity SecurityTestType = "network_security"
)

// SecurityExpectedResult representa o resultado esperado de um teste de segurança
type SecurityExpectedResult struct {
	ShouldFail    bool   `json:"should_fail"`
	ExpectedError string `json:"expected_error,omitempty"`
	Blocked       bool   `json:"blocked"`
}

// LoadTest representa um teste de carga
type LoadTest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Concurrency int                    `json:"concurrency"`
	Duration    time.Duration          `json:"duration"`
	Iterations  int                    `json:"iterations"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

// IntegrationTest representa um teste de integração
type IntegrationTest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Components  []string               `json:"components"`
	Setup       func() error           `json:"-"`
	Test        func() error           `json:"-"`
	Teardown    func() error           `json:"-"`
	Expected    ExpectedResult         `json:"expected"`
	Timeout     time.Duration          `json:"timeout"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

// E2ETest representa um teste end-to-end
type E2ETest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Workflow    []WorkflowStep         `json:"workflow"`
	Setup       func() error           `json:"-"`
	Teardown    func() error           `json:"-"`
	Expected    ExpectedResult         `json:"expected"`
	Timeout     time.Duration          `json:"timeout"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

// WorkflowStep representa um passo em um workflow de teste
type WorkflowStep struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Action      func() error           `json:"-"`
	Expected    ExpectedResult         `json:"expected"`
	Timeout     time.Duration          `json:"timeout"`
	Data        map[string]interface{} `json:"data,omitempty"`
}
