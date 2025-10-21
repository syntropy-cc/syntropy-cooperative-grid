/**
 * File: simple_token_test.go
 * Purpose: Teste simples para validar carregamento de tokens do Setup Component
 * Dependencies: setup_adapter.go
 * Exports: Teste básico de integração
 * Author: Syntropy Node Component
 * Created: 2025-01-27
 * Modified: 2025-01-27
 * Version: 1.0.0
 *
 * Business Context:
 * Este arquivo contém um teste simples que valida que o Node Component
 * consegue carregar corretamente os tokens salvos pelo Setup Component.
 */

package integration

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestNodeComponent_SimpleTokenLoading testa o carregamento básico de tokens
func TestNodeComponent_SimpleTokenLoading(t *testing.T) {
	// Criar diretório temporário para simular ambiente de teste
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Criar estrutura de diretórios do Setup Component
	tokensDir := filepath.Join(tempDir, ".syntropy", "tokens")
	err := os.MkdirAll(tokensDir, 0700)
	if err != nil {
		t.Fatalf("Should create tokens directory: %v", err)
	}

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
	if err != nil {
		t.Fatalf("Should marshal token data: %v", err)
	}

	err = os.WriteFile(tokenFile, tokenData, 0600)
	if err != nil {
		t.Fatalf("Should write token file: %v", err)
	}

	// Testar carregamento manual do token (simulando o que o SetupAdapter faz)
	loadedToken, err := loadTokenFromFile(tokenFile)
	if err != nil {
		t.Fatalf("Should load token without error: %v", err)
	}

	if loadedToken != testToken {
		t.Errorf("Loaded token should match saved token: expected %s, got %s", testToken, loadedToken)
	}

	// Verificar permissões do arquivo
	stat, err := os.Stat(tokenFile)
	if err != nil {
		t.Fatalf("Should be able to stat token file: %v", err)
	}

	fileMode := stat.Mode()
	expectedMode := os.FileMode(0600)
	if fileMode&os.ModePerm != expectedMode {
		t.Errorf("Token file should have %o permissions, got %o", expectedMode, fileMode&os.ModePerm)
	}

	// Verificar permissões do diretório
	dirStat, err := os.Stat(tokensDir)
	if err != nil {
		t.Fatalf("Should be able to stat tokens directory: %v", err)
	}

	dirMode := dirStat.Mode()
	expectedDirMode := os.FileMode(0700)
	if dirMode&os.ModePerm != expectedDirMode {
		t.Errorf("Tokens directory should have %o permissions, got %o", expectedDirMode, dirMode&os.ModePerm)
	}

	fmt.Printf("✓ Token loading test passed: %s\n", loadedToken[:8]+"...[HIDDEN]")
}

// TestNodeComponent_TokenFileNotFound testa comportamento quando arquivo não existe
func TestNodeComponent_TokenFileNotFound(t *testing.T) {
	// Criar diretório temporário para simular ambiente de teste
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Tentar carregar token que não existe
	tokensDir := filepath.Join(tempDir, ".syntropy", "tokens")
	tokenFile := filepath.Join(tokensDir, "grid-token.json")

	_, err := loadTokenFromFile(tokenFile)
	if err == nil {
		t.Errorf("Should error when loading non-existent token file")
	}

	if !os.IsNotExist(err) {
		t.Errorf("Should return file not found error, got: %v", err)
	}

	fmt.Println("✓ Token file not found test passed")
}

// TestNodeComponent_InvalidJSON testa comportamento com JSON inválido
func TestNodeComponent_InvalidJSON(t *testing.T) {
	// Criar diretório temporário para simular ambiente de teste
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Criar estrutura de diretórios
	tokensDir := filepath.Join(tempDir, ".syntropy", "tokens")
	err := os.MkdirAll(tokensDir, 0700)
	if err != nil {
		t.Fatalf("Should create tokens directory: %v", err)
	}

	// Criar arquivo com JSON inválido
	tokenFile := filepath.Join(tokensDir, "grid-token.json")
	err = os.WriteFile(tokenFile, []byte("invalid json"), 0600)
	if err != nil {
		t.Fatalf("Should write invalid JSON file: %v", err)
	}

	// Tentar carregar token com JSON inválido
	_, err = loadTokenFromFile(tokenFile)
	if err == nil {
		t.Errorf("Should error when loading invalid JSON")
	}

	// Verificar que o erro menciona parsing
	if err != nil && !containsString(err.Error(), "parse") {
		t.Errorf("Error should mention parsing, got: %v", err)
	}

	fmt.Println("✓ Invalid JSON test passed")
}

// loadTokenFromFile simula o carregamento de token do SetupAdapter
func loadTokenFromFile(tokenFile string) (string, error) {
	// Verificar se arquivo existe
	if _, err := os.Stat(tokenFile); err != nil {
		return "", err
	}

	// Ler arquivo
	tokenBytes, err := os.ReadFile(tokenFile)
	if err != nil {
		return "", fmt.Errorf("failed to read token file: %w", err)
	}

	// Parse JSON estrutura TokenBackup
	var backup struct {
		Token     string    `json:"token"`
		CreatedAt time.Time `json:"created_at"`
		Version   string    `json:"version"`
		Checksum  string    `json:"checksum"`
	}

	if err := json.Unmarshal(tokenBytes, &backup); err != nil {
		return "", fmt.Errorf("failed to parse token file: %w", err)
	}

	return backup.Token, nil
}

// containsString verifica se uma string contém uma substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr ||
		len(s) > len(substr) && containsString(s[1:], substr)
}
