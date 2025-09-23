package usb

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// generateCloudInitConfig gera a configuração do cloud-init
func generateCloudInitConfig(config *Config, workDir string, certs *Certificates) (*CloudInitConfig, error) {
	// Gerar ID único da instância
	instanceID := fmt.Sprintf("%d", time.Now().UnixNano())

	// Detectar interface de rede primária
	networkInterface := "eth0" // Padrão, será detectado automaticamente

	// Configurar gateway padrão
	gateway := "192.168.1.1" // Padrão, será detectado automaticamente

	// Gerar sufixo IP único
	nodeIPSuffix := fmt.Sprintf("%d", time.Now().Unix()%254+2)

	// Caminhos dos certificados
	certDir := filepath.Join(workDir, "certs")
	nodeCertPath := filepath.Join(certDir, "node.crt")
	nodeKeyPath := filepath.Join(certDir, "node.key")
	caCertPath := filepath.Join(certDir, "ca.crt")

	return &CloudInitConfig{
		NodeName:         config.NodeName,
		NodeDescription:  config.NodeDescription,
		Coordinates:      config.Coordinates,
		OwnerKey:         config.OwnerKeyFile,
		DiscoveryServer:  config.DiscoveryServer,
		SSHPublicKey:     config.SSHPublicKey,
		CreatedBy:        config.CreatedBy,
		CreatedAt:        time.Now().Format(time.RFC3339),
		InstanceID:       instanceID,
		Interface:        networkInterface,
		Gateway:          gateway,
		NodeIPSuffix:     nodeIPSuffix,
		PrimaryInterface: networkInterface,
		MeshGateway:      "172.20.0.1",
		MgmtGateway:      "192.168.100.1",
		HTTPProxy:        "",
		HTTPSProxy:       "",
		NodeType:         "worker", // Padrão
		HardwareType:     "generic",
		CPUCores:         4,   // Será detectado automaticamente
		MemoryGB:         8,   // Será detectado automaticamente
		StorageGB:        100, // Será detectado automaticamente
		InitialRole:      "worker",
		CanBeLeader:      true,
		CanBeWorker:      true,
		NodeCertPath:     nodeCertPath,
		NodeKeyPath:      nodeKeyPath,
		CACertPath:       caCertPath,
		// Certificados PEM para write_files
		CACertPEM:   string(certs.CACert),
		NodeCertPEM: string(certs.NodeCert),
		NodeKeyPEM:  string(certs.NodeKey),
	}, nil
}

// renderTemplate renderiza um template com os dados fornecidos
func renderTemplate(templateStr string, data interface{}) (string, error) {
	tmpl, err := template.New("config").Option("missingkey=error").Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf("erro ao fazer parse do template: %w", err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("erro ao executar template: %w", err)
	}

	return buf.String(), nil
}

// createCloudInitFiles cria os arquivos de configuração do cloud-init
func createCloudInitFiles(config *CloudInitConfig, workDir string, certPaths map[string]string) error {
	cloudInitDir := filepath.Join(workDir, "cloud-init")
	if err := os.MkdirAll(cloudInitDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório cloud-init: %w", err)
	}

	// Template do user-data
	userDataTemplate := `#cloud-config
# Ubuntu Server 24.04 LTS - Syntropy Cooperative Grid Node
# Este arquivo configura automaticamente um nó da rede Syntropy
# Gerado automaticamente pelo PC de trabalho (Quartel General)

# Configurações básicas do sistema
locale: pt_BR.UTF-8
timezone: America/Sao_Paulo
hostname: {{.NodeName}}

# Usuário padrão (será substituído pela configuração Syntropy)
users:
  - name: syntropy
    groups: [adm, sudo, docker]
    shell: /bin/bash
    sudo: ALL=(ALL) NOPASSWD:ALL
    lock_passwd: true

# Configuração de rede
network:
  version: 2
  ethernets:
    {{.Interface}}:
      dhcp4: true
      dhcp6: false

# Configuração SSH
ssh_pwauth: false
disable_root: true
ssh_authorized_keys:
  - {{.SSHPublicKey}}

# Escrever arquivos de certificados
write_files:
  - path: /opt/syntropy/certs/ca.crt
    owner: root:root
    permissions: "0644"
    content: |
      {{.CACertPEM}}
  - path: /opt/syntropy/certs/node.crt
    owner: root:root
    permissions: "0644"
    content: |
      {{.NodeCertPEM}}
  - path: /opt/syntropy/certs/node.key
    owner: root:root
    permissions: "0600"
    content: |
      {{.NodeKeyPEM}}

# Pacotes a serem instalados
packages:
  - curl
  - wget
  - git
  - htop
  - vim
  - net-tools
  - dnsutils
  - fail2ban
  - ufw
  - docker.io
  - docker-compose-plugin
  - containerd
  - kubectl
  - wireguard
  - jq
  - openssl
  - ca-certificates
  - gnupg
  - lsb-release
  - prometheus-node-exporter
  - ntp
  - rsync
  - unzip
  - tree
  - tmux

# Configuração do Docker
runcmd:
  # Configurar Docker
  - systemctl enable docker
  - systemctl start docker
  - usermod -aG docker syntropy
  
  # Configurar firewall
  - ufw --force enable
  - ufw default deny incoming
  - ufw default allow outgoing
  - ufw allow ssh
  - ufw allow 6443/tcp
  - ufw allow 2379:2380/tcp
  - ufw allow 10250/tcp
  - ufw allow 10251/tcp
  - ufw allow 10252/tcp
  - ufw allow 30000:32767/tcp
  - ufw allow 51820/udp
  - ufw allow 8080/tcp
  - ufw allow 9090/tcp
  - ufw allow 9100/tcp
  
  # Configurar fail2ban
  - systemctl enable fail2ban
  - systemctl start fail2ban
  
  # Criar diretórios Syntropy
  - mkdir -p /opt/syntropy/{bin,config,logs,certs,data,scripts,backups,audit}
  - chown -R syntropy:syntropy /opt/syntropy
  
  # Download e instalação do Syntropy Agent
  - curl -L https://github.com/syntropy-cooperative-grid/agent/releases/latest/download/syntropy-agent-linux-amd64 -o /opt/syntropy/bin/syntropy-agent
  - chmod +x /opt/syntropy/bin/syntropy-agent
  
  # Configurar certificados (já criados via write_files)
  - chmod 600 /opt/syntropy/certs/*
  - chown syntropy:syntropy /opt/syntropy/certs/*
  
  # Configurar Syntropy Agent
  - cat > /opt/syntropy/config/agent.yaml << EOF
node:
  name: "{{.NodeName}}"
  type: "{{.NodeType}}"
  description: "{{.NodeDescription}}"
  coordinates: "{{.Coordinates}}"
  owner_key: "{{.OwnerKey}}"

network:
  discovery_endpoints:
    - "https://{{.DiscoveryServer}}:8443"
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
EOF
  
  # Configurar systemd service
  - cat > /etc/systemd/system/syntropy-agent.service << EOF
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

[Install]
WantedBy=multi-user.target
EOF
  
  # Iniciar Syntropy Agent
  - systemctl daemon-reload
  - systemctl enable syntropy-agent
  - systemctl start syntropy-agent
  
  # Configurar logrotate
  - cat > /etc/logrotate.d/syntropy << EOF
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

# Configurações finais
final_message: |
  ✅ Nó Syntropy configurado com sucesso!
  
  📊 Informações do Nó:
  - Nome: {{.NodeName}}
  - Tipo: {{.NodeType}}
  - Descrição: {{.NodeDescription}}
  - Coordenadas: {{.Coordinates}}
  
  🔐 Segurança:
  - SSH configurado com chave pública
  - Firewall ativo
  - Fail2ban configurado
  - Certificados TLS instalados
  
  🌐 Rede:
  - Descoberta: {{.DiscoveryServer}}
  - Mesh: porta 51820
  - API: porta 8080
  
  📝 Logs:
  - Agent: /opt/syntropy/logs/agent.log
  - System: journalctl -u syntropy-agent
  
  🚀 Status: systemctl status syntropy-agent`

	// Renderizar e criar arquivo user-data
	userDataContent, err := renderTemplate(userDataTemplate, config)
	if err != nil {
		return fmt.Errorf("erro ao renderizar user-data: %w", err)
	}

	userDataPath := filepath.Join(cloudInitDir, "user-data")
	if err := os.WriteFile(userDataPath, []byte(userDataContent), 0644); err != nil {
		return fmt.Errorf("erro ao criar user-data: %w", err)
	}

	// Template do meta-data
	metaDataTemplate := `# Syntropy Cooperative Grid - Node Metadata
# Gerado automaticamente pelo PC de trabalho (Quartel General)

instance-id: {{.NodeName}}-{{.InstanceID}}
local-hostname: {{.NodeName}}

# Informações do nó Syntropy
syntropy:
  node:
    name: {{.NodeName}}
    type: {{.NodeType}}
    description: {{.NodeDescription}}
    coordinates: {{.Coordinates}}
    owner_key: {{.OwnerKey}}
    created_at: {{.CreatedAt}}
    created_by: {{.CreatedBy}}
  
  network:
    discovery_server: {{.DiscoveryServer}}
    mesh_port: 51820
    api_port: 8080
    metrics_port: 9090
  
  security:
    ca_cert: {{.CACertPath}}
    node_cert: {{.NodeCertPath}}
    node_key: {{.NodeKeyPath}}
    ssh_public_key: {{.SSHPublicKey}}
  
  hardware:
    detected_type: {{.HardwareType}}
    cpu_cores: {{.CPUCores}}
    memory_gb: {{.MemoryGB}}
    storage_gb: {{.StorageGB}}
  
  role:
    initial_role: {{.InitialRole}}
    can_be_leader: {{.CanBeLeader}}
    can_be_worker: {{.CanBeWorker}}
  
  audit:
    enabled: true
    log_level: info
    retention_days: 90`

	// Renderizar e criar arquivo meta-data
	metaDataContent, err := renderTemplate(metaDataTemplate, config)
	if err != nil {
		return fmt.Errorf("erro ao renderizar meta-data: %w", err)
	}

	metaDataPath := filepath.Join(cloudInitDir, "meta-data")
	if err := os.WriteFile(metaDataPath, []byte(metaDataContent), 0644); err != nil {
		return fmt.Errorf("erro ao criar meta-data: %w", err)
	}

	// Template do network-config
	networkConfigTemplate := `# Syntropy Cooperative Grid - Network Configuration
# Gerado automaticamente pelo PC de trabalho (Quartel General)

version: 2
ethernets:
  # Configuração automática para interfaces Ethernet
  en*:
    dhcp4: true
    dhcp6: false
    dhcp4-overrides:
      hostname: {{.NodeName}}
      use-dns: true
      use-routes: true
      use-domains: true
    nameservers:
      addresses:
        - 8.8.8.8
        - 8.8.4.4
        - 1.1.1.1
        - 1.0.0.1
    routes:
      - to: 0.0.0.0/0
        via: {{.Gateway}}
        metric: 100
  
  # Configuração automática para interfaces Ethernet (nomenclatura alternativa)
  eth*:
    dhcp4: true
    dhcp6: false
    dhcp4-overrides:
      hostname: {{.NodeName}}
      use-dns: true
      use-routes: true
      use-domains: true
    nameservers:
      addresses:
        - 8.8.8.8
        - 8.8.4.4
        - 1.1.1.1
        - 1.0.0.1
    routes:
      - to: 0.0.0.0/0
        via: {{.Gateway}}
        metric: 100
  
  # Configuração automática para interfaces Ethernet (nomenclatura moderna)
  enp*:
    dhcp4: true
    dhcp6: false
    dhcp4-overrides:
      hostname: {{.NodeName}}
      use-dns: true
      use-routes: true
      use-domains: true
    nameservers:
      addresses:
        - 8.8.8.8
        - 8.8.4.4
        - 1.1.1.1
        - 1.0.0.1
    routes:
      - to: 0.0.0.0/0
        via: {{.Gateway}}
        metric: 100

# Configuração de bridge para virtualização (se necessário)
bridges:
  br0:
    interfaces: []
    dhcp4: false
    dhcp6: false
    addresses:
      - 172.20.0.{{.NodeIPSuffix}}/24
    gateway4: 172.20.0.1
    nameservers:
      addresses:
        - 8.8.8.8
        - 8.8.4.4
        - 1.1.1.1
        - 1.0.0.1
    parameters:
      stp: false
      forward-delay: 0

# Configuração de VLAN (se necessário)
vlans:
  vlan100:
    id: 100
    link: {{.PrimaryInterface}}
    dhcp4: false
    dhcp6: false
    addresses:
      - 192.168.100.{{.NodeIPSuffix}}/24
    gateway4: 192.168.100.1
    nameservers:
      addresses:
        - 8.8.8.8
        - 8.8.4.4
        - 1.1.1.1
        - 1.0.0.1

# Configuração de roteamento estático para redes Syntropy
routes:
  # Rota para rede mesh Syntropy
  - to: 172.20.0.0/12
    via: {{.MeshGateway}}
    metric: 50
    table: 100
  # Rota para rede de gerenciamento
  - to: 192.168.100.0/24
    via: {{.MgmtGateway}}
    metric: 75
    table: 100

# Configuração de regras de roteamento
routing-policy:
  - from: 172.20.0.0/12
    table: 100
  - from: 192.168.100.0/24
    table: 100

# Configuração de proxy (se necessário)
proxy:
  http: {{.HTTPProxy}}
  https: {{.HTTPSProxy}}
  no_proxy:
    - localhost
    - 127.0.0.1
    - 172.20.0.0/12
    - 192.168.100.0/24`

	// Renderizar e criar arquivo network-config
	networkConfigContent, err := renderTemplate(networkConfigTemplate, config)
	if err != nil {
		return fmt.Errorf("erro ao renderizar network-config: %w", err)
	}

	networkConfigPath := filepath.Join(cloudInitDir, "network-config")
	if err := os.WriteFile(networkConfigPath, []byte(networkConfigContent), 0644); err != nil {
		return fmt.Errorf("erro ao criar network-config: %w", err)
	}

	return nil
}

// copyScripts copia os scripts de instalação para o diretório de trabalho
func copyScripts(workDir string) error {
	scriptsDir := filepath.Join(workDir, "scripts")
	if err := os.MkdirAll(scriptsDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório de scripts: %w", err)
	}

	// Determinar o caminho correto para os scripts no projeto
	// Primeiro tenta o caminho relativo atual (para compatibilidade)
	projectScriptsDir := "infrastructure/cloud-init/scripts"

	// Se não existir, tenta encontrar o diretório raiz do projeto
	if _, err := os.Stat(projectScriptsDir); os.IsNotExist(err) {
		// Procura pelo diretório raiz do projeto (contém go.mod)
		currentDir, _ := os.Getwd()
		searchDir := currentDir

		// Sobe na hierarquia até encontrar o diretório raiz do projeto
		for i := 0; i < 10; i++ { // Máximo 10 níveis para evitar loop infinito
			// Verifica se este diretório contém o diretório infrastructure/cloud-init/scripts
			testPath := filepath.Join(searchDir, "infrastructure", "cloud-init", "scripts")
			if _, err := os.Stat(testPath); err == nil {
				// Encontrou o diretório correto
				projectScriptsDir = testPath
				break
			}
			parentDir := filepath.Dir(searchDir)
			if parentDir == searchDir {
				// Chegou na raiz do filesystem
				break
			}
			searchDir = parentDir
		}

		// Verifica se o diretório encontrado realmente existe
		if _, err := os.Stat(projectScriptsDir); os.IsNotExist(err) {
			return fmt.Errorf("não foi possível encontrar o diretório de scripts: %s", projectScriptsDir)
		}
	}

	// Lista de scripts para copiar
	scripts := []string{
		"hardware-detection.sh",
		"network-discovery.sh",
		"syntropy-install.sh",
		"cluster-join.sh",
	}

	for _, script := range scripts {
		srcPath := filepath.Join(projectScriptsDir, script)
		dstPath := filepath.Join(scriptsDir, script)

		// Verificar se o arquivo fonte existe
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			return fmt.Errorf("script %s não encontrado em %s", script, srcPath)
		}

		// Copiar arquivo
		srcData, err := os.ReadFile(srcPath)
		if err != nil {
			return fmt.Errorf("erro ao ler script %s: %w", script, err)
		}

		if err := os.WriteFile(dstPath, srcData, 0755); err != nil {
			return fmt.Errorf("erro ao copiar script %s: %w", script, err)
		}

		fmt.Printf("✅ Script %s copiado com sucesso\n", script)
	}

	return nil
}
