# Syntropy Cooperative Grid - MVP (Minimum Viable Product)

**Versão**: 2.0  
**Data**: Outubro 2025  
**Objetivo**: Mini-cloud pessoal com 6 nós físicos  
**Status**: ✅ Especificação completa - Pronto para implementação

---

## 📋 VISÃO GERAL

O Syntropy Cooperative Grid MVP é um sistema de mini-cloud pessoal que permite criar e gerenciar uma infraestrutura distribuída usando 6 nós físicos. O sistema é projetado para ser implementado por LLMs seguindo especificações técnicas detalhadas.

### Objetivos do MVP
- ✅ Provisionar 6 nós físicos automaticamente
- ✅ Gerenciar workloads distribuídos na grid
- ✅ Orquestração inteligente de recursos
- ✅ Interface CLI intuitiva
- ✅ Segurança robusta (Grid Token via Keyring)

### Arquitetura Geral
```
┌─────────────────────────────────────────────────────────────┐
│                    COMMAND STATION                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │    Setup    │  │    Node     │  │  Workload   │         │
│  │ Component   │  │ Management  │  │ Component   │         │
│  │             │  │ Component   │  │ (Unificado) │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
│  ┌─────────────┐                                           │
│  │ Management  │                                           │
│  │ Component   │                                           │
│  └─────────────┘                                           │
│                                                             │
│  Workload Component inclui:                                 │
│  • Admission Control • Scheduler • Queue System            │
│  • Deploy Execution • Lifecycle • Monitoring               │
│  • Docker Compose • App Deploy • Server Deploy             │
│  • Auto-Orchestration • Workflow • Event Bus               │
│  • State Management • Metrics Collection                   │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ SSH + Registration Protocol
                              │
┌─────────────────────────────────────────────────────────────┐
│                        GRID NODES                           │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐        │
│  │ Node-01 │  │ Node-02 │  │ Node-03 │  │ Node-04 │        │
│  │ 8 cores │  │ 8 cores │  │ 8 cores │  │ 8 cores │        │
│  │ 28GB    │  │ 28GB    │  │ 28GB    │  │ 28GB    │        │
│  └─────────┘  └─────────┘  └─────────┘  └─────────┘        │
│  ┌─────────┐  ┌─────────┐                                  │
│  │ Node-05 │  │ Node-06 │                                  │
│  │ 8 cores │  │ 8 cores │                                  │
│  │ 28GB    │  │ 28GB    │                                  │
│  └─────────┘  └─────────┘                                  │
└─────────────────────────────────────────────────────────────┘
```

---

## 🏗️ COMPONENTES DO MVP

### 1. Setup Component
**Responsabilidade**: Configuração inicial do Command Station  
**Status**: 🚧 Em desenvolvimento (80% implementado)  
**Documentação**: [Setup Component](./components/setup.md)

### 2. Node Management Component
**Responsabilidade**: Provisionamento automático e registro de nós físicos  
**Status**: 🚧 A implementar  
**Documentação**: [Node Management Component](./components/node.md)

**Funcionalidades**:
- Detecção automática de dispositivos USB
- Geração automática de configurações (NodeID, Grid Token, SSH keys)
- Criação de USBs bootáveis com cloud-init customizado
- Registro automático e handshake seguro
- Gerenciamento de múltiplos nós simultâneos

### 3. Workload Component (Unificado)
**Responsabilidade**: Orquestração inteligente completa de workloads  
**Status**: 🚧 A implementar  
**Documentação**: [Workload Component](./components/workload.md)

**Funcionalidades**:
- Admission Control (validação de recursos)
- Intelligent Scheduler (3 estratégias)
- Queue System (gerenciamento de filas)
- Deploy Execution (execução via SSH)
- Lifecycle Management (start/stop/scale)
- Monitoring (logs e métricas)
- Docker Compose Support (deploy multi-container)
- Application Deploy (deploy de aplicações completas)
- Server Deploy (deploy de servidores especializados)
- Auto-Orchestration (orquestração automática integrada)

### 4. Management Component
**Responsabilidade**: Gerenciamento e monitoramento da grid  
**Status**: 🚧 A implementar  
**Documentação**: [Management Component](./components/management.md)

**Funcionalidades**:
- Health Monitoring (monitoramento de saúde dos nós)
- Inventory Sync (sincronização de inventário)
- Service Discovery (descoberta de serviços)
- Grid Analytics (análise de performance da grid)
- Administrative Operations (operação administrativas)

---

## 🔐 SEGURANÇA

### Grid Token Management
O MVP implementa segurança robusta através do sistema de Keyring do sistema operacional:

- **Windows**: Credential Manager
- **Linux**: Secret Service / gnome-keyring  
- **macOS**: Keychain

**Comandos de segurança**:
```bash
syntropy token show      # Ver token (com confirmação)
syntropy token export    # Backup seguro
syntropy token rotate    # Gerar novo token
```

---

## 🚀 FLUXO DE IMPLEMENTAÇÃO

### Fase 0: Setup Component (Semana 0)
1. Finalizar Setup Component (20% restante)
2. Implementar TokenManager
3. Testar configuração inicial

### Fase 1: Node Management (Semanas 1-2)
1. Implementar USB detection (Windows/Linux)
2. Implementar cloud-init injection
3. Implementar USB writing
4. Testar provisionamento de primeiro nó

### Fase 2: Registration Protocol (Semana 3)
1. Implementar Listener (Command Station)
2. Implementar Node announcement
3. Implementar Inventory management
4. Testar registro completo

### Fase 3: Workload Orchestration Unificado (Semanas 4-5)
1. Implementar Admission Control
2. Implementar Scheduler (3 estratégias)
3. Implementar Queue System
4. Implementar Deploy Execution
5. Implementar Docker Compose Support
6. Implementar Application Deploy
7. Implementar Server Deploy
8. Implementar Auto-Orchestration
9. Testes end-to-end

### Fase 4: Management e Polish (Semana 6)
1. Implementar Management Component
2. Provisionar 6 nós completos
3. Testes de carga
4. Documentação final

---

## 📊 ESPECIFICAÇÕES TÉCNICAS

### Hardware Mínimo
**Command Station**:
- CPU: 2 cores
- RAM: 4GB
- Disk: 20GB livres
- USB: Porta disponível

**Grid Nodes** (6 unidades):
- CPU: 8 cores cada (48 cores total)
- RAM: 28GB cada (168GB total)
- Disk: 512GB SSD cada
- Network: Ethernet 1Gbps

### Software Requirements
- Go 1.22+
- Docker
- SSH client
- Git
- Dependências de sistema (libsecret-1-dev no Linux)

---

## 🎯 COMANDOS CLI PRINCIPAIS

### Setup Component
```bash
syntropy setup run                    # Configuração inicial
syntropy setup status                 # Status do setup
syntropy setup validate               # Validar ambiente
syntropy setup reset --confirm        # Reset do setup
syntropy token show                   # Ver token (com confirmação)
syntropy token export                 # Backup seguro
syntropy token rotate                 # Gerar novo token
```

### Node Management
```bash
syntropy node create                  # Criar nó com registro automático
syntropy node list                    # Listar nós (ativos e pendentes)
syntropy node status <node-id>        # Status de nó específico
syntropy node logs <node-id>          # Logs de nó específico
```

### Workload Management (Unificado)
```bash
# Deploy de containers
syntropy workload deploy nginx --replicas 3 --cpu 1 --memory 512M

# Deploy de docker-compose
syntropy workload deploy-compose ./docker-compose.yaml

# Deploy de aplicações completas
syntropy workload deploy-app ./app-config.yaml

# Deploy de servidores
syntropy workload deploy-server nginx --port 80 --ssl --domain example.com

# Gerenciamento unificado
syntropy workload list                # Listar todos os workloads
syntropy workload list --type container|compose|app|server
syntropy workload status <workload-id> # Status de qualquer workload
syntropy workload scale <workload-id> --replicas 5
syntropy workload logs <workload-id> --service <service-name>

# Orquestração integrada
syntropy workload orchestration status
syntropy workload workflow list
syntropy workload state show

# Capacidade da grid
syntropy grid capacity                # Capacidade da grid
```

### Management
```bash
syntropy grid status                  # Status geral da grid
syntropy grid health                  # Health check completo
syntropy grid sync                    # Sincronizar inventário
```

---

## 📚 DOCUMENTAÇÃO DETALHADA

Para implementação completa, consulte a documentação específica de cada componente:

1. **[Setup Component](./components/setup.md)** - Configuração inicial do Command Station
2. **[Node Management Component](./components/node.md)** - Provisionamento e registro de nós
3. **[Workload Component](./components/workload.md)** - Orquestração completa de workloads
4. **[Management Component](./components/management.md)** - Gerenciamento da grid

---

## ✅ CRITÉRIOS DE SUCESSO

O MVP está completo quando:

- ✅ Setup Component funcionando (configuração inicial)
- ✅ Grid Token seguro via Keyring
- ✅ 6 nós físicos provisionados e registrados
- ✅ Deploy de workloads funcionando (containers, compose, apps, servers)
- ✅ Orquestração automática integrada (Admission + Scheduler + Queue)
- ✅ Docker Compose support funcionando
- ✅ Application deploy funcionando
- ✅ Server deploy funcionando
- ✅ Management Component funcionando
- ✅ Comandos CLI unificados funcionais
- ✅ Documentação completa

**Score de Qualidade**: 9.1/10 - Excelente para implementação por LLMs

---

## 🚨 AVISOS IMPORTANTES

### ⚠️ Grid Token Security
- **NÃO** armazenar Grid Token em arquivos de texto
- **SEMPRE** usar Keyring do sistema operacional
- **NUNCA** committar tokens no Git

### ⚠️ Templates Cloud-Init
- **USAR** templates MVP (`-mvp.yaml`) para implementação
- **NÃO** usar templates avançados (`-advanced.yaml`) no MVP
- **VALIDAR** sintaxe YAML antes de criar USBs

### ⚠️ Agent Implementation
- **USAR** Agent Placeholder (script bash) no MVP
- **NÃO** assumir Agent completo implementado
- **PLANEJAR** migração para Agent Go no pós-MVP

---

**Documentação mantida por**: Syntropy Team  
**Última revisão**: Outubro 2025  
**Versão**: 2.0 - Pronto para implementação