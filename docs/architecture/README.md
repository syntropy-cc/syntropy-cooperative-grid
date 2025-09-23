# 🏗️ Arquitetura do Sistema - Syntropy Cooperative Grid

> **Documentação Técnica Completa da Arquitetura**

## 📋 **Índice**

1. [Visão Geral](#visão-geral)
2. [Arquitetura de Alto Nível](#arquitetura-de-alto-nível)
3. [Core Layer](#core-layer)
4. [Interfaces Layer](#interfaces-layer)
5. [Camada de Dados](#camada-de-dados)
6. [Segurança](#segurança)
7. [Monitoramento](#monitoramento)
8. [Deployment](#deployment)

---

## 🎯 **Visão Geral**

O Syntropy Cooperative Grid Management System é uma plataforma distribuída que permite gerenciar nós, containers, redes e serviços cooperativos através de múltiplas interfaces. A arquitetura segue princípios de **separação de responsabilidades**, **escalabilidade horizontal** e **alta disponibilidade**.

### **Princípios Arquiteturais**

1. **Separação Core/Interfaces**: Core contém lógica de negócio, Interfaces são formas de acesso
2. **Microserviços**: Serviços independentes e especializados
3. **API-First**: Todas as funcionalidades expostas via APIs
4. **Cross-Platform**: Suporte nativo a Windows, Linux, macOS
5. **Real-time**: Updates em tempo real via WebSocket
6. **Scalable**: Escalabilidade horizontal automática
7. **Secure**: Segurança em múltiplas camadas
8. **Observable**: Monitoramento e observabilidade completos

---

## 🏗️ **Arquitetura de Alto Nível**

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
│  Rate Limiting │  Auth │  CORS │  Logging │  Monitoring    │
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

---

## 🧠 **Core Layer**

O **Core Layer** contém toda a lógica de negócio do sistema, independente de como o usuário acessa.

### **Estrutura do Core**

```
core/
├── services/                   # Lógica de negócio
│   ├── node/                  # Gerenciamento de nós
│   ├── container/             # Orquestração de containers
│   ├── network/               # Gerenciamento de rede
│   ├── cooperative/           # Serviços cooperativos
│   └── usb/                   # Gerenciamento de USB
├── platform/                  # Código específico de plataforma
│   ├── windows/               # Implementações Windows
│   ├── linux/                 # Implementações Linux
│   └── common/                # Código comum
├── storage/                   # Camada de dados
│   ├── postgres/              # Repositórios PostgreSQL
│   ├── redis/                 # Cache Redis
│   └── files/                 # Sistema de arquivos
├── config/                    # Configuração
├── utils/                     # Utilitários
└── types/                     # Tipos e modelos
    ├── models/                # Modelos de dados
    ├── errors/                # Tratamento de erros
    └── constants/             # Constantes
```

---

## 🖥️ **Interfaces Layer**

As **Interfaces** são diferentes formas de acesso ao Core, cada uma otimizada para seu contexto de uso.

### **Estrutura das Interfaces**

```
interfaces/
├── cli/                        # Interface de linha de comando
│   ├── cmd/                    # Comandos CLI
│   ├── internal/               # Lógica interna da CLI
│   └── pkg/                    # Pacotes da CLI
├── web/                        # Interface web
│   ├── frontend/               # React/Next.js
│   ├── backend/                # API Gateway
│   └── static/                 # Assets estáticos
├── mobile/                     # App mobile
│   └── flutter/                # Flutter app
├── desktop/                    # App desktop
│   └── electron/               # Electron app
└── api/                        # Definições de API
    ├── openapi/                # OpenAPI specs
    ├── graphql/                # GraphQL schemas
    └── proto/                  # Protocol buffers
```

### **1. CLI Interface**

#### **Características**
- **Linguagem**: Go com Cobra CLI framework
- **Uso**: Automação, scripts, administração
- **Distribuição**: Binários nativos para cada plataforma

#### **Comandos Principais**
```bash
# Gerenciamento de Nós
syntropy node list
syntropy node create --usb /dev/sdb --name "node-01"
syntropy node status <node-id>
syntropy node update <node-id> --name "new-name"
syntropy node delete <node-id>
syntropy node restart <node-id>

# Gerenciamento de Containers
syntropy container list
syntropy container deploy --image nginx --node node-01
syntropy container status <container-id>
syntropy container logs <container-id>
syntropy container start/stop <container-id>

# Gerenciamento de Rede
syntropy network status
syntropy network topology
syntropy network routes list
syntropy network routes create --source node-01 --destination node-02
syntropy network mesh status

# Serviços Cooperativos
syntropy cooperative credits balance
syntropy cooperative credits transfer --from node-01 --to node-02 --amount 100
syntropy cooperative governance proposals
syntropy cooperative governance vote --proposal prop-01 --vote yes
syntropy cooperative reputation show
```

### **2. Web Interface**

#### **Frontend (React/Next.js)**
- **Tecnologia**: React 18, Next.js 14, TypeScript
- **Uso**: Dashboard interativo, gerenciamento visual
- **Características**:
  - Dashboard interativo
  - Gerenciamento visual de nós
  - Deploy wizard para containers
  - Monitoramento em tempo real
  - Interface responsiva

#### **Backend (API Gateway)**
- **Tecnologia**: Go com Gin/Echo framework
- **Uso**: API Gateway para todas as interfaces
- **Características**:
  - Roteamento de requisições
  - Autenticação e autorização
  - Rate limiting
  - Logging e monitoramento
  - CORS e segurança

### **3. Mobile Interface**

#### **Flutter App**
- **Tecnologia**: Flutter com Dart
- **Uso**: Monitoramento remoto, ações básicas
- **Características**:
  - Monitoramento remoto
  - Notificações push
  - Ações básicas de gerenciamento
  - Biometric authentication
  - Interface nativa

### **4. Desktop Interface**

#### **Electron App**
- **Tecnologia**: Electron com React
- **Uso**: Acesso offline, notificações do sistema
- **Características**:
  - Acesso offline básico
  - Notificações do sistema
  - Tray icon para acesso rápido
  - Auto-updater
  - Interface nativa

---

## 🗄️ **Camada de Dados**

### **PostgreSQL (Primary Database)**

#### **Schema Principal**
```sql
-- Nodes table
CREATE TABLE nodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    status VARCHAR(50) NOT NULL,
    hardware_info JSONB,
    network_config JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Containers table
CREATE TABLE containers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    node_id UUID REFERENCES nodes(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    image VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    config JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Network routes table
CREATE TABLE network_routes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_node_id UUID REFERENCES nodes(id),
    destination_node_id UUID REFERENCES nodes(id),
    config JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Cooperative credits table
CREATE TABLE cooperative_credits (
    node_id UUID REFERENCES nodes(id) PRIMARY KEY,
    balance DECIMAL(15,2) DEFAULT 0,
    last_transaction_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

---

## 🔒 **Segurança**

### **Authentication & Authorization**

#### **JWT Token Structure**
```json
{
  "header": {
    "alg": "RS256",
    "typ": "JWT"
  },
  "payload": {
    "sub": "user_123",
    "iat": 1640995200,
    "exp": 1641081600,
    "roles": ["admin", "node_manager"],
    "permissions": ["nodes:read", "nodes:write", "containers:deploy"]
  }
}
```

---

## 📊 **Monitoramento**

### **Application Monitoring**

#### **Metrics Collection**
```go
// Prometheus metrics
var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    nodeStatusGauge = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "node_status",
            Help: "Current status of nodes",
        },
        []string{"node_id", "status"},
    )
)
```

---

## 🚀 **Deployment**

### **Container Orchestration**

#### **Kubernetes Manifests**
```yaml
# Core service deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: syntropy-core
spec:
  replicas: 3
  selector:
    matchLabels:
      app: syntropy-core
  template:
    metadata:
      labels:
        app: syntropy-core
    spec:
      containers:
      - name: core
        image: syntropy/core:latest
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: url
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
```

---

**Esta documentação é um documento vivo que será atualizado conforme a evolução do sistema.**