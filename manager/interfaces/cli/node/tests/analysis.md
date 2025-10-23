# Node Component - Análise Completa do Código Fonte

**Data**: 2025-01-28
**Versão**: 2.0 (Revisão Profunda)
**Status**: Pronto para Implementação

---

## 1. MAPEAMENTO DE ARQUIVOS FONTE

### Arquivos Principais (src/)

#### 1.1 Componente Principal
- **node.go** (500+ linhas)
  - `NodeManager`: Orquestrador principal
  - Métodos públicos: CreateNode(), ListNodes(), GetNodeStatus(), DeleteNode(), etc.
  - Dependências: EventBus, NodeStateManager, Logger, Configuration
  - Padrão: Facade/Orchestrator

#### 1.2 Subcomponente Create
- **create.go** (576 linhas)
  - `CreateSubcomponent`: Orquestrador de criação de nós
  - `CreateNodeManager`: Manager para criar nós
  - Métodos: CreateNode(), CreateNodeInteractive()
  - Tipos: CreateOptions, CreateResult
  - **CRÍTICO**: validateSetupComponent(), validatePrerequisites()

#### 1.3 Gerenciamento de Estado
- **node_state.go** (existente ou a criar)
  - `NodeStateManager`: Gerenciador de estado thread-safe
  - Responsabilidade: Cache em memória + Persistência em arquivo
  - Fonte de verdade: `~/.syntropy/nodes/`

#### 1.4 Geração de Configuração
- **auto_config_generator.go** (285+ linhas)
  - `AutoConfigGenerator`: Geração automática de config
  - Métodos: GenerateNodeConfig(), GenerateNodeID(), GenerateSSHKeys()
  - Integração com Setup Component via TokenIntegration

#### 1.5 Geração Cloud-Init
- **cloud_init_generator.go** (464+ linhas)
  - `CloudInitGenerator`: Geração de cloud-init
  - Métodos: GenerateCloudInit(), generateUserData(), etc.
  - Templates: user-data.yaml, network-config.yaml, meta-data.yaml
  - **CRÍTICO**: Criptografia de token antes de incluir

#### 1.6 Integração com Setup Component
- **token_integration.go** (279+ linhas)
  - `TokenIntegration`: Integração com Setup Component TokenManager
  - Métodos: GetGridToken(), ValidateToken(), RefreshToken()
  - Cache com expiração de token

#### 1.7 Detecção USB
- **usb_detector.go** (400+ linhas)
  - Interface: `USBDetector`
  - Implementações: USBDetectorLinux, USBDetectorWindows, USBDetectorMacOS
  - Factory pattern: `NewUSBDetectorFactory()`
  - Validação de segurança (não gravar em disco sistema)

#### 1.8 Download de ISO
- **iso_downloader.go** (350+ linhas)
  - `ISODownloader`: Download e cache de ISOs Ubuntu
  - Verificação SHA256
  - Cache em `~/.syntropy/cache/isos/`

#### 1.9 Gravação USB
- **usb_writer.go** + implementações plataforma-específicas
  - Interface: `USBWriter`
  - Injeção de cloud-init no ISO
  - Multi-plataforma (Windows/Linux/macOS)

#### 1.10 Protocolos de Comunicação
- **handshake.go**: Protocolo de handshake seguro
- **listener.go**: Listener TCP na porta 51000
- **heartbeat.go**: Manutenção de conexão contínua

#### 1.11 Implementações Base
- **implementations.go**: Implementações de EventBus, Logger, NodeStateManager
- **setup_adapter.go**: Adapter para Setup Component
- **types.go**: Tipos públicos exportados

#### 1.12 Interface CLI
- **cli.go**: Comandos CLI (create, list, status, logs, etc.)

---

## 2. ARQUITETURA DE INTERNO (src/internal/)

### 2.1 Tipos (src/internal/types/)

#### Tipos de Dados
```
NodeStatus:
  - NodeID: string (imutável)
  - Status: string (pending → active → inactive)
  - CreatedAt: time.Time (imutável)
  - RegisteredAt: time.Time (preenchido no handshake)
  - IPAddress: string
  - LastHeartbeat: time.Time
  - Hardware: *HardwareInfo (nil até conexão)

HardwareInfo:
  - CPUCores: int
  - MemoryGB: int
  - DiskGB: int
  - Hostname: string
  - IPAddress: string
  - OSVersion: string
  - KernelVersion: string

NodeConfig:
  - NodeID: string
  - GridToken: string (do Setup Component)
  - SSHPublicKey: string
  - SSHPrivateKey: string
  - NodeCertificate: string
  - CommandStationIP: string
  - CreatedAt: time.Time
  - ExpiresAt: time.Time

CloudInitConfig:
  - NodeID: string
  - UserData: string
  - NetworkConfig: string
  - MetaData: string
  - Valid: bool
```

#### Interfaces Críticas
```
EventBus: Subscribe, Unsubscribe, Publish, Close
Logger: Debug, Info, Warn, Error, Fatal, SetLevel, WithFields
USBDetector: DetectDevices, DetectRemovableDevices, ValidateDevice, etc.
ISODownloader: DownloadISO, GetCachedISO, VerifyISO
USBWriter: WriteToUSB, ValidateWriteOperation, VerifyWriteOperation
CloudInitGenerator: GenerateCloudInit, ValidateCloudInit, LoadTemplate
ConfigGenerator: GenerateNodeConfig, GenerateNodeID, GenerateSSHKeys
TokenIntegration: GetGridToken, ValidateToken, RefreshToken, IsTokenValid
NodeStateManager: CreateNode, GetNode, UpdateNodeStatus, TransitionToActive, etc.
```

### 2.2 Helpers (src/internal/helpers/)

#### Funções Disponíveis
```
GenerateNodeID(existingNodes []string): string
  - Gera NodeID sequencial (node-01, node-02, etc.)
  
ValidateNodeID(nodeID string): error
  - Valida formato do NodeID
  
GenerateSSHKeys(): (*SSHKeys, error)
  - Gera par de chaves SSH RSA 2048 bits
  
[... outras funções de criptografia, validação, etc ...]
```

### 2.3 Constantes (src/internal/constants/)

#### Valores Críticos
```
DefaultMinUSBCapacity = 8GB
DefaultMaxUSBCapacity = 64GB
DefaultRegistrationPort = 51000
DefaultRegistrationTimeout = 30 minutos
DefaultHeartbeatInterval = 30 segundos
DefaultHeartbeatTimeout = 10 segundos
DefaultMaxHeartbeatFailures = 3
DefaultTokenExpiry = 24 horas
DefaultSSHKeySize = 2048

Caminhos:
  DefaultConfigDir = "~/.syntropy"
  DefaultNodeConfigDir = "~/.syntropy/nodes"
  DefaultCacheDir = "~/.syntropy/cache"
  DefaultLogDir = "~/.syntropy/logs"
```

---

## 3. ARQUITETURA DE PERSISTÊNCIA

### 3.1 Estrutura de Diretórios
```
$HOME/.syntropy/
├── config/                 # Criado por Setup Component
│   ├── setup_state.json   # Estado do Setup
│   └── ...
├── cache/                  # Criado por Setup Component
│   └── isos/              # ISOs em cache
├── logs/                   # Criado por Setup Component
└── nodes/                  # Criado por Setup Component OU Node Component
    ├── node-01/
    │   ├── status.json    # Evolui conforme nó conecta
    │   └── config.json    # Imutável (NodeID, token, etc)
    ├── node-02/
    └── ...
```

### 3.2 Evolução do status.json

**Fase 1 - Criação do nó (Step 7 do CreateNode):**
```json
{
  "node_id": "node-01",
  "status": "pending",
  "created_at": "2025-01-28T10:00:00Z",
  "command_station_ip": "192.168.1.100"
}
```

**Fase 2 - Após conexão/handshake (via heartbeat):**
```json
{
  "node_id": "node-01",
  "status": "active",
  "created_at": "2025-01-28T10:00:00Z",
  "registered_at": "2025-01-28T10:05:30Z",
  "command_station_ip": "192.168.1.100",
  "ip_address": "192.168.1.101",
  "last_heartbeat": "2025-01-28T10:06:00Z",
  "hardware": {
    "cpu_cores": 8,
    "memory_gb": 28,
    "disk_gb": 512,
    "hostname": "node-01",
    "ip_address": "192.168.1.101",
    "os_version": "Ubuntu 24.04 LTS",
    "kernel_version": "6.6.x"
  }
}
```

---

## 4. PROBLEMA DE $HOME (CRÍTICO)

### 4.1 Problema Identificado
- `os.Getenv("HOME")` pode não funcionar em:
  - Docker containers
  - CI/CD environments
  - Usuários especiais do sistema
  - Windows (usa USERPROFILE)

### 4.2 Solução
```go
// Strategy: Use os.UserHomeDir() PRIMEIRO (Go 1.12+)
// Fallback: Tentar $HOME/$USERPROFILE conforme SO
// Last resort: /tmp/.syntropy (não ideal)
```

---

## 5. INTEGRAÇÃO COM SETUP COMPONENT

### 5.1 O que Setup Component Cria
1. Estrutura de diretórios ($HOME/.syntropy/*)
2. Chaves Ed25519 (keyManager)
3. Grid Token (tokenManager) - CRÍTICO
4. Arquivo de estado (stateManager)

### 5.2 Validação Necessária em Node Component
1. TokenIntegration != nil
2. Grid Token acessível e válido (não placeholder)
3. Estrutura de diretórios existente
4. Arquivo de estado do Setup válido
5. Permissões corretas (0700 dirs, 0600 arquivos)

---

## 6. FLUXO DE CRIAÇÃO DE NÓ (CreateNode)

### Step 1: Validar Setup Component
- Verificar TokenIntegration
- Verificar Grid Token
- Verificar estrutura de diretórios
- Verificar permissões

### Step 2: Validar Pré-Requisitos
- Validar Setup Component (já feito em Step 1)
- Validar ferramentas requeridas
- Validar permissões de escrita

### Step 3: Gerar Configuração
- NodeID sequencial
- SSH keys
- Node certificate
- IP da Command Station
- Grid Token (do Setup Component)

### Step 4: Detectar USB
- Se não especificado, detectar automaticamente
- Validar: removível, tamanho mínimo, não é disco sistema

### Step 5: Download ISO
- Usar cache se disponível
- Verificar SHA256
- Download com progress bar

### Step 6: Gerar Cloud-Init
- user-data.yaml (instalação e config)
- network-config.yaml (DHCP)
- meta-data.yaml (hostname, instance-id)
- Criptografar token antes de incluir

### Step 7: Gravar USB
- Injetar cloud-init no ISO
- Gravar em dispositivo USB
- Validar gravação

### Step 8: Salvar Estado
- Usar NodeStateManager APENAS (não duplicar persistência)
- Persistir em `~/.syntropy/nodes/node-XX/`
- Iniciar sincronização automática

### Step 9: Iniciar Listener
- TCP na porta 51000
- Aguardar conexão do nó específico
- Timeout: 30 minutos

---

## 7. THREAD-SAFETY

### Sincronização Crítica
- `NodeStateManager`: sync.RWMutex para proteger nodeStates
- Operações de arquivo: Usar file locks (não implementado ainda)
- Operações de token: Cache com sincronização

---

## 8. SEGURANÇA

### Proteção de Dados Sensíveis
- ❌ NÃO logar tokens/chaves
- ✅ Permissões 0700 (dirs) e 0600 (arquivos)
- ✅ NÃO salvar chaves privadas em arquivo
- ✅ Criptografar token antes de injetar no cloud-init
- ✅ Validação de integridade em checkpoints

---

## 9. COMPONENTES A IMPLEMENTAR/MODIFICAR

### Modific ações Necessárias
1. **src/create.go**
   - validateSetupComponent() - ROBUSTO
   - getSyntropyDir() - CROSS-PLATFORM
   - Refatorar saveNodeState() para usar APENAS NodeStateManager

2. **src/node_state.go** (criar ou refatorar)
   - NodeStateManager com sincronização arquivo↔memória
   - LoadState(), CreateNode(), UpdateNodeStatus()
   - StartAutoSync(), detectExternalChanges()
   - Persistência em ~/.syntropy/nodes/

3. **src/types.go**
   - HardwareInfo com CollectedAt e Source
   - NodeStatus com campos progressivos

4. **Testes** (criar em tests/)
   - create_integration_setup_test.go
   - node_state_sync_test.go
   - node_hardware_test.go
   - create_recovery_test.go
   - create_security_test.go

---

## 10. RESUMO EXECUTIVO

**Fonte de Verdade**: Arquivos em `~/.syntropy/nodes/` (persistência)
**Cache**: NodeStateManager em memória
**Sincronização**: Automática a cada 10 segundos
**Resolução de $HOME**: os.UserHomeDir() → env vars → fallback
**Validação Setup**: Multi-componente (token, dirs, state, permissions)
**Hardware**: Progressivo (nil → coletado após conexão)
**Segurança**: 0700/0600 permissions, sem logs sensíveis, validação integridade
