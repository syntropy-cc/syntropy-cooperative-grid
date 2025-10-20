# Estratégia de Testes - Node Component

## Visão Geral

Esta estratégia de testes foi desenvolvida para garantir **100% de cobertura** do código fonte do Node Component, seguindo rigorosamente as regras do `test-suite.rules.md`.

### Filosofia de Testes

- **Test Pyramid**: 70% unitários, 25% integração, 5% end-to-end
- **Zero Duplication**: Nenhum utilitário de `src/internal/` será duplicado em `tests/`
- **Immutability**: Nenhum arquivo em `src/` será modificado durante testes
- **Platform Coverage**: Testes específicos para Linux, Windows e macOS
- **Error Coverage**: 100% dos caminhos de erro testados

## Estrutura de Testes

### Testes Unitários (70% - ~4200 linhas de teste)

#### 1. Testes de Componentes Principais

**tests/unit/node_test.go**
- Testa `src/node.go` (NodeManager)
- Cobertura: 100% de todas as funções e métodos
- Cenários: criação, listagem, status, logs, remoção de nós
- Mocks: EventBus, NodeStateManager, Logger, Configuration

**tests/unit/create_test.go**
- Testa `src/create.go` (CreateSubcomponent)
- Cobertura: 100% do fluxo de criação
- Cenários: criação normal, interativa, com erros
- Mocks: todos os subcomponentes (ConfigGenerator, USBDetector, etc.)

**tests/unit/handshake_test.go**
- Testa `src/handshake.go` (HandshakeManager)
- Cobertura: 100% do protocolo de handshake
- Cenários: handshake válido, inválido, timeout, falhas
- Mocks: TokenIntegration, NodeStateManager, CertificateManager

**tests/unit/listener_test.go**
- Testa `src/listener.go` (Listener)
- Cobertura: 100% do gerenciamento de conexões
- Cenários: start/stop, conexões, timeouts, estatísticas
- Mocks: HandshakeManager, Logger

**tests/unit/heartbeat_test.go**
- Testa `src/heartbeat.go` (HeartbeatManager)
- Cobertura: 100% do sistema de heartbeat
- Cenários: heartbeat normal, falhas, reconexão, métricas
- Mocks: NodeStateManager, Logger

**tests/unit/node_state_test.go**
- Testa `src/node_state.go` (NodeStateManager)
- Cobertura: 100% do gerenciamento de estado
- Cenários: transições de estado, persistência, recuperação
- Mocks: Logger

#### 2. Testes de Componentes de Suporte

**tests/unit/auto_config_generator_test.go**
- Testa `src/auto_config_generator.go`
- Cobertura: 100% da geração de configurações
- Cenários: NodeID, SSH keys, certificados, IP detection
- Mocks: TokenIntegration, Logger

**tests/unit/cloud_init_generator_test.go**
- Testa `src/cloud_init_generator.go`
- Cobertura: 100% da geração de cloud-init
- Cenários: templates, validação, estatísticas
- Mocks: Logger

**tests/unit/usb_detector_test.go**
- Testa `src/usb_detector.go` (interface base)
- Cobertura: 100% da interface USBDetector
- Cenários: factory pattern, validação de interface
- Mocks: sistema operacional

**tests/unit/usb_detector_linux_test.go**
- Testa `src/usb_detector_linux.go`
- Cobertura: 100% da implementação Linux
- Cenários: lsblk, fdisk, udevadm, permissões
- Build constraint: `// +build linux`

**tests/unit/usb_detector_windows_test.go**
- Testa `src/usb_detector_windows.go`
- Cobertura: 100% da implementação Windows
- Cenários: WMIC, PowerShell, permissões
- Build constraint: `// +build windows`

**tests/unit/iso_downloader_test.go**
- Testa `src/iso_downloader.go`
- Cobertura: 100% do download e cache de ISOs
- Cenários: download, cache, verificação SHA256, progress
- Mocks: HTTP client, filesystem

**tests/unit/usb_writer_test.go**
- Testa `src/usb_writer.go` (interface base)
- Cobertura: 100% da interface USBWriter
- Cenários: factory pattern, validação de interface
- Mocks: sistema operacional

**tests/unit/usb_writer_linux_test.go**
- Testa `src/usb_writer_linux.go`
- Cobertura: 100% da implementação Linux
- Cenários: dd, sync, umount, cloud-init injection
- Build constraint: `// +build linux`

**tests/unit/usb_writer_windows_test.go**
- Testa `src/usb_writer_windows.go`
- Cobertura: 100% da implementação Windows
- Cenários: diskpart, PowerShell, cloud-init injection
- Build constraint: `// +build windows`

**tests/unit/usb_writer_macos_test.go**
- Testa `src/usb_writer_macos.go`
- Cobertura: 100% da implementação macOS
- Cenários: diskutil, dd, cloud-init injection
- Build constraint: `// +build darwin`

**tests/unit/token_integration_test.go**
- Testa `src/token_integration.go`
- Cobertura: 100% da integração com tokens
- Cenários: get, validate, refresh tokens
- Mocks: Setup TokenManager

**tests/unit/cli_test.go**
- Testa `src/cli.go` (CLICommands)
- Cobertura: 100% dos comandos CLI
- Cenários: todos os comandos, flags, output formatting
- Mocks: NodeManager

**tests/unit/implementations_test.go**
- Testa `src/implementations.go`
- Cobertura: 100% das implementações básicas
- Cenários: Loggers, EventBus, Configuration
- Mocks: sistema operacional

**tests/unit/types_test.go**
- Testa `src/types.go`
- Cobertura: 100% dos tipos públicos
- Cenários: serialização, validação, métodos
- Mocks: nenhum necessário

### Testes de Integração (25% - ~1500 linhas de teste)

#### 1. Integração com Setup Component

**tests/integration/api/token_integration_test.go**
- Testa integração real com Setup TokenManager
- Cenários: obtenção de Grid Token, validação, refresh
- Dependencies: Setup Component real

#### 2. Integração de Criação

**tests/integration/create/usb_creation_flow_test.go**
- Testa fluxo completo de criação de USB
- Cenários: detecção USB, download ISO, gravação, cloud-init
- Dependencies: USB real, sistema operacional

**tests/integration/create/cloud_init_injection_test.go**
- Testa injeção de cloud-init em ISO
- Cenários: modificação de ISO, validação de cloud-init
- Dependencies: arquivos ISO reais

#### 3. Integração de Registro

**tests/integration/registration/handshake_flow_test.go**
- Testa fluxo completo de handshake
- Cenários: conexão TCP, validação de token, resposta
- Dependencies: rede real, TCP sockets

**tests/integration/registration/node_lifecycle_test.go**
- Testa lifecycle completo do nó
- Cenários: criação → registro → ativo → inativo
- Dependencies: todos os componentes integrados

#### 4. Integração End-to-End

**tests/integration/plug_and_play_test.go**
- Testa fluxo completo plug-and-play
- Cenários: criação → boot → registro → heartbeat
- Dependencies: hardware real (quando possível)

### Testes End-to-End (5% - ~300 linhas de teste)

#### 1. Cenários Completos

**tests/e2e/scenarios/single_node_creation_test.go**
- Testa criação de nó único
- Cenários: fluxo completo sem falhas
- Dependencies: hardware real

**tests/e2e/scenarios/multiple_nodes_test.go**
- Testa criação de múltiplos nós
- Cenários: concorrência, recursos compartilhados
- Dependencies: hardware real

**tests/e2e/scenarios/node_failure_recovery_test.go**
- Testa recuperação de falhas
- Cenários: falhas de hardware, rede, software
- Dependencies: hardware real

### Testes de Performance

**tests/performance/usb_write_performance_test.go**
- Testa performance de gravação USB
- Métricas: velocidade, throughput, latência
- Benchmarks: diferentes tamanhos de ISO

**tests/performance/handshake_latency_test.go**
- Testa latência de handshake
- Métricas: tempo de resposta, throughput
- Benchmarks: diferentes cargas

**tests/performance/concurrent_nodes_test.go**
- Testa performance com múltiplos nós
- Métricas: escalabilidade, recursos
- Benchmarks: 1, 10, 100, 1000 nós

### Testes de Segurança

**tests/security/token_validation_test.go**
- Testa validação de tokens
- Cenários: tokens válidos, inválidos, expirados
- Security: injeção, manipulação

**tests/security/certificate_validation_test.go**
- Testa validação de certificados
- Cenários: certificados válidos, inválidos, expirados
- Security: man-in-the-middle, replay attacks

**tests/security/ssh_security_test.go**
- Testa segurança SSH
- Cenários: chaves válidas, inválidas, comprometidas
- Security: brute force, key compromise

## Mocks e Fixtures

### Mocks Necessários

#### tests/mocks/setup_token_manager.mock.go
```go
type MockTokenManager struct {
    GetGridTokenFunc    func() (string, error)
    ValidateTokenFunc   func(string) error
    RefreshTokenFunc    func() error
}
```

#### tests/mocks/usb_device.mock.go
```go
type MockUSBDevice struct {
    GetSuitableDevicesFunc func() ([]types.USBDevice, error)
    GetDeviceInfoFunc      func(string) (*types.USBDevice, error)
    ValidateDeviceFunc     func(*types.USBDevice) error
}
```

#### tests/mocks/network.mock.go
```go
type MockNetwork struct {
    ListenFunc func(int) (net.Listener, error)
    DialFunc   func(string) (net.Conn, error)
}
```

#### tests/mocks/filesystem.mock.go
```go
type MockFilesystem struct {
    CreateFileFunc      func(string, []byte) error
    ReadFileFunc        func(string) ([]byte, error)
    FileExistsFunc      func(string) bool
    CreateDirectoryFunc func(string) error
}
```

### Fixtures Necessários

#### tests/fixtures/valid/
- **node-configs/**: Configurações válidas de nós
- **cloud-init/**: Templates válidos de cloud-init
- **usb-devices/**: Informações válidas de dispositivos USB
- **iso-files/**: Arquivos ISO válidos para teste

#### tests/fixtures/invalid/
- **malformed-configs/**: Configurações malformadas
- **invalid-tokens/**: Tokens inválidos
- **corrupted-files/**: Arquivos corrompidos
- **invalid-certificates/**: Certificados inválidos

#### tests/fixtures/edge-cases/
- **boundary-values/**: Valores nos limites
- **empty-data/**: Dados vazios
- **null-pointers/**: Ponteiros nulos
- **extreme-values/**: Valores extremos

## Helpers de Teste

### tests/helpers/assertions.go
```go
func AssertNodeCreated(t *testing.T, result *CreateResult)
func AssertNodeStatus(t *testing.T, status *NodeStatus, expected string)
func AssertHandshakeValid(t *testing.T, request *HandshakeRequest)
func AssertUSBDeviceValid(t *testing.T, device *types.USBDevice)
```

### tests/helpers/builders.go
```go
func BuildValidNodeConfig() *types.NodeConfig
func BuildValidUSBDevice() *types.USBDevice
func BuildValidHandshakeRequest() *HandshakeRequest
func BuildValidCloudInitConfig() *types.CloudInitConfig
```

### tests/helpers/platform.go
```go
func SkipIfNotLinux(t *testing.T)
func SkipIfNotWindows(t *testing.T)
func SkipIfNotMacOS(t *testing.T)
func GetTestUSBDevice() string
```

## Cobertura de Código

### Ferramentas de Cobertura
- **go test -cover**: Cobertura básica
- **go test -coverprofile**: Profile detalhado
- **go tool cover**: Visualização de cobertura
- **gocov**: Cobertura avançada
- **gocov-xml**: Relatórios XML

### Métricas de Cobertura
- **Line Coverage**: 100% obrigatório
- **Branch Coverage**: 100% obrigatório
- **Function Coverage**: 100% obrigatório
- **Path Coverage**: 100% obrigatório
- **Error Coverage**: 100% obrigatório
- **Edge Case Coverage**: 100% obrigatório

### Relatórios de Cobertura
- **HTML Report**: Visualização interativa
- **XML Report**: Para CI/CD
- **JSON Report**: Para análise programática
- **Console Report**: Para desenvolvimento

## Execução de Testes

### Comandos de Teste

#### Testes Unitários
```bash
# Todos os testes unitários
go test ./tests/unit/... -v

# Testes específicos por plataforma
go test ./tests/unit/... -tags linux -v
go test ./tests/unit/... -tags windows -v
go test ./tests/unit/... -tags darwin -v

# Com cobertura
go test ./tests/unit/... -cover -coverprofile=unit.out
```

#### Testes de Integração
```bash
# Todos os testes de integração
go test ./tests/integration/... -v

# Testes específicos
go test ./tests/integration/api/... -v
go test ./tests/integration/create/... -v
go test ./tests/integration/registration/... -v
```

#### Testes End-to-End
```bash
# Todos os testes e2e
go test ./tests/e2e/... -v

# Testes específicos
go test ./tests/e2e/scenarios/... -v
```

#### Testes de Performance
```bash
# Benchmarks
go test ./tests/performance/... -bench=.

# Com profiling
go test ./tests/performance/... -bench=. -cpuprofile=cpu.out
```

#### Testes de Segurança
```bash
# Todos os testes de segurança
go test ./tests/security/... -v
```

### CI/CD Integration

#### GitHub Actions
```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, windows-latest, macos-latest]
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: '1.21'
      - run: go test ./tests/... -cover -coverprofile=coverage.out
      - run: go tool cover -html=coverage.out -o coverage.html
      - uses: codecov/codecov-action@v1
        with:
          file: coverage.out
```

#### Jenkins Pipeline
```groovy
pipeline {
    agent any
    stages {
        stage('Test') {
            parallel {
                stage('Unit Tests') {
                    steps {
                        sh 'go test ./tests/unit/... -cover'
                    }
                }
                stage('Integration Tests') {
                    steps {
                        sh 'go test ./tests/integration/... -v'
                    }
                }
                stage('E2E Tests') {
                    steps {
                        sh 'go test ./tests/e2e/... -v'
                    }
                }
            }
        }
    }
}
```

## Validação de Cobertura

### Checklist de Completude

#### Pré-Implementação
- [x] Análise completa do código fonte
- [x] Catalogação de todos os utilitários em `src/internal/`
- [x] Identificação de todos os componentes
- [x] Mapeamento de implementações OS-específicas
- [x] Lista de tipos e interfaces
- [x] Verificação de não duplicação de utilitários
- [x] Documentação de análise completa

#### Durante Implementação
- [ ] Todos os arquivos de teste em `tests/`
- [ ] Nenhum arquivo modificado em `src/`
- [ ] Uso de utilitários de `src/internal/`
- [ ] Cada arquivo `src/` tem arquivo de teste correspondente
- [ ] Testes OS-específicos com build constraints corretos
- [ ] Mocks funcionando
- [ ] Fixtures organizados
- [ ] Helpers apenas para funcionalidades não existentes

#### Pós-Implementação
- [ ] 100% line coverage de todos os arquivos `src/`
- [ ] 100% branch coverage de todos os arquivos `src/`
- [ ] 100% function coverage de todos os arquivos `src/`
- [ ] Todos os subcomponentes com testes dedicados
- [ ] Todos os códigos OS-específicos com testes de plataforma
- [ ] Integração com Setup TokenManager testada
- [ ] Nenhum utilitário de `src/internal/` duplicado
- [ ] Nenhuma modificação em `src/`
- [ ] Todos os testes em `tests/`
- [ ] Testes independentes e determinísticos
- [ ] Suite executa em < 10 minutos

## Conclusão

Esta estratégia de testes garante **100% de cobertura** do código fonte do Node Component, seguindo rigorosamente as regras do `test-suite.rules.md`. A implementação será feita de forma incremental, começando pelos testes unitários e progredindo para testes de integração e end-to-end.

**Próximos passos:**
1. Implementar mocks necessários
2. Criar fixtures de teste
3. Implementar testes unitários (70%)
4. Implementar testes de integração (25%)
5. Implementar testes end-to-end (5%)
6. Validar cobertura 100%
7. Integrar com CI/CD
