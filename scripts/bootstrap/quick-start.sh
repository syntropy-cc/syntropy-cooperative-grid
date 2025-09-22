#!/bin/bash

# Syntropy Cooperative Grid - Quick Start Script
# Version: 1.0.0

set -e

# Colors
PURPLE='\033[0;35m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Show banner
echo -e "${PURPLE}"
cat << 'EOF'
╔════════════════════════════════════════════════════════════════════════════╗
║                    SYNTROPY COOPERATIVE GRID                              ║
║                         Quick Start Setup                                 ║
║                                                                           ║
║  This script will guide you through setting up your first Syntropy node  ║
╚════════════════════════════════════════════════════════════════════════════╝
EOF
echo -e "${NC}"

# Function to check step completion
check_step() {
    local step=$1
    local description=$2
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Step $step: $description - Completed${NC}"
        return 0
    else
        echo -e "${RED}✗ Step $step: $description - Failed${NC}"
        return 1
    fi
}

# Function to pause and confirm
confirm_continue() {
    echo
    read -p "Press Enter to continue or Ctrl+C to abort..."
    echo
}

# Step 1: Install prerequisites
echo -e "${BLUE}[Step 1/5] Installing prerequisites...${NC}"
if [ -f "$SCRIPT_DIR/install-prerequisites.sh" ]; then
    bash "$SCRIPT_DIR/install-prerequisites.sh"
    check_step 1 "Prerequisites installation" || exit 1
else
    echo -e "${RED}Prerequisites script not found at: $SCRIPT_DIR/install-prerequisites.sh${NC}"
    exit 1
fi

confirm_continue

# Step 2: Setup management environment
echo -e "${BLUE}[Step 2/5] Setting up management environment...${NC}"
if [ -f "$SCRIPT_DIR/setup-syntropy-management.sh" ]; then
    bash "$SCRIPT_DIR/setup-syntropy-management.sh"
    check_step 2 "Management environment setup" || exit 1
else
    echo -e "${RED}Management setup script not found!${NC}"
    exit 1
fi

confirm_continue

# Step 3: Check for USB devices
# Step 3: Check for USB devices
echo -e "${BLUE}[Step 3/5] Checking for USB devices...${NC}"
echo "Available USB devices:"

# First source required modules
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/lib/colors.sh"
source "$SCRIPT_DIR/lib/logging.sh"

# Function to detect WSL
is_wsl() {
    if grep -qi microsoft /proc/version; then
        return 0  # True, is WSL
    else
        return 1  # False, not WSL
    fi
}

# Function to get USB devices in WSL
get_wsl_usb_devices() {
    # Basic PowerShell command that we know works
    local ps_command='
    $ErrorActionPreference = "Stop"
    try {
        $usbDisks = Get-PnpDevice -Class "DiskDrive" | Where-Object { $_.FriendlyName -like "*USB*" }
        foreach ($disk in $usbDisks) {
            # Basic disk info
            Write-Host ($disk.FriendlyName + "|" + $disk.Status)
        }
    } catch {
        Write-Error "Failed to get USB devices: $_"
        exit 1
    }'
    
    # Run PowerShell command with better error handling
    local result
    result=$(powershell.exe -Command "$ps_command" 2>&1)
    local exit_code=$?
    
    if [ $exit_code -ne 0 ]; then
        echo "PowerShell Error: $result" >&2
        return 1
    fi
    
    # Check if we got any output
    if [ -z "$result" ]; then
        echo "No USB devices found" >&2
        return 1
    fi
    
    echo "$result"
    return 0
}

if is_wsl; then
    echo -e "${YELLOW}WSL detected - using Windows USB detection${NC}"
    
    # Get USB devices using PowerShell
    echo -e "\nQuerying Windows for USB devices..."
    if ! devices=$(get_wsl_usb_devices); then
        echo "[ERROR] Failed to detect USB devices"
        exit 1
    fi
    
    echo -e "\nFound USB devices:"
    echo "$devices" | while IFS='|' read -r name status; do
        echo "Name: $name"
        echo "Status: $status"
        echo "-------------------"
    done
else
    # Standard Linux USB detection
    lsblk -d -o NAME,SIZE,TYPE,RM | grep -E "disk.*1$" || echo "No removable devices found"
fi

# Validate USB devices found
if is_wsl; then
    # Para WSL, já validamos na detecção anterior
    # Se chegamos até aqui, significa que encontramos dispositivos
    WSL_USB_DETECTED=true
else
    # Para Linux normal, usamos lsblk
    if ! lsblk -d -o NAME,RM | grep -q "1$"; then
        echo -e "${RED}No USB devices detected. Please insert a USB drive and try again.${NC}"
        exit 1
    fi
fi

confirm_continue

# Step 4: Node creation preparation
echo -e "${BLUE}[Step 4/5] Preparing for node creation${NC}"

# Function to validate node name
validate_node_name() {
    local name=$1
    if [[ ! "$name" =~ ^[a-zA-Z0-9-]+$ ]]; then
        echo -e "${RED}Invalid node name. Use only letters, numbers, and hyphens.${NC}"
        return 1
    fi
    return 0
}

# Get and validate node name with retry
while true; do
    echo
    echo -e "${YELLOW}Please choose your node name:${NC}"
    read -p "Enter node name (e.g., syntropy-node-01): " NODE_NAME
    
    if validate_node_name "$NODE_NAME"; then
        break
    fi
    echo "Please try again."
done

# Function to get USB devices with drive letters in WSL
get_wsl_usb_devices_with_letters() {
    local ps_command='
    $ErrorActionPreference = "Stop"
    try {
        $usbDisks = Get-PnpDevice -Class "DiskDrive" | Where-Object { $_.FriendlyName -like "*USB*" }
        foreach ($disk in $usbDisks) {
            $diskNumber = ($disk.InstanceId -split "\\",-1)[2] -replace ".*#",""
            
            # Get associated drive letter
            $volumes = Get-WmiObject Win32_LogicalDisk | Where-Object { $_.DriveType -eq 2 }
            foreach ($volume in $volumes) {
                Write-Host ($disk.FriendlyName + "|" + $volume.DeviceID.Replace(":", "") + "|" + $disk.Status)
            }
        }
    } catch {
        Write-Error "Failed to get USB devices: $_"
        exit 1
    }'
    
    powershell.exe -Command "$ps_command" 2>/dev/null
}

# Function to list available USB devices
list_usb_devices() {
    echo -e "${YELLOW}Available USB devices:${NC}"
    if is_wsl; then
        local devices
        devices=$(get_wsl_usb_devices_with_letters)
        local i=1
        echo "$devices" | while IFS='|' read -r name letter status; do
            echo "[$i] Name: $name"
            echo "    Drive: ${letter}:"
            echo "    Status: $status"
            echo "-------------------"
            i=$((i + 1))
        done
    else
        local i=1
        while IFS= read -r line; do
            echo "[$i] $line"
            echo "-------------------"
            i=$((i + 1))
        done < <(lsblk -d -o NAME,SIZE,TYPE,RM,MODEL | grep -E "disk.*1" || echo "No removable devices found")
    fi
    echo
}

# Function to validate USB device
validate_usb_device() {
    local device=$1
    if is_wsl; then
        # Primeiro, vamos verificar se temos o device montado
        if ! mountpoint -q "/mnt/$device"; then
            echo -e "${RED}Drive ${device}: is not mounted in WSL. Please mount it first.${NC}"
            echo "Try running: 'sudo mkdir -p /mnt/$device && sudo mount -t drvfs ${device}: /mnt/$device'"
            return 1
        fi
        
        # Verificar se temos permissão de acesso
        if ! [[ -r "/mnt/$device" && -w "/mnt/$device" ]]; then
            echo -e "${RED}No permission to access /mnt/$device. Try remounting with correct permissions.${NC}"
            return 1
        fi
        
        # Verificar se é um dispositivo removível usando PowerShell
        local ps_command='
        $drive = Get-WmiObject Win32_LogicalDisk | Where-Object { $_.DeviceID -eq "'${device}':"}
        if ($drive.DriveType -eq 2) { exit 0 } else { exit 1 }'
        
        if ! powershell.exe -Command "$ps_command" 2>/dev/null; then
            echo -e "${RED}Drive ${device}: is not a removable device.${NC}"
            return 1
        fi
    else
        if [[ ! -b "/dev/$device" ]]; then
            echo -e "${RED}Invalid USB device: /dev/$device${NC}"
            return 1
        fi
    fi
    return 0
}

# Function to mount WSL drive if needed
mount_wsl_drive() {
    local drive=$1
    if ! mountpoint -q "/mnt/$drive"; then
        echo -e "${YELLOW}Attempting to mount $drive: to /mnt/$drive...${NC}"
        sudo mkdir -p "/mnt/$drive" 2>/dev/null
        if ! sudo mount -t drvfs "$drive:" "/mnt/$drive" 2>/dev/null; then
            echo -e "${RED}Failed to mount $drive: automatically${NC}"
            return 1
        fi
        echo -e "${GREEN}Successfully mounted $drive: to /mnt/$drive${NC}"
    fi
    return 0
}

# Function to select USB device from menu
select_usb_device() {
    if is_wsl; then
        local devices
        devices=$(get_wsl_usb_devices_with_letters)
        local count=0
        local options=()
        local letters=()
        
        # Criar arrays com as opções e letras dos drives
        while IFS='|' read -r name letter status; do
            count=$((count + 1))
            options+=("$name")
            letters+=("$letter")
        done < <(echo "$devices")
        
        # Se não houver dispositivos, retornar erro
        if [ $count -eq 0 ]; then
            echo -e "${RED}No USB devices found${NC}"
            return 1
        fi
        
        # Mostrar menu e pegar seleção
        while true; do
            list_usb_devices
            echo -e "${YELLOW}Select a USB device (1-$count):${NC}"
            read -p "> " selection
            
            if [[ "$selection" =~ ^[0-9]+$ ]] && [ "$selection" -ge 1 ] && [ "$selection" -le $count ]; then
                local index=$((selection - 1))
                SELECTED_DRIVE="${letters[$index]}"
                break
            else
                echo -e "${RED}Invalid selection. Please choose a number between 1 and $count${NC}"
            fi
        done
        return 0
    else
        local devices
        mapfile -t devices < <(lsblk -d -o NAME,SIZE,TYPE,RM,MODEL | grep -E "disk.*1")
        local count=${#devices[@]}
        
        if [ $count -eq 0 ]; then
            echo -e "${RED}No USB devices found${NC}"
            return 1
        fi
        
        while true; do
            list_usb_devices
            echo -e "${YELLOW}Select a USB device (1-$count):${NC}"
            read -p "> " selection
            
            if [[ "$selection" =~ ^[0-9]+$ ]] && [ "$selection" -ge 1 ] && [ "$selection" -le $count ]; then
                local selected_line=${devices[$((selection - 1))]}
                SELECTED_DRIVE=$(echo "$selected_line" | awk '{print $1}')
                break
            else
                echo -e "${RED}Invalid selection. Please choose a number between 1 and $count${NC}"
            fi
        done
        return 0
    fi
}

# Function to get USB device with validation and retry
get_usb_device() {
    while true; do
        if ! select_usb_device; then
            echo "No USB devices available. Please insert a USB drive and try again."
            read -p "Press Enter to retry or Ctrl+C to abort..."
            continue
        fi
        
        if is_wsl; then
            DRIVE_LETTER=$(echo "$SELECTED_DRIVE" | tr '[:upper:]' '[:lower:]')
            
            # Tentar montar o drive se necessário
            if ! mount_wsl_drive "$DRIVE_LETTER"; then
                echo "Please try again or choose a different drive."
                continue
            fi
            
            if validate_usb_device "$DRIVE_LETTER"; then
                USB_DEVICE="/mnt/$DRIVE_LETTER"
                break
            fi
        else
            if validate_usb_device "$SELECTED_DRIVE"; then
                USB_DEVICE="/dev/$SELECTED_DRIVE"
                break
            fi
        fi
        echo "Please try again."
    done
}

# Function to confirm USB device with retry
confirm_usb_device() {
    while true; do
        echo -e "${RED}WARNING: This will erase all data on $USB_DEVICE${NC}"
        read -p "Are you sure you want to continue? (y/N): " -n 1 -r
        echo
        
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            return 0
        elif [[ $REPLY =~ ^[Nn]$ ]] || [[ -z $REPLY ]]; then
            echo "Operation cancelled by user. Please select a different device."
            return 1
        fi
        echo "Please answer 'y' or 'n'."
    done
}

# Get and confirm USB device
while true; do
    get_usb_device
    
    if confirm_usb_device; then
        break
    fi
done

# Create USB
echo -e "${BLUE}Creating Syntropy node USB...${NC}"

# No WSL, precisamos obter o dispositivo físico equivalente
if is_wsl; then
    # Função para obter o dispositivo físico do Windows
    get_physical_device() {
        local drive_letter=$1
        echo -e "${YELLOW}Detecting physical device for drive ${drive_letter}:${NC}" >&2
        
        # Criar um arquivo temporário para o script PowerShell
        local ps_script=$(mktemp)
        
        # Escrever o script PowerShell com a variável corretamente interpolada
        cat > "$ps_script" << POWERSHELL
# Configurar tratamento de erro
\$ErrorActionPreference = "Stop"
Write-Warning "Starting device detection for drive ${drive_letter}:"

try {
    # Primeiro, mostrar todos os discos USB disponíveis
    Write-Warning "Available USB devices:"
    Get-CimInstance Win32_DiskDrive | Where-Object { \$_.InterfaceType -eq "USB" } | ForEach-Object {
        Write-Warning ("  Disk: \$(\$_.Caption)")
        Write-Warning ("  Device ID: \$(\$_.DeviceID)")
    }

    # Pegar informações do drive especificado
    \$targetDrive = "${drive_letter}:"
    Write-Warning "Looking for drive \$targetDrive"
    
    \$volumes = Get-WmiObject Win32_Volume | Where-Object { \$_.DriveLetter -eq \$targetDrive }
    if (-not \$volumes) {
        Write-Error "Drive \$targetDrive not found"
        exit 1
    }
    
    Write-Warning "Found drive \$targetDrive"
    
    # Pegar o número do disco físico
    \$diskNumber = (Get-Partition -DriveLetter ${drive_letter}).DiskNumber
    if (-not \$diskNumber) {
        Write-Error "Could not find disk number for drive \$targetDrive"
        exit 1
    }
    
    Write-Warning "Found disk number: \$diskNumber"
    
    # Verificar se é um disco USB
    \$physicalDisk = Get-CimInstance Win32_DiskDrive | Where-Object { \$_.Index -eq \$diskNumber }
    if (-not \$physicalDisk) {
        Write-Error "Could not find physical disk with index \$diskNumber"
        exit 1
    }
    
    if (\$physicalDisk.InterfaceType -ne "USB") {
        Write-Error "Drive \$targetDrive is not a USB device (Interface: \$(\$physicalDisk.InterfaceType))"
        exit 1
    }
    
    # Se chegamos até aqui, encontramos o dispositivo USB correto
    Write-Output "PhysicalDrive\$diskNumber"
    Write-Warning "Device: \$(\$physicalDisk.Caption)"
    Write-Warning "Size: \$([math]::Round(\$physicalDisk.Size/1GB, 2)) GB"
    Write-Warning "Model: \$(\$physicalDisk.Model)"
    Write-Warning "Interface: \$(\$physicalDisk.InterfaceType)"
    
} catch {
    Write-Warning "Detailed error information:"
    Write-Warning \$_.Exception.Message
    Write-Warning \$_.Exception.ItemName
    Write-Error "Failed to get physical device: \$_"
    exit 1
}
POWERSHELL
        
        # Executar o script PowerShell
        local result
        result=$(powershell.exe -File "$ps_script" 2>&1)
        local exit_code=$?
        
        # Limpar arquivo temporário
        rm -f "$ps_script"
        
        if [ $exit_code -ne 0 ]; then
            echo -e "${RED}PowerShell Error: $(echo "$result" | grep -v "WARNING: ")${NC}" >&2
            return 1
        fi
        
        # Filtrar e mostrar resultados
        local device
        device=$(echo "$result" | grep -i "PhysicalDrive" | head -n1)
        
        if [ -n "$device" ]; then
            # Mostrar mensagens de diagnóstico (WARNING)
            echo "$result" | grep "WARNING: " | sed 's/WARNING: //' >&2
            # Retornar o dispositivo
            echo "$device"
        else
            echo -e "${RED}No physical device information found in output${NC}" >&2
            return 1
        fi
    }
    
    # Obter letra do drive do caminho WSL (ex: /mnt/d -> D:)
    drive_letter=$(echo "$USB_DEVICE" | sed -n 's/.*\/mnt\/\(.\).*/\1/p' | tr '[:lower:]' '[:upper:]')
    
    if [ -n "$drive_letter" ]; then
        echo -e "${YELLOW}Detecting USB device for drive ${drive_letter}:${NC}"
        
        # Verificar se o drive está montado
        if ! mountpoint -q "/mnt/${drive_letter,,}"; then
            echo -e "${RED}Drive ${drive_letter}: is not mounted in WSL${NC}"
            exit 1
        fi
        
        # Tentar obter o dispositivo físico
        if ! physical_device=$(get_physical_device "$drive_letter"); then
            echo -e "${RED}Failed to get physical device for drive ${drive_letter}:${NC}"
            echo "Please ensure the USB drive is properly connected and recognized by Windows"
            exit 1
        fi
        
        if [ -n "$physical_device" ]; then
            echo -e "${GREEN}Found USB device: $physical_device${NC}"
            if ! "$SCRIPT_DIR/create-syntropy-usb-enhanced.sh" "$physical_device" --node-name "$NODE_NAME"; then
                echo -e "${RED}Failed to create USB${NC}"
                exit 1
            fi
        else
            echo -e "${RED}No physical device found for drive ${drive_letter}:${NC}"
            echo "Debug information:"
            echo "• Drive letter: ${drive_letter}"
            echo "• WSL path: /mnt/${drive_letter,,}"
            echo "• Mount status: $(mountpoint "/mnt/${drive_letter,,}" 2>&1)"
            echo "• Windows drives:"
            powershell.exe -Command "Get-CimInstance Win32_DiskDrive | Format-List DeviceID, Caption, InterfaceType, Size"
            exit 1
        fi
    else
        echo -e "${RED}Failed to extract drive letter from path: $USB_DEVICE${NC}"
        exit 1
    fi
else
    # Em sistemas Linux normais, usa o dispositivo diretamente
    if ! "$SCRIPT_DIR/create-syntropy-usb-enhanced.sh" "/dev/$USB_DEVICE" --node-name "$NODE_NAME"; then
        echo -e "${RED}Failed to create USB${NC}"
        exit 1
    fi
fi

# Step 5: Next steps
echo -e "${BLUE}[Step 5/5] Next steps${NC}"
echo
echo -e "${GREEN}USB creation successful! Here are your next steps:${NC}"
echo
echo "1. Remove the USB drive safely"
echo "2. Insert the USB drive into your target hardware"
echo "3. Boot from USB (you may need to configure BIOS/UEFI boot order)"
echo "4. Wait for the automated installation (~30 minutes)"
echo "5. After installation, the node will automatically register"
echo "6. Use 'syntropy-connect $NODE_NAME' to connect to your node"
echo
echo -e "${YELLOW}Important:${NC}"
echo "- Default SSH user: admin"
echo "- SSH key authentication is enforced (password login disabled)"
echo "- Your SSH key has been automatically configured"
echo
echo -e "${GREEN}Your Syntropy node is ready to be deployed!${NC}"

# Source bashrc for immediate use
if [ -f "$HOME/.syntropy/config/syntropy.bashrc" ]; then
    source "$HOME/.syntropy/config/syntropy.bashrc"
    echo -e "${GREEN}Syntropy commands loaded. Try: syntropy-info${NC}"
fi