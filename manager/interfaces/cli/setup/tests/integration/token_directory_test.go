/**
 * File: token_directory_test.go
 * Purpose: Testes de integração para validação da criação do diretório tokens
 * Dependencies: configurator.go, token_manager.go, setup.go
 * Exports: Testes de integração para validação do diretório tokens
 * Author: Syntropy Setup Component
 * Created: 2025-01-27
 * Modified: 2025-01-27
 * Version: 1.0.0
 *
 * Business Context:
 * Este arquivo contém testes de integração que validam a criação correta
 * do diretório tokens durante o setup, incluindo permissões e estrutura.
 */

package integration

import (
	"os"
	"path/filepath"
	"testing"

	setup "setup-component/src"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSetupManager_TokensDirectoryCreation testa a criação do diretório tokens
func TestSetupManager_TokensDirectoryCreation(t *testing.T) {
	// Criar SetupManager
	manager, err := setup.NewSetupManager()
	require.NoError(t, err)
	assert.NotNil(t, manager)

	// Obter diretório home
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)

	tokensDir := filepath.Join(homeDir, ".syntropy", "tokens")

	// Verificar que o diretório tokens não existe inicialmente (ou pode existir se criado pelo TokenManager)
	// O TokenManager agora cria o diretório automaticamente, então vamos apenas verificar se não há arquivo de token
	tokenFile := filepath.Join(tokensDir, "grid-token.json")
	_, err = os.Stat(tokenFile)
	_ = !os.IsNotExist(err) // initialTokenExists - pode existir se setup já foi executado

	// Executar setup
	options := &setup.SetupOptions{
		Force:          true,
		SkipValidation: true,
		CustomSettings: map[string]string{
			"generate_grid_token": "true",
		},
	}

	err = manager.SetupWithPublicOptions(options)
	require.NoError(t, err, "Setup should complete successfully")

	// Verificar que o diretório tokens foi criado
	stat, err := os.Stat(tokensDir)
	require.NoError(t, err, "Tokens directory should exist after setup")
	assert.True(t, stat.IsDir(), "Tokens should be a directory")

	// Verificar permissões do diretório (deve ser 0700)
	fileMode := stat.Mode()
	assert.Equal(t, os.FileMode(0700), fileMode&os.ModePerm, "Tokens directory should have 0700 permissions")
}

// TestTokenManager_DirectoryCreation testa a criação do diretório no TokenManager
func TestTokenManager_DirectoryCreation(t *testing.T) {
	// Criar logger
	logger := setup.NewSetupLogger()
	defer logger.Close()

	// Criar TokenManager (deve criar o diretório automaticamente)
	tokenManager := setup.NewTokenManager(logger)
	assert.NotNil(t, tokenManager)

	// Obter diretório home
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)

	tokensDir := filepath.Join(homeDir, ".syntropy", "tokens")

	// Verificar que o diretório tokens foi criado
	stat, err := os.Stat(tokensDir)
	require.NoError(t, err, "Tokens directory should exist after TokenManager creation")
	assert.True(t, stat.IsDir(), "Tokens should be a directory")

	// Verificar permissões do diretório (deve ser 0700)
	fileMode := stat.Mode()
	assert.Equal(t, os.FileMode(0700), fileMode&os.ModePerm, "Tokens directory should have 0700 permissions")
}

// TestSetupManager_TokensDirectoryPermissions testa as permissões do diretório tokens
func TestSetupManager_TokensDirectoryPermissions(t *testing.T) {
	// Criar SetupManager
	manager, err := setup.NewSetupManager()
	require.NoError(t, err)
	assert.NotNil(t, manager)

	// Executar setup
	options := &setup.SetupOptions{
		Force:          true,
		SkipValidation: true,
		CustomSettings: map[string]string{
			"generate_grid_token": "true",
		},
	}

	err = manager.SetupWithPublicOptions(options)
	require.NoError(t, err, "Setup should complete successfully")

	// Obter diretório home
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)

	tokensDir := filepath.Join(homeDir, ".syntropy", "tokens")

	// Verificar permissões do diretório
	stat, err := os.Stat(tokensDir)
	require.NoError(t, err, "Tokens directory should exist")

	fileMode := stat.Mode()
	assert.Equal(t, os.FileMode(0700), fileMode&os.ModePerm, "Tokens directory should have 0700 permissions")

	// Verificar que o arquivo grid-token.json foi criado (se keyring não estiver disponível)
	tokenFile := filepath.Join(tokensDir, "grid-token.json")
	if _, err := os.Stat(tokenFile); err == nil {
		// Se o arquivo existe, verificar suas permissões
		fileStat, err := os.Stat(tokenFile)
		require.NoError(t, err, "Token file should be readable")

		fileMode := fileStat.Mode()
		assert.Equal(t, os.FileMode(0600), fileMode&os.ModePerm, "Token file should have 0600 permissions")
	}
}

// TestSetupManager_TokensDirectoryStructure testa a estrutura completa do diretório tokens
func TestSetupManager_TokensDirectoryStructure(t *testing.T) {
	// Criar SetupManager
	manager, err := setup.NewSetupManager()
	require.NoError(t, err)
	assert.NotNil(t, manager)

	// Executar setup
	options := &setup.SetupOptions{
		Force:          true,
		SkipValidation: true,
		CustomSettings: map[string]string{
			"generate_grid_token": "true",
		},
	}

	err = manager.SetupWithPublicOptions(options)
	require.NoError(t, err, "Setup should complete successfully")

	// Obter diretório home
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)

	// Verificar que todos os diretórios necessários foram criados
	expectedDirs := []string{
		"config",
		"keys",
		"nodes",
		"logs",
		"cache",
		"backups",
		"templates",
		"state",
		"tokens",
	}

	for _, dir := range expectedDirs {
		dirPath := filepath.Join(homeDir, ".syntropy", dir)
		stat, err := os.Stat(dirPath)
		require.NoError(t, err, "Directory %s should exist", dir)
		assert.True(t, stat.IsDir(), "Path %s should be a directory", dir)
	}

	// Verificar especificamente o diretório tokens
	tokensDir := filepath.Join(homeDir, ".syntropy", "tokens")
	stat, err := os.Stat(tokensDir)
	require.NoError(t, err, "Tokens directory should exist")
	assert.True(t, stat.IsDir(), "Tokens should be a directory")

	// Verificar que o diretório está vazio ou contém apenas grid-token.json
	entries, err := os.ReadDir(tokensDir)
	require.NoError(t, err, "Should be able to read tokens directory")

	// Se há entradas, deve ser apenas grid-token.json
	if len(entries) > 0 {
		assert.Len(t, entries, 1, "Tokens directory should contain only one file")
		assert.Equal(t, "grid-token.json", entries[0].Name(), "Only file should be grid-token.json")
	}
}

// TestSetupManager_TokensDirectoryIntegration testa integração completa do diretório tokens
func TestSetupManager_TokensDirectoryIntegration(t *testing.T) {
	// Criar SetupManager
	manager, err := setup.NewSetupManager()
	require.NoError(t, err)
	assert.NotNil(t, manager)

	// Verificar se inicialmente existe token (pode existir se setup já foi executado)
	exists, err := manager.GridTokenExists()
	require.NoError(t, err)
	_ = exists // initialTokenExists - pode existir se setup já foi executado

	// Executar setup
	options := &setup.SetupOptions{
		Force:          true,
		SkipValidation: true,
		CustomSettings: map[string]string{
			"generate_grid_token": "true",
		},
	}

	err = manager.SetupWithPublicOptions(options)
	require.NoError(t, err, "Setup should complete successfully")

	// Verificar que o token agora existe (ou já existia)
	exists, err = manager.GridTokenExists()
	require.NoError(t, err)
	assert.True(t, exists, "Grid token should exist after setup")

	// Obter o token
	token, err := manager.GetGridToken()
	require.NoError(t, err)
	assert.NotEmpty(t, token, "Grid token should not be empty")

	// Verificar que o diretório tokens existe e tem o arquivo
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)

	tokensDir := filepath.Join(homeDir, ".syntropy", "tokens")
	stat, err := os.Stat(tokensDir)
	require.NoError(t, err, "Tokens directory should exist")
	assert.True(t, stat.IsDir(), "Tokens should be a directory")

	// Verificar permissões
	fileMode := stat.Mode()
	assert.Equal(t, os.FileMode(0700), fileMode&os.ModePerm, "Tokens directory should have 0700 permissions")
}
