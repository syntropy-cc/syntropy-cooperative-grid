package usb

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// detectPlatform detecta a plataforma atual (linux, windows, wsl)
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

// ListDevices lista dispositivos USB disponíveis baseado na plataforma
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

// validateDevice valida se um dispositivo é seguro para usar
func validateDevice(devicePath string) error {
	platform := detectPlatform()

	// Validação específica por plataforma
	switch platform {
	case "linux":
		return validateDeviceLinux(devicePath)
	case "wsl", "windows":
		return validateDeviceWindows(devicePath)
	default:
		return fmt.Errorf("plataforma não suportada: %s", platform)
	}
}

// validateDeviceLinux valida dispositivos no Linux
func validateDeviceLinux(devicePath string) error {
	// Verificar se o dispositivo existe
	if _, err := os.Stat(devicePath); err != nil {
		return fmt.Errorf("dispositivo não encontrado: %s", devicePath)
	}

	// Verificar se não é um dispositivo do sistema
	if isSystemDeviceLinux(devicePath) {
		return fmt.Errorf("dispositivo %s parece ser um dispositivo do sistema - operação cancelada por segurança", devicePath)
	}

	return nil
}

// validateDeviceWindows valida dispositivos no Windows/WSL
func validateDeviceWindows(devicePath string) error {
	// Extrair número do disco
	var diskNum int
	if strings.HasPrefix(devicePath, "PHYSICALDRIVE") {
		fmt.Sscanf(devicePath, "PHYSICALDRIVE%d", &diskNum)
	} else if strings.HasPrefix(devicePath, "\\\\.\\PHYSICALDRIVE") {
		fmt.Sscanf(devicePath, "\\\\.\\PHYSICALDRIVE%d", &diskNum)
	} else {
		return fmt.Errorf("formato de dispositivo inválido para Windows: %s", devicePath)
	}

	// Verificar se não é um dispositivo do sistema usando PowerShell
	if isSystemDeviceWindows(diskNum) {
		return fmt.Errorf("dispositivo %s parece ser um dispositivo do sistema - operação cancelada por segurança", devicePath)
	}

	return nil
}

// isSystemDeviceLinux verifica se um dispositivo é do sistema no Linux
func isSystemDeviceLinux(devicePath string) bool {
	// Verificar se contém partições do sistema
	if strings.Contains(devicePath, "/dev/sda") || strings.Contains(devicePath, "/dev/nvme0n1") {
		// Verificar se tem partições do sistema
		cmd := exec.Command("lsblk", "-n", "-o", "MOUNTPOINT", devicePath)
		output, err := cmd.Output()
		if err == nil {
			mountpoints := strings.Split(strings.TrimSpace(string(output)), "\n")
			for _, mp := range mountpoints {
				if mp == "/" || mp == "/boot" || mp == "/home" {
					return true
				}
			}
		}
	}
	return false
}

// isSystemDeviceWindows verifica se um dispositivo é do sistema no Windows
func isSystemDeviceWindows(diskNum int) bool {
	// Usar PowerShell para verificar se é disco do sistema
	psScript := fmt.Sprintf(`
$disk = Get-Disk -Number %d -ErrorAction SilentlyContinue
if ($disk) {
    $isSystem = $disk.IsSystem -or $disk.IsBoot -or $disk.IsOffline
    $partitions = Get-Partition -DiskNumber %d -ErrorAction SilentlyContinue
    foreach ($part in $partitions) {
        if ($part.IsSystem -or $part.IsBoot -or $part.DriveLetter -eq "C") {
            $isSystem = $true
            break
        }
    }
    if ($isSystem) { exit 1 } else { exit 0 }
} else {
    exit 1
}
`, diskNum, diskNum)

	cmd := exec.Command("powershell.exe", "-NoProfile", "-Command", psScript)
	err := cmd.Run()
	return err != nil // Se houver erro, é dispositivo do sistema
}

// formatUSB formata um dispositivo USB
func formatUSB(devicePath, label string, force bool) error {
	platform := detectPlatform()

	// Validar dispositivo
	if err := validateDevice(devicePath); err != nil {
		return err
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

// listUSBDevices lista dispositivos USB e formata a saída
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

// SelectDevice permite ao usuário selecionar um dispositivo USB
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

	fmt.Printf("\nSelecione o dispositivo (1-%d): ", len(devices))
	var choice int
	fmt.Scanln(&choice)

	if choice < 1 || choice > len(devices) {
		return nil, fmt.Errorf("seleção inválida")
	}

	return &devices[choice-1], nil
}
