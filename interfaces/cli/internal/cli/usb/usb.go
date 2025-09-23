package usb

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

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

			config := &Config{
				NodeName:        nodeName,
				NodeDescription: nodeDescription,
				Coordinates:     coordinates,
				OwnerKeyFile:    ownerKeyFile,
				Label:           label,
				ISOPath:         isoPath,
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
	cmd.Flags().StringVar(&workDir, "work-dir", "", "Diretório de trabalho (padrão: /tmp/syntropy-work)")
	cmd.Flags().StringVar(&cacheDir, "cache-dir", "", "Diretório de cache (padrão: ~/.syntropy/cache)")
	cmd.Flags().StringVar(&isoPath, "iso", "", "Caminho para ISO Ubuntu (baixa automaticamente se não especificado)")

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

// Estruturas
type USBDevice struct {
	Path        string `json:"path"`
	Size        string `json:"size"`
	SizeGB      int    `json:"size_gb"`
	Model       string `json:"model"`
	Vendor      string `json:"vendor"`
	Serial      string `json:"serial"`
	Removable   bool   `json:"removable"`
	Platform    string `json:"platform"`
	DiskNumber  int    `json:"disk_number,omitempty"`
	WindowsPath string `json:"windows_path,omitempty"`
}

type Config struct {
	NodeName        string `json:"node_name"`
	NodeDescription string `json:"node_description"`
	Coordinates     string `json:"coordinates"`
	OwnerKeyFile    string `json:"owner_key_file"`
	Label           string `json:"label"`
	ISOPath         string `json:"iso_path"`
}

// WindowsDisk estrutura para parse do JSON do PowerShell
type WindowsDisk struct {
	Number       int    `json:"Number"`
	FriendlyName string `json:"FriendlyName"`
	Size         int64  `json:"Size"`
	SerialNumber string `json:"SerialNumber"`
	BusType      string `json:"BusType"`
	Model        string `json:"Model"`
}

// Funções principais

func listUSBDevices(format string) error {
	platform := detectPlatform()

	// Mostrar aviso para WSL
	if platform == "wsl" {
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println("🖥️  WSL Detectado - Acessando dispositivos via Windows")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println()
	}

	devices, err := ListDevices()
	if err != nil {
		return fmt.Errorf("falha ao detectar dispositivos: %w", err)
	}

	if len(devices) == 0 {
		fmt.Println("❌ Nenhum dispositivo USB encontrado.")
		if platform == "wsl" {
			fmt.Println("\n💡 Dicas para WSL:")
			fmt.Println("  • Certifique-se de que o USB está conectado")
			fmt.Println("  • Execute o PowerShell como Administrador")
			fmt.Println("  • Tente executar: powershell.exe Get-Disk")
		}
		return nil
	}

	switch format {
	case "json":
		return outputJSON(devices)
	case "yaml":
		return outputYAML(devices)
	default:
		return outputTable(devices)
	}
}

func createUSB(devicePath string, config *Config, workDir, cacheDir string) error {
	platform := detectPlatform()

	// Configurar diretórios padrão
	if workDir == "" {
		if platform == "windows" || platform == "wsl" {
			workDir = os.TempDir()
		} else {
			workDir = "/tmp/syntropy-work"
		}
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

	switch platform {
	case "wsl":
		return createUSBWSL(devicePath, config, workDir, cacheDir)
	case "windows":
		return createUSBWindows(devicePath, config, workDir, cacheDir)
	default:
		return createUSBLinux(devicePath, config, workDir, cacheDir)
	}
}

func formatUSB(devicePath, label string, force bool) error {
	platform := detectPlatform()

	// Confirmação do usuário
	if !force {
		fmt.Printf("⚠️  ATENÇÃO: Esta operação apagará TODOS os dados em %s!\n", devicePath)
		fmt.Print("Tem certeza que deseja continuar? (y/N): ")

		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" && response != "yes" {
			fmt.Println("Operação cancelada.")
			return nil
		}
	}

	fmt.Printf("🔧 Formatando dispositivo %s...\n", devicePath)

	switch platform {
	case "wsl":
		return formatUSBWSL(devicePath, label)
	case "windows":
		return formatUSBWindows(devicePath, label)
	default:
		return formatUSBLinux(devicePath, label)
	}
}

// Funções específicas por plataforma

func detectPlatform() string {
	if runtime.GOOS == "windows" {
		return "windows"
	}

	// Verificar se está rodando no WSL
	if _, err := os.Stat("/proc/sys/fs/binfmt_misc/WSLInterop"); err == nil {
		return "wsl"
	}

	// Verificar outra forma de detectar WSL
	if data, err := os.ReadFile("/proc/version"); err == nil {
		if strings.Contains(strings.ToLower(string(data)), "microsoft") {
			return "wsl"
		}
	}

	return "linux"
}

func ListDevices() ([]USBDevice, error) {
	platform := detectPlatform()

	switch platform {
	case "windows":
		return listDevicesWindows()
	case "wsl":
		return listDevicesWSL()
	case "linux":
		return listDevicesLinux()
	default:
		return nil, fmt.Errorf("plataforma não suportada: %s", platform)
	}
}

// WSL Functions

func listDevicesWSL() ([]USBDevice, error) {
	// PowerShell script para listar USBs físicos
	psScript := `
	Get-Disk | Where-Object {
		$_.BusType -eq 'USB' -or 
		($_.BusType -eq 'SCSI' -and $_.Size -lt 500GB -and $_.Size -gt 1GB)
	} | Select-Object Number, FriendlyName, Size, SerialNumber, BusType, Model | 
	ConvertTo-Json -Compress
	`

	cmd := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-Command", psScript)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Tentar método alternativo
		return listDevicesWSLAlternative()
	}

	// Limpar output do PowerShell
	jsonStr := strings.TrimSpace(string(output))
	if jsonStr == "" {
		return []USBDevice{}, nil
	}

	// Parse JSON
	var disks []WindowsDisk

	// Verificar se é array ou objeto único
	if strings.HasPrefix(jsonStr, "[") {
		if err := json.Unmarshal([]byte(jsonStr), &disks); err != nil {
			return nil, fmt.Errorf("erro ao fazer parse do JSON (array): %w", err)
		}
	} else {
		var disk WindowsDisk
		if err := json.Unmarshal([]byte(jsonStr), &disk); err != nil {
			return nil, fmt.Errorf("erro ao fazer parse do JSON (objeto): %w", err)
		}
		disks = []WindowsDisk{disk}
	}

	// Converter para USBDevice
	var devices []USBDevice
	for _, disk := range disks {
		device := USBDevice{
			Path:        fmt.Sprintf("PHYSICALDRIVE%d", disk.Number),
			WindowsPath: fmt.Sprintf("\\\\.\\PHYSICALDRIVE%d", disk.Number),
			DiskNumber:  disk.Number,
			Size:        formatSize(disk.Size),
			SizeGB:      int(disk.Size / (1024 * 1024 * 1024)),
			Model:       disk.Model,
			Vendor:      "Unknown",
			Serial:      disk.SerialNumber,
			Removable:   true,
			Platform:    "wsl",
		}

		if disk.FriendlyName != "" {
			device.Model = disk.FriendlyName
		}

		devices = append(devices, device)
	}

	return devices, nil
}

func listDevicesWSLAlternative() ([]USBDevice, error) {
	// Método alternativo usando WMIC
	cmd := exec.Command("cmd.exe", "/c", "wmic diskdrive where \"InterfaceType='USB'\" get Model,Size,SerialNumber,Index /format:csv")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("erro ao executar WMIC: %w", err)
	}

	var devices []USBDevice
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) >= 5 && parts[1] != "Index" && parts[1] != "" {
			index, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
			size, _ := strconv.ParseInt(strings.TrimSpace(parts[3]), 10, 64)

			devices = append(devices, USBDevice{
				Path:        fmt.Sprintf("PHYSICALDRIVE%d", index),
				WindowsPath: fmt.Sprintf("\\\\.\\PHYSICALDRIVE%d", index),
				DiskNumber:  index,
				Model:       strings.TrimSpace(parts[2]),
				Size:        formatSize(size),
				SizeGB:      int(size / (1024 * 1024 * 1024)),
				Serial:      strings.TrimSpace(parts[4]),
				Removable:   true,
				Platform:    "wsl",
			})
		}
	}

	return devices, nil
}

// createUSBWSL reescrito: grava a ISO bit-a-bit usando wsl --mount --bare + dd
func createUSBWSL(devicePath string, config *Config, workDir, cacheDir string) error {
	// 1) Extrai número do disco Windows (PHYSICALDRIVE<N>)
	var diskNum int
	switch {
	case strings.HasPrefix(devicePath, "PHYSICALDRIVE"):
		fmt.Sscanf(devicePath, "PHYSICALDRIVE%d", &diskNum)
	case strings.HasPrefix(devicePath, "\\\\.\\PHYSICALDRIVE"):
		fmt.Sscanf(devicePath, "\\\\.\\PHYSICALDRIVE%d", &diskNum)
	default:
		return fmt.Errorf("formato de dispositivo inválido para WSL: %s", devicePath)
	}
	winPhysical := fmt.Sprintf("\\\\.\\PHYSICALDRIVE%d", diskNum)

	// 2) Garante ISO no cache (ou usa a fornecida) e converte caminho p/ WSL
	isoPath := config.ISOPath
	if isoPath == "" {
		var err error
		isoPath, err = manageISOCache(cacheDir)
		if err != nil {
			return fmt.Errorf("erro ao gerenciar ISO: %w", err)
		}
	}
	isoWSL := convertAnyToWSLPath(isoPath) // aceita /mnt/c/... ou C:\...

	fmt.Printf("📀 ISO (WSL): %s\n", isoWSL)
	fmt.Printf("🧱 Disco: %s (nº %d)\n\n", winPhysical, diskNum)

	// 3) Script PowerShell elevando privilégios e orquestrando o fluxo
	psScript := fmt.Sprintf(`
$ErrorActionPreference = "Stop"

Write-Host "Colocando disco offline no Windows..." -ForegroundColor Cyan
Set-Disk -Number %d -IsReadOnly $false -IsOffline $true

try {
    Write-Host "Montando o disco cru no WSL (--bare)..." -ForegroundColor Cyan
    wsl --mount %s --bare

    Write-Host "Gravando ISO no dispositivo via WSL (dd)..." -ForegroundColor Cyan
    # Descobrir o device recém-montado com segurança: compara antes/depois
    wsl bash -lc 'set -euo pipefail
before=($(ls /dev/sd? 2>/dev/null || true))
sleep 0.5
# Confirma que o device apareceu (em algumas máquinas demora um pouco)
tries=0
while [ $tries -lt 20 ]; do
  after=($(ls /dev/sd? 2>/dev/null || true))
  # encontra item de after que não está em before
  dev=""
  for d in "${after[@]}"; do
    found=0
    for b in "${before[@]}"; do [ "$d" = "$b" ] && found=1 && break; done
    [ $found -eq 0 ] && dev="$d" && break
  done
  if [ -n "$dev" ]; then
    echo "Dispositivo WSL detectado: $dev"
    ISO="%s"
    sudo dd if="$ISO" of="$dev" bs=4M status=progress conv=fsync
    sync
    exit 0
  fi
  tries=$((tries+1))
  sleep 0.5
done
echo "Falha ao detectar o device no WSL." 1>&2
exit 1
'

} finally {
    Write-Host "Desmontando do WSL e voltando disco online no Windows..." -ForegroundColor Cyan
    try { wsl --unmount %s } catch { Write-Host "Aviso: unmount falhou ou já desmontado." -ForegroundColor Yellow }
    Set-Disk -Number %d -IsOffline $false
}

Write-Host "✅ USB criado (modo dd). Pronto para boot." -ForegroundColor Green
`, diskNum, winPhysical, isoWSL, winPhysical, diskNum)

	// 4) Grava e executa o script elevado
	os.MkdirAll(workDir, 0755)
	scriptPath := filepath.Join(workDir, "create_usb_dd.ps1")
	if err := os.WriteFile(scriptPath, []byte(psScript), 0644); err != nil {
		return fmt.Errorf("erro ao criar script temporário: %w", err)
	}
	winScriptPath := convertWSLToWindowsPath(scriptPath)

	fmt.Println("📝 Solicitando permissões de administrador...")
	cmd := exec.Command("powershell.exe", "-Command",
		fmt.Sprintf(`Start-Process powershell -ArgumentList '-NoProfile -ExecutionPolicy Bypass -File "%s"' -Verb RunAs -Wait`, winScriptPath))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("erro ao executar criação do USB (dd): %w", err)
	}

	return nil
}

// Helper novo: aceita caminho Windows (C:\...) ou já em WSL (/mnt/c/...)
// e devolve SEMPRE um caminho válido no WSL.
func convertAnyToWSLPath(p string) string {
	// Se já parece WSL (começa com /), mantém
	if strings.HasPrefix(p, "/") {
		return p
	}
	// Tenta converter via wslpath -u
	out, err := exec.Command("wslpath", "-u", p).Output()
	if err == nil && len(out) > 0 {
		return strings.TrimSpace(string(out))
	}
	// Fallback: tentativa simples C:\ -> /mnt/c/
	if len(p) >= 3 && p[1] == ':' {
		drive := strings.ToLower(string(p[0]))
		rest := strings.ReplaceAll(p[2:], `\`, `/`)
		return fmt.Sprintf("/mnt/%s/%s", drive, strings.TrimLeft(rest, `/`))
	}
	return p
}

func formatUSBWSL(devicePath, label string) error {
	// Extrair número do disco
	var diskNum int
	if strings.HasPrefix(devicePath, "PHYSICALDRIVE") {
		fmt.Sscanf(devicePath, "PHYSICALDRIVE%d", &diskNum)
	} else {
		return fmt.Errorf("formato de dispositivo inválido: %s", devicePath)
	}

	// Script PowerShell para formatar
	psScript := fmt.Sprintf(`
	$ErrorActionPreference = "Stop"
	Clear-Disk -Number %d -RemoveData -Confirm:$false
	New-Partition -DiskNumber %d -UseMaximumSize -AssignDriveLetter |
		Format-Volume -FileSystem FAT32 -NewFileSystemLabel "%s" -Confirm:$false
	Write-Host "Formatação concluída!" -ForegroundColor Green
	`, diskNum, diskNum, label)

	cmd := exec.Command("powershell.exe", "-Command",
		fmt.Sprintf(`Start-Process powershell -ArgumentList '-NoProfile -Command "%s"' -Verb RunAs -Wait`,
			strings.ReplaceAll(psScript, "\n", "; ")))

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("erro ao formatar: %w", err)
	}

	fmt.Printf("✅ Dispositivo %s formatado com sucesso!\n", devicePath)
	return nil
}

// Windows Functions

func listDevicesWindows() ([]USBDevice, error) {
	// Similar ao WSL mas sem necessidade de conversões
	psScript := `
	Get-Disk | Where-Object {$_.BusType -eq 'USB'} | 
	Select-Object Number, FriendlyName, Size, SerialNumber, BusType, Model | 
	ConvertTo-Json -Compress
	`

	cmd := exec.Command("powershell", "-NoProfile", "-Command", psScript)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("erro ao executar PowerShell: %w", err)
	}

	var disks []WindowsDisk
	if err := json.Unmarshal(output, &disks); err != nil {
		// Tentar parse de objeto único
		var disk WindowsDisk
		if err := json.Unmarshal(output, &disk); err != nil {
			return nil, fmt.Errorf("erro ao fazer parse: %w", err)
		}
		disks = []WindowsDisk{disk}
	}

	var devices []USBDevice
	for _, disk := range disks {
		devices = append(devices, USBDevice{
			Path:       fmt.Sprintf("\\\\.\\PHYSICALDRIVE%d", disk.Number),
			DiskNumber: disk.Number,
			Size:       formatSize(disk.Size),
			SizeGB:     int(disk.Size / (1024 * 1024 * 1024)),
			Model:      disk.FriendlyName,
			Serial:     disk.SerialNumber,
			Removable:  true,
			Platform:   "windows",
		})
	}

	return devices, nil
}

func createUSBWindows(devicePath string, config *Config, workDir, cacheDir string) error {
	// Implementação similar ao WSL mas sem conversão de caminhos
	return fmt.Errorf("implementação Windows nativa pendente")
}

func formatUSBWindows(devicePath, label string) error {
	// Implementação similar ao WSL mas sem conversão
	return fmt.Errorf("implementação Windows nativa pendente")
}

// Linux Functions

func listDevicesLinux() ([]USBDevice, error) {
	var devices []USBDevice

	// Listar dispositivos de bloco
	blockDevs, _ := filepath.Glob("/sys/block/sd*")
	for _, blockDev := range blockDevs {
		devName := filepath.Base(blockDev)

		// Verificar se é removível
		removableData, _ := os.ReadFile(filepath.Join(blockDev, "removable"))
		if strings.TrimSpace(string(removableData)) != "1" {
			continue
		}

		// Obter informações do dispositivo
		device := USBDevice{
			Path:      "/dev/" + devName,
			Removable: true,
			Platform:  "linux",
		}

		// Tamanho
		if sizeData, err := os.ReadFile(filepath.Join(blockDev, "size")); err == nil {
			if sectors, err := strconv.ParseInt(strings.TrimSpace(string(sizeData)), 10, 64); err == nil {
				sizeBytes := sectors * 512
				device.Size = formatSize(sizeBytes)
				device.SizeGB = int(sizeBytes / (1024 * 1024 * 1024))
			}
		}

		// Modelo e Vendor
		if modelData, err := os.ReadFile(filepath.Join(blockDev, "device/model")); err == nil {
			device.Model = strings.TrimSpace(string(modelData))
		}
		if vendorData, err := os.ReadFile(filepath.Join(blockDev, "device/vendor")); err == nil {
			device.Vendor = strings.TrimSpace(string(vendorData))
		}

		devices = append(devices, device)
	}

	return devices, nil
}

func createUSBLinux(devicePath string, config *Config, workDir, cacheDir string) error {
	// Verificar se tem permissões de root
	if os.Geteuid() != 0 {
		return fmt.Errorf("este comando requer privilégios de root (use sudo)")
	}

	// Gerenciar ISO com cache inteligente
	isoPath := config.ISOPath
	if isoPath == "" {
		var err error
		isoPath, err = manageISOCache(cacheDir)
		if err != nil {
			return fmt.Errorf("erro ao gerenciar ISO: %w", err)
		}
	}

	fmt.Println("📝 Criando USB bootável com dd...")

	// Usar dd para gravar ISO no USB
	cmd := exec.Command("dd",
		fmt.Sprintf("if=%s", isoPath),
		fmt.Sprintf("of=%s", devicePath),
		"bs=4M",
		"status=progress",
		"oflag=sync")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("erro ao gravar ISO: %w", err)
	}

	// Sync para garantir que tudo foi gravado
	exec.Command("sync").Run()

	fmt.Println("\n✅ USB criado com sucesso!")
	return nil
}

func formatUSBLinux(devicePath, label string) error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("este comando requer privilégios de root (use sudo)")
	}

	// Desmontar partições se estiverem montadas
	exec.Command("umount", devicePath+"*").Run()

	// Criar nova tabela de partições
	cmd := exec.Command("parted", "-s", devicePath, "mklabel", "msdos")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("erro ao criar tabela de partições: %w", err)
	}

	// Criar partição primária
	cmd = exec.Command("parted", "-s", devicePath, "mkpart", "primary", "fat32", "0%", "100%")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("erro ao criar partição: %w", err)
	}

	// Formatar partição
	partition := devicePath + "1"
	if strings.Contains(devicePath, "nvme") || strings.Contains(devicePath, "mmcblk") {
		partition = devicePath + "p1"
	}

	cmd = exec.Command("mkfs.vfat", "-F", "32", "-n", label, partition)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("erro ao formatar partição: %w", err)
	}

	fmt.Printf("✅ Dispositivo %s formatado com sucesso!\n", devicePath)
	return nil
}

// Funções auxiliares

// manageISOCache gerencia o download e cache de ISOs Ubuntu
func manageISOCache(cacheDir string) (string, error) {
	// Criar diretório de cache para ISOs
	isoDir := filepath.Join(cacheDir, "iso")
	if err := os.MkdirAll(isoDir, 0755); err != nil {
		return "", fmt.Errorf("erro ao criar diretório de cache: %w", err)
	}

	// Lista de URLs para tentar (do mais recente para o mais antigo)
	isoOptions := []struct {
		url      string
		filename string
		name     string
	}{
		{
			url:      "https://releases.ubuntu.com/24.04/ubuntu-24.04-live-server-amd64.iso",
			filename: "ubuntu-24.04-server.iso",
			name:     "Ubuntu 24.04 LTS Server",
		},
		{
			url:      "https://releases.ubuntu.com/22.04.5/ubuntu-22.04.4-live-server-amd64.iso",
			filename: "ubuntu-22.04.5-server-amd64.iso",
			name:     "Ubuntu 22.04.5 LTS Server",
		},
		{
			url:      "https://releases.ubuntu.com/22.04/ubuntu-22.04.4-live-server-amd64.iso",
			filename: "ubuntu-22.04-server.iso",
			name:     "Ubuntu 22.04 LTS Server",
		},
		{
			url:      "https://releases.ubuntu.com/20.04/ubuntu-20.04.6-live-server-amd64.iso",
			filename: "ubuntu-20.04-server.iso",
			name:     "Ubuntu 20.04.6 LTS Server",
		},
	}

	// Verificar se alguma ISO já existe no cache
	fmt.Println("🔍 Verificando cache de ISOs...")
	for _, iso := range isoOptions {
		isoPath := filepath.Join(isoDir, iso.filename)
		if fileInfo, err := os.Stat(isoPath); err == nil {
			// Verificar se o arquivo tem tamanho razoável (> 500MB)
			if fileInfo.Size() > 500*1024*1024 {
				fmt.Printf("✅ ISO encontrada no cache: %s\n", iso.name)
				fmt.Printf("   Arquivo: %s\n", isoPath)
				fmt.Printf("   Tamanho: %.2f GB\n", float64(fileInfo.Size())/(1024*1024*1024))
				return isoPath, nil
			} else {
				fmt.Printf("⚠️  ISO corrompida encontrada (tamanho: %d bytes), removendo...\n", fileInfo.Size())
				os.Remove(isoPath)
			}
		}
	}

	// Nenhuma ISO no cache, tentar baixar
	fmt.Println("\n📥 Nenhuma ISO encontrada no cache. Iniciando download...")
	fmt.Println("   Cache: " + isoDir)

	for _, iso := range isoOptions {
		isoPath := filepath.Join(isoDir, iso.filename)
		fmt.Printf("\n🌐 Tentando baixar: %s\n", iso.name)
		fmt.Printf("   URL: %s\n", iso.url)

		// Verificar se a URL existe antes de baixar (HEAD request)
		if !checkURLExists(iso.url) {
			fmt.Printf("   ❌ URL não disponível, tentando próxima opção...\n")
			continue
		}

		// Baixar com progresso
		if err := downloadWithProgress(iso.url, isoPath); err != nil {
			fmt.Printf("   ❌ Erro ao baixar: %v\n", err)
			// Remover arquivo parcial se existir
			os.Remove(isoPath)
			continue
		}

		// Verificar tamanho do arquivo baixado
		if fileInfo, err := os.Stat(isoPath); err == nil {
			if fileInfo.Size() > 500*1024*1024 {
				fmt.Printf("\n✅ ISO baixada com sucesso!\n")
				fmt.Printf("   Arquivo: %s\n", isoPath)
				fmt.Printf("   Tamanho: %.2f GB\n", float64(fileInfo.Size())/(1024*1024*1024))
				return isoPath, nil
			} else {
				fmt.Printf("   ❌ Arquivo baixado muito pequeno, tentando próxima opção...\n")
				os.Remove(isoPath)
			}
		}
	}

	return "", fmt.Errorf("não foi possível baixar nenhuma ISO Ubuntu. Verifique sua conexão ou forneça uma ISO com --iso")
}

// checkURLExists verifica se uma URL existe fazendo um HEAD request
func checkURLExists(url string) bool {
	cmd := exec.Command("curl", "-I", "--silent", "--head", "--fail", url)
	err := cmd.Run()
	return err == nil
}

// downloadWithProgress baixa um arquivo com indicador de progresso
func downloadWithProgress(url, destPath string) error {
	// Criar arquivo temporário
	tmpPath := destPath + ".tmp"

	// Usar wget com barra de progresso
	cmd := exec.Command("wget",
		"--progress=bar:force:noscroll",
		"--tries=3",
		"--timeout=30",
		"--continue",
		"-O", tmpPath,
		url)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	// Mover arquivo temporário para destino final
	if err := os.Rename(tmpPath, destPath); err != nil {
		return fmt.Errorf("erro ao mover arquivo: %w", err)
	}

	return nil
}

func SelectDevice() (*USBDevice, error) {
	devices, err := ListDevices()
	if err != nil {
		return nil, err
	}

	if len(devices) == 0 {
		return nil, fmt.Errorf("nenhum dispositivo USB encontrado")
	}

	// Se houver apenas um dispositivo, seleciona automaticamente
	if len(devices) == 1 {
		fmt.Printf("Auto-selecionado: %s (%s)\n", devices[0].Path, devices[0].Model)
		return &devices[0], nil
	}

	// Mostrar opções para o usuário
	fmt.Println("\n🔍 Dispositivos USB detectados:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	for i, device := range devices {
		fmt.Printf("[%d] %s - %s (%s)\n", i+1, device.Path, device.Model, device.Size)
	}

	fmt.Print("\nSelecione o dispositivo (1-", len(devices), "): ")
	var choice int
	fmt.Scanln(&choice)

	if choice < 1 || choice > len(devices) {
		return nil, fmt.Errorf("seleção inválida")
	}

	return &devices[choice-1], nil
}

func convertWSLToWindowsPath(wslPath string) string {
	// Converter caminho WSL para Windows
	cmd := exec.Command("wslpath", "-w", wslPath)
	output, err := cmd.Output()
	if err != nil {
		// Fallback: conversão manual
		if strings.HasPrefix(wslPath, "/mnt/") {
			parts := strings.Split(wslPath, "/")
			if len(parts) > 2 {
				drive := strings.ToUpper(parts[2])
				remainingPath := strings.Join(parts[3:], "\\")
				return fmt.Sprintf("%s:\\%s", drive, remainingPath)
			}
		}
		return wslPath
	}
	return strings.TrimSpace(string(output))
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	units := []string{"KB", "MB", "GB", "TB"}
	if exp >= len(units) {
		exp = len(units) - 1
	}

	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp])
}

// Funções de saída

func outputTable(devices []USBDevice) error {
	platform := detectPlatform()

	if platform == "wsl" || platform == "windows" {
		fmt.Printf("%-15s %-8s %-30s %-15s %-10s\n",
			"DISPOSITIVO", "TAMANHO", "MODELO", "SERIAL", "PLATAFORMA")
		fmt.Println(strings.Repeat("─", 80))

		for _, device := range devices {
			fmt.Printf("%-15s %-8s %-30s %-15s %-10s\n",
				device.Path, device.Size, device.Model, device.Serial, device.Platform)
		}
	} else {
		fmt.Printf("%-12s %-8s %-20s %-15s %-10s %-10s\n",
			"DISPOSITIVO", "TAMANHO", "MODELO", "FABRICANTE", "REMOVÍVEL", "PLATAFORMA")
		fmt.Println(strings.Repeat("─", 80))

		for _, device := range devices {
			removable := "Não"
			if device.Removable {
				removable = "Sim"
			}
			fmt.Printf("%-12s %-8s %-20s %-15s %-10s %-10s\n",
				device.Path, device.Size, device.Model, device.Vendor, removable, device.Platform)
		}
	}

	return nil
}

func outputJSON(devices []USBDevice) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(devices)
}

func outputYAML(devices []USBDevice) error {
	fmt.Println("devices:")
	for _, device := range devices {
		fmt.Printf("  - path: %s\n", device.Path)
		fmt.Printf("    size: %s\n", device.Size)
		fmt.Printf("    size_gb: %d\n", device.SizeGB)
		fmt.Printf("    model: %s\n", device.Model)
		fmt.Printf("    vendor: %s\n", device.Vendor)
		fmt.Printf("    serial: %s\n", device.Serial)
		fmt.Printf("    removable: %t\n", device.Removable)
		fmt.Printf("    platform: %s\n", device.Platform)
		if device.DiskNumber > 0 {
			fmt.Printf("    disk_number: %d\n", device.DiskNumber)
		}
		if device.WindowsPath != "" {
			fmt.Printf("    windows_path: %s\n", device.WindowsPath)
		}
	}
	return nil
}
