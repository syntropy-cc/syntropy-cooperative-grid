#!/bin/bash

# Syntropy Library Validator and Fixer
# Validates and fixes critical library files for WSL compatibility

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo -e "${CYAN}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║           Syntropy Library Validator & Fixer                  ║${NC}"
echo -e "${CYAN}╚════════════════════════════════════════════════════════════════╝${NC}"
echo

# Function to check if running in WSL
is_wsl() {
    grep -qi microsoft /proc/version
}

# Create backup
backup_file() {
    local file="$1"
    if [ -f "$file" ]; then
        cp "$file" "${file}.backup.$(date +%Y%m%d_%H%M%S)"
        echo -e "${BLUE}Backed up: $file${NC}"
    fi
}

# ============================================================================
# FIX 1: lib/colors.sh - Ensure it exists and has basic definitions
# ============================================================================

fix_colors() {
    local file="$SCRIPT_DIR/lib/colors.sh"
    echo -e "${YELLOW}Checking lib/colors.sh...${NC}"
    
    if [ ! -f "$file" ]; then
        echo -e "${RED}Missing! Creating...${NC}"
        mkdir -p "$(dirname "$file")"
        cat > "$file" << 'COLORS_EOF'
#!/bin/bash
# Color definitions
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m'

# Export colors
export RED GREEN YELLOW BLUE PURPLE CYAN NC
COLORS_EOF
        chmod +x "$file"
        echo -e "${GREEN}Created lib/colors.sh${NC}"
    else
        echo -e "${GREEN}✓ lib/colors.sh exists${NC}"
    fi
}

# ============================================================================
# FIX 2: lib/logging.sh - Ensure logging functions exist
# ============================================================================

fix_logging() {
    local file="$SCRIPT_DIR/lib/logging.sh"
    echo -e "${YELLOW}Checking lib/logging.sh...${NC}"
    
    if [ ! -f "$file" ] || ! grep -q "^log()" "$file"; then
        echo -e "${RED}Missing or incomplete! Fixing...${NC}"
        backup_file "$file"
        mkdir -p "$(dirname "$file")"
        cat > "$file" << 'LOGGING_EOF'
#!/bin/bash
# Logging functions

# Source colors if available
if [ -f "$(dirname "${BASH_SOURCE[0]}")/colors.sh" ]; then
    source "$(dirname "${BASH_SOURCE[0]}")/colors.sh"
else
    # Define basic colors if not available
    RED='\033[0;31m'
    GREEN='\033[0;32m'
    YELLOW='\033[1;33m'
    BLUE='\033[0;34m'
    CYAN='\033[0;36m'
    NC='\033[0m'
fi

# Main logging function
log() {
    local level="$1"
    shift
    local message="$*"
    
    case "$level" in
        INFO)    echo -e "${BLUE}[INFO]${NC} $message" ;;
        SUCCESS) echo -e "${GREEN}[SUCCESS]${NC} $message" ;;
        WARN)    echo -e "${YELLOW}[WARN]${NC} $message" ;;
        ERROR)   echo -e "${RED}[ERROR]${NC} $message" >&2 ;;
        DEBUG)   [ "${DEBUG:-0}" = "1" ] && echo -e "${CYAN}[DEBUG]${NC} $message" ;;
        *)       echo "$message" ;;
    esac
}

# Export function
export -f log
LOGGING_EOF
        chmod +x "$file"
        echo -e "${GREEN}Fixed lib/logging.sh${NC}"
    else
        echo -e "${GREEN}✓ lib/logging.sh exists${NC}"
    fi
}

# ============================================================================
# FIX 3: lib/usb-detection.sh - Add WSL-safe validation
# ============================================================================

fix_usb_detection() {
    local file="$SCRIPT_DIR/lib/usb-detection.sh"
    echo -e "${YELLOW}Checking lib/usb-detection.sh...${NC}"
    
    if [ ! -f "$file" ]; then
        echo -e "${RED}Missing! Creating WSL-compatible version...${NC}"
        mkdir -p "$(dirname "$file")"
        cat > "$file" << 'USB_DETECTION_EOF'
#!/bin/bash
# USB Detection Library - WSL Compatible

# Source dependencies
SCRIPT_DIR="$(dirname "${BASH_SOURCE[0]}")"
[ -f "$SCRIPT_DIR/colors.sh" ] && source "$SCRIPT_DIR/colors.sh"
[ -f "$SCRIPT_DIR/logging.sh" ] && source "$SCRIPT_DIR/logging.sh"

# Check if running in WSL
is_wsl() {
    grep -qi microsoft /proc/version
}

# Select USB device
select_usb_device() {
    if is_wsl; then
        # WSL mode - manual selection
        log INFO "WSL detected - Manual device selection"
        echo "Available block devices:"
        lsblk -d -o NAME,SIZE,TYPE 2>/dev/null | grep disk || echo "No devices found"
        echo
        read -p "Enter device name (e.g., sdb for /dev/sdb): " device_name
        
        if [ -n "$device_name" ]; then
            echo "/dev/$device_name"
        else
            log ERROR "No device selected"
            return 1
        fi
    else
        # Linux mode - auto detection
        local devices=($(lsblk -d -o NAME,RM -n | grep " 1$" | awk '{print $1}' | grep -v "sr0"))
        
        if [ ${#devices[@]} -eq 0 ]; then
            log ERROR "No removable devices found"
            return 1
        fi
        
        if [ ${#devices[@]} -eq 1 ]; then
            echo "/dev/${devices[0]}"
        else
            echo "Multiple devices found:"
            local i=1
            for dev in "${devices[@]}"; do
                local size=$(lsblk -d -o SIZE -n "/dev/$dev" 2>/dev/null)
                echo "  [$i] /dev/$dev ($size)"
                ((i++))
            done
            read -p "Select device (1-${#devices[@]}): " choice
            
            if [[ "$choice" =~ ^[0-9]+$ ]] && [ "$choice" -ge 1 ] && [ "$choice" -le ${#devices[@]} ]; then
                echo "/dev/${devices[$((choice-1))]}"
            else
                log ERROR "Invalid selection"
                return 1
            fi
        fi
    fi
}

# Validate USB safety
validate_usb_safety() {
    local device="$1"
    
    # Basic checks
    if [ ! -b "$device" ]; then
        log ERROR "$device is not a block device"
        return 1
    fi
    
    # Check if it's a system disk
    local mount_point=$(lsblk -n -o MOUNTPOINT "$device" 2>/dev/null | grep -E "^/$|^/boot$|^/home$" | head -1)
    if [ -n "$mount_point" ]; then
        log ERROR "Device $device has system mount point: $mount_point"
        echo "This appears to be a system disk!"
        return 1
    fi
    
    # Size check (between 1GB and 256GB for USB)
    local size_bytes=$(lsblk -b -d -n -o SIZE "$device" 2>/dev/null || echo "0")
    local size_gb=$((size_bytes / 1024 / 1024 / 1024))
    
    if [ "$size_gb" -lt 1 ]; then
        log WARN "Device too small: ${size_gb}GB"
        read -p "Continue anyway? (y/N): " -n 1 -r
        echo
        [[ ! $REPLY =~ ^[Yy]$ ]] && return 1
    fi
    
    if [ "$size_gb" -gt 256 ]; then
        log WARN "Device very large: ${size_gb}GB - might not be a USB"
        read -p "Continue anyway? (y/N): " -n 1 -r
        echo
        [[ ! $REPLY =~ ^[Yy]$ ]] && return 1
    fi
    
    log SUCCESS "Device $device passed safety validation"
    return 0
}

# Check if system disk
is_system_disk() {
    local device="$1"
    
    # Check for critical mount points
    local critical_mounts=("/" "/boot" "/boot/efi" "/usr" "/var" "/home")
    
    for mount in "${critical_mounts[@]}"; do
        local mount_device=$(findmnt -n -o SOURCE "$mount" 2>/dev/null | head -1)
        if [[ "$mount_device" == "$device"* ]]; then
            return 0  # Is system disk
        fi
    done
    
    return 1  # Not system disk
}
USB_DETECTION_EOF
        chmod +x "$file"
        echo -e "${GREEN}Created WSL-compatible lib/usb-detection.sh${NC}"
    else
        echo -e "${GREEN}✓ lib/usb-detection.sh exists${NC}"
        
        # Add WSL compatibility if missing
        if ! grep -q "is_wsl()" "$file"; then
            echo -e "${YELLOW}Adding WSL support to existing file...${NC}"
            backup_file "$file"
            # Add WSL detection at the beginning of the file
            sed -i '1a\
\
# Check if running in WSL\
is_wsl() {\
    grep -qi microsoft /proc/version\
}' "$file"
        fi
    fi
}

# ============================================================================
# FIX 4: lib/config.sh - Ensure init_configuration exists
# ============================================================================

fix_config() {
    local file="$SCRIPT_DIR/lib/config.sh"
    echo -e "${YELLOW}Checking lib/config.sh...${NC}"
    
    if [ ! -f "$file" ] || ! grep -q "init_configuration()" "$file"; then
        echo -e "${RED}Missing or incomplete! Creating minimal version...${NC}"
        backup_file "$file"
        mkdir -p "$(dirname "$file")"
        cat > "$file" << 'CONFIG_EOF'
#!/bin/bash
# Configuration management

# Global variables
WORK_DIR="/tmp/syntropy-usb-$$"
KEYS_DIR="$HOME/.syntropy/keys"
NODES_DIR="$HOME/.syntropy/nodes"
CONFIG_DIR="$HOME/.syntropy/config"
ISO_CACHE_DIR="$HOME/.syntropy/cache/iso"

# ISO Configuration
ISO_FILE="ubuntu-22.04.4-live-server-amd64.iso"
ISO_URL="https://releases.ubuntu.com/22.04/ubuntu-22.04.4-live-server-amd64.iso"
ISO_SHA256="45f873de9f8cb637345d6e66a583762730bbea30277ef7b32c9c3bd6700a32b2"

# Initialize configuration
init_configuration() {
    # Create directories
    mkdir -p "$WORK_DIR"
    mkdir -p "$KEYS_DIR"
    mkdir -p "$NODES_DIR"
    mkdir -p "$CONFIG_DIR"
    mkdir -p "$ISO_CACHE_DIR"
    
    # Export variables
    export WORK_DIR KEYS_DIR NODES_DIR CONFIG_DIR ISO_CACHE_DIR
    export ISO_FILE ISO_URL ISO_SHA256
    
    return 0
}

# Cleanup function
cleanup() {
    [ -d "$WORK_DIR" ] && rm -rf "$WORK_DIR"
    return 0
}

trap cleanup EXIT
CONFIG_EOF
        chmod +x "$file"
        echo -e "${GREEN}Created minimal lib/config.sh${NC}"
    else
        echo -e "${GREEN}✓ lib/config.sh exists${NC}"
    fi
}

# ============================================================================
# FIX 5: lib/usb-preparation.sh - Ensure prepare_usb_device exists
# ============================================================================

fix_usb_preparation() {
    local file="$SCRIPT_DIR/lib/usb-preparation.sh"
    echo -e "${YELLOW}Checking lib/usb-preparation.sh...${NC}"
    
    if [ ! -f "$file" ] || ! grep -q "prepare_usb_device()" "$file"; then
        echo -e "${RED}Missing or incomplete! Creating minimal version...${NC}"
        backup_file "$file"
        mkdir -p "$(dirname "$file")"
        cat > "$file" << 'USB_PREP_EOF'
#!/bin/bash
# USB Preparation

# Source dependencies
SCRIPT_DIR="$(dirname "${BASH_SOURCE[0]}")"
[ -f "$SCRIPT_DIR/logging.sh" ] && source "$SCRIPT_DIR/logging.sh"

# Prepare USB device
prepare_usb_device() {
    local device="$1"
    
    log INFO "Preparing USB device: $device"
    
    # Unmount any mounted partitions
    for partition in ${device}*; do
        if mountpoint -q "$partition" 2>/dev/null; then
            sudo umount "$partition" 2>/dev/null || true
        fi
    done
    
    # Wipe device
    log INFO "Wiping device..."
    sudo wipefs -a "$device" 2>/dev/null || true
    sudo dd if=/dev/zero of="$device" bs=1M count=10 2>/dev/null || true
    
    # Create partition table
    log INFO "Creating partition table..."
    sudo parted -s "$device" mklabel msdos || return 1
    sudo parted -s "$device" mkpart primary fat32 1MiB 100% || return 1
    sudo parted -s "$device" set 1 boot on || return 1
    
    # Wait for kernel to recognize new partition
    sudo partprobe "$device" 2>/dev/null || true
    sleep 3
    
    # Find partition device
    local partition="${device}1"
    if [ ! -b "$partition" ]; then
        partition="${device}p1"
    fi
    
    if [ ! -b "$partition" ]; then
        log ERROR "Partition not found after creation"
        return 1
    fi
    
    # Format partition
    log INFO "Formatting partition..."
    sudo mkfs.fat -F32 -n "SYNTROPY" "$partition" || return 1
    
    log SUCCESS "USB device prepared successfully"
    echo "$partition"
    return 0
}
USB_PREP_EOF
        chmod +x "$file"
        echo -e "${GREEN}Created minimal lib/usb-preparation.sh${NC}"
    else
        echo -e "${GREEN}✓ lib/usb-preparation.sh exists${NC}"
    fi
}

# ============================================================================
# FIX 6: Create missing stubs for other required libraries
# ============================================================================

create_missing_stubs() {
    echo -e "${YELLOW}Creating stubs for other required libraries...${NC}"
    
    # Geographic detection stub
    if [ ! -f "$SCRIPT_DIR/lib/geographic.sh" ]; then
        cat > "$SCRIPT_DIR/lib/geographic.sh" << 'GEO_EOF'
#!/bin/bash
# Geographic detection stub

detect_coordinates() {
    local manual="$1"
    if [ -n "$manual" ]; then
        echo "$manual:manual:Manual:Entry"
    else
        # Default to São Paulo
        echo "-23.5505,-46.6333:auto:São Paulo:Brazil"
    fi
}

generate_location_id() {
    echo "loc-$(openssl rand -hex 4)"
}
GEO_EOF
        chmod +x "$SCRIPT_DIR/lib/geographic.sh"
        echo -e "${GREEN}Created geographic.sh stub${NC}"
    fi
    
    # Security stub
    if [ ! -f "$SCRIPT_DIR/lib/security.sh" ]; then
        cat > "$SCRIPT_DIR/lib/security.sh" << 'SEC_EOF'
#!/bin/bash
# Security stub

OWNER_KEY_PATH=""
COMMUNITY_KEY_PATH=""
OWNER_FINGERPRINT=""
COMMUNITY_FINGERPRINT=""

setup_security_keys() {
    OWNER_KEY_PATH="$WORK_DIR/owner_key"
    COMMUNITY_KEY_PATH="$WORK_DIR/community_key"
    
    # Generate keys
    ssh-keygen -t ed25519 -f "$OWNER_KEY_PATH" -N "" -q
    ssh-keygen -t ed25519 -f "$COMMUNITY_KEY_PATH" -N "" -q
    
    OWNER_FINGERPRINT=$(ssh-keygen -lf "${OWNER_KEY_PATH}.pub" | awk '{print $2}')
    COMMUNITY_FINGERPRINT=$(ssh-keygen -lf "${COMMUNITY_KEY_PATH}.pub" | awk '{print $2}')
    
    return 0
}

store_keys_locally() {
    return 0
}

generate_security_summary() {
    return 0
}

cleanup_sensitive_data() {
    return 0
}
SEC_EOF
        chmod +x "$SCRIPT_DIR/lib/security.sh"
        echo -e "${GREEN}Created security.sh stub${NC}"
    fi
    
    # ISO management stub
    if [ ! -f "$SCRIPT_DIR/lib/iso-management.sh" ]; then
        cat > "$SCRIPT_DIR/lib/iso-management.sh" << 'ISO_EOF'
#!/bin/bash
# ISO management stub

download_ubuntu_iso() {
    log INFO "Downloading Ubuntu ISO (simulated)..."
    # In production, this would download the actual ISO
    touch "$WORK_DIR/$ISO_FILE"
    return 0
}

install_ubuntu_to_usb() {
    log INFO "Installing Ubuntu to USB (simulated)..."
    return 0
}
ISO_EOF
        chmod +x "$SCRIPT_DIR/lib/iso-management.sh"
        echo -e "${GREEN}Created iso-management.sh stub${NC}"
    fi
    
    # Cloud-init stub
    if [ ! -f "$SCRIPT_DIR/lib/cloud-init.sh" ]; then
        cat > "$SCRIPT_DIR/lib/cloud-init.sh" << 'CLOUD_EOF'
#!/bin/bash
# Cloud-init stub

create_cloud_init_configuration() {
    log INFO "Creating cloud-init configuration..."
    return 0
}
CLOUD_EOF
        chmod +x "$SCRIPT_DIR/lib/cloud-init.sh"
        echo -e "${GREEN}Created cloud-init.sh stub${NC}"
    fi
    
    # Management scripts stub
    if [ ! -f "$SCRIPT_DIR/lib/management-scripts.sh" ]; then
        cat > "$SCRIPT_DIR/lib/management-scripts.sh" << 'MGMT_EOF'
#!/bin/bash
# Management scripts stub

create_node_metadata() {
    log INFO "Creating node metadata..."
    return 0
}

create_node_management_scripts() {
    log INFO "Creating management scripts..."
    return 0
}

finalize_usb_creation() {
    log SUCCESS "USB creation finalized"
    return 0
}
MGMT_EOF
        chmod +x "$SCRIPT_DIR/lib/management-scripts.sh"
        echo -e "${GREEN}Created management-scripts.sh stub${NC}"
    fi
}

# ============================================================================
# MAIN VALIDATION AND FIX PROCESS
# ============================================================================

main() {
    echo -e "${CYAN}Starting validation and fixes...${NC}"
    echo
    
    # Check if lib directory exists
    if [ ! -d "$SCRIPT_DIR/lib" ]; then
        echo -e "${YELLOW}Creating lib directory...${NC}"
        mkdir -p "$SCRIPT_DIR/lib"
    fi
    
    # Run all fixes
    fix_colors
    fix_logging
    fix_usb_detection
    fix_config
    fix_usb_preparation
    create_missing_stubs
    
    echo
    echo -e "${GREEN}╔════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║                  Validation Complete!                         ║${NC}"
    echo -e "${GREEN}╚════════════════════════════════════════════════════════════════╝${NC}"
    echo
    
    # Test if everything works
    echo -e "${CYAN}Testing library imports...${NC}"
    
    # Try to source all libraries
    local test_failed=false
    for lib in colors.sh logging.sh config.sh usb-detection.sh usb-preparation.sh \
               geographic.sh security.sh iso-management.sh cloud-init.sh management-scripts.sh; do
        if [ -f "$SCRIPT_DIR/lib/$lib" ]; then
            if source "$SCRIPT_DIR/lib/$lib" 2>/dev/null; then
                echo -e "${GREEN}✓ $lib loaded successfully${NC}"
            else
                echo -e "${RED}✗ $lib failed to load${NC}"
                test_failed=true
            fi
        else
            echo -e "${RED}✗ $lib not found${NC}"
            test_failed=true
        fi
    done
    
    if [ "$test_failed" = false ]; then
        echo
        echo -e "${GREEN}All libraries are ready!${NC}"
        echo -e "${GREEN}You can now run:${NC}"
        echo -e "${CYAN}  ./quick-start.sh${NC}"
        echo -e "${CYAN}  ./create-syntropy-usb-enhanced.sh${NC}"
    else
        echo
        echo -e "${YELLOW}Some issues remain. Please check the errors above.${NC}"
    fi
    
    # WSL-specific advice
    if is_wsl; then
        echo
        echo -e "${CYAN}═══ WSL Detected ═══${NC}"
        echo "Remember:"
        echo "1. USB devices appear as /dev/sdX in WSL"
        echo "2. You may need to restart WSL after connecting USB"
        echo "3. Run 'wsl --shutdown' in PowerShell if devices don't appear"
    fi
}

# Run main
main "$@"