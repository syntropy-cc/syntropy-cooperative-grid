#!/bin/bash
# Syntropy Cooperative Grid - Network Discovery Script
# Descobre automaticamente a rede Syntropy e se conecta a ela

set -euo pipefail

# Cores de saída
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Função de log
log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}" | tee -a /opt/syntropy/logs/network-discovery.log
}

warn() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] WARNING: $1${NC}" | tee -a /opt/syntropy/logs/network-discovery.log
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR: $1${NC}" | tee -a /opt/syntropy/logs/network-discovery.log
}

# Carregar configurações de hardware
if [ -f /opt/syntropy/config/hardware.env ]; then
    source /opt/syntropy/config/hardware.env
else
    error "Arquivo de configuração de hardware não encontrado"
    exit 1
fi

log "Iniciando descoberta de rede Syntropy..."

# Configurações de descoberta
DISCOVERY_TIMEOUT=30
DISCOVERY_RETRIES=3
DISCOVERY_INTERVAL=5

# Métodos de descoberta
discover_via_dns() {
    local dns_hostnames=(
        "syntropy-discovery.local"
        "syntropy-mesh.local"
        "syntropy-leader.local"
        "syntropy.local"
    )
    
    log "Tentando descoberta via DNS..."
    
    for hostname in "${dns_hostnames[@]}"; do
        log "Testando hostname: $hostname"
        
        if nslookup "$hostname" &> /dev/null; then
            local ip=$(nslookup "$hostname" | awk '/^Address: / { print $2 }' | tail -1)
            if [ -n "$ip" ] && [ "$ip" != "127.0.0.1" ]; then
                log "Hostname $hostname resolve para $ip"
                
                # Testar conectividade
                if ping -c 1 -W 2 "$ip" &> /dev/null; then
                    log "Conectividade confirmada com $ip"
                    echo "$ip"
                    return 0
                fi
            fi
        fi
    done
    
    return 1
}

discover_via_broadcast() {
    log "Tentando descoberta via broadcast..."
    
    # Criar arquivo temporário para capturar respostas
    local temp_file=$(mktemp)
    
    # Enviar broadcast na rede local
    local network=$(ip route | grep default | awk '{print $3}' | head -1)
    local subnet=$(ip route | grep "$network" | awk '{print $1}' | head -1)
    
    if [ -z "$subnet" ]; then
        warn "Não foi possível determinar a sub-rede local"
        return 1
    fi
    
    log "Enviando broadcast na sub-rede: $subnet"
    
    # Usar nmap para descobrir hosts ativos na porta 8080 (API Syntropy)
    if command -v nmap &> /dev/null; then
        nmap -p 8080 --open "$subnet" -oG "$temp_file" &> /dev/null
        
        if [ -s "$temp_file" ]; then
            local hosts=$(grep "8080/open" "$temp_file" | awk '{print $2}' | head -5)
            
            for host in $hosts; do
                log "Testando host descoberto: $host"
                
                # Testar se é um nó Syntropy
                if curl -s -k --connect-timeout 5 "https://$host:8080/health" &> /dev/null; then
                    log "Nó Syntropy encontrado: $host"
                    rm -f "$temp_file"
                    echo "$host"
                    return 0
                fi
            done
        fi
    fi
    
    rm -f "$temp_file"
    return 1
}

discover_via_multicast() {
    log "Tentando descoberta via multicast..."
    
    # Configurar multicast para descoberta Syntropy
    local multicast_group="224.0.0.251"
    local multicast_port="5353"
    
    # Criar arquivo temporário para capturar respostas
    local temp_file=$(mntemp)
    
    # Enviar requisição de descoberta via multicast
    echo "SYNTROPY_DISCOVERY_REQUEST" | nc -u -w 2 "$multicast_group" "$multicast_port" &> /dev/null &
    
    # Aguardar respostas
    sleep 2
    
    # Capturar respostas (implementação simplificada)
    if [ -f "$temp_file" ] && [ -s "$temp_file" ]; then
        local host=$(head -1 "$temp_file")
        log "Nó Syntropy encontrado via multicast: $host"
        rm -f "$temp_file"
        echo "$host"
        return 0
    fi
    
    rm -f "$temp_file"
    return 1
}

discover_via_manual_config() {
    log "Tentando descoberta via configuração manual..."
    
    # Verificar se há configuração manual no arquivo de configuração
    if [ -f /opt/syntropy/config/network.yaml ]; then
        local manual_hosts=$(grep -E "^\s*hosts:" /opt/syntropy/config/network.yaml | awk '{print $2}' | tr -d '[]' | tr ',' '\n')
        
        for host in $manual_hosts; do
            if [ -n "$host" ]; then
                log "Testando host manual: $host"
                
                if ping -c 1 -W 2 "$host" &> /dev/null; then
                    log "Host manual acessível: $host"
                    echo "$host"
                    return 0
                fi
            fi
        done
    fi
    
    return 1
}

# Função principal de descoberta
discover_network() {
    local discovered_host=""
    
    # Tentar diferentes métodos de descoberta
    if discovered_host=$(discover_via_dns); then
        log "Descoberta bem-sucedida via DNS: $discovered_host"
        echo "$discovered_host"
        return 0
    fi
    
    if discovered_host=$(discover_via_broadcast); then
        log "Descoberta bem-sucedida via broadcast: $discovered_host"
        echo "$discovered_host"
        return 0
    fi
    
    if discovered_host=$(discover_via_multicast); then
        log "Descoberta bem-sucedida via multicast: $discovered_host"
        echo "$discovered_host"
        return 0
    fi
    
    if discovered_host=$(discover_via_manual_config); then
        log "Descoberta bem-sucedida via configuração manual: $discovered_host"
        echo "$discovered_host"
        return 0
    fi
    
    # Se nenhum método funcionou, verificar se este é o primeiro nó
    if [ "$INITIAL_ROLE" = "leader" ]; then
        log "Este é o primeiro nó da rede, iniciando como líder"
        echo "self"
        return 0
    fi
    
    error "Não foi possível descobrir a rede Syntropy"
    return 1
}

# Função para validar conectividade com um host
validate_connectivity() {
    local host="$1"
    
    if [ "$host" = "self" ]; then
        return 0
    fi
    
    log "Validando conectividade com $host..."
    
    # Testar ping
    if ! ping -c 1 -W 5 "$host" &> /dev/null; then
        warn "Ping falhou para $host"
        return 1
    fi
    
    # Testar porta da API
    if ! nc -z -w 5 "$host" 8080 &> /dev/null; then
        warn "Porta 8080 não acessível em $host"
        return 1
    fi
    
    # Testar endpoint de saúde
    if ! curl -s -k --connect-timeout 10 "https://$host:8080/health" &> /dev/null; then
        warn "Endpoint de saúde não acessível em $host"
        return 1
    fi
    
    log "Conectividade validada com $host"
    return 0
}

# Função para configurar conexão com a rede
configure_network_connection() {
    local leader_host="$1"
    
    if [ "$leader_host" = "self" ]; then
        log "Configurando como primeiro nó da rede (líder)"
        
        # Criar configuração de rede para primeiro nó
        cat > /opt/syntropy/config/network-connection.yaml << EOF
network:
  role: "leader"
  leader_host: "self"
  mesh_config:
    enabled: true
    port: 51820
    interface: "wg0"
    subnet: "172.20.0.0/12"
    node_ip: "172.20.0.1/12"
  
  api_config:
    enabled: true
    port: 8080
    ssl: true
    cert_file: "/opt/syntropy/certs/node.crt"
    key_file: "/opt/syntropy/certs/node.key"
  
  discovery_config:
    enabled: true
    method: "dns"
    hostname: "syntropy-discovery.local"
    port: 8443
EOF
        
        # Configurar Wireguard para primeiro nó
        configure_wireguard_leader
        
    else
        log "Configurando conexão com líder: $leader_host"
        
        # Criar configuração de rede para nó worker
        cat > /opt/syntropy/config/network-connection.yaml << EOF
network:
  role: "worker"
  leader_host: "$leader_host"
  mesh_config:
    enabled: true
    port: 51820
    interface: "wg0"
    subnet: "172.20.0.0/12"
    node_ip: "172.20.0.2/12"
  
  api_config:
    enabled: true
    port: 8080
    ssl: true
    cert_file: "/opt/syntropy/certs/node.crt"
    key_file: "/opt/syntropy/certs/node.key"
  
  discovery_config:
    enabled: true
    method: "leader"
    leader_host: "$leader_host"
    port: 8443
EOF
        
        # Configurar Wireguard para nó worker
        configure_wireguard_worker "$leader_host"
    fi
    
    # Configurar permissões
    chown syntropy:syntropy /opt/syntropy/config/network-connection.yaml
    chmod 644 /opt/syntropy/config/network-connection.yaml
}

# Função para configurar Wireguard como líder
configure_wireguard_leader() {
    log "Configurando Wireguard como líder..."
    
    # Gerar chaves Wireguard
    wg genkey | tee /opt/syntropy/certs/wg-private.key | wg pubkey > /opt/syntropy/certs/wg-public.key
    chmod 600 /opt/syntropy/certs/wg-private.key
    chmod 644 /opt/syntropy/certs/wg-public.key
    chown syntropy:syntropy /opt/syntropy/certs/wg-*.key
    
    # Criar configuração Wireguard
    cat > /etc/wireguard/wg0.conf << EOF
[Interface]
PrivateKey = $(cat /opt/syntropy/certs/wg-private.key)
Address = 172.20.0.1/12
ListenPort = 51820
SaveConfig = true

# Permitir tráfego através da interface
PostUp = iptables -A FORWARD -i wg0 -j ACCEPT; iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
PostDown = iptables -D FORWARD -i wg0 -j ACCEPT; iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE
EOF
    
    # Iniciar Wireguard
    systemctl enable wg-quick@wg0
    systemctl start wg-quick@wg0
    
    log "Wireguard configurado como líder"
}

# Função para configurar Wireguard como worker
configure_wireguard_worker() {
    local leader_host="$1"
    
    log "Configurando Wireguard como worker para $leader_host..."
    
    # Gerar chaves Wireguard
    wg genkey | tee /opt/syntropy/certs/wg-private.key | wg pubkey > /opt/syntropy/certs/wg-public.key
    chmod 600 /opt/syntropy/certs/wg-private.key
    chmod 644 /opt/syntropy/certs/wg-public.key
    chown syntropy:syntropy /opt/syntropy/certs/wg-*.key
    
    # Obter chave pública do líder (implementação simplificada)
    local leader_public_key="placeholder_leader_public_key"
    
    # Criar configuração Wireguard
    cat > /etc/wireguard/wg0.conf << EOF
[Interface]
PrivateKey = $(cat /opt/syntropy/certs/wg-private.key)
Address = 172.20.0.2/12
ListenPort = 51820
SaveConfig = true

[Peer]
PublicKey = $leader_public_key
Endpoint = $leader_host:51820
AllowedIPs = 172.20.0.0/12
PersistentKeepalive = 25
EOF
    
    # Iniciar Wireguard
    systemctl enable wg-quick@wg0
    systemctl start wg-quick@wg0
    
    log "Wireguard configurado como worker"
}

# Função principal
main() {
    log "Iniciando processo de descoberta de rede..."
    
    # Tentar descobrir a rede
    local leader_host
    if leader_host=$(discover_network); then
        log "Rede descoberta com sucesso: $leader_host"
        
        # Validar conectividade
        if validate_connectivity "$leader_host"; then
            log "Conectividade validada com $leader_host"
            
            # Configurar conexão
            configure_network_connection "$leader_host"
            
            log "Configuração de rede concluída com sucesso!"
            
            # Salvar informações de descoberta
            echo "LEADER_HOST=$leader_host" > /opt/syntropy/config/network-discovery.env
            echo "DISCOVERY_SUCCESS=true" >> /opt/syntropy/config/network-discovery.env
            echo "DISCOVERY_TIMESTAMP=$(date -Iseconds)" >> /opt/syntropy/config/network-discovery.env
            
            chown syntropy:syntropy /opt/syntropy/config/network-discovery.env
            chmod 600 /opt/syntropy/config/network-discovery.env
            
            exit 0
        else
            error "Falha na validação de conectividade"
            exit 1
        fi
    else
        error "Falha na descoberta de rede"
        exit 1
    fi
}

# Executar função principal
main "$@"
