#!/bin/bash

# Syntropy Cooperative Grid - ISO Download and Management
# Version: 2.0.0

# Download and cache Ubuntu ISO
download_ubuntu_iso() {
    local iso_path="$ISO_CACHE_DIR/$ISO_FILE"
    
    log INFO "Setting up Ubuntu Server ISO..."
    
    # Ensure cache directory exists
    mkdir -p "$ISO_CACHE_DIR"
    
    # Check if ISO exists in cache and verify
    if [ -f "$iso_path" ]; then
        if verify_iso_integrity "$iso_path"; then
            log SUCCESS "Using cached ISO: $iso_path"
            copy_iso_to_work_dir "$iso_path"
            return 0
        else
            log WARN "Cached ISO failed verification, removing..."
            rm -f "$iso_path"
        fi
    fi
    
    # Download ISO
    if download_iso_from_server "$iso_path"; then
        if verify_iso_integrity "$iso_path"; then
            log SUCCESS "ISO downloaded and verified successfully"
            copy_iso_to_work_dir "$iso_path"
            return 0
        else
            log ERROR "Downloaded ISO failed verification"
            rm -f "$iso_path"
            return 1
        fi
    else
        log ERROR "Failed to download Ubuntu ISO"
        return 1
    fi
}

# Download ISO from Ubuntu servers
download_iso_from_server() {
    local iso_path="$1"
    
    log INFO "Downloading Ubuntu 22.04.4 Server (~1.5GB)..."
    echo "This may take several minutes depending on your connection..."
    
    # Try primary URL first
    if download_with_progress "$ISO_URL" "$iso_path"; then
        return 0
    fi
    
    # Try alternative mirrors if primary fails
    log WARN "Primary download failed, trying alternative mirrors..."
    
    local alternative_urls=(
        "http://archive.ubuntu.com/ubuntu/dists/jammy/main/installer-amd64/current/legacy-images/netboot/ubuntu-installer/amd64/linux"
        "http://cdimage.ubuntu.com/ubuntu/releases/22.04/release/$ISO_FILE"
        "http://mirror.math.princeton.edu/pub/ubuntu-iso/22.04/$ISO_FILE"
    )
    
    for url in "${alternative_urls[@]}"; do
        log INFO "Trying mirror: $url"
        if download_with_progress "$url" "$iso_path"; then
            return 0
        fi
        log WARN "Mirror failed, trying next..."
    done
    
    return 1
}

# Download with progress bar
download_with_progress() {
    local url="$1"
    local output_path="$2"
    
    # Use wget with progress bar
    if command -v wget >/dev/null 2>&1; then
        if wget --progress=bar:force:noscroll --timeout=30 --tries=3 -c -O "$output_path" "$url"; then
            return 0
        fi
    fi
    
    # Fallback to curl
    if command -v curl >/dev/null 2>&1; then
        if curl -L --progress-bar --connect-timeout 30 --max-time 3600 -C - -o "$output_path" "$url"; then
            return 0
        fi
    fi
    
    return 1
}

# Verify ISO integrity using SHA256
verify_iso_integrity() {
    local iso_path="$1"
    
    log INFO "Verifying ISO integrity..."
    
    if [ ! -f "$iso_path" ]; then
        log ERROR "ISO file not found: $iso_path"
        return 1
    fi
    
    # Calculate SHA256 checksum
    local calculated_sha256=""
    if command -v sha256sum >/dev/null 2>&1; then
        calculated_sha256=$(sha256sum "$iso_path" | cut -d' ' -f1)
    elif command -v shasum >/dev/null 2>&1; then
        calculated_sha256=$(shasum -a 256 "$iso_path" | cut -d' ' -f1)
    else
        log ERROR "No SHA256 checksum tool available"
        return 1
    fi
    
    if [ "$calculated_sha256" = "$ISO_SHA256" ]; then
        log SUCCESS "ISO integrity verified (SHA256 match)"
        return 0
    else
        log ERROR "ISO checksum verification failed!"
        echo "Expected: $ISO_SHA256"
        echo "Got:      $calculated_sha256"
        return 1
    fi
}

# Copy ISO from cache to work directory
copy_iso_to_work_dir() {
    local iso_path="$1"
    
    log INFO "Copying ISO to work directory..."
    if cp "$iso_path" "$WORK_DIR/$ISO_FILE"; then
        log SUCCESS "ISO ready for use: $WORK_DIR/$ISO_FILE"
        return 0
    else
        log ERROR "Failed to copy ISO to work directory"
        return 1
    fi
}

# Install Ubuntu ISO contents to USB
install_ubuntu_to_usb() {
    local usb_partition="$1"
    
    log INFO "Installing Ubuntu to USB device..."
    
    # Mount USB
    USB_MOUNT="/mnt/syntropy-usb"
    sudo mkdir -p "$USB_MOUNT"
    if ! sudo mount "$usb_partition" "$USB_MOUNT"; then
        log ERROR "Failed to mount USB partition: $usb_partition"
        return 1
    fi
    
    # Extract ISO
    extract_iso_to_usb
    
    # Make bootable
    make_usb_bootable
    
    return 0
}

# Extract ISO contents to USB
extract_iso_to_usb() {
    log INFO "Extracting Ubuntu ISO contents..."
    
    local iso_mount_dir="$WORK_DIR/iso-mount"
    mkdir -p "$iso_mount_dir"
    
    # Mount ISO
    if ! sudo mount -o loop "$WORK_DIR/$ISO_FILE" "$iso_mount_dir"; then
        log ERROR "Failed to mount ISO file"
        return 1
    fi
    
    # Copy all contents
    if ! sudo cp -r "$iso_mount_dir"/* "$USB_MOUNT/"; then
        log ERROR "Failed to copy ISO contents to USB"
        sudo umount "$iso_mount_dir"
        return 1
    fi
    
    # Unmount ISO
    sudo umount "$iso_mount_dir"
    rm -rf "$iso_mount_dir"
    
    log SUCCESS "ISO contents extracted to USB"
    return 0
}

# Make USB bootable
make_usb_bootable() {
    log INFO "Configuring USB as bootable..."
    
    # Ensure proper permissions
    sudo chmod -R 755 "$USB_MOUNT"
    
    # Update boot configuration if needed
    if [ -f "$USB_MOUNT/isolinux/isolinux.cfg" ]; then
        log DEBUG "Found isolinux configuration"
    fi
    
    if [ -f "$USB_MOUNT/boot/grub/grub.cfg" ]; then
        log DEBUG "Found GRUB configuration"
    fi
    
    # Create boot flag file
    sudo touch "$USB_MOUNT/.syntropy_boot"
    
    log SUCCESS "USB configured as bootable"
    return 0
}

# Add documentation to USB
add_usb_documentation() {
    local node_name="$1"
    
    log INFO "Adding documentation to USB..."
    
    # Create documentation directory
    sudo mkdir -p "$USB_MOUNT/syntropy-docs"
    
    # Copy node summary if it exists
    if [ -f "$HOME/.syntropy/nodes/${node_name}_summary.md" ]; then
        sudo cp "$HOME/.syntropy/nodes/${node_name}_summary.md" "$USB_MOUNT/syntropy-docs/" 2>/dev/null || true
    fi
    
    # Create USB-specific README
    create_usb_readme "$node_name"
    
    log SUCCESS "Documentation added to USB"
}

# Create README for USB
create_usb_readme() {
    local node_name="$1"
    
    sudo tee "$USB_MOUNT/syntropy-docs/README.md" > /dev/null << README_EOF
# Syntropy Cooperative Grid - Installation USB

This USB contains an automated Ubuntu installation configured for the Syntropy Cooperative Grid.

## Node Information
- **Node Name**: $node_name
- **Created**: $(date)
- **Platform Version**: 2.0.0
- **ISO Version**: $ISO_FILE

## Installation Process

### 1. Hardware Preparation
- Ensure target hardware meets minimum requirements:
  - x86_64 architecture (Intel/AMD 64-bit)
  - Minimum 8GB RAM (4GB+ recommended)
  - Minimum 32GB storage (64GB+ recommended)
  - Network interface (Ethernet or WiFi)

### 2. BIOS/UEFI Configuration
- Insert this USB into target hardware
- Boot into BIOS/UEFI setup (usually F2, F12, or DEL during startup)
- Configure boot order to prioritize USB devices
- Disable Secure Boot if necessary
- Save settings and exit

### 3. Installation
- Boot from this USB device
- Installation will proceed automatically (~20-30 minutes)
- No user interaction required
- System will reboot when installation is complete

### 4. Post-Installation
- Node will be accessible via SSH using the generated keys
- Use connection tools on the creation system to access node
- All services will be configured and running automatically

## What Gets Installed

- **Operating System**: Ubuntu 22.04.4 Server
- **Container Runtime**: Docker with automatic startup
- **Security**: SSH key-only authentication, UFW firewall, fail2ban
- **Monitoring**: Prometheus node exporter on port 9100
- **Platform**: Syntropy templates and tools in /opt/syntropy/

## Management

### From Creation System
Use the management tools created when this USB was generated:

\`\`\`bash
# Quick connection
$HOME/.syntropy/connect-${node_name}.sh

# Node manager
syntropy-manager.sh connect $node_name
syntropy-manager.sh status $node_name
\`\`\`

### Direct SSH Access
Once you have the node's IP address:

\`\`\`bash
ssh -i $HOME/.syntropy/keys/${node_name}_owner.key admin@<NODE_IP>
\`\`\`

## Troubleshooting

### Boot Issues
- Verify BIOS/UEFI boot order settings
- Check that Secure Boot is disabled
- Ensure USB device is properly connected
- Try different USB ports

### Installation Issues
- Verify hardware meets minimum requirements
- Check network connectivity (DHCP required)
- Ensure sufficient storage space
- Monitor installation logs if accessible

### Network Issues
- Verify DHCP is available on target network
- Check firewall settings on network
- Ensure SSH port (22) is not blocked

## Support

For additional help:
- Check the Syntropy documentation
- Review system logs: \`journalctl -u syntropy-first-boot\`
- Test connectivity: \`ping <node-ip>\` and \`nc -zv <node-ip> 22\`

---
**Created by**: $(hostname) on $(date)  
**USB Creator Version**: 2.0.0
README_EOF
}

# Clean up ISO mount points
cleanup_iso_mounts() {
    log INFO "Cleaning up ISO mount points..."
    
    # Clean up any remaining mounts
    local mount_points=$(mount | grep "$WORK_DIR" | awk '{print $3}' || true)
    for mount_point in $mount_points; do
        log DEBUG "Unmounting: $mount_point"
        sudo umount "$mount_point" 2>/dev/null || true
    done
    
    # Remove temporary directories
    rm -rf "$WORK_DIR/iso-mount" 2>/dev/null || true
    
    log SUCCESS "ISO cleanup completed"
}

# Get ISO information
get_iso_info() {
    echo "Ubuntu ISO Information:"
    echo "  File: $ISO_FILE"
    echo "  URL: $ISO_URL"
    echo "  SHA256: $ISO_SHA256"
    echo "  Cache Directory: $ISO_CACHE_DIR"
    
    if [ -f "$ISO_CACHE_DIR/$ISO_FILE" ]; then
        local file_size=$(du -h "$ISO_CACHE_DIR/$ISO_FILE" | cut -f1)
        echo "  Cached: Yes ($file_size)"
    else
        echo "  Cached: No"
    fi
}