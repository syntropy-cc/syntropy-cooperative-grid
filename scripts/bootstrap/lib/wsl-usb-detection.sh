#!/bin/bash

# Syntropy Cooperative Grid - WSL USB Detection
# Version: 1.0.0

# Ensure we have required functions
if ! declare -F log >/dev/null; then
    # If log function is not available, create a simple version
    log() {
        local level="$1"
        shift
        echo "[$level] $*"
    }
fi

# Colors if not already defined
if [ -z "$YELLOW" ]; then
    YELLOW='\033[1;33m'
    RED='\033[0;31m'
    GREEN='\033[0;32m'
    BLUE='\033[0;34m'
    NC='\033[0m'
fi

# Detect if running in WSL
is_wsl() {
    if grep -qi microsoft /proc/version; then
        return 0  # True, is WSL
    else
        return 1  # False, not WSL
    fi
}

# Get Windows drive letter from PowerShell
get_windows_usb_devices() {
    log DEBUG "Detecting USB devices using PowerShell..."
    
    # PowerShell command to get USB devices
    local ps_command='Get-WmiObject Win32_DiskDrive | Where-Object{$_.InterfaceType -eq "USB"} | ForEach-Object{
        $disk = $_
        $partitions = "ASSOCIATORS OF {Win32_DiskDrive.DeviceID=`"$($disk.DeviceID.replace(\"\\\\\",\"\\\\\\\\\"))`"} WHERE AssocClass = Win32_DiskDriveToDiskPartition"
        $query = Get-WmiObject -Query $partitions
        $query | ForEach-Object{
            $partition = $_
            $drives = "ASSOCIATORS OF {Win32_DiskPartition.DeviceID=`"$($partition.DeviceID)`"} WHERE AssocClass = Win32_LogicalDiskToPartition"
            Get-WmiObject -Query $drives | ForEach-Object{
                New-Object PSObject -Property @{
                    Drive = $_.DeviceID
                    Size = $disk.Size
                    Model = $disk.Model
                    Path = $disk.DeviceID
                }
            }
        }
    } | ConvertTo-Json'
    
    # Run PowerShell command and capture output
    powershell.exe -Command "$ps_command" 2>/dev/null
}

# Convert Windows path to WSL path
convert_to_wsl_path() {
    local win_drive="$1"
    echo "/mnt/${win_drive,,}"  # Convert to lowercase
}

# Mount USB device in WSL
mount_usb_wsl() {
    local win_drive="$1"
    local wsl_path=$(convert_to_wsl_path "$win_drive")
    
    log INFO "Mounting Windows USB device ($win_drive) to WSL path ($wsl_path)"
    
    # Ensure mount point exists
    sudo mkdir -p "$wsl_path"
    
    # Mount with proper permissions
    if ! sudo mount -t drvfs "${win_drive}:" "$wsl_path" -o uid=$(id -u),gid=$(id -g); then
        log ERROR "Failed to mount USB device"
        return 1
    fi
    
    echo "$wsl_path"
    return 0
}

# Enhanced USB detection for WSL
detect_usb_devices_wsl() {
    if ! is_wsl; then
        # Not in WSL, use standard detection
        detect_usb_devices
        return
    fi
    
    log INFO "Detecting USB devices in WSL environment..."
    
    # Get USB devices from Windows
    local devices_json=$(get_windows_usb_devices)
    
    if [ -z "$devices_json" ]; then
        log ERROR "No USB devices found or PowerShell command failed"
        return 1
    fi
    
    # Parse JSON output (requires jq)
    if ! command -v jq &> /dev/null; then
        log ERROR "jq is required for WSL USB detection"
        return 1
    fi
    
    # Convert JSON to our format
    echo "$devices_json" | jq -r '.[] | "\(.Drive):\(.Size):\(.Model)"' | while read -r line; do
        if [ -n "$line" ]; then
            local drive=$(echo "$line" | cut -d: -f1)
            local size=$(echo "$line" | cut -d: -f2)
            local model=$(echo "$line" | cut -d: -f3)
            
            # Convert size to human readable format
            local size_gb=$(echo "scale=1; $size/1024/1024/1024" | bc)
            
            # Get WSL path
            local wsl_path=$(convert_to_wsl_path "$drive")
            
            echo "$wsl_path:${size_gb}GB:$model"
        fi
    done
}

# Select USB device with WSL support
select_usb_device_wsl() {
    if ! is_wsl; then
        # Not in WSL, use standard selection
        select_usb_device
        return
    fi
    
    log INFO "Selecting USB device in WSL environment..."
    
    # Get list of devices
    local devices=()
    while IFS= read -r line; do
        [ -n "$line" ] && devices+=("$line")
    done < <(detect_usb_devices_wsl)
    
    if [ ${#devices[@]} -eq 0 ]; then
        log ERROR "No USB devices found"
        return 1
    fi
    
    # Show devices
    echo "Available USB devices:"
    local i=1
    for device in "${devices[@]}"; do
        local path=$(echo "$device" | cut -d: -f1)
        local size=$(echo "$device" | cut -d: -f2)
        local model=$(echo "$device" | cut -d: -f3)
        echo "[$i] $path ($size) - $model"
        ((i++))
    done
    
    # Get user selection
    local selection
    read -p "Select USB device [1-${#devices[@]}]: " selection
    
    if ! [[ "$selection" =~ ^[0-9]+$ ]] || [ "$selection" -lt 1 ] || [ "$selection" -gt ${#devices[@]} ]; then
        log ERROR "Invalid selection"
        return 1
    fi
    
    # Return selected device
    echo "${devices[$((selection-1))]}" | cut -d: -f1
}