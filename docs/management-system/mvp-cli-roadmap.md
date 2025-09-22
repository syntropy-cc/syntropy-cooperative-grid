# 🎯 MVP CLI Roadmap - Management System

> **Roadmap Técnico Detalhado para o MVP CLI do Syntropy Cooperative Grid Management System**

## 📋 **Índice**

1. [Visão Geral do MVP](#visão-geral-do-mvp)
2. [Objetivos e Escopo](#objetivos-e-escopo)
3. [Sprint 1: Foundation & Setup](#sprint-1-foundation--setup)
4. [Sprint 2: USB Detection & Node Creation](#sprint-2-usb-detection--node-creation)
5. [Sprint 3: Node Management](#sprint-3-node-management)
6. [Sprint 4: Container Basics](#sprint-4-container-basics)
7. [Critérios de Aceitação](#critérios-de-aceitação)
8. [Métricas de Sucesso](#métricas-de-sucesso)
9. [Riscos e Mitigações](#riscos-e-mitigações)
10. [Próximos Passos](#próximos-passos)

---

## 🎯 **Visão Geral do MVP**

### **Definição do MVP**
O **MVP CLI** é a primeira versão funcional do Management System, focada em fornecer uma interface de linha de comando robusta e intuitiva para gerenciar nós e containers básicos da Syntropy Cooperative Grid.

### **Objetivo Principal**
Criar uma ferramenta CLI que permita aos usuários:
- **Detectar e configurar** dispositivos USB automaticamente
- **Criar e gerenciar** nós da grid
- **Fazer deploy** de containers básicos
- **Monitorar** status e logs
- **Automatizar** operações repetitivas

### **Duração**
- **4 Sprints** de 2 semanas cada
- **8 semanas** total
- **Foco**: Funcionalidade core, não perfeição

### **Público-Alvo**
- **Administradores de Sistema**: Gerenciamento de infraestrutura
- **Desenvolvedores**: Automação e scripts
- **Operadores**: Troubleshooting e manutenção
- **Usuários Técnicos**: Operações básicas da grid

---

## 🎯 **Objetivos e Escopo**

### **Objetivos do MVP**

#### **1. Funcionalidade Core**
- ✅ Detecção automática de USB cross-platform
- ✅ Criação e configuração de nós
- ✅ Gerenciamento básico de containers
- ✅ Interface CLI intuitiva e consistente

#### **2. Qualidade**
- ✅ Código limpo e bem documentado
- ✅ Testes unitários básicos
- ✅ Tratamento de erros robusto
- ✅ Logs estruturados

#### **3. Usabilidade**
- ✅ Comandos intuitivos e consistentes
- ✅ Help contextual e documentação
- ✅ Output formatável (table, json, yaml)
- ✅ Progress indicators para operações longas

#### **4. Confiabilidade**
- ✅ Funciona em Windows, Linux e WSL
- ✅ Tratamento de permissões adequado
- ✅ Recuperação de erros
- ✅ Validação de inputs

### **Escopo do MVP**

#### **✅ Incluído**
- Detecção de USB devices
- Criação de nós via USB
- Listagem e status de nós
- Deploy básico de containers
- Logs de containers
- Configuração básica
- Help e documentação

#### **❌ Não Incluído (Futuras Versões)**
- Interface web
- Gerenciamento avançado de rede
- Serviços cooperativos
- Monitoramento avançado
- APIs REST/GraphQL
- Autenticação complexa
- Multi-tenant support

---

## 🚀 **Sprint 1: Foundation & Setup (2 semanas)**

### **Objetivos do Sprint**
- Configurar base sólida do projeto
- Implementar sistema de build e CI/CD
- Criar documentação de setup
- Estabelecer padrões de desenvolvimento

### **User Stories**

#### **US-001: Estrutura de Projeto**
**Como** desenvolvedor  
**Quero** ter uma estrutura de projeto bem definida  
**Para que** eu possa organizar o código de forma clara e escalável

**Critérios de Aceitação:**
- [ ] Estrutura de diretórios seguindo convenções Go
- [ ] Separação clara entre core e interfaces
- [ ] Arquivos de configuração organizados
- [ ] Documentação da estrutura

#### **US-002: Sistema de Build**
**Como** desenvolvedor  
**Quero** ter um sistema de build automatizado  
**Para que** eu possa compilar e testar o projeto facilmente

**Critérios de Aceitação:**
- [ ] Makefile com comandos essenciais
- [ ] Build cross-platform (Windows, Linux, macOS)
- [ ] Testes automatizados
- [ ] Linting e formatação automática

#### **US-003: CI/CD Pipeline**
**Como** desenvolvedor  
**Quero** ter CI/CD básico funcionando  
**Para que** eu possa garantir qualidade do código

**Critérios de Aceitação:**
- [ ] GitHub Actions configurado
- [ ] Testes executados em PRs
- [ ] Build executado em PRs
- [ ] Linting executado em PRs

### **Tarefas Técnicas Detalhadas**

#### **1. Estrutura de Projeto**
```bash
# Criar estrutura de diretórios
mkdir -p {core,interfaces/cli/{cmd,internal/cli,pkg},deployments/{docker,kubernetes},docs,scripts}

# Configurar go.mod
go mod init github.com/syntropy-cc/cooperative-grid
go mod tidy
```

#### **2. Sistema de Build**
```makefile
# Makefile principal
.PHONY: build test lint clean

build:
	go build -o bin/syntropy-cli interfaces/cli/cmd/main.go

test:
	go test ./...

lint:
	golangci-lint run

clean:
	rm -rf bin/
```

#### **3. CI/CD Pipeline**
```yaml
# .github/workflows/ci.yml
name: CI
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - name: Run tests
      run: make test
    - name: Run linting
      run: make lint
```

#### **4. Docker para Desenvolvimento**
```dockerfile
# Dockerfile para desenvolvimento
FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go mod download
CMD ["go", "run", "interfaces/cli/cmd/main.go"]
```

### **Entregáveis**
- [ ] Estrutura de projeto completa
- [ ] Sistema de build funcionando
- [ ] CI/CD pipeline básico
- [ ] Documentação de setup
- [ ] Ambiente de desenvolvimento configurado

### **Métricas de Sucesso**
- ✅ Projeto compila sem erros
- ✅ Docker Compose sobe todos os serviços
- ✅ CI passa em todos os PRs
- ✅ Documentação de setup completa

---

## 🔍 **Sprint 2: USB Detection & Node Creation (2 semanas)**

### **Objetivos do Sprint**
- Implementar detecção de USB cross-platform
- Criar sistema de criação de nós
- Implementar formatação de USB
- Resolver problemas de permissões

### **User Stories**

#### **US-004: Detecção de USB**
**Como** usuário  
**Quero** listar dispositivos USB disponíveis  
**Para que** eu possa escolher qual usar para criar um nó

**Critérios de Aceitação:**
- [ ] Lista USB devices em Windows (PowerShell)
- [ ] Lista USB devices em Linux (lsblk)
- [ ] Lista USB devices em WSL (hybrid approach)
- [ ] Mostra informações relevantes (tamanho, tipo, status)
- [ ] Filtra apenas dispositivos removíveis

#### **US-005: Criação de Nó**
**Como** usuário  
**Quero** criar um nó a partir de um USB  
**Para que** eu possa adicionar um novo nó à grid

**Critérios de Aceitação:**
- [ ] Aceita parâmetros: --usb, --name, --auto-detect
- [ ] Valida nome do nó (único, formato válido)
- [ ] Formata USB automaticamente
- [ ] Gera chaves SSH
- [ ] Cria configuração do nó
- [ ] Mostra progresso da operação

#### **US-006: Progresso e Logs**
**Como** usuário  
**Quero** ver o progresso da criação do nó  
**Para que** eu possa acompanhar a operação

**Critérios de Aceitação:**
- [ ] Progress bar para operações longas
- [ ] Logs estruturados com níveis
- [ ] Mensagens de erro claras
- [ ] Possibilidade de cancelar operação

### **Tarefas Técnicas Detalhadas**

#### **1. Detecção de USB Cross-Platform**

##### **Windows (PowerShell)**
```go
// internal/platform/windows/usb_detector.go
func DetectUSBDevices() ([]USBDevice, error) {
    cmd := exec.Command("powershell", "-Command", `
        Get-WmiObject -Class Win32_DiskDrive | 
        Where-Object { $_.MediaType -eq "Removable Media" } |
        Select-Object DeviceID, Size, Model |
        ConvertTo-Json
    `)
    
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    var devices []USBDevice
    err = json.Unmarshal(output, &devices)
    return devices, err
}
```

##### **Linux (lsblk)**
```go
// internal/platform/linux/usb_detector.go
func DetectUSBDevices() ([]USBDevice, error) {
    cmd := exec.Command("lsblk", "-J", "-o", "NAME,SIZE,TYPE,MOUNTPOINT")
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    var lsblkOutput LSBlkOutput
    err = json.Unmarshal(output, &lsblkOutput)
    if err != nil {
        return nil, err
    }
    
    return filterRemovableDevices(lsblkOutput.BlockDevices), nil
}
```

##### **WSL (Hybrid)**
```go
// internal/platform/wsl/usb_detector.go
func DetectUSBDevices() ([]USBDevice, error) {
    // Primeiro tenta via Windows PowerShell
    windowsDevices, err := detectViaWindows()
    if err == nil && len(windowsDevices) > 0 {
        return convertToWSLDevices(windowsDevices)
    }
    
    // Fallback para detecção Linux
    return detectViaLinux()
}
```

#### **2. Formatação de USB**

##### **Windows (diskpart)**
```go
// internal/platform/windows/usb_formatter.go
func FormatUSBDevice(deviceID string) error {
    script := fmt.Sprintf(`
        select disk %s
        clean
        create partition primary
        active
        format fs=fat32 quick
        assign
    `, extractDiskNumber(deviceID))
    
    return executeDiskpartScript(script)
}
```

##### **Linux (mkfs)**
```go
// internal/platform/linux/usb_formatter.go
func FormatUSBDevice(devicePath string) error {
    // Unmount se estiver montado
    unmountDevice(devicePath)
    
    // Criar nova tabela de partições
    cmd := exec.Command("parted", devicePath, "mklabel", "msdos")
    if err := cmd.Run(); err != nil {
        return err
    }
    
    // Criar partição
    cmd = exec.Command("parted", devicePath, "mkpart", "primary", "fat32", "1MiB", "100%")
    if err := cmd.Run(); err != nil {
        return err
    }
    
    // Formatar
    partitionPath := devicePath + "1"
    cmd = exec.Command("mkfs.fat", "-F32", partitionPath)
    return cmd.Run()
}
```

#### **3. Geração de Chaves SSH**
```go
// internal/crypto/ssh_keygen.go
func GenerateSSHKeys(nodeName string) (*SSHKeyPair, error) {
    privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        return nil, err
    }
    
    publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
    if err != nil {
        return nil, err
    }
    
    return &SSHKeyPair{
        PrivateKey: privateKey,
        PublicKey:  publicKey,
        NodeName:   nodeName,
    }, nil
}
```

#### **4. Progress Bar e Logs**
```go
// internal/ui/progress.go
type ProgressBar struct {
    bar    *pb.ProgressBar
    logger *logrus.Logger
}

func (p *ProgressBar) Update(step string, progress int) {
    p.bar.SetCurrent(int64(progress))
    p.bar.Set("step", step)
    p.logger.WithField("step", step).Info("Progress update")
}
```

### **Entregáveis**
- [ ] Comando `syntropy-cli node create` funcional
- [ ] Detecção de USB cross-platform
- [ ] Formatação de USB com tratamento de permissões
- [ ] Geração de chaves SSH
- [ ] Progress bar e logs estruturados

### **Métricas de Sucesso**
- ✅ Detecta USB em Windows, Linux e WSL
- ✅ Formata USB com sucesso
- ✅ Cria nó funcional
- ✅ Interface CLI intuitiva

---

## 🖥️ **Sprint 3: Node Management (2 semanas)**

### **Objetivos do Sprint**
- Implementar CRUD completo de nós
- Criar sistema de persistência
- Implementar validações
- Adicionar monitoramento básico

### **User Stories**

#### **US-007: Listagem de Nós**
**Como** usuário  
**Quero** listar todos os meus nós  
**Para que** eu possa ver o status da minha grid

**Critérios de Aceitação:**
- [ ] Lista todos os nós com informações básicas
- [ ] Suporte a filtros (status, nome, data)
- [ ] Output formatável (table, json, yaml)
- [ ] Paginação para muitos nós
- [ ] Ordenação por diferentes campos

#### **US-008: Status de Nó**
**Como** usuário  
**Quero** ver status detalhado de um nó  
**Para que** eu possa diagnosticar problemas

**Critérios de Aceitação:**
- [ ] Mostra informações detalhadas do nó
- [ ] Inclui métricas de performance
- [ ] Mostra histórico de eventos
- [ ] Suporte a --watch para updates em tempo real
- [ ] Formato de saída configurável

#### **US-009: Atualização de Nó**
**Como** usuário  
**Quero** atualizar configuração de um nó  
**Para que** eu possa ajustar parâmetros

**Critérios de Aceitação:**
- [ ] Permite atualizar nome e descrição
- [ ] Permite atualizar configurações de rede
- [ ] Valida configurações antes de aplicar
- [ ] Aplica mudanças sem reiniciar nó
- [ ] Log de auditoria das mudanças

#### **US-010: Remoção de Nó**
**Como** usuário  
**Quero** remover um nó  
**Para que** eu possa limpar a grid

**Critérios de Aceitação:**
- [ ] Confirmação antes de remover
- [ ] Remove configurações e dados
- [ ] Limpa recursos associados
- [ ] Log de auditoria
- [ ] Opção --force para pular confirmação

### **Tarefas Técnicas Detalhadas**

#### **1. Sistema de Persistência**

##### **SQLite para MVP**
```go
// internal/storage/sqlite/node_repository.go
type NodeRepository struct {
    db *sql.DB
}

func (r *NodeRepository) Create(node *models.Node) error {
    query := `
        INSERT INTO nodes (id, name, description, status, hardware_info, network_config, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `
    
    _, err := r.db.Exec(query, 
        node.ID, node.Name, node.Description, node.Status,
        node.HardwareInfo, node.NetworkConfig, node.CreatedAt, node.UpdatedAt)
    
    return err
}

func (r *NodeRepository) GetByID(id string) (*models.Node, error) {
    query := `SELECT * FROM nodes WHERE id = ?`
    row := r.db.QueryRow(query, id)
    
    node := &models.Node{}
    err := row.Scan(&node.ID, &node.Name, &node.Description, &node.Status,
                   &node.HardwareInfo, &node.NetworkConfig, &node.CreatedAt, &node.UpdatedAt)
    
    return node, err
}
```

##### **Schema do Banco**
```sql
-- internal/storage/sqlite/schema.sql
CREATE TABLE nodes (
    id TEXT PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    description TEXT,
    status TEXT NOT NULL DEFAULT 'creating',
    hardware_info TEXT, -- JSON
    network_config TEXT, -- JSON
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_nodes_status ON nodes(status);
CREATE INDEX idx_nodes_name ON nodes(name);
CREATE INDEX idx_nodes_created_at ON nodes(created_at);
```

#### **2. Validações**

##### **Validação de Nome**
```go
// internal/validation/node_validator.go
func ValidateNodeName(name string) error {
    if len(name) < 3 {
        return errors.New("node name must be at least 3 characters")
    }
    
    if len(name) > 50 {
        return errors.New("node name must be less than 50 characters")
    }
    
    // Apenas alfanumérico, hífen e underscore
    matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, name)
    if !matched {
        return errors.New("node name can only contain alphanumeric characters, hyphens, and underscores")
    }
    
    return nil
}
```

##### **Validação de Configuração**
```go
// internal/validation/config_validator.go
func ValidateNetworkConfig(config *models.NetworkConfig) error {
    if config.IPAddress != "" {
        if net.ParseIP(config.IPAddress) == nil {
            return errors.New("invalid IP address")
        }
    }
    
    if config.Port < 1 || config.Port > 65535 {
        return errors.New("port must be between 1 and 65535")
    }
    
    return nil
}
```

#### **3. Monitoramento Básico**

##### **Health Check**
```go
// internal/monitoring/health_checker.go
type HealthChecker struct {
    nodeRepo NodeRepository
}

func (h *HealthChecker) CheckNodeHealth(nodeID string) (*HealthStatus, error) {
    node, err := h.nodeRepo.GetByID(nodeID)
    if err != nil {
        return nil, err
    }
    
    status := &HealthStatus{
        NodeID:    nodeID,
        Timestamp: time.Now(),
        Status:    "unknown",
    }
    
    // Verificar conectividade SSH
    if h.checkSSHConnectivity(node) {
        status.Status = "healthy"
    } else {
        status.Status = "unhealthy"
    }
    
    // Verificar recursos
    resources, err := h.checkResources(node)
    if err == nil {
        status.Resources = resources
    }
    
    return status, nil
}
```

#### **4. Comandos CLI**

##### **Listar Nós**
```go
// interfaces/cli/internal/cli/node_list.go
func newNodeListCommand() *cobra.Command {
    var format, filter string
    
    cmd := &cobra.Command{
        Use:   "list",
        Short: "List all nodes",
        RunE: func(cmd *cobra.Command, args []string) error {
            nodes, err := nodeService.ListNodes(cmd.Context(), &node.Filter{
                Status: filter,
            })
            if err != nil {
                return err
            }
            
            return outputNodes(nodes, format)
        },
    }
    
    cmd.Flags().StringVarP(&format, "format", "f", "table", "Output format")
    cmd.Flags().StringVar(&filter, "filter", "", "Filter by status")
    
    return cmd
}
```

##### **Status de Nó**
```go
// interfaces/cli/internal/cli/node_status.go
func newNodeStatusCommand() *cobra.Command {
    var format string
    var watch bool
    
    cmd := &cobra.Command{
        Use:   "status [node-id]",
        Short: "Show node status",
        Args:  cobra.MaximumNArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            if len(args) == 0 {
                return showAllNodesStatus(cmd.Context(), format, watch)
            }
            
            return showNodeStatus(cmd.Context(), args[0], format, watch)
        },
    }
    
    cmd.Flags().StringVarP(&format, "format", "f", "table", "Output format")
    cmd.Flags().BoolVarP(&watch, "watch", "w", false, "Watch for changes")
    
    return cmd
}
```

### **Entregáveis**
- [ ] CRUD completo de nós
- [ ] Interface intuitiva
- [ ] Dados persistidos corretamente
- [ ] Validações funcionando
- [ ] Monitoramento básico

### **Métricas de Sucesso**
- ✅ CRUD completo de nós
- ✅ Interface intuitiva
- ✅ Dados persistidos corretamente
- ✅ Validações funcionando

---

## 🐳 **Sprint 4: Container Basics (2 semanas)**

### **Objetivos do Sprint**
- Implementar funcionalidades básicas de container
- Integrar com Docker API
- Criar templates de containers
- Implementar logs e monitoramento

### **User Stories**

#### **US-011: Listagem de Containers**
**Como** usuário  
**Quero** listar containers em um nó  
**Para que** eu possa ver o que está rodando

**Critérios de Aceitação:**
- [ ] Lista containers por nó
- [ ] Mostra informações básicas (nome, imagem, status)
- [ ] Suporte a filtros (status, imagem)
- [ ] Output formatável
- [ ] Ordenação por diferentes campos

#### **US-012: Deploy de Container**
**Como** usuário  
**Quero** fazer deploy de um container  
**Para que** eu possa executar aplicações na grid

**Critérios de Aceitação:**
- [ ] Aceita parâmetros: --image, --node, --name, --ports, --env
- [ ] Valida imagem Docker
- [ ] Configura portas e variáveis de ambiente
- [ ] Deploy no nó especificado
- [ ] Verifica se deploy foi bem-sucedido

#### **US-013: Logs de Container**
**Como** usuário  
**Quero** ver logs de um container  
**Para que** eu possa diagnosticar problemas

**Critérios de Aceitação:**
- [ ] Mostra logs em tempo real
- [ ] Suporte a --follow para acompanhar
- [ ] Suporte a --tail para mostrar últimas linhas
- [ ] Filtros por nível de log
- [ ] Formatação colorida

#### **US-014: Controle de Container**
**Como** usuário  
**Quero** parar/iniciar containers  
**Para que** eu possa gerenciar aplicações

**Critérios de Aceitação:**
- [ ] Comando start para iniciar container
- [ ] Comando stop para parar container
- [ ] Comando restart para reiniciar
- [ ] Comando remove para remover
- [ ] Confirmação para operações destrutivas

### **Tarefas Técnicas Detalhadas**

#### **1. Integração com Docker API**

##### **Docker Client**
```go
// internal/container/docker_client.go
type DockerClient struct {
    client *docker.Client
}

func NewDockerClient() (*DockerClient, error) {
    client, err := docker.NewClientWithOpts(docker.FromEnv)
    if err != nil {
        return nil, err
    }
    
    return &DockerClient{client: client}, nil
}

func (d *DockerClient) ListContainers(nodeID string) ([]Container, error) {
    containers, err := d.client.ContainerList(context.Background(), types.ContainerListOptions{
        All: true,
    })
    if err != nil {
        return nil, err
    }
    
    var result []Container
    for _, container := range containers {
        if d.isContainerOnNode(container, nodeID) {
            result = append(result, convertToContainer(container))
        }
    }
    
    return result, nil
}
```

##### **Deploy de Container**
```go
// internal/container/deployer.go
type ContainerDeployer struct {
    dockerClient *DockerClient
    nodeRepo     NodeRepository
}

func (d *ContainerDeployer) Deploy(req *DeployRequest) (*Container, error) {
    // Validar nó
    node, err := d.nodeRepo.GetByID(req.NodeID)
    if err != nil {
        return nil, fmt.Errorf("node not found: %w", err)
    }
    
    // Configurar container
    config := &container.Config{
        Image: req.Image,
        Env:   req.EnvVars,
    }
    
    hostConfig := &container.HostConfig{
        PortBindings: d.buildPortBindings(req.Ports),
    }
    
    // Criar container
    resp, err := d.dockerClient.client.ContainerCreate(
        context.Background(),
        config,
        hostConfig,
        nil,
        nil,
        req.Name,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to create container: %w", err)
    }
    
    // Iniciar container
    err = d.dockerClient.client.ContainerStart(
        context.Background(),
        resp.ID,
        types.ContainerStartOptions{},
    )
    if err != nil {
        return nil, fmt.Errorf("failed to start container: %w", err)
    }
    
    return d.getContainerInfo(resp.ID)
}
```

#### **2. Templates de Containers**

##### **Template System**
```go
// internal/container/templates.go
type ContainerTemplate struct {
    Name        string            `json:"name"`
    Description string            `json:"description"`
    Image       string            `json:"image"`
    Ports       []PortMapping     `json:"ports"`
    EnvVars     map[string]string `json:"env_vars"`
    Volumes     []VolumeMapping   `json:"volumes"`
}

var DefaultTemplates = map[string]ContainerTemplate{
    "nginx": {
        Name:        "Nginx Web Server",
        Description: "Lightweight web server",
        Image:       "nginx:latest",
        Ports: []PortMapping{
            {HostPort: 80, ContainerPort: 80},
        },
    },
    "postgres": {
        Name:        "PostgreSQL Database",
        Description: "Relational database",
        Image:       "postgres:15",
        Ports: []PortMapping{
            {HostPort: 5432, ContainerPort: 5432},
        },
        EnvVars: map[string]string{
            "POSTGRES_DB":       "syntropy",
            "POSTGRES_USER":     "syntropy",
            "POSTGRES_PASSWORD": "syntropy",
        },
    },
}
```

##### **Template Deploy**
```go
// interfaces/cli/internal/cli/container_deploy.go
func newContainerDeployCommand() *cobra.Command {
    var image, nodeID, name string
    var ports, envVars []string
    
    cmd := &cobra.Command{
        Use:   "deploy",
        Short: "Deploy a container",
        RunE: func(cmd *cobra.Command, args []string) error {
            // Se não especificou imagem, mostrar templates
            if image == "" {
                return showAvailableTemplates()
            }
            
            req := &container.DeployRequest{
                Image:   image,
                NodeID:  nodeID,
                Name:    name,
                Ports:   parsePorts(ports),
                EnvVars: parseEnvVars(envVars),
            }
            
            deployedContainer, err := containerService.Deploy(cmd.Context(), req)
            if err != nil {
                return err
            }
            
            fmt.Printf("✅ Container deployed successfully!\n")
            fmt.Printf("   ID: %s\n", deployedContainer.ID)
            fmt.Printf("   Name: %s\n", deployedContainer.Name)
            fmt.Printf("   Status: %s\n", deployedContainer.Status)
            
            return nil
        },
    }
    
    cmd.Flags().StringVarP(&image, "image", "i", "", "Container image")
    cmd.Flags().StringVarP(&nodeID, "node", "n", "", "Target node ID")
    cmd.Flags().StringVar(&name, "name", "", "Container name")
    cmd.Flags().StringSliceVarP(&ports, "port", "p", []string{}, "Port mappings")
    cmd.Flags().StringSliceVarP(&envVars, "env", "e", []string{}, "Environment variables")
    
    return cmd
}
```

#### **3. Logs e Monitoramento**

##### **Log Viewer**
```go
// internal/container/log_viewer.go
type LogViewer struct {
    dockerClient *DockerClient
}

func (l *LogViewer) GetLogs(containerID string, options *LogOptions) (io.ReadCloser, error) {
    return l.dockerClient.client.ContainerLogs(context.Background(), containerID, types.ContainerLogsOptions{
        ShowStdout: true,
        ShowStderr: true,
        Follow:     options.Follow,
        Tail:       options.Tail,
        Since:      options.Since,
    })
}

func (l *LogViewer) StreamLogs(containerID string, options *LogOptions, writer io.Writer) error {
    logs, err := l.GetLogs(containerID, options)
    if err != nil {
        return err
    }
    defer logs.Close()
    
    scanner := bufio.NewScanner(logs)
    for scanner.Scan() {
        line := scanner.Text()
        // Remover header do Docker (8 bytes)
        if len(line) > 8 {
            line = line[8:]
        }
        fmt.Fprintln(writer, line)
    }
    
    return scanner.Err()
}
```

##### **Comando de Logs**
```go
// interfaces/cli/internal/cli/container_logs.go
func newContainerLogsCommand() *cobra.Command {
    var follow bool
    var tail int
    
    cmd := &cobra.Command{
        Use:   "logs <container-id>",
        Short: "Show container logs",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            containerID := args[0]
            
            options := &container.LogOptions{
                Follow: follow,
                Tail:   tail,
            }
            
            return containerService.StreamLogs(cmd.Context(), containerID, options, os.Stdout)
        },
    }
    
    cmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow log output")
    cmd.Flags().IntVar(&tail, "tail", 100, "Number of lines to show")
    
    return cmd
}
```

#### **4. Controle de Containers**

##### **Container Controller**
```go
// internal/container/controller.go
type ContainerController struct {
    dockerClient *DockerClient
}

func (c *ContainerController) Start(containerID string) error {
    return c.dockerClient.client.ContainerStart(
        context.Background(),
        containerID,
        types.ContainerStartOptions{},
    )
}

func (c *ContainerController) Stop(containerID string, timeout *int) error {
    timeoutDuration := 30 * time.Second
    if timeout != nil {
        timeoutDuration = time.Duration(*timeout) * time.Second
    }
    
    return c.dockerClient.client.ContainerStop(
        context.Background(),
        containerID,
        &timeoutDuration,
    )
}

func (c *ContainerController) Restart(containerID string, timeout *int) error {
    timeoutDuration := 30 * time.Second
    if timeout != nil {
        timeoutDuration = time.Duration(*timeout) * time.Second
    }
    
    return c.dockerClient.client.ContainerRestart(
        context.Background(),
        containerID,
        &timeoutDuration,
    )
}

func (c *ContainerController) Remove(containerID string, force bool) error {
    return c.dockerClient.client.ContainerRemove(
        context.Background(),
        containerID,
        types.ContainerRemoveOptions{
            Force: force,
        },
    )
}
```

### **Entregáveis**
- [ ] Deploy de containers funciona
- [ ] Gerenciamento básico de containers
- [ ] Logs acessíveis
- [ ] Templates úteis
- [ ] Controle de containers

### **Métricas de Sucesso**
- ✅ Deploy de containers funciona
- ✅ Gerenciamento básico de containers
- ✅ Logs acessíveis
- ✅ Templates úteis

---

## ✅ **Critérios de Aceitação do MVP**

### **Funcionalidade Core**
- [ ] **USB Detection**: Detecta USB em Windows, Linux e WSL
- [ ] **Node Creation**: Cria nós via USB com formatação automática
- [ ] **Node Management**: CRUD completo de nós
- [ ] **Container Deploy**: Deploy básico de containers
- [ ] **Container Control**: Start, stop, restart, remove containers
- [ ] **Logs**: Visualização de logs de containers

### **Qualidade**
- [ ] **Cross-Platform**: Funciona em Windows, Linux e WSL
- [ ] **Error Handling**: Tratamento robusto de erros
- [ ] **Validation**: Validação de inputs e configurações
- [ ] **Logging**: Logs estruturados e informativos
- [ ] **Testing**: Testes unitários básicos (> 70% cobertura)

### **Usabilidade**
- [ ] **Intuitive Commands**: Comandos intuitivos e consistentes
- [ ] **Help System**: Help contextual e documentação
- [ ] **Output Formats**: Suporte a table, json, yaml
- [ ] **Progress Indicators**: Progress bars para operações longas
- [ ] **Error Messages**: Mensagens de erro claras e acionáveis

### **Performance**
- [ ] **Response Time**: Comandos respondem em < 2 segundos
- [ ] **Memory Usage**: Uso de memória < 100MB
- [ ] **CPU Usage**: Uso de CPU < 10% em idle
- [ ] **Disk Usage**: Binário < 50MB

---

## 📊 **Métricas de Sucesso**

### **Métricas Técnicas**
- **Test Coverage**: > 70%
- **Build Time**: < 2 minutos
- **Binary Size**: < 50MB
- **Memory Usage**: < 100MB
- **Response Time**: < 2 segundos

### **Métricas de Qualidade**
- **Bug Rate**: < 5 bugs críticos
- **Code Review**: 100% dos PRs revisados
- **Documentation**: 100% das funcionalidades documentadas
- **Cross-Platform**: Funciona em 3 plataformas

### **Métricas de Usabilidade**
- **Command Success Rate**: > 95%
- **User Satisfaction**: > 4.0/5
- **Help Usage**: < 20% dos usuários precisam de help
- **Error Recovery**: > 90% dos erros são recuperáveis

---

## ⚠️ **Riscos e Mitigações**

### **Riscos Técnicos**

#### **Complexidade de USB Detection**
- **Risco**: Dificuldade em detectar USB em diferentes plataformas
- **Probabilidade**: Alta
- **Impacto**: Alto
- **Mitigação**: Prototipagem early, fallbacks para cada plataforma

#### **Permissões de Sistema**
- **Risco**: Problemas de permissão para formatação de USB
- **Probabilidade**: Média
- **Impacto**: Alto
- **Mitigação**: Documentação clara, instruções de elevação

#### **Performance de Docker API**
- **Risco**: Docker API lenta ou indisponível
- **Probabilidade**: Baixa
- **Impacto**: Médio
- **Mitigação**: Timeouts, retry logic, fallbacks

### **Riscos de Negócio**

#### **Mudança de Requisitos**
- **Risco**: Requisitos mudam durante desenvolvimento
- **Probabilidade**: Média
- **Impacto**: Médio
- **Mitigação**: Sprints curtos, feedback contínuo

#### **Complexidade de Uso**
- **Risco**: CLI muito complexa para usuários
- **Probabilidade**: Baixa
- **Impacto**: Alto
- **Mitigação**: User testing, design simples

### **Riscos de Recursos**

#### **Tempo de Desenvolvimento**
- **Risco**: Desenvolvimento demora mais que esperado
- **Probabilidade**: Média
- **Impacto**: Médio
- **Mitigação**: Buffer time, scope management

#### **Qualidade do Código**
- **Risco**: Código de baixa qualidade
- **Probabilidade**: Baixa
- **Impacto**: Alto
- **Mitigação**: Code reviews, testing, linting

---

## 🚀 **Próximos Passos**

### **Imediatos (Próxima Semana)**
1. **Aprovação do Roadmap**: Validar com stakeholders
2. **Setup do Ambiente**: Configurar ferramentas de desenvolvimento
3. **Formação da Equipe**: Recrutar desenvolvedores se necessário
4. **Início do Sprint 1**: Começar desenvolvimento

### **Curto Prazo (1-2 Meses)**
1. **Completar MVP CLI**: 4 sprints de desenvolvimento
2. **Testes com Usuários**: Validar usabilidade
3. **Documentação**: Completar documentação do MVP
4. **Preparação para Fase 2**: Planejar API Foundation

### **Médio Prazo (3-6 Meses)**
1. **API Foundation**: Desenvolver backend robusto
2. **Web Interface**: Criar dashboard web
3. **Advanced Features**: Implementar funcionalidades avançadas
4. **Production Ready**: Preparar para produção

### **Longo Prazo (6-12 Meses)**
1. **Mobile & Desktop**: Interfaces adicionais
2. **Enterprise Features**: Funcionalidades empresariais
3. **Ecosystem**: Integrações e plugins
4. **Community**: Comunidade de desenvolvedores

---

## 🎯 **Conclusão**

O **MVP CLI Roadmap** fornece um plano detalhado e técnico para criar a primeira versão funcional do Management System. Com 4 sprints focados em funcionalidade core, o MVP estabelece uma base sólida para o desenvolvimento futuro.

### **Principais Destaques**

1. **Foco na Essência**: Funcionalidades core sem complexidade desnecessária
2. **Cross-Platform**: Suporte nativo a Windows, Linux e WSL
3. **Qualidade**: Padrões altos de qualidade e usabilidade
4. **Escalabilidade**: Arquitetura preparada para crescimento
5. **Pragmatismo**: MVP funcional em 8 semanas

### **Valor do MVP**

- **Validação Rápida**: Testa conceitos principais rapidamente
- **Feedback Early**: Permite feedback de usuários cedo
- **Base Sólida**: Estabelece fundação para desenvolvimento futuro
- **Time-to-Market**: Entrega valor rapidamente aos usuários

**O MVP CLI está pronto para transformar a forma como gerenciamos a Syntropy Cooperative Grid!** 🚀
