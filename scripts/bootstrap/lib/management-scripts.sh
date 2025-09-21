#!/bin/bash

# Syntropy Cooperative Grid - Node Management Scripts Generator
# Version: 2.0.0

# Create comprehensive node metadata
create_node_metadata() {
    local node_name="$1"
    local location_node_id="$2"
    local coordinates="$3"
    local detection_method="$4"
    local detected_city="$5"
    local detected_country="$6"
    local node_description="$7"
    
    log INFO "Creating comprehensive node metadata..."
    
    # Ensure metadata directory exists
    mkdir -p "$NODES_DIR"
    
    local metadata_file="$NODES_DIR/${node_name}.json"
    
    # Create comprehensive metadata file
    cat > "$metadata_file" << METADATA_EOF
{
  "metadata_version": "2.0",
  "node_info": {
    "node_id": "$location_node_id",
    "node_name": "$node_name",
    "description": "$node_description",
    "creation_time": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
    "platform_version": "2.0.0-genesis",
    "status": "usb_created"
  },
  "geographic_info": {
    "coordinates": {
      "latitude": $(echo "$coordinates" | cut -d',' -f1),
      "longitude": $(echo "$coordinates" | cut -d',' -f2),
      "formatted": "$coordinates"
    },
    "location": {
      "city": "$detected_city",
      "country": "$detected_country",
      "timezone": "$(timedatectl show --property=Timezone --value 2>/dev/null || echo UTC)"
    },
    "detection": {
      "method": "$detection_method",
      "timestamp": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
      "accuracy": "$(case "$detection_method" in *ip_geolocation*) echo "high" ;; *timezone*) echo "medium" ;; *manual*) echo "exact" ;; *) echo "low" ;; esac)"
    },
    "location_id": "$location_node_id"
  },
  "security": {
    "owner_key_fingerprint": "$OWNER_FINGERPRINT",
    "community_key_fingerprint": "$COMMUNITY_FINGERPRINT",
    "ssh_port": 22,
    "authentication_method": "key_only",
    "key_created": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
    "local_key_path": "$KEYS_DIR/${node_name}_owner.key",
    "key_reused": $([ -n "$OWNER_KEY_FILE" ] && echo "true" || echo "false")
  },
  "platform": {
    "type": "syntropy_cooperative_grid",
    "capabilities": [
      "container_orchestration",
      "resource_sharing",
      "cooperative_computing",
      "distributed_storage",
      "service_mesh",
      "universal_applications"
    ],
    "expected_services": {
      "docker": "enabled",
      "ssh": "enabled",
      "prometheus_exporter": "enabled",
      "firewall": "enabled"
    }
  },
  "usb_creation": {
    "created_on": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
    "created_by": "$(whoami)@$(hostname)",
    "iso_version": "$ISO_FILE",
    "script_version": "2.0.0"
  },
  "management": {
    "status": "usb_created",
    "installation_complete": false,
    "ssh_tested": false,
    "last_contact": null,
    "ip_address": null,
    "first_boot": null
  }
}
METADATA_EOF
    
    log SUCCESS "Node metadata created: $metadata_file"
}

# Create node management scripts
create_node_management_scripts() {
    local node_name="$1"
    local location_node_id="$2"
    local coordinates="$3"
    local detected_city="$4"
    local detected_country="$5"
    local node_description="$6"
    
    log INFO "Creating node management scripts..."
    
    # Create individual connection script
    create_connection_script "$node_name"
    
    # Create discovery script
    create_discovery_script "$node_name"
    
    # Create status monitoring script
    create_status_script "$node_name"
    
    # Create comprehensive manager script if it doesn't exist
    create_main_manager_script
    
    # Store keys locally
    store_keys_locally "$node_name"
    
    # Generate security summary
    generate_security_summary "$node_name"
    
    log SUCCESS "Management scripts created for node: $node_name"
}

# Create individual node connection script
create_connection_script() {
    local node_name="$1"
    local script_path="$HOME/.syntropy/connect-${node_name}.sh"
    
    cat > "$script_path" << 'CONNECT_SCRIPT_EOF'
#!/bin/bash

# Auto-generated connection script for node: %NODE_NAME%
# Created: %CREATION_DATE%

NODE_NAME="%NODE_NAME%"
KEY_FILE="$HOME/.syntropy/keys/${NODE_NAME}_owner.key"
METADATA_FILE="$HOME/.syntropy/nodes/${NODE_NAME}.json"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log() {
    local level="$1"
    shift
    case "$level" in
        INFO) echo -e "${BLUE}[INFO]${NC} $*" ;;
        SUCCESS) echo -e "${GREEN}[SUCCESS]${NC} $*" ;;
        WARN) echo -e "${YELLOW}[WARN]${NC} $*" ;;
        ERROR) echo -e "${RED}[ERROR]${NC} $*" ;;
    esac
}

# Discover node IP address
discover_node_ip() {
    log INFO "Discovering IP address for node: $NODE_NAME"
    
    # Method 1: Check if we have a cached IP
    if [ -f "$METADATA_FILE" ]; then
        local cached_ip=$(jq -r '.management.ip_address // empty' "$METADATA_FILE" 2>/dev/null)
        if [ -n "$cached_ip" ] && [ "$cached_ip" != "null" ]; then
            log INFO "Trying cached IP: $cached_ip"
            if test_ssh_connection "$cached_ip"; then
                echo "$cached_ip"
                return 0
            fi
        fi
    fi
    
    # Method 2: Network scan
    log INFO "Scanning local network for node..."
    local network=$(ip route | grep "scope link" | head -1 | awk '{print $1}' | head -1)
    if [ -n "$network" ]; then
        log INFO "Scanning network: $network"
        
        # Use nmap to find hosts with SSH
        local ssh_hosts=$(nmap -p 22 --open -T4 "$network" 2>/dev/null | grep -B4 "22/tcp open" | grep "Nmap scan report" | awk '{print $5}')
        
        for ip in $ssh_hosts; do
            log INFO "Testing SSH connection to: $ip"
            if test_ssh_connection "$ip"; then
                # Update metadata with found IP
                update_node_ip "$ip"
                echo "$ip"
                return 0
            fi
        done
    fi
    
    # Method 3: Manual input
    echo ""
    log WARN "Could not automatically discover node IP"
    echo "Please enter the IP address manually, or check:"
    echo "1. Node is powered on and connected to network"
    echo "2. Installation completed successfully"
    echo "3. Network allows SSH connections"
    echo ""
    read -p "Enter IP address (or 'q' to quit): " manual_ip
    
    if [ "$manual_ip" = "q" ]; then
        exit 0
    fi
    
    if test_ssh_connection "$manual_ip"; then
        update_node_ip "$manual_ip"
        echo "$manual_ip"
        return 0
    else
        log ERROR "Could not connect to manually entered IP: $manual_ip"
        return 1
    fi
}

# Test SSH connection
test_ssh_connection() {
    local ip="$1"
    timeout 5 ssh -i "$KEY_FILE" -o ConnectTimeout=3 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=quiet admin@"$ip" 'echo "connected"' >/dev/null 2>&1
}

# Update node metadata with IP
update_node_ip() {
    local ip="$1"
    if [ -f "$METADATA_FILE" ] && command -v jq >/dev/null 2>&1; then
        local temp_file=$(mktemp)
        jq --arg ip "$ip" --arg timestamp "$(date -u +%Y-%m-%dT%H:%M:%SZ)" '
            .management.ip_address = $ip |
            .management.last_contact = $timestamp |
            .management.ssh_tested = true
        ' "$METADATA_FILE" > "$temp_file" && mv "$temp_file" "$METADATA_FILE"
        log SUCCESS "Updated node metadata with IP: $ip"
    fi
}

# Main connection function
connect_to_node() {
    echo "Syntropy Node Connection - $NODE_NAME"
    echo "========================================"
    
    if [ ! -f "$KEY_FILE" ]; then
        log ERROR "SSH key not found: $KEY_FILE"
        exit 1
    fi
    
    local node_ip=$(discover_node_ip)
    if [ $? -ne 0 ] || [ -z "$node_ip" ]; then
        log ERROR "Could not discover or connect to node"
        exit 1
    fi
    
    log SUCCESS "Connecting to $NODE_NAME at $node_ip"
    echo ""
    
    # Connect with proper SSH options
    ssh -i "$KEY_FILE" \
        -o ConnectTimeout=10 \
        -o StrictHostKeyChecking=no \
        -o UserKnownHostsFile=/dev/null \
        admin@"$node_ip"
}

# Show usage help
show_help() {
    echo "Usage: $0 [options]"
    echo ""
    echo "Options:"
    echo "  --discover-only    Only discover IP, don't connect"
    echo "  --status          Show node status"
    echo "  --help            Show this help"
    echo ""
    echo "Node: $NODE_NAME"
    echo "Key:  $KEY_FILE"
}

# Parse arguments
case "${1:-}" in
    --discover-only)
        discover_node_ip
        ;;
    --status)
        if command -v jq >/dev/null 2>&1 && [ -f "$METADATA_FILE" ]; then
            echo "Node Status for: $NODE_NAME"
            echo "=========================="
            jq -r '
                "Status: " + .management.status + 
                "\nLast Contact: " + (.management.last_contact // "never") +
                "\nIP Address: " + (.management.ip_address // "unknown") +
                "\nSSH Tested: " + (.management.ssh_tested | tostring)
            ' "$METADATA_FILE"
        else
            log WARN "Cannot show status - jq not available or metadata missing"
        fi
        ;;
    --help)
        show_help
        ;;
    *)
        connect_to_node
        ;;
esac
CONNECT_SCRIPT_EOF

    # Replace placeholders
    sed -i.bak "s/%NODE_NAME%/$node_name/g" "$script_path"
    sed -i.bak "s/%CREATION_DATE%/$(date)/g" "$script_path"
    rm -f "${script_path}.bak"
    
    chmod +x "$script_path"
    log SUCCESS "Connection script created: $script_path"
}

# Create network discovery script
create_discovery_script() {
    local node_name="$1"
    local script_path="$HOME/.syntropy/discover-${node_name}.sh"
    
    cat > "$script_path" << 'DISCOVERY_SCRIPT_EOF'
#!/bin/bash

# Network discovery script for node: %NODE_NAME%
# Finds Syntropy nodes on the local network

NODE_NAME="%NODE_NAME%"
METADATA_FILE="$HOME/.syntropy/nodes/${NODE_NAME}.json"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo "Syntropy Network Discovery - $NODE_NAME"
echo "========================================"

# Discover local networks
echo -e "${BLUE}Scanning for Syntropy nodes...${NC}"

local_networks=$(ip route | grep -E "192\.168\.|10\.|172\." | grep "scope link" | awk '{print $1}' | sort -u)

for network in $local_networks; do
    echo -e "${YELLOW}Scanning network: $network${NC}"
    
    # Find hosts with SSH and Prometheus node exporter
    nmap -p 22,9100 --open -T4 "$network" 2>/dev/null | \
    awk '/Nmap scan report/{ip=$5} /22\/tcp open/{ssh=1} /9100\/tcp open/{prom=1} ssh&&prom{print ip; ssh=0; prom=0}' | \
    while read -r host_ip; do
        echo -e "${GREEN}Found potential Syntropy node: $host_ip${NC}"
        
        # Test if it's actually a Syntropy node
        if timeout 3 curl -s "http://$host_ip:9100/metrics" | grep -q "node_exporter"; then
            echo "  ✓ Confirmed Syntropy node with metrics endpoint"
            
            # Try to get hostname
            hostname=$(timeout 3 ssh -o ConnectTimeout=2 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=quiet admin@"$host_ip" 'hostname' 2>/dev/null || echo "unknown")
            echo "  ✓ Hostname: $hostname"
            
            # Update metadata if this is our node
            if [ "$hostname" = "$NODE_NAME" ]; then
                echo "  ✓ This is our target node!"
                if command -v jq >/dev/null 2>&1 && [ -f "$METADATA_FILE" ]; then
                    temp_file=$(mktemp)
                    jq --arg ip "$host_ip" --arg timestamp "$(date -u +%Y-%m-%dT%H:%M:%SZ)" '
                        .management.ip_address = $ip |
                        .management.last_contact = $timestamp |
                        .management.ssh_tested = true
                    ' "$METADATA_FILE" > "$temp_file" && mv "$temp_file" "$METADATA_FILE"
                    echo "  ✓ Updated metadata with discovered IP"
                fi
            fi
        fi
        echo ""
    done
done

echo "Discovery complete!"
DISCOVERY_SCRIPT_EOF

    # Replace placeholders
    sed -i.bak "s/%NODE_NAME%/$node_name/g" "$script_path"
    rm -f "${script_path}.bak"
    
    chmod +x "$script_path"
    log SUCCESS "Discovery script created: $script_path"
}

# Create status monitoring script
create_status_script() {
    local node_name="$1"
    local script_path="$HOME/.syntropy/status-${node_name}.sh"
    
    cat > "$script_path" << 'STATUS_SCRIPT_EOF'
#!/bin/bash

# Status monitoring script for node: %NODE_NAME%

NODE_NAME="%NODE_NAME%"
KEY_FILE="$HOME/.syntropy/keys/${NODE_NAME}_owner.key"
METADATA_FILE="$HOME/.syntropy/nodes/${NODE_NAME}.json"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

show_local_status() {
    echo "Local Metadata Status"
    echo "===================="
    
    if [ -f "$METADATA_FILE" ] && command -v jq >/dev/null 2>&1; then
        jq -r '
            "Node ID: " + .node_info.node_id + 
            "\nNode Name: " + .node_info.node_name +
            "\nStatus: " + .management.status +
            "\nCreated: " + .node_info.creation_time +
            "\nLocation: " + .geographic_info.location.city + ", " + .geographic_info.location.country +
            "\nCoordinates: " + .geographic_info.coordinates.formatted +
            "\nLast Contact: " + (.management.last_contact // "never") +
            "\nIP Address: " + (.management.ip_address // "unknown")
        ' "$METADATA_FILE"
    else
        echo "Metadata file not found or jq not available"
    fi
}

show_remote_status() {
    local ip="$1"
    
    echo ""
    echo "Remote System Status"
    echo "==================="
    
    if ! timeout 5 ssh -i "$KEY_FILE" -o ConnectTimeout=3 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=quiet admin@"$ip" 'echo "connected"' >/dev/null 2>&1; then
        echo -e "${RED}Cannot connect to node via SSH${NC}"
        return 1
    fi
    
    # Get system information
    ssh -i "$KEY_FILE" -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=quiet admin@"$ip" '
        echo "Hostname: $(hostname)"
        echo "Uptime: $(uptime -p 2>/dev/null || uptime)"
        echo "Load: $(cat /proc/loadavg | cut -d" " -f1-3)"
        echo "Memory: $(free -h | grep Mem | awk "{print \$3\"/\"\$2}")"
        echo "Disk: $(df -h / | tail -1 | awk "{print \$3\"/\"\$2\" (\"\$5\" used)\"}")"
        echo ""
        echo "Docker Status:"
        sudo systemctl is-active docker || echo "Docker not running"
        echo ""
        echo "Active Containers:"
        sudo docker ps --format "table {{.Names}}\t{{.Status}}" 2>/dev/null || echo "Cannot access docker"
        echo ""
        echo "Syntropy Services:"
        sudo systemctl is-active prometheus-node-exporter 2>/dev/null || echo "Node exporter not running"
        echo ""
        echo "Firewall Status:"
        sudo ufw status 2>/dev/null || echo "UFW not available"
    ' 2>/dev/null
}

# Main status function
show_status() {
    echo "Syntropy Node Status - $NODE_NAME"
    echo "=================================="
    echo ""
    
    show_local_status
    
    # Try to get IP and show remote status
    if [ -f "$METADATA_FILE" ] && command -v jq >/dev/null 2>&1; then
        local ip=$(jq -r '.management.ip_address // empty' "$METADATA_FILE" 2>/dev/null)
        if [ -n "$ip" ] && [ "$ip" != "null" ]; then
            show_remote_status "$ip"
        else
            echo ""
            echo -e "${YELLOW}No IP address known - run discovery first${NC}"
        fi
    fi
}

show_status
STATUS_SCRIPT_EOF

    # Replace placeholders
    sed -i.bak "s/%NODE_NAME%/$node_name/g" "$script_path"
    rm -f "${script_path}.bak"
    
    chmod +x "$script_path"
    log SUCCESS "Status script created: $script_path"
}

# Create main manager script
create_main_manager_script() {
    local manager_script="$HOME/.syntropy/syntropy-manager.sh"
    
    if [ -f "$manager_script" ]; then
        log DEBUG "Manager script already exists, skipping creation"
        return 0
    fi
    
    cat > "$manager_script" << 'MANAGER_SCRIPT_EOF'
#!/bin/bash

# Syntropy Cooperative Grid - Node Manager
# Version: 2.0.0

SYNTROPY_DIR="$HOME/.syntropy"
NODES_DIR="$SYNTROPY_DIR/nodes"
KEYS_DIR="$SYNTROPY_DIR/keys"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m'

log() {
    local level="$1"
    shift
    case "$level" in
        INFO) echo -e "${BLUE}[INFO]${NC} $*" ;;
        SUCCESS) echo -e "${GREEN}[SUCCESS]${NC} $*" ;;
        WARN) echo -e "${YELLOW}[WARN]${NC} $*" ;;
        ERROR) echo -e "${RED}[ERROR]${NC} $*" ;;
    esac
}

show_banner() {
    echo -e "${PURPLE}"
    echo "╔══════════════════════════════════════════════════════════════════════════════╗"
    echo "║                        SYNTROPY NODE MANAGER                                ║"
    echo "║                     Cooperative Grid Management                             ║"
    echo "╚══════════════════════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

list_nodes() {
    echo "Syntropy Nodes"
    echo "=============="
    
    if [ ! -d "$NODES_DIR" ] || [ -z "$(ls -A "$NODES_DIR" 2>/dev/null)" ]; then
        echo "No nodes found. Create your first node with create-syntropy-usb-enhanced.sh"
        return 0
    fi
    
    printf "%-20s %-15s %-30s %-15s %s\n" "NODE NAME" "STATUS" "LOCATION" "IP ADDRESS" "LAST CONTACT"
    printf "%-20s %-15s %-30s %-15s %s\n" "$(printf '%*s' 20 | tr ' ' '-')" "$(printf '%*s' 15 | tr ' ' '-')" "$(printf '%*s' 30 | tr ' ' '-')" "$(printf '%*s' 15 | tr ' ' '-')" "$(printf '%*s' 15 | tr ' ' '-')"
    
    for metadata_file in "$NODES_DIR"/*.json; do
        if [ -f "$metadata_file" ] && command -v jq >/dev/null 2>&1; then
            local node_name=$(jq -r '.node_info.node_name' "$metadata_file" 2>/dev/null)
            local status=$(jq -r '.management.status' "$metadata_file" 2>/dev/null)
            local city=$(jq -r '.geographic_info.location.city' "$metadata_file" 2>/dev/null)
            local country=$(jq -r '.geographic_info.location.country' "$metadata_file" 2>/dev/null)
            local ip=$(jq -r '.management.ip_address // "unknown"' "$metadata_file" 2>/dev/null)
            local last_contact=$(jq -r '.management.last_contact // "never"' "$metadata_file" 2>/dev/null)
            
            local location="$city, $country"
            if [ ${#location} -gt 28 ]; then
                location="${location:0:27}…"
            fi
            
            printf "%-20s %-15s %-30s %-15s %s\n" "$node_name" "$status" "$location" "$ip" "$last_contact"
        fi
    done
}

show_node_details() {
    local node_name="$1"
    local metadata_file="$NODES_DIR/${node_name}.json"
    
    if [ ! -f "$metadata_file" ]; then
        log ERROR "Node not found: $node_name"
        return 1
    fi
    
    if ! command -v jq >/dev/null 2>&1; then
        log ERROR "jq is required for node details"
        return 1
    fi
    
    echo "Node Details: $node_name"
    echo "========================"
    echo ""
    
    jq -r '
        "Basic Information:",
        "  Node ID: " + .node_info.node_id,
        "  Node Name: " + .node_info.node_name,
        "  Description: " + .node_info.description,
        "  Created: " + .node_info.creation_time,
        "  Platform Version: " + .node_info.platform_version,
        "",
        "Geographic Information:",
        "  Location: " + .geographic_info.location.city + ", " + .geographic_info.location.country,
        "  Coordinates: " + .geographic_info.coordinates.formatted,
        "  Timezone: " + .geographic_info.location.timezone,
        "  Detection Method: " + .geographic_info.detection.method,
        "",
        "Security:",
        "  Owner Key Fingerprint: " + .security.owner_key_fingerprint,
        "  Community Key Fingerprint: " + .security.community_key_fingerprint,
        "  SSH Port: " + (.security.ssh_port | tostring),
        "  Key Path: " + .security.local_key_path,
        "",
        "Management:",
        "  Status: " + .management.status,
        "  IP Address: " + (.management.ip_address // "unknown"),
        "  Last Contact: " + (.management.last_contact // "never"),
        "  SSH Tested: " + (.management.ssh_tested | tostring)
    ' "$metadata_file"
}

connect_to_node() {
    local node_name="$1"
    local connect_script="$SYNTROPY_DIR/connect-${node_name}.sh"
    
    if [ ! -f "$connect_script" ]; then
        log ERROR "Connection script not found for node: $node_name"
        log INFO "Available nodes: $(ls "$NODES_DIR"/*.json 2>/dev/null | xargs -I {} basename {} .json | tr '\n' ' ')"
        return 1
    fi
    
    log INFO "Connecting to node: $node_name"
    exec "$connect_script"
}

discover_nodes() {
    echo "Network Discovery"
    echo "================"
    
    if [ ! -d "$NODES_DIR" ] || [ -z "$(ls -A "$NODES_DIR" 2>/dev/null)" ]; then
        echo "No nodes to discover. Create nodes first."
        return 0
    fi
    
    for metadata_file in "$NODES_DIR"/*.json; do
        if [ -f "$metadata_file" ] && command -v jq >/dev/null 2>&1; then
            local node_name=$(jq -r '.node_info.node_name' "$metadata_file" 2>/dev/null)
            local discovery_script="$SYNTROPY_DIR/discover-${node_name}.sh"
            
            if [ -f "$discovery_script" ]; then
                echo ""
                echo "Discovering: $node_name"
                echo "========================"
                "$discovery_script"
            fi
        fi
    done
}

show_help() {
    show_banner
    echo "Usage: syntropy-manager.sh <command> [arguments]"
    echo ""
    echo "Commands:"
    echo "  list                    List all managed nodes"
    echo "  show <node>            Show detailed node information"
    echo "  connect <node>         Connect to node via SSH"
    echo "  status <node>          Show node status"
    echo "  discover               Discover all nodes on network"
    echo "  backup                 Backup all node configurations"
    echo "  help                   Show this help"
    echo ""
    echo "Examples:"
    echo "  syntropy-manager.sh list"
    echo "  syntropy-manager.sh show my-home-server"
    echo "  syntropy-manager.sh connect my-home-server"
    echo "  syntropy-manager.sh discover"
    echo ""
}

# Main command parsing
case "${1:-help}" in
    list)
        list_nodes
        ;;
    show)
        if [ -z "$2" ]; then
            log ERROR "Node name required"
            echo "Usage: syntropy-manager.sh show <node_name>"
            exit 1
        fi
        show_node_details "$2"
        ;;
    connect)
        if [ -z "$2" ]; then
            log ERROR "Node name required"
            echo "Usage: syntropy-manager.sh connect <node_name>"
            exit 1
        fi
        connect_to_node "$2"
        ;;
    status)
        if [ -z "$2" ]; then
            log ERROR "Node name required"
            echo "Usage: syntropy-manager.sh status <node_name>"
            exit 1
        fi
        status_script="$SYNTROPY_DIR/status-$2.sh"
        if [ -f "$status_script" ]; then
            "$status_script"
        else
            log ERROR "Status script not found for node: $2"
        fi
        ;;
    discover)
        discover_nodes
        ;;
    backup)
        log INFO "Creating backup of Syntropy configuration..."
        backup_dir="$SYNTROPY_DIR/backups/backup-$(date +%Y%m%d-%H%M%S)"
        mkdir -p "$backup_dir"
        cp -r "$NODES_DIR" "$KEYS_DIR" "$backup_dir/" 2>/dev/null
        log SUCCESS "Backup created: $backup_dir"
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        echo "Unknown command: $1"
        echo "Use 'syntropy-manager.sh help' for available commands"
        exit 1
        ;;
esac
MANAGER_SCRIPT_EOF

    chmod +x "$manager_script"
    log SUCCESS "Main manager script created: $manager_script"
    
    # Add to PATH if not already there
    if ! grep -q "/.syntropy" "$HOME/.bashrc" 2>/dev/null; then
        echo 'export PATH="$HOME/.syntropy:$PATH"' >> "$HOME/.bashrc"
        log INFO "Added Syntropy tools to PATH in ~/.bashrc"
    fi
}

# Finalize USB creation with comprehensive summary
finalize_usb_creation() {
    local node_name="$1"
    local location_node_id="$2"
    local coordinates="$3"
    local detected_city="$4"
    local detected_country="$5"
    
    log INFO "Finalizing USB creation..."
    
    # Create node summary document
    create_node_summary_document "$node_name" "$location_node_id" "$coordinates" \
        "$detected_city" "$detected_country"
    
    # Generate quick reference card
    create_quick_reference_card "$node_name"
    
    # Clean up sensitive data
    cleanup_sensitive_data
    
    # Show final success message and instructions
    show_final_success_message "$node_name" "$location_node_id" "$coordinates" \
        "$detected_city" "$detected_country"
}

# Create comprehensive node summary document
create_node_summary_document() {
    local node_name="$1"
    local location_node_id="$2"
    local coordinates="$3"
    local detected_city="$4"
    local detected_country="$5"
    
    local summary_file="$NODES_DIR/${node_name}_summary.md"
    
    cat > "$summary_file" << SUMMARY_EOF
# Syntropy Node: $node_name

**Created**: $(date)  
**Platform Version**: 2.0.0  
**Node ID**: $location_node_id  

## Geographic Information

- **Location**: $detected_city, $detected_country
- **Coordinates**: $coordinates
- **Timezone**: $(timedatectl show --property=Timezone --value 2>/dev/null || echo UTC)

## Security Configuration

- **Owner Key Fingerprint**: $OWNER_FINGERPRINT
- **Community Key Fingerprint**: $COMMUNITY_FINGERPRINT
- **SSH Authentication**: Key-only (passwords disabled)
- **Local Key Path**: $KEYS_DIR/${node_name}_owner.key

## Installation Details

### Hardware Requirements
- **Minimum RAM**: 4GB (8GB+ recommended)
- **Minimum Storage**: 32GB (64GB+ recommended)
- **Architecture**: x86_64 (Intel/AMD 64-bit)
- **Network**: Ethernet or WiFi with DHCP

### Installation Process
1. Insert USB into target hardware
2. Configure BIOS/UEFI to boot from USB
3. Boot from USB (installation is fully automated)
4. Wait ~30 minutes for installation completion
5. System will reboot automatically when ready

### Post-Installation Access
```bash
# Quick connection (auto-discovers IP)
$HOME/.syntropy/connect-${node_name}.sh

# Manual connection (if IP is known)
ssh -i $KEYS_DIR/${node_name}_owner.key admin@<NODE_IP>

# Node management
syntropy-manager.sh connect $node_name
syntropy-manager.sh status $node_name
```

## Platform Capabilities

This node will support:
- **Container Orchestration**: Docker with Kubernetes-ready configuration
- **Universal Applications**: Any containerized application
- **Scientific Computing**: Python, R, Julia, Fortran, MATLAB
- **Web Applications**: Node.js, Python, Java, Go, PHP, Ruby
- **Machine Learning**: TensorFlow, PyTorch, scikit-learn
- **Databases**: PostgreSQL, MongoDB, Redis, MySQL
- **Monitoring**: Prometheus metrics on port 9100

## Network Configuration

- **SSH Access**: Port 22 (firewall configured)
- **Monitoring**: Port 9100 (Prometheus metrics)
- **DHCP**: Automatic IP assignment
- **Firewall**: UFW enabled with SSH and monitoring access

## Management Tools

### Local Scripts
- **Connection**: \`$HOME/.syntropy/connect-${node_name}.sh\`
- **Discovery**: \`$HOME/.syntropy/discover-${node_name}.sh\`
- **Status**: \`$HOME/.syntropy/status-${node_name}.sh\`

### Central Management
- **Manager**: \`syntropy-manager.sh\`
- **List Nodes**: \`syntropy-manager.sh list\`
- **Node Details**: \`syntropy-manager.sh show $node_name\`

## Multi-Node Management

$(if [ -n "$OWNER_KEY_FILE" ]; then
    echo "This node uses an existing owner key and can be managed alongside other nodes"
    echo "with the same owner key. All nodes sharing this key form a cooperative fleet."
else
    echo "To create additional nodes that can be managed together:"
    echo "\`\`\`bash"
    echo "create-syntropy-usb-enhanced.sh --owner-key $KEYS_DIR/${node_name}_owner.key --node-name second-node"
    echo "\`\`\`"
fi)

## Troubleshooting

### Installation Issues
- Verify hardware meets minimum requirements
- Check BIOS/UEFI boot order and Secure Boot settings
- Ensure USB device is properly connected
- Verify network connectivity (DHCP required)

### Connection Issues
- Run discovery: \`$HOME/.syntropy/discover-${node_name}.sh\`
- Check network connectivity: \`ping <node-ip>\`
- Test SSH port: \`nc -zv <node-ip> 22\`
- Verify firewall allows SSH

### System Logs
```bash
# SSH to node and check logs
journalctl -u syntropy-first-boot
journalctl -u ssh
journalctl -u docker
```

## Support Information

- **Creation System**: $(hostname)
- **Creation User**: $(whoami)
- **Creation Date**: $(date)
- **USB Creator Version**: 2.0.0

---

*This document was automatically generated by the Syntropy USB Creator.*
SUMMARY_EOF

    log SUCCESS "Node summary created: $summary_file"
}

# Create quick reference card
create_quick_reference_card() {
    local node_name="$1"
    local ref_file="$HOME/.syntropy/${node_name}_quick_ref.txt"
    
    cat > "$ref_file" << REF_EOF
SYNTROPY NODE QUICK REFERENCE
=============================
Node: $node_name
Key:  $KEYS_DIR/${node_name}_owner.key

CONNECT:
  $HOME/.syntropy/connect-${node_name}.sh

DISCOVER:
  $HOME/.syntropy/discover-${node_name}.sh

STATUS:
  $HOME/.syntropy/status-${node_name}.sh

MANAGER:
  syntropy-manager.sh list
  syntropy-manager.sh connect $node_name
  syntropy-manager.sh status $node_name

SSH DIRECT:
  ssh -i $KEYS_DIR/${node_name}_owner.key admin@<IP>

TROUBLESHOOT:
  - Check BIOS boot order
  - Verify network (DHCP)
  - Run discovery script
  - Check SSH connectivity
REF_EOF

    log SUCCESS "Quick reference created: $ref_file"
}

# Show final success message
show_final_success_message() {
    local node_name="$1"
    local location_node_id="$2"
    local coordinates="$3"
    local detected_city="$4"
    local detected_country="$5"
    
    echo ""
    echo -e "${GREEN}"
    cat << 'SUCCESS_EOF'
╔════════════════════════════════════════════════════════════════════════════╗
║                          USB CREATION COMPLETE!                           ║
╚════════════════════════════════════════════════════════════════════════════╝
SUCCESS_EOF
    echo -e "${NC}"
    
    echo -e "${PURPLE}═══ NODE SUMMARY ═══${NC}"
    echo "Node Name: $node_name"
    echo "Node ID: $location_node_id"  
    echo "Location: $coordinates ($detected_city, $detected_country)"
    echo "Owner Key: $KEYS_DIR/${node_name}_owner.key"
    echo "Owner Fingerprint: $OWNER_FINGERPRINT"
    echo "Metadata: $NODES_DIR/${node_name}.json"
    echo ""
    
    echo -e "${CYAN}═══ MULTI-NODE CAPABILITY ═══${NC}"
    if [ -n "$OWNER_KEY_FILE" ]; then
        echo "Using existing owner key - this node joins your existing fleet"
        echo "All nodes with this owner key can be managed together"
    else
        echo "New owner key created - save for additional nodes:"
        echo "   Owner key: $KEYS_DIR/${node_name}_owner.key"
        echo "Create additional nodes with: --owner-key $KEYS_DIR/${node_name}_owner.key"
    fi
    echo ""
    
    echo -e "${YELLOW}═══ INSTALLATION PROCESS ═══${NC}"
    echo "1. Insert USB into target hardware"
    echo "2. Configure BIOS/UEFI to boot from USB (F12/F2/DEL)"
    echo "3. Boot and wait for automatic installation (~30 minutes)"
    echo "4. Node will reboot when installation completes"
    echo "5. Use connection tools to access the node"
    echo ""
    
    echo -e "${BLUE}═══ POST-INSTALLATION ACCESS ═══${NC}"
    echo "Quick connect (auto-discovers IP):"
    echo "   $HOME/.syntropy/connect-${node_name}.sh"
    echo ""
    echo "Node management:"
    echo "   syntropy-manager.sh list"
    echo "   syntropy-manager.sh connect $node_name"
    echo "   syntropy-manager.sh status $node_name"
    echo ""
    echo "Discovery and troubleshooting:"
    echo "   $HOME/.syntropy/discover-${node_name}.sh"
    echo "   $HOME/.syntropy/status-${node_name}.sh"
    echo ""
    
    echo -e "${GREEN}═══ UNIVERSAL APPLICATION SUPPORT ═══${NC}"
    echo "This node supports any containerized application including:"
    echo "• Scientific Computing: Python, R, Julia, Fortran, MATLAB"
    echo "• Web Applications: Node.js, Python, Java, Go, PHP, Ruby"
    echo "• Machine Learning: TensorFlow, PyTorch, scikit-learn"
    echo "• Databases: PostgreSQL, MongoDB, Redis, MySQL"
    echo "• Custom Applications: Any Docker container"
    echo ""
    
    echo -e "${PURPLE}═══ DOCUMENTATION ═══${NC}"
    echo "Node summary: $NODES_DIR/${node_name}_summary.md"
    echo "Quick reference: $HOME/.syntropy/${node_name}_quick_ref.txt"
    echo ""
    
    echo -e "${GREEN}Ready for installation! Insert USB and boot target hardware.${NC}"
}