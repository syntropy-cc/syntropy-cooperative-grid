#!/bin/bash

# Syntropy Cooperative Grid - Quick Start Script
# Version: 3.0.0 - Complete WSL-Ready Plug & Play Version
# This script is fully self-contained and handles all WSL USB detection/conversion automatically

set -e

# ============================================================================
# CONFIGURATION & COLORS
# ============================================================================

PURPLE='\033[0;35m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
CYAN='\033[0;36m'
NC='\033[0m'

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# ============================================================================
# CORE FUNCTIONS
# ============================================================================

# Function to detect if running in WSL
is_wsl() {
    if grep -qi microsoft /proc/version; then
        return 0  # True, is WSL
    else
        return 1  # False, not WSL
    fi
}

# Logging functions (embedded to avoid dependencies)
log() {
    local level="$1"
    shift
    local message="$*"
    
    case "$level" in
        INFO)    echo -e "${BLUE}[INFO]${NC} $message" ;;
        SUCCESS) echo -e "${GREEN}[SUCCESS]${NC} $message" ;;
        WARN)    echo -e "${YELLOW}[WARN]${NC} $message" ;;
        ERROR)   echo -e "${RED}[ERROR]${NC} $message" ;;
        DEBUG)   [ "${DEBUG:-0}" = "1" ] && echo -e "${CYAN}[DEBUG]${NC} $message" ;;
        *)       echo "$message" ;;
    esac
}

# Show banner
show_banner() {
    echo -e "${PURPLE}"
    cat << 'EOF'
╔════════════════════════════════════════════════════════════════════════════╗
║                    SYNTROPY COOPERATIVE GRID                              ║
║                         Quick Start Setup                                 ║
║                          Version 3.0.0                                   ║
║                                                                           ║
║  This script will guide you through setting up your first Syntropy node  ║
╚════════════════════════════════════════════════════════════════════════════╝
EOF
    echo -e "${NC}"
}

# Function to check and install prerequisites
check_and_install_prerequisites() {
    log INFO "Checking prerequisites..."
    
    local missing_tools=()
    local required_tools=(jq curl wget git python3 ssh-keygen openssl bc)
    
    for tool in "${required_tools[@]}"; do
        if ! command -v "$tool" >/dev/null 2>&1; then
            missing_tools+=("$tool")
        fi
    done
    
    if [ ${#missing_tools[@]} -gt 0 ]; then
        log WARN "Missing tools: ${missing_tools[*]}"
        log INFO "Installing missing prerequisites..."
        
        # Update package list
        sudo apt-get update >/dev/null 2>&1 || true
        
        # Map tool names to package names
        for tool in "${missing_tools[@]}"; do
            case "$tool" in
                ssh-keygen) sudo apt-get install -y openssh-client >/dev/null 2>&1 ;;
                *)          sudo apt-get install -y "$tool" >/dev/null 2>&1 ;;
            esac
        done
        
        log SUCCESS "Prerequisites installed"
    else
        log SUCCESS "All prerequisites are already installed"
    fi
}

# ============================================================================
# WSL USB DETECTION & CONVERSION FUNCTIONS
# ============================================================================

# Convert Windows PhysicalDrive to WSL /dev/sdX format (enhanced)
convert_physical_drive_to_wsl() {
    local physical_drive="$1"
    local disk_number="$2"  # Optional: specific disk number for cross-reference
    
    # Extract disk number from PhysicalDrive format
    local extracted_disk_number
    if [[ "$physical_drive" =~ PhysicalDrive([0-9]+) ]]; then
        extracted_disk_number="${BASH_REMATCH[1]}"
    else
        log ERROR "Invalid PhysicalDrive format: $physical_drive"
        return 1
    fi
    
    # Use provided disk number if available, otherwise use extracted
    local target_disk_number="${disk_number:-$extracted_disk_number}"
    
    # Convert to WSL device path
    # PhysicalDrive0 = /dev/sda, PhysicalDrive1 = /dev/sdb, etc.
    local letter_ascii=$((97 + target_disk_number))
    local device_letter=$(printf "\\$(printf '%03o' $letter_ascii)")
    local wsl_device="/dev/sd${device_letter}"
    
    log DEBUG "Converting $physical_drive (disk $target_disk_number) to $wsl_device"
    
    # Validate that the device exists
    if [ -b "$wsl_device" ]; then
        log DEBUG "Device $wsl_device exists and is accessible"
    echo "$wsl_device"
        return 0
    else
        log WARN "Converted device $wsl_device does not exist or is not accessible"
        
        # Try to find alternative devices
        log INFO "Searching for alternative USB devices..."
        for dev in /dev/sd[a-z]; do
            if [ -b "$dev" ] && [ "$dev" != "/dev/sda" ]; then
                local dev_name=$(basename "$dev")
                local removable=$(cat "/sys/block/$dev_name/removable" 2>/dev/null || echo "0")
                if [ "$removable" = "1" ]; then
                    log INFO "Found alternative removable device: $dev"
                    echo "$dev"
                    return 0
                fi
            fi
        done
        
        log ERROR "No suitable alternative device found"
        return 1
    fi
}

# Get USB devices from Windows
get_windows_usb_devices() {
    local ps_command='
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
        
        # Get drive letters
        $partitions = Get-WmiObject -Query "ASSOCIATORS OF {Win32_DiskDrive.DeviceID=\"$($disk.DeviceID)\"} WHERE AssocClass = Win32_DiskDriveToDiskPartition"
        $letters = @()
        
        foreach ($partition in $partitions) {
            $logicalDisks = Get-WmiObject -Query "ASSOCIATORS OF {Win32_DiskPartition.DeviceID=\"$($partition.DeviceID)\"} WHERE AssocClass = Win32_LogicalDiskToPartition"
            foreach ($ld in $logicalDisks) {
                if ($ld.DeviceID) {
                    $letters += $ld.DeviceID.Replace(":", "")
                }
            }
        }
        
        $letterStr = if ($letters.Count -gt 0) { $letters -join "," } else { "NONE" }
        Write-Output "USB|$diskNum|$model|$size|$letterStr"
    }'
    
    powershell.exe -ExecutionPolicy Bypass -Command "$ps_command" 2>/dev/null
}

# Verify WSL device exists and is accessible (enhanced)
verify_wsl_device() {
    local device="$1"
    local expected_size_gb="$2"  # Optional: expected size for validation
    
    # Check if device exists
    if [ ! -e "$device" ]; then
        log WARN "Device $device not found in WSL filesystem"
        
        # Try to create device node if missing
        local device_name=$(basename "$device")
        local major_minor=$(ls -la /dev/sd* 2>/dev/null | grep -m1 "sd" | awk '{print $5 $6}' | tr -d ',')
        
        if [ -n "$major_minor" ]; then
            log INFO "Attempting to create device node..."
            sudo mknod "$device" b 8 $((16 * ($(echo $device_name | sed 's/sd//' | od -An -N1 -tc | tr -d ' ') - 97))) 2>/dev/null || true
        fi
    fi
    
    # Check if it's a block device
    if [ -b "$device" ]; then
        log DEBUG "Device $device exists and is a block device"
        
        # Additional validation if expected size is provided
        if [ -n "$expected_size_gb" ] && [ "$expected_size_gb" != "0" ]; then
            local actual_size_bytes=$(sudo blockdev --getsize64 "$device" 2>/dev/null || echo "0")
            local actual_size_gb=$((actual_size_bytes / 1024 / 1024 / 1024))
            local size_diff=$((actual_size_gb - expected_size_gb))
            local size_diff_abs=${size_diff#-}  # Absolute value
            local size_tolerance=$((expected_size_gb / 10))  # 10% tolerance
            
            if [ "$size_diff_abs" -le "$size_tolerance" ]; then
                log DEBUG "Device size matches expected (${actual_size_gb}GB ≈ ${expected_size_gb}GB)"
            else
                log WARN "Device size doesn't match expected (${actual_size_gb}GB vs ${expected_size_gb}GB)"
            fi
        fi
        
        return 0
    else
        log WARN "Device $device is not accessible as a block device"
        return 1
    fi
}

# Enhanced USB device detection for WSL
detect_wsl_usb_devices() {
    log INFO "Detecting USB devices in WSL environment..."
    
    local devices=()
    local device_info=""
    
    # Get USB devices from Windows for cross-reference
    local windows_devices=$(get_windows_usb_devices)
    
    if [ -n "$windows_devices" ]; then
        log INFO "Windows USB devices detected:"
        while IFS='|' read -r type num model size letters; do
            if [ "$type" = "USB" ]; then
                echo "  Disk $num: $model (${size}GB)"
                [ "$letters" != "NONE" ] && echo "      Drive letters: $letters"
                
                # Try to find corresponding WSL device
                local wsl_device=$(convert_physical_drive_to_wsl "PhysicalDrive$num" "$num")
                if [ $? -eq 0 ] && [ -b "$wsl_device" ]; then
                    device_info="$wsl_device:$size:$model:PhysicalDrive$num"
                    devices+=("$device_info")
                    log SUCCESS "Mapped PhysicalDrive$num → $wsl_device"
                else
                    log WARN "Could not map PhysicalDrive$num to WSL device"
                fi
            fi
        done <<< "$windows_devices"
    fi
    
    # If no devices found via Windows mapping, try direct WSL detection
    if [ ${#devices[@]} -eq 0 ]; then
        log INFO "No devices found via Windows mapping, trying direct WSL detection..."
        
        for dev in /dev/sd[a-z]; do
            if [ -b "$dev" ] && [ "$dev" != "/dev/sda" ]; then
                local dev_name=$(basename "$dev")
                local removable=$(cat "/sys/block/$dev_name/removable" 2>/dev/null || echo "0")
                local size_bytes=$(sudo blockdev --getsize64 "$dev" 2>/dev/null || echo "0")
                local size_gb=$((size_bytes / 1024 / 1024 / 1024))
                local model=$(lsblk -d -n -o MODEL "$dev" 2>/dev/null || echo "Unknown")
                
                # Consider devices that are removable or in typical USB size range
                if [ "$removable" = "1" ] || ([ "$size_gb" -ge 1 ] && [ "$size_gb" -le 1024 ]); then
                    device_info="$dev:$size_gb:$model:WSL-Detected"
                    devices+=("$device_info")
                    log INFO "Found WSL device: $dev (${size_gb}GB, $model)"
                fi
            fi
        done
    fi
    
    # Output devices in format: device:size:model:source
    for device in "${devices[@]}"; do
        echo "$device"
    done
    
    return ${#devices[@]}
}

# ============================================================================
# ENHANCED USB CREATION WRAPPER
# ============================================================================

# Create wrapper script that handles WSL device conversion
create_wsl_wrapper_script() {
    local wrapper_script="/tmp/syntropy-usb-wrapper-$$.sh"
    
    cat > "$wrapper_script" << 'WRAPPER_EOF'
#!/bin/bash

# Temporary wrapper for WSL USB creation
set -e

PHYSICAL_DRIVE="$1"
NODE_NAME="$2"
SCRIPT_DIR="$3"
SELECTED_DISK_NUM="$4"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

echo -e "${CYAN}WSL USB Creation Wrapper${NC}"

# Enhanced device detection and mapping
echo -e "${BLUE}Device Detection and Mapping:${NC}"
echo "  Windows Physical Drive: $PHYSICAL_DRIVE"
echo "  Selected Disk Number: $SELECTED_DISK_NUM"

# First, try to get the actual WSL device using PowerShell to cross-reference
echo -e "${BLUE}Cross-referencing with Windows USB detection...${NC}"

# Get USB devices from Windows again to ensure consistency
PS_COMMAND='
$ErrorActionPreference = "SilentlyContinue"
$usbDrives = Get-WmiObject Win32_DiskDrive | Where-Object { $_.InterfaceType -eq "USB" }
foreach ($disk in $usbDrives) {
    if ($disk.Index -eq '$SELECTED_DISK_NUM') {
        $diskNum = $disk.Index
        $model = $disk.Model
        $size = [math]::Round($disk.Size / 1GB, 2)
        Write-Output "USB|$diskNum|$model|$size"
        break
    }
}'

WINDOWS_USB_INFO=$(powershell.exe -ExecutionPolicy Bypass -Command "$PS_COMMAND" 2>/dev/null)

if [ -n "$WINDOWS_USB_INFO" ]; then
    echo "  Windows USB Info: $WINDOWS_USB_INFO"
    WINDOWS_SIZE=$(echo "$WINDOWS_USB_INFO" | cut -d'|' -f4)
    WINDOWS_MODEL=$(echo "$WINDOWS_USB_INFO" | cut -d'|' -f3)
    echo "  Expected Size: ${WINDOWS_SIZE}GB"
    echo "  Expected Model: $WINDOWS_MODEL"
else
    echo -e "${YELLOW}Warning: Could not cross-reference with Windows USB info${NC}"
fi

# List all available WSL devices with detailed info
echo -e "${BLUE}Available WSL block devices:${NC}"
lsblk -d -o NAME,SIZE,TYPE,MODEL 2>/dev/null | grep disk || true

# Try to find the matching device by size and characteristics
WSL_DEVICE=""
BEST_MATCH=""

echo -e "${BLUE}Searching for matching USB device...${NC}"

# Check each available device
for dev in /dev/sd[a-z]; do
            if [ -b "$dev" ]; then
        DEVICE_NAME=$(basename "$dev")
        SIZE_BYTES=$(sudo blockdev --getsize64 "$dev" 2>/dev/null || echo "0")
        SIZE_GB=$((SIZE_BYTES / 1024 / 1024 / 1024))
        MODEL=$(lsblk -d -n -o MODEL "$dev" 2>/dev/null || echo "Unknown")
        REMOVABLE=$(cat "/sys/block/$DEVICE_NAME/removable" 2>/dev/null || echo "0")
        
        echo "  Checking $dev: ${SIZE_GB}GB, Model: $MODEL, Removable: $REMOVABLE"
        
        # Skip system disk
        if [ "$dev" = "/dev/sda" ]; then
            echo "    → Skipping system disk"
            continue
        fi
        
        # Check if size matches (within 10% tolerance)
        if [ -n "$WINDOWS_SIZE" ] && [ "$WINDOWS_SIZE" != "0" ]; then
            # Convert Windows size to integer (remove decimal part)
            WINDOWS_SIZE_INT=$(echo "$WINDOWS_SIZE" | cut -d'.' -f1)
            
            if [ -n "$WINDOWS_SIZE_INT" ] && [ "$WINDOWS_SIZE_INT" -gt 0 ]; then
                SIZE_DIFF=$((SIZE_GB - WINDOWS_SIZE_INT))
                SIZE_DIFF_ABS=${SIZE_DIFF#-}  # Absolute value
                SIZE_TOLERANCE=$((WINDOWS_SIZE_INT / 10))  # 10% tolerance
                
                if [ "$SIZE_DIFF_ABS" -le "$SIZE_TOLERANCE" ]; then
                    echo "    → Size matches Windows USB (${SIZE_GB}GB ≈ ${WINDOWS_SIZE_INT}GB)"
                    BEST_MATCH="$dev"
                else
                    echo "    → Size doesn't match (${SIZE_GB}GB vs ${WINDOWS_SIZE_INT}GB)"
                fi
            else
                echo "    → Could not parse Windows size: $WINDOWS_SIZE"
            fi
        fi
        
        # Prefer removable devices
        if [ "$REMOVABLE" = "1" ]; then
            echo "    → Marked as removable device"
            if [ -z "$BEST_MATCH" ]; then
                BEST_MATCH="$dev"
            fi
        fi
        
        # Check for typical USB size range (1GB - 1TB)
        if [ "$SIZE_GB" -ge 1 ] && [ "$SIZE_GB" -le 1024 ]; then
            echo "    → Size in typical USB range"
            if [ -z "$BEST_MATCH" ]; then
                BEST_MATCH="$dev"
            fi
                fi
            fi
        done

# Select the best match
if [ -n "$BEST_MATCH" ]; then
    WSL_DEVICE="$BEST_MATCH"
    echo -e "${GREEN}Selected device: $WSL_DEVICE${NC}"
else
    echo -e "${YELLOW}No automatic match found. Manual selection required.${NC}"
    
    # Manual selection fallback
    echo "Available devices for manual selection:"
    for dev in /dev/sd[a-z]; do
        if [ -b "$dev" ] && [ "$dev" != "/dev/sda" ]; then
            SIZE_BYTES=$(sudo blockdev --getsize64 "$dev" 2>/dev/null || echo "0")
            SIZE_GB=$((SIZE_BYTES / 1024 / 1024 / 1024))
            MODEL=$(lsblk -d -n -o MODEL "$dev" 2>/dev/null || echo "Unknown")
            echo "  $dev - ${SIZE_GB}GB - $MODEL"
        fi
    done
    
    echo ""
    read -p "Enter the device path (e.g., /dev/sdb): " WSL_DEVICE
    
    if [ ! -b "$WSL_DEVICE" ]; then
        echo -e "${RED}Error: Device $WSL_DEVICE not found or not accessible${NC}"
        exit 1
    fi
fi

# Final validation
echo -e "${BLUE}Final device validation:${NC}"
echo "  Selected WSL Device: $WSL_DEVICE"
echo "  Device exists: $([ -b "$WSL_DEVICE" ] && echo "Yes" || echo "No")"

if [ -b "$WSL_DEVICE" ]; then
    SIZE_BYTES=$(sudo blockdev --getsize64 "$WSL_DEVICE" 2>/dev/null || echo "0")
    SIZE_GB=$((SIZE_BYTES / 1024 / 1024 / 1024))
    MODEL=$(lsblk -d -n -o MODEL "$WSL_DEVICE" 2>/dev/null || echo "Unknown")
    echo "  Size: ${SIZE_GB}GB"
    echo "  Model: $MODEL"
    
    # Show device details
    echo -e "${BLUE}Device details:${NC}"
    lsblk "$WSL_DEVICE" 2>/dev/null || echo "  Unable to show device details"
    
    echo -e "${GREEN}Using device: $WSL_DEVICE${NC}"
    
    # Call the actual creation script with WSL device
    if [ -f "$SCRIPT_DIR/create-syntropy-usb-enhanced.sh" ]; then
        # Source required libraries
        if [ -f "$SCRIPT_DIR/lib/colors.sh" ]; then
            source "$SCRIPT_DIR/lib/colors.sh"
        fi
        if [ -f "$SCRIPT_DIR/lib/logging.sh" ]; then
            source "$SCRIPT_DIR/lib/logging.sh"
        fi
        
        # Set required environment variables
        export WSL_MODE=false  # Tell the script we've already converted the device
        
        # Execute with WSL device
        "$SCRIPT_DIR/create-syntropy-usb-enhanced.sh" "$WSL_DEVICE" --node-name "$NODE_NAME"
    else
        echo -e "${RED}Creation script not found at: $SCRIPT_DIR/create-syntropy-usb-enhanced.sh${NC}"
        exit 1
    fi
else
    echo -e "${RED}Error: No suitable USB device found${NC}"
    echo -e "${YELLOW}Troubleshooting tips:${NC}"
    echo "1. Ensure USB is connected and recognized by Windows"
    echo "2. Try restarting WSL: wsl --shutdown (from PowerShell)"
    echo "3. Check if USB appears in Windows Disk Management"
    echo "4. Ensure USB is not in use by Windows applications"
    exit 1
fi
WRAPPER_EOF
    
    chmod +x "$wrapper_script"
    echo "$wrapper_script"
}

# ============================================================================
# MAIN SETUP STEPS
# ============================================================================

# Step 1: Install prerequisites
step1_prerequisites() {
    echo -e "${BLUE}═══ Step 1/5: Installing Prerequisites ═══${NC}"
    
    # First try the official script
    if [ -f "$SCRIPT_DIR/install-prerequisites.sh" ]; then
        if bash "$SCRIPT_DIR/install-prerequisites.sh"; then
            log SUCCESS "Prerequisites installed via official script"
            return 0
        fi
    fi
    
    # Fallback to embedded installation
    check_and_install_prerequisites
}

# Step 2: Setup management environment
step2_management() {
    echo -e "${BLUE}═══ Step 2/5: Setting up Management Environment ═══${NC}"
    
    local syntropy_dir="$HOME/.syntropy"
    
    # Create directory structure
    mkdir -p "$syntropy_dir"/{nodes,keys,config,cache,scripts,backups}
    
    # Try official script first
    if [ -f "$SCRIPT_DIR/setup-syntropy-management.sh" ]; then
        if bash "$SCRIPT_DIR/setup-syntropy-management.sh"; then
            log SUCCESS "Management environment configured"
            return 0
        fi
    fi
    
    # Fallback: Create minimal configuration
    if [ ! -f "$syntropy_dir/config/manager.json" ]; then
        cat > "$syntropy_dir/config/manager.json" << EOF
{
  "version": "1.0.0",
  "created": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "manager_id": "mgr-$(openssl rand -hex 8)"
}
EOF
        log SUCCESS "Basic management environment created"
    else
        log SUCCESS "Management environment already exists"
    fi
}

# Step 3: Detect USB devices
step3_detect_usb() {
    echo -e "${BLUE}═══ Step 3/5: Detecting USB Devices ═══${NC}"
    
    if is_wsl; then
        log INFO "Running in WSL environment"
        
        # Get USB devices from Windows
        local devices=$(get_windows_usb_devices)
        
        if echo "$devices" | grep -q "NO_USB_FOUND"; then
            log ERROR "No USB devices detected by Windows"
            echo "Please ensure:"
            echo "  1. USB drive is connected"
            echo "  2. Windows recognizes the drive"
            echo "  3. The drive appears in File Explorer"
            return 1
        fi
        
        # Parse and display devices
        log SUCCESS "USB devices detected:"
        local count=0
        while IFS='|' read -r type num model size letters; do
            if [ "$type" = "USB" ]; then
                ((count++))
                echo "  [$count] Disk $num: $model (${size}GB)"
                [ "$letters" != "NONE" ] && echo "      Drive letters: $letters"
            fi
        done <<< "$devices"
        
        if [ $count -eq 0 ]; then
            log ERROR "No USB devices found"
            return 1
        fi
        
        # Store for later use
        USB_DEVICES_DATA="$devices"
        USB_DEVICE_COUNT=$count
    else
        # Linux detection
        log INFO "Running in standard Linux environment"
        
        local removable_devices=$(lsblk -d -o NAME,SIZE,RM | grep " 1$" | grep -v "sr0")
        
        if [ -z "$removable_devices" ]; then
            log ERROR "No removable devices detected"
            return 1
        fi
        
        log SUCCESS "Removable devices detected:"
        echo "$removable_devices"
    fi
}

# Step 4: Configure node
step4_configure_node() {
    echo -e "${BLUE}═══ Step 4/5: Node Configuration ═══${NC}"
    
    # Get node name
    while true; do
        echo
        read -p "Enter node name (e.g., syntropy-node-01): " NODE_NAME
        
        if [[ "$NODE_NAME" =~ ^[a-zA-Z0-9-]+$ ]]; then
            log SUCCESS "Node name: $NODE_NAME"
            break
        else
            log ERROR "Invalid name. Use only letters, numbers, and hyphens"
        fi
    done
    
    # Select USB device
    echo
    log INFO "Selecting USB device..."
    
    if is_wsl; then
        # WSL device selection
        local devices=()
        local disk_nums=()
        
        while IFS='|' read -r type num model size letters; do
            if [ "$type" = "USB" ]; then
                devices+=("Disk $num: $model (${size}GB)")
                disk_nums+=("$num")
            fi
        done <<< "$USB_DEVICES_DATA"
        
        if [ ${#devices[@]} -eq 1 ]; then
            log INFO "Using the only USB device available"
            SELECTED_DISK_NUM="${disk_nums[0]}"
        else
            echo "Available USB devices:"
            for i in "${!devices[@]}"; do
                echo "  [$((i+1))] ${devices[$i]}"
            done
            
            while true; do
                read -p "Select device (1-${#devices[@]}): " selection
                
                if [[ "$selection" =~ ^[0-9]+$ ]] && [ "$selection" -ge 1 ] && [ "$selection" -le ${#devices[@]} ]; then
                    SELECTED_DISK_NUM="${disk_nums[$((selection-1))]}"
                    break
                else
                    log ERROR "Invalid selection"
                fi
            done
        fi
        
        PHYSICAL_DEVICE="PhysicalDrive$SELECTED_DISK_NUM"
        log SUCCESS "Selected: $PHYSICAL_DEVICE"
    else
        # Linux device selection
        local devices=($(lsblk -d -o NAME -n | grep -E "^sd[b-z]$"))
        
        echo "Available devices:"
        for i in "${!devices[@]}"; do
            local dev_info=$(lsblk -d -o SIZE,MODEL -n "/dev/${devices[$i]}" 2>/dev/null)
            echo "  [$((i+1))] /dev/${devices[$i]} - $dev_info"
        done
        
        read -p "Select device (1-${#devices[@]}): " selection
        
        if [[ "$selection" =~ ^[0-9]+$ ]] && [ "$selection" -ge 1 ] && [ "$selection" -le ${#devices[@]} ]; then
            USB_DEVICE="/dev/${devices[$((selection-1))]}"
            log SUCCESS "Selected: $USB_DEVICE"
        else
            log ERROR "Invalid selection"
            return 1
        fi
    fi
    
    # Final confirmation
    echo
    echo -e "${RED}╔════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${RED}║                         WARNING                                ║${NC}"
    echo -e "${RED}║     ALL DATA ON THE SELECTED DEVICE WILL BE ERASED!          ║${NC}"
    echo -e "${RED}╚════════════════════════════════════════════════════════════════╝${NC}"
    echo
    
    read -p "Type 'yes' to confirm: " confirmation
    if [ "$confirmation" != "yes" ]; then
        log WARN "Operation cancelled"
        return 1
    fi
}

# Step 5: Create USB
step5_create_usb() {
    echo -e "${BLUE}═══ Step 5/5: Creating USB ═══${NC}"
    
    # Detect environment and choose appropriate strategy
    local environment=$(detect_environment)
    
    case "$environment" in
        "WSL")
            log INFO "Environment: WSL - Using hybrid Windows/Linux approach"
            create_usb_wsl_hybrid
            ;;
        "LINUX")
            log INFO "Environment: Linux native - Using standard Linux tools"
            create_usb_linux_native
            ;;
        "WINDOWS")
            log INFO "Environment: Windows native - Using PowerShell tools"
            create_usb_windows_native
            ;;
        *)
            log ERROR "Unknown environment: $environment"
            return 1
            ;;
    esac
}

# Detect the current environment
detect_environment() {
    # Check if running in WSL
    if is_wsl; then
        echo "WSL"
        return 0
    fi
    
    # Check if running on Windows (PowerShell available)
    if command -v powershell.exe >/dev/null 2>&1; then
        echo "WINDOWS"
        return 0
    fi
    
    # Default to Linux
    echo "LINUX"
    return 0
}

# Create USB using WSL hybrid approach (Windows PowerShell + Linux processing)
create_usb_wsl_hybrid() {
    log INFO "Using WSL hybrid approach for USB creation..."
    
    # In WSL, we know USB was detected via PowerShell in Step 3
    # So we should always try Windows approach first
    log INFO "USB was detected via PowerShell in Step 3 - using Windows USB creation"
    create_usb_windows_powershell
}

# Create USB using Linux native tools
create_usb_linux_native() {
    log INFO "Using Linux native tools for USB creation..."
    
    if [ -f "$SCRIPT_DIR/create-syntropy-usb-enhanced.sh" ]; then
        if "$SCRIPT_DIR/create-syntropy-usb-enhanced.sh" "$USB_DEVICE" --node-name "$NODE_NAME"; then
            log SUCCESS "USB creation completed!"
            return 0
        else
            log ERROR "USB creation failed"
            return 1
        fi
    else
        log ERROR "USB creation script not found"
        return 1
    fi
}

# Create USB using Windows native PowerShell
create_usb_windows_native() {
    log INFO "Using Windows native PowerShell for USB creation..."
    create_usb_windows_powershell
}

# Check if we can access USB device via Windows
can_access_usb_via_windows() {
    # Try to get USB device info via PowerShell
    local ps_command='
    $ErrorActionPreference = "SilentlyContinue"
    $usbDrives = Get-WmiObject Win32_DiskDrive | Where-Object { $_.InterfaceType -eq "USB" }
    if ($usbDrives.Count -gt 0) {
        Write-Output "USB_ACCESSIBLE"
    } else {
        Write-Output "USB_NOT_ACCESSIBLE"
    }'
    
    local result=$(powershell.exe -ExecutionPolicy Bypass -Command "$ps_command" 2>/dev/null)
    
    if [ "$result" = "USB_ACCESSIBLE" ]; then
                return 0
            else
                return 1
            fi
}

# Create USB using Windows PowerShell (hybrid approach)
create_usb_windows_powershell() {
    log INFO "Creating USB using Windows PowerShell approach..."
    
    # Validate required variables
    if [ -z "$PHYSICAL_DEVICE" ] || [ -z "$NODE_NAME" ] || [ -z "$SELECTED_DISK_NUM" ]; then
        log ERROR "Missing required variables for Windows USB creation"
        log ERROR "PHYSICAL_DEVICE: $PHYSICAL_DEVICE"
        log ERROR "NODE_NAME: $NODE_NAME"
        log ERROR "SELECTED_DISK_NUM: $SELECTED_DISK_NUM"
        return 1
    fi
    
    # Create PowerShell script for USB creation
    local ps_script=$(create_windows_usb_script)
    
    if [ -n "$ps_script" ]; then
        log INFO "Executing Windows USB creation script..."
        log INFO "Parameters:"
        log INFO "  Physical Drive: $PHYSICAL_DEVICE"
        log INFO "  Node Name: $NODE_NAME"
        log INFO "  Selected Disk Number: $SELECTED_DISK_NUM"
        log INFO "  Script Directory: $SCRIPT_DIR"
        
        # Execute the PowerShell script with parameters
        local ps_output
        if ps_output=$(powershell.exe -ExecutionPolicy Bypass -File "$ps_script" -PhysicalDrive "$PHYSICAL_DEVICE" -NodeName "$NODE_NAME" -ScriptDir "$SCRIPT_DIR" -SelectedDiskNum "$SELECTED_DISK_NUM" 2>&1); then
            log SUCCESS "USB creation completed via Windows PowerShell!"
            echo "$ps_output" | while IFS= read -r line; do
                log INFO "PS: $line"
            done
            rm -f "$ps_script"
            return 0
        else
            log ERROR "USB creation failed via Windows PowerShell"
            log ERROR "PowerShell output:"
            echo "$ps_output" | while IFS= read -r line; do
                log ERROR "PS: $line"
            done
            rm -f "$ps_script"
            return 1
        fi
    else
        log ERROR "Failed to create Windows USB script"
        return 1
    fi
}

# Create USB using WSL fallback (when Windows approach fails)
create_usb_wsl_fallback() {
    log INFO "Using WSL fallback approach..."
    
    # Create and use wrapper script
    local wrapper=$(create_wsl_wrapper_script)
    
    if bash "$wrapper" "$PHYSICAL_DEVICE" "$NODE_NAME" "$SCRIPT_DIR" "$SELECTED_DISK_NUM"; then
        log SUCCESS "USB creation completed via WSL fallback!"
        rm -f "$wrapper"
        return 0
    else
        log ERROR "USB creation failed via WSL fallback"
        rm -f "$wrapper"
        return 1
    fi
}

# Create Windows PowerShell script for USB creation
create_windows_usb_script() {
    local ps_script="/tmp/syntropy-usb-windows-$$.ps1"
    
    cat > "$ps_script" << 'PS_EOF'
# Syntropy Cooperative Grid - Windows USB Creator
# PowerShell script for creating USB drives on Windows

param(
    [string]$PhysicalDrive = "",
    [string]$NodeName = "",
    [string]$ScriptDir = "",
    [string]$SelectedDiskNum = ""
)

# Set error action preference
$ErrorActionPreference = "Stop"

# Colors for output
$Colors = @{
    Red = "Red"
    Green = "Green"
    Yellow = "Yellow"
    Blue = "Blue"
    Cyan = "Cyan"
    White = "White"
}

function Write-ColorOutput {
    param(
        [string]$Message,
        [string]$Color = "White"
    )
    Write-Host $Message -ForegroundColor $Colors[$Color]
}

function Get-USBDeviceInfo {
    param([string]$DiskNum)
    
    try {
        Write-ColorOutput "Searching for USB device with disk number: $DiskNum" "Blue"
        
        # Get all USB drives
        $usbDrives = Get-WmiObject Win32_DiskDrive | Where-Object { $_.InterfaceType -eq "USB" }
        
        Write-ColorOutput "Found $($usbDrives.Count) USB drives" "Cyan"
        
        foreach ($drive in $usbDrives) {
            Write-ColorOutput "  Disk $($drive.Index): $($drive.Model) ($([math]::Round($drive.Size / 1GB, 2)) GB)" "White"
        }
        
        # Find the specific drive
        $usbDrive = $usbDrives | Where-Object { $_.Index -eq [int]$DiskNum }
        
        if ($usbDrive) {
            Write-ColorOutput "Found target USB device: $($usbDrive.Model)" "Green"
            return @{
                DeviceID = $usbDrive.DeviceID
                Model = $usbDrive.Model
                Size = [math]::Round($usbDrive.Size / 1GB, 2)
                Index = $usbDrive.Index
            }
        } else {
            Write-ColorOutput "USB device with disk number $DiskNum not found" "Red"
            return $null
        }
    }
    catch {
        Write-ColorOutput "Error getting USB device info: $($_.Exception.Message)" "Red"
        return $null
    }
}

function Test-USBDeviceAccess {
    param([string]$DeviceID)
    
    try {
        # Try to get partition information
        $partitions = Get-WmiObject -Query "ASSOCIATORS OF {Win32_DiskDrive.DeviceID='$DeviceID'} WHERE AssocClass = Win32_DiskDriveToDiskPartition"
        
        if ($partitions) {
            return $true
        }
        return $false
    }
    catch {
        return $false
    }
}

function Format-USBDevice {
    param(
        [string]$DeviceID,
        [string]$NodeName,
        [bool]$IsAdmin = $false
    )
    
    try {
        Write-ColorOutput "Formatting USB device: $DeviceID" "Blue"
        
        # Get the disk number from DeviceID (handle different formats)
        Write-ColorOutput "Parsing DeviceID: $DeviceID" "Cyan"
        
        # Try different regex patterns to extract disk number
        $diskNumber = $null
        
        # Pattern 1: PhysicalDriveX
        if (-not $diskNumber) {
            $match = [regex]::Match($DeviceID, 'PhysicalDrive(\d+)', [System.Text.RegularExpressions.RegexOptions]::IgnoreCase)
            if ($match.Success) {
                $diskNumber = $match.Groups[1].Value
                Write-ColorOutput "Extracted disk number using Pattern 1: $diskNumber" "Green"
            }
        }
        
        # Pattern 2: \\.\PHYSICALDRIVEX (with escape characters)
        if (-not $diskNumber) {
            $match = [regex]::Match($DeviceID, '\\\\.\\\\(PhysicalDrive\d+)', [System.Text.RegularExpressions.RegexOptions]::IgnoreCase)
            if ($match.Success) {
                $diskNumber = [regex]::Match($match.Groups[1].Value, 'PhysicalDrive(\d+)', [System.Text.RegularExpressions.RegexOptions]::IgnoreCase).Groups[1].Value
                Write-ColorOutput "Extracted disk number using Pattern 2: $diskNumber" "Green"
            }
        }
        
        # Pattern 3: Just extract any number at the end
        if (-not $diskNumber) {
            $match = [regex]::Match($DeviceID, '(\d+)$')
            if ($match.Success) {
                $diskNumber = $match.Groups[1].Value
                Write-ColorOutput "Extracted disk number using Pattern 3: $diskNumber" "Green"
            }
        }
        
        if (-not $diskNumber) {
            throw "Could not extract disk number from DeviceID: $DeviceID"
        }
        
        Write-ColorOutput "Final disk number: $diskNumber" "Cyan"
        
        Write-ColorOutput "Disk Number: $diskNumber" "Cyan"
        
        # Check if we have admin privileges for disk operations
        if (-not $IsAdmin) {
            Write-ColorOutput "No administrative privileges - attempting limited operations..." "Yellow"
            
            # Try to find the USB drive letter and work with it directly
            $usbDrive = Get-WmiObject Win32_LogicalDisk | Where-Object { 
                $_.DriveType -eq 2 -and 
                $_.Size -gt 0 -and 
                $_.Size -lt 200GB  # Typical USB size range
            }
            
            if ($usbDrive) {
                $driveLetter = $usbDrive.DeviceID
                Write-ColorOutput "Found USB drive: $driveLetter" "Green"
                
                # Try to format using PowerShell (may work without admin in some cases)
                try {
                    Write-ColorOutput "Attempting to format $driveLetter using PowerShell..." "Blue"
                    Format-Volume -DriveLetter $driveLetter.Replace(':', '') -FileSystem FAT32 -NewFileSystemLabel "SYNTROPY-$NodeName" -Confirm:$false -Force
                    Write-ColorOutput "USB drive formatted successfully using PowerShell" "Green"
                    return $true
                }
                catch {
                    Write-ColorOutput "PowerShell formatting failed: $($_.Exception.Message)" "Red"
                    Write-ColorOutput "Administrative privileges required for USB formatting." "Yellow"
                    return $false
                }
            } else {
                Write-ColorOutput "Could not find USB drive for formatting" "Red"
                return $false
            }
        }
        
        # First, try to unmount any existing partitions
        Write-ColorOutput "Unmounting existing partitions..." "Blue"
        $partitions = Get-WmiObject -Query "ASSOCIATORS OF {Win32_DiskDrive.DeviceID='$DeviceID'} WHERE AssocClass = Win32_DiskDriveToDiskPartition"
        
        foreach ($partition in $partitions) {
            $logicalDisks = Get-WmiObject -Query "ASSOCIATORS OF {Win32_DiskPartition.DeviceID='$($partition.DeviceID)'} WHERE AssocClass = Win32_LogicalDiskToPartition"
            foreach ($disk in $logicalDisks) {
                if ($disk.DeviceID) {
                    Write-ColorOutput "Unmounting drive $($disk.DeviceID)" "Yellow"
                    try {
                        $disk | Invoke-WmiMethod -Name "Eject" -ErrorAction SilentlyContinue
                    } catch {
                        Write-ColorOutput "Could not unmount $($disk.DeviceID): $($_.Exception.Message)" "Yellow"
                    }
                }
            }
        }
        
        # Use diskpart to clean and format the disk
        $diskpartScript = @"
select disk $diskNumber
clean
create partition primary
active
format fs=fat32 quick label="SYNTROPY-$NodeName"
assign
exit
"@
        
        # Write diskpart script to temporary file
        $diskpartFile = [System.IO.Path]::GetTempFileName()
        $diskpartScript | Out-File -FilePath $diskpartFile -Encoding ASCII
        
        Write-ColorOutput "Executing diskpart script..." "Blue"
        Write-ColorOutput "Script content:" "Cyan"
        Write-ColorOutput $diskpartScript "White"
        
        # Execute diskpart
        $result = & diskpart /s $diskpartFile 2>&1
        
        # Clean up
        Remove-Item $diskpartFile -Force
        
        Write-ColorOutput "diskpart output:" "Cyan"
        Write-ColorOutput $result "White"
        
        if ($LASTEXITCODE -eq 0) {
            Write-ColorOutput "USB device formatted successfully" "Green"
            return $true
        } else {
            Write-ColorOutput "diskpart failed with exit code: $LASTEXITCODE" "Red"
            Write-ColorOutput "Output: $result" "Red"
            
            # Try alternative approach using PowerShell
            Write-ColorOutput "Trying alternative formatting approach..." "Yellow"
            return Format-USBDeviceAlternative -DeviceID $DeviceID -NodeName $NodeName
        }
    }
    catch {
        Write-ColorOutput "Error formatting USB device: $($_.Exception.Message)" "Red"
        return $false
    }
}

function Format-USBDeviceAlternative {
    param(
        [string]$DeviceID,
        [string]$NodeName
    )
    
    try {
        Write-ColorOutput "Using alternative formatting approach..." "Blue"
        
        # Get the disk number from DeviceID (handle different formats)
        Write-ColorOutput "Parsing DeviceID: $DeviceID" "Cyan"
        
        # Try different regex patterns to extract disk number
        $diskNumber = $null
        
        # Pattern 1: PhysicalDriveX
        if (-not $diskNumber) {
            $match = [regex]::Match($DeviceID, 'PhysicalDrive(\d+)', [System.Text.RegularExpressions.RegexOptions]::IgnoreCase)
            if ($match.Success) {
                $diskNumber = $match.Groups[1].Value
                Write-ColorOutput "Extracted disk number using Pattern 1: $diskNumber" "Green"
            }
        }
        
        # Pattern 2: \\.\PHYSICALDRIVEX (with escape characters)
        if (-not $diskNumber) {
            $match = [regex]::Match($DeviceID, '\\\\.\\\\(PhysicalDrive\d+)', [System.Text.RegularExpressions.RegexOptions]::IgnoreCase)
            if ($match.Success) {
                $diskNumber = [regex]::Match($match.Groups[1].Value, 'PhysicalDrive(\d+)', [System.Text.RegularExpressions.RegexOptions]::IgnoreCase).Groups[1].Value
                Write-ColorOutput "Extracted disk number using Pattern 2: $diskNumber" "Green"
            }
        }
        
        # Pattern 3: Just extract any number at the end
        if (-not $diskNumber) {
            $match = [regex]::Match($DeviceID, '(\d+)$')
            if ($match.Success) {
                $diskNumber = $match.Groups[1].Value
                Write-ColorOutput "Extracted disk number using Pattern 3: $diskNumber" "Green"
            }
        }
        
        if (-not $diskNumber) {
            throw "Could not extract disk number from DeviceID: $DeviceID"
        }
        
        Write-ColorOutput "Final disk number: $diskNumber" "Cyan"
        
        # Try using PowerShell to format
        $disk = Get-Disk -Number $diskNumber -ErrorAction Stop
        
        Write-ColorOutput "Clearing disk..." "Blue"
        Clear-Disk -Number $diskNumber -RemoveData -Confirm:$false -ErrorAction Stop
        
        Write-ColorOutput "Creating partition..." "Blue"
        $partition = New-Partition -DiskNumber $diskNumber -UseMaximumSize -IsActive -ErrorAction Stop
        
        Write-ColorOutput "Formatting partition..." "Blue"
        Format-Volume -Partition $partition -FileSystem FAT32 -NewFileSystemLabel "SYNTROPY-$NodeName" -Confirm:$false -ErrorAction Stop
        
        Write-ColorOutput "USB device formatted successfully using alternative method" "Green"
        return $true
    }
    catch {
        Write-ColorOutput "Alternative formatting also failed: $($_.Exception.Message)" "Red"
        return $false
    }
}

function Copy-SyntropyFiles {
    param(
        [string]$NodeName,
        [string]$ScriptDir
    )
    
    try {
        Write-ColorOutput "Copying Syntropy files to USB..." "Blue"
        
        # Find the USB drive letter
        $usbDrives = Get-WmiObject Win32_LogicalDisk | Where-Object { $_.DriveType -eq 2 -and $_.VolumeLabel -like "SYNTROPY-*" }
        
        if (-not $usbDrives) {
            throw "Could not find formatted USB drive"
        }
        
        $usbDrive = $usbDrives[0]
        $driveLetter = $usbDrive.DeviceID
        
        Write-ColorOutput "Found USB drive: $driveLetter" "Cyan"
        
        # Create directory structure
        $syntropyDir = "$driveLetter\syntropy"
        $keysDir = "$syntropyDir\keys"
        $configDir = "$syntropyDir\config"
        
        New-Item -ItemType Directory -Path $syntropyDir -Force | Out-Null
        New-Item -ItemType Directory -Path $keysDir -Force | Out-Null
        New-Item -ItemType Directory -Path $configDir -Force | Out-Null
        
        # Generate node configuration
        $nodeConfig = @{
            node_name = $NodeName
            created = (Get-Date).ToString("yyyy-MM-ddTHH:mm:ssZ")
            version = "1.0.0"
            environment = "windows_created"
        } | ConvertTo-Json -Depth 3
        
        $nodeConfig | Out-File -FilePath "$configDir\node.json" -Encoding UTF8
        
        # Generate SSH keys (if OpenSSH is available)
        if (Get-Command ssh-keygen -ErrorAction SilentlyContinue) {
            Write-ColorOutput "Generating SSH keys..." "Blue"
            
            $privateKeyPath = "$keysDir\node.key"
            $publicKeyPath = "$keysDir\node.pub"
            
            # Generate SSH key pair
            & ssh-keygen -t rsa -b 4096 -f $privateKeyPath -N '""' -C "syntropy-node-$NodeName"
            
            if ($LASTEXITCODE -eq 0) {
                Write-ColorOutput "SSH keys generated successfully" "Green"
            } else {
                Write-ColorOutput "SSH key generation failed" "Yellow"
            }
        } else {
            Write-ColorOutput "OpenSSH not available - skipping SSH key generation" "Yellow"
        }
        
        # Copy any additional files from script directory
        if (Test-Path "$ScriptDir\files") {
            Write-ColorOutput "Copying additional files..." "Blue"
            Copy-Item "$ScriptDir\files\*" -Destination $syntropyDir -Recurse -Force
        }
        
        Write-ColorOutput "Files copied successfully to $driveLetter" "Green"
        return $true
    }
    catch {
        Write-ColorOutput "Error copying files: $($_.Exception.Message)" "Red"
        return $false
    }
}

# Check if running with administrative privileges
function Test-Administrator {
    $currentUser = [Security.Principal.WindowsIdentity]::GetCurrent()
    $principal = New-Object Security.Principal.WindowsPrincipal($currentUser)
    return $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

# Request elevation if not running as administrator
function Request-Elevation {
    param([string]$ScriptPath, [string]$Arguments)
    
    Write-ColorOutput "This operation requires administrative privileges." "Yellow"
    Write-ColorOutput "Requesting elevation..." "Blue"
    
    try {
        Start-Process -FilePath "powershell.exe" -ArgumentList "-ExecutionPolicy Bypass -File `"$ScriptPath`" $Arguments" -Verb RunAs -Wait
        return $true
    }
    catch {
        Write-ColorOutput "Failed to request elevation: $($_.Exception.Message)" "Red"
        return $false
    }
}

# Main execution
try {
    Write-ColorOutput "=== Syntropy Cooperative Grid - Windows USB Creator ===" "Cyan"
    Write-ColorOutput "Physical Drive: $PhysicalDrive" "Blue"
    Write-ColorOutput "Node Name: $NodeName" "Blue"
    Write-ColorOutput "Selected Disk Number: $SelectedDiskNum" "Blue"
    
    # Check if running with administrative privileges
    $isAdmin = Test-Administrator
    Write-ColorOutput "Running as Administrator: $isAdmin" "Cyan"
    
    if (-not $isAdmin) {
        Write-ColorOutput "Administrative privileges required for USB formatting." "Yellow"
        Write-ColorOutput "Attempting to continue with limited functionality..." "Yellow"
    }
    
    # Get USB device information
    $usbInfo = Get-USBDeviceInfo -DiskNum $SelectedDiskNum
    
    if (-not $usbInfo) {
        throw "Could not find USB device with disk number: $SelectedDiskNum"
    }
    
    Write-ColorOutput "USB Device Information:" "Cyan"
    Write-ColorOutput "  Device ID: $($usbInfo.DeviceID)" "White"
    Write-ColorOutput "  Model: $($usbInfo.Model)" "White"
    Write-ColorOutput "  Size: $($usbInfo.Size) GB" "White"
    
    # Test device access
    if (-not (Test-USBDeviceAccess -DeviceID $usbInfo.DeviceID)) {
        throw "Cannot access USB device: $($usbInfo.DeviceID)"
    }
    
    # Format the USB device
    if (-not (Format-USBDevice -DeviceID $usbInfo.DeviceID -NodeName $NodeName -IsAdmin $isAdmin)) {
        if (-not $isAdmin) {
            Write-ColorOutput "Formatting failed due to insufficient privileges." "Red"
            Write-ColorOutput "Please run this script as Administrator or use Windows Disk Management." "Yellow"
            Write-ColorOutput "Manual steps:" "Cyan"
            Write-ColorOutput "1. Open Windows Disk Management (diskmgmt.msc)" "White"
            Write-ColorOutput "2. Right-click on the USB drive (Disk $($usbInfo.Index))" "White"
            Write-ColorOutput "3. Select 'Format...'" "White"
            Write-ColorOutput "4. Choose FAT32 file system" "White"
            Write-ColorOutput "5. Set label to 'SYNTROPY-$NodeName'" "White"
            Write-ColorOutput "6. Click OK to format" "White"
            Write-ColorOutput "7. Run this script again after formatting" "White"
            throw "Failed to format USB device - insufficient privileges"
        } else {
            throw "Failed to format USB device"
        }
    }
    
    # Copy Syntropy files (this should work without admin privileges)
    if (-not (Copy-SyntropyFiles -NodeName $NodeName -ScriptDir $ScriptDir)) {
        Write-ColorOutput "Failed to copy Syntropy files to USB" "Red"
        Write-ColorOutput "This may be due to insufficient permissions or USB not being formatted." "Yellow"
        
        # Try to create files in a temporary location and provide instructions
        Write-ColorOutput "Creating Syntropy files in temporary location..." "Blue"
        $tempDir = "$env:TEMP\Syntropy-$NodeName"
        New-Item -ItemType Directory -Path $tempDir -Force | Out-Null
        
        # Create basic node configuration
        $nodeConfig = @{
            node_name = $NodeName
            created = (Get-Date).ToString("yyyy-MM-ddTHH:mm:ssZ")
            version = "1.0.0"
            environment = "windows_created"
            status = "requires_manual_setup"
        } | ConvertTo-Json -Depth 3
        
        $nodeConfig | Out-File -FilePath "$tempDir\node.json" -Encoding UTF8
        
        Write-ColorOutput "Files created in: $tempDir" "Green"
        Write-ColorOutput "Manual setup required:" "Yellow"
        Write-ColorOutput "1. Format the USB drive manually (FAT32, label: SYNTROPY-$NodeName)" "White"
        Write-ColorOutput "2. Copy files from $tempDir to the USB drive" "White"
        Write-ColorOutput "3. Create 'syntropy' folder on USB and copy files there" "White"
        
        throw "Failed to copy Syntropy files - manual setup required"
    }
    
    Write-ColorOutput "=== USB Creation Completed Successfully! ===" "Green"
    Write-ColorOutput "Node Name: $NodeName" "Cyan"
    Write-ColorOutput "USB Drive: $($usbInfo.Model)" "Cyan"
    Write-ColorOutput "Ready for deployment!" "Green"
    
    exit 0
}
catch {
    Write-ColorOutput "Error: $($_.Exception.Message)" "Red"
    exit 1
}
PS_EOF
    
    # Make the script executable and return the path
    chmod +x "$ps_script" 2>/dev/null || true
    echo "$ps_script"
}

# ============================================================================
# MAIN EXECUTION
# ============================================================================

main() {
    show_banner
    
    # Detect environment
    if is_wsl; then
        echo -e "${CYAN}╔════════════════════════════════════════════════════════════════╗${NC}"
        echo -e "${CYAN}║              Running in WSL (Windows Subsystem for Linux)     ║${NC}"
        echo -e "${CYAN}║              USB operations will use Windows devices          ║${NC}"
        echo -e "${CYAN}╚════════════════════════════════════════════════════════════════╝${NC}"
        echo
    fi
    
    # Execute steps
    if ! step1_prerequisites; then
        log ERROR "Prerequisites installation failed"
        exit 1
    fi
    
    echo
    read -p "Press Enter to continue..."
    echo
    
    if ! step2_management; then
        log ERROR "Management setup failed"
        exit 1
    fi
    
    echo
    read -p "Press Enter to continue..."
    echo
    
    if ! step3_detect_usb; then
        log ERROR "USB detection failed"
        exit 1
    fi
    
    echo
    read -p "Press Enter to continue..."
    echo
    
    if ! step4_configure_node; then
        log ERROR "Node configuration failed"
        exit 1
    fi
    
    echo
    log INFO "Starting USB creation process..."
    echo
    
    if ! step5_create_usb; then
        log ERROR "USB creation failed"
        exit 1
    fi
    
    # Success message
    echo
    echo -e "${GREEN}╔════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║                    SETUP COMPLETED SUCCESSFULLY!              ║${NC}"
    echo -e "${GREEN}╚════════════════════════════════════════════════════════════════╝${NC}"
    echo
    echo -e "${CYAN}Next Steps:${NC}"
    echo "1. Safely remove the USB drive"
    echo "2. Insert it into your target hardware"
    echo "3. Boot from USB (may need BIOS/UEFI configuration)"
    echo "4. Wait for automated installation (~30 minutes)"
    echo "5. The node will register automatically when ready"
    echo
    echo -e "${YELLOW}Node Information:${NC}"
    echo "• Node name: $NODE_NAME"
    echo "• Default user: admin"
    echo "• Authentication: SSH key only"
    echo
    echo -e "${GREEN}Your Syntropy node USB is ready!${NC}"
    
    # Load management commands if available
    if [ -f "$HOME/.syntropy/config/syntropy.bashrc" ]; then
        source "$HOME/.syntropy/config/syntropy.bashrc" 2>/dev/null || true
    fi
}

# Run main function
main "$@"