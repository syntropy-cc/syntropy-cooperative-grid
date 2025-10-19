# Syntropy Cooperative Grid - MVP Technical Specification v2.0

**Versão**: 2.0.0  
**Data**: Outubro 2025  
**Objetivo**: Especificação técnica completa para construção de mini-cloud pessoal com 6 nós físicos  
**Audiência**: LLMs e desenvolvimento iterativo  
**Status de Implementação**:
- ✅ **Setup Component** (80% - em `manager/interfaces/cli/setup/`)
- ✅ **Infrastructure Base** (100% - em `infrastructure/`)
- 🚧 **Node Creation Component** (0% - a implementar)
- 🚧 **Workload Component** (0% - a implementar)
- 🚧 **Management Component** (0% - a implementar)

---

## GLOSSÁRIO TÉCNICO (USE ESTA TERMINOLOGIA)

| Termo | Definição | Exemplo de Uso |
|-------|-----------|----------------|
| **syntropy** | Binário CLI (sempre lowercase) | `syntropy node create` |
| **Command Station** | PC de trabalho que gerencia toda a Grid | Seu laptop/desktop principal |
| **Node** | Servidor físico que executa workloads | Hardware com Ubuntu instalado |
| **Agent** | Daemon em Go rodando em cada Node | `/opt/syntropy/bin/syntropy-agent` |
| **Workload** | Container Docker + configuração YAML | `nginx:latest` com 512MB RAM |
| **Grid** | Conjunto de todos os Nodes | Sua mini-cloud pessoal |
| **Grid Token** | Token secreto compartilhado para registro | UUID gerado no setup |
| **USB Provisioner** | USB bootável com cloud-init | Pendrive com Ubuntu + templates |
| **Registration Protocol** | Handshake entre Node e Command Station | Node → announce → Command Station → ack |
| **Hardware Manifest** | Inventário automático de recursos do Node | CPU: 8 cores, RAM: 28GB, Disk: 500GB |

---

## 1. VISÃO GERAL DO MVP

### 1.1 Objetivo
Criar uma **mini-cloud pessoal** de 6 Nodes físicos gerenciados por CLI a partir de uma Command Station, com provisionamento via USB, registro automático e gerenciamento de workloads containerizados.

### 1.2 Arquitetura Macro

```
┌───────────────────────────────────────────────────────────────┐
│              COMMAND STATION (PC de Trabalho)                 │
│  ─────────────────────────────────────────────────────────    │
│  • syntropy CLI (Go)                                          │
│  • ~/.syntropy/ (estado local, inventário, logs)             │
│  • Projetos: infrastructure/ (templates, chaves)             │
│  • SSH para comunicação com Nodes                            │
│  • Grid Token (segredo compartilhado)                        │
└───────────────────────────────────────────────────────────────┘
                          │
                          │ SSH (10.0.100.0/24)
                          │ + Registration Protocol
                          │
        ┌─────────────────┴────────────────────┐
        │                                       │
   ┌────▼─────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐
   │  Node 1  │  │  Node 2  │  │  Node 3  │  │  Node 4-6│
   │          │  │          │  │          │  │          │
   │ DHCP IP  │  │ DHCP IP  │  │ DHCP IP  │  │ DHCP IP  │
   │ Auto     │  │ Auto     │  │ Auto     │  │ Auto     │
   │ Detect   │  │ Detect   │  │ Detect   │  │ Detect   │
   │          │  │          │  │          │  │          │
   │ Ubuntu   │  │ Ubuntu   │  │ Ubuntu   │  │ Ubuntu   │
   │ Docker   │  │ Docker   │  │ Docker   │  │ Docker   │
   │ Agent    │  │ Agent    │  │ Agent    │  │ Agent    │
   │          │  │          │  │          │  │          │
   │ Hardware │  │ Hardware │  │ Hardware │  │ Hardware │
   │ Manifest │  │ Manifest │  │ Manifest │  │ Manifest │
   │ Auto     │  │ Auto     │  │ Auto     │  │ Auto     │
   └──────────┘  └──────────┘  └──────────┘  └──────────┘
```

**Princípios Fundamentais**:
1. **Templates Centralizados**: Todos os templates em `infrastructure/cloud-init/`
2. **Auto-Detecção**: Nodes detectam hardware automaticamente
3. **Registro Automático**: Nodes se registram via Grid Token
4. **Configuração Mínima**: Usuário define apenas nome do Node
5. **Sincronização Bidirecional**: Command Station ↔ Nodes sempre sincronizados

---

## 2. COMPONENTES DO SISTEMA

### 2.1 Estrutura de Diretórios do Projeto

```
syntropy-cooperative-grid/
├── infrastructure/                    # ✅ EXISTENTE
│   ├── cloud-init/                   # Templates cloud-init
│   │   ├── user-data-template.yaml   # Template principal
│   │   ├── meta-data-template.yaml   # Metadados
│   │   └── network-config-template.yaml
│   ├── key_manager.go                # Gerenciador de chaves SSH
│   ├── template_manager.go           # Renderizador de templates
│   └── README.md                     # Documentação
│
├── manager/interfaces/cli/           # CLI Manager
│   ├── main.go                       # Entry point
│   ├── setup/                        # ✅ Setup Component (80%)
│   │   ├── src/setup.go
│   │   ├── src/configurator.go
│   │   ├── src/key_manager.go
│   │   └── docs/LEARN.md
│   │
│   ├── node/                         # 🚧 Node Component (a criar)
│   │   ├── README.md                 # Documentação do componente
│   │   ├── ARCHITECTURE.md           # Arquitetura detalhada
│   │   ├── node.go                   # Orquestrador principal
│   │   ├── create/                   # Subcomponente: Creation
│   │   │   ├── create.go             # Interface
│   │   │   ├── create_windows.go     # Implementação Windows
│   │   │   ├── create_linux.go       # Implementação Linux
│   │   │   ├── usb_detector.go       # Detecção de USBs
│   │   │   ├── usb_detector_windows.go
│   │   │   ├── usb_detector_linux.go
│   │   │   ├── iso_downloader.go     # Download Ubuntu ISO
│   │   │   ├── cloud_init_generator.go # Geração cloud-init
│   │   │   └── usb_writer.go         # Gravação de USB
│   │   ├── registration/             # Subcomponente: Registration
│   │   │   ├── registration.go       # Protocol de registro
│   │   │   ├── token_manager.go      # Gerenciamento de tokens
│   │   │   └── handshake.go          # Handshake automático
│   │   └── inventory/                # Subcomponente: Inventory
│   │       ├── inventory.go          # Gerenciamento de inventário
│   │       ├── sync.go               # Sincronização Command Station ↔ Node
│   │       └── hardware_manifest.go  # Recebe manifesto de hardware
│   │
│   ├── workload/                     # 🚧 Workload Component
│   │   ├── README.md
│   │   ├── ARCHITECTURE.md
│   │   ├── workload.go               # Orquestrador principal
│   │   ├── deploy/                   # Subcomponente: Deployment
│   │   │   ├── deploy.go
│   │   │   ├── scheduler.go          # Scheduler simples
│   │   │   └── docker_executor.go    # Execução via Docker
│   │   ├── lifecycle/                # Subcomponente: Lifecycle
│   │   │   ├── lifecycle.go
│   │   │   ├── start.go
│   │   │   ├── stop.go
│   │   │   └── restart.go
│   │   └── monitoring/               # Subcomponente: Monitoring
│   │       ├── monitoring.go
│   │       ├── logs.go
│   │       └── metrics.go
│   │
│   └── management/                   # 🚧 Management Component
│       ├── README.md
│       ├── ARCHITECTURE.md
│       ├── management.go             # Orquestrador principal
│       ├── discovery/                # Subcomponente: Discovery
│       │   ├── discovery.go
│       │   └── network_scanner.go
│       ├── health/                   # Subcomponente: Health
│       │   ├── health.go
│       │   ├── healthcheck.go
│       │   └── diagnostics.go
│       └── sync/                     # Subcomponente: Sync
│           ├── sync.go
│           ├── state_sync.go
│           └── manifest_sync.go
│
└── ~/.syntropy/                      # Estado local (runtime)
    ├── config/
    │   ├── manager.yaml              # Config Command Station
    │   └── grid-token.txt            # Grid Token (secreto)
    ├── nodes/                        # Inventário de Nodes
    │   ├── node-01.yaml              # Metadados + Hardware Manifest
    │   └── node-02.yaml
    ├── keys/                         # Chaves SSH
    │   ├── command-station.key
    │   ├── command-station.pub
    │   └── nodes/                    # Chaves dos Nodes
    │       ├── node-01-owner.pub
    │       └── node-01-community.pub
    ├── workloads/                    # Workloads deployados
    │   └── nginx-001.yaml
    ├── cache/                        # Cache
    │   └── ubuntu-24.04.iso          # ISO Ubuntu em cache
    └── logs/                         # Logs
        └── syntropy.log
```

---

## 3. PILAR 1: SETUP DA COMMAND STATION ✅

### 3.1 Status Atual
**Implementado**: 80% em `manager/interfaces/cli/setup/`

**O que já funciona**:
```bash
$ syntropy setup run --force
✅ Criar ~/.syntropy/
✅ Criar config/manager.yaml
✅ Gerar estrutura de diretórios
```

### 3.2 O que precisa ser adicionado

#### 3.2.1 Geração SEGURA de Grid Token

**⚠️ SEGURANÇA CRÍTICA**: Grid Token NÃO deve ser armazenado em texto plano!

**Arquivo**: `manager/interfaces/cli/setup/src/token_manager.go`

```go
package setup

import (
	"crypto/rand"
	"fmt"
	
	"github.com/zalando/go-keyring"
)

const (
	KeyringService = "syntropy-grid"
	KeyringUser    = "grid-token"
)

// TokenManager gerencia Grid Token de forma SEGURA
type TokenManager struct{}

// GenerateToken gera novo Grid Token (UUID v4)
func (tm *TokenManager) GenerateToken() (string, error) {
	tokenBytes := make([]byte, 16)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	
	token := fmt.Sprintf("%x-%x-%x-%x-%x",
		tokenBytes[0:4],
		tokenBytes[4:6],
		tokenBytes[6:8],
		tokenBytes[8:10],
		tokenBytes[10:16],
	)
	
	return token, nil
}

// SaveToken salva token SEGURAMENTE no keyring do sistema
// Windows: Credential Manager (criptografado)
// Linux: Secret Service API / gnome-keyring
// macOS: Keychain (criptografado por hardware)
func (tm *TokenManager) SaveToken(token string) error {
	if err := keyring.Set(KeyringService, KeyringUser, token); err != nil {
		return fmt.Errorf("failed to save token to system keyring: %w", err)
	}
	
	fmt.Println("✅ Grid Token saved securely in system keyring")
	fmt.Printf("   Service: %s\n", KeyringService)
	fmt.Printf("   Token: %s...[HIDDEN]\n", token[:8])
	
	return nil
}

// LoadToken carrega token do keyring
func (tm *TokenManager) LoadToken() (string, error) {
	token, err := keyring.Get(KeyringService, KeyringUser)
	if err != nil {
		return "", fmt.Errorf("failed to load token from keyring: %w", err)
	}
	
	return token, nil
}
```

**Dependências**:
```bash
# Adicionar ao go.mod
go get github.com/zalando/go-keyring@latest

# Linux: instalar dependências do sistema
sudo apt-get install libsecret-1-dev  # Ubuntu/Debian
sudo dnf install libsecret-devel      # Fedora/RHEL

# Windows/macOS: sem dependências extras
```

**Integração no setup**:
```go
// No final de setup.go, adicionar:
func (s *Setup) Run() error {
    // ... existing code ...
    
    // Gerar Grid Token
    log.Println("🔐 Generating Grid Token...")
    token, err := GenerateGridToken()
    if err != nil {
        return err
    }
    
    if err := SaveGridToken(token, s.configDir); err != nil {
        return err
    }
    
    log.Printf("✅ Grid Token: %s", token[:8]+"..."+"[HIDDEN]")
    
    return nil
}
```

#### 3.2.2 Melhorar Geração de Chaves SSH

Usar `infrastructure/key_manager.go` existente:

```go
import (
    infra "github.com/syntropy-cooperative-grid/infrastructure"
)

func (s *Setup) GenerateSSHKeys() error {
    keyManager := infra.NewKeyManager(filepath.Join(s.baseDir, "keys"))
    
    // Gerar chave da Command Station (ED25519)
    log.Println("🔑 Generating Command Station SSH key...")
    keyPair, err := keyManager.GenerateKeyPair(infra.OwnerKey, "command-station")
    if err != nil {
        return fmt.Errorf("failed to generate SSH key: %w", err)
    }
    
    // Salvar chave
    if err := keyManager.SaveKeyPair(keyPair, infra.OwnerKey, "command-station"); err != nil {
        return fmt.Errorf("failed to save SSH key: %w", err)
    }
    
    log.Printf("✅ SSH Key generated: %s", keyPair.Fingerprint)
    
    return nil
}
```

### 3.3 Configuração Atualizada (manager.yaml)

```yaml
# ~/.syntropy/config/manager.yaml
version: "2.0"

command_station:
  name: "my-grid-hq"
  id: "command-station-001"
  created_at: "2025-10-10T10:00:00Z"
  
# Grid Token (armazenado SEGURAMENTE no keyring do sistema)
# Windows: Credential Manager
# Linux: Secret Service / gnome-keyring
# macOS: Keychain
grid:
  token_storage: "keyring"  # NUNCA "file" em produção!
  keyring_service: "syntropy-grid"
  keyring_user: "grid-token"
  registration_timeout: 300s  # 5 minutos para Node se registrar
  
network:
  discovery_mode: "dhcp"  # Nodes usam DHCP
  subnet: "10.0.100.0/24"  # Subnet esperada
  ssh_port: 22
  
ssh:
  user: "syntropy"
  key: "~/.syntropy/keys/command-station-owner.key"
  timeout: 30s
  
nodes:
  auto_accept: true  # Aceitar Nodes automaticamente se tiverem Grid Token
  default_os: "ubuntu-24.04"
  
# Paths
paths:
  infrastructure: "../../../infrastructure"  # Relativo a ~/.syntropy/config
  templates: "../../../infrastructure/cloud-init"
  cache: "../cache"
  logs: "../logs"
```

---

## 4. PILAR 2: CRIAÇÃO DE NODES VIA USB 🚧

### ⚠️ IMPORTANTE: Templates Cloud-Init

**STATUS ATUAL DOS TEMPLATES**: Os templates em `infrastructure/cloud-init/` foram criados para uma arquitetura mais avançada e **NÃO estão alinhados com este MVP**.

**Inconsistências Identificadas**:
- ❌ Falta variável `${GRID_TOKEN}` (crítico para registro)
- ❌ Falta script de auto-detecção de hardware
- ❌ Falta script de registration protocol
- ❌ Agent tenta baixar do GitHub (não existe ainda)
- ❌ Usa variáveis não definidas no MVP (`${NODE_CERT_PATH}`, `${DISCOVERY_SERVER}`, etc.)
- ❌ Network config muito complexo (bridges, VLANs desnecessários para MVP)

**Solução**: Ver `docs/architecture/MVP-CORRECTIONS.md` para templates corrigidos:
- `user-data-mvp.yaml` - Template simplificado e alinhado com MVP
- `network-config-mvp.yaml` - Configuração DHCP simples
- `meta-data-mvp.yaml` - Metadados mínimos

**Durante Implementação**:
1. Usar templates `-mvp.yaml` (a serem criados)
2. Manter templates atuais como `-advanced.yaml` (pós-MVP)
3. Seguir especificação deste documento MVP

### 4.1 Arquitetura do Componente

```
node/
├── README.md                    # Documentação: O que é, como usar
├── ARCHITECTURE.md              # Arquitetura: Design decisions, fluxos
├── node.go                      # Orquestrador principal (< 500 linhas)
│
├── create/                      # Subcomponente: Criação de Nodes
│   ├── create.go                # Interface e lógica comum (< 300 linhas)
│   ├── create_windows.go        # Implementação Windows (< 500 linhas)
│   ├── create_linux.go          # Implementação Linux (< 500 linhas)
│   ├── usb_detector.go          # Interface de detecção (< 200 linhas)
│   ├── usb_detector_windows.go  # Detecção USB Windows (< 300 linhas)
│   ├── usb_detector_linux.go    # Detecção USB Linux (< 300 linhas)
│   ├── iso_downloader.go        # Download Ubuntu ISO (< 400 linhas)
│   ├── cloud_init_generator.go  # Geração cloud-init (< 400 linhas)
│   └── usb_writer.go            # Gravação USB (< 400 linhas)
│
├── registration/                # Subcomponente: Registro de Nodes
│   ├── registration.go          # Protocol de registro (< 400 linhas)
│   ├── token_manager.go         # Validação de Grid Token (< 300 linhas)
│   └── handshake.go             # Handshake Node ↔ Command Station (< 400 linhas)
│
└── inventory/                   # Subcomponente: Inventário
    ├── inventory.go             # Gerenciamento de inventário (< 400 linhas)
    ├── sync.go                  # Sincronização bidirecional (< 400 linhas)
    └── hardware_manifest.go     # Recebe/processa manifesto (< 300 linhas)
```

### 4.2 Fluxo Completo de Criação

```
╔════════════════════════════════════════════════════════════════╗
║  FASE 1: CRIAÇÃO DO USB NA COMMAND STATION                     ║
╚════════════════════════════════════════════════════════════════╝

$ syntropy node create

┌─────────────────────────────────────────────────────────────┐
│ STEP 1: Input do Usuário (mínimo)                           │
├─────────────────────────────────────────────────────────────┤
│ Prompt: "Enter node name (e.g., node-01): "                 │
│ Input: node-01                                               │
│                                                               │
│ ✅ APENAS ISSO! Todo resto é automático                     │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ STEP 2: Detecção Automática de USBs                         │
├─────────────────────────────────────────────────────────────┤
│ Scanning USB devices...                                      │
│                                                               │
│ Found USBs:                                                   │
│ [1] E:\ (SanDisk 32GB) - ⚠️  Contains data                  │
│ [2] F:\ (Kingston 64GB) - ✅ Empty                          │
│ [3] G:\ (Samsung 128GB) - ⚠️  System partition (BLOCKED)   │
│                                                               │
│ ⚠️  USB will be FORMATTED! All data will be LOST!           │
│ Select USB [1-2]: 2                                          │
│                                                               │
│ Confirm format F:\ Kingston 64GB? (type YES): YES            │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ STEP 3: Geração Automática de Cloud-Init                    │
├─────────────────────────────────────────────────────────────┤
│ [1/8] Loading Grid Token...                                  │
│   ✅ Token: a4f3b2e7-...                                     │
│                                                               │
│ [2/8] Generating Node SSH keys...                            │
│   ✅ Owner key: node-01-owner.key                           │
│   ✅ Community key: node-01-community.key                   │
│                                                               │
│ [3/8] Loading cloud-init templates...                        │
│   ✅ infrastructure/cloud-init/user-data-template.yaml      │
│   ✅ infrastructure/cloud-init/meta-data-template.yaml      │
│                                                               │
│ [4/8] Rendering templates with variables...                  │
│   Variables:                                                  │
│   - NODE_NAME: node-01                                       │
│   - GRID_TOKEN: a4f3b2e7-... (from ~/.syntropy/)            │
│   - SSH_PUBLIC_KEY: ssh-ed25519 AAAA... (generated)         │
│   - COMMAND_STATION_IP: 10.0.100.1 (from config)            │
│   ✅ Templates rendered                                      │
│                                                               │
│ [5/8] Downloading Ubuntu Server ISO (if needed)...           │
│   ✅ Cached: ~/.syntropy/cache/ubuntu-24.04.iso             │
│                                                               │
│ [6/8] Injecting cloud-init into ISO...                       │
│   ⏱️  Creating custom ISO... (2-3 min)                       │
│   ✅ Custom ISO ready                                        │
│                                                               │
│ [7/8] Writing bootable USB...                                │
│   ⏱️  Writing to F:\ ... (5-7 min)                           │
│   ✅ USB bootable created                                    │
│                                                               │
│ [8/8] Registering Node in inventory...                       │
│   ✅ Saved: ~/.syntropy/nodes/node-01.yaml                  │
│   Status: waiting_provisioning                               │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ ✅ USB BOOTÁVEL CRIADO COM SUCESSO!                         │
├─────────────────────────────────────────────────────────────┤
│ Node: node-01                                                 │
│ USB: F:\ (Kingston 64GB)                                      │
│                                                               │
│ 🔄 PRÓXIMOS PASSOS:                                          │
│ 1. Remova o USB da Command Station                           │
│ 2. Insira no hardware virgem                                 │
│ 3. Ligue o hardware (BIOS deve bootar do USB)               │
│ 4. Aguarde ~15 minutos (instalação automática)              │
│ 5. Node se registrará automaticamente!                       │
│                                                               │
│ 📊 Acompanhe o registro:                                     │
│    $ syntropy node list --watch                              │
│    $ syntropy node status node-01                            │
└─────────────────────────────────────────────────────────────┘


╔════════════════════════════════════════════════════════════════╗
║  FASE 2: BOOT AUTOMÁTICO DO NODE (SEM INTERVENÇÃO)            ║
╚════════════════════════════════════════════════════════════════╝

[Hardware virgem com USB inserido]

⏱️  00:00 - BIOS detecta USB bootável
⏱️  00:30 - GRUB inicia Ubuntu installer
⏱️  01:00 - Cloud-init lê user-data.yaml
⏱️  02:00 - Instalando Ubuntu Server 24.04...
⏱️  05:00 - Instalando pacotes (Docker, tools)...
⏱️  08:00 - Configurando firewall (UFW)...
⏱️  09:00 - Instalando Syntropy Agent...
⏱️  10:00 - Detectando hardware automaticamente...
         ├─ CPU: 8 cores (AMD Ryzen 7)
         ├─ RAM: 28GB
         ├─ Disk: 512GB NVMe
         └─ Network: Intel I225-V (DHCP: 10.0.100.45)
⏱️  11:00 - Criando Hardware Manifest...
⏱️  12:00 - Reiniciando...
⏱️  13:00 - Boot completo
⏱️  14:00 - Syntropy Agent iniciando...
⏱️  14:30 - Iniciando Registration Protocol...


╔════════════════════════════════════════════════════════════════╗
║  FASE 3: REGISTRATION PROTOCOL (AUTOMÁTICO)                    ║
╚════════════════════════════════════════════════════════════════╝

Node (10.0.100.45) → Command Station (10.0.100.1)

┌─────────────────────────────────────────────────────────────┐
│ STEP 1: Node Announcement (via mDNS/broadcast)              │
├─────────────────────────────────────────────────────────────┤
│ Node: "Hello! I am node-01, Grid Token: a4f3b2e7-..."       │
│                                                               │
│ Broadcast: {                                                  │
│   "type": "node_announcement",                               │
│   "node_name": "node-01",                                    │
│   "grid_token": "a4f3b2e7-1234-5678-90ab-cdef12345678",    │
│   "ip": "10.0.100.45",                                       │
│   "ssh_port": 22,                                            │
│   "public_key": "ssh-ed25519 AAAA...",                      │
│   "hardware_manifest": {                                     │
│     "cpu_cores": 8,                                          │
│     "cpu_model": "AMD Ryzen 7 5800X",                       │
│     "ram_gb": 28,                                            │
│     "disk_gb": 512,                                          │
│     "disk_type": "nvme",                                     │
│     "network_interfaces": [                                  │
│       {"name": "enp0s1", "mac": "00:1a:2b:3c:4d:5e"}       │
│     ]                                                         │
│   }                                                           │
│ }                                                             │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ STEP 2: Command Station Validation                          │
├─────────────────────────────────────────────────────────────┤
│ Command Station: Received announcement from 10.0.100.45      │
│                                                               │
│ [✅ CHECK 1] Grid Token validation                          │
│   Expected: a4f3b2e7-1234-5678-90ab-cdef12345678            │
│   Received: a4f3b2e7-1234-5678-90ab-cdef12345678            │
│   ✅ MATCH                                                   │
│                                                               │
│ [✅ CHECK 2] Node name validation                           │
│   Inventory: node-01 (status: waiting_provisioning)          │
│   Announced: node-01                                          │
│   ✅ MATCH                                                   │
│                                                               │
│ [✅ CHECK 3] SSH key validation                             │
│   Expected pub key: ~/.syntropy/keys/nodes/node-01-owner.pub│
│   Announced: ssh-ed25519 AAAA...                            │
│   ✅ MATCH                                                   │
│                                                               │
│ Decision: ✅ ACCEPT NODE                                     │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ STEP 3: Command Station → Node Acknowledgment               │
├─────────────────────────────────────────────────────────────┤
│ Command Station → Node (SSH 10.0.100.45:22)                 │
│                                                               │
│ Response: {                                                   │
│   "type": "registration_ack",                                │
│   "status": "accepted",                                       │
│   "node_id": "node-01",                                      │
│   "command_station_ip": "10.0.100.1",                       │
│   "command_station_public_key": "ssh-ed25519 BBBB...",      │
│   "grid_config": {                                           │
│     "subnet": "10.0.100.0/24",                              │
│     "dns_servers": ["1.1.1.1", "8.8.8.8"],                  │
│     "ntp_server": "pool.ntp.org"                             │
│   }                                                           │
│ }                                                             │
│                                                               │
│ Node: ✅ Registration ACK received                           │
│ Node: Saving Command Station config...                       │
│ Node: Status: REGISTERED                                     │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ STEP 4: Atualização do Inventário (Command Station)         │
├─────────────────────────────────────────────────────────────┤
│ Updating ~/.syntropy/nodes/node-01.yaml                      │
│                                                               │
│ Changes:                                                      │
│   status: waiting_provisioning → online                      │
│   ip: null → 10.0.100.45                                    │
│   registered_at: 2025-10-10T14:35:00Z                       │
│   hardware: (added from manifest)                            │
│     cpu_cores: 8                                             │
│     cpu_model: "AMD Ryzen 7 5800X"                          │
│     ram_gb: 28                                               │
│     disk_gb: 512                                             │
│                                                               │
│ ✅ Inventory updated                                         │
│ ✅ Node node-01 is now part of the Grid!                    │
└─────────────────────────────────────────────────────────────┘


╔════════════════════════════════════════════════════════════════╗
║  FASE 4: VERIFICAÇÃO NA COMMAND STATION                        ║
╚════════════════════════════════════════════════════════════════╝

$ syntropy node list

ID       IP            STATUS   CPU    RAM    DISK   UPTIME
node-01  10.0.100.45   online   8c     28GB   512GB  2m

$ syntropy node status node-01

╔═══════════════════════════════════════════════════════════════╗
║  NODE: node-01                                                 ║
╚═══════════════════════════════════════════════════════════════╝

🔌 Conectividade:
   IP: 10.0.100.45 (DHCP)
   Status: ✅ Online
   Last seen: 5 seconds ago
   
💻 Hardware (auto-detected):
   CPU: AMD Ryzen 7 5800X (8 cores)
   RAM: 28GB
   Disk: 512GB NVMe
   Network: Intel I225-V (00:1a:2b:3c:4d:5e)
   
🔐 Segurança:
   SSH: ✅ Key-based only
   Firewall: ✅ Active (UFW)
   Fail2ban: ✅ Active
   
📦 Workloads:
   Running: 0
   Total deployed: 0
   
⏰ Uptime:
   2 minutes
```

### 4.3 Implementação: Subcomponente USB Detection

**Arquivo**: `manager/interfaces/cli/node/create/usb_detector.go`

```go
// Package create handles Node creation via USB
package create

import (
	"fmt"
)

// USBDevice represents a USB storage device
type USBDevice struct {
	Path         string  // Device path (e.g., "E:\" or "/dev/sdb")
	Name         string  // Device name
	Vendor       string  // Vendor name
	Model        string  // Model name
	SizeGB       float64 // Size in GB
	Filesystem   string  // Filesystem type (if mounted)
	MountPoint   string  // Mount point (if mounted)
	HasData      bool    // Whether device contains data
	IsSystemDisk bool    // Safety: block system disks
	Removable    bool    // Whether device is removable
}

// USBDetector interface (multi-platform)
type USBDetector interface {
	// ListDevices lists all USB storage devices
	ListDevices() ([]*USBDevice, error)
	
	// ValidateDevice validates if device is safe to format
	ValidateDevice(device *USBDevice) error
	
	// FormatDevice formats device (WARNING: destructive)
	FormatDevice(device *USBDevice) error
}

// DetectUSBDevices detects USB devices (platform-agnostic wrapper)
func DetectUSBDevices() ([]*USBDevice, error) {
	detector := newPlatformUSBDetector()
	return detector.ListDevices()
}

// PromptUserSelection prompts user to select USB device
func PromptUserSelection(devices []*USBDevice) (*USBDevice, error) {
	if len(devices) == 0 {
		return nil, fmt.Errorf("no USB devices found")
	}
	
	fmt.Println("\n📀 Available USB Devices:")
	fmt.Println("─────────────────────────────────────────────────────────")
	
	for i, dev := range devices {
		statusIcon := "✅"
		statusText := "Empty"
		
		if dev.IsSystemDisk {
			statusIcon = "🚫"
			statusText = "BLOCKED (system disk)"
		} else if dev.HasData {
			statusIcon = "⚠️ "
			statusText = "Contains data"
		}
		
		fmt.Printf("[%d] %s (%s %s - %.0fGB) - %s %s\n",
			i+1,
			dev.Path,
			dev.Vendor,
			dev.Model,
			dev.SizeGB,
			statusIcon,
			statusText,
		)
	}
	
	fmt.Println("\n⚠️  WARNING: Selected USB will be FORMATTED!")
	fmt.Println("   All data will be PERMANENTLY LOST!")
	
	// Get user input
	var selection int
	fmt.Print("\nSelect USB [1-" + fmt.Sprintf("%d", len(devices)) + "]: ")
	fmt.Scanln(&selection)
	
	if selection < 1 || selection > len(devices) {
		return nil, fmt.Errorf("invalid selection")
	}
	
	selectedDevice := devices[selection-1]
	
	// Block system disks
	if selectedDevice.IsSystemDisk {
		return nil, fmt.Errorf("cannot use system disk - blocked for safety")
	}
	
	// Confirm if device has data
	if selectedDevice.HasData {
		fmt.Printf("\n⚠️  Device %s contains data!\n", selectedDevice.Path)
		fmt.Print("   Type 'DELETE ALL DATA' to confirm: ")
		
		var confirm string
		fmt.Scanln(&confirm)
		
		if confirm != "DELETE ALL DATA" {
			return nil, fmt.Errorf("operation cancelled")
		}
	} else {
		// Still confirm for empty devices
		fmt.Printf("\nConfirm format %s %s %.0fGB? (type YES): ",
			selectedDevice.Path,
			selectedDevice.Model,
			selectedDevice.SizeGB,
		)
		
		var confirm string
		fmt.Scanln(&confirm)
		
		if confirm != "YES" {
			return nil, fmt.Errorf("operation cancelled")
		}
	}
	
	return selectedDevice, nil
}
```

**Implementação Windows**: `usb_detector_windows.go`

```go
//go:build windows

package create

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

type windowsUSBDetector struct{}

func newPlatformUSBDetector() USBDetector {
	return &windowsUSBDetector{}
}

func (d *windowsUSBDetector) ListDevices() ([]*USBDevice, error) {
	// Use PowerShell para detectar USBs
	cmd := exec.Command("powershell", "-Command",
		"Get-Disk | Where-Object {$_.BusType -eq 'USB'} | ForEach-Object {"+
			"$partitions = Get-Partition -DiskNumber $_.Number;"+
			"$_.Number, $_.Model, $_.Size, $_.BusType, ($partitions | Select-Object -ExpandProperty DriveLetter)"+
		"}")
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list USB devices: %w\nOutput: %s", err, output)
	}
	
	var devices []*USBDevice
	lines := strings.Split(string(output), "\n")
	
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		
		// Parse output
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		
		// Extract info
		diskNumber := fields[0]
		model := fields[1]
		sizeBytes := parseSize(fields[2])
		driveLetter := extractDriveLetter(fields)
		
		// Check if drive is removable and not system
		isSystem := isSystemDisk(driveLetter)
		hasData := checkHasData(driveLetter)
		
		device := &USBDevice{
			Path:         driveLetter + ":\\",
			Name:         fmt.Sprintf("Disk %s", diskNumber),
			Model:        model,
			SizeGB:       float64(sizeBytes) / (1024 * 1024 * 1024),
			IsSystemDisk: isSystem,
			HasData:      hasData,
			Removable:    true,
		}
		
		devices = append(devices, device)
	}
	
	return devices, nil
}

func (d *windowsUSBDetector) ValidateDevice(device *USBDevice) error {
	if device.IsSystemDisk {
		return fmt.Errorf("cannot use system disk: %s", device.Path)
	}
	
	if device.SizeGB < 8 {
		return fmt.Errorf("USB too small (need ≥8GB): %.1fGB", device.SizeGB)
	}
	
	return nil
}

func (d *windowsUSBDetector) FormatDevice(device *USBDevice) error {
	// Use diskpart for formatting
	diskpartScript := fmt.Sprintf(`
		select volume %s
		clean
		create partition primary
		format fs=fat32 quick
		assign
		exit
	`, strings.TrimSuffix(device.Path, ":\\"))
	
	cmd := exec.Command("diskpart")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	
	// TODO: Implement safe formatting
	return fmt.Errorf("format not implemented yet")
}

// Helper functions
func isSystemDisk(driveLetter string) bool {
	// Block C:\ and system drives
	return driveLetter == "C" || driveLetter == ""
}

func checkHasData(driveLetter string) bool {
	if driveLetter == "" {
		return false
	}
	
	path := driveLetter + ":\\"
	cmd := exec.Command("cmd", "/C", "dir", path)
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		return false
	}
	
	// Check if directory listing shows files
	return !strings.Contains(string(output), "0 File(s)")
}

func parseSize(sizeStr string) int64 {
	// Parse size from PowerShell output
	// TODO: Implement proper parsing
	return 32 * 1024 * 1024 * 1024 // Default 32GB
}

func extractDriveLetter(fields []string) string {
	// Extract drive letter from fields
	for _, field := range fields {
		if len(field) == 1 && field >= "A" && field <= "Z" {
			return field
		}
	}
	return ""
}
```

**Implementação Linux**: `usb_detector_linux.go`

```go
//go:build linux

package create

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type linuxUSBDetector struct{}

func newPlatformUSBDetector() USBDetector {
	return &linuxUSBDetector{}
}

func (d *linuxUSBDetector) ListDevices() ([]*USBDevice, error) {
	// Use lsblk to detect USB devices
	cmd := exec.Command("lsblk", "-o", "NAME,SIZE,VENDOR,MODEL,TRAN,TYPE,MOUNTPOINT", "-J")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list USB devices: %w\nOutput: %s", err, output)
	}
	
	// Parse JSON output
	// TODO: Implement JSON parsing
	
	// Alternative: simple parsing
	cmd = exec.Command("lsblk", "-o", "NAME,SIZE,VENDOR,MODEL,TRAN,TYPE,MOUNTPOINT", "-n")
	output, err = cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list devices: %w", err)
	}
	
	var devices []*USBDevice
	lines := strings.Split(string(output), "\n")
	
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}
		
		name := fields[0]
		size := fields[1]
		vendor := fields[2]
		model := fields[3]
		transport := fields[4]
		devType := fields[5]
		mountPoint := ""
		if len(fields) > 6 {
			mountPoint = fields[6]
		}
		
		// Only USB disks (not partitions)
		if transport != "usb" || devType != "disk" {
			continue
		}
		
		path := "/dev/" + name
		sizeGB := parseSizeLinux(size)
		isSystem := isSystemDiskLinux(path)
		hasData := checkHasDataLinux(path)
		
		device := &USBDevice{
			Path:         path,
			Name:         name,
			Vendor:       vendor,
			Model:        model,
			SizeGB:       sizeGB,
			MountPoint:   mountPoint,
			IsSystemDisk: isSystem,
			HasData:      hasData,
			Removable:    true,
		}
		
		devices = append(devices, device)
	}
	
	return devices, nil
}

func (d *linuxUSBDetector) ValidateDevice(device *USBDevice) error {
	if device.IsSystemDisk {
		return fmt.Errorf("cannot use system disk: %s", device.Path)
	}
	
	if device.SizeGB < 8 {
		return fmt.Errorf("USB too small (need ≥8GB): %.1fGB", device.SizeGB)
	}
	
	// Check if device is writable
	info, err := os.Stat(device.Path)
	if err != nil {
		return fmt.Errorf("cannot access device: %w", err)
	}
	
	if info.Mode().Perm()&0200 == 0 {
		return fmt.Errorf("device is not writable (check permissions)")
	}
	
	return nil
}

func (d *linuxUSBDetector) FormatDevice(device *USBDevice) error {
	// Use dd and mkfs for formatting
	// WARNING: DESTRUCTIVE OPERATION
	
	// Unmount if mounted
	if device.MountPoint != "" {
		cmd := exec.Command("umount", device.MountPoint)
		cmd.Run() // Ignore errors
	}
	
	// TODO: Implement safe formatting
	return fmt.Errorf("format not implemented yet")
}

// Helper functions
func isSystemDiskLinux(path string) bool {
	// Block /dev/sda (usually system disk)
	// Block mounted root partitions
	return strings.HasPrefix(path, "/dev/sda") ||
		strings.HasPrefix(path, "/dev/nvme0n1")
}

func checkHasDataLinux(path string) bool {
	// Try to mount and check for files
	cmd := exec.Command("blkid", path)
	output, err := cmd.CombinedOutput()
	
	// If blkid returns filesystem info, device likely has data
	return err == nil && len(output) > 0
}

func parseSizeLinux(sizeStr string) float64 {
	// Parse size like "32G", "512M", etc.
	sizeStr = strings.TrimSpace(sizeStr)
	
	if sizeStr == "" {
		return 0
	}
	
	unit := sizeStr[len(sizeStr)-1]
	valueStr := sizeStr[:len(sizeStr)-1]
	
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return 0
	}
	
	switch unit {
	case 'T':
		return value * 1024
	case 'G':
		return value
	case 'M':
		return value / 1024
	case 'K':
		return value / (1024 * 1024)
	default:
		return value / (1024 * 1024 * 1024) // Assume bytes
	}
}
```

### 4.4 Implementação: Cloud-Init Generator

**Arquivo**: `manager/interfaces/cli/node/create/cloud_init_generator.go`

```go
package create

import (
	"fmt"
	"os"
	"path/filepath"
	
	infra "github.com/syntropy-cooperative-grid/infrastructure"
)

// CloudInitGenerator generates cloud-init configuration for Nodes
type CloudInitGenerator struct {
	templateManager *infra.TemplateManager
	keyManager      *infra.KeyManager
	configDir       string
}

// NewCloudInitGenerator creates a new cloud-init generator
func NewCloudInitGenerator(projectRoot, configDir string) *CloudInitGenerator {
	templatesPath := filepath.Join(projectRoot, "infrastructure", "cloud-init")
	keysPath := filepath.Join(configDir, "../keys/nodes")
	
	return &CloudInitGenerator{
		templateManager: infra.NewTemplateManager(templatesPath),
		keyManager:      infra.NewKeyManager(keysPath),
		configDir:       configDir,
	}
}

// GenerateForNode generates cloud-init files for a specific Node
func (g *CloudInitGenerator) GenerateForNode(nodeName string, outputDir string) error {
	fmt.Printf("⚙️  Generating cloud-init for %s...\n", nodeName)
	
	// STEP 1: Load Grid Token
	fmt.Println("  [1/7] Loading Grid Token...")
	gridToken, err := g.loadGridToken()
	if err != nil {
		return fmt.Errorf("failed to load grid token: %w", err)
	}
	fmt.Printf("  ✅ Token: %s...%s\n", gridToken[:8], "[HIDDEN]")
	
	// STEP 2: Generate SSH keys for Node
	fmt.Println("  [2/7] Generating Node SSH keys...")
	ownerKey, communityKey, err := g.generateNodeKeys(nodeName)
	if err != nil {
		return fmt.Errorf("failed to generate keys: %w", err)
	}
	fmt.Printf("  ✅ Owner key: %s\n", ownerKey.Fingerprint)
	fmt.Printf("  ✅ Community key: %s\n", communityKey.Fingerprint)
	
	// STEP 3: Load Command Station SSH public key
	fmt.Println("  [3/7] Loading Command Station SSH key...")
	stationKey, err := g.keyManager.LoadKeyPair(infra.OwnerKey, "command-station")
	if err != nil {
		return fmt.Errorf("failed to load command station key: %w", err)
	}
	fmt.Println("  ✅ Command Station key loaded")
	
	// STEP 4: Load configuration
	fmt.Println("  [4/7] Loading configuration...")
	commandStationIP, err := g.getCommandStationIP()
	if err != nil {
		return fmt.Errorf("failed to get command station IP: %w", err)
	}
	fmt.Printf("  ✅ Command Station IP: %s\n", commandStationIP)
	
	// STEP 5: Prepare template data
	fmt.Println("  [5/7] Preparing template variables...")
	data := &infra.TemplateData{
		NodeName:           nodeName,
		NodeDescription:    fmt.Sprintf("Syntropy Node %s", nodeName),
		Coordinates:        "0,0", // Default, will be updated after provisioning
		CreatedAt:          time.Now().Format(time.RFC3339),
		AdminPasswordHash:  "$6$rounds=4096$salt$hash", // Disabled (key-only auth)
		OwnerPublicKey:     ownerKey.PublicKey,
		CommunityPublicKey: communityKey.PublicKey,
	}
	
	// Add custom fields for our template
	additionalData := map[string]string{
		"GRID_TOKEN":            gridToken,
		"COMMAND_STATION_IP":    commandStationIP,
		"SSH_PUBLIC_KEY":        stationKey.PublicKey,
		"OWNER_KEY_FINGERPRINT": ownerKey.Fingerprint,
		"COMM_KEY_FINGERPRINT":  communityKey.Fingerprint,
	}
	
	fmt.Println("  ✅ Variables prepared")
	
	// STEP 6: Render templates
	fmt.Println("  [6/7] Rendering cloud-init templates...")
	
	// Render user-data
	userDataPath := filepath.Join(outputDir, "user-data")
	if err := g.renderUserData(data, additionalData, userDataPath); err != nil {
		return fmt.Errorf("failed to render user-data: %w", err)
	}
	fmt.Println("  ✅ user-data generated")
	
	// Render meta-data
	metaDataPath := filepath.Join(outputDir, "meta-data")
	if err := g.renderMetaData(nodeName, metaDataPath); err != nil {
		return fmt.Errorf("failed to render meta-data: %w", err)
	}
	fmt.Println("  ✅ meta-data generated")
	
	// Render network-config
	networkConfigPath := filepath.Join(outputDir, "network-config")
	if err := g.renderNetworkConfig(networkConfigPath); err != nil {
		return fmt.Errorf("failed to render network-config: %w", err)
	}
	fmt.Println("  ✅ network-config generated")
	
	// STEP 7: Save additional files (Agent, scripts, etc.)
	fmt.Println("  [7/7] Preparing additional files...")
	if err := g.prepareAgentBinary(outputDir); err != nil {
		return fmt.Errorf("failed to prepare agent: %w", err)
	}
	fmt.Println("  ✅ Agent binary ready")
	
	fmt.Println("\n✅ Cloud-init generation complete!")
	return nil
}

func (g *CloudInitGenerator) loadGridToken() (string, error) {
	tokenPath := filepath.Join(g.configDir, "../grid-token.txt")
	data, err := os.ReadFile(tokenPath)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func (g *CloudInitGenerator) generateNodeKeys(nodeName string) (*infra.KeyPair, *infra.KeyPair, error) {
	// Generate owner key (for SSH access)
	ownerKey, err := g.keyManager.GenerateKeyPair(infra.OwnerKey, nodeName)
	if err != nil {
		return nil, nil, err
	}
	
	if err := g.keyManager.SaveKeyPair(ownerKey, infra.OwnerKey, nodeName); err != nil {
		return nil, nil, err
	}
	
	// Generate community key (for inter-node communication)
	communityKey, err := g.keyManager.GenerateKeyPair(infra.CommunityKey, nodeName)
	if err != nil {
		return nil, nil, err
	}
	
	if err := g.keyManager.SaveKeyPair(communityKey, infra.CommunityKey, nodeName); err != nil {
		return nil, nil, err
	}
	
	return ownerKey, communityKey, nil
}

func (g *CloudInitGenerator) getCommandStationIP() (string, error) {
	// Read from config.yaml
	configPath := filepath.Join(g.configDir, "manager.yaml")
	// TODO: Parse YAML and extract command_station.ip
	// For now, return default
	return "10.0.100.1", nil
}

func (g *CloudInitGenerator) renderUserData(data *infra.TemplateData, additional map[string]string, outputPath string) error {
	// Load template
	content, err := g.templateManager.RenderTemplate("user-data-template.yaml", data)
	if err != nil {
		return err
	}
	
	// Replace additional placeholders
	for key, value := range additional {
		placeholder := "${" + key + "}"
		content = strings.ReplaceAll(content, placeholder, value)
	}
	
	// Save to output
	return os.WriteFile(outputPath, []byte(content), 0644)
}

func (g *CloudInitGenerator) renderMetaData(nodeName, outputPath string) error {
	// Simple meta-data
	content := fmt.Sprintf(`instance-id: %s-%s
local-hostname: syntropy-%s
`, nodeName, generateInstanceID(), nodeName)
	
	return os.WriteFile(outputPath, []byte(content), 0644)
}

func (g *CloudInitGenerator) renderNetworkConfig(outputPath string) error {
	// Use DHCP (network-config-template.yaml)
	return g.templateManager.SaveTemplate("network-config-template.yaml", nil, outputPath)
}

func (g *CloudInitGenerator) prepareAgentBinary(outputDir string) error {
	// TODO: Copy Agent binary to output directory
	// For now, create placeholder
	agentDir := filepath.Join(outputDir, "syntropy")
	if err := os.MkdirAll(agentDir, 0755); err != nil {
		return err
	}
	
	placeholderPath := filepath.Join(agentDir, "agent")
	placeholder := "#!/bin/bash\necho 'Syntropy Agent Placeholder'\n"
	return os.WriteFile(placeholderPath, []byte(placeholder), 0755)
}

func generateInstanceID() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}
```

### 4.5 Implementação: Registration Protocol

**Arquivo**: `manager/interfaces/cli/node/registration/registration.go`

```go
package registration

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// RegistrationProtocol handles Node registration
type RegistrationProtocol struct {
	gridToken         string
	commandStationIP  string
	inventoryManager  *InventoryManager
	listener          net.Listener
}

// NodeAnnouncement represents Node announcement message
type NodeAnnouncement struct {
	Type             string            `json:"type"`
	NodeName         string            `json:"node_name"`
	GridToken        string            `json:"grid_token"`
	IP               string            `json:"ip"`
	SSHPort          int               `json:"ssh_port"`
	PublicKey        string            `json:"public_key"`
	HardwareManifest *HardwareManifest `json:"hardware_manifest"`
	Timestamp        time.Time         `json:"timestamp"`
}

// HardwareManifest represents Node hardware information
type HardwareManifest struct {
	CPUCores        int                  `json:"cpu_cores"`
	CPUModel        string               `json:"cpu_model"`
	RAMGB           float64              `json:"ram_gb"`
	DiskGB          float64              `json:"disk_gb"`
	DiskType        string               `json:"disk_type"`
	NetworkInterfaces []NetworkInterface `json:"network_interfaces"`
}

// NetworkInterface represents a network interface
type NetworkInterface struct {
	Name       string `json:"name"`
	MAC        string `json:"mac"`
	Speed      string `json:"speed"`
	LinkStatus string `json:"link_status"`
}

// RegistrationAck represents Command Station acknowledgment
type RegistrationAck struct {
	Type                   string                 `json:"type"`
	Status                 string                 `json:"status"` // "accepted" or "rejected"
	NodeID                 string                 `json:"node_id"`
	CommandStationIP       string                 `json:"command_station_ip"`
	CommandStationPublicKey string                 `json:"command_station_public_key"`
	GridConfig             map[string]interface{} `json:"grid_config"`
	Reason                 string                 `json:"reason,omitempty"` // If rejected
}

// NewRegistrationProtocol creates a new registration protocol handler
func NewRegistrationProtocol(gridToken, commandStationIP string, inventoryMgr *InventoryManager) *RegistrationProtocol {
	return &RegistrationProtocol{
		gridToken:         gridToken,
		commandStationIP:  commandStationIP,
		inventoryManager:  inventoryMgr,
	}
}

// StartListener starts listening for Node announcements
func (rp *RegistrationProtocol) StartListener() error {
	// Listen on port 51000 for announcements
	listener, err := net.Listen("tcp", ":51000")
	if err != nil {
		return fmt.Errorf("failed to start listener: %w", err)
	}
	
	rp.listener = listener
	
	fmt.Println("🔊 Registration Protocol: Listening for Node announcements on :51000")
	
	go rp.acceptConnections()
	
	return nil
}

func (rp *RegistrationProtocol) acceptConnections() {
	for {
		conn, err := rp.listener.Accept()
		if err != nil {
			fmt.Printf("⚠️  Error accepting connection: %v\n", err)
			continue
		}
		
		go rp.handleAnnouncement(conn)
	}
}

func (rp *RegistrationProtocol) handleAnnouncement(conn net.Conn) {
	defer conn.Close()
	
	// Read announcement
	decoder := json.NewDecoder(conn)
	var announcement NodeAnnouncement
	
	if err := decoder.Decode(&announcement); err != nil {
		fmt.Printf("❌ Failed to decode announcement: %v\n", err)
		return
	}
	
	fmt.Printf("\n📢 Received announcement from %s (%s)\n", announcement.NodeName, announcement.IP)
	
	// Validate announcement
	ack, err := rp.validateAndRegister(&announcement)
	if err != nil {
		fmt.Printf("❌ Registration failed: %v\n", err)
		ack = &RegistrationAck{
			Type:   "registration_ack",
			Status: "rejected",
			Reason: err.Error(),
		}
	}
	
	// Send acknowledgment
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(ack); err != nil {
		fmt.Printf("⚠️  Failed to send ack: %v\n", err)
		return
	}
	
	if ack.Status == "accepted" {
		fmt.Printf("✅ Node %s registered successfully!\n", announcement.NodeName)
	}
}

func (rp *RegistrationProtocol) validateAndRegister(announcement *NodeAnnouncement) (*RegistrationAck, error) {
	// VALIDATION 1: Grid Token
	fmt.Println("  [CHECK 1] Validating Grid Token...")
	if announcement.GridToken != rp.gridToken {
		return nil, fmt.Errorf("invalid grid token")
	}
	fmt.Println("  ✅ Grid Token valid")
	
	// VALIDATION 2: Node exists in inventory
	fmt.Println("  [CHECK 2] Validating Node in inventory...")
	node, err := rp.inventoryManager.GetNode(announcement.NodeName)
	if err != nil {
		return nil, fmt.Errorf("node not in inventory: %w", err)
	}
	
	if node.Status != "waiting_provisioning" {
		return nil, fmt.Errorf("node already registered (status: %s)", node.Status)
	}
	fmt.Println("  ✅ Node found in inventory")
	
	// VALIDATION 3: SSH Public Key
	fmt.Println("  [CHECK 3] Validating SSH key...")
	expectedKey := node.SSHPublicKey
	if announcement.PublicKey != expectedKey {
		return nil, fmt.Errorf("SSH key mismatch")
	}
	fmt.Println("  ✅ SSH key valid")
	
	// REGISTER NODE
	fmt.Println("  [UPDATE] Registering Node...")
	node.Status = "online"
	node.IP = announcement.IP
	node.RegisteredAt = time.Now()
	node.Hardware = announcement.HardwareManifest
	
	if err := rp.inventoryManager.UpdateNode(node); err != nil {
		return nil, fmt.Errorf("failed to update inventory: %w", err)
	}
	fmt.Println("  ✅ Inventory updated")
	
	// Prepare ACK
	ack := &RegistrationAck{
		Type:                   "registration_ack",
		Status:                 "accepted",
		NodeID:                 announcement.NodeName,
		CommandStationIP:       rp.commandStationIP,
		CommandStationPublicKey: rp.getCommandStationPublicKey(),
		GridConfig: map[string]interface{}{
			"subnet":      "10.0.100.0/24",
			"dns_servers": []string{"1.1.1.1", "8.8.8.8"},
			"ntp_server":  "pool.ntp.org",
		},
	}
	
	return ack, nil
}

func (rp *RegistrationProtocol) getCommandStationPublicKey() string {
	// TODO: Load from ~/.syntropy/keys/command-station-owner.pub
	return "ssh-ed25519 AAAA... command-station"
}

// Stop stops the listener
func (rp *RegistrationProtocol) Stop() error {
	if rp.listener != nil {
		return rp.listener.Close()
	}
	return nil
}
```

---

## 5. SINCRONIZAÇÃO COMMAND STATION ↔ NODES

### 5.1 Protocolo de Sincronização

```
┌────────────────────────────────────────────────────────┐
│ SINCRONIZAÇÃO BIDIRECIONAL                             │
├────────────────────────────────────────────────────────┤
│                                                         │
│ Command Station              Node                      │
│       │                       │                         │
│       │─────── (1) Poll ─────►│  (a cada 30s)          │
│       │   "Send me updates"   │                         │
│       │                       │                         │
│       │◄─── (2) Manifest ─────│                         │
│       │   { hardware, status,│                         │
│       │     workloads, metrics }                        │
│       │                       │                         │
│       │─── (3) Commands ─────►│                         │
│       │   [ deploy, stop,    │                         │
│       │     restart, update ]│                         │
│       │                       │                         │
│       │◄───── (4) ACK ────────│                         │
│       │   { status: ok }     │                         │
│       │                       │                         │
└────────────────────────────────────────────────────────┘
```

### 5.2 Implementação: Sync Manager

**Arquivo**: `manager/interfaces/cli/management/sync/sync.go`

```go
package sync

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"
)

// SyncManager handles Command Station ↔ Node synchronization
type SyncManager struct {
	keyPath    string
	sshUser    string
	syncInterval time.Duration
}

// NodeStatus represents current Node status
type NodeStatus struct {
	NodeID            string                 `json:"node_id"`
	IP                string                 `json:"ip"`
	Status            string                 `json:"status"` // online, degraded, offline
	Timestamp         time.Time              `json:"timestamp"`
	Hardware          *HardwareManifest      `json:"hardware"`
	Workloads         []WorkloadStatus       `json:"workloads"`
	Metrics           *NodeMetrics           `json:"metrics"`
}

// WorkloadStatus represents status of a workload
type WorkloadStatus struct {
	ID        string    `json:"id"`
	Image     string    `json:"image"`
	Status    string    `json:"status"` // running, stopped, failed
	StartedAt time.Time `json:"started_at"`
}

// NodeMetrics represents Node metrics
type NodeMetrics struct {
	CPUPercent    float64 `json:"cpu_percent"`
	RAMUsedGB     float64 `json:"ram_used_gb"`
	RAMTotalGB    float64 `json:"ram_total_gb"`
	DiskUsedGB    float64 `json:"disk_used_gb"`
	DiskTotalGB   float64 `json:"disk_total_gb"`
	NetworkRxMB   float64 `json:"network_rx_mb"`
	NetworkTxMB   float64 `json:"network_tx_mb"`
	Uptime        string  `json:"uptime"`
}

// NewSyncManager creates a new sync manager
func NewSyncManager(keyPath, sshUser string, syncInterval time.Duration) *SyncManager {
	return &SyncManager{
		keyPath:      keyPath,
		sshUser:      sshUser,
		syncInterval: syncInterval,
	}
}

// PollNode polls Node for current status
func (sm *SyncManager) PollNode(nodeIP string) (*NodeStatus, error) {
	// Execute remote command to get Node status
	cmd := sm.buildSSHCommand(nodeIP, "/opt/syntropy/bin/syntropy-agent status --json")
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to poll node: %w\nOutput: %s", err, output)
	}
	
	// Parse JSON response
	var status NodeStatus
	if err := json.Unmarshal(output, &status); err != nil {
		return nil, fmt.Errorf("failed to parse status: %w", err)
	}
	
	return &status, nil
}

// SendCommand sends a command to Node
func (sm *SyncManager) SendCommand(nodeIP string, command string, args map[string]interface{}) error {
	// Serialize command
	payload, err := json.Marshal(map[string]interface{}{
		"command": command,
		"args":    args,
	})
	if err != nil {
		return fmt.Errorf("failed to serialize command: %w", err)
	}
	
	// Send command via SSH
	remoteCmd := fmt.Sprintf("/opt/syntropy/bin/syntropy-agent exec --json '%s'", string(payload))
	cmd := sm.buildSSHCommand(nodeIP, remoteCmd)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to execute command: %w\nOutput: %s", err, output)
	}
	
	// Parse ACK
	var ack map[string]interface{}
	if err := json.Unmarshal(output, &ack); err != nil {
		return fmt.Errorf("failed to parse ack: %w", err)
	}
	
	if ack["status"] != "ok" {
		return fmt.Errorf("command failed: %v", ack["error"])
	}
	
	return nil
}

// StartPeriodicSync starts periodic synchronization
func (sm *SyncManager) StartPeriodicSync(nodes []string, inventoryMgr *InventoryManager) {
	ticker := time.NewTicker(sm.syncInterval)
	
	go func() {
		for range ticker.C {
			for _, nodeID := range nodes {
				go sm.syncNode(nodeID, inventoryMgr)
			}
		}
	}()
}

func (sm *SyncManager) syncNode(nodeID string, inventoryMgr *InventoryManager) {
	// Get Node from inventory
	node, err := inventoryMgr.GetNode(nodeID)
	if err != nil {
		return
	}
	
	// Poll Node
	status, err := sm.PollNode(node.IP)
	if err != nil {
		// Mark Node as offline if poll fails
		node.Status = "offline"
		node.LastSeen = time.Now()
		inventoryMgr.UpdateNode(node)
		return
	}
	
	// Update inventory with latest status
	node.Status = status.Status
	node.LastSeen = time.Now()
	node.Hardware = status.Hardware
	node.Metrics = status.Metrics
	
	inventoryMgr.UpdateNode(node)
}

func (sm *SyncManager) buildSSHCommand(nodeIP, remoteCommand string) *exec.Cmd {
	return exec.Command("ssh",
		"-i", sm.keyPath,
		"-o", "StrictHostKeyChecking=no",
		"-o", "ConnectTimeout=10",
		fmt.Sprintf("%s@%s", sm.sshUser, nodeIP),
		remoteCommand,
	)
}
```

---

## 6. INTELLIGENT WORKLOAD ORCHESTRATION

### 6.1 Problema Identificado

**QUESTÃO CRÍTICA NÃO ABORDADA**: Como distribuir workloads automaticamente pela Grid sem especificar Nodes manualmente?

**Cenários Reais**:
```bash
# ❌ RUIM: Usuário tem que escolher Node
syntropy workload deploy nginx --node node-01

# ✅ BOM: Usuário especifica requisitos, Grid distribui
syntropy workload deploy nginx --replicas 3 --cpu 1 --memory 512M
# Grid decide automaticamente: node-01, node-03, node-05
```

**Desafios**:
1. **Validação de Capacidade**: Grid suporta essa carga?
2. **Distribuição Inteligente**: Quais Nodes usar?
3. **Sobrecarga**: O que fazer se Grid estiver cheia?
4. **Balanceamento**: Como distribuir uniformemente?
5. **Rejeição**: Como lidar com workloads impossíveis?

### 6.2 Arquitetura de Orquestração

```
┌─────────────────────────────────────────────────────────────┐
│  WORKLOAD DEPLOYMENT FLOW                                   │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  1. USER INPUT                                               │
│     syntropy workload deploy nginx \                         │
│       --replicas 3 \                                         │
│       --cpu 1 --memory 512M \                                │
│       --strategy spread                                      │
│                                                              │
│  2. ADMISSION CONTROL (Valida se é possível)                │
│     ├─ Total resources needed: 3 CPUs, 1.5GB RAM           │
│     ├─ Grid capacity check                                  │
│     ├─ Node availability check                              │
│     └─ Decision: ACCEPT / REJECT / QUEUE                    │
│                                                              │
│  3. SCHEDULER (Decide onde deployar)                        │
│     ├─ Strategy: spread / binpack / resource-optimized      │
│     ├─ Health filtering (only healthy nodes)                │
│     ├─ Resource scoring                                      │
│     └─ Placement: [node-01, node-03, node-05]              │
│                                                              │
│  4. DEPLOYMENT EXECUTOR                                      │
│     ├─ Deploy to node-01 (replica 1)                        │
│     ├─ Deploy to node-03 (replica 2)                        │
│     ├─ Deploy to node-05 (replica 3)                        │
│     └─ Monitor deployment status                            │
│                                                              │
│  5. RESULT                                                   │
│     ✅ Workload nginx deployed: 3/3 replicas running        │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 6.3 Modos de Deployment

#### 6.3.1 Grid-Wide Deployment (Padrão)

**Uso**: Quando usuário NÃO especifica Nodes.

```bash
# Deploy em TODA a Grid (Grid decide distribuição)
syntropy workload deploy nginx --replicas 3 --cpu 1 --memory 512M
```

**Comportamento**:
- Grid calcula: 3 réplicas × (1 CPU + 512M) = 3 CPUs + 1.5GB total
- Valida se Grid tem recursos
- Scheduler escolhe 3 Nodes automaticamente
- Distribui uniformemente (strategy: spread)

#### 6.3.2 Node-Specific Deployment

**Uso**: Quando usuário especifica Node(s).

```bash
# Deploy em Node específico
syntropy workload deploy postgres --node node-01 --cpu 2 --memory 4G

# Deploy em subset de Nodes
syntropy workload deploy worker --nodes node-02,node-03,node-04 --replicas 3
```

**Comportamento**:
- Valida apenas os Nodes especificados
- Ignora outros Nodes
- Scheduler limitado ao subset

#### 6.3.3 Constraint-Based Deployment

**Uso**: Deploy com restrições (tags, labels, capacidades).

```bash
# Deploy apenas em Nodes com GPU
syntropy workload deploy ml-training --tag gpu --replicas 2

# Deploy em Nodes com SSD
syntropy workload deploy database --disk-type ssd --replicas 1

# Deploy em Nodes com mínimo de RAM
syntropy workload deploy cache --min-ram 16G --replicas 2
```

### 6.4 Admission Control (Validação de Capacidade)

**Objetivo**: Validar se workload é POSSÍVEL antes de tentar deployar.

#### 6.4.1 Resource Validation

**Arquivo**: `workload/admission/admission_controller.go`

```go
package admission

import (
	"fmt"
)

// AdmissionController valida workloads antes de aceitar
type AdmissionController struct {
	inventoryMgr *InventoryManager
}

// WorkloadRequest representa requisição de deployment
type WorkloadRequest struct {
	Image         string
	Replicas      int
	CPUPerReplica float64 // Cores
	RAMPerReplica float64 // GB
	Strategy      string  // spread, binpack, resource-optimized
	Constraints   []Constraint
}

// Constraint representa restrição de placement
type Constraint struct {
	Type  string // tag, disk-type, min-ram, min-cpu
	Value string
}

// AdmissionResult representa resultado da validação
type AdmissionResult struct {
	Admitted      bool
	Reason        string
	GridCapacity  *CapacityInfo
	Recommendation string
}

// CapacityInfo representa capacidade da Grid
type CapacityInfo struct {
	TotalNodes       int
	HealthyNodes     int
	TotalCPU         float64
	AvailableCPU     float64
	TotalRAM         float64
	AvailableRAM     float64
	CPUUtilization   float64 // Percentage
	RAMUtilization   float64 // Percentage
}

// Validate valida se workload pode ser aceito
func (ac *AdmissionController) Validate(req *WorkloadRequest) (*AdmissionResult, error) {
	fmt.Printf("🔍 Validating workload request...\n")
	fmt.Printf("   Image: %s\n", req.Image)
	fmt.Printf("   Replicas: %d\n", req.Replicas)
	fmt.Printf("   Resources per replica: %.1f CPU, %.1f GB RAM\n", 
		req.CPUPerReplica, req.RAMPerReplica)
	
	// STEP 1: Calcular recursos totais necessários
	totalCPU := float64(req.Replicas) * req.CPUPerReplica
	totalRAM := float64(req.Replicas) * req.RAMPerReplica
	
	fmt.Printf("\n📊 Total resources needed:\n")
	fmt.Printf("   CPU: %.1f cores\n", totalCPU)
	fmt.Printf("   RAM: %.1f GB\n", totalRAM)
	
	// STEP 2: Obter capacidade atual da Grid
	capacity, err := ac.getGridCapacity()
	if err != nil {
		return nil, fmt.Errorf("failed to get grid capacity: %w", err)
	}
	
	fmt.Printf("\n📈 Grid capacity:\n")
	fmt.Printf("   Nodes: %d total, %d healthy\n", capacity.TotalNodes, capacity.HealthyNodes)
	fmt.Printf("   CPU: %.1f total, %.1f available (%.0f%% used)\n", 
		capacity.TotalCPU, capacity.AvailableCPU, capacity.CPUUtilization)
	fmt.Printf("   RAM: %.1f GB total, %.1f GB available (%.0f%% used)\n", 
		capacity.TotalRAM, capacity.AvailableRAM, capacity.RAMUtilization)
	
	// STEP 3: Validações
	
	// 3.1: Validar Nodes suficientes
	if capacity.HealthyNodes < req.Replicas {
		return &AdmissionResult{
			Admitted: false,
			Reason: fmt.Sprintf(
				"Not enough healthy nodes. Need %d, have %d", 
				req.Replicas, capacity.HealthyNodes,
			),
			GridCapacity: capacity,
			Recommendation: fmt.Sprintf(
				"Reduce replicas to %d or provision more nodes", 
				capacity.HealthyNodes,
			),
		}, nil
	}
	
	// 3.2: Validar CPU disponível
	if totalCPU > capacity.AvailableCPU {
		return &AdmissionResult{
			Admitted: false,
			Reason: fmt.Sprintf(
				"Insufficient CPU. Need %.1f cores, have %.1f available", 
				totalCPU, capacity.AvailableCPU,
			),
			GridCapacity: capacity,
			Recommendation: fmt.Sprintf(
				"Reduce CPU to %.1f per replica or reduce replicas to %d", 
				capacity.AvailableCPU / float64(req.Replicas),
				int(capacity.AvailableCPU / req.CPUPerReplica),
			),
		}, nil
	}
	
	// 3.3: Validar RAM disponível
	if totalRAM > capacity.AvailableRAM {
		return &AdmissionResult{
			Admitted: false,
			Reason: fmt.Sprintf(
				"Insufficient RAM. Need %.1f GB, have %.1f GB available", 
				totalRAM, capacity.AvailableRAM,
			),
			GridCapacity: capacity,
			Recommendation: fmt.Sprintf(
				"Reduce RAM to %.1f GB per replica or reduce replicas to %d", 
				capacity.AvailableRAM / float64(req.Replicas),
				int(capacity.AvailableRAM / req.RAMPerReplica),
			),
		}, nil
	}
	
	// 3.4: Validar constraints (tags, disk-type, etc.)
	if len(req.Constraints) > 0 {
		constraintResult := ac.validateConstraints(req.Constraints)
		if !constraintResult.Valid {
			return &AdmissionResult{
				Admitted: false,
				Reason: constraintResult.Reason,
				GridCapacity: capacity,
				Recommendation: constraintResult.Recommendation,
			}, nil
		}
	}
	
	// STEP 4: Validar limites de utilização (evitar sobrecarga)
	// Não permitir deploy se Grid ficaria >90% utilizada
	
	projectedCPU := capacity.CPUUtilization + ((totalCPU / capacity.TotalCPU) * 100)
	projectedRAM := capacity.RAMUtilization + ((totalRAM / capacity.TotalRAM) * 100)
	
	if projectedCPU > 90.0 {
		return &AdmissionResult{
			Admitted: false,
			Reason: fmt.Sprintf(
				"Deployment would overload Grid. CPU utilization would reach %.0f%% (limit: 90%%)", 
				projectedCPU,
			),
			GridCapacity: capacity,
			Recommendation: "Wait for current workloads to complete or scale down other workloads",
		}, nil
	}
	
	if projectedRAM > 90.0 {
		return &AdmissionResult{
			Admitted: false,
			Reason: fmt.Sprintf(
				"Deployment would overload Grid. RAM utilization would reach %.0f%% (limit: 90%%)", 
				projectedRAM,
			),
			GridCapacity: capacity,
			Recommendation: "Wait for current workloads to complete or scale down other workloads",
		}, nil
	}
	
	// STEP 5: ACEITO!
	fmt.Printf("\n✅ Workload ADMITTED\n")
	fmt.Printf("   Projected CPU utilization: %.0f%%\n", projectedCPU)
	fmt.Printf("   Projected RAM utilization: %.0f%%\n", projectedRAM)
	
	return &AdmissionResult{
		Admitted:     true,
		Reason:       "Workload meets all requirements",
		GridCapacity: capacity,
	}, nil
}

func (ac *AdmissionController) getGridCapacity() (*CapacityInfo, error) {
	nodes, err := ac.inventoryMgr.ListNodes()
	if err != nil {
		return nil, err
	}
	
	capacity := &CapacityInfo{
		TotalNodes: len(nodes),
	}
	
	for _, node := range nodes {
		// Contar apenas Nodes healthy
		if node.Status == "online" {
			capacity.HealthyNodes++
			
			// Somar recursos totais
			capacity.TotalCPU += float64(node.Hardware.CPUCores)
			capacity.TotalRAM += node.Hardware.RAMGB
			
			// Calcular recursos disponíveis (total - usado)
			cpuUsed := (node.Metrics.CPUPercent / 100) * float64(node.Hardware.CPUCores)
			capacity.AvailableCPU += float64(node.Hardware.CPUCores) - cpuUsed
			capacity.AvailableRAM += node.Hardware.RAMGB - node.Metrics.RAMUsedGB
		}
	}
	
	// Calcular utilização geral da Grid
	if capacity.TotalCPU > 0 {
		capacity.CPUUtilization = ((capacity.TotalCPU - capacity.AvailableCPU) / capacity.TotalCPU) * 100
	}
	if capacity.TotalRAM > 0 {
		capacity.RAMUtilization = ((capacity.TotalRAM - capacity.AvailableRAM) / capacity.TotalRAM) * 100
	}
	
	return capacity, nil
}

func (ac *AdmissionController) validateConstraints(constraints []Constraint) *ConstraintValidationResult {
	// TODO: Implementar validação de constraints
	// - Tags (ex: gpu, ssd)
	// - Disk type
	// - Min RAM/CPU
	
	return &ConstraintValidationResult{
		Valid: true,
	}
}

type ConstraintValidationResult struct {
	Valid          bool
	Reason         string
	Recommendation string
}
```

### 6.5 Intelligent Scheduler

**Objetivo**: Decidir ONDE deployar cada réplica.

#### 6.5.1 Estratégias de Placement

**Arquivo**: `workload/scheduler/scheduler.go`

```go
package scheduler

import (
	"fmt"
	"sort"
)

// Scheduler decide placement de workloads
type Scheduler struct {
	inventoryMgr *InventoryManager
}

// PlacementDecision representa decisão de placement
type PlacementDecision struct {
	NodeID       string
	ReplicaIndex int
	Score        float64
	Reason       string
}

// Schedule decide placement de workload
func (s *Scheduler) Schedule(req *WorkloadRequest) ([]*PlacementDecision, error) {
	fmt.Printf("\n🎯 Scheduling %d replicas...\n", req.Replicas)
	
	// STEP 1: Filtrar Nodes elegíveis
	nodes, err := s.getEligibleNodes(req)
	if err != nil {
		return nil, err
	}
	
	fmt.Printf("   Eligible nodes: %d\n", len(nodes))
	
	if len(nodes) < req.Replicas {
		return nil, fmt.Errorf("not enough eligible nodes: need %d, have %d", req.Replicas, len(nodes))
	}
	
	// STEP 2: Aplicar estratégia de placement
	var decisions []*PlacementDecision
	
	switch req.Strategy {
	case "spread":
		decisions = s.strategySpread(nodes, req)
	case "binpack":
		decisions = s.strategyBinpack(nodes, req)
	case "resource-optimized":
		decisions = s.strategyResourceOptimized(nodes, req)
	default:
		decisions = s.strategySpread(nodes, req) // Default
	}
	
	// STEP 3: Log decisões
	fmt.Printf("\n📍 Placement decisions:\n")
	for _, decision := range decisions {
		fmt.Printf("   Replica %d → %s (score: %.2f) - %s\n", 
			decision.ReplicaIndex, decision.NodeID, decision.Score, decision.Reason)
	}
	
	return decisions, nil
}

// strategySpread: Distribuir uniformemente pela Grid
func (s *Scheduler) strategySpread(nodes []*Node, req *WorkloadRequest) []*PlacementDecision {
	// Ordenar Nodes por MENOR carga atual (workloads count)
	sort.Slice(nodes, func(i, j int) bool {
		return len(nodes[i].Workloads) < len(nodes[j].Workloads)
	})
	
	decisions := make([]*PlacementDecision, req.Replicas)
	
	for i := 0; i < req.Replicas; i++ {
		node := nodes[i % len(nodes)]
		
		decisions[i] = &PlacementDecision{
			NodeID:       node.ID,
			ReplicaIndex: i + 1,
			Score:        calculateSpreadScore(node),
			Reason:       fmt.Sprintf("Spread strategy - least loaded (%d workloads)", len(node.Workloads)),
		}
		
		// Simular adição de workload para próxima iteração
		node.Workloads = append(node.Workloads, WorkloadInfo{})
	}
	
	return decisions
}

// strategyBinpack: Preencher Nodes até capacidade antes de usar próximo
func (s *Scheduler) strategyBinpack(nodes []*Node, req *WorkloadRequest) []*PlacementDecision {
	// Ordenar Nodes por MAIOR utilização (preencher os que já têm carga)
	sort.Slice(nodes, func(i, j int) bool {
		utilizationI := nodes[i].Metrics.CPUPercent + nodes[i].Metrics.RAMUtilization
		utilizationJ := nodes[j].Metrics.CPUPercent + nodes[j].Metrics.RAMUtilization
		return utilizationI > utilizationJ
	})
	
	decisions := make([]*PlacementDecision, req.Replicas)
	nodeIndex := 0
	
	for i := 0; i < req.Replicas; i++ {
		// Verificar se Node atual ainda tem capacidade
		for !s.nodeHasCapacity(nodes[nodeIndex], req.CPUPerReplica, req.RAMPerReplica) {
			nodeIndex++
			if nodeIndex >= len(nodes) {
				// Sem Nodes com capacidade suficiente
				return decisions[:i] // Retornar parcial
			}
		}
		
		node := nodes[nodeIndex]
		
		decisions[i] = &PlacementDecision{
			NodeID:       node.ID,
			ReplicaIndex: i + 1,
			Score:        calculateBinpackScore(node),
			Reason:       fmt.Sprintf("Binpack strategy - fill node (%.0f%% CPU, %.0f%% RAM)", 
				node.Metrics.CPUPercent, node.Metrics.RAMUtilization),
		}
		
		// Atualizar métricas simuladas
		node.Metrics.CPUPercent += (req.CPUPerReplica / float64(node.Hardware.CPUCores)) * 100
		node.Metrics.RAMUsedGB += req.RAMPerReplica
	}
	
	return decisions
}

// strategyResourceOptimized: Balancear CPU e RAM de forma otimizada
func (s *Scheduler) strategyResourceOptimized(nodes []*Node, req *WorkloadRequest) []*PlacementDecision {
	decisions := make([]*PlacementDecision, req.Replicas)
	
	for i := 0; i < req.Replicas; i++ {
		// Calcular score de otimização para cada Node
		scores := make([]NodeScore, len(nodes))
		
		for j, node := range nodes {
			scores[j] = NodeScore{
				Node:  node,
				Score: s.calculateResourceScore(node, req),
			}
		}
		
		// Ordenar por melhor score
		sort.Slice(scores, func(i, j int) bool {
			return scores[i].Score > scores[j].Score
		})
		
		bestNode := scores[0].Node
		
		decisions[i] = &PlacementDecision{
			NodeID:       bestNode.ID,
			ReplicaIndex: i + 1,
			Score:        scores[0].Score,
			Reason:       fmt.Sprintf("Resource-optimized - best fit (score: %.2f)", scores[0].Score),
		}
		
		// Simular alocação
		bestNode.Metrics.CPUPercent += (req.CPUPerReplica / float64(bestNode.Hardware.CPUCores)) * 100
		bestNode.Metrics.RAMUsedGB += req.RAMPerReplica
	}
	
	return decisions
}

// calculateResourceScore: Score baseado em balanceamento de recursos
func (s *Scheduler) calculateResourceScore(node *Node, req *WorkloadRequest) float64 {
	// Calcular utilização projetada
	cpuUtilization := node.Metrics.CPUPercent + ((req.CPUPerReplica / float64(node.Hardware.CPUCores)) * 100)
	ramUtilization := ((node.Metrics.RAMUsedGB + req.RAMPerReplica) / node.Hardware.RAMGB) * 100
	
	// Preferir Nodes com utilização balanceada
	balanceDiff := abs(cpuUtilization - ramUtilization)
	
	// Score maior = melhor balanceamento
	// Penalizar Nodes muito cheios ou muito vazios
	score := 100.0
	score -= balanceDiff * 0.5 // Penalizar desbalanceamento
	score -= cpuUtilization * 0.3 // Preferir Nodes menos utilizados
	score -= ramUtilization * 0.2
	
	// Bônus se Node está no sweet spot (40-60% utilizado)
	avgUtilization := (cpuUtilization + ramUtilization) / 2
	if avgUtilization >= 40 && avgUtilization <= 60 {
		score += 20
	}
	
	return score
}

func (s *Scheduler) nodeHasCapacity(node *Node, cpuNeeded, ramNeeded float64) bool {
	cpuAvailable := float64(node.Hardware.CPUCores) * ((100 - node.Metrics.CPUPercent) / 100)
	ramAvailable := node.Hardware.RAMGB - node.Metrics.RAMUsedGB
	
	return cpuAvailable >= cpuNeeded && ramAvailable >= ramNeeded
}

func (s *Scheduler) getEligibleNodes(req *WorkloadRequest) ([]*Node, error) {
	allNodes, err := s.inventoryMgr.ListNodes()
	if err != nil {
		return nil, err
	}
	
	var eligible []*Node
	
	for _, node := range allNodes {
		// Filtro 1: Apenas Nodes online
		if node.Status != "online" {
			continue
		}
		
		// Filtro 2: Constraints (tags, disk-type, etc.)
		if !s.matchesConstraints(node, req.Constraints) {
			continue
		}
		
		// Filtro 3: Capacidade mínima (pelo menos 1 réplica)
		if !s.nodeHasCapacity(node, req.CPUPerReplica, req.RAMPerReplica) {
			continue
		}
		
		eligible = append(eligible, node)
	}
	
	return eligible, nil
}

func (s *Scheduler) matchesConstraints(node *Node, constraints []Constraint) bool {
	for _, c := range constraints {
		switch c.Type {
		case "tag":
			if !nodeHasTag(node, c.Value) {
				return false
			}
		case "disk-type":
			if node.Hardware.DiskType != c.Value {
				return false
			}
		case "min-ram":
			// TODO: Parse c.Value and compare
		}
	}
	return true
}

type NodeScore struct {
	Node  *Node
	Score float64
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func calculateSpreadScore(node *Node) float64 {
	// Score inversamente proporcional ao número de workloads
	return 100.0 / (float64(len(node.Workloads)) + 1)
}

func calculateBinpackScore(node *Node) float64 {
	// Score proporcional à utilização
	return (node.Metrics.CPUPercent + node.Metrics.RAMUtilization) / 2
}

func nodeHasTag(node *Node, tag string) bool {
	for _, t := range node.Tags {
		if t == tag {
			return true
		}
	}
	return false
}
```

### 6.6 Exemplos Práticos

#### 6.6.1 Deploy Grid-Wide com Validação

```bash
# Tentar deployar workload muito grande
$ syntropy workload deploy heavy-app --replicas 10 --cpu 4 --memory 8G

🔍 Validating workload request...
   Image: heavy-app
   Replicas: 10
   Resources per replica: 4.0 CPU, 8.0 GB RAM

📊 Total resources needed:
   CPU: 40.0 cores
   RAM: 80.0 GB

📈 Grid capacity:
   Nodes: 6 total, 6 healthy
   CPU: 48.0 total, 30.0 available (37% used)
   RAM: 168.0 GB total, 120.0 GB available (29% used)

❌ ADMISSION DENIED

Reason: Deployment would overload Grid. RAM utilization would reach 77% (limit: 90%)

Recommendation: Reduce replicas to 7 or reduce RAM to 6.0 GB per replica

Grid Status:
  CPU: 37% used
  RAM: 29% used
  Projected after deployment: CPU 74%, RAM 77%
```

#### 6.6.2 Deploy com Auto-Ajuste

```bash
# CLI sugere ajuste automático
$ syntropy workload deploy app --replicas 5 --cpu 2 --memory 4G --auto-adjust

⚠️  Workload as specified would be REJECTED (insufficient RAM)

🔧 Auto-adjusting parameters...

Option 1: Reduce replicas
  Replicas: 5 → 4
  Resources: 2 CPU, 4 GB RAM (unchanged)
  
Option 2: Reduce RAM per replica
  Replicas: 5 (unchanged)
  Resources: 2 CPU, 4 GB → 3.5 GB RAM
  
Select option [1/2] or cancel: 1

✅ Adjusted deployment:
   Replicas: 4
   Resources: 2 CPU, 4 GB RAM per replica
   
Proceed? [y/N]: y

🎯 Scheduling 4 replicas...
   Eligible nodes: 6

📍 Placement decisions:
   Replica 1 → node-01 (score: 85.3) - Spread strategy
   Replica 2 → node-03 (score: 82.1) - Spread strategy
   Replica 3 → node-02 (score: 79.8) - Spread strategy
   Replica 4 → node-05 (score: 77.5) - Spread strategy

✅ Workload deployed: 4/4 replicas running
```

#### 6.6.3 Deploy com Constraints

```bash
# Deploy apenas em Nodes com SSD
$ syntropy workload deploy database \
    --replicas 2 \
    --cpu 2 --memory 8G \
    --disk-type ssd \
    --strategy binpack

🔍 Validating workload request...
   Constraints: disk-type=ssd

📊 Eligible nodes after filtering:
   Total nodes: 6
   After constraints: 3 (node-01, node-03, node-06 have SSD)
   
✅ Workload ADMITTED

🎯 Scheduling 2 replicas (strategy: binpack)...

📍 Placement decisions:
   Replica 1 → node-01 (score: 68.0) - Binpack - fill node (45% CPU, 38% RAM)
   Replica 2 → node-01 (score: 78.0) - Binpack - fill node (78% CPU, 67% RAM)
   
✅ Workload deployed: 2/2 replicas running on node-01
```

### 6.7 Queue System (Workloads Pendentes)

**Problema**: O que fazer quando Grid está cheia mas workload é válido?

**Solução**: Sistema de fila para workloads pendentes.

```bash
# Grid está 95% utilizada
$ syntropy workload deploy app --replicas 3 --cpu 1 --memory 1G

❌ ADMISSION DENIED: Grid would be overloaded (projected CPU: 97%)

💡 Options:
  1. QUEUE workload (deploy when resources available)
  2. FORCE deploy (override 90% limit - NOT RECOMMENDED)
  3. CANCEL
  
Select option [1/2/3]: 1

✅ Workload QUEUED
   Position: 2 in queue
   Estimated wait: ~30 minutes (when current jobs complete)
   
Track status: syntropy workload queue status
```

**Implementação**:

```go
// QueueManager gerencia workloads pendentes
type QueueManager struct {
	queue []*QueuedWorkload
}

type QueuedWorkload struct {
	ID           string
	Request      *WorkloadRequest
	QueuedAt     time.Time
	Priority     int
	EstimatedWait time.Duration
}

func (qm *QueueManager) Enqueue(req *WorkloadRequest) (*QueuedWorkload, error) {
	queued := &QueuedWorkload{
		ID:       generateID(),
		Request:  req,
		QueuedAt: time.Now(),
		Priority: 0, // FIFO por padrão
	}
	
	qm.queue = append(qm.queue, queued)
	
	// Estimar tempo de espera baseado em workloads atuais
	queued.EstimatedWait = qm.estimateWaitTime(req)
	
	return queued, nil
}

// ProcessQueue: Executar periodicamente (a cada 1 minuto)
func (qm *QueueManager) ProcessQueue(admissionCtrl *AdmissionController) {
	for i := 0; i < len(qm.queue); i++ {
		queued := qm.queue[i]
		
		// Tentar admitir novamente
		result, _ := admissionCtrl.Validate(queued.Request)
		
		if result.Admitted {
			// Deploy automaticamente
			fmt.Printf("✅ Queued workload %s now deploying...\n", queued.ID)
			
			// TODO: Deployar
			
			// Remover da fila
			qm.queue = append(qm.queue[:i], qm.queue[i+1:]...)
			i--
		}
	}
}
```

### 6.8 CLI Commands Atualizados

```bash
# Deploy com validação inteligente
syntropy workload deploy <image> \
  --replicas <n> \
  --cpu <cores> \
  --memory <size> \
  --strategy <spread|binpack|resource-optimized> \
  [--node <node-id>] \
  [--nodes <node1,node2>] \
  [--tag <tag>] \
  [--disk-type <ssd|hdd|nvme>] \
  [--min-ram <size>] \
  [--auto-adjust] \
  [--queue-if-full]

# Ver capacidade da Grid
syntropy grid capacity

# Ver workloads na fila
syntropy workload queue list

# Forçar deploy (override limites)
syntropy workload deploy <image> --force --ignore-limits
```

### 6.9 Integração com Deployment Workflow

**Fluxo Completo**:

```go
// workload/deploy/deploy.go

func (d *Deployer) Deploy(req *WorkloadRequest) error {
	// STEP 1: Admission Control
	admissionResult, err := d.admissionCtrl.Validate(req)
	if err != nil {
		return err
	}
	
	if !admissionResult.Admitted {
		fmt.Printf("❌ ADMISSION DENIED\n")
		fmt.Printf("   Reason: %s\n", admissionResult.Reason)
		fmt.Printf("   Recommendation: %s\n", admissionResult.Recommendation)
		
		// Oferecer fila
		if req.QueueIfFull {
			return d.queueMgr.Enqueue(req)
		}
		
		return fmt.Errorf("workload rejected: %s", admissionResult.Reason)
	}
	
	fmt.Printf("✅ ADMISSION APPROVED\n")
	
	// STEP 2: Scheduling
	placements, err := d.scheduler.Schedule(req)
	if err != nil {
		return err
	}
	
	// STEP 3: Deployment Execution
	for _, placement := range placements {
		if err := d.deployToNode(placement, req); err != nil {
			// Rollback em caso de falha
			d.rollback(placements[:placement.ReplicaIndex])
			return err
		}
	}
	
	fmt.Printf("✅ Workload deployed: %d/%d replicas running\n", 
		len(placements), req.Replicas)
	
	return nil
}
```

---

## 7. COMPONENTE WORKLOAD - ESTRUTURA DETALHADA

### 7.1 Arquitetura do Componente Workload

```
workload/
├── README.md                  # Documentação: O que é, como usar
├── ARCHITECTURE.md            # Arquitetura: Fluxo de deployment
├── workload.go                # Orquestrador principal (< 500 linhas)
│
├── admission/                 # Subcomponente: Admission Control
│   ├── README.md              # Documentação do subcomponente
│   ├── admission_controller.go   # Controller principal (< 400 linhas)
│   ├── capacity_calculator.go    # Cálculo de capacidade Grid (< 300 linhas)
│   ├── constraint_validator.go   # Validação de constraints (< 300 linhas)
│   ├── resource_validator.go     # Validação de recursos (< 300 linhas)
│   └── tests/
│       ├── admission_test.go
│       └── capacity_test.go
│
├── scheduler/                 # Subcomponente: Intelligent Scheduler
│   ├── README.md
│   ├── scheduler.go           # Scheduler principal (< 400 linhas)
│   ├── node_filter.go         # Filtragem de Nodes elegíveis (< 300 linhas)
│   ├── node_scorer.go         # Cálculo de scores (< 300 linhas)
│   ├── strategy_spread.go     # Estratégia spread (< 300 linhas)
│   ├── strategy_binpack.go    # Estratégia binpack (< 300 linhas)
│   ├── strategy_optimized.go  # Estratégia resource-optimized (< 400 linhas)
│   └── tests/
│       ├── scheduler_test.go
│       ├── spread_test.go
│       ├── binpack_test.go
│       └── optimized_test.go
│
├── queue/                     # Subcomponente: Queue Management
│   ├── README.md
│   ├── queue_manager.go       # Gerenciador de fila (< 300 linhas)
│   ├── queue_processor.go     # Processador periódico (< 300 linhas)
│   ├── wait_estimator.go      # Estimativa de tempo (< 200 linhas)
│   ├── priority_manager.go    # Gerenciamento de prioridades (< 200 linhas)
│   └── tests/
│       └── queue_test.go
│
├── deploy/                    # Subcomponente: Deployment Execution
│   ├── README.md
│   ├── deployer.go            # Orquestrador de deployment (< 400 linhas)
│   ├── executor.go            # Execução via SSH/Docker (< 400 linhas)
│   ├── executor_windows.go    # Implementação Windows (< 300 linhas)
│   ├── executor_linux.go      # Implementação Linux (< 300 linhas)
│   ├── rollback.go            # Rollback em falhas (< 300 linhas)
│   ├── docker_client.go       # Cliente Docker (< 400 linhas)
│   └── tests/
│       ├── deployer_test.go
│       └── rollback_test.go
│
├── lifecycle/                 # Subcomponente: Lifecycle Management
│   ├── README.md
│   ├── lifecycle.go           # Manager de lifecycle (< 300 linhas)
│   ├── start.go               # Start workload (< 200 linhas)
│   ├── stop.go                # Stop workload (< 200 linhas)
│   ├── restart.go             # Restart workload (< 200 linhas)
│   ├── scale.go               # Scale workload (< 300 linhas)
│   └── tests/
│       └── lifecycle_test.go
│
└── monitoring/                # Subcomponente: Monitoring
    ├── README.md
    ├── monitoring.go          # Monitor principal (< 300 linhas)
    ├── logs.go                # Agregação de logs (< 400 linhas)
    ├── metrics.go             # Coleta de métricas (< 400 linhas)
    └── tests/
        └── monitoring_test.go
```

### 7.2 Fluxo de Integração dos Subcomponentes

```
┌─────────────────────────────────────────────────────────────┐
│  WORKLOAD COMPONENT - INTEGRATION FLOW                      │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  User Command:                                               │
│  $ syntropy workload deploy nginx --replicas 3               │
│                                                              │
│         ↓                                                    │
│  ┌──────────────────────────────────────────┐              │
│  │  workload.go (Orquestrador Principal)    │              │
│  │  - Parse CLI arguments                    │              │
│  │  - Create WorkloadRequest                 │              │
│  │  - Coordena subcomponentes                │              │
│  └──────────────────────────────────────────┘              │
│         ↓                                                    │
│  ┌──────────────────────────────────────────┐              │
│  │  admission/ (Validação)                   │              │
│  │  ├─ admission_controller.go               │              │
│  │  ├─ capacity_calculator.go                │              │
│  │  └─ resource_validator.go                 │              │
│  │                                            │              │
│  │  Output: AdmissionResult                  │              │
│  │    - Admitted: true/false                 │              │
│  │    - GridCapacity                         │              │
│  │    - Recommendation (se rejected)         │              │
│  └──────────────────────────────────────────┘              │
│         ↓ (if admitted)                                      │
│  ┌──────────────────────────────────────────┐              │
│  │  scheduler/ (Placement Decision)          │              │
│  │  ├─ scheduler.go                          │              │
│  │  ├─ node_filter.go                        │              │
│  │  ├─ node_scorer.go                        │              │
│  │  └─ strategy_spread.go (ou binpack/opt)  │              │
│  │                                            │              │
│  │  Output: []PlacementDecision              │              │
│  │    - NodeID, ReplicaIndex, Score          │              │
│  └──────────────────────────────────────────┘              │
│         ↓                                                    │
│  ┌──────────────────────────────────────────┐              │
│  │  deploy/ (Execution)                      │              │
│  │  ├─ deployer.go                           │              │
│  │  ├─ executor.go                           │              │
│  │  ├─ docker_client.go                      │              │
│  │  └─ rollback.go (se falhar)               │              │
│  │                                            │              │
│  │  Output: DeploymentResult                 │              │
│  │    - Success: true/false                  │              │
│  │    - ReplicasRunning: 3/3                 │              │
│  └──────────────────────────────────────────┘              │
│         ↓                                                    │
│  ┌──────────────────────────────────────────┐              │
│  │  monitoring/ (Observability)              │              │
│  │  ├─ monitoring.go                         │              │
│  │  ├─ logs.go                               │              │
│  │  └─ metrics.go                            │              │
│  │                                            │              │
│  │  Background: Collect logs/metrics         │              │
│  └──────────────────────────────────────────┘              │
│                                                              │
│  ┌──────────────────────────────────────────┐              │
│  │  queue/ (se admission rejected)           │              │
│  │  ├─ queue_manager.go                      │              │
│  │  └─ queue_processor.go (periodic)         │              │
│  │                                            │              │
│  │  Background: Retry quando resources free  │              │
│  └──────────────────────────────────────────┘              │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 7.3 Subcomponente: Admission

**Responsabilidade**: Validar se workload pode ser aceito pela Grid.

**Arquivos**:

#### `admission/admission_controller.go` (< 400 linhas)
```go
package admission

// Funções principais:
// - Validate(req *WorkloadRequest) (*AdmissionResult, error)
// - ValidateGridWide(req) - valida toda Grid
// - ValidateNodeSpecific(req, nodeID) - valida Node específico

// Validações implementadas:
// 1. Nodes suficientes (healthy_nodes >= replicas)
// 2. CPU disponível (total_cpu_needed <= available_cpu)
// 3. RAM disponível (total_ram_needed <= available_ram)
// 4. Limite de sobrecarga (utilização projetada <= 90%)
// 5. Constraints (tags, disk-type, min-ram)
```

#### `admission/capacity_calculator.go` (< 300 linhas)
```go
package admission

// Funções principais:
// - GetGridCapacity() (*CapacityInfo, error)
// - CalculateNodeCapacity(node) (*NodeCapacityInfo, error)
// - ProjectUtilization(current, needed) (float64, error)

// Cálculos:
// - Total CPU/RAM da Grid
// - CPU/RAM disponível (total - usado)
// - Utilização atual (%)
// - Utilização projetada após deployment
```

#### `admission/constraint_validator.go` (< 300 linhas)
```go
package admission

// Funções principais:
// - ValidateConstraints(nodes, constraints) (*ConstraintResult, error)
// - ValidateTag(node, tag) bool
// - ValidateDiskType(node, diskType) bool
// - ValidateMinRAM(node, minRAM) bool
// - ValidateMinCPU(node, minCPU) bool

// Constraints suportados:
// - tag=<value> (ex: gpu, ssd, compute)
// - disk-type=<ssd|hdd|nvme>
// - min-ram=<size> (ex: 16G, 32G)
// - min-cpu=<cores> (ex: 4, 8)
// - max-workloads=<n> (limite de workloads por Node)
```

#### `admission/resource_validator.go` (< 300 linhas)
```go
package admission

// Funções principais:
// - ValidateCPU(needed, available) error
// - ValidateRAM(needed, available) error
// - ValidateDisk(needed, available) error
// - ValidateNetwork(bandwidth) error

// Limites configuráveis:
const (
    MaxCPUUtilization  = 90.0  // %
    MaxRAMUtilization  = 90.0  // %
    MaxDiskUtilization = 85.0  // %
    MinReservedCPU     = 0.5   // cores (reserva do sistema)
    MinReservedRAM     = 1.0   // GB (reserva do sistema)
)
```

### 7.4 Subcomponente: Scheduler

**Responsabilidade**: Decidir em quais Nodes deployar cada réplica.

**Arquivos**:

#### `scheduler/scheduler.go` (< 400 linhas)
```go
package scheduler

// Interface principal
type Scheduler interface {
    Schedule(req *WorkloadRequest) ([]*PlacementDecision, error)
}

// Funções principais:
// - Schedule(req) - orquestra todo processo de scheduling
// - SelectStrategy(strategyName) Strategy
// - ApplyConstraints(nodes, constraints) []*Node

// Decisões de placement:
type PlacementDecision struct {
    NodeID       string
    ReplicaIndex int
    Score        float64
    Reason       string
    Resources    *AllocatedResources
}
```

#### `scheduler/node_filter.go` (< 300 linhas)
```go
package scheduler

// Funções de filtragem:
// - FilterHealthyNodes(nodes) []*Node
// - FilterByConstraints(nodes, constraints) []*Node
// - FilterByCapacity(nodes, cpu, ram) []*Node
// - FilterByTags(nodes, tags) []*Node
// - FilterByDiskType(nodes, diskType) []*Node

// Pipeline de filtragem:
// AllNodes → Healthy → Constraints → Capacity → EligibleNodes
```

#### `scheduler/node_scorer.go` (< 300 linhas)
```go
package scheduler

// Funções de scoring:
// - CalculateSpreadScore(node) float64
// - CalculateBinpackScore(node) float64
// - CalculateOptimizedScore(node, cpu, ram) float64
// - CalculateLoadBalanceScore(node) float64

// Fatores considerados:
// - Número de workloads atual
// - Utilização de CPU/RAM
// - Balanceamento CPU vs RAM
// - Sweet spot (40-60% utilização)
```

#### `scheduler/strategy_spread.go` (< 300 linhas)
```go
package scheduler

// Estratégia: Distribuição Uniforme
type SpreadStrategy struct {
    inventoryMgr *InventoryManager
}

func (s *SpreadStrategy) Schedule(nodes []*Node, req *WorkloadRequest) []*PlacementDecision {
    // 1. Ordenar Nodes por MENOR número de workloads
    sort.Slice(nodes, func(i, j int) bool {
        return len(nodes[i].Workloads) < len(nodes[j].Workloads)
    })
    
    // 2. Round-robin allocation
    decisions := make([]*PlacementDecision, req.Replicas)
    for i := 0; i < req.Replicas; i++ {
        nodeIndex := i % len(nodes)
        decisions[i] = createDecision(nodes[nodeIndex], i+1, "spread")
    }
    
    return decisions
}

// Casos de uso:
// - Aplicações web stateless
// - Microservices
// - APIs REST
// - Alta disponibilidade
```

#### `scheduler/strategy_binpack.go` (< 300 linhas)
```go
package scheduler

// Estratégia: Preenchimento Denso
type BinpackStrategy struct {
    inventoryMgr *InventoryManager
}

func (s *BinpackStrategy) Schedule(nodes []*Node, req *WorkloadRequest) []*PlacementDecision {
    // 1. Ordenar Nodes por MAIOR utilização
    sort.Slice(nodes, func(i, j int) bool {
        utilI := nodes[i].Metrics.CPUPercent + nodes[i].Metrics.RAMUtilization
        utilJ := nodes[j].Metrics.CPUPercent + nodes[j].Metrics.RAMUtilization
        return utilI > utilJ
    })
    
    // 2. Preencher Nodes até capacidade máxima
    decisions := make([]*PlacementDecision, req.Replicas)
    nodeIndex := 0
    
    for i := 0; i < req.Replicas; i++ {
        // Procurar próximo Node com capacidade
        for !nodeHasCapacity(nodes[nodeIndex], req) {
            nodeIndex++
        }
        decisions[i] = createDecision(nodes[nodeIndex], i+1, "binpack")
    }
    
    return decisions
}

// Casos de uso:
// - Batch jobs
// - Jobs finitos
// - Economia de energia
// - Deixar Nodes livres para outras tarefas
```

#### `scheduler/strategy_optimized.go` (< 400 linhas)
```go
package scheduler

// Estratégia: Otimização de Recursos
type OptimizedStrategy struct {
    inventoryMgr *InventoryManager
}

func (s *OptimizedStrategy) Schedule(nodes []*Node, req *WorkloadRequest) []*PlacementDecision {
    decisions := make([]*PlacementDecision, req.Replicas)
    
    for i := 0; i < req.Replicas; i++ {
        // Calcular score para CADA Node
        scores := make([]NodeScore, len(nodes))
        
        for j, node := range nodes {
            scores[j] = NodeScore{
                Node:  node,
                Score: s.calculateOptimizedScore(node, req),
            }
        }
        
        // Escolher melhor score
        sort.Slice(scores, func(i, j int) bool {
            return scores[i].Score > scores[j].Score
        })
        
        bestNode := scores[0].Node
        decisions[i] = createDecision(bestNode, i+1, "resource-optimized")
        
        // Simular alocação para próxima iteração
        updateNodeUtilization(bestNode, req)
    }
    
    return decisions
}

func (s *OptimizedStrategy) calculateOptimizedScore(node *Node, req *WorkloadRequest) float64 {
    // Score baseado em:
    // 1. Balanceamento CPU/RAM (peso 50%)
    // 2. Utilização atual (peso 30%)
    // 3. Sweet spot 40-60% (peso 20%)
    
    cpuUtil := projectCPUUtilization(node, req.CPUPerReplica)
    ramUtil := projectRAMUtilization(node, req.RAMPerReplica)
    
    balanceDiff := abs(cpuUtil - ramUtil)
    avgUtil := (cpuUtil + ramUtil) / 2
    
    score := 100.0
    score -= balanceDiff * 0.5  // Penalizar desbalanceamento
    score -= cpuUtil * 0.3       // Preferir Nodes menos utilizados
    score -= ramUtil * 0.2
    
    // Bônus sweet spot
    if avgUtil >= 40 && avgUtil <= 60 {
        score += 20
    }
    
    return score
}

// Casos de uso:
// - Workloads balanceados (CPU e RAM)
// - Aplicações de longa duração
// - Otimização de recursos gerais
```

### 7.5 Subcomponente: Queue

**Responsabilidade**: Gerenciar workloads quando Grid está cheia.

**Arquivos**:

#### `queue/queue_manager.go` (< 300 linhas)
```go
package queue

// Funções principais:
// - Enqueue(req *WorkloadRequest) (*QueuedWorkload, error)
// - Dequeue(id string) error
// - List() []*QueuedWorkload
// - Get(id string) (*QueuedWorkload, error)
// - Cancel(id string) error

type QueuedWorkload struct {
    ID            string
    Request       *WorkloadRequest
    QueuedAt      time.Time
    Priority      int
    EstimatedWait time.Duration
    Status        string // queued, deploying, failed
}

// Persistência:
// - Salvar em ~/.syntropy/queue/
// - Formato: YAML
// - 1 arquivo por workload queued
```

#### `queue/queue_processor.go` (< 300 linhas)
```go
package queue

// Processador periódico (background job)
type QueueProcessor struct {
    queueMgr       *QueueManager
    admissionCtrl  *AdmissionController
    deployer       *Deployer
    processInterval time.Duration
}

// Funções principais:
// - Start() - inicia processamento periódico
// - Stop() - para processamento
// - ProcessOnce() error - processa fila uma vez

// Lógica:
// A cada 1 minuto:
//   1. Listar workloads na fila
//   2. Para cada workload:
//      - Tentar admission novamente
//      - Se aprovado: deployar automaticamente
//      - Se rejeitado: atualizar estimated_wait
//   3. Ordenar fila por prioridade

// Background job usando goroutine + ticker
```

#### `queue/wait_estimator.go` (< 200 linhas)
```go
package queue

// Funções principais:
// - EstimateWaitTime(req *WorkloadRequest) time.Duration
// - CalculateAverageJobDuration() time.Duration
// - PredictResourceAvailability() time.Time

// Estimativa baseada em:
// 1. Workloads atuais (quando vão terminar?)
// 2. Duração média histórica
// 3. Recursos necessários vs. utilizados
// 4. Posição na fila

// Algoritmo simplificado:
// estimated_wait = queue_position × avg_job_duration
```

#### `queue/priority_manager.go` (< 200 linhas)
```go
package queue

// Gerenciamento de prioridades
const (
    PriorityLow    = 0
    PriorityNormal = 5
    PriorityHigh   = 10
    PriorityUrgent = 15
)

// Funções principais:
// - SetPriority(workloadID, priority) error
// - ReorderQueue() - reordena por prioridade
// - GetNextWorkload() *QueuedWorkload

// Regras de prioridade (MVP):
// - FIFO por padrão (PriorityNormal)
// - Usuário pode aumentar prioridade
// - Workloads pequenos têm bônus
```

### 7.6 Subcomponente: Deploy

**Responsabilidade**: Executar deployment nos Nodes escolhidos.

**Arquivos**:

#### `deploy/deployer.go` (< 400 linhas)
```go
package deploy

// Orquestrador de deployment
type Deployer struct {
    admissionCtrl *AdmissionController
    scheduler     *Scheduler
    executor      *Executor
    rollbackMgr   *RollbackManager
    queueMgr      *QueueManager
}

// Funções principais:
// - Deploy(req *WorkloadRequest) (*DeploymentResult, error)
// - DeployGridWide(req) - deploy em toda Grid
// - DeployToNodes(req, nodeIDs) - deploy em Nodes específicos
// - Rollback(deploymentID) error

// Fluxo completo:
func (d *Deployer) Deploy(req *WorkloadRequest) (*DeploymentResult, error) {
    // 1. Admission Control
    admission, err := d.admissionCtrl.Validate(req)
    if !admission.Admitted {
        // Oferecer fila ou rejeitar
        if req.QueueIfFull {
            return d.queueMgr.Enqueue(req)
        }
        return nil, fmt.Errorf("admission denied: %s", admission.Reason)
    }
    
    // 2. Scheduling
    placements, err := d.scheduler.Schedule(req)
    
    // 3. Execution (com rollback em caso de falha)
    results := make([]*ReplicaResult, len(placements))
    for i, placement := range placements {
        result, err := d.executor.DeployReplica(placement, req)
        if err != nil {
            // ROLLBACK: desfazer deployments anteriores
            d.rollbackMgr.Rollback(results[:i])
            return nil, fmt.Errorf("deployment failed: %w", err)
        }
        results[i] = result
    }
    
    // 4. Salvar metadados
    return d.saveDeployment(req, placements, results)
}
```

#### `deploy/executor.go` (< 400 linhas)
```go
package deploy

// Executor de deployment
type Executor struct {
    sshClient    *SSHClient
    dockerClient *DockerClient
}

// Funções principais:
// - DeployReplica(placement, req) (*ReplicaResult, error)
// - ValidateImage(image) error
// - PullImage(node, image) error
// - CreateContainer(node, config) (containerID, error)
// - StartContainer(node, containerID) error
// - VerifyRunning(node, containerID) error

// Fluxo de deployment de 1 réplica:
func (e *Executor) DeployReplica(placement *PlacementDecision, req *WorkloadRequest) error {
    node := getNode(placement.NodeID)
    
    // 1. Validar imagem
    // 2. Pull image (se necessário)
    // 3. Criar container com configurações
    // 4. Iniciar container
    // 5. Verificar se está running
    // 6. Salvar em inventory
    
    return nil
}
```

#### `deploy/executor_windows.go` (< 300 linhas)
```go
//go:build windows

package deploy

// Implementação específica Windows
type WindowsExecutor struct {
    baseExecutor *Executor
}

// Diferenças Windows:
// - Paths: C:\syntropy\ vs /opt/syntropy/
// - SSH client: usar golang.org/x/crypto/ssh
// - Comandos: ajustar para Windows PowerShell se necessário
```

#### `deploy/executor_linux.go` (< 300 linhas)
```go
//go:build linux

package deploy

// Implementação específica Linux
type LinuxExecutor struct {
    baseExecutor *Executor
}

// Otimizações Linux:
// - SSH nativo
// - Docker socket direto
```

#### `deploy/rollback.go` (< 300 linhas)
```go
package deploy

// Gerenciador de rollback
type RollbackManager struct {
    executor *Executor
}

// Funções principais:
// - Rollback(deployments) error
// - RollbackSingle(deployment) error
// - SaveRollbackState(deployment) error
// - RestorePreviousVersion(workloadID) error

// Lógica de rollback:
func (rm *RollbackManager) Rollback(results []*ReplicaResult) error {
    fmt.Printf("🔄 Rolling back %d deployments...\n", len(results))
    
    for i := len(results) - 1; i >= 0; i-- {
        result := results[i]
        
        // Parar e remover container
        cmd := fmt.Sprintf("docker stop %s && docker rm %s", 
            result.ContainerID, result.ContainerID)
        
        if err := executeSSH(result.NodeIP, cmd); err != nil {
            fmt.Printf("⚠️  Warning: failed to rollback on %s: %v\n", 
                result.NodeID, err)
        } else {
            fmt.Printf("   ✅ Rolled back on %s\n", result.NodeID)
        }
    }
    
    return nil
}

// Casos de rollback:
// 1. Falha em qualquer réplica → rollback todas
// 2. Timeout no deployment → rollback
// 3. Validação pós-deployment falha → rollback
// 4. Comando explícito do usuário → rollback
```

#### `deploy/docker_client.go` (< 400 linhas)
```go
package deploy

// Cliente Docker (via SSH ou API)
type DockerClient struct {
    sshClient *SSHClient
}

// Funções principais:
// - PullImage(nodeIP, image) error
// - CreateContainer(nodeIP, config) (string, error)
// - StartContainer(nodeIP, containerID) error
// - StopContainer(nodeIP, containerID) error
// - RemoveContainer(nodeIP, containerID) error
// - GetContainerStatus(nodeIP, containerID) (string, error)
// - GetContainerLogs(nodeIP, containerID) (string, error)

// Comandos Docker executados via SSH:
// - docker pull <image>
// - docker create --name <name> --cpu <cpu> --memory <mem> <image>
// - docker start <container>
// - docker ps --filter id=<container>
// - docker logs <container>
```

### 7.7 Subcomponente: Lifecycle

**Responsabilidade**: Gerenciar lifecycle de workloads (start, stop, restart, scale).

**Arquivos**:

#### `lifecycle/lifecycle.go` (< 300 linhas)
```go
package lifecycle

// Manager de lifecycle
type LifecycleManager struct {
    inventoryMgr *InventoryManager
    executor     *Executor
}

// Funções principais:
// - Start(workloadID) error
// - Stop(workloadID) error
// - Restart(workloadID) error
// - Scale(workloadID, replicas) error
// - Remove(workloadID) error

// Orquestra subcomponentes
```

#### `lifecycle/start.go` (< 200 linhas)
```go
package lifecycle

// Start workload parado
func (lm *LifecycleManager) Start(workloadID string) error {
    // 1. Carregar workload de ~/.syntropy/workloads/
    // 2. Verificar status atual (se já running, retornar)
    // 3. Para cada réplica:
    //    - docker start <container>
    // 4. Atualizar status
    // 5. Salvar
}
```

#### `lifecycle/stop.go` (< 200 linhas)
```go
package lifecycle

// Stop workload
func (lm *LifecycleManager) Stop(workloadID string) error {
    // 1. Carregar workload
    // 2. Para cada réplica:
    //    - docker stop <container>
    // 3. Atualizar status: stopped
    // 4. Salvar
}
```

#### `lifecycle/restart.go` (< 200 linhas)
```go
package lifecycle

// Restart workload
func (lm *LifecycleManager) Restart(workloadID string) error {
    // 1. Stop
    // 2. Start
    // 3. Verificar saúde
}
```

#### `lifecycle/scale.go` (< 300 linhas)
```go
package lifecycle

// Scale workload (up ou down)
func (lm *LifecycleManager) Scale(workloadID string, newReplicas int) error {
    // 1. Carregar workload atual
    currentReplicas := len(workload.Replicas)
    
    if newReplicas > currentReplicas {
        // SCALE UP
        delta := newReplicas - currentReplicas
        
        // 1. Criar WorkloadRequest para novas réplicas
        // 2. Passar por Admission Control
        // 3. Scheduler decide placement
        // 4. Deploy novas réplicas
        
    } else if newReplicas < currentReplicas {
        // SCALE DOWN
        delta := currentReplicas - newReplicas
        
        // 1. Selecionar réplicas a remover (últimas criadas)
        // 2. Parar containers
        // 3. Remover containers
        // 4. Atualizar inventory
    }
    
    return nil
}

// Integração com Admission:
// - Scale UP passa por validação (como novo deployment)
// - Scale DOWN libera recursos (atualiza capacity)
```

### 7.8 Subcomponente: Monitoring

**Responsabilidade**: Logs e métricas dos workloads.

**Arquivos**:

#### `monitoring/monitoring.go` (< 300 linhas)
```go
package monitoring

// Monitor principal
type WorkloadMonitor struct {
    inventoryMgr *InventoryManager
    sshClient    *SSHClient
}

// Funções principais:
// - GetWorkloadStatus(workloadID) (*WorkloadStatus, error)
// - GetWorkloadLogs(workloadID, options) ([]LogEntry, error)
// - GetWorkloadMetrics(workloadID) (*WorkloadMetrics, error)
// - StreamLogs(workloadID, follow bool) (chan LogEntry, error)
```

#### `monitoring/logs.go` (< 400 linhas)
```go
package monitoring

// Agregação de logs
type LogAggregator struct {
    sshClient *SSHClient
}

// Funções principais:
// - GetLogs(nodeIP, containerID, options) ([]LogEntry, error)
// - StreamLogs(nodeIP, containerID) (chan LogEntry, error)
// - FilterLogs(logs, level, since, until) []LogEntry
// - AggregateLogs(workloadID) []LogEntry

type LogEntry struct {
    Timestamp   time.Time
    NodeID      string
    ContainerID string
    Level       string
    Message     string
}

// Implementação:
// - Executar: docker logs <container> via SSH
// - Parse de logs
// - Agregação de múltiplos Nodes
// - Streaming em tempo real (follow)
```

#### `monitoring/metrics.go` (< 400 linhas)
```go
package monitoring

// Coleta de métricas
type MetricsCollector struct {
    sshClient *SSHClient
}

// Funções principais:
// - CollectWorkloadMetrics(workloadID) (*WorkloadMetrics, error)
// - CollectContainerMetrics(nodeIP, containerID) (*ContainerMetrics, error)
// - AggregateMetrics(workloadID) (*AggregatedMetrics, error)

type WorkloadMetrics struct {
    WorkloadID    string
    Replicas      int
    TotalCPU      float64
    TotalRAM      float64
    NetworkRx     float64
    NetworkTx     float64
    PerReplica    []*ContainerMetrics
}

type ContainerMetrics struct {
    ContainerID string
    NodeID      string
    CPUPercent  float64
    RAMUsedMB   float64
    NetworkRx   float64
    NetworkTx   float64
    BlockIO     float64
}

// Coleta via:
// - docker stats <container> --no-stream (via SSH)
// - Parse de output
// - Agregação de todas as réplicas
```

### 7.9 Integração do Componente Workload

**Arquivo Principal**: `workload/workload.go` (< 500 linhas)

```go
package workload

import (
    "github.com/syntropy/workload/admission"
    "github.com/syntropy/workload/scheduler"
    "github.com/syntropy/workload/queue"
    "github.com/syntropy/workload/deploy"
    "github.com/syntropy/workload/lifecycle"
    "github.com/syntropy/workload/monitoring"
)

// WorkloadComponent orquestra todos os subcomponentes
type WorkloadComponent struct {
    admissionCtrl *admission.AdmissionController
    scheduler     *scheduler.Scheduler
    queueMgr      *queue.QueueManager
    queueProc     *queue.QueueProcessor
    deployer      *deploy.Deployer
    lifecycleMgr  *lifecycle.LifecycleManager
    monitor       *monitoring.WorkloadMonitor
    inventoryMgr  *InventoryManager
}

// NewWorkloadComponent cria componente completo
func NewWorkloadComponent(inventoryMgr *InventoryManager, sshClient *SSHClient) *WorkloadComponent {
    // Inicializar subcomponentes
    admissionCtrl := admission.NewAdmissionController(inventoryMgr)
    sched := scheduler.NewScheduler(inventoryMgr)
    queueMgr := queue.NewQueueManager()
    executor := deploy.NewExecutor(sshClient)
    deployer := deploy.NewDeployer(admissionCtrl, sched, executor, queueMgr)
    lifecycleMgr := lifecycle.NewLifecycleManager(inventoryMgr, executor)
    monitor := monitoring.NewWorkloadMonitor(inventoryMgr, sshClient)
    queueProc := queue.NewQueueProcessor(queueMgr, admissionCtrl, deployer, 1*time.Minute)
    
    return &WorkloadComponent{
        admissionCtrl: admissionCtrl,
        scheduler:     sched,
        queueMgr:      queueMgr,
        queueProc:     queueProc,
        deployer:      deployer,
        lifecycleMgr:  lifecycleMgr,
        monitor:       monitor,
        inventoryMgr:  inventoryMgr,
    }
}

// Deploy workload (interface pública)
func (wc *WorkloadComponent) Deploy(req *WorkloadRequest) (*DeploymentResult, error) {
    return wc.deployer.Deploy(req)
}

// List workloads
func (wc *WorkloadComponent) List() ([]*Workload, error) {
    // Listar de ~/.syntropy/workloads/
}

// Get workload status
func (wc *WorkloadComponent) Status(workloadID string) (*WorkloadStatus, error) {
    return wc.monitor.GetWorkloadStatus(workloadID)
}

// Lifecycle operations
func (wc *WorkloadComponent) Start(workloadID string) error {
    return wc.lifecycleMgr.Start(workloadID)
}

func (wc *WorkloadComponent) Stop(workloadID string) error {
    return wc.lifecycleMgr.Stop(workloadID)
}

func (wc *WorkloadComponent) Scale(workloadID string, replicas int) error {
    return wc.lifecycleMgr.Scale(workloadID, replicas)
}

// Monitoring
func (wc *WorkloadComponent) Logs(workloadID string, follow bool) (chan LogEntry, error) {
    return wc.monitor.StreamLogs(workloadID, follow)
}

// Queue management
func (wc *WorkloadComponent) QueueList() ([]*QueuedWorkload, error) {
    return wc.queueMgr.List()
}

// Start background jobs
func (wc *WorkloadComponent) Start() error {
    // Iniciar queue processor
    return wc.queueProc.Start()
}
```

### 7.10 Estrutura Completa de Componentes

```
manager/interfaces/cli/
│
├── setup/                     # Componente: Setup ✅
│   ├── README.md
│   ├── ARCHITECTURE.md
│   └── src/
│       ├── setup.go           # Orquestrador (< 500 linhas)
│       ├── configurator.go    # Subcomponente: Configuration
│       ├── key_manager.go     # Subcomponente: Keys
│       └── token_manager.go   # Subcomponente: Token (🔧 a criar)
│
├── node/                      # Componente: Node 🚧
│   ├── README.md
│   ├── ARCHITECTURE.md
│   ├── node.go                # Orquestrador (< 500 linhas)
│   ├── create/                # Subcomponente: Creation
│   ├── registration/          # Subcomponente: Registration
│   └── inventory/             # Subcomponente: Inventory
│
├── workload/                  # Componente: Workload 🚧 (ATUALIZADO)
│   ├── README.md              # 📖 Documentação do componente
│   ├── ARCHITECTURE.md        # 🏗️  Arquitetura e fluxos
│   ├── workload.go            # Orquestrador principal (< 500 linhas)
│   │
│   ├── admission/             # Subcomponente: Admission Control
│   │   ├── README.md          # Doc: O que é Admission Control
│   │   ├── admission_controller.go  (< 400 linhas)
│   │   ├── capacity_calculator.go   (< 300 linhas)
│   │   ├── constraint_validator.go  (< 300 linhas)
│   │   ├── resource_validator.go    (< 300 linhas)
│   │   └── tests/
│   │       ├── admission_test.go
│   │       └── capacity_test.go
│   │
│   ├── scheduler/             # Subcomponente: Scheduler
│   │   ├── README.md          # Doc: Estratégias de scheduling
│   │   ├── scheduler.go       (< 400 linhas)
│   │   ├── node_filter.go     (< 300 linhas)
│   │   ├── node_scorer.go     (< 300 linhas)
│   │   ├── strategy_spread.go (< 300 linhas)
│   │   ├── strategy_binpack.go (< 300 linhas)
│   │   ├── strategy_optimized.go (< 400 linhas)
│   │   └── tests/
│   │       ├── scheduler_test.go
│   │       ├── spread_test.go
│   │       ├── binpack_test.go
│   │       └── optimized_test.go
│   │
│   ├── queue/                 # Subcomponente: Queue Management
│   │   ├── README.md          # Doc: Sistema de filas
│   │   ├── queue_manager.go   (< 300 linhas)
│   │   ├── queue_processor.go (< 300 linhas)
│   │   ├── wait_estimator.go  (< 200 linhas)
│   │   ├── priority_manager.go (< 200 linhas)
│   │   └── tests/
│   │       └── queue_test.go
│   │
│   ├── deploy/                # Subcomponente: Deployment Execution
│   │   ├── README.md          # Doc: Execução de deployments
│   │   ├── deployer.go        (< 400 linhas)
│   │   ├── executor.go        (< 400 linhas)
│   │   ├── executor_windows.go (< 300 linhas)
│   │   ├── executor_linux.go  (< 300 linhas)
│   │   ├── rollback.go        (< 300 linhas)
│   │   ├── docker_client.go   (< 400 linhas)
│   │   └── tests/
│   │       ├── deployer_test.go
│   │       └── rollback_test.go
│   │
│   ├── lifecycle/             # Subcomponente: Lifecycle
│   │   ├── README.md          # Doc: Gerenciamento de lifecycle
│   │   ├── lifecycle.go       (< 300 linhas)
│   │   ├── start.go           (< 200 linhas)
│   │   ├── stop.go            (< 200 linhas)
│   │   ├── restart.go         (< 200 linhas)
│   │   ├── scale.go           (< 300 linhas)
│   │   └── tests/
│   │       └── lifecycle_test.go
│   │
│   └── monitoring/            # Subcomponente: Monitoring
│       ├── README.md          # Doc: Observabilidade
│       ├── monitoring.go      (< 300 linhas)
│       ├── logs.go            (< 400 linhas)
│       ├── metrics.go         (< 400 linhas)
│       └── tests/
│           └── monitoring_test.go
│
└── management/                # Componente: Management 🚧
    ├── README.md
    ├── ARCHITECTURE.md
    ├── management.go          # Orquestrador (< 500 linhas)
    ├── discovery/             # Subcomponente: Discovery
    ├── health/                # Subcomponente: Health
    └── sync/                  # Subcomponente: Sync
```

### 7.11 Padrão de Componentes (Best Practices)

Cada componente segue esta estrutura:

```
<component>/
├── README.md              # O QUÊ é o componente
├── ARCHITECTURE.md        # COMO funciona internamente
├── <component>.go         # Orquestrador principal (< 500 linhas)
├── <subcomponent>/        # Subcomponentes
│   ├── README.md          # Documentação do subcomponente
│   ├── <sub>.go           # Interface comum (< 300 linhas)
│   ├── <sub>_windows.go   # Implementação Windows (< 500 linhas)
│   ├── <sub>_linux.go     # Implementação Linux (< 500 linhas)
│   └── tests/             # Testes unitários
│       └── <sub>_test.go
└── tests/                 # Testes de integração
    └── <component>_integration_test.go
```

### 6.3 Implementação Multi-Plataforma

Usar build tags do Go:

```go
//go:build windows

package create

// Implementação específica Windows
func createUSB() error {
    // Use diskpart, PowerShell, etc.
}
```

```go
//go:build linux

package create

// Implementação específica Linux
func createUSB() error {
    // Use dd, mkfs, etc.
}
```

---

## 7. ROADMAP DE IMPLEMENTAÇÃO (REVISADO)

### 7.1 Semana 1: Setup + Foundations
```
[✅ DONE] Setup Component (80%)
[🚧 TODO] Completar Setup:
  - Gerar Grid Token
  - Integrar KeyManager de infrastructure/
  - Melhorar geração de chaves SSH

[🚧 TODO] Estrutura de Componentes:
  - Criar diretórios node/, workload/, management/
  - Criar README.md e ARCHITECTURE.md de cada componente
  - Definir interfaces principais
```

### 7.2 Semana 2: Node Creation (Part 1)
```
[🚧 TODO] USB Detection:
  - Implementar usb_detector.go (interface)
  - Implementar usb_detector_windows.go
  - Implementar usb_detector_linux.go
  - Testar detecção em ambas plataformas

[🚧 TODO] ISO Download:
  - Implementar iso_downloader.go
  - Download Ubuntu 24.04 Server
  - Caching em ~/.syntropy/cache/
  - Verificação de checksum
```

### 7.3 Semana 3: Node Creation (Part 2)
```
[🚧 TODO] Cloud-Init Generation:
  - Implementar cloud_init_generator.go
  - Integrar com infrastructure/template_manager.go
  - Gerar user-data, meta-data, network-config
  - Injetar Grid Token e SSH keys

[🚧 TODO] USB Writing:
  - Implementar usb_writer.go
  - Injetar cloud-init no ISO
  - Gravar USB bootável
  - Implementação Windows + Linux
```

### 7.4 Semana 4: Registration Protocol
```
[🚧 TODO] Registration:
  - Implementar registration.go
  - Listener para Node announcements
  - Validação de Grid Token
  - Handshake automático

[🚧 TODO] Inventory Management:
  - Implementar inventory.go
  - CRUD de Nodes
  - Atualização com Hardware Manifest
  - Persistência em YAML

[📋 TESTE] Provisionar primeiro Node físico completo
```

### 7.5 Semana 5: Workload Deployment + Orquestração

**ATUALIZADO**: Integração completa de Admission Control + Scheduler + Queue System

```
[🚧 TODO] Subcomponente: Admission
  Dia 1-2: Implementação
    - workload/admission/admission_controller.go (< 400 linhas)
    - workload/admission/capacity_calculator.go (< 300 linhas)
    - workload/admission/constraint_validator.go (< 300 linhas)
    - workload/admission/resource_validator.go (< 300 linhas)
  
  Validações a implementar:
    - Nodes suficientes
    - CPU disponível
    - RAM disponível
    - Limite de sobrecarga (90%)
    - Constraints (tags, disk-type, min-ram)
  
  Testes:
    - Workload válido (aceito)
    - Workload muito grande (rejeitado)
    - Grid cheia (rejeitado com sugestões)
    - Constraints não atendidas (rejeitado)

[🚧 TODO] Subcomponente: Scheduler
  Dia 3-4: Implementação
    - workload/scheduler/scheduler.go (< 400 linhas)
    - workload/scheduler/node_filter.go (< 300 linhas)
    - workload/scheduler/node_scorer.go (< 300 linhas)
    - workload/scheduler/strategy_spread.go (< 300 linhas)
    - workload/scheduler/strategy_binpack.go (< 300 linhas)
    - workload/scheduler/strategy_optimized.go (< 400 linhas)
  
  Estratégias a implementar:
    - spread (distribuição uniforme) - PADRÃO
    - binpack (preenchimento denso)
    - resource-optimized (balanceamento CPU/RAM)
  
  Testes:
    - Spread com 3 réplicas em 6 Nodes
    - Binpack com Nodes parcialmente cheios
    - Resource-optimized com carga balanceada
    - Filtragem por constraints

[🚧 TODO] Subcomponente: Queue
  Dia 5: Implementação
    - workload/queue/queue_manager.go (< 300 linhas)
    - workload/queue/queue_processor.go (< 300 linhas)
    - workload/queue/wait_estimator.go (< 200 linhas)
    - workload/queue/priority_manager.go (< 200 linhas)
  
  Funcionalidades:
    - Enfileirar workloads rejeitados
    - Processar fila periodicamente (1 min)
    - Estimar tempo de espera
    - Deploy automático quando recursos liberarem
    - Gerenciamento de prioridades
  
  Testes:
    - Enfileirar workload
    - Processamento automático
    - Estimativa de tempo
    - Cancelamento de fila

[🚧 TODO] Subcomponente: Deploy
  Dia 6: Implementação
    - workload/deploy/deployer.go (< 400 linhas)
    - workload/deploy/executor.go (< 400 linhas)
    - workload/deploy/executor_windows.go (< 300 linhas)
    - workload/deploy/executor_linux.go (< 300 linhas)
    - workload/deploy/rollback.go (< 300 linhas)
    - workload/deploy/docker_client.go (< 400 linhas)
  
  Funcionalidades:
    - Orquestrar Admission + Scheduler
    - Executar deployment via SSH/Docker
    - Rollback automático em falhas
    - Multi-plataforma (Windows/Linux)
  
  Testes:
    - Deploy bem-sucedido
    - Deploy com falha (+ rollback)
    - Deploy em múltiplos Nodes
    - Verificação pós-deployment

[🚧 TODO] Orquestrador Principal
  Dia 7: Integração
    - workload/workload.go (< 500 linhas)
    - Integrar todos os subcomponentes
    - Interface pública do componente
    - Background jobs (queue processor)
  
  Testes de Integração:
    - Deploy grid-wide completo
    - Admission → Scheduler → Deploy
    - Rollback em falha
    - Queue processing automático

[📋 TESTE] Cenários End-to-End
  - Deploy Nginx (3 réplicas, spread)
  - Deploy PostgreSQL (1 réplica, Node específico)
  - Deploy ML training (2 réplicas, tag=gpu, binpack)
  - Deploy com Grid cheia (queue)
  - Deploy com auto-ajuste
```

### 7.6 Semana 6: Management & Finalization
```
[🚧 TODO] Management Component:
  - Implementar management/health/health.go
  - Implementar management/sync/sync.go
  - Comando: syntropy node list
  - Comando: syntropy node status

[🚧 TODO] Provisionar 6 Nodes completos

[📋 TESTE] Grid completa funcionando
```

---

## 8. CHECKLIST DE IMPLEMENTAÇÃO PARA LLM

### ✅ Pilar 1: Setup (Completar)
```
[ ] Implementar GenerateGridToken() em configurator.go
[ ] Implementar SaveGridToken()
[ ] Integrar infrastructure/key_manager.go
[ ] Testar: syntropy setup run --force
[ ] Validar: ~/.syntropy/grid-token.txt gerado
[ ] Validar: Chaves SSH geradas corretamente
```

### ✅ Pilar 2: Node Creation (Implementar)
```
[ ] Criar estrutura de diretórios node/
[ ] Criar README.md e ARCHITECTURE.md
[ ] Implementar usb_detector.go (interface)
[ ] Implementar usb_detector_windows.go
[ ] Implementar usb_detector_linux.go
[ ] Testar: Detecção de USBs em Windows
[ ] Testar: Detecção de USBs em Linux
[ ] Implementar iso_downloader.go
[ ] Testar: Download + cache de ISO
[ ] Implementar cloud_init_generator.go
[ ] Integrar infrastructure/template_manager.go
[ ] Testar: Geração de cloud-init
[ ] Implementar usb_writer_windows.go
[ ] Implementar usb_writer_linux.go
[ ] Testar: Criação de USB bootável
[ ] Implementar registration.go
[ ] Implementar token_manager.go
[ ] Implementar handshake.go
[ ] Testar: Registration Protocol completo
[ ] Implementar inventory.go
[ ] Implementar hardware_manifest.go
[ ] Testar: Provisionar Node físico real
```

### ✅ Pilar 3: Workload Deployment (ATUALIZADO COM ORQUESTRAÇÃO)

```
[ ] Criar estrutura de componente workload/
    [ ] README.md (documentação do componente)
    [ ] ARCHITECTURE.md (fluxos e decisões de design)
    [ ] workload.go (orquestrador principal)

[ ] Subcomponente: Admission Control
    [ ] admission/admission_controller.go
    [ ] admission/capacity_calculator.go
    [ ] admission/constraint_validator.go
    [ ] admission/resource_validator.go
    [ ] admission/tests/admission_test.go
    [ ] Testar: Validação de recursos
    [ ] Testar: Limite de sobrecarga (90%)
    [ ] Testar: Constraints (tags, disk-type)
    [ ] Testar: Rejeição com sugestões

[ ] Subcomponente: Scheduler
    [ ] scheduler/scheduler.go
    [ ] scheduler/node_filter.go
    [ ] scheduler/node_scorer.go
    [ ] scheduler/strategy_spread.go
    [ ] scheduler/strategy_binpack.go
    [ ] scheduler/strategy_optimized.go
    [ ] scheduler/tests/scheduler_test.go
    [ ] scheduler/tests/spread_test.go
    [ ] scheduler/tests/binpack_test.go
    [ ] scheduler/tests/optimized_test.go
    [ ] Testar: Estratégia spread (distribuição uniforme)
    [ ] Testar: Estratégia binpack (preenchimento denso)
    [ ] Testar: Estratégia resource-optimized
    [ ] Testar: Filtragem de Nodes por constraints

[ ] Subcomponente: Queue Management
    [ ] queue/queue_manager.go
    [ ] queue/queue_processor.go
    [ ] queue/wait_estimator.go
    [ ] queue/priority_manager.go
    [ ] queue/tests/queue_test.go
    [ ] Testar: Enfileirar workload
    [ ] Testar: Processamento periódico
    [ ] Testar: Estimativa de tempo
    [ ] Testar: Deploy automático da fila

[ ] Subcomponente: Deployment Execution
    [ ] deploy/deployer.go
    [ ] deploy/executor.go
    [ ] deploy/executor_windows.go
    [ ] deploy/executor_linux.go
    [ ] deploy/rollback.go
    [ ] deploy/docker_client.go
    [ ] deploy/tests/deployer_test.go
    [ ] deploy/tests/rollback_test.go
    [ ] Testar: Deploy via SSH
    [ ] Testar: Deploy com Docker API
    [ ] Testar: Rollback em falha
    [ ] Testar: Deploy em Windows
    [ ] Testar: Deploy em Linux

[ ] Subcomponente: Lifecycle
    [ ] lifecycle/lifecycle.go
    [ ] lifecycle/start.go
    [ ] lifecycle/stop.go
    [ ] lifecycle/restart.go
    [ ] lifecycle/scale.go
    [ ] lifecycle/tests/lifecycle_test.go
    [ ] Testar: Start workload
    [ ] Testar: Stop workload
    [ ] Testar: Restart workload
    [ ] Testar: Scale up (com Admission Control)
    [ ] Testar: Scale down

[ ] Subcomponente: Monitoring
    [ ] monitoring/monitoring.go
    [ ] monitoring/logs.go
    [ ] monitoring/metrics.go
    [ ] monitoring/tests/monitoring_test.go
    [ ] Testar: Agregação de logs
    [ ] Testar: Streaming de logs (follow)
    [ ] Testar: Coleta de métricas
    [ ] Testar: Métricas agregadas de múltiplas réplicas

[ ] Integração Completa
    [ ] Integrar Admission + Scheduler + Deploy
    [ ] Background job: Queue processor
    [ ] Commands CLI: deploy, list, logs, status, scale
    [ ] Testar: Deploy grid-wide (sem especificar Node)
    [ ] Testar: Deploy com constraints
    [ ] Testar: Deploy com Grid cheia (queue)
    [ ] Testar: Auto-ajuste de recursos

[ ] Testes End-to-End
    [ ] Deploy Nginx (3 réplicas, spread)
    [ ] Deploy PostgreSQL (1 réplica, Node específico)
    [ ] Deploy com Grid 95% cheia (rejeição)
    [ ] Deploy em fila (deployment automático)
    [ ] Scale workload (up e down)
    [ ] Rollback após falha
```

### ✅ Pilar 4: Management
```
[ ] Criar estrutura de diretórios management/
[ ] Implementar health/health.go
[ ] Implementar sync/sync.go
[ ] Implementar node list command
[ ] Implementar node status command
[ ] Testar: syntropy node list (6 nodes)
[ ] Testar: syntropy node status node-01
[ ] Testar: Sincronização automática
```

---

## 9. ANEXO: TEMPLATES CLOUD-INIT

### 9.1 Template user-data Atualizado

```yaml
#cloud-config
# Syntropy Cooperative Grid - Node Provisioning
# Generated by Command Station

# Basic system configuration
locale: pt_BR.UTF-8
timezone: America/Sao_Paulo
hostname: syntropy-${NODE_NAME}

# User configuration
users:
  - name: syntropy
    groups: [adm, sudo, docker]
    shell: /bin/bash
    sudo: ALL=(ALL) NOPASSWD:ALL
    lock_passwd: true
    ssh_authorized_keys:
      - ${SSH_PUBLIC_KEY}  # Command Station public key

# SSH configuration
ssh_pwauth: false
disable_root: true

# Packages to install
packages:
  - curl
  - wget
  - git
  - htop
  - vim
  - docker.io
  - docker-compose-plugin
  - fail2ban
  - ufw
  - jq
  - prometheus-node-exporter

# Commands to run after installation
runcmd:
  # Docker setup
  - systemctl enable docker
  - systemctl start docker
  - usermod -aG docker syntropy
  
  # Firewall setup
  - ufw --force enable
  - ufw default deny incoming
  - ufw default allow outgoing
  - ufw allow from ${COMMAND_STATION_IP} to any port 22
  - ufw allow 51000/tcp  # Registration Protocol
  - ufw allow 8080/tcp   # Agent API
  - ufw allow 9100/tcp   # Prometheus metrics
  
  # Fail2ban setup
  - systemctl enable fail2ban
  - systemctl start fail2ban
  
  # Create Syntropy directories
  - mkdir -p /opt/syntropy/{bin,config,logs,identity,metadata}
  - chown -R syntropy:syntropy /opt/syntropy
  
  # Install Agent (placeholder - will be in USB)
  - cp /media/usb/syntropy/agent /opt/syntropy/bin/syntropy-agent
  - chmod +x /opt/syntropy/bin/syntropy-agent
  
  # Create Agent configuration
  - |
    cat > /opt/syntropy/config/agent.yaml << EOF
    node:
      name: "${NODE_NAME}"
      grid_token: "${GRID_TOKEN}"
    
    command_station:
      ip: "${COMMAND_STATION_IP}"
      port: 51000
    
    registration:
      auto_register: true
      timeout: 300s
    
    hardware:
      auto_detect: true
      report_interval: 30s
    
    logging:
      level: "info"
      file: "/opt/syntropy/logs/agent.log"
    EOF
  
  # Hardware detection script
  - |
    cat > /opt/syntropy/bin/detect-hardware << 'EOF'
    #!/bin/bash
    # Auto-detect hardware and generate manifest
    
    CPU_CORES=$(nproc)
    CPU_MODEL=$(lscpu | grep "Model name" | sed 's/Model name:\\s*//')
    RAM_GB=$(free -g | awk '/^Mem:/{print $2}')
    DISK_GB=$(df -BG / | awk 'NR==2{print $2}' | sed 's/G//')
    DISK_TYPE=$(lsblk -d -o name,rota | awk '$2==0{print "nvme"; exit} $2==1{print "hdd"}')
    
    # Network interfaces
    INTERFACES=$(ip -j link show | jq -c '[.[] | select(.link_type=="ether") | {name:.ifname, mac:.address}]')
    
    # Generate manifest
    cat > /opt/syntropy/metadata/hardware-manifest.json << MANIFEST
    {
      "cpu_cores": $CPU_CORES,
      "cpu_model": "$CPU_MODEL",
      "ram_gb": $RAM_GB,
      "disk_gb": $DISK_GB,
      "disk_type": "$DISK_TYPE",
      "network_interfaces": $INTERFACES,
      "detected_at": "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
    }
    MANIFEST
    
    chown syntropy:syntropy /opt/syntropy/metadata/hardware-manifest.json
    EOF
  
  - chmod +x /opt/syntropy/bin/detect-hardware
  - /opt/syntropy/bin/detect-hardware
  
  # Create systemd service for Agent
  - |
    cat > /etc/systemd/system/syntropy-agent.service << 'EOF'
    [Unit]
    Description=Syntropy Cooperative Grid Agent
    After=network.target docker.service
    Wants=docker.service
    
    [Service]
    Type=simple
    User=syntropy
    Group=syntropy
    WorkingDirectory=/opt/syntropy
    ExecStart=/opt/syntropy/bin/syntropy-agent --config=/opt/syntropy/config/agent.yaml
    Restart=always
    RestartSec=10
    StandardOutput=journal
    StandardError=journal
    
    [Install]
    WantedBy=multi-user.target
    EOF
  
  # Enable and start Agent
  - systemctl daemon-reload
  - systemctl enable syntropy-agent
  - systemctl start syntropy-agent

# Final message
final_message: |
  ✅ Syntropy Node ${NODE_NAME} configured successfully!
  
  📊 Hardware Detection: /opt/syntropy/metadata/hardware-manifest.json
  🔐 Grid Token: ${GRID_TOKEN:0:8}...
  🌐 Command Station: ${COMMAND_STATION_IP}
  
  🚀 Agent starting... will auto-register in ~30 seconds
  📝 Logs: journalctl -u syntropy-agent -f
```

---

## 10. CONSIDERAÇÕES FINAIS

### 10.1 Segurança
- **Grid Token**: UUID único, armazenado com permissões 600
- **SSH Keys**: ED25519, geradas automaticamente
- **Firewall**: UFW habilitado, regras restritivas
- **Fail2ban**: Proteção contra brute-force
- **Registration**: Validação em 3 camadas (token + nome + chave)

### 10.2 Resiliência
- **Auto-detecção**: Hardware detectado automaticamente
- **Auto-registro**: Nodes se registram sem intervenção
- **Sincronização**: Polling a cada 30s mantém estado atualizado
- **Healthcheck**: Detecta Nodes offline automaticamente
- **DHCP**: Rede configurada automaticamente

### 10.3 Extensibilidade
- **Multi-plataforma**: Windows e Linux suportados
- **Componentes Independentes**: Cada componente é isolado
- **Templates Centralizados**: Fácil atualização em `infrastructure/`
- **Hardware Manifest**: Suporta novos tipos de hardware

---

---

## ⚠️ CORREÇÕES CRÍTICAS IDENTIFICADAS

**IMPORTANTE**: Durante análise crítica, foram identificadas inconsistências entre:
1. Este documento MVP
2. Templates existentes em `infrastructure/cloud-init/`
3. Código já implementado em `manager/interfaces/cli/setup/`

**DOCUMENTO DE CORREÇÕES**: Ver `docs/architecture/MVP-CORRECTIONS.md`

### Correções Críticas (Implementar ANTES de começar):

#### 1. ��� Segurança do Grid Token
- ❌ **Problema**: Grid Token em texto plano (`grid-token.txt`)
- ✅ **Solução**: Usar Keyring do sistema operacional
- 📁 **Arquivo**: `setup/src/token_manager.go` (criar)
- 📦 **Dependência**: `github.com/zalando/go-keyring`

#### 2. 📋 Templates Cloud-Init Desatualizados
- ❌ **Problema**: Templates atuais não batem com MVP
  - Falta `${GRID_TOKEN}`
  - Falta script `detect-hardware`
  - Falta script `register-node`
  - Agent tenta baixar do GitHub (não existe)
- ✅ **Solução**: Criar templates `-mvp.yaml` simplificados
- 📁 **Arquivos**: Ver `MVP-CORRECTIONS.md` seção 2

#### 3. 🤖 Agent Não Existe
- ❌ **Problema**: MVP assume Agent funcional
- ✅ **Solução**: Usar Agent Placeholder (script bash)
- 📁 **Arquivo**: USB deve conter `syntropy/agent` (script)
- 🎯 **Estratégia**: Implementação faseada
  - Fase 1: Placeholder (status reporting)
  - Fase 2: Registration
  - Fase 3: Deploy real
  - Fase 4: Agent completo (pós-MVP)

### Status de Implementação Atualizado

```
✅ IMPLEMENTADO (pode usar):
  - Setup Component (80%)
  - KeyManager (infrastructure/)
  - TemplateManager (infrastructure/)
  - Estrutura de diretórios (~/.syntropy/)

🔧 PRECISA CORREÇÃO:
  - Grid Token (migrar para Keyring)
  - Templates cloud-init (criar versões -mvp)
  - Setup (integrar TokenManager)

🚧 A IMPLEMENTAR (seguir MVP):
  - Node Creation Component
  - Workload Component
  - Management Component
  - Registration Protocol
  - Sync Manager

❌ NÃO IMPLEMENTADO (placeholder necessário):
  - Syntropy Agent (usar script bash por enquanto)
```

### Ordem de Implementação Recomendada

```
1. [CRÍTICO] Implementar TokenManager
   - setup/src/token_manager.go
   - Integrar com setup.go
   - Testar em Windows/Linux/macOS
   
2. [CRÍTICO] Criar templates MVP
   - user-data-mvp.yaml
   - network-config-mvp.yaml
   - meta-data-mvp.yaml
   - Incluir scripts: detect-hardware, register-node
   
3. [CRÍTICO] Criar Agent Placeholder
   - Script bash em USB/syntropy/agent
   - Status reporting funcional
   
4. [NORMAL] Implementar Node Creation
   - Seguir estrutura de componentes
   - USB detection + writing
   - Cloud-init generation
   
5. [NORMAL] Implementar Registration Protocol
   - Listener na Command Station
   - Validação de Grid Token
   - Inventory management
   
6. [NORMAL] Implementar Workload + Management
   - Deploy via SSH (fallback)
   - Node status/health
   - Sync básico
```

---

**FIM DA ESPECIFICAÇÃO MVP v2.0**

Este documento fornece **TODAS as informações necessárias** para um LLM implementar o MVP completo do Syntropy Cooperative Grid.

**⚠️ ATENÇÃO**: Antes de implementar, **LEIA** `docs/architecture/MVP-CORRECTIONS.md` para correções críticas e templates atualizados.

**Características do MVP**:
- ✅ Templates centralizados em `infrastructure/`
- ✅ Auto-detecção de hardware (via script)
- ✅ Protocol de registro automático
- ✅ Sincronização bidirecional
- ✅ Configuração mínima (apenas nome do Node)
- ✅ Cloud-init dinâmico
- ✅ Multi-plataforma (Windows/Linux)
- ✅ Detecção automática de USBs
- ✅ Estrutura por componentes/subcomponentes
- ✅ Best practices de código
- ✅ Segurança: Grid Token em Keyring, SSH key-only, Firewall
- ✅ Implementação faseada com placeholders funcionais

**Próximos passos**: 
1. Ler `MVP-CORRECTIONS.md`
2. Implementar TokenManager
3. Criar templates MVP
4. Seguir Roadmap da seção 7
