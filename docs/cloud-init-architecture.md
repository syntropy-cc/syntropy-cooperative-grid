# Syntropy Cooperative Grid - Arquitetura Cloud-Init

## Visão Geral

Este documento detalha a arquitetura completa do sistema de cloud-init implementado para o MVP do Syntropy Cooperative Grid. O sistema transforma o PC de trabalho em um "quartel general" que gera USBs bootáveis 100% plug-and-play para nós da rede Syntropy.

## Índice

1. [Arquitetura Geral](#arquitetura-geral)
2. [Componentes do Sistema](#componentes-do-sistema)
3. [Templates Cloud-Init](#templates-cloud-init)
4. [Scripts de Instalação](#scripts-de-instalação)
5. [Sistema de Segurança](#sistema-de-segurança)
6. [Processo de Descoberta](#processo-de-descoberta)
7. [Integração com CLI](#integração-com-cli)
8. [Fluxo de Funcionamento](#fluxo-de-funcionamento)
9. [Configurações e Personalização](#configurações-e-personalização)
10. [Troubleshooting](#troubleshooting)

---

## Arquitetura Geral

### Conceito Central

O sistema implementa uma arquitetura de **"Quartel General"** onde:

- **PC de Trabalho**: Centro de comando que gera USBs personalizados
- **USB Bootável**: DNA digital do nó com configurações específicas
- **Hardware Virgem**: Boot automático e configuração completa
- **Rede Syntropy**: Descoberta e conexão automática

### Diagrama de Arquitetura

```
┌─────────────────────────────────────────────────────────────┐
│                    PC DE TRABALHO                          │
│                   (Quartel General)                        │
│ ─────────────────────────────────────────────────────────── │
│ • CLI Go (syntropy usb create)                            │
│ • Geração de Certificados TLS                              │
│ • Geração de Chaves SSH                                    │
│ • Templates Cloud-Init                                     │
│ • Scripts de Instalação                                    │
│ • Estrutura ~/.syntropy/                                   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    USB BOOTÁVEL                            │
│                   (DNA do Nó)                              │
│ ─────────────────────────────────────────────────────────── │
│ • ISO Ubuntu Personalizada                                 │
│ • Cloud-Init Configurado                                   │
│ • Certificados TLS                                         │
│ • Chaves SSH                                               │
│ • Scripts de Instalação                                    │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                 HARDWARE VIRGEM                            │
│                   (Nó Syntropy)                            │
│ ─────────────────────────────────────────────────────────── │
│ • Boot Automático                                          │
│ • Detecção de Hardware                                     │
│ • Configuração de Rede                                     │
│ • Instalação do Syntropy Agent                             │
│ • Conexão ao Cluster                                       │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   REDE SYNTROPY                            │
│                 (Cooperative Grid)                         │
│ ─────────────────────────────────────────────────────────── │
│ • Descoberta Automática                                    │
│ • Mesh Network (Wireguard)                                 │
│ • Kubernetes Cluster                                       │
│ • Sistema de Créditos                                      │
│ • Monitoramento                                            │
└─────────────────────────────────────────────────────────────┘
```

---

## Componentes do Sistema

### 1. Templates Cloud-Init

#### user-data-template.yaml
**Propósito**: Configuração principal do sistema operacional

**Características**:
- Configuração de usuário e grupos
- Instalação de pacotes essenciais
- Configuração de serviços (Docker, SSH, firewall)
- Scripts de inicialização do Syntropy Agent

**Estrutura**:
```yaml
# Configurações básicas
locale: pt_BR.UTF-8
timezone: America/Sao_Paulo
hostname: ${NODE_NAME}

# Usuário Syntropy
users:
  - name: syntropy
    groups: [adm, sudo, docker]
    shell: /bin/bash
    sudo: ALL=(ALL) NOPASSWD:ALL
    lock_passwd: true
    ssh_authorized_keys:
      - ${SSH_PUBLIC_KEY}

# Pacotes instalados
packages:
  - curl, wget, git, htop, vim
  - docker.io, containerd, kubectl
  - wireguard, fail2ban, ufw
  - prometheus-node-exporter

# Comandos de inicialização
runcmd:
  - Configuração do Docker
  - Configuração do firewall
  - Criação de diretórios Syntropy
  - Download do Syntropy Agent
  - Configuração de certificados
  - Inicialização de serviços
```

#### meta-data-template.yaml
**Propósito**: Metadados do nó e configurações específicas

**Características**:
- Identificação única do nó
- Informações de hardware
- Configurações de rede
- Dados de segurança

**Estrutura**:
```yaml
instance-id: ${NODE_NAME}-${INSTANCE_ID}
local-hostname: ${NODE_NAME}

syntropy:
  node:
    name: ${NODE_NAME}
    type: ${NODE_TYPE}
    description: ${NODE_DESCRIPTION}
    coordinates: ${COORDINATES}
    owner_key: ${OWNER_KEY}
  
  hardware:
    detected_type: ${DETECTED_HARDWARE_TYPE}
    cpu_cores: ${CPU_CORES}
    memory_gb: ${MEMORY_GB}
    storage_gb: ${STORAGE_GB}
  
  role:
    initial_role: ${INITIAL_ROLE}
    can_be_leader: ${CAN_BE_LEADER}
    can_be_worker: ${CAN_BE_WORKER}
```

#### network-config-template.yaml
**Propósito**: Configuração avançada de rede

**Características**:
- Suporte a múltiplas interfaces
- Configuração de bridges e VLANs
- Roteamento estático
- Configuração de proxy

**Estrutura**:
```yaml
version: 2
ethernets:
  en*:
    dhcp4: true
    dhcp6: false
    dhcp4-overrides:
      hostname: ${NODE_NAME}
  
bridges:
  br0:
    interfaces: []
    addresses:
      - 172.20.0.${NODE_IP_SUFFIX}/24
    gateway4: 172.20.0.1

vlans:
  vlan100:
    id: 100
    link: ${PRIMARY_INTERFACE}
    addresses:
      - 192.168.100.${NODE_IP_SUFFIX}/24
```

### 2. Scripts de Instalação

#### hardware-detection.sh
**Propósito**: Detecção automática de hardware e classificação

**Funcionalidades**:
- Detecção de CPU, RAM, storage
- Classificação do tipo de hardware
- Determinação de papéis (líder/worker)
- Configuração de capacidades

**Algoritmo de Classificação**:
```bash
# Servidor Dedicado
if [ CPU_CORES >= 16 ] && [ MEMORY_GB >= 32 ] && [ STORAGE_GB >= 1000 ]; then
    HARDWARE_TYPE="server"
    CAN_BE_LEADER="true"
    INITIAL_ROLE="leader"

# Home Server
elif [ CPU_CORES >= 4 ] && [ MEMORY_GB >= 8 ] && [ STORAGE_GB >= 500 ]; then
    HARDWARE_TYPE="home_server"
    CAN_BE_LEADER="true"
    INITIAL_ROLE="worker"

# Computador Pessoal
elif [ CPU_CORES >= 2 ] && [ MEMORY_GB >= 4 ] && [ STORAGE_GB >= 100 ]; then
    HARDWARE_TYPE="personal_computer"
    CAN_BE_LEADER="false"
    INITIAL_ROLE="worker"

# Mobile/IoT
else
    HARDWARE_TYPE="mobile_iot"
    CAN_BE_LEADER="false"
    INITIAL_ROLE="worker"
fi
```

#### network-discovery.sh
**Propósito**: Descoberta inteligente da rede Syntropy

**Métodos de Descoberta**:
1. **DNS**: Resolve hostnames como `syntropy-discovery.local`
2. **Broadcast**: Envia broadcast na rede local
3. **Multicast**: Usa multicast para descoberta
4. **Configuração Manual**: Usa hosts pré-configurados

**Algoritmo de Descoberta**:
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

#### syntropy-install.sh
**Propósito**: Instalação completa do Syntropy Agent

**Funcionalidades**:
- Atualização do sistema
- Instalação de dependências
- Configuração do Docker
- Instalação do Kubernetes
- Configuração de segurança
- Instalação do Syntropy Agent

**Processo de Instalação**:
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

#### cluster-join.sh
**Propósito**: Conexão ao cluster Syntropy

**Funcionalidades**:
- Obtenção de informações do cluster
- Registro do nó
- Configuração do Kubernetes
- Configuração do mesh network
- Verificação de conectividade

**Processo de Conexão**:
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

---

## Sistema de Segurança

### 1. Geração de Certificados

#### Certificados TLS
**Propósito**: Comunicação segura entre nós

**Estrutura**:
- **CA (Certificate Authority)**: Certificado raiz válido por 10 anos
- **Node Certificate**: Certificado do nó válido por 1 ano
- **Algoritmo**: RSA 4096 bits para CA, RSA 2048 bits para nós

**Geração**:
```go
// Gerar chave CA
caKey, err := rsa.GenerateKey(rand.Reader, 4096)

// Criar certificado CA
caTemplate := x509.Certificate{
    SerialNumber: big.NewInt(1),
    Subject: pkix.Name{
        Organization: []string{"Syntropy Cooperative Grid"},
        Country:      []string{"BR"},
        Locality:     []string{"São Paulo"},
    },
    NotBefore:             time.Now(),
    NotAfter:              time.Now().AddDate(10, 0, 0),
    KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
    ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
    BasicConstraintsValid: true,
    IsCA:                  true,
}

// Gerar certificado do nó
nodeTemplate := x509.Certificate{
    SerialNumber: big.NewInt(2),
    Subject: pkix.Name{
        CommonName: nodeName,
    },
    NotBefore:   time.Now(),
    NotAfter:    time.Now().AddDate(1, 0, 0),
    KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
    ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
    IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1)},
    DNSNames:    []string{nodeName, "localhost"},
}
```

#### Chaves SSH
**Propósito**: Acesso seguro ao nó

**Características**:
- Algoritmo: RSA 2048 bits
- Formato: PEM
- Usuário: syntropy
- Acesso: Apenas por chave pública

**Geração**:
```go
// Gerar chave privada RSA
privateKey, err := rsa.GenerateKey(rand.Reader, 2048)

// Codificar em PEM
privateKeyPEM := &pem.Block{
    Type:  "RSA PRIVATE KEY",
    Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
}

// Gerar chave pública
publicKey, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
publicKeyPEM := &pem.Block{
    Type:  "PUBLIC KEY",
    Bytes: publicKey,
}
```

### 2. Configuração de Firewall

#### Regras UFW
**Propósito**: Proteção do nó

**Regras Implementadas**:
```bash
# Política padrão
ufw default deny incoming
ufw default allow outgoing

# Serviços essenciais
ufw allow ssh
ufw allow 6443/tcp    # Kubernetes API
ufw allow 2379:2380/tcp  # etcd
ufw allow 10250/tcp   # Kubelet
ufw allow 10251/tcp   # Kube-scheduler
ufw allow 10252/tcp   # Kube-controller-manager
ufw allow 30000:32767/tcp  # NodePort services

# Serviços Syntropy
ufw allow 51820/udp   # Wireguard
ufw allow 8080/tcp    # Syntropy API
ufw allow 9090/tcp    # Metrics
ufw allow 9100/tcp    # Node Exporter
```

### 3. Fail2ban

#### Configuração
**Propósito**: Proteção contra ataques de força bruta

**Regras**:
- SSH: 3 tentativas em 10 minutos
- HTTP: 5 tentativas em 10 minutos
- Ban: 1 hora

---

## Processo de Descoberta

### 1. Descoberta via DNS

#### Hostnames Suportados
- `syntropy-discovery.local`
- `syntropy-mesh.local`
- `syntropy-leader.local`
- `syntropy.local`

#### Processo
```bash
discover_via_dns() {
    local dns_hostnames=(
        "syntropy-discovery.local"
        "syntropy-mesh.local"
        "syntropy-leader.local"
        "syntropy.local"
    )
    
    for hostname in "${dns_hostnames[@]}"; do
        if nslookup "$hostname" &> /dev/null; then
            local ip=$(nslookup "$hostname" | awk '/^Address: / { print $2 }' | tail -1)
            if [ -n "$ip" ] && [ "$ip" != "127.0.0.1" ]; then
                if ping -c 1 -W 2 "$ip" &> /dev/null; then
                    echo "$ip"
                    return 0
                fi
            fi
        fi
    done
    
    return 1
}
```

### 2. Descoberta via Broadcast

#### Processo
```bash
discover_via_broadcast() {
    # Determinar sub-rede local
    local subnet=$(ip route | grep "$network" | awk '{print $1}' | head -1)
    
    # Usar nmap para descobrir hosts ativos na porta 8080
    nmap -p 8080 --open "$subnet" -oG "$temp_file"
    
    # Testar cada host descoberto
    for host in $hosts; do
        if curl -s -k --connect-timeout 5 "https://$host:8080/health" &> /dev/null; then
            echo "$host"
            return 0
        fi
    done
    
    return 1
}
```

### 3. Descoberta via Multicast

#### Configuração
- **Grupo Multicast**: 224.0.0.251
- **Porta**: 5353
- **Mensagem**: "SYNTROPY_DISCOVERY_REQUEST"

#### Processo
```bash
discover_via_multicast() {
    # Enviar requisição via multicast
    echo "SYNTROPY_DISCOVERY_REQUEST" | nc -u -w 2 "$multicast_group" "$multicast_port"
    
    # Aguardar respostas
    sleep 2
    
    # Processar respostas recebidas
    if [ -f "$temp_file" ] && [ -s "$temp_file" ]; then
        local host=$(head -1 "$temp_file")
        echo "$host"
        return 0
    fi
    
    return 1
}
```

---

## Integração com CLI

### 1. Comando USB Create

#### Sintaxe
```bash
syntropy usb create [device] [flags]
```

#### Flags Disponíveis
```bash
--node-name string          Nome do nó (obrigatório)
--description string        Descrição do nó
--coordinates string        Coordenadas geográficas (lat,lon)
--owner-key string          Arquivo de chave de proprietário existente
--auto-detect              Detectar automaticamente dispositivo USB
--label string             Rótulo do sistema de arquivos (padrão: SYNTROPY)
--work-dir string          Diretório de trabalho (padrão: /tmp/syntropy-work)
--cache-dir string         Diretório de cache (padrão: ~/.syntropy/cache)
--iso string               Caminho para ISO Ubuntu (baixa automaticamente se não especificado)
--discovery-server string  Servidor de descoberta da rede (padrão: syntropy-discovery.local)
--created-by string        Usuário que criou o nó (padrão: usuário atual)
```

#### Exemplos de Uso
```bash
# Criar USB com auto-detecção
syntropy usb create --auto-detect --node-name "node-01"

# Criar USB especificando dispositivo
syntropy usb create /dev/sdb --node-name "node-02" --description "Servidor principal"

# Criar USB com coordenadas
syntropy usb create --auto-detect --node-name "node-03" --coordinates "-23.5505,-46.6333"

# Criar USB com servidor de descoberta personalizado
syntropy usb create --auto-detect --node-name "node-04" --discovery-server "192.168.1.100"
```

### 2. Estrutura de Arquivos

#### PC de Trabalho
```
~/.syntropy/
├── backups/          # Backups de configurações
├── cache/           # Cache de ISOs e downloads
├── config/          # Configurações globais
├── diagnostics/     # Logs de diagnóstico
├── keys/            # Chaves SSH e certificados
├── logs/            # Logs da CLI
├── nodes/           # Configurações dos nós
└── scripts/         # Scripts auxiliares
```

#### USB Bootável
```
/
├── cloud-init/
│   ├── user-data
│   ├── meta-data
│   └── network-config
├── scripts/
│   ├── hardware-detection.sh
│   ├── network-discovery.sh
│   ├── syntropy-install.sh
│   └── cluster-join.sh
└── certs/
    ├── ca.crt
    ├── ca.key
    ├── node.crt
    └── node.key
```

### 3. Processo de Criação

#### Fluxo de Execução
```go
func createUSB(devicePath string, config *Config, workDir, cacheDir string) error {
    // 1. Configurar diretórios
    os.MkdirAll(workDir, 0755)
    os.MkdirAll(cacheDir, 0755)

    // 2. Gerar chaves SSH se não fornecidas
    if config.SSHPublicKey == "" || config.SSHPrivateKey == "" {
        privateKey, publicKey, err := generateSSHKeyPair()
        config.SSHPrivateKey = privateKey
        config.SSHPublicKey = publicKey
    }

    // 3. Gerar certificados TLS
    certs, err := generateCertificates(config.NodeName, config.OwnerKeyFile)

    // 4. Salvar certificados
    caKeyPath, caCertPath, nodeKeyPath, nodeCertPath, err := saveCertificates(certs, workDir)

    // 5. Gerar configuração do cloud-init
    cloudInitConfig, err := generateCloudInitConfig(config, workDir)

    // 6. Criar arquivos do cloud-init
    createCloudInitFiles(cloudInitConfig, workDir, certPaths)

    // 7. Copiar scripts
    copyScripts(workDir)

    // 8. Criar ISO personalizada
    isoPath, err := createCustomISO(workDir, config)

    // 9. Criar USB com a ISO
    return createUSBWithISO(devicePath, isoPath)
}
```

---

## Fluxo de Funcionamento

### 1. Preparação no PC de Trabalho

#### Etapas
1. **Configuração**: Usuário define parâmetros do nó
2. **Geração de Chaves**: Sistema gera chaves SSH e certificados TLS
3. **Templates**: Geração dos arquivos cloud-init personalizados
4. **Scripts**: Cópia dos scripts de instalação
5. **ISO**: Criação da ISO Ubuntu personalizada
6. **USB**: Gravação da ISO no dispositivo USB

#### Comando Executado
```bash
syntropy usb create --auto-detect --node-name "node-01" --description "Servidor principal"
```

#### Saída Esperada
```
🚀 Iniciando criação de USB para nó: node-01
📍 Plataforma: linux
💾 Dispositivo: /dev/sdb
📂 Diretório de trabalho: /tmp/syntropy-work
📂 Diretório de cache: ~/.syntropy/cache

🔑 Gerando par de chaves SSH...
✅ Chaves SSH geradas com sucesso

🔐 Gerando certificados TLS...
✅ Certificados TLS gerados com sucesso

📝 Gerando configuração do cloud-init...
✅ Configuração do cloud-init criada com sucesso

📜 Copiando scripts de instalação...
✅ Scripts copiados com sucesso

💿 Criando ISO personalizada...
📥 Baixando ISO Ubuntu base...
🔗 Montando ISO base...
📋 Copiando conteúdo da ISO base...
📝 Copiando arquivos do cloud-init...
📜 Copiando scripts...
🔐 Copiando certificados...
🔓 Desmontando ISO base...
💿 Criando ISO personalizada...
✅ ISO personalizada criada: /tmp/syntropy-work/syntropy-node-01.iso

✅ USB criado com sucesso!
```

### 2. Boot do Hardware Virgem

#### Etapas
1. **Boot**: Hardware inicia pelo USB
2. **Cloud-Init**: Sistema executa configurações
3. **Detecção**: Script detecta tipo de hardware
4. **Instalação**: Instalação do Syntropy Agent
5. **Descoberta**: Descoberta da rede Syntropy
6. **Conexão**: Conexão ao cluster

#### Logs de Boot
```
[   0.000] Boot iniciado
[   5.234] Cloud-init iniciado
[   8.456] Usuário syntropy criado
[  12.789] Pacotes instalados
[  18.234] Docker configurado
[  22.567] Firewall configurado
[  25.890] Syntropy Agent baixado
[  28.123] Certificados instalados
[  30.456] Hardware detectado: server
[  32.789] Rede descoberta: 192.168.1.100
[  35.012] Cluster conectado
[  37.345] Nó operacional
```

### 3. Operação na Rede

#### Status do Nó
```bash
# Verificar status do Syntropy Agent
systemctl status syntropy-agent

# Verificar logs
journalctl -u syntropy-agent -f

# Verificar conectividade
curl -k https://localhost:8080/health

# Verificar métricas
curl http://localhost:9090/metrics
```

#### Comandos de Gerenciamento
```bash
# Reiniciar agente
sudo systemctl restart syntropy-agent

# Verificar configuração
cat /opt/syntropy/config/agent.yaml

# Verificar certificados
openssl x509 -in /opt/syntropy/certs/node.crt -text -noout

# Verificar conectividade de rede
ping syntropy-discovery.local
```

---

## Configurações e Personalização

### 1. Variáveis de Ambiente

#### Template Variables
```yaml
# Identificação do nó
NODE_NAME: "node-01"
NODE_DESCRIPTION: "Servidor principal"
COORDINATES: "-23.5505,-46.6333"

# Segurança
OWNER_KEY: "owner_key_here"
SSH_PUBLIC_KEY: "ssh-rsa AAAAB3NzaC1yc2E..."
SSH_PRIVATE_KEY: "-----BEGIN RSA PRIVATE KEY-----"

# Rede
DISCOVERY_SERVER: "syntropy-discovery.local"
INTERFACE: "eth0"
GATEWAY: "192.168.1.1"
NODE_IP_SUFFIX: "2"

# Hardware
DETECTED_HARDWARE_TYPE: "server"
CPU_CORES: "16"
MEMORY_GB: "32"
STORAGE_GB: "1000"

# Papel
INITIAL_ROLE: "leader"
CAN_BE_LEADER: "true"
CAN_BE_WORKER: "true"
```

### 2. Personalização de Scripts

#### Modificando Detecção de Hardware
```bash
# Adicionar novo tipo de hardware
if [ "$CPU_CORES" -ge 32 ] && [ "$MEMORY_GB" -ge 64 ]; then
    HARDWARE_TYPE="super_server"
    CAN_BE_LEADER="true"
    INITIAL_ROLE="leader"
    log "Hardware detectado: SUPER SERVER (Capacidade máxima)"
```

#### Modificando Descoberta de Rede
```bash
# Adicionar novo método de descoberta
discover_via_zeroconf() {
    log "Tentando descoberta via Zeroconf..."
    
    # Implementar descoberta Zeroconf
    if command -v avahi-browse &> /dev/null; then
        local hosts=$(avahi-browse -t _syntropy._tcp | grep -o '[0-9.]*')
        for host in $hosts; do
            if ping -c 1 -W 2 "$host" &> /dev/null; then
                echo "$host"
                return 0
            fi
        done
    fi
    
    return 1
}
```

### 3. Configurações Avançadas

#### Personalizando Firewall
```bash
# Adicionar regras específicas
ufw allow from 192.168.100.0/24 to any port 8080
ufw allow from 172.20.0.0/12 to any port 51820
ufw deny from 10.0.0.0/8
```

#### Personalizando Kubernetes
```bash
# Configurar recursos específicos
kubectl taint nodes node-01 node-role.kubernetes.io/control-plane:NoSchedule
kubectl label nodes node-01 node-type=server
kubectl label nodes node-01 location=datacenter-1
```

---

## Troubleshooting

### 1. Problemas Comuns

#### USB Não Boota
**Sintomas**: Hardware não inicia pelo USB

**Diagnóstico**:
```bash
# Verificar se USB é bootável
file /dev/sdb
fdisk -l /dev/sdb

# Verificar ISO
file /tmp/syntropy-work/syntropy-node-01.iso
```

**Soluções**:
1. Verificar se hardware suporta boot USB
2. Verificar se USB está formatado corretamente
3. Verificar se ISO foi gravada corretamente
4. Tentar outro dispositivo USB

#### Cloud-Init Falha
**Sintomas**: Sistema não configura automaticamente

**Diagnóstico**:
```bash
# Verificar logs do cloud-init
journalctl -u cloud-init
tail -f /var/log/cloud-init.log

# Verificar arquivos de configuração
ls -la /var/lib/cloud/seed/nocloud/
cat /var/lib/cloud/seed/nocloud/user-data
```

**Soluções**:
1. Verificar sintaxe dos arquivos YAML
2. Verificar se variáveis estão definidas
3. Verificar se arquivos estão no local correto
4. Verificar permissões dos arquivos

#### Falha na Descoberta de Rede
**Sintomas**: Nó não encontra a rede Syntropy

**Diagnóstico**:
```bash
# Verificar conectividade de rede
ping 8.8.8.8
nslookup syntropy-discovery.local

# Verificar portas
nmap -p 8080 192.168.1.100
telnet 192.168.1.100 8080

# Verificar logs de descoberta
tail -f /opt/syntropy/logs/network-discovery.log
```

**Soluções**:
1. Verificar se rede está configurada corretamente
2. Verificar se servidor de descoberta está ativo
3. Verificar se firewall permite comunicação
4. Tentar configuração manual

#### Falha na Instalação do Syntropy Agent
**Sintomas**: Syntropy Agent não inicia

**Diagnóstico**:
```bash
# Verificar status do serviço
systemctl status syntropy-agent

# Verificar logs
journalctl -u syntropy-agent -f

# Verificar configuração
cat /opt/syntropy/config/agent.yaml

# Verificar certificados
openssl x509 -in /opt/syntropy/certs/node.crt -text -noout
```

**Soluções**:
1. Verificar se certificados estão válidos
2. Verificar se configuração está correta
3. Verificar se dependências estão instaladas
4. Verificar se portas estão disponíveis

### 2. Logs e Diagnóstico

#### Localização dos Logs
```bash
# Logs do sistema
/var/log/syslog
/var/log/cloud-init.log

# Logs do Syntropy
/opt/syntropy/logs/agent.log
/opt/syntropy/logs/hardware-detection.log
/opt/syntropy/logs/network-discovery.log
/opt/syntropy/logs/installation.log
/opt/syntropy/logs/cluster-join.log

# Logs de auditoria
/opt/syntropy/audit/audit.log
```

#### Comandos de Diagnóstico
```bash
# Verificar status geral
systemctl status syntropy-agent docker ssh

# Verificar conectividade
ping syntropy-discovery.local
curl -k https://localhost:8080/health

# Verificar recursos
htop
df -h
free -h

# Verificar rede
ip addr show
ip route show
ss -tulpn
```

### 3. Recuperação e Manutenção

#### Backup de Configuração
```bash
# Criar backup
tar -czf /opt/syntropy/backups/config-backup-$(date +%Y%m%d).tar.gz \
    /opt/syntropy/config \
    /opt/syntropy/certs \
    /etc/wireguard

# Restaurar backup
tar -xzf /opt/syntropy/backups/config-backup-20240101.tar.gz -C /
```

#### Atualização do Sistema
```bash
# Atualizar pacotes
apt update && apt upgrade -y

# Atualizar Syntropy Agent
curl -L https://github.com/syntropy-cooperative-grid/agent/releases/latest/download/syntropy-agent-linux-amd64 -o /opt/syntropy/bin/syntropy-agent
chmod +x /opt/syntropy/bin/syntropy-agent
systemctl restart syntropy-agent
```

#### Limpeza de Logs
```bash
# Limpar logs antigos
find /opt/syntropy/logs -name "*.log" -mtime +30 -delete
find /opt/syntropy/audit -name "*.log" -mtime +90 -delete

# Rotacionar logs
logrotate -f /etc/logrotate.d/syntropy
```

---

## Conclusão

A arquitetura de cloud-init implementada para o Syntropy Cooperative Grid oferece:

1. **Automação Completa**: Boot para operação em minutos
2. **Segurança Robusta**: Certificados TLS e chaves SSH automáticas
3. **Descoberta Inteligente**: Múltiplos métodos de descoberta de rede
4. **Flexibilidade**: Suporte a diferentes tipos de hardware
5. **Auditoria**: Logs completos de todas as operações
6. **Manutenibilidade**: Sistema de backup e recuperação

Esta implementação estabelece uma base sólida para o MVP do Syntropy Cooperative Grid, permitindo que usuários criem e gerenciem nós da rede de forma simples e segura, enquanto mantém a visão de descentralização para o futuro.
