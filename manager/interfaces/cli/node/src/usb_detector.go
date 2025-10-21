package node

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/node/src/internal/constants"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/node/src/internal/types"
)

// USBDetector defines the interface for USB device detection
type USBDetector interface {
	// DetectDevices detects all available USB devices
	DetectDevices(ctx context.Context) ([]types.USBDevice, error)

	// DetectRemovableDevices detects only removable USB devices
	DetectRemovableDevices(ctx context.Context) ([]types.USBDevice, error)

	// GetDeviceInfo gets detailed information about a specific device
	GetDeviceInfo(ctx context.Context, devicePath string) (*types.USBDevice, error)

	// IsDeviceRemovable checks if a device is removable
	IsDeviceRemovable(ctx context.Context, devicePath string) (bool, error)

	// ValidateDevice validates if a device is suitable for node creation
	ValidateDevice(ctx context.Context, device types.USBDevice) error

	// GetSystemDevices returns all system devices (to avoid writing to them)
	GetSystemDevices(ctx context.Context) ([]string, error)

	// StartMonitoring starts monitoring for USB device changes
	StartMonitoring(ctx context.Context, callback func([]types.USBDevice)) error

	// StopMonitoring stops monitoring for USB device changes
	StopMonitoring() error

	// GetPlatform returns the platform this detector is for
	GetPlatform() string
}

// USBDetectorFactory creates platform-specific USB detectors
type USBDetectorFactory struct{}

// NewUSBDetectorFactory creates a new USB detector factory
func NewUSBDetectorFactory() *USBDetectorFactory {
	return &USBDetectorFactory{}
}

// CreateUSBDetector creates a platform-specific USB detector
func (f *USBDetectorFactory) CreateUSBDetector(logger types.Logger) (USBDetector, error) {
	switch runtime.GOOS {
	case "windows":
		return f.createWindowsDetector(logger), nil
	case "linux":
		return f.createLinuxDetector(logger), nil
	case "darwin":
		return f.createMacOSDetector(logger), nil
	default:
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// createWindowsDetector creates a Windows USB detector
func (f *USBDetectorFactory) createWindowsDetector(logger types.Logger) USBDetector {
	// This will be implemented by platform-specific files
	return nil
}

// createLinuxDetector creates a Linux USB detector
func (f *USBDetectorFactory) createLinuxDetector(logger types.Logger) USBDetector {
	base := NewUSBDetectorBase("linux", logger)
	return &USBDetectorLinux{
		USBDetectorBase: base,
	}
}

// createMacOSDetector creates a macOS USB detector
func (f *USBDetectorFactory) createMacOSDetector(logger types.Logger) USBDetector {
	// This will be implemented by platform-specific files
	return nil
}

// Platform-specific detector constructors will be implemented in separate files
// with build constraints for each platform

// USBDetectorBase provides common functionality for USB detectors
type USBDetectorBase struct {
	logger     types.Logger
	platform   string
	monitoring bool
	stopChan   chan struct{}
}

// NewUSBDetectorBase creates a new base USB detector
func NewUSBDetectorBase(platform string, logger types.Logger) *USBDetectorBase {
	return &USBDetectorBase{
		logger:     logger,
		platform:   platform,
		monitoring: false,
		stopChan:   make(chan struct{}),
	}
}

// GetPlatform returns the platform
func (udb *USBDetectorBase) GetPlatform() string {
	return udb.platform
}

// ValidateDevice validates if a device is suitable for node creation
func (udb *USBDetectorBase) ValidateDevice(ctx context.Context, device types.USBDevice) error {
	// Check minimum capacity
	if device.Capacity < constants.DefaultMinUSBCapacity {
		return fmt.Errorf("device capacity too small: %d bytes (minimum: %d bytes)",
			device.Capacity, constants.DefaultMinUSBCapacity)
	}

	// Check maximum capacity
	if device.Capacity > constants.DefaultMaxUSBCapacity {
		return fmt.Errorf("device capacity too large: %d bytes (maximum: %d bytes)",
			device.Capacity, constants.DefaultMaxUSBCapacity)
	}

	// Check if device is removable
	if !device.IsRemovable {
		return fmt.Errorf("device is not removable: %s", device.Path)
	}

	// Check if device is system device
	if device.IsSystem {
		return fmt.Errorf("device appears to be a system device: %s", device.Path)
	}

	// Check device path
	if device.Path == "" {
		return fmt.Errorf("device path cannot be empty")
	}

	// Check if device has reasonable speed
	if device.Speed < 1 {
		udb.logger.Warn("Device has very low speed", "device", device.Path, "speed", device.Speed)
	}

	udb.logger.Debug("Device validation passed", "device", device.Path, "capacity", device.Capacity)
	return nil
}

// StartMonitoring starts monitoring for USB device changes
func (udb *USBDetectorBase) StartMonitoring(ctx context.Context, callback func([]types.USBDevice)) error {
	if udb.monitoring {
		return fmt.Errorf("monitoring already started")
	}

	udb.monitoring = true
	udb.stopChan = make(chan struct{})

	go udb.monitorDevices(ctx, callback)

	udb.logger.Info("USB device monitoring started")
	return nil
}

// StopMonitoring stops monitoring for USB device changes
func (udb *USBDetectorBase) StopMonitoring() error {
	if !udb.monitoring {
		return fmt.Errorf("monitoring not started")
	}

	close(udb.stopChan)
	udb.monitoring = false

	udb.logger.Info("USB device monitoring stopped")
	return nil
}

// monitorDevices monitors for USB device changes
func (udb *USBDetectorBase) monitorDevices(ctx context.Context, callback func([]types.USBDevice)) {
	ticker := time.NewTicker(constants.DefaultUSBScanInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			udb.logger.Debug("USB monitoring context cancelled")
			return
		case <-udb.stopChan:
			udb.logger.Debug("USB monitoring stopped by user")
			return
		case <-ticker.C:
			// Detect current devices - this will be implemented by the concrete detector
			// For now, we'll skip this functionality in the base class
			continue
		}
	}
}

// devicesEqual compares two device lists for equality
func (udb *USBDetectorBase) devicesEqual(devices1, devices2 []types.USBDevice) bool {
	if len(devices1) != len(devices2) {
		return false
	}

	// Create maps for comparison
	map1 := make(map[string]types.USBDevice)
	map2 := make(map[string]types.USBDevice)

	for _, device := range devices1 {
		map1[device.Path] = device
	}

	for _, device := range devices2 {
		map2[device.Path] = device
	}

	// Compare devices
	for path, device1 := range map1 {
		device2, exists := map2[path]
		if !exists {
			return false
		}

		// Compare key properties
		if device1.Capacity != device2.Capacity ||
			device1.Speed != device2.Speed ||
			device1.IsRemovable != device2.IsRemovable ||
			device1.IsSystem != device2.IsSystem {
			return false
		}
	}

	return true
}

// USBDeviceFilter provides filtering capabilities for USB devices
type USBDeviceFilter struct {
	MinCapacity      int64
	MaxCapacity      int64
	MinSpeed         int
	RequireRemovable bool
	ExcludeSystem    bool
}

// NewUSBDeviceFilter creates a new USB device filter
func NewUSBDeviceFilter() *USBDeviceFilter {
	return &USBDeviceFilter{
		MinCapacity:      constants.DefaultMinUSBCapacity,
		MaxCapacity:      constants.DefaultMaxUSBCapacity,
		MinSpeed:         constants.DefaultPreferredUSBSpeed,
		RequireRemovable: true,
		ExcludeSystem:    true,
	}
}

// FilterDevices filters a list of USB devices based on criteria
func (f *USBDeviceFilter) FilterDevices(devices []types.USBDevice) []types.USBDevice {
	var filtered []types.USBDevice

	for _, device := range devices {
		if f.MatchesFilter(device) {
			filtered = append(filtered, device)
		}
	}

	return filtered
}

// MatchesFilter checks if a device matches the filter criteria
func (f *USBDeviceFilter) MatchesFilter(device types.USBDevice) bool {
	// Check capacity
	if device.Capacity < f.MinCapacity || device.Capacity > f.MaxCapacity {
		return false
	}

	// Check speed
	if device.Speed < f.MinSpeed {
		return false
	}

	// Check removable requirement
	if f.RequireRemovable && !device.IsRemovable {
		return false
	}

	// Check system exclusion
	if f.ExcludeSystem && device.IsSystem {
		return false
	}

	return true
}

// USBDeviceManager manages USB device operations
type USBDeviceManager struct {
	detector USBDetector
	filter   *USBDeviceFilter
	logger   types.Logger
}

// NewUSBDeviceManager creates a new USB device manager
func NewUSBDeviceManager(detector USBDetector, logger types.Logger) *USBDeviceManager {
	return &USBDeviceManager{
		detector: detector,
		filter:   NewUSBDeviceFilter(),
		logger:   logger,
	}
}

// GetSuitableDevices returns devices suitable for node creation
func (udm *USBDeviceManager) GetSuitableDevices(ctx context.Context) ([]types.USBDevice, error) {
	// Detect all removable devices
	devices, err := udm.detector.DetectRemovableDevices(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to detect removable devices: %w", err)
	}

	// Filter devices
	filtered := udm.filter.FilterDevices(devices)

	// Validate each device
	var suitable []types.USBDevice
	for _, device := range filtered {
		if err := udm.detector.ValidateDevice(ctx, device); err != nil {
			udm.logger.Debug("Device validation failed", "device", device.Path, "error", err)
			continue
		}

		suitable = append(suitable, device)
	}

	udm.logger.Info("Found suitable USB devices", "count", len(suitable), "total", len(devices))
	return suitable, nil
}

// GetDeviceDetails returns detailed information about a device
func (udm *USBDeviceManager) GetDeviceDetails(ctx context.Context, devicePath string) (*types.USBDevice, error) {
	device, err := udm.detector.GetDeviceInfo(ctx, devicePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get device info: %w", err)
	}

	// Validate device
	if err := udm.detector.ValidateDevice(ctx, *device); err != nil {
		return nil, fmt.Errorf("device validation failed: %w", err)
	}

	return device, nil
}

// StartDeviceMonitoring starts monitoring for device changes
func (udm *USBDeviceManager) StartDeviceMonitoring(ctx context.Context, callback func([]types.USBDevice)) error {
	return udm.detector.StartMonitoring(ctx, callback)
}

// StopDeviceMonitoring stops monitoring for device changes
func (udm *USBDeviceManager) StopDeviceMonitoring() error {
	return udm.detector.StopMonitoring()
}

// SetFilter updates the device filter
func (udm *USBDeviceManager) SetFilter(filter *USBDeviceFilter) {
	udm.filter = filter
}

// GetFilter returns the current device filter
func (udm *USBDeviceManager) GetFilter() *USBDeviceFilter {
	return udm.filter
}

// GetPlatform returns the platform of the detector
func (udm *USBDeviceManager) GetPlatform() string {
	return udm.detector.GetPlatform()
}
