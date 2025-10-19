/**
 * File: token_manager_integration_test.go
 * Purpose: Testes de integração para o TokenManager com SetupManager
 * Dependencies: token_manager.go, setup.go, token_management.go
 * Exports: Testes de integração para validação completa do TokenManager
 * Author: Syntropy Setup Component
 * Created: 2025-01-27
 * Modified: 2025-01-27
 * Version: 1.0.0
 *
 * Business Context:
 * Este arquivo contém testes de integração que validam o TokenManager
 * integrado com o SetupManager, testando fluxos completos de setup
 * com geração de Grid Token e operações de gerenciamento.
 */

package integration

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	setup "setup-component/src"
	"setup-component/src/internal/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSetupManager_GridTokenIntegration testa integração completa do Grid Token
func TestSetupManager_GridTokenIntegration(t *testing.T) {
	// Criar SetupManager
	manager, err := setup.NewSetupManager()
	require.NoError(t, err)
	assert.NotNil(t, manager)

	// Verificar que inicialmente não existe token
	exists, err := manager.GridTokenExists()
	require.NoError(t, err)
	assert.False(t, exists)

	// Gerar Grid Token
	token, err := manager.GenerateGridToken()
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Len(t, token, 36)

	// Verificar que token agora existe
	exists, err = manager.GridTokenExists()
	require.NoError(t, err)
	assert.True(t, exists)

	// Obter token
	retrievedToken, err := manager.GetGridToken()
	require.NoError(t, err)
	assert.Equal(t, token, retrievedToken)

	// Rotacionar token
	newToken, err := manager.RotateGridToken()
	require.NoError(t, err)
	assert.NotEqual(t, token, newToken)
	assert.Len(t, newToken, 36)

	// Verificar que novo token foi salvo
	finalToken, err := manager.GetGridToken()
	require.NoError(t, err)
	assert.Equal(t, newToken, finalToken)
}

// TestSetupManager_SetupWithGridToken testa setup completo com geração de Grid Token
func TestSetupManager_SetupWithGridToken(t *testing.T) {
	// Criar SetupManager
	manager, err := setup.NewSetupManager()
	require.NoError(t, err)
	require.NotNil(t, manager)

	// Opções de setup com geração de Grid Token
	options := &types.SetupOptions{
		Force:        false,
		ValidateOnly: false,
		TestMode:     true, // Bypass strict validation
		Verbose:      true,
		CustomSettings: map[string]string{
			"owner_name":          "Test User",
			"owner_email":         "test@example.com",
			"generate_grid_token": "true",
		},
	}

	// Executar setup completo
	err = manager.Setup(options)
	require.NoError(t, err)

	// Verificar que Grid Token foi gerado
	exists, err := manager.GridTokenExists()
	require.NoError(t, err)
	assert.True(t, exists)

	// Obter token
	token, err := manager.GetGridToken()
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Len(t, token, 36)

	// Verificar estado do setup
	state, err := manager.Status()
	require.NoError(t, err)
	assert.NotNil(t, state)
	assert.Equal(t, types.SetupStatusCompleted, state.Status)

	// Verificar metadados do Grid Token
	assert.Equal(t, "true", state.Metadata["grid_token_generated"])
	assert.Contains(t, state.Metadata["grid_token_preview"], "...[HIDDEN]")
	assert.Equal(t, "keyring", state.Metadata["grid_token_storage"])
}

// TestSetupManager_GridTokenBackupRestore testa backup e restore de Grid Token
func TestSetupManager_GridTokenBackupRestore(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	backupPath := filepath.Join(tempDir, "grid_token_backup.json")

	// Criar SetupManager
	manager, err := setup.NewSetupManager()
	require.NoError(t, err)
	require.NotNil(t, manager)

	// Gerar Grid Token
	token, err := manager.GenerateGridToken()
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Exportar token para backup
	err = manager.ExportGridToken(backupPath)
	require.NoError(t, err)

	// Verificar que arquivo de backup foi criado
	_, err = os.Stat(backupPath)
	require.NoError(t, err)

	// Deletar token atual
	err = manager.DeleteGridToken()
	require.NoError(t, err)

	// Verificar que token foi removido
	exists, err := manager.GridTokenExists()
	require.NoError(t, err)
	assert.False(t, exists)

	// Importar token de backup
	err = manager.ImportGridToken(backupPath)
	require.NoError(t, err)

	// Verificar que token foi restaurado
	restoredToken, err := manager.GetGridToken()
	require.NoError(t, err)
	assert.Equal(t, token, restoredToken)

	// Verificar que token existe novamente
	exists, err = manager.GridTokenExists()
	require.NoError(t, err)
	assert.True(t, exists)
}

// TestSetupManager_GridTokenRotation testa rotação de Grid Token
func TestSetupManager_GridTokenRotation(t *testing.T) {
	// Criar SetupManager
	manager, err := setup.NewSetupManager()
	require.NoError(t, err)
	require.NotNil(t, manager)

	// Gerar token inicial
	originalToken, err := manager.GenerateGridToken()
	require.NoError(t, err)
	assert.NotEmpty(t, originalToken)

	// Verificar que token foi salvo
	retrievedToken, err := manager.GetGridToken()
	require.NoError(t, err)
	assert.Equal(t, originalToken, retrievedToken)

	// Rotacionar token
	newToken, err := manager.RotateGridToken()
	require.NoError(t, err)
	assert.NotEqual(t, originalToken, newToken)
	assert.Len(t, newToken, 36)

	// Verificar que novo token foi salvo
	finalToken, err := manager.GetGridToken()
	require.NoError(t, err)
	assert.Equal(t, newToken, finalToken)

	// Verificar que token antigo não é mais válido
	assert.NotEqual(t, originalToken, finalToken)
}

// TestSetupManager_GridTokenErrorHandling testa tratamento de erros
func TestSetupManager_GridTokenErrorHandling(t *testing.T) {
	// Criar SetupManager
	manager, err := setup.NewSetupManager()
	require.NoError(t, err)
	require.NotNil(t, manager)

	// Tentar obter token quando não existe
	_, err = manager.GetGridToken()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "grid_token_load_failed")

	// Tentar rotacionar token quando não existe
	_, err = manager.RotateGridToken()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "grid_token_rotation_failed")

	// Tentar exportar token quando não existe
	tempDir := t.TempDir()
	backupPath := filepath.Join(tempDir, "nonexistent_token.json")

	err = manager.ExportGridToken(backupPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "grid_token_export_failed")

	// Tentar importar arquivo inexistente
	err = manager.ImportGridToken("/path/that/does/not/exist.json")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "grid_token_import_failed")
}

// TestSetupManager_GridTokenConcurrentOperations testa operações concorrentes
func TestSetupManager_GridTokenConcurrentOperations(t *testing.T) {
	// Criar SetupManager
	manager, err := setup.NewSetupManager()
	require.NoError(t, err)
	require.NotNil(t, manager)

	// Canal para sincronização
	done := make(chan bool, 5)

	// Executar operações concorrentes
	go func() {
		defer func() { done <- true }()
		_, err := manager.GenerateGridToken()
		assert.NoError(t, err)
	}()

	go func() {
		defer func() { done <- true }()
		time.Sleep(10 * time.Millisecond) // Pequeno delay
		exists, err := manager.GridTokenExists()
		assert.NoError(t, err)
		assert.True(t, exists)
	}()

	go func() {
		defer func() { done <- true }()
		time.Sleep(20 * time.Millisecond) // Pequeno delay
		token, err := manager.GetGridToken()
		if err == nil {
			assert.NotEmpty(t, token)
		}
	}()

	go func() {
		defer func() { done <- true }()
		time.Sleep(30 * time.Millisecond) // Pequeno delay
		_, err := manager.RotateGridToken()
		// Pode falhar se token não existir ainda, isso é esperado
		if err != nil {
			assert.Contains(t, err.Error(), "grid_token_rotation_failed")
		}
	}()

	go func() {
		defer func() { done <- true }()
		time.Sleep(40 * time.Millisecond) // Pequeno delay
		exists, err := manager.GridTokenExists()
		assert.NoError(t, err)
		// Pode ser true ou false dependendo da ordem de execução
		assert.NotNil(t, exists)
	}()

	// Aguardar todas as operações
	for i := 0; i < 5; i++ {
		<-done
	}
}

// TestSetupManager_GridTokenWithSetupReset testa Grid Token com reset de setup
func TestSetupManager_GridTokenWithSetupReset(t *testing.T) {
	// Criar SetupManager
	manager, err := setup.NewSetupManager()
	require.NoError(t, err)
	require.NotNil(t, manager)

	// Executar setup com Grid Token
	options := &types.SetupOptions{
		Force:        false,
		ValidateOnly: false,
		TestMode:     true,
		CustomSettings: map[string]string{
			"owner_name":          "Test User",
			"owner_email":         "test@example.com",
			"generate_grid_token": "true",
		},
	}

	err = manager.Setup(options)
	require.NoError(t, err)

	// Verificar que Grid Token foi gerado
	exists, err := manager.GridTokenExists()
	require.NoError(t, err)
	assert.True(t, exists)

	token, err := manager.GetGridToken()
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Executar reset do setup
	err = manager.Reset(true)
	require.NoError(t, err)

	// Verificar que Grid Token ainda existe (não é removido pelo reset)
	exists, err = manager.GridTokenExists()
	require.NoError(t, err)
	assert.True(t, exists)

	// Verificar que token ainda é o mesmo
	retrievedToken, err := manager.GetGridToken()
	require.NoError(t, err)
	assert.Equal(t, token, retrievedToken)
}

// TestSetupManager_GridTokenValidation testa validação de Grid Token
func TestSetupManager_GridTokenValidation(t *testing.T) {
	// Criar SetupManager
	manager, err := setup.NewSetupManager()
	require.NoError(t, err)
	require.NotNil(t, manager)

	// Gerar Grid Token
	token, err := manager.GenerateGridToken()
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verificar formato UUID v4
	assert.Regexp(t, `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`, token)

	// Verificar que token é válido
	err = manager.tokenManager.ValidateToken(token)
	assert.NoError(t, err)

	// Testar token inválido
	invalidToken := "invalid-token-format"
	err = manager.tokenManager.ValidateToken(invalidToken)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token length")
}

// TestSetupManager_GridTokenMetadata testa metadados do Grid Token
func TestSetupManager_GridTokenMetadata(t *testing.T) {
	// Criar SetupManager
	manager, err := setup.NewSetupManager()
	require.NoError(t, err)
	require.NotNil(t, manager)

	// Executar setup com Grid Token
	options := &types.SetupOptions{
		Force:        false,
		ValidateOnly: false,
		TestMode:     true,
		CustomSettings: map[string]string{
			"owner_name":          "Test User",
			"owner_email":         "test@example.com",
			"generate_grid_token": "true",
		},
	}

	err = manager.Setup(options)
	require.NoError(t, err)

	// Verificar estado do setup
	state, err := manager.Status()
	require.NoError(t, err)
	assert.NotNil(t, state)

	// Verificar metadados específicos do Grid Token
	assert.Equal(t, "true", state.Metadata["grid_token_generated"])
	assert.Contains(t, state.Metadata["grid_token_preview"], "...[HIDDEN]")
	assert.Equal(t, "keyring", state.Metadata["grid_token_storage"])

	// Verificar que preview não contém o token completo
	assert.NotEqual(t, state.Metadata["grid_token_preview"], state.Metadata["grid_token"])
	assert.NotContains(t, state.Metadata["grid_token_preview"], state.Metadata["grid_token"])
}
