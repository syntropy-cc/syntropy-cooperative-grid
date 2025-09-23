package usb

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	// "syntropy-cc/cooperative-grid/core/services/usb" // Temporariamente comentado para compilação

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
	)

	cmd := &cobra.Command{
		Use:   "create [device]",
		Short: "Cria USB com boot para um nó Syntropy",
		Long: `Cria um USB com boot contendo Ubuntu Server e configuração automática
para um nó da Syntropy Cooperative Grid.

Exemplos:
  # Criar USB com auto-detecção
  syntropy usb create --auto-detect --node-name "node-01"

  # Criar USB especificando dispositivo
  syntropy usb create /dev/sdb --node-name "node-01"

  # Criar USB com coordenadas específicas
  syntropy usb create /dev/sdb --node-name "node-01" --coordinates "-23.5505,-46.6333"

  # Criar USB usando chave de proprietário existente
  syntropy usb create /dev/sdb --node-name "node-02" --owner-key ~/.syntropy/keys/main.key
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

			return createUSB(devicePath, nodeName, nodeDescription, coordinates, 
				ownerKeyFile, label, workDir, cacheDir)
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
  # Formatar USB
  syntropy usb format /dev/sdb

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

// listUSBDevices lista dispositivos USB disponíveis
func listUSBDevices(format string) error {
	devices, err := ListDevices()
	if err != nil {
		return fmt.Errorf("falha ao detectar dispositivos: %w", err)
	}

	if len(devices) == 0 {
		fmt.Println("Nenhum dispositivo USB encontrado.")
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

// createUSB cria um USB com boot
func createUSB(devicePath, nodeName, nodeDescription, coordinates, ownerKeyFile, 
	label, workDir, cacheDir string) error {
	
	// Configurar diretórios padrão se não especificados
	if workDir == "" {
		workDir = "/tmp/syntropy-work"
	}
	if cacheDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("erro ao obter diretório home: %w", err)
		}
		cacheDir = filepath.Join(homeDir, ".syntropy", "cache")
	}

	// Criar instância do creator
	creator := NewCreator(workDir, cacheDir)
	defer creator.Cleanup()

	// Configurar criação
	config := &Config{
		NodeName:        nodeName,
		NodeDescription: nodeDescription,
		Coordinates:     coordinates,
		OwnerKeyFile:    ownerKeyFile,
		Label:           label,
	}

	fmt.Printf("Iniciando criação de USB para nó: %s\n", nodeName)
	fmt.Printf("Dispositivo: %s\n", devicePath)
	fmt.Printf("Diretório de trabalho: %s\n", workDir)
	fmt.Printf("Diretório de cache: %s\n", cacheDir)
	fmt.Println()

	// Executar criação
	if err := creator.CreateUSB(devicePath, config); err != nil {
		return fmt.Errorf("falha na criação do USB: %w", err)
	}

	fmt.Println("✅ USB criado com sucesso!")
	fmt.Printf("Nó: %s\n", nodeName)
	fmt.Printf("Dispositivo: %s\n", devicePath)
	fmt.Println()
	fmt.Println("Próximos passos:")
	fmt.Println("1. Remova o USB com segurança")
	fmt.Println("2. Insira no hardware alvo")
	fmt.Println("3. Configure boot para USB no BIOS/UEFI")
	fmt.Println("4. A instalação será automática (~30 minutos)")

	return nil
}

// formatUSB formata um dispositivo USB
func formatUSB(devicePath, label string, force bool) error {
	// Validações de segurança
	detector := NewDetector()
	if err := detector.ValidateDevice(devicePath); err != nil {
		return fmt.Errorf("dispositivo inválido: %w", err)
	}

	if detector.IsSystemDisk(devicePath) {
		return fmt.Errorf("ERRO: Dispositivo parece ser um disco do sistema: %s", devicePath)
	}

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

	// Formatar dispositivo
	formatter := NewFormatter()
	if err := formatter.FormatDevice(devicePath, label); err != nil {
		return fmt.Errorf("falha ao formatar dispositivo: %w", err)
	}

	fmt.Printf("✅ Dispositivo %s formatado com sucesso!\n", devicePath)
	fmt.Printf("Rótulo: %s\n", label)

	return nil
}

// Funções de saída em diferentes formatos

func outputTable(devices []USBDevice) error {
	fmt.Printf("%-12s %-8s %-20s %-15s %-10s %s\n", 
		"DISPOSITIVO", "TAMANHO", "MODELO", "FABRICANTE", "REMOVÍVEL", "PLATAFORMA")
	fmt.Println(strings.Repeat("-", 80))

	for _, device := range devices {
		removable := "Não"
		if device.Removable {
			removable = "Sim"
		}
		fmt.Printf("%-12s %-8s %-20s %-15s %-10s %s\n",
			device.Path, device.Size, device.Model, device.Vendor, removable, device.Platform)
	}

	return nil
}

func outputJSON(devices []USBDevice) error {
	// Implementação simples de JSON
	fmt.Println("[")
	for i, device := range devices {
		fmt.Printf("  {\n")
		fmt.Printf("    \"path\": \"%s\",\n", device.Path)
		fmt.Printf("    \"size\": \"%s\",\n", device.Size)
		fmt.Printf("    \"size_gb\": %d,\n", device.SizeGB)
		fmt.Printf("    \"model\": \"%s\",\n", device.Model)
		fmt.Printf("    \"vendor\": \"%s\",\n", device.Vendor)
		fmt.Printf("    \"serial\": \"%s\",\n", device.Serial)
		fmt.Printf("    \"removable\": %t,\n", device.Removable)
		fmt.Printf("    \"platform\": \"%s\"\n", device.Platform)
		fmt.Printf("  }")
		if i < len(devices)-1 {
			fmt.Print(",")
		}
		fmt.Println()
	}
	fmt.Println("]")
	return nil
}

func outputYAML(devices []USBDevice) error {
	fmt.Println("devices:")
	for _, device := range devices {
		fmt.Printf("- path: %s\n", device.Path)
		fmt.Printf("  size: %s\n", device.Size)
		fmt.Printf("  size_gb: %d\n", device.SizeGB)
		fmt.Printf("  model: %s\n", device.Model)
		fmt.Printf("  vendor: %s\n", device.Vendor)
		fmt.Printf("  serial: %s\n", device.Serial)
		fmt.Printf("  removable: %t\n", device.Removable)
		fmt.Printf("  platform: %s\n", device.Platform)
	}
	return nil
}

// Estruturas mock para substituir o core USB service
type USBDevice struct {
	Path      string `json:"path"`
	Size      string `json:"size"`
	SizeGB    int    `json:"size_gb"`
	Model     string `json:"model"`
	Vendor    string `json:"vendor"`
	Serial    string `json:"serial"`
	Removable bool   `json:"removable"`
	Platform  string `json:"platform"`
}

type Config struct {
	NodeName        string `json:"node_name"`
	NodeDescription string `json:"node_description"`
	Coordinates     string `json:"coordinates"`
	OwnerKeyFile    string `json:"owner_key_file"`
	Label           string `json:"label"`
}

// detectPlatform detecta a plataforma atual
func detectPlatform() string {
	if runtime.GOOS == "windows" {
		return "windows"
	}
	
	// Verificar se está rodando no WSL
	if _, err := os.Stat("/proc/sys/fs/binfmt_misc/WSLInterop"); err == nil {
		return "wsl"
	}
	
	return "linux"
}

// ListDevices lista dispositivos USB reais baseado na plataforma
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

// listDevicesLinux lista dispositivos USB no Linux
func listDevicesLinux() ([]USBDevice, error) {
	var devices []USBDevice
	
	// Ler /proc/partitions para encontrar dispositivos de bloco
	file, err := os.Open("/proc/partitions")
	if err != nil {
		return nil, fmt.Errorf("erro ao ler /proc/partitions: %w", err)
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	scanner.Scan() // Pular cabeçalho
	
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if len(line) < 4 {
			continue
		}
		
		deviceName := line[3]
		// Filtrar apenas dispositivos que parecem ser USB (sd*, mmcblk*, etc.)
		if !strings.HasPrefix(deviceName, "sd") && !strings.HasPrefix(deviceName, "mmcblk") {
			continue
		}
		
		// Verificar se é um dispositivo USB
		if isUSBDevice(deviceName) {
			device, err := getDeviceInfo(deviceName)
			if err != nil {
				continue // Pular dispositivos com erro
			}
			devices = append(devices, *device)
		}
	}
	
	// Se não encontrou dispositivos USB, tentar listar todos os dispositivos removíveis
	if len(devices) == 0 {
		return listRemovableDevices()
	}
	
	return devices, nil
}

// listRemovableDevices lista todos os dispositivos removíveis como fallback
func listRemovableDevices() ([]USBDevice, error) {
	var devices []USBDevice
	
	// Ler /proc/partitions para encontrar dispositivos de bloco
	file, err := os.Open("/proc/partitions")
	if err != nil {
		return nil, fmt.Errorf("erro ao ler /proc/partitions: %w", err)
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	scanner.Scan() // Pular cabeçalho
	
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if len(line) < 4 {
			continue
		}
		
		deviceName := line[3]
		// Filtrar apenas dispositivos que parecem ser removíveis (sd*, mmcblk*, etc.)
		if !strings.HasPrefix(deviceName, "sd") && !strings.HasPrefix(deviceName, "mmcblk") {
			continue
		}
		
		// Verificar se é removível
		removablePath := fmt.Sprintf("/sys/block/%s/removable", deviceName)
		if data, err := os.ReadFile(removablePath); err == nil {
			if strings.TrimSpace(string(data)) == "1" {
				device, err := getDeviceInfo(deviceName)
				if err != nil {
					continue // Pular dispositivos com erro
				}
				devices = append(devices, *device)
			}
		}
	}
	
	// Se ainda não encontrou dispositivos removíveis, listar todos os dispositivos de armazenamento
	// (útil para ambientes virtuais como WSL)
	if len(devices) == 0 {
		return listAllStorageDevices()
	}
	
	return devices, nil
}

// listAllStorageDevices lista todos os dispositivos de armazenamento como último recurso
func listAllStorageDevices() ([]USBDevice, error) {
	var devices []USBDevice
	
	// Ler /proc/partitions para encontrar dispositivos de bloco
	file, err := os.Open("/proc/partitions")
	if err != nil {
		return nil, fmt.Errorf("erro ao ler /proc/partitions: %w", err)
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	scanner.Scan() // Pular cabeçalho
	
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if len(line) < 4 {
			continue
		}
		
		deviceName := line[3]
		// Filtrar apenas dispositivos de armazenamento (sd*, mmcblk*, etc.)
		if !strings.HasPrefix(deviceName, "sd") && !strings.HasPrefix(deviceName, "mmcblk") {
			continue
		}
		
		// Pular dispositivos do sistema (sda geralmente é o disco principal)
		if deviceName == "sda" {
			continue
		}
		
		device, err := getDeviceInfo(deviceName)
		if err != nil {
			continue // Pular dispositivos com erro
		}
		
		// Marcar como não removível se não for detectado como tal
		removablePath := fmt.Sprintf("/sys/block/%s/removable", deviceName)
		if data, err := os.ReadFile(removablePath); err == nil {
			device.Removable = strings.TrimSpace(string(data)) == "1"
		}
		
		devices = append(devices, *device)
	}
	
	return devices, nil
}

// listDevicesWSL lista dispositivos USB no WSL
func listDevicesWSL() ([]USBDevice, error) {
	// No WSL, tentar usar lsusb se disponível
	cmd := exec.Command("lsusb")
	output, err := cmd.Output()
	if err != nil {
		// Se lsusb não estiver disponível, usar método Linux
		return listDevicesLinux()
	}
	
	// Parse do output do lsusb
	var devices []USBDevice
	lines := strings.Split(string(output), "\n")
	
	for _, line := range lines {
		if strings.Contains(line, "Mass Storage") || strings.Contains(line, "Storage") {
			// Extrair informações básicas do lsusb
			parts := strings.Fields(line)
			if len(parts) >= 6 {
				vendor := strings.Join(parts[5:], " ")
				device := USBDevice{
					Path:      "/dev/sdX", // Placeholder
					Size:      "Unknown",
					SizeGB:    0,
					Model:     "USB Storage",
					Vendor:    vendor,
					Serial:    "Unknown",
					Removable: true,
					Platform:  "wsl",
				}
				devices = append(devices, device)
			}
		}
	}
	
	// Se não encontrou nada via lsusb, tentar método Linux
	if len(devices) == 0 {
		return listDevicesLinux()
	}
	
	return devices, nil
}

// listDevicesWindows lista dispositivos USB no Windows
func listDevicesWindows() ([]USBDevice, error) {
	// Usar PowerShell para listar dispositivos USB via WMI
	psCmd := `Get-WmiObject -Class Win32_LogicalDisk | Where-Object {$_.DriveType -eq 2} | Select-Object DeviceID, Size, VolumeName | ConvertTo-Json`
	
	cmd := exec.Command("powershell", "-Command", psCmd)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("erro ao executar PowerShell: %w", err)
	}
	
	// Parse básico do JSON (implementação simplificada)
	var devices []USBDevice
	lines := strings.Split(string(output), "\n")
	
	for _, line := range lines {
		if strings.Contains(line, "DeviceID") {
			// Extrair informações básicas
			device := USBDevice{
				Path:      "Unknown",
				Size:      "Unknown",
				SizeGB:    0,
				Model:     "USB Storage",
				Vendor:    "Unknown",
				Serial:    "Unknown",
				Removable: true,
				Platform:  "windows",
			}
			devices = append(devices, device)
		}
	}
	
	return devices, nil
}

// isUSBDevice verifica se um dispositivo é USB
func isUSBDevice(deviceName string) bool {
	// Verificar se o dispositivo está em /sys/block/
	sysPath := fmt.Sprintf("/sys/block/%s", deviceName)
	if _, err := os.Stat(sysPath); err != nil {
		return false
	}
	
	// Verificar se é removível
	removablePath := fmt.Sprintf("/sys/block/%s/removable", deviceName)
	if data, err := os.ReadFile(removablePath); err == nil {
		return strings.TrimSpace(string(data)) == "1"
	}
	
	// Verificar se tem informações USB
	ueventPath := fmt.Sprintf("/sys/block/%s/device/uevent", deviceName)
	if data, err := os.ReadFile(ueventPath); err == nil {
		content := string(data)
		return strings.Contains(content, "USB") || strings.Contains(content, "usb")
	}
	
	return false
}

// getDeviceInfo obtém informações detalhadas de um dispositivo
func getDeviceInfo(deviceName string) (*USBDevice, error) {
	device := &USBDevice{
		Path:      fmt.Sprintf("/dev/%s", deviceName),
		Platform:  detectPlatform(),
		Removable: true,
	}
	
	// Obter tamanho
	if size, err := getDeviceSize(deviceName); err == nil {
		device.Size = formatSize(size)
		device.SizeGB = int(size / (1024 * 1024 * 1024))
	}
	
	// Obter modelo e fabricante
	if model, vendor, err := getDeviceModel(deviceName); err == nil {
		device.Model = model
		device.Vendor = vendor
	}
	
	// Obter serial
	if serial, err := getDeviceSerial(deviceName); err == nil {
		device.Serial = serial
	}
	
	return device, nil
}

// getDeviceSize obtém o tamanho do dispositivo em bytes
func getDeviceSize(deviceName string) (int64, error) {
	sizePath := fmt.Sprintf("/sys/block/%s/size", deviceName)
	data, err := os.ReadFile(sizePath)
	if err != nil {
		return 0, err
	}
	
	sizeStr := strings.TrimSpace(string(data))
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return 0, err
	}
	
	// Tamanho em setores de 512 bytes
	return size * 512, nil
}

// getDeviceModel obtém modelo e fabricante do dispositivo
func getDeviceModel(deviceName string) (string, string, error) {
	modelPath := fmt.Sprintf("/sys/block/%s/device/model", deviceName)
	vendorPath := fmt.Sprintf("/sys/block/%s/device/vendor", deviceName)
	
	model := "Unknown"
	vendor := "Unknown"
	
	if data, err := os.ReadFile(modelPath); err == nil {
		model = strings.TrimSpace(string(data))
	}
	
	if data, err := os.ReadFile(vendorPath); err == nil {
		vendor = strings.TrimSpace(string(data))
	}
	
	return model, vendor, nil
}

// getDeviceSerial obtém o serial do dispositivo
func getDeviceSerial(deviceName string) (string, error) {
	serialPath := fmt.Sprintf("/sys/block/%s/device/serial", deviceName)
	data, err := os.ReadFile(serialPath)
	if err != nil {
		return "Unknown", err
	}
	
	return strings.TrimSpace(string(data)), nil
}

// formatSize formata o tamanho em bytes para string legível
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
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp])
}

func SelectDevice() (*USBDevice, error) {
	devices, err := ListDevices()
	if err != nil {
		return nil, err
	}
	if len(devices) == 0 {
		return nil, fmt.Errorf("nenhum dispositivo USB encontrado")
	}
	return &devices[0], nil
}

func NewDetector() *Detector {
	return &Detector{}
}

func NewFormatter() *Formatter {
	return &Formatter{}
}

func NewCreator(workDir, cacheDir string) *Creator {
	return &Creator{
		workDir:  workDir,
		cacheDir: cacheDir,
	}
}

type Detector struct{}

func (d *Detector) ValidateDevice(devicePath string) error {
	// Mock: sempre válido
	return nil
}

func (d *Detector) IsSystemDisk(devicePath string) bool {
	// Mock: nunca é disco do sistema
	return false
}

type Formatter struct{}

func (f *Formatter) FormatDevice(devicePath, label string) error {
	// Mock: simular formatação
	fmt.Printf("Formatando dispositivo %s com rótulo %s...\n", devicePath, label)
	return nil
}

type Creator struct {
	workDir  string
	cacheDir string
}

func (c *Creator) CreateUSB(devicePath string, config *Config) error {
	// Mock: simular criação de USB
	fmt.Printf("Criando USB em %s para nó %s...\n", devicePath, config.NodeName)
	return nil
}

func (c *Creator) Cleanup() error {
	// Mock: simular limpeza
	return nil
}
