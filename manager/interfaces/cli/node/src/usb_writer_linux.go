//go:build linux
// +build linux

package node

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"node-component/src/internal/helpers"
	"node-component/src/internal/types"
)

// USBWriterLinux implements USB writing for Linux systems
type USBWriterLinux struct {
	*USBWriterBase
}

// NewUSBWriterLinux creates a new Linux USB writer
func NewUSBWriterLinux(logger types.Logger) *USBWriterLinux {
	base := NewUSBWriterBase("linux", logger)
	return &USBWriterLinux{
		USBWriterBase: base,
	}
}

// WriteISO writes an ISO to a USB device with cloud-init injection
func (uwl *USBWriterLinux) WriteISO(ctx context.Context, isoPath string, devicePath string, cloudInitConfig *types.CloudInitConfig) (*types.WriteResult, error) {
	uwl.logger.Info("Writing ISO to USB device", "iso", isoPath, "device", devicePath)

	startTime := time.Now()

	// Validate inputs
	if err := uwl.validateInputs(isoPath, devicePath); err != nil {
		return nil, fmt.Errorf("input validation failed: %w", err)
	}

	// Get ISO file size
	isoFile, err := os.Open(isoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open ISO file: %w", err)
	}
	defer isoFile.Close()

	isoInfo, err := isoFile.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get ISO file info: %w", err)
	}

	// Create progress tracker
	progressTracker := NewWriteProgressTracker(uwl.logger, isoInfo.Size())

	// Write ISO to device using dd
	actualISOPath := isoPath

	// Write ISO to device using dd
	if err := uwl.writeISOWithDD(ctx, actualISOPath, devicePath, progressTracker); err != nil {
		return &types.WriteResult{
			DevicePath:   devicePath,
			ISOPath:      isoPath,
			BytesWritten: 0,
			Duration:     time.Since(startTime),
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	// Sync to ensure data is written to disk
	if err := uwl.syncDevice(devicePath); err != nil {
		uwl.logger.Warn("Failed to sync device", "device", devicePath, "error", err)
	}

	// Create NoCloud partition for cloud-init if provided
	if cloudInitConfig != nil {
		uwl.logger.Info("Creating NoCloud partition for cloud-init")
		cloudInitInjector := NewCloudInitInjector(uwl.logger)
		if err := cloudInitInjector.CreateNoCloudPartition(devicePath, cloudInitConfig); err != nil {
			uwl.logger.Error("Failed to create NoCloud partition", "error", err)
			return nil, fmt.Errorf("failed to create cloud-init partition: %w", err)
		}
	}

	// Clean up injected ISO if created
	if actualISOPath != isoPath {
		os.Remove(actualISOPath)
	}

	duration := time.Since(startTime)
	uwl.logger.Info("ISO write completed successfully",
		"device", devicePath,
		"duration", duration,
		"bytes_written", isoInfo.Size())

	return &types.WriteResult{
		DevicePath:        devicePath,
		ISOPath:           isoPath,
		BytesWritten:      isoInfo.Size(),
		Duration:          duration,
		Success:           true,
		CloudInitInjected: cloudInitConfig != nil,
	}, nil
}

// UnmountDevice unmounts a device before writing
func (uwl *USBWriterLinux) UnmountDevice(ctx context.Context, devicePath string) error {
	uwl.logger.Debug("Unmounting device", "device", devicePath)

	// Get device name from path (for future use)
	_ = helpers.GetDeviceNameFromPath(devicePath)

	// Check if device has mounted partitions
	cmd := exec.CommandContext(ctx, "lsblk", "-o", "MOUNTPOINT", "-n", "-r", devicePath)
	output, err := cmd.Output()
	if err != nil {
		uwl.logger.Warn("Failed to check device mounts", "device", devicePath, "error", err)
		return nil // Continue even if check fails
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && line != "MOUNTPOINT" {
			// Device has mounted partitions, unmount them
			uwl.logger.Debug("Unmounting partition", "mountpoint", line)
			cmd := exec.CommandContext(ctx, "umount", line)
			if err := cmd.Run(); err != nil {
				uwl.logger.Warn("Failed to unmount partition", "mountpoint", line, "error", err)
			}
		}
	}

	// Also try to unmount the device itself
	cmd = exec.CommandContext(ctx, "umount", devicePath)
	if err := cmd.Run(); err != nil {
		uwl.logger.Debug("Device not mounted or already unmounted", "device", devicePath)
	}

	uwl.logger.Debug("Device unmounted successfully", "device", devicePath)
	return nil
}

// MountDevice mounts a device after writing
func (uwl *USBWriterLinux) MountDevice(ctx context.Context, devicePath string) error {
	uwl.logger.Debug("Mounting device", "device", devicePath)

	// Create mount point
	mountPoint := filepath.Join("/tmp", "syntropy-"+helpers.GetDeviceNameFromPath(devicePath))
	if err := os.MkdirAll(mountPoint, 0755); err != nil {
		return fmt.Errorf("failed to create mount point: %w", err)
	}

	// Mount device
	cmd := exec.CommandContext(ctx, "mount", devicePath, mountPoint)
	if err := cmd.Run(); err != nil {
		uwl.logger.Warn("Failed to mount device", "device", devicePath, "mountpoint", mountPoint, "error", err)
		return nil // Continue even if mount fails
	}

	uwl.logger.Debug("Device mounted successfully", "device", devicePath, "mountpoint", mountPoint)
	return nil
}

// Private helper methods

// validateInputs validates the input parameters
func (uwl *USBWriterLinux) validateInputs(isoPath string, devicePath string) error {
	// Validate ISO file
	if _, err := os.Stat(isoPath); err != nil {
		return fmt.Errorf("ISO file does not exist: %s", isoPath)
	}

	// Validate device
	if _, err := os.Stat(devicePath); err != nil {
		return fmt.Errorf("device does not exist: %s", devicePath)
	}

	// Check if device is a block device
	cmd := exec.Command("lsblk", "-n", "-o", "TYPE", devicePath)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to check device type: %w", err)
	}

	deviceType := strings.TrimSpace(string(output))
	if deviceType != "disk" {
		return fmt.Errorf("device is not a disk: %s (type: %s)", devicePath, deviceType)
	}

	return nil
}

// writeISOWithDD writes ISO to device using dd command
func (uwl *USBWriterLinux) writeISOWithDD(ctx context.Context, isoPath string, devicePath string, progressTracker *WriteProgressTracker) error {
	uwl.logger.Debug("Writing ISO with dd", "iso", isoPath, "device", devicePath)

	// Create dd command with progress tracking
	cmd := exec.CommandContext(ctx, "dd",
		"if="+isoPath,
		"of="+devicePath,
		"bs=4M",
		"status=progress",
		"conv=fsync")

	// Set up progress monitoring
	go uwl.monitorDDProgress(ctx, devicePath, progressTracker)

	// Execute dd command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("dd command failed: %w", err)
	}

	uwl.logger.Debug("dd command completed successfully")
	return nil
}

// monitorDDProgress monitors the progress of dd command
func (uwl *USBWriterLinux) monitorDDProgress(ctx context.Context, devicePath string, progressTracker *WriteProgressTracker) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-uwl.cancelChan:
			return
		case <-ticker.C:
			// Check device write progress using iostat or similar
			// This is a simplified implementation
			// In production, parse dd status output or use iostat
			progressTracker.UpdateProgress(progressTracker.bytesWritten + 1024*1024) // Simulate progress
		}
	}
}

// syncDevice syncs the device to ensure data is written to disk
func (uwl *USBWriterLinux) syncDevice(devicePath string) error {
	uwl.logger.Debug("Syncing device", "device", devicePath)

	// Use sync command to flush buffers
	cmd := exec.Command("sync")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("sync command failed: %w", err)
	}

	// Also try to sync the specific device
	cmd = exec.Command("sync", devicePath)
	if err := cmd.Run(); err != nil {
		uwl.logger.Warn("Failed to sync specific device", "device", devicePath, "error", err)
	}

	uwl.logger.Debug("Device synced successfully", "device", devicePath)
	return nil
}

// ValidateDevice validates if a device is suitable for writing (override base)
func (uwl *USBWriterLinux) ValidateDevice(ctx context.Context, devicePath string) error {
	uwl.logger.Debug("Validating Linux device", "device", devicePath)

	// Check if device exists
	if _, err := os.Stat(devicePath); err != nil {
		return fmt.Errorf("device does not exist: %s", devicePath)
	}

	// Check if it's a block device (not a partition)
	cmd := exec.CommandContext(ctx, "lsblk", "-n", "-o", "TYPE", devicePath)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to check device type: %w", err)
	}

	deviceType := strings.TrimSpace(string(output))
	if deviceType != "disk" {
		return fmt.Errorf("device must be a disk, not a partition: %s (type: %s)", devicePath, deviceType)
	}

	// Check if we have write permissions
	file, err := os.OpenFile(devicePath, os.O_WRONLY, 0)
	if err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("insufficient permissions to write to device: %s (try running with sudo)", devicePath)
		}
		return fmt.Errorf("device is not writable: %s", devicePath)
	}
	file.Close()

	uwl.logger.Debug("Device validation passed", "device", devicePath)
	return nil
}

// Linux-specific implementation will override the base method

// NewPlatformUSBWriter creates the platform-specific USB writer
func NewPlatformUSBWriter(logger types.Logger) USBWriter {
	return NewUSBWriterLinux(logger)
}
