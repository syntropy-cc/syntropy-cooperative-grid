//go:build darwin
// +build darwin

package node

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"node-component/src/internal/helpers"
	"node-component/src/internal/types"
)

// USBWriterMacOS implements USB writing for macOS systems
type USBWriterMacOS struct {
	*USBWriterBase
}

// NewUSBWriterMacOS creates a new macOS USB writer
func NewUSBWriterMacOS(logger types.Logger) *USBWriterMacOS {
	base := NewUSBWriterBase("darwin", logger)
	return &USBWriterMacOS{
		USBWriterBase: base,
	}
}

// WriteISO writes an ISO to a USB device with cloud-init injection
func (uwm *USBWriterMacOS) WriteISO(ctx context.Context, isoPath string, devicePath string, cloudInitConfig *types.CloudInitConfig) (*types.WriteResult, error) {
	uwm.logger.Info("Writing ISO to USB device", "iso", isoPath, "device", devicePath)

	startTime := time.Now()

	// Validate inputs
	if err := uwm.validateInputs(isoPath, devicePath); err != nil {
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
	progressTracker := NewWriteProgressTracker(uwm.logger, isoInfo.Size())

	// Write ISO to device using dd
	actualISOPath := isoPath

	// Write ISO to device using dd
	if err := uwm.writeISOWithDD(ctx, actualISOPath, devicePath, progressTracker); err != nil {
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
	if err := uwm.syncDevice(devicePath); err != nil {
		uwm.logger.Warn("Failed to sync device", "device", devicePath, "error", err)
	}

	// Create NoCloud partition for cloud-init if provided
	if cloudInitConfig != nil {
		uwm.logger.Info("Creating NoCloud partition for cloud-init")
		cloudInitInjector := NewCloudInitInjector(uwm.logger)
		if err := cloudInitInjector.CreateNoCloudPartition(devicePath, cloudInitConfig); err != nil {
			uwm.logger.Error("Failed to create NoCloud partition", "error", err)
			return nil, fmt.Errorf("failed to create cloud-init partition: %w", err)
		}
	}

	// Clean up injected ISO if created
	if actualISOPath != isoPath {
		os.Remove(actualISOPath)
	}

	duration := time.Since(startTime)
	uwm.logger.Info("ISO write completed successfully",
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
func (uwm *USBWriterMacOS) UnmountDevice(ctx context.Context, devicePath string) error {
	uwm.logger.Debug("Unmounting device", "device", devicePath)

	// Get device name from path (for future use)
	_ = helpers.GetDeviceNameFromPath(devicePath)

	// Check if device has mounted volumes
	cmd := exec.CommandContext(ctx, "diskutil", "list", devicePath)
	output, err := cmd.Output()
	if err != nil {
		uwm.logger.Warn("Failed to check device mounts", "device", devicePath, "error", err)
		return nil // Continue even if check fails
	}

	// Parse output to find mounted volumes
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "disk") && strings.Contains(line, "mounted") {
			// Extract volume identifier
			fields := strings.Fields(line)
			if len(fields) > 0 {
				volumeID := fields[0]
				uwm.logger.Debug("Unmounting volume", "volume", volumeID)

				cmd := exec.CommandContext(ctx, "diskutil", "unmount", volumeID)
				if err := cmd.Run(); err != nil {
					uwm.logger.Warn("Failed to unmount volume", "volume", volumeID, "error", err)
				}
			}
		}
	}

	// Also try to unmount the device itself
	cmd = exec.CommandContext(ctx, "diskutil", "unmount", devicePath)
	if err := cmd.Run(); err != nil {
		uwm.logger.Debug("Device not mounted or already unmounted", "device", devicePath)
	}

	uwm.logger.Debug("Device unmounted successfully", "device", devicePath)
	return nil
}

// MountDevice mounts a device after writing
func (uwm *USBWriterMacOS) MountDevice(ctx context.Context, devicePath string) error {
	uwm.logger.Debug("Mounting device", "device", devicePath)

	// Mount device using diskutil
	cmd := exec.CommandContext(ctx, "diskutil", "mount", devicePath)
	if err := cmd.Run(); err != nil {
		uwm.logger.Warn("Failed to mount device", "device", devicePath, "error", err)
		return nil // Continue even if mount fails
	}

	uwm.logger.Debug("Device mounted successfully", "device", devicePath)
	return nil
}

// Private helper methods

// validateInputs validates the input parameters
func (uwm *USBWriterMacOS) validateInputs(isoPath string, devicePath string) error {
	// Validate ISO file
	if _, err := os.Stat(isoPath); err != nil {
		return fmt.Errorf("ISO file does not exist: %s", isoPath)
	}

	// Validate device
	if _, err := os.Stat(devicePath); err != nil {
		return fmt.Errorf("device does not exist: %s", devicePath)
	}

	// Check if device is a disk using diskutil
	cmd := exec.Command("diskutil", "info", devicePath)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get device info: %w", err)
	}

	// Check if it's a disk device
	if !strings.Contains(string(output), "Type:") || !strings.Contains(string(output), "Disk") {
		return fmt.Errorf("device is not a disk: %s", devicePath)
	}

	return nil
}

// writeISOWithDD writes ISO to device using dd command
func (uwm *USBWriterMacOS) writeISOWithDD(ctx context.Context, isoPath string, devicePath string, progressTracker *WriteProgressTracker) error {
	uwm.logger.Debug("Writing ISO with dd", "iso", isoPath, "device", devicePath)

	// Create dd command with progress tracking
	cmd := exec.CommandContext(ctx, "dd",
		"if="+isoPath,
		"of="+devicePath,
		"bs=4m",
		"status=progress")

	// Set up progress monitoring
	go uwm.monitorDDProgress(ctx, devicePath, progressTracker)

	// Execute dd command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("dd command failed: %w", err)
	}

	uwm.logger.Debug("dd command completed successfully")
	return nil
}

// monitorDDProgress monitors the progress of dd command
func (uwm *USBWriterMacOS) monitorDDProgress(ctx context.Context, devicePath string, progressTracker *WriteProgressTracker) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-uwm.cancelChan:
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
func (uwm *USBWriterMacOS) syncDevice(devicePath string) error {
	uwm.logger.Debug("Syncing device", "device", devicePath)

	// Use sync command to flush buffers
	cmd := exec.Command("sync")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("sync command failed: %w", err)
	}

	uwm.logger.Debug("Device synced successfully", "device", devicePath)
	return nil
}

// macOS-specific implementation will override the base method
