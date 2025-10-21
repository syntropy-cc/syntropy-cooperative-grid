package setup

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"time"
)

// SecureRandom provides cryptographically secure random generation
type SecureRandom struct{}

// NewSecureRandom creates a new secure random generator
func NewSecureRandom() *SecureRandom {
	return &SecureRandom{}
}

// GenerateBytes generates cryptographically secure random bytes
func (sr *SecureRandom) GenerateBytes(n int) ([]byte, error) {
	if n <= 0 {
		return nil, fmt.Errorf("invalid length: %d", n)
	}

	b := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return b, nil
}

// GenerateString generates a random hex string
func (sr *SecureRandom) GenerateString(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("invalid length: %d", length)
	}

	bytes, err := sr.GenerateBytes(length / 2)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

// GenerateUUID generates a random UUID v4
func (sr *SecureRandom) GenerateUUID() (string, error) {
	bytes, err := sr.GenerateBytes(16)
	if err != nil {
		return "", err
	}

	// Set version (4) and variant bits
	bytes[6] = (bytes[6] & 0x0f) | 0x40 // Version 4
	bytes[8] = (bytes[8] & 0x3f) | 0x80 // Variant bits

	// Format as UUID
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		bytes[0:4],
		bytes[4:6],
		bytes[6:8],
		bytes[8:10],
		bytes[10:16],
	), nil
}

// GenerateToken generates a secure token
func (sr *SecureRandom) GenerateToken(length int) (string, error) {
	if length < 16 {
		return "", fmt.Errorf("token length must be at least 16 characters")
	}

	// Use UUID for tokens of 36 characters or less
	if length <= 36 {
		return sr.GenerateUUID()
	}

	// Generate hex string for longer tokens
	return sr.GenerateString(length)
}

// CryptoUtils provides cryptographic utility functions
type CryptoUtils struct {
	random *SecureRandom
}

// NewCryptoUtils creates a new crypto utils instance
func NewCryptoUtils() *CryptoUtils {
	return &CryptoUtils{
		random: NewSecureRandom(),
	}
}

// GenerateSalt generates a random salt
func (cu *CryptoUtils) GenerateSalt(size int) ([]byte, error) {
	if size < 16 {
		return nil, fmt.Errorf("salt size must be at least 16 bytes")
	}
	return cu.random.GenerateBytes(size)
}

// GenerateNonce generates a random nonce for encryption
func (cu *CryptoUtils) GenerateNonce(size int) ([]byte, error) {
	if size < 12 {
		return nil, fmt.Errorf("nonce size must be at least 12 bytes for GCM")
	}
	return cu.random.GenerateBytes(size)
}

// GenerateIV generates a random initialization vector
func (cu *CryptoUtils) GenerateIV(size int) ([]byte, error) {
	if size < 16 {
		return nil, fmt.Errorf("IV size must be at least 16 bytes")
	}
	return cu.random.GenerateBytes(size)
}

// GenerateKey generates a random encryption key
func (cu *CryptoUtils) GenerateKey(size int) ([]byte, error) {
	if size < 16 {
		return nil, fmt.Errorf("key size must be at least 16 bytes")
	}
	return cu.random.GenerateBytes(size)
}

// GeneratePassphrase generates a random passphrase
func (cu *CryptoUtils) GeneratePassphrase(length int) (string, error) {
	if length < 12 {
		return "", fmt.Errorf("passphrase length must be at least 12 characters")
	}

	// Use a mix of characters for passphrase
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*"

	passphrase := make([]byte, length)
	for i := range passphrase {
		randomBytes, err := cu.random.GenerateBytes(1)
		if err != nil {
			return "", err
		}
		passphrase[i] = chars[randomBytes[0]%byte(len(chars))]
	}

	return string(passphrase), nil
}

// GenerateTimestamp generates a timestamp-based random value
func (cu *CryptoUtils) GenerateTimestamp() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// GenerateSessionID generates a secure session ID
func (cu *CryptoUtils) GenerateSessionID() (string, error) {
	// Generate 32 bytes for session ID
	bytes, err := cu.random.GenerateBytes(32)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateAPIKey generates a secure API key
func (cu *CryptoUtils) GenerateAPIKey() (string, error) {
	// Generate 64 bytes for API key
	bytes, err := cu.random.GenerateBytes(64)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateChallenge generates a random challenge for authentication
func (cu *CryptoUtils) GenerateChallenge() (string, error) {
	// Generate 16 bytes for challenge
	bytes, err := cu.random.GenerateBytes(16)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// ConstantTimeCompare performs constant-time comparison
func (cu *CryptoUtils) ConstantTimeCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	result := byte(0)
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}

	return result == 0
}

// TimingSafeEqual performs timing-safe string comparison
func (cu *CryptoUtils) TimingSafeEqual(a, b string) bool {
	return cu.ConstantTimeCompare([]byte(a), []byte(b))
}
