package setup

import (
	"runtime"
	"unsafe"
)

// SecureBytes protects sensitive data in memory
type SecureBytes struct {
	data []byte
}

// NewSecureBytes creates a new secure bytes container
func NewSecureBytes(size int) *SecureBytes {
	sb := &SecureBytes{
		data: make([]byte, size),
	}
	runtime.SetFinalizer(sb, (*SecureBytes).Destroy)
	return sb
}

// Destroy zeroizes memory before GC
func (sb *SecureBytes) Destroy() {
	if sb.data != nil {
		// Zeroize memory before GC
		for i := range sb.data {
			sb.data[i] = 0
		}
		sb.data = nil
	}
}

// Copy copies data into secure container
func (sb *SecureBytes) Copy(data []byte) {
	if len(data) <= len(sb.data) {
		copy(sb.data, data)
	}
}

// Bytes returns the secure data
func (sb *SecureBytes) Bytes() []byte {
	return sb.data
}

// Size returns the size of the secure container
func (sb *SecureBytes) Size() int {
	return len(sb.data)
}

// Zeroize manually zeroizes the memory
func (sb *SecureBytes) Zeroize() {
	for i := range sb.data {
		sb.data[i] = 0
	}
}

// SecureString provides secure string handling
type SecureString struct {
	data []byte
}

// NewSecureString creates a new secure string
func NewSecureString(s string) *SecureString {
	data := []byte(s)
	ss := &SecureString{
		data: make([]byte, len(data)),
	}
	copy(ss.data, data)

	// Zeroize the original string
	for i := range data {
		data[i] = 0
	}

	runtime.SetFinalizer(ss, (*SecureString).Destroy)
	return ss
}

// Destroy zeroizes memory before GC
func (ss *SecureString) Destroy() {
	if ss.data != nil {
		for i := range ss.data {
			ss.data[i] = 0
		}
		ss.data = nil
	}
}

// String returns the secure string
func (ss *SecureString) String() string {
	return string(ss.data)
}

// Zeroize manually zeroizes the memory
func (ss *SecureString) Zeroize() {
	for i := range ss.data {
		ss.data[i] = 0
	}
}

// SecureCompare performs constant-time comparison
func SecureCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	result := byte(0)
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}

	return result == 0
}

// SecureEqual performs constant-time string comparison
func SecureEqual(a, b string) bool {
	if len(a) != len(b) {
		return false
	}

	result := byte(0)
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}

	return result == 0
}

// ZeroizeBytes zeroizes a byte slice
func ZeroizeBytes(data []byte) {
	if data != nil {
		for i := range data {
			data[i] = 0
		}
	}
}

// ZeroizeString zeroizes a string (creates a new zeroed slice)
func ZeroizeString(s string) {
	data := unsafe.Slice(unsafe.StringData(s), len(s))
	ZeroizeBytes(data)
}
