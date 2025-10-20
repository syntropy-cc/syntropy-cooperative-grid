package mocks

import (
	"fmt"
	"sync"

	node "github.com/syntropy-grid/syntropy-cooperative-grid/manager/interfaces/cli/node/src"
)

// MockUSBDevice simula dispositivos USB para testes
type MockUSBDevice struct {
	mu sync.RWMutex

	// Configuração de comportamento
	Devices           []node.USBDevice
	DeviceInfo        map[string]*node.USBDevice
	ValidationResults map[string]error
	DetectionError    error

	// Histórico de chamadas
	GetSuitableDevicesCalls     []GetSuitableDevicesCall
	GetDeviceInfoCalls          []GetDeviceInfoCall
	ValidateDeviceCalls         []ValidateDeviceCall
	GetSuitableDevicesCallCount int
	GetDeviceInfoCallCount      int
	ValidateDeviceCallCount     int
}

// GetSuitableDevicesCall representa uma chamada para GetSuitableDevices
type GetSuitableDevicesCall struct {
	Timestamp int64
	Devices   []node.USBDevice
	Error     error
}

// GetDeviceInfoCall representa uma chamada para GetDeviceInfo
type GetDeviceInfoCall struct {
	Timestamp int64
	Path      string
	Device    *node.USBDevice
	Error     error
}

// ValidateDeviceCall representa uma chamada para ValidateDevice
type ValidateDeviceCall struct {
	Timestamp int64
	Device    *node.USBDevice
	Error     error
}

// NewMockUSBDevice cria um novo mock do USBDevice
func NewMockUSBDevice() *MockUSBDevice {
	return &MockUSBDevice{
		Devices:                 make([]node.USBDevice, 0),
		DeviceInfo:              make(map[string]*node.USBDevice),
		ValidationResults:       make(map[string]error),
		GetSuitableDevicesCalls: make([]GetSuitableDevicesCall, 0),
		GetDeviceInfoCalls:      make([]GetDeviceInfoCall, 0),
		ValidateDeviceCalls:     make([]ValidateDeviceCall, 0),
	}
}

// GetSuitableDevices simula a detecção de dispositivos USB adequados
func (m *MockUSBDevice) GetSuitableDevices() ([]node.USBDevice, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.GetSuitableDevicesCallCount++
	call := GetSuitableDevicesCall{
		Timestamp: getCurrentTimestamp(),
		Devices:   m.Devices,
		Error:     m.DetectionError,
	}
	m.GetSuitableDevicesCalls = append(m.GetSuitableDevicesCalls, call)

	if m.DetectionError != nil {
		return nil, m.DetectionError
	}

	// Retorna uma cópia dos dispositivos
	devices := make([]node.USBDevice, len(m.Devices))
	copy(devices, m.Devices)
	return devices, nil
}

// GetDeviceInfo simula a obtenção de informações de um dispositivo específico
func (m *MockUSBDevice) GetDeviceInfo(path string) (*node.USBDevice, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.GetDeviceInfoCallCount++
	call := GetDeviceInfoCall{
		Timestamp: getCurrentTimestamp(),
		Path:      path,
		Device:    nil,
		Error:     nil,
	}

	device, exists := m.DeviceInfo[path]
	if !exists {
		call.Error = fmt.Errorf("device not found: %s", path)
		m.GetDeviceInfoCalls = append(m.GetDeviceInfoCalls, call)
		return nil, call.Error
	}

	call.Device = device
	m.GetDeviceInfoCalls = append(m.GetDeviceInfoCalls, call)
	return device, nil
}

// ValidateDevice simula a validação de um dispositivo
func (m *MockUSBDevice) ValidateDevice(device *node.USBDevice) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.ValidateDeviceCallCount++
	call := ValidateDeviceCall{
		Timestamp: getCurrentTimestamp(),
		Device:    device,
		Error:     nil,
	}

	if device == nil {
		call.Error = fmt.Errorf("device cannot be nil")
		m.ValidateDeviceCalls = append(m.ValidateDeviceCalls, call)
		return call.Error
	}

	// Verifica se há um resultado de validação específico para este dispositivo
	if err, exists := m.ValidationResults[device.Path]; exists {
		call.Error = err
		m.ValidateDeviceCalls = append(m.ValidateDeviceCalls, call)
		return err
	}

	// Validação básica
	if device.Path == "" {
		call.Error = fmt.Errorf("device path cannot be empty")
		m.ValidateDeviceCalls = append(m.ValidateDeviceCalls, call)
		return call.Error
	}

	if device.Capacity < 8000000000 { // 8GB em bytes
		call.Error = fmt.Errorf("device capacity too small: %d bytes", device.Capacity)
		m.ValidateDeviceCalls = append(m.ValidateDeviceCalls, call)
		return call.Error
	}

	m.ValidateDeviceCalls = append(m.ValidateDeviceCalls, call)
	return nil
}

// Configuração de comportamento do mock

// SetDevices define os dispositivos que serão retornados
func (m *MockUSBDevice) SetDevices(devices []node.USBDevice) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Devices = devices
}

// AddDevice adiciona um dispositivo à lista
func (m *MockUSBDevice) AddDevice(device node.USBDevice) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Devices = append(m.Devices, device)
	m.DeviceInfo[device.Path] = &device
}

// SetDeviceInfo define as informações de um dispositivo específico
func (m *MockUSBDevice) SetDeviceInfo(path string, device *node.USBDevice) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.DeviceInfo[path] = device
}

// SetValidationResult define o resultado de validação para um dispositivo
func (m *MockUSBDevice) SetValidationResult(path string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ValidationResults[path] = err
}

// SetDetectionError define o erro que será retornado em GetSuitableDevices
func (m *MockUSBDevice) SetDetectionError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.DetectionError = err
}

// Verificação de chamadas

// GetGetSuitableDevicesCalls retorna o histórico de chamadas para GetSuitableDevices
func (m *MockUSBDevice) GetGetSuitableDevicesCalls() []GetSuitableDevicesCall {
	m.mu.RLock()
	defer m.mu.RUnlock()

	calls := make([]GetSuitableDevicesCall, len(m.GetSuitableDevicesCalls))
	copy(calls, m.GetSuitableDevicesCalls)
	return calls
}

// GetGetDeviceInfoCalls retorna o histórico de chamadas para GetDeviceInfo
func (m *MockUSBDevice) GetGetDeviceInfoCalls() []GetDeviceInfoCall {
	m.mu.RLock()
	defer m.mu.RUnlock()

	calls := make([]GetDeviceInfoCall, len(m.GetDeviceInfoCalls))
	copy(calls, m.GetDeviceInfoCalls)
	return calls
}

// GetValidateDeviceCalls retorna o histórico de chamadas para ValidateDevice
func (m *MockUSBDevice) GetValidateDeviceCalls() []ValidateDeviceCall {
	m.mu.RLock()
	defer m.mu.RUnlock()

	calls := make([]ValidateDeviceCall, len(m.ValidateDeviceCalls))
	copy(calls, m.ValidateDeviceCalls)
	return calls
}

// GetCallCounts retorna os contadores de chamadas
func (m *MockUSBDevice) GetCallCounts() (getSuitable, getInfo, validate int) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.GetSuitableDevicesCallCount, m.GetDeviceInfoCallCount, m.ValidateDeviceCallCount
}

// Reset reseta o mock para o estado inicial
func (m *MockUSBDevice) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Devices = make([]node.USBDevice, 0)
	m.DeviceInfo = make(map[string]*node.USBDevice)
	m.ValidationResults = make(map[string]error)
	m.DetectionError = nil
	m.GetSuitableDevicesCalls = make([]GetSuitableDevicesCall, 0)
	m.GetDeviceInfoCalls = make([]GetDeviceInfoCall, 0)
	m.ValidateDeviceCalls = make([]ValidateDeviceCall, 0)
	m.GetSuitableDevicesCallCount = 0
	m.GetDeviceInfoCallCount = 0
	m.ValidateDeviceCallCount = 0
}

// Verificação de comportamento

// WasGetSuitableDevicesCalled verifica se GetSuitableDevices foi chamado
func (m *MockUSBDevice) WasGetSuitableDevicesCalled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.GetSuitableDevicesCallCount > 0
}

// WasGetDeviceInfoCalled verifica se GetDeviceInfo foi chamado
func (m *MockUSBDevice) WasGetDeviceInfoCalled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.GetDeviceInfoCallCount > 0
}

// WasValidateDeviceCalled verifica se ValidateDevice foi chamado
func (m *MockUSBDevice) WasValidateDeviceCalled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.ValidateDeviceCallCount > 0
}

// WasGetDeviceInfoCalledWith verifica se GetDeviceInfo foi chamado com um path específico
func (m *MockUSBDevice) WasGetDeviceInfoCalledWith(path string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, call := range m.GetDeviceInfoCalls {
		if call.Path == path {
			return true
		}
	}
	return false
}

// WasValidateDeviceCalledWith verifica se ValidateDevice foi chamado com um dispositivo específico
func (m *MockUSBDevice) WasValidateDeviceCalledWith(device *node.USBDevice) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, call := range m.ValidateDeviceCalls {
		if call.Device != nil && device != nil && call.Device.Path == device.Path {
			return true
		}
	}
	return false
}

// GetLastGetSuitableDevicesCall retorna a última chamada para GetSuitableDevices
func (m *MockUSBDevice) GetLastGetSuitableDevicesCall() *GetSuitableDevicesCall {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.GetSuitableDevicesCalls) == 0 {
		return nil
	}

	lastCall := m.GetSuitableDevicesCalls[len(m.GetSuitableDevicesCalls)-1]
	return &lastCall
}

// GetLastGetDeviceInfoCall retorna a última chamada para GetDeviceInfo
func (m *MockUSBDevice) GetLastGetDeviceInfoCall() *GetDeviceInfoCall {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.GetDeviceInfoCalls) == 0 {
		return nil
	}

	lastCall := m.GetDeviceInfoCalls[len(m.GetDeviceInfoCalls)-1]
	return &lastCall
}

// GetLastValidateDeviceCall retorna a última chamada para ValidateDevice
func (m *MockUSBDevice) GetLastValidateDeviceCall() *ValidateDeviceCall {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.ValidateDeviceCalls) == 0 {
		return nil
	}

	lastCall := m.ValidateDeviceCalls[len(m.ValidateDeviceCalls)-1]
	return &lastCall
}

// MockUSBDeviceFactory cria uma factory para mocks do USBDevice
type MockUSBDeviceFactory struct {
	defaultDevices []node.USBDevice
	defaultError   error
}

// NewMockUSBDeviceFactory cria uma nova factory
func NewMockUSBDeviceFactory() *MockUSBDeviceFactory {
	return &MockUSBDeviceFactory{
		defaultDevices: []node.USBDevice{
			{
				Path:        "/dev/sdb",
				RawPath:     "/dev/sdb",
				Capacity:    16000000000, // 16GB
				Speed:       3,
				Vendor:      "ACME",
				Model:       "USB-DRIVE",
				Serial:      "XYZ123",
				IsSystem:    false,
				IsRemovable: true,
			},
		},
		defaultError: nil,
	}
}

// SetDefaultDevices define os dispositivos padrão
func (f *MockUSBDeviceFactory) SetDefaultDevices(devices []node.USBDevice) {
	f.defaultDevices = devices
}

// SetDefaultError define o erro padrão
func (f *MockUSBDeviceFactory) SetDefaultError(err error) {
	f.defaultError = err
}

// Create cria um novo mock com as configurações padrão
func (f *MockUSBDeviceFactory) Create() *MockUSBDevice {
	mock := NewMockUSBDevice()
	mock.SetDevices(f.defaultDevices)
	mock.SetDetectionError(f.defaultError)
	return mock
}

// CreateWithError cria um mock que retorna erro em GetSuitableDevices
func (f *MockUSBDeviceFactory) CreateWithError(err error) *MockUSBDevice {
	mock := f.Create()
	mock.SetDetectionError(err)
	return mock
}

// CreateWithNoDevices cria um mock que não retorna dispositivos
func (f *MockUSBDeviceFactory) CreateWithNoDevices() *MockUSBDevice {
	mock := f.Create()
	mock.SetDevices([]node.USBDevice{})
	return mock
}

// CreateWithMultipleDevices cria um mock com múltiplos dispositivos
func (f *MockUSBDeviceFactory) CreateWithMultipleDevices() *MockUSBDevice {
	devices := []node.USBDevice{
		{
			Path:        "/dev/sdb",
			RawPath:     "/dev/sdb",
			Capacity:    16000000000, // 16GB
			Speed:       3,
			Vendor:      "ACME",
			Model:       "USB-DRIVE-1",
			Serial:      "XYZ111",
			IsSystem:    false,
			IsRemovable: true,
		},
		{
			Path:        "/dev/sdc",
			RawPath:     "/dev/sdc",
			Capacity:    32000000000, // 32GB
			Speed:       3,
			Vendor:      "ACME",
			Model:       "USB-DRIVE-2",
			Serial:      "XYZ222",
			IsSystem:    false,
			IsRemovable: true,
		},
		{
			Path:        "/dev/sdd",
			RawPath:     "/dev/sdd",
			Capacity:    64000000000, // 64GB
			Speed:       3,
			Vendor:      "ACME",
			Model:       "USB-DRIVE-3",
			Serial:      "XYZ333",
			IsSystem:    false,
			IsRemovable: true,
		},
	}

	mock := f.Create()
	mock.SetDevices(devices)
	return mock
}

// CreateWithInvalidDevice cria um mock com dispositivo inválido
func (f *MockUSBDeviceFactory) CreateWithInvalidDevice() *MockUSBDevice {
	devices := []node.USBDevice{
		{
			Path:        "/dev/sdb",
			RawPath:     "/dev/sdb",
			Capacity:    4000000000, // 4GB (muito pequeno)
			Speed:       2,
			Vendor:      "ACME",
			Model:       "USB-DRIVE-SMALL",
			Serial:      "XYZ444",
			IsSystem:    false,
			IsRemovable: true,
		},
	}

	mock := f.Create()
	mock.SetDevices(devices)
	return mock
}

// CreateWithSystemDevice cria um mock com dispositivo do sistema
func (f *MockUSBDeviceFactory) CreateWithSystemDevice() *MockUSBDevice {
	devices := []node.USBDevice{
		{
			Path:        "/dev/sda",
			RawPath:     "/dev/sda",
			Capacity:    1000000000000, // 1TB
			Speed:       0,
			Vendor:      "ACME",
			Model:       "SYSTEM-DISK",
			Serial:      "XYZ999",
			IsSystem:    true,
			IsRemovable: false,
		},
	}

	mock := f.Create()
	mock.SetDevices(devices)
	return mock
}
