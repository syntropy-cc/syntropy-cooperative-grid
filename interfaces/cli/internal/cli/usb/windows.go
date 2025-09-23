package usb

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// listDevicesWSL lista dispositivos USB no WSL
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

// listDevicesWSLAlternative método alternativo usando WMIC
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

// listDevicesWindows lista dispositivos USB no Windows nativo
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

// createUSBWithNoCloudWSL cria USB usando estratégia NoCloud no WSL
func createUSBWithNoCloudWSL(devicePath string, config *Config, workDir, cacheDir string) error {
	// Extrair número do disco Windows
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

	// Gerenciar ISO
	isoPath := config.ISOPath
	if isoPath == "" {
		var err error
		isoPath, err = manageISOCache(cacheDir)
		if err != nil {
			return fmt.Errorf("erro ao gerenciar ISO: %w", err)
		}
	}
	isoWSL := convertAnyToWSLPath(isoPath)

	fmt.Printf("📀 ISO (WSL): %s\n", isoWSL)
	fmt.Printf("🧱 Disco: %s (nº %d)\n\n", winPhysical, diskNum)

	// Criar script bash separado para melhor debugging
	bashScript := fmt.Sprintf(`#!/bin/bash
set -euo pipefail

echo "🔍 Detectando dispositivo WSL..."

# Listar dispositivos antes
before=($(ls /dev/sd? /dev/hd? /dev/nvme?n? 2>/dev/null || true))
echo "Dispositivos antes: ${before[*]}"

# Aguardar um pouco para o dispositivo aparecer
sleep 2

# Listar dispositivos depois
after=($(ls /dev/sd? /dev/hd? /dev/nvme?n? 2>/dev/null || true))
echo "Dispositivos depois: ${after[*]}"

# Encontrar novo dispositivo
dev=""
for d in "${after[@]}"; do
  found=0
  for b in "${before[@]}"; do 
    [ "$d" = "$b" ] && found=1 && break
  done
  [ $found -eq 0 ] && dev="$d" && break
done

if [ -z "$dev" ]; then
  echo "❌ ERRO: Não foi possível detectar o dispositivo no WSL" >&2
  echo "Dispositivos disponíveis:" >&2
  ls -la /dev/sd* /dev/hd* /dev/nvme* 2>/dev/null || true >&2
  exit 1
fi

echo "✅ Dispositivo WSL detectado: $dev"

# Verificar se ISO existe
ISO="%s"
if [ ! -f "$ISO" ]; then
  echo "❌ ERRO: ISO não encontrada: $ISO" >&2
  exit 1
fi

echo "📀 Gravando ISO: $ISO -> $dev"
sudo dd if="$ISO" of="$dev" bs=4M status=progress conv=fsync
sync

echo "⏳ Aguardando gravação finalizar..."
sleep 3

echo "🔧 Criando partição CIDATA..."
sudo sgdisk -e "$dev"
sudo sgdisk -n 0:0:+128MiB -t 0:0700 -c 0:CIDATA "$dev"
sleep 2

# Determinar nome da partição CIDATA
cidata_part=""
if [[ "$dev" =~ nvme ]]; then
  cidata_part="${dev}p2"
else
  cidata_part="${dev}2"
fi

echo "📁 Partição CIDATA: $cidata_part"

# Verificar se partição existe
if [ ! -b "$cidata_part" ]; then
  echo "❌ ERRO: Partição CIDATA não encontrada: $cidata_part" >&2
  echo "Partições disponíveis:" >&2
  ls -la ${dev}* 2>/dev/null || true >&2
  exit 1
fi

echo "💾 Formatando partição CIDATA..."
sudo mkfs.vfat -F 32 -n CIDATA "$cidata_part"

echo "📂 Montando partição CIDATA..."
mount_point="$HOME/.syntropy/work/cidata-mount"
sudo mkdir -p "$mount_point"
sudo mount "$cidata_part" "$mount_point"

echo "📋 Copiando arquivos cloud-init..."
cloud_init_dir="%s/cloud-init"

# Verificar se diretório cloud-init existe
if [ ! -d "$cloud_init_dir" ]; then
  echo "❌ ERRO: Diretório cloud-init não encontrado: $cloud_init_dir" >&2
  sudo umount "$mount_point" || true
  sudo rmdir "$mount_point" || true
  exit 1
fi

# Verificar se arquivos existem
for file in user-data meta-data network-config; do
  if [ ! -f "$cloud_init_dir/$file" ]; then
    echo "❌ ERRO: Arquivo cloud-init não encontrado: $cloud_init_dir/$file" >&2
    sudo umount "$mount_point" || true
    sudo rmdir "$mount_point" || true
    exit 1
  fi
done

sudo cp "$cloud_init_dir/user-data" "$mount_point/"
sudo cp "$cloud_init_dir/meta-data" "$mount_point/"
sudo cp "$cloud_init_dir/network-config" "$mount_point/"

echo "🔍 Verificando arquivos copiados..."
ls -la "$mount_point/"

echo "🔓 Desmontando partição..."
sudo umount "$mount_point"
sudo rmdir "$mount_point"
sync

echo "✅ USB criado com sucesso usando estratégia NoCloud!"
`, isoWSL, convertAnyToWSLPath(workDir))

	// Salvar script bash
	bashScriptPath := filepath.Join(workDir, "create_usb_wsl.sh")
	if err := os.WriteFile(bashScriptPath, []byte(bashScript), 0755); err != nil {
		return fmt.Errorf("erro ao criar script bash: %w", err)
	}

	// Script PowerShell simplificado
	psScript := fmt.Sprintf(`
$ErrorActionPreference = "Stop"

Write-Host "🚀 Iniciando criação de USB com estratégia NoCloud..." -ForegroundColor Green
Write-Host "💾 Disco: %s (nº %d)" -ForegroundColor Cyan
Write-Host "📀 ISO: %s" -ForegroundColor Cyan

try {
    Write-Host "📴 Colocando disco offline no Windows..." -ForegroundColor Yellow
    Set-Disk -Number %d -IsReadOnly $false -IsOffline $true
    
    Write-Host "🔗 Montando disco cru no WSL..." -ForegroundColor Yellow
    $mountResult = wsl --mount %s --bare 2>&1
    if ($LASTEXITCODE -ne 0) {
        throw "Falha ao montar disco no WSL: $mountResult"
    }
    
    Write-Host "🐧 Executando script de criação no WSL..." -ForegroundColor Yellow
    $bashScript = "%s"
    $wslResult = wsl bash -lc "bash '$bashScript'" 2>&1
    $exitCode = $LASTEXITCODE
    
    if ($exitCode -ne 0) {
        Write-Host "❌ ERRO no WSL:" -ForegroundColor Red
        Write-Host $wslResult -ForegroundColor Red
        throw "Script WSL falhou com código: $exitCode"
    }
    
    Write-Host "✅ Script WSL executado com sucesso!" -ForegroundColor Green
    Write-Host $wslResult -ForegroundColor White
    
} catch {
    Write-Host "❌ ERRO: $($_.Exception.Message)" -ForegroundColor Red
    throw
} finally {
    Write-Host "🔄 Limpando recursos..." -ForegroundColor Yellow
    try { 
        wsl --unmount %s 2>$null
        Write-Host "✅ Dispositivo desmontado do WSL" -ForegroundColor Green
    } catch { 
        Write-Host "⚠️  Aviso: Falha ao desmontar do WSL (pode já estar desmontado)" -ForegroundColor Yellow 
    }
    
    try {
        Set-Disk -Number %d -IsOffline $false
        Write-Host "✅ Disco voltou online no Windows" -ForegroundColor Green
    } catch {
        Write-Host "⚠️  Aviso: Falha ao voltar disco online" -ForegroundColor Yellow
    }
}

Write-Host "🎉 USB criado com sucesso usando estratégia NoCloud!" -ForegroundColor Green
Write-Host "🔧 O USB agora contém:" -ForegroundColor Cyan
Write-Host "   • ISO Ubuntu original (bootável)" -ForegroundColor White
Write-Host "   • Partição CIDATA com configuração cloud-init" -ForegroundColor White
Write-Host "   • Configuração será aplicada automaticamente no boot" -ForegroundColor White
`, winPhysical, diskNum, isoWSL, diskNum, winPhysical, convertAnyToWSLPath(bashScriptPath), winPhysical, diskNum)

	// Gravar e executar o script elevado
	os.MkdirAll(workDir, 0755)
	scriptPath := filepath.Join(workDir, "create_usb_nocloud.ps1")
	if err := os.WriteFile(scriptPath, []byte(psScript), 0644); err != nil {
		return fmt.Errorf("erro ao criar script temporário: %w", err)
	}
	winScriptPath := convertWSLToWindowsPath(scriptPath)

	fmt.Println("📝 Solicitando permissões de administrador...")
	fmt.Println("⚠️  IMPORTANTE: O PowerShell será aberto com privilégios elevados.")
	fmt.Println("   Se aparecer uma mensagem de erro vermelha, verifique:")
	fmt.Println("   1. Se o dispositivo USB está conectado")
	fmt.Println("   2. Se não há outros programas usando o USB")
	fmt.Println("   3. Se você tem privilégios de administrador")
	fmt.Println()

	// Executar diagnóstico antes de tentar criar o USB
	if err := debugWSLEnvironment(workDir); err != nil {
		fmt.Printf("⚠️  Aviso: Falha no diagnóstico: %v\n", err)
	}

	cmd := exec.Command("powershell.exe", "-Command",
		fmt.Sprintf(`Start-Process powershell -ArgumentList '-NoProfile -ExecutionPolicy Bypass -File "%s"' -Verb RunAs -Wait`, winScriptPath))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("❌ Falha na criação do USB. Executando diagnóstico adicional...")
		if debugErr := debugWSLEnvironment(workDir); debugErr != nil {
			fmt.Printf("⚠️  Diagnóstico adicional também falhou: %v\n", debugErr)
		}
		return fmt.Errorf("erro ao executar criação do USB (NoCloud): %w", err)
	}

	return nil
}

// createUSBWithNoCloudWindows cria USB usando estratégia NoCloud no Windows nativo
func createUSBWithNoCloudWindows(devicePath string, config *Config, workDir, cacheDir string) error {
	// Extrair número do disco
	var diskNum int
	if strings.HasPrefix(devicePath, "\\\\.\\PHYSICALDRIVE") {
		fmt.Sscanf(devicePath, "\\\\.\\PHYSICALDRIVE%d", &diskNum)
	} else {
		return fmt.Errorf("formato de dispositivo inválido: %s", devicePath)
	}

	// Gerenciar ISO
	isoPath := config.ISOPath
	if isoPath == "" {
		var err error
		isoPath, err = manageISOCache(cacheDir)
		if err != nil {
			return fmt.Errorf("erro ao gerenciar ISO: %w", err)
		}
	}

	// Script PowerShell para criar USB com NoCloud
	psScript := fmt.Sprintf(`
$ErrorActionPreference = "Stop"

Write-Host "Criando USB com estratégia NoCloud..." -ForegroundColor Cyan
Write-Host "Disco: %s (nº %d)" -ForegroundColor Cyan
Write-Host "ISO: %s" -ForegroundColor Cyan

# Colocar disco offline
Write-Host "Colocando disco offline..." -ForegroundColor Cyan
Set-Disk -Number %d -IsReadOnly $false -IsOffline $true

try {
    # Usar dd para gravar ISO (via WSL ou ferramenta Windows)
    Write-Host "Gravando ISO no dispositivo..." -ForegroundColor Cyan
    
    # Verificar se WSL está disponível
    $wslAvailable = $false
    try {
        wsl --version | Out-Null
        $wslAvailable = $true
    } catch {
        Write-Host "WSL não disponível, usando método alternativo..." -ForegroundColor Yellow
    }
    
    if ($wslAvailable) {
        # Usar WSL para todo o processo
        Write-Host "Usando WSL para gravação completa..." -ForegroundColor Cyan
        wsl bash -lc 'set -euo pipefail
before=($(ls /dev/sd? /dev/hd? /dev/nvme?n? 2>/dev/null || true))
sleep 0.5
tries=0
while [ $tries -lt 20 ]; do
  after=($(ls /dev/sd? /dev/hd? /dev/nvme?n? 2>/dev/null || true))
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
    
    # Criar partição CIDATA
    sleep 2
    echo "Criando partição CIDATA..."
    sudo sgdisk -e "$dev"
    sudo sgdisk -n 0:0:+128MiB -t 0:0700 -c 0:CIDATA "$dev"
    sleep 1
    
    # Determinar nome da partição CIDATA
    cidata_part=""
    if [[ "$dev" =~ nvme ]]; then
      cidata_part="${dev}p2"
    else
      cidata_part="${dev}2"
    fi
    
    # Formatar partição CIDATA
    sudo mkfs.vfat -F 32 -n CIDATA "$cidata_part"
    
    # Montar e copiar arquivos cloud-init
    mount_point="$HOME/.syntropy/work/cidata-mount"
    sudo mkdir -p "$mount_point"
    sudo mount "$cidata_part" "$mount_point"
    
    # Copiar arquivos cloud-init
    cloud_init_dir="%s/cloud-init"
    sudo cp "$cloud_init_dir/user-data" "$mount_point/"
    sudo cp "$cloud_init_dir/meta-data" "$mount_point/"
    sudo cp "$cloud_init_dir/network-config" "$mount_point/"
    
    # Desmontar e limpar
    sudo umount "$mount_point"
    sudo rmdir "$mount_point"
    sync
    
    echo "USB criado com sucesso usando estratégia NoCloud!"
    exit 0
  fi
  tries=$((tries+1))
  sleep 0.5
done
echo "Falha ao detectar o device no WSL." 1>&2
exit 1
'
    } else {
        # Método alternativo usando PowerShell (limitado)
        Write-Host "Método alternativo não implementado ainda." -ForegroundColor Red
        throw "WSL necessário para gravação de ISO"
    }
    
    Write-Host "USB criado com sucesso via WSL!" -ForegroundColor Green
    
} finally {
    Write-Host "Voltando disco online..." -ForegroundColor Cyan
    Set-Disk -Number %d -IsOffline $false
}

Write-Host "✅ USB criado com sucesso usando estratégia NoCloud!" -ForegroundColor Green
Write-Host "🔧 O USB agora contém:" -ForegroundColor Cyan
Write-Host "   • ISO Ubuntu original (bootável)" -ForegroundColor White
Write-Host "   • Partição CIDATA com configuração cloud-init" -ForegroundColor White
Write-Host "   • Configuração será aplicada automaticamente no boot" -ForegroundColor White
`, devicePath, diskNum, isoPath, diskNum, isoPath, workDir, diskNum)

	// Gravar e executar o script elevado
	os.MkdirAll(workDir, 0755)
	scriptPath := filepath.Join(workDir, "create_usb_nocloud_windows.ps1")
	if err := os.WriteFile(scriptPath, []byte(psScript), 0644); err != nil {
		return fmt.Errorf("erro ao criar script temporário: %w", err)
	}

	fmt.Println("📝 Solicitando permissões de administrador...")
	cmd := exec.Command("powershell.exe", "-Command",
		fmt.Sprintf(`Start-Process powershell -ArgumentList '-NoProfile -ExecutionPolicy Bypass -File "%s"' -Verb RunAs -Wait`, scriptPath))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("erro ao executar criação do USB (NoCloud Windows): %w", err)
	}

	return nil
}

// createUSBWindows função legada - redirecionar para nova implementação
func createUSBWindows(devicePath string, config *Config, workDir, cacheDir string) error {
	// Função legada - redirecionar para nova implementação
	return createUSBWithNoCloudWindows(devicePath, config, workDir, cacheDir)
}

// formatUSBWSL formata um dispositivo USB no WSL
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

// formatUSBWindows formata um dispositivo USB no Windows nativo
func formatUSBWindows(devicePath, label string) error {
	// Implementação similar ao WSL mas sem conversão
	return fmt.Errorf("implementação Windows nativa pendente")
}

// convertAnyToWSLPath aceita caminho Windows (C:\...) ou já em WSL (/mnt/c/...)
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

// convertWSLToWindowsPath converte caminho WSL para Windows
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

// debugWSLEnvironment função para diagnosticar problemas no WSL
func debugWSLEnvironment(workDir string) error {
	fmt.Println("🔍 Executando diagnóstico do ambiente WSL...")

	// Script de diagnóstico
	debugScript := `#!/bin/bash
set -euo pipefail

echo "=== DIAGNÓSTICO DO AMBIENTE WSL ==="
echo "Data/Hora: $(date)"
echo "Usuário: $(whoami)"
echo "Diretório atual: $(pwd)"
echo ""

echo "=== VERIFICAÇÃO DE COMANDOS ==="
for cmd in dd sgdisk mkfs.vfat mount umount; do
  if command -v $cmd >/dev/null 2>&1; then
    echo "✅ $cmd: $(which $cmd)"
  else
    echo "❌ $cmd: NÃO ENCONTRADO"
  fi
done
echo ""

echo "=== PERMISSÕES SUDO ==="
if sudo -n true 2>/dev/null; then
  echo "✅ Sudo disponível sem senha"
else
  echo "⚠️  Sudo requer senha ou não disponível"
fi
echo ""

echo "=== DISPOSITIVOS DISPONÍVEIS ==="
echo "Dispositivos de bloco:"
ls -la /dev/sd* /dev/hd* /dev/nvme* 2>/dev/null || echo "Nenhum dispositivo encontrado"
echo ""

echo "=== MONTAGENS ATIVAS ==="
mount | grep -E "(sd|hd|nvme)" || echo "Nenhuma montagem de dispositivo encontrada"
echo ""

echo "=== ESPAÇO EM DISCO ==="
df -h /tmp /home 2>/dev/null || echo "Erro ao verificar espaço em disco"
echo ""

echo "=== INFORMAÇÕES DO SISTEMA ==="
uname -a
echo ""

echo "=== DIAGNÓSTICO CONCLUÍDO ==="
`

	// Salvar script de debug
	debugScriptPath := filepath.Join(workDir, "debug_wsl.sh")
	if err := os.WriteFile(debugScriptPath, []byte(debugScript), 0755); err != nil {
		return fmt.Errorf("erro ao criar script de debug: %w", err)
	}

	// Executar diagnóstico
	fmt.Println("🐧 Executando diagnóstico no WSL...")
	cmd := exec.Command("wsl", "bash", "-lc", fmt.Sprintf("bash '%s'", convertAnyToWSLPath(debugScriptPath)))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("⚠️  Aviso: Diagnóstico falhou: %v\n", err)
	} else {
		fmt.Println("✅ Diagnóstico concluído com sucesso")
	}

	return nil
}
