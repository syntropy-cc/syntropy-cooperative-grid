package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	setup "setup-component/src"

	"github.com/spf13/cobra"
)

var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "syntropy",
	Short: "🌐 Syntropy Cooperative Grid - Gerenciador de Rede Cooperativa",
	Long: `Syntropy Cooperative Grid CLI Manager

O Syntropy CLI permite que você configure e gerencie sua participação na
rede cooperativa descentralizada Syntropy.

🚀 Começando:
   syntropy setup run          Configure sua Command Station
   syntropy setup status       Verifique o status da instalação

📚 Principais Comandos:
   setup     Configure e gerencie seu ambiente Syntropy
   node      Crie e gerencie nós da rede cooperativa

💡 Dica: Use 'syntropy [comando] --help' para mais informações sobre cada comando.

Para mais informações visite: https://syntropy.network`,
	Version: fmt.Sprintf("%s (build: %s, commit: %s, platform: %s/%s)",
		version, buildTime, gitCommit, runtime.GOOS, runtime.GOARCH),
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.syntropy/config/manager.yaml)")
	rootCmd.PersistentFlags().Bool("verbose", false, "verbose output")
	rootCmd.PersistentFlags().Bool("quiet", false, "quiet output (suppress non-error messages)")

	// Add commands
	addCommands()
}

// addCommands adds all CLI subcommands
func addCommands() {
	// Setup commands
	rootCmd.AddCommand(setupCmd)

	// Node commands
	rootCmd.AddCommand(nodeCmd)

	// Future component commands will be added here:
	// rootCmd.AddCommand(workloadCmd)
	// rootCmd.AddCommand(configCmd)
	// rootCmd.AddCommand(stateCmd)
}

func main() {
	Execute()
}

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "⚙️  Configurar e gerenciar a Command Station",
	Long: `Configure seu computador como uma Command Station para gerenciar
a rede Syntropy Cooperative Grid.

A Command Station é seu "quartel general" de onde você:
   • Cria e gerencia nós da rede
   • Monitora o status da rede cooperativa
   • Configura segurança e autenticação
   • Gerencia tokens e credenciais

🎯 Comandos Disponíveis:
   run        Executar configuração inicial
   status     Verificar status da configuração
   validate   Validar ambiente antes do setup
   reset      Remover configuração (com backup)
   token      Gerenciar Grid Token de autenticação
   key        Gerenciar Owner Keys (Ed25519)`,
}

// setupRunCmd represents the setup run command
var setupRunCmd = &cobra.Command{
	Use:   "run",
	Short: "🚀 Executar configuração da Command Station",
	Long: `Configura seu computador como Command Station da rede Syntropy.

Este comando irá:
   ✓ Validar requisitos do sistema
   ✓ Criar estrutura de diretórios (~/.syntropy/)
   ✓ Gerar chaves criptográficas Ed25519
   ✓ Criar Grid Token de autenticação
   ✓ Configurar ambiente completo

O Grid Token será gerado automaticamente e armazenado de forma
segura no keyring do sistema operacional.

⚠️  Se um setup já existir, você será notificado e poderá:
   • Manter o setup existente
   • Sobrescrever com backup automático (--force)
   • Resetar completamente (syntropy setup reset)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")

		fmt.Println("🚀 Iniciando configuração do Syntropy Manager...")
		fmt.Println()

		manager, err := setup.NewSetupManager()
		if err != nil {
			return fmt.Errorf("❌ Falha ao inicializar: %w", err)
		}

		setupOptions := &setup.SetupOptions{
			Force:          force,
			CustomSettings: make(map[string]string),
		}

		err = manager.SetupWithPublicOptions(setupOptions)
		if err != nil {
			// Verificar se é erro de setup existente
			if setupErr, ok := err.(*setup.SetupError); ok && setupErr.Code == "SETUP_022" {
				fmt.Println("ℹ️  Setup já existe!")
				fmt.Println()
				fmt.Printf("   Criado em: %v\n", setupErr.Context["created_at"])
				fmt.Printf("   Versão: %v\n", setupErr.Context["version"])
				fmt.Println()
				fmt.Println("📋 Opções disponíveis:")
				fmt.Println("   • Ver status:     syntropy setup status")
				fmt.Println("   • Sobrescrever:   syntropy setup run --force")
				fmt.Println("   • Resetar tudo:   syntropy setup reset")
				fmt.Println()
				fmt.Println("⚠️  Sobrescrever ou resetar criará backup automático em ~/.syntropy/backups/")
				return nil
			}
			return fmt.Errorf("❌ Setup falhou: %w", err)
		}

		fmt.Println()
		fmt.Println("✅ Setup concluído com sucesso!")
		fmt.Println()

		// Verificar token
		exists, _ := manager.GridTokenExists()
		if exists {
			fmt.Println("🔐 Grid Token gerado e armazenado com segurança")
			fmt.Printf("   📍 Localização: Keyring do Sistema (%s)\n", runtime.GOOS)
			fmt.Println("   💡 Use 'syntropy setup token show' para visualizar")
		}

		fmt.Println()
		fmt.Println("🎉 Próximos passos:")
		fmt.Println("   1. Verifique o status:  syntropy setup status")
		fmt.Println("   2. Crie um nó:          syntropy node create")
		fmt.Println()

		return nil
	},
}

// setupStatusCmd represents the setup status command
var setupStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "📊 Verificar status da configuração",
	Long: `Verifica o status atual da configuração do Syntropy Manager.

Mostra informações sobre:
   • Estado da configuração
   • Chaves criptográficas
   • Grid Token
   • Configurações de nós
   • Logs e diagnósticos`,
	RunE: func(cmd *cobra.Command, args []string) error {
		manager, err := setup.NewSetupManager()
		if err != nil {
			return fmt.Errorf("❌ Falha ao criar gerenciador: %w", err)
		}

		status, err := manager.GetStatus()
		if err != nil {
			return fmt.Errorf("❌ Falha ao obter status: %w", err)
		}

		fmt.Println("📊 Status do Syntropy Manager")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Printf("   Status:        %s\n", *status)
		fmt.Println()

		// Verificar componentes
		fmt.Println("🔧 Componentes:")

		// Verificar Grid Token
		tokenExists, _ := manager.GridTokenExists()
		if tokenExists {
			fmt.Println("   ✅ Grid Token: Disponível")
		} else {
			fmt.Println("   ❌ Grid Token: Não encontrado")
		}

		// Verificar Owner Keys
		keyInfo, err := manager.GetOwnerKeyInfo()
		if err == nil {
			fmt.Printf("   ✅ Owner Keys: %s (%s)\n", keyInfo.Algorithm, keyInfo.Fingerprint[:16]+"...")
		} else {
			fmt.Println("   ❌ Owner Keys: Não encontradas")
		}

		fmt.Println()
		return nil
	},
}

// setupResetCmd represents the setup reset command
var setupResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "🗑️  Remover configuração do setup (com backup)",
	Long: `Remove completamente a configuração do Syntropy Manager.

⚠️  AVISO IMPORTANTE:
   Esta operação irá remover:
   • Todas as configurações em ~/.syntropy/
   • Chaves criptográficas
   • Grid Token
   • Configurações de nós (se houver)
   • Logs e estados salvos

🔒 SEGURANÇA:
   • Um backup automático será criado em ~/.syntropy/backups/
   • Você pode restaurar manualmente se necessário
   • O backup contém dados sensíveis - gerencie com cuidado!

💡 Recomendações de segurança:
   1. Exporte seu Grid Token antes: syntropy setup token export
   2. Faça backup manual das chaves se necessário
   3. Exclua backups antigos regularmente
   4. Nunca compartilhe backups não criptografados`,
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")

		if !force {
			fmt.Println("⚠️  AVISO: Esta operação irá remover TODAS as configurações do Syntropy!")
			fmt.Println()
			fmt.Println("📦 Um backup automático será criado em ~/.syntropy/backups/")
			fmt.Println()
			fmt.Print("   Tem certeza que deseja continuar? Digite 'yes' para confirmar: ")
			var response string
			fmt.Scanln(&response)
			if response != "yes" {
				fmt.Println("❌ Operação cancelada")
				return nil
			}
		}

		manager, err := setup.NewSetupManager()
		if err != nil {
			return fmt.Errorf("❌ Falha ao inicializar: %w", err)
		}

		fmt.Println("🗑️  Removendo configuração do Syntropy...")
		fmt.Println("📦 Criando backup automático...")

		err = manager.Reset(true)
		if err != nil {
			return fmt.Errorf("❌ Reset falhou: %w", err)
		}

		homeDir, _ := os.UserHomeDir()
		backupPath := fmt.Sprintf("%s/.syntropy/backups", homeDir)

		fmt.Println()
		fmt.Println("✅ Configuração removida com sucesso!")
		fmt.Println()
		fmt.Printf("📦 Backup salvo em: %s\n", backupPath)
		fmt.Println()
		fmt.Println("🔒 LEMBRETE DE SEGURANÇA:")
		fmt.Println("   • Backups contêm chaves e tokens sensíveis")
		fmt.Println("   • Gerencie os backups com cuidado")
		fmt.Println("   • Considere criptografar backups importantes")
		fmt.Println("   • Remova backups antigos regularmente")
		fmt.Println()
		fmt.Println("💡 Para reconfigurar: syntropy setup run")
		fmt.Println()

		return nil
	},
}

// setupValidateCmd represents the setup validate command
var setupValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "🔍 Validar ambiente antes do setup",
	Long: `Valida o ambiente do sistema antes de executar o setup.

Verifica:
   • Sistema operacional compatível
   • Permissões necessárias
   • Dependências do sistema
   • Espaço em disco disponível
   • Conectividade de rede

Este comando não faz alterações no sistema.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, _ := cmd.Flags().GetString("config-path")

		manager, err := setup.NewSetupManager()
		if err != nil {
			return fmt.Errorf("❌ Falha ao criar gerenciador: %w", err)
		}

		options := &setup.SetupOptions{
			ConfigPath:     configPath,
			CustomSettings: make(map[string]string),
		}

		fmt.Println("🔍 Validando ambiente do sistema...")
		fmt.Println()

		envInfo, err := manager.ValidateEnvironmentWithOptions(options)
		if err != nil {
			return fmt.Errorf("❌ Validação falhou: %w", err)
		}

		fmt.Println("✅ Validação concluída!")
		fmt.Println()
		fmt.Printf("   Sistema:       %s %s\n", envInfo.OS, envInfo.Architecture)
		fmt.Printf("   Pode prosseguir: %t\n", envInfo.CanProceed)
		fmt.Println()

		if !envInfo.CanProceed {
			fmt.Println("⚠️  Problemas encontrados:")
			for _, issue := range envInfo.Issues {
				fmt.Printf("   • %s\n", issue)
			}
			fmt.Println()
			fmt.Println("💡 Corrija os problemas antes de executar o setup")
		} else {
			fmt.Println("🎉 Sistema pronto para o setup!")
			fmt.Println("   Execute: syntropy setup run")
		}

		return nil
	},
}

// setupTokenCmd represents comandos de gerenciamento de token
var setupTokenCmd = &cobra.Command{
	Use:   "token",
	Short: "🔐 Gerenciar Grid Token do setup",
	Long: `Gerenciar o Grid Token usado para autenticação na rede Syntropy.

O Grid Token é gerado automaticamente durante o setup e armazenado
de forma segura no keyring do sistema operacional.`,
}

// setupTokenShowCmd represents the token show command
var setupTokenShowCmd = &cobra.Command{
	Use:   "show",
	Short: "👁️  Exibir o Grid Token atual",
	Long: `Exibe o Grid Token atual com preview seguro.

Use --full para exibir o token completo (requer --confirm por segurança).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		full, _ := cmd.Flags().GetBool("full")
		confirm, _ := cmd.Flags().GetBool("confirm")

		manager, err := setup.NewSetupManager()
		if err != nil {
			return fmt.Errorf("❌ Falha ao criar gerenciador: %w", err)
		}

		token, err := manager.GetGridToken()
		if err != nil {
			return fmt.Errorf("❌ Falha ao obter token: %w", err)
		}

		if full && confirm {
			fmt.Printf("🔑 Grid Token (COMPLETO):\n%s\n", token)
			fmt.Println("\n⚠️  AVISO: Mantenha este token em segurança!")
		} else if full && !confirm {
			fmt.Println("❌ Para exibir o token completo, use: --full --confirm")
		} else {
			fmt.Printf("🔑 Grid Token Preview: %s...[OCULTO]\n", token[:8])
			fmt.Println("💡 Use --full --confirm para ver o token completo")
		}

		return nil
	},
}

// setupTokenGenerateCmd represents the token generate command
var setupTokenGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "🔄 Gerar novo Grid Token",
	Long: `Gera um novo Grid Token substituindo o existente.

⚠️  AVISO: Isto irá invalidar o token anterior!`,
	RunE: func(cmd *cobra.Command, args []string) error {
		show, _ := cmd.Flags().GetBool("show")

		// Confirmação
		fmt.Print("⚠️  Gerar novo token irá invalidar o anterior. Continuar? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("❌ Operação cancelada")
			return nil
		}

		manager, err := setup.NewSetupManager()
		if err != nil {
			return fmt.Errorf("❌ Falha ao criar gerenciador: %w", err)
		}

		fmt.Println("🔄 Gerando novo Grid Token...")
		token, err := manager.GenerateGridToken()
		if err != nil {
			return fmt.Errorf("❌ Falha ao gerar token: %w", err)
		}

		fmt.Println("✅ Novo Grid Token gerado com sucesso!")
		if show {
			fmt.Printf("🔑 Token: %s\n", token)
		} else {
			fmt.Printf("🔑 Preview: %s...[OCULTO]\n", token[:8])
		}

		return nil
	},
}

// setupTokenRotateCmd represents the token rotate command
var setupTokenRotateCmd = &cobra.Command{
	Use:   "rotate",
	Short: "🔄 Rotacionar Grid Token",
	Long: `Rotaciona o Grid Token atual para maior segurança.

⚠️  AVISO: Isto irá invalidar o token anterior!`,
	RunE: func(cmd *cobra.Command, args []string) error {
		show, _ := cmd.Flags().GetBool("show")

		// Confirmação
		fmt.Print("⚠️  Rotacionar token irá invalidar o anterior. Continuar? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("❌ Operação cancelada")
			return nil
		}

		manager, err := setup.NewSetupManager()
		if err != nil {
			return fmt.Errorf("❌ Falha ao criar gerenciador: %w", err)
		}

		fmt.Println("🔄 Rotacionando Grid Token...")
		token, err := manager.RotateGridToken()
		if err != nil {
			return fmt.Errorf("❌ Falha ao rotacionar token: %w", err)
		}

		fmt.Println("✅ Grid Token rotacionado com sucesso!")
		if show {
			fmt.Printf("🔑 Novo Token: %s\n", token)
		} else {
			fmt.Printf("🔑 Preview: %s...[OCULTO]\n", token[:8])
		}

		return nil
	},
}

// setupTokenExportCmd represents the token export command
var setupTokenExportCmd = &cobra.Command{
	Use:   "export <arquivo-saida>",
	Short: "💾 Exportar Grid Token para backup",
	Long: `Exporta o Grid Token para arquivo de backup criptografado.

FORMATO DO ARQUIVO:
   • JSON estruturado com metadados
   • Token criptografado com AES-256
   • Inclui timestamp e versão
   • Verificação de integridade

EXEMPLO:
   syntropy setup token export backup-token.json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		outputPath := args[0]

		manager, err := setup.NewSetupManager()
		if err != nil {
			return fmt.Errorf("❌ Falha ao criar gerenciador: %w", err)
		}

		fmt.Println("💾 Exportando Grid Token...")

		err = manager.ExportGridToken(outputPath)
		if err != nil {
			return fmt.Errorf("❌ Falha ao exportar: %w", err)
		}

		fmt.Println()
		fmt.Println("✅ Grid Token exportado com sucesso!")
		fmt.Printf("   📁 Arquivo: %s\n", outputPath)
		fmt.Println()
		fmt.Println("🔒 BACKUP CONTÉM TOKEN CRIPTOGRAFADO")
		fmt.Println()
		fmt.Println("⚠️  INSTRUÇÕES DE SEGURANÇA:")
		fmt.Println("   1. Guarde este arquivo em local MUITO SEGURO")
		fmt.Println("   2. Considere armazenar cópias em locais diferentes")
		fmt.Println("   3. NUNCA envie por email ou mensagem")
		fmt.Println("   4. Use criptografia adicional se armazenar na nuvem")
		fmt.Println()
		fmt.Printf("💡 Para restaurar: syntropy setup token import %s\n", outputPath)

		return nil
	},
}

// setupTokenImportCmd represents the token import command
var setupTokenImportCmd = &cobra.Command{
	Use:   "import <arquivo-backup>",
	Short: "📥 Importar Grid Token de backup",
	Long: `Importa Grid Token de arquivo de backup.

⚠️  ATENÇÃO:
   • Esta operação irá SOBRESCREVER o token atual
   • Certifique-se de ter backup do token existente
   • O arquivo deve estar no formato correto

PROCESSO:
   1. Cria backup automático do token atual
   2. Valida o arquivo de backup
   3. Descriptografa e instala o novo token
   4. Valida a instalação`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputPath := args[0]

		fmt.Println("📥 Importando Grid Token")
		fmt.Println()
		fmt.Println("⚠️  AVISO: Isto irá sobrescrever o token atual!")
		fmt.Println()
		fmt.Print("   Digite 'IMPORT TOKEN' para confirmar: ")
		var response string
		fmt.Scanln(&response)
		if response != "IMPORT TOKEN" {
			fmt.Println("❌ Operação cancelada")
			return nil
		}

		manager, err := setup.NewSetupManager()
		if err != nil {
			return fmt.Errorf("❌ Falha ao criar gerenciador: %w", err)
		}

		fmt.Println()
		fmt.Println("📥 Importando token...")

		err = manager.ImportGridToken(inputPath)
		if err != nil {
			return fmt.Errorf("❌ Falha ao importar: %w", err)
		}

		fmt.Println()
		fmt.Println("✅ Grid Token importado com sucesso!")
		fmt.Println()
		fmt.Println("💡 Verifique o token: syntropy setup token show")

		return nil
	},
}

// setupTokenDeleteCmd represents the token delete command
var setupTokenDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "🗑️  Deletar Grid Token",
	Long: `Deleta o Grid Token atual do sistema.

⚠️  AVISO IMPORTANTE:
   • Esta operação é IRREVERSÍVEL
   • Você perderá acesso ao sistema
   • Faça backup antes de deletar
   • Confirmação dupla é obrigatória`,
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")

		if !force {
			fmt.Println("🚨 ATENÇÃO: Você está prestes a DELETAR seu Grid Token!")
			fmt.Println()
			fmt.Println("⚠️  CONSEQUÊNCIAS:")
			fmt.Println("   • Perda total de acesso ao sistema")
			fmt.Println("   • Necessidade de reconfiguração completa")
			fmt.Println("   • Operação IRREVERSÍVEL")
			fmt.Println()
			fmt.Print("   Digite 'DELETE TOKEN' para confirmar: ")
			var response string
			fmt.Scanln(&response)
			if response != "DELETE TOKEN" {
				fmt.Println("❌ Operação cancelada")
				return nil
			}
		}

		manager, err := setup.NewSetupManager()
		if err != nil {
			return fmt.Errorf("❌ Falha ao criar gerenciador: %w", err)
		}

		fmt.Println()
		fmt.Println("🗑️  Deletando Grid Token...")

		err = manager.DeleteGridToken()
		if err != nil {
			return fmt.Errorf("❌ Falha ao deletar: %w", err)
		}

		fmt.Println()
		fmt.Println("✅ Grid Token deletado com sucesso!")
		fmt.Println()
		fmt.Println("⚠️  PRÓXIMOS PASSOS:")
		fmt.Println("   1. Execute: syntropy setup run")
		fmt.Println("   2. Ou restaure de backup: syntropy setup token import")
		fmt.Println()

		return nil
	},
}

// setupKeyCmd represents comandos de gerenciamento de Owner Keys
var setupKeyCmd = &cobra.Command{
	Use:   "key",
	Short: "🔑 Gerenciar Owner Keys (Ed25519)",
	Long: `Gerenciar as Owner Keys criptográficas do Command Station.

As Owner Keys são chaves Ed25519 usadas para identificação criptográfica
do seu Command Station na rede Syntropy. Elas são PERMANENTES e críticas
para a segurança da sua infraestrutura.

⚠️  IMPORTANTE:
   • Owner Keys NÃO podem ser rotacionadas
   • A perda das chaves resulta em perda de acesso ao Command Station
   • Faça backup regular das chaves
   • Mantenha o backup em local seguro

🔐 Comandos Disponíveis:
   info       Ver informações das Owner Keys
   show       Exibir chaves (pública livre, privada com confirmação)
   export     Exportar chaves para backup seguro
   import     Importar chaves de backup`,
}

// setupKeyInfoCmd represents the key info command
var setupKeyInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "📋 Exibir informações das Owner Keys",
	Long: `Exibe informações gerais das Owner Keys sem expor dados sensíveis.

Mostra:
   • Algoritmo criptográfico (Ed25519)
   • Fingerprint (hash SHA256)
   • Data de criação
   • Localização dos arquivos`,
	RunE: func(cmd *cobra.Command, args []string) error {
		manager, err := setup.NewSetupManager()
		if err != nil {
			return fmt.Errorf("❌ Falha ao criar gerenciador: %w", err)
		}

		keyInfo, err := manager.GetOwnerKeyInfo()
		if err != nil {
			return fmt.Errorf("❌ Falha ao obter informações: %w", err)
		}

		fmt.Println("🔑 Owner Keys - Informações")
		fmt.Println()
		fmt.Printf("   Algoritmo:    %s\n", keyInfo.Algorithm)
		fmt.Printf("   Fingerprint:  %s\n", keyInfo.Fingerprint)
		fmt.Printf("   Criada em:    %s\n", keyInfo.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("   Localização:  %s\n", keyInfo.Path)
		fmt.Println()
		fmt.Println("💡 Use 'syntropy setup key show' para ver as chaves")
		fmt.Println("💡 Use 'syntropy setup key export' para fazer backup")

		return nil
	},
}

// setupKeyShowCmd represents the key show command
var setupKeyShowCmd = &cobra.Command{
	Use:   "show",
	Short: "👁️  Exibir Owner Keys",
	Long: `Exibe as Owner Keys do Command Station.

CHAVE PÚBLICA:
   • Pode ser compartilhada livremente
   • Usada por outros para verificar sua identidade
   • Sempre visível sem restrições

CHAVE PRIVADA:
   • EXTREMAMENTE SENSÍVEL - nunca compartilhe!
   • Requer confirmação dupla (--force --confirm)
   • Acesso é registrado em log de auditoria
   • Use apenas quando absolutamente necessário`,
	RunE: func(cmd *cobra.Command, args []string) error {
		showPrivate, _ := cmd.Flags().GetBool("force")
		confirm, _ := cmd.Flags().GetBool("confirm")

		manager, err := setup.NewSetupManager()
		if err != nil {
			return fmt.Errorf("❌ Falha ao criar gerenciador: %w", err)
		}

		// Sempre mostrar chave pública
		publicKey, err := manager.GetOwnerPublicKey()
		if err != nil {
			return fmt.Errorf("❌ Falha ao obter chave pública: %w", err)
		}

		fmt.Println("🔓 Owner Public Key")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println(publicKey)
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println()
		fmt.Println("✅ Esta chave pode ser compartilhada livremente")
		fmt.Println()

		// Mostrar preview da chave privada ou chave completa
		if showPrivate {
			if !confirm {
				// Mostrar preview da chave privada
				privateKey, err := manager.GetOwnerPrivateKey("default_passphrase")
				if err != nil {
					return fmt.Errorf("❌ Falha ao obter chave privada: %w", err)
				}

				fmt.Println("🔐 Owner Private Key (Preview)")
				fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
				fmt.Printf("%s...[OCULTO]\n", privateKey[:8])
				fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
				fmt.Println()
				fmt.Println("🚨 AVISO DE SEGURANÇA:")
				fmt.Println("   • A chave privada é EXTREMAMENTE SENSÍVEL")
				fmt.Println("   • NUNCA compartilhe sua chave privada")
				fmt.Println("   • Use --force --confirm para ver a chave completa")
				return nil
			}

			// Dupla confirmação
			fmt.Println("🚨 ATENÇÃO: Você está prestes a visualizar sua CHAVE PRIVADA!")
			fmt.Print("   Digite 'SHOW PRIVATE KEY' para confirmar: ")

			// Usar bufio.Reader para melhor compatibilidade com Windows/WSL
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("❌ Erro ao ler entrada:", err)
				return nil
			}

			// Remover quebras de linha e espaços
			response = strings.TrimSpace(response)

			// Verificar se a resposta está correta (case-insensitive para melhor UX)
			expected := "SHOW PRIVATE KEY"
			if strings.ToUpper(response) != strings.ToUpper(expected) {
				fmt.Println("❌ Operação cancelada - confirmação incorreta")
				fmt.Printf("   Esperado: '%s'\n", expected)
				fmt.Printf("   Recebido: '%s'\n", response)
				fmt.Println("   💡 Dica: Digite exatamente: SHOW PRIVATE KEY")
				return nil
			}

			fmt.Println("✅ Confirmação aceita!")

			privateKey, err := manager.GetOwnerPrivateKey("default_passphrase")
			if err != nil {
				return fmt.Errorf("❌ Falha ao obter chave privada: %w", err)
			}

			fmt.Println()
			fmt.Println("🔐 Owner Private Key")
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			fmt.Println(privateKey)
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			fmt.Println()
			fmt.Println("🚨 MANTENHA ESTA CHAVE EM SEGURANÇA ABSOLUTA!")
			fmt.Println("   • NUNCA compartilhe esta chave")
			fmt.Println("   • Apague esta tela após copiar")
			fmt.Println("   • Considere fazer backup: syntropy setup key export")
		}

		return nil
	},
}

// setupKeyExportCmd represents the key export command
var setupKeyExportCmd = &cobra.Command{
	Use:   "export <arquivo-saida>",
	Short: "💾 Exportar Owner Keys para backup",
	Long: `Exporta Owner Keys para arquivo de backup criptografado.

EXPORTAÇÃO:
   • Chave pública: sempre incluída (não criptografada)
   • Chave privada: opcional, criptografada com senha

FLAGS:
   --include-private    Incluir chave privada no backup
   --password          Senha para criptografar chave privada

FORMATO DO ARQUIVO:
   • JSON estruturado com metadados
   • Chave pública em texto claro
   • Chave privada criptografada com AES-256
   • Inclui timestamp e versão

EXEMPLO:
   # Exportar apenas chave pública
   syntropy setup key export backup-public.json

   # Exportar com chave privada (recomendado)
   syntropy setup key export backup-full.json --include-private

   💡 Você será solicitado a criar uma senha forte para o backup`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		outputPath := args[0]
		includePrivate, _ := cmd.Flags().GetBool("include-private")

		manager, err := setup.NewSetupManager()
		if err != nil {
			return fmt.Errorf("❌ Falha ao criar gerenciador: %w", err)
		}

		var passphrase string
		if includePrivate {
			fmt.Println("🔐 Exportando Owner Keys com chave privada")
			fmt.Println()
			fmt.Println("📝 Você precisará criar uma SENHA FORTE para criptografar o backup.")
			fmt.Println("   Esta senha será necessária para restaurar as chaves.")
			fmt.Println()
			fmt.Print("   Digite a senha do backup: ")
			fmt.Scanln(&passphrase)

			if len(passphrase) < 12 {
				return fmt.Errorf("❌ Senha muito curta (mínimo 12 caracteres)")
			}

			fmt.Print("   Confirme a senha: ")
			var confirm string
			fmt.Scanln(&confirm)

			if passphrase != confirm {
				return fmt.Errorf("❌ Senhas não correspondem")
			}
		} else {
			fmt.Println("🔓 Exportando apenas chave pública (sem criptografia)")
		}

		fmt.Println()
		fmt.Println("💾 Criando backup...")

		err = manager.ExportOwnerKeys(outputPath, includePrivate, passphrase)
		if err != nil {
			return fmt.Errorf("❌ Falha ao exportar: %w", err)
		}

		fmt.Println()
		fmt.Println("✅ Owner Keys exportadas com sucesso!")
		fmt.Printf("   📁 Arquivo: %s\n", outputPath)
		fmt.Println()

		if includePrivate {
			fmt.Println("🔒 BACKUP CONTÉM CHAVE PRIVADA CRIPTOGRAFADA")
			fmt.Println()
			fmt.Println("⚠️  INSTRUÇÕES DE SEGURANÇA:")
			fmt.Println("   1. Guarde este arquivo em local MUITO SEGURO")
			fmt.Println("   2. Anote a senha em local separado")
			fmt.Println("   3. Considere armazenar cópias em locais diferentes")
			fmt.Println("   4. NUNCA envie por email ou mensagem")
			fmt.Println("   5. Use criptografia adicional se armazenar na nuvem")
			fmt.Println()
			fmt.Printf("💡 Para restaurar: syntropy setup key import %s\n", outputPath)
		} else {
			fmt.Println("ℹ️  Este backup contém apenas a chave PÚBLICA")
			fmt.Println("   Para backup completo, use: --include-private")
		}

		return nil
	},
}

// setupKeyImportCmd represents the key import command
var setupKeyImportCmd = &cobra.Command{
	Use:   "import <arquivo-backup>",
	Short: "📥 Importar Owner Keys de backup",
	Long: `Importa Owner Keys de arquivo de backup.

⚠️  ATENÇÃO:
   • Esta operação irá SOBRESCREVER as chaves atuais
   • Certifique-se de ter backup das chaves existentes
   • A senha usada na exportação será necessária

PROCESSO:
   1. Cria backup automático das chaves atuais
   2. Valida o arquivo de backup
   3. Descriptografa chave privada (se incluída)
   4. Instala as novas chaves
   5. Valida a instalação`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputPath := args[0]

		fmt.Println("📥 Importando Owner Keys")
		fmt.Println()
		fmt.Println("⚠️  AVISO: Isto irá sobrescrever as chaves atuais!")
		fmt.Println()
		fmt.Print("   Digite 'IMPORT KEYS' para confirmar: ")
		var response string
		fmt.Scanln(&response)
		if response != "IMPORT KEYS" {
			fmt.Println("❌ Operação cancelada")
			return nil
		}

		manager, err := setup.NewSetupManager()
		if err != nil {
			return fmt.Errorf("❌ Falha ao criar gerenciador: %w", err)
		}

		// Criar backup das chaves atuais
		fmt.Println()
		fmt.Println("📦 Criando backup das chaves atuais...")
		backupPath := fmt.Sprintf("keys-backup-%d.json", time.Now().Unix())
		if err := manager.ExportOwnerKeys(backupPath, true, "auto-backup"); err != nil {
			fmt.Printf("⚠️  Aviso: Não foi possível fazer backup: %v\n", err)
		} else {
			fmt.Printf("✅ Backup criado: %s\n", backupPath)
		}

		// Solicitar senha se backup incluir chave privada
		fmt.Println()
		fmt.Print("   Digite a senha do backup: ")
		var passphrase string
		fmt.Scanln(&passphrase)

		fmt.Println()
		fmt.Println("📥 Importando chaves...")

		err = manager.ImportOwnerKeys(inputPath, passphrase)
		if err != nil {
			return fmt.Errorf("❌ Falha ao importar: %w", err)
		}

		fmt.Println()
		fmt.Println("✅ Owner Keys importadas com sucesso!")
		fmt.Println()
		fmt.Println("💡 Verifique as chaves: syntropy setup key info")

		return nil
	},
}

// nodeCmd represents the node command
var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "🖥️  Gerenciar nós da rede cooperativa",
	Long: `Gerencie nós da rede cooperativa Syntropy.

Este comando permite:
   • Criar e provisionar novos nós automaticamente
   • Listar e monitorar nós existentes
   • Visualizar status e logs dos nós
   • Remover nós da rede
   • Iniciar/parar listeners dos nós

O componente node implementa um sistema plug-and-play que cria dispositivos
USB bootáveis e registra nós automaticamente sem intervenção do usuário.`,
}

func init() {
	// Setup subcommands
	setupCmd.AddCommand(setupRunCmd)
	setupCmd.AddCommand(setupStatusCmd)
	setupCmd.AddCommand(setupValidateCmd)
	setupCmd.AddCommand(setupResetCmd)
	setupCmd.AddCommand(setupTokenCmd)
	setupCmd.AddCommand(setupKeyCmd)

	// Token subcommands
	setupTokenCmd.AddCommand(setupTokenShowCmd)
	setupTokenCmd.AddCommand(setupTokenGenerateCmd)
	setupTokenCmd.AddCommand(setupTokenRotateCmd)
	setupTokenCmd.AddCommand(setupTokenExportCmd)
	setupTokenCmd.AddCommand(setupTokenImportCmd)
	setupTokenCmd.AddCommand(setupTokenDeleteCmd)

	// Key subcommands
	setupKeyCmd.AddCommand(setupKeyInfoCmd)
	setupKeyCmd.AddCommand(setupKeyShowCmd)
	setupKeyCmd.AddCommand(setupKeyExportCmd)
	setupKeyCmd.AddCommand(setupKeyImportCmd)

	// Setup flags
	setupRunCmd.Flags().Bool("force", false, "forçar sobrescrever setup existente")
	setupResetCmd.Flags().Bool("force", false, "pular confirmação de reset")

	// Token flags
	setupTokenShowCmd.Flags().Bool("full", false, "exibir token completo (requer --confirm)")
	setupTokenShowCmd.Flags().Bool("confirm", false, "confirmar exibição de token completo")
	setupTokenGenerateCmd.Flags().Bool("show", false, "exibir token após gerar")
	setupTokenRotateCmd.Flags().Bool("show", false, "exibir token após rotacionar")
	setupTokenDeleteCmd.Flags().Bool("force", false, "pular confirmação de deleção")

	// Key flags
	setupKeyShowCmd.Flags().Bool("force", false, "exibir chave privada (requer --confirm)")
	setupKeyShowCmd.Flags().Bool("confirm", false, "confirmar exibição de chave privada")
	setupKeyExportCmd.Flags().Bool("include-private", false, "incluir chave privada no backup")

	// Setup validate flags
	setupValidateCmd.Flags().String("config-path", "", "caminho personalizado do arquivo de configuração")

	// Add placeholder node subcommands
	nodeCmd.AddCommand(&cobra.Command{
		Use:   "create",
		Short: "Criar novo nó",
		Long:  "Cria um novo nó na rede cooperativa",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("🖥️  Node component is not yet fully integrated.")
			fmt.Println("This is a placeholder command for build testing.")
			return nil
		},
	})

	nodeCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "Listar nós existentes",
		Long:  "Lista todos os nós registrados na rede",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("🖥️  Node component is not yet fully integrated.")
			fmt.Println("This is a placeholder command for build testing.")
			return nil
		},
	})
}
