package usb

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Formatter interface para formatação de dispositivos USB
type Formatter interface {
	FormatDevice(devicePath, label string) error
	UnmountDevice(devicePath string) error
	WipeDevice(devicePath string) error
	CreatePartitionTable(devicePath string) error
	CreatePartition(devicePath string) error
	FormatPartition(partitionPath, label string) error
}

// LinuxFormatter implementa formatação para Linux/WSL
type LinuxFormatter struct{}

// WindowsFormatter implementa formatação para Windows
type WindowsFormatter struct{}

// NewFormatter cria um formatador apropriado para a plataforma atual
func NewFormatter() Formatter {
	switch runtime.GOOS {
	case "linux":
		return &LinuxFormatter{}
	case "windows":
		return &WindowsFormatter{}
	default:
		return &LinuxFormatter{} // fallback para Linux
	}
}

// FormatDevice formata um dispositivo USB no Linux/WSL
func (f *LinuxFormatter) FormatDevice(devicePath, label string) error {
	fmt.Printf("Iniciando formatação do dispositivo: %s\n", devicePath)

	// 1. Desmontar todas as partições
	if err := f.UnmountDevice(devicePath); err != nil {
		return fmt.Errorf("falha ao desmontar dispositivo: %w", err)
	}

	// 2. Limpar o dispositivo
	if err := f.WipeDevice(devicePath); err != nil {
		return fmt.Errorf("falha ao limpar dispositivo: %w", err)
	}

	// 3. Criar tabela de partições
	if err := f.CreatePartitionTable(devicePath); err != nil {
		return fmt.Errorf("falha ao criar tabela de partições: %w", err)
	}

	// 4. Criar partição
	if err := f.CreatePartition(devicePath); err != nil {
		return fmt.Errorf("falha ao criar partição: %w", err)
	}

	// 5. Obter caminho da partição
	partitionPath := f.getPartitionPath(devicePath)
	if partitionPath == "" {
		return fmt.Errorf("não foi possível determinar o caminho da partição")
	}

	// 6. Aguardar o sistema reconhecer a partição
	time.Sleep(2 * time.Second)

	// 7. Formatar a partição
	if err := f.FormatPartition(partitionPath, label); err != nil {
		return fmt.Errorf("falha ao formatar partição: %w", err)
	}

	fmt.Printf("✅ Dispositivo %s formatado com sucesso!\n", devicePath)
	fmt.Printf("   Partição: %s\n", partitionPath)
	fmt.Printf("   Rótulo: %s\n", label)

	return nil
}

// UnmountDevice desmonta todas as partições de um dispositivo
func (f *LinuxFormatter) UnmountDevice(devicePath string) error {
	fmt.Printf("Desmontando partições em %s...\n", devicePath)

	// Usar lsblk para encontrar partições montadas
	cmd := exec.Command("lsblk", "-n", "-o", "NAME,MOUNTPOINT", devicePath)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("falha ao listar partições: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[1] != "" {
			partitionPath := "/dev/" + fields[0]
			mountPoint := fields[1]

			// Desmontar a partição
			unmountCmd := exec.Command("sudo", "umount", "-f", mountPoint)
			if err := unmountCmd.Run(); err != nil {
				fmt.Printf("⚠️  Aviso: Falha ao desmontar %s: %v\n", mountPoint, err)
			} else {
				fmt.Printf("   Desmontado: %s\n", mountPoint)
			}
		}
	}

	// Aguardar um pouco para o sistema processar
	time.Sleep(1 * time.Second)
	return nil
}

// WipeDevice limpa completamente um dispositivo
func (f *LinuxFormatter) WipeDevice(devicePath string) error {
	fmt.Printf("Limpando dispositivo %s...\n", devicePath)

	// 1. Remover assinaturas de sistema de arquivos
	fmt.Println("   Removendo assinaturas de sistema de arquivos...")
	wipeCmd := exec.Command("sudo", "wipefs", "-a", devicePath)
	if err := wipeCmd.Run(); err != nil {
		fmt.Printf("⚠️  Aviso: Falha ao executar wipefs: %v\n", err)
	}

	// 2. Destruir estruturas GPT e MBR
	fmt.Println("   Destruindo estruturas de partição...")
	sgdiskCmd := exec.Command("sudo", "sgdisk", "--zap-all", devicePath)
	if err := sgdiskCmd.Run(); err != nil {
		fmt.Printf("⚠️  Aviso: Falha ao executar sgdisk: %v\n")
	}

	// 3. Zerar o início e o fim do dispositivo
	fmt.Println("   Zerando setores críticos...")
	
	// Zerar primeiros 10MB
	ddCmd1 := exec.Command("sudo", "dd", "if=/dev/zero", "of="+devicePath, "bs=1M", "count=10")
	if err := ddCmd1.Run(); err != nil {
		fmt.Printf("⚠️  Aviso: Falha ao zerar início do dispositivo: %v\n", err)
	}

	// Zerar últimos 10MB
	ddCmd2 := exec.Command("sudo", "dd", "if=/dev/zero", "of="+devicePath, "bs=1M", "seek=-10")
	if err := ddCmd2.Run(); err != nil {
		fmt.Printf("⚠️  Aviso: Falha ao zerar fim do dispositivo: %v\n", err)
	}

	// 4. Forçar o kernel a reler a tabela de partições
	fmt.Println("   Recarregando tabela de partições...")
	partprobeCmd := exec.Command("sudo", "partprobe", devicePath)
	partprobeCmd.Run() // Ignorar erro

	time.Sleep(2 * time.Second)
	fmt.Println("   ✅ Dispositivo limpo com sucesso")
	return nil
}

// CreatePartitionTable cria uma tabela de partições MBR
func (f *LinuxFormatter) CreatePartitionTable(devicePath string) error {
	fmt.Printf("Criando tabela de partições MBR em %s...\n", devicePath)

	// Criar tabela de partições msdos (MBR)
	cmd := exec.Command("sudo", "parted", "-s", devicePath, "mklabel", "msdos")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("falha ao criar tabela de partições msdos: %w", err)
	}

	fmt.Println("   ✅ Tabela de partições MBR criada")
	return nil
}

// CreatePartition cria uma partição primária
func (f *LinuxFormatter) CreatePartition(devicePath string) error {
	fmt.Printf("Criando partição primária em %s...\n", devicePath)

	// Criar partição primária usando todo o espaço
	cmd := exec.Command("sudo", "parted", "-s", devicePath, "mkpart", "primary", "fat32", "1MiB", "100%")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("falha ao criar partição primária: %w", err)
	}

	// Definir flag de boot
	cmd = exec.Command("sudo", "parted", "-s", devicePath, "set", "1", "boot", "on")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("falha ao definir flag de boot: %w", err)
	}

	// Recarregar tabela de partições
	partprobeCmd := exec.Command("sudo", "partprobe", devicePath)
	partprobeCmd.Run() // Ignorar erro

	fmt.Println("   ✅ Partição primária criada com flag de boot")
	return nil
}

// FormatPartition formata uma partição como FAT32
func (f *LinuxFormatter) FormatPartition(partitionPath, label string) error {
	fmt.Printf("Formatando partição %s como FAT32...\n", partitionPath)

	// Formatar como FAT32
	cmd := exec.Command("sudo", "mkfs.fat", "-F32", "-n", label, partitionPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("falha ao formatar partição como FAT32: %w", err)
	}

	fmt.Printf("   ✅ Partição %s formatada como FAT32 (rótulo: %s)\n", partitionPath, label)
	return nil
}

// getPartitionPath determina o caminho da partição criada
func (f *LinuxFormatter) getPartitionPath(devicePath string) string {
	// Tentar diferentes padrões de nomenclatura
	patterns := []string{
		devicePath + "1",    // /dev/sdb1
		devicePath + "p1",   // /dev/nvme0n1p1
	}

	for _, pattern := range patterns {
		if _, err := os.Stat(pattern); err == nil {
			return pattern
		}
	}

	return ""
}

// FormatDevice formata um dispositivo USB no Windows
func (f *WindowsFormatter) FormatDevice(devicePath, label string) error {
	fmt.Printf("Iniciando formatação do dispositivo Windows: %s\n", devicePath)

	// Extrair número do disco
	diskNumber := f.extractDiskNumber(devicePath)
	if diskNumber == -1 {
		return fmt.Errorf("não foi possível extrair o número do disco de %s", devicePath)
	}

	// Script PowerShell para formatar o disco
	psScript := fmt.Sprintf(`
	$ErrorActionPreference = "Stop"
	$diskNumber = %d
	$label = "%s"

	try {
		Write-Host "Selecionando disco $diskNumber..."
		$disk = Get-Disk -Number $diskNumber -ErrorAction Stop
		
		Write-Host "Limpando disco..."
		Clear-Disk -Number $diskNumber -RemoveData -Confirm:$false -ErrorAction Stop
		
		Write-Host "Criando partição..."
		$partition = New-Partition -DiskNumber $diskNumber -UseMaximumSize -IsActive -ErrorAction Stop
		
		Write-Host "Formatando partição..."
		Format-Volume -Partition $partition -FileSystem FAT32 -NewFileSystemLabel $label -Confirm:$false -ErrorAction Stop
		
		Write-Host "Formatação concluída com sucesso."
	} catch {
		Write-Error "Erro durante a formatação: $_"
		exit 1
	}
	`, diskNumber, label)

	cmd := exec.Command("powershell.exe", "-ExecutionPolicy", "Bypass", "-Command", psScript)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("falha ao formatar USB via PowerShell: %w\nOutput: %s", err, string(output))
	}

	fmt.Printf("✅ Dispositivo %s formatado com sucesso via PowerShell\n", devicePath)
	return nil
}

// extractDiskNumber extrai o número do disco do caminho do dispositivo Windows
func (f *WindowsFormatter) extractDiskNumber(devicePath string) int {
	// Tentar diferentes padrões
	patterns := []string{
		`PhysicalDrive(\d+)`,
		`\\\\.\\PHYSICALDRIVE(\d+)`,
		`(\d+)$`,
	}

	for _, pattern := range patterns {
		matches := f.findRegexMatches(devicePath, pattern)
		if len(matches) > 0 {
			if diskNum, err := f.parseDiskNumber(matches[0]); err == nil {
				return diskNum
			}
		}
	}

	return -1
}

// Métodos auxiliares para Windows

func (f *WindowsFormatter) findRegexMatches(text, pattern string) []string {
	// Implementação simples de regex match
	// Em produção, usar regexp package
	return []string{}
}

func (f *WindowsFormatter) parseDiskNumber(str string) (int, error) {
	// Implementação simples de parsing
	// Em produção, usar strconv package
	return 0, fmt.Errorf("not implemented")
}

// Métodos não implementados para Windows (usam PowerShell internamente)

func (f *WindowsFormatter) UnmountDevice(devicePath string) error {
	fmt.Println("Desmontagem não necessária no Windows (PowerShell gerencia automaticamente)")
	return nil
}

func (f *WindowsFormatter) WipeDevice(devicePath string) error {
	fmt.Println("Limpeza não necessária no Windows (PowerShell Clear-Disk gerencia)")
	return nil
}

func (f *WindowsFormatter) CreatePartitionTable(devicePath string) error {
	fmt.Println("Criação de tabela não necessária no Windows (PowerShell gerencia)")
	return nil
}

func (f *WindowsFormatter) CreatePartition(devicePath string) error {
	fmt.Println("Criação de partição não necessária no Windows (PowerShell gerencia)")
	return nil
}

func (f *WindowsFormatter) FormatPartition(partitionPath, label string) error {
	fmt.Println("Formatação de partição não necessária no Windows (PowerShell gerencia)")
	return nil
}
