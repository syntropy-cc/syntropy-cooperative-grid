package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"syntropy-cc/cooperative-grid/interfaces/cli/internal/cli/usb"

	"github.com/spf13/cobra"
)

// NewRootCommand cria o comando raiz da CLI
func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "syntropy",
		Short: "Syntropy Cooperative Grid Management CLI",
		Long: `Syntropy Cooperative Grid Management CLI

Uma ferramenta completa de linha de comando para gerenciar a Syntropy Cooperative Grid,
incluindo setup do ambiente, gerenciamento de nós, criação de USB bootável, 
deploy de containers e operações de rede cooperativa.

Exemplos:
  # Setup inicial do ambiente de gerenciamento
  syntropy setup

  # Listar nós gerenciados
  syntropy manager list

  # Conectar a um nó específico
  syntropy manager connect node-01

  # Descobrir nós na rede
  syntropy manager discover

  # Criar USB com boot para um nó
  syntropy usb create /dev/sdb --node-name "node-01"

  # Health check de todos os nós
  syntropy manager health

  # Backup das configurações
  syntropy manager backup
`,
		Version: "1.0.0",
	}

	// Configurar diretórios padrão
	setupDirectories()

	// Adicionar comandos
	rootCmd.AddCommand(NewSetupCommand())
	rootCmd.AddCommand(NewManagerCommand())
	rootCmd.AddCommand(NewTemplatesCommand())
	rootCmd.AddCommand(usb.NewUSBCommand())
	rootCmd.AddCommand(NewNodeCommand())
	rootCmd.AddCommand(NewContainerCommand())
	rootCmd.AddCommand(NewNetworkCommand())
	rootCmd.AddCommand(NewCooperativeCommand())

	return rootCmd
}

// setupDirectories cria diretórios necessários
func setupDirectories() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao obter diretório home: %v\n", err)
		os.Exit(1)
	}

	syntropyDir := filepath.Join(homeDir, ".syntropy")
	dirs := []string{
		syntropyDir,
		filepath.Join(syntropyDir, "nodes"),
		filepath.Join(syntropyDir, "keys"),
		filepath.Join(syntropyDir, "config"),
		filepath.Join(syntropyDir, "cache"),
		filepath.Join(syntropyDir, "work"),
		filepath.Join(syntropyDir, "scripts"),
		filepath.Join(syntropyDir, "backups"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao criar diretório %s: %v\n", dir, err)
			os.Exit(1)
		}
	}
}
