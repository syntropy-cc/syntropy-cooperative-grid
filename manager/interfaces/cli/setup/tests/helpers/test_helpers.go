package helpers

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src/internal/types"
)

// LoadFixture carrega um fixture JSON do diretório de fixtures
func LoadFixture(t *testing.T, category, filename string, target interface{}) {
	t.Helper()

	fixturePath := filepath.Join("fixtures", category, filename)
	data, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("Failed to read fixture %s: %v", fixturePath, err)
	}

	if err := json.Unmarshal(data, target); err != nil {
		t.Fatalf("Failed to unmarshal fixture %s: %v", fixturePath, err)
	}
}

// CreateTempDir cria um diretório temporário para testes
func CreateTempDir(t *testing.T, prefix string) string {
	t.Helper()

	tempDir, err := os.MkdirTemp("", prefix)
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})

	return tempDir
}

// CreateTempFile cria um arquivo temporário para testes
func CreateTempFile(t *testing.T, dir, filename, content string) string {
	t.Helper()

	filePath := filepath.Join(dir, filename)
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create temp file %s: %v", filePath, err)
	}

	return filePath
}

// AssertFileExists verifica se um arquivo existe
func AssertFileExists(t *testing.T, filePath string) {
	t.Helper()

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("Expected file to exist: %s", filePath)
	}
}

// AssertFileNotExists verifica se um arquivo não existe
func AssertFileNotExists(t *testing.T, filePath string) {
	t.Helper()

	if _, err := os.Stat(filePath); err == nil {
		t.Errorf("Expected file to not exist: %s", filePath)
	}
}

// AssertDirExists verifica se um diretório existe
func AssertDirExists(t *testing.T, dirPath string) {
	t.Helper()

	info, err := os.Stat(dirPath)
	if err != nil {
		t.Errorf("Expected directory to exist: %s", dirPath)
		return
	}

	if !info.IsDir() {
		t.Errorf("Expected %s to be a directory", dirPath)
	}
}

// AssertErrorContains verifica se um erro contém uma string específica
func AssertErrorContains(t *testing.T, err error, expected string) {
	t.Helper()

	if err == nil {
		t.Errorf("Expected error to contain '%s', but got nil", expected)
		return
	}

	if !contains(err.Error(), expected) {
		t.Errorf("Expected error to contain '%s', but got: %s", expected, err.Error())
	}
}

// AssertSetupError verifica se um erro é do tipo SetupError com código específico
func AssertSetupError(t *testing.T, err error, expectedCode string) {
	t.Helper()

	if err == nil {
		t.Errorf("Expected SetupError with code '%s', but got nil", expectedCode)
		return
	}

	setupErr, ok := err.(*types.SetupError)
	if !ok {
		t.Errorf("Expected SetupError, but got: %T", err)
		return
	}

	if setupErr.Code != expectedCode {
		t.Errorf("Expected error code '%s', but got: %s", expectedCode, setupErr.Code)
	}
}

// CreateValidSetupOptions cria opções de setup válidas para testes
func CreateValidSetupOptions() *types.SetupOptions {
	return &types.SetupOptions{
		Force:        false,
		ValidateOnly: false,
		Verbose:      true,
		Quiet:        false,
		ConfigPath:   "",
		CustomSettings: map[string]string{
			"owner_name":  "Test User",
			"owner_email": "test@example.com",
		},
	}
}

// CreateValidConfigOptions cria opções de configuração válidas para testes
func CreateValidConfigOptions() *types.ConfigOptions {
	return &types.ConfigOptions{
		OwnerName:      "Test User",
		OwnerEmail:     "test@example.com",
		NetworkConfig:  nil,
		SecurityConfig: nil,
		CustomSettings: map[string]string{},
	}
}

// CreateValidKeyPair cria um par de chaves válido para testes
func CreateValidKeyPair() *types.KeyPair {
	return &types.KeyPair{
		ID:          "test_key_12345",
		Algorithm:   "ed25519",
		PrivateKey:  []byte("test_private_key"),
		PublicKey:   []byte("test_public_key"),
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().AddDate(1, 0, 0),
		Fingerprint: "test_fingerprint",
		Metadata: map[string]string{
			"generated_by": "test",
			"version":      "1.0.0",
		},
	}
}

// CreateValidSetupState cria um estado de setup válido para testes
func CreateValidSetupState() *types.SetupState {
	return &types.SetupState{
		Version:   "1.0.0",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    types.SetupStatusCompleted,
		Environment: &types.EnvironmentInfo{
			OS:              "linux",
			OSVersion:       "20.04",
			Architecture:    "amd64",
			HasAdminRights:  false,
			AvailableDiskGB: 100.0,
			HasInternet:     true,
			HomeDir:         "/home/testuser",
			CanProceed:      true,
			Issues:          []string{},
		},
		Configuration: &types.ConfigInfo{
			Path:      "/home/testuser/.syntropy/config/manager.yaml",
			Valid:     true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Keys: &types.KeyInfo{
			OwnerKeyID: "test_key_12345",
			Algorithm:  "ed25519",
			CreatedAt:  time.Now(),
			ExpiresAt:  time.Now().AddDate(1, 0, 0),
		},
		LastBackup: nil,
		Metadata:   make(map[string]string),
	}
}

// CreateValidEnvironmentInfo cria informações de ambiente válidas para testes
func CreateValidEnvironmentInfo() *types.EnvironmentInfo {
	return &types.EnvironmentInfo{
		OS:              "linux",
		OSVersion:       "20.04",
		Architecture:    "amd64",
		HasAdminRights:  false,
		AvailableDiskGB: 100.0,
		HasInternet:     true,
		HomeDir:         "/home/testuser",
		CanProceed:      true,
		Issues:          []string{},
	}
}

// CreateValidValidationResult cria um resultado de validação válido para testes
func CreateValidValidationResult() *types.ValidationResult {
	return &types.ValidationResult{
		Environment: CreateValidEnvironmentInfo(),
		Dependencies: &types.DependencyStatus{
			Required:  []types.Dependency{},
			Installed: []types.Dependency{},
			Missing:   []types.Dependency{},
			Outdated:  []types.Dependency{},
		},
		Network: &types.NetworkInfo{
			HasInternet:     true,
			Connectivity:    true,
			ProxyConfigured: false,
			FirewallActive:  false,
			PortsOpen:       []int{8080, 9090},
		},
		Permissions: &types.PermissionStatus{
			FileSystem: true,
			Network:    true,
			Service:    false,
			Admin:      false,
			Issues:     []string{},
		},
		CanProceed: true,
		Issues:     []types.ValidationIssue{},
		Warnings:   []string{},
	}
}

// CreateInvalidSetupOptions cria opções de setup inválidas para testes
func CreateInvalidSetupOptions() *types.SetupOptions {
	return &types.SetupOptions{
		Force:        true,
		ValidateOnly: true, // Conflito: não pode ser force e validate_only
		Verbose:      true,
		Quiet:        true, // Conflito: não pode ser verbose e quiet
		ConfigPath:   "/nonexistent/path",
		CustomSettings: map[string]string{
			"invalid_key": "invalid_value",
		},
	}
}

// CreateInvalidConfigOptions cria opções de configuração inválidas para testes
func CreateInvalidConfigOptions() *types.ConfigOptions {
	return &types.ConfigOptions{
		OwnerName:      "",              // Nome vazio
		OwnerEmail:     "invalid-email", // Email inválido
		NetworkConfig:  nil,
		SecurityConfig: nil,
		CustomSettings: map[string]string{
			"invalid_setting": "invalid_value",
		},
	}
}

// CreateInvalidKeyPair cria um par de chaves inválido para testes
func CreateInvalidKeyPair() *types.KeyPair {
	return &types.KeyPair{
		ID:          "",                  // ID vazio
		Algorithm:   "invalid_algorithm", // Algoritmo inválido
		PrivateKey:  []byte{},            // Chave privada vazia
		PublicKey:   []byte{},            // Chave pública vazia
		CreatedAt:   time.Time{},         // Timestamp zero
		ExpiresAt:   time.Time{},         // Timestamp zero
		Fingerprint: "",                  // Fingerprint vazio
		Metadata:    nil,                 // Metadata nil
	}
}

// CreateInvalidSetupState cria um estado de setup inválido para testes
func CreateInvalidSetupState() *types.SetupState {
	return &types.SetupState{
		Version:       "",          // Versão vazia
		CreatedAt:     time.Time{}, // Timestamp zero
		UpdatedAt:     time.Time{}, // Timestamp zero
		Status:        "",          // Status vazio
		Environment:   nil,         // Environment nil
		Configuration: nil,         // Configuration nil
		Keys:          nil,         // Keys nil
		LastBackup:    nil,
		Metadata:      nil, // Metadata nil
	}
}

// CreateInvalidEnvironmentInfo cria informações de ambiente inválidas para testes
func CreateInvalidEnvironmentInfo() *types.EnvironmentInfo {
	return &types.EnvironmentInfo{
		OS:              "", // OS vazio
		OSVersion:       "", // Versão vazia
		Architecture:    "", // Arquitetura vazia
		HasAdminRights:  false,
		AvailableDiskGB: -1.0, // Espaço negativo
		HasInternet:     false,
		HomeDir:         "", // Home dir vazio
		CanProceed:      false,
		Issues:          []string{"Invalid environment"},
	}
}

// CreateInvalidValidationResult cria um resultado de validação inválido para testes
func CreateInvalidValidationResult() *types.ValidationResult {
	return &types.ValidationResult{
		Environment:  CreateInvalidEnvironmentInfo(),
		Dependencies: nil, // Dependencies nil
		Network:      nil, // Network nil
		Permissions:  nil, // Permissions nil
		CanProceed:   false,
		Issues: []types.ValidationIssue{
			{
				Type:        "environment",
				Severity:    "error",
				Message:     "Invalid environment",
				Suggestions: []string{"Fix environment issues"},
			},
		},
		Warnings: []string{"Multiple validation issues found"},
	}
}

// Função auxiliar para verificar se uma string contém outra
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// AssertTimeWithinRange verifica se um tempo está dentro de um intervalo
func AssertTimeWithinRange(t *testing.T, actual time.Time, expected time.Time, tolerance time.Duration) {
	t.Helper()

	diff := actual.Sub(expected)
	if diff < -tolerance || diff > tolerance {
		t.Errorf("Expected time to be within %v of %v, but got %v (diff: %v)",
			tolerance, expected, actual, diff)
	}
}

// AssertStringNotEmpty verifica se uma string não está vazia
func AssertStringNotEmpty(t *testing.T, str, name string) {
	t.Helper()

	if str == "" {
		t.Errorf("Expected %s to not be empty", name)
	}
}

// AssertStringEmpty verifica se uma string está vazia
func AssertStringEmpty(t *testing.T, str, name string) {
	t.Helper()

	if str != "" {
		t.Errorf("Expected %s to be empty, but got: %s", name, str)
	}
}

// AssertIntGreaterThan verifica se um inteiro é maior que um valor
func AssertIntGreaterThan(t *testing.T, actual, expected int, name string) {
	t.Helper()

	if actual <= expected {
		t.Errorf("Expected %s to be greater than %d, but got: %d", name, expected, actual)
	}
}

// AssertIntEqual verifica se dois inteiros são iguais
func AssertIntEqual(t *testing.T, actual, expected int, name string) {
	t.Helper()

	if actual != expected {
		t.Errorf("Expected %s to be %d, but got: %d", name, expected, actual)
	}
}

// AssertBoolEqual verifica se dois booleanos são iguais
func AssertBoolEqual(t *testing.T, actual, expected bool, name string) {
	t.Helper()

	if actual != expected {
		t.Errorf("Expected %s to be %t, but got: %t", name, expected, actual)
	}
}

// AssertFloat64Equal verifica se dois float64 são iguais (com tolerância)
func AssertFloat64Equal(t *testing.T, actual, expected, tolerance float64, name string) {
	t.Helper()

	diff := actual - expected
	if diff < -tolerance || diff > tolerance {
		t.Errorf("Expected %s to be %f (tolerance: %f), but got: %f (diff: %f)",
			name, expected, tolerance, actual, diff)
	}
}

// AssertSliceLength verifica se um slice tem o comprimento esperado
func AssertSliceLength(t *testing.T, actual interface{}, expected int, name string) {
	t.Helper()

	switch v := actual.(type) {
	case []string:
		if len(v) != expected {
			t.Errorf("Expected %s to have length %d, but got: %d", name, expected, len(v))
		}
	case []types.Dependency:
		if len(v) != expected {
			t.Errorf("Expected %s to have length %d, but got: %d", name, expected, len(v))
		}
	case []types.ValidationIssue:
		if len(v) != expected {
			t.Errorf("Expected %s to have length %d, but got: %d", name, expected, len(v))
		}
	case []int:
		if len(v) != expected {
			t.Errorf("Expected %s to have length %d, but got: %d", name, expected, len(v))
		}
	default:
		t.Errorf("Unsupported type for slice length assertion: %T", actual)
	}
}

// AssertMapNotEmpty verifica se um map não está vazio
func AssertMapNotEmpty(t *testing.T, actual map[string]string, name string) {
	t.Helper()

	if len(actual) == 0 {
		t.Errorf("Expected %s to not be empty", name)
	}
}

// AssertMapEmpty verifica se um map está vazio
func AssertMapEmpty(t *testing.T, actual map[string]string, name string) {
	t.Helper()

	if len(actual) != 0 {
		t.Errorf("Expected %s to be empty, but got: %v", name, actual)
	}
}
