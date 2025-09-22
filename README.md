# Syntropy Cooperative Grid Management System

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)](https://github.com/syntropy-cc/cooperative-grid/actions)

> **Sistema de Gerenciamento Completo para a Syntropy Cooperative Grid**

Um sistema unificado para gerenciar nós, containers, redes e serviços cooperativos da Syntropy Cooperative Grid através de múltiplas interfaces (CLI, Web, Desktop, Mobile).

## 🎯 **Visão Geral**

O Syntropy Cooperative Grid Management System é uma plataforma abrangente que permite:

- **Criação e Gerenciamento de Nós**: Detecção automática de hardware, configuração de USB, setup de nós
- **Orquestração de Containers**: Deploy, gerenciamento e monitoramento de containers Docker/Kubernetes
- **Gerenciamento de Rede**: Configuração de service mesh, rotas e conectividade
- **Serviços Cooperativos**: Sistema de créditos, governança e economia distribuída
- **Múltiplas Interfaces**: CLI, Web, Desktop e Mobile para diferentes casos de uso

## 🏗️ **Arquitetura**

```
┌─────────────────────────────────────────────────────────────┐
│                    INTERFACES LAYER                         │
├─────────────────┬─────────────────┬─────────────────────────┤
│   CLI (Go)      │   Web (React)   │  Desktop (Electron)     │
│   Mobile (Flutter) │  API Client  │  Future Interfaces      │
└─────────────────┴─────────────────┴─────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                    API GATEWAY LAYER                        │
├─────────────────────────────────────────────────────────────┤
│  REST API  │  GraphQL  │  WebSocket  │  gRPC  │  CLI Direct │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                      CORE LAYER                             │
├─────────────┬─────────────┬─────────────┬───────────────────┤
│ Node Mgmt   │ Container   │ Network     │ Cooperative       │
│ Service     │ Service     │ Service     │ Service           │
├─────────────┼─────────────┼─────────────┼───────────────────┤
│ USB Creator │ K8s Mgmt    │ Mesh Mgmt   │ Credit System     │
│ Device Mgmt │ Runtime     │ Routing     │ Governance        │
│ Monitoring  │ Security    │ Discovery   │ Economics         │
└─────────────┴─────────────┴─────────────┴───────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                    DATA LAYER                               │
├─────────────┬─────────────┬─────────────┬───────────────────┤
│ PostgreSQL  │ Redis       │ InfluxDB    │ File Storage      │
│ (Metadata)  │ (Cache)     │ (Metrics)   │ (Configs/Logs)    │
└─────────────┴─────────────┴─────────────┴───────────────────┘
```

## 🚀 **Quick Start**

### Pré-requisitos

- Go 1.21+
- Docker & Docker Compose
- Node.js 18+ (para interfaces web)
- Git

### Instalação

```bash
# Clone o repositório
git clone https://github.com/syntropy-cc/cooperative-grid.git
cd cooperative-grid

# Instale dependências
make install

# Inicie os serviços
make up

# Execute a CLI
./interfaces/cli/bin/syntropy-cli --help
```

### Uso Básico

```bash
# Listar nós disponíveis
syntropy-cli node list

# Criar um novo nó
syntropy-cli node create --usb /dev/sdb --name "node-01"

# Deploy de container
syntropy-cli container deploy --image nginx --node node-01

# Status da rede
syntropy-cli network status
```

## 📁 **Estrutura do Projeto**

```
syntropy-cooperative-grid/
├── core/                    # Core do sistema (lógica de negócio)
├── interfaces/              # Interfaces de usuário
│   ├── cli/                # Interface CLI
│   ├── web/                # Interface web
│   ├── mobile/             # App mobile
│   ├── desktop/            # App desktop
│   └── api/                # Definições de API
├── deployments/            # Configurações de deploy
├── docs/                   # Documentação
└── scripts/                # Scripts de build e utilitários
```

## 🛠️ **Tecnologias**

### Core
- **Go 1.21+**: Linguagem principal
- **PostgreSQL**: Banco de dados principal
- **Redis**: Cache e sessões
- **InfluxDB**: Métricas e time-series

### Interfaces
- **CLI**: Go com Cobra
- **Web**: React/Next.js + TypeScript
- **Mobile**: Flutter com Dart
- **Desktop**: Electron com React

### DevOps
- **Docker Compose**: Desenvolvimento local
- **Kubernetes**: Orquestração
- **Helm**: Gerenciamento de charts
- **GitHub Actions**: CI/CD

## 📚 **Documentação**

- [Arquitetura](docs/architecture/README.md)
- [API Reference](docs/api/README.md)
- [Guia de Desenvolvimento](docs/development/README.md)
- [Deployment](docs/deployment/README.md)
- [Roadmap](docs/roadmap/README.md)

## 🤝 **Contribuição**

Veja nosso [Guia de Contribuição](CONTRIBUTING.md) para detalhes sobre como contribuir.

## 📄 **Licença**

Este projeto está licenciado sob a [Licença MIT](LICENSE).

## 🔗 **Links Úteis**

- [Website](https://syntropy.coop)
- [Documentação](https://docs.syntropy.coop)
- [Discord](https://discord.gg/syntropy)
- [Twitter](https://twitter.com/syntropy_coop)

---

**Syntropy Cooperative Grid** - Construindo o futuro da computação cooperativa 🌐