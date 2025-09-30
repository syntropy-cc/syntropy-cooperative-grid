package types

import (
	"time"
)

// SetupManager interface principal do setup
type SetupManager interface {
	// Validação do ambiente
	Validate() (*ValidationResult, error)

	// Execução do setup completo
	Setup(options *SetupOptions) error

	// Verificação de status
	Status() (*SetupStatus, error)

	// Reset do setup
	Reset(confirm bool) error

	// Reparo automático
	Repair() error
}

// Validator interface de validação
type Validator interface {
	// Validação completa do ambiente
	ValidateEnvironment() (*EnvironmentInfo, error)

	// Validação de dependências
	ValidateDependencies() (*DependencyStatus, error)

	// Validação de rede
	ValidateNetwork() (*NetworkInfo, error)

	// Validação de permissões
	ValidatePermissions() (*PermissionStatus, error)

	// Correção automática de problemas
	FixIssues(issues []ValidationIssue) error
}

// Configurator interface de configuração
type Configurator interface {
	// Geração de configuração principal
	GenerateConfig(options *ConfigOptions) error

	// Criação da estrutura de diretórios
	CreateStructure() error

	// Geração de chaves criptográficas
	GenerateKeys() (*KeyPair, error)

	// Validação de configuração
	ValidateConfig() error

	// Backup de configuração
	BackupConfig(name string) error

	// Restauração de configuração
	RestoreConfig(backupPath string) error
}

// StateManager interface de gerenciamento de estado
type StateManager interface {
	// Carregamento do estado atual
	LoadState() (*SetupState, error)

	// Salvamento atômico do estado
	SaveState(state *SetupState) error

	// Atualização atômica do estado
	UpdateState(update func(*SetupState) error) error

	// Backup do estado
	BackupState(name string) error

	// Restauração do estado
	RestoreState(backupPath string) error

	// Verificação de integridade
	VerifyIntegrity() error
}

// KeyManager interface de gerenciamento de chaves
type KeyManager interface {
	// Geração de par de chaves
	GenerateKeyPair(algorithm string) (*KeyPair, error)

	// Geração ou carregamento de chaves existentes
	GenerateOrLoadKeyPair(algorithm string) (*KeyPair, error)

	// Armazenamento seguro de chaves
	StoreKeyPair(keyPair *KeyPair, passphrase string) error

	// Carregamento de chaves
	LoadKeyPair(keyID string, passphrase string) (*KeyPair, error)

	// Rotação de chaves
	RotateKeys(keyID string) error

	// Verificação de integridade
	VerifyKeyIntegrity(keyID string) error

	// Backup de chaves
	BackupKeys(keyID string, passphrase string) ([]byte, error)

	// Restauração de chaves
	RestoreKeys(backupData []byte, passphrase string) error
}

// SetupLogger interface de logging
type SetupLogger interface {
	// Logging de etapas do setup
	LogStep(step string, data map[string]interface{})

	// Logging de erros
	LogError(err error, context map[string]interface{})

	// Logging de warnings
	LogWarning(message string, data map[string]interface{})

	// Logging de informações
	LogInfo(message string, data map[string]interface{})

	// Logging de debug
	LogDebug(message string, data map[string]interface{})

	// Exportação de logs
	ExportLogs(format string, outputPath string) error

	// Fechar logger
	Close() error
}

// OSValidator interface de validação por SO
type OSValidator interface {
	// Detecção do sistema operacional
	DetectOS() (*OSInfo, error)

	// Validação de recursos do sistema
	ValidateResources() (*ResourceInfo, error)

	// Validação de permissões
	ValidatePermissions() (*PermissionInfo, error)

	// Instalação de dependências
	InstallDependencies(deps []Dependency) error

	// Configuração de ambiente
	ConfigureEnvironment() error
}

// Estruturas de dados principais

// SetupOptions define as opções do setup
type SetupOptions struct {
	Force          bool              `json:"force"`
	ValidateOnly   bool              `json:"validate_only"`
	Verbose        bool              `json:"verbose"`
	Quiet          bool              `json:"quiet"`
	ConfigPath     string            `json:"config_path"`
	CustomSettings map[string]string `json:"custom_settings"`
}

// ConfigOptions define as opções de configuração
type ConfigOptions struct {
	OwnerName      string            `json:"owner_name"`
	OwnerEmail     string            `json:"owner_email"`
	NetworkConfig  *NetworkConfig    `json:"network_config"`
	SecurityConfig *SecurityConfig   `json:"security_config"`
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
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	Status        SetupStatus       `json:"status"`
	Environment   *EnvironmentInfo  `json:"environment"`
	Configuration *ConfigInfo       `json:"configuration"`
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

// LogEntry entrada de log
type LogEntry struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Step      string                 `json:"step,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Error     string                 `json:"error,omitempty"`
}

// OSInfo informações do SO
type OSInfo struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	Architecture string `json:"architecture"`
	Build        string `json:"build"`
	Kernel       string `json:"kernel"`
}

// ResourceInfo informações de recursos
type ResourceInfo struct {
	TotalMemoryGB  float64 `json:"total_memory_gb"`
	AvailableMemGB float64 `json:"available_mem_gb"`
	CPUCores       int     `json:"cpu_cores"`
	DiskSpaceGB    float64 `json:"disk_space_gb"`
}

// PermissionInfo informações de permissões
type PermissionInfo struct {
	HasAdminRights bool     `json:"has_admin_rights"`
	UserID         string   `json:"user_id"`
	GroupID        string   `json:"group_id"`
	Capabilities   []string `json:"capabilities"`
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
