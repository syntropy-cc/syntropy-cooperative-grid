package usb

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewWindowsOnlyUSBCommand cria comandos específicos para Windows
func NewWindowsOnlyUSBCommand() *cobra.Command {
	winCmd := &cobra.Command{
		Use:   "usb-win",
		Short: "Comandos USB específicos para Windows",
		Long: `Comandos USB específicos para Windows com validações robustas
e tratamento de erros otimizado para criação de nós Syntropy.

Este módulo inclui:
- Validação completa de permissões e pré-requisitos
- Tratamento de erros específicos do Windows
- Verificação de UAC e privilégios de administrador
- Diagnóstico automático do ambiente
- Criação otimizada de USBs bootáveis

REQUISITOS:
- Windows 10/11 com WSL instalado
- Privilégios de Administrador
- PowerShell com política de execução adequada
`,
	}

	winCmd.AddCommand(newWindowsOnlyListCommand())
	winCmd.AddCommand(newWindowsOnlyCreateCommand())
	winCmd.AddCommand(newWindowsOnlyFormatCommand())
	winCmd.AddCommand(newWindowsOnlyDebugCommand())

	return winCmd
}

// newWindowsOnlyListCommand cria comando para listar dispositivos no Windows
func newWindowsOnlyListCommand() *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lista dispositivos USB (Windows Only)",
		Long: `Lista dispositivos USB com validações específicas do Windows.

Este comando executa verificações de ambiente antes de listar
os dispositivos, garantindo que tudo esteja configurado corretamente.

Exemplos:
  # Listar dispositivos em formato tabela
  syntropy usb-win list

  # Listar dispositivos em formato JSON
  syntropy usb-win list --format json

  # Listar dispositivos em formato YAML
  syntropy usb-win list --format yaml
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listWindowsOnlyDevicesFormatted(format)
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "table", "Formato de saída (table, json, yaml)")

	return cmd
}

// newWindowsOnlyCreateCommand cria comando para criar USB no Windows
func newWindowsOnlyCreateCommand() *cobra.Command {
	var (
		nodeName        string
		nodeDescription string
		coordinates     string
		ownerKeyFile    string
		label           string
		isoPath         string
		discoveryServer string
		createdBy       string
		tempDir         string
		logLevel        string
	)

	cmd := &cobra.Command{
		Use:   "create [device]",
		Short: "Cria USB bootável para nó Syntropy (Windows Only)",
		Long: `Cria um USB bootável otimizado para Windows com validações robustas.

Este comando executa verificações completas de ambiente e permissões
antes de criar o USB, garantindo o melhor resultado possível.

VALIDAÇÕES INCLUÍDAS:
- Privilégios de Administrador
- WSL disponível e configurado
- Política de execução do PowerShell
- Ferramentas necessárias instaladas
- Dispositivo USB válido e seguro
- Espaço em disco suficiente

Exemplos:
  # Criar USB com auto-detecção
  syntropy usb-win create --node-name "node-01"

  # Criar USB especificando dispositivo
  syntropy usb-win create PHYSICALDRIVE1 --node-name "node-01"

  # Criar USB com configurações completas
  syntropy usb-win create PHYSICALDRIVE1 \
    --node-name "node-01" \
    --description "Nó principal" \
    --coordinates "-23.5505,-46.6333" \
    --label "SYNTROPY-NODE"
`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var devicePath string

			if len(args) > 0 {
				devicePath = args[0]
			} else {
				// Auto-detectar dispositivo
				devices, err := listWindowsOnlyDevices()
				if err != nil {
					return fmt.Errorf("falha na detecção de dispositivos: %w", err)
				}

				if len(devices) == 0 {
					return fmt.Errorf("nenhum dispositivo USB encontrado")
				}

				if len(devices) == 1 {
					devicePath = fmt.Sprintf("PHYSICALDRIVE%d", devices[0].DiskNumber)
					fmt.Printf("Auto-selecionado: %s (%s)\n", devicePath, devices[0].FriendlyName)
				} else {
					return fmt.Errorf("múltiplos dispositivos encontrados. Especifique um: %s", devicePath)
				}
			}

			// Configurar usuário atual se não especificado
			if createdBy == "" {
				createdBy = os.Getenv("USERNAME")
				if createdBy == "" {
					createdBy = "unknown"
				}
			}

			config := &WindowsOnlyConfig{
				NodeName:        nodeName,
				NodeDescription: nodeDescription,
				Coordinates:     coordinates,
				OwnerKeyFile:    ownerKeyFile,
				Label:           label,
				ISOPath:         isoPath,
				DiscoveryServer: discoveryServer,
				CreatedBy:       createdBy,
				TempDir:         tempDir,
				LogLevel:        logLevel,
			}

			return createWindowsOnlyUSB(devicePath, config)
		},
	}

	cmd.Flags().StringVar(&nodeName, "node-name", "", "Nome do nó (obrigatório)")
	cmd.Flags().StringVar(&nodeDescription, "description", "", "Descrição do nó")
	cmd.Flags().StringVar(&coordinates, "coordinates", "", "Coordenadas geográficas (lat,lon)")
	cmd.Flags().StringVar(&ownerKeyFile, "owner-key", "", "Arquivo de chave de proprietário existente")
	cmd.Flags().StringVar(&label, "label", "SYNTROPY", "Rótulo do sistema de arquivos")
	cmd.Flags().StringVar(&isoPath, "iso", "", "Caminho para ISO Ubuntu (baixa automaticamente se não especificado)")
	cmd.Flags().StringVar(&discoveryServer, "discovery-server", "syntropy-discovery.local", "Servidor de descoberta da rede")
	cmd.Flags().StringVar(&createdBy, "created-by", "", "Usuário que criou o nó (padrão: usuário atual)")
	cmd.Flags().StringVar(&tempDir, "temp-dir", "", "Diretório temporário (padrão: %TEMP%\\syntropy-usb)")
	cmd.Flags().StringVar(&logLevel, "log-level", "info", "Nível de log (debug, info, warn, error)")

	cmd.MarkFlagRequired("node-name")

	return cmd
}

// newWindowsOnlyFormatCommand cria comando para formatar USB no Windows
func newWindowsOnlyFormatCommand() *cobra.Command {
	var (
		label string
		force bool
	)

	cmd := &cobra.Command{
		Use:   "format [device]",
		Short: "Formata dispositivo USB (Windows Only)",
		Long: `Formata um dispositivo USB com validações específicas do Windows.

⚠️  ATENÇÃO: Esta operação apagará TODOS os dados do dispositivo!

VALIDAÇÕES INCLUÍDAS:
- Privilégios de Administrador
- Dispositivo não é do sistema
- Dispositivo não contém partições críticas
- Confirmação do usuário (a menos que --force seja usado)

Exemplos:
  # Formatar USB
  syntropy usb-win format PHYSICALDRIVE1

  # Formatar USB com rótulo personalizado
  syntropy usb-win format PHYSICALDRIVE1 --label "MYUSB"

  # Formatar USB sem confirmação
  syntropy usb-win format PHYSICALDRIVE1 --force
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !force {
				fmt.Printf("⚠️  ATENÇÃO: Esta operação apagará TODOS os dados em %s!\n", args[0])
				fmt.Print("Tem certeza que deseja continuar? (y/N): ")

				var response string
				fmt.Scanln(&response)
				if response != "y" && response != "Y" && response != "yes" {
					fmt.Println("Operação cancelada.")
					return nil
				}
			}

			return formatWindowsOnlyUSB(args[0], label)
		},
	}

	cmd.Flags().StringVar(&label, "label", "SYNTROPY", "Rótulo do sistema de arquivos")
	cmd.Flags().BoolVar(&force, "force", false, "Não pedir confirmação")

	return cmd
}

// newWindowsOnlyDebugCommand cria comando para debug no Windows
func newWindowsOnlyDebugCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "debug",
		Short: "Executa diagnóstico completo do ambiente Windows",
		Long: `Executa diagnóstico completo do ambiente Windows para identificar
problemas com a criação de USBs bootáveis.

Este comando verifica:
- Privilégios de Administrador
- WSL disponível e configurado
- Política de execução do PowerShell
- Ferramentas necessárias instaladas
- Dispositivos USB disponíveis
- Espaço em disco e permissões

Exemplos:
  # Executar diagnóstico completo
  syntropy usb-win debug
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return debugWindowsOnlyEnvironment()
		},
	}

	return cmd
}

// listWindowsOnlyDevicesFormatted lista dispositivos e formata a saída
func listWindowsOnlyDevicesFormatted(format string) error {
	devices, err := listWindowsOnlyDevices()
	if err != nil {
		return fmt.Errorf("falha ao listar dispositivos: %w", err)
	}

	if len(devices) == 0 {
		fmt.Println("❌ Nenhum dispositivo USB encontrado.")
		fmt.Println("\n💡 Dicas:")
		fmt.Println("  • Certifique-se de que o USB está conectado")
		fmt.Println("  • Execute o PowerShell como Administrador")
		fmt.Println("  • Execute 'syntropy usb-win debug' para diagnóstico")
		return nil
	}

	switch format {
	case "json":
		return outputWindowsOnlyJSON(devices)
	case "yaml":
		return outputWindowsOnlyYAML(devices)
	default:
		return outputWindowsOnlyTable(devices)
	}
}

// outputWindowsOnlyTable exibe dispositivos em formato tabela
func outputWindowsOnlyTable(devices []WindowsOnlyUSBDevice) error {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("💾 Dispositivos USB Detectados (Windows Only)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for i, device := range devices {
		fmt.Printf("[%d] %s\n", i+1, fmt.Sprintf("PHYSICALDRIVE%d", device.DiskNumber))
		fmt.Printf("     Nome: %s\n", device.FriendlyName)
		fmt.Printf("     Tamanho: %s\n", device.SizeFormatted)
		fmt.Printf("     Modelo: %s\n", device.Model)
		fmt.Printf("     Serial: %s\n", device.SerialNumber)
		fmt.Printf("     Status: %s\n", device.Status)
		fmt.Printf("     Sistema: %t | Boot: %t | Offline: %t\n",
			device.IsSystem, device.IsBoot, device.IsOffline)
		fmt.Printf("     Partições: %d\n", device.PartitionCount)

		if i < len(devices)-1 {
			fmt.Println("     ───────────────────────────────────────────")
		}
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	return nil
}

// outputWindowsOnlyJSON exibe dispositivos em formato JSON
func outputWindowsOnlyJSON(devices []WindowsOnlyUSBDevice) error {
	jsonData, err := json.MarshalIndent(devices, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao converter para JSON: %w", err)
	}

	fmt.Println(string(jsonData))
	return nil
}

// outputWindowsOnlyYAML exibe dispositivos em formato YAML
func outputWindowsOnlyYAML(devices []WindowsOnlyUSBDevice) error {
	fmt.Println("devices:")
	for _, device := range devices {
		fmt.Printf("- disk_number: %d\n", device.DiskNumber)
		fmt.Printf("  friendly_name: \"%s\"\n", device.FriendlyName)
		fmt.Printf("  size: %d\n", device.Size)
		fmt.Printf("  size_formatted: \"%s\"\n", device.SizeFormatted)
		fmt.Printf("  serial_number: \"%s\"\n", device.SerialNumber)
		fmt.Printf("  bus_type: \"%s\"\n", device.BusType)
		fmt.Printf("  model: \"%s\"\n", device.Model)
		fmt.Printf("  is_system: %t\n", device.IsSystem)
		fmt.Printf("  is_boot: %t\n", device.IsBoot)
		fmt.Printf("  is_offline: %t\n", device.IsOffline)
		fmt.Printf("  partition_count: %d\n", device.PartitionCount)
		fmt.Printf("  status: \"%s\"\n", device.Status)
	}
	return nil
}
