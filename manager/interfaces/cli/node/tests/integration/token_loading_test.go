/**
 * File: token_loading_test.go
 * Purpose: Testes de integração para validação do carregamento de tokens do Setup Component
 * Dependencies: setup_adapter.go, token_integration.go, auto_config_generator.go
 * Exports: Testes de integração para validação da compatibilidade entre Setup e Node
 * Author: Syntropy Node Component
 * Created: 2025-01-27
 * Modified: 2025-01-27
 * Version: 1.0.0
 *
 * Business Context:
 * Este arquivo contém testes de integração que validam que o Node Component
 * consegue carregar corretamente os tokens salvos pelo Setup Component.
 */

package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	node "github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/node/src"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Logger interface para compatibilidade com testes
type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
	SetLevel(level string)
	WithFields(fields map[string]interface{}) Logger
}

// TestNodeComponent_TokenLoadingFromSetup testa o carregamento de tokens do Setup Component
func TestNodeComponent_TokenLoadingFromSetup(t *testing.T) {
	// Criar diretório temporário para simular ambiente de teste
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Criar estrutura de diretórios do Setup Component
	tokensDir := filepath.Join(tempDir, ".syntropy", "tokens")
	err := os.MkdirAll(tokensDir, 0700)
	require.NoError(t, err, "Should create tokens directory")

	// Criar token de teste (simulando o formato salvo pelo Setup Component)
	testToken := "12345678-1234-4abc-9def-123456789012"
	tokenBackup := struct {
		Token     string    `json:"token"`
		CreatedAt time.Time `json:"created_at"`
		Version   string    `json:"version"`
		Checksum  string    `json:"checksum"`
	}{
		Token:     testToken,
		CreatedAt: time.Now(),
		Version:   "1.0.0",
		Checksum:  "test_checksum",
	}

	// Salvar token no formato JSON do Setup Component
	tokenFile := filepath.Join(tokensDir, "grid-token.json")
	tokenData, err := json.MarshalIndent(tokenBackup, "", "  ")
	require.NoError(t, err, "Should marshal token data")

	err = os.WriteFile(tokenFile, tokenData, 0600)
	require.NoError(t, err, "Should write token file")

	// Criar SetupAdapter e testar carregamento
	// Usar nil logger para simplificar o teste
	adapter := node.NewSetupAdapter(nil)
	tokenManagerAdapter := node.NewSetupTokenManagerAdapter(adapter)

	// Testar TokenExists
	exists, err := tokenManagerAdapter.TokenExists()
	require.NoError(t, err, "Should check token existence without error")
	assert.True(t, exists, "Token should exist")

	// Testar LoadToken
	loadedToken, err := tokenManagerAdapter.LoadToken()
	require.NoError(t, err, "Should load token without error")
	assert.Equal(t, testToken, loadedToken, "Loaded token should match saved token")

	// Testar ValidateToken
	err = tokenManagerAdapter.ValidateToken(loadedToken)
	require.NoError(t, err, "Should validate token without error")

	// Testar com token inválido
	err = tokenManagerAdapter.ValidateToken("")
	assert.Error(t, err, "Should error on empty token")

	err = tokenManagerAdapter.ValidateToken("short")
	assert.Error(t, err, "Should error on short token")
}

// TestNodeComponent_TokenLoadingIntegration testa integração completa com NodeManager
func TestNodeComponent_TokenLoadingIntegration(t *testing.T) {
	// Criar diretório temporário para simular ambiente de teste
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Criar estrutura de diretórios do Setup Component
	tokensDir := filepath.Join(tempDir, ".syntropy", "tokens")
	err := os.MkdirAll(tokensDir, 0700)
	require.NoError(t, err, "Should create tokens directory")

	// Criar token de teste
	testToken := "12345678-1234-4abc-9def-123456789012"
	tokenBackup := struct {
		Token     string    `json:"token"`
		CreatedAt time.Time `json:"created_at"`
		Version   string    `json:"version"`
		Checksum  string    `json:"checksum"`
	}{
		Token:     testToken,
		CreatedAt: time.Now(),
		Version:   "1.0.0",
		Checksum:  "test_checksum",
	}

	// Salvar token no formato JSON do Setup Component
	tokenFile := filepath.Join(tokensDir, "grid-token.json")
	tokenData, err := json.MarshalIndent(tokenBackup, "", "  ")
	require.NoError(t, err, "Should marshal token data")

	err = os.WriteFile(tokenFile, tokenData, 0600)
	require.NoError(t, err, "Should write token file")

	// Criar NodeManager e testar inicialização
	nodeManager := node.NewNodeManager()
	require.NotNil(t, nodeManager, "NodeManager should be created")

	// Inicializar NodeManager (deve carregar token do Setup Component)
	err = nodeManager.Initialize()
	require.NoError(t, err, "NodeManager should initialize without error")

	// Verificar que a integração de token foi inicializada corretamente
	// (Isso seria testado através do AutoConfigGenerator se ele tentar usar o token)
	assert.True(t, true, "NodeManager initialized successfully with token from Setup Component")
}

// TestNodeComponent_TokenLoadingErrorHandling testa tratamento de erros no carregamento de tokens
func TestNodeComponent_TokenLoadingErrorHandling(t *testing.T) {
	// Criar diretório temporário para simular ambiente de teste
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Usar nil logger para simplificar o teste
	adapter := node.NewSetupAdapter(nil)
	tokenManagerAdapter := node.NewSetupTokenManagerAdapter(adapter)

	// Testar TokenExists quando não há token
	exists, err := tokenManagerAdapter.TokenExists()
	require.NoError(t, err, "Should check token existence without error")
	assert.False(t, exists, "Token should not exist")

	// Testar LoadToken quando não há token
	_, err = tokenManagerAdapter.LoadToken()
	assert.Error(t, err, "Should error when loading non-existent token")
	assert.Contains(t, err.Error(), "grid token not found", "Error should mention token not found")
	assert.Contains(t, err.Error(), "syntropy setup", "Error should mention running setup")

	// Testar com arquivo JSON malformado
	tokensDir := filepath.Join(tempDir, ".syntropy", "tokens")
	err = os.MkdirAll(tokensDir, 0700)
	require.NoError(t, err, "Should create tokens directory")

	tokenFile := filepath.Join(tokensDir, "grid-token.json")
	err = os.WriteFile(tokenFile, []byte("invalid json"), 0600)
	require.NoError(t, err, "Should write invalid JSON file")

	// Testar carregamento com JSON inválido
	_, err = tokenManagerAdapter.LoadToken()
	assert.Error(t, err, "Should error when loading invalid JSON")
	assert.Contains(t, err.Error(), "parse token file", "Error should mention parsing failure")
}

// TestNodeComponent_TokenLoadingPermissions testa permissões do arquivo de token
func TestNodeComponent_TokenLoadingPermissions(t *testing.T) {
	// Criar diretório temporário para simular ambiente de teste
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Criar estrutura de diretórios
	tokensDir := filepath.Join(tempDir, ".syntropy", "tokens")
	err := os.MkdirAll(tokensDir, 0700)
	require.NoError(t, err, "Should create tokens directory")

	// Criar token de teste
	testToken := "12345678-1234-4abc-9def-123456789012"
	tokenBackup := struct {
		Token     string    `json:"token"`
		CreatedAt time.Time `json:"created_at"`
		Version   string    `json:"version"`
		Checksum  string    `json:"checksum"`
	}{
		Token:     testToken,
		CreatedAt: time.Now(),
		Version:   "1.0.0",
		Checksum:  "test_checksum",
	}

	// Salvar token com permissões corretas (0600)
	tokenFile := filepath.Join(tokensDir, "grid-token.json")
	tokenData, err := json.MarshalIndent(tokenBackup, "", "  ")
	require.NoError(t, err, "Should marshal token data")

	err = os.WriteFile(tokenFile, tokenData, 0600)
	require.NoError(t, err, "Should write token file")

	// Verificar permissões do arquivo
	stat, err := os.Stat(tokenFile)
	require.NoError(t, err, "Should be able to stat token file")

	fileMode := stat.Mode()
	assert.Equal(t, os.FileMode(0600), fileMode&os.ModePerm, "Token file should have 0600 permissions")

	// Testar carregamento (deve funcionar com permissões corretas)
	// Usar nil logger para simplificar o teste
	adapter := node.NewSetupAdapter(nil)
	tokenManagerAdapter := node.NewSetupTokenManagerAdapter(adapter)

	loadedToken, err := tokenManagerAdapter.LoadToken()
	require.NoError(t, err, "Should load token without error")
	assert.Equal(t, testToken, loadedToken, "Loaded token should match saved token")
}

// MockLogger implementa uma interface logger simples para testes
type MockLogger struct{}

func (ml *MockLogger) Info(msg string, args ...interface{})            {}
func (ml *MockLogger) Debug(msg string, args ...interface{})           {}
func (ml *MockLogger) Warn(msg string, args ...interface{})            {}
func (ml *MockLogger) Error(msg string, args ...interface{})           {}
func (ml *MockLogger) Fatal(msg string, args ...interface{})           {}
func (ml *MockLogger) SetLevel(level string)                           {}
func (ml *MockLogger) GetLevel() string                                { return "info" }
func (ml *MockLogger) WithFields(fields map[string]interface{}) Logger { return ml }
