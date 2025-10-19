# Workload Component - Especificação Completa

**Componente**: Workload  
**Versão**: 1.0  
**Objetivo**: Orquestração inteligente de workloads na Grid  
**Status**: 🚧 A implementar  

---

## 📊 VISÃO GERAL DO COMPONENTE

### Responsabilidade Principal
Gerenciar todo lifecycle de workloads: validação, placement, deployment, execução, monitoramento e scaling.

### 6 Subcomponentes
1. **Admission** - Valida se workload é possível
2. **Scheduler** - Decide onde deployar
3. **Queue** - Gerencia workloads pendentes
4. **Deploy** - Executa deployment
5. **Lifecycle** - Gerencia start/stop/scale
6. **Monitoring** - Logs e métricas

---

## 🏗️ ESTRUTURA DE ARQUIVOS (26 arquivos)

```
workload/
├── README.md                           # Documentação do componente
├── ARCHITECTURE.md                     # Arquitetura e fluxos
├── workload.go                         # Orquestrador (< 500 linhas)
│
├── admission/                          # Subcomponente 1: Admission Control
│   ├── README.md                       # O que é Admission Control
│   ├── admission_controller.go         # [400 linhas] Validação principal
│   ├── capacity_calculator.go          # [300 linhas] Cálculo de capacidade
│   ├── constraint_validator.go         # [300 linhas] Validação de constraints
│   ├── resource_validator.go           # [300 linhas] Validação de recursos
│   └── tests/
│       ├── admission_test.go           # Testes unitários
│       └── capacity_test.go            # Testes de cálculo
│
├── scheduler/                          # Subcomponente 2: Scheduler
│   ├── README.md                       # Estratégias de scheduling
│   ├── scheduler.go                    # [400 linhas] Scheduler principal
│   ├── node_filter.go                  # [300 linhas] Filtragem de Nodes
│   ├── node_scorer.go                  # [300 linhas] Cálculo de scores
│   ├── strategy_spread.go              # [300 linhas] Estratégia spread
│   ├── strategy_binpack.go             # [300 linhas] Estratégia binpack
│   ├── strategy_optimized.go           # [400 linhas] Estratégia optimized
│   └── tests/
│       ├── scheduler_test.go           # Testes gerais
│       ├── spread_test.go              # Testes spread
│       ├── binpack_test.go             # Testes binpack
│       └── optimized_test.go           # Testes optimized
│
├── queue/                              # Subcomponente 3: Queue
│   ├── README.md                       # Sistema de filas
│   ├── queue_manager.go                # [300 linhas] Gerenciador
│   ├── queue_processor.go              # [300 linhas] Processador periódico
│   ├── wait_estimator.go               # [200 linhas] Estimativa de tempo
│   ├── priority_manager.go             # [200 linhas] Prioridades
│   └── tests/
│       └── queue_test.go               # Testes
│
├── deploy/                             # Subcomponente 4: Deploy
│   ├── README.md                       # Execução de deployments
│   ├── deployer.go                     # [400 linhas] Orquestrador
│   ├── executor.go                     # [400 linhas] Executor base
│   ├── executor_windows.go             # [300 linhas] Windows
│   ├── executor_linux.go               # [300 linhas] Linux
│   ├── rollback.go                     # [300 linhas] Rollback
│   ├── docker_client.go                # [400 linhas] Cliente Docker
│   └── tests/
│       ├── deployer_test.go            # Testes deployer
│       └── rollback_test.go            # Testes rollback
│
├── lifecycle/                          # Subcomponente 5: Lifecycle
│   ├── README.md                       # Gerenciamento de lifecycle
│   ├── lifecycle.go                    # [300 linhas] Manager
│   ├── start.go                        # [200 linhas] Start
│   ├── stop.go                         # [200 linhas] Stop
│   ├── restart.go                      # [200 linhas] Restart
│   ├── scale.go                        # [300 linhas] Scale
│   └── tests/
│       └── lifecycle_test.go           # Testes
│
└── monitoring/                         # Subcomponente 6: Monitoring
    ├── README.md                       # Observabilidade
    ├── monitoring.go                   # [300 linhas] Monitor
    ├── logs.go                         # [400 linhas] Logs
    ├── metrics.go                      # [400 linhas] Métricas
    └── tests/
        └── monitoring_test.go          # Testes

TOTAL: 26 arquivos, ~7,500 linhas (estimado)
Seguindo best practices: cada arquivo < 500 linhas
```

---

## 🔄 FLUXO DE DEPLOYMENT INTEGRADO

```
┌─────────────────────────────────────────────────────────────┐
│  COMANDO DO USUÁRIO                                         │
└─────────────────────────────────────────────────────────────┘
   $ syntropy workload deploy nginx --replicas 3 --cpu 1 --memory 512M
                              ↓
┌─────────────────────────────────────────────────────────────┐
│  workload.go (Orquestrador Principal)                        │
│  - ParseRequest()                                            │
│  - CreateWorkloadRequest{                                    │
│      Image: "nginx",                                         │
│      Replicas: 3,                                            │
│      CPUPerReplica: 1,                                       │
│      RAMPerReplica: 0.5,                                     │
│      Strategy: "spread"  // padrão                           │
│    }                                                         │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│  SUBCOMPONENTE 1: admission/                                 │
│  admission_controller.go                                     │
│                                                              │
│  [1] capacity_calculator.go                                  │
│      → GetGridCapacity()                                     │
│      → Total: 48 cores, 168 GB RAM                          │
│      → Available: 42 cores, 155 GB RAM                      │
│      → Utilization: CPU 13%, RAM 8%                         │
│                                                              │
│  [2] resource_validator.go                                   │
│      → ValidateCPU(needed: 3, available: 42)  ✅            │
│      → ValidateRAM(needed: 1.5, available: 155) ✅          │
│                                                              │
│  [3] Projetar utilização                                     │
│      → Projected CPU: 19% (13% + 6%)                        │
│      → Projected RAM: 9% (8% + 1%)                          │
│      → Check: < 90%? ✅ YES                                 │
│                                                              │
│  [4] constraint_validator.go                                 │
│      → Constraints: none                                     │
│      → ✅ PASS                                              │
│                                                              │
│  OUTPUT: AdmissionResult{                                    │
│    Admitted: true,                                           │
│    Reason: "Workload meets all requirements",               │
│    GridCapacity: {...}                                       │
│  }                                                           │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│  SUBCOMPONENTE 2: scheduler/                                 │
│  scheduler.go                                                │
│                                                              │
│  [1] node_filter.go                                          │
│      → FilterHealthyNodes() → 6 nodes                       │
│      → FilterByCapacity(cpu: 1, ram: 0.5) → 6 nodes         │
│      → Eligible: 6 nodes                                     │
│                                                              │
│  [2] strategy_spread.go (estratégia selecionada)             │
│      → OrderByLeastWorkloads()                              │
│      → node-01: 0 workloads (score: 100.0)                  │
│      → node-02: 0 workloads (score: 100.0)                  │
│      → node-03: 0 workloads (score: 100.0)                  │
│      → ...                                                   │
│                                                              │
│  [3] node_scorer.go                                          │
│      → CalculateSpreadScore() para cada Node                │
│      → Round-robin allocation                               │
│                                                              │
│  OUTPUT: []PlacementDecision{                                │
│    {NodeID: "node-01", ReplicaIndex: 1, Score: 100.0},      │
│    {NodeID: "node-02", ReplicaIndex: 2, Score: 100.0},      │
│    {NodeID: "node-03", ReplicaIndex: 3, Score: 100.0}       │
│  }                                                           │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│  SUBCOMPONENTE 4: deploy/                                    │
│  deployer.go                                                 │
│                                                              │
│  Para cada PlacementDecision:                                │
│                                                              │
│  [Replica 1 → node-01]                                       │
│  ├─ executor.go                                              │
│  │  ├─ ValidateImage("nginx") ✅                            │
│  │  ├─ PullImage(node-01, "nginx")                          │
│  │  │  └─ docker_client.go                                  │
│  │  │     └─ SSH: docker pull nginx                         │
│  │  ├─ CreateContainer(config)                              │
│  │  │  └─ SSH: docker create --cpu 1 --memory 512M nginx    │
│  │  ├─ StartContainer(containerID)                          │
│  │  │  └─ SSH: docker start <container>                     │
│  │  └─ VerifyRunning() ✅                                   │
│  └─ ✅ Replica 1 deployed                                   │
│                                                              │
│  [Replica 2 → node-02] ... (mesmo processo)                 │
│  [Replica 3 → node-03] ... (mesmo processo)                 │
│                                                              │
│  Se FALHA em qualquer réplica:                               │
│  └─ rollback.go                                              │
│     └─ Rollback de todas as réplicas já deployadas          │
│                                                              │
│  OUTPUT: DeploymentResult{                                   │
│    Success: true,                                            │
│    WorkloadID: "nginx-153045",                              │
│    ReplicasRunning: 3,                                       │
│    Placements: [...]                                         │
│  }                                                           │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│  BACKGROUND JOBS                                             │
│                                                              │
│  [1] queue/ (queue_processor.go)                             │
│      Goroutine rodando a cada 1 minuto:                      │
│      - Processar fila de workloads pendentes                 │
│      - Tentar admission novamente                            │
│      - Deploy automático se aprovado                         │
│                                                              │
│  [2] monitoring/ (monitoring.go)                             │
│      Goroutine rodando a cada 30 segundos:                   │
│      - Coletar logs dos containers                           │
│      - Coletar métricas (CPU, RAM, Network)                  │
│      - Atualizar status dos workloads                        │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│  RESULTADO FINAL                                             │
│                                                              │
│  ✅ Workload nginx-153045 deployed: 3/3 replicas running    │
│                                                              │
│  Placement:                                                  │
│    node-01: nginx-153045-1 (running)                        │
│    node-02: nginx-153045-2 (running)                        │
│    node-03: nginx-153045-3 (running)                        │
│                                                              │
│  Resources:                                                  │
│    Total: 3 CPU cores, 1.5 GB RAM                           │
│    Grid utilization: CPU 19%, RAM 9%                        │
└─────────────────────────────────────────────────────────────┘
```

---

## 📋 CONTAGEM DE ARQUIVOS POR SUBCOMPONENTE

| Subcomponente | Arquivos | Linhas Estimadas | Testes |
|---------------|----------|------------------|--------|
| **admission/** | 4 arquivos | ~1,300 linhas | 2 arquivos |
| **scheduler/** | 6 arquivos | ~2,200 linhas | 4 arquivos |
| **queue/** | 4 arquivos | ~1,000 linhas | 1 arquivo |
| **deploy/** | 6 arquivos | ~2,200 linhas | 2 arquivos |
| **lifecycle/** | 5 arquivos | ~1,200 linhas | 1 arquivo |
| **monitoring/** | 3 arquivos | ~1,100 linhas | 1 arquivo |
| **workload.go** | 1 arquivo | ~500 linhas | - |
| **README.md + ARCH** | 2 arquivos | - | - |
| **TOTAL** | **31 arquivos** | **~9,500 linhas** | **11 arquivos** |

**Seguindo Best Practices**:
- ✅ Cada arquivo < 500 linhas
- ✅ Separação clara de responsabilidades
- ✅ Testes para cada subcomponente
- ✅ Documentação por subcomponente
- ✅ Multi-plataforma (Windows/Linux)

---

## 🎯 SUBCOMPONENTE 1: ADMISSION

### Responsabilidade
Validar se workload PODE ser aceito pela Grid antes de tentar deployar.

### Arquivos (4)

```
admission/
├── admission_controller.go    [400 linhas]
│   └─ Validate()
│      ├─ Calcular recursos totais
│      ├─ Obter capacidade da Grid
│      ├─ Validar CPU/RAM disponível
│      ├─ Validar limite de sobrecarga (90%)
│      └─ Validar constraints
│
├── capacity_calculator.go     [300 linhas]
│   └─ GetGridCapacity()
│      ├─ Listar Nodes
│      ├─ Somar recursos totais
│      ├─ Calcular recursos disponíveis
│      └─ Calcular utilização (%)
│
├── constraint_validator.go    [300 linhas]
│   └─ ValidateConstraints()
│      ├─ ValidateTag()
│      ├─ ValidateDiskType()
│      ├─ ValidateMinRAM()
│      └─ ValidateMinCPU()
│
└── resource_validator.go      [300 linhas]
    └─ Validadores específicos
       ├─ ValidateCPU()
       ├─ ValidateRAM()
       ├─ ValidateDisk()
       └─ Constantes (MaxUtilization, MinReserved)
```

### Entrada/Saída
```go
// Input
type WorkloadRequest struct {
    Image         string
    Replicas      int
    CPUPerReplica float64
    RAMPerReplica float64
    Strategy      string
    Constraints   []Constraint
}

// Output
type AdmissionResult struct {
    Admitted      bool
    Reason        string
    GridCapacity  *CapacityInfo
    Recommendation string
}
```

---

## 🧠 SUBCOMPONENTE 2: SCHEDULER

### Responsabilidade
Decidir ONDE deployar cada réplica (placement decisions).

### Arquivos (6)

```
scheduler/
├── scheduler.go               [400 linhas]
│   └─ Schedule()
│      ├─ GetEligibleNodes()
│      ├─ SelectStrategy()
│      ├─ ApplyStrategy()
│      └─ LogDecisions()
│
├── node_filter.go             [300 linhas]
│   └─ Filtros de Nodes
│      ├─ FilterHealthyNodes()
│      ├─ FilterByConstraints()
│      ├─ FilterByCapacity()
│      └─ FilterByTags()
│
├── node_scorer.go             [300 linhas]
│   └─ Cálculo de scores
│      ├─ CalculateSpreadScore()
│      ├─ CalculateBinpackScore()
│      ├─ CalculateOptimizedScore()
│      └─ CalculateLoadBalanceScore()
│
├── strategy_spread.go         [300 linhas]
│   └─ SpreadStrategy
│      └─ Algoritmo: round-robin por menor carga
│
├── strategy_binpack.go        [300 linhas]
│   └─ BinpackStrategy
│      └─ Algoritmo: preencher Nodes até capacidade
│
└── strategy_optimized.go      [400 linhas]
    └─ OptimizedStrategy
       └─ Algoritmo: balancear CPU/RAM, sweet spot 40-60%
```

### Entrada/Saída
```go
// Input
WorkloadRequest + []EligibleNodes

// Output
type PlacementDecision struct {
    NodeID       string
    ReplicaIndex int
    Score        float64
    Reason       string
    Resources    *AllocatedResources
}
```

---

## 📥 SUBCOMPONENTE 3: QUEUE

### Responsabilidade
Gerenciar workloads que não podem ser deployados imediatamente.

### Arquivos (4)

```
queue/
├── queue_manager.go           [300 linhas]
│   └─ CRUD de fila
│      ├─ Enqueue()
│      ├─ Dequeue()
│      ├─ List()
│      ├─ Cancel()
│      └─ Persistência em ~/.syntropy/queue/
│
├── queue_processor.go         [300 linhas]
│   └─ Processamento periódico
│      ├─ Start() - inicia goroutine
│      ├─ ProcessOnce() - processa fila
│      ├─ TryAdmission() - tenta novamente
│      └─ AutoDeploy() - se aprovado
│
├── wait_estimator.go          [200 linhas]
│   └─ Estimativa de tempo
│      ├─ EstimateWaitTime()
│      ├─ CalculateAvgDuration()
│      └─ PredictAvailability()
│
└── priority_manager.go        [200 linhas]
    └─ Gerenciamento de prioridades
       ├─ SetPriority()
       ├─ ReorderQueue()
       └─ GetNextWorkload()
```

### Background Job
```go
// Goroutine rodando a cada 1 minuto
go func() {
    ticker := time.NewTicker(1 * time.Minute)
    for range ticker.C {
        processor.ProcessOnce()
    }
}()
```

---

## 🚀 SUBCOMPONENTE 4: DEPLOY

### Responsabilidade
Executar deployment nos Nodes (criar e iniciar containers).

### Arquivos (6)

```
deploy/
├── deployer.go                [400 linhas]
│   └─ Orquestração
│      ├─ Deploy()
│      │  ├─ Call Admission
│      │  ├─ Call Scheduler
│      │  ├─ ExecuteDeployments()
│      │  └─ SaveMetadata()
│      └─ DeployReplica()
│
├── executor.go                [400 linhas]
│   └─ Execução base
│      ├─ DeployReplica()
│      ├─ ValidateImage()
│      ├─ PullImage()
│      ├─ CreateContainer()
│      ├─ StartContainer()
│      └─ VerifyRunning()
│
├── executor_windows.go        [300 linhas]
│   └─ Implementação Windows
│      └─ Paths: C:\syntropy\
│
├── executor_linux.go          [300 linhas]
│   └─ Implementação Linux
│      └─ Paths: /opt/syntropy/
│
├── rollback.go                [300 linhas]
│   └─ Rollback Manager
│      ├─ Rollback()
│      ├─ StopContainer()
│      ├─ RemoveContainer()
│      └─ RestorePreviousState()
│
└── docker_client.go           [400 linhas]
    └─ Cliente Docker
       ├─ Pull()
       ├─ Create()
       ├─ Start()
       ├─ Stop()
       ├─ Remove()
       └─ Stats()
```

### Comandos Docker Executados
```bash
# Via SSH para cada Node
ssh syntropy@node-01 "docker pull nginx"
ssh syntropy@node-01 "docker create --name nginx-153045-1 --cpu 1 --memory 512M nginx"
ssh syntropy@node-01 "docker start nginx-153045-1"
ssh syntropy@node-01 "docker ps --filter id=nginx-153045-1"
```

---

## 🔄 SUBCOMPONENTE 5: LIFECYCLE

### Responsabilidade
Gerenciar lifecycle de workloads após deployment.

### Arquivos (5)

```
lifecycle/
├── lifecycle.go               [300 linhas]
│   └─ Manager principal
│      └─ Orquestra start/stop/restart/scale
│
├── start.go                   [200 linhas]
│   └─ StartWorkload()
│      ├─ LoadWorkload()
│      ├─ CheckStatus()
│      └─ docker start (para cada réplica)
│
├── stop.go                    [200 linhas]
│   └─ StopWorkload()
│      └─ docker stop (para cada réplica)
│
├── restart.go                 [200 linhas]
│   └─ RestartWorkload()
│      ├─ Stop()
│      ├─ Start()
│      └─ VerifyHealth()
│
└── scale.go                   [300 linhas]
    └─ ScaleWorkload()
       ├─ ScaleUp()
       │  ├─ CreateWorkloadRequest (delta réplicas)
       │  ├─ Call Admission Control
       │  ├─ Call Scheduler
       │  └─ Deploy novas réplicas
       └─ ScaleDown()
          ├─ SelectReplicasToRemove()
          └─ Stop + Remove containers
```

### Integração com Orquestração
```
Scale UP:
  ┌─> Admission Control (validar recursos para novas réplicas)
  ├─> Scheduler (decidir onde deployar)
  └─> Deploy (executar)

Scale DOWN:
  ├─> Selecionar réplicas
  ├─> Parar containers
  └─> Atualizar capacity (liberar recursos)
```

---

## 📊 SUBCOMPONENTE 6: MONITORING

### Responsabilidade
Coletar logs e métricas dos workloads.

### Arquivos (3)

```
monitoring/
├── monitoring.go              [300 linhas]
│   └─ Monitor principal
│      ├─ GetWorkloadStatus()
│      ├─ GetWorkloadLogs()
│      ├─ GetWorkloadMetrics()
│      └─ StreamLogs()
│
├── logs.go                    [400 linhas]
│   └─ Agregação de logs
│      ├─ GetLogs() - logs de 1 container
│      ├─ AggregateLogs() - logs de todas réplicas
│      ├─ StreamLogs() - streaming em tempo real
│      └─ FilterLogs() - por level, timestamp
│
└── metrics.go                 [400 linhas]
    └─ Coleta de métricas
       ├─ CollectWorkloadMetrics() - agregado
       ├─ CollectContainerMetrics() - por réplica
       └─ AggregateMetrics() - total do workload
```

### Coleta de Dados
```bash
# Para cada Node/Container via SSH:
docker stats <container> --no-stream --format "{{json .}}"

# Parse JSON:
{
  "CPUPerc": "15.34%",
  "MemUsage": "234.5MiB / 512MiB",
  "NetIO": "1.2kB / 2.5kB",
  "BlockIO": "0B / 0B"
}
```

---

## 🔧 COMANDOS CLI MAPEADOS

### Deploy
```bash
syntropy workload deploy <image> \
  --replicas <n> \
  --cpu <cores> \
  --memory <size> \
  [--strategy <spread|binpack|resource-optimized>] \
  [--node <node-id>] \
  [--tag <tag>] \
  [--disk-type <type>] \
  [--auto-adjust] \
  [--queue-if-full]

Chama:
  workload.Deploy() 
    → admission.Validate()
    → scheduler.Schedule()
    → deploy.Execute()
```

### List
```bash
syntropy workload list

Chama:
  workload.List()
    → Load de ~/.syntropy/workloads/*.yaml
    → Format output
```

### Status
```bash
syntropy workload status <workload-id>

Chama:
  workload.Status()
    → monitoring.GetWorkloadStatus()
    → SSH para cada Node/réplica
    → Agregar status
```

### Logs
```bash
syntropy workload logs <workload-id> [--follow]

Chama:
  workload.Logs()
    → monitoring.StreamLogs() (se --follow)
    → monitoring.GetLogs() (se não)
    → SSH: docker logs <container>
    → Agregar de todas réplicas
```

### Scale
```bash
syntropy workload scale <workload-id> --replicas <n>

Chama:
  workload.Scale()
    → lifecycle.Scale()
      → Se UP: admission.Validate() + scheduler.Schedule() + deploy.Execute()
      → Se DOWN: select réplicas + stop + remove
```

### Grid Capacity
```bash
syntropy grid capacity

Chama:
  admission.GetGridCapacity()
    → Listar Nodes
    → Calcular total/disponível
    → Format output
```

### Queue
```bash
syntropy workload queue list
syntropy workload queue status <queue-id>
syntropy workload queue cancel <queue-id>

Chama:
  workload.QueueList()
    → queue.List()
    → Load de ~/.syntropy/queue/*.yaml
```

---

## 📦 PERSISTÊNCIA DE DADOS

### ~/.syntropy/workloads/
```yaml
# nginx-153045.yaml
id: nginx-153045
image: nginx:latest
replicas: 3
resources:
  cpu_per_replica: 1.0
  ram_per_replica: 0.5
strategy: spread
status: running
deployed_at: "2025-10-10T15:30:00Z"
placements:
  - node_id: node-01
    replica_index: 1
    container_id: abc123
    status: running
  - node_id: node-02
    replica_index: 2
    container_id: def456
    status: running
  - node_id: node-03
    replica_index: 3
    container_id: ghi789
    status: running
```

### ~/.syntropy/queue/
```yaml
# queue-001.yaml
id: queue-001
workload_request:
  image: heavy-app
  replicas: 10
  cpu_per_replica: 4
  ram_per_replica: 8
queued_at: "2025-10-10T16:00:00Z"
priority: 5
estimated_wait: 15m
status: queued
attempts: 3
last_attempt: "2025-10-10T16:15:00Z"
```

---

## 🧪 TESTES POR SUBCOMPONENTE

### Admission Tests
```go
// admission/tests/admission_test.go

TestValidate_Success()
  - Workload válido
  - Grid tem recursos
  - Deve: ACEITAR

TestValidate_InsufficientCPU()
  - Workload precisa 50 cores
  - Grid tem 48 cores
  - Deve: REJEITAR com sugestão

TestValidate_GridOverload()
  - Utilização projetada: 95%
  - Limite: 90%
  - Deve: REJEITAR

TestValidate_Constraints()
  - Constraint: tag=gpu
  - Apenas 2 Nodes têm GPU
  - Replicas: 5
  - Deve: REJEITAR (não há Nodes suficientes)
```

### Scheduler Tests
```go
// scheduler/tests/spread_test.go

TestSpreadStrategy()
  - 3 réplicas
  - 6 Nodes (todos vazios)
  - Deve: distribuir em node-01, node-02, node-03

TestSpreadStrategy_UnbalancedLoad()
  - 3 réplicas
  - Nodes: [0, 0, 2, 3, 5, 7] workloads
  - Deve: preferir node-01, node-02, node-03 (menor carga)

// scheduler/tests/binpack_test.go

TestBinpackStrategy()
  - 3 réplicas
  - Nodes: [30%, 45%, 60%, 20%, 10%, 5%] utilização
  - Deve: preferir node-03 (60%), depois node-02 (45%)

TestBinpackStrategy_FillUntilFull()
  - 10 réplicas (1 CPU cada)
  - Node com 8 cores
  - Deve: colocar 8 no primeiro Node, 2 no segundo
```

### Queue Tests
```go
// queue/tests/queue_test.go

TestEnqueue()
  - Workload rejeitado
  - Deve: adicionar à fila, gerar ID, estimar espera

TestProcessQueue()
  - 3 workloads na fila
  - Grid libera recursos
  - Deve: deployar automaticamente primeiro da fila

TestCancelQueue()
  - Workload na fila
  - Usuário cancela
  - Deve: remover da fila, não deployar
```

### Deploy Tests
```go
// deploy/tests/deployer_test.go

TestDeploy_Success()
  - Admission: ✅
  - Scheduler: 3 placements
  - Deve: deployar 3 réplicas, salvar metadados

TestDeploy_FailureWithRollback()
  - Deploy réplica 1: ✅
  - Deploy réplica 2: ❌ (falha)
  - Deve: rollback réplica 1, retornar erro

TestDeploy_NodeOffline()
  - Scheduler escolhe node-01
  - node-01 offline durante deploy
  - Deve: falhar, rollback
```

---

## 📊 MATRIZ DE DEPENDÊNCIAS

```
workload.go (Orquestrador)
    ↓
    ├─> admission/ ─────────┐
    │   ├─> inventory       │
    │   └─> sync            │
    │                        ↓
    ├─> scheduler/ ────────> deploy/
    │   ├─> inventory          ├─> docker_client
    │   └─> admission          ├─> ssh_client
    │                          └─> rollback
    ├─> queue/
    │   ├─> admission
    │   └─> deployer
    │
    ├─> lifecycle/
    │   ├─> executor
    │   └─> inventory
    │
    └─> monitoring/
        ├─> ssh_client
        └─> inventory
```

---

## ⏱️ ESTIMATIVA DE TEMPO (REALISTA)

### Semana 5: Workload Component (7 dias)

```
Dia 1: Admission Control (Parte 1)
  - admission_controller.go
  - capacity_calculator.go
  - Testes básicos
  (8 horas)

Dia 2: Admission Control (Parte 2)
  - constraint_validator.go
  - resource_validator.go
  - Testes completos
  (8 horas)

Dia 3: Scheduler (Parte 1)
  - scheduler.go
  - node_filter.go
  - strategy_spread.go
  (8 horas)

Dia 4: Scheduler (Parte 2)
  - strategy_binpack.go
  - strategy_optimized.go
  - node_scorer.go
  - Testes completos
  (8 horas)

Dia 5: Queue System
  - queue_manager.go
  - queue_processor.go
  - wait_estimator.go
  - priority_manager.go
  - Testes
  (8 horas)

Dia 6: Deploy Execution
  - deployer.go
  - executor.go + Windows/Linux
  - docker_client.go
  - rollback.go
  - Testes
  (8 horas)

Dia 7: Integração + Lifecycle + Monitoring
  - workload.go (orquestrador)
  - lifecycle/* (5 arquivos)
  - monitoring/* (3 arquivos)
  - Testes de integração
  (8 horas)

Total: 56 horas (7 dias × 8h)
```

---

## ✅ CHECKLIST DE VALIDAÇÃO

### Admission Control
```
[ ] Valida CPU disponível corretamente
[ ] Valida RAM disponível corretamente
[ ] Detecta sobrecarga (>90%)
[ ] Sugere ajustes quando rejeitado
[ ] Valida constraints (tags, disk-type, min-ram)
[ ] GetGridCapacity() retorna dados corretos
[ ] Testes passam (4 cenários mínimos)
```

### Scheduler
```
[ ] Spread distribui uniformemente
[ ] Binpack preenche Nodes
[ ] Resource-optimized balanceia CPU/RAM
[ ] Filtra Nodes healthy
[ ] Filtra por constraints
[ ] Calcula scores corretamente
[ ] Testes passam (3 estratégias)
```

### Queue
```
[ ] Enfileira workloads rejeitados
[ ] Persiste em ~/.syntropy/queue/
[ ] Processa periodicamente (1 min)
[ ] Estima tempo de espera
[ ] Deploy automático quando aprovado
[ ] Permite cancelamento
[ ] Testes passam (4 cenários)
```

### Deploy
```
[ ] Executa deployment via SSH
[ ] Pull de imagem funciona
[ ] Cria container com recursos corretos
[ ] Inicia container
[ ] Verifica se está running
[ ] Rollback em falha funciona
[ ] Multi-plataforma (Windows/Linux)
[ ] Testes passam (4 cenários)
```

### Lifecycle
```
[ ] Start workload funciona
[ ] Stop workload funciona
[ ] Restart workload funciona
[ ] Scale up passa por Admission
[ ] Scale down libera recursos
[ ] Testes passam
```

### Monitoring
```
[ ] Coleta logs via docker logs
[ ] Agrega logs de múltiplas réplicas
[ ] Streaming funciona (--follow)
[ ] Coleta métricas via docker stats
[ ] Agrega métricas
[ ] Testes passam
```

### Integração
```
[ ] workload.Deploy() orquestra tudo
[ ] Background jobs iniciam
[ ] Queue processor funciona
[ ] Todos os comandos CLI funcionam:
    [ ] syntropy workload deploy
    [ ] syntropy workload list
    [ ] syntropy workload status
    [ ] syntropy workload logs
    [ ] syntropy workload scale
    [ ] syntropy grid capacity
    [ ] syntropy workload queue list
```

---

## 🎯 CRITÉRIOS DE SUCESSO

### MVP do Componente Workload está completo quando:

```
✅ Deploy grid-wide funciona
   $ syntropy workload deploy nginx --replicas 3
   → Admission valida
   → Scheduler escolhe 3 Nodes automaticamente
   → Deploy em 3 Nodes
   → 3/3 réplicas running

✅ Validação de capacidade funciona
   $ syntropy workload deploy huge-app --replicas 20
   → ❌ Rejeitado (insuficiente recursos)
   → Mostra reason + recommendation

✅ Estratégias de scheduling funcionam
   $ syntropy workload deploy app --strategy binpack --replicas 5
   → Scheduler usa binpack
   → Preenche Nodes sequencialmente

✅ Queue system funciona
   $ syntropy workload deploy job --queue-if-full
   → Rejeitado (Grid cheia)
   → Enfileirado automaticamente
   → Deploy automático quando liberar

✅ Constraints funcionam
   $ syntropy workload deploy ml --tag gpu --replicas 2
   → Filtra apenas Nodes com GPU
   → Deploy apenas em Nodes elegíveis

✅ Rollback funciona
   → Deploy falha na réplica 2
   → Rollback automático da réplica 1
   → Erro reportado ao usuário

✅ Lifecycle funciona
   $ syntropy workload scale nginx-001 --replicas 6
   → Scale de 3 → 6
   → Admission valida 3 novas réplicas
   → Deploy automático

✅ Monitoring funciona
   $ syntropy workload logs nginx-001 --follow
   → Logs agregados de 3 réplicas
   → Streaming em tempo real

✅ Grid capacity funciona
   $ syntropy grid capacity
   → Mostra recursos totais/disponíveis
   → Mostra utilização por Node
```

---

**Especificação completa do Componente Workload**  
**Integrado ao**: MVP.md seção 6-7  
**Ver também**: ORCHESTRATION-SUMMARY.md


