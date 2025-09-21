#!/bin/bash

# Syntropy Cooperative Grid - USB Device Preparation
# Version: 2.0.0

# Prepare USB device for installation
prepare_usb_device() {
    local device="$1"
    
    log INFO "Preparing USB device: $device"
    
    # Get and display device information
    display_device_info "$device"
    
    # Validate device requirements
    validate_device_requirements "$device"
    
    # Unmount existing partitions
    unmount_device_partitions "$device"
    
    # Wipe device completely
    wipe_device_clean "$device"
    
    # Create new partition structure
    create_partition_structure "$device"
    
    # Format the partition
    local partition=$(format_usb_partition "$device")
    
    log SUCCESS "USB device prepared successfully: $partition"
    echo "$partition"
    return 0
}

# Display comprehensive device information
display_device_info() {
    local device="$1"
    
    log INFO "Device Information Analysis:"
    
    # Basic device info
    local device_info=$(lsblk -d -n -o SIZE,MODEL "$device" 2>/dev/null || echo "Unknown Unknown")
    local device_size=$(echo "$device_info" | awk '{print $1}')
    local device_model=$(echo "$device_info" | awk '{for(i=2;i<=NF;i++) printf "%s ", $i; print ""}' | sed 's/[[:space:]]*$//')
    
    echo "  Device: $device"
    echo "  Size: $device_size"
    echo "  Model: $device_model"
    
    # Detailed information
    echo ""
    echo "  Current partition table:"
    sudo parted "$device" print 2>/dev/null || echo "    No partition table or unreadable"
    
    echo ""
    echo "  Current partitions:"
    lsblk -o NAME,SIZE,TYPE,FSTYPE,MOUNTPOINT "$device" 2>/dev/null || echo "    Unable to read partitions"
}

# Validate device meets requirements
validate_device_requirements() {
    local device="$1"
    
    log INFO "Validating device requirements..."
    
    # Check minimum size (8GB)
    local device_size_bytes=$(lsblk -b -d -n -o SIZE "$device" 2>/dev/null || echo "0")
    local min_size_bytes=$((8 * 1024 * 1024 * 1024))  # 8GB
    local device_size_gb=$((device_size_bytes / 1024 / 1024 / 1024))
    
    if [ "$device_size_bytes" -lt "$min_size_bytes" ]; then
        log ERROR "Device too small: ${device_size_gb}GB (minimum 8GB required)"
        return 1
    fi
    
    log SUCCESS "Size requirement met: ${device_size_gb}GB"
    
    # Check if device is writable
    if [ ! -w "$device" ]; then
        log ERROR "Device is not writable: $device"
        echo "This may require sudo privileges or the device may be write-protected"
        return 1
    fi
    
    log SUCCESS "Device is writable"
    
    # Check if device is busy
    if fuser "$device" 2>/dev/null; then
        log WARN "Device appears to be in use by other processes"
        echo "Attempting to free device..."
        # Try to kill processes using the device
        sudo fuser -k "$device" 2>/dev/null || true
        sleep 2
    fi
    
    return 0
}

# Unmount all partitions on device
unmount_device_partitions() {
    local device="$1"
    
    log INFO "Unmounting existing partitions..."
    
    # Get all mounted partitions for this device
    local mounted_partitions=$(mount | grep "^$device" | awk '{print $1}' || true)
    
    if [ -n "$mounted_partitions" ]; then
        echo "Found mounted partitions:"
        for partition in $mounted_partitions; do
            local mount_point=$(mount | grep "^$partition " | awk '{print $3}')
            echo "  $partition mounted at $mount_point"
            
            log INFO "Unmounting $partition..."
            if sudo umount "$partition" 2>/dev/null; then
                log SUCCESS "Unmounted $partition"
            else
                log WARN "Failed to unmount $partition, trying force unmount..."
                if sudo umount -f "$partition" 2>/dev/null; then
                    log SUCCESS "Force unmounted $partition"
                else
                    log WARN "Could not unmount $partition"
                fi
            fi
        done
    else
        log INFO "No mounted partitions found"
    fi
    
    # Additional cleanup - unmount any partition that might be numbered
    for i in {1..9}; do
        local partition="${device}${i}"
        if [ -b "$partition" ]; then
            sudo umount "$partition" 2>/dev/null || true
        fi
        
        # Also try p-style naming (for NVMe, etc.)
        partition="${device}p${i}"
        if [ -b "$partition" ]; then
            sudo umount "$partition" 2>/dev/null || true
        fi
    done
    
    # Wait for unmounting to complete
    sleep 2
    log SUCCESS "Partition unmounting completed"
}

# Completely wipe device clean
wipe_device_clean() {
    local device="$1"
    
    log INFO "Wiping device clean (this may take a moment)..."
    
    # Method 1: wipefs - Remove filesystem signatures
    if command -v wipefs >/dev/null 2>&1; then
        log DEBUG "Removing filesystem signatures..."
        sudo wipefs -a "$device" >/dev/null 2>&1 || true
    fi
    
    # Method 2: sgdisk - Zap GPT and MBR structures
    if command -v sgdisk >/dev/null 2>&1; then
        log DEBUG "Destroying GPT and MBR structures..."
        sudo sgdisk --zap-all "$device" >/dev/null 2>&1 || true
    fi
    
    # Method 3: dd - Zero out beginning and end of device
    log DEBUG "Clearing partition table areas..."
    sudo dd if=/dev/zero of="$device" bs=1M count=10 >/dev/null 2>&1 || true
    sudo dd if=/dev/zero of="$device" bs=1M seek=$(($(blockdev --getsz "$device") / 2048 - 10)) count=10 >/dev/null 2>&1 || true
    
    # Force kernel to re-read partition table
    sudo partprobe "$device" 2>/dev/null || true
    sleep 2
    
    log SUCCESS "Device wiped clean"
}

# Create new partition structure
create_partition_structure() {
    local device="$1"
    
    log INFO "Creating new partition structure..."
    
    # Create DOS/MBR partition table
    if ! sudo parted -s "$device" mklabel msdos; then
        log ERROR "Failed to create partition table"
        return 1
    fi
    
    log DEBUG "Created DOS partition table"
    
    # Create primary partition using all available space
    if ! sudo parted -s "$device" mkpart primary fat32 1MiB 100%; then
        log ERROR "Failed to create primary partition"
        return 1
    fi
    
    log DEBUG "Created primary partition"
    
    # Set boot flag
    if ! sudo parted -s "$device" set 1 boot on; then
        log ERROR "Failed to set boot flag"
        return 1
    fi
    
    log DEBUG "Set boot flag"
    
    # Force kernel to re-read partition table
    sudo partprobe "$device" 2>/dev/null || true
    sleep 3
    
    log SUCCESS "Partition structure created"
}

# Format USB partition
format_usb_partition() {
    local device="$1"
    
    log INFO "Formatting USB partition..."
    
    # Determine partition device name
    local partition=""
    local partition_candidates=("${device}1" "${device}p1")
    
    for candidate in "${partition_candidates[@]}"; do
        if [ -b "$candidate" ]; then
            partition="$candidate"
            break
        fi
    done
    
    if [ -z "$partition" ]; then
        log ERROR "Cannot find partition after creation"
        echo "Expected partition devices:"
        for candidate in "${partition_candidates[@]}"; do
            echo "  $candidate (exists: $([ -b "$candidate" ] && echo "yes" || echo "no"))"
        done
        return 1
    fi
    
    log DEBUG "Found partition: $partition"
    
    # Wait a bit more for partition to be fully recognized
    sleep 2
    
    # Format as FAT32 with SYNTROPY label
    log INFO "Formatting $partition as FAT32..."
    if ! sudo mkfs.fat -F32 -n "SYNTROPY" "$partition" >/dev/null 2>&1; then
        log ERROR "Failed to format partition $partition"
        
        # Try to diagnose the issue
        echo "Diagnostic information:"
        echo "  Partition exists: $([ -b "$partition" ] && echo "yes" || echo "no")"
        echo "  Partition info:"
        lsblk "$partition" 2>/dev/null || echo "    Cannot read partition info"
        
        return 1
    fi
    
    log SUCCESS "Partition formatted successfully"
    
    # Verify filesystem
    if ! sudo fsck.fat -v "$partition" >/dev/null 2>&1; then
        log WARN "Filesystem verification had warnings (this is usually normal for new filesystems)"
    else
        log SUCCESS "Filesystem verification passed"
    fi
    
    echo "$partition"
    return 0
}

# Mount USB partition for writing
mount_usb_partition() {
    local partition="$1"
    local mount_point="$2"
    
    log INFO "Mounting USB partition..."
    
    # Create mount point
    sudo mkdir -p "$mount_point"
    
    # Mount the partition
    if ! sudo mount "$partition" "$mount_point"; then
        log ERROR "Failed to mount $partition at $mount_point"
        return 1
    fi
    
    # Verify mount
    if ! mountpoint -q "$mount_point"; then
        log ERROR "Mount verification failed"
        return 1
    fi
    
    # Check available space
    local available_space=$(df -BM "$mount_point" | tail -1 | awk '{print $4}' | tr -d 'M')
    log INFO "Available space on USB: ${available_space}MB"
    
    # Verify we have enough space (need at least 2GB for Ubuntu)
    if [ "$available_space" -lt 2048 ]; then
        log ERROR "Insufficient space on USB: ${available_space}MB (need at least 2048MB)"
        sudo umount "$mount_point"
        return 1
    fi
    
    log SUCCESS "USB partition mounted at $mount_point"
    return 0
}

# Safely unmount USB partition
unmount_usb_partition() {
    local mount_point="$1"
    
    log INFO "Safely unmounting USB partition..."
    
    # Sync filesystem
    log DEBUG "Syncing filesystem..."
    sync
    sleep 2
    
    # Unmount
    if sudo umount "$mount_point" 2>/dev/null; then
        log SUCCESS "USB unmounted successfully"
    else
        log WARN "Normal unmount failed, trying force unmount..."
        if sudo umount -f "$mount_point" 2>/dev/null; then
            log SUCCESS "USB force unmounted successfully"
        else
            log ERROR "Failed to unmount USB"
            return 1
        fi
    fi
    
    # Remove mount point
    sudo rmdir "$mount_point" 2>/dev/null || true
    
    return 0
}

# Verify USB preparation
verify_usb_preparation() {
    local device="$1"
    
    log INFO "Verifying USB preparation..."
    
    # Check partition table
    local partition_table=$(sudo parted "$device" print 2>/dev/null | grep "Partition Table" | awk '{print $3}')
    if [ "$partition_table" != "msdos" ]; then
        log ERROR "Incorrect partition table type: $partition_table (expected msdos)"
        return 1
    fi
    
    # Check partition exists and is bootable
    local partition_info=$(sudo parted "$device" print 2>/dev/null | grep "^ 1")
    if [ -z "$partition_info" ]; then
        log ERROR "Primary partition not found"
        return 1
    fi
    
    if [[ "$partition_info" != *"boot"* ]]; then
        log ERROR "Boot flag not set on partition"
        return 1
    fi
    
    # Check filesystem
    local partition="${device}1"
    if [ ! -b "$partition" ]; then
        partition="${device}p1"
    fi
    
    if [ ! -b "$partition" ]; then
        log ERROR "Cannot find formatted partition"
        return 1
    fi
    
    local fstype=$(lsblk -n -o FSTYPE "$partition" 2>/dev/null)
    if [ "$fstype" != "vfat" ]; then
        log ERROR "Incorrect filesystem type: $fstype (expected vfat)"
        return 1
    fi
    
    log SUCCESS "USB preparation verified successfully"
    return 0
}

# Get USB device summary
get_usb_summary() {
    local device="$1"
    
    echo "USB Device Summary:"
    echo "  Device: $device"
    
    local device_info=$(lsblk -d -n -o SIZE,MODEL "$device" 2>/dev/null || echo "Unknown Unknown")
    local device_size=$(echo "$device_info" | awk '{print $1}')
    local device_model=$(echo "$device_info" | awk '{for(i=2;i<=NF;i++) printf "%s ", $i; print ""}' | sed 's/[[:space:]]*$//')
    
    echo "  Size: $device_size"
    echo "  Model: $device_model"
    
    # Show partition information
    echo "  Partitions:"
    lsblk -o NAME,SIZE,TYPE,FSTYPE "$device" 2>/dev/null | grep -v "^$device " | while read line; do
        echo "    $line"
    done
}