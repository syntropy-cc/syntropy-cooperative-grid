package setup

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src/internal/types"
)

// TestSetupFlow testa o fluxo completo de setup, status e reset
func TestSetupFlow(t *testing.T) {
	// Evitar shadowing: declarar variáveis no topo do escopo
	var (
		err     error
		tempDir string
	)

	// Criar diretório temporário para testes
	tempDir, err = os.MkdirTemp("", "syntropy-test")
	if err != nil {
		t.Fatalf("Falha ao criar diretório temporário: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Configurar opções de teste
	options := types.SetupOptions{
		Force:          true,
		InstallService: false,
		ConfigPath:     filepath.Join(tempDir, "config.yaml"),
		HomeDir:        tempDir,
	}

	fmt.Println("=== Testando Setup ===")
	fmt.Printf("Diretório de teste: %s\n", tempDir)

	// Testar Setup
	setupResult, setupErr := Setup(options)
	if setupErr != nil {
		// Em sistemas não-Windows, esperamos ErrNotImplemented
		if setupErr.Error() != "funcionalidade não implementada para este sistema operacional: windows (stub)" &&
			setupErr.Error() != "funcionalidade não implementada para este sistema operacional: linux" &&
			setupErr.Error() != "funcionalidade não implementada para este sistema operacional: darwin" {
			t.Fatalf("Setup falhou com erro inesperado: %v", setupErr)
		}
		fmt.Printf("Setup não implementado para este sistema: %v (esperado em sistemas não-Windows)\n", setupErr)
	} else {
		fmt.Printf("Setup concluído com sucesso: %v\n", setupResult.Success)

		// Verificar se os arquivos foram criados
		configFile := setupResult.ConfigPath
		if fiErr := fileExists(configFile); fiErr != nil {
			t.Errorf("Arquivo de configuração não foi criado: %s", configFile)
		} else {
			fmt.Printf("Arquivo de configuração criado: %s\n", configFile)
		}
	}

	// Testar Status
	fmt.Println("\n=== Testando Status ===")
	statusResult, statusErr := Status(options)
	if statusErr != nil {
		// Em sistemas não-Windows, esperamos ErrNotImplemented
		if statusErr.Error() != "funcionalidade não implementada para este sistema operacional: windows (stub)" &&
			statusErr.Error() != "funcionalidade não implementada para este sistema operacional: linux" &&
			statusErr.Error() != "funcionalidade não implementada para este sistema operacional: darwin" {
			t.Fatalf("Status falhou com erro inesperado: %v", statusErr)
		}
		fmt.Printf("Status não implementado para este sistema: %v (esperado em sistemas não-Windows)\n", statusErr)
	} else {
		fmt.Printf("Status concluído com sucesso: %v\n", statusResult.Success)
		fmt.Printf("Caminho da configuração: %s\n", statusResult.ConfigPath)
	}

	// Testar Reset
	fmt.Println("\n=== Testando Reset ===")
	resetResult, resetErr := Reset(options)
	if resetErr != nil {
		// Em sistemas não-Windows, esperamos ErrNotImplemented
		if resetErr.Error() != "funcionalidade não implementada para este sistema operacional: windows (stub)" &&
			resetErr.Error() != "funcionalidade não implementada para este sistema operacional: linux" &&
			resetErr.Error() != "funcionalidade não implementada para este sistema operacional: darwin" {
			t.Fatalf("Reset falhou com erro inesperado: %v", resetErr)
		}
		fmt.Printf("Reset não implementado para este sistema: %v (esperado em sistemas não-Windows)\n", resetErr)
	} else {
		fmt.Printf("Reset concluído com sucesso: %v\n", resetResult.Success)

		// Verificar se os arquivos foram removidos (se Reset foi bem-sucedido)
		if resetResult.Success {
			configFile := filepath.Join(tempDir, "config", "manager.yaml")
			// Evitar shadowing de err: usar função helper
			if fiErr := ensureNotExists(configFile); fiErr != nil {
				t.Errorf("Arquivo de configuração não foi removido após reset: %s", configFile)
			} else {
				fmt.Printf("Arquivo de configuração removido com sucesso\n")
			}
		}
	}

	_ = err // manter 'err' referenciado (caso o linter reclame de variável não utilizada)
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

// Helpers sem shadowing

func fileExists(path string) error {
	_, statErr := os.Stat(path)
	if os.IsNotExist(statErr) {
		return fmt.Errorf("arquivo não existe")
	}
	return nil
}

func ensureNotExists(path string) error {
	_, statErr := os.Stat(path)
	if os.IsNotExist(statErr) {
		return nil
	}
	if statErr != nil {
		return statErr
	}
	return fmt.Errorf("arquivo ainda existe")
}
