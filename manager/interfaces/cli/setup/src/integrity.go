package setup

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
)

// IntegrityVerifier provides integrity verification capabilities
type IntegrityVerifier struct{}

// NewIntegrityVerifier creates a new integrity verifier
func NewIntegrityVerifier() *IntegrityVerifier {
	return &IntegrityVerifier{}
}

// ComputeMAC computes HMAC-SHA256 for data integrity
func (iv *IntegrityVerifier) ComputeMAC(data []byte, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

// ComputeMAC512 computes HMAC-SHA512 for data integrity
func (iv *IntegrityVerifier) ComputeMAC512(data []byte, key []byte) []byte {
	h := hmac.New(sha512.New, key)
	h.Write(data)
	return h.Sum(nil)
}

// VerifyMAC verifies HMAC-SHA256 integrity
func (iv *IntegrityVerifier) VerifyMAC(data []byte, mac []byte, key []byte) bool {
	expectedMAC := iv.ComputeMAC(data, key)
	return hmac.Equal(mac, expectedMAC)
}

// VerifyMAC512 verifies HMAC-SHA512 integrity
func (iv *IntegrityVerifier) VerifyMAC512(data []byte, mac []byte, key []byte) bool {
	expectedMAC := iv.ComputeMAC512(data, key)
	return hmac.Equal(mac, expectedMAC)
}

// ComputeChecksum computes SHA256 checksum
func (iv *IntegrityVerifier) ComputeChecksum(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

// ComputeChecksum512 computes SHA512 checksum
func (iv *IntegrityVerifier) ComputeChecksum512(data []byte) []byte {
	hash := sha512.Sum512(data)
	return hash[:]
}

// VerifyChecksum verifies SHA256 checksum
func (iv *IntegrityVerifier) VerifyChecksum(data []byte, checksum []byte) bool {
	expectedChecksum := iv.ComputeChecksum(data)
	return hmac.Equal(checksum, expectedChecksum)
}

// VerifyChecksum512 verifies SHA512 checksum
func (iv *IntegrityVerifier) VerifyChecksum512(data []byte, checksum []byte) bool {
	expectedChecksum := iv.ComputeChecksum512(data)
	return hmac.Equal(checksum, expectedChecksum)
}

// ComputeChecksumB64 computes SHA256 checksum and returns base64 encoded
func (iv *IntegrityVerifier) ComputeChecksumB64(data []byte) string {
	checksum := iv.ComputeChecksum(data)
	return base64.StdEncoding.EncodeToString(checksum)
}

// VerifyChecksumB64 verifies base64 encoded SHA256 checksum
func (iv *IntegrityVerifier) VerifyChecksumB64(data []byte, checksumB64 string) bool {
	checksum, err := base64.StdEncoding.DecodeString(checksumB64)
	if err != nil {
		return false
	}
	return iv.VerifyChecksum(data, checksum)
}

// IntegrityData represents data with integrity protection
type IntegrityData struct {
	Data      []byte `json:"data"`
	MAC       []byte `json:"mac"`
	Algorithm string `json:"algorithm"`
}

// ProtectData protects data with HMAC integrity
func (iv *IntegrityVerifier) ProtectData(data []byte, key []byte) (*IntegrityData, error) {
	if len(key) < 32 {
		return nil, fmt.Errorf("key must be at least 32 bytes for integrity protection")
	}

	mac := iv.ComputeMAC(data, key)

	return &IntegrityData{
		Data:      data,
		MAC:       mac,
		Algorithm: "HMAC-SHA256",
	}, nil
}

// VerifyProtectedData verifies integrity of protected data
func (iv *IntegrityVerifier) VerifyProtectedData(protected *IntegrityData, key []byte) ([]byte, error) {
	if len(key) < 32 {
		return nil, fmt.Errorf("key must be at least 32 bytes for integrity verification")
	}

	if protected.Algorithm != "HMAC-SHA256" {
		return nil, fmt.Errorf("unsupported algorithm: %s", protected.Algorithm)
	}

	if !iv.VerifyMAC(protected.Data, protected.MAC, key) {
		return nil, fmt.Errorf("integrity verification failed")
	}

	return protected.Data, nil
}

// KeyedHash provides keyed hashing for integrity
type KeyedHash struct {
	key []byte
}

// NewKeyedHash creates a new keyed hash
func NewKeyedHash(key []byte) *KeyedHash {
	return &KeyedHash{key: key}
}

// Hash computes keyed hash
func (kh *KeyedHash) Hash(data []byte) []byte {
	h := hmac.New(sha256.New, kh.key)
	h.Write(data)
	return h.Sum(nil)
}

// HashB64 computes keyed hash and returns base64 encoded
func (kh *KeyedHash) HashB64(data []byte) string {
	hash := kh.Hash(data)
	return base64.StdEncoding.EncodeToString(hash)
}

// Verify verifies keyed hash
func (kh *KeyedHash) Verify(data []byte, hash []byte) bool {
	expectedHash := kh.Hash(data)
	return hmac.Equal(hash, expectedHash)
}

// VerifyB64 verifies base64 encoded keyed hash
func (kh *KeyedHash) VerifyB64(data []byte, hashB64 string) bool {
	hash, err := base64.StdEncoding.DecodeString(hashB64)
	if err != nil {
		return false
	}
	return kh.Verify(data, hash)
}
