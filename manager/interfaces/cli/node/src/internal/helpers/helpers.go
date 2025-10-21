package helpers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ssh"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/node/src/internal/constants"
)

// GenerateNodeID generates a sequential node ID
func GenerateNodeID(existingNodes []string) string {
	maxNum := 0

	// Parse existing node IDs to find the highest number
	for _, nodeID := range existingNodes {
		if strings.HasPrefix(nodeID, "node-") {
			if numStr := strings.TrimPrefix(nodeID, "node-"); numStr != "" {
				if num, err := strconv.Atoi(numStr); err == nil && num > maxNum {
					maxNum = num
				}
			}
		}
	}

	// Generate next sequential ID
	nextNum := maxNum + 1
	return fmt.Sprintf("node-%02d", nextNum)
}

// ValidateNodeID validates the format of a node ID
func ValidateNodeID(nodeID string) error {
	if len(nodeID) < constants.MinNodeIDLength || len(nodeID) > constants.MaxNodeIDLength {
		return fmt.Errorf("node ID length must be between %d and %d characters",
			constants.MinNodeIDLength, constants.MaxNodeIDLength)
	}

	pattern := regexp.MustCompile(constants.NodeIDPattern)
	if !pattern.MatchString(nodeID) {
		return fmt.Errorf("invalid node ID format: must match pattern %s", constants.NodeIDPattern)
	}

	return nil
}

// SSHKeys represents a pair of SSH keys
type SSHKeys struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

// GenerateSSHKeys generates a pair of SSH keys (RSA 2048)
func GenerateSSHKeys() (*SSHKeys, error) {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, constants.RSAKeySize)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	// Get public key
	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get public key: %w", err)
	}

	// Encode private key
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	// Encode public key
	publicKeySSH := ssh.MarshalAuthorizedKey(publicKey)

	return &SSHKeys{
		PrivateKey: string(privateKeyPEM),
		PublicKey:  string(publicKeySSH),
	}, nil
}

// GenerateNodeCertificate generates an Ed25519 certificate for a node
func GenerateNodeCertificate(nodeID string) (string, error) {
	// Generate Ed25519 key pair
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", fmt.Errorf("failed to generate Ed25519 key pair: %w", err)
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   fmt.Sprintf("syntropy-node-%s", nodeID),
			Organization: []string{"syntropy-grid"},
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(constants.DefaultCertValidity),
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
	}

	// Create certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey, privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to create certificate: %w", err)
	}

	// Encode certificate to PEM
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})

	return string(certPEM), nil
}

// DetectCommandStationIP detects the local IP address of the command station
func DetectCommandStationIP() (string, error) {
	// Get local IP address
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", fmt.Errorf("failed to detect local IP: %w", err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

// ExpandPath expands a path with tilde and environment variables
func ExpandPath(path string) (string, error) {
	// Expand tilde
	if strings.HasPrefix(path, "~") {
		usr, err := user.Current()
		if err != nil {
			return "", fmt.Errorf("failed to get current user: %w", err)
		}
		path = filepath.Join(usr.HomeDir, path[1:])
	}

	// Expand environment variables
	path = os.ExpandEnv(path)

	// Clean the path
	path = filepath.Clean(path)

	return path, nil
}

// EnsureDirectory creates a directory if it doesn't exist
func EnsureDirectory(path string) error {
	expandedPath, err := ExpandPath(path)
	if err != nil {
		return fmt.Errorf("failed to expand path: %w", err)
	}

	if err := os.MkdirAll(expandedPath, constants.FileModeDirectory); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", expandedPath, err)
	}

	return nil
}

// WriteFileWithPermissions writes a file with specific permissions
func WriteFileWithPermissions(path string, data []byte, mode os.FileMode) error {
	expandedPath, err := ExpandPath(path)
	if err != nil {
		return fmt.Errorf("failed to expand path: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(expandedPath)
	if err := EnsureDirectory(dir); err != nil {
		return fmt.Errorf("failed to ensure directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(expandedPath, data, mode); err != nil {
		return fmt.Errorf("failed to write file %s: %w", expandedPath, err)
	}

	return nil
}

// ReadFileSafely reads a file with error handling
func ReadFileSafely(path string) ([]byte, error) {
	expandedPath, err := ExpandPath(path)
	if err != nil {
		return nil, fmt.Errorf("failed to expand path: %w", err)
	}

	data, err := os.ReadFile(expandedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", expandedPath, err)
	}

	return data, nil
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	expandedPath, err := ExpandPath(path)
	if err != nil {
		return false
	}

	_, err = os.Stat(expandedPath)
	return !os.IsNotExist(err)
}

// GenerateRandomBytes generates cryptographically secure random bytes
func GenerateRandomBytes(size int) ([]byte, error) {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return bytes, nil
}

// GenerateSecureToken generates a secure random token
func GenerateSecureToken(length int) (string, error) {
	bytes, err := GenerateRandomBytes(length)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}

// HashData creates a SHA256 hash of the input data
func HashData(data []byte) string {
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}

// ValidateHash validates a hash against data
func ValidateHash(data []byte, hash string) bool {
	expectedHash := HashData(data)
	return expectedHash == hash
}

// IsPortAvailable checks if a port is available
func IsPortAvailable(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

// GetNextAvailablePort finds the next available port starting from a given port
func GetNextAvailablePort(startPort int) int {
	for port := startPort; port < startPort+100; port++ {
		if IsPortAvailable(port) {
			return port
		}
	}
	return -1
}

// TestConnectivity tests network connectivity to a host and port
func TestConnectivity(host string, port int) error {
	timeout := time.Duration(5) * time.Second
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), timeout)
	if err != nil {
		return fmt.Errorf("failed to connect to %s:%d: %w", host, port, err)
	}
	conn.Close()
	return nil
}

// FormatDuration formats a duration in a human-readable format
func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0fs", d.Seconds())
	} else if d < time.Hour {
		return fmt.Sprintf("%.0fm", d.Minutes())
	} else if d < 24*time.Hour {
		return fmt.Sprintf("%.0fh", d.Hours())
	} else {
		return fmt.Sprintf("%.0fd", d.Hours()/24)
	}
}

// FormatBytes formats bytes in a human-readable format
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// CalculatePercentage calculates percentage with proper bounds
func CalculatePercentage(value, total int64) float64 {
	if total == 0 {
		return 0
	}
	percentage := float64(value) / float64(total) * 100
	if percentage > 100 {
		return 100
	}
	if percentage < 0 {
		return 0
	}
	return percentage
}

// RetryOperation retries an operation with exponential backoff
func RetryOperation(operation func() error, maxRetries int, initialDelay time.Duration) error {
	var err error
	delay := initialDelay

	for attempt := 0; attempt <= maxRetries; attempt++ {
		err = operation()
		if err == nil {
			return nil
		}

		if attempt == maxRetries {
			break
		}

		time.Sleep(delay)
		delay = time.Duration(float64(delay) * constants.RetryBackoffMultiplier)
		if delay > constants.MaxRetryDelay {
			delay = constants.MaxRetryDelay
		}
	}

	return fmt.Errorf("operation failed after %d attempts: %w", maxRetries+1, err)
}

// ValidateJSON validates if a string is valid JSON
func ValidateJSON(jsonStr string) error {
	var jsonData interface{}
	return json.Unmarshal([]byte(jsonStr), &jsonData)
}

// PrettyPrintJSON formats JSON in a pretty format
func PrettyPrintJSON(data interface{}) (string, error) {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return string(jsonBytes), nil
}

// GetCurrentTimestamp returns the current timestamp in RFC3339 format
func GetCurrentTimestamp() string {
	return time.Now().UTC().Format(constants.TimeFormatRFC3339)
}

// ParseTimestamp parses a timestamp string
func ParseTimestamp(timestamp string) (time.Time, error) {
	return time.Parse(constants.TimeFormatRFC3339, timestamp)
}

// GetPlatform returns the current platform
func GetPlatform() string {
	return runtime.GOOS
}

// GetArchitecture returns the current architecture
func GetArchitecture() string {
	return runtime.GOARCH
}

// IsWindows checks if running on Windows
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// IsLinux checks if running on Linux
func IsLinux() bool {
	return runtime.GOOS == "linux"
}

// IsMacOS checks if running on macOS
func IsMacOS() bool {
	return runtime.GOOS == "darwin"
}

// CopyFile copies a file from source to destination
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}

// RemoveFile removes a file if it exists
func RemoveFile(path string) error {
	if FileExists(path) {
		if err := os.Remove(path); err != nil {
			return fmt.Errorf("failed to remove file %s: %w", path, err)
		}
	}
	return nil
}

// GetFileSize returns the size of a file
func GetFileSize(path string) (int64, error) {
	expandedPath, err := ExpandPath(path)
	if err != nil {
		return 0, fmt.Errorf("failed to expand path: %w", err)
	}

	fileInfo, err := os.Stat(expandedPath)
	if err != nil {
		return 0, fmt.Errorf("failed to get file info: %w", err)
	}

	return fileInfo.Size(), nil
}

// ChmodFile changes file permissions
func ChmodFile(path string, mode os.FileMode) error {
	expandedPath, err := ExpandPath(path)
	if err != nil {
		return fmt.Errorf("failed to expand path: %w", err)
	}

	if err := os.Chmod(expandedPath, mode); err != nil {
		return fmt.Errorf("failed to change file permissions: %w", err)
	}

	return nil
}

// CreateTempFile creates a temporary file with the given content
func CreateTempFile(content string) (string, error) {
	tmpFile, err := os.CreateTemp("", "syntropy-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}

	if _, err := tmpFile.WriteString(content); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("failed to write to temp file: %w", err)
	}

	tmpFile.Close()
	return tmpFile.Name(), nil
}

// CleanupTempFile removes a temporary file
func CleanupTempFile(path string) error {
	return os.Remove(path)
}

// GetDirectorySize returns the total size of a directory
func GetDirectorySize(dirPath string) (int64, error) {
	var totalSize int64

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			totalSize += info.Size()
		}

		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to walk directory: %w", err)
	}

	return totalSize, nil
}

// GetDeviceNameFromPath extracts device name from device path
func GetDeviceNameFromPath(devicePath string) string {
	// Extract device name from path like /dev/sda -> sda
	if strings.HasPrefix(devicePath, "/dev/") {
		return strings.TrimPrefix(devicePath, "/dev/")
	}

	// For Windows paths like C:, return as is
	if strings.HasSuffix(devicePath, ":") {
		return devicePath
	}

	// For other paths, return the last component
	return filepath.Base(devicePath)
}

// ReadJSONFile reads a JSON file and unmarshals it into the target
func ReadJSONFile(filePath string, target interface{}) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to unmarshal JSON from file %s: %w", filePath, err)
	}

	return nil
}

// WriteJSONFile marshals the source to JSON and writes it to a file
func WriteJSONFile(filePath string, source interface{}) error {
	data, err := json.MarshalIndent(source, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filePath, err)
	}

	return nil
}
