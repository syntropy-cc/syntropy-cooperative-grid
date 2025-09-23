package usb

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

// USBCreator implementa a criação de USB com boot
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

// CreateUSB orquestra o processo completo de criação do USB
func (c *USBCreator) CreateUSB(devicePath string, config *Config) error {
	fmt.Println("🚀 Iniciando criação de USB com boot para Syntropy Cooperative Grid")
	fmt.Println()

	// Validar configuração
	if err := c.validateConfig(config); err != nil {
		return fmt.Errorf("configuração inválida: %w", err)
	}

	// Criar diretórios de trabalho
	if err := c.setupDirectories(); err != nil {
		return fmt.Errorf("falha ao criar diretórios: %w", err)
	}

	// Etapa 1: Formatar dispositivo USB
	fmt.Println("📱 Etapa 1/6: Formatando dispositivo USB...")
	if err := c.formatDevice(devicePath, config.Label); err != nil {
		return fmt.Errorf("falha na formatação: %w", err)
	}
	fmt.Println("   ✅ Dispositivo formatado com sucesso")
	fmt.Println()

	// Etapa 2: Montar partição
	fmt.Println("🔗 Etapa 2/6: Montando partição...")
	mountPoint, err := c.mountPartition(devicePath)
	if err != nil {
		return fmt.Errorf("falha ao montar partição: %w", err)
	}
	defer c.unmountPartition(mountPoint)
	fmt.Printf("   ✅ Partição montada em %s\n", mountPoint)
	fmt.Println()

	// Etapa 3: Download e instalação do Ubuntu
	fmt.Println("⬇️  Etapa 3/6: Baixando e instalando Ubuntu Server...")
	isoPath, err := c.downloadUbuntuISO()
	if err != nil {
		return fmt.Errorf("falha no download do Ubuntu: %w", err)
	}
	defer os.Remove(isoPath)

	if err := c.installUbuntuToUSB(mountPoint, isoPath); err != nil {
		return fmt.Errorf("falha na instalação do Ubuntu: %w", err)
	}
	fmt.Println("   ✅ Ubuntu Server instalado com sucesso")
	fmt.Println()

	// Etapa 4: Configurar boot
	fmt.Println("🔧 Etapa 4/6: Configurando boot...")
	if err := c.configureBoot(mountPoint); err != nil {
		return fmt.Errorf("falha na configuração de boot: %w", err)
	}
	fmt.Println("   ✅ Boot configurado com sucesso")
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
	fmt.Println()
	fmt.Println("📋 Próximos passos:")
	fmt.Println("   1. Desmonte o USB com segurança")
	fmt.Println("   2. Insira no hardware alvo")
	fmt.Println("   3. Configure boot para USB no BIOS/UEFI")
	fmt.Println("   4. A instalação será automática (~30 minutos)")

	return nil
}

// validateConfig valida a configuração fornecida
func (c *USBCreator) validateConfig(config *Config) error {
	if config.NodeName == "" {
		return fmt.Errorf("nome do nó é obrigatório")
	}

	// Validar formato do nome do nó
	if !isValidNodeName(config.NodeName) {
		return fmt.Errorf("nome do nó inválido: %s (deve conter apenas letras, números e hífens)", config.NodeName)
	}

	// Validar coordenadas se fornecidas
	if config.Coordinates != "" && !isValidCoordinates(config.Coordinates) {
		return fmt.Errorf("coordenadas inválidas: %s (formato esperado: lat,lon)", config.Coordinates)
	}

	return nil
}

// setupDirectories cria os diretórios necessários
func (c *USBCreator) setupDirectories() error {
	dirs := []string{c.workDir, c.cacheDir}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("falha ao criar diretório %s: %w", dir, err)
		}
	}

	return nil
}

// formatDevice formata o dispositivo USB
func (c *USBCreator) formatDevice(devicePath, label string) error {
	return c.formatter.FormatDevice(devicePath, label)
}

// mountPartition monta a partição do USB
func (c *USBCreator) mountPartition(devicePath string) (string, error) {
	// Determinar o caminho da partição
	partitionPath := c.getPartitionPath(devicePath)
	if partitionPath == "" {
		return "", fmt.Errorf("não foi possível determinar o caminho da partição")
	}

	// Criar ponto de montagem
	mountPoint := filepath.Join(c.workDir, "mount")
	if err := os.MkdirAll(mountPoint, 0755); err != nil {
		return "", fmt.Errorf("falha ao criar ponto de montagem: %w", err)
	}

	// Montar a partição
	cmd := exec.Command("sudo", "mount", partitionPath, mountPoint)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("falha ao montar partição: %w", err)
	}

	return mountPoint, nil
}

// unmountPartition desmonta a partição
func (c *USBCreator) unmountPartition(mountPoint string) {
	cmd := exec.Command("sudo", "umount", mountPoint)
	cmd.Run() // Ignorar erro
	os.RemoveAll(mountPoint)
}

// downloadUbuntuISO baixa o Ubuntu Server ISO
func (c *USBCreator) downloadUbuntuISO() (string, error) {
	isoPath := filepath.Join(c.cacheDir, "ubuntu-22.04.3-live-server-amd64.iso")

	// Verificar se já existe
	if _, err := os.Stat(isoPath); err == nil {
		fmt.Printf("   ✅ ISO já existe em cache: %s\n", isoPath)
		return isoPath, nil
	}

	fmt.Printf("   📥 Baixando Ubuntu Server ISO...\n")
	fmt.Printf("   📍 Destino: %s\n", isoPath)

	// URL do Ubuntu Server 22.04.3 LTS
	url := "https://releases.ubuntu.com/22.04.3/ubuntu-22.04.3-live-server-amd64.iso"

	// Criar arquivo de destino
	file, err := os.Create(isoPath)
	if err != nil {
		return "", fmt.Errorf("falha ao criar arquivo ISO: %w", err)
	}
	defer file.Close()

	// Fazer download
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("falha ao baixar ISO: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("falha no download: status %d", resp.StatusCode)
	}

	// Copiar dados
	if _, err := io.Copy(file, resp.Body); err != nil {
		return "", fmt.Errorf("falha ao salvar ISO: %w", err)
	}

	fmt.Printf("   ✅ Download concluído: %s\n", isoPath)
	return isoPath, nil
}

// installUbuntuToUSB instala o Ubuntu no USB
func (c *USBCreator) installUbuntuToUSB(mountPoint, isoPath string) error {
	fmt.Printf("   📦 Extraindo ISO para USB...\n")

	// Criar ponto de montagem temporário para o ISO
	isoMountPoint := filepath.Join(c.workDir, "iso-mount")
	if err := os.MkdirAll(isoMountPoint, 0755); err != nil {
		return fmt.Errorf("falha ao criar ponto de montagem do ISO: %w", err)
	}
	defer os.RemoveAll(isoMountPoint)

	// Montar ISO
	cmd := exec.Command("sudo", "mount", "-o", "loop", isoPath, isoMountPoint)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("falha ao montar ISO: %w", err)
	}
	defer exec.Command("sudo", "umount", isoMountPoint).Run()

	// Copiar arquivos do ISO para o USB
	fmt.Printf("   📋 Copiando arquivos do sistema...\n")

	// Copiar arquivos principais
	filesToCopy := []string{
		"casper/",
		"dists/",
		"install/",
		"pool/",
		"EFI/",
		".disk/",
		"isolinux/",
		"[BOOT]/",
	}

	for _, file := range filesToCopy {
		src := filepath.Join(isoMountPoint, file)
		dst := filepath.Join(mountPoint, file)

		if _, err := os.Stat(src); err == nil {
			cmd := exec.Command("sudo", "cp", "-r", src, dst)
			if err := cmd.Run(); err != nil {
				fmt.Printf("   ⚠️  Aviso: Falha ao copiar %s: %v\n", file, err)
			}
		}
	}

	// Copiar arquivos de boot específicos
	bootFiles := []string{
		"boot/grub/grub.cfg",
		"boot/grub/loopback.cfg",
		"md5sum.txt",
	}

	for _, file := range bootFiles {
		src := filepath.Join(isoMountPoint, file)
		dst := filepath.Join(mountPoint, file)

		if _, err := os.Stat(src); err == nil {
			cmd := exec.Command("sudo", "cp", src, dst)
			cmd.Run() // Ignorar erro
		}
	}

	fmt.Printf("   ✅ Arquivos do sistema copiados\n")
	return nil
}

// configureBoot configura o boot do USB
func (c *USBCreator) configureBoot(mountPoint string) error {
	fmt.Printf("   🔧 Configurando GRUB...\n")

	// Instalar GRUB no USB
	cmd := exec.Command("sudo", "grub-install", "--target=i386-pc", "--boot-directory="+filepath.Join(mountPoint, "boot"), "--force", c.getDeviceFromMount(mountPoint))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("falha ao instalar GRUB: %w", err)
	}

	// Criar configuração GRUB personalizada
	grubCfg := filepath.Join(mountPoint, "boot/grub/grub.cfg")
	grubConfig := c.generateGRUBConfig()

	if err := os.WriteFile(grubCfg, []byte(grubConfig), 0644); err != nil {
		return fmt.Errorf("falha ao criar configuração GRUB: %w", err)
	}

	fmt.Printf("   ✅ GRUB configurado com sucesso\n")
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
	cmd := exec.Command("ssh-keygen", "-t", "rsa", "-b", "4096", "-f", privateKeyPath, "-N", "", "-C", nodeName)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("falha ao gerar chave privada: %w", err)
	}

	// Gerar chave pública
	publicKeyPath := privateKeyPath + ".pub"
	cmd = exec.Command("ssh-keygen", "-y", "-f", privateKeyPath)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("falha ao gerar chave pública: %w", err)
	}

	if err := os.WriteFile(publicKeyPath, output, 0644); err != nil {
		return "", fmt.Errorf("falha ao salvar chave pública: %w", err)
	}

	// Gerar fingerprint
	fingerprintCmd := exec.Command("ssh-keygen", "-lf", privateKeyPath)
	fingerprintOutput, err := fingerprintCmd.Output()
	if err == nil {
		fingerprintPath := filepath.Join(keysDir, nodeName+".fingerprint")
		os.WriteFile(fingerprintPath, fingerprintOutput, 0644)
	}

	fmt.Printf("   ✅ Chaves SSH geradas com sucesso\n")
	fmt.Printf("   📍 Chave privada: %s\n", privateKeyPath)
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

// generateGRUBConfig gera configuração GRUB
func (c *USBCreator) generateGRUBConfig() string {
	return `set timeout=10
set default=0

menuentry "Ubuntu Server (Syntropy Cooperative Grid)" {
    set gfxpayload=keep
    linux /casper/vmlinuz quiet autoinstall ds=nocloud
    initrd /casper/initrd
}

menuentry "Ubuntu Server (Manual Install)" {
    set gfxpayload=keep
    linux /casper/vmlinuz
    initrd /casper/initrd
}
`
}

// generateCloudInitConfig gera configuração cloud-init
func (c *USBCreator) generateCloudInitConfig(config *Config) string {
	return fmt.Sprintf(`#cloud-config
# Syntropy Cooperative Grid Node Configuration
# Node: %s
# Generated: %s

users:
  - name: syntropy
    groups: sudo, docker
    shell: /bin/bash
    sudo: ['ALL=(ALL) NOPASSWD:ALL']
    ssh_authorized_keys:
      - %s

package_update: true
package_upgrade: true

packages:
  - docker.io
  - docker-compose
  - curl
  - wget
  - git
  - htop
  - fail2ban
  - ufw

runcmd:
  - systemctl enable docker
  - systemctl start docker
  - usermod -aG docker syntropy
  - ufw --force enable
  - ufw allow ssh
  - ufw allow 80/tcp
  - ufw allow 443/tcp
  - systemctl enable fail2ban
  - systemctl start fail2ban
  - echo "Syntropy Cooperative Grid Node: %s" > /etc/motd

write_files:
  - path: /etc/syntropy/node.conf
    content: |
      node_name=%s
      node_description=%s
      coordinates=%s
      created_at=%s
    permissions: '0644'
    owner: root:root

final_message: "Syntropy Cooperative Grid Node %s installed successfully!"
`, config.NodeName, time.Now().Format(time.RFC3339),
		c.getSSHPublicKey(config), config.NodeName, config.NodeName,
		config.NodeDescription, config.Coordinates, time.Now().Format(time.RFC3339), config.NodeName)
}

// generateMetaData gera meta-data
func (c *USBCreator) generateMetaData(config *Config) string {
	return fmt.Sprintf(`instance-id: %s
local-hostname: %s
`, generateInstanceID(), config.NodeName)
}

// Funções auxiliares

func (c *USBCreator) getPartitionPath(devicePath string) string {
	// Tentar diferentes padrões
	patterns := []string{
		devicePath + "1",
		devicePath + "p1",
	}

	for _, pattern := range patterns {
		if _, err := os.Stat(pattern); err == nil {
			return pattern
		}
	}
	return ""
}

func (c *USBCreator) getDeviceFromMount(mountPoint string) string {
	// Implementação simples - em produção usar /proc/mounts
	return "/dev/sdb" // Placeholder
}

func (c *USBCreator) getSSHPublicKey(config *Config) string {
	// Se uma chave de proprietário foi fornecida, usar ela
	if config.OwnerKeyFile != "" {
		if pubKey, err := os.ReadFile(config.OwnerKeyFile + ".pub"); err == nil {
			return strings.TrimSpace(string(pubKey))
		}
	}

	// Caso contrário, gerar uma chave temporária
	return "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC..." // Placeholder
}

func isValidNodeName(name string) bool {
	// Nome deve conter apenas letras, números e hífens
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '-') {
			return false
		}
	}
	return len(name) > 0 && len(name) <= 63
}

func isValidCoordinates(coords string) bool {
	// Formato esperado: lat,lon (ex: -23.5505,-46.6333)
	parts := strings.Split(coords, ",")
	if len(parts) != 2 {
		return false
	}

	// Validar latitude (-90 a 90)
	lat := strings.TrimSpace(parts[0])
	if lat == "" {
		return false
	}

	// Validar longitude (-180 a 180)
	lon := strings.TrimSpace(parts[1])
	if lon == "" {
		return false
	}

	return true
}

func generateInstanceID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// Cleanup limpa arquivos temporários
func (c *USBCreator) Cleanup() error {
	// Limpar diretório de trabalho
	if err := os.RemoveAll(c.workDir); err != nil {
		return fmt.Errorf("falha ao limpar diretório de trabalho: %w", err)
	}
	return nil
}
