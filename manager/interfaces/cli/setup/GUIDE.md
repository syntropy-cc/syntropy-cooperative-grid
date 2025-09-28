# Setup Component - Guia de Implementação para LLMs

## Contexto e Objetivos

### Syntropy Cooperative Grid
O **Syntropy Cooperative Grid** é uma rede descentralizada que permite a criação de uma infraestrutura computacional cooperativa. A rede opera de forma autônoma, permitindo que participantes compartilhem recursos computacionais de forma segura e eficiente através de um sistema de créditos e reputação.

### Syntropy Manager
O **Syntropy Manager** é a interface de controle para o Syntropy Cooperative Grid. Ele atua como um **controlador de estado** que modifica a rede descentralizada sem ser parte dela. A rede opera autonomamente; o manager apenas altera seu estado através de múltiplas interfaces (CLI, Web, Mobile, Desktop).

### Setup Component
O **Setup Component** é o componente responsável por configurar o **computador de trabalho** como um "quartel geral" para criação e gestão de nós da rede Syntropy. Este componente estabelece o ambiente inicial necessário para que o usuário possa criar, gerenciar e monitorar nós da rede através da CLI, funcionando como uma estação de controle centralizada.

## Princípios de Implementação

- **Simplicidade**: Arquitetura simples e direta, evitando over-engineering
- **Multiplataforma**: Suporte a Windows, Linux e macOS usando interfaces Go
- **Thread-Safe**: Operações atômicas e controle de concorrência
- **Segurança**: Sistema de chaves criptográficas robusto
- **Observabilidade**: Logging estruturado e métricas
- **Testabilidade**: Componentes desacoplados e testáveis
- **Manutenibilidade**: Código limpo e bem documentado

## Arquitetura Simplificada

```
┌─────────────────────────────────────────────────────────────┐
│ Setup Component (Orquestrador Principal)                   │
│ ─────────────────────────────────────────────────────────── │
│ • setup.go          • validator.go      • configurator.go   │
│ • state_manager.go  • key_manager.go    • logger.go         │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│ Serviços Internos (Implementação por Interface)           │
│ ─────────────────────────────────────────────────────────── │
│ • OSValidator       • DependencyManager • FileSystemOps     │
│ • CryptoProvider    • NetworkValidator  • ConfigGenerator   │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│ Sistema de Estado Local (~/.syntropy/)                     │
│ ─────────────────────────────────────────────────────────── │
│ • config/manager.yaml    • keys/owner.key*                 │
│ • logs/setup.log         • backups/backup_*.tar.gz         │
│ • nodes/ (estrutura)     • cache/ (temporários)            │
└─────────────────────────────────────────────────────────────┘
```

## Explicação da Arquitetura por Níveis

### Nível 1: Componentes Principais (Orquestração)

#### **setup.go** - Orquestrador Principal
- **Função**: Coordena todo o processo de setup, gerenciando o fluxo de execução
- **Responsabilidades**: Validação inicial, coordenação de subcomponentes, tratamento de erros globais
- **Interface**: `SetupManager` - ponto de entrada único para todas as operações

#### **validator.go** - Validador Unificado
- **Função**: Verifica se o ambiente está pronto para o setup
- **Responsabilidades**: Validação de SO, recursos, permissões, dependências e rede
- **Interface**: `Validator` - centraliza todas as verificações de pré-requisitos

#### **configurator.go** - Configurador Unificado
- **Função**: Executa a configuração propriamente dita do sistema
- **Responsabilidades**: Criação de estrutura de diretórios, geração de configurações, templates
- **Interface**: `Configurator` - gerencia toda a configuração do ambiente

#### **state_manager.go** - Gerenciador de Estado Atômico
- **Função**: Controla o estado do setup de forma thread-safe
- **Responsabilidades**: Operações atômicas, locks, backup/restore de estado
- **Interface**: `StateManager` - garante consistência e integridade do estado

#### **key_manager.go** - Gerenciador de Chaves Criptográficas
- **Função**: Gerencia chaves de segurança do sistema
- **Responsabilidades**: Geração, armazenamento seguro, rotação e backup de chaves
- **Interface**: `KeyManager` - centraliza toda a criptografia do sistema

#### **logger.go** - Sistema de Logging Estruturado
- **Função**: Fornece observabilidade e debugging do sistema
- **Responsabilidades**: Logging estruturado, rotação de logs, exportação
- **Interface**: `SetupLogger` - padroniza todo o logging do componente

### Nível 2: Serviços Internos (Implementação)

#### **OSValidator** - Validação por Sistema Operacional
- **Função**: Implementa validações específicas para cada SO (Windows, Linux, macOS)
- **Responsabilidades**: Detecção de SO, validação de recursos, permissões específicas
- **Implementações**: `WindowsValidator`, `LinuxValidator`, `DarwinValidator`

#### **DependencyManager** - Gerenciamento de Dependências
- **Função**: Verifica e instala dependências necessárias por SO
- **Responsabilidades**: Detecção de ferramentas, instalação automática, validação de versões
- **Integração**: Usa gerenciadores nativos (winget, apt, brew)

#### **FileSystemOps** - Operações de Sistema de Arquivos
- **Função**: Gerencia operações atômicas de arquivo e diretório
- **Responsabilidades**: Criação de estrutura, permissões, validação de integridade
- **Segurança**: Operações atômicas com locks para evitar corrupção

#### **CryptoProvider** - Provedor Criptográfico
- **Função**: Fornece serviços criptográficos seguros
- **Responsabilidades**: Geração de entropia, criptografia de chaves, derivação PBKDF2
- **Segurança**: Fonte de entropia criptograficamente segura

#### **NetworkValidator** - Validador de Rede
- **Função**: Verifica conectividade e configurações de rede
- **Responsabilidades**: Teste de conectividade, validação de firewall, proxy
- **Integração**: Testa conectividade com serviços externos

#### **ConfigGenerator** - Gerador de Configurações
- **Função**: Gera configurações a partir de templates
- **Responsabilidades**: Processamento de templates, validação de schemas, personalização
- **Templates**: YAML com variáveis dinâmicas

### Nível 3: Sistema de Estado Local (Persistência)

#### **~/.syntropy/config/** - Configurações Principais
- **Função**: Armazena configurações do sistema
- **Conteúdo**: `manager.yaml` (configuração principal), templates processados
- **Segurança**: Permissões 644, validação por schema JSON

#### **~/.syntropy/keys/** - Chaves Criptográficas
- **Função**: Armazena chaves de segurança do sistema
- **Conteúdo**: `owner.key*` (chaves criptografadas), metadados, fingerprints
- **Segurança**: Permissões 600, criptografia AES-256, backup automático

#### **~/.syntropy/nodes/** - Estrutura para Nós
- **Função**: Preparação para gerenciamento de nós da rede
- **Conteúdo**: Diretórios vazios prontos para nós futuros
- **Estrutura**: Uma pasta por nó com metadados e configurações

#### **~/.syntropy/logs/** - Logs do Sistema
- **Função**: Armazena logs estruturados do setup
- **Conteúdo**: `setup.log`, logs de validação, erros, auditoria
- **Gerenciamento**: Rotação automática, compressão, retenção configurável

#### **~/.syntropy/cache/** - Cache Temporário
- **Função**: Armazena dados temporários e cache
- **Conteúdo**: Downloads temporários, cache de validações, ISOs
- **Limpeza**: Limpeza automática, TTL configurável

#### **~/.syntropy/backups/** - Backups Automáticos
- **Função**: Armazena backups de segurança do sistema
- **Conteúdo**: Backups de configuração, chaves, estado completo
- **Gerenciamento**: Rotação automática, compressão, retenção por tempo

### Fluxo de Interação entre Níveis

```
1. setup.go (Nível 1) → coordena o processo
   ↓
2. validator.go (Nível 1) → chama OSValidator (Nível 2)
   ↓
3. OSValidator (Nível 2) → valida ~/.syntropy/ (Nível 3)
   ↓
4. configurator.go (Nível 1) → chama ConfigGenerator (Nível 2)
   ↓
5. ConfigGenerator (Nível 2) → cria arquivos em ~/.syntropy/ (Nível 3)
   ↓
6. state_manager.go (Nível 1) → persiste estado em ~/.syntropy/ (Nível 3)
```

### Princípios de Design

- **Separação de Responsabilidades**: Cada nível tem responsabilidades bem definidas
- **Desacoplamento**: Interfaces permitem troca de implementações
- **Atomicidade**: Operações críticas são atômicas e thread-safe
- **Observabilidade**: Logging estruturado em todos os níveis
- **Segurança**: Criptografia e validação em múltiplas camadas

## Estrutura de Projeto Otimizada

```
manager/interfaces/cli/setup/
├── setup.go                     # Orquestrador principal (200-300 linhas)
├── validator.go                 # Validação unificada (200-300 linhas)
├── configurator.go              # Configuração unificada (200-300 linhas)
├── state_manager.go             # Gerenciamento de estado atômico (150-200 linhas)
├── key_manager.go               # Gerenciamento de chaves criptográficas (200-250 linhas)
├── logger.go                    # Sistema de logging estruturado (100-150 linhas)
├── internal/
│   ├── types/
│   │   ├── setup.go            # Estruturas de dados principais
│   │   ├── errors.go           # Erros estruturados e códigos
│   │   └── interfaces.go       # Interfaces para desacoplamento
│   ├── services/
│   │   ├── validator/          # Serviço de validação
│   │   │   ├── os_validator.go # Validação por SO
│   │   │   ├── dependency.go   # Validação de dependências
│   │   │   └── network.go      # Validação de rede
│   │   ├── configurator/       # Serviço de configuração
│   │   │   ├── config_gen.go   # Geração de configurações
│   │   │   ├── filesystem.go   # Criação de estrutura
│   │   │   └── templates.go    # Processamento de templates
│   │   ├── keystore/           # Gerenciamento de chaves
│   │   │   ├── generator.go    # Geração de chaves
│   │   │   ├── storage.go      # Armazenamento seguro
│   │   │   └── rotation.go     # Rotação de chaves
│   │   └── state/              # Gerenciamento de estado
│   │       ├── manager.go      # Gerenciador de estado
│   │       ├── atomic.go       # Operações atômicas
│   │       └── backup.go       # Backup e recuperação
│   └── utils/
│       ├── os/                 # Utilitários por SO
│       │   ├── windows.go      # Utilitários Windows
│       │   ├── linux.go        # Utilitários Linux
│       │   └── darwin.go       # Utilitários macOS
│       ├── crypto/             # Utilitários criptográficos
│       │   ├── entropy.go      # Fonte de entropia
│       │   ├── keygen.go       # Geração de chaves
│       │   └── keystore.go     # Armazenamento seguro
│       └── filesystem/         # Operações de arquivo
│           ├── atomic.go       # Operações atômicas
│           ├── permissions.go  # Gerenciamento de permissões
│           └── validation.go   # Validação de arquivos
├── config/
│   ├── templates/
│   │   ├── manager.yaml       # Template de configuração principal
│   │   ├── security.yaml      # Template de configuração de segurança
│   │   └── network.yaml       # Template de configuração de rede
│   ├── defaults/
│   │   ├── windows.yaml       # Padrões para Windows
│   │   ├── linux.yaml         # Padrões para Linux
│   │   └── darwin.yaml        # Padrões para macOS
│   └── schemas/
│       ├── config.schema.json # Schema de configuração
│       └── environment.schema.json # Schema de ambiente
└── tests/
    ├── unit/                   # Testes unitários
    ├── integration/            # Testes de integração
    └── fixtures/               # Dados de teste
```

## Interfaces e Contratos

### Interface Principal do Setup
```go
// setup.go
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

type SetupOptions struct {
    Force           bool              `json:"force"`
    ValidateOnly    bool              `json:"validate_only"`
    Verbose         bool              `json:"verbose"`
    Quiet           bool              `json:"quiet"`
    ConfigPath      string            `json:"config_path"`
    CustomSettings  map[string]string `json:"custom_settings"`
}
```

### Interface de Validação
```go
// validator.go
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

type ValidationResult struct {
    Environment   *EnvironmentInfo   `json:"environment"`
    Dependencies  *DependencyStatus  `json:"dependencies"`
    Network       *NetworkInfo       `json:"network"`
    Permissions   *PermissionStatus  `json:"permissions"`
    Issues        []ValidationIssue  `json:"issues"`
    CanProceed    bool               `json:"can_proceed"`
    Warnings      []string           `json:"warnings"`
}
```

### Interface de Configuração
```go
// configurator.go
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

type ConfigOptions struct {
    OwnerName      string            `json:"owner_name"`
    OwnerEmail     string            `json:"owner_email"`
    NetworkConfig  *NetworkConfig    `json:"network_config"`
    SecurityConfig *SecurityConfig   `json:"security_config"`
    CustomSettings map[string]string `json:"custom_settings"`
}
```

### Interface de Gerenciamento de Estado
```go
// state_manager.go
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

type SetupState struct {
    Version        string            `json:"version"`
    CreatedAt      time.Time         `json:"created_at"`
    UpdatedAt      time.Time         `json:"updated_at"`
    Status         SetupStatus       `json:"status"`
    Environment    *EnvironmentInfo  `json:"environment"`
    Configuration  *ConfigInfo       `json:"configuration"`
    Keys           *KeyInfo          `json:"keys"`
    LastBackup     *BackupInfo       `json:"last_backup"`
    Metadata       map[string]string `json:"metadata"`
}
```

### Interface de Gerenciamento de Chaves
```go
// key_manager.go
type KeyManager interface {
    // Geração de par de chaves
    GenerateKeyPair(algorithm string) (*KeyPair, error)
    
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

type KeyPair struct {
    ID           string    `json:"id"`
    Algorithm    string    `json:"algorithm"`
    PrivateKey   []byte    `json:"private_key"`
    PublicKey    []byte    `json:"public_key"`
    CreatedAt    time.Time `json:"created_at"`
    ExpiresAt    time.Time `json:"expires_at"`
    Fingerprint  string    `json:"fingerprint"`
    Metadata     map[string]string `json:"metadata"`
}
```

## Implementação por Sistema Operacional

### Interface de Validação por SO
```go
// internal/services/validator/os_validator.go
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

type OSInfo struct {
    Name         string `json:"name"`
    Version      string `json:"version"`
    Architecture string `json:"architecture"`
    Build        string `json:"build"`
    Kernel       string `json:"kernel"`
}

// Implementações específicas por SO
type WindowsValidator struct {
    logger *logrus.Logger
}

type LinuxValidator struct {
    logger *logrus.Logger
}

type DarwinValidator struct {
    logger *logrus.Logger
}

func NewOSValidator(logger *logrus.Logger) OSValidator {
    switch runtime.GOOS {
    case "windows":
        return &WindowsValidator{logger: logger}
    case "linux":
        return &LinuxValidator{logger: logger}
    case "darwin":
        return &DarwinValidator{logger: logger}
    default:
        return nil
    }
}
```

## Sistema de Logging Estruturado

### Interface de Logging
```go
// logger.go
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
}

type LogEntry struct {
    Timestamp time.Time              `json:"timestamp"`
    Level     string                 `json:"level"`
    Message   string                 `json:"message"`
    Step      string                 `json:"step,omitempty"`
    Data      map[string]interface{} `json:"data,omitempty"`
    Error     string                 `json:"error,omitempty"`
}
```

## Sistema de Erros Estruturado

### Códigos de Erro e Contexto
```go
// internal/types/errors.go
type SetupError struct {
    Code        string                 `json:"code"`
    Message     string                 `json:"message"`
    Context     map[string]interface{} `json:"context"`
    Suggestions []string               `json:"suggestions"`
    Timestamp   time.Time              `json:"timestamp"`
    Cause       error                  `json:"-"`
}

// Códigos de erro específicos
const (
    ErrOSNotSupported     = "SETUP_001"
    ErrInsufficientPerms  = "SETUP_002"
    ErrMissingDependency  = "SETUP_003"
    ErrInsufficientSpace  = "SETUP_004"
    ErrKeyGeneration      = "SETUP_005"
    ErrNetworkConnectivity = "SETUP_006"
    ErrConfigCorrupted    = "SETUP_007"
    ErrStateCorrupted     = "SETUP_008"
    ErrBackupFailed       = "SETUP_009"
    ErrRestoreFailed      = "SETUP_010"
)

func (e *SetupError) Error() string {
    return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *SetupError) Unwrap() error {
    return e.Cause
}
```

## Comandos CLI Simplificados

### Estrutura de Comandos
```bash
# Comando principal
syntropy setup [OPTIONS]

# Opções disponíveis
--validate-only     # Apenas validar, não configurar
--force            # Forçar setup mesmo com warnings
--verbose          # Saída detalhada
--quiet            # Saída silenciosa
--config PATH      # Caminho para arquivo de configuração
--backup PATH      # Caminho para backup
--restore PATH     # Restaurar de backup

# Subcomandos
syntropy setup status    # Status do setup
syntropy setup reset     # Reset completo
syntropy setup repair    # Reparo automático
syntropy setup backup    # Backup manual
syntropy setup restore   # Restauração
```

### Exemplo de Uso
```bash
# Setup completo
syntropy setup

# Setup com validação detalhada
syntropy setup --verbose

# Apenas validar ambiente
syntropy setup --validate-only

# Setup forçado
syntropy setup --force

# Verificar status
syntropy setup status

# Reparar problemas
syntropy setup repair

# Reset completo
syntropy setup reset --confirm
```

## Fluxo de Implementação

### 1. Implementação do Orquestrador Principal
```go
// setup.go
type SetupManager struct {
    validator    Validator
    configurator Configurator
    stateManager StateManager
    keyManager   KeyManager
    logger       SetupLogger
}

func NewSetupManager() (*SetupManager, error) {
    logger := NewSetupLogger()
    
    return &SetupManager{
        validator:    NewValidator(logger),
        configurator: NewConfigurator(logger),
        stateManager: NewStateManager(logger),
        keyManager:   NewKeyManager(logger),
        logger:       logger,
    }, nil
}

func (sm *SetupManager) Setup(options *SetupOptions) error {
    sm.logger.LogStep("setup_start", map[string]interface{}{
        "options": options,
    })
    
    // 1. Validar ambiente
    result, err := sm.validator.ValidateEnvironment()
    if err != nil {
        return sm.handleError(err, "validation_failed")
    }
    
    if !result.CanProceed && !options.Force {
        return sm.handleError(ErrValidationFailed, "validation_failed")
    }
    
    // 2. Criar estrutura de diretórios
    if err := sm.configurator.CreateStructure(); err != nil {
        return sm.handleError(err, "structure_creation_failed")
    }
    
    // 3. Gerar chaves
    keyPair, err := sm.keyManager.GenerateKeyPair("ed25519")
    if err != nil {
        return sm.handleError(err, "key_generation_failed")
    }
    
    // 4. Gerar configuração
    if err := sm.configurator.GenerateConfig(&ConfigOptions{
        OwnerName: options.CustomSettings["owner_name"],
        OwnerEmail: options.CustomSettings["owner_email"],
    }); err != nil {
        return sm.handleError(err, "config_generation_failed")
    }
    
    // 5. Salvar estado
    state := &SetupState{
        Version: "1.0.0",
        CreatedAt: time.Now(),
        Status: SetupStatusCompleted,
        Keys: &KeyInfo{
            OwnerKeyID: keyPair.ID,
            Algorithm: keyPair.Algorithm,
        },
    }
    
    if err := sm.stateManager.SaveState(state); err != nil {
        return sm.handleError(err, "state_save_failed")
    }
    
    sm.logger.LogStep("setup_completed", map[string]interface{}{
        "key_id": keyPair.ID,
    })
    
    return nil
}
```

### 2. Implementação do Validador
```go
// validator.go
type Validator struct {
    osValidator OSValidator
    logger      SetupLogger
}

func NewValidator(logger SetupLogger) *Validator {
    return &Validator{
        osValidator: NewOSValidator(logger),
        logger:      logger,
    }
}

func (v *Validator) ValidateEnvironment() (*ValidationResult, error) {
    v.logger.LogStep("validation_start", nil)
    
    // Validar SO
    osInfo, err := v.osValidator.DetectOS()
    if err != nil {
        return nil, err
    }
    
    // Validar recursos
    resources, err := v.osValidator.ValidateResources()
    if err != nil {
        return nil, err
    }
    
    // Validar permissões
    permissions, err := v.osValidator.ValidatePermissions()
    if err != nil {
        return nil, err
    }
    
    // Validar dependências
    dependencies, err := v.validateDependencies()
    if err != nil {
        return nil, err
    }
    
    // Validar rede
    network, err := v.validateNetwork()
    if err != nil {
        return nil, err
    }
    
    result := &ValidationResult{
        Environment: &EnvironmentInfo{
            OS: osInfo,
            Resources: resources,
            Permissions: permissions,
        },
        Dependencies: dependencies,
        Network: network,
        CanProceed: true,
    }
    
    v.logger.LogStep("validation_completed", map[string]interface{}{
        "can_proceed": result.CanProceed,
        "issues_count": len(result.Issues),
    })
    
    return result, nil
}
```

### 3. Implementação do Configurador
```go
// configurator.go
type Configurator struct {
    logger SetupLogger
}

func NewConfigurator(logger SetupLogger) *Configurator {
    return &Configurator{
        logger: logger,
    }
}

func (c *Configurator) CreateStructure() error {
    c.logger.LogStep("structure_creation_start", nil)
    
    baseDir := filepath.Join(os.Getenv("HOME"), ".syntropy")
    
    directories := []string{
        filepath.Join(baseDir, "config"),
        filepath.Join(baseDir, "keys"),
        filepath.Join(baseDir, "nodes"),
        filepath.Join(baseDir, "logs"),
        filepath.Join(baseDir, "cache"),
        filepath.Join(baseDir, "backups"),
    }
    
    for _, dir := range directories {
        if err := os.MkdirAll(dir, 0755); err != nil {
            return fmt.Errorf("failed to create directory %s: %w", dir, err)
        }
    }
    
    c.logger.LogStep("structure_creation_completed", map[string]interface{}{
        "base_dir": baseDir,
        "directories": directories,
    })
    
    return nil
}

func (c *Configurator) GenerateConfig(options *ConfigOptions) error {
    c.logger.LogStep("config_generation_start", map[string]interface{}{
        "options": options,
    })
    
    config := &ManagerConfig{
        Version: "1.0.0",
        Owner: OwnerConfig{
            Name:  options.OwnerName,
            Email: options.OwnerEmail,
        },
        Network: options.NetworkConfig,
        Security: options.SecurityConfig,
        CreatedAt: time.Now(),
    }
    
    configPath := filepath.Join(os.Getenv("HOME"), ".syntropy", "config", "manager.yaml")
    
    data, err := yaml.Marshal(config)
    if err != nil {
        return fmt.Errorf("failed to marshal config: %w", err)
    }
    
    if err := os.WriteFile(configPath, data, 0644); err != nil {
        return fmt.Errorf("failed to write config file: %w", err)
    }
    
    c.logger.LogStep("config_generation_completed", map[string]interface{}{
        "config_path": configPath,
    })
    
    return nil
}
```

## Sistema de Segurança

### Geração Segura de Chaves
```go
// internal/services/keystore/generator.go
type KeyGenerator struct {
    entropySource *crypto.EntropySource
    logger        SetupLogger
}

func NewKeyGenerator(logger SetupLogger) *KeyGenerator {
    return &KeyGenerator{
        entropySource: crypto.NewEntropySource(),
        logger:        logger,
    }
}

func (kg *KeyGenerator) GenerateEd25519KeyPair() (*KeyPair, error) {
    kg.logger.LogStep("key_generation_start", map[string]interface{}{
        "algorithm": "ed25519",
    })
    
    // Gerar chave privada usando fonte de entropia segura
    privateKey, err := kg.entropySource.GeneratePrivateKey(ed25519.PrivateKeySize)
    if err != nil {
        return nil, fmt.Errorf("failed to generate private key: %w", err)
    }
    
    // Gerar chave pública
    publicKey := privateKey.Public().(ed25519.PublicKey)
    
    // Criar fingerprint
    fingerprint := kg.generateFingerprint(publicKey)
    
    keyPair := &KeyPair{
        ID:          kg.generateKeyID(),
        Algorithm:   "ed25519",
        PrivateKey:  privateKey,
        PublicKey:   publicKey,
        CreatedAt:   time.Now(),
        ExpiresAt:   time.Now().AddDate(1, 0, 0), // 1 ano
        Fingerprint: fingerprint,
        Metadata: map[string]string{
            "generated_by": "syntropy-setup",
            "version": "1.0.0",
        },
    }
    
    kg.logger.LogStep("key_generation_completed", map[string]interface{}{
        "key_id": keyPair.ID,
        "fingerprint": keyPair.Fingerprint,
    })
    
    return keyPair, nil
}
```

### Armazenamento Seguro
```go
// internal/services/keystore/storage.go
type SecureKeyStorage struct {
    keystorePath string
    logger       SetupLogger
}

func NewSecureKeyStorage(logger SetupLogger) *SecureKeyStorage {
    return &SecureKeyStorage{
        keystorePath: filepath.Join(os.Getenv("HOME"), ".syntropy", "keys"),
        logger:       logger,
    }
}

func (sks *SecureKeyStorage) StoreKeyPair(keyPair *KeyPair, passphrase string) error {
    sks.logger.LogStep("key_storage_start", map[string]interface{}{
        "key_id": keyPair.ID,
    })
    
    // Criptografar chave privada
    encryptedPrivateKey, err := sks.encryptPrivateKey(keyPair.PrivateKey, passphrase)
    if err != nil {
        return fmt.Errorf("failed to encrypt private key: %w", err)
    }
    
    // Salvar chave privada criptografada
    privateKeyPath := filepath.Join(sks.keystorePath, fmt.Sprintf("%s.key", keyPair.ID))
    if err := os.WriteFile(privateKeyPath, encryptedPrivateKey, 0600); err != nil {
        return fmt.Errorf("failed to write private key: %w", err)
    }
    
    // Salvar chave pública
    publicKeyPath := filepath.Join(sks.keystorePath, fmt.Sprintf("%s.key.pub", keyPair.ID))
    if err := os.WriteFile(publicKeyPath, keyPair.PublicKey, 0644); err != nil {
        return fmt.Errorf("failed to write public key: %w", err)
    }
    
    // Salvar metadados
    metadataPath := filepath.Join(sks.keystorePath, fmt.Sprintf("%s.meta", keyPair.ID))
    metadata, err := json.Marshal(keyPair.Metadata)
    if err != nil {
        return fmt.Errorf("failed to marshal metadata: %w", err)
    }
    
    if err := os.WriteFile(metadataPath, metadata, 0644); err != nil {
        return fmt.Errorf("failed to write metadata: %w", err)
    }
    
    sks.logger.LogStep("key_storage_completed", map[string]interface{}{
        "key_id": keyPair.ID,
        "private_key_path": privateKeyPath,
        "public_key_path": publicKeyPath,
    })
    
    return nil
}
```

## Sistema de Backup e Recuperação

### Backup Automático
```go
// internal/services/state/backup.go
type BackupManager struct {
    backupPath string
    logger     SetupLogger
}

func NewBackupManager(logger SetupLogger) *BackupManager {
    return &BackupManager{
        backupPath: filepath.Join(os.Getenv("HOME"), ".syntropy", "backups"),
        logger:     logger,
    }
}

func (bm *BackupManager) CreateBackup(name string) error {
    bm.logger.LogStep("backup_start", map[string]interface{}{
        "backup_name": name,
    })
    
    timestamp := time.Now().Format("20060102_150405")
    backupName := fmt.Sprintf("%s_%s.tar.gz", name, timestamp)
    backupPath := filepath.Join(bm.backupPath, backupName)
    
    // Criar arquivo de backup
    file, err := os.Create(backupPath)
    if err != nil {
        return fmt.Errorf("failed to create backup file: %w", err)
    }
    defer file.Close()
    
    // Criar compressor
    gzWriter := gzip.NewWriter(file)
    defer gzWriter.Close()
    
    tarWriter := tar.NewWriter(gzWriter)
    defer tarWriter.Close()
    
    // Adicionar arquivos ao backup
    baseDir := filepath.Join(os.Getenv("HOME"), ".syntropy")
    if err := bm.addDirectoryToTar(tarWriter, baseDir, ""); err != nil {
        return fmt.Errorf("failed to add directory to tar: %w", err)
    }
    
    bm.logger.LogStep("backup_completed", map[string]interface{}{
        "backup_path": backupPath,
        "backup_name": backupName,
    })
    
    return nil
}
```

## Testes e Validação

### Estrutura de Testes
```go
// tests/unit/setup_test.go
func TestSetupManager_Setup(t *testing.T) {
    // Setup do teste
    logger := NewTestLogger()
    manager := NewSetupManager(logger)
    
    // Teste de setup completo
    options := &SetupOptions{
        Force: false,
        Verbose: true,
        CustomSettings: map[string]string{
            "owner_name": "Test User",
            "owner_email": "test@example.com",
        },
    }
    
    err := manager.Setup(options)
    assert.NoError(t, err)
    
    // Verificar se o estado foi salvo
    state, err := manager.stateManager.LoadState()
    assert.NoError(t, err)
    assert.Equal(t, SetupStatusCompleted, state.Status)
    
    // Verificar se as chaves foram geradas
    assert.NotEmpty(t, state.Keys.OwnerKeyID)
    
    // Verificar se a configuração foi criada
    configPath := filepath.Join(os.Getenv("HOME"), ".syntropy", "config", "manager.yaml")
    assert.FileExists(t, configPath)
}
```

## Configurações e Templates

### Template de Configuração Principal
```yaml
# config/templates/manager.yaml
version: "1.0.0"
owner:
  name: "{{.OwnerName}}"
  email: "{{.OwnerEmail}}"
  key_id: "{{.OwnerKeyID}}"
network:
  discovery:
    enabled: true
    port: 8080
  api:
    enabled: true
    port: 9090
    tls:
      enabled: true
      cert_path: "{{.CertPath}}"
      key_path: "{{.KeyPath}}"
security:
  encryption:
    algorithm: "ed25519"
    key_rotation:
      enabled: true
      interval: "30d"
  backup:
    enabled: true
    interval: "24h"
    retention: "30d"
logging:
  level: "info"
  format: "json"
  output: "file"
  file_path: "{{.LogPath}}"
  rotation:
    enabled: true
    max_size: "100MB"
    max_age: "7d"
    max_backups: 10
```

### Schema de Validação
```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["version", "owner", "network", "security"],
  "properties": {
    "version": {
      "type": "string",
      "pattern": "^[0-9]+\\.[0-9]+\\.[0-9]+$"
    },
    "owner": {
      "type": "object",
      "required": ["name", "email", "key_id"],
      "properties": {
        "name": {"type": "string", "minLength": 1},
        "email": {"type": "string", "format": "email"},
        "key_id": {"type": "string", "minLength": 1}
      }
    },
    "network": {
      "type": "object",
      "required": ["discovery", "api"],
      "properties": {
        "discovery": {
          "type": "object",
          "properties": {
            "enabled": {"type": "boolean"},
            "port": {"type": "integer", "minimum": 1, "maximum": 65535}
          }
        },
        "api": {
          "type": "object",
          "properties": {
            "enabled": {"type": "boolean"},
            "port": {"type": "integer", "minimum": 1, "maximum": 65535},
            "tls": {
              "type": "object",
              "properties": {
                "enabled": {"type": "boolean"},
                "cert_path": {"type": "string"},
                "key_path": {"type": "string"}
              }
            }
          }
        }
      }
    }
  }
}
```

## Considerações de Implementação

### 1. Tratamento de Erros Robusto
- Todos os erros devem ser estruturados com código, contexto e sugestões
- Implementar retry automático para operações de rede
- Logging detalhado de todos os erros para debugging

### 2. Operações Atômicas
- Todas as operações de estado devem ser atômicas
- Usar locks de arquivo para evitar corrupção
- Implementar rollback automático em caso de falha

### 3. Segurança
- Usar fontes de entropia criptograficamente seguras
- Implementar derivação de chaves PBKDF2
- Armazenar chaves privadas criptografadas
- Implementar rotação automática de chaves

### 4. Performance
- Operações de I/O assíncronas onde possível
- Cache de validações para evitar verificações repetidas
- Compressão de backups para economizar espaço

### 5. Observabilidade
- Logging estruturado com níveis apropriados
- Métricas de performance e uso
- Rastreamento de operações com correlation IDs

### 6. Testabilidade
- Interfaces bem definidas para mocking
- Testes unitários com cobertura > 80%
- Testes de integração para fluxos completos
- Testes de regressão para diferentes SOs

## Conclusão

Este guia fornece uma arquitetura simplificada e bem estruturada para o componente setup, focando em:

1. **Simplicidade**: Arquitetura direta sem over-engineering
2. **Robustez**: Tratamento de erros e operações atômicas
3. **Segurança**: Sistema de chaves criptográficas robusto
4. **Manutenibilidade**: Código limpo e bem testado
5. **Observabilidade**: Logging estruturado e métricas

A implementação deve seguir os padrões Go estabelecidos e garantir que o componente seja confiável, seguro e fácil de manter.
