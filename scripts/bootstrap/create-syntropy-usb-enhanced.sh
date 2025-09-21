#!/bin/bash

# Syntropy Cooperative Grid - Enhanced USB Creator with Advanced Node Management
# Creates production-ready USB with automated node management setup
# Version: 2.0.0

set -e

# Configuration
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
WORK_DIR="/tmp/syntropy-usb-enhanced-$$"
NODES_DIR="$HOME/.syntropy/nodes"
KEYS_DIR="$HOME/.syntropy/keys"
CONFIG_DIR="$HOME/.syntropy/config"
ISO_CACHE_DIR="$HOME/.syntropy/cache/iso"
TIMESTAMP=$(date +%s)

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m'

# Ubuntu ISO Configuration
ISO_FILE="ubuntu-22.04.4-live-server-amd64.iso"
ISO_URL="https://releases.ubuntu.com/22.04/ubuntu-22.04.4-live-server-amd64.iso"
ISO_SHA256="45f873de9f8cb637345d6e66a583762730bbea30277ef7b32c9c3bd6700a32b2"

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

# Logging function
log() {
    local level="$1"
    shift
    local message="$*"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    
    case "$level" in
        INFO)  echo -e "${BLUE}[INFO]${NC} $message" ;;
        WARN)  echo -e "${YELLOW}[WARN]${NC} $message" ;;
        ERROR) echo -e "${RED}[ERROR]${NC} $message" ;;
        SUCCESS) echo -e "${GREEN}[SUCCESS]${NC} $message" ;;
        *) echo "[$level] $message" ;;
    esac
}

# Error handling
cleanup() {
    local exit_code=$?
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
    
    if [ $exit_code -ne 0 ]; then
        log ERROR "Script failed with exit code $exit_code"
    fi
    exit $exit_code
}

trap cleanup EXIT

# Enhanced USB device detection
detect_usb_devices() {
    log INFO "Scanning for USB storage devices..."
    
    local usb_devices=()
    
    # Find removable devices that are connected via USB
    while IFS= read -r line; do
        local device=$(echo "$line" | awk '{print $1}')
        local size=$(echo "$line" | awk '{print $2}')
        local type=$(echo "$line" | awk '{print $3}')
        local removable=$(echo "$line" | awk '{print $4}')
        local model=$(echo "$line" | awk '{for(i=5;i<=NF;i++) printf "%s ", $i; print ""}' | sed 's/[[:space:]]*$//')
        
        # Only include removable disks
        if [ "$removable" = "1" ] && [ "$type" = "disk" ]; then
            # Verify it's actually USB by checking the device path
            local device_name=$(basename "/dev/$device")
            if [ -e "/sys/block/$device_name" ]; then
                local device_path=$(readlink -f "/sys/block/$device_name")
                if [[ "$device_path" == *"/usb"* ]]; then
                    usb_devices+=("/dev/$device:$size:$model")
                fi
            fi
        fi
    done < <(lsblk -d -n -o NAME,SIZE,TYPE,RM,MODEL 2>/dev/null | grep -E "disk.*1" || true)
    
    printf '%s\n' "${usb_devices[@]}"
}

# Interactive USB device selection
select_usb_device() {
    local usb_devices=($(detect_usb_devices))
    
    if [ ${#usb_devices[@]} -eq 0 ]; then
        log ERROR "No USB storage devices detected."
        echo ""
        echo "Please:"
        echo "1. Insert a USB drive (minimum 8GB)"
        echo "2. Wait a few seconds for detection"
        echo "3. Run the script again"
        echo ""
        echo "If you have a USB connected but it's not detected, you can specify it manually:"
        echo "Available storage devices:"
        lsblk -d -o NAME,SIZE,TYPE,MODEL 2>/dev/null | grep disk | while read line; do
            echo "  /dev/$(echo $line | awk '{print $1}') ($(echo $line | awk '{print $2}')) - $(echo $line | awk '{for(i=4;i<=NF;i++) printf "%s ", $i; print ""}' | sed 's/[[:space:]]*$//')"
        done
        exit 1
    fi
    
    if [ ${#usb_devices[@]} -eq 1 ]; then
        # Only one USB device found
        local device_info="${usb_devices[0]}"
        local device=$(echo "$device_info" | cut -d':' -f1)
        local size=$(echo "$device_info" | cut -d':' -f2)
        local model=$(echo "$device_info" | cut -d':' -f3)
        
        log INFO "Auto-detected USB device:"
        echo "  Device: $device"
        echo "  Size: $size"
        echo "  Model: $model"
        echo ""
        
        read -p "Use this device? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            echo "$device"
            return 0
        else
            log INFO "Operation cancelled by user."
            exit 0
        fi
    else
        # Multiple USB devices found
        log INFO "Multiple USB devices detected:"
        echo ""
        
        local i=1
        for device_info in "${usb_devices[@]}"; do
            local device=$(echo "$device_info" | cut -d':' -f1)
            local size=$(echo "$device_info" | cut -d':' -f2)
            local model=$(echo "$device_info" | cut -d':' -f3)
            echo "  $i) $device ($size) - $model"
            ((i++))
        done
        echo "  0) Cancel"
        echo ""
        
        while true; do
            read -p "Select USB device (0-$((${#usb_devices[@]}))): " choice
            
            if [ "$choice" = "0" ]; then
                log INFO "Operation cancelled by user."
                exit 0
            elif [ "$choice" -ge 1 ] && [ "$choice" -le ${#usb_devices[@]} ] 2>/dev/null; then
                local selected_info="${usb_devices[$((choice-1))]}"
                local selected_device=$(echo "$selected_info" | cut -d':' -f1)
                echo "$selected_device"
                return 0
            else
                echo "Invalid selection. Please choose 0-$((${#usb_devices[@]}))"
            fi
        done
    fi
}

# Enhanced geographic coordinate detection
detect_coordinates() {
    local manual_coords="$1"
    
    if [ -n "$manual_coords" ]; then
        echo "$manual_coords:manual:Manual:Entry"
        return 0
    fi
    
    log INFO "Detecting geographic coordinates using multiple methods..."
    
    # Method 1: ipapi.co (most accurate)
    local result=$(timeout 15 curl -s "http://ipapi.co/json" 2>/dev/null || true)
    if [ -n "$result" ]; then
        local coords_result=$(echo "$result" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    lat = data.get('latitude')
    lon = data.get('longitude')
    city = data.get('city', 'unknown')
    country = data.get('country_name', 'unknown')
    if lat and lon and str(lat) != 'None' and str(lon) != 'None':
        print(f'{lat},{lon}:ip_geolocation_ipapi:{city}:{country}')
except Exception:
    pass
" 2>/dev/null || true)
        if [ -n "$coords_result" ] && [[ "$coords_result" != *"None"* ]]; then
            echo "$coords_result"
            return 0
        fi
    fi
    
    # Method 2: ipinfo.io (backup)
    result=$(timeout 15 curl -s "http://ipinfo.io/json" 2>/dev/null || true)
    if [ -n "$result" ]; then
        local coords_result=$(echo "$result" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    loc = data.get('loc', '')
    city = data.get('city', 'unknown')
    country = data.get('country', 'unknown')
    if loc and ',' in loc:
        print(f'{loc}:ip_geolocation_ipinfo:{city}:{country}')
except Exception:
    pass
" 2>/dev/null || true)
        if [ -n "$coords_result" ] && [[ "$coords_result" == *","* ]]; then
            echo "$coords_result"
            return 0
        fi
    fi
    
    # Method 3: Enhanced timezone-based approximation
    local timezone=$(timedatectl show --property=Timezone --value 2>/dev/null || echo "UTC")
    case "$timezone" in
        # Brazil
        "America/Sao_Paulo"|"America/Bahia"|"America/Fortaleza"|"America/Recife") 
            echo "-23.5505,-46.6333:timezone_brazil:São Paulo:Brazil" ;;
        "America/Manaus") 
            echo "-3.1190,-60.0217:timezone_brazil:Manaus:Brazil" ;;
        
        # North America
        "America/New_York"|"America/Detroit"|"America/Toronto") 
            echo "40.7128,-74.0060:timezone_us_eastern:New York:United States" ;;
        "America/Chicago"|"America/Winnipeg") 
            echo "41.8781,-87.6298:timezone_us_central:Chicago:United States" ;;
        "America/Denver"|"America/Edmonton") 
            echo "39.7392,-104.9903:timezone_us_mountain:Denver:United States" ;;
        "America/Los_Angeles"|"America/Vancouver") 
            echo "34.0522,-118.2437:timezone_us_pacific:Los Angeles:United States" ;;
        "America/Mexico_City") 
            echo "19.4326,-99.1332:timezone_mexico:Mexico City:Mexico" ;;
        
        # Europe
        "Europe/London"|"Europe/Belfast") 
            echo "51.5074,-0.1278:timezone_uk:London:United Kingdom" ;;
        "Europe/Paris") 
            echo "48.8566,2.3522:timezone_france:Paris:France" ;;
        "Europe/Berlin") 
            echo "52.5200,13.4050:timezone_germany:Berlin:Germany" ;;
        "Europe/Madrid") 
            echo "40.4168,-3.7038:timezone_spain:Madrid:Spain" ;;
        "Europe/Rome") 
            echo "41.9028,12.4964:timezone_italy:Rome:Italy" ;;
        "Europe/Amsterdam") 
            echo "52.3676,4.9041:timezone_netherlands:Amsterdam:Netherlands" ;;
        
        # Asia
        "Asia/Tokyo") 
            echo "35.6762,139.6503:timezone_japan:Tokyo:Japan" ;;
        "Asia/Seoul") 
            echo "37.5665,126.9780:timezone_korea:Seoul:South Korea" ;;
        "Asia/Shanghai"|"Asia/Beijing") 
            echo "31.2304,121.4737:timezone_china:Shanghai:China" ;;
        "Asia/Hong_Kong") 
            echo "22.3193,114.1694:timezone_hongkong:Hong Kong:Hong Kong" ;;
        "Asia/Singapore") 
            echo "1.3521,103.8198:timezone_singapore:Singapore:Singapore" ;;
        "Asia/Mumbai"|"Asia/Kolkata") 
            echo "19.0760,72.8777:timezone_india:Mumbai:India" ;;
        "Asia/Dubai") 
            echo "25.2048,55.2708:timezone_uae:Dubai:United Arab Emirates" ;;
        
        # Australia & Oceania
        "Australia/Sydney") 
            echo "-33.8688,151.2093:timezone_australia:Sydney:Australia" ;;
        "Australia/Melbourne") 
            echo "-37.8136,144.9631:timezone_australia:Melbourne:Australia" ;;
        "Pacific/Auckland") 
            echo "-36.8485,174.7633:timezone_newzealand:Auckland:New Zealand" ;;
        
        # Africa
        "Africa/Johannesburg") 
            echo "-26.2041,28.0473:timezone_southafrica:Johannesburg:South Africa" ;;
        "Africa/Cairo") 
            echo "30.0444,31.2357:timezone_egypt:Cairo:Egypt" ;;
        "Africa/Lagos") 
            echo "6.5244,3.3792:timezone_nigeria:Lagos:Nigeria" ;;
        
        # Default fallback
        *) 
            echo "0.0000,0.0000:timezone_unknown:Unknown:Unknown" ;;
    esac
}

# Generate location-based node ID
generate_location_id() {
    local coords_with_method="$1"
    local coords=$(echo "$coords_with_method" | cut -d':' -f1)
    local method=$(echo "$coords_with_method" | cut -d':' -f2)
    local city=$(echo "$coords_with_method" | cut -d':' -f3)
    
    # Create readable location ID
    local lat=$(echo "$coords" | cut -d',' -f1 | tr -d '-.' | cut -c1-4)
    local lon=$(echo "$coords" | cut -d',' -f2 | tr -d '-.' | cut -c1-4)
    
    # Clean city name for ID - remove special characters and limit length
    local city_clean=$(echo "$city" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9]//g' | cut -c1-8)
    
    # Method prefix
    local method_prefix=""
    case "$method" in
        "ip_geolocation"*) method_prefix="geo" ;;
        "timezone"*) method_prefix="tz" ;;
        "manual") method_prefix="man" ;;
        *) method_prefix="unk" ;;
    esac
    
    # Generate short random suffix
    local random_suffix=$(openssl rand -hex 3)
    
    # Ensure valid city name, fallback to coordinates
    if [ -z "$city_clean" ] || [ "$city_clean" = "unknown" ]; then
        city_clean="loc${lat}${lon}"
    fi
    
    echo "${method_prefix}-${city_clean}-${random_suffix}"
}

# Validate USB device safety
validate_usb_safety() {
    local device="$1"
    
    log INFO "Performing safety validation for $device..."
    
    # Critical: Check if device contains root filesystem
    if mount | grep -q "^$device.* / "; then
        log ERROR "CRITICAL: $device contains the root filesystem!"
        echo "This would destroy your operating system. Aborting immediately."
        exit 1
    fi
    
    # Check if any partition is mounted as system critical
    for critical_mount in / /boot /usr /var; do
        if mount | grep -q "^$device.*$critical_mount "; then
            log ERROR "CRITICAL: $device contains system partition: $critical_mount"
            echo "This would damage your system. Aborting immediately."
            exit 1
        fi
    done
    
    # Check if device is marked as removable
    local device_name=$(basename "$device")
    local removable_flag="/sys/block/$device_name/removable"
    if [ -f "$removable_flag" ]; then
        local is_removable=$(cat "$removable_flag" 2>/dev/null || echo "0")
        if [ "$is_removable" != "1" ]; then
            log WARN "$device is not marked as removable."
            echo "This might be a system disk. Device information:"
            lsblk -o NAME,SIZE,TYPE,MOUNTPOINT,MODEL "$device" 2>/dev/null || true
            echo ""
            read -p "Are you absolutely sure you want to continue? Type 'yes' to confirm: " confirm
            if [ "$confirm" != "yes" ]; then
                log INFO "Operation cancelled for safety."
                exit 0
            fi
        fi
    fi
    
    # Check device size (unusually large devices are suspicious)
    local device_size_gb=$(lsblk -b -d -n -o SIZE "$device" 2>/dev/null | awk '{printf "%.0f", $1/(1024*1024*1024)}' || echo "0")
    if [ "$device_size_gb" -gt 512 ]; then
        log WARN "$device is very large (${device_size_gb}GB)."
        echo "This is unusually large for a USB drive and might be a system disk."
        echo ""
        read -p "Are you absolutely sure this is a USB drive? Type 'yes' to confirm: " confirm
        if [ "$confirm" != "yes" ]; then
            log INFO "Operation cancelled for safety."
            exit 0
        fi
    fi
    
    # Show device information
    echo ""
    log INFO "Device to be erased:"
    lsblk -o NAME,SIZE,TYPE,MOUNTPOINT,MODEL "$device" 2>/dev/null || true
    
    # Check for mounted partitions and show them
    local mounted_parts=$(mount | grep "^$device" | wc -l)
    if [ "$mounted_parts" -gt 0 ]; then
        echo ""
        log WARN "Device has mounted partitions that will be unmounted:"
        mount | grep "^$device" || true
    fi
}

# Download and cache Ubuntu ISO
download_ubuntu_iso() {
    local iso_path="$ISO_CACHE_DIR/$ISO_FILE"
    
    log INFO "Setting up Ubuntu Server ISO..."
    mkdir -p "$ISO_CACHE_DIR"
    
    # Check if ISO exists in cache and verify
    if [ -f "$iso_path" ]; then
        log INFO "Found cached Ubuntu ISO, verifying integrity..."
        local cached_sha256=$(sha256sum "$iso_path" | cut -d' ' -f1)
        if [ "$cached_sha256" = "$ISO_SHA256" ]; then
            log SUCCESS "Cached ISO verified successfully"
            cp "$iso_path" "$WORK_DIR/$ISO_FILE"
            return 0
        else
            log WARN "Cached ISO checksum mismatch, removing..."
            rm -f "$iso_path"
        fi
    fi
    
    # Download ISO
    log INFO "Downloading Ubuntu 22.04.4 Server (~1.5GB)..."
    echo "This may take several minutes depending on your connection..."
    
    if wget --progress=bar:force:noscroll -c -O "$iso_path" "$ISO_URL"; then
        log INFO "Download completed, verifying integrity..."
        local downloaded_sha256=$(sha256sum "$iso_path" | cut -d' ' -f1)
        if [ "$downloaded_sha256" != "$ISO_SHA256" ]; then
            log ERROR "ISO checksum verification failed!"
            echo "Expected: $ISO_SHA256"
            echo "Got:      $downloaded_sha256"
            rm -f "$iso_path"
            return 1
        fi
        log SUCCESS "ISO downloaded and verified successfully"
        cp "$iso_path" "$WORK_DIR/$ISO_FILE"
        return 0
    else
        log ERROR "Failed to download Ubuntu ISO"
        echo "Please check your internet connection and try again"
        return 1
    fi
}

# Prepare USB device
prepare_usb_device() {
    local device="$1"
    
    log INFO "Preparing USB device: $device"
    
    # Get device information for logging
    local device_info=$(lsblk -d -n -o SIZE,MODEL "$device" 2>/dev/null || echo "Unknown Unknown")
    local device_size=$(echo "$device_info" | awk '{print $1}')
    local device_model=$(echo "$device_info" | awk '{for(i=2;i<=NF;i++) printf "%s ", $i; print ""}' | sed 's/[[:space:]]*$//')
    
    log INFO "Device information:"
    echo "  Device: $device"
    echo "  Size: $device_size"
    echo "  Model: $device_model"
    
    # Check minimum size (8GB)
    local device_size_bytes=$(lsblk -b -d -n -o SIZE "$device" 2>/dev/null || echo "0")
    local min_size_bytes=$((8 * 1024 * 1024 * 1024))
    
    if [ "$device_size_bytes" -lt "$min_size_bytes" ]; then
        log ERROR "USB device is too small (minimum 8GB required)"
        echo "Device size: $device_size"
        return 1
    fi
    
    # Unmount any existing partitions
    log INFO "Unmounting existing partitions..."
    local mounted_partitions=$(mount | grep "^$device" | awk '{print $1}' || true)
    if [ -n "$mounted_partitions" ]; then
        for partition in $mounted_partitions; do
            log INFO "Unmounting $partition"
            sudo umount "$partition" 2>/dev/null || true
        done
    fi
    
    # Additional cleanup
    sudo umount ${device}* 2>/dev/null || true
    sleep 2
    
    # Completely wipe device
    log INFO "Wiping device and creating new partition table..."
    sudo wipefs -a "$device" >/dev/null 2>&1 || true
    sudo sgdisk --zap-all "$device" >/dev/null 2>&1 || true
    
    # Create new partition table
    log INFO "Creating new partition structure..."
    if ! sudo parted -s "$device" mklabel msdos; then
        log ERROR "Failed to create partition table"
        return 1
    fi
    
    if ! sudo parted -s "$device" mkpart primary fat32 1MiB 100%; then
        log ERROR "Failed to create partition"
        return 1
    fi
    
    if ! sudo parted -s "$device" set 1 boot on; then
        log ERROR "Failed to set boot flag"
        return 1
    fi
    
    # Wait for partition to be recognized
    sleep 3
    
    # Determine partition name
    local partition="${device}1"
    if [ ! -b "$partition" ]; then
        partition="${device}p1"
        if [ ! -b "$partition" ]; then
            log ERROR "Cannot find USB partition after creation"
            return 1
        fi
    fi
    
    # Format partition
    log INFO "Formatting USB partition..."
    if ! sudo mkfs.fat -F32 -n "SYNTROPY" "$partition" >/dev/null 2>&1; then
        log ERROR "Failed to format USB partition"
        return 1
    fi
    
    echo "$partition"
    return 0
}

# Create cloud-init configuration
create_cloud_init_config() {
    local usb_mount="$1"
    local node_name="$2"
    local location_node_id="$3"
    local coordinates="$4"
    local detection_method="$5"
    local detected_city="$6"
    local detected_country="$7"
    local owner_key_file="$8"
    local community_key_file="$9"
    local node_description="${10}"
    
    log INFO "Creating cloud-init configuration..."
    
    # Read key contents
    local owner_key_content=$(cat "$owner_key_file")
    local owner_pub_content=$(cat "${owner_key_file}.pub")
    local community_key_content=$(cat "$community_key_file")
    local community_pub_content=$(cat "${community_key_file}.pub")
    
    # Generate owner and community fingerprints
    local owner_fingerprint=$(ssh-keygen -lf "${owner_key_file}.pub" | awk '{print $2}')
    local community_fingerprint=$(ssh-keygen -lf "${community_key_file}.pub" | awk '{print $2}')
    
    # Create user-data file
    sudo tee "$usb_mount/user-data" > /dev/null << USER_DATA_EOF
#cloud-config
# Syntropy Cooperative Grid - Enhanced Auto-Installation
# Node: $node_name | Location: $coordinates | Created: $(date)

autoinstall:
  version: 1
  interactive-sections: []
  locale: en_US.UTF-8
  keyboard:
    layout: us
    
  network:
    network:
      version: 2
      ethernets:
        "en*":
          dhcp4: true
          dhcp6: false
          dhcp4-overrides:
            hostname: $node_name
        "eth*":
          dhcp4: true
          dhcp6: false
          dhcp4-overrides:
            hostname: $node_name
        "enp*":
          dhcp4: true
          dhcp6: false
          dhcp4-overrides:
            hostname: $node_name

  identity:
    hostname: $node_name
    username: admin
    password: "\$6\$rounds=4096\$syntropy\$N8mVzFK0Y1OelT1SKEjg0jIXzKMzL3ZcOGcE5xR8nS6E8qSO5qFV6eJs1g7T6E0cC7w.kfNO3FqC3YhE9Gz19."

  ssh:
    install-server: true
    allow-pw: false
    
  storage:
    layout:
      name: lvm
      sizing-policy: all
    swap:
      size: 0

  packages:
    - curl
    - wget
    - git
    - vim
    - htop
    - jq
    - python3
    - python3-pip
    - docker.io
    - fail2ban
    - ufw
    - prometheus-node-exporter
    - openssh-server
    - nmap
    - ncdu
    - tree
    - tmux
    - net-tools

  late-commands:
    # Create comprehensive directory structure
    - curtin in-target -- mkdir -p /opt/syntropy/{identity/{owner,community},platform/{templates,scripts,data},metadata,logs,backups}
    
    # Install security keys and metadata
    - |
      curtin in-target -- bash -c '
      # Install owner key (SSH access and management)
      cat > /opt/syntropy/identity/owner/private.key << "OWNER_KEY_EOF"
$owner_key_content
OWNER_KEY_EOF
      
      cat > /opt/syntropy/identity/owner/public.key << "OWNER_PUB_EOF"
$owner_pub_content
OWNER_PUB_EOF
      
      # Install community key (inter-node communication)
      cat > /opt/syntropy/identity/community/private.key << "COMMUNITY_KEY_EOF"
$community_key_content
COMMUNITY_KEY_EOF
      
      cat > /opt/syntropy/identity/community/public.key << "COMMUNITY_PUB_EOF"
$community_pub_content
COMMUNITY_PUB_EOF
      
      # Create key metadata
      cat > /opt/syntropy/identity/key_info.json << "KEY_INFO_EOF"
{
  "owner_key": {
    "fingerprint": "$owner_fingerprint",
    "algorithm": "ed25519",
    "purpose": "ssh_access_and_management",
    "created": "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
  },
  "community_key": {
    "fingerprint": "$community_fingerprint", 
    "algorithm": "ed25519",
    "purpose": "inter_node_communication",
    "created": "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
  }
}
KEY_INFO_EOF
      
      # Set proper permissions
      chmod 600 /opt/syntropy/identity/owner/private.key
      chmod 600 /opt/syntropy/identity/community/private.key
      chmod 644 /opt/syntropy/identity/owner/public.key
      chmod 644 /opt/syntropy/identity/community/public.key
      chmod 644 /opt/syntropy/identity/key_info.json
      chown -R admin:admin /opt/syntropy/
      
      # Configure SSH access
      mkdir -p /home/admin/.ssh
      cp /opt/syntropy/identity/owner/public.key /home/admin/.ssh/authorized_keys
      chmod 600 /home/admin/.ssh/authorized_keys
      chown admin:admin /home/admin/.ssh/authorized_keys
      '
    
    # Create comprehensive node metadata with hardware detection
    - |
      curtin in-target -- bash -c '
      # Detect hardware specifications
      CPU_CORES=\$(nproc)
      RAM_GB=\$(free -g | awk "/^Mem:/{\print \\\$2}")
      STORAGE_GB=\$(df / --output=avail -BG 2>/dev/null | tail -1 | sed "s/G//" | xargs)
      ARCHITECTURE=\$(uname -m)
      
      # Enhanced hardware classification
      if [ \$RAM_GB -le 2 ]; then
        HW_CLASS="edge"
        K8S_ROLE="worker-light"
        CAPABILITIES="[\"edge_computing\", \"sensor_data\", \"lightweight_services\"]"
      elif [ \$RAM_GB -le 8 ]; then
        HW_CLASS="home-server"
        K8S_ROLE="worker"
        CAPABILITIES="[\"container_hosting\", \"development\", \"personal_services\"]"
      elif [ \$RAM_GB -le 32 ]; then
        HW_CLASS="workstation"
        K8S_ROLE="worker-heavy"
        CAPABILITIES="[\"compute_intensive\", \"ai_inference\", \"database_hosting\", \"media_processing\"]"
      else
        HW_CLASS="server"
        K8S_ROLE="control-plane-capable"
        CAPABILITIES="[\"distributed_computing\", \"ai_training\", \"high_performance\", \"cluster_management\"]"
      fi
      
      # Create comprehensive node metadata
      cat > /opt/syntropy/metadata/node.json << "NODE_METADATA_EOF"
{
  "node_identity": {
    "name": "$node_name",
    "location_id": "$location_node_id",
    "description": "$node_description",
    "created_at": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
    "version": "2.0.0"
  },
  "geographic_info": {
    "coordinates": "$coordinates",
    "detection_method": "$detection_method",
    "city": "$detected_city",
    "country": "$detected_country",
    "timezone": "\$(timedatectl show --property=Timezone --value 2>/dev/null || echo UTC)"
  },
  "hardware_profile": {
    "class": "\$HW_CLASS",
    "cpu_cores": \$CPU_CORES,
    "ram_gb": \$RAM_GB,
    "storage_gb": \$STORAGE_GB,
    "architecture": "\$ARCHITECTURE",
    "kubernetes_role": "\$K8S_ROLE",
    "capabilities": \$CAPABILITIES
  },
  "network_config": {
    "hostname": "$node_name",
    "ports": {
      "ssh": 22,
      "kubernetes_api": 6443,
      "node_exporter": 9100,
      "syntropy_api": 8080
    }
  },
  "security": {
    "owner_key_fingerprint": "$owner_fingerprint",
    "community_key_fingerprint": "$community_fingerprint",
    "ssh_password_auth": false,
    "firewall_enabled": true
  }
}
NODE_METADATA_EOF
      '
    
    # Configure security and firewall
    - |
      curtin in-target -- bash -c '
      # Configure UFW firewall
      ufw --force enable
      ufw default deny incoming
      ufw default allow outgoing
      
      # Allow SSH
      ufw allow 22/tcp
      
      # Allow Kubernetes ports (conditional based on role)
      if [ "\$(cat /opt/syntropy/metadata/node.json | jq -r .hardware_profile.kubernetes_role)" != "worker-light" ]; then
        ufw allow 6443/tcp  # Kubernetes API
        ufw allow 2379:2380/tcp  # etcd
        ufw allow 10250/tcp  # Kubelet API
        ufw allow 10251/tcp  # kube-scheduler
        ufw allow 10252/tcp  # kube-controller-manager
      fi
      
      # Allow worker node ports
      ufw allow 10250/tcp  # Kubelet API
      ufw allow 30000:32767/tcp  # NodePort services
      
      # Allow monitoring
      ufw allow 9100/tcp  # Node Exporter
      
      # Allow Syntropy communication
      ufw allow 8080/tcp  # Syntropy API
      ufw allow 51820/udp  # WireGuard VPN
      
      # Configure fail2ban
      systemctl enable fail2ban
      
      # Create fail2ban SSH jail
      cat > /etc/fail2ban/jail.local << "FAIL2BAN_EOF"
[sshd]
enabled = true
port = ssh
filter = sshd
logpath = /var/log/auth.log
maxretry = 3
bantime = 3600
findtime = 600
FAIL2BAN_EOF
      '
    
    # Install Docker and Kubernetes prerequisites
    - |
      curtin in-target -- bash -c '
      # Add Docker to admin user group
      usermod -aG docker admin
      
      # Configure Docker daemon
      mkdir -p /etc/docker
      cat > /etc/docker/daemon.json << "DOCKER_CONFIG_EOF"
{
  "exec-opts": ["native.cgroupdriver=systemd"],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m",
    "max-file": "3"
  },
  "storage-driver": "overlay2",
  "insecure-registries": ["registry.syntropy.local:5000"]
}
DOCKER_CONFIG_EOF
      
      # Enable and start Docker
      systemctl enable docker
      systemctl start docker
      
      # Install Kubernetes packages
      curl -fsSL https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
      echo "deb https://apt.kubernetes.io/ kubernetes-xenial main" > /etc/apt/sources.list.d/kubernetes.list
      apt-get update
      apt-get install -y kubelet kubeadm kubectl
      apt-mark hold kubelet kubeadm kubectl
      
      # Configure kubelet
      echo "KUBELET_EXTRA_ARGS=--node-ip=\$(hostname -I | awk \"{print \\\$1}\")" > /etc/default/kubelet
      systemctl enable kubelet
      '
    
    # Install Syntropy platform components
    - |
      curtin in-target -- bash -c '
      # Create Syntropy management scripts
      cat > /opt/syntropy/platform/scripts/bootstrap.sh << "BOOTSTRAP_EOF"
#!/bin/bash
# Syntropy Node Bootstrap Script
set -e

echo "Starting Syntropy Node Bootstrap..."
echo "Node: $node_name"
echo "Location: $coordinates"
echo "Hardware Class: \$(cat /opt/syntropy/metadata/node.json | jq -r .hardware_profile.class)"

# Wait for network connectivity
echo "Waiting for network connectivity..."
while ! ping -c 1 8.8.8.8 >/dev/null 2>&1; do
  sleep 5
done

# Register with Syntropy network
echo "Registering with Syntropy Cooperative Grid..."
curl -X POST https://api.syntropy.coop/v1/nodes/register \\
  -H "Content-Type: application/json" \\
  -d @/opt/syntropy/metadata/node.json \\
  > /opt/syntropy/logs/registration.log 2>&1 || true

# Start monitoring
systemctl enable prometheus-node-exporter
systemctl start prometheus-node-exporter

# Initialize based on hardware class
HW_CLASS=\$(cat /opt/syntropy/metadata/node.json | jq -r .hardware_profile.class)
case "\$HW_CLASS" in
  "edge")
    echo "Configuring edge node..."
    /opt/syntropy/platform/scripts/setup-edge.sh
    ;;
  "home-server"|"workstation")
    echo "Configuring worker node..."
    /opt/syntropy/platform/scripts/setup-worker.sh
    ;;
  "server")
    echo "Configuring server node..."
    /opt/syntropy/platform/scripts/setup-server.sh
    ;;
esac

echo "Bootstrap completed successfully!"
BOOTSTRAP_EOF

      chmod +x /opt/syntropy/platform/scripts/bootstrap.sh
      
      # Create edge node setup script
      cat > /opt/syntropy/platform/scripts/setup-edge.sh << "SETUP_EDGE_EOF"
#!/bin/bash
# Edge Node Setup
echo "Setting up edge computing capabilities..."

# Install lightweight container runtime
docker pull alpine:latest
docker pull nginx:alpine
docker pull node:alpine

# Create edge services directory
mkdir -p /opt/syntropy/platform/data/edge-services

echo "Edge node setup completed"
SETUP_EDGE_EOF

      chmod +x /opt/syntropy/platform/scripts/setup-edge.sh
      
      # Create worker node setup script
      cat > /opt/syntropy/platform/scripts/setup-worker.sh << "SETUP_WORKER_EOF"
#!/bin/bash
# Worker Node Setup
echo "Setting up Kubernetes worker node..."

# Pull common container images
docker pull k8s.gcr.io/pause:3.7
docker pull k8s.gcr.io/coredns/coredns:v1.8.6

# Wait for cluster join instructions
echo "Worker node ready for cluster joining"
echo "Use: kubeadm join <master-ip>:6443 --token <token> --discovery-token-ca-cert-hash <hash>"
SETUP_WORKER_EOF

      chmod +x /opt/syntropy/platform/scripts/setup-worker.sh
      
      # Create server node setup script
      cat > /opt/syntropy/platform/scripts/setup-server.sh << "SETUP_SERVER_EOF"
#!/bin/bash
# Server Node Setup
echo "Setting up Kubernetes control plane..."

# Initialize Kubernetes cluster
kubeadm init --pod-network-cidr=10.244.0.0/16 \\
  --service-cidr=10.96.0.0/12 \\
  --node-name=$node_name

# Configure kubectl for admin user
mkdir -p /home/admin/.kube
cp /etc/kubernetes/admin.conf /home/admin/.kube/config
chown admin:admin /home/admin/.kube/config

# Install Calico CNI
kubectl --kubeconfig=/home/admin/.kube/config apply -f https://docs.projectcalico.org/manifests/calico.yaml

echo "Control plane setup completed"
echo "Join tokens available in /var/log/kubeadm.log"
SETUP_SERVER_EOF

      chmod +x /opt/syntropy/platform/scripts/setup-server.sh
      '
    
    # Create systemd service for Syntropy bootstrap
    - |
      curtin in-target -- bash -c '
      cat > /etc/systemd/system/syntropy-bootstrap.service << "SERVICE_EOF"
[Unit]
Description=Syntropy Cooperative Grid Bootstrap
After=network-online.target docker.service
Wants=network-online.target
Requires=docker.service

[Service]
Type=oneshot
ExecStart=/opt/syntropy/platform/scripts/bootstrap.sh
RemainAfterExit=yes
StandardOutput=journal
StandardError=journal
User=root

[Install]
WantedBy=multi-user.target
SERVICE_EOF

      systemctl enable syntropy-bootstrap.service
      '
    
    # Final system configuration
    - |
      curtin in-target -- bash -c '
      # Create status script
      cat > /opt/syntropy/platform/scripts/status.sh << "STATUS_EOF"
#!/bin/bash
# Syntropy Node Status
echo "=== Syntropy Cooperative Grid Node Status ==="
echo "Node Name: $node_name"
echo "Location: $coordinates"
echo "Hardware Class: \$(cat /opt/syntropy/metadata/node.json | jq -r .hardware_profile.class)"
echo "Created: \$(cat /opt/syntropy/metadata/node.json | jq -r .node_identity.created_at)"
echo ""
echo "=== System Information ==="
echo "Hostname: \$(hostname)"
echo "IP Address: \$(hostname -I | awk \"{print \\\$1}\")"
echo "Uptime: \$(uptime -p)"
echo "Load: \$(cat /proc/loadavg | awk \"{print \\\$1, \\\$2, \\\$3}\")"
echo ""
echo "=== Docker Status ==="
systemctl is-active docker
docker version --format "{{.Server.Version}}" 2>/dev/null || echo "Not available"
echo ""
echo "=== Kubernetes Status ==="
systemctl is-active kubelet
kubectl version --client --short 2>/dev/null || echo "Not configured"
echo ""
echo "=== Security Status ==="
echo "UFW: \$(ufw status | head -1)"
echo "Fail2ban: \$(systemctl is-active fail2ban)"
echo "SSH Keys: \$(ls -la /opt/syntropy/identity/*/public.key | wc -l) configured"
STATUS_EOF

      chmod +x /opt/syntropy/platform/scripts/status.sh
      
      # Set proper ownership
      chown -R admin:admin /opt/syntropy/
      
      # Create welcome message
      cat > /etc/motd << "MOTD_EOF"

╔════════════════════════════════════════════════════════════════════════════╗
║                    SYNTROPY COOPERATIVE GRID NODE                         ║
║                                                                            ║
║  Node: $node_name                                                         ║
║  Location: $coordinates                                                   ║
║  Created: $(date)                                            ║
║                                                                            ║
║  Quick Commands:                                                           ║
║    sudo /opt/syntropy/platform/scripts/status.sh  - Node status           ║
║    docker ps                                       - Running containers   ║
║    kubectl get nodes                              - Kubernetes cluster    ║
║                                                                            ║
║  Configuration: /opt/syntropy/                                             ║
║  Logs: journalctl -u syntropy-bootstrap                                    ║
╚════════════════════════════════════════════════════════════════════════════╝

MOTD_EOF
      '

USER_DATA_EOF

    # Create meta-data file (required for cloud-init)
    sudo tee "$usb_mount/meta-data" > /dev/null << META_DATA_EOF
instance-id: syntropy-$location_node_id
local-hostname: $node_name
META_DATA_EOF

    log SUCCESS "Cloud-init configuration created successfully"
}

# Setup SSH keys
setup_ssh_keys() {
    local node_name="$1"
    
    log INFO "Setting up SSH keys for node management..."
    
    mkdir -p "$KEYS_DIR"
    
    # Generate owner SSH key (for human access)
    local owner_key_file="$KEYS_DIR/owner_${node_name}_$(date +%Y%m%d)"
    if [ ! -f "$owner_key_file" ]; then
        log INFO "Generating owner SSH key..."
        ssh-keygen -t ed25519 -f "$owner_key_file" -N "" -C "syntropy-owner-$node_name-$(date +%Y%m%d)" >/dev/null 2>&1
        log SUCCESS "Owner key generated: $owner_key_file"
    else
        log INFO "Using existing owner key: $owner_key_file"
    fi
    
    # Generate community SSH key (for inter-node communication)
    local community_key_file="$KEYS_DIR/community_${node_name}_$(date +%Y%m%d)"
    if [ ! -f "$community_key_file" ]; then
        log INFO "Generating community SSH key..."
        ssh-keygen -t ed25519 -f "$community_key_file" -N "" -C "syntropy-community-$node_name-$(date +%Y%m%d)" >/dev/null 2>&1
        log SUCCESS "Community key generated: $community_key_file"
    else
        log INFO "Using existing community key: $community_key_file"
    fi
    
    echo "$owner_key_file:$community_key_file"
}

# Create USB bootloader
create_bootloader() {
    local usb_mount="$1"
    local iso_file="$2"
    
    log INFO "Creating GRUB bootloader configuration..."
    
    # Create GRUB directory structure
    sudo mkdir -p "$usb_mount/boot/grub"
    
    # Create GRUB configuration
    sudo tee "$usb_mount/boot/grub/grub.cfg" > /dev/null << 'GRUB_EOF'
set timeout=10
set default=0

menuentry "Syntropy Cooperative Grid - Auto Install" {
    set isofile="/ubuntu-22.04.4-live-server-amd64.iso"
    loopback loop (hd0,msdos1)$isofile
    linux (loop)/casper/vmlinuz boot=casper iso-scan/filename=$isofile autoinstall quiet splash ---
    initrd (loop)/casper/initrd
}

menuentry "Syntropy Cooperative Grid - Manual Install" {
    set isofile="/ubuntu-22.04.4-live-server-amd64.iso"
    loopback loop (hd0,msdos1)$isofile
    linux (loop)/casper/vmlinuz boot=casper iso-scan/filename=$isofile quiet splash ---
    initrd (loop)/casper/initrd
}

menuentry "Boot from Hard Drive" {
    set root=(hd1)
    chainloader +1
}
GRUB_EOF

    # Install GRUB to USB
    log INFO "Installing GRUB bootloader to USB..."
    local usb_device=$(echo "$usb_mount" | sed 's/[0-9]*$//')
    sudo grub-install --target=i386-pc --boot-directory="$usb_mount/boot" "$usb_device" >/dev/null 2>&1
    
    log SUCCESS "Bootloader created successfully"
}

# Copy ISO to USB
copy_iso_to_usb() {
    local usb_mount="$1"
    local iso_file="$2"
    
    log INFO "Copying Ubuntu ISO to USB (this may take several minutes)..."
    if sudo cp "$iso_file" "$usb_mount/"; then
        log SUCCESS "ISO copied successfully"
        return 0
    else
        log ERROR "Failed to copy ISO to USB"
        return 1
    fi
}

# Generate node documentation
generate_documentation() {
    local usb_mount="$1"
    local node_name="$2"
    local coordinates="$3"
    local owner_key_file="$4"
    local community_key_file="$5"
    
    log INFO "Generating node documentation..."
    
    sudo tee "$usb_mount/README.txt" > /dev/null << DOC_EOF
SYNTROPY COOPERATIVE GRID - Node Installation USB
================================================

Node Information:
- Name: $node_name
- Coordinates: $coordinates
- Created: $(date)
- Version: 2.0.0

INSTALLATION INSTRUCTIONS:
=========================

1. BEFORE BOOTING:
   - Ensure target computer has minimum 8GB RAM and 64GB storage
   - Connect to internet via Ethernet (recommended) or WiFi
   - Backup any important data (installation will erase the disk)

2. BOOT FROM USB:
   - Insert this USB into target computer
   - Boot from USB (usually F12, F2, or DEL during startup)
   - Select "Syntropy Cooperative Grid - Auto Install"

3. INSTALLATION PROCESS:
   - Installation is fully automated (15-30 minutes)
   - System will reboot automatically when complete
   - Remove USB after reboot

4. FIRST LOGIN:
   - Username: admin
   - Authentication: SSH key only (no password)
   - Use the private key: $(basename "$owner_key_file")

5. POST-INSTALLATION:
   - Check status: sudo /opt/syntropy/platform/scripts/status.sh
   - View logs: journalctl -u syntropy-bootstrap
   - Node configuration: /opt/syntropy/metadata/node.json

SSH ACCESS:
===========
Private key file: $(basename "$owner_key_file")
Public key fingerprint: $(ssh-keygen -lf "${owner_key_file}.pub" 2>/dev/null | awk '{print $2}' || echo "N/A")

To connect:
ssh -i $(basename "$owner_key_file") admin@<node-ip-address>

SECURITY NOTES:
===============
- SSH password authentication is DISABLED
- Firewall (UFW) is enabled with restrictive rules
- Fail2ban protects against brute force attacks
- All keys are ed25519 for enhanced security

TROUBLESHOOTING:
================
- If boot fails: Try "Manual Install" option in GRUB menu
- If network issues: Check Ethernet connection
- If SSH fails: Verify private key permissions (chmod 600)
- Support: https://docs.syntropy.coop/troubleshooting

For more information: https://syntropy.coop/cooperative-grid
DOC_EOF

    log SUCCESS "Documentation generated successfully"
}

# Main execution function
main() {
    log INFO "Starting Syntropy USB Creator..."
    
    # Check dependencies
    for cmd in wget curl ssh-keygen grub-install parted mkfs.fat lsblk python3; do
        if ! command -v "$cmd" >/dev/null 2>&1; then
            log ERROR "Required command not found: $cmd"
            echo "Please install missing dependencies and try again."
            exit 1
        fi
    done
    
    # Collect user input
    echo ""
    echo "=== Node Configuration ==="
    
    # Get node name
    while true; do
        read -p "Enter node name (e.g., syntropy-home-01): " node_name
        if [[ "$node_name" =~ ^[a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]$ ]] && [ ${#node_name} -ge 3 ] && [ ${#node_name} -le 63 ]; then
            break
        else
            echo "Invalid name. Use 3-63 characters, alphanumeric and hyphens only."
        fi
    done
    
    # Get coordinates
    echo ""
    echo "Geographic location detection:"
    echo "1) Auto-detect from IP geolocation"
    echo "2) Enter coordinates manually"
    read -p "Choose option (1-2): " coord_option
    
    local manual_coords=""
    if [ "$coord_option" = "2" ]; then
        read -p "Enter coordinates (lat,lon e.g. -23.5505,-46.6333): " manual_coords
        if ! [[ "$manual_coords" =~ ^-?[0-9]+\.?[0-9]*,-?[0-9]+\.?[0-9]*$ ]]; then
            log ERROR "Invalid coordinates format"
            exit 1
        fi
    fi
    
    # Detect coordinates
    local coords_info=$(detect_coordinates "$manual_coords")
    local coordinates=$(echo "$coords_info" | cut -d':' -f1)
    local detection_method=$(echo "$coords_info" | cut -d':' -f2)
    local detected_city=$(echo "$coords_info" | cut -d':' -f3)
    local detected_country=$(echo "$coords_info" | cut -d':' -f4)
    
    log INFO "Location detected: $detected_city, $detected_country ($coordinates)"
    
    # Generate location-based ID
    local location_node_id=$(generate_location_id "$coords_info")
    log INFO "Generated location ID: $location_node_id"
    
    # Get node description
    read -p "Enter node description (optional): " node_description
    [ -z "$node_description" ] && node_description="Syntropy Cooperative Grid Node"
    
    # Create working directory
    mkdir -p "$WORK_DIR"
    cd "$WORK_DIR"
    
    # Setup SSH keys
    local key_files=$(setup_ssh_keys "$node_name")
    local owner_key_file=$(echo "$key_files" | cut -d':' -f1)
    local community_key_file=$(echo "$key_files" | cut -d':' -f2)
    
    # Select USB device
    echo ""
    echo "=== USB Device Selection ==="
    local usb_device=$(select_usb_device)
    
    # Validate USB safety
    validate_usb_safety "$usb_device"
    
    # Final confirmation
    echo ""
    echo "=== FINAL CONFIRMATION ==="
    echo "Node Name: $node_name"
    echo "Location: $detected_city, $detected_country ($coordinates)"
    echo "USB Device: $usb_device"
    echo ""
    echo "WARNING: This will completely erase the USB device!"
    echo ""
    read -p "Continue with USB creation? Type 'yes' to confirm: " final_confirm
    
    if [ "$final_confirm" != "yes" ]; then
        log INFO "Operation cancelled by user."
        exit 0
    fi
    
    # Download Ubuntu ISO
    download_ubuntu_iso
    
    # Prepare USB device
    local usb_partition=$(prepare_usb_device "$usb_device")
    
    # Mount USB
    USB_MOUNT="/tmp/syntropy-usb-mount-$"
    mkdir -p "$USB_MOUNT"
    sudo mount "$usb_partition" "$USB_MOUNT"
    
    # Copy ISO to USB
    copy_iso_to_usb "$USB_MOUNT" "$WORK_DIR/$ISO_FILE"
    
    # Create cloud-init configuration
    create_cloud_init_config "$USB_MOUNT" "$node_name" "$location_node_id" "$coordinates" "$detection_method" "$detected_city" "$detected_country" "$owner_key_file" "$community_key_file" "$node_description"
    
    # Create bootloader
    create_bootloader "$USB_MOUNT" "$ISO_FILE"
    
    # Generate documentation
    generate_documentation "$USB_MOUNT" "$node_name" "$coordinates" "$owner_key_file" "$community_key_file"
    
    # Copy SSH keys to USB for user
    sudo cp "$owner_key_file" "$USB_MOUNT/"
    sudo cp "${owner_key_file}.pub" "$USB_MOUNT/"
    
    # Sync and unmount
    log INFO "Finalizing USB creation..."
    sync
    sudo umount "$USB_MOUNT"
    rmdir "$USB_MOUNT"
    USB_MOUNT=""
    
    # Success summary
    echo ""
    log SUCCESS "Syntropy USB created successfully!"
    echo ""
    echo "=== SUMMARY ==="
    echo "USB Device: $usb_device"
    echo "Node Name: $node_name"
    echo "Location: $detected_city, $detected_country"
    echo "Coordinates: $coordinates"
    echo "SSH Key: $owner_key_file"
    echo ""
    echo "=== NEXT STEPS ==="
    echo "1. Safely remove the USB device"
    echo "2. Boot target computer from USB"
    echo "3. Select 'Auto Install' in GRUB menu"
    echo "4. Wait for installation to complete (~30 minutes)"
    echo "5. Connect via SSH using the generated key"
    echo ""
    echo "SSH Connection:"
    echo "  ssh -i $owner_key_file admin@<node-ip>"
    echo ""
    echo "Documentation and keys are saved on the USB drive."
    echo "For support: https://docs.syntropy.coop"
}

# Run main function
main "$@"