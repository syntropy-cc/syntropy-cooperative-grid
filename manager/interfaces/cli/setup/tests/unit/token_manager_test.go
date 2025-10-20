/**
 * File: token_manager_test.go
 * Purpose: Testes unitários para o TokenManager
 * Dependencies: token_manager.go, internal/types
 * Exports: Testes para validação do TokenManager
 * Author: Syntropy Setup Component
 * Created: 2025-01-27
 * Modified: 2025-01-27
 * Version: 1.0.0
 *
 * Business Context:
 * Este arquivo contém testes unitários abrangentes para o TokenManager,
 * validando geração, armazenamento, carregamento, rotação e backup/restore
 * de Grid Tokens com diferentes cenários de teste.
 */

package unit

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	setup "setup-component/src"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockLogger implementa SetupLogger para testes
type MockLogger struct {
	steps    []map[string]interface{}
	errors   []error
	warnings []string
	infos    []string
}

func NewMockLogger() *MockLogger {
	return &MockLogger{
		steps:    make([]map[string]interface{}, 0),
		errors:   make([]error, 0),
		warnings: make([]string, 0),
		infos:    make([]string, 0),
	}
}

func (ml *MockLogger) LogStep(step string, data map[string]interface{}) {
	ml.steps = append(ml.steps, map[string]interface{}{
		"step": step,
		"data": data,
	})
}

func (ml *MockLogger) LogError(err error, context map[string]interface{}) {
	ml.errors = append(ml.errors, err)
}

func (ml *MockLogger) LogWarning(message string, data map[string]interface{}) {
	ml.warnings = append(ml.warnings, message)
}

func (ml *MockLogger) LogInfo(message string, data map[string]interface{}) {
	ml.infos = append(ml.infos, message)
}

func (ml *MockLogger) LogDebug(message string, data map[string]interface{}) {
	// Implementação vazia para testes
}

func (ml *MockLogger) ExportLogs(format string, outputPath string) error {
	return nil
}

func (ml *MockLogger) Close() error {
	return nil
}

// TestTokenManager_NewTokenManager testa a criação do TokenManager
func TestTokenManager_NewTokenManager(t *testing.T) {
	logger := NewMockLogger()
	tm := setup.NewTokenManager(logger)

	assert.NotNil(t, tm)
	assert.NotNil(t, logger)
}

// TestTokenManager_GenerateToken testa a geração de tokens
func TestTokenManager_GenerateToken(t *testing.T) {
	logger := NewMockLogger()
	tm := setup.NewTokenManager(logger)

	// Testar geração de token
	token, err := tm.GenerateToken()
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Len(t, token, 36) // UUID v4 tem 36 caracteres

	// Verificar formato UUID v4
	assert.Regexp(t, `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`, token)

	// Verificar que foi logado
	assert.Greater(t, len(logger.steps), 0)
	assert.Equal(t, "token_generation_start", logger.steps[0]["step"])
}

// TestTokenManager_ValidateToken testa a validação de tokens
func TestTokenManager_ValidateToken(t *testing.T) {
	logger := NewMockLogger()
	tm := setup.NewTokenManager(logger)

	// Testar token válido
	validToken := "12345678-1234-4567-8901-123456789012"
	err := tm.ValidateToken(validToken)
	assert.NoError(t, err)

	// Testar token inválido (comprimento errado)
	invalidToken := "12345678-1234-4567-8901"
	err = tm.ValidateToken(invalidToken)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "invalid token length")
	}

	// Testar token inválido (formato errado)
	invalidToken2 := "12345678-1234-4567-8901-12345678901g"
	err = tm.ValidateToken(invalidToken2)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "invalid UUID v4 format")
	}
}

// TestTokenManager_SaveAndLoadToken testa salvamento e carregamento de tokens
func TestTokenManager_SaveAndLoadToken(t *testing.T) {
	logger := NewMockLogger()
	tm := setup.NewTokenManager(logger)

	// Gerar token
	token, err := tm.GenerateToken()
	require.NoError(t, err)

	// Salvar token (vai usar fallback para arquivo)
	err = tm.SaveToken(token)
	require.NoError(t, err)

	// Carregar token
	loadedToken, err := tm.LoadToken()
	require.NoError(t, err)
	assert.Equal(t, token, loadedToken)

	// Verificar logs
	assert.Greater(t, len(logger.steps), 0)
}

// TestTokenManager_TokenExists testa verificação de existência de token
func TestTokenManager_TokenExists(t *testing.T) {
	logger := NewMockLogger()
	tm := setup.NewTokenManager(logger)

	// Inicialmente não deve existir token
	exists, err := tm.TokenExists()
	require.NoError(t, err)
	assert.False(t, exists)

	// Gerar e salvar token
	token, err := tm.GenerateToken()
	require.NoError(t, err)

	err = tm.SaveToken(token)
	require.NoError(t, err)

	// Agora deve existir
	exists, err = tm.TokenExists()
	require.NoError(t, err)
	assert.True(t, exists)
}

// TestTokenManager_RotateToken testa rotação de tokens
func TestTokenManager_RotateToken(t *testing.T) {
	logger := NewMockLogger()
	tm := setup.NewTokenManager(logger)

	// Gerar token inicial
	originalToken, err := tm.GenerateToken()
	require.NoError(t, err)

	err = tm.SaveToken(originalToken)
	require.NoError(t, err)

	// Rotacionar token
	newToken, err := tm.RotateToken()
	require.NoError(t, err)
	assert.NotEqual(t, originalToken, newToken)
	assert.Len(t, newToken, 36)

	// Verificar que o novo token foi salvo
	loadedToken, err := tm.LoadToken()
	require.NoError(t, err)
	assert.Equal(t, newToken, loadedToken)
}

// TestTokenManager_ExportImportToken testa exportação e importação de tokens
func TestTokenManager_ExportImportToken(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	exportPath := filepath.Join(tempDir, "token_backup.json")

	logger := NewMockLogger()
	tm := setup.NewTokenManager(logger)

	// Gerar e salvar token
	token, err := tm.GenerateToken()
	require.NoError(t, err)

	err = tm.SaveToken(token)
	require.NoError(t, err)

	// Exportar token
	err = tm.ExportToken(exportPath)
	require.NoError(t, err)

	// Verificar que arquivo foi criado
	_, err = os.Stat(exportPath)
	require.NoError(t, err)

	// Ler e verificar conteúdo do backup
	backupData, err := os.ReadFile(exportPath)
	require.NoError(t, err)

	var backup map[string]interface{}
	err = json.Unmarshal(backupData, &backup)
	require.NoError(t, err)
	assert.Equal(t, token, backup["token"])
	assert.NotEmpty(t, backup["checksum"])

	// Deletar token atual
	err = tm.DeleteToken()
	require.NoError(t, err)

	// Verificar que token foi removido
	exists, err := tm.TokenExists()
	require.NoError(t, err)
	assert.False(t, exists)

	// Importar token
	err = tm.ImportToken(exportPath)
	require.NoError(t, err)

	// Verificar que token foi restaurado
	loadedToken, err := tm.LoadToken()
	require.NoError(t, err)
	assert.Equal(t, token, loadedToken)
}

// TestTokenManager_DeleteToken testa remoção de tokens
func TestTokenManager_DeleteToken(t *testing.T) {
	logger := NewMockLogger()
	tm := setup.NewTokenManager(logger)

	// Gerar e salvar token
	token, err := tm.GenerateToken()
	require.NoError(t, err)

	err = tm.SaveToken(token)
	require.NoError(t, err)

	// Verificar que token existe
	exists, err := tm.TokenExists()
	require.NoError(t, err)
	assert.True(t, exists)

	// Deletar token
	err = tm.DeleteToken()
	require.NoError(t, err)

	// Verificar que token foi removido
	exists, err = tm.TokenExists()
	require.NoError(t, err)
	assert.False(t, exists)
}

// TestTokenManager_ChecksumValidation testa validação de checksum
func TestTokenManager_ChecksumValidation(t *testing.T) {
	tempDir := t.TempDir()
	exportPath := filepath.Join(tempDir, "token_backup.json")

	logger := NewMockLogger()
	tm := setup.NewTokenManager(logger)

	// Gerar e salvar token
	token, err := tm.GenerateToken()
	require.NoError(t, err)

	err = tm.SaveToken(token)
	require.NoError(t, err)

	// Exportar token
	err = tm.ExportToken(exportPath)
	require.NoError(t, err)

	// Corromper o arquivo de backup
	corruptedData := []byte(`{"token":"corrupted-token","created_at":"2025-01-27T00:00:00Z","version":"1.0.0","checksum":"invalid-checksum"}`)
	err = os.WriteFile(exportPath, corruptedData, 0600)
	require.NoError(t, err)

	// Tentar importar token corrompido
	err = tm.ImportToken(exportPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "checksum validation failed")
}

// TestTokenManager_MultipleTokens testa geração de múltiplos tokens únicos
func TestTokenManager_MultipleTokens(t *testing.T) {
	logger := NewMockLogger()
	tm := setup.NewTokenManager(logger)

	tokens := make(map[string]bool)

	// Gerar 100 tokens e verificar que são únicos
	for i := 0; i < 100; i++ {
		token, err := tm.GenerateToken()
		require.NoError(t, err)

		// Verificar que token não foi gerado antes
		assert.False(t, tokens[token], "Token duplicado gerado: %s", token)
		tokens[token] = true

		// Verificar formato
		assert.Len(t, token, 36)
		assert.Regexp(t, `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`, token)
	}

	assert.Len(t, tokens, 100)
}

// TestTokenManager_ConcurrentAccess testa acesso concorrente ao TokenManager
func TestTokenManager_ConcurrentAccess(t *testing.T) {
	logger := NewMockLogger()
	tm := setup.NewTokenManager(logger)

	// Canal para sincronização
	done := make(chan bool, 10)

	// Executar 10 operações concorrentes
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()

			// Gerar token
			token, err := tm.GenerateToken()
			require.NoError(t, err)

			// Validar token
			err = tm.ValidateToken(token)
			require.NoError(t, err)
		}()
	}

	// Aguardar todas as operações
	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestTokenManager_ErrorHandling testa tratamento de erros
func TestTokenManager_ErrorHandling(t *testing.T) {
	logger := NewMockLogger()
	tm := setup.NewTokenManager(logger)

	// Testar importação de arquivo inexistente
	err := tm.ImportToken("/path/that/does/not/exist.json")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read backup file")

	// Testar importação de arquivo com JSON inválido
	tempDir := t.TempDir()
	invalidPath := filepath.Join(tempDir, "invalid.json")

	err = os.WriteFile(invalidPath, []byte("invalid json"), 0600)
	require.NoError(t, err)

	err = tm.ImportToken(invalidPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal backup")
}

// TestTokenManager_BackupIntegrity testa integridade de backups
func TestTokenManager_BackupIntegrity(t *testing.T) {
	tempDir := t.TempDir()
	exportPath := filepath.Join(tempDir, "token_backup.json")

	logger := NewMockLogger()
	tm := setup.NewTokenManager(logger)

	// Gerar e salvar token
	token, err := tm.GenerateToken()
	require.NoError(t, err)

	err = tm.SaveToken(token)
	require.NoError(t, err)

	// Exportar token
	err = tm.ExportToken(exportPath)
	require.NoError(t, err)

	// Ler backup e verificar estrutura
	backupData, err := os.ReadFile(exportPath)
	require.NoError(t, err)

	var backup map[string]interface{}
	err = json.Unmarshal(backupData, &backup)
	require.NoError(t, err)

	// Verificar campos obrigatórios
	assert.Equal(t, token, backup["token"])
	assert.NotNil(t, backup["created_at"])
	assert.Equal(t, "1.0.0", backup["version"])
	assert.NotEmpty(t, backup["checksum"])

	// Verificar que CreatedAt é recente
	createdAtStr, ok := backup["created_at"].(string)
	require.True(t, ok)
	createdAt, err := time.Parse(time.RFC3339, createdAtStr)
	require.NoError(t, err)
	assert.WithinDuration(t, time.Now(), createdAt, time.Minute)
}
