#!/bin/bash

# Syntropy Cooperative Grid - Management Environment Setup
# Initializes the complete node management system

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m'

SYNTROPY_DIR="$HOME/.syntropy"
NODES_DIR="$SYNTROPY_DIR/nodes"
KEYS_DIR="$SYNTROPY_DIR/keys"
CONFIG_DIR="$SYNTROPY_DIR/config"
CACHE_DIR="$SYNTROPY_DIR/cache"
SCRIPTS_DIR="$SYNTROPY_DIR/scripts"
BACKUPS_DIR="$SYNTROPY_DIR/backups"

echo -e "${PURPLE}"
cat << 'EOF'
╔════════════════════════════════════════════════════════════════════════════╗
║                    SYNTROPY COOPERATIVE GRID                              ║
║                   Management Environment Setup                            ║
║                                                                            ║
║  Setting up complete node management infrastructure                       ║
╚════════════════════════════════════════════════════════════════════════════╝
EOF
echo -e "${NC}"

# Check if already initialized
if [ -f "$CONFIG_DIR/manager.json" ]; then
    echo -e "${YELLOW}Syntropy management environment already exists.${NC}"
    echo "Manager config: $CONFIG_DIR/manager.json"
    echo "Nodes directory: $NODES_DIR"
    echo ""
    read -p "Reinitialize? This will preserve existing nodes but reset configuration (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Setup cancelled. Use 'syntropy-manager.sh' to manage existing nodes."
        exit 0
    fi
fi

echo -e "${BLUE}[1/6] Creating directory structure...${NC}"

# Create directory structure
mkdir -p "$NODES_DIR" "$KEYS_DIR" "$CONFIG_DIR" "$CACHE_DIR" "$SCRIPTS_DIR" "$BACKUPS_DIR"
mkdir -p "$CONFIG_DIR/removed" "$CONFIG_DIR/templates" "$CONFIG_DIR/logs"

echo "Directory structure created:"
echo "├── $SYNTROPY_DIR"
echo "│   ├── nodes/           # Node metadata and configurations"
echo "│   ├── keys/            # SSH keys for all managed nodes"
echo "│   ├── config/          # Manager configuration and templates"
echo "│   ├── cache/           # Temporary files and discovery cache"
echo "│   ├── scripts/         # Custom scripts and tools"
echo "│   └── backups/         # Node configuration backups"

echo -e "${BLUE}[2/6] Installing dependencies...${NC}"

# Check for required tools
MISSING_TOOLS=()

command -v jq >/dev/null 2>&1 || MISSING_TOOLS+=("jq")
command -v nmap >/dev/null 2>&1 || MISSING_TOOLS+=("nmap")
command -v python3 >/dev/null 2>&1 || MISSING_TOOLS+=("python3")
command -v ssh-keygen >/dev/null 2>&1 || MISSING_TOOLS+=("openssh-client")
command -v curl >/dev/null 2>&1 || MISSING_TOOLS+=("curl")

if [ ${#MISSING_TOOLS[@]} -gt 0 ]; then
    echo "Installing missing dependencies: ${MISSING_TOOLS[*]}"
    
    # Detect package manager and install
    if command -v apt-get >/dev/null 2>&1; then
        sudo apt-get update && sudo apt-get install -y "${MISSING_TOOLS[@]}"
    elif command -v yum >/dev/null 2>&1; then
        sudo yum install -y "${MISSING_TOOLS[@]}"
    elif command -v brew >/dev/null 2>&1; then
        brew install "${MISSING_TOOLS[@]}"
    else
        echo -e "${YELLOW}Please install manually: ${MISSING_TOOLS[*]}${NC}"
    fi
fi

echo -e "${BLUE}[3/6] Creating manager configuration...${NC}"

# Generate unique manager ID
MANAGER_ID="mgr-$(openssl rand -hex 8)"

# Detect local networks for discovery
LOCAL_NETWORKS=$(ip route | grep -E "192\.168\.|10\.|172\." | grep "/" | awk '{print $1}' | sort -u | head -5 | tr '\n' ',' | sed 's/,$//')

# Create comprehensive manager configuration
cat > "$CONFIG_DIR/manager.json" << CONFIG_EOF
{
  "version": "1.0.0",
  "created": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "manager_id": "$MANAGER_ID",
  "system_info": {
    "hostname": "$(hostname)",
    "user": "$USER",
    "os": "$(uname -s)",
    "architecture": "$(uname -m)"
  },
  "discovery": {
    "enabled": true,
    "scan_networks": [$(echo "$LOCAL_NETWORKS" | sed 's/,/","/g' | sed 's/^/"/' | sed 's/$/"/'),
                      "192.168.1.0/24", "10.0.0.0/24", "172.16.0.0/16"],
    "default_ssh_port": 22,
    "connection_timeout": 10,
    "parallel_scans": 5,
    "cache_results": true,
    "cache_ttl_minutes": 30
  },
  "security": {
    "key_rotation_days": 90,
    "require_confirmation": true,
    "audit_log": true,
    "backup_keys": true,
    "verify_fingerprints": true
  },
  "preferences": {
    "default_editor": "nano",
    "show_coordinates": true,
    "auto_update_metadata": true,
    "concurrent_connections": 3,
    "verbose_output": false,
    "color_output": true
  },
  "notifications": {
    "enabled": false,
    "email": "",
    "webhook_url": "",
    "slack_channel": ""
  },
  "backup": {
    "auto_backup": true,
    "backup_frequency_days": 7,
    "max_backups": 30,
    "compress_backups": true
  }
}
CONFIG_EOF

echo "Manager configuration created with ID: $MANAGER_ID"

echo -e "${BLUE}[4/6] Creating helper scripts and tools...${NC}"

# Create network discovery script
cat > "$SCRIPTS_DIR/discover-network.sh" << 'DISCOVER_EOF'
#!/bin/bash

# Quick network discovery for Syntropy nodes
echo "Discovering devices on local networks..."

# Get network interfaces and their networks
NETWORKS=$(ip route | grep -E "192\.168\.|10\.|172\." | grep "/" | awk '{print $1}' | sort -u)

for network in $NETWORKS; do
    echo "Scanning $network..."
    nmap -sn "$network" 2>/dev/null | grep "Nmap scan report" | awk '{print $5}' | while read ip; do
        if [[ $ip =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
            # Check if SSH is open
            if nc -z -w2 "$ip" 22 2>/dev/null; then
                echo "  $ip - SSH available"
            fi
        fi
    done
done
DISCOVER_EOF

chmod +x "$SCRIPTS_DIR/discover-network.sh"

# Create backup script
cat > "$SCRIPTS_DIR/backup-all-nodes.sh" << 'BACKUP_EOF'
#!/bin/bash

# Backup all node configurations
BACKUP_DIR="$HOME/.syntropy/backups/$(date +%Y%m%d_%H%M%S)"
mkdir -p "$BACKUP_DIR"

echo "Creating backup in: $BACKUP_DIR"

# Copy all node metadata
cp -r "$HOME/.syntropy/nodes" "$BACKUP_DIR/"
cp -r "$HOME/.syntropy/keys" "$BACKUP_DIR/"
cp -r "$HOME/.syntropy/config" "$BACKUP_DIR/"

# Create backup info
cat > "$BACKUP_DIR/backup_info.json" << EOF
{
  "created": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "system": "$(hostname)",
  "user": "$USER",
  "nodes_count": $(ls -1 "$HOME/.syntropy/nodes"/*.json 2>/dev/null | wc -l),
  "backup_type": "full"
}
EOF

# Compress backup
tar -czf "$BACKUP_DIR.tar.gz" -C "$(dirname "$BACKUP_DIR")" "$(basename "$BACKUP_DIR")"
rm -rf "$BACKUP_DIR"

echo "Backup created: $BACKUP_DIR.tar.gz"
BACKUP_EOF

chmod +x "$SCRIPTS_DIR/backup-all-nodes.sh"

# Create status check script
cat > "$SCRIPTS_DIR/health-check-all.sh" << 'HEALTH_EOF'
#!/bin/bash

# Check health of all managed nodes
TOTAL=0
ONLINE=0
OFFLINE=0

echo "=== SYNTROPY NETWORK HEALTH CHECK ==="
echo ""

for node_file in "$HOME/.syntropy/nodes"/*.json; do
    if [ -f "$node_file" ]; then
        ((TOTAL++))
        
        node_name=$(basename "$node_file" .json)
        ip=$(jq -r '.network.ip_address // "unknown"' "$node_file")
        
        if [ "$ip" != "unknown" ] && [ "$ip" != "null" ]; then
            if nc -z -w3 "$ip" 22 2>/dev/null; then
                echo "✓ $node_name ($ip) - Online"
                ((ONLINE++))
            else
                echo "✗ $node_name ($ip) - Offline"
                ((OFFLINE++))
            fi
        else
            echo "? $node_name - No IP address"
            ((OFFLINE++))
        fi
    fi
done

echo ""
echo "Summary: $ONLINE/$TOTAL nodes online"
HEALTH_EOF

chmod +x "$SCRIPTS_DIR/health-check-all.sh"

echo -e "${BLUE}[5/6] Creating application templates...${NC}"

# Create application deployment templates
mkdir -p "$CONFIG_DIR/templates/applications"

# Fortran scientific computing template
cat > "$CONFIG_DIR/templates/applications/fortran-computation.yaml" << 'FORTRAN_EOF'
apiVersion: batch/v1
kind: Job
metadata:
  name: fortran-simulation
  labels:
    app: scientific-computing
    language: fortran
spec:
  template:
    spec:
      containers:
      - name: fortran-runner
        image: gcc:latest
        workingDir: /workspace
        command: ["/bin/bash"]
        args:
          - -c
          - |
            # Install Fortran compiler
            apt-get update && apt-get install -y gfortran
            
            # Compile and run Fortran program
            echo "program hello
              print *, 'Hello from Syntropy Cooperative Grid!'
              print *, 'Running Fortran simulation...'
            end program hello" > simulation.f90
            
            gfortran -o simulation simulation.f90
            ./simulation
            
            echo "Computation complete"
        resources:
          requests:
            cpu: "500m"
            memory: "512Mi"
          limits:
            cpu: "2000m"
            memory: "2Gi"
        volumeMounts:
        - name: workspace
          mountPath: /workspace
      volumes:
      - name: workspace
        emptyDir: {}
      restartPolicy: Never
  backoffLimit: 3
FORTRAN_EOF

# Python data science template
cat > "$CONFIG_DIR/templates/applications/python-datascience.yaml" << 'PYTHON_EOF'
apiVersion: apps/v1
kind: Deployment
metadata:
  name: jupyter-lab
  labels:
    app: data-science
    language: python
spec:
  replicas: 1
  selector:
    matchLabels:
      app: jupyter-lab
  template:
    metadata:
      labels:
        app: jupyter-lab
    spec:
      containers:
      - name: jupyter
        image: jupyter/datascience-notebook:latest
        ports:
        - containerPort: 8888
        env:
        - name: JUPYTER_ENABLE_LAB
          value: "yes"
        - name: JUPYTER_TOKEN
          value: "syntropy123"
        resources:
          requests:
            cpu: "200m"
            memory: "512Mi"
          limits:
            cpu: "2000m"
            memory: "4Gi"
        volumeMounts:
        - name: notebooks
          mountPath: /home/jovyan/work
      volumes:
      - name: notebooks
        emptyDir: {}
---
apiVersion: v1
kind: Service
metadata:
  name: jupyter-service
spec:
  selector:
    app: jupyter-lab
  ports:
  - port: 8888
    targetPort: 8888
  type: NodePort
PYTHON_EOF

echo -e "${BLUE}[6/6] Setting up command aliases and PATH...${NC}"

# Add syntropy commands to PATH
SYNTROPY_BASHRC="$CONFIG_DIR/syntropy.bashrc"
cat > "$SYNTROPY_BASHRC" << 'BASHRC_EOF'
# Syntropy Cooperative Grid - Management Aliases and Functions

# Core commands
alias syntropy-list='syntropy-manager.sh list'
alias syntropy-status='syntropy-manager.sh status-all'
alias syntropy-discover='syntropy-manager.sh discover'
alias syntropy-health='~/.syntropy/scripts/health-check-all.sh'
alias syntropy-backup='~/.syntropy/scripts/backup-all-nodes.sh'

# Quick node access
syntropy-connect() {
    if [ -z "$1" ]; then
        echo "Usage: syntropy-connect <node-name>"
        syntropy-manager.sh list
        return 1
    fi
    syntropy-manager.sh connect "$1"
}

syntropy-ssh() {
    if [ -z "$1" ]; then
        echo "Usage: syntropy-ssh <node-name>"
        return 1
    fi
    
    local node_file="$HOME/.syntropy/nodes/${1}.json"
    local key_file="$HOME/.syntropy/keys/${1}_owner.key"
    
    if [ ! -f "$node_file" ]; then
        echo "Node $1 not found"
        return 1
    fi
    
    local ip=$(jq -r '.network.ip_address' "$node_file" 2>/dev/null)
    if [ "$ip" = "null" ] || [ -z "$ip" ]; then
        echo "No IP address for node $1. Try: syntropy-discover"
        return 1
    fi
    
    ssh -i "$key_file" admin@"$ip"
}

# Node creation shortcut
syntropy-create() {
    local usb_device="$1"
    shift
    create-syntropy-usb-enhanced.sh "$usb_device" "$@"
}

# Show syntropy status
syntropy-info() {
    echo "=== SYNTROPY MANAGEMENT STATUS ==="
    echo "Managed nodes: $(ls -1 "$HOME/.syntropy/nodes"/*.json 2>/dev/null | wc -l)"
    echo "SSH keys: $(ls -1 "$HOME/.syntropy/keys"/*_owner.key 2>/dev/null | wc -l)"
    echo "Config dir: $HOME/.syntropy"
    echo ""
    echo "Recent nodes:"
    ls -t "$HOME/.syntropy/nodes"/*.json 2>/dev/null | head -3 | while read f; do
        local name=$(basename "$f" .json)
        local ip=$(jq -r '.network.ip_address // "unknown"' "$f")
        echo "  $name ($ip)"
    done
}

export PATH="$HOME/.syntropy/scripts:$PATH"
BASHRC_EOF

# Add to user's bashrc if not already there
if ! grep -q "syntropy.bashrc" "$HOME/.bashrc" 2>/dev/null; then
    echo "" >> "$HOME/.bashrc"
    echo "# Syntropy Cooperative Grid Management" >> "$HOME/.bashrc"
    echo "source $SYNTROPY_BASHRC" >> "$HOME/.bashrc"
    echo "Syntropy aliases added to ~/.bashrc"
else
    echo "Syntropy aliases already configured in ~/.bashrc"
fi

echo -e "${GREEN}"
cat << 'COMPLETE_EOF'
╔════════════════════════════════════════════════════════════════════════════╗
║                    SETUP COMPLETE - READY TO USE!                         ║
╚════════════════════════════════════════════════════════════════════════════╝
COMPLETE_EOF
echo -e "${NC}"

echo -e "${PURPLE}═══ SYNTROPY MANAGEMENT ENVIRONMENT ═══${NC}"
echo "🏠 Base directory: $SYNTROPY_DIR"
echo "🔧 Manager ID: $MANAGER_ID"
echo "📋 Configuration: $CONFIG_DIR/manager.json"
echo ""

echo -e "${CYAN}═══ AVAILABLE COMMANDS ═══${NC}"
echo "Core Management:"
echo "  syntropy-manager.sh list              # List all nodes"
echo "  syntropy-manager.sh connect <node>    # Connect to specific node"
echo "  syntropy-manager.sh status <node>     # Check node status"
echo "  syntropy-manager.sh discover          # Discover nodes on network"
echo ""
echo "Quick Aliases (after source ~/.bashrc):"
echo "  syntropy-list                         # List all nodes"
echo "  syntropy-connect <node>               # Connect to node"
echo "  syntropy-ssh <node>                   # Direct SSH to node"
echo "  syntropy-status                       # Status of all nodes"
echo "  syntropy-discover                     # Discover network nodes"
echo "  syntropy-health                       # Health check all nodes"
echo "  syntropy-backup                       # Backup all configurations"
echo "  syntropy-info                         # Show management summary"
echo ""
echo "Node Creation:"
echo "  create-syntropy-usb-enhanced.sh /dev/sdb              # Create first node"
echo "  syntropy-create /dev/sdb --node-name my-server       # Quick create"
echo ""

echo -e "${YELLOW}═══ NEXT STEPS ═══${NC}"
echo "1. Reload shell: source ~/.bashrc"
echo "2. Create your first node:"
echo "   create-syntropy-usb-enhanced.sh /dev/sdb --node-name main-server"
echo "3. Boot USB on target hardware"
echo "4. Wait for installation (~30 minutes)"
echo "5. Connect to node: syntropy-connect main-server"
echo ""

echo -e "${BLUE}═══ HELPER SCRIPTS CREATED ═══${NC}"
echo "🔍 Network discovery: $SCRIPTS_DIR/discover-network.sh"
echo "💾 Backup utility: $SCRIPTS_DIR/backup-all-nodes.sh"
echo "❤️  Health checker: $SCRIPTS_DIR/health-check-all.sh"
echo ""

echo -e "${GREEN}Management environment ready! Start creating nodes with create-syntropy-usb-enhanced.sh${NC}"

# Source the aliases immediately if running interactively
if [[ $- == *i* ]]; then
    source "$SYNTROPY_BASHRC"
    echo ""
    echo -e "${GREEN}Aliases loaded! Try: syntropy-info${NC}"
fi