// Package setup provides security validation functionality
package setup

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

// SecurityValidator provides security validation functionality
type SecurityValidator struct {
	attemptCounts map[string]int
	lastAttempt   map[string]time.Time
	blockedIPs    map[string]time.Time
}

// NewSecurityValidator creates a new security validator
func NewSecurityValidator() *SecurityValidator {
	return &SecurityValidator{
		attemptCounts: make(map[string]int),
		lastAttempt:   make(map[string]time.Time),
		blockedIPs:    make(map[string]time.Time),
	}
}

// ValidatePath validates file paths for security issues
func (sv *SecurityValidator) ValidatePath(path string, baseDir string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	// Check for directory traversal
	if strings.Contains(path, "..") {
		return fmt.Errorf("invalid path: directory traversal detected")
	}

	// Check for absolute paths to sensitive directories
	sensitivePaths := []string{
		"/etc/", "/root/", "/sys/", "/proc/", "/dev/",
		"C:\\Windows\\", "C:\\System32\\", "C:\\Program Files\\",
	}

	for _, sensitive := range sensitivePaths {
		if strings.HasPrefix(strings.ToLower(path), strings.ToLower(sensitive)) {
			return fmt.Errorf("invalid path: access to sensitive directory denied")
		}
	}

	// Ensure path is within base directory (only if baseDir is specified and path is relative)
	if baseDir != "" && !filepath.IsAbs(path) {
		cleanPath := filepath.Clean(path)
		cleanBase := filepath.Clean(baseDir)

		relPath, err := filepath.Rel(cleanBase, cleanPath)
		if err != nil || strings.HasPrefix(relPath, "..") {
			return fmt.Errorf("invalid path: outside allowed directory")
		}
	}

	return nil
}

// ValidateCryptoAlgorithm validates cryptographic algorithms
func (sv *SecurityValidator) ValidateCryptoAlgorithm(algorithm string) error {
	weakAlgorithms := map[string]bool{
		"DES":      true,
		"3DES":     true,
		"RC4":      true,
		"MD5":      true,
		"SHA1":     true,
		"MD4":      true,
		"MD2":      true,
		"RC2":      true,
		"BLOWFISH": true,
	}

	algorithm = strings.ToUpper(algorithm)
	if weakAlgorithms[algorithm] {
		return fmt.Errorf("weak cryptographic algorithm not allowed: %s", algorithm)
	}

	// Check for strong algorithms
	strongAlgorithms := []string{"AES", "SHA256", "SHA512", "RSA", "ECDSA", "ED25519"}
	isStrong := false
	for _, strong := range strongAlgorithms {
		if strings.Contains(algorithm, strong) {
			isStrong = true
			break
		}
	}

	if !isStrong {
		return fmt.Errorf("unknown or potentially weak algorithm: %s", algorithm)
	}

	return nil
}

// ValidateCredentials validates authentication credentials
func (sv *SecurityValidator) ValidateCredentials(credType, credential string) error {
	if credential == "" {
		return fmt.Errorf("credential cannot be empty")
	}

	// Check for weak passwords
	weakPasswords := []string{
		"123456", "password", "admin", "root", "user", "guest",
		"qwerty", "abc123", "password123", "123456789", "letmein",
	}

	credLower := strings.ToLower(credential)
	for _, weak := range weakPasswords {
		if credLower == weak {
			return fmt.Errorf("weak %s not allowed", credType)
		}
	}

	// Check minimum length
	if len(credential) < 8 {
		return fmt.Errorf("%s too short: minimum 8 characters required", credType)
	}

	// Check for API keys (should be longer)
	if credType == "key" && len(credential) < 16 {
		return fmt.Errorf("API key too short: minimum 16 characters required")
	}

	return nil
}

// ValidateURL validates URLs for SSRF prevention
func (sv *SecurityValidator) ValidateURL(urlStr string) error {
	if urlStr == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL format: %v", err)
	}

	// Check for dangerous schemes
	dangerousSchemes := []string{"file", "ftp", "gopher", "jar"}
	for _, scheme := range dangerousSchemes {
		if parsedURL.Scheme == scheme {
			return fmt.Errorf("dangerous URL scheme not allowed: %s", scheme)
		}
	}

	// Check for localhost/internal IPs
	host := strings.ToLower(parsedURL.Hostname())
	dangerousHosts := []string{
		"localhost", "127.0.0.1", "0.0.0.0", "::1",
		"169.254.169.254",                       // AWS metadata
		"10.0.0.0", "172.16.0.0", "192.168.0.0", // Private networks
	}

	for _, dangerous := range dangerousHosts {
		if host == dangerous || strings.HasPrefix(host, dangerous+".") {
			return fmt.Errorf("internal/localhost URL not allowed: %s", host)
		}
	}

	// Check for common internal ports
	if parsedURL.Port() != "" {
		internalPorts := []string{"22", "23", "25", "53", "80", "135", "139", "445", "1433", "3306", "3389", "5432", "6379", "27017"}
		for _, port := range internalPorts {
			if parsedURL.Port() == port {
				return fmt.Errorf("internal port not allowed: %s", port)
			}
		}
	}

	return nil
}

// ValidateConfiguration validates configuration content for security issues
func (sv *SecurityValidator) ValidateConfiguration(config string) error {
	if config == "" {
		return fmt.Errorf("configuration cannot be empty")
	}

	// Check for command injection patterns
	injectionPatterns := []string{
		";", "&&", "||", "|", "`", "$(", "$", "\\",
		"rm -rf", "del /f", "format", "shutdown", "reboot",
		"curl", "wget", "nc", "netcat", "telnet",
	}

	configLower := strings.ToLower(config)
	for _, pattern := range injectionPatterns {
		if strings.Contains(configLower, pattern) {
			return fmt.Errorf("potentially malicious configuration detected: %s", pattern)
		}
	}

	// Check for debug/insecure settings
	insecureSettings := []string{
		"debug: true", "log_level: debug", "expose_internal_apis: true",
		"disable_ssl: true", "allow_weak_ciphers: true",
		"world_readable: true", "allow_all: true",
	}

	for _, setting := range insecureSettings {
		if strings.Contains(configLower, setting) {
			return fmt.Errorf("insecure configuration detected: %s", setting)
		}
	}

	return nil
}

// CheckBruteForceProtection checks for brute force attempts
func (sv *SecurityValidator) CheckBruteForceProtection(identifier string) error {
	const maxAttempts = 5
	const blockDuration = 15 * time.Minute

	now := time.Now()

	// Check if IP is currently blocked
	if blockTime, exists := sv.blockedIPs[identifier]; exists {
		if now.Sub(blockTime) < blockDuration {
			return fmt.Errorf("too many failed attempts: blocked for %v", blockDuration-now.Sub(blockTime))
		}
		// Unblock if time has passed
		delete(sv.blockedIPs, identifier)
		delete(sv.attemptCounts, identifier)
		delete(sv.lastAttempt, identifier)
	}

	// Check attempt count
	if count, exists := sv.attemptCounts[identifier]; exists {
		if count >= maxAttempts {
			sv.blockedIPs[identifier] = now
			return fmt.Errorf("too many failed attempts: blocked for %v", blockDuration)
		}
	}

	return nil
}

// RecordFailedAttempt records a failed authentication attempt
func (sv *SecurityValidator) RecordFailedAttempt(identifier string) {
	sv.attemptCounts[identifier]++
	sv.lastAttempt[identifier] = time.Now()
}

// RecordSuccessfulAttempt clears failed attempts for an identifier
func (sv *SecurityValidator) RecordSuccessfulAttempt(identifier string) {
	delete(sv.attemptCounts, identifier)
	delete(sv.lastAttempt, identifier)
	delete(sv.blockedIPs, identifier)
}

// ValidateFilePermissions validates file permissions for security
func (sv *SecurityValidator) ValidateFilePermissions(filePath string) error {
	if runtime.GOOS == "windows" {
		// Windows permission validation would be different
		return nil
	}

	info, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("cannot stat file: %v", err)
	}

	mode := info.Mode().Perm()

	// Check if file has group or other permissions
	if mode&0077 != 0 {
		return fmt.Errorf("insecure file permissions: file should not be readable/writable by group or other")
	}

	return nil
}

// SetSecureFilePermissions sets secure file permissions
func (sv *SecurityValidator) SetSecureFilePermissions(filePath string) error {
	if runtime.GOOS == "windows" {
		// Windows permission setting would be different
		return nil
	}

	return os.Chmod(filePath, 0600)
}

// GenerateSecureKey generates a cryptographically secure key
func (sv *SecurityValidator) GenerateSecureKey(length int) ([]byte, error) {
	if length < 16 {
		return nil, fmt.Errorf("key length too short: minimum 16 bytes required")
	}

	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate secure key: %v", err)
	}

	return key, nil
}

// ValidateKeyStrength validates the strength of a cryptographic key
func (sv *SecurityValidator) ValidateKeyStrength(key []byte) error {
	if len(key) < 16 {
		return fmt.Errorf("key too short: minimum 16 bytes required")
	}

	// Check for sufficient entropy (basic test)
	zeroCount := 0
	for _, b := range key {
		if b == 0 {
			zeroCount++
		}
	}

	// Key should not have too many zero bytes
	if zeroCount > len(key)/4 {
		return fmt.Errorf("key has insufficient entropy")
	}

	return nil
}

// SanitizeLogInput sanitizes input for logging to prevent log injection
func (sv *SecurityValidator) SanitizeLogInput(input string) string {
	// Remove control characters and newlines
	sanitized := strings.ReplaceAll(input, "\n", "")
	sanitized = strings.ReplaceAll(sanitized, "\r", "")
	sanitized = strings.ReplaceAll(sanitized, "\x00", "")
	sanitized = strings.ReplaceAll(sanitized, "\t", " ")

	// Remove other control characters
	re := regexp.MustCompile(`[\x00-\x1F\x7F]`)
	sanitized = re.ReplaceAllString(sanitized, "")

	// Mask sensitive patterns
	sensitivePatterns := []string{"password", "key", "token", "secret", "api_key"}
	for _, pattern := range sensitivePatterns {
		if strings.Contains(strings.ToLower(sanitized), pattern) {
			sanitized = "***REDACTED***"
			break
		}
	}

	return sanitized
}

// ValidateFileType validates if a file type is allowed
func (sv *SecurityValidator) ValidateFileType(filename string) error {
	allowedExtensions := []string{".yaml", ".yml", ".json", ".txt", ".conf", ".toml"}
	dangerousExtensions := []string{".sh", ".exe", ".bat", ".ps1", ".cmd", ".com", ".scr", ".pif"}

	ext := strings.ToLower(filepath.Ext(filename))

	// Check for dangerous extensions
	for _, dangerous := range dangerousExtensions {
		if ext == dangerous {
			return fmt.Errorf("dangerous file type not allowed: %s", ext)
		}
	}

	// Check for allowed extensions
	for _, allowed := range allowedExtensions {
		if ext == allowed {
			return nil
		}
	}

	return fmt.Errorf("file type not allowed: %s", ext)
}

// CalculateFileChecksum calculates SHA256 checksum of a file
func (sv *SecurityValidator) CalculateFileChecksum(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	hash := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(hash[:]), nil
}

// ValidateChecksum validates a file against its expected checksum
func (sv *SecurityValidator) ValidateChecksum(filePath, expectedChecksum string) error {
	actualChecksum, err := sv.CalculateFileChecksum(filePath)
	if err != nil {
		return err
	}

	if actualChecksum != expectedChecksum {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedChecksum, actualChecksum)
	}

	return nil
}

// CheckAdminRights checks if the current process has administrative rights
func (sv *SecurityValidator) CheckAdminRights() bool {
	if runtime.GOOS == "windows" {
		// Windows admin check would be more complex
		return false
	}
	return os.Getuid() == 0
}

// ValidateEnvironment validates the current environment for security
func (sv *SecurityValidator) ValidateEnvironment() error {
	// Check if running as root (security risk)
	if sv.CheckAdminRights() {
		return fmt.Errorf("running with administrative privileges is not recommended for security")
	}

	// Check for dangerous environment variables
	dangerousVars := []string{"LD_PRELOAD", "LD_LIBRARY_PATH", "PATH"}
	for _, varName := range dangerousVars {
		if value := os.Getenv(varName); value != "" {
			// Check for suspicious values
			if strings.Contains(value, "..") || strings.Contains(value, "/tmp") {
				return fmt.Errorf("suspicious environment variable %s: %s", varName, value)
			}
		}
	}

	return nil
}
