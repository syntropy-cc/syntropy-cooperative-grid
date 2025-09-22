#!/bin/bash

# Syntropy Cooperative Grid - Enhanced Cloud-Init Configuration
# Version: 2.1.0

# Create cloud-init configuration with proper Ubuntu 22.04 support
create_cloud_init_configuration_fixed() {
    local usb_partition="$1"
    local node_name="$2"
    local location_node_id="$3"
    local coordinates="$4"
    local detection_method="$5"
    local detected_city="$6"
    local detected_country="$7"
    local node_description="$8"
    
    log INFO "Creating enhanced cloud-init configuration..."
    
    # Mount points
    USB_MOUNT="/mnt/syntropy-usb"
    mkdir -p "$USB_MOUNT"
    
    # Mount USB
    if ! sudo mount "$usb_partition" "$USB_MOUNT"; then
        log ERROR "Failed to mount USB for cloud-init configuration"
        return 1
    fi
    
    # Create autoinstall directory
    sudo mkdir -p "$USB_MOUNT/autoinstall"
    
    # Generate configurations
    create_user_data_fixed "$USB_MOUNT" "$node_name" "$location_node_id" \
        "$coordinates" "$detection_method" "$detected_city" "$detected_country" "$node_description"
    
    create_meta_data_fixed "$USB_MOUNT" "$node_name"
    create_vendor_data_fixed "$USB_MOUNT"
    
    # Create documentation
    create_readme_fixed "$USB_MOUNT" "$node_name"
    
    # Unmount
    sudo umount "$USB_MOUNT"
    
    log SUCCESS "Cloud-init configuration created successfully"
    return 0
}

# Create enhanced user-data configuration
create_user_data_fixed() {
    local usb_mount="$1"
    local node_name="$2"
    local node_id="$3"
    local coordinates="$4"
    local detection_method="$5"
    local city="$6"
    local country="$7"
    local description="$8"
    
    log DEBUG "Generating user-data configuration..."
    
    # Generate password hash for initial password
    local password_hash=$(mkpasswd -m sha-512 "syntropy" 2>/dev/null || echo '$6$rounds=4096$syntropy$N8mVzFK0Y1OelT1SKEjg0jIXzKMzL3ZcOGcE5xR8nS6E8qSO5qFV6eJs1g7T6E0cC7w.kfNO3FqC3YhE9Gz19.')
    
    # Read and escape key contents
    local owner_key=$(cat "$OWNER_KEY_PATH" | sed 's/$/\\n/g' | tr -d '\n')
    local owner_pub=$(cat "${OWNER_KEY_PATH}.pub")
    local community_key=$(cat "$COMMUNITY_KEY_PATH" | sed 's/$/\\n/g' | tr -d '\n')
    local community_pub=$(cat "${COMMUNITY_KEY_PATH}.pub")
    
    # Create user-data
    cat > "$usb_mount/autoinstall/user-data" << EOF
#cloud-config
autoinstall:
  version: 1
  
  # Locale and keyboard
  locale: en_US.UTF-8
  keyboard:
    layout: us
  
  # Network configuration - Use DHCP initially
  network:
    network:
      version: 2
      ethernets:
        all-eth:
          match:
            name: eth*
          dhcp4: true
  
  # Storage configuration
  storage:
    layout:
      name: direct
    config:
      - type: disk
        id: disk0
        match:
          size: largest
      - type: partition
        id: boot-partition
        device: disk0
        size: 512M
        flag: boot
        grub_device: true
      - type: partition
        id: root-partition
        device: disk0
        size: -1
  
  # User configuration
  identity:
    hostname: ${node_name}
    username: admin
    password: ${password_hash}
  
  # SSH configuration
  ssh:
    install-server: true
    allow-pw: false
    authorized-keys:
      - ${owner_pub}
      - ${community_pub}
  
  # Package configuration
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
    - net-tools
    - bc
    - tmux
    - unzip
    - ca-certificates
    - gnupg
    - lsb-release
  
  # Late commands for final setup
  late-commands:
    # Create Syntropy directories
    - curtin in-target -- mkdir -p /opt/syntropy/{identity,platform,metadata,logs,config,scripts}
    - curtin in-target -- mkdir -p /opt/syntropy/identity/{owner,community}
    
    # Install identity keys
    - |
      echo "${owner_key}" | curtin in-target -- tee /opt/syntropy/identity/owner/private.key
    - |
      echo "${owner_pub}" | curtin in-target -- tee /opt/syntropy/identity/owner/public.key
    - |
      echo "${community_key}" | curtin in-target -- tee /opt/syntropy/identity/community/private.key
    - |
      echo "${community_pub}" | curtin in-target -- tee /opt/syntropy/identity/community/public.key
    
    # Set proper permissions
    - curtin in-target -- chmod 600 /opt/syntropy/identity/owner/private.key
    - curtin in-target -- chmod 600 /opt/syntropy/identity/community/private.key
    - curtin in-target -- chmod 644 /opt/syntropy/identity/owner/public.key
    - curtin in-target -- chmod 644 /opt/syntropy/identity/community/public.key
    
    # Create node metadata
    - |
      cat > /target/opt/syntropy/metadata/node.json << 'EOJ'
      {
        "node_info": {
          "node_id": "${node_id}",
          "node_name": "${node_name}",
          "description": "${description}",
          "created_at": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
          "location": {
            "city": "${city}",
            "country": "${country}",
            "coordinates": "${coordinates}",
            "detection_method": "${detection_method}"
          }
        }
      }
      EOJ
    
    # Configure Docker
    - curtin in-target -- systemctl enable docker
    - curtin in-target -- usermod -aG docker admin
    
    # Configure firewall
    - curtin in-target -- ufw default deny incoming
    - curtin in-target -- ufw default allow outgoing
    - curtin in-target -- ufw allow ssh
    - curtin in-target -- ufw allow 9100/tcp  # Prometheus node exporter
    - curtin in-target -- ufw --force enable
    
    # Set up fail2ban
    - |
      curtin in-target -- tee /etc/fail2ban/jail.local << 'EOB'
      [sshd]
      enabled = true
      bantime = 3600
      findtime = 600
      maxretry = 5
      EOB
    
    # Enable services
    - curtin in-target -- systemctl enable ssh
    - curtin in-target -- systemctl enable fail2ban
    - curtin in-target -- systemctl enable prometheus-node-exporter
    
    # Create startup script
    - |
      curtin in-target -- tee /opt/syntropy/scripts/node-startup.sh << 'EOS'
      #!/bin/bash
      # Syntropy Node Startup Script
      
      # Update node status
      echo "$(date -u +%Y-%m-%dT%H:%M:%SZ): Node startup" >> /opt/syntropy/logs/startup.log
      
      # Additional startup tasks can be added here
      EOS
    - curtin in-target -- chmod +x /opt/syntropy/scripts/node-startup.sh
    
    # Create systemd service for startup script
    - |
      curtin in-target -- tee /etc/systemd/system/syntropy-node.service << 'EOU'
      [Unit]
      Description=Syntropy Node Startup Service
      After=network.target
      
      [Service]
      Type=oneshot
      ExecStart=/opt/syntropy/scripts/node-startup.sh
      RemainAfterExit=yes
      
      [Install]
      WantedBy=multi-user.target
      EOU
    
    # Enable startup service
    - curtin in-target -- systemctl enable syntropy-node
    
    # Final ownership adjustment
    - curtin in-target -- chown -R admin:admin /opt/syntropy

  # User data for first boot
  user-data:
    runcmd:
      - echo "Syntropy Node $(cat /opt/syntropy/metadata/node.json | jq -r .node_info.node_name) initialized at $(date -u)" >> /opt/syntropy/logs/init.log
EOF

    log DEBUG "User-data configuration created"
}

# Create meta-data configuration
create_meta_data_fixed() {
    local usb_mount="$1"
    local node_name="$2"
    
    log DEBUG "Creating meta-data configuration..."
    
    cat > "$usb_mount/autoinstall/meta-data" << EOF
instance-id: syntropy-${node_name}-$(date +%s)
local-hostname: ${node_name}
EOF

    log DEBUG "Meta-data configuration created"
}

# Create vendor-data configuration
create_vendor_data_fixed() {
    local usb_mount="$1"
    
    log DEBUG "Creating vendor-data configuration..."
    
    cat > "$usb_mount/autoinstall/vendor-data" << 'EOF'
#cloud-config
# Syntropy-specific vendor configuration
runcmd:
  - echo "Syntropy vendor configuration applied" >> /opt/syntropy/logs/vendor-init.log
EOF

    log DEBUG "Vendor-data configuration created"
}

# Create README with installation instructions
create_readme_fixed() {
    local usb_mount="$1"
    local node_name="$2"
    
    log DEBUG "Creating installation documentation..."
    
    cat > "$usb_mount/README.txt" << EOF
Syntropy Cooperative Grid - Node Installation
===========================================

Node Name: ${node_name}
Created: $(date -u +%Y-%m-%d\ %H:%M:%S\ UTC)

Installation Instructions:
------------------------
1. Insert this USB drive into the target machine
2. Boot from USB (you may need to change BIOS/UEFI boot order)
3. The installation will proceed automatically
4. When complete, the machine will reboot
5. The node will automatically register with the network

Initial Access:
-------------
- Username: admin
- Authentication: SSH key only (password login disabled)
- Default SSH port: 22

Important Locations:
-----------------
/opt/syntropy/
  ├── identity/    # Node identity and keys
  ├── platform/    # Platform-specific files
  ├── metadata/    # Node metadata
  ├── logs/        # System logs
  ├── config/      # Configuration files
  └── scripts/     # Management scripts

For more information, visit: https://docs.syntropy.network
EOF

    log DEBUG "Documentation created"
}

# Validate cloud-init configuration
validate_cloud_init() {
    local usb_mount="$1"
    
    log INFO "Validating cloud-init configuration..."
    
    local errors=0
    
    # Check required files
    for file in user-data meta-data vendor-data; do
        if [ ! -f "$usb_mount/autoinstall/$file" ]; then
            log ERROR "Missing required file: $file"
            errors=$((errors + 1))
        fi
    done
    
    # Validate YAML syntax
    if command -v python3 >/dev/null 2>&1; then
        if ! python3 -c "import yaml; yaml.safe_load(open('$usb_mount/autoinstall/user-data'))" 2>/dev/null; then
            log ERROR "Invalid YAML in user-data"
            errors=$((errors + 1))
        fi
    fi
    
    # Check for required sections in user-data
    local required_sections=("autoinstall" "network" "storage" "identity" "ssh")
    for section in "${required_sections[@]}"; do
        if ! grep -q "^[[:space:]]*$section:" "$usb_mount/autoinstall/user-data"; then
            log ERROR "Missing required section in user-data: $section"
            errors=$((errors + 1))
        fi
    done
    
    if [ $errors -eq 0 ]; then
        log SUCCESS "Cloud-init configuration validated successfully"
        return 0
    else
        log ERROR "Cloud-init validation failed with $errors errors"
        return 1
    fi
}