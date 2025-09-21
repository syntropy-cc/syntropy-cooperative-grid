#!/bin/bash

# Syntropy Cooperative Grid - Cloud-Init Configuration Generator
# Version: 2.0.0

# Create complete cloud-init configuration
create_cloud_init_configuration() {
    local usb_partition="$1"
    local node_name="$2"
    local location_node_id="$3"
    local coordinates="$4"
    local detection_method="$5"
    local detected_city="$6"
    local detected_country="$7"
    local node_description="$8"
    
    log INFO "Creating cloud-init configuration..."
    
    # Mount USB for writing
    USB_MOUNT="/mnt/syntropy-usb"
    if ! mount_usb_partition "$usb_partition" "$USB_MOUNT"; then
        log ERROR "Failed to mount USB for cloud-init configuration"
        return 1
    fi
    
    # Create user-data file
    create_user_data_file "$USB_MOUNT" "$node_name" "$location_node_id" \
        "$coordinates" "$detection_method" "$detected_city" "$detected_country" "$node_description"
    
    # Create meta-data file
    create_meta_data_file "$USB_MOUNT" "$node_name"
    
    # Add documentation
    add_usb_documentation "$node_name"
    
    # Unmount USB
    unmount_usb_partition "$USB_MOUNT"
    
    log SUCCESS "Cloud-init configuration created successfully"
    return 0
}

# Create the main user-data file
create_user_data_file() {
    local usb_mount="$1"
    local node_name="$2"
    local location_node_id="$3"
    local coordinates="$4"
    local detection_method="$5"
    local detected_city="$6"
    local detected_country="$7"
    local node_description="$8"
    
    log DEBUG "Generating user-data file..."
    
    # Read key contents
    local owner_key_content=$(cat "$OWNER_KEY_PATH")
    local owner_pub_content=$(cat "${OWNER_KEY_PATH}.pub")
    local community_key_content=$(cat "$COMMUNITY_KEY_PATH")
    local community_pub_content=$(cat "${COMMUNITY_KEY_PATH}.pub")
    
    # Create the comprehensive user-data file
    sudo tee "$usb_mount/user-data" > /dev/null << USER_DATA_EOF
#cloud-config
# Syntropy Cooperative Grid - Enhanced Auto-Installation
# Node: $node_name | Location: $coordinates | Created: $(date)
# Platform Version: 2.0.0

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
    - bc

  late-commands:
    # Create comprehensive directory structure
    - curtin in-target -- mkdir -p /opt/syntropy/{identity/{owner,community},platform/{templates,scripts,data},metadata,logs,backups}
    
    # Install security keys and configuration
    - $(generate_key_installation_commands "$owner_key_content" "$owner_pub_content" "$community_key_content" "$community_pub_content")
    
    # Create comprehensive node metadata
    - $(generate_metadata_creation_commands "$node_name" "$location_node_id" "$coordinates" "$detection_method" "$detected_city" "$detected_country" "$node_description")
    
    # Create platform templates
    - $(generate_template_creation_commands)
    
    # Configure services and security
    - curtin in-target -- systemctl enable ssh docker prometheus-node-exporter
    - curtin in-target -- systemctl disable snapd
    - curtin in-target -- ufw default deny incoming
    - curtin in-target -- ufw default allow outgoing  
    - curtin in-target -- ufw allow ssh
    - curtin in-target -- ufw allow 9100/tcp
    - curtin in-target -- ufw --force enable
    - curtin in-target -- usermod -aG docker admin
    
    # Create startup services
    - $(generate_startup_service_commands "$node_name")

  power_state:
    mode: reboot
    timeout: 30
USER_DATA_EOF

    log DEBUG "User-data file created successfully"
}

# Generate key installation commands
generate_key_installation_commands() {
    local owner_key="$1"
    local owner_pub="$2"
    local community_key="$3"
    local community_pub="$4"
    
    cat << 'KEY_INSTALL_EOF'
|
      curtin in-target -- bash -c '
      # Install owner key (SSH access and management)
      cat > /opt/syntropy/identity/owner/private.key << "OWNER_KEY_EOF"
$owner_key
OWNER_KEY_EOF
      
      cat > /opt/syntropy/identity/owner/public.key << "OWNER_PUB_EOF"
$owner_pub
OWNER_PUB_EOF
      
      # Install community key (inter-node communication)
      cat > /opt/syntropy/identity/community/private.key << "COMMUNITY_KEY_EOF"
$community_key
COMMUNITY_KEY_EOF
      
      cat > /opt/syntropy/identity/community/public.key << "COMMUNITY_PUB_EOF"
$community_pub
COMMUNITY_PUB_EOF
      
      # Create key metadata
      cat > /opt/syntropy/identity/key_info.json << "KEY_INFO_EOF"
{
  "owner_key": {
    "fingerprint": "$OWNER_FINGERPRINT",
    "algorithm": "ed25519",
    "purpose": "ssh_access_and_management",
    "created": "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
  },
  "community_key": {
    "fingerprint": "$COMMUNITY_FINGERPRINT", 
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
KEY_INSTALL_EOF
}

# Generate metadata creation commands
generate_metadata_creation_commands() {
    local node_name="$1"
    local location_node_id="$2"
    local coordinates="$3"
    local detection_method="$4"
    local detected_city="$5"
    local detected_country="$6"
    local node_description="$7"
    
    cat << 'METADATA_EOF'
|
      curtin in-target -- bash -c '
      # Detect hardware specifications
      CPU_CORES=$(nproc)
      RAM_GB=$(free -g | awk "/^Mem:/{print \$2}")
      STORAGE_GB=$(df / --output=avail -BG 2>/dev/null | tail -1 | sed "s/G//" | xargs)
      ARCHITECTURE=$(uname -m)
      
      # Enhanced hardware classification
      if [ $RAM_GB -le 2 ]; then
        HW_CLASS="edge"
        K8S_ROLE="worker-light"
        CAPABILITIES="[\"edge_computing\", \"sensor_data\", \"lightweight_services\"]"
      elif [ $RAM_GB -le 8 ]; then
        HW_CLASS="home-server"
        K8S_ROLE="worker"
        CAPABILITIES="[\"container_hosting\", \"development\", \"personal_services\"]"
      elif [ $RAM_GB -le 32 ]; then
        HW_CLASS="server"
        K8S_ROLE="worker-heavy"
        CAPABILITIES="[\"production_workloads\", \"databases\", \"ai_inference\"]"
      else
        HW_CLASS="high-end-server"
        K8S_ROLE="master-capable"
        CAPABILITIES="[\"cluster_management\", \"ai_training\", \"high_performance_computing\"]"
      fi
      
      # Detect network interfaces and current IP
      INTERFACES=$(ip link show | grep "^[0-9]" | awk -F: "{print \$2}" | grep -v lo | tr "\n" "," | sed "s/,$//" | tr " " "")
      CURRENT_IP=$(hostname -I | awk "{print \$1}")
      
      # Create comprehensive installation metadata
      cat > /opt/syntropy/metadata/node.json << "NODE_METADATA_EOF"
{
  "metadata_version": "2.0",
  "node_info": {
    "node_id": "$location_node_id",
    "node_name": "$node_name",
    "hostname": "$(hostname)",
    "description": "$node_description",
    "installation_time": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
    "platform_version": "2.0.0-genesis",
    "platform_type": "syntropy_cooperative_grid"
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
  "hardware": {
    "cpu_cores": $CPU_CORES,
    "ram_gb": $RAM_GB,
    "storage_gb": $STORAGE_GB,
    "architecture": "$ARCHITECTURE",
    "classification": "$HW_CLASS",
    "kubernetes_role": "$K8S_ROLE",
    "capabilities": $CAPABILITIES,
    "detection_time": "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
  },
  "network": {
    "ip_address": "$CURRENT_IP",
    "interfaces": "$INTERFACES",
    "hostname": "$(hostname)",
    "dhcp_configured": true
  },
  "security": {
    "owner_key_fingerprint": "$OWNER_FINGERPRINT",
    "community_key_fingerprint": "$COMMUNITY_FINGERPRINT",
    "ssh_port": 22,
    "authentication_method": "key_only",
    "firewall_enabled": true
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
    "status": "installed",
    "services": {
      "docker": "enabled",
      "ssh": "enabled",
      "prometheus_exporter": "enabled",
      "firewall": "enabled"
    },
    "universal_support": {
      "scientific_computing": ["fortran", "python", "r", "julia", "matlab"],
      "web_applications": ["nodejs", "python", "java", "go", "php", "ruby"],
      "machine_learning": ["tensorflow", "pytorch", "scikit-learn", "keras"],
      "databases": ["postgresql", "mongodb", "redis", "mysql", "cassandra"],
      "custom_applications": "any_containerized_application"
    }
  },
  "management": {
    "status": "installed",
    "installation_complete": true,
    "ssh_ready": true,
    "first_boot": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
    "last_update": "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
  }
}
NODE_METADATA_EOF
      '
METADATA_EOF
}

# Generate template creation commands
generate_template_creation_commands() {
    cat << 'TEMPLATE_EOF'
|
      curtin in-target -- bash -c '
      # Scientific computing template
      cat > /opt/syntropy/platform/templates/batch-job-template.yaml << "BATCH_TEMPLATE_EOF"
apiVersion: batch/v1
kind: Job
metadata:
  name: scientific-computation
  namespace: default
  labels:
    app: scientific-computing
    platform: syntropy
spec:
  template:
    spec:
      containers:
      - name: computation
        image: ubuntu:22.04
        command: ["/bin/bash"]
        args: ["-c", "echo \"Starting computation...\" && sleep 30 && echo \"Computation complete\""]
        resources:
          requests:
            cpu: "100m"
            memory: "128Mi"
          limits:
            cpu: "2000m"
            memory: "2Gi"
      restartPolicy: Never
  backoffLimit: 3
BATCH_TEMPLATE_EOF

      # Web service template
      cat > /opt/syntropy/platform/templates/web-service-template.yaml << "WEB_TEMPLATE_EOF"
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web-application
  namespace: default
  labels:
    app: web-application
    platform: syntropy
spec:
  replicas: 1
  selector:
    matchLabels:
      app: web-application
  template:
    metadata:
      labels:
        app: web-application
    spec:
      containers:
      - name: web
        image: nginx:alpine
        ports:
        - containerPort: 80
        resources:
          requests:
            cpu: "50m"
            memory: "64Mi"
          limits:
            cpu: "500m"
            memory: "512Mi"
---
apiVersion: v1
kind: Service
metadata:
  name: web-service
  labels:
    platform: syntropy
spec:
  selector:
    app: web-application
  ports:
  - port: 80
    targetPort: 80
  type: ClusterIP
WEB_TEMPLATE_EOF

      # Database template
      cat > /opt/syntropy/platform/templates/database-template.yaml << "DB_TEMPLATE_EOF"
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: database-service
  namespace: default
  labels:
    app: database
    platform: syntropy
spec:
  serviceName: database
  replicas: 1
  selector:
    matchLabels:
      app: database
  template:
    metadata:
      labels:
        app: database
    spec:
      containers:
      - name: database
        image: postgres:15-alpine
        env:
        - name: POSTGRES_DB
          value: "syntropy"
        - name: POSTGRES_USER
          value: "admin"
        - name: POSTGRES_PASSWORD
          value: "changeme"
        ports:
        - containerPort: 5432
        volumeMounts:
        - name: data
          mountPath: /var/lib/postgresql/data
        resources:
          requests:
            cpu: "100m"
            memory: "256Mi"
          limits:
            cpu: "1000m"
            memory: "1Gi"
  volumeClaimTemplates:
  - metadata:
      name: data
    spec:
      accessModes: ["ReadWriteOnce"]
      resources:
        requests:
          storage: 10Gi
DB_TEMPLATE_EOF
      '
TEMPLATE_EOF
}

# Generate startup service commands
generate_startup_service_commands() {
    local node_name="$1"
    
    cat << 'STARTUP_EOF'
|
      curtin in-target -- bash -c '
      # Create startup script for first boot
      cat > /opt/syntropy/scripts/first-boot.sh << "FIRST_BOOT_EOF"
#!/bin/bash
# Syntropy first boot setup script
echo "Syntropy node first boot setup starting..." | logger
sleep 30  # Wait for network to be stable
echo "Syntropy node ready for management" | logger
FIRST_BOOT_EOF
      chmod +x /opt/syntropy/scripts/first-boot.sh
      
      # Create systemd service for first boot
      cat > /etc/systemd/system/syntropy-first-boot.service << "SERVICE_EOF"
[Unit]
Description=Syntropy First Boot Setup
After=network-online.target
Wants=network-online.target

[Service]
Type=oneshot
ExecStart=/opt/syntropy/scripts/first-boot.sh
RemainAfterExit=yes

[Install]
WantedBy=multi-user.target
SERVICE_EOF
      
      systemctl enable syntropy-first-boot.service
      '
STARTUP_EOF
}

# Create meta-data file
create_meta_data_file() {
    local usb_mount="$1"
    local node_name="$2"
    
    log DEBUG "Creating meta-data file..."
    
    sudo tee "$usb_mount/meta-data" > /dev/null << META_DATA_EOF
instance-id: syntropy-node-$RANDOM
local-hostname: $node_name
META_DATA_EOF

    log DEBUG "Meta-data file created"
}

# Validate cloud-init configuration
validate_cloud_init_config() {
    local usb_mount="$1"
    
    log INFO "Validating cloud-init configuration..."
    
    # Check if files exist
    if [ ! -f "$usb_mount/user-data" ]; then
        log ERROR "user-data file not found"
        return 1
    fi
    
    if [ ! -f "$usb_mount/meta-data" ]; then
        log ERROR "meta-data file not found"
        return 1
    fi
    
    # Basic syntax validation
    if command -v cloud-init >/dev/null 2>&1; then
        log DEBUG "Validating user-data syntax with cloud-init..."
        if sudo cloud-init devel schema --config-file "$usb_mount/user-data" >/dev/null 2>&1; then
            log SUCCESS "Cloud-init configuration syntax is valid"
        else
            log WARN "Cloud-init syntax validation warnings (may still work)"
        fi
    else
        log DEBUG "cloud-init not available for validation, skipping syntax check"
    fi
    
    log SUCCESS "Cloud-init configuration validated"
    return 0
}

# Complementos para o arquivo cloud-init.sh
# Estas funções devem ser adicionadas ao final do arquivo lib/cloud-init.sh

# Generate enhanced key installation commands with proper escaping
generate_key_installation_commands() {
    local owner_key="$1"
    local owner_pub="$2"
    local community_key="$3"
    local community_pub="$4"
    
    # Escape keys for safe embedding in cloud-init
    local escaped_owner_key=$(echo "$owner_key" | sed 's/$/\\/')
    local escaped_owner_pub=$(echo "$owner_pub" | sed 's/$/\\/')
    local escaped_community_key=$(echo "$community_key" | sed 's/$/\\/')
    local escaped_community_pub=$(echo "$community_pub" | sed 's/$/\\/')
    
    cat << 'KEY_INSTALL_EOF'
curtin in-target -- bash -c "
# Create Syntropy identity structure
mkdir -p /opt/syntropy/identity/{owner,community}
chown -R 1000:1000 /opt/syntropy

# Install owner key (SSH access and management)
cat > /opt/syntropy/identity/owner/private.key << 'OWNER_KEY_END'
KEY_INSTALL_EOF
    
    echo "$owner_key"
    
    cat << 'KEY_INSTALL_MIDDLE'
OWNER_KEY_END

cat > /opt/syntropy/identity/owner/public.key << 'OWNER_PUB_END'
KEY_INSTALL_MIDDLE
    
    echo "$owner_pub"
    
    cat << 'KEY_INSTALL_CONTINUE'
OWNER_PUB_END

# Install community key (inter-node communication)
cat > /opt/syntropy/identity/community/private.key << 'COMMUNITY_KEY_END'
KEY_INSTALL_CONTINUE
    
    echo "$community_key"
    
    cat << 'KEY_INSTALL_FINAL'
COMMUNITY_KEY_END

cat > /opt/syntropy/identity/community/public.key << 'COMMUNITY_PUB_END'
KEY_INSTALL_FINAL
    
    echo "$community_pub"
    
    cat << 'KEY_INSTALL_FINISH'
COMMUNITY_PUB_END

# Set proper permissions
chmod 600 /opt/syntropy/identity/owner/private.key
chmod 600 /opt/syntropy/identity/community/private.key
chmod 644 /opt/syntropy/identity/owner/public.key
chmod 644 /opt/syntropy/identity/community/public.key

# Configure SSH access
mkdir -p /home/admin/.ssh
cp /opt/syntropy/identity/owner/public.key /home/admin/.ssh/authorized_keys
chmod 600 /home/admin/.ssh/authorized_keys
chown admin:admin /home/admin/.ssh/authorized_keys
chown -R admin:admin /opt/syntropy/

# Create key metadata
cat > /opt/syntropy/identity/key_info.json << 'KEY_INFO_END'
{
  \"owner_key\": {
    \"fingerprint\": \"$OWNER_FINGERPRINT\",
    \"algorithm\": \"ed25519\",
    \"purpose\": \"ssh_access_and_management\",
    \"created\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"
  },
  \"community_key\": {
    \"fingerprint\": \"$COMMUNITY_FINGERPRINT\", 
    \"algorithm\": \"ed25519\",
    \"purpose\": \"inter_node_communication\",
    \"created\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"
  }
}
KEY_INFO_END

chmod 644 /opt/syntropy/identity/key_info.json
"
KEY_INSTALL_FINISH
}

# Generate startup service commands
generate_startup_service_commands() {
    local node_name="$1"
    
    cat << 'STARTUP_COMMANDS_EOF'
curtin in-target -- bash -c "
# Create startup script for first boot
cat > /opt/syntropy/scripts/first-boot.sh << 'FIRST_BOOT_END'
#!/bin/bash
# Syntropy first boot setup script

# Wait for network to be stable
sleep 30

# Log startup
echo \"Syntropy node $node_name first boot setup starting...\" | logger -t syntropy

# Detect hardware and update metadata
CPU_CORES=\$(nproc)
RAM_GB=\$(free -g | awk \"/^Mem:/{print \\\$2}\")
STORAGE_GB=\$(df / --output=avail -BG 2>/dev/null | tail -1 | sed \"s/G//\" | xargs)
ARCHITECTURE=\$(uname -m)

# Get current IP
CURRENT_IP=\$(hostname -I | awk \"{print \\\$1}\")

# Update node metadata with runtime information
if [ -f /opt/syntropy/metadata/node.json ]; then
    # Create updated metadata with runtime info
    python3 -c \"
import json
import sys
from datetime import datetime

try:
    with open('/opt/syntropy/metadata/node.json', 'r') as f:
        data = json.load(f)
    
    # Update hardware info with actual detected values
    data['hardware']['cpu_cores'] = $CPU_CORES
    data['hardware']['ram_gb'] = $RAM_GB
    data['hardware']['storage_gb'] = $STORAGE_GB
    data['hardware']['architecture'] = '$ARCHITECTURE'
    
    # Update network info
    data['network']['ip_address'] = '$CURRENT_IP'
    
    # Update management status
    data['management']['status'] = 'installed'
    data['management']['installation_complete'] = True
    data['management']['first_boot'] = datetime.utcnow().isoformat() + 'Z'
    data['management']['last_update'] = datetime.utcnow().isoformat() + 'Z'
    
    with open('/opt/syntropy/metadata/node.json', 'w') as f:
        json.dump(data, f, indent=2)
    
    print('Metadata updated successfully')
except Exception as e:
    print(f'Error updating metadata: {e}', file=sys.stderr)
\"
fi

# Ensure all services are running
systemctl enable docker
systemctl start docker

systemctl enable ssh
systemctl start ssh

systemctl enable prometheus-node-exporter
systemctl start prometheus-node-exporter

# Configure firewall
ufw --force reset
ufw default deny incoming
ufw default allow outgoing
ufw allow ssh
ufw allow 9100/tcp
ufw --force enable

# Configure fail2ban
systemctl enable fail2ban
systemctl start fail2ban

echo \"Syntropy node $node_name ready for management\" | logger -t syntropy

# Create ready indicator
touch /opt/syntropy/.ready

exit 0
FIRST_BOOT_END

chmod +x /opt/syntropy/scripts/first-boot.sh

# Create systemd service for first boot
cat > /etc/systemd/system/syntropy-first-boot.service << 'SERVICE_END'
[Unit]
Description=Syntropy First Boot Setup
After=network-online.target cloud-init.service
Wants=network-online.target

[Service]
Type=oneshot
ExecStart=/opt/syntropy/scripts/first-boot.sh
RemainAfterExit=yes
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
SERVICE_END

systemctl enable syntropy-first-boot.service
"
STARTUP_COMMANDS_EOF
}

# Enhanced metadata creation with runtime detection capabilities
generate_metadata_creation_commands() {
    local node_name="$1"
    local location_node_id="$2"
    local coordinates="$3"
    local detection_method="$4"
    local detected_city="$5"
    local detected_country="$6"
    local node_description="$7"
    
    cat << 'METADATA_COMMANDS_EOF'
curtin in-target -- bash -c "
# Create comprehensive installation metadata
cat > /opt/syntropy/metadata/node.json << 'NODE_METADATA_END'
{
  \"metadata_version\": \"2.0\",
  \"node_info\": {
    \"node_id\": \"$location_node_id\",
    \"node_name\": \"$node_name\",
    \"hostname\": \"$(hostname)\",
    \"description\": \"$node_description\",
    \"installation_time\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
    \"platform_version\": \"2.0.0-genesis\",
    \"platform_type\": \"syntropy_cooperative_grid\"
  },
  \"geographic_info\": {
    \"coordinates\": {
      \"latitude\": LATITUDE_PLACEHOLDER,
      \"longitude\": LONGITUDE_PLACEHOLDER,
      \"formatted\": \"$coordinates\"
    },
    \"location\": {
      \"city\": \"$detected_city\",
      \"country\": \"$detected_country\",
      \"timezone\": \"$(timedatectl show --property=Timezone --value 2>/dev/null || echo UTC)\"
    },
    \"detection\": {
      \"method\": \"$detection_method\",
      \"timestamp\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
      \"accuracy\": \"ACCURACY_PLACEHOLDER\"
    },
    \"location_id\": \"$location_node_id\"
  },
  \"hardware\": {
    \"cpu_cores\": \"RUNTIME_DETECTION\",
    \"ram_gb\": \"RUNTIME_DETECTION\",
    \"storage_gb\": \"RUNTIME_DETECTION\",
    \"architecture\": \"RUNTIME_DETECTION\",
    \"classification\": \"RUNTIME_CLASSIFICATION\",
    \"kubernetes_role\": \"RUNTIME_K8S_ROLE\",
    \"capabilities\": \"RUNTIME_CAPABILITIES\",
    \"detection_time\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"
  },
  \"network\": {
    \"ip_address\": \"RUNTIME_DETECTION\",
    \"interfaces\": \"RUNTIME_DETECTION\",
    \"hostname\": \"$(hostname)\",
    \"dhcp_configured\": true
  },
  \"security\": {
    \"owner_key_fingerprint\": \"$OWNER_FINGERPRINT\",
    \"community_key_fingerprint\": \"$COMMUNITY_FINGERPRINT\",
    \"ssh_port\": 22,
    \"authentication_method\": \"key_only\",
    \"firewall_enabled\": true
  },
  \"platform\": {
    \"type\": \"syntropy_cooperative_grid\",
    \"capabilities\": [
      \"container_orchestration\",
      \"resource_sharing\",
      \"cooperative_computing\",
      \"distributed_storage\",
      \"service_mesh\",
      \"universal_applications\"
    ],
    \"status\": \"installing\",
    \"services\": {
      \"docker\": \"enabled\",
      \"ssh\": \"enabled\",
      \"prometheus_exporter\": \"enabled\",
      \"firewall\": \"enabled\"
    },
    \"universal_support\": {
      \"scientific_computing\": [\"fortran\", \"python\", \"r\", \"julia\", \"matlab\"],
      \"web_applications\": [\"nodejs\", \"python\", \"java\", \"go\", \"php\", \"ruby\"],
      \"machine_learning\": [\"tensorflow\", \"pytorch\", \"scikit-learn\", \"keras\"],
      \"databases\": [\"postgresql\", \"mongodb\", \"redis\", \"mysql\", \"cassandra\"],
      \"custom_applications\": \"any_containerized_application\"
    }
  },
  \"management\": {
    \"status\": \"installing\",
    \"installation_complete\": false,
    \"ssh_ready\": false,
    \"first_boot\": null,
    \"last_update\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"
  }
}
NODE_METADATA_END

# Replace placeholders with actual coordinate values
python3 -c \"
import json
import sys

coords = '$coordinates'.split(',')
lat = float(coords[0])
lon = float(coords[1])

detection_method = '$detection_method'
if 'ip_geolocation' in detection_method:
    accuracy = 'high'
elif 'timezone' in detection_method:
    accuracy = 'medium'
elif 'manual' in detection_method:
    accuracy = 'exact'
else:
    accuracy = 'low'

try:
    with open('/opt/syntropy/metadata/node.json', 'r') as f:
        content = f.read()
    
    content = content.replace('LATITUDE_PLACEHOLDER', str(lat))
    content = content.replace('LONGITUDE_PLACEHOLDER', str(lon))
    content = content.replace('ACCURACY_PLACEHOLDER', accuracy)
    
    with open('/opt/syntropy/metadata/node.json', 'w') as f:
        f.write(content)
        
    print('Coordinates updated in metadata')
except Exception as e:
    print(f'Error processing metadata: {e}', file=sys.stderr)
\"

chmod 644 /opt/syntropy/metadata/node.json
chown admin:admin /opt/syntropy/metadata/node.json
"
METADATA_COMMANDS_EOF
}