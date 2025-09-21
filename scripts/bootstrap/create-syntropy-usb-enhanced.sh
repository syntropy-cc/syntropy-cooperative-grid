#!/bin/bash

# Syntropy Cooperative Grid - Enhanced USB Creator with Advanced Node Management
# Creates production-ready USB with automated node management setup

set -e

USB_DEVICE=$1
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
WORK_DIR="/tmp/syntropy-usb-enhanced"
NODES_DIR="$HOME/.syntropy/nodes"
KEYS_DIR="$HOME/.syntropy/keys"
CONFIG_DIR="$HOME/.syntropy/config"
TIMESTAMP=$(date +%s)

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m'

echo -e "${PURPLE}"
cat << 'EOF'
╔════════════════════════════════════════════════════════════════════════════╗
║                    SYNTROPY COOPERATIVE GRID                               ║
║                  Enhanced USB Creator with Node Management                 ║
║                                                                            ║
║  Creates production-ready USB with automated node management setup         ║
╚════════════════════════════════════════════════════════════════════════════╝
EOF
echo -e "${NC}"

# Function to show help
show_help() {
    echo "Usage: $0 <usb_device> [options]"
    echo ""
    echo "Options:"
    echo "  --owner-key <file>        Use existing owner private key (enables multi-node management)"
    echo "  --node-name <name>        Custom node name (default: auto-generated from location)"
    echo "  --description <desc>      Node description"
    echo "  --coordinates <lat,lon>   Manual coordinates override (default: auto-detected)"
    echo "  --help                    Show this help"
    echo ""
    echo "Examples:"
    echo "  $0 /dev/sdb                                              # Auto-detected setup"
    echo "  $0 /dev/sdb --node-name home-server-01                  # Custom name"
    echo "  $0 /dev/sdb --owner-key ~/.syntropy/keys/main_owner.key # Use existing key (multi-node)"
    echo "  $0 /dev/sdb --coordinates \"-23.5505,-46.6333\"           # Manual coordinates"
    echo ""
    echo "Multi-node workflow:"
    echo "  1. First node:  $0 /dev/sdb --node-name main-server"
    echo "  2. Second node: $0 /dev/sdb --owner-key ~/.syntropy/keys/main-server_owner.key --node-name edge-01"
    echo "  3. Third node:  $0 /dev/sdb --owner-key ~/.syntropy/keys/main-server_owner.key --node-name edge-02"
    echo ""
    echo "Note: Using the same owner key allows unified management of multiple nodes"
}

# Enhanced geographic coordinate detection with multiple fallback methods
detect_coordinates() {
    local manual_coords="$1"
    
    if [ -n "$manual_coords" ]; then
        echo "$manual_coords:manual"
        return 0
    fi
    
    echo "Detecting geographic coordinates using multiple methods..."
    
    # Method 1: ipapi.co (most accurate, includes ISP info)
    local result=$(timeout 10 curl -s "http://ipapi.co/json" 2>/dev/null)
    if [ -n "$result" ]; then
        local coords=$(echo "$result" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    lat = data.get('latitude')
    lon = data.get('longitude')
    city = data.get('city', 'unknown')
    country = data.get('country_name', 'unknown')
    if lat and lon and str(lat) != 'None' and str(lon) != 'None':
        print(f'{lat},{lon}:ip_geolocation_ipapi:{city}:{country}')
except:
    pass
" 2>/dev/null)
        if [ -n "$coords" ] && [[ "$coords" != *"None"* ]]; then
            echo "$coords"
            return 0
        fi
    fi
    
    # Method 2: ipinfo.io (backup with different data source)
    result=$(timeout 10 curl -s "http://ipinfo.io/json" 2>/dev/null)
    if [ -n "$result" ]; then
        coords=$(echo "$result" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    loc = data.get('loc', '')
    city = data.get('city', 'unknown')
    country = data.get('country', 'unknown')
    if loc and ',' in loc:
        print(f'{loc}:ip_geolocation_ipinfo:{city}:{country}')
except:
    pass
" 2>/dev/null)
        if [ -n "$coords" ] && [[ "$coords" != *","* ]]; then
            echo "$coords"
            return 0
        fi
    fi
    
    # Method 3: Enhanced timezone-based approximation with major cities
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

# Generate location-based node ID with better encoding
generate_location_id() {
    local coords_with_method="$1"
    local coords=$(echo "$coords_with_method" | cut -d':' -f1)
    local method=$(echo "$coords_with_method" | cut -d':' -f2)
    local city=$(echo "$coords_with_method" | cut -d':' -f3)
    
    # Create more readable location ID
    local lat=$(echo "$coords" | cut -d',' -f1 | tr -d '-' | cut -c1-4)
    local lon=$(echo "$coords" | cut -d',' -f2 | tr -d '-' | cut -c1-4)
    
    # Clean city name for ID
    local city_clean=$(echo "$city" | tr '[:upper:]' '[:lower:]' | tr -cd '[:alnum:]' | cut -c1-8)
    
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
    
    echo "${method_prefix}-${city_clean}-${lat}${lon}-${random_suffix}"
}

# Parse command line arguments
OWNER_KEY_FILE=""
NODE_NAME=""
NODE_DESCRIPTION=""
MANUAL_COORDINATES=""

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
            echo "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Validate USB device
if [ -z "$USB_DEVICE" ]; then
    echo -e "${RED}Error: USB device not specified${NC}"
    show_help
    exit 1
fi

if [ ! -b "$USB_DEVICE" ]; then
    echo -e "${RED}Error: $USB_DEVICE is not a valid block device${NC}"
    echo "Available devices:"
    lsblk | grep disk
    exit 1
fi

# Create management directories
mkdir -p "$NODES_DIR" "$KEYS_DIR" "$CONFIG_DIR"
mkdir -p "$WORK_DIR"
cd "$WORK_DIR"

echo -e "${BLUE}[1/8] Detecting location and generating node identity...${NC}"

# Enhanced coordinate detection with location information
COORDS_WITH_INFO=$(detect_coordinates "$MANUAL_COORDINATES")
COORDINATES=$(echo "$COORDS_WITH_INFO" | cut -d':' -f1)
DETECTION_METHOD=$(echo "$COORDS_WITH_INFO" | cut -d':' -f2)
DETECTED_CITY=$(echo "$COORDS_WITH_INFO" | cut -d':' -f3)
DETECTED_COUNTRY=$(echo "$COORDS_WITH_INFO" | cut -d':' -f4)

echo "Location detected:"
echo "  Coordinates: $COORDINATES"
echo "  City: $DETECTED_CITY"
echo "  Country: $DETECTED_COUNTRY"
echo "  Method: $DETECTION_METHOD"

# Generate location-based node ID
LOCATION_NODE_ID=$(generate_location_id "$COORDS_WITH_INFO")

echo -e "${BLUE}[2/8] Setting up security keys...${NC}"

# Enhanced key management with validation
if [ -n "$OWNER_KEY_FILE" ]; then
    if [ ! -f "$OWNER_KEY_FILE" ]; then
        echo -e "${RED}Error: Owner key file not found: $OWNER_KEY_FILE${NC}"
        exit 1
    fi
    
    # Validate key format
    if ! ssh-keygen -l -f "$OWNER_KEY_FILE" >/dev/null 2>&1; then
        echo -e "${RED}Error: Invalid SSH private key format: $OWNER_KEY_FILE${NC}"
        exit 1
    fi
    
    echo "Using existing owner key: $OWNER_KEY_FILE"
    cp "$OWNER_KEY_FILE" owner_key
    if [ -f "${OWNER_KEY_FILE}.pub" ]; then
        cp "${OWNER_KEY_FILE}.pub" owner_key.pub
    else
        ssh-keygen -y -f owner_key > owner_key.pub
    fi
    OWNER_KEY_FINGERPRINT=$(ssh-keygen -lf owner_key.pub | awk '{print $2}')
    echo -e "${GREEN}✅ Using existing owner key (fingerprint: $OWNER_KEY_FINGERPRINT)${NC}"
else
    echo "Generating new owner key for multi-node management..."
    ssh-keygen -t ed25519 -f owner_key -N "" -C "syntropy-owner-$(date +%Y%m%d)-$(hostname)"
    OWNER_KEY_FINGERPRINT=$(ssh-keygen -lf owner_key.pub | awk '{print $2}')
    echo -e "${GREEN}✅ New owner key generated (fingerprint: $OWNER_KEY_FINGERPRINT)${NC}"
    echo -e "${YELLOW}💡 Save this key for future nodes: $(pwd)/owner_key${NC}"
fi

# Generate unique community key for this specific node
COMMUNITY_KEY_NAME="community-$LOCATION_NODE_ID"
ssh-keygen -t ed25519 -f "$COMMUNITY_KEY_NAME" -N "" -C "syntropy-community-$LOCATION_NODE_ID"
COMMUNITY_KEY_FINGERPRINT=$(ssh-keygen -lf "${COMMUNITY_KEY_NAME}.pub" | awk '{print $2}')

# Generate node name if not provided
if [ -z "$NODE_NAME" ]; then
    NODE_NAME="syntropy-$LOCATION_NODE_ID"
fi

# Set default description with location info
if [ -z "$NODE_DESCRIPTION" ]; then
    NODE_DESCRIPTION="Syntropy Cooperative Grid Node in $DETECTED_CITY, $DETECTED_COUNTRY ($COORDINATES)"
fi

echo -e "${BLUE}[3/8] Creating comprehensive node metadata...${NC}"

# Create enhanced node metadata with structured geographic information
NODE_METADATA_FILE="$NODES_DIR/${NODE_NAME}.json"
cat > "$NODE_METADATA_FILE" << METADATA_EOF
{
  "metadata_version": "1.1",
  "node_info": {
    "node_id": "$LOCATION_NODE_ID",
    "node_name": "$NODE_NAME",
    "description": "$NODE_DESCRIPTION",
    "creation_time": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
    "creator_system": "$(hostname)",
    "creator_user": "$USER",
    "platform_version": "0.1.0-genesis",
    "platform_type": "syntropy_cooperative_grid"
  },
  "geographic_info": {
    "coordinates": {
      "latitude": $(echo "$COORDINATES" | cut -d',' -f1),
      "longitude": $(echo "$COORDINATES" | cut -d',' -f2),
      "formatted": "$COORDINATES"
    },
    "location": {
      "city": "$DETECTED_CITY",
      "country": "$DETECTED_COUNTRY",
      "timezone": "$(timedatectl show --property=Timezone --value 2>/dev/null || echo "UTC")"
    },
    "detection": {
      "method": "$DETECTION_METHOD",
      "timestamp": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
      "accuracy": "$(case "$DETECTION_METHOD" in *ip_geolocation*) echo "high" ;; *timezone*) echo "medium" ;; *manual*) echo "exact" ;; *) echo "low" ;; esac)"
    },
    "location_id": "$LOCATION_NODE_ID"
  },
  "security": {
    "owner_key": {
      "fingerprint": "$OWNER_KEY_FINGERPRINT",
      "algorithm": "ed25519",
      "created": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
      "local_path": "$KEYS_DIR/${NODE_NAME}_owner.key",
      "reused": $([ -n "$OWNER_KEY_FILE" ] && echo "true" || echo "false")
    },
    "community_key": {
      "fingerprint": "$COMMUNITY_KEY_FINGERPRINT",
      "algorithm": "ed25519",
      "created": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
      "local_path": "$KEYS_DIR/${NODE_NAME}_community.key"
    },
    "ssh_access": {
      "user": "admin",
      "port": 22,
      "authentication": "key_only",
      "password_disabled": true
    }
  },
  "usb_creation": {
    "created_on": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
    "usb_device": "$USB_DEVICE",
    "creator_ip": "$(timeout 5 curl -s ifconfig.me 2>/dev/null || echo 'unknown')",
    "iso_version": "ubuntu-22.04.3-live-server-amd64",
    "script_version": "enhanced-v1.1"
  },
  "expected_installation": {
    "hostname": "$NODE_NAME",
    "platform_capabilities": [
      "container_orchestration",
      "resource_sharing", 
      "cooperative_computing",
      "distributed_storage",
      "service_mesh",
      "universal_applications"
    ],
    "auto_configuration": true,
    "hardware_adaptive": true
  },
  "management": {
    "status": "usb_created",
    "installation_complete": false,
    "ssh_tested": false,
    "last_contact": null,
    "manager_system": "$(hostname)",
    "manager_user": "$USER",
    "notes": []
  },
  "network": {
    "ip_address": null,
    "interfaces": null,
    "discovered_peers": [],
    "connectivity_status": "unknown",
    "discovery_methods": ["automatic_scan", "manual_entry"]
  },
  "hardware": {
    "cpu_cores": null,
    "ram_gb": null,
    "storage_gb": null,
    "architecture": null,
    "classification": null,
    "capabilities": [],
    "detection_pending": true
  },
  "platform": {
    "universal_support": {
      "scientific_computing": ["fortran", "python", "r", "julia", "matlab"],
      "web_applications": ["nodejs", "python", "java", "go", "php"],
      "machine_learning": ["tensorflow", "pytorch", "scikit-learn", "keras"],
      "databases": ["postgresql", "mongodb", "redis", "mysql", "cassandra"],
      "custom_applications": "any_containerized_application"
    },
    "templates_available": [
      "batch-job-template.yaml",
      "web-service-template.yaml", 
      "persistent-service-template.yaml",
      "ml-training-template.yaml"
    ]
  }
}
METADATA_EOF

echo -e "${GREEN}✅ Enhanced node metadata created: $NODE_METADATA_FILE${NC}"

# Setup local SSH key management for easy access
mkdir -p "$KEYS_DIR"
cp owner_key "$KEYS_DIR/${NODE_NAME}_owner.key"
cp owner_key.pub "$KEYS_DIR/${NODE_NAME}_owner.pub"
cp "$COMMUNITY_KEY_NAME" "$KEYS_DIR/${NODE_NAME}_community.key"
cp "${COMMUNITY_KEY_NAME}.pub" "$KEYS_DIR/${NODE_NAME}_community.pub"
chmod 600 "$KEYS_DIR/${NODE_NAME}_owner.key"
chmod 600 "$KEYS_DIR/${NODE_NAME}_community.key"

echo -e "${BLUE}[4/8] Downloading Ubuntu Server ISO...${NC}"
ISO_FILE="ubuntu-22.04.3-live-server-amd64.iso"
ISO_URL="https://releases.ubuntu.com/22.04/ubuntu-22.04.3-live-server-amd64.iso"
ISO_SHA256="a4acfda10b18da50e2ec50ccaf860d7f20b389df8765611142305c0e911d16fd"

if [ ! -f "$ISO_FILE" ]; then
    echo "Downloading Ubuntu 22.04.3 Server (~1.5GB)..."
    wget --progress=bar:force:noscroll -c -O "$ISO_FILE" "$ISO_URL"
    
    echo "Verifying ISO integrity..."
    DOWNLOADED_SHA256=$(sha256sum "$ISO_FILE" | cut -d' ' -f1)
    if [ "$DOWNLOADED_SHA256" != "$ISO_SHA256" ]; then
        echo -e "${RED}ISO checksum verification failed!${NC}"
        exit 1
    fi
fi

echo -e "${BLUE}[5/8] Preparing USB device...${NC}"
sudo umount ${USB_DEVICE}* 2>/dev/null || true
sudo parted -s $USB_DEVICE mklabel msdos
sudo parted -s $USB_DEVICE mkpart primary fat32 1MiB 100%
sudo parted -s $USB_DEVICE set 1 boot on
sudo mkfs.fat -F32 -n "SYNTROPY" ${USB_DEVICE}1

echo -e "${BLUE}[6/8] Installing Ubuntu and creating enhanced configuration...${NC}"
USB_MOUNT="/mnt/syntropy-usb"
sudo mkdir -p "$USB_MOUNT"
sudo mount ${USB_DEVICE}1 "$USB_MOUNT"

# Extract ISO
mkdir -p iso-mount
sudo mount -o loop "$ISO_FILE" iso-mount
sudo cp -r iso-mount/* "$USB_MOUNT/"
sudo umount iso-mount

echo -e "${BLUE}[6/8] Creating enhanced cloud-init configuration...${NC}"

# Create enhanced user-data with pre-configured values
sudo tee "$USB_MOUNT/user-data" > /dev/null << USER_DATA_EOF
#cloud-config
# Syntropy Cooperative Grid - Enhanced Pre-configured Installation

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
        "eth*":
          dhcp4: true
          dhcp6: false
        "enp*":
          dhcp4: true
          dhcp6: false

  identity:
    hostname: $NODE_NAME
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

  late-commands:
    # Create directory structure
    - curtin in-target -- mkdir -p /opt/syntropy/{identity,platform,scripts,metadata}
    
    # Install pre-configured keys
    - |
      curtin in-target -- bash -c '
      mkdir -p /opt/syntropy/identity/{owner,community}
      
      # Install owner key
      cat > /opt/syntropy/identity/owner/private.key << "OWNER_KEY_EOF"
$(cat owner_key)
OWNER_KEY_EOF
      
      cat > /opt/syntropy/identity/owner/public.key << "OWNER_PUB_EOF"
$(cat owner_key.pub)
OWNER_PUB_EOF
      
      # Install community key
      cat > /opt/syntropy/identity/community/private.key << "COMMUNITY_KEY_EOF"
$(cat $COMMUNITY_KEY_NAME)
COMMUNITY_KEY_EOF
      
      cat > /opt/syntropy/identity/community/public.key << "COMMUNITY_PUB_EOF"
$(cat ${COMMUNITY_KEY_NAME}.pub)
COMMUNITY_PUB_EOF
      
      # Set permissions
      chmod 600 /opt/syntropy/identity/owner/private.key
      chmod 600 /opt/syntropy/identity/community/private.key
      chmod 644 /opt/syntropy/identity/owner/public.key
      chmod 644 /opt/syntropy/identity/community/public.key
      chown -R admin:admin /opt/syntropy/identity
      
      # Configure SSH
      mkdir -p /home/admin/.ssh
      cp /opt/syntropy/identity/owner/public.key /home/admin/.ssh/authorized_keys
      chmod 600 /home/admin/.ssh/authorized_keys
      chown admin:admin /home/admin/.ssh/authorized_keys
      '
    
    # Create node metadata
    - |
      curtin in-target -- bash -c '
      # Detect hardware
      CPU_CORES=\$(nproc)
      RAM_GB=\$(free -g | awk "/^Mem:/{\print \\\$2}")
      STORAGE_GB=\$(df / --output=avail -BG 2>/dev/null | tail -1 | sed "s/G//" || echo "50")
      
      # Hardware classification
      if [ \$RAM_GB -le 2 ]; then
        HW_CLASS="edge"
      elif [ \$RAM_GB -le 8 ]; then
        HW_CLASS="home-server"
      elif [ \$RAM_GB -le 32 ]; then
        HW_CLASS="server"
      else
        HW_CLASS="high-end-server"
      fi
      
      # Create comprehensive metadata
      cat > /opt/syntropy/metadata/node.json << "NODE_METADATA_EOF"
{
  "node_info": {
    "node_id": "$NODE_ID",
    "node_name": "$NODE_NAME",
    "hostname": "$(hostname)",
    "description": "$NODE_DESCRIPTION",
    "installation_time": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
    "version": "0.1.0-genesis"
  },
  "geographic_info": {
    "coordinates": "$COORDINATES",
    "detection_method": "$LOCATION_METHOD",
    "timezone": "$(timedatectl show --property=Timezone --value 2>/dev/null || echo UTC)"
  },
  "hardware": {
    "cpu_cores": \$CPU_CORES,
    "ram_gb": \$RAM_GB,
    "storage_gb": \$STORAGE_GB,
    "architecture": "$(uname -m)",
    "classification": "\$HW_CLASS"
  },
  "network": {
    "interfaces": "$(ip link show | grep "^[0-9]" | awk -F: "{print \\\$2}" | grep -v lo | tr "\n" " " | sed "s/ \$//")",
    "ip_address": "$(hostname -I | awk "{print \\\$1}")"
  },
  "security": {
    "owner_key_fingerprint": "$(ssh-keygen -lf /opt/syntropy/identity/owner/public.key | awk "{print \\\$2}")",
    "community_key_fingerprint": "$(ssh-keygen -lf /opt/syntropy/identity/community/public.key | awk "{print \\\$2}")",
    "ssh_port": 22
  },
  "platform": {
    "type": "syntropy_cooperative_grid",
    "capabilities": ["container_orchestration", "resource_sharing", "cooperative_computing"],
    "status": "installed"
  }
}
NODE_METADATA_EOF
      '
    
    # Configure services
    - curtin in-target -- systemctl enable ssh docker prometheus-node-exporter
    - curtin in-target -- ufw default deny incoming
    - curtin in-target -- ufw default allow outgoing
    - curtin in-target -- ufw allow ssh
    - curtin in-target -- ufw --force enable

  power_state:
    mode: reboot
    timeout: 30
USER_DATA_EOF

# Create meta-data
sudo tee "$USB_MOUNT/meta-data" > /dev/null << 'EOF'
instance-id: syntropy-node-auto
local-hostname: syntropy-node
EOF

echo -e "${BLUE}[7/8] Creating management tools...${NC}"

# Create node connection script
cat > "$HOME/.syntropy/connect-${NODE_NAME}.sh" << CONNECT_EOF
#!/bin/bash

# Syntropy Node Connection Script for $NODE_NAME
# Auto-generated on $(date)

NODE_NAME="$NODE_NAME"
KEY_FILE="$HOME/.syntropy/keys/\${NODE_NAME}_owner.key"
METADATA_FILE="$HOME/.syntropy/nodes/\${NODE_NAME}.json"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo "=== SYNTROPY NODE CONNECTION ==="
echo "Node: \$NODE_NAME"

# Check if metadata exists
if [ ! -f "\$METADATA_FILE" ]; then
    echo -e "\${RED}Error: Node metadata not found: \$METADATA_FILE\${NC}"
    exit 1
fi

# Check if key exists
if [ ! -f "\$KEY_FILE" ]; then
    echo -e "\${RED}Error: SSH key not found: \$KEY_FILE\${NC}"
    exit 1
fi

# Discover node IP
echo "Discovering node IP address..."
NODE_IP=""

# Method 1: Try known IP from metadata
KNOWN_IP=\$(jq -r '.network.ip_address // empty' "\$METADATA_FILE" 2>/dev/null)
if [ -n "\$KNOWN_IP" ] && [ "\$KNOWN_IP" != "null" ]; then
    if ssh -i "\$KEY_FILE" -o ConnectTimeout=5 -o BatchMode=yes admin@\$KNOWN_IP exit 2>/dev/null; then
        NODE_IP="\$KNOWN_IP"
    fi
fi

# Method 2: Network scan for hostname
if [ -z "\$NODE_IP" ]; then
    echo "Scanning network for \$NODE_NAME..."
    for ip in \$(nmap -sn 192.168.1.0/24 2>/dev/null | grep -E "Nmap scan report|MAC Address" | grep -B1 "\$NODE_NAME" | grep "Nmap scan report" | awk '{print \$5}'); do
        if ssh -i "\$KEY_FILE" -o ConnectTimeout=5 -o BatchMode=yes admin@\$ip exit 2>/dev/null; then
            NODE_IP="\$ip"
            break
        fi
    done
fi

# Method 3: Interactive IP entry
if [ -z "\$NODE_IP" ]; then
    echo -e "\${YELLOW}Could not automatically discover node IP.\${NC}"
    echo "Available IPs on network:"
    nmap -sn 192.168.1.0/24 2>/dev/null | grep "Nmap scan report" | awk '{print \$5}'
    echo ""
    read -p "Enter node IP address: " NODE_IP
fi

if [ -z "\$NODE_IP" ]; then
    echo -e "\${RED}No IP address provided\${NC}"
    exit 1
fi

# Test connection
echo "Testing SSH connection to \$NODE_IP..."
if ssh -i "\$KEY_FILE" -o ConnectTimeout=10 -o BatchMode=yes admin@\$NODE_IP exit; then
    echo -e "\${GREEN}✅ SSH connection successful\${NC}"
    
    # Update metadata with current IP
    jq ".network.ip_address = \"\$NODE_IP\" | .management.last_contact = \"\$(date -u +%Y-%m-%dT%H:%M:%SZ)\" | .management.ssh_tested = true" "\$METADATA_FILE" > "\$METADATA_FILE.tmp" && mv "\$METADATA_FILE.tmp" "\$METADATA_FILE"
    
    echo ""
    echo "Connection details:"
    echo "• Host: \$NODE_IP"
    echo "• User: admin"
    echo "• Key: \$KEY_FILE"
    echo ""
    echo "Commands:"
    echo "• Connect: ssh -i \$KEY_FILE admin@\$NODE_IP"
    echo "• Status: ssh -i \$KEY_FILE admin@\$NODE_IP 'cat /opt/syntropy/metadata/node.json | jq .'"
    echo "• Copy file: scp -i \$KEY_FILE file.txt admin@\$NODE_IP:/tmp/"
    echo ""
    
    # Offer direct connection
    read -p "Connect now? (y/N): " -n 1 -r
    echo
    if [[ \$REPLY =~ ^[Yy]\$ ]]; then
        ssh -i "\$KEY_FILE" admin@\$NODE_IP
    fi
else
    echo -e "\${RED}❌ SSH connection failed\${NC}"
    echo "Troubleshooting:"
    echo "1. Ensure the node has finished installing and rebooted"
    echo "2. Check if the IP address is correct"
    echo "3. Verify network connectivity"
    echo "4. Try manual connection: ssh -i \$KEY_FILE admin@\$NODE_IP"
fi
CONNECT_EOF

chmod +x "$HOME/.syntropy/connect-${NODE_NAME}.sh"

# Create management summary
cat > "$HOME/.syntropy/nodes/${NODE_NAME}_summary.md" << SUMMARY_EOF
# Syntropy Node: $NODE_NAME

## Node Information
- **Node ID**: $NODE_ID
- **Name**: $NODE_NAME
- **Description**: $NODE_DESCRIPTION
- **Created**: $(date)
- **Coordinates**: $COORDINATES

## Connection Details
- **SSH Key**: $HOME/.syntropy/keys/${NODE_NAME}_owner.key
- **Username**: admin
- **SSH Port**: 22

## Quick Commands
\`\`\`bash
# Connect to node
$HOME/.syntropy/connect-${NODE_NAME}.sh

# Direct SSH (after getting IP)
ssh -i $HOME/.syntropy/keys/${NODE_NAME}_owner.key admin@<NODE_IP>

# Check node status
ssh -i $HOME/.syntropy/keys/${NODE_NAME}_owner.key admin@<NODE_IP> 'cat /opt/syntropy/metadata/node.json | jq .'

# Copy files to node
scp -i $HOME/.syntropy/keys/${NODE_NAME}_owner.key file.txt admin@<NODE_IP>:/tmp/
\`\`\`

## Files Created
- Node metadata: $HOME/.syntropy/nodes/${NODE_NAME}.json
- SSH keys: $HOME/.syntropy/keys/${NODE_NAME}_*.key
- Connection script: $HOME/.syntropy/connect-${NODE_NAME}.sh

## Next Steps
1. Boot the created USB on target hardware
2. Wait for installation to complete (~20-30 minutes)
3. Run connection script: $HOME/.syntropy/connect-${NODE_NAME}.sh
4. Verify node is working and accessible
SUMMARY_EOF

echo -e "${BLUE}[8/8] Finalizing USB...${NC}"

# Add documentation to USB
sudo mkdir -p "$USB_MOUNT/syntropy-docs"
sudo cp "$HOME/.syntropy/nodes/${NODE_NAME}_summary.md" "$USB_MOUNT/syntropy-docs/"

# Sync and unmount
sync
sudo umount "$USB_MOUNT"
sudo rmdir "$USB_MOUNT"

# Cleanup
cd /
rm -rf "$WORK_DIR"

echo -e "${GREEN}✅ Enhanced Syntropy USB created successfully!${NC}"
echo ""
echo -e "${PURPLE}═══ NODE MANAGEMENT SUMMARY ═══${NC}"
echo "🏷️  Node Name: $NODE_NAME"
echo "🆔 Node ID: $NODE_ID"
echo "📍 Coordinates: $COORDINATES"
echo "🔑 SSH Key: $HOME/.syntropy/keys/${NODE_NAME}_owner.key"
echo "📄 Metadata: $HOME/.syntropy/nodes/${NODE_NAME}.json"
echo ""
echo -e "${CYAN}═══ INSTALLATION PROCESS ═══${NC}"
echo "1. 🔌 Insert USB into target hardware"
echo "2. ⚙️  Configure BIOS to boot from USB"
echo "3. 🚀 Boot and wait for automatic installation (~30 minutes)"
echo "4. 🔍 Run connection script: $HOME/.syntropy/connect-${NODE_NAME}.sh"
echo ""
echo -e "${YELLOW}═══ MANAGEMENT FILES CREATED ═══${NC}"
echo "📁 Node management directory: $HOME/.syntropy/"
echo "├── 📄 nodes/${NODE_NAME}.json (metadata)"
echo "├── 🔑 keys/${NODE_NAME}_owner.key (SSH access)"
echo "├── 🔑 keys/${NODE_NAME}_community.key (network identity)"
echo "├── 📜 connect-${NODE_NAME}.sh (connection script)"
echo "└── 📋 nodes/${NODE_NAME}_summary.md (documentation)"
echo ""
echo -e "${GREEN}Ready for hardware installation! Your management setup is complete. 🚀${NC}"