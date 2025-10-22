package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	node "manager/interfaces/cli/node/src"
)

func main() {
	fmt.Println("🔍 Testando correção do problema de token...")

	// Criar um logger simples
	logger := &testLogger{}

	// Criar o setup adapter
	setupAdapter := node.NewSetupAdapter(logger)

	// Criar a integração de token
	tokenIntegration := setupAdapter.CreateTokenIntegration()

	// Tentar inicializar a integração de token
	fmt.Println("📋 Inicializando integração de token...")
	err := setupAdapter.InitializeTokenIntegration(tokenIntegration)
	if err != nil {
		fmt.Printf("❌ Erro na inicialização: %v\n", err)

		// Verificar se o arquivo de token existe
		homeDir, _ := os.UserHomeDir()
		tokenFile := filepath.Join(homeDir, ".syntropy", "tokens", "grid-token.json")
		if _, err := os.Stat(tokenFile); err == nil {
			fmt.Printf("✅ Arquivo de token existe em: %s\n", tokenFile)
			fmt.Println("🔧 O problema está na lógica de verificação de existência do token")
		} else {
			fmt.Printf("❌ Arquivo de token não encontrado em: %s\n", tokenFile)
		}
		return
	}

	fmt.Println("✅ Integração de token inicializada com sucesso!")

	// Tentar obter o token
	fmt.Println("🔑 Obtendo token...")
	token, err := tokenIntegration.GetGridToken()
	if err != nil {
		fmt.Printf("❌ Erro ao obter token: %v\n", err)
		return
	}

	fmt.Printf("✅ Token obtido com sucesso: %s...[OCULTO]\n", token[:8])
}

// Logger simples para teste
type testLogger struct{}

func (l *testLogger) Debug(msg string, args ...interface{}) {
	log.Printf("[DEBUG] %s", msg)
}

func (l *testLogger) Info(msg string, args ...interface{}) {
	log.Printf("[INFO] %s", msg)
}

func (l *testLogger) Warn(msg string, args ...interface{}) {
	log.Printf("[WARN] %s", msg)
}

func (l *testLogger) Error(msg string, args ...interface{}) {
	log.Printf("[ERROR] %s", msg)
}
