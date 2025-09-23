# Scripts de Instalação - Syntropy Cooperative Grid

## Visão Geral

Este documento detalha os scripts bash implementados para a instalação automática de nós da rede Syntropy Cooperative Grid. Cada script tem uma função específica no processo de boot e configuração.

## Índice

1. [hardware-detection.sh](#hardware-detectionsh)
2. [network-discovery.sh](#network-discoverysh)
3. [syntropy-install.sh](#syntropy-installsh)
4. [cluster-join.sh](#cluster-joinsh)
5. [Fluxo de Execução](#fluxo-de-execução)

---

## hardware-detection.sh

### Propósito
Detecta automaticamente o tipo de hardware e configura o nó adequadamente.

### Funcionalidades
- Detecção de CPU, RAM, storage
- Classificação do tipo de hardware
- Determinação de papéis (líder/worker)
- Configuração de capacidades

### Algoritmo de Classificação
```bash
# Servidor Dedicado (Alta Capacidade)
if [ CPU_CORES >= 16 ] && [ MEMORY_GB >= 32 ] && [ STORAGE_GB >= 1000 ]; then
    HARDWARE_TYPE="server"
    CAN_BE_LEADER="true"
    INITIAL_ROLE="leader"

# Home Server (Média Capacidade)
elif [ CPU_CORES >= 4 ] && [ MEMORY_GB >= 8 ] && [ STORAGE_GB >= 500 ]; then
    HARDWARE_TYPE="home_server"
    CAN_BE_LEADER="true"
    INITIAL_ROLE="worker"

# Computador Pessoal (Capacidade Variável)
elif [ CPU_CORES >= 2 ] && [ MEMORY_GB >= 4 ] && [ STORAGE_GB >= 100 ]; then
    HARDWARE_TYPE="personal_computer"
    CAN_BE_LEADER="false"
    INITIAL_ROLE="worker"

# Mobile/IoT (Baixa Capacidade)
else
    HARDWARE_TYPE="mobile_iot"
    CAN_BE_LEADER="false"
    INITIAL_ROLE="worker"
fi
```

### Saída
- Arquivo `/opt/syntropy/config/hardware.yaml`
- Arquivo `/opt/syntropy/config/hardware.env`
- Logs em `/opt/syntropy/logs/hardware-detection.log`

---

## network-discovery.sh

### Propósito
Descobre automaticamente a rede Syntropy e se conecta a ela.

### Métodos de Descoberta
1. **DNS**: Resolve hostnames como `syntropy-discovery.local`
2. **Broadcast**: Envia broadcast na rede local
3. **Multicast**: Usa multicast para descoberta
4. **Configuração Manual**: Usa hosts pré-configurados

### Algoritmo de Descoberta
```bash
discover_network() {
    # Tentar DNS primeiro
    if discovered_host=$(discover_via_dns); then
        return 0
    fi
    
    # Tentar broadcast
    if discovered_host=$(discover_via_broadcast); then
        return 0
    fi
    
    # Tentar multicast
    if discovered_host=$(discover_via_multicast); then
        return 0
    fi
    
    # Tentar configuração manual
    if discovered_host=$(discover_via_manual_config); then
        return 0
    fi
    
    # Se é o primeiro nó, vira líder
    if [ "$INITIAL_ROLE" = "leader" ]; then
        echo "self"
        return 0
    fi
    
    return 1
}
```

### Configuração de Rede
- **Wireguard**: Configuração automática de mesh network
- **Kubernetes**: Configuração de cluster
- **API**: Configuração de endpoints

### Saída
- Arquivo `/opt/syntropy/config/network-connection.yaml`
- Arquivo `/opt/syntropy/config/network-discovery.env`
- Logs em `/opt/syntropy/logs/network-discovery.log`

---

## syntropy-install.sh

### Propósito
Instala e configura o Syntropy Agent e todos os componentes necessários.

### Funcionalidades
- Atualização do sistema
- Instalação de dependências
- Configuração do Docker
- Instalação do Kubernetes
- Configuração de segurança
- Instalação do Syntropy Agent

### Processo de Instalação
```bash
# 1. Atualizar sistema
apt-get update -y
apt-get upgrade -y

# 2. Instalar dependências
apt-get install -y docker.io kubectl wireguard

# 3. Configurar serviços
systemctl enable docker
systemctl start docker

# 4. Configurar firewall
ufw --force enable
ufw default deny incoming
ufw allow ssh

# 5. Instalar Syntropy Agent
curl -L https://github.com/syntropy-cooperative-grid/agent/releases/latest/download/syntropy-agent-linux-amd64 -o /opt/syntropy/bin/syntropy-agent
chmod +x /opt/syntropy/bin/syntropy-agent

# 6. Configurar systemd service
systemctl enable syntropy-agent
systemctl start syntropy-agent
```

### Configurações de Segurança
- **Firewall**: Regras específicas para Syntropy
- **Fail2ban**: Proteção contra ataques
- **Certificados**: Configuração de TLS
- **Auditoria**: Sistema de logs

### Saída
- Arquivo `/opt/syntropy/config/agent.yaml`
- Serviço systemd `syntropy-agent.service`
- Logs em `/opt/syntropy/logs/installation.log`

---

## cluster-join.sh

### Propósito
Conecta o nó ao cluster Syntropy e configura a participação.

### Funcionalidades
- Obtenção de informações do cluster
- Registro do nó
- Configuração do Kubernetes
- Configuração do mesh network
- Verificação de conectividade

### Processo de Conexão
```bash
# 1. Obter informações do cluster
get_cluster_info "$LEADER_HOST"

# 2. Registrar nó
register_node "$LEADER_HOST"

# 3. Configurar Kubernetes
configure_kubernetes "$LEADER_HOST"

# 4. Configurar mesh network
configure_mesh_network "$LEADER_HOST"

# 5. Iniciar serviços
start_services

# 6. Verificar conectividade
verify_connectivity "$LEADER_HOST"
```

### Configurações de Cluster
- **Kubernetes**: Master ou worker
- **Mesh Network**: Líder ou worker
- **Monitoramento**: Prometheus e Grafana
- **Backup**: Sistema automático

### Saída
- Arquivo `/opt/syntropy/config/cluster-connection.env`
- Configuração de Kubernetes
- Configuração de Wireguard
- Logs em `/opt/syntropy/logs/cluster-join.log`

---

## Fluxo de Execução

### Ordem de Execução
1. **hardware-detection.sh**: Detecta hardware e configura papéis
2. **network-discovery.sh**: Descobre rede e configura conectividade
3. **syntropy-install.sh**: Instala componentes e serviços
4. **cluster-join.sh**: Conecta ao cluster e inicia operação

### Dependências
- **hardware-detection.sh**: Nenhuma
- **network-discovery.sh**: Requer hardware-detection.sh
- **syntropy-install.sh**: Requer hardware-detection.sh
- **cluster-join.sh**: Requer todos os anteriores

### Configuração Automática
- **Cron Jobs**: Scripts de auditoria e backup
- **Systemd Services**: Serviços do Syntropy
- **Logrotate**: Rotação de logs
- **Monitoramento**: Verificação de saúde

### Logs e Auditoria
- **Logs Individuais**: Cada script tem seu próprio log
- **Log de Auditoria**: `/opt/syntropy/audit/audit.log`
- **Log de Saúde**: `/opt/syntropy/logs/health.log`
- **Retenção**: 30 dias para logs, 90 dias para auditoria

---

## Conclusão

Os scripts de instalação fornecem uma base sólida para a configuração automática de nós da rede Syntropy Cooperative Grid. Com detecção inteligente de hardware, descoberta robusta de rede e instalação completa de componentes, o sistema permite que nós sejam configurados automaticamente desde o boot até a operação completa.
