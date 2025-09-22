#!/bin/bash

# Syntropy Cooperative Grid - Enhanced USB Creator (Main Script)
# Version: 2.0.0

set -e

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Check for required dependencies
check_dependencies() {
    local missing_deps=()
    
    # Required commands
    local required_commands=(
        "dd"        # For writing to USB
        "mkfs.fat"  # For formatting USB
        "parted"    # For partition management
        "openssl"   # For key generation
        "curl"      # For downloads
        "jq"        # For JSON processing
        "lsblk"     # For disk detection
        "sudo"      # For elevated privileges
    )
    
    # Check each command
    for cmd in "${required_commands[@]}"; do
        if ! command -v "$cmd" &> /dev/null; then
            missing_deps+=("$cmd")
        fi
    done
    
    # Check for sudo access
    if ! sudo -n true 2>/dev/null; then
        echo -e "${RED}Error: This script requires sudo privileges${NC}"
        echo "Please ensure you have sudo access and try again"
        exit 1
    fi
    
    # If we found missing dependencies, error out with instructions
    if [ ${#missing_deps[@]} -ne 0 ]; then
        echo -e "${RED}Error: Missing required dependencies${NC}"
        echo "Please install the following packages:"
        printf '%s\n' "${missing_deps[@]}"
        
        # Provide installation instructions based on common package managers
        echo -e "\nOn Ubuntu/Debian:"
        echo "sudo apt-get install ${missing_deps[*]}"
        echo -e "\nOn RHEL/CentOS:"
        echo "sudo yum install ${missing_deps[*]}"
        exit 1
    fi
    
    # Check for minimum disk space (1GB in /tmp)
    local available_space=$(df -BG /tmp | awk 'NR==2 {print $4}' | sed 's/G//')
    if [ "$available_space" -lt 1 ]; then
        echo -e "${RED}Error: Insufficient disk space${NC}"
        echo "At least 1GB of free space is required in /tmp"
        exit 1
    fi
}

# Source all modules
source "$SCRIPT_DIR/lib/colors.sh"
source "$SCRIPT_DIR/lib/logging.sh"
source "$SCRIPT_DIR/lib/config.sh"
source "$SCRIPT_DIR/lib/usb-detection.sh"
source "$SCRIPT_DIR/lib/wsl-usb-detection.sh"  # Add WSL support
source "$SCRIPT_DIR/lib/geographic.sh"
source "$SCRIPT_DIR/lib/security.sh"
source "$SCRIPT_DIR/lib/iso-management.sh"
source "$SCRIPT_DIR/lib/usb-preparation.sh"
source "$SCRIPT_DIR/lib/cloud-init.sh"
source "$SCRIPT_DIR/lib/management-scripts.sh"

# Check if running in WSL
WSL_MODE=false
if is_wsl; then
    WSL_MODE=true
    log INFO "Running in WSL mode - using Windows USB detection"
fi
source "$SCRIPT_DIR/lib/progress.sh"

# Check dependencies first
check_dependencies

# Main banner
show_banner() {
    echo -e "${PURPLE}"
    cat << 'EOF'
╔════════════════════════════════════════════════════════════════════════════╗
║                    SYNTROPY COOPERATIVE GRID                              ║
║                Enhanced USB Creator with Node Management                  ║
║                          Version 2.0.0                                   ║
║                                                                           ║
║  Creates production-ready USB with automated node management setup       ║
╚════════════════════════════════════════════════════════════════════════════╝
EOF
    echo -e "${NC}"
}

# Show help function
show_help() {
    echo "Usage: $0 [usb_device] [options]"
    echo ""
    echo "Options:"
    echo "  --owner-key <file>        Use existing owner private key (enables multi-node management)"
    echo "  --node-name <name>        Custom node name (default: auto-generated from location)"
    echo "  --description <desc>      Node description"
    echo "  --coordinates <lat,lon>   Manual coordinates override (default: auto-detected)"
    echo "  --auto-detect             Automatically detect and select USB device"
    echo "  --help                    Show this help"
    echo ""
    echo "USB Device:"
    echo "  If no USB device is specified, the script will automatically detect and"
    echo "  prompt you to select from available USB storage devices."
    echo ""
    echo "Examples:"
    echo "  $0                                                       # Auto-detect USB"
    echo "  $0 --auto-detect --node-name home-server-01             # Auto-detect with custom name"
    echo "  $0 /dev/sdb --node-name home-server-01                  # Specify USB device"
    echo "  $0 --owner-key ~/.syntropy/keys/main_owner.key          # Use existing key"
    echo "  $0 --coordinates \"-23.5505,-46.6333\"                   # Manual coordinates"
    echo ""
    echo "Multi-node workflow:"
    echo "  1. First node:  $0 --node-name main-server"
    echo "  2. Second node: $0 --owner-key ~/.syntropy/keys/main-server_owner.key --node-name edge-01"
    echo "  3. Third node:  $0 --owner-key ~/.syntropy/keys/main-server_owner.key --node-name edge-02"
    echo ""
    echo "Safety features:"
    echo "  - Automatic USB device detection and validation"
    echo "  - Multiple confirmation prompts for destructive operations"
    echo "  - Device size and type verification"
    echo "  - Protection against overwriting system disks"
    echo "  - Comprehensive error handling and recovery"
}

# Parse command line arguments
parse_arguments() {
    OWNER_KEY_FILE=""
    NODE_NAME=""
    NODE_DESCRIPTION=""
    MANUAL_COORDINATES=""
    AUTO_DETECT=false
    USB_DEVICE=""

    # If first argument is not a flag and looks like a device, use it
    if [ $# -gt 0 ] && [[ "$1" =~ ^/dev/ ]]; then
        USB_DEVICE="$1"
        shift
    fi

    while [[ $# -gt 0 ]]; do
        case $1 in
            --owner-key)
                OWNER_KEY_FILE="$2"
                shift 2
                ;;
            --node-name)
                NODE_NAME="$2"
                shift 2
                ;;
            --description)
                NODE_DESCRIPTION="$2"
                shift 2
                ;;
            --coordinates)
                MANUAL_COORDINATES="$2"
                shift 2
                ;;
            --auto-detect)
                AUTO_DETECT=true
                shift
                ;;
            --help)
                show_help
                exit 0
                ;;
            /dev/*)
                if [ -z "$USB_DEVICE" ]; then
                    USB_DEVICE="$1"
                fi
                shift
                ;;
            *)
                log ERROR "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
}

# Main function
main() {
    show_banner
    
    # Parse arguments
    parse_arguments "$@"
    
    # Initialize configuration
    init_configuration
    
    # Create work directory
    mkdir -p "$WORK_DIR"
    cd "$WORK_DIR"
    
    # USB device detection and validation
    if [ -z "$USB_DEVICE" ] || [ "$AUTO_DETECT" = true ]; then
        log INFO "Auto-detecting USB devices..."
        USB_DEVICE=$(select_usb_device)
    fi

    # Para WSL, validar o dispositivo físico do Windows
    if [ "$WSL_MODE" = true ]; then
        if [[ ! "$USB_DEVICE" =~ ^PhysicalDrive[0-9]+$ ]] && [[ ! "$USB_DEVICE" =~ ^\\\\\.\\PhysicalDrive[0-9]+$ ]]; then
            log ERROR "Invalid Windows physical device format: $USB_DEVICE"
            log ERROR "Expected format: PhysicalDrive# or \\\\.\\PhysicalDrive#"
            exit 1
        fi
        # Converter formato Windows para formato WSL
        DISK_NUMBER=$(echo "$USB_DEVICE" | grep -o '[0-9]\+')
        if [ -n "$DISK_NUMBER" ]; then
            USB_DEVICE="/dev/sd$(printf "\\$(printf '%03o' $((97 + $DISK_NUMBER)))")"
            log INFO "Converted Windows device $USB_DEVICE to WSL device $USB_DEVICE"
        else
            log ERROR "Failed to extract disk number from device path"
            exit 1
        fi
    fi

    if [ ! -b "$USB_DEVICE" ]; then
        log ERROR "$USB_DEVICE is not a valid block device"
        echo "Available devices:"
        lsblk | grep disk
        exit 1
    fi

    # Perform comprehensive safety validation
    validate_usb_safety "$USB_DEVICE"

    # Final confirmation
    show_final_confirmation

    # Execute main workflow
    execute_usb_creation_workflow
}

# Show final confirmation
show_final_confirmation() {
    echo ""
    echo -e "${YELLOW}WARNING: This will completely erase $USB_DEVICE${NC}"
    echo "Device information:"
    lsblk "$USB_DEVICE"
    echo ""
    echo -e "${CYAN}Node Configuration:${NC}"
    echo "• Owner Key: ${OWNER_KEY_FILE:-"Will be generated"}"
    echo "• Node Name: ${NODE_NAME:-"Auto-generated from location"}"
    echo "• Description: ${NODE_DESCRIPTION:-"Auto-generated"}"
    echo "• Coordinates: ${MANUAL_COORDINATES:-"Auto-detected"}"
    echo ""
    read -p "Continue with USB creation? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log INFO "Operation cancelled by user."
        exit 0
    fi
}

# Execute the main USB creation workflow
execute_usb_creation_workflow() {
    log INFO "[1/8] Detecting location and generating node identity..."
    
    # Enhanced coordinate detection
    local coords_with_info=$(detect_coordinates "$MANUAL_COORDINATES")
    local coordinates=$(echo "$coords_with_info" | cut -d':' -f1)
    local detection_method=$(echo "$coords_with_info" | cut -d':' -f2)
    local detected_city=$(echo "$coords_with_info" | cut -d':' -f3)
    local detected_country=$(echo "$coords_with_info" | cut -d':' -f4)
    
    echo "Location detected:"
    echo "  Coordinates: $coordinates"
    echo "  City: $detected_city"
    echo "  Country: $detected_country"
    echo "  Method: $detection_method"
    
    # Generate location-based node ID
    local location_node_id=$(generate_location_id "$coords_with_info")
    
    log INFO "[2/8] Setting up security keys..."
    
    # Enhanced key management
    setup_security_keys "$OWNER_KEY_FILE" "$location_node_id"
    
    # Generate node name if not provided
    if [ -z "$NODE_NAME" ]; then
        NODE_NAME="syntropy-$location_node_id"
    fi
    
    # Set default description
    if [ -z "$NODE_DESCRIPTION" ]; then
        NODE_DESCRIPTION="Syntropy Cooperative Grid Node in $detected_city, $detected_country ($coordinates)"
    fi
    
    log INFO "[3/8] Creating comprehensive node metadata..."
    
    create_node_metadata "$NODE_NAME" "$location_node_id" "$coordinates" \
        "$detection_method" "$detected_city" "$detected_country" "$NODE_DESCRIPTION"
    
    log INFO "[4/8] Downloading Ubuntu Server ISO..."
    
    if ! download_ubuntu_iso; then
        log ERROR "Failed to download Ubuntu ISO"
        exit 1
    fi
    
    log INFO "[5/8] Preparing USB device..."
    
    local usb_partition=$(prepare_usb_device "$USB_DEVICE")
    if [ $? -ne 0 ]; then
        log ERROR "Failed to prepare USB device"
        exit 1
    fi
    
    log INFO "[6/8] Installing Ubuntu and creating configuration..."
    
    install_ubuntu_to_usb "$usb_partition"
    
    log INFO "[7/8] Creating cloud-init configuration..."
    
    create_cloud_init_configuration "$usb_partition" "$NODE_NAME" "$location_node_id" \
        "$coordinates" "$detection_method" "$detected_city" "$detected_country" "$NODE_DESCRIPTION"
    
    log INFO "[8/8] Creating management scripts and finalizing..."
    
    create_node_management_scripts "$NODE_NAME" "$location_node_id" "$coordinates" \
        "$detected_city" "$detected_country" "$NODE_DESCRIPTION"
    
    finalize_usb_creation "$NODE_NAME" "$location_node_id" "$coordinates" \
        "$detected_city" "$detected_country"
}

# Run main function with all arguments
main "$@"