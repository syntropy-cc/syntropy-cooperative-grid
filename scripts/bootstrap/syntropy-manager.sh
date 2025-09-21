#!/bin/bash

# Syntropy Cooperative Grid - Node Management CLI
# Centralized management for all Syntropy nodes

set -e

NODES_DIR="$HOME/.syntropy/nodes"
KEYS_DIR="$HOME/.syntropy/keys"
CONFIG_DIR="$HOME/.syntropy/config"
CACHE_DIR="$HOME/.syntropy/cache"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m'

# Create directories if they don't exist
mkdir -p "$NODES_DIR" "$KEYS_DIR" "$CONFIG_DIR" "$CACHE_DIR"

# Initialize manager config if not exists
MANAGER_CONFIG="$CONFIG_DIR/manager.json"
if [ ! -f "$MANAGER_CONFIG" ]; then
    cat > "$MANAGER_CONFIG" << EOF
{
  "version": "0.1.0",
  "created": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "manager_id": "mgr-$(openssl rand -hex 8)",
  "discovery": {
    "enabled": true,
    "scan_networks": ["192.168.1.0/24", "10.0.0.0/24"],
    "default_ssh_port": 22,
    "connection_timeout": 10
  },
  "security": {
    "key_rotation_days": 90,
    "require_confirmation": true,
    "audit_log": true
  },
  "preferences": {
    "default_editor": "nano",
    "show_coordinates": true,
    "auto_update_metadata": true
  }
}
EOF
fi

show_help() {
    echo -e "${PURPLE}"
    cat << 'EOF'
╔════════════════════════════════════════════════════════════════════════════╗
║                    SYNTROPY COOPERATIVE GRID                              ║
║                        Node Management CLI                                ║
╚════════════════════════════════════════════════════════════════════════════╝
EOF
    echo -e "${NC}"
    echo "Usage: $0 <command> [options]"
    echo ""
    echo -e "${CYAN}Management Commands:${NC}"
    echo "  list                     List all managed nodes"
    echo "  show <node-name>         Show detailed node information"
    echo "  connect <node-name>      Connect to a node via SSH"
    echo "  status <node-name>       Check node status remotely"
    echo "  update <node-name>       Update node metadata"
    echo "  edit <node-name>         Edit node configuration"
    echo ""
    echo -e "${CYAN}Discovery Commands:${NC}"
    echo "  discover                 Discover nodes on local network"
    echo "  scan <network>           Scan specific network for nodes"
    echo "  auto-discover            Auto-discover and register found nodes"
    echo ""
    echo -e "${CYAN}Maintenance Commands:${NC}"
    echo "  backup <node-name>       Backup node configuration"
    echo "  restore <node-name>      Restore node from backup"
    echo "  export <node-name>       Export node configuration"
    echo "  import <file>            Import node configuration"
    echo "  remove <node-name>       Remove node from management"
    echo ""
    echo -e "${CYAN}Security Commands:${NC}"
    echo "  rotate-keys <node-name>  Rotate SSH keys for node"
    echo "  audit                    Show security audit log"
    echo "  verify <node-name>       Verify node integrity"
    echo ""
    echo -e "${CYAN}Batch Operations:${NC}"
    echo "  status-all               Check status of all nodes"
    echo "  update-all               Update metadata for all nodes"
    echo "  health-check             Run health check on all nodes"
    echo ""
    echo -e "${CYAN}Configuration:${NC}"
    echo "  config                   Show manager configuration"
    echo "  set-config <key> <value> Update configuration setting"
    echo ""
    echo -e "${CYAN}Examples:${NC}"
    echo "  $0 list"
    echo "  $0 show syntropy-sp-node-01"
    echo "  $0 connect syntropy-sp-node-01"
    echo "  $0 discover"
    echo "  $0 status-all"
    echo ""
}

log_action() {
    local action="$1"
    local node="$2"
    local details="$3"
    local timestamp="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
    
    if [ "$(jq -r '.security.audit_log' "$MANAGER_CONFIG")" = "true" ]; then
        local log_file="$CONFIG_DIR/audit.log"
        echo "{\"timestamp\":\"$timestamp\",\"action\":\"$action\",\"node\":\"$node\",\"details\":\"$details\"}" >> "$log_file"
    fi
}

require_confirmation() {
    local message="$1"
    if [ "$(jq -r '.security.require_confirmation' "$MANAGER_CONFIG")" = "true" ]; then
        echo -e "${YELLOW}$message${NC}"
        read -p "Continue? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo "Operation cancelled."
            exit 0
        fi
    fi
}

list_nodes() {
    echo -e "${PURPLE}═══ SYNTROPY NODES ═══${NC}"
    
    if [ ! -d "$NODES_DIR" ] || [ -z "$(ls -A "$NODES_DIR" 2>/dev/null)" ]; then
        echo -e "${YELLOW}No nodes found. Create nodes using create-syntropy-usb-enhanced.sh${NC}"
        return
    fi
    
    echo -e "${CYAN}Node Name                    Status      IP Address      Hardware Class     Last Contact${NC}"
    echo "────────────────────────────────────────────────────────────────────────────────────────────"
    
    for node_file in "$NODES_DIR"/*.json; do
        if [ -f "$node_file" ]; then
            local node_name=$(basename "$node_file" .json)
            local node_data=$(cat "$node_file")
            
            local status=$(echo "$node_data" | jq -r '.management.status // "unknown"')
            local ip=$(echo "$node_data" | jq -r '.network.ip_address // "unknown"')
            local hw_class=$(echo "$node_data" | jq -r '.hardware.classification // "unknown"')
            local last_contact=$(echo "$node_data" | jq -r '.management.last_contact // "never"')
            
            # Truncate last_contact for display
            if [ "$last_contact" != "never" ] && [ "$last_contact" != "null" ]; then
                last_contact=$(echo "$last_contact" | cut -d'T' -f1)
            fi
            
            printf "%-28s %-11s %-15s %-18s %s\n" "$node_name" "$status" "$ip" "$hw_class" "$last_contact"
        fi
    done
    
    echo ""
    echo "Total nodes: $(ls -1 "$NODES_DIR"/*.json 2>/dev/null | wc -l)"
}

show_node() {
    local node_name="$1"
    if [ -z "$node_name" ]; then
        echo -e "${RED}Error: Node name required${NC}"
        exit 1
    fi
    
    local node_file="$NODES_DIR/${node_name}.json"
    if [ ! -f "$node_file" ]; then
        echo -e "${RED}Error: Node $node_name not found${NC}"
        exit 1
    fi
    
    local node_data=$(cat "$node_file")
    
    echo -e "${PURPLE}═══ NODE DETAILS: $node_name ═══${NC}"
    echo ""
    
    echo -e "${CYAN}Basic Information:${NC}"
    echo "• Node ID: $(echo "$node_data" | jq -r '.node_info.node_id')"
    echo "• Hostname: $(echo "$node_data" | jq -r '.node_info.hostname // .node_info.node_name')"
    echo "• Description: $(echo "$node_data" | jq -r '.node_info.description')"
    echo "• Created: $(echo "$node_data" | jq -r '.node_info.creation_time // .node_info.installation_time')"
    echo ""
    
    echo -e "${CYAN}Geographic Information:${NC}"
    local coords=$(echo "$node_data" | jq -r '.geographic_info.coordinates')
    echo "• Coordinates: $coords"
    echo "• Detection Method: $(echo "$node_data" | jq -r '.geographic_info.detection_method')"
    echo "• Timezone: $(echo "$node_data" | jq -r '.geographic_info.timezone')"
    echo ""
    
    echo -e "${CYAN}Hardware Specifications:${NC}"
    echo "• CPU Cores: $(echo "$node_data" | jq -r '.hardware.cpu_cores // "unknown"')"
    echo "• RAM: $(echo "$node_data" | jq -r '.hardware.ram_gb // "unknown"') GB"
    echo "• Storage: $(echo "$node_data" | jq -r '.hardware.storage_gb // "unknown"') GB"
    echo "• Architecture: $(echo "$node_data" | jq -r '.hardware.architecture // "unknown"')"
    echo "• Classification: $(echo "$node_data" | jq -r '.hardware.classification // "unknown"')"
    echo ""
    
    echo -e "${CYAN}Network Information:${NC}"
    echo "• IP Address: $(echo "$node_data" | jq -r '.network.ip_address // "unknown"')"
    echo "• Interfaces: $(echo "$node_data" | jq -r '.network.interfaces // "unknown"')"
    echo ""
    
    echo -e "${CYAN}Security & Access:${NC}"
    echo "• SSH Port: $(echo "$node_data" | jq -r '.security.ssh_port // 22')"
    echo "• Owner Key: $KEYS_DIR/${node_name}_owner.key"
    echo "• Community Key: $KEYS_DIR/${node_name}_community.key"
    echo ""
    
    echo -e "${CYAN}Management Status:${NC}"
    echo "• Current Status: $(echo "$node_data" | jq -r '.management.status // "unknown"')"
    echo "• SSH Tested: $(echo "$node_data" | jq -r '.management.ssh_tested // false')"
    echo "• Last Contact: $(echo "$node_data" | jq -r '.management.last_contact // "never"')"
    echo ""
}

connect_node() {
    local node_name="$1"
    if [ -z "$node_name" ]; then
        echo -e "${RED}Error: Node name required${NC}"
        exit 1
    fi
    
    local node_file="$NODES_DIR/${node_name}.json"
    local key_file="$KEYS_DIR/${node_name}_owner.key"
    
    if [ ! -f "$node_file" ]; then
        echo -e "${RED}Error: Node $node_name not found${NC}"
        exit 1
    fi
    
    if [ ! -f "$key_file" ]; then
        echo -e "${RED}Error: SSH key not found for node $node_name${NC}"
        exit 1
    fi
    
    local node_data=$(cat "$node_file")
    local ip=$(echo "$node_data" | jq -r '.network.ip_address')
    local ssh_port=$(echo "$node_data" | jq -r '.security.ssh_port // 22')
    
    if [ "$ip" = "null" ] || [ -z "$ip" ]; then
        echo -e "${YELLOW}No IP address found for node. Attempting discovery...${NC}"
        if discover_node_ip "$node_name"; then
            node_data=$(cat "$node_file")
            ip=$(echo "$node_data" | jq -r '.network.ip_address')
        else
            echo -e "${RED}Could not discover node IP. Try: $0 discover${NC}"
            exit 1
        fi
    fi
    
    echo -e "${BLUE}Connecting to $node_name ($ip)...${NC}"
    
    # Test connection first
    if ssh -i "$key_file" -p "$ssh_port" -o ConnectTimeout=10 -o BatchMode=yes admin@"$ip" exit 2>/dev/null; then
        # Update last contact
        local updated_data=$(echo "$node_data" | jq ".management.last_contact = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\" | .management.ssh_tested = true")
        echo "$updated_data" > "$node_file"
        log_action "connect" "$node_name" "SSH connection successful"
        
        # Connect
        ssh -i "$key_file" -p "$ssh_port" admin@"$ip"
    else
        echo -e "${RED}Connection failed to $node_name ($ip:$ssh_port)${NC}"
        echo "Troubleshooting steps:"
        echo "1. Verify node is powered on and network connected"
        echo "2. Check IP address: $0 discover"
        echo "3. Verify SSH service: nc -zv $ip $ssh_port"
        log_action "connect_failed" "$node_name" "SSH connection failed"
        exit 1
    fi
}

check_node_status() {
    local node_name="$1"
    if [ -z "$node_name" ]; then
        echo -e "${RED}Error: Node name required${NC}"
        exit 1
    fi
    
    local node_file="$NODES_DIR/${node_name}.json"
    local key_file="$KEYS_DIR/${node_name}_owner.key"
    
    if [ ! -f "$node_file" ]; then
        echo -e "${RED}Error: Node $node_name not found${NC}"
        exit 1
    fi
    
    local node_data=$(cat "$node_file")
    local ip=$(echo "$node_data" | jq -r '.network.ip_address')
    local ssh_port=$(echo "$node_data" | jq -r '.security.ssh_port // 22')
    
    if [ "$ip" = "null" ] || [ -z "$ip" ]; then
        echo -e "${YELLOW}No IP address for $node_name - Status: Unknown${NC}"
        return 1
    fi
    
    echo -e "${BLUE}Checking status of $node_name ($ip)...${NC}"
    
    # Check basic connectivity
    if ! nc -z -w5 "$ip" "$ssh_port" 2>/dev/null; then
        echo -e "${RED}✗ $node_name: Network unreachable${NC}"
        return 1
    fi
    
    # Check SSH and get system info
    if ssh -i "$key_file" -p "$ssh_port" -o ConnectTimeout=10 -o BatchMode=yes admin@"$ip" exit 2>/dev/null; then
        echo -e "${GREEN}✓ $node_name: Online and accessible${NC}"
        
        # Get detailed status
        local uptime=$(ssh -i "$key_file" -p "$ssh_port" -o ConnectTimeout=10 admin@"$ip" "uptime -p" 2>/dev/null || echo "unknown")
        local load=$(ssh -i "$key_file" -p "$ssh_port" -o ConnectTimeout=10 admin@"$ip" "uptime | awk -F'load average:' '{print \$2}'" 2>/dev/null || echo "unknown")
        local disk=$(ssh -i "$key_file" -p "$ssh_port" -o ConnectTimeout=10 admin@"$ip" "df / | tail -1 | awk '{print \$5}'" 2>/dev/null || echo "unknown")
        
        echo "  • Uptime: $uptime"
        echo "  • Load Average: $load"
        echo "  • Disk Usage: $disk"
        
        # Update metadata
        local updated_data=$(echo "$node_data" | jq ".management.last_contact = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\" | .management.status = \"online\"")
        echo "$updated_data" > "$node_file"
        
        log_action "status_check" "$node_name" "Node online and responsive"
        return 0
    else
        echo -e "${YELLOW}✗ $node_name: SSH connection failed${NC}"
        local updated_data=$(echo "$node_data" | jq ".management.status = \"ssh_failed\"")
        echo "$updated_data" > "$node_file"
        return 1
    fi
}

discover_node_ip() {
    local node_name="$1"
    local node_file="$NODES_DIR/${node_name}.json"
    local key_file="$KEYS_DIR/${node_name}_owner.key"
    
    if [ ! -f "$node_file" ] || [ ! -f "$key_file" ]; then
        return 1
    fi
    
    echo "Discovering IP for $node_name..."
    
    # Get scan networks from config
    local scan_networks=$(jq -r '.discovery.scan_networks[]' "$MANAGER_CONFIG" 2>/dev/null || echo "192.168.1.0/24")
    
    for network in $scan_networks; do
        echo "Scanning network: $network"
        
        # Use nmap to find live hosts
        local live_hosts=$(nmap -sn "$network" 2>/dev/null | grep "Nmap scan report" | awk '{print $5}' | grep -E '^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$')
        
        for ip in $live_hosts; do
            # Test SSH connection with node key
            if timeout 5 ssh -i "$key_file" -o ConnectTimeout=3 -o BatchMode=yes admin@"$ip" exit 2>/dev/null; then
                echo -e "${GREEN}Found $node_name at $ip${NC}"
                
                # Update node metadata
                local node_data=$(cat "$node_file")
                local updated_data=$(echo "$node_data" | jq ".network.ip_address = \"$ip\" | .management.last_contact = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"")
                echo "$updated_data" > "$node_file"
                
                log_action "discover" "$node_name" "IP discovered: $ip"
                return 0
            fi
        done
    done
    
    return 1
}

discover_network() {
    echo -e "${PURPLE}═══ NETWORK DISCOVERY ═══${NC}"
    
    local scan_networks=$(jq -r '.discovery.scan_networks[]' "$MANAGER_CONFIG" 2>/dev/null || echo "192.168.1.0/24")
    
    for network in $scan_networks; do
        echo -e "${CYAN}Scanning network: $network${NC}"
        
        # Find Syntropy nodes by trying to connect with known keys
        for node_file in "$NODES_DIR"/*.json; do
            if [ -f "$node_file" ]; then
                local node_name=$(basename "$node_file" .json)
                echo -n "  Searching for $node_name... "
                
                if discover_node_ip "$node_name"; then
                    echo -e "${GREEN}Found${NC}"
                else
                    echo -e "${YELLOW}Not found${NC}"
                fi
            fi
        done
    done
}

status_all() {
    echo -e "${PURPLE}═══ ALL NODES STATUS ═══${NC}"
    
    local online=0
    local offline=0
    local total=0
    
    for node_file in "$NODES_DIR"/*.json; do
        if [ -f "$node_file" ]; then
            local node_name=$(basename "$node_file" .json)
            ((total++))
            
            if check_node_status "$node_name"; then
                ((online++))
            else
                ((offline++))
            fi
            echo ""
        fi
    done
    
    echo -e "${CYAN}Summary:${NC}"
    echo "• Total nodes: $total"
    echo -e "• Online: ${GREEN}$online${NC}"
    echo -e "• Offline: ${RED}$offline${NC}"
}

backup_node() {
    local node_name="$1"
    if [ -z "$node_name" ]; then
        echo -e "${RED}Error: Node name required${NC}"
        exit 1
    fi
    
    local node_file="$NODES_DIR/${node_name}.json"
    if [ ! -f "$node_file" ]; then
        echo -e "${RED}Error: Node $node_name not found${NC}"
        exit 1
    fi
    
    local backup_dir="$CONFIG_DIR/backups"
    mkdir -p "$backup_dir"
    
    local timestamp=$(date +%Y%m%d_%H%M%S)
    local backup_file="$backup_dir/${node_name}_${timestamp}.tar.gz"
    
    echo -e "${BLUE}Creating backup for $node_name...${NC}"
    
    # Create temporary directory for backup
    local temp_dir=$(mktemp -d)
    cp "$node_file" "$temp_dir/"
    cp "$KEYS_DIR/${node_name}_"*.key "$temp_dir/" 2>/dev/null || true
    cp "$KEYS_DIR/${node_name}_"*.pub "$temp_dir/" 2>/dev/null || true
    
    # Create tarball
    tar -czf "$backup_file" -C "$temp_dir" .
    rm -rf "$temp_dir"
    
    echo -e "${GREEN}Backup created: $backup_file${NC}"
    log_action "backup" "$node_name" "Backup created: $backup_file"
}

export_node() {
    local node_name="$1"
    if [ -z "$node_name" ]; then
        echo -e "${RED}Error: Node name required${NC}"
        exit 1
    fi
    
    local node_file="$NODES_DIR/${node_name}.json"
    if [ ! -f "$node_file" ]; then
        echo -e "${RED}Error: Node $node_name not found${NC}"
        exit 1
    fi
    
    local export_file="${node_name}_export.json"
    local node_data=$(cat "$node_file")
    
    # Create export with connection instructions
    local export_data=$(jq ". + {
        \"export_info\": {
            \"exported_on\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
            \"exported_by\": \"$USER@$(hostname)\",
            \"connection_command\": \"ssh -i ${node_name}_owner.key admin@$(echo "$node_data" | jq -r '.network.ip_address')\",
            \"required_files\": [\"${node_name}_owner.key\", \"${node_name}_owner.pub\"]
        }
    }" <<< "$node_data")
    
    echo "$export_data" > "$export_file"
    
    # Also export keys
    cp "$KEYS_DIR/${node_name}_owner.key" "${node_name}_owner.key" 2>/dev/null || true
    cp "$KEYS_DIR/${node_name}_owner.pub" "${node_name}_owner.pub" 2>/dev/null || true
    
    echo -e "${GREEN}Node exported to: $export_file${NC}"
    echo -e "${YELLOW}Don't forget to copy the SSH keys: ${node_name}_owner.key${NC}"
    
    log_action "export" "$node_name" "Node configuration exported"
}

remove_node() {
    local node_name="$1"
    if [ -z "$node_name" ]; then
        echo -e "${RED}Error: Node name required${NC}"
        exit 1
    fi
    
    require_confirmation "This will remove node $node_name from management. Keys and backups will be preserved."
    
    local node_file="$NODES_DIR/${node_name}.json"
    if [ ! -f "$node_file" ]; then
        echo -e "${RED}Error: Node $node_name not found${NC}"
        exit 1
    fi
    
    # Create backup before removal
    backup_node "$node_name"
    
    # Move files to removed directory instead of deleting
    local removed_dir="$CONFIG_DIR/removed"
    mkdir -p "$removed_dir"
    
    mv "$node_file" "$removed_dir/"
    echo -e "${GREEN}Node $node_name removed from active management${NC}"
    echo -e "${CYAN}Files preserved in: $removed_dir${NC}"
    
    log_action "remove" "$node_name" "Node removed from management"
}

show_config() {
    echo -e "${PURPLE}═══ MANAGER CONFIGURATION ═══${NC}"
    cat "$MANAGER_CONFIG" | jq '.'
}

set_config() {
    local key="$1"
    local value="$2"
    
    if [ -z "$key" ] || [ -z "$value" ]; then
        echo -e "${RED}Error: Both key and value required${NC}"
        echo "Usage: $0 set-config <key> <value>"
        echo "Example: $0 set-config discovery.enabled true"
        exit 1
    fi
    
    # Update configuration
    local updated_config=$(jq --arg key "$key" --arg value "$value" 'setpath($key | split("."); $value)' "$MANAGER_CONFIG")
    echo "$updated_config" > "$MANAGER_CONFIG"
    
    echo -e "${GREEN}Configuration updated: $key = $value${NC}"
    log_action "config_update" "manager" "$key = $value"
}

# Main command processing
case "${1:-}" in
    list)
        list_nodes
        ;;
    show)
        show_node "$2"
        ;;
    connect)
        connect_node "$2"
        ;;
    status)
        check_node_status "$2"
        ;;
    discover)
        discover_network
        ;;
    status-all)
        status_all
        ;;
    backup)
        backup_node "$2"
        ;;
    export)
        export_node "$2"
        ;;
    remove)
        remove_node "$2"
        ;;
    config)
        show_config
        ;;
    set-config)
        set_config "$2" "$3"
        ;;
    audit)
        if [ -f "$CONFIG_DIR/audit.log" ]; then
            echo -e "${PURPLE}═══ AUDIT LOG ═══${NC}"
            tail -20 "$CONFIG_DIR/audit.log" | jq -r '[.timestamp, .action, .node, .details] | @tsv' | column -t
        else
            echo -e "${YELLOW}No audit log found${NC}"
        fi
        ;;
    help|--help|-h)
        show_help
        ;;
    "")
        show_help
        ;;
    *)
        echo -e "${RED}Unknown command: $1${NC}"
        echo "Use '$0 help' for usage information."
        exit 1
        ;;
esac