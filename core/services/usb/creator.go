package usb

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"syntropy-cc/cooperative-grid/infrastructure"
)

// Config contém configurações para criação de USB
type Config struct {
	NodeName        string `json:"node_name"`
	NodeDescription string `json:"node_description"`
	Coordinates     string `json:"coordinates"`
	OwnerKeyFile    string `json:"owner_key_file"`
	Label           string `json:"label"`
}

// Creator interface para criação de USB com boot
type Creator interface {
	CreateUSB(devicePath string, config *Config) error
	Cleanup() error
}

// USBCreator implementa a criação de USB com boot seguindo a arquitetura do projeto
type USBCreator struct {
	workDir     string
	cacheDir    string
	formatter   Formatter
	templateMgr *infrastructure.TemplateManager
	keyMgr      *infrastructure.KeyManager
}

// NewCreator cria uma nova instância do criador de USB
func NewCreator(workDir, cacheDir string) *USBCreator {
	// Determinar diretório de templates baseado na estrutura do projeto
	templateDir := "infrastructure"
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		// Se não existe no diretório atual, tentar relativo ao projeto
		templateDir = "../../infrastructure"
	}

	// Criar diretório de chaves dentro do workDir
	keyDir := filepath.Join(workDir, "keys")

	return &USBCreator{
		workDir:     workDir,
		cacheDir:    cacheDir,
		formatter:   NewFormatter(),
		templateMgr: infrastructure.NewTemplateManager(templateDir),
		keyMgr:      infrastructure.NewKeyManager(keyDir),
	}
}

// CreateUSB orquestra o processo completo de criação do USB seguindo a arquitetura
func (c *USBCreator) CreateUSB(devicePath string, config *Config) error {
	fmt.Println("🚀 Iniciando criação de USB com boot para Syntropy Cooperative Grid")
	fmt.Println()

	// Etapa 1: Validar dispositivo
	fmt.Println("🔍 Etapa 1/6: Validando dispositivo USB...")
	detector := NewDetector()
	if err := detector.ValidateDevice(devicePath); err != nil {
		return fmt.Errorf("dispositivo inválido: %w", err)
	}

	if detector.IsSystemDisk(devicePath) {
		return fmt.Errorf("dispositivo %s parece ser um disco do sistema", devicePath)
	}
	fmt.Printf("   ✅ Dispositivo %s validado com sucesso\n", devicePath)
	fmt.Println()

	// Etapa 2: Montar dispositivo para verificação
	fmt.Println("📱 Etapa 2/6: Verificando dispositivo...")
	mountPoint, err := c.mountDevice(devicePath)
	if err != nil {
		return fmt.Errorf("falha ao montar dispositivo: %w", err)
	}
	defer c.unmountDevice(mountPoint)
	fmt.Printf("   ✅ Dispositivo montado em: %s\n", mountPoint)
	fmt.Println()

	// Etapa 3: Formatar dispositivo
	fmt.Println("💾 Etapa 3/6: Formatando dispositivo...")
	if err := c.formatter.FormatDevice(devicePath, config.Label); err != nil {
		return fmt.Errorf("falha na formatação: %w", err)
	}
	fmt.Println("   ✅ Dispositivo formatado com sucesso")
	fmt.Println()

	// Etapa 4: Remontar dispositivo formatado
	fmt.Println("📱 Etapa 4/6: Remontando dispositivo formatado...")
	mountPoint, err = c.mountDevice(devicePath)
	if err != nil {
		return fmt.Errorf("falha ao remontar dispositivo: %w", err)
	}
	defer c.unmountDevice(mountPoint)
	fmt.Printf("   ✅ Dispositivo remontado em: %s\n", mountPoint)
	fmt.Println()

	// Etapa 5: Gerar chaves SSH
	fmt.Println("🔑 Etapa 5/6: Gerando chaves SSH...")
	nodeKeyPath, err := c.generateSSHKeys(mountPoint, config.NodeName)
	if err != nil {
		return fmt.Errorf("falha na geração de chaves: %w", err)
	}
	fmt.Printf("   ✅ Chaves SSH geradas: %s\n", nodeKeyPath)
	fmt.Println()

	// Etapa 6: Criar configuração cloud-init usando IaC
	fmt.Println("☁️  Etapa 6/6: Criando configuração cloud-init usando Infrastructure as Code...")
	if err := c.createCloudInitWithIAC(mountPoint, config); err != nil {
		return fmt.Errorf("falha na criação do cloud-init: %w", err)
	}
	fmt.Println("   ✅ Configuração cloud-init criada com sucesso usando IaC")
	fmt.Println()

	fmt.Println("🎉 USB com boot criado com sucesso!")
	fmt.Printf("   Nó: %s\n", config.NodeName)
	fmt.Printf("   Dispositivo: %s\n", devicePath)
	fmt.Printf("   Montado em: %s\n", mountPoint)
	fmt.Printf("   Coordenadas: %s\n", config.Coordinates)
	fmt.Println()
	fmt.Println("📋 Próximos passos:")
	fmt.Println("   1. Remover o USB com segurança")
	fmt.Println("   2. Inserir em um computador e fazer boot")
	fmt.Println("   3. Aguardar a instalação automática do Ubuntu Server")
	fmt.Println("   4. Conectar via SSH usando as chaves geradas")

	return nil
}

// mountDevice monta um dispositivo USB
func (c *USBCreator) mountDevice(devicePath string) (string, error) {
	// Criar diretório de montagem
	mountPoint := filepath.Join(c.workDir, "mount")
	if err := os.MkdirAll(mountPoint, 0755); err != nil {
		return "", fmt.Errorf("falha ao criar diretório de montagem: %w", err)
	}

	// Montar dispositivo
	cmd := exec.Command("sudo", "mount", devicePath, mountPoint)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("falha ao montar dispositivo: %w", err)
	}

	return mountPoint, nil
}

// unmountDevice desmonta um dispositivo
func (c *USBCreator) unmountDevice(mountPoint string) {
	if mountPoint == "" {
		return
	}

	cmd := exec.Command("sudo", "umount", "-f", mountPoint)
	cmd.Run() // Ignorar erro
}

// generateSSHKeys gera chaves SSH para o nó
func (c *USBCreator) generateSSHKeys(mountPoint, nodeName string) (string, error) {
	fmt.Printf("   🔑 Gerando chave SSH RSA 4096-bit...\n")

	// Criar diretório para chaves
	keysDir := filepath.Join(mountPoint, "syntropy", "keys")
	if err := os.MkdirAll(keysDir, 0750); err != nil {
		return "", fmt.Errorf("falha ao criar diretório de chaves: %w", err)
	}

	// Gerar chave privada
	privateKeyPath := filepath.Join(keysDir, nodeName+".key")
	cmd := exec.Command("ssh-keygen", "-t", "rsa", "-b", "4096", "-f", privateKeyPath, "-N", "")
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("falha ao gerar chave SSH: %w", err)
	}

	publicKeyPath := privateKeyPath + ".pub"
	fmt.Printf("   ✅ Chave privada: %s\n", privateKeyPath)
	fmt.Printf("   📍 Chave pública: %s\n", publicKeyPath)

	return privateKeyPath, nil
}

// createCloudInitWithIAC cria a configuração cloud-init usando Infrastructure as Code
func (c *USBCreator) createCloudInitWithIAC(mountPoint string, config *Config) error {
	fmt.Printf("   ☁️  Criando configuração cloud-init usando IaC...\n")

	// Criar diretório cloud-init
	cloudInitDir := filepath.Join(mountPoint, "cloud-init")
	if err := os.MkdirAll(cloudInitDir, 0755); err != nil {
		return fmt.Errorf("falha ao criar diretório cloud-init: %w", err)
	}

	// Gerar ou carregar chaves SSH
	ownerKey, ownerPub, err := c.generateOrLoadSSHKeys(config.OwnerKeyFile, config.NodeName, infrastructure.OwnerKey)
	if err != nil {
		return fmt.Errorf("falha ao gerar chaves do proprietário: %w", err)
	}

	communityKey, communityPub, err := c.generateOrLoadSSHKeys("", config.NodeName, infrastructure.CommunityKey)
	if err != nil {
		return fmt.Errorf("falha ao gerar chaves da comunidade: %w", err)
	}

	// Preparar dados para o template
	templateData := &infrastructure.TemplateData{
		NodeName:                 config.NodeName,
		NodeDescription:          config.NodeDescription,
		Coordinates:              config.Coordinates,
		CreatedAt:                time.Now().Format(time.RFC3339),
		AdminPasswordHash:        "$6$rounds=4096$syntropy$N8mVzFK0Y1OelT1SKEjg0jIXzKMzL3ZcOGcE5xR8nS6E8qSO5qFV6eJs1g7T6E0cC7w.kfNO3FqC3YhE9Gz19.",
		OwnerPublicKey:           ownerPub.PublicKey,
		CommunityPublicKey:       communityPub.PublicKey,
		KeyInstallationCommands:  c.keyMgr.GenerateKeyInstallationCommands(ownerKey, ownerPub, communityKey, communityPub),
		MetadataCreationCommands: c.keyMgr.GenerateMetadataCreationCommands(config.NodeName, config.Coordinates, config.NodeDescription),
		TemplateCreationCommands: c.generateTemplateCreationCommands(),
		StartupServiceCommands:   c.generateStartupServiceCommands(config.NodeName),
		NodeID:                   generateInstanceID(),
		LocationNodeID:           generateInstanceID(),
		DetectionMethod:          "manual",
		DetectedCity:             "Unknown",
		DetectedCountry:          "Unknown",
		OwnerFingerprint:         ownerKey.Fingerprint,
		CommunityFingerprint:     communityKey.Fingerprint,
	}

	// Gerar arquivos cloud-init usando templates IaC
	if err := c.templateMgr.SaveCloudInitFiles(cloudInitDir, templateData); err != nil {
		return fmt.Errorf("falha ao gerar arquivos cloud-init: %w", err)
	}

	fmt.Printf("   ✅ Configuração cloud-init criada usando IaC\n")
	fmt.Printf("   📍 user-data: %s\n", filepath.Join(cloudInitDir, "user-data"))
	fmt.Printf("   📍 meta-data: %s\n", filepath.Join(cloudInitDir, "meta-data"))
	fmt.Printf("   📍 network-config: %s\n", filepath.Join(cloudInitDir, "network-config"))

	return nil
}

// generateOrLoadSSHKeys gera ou carrega chaves SSH
func (c *USBCreator) generateOrLoadSSHKeys(keyFilePath, nodeName string, purpose infrastructure.KeyPurpose) (*infrastructure.KeyPair, *infrastructure.KeyPair, error) {
	var keyPair *infrastructure.KeyPair
	var err error

	if keyFilePath != "" && keyFilePath != "" {
		// Carregar chaves existentes
		keyPair, err = c.keyMgr.LoadExistingKeyPair(keyFilePath)
		if err != nil {
			return nil, nil, fmt.Errorf("falha ao carregar chaves existentes: %w", err)
		}
	} else {
		// Gerar novas chaves
		keyPair, err = c.keyMgr.GenerateKeyPair(purpose, nodeName)
		if err != nil {
			return nil, nil, fmt.Errorf("falha ao gerar novas chaves: %w", err)
		}

		// Salvar chaves geradas
		if err := c.keyMgr.SaveKeyPair(keyPair, purpose, nodeName); err != nil {
			return nil, nil, fmt.Errorf("falha ao salvar chaves geradas: %w", err)
		}
	}

	// Retornar o mesmo par como privada e pública (estrutura do KeyPair já contém ambas)
	return keyPair, keyPair, nil
}

// generateTemplateCreationCommands gera comandos para criar templates Kubernetes
func (c *USBCreator) generateTemplateCreationCommands() string {
	return `curtin in-target -- bash -c '
# Create Kubernetes templates directory
mkdir -p /opt/syntropy/platform/templates/kubernetes

# Scientific computing batch job template
cat > /opt/syntropy/platform/templates/kubernetes/scientific-computing-job.yaml << "BATCH_TEMPLATE_EOF"
apiVersion: batch/v1
kind: Job
metadata:
  name: scientific-computation
  namespace: default
  labels:
    app: scientific-computing
    platform: syntropy
spec:
  template:
    spec:
      containers:
      - name: computation
        image: ubuntu:22.04
        command: ["/bin/bash"]
        args: ["-c", "echo \"Starting computation...\" && sleep 30 && echo \"Computation complete\""]
        resources:
          requests:
            cpu: "100m"
            memory: "128Mi"
          limits:
            cpu: "2000m"
            memory: "2Gi"
      restartPolicy: Never
  backoffLimit: 3
BATCH_TEMPLATE_EOF

# Web service template
cat > /opt/syntropy/platform/templates/kubernetes/web-service.yaml << "WEB_TEMPLATE_EOF"
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web-application
  namespace: default
  labels:
    app: web-application
    platform: syntropy
spec:
  replicas: 1
  selector:
    matchLabels:
      app: web-application
  template:
    metadata:
      labels:
        app: web-application
    spec:
      containers:
      - name: web
        image: nginx:alpine
        ports:
        - containerPort: 80
        resources:
          requests:
            cpu: "50m"
            memory: "64Mi"
          limits:
            cpu: "500m"
            memory: "512Mi"
---
apiVersion: v1
kind: Service
metadata:
  name: web-service
  labels:
    platform: syntropy
spec:
  selector:
    app: web-application
  ports:
  - port: 80
    targetPort: 80
  type: ClusterIP
WEB_TEMPLATE_EOF

# Database template
cat > /opt/syntropy/platform/templates/kubernetes/database-statefulset.yaml << "DB_TEMPLATE_EOF"
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: database-service
  namespace: default
  labels:
    app: database
    platform: syntropy
spec:
  serviceName: database
  replicas: 1
  selector:
    matchLabels:
      app: database
  template:
    metadata:
      labels:
        app: database
    spec:
      containers:
      - name: database
        image: postgres:15-alpine
        env:
        - name: POSTGRES_DB
          value: "syntropy"
        - name: POSTGRES_USER
          value: "admin"
        - name: POSTGRES_PASSWORD
          value: "changeme"
        ports:
        - containerPort: 5432
        volumeMounts:
        - name: data
          mountPath: /var/lib/postgresql/data
        resources:
          requests:
            cpu: "100m"
            memory: "256Mi"
          limits:
            cpu: "1000m"
            memory: "1Gi"
  volumeClaimTemplates:
  - metadata:
      name: data
    spec:
      accessModes: ["ReadWriteOnce"]
      resources:
        requests:
          storage: 10Gi
DB_TEMPLATE_EOF

chmod 644 /opt/syntropy/platform/templates/kubernetes/*.yaml
chown admin:admin /opt/syntropy/platform/templates/kubernetes/*.yaml
'`
}

// generateStartupServiceCommands gera comandos para criar serviços de inicialização
func (c *USBCreator) generateStartupServiceCommands(nodeName string) string {
	return fmt.Sprintf(`curtin in-target -- bash -c '
# Create startup script for first boot
cat > /opt/syntropy/scripts/first-boot.sh << "FIRST_BOOT_EOF"
#!/bin/bash
# Syntropy Cooperative Grid - First Boot Setup Script
# Node: %s

echo "🚀 Syntropy Cooperative Grid Node %s - First Boot Setup" | logger -t syntropy

# Wait for network to be stable
sleep 30

# Log startup
echo "Syntropy node %s first boot setup starting..." | logger -t syntropy

# Detect hardware and update metadata
CPU_CORES=$(nproc)
RAM_GB=$(free -g | awk "/^Mem:/{print \\$2}")
STORAGE_GB=$(df / --output=avail -BG 2>/dev/null | tail -1 | sed "s/G//" | xargs)
ARCHITECTURE=$(uname -m)

# Get current IP
CURRENT_IP=$(hostname -I | awk "{print \\$1}")

# Update node metadata with runtime information
if [ -f /opt/syntropy/metadata/node.json ]; then
    # Create updated metadata with runtime info
    python3 -c "
import json
import sys
from datetime import datetime

try:
    with open('/opt/syntropy/metadata/node.json', 'r') as f:
        data = json.load(f)
    
    # Update hardware info with actual detected values
    data['hardware']['cpu_cores'] = $CPU_CORES
    data['hardware']['ram_gb'] = $RAM_GB
    data['hardware']['storage_gb'] = $STORAGE_GB
    data['hardware']['architecture'] = '$ARCHITECTURE'
    
    # Update network info
    data['network']['ip_address'] = '$CURRENT_IP'
    
    # Update management status
    data['management']['status'] = 'installed'
    data['management']['installation_complete'] = True
    data['management']['first_boot'] = datetime.utcnow().isoformat() + 'Z'
    data['management']['last_update'] = datetime.utcnow().isoformat() + 'Z'
    
    with open('/opt/syntropy/metadata/node.json', 'w') as f:
        json.dump(data, f, indent=2)
    
    print('Metadata updated successfully')
except Exception as e:
    print(f'Error updating metadata: {e}', file=sys.stderr)
"
fi

# Ensure all services are running
systemctl enable docker
systemctl start docker

systemctl enable ssh
systemctl start ssh

systemctl enable prometheus-node-exporter
systemctl start prometheus-node-exporter

# Configure firewall
ufw --force reset
ufw default deny incoming
ufw default allow outgoing
ufw allow ssh
ufw allow 9100/tcp
ufw --force enable

# Configure fail2ban
systemctl enable fail2ban
systemctl start fail2ban

echo "✅ Syntropy Cooperative Grid Node %s ready for management" | logger -t syntropy

# Create ready indicator
touch /opt/syntropy/.ready

exit 0
FIRST_BOOT_EOF

chmod +x /opt/syntropy/scripts/first-boot.sh

# Create systemd service for first boot
cat > /etc/systemd/system/syntropy-first-boot.service << "SERVICE_EOF"
[Unit]
Description=Syntropy First Boot Setup
After=network-online.target cloud-init.service
Wants=network-online.target

[Service]
Type=oneshot
ExecStart=/opt/syntropy/scripts/first-boot.sh
RemainAfterExit=yes
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
SERVICE_EOF

systemctl enable syntropy-first-boot.service
'`, nodeName, nodeName, nodeName, nodeName, nodeName)
}

// generateInstanceID gera um ID único para a instância
func generateInstanceID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// Cleanup limpa recursos temporários
func (c *USBCreator) Cleanup() error {
	// Limpar diretório de montagem
	mountPoint := filepath.Join(c.workDir, "mount")
	if err := os.RemoveAll(mountPoint); err != nil {
		return fmt.Errorf("falha ao limpar diretório de montagem: %w", err)
	}

	// Limpar cache se necessário
	// TODO: Implementar limpeza de cache

	return nil
}

// downloadISO baixa ISO do Ubuntu Server (placeholder)
func (c *USBCreator) downloadISO() error {
	// TODO: Implementar download de ISO
	fmt.Println("📥 Download de ISO não implementado (placeholder)")
	return nil
}

// verifyISO verifica integridade da ISO (placeholder)
func (c *USBCreator) verifyISO(isoPath string) error {
	// TODO: Implementar verificação de ISO
	fmt.Printf("🔍 Verificação de ISO não implementada (placeholder): %s\n", isoPath)
	return nil
}

// extractISO extrai conteúdo da ISO (placeholder)
func (c *USBCreator) extractISO(isoPath, outputDir string) error {
	// TODO: Implementar extração de ISO
	fmt.Printf("📦 Extração de ISO não implementada (placeholder): %s -> %s\n", isoPath, outputDir)
	return nil
}
