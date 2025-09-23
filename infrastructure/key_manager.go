package infrastructure

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// KeyManager gerencia chaves SSH para Syntropy nodes
type KeyManager struct {
	keyDir string
}

// KeyPair representa um par de chaves SSH
type KeyPair struct {
	PrivateKey  string
	PublicKey   string
	Fingerprint string
	Algorithm   string
	CreatedAt   time.Time
}

// KeyPurpose define o propósito da chave
type KeyPurpose string

const (
	OwnerKey     KeyPurpose = "owner"     // Chave do proprietário (SSH access)
	CommunityKey KeyPurpose = "community" // Chave da comunidade (inter-node communication)
	NodeKey      KeyPurpose = "node"      // Chave específica do nó
)

// NewKeyManager cria um novo gerenciador de chaves
func NewKeyManager(keyDir string) *KeyManager {
	return &KeyManager{
		keyDir: keyDir,
	}
}

// GenerateKeyPair gera um novo par de chaves SSH
func (km *KeyManager) GenerateKeyPair(purpose KeyPurpose, nodeName string) (*KeyPair, error) {
	// Gerar chave ED25519 (recomendada para SSH)
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("falha ao gerar chave ED25519: %w", err)
	}

	// Converter chave privada para formato PEM
	privateKeyDER, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("falha ao codificar chave privada: %w", err)
	}

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyDER,
	})

	// Converter chave pública para formato SSH
	sshPublicKey, err := ssh.NewPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("falha ao converter chave pública: %w", err)
	}

	publicKeySSH := string(ssh.MarshalAuthorizedKey(sshPublicKey))

	// Gerar fingerprint
	fingerprint := ssh.FingerprintSHA256(sshPublicKey)

	return &KeyPair{
		PrivateKey:  string(privateKeyPEM),
		PublicKey:   publicKeySSH,
		Fingerprint: fingerprint,
		Algorithm:   "ed25519",
		CreatedAt:   time.Now(),
	}, nil
}

// GenerateRSAKeyPair gera um par de chaves RSA (para compatibilidade)
func (km *KeyManager) GenerateRSAKeyPair(purpose KeyPurpose, nodeName string, bits int) (*KeyPair, error) {
	if bits == 0 {
		bits = 4096 // Tamanho padrão
	}

	// Gerar chave RSA
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, fmt.Errorf("falha ao gerar chave RSA: %w", err)
	}

	// Converter chave privada para formato PEM
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	// Gerar chave pública SSH
	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("falha ao gerar chave pública SSH: %w", err)
	}

	publicKeySSH := string(ssh.MarshalAuthorizedKey(publicKey))

	// Gerar fingerprint
	fingerprint := ssh.FingerprintSHA256(publicKey)

	return &KeyPair{
		PrivateKey:  string(privateKeyPEM),
		PublicKey:   publicKeySSH,
		Fingerprint: fingerprint,
		Algorithm:   fmt.Sprintf("rsa-%d", bits),
		CreatedAt:   time.Now(),
	}, nil
}

// SaveKeyPair salva um par de chaves em arquivos
func (km *KeyManager) SaveKeyPair(keyPair *KeyPair, purpose KeyPurpose, nodeName string) error {
	// Criar diretório se não existir
	if err := os.MkdirAll(km.keyDir, 0700); err != nil {
		return fmt.Errorf("falha ao criar diretório de chaves: %w", err)
	}

	// Nomes dos arquivos
	keyFileName := fmt.Sprintf("%s-%s.key", nodeName, purpose)
	pubFileName := fmt.Sprintf("%s-%s.key.pub", nodeName, purpose)
	fingerprintFileName := fmt.Sprintf("%s-%s.fingerprint", nodeName, purpose)

	keyPath := filepath.Join(km.keyDir, keyFileName)
	pubPath := filepath.Join(km.keyDir, pubFileName)
	fingerprintPath := filepath.Join(km.keyDir, fingerprintFileName)

	// Salvar chave privada
	if err := os.WriteFile(keyPath, []byte(keyPair.PrivateKey), 0600); err != nil {
		return fmt.Errorf("falha ao salvar chave privada: %w", err)
	}

	// Salvar chave pública
	if err := os.WriteFile(pubPath, []byte(keyPair.PublicKey), 0644); err != nil {
		return fmt.Errorf("falha ao salvar chave pública: %w", err)
	}

	// Salvar fingerprint
	fingerprintData := fmt.Sprintf("%s %s %s\n", keyPair.Fingerprint, keyPair.Algorithm, keyFileName)
	if err := os.WriteFile(fingerprintPath, []byte(fingerprintData), 0644); err != nil {
		return fmt.Errorf("falha ao salvar fingerprint: %w", err)
	}

	return nil
}

// LoadKeyPair carrega um par de chaves de arquivos
func (km *KeyManager) LoadKeyPair(purpose KeyPurpose, nodeName string) (*KeyPair, error) {
	keyFileName := fmt.Sprintf("%s-%s.key", nodeName, purpose)
	pubFileName := fmt.Sprintf("%s-%s.key.pub", nodeName, purpose)
	fingerprintFileName := fmt.Sprintf("%s-%s.fingerprint", nodeName, purpose)

	keyPath := filepath.Join(km.keyDir, keyFileName)
	pubPath := filepath.Join(km.keyDir, pubFileName)
	fingerprintPath := filepath.Join(km.keyDir, fingerprintFileName)

	// Carregar chave privada
	privateKeyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("falha ao carregar chave privada: %w", err)
	}

	// Carregar chave pública
	publicKeyData, err := os.ReadFile(pubPath)
	if err != nil {
		return nil, fmt.Errorf("falha ao carregar chave pública: %w", err)
	}

	// Carregar fingerprint
	fingerprintData, err := os.ReadFile(fingerprintPath)
	if err != nil {
		return nil, fmt.Errorf("falha ao carregar fingerprint: %w", err)
	}

	// Parse do fingerprint
	fingerprintInfo := string(fingerprintData)
	fingerprint := strings.Fields(fingerprintInfo)[0]
	algorithm := strings.Fields(fingerprintInfo)[1]

	return &KeyPair{
		PrivateKey:  string(privateKeyData),
		PublicKey:   string(publicKeyData),
		Fingerprint: fingerprint,
		Algorithm:   algorithm,
		CreatedAt:   time.Now(), // Placeholder
	}, nil
}

// LoadExistingKeyPair carrega chaves existentes de um arquivo
func (km *KeyManager) LoadExistingKeyPair(keyFilePath string) (*KeyPair, error) {
	// Carregar chave privada
	privateKeyData, err := os.ReadFile(keyFilePath)
	if err != nil {
		return nil, fmt.Errorf("falha ao carregar chave existente: %w", err)
	}

	// Carregar chave pública correspondente
	pubFilePath := keyFilePath + ".pub"
	publicKeyData, err := os.ReadFile(pubFilePath)
	if err != nil {
		return nil, fmt.Errorf("falha ao carregar chave pública correspondente: %w", err)
	}

	// Parse da chave privada para obter informações
	block, _ := pem.Decode(privateKeyData)
	if block == nil {
		return nil, fmt.Errorf("falha ao fazer parse da chave privada")
	}

	// Gerar fingerprint da chave pública
	sshPublicKey, _, _, _, err := ssh.ParseAuthorizedKey(publicKeyData)
	if err != nil {
		return nil, fmt.Errorf("falha ao fazer parse da chave pública: %w", err)
	}

	fingerprint := ssh.FingerprintSHA256(sshPublicKey)

	// Determinar algoritmo
	var algorithm string
	switch block.Type {
	case "RSA PRIVATE KEY":
		algorithm = "rsa"
	case "PRIVATE KEY":
		// Pode ser ED25519 ou outro formato
		algorithm = "ed25519"
	default:
		algorithm = "unknown"
	}

	return &KeyPair{
		PrivateKey:  string(privateKeyData),
		PublicKey:   string(publicKeyData),
		Fingerprint: fingerprint,
		Algorithm:   algorithm,
		CreatedAt:   time.Now(), // Placeholder
	}, nil
}

// ListKeys lista todas as chaves disponíveis
func (km *KeyManager) ListKeys() ([]string, error) {
	var keys []string

	files, err := os.ReadDir(km.keyDir)
	if err != nil {
		return nil, fmt.Errorf("falha ao listar diretório de chaves: %w", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".key") && !strings.HasSuffix(file.Name(), ".pub") {
			keys = append(keys, file.Name())
		}
	}

	return keys, nil
}

// DeleteKeyPair remove um par de chaves
func (km *KeyManager) DeleteKeyPair(purpose KeyPurpose, nodeName string) error {
	keyFileName := fmt.Sprintf("%s-%s.key", nodeName, purpose)
	pubFileName := fmt.Sprintf("%s-%s.key.pub", nodeName, purpose)
	fingerprintFileName := fmt.Sprintf("%s-%s.fingerprint", nodeName, purpose)

	keyPath := filepath.Join(km.keyDir, keyFileName)
	pubPath := filepath.Join(km.keyDir, pubFileName)
	fingerprintPath := filepath.Join(km.keyDir, fingerprintFileName)

	// Remover arquivos
	files := []string{keyPath, pubPath, fingerprintPath}
	for _, file := range files {
		if err := os.Remove(file); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("falha ao remover arquivo %s: %w", file, err)
		}
	}

	return nil
}

// GenerateKeyInstallationCommands gera comandos para instalar chaves no cloud-init
func (km *KeyManager) GenerateKeyInstallationCommands(ownerKey, ownerPub, communityKey, communityPub *KeyPair) string {
	return fmt.Sprintf(`curtin in-target -- bash -c '
# Install owner key (SSH access and management)
mkdir -p /opt/syntropy/identity/owner
cat > /opt/syntropy/identity/owner/private.key << "OWNER_KEY_EOF"
%s
OWNER_KEY_EOF

cat > /opt/syntropy/identity/owner/public.key << "OWNER_PUB_EOF"
%s
OWNER_PUB_EOF

# Install community key (inter-node communication)
mkdir -p /opt/syntropy/identity/community
cat > /opt/syntropy/identity/community/private.key << "COMMUNITY_KEY_EOF"
%s
COMMUNITY_KEY_EOF

cat > /opt/syntropy/identity/community/public.key << "COMMUNITY_PUB_EOF"
%s
COMMUNITY_PUB_EOF

# Set proper permissions
chmod 600 /opt/syntropy/identity/owner/private.key
chmod 600 /opt/syntropy/identity/community/private.key
chmod 644 /opt/syntropy/identity/owner/public.key
chmod 644 /opt/syntropy/identity/community/public.key

# Configure SSH access
mkdir -p /home/admin/.ssh
cp /opt/syntropy/identity/owner/public.key /home/admin/.ssh/authorized_keys
chmod 600 /home/admin/.ssh/authorized_keys
chown admin:admin /home/admin/.ssh/authorized_keys
chown -R admin:admin /opt/syntropy/
'`,
		ownerKey.PrivateKey, ownerPub.PublicKey, communityKey.PrivateKey, communityPub.PublicKey)
}

// GenerateMetadataCreationCommands gera comandos para criar metadados do nó
func (km *KeyManager) GenerateMetadataCreationCommands(nodeName, coordinates, description string) string {
	return fmt.Sprintf(`curtin in-target -- bash -c '
# Create comprehensive node metadata
mkdir -p /opt/syntropy/metadata
cat > /opt/syntropy/metadata/node.json << "NODE_METADATA_EOF"
{
  "metadata_version": "2.0",
  "node_info": {
    "node_name": "%s",
    "description": "%s",
    "installation_time": "$(date -u +%%Y-%%m-%%dT%%H:%%M:%%SZ)",
    "platform_version": "2.0.0",
    "platform_type": "syntropy_cooperative_grid"
  },
  "geographic_info": {
    "coordinates": {
      "formatted": "%s"
    }
  },
  "security": {
    "owner_key_fingerprint": "%s",
    "community_key_fingerprint": "%s",
    "ssh_port": 22,
    "authentication_method": "key_only",
    "firewall_enabled": true
  },
  "platform": {
    "type": "syntropy_cooperative_grid",
    "status": "installed",
    "services": {
      "docker": "enabled",
      "ssh": "enabled",
      "prometheus_exporter": "enabled",
      "firewall": "enabled"
    }
  }
}
NODE_METADATA_EOF

chmod 644 /opt/syntropy/metadata/node.json
chown admin:admin /opt/syntropy/metadata/node.json
'`,
		nodeName, description, coordinates, "OWNER_FINGERPRINT_PLACEHOLDER", "COMMUNITY_FINGERPRINT_PLACEHOLDER")
}
