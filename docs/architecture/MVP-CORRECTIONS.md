# MVP Corrections & Improvements

**Data**: Outubro 2025  
**Versão**: 2.1  
**Objetivo**: Corrigir inconsistências críticas identificadas na análise

---

## 🔐 1. CORREÇÃO: SEGURANÇA DO GRID TOKEN

### 1.1 Problema
Armazenar Grid Token em `~/.syntropy/grid-token.txt` é inseguro:
- Texto plano no disco
- Permissões 600 não impedem root
- Backups podem vazar
- Git pode committar acidentalmente

### 1.2 Solução: Sistema de Keyring

**Implementação**: `manager/interfaces/cli/setup/src/token_manager.go`

```go
package setup

import (
	"crypto/rand"
	"fmt"
	
	"github.com/zalando/go-keyring"
)

const (
	KeyringService = "syntropy-grid"
	KeyringUser    = "grid-token"
)

// TokenManager gerencia Grid Token de forma segura
type TokenManager struct{}

// NewTokenManager cria novo gerenciador de tokens
func NewTokenManager() *TokenManager {
	return &TokenManager{}
}

// GenerateToken gera novo Grid Token (UUID v4)
func (tm *TokenManager) GenerateToken() (string, error) {
	tokenBytes := make([]byte, 16)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	
	token := fmt.Sprintf("%x-%x-%x-%x-%x",
		tokenBytes[0:4],
		tokenBytes[4:6],
		tokenBytes[6:8],
		tokenBytes[8:10],
		tokenBytes[10:16],
	)
	
	return token, nil
}

// SaveToken salva token no keyring do sistema
func (tm *TokenManager) SaveToken(token string) error {
	if err := keyring.Set(KeyringService, KeyringUser, token); err != nil {
		return fmt.Errorf("failed to save token to keyring: %w", err)
	}
	
	fmt.Println("✅ Grid Token saved securely in system keyring")
	fmt.Printf("   Service: %s\n", KeyringService)
	fmt.Printf("   Token: %s...[HIDDEN]\n", token[:8])
	
	return nil
}

// LoadToken carrega token do keyring
func (tm *TokenManager) LoadToken() (string, error) {
	token, err := keyring.Get(KeyringService, KeyringUser)
	if err != nil {
		return "", fmt.Errorf("failed to load token from keyring: %w", err)
	}
	
	return token, nil
}

// DeleteToken remove token do keyring
func (tm *TokenManager) DeleteToken() error {
	if err := keyring.Delete(KeyringService, KeyringUser); err != nil {
		return fmt.Errorf("failed to delete token from keyring: %w", err)
	}
	
	return nil
}

// ExportToken exporta token para arquivo temporário (com confirmação)
// Usado apenas para recovery ou migração
func (tm *TokenManager) ExportToken(outputPath string) error {
	fmt.Println("⚠️  WARNING: Exporting Grid Token to file!")
	fmt.Println("   This should ONLY be done for backup/recovery")
	fmt.Print("   Type 'EXPORT' to confirm: ")
	
	var confirm string
	fmt.Scanln(&confirm)
	
	if confirm != "EXPORT" {
		return fmt.Errorf("export cancelled")
	}
	
	token, err := tm.LoadToken()
	if err != nil {
		return err
	}
	
	// Salvar com permissões muito restritivas
	file, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0400)
	if err != nil {
		return fmt.Errorf("failed to create export file: %w", err)
	}
	defer file.Close()
	
	if _, err := file.WriteString(token); err != nil {
		return fmt.Errorf("failed to write token: %w", err)
	}
	
	fmt.Printf("✅ Token exported to: %s (READ-ONLY)\n", outputPath)
	fmt.Println("⚠️  DELETE this file after backup!")
	
	return nil
}
```

### 1.3 Integração no Setup

**Atualizar**: `manager/interfaces/cli/setup/src/setup.go`

```go
func (s *Setup) Run() error {
	// ... existing code ...
	
	// Gerar e salvar Grid Token SEGURAMENTE
	log.Println("🔐 Generating Grid Token...")
	tokenMgr := NewTokenManager()
	
	token, err := tokenMgr.GenerateToken()
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}
	
	// Salvar no keyring do sistema (NÃO em arquivo)
	if err := tokenMgr.SaveToken(token); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}
	
	log.Println("✅ Grid Token saved securely")
	log.Printf("   Location: System Keyring (%s)\n", runtime.GOOS)
	log.Printf("   Service: syntropy-grid")
	
	return nil
}
```

### 1.4 Dependências

**Adicionar ao** `go.mod`:
```go
require (
    github.com/zalando/go-keyring v0.2.3
)
```

**Instalação de dependências do sistema**:

**Linux**:
```bash
# Ubuntu/Debian
sudo apt-get install libsecret-1-dev

# Fedora/RHEL
sudo dnf install libsecret-devel
```

**Windows**: Nenhuma dependência (usa Credential Manager nativo)

**macOS**: Nenhuma dependência (usa Keychain nativo)

### 1.5 Comandos da CLI

**Adicionar comandos** para gerenciar token:

```bash
# Ver token (com confirmação de segurança)
syntropy token show

# Exportar token (backup)
syntropy token export --output backup-token.txt

# Importar token (recovery)
syntropy token import --file backup-token.txt

# Rotacionar token (gerar novo)
syntropy token rotate
```

---

## 📋 2. CORREÇÃO: TEMPLATES CLOUD-INIT

### 2.1 Problemas Identificados

1. ❌ Falta `${GRID_TOKEN}`
2. ❌ Falta script de auto-detecção de hardware
3. ❌ Falta Registration Protocol
4. ❌ Agent não existe ainda (template tenta baixar do GitHub)
5. ❌ Variáveis inconsistentes com MVP
6. ❌ Sobre-complexidade no network-config

### 2.2 Template Corrigido: user-data

**Criar**: `infrastructure/cloud-init/user-data-mvp.yaml`

```yaml
#cloud-config
# Syntropy Cooperative Grid - MVP Node Provisioning
# Generated by Command Station

locale: pt_BR.UTF-8
timezone: America/Sao_Paulo
hostname: syntropy-${NODE_NAME}

# User configuration
users:
  - name: syntropy
    groups: [adm, sudo, docker]
    shell: /bin/bash
    sudo: ALL=(ALL) NOPASSWD:ALL
    lock_passwd: true
    ssh_authorized_keys:
      - ${SSH_PUBLIC_KEY}

# SSH configuration
ssh_pwauth: false
disable_root: true

# Essential packages only (MVP)
packages:
  - curl
  - wget
  - docker.io
  - docker-compose-plugin
  - fail2ban
  - ufw
  - jq
  - prometheus-node-exporter

runcmd:
  # Docker setup
  - systemctl enable docker
  - systemctl start docker
  - usermod -aG docker syntropy
  
  # Firewall (MVP essentials only)
  - ufw --force enable
  - ufw default deny incoming
  - ufw default allow outgoing
  - ufw allow from ${COMMAND_STATION_IP} to any port 22
  - ufw allow 51000/tcp  # Registration Protocol
  - ufw allow 8080/tcp   # Agent API
  - ufw allow 9100/tcp   # Prometheus
  
  # Fail2ban
  - systemctl enable fail2ban
  - systemctl start fail2ban
  
  # Create Syntropy directories
  - mkdir -p /opt/syntropy/{bin,config,logs,metadata}
  - chown -R syntropy:syntropy /opt/syntropy
  
  # Copy Agent from USB (NOT from GitHub - doesn't exist yet)
  - |
    if [ -f /media/cdrom/syntropy/agent ]; then
      cp /media/cdrom/syntropy/agent /opt/syntropy/bin/syntropy-agent
      chmod +x /opt/syntropy/bin/syntropy-agent
      echo "✅ Agent copied from USB"
    else
      echo "⚠️  Agent not found on USB - will use placeholder"
      cat > /opt/syntropy/bin/syntropy-agent << 'AGENT_PLACEHOLDER'
#!/bin/bash
# Syntropy Agent Placeholder
# TODO: Implement actual agent

case "$1" in
  status)
    echo '{"status":"online","node":"${NODE_NAME}"}'
    ;;
  *)
    echo "Syntropy Agent Placeholder - No-op"
    ;;
esac
AGENT_PLACEHOLDER
      chmod +x /opt/syntropy/bin/syntropy-agent
    fi
  
  # Hardware auto-detection script
  - |
    cat > /opt/syntropy/bin/detect-hardware << 'DETECT_SCRIPT'
#!/bin/bash
# Auto-detect hardware and generate manifest

set -e

echo "🔍 Detecting hardware..."

# CPU detection
CPU_CORES=$(nproc)
CPU_MODEL=$(lscpu | grep "Model name" | sed 's/Model name:\s*//' | xargs)
CPU_FREQ=$(lscpu | grep "CPU MHz" | head -1 | awk '{print $3}' | cut -d. -f1)

# RAM detection
RAM_BYTES=$(free -b | awk '/^Mem:/{print $2}')
RAM_GB=$(echo "scale=2; $RAM_BYTES / 1024 / 1024 / 1024" | bc)

# Disk detection
DISK_DEVICE=$(df / | awk 'NR==2{print $1}' | sed 's/[0-9]*$//')
DISK_BYTES=$(lsblk -b -d -o SIZE -n $DISK_DEVICE 2>/dev/null || echo "0")
DISK_GB=$(echo "scale=0; $DISK_BYTES / 1024 / 1024 / 1024" | bc)

# Disk type (NVMe, SSD, HDD)
if [[ $DISK_DEVICE == *"nvme"* ]]; then
    DISK_TYPE="nvme"
elif [ -f /sys/block/$(basename $DISK_DEVICE)/queue/rotational ]; then
    ROTATIONAL=$(cat /sys/block/$(basename $DISK_DEVICE)/queue/rotational)
    if [ "$ROTATIONAL" = "0" ]; then
        DISK_TYPE="ssd"
    else
        DISK_TYPE="hdd"
    fi
else
    DISK_TYPE="unknown"
fi

# Network interfaces
INTERFACES=$(ip -j link show | jq -c '[.[] | select(.link_type=="ether" and .ifname!="docker0") | {name:.ifname, mac:.address, state:.operstate}]')

# Current IP
CURRENT_IP=$(ip -4 addr show | grep -oP '(?<=inet\s)\d+(\.\d+){3}' | grep -v '127.0.0.1' | head -1)

# System info
KERNEL=$(uname -r)
OS_VERSION=$(lsb_release -d | cut -f2)

# Timestamp
TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Generate Hardware Manifest JSON
cat > /opt/syntropy/metadata/hardware-manifest.json << MANIFEST
{
  "version": "1.0",
  "detected_at": "$TIMESTAMP",
  "node_name": "${NODE_NAME}",
  "cpu": {
    "cores": $CPU_CORES,
    "model": "$CPU_MODEL",
    "frequency_mhz": $CPU_FREQ
  },
  "memory": {
    "total_gb": $RAM_GB,
    "total_bytes": $RAM_BYTES
  },
  "disk": {
    "total_gb": $DISK_GB,
    "type": "$DISK_TYPE",
    "device": "$DISK_DEVICE"
  },
  "network": {
    "interfaces": $INTERFACES,
    "current_ip": "$CURRENT_IP"
  },
  "system": {
    "kernel": "$KERNEL",
    "os": "$OS_VERSION"
  }
}
MANIFEST

chmod 644 /opt/syntropy/metadata/hardware-manifest.json
chown syntropy:syntropy /opt/syntropy/metadata/hardware-manifest.json

echo "✅ Hardware manifest generated"
cat /opt/syntropy/metadata/hardware-manifest.json
DETECT_SCRIPT
  
  - chmod +x /opt/syntropy/bin/detect-hardware
  
  # Run hardware detection
  - /opt/syntropy/bin/detect-hardware
  
  # Agent configuration
  - |
    cat > /opt/syntropy/config/agent.yaml << 'AGENT_CONFIG'
node:
  name: "${NODE_NAME}"
  grid_token: "${GRID_TOKEN}"

command_station:
  ip: "${COMMAND_STATION_IP}"
  port: 51000
  registration_timeout: 300s

hardware:
  manifest_path: "/opt/syntropy/metadata/hardware-manifest.json"
  auto_update: true
  update_interval: 300s

logging:
  level: "info"
  file: "/opt/syntropy/logs/agent.log"

metrics:
  enabled: true
  port: 9100
AGENT_CONFIG
  
  # Registration script (announces to Command Station)
  - |
    cat > /opt/syntropy/bin/register-node << 'REGISTER_SCRIPT'
#!/bin/bash
# Registration Protocol - Announce to Command Station

set -e

CONFIG="/opt/syntropy/config/agent.yaml"
MANIFEST="/opt/syntropy/metadata/hardware-manifest.json"

# Parse config
NODE_NAME=$(grep "name:" $CONFIG | awk '{print $2}' | tr -d '"')
GRID_TOKEN=$(grep "grid_token:" $CONFIG | awk '{print $2}' | tr -d '"')
STATION_IP=$(grep "ip:" $CONFIG | awk '{print $2}' | tr -d '"')
STATION_PORT=$(grep "port:" $CONFIG | awk '{print $2}')

# Get current IP
CURRENT_IP=$(ip -4 addr show | grep -oP '(?<=inet\s)\d+(\.\d+){3}' | grep -v '127.0.0.1' | head -1)

# Load hardware manifest
HARDWARE=$(cat $MANIFEST)

# Build announcement JSON
ANNOUNCEMENT=$(cat << ANNOUNCE
{
  "type": "node_announcement",
  "node_name": "$NODE_NAME",
  "grid_token": "$GRID_TOKEN",
  "ip": "$CURRENT_IP",
  "ssh_port": 22,
  "public_key": "$(cat /home/syntropy/.ssh/authorized_keys | head -1)",
  "hardware_manifest": $HARDWARE,
  "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ)"
}
ANNOUNCE
)

echo "📢 Announcing to Command Station..."
echo "   Station: $STATION_IP:$STATION_PORT"
echo "   Node: $NODE_NAME"
echo "   IP: $CURRENT_IP"

# Send announcement (retry logic)
MAX_RETRIES=10
RETRY_DELAY=30

for i in $(seq 1 $MAX_RETRIES); do
    echo "   Attempt $i/$MAX_RETRIES..."
    
    # Try to connect and send announcement
    RESPONSE=$(echo "$ANNOUNCEMENT" | nc -w 10 $STATION_IP $STATION_PORT 2>/dev/null || echo "")
    
    if [ -n "$RESPONSE" ]; then
        echo "$RESPONSE" > /opt/syntropy/metadata/registration-ack.json
        
        # Check if accepted
        STATUS=$(echo "$RESPONSE" | jq -r '.status' 2>/dev/null || echo "unknown")
        
        if [ "$STATUS" = "accepted" ]; then
            echo "✅ Registration accepted!"
            echo "$RESPONSE" | jq '.'
            exit 0
        else
            REASON=$(echo "$RESPONSE" | jq -r '.reason' 2>/dev/null || echo "unknown")
            echo "❌ Registration rejected: $REASON"
            exit 1
        fi
    fi
    
    echo "   No response, retrying in ${RETRY_DELAY}s..."
    sleep $RETRY_DELAY
done

echo "❌ Registration failed after $MAX_RETRIES attempts"
exit 1
REGISTER_SCRIPT
  
  - chmod +x /opt/syntropy/bin/register-node
  
  # Systemd service for Agent
  - |
    cat > /etc/systemd/system/syntropy-agent.service << 'SERVICE'
[Unit]
Description=Syntropy Cooperative Grid Agent
After=network-online.target docker.service
Wants=network-online.target docker.service

[Service]
Type=simple
User=syntropy
Group=syntropy
WorkingDirectory=/opt/syntropy
ExecStart=/opt/syntropy/bin/syntropy-agent
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
SERVICE
  
  # Systemd service for Registration (one-shot)
  - |
    cat > /etc/systemd/system/syntropy-register.service << 'REGISTER_SERVICE'
[Unit]
Description=Syntropy Node Registration
After=network-online.target syntropy-agent.service
Wants=network-online.target

[Service]
Type=oneshot
User=syntropy
Group=syntropy
ExecStart=/opt/syntropy/bin/register-node
RemainAfterExit=yes
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
REGISTER_SERVICE
  
  # Enable and start services
  - systemctl daemon-reload
  - systemctl enable syntropy-agent
  - systemctl enable syntropy-register
  - systemctl start syntropy-agent
  
  # Start registration (with delay to ensure network is up)
  - sleep 10
  - systemctl start syntropy-register

# Final message
final_message: |
  ✅ Syntropy Node ${NODE_NAME} provisioned!
  
  📊 Hardware: /opt/syntropy/metadata/hardware-manifest.json
  🔐 Grid Token: ${GRID_TOKEN:0:8}...[HIDDEN]
  🌐 Command Station: ${COMMAND_STATION_IP}:51000
  
  🚀 Registration in progress...
  📝 Status: journalctl -u syntropy-register -f
```

### 2.3 Template Simplificado: network-config

**Criar**: `infrastructure/cloud-init/network-config-mvp.yaml`

```yaml
# Syntropy Cooperative Grid - Simple Network Config (MVP)

version: 2
ethernets:
  # Auto-detect all Ethernet interfaces
  all:
    match:
      name: "en*"
    dhcp4: true
    dhcp6: false
    dhcp4-overrides:
      use-dns: true
      use-routes: true
    nameservers:
      addresses:
        - 1.1.1.1
        - 8.8.8.8
```

### 2.4 Template Simplificado: meta-data

**Criar**: `infrastructure/cloud-init/meta-data-mvp.yaml`

```yaml
# Syntropy Cooperative Grid - Node Metadata (MVP)

instance-id: ${NODE_NAME}-${INSTANCE_ID}
local-hostname: syntropy-${NODE_NAME}
```

### 2.5 Migração dos Templates

**Estratégia**:
1. Manter templates complexos atuais como `-advanced`
2. Criar versões `-mvp` simplificadas
3. MVP usa versões simplificadas
4. Pós-MVP pode migrar para advanced

```
infrastructure/cloud-init/
├── user-data-template.yaml           → user-data-advanced.yaml
├── user-data-mvp.yaml                → NOVO (usar no MVP)
├── network-config-template.yaml      → network-config-advanced.yaml
├── network-config-mvp.yaml           → NOVO (usar no MVP)
├── meta-data-template.yaml           → meta-data-advanced.yaml
└── meta-data-mvp.yaml                → NOVO (usar no MVP)
```

---

## 📄 3. MELHORIAS NO DOCUMENTO MVP

### 3.1 Adicionar Seção: Agent Placeholder

**Inserir após seção 4.4**:

```markdown
### 4.5 Agent Placeholder (MVP Fase 1)

**IMPORTANTE**: Para o MVP funcionar SEM o Agent completo, usamos um placeholder.

#### Por que Placeholder?
- Agent completo é complexo (Registration, Sync, Docker management)
- MVP prioriza provisionamento funcional
- Placeholder permite testar todo o fluxo

#### Implementação do Placeholder

**Arquivo**: USB deve conter `syntropy/agent` (script bash)

```bash
#!/bin/bash
# Syntropy Agent Placeholder (MVP Phase 1)

case "$1" in
  status)
    # Return fake status (allows Command Station to poll)
    MANIFEST=$(cat /opt/syntropy/metadata/hardware-manifest.json 2>/dev/null || echo '{}')
    cat << STATUS
{
  "status": "online",
  "node_id": "$(hostname)",
  "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "hardware": $MANIFEST
}
STATUS
    ;;
    
  exec)
    # Execute JSON command
    COMMAND=$(echo "$2" | jq -r '.command')
    case "$COMMAND" in
      deploy)
        # TODO: Implement actual deployment
        echo '{"status":"ok","message":"Placeholder - not implemented"}'
        ;;
      *)
        echo '{"status":"error","message":"Unknown command"}'
        ;;
    esac
    ;;
    
  *)
    echo "Syntropy Agent Placeholder v1.0"
    echo "Usage: $0 {status|exec}"
    ;;
esac
```

#### Roadmap do Agent

**Fase 1 (Semana 1-2)**: Placeholder
- Status reporting
- Responde a polling

**Fase 2 (Semana 3-4)**: Registration
- Node announcement
- Handshake completo

**Fase 3 (Semana 5-6)**: Workload Management
- Deploy via Docker
- Lifecycle management

**Fase 4 (Pós-MVP)**: Features Avançadas
- Metrics collection
- Log aggregation
- Auto-healing
```

### 3.2 Adicionar Seção: Pré-requisitos do Sistema

**Inserir após seção 2 (Componentes)**:

```markdown
### 2.X Pré-requisitos do Sistema

#### Command Station (PC de Trabalho)

**Sistema Operacional**:
- ✅ Windows 10/11 (64-bit)
- ✅ Ubuntu 22.04+ / Debian 11+
- ✅ macOS 12+

**Software Necessário**:
```bash
# Go 1.22+
go version  # deve mostrar 1.22 ou superior

# Git
git --version

# SSH client
ssh -V

# Dependências de sistema (Linux)
sudo apt-get install libsecret-1-dev  # Para keyring

# PowerShell (Windows) - já incluído
# diskpart (Windows) - já incluído

# dd (Linux) - já incluído
# lsblk (Linux) - já incluído
```

**Hardware Mínimo**:
- CPU: 2 cores
- RAM: 4GB
- Disk: 20GB livres (para cache de ISOs)
- USB: Porta USB 2.0+ disponível

#### Nodes (Hardware Físico)

**Hardware Mínimo por Node**:
- CPU: 4 cores (8 recomendado)
- RAM: 16GB (28GB recomendado para MVP)
- Disk: 256GB (512GB recomendado, SSD/NVMe)
- Network: Ethernet 1Gbps
- BIOS: Suporte a boot via USB

**Requisitos de Rede**:
- Todos os Nodes na mesma subnet
- DHCP disponível
- Roteador/Switch com 6+ portas
- Latência <10ms entre Nodes
- Acesso à Internet (para apt-get, Docker pulls)
```

### 3.3 Adicionar Seção: Troubleshooting Avançado

**Adicionar ao final do documento (antes de Referências)**:

```markdown
## 11. TROUBLESHOOTING AVANÇADO

### 11.1 Grid Token não salva no Keyring

**Sintoma**:
```bash
$ syntropy setup run
❌ failed to save token to keyring: secret service not available
```

**Causa**: Keyring do sistema não disponível (comum em servidores sem GUI)

**Solução**:
```bash
# Opção 1: Instalar gnome-keyring
sudo apt-get install gnome-keyring
gnome-keyring-daemon --start

# Opção 2: Usar fallback para arquivo criptografado
syntropy setup run --token-storage=file
# Sistema pedirá uma passphrase
```

### 11.2 USB não detectado no Linux

**Sintoma**:
```bash
$ syntropy node create
No USB devices found
```

**Causa**: Permissões insuficientes

**Solução**:
```bash
# Adicionar usuário ao grupo disk
sudo usermod -aG disk $USER

# Re-login necessário
# Ou: newgrp disk

# Verificar dispositivos
lsblk -o NAME,SIZE,TYPE,MOUNTPOINT
```

### 11.3 Registration falha (Node não aparece)

**Sintoma**:
```bash
# No Node
$ journalctl -u syntropy-register -f
❌ Registration failed after 10 attempts
```

**Diagnóstico**:
```bash
# 1. Verificar conectividade
ping <COMMAND_STATION_IP>

# 2. Verificar porta 51000 aberta
nc -zv <COMMAND_STATION_IP> 51000

# 3. Verificar Grid Token
grep grid_token /opt/syntropy/config/agent.yaml

# 4. Na Command Station, verificar listener
ss -tlnp | grep 51000
```

**Soluções**:
```bash
# Se porta 51000 não está aberta na Command Station:
syntropy node listen  # Inicia listener manualmente

# Se Grid Token não bate:
# Recriar USB com token correto

# Se firewall bloqueando:
# No Node:
sudo ufw allow from <COMMAND_STATION_IP> to any port 51000
```

### 11.4 Cloud-init não executa

**Sintoma**: Node boota, mas nenhuma configuração foi aplicada

**Diagnóstico**:
```bash
# No Node (via monitor/teclado físico ou SSH de recuperação)
sudo cloud-init status
sudo cloud-init analyze show

# Ver logs
sudo cat /var/log/cloud-init.log
sudo cat /var/log/cloud-init-output.log
```

**Causas Comuns**:
1. Cloud-init não presente no ISO
2. Sintaxe YAML inválida
3. Variáveis não substituídas

**Solução**:
```bash
# Validar YAML antes de criar USB
cloud-init schema --config-file user-data-mvp.yaml

# Verificar se variáveis foram substituídas
grep '\${' user-data-mvp.yaml  # Não deve retornar nada
```
```

### 3.4 Atualizar Roadmap com Fases Realistas

**Substituir Seção 7 (Roadmap)**:

```markdown
## 7. ROADMAP DE IMPLEMENTAÇÃO (REALISTA)

### Semana 1: Setup + Token Security
```
[🔧 TODO] Implementar TokenManager
  - Integrar go-keyring
  - Testar em Windows/Linux/macOS
  - Comandos: token show, export, import

[🔧 TODO] Completar Setup
  - Gerar Grid Token via Keyring
  - Atualizar testes
  
[📋 TESTE] Setup completo em 3 plataformas
```

### Semana 2: Cloud-Init Templates
```
[🔧 TODO] Criar templates MVP
  - user-data-mvp.yaml
  - network-config-mvp.yaml
  - meta-data-mvp.yaml
  
[🔧 TODO] Script detect-hardware
  - Implementar no template
  - Testar detecção em hardware real
  
[🔧 TODO] Script register-node
  - Implementar no template
  - Testar comunicação com Command Station

[📋 TESTE] Boot manual em VM com templates
```

### Semana 3: USB Creation (Windows)
```
[🔧 TODO] Implementar USB detection Windows
  - usb_detector_windows.go
  - PowerShell integration
  - Safety checks (block C:\)
  
[🔧 TODO] Implementar USB writing Windows
  - Inject cloud-init no ISO
  - Gravar via diskpart/PowerShell
  
[🔧 TODO] Agent Placeholder
  - Script bash simples
  - Status reporting
  - Copiar para USB

[📋 TESTE] Criar USB no Windows
```

### Semana 4: USB Creation (Linux) + First Node
```
[🔧 TODO] Implementar USB detection Linux
  - usb_detector_linux.go
  - lsblk integration
  - Safety checks (block /dev/sda)
  
[🔧 TODO] Implementar USB writing Linux
  - Inject cloud-init no ISO
  - Gravar via dd
  
[📋 TESTE] Provisionar PRIMEIRO Node físico real
  - Criar USB
  - Boot hardware virgem
  - Verificar auto-detecção
```

### Semana 5: Registration Protocol
```
[🔧 TODO] Implementar Listener (Command Station)
  - registration.go
  - Porta 51000
  - Validação de Grid Token
  
[🔧 TODO] Implementar Inventory Manager
  - inventory.go
  - CRUD de Nodes
  - Hardware Manifest storage
  
[📋 TESTE] Registration completo
  - Node anuncia
  - Command Station valida
  - Inventory atualizado
```

### Semana 6: Provisionar 6 Nodes + Básico de Workload
```
[🔧 TODO] Provisionar 6 Nodes completos
  - Criar 6 USBs
  - Provisionar hardware
  - Todos registrados

[🔧 TODO] Deploy básico via SSH
  - workload/deploy/deploy.go
  - Execução: ssh node "docker run..."
  - Salvar em workloads/

[📋 TESTE MVP COMPLETO]
  - 6 Nodes online
  - syntropy node list (mostra 6)
  - Deploy Nginx em 3 Nodes
  - Acessar via http://node-ip
```

### Pós-MVP (Semana 7+)
```
- Agent completo (substituir placeholder)
- Sincronização automática
- Workload scheduler avançado
- Management commands (start/stop/restart)
- Metrics collection
- Dashboard TUI
```
```

---

## 🎯 4. RESUMO DAS CORREÇÕES

### Críticas (Implementar Imediatamente)
1. ✅ **Grid Token**: Migrar para Keyring do sistema
2. ✅ **Templates**: Criar versões MVP simplificadas
3. ✅ **Agent**: Documentar uso de placeholder
4. ✅ **Hardware Detection**: Script completo no template
5. ✅ **Registration**: Script completo no template
6. ✅ **Pré-requisitos**: Documentar dependências

### Importantes (Implementar Durante MVP)
7. ✅ **Troubleshooting**: Seção expandida
8. ✅ **Roadmap**: Fases mais realistas
9. ✅ **Consistência**: Alinhar MVP com templates
10. ✅ **Variáveis**: Padronizar nomenclatura

### Melhorias Futuras (Pós-MVP)
11. Agent completo em Go
12. TLS além de SSH
13. Metrics collection automática
14. Dashboard web

---

**Próximo Passo**: Implementar TokenManager e atualizar templates para MVP.

