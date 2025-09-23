#!/bin/bash
# Syntropy Cooperative Grid - Hardware Detection Script
# Detecta automaticamente o tipo de hardware e configura o nó adequadamente

set -euo pipefail

# Cores de saída
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Função de log
log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}" | tee -a /opt/syntropy/logs/hardware-detection.log
}

warn() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] WARNING: $1${NC}" | tee -a /opt/syntropy/logs/hardware-detection.log
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR: $1${NC}" | tee -a /opt/syntropy/logs/hardware-detection.log
}

# Criar diretório de logs se não existir
mkdir -p /opt/syntropy/logs

log "Iniciando detecção de hardware..."

# Detectar informações básicas do sistema
CPU_CORES=$(nproc)
MEMORY_GB=$(free -g | awk 'NR==2{print $2}')
STORAGE_GB=$(df -BG / | awk 'NR==2{print $2}' | sed 's/G//')
ARCH=$(uname -m)
KERNEL=$(uname -r)

log "Informações básicas detectadas:"
log "  CPU Cores: $CPU_CORES"
log "  Memory: ${MEMORY_GB}GB"
log "  Storage: ${STORAGE_GB}GB"
log "  Architecture: $ARCH"
log "  Kernel: $KERNEL"

# Detectar tipo de hardware baseado em características
HARDWARE_TYPE="unknown"
CAN_BE_LEADER="false"
CAN_BE_WORKER="true"
INITIAL_ROLE="worker"

# Detectar se é um servidor dedicado
if [ "$CPU_CORES" -ge 16 ] && [ "$MEMORY_GB" -ge 32 ] && [ "$STORAGE_GB" -ge 1000 ]; then
    HARDWARE_TYPE="server"
    CAN_BE_LEADER="true"
    CAN_BE_WORKER="true"
    INITIAL_ROLE="leader"
    log "Hardware detectado: SERVER (Alta capacidade)"
    
# Detectar se é um home server
elif [ "$CPU_CORES" -ge 4 ] && [ "$MEMORY_GB" -ge 8 ] && [ "$STORAGE_GB" -ge 500 ]; then
    HARDWARE_TYPE="home_server"
    CAN_BE_LEADER="true"
    CAN_BE_WORKER="true"
    INITIAL_ROLE="worker"
    log "Hardware detectado: HOME SERVER (Média capacidade)"
    
# Detectar se é um computador pessoal
elif [ "$CPU_CORES" -ge 2 ] && [ "$MEMORY_GB" -ge 4 ] && [ "$STORAGE_GB" -ge 100 ]; then
    HARDWARE_TYPE="personal_computer"
    CAN_BE_LEADER="false"
    CAN_BE_WORKER="true"
    INITIAL_ROLE="worker"
    log "Hardware detectado: PERSONAL COMPUTER (Capacidade variável)"
    
# Detectar se é um dispositivo móvel ou IoT
elif [ "$CPU_CORES" -ge 1 ] && [ "$MEMORY_GB" -ge 1 ] && [ "$STORAGE_GB" -ge 16 ]; then
    HARDWARE_TYPE="mobile_iot"
    CAN_BE_LEADER="false"
    CAN_BE_WORKER="true"
    INITIAL_ROLE="worker"
    log "Hardware detectado: MOBILE/IoT (Baixa capacidade)"
    
else
    HARDWARE_TYPE="unknown"
    CAN_BE_LEADER="false"
    CAN_BE_WORKER="true"
    INITIAL_ROLE="worker"
    warn "Hardware não classificado, configurando como worker genérico"
fi

# Detectar características específicas
GPU_PRESENT="false"
if command -v nvidia-smi &> /dev/null; then
    GPU_PRESENT="true"
    log "GPU NVIDIA detectada"
elif command -v lspci &> /dev/null && lspci | grep -i vga | grep -i amd &> /dev/null; then
    GPU_PRESENT="true"
    log "GPU AMD detectada"
elif command -v lspci &> /dev/null && lspci | grep -i vga | grep -i intel &> /dev/null; then
    GPU_PRESENT="true"
    log "GPU Intel detectada"
fi

# Detectar se é um ambiente virtualizado
VIRTUALIZED="false"
if [ -f /proc/cpuinfo ] && grep -q "hypervisor" /proc/cpuinfo; then
    VIRTUALIZED="true"
    log "Ambiente virtualizado detectado"
fi

# Detectar se é um ambiente containerizado
CONTAINERIZED="false"
if [ -f /.dockerenv ] || [ -f /run/.containerenv ]; then
    CONTAINERIZED="true"
    log "Ambiente containerizado detectado"
fi

# Detectar características de rede
NETWORK_INTERFACES=$(ip -o link show | grep -v lo | wc -l)
log "Interfaces de rede detectadas: $NETWORK_INTERFACES"

# Detectar se tem conectividade com internet
INTERNET_ACCESS="false"
if ping -c 1 8.8.8.8 &> /dev/null; then
    INTERNET_ACCESS="true"
    log "Conectividade com internet confirmada"
else
    warn "Sem conectividade com internet"
fi

# Salvar informações de hardware em arquivo de configuração
cat > /opt/syntropy/config/hardware.yaml << EOF
hardware:
  type: "$HARDWARE_TYPE"
  cpu_cores: $CPU_CORES
  memory_gb: $MEMORY_GB
  storage_gb: $STORAGE_GB
  architecture: "$ARCH"
  kernel: "$KERNEL"
  gpu_present: $GPU_PRESENT
  virtualized: $VIRTUALIZED
  containerized: $CONTAINERIZED
  network_interfaces: $NETWORK_INTERFACES
  internet_access: $INTERNET_ACCESS

capabilities:
  can_be_leader: $CAN_BE_LEADER
  can_be_worker: $CAN_BE_WORKER
  initial_role: "$INITIAL_ROLE"
  max_containers: $((CPU_CORES * 2))
  max_memory_usage: $((MEMORY_GB * 80 / 100))
  max_storage_usage: $((STORAGE_GB * 70 / 100))

performance:
  estimated_throughput: "$(echo "scale=2; $CPU_CORES * $MEMORY_GB / 10" | bc)"
  estimated_latency: "$(if [ $HARDWARE_TYPE = "server" ]; then echo "low"; elif [ $HARDWARE_TYPE = "home_server" ]; then echo "medium"; else echo "high"; fi)"
  reliability_score: "$(if [ $HARDWARE_TYPE = "server" ]; then echo "high"; elif [ $HARDWARE_TYPE = "home_server" ]; then echo "medium"; else echo "low"; fi)"
EOF

# Configurar permissões
chown syntropy:syntropy /opt/syntropy/config/hardware.yaml
chmod 644 /opt/syntropy/config/hardware.yaml

log "Detecção de hardware concluída com sucesso!"
log "Tipo de hardware: $HARDWARE_TYPE"
log "Papel inicial: $INITIAL_ROLE"
log "Pode ser líder: $CAN_BE_LEADER"
log "Pode ser worker: $CAN_BE_WORKER"

# Exportar variáveis para uso em outros scripts
echo "HARDWARE_TYPE=$HARDWARE_TYPE" >> /opt/syntropy/config/hardware.env
echo "CAN_BE_LEADER=$CAN_BE_LEADER" >> /opt/syntropy/config/hardware.env
echo "CAN_BE_WORKER=$CAN_BE_WORKER" >> /opt/syntropy/config/hardware.env
echo "INITIAL_ROLE=$INITIAL_ROLE" >> /opt/syntropy/config/hardware.env
echo "CPU_CORES=$CPU_CORES" >> /opt/syntropy/config/hardware.env
echo "MEMORY_GB=$MEMORY_GB" >> /opt/syntropy/config/hardware.env
echo "STORAGE_GB=$STORAGE_GB" >> /opt/syntropy/config/hardware.env

chown syntropy:syntropy /opt/syntropy/config/hardware.env
chmod 600 /opt/syntropy/config/hardware.env

exit 0
