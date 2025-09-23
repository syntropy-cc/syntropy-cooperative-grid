package usb

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

// NewUSBCommand cria o comando USB
func NewUSBCommand() *cobra.Command {
	usbCmd := &cobra.Command{
		Use:   "usb",
		Short: "Comandos para gerenciamento de dispositivos USB",
		Long: `Comandos para gerenciamento de dispositivos USB

Este grupo de comandos permite detectar, formatar e criar USBs com boot
para nós da Syntropy Cooperative Grid.

NOTA: No WSL, os comandos utilizam PowerShell para acessar dispositivos físicos.
      Para melhor desempenho, execute diretamente no Windows.
`,
	}

	usbCmd.AddCommand(newUSBListCommand())
	usbCmd.AddCommand(newUSBCreateCommand())
	usbCmd.AddCommand(newUSBFormatCommand())
	usbCmd.AddCommand(newUSBDebugCommand())

	return usbCmd
}

// newUSBListCommand cria o comando para listar dispositivos USB
func newUSBListCommand() *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lista dispositivos USB disponíveis",
		Long: `Lista todos os dispositivos USB disponíveis no sistema.

Exemplos:
  # Listar dispositivos em formato tabela
  syntropy usb list

  # Listar dispositivos em formato JSON
  syntropy usb list --format json

  # Listar dispositivos em formato YAML
  syntropy usb list --format yaml
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listUSBDevices(format)
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "table", "Formato de saída (table, json, yaml)")

	return cmd
}

// newUSBCreateCommand cria o comando para criar USB com boot
func newUSBCreateCommand() *cobra.Command {
	var (
		nodeName        string
		nodeDescription string
		coordinates     string
		ownerKeyFile    string
		autoDetect      bool
		label           string
		workDir         string
		cacheDir        string
		isoPath         string
		discoveryServer string
		createdBy       string
	)

	cmd := &cobra.Command{
		Use:   "create [device]",
		Short: "Cria USB com boot para um nó Syntropy",
		Long: `Cria um USB com boot contendo Ubuntu Server e configuração automática
para um nó da Syntropy Cooperative Grid.

Exemplos:
  # Criar USB com auto-detecção
  syntropy usb create --auto-detect --node-name "node-01"

  # Criar USB especificando dispositivo (Linux)
  syntropy usb create /dev/sdb --node-name "node-01"

  # Criar USB especificando dispositivo (Windows/WSL)
  syntropy usb create PHYSICALDRIVE1 --node-name "node-01"

  # Criar USB com ISO personalizada
  syntropy usb create --auto-detect --node-name "node-01" --iso /path/to/ubuntu.iso
`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var devicePath string

			if autoDetect {
				device, err := SelectDevice()
				if err != nil {
					return fmt.Errorf("falha na auto-detecção: %w", err)
				}
				devicePath = device.Path
			} else if len(args) > 0 {
				devicePath = args[0]
			} else {
				return fmt.Errorf("especifique um dispositivo ou use --auto-detect")
			}

			// Configurar usuário atual se não especificado
			if createdBy == "" {
				createdBy = os.Getenv("USER")
				if createdBy == "" {
					createdBy = "unknown"
				}
			}

			config := &Config{
				NodeName:        nodeName,
				NodeDescription: nodeDescription,
				Coordinates:     coordinates,
				OwnerKeyFile:    ownerKeyFile,
				Label:           label,
				ISOPath:         isoPath,
				DiscoveryServer: discoveryServer,
				CreatedBy:       createdBy,
			}

			return createUSB(devicePath, config, workDir, cacheDir)
		},
	}

	cmd.Flags().StringVar(&nodeName, "node-name", "", "Nome do nó (obrigatório)")
	cmd.Flags().StringVar(&nodeDescription, "description", "", "Descrição do nó")
	cmd.Flags().StringVar(&coordinates, "coordinates", "", "Coordenadas geográficas (lat,lon)")
	cmd.Flags().StringVar(&ownerKeyFile, "owner-key", "", "Arquivo de chave de proprietário existente")
	cmd.Flags().BoolVar(&autoDetect, "auto-detect", false, "Detectar automaticamente dispositivo USB")
	cmd.Flags().StringVar(&label, "label", "SYNTROPY", "Rótulo do sistema de arquivos")
	cmd.Flags().StringVar(&workDir, "work-dir", "", "Diretório de trabalho (padrão: ~/.syntropy/work)")
	cmd.Flags().StringVar(&cacheDir, "cache-dir", "", "Diretório de cache (padrão: ~/.syntropy/cache)")
	cmd.Flags().StringVar(&isoPath, "iso", "", "Caminho para ISO Ubuntu (baixa automaticamente se não especificado)")
	cmd.Flags().StringVar(&discoveryServer, "discovery-server", "syntropy-discovery.local", "Servidor de descoberta da rede")
	cmd.Flags().StringVar(&createdBy, "created-by", "", "Usuário que criou o nó (padrão: usuário atual)")

	cmd.MarkFlagRequired("node-name")

	return cmd
}

// newUSBFormatCommand cria o comando para formatar USB
func newUSBFormatCommand() *cobra.Command {
	var (
		label string
		force bool
	)

	cmd := &cobra.Command{
		Use:   "format [device]",
		Short: "Formata um dispositivo USB",
		Long: `Formata um dispositivo USB com sistema de arquivos FAT32.

⚠️  ATENÇÃO: Esta operação apagará TODOS os dados do dispositivo!

Exemplos:
  # Formatar USB (Linux)
  syntropy usb format /dev/sdb

  # Formatar USB (Windows/WSL)
  syntropy usb format PHYSICALDRIVE1

  # Formatar USB com rótulo personalizado
  syntropy usb format /dev/sdb --label "MYUSB"

  # Formatar USB sem confirmação
  syntropy usb format /dev/sdb --force
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return formatUSB(args[0], label, force)
		},
	}

	cmd.Flags().StringVar(&label, "label", "SYNTROPY", "Rótulo do sistema de arquivos")
	cmd.Flags().BoolVar(&force, "force", false, "Não pedir confirmação")

	return cmd
}

// newUSBDebugCommand cria o comando para debug do ambiente
func newUSBDebugCommand() *cobra.Command {
	var workDir string

	cmd := &cobra.Command{
		Use:   "debug",
		Short: "Executa diagnóstico do ambiente WSL/Windows",
		Long: `Executa diagnóstico completo do ambiente para identificar problemas
com a criação de USBs bootáveis.

Este comando verifica:
- Comandos necessários (dd, sgdisk, mkfs.vfat, etc.)
- Permissões sudo
- Dispositivos disponíveis
- Montagens ativas
- Espaço em disco
- Informações do sistema

Exemplos:
  # Executar diagnóstico completo
  syntropy usb debug

  # Executar diagnóstico em diretório específico
  syntropy usb debug --work-dir /tmp/debug
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Configurar diretório de trabalho
			if workDir == "" {
				homeDir, _ := os.UserHomeDir()
				workDir = filepath.Join(homeDir, ".syntropy", "work", "debug-"+time.Now().Format("20060102-150405"))
			}

			// Criar diretório se não existir
			os.MkdirAll(workDir, 0755)

			platform := detectPlatform()
			fmt.Printf("🔍 Executando diagnóstico para plataforma: %s\n", platform)
			fmt.Printf("📂 Diretório de trabalho: %s\n\n", workDir)

			if platform == "wsl" || platform == "windows" {
				return debugWSLEnvironment(workDir)
			} else {
				return fmt.Errorf("diagnóstico não implementado para plataforma: %s", platform)
			}
		},
	}

	cmd.Flags().StringVar(&workDir, "work-dir", "", "Diretório de trabalho para arquivos temporários")

	return cmd
}

// createUSB função principal para criar USB com boot
func createUSB(devicePath string, config *Config, workDir, cacheDir string) error {
	platform := detectPlatform()

	// Validar dispositivo
	if err := validateDevice(devicePath); err != nil {
		return err
	}

	// Configurar diretórios padrão com timestamp único
	if workDir == "" {
		homeDir, _ := os.UserHomeDir()
		workDir = filepath.Join(homeDir, ".syntropy", "work", "usb-"+time.Now().Format("20060102-150405"))
	}
	if cacheDir == "" {
		homeDir, _ := os.UserHomeDir()
		cacheDir = filepath.Join(homeDir, ".syntropy", "cache")
	}

	fmt.Printf("🚀 Iniciando criação de USB para nó: %s\n", config.NodeName)
	fmt.Printf("📍 Plataforma: %s\n", platform)
	fmt.Printf("💾 Dispositivo: %s\n", devicePath)
	fmt.Printf("📂 Diretório de trabalho: %s\n", workDir)
	fmt.Printf("📂 Diretório de cache: %s\n", cacheDir)
	fmt.Println()

	// Criar diretórios necessários
	os.MkdirAll(workDir, 0755)
	os.MkdirAll(cacheDir, 0755)

	// Gerar ou carregar chaves SSH
	if config.SSHPublicKey == "" {
		fmt.Println("🔑 Verificando chaves SSH existentes...")

		// Tentar carregar chaves existentes primeiro
		privateKey, publicKey, err := loadExistingSSHKeyPair(config.NodeName)
		if err != nil {
			// Se não existir, gerar novas chaves
			fmt.Println("🔑 Gerando novo par de chaves SSH...")
			privateKey, publicKey, err = generateSSHKeyPair(config.NodeName)
			if err != nil {
				return fmt.Errorf("erro ao gerar chaves SSH: %w", err)
			}
			fmt.Printf("✅ Chaves SSH geradas com sucesso\n")
		} else {
			fmt.Printf("✅ Chaves SSH existentes carregadas para nó: %s\n", config.NodeName)
		}

		// Configurar chaves
		config.SSHPrivateKey = privateKey
		config.SSHPublicKey = publicKey

		// Mostrar localização das chaves
		homeDir, _ := os.UserHomeDir()
		keyDir := filepath.Join(homeDir, ".syntropy", "keys")
		fmt.Printf("🔐 Chaves em: %s\n", keyDir)
		fmt.Printf("📁 Arquivos: %s-node.key, %s-node.key.pub, %s-node.fingerprint\n",
			config.NodeName, config.NodeName, config.NodeName)
		fmt.Println("⚠️  IMPORTANTE: A chave privada NÃO será enviada para o nó por segurança")
	}

	// Gerar certificados TLS
	fmt.Println("🔐 Gerando certificados TLS...")
	certs, err := generateCertificates(config.NodeName, config.OwnerKeyFile)
	if err != nil {
		return fmt.Errorf("erro ao gerar certificados: %w", err)
	}
	fmt.Println("✅ Certificados TLS gerados com sucesso")

	// Salvar certificados
	caKeyPath, caCertPath, nodeKeyPath, nodeCertPath, err := saveCertificates(certs, workDir)
	if err != nil {
		return fmt.Errorf("erro ao salvar certificados: %w", err)
	}

	// Gerar configuração do cloud-init
	fmt.Println("📝 Gerando configuração do cloud-init...")
	cloudInitConfig, err := generateCloudInitConfig(config, workDir, certs)
	if err != nil {
		return fmt.Errorf("erro ao gerar configuração do cloud-init: %w", err)
	}

	// Criar arquivos do cloud-init
	certPaths := map[string]string{
		"CAKey":    caKeyPath,
		"CACert":   caCertPath,
		"NodeKey":  nodeKeyPath,
		"NodeCert": nodeCertPath,
	}

	if err := createCloudInitFiles(cloudInitConfig, workDir, certPaths); err != nil {
		return fmt.Errorf("erro ao criar arquivos do cloud-init: %w", err)
	}
	fmt.Println("✅ Configuração do cloud-init criada com sucesso")

	// Copiar scripts para o diretório de trabalho
	fmt.Println("📜 Copiando scripts de instalação...")
	if err := copyScripts(workDir); err != nil {
		return fmt.Errorf("erro ao copiar scripts: %w", err)
	}
	fmt.Println("✅ Scripts copiados com sucesso")

	// Criar USB com estratégia NoCloud (dd + partição CIDATA)
	switch platform {
	case "wsl":
		return createUSBWithNoCloudWSL(devicePath, config, workDir, cacheDir)
	case "windows":
		return createUSBWithNoCloudWindows(devicePath, config, workDir, cacheDir)
	default:
		return createUSBWithNoCloudLinux(devicePath, config, workDir, cacheDir)
	}
}
