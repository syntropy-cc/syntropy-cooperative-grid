package usb

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

// USBDevice representa um dispositivo USB detectado
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

// Detector interface para detecção de dispositivos USB
type Detector interface {
	DetectDevices() ([]USBDevice, error)
	ValidateDevice(devicePath string) error
	IsSystemDisk(devicePath string) bool
}

// LinuxDetector implementa detecção para Linux nativo
type LinuxDetector struct{}

// WSLDetector implementa detecção para WSL (Windows Subsystem for Linux)
type WSLDetector struct{}

// WindowsDetector implementa detecção para Windows nativo
type WindowsDetector struct{}

// NewDetector cria um detector apropriado para a plataforma atual
func NewDetector() Detector {
	switch runtime.GOOS {
	case "linux":
		if isWSL() {
			return &WSLDetector{}
		}
		return &LinuxDetector{}
	case "windows":
		return &WindowsDetector{}
	default:
		return &LinuxDetector{} // fallback para Linux
	}
}

// isWSL verifica se está rodando no WSL
func isWSL() bool {
	data, err := os.ReadFile("/proc/version")
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(data)), "microsoft")
}

// DetectDevices implementa a detecção para Linux nativo
func (d *LinuxDetector) DetectDevices() ([]USBDevice, error) {
	var devices []USBDevice

	// Usar lsblk para obter informações dos dispositivos
	cmd := exec.Command("lsblk", "-d", "-n", "-o", "NAME,SIZE,TYPE,RM,MODEL,SERIAL,VENDOR")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("falha ao executar lsblk: %w", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		
		if len(fields) < 4 {
			continue
		}

		// Verificar se é um disco
		if fields[2] != "disk" {
			continue
		}

		devicePath := "/dev/" + fields[0]
		
		// Verificar se é removível
		removable := fields[3] == "1"
		
		// Verificar se não é disco do sistema
		if !removable && fields[0] == "sda" {
			continue
		}

		// Calcular tamanho em GB
		sizeGB := parseSizeToGB(fields[1])
		
		// Filtrar dispositivos muito pequenos ou muito grandes
		if sizeGB < 1 || sizeGB > 1024 {
			continue
		}

		// Construir informações do dispositivo
		device := USBDevice{
			Path:      devicePath,
			Size:      fields[1],
			SizeGB:    sizeGB,
			Model:     getField(fields, 4, "Unknown"),
			Vendor:    getField(fields, 5, "Unknown"),
			Serial:    getField(fields, 6, "Unknown"),
			Removable: removable,
			Platform:  "linux",
		}

		// Validar se o dispositivo existe e é acessível
		if d.ValidateDevice(devicePath) == nil {
			devices = append(devices, device)
		}
	}

	return devices, scanner.Err()
}

// DetectDevices implementa a detecção para WSL
func (d *WSLDetector) DetectDevices() ([]USBDevice, error) {
	var devices []USBDevice

	// Primeiro, tentar obter dispositivos via PowerShell (Windows)
	windowsDevices, err := d.getWindowsUSBDevices()
	if err == nil && len(windowsDevices) > 0 {
		// Mapear dispositivos Windows para WSL
		for _, winDevice := range windowsDevices {
			wslDevice := d.convertWindowsToWSL(winDevice)
			if wslDevice != nil {
				devices = append(devices, *wslDevice)
			}
		}
	}

	// Se não encontrou via Windows, usar detecção Linux como fallback
	if len(devices) == 0 {
		linuxDetector := &LinuxDetector{}
		linuxDevices, err := linuxDetector.DetectDevices()
		if err == nil {
			// Marcar como WSL
			for _, device := range linuxDevices {
				device.Platform = "wsl"
				devices = append(devices, device)
			}
		}
	}

	return devices, nil
}

// getWindowsUSBDevices obtém dispositivos USB via PowerShell
func (d *WSLDetector) getWindowsUSBDevices() ([]WindowsUSBDevice, error) {
	psCommand := `
	$ErrorActionPreference = "SilentlyContinue"
	$usbDrives = Get-WmiObject Win32_DiskDrive | Where-Object { $_.InterfaceType -eq "USB" }
	
	if ($usbDrives.Count -eq 0) {
		Write-Output "NO_USB_FOUND"
		exit 1
	}
	
	foreach ($disk in $usbDrives) {
		$diskNum = $disk.Index
		$model = $disk.Model
		$size = [math]::Round($disk.Size / 1GB, 2)
		Write-Output "USB|$diskNum|$model|$size"
	}`

	cmd := exec.Command("powershell.exe", "-ExecutionPolicy", "Bypass", "-Command", psCommand)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("falha ao executar PowerShell: %w", err)
	}

	var devices []WindowsUSBDevice
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	
	for _, line := range lines {
		if strings.HasPrefix(line, "USB|") {
			parts := strings.Split(line, "|")
			if len(parts) >= 4 {
				diskNum, _ := strconv.Atoi(parts[1])
				size, _ := strconv.ParseFloat(parts[3], 64)
				
				device := WindowsUSBDevice{
					DiskNumber: diskNum,
					Model:      parts[2],
					SizeGB:     int(size),
				}
				devices = append(devices, device)
			}
		}
	}

	return devices, nil
}

// WindowsUSBDevice representa um dispositivo USB detectado no Windows
type WindowsUSBDevice struct {
	DiskNumber int    `json:"disk_number"`
	Model      string `json:"model"`
	SizeGB     int    `json:"size_gb"`
}

// convertWindowsToWSL converte dispositivo Windows para formato WSL
func (d *WSLDetector) convertWindowsToWSL(winDevice WindowsUSBDevice) *USBDevice {
	// PhysicalDrive0 = /dev/sda, PhysicalDrive1 = /dev/sdb, etc.
	letterASCII := 97 + winDevice.DiskNumber // 'a' + disk_number
	deviceLetter := string(rune(letterASCII))
	wslPath := "/dev/sd" + deviceLetter

	// Verificar se o dispositivo existe no WSL
	if _, err := os.Stat(wslPath); os.IsNotExist(err) {
		return nil
	}

	return &USBDevice{
		Path:      wslPath,
		Size:      fmt.Sprintf("%dG", winDevice.SizeGB),
		SizeGB:    winDevice.SizeGB,
		Model:     winDevice.Model,
		Vendor:    "Unknown",
		Serial:    "Unknown",
		Removable: true, // Assumir que dispositivos USB são removíveis
		Platform:  "wsl",
	}
}

// DetectDevices implementa a detecção para Windows nativo
func (d *WindowsDetector) DetectDevices() ([]USBDevice, error) {
	// Implementação para Windows nativo usando PowerShell
	psCommand := `
	$ErrorActionPreference = "SilentlyContinue"
	$usbDrives = Get-WmiObject Win32_DiskDrive | Where-Object { $_.InterfaceType -eq "USB" }
	
	foreach ($disk in $usbDrives) {
		$device = @{
			Path = $disk.DeviceID
			Size = [math]::Round($disk.Size / 1GB, 2)
			Model = $disk.Model
			Vendor = $disk.Manufacturer
			Serial = $disk.SerialNumber
			Removable = $true
			Platform = "windows"
		}
		$device | ConvertTo-Json -Compress
	}`

	cmd := exec.Command("powershell.exe", "-ExecutionPolicy", "Bypass", "-Command", psCommand)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("falha ao executar PowerShell: %w", err)
	}

	var devices []USBDevice
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		
		var device USBDevice
		if err := json.Unmarshal([]byte(line), &device); err == nil {
			devices = append(devices, device)
		}
	}

	return devices, nil
}

// ValidateDevice valida se um dispositivo é acessível e seguro para uso
func (d *LinuxDetector) ValidateDevice(devicePath string) error {
	// Verificar se o dispositivo existe
	if _, err := os.Stat(devicePath); os.IsNotExist(err) {
		return fmt.Errorf("dispositivo não encontrado: %s", devicePath)
	}

	// Verificar se é um dispositivo de bloco
	var stat syscall.Stat_t
	if err := syscall.Stat(devicePath, &stat); err != nil {
		return fmt.Errorf("erro ao acessar dispositivo: %w", err)
	}

	if stat.Mode&syscall.S_IFBLK == 0 {
		return fmt.Errorf("não é um dispositivo de bloco: %s", devicePath)
	}

	// Verificar se não é disco do sistema
	if d.IsSystemDisk(devicePath) {
		return fmt.Errorf("dispositivo parece ser um disco do sistema: %s", devicePath)
	}

	return nil
}

// ValidateDevice para WSL
func (d *WSLDetector) ValidateDevice(devicePath string) error {
	// Usar validação Linux como base
	linuxDetector := &LinuxDetector{}
	if err := linuxDetector.ValidateDevice(devicePath); err != nil {
		return err
	}

	// Validações específicas do WSL
	deviceName := strings.TrimPrefix(devicePath, "/dev/")
	removablePath := fmt.Sprintf("/sys/block/%s/removable", deviceName)
	
	if removableData, err := os.ReadFile(removablePath); err == nil {
		if strings.TrimSpace(string(removableData)) == "0" {
			// Dispositivo não é removível, aplicar validação mais rigorosa
			if deviceName == "sda" {
				return fmt.Errorf("dispositivo sda não pode ser usado no WSL")
			}
		}
	}

	return nil
}

// ValidateDevice para Windows
func (d *WindowsDetector) ValidateDevice(devicePath string) error {
	// Implementação básica para Windows
	if !strings.HasPrefix(devicePath, "\\\\.\\PHYSICALDRIVE") && 
	   !strings.HasPrefix(devicePath, "PhysicalDrive") {
		return fmt.Errorf("formato de dispositivo inválido para Windows: %s", devicePath)
	}
	return nil
}

// IsSystemDisk verifica se um dispositivo é um disco do sistema
func (d *LinuxDetector) IsSystemDisk(devicePath string) bool {
	// Verificar se está em /etc/fstab
	if fstabData, err := os.ReadFile("/etc/fstab"); err == nil {
		if strings.Contains(string(fstabData), devicePath) {
			return true
		}
	}

	// Verificar pontos de montagem críticos
	criticalMounts := []string{"/", "/boot", "/boot/efi", "/usr", "/var", "/home", "/opt"}
	
	cmd := exec.Command("lsblk", "-n", "-o", "NAME,MOUNTPOINT", devicePath)
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			for _, mount := range criticalMounts {
				if fields[1] == mount {
					return true
				}
			}
		}
	}

	return false
}

// IsSystemDisk para WSL
func (d *WSLDetector) IsSystemDisk(devicePath string) bool {
	// Usar validação Linux como base
	linuxDetector := &LinuxDetector{}
	if linuxDetector.IsSystemDisk(devicePath) {
		return true
	}

	// Validações específicas do WSL
	deviceName := strings.TrimPrefix(devicePath, "/dev/")
	
	// Verificar se é sda (disco principal do WSL)
	if deviceName == "sda" {
		return true
	}

	// Verificar se é um dispositivo de swap
	cmd := exec.Command("lsblk", "-n", "-o", "FSTYPE", devicePath)
	output, err := cmd.Output()
	if err == nil && strings.Contains(string(output), "swap") {
		return true
	}

	return false
}

// IsSystemDisk para Windows
func (d *WindowsDetector) IsSystemDisk(devicePath string) bool {
	// Implementação básica para Windows
	// Verificar se é PhysicalDrive0 (disco principal)
	if strings.Contains(devicePath, "PhysicalDrive0") {
		return true
	}
	return false
}

// Funções auxiliares

// parseSizeToGB converte string de tamanho para GB
func parseSizeToGB(sizeStr string) int {
	// Remover espaços e converter para minúsculas
	sizeStr = strings.ToLower(strings.TrimSpace(sizeStr))
	
	// Regex para extrair número e unidade
	re := regexp.MustCompile(`^(\d+(?:\.\d+)?)\s*([kmgt]?b?)$`)
	matches := re.FindStringSubmatch(sizeStr)
	
	if len(matches) != 3 {
		return 0
	}
	
	number, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0
	}
	
	unit := matches[2]
	switch unit {
	case "k", "kb":
		return int(number / 1024 / 1024)
	case "m", "mb":
		return int(number / 1024)
	case "g", "gb":
		return int(number)
	case "t", "tb":
		return int(number * 1024)
	default:
		return int(number) // Assumir GB se não especificado
	}
}

// getField obtém um campo de um slice com fallback
func getField(fields []string, index int, fallback string) string {
	if index < len(fields) {
		return fields[index]
	}
	return fallback
}

// ListDevices lista todos os dispositivos USB disponíveis
func ListDevices() ([]USBDevice, error) {
	detector := NewDetector()
	return detector.DetectDevices()
}

// SelectDevice permite ao usuário selecionar um dispositivo
func SelectDevice() (*USBDevice, error) {
	devices, err := ListDevices()
	if err != nil {
		return nil, fmt.Errorf("falha ao detectar dispositivos: %w", err)
	}

	if len(devices) == 0 {
		return nil, fmt.Errorf("nenhum dispositivo USB encontrado")
	}

	if len(devices) == 1 {
		return &devices[0], nil
	}

	// Implementar seleção interativa
	fmt.Println("Dispositivos USB disponíveis:")
	for i, device := range devices {
		fmt.Printf("[%d] %s (%s) - %s - %s\n", 
			i+1, device.Path, device.Size, device.Model, device.Vendor)
	}

	// Por enquanto, retornar o primeiro dispositivo
	// TODO: Implementar seleção interativa
	return &devices[0], nil
}
