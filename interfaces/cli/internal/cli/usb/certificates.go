package usb

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"

	"syntropy-cc/cooperative-grid/infrastructure"
)

// generateSSHKeyPair gera um par de chaves SSH usando o KeyManager centralizado
func generateSSHKeyPair(nodeName string) (string, string, error) {
	// Obter diretório de chaves
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", "", fmt.Errorf("erro ao obter diretório home: %w", err)
	}
	keyDir := filepath.Join(homeDir, ".syntropy", "keys")

	// Criar KeyManager
	keyManager := infrastructure.NewKeyManager(keyDir)

	// Gerar par de chaves ED25519 (mais seguro)
	keyPair, err := keyManager.GenerateKeyPair(infrastructure.NodeKey, nodeName)
	if err != nil {
		// Fallback para RSA se ED25519 falhar
		keyPair, err = keyManager.GenerateRSAKeyPair(infrastructure.NodeKey, nodeName, 2048)
		if err != nil {
			return "", "", fmt.Errorf("erro ao gerar chaves SSH: %w", err)
		}
	}

	// Salvar chaves no diretório centralizado
	if err := keyManager.SaveKeyPair(keyPair, infrastructure.NodeKey, nodeName); err != nil {
		return "", "", fmt.Errorf("erro ao salvar chaves SSH: %w", err)
	}

	return keyPair.PrivateKey, keyPair.PublicKey, nil
}

// loadExistingSSHKeyPair carrega um par de chaves SSH existente
func loadExistingSSHKeyPair(nodeName string) (string, string, error) {
	// Obter diretório de chaves
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", "", fmt.Errorf("erro ao obter diretório home: %w", err)
	}
	keyDir := filepath.Join(homeDir, ".syntropy", "keys")

	// Criar KeyManager
	keyManager := infrastructure.NewKeyManager(keyDir)

	// Carregar par de chaves existente
	keyPair, err := keyManager.LoadKeyPair(infrastructure.NodeKey, nodeName)
	if err != nil {
		return "", "", fmt.Errorf("erro ao carregar chaves SSH existentes: %w", err)
	}

	return keyPair.PrivateKey, keyPair.PublicKey, nil
}

// generateCertificates gera certificados TLS para o nó
func generateCertificates(nodeName string, ownerKey string) (*Certificates, error) {
	// Gerar chave CA
	caKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar chave CA: %w", err)
	}

	// Criar certificado CA
	caTemplate := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:  []string{"Syntropy Cooperative Grid"},
			Country:       []string{"BR"},
			Province:      []string{""},
			Locality:      []string{"São Paulo"},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0), // 10 anos
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	caCert, err := x509.CreateCertificate(rand.Reader, &caTemplate, &caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar certificado CA: %w", err)
	}

	// Gerar chave do nó
	nodeKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar chave do nó: %w", err)
	}

	// Criar certificado do nó
	nodeTemplate := x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			Organization:  []string{"Syntropy Cooperative Grid"},
			Country:       []string{"BR"},
			Province:      []string{""},
			Locality:      []string{"São Paulo"},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
			CommonName:    nodeName,
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(1, 0, 0), // 1 ano
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1)},
		DNSNames:    []string{nodeName, "localhost"},
	}

	nodeCert, err := x509.CreateCertificate(rand.Reader, &nodeTemplate, &caTemplate, &nodeKey.PublicKey, caKey)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar certificado do nó: %w", err)
	}

	// Codificar certificados em PEM
	caKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caKey),
	}

	caCertPEM := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caCert,
	}

	nodeKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(nodeKey),
	}

	nodeCertPEM := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: nodeCert,
	}

	return &Certificates{
		CAKey:    pem.EncodeToMemory(caKeyPEM),
		CACert:   pem.EncodeToMemory(caCertPEM),
		NodeKey:  pem.EncodeToMemory(nodeKeyPEM),
		NodeCert: pem.EncodeToMemory(nodeCertPEM),
	}, nil
}

// saveCertificates salva os certificados no diretório de trabalho
func saveCertificates(certs *Certificates, workDir string) (string, string, string, string, error) {
	certDir := filepath.Join(workDir, "certs")
	if err := os.MkdirAll(certDir, 0755); err != nil {
		return "", "", "", "", fmt.Errorf("erro ao criar diretório de certificados: %w", err)
	}

	caKeyPath := filepath.Join(certDir, "ca.key")
	caCertPath := filepath.Join(certDir, "ca.crt")
	nodeKeyPath := filepath.Join(certDir, "node.key")
	nodeCertPath := filepath.Join(certDir, "node.crt")

	if err := os.WriteFile(caKeyPath, certs.CAKey, 0600); err != nil {
		return "", "", "", "", fmt.Errorf("erro ao salvar chave CA: %w", err)
	}

	if err := os.WriteFile(caCertPath, certs.CACert, 0644); err != nil {
		return "", "", "", "", fmt.Errorf("erro ao salvar certificado CA: %w", err)
	}

	if err := os.WriteFile(nodeKeyPath, certs.NodeKey, 0600); err != nil {
		return "", "", "", "", fmt.Errorf("erro ao salvar chave do nó: %w", err)
	}

	if err := os.WriteFile(nodeCertPath, certs.NodeCert, 0644); err != nil {
		return "", "", "", "", fmt.Errorf("erro ao salvar certificado do nó: %w", err)
	}

	return caKeyPath, caCertPath, nodeKeyPath, nodeCertPath, nil
}
