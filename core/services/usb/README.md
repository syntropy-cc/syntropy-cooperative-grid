# 🔌 USB Service - Syntropy Cooperative Grid

> **Serviço de Gerenciamento de USB para Criação de Nós Bootáveis**

## 📋 **Visão Geral**

O **USB Service** é um componente do Core Layer da Syntropy Cooperative Grid responsável por gerenciar dispositivos USB e criar USBs bootáveis para deployment de nós. Este serviço segue os princípios arquiteturais do projeto, fornecendo uma interface limpa e abstraindo a complexidade da formatação e configuração de dispositivos USB.

## 🏗️ **Arquitetura do Serviço**

### **Estrutura do Módulo**
```
core/services/usb/
├── detector.go          # Detecção de dispositivos USB
├── formatter.go         # Formatação de dispositivos USB
├── creator.go           # Criação de USBs bootáveis
├── go.mod              # Dependências do módulo
└── README.md           # Documentação
```

### **Interfaces Principais**

#### **1. Detector Interface**
```go
type Detector interface {
    DetectDevices() ([]USBDevice, error)
    ValidateDevice(devicePath string) error
    IsSystemDisk(devicePath string) bool
}
```

#### **2. Formatter Interface**
```go
type Formatter interface {
    FormatDevice(devicePath, label string) error
    UnmountDevice(devicePath string) error
    WipeDevice(devicePath string) error
    CreatePartitionTable(devicePath string) error
    CreatePartition(devicePath string) error
    FormatPartition(partitionPath, label string) error
}
```

#### **3. Creator Interface**
```go
type Creator interface {
    CreateUSB(devicePath string, config *Config) error
    Cleanup() error
}
```

## 🔧 **Funcionalidades**

### **1. Detecção Multi-Plataforma**
- **Linux Nativo**: Usa `lsblk` para detectar dispositivos USB
- **WSL (Windows Subsystem for Linux)**: Integração híbrida com PowerShell
- **Windows Nativo**: Usa PowerShell para detecção via WMI

### **2. Formatação Segura**
- **Limpeza completa**: Remove assinaturas de sistema de arquivos
- **Tabela de partições MBR**: Compatibilidade máxima
- **Sistema FAT32**: Suporte universal para boot
- **Validação de segurança**: Proteção contra formatação de discos do sistema

### **3. Criação de USB Bootável**
- **Infrastructure as Code**: Usa templates cloud-init parametrizáveis
- **Geração de chaves SSH**: ED25519 e RSA 4096-bit
- **Configuração automática**: Ubuntu Server com auto-install
- **Metadados estruturados**: JSON com informações completas do nó

## 🚀 **Como Usar**

### **Via CLI Interface**
```bash
# Listar dispositivos USB disponíveis
syntropy usb list

# Criar USB com auto-detecção
syntropy usb create --auto-detect --node-name "node-01"

# Criar USB especificando dispositivo
syntropy usb create /dev/sdb --node-name "node-01" --coordinates "-23.5505,-46.6333"

# Formatar USB
syntropy usb format /dev/sdb --label "SYNTROPY"
```

### **Via Programmatic Interface**
```go
package main

import (
    "fmt"
    "syntropy-cc/cooperative-grid/core/services/usb"
)

func main() {
    // Criar detector
    detector := usb.NewDetector()
    
    // Detectar dispositivos
    devices, err := detector.DetectDevices()
    if err != nil {
        panic(err)
    }
    
    // Criar config
    config := &usb.Config{
        NodeName:        "node-01",
        NodeDescription: "Production node",
        Coordinates:     "-23.5505,-46.6333",
        Label:           "SYNTROPY",
    }
    
    // Criar USB
    creator := usb.NewCreator("/tmp/work", "/tmp/cache")
    defer creator.Cleanup()
    
    err = creator.CreateUSB("/dev/sdb", config)
    if err != nil {
        panic(err)
    }
    
    fmt.Println("USB criado com sucesso!")
}
```

## 🔐 **Segurança**

### **Validações de Segurança**
1. **Detecção de disco do sistema**: Previne formatação acidental
2. **Validação de dispositivo**: Verifica se é um dispositivo de bloco válido
3. **Confirmação do usuário**: Para operações destrutivas
4. **Permissões adequadas**: Requer privilégios sudo no Linux

### **Geração de Chaves SSH**
- **Algoritmos suportados**: ED25519 (recomendado), RSA 4096-bit
- **Chaves do proprietário**: Acesso SSH e gerenciamento
- **Chaves da comunidade**: Comunicação entre nós
- **Armazenamento seguro**: Permissões 600 para chaves privadas

## 🌍 **Suporte Multi-Plataforma**

### **Linux**
- **Ferramentas**: `lsblk`, `parted`, `mkfs.fat`, `wipefs`, `sgdisk`
- **Montagem**: `/tmp/syntropy-work/mount`
- **Permissões**: Requer sudo para operações de baixo nível

### **WSL (Windows Subsystem for Linux)**
- **Detecção híbrida**: PowerShell + lsblk
- **Mapeamento**: PhysicalDriveN → /dev/sdX
- **Validações específicas**: Proteção contra formatação do disco WSL

### **Windows**
- **PowerShell**: WMI queries para detecção
- **Formatação**: Clear-Disk, New-Partition, Format-Volume
- **Dispositivos**: PhysicalDriveN notation

## 📊 **Estrutura de Dados**

### **USBDevice**
```go
type USBDevice struct {
    Path      string `json:"path"`       // Caminho do dispositivo
    Size      string `json:"size"`       // Tamanho em formato legível
    SizeGB    int    `json:"size_gb"`    // Tamanho em GB
    Model     string `json:"model"`      // Modelo do dispositivo
    Vendor    string `json:"vendor"`     // Fabricante
    Serial    string `json:"serial"`     // Número de série
    Removable bool   `json:"removable"`  // É removível?
    Platform  string `json:"platform"`  // Plataforma (linux/wsl/windows)
}
```

### **Config**
```go
type Config struct {
    NodeName        string `json:"node_name"`        // Nome do nó
    NodeDescription string `json:"node_description"` // Descrição
    Coordinates     string `json:"coordinates"`      // Coordenadas geográficas
    OwnerKeyFile    string `json:"owner_key_file"`   // Chave do proprietário
    Label           string `json:"label"`            // Rótulo do sistema de arquivos
}
```

## 🔄 **Workflow de Criação**

### **Processo Completo**
1. **Validação**: Verificar se dispositivo é válido e seguro
2. **Montagem**: Montar dispositivo para verificação
3. **Formatação**: Limpar, particionar e formatar como FAT32
4. **Remontagem**: Remontar dispositivo formatado
5. **Geração de chaves**: Criar chaves SSH (owner + community)
6. **Cloud-init**: Gerar configuração usando templates IaC

### **Infrastructure as Code**
- **Templates**: Usa templates parametrizáveis em `infrastructure/cloud-init/`
- **Geração dinâmica**: user-data, meta-data, network-config
- **Chaves SSH**: Integração com `infrastructure/key_manager`
- **Metadados**: JSON estruturado com informações do nó

## 📁 **Arquivos Gerados**

### **Estrutura no USB**
```
USB_ROOT/
├── cloud-init/
│   ├── user-data          # Configuração principal
│   ├── meta-data          # Metadados do nó
│   └── network-config     # Configuração de rede
├── syntropy/
│   ├── keys/
│   │   ├── node-01.key    # Chave privada SSH
│   │   └── node-01.key.pub # Chave pública SSH
│   └── metadata/
│       └── node.json      # Metadados estruturados
└── platform/
    └── templates/         # Templates Kubernetes
```

## 🐛 **Troubleshooting**

### **Problemas Comuns**

**Dispositivo não detectado:**
```bash
# Verificar dispositivos de bloco
lsblk

# Verificar permissões
sudo ls -la /dev/sd*

# Verificar se está montado
mount | grep sd
```

**Falha na formatação:**
```bash
# Desmontar todas as partições
sudo umount /dev/sdb*

# Verificar processos usando o dispositivo
sudo lsof | grep sdb

# Limpar manualmente
sudo wipefs -a /dev/sdb
```

**Permissões insuficientes:**
```bash
# Adicionar usuário ao grupo disk
sudo usermod -a -G disk $USER

# Verificar grupos
groups $USER
```

### **Logs e Debug**
```bash
# Logs do sistema
sudo journalctl -f

# Logs específicos do USB
dmesg | grep -i usb

# Verificar montagens
cat /proc/mounts | grep sd
```

## 🔧 **Desenvolvimento**

### **Adicionar Nova Plataforma**
1. Implementar interfaces `Detector`, `Formatter`
2. Adicionar case no `NewDetector()`, `NewFormatter()`
3. Testar detecção e formatação
4. Atualizar documentação

### **Testes**
```bash
# Executar testes unitários
go test ./...

# Testar detecção
go run detector_test.go

# Testar formatação (CUIDADO!)
go run formatter_test.go
```

### **Dependências**
- `syntropy-cc/cooperative-grid/infrastructure`: Gerenciamento de templates e chaves
- `golang.org/x/crypto/ssh`: Geração de chaves SSH
- Sistema: `lsblk`, `parted`, `mkfs.fat`, `wipefs`, `sgdisk`

## 📚 **Referências**

- [Cloud-Init Documentation](https://cloudinit.readthedocs.io/)
- [Ubuntu Server Auto-Install](https://ubuntu.com/server/docs/install/autoinstall)
- [USB Boot Standards](https://en.wikipedia.org/wiki/USB_boot)
- [GPT vs MBR](https://en.wikipedia.org/wiki/GUID_Partition_Table)
- [SSH Key Generation](https://www.ssh.com/academy/ssh/keygen)

---

**Este serviço é parte integrante da arquitetura Syntropy Cooperative Grid e segue os princípios de separação de responsabilidades, cross-platform compatibility e Infrastructure as Code.**
