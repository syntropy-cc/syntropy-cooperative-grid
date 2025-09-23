#!/bin/bash
# Syntropy Cooperative Grid - Installation Script
# Instala e configura o Syntropy Agent e todos os componentes necessários

set -euo pipefail

# Cores de saída
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Função de log
log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}" | tee -a /opt/syntropy/logs/installation.log
}

warn() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] WARNING: $1${NC}" | tee -a /opt/syntropy/logs/installation.log
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR: $1${NC}" | tee -a /opt/syntropy/logs/installation.log
}

# Carregar configurações
if [ -f /opt/syntropy/config/hardware.env ]; then
    source /opt/syntropy/config/hardware.env
else
    error "Arquivo de configuração de hardware não encontrado"
    exit 1
fi

log "Iniciando instalação do Syntropy Cooperative Grid..."

# Atualizar sistema
log "Atualizando sistema..."
apt-get update -y
apt-get upgrade -y

# Instalar dependências
log "Instalando dependências..."
apt-get install -y \
    curl \
    wget \
    git \
    htop \
    vim \
    net-tools \
    dnsutils \
    fail2ban \
    ufw \
    jq \
    openssl \
    ca-certificates \
    gnupg \
    lsb-release \
    apt-transport-https \
    ntp \
    rsync \
    unzip \
    tree \
    tmux \
    bc \
    nmap \
    netcat-openbsd

# Configurar Docker
log "Configurando Docker..."
if ! command -v docker &> /dev/null; then
    # Adicionar repositório Docker
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
    
    # Instalar Docker
    apt-get update -y
    apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin
fi

# Configurar Docker
systemctl enable docker
systemctl start docker
usermod -aG docker syntropy

# Configurar Kubernetes (kubectl)
log "Instalando kubectl..."
if ! command -v kubectl &> /dev/null; then
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
    install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
    rm kubectl
fi

# Instalar Wireguard
log "Instalando Wireguard..."
if ! command -v wg &> /dev/null; then
    apt-get install -y wireguard
fi

# Configurar firewall
log "Configurando firewall..."
ufw --force reset
ufw default deny incoming
ufw default allow outgoing
ufw allow ssh
ufw allow 6443/tcp
ufw allow 2379:2380/tcp
ufw allow 10250/tcp
ufw allow 10251/tcp
ufw allow 10252/tcp
ufw allow 30000:32767/tcp
ufw allow 51820/udp
ufw allow 8080/tcp
ufw allow 9090/tcp
ufw allow 9100/tcp
ufw --force enable

# Configurar fail2ban
log "Configurando fail2ban..."
systemctl enable fail2ban
systemctl start fail2ban

# Configurar NTP
log "Configurando NTP..."
systemctl enable ntp
systemctl start ntp

# Instalar Prometheus Node Exporter
log "Instalando Prometheus Node Exporter..."
if ! command -v node_exporter &> /dev/null; then
    wget https://github.com/prometheus/node_exporter/releases/latest/download/node_exporter-1.6.1.linux-amd64.tar.gz
    tar xvfz node_exporter-1.6.1.linux-amd64.tar.gz
    cp node_exporter-1.6.1.linux-amd64/node_exporter /usr/local/bin/
    chmod +x /usr/local/bin/node_exporter
    rm -rf node_exporter-1.6.1.linux-amd64*
fi

# Configurar Prometheus Node Exporter
cat > /etc/systemd/system/node_exporter.service << EOF
[Unit]
Description=Prometheus Node Exporter
After=network.target

[Service]
Type=simple
User=prometheus
Group=prometheus
ExecStart=/usr/local/bin/node_exporter
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

# Criar usuário prometheus
useradd --no-create-home --shell /bin/false prometheus
chown prometheus:prometheus /usr/local/bin/node_exporter

systemctl daemon-reload
systemctl enable node_exporter
systemctl start node_exporter

# Download e instalação do Syntropy Agent
log "Baixando Syntropy Agent..."
AGENT_VERSION="latest"
AGENT_URL="https://github.com/syntropy-cooperative-grid/agent/releases/latest/download/syntropy-agent-linux-amd64"

# Tentar baixar a versão mais recente
if ! curl -L "$AGENT_URL" -o /opt/syntropy/bin/syntropy-agent; then
    warn "Falha ao baixar versão mais recente, tentando versão estável..."
    AGENT_VERSION="v1.0.0"
    AGENT_URL="https://github.com/syntropy-cooperative-grid/agent/releases/download/$AGENT_VERSION/syntropy-agent-linux-amd64"
    
    if ! curl -L "$AGENT_URL" -o /opt/syntropy/bin/syntropy-agent; then
        error "Falha ao baixar Syntropy Agent"
        exit 1
    fi
fi

chmod +x /opt/syntropy/bin/syntropy-agent
chown syntropy:syntropy /opt/syntropy/bin/syntropy-agent

# Configurar certificados
log "Configurando certificados..."
if [ ! -f /opt/syntropy/certs/node.crt ] || [ ! -f /opt/syntropy/certs/node.key ] || [ ! -f /opt/syntropy/certs/ca.crt ]; then
    error "Certificados não encontrados"
    exit 1
fi

chmod 600 /opt/syntropy/certs/*
chown syntropy:syntropy /opt/syntropy/certs/*

# Configurar Syntropy Agent
log "Configurando Syntropy Agent..."
cat > /opt/syntropy/config/agent.yaml << EOF
node:
  name: "${NODE_NAME}"
  type: "${HARDWARE_TYPE}"
  description: "${NODE_DESCRIPTION}"
  coordinates: "${COORDINATES}"
  owner_key: "${OWNER_KEY}"

network:
  discovery_endpoints:
    - "https://${DISCOVERY_SERVER}:8443"
  mesh_port: 51820
  api_port: 8080

security:
  tls:
    enabled: true
    cert_file: "/opt/syntropy/certs/node.crt"
    key_file: "/opt/syntropy/certs/node.key"
    ca_file: "/opt/syntropy/certs/ca.crt"
  
  firewall:
    enabled: true
    default_policy: "deny"
    allow_ssh: true
    allow_management: true

logging:
  level: "info"
  file: "/opt/syntropy/logs/agent.log"
  max_size: "100MB"
  max_files: 5

metrics:
  enabled: true
  port: 9090
  path: "/metrics"

monitoring:
  prometheus:
    enabled: true
    port: 9100
    path: "/metrics"
  
  health_checks:
    - name: "docker_service"
      command: "systemctl is-active docker"
      interval: "30s"
    - name: "ssh_service"
      command: "systemctl is-active ssh"
      interval: "30s"
    - name: "disk_space"
      command: "df -h / | awk 'NR==2{print \$5}' | sed 's/%//'"
      threshold: 90
      interval: "60s"
EOF

chown syntropy:syntropy /opt/syntropy/config/agent.yaml
chmod 644 /opt/syntropy/config/agent.yaml

# Configurar systemd service
log "Configurando systemd service..."
cat > /etc/systemd/system/syntropy-agent.service << EOF
[Unit]
Description=Syntropy Cooperative Grid Agent
After=network.target docker.service
Wants=docker.service

[Service]
Type=simple
User=syntropy
Group=syntropy
WorkingDirectory=/opt/syntropy
ExecStart=/opt/syntropy/bin/syntropy-agent --config=/opt/syntropy/config/agent.yaml
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

# Configurações de segurança
NoNewPrivileges=yes
PrivateTmp=yes
ProtectSystem=strict
ProtectHome=yes
ReadWritePaths=/opt/syntropy

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable syntropy-agent

# Configurar logrotate
log "Configurando logrotate..."
cat > /etc/logrotate.d/syntropy << EOF
/opt/syntropy/logs/*.log {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    notifempty
    create 644 syntropy syntropy
    postrotate
        systemctl reload syntropy-agent
    endscript
}
EOF

# Configurar auditoria
log "Configurando sistema de auditoria..."
mkdir -p /opt/syntropy/audit
chown syntropy:syntropy /opt/syntropy/audit
chmod 755 /opt/syntropy/audit

# Criar script de auditoria
cat > /opt/syntropy/scripts/audit.sh << 'EOF'
#!/bin/bash
# Script de auditoria do Syntropy

AUDIT_LOG="/opt/syntropy/audit/audit.log"
TIMESTAMP=$(date -Iseconds)

# Função de log de auditoria
audit_log() {
    echo "[$TIMESTAMP] $1" >> "$AUDIT_LOG"
}

# Auditar eventos do sistema
audit_log "SYSTEM: Boot completed"
audit_log "SYSTEM: Syntropy Agent started"
audit_log "NETWORK: Discovery completed"
audit_log "SECURITY: Certificates loaded"
audit_log "MONITORING: Metrics collection started"

# Rotacionar log de auditoria se necessário
if [ -f "$AUDIT_LOG" ] && [ $(stat -c%s "$AUDIT_LOG") -gt 10485760 ]; then
    mv "$AUDIT_LOG" "$AUDIT_LOG.$(date +%Y%m%d%H%M%S)"
    touch "$AUDIT_LOG"
    chown syntropy:syntropy "$AUDIT_LOG"
fi
EOF

chmod +x /opt/syntropy/scripts/audit.sh
chown syntropy:syntropy /opt/syntropy/scripts/audit.sh

# Executar auditoria inicial
/opt/syntropy/scripts/audit.sh

# Configurar cron para auditoria
echo "*/5 * * * * syntropy /opt/syntropy/scripts/audit.sh" >> /etc/crontab

# Configurar backup automático
log "Configurando backup automático..."
cat > /opt/syntropy/scripts/backup.sh << 'EOF'
#!/bin/bash
# Script de backup do Syntropy

BACKUP_DIR="/opt/syntropy/backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/syntropy_backup_$TIMESTAMP.tar.gz"

mkdir -p "$BACKUP_DIR"

# Criar backup
tar -czf "$BACKUP_FILE" \
    /opt/syntropy/config \
    /opt/syntropy/certs \
    /opt/syntropy/logs \
    /opt/syntropy/audit \
    /etc/wireguard \
    /etc/systemd/system/syntropy-agent.service

# Manter apenas os 7 backups mais recentes
ls -t "$BACKUP_DIR"/syntropy_backup_*.tar.gz | tail -n +8 | xargs -r rm

chown syntropy:syntropy "$BACKUP_FILE"
EOF

chmod +x /opt/syntropy/scripts/backup.sh
chown syntropy:syntropy /opt/syntropy/scripts/backup.sh

# Configurar cron para backup
echo "0 2 * * * syntropy /opt/syntropy/scripts/backup.sh" >> /etc/crontab

# Configurar monitoramento de saúde
log "Configurando monitoramento de saúde..."
cat > /opt/syntropy/scripts/health-check.sh << 'EOF'
#!/bin/bash
# Script de verificação de saúde do Syntropy

HEALTH_LOG="/opt/syntropy/logs/health.log"
TIMESTAMP=$(date -Iseconds)

# Função de log de saúde
health_log() {
    echo "[$TIMESTAMP] $1" >> "$HEALTH_LOG"
}

# Verificar serviços
if systemctl is-active --quiet docker; then
    health_log "HEALTH: Docker service OK"
else
    health_log "HEALTH: Docker service FAILED"
    systemctl restart docker
fi

if systemctl is-active --quiet syntropy-agent; then
    health_log "HEALTH: Syntropy Agent OK"
else
    health_log "HEALTH: Syntropy Agent FAILED"
    systemctl restart syntropy-agent
fi

if systemctl is-active --quiet ssh; then
    health_log "HEALTH: SSH service OK"
else
    health_log "HEALTH: SSH service FAILED"
    systemctl restart ssh
fi

# Verificar espaço em disco
DISK_USAGE=$(df -h / | awk 'NR==2{print $5}' | sed 's/%//')
if [ "$DISK_USAGE" -gt 90 ]; then
    health_log "HEALTH: Disk usage critical: ${DISK_USAGE}%"
else
    health_log "HEALTH: Disk usage OK: ${DISK_USAGE}%"
fi

# Verificar memória
MEMORY_USAGE=$(free | awk 'NR==2{printf "%.0f", $3*100/$2}')
if [ "$MEMORY_USAGE" -gt 90 ]; then
    health_log "HEALTH: Memory usage critical: ${MEMORY_USAGE}%"
else
    health_log "HEALTH: Memory usage OK: ${MEMORY_USAGE}%"
fi
EOF

chmod +x /opt/syntropy/scripts/health-check.sh
chown syntropy:syntropy /opt/syntropy/scripts/health-check.sh

# Configurar cron para verificação de saúde
echo "*/10 * * * * syntropy /opt/syntropy/scripts/health-check.sh" >> /etc/crontab

# Criar arquivo de status
log "Criando arquivo de status..."
cat > /opt/syntropy/.status << EOF
SYNTROPY_STATUS=installed
INSTALLATION_DATE=$(date -Iseconds)
NODE_NAME=${NODE_NAME}
NODE_TYPE=${HARDWARE_TYPE}
VERSION=${AGENT_VERSION}
EOF

chown syntropy:syntropy /opt/syntropy/.status
chmod 644 /opt/syntropy/.status

log "Instalação do Syntropy Cooperative Grid concluída com sucesso!"
log "Status: $(cat /opt/syntropy/.status)"

# Executar verificação de saúde inicial
/opt/syntropy/scripts/health-check.sh

exit 0
