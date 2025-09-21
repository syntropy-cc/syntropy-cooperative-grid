#!/bin/bash

# Syntropy Cooperative Grid - Security and Key Management
# Version: 2.0.0

# Global variables for key management
OWNER_KEY_PATH=""
COMMUNITY_KEY_PATH=""
OWNER_FINGERPRINT=""
COMMUNITY_FINGERPRINT=""

# Setup security keys (owner and community)
setup_security_keys() {
    local owner_key_file="$1"
    local location_node_id="$2"
    
    log INFO "Setting up security keys..."
    
    # Setup owner key
    setup_owner_key "$owner_key_file"
    
    # Setup community key
    setup_community_key "$location_node_id"
    
    # Validate both keys
    validate_generated_keys
    
    log SUCCESS "Security keys configured successfully"
}

# Setup owner key (for SSH access and management)
setup_owner_key() {
    local owner_key_file="$1"
    
    OWNER_KEY_PATH="$WORK_DIR/owner_key"
    
    if [ -n "$owner_key_file" ]; then
        if [ ! -f "$owner_key_file" ]; then
            log ERROR "Owner key file not found: $owner_key_file"
            exit 1
        fi
        
        # Validate key format
        if ! validate_ssh_key "$owner_key_file"; then
            log ERROR "Invalid SSH private key format: $owner_key_file"
            exit 1
        fi
        
        log INFO "Using existing owner key: $owner_key_file"
        cp "$owner_key_file" "$OWNER_KEY_PATH"
        
        if [ -f "${owner_key_file}.pub" ]; then
            cp "${owner_key_file}.pub" "${OWNER_KEY_PATH}.pub"
        else
            ssh-keygen -y -f "$OWNER_KEY_PATH" > "${OWNER_KEY_PATH}.pub"
        fi
        
        OWNER_FINGERPRINT=$(ssh-keygen -lf "${OWNER_KEY_PATH}.pub" | awk '{print $2}')
        log SUCCESS "Using existing owner key (fingerprint: $OWNER_FINGERPRINT)"
    else
        log INFO "Generating new owner key for multi-node management..."
        
        # Generate new ed25519 key
        if ! ssh-keygen -t ed25519 -f "$OWNER_KEY_PATH" -N "" -C "syntropy-owner-$(date +%Y%m%d)-$(hostname)"; then
            log ERROR "Failed to generate owner key"
            exit 1
        fi
        
        OWNER_FINGERPRINT=$(ssh-keygen -lf "${OWNER_KEY_PATH}.pub" | awk '{print $2}')
        log SUCCESS "New owner key generated (fingerprint: $OWNER_FINGERPRINT)"
        log INFO "Save this key for future nodes: $OWNER_KEY_PATH"
    fi
}

# Setup community key (for inter-node communication)
setup_community_key() {
    local location_node_id="$1"
    
    COMMUNITY_KEY_PATH="$WORK_DIR/community-$location_node_id"
    
    log INFO "Generating community key for node communication..."
    
    if ! ssh-keygen -t ed25519 -f "$COMMUNITY_KEY_PATH" -N "" -C "syntropy-community-$location_node_id"; then
        log ERROR "Failed to generate community key"
        exit 1
    fi
    
    COMMUNITY_FINGERPRINT=$(ssh-keygen -lf "${COMMUNITY_KEY_PATH}.pub" | awk '{print $2}')
    log SUCCESS "Community key generated (fingerprint: $COMMUNITY_FINGERPRINT)"
}

# Validate SSH key format and integrity
validate_ssh_key() {
    local key_file="$1"
    
    # Check if file exists and is readable
    if [ ! -r "$key_file" ]; then
        log ERROR "Key file is not readable: $key_file"
        return 1
    fi
    
    # Check if it's a valid SSH private key
    if ! ssh-keygen -l -f "$key_file" >/dev/null 2>&1; then
        log ERROR "Invalid SSH private key format"
        return 1
    fi
    
    # Check key algorithm (prefer ed25519)
    local key_type=$(ssh-keygen -l -f "$key_file" | awk '{print $4}' | tr -d '()')
    case "$key_type" in
        "ED25519")
            log INFO "Key algorithm: ed25519 (recommended)"
            ;;
        "RSA")
            local key_bits=$(ssh-keygen -l -f "$key_file" | awk '{print $1}')
            if [ "$key_bits" -lt 2048 ]; then
                log WARN "RSA key is less than 2048 bits, consider using ed25519"
            else
                log INFO "Key algorithm: RSA $key_bits bits"
            fi
            ;;
        "ECDSA")
            log INFO "Key algorithm: ECDSA"
            ;;
        *)
            log WARN "Unknown key algorithm: $key_type"
            ;;
    esac
    
    return 0
}

# Validate that generated keys are correct
validate_generated_keys() {
    log INFO "Validating generated keys..."
    
    # Validate owner key
    if [ ! -f "$OWNER_KEY_PATH" ] || [ ! -f "${OWNER_KEY_PATH}.pub" ]; then
        log ERROR "Owner key files missing"
        exit 1
    fi
    
    if ! validate_ssh_key "$OWNER_KEY_PATH"; then
        log ERROR "Generated owner key is invalid"
        exit 1
    fi
    
    # Validate community key
    if [ ! -f "$COMMUNITY_KEY_PATH" ] || [ ! -f "${COMMUNITY_KEY_PATH}.pub" ]; then
        log ERROR "Community key files missing"
        exit 1
    fi
    
    if ! validate_ssh_key "$COMMUNITY_KEY_PATH"; then
        log ERROR "Generated community key is invalid"
        exit 1
    fi
    
    # Ensure keys are different
    if [ "$OWNER_FINGERPRINT" = "$COMMUNITY_FINGERPRINT" ]; then
        log ERROR "Owner and community keys have same fingerprint"
        exit 1
    fi
    
    log SUCCESS "All keys validated successfully"
}

# Store keys in local management directory
store_keys_locally() {
    local node_name="$1"
    
    log INFO "Storing keys in local management directory..."
    
    # Copy owner keys
    cp "$OWNER_KEY_PATH" "$KEYS_DIR/${node_name}_owner.key"
    cp "${OWNER_KEY_PATH}.pub" "$KEYS_DIR/${node_name}_owner.pub"
    
    # Copy community keys
    cp "$COMMUNITY_KEY_PATH" "$KEYS_DIR/${node_name}_community.key"
    cp "${COMMUNITY_KEY_PATH}.pub" "$KEYS_DIR/${node_name}_community.pub"
    
    # Set proper permissions
    chmod 600 "$KEYS_DIR/${node_name}_owner.key"
    chmod 600 "$KEYS_DIR/${node_name}_community.key"
    chmod 644 "$KEYS_DIR/${node_name}_owner.pub"
    chmod 644 "$KEYS_DIR/${node_name}_community.pub"
    
    log SUCCESS "Keys stored locally with proper permissions"
}

# Generate key information metadata
generate_key_metadata() {
    local node_name="$1"
    
    cat << KEY_METADATA_EOF
{
  "owner_key": {
    "fingerprint": "$OWNER_FINGERPRINT",
    "algorithm": "ed25519",
    "purpose": "ssh_access_and_management",
    "created": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
    "local_path": "$KEYS_DIR/${node_name}_owner.key"
  },
  "community_key": {
    "fingerprint": "$COMMUNITY_FINGERPRINT", 
    "algorithm": "ed25519",
    "purpose": "inter_node_communication",
    "created": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
    "local_path": "$KEYS_DIR/${node_name}_community.key"
  },
  "security_notes": {
    "owner_key_reused": $([ -n "$OWNER_KEY_FILE" ] && echo "true" || echo "false"),
    "password_authentication": "disabled",
    "key_rotation_recommended": "every_90_days"
  }
}
KEY_METADATA_EOF
}

# Generate SSH authorized_keys content for the node
generate_authorized_keys() {
    cat "${OWNER_KEY_PATH}.pub"
}

# Create security configuration for cloud-init
create_security_config() {
    cat << SECURITY_CONFIG_EOF
# SSH Configuration
ssh:
  install-server: true
  allow-pw: false
  authorized-keys:
    - $(cat "${OWNER_KEY_PATH}.pub")

# Disable password authentication
chpasswd:
  expire: false

# Lock root account
disable_root: true

# Security packages
packages:
  - fail2ban
  - ufw
  - unattended-upgrades

# Firewall configuration
runcmd:
  - ufw default deny incoming
  - ufw default allow outgoing
  - ufw allow ssh
  - ufw allow 9100/tcp
  - ufw --force enable
  - systemctl enable fail2ban
  - systemctl start fail2ban
SECURITY_CONFIG_EOF
}

# Validate key file permissions
check_key_permissions() {
    local key_file="$1"
    
    if [ ! -f "$key_file" ]; then
        log ERROR "Key file does not exist: $key_file"
        return 1
    fi
    
    local permissions=$(stat -c "%a" "$key_file" 2>/dev/null || stat -f "%A" "$key_file" 2>/dev/null)
    
    # Private keys should be 600, public keys can be 644
    if [[ "$key_file" == *.pub ]]; then
        if [ "$permissions" != "644" ] && [ "$permissions" != "600" ]; then
            log WARN "Public key has unusual permissions: $permissions (expected 644)"
        fi
    else
        if [ "$permissions" != "600" ]; then
            log WARN "Private key has unsafe permissions: $permissions (should be 600)"
            log INFO "Fixing permissions..."
            chmod 600 "$key_file"
        fi
    fi
    
    return 0
}

# Generate security summary for user
generate_security_summary() {
    local node_name="$1"
    
    echo ""
    log INFO "Security Configuration Summary:"
    echo "  Owner Key Fingerprint: $OWNER_FINGERPRINT"
    echo "  Community Key Fingerprint: $COMMUNITY_FINGERPRINT"
    echo "  SSH Authentication: Key-only (passwords disabled)"
    echo "  Firewall: Enabled with SSH access only"
    echo "  Intrusion Detection: fail2ban enabled"
    echo "  Local Key Storage: $KEYS_DIR/${node_name}_*.key"
    echo ""
    
    if [ -n "$OWNER_KEY_FILE" ]; then
        echo "  Multi-Node Management: Enabled (reusing owner key)"
        echo "  Additional nodes can use: --owner-key $KEYS_DIR/${node_name}_owner.key"
    else
        echo "  Multi-Node Management: Available"
        echo "  For additional nodes use: --owner-key $KEYS_DIR/${node_name}_owner.key"
    fi
    echo ""
}

# Cleanup sensitive data from work directory
cleanup_sensitive_data() {
    log INFO "Cleaning up sensitive data from work directory..."
    
    # Securely remove private keys from work directory
    if [ -f "$OWNER_KEY_PATH" ]; then
        shred -vfz -n 3 "$OWNER_KEY_PATH" 2>/dev/null || rm -f "$OWNER_KEY_PATH"
    fi
    
    if [ -f "$COMMUNITY_KEY_PATH" ]; then
        shred -vfz -n 3 "$COMMUNITY_KEY_PATH" 2>/dev/null || rm -f "$COMMUNITY_KEY_PATH"
    fi
    
    log SUCCESS "Sensitive data cleanup completed"
}