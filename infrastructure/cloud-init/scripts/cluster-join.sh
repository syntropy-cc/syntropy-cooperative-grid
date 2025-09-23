#!/bin/bash
# Syntropy Cooperative Grid - Cluster Join Script
# Conecta o nó ao cluster Syntropy e configura a participação

set -euo pipefail

# Cores de saída
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Função de log
log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}" | tee -a /opt/syntropy/logs/cluster-join.log
}

warn() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] WARNING: $1${NC}" | tee -a /opt/syntropy/logs/cluster-join.log
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR: $1${NC}" | tee -a /opt/syntropy/logs/cluster-join.log
}

# Carregar configurações
if [ -f /opt/syntropy/config/hardware.env ]; then
    source /opt/syntropy/config/hardware.env
else
    error "Arquivo de configuração de hardware não encontrado"
    exit 1
fi

if [ -f /opt/syntropy/config/network-discovery.env ]; then
    source /opt/syntropy/config/network-discovery.env
else
    error "Arquivo de configuração de descoberta de rede não encontrado"
    exit 1
fi

log "Iniciando processo de conexão ao cluster Syntropy..."

# Função para obter informações do cluster
get_cluster_info() {
    local leader_host="$1"
    
    if [ "$leader_host" = "self" ]; then
        log "Este nó é o líder do cluster"
        return 0
    fi
    
    log "Obtendo informações do cluster do líder: $leader_host"
    
    # Tentar obter informações via API
    local cluster_info
    if cluster_info=$(curl -s -k --connect-timeout 10 "https://$leader_host:8080/api/v1/cluster/info"); then
        log "Informações do cluster obtidas com sucesso"
        echo "$cluster_info" > /opt/syntropy/config/cluster-info.json
        chown syntropy:syntropy /opt/syntropy/config/cluster-info.json
        chmod 644 /opt/syntropy/config/cluster-info.json
        return 0
    else
        warn "Falha ao obter informações do cluster via API"
        return 1
    fi
}

# Função para registrar nó no cluster
register_node() {
    local leader_host="$1"
    
    if [ "$leader_host" = "self" ]; then
        log "Registrando como primeiro nó do cluster (líder)"
        return 0
    fi
    
    log "Registrando nó no cluster via líder: $leader_host"
    
    # Preparar dados de registro
    local node_data
    node_data=$(cat << EOF
{
    "name": "${NODE_NAME}",
    "type": "${HARDWARE_TYPE}",
    "description": "${NODE_DESCRIPTION}",
    "coordinates": "${COORDINATES}",
    "owner_key": "${OWNER_KEY}",
    "capabilities": {
        "can_be_leader": ${CAN_BE_LEADER},
        "can_be_worker": ${CAN_BE_WORKER},
        "cpu_cores": ${CPU_CORES},
        "memory_gb": ${MEMORY_GB},
        "storage_gb": ${STORAGE_GB}
    },
    "network": {
        "mesh_port": 51820,
        "api_port": 8080,
        "metrics_port": 9090
    }
}
EOF
)
    
    # Registrar nó
    local registration_response
    if registration_response=$(curl -s -k --connect-timeout 10 \
        -X POST \
        -H "Content-Type: application/json" \
        -d "$node_data" \
        "https://$leader_host:8080/api/v1/cluster/nodes/register"); then
        
        log "Nó registrado com sucesso no cluster"
        echo "$registration_response" > /opt/syntropy/config/node-registration.json
        chown syntropy:syntropy /opt/syntropy/config/node-registration.json
        chmod 644 /opt/syntropy/config/node-registration.json
        return 0
    else
        error "Falha ao registrar nó no cluster"
        return 1
    fi
}

# Função para configurar Kubernetes (se necessário)
configure_kubernetes() {
    local leader_host="$1"
    
    if [ "$leader_host" = "self" ]; then
        log "Configurando Kubernetes como master"
        configure_kubernetes_master
    else
        log "Configurando Kubernetes como worker"
        configure_kubernetes_worker "$leader_host"
    fi
}

# Função para configurar Kubernetes master
configure_kubernetes_master() {
    log "Configurando Kubernetes master..."
    
    # Inicializar cluster Kubernetes
    kubeadm init --pod-network-cidr=10.244.0.0/16 --apiserver-advertise-address=172.20.0.1
    
    # Configurar kubectl para usuário syntropy
    mkdir -p /home/syntropy/.kube
    cp /etc/kubernetes/admin.conf /home/syntropy/.kube/config
    chown syntropy:syntropy /home/syntropy/.kube/config
    chmod 600 /home/syntropy/.kube/config
    
    # Permitir pods no master (para ambiente de teste)
    kubectl taint nodes --all node-role.kubernetes.io/control-plane-
    
    # Instalar CNI (Calico)
    kubectl apply -f https://raw.githubusercontent.com/projectcalico/calico/v3.26.1/manifests/tigera-operator.yaml
    kubectl apply -f https://raw.githubusercontent.com/projectcalico/calico/v3.26.1/manifests/custom-resources.yaml
    
    # Aguardar pods estarem prontos
    kubectl wait --for=condition=ready pod -l k8s-app=calico-node -n kube-system --timeout=300s
    
    log "Kubernetes master configurado com sucesso"
}

# Função para configurar Kubernetes worker
configure_kubernetes_worker() {
    local leader_host="$1"
    
    log "Configurando Kubernetes worker para $leader_host"
    
    # Obter token de join do master
    local join_token
    if join_token=$(curl -s -k --connect-timeout 10 "https://$leader_host:8080/api/v1/kubernetes/join-token"); then
        log "Token de join obtido com sucesso"
        
        # Executar comando de join
        if echo "$join_token" | bash; then
            log "Kubernetes worker configurado com sucesso"
            return 0
        else
            error "Falha ao executar comando de join do Kubernetes"
            return 1
        fi
    else
        error "Falha ao obter token de join do Kubernetes"
        return 1
    fi
}

# Função para configurar mesh network
configure_mesh_network() {
    local leader_host="$1"
    
    log "Configurando mesh network..."
    
    if [ "$leader_host" = "self" ]; then
        log "Configurando mesh network como líder"
        configure_mesh_leader
    else
        log "Configurando mesh network como worker"
        configure_mesh_worker "$leader_host"
    fi
}

# Função para configurar mesh como líder
configure_mesh_leader() {
    log "Configurando mesh network como líder..."
    
    # Criar configuração de mesh
    cat > /opt/syntropy/config/mesh.yaml << EOF
mesh:
  role: "leader"
  interface: "wg0"
  subnet: "172.20.0.0/12"
  node_ip: "172.20.0.1/12"
  port: 51820
  
  peers: []
  
  routing:
    enabled: true
    method: "ospf"
  
  discovery:
    enabled: true
    method: "dns"
    hostname: "syntropy-discovery.local"
EOF
    
    chown syntropy:syntropy /opt/syntropy/config/mesh.yaml
    chmod 644 /opt/syntropy/config/mesh.yaml
    
    log "Mesh network configurado como líder"
}

# Função para configurar mesh como worker
configure_mesh_worker() {
    local leader_host="$1"
    
    log "Configurando mesh network como worker para $leader_host"
    
    # Obter configuração de mesh do líder
    local mesh_config
    if mesh_config=$(curl -s -k --connect-timeout 10 "https://$leader_host:8080/api/v1/mesh/config"); then
        log "Configuração de mesh obtida com sucesso"
        
        # Salvar configuração
        echo "$mesh_config" > /opt/syntropy/config/mesh.yaml
        chown syntropy:syntropy /opt/syntropy/config/mesh.yaml
        chmod 644 /opt/syntropy/config/mesh.yaml
        
        log "Mesh network configurado como worker"
        return 0
    else
        error "Falha ao obter configuração de mesh"
        return 1
    fi
}

# Função para configurar monitoramento
configure_monitoring() {
    log "Configurando monitoramento..."
    
    # Configurar Prometheus
    cat > /opt/syntropy/config/prometheus.yaml << EOF
prometheus:
  enabled: true
  port: 9090
  retention: "30d"
  
  targets:
    - "localhost:9100"  # Node Exporter
    - "localhost:8080"  # Syntropy Agent
  
  rules:
    - name: "syntropy.rules"
      interval: "30s"
      rules:
        - alert: "NodeDown"
          expr: "up == 0"
          for: "1m"
          labels:
            severity: "critical"
          annotations:
            summary: "Node {{ \$labels.instance }} is down"
        
        - alert: "HighCPUUsage"
          expr: "100 - (avg by(instance) (irate(node_cpu_seconds_total{mode=\"idle\"}[5m])) * 100) > 80"
          for: "5m"
          labels:
            severity: "warning"
          annotations:
            summary: "High CPU usage on {{ \$labels.instance }}"
        
        - alert: "HighMemoryUsage"
          expr: "(1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100 > 90"
          for: "5m"
          labels:
            severity: "warning"
          annotations:
            summary: "High memory usage on {{ \$labels.instance }}"
EOF
    
    chown syntropy:syntropy /opt/syntropy/config/prometheus.yaml
    chmod 644 /opt/syntropy/config/prometheus.yaml
    
    log "Monitoramento configurado com sucesso"
}

# Função para configurar backup
configure_backup() {
    log "Configurando backup..."
    
    # Configurar backup automático
    cat > /opt/syntropy/config/backup.yaml << EOF
backup:
  enabled: true
  schedule: "0 2 * * *"  # Diariamente às 2h
  
  destinations:
    - type: "local"
      path: "/opt/syntropy/backups"
      retention: "7d"
  
  sources:
    - "/opt/syntropy/config"
    - "/opt/syntropy/certs"
    - "/opt/syntropy/logs"
    - "/opt/syntropy/audit"
    - "/etc/wireguard"
    - "/etc/systemd/system/syntropy-agent.service"
  
  encryption:
    enabled: true
    algorithm: "aes256"
EOF
    
    chown syntropy:syntropy /opt/syntropy/config/backup.yaml
    chmod 644 /opt/syntropy/config/backup.yaml
    
    log "Backup configurado com sucesso"
}

# Função para iniciar serviços
start_services() {
    log "Iniciando serviços..."
    
    # Iniciar Syntropy Agent
    if systemctl start syntropy-agent; then
        log "Syntropy Agent iniciado com sucesso"
    else
        error "Falha ao iniciar Syntropy Agent"
        return 1
    fi
    
    # Aguardar agent estar pronto
    sleep 10
    
    # Verificar status
    if systemctl is-active --quiet syntropy-agent; then
        log "Syntropy Agent está ativo"
    else
        error "Syntropy Agent não está ativo"
        return 1
    fi
    
    # Iniciar monitoramento
    if systemctl start node_exporter; then
        log "Node Exporter iniciado com sucesso"
    else
        warn "Falha ao iniciar Node Exporter"
    fi
    
    log "Serviços iniciados com sucesso"
}

# Função para verificar conectividade
verify_connectivity() {
    local leader_host="$1"
    
    log "Verificando conectividade..."
    
    if [ "$leader_host" = "self" ]; then
        log "Verificando conectividade como líder"
        
        # Verificar se a API está respondendo
        if curl -s -k --connect-timeout 10 "https://localhost:8080/health" &> /dev/null; then
            log "API do líder está respondendo"
        else
            error "API do líder não está respondendo"
            return 1
        fi
    else
        log "Verificando conectividade com líder: $leader_host"
        
        # Verificar conectividade com líder
        if ping -c 1 -W 5 "$leader_host" &> /dev/null; then
            log "Conectividade com líder confirmada"
        else
            error "Falha na conectividade com líder"
            return 1
        fi
        
        # Verificar API do líder
        if curl -s -k --connect-timeout 10 "https://$leader_host:8080/health" &> /dev/null; then
            log "API do líder está respondendo"
        else
            error "API do líder não está respondendo"
            return 1
        fi
    fi
    
    log "Conectividade verificada com sucesso"
}

# Função principal
main() {
    log "Iniciando processo de conexão ao cluster..."
    
    # Verificar se a descoberta foi bem-sucedida
    if [ "$DISCOVERY_SUCCESS" != "true" ]; then
        error "Descoberta de rede não foi bem-sucedida"
        exit 1
    fi
    
    # Obter informações do cluster
    if ! get_cluster_info "$LEADER_HOST"; then
        error "Falha ao obter informações do cluster"
        exit 1
    fi
    
    # Registrar nó no cluster
    if ! register_node "$LEADER_HOST"; then
        error "Falha ao registrar nó no cluster"
        exit 1
    fi
    
    # Configurar Kubernetes
    if ! configure_kubernetes "$LEADER_HOST"; then
        error "Falha ao configurar Kubernetes"
        exit 1
    fi
    
    # Configurar mesh network
    if ! configure_mesh_network "$LEADER_HOST"; then
        error "Falha ao configurar mesh network"
        exit 1
    fi
    
    # Configurar monitoramento
    configure_monitoring
    
    # Configurar backup
    configure_backup
    
    # Iniciar serviços
    if ! start_services; then
        error "Falha ao iniciar serviços"
        exit 1
    fi
    
    # Verificar conectividade
    if ! verify_connectivity "$LEADER_HOST"; then
        error "Falha na verificação de conectividade"
        exit 1
    fi
    
    # Salvar informações de conexão
    cat > /opt/syntropy/config/cluster-connection.env << EOF
CLUSTER_CONNECTED=true
LEADER_HOST=$LEADER_HOST
CONNECTION_DATE=$(date -Iseconds)
NODE_ROLE=$INITIAL_ROLE
CLUSTER_STATUS=active
EOF
    
    chown syntropy:syntropy /opt/syntropy/config/cluster-connection.env
    chmod 600 /opt/syntropy/config/cluster-connection.env
    
    log "Conexão ao cluster concluída com sucesso!"
    log "Líder: $LEADER_HOST"
    log "Papel: $INITIAL_ROLE"
    log "Status: Ativo"
    
    # Executar auditoria
    /opt/syntropy/scripts/audit.sh
    
    exit 0
}

# Executar função principal
main "$@"
