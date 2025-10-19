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
│  │   Setup     │  │    Node     │  │  Workload   │         │
│  │ Component   │  │ Component   │  │ Component   │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ Management  │  │ Registration│  │ Orchestration│         │
│  │ Component   │  │ Protocol    │  │ Engine      │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
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
**Responsabilidade**: Configuração inicial e gerenciamento de Grid Token  
**Status**: ✅ 80% implementado  
**Documentação**: [Setup Component](./components/setup.md)

**Funcionalidades**:
- Geração segura de Grid Token (via Keyring do sistema)
- Configuração inicial do Command Station
- Gerenciamento de chaves SSH
- Validação de pré-requisitos

### 2. Node Creation Component
**Responsabilidade**: Provisionamento automático de nós físicos  
**Status**: 🚧 A implementar  
**Documentação**: [Node Creation Component](./components/node-creation.md)

**Funcionalidades**:
- Detecção automática de dispositivos USB
- Download e injeção de cloud-init
- Criação de USBs bootáveis
- Provisionamento de hardware virgem

### 3. Workload Component
**Responsabilidade**: Orquestração inteligente de workloads  
**Status**: 🚧 A implementar  
**Documentação**: [Workload Component](./components/workload.md)

**Funcionalidades**:
- Admission Control (validação de recursos)
- Intelligent Scheduler (3 estratégias)
- Queue System (gerenciamento de filas)
- Deploy Execution (execução via SSH)
- Lifecycle Management (start/stop/scale)
- Monitoring (logs e métricas)

### 4. Management Component
**Responsabilidade**: Gerenciamento e monitoramento da grid  
**Status**: 🚧 A implementar  
**Documentação**: [Management Component](./components/management.md)

**Funcionalidades**:
- Listagem de nós
- Status e health checks
- Sincronização de inventário
- Descoberta de serviços

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

### Fase 1: Setup e Segurança (Semana 1)
1. Implementar TokenManager com Keyring
2. Completar Setup Component
3. Testes em múltiplas plataformas

### Fase 2: Node Creation (Semanas 2-3)
1. Implementar USB detection (Windows/Linux)
2. Implementar cloud-init injection
3. Implementar USB writing
4. Testar provisionamento de primeiro nó

### Fase 3: Registration Protocol (Semana 4)
1. Implementar Listener (Command Station)
2. Implementar Node announcement
3. Implementar Inventory management
4. Testar registro completo

### Fase 4: Workload Orchestration (Semana 5)
1. Implementar Admission Control
2. Implementar Scheduler (3 estratégias)
3. Implementar Queue System
4. Implementar Deploy Execution
5. Testes end-to-end

### Fase 5: Management e Polish (Semana 6)
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

### Setup
```bash
syntropy setup run                    # Configuração inicial
syntropy token show                   # Ver Grid Token
syntropy token rotate                 # Gerar novo token
```

### Node Management
```bash
syntropy node create                  # Criar USB para nó
syntropy node listen                  # Iniciar listener
syntropy node list                    # Listar nós registrados
syntropy node status <node-id>        # Status de nó específico
```

### Workload Management
```bash
syntropy workload deploy nginx --replicas 3 --cpu 1 --memory 512M
syntropy workload list                # Listar workloads
syntropy workload status <workload-id> # Status de workload
syntropy workload scale <workload-id> --replicas 5
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

1. **[Setup Component](./components/setup.md)** - Configuração e segurança
2. **[Node Creation Component](./components/node-creation.md)** - Provisionamento de nós
3. **[Workload Component](./components/workload.md)** - Orquestração de workloads
4. **[Management Component](./components/management.md)** - Gerenciamento da grid
5. **[Registration Protocol](./components/registration.md)** - Protocolo de registro
6. **[Orchestration Engine](./components/orchestration.md)** - Motor de orquestração

---

## ✅ CRITÉRIOS DE SUCESSO

O MVP está completo quando:

- ✅ 6 nós físicos provisionados e registrados
- ✅ Grid Token seguro via Keyring
- ✅ Deploy de workloads funcionando
- ✅ Orquestração automática (Admission + Scheduler)
- ✅ Queue system para grid cheia
- ✅ Comandos CLI funcionais
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