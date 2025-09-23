# 🎛️ Regras para Management System

> **Regras técnicas e sucintas para LLMs trabalharem com o Syntropy Cooperative Grid Management System**

## 📋 **Visão Geral**

O Management System é um sistema unificado para gerenciar a Syntropy Cooperative Grid, fornecendo interfaces CLI, Web, Mobile e Desktop para administrar nós, containers, redes e serviços cooperativos.

## 🏗️ **Arquitetura**

### **Estrutura de Diretórios do Projeto**
```
syntropy-cooperative-grid/
├── interfaces/                    # Interfaces do Management System
│   ├── cli/                      # CLI Interface (Go + Cobra)
│   │   ├── cmd/                  # Comandos principais
│   │   │   └── main.go          # Entry point do CLI
│   │   ├── internal/             # Lógica interna do CLI
│   │   │   └── cli/             # Implementação dos comandos
│   │   │       ├── config.go    # Configurações
│   │   │       ├── container.go # Comandos de container
│   │   │       ├── cooperative.go # Comandos cooperativos
│   │   │       ├── network.go   # Comandos de rede
│   │   │       └── node.go      # Comandos de nó
│   │   └── go.mod               # Dependências Go
│   ├── web/                      # Web Interface
│   │   ├── frontend/             # React + Next.js
│   │   │   └── package.json     # Dependências frontend
│   │   └── backend/              # Go backend
│   │       └── go.mod           # Dependências backend
│   ├── mobile/                   # Mobile Interface
│   │   └── flutter/             # Flutter app
│   │       └── pubspec.yaml     # Dependências Flutter
│   └── desktop/                  # Desktop Interface
│       └── electron/            # Electron app
│           └── package.json     # Dependências Electron
├── core/                         # Management Core (lógica de negócio)
├── internal/                     # Código interno compartilhado
├── services/                     # Microserviços
├── docs/                         # Documentação
├── rules/                        # Regras para LLMs
└── scripts/                      # Scripts de automação
```

### **Componentes Principais**
- **Management Core**: Lógica de gerenciamento (Go) - `core/`
- **Management Interfaces**: CLI, Web, Mobile, Desktop - `interfaces/`
- **Grid Infrastructure**: Nós físicos, virtuais, cloud e edge

### **Engines do Core**
- **Node Management Engine**: Detecção, configuração, monitoramento de nós
- **Container Management Engine**: Deploy, orquestração, escalabilidade
- **Network Management Engine**: Service mesh, roteamento, conectividade
- **Cooperative Management Engine**: Créditos, governança, reputação

## 🖥️ **Interfaces Disponíveis**

### **CLI (Go + Cobra)**
```bash
# Comandos principais
syntropy node create --usb /dev/sdb --name "node-01"
syntropy node list --format table --filter running
syntropy container deploy --template nginx --node node-01 --scale 3
syntropy network mesh enable --encryption --monitoring
syntropy cooperative credits balance --node node-01
```

### **Web Interface (React + Next.js)**
- Dashboard em tempo real
- Gerenciamento visual de nós e containers
- Visualização de topologia de rede
- Interface de governança cooperativa

### **Mobile (Flutter)**
- Monitoramento remoto
- Notificações push
- Ações básicas de gerenciamento

### **Desktop (Electron)**
- Aplicação nativa
- Tray icon para acesso rápido
- Notificações do sistema

## 🔧 **Funcionalidades Core**

### **Gerenciamento de Nós**
- **Auto-discovery**: Detecção automática de dispositivos USB
- **Cross-platform**: Windows, Linux, WSL
- **Configuração automática**: Formatação, chaves SSH, configurações
- **Monitoramento**: Health checks, métricas, alertas

### **Gerenciamento de Containers**
- **Deploy**: Templates, multi-cluster, service discovery
- **Escalabilidade**: Auto-scaling baseado em métricas
- **Segurança**: Container security, image scanning
- **Orquestração**: Load balancing, rolling updates

### **Gerenciamento de Rede**
- **Service Mesh**: Configuração automática, traffic management
- **Roteamento**: Dynamic routing, failover, QoS
- **Monitoramento**: Network analytics, anomaly detection

### **Gerenciamento Cooperativo**
- **Sistema de Créditos**: Credit management, transações
- **Governança**: Propostas, votação, compliance
- **Reputação**: Scoring, trust metrics, behavioral analysis

## 🛠️ **Stack Tecnológico**

### **Backend**
- **Go 1.21+**: Linguagem principal
- **Gin/Echo**: Web framework
- **PostgreSQL**: Banco principal
- **Redis**: Cache e sessões
- **Docker API**: Integração com containers

### **Frontend Web**
- **Next.js 14**: React framework
- **TypeScript**: Type safety
- **Tailwind CSS**: Styling
- **Zustand**: State management
- **Recharts**: Visualizações

### **Mobile**
- **Flutter 3.10+**: Cross-platform
- **Riverpod**: State management
- **Dio**: HTTP client

### **Desktop**
- **Electron 27**: Desktop framework
- **React**: UI reutilizado

## 📊 **Casos de Uso Principais**

### **Administrador de Sistema**
- Monitorar 1000+ nós distribuídos
- Aplicar atualizações de segurança
- Gerenciar recursos e performance
- Responder a incidentes

### **Desenvolvedor**
- Deploy de aplicações distribuídas
- Configurar múltiplos serviços
- Monitorar performance
- Escalar conforme demanda

### **Operador de Rede**
- Otimizar conectividade
- Gerenciar largura de banda
- Implementar políticas de segurança
- Monitorar QoS

### **Participante Cooperativo**
- Monitorar créditos e transações
- Participar de votações
- Gerenciar reputação
- Otimizar contribuições

## 🎯 **Roadmap de Desenvolvimento**

### **Fase 1: MVP CLI (Sprints 1-4)**
- ✅ Foundation & Setup
- ✅ USB Detection & Node Creation
- ✅ Node Management
- ✅ Container Basics

### **Fase 2: API Foundation (Sprints 5-8)**
- 🔄 API Gateway
- 🔄 Database & Models
- 🔄 Microservices Architecture
- 🔄 Real-time Features

### **Fase 3: Web Interface (Sprints 9-12)**
- ⏳ Web Foundation
- ⏳ Node Management UI
- ⏳ Container Management UI
- ⏳ Dashboard & Monitoring

### **Fases 4-6: Advanced Features**
- Network Management
- Cooperative Services
- Mobile & Desktop
- Production Ready

## ⚡ **Comandos Essenciais**

### **Nós**
```bash
# Criar nó
syntropy node create --usb /dev/sdb --name "prod-node-01" --auto-config

# Listar nós
syntropy node list --format table --filter running

# Status detalhado
syntropy node status node-01 --watch --format json

# Atualizar configuração
syntropy node update node-01 --config-file production.yaml
```

### **Containers**
```bash
# Deploy container
syntropy container deploy --template nginx --node node-01 --scale 3

# Listar containers
syntropy container list --node node-01 --status running

# Ver logs
syntropy container logs container-01 --follow --tail 100

# Escalar
syntropy container scale container-01 --replicas 5
```

### **Rede**
```bash
# Habilitar service mesh
syntropy network mesh enable --encryption --monitoring

# Criar rotas
syntropy network routes create --source node-01 --dest node-02 --priority 1

# Ver topologia
syntropy network topology --format graphviz
```

### **Cooperativo**
```bash
# Ver saldo
syntropy cooperative credits balance --node node-01

# Transferir créditos
syntropy cooperative credits transfer --from node-01 --to node-02 --amount 100

# Votar em proposta
syntropy cooperative governance vote --proposal prop-01 --vote yes
```

## 🔒 **Segurança e Compliance**

### **Autenticação**
- JWT tokens
- 2FA/MFA
- Biometric authentication (mobile)

### **Autorização**
- RBAC (Role-Based Access Control)
- Permissões granulares
- Audit logs

### **Criptografia**
- End-to-end encryption
- Encryption at rest
- Secure key management

## 📈 **Métricas e Monitoramento**

### **Performance**
- Response time < 200ms
- 99.9% uptime
- Suporte a 1000+ usuários simultâneos

### **Qualidade**
- Test coverage > 90%
- < 5 bugs críticos por sprint
- 100% code review coverage

### **Negócio**
- 50% redução em custos operacionais
- 60% redução em tempo de setup
- 80% redução em erros

## 🚨 **Tratamento de Erros**

### **Padrões de Erro**
- Códigos de erro consistentes
- Mensagens claras e acionáveis
- Logs estruturados
- Recovery automático quando possível

### **Validações**
- Input validation rigorosa
- Configuração validation
- Cross-platform compatibility checks

## 🔄 **Integração e APIs**

### **APIs Disponíveis**
- REST API (Gin/Echo)
- GraphQL (gqlgen)
- WebSocket para real-time
- gRPC para microserviços

### **Integrações**
- Docker API
- Kubernetes API
- Prometheus metrics
- Slack/Discord notifications

## 📝 **Regras para LLMs**

### **Estrutura de Diretórios - Regras Críticas:**
1. **NUNCA crie arquivos fora dos diretórios corretos**
2. **CLI**: Use apenas `interfaces/cli/` para código CLI
3. **Web Frontend**: Use apenas `interfaces/web/frontend/` para React/Next.js
4. **Web Backend**: Use apenas `interfaces/web/backend/` para Go backend
5. **Mobile**: Use apenas `interfaces/mobile/flutter/` para Flutter
6. **Desktop**: Use apenas `interfaces/desktop/electron/` para Electron
7. **Core Logic**: Use `core/` para lógica de negócio compartilhada
8. **Shared Code**: Use `internal/` para código interno compartilhado
9. **Services**: Use `services/` para microserviços
10. **Documentation**: Use `docs/` para documentação

### **Ao trabalhar com Management System:**
1. **Sempre use comandos CLI** quando possível para automação
2. **Valide inputs** antes de executar operações
3. **Use formatos apropriados** (table, json, yaml) conforme contexto
4. **Implemente error handling** robusto
5. **Siga padrões de nomenclatura** consistentes
6. **Use templates** para deployments comuns
7. **Monitore operações** com logs estruturados
8. **Respeite permissões** e segurança
9. **Documente configurações** importantes
10. **Teste em ambiente isolado** antes de produção

### **Exemplos de Estrutura de Arquivos por Interface:**

#### **CLI Interface (`interfaces/cli/`)**
```
interfaces/cli/
├── cmd/
│   └── main.go                    # Entry point
├── internal/
│   └── cli/
│       ├── config.go             # Configurações
│       ├── container.go          # Comandos de container
│       ├── cooperative.go        # Comandos cooperativos
│       ├── network.go            # Comandos de rede
│       └── node.go               # Comandos de nó
└── go.mod                        # Dependências
```

#### **Web Interface (`interfaces/web/`)**
```
interfaces/web/
├── frontend/                     # React + Next.js
│   ├── src/
│   │   ├── components/           # Componentes React
│   │   ├── pages/               # Páginas Next.js
│   │   ├── hooks/               # Custom hooks
│   │   └── utils/               # Utilitários
│   ├── public/                  # Assets estáticos
│   └── package.json             # Dependências
└── backend/                     # Go backend
    ├── cmd/
    │   └── server/
    │       └── main.go          # Entry point
    ├── internal/
    │   ├── handlers/            # HTTP handlers
    │   ├── services/            # Business logic
    │   └── models/              # Data models
    └── go.mod                   # Dependências
```

#### **Mobile Interface (`interfaces/mobile/flutter/`)**
```
interfaces/mobile/flutter/
├── lib/
│   ├── main.dart                # Entry point
│   ├── screens/                 # Telas da aplicação
│   ├── widgets/                 # Widgets customizados
│   ├── services/                # Serviços e APIs
│   └── models/                  # Modelos de dados
├── assets/                      # Imagens, fontes, etc.
└── pubspec.yaml                 # Dependências
```

#### **Desktop Interface (`interfaces/desktop/electron/`)**
```
interfaces/desktop/electron/
├── src/
│   ├── main.js                  # Processo principal
│   ├── renderer/                # Processo de renderização
│   │   ├── components/          # Componentes React
│   │   └── pages/               # Páginas
│   └── preload.js               # Script de preload
├── public/                      # Assets estáticos
└── package.json                 # Dependências
```

### **Comandos de Troubleshooting:**
```bash
# Verificar saúde geral
syntropy node status --all --format json

# Ver logs de sistema
syntropy logs system --tail 100 --follow

# Verificar conectividade
syntropy network health --detailed

# Backup de configurações
syntropy backup create --include-nodes --include-configs
```

### **Padrões de Configuração:**
- Use nomes descritivos para nós
- Configure monitoring desde o início
- Implemente backup automático
- Use templates para consistência
- Documente mudanças importantes
