package helpers

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	setup "setup-component/src"
)

// CreateTempDir cria um diretório temporário para testes
func CreateTempDir(t *testing.T) string {
	tempDir := t.TempDir()
	return tempDir
}

// SetupTestEnvironment configura o ambiente de teste
func SetupTestEnvironment(t *testing.T) (string, func()) {
	tempDir := CreateTempDir(t)
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)

	cleanup := func() {
		os.Setenv("HOME", originalHome)
	}

	return tempDir, cleanup
}

// CreateTestLogger cria um logger para testes
func CreateTestLogger() *setup.SetupLogger {
	return setup.NewSetupLogger()
}

// CreateTestConfigOptions cria opções de configuração para testes
func CreateTestConfigOptions() *setup.ConfigOptions {
	return &setup.ConfigOptions{
		OwnerName:      "Test User",
		OwnerEmail:     "test@example.com",
		CustomSettings: make(map[string]string),
	}
}

// CreateTestContext cria um contexto para testes
func CreateTestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 30*time.Second)
}

// GetTestDataPath retorna o caminho para dados de teste
func GetTestDataPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "testdata")
}

// CreateTestFile cria um arquivo de teste
func CreateTestFile(t *testing.T, dir, filename, content string) string {
	filePath := filepath.Join(dir, filename)
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	return filePath
}

// CreateTestDir cria um diretório de teste
func CreateTestDir(t *testing.T, parentDir, dirName string) string {
	dirPath := filepath.Join(parentDir, dirName)
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	return dirPath
}

// CleanupTestFiles remove arquivos de teste
func CleanupTestFiles(t *testing.T, paths ...string) {
	for _, path := range paths {
		if err := os.RemoveAll(path); err != nil {
			t.Logf("Warning: Failed to cleanup %s: %v", path, err)
		}
	}
}

// AssertFileExists verifica se um arquivo existe
func AssertFileExists(t *testing.T, filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("Expected file %s to exist, but it doesn't", filePath)
	}
}

// AssertDirExists verifica se um diretório existe
func AssertDirExists(t *testing.T, dirPath string) {
	if info, err := os.Stat(dirPath); os.IsNotExist(err) || !info.IsDir() {
		t.Errorf("Expected directory %s to exist, but it doesn't", dirPath)
	}
}

// GetPlatformSpecificTestData retorna dados de teste específicos da plataforma
func GetPlatformSpecificTestData() map[string]interface{} {
	return map[string]interface{}{
		"os":   runtime.GOOS,
		"arch": runtime.GOARCH,
	}
}
