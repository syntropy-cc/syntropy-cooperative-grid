// Package setup provides public types and interfaces for the setup component
package setup

import (
	"time"
)

// ConfigOptions define as opções de configuração
type ConfigOptions struct {
	OwnerName      string            `json:"owner_name"`
	OwnerEmail     string            `json:"owner_email"`
	NetworkConfig  *NetworkConfig    `json:"network_config"`
	SecurityConfig *SecurityConfig   `json:"security_config"`
	CustomSettings map[string]string `json:"custom_settings"`
}

// SetupOptions define as opções do setup
type SetupOptions struct {
	Force          bool              `json:"force"`
	ValidateOnly   bool              `json:"validate_only"`
	SkipValidation bool              `json:"skip_validation"`
	TestMode       bool              `json:"test_mode"`       // Bypass strict validation for unit tests
	Verbose        bool              `json:"verbose"`
	Quiet          bool              `json:"quiet"`
	ConfigPath     string            `json:"config_path"`
	CustomSettings map[string]string `json:"custom_settings"`
}

// ValidationResult resultado da validação
type ValidationResult struct {
	Environment  *EnvironmentInfo  `json:"environment"`
	Dependencies *DependencyStatus `json:"dependencies"`
	Network      *NetworkInfo      `json:"network"`
	Permissions  *PermissionStatus `json:"permissions"`
	Issues       []ValidationIssue `json:"issues"`
	CanProceed   bool              `json:"can_proceed"`
	Warnings     []string          `json:"warnings"`
}

// SetupState estado do setup
type SetupState struct {
	Version       string            `json:"version"`
	Timestamp     time.Time         `json:"timestamp"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	Status        SetupStatus       `json:"status"`
	Environment   *EnvironmentInfo  `json:"environment"`
	Configuration *ConfigInfo       `json:"configuration"`
	Config        *ConfigInfo       `json:"config"`
	Keys          *KeyInfo          `json:"keys"`
	LastBackup    *BackupInfo       `json:"last_backup"`
	Metadata      map[string]string `json:"metadata"`
}

// KeyPair par de chaves
type KeyPair struct {
	ID          string            `json:"id"`
	Algorithm   string            `json:"algorithm"`
	PrivateKey  []byte            `json:"private_key"`
	PublicKey   []byte            `json:"public_key"`
	CreatedAt   time.Time         `json:"created_at"`
	ExpiresAt   time.Time         `json:"expires_at"`
	Fingerprint string            `json:"fingerprint"`
	Metadata    map[string]string `json:"metadata"`
}

// EnvironmentInfo informações do ambiente
type EnvironmentInfo struct {
	OS              string   `json:"os"`
	OSVersion       string   `json:"os_version"`
	Architecture    string   `json:"architecture"`
	HasAdminRights  bool     `json:"has_admin_rights"`
	PowerShellVer   string   `json:"powershell_ver"`
	AvailableDiskGB float64  `json:"available_disk_gb"`
	HasInternet     bool     `json:"has_internet"`
	HomeDir         string   `json:"home_dir"`
	CanProceed      bool     `json:"can_proceed"`
	Issues          []string `json:"issues"`
}

// DependencyStatus status das dependências
type DependencyStatus struct {
	Required  []Dependency `json:"required"`
	Installed []Dependency `json:"installed"`
	Missing   []Dependency `json:"missing"`
	Outdated  []Dependency `json:"outdated"`
}

// Dependency dependência
type Dependency struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Required  bool   `json:"required"`
	Installed bool   `json:"installed"`
	Path      string `json:"path,omitempty"`
}

// NetworkInfo informações de rede
type NetworkInfo struct {
	HasInternet     bool  `json:"has_internet"`
	Connectivity    bool  `json:"connectivity"`
	ProxyConfigured bool  `json:"proxy_configured"`
	FirewallActive  bool  `json:"firewall_active"`
	PortsOpen       []int `json:"ports_open"`
}

// PermissionStatus status das permissões
type PermissionStatus struct {
	FileSystem bool     `json:"file_system"`
	Network    bool     `json:"network"`
	Service    bool     `json:"service"`
	Admin      bool     `json:"admin"`
	Issues     []string `json:"issues"`
}

// ValidationIssue issue de validação
type ValidationIssue struct {
	Type        string   `json:"type"`
	Severity    string   `json:"severity"`
	Message     string   `json:"message"`
	Suggestions []string `json:"suggestions"`
}

// SetupStatus status do setup
type SetupStatus string

const (
	SetupStatusNotStarted SetupStatus = "not_started"
	SetupStatusInProgress SetupStatus = "in_progress"
	SetupStatusCompleted  SetupStatus = "completed"
	SetupStatusFailed     SetupStatus = "failed"
	SetupStatusCorrupted  SetupStatus = "corrupted"
)

// NetworkConfig configuração de rede
type NetworkConfig struct {
	Discovery *DiscoveryConfig `json:"discovery"`
	API       *APIConfig       `json:"api"`
}

// DiscoveryConfig configuração de descoberta
type DiscoveryConfig struct {
	Enabled bool `json:"enabled"`
	Port    int  `json:"port"`
}

// APIConfig configuração da API
type APIConfig struct {
	Enabled bool       `json:"enabled"`
	Port    int        `json:"port"`
	TLS     *TLSConfig `json:"tls"`
}

// TLSConfig configuração TLS
type TLSConfig struct {
	Enabled  bool   `json:"enabled"`
	CertPath string `json:"cert_path"`
	KeyPath  string `json:"key_path"`
}

// SecurityConfig configuração de segurança
type SecurityConfig struct {
	Encryption *EncryptionConfig `json:"encryption"`
	Backup     *BackupConfig     `json:"backup"`
}

// EncryptionConfig configuração de criptografia
type EncryptionConfig struct {
	Algorithm   string             `json:"algorithm"`
	KeyRotation *KeyRotationConfig `json:"key_rotation"`
}

// KeyRotationConfig configuração de rotação de chaves
type KeyRotationConfig struct {
	Enabled  bool   `json:"enabled"`
	Interval string `json:"interval"`
}

// BackupConfig configuração de backup
type BackupConfig struct {
	Enabled   bool   `json:"enabled"`
	Interval  string `json:"interval"`
	Retention string `json:"retention"`
}

// ConfigInfo informações de configuração
type ConfigInfo struct {
	Path      string    `json:"path"`
	Valid     bool      `json:"valid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// KeyInfo informações de chaves
type KeyInfo struct {
	OwnerKeyID string    `json:"owner_key_id"`
	Algorithm  string    `json:"algorithm"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiresAt  time.Time `json:"expires_at"`
}

// BackupInfo informações de backup
type BackupInfo struct {
	ID        string    `json:"id"`
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
}
