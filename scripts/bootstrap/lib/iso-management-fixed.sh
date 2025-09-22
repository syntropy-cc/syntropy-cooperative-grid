#!/bin/bash

# Syntropy Cooperative Grid - Enhanced ISO Management
# Version: 2.1.0

# Fixed ISO management functions with enhanced validation and error handling

# Global variables
ISO_CACHE_DIR="$HOME/.syntropy/cache/iso"
ISO_VERSION="22.04.3"  # Ubuntu Server LTS version
ISO_NAME="ubuntu-${ISO_VERSION}-live-server-amd64.iso"
ISO_URL="https://releases.ubuntu.com/${ISO_VERSION}/${ISO_NAME}"
ISO_CHECKSUM_URL="https://releases.ubuntu.com/${ISO_VERSION}/SHA256SUMS"

# Initialize ISO environment
init_iso_environment() {
    mkdir -p "$ISO_CACHE_DIR"
    log DEBUG "ISO cache directory: $ISO_CACHE_DIR"
}

# Download Ubuntu ISO with validation
download_ubuntu_iso() {
    local iso_path="$ISO_CACHE_DIR/$ISO_NAME"
    local checksum_file="$ISO_CACHE_DIR/SHA256SUMS"
    
    # Skip if ISO exists and is valid
    if [ -f "$iso_path" ] && validate_iso "$iso_path"; then
        log INFO "Using cached ISO: $iso_path"
        return 0
    fi
    
    log INFO "Downloading Ubuntu Server ${ISO_VERSION} ISO..."
    
    # Download checksum file
    if ! wget -q "$ISO_CHECKSUM_URL" -O "$checksum_file"; then
        log ERROR "Failed to download checksum file"
        return 1
    fi
    
    # Download ISO with progress
    if ! wget --progress=bar:force "$ISO_URL" -O "$iso_path"; then
        log ERROR "Failed to download ISO"
        rm -f "$iso_path"
        return 1
    fi
    
    # Validate downloaded ISO
    if validate_iso "$iso_path"; then
        log SUCCESS "ISO downloaded and verified successfully"
        return 0
    else
        log ERROR "ISO validation failed"
        rm -f "$iso_path"
        return 1
    fi
}

# Validate ISO checksum
validate_iso() {
    local iso_path="$1"
    local checksum_file="$ISO_CACHE_DIR/SHA256SUMS"
    
    log DEBUG "Validating ISO checksum..."
    
    if [ ! -f "$checksum_file" ]; then
        log ERROR "Checksum file not found"
        return 1
    fi
    
    # Extract expected checksum
    local expected_sum=$(grep "$ISO_NAME" "$checksum_file" | cut -d' ' -f1)
    if [ -z "$expected_sum" ]; then
        log ERROR "Could not find checksum for ISO"
        return 1
    }
    
    # Calculate actual checksum
    local actual_sum=$(sha256sum "$iso_path" | cut -d' ' -f1)
    
    # Compare checksums
    if [ "$expected_sum" = "$actual_sum" ]; then
        log DEBUG "ISO checksum verified"
        return 0
    else
        log ERROR "ISO checksum mismatch"
        return 1
    fi
}

# Install Ubuntu ISO contents to USB with enhanced error handling
install_ubuntu_to_usb_fixed() {
    local usb_partition="$1"
    
    log INFO "Installing Ubuntu to USB device..."
    
    # Prepare mount points
    USB_MOUNT="/mnt/syntropy-usb"
    ISO_MOUNT="/mnt/syntropy-iso"
    
    sudo mkdir -p "$USB_MOUNT" "$ISO_MOUNT"
    
    # Mount USB
    if ! sudo mount "$usb_partition" "$USB_MOUNT"; then
        log ERROR "Failed to mount USB partition: $usb_partition"
        cleanup_mounts
        return 1
    fi
    
    # Mount ISO
    if ! sudo mount -o loop "$ISO_CACHE_DIR/$ISO_NAME" "$ISO_MOUNT"; then
        log ERROR "Failed to mount ISO"
        cleanup_mounts
        return 1
    fi
    
    # Extract ISO with progress tracking
    local total_files=$(find "$ISO_MOUNT" -type f | wc -l)
    local current_file=0
    
    log INFO "Copying ISO contents (this may take several minutes)..."
    
    find "$ISO_MOUNT" -type f | while read -r file; do
        current_file=$((current_file + 1))
        progress=$((current_file * 100 / total_files))
        
        # Get relative path
        rel_path="${file#$ISO_MOUNT/}"
        target="$USB_MOUNT/$rel_path"
        
        # Create directory structure
        sudo mkdir -p "$(dirname "$target")"
        
        # Copy file
        if ! sudo cp -a "$file" "$target"; then
            log WARN "Failed to copy: $rel_path"
        fi
        
        # Update progress
        if [ $((current_file % 100)) -eq 0 ]; then
            log DEBUG "Progress: $progress% ($current_file/$total_files files)"
        fi
    done
    
    # Copy boot files specifically
    log INFO "Installing boot files..."
    
    # Copy UEFI boot files
    if [ -d "$ISO_MOUNT/EFI" ]; then
        sudo cp -ar "$ISO_MOUNT/EFI" "$USB_MOUNT/"
    else
        log WARN "UEFI boot files not found in ISO"
    fi
    
    # Install GRUB for UEFI
    if [ -d "$USB_MOUNT/EFI" ]; then
        log DEBUG "Configuring UEFI boot..."
        sudo grub-install --target=x86_64-efi --efi-directory="$USB_MOUNT/EFI" \
            --boot-directory="$USB_MOUNT/boot" --removable 2>/dev/null || \
            log WARN "UEFI boot setup had warnings"
    fi
    
    # Install Syslinux for Legacy BIOS
    if command -v syslinux >/dev/null 2>&1; then
        log DEBUG "Installing Legacy BIOS boot..."
        local device=$(echo "$usb_partition" | sed 's/[0-9]*$//')
        
        sudo syslinux --install "$usb_partition" 2>/dev/null || \
            log WARN "Syslinux installation had warnings"
            
        if [ -f /usr/lib/syslinux/mbr/mbr.bin ]; then
            sudo dd if=/usr/lib/syslinux/mbr/mbr.bin of="$device" bs=440 count=1 conv=notrunc 2>/dev/null
        elif [ -f /usr/share/syslinux/mbr.bin ]; then
            sudo dd if=/usr/share/syslinux/mbr.bin of="$device" bs=440 count=1 conv=notrunc 2>/dev/null
        fi
    fi
    
    # Cleanup and finalize
    cleanup_mounts
    
    log SUCCESS "Ubuntu installation completed successfully"
    return 0
}

# Cleanup mount points
cleanup_mounts() {
    log DEBUG "Cleaning up mount points..."
    sudo umount "$ISO_MOUNT" 2>/dev/null || true
    sudo umount "$USB_MOUNT" 2>/dev/null || true
    sudo rm -rf "$ISO_MOUNT" "$USB_MOUNT"
}

# Extract ISO contents properly with progress
extract_iso_to_usb_proper() {
    log INFO "Extracting Ubuntu ISO contents..."
    
    local iso_path="$ISO_CACHE_DIR/$ISO_NAME"
    
    if [ ! -f "$iso_path" ]; then
        log ERROR "ISO file not found: $iso_path"
        return 1
    fi
    
    # Calculate total size for progress
    local total_size=$(du -b "$iso_path" | cut -f1)
    local processed_size=0
    
    # Use dd with status=progress for better feedback
    if ! sudo dd if="$iso_path" of="$USB_DEVICE" bs=4M status=progress conv=fsync; then
        log ERROR "Failed to write ISO to USB"
        return 1
    fi
    
    # Sync to ensure all writes are complete
    sync
    
    log SUCCESS "ISO contents extracted successfully"
    return 0
}

# Make USB bootable with hybrid UEFI/Legacy support
make_usb_bootable_hybrid() {
    log INFO "Configuring hybrid boot support..."
    
    local device="$1"
    local efi_partition="${device}1"
    
    # Create hybrid MBR/GPT
    if ! sudo sgdisk --hybrid 1 "$device"; then
        log WARN "Hybrid MBR creation failed (non-critical)"
    fi
    
    # Install GRUB for UEFI
    if [ -d "/usr/lib/grub/x86_64-efi" ]; then
        log DEBUG "Installing GRUB for UEFI..."
        sudo grub-install --target=x86_64-efi --efi-directory="$USB_MOUNT/EFI" \
            --boot-directory="$USB_MOUNT/boot" --removable || \
            log WARN "GRUB UEFI installation had warnings"
    fi
    
    # Install Syslinux for Legacy BIOS
    if command -v syslinux >/dev/null 2>&1; then
        log DEBUG "Installing Syslinux for Legacy BIOS..."
        sudo syslinux --install "$efi_partition" || \
            log WARN "Syslinux installation had warnings"
            
        # Install MBR
        if [ -f /usr/lib/syslinux/mbr/mbr.bin ]; then
            sudo dd if=/usr/lib/syslinux/mbr/mbr.bin of="$device" bs=440 count=1 conv=notrunc
        elif [ -f /usr/share/syslinux/mbr.bin ]; then
            sudo dd if=/usr/share/syslinux/mbr.bin of="$device" bs=440 count=1 conv=notrunc
        fi
    fi
    
    # Mark partition as bootable
    sudo parted "$device" set 1 boot on
    sudo parted "$device" set 1 esp on
    
    log SUCCESS "Hybrid boot configuration completed"
    return 0
}

# Create proper bootable USB structure
create_bootable_usb_structure() {
    local usb_partition="$1"
    
    log INFO "Creating bootable USB structure..."
    
    # Mount USB
    USB_MOUNT="/mnt/syntropy-usb"
    sudo mkdir -p "$USB_MOUNT"
    
    if ! sudo mount "$usb_partition" "$USB_MOUNT"; then
        log ERROR "Failed to mount USB partition"
        return 1
    fi
    
    # Create necessary directories
    sudo mkdir -p "$USB_MOUNT"/{boot/{grub,syslinux},EFI/BOOT,autoinstall}
    
    # Create basic files
    echo "Syntropy Cooperative Grid Node Installer" | \
        sudo tee "$USB_MOUNT/.disk/info" > /dev/null
    
    # Unmount
    sudo umount "$USB_MOUNT"
    sudo rm -rf "$USB_MOUNT"
    
    log SUCCESS "USB structure created successfully"
    return 0
}