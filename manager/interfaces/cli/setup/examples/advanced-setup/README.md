# Exemplo de Setup Avançado - Syntropy CLI

Este diretório demonstra configurações avançadas do setup component, incluindo customizações específicas por ambiente e topologias de rede complexas.

## Visão Geral

O setup avançado oferece controles detalhados para cenários empresariais e ambientes de produção, incluindo:
- Configurações customizadas por arquivo
- Variáveis de ambiente para automação CI/CD
- Topologias de rede específicas
- Configurações de segurança avançadas
- Integração com sistemas de monitoramento

## Estrutura dos Arquivos

- `custom-config.yaml` - Configuração avançada personalizável
- `environment-variables.env` - Variáveis de ambiente para automação
- `network-topology.yaml` - Definir topologia de rede específica
- `README.md` - Esta documentação de referência

## Ambientes Suportados

### Desenvolvimento Local
```bash
export SYNTHROPY_ENV=development
export SYNTHROPY_LOG_LEVEL=debug
```

### Staging/Produção
```bash
export SYNTHROPY_ENV=production
export SYNTHROPY_LOG_LEVEL=warn
export SYNTHROPY_API_ENDPOINT=https://api.production.syntropy.io
```

### CI/CD Pipeline
```bash
export SYNTHROPY_CI=true
export SYNTHROPY_SKIP_SERVICE=true
export SYNTHROPY_CONFIG_PATH=/shared/syntropy-config.yaml
```

## Configuração Avançada

### 1. Configuração com Variáveis de Ambiente

**Linux/macOS:**
```bash
# Carregar variáveis de ambiente
source environment-variables.env

# Executar setup com configuração avançada
SYNTHROPY_ENV=production ./../basic-setup/setup-basic.sh --config custom-config.yaml
```

**Windows:**
```powershell
# Carregar variáveis de ambiente
Get-Content environment-variables.env | ForEach-Object { 
    $key, $value = $_.Split('=')  
    [Environment]::SetEnvironmentVariable($key, $value, "Process")
}

# Executar setup
.\..\basic-setup\setup-basic.ps1 -ConfigPath custom-config.yaml
```

### 2. Configuração Multi-Plataforma

**Clusters Kubernetes:**
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: syntropy-config
data:
  manager.yaml: |
    # Conteúdo customizado para Kubernetes
    manager:
      home_dir: "/syntropy/shared"
      log_level: "info"
```

**Docker Swarm:**
```yaml
version: "3.8"
services:
  syntropy-manager:
    image: syntropy/cli:latest
    configs:
      - source: syntropy-config
        target: /etc/syntropy/manager.yaml
```

### 3. Configuração com Monitoramento

```yaml
# custom-config.yaml com monitoramento avançado
manager:
  debug: true
  metrics:
    enabled: true
    endpoint: "http://prometheus:9090"
```

## Cenários de Integração

### Integração com Infraestrutura Existente

#### 1. Integração com Service Discovery
```yaml
# Para ambientes com Consul/etcd/Zookeeper
discovery:
  consul:
    enabled: true
    endpoint: "consul.service.internal:8500"
    service_name: "syntropy-manager"
    check_interval: "30s"
```

#### 2. Integração com Load Balancer
```yaml
# Para ambientes com HAProxy/Nginx
load_balancer:
  primary: "https://lb.internal:443"
  failover: "https://lb-secondary.internal:443"
  health_check: "/health"
  retry_interval: 5s
```

#### 3. Integração com Message Broker
```yaml
# Para ambientes com RabbitMQ/Apache Kafka
messaging:
  broker:
    type: "kafka"
    endpoints:
      - "kafka1.internal:9092"
      - "kafka2.internal:9092"
    topic: "syntropy.events"
```

### Configurações de Segurança

#### Autenticação Multi-Relatório
```yaml
security:
  multi_tenant: true
  ldap:
    enabled: true
    server: "ldap.internal:389"
    base_dn: "DC=internal,DC=company"
  oauth:
    provider: "auth0"
    issuer: "https://company.auth0.com"
    client_id: "${OAUTH_CLIENT_ID}"
```

#### Criptografia de Comunicação
```yaml
security:
  encryption:
    tls_min_version: "TLS1.2"
    cipher_suites:
      - "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"
      - "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384"
```

#### Rotação de Chaves
```yaml
security:
  key_rotation:
    enabled: true
    rotation_interval: "24h"
    backup_before_rotation: true
```

## Troubleshooting Avançado

### Problemas de Conectividade Complexa

#### Firewall Corporativo
```bash
# Verificar conectividade necessária
./validate-network.sh --check-firewall
./validate-network.sh --check-proxy
```

#### Proxy Corporate
```yaml
network:
  proxy:
    enabled: true
    http_proxy: "http://proxy.corp:8080"
    https_proxy: "https://proxy.corp:8080"
    no_proxy: "*.internal,localhost"
```

### Problemas de Performance

#### Altas Cargas de Trabalho
```yaml
performance:
  concurrency:
    max_connections: 1000
    conn_per_owner: 10
  resources:
    cpu_limit: "2000m"
    memory_limit: "4Gi"
  queue:
    backlog_size: 10000
```

#### Problemas de Recursos
```bash
# Monitorar uso de recursos
./monitor-resources.sh --metrics-config=advanced-metrics.yaml
```

## Scripts de Automação Avançada

### Setup em Máquinas Múltiplas

#### Cluster de Alta Disponibilidade
```bash
# Cada nó do cluster
./setup-cluster-node.sh \
  --role manager \
  --cluster-name production \
  --seed-nodes "node1,node2,node3"
```

#### Load Balancing Automático
```yaml
cluster:
  nodes:
    - host: "manager1.syntropy.local"
      priority: 100
    - host: "manager2.syntropy.local"  
      priority: 100
    - host: "manager3.syntropy.local"
      priority: 50
```

### Integração DevOps

#### Ansible Playbook
```yaml
# deploy-syntropy.yml
---
- hosts: syntropy_managers
  tasks:
    - include_vars: "vars/{{ ansible_env.SYNTHROPY_ENV }}.yml"
    - service: name=syntropy state=started enabled=yes
```

#### Terraform Resource
```hcl
# main.tf
resource "syntropy_network" "main" {
  config_path = var.config_file
  environment {
    SYNTHROPY_ENV = var.environment
  }
}
```

## Exemplos de Configuração

### Enterprise Deployment
```bash
#!/bin/bash  
# Enterprise setup script
export SYNTHROPY_ENV=production
export SYNTHROPY_COMPANY=AcmeCorp
export SYNTHROPY_REQUIRE_AUTH=true

./setup-basic.sh \
  --config=/etc/syntropy/enterprise-conf.yaml \
  --multi-tenant \
  --audit-logging \
  --high-availability
```

### Container Microservices
```yaml
# docker-compose.yml example
version: "3.9"
services:
  syntropy-manager-1:
    image: syntropy/cli:latest
    environment:
      - SYNTHROPY_CLUSTER_ID=prod-cluster-1
  syntropy-manager-2:
    image: syntropy/cli:latest  
    environment:
      - SYNTHROPY_CLUSTER_ID=prod-cluster-1
```

### IoT Edge Deployment
```yaml
# IoT configuration minimal resource usage
performance:
  cpu_usage: "low"
  memory_footprint: "128MB" 
  storage_requirement: "512MB"

networking:
  low_bandwidth: true
  compression: "gzip"
  connection_pooling: false
```

## Monitoramento e Observabilidade

### Métricas Customizadas
```yaml
monitoring:
  metrics:
    prometheus:
      enabled: true
      path: "/metrics"
      port: 9090
    grafana:
      dashboard_id: "syntropy-overview"
  logging:
    levels:
      - "audit:info"
      - "security:warn"
      - "performance:debug"
```

### Alertas Avançados
```yaml
alerts:
  connection_failure:
    threshold: 5
    window: "5m" 
    receiver: "pagerduty"
  resource_exhaustion:
    cpu_threshold: 80
    memory_threshold: 85
    disk_threshold: 90
```

## Manutenção e Atualizações

### Rolling Updates
```bash
#!/bin/bash
# rolling-update.sh para atualização sem interrupção

for node in $(cat nodes-list.txt); do
  echo "Atualizando node: $node"
  ./update-node.sh --name $node --drain-first
  sleep 30
done
```

### Rollback Automático
```bash
#!/bin/bash
# auto-rollback.sh para reverter alterações problemáticas
if ./health-check.sh; then
  echo "Update successful"
else
  echo "Rolling back..."
  ./restore-config.sh --backup-dir backups/latest
fi
```

## Próximos Passos

1. **Leia a documentação de referência:** `../../GUIDE.md`
2. **Execute testes de validação:** `../validation-tests/`
3. **Configure monitoramento:** Adicione métricas personalizadas
4. **Implemente backups:** Use scripts em `../../scripts/automation/`
5. **Personalize configuração:** Adapte `custom-config.yaml` aos seu ambiente

## Suporte e Resolução de Problemas

- **Configuração Complexa:** Use logs detalhados com `--log-level=debug`
- **Problemas de Rede:** Execute diagnósticos de rede: `../../scripts/network/diagnose.sh`
- **Performance Issues:** Monitore recursos: `../../scripts/monitoring/resource-monitor.sh`
- **Segurança Avançada:** Consulte: `../../docs/security.md`
