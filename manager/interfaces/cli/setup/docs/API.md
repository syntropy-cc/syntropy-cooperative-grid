# Setup Component - API Documentation

## API Overview

[MACRO_VIEW]
A API do componente Setup fornece uma interface unificada para configuração e gerenciamento do ambiente Syntropy CLI, integrando-se com o ecossistema de APIs do projeto Syntropy Cooperative Grid através de contratos bem definidos.
[/MACRO_VIEW]

[MESO_VIEW]
Esta API se relaciona com as APIs de nível de módulo através de interfaces padronizadas, mantendo compatibilidade com o sistema de gerenciamento de estado centralizado e seguindo os padrões de versionamento do projeto.
[/MESO_VIEW]

[MICRO_VIEW]
A API específica expõe capacidades de setup, validação, gerenciamento de estado e operações de manutenção através de métodos bem tipados e tratamento de erros consistente.
[/MICRO_VIEW]

## API Principles
- **Consistency**: Todos os métodos seguem o padrão de nomenclatura verb-noun
- **Idempotency**: Operações são seguras para retry e execução múltipla
- **Versioning**: Compatibilidade retroativa garantida dentro de versões major
- **Error Handling**: Tratamento consistente de erros com códigos estruturados
- **Logging**: Todas as operações são registradas com contexto estruturado

## Authentication & Authorization

### Authentication Methods

| Method | Use Case | Example |
|--------|----------|---------|
| None | Operações locais de setup | Todas as operações da API |

### Required Permissions

| Endpoint/Method | Permission Level | Scope |
|-----------------|------------------|--------|
| Setup | File System Write | Diretório home do usuário |
| Reset | File System Write | Diretório de configuração |
| Validate | File System Read | Sistema de arquivos local |

## Core API Reference

### Setup Operations

#### Method: NewSetupManager
**Purpose**: Cria uma nova instância do gerenciador de setup

**Signature**:
```go
func NewSetupManager() (*SetupManager, error)
```

**Parameters**: Nenhum

**Returns**:
```go
*SetupManager // Instância configurada do gerenciador
error         // Erro, se houver
```

**Errors**:
| Error Code | Condition | Resolution |
|------------|-----------|------------|
| ErrManagerCreation | Falha na inicialização | Verificar dependências e permissões |

**Example**:
```go
// Request
manager, err := setup.NewSetupManager()
if err != nil {
    log.Fatal(err)
}

// Response - Success
// manager é uma instância válida de SetupManager

// Response - Error
// err contém detalhes do erro de inicialização
```

**Notes**:
- Deve ser chamado antes de qualquer operação de setup
- Inicializa todos os subcomponentes necessários
- Thread-safe para uso concorrente

#### Method: Setup
**Purpose**: Executa o processo completo de setup do ambiente Syntropy

**Signature**:
```go
func (sm *SetupManager) Setup(options *types.SetupOptions) error
```

**Parameters**:
| Parameter | Type | Required | Default | Description | Constraints |
|-----------|------|----------|---------|-------------|-------------|
| options | *types.SetupOptions | Yes | nil | Opções de configuração do setup | Não pode ser nil |

**Returns**:
```go
error // Erro, se houver
```

**Errors**:
| Error Code | Condition | Resolution |
|------------|-----------|------------|
| ErrValidationFailed | Validação do ambiente falhou | Verificar requisitos do sistema |
| ErrStructureCreation | Falha na criação de diretórios | Verificar permissões de escrita |
| ErrKeyGeneration | Falha na geração de chaves | Verificar espaço em disco |
| ErrConfigGeneration | Falha na geração de configuração | Verificar permissões |

**Example**:
```go
// Request
options := &types.SetupOptions{
    Force:        false,
    Verbose:      true,
    CustomSettings: map[string]string{
        "owner_name":  "João Silva",
        "owner_email": "joao@syntropy.network",
    },
}

err := manager.Setup(options)

// Response - Success
// err == nil, setup concluído

// Response - Error
// err contém detalhes do erro específico
```

**Notes**:
- Operação atômica - ou completa com sucesso ou falha completamente
- Cria backup automático de configurações existentes
- Registra todas as operações no log estruturado

#### Method: SetupWithPublicOptions
**Purpose**: Executa setup usando opções públicas simplificadas

**Signature**:
```go
func (sm *SetupManager) SetupWithPublicOptions(options *SetupOptions) error
```

**Parameters**:
| Parameter | Type | Required | Default | Description | Constraints |
|-----------|------|----------|---------|-------------|-------------|
| options | *SetupOptions | Yes | nil | Opções públicas de setup | Não pode ser nil |

**Returns**:
```go
error // Erro, se houver
```

**Example**:
```go
// Request
options := &SetupOptions{
    Force:          false,
    ValidateOnly:   false,
    Verbose:        true,
    ConfigPath:     "/custom/path/config.yaml",
    CustomSettings: map[string]string{
        "owner_name": "Maria Santos",
    },
}

err := manager.SetupWithPublicOptions(options)
```

### Validation Operations

#### Method: Validate
**Purpose**: Valida o ambiente sem executar o setup

**Signature**:
```go
func (sm *SetupManager) Validate() (*types.ValidationResult, error)
```

**Parameters**: Nenhum

**Returns**:
```go
*types.ValidationResult // Resultado detalhado da validação
error                   // Erro, se houver
```

**Errors**:
| Error Code | Condition | Resolution |
|------------|-----------|------------|
| ErrValidationFailed | Validação crítica falhou | Corrigir problemas identificados |

**Example**:
```go
// Request
result, err := manager.Validate()

// Response - Success
// result.Environment.OS == "linux"
// result.CanProceed == true
// result.Issues == [] (lista vazia se sem problemas)

// Response - Error
// err contém detalhes do erro de validação
```

### Status Operations

#### Method: Status
**Purpose**: Verifica o status atual do setup

**Signature**:
```go
func (sm *SetupManager) Status() (*types.SetupStatus, error)
```

**Parameters**: Nenhum

**Returns**:
```go
*types.SetupStatus // Status atual do setup
error              // Erro, se houver
```

**Errors**:
| Error Code | Condition | Resolution |
|------------|-----------|------------|
| ErrStateLoad | Estado não encontrado ou corrompido | Executar setup ou reset |

**Example**:
```go
// Request
status, err := manager.Status()

// Response - Success
// status == SetupStatusCompleted

// Response - Error
// err indica que setup não foi executado
```

### Maintenance Operations

#### Method: Reset
**Purpose**: Remove todas as configurações e estado do setup

**Signature**:
```go
func (sm *SetupManager) Reset(confirm bool) error
```

**Parameters**:
| Parameter | Type | Required | Default | Description | Constraints |
|-----------|------|----------|---------|-------------|-------------|
| confirm | bool | Yes | - | Confirmação de reset | Deve ser true para executar |

**Returns**:
```go
error // Erro, se houver
```

**Errors**:
| Error Code | Condition | Resolution |
|------------|-----------|------------|
| ErrResetConfirmation | Confirmação não fornecida | Chamar com confirm=true |

**Example**:
```go
// Request
err := manager.Reset(true)

// Response - Success
// err == nil, reset concluído

// Response - Error
// err contém detalhes do erro de reset
```

**Notes**:
- Operação irreversível
- Remove todos os arquivos de configuração e chaves
- Requer confirmação explícita

#### Method: Repair
**Purpose**: Repara problemas detectados no setup

**Signature**:
```go
func (sm *SetupManager) Repair() error
```

**Parameters**: Nenhum

**Returns**:
```go
error // Erro, se houver
```

**Example**:
```go
// Request
err := manager.Repair()

// Response - Success
// err == nil, reparo concluído

// Response - Error
// err contém detalhes do erro de reparo
```

### Legacy API Operations

#### Method: SetupLegacy
**Purpose**: Executa setup usando interface legacy para compatibilidade

**Signature**:
```go
func SetupLegacy(options LegacySetupOptions) (*LegacySetupResult, error)
```

**Parameters**:
| Parameter | Type | Required | Default | Description | Constraints |
|-----------|------|----------|---------|-------------|-------------|
| options | LegacySetupOptions | Yes | - | Opções legacy de setup | - |

**Returns**:
```go
*LegacySetupResult // Resultado em formato legacy
error              // Erro, se houver
```

**Example**:
```go
// Request
options := LegacySetupOptions{
    Force:          false,
    InstallService: true,
    ConfigPath:     "/custom/config.yaml",
}

result, err := setup.SetupLegacy(options)

// Response - Success
// result.Success == true
// result.ConfigPath contém caminho da configuração

// Response - Error
// result.Success == false
// result.Error contém detalhes do erro
```

## Event API / Callbacks

### Event: SetupStarted
**Triggered When**: Processo de setup é iniciado

**Payload**:
```go
{
    EventType: "setup_started",
    Timestamp: time.Time,
    Data: {
        Options: *types.SetupOptions,
        CorrelationID: string,
    }
}
```

**Subscribe Example**:
```go
// Logging automático através do SetupLogger
manager.logger.LogStep("setup_start", map[string]interface{}{
    "options": options,
})
```

### Event: SetupCompleted
**Triggered When**: Processo de setup é concluído com sucesso

**Payload**:
```go
{
    EventType: "setup_completed",
    Timestamp: time.Time,
    Data: {
        State: *types.SetupState,
        Duration: time.Duration,
        CorrelationID: string,
    }
}
```

### Event: SetupFailed
**Triggered When**: Processo de setup falha

**Payload**:
```go
{
    EventType: "setup_failed",
    Timestamp: time.Time,
    Data: {
        Error: error,
        Step: string,
        CorrelationID: string,
    }
}
```

## Rate Limiting

| Endpoint Category | Requests/Minute | Burst Limit | Retry-After |
|-------------------|-----------------|-------------|-------------|
| Setup Operations | 10 | 2 | 60 seconds |
| Validation | 60 | 10 | 10 seconds |
| Status/Repair | 120 | 20 | 5 seconds |

## API Versioning

### Version Strategy
**Current Version**: 1.0.0
**Supported Versions**: 1.0.x
**Deprecation Policy**: 6 meses de aviso para mudanças breaking

### Version Differences

| Feature | v1.0 | v1.1 (planned) | v2.0 (planned) |
|---------|------|----------------|----------------|
| Basic Setup | ✅ | ✅ | ✅ |
| Legacy API | ✅ | ✅ | ❌ Deprecated |
| Advanced Validation | ❌ | ✅ | ✅ |
| Plugin Support | ❌ | ❌ | ✅ |

## SDK/Client Libraries

### Official Libraries

| Language | Package | Version | Documentation |
|----------|---------|---------|---------------|
| Go | setup-component | 1.0.0 | [API.md](./API.md) |

### Code Generation
```bash
# Gerar clientes a partir da especificação da API
go generate ./...
```

## Testing the API

### Test Environment
**Base Path**: `setup-component/src`
**Test Credentials**: Nenhuma autenticação necessária
**Limitations**: Requer permissões de arquivo para operações de setup

### Example Test Flow
```go
func TestSetupFlow(t *testing.T) {
    // 1. Setup - Create manager
    manager, err := NewSetupManager()
    require.NoError(t, err)
    
    // 2. Execute - Run setup
    options := &types.SetupOptions{
        Force: true,
        TestMode: true,
    }
    err = manager.Setup(options)
    require.NoError(t, err)
    
    // 3. Verify - Check status
    status, err := manager.Status()
    require.NoError(t, err)
    assert.Equal(t, types.SetupStatusCompleted, *status)
}
```

## Migration Guides

### Migrating from v0.x to v1.0
#### Breaking Changes
- `SetupOptions` struct redesenhada: Migrar para nova estrutura
- Error handling: Usar novos tipos de erro tipados

#### Deprecated Features

| Feature | Deprecated In | Removed In | Alternative |
|---------|---------------|------------|-------------|
| LegacySetupOptions | 1.0.0 | 2.0.0 | SetupOptions |
| Simple Setup API | 1.0.0 | 2.0.0 | SetupWithPublicOptions |

## API Best Practices

### DO
- Sempre verificar erros retornados pelos métodos
- Usar TestMode=true durante desenvolvimento
- Configurar logging verbose para debugging
- Executar Validate() antes de Setup() em produção

### DON'T
- Não chamar Setup() múltiplas vezes sem Reset()
- Não ignorar erros de validação
- Não executar Reset() sem confirmação explícita
- Não misturar APIs legacy com novas

## Troubleshooting

### Common Integration Issues

| Issue | Symptoms | Solution |
|-------|----------|----------|
| Permission Denied | Setup falha na criação de diretórios | Executar com privilégios de administrador |
| Disk Space | Key generation falha | Liberar espaço em disco ou especificar diretório alternativo |
| Network Issues | Validation falha em checks de conectividade | Verificar conectividade ou usar Force=true |
| State Corruption | Status() retorna erro | Executar Reset() seguido de novo Setup() |

## API Metrics

### SLA

| Metric | Target | Measurement |
|--------|--------|-------------|
| Availability | 99.9% | Uptime do componente |
| Response Time (p95) | <5s | Tempo de execução do Setup() |
| Success Rate | >95% | Taxa de sucesso de operações |

## Changelog

### 1.0.0 - 2024-01-15
#### Added
- Setup completo do ambiente Syntropy
- Validação de ambiente e dependências
- Gerenciamento de chaves criptográficas
- Sistema de logging estruturado

#### Changed
- Migração de API legacy para nova arquitetura
- Melhoria no tratamento de erros

#### Deprecated
- Nenhuma funcionalidade deprecated nesta versão

#### Removed
- Nenhuma funcionalidade removida nesta versão

## Educational Resources
Para compreensão conceitual e exercícios de aprendizado, veja [LEARN.md](./LEARN.md).






