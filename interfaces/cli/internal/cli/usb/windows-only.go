package usb

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// WindowsOnlyUSBDevice representa um dispositivo USB específico para Windows
type WindowsOnlyUSBDevice struct {
	DiskNumber     int    `json:"disk_number"`
	FriendlyName   string `json:"friendly_name"`
	Size           int64  `json:"size"`
	SizeFormatted  string `json:"size_formatted"`
	SerialNumber   string `json:"serial_number"`
	BusType        string `json:"bus_type"`
	Model          string `json:"model"`
	IsSystem       bool   `json:"is_system"`
	IsBoot         bool   `json:"is_boot"`
	IsOffline      bool   `json:"is_offline"`
	PartitionCount int    `json:"partition_count"`
	Status         string `json:"status"`
}

// WindowsOnlyConfig configuração específica para Windows
type WindowsOnlyConfig struct {
	NodeName        string `json:"node_name"`
	NodeDescription string `json:"node_description"`
	Coordinates     string `json:"coordinates"`
	OwnerKeyFile    string `json:"owner_key_file"`
	Label           string `json:"label"`
	ISOPath         string `json:"iso_path"`
	DiscoveryServer string `json:"discovery_server"`
	CreatedBy       string `json:"created_by"`
	// Configurações específicas do Windows
	ExecutionPolicy string `json:"execution_policy"`
	PowerShellPath  string `json:"powershell_path"`
	WSLDistro       string `json:"wsl_distro"`
	TempDir         string `json:"temp_dir"`
	LogLevel        string `json:"log_level"`
}

// WindowsOnlyError representa um erro específico do Windows
type WindowsOnlyError struct {
	Code        string `json:"code"`
	Message     string `json:"message"`
	Suggestion  string `json:"suggestion"`
	ErrorType   string `json:"error_type"` // "permission", "device", "wsl", "powershell", "system"
	Recoverable bool   `json:"recoverable"`
}

func (e *WindowsOnlyError) Error() string {
	return fmt.Sprintf("[%s] %s - %s", e.Code, e.Message, e.Suggestion)
}

// validateWindowsEnvironment valida o ambiente Windows antes de executar
func validateWindowsEnvironment() error {
	// Verificar se estamos no Windows
	if runtime.GOOS != "windows" {
		return &WindowsOnlyError{
			Code:        "NOT_WINDOWS",
			Message:     "Este comando só pode ser executado no Windows",
			Suggestion:  "Execute este comando em um sistema Windows",
			ErrorType:   "system",
			Recoverable: false,
		}
	}

	// Verificar privilégios de administrador
	if !isRunningAsAdministrator() {
		return &WindowsOnlyError{
			Code:        "NO_ADMIN_PRIVILEGES",
			Message:     "Privilégios de administrador são necessários",
			Suggestion:  "Execute o PowerShell como Administrador (botão direito → Executar como administrador)",
			ErrorType:   "permission",
			Recoverable: true,
		}
	}

	// Verificar se WSL está disponível
	if !isWSLAvailable() {
		return &WindowsOnlyError{
			Code:        "WSL_NOT_AVAILABLE",
			Message:     "WSL não está disponível ou configurado",
			Suggestion:  "Instale o WSL: wsl --install ou wsl --install -d Ubuntu",
			ErrorType:   "wsl",
			Recoverable: true,
		}
	}

	// Verificar política de execução do PowerShell
	if err := checkPowerShellExecutionPolicy(); err != nil {
		return err
	}

	// Verificar ferramentas necessárias
	if err := checkRequiredTools(); err != nil {
		return err
	}

	return nil
}

// isRunningAsAdministrator verifica se o processo está executando como administrador
func isRunningAsAdministrator() bool {
	// Método 1: Verificar via PowerShell
	psScript := `
	$currentPrincipal = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
	$isAdmin = $currentPrincipal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
	if ($isAdmin) { exit 0 } else { exit 1 }
	`

	cmd := exec.Command("powershell.exe", "-NoProfile", "-Command", psScript)
	err := cmd.Run()
	return err == nil
}

// isWSLAvailable verifica se o WSL está disponível e configurado
func isWSLAvailable() bool {
	// Verificar se comando wsl existe
	cmd := exec.Command("wsl", "--version")
	err := cmd.Run()
	if err != nil {
		return false
	}

	// Verificar se há distribuições instaladas
	cmd = exec.Command("wsl", "--list", "--quiet")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	// Verificar se há pelo menos uma distribuição
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	return len(lines) > 0 && strings.TrimSpace(lines[0]) != ""
}

// checkPowerShellExecutionPolicy verifica e configura a política de execução
func checkPowerShellExecutionPolicy() error {
	// Verificar política atual
	cmd := exec.Command("powershell.exe", "-NoProfile", "-Command", "Get-ExecutionPolicy")
	output, err := cmd.Output()
	if err != nil {
		return &WindowsOnlyError{
			Code:        "POWERSHELL_ERROR",
			Message:     "Erro ao verificar política de execução do PowerShell",
			Suggestion:  "Verifique se o PowerShell está funcionando corretamente",
			ErrorType:   "powershell",
			Recoverable: true,
		}
	}

	policy := strings.TrimSpace(string(output))
	if policy == "Restricted" {
		return &WindowsOnlyError{
			Code:        "EXECUTION_POLICY_RESTRICTED",
			Message:     "Política de execução do PowerShell está restrita",
			Suggestion:  "Execute: Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser",
			ErrorType:   "permission",
			Recoverable: true,
		}
	}

	return nil
}

// checkRequiredTools verifica se as ferramentas necessárias estão disponíveis
func checkRequiredTools() error {
	tools := []string{
		"powershell.exe",
		"wsl.exe",
		"diskpart.exe",
	}

	for _, tool := range tools {
		cmd := exec.Command("where", tool)
		if err := cmd.Run(); err != nil {
			return &WindowsOnlyError{
				Code:        "TOOL_NOT_FOUND",
				Message:     fmt.Sprintf("Ferramenta necessária não encontrada: %s", tool),
				Suggestion:  fmt.Sprintf("Instale ou adicione %s ao PATH do sistema", tool),
				ErrorType:   "system",
				Recoverable: true,
			}
		}
	}

	return nil
}

// listWindowsOnlyDevices lista dispositivos USB específicos para Windows
func listWindowsOnlyDevices() ([]WindowsOnlyUSBDevice, error) {
	if err := validateWindowsEnvironment(); err != nil {
		return nil, err
	}

	// Script PowerShell para listar dispositivos com informações detalhadas
	psScript := `
	Get-Disk | Where-Object {
		$_.BusType -eq 'USB' -or 
		($_.BusType -eq 'SCSI' -and $_.Size -lt 500GB -and $_.Size -gt 1GB)
	} | ForEach-Object {
		$disk = $_
		$partitions = Get-Partition -DiskNumber $disk.Number -ErrorAction SilentlyContinue
		
		[PSCustomObject]@{
			DiskNumber = $disk.Number
			FriendlyName = $disk.FriendlyName
			Size = $disk.Size
			SizeFormatted = [math]::Round($disk.Size / 1GB, 2).ToString() + " GB"
			SerialNumber = $disk.SerialNumber
			BusType = $disk.BusType
			Model = $disk.Model
			IsSystem = $disk.IsSystem
			IsBoot = $disk.IsBoot
			IsOffline = $disk.IsOffline
			PartitionCount = if ($partitions) { $partitions.Count } else { 0 }
			Status = $disk.HealthStatus
		}
	} | ConvertTo-Json -Compress
	`

	cmd := exec.Command("powershell.exe", "-NoProfile", "-Command", psScript)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, &WindowsOnlyError{
			Code:        "DEVICE_LIST_ERROR",
			Message:     "Erro ao listar dispositivos USB",
			Suggestion:  "Verifique se há dispositivos USB conectados e se você tem permissões adequadas",
			ErrorType:   "device",
			Recoverable: true,
		}
	}

	// Parse JSON
	jsonStr := strings.TrimSpace(string(output))
	if jsonStr == "" {
		return []WindowsOnlyUSBDevice{}, nil
	}

	var devices []WindowsOnlyUSBDevice
	if strings.HasPrefix(jsonStr, "[") {
		if err := json.Unmarshal([]byte(jsonStr), &devices); err != nil {
			return nil, &WindowsOnlyError{
				Code:        "JSON_PARSE_ERROR",
				Message:     "Erro ao fazer parse da lista de dispositivos",
				Suggestion:  "Tente novamente ou execute o diagnóstico",
				ErrorType:   "system",
				Recoverable: true,
			}
		}
	} else {
		var device WindowsOnlyUSBDevice
		if err := json.Unmarshal([]byte(jsonStr), &device); err != nil {
			return nil, &WindowsOnlyError{
				Code:        "JSON_PARSE_ERROR",
				Message:     "Erro ao fazer parse do dispositivo",
				Suggestion:  "Tente novamente ou execute o diagnóstico",
				ErrorType:   "system",
				Recoverable: true,
			}
		}
		devices = []WindowsOnlyUSBDevice{device}
	}

	return devices, nil
}

// createWindowsOnlyUSB cria USB bootável específico para Windows
func createWindowsOnlyUSB(devicePath string, config *WindowsOnlyConfig) error {
	if err := validateWindowsEnvironment(); err != nil {
		return err
	}

	// Extrair número do disco
	diskNum, err := extractDiskNumber(devicePath)
	if err != nil {
		return err
	}

	// Validar dispositivo
	if err := validateWindowsOnlyDevice(diskNum); err != nil {
		return err
	}

	// Configurar diretórios
	if err := setupWindowsOnlyDirectories(config); err != nil {
		return err
	}

	// Gerenciar ISO
	isoPath, err := manageWindowsOnlyISO(config)
	if err != nil {
		return err
	}

	fmt.Printf("🚀 Criando USB Syntropy (Windows Only)\n")
	fmt.Printf("📍 Nó: %s\n", config.NodeName)
	fmt.Printf("💾 Dispositivo: %s (nº %d)\n", devicePath, diskNum)
	fmt.Printf("📀 ISO: %s\n", isoPath)
	fmt.Printf("📂 Diretório temporário: %s\n\n", config.TempDir)

	// Gerar chaves SSH
	fmt.Println("🔑 Gerando chaves SSH...")
	sshPrivateKey, sshPublicKey, err := generateSSHKeyPair(config.NodeName)
	if err != nil {
		return fmt.Errorf("erro ao gerar chaves SSH: %w", err)
	}
	fmt.Println("✅ Chaves SSH geradas com sucesso")

	// Gerar certificados TLS
	fmt.Println("🔐 Gerando certificados TLS...")
	certs, err := generateCertificates(config.NodeName, config.OwnerKeyFile)
	if err != nil {
		return fmt.Errorf("erro ao gerar certificados: %w", err)
	}
	fmt.Println("✅ Certificados TLS gerados com sucesso")

	// Salvar certificados
	certsDir := filepath.Join(config.TempDir, "certs")
	os.MkdirAll(certsDir, 0755)
	_, _, _, _, err = saveCertificates(certs, certsDir)
	if err != nil {
		return fmt.Errorf("erro ao salvar certificados: %w", err)
	}

	// Criar configuração cloud-init
	fmt.Println("📝 Criando configuração cloud-init...")
	cloudInitDir := filepath.Join(config.TempDir, "cloud-init")
	os.MkdirAll(cloudInitDir, 0755)

	// Configuração para cloud-init
	configForCloudInit := &Config{
		NodeName:        config.NodeName,
		NodeDescription: config.NodeDescription,
		Coordinates:     config.Coordinates,
		OwnerKeyFile:    config.OwnerKeyFile,
		Label:           config.Label,
		ISOPath:         isoPath,
		DiscoveryServer: config.DiscoveryServer,
		SSHPublicKey:    sshPublicKey,
		SSHPrivateKey:   sshPrivateKey,
		CreatedBy:       config.CreatedBy,
	}

	cloudInitConfig, err := generateCloudInitConfig(configForCloudInit, config.TempDir, certs)
	if err != nil {
		return fmt.Errorf("erro ao gerar configuração cloud-init: %w", err)
	}

	certPaths := map[string]string{
		"CAKey":    filepath.Join(certsDir, "ca.key"),
		"CACert":   filepath.Join(certsDir, "ca.crt"),
		"NodeKey":  filepath.Join(certsDir, "node.key"),
		"NodeCert": filepath.Join(certsDir, "node.crt"),
	}

	if err := createCloudInitFiles(cloudInitConfig, config.TempDir, certPaths); err != nil {
		return fmt.Errorf("erro ao criar arquivos cloud-init: %w", err)
	}

	// Copiar scripts
	if err := copyScripts(config.TempDir); err != nil {
		return fmt.Errorf("erro ao copiar scripts: %w", err)
	}

	fmt.Println("✅ Configuração preparada com sucesso")

	// Executar criação do USB
	return executeWindowsOnlyUSBCreation(diskNum, isoPath, config)
}

// extractDiskNumber extrai o número do disco do caminho do dispositivo
func extractDiskNumber(devicePath string) (int, error) {
	var diskNum int
	switch {
	case strings.HasPrefix(devicePath, "PHYSICALDRIVE"):
		_, err := fmt.Sscanf(devicePath, "PHYSICALDRIVE%d", &diskNum)
		if err != nil {
			return 0, &WindowsOnlyError{
				Code:        "INVALID_DEVICE_FORMAT",
				Message:     "Formato de dispositivo inválido",
				Suggestion:  "Use o formato PHYSICALDRIVEn (ex: PHYSICALDRIVE1)",
				ErrorType:   "device",
				Recoverable: true,
			}
		}
	case strings.HasPrefix(devicePath, "\\\\.\\PHYSICALDRIVE"):
		_, err := fmt.Sscanf(devicePath, "\\\\.\\PHYSICALDRIVE%d", &diskNum)
		if err != nil {
			return 0, &WindowsOnlyError{
				Code:        "INVALID_DEVICE_FORMAT",
				Message:     "Formato de dispositivo inválido",
				Suggestion:  "Use o formato \\\\.\\PHYSICALDRIVEn (ex: \\\\.\\PHYSICALDRIVE1)",
				ErrorType:   "device",
				Recoverable: true,
			}
		}
	default:
		return 0, &WindowsOnlyError{
			Code:        "UNSUPPORTED_DEVICE_FORMAT",
			Message:     "Formato de dispositivo não suportado",
			Suggestion:  "Use PHYSICALDRIVEn ou \\\\.\\PHYSICALDRIVEn",
			ErrorType:   "device",
			Recoverable: true,
		}
	}

	return diskNum, nil
}

// validateWindowsOnlyDevice valida um dispositivo específico para Windows
func validateWindowsOnlyDevice(diskNum int) error {
	// Script PowerShell para validação detalhada
	psScript := fmt.Sprintf(`
	$disk = Get-Disk -Number %d -ErrorAction SilentlyContinue
	if (-not $disk) {
		Write-Output "DISK_NOT_FOUND"
		exit 1
	}
	
	if ($disk.IsSystem) {
		Write-Output "SYSTEM_DISK"
		exit 2
	}
	
	if ($disk.IsBoot) {
		Write-Output "BOOT_DISK"
		exit 3
	}
	
	if ($disk.Size -lt 1GB) {
		Write-Output "TOO_SMALL"
		exit 4
	}
	
	if ($disk.Size -gt 2TB) {
		Write-Output "TOO_LARGE"
		exit 5
	}
	
	$partitions = Get-Partition -DiskNumber %d -ErrorAction SilentlyContinue
	if ($partitions) {
		foreach ($part in $partitions) {
			if ($part.DriveLetter -eq "C" -or $part.IsSystem -or $part.IsBoot) {
				Write-Output "SYSTEM_PARTITION"
				exit 6
			}
		}
	}
	
	Write-Output "VALID"
	exit 0
	`, diskNum, diskNum)

	cmd := exec.Command("powershell.exe", "-NoProfile", "-Command", psScript)
	output, err := cmd.CombinedOutput()

	if err != nil {
		result := strings.TrimSpace(string(output))
		switch result {
		case "DISK_NOT_FOUND":
			return &WindowsOnlyError{
				Code:        "DISK_NOT_FOUND",
				Message:     fmt.Sprintf("Dispositivo %d não encontrado", diskNum),
				Suggestion:  "Verifique se o USB está conectado e tente novamente",
				ErrorType:   "device",
				Recoverable: true,
			}
		case "SYSTEM_DISK":
			return &WindowsOnlyError{
				Code:        "SYSTEM_DISK",
				Message:     "Dispositivo é um disco do sistema",
				Suggestion:  "Use um dispositivo USB removível, não o disco do sistema",
				ErrorType:   "device",
				Recoverable: true,
			}
		case "BOOT_DISK":
			return &WindowsOnlyError{
				Code:        "BOOT_DISK",
				Message:     "Dispositivo é um disco de boot",
				Suggestion:  "Use um dispositivo USB removível, não o disco de boot",
				ErrorType:   "device",
				Recoverable: true,
			}
		case "TOO_SMALL":
			return &WindowsOnlyError{
				Code:        "TOO_SMALL",
				Message:     "Dispositivo é muito pequeno",
				Suggestion:  "Use um dispositivo USB com pelo menos 2GB",
				ErrorType:   "device",
				Recoverable: true,
			}
		case "TOO_LARGE":
			return &WindowsOnlyError{
				Code:        "TOO_LARGE",
				Message:     "Dispositivo é muito grande",
				Suggestion:  "Use um dispositivo USB com no máximo 2TB",
				ErrorType:   "device",
				Recoverable: true,
			}
		case "SYSTEM_PARTITION":
			return &WindowsOnlyError{
				Code:        "SYSTEM_PARTITION",
				Message:     "Dispositivo contém partições do sistema",
				Suggestion:  "Use um dispositivo USB sem partições do sistema",
				ErrorType:   "device",
				Recoverable: true,
			}
		default:
			return &WindowsOnlyError{
				Code:        "VALIDATION_ERROR",
				Message:     "Erro na validação do dispositivo",
				Suggestion:  "Execute o diagnóstico para mais informações",
				ErrorType:   "device",
				Recoverable: true,
			}
		}
	}

	return nil
}

// setupWindowsOnlyDirectories configura os diretórios necessários
func setupWindowsOnlyDirectories(config *WindowsOnlyConfig) error {
	if config.TempDir == "" {
		tempDir := os.TempDir()
		config.TempDir = filepath.Join(tempDir, "syntropy-usb", time.Now().Format("20060102-150405"))
	}

	// Criar diretórios necessários
	dirs := []string{
		config.TempDir,
		filepath.Join(config.TempDir, "certs"),
		filepath.Join(config.TempDir, "cloud-init"),
		filepath.Join(config.TempDir, "scripts"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return &WindowsOnlyError{
				Code:        "DIRECTORY_CREATE_ERROR",
				Message:     fmt.Sprintf("Erro ao criar diretório: %s", dir),
				Suggestion:  "Verifique permissões e espaço em disco",
				ErrorType:   "system",
				Recoverable: true,
			}
		}
	}

	return nil
}

// manageWindowsOnlyISO gerencia o download e cache da ISO Ubuntu
func manageWindowsOnlyISO(config *WindowsOnlyConfig) (string, error) {
	if config.ISOPath != "" {
		// Verificar se arquivo ISO existe
		if _, err := os.Stat(config.ISOPath); err != nil {
			return "", &WindowsOnlyError{
				Code:        "ISO_NOT_FOUND",
				Message:     fmt.Sprintf("Arquivo ISO não encontrado: %s", config.ISOPath),
				Suggestion:  "Verifique o caminho do arquivo ISO",
				ErrorType:   "system",
				Recoverable: true,
			}
		}
		return config.ISOPath, nil
	}

	// Usar função existente para gerenciar cache
	cacheDir := filepath.Join(os.Getenv("USERPROFILE"), ".syntropy", "cache")
	os.MkdirAll(cacheDir, 0755)

	isoPath, err := manageISOCache(cacheDir)
	if err != nil {
		return "", &WindowsOnlyError{
			Code:        "ISO_DOWNLOAD_ERROR",
			Message:     "Erro ao baixar ISO Ubuntu",
			Suggestion:  "Verifique sua conexão com a internet e tente novamente",
			ErrorType:   "system",
			Recoverable: true,
		}
	}

	return isoPath, nil
}

// executeWindowsOnlyUSBCreation executa a criação do USB específica para Windows
func executeWindowsOnlyUSBCreation(diskNum int, isoPath string, config *WindowsOnlyConfig) error {
	// Criar script PowerShell robusto para criação do USB
	psScript := fmt.Sprintf(`
	$ErrorActionPreference = "Stop"
	$ProgressPreference = "SilentlyContinue"
	
	Write-Host "[INFO] Iniciando criacao de USB Syntropy (Windows Only)" -ForegroundColor Green
	Write-Host "[NODE] No: %s" -ForegroundColor Cyan
	Write-Host "[DISK] Dispositivo: PHYSICALDRIVE%d" -ForegroundColor Cyan
	Write-Host "[ISO] ISO: %s" -ForegroundColor Cyan
	Write-Host ""
	
	try {
		# Verificar se dispositivo ainda existe
		$disk = Get-Disk -Number %d -ErrorAction SilentlyContinue
		if (-not $disk) {
			throw "Dispositivo %d nao encontrado. Verifique se o USB esta conectado."
		}
		
		Write-Host "[OK] Dispositivo verificado: $($disk.FriendlyName)" -ForegroundColor Green
		
		# Verificar se ISO existe
		if (-not (Test-Path "%s")) {
			throw "Arquivo ISO nao encontrado: %s"
		}
		
		Write-Host "[OK] ISO verificada: $(Get-Item \"%s\").Length bytes" -ForegroundColor Green
		
		# Colocar disco offline
		Write-Host "[OFFLINE] Colocando disco offline..." -ForegroundColor Yellow
		Set-Disk -Number %d -IsReadOnly $false -IsOffline $true
		
		# Montar no WSL
		Write-Host "[MOUNT] Montando disco no WSL..." -ForegroundColor Yellow
		$mountResult = wsl --mount PHYSICALDRIVE%d --bare 2>&1
		if ($LASTEXITCODE -ne 0) {
			throw "Falha ao montar disco no WSL: $mountResult"
		}
		
		Write-Host "[OK] Disco montado no WSL com sucesso" -ForegroundColor Green
		
		# Executar script de criacao no WSL
		Write-Host "[WSL] Executando criacao no WSL..." -ForegroundColor Yellow
		$wslScript = @"
#!/bin/bash
set -euo pipefail

echo "[SEARCH] Detectando dispositivo WSL..."

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
  echo "[ERROR] ERRO: Nao foi possivel detectar o dispositivo no WSL" >&2
  echo "Dispositivos disponiveis:" >&2
  ls -la /dev/sd* /dev/hd* /dev/nvme* 2>/dev/null || true >&2
  exit 1
fi

echo "[OK] Dispositivo WSL detectado: $dev"

# Verificar se ISO existe
ISO="%s"
if [ ! -f "$ISO" ]; then
  echo "[ERROR] ERRO: ISO nao encontrada: $ISO" >&2
  exit 1
fi

echo "[ISO] Gravando ISO: $ISO -> $dev"
sudo dd if="$ISO" of="$dev" bs=4M status=progress conv=fsync
sync

echo "Aguardando gravacao finalizar..."
sleep 3

echo "[CONFIG] Criando particao CIDATA..."
sudo sgdisk -e "$dev"
sudo sgdisk -n 0:0:+128MiB -t 0:0700 -c 0:CIDATA "$dev"
sleep 2

# Determinar nome da particao CIDATA
cidata_part=""
if [[ "$dev" =~ nvme ]]; then
  cidata_part="${dev}p2"
else
  cidata_part="${dev}2"
fi

echo "[PARTITION] Particao CIDATA: $cidata_part"

# Verificar se particao existe
if [ ! -b "$cidata_part" ]; then
  echo "[ERROR] ERRO: Particao CIDATA nao encontrada: $cidata_part" >&2
  echo "Particoes disponiveis:" >&2
  ls -la ${dev}* 2>/dev/null || true >&2
  exit 1
fi

echo "[DISK] Formatando particao CIDATA..."
sudo mkfs.vfat -F 32 -n CIDATA "$cidata_part"

echo "[MOUNT] Montando particao CIDATA..."
mount_point="$HOME/.syntropy/work/cidata-mount"
sudo mkdir -p "$mount_point"
sudo mount "$cidata_part" "$mount_point"

echo "[COPY] Copiando arquivos cloud-init..."
cloud_init_dir="%s/cloud-init"

# Verificar se diretorio cloud-init existe
if [ ! -d "$cloud_init_dir" ]; then
  echo "[ERROR] ERRO: Diretorio cloud-init nao encontrado: $cloud_init_dir" >&2
  sudo umount "$mount_point" || true
  sudo rmdir "$mount_point" || true
  exit 1
fi

# Verificar se arquivos existem
for file in user-data meta-data network-config; do
  if [ ! -f "$cloud_init_dir/$file" ]; then
    echo "[ERROR] ERRO: Arquivo cloud-init nao encontrado: $cloud_init_dir/$file" >&2
    sudo umount "$mount_point" || true
    sudo rmdir "$mount_point" || true
    exit 1
  fi
done

sudo cp "$cloud_init_dir/user-data" "$mount_point/"
sudo cp "$cloud_init_dir/meta-data" "$mount_point/"
sudo cp "$cloud_init_dir/network-config" "$mount_point/"

echo "[SEARCH] Verificando arquivos copiados..."
ls -la "$mount_point/"

echo "[UNMOUNT] Desmontando particao..."
sudo umount "$mount_point"
sudo rmdir "$mount_point"
sync

echo "[OK] USB criado com sucesso usando estrategia NoCloud!"
"@

		# Converter caminho ISO para WSL
		$isoWSL = wsl wslpath -u "%s"
		$workDirWSL = wsl wslpath -u "%s"
		
		# Substituir caminhos no script
		$wslScript = $wslScript -replace "%s", $isoWSL
		$wslScript = $wslScript -replace "%s/cloud-init", "$workDirWSL/cloud-init"
		
		# Executar script no WSL
		$wslResult = wsl bash -lc $wslScript 2>&1
		$exitCode = $LASTEXITCODE
		
		if ($exitCode -ne 0) {
			Write-Host "[ERROR] ERRO no WSL:" -ForegroundColor Red
			Write-Host $wslResult -ForegroundColor Red
			throw "Script WSL falhou com codigo: $exitCode"
		}
		
		Write-Host "[OK] Script WSL executado com sucesso!" -ForegroundColor Green
		Write-Host $wslResult -ForegroundColor White
		
	} catch {
		Write-Host "[ERROR] ERRO: $($_.Exception.Message)" -ForegroundColor Red
		throw
	} finally {
		Write-Host "[CLEANUP] Limpando recursos..." -ForegroundColor Yellow
		try { 
			wsl --unmount PHYSICALDRIVE%d 2>$null
			Write-Host "[OK] Dispositivo desmontado do WSL" -ForegroundColor Green
		} catch { 
			Write-Host "[WARNING] Aviso: Falha ao desmontar do WSL (pode ja estar desmontado)" -ForegroundColor Yellow 
		}
		
		try {
			Set-Disk -Number %d -IsOffline $false
			Write-Host "[OK] Disco voltou online no Windows" -ForegroundColor Green
		} catch {
			Write-Host "[WARNING] Aviso: Falha ao voltar disco online" -ForegroundColor Yellow
		}
	}
	
	Write-Host "[SUCCESS] USB criado com sucesso usando estrategia NoCloud!" -ForegroundColor Green
	Write-Host "[CONFIG] O USB agora contem:" -ForegroundColor Cyan
	Write-Host "   - ISO Ubuntu original (bootavel)" -ForegroundColor White
	Write-Host "   - Particao CIDATA com configuracao cloud-init" -ForegroundColor White
	Write-Host "   - Configuracao sera aplicada automaticamente no boot" -ForegroundColor White
	Write-Host ""
	Write-Host "[INFO] Informacoes do no:" -ForegroundColor Cyan
	Write-Host "   - Nome: %s" -ForegroundColor White
	Write-Host "   - Descricao: %s" -ForegroundColor White
	Write-Host "   - Criado por: %s" -ForegroundColor White
	Write-Host "   - Data: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')" -ForegroundColor White
	`, config.NodeName, diskNum, isoPath, diskNum, diskNum, isoPath, isoPath, diskNum, diskNum, diskNum, diskNum, isoPath, config.TempDir, isoPath, config.TempDir, diskNum, diskNum, config.NodeName, config.NodeDescription, config.CreatedBy)

	// Salvar script PowerShell
	scriptPath := filepath.Join(config.TempDir, "create_usb_windows_only.ps1")
	// Limpar script PowerShell para evitar problemas de encoding
	cleanScript := cleanPowerShellString(psScript)
	if err := os.WriteFile(scriptPath, []byte(cleanScript), 0644); err != nil {
		return fmt.Errorf("erro ao criar script PowerShell: %w", err)
	}

	// Executar script no PowerShell atual
	fmt.Println("📝 Executando script PowerShell...")
	cmd := exec.Command("powershell.exe", "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return &WindowsOnlyError{
			Code:        "USB_CREATION_FAILED",
			Message:     "Falha na criação do USB",
			Suggestion:  "Verifique os logs de erro e execute o diagnóstico",
			ErrorType:   "device",
			Recoverable: true,
		}
	}

	return nil
}

// formatWindowsOnlyUSB formata um dispositivo USB específico para Windows
func formatWindowsOnlyUSB(devicePath, label string) error {
	if err := validateWindowsEnvironment(); err != nil {
		return err
	}

	diskNum, err := extractDiskNumber(devicePath)
	if err != nil {
		return err
	}

	// Script PowerShell para formatação
	psScript := fmt.Sprintf(`
	$ErrorActionPreference = "Stop"
	
	Write-Host "🔧 Formatando dispositivo USB (Windows Only)" -ForegroundColor Yellow
	Write-Host "💾 Dispositivo: %s (nº %d)" -ForegroundColor Cyan
	Write-Host "🏷️  Rótulo: %s" -ForegroundColor Cyan
	Write-Host ""
	
	try {
		# Verificar se dispositivo existe
		$disk = Get-Disk -Number %d -ErrorAction SilentlyContinue
		if (-not $disk) {
			throw "Dispositivo %d não encontrado"
		}
		
		Write-Host "✅ Dispositivo encontrado: $($disk.FriendlyName)" -ForegroundColor Green
		
		# Limpar disco
		Write-Host "🧹 Limpando disco..." -ForegroundColor Yellow
		Clear-Disk -Number %d -RemoveData -Confirm:$false
		
		# Criar nova partição
		Write-Host "📁 Criando nova partição..." -ForegroundColor Yellow
		New-Partition -DiskNumber %d -UseMaximumSize -AssignDriveLetter | Out-Null
		
		# Formatar volume
		Write-Host "💾 Formatando com FAT32..." -ForegroundColor Yellow
		Format-Volume -DriveLetter (Get-Partition -DiskNumber %d | Select-Object -First 1).DriveLetter -FileSystem FAT32 -NewFileSystemLabel "%s" -Confirm:$false
		
		Write-Host "✅ Formatação concluída com sucesso!" -ForegroundColor Green
		
	} catch {
		Write-Host "❌ ERRO: $($_.Exception.Message)" -ForegroundColor Red
		throw
	}
	`, devicePath, diskNum, label, diskNum, diskNum, diskNum, diskNum, diskNum, diskNum, label)

	// Salvar e executar script
	tempDir := filepath.Join(os.TempDir(), "syntropy-format")
	os.MkdirAll(tempDir, 0755)

	scriptPath := filepath.Join(tempDir, "format_usb.ps1")
	// Limpar script PowerShell para evitar problemas de encoding
	cleanScript := cleanPowerShellString(psScript)
	if err := os.WriteFile(scriptPath, []byte(cleanScript), 0644); err != nil {
		return fmt.Errorf("erro ao criar script de formatação: %w", err)
	}

	cmd := exec.Command("powershell.exe", "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return &WindowsOnlyError{
			Code:        "FORMAT_FAILED",
			Message:     "Falha na formatação do dispositivo",
			Suggestion:  "Verifique se o dispositivo não está em uso e tente novamente",
			ErrorType:   "device",
			Recoverable: true,
		}
	}

	return nil
}

// debugWindowsOnlyEnvironment executa diagnóstico específico para Windows
func debugWindowsOnlyEnvironment() error {
	fmt.Println("🔍 Executando diagnóstico Windows Only...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Verificar ambiente básico
	fmt.Println("🖥️  Verificando ambiente Windows...")
	if runtime.GOOS != "windows" {
		fmt.Println("❌ Não está executando no Windows")
		return fmt.Errorf("ambiente não é Windows")
	}
	fmt.Println("✅ Executando no Windows")

	// Verificar privilégios
	fmt.Println("🔐 Verificando privilégios...")
	if !isRunningAsAdministrator() {
		fmt.Println("❌ Não está executando como Administrador")
		fmt.Println("💡 Solução: Execute o PowerShell como Administrador")
		return fmt.Errorf("privilégios de administrador necessários")
	}
	fmt.Println("✅ Executando como Administrador")

	// Verificar WSL
	fmt.Println("🐧 Verificando WSL...")
	if !isWSLAvailable() {
		fmt.Println("❌ WSL não está disponível")
		fmt.Println("💡 Solução: Execute 'wsl --install' ou 'wsl --install -d Ubuntu'")
		return fmt.Errorf("WSL não disponível")
	}
	fmt.Println("✅ WSL está disponível")

	// Verificar política de execução
	fmt.Println("⚙️  Verificando política de execução do PowerShell...")
	if err := checkPowerShellExecutionPolicy(); err != nil {
		fmt.Printf("❌ %v\n", err)
		return err
	}
	fmt.Println("✅ Política de execução OK")

	// Verificar ferramentas
	fmt.Println("🛠️  Verificando ferramentas necessárias...")
	if err := checkRequiredTools(); err != nil {
		fmt.Printf("❌ %v\n", err)
		return err
	}
	fmt.Println("✅ Todas as ferramentas estão disponíveis")

	// Listar dispositivos
	fmt.Println("💾 Verificando dispositivos USB...")
	devices, err := listWindowsOnlyDevices()
	if err != nil {
		fmt.Printf("❌ Erro ao listar dispositivos: %v\n", err)
		return err
	}

	if len(devices) == 0 {
		fmt.Println("⚠️  Nenhum dispositivo USB encontrado")
		fmt.Println("💡 Solução: Conecte um dispositivo USB e tente novamente")
	} else {
		fmt.Printf("✅ %d dispositivo(s) USB encontrado(s):\n", len(devices))
		for _, device := range devices {
			fmt.Printf("   • %s - %s (%s)\n",
				fmt.Sprintf("PHYSICALDRIVE%d", device.DiskNumber),
				device.FriendlyName,
				device.SizeFormatted)
		}
	}

	fmt.Println("\n🎉 Diagnóstico concluído com sucesso!")
	return nil
}
