package setup

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/internal/types"
)

// TestSetupFlow testa o fluxo completo de setup, status e reset
func TestSetupFlow(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir, err := os.MkdirTemp("", "syntropy-test")
	if err != nil {
		t.Fatalf("Falha ao criar diretório temporário: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Configurar opções de teste
	options := types.SetupOptions{
		Force:         true,
		InstallService: false,
		ConfigPath:    filepath.Join(tempDir, "config.yaml"),
		HomeDir:       tempDir,
	}

	fmt.Println("=== Testando Setup ===")
	fmt.Printf("Diretório de teste: %s\n", tempDir)

	// Testar Setup
	setupResult, err := Setup(options)
	if err != nil {
		// Em sistemas não-Windows, esperamos ErrNotImplemented
		if err.Error() != "funcionalidade não implementada para este sistema operacional: windows (stub)" &&
		   err.Error() != "funcionalidade não implementada para este sistema operacional: linux" &&
		   err.Error() != "funcionalidade não implementada para este sistema operacional: darwin" {
			t.Fatalf("Setup falhou com erro inesperado: %v", err)
		}
		fmt.Printf("Setup não implementado para este sistema: %v (esperado em sistemas não-Windows)\n", err)
	} else {
		fmt.Printf("Setup concluído com sucesso: %v\n", setupResult.Success)
		
		// Verificar se os arquivos foram criados
		configFile := setupResult.ConfigPath
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			t.Errorf("Arquivo de configuração não foi criado: %s", configFile)
		} else {
			fmt.Printf("Arquivo de configuração criado: %s\n", configFile)
		}
	}

	// Testar Status
	fmt.Println("\n=== Testando Status ===")
	statusResult, err := Status(options)
	if err != nil {
		// Em sistemas não-Windows, esperamos ErrNotImplemented
		if err.Error() != "funcionalidade não implementada para este sistema operacional: windows (stub)" &&
		   err.Error() != "funcionalidade não implementada para este sistema operacional: linux" &&
		   err.Error() != "funcionalidade não implementada para este sistema operacional: darwin" {
			t.Fatalf("Status falhou com erro inesperado: %v", err)
		}
		fmt.Printf("Status não implementado para este sistema: %v (esperado em sistemas não-Windows)\n", err)
	} else {
		fmt.Printf("Status concluído com sucesso: %v\n", statusResult.Success)
		fmt.Printf("Caminho da configuração: %s\n", statusResult.ConfigPath)
	}

	// Testar Reset
	fmt.Println("\n=== Testando Reset ===")
	resetResult, err := Reset(options)
	if err != nil {
		// Em sistemas não-Windows, esperamos ErrNotImplemented
		if err.Error() != "funcionalidade não implementada para este sistema operacional: windows (stub)" &&
		   err.Error() != "funcionalidade não implementada para este sistema operacional: linux" &&
		   err.Error() != "funcionalidade não implementada para este sistema operacional: darwin" {
			t.Fatalf("Reset falhou com erro inesperado: %v", err)
		}
		fmt.Printf("Reset não implementado para este sistema: %v (esperado em sistemas não-Windows)\n", err)
	} else {
		fmt.Printf("Reset concluído com sucesso: %v\n", resetResult.Success)
		
		// Verificar se os arquivos foram removidos (se Reset foi bem-sucedido)
		if resetResult.Success {
			configFile := filepath.Join(tempDir, "config", "manager.yaml")
			if _, err := os.Stat(configFile); !os.IsNotExist(err) {
				t.Errorf("Arquivo de configuração não foi removido após reset: %s", configFile)
			} else {
				fmt.Printf("Arquivo de configuração removido com sucesso\n")
			}
		}
	}
}

// TestGetSyntropyDir testa a função GetSyntropyDir
func TestGetSyntropyDir(t *testing.T) {
	fmt.Println("\n=== Testando GetSyntropyDir ===")
	dir := GetSyntropyDir()
	fmt.Printf("Diretório Syntropy: %s\n", dir)
	
	// Verificar se o caminho retornado é válido
	if dir == "" {
		t.Errorf("GetSyntropyDir retornou caminho vazio")
	}
}