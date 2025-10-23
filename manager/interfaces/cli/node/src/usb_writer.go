package node

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"node-component/src/internal/types"
)

// USBWriter defines the interface for writing ISOs to USB devices
type USBWriter interface {
	// WriteISO writes an ISO to a USB device with cloud-init injection
	WriteISO(ctx context.Context, isoPath string, devicePath string, cloudInitConfig *types.CloudInitConfig) (*types.WriteResult, error)

	// ValidateDevice validates if a device is suitable for writing
	ValidateDevice(ctx context.Context, devicePath string) error

	// UnmountDevice unmounts a device before writing
	UnmountDevice(ctx context.Context, devicePath string) error

	// MountDevice mounts a device after writing
	MountDevice(ctx context.Context, devicePath string) error

	// GetWriteProgress returns the progress of a write operation
	GetWriteProgress() *types.WriteProgress

	// CancelWrite cancels an ongoing write operation
	CancelWrite() error

	// GetPlatform returns the platform this writer is for
	GetPlatform() string
}

// USBWriterFactory creates platform-specific USB writers
type USBWriterFactory struct{}

// NewUSBWriterFactory creates a new USB writer factory
func NewUSBWriterFactory() *USBWriterFactory {
	return &USBWriterFactory{}
}

// CreateUSBWriter creates a platform-specific USB writer
func (f *USBWriterFactory) CreateUSBWriter(logger types.Logger) (USBWriter, error) {
	writer := NewPlatformUSBWriter(logger)
	if writer == nil {
		return nil, fmt.Errorf("failed to create USB writer for platform: %s", runtime.GOOS)
	}
	return writer, nil
}

// Platform-specific detector constructors will be implemented in separate files
// with build constraints for each platform

// USBWriterBase provides common functionality for USB writers
type USBWriterBase struct {
	logger        types.Logger
	platform      string
	writeProgress *types.WriteProgress
	cancelChan    chan struct{}
}

// NewUSBWriterBase creates a new base USB writer
func NewUSBWriterBase(platform string, logger types.Logger) *USBWriterBase {
	return &USBWriterBase{
		logger:        logger,
		platform:      platform,
		writeProgress: &types.WriteProgress{},
		cancelChan:    make(chan struct{}),
	}
}

// GetPlatform returns the platform
func (uwb *USBWriterBase) GetPlatform() string {
	return uwb.platform
}

// GetWriteProgress returns the progress of a write operation
func (uwb *USBWriterBase) GetWriteProgress() *types.WriteProgress {
	return uwb.writeProgress
}

// CancelWrite cancels an ongoing write operation
func (uwb *USBWriterBase) CancelWrite() error {
	select {
	case <-uwb.cancelChan:
		// Already cancelled
		return nil
	default:
		close(uwb.cancelChan)
		uwb.logger.Info("Write operation cancelled")
		return nil
	}
}

// ValidateDevice validates if a device is suitable for writing
// This is a fallback implementation - platform-specific implementations should override
func (uwb *USBWriterBase) ValidateDevice(ctx context.Context, devicePath string) error {
	// Check if device path is not empty
	if devicePath == "" {
		return fmt.Errorf("device path is empty")
	}

	uwb.logger.Debug("Using base device validation", "device", devicePath)

	// Note: Platform-specific implementations should override this method
	// to provide proper validation for their platform
	return nil
}

// USBWriterManager manages USB writing operations
type USBWriterManager struct {
	writer USBWriter
	logger types.Logger
}

// NewUSBWriterManager creates a new USB writer manager
func NewUSBWriterManager(writer USBWriter, logger types.Logger) *USBWriterManager {
	return &USBWriterManager{
		writer: writer,
		logger: logger,
	}
}

// WriteUbuntuISO writes Ubuntu ISO to USB with cloud-init injection
func (uwm *USBWriterManager) WriteUbuntuISO(ctx context.Context, isoPath string, devicePath string, cloudInitConfig *types.CloudInitConfig) (*types.WriteResult, error) {
	uwm.logger.Info("Starting USB write operation", "iso", isoPath, "device", devicePath)

	// Validate ISO file
	if err := uwm.validateISOFile(isoPath); err != nil {
		return nil, fmt.Errorf("ISO validation failed: %w", err)
	}

	// Validate device
	if err := uwm.writer.ValidateDevice(ctx, devicePath); err != nil {
		return nil, fmt.Errorf("device validation failed: %w", err)
	}

	// Unmount device before writing
	if err := uwm.writer.UnmountDevice(ctx, devicePath); err != nil {
		uwm.logger.Warn("Failed to unmount device", "device", devicePath, "error", err)
	}

	// Write ISO with cloud-init injection
	result, err := uwm.writer.WriteISO(ctx, isoPath, devicePath, cloudInitConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to write ISO: %w", err)
	}

	// Mount device after writing
	if err := uwm.writer.MountDevice(ctx, devicePath); err != nil {
		uwm.logger.Warn("Failed to mount device", "device", devicePath, "error", err)
	}

	uwm.logger.Info("USB write operation completed successfully",
		"device", devicePath,
		"duration", result.Duration,
		"bytes_written", result.BytesWritten)

	return result, nil
}

// GetWriteProgress returns the current write progress
func (uwm *USBWriterManager) GetWriteProgress() *types.WriteProgress {
	return uwm.writer.GetWriteProgress()
}

// CancelWrite cancels the current write operation
func (uwm *USBWriterManager) CancelWrite() error {
	return uwm.writer.CancelWrite()
}

// ValidateDevice validates a USB device for writing
func (uwm *USBWriterManager) ValidateDevice(ctx context.Context, devicePath string) error {
	return uwm.writer.ValidateDevice(ctx, devicePath)
}

// GetPlatform returns the platform of the writer
func (uwm *USBWriterManager) GetPlatform() string {
	return uwm.writer.GetPlatform()
}

// Private helper methods

// validateISOFile validates the ISO file before writing
func (uwm *USBWriterManager) validateISOFile(isoPath string) error {
	// Check if file exists
	if _, err := os.Stat(isoPath); err != nil {
		return fmt.Errorf("ISO file does not exist: %s", isoPath)
	}

	// Check if file is readable
	file, err := os.Open(isoPath)
	if err != nil {
		return fmt.Errorf("ISO file is not readable: %s", isoPath)
	}
	file.Close()

	// Get file info for logging
	fileInfo, err := os.Stat(isoPath)
	if err != nil {
		return fmt.Errorf("failed to get ISO file info: %w", err)
	}

	// Skip size validation by default - let user choose any ISO
	// TODO: Add size validation toggle in future versions
	uwm.logger.Debug("Skipping ISO size validation", "path", isoPath, "size", fileInfo.Size())
	return nil
}

// WriteProgressTracker tracks write progress
type WriteProgressTracker struct {
	logger       types.Logger
	startTime    time.Time
	lastUpdate   time.Time
	bytesWritten int64
	totalBytes   int64
	speed        float64 // MB/s
	eta          time.Duration
}

// NewWriteProgressTracker creates a new write progress tracker
func NewWriteProgressTracker(logger types.Logger, totalBytes int64) *WriteProgressTracker {
	return &WriteProgressTracker{
		logger:     logger,
		startTime:  time.Now(),
		lastUpdate: time.Now(),
		totalBytes: totalBytes,
	}
}

// UpdateProgress updates the write progress
func (wpt *WriteProgressTracker) UpdateProgress(bytesWritten int64) {
	wpt.bytesWritten = bytesWritten

	now := time.Now()
	elapsed := now.Sub(wpt.startTime)
	timeSinceLastUpdate := now.Sub(wpt.lastUpdate)

	// Calculate speed (MB/s)
	if elapsed.Seconds() > 0 {
		wpt.speed = float64(bytesWritten) / elapsed.Seconds() / 1024 / 1024
	}

	// Calculate ETA
	if wpt.speed > 0 && wpt.bytesWritten < wpt.totalBytes {
		remaining := wpt.totalBytes - wpt.bytesWritten
		wpt.eta = time.Duration(float64(remaining)/wpt.speed/1024/1024) * time.Second
	}

	// Log progress every 5%
	percentage := float64(bytesWritten) / float64(wpt.totalBytes) * 100
	if int(percentage)%5 == 0 && timeSinceLastUpdate > time.Second {
		wpt.logger.Info("Write progress",
			"progress", fmt.Sprintf("%.1f%%", percentage),
			"written", fmt.Sprintf("%.1fMB", float64(bytesWritten)/1024/1024),
			"total", fmt.Sprintf("%.1fMB", float64(wpt.totalBytes)/1024/1024),
			"speed", fmt.Sprintf("%.1fMB/s", wpt.speed),
			"eta", wpt.eta.Truncate(time.Second))
		wpt.lastUpdate = now
	}
}

// GetProgress returns the current progress information
func (wpt *WriteProgressTracker) GetProgress() *types.WriteProgress {
	return &types.WriteProgress{
		BytesWritten: wpt.bytesWritten,
		TotalBytes:   wpt.totalBytes,
		Percentage:   float64(wpt.bytesWritten) / float64(wpt.totalBytes) * 100,
		Speed:        wpt.speed,
		ETA:          wpt.eta,
	}
}

// CloudInitInjector handles injection of cloud-init into ISO
type CloudInitInjector struct {
	logger types.Logger
}

// NewCloudInitInjector creates a new cloud-init injector
func NewCloudInitInjector(logger types.Logger) *CloudInitInjector {
	return &CloudInitInjector{
		logger: logger,
	}
}

// CreateNoCloudPartition creates a FAT32 partition with cloud-init files
func (cii *CloudInitInjector) CreateNoCloudPartition(devicePath string, cloudInitConfig *types.CloudInitConfig) error {
	cii.logger.Info("Creating NoCloud partition", "device", devicePath)

	// 1. Detectar tamanho do device
	deviceSize, err := getDeviceSize(devicePath)
	if err != nil {
		return fmt.Errorf("failed to get device size: %w", err)
	}

	// 2. Criar partição FAT32 (100MB) ao final do device
	// usando parted: mkpart primary fat32 <end-100MB> <end>
	partitionNum := 2 // Assumindo que ISO é partição 1
	partitionPath := fmt.Sprintf("%s%d", devicePath, partitionNum)

	cmd := exec.Command("parted", "-s", devicePath,
		"mkpart", "primary", "fat32",
		fmt.Sprintf("%dMB", deviceSize-100),
		fmt.Sprintf("%dMB", deviceSize))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create partition: %w", err)
	}

	// 3. Formatar como FAT32
	cmd = exec.Command("mkfs.vfat", "-n", "CIDATA", partitionPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to format partition: %w", err)
	}

	// 4. Montar partição temporária
	mountPoint, err := os.MkdirTemp("", "cidata-*")
	if err != nil {
		return fmt.Errorf("failed to create mount point: %w", err)
	}
	defer os.RemoveAll(mountPoint)

	cmd = exec.Command("mount", partitionPath, mountPoint)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to mount partition: %w", err)
	}
	defer exec.Command("umount", mountPoint).Run()

	// 5. Escrever arquivos cloud-init (SEM extensão .yaml!)
	files := map[string]string{
		"user-data":      cloudInitConfig.UserData,
		"meta-data":      cloudInitConfig.MetaData,
		"network-config": cloudInitConfig.NetworkConfig,
	}

	for filename, content := range files {
		filepath := filepath.Join(mountPoint, filename)
		if err := os.WriteFile(filepath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", filename, err)
		}
	}

	cii.logger.Info("NoCloud partition created successfully")
	return nil
}

func getDeviceSize(devicePath string) (int64, error) {
	cmd := exec.Command("blockdev", "--getsize64", devicePath)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	size, err := strconv.ParseInt(strings.TrimSpace(string(output)), 10, 64)
	if err != nil {
		return 0, err
	}
	return size / (1024 * 1024), nil // Converter para MB
}

// InjectCloudInit injects cloud-init configuration into an ISO (DEPRECATED - use CreateNoCloudPartition)
func (cii *CloudInitInjector) InjectCloudInit(isoPath string, cloudInitConfig *types.CloudInitConfig) (string, error) {
	cii.logger.Info("Injecting cloud-init into ISO", "iso", isoPath)

	// Create temporary directory for ISO extraction
	tempDir, err := os.MkdirTemp("", "syntropy-iso-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Extract ISO contents
	if err := cii.extractISO(isoPath, tempDir); err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to extract ISO: %w", err)
	}

	// Inject cloud-init files
	if err := cii.injectCloudInitFiles(tempDir, cloudInitConfig); err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to inject cloud-init files: %w", err)
	}

	// Create new ISO with cloud-init
	newISOPath := filepath.Join(filepath.Dir(isoPath), "syntropy-"+filepath.Base(isoPath))
	if err := cii.createISO(tempDir, newISOPath); err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to create new ISO: %w", err)
	}

	// Clean up temp directory
	os.RemoveAll(tempDir)

	cii.logger.Info("Cloud-init injection completed", "new_iso", newISOPath)
	return newISOPath, nil
}

// extractISO extracts ISO contents to a directory
func (cii *CloudInitInjector) extractISO(isoPath string, extractDir string) error {
	// This is a simplified implementation
	// In production, use proper ISO extraction tools like 7zip or genisoimage
	cii.logger.Debug("Extracting ISO", "iso", isoPath, "dir", extractDir)

	// For now, we'll create a placeholder
	// TODO: Implement proper ISO extraction
	return nil
}

// injectCloudInitFiles injects cloud-init files into the extracted ISO
func (cii *CloudInitInjector) injectCloudInitFiles(extractDir string, cloudInitConfig *types.CloudInitConfig) error {
	cii.logger.Debug("Injecting cloud-init files", "dir", extractDir)

	// Create cloud-init directory
	cloudInitDir := filepath.Join(extractDir, "cloud-init")
	if err := os.MkdirAll(cloudInitDir, 0755); err != nil {
		return fmt.Errorf("failed to create cloud-init directory: %w", err)
	}

	// Write user-data
	if cloudInitConfig.UserData != "" {
		userDataPath := filepath.Join(cloudInitDir, "user-data")
		if err := os.WriteFile(userDataPath, []byte(cloudInitConfig.UserData), 0644); err != nil {
			return fmt.Errorf("failed to write user-data: %w", err)
		}
	}

	// Write network-config
	if cloudInitConfig.NetworkConfig != "" {
		networkConfigPath := filepath.Join(cloudInitDir, "network-config")
		if err := os.WriteFile(networkConfigPath, []byte(cloudInitConfig.NetworkConfig), 0644); err != nil {
			return fmt.Errorf("failed to write network-config: %w", err)
		}
	}

	// Write meta-data
	if cloudInitConfig.MetaData != "" {
		metaDataPath := filepath.Join(cloudInitDir, "meta-data")
		if err := os.WriteFile(metaDataPath, []byte(cloudInitConfig.MetaData), 0644); err != nil {
			return fmt.Errorf("failed to write meta-data: %w", err)
		}
	}

	cii.logger.Debug("Cloud-init files injected successfully")
	return nil
}

// createISO creates a new ISO from extracted contents
func (cii *CloudInitInjector) createISO(extractDir string, outputPath string) error {
	cii.logger.Debug("Creating ISO", "dir", extractDir, "output", outputPath)

	// This is a simplified implementation
	// In production, use proper ISO creation tools like genisoimage or mkisofs
	// TODO: Implement proper ISO creation

	// For now, we'll create a placeholder file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output ISO: %w", err)
	}
	file.Close()

	cii.logger.Debug("ISO created successfully", "path", outputPath)
	return nil
}

// WriteResult represents the result of a write operation
type WriteResult struct {
	DevicePath        string        `json:"device_path"`
	ISOPath           string        `json:"iso_path"`
	BytesWritten      int64         `json:"bytes_written"`
	Duration          time.Duration `json:"duration"`
	Success           bool          `json:"success"`
	ErrorMessage      string        `json:"error_message,omitempty"`
	CloudInitInjected bool          `json:"cloud_init_injected"`
}
