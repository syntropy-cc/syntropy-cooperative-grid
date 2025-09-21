#!/bin/bash

# Syntropy Cooperative Grid - Configuration Management
# Version: 2.0.0

# Default configuration values
DEFAULT_NODES_DIR="$HOME/.syntropy/nodes"
DEFAULT_KEYS_DIR="$HOME/.syntropy/keys"
DEFAULT_CONFIG_DIR="$HOME/.syntropy/config"
DEFAULT_ISO_CACHE_DIR="$HOME/.syntropy/cache/iso"
DEFAULT_WORK_DIR="/tmp/syntropy-usb-enhanced-$$"

# Ubuntu ISO Configuration
ISO_FILE="ubuntu-22.04.4-live-server-amd64.iso"
ISO_URL="https://releases.ubuntu.com/22.04/ubuntu-22.04.4-live-server-amd64.iso"
ISO_SHA256="45f873de9f8cb637345d6e66a583762730bbea30277ef7b32c9c3bd6700a32b2"

# Global configuration variables
NODES_DIR=""
KEYS_DIR=""
CONFIG_DIR=""
ISO_CACHE_DIR=""
WORK_DIR=""
USB_MOUNT=""

# Initialize configuration directories and variables
init_configuration() {
    log INFO "Initializing Syntropy configuration..."
    
    # Set directory paths
    NODES_DIR="${SYNTROPY_NODES_DIR:-$DEFAULT_NODES_DIR}"
    KEYS_DIR="${SYNTROPY_KEYS_DIR:-$DEFAULT_KEYS_DIR}"
    CONFIG_DIR="${SYNTROPY_CONFIG_DIR:-$DEFAULT_CONFIG_DIR}"
    ISO_CACHE_DIR="${SYNTROPY_ISO_CACHE_DIR:-$DEFAULT_ISO_CACHE_DIR}"
    WORK_DIR="${SYNTROPY_WORK_DIR:-$DEFAULT_WORK_DIR}"
    
    # Create all necessary directories
    create_directory_structure
    
    # Initialize logging
    init_logging
    
    # Validate system requirements
    if ! validate_system_requirements; then
        log ERROR "System requirements validation failed"
        return 1
    fi
    
    # Set up cleanup trap
    setup_cleanup_trap
    
    log SUCCESS "Configuration initialized successfully"
    log DEBUG "Nodes directory: $NODES_DIR"
    log DEBUG "Keys directory: $KEYS_DIR"
    log DEBUG "Config directory: $CONFIG_DIR"
    log DEBUG "ISO cache directory: $ISO_CACHE_DIR"
    log DEBUG "Work directory: $WORK_DIR"
    
    return 0
}

# Create directory structure
create_directory_structure() {
    local directories=(
        "$NODES_DIR"
        "$KEYS_DIR"
        "$CONFIG_DIR"
        "$ISO_CACHE_DIR"
        "$CONFIG_DIR/templates"
        "$CONFIG_DIR/backups"
        "$CONFIG_DIR/removed"
    )
    
    for dir in "${directories[@]}"; do
        if [ ! -d "$dir" ]; then
            if ! mkdir -p "$dir"; then
                log ERROR "Failed to create directory: $dir"
                return 1
            fi
            log DEBUG "Created directory: $dir"
        fi
    done
    
    return 0
}

# Cleanup function for graceful exit
cleanup() {
    local exit_code=$?
    
    log INFO "Starting cleanup process..."
    
    if [ -d "$WORK_DIR" ]; then
        log INFO "Cleaning up temporary directory: $WORK_DIR"
        cd /
        rm -rf "$WORK_DIR"
    fi
    
    # Unmount any mounted USBs
    if [ -n "${USB_MOUNT:-}" ] && mountpoint -q "$USB_MOUNT" 2>/dev/null; then
        log INFO "Unmounting USB device"
        sudo umount "$USB_MOUNT" 2>/dev/null || true
        sudo rmdir "$USB_MOUNT" 2>/dev/null || true
    fi
    
    # Cleanup any remaining ISO mounts
    cleanup_iso_mounts 2>/dev/null || true
    
    if [ $exit_code -ne 0 ]; then
        log ERROR "Script failed with exit code $exit_code"
    else
        log SUCCESS "Script completed successfully"
    fi
    
    exit $exit_code
}

# Set up cleanup trap
setup_cleanup_trap() {
    trap cleanup EXIT
    trap 'cleanup; exit 130' INT  # Ctrl+C
    trap 'cleanup; exit 143' TERM # Termination
}

# Validate required tools
validate_dependencies() {
    local required_tools=(
        "curl"
        "wget"
        "python3"
        "ssh-keygen"
        "openssl"
        "lsblk"
        "mount"
        "umount"
        "parted"
        "mkfs.fat"
        "wipefs"
        "sgdisk"
        "nmap"
        "jq"
    )
    
    local missing_tools=()
    local missing_packages=()
    
    for tool in "${required_tools[@]}"; do
        if ! command -v "$tool" >/dev/null 2>&1; then
            missing_tools+=("$tool")
            # Map tool names to actual package names
            case "$tool" in
                "mkfs.fat") missing_packages+=("dosfstools") ;;
                "sgdisk") missing_packages+=("gdisk") ;;
                "ssh-keygen") missing_packages+=("openssh-client") ;;
                *) missing_packages+=("$tool") ;;
            esac
        fi
    done
    
    if [ ${#missing_tools[@]} -gt 0 ]; then
        log WARN "Missing required tools: ${missing_tools[*]}"
        echo ""
        echo "The following packages need to be installed:"
        
        # Show package mapping
        local i=0
        for tool in "${missing_tools[@]}"; do
            echo "  - ${missing_packages[$i]} (provides: $tool)"
            ((i++))
        done
        
        echo ""
        
        # Auto-install for supported package managers
        if command -v apt-get >/dev/null 2>&1; then
            local install_cmd="sudo apt-get update && sudo apt-get install -y ${missing_packages[*]}"
            echo "Ubuntu/Debian detected. Installation command:"
            echo "  $install_cmd"
            echo ""
            read -p "Install missing packages automatically? (y/N): " -n 1 -r
            echo
            
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                log INFO "Installing missing packages..."
                
                if sudo apt-get update; then
                    if sudo apt-get install -y "${missing_packages[@]}"; then
                        log SUCCESS "All packages installed successfully"
                        
                        # Verify installation
                        local still_missing=()
                        for tool in "${missing_tools[@]}"; do
                            if ! command -v "$tool" >/dev/null 2>&1; then
                                still_missing+=("$tool")
                            fi
                        done
                        
                        if [ ${#still_missing[@]} -gt 0 ]; then
                            log ERROR "Some tools are still missing after installation: ${still_missing[*]}"
                            echo "Please check your package manager or install manually"
                            return 1
                        fi
                        
                        log SUCCESS "All required dependencies are now available"
                        return 0
                    else
                        log ERROR "Package installation failed"
                        echo "Please install manually or check your package manager"
                        return 1
                    fi
                else
                    log ERROR "Failed to update package lists"
                    return 1
                fi
            else
                log INFO "Automatic installation declined"
                echo ""
                echo "Please install manually with:"
                echo "  $install_cmd"
                return 1
            fi
            
        elif command -v yum >/dev/null 2>&1; then
            local install_cmd="sudo yum install -y ${missing_packages[*]}"
            echo "RHEL/CentOS detected. Installation command:"
            echo "  $install_cmd"
            echo ""
            read -p "Install missing packages automatically? (y/N): " -n 1 -r
            echo
            
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                log INFO "Installing missing packages..."
                if sudo yum install -y "${missing_packages[@]}"; then
                    log SUCCESS "All packages installed successfully"
                    return 0
                else
                    log ERROR "Package installation failed"
                    return 1
                fi
            else
                echo "Please install manually with: $install_cmd"
                return 1
            fi
            
        elif command -v dnf >/dev/null 2>&1; then
            local install_cmd="sudo dnf install -y ${missing_packages[*]}"
            echo "Fedora detected. Installation command:"
            echo "  $install_cmd"
            echo ""
            read -p "Install missing packages automatically? (y/N): " -n 1 -r
            echo
            
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                log INFO "Installing missing packages..."
                if sudo dnf install -y "${missing_packages[@]}"; then
                    log SUCCESS "All packages installed successfully"
                    return 0
                else
                    log ERROR "Package installation failed"
                    return 1
                fi
            else
                echo "Please install manually with: $install_cmd"
                return 1
            fi
            
        elif command -v brew >/dev/null 2>&1; then
            local install_cmd="brew install ${missing_packages[*]}"
            echo "macOS (Homebrew) detected. Installation command:"
            echo "  $install_cmd"
            echo ""
            read -p "Install missing packages automatically? (y/N): " -n 1 -r
            echo
            
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                log INFO "Installing missing packages..."
                if brew install "${missing_packages[@]}"; then
                    log SUCCESS "All packages installed successfully"
                    return 0
                else
                    log ERROR "Package installation failed"
                    return 1
                fi
            else
                echo "Please install manually with: $install_cmd"
                return 1
            fi
            
        else
            echo "Unsupported package manager. Please install the following packages manually:"
            for package in "${missing_packages[@]}"; do
                echo "  - $package"
            done
            return 1
        fi
    fi
    
    log SUCCESS "All required dependencies are available"
    return 0
}

# Check if running as root (should not be)
check_root_user() {
    if [ "$EUID" -eq 0 ]; then
        log ERROR "This script should not be run as root"
        echo ""
        echo "Running as root can be dangerous and is not necessary."
        echo "The script will use sudo only when needed for specific operations."
        echo "Please run as a regular user."
        echo ""
        return 1
    fi
    
    return 0
}

# Check sudo access
check_sudo_access() {
    if ! command -v sudo >/dev/null 2>&1; then
        log ERROR "sudo is required but not available"
        echo ""
        echo "This script requires sudo for USB device operations."
        echo "Please install sudo or run on a system where it's available."
        echo ""
        return 1
    fi
    
    # Test sudo access
    if ! sudo -n true 2>/dev/null; then
        log INFO "Testing sudo access..."
        if ! sudo true; then
            log ERROR "sudo access required for USB operations"
            echo ""
            echo "This script needs sudo privileges to:"
            echo "- Access block devices"
            echo "- Mount/unmount filesystems"
            echo "- Write to USB devices"
            echo ""
            return 1
        fi
    fi
    
    log SUCCESS "sudo access confirmed"
    return 0
}

# Validate system requirements
validate_system_requirements() {
    log INFO "Validating system requirements..."
    
    # Check if not running as root
    if ! check_root_user; then
        return 1
    fi
    
    # Check sudo access
    if ! check_sudo_access; then
        return 1
    fi
    
    # Check required tools
    if ! validate_dependencies; then
        return 1
    fi
    
    # Check disk space
    local available_space=$(df "$HOME" | tail -1 | awk '{print $4}')
    local required_space=$((2 * 1024 * 1024))  # 2GB in KB
    
    if [ "$available_space" -lt "$required_space" ]; then
        log ERROR "Insufficient disk space"
        echo "Available: $(($available_space / 1024 / 1024))GB"
        echo "Required: 2GB minimum"
        return 1
    fi
    
    # Check if Python 3 has required modules
    if ! python3 -c "import json, sys" 2>/dev/null; then
        log ERROR "Python 3 json module not available"
        return 1
    fi
    
    # Check internet connectivity
    if ! timeout 10 curl -s "http://releases.ubuntu.com" >/dev/null 2>&1; then
        log WARN "Internet connectivity check failed"
        echo "Internet access is required to download Ubuntu ISO"
        echo "Please check your internet connection"
    fi
    
    log SUCCESS "System requirements validated"
    return 0
}

# Get configuration value with fallback
get_config() {
    local key="$1"
    local default="$2"
    local config_file="$CONFIG_DIR/syntropy.conf"
    
    if [ -f "$config_file" ]; then
        local value=$(grep "^$key=" "$config_file" 2>/dev/null | cut -d'=' -f2- | tr -d '"'"'")
        if [ -n "$value" ]; then
            echo "$value"
            return 0
        fi
    fi
    
    echo "$default"
}

# Set configuration value
set_config() {
    local key="$1"
    local value="$2"
    local config_file="$CONFIG_DIR/syntropy.conf"
    
    # Create config file if it doesn't exist
    if [ ! -f "$config_file" ]; then
        cat > "$config_file" << 'EOF'
# Syntropy Cooperative Grid Configuration
# This file contains local configuration overrides

EOF
    fi
    
    # Update or add configuration value
    if grep -q "^$key=" "$config_file" 2>/dev/null; then
        # Update existing value
        sed -i "s/^$key=.*/$key=\"$value\"/" "$config_file"
    else
        # Add new value
        echo "$key=\"$value\"" >> "$config_file"
    fi
    
    log DEBUG "Configuration updated: $key=$value"
}