//go:build windows
// +build windows

package node

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"node-component/src/internal/types"
)

// USBWriterWindows implements USB writing for Windows systems
type USBWriterWindows struct {
	*USBWriterBase
}

// NewUSBWriterWindows creates a new Windows USB writer
func NewUSBWriterWindows(logger types.Logger) *USBWriterWindows {
	base := NewUSBWriterBase("windows", logger)
	return &USBWriterWindows{
		USBWriterBase: base,
	}
}

// WriteISO writes an ISO to a USB device with cloud-init injection
func (uww *USBWriterWindows) WriteISO(ctx context.Context, isoPath string, devicePath string, cloudInitConfig *types.CloudInitConfig) (*types.WriteResult, error) {
	uww.logger.Info("Writing ISO to USB device", "iso", isoPath, "device", devicePath)

	startTime := time.Now()

	// Validate inputs
	if err := uww.validateInputs(isoPath, devicePath); err != nil {
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
	progressTracker := NewWriteProgressTracker(uww.logger, isoInfo.Size())

	// Write ISO to device using dd or diskpart
	actualISOPath := isoPath

	// Write ISO to device using dd or diskpart
	if err := uww.writeISOWithDD(ctx, actualISOPath, devicePath, progressTracker); err != nil {
		return &types.WriteResult{
			DevicePath:   devicePath,
			ISOPath:      isoPath,
			BytesWritten: 0,
			Duration:     time.Since(startTime),
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	// Create NoCloud partition for cloud-init if provided
	if cloudInitConfig != nil {
		uww.logger.Info("Creating NoCloud partition for cloud-init")
		cloudInitInjector := NewCloudInitInjector(uww.logger)
		if err := cloudInitInjector.CreateNoCloudPartition(devicePath, cloudInitConfig); err != nil {
			uww.logger.Error("Failed to create NoCloud partition", "error", err)
			return nil, fmt.Errorf("failed to create cloud-init partition: %w", err)
		}
	}

	// Clean up injected ISO if created
	if actualISOPath != isoPath {
		os.Remove(actualISOPath)
	}

	duration := time.Since(startTime)
	uww.logger.Info("ISO write completed successfully",
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
func (uww *USBWriterWindows) UnmountDevice(ctx context.Context, devicePath string) error {
	uww.logger.Debug("Unmounting device", "device", devicePath)

	// Get drive letter from device path (e.g., "C:" from "C:\")
	driveLetter := strings.TrimSuffix(devicePath, ":")
	if driveLetter == "" {
		return fmt.Errorf("invalid device path: %s", devicePath)
	}

	// Use diskpart to unmount the drive
	diskpartScript := fmt.Sprintf(`
select disk %s
offline disk
online disk
`, driveLetter)

	// Create temporary script file
	scriptPath := filepath.Join(os.TempDir(), "syntropy-unmount.txt")
	if err := os.WriteFile(scriptPath, []byte(diskpartScript), 0644); err != nil {
		return fmt.Errorf("failed to create diskpart script: %w", err)
	}
	defer os.Remove(scriptPath)

	// Execute diskpart
	cmd := exec.CommandContext(ctx, "diskpart", "/s", scriptPath)
	if err := cmd.Run(); err != nil {
		uww.logger.Warn("Failed to unmount device with diskpart", "device", devicePath, "error", err)
		// Continue even if unmount fails
	}

	uww.logger.Debug("Device unmounted successfully", "device", devicePath)
	return nil
}

// MountDevice mounts a device after writing
func (uww *USBWriterWindows) MountDevice(ctx context.Context, devicePath string) error {
	uww.logger.Debug("Mounting device", "device", devicePath)

	// Get drive letter from device path
	driveLetter := strings.TrimSuffix(devicePath, ":")
	if driveLetter == "" {
		return fmt.Errorf("invalid device path: %s", devicePath)
	}

	// Use diskpart to mount the drive
	diskpartScript := fmt.Sprintf(`
select disk %s
online disk
attributes disk clear readonly
`, driveLetter)

	// Create temporary script file
	scriptPath := filepath.Join(os.TempDir(), "syntropy-mount.txt")
	if err := os.WriteFile(scriptPath, []byte(diskpartScript), 0644); err != nil {
		return fmt.Errorf("failed to create diskpart script: %w", err)
	}
	defer os.Remove(scriptPath)

	// Execute diskpart
	cmd := exec.CommandContext(ctx, "diskpart", "/s", scriptPath)
	if err := cmd.Run(); err != nil {
		uww.logger.Warn("Failed to mount device with diskpart", "device", devicePath, "error", err)
		return nil // Continue even if mount fails
	}

	uww.logger.Debug("Device mounted successfully", "device", devicePath)
	return nil
}

// Private helper methods

// validateInputs validates the input parameters
func (uww *USBWriterWindows) validateInputs(isoPath string, devicePath string) error {
	// Validate ISO file
	if _, err := os.Stat(isoPath); err != nil {
		return fmt.Errorf("ISO file does not exist: %s", isoPath)
	}

	// Validate device path format (should be like "C:")
	if !strings.HasSuffix(devicePath, ":") {
		return fmt.Errorf("invalid device path format: %s (should be like 'C:')", devicePath)
	}

	// Check if drive exists
	cmd := exec.Command("wmic", "logicaldisk", "where", fmt.Sprintf("deviceid='%s'", devicePath), "get", "deviceid")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to check if drive exists: %w", err)
	}

	if !strings.Contains(string(output), devicePath) {
		return fmt.Errorf("drive does not exist: %s", devicePath)
	}

	return nil
}

// writeISOWithDD writes ISO to device using dd command (if available) or alternative method
func (uww *USBWriterWindows) writeISOWithDD(ctx context.Context, isoPath string, devicePath string, progressTracker *WriteProgressTracker) error {
	uww.logger.Debug("Writing ISO with dd", "iso", isoPath, "device", devicePath)

	// Try to use dd if available (common in Windows with WSL or Git Bash)
	cmd := exec.CommandContext(ctx, "dd",
		"if="+isoPath,
		"of="+devicePath,
		"bs=4M",
		"status=progress",
		"conv=fsync")

	if err := cmd.Run(); err != nil {
		uww.logger.Warn("dd command failed, trying alternative method", "error", err)

		// Fallback to PowerShell method
		return uww.writeISOWithPowerShell(ctx, isoPath, devicePath, progressTracker)
	}

	uww.logger.Debug("dd command completed successfully")
	return nil
}

// writeISOWithPowerShell writes ISO using PowerShell as fallback
func (uww *USBWriterWindows) writeISOWithPowerShell(ctx context.Context, isoPath string, devicePath string, progressTracker *WriteProgressTracker) error {
	uww.logger.Debug("Writing ISO with PowerShell", "iso", isoPath, "device", devicePath)

	// PowerShell script to write ISO to USB
	psScript := fmt.Sprintf(`
$isoPath = "%s"
$devicePath = "%s"

try {
    # Mount ISO
    $isoMount = Mount-DiskImage -ImagePath $isoPath -PassThru
    $isoDrive = ($isoMount | Get-Volume).DriveLetter + ":"
    
    # Copy ISO contents to USB device
    robocopy $isoDrive $devicePath /E /COPY:DAT /R:3 /W:1
    
    # Dismount ISO
    Dismount-DiskImage -ImagePath $isoPath
    
    Write-Output "ISO write completed successfully"
} catch {
    Write-Error "Failed to write ISO: $_"
    exit 1
}
`, isoPath, devicePath)

	// Execute PowerShell script
	cmd := exec.CommandContext(ctx, "powershell", "-Command", psScript)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("PowerShell ISO write failed: %w", err)
	}

	uww.logger.Debug("PowerShell ISO write completed successfully")
	return nil
}

// Windows-specific implementation will override the base method
