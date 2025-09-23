# 📱 Módulo USB - Syntropy Cooperative Grid

Este módulo fornece funcionalidades completas para gerenciamento de dispositivos USB, criação de USBs bootáveis e configuração automática de nós da Syntropy Cooperative Grid.

## 📋 Visão Geral

O módulo USB é responsável por:
- **Detecção** de dispositivos USB em diferentes plataformas (Linux, Windows, WSL)
- **Formatação** segura de dispositivos USB
- **Criação** de USBs bootáveis com Ubuntu Server
- **Configuração** automática via cloud-init para nós Syntropy
- **Geração** de certificados TLS e chaves SSH
- **Gerenciamento** de cache de ISOs Ubuntu

## 🏗️ Arquitetura do Módulo

```
usb/
├── usb.go           # Arquivo principal e ponto de entrada
├── types.go         # Estruturas de dados e tipos
├── commands.go      # Comandos CLI e lógica principal
├── platform.go      # Detecção de plataforma e funções comuns
├── linux.go         # Implementações específicas do Linux
├── windows.go       # Implementações específicas do Windows/WSL
├── certificates.go  # Geração de certificados e chaves SSH
├── cloudinit.go     # Configuração do cloud-init
└── utils.go         # Funções auxiliares e utilitários
```

## 📁 Estrutura de Diretórios

O módulo USB utiliza a estrutura padrão `$HOME/.syntropy` para organização:

```
~/.syntropy/
├── cache/                    # Cache de ISOs Ubuntu
│   └── iso/                 # ISOs baixadas automaticamente
├── work/                    # Diretórios de trabalho temporários
│   ├── usb-{timestamp}/     # Trabalho específico de criação USB
│   │   ├── certs/           # Certificados TLS gerados
│   │   ├── cloud-init/      # Arquivos de configuração cloud-init
│   │   └── scripts/         # Scripts copiados
│   └── cidata-mount/        # Pontos de montagem (WSL)
├── nodes/                   # Configurações dos nós
├── keys/                    # Chaves SSH centralizadas ✅
│   ├── {nodeName}-node.key      # Chave privada do nó
│   ├── {nodeName}-node.key.pub  # Chave pública do nó
│   ├── {nodeName}-node.fingerprint # Fingerprint para auditoria
│   ├── {nodeName}-owner.key      # Chave do proprietário
│   └── {nodeName}-community.key  # Chave da comunidade
├── config/                  # Configurações globais
├── scripts/                 # Scripts auxiliares
└── backups/                 # Backups
```

---

## 📄 Documentação por Arquivo

### 🔧 `usb.go` (14 linhas)
**Arquivo principal e ponto de entrada do módulo**

Este arquivo serve como documentação da estrutura do módulo e contém apenas comentários explicativos sobre a organização dos arquivos.

**Funções:**
- Nenhuma função específica - apenas documentação

**Responsabilidades:**
- Documentar a estrutura modular do pacote
- Servir como ponto de referência para desenvolvedores

---

### 🏷️ `types.go` (84 linhas)
**Definições de tipos e estruturas de dados**

Contém todas as estruturas de dados utilizadas pelo módulo USB.

#### Estruturas Principais:

**`USBDevice`**
```go
type USBDevice struct {
    Path        string `json:"path"`         // Caminho do dispositivo
    Size        string `json:"size"`         // Tamanho formatado (ex: "8.5 GB")
    SizeGB      int    `json:"size_gb"`      // Tamanho em GB
    Model       string `json:"model"`        // Modelo do dispositivo
    Vendor      string `json:"vendor"`       // Fabricante
    Serial      string `json:"serial"`       // Número de série
    Removable   bool   `json:"removable"`    // Se é removível
    Platform    string `json:"platform"`     // Plataforma (linux/windows/wsl)
    DiskNumber  int    `json:"disk_number,omitempty"`  // Número do disco (Windows)
    WindowsPath string `json:"windows_path,omitempty"` // Caminho Windows
}
```

**`Config`**
```go
type Config struct {
    NodeName        string `json:"node_name"`         // Nome do nó Syntropy
    NodeDescription string `json:"node_description"`  // Descrição do nó
    Coordinates     string `json:"coordinates"`       // Coordenadas geográficas
    OwnerKeyFile    string `json:"owner_key_file"`    // Arquivo de chave do proprietário
    Label           string `json:"label"`             // Rótulo do sistema de arquivos
    ISOPath         string `json:"iso_path"`          // Caminho da ISO Ubuntu
    DiscoveryServer string `json:"discovery_server"`  // Servidor de descoberta
    SSHPublicKey    string `json:"ssh_public_key"`    // Chave pública SSH
    SSHPrivateKey   string `json:"ssh_private_key"`   // Chave privada SSH
    CreatedBy       string `json:"created_by"`        // Usuário criador
}
```

**`CloudInitConfig`**
```go
type CloudInitConfig struct {
    // Configurações básicas do nó
    NodeName         string
    NodeDescription  string
    Coordinates      string
    OwnerKey         string
    DiscoveryServer  string
    SSHPublicKey     string
    CreatedBy        string
    CreatedAt        string
    
    // Configurações de rede
    InstanceID       string
    Interface        string
    Gateway          string
    NodeIPSuffix     string
    PrimaryInterface string
    MeshGateway      string
    MgmtGateway      string
    
    // Configurações de proxy
    HTTPProxy        string
    HTTPSProxy       string
    
    // Configurações de hardware
    NodeType         string
    HardwareType     string
    CPUCores         int
    MemoryGB         int
    StorageGB        int
    
    // Configurações de papel
    InitialRole      string
    CanBeLeader      bool
    CanBeWorker      bool
    
    // Caminhos de certificados
    NodeCertPath     string
    NodeKeyPath      string
    CACertPath       string
    
    // Certificados PEM para write_files
    CACertPEM        string
    NodeCertPEM      string
    NodeKeyPEM       string
}
```

**`Certificates`**
```go
type Certificates struct {
    CAKey    []byte  // Chave privada da CA
    CACert   []byte  // Certificado da CA
    NodeKey  []byte  // Chave privada do nó
    NodeCert []byte  // Certificado do nó
}
```

**`WindowsDisk`**
```go
type WindowsDisk struct {
    Number       int    `json:"Number"`        // Número do disco
    FriendlyName string `json:"FriendlyName"`  // Nome amigável
    Size         int64  `json:"Size"`          // Tamanho em bytes
    SerialNumber string `json:"SerialNumber"`  // Número de série
    BusType      string `json:"BusType"`       // Tipo de barramento
    Model        string `json:"Model"`         // Modelo do disco
}
```

---

### ⚡ `commands.go` (298 linhas)
**Comandos CLI e lógica principal de criação de USB**

Contém a implementação dos comandos CLI e a função principal de criação de USB.

#### Funções Principais:

**`NewUSBCommand()` → `*cobra.Command`**
- **Propósito**: Cria o comando principal USB e seus subcomandos
- **Retorna**: Comando Cobra configurado com subcomandos
- **Subcomandos**: `list`, `create`, `format`

**`newUSBListCommand()` → `*cobra.Command`**
- **Propósito**: Cria comando para listar dispositivos USB
- **Flags**: `--format` (table/json/yaml)
- **Funcionalidade**: Detecta e lista dispositivos USB disponíveis

**`newUSBCreateCommand()` → `*cobra.Command`**
- **Propósito**: Cria comando para criar USB bootável
- **Flags**:
  - `--node-name` (obrigatório): Nome do nó Syntropy
  - `--description`: Descrição do nó
  - `--coordinates`: Coordenadas geográficas
  - `--owner-key`: Arquivo de chave do proprietário
  - `--auto-detect`: Auto-detectar dispositivo USB
  - `--label`: Rótulo do sistema de arquivos
  - `--work-dir`: Diretório de trabalho (padrão: ~/.syntropy/work)
  - `--cache-dir`: Diretório de cache (padrão: ~/.syntropy/cache)
  - `--iso`: Caminho da ISO Ubuntu
  - `--discovery-server`: Servidor de descoberta
  - `--created-by`: Usuário criador

**`newUSBFormatCommand()` → `*cobra.Command`**
- **Propósito**: Cria comando para formatar USB
- **Flags**:
  - `--label`: Rótulo do sistema de arquivos
  - `--force`: Não pedir confirmação
- **Funcionalidade**: Formata dispositivo USB com FAT32

**`createUSB(devicePath, config, workDir, cacheDir)` → `error`**
- **Propósito**: Função principal de criação de USB bootável
- **Parâmetros**:
  - `devicePath`: Caminho do dispositivo USB
  - `config`: Configuração do nó Syntropy
  - `workDir`: Diretório de trabalho
  - `cacheDir`: Diretório de cache
- **Fluxo**:
  1. Valida dispositivo
  2. Configura diretórios
  3. Gera chaves SSH
  4. Gera certificados TLS
  5. Cria configuração cloud-init
  6. Copia scripts
  7. Cria USB com estratégia NoCloud

---

### 🌐 `platform.go` (248 linhas)
**Detecção de plataforma e funções comuns**

Gerencia a detecção de plataforma e fornece funções comuns para todas as plataformas.

#### Funções Principais:

**`detectPlatform()` → `string`**
- **Propósito**: Detecta a plataforma atual (linux/windows/wsl)
- **Lógica**:
  - Verifica `runtime.GOOS`
  - Detecta WSL via `/proc/sys/fs/binfmt_misc/WSLInterop`
  - Verifica `/proc/version` para "microsoft"
- **Retorna**: "linux", "windows" ou "wsl"

**`ListDevices()` → `[]USBDevice, error`**
- **Propósito**: Lista dispositivos USB baseado na plataforma
- **Delegação**: Chama função específica da plataforma
- **Retorna**: Slice de dispositivos USB detectados

**`validateDevice(devicePath)` → `error`**
- **Propósito**: Valida se dispositivo é seguro para usar
- **Validações**:
  - Verifica se dispositivo existe
  - Verifica se não é dispositivo do sistema
  - Implementação específica por plataforma

**`validateDeviceLinux(devicePath)` → `error`**
- **Propósito**: Validação específica para Linux
- **Verificações**:
  - Existência do dispositivo via `os.Stat`
  - Verificação de partições do sistema via `lsblk`

**`validateDeviceWindows(devicePath)` → `error`**
- **Propósito**: Validação específica para Windows/WSL
- **Verificações**:
  - Parse do número do disco
  - Verificação via PowerShell se é dispositivo do sistema

**`isSystemDeviceLinux(devicePath)` → `bool`**
- **Propósito**: Verifica se dispositivo é do sistema no Linux
- **Lógica**: Verifica mountpoints como "/", "/boot", "/home"

**`isSystemDeviceWindows(diskNum)` → `bool`**
- **Propósito**: Verifica se dispositivo é do sistema no Windows
- **Método**: Executa script PowerShell para verificar propriedades do disco

**`formatUSB(devicePath, label, force)` → `error`**
- **Propósito**: Formata dispositivo USB
- **Parâmetros**:
  - `devicePath`: Caminho do dispositivo
  - `label`: Rótulo do sistema de arquivos
  - `force`: Pular confirmação
- **Fluxo**:
  1. Valida dispositivo
  2. Solicita confirmação (se não `force`)
  3. Delega para implementação específica da plataforma

**`listUSBDevices(format)` → `error`**
- **Propósito**: Lista dispositivos e formata saída
- **Formatos**: table, json, yaml
- **Funcionalidades**:
  - Detecta plataforma
  - Mostra avisos específicos (ex: WSL)
  - Formata saída conforme solicitado

**`SelectDevice()` → `*USBDevice, error`**
- **Propósito**: Seleção interativa de dispositivo USB
- **Comportamento**:
  - Auto-seleciona se apenas 1 dispositivo
  - Mostra menu se múltiplos dispositivos
  - Valida seleção do usuário

---

### 🐧 `linux.go` (381 linhas)
**Implementações específicas do Linux**

Contém todas as implementações específicas para sistemas Linux.

#### Funções Principais:

**`listDevicesLinux()` → `[]USBDevice, error`**
- **Propósito**: Lista dispositivos USB no Linux
- **Método Principal**: Usa `lsblk -J` para detecção robusta
- **Fallback**: `listDevicesLinuxFallback()` se lsblk falhar
- **Filtros**:
  - Apenas dispositivos de disco (não partições)
  - Dispositivos removíveis ou com transport USB
  - Heurística para dispositivos pequenos (< 2TB)

**`listDevicesLinuxFallback()` → `[]USBDevice, error`**
- **Propósito**: Método alternativo de detecção
- **Método**: Varre diretórios `/sys/block/`
- **Padrões**: `sd*`, `nvme*`, `mmcblk*`
- **Verificações**:
  - Arquivo `removable`
  - Tamanho via `size` (setores × 512)
  - Modelo e vendor via `/device/`

**`parseSizeToGB(sizeStr)` → `int`**
- **Propósito**: Converte string de tamanho para GB
- **Suporta**: "1.5G", "1024M", "2T"
- **Retorna**: Tamanho em GB como inteiro

**`createUSBWithNoCloudLinux(devicePath, config, workDir, cacheDir)` → `error`**
- **Propósito**: Cria USB bootável no Linux usando estratégia NoCloud
- **Requisitos**: Privilégios de root
- **Fluxo**:
  1. Gerencia cache de ISO
  2. Grava ISO via `dd` com timeout
  3. Cria partição CIDATA via `sgdisk`
  4. Formata partição CIDATA com `mkfs.vfat`
  5. Monta partição e copia arquivos cloud-init
  6. Desmonta e sincroniza

**`createUSBLinux(devicePath, config, workDir, cacheDir)` → `error`**
- **Propósito**: Função legada que redireciona para nova implementação

**`formatUSBLinux(devicePath, label)` → `error`**
- **Propósito**: Formata dispositivo USB no Linux
- **Requisitos**: Privilégios de root
- **Fluxo**:
  1. Desmonta todas as partições existentes
  2. Cria tabela de partições GPT (fallback para msdos)
  3. Cria partição primária com alinhamento 1MiB
  4. Formata com FAT32

**`runCommandWithTimeout(timeout, name, args...)` → `error`**
- **Propósito**: Executa comando com timeout
- **Parâmetros**:
  - `timeout`: Duração máxima
  - `name`: Nome do comando
  - `args`: Argumentos do comando
- **Funcionalidades**:
  - Contexto com timeout
  - Captura de stdout/stderr
  - Tratamento de timeout vs erro

**`umountAll(device)`**
- **Propósito**: Desmonta todas as partições de um dispositivo
- **Método**: Usa `mount | awk` para encontrar partições montadas
- **Execução**: Executa `umount` para cada partição encontrada

---

### 🪟 `windows.go` (521 linhas)
**Implementações específicas do Windows/WSL**

Contém implementações para Windows nativo e WSL (Windows Subsystem for Linux).

#### Funções WSL:

**`listDevicesWSL()` → `[]USBDevice, error`**
- **Propósito**: Lista dispositivos USB no WSL
- **Método**: Executa PowerShell via `powershell.exe`
- **Script**: `Get-Disk` com filtros para USB e SCSI pequenos
- **Parse**: JSON do PowerShell convertido para `WindowsDisk`
- **Fallback**: `listDevicesWSLAlternative()` se PowerShell falhar

**`listDevicesWSLAlternative()` → `[]USBDevice, error`**
- **Propósito**: Método alternativo usando WMIC
- **Comando**: `wmic diskdrive where "InterfaceType='USB'"`
- **Parse**: CSV output do WMIC

**`createUSBWithNoCloudWSL(devicePath, config, workDir, cacheDir)` → `error`**
- **Propósito**: Cria USB bootável no WSL usando estratégia NoCloud
- **Fluxo**:
  1. Extrai número do disco Windows
  2. Gera script PowerShell elevado
  3. Coloca disco offline no Windows
  4. Monta disco cru no WSL (`wsl --mount --bare`)
  5. Detecta device no WSL e grava ISO via `dd`
  6. Cria partição CIDATA via `sgdisk`
  7. Formata e monta partição CIDATA
  8. Copia arquivos cloud-init
  9. Desmonta e volta disco online

#### Funções Windows:

**`listDevicesWindows()` → `[]USBDevice, error`**
- **Propósito**: Lista dispositivos USB no Windows nativo
- **Método**: Similar ao WSL mas sem conversões de caminho

**`createUSBWithNoCloudWindows(devicePath, config, workDir, cacheDir)` → `error`**
- **Propósito**: Cria USB bootável no Windows nativo
- **Fluxo**: Similar ao WSL mas executa diretamente no Windows

**`formatUSBWSL(devicePath, label)` → `error`**
- **Propósito**: Formata dispositivo USB no WSL
- **Método**: PowerShell com `Clear-Disk`, `New-Partition`, `Format-Volume`

**`formatUSBWindows(devicePath, label)` → `error`**
- **Propósito**: Formata dispositivo USB no Windows nativo
- **Status**: Implementação pendente

#### Funções de Conversão:

**`convertAnyToWSLPath(p)` → `string`**
- **Propósito**: Converte qualquer caminho para formato WSL
- **Suporta**: Caminhos Windows (`C:\...`) e WSL (`/mnt/c/...`)
- **Método**: Tenta `wslpath -u`, fallback manual

**`convertWSLToWindowsPath(wslPath)` → `string`**
- **Propósito**: Converte caminho WSL para Windows
- **Método**: Usa `wslpath -w`, fallback manual

---

### 🔐 `certificates.go` (169 linhas)
**Geração de certificados TLS e chaves SSH**

Responsável pela geração e gerenciamento de certificados e chaves.

#### Funções Principais:

**`generateSSHKeyPair(nodeName)` → `(string, string, error)`**
- **Propósito**: Gera par de chaves SSH usando KeyManager centralizado
- **Algoritmo**: ED25519 (padrão), RSA 2048 bits (fallback)
- **Armazenamento**: `~/.syntropy/keys/{nodeName}-node.key`
- **Formato**: OpenSSH (chave privada PEM, chave pública autorizada)
- **Persistência**: Chaves salvas automaticamente no diretório centralizado
- **Retorna**: `(chavePrivada, chavePublica, erro)`

**`loadExistingSSHKeyPair(nodeName)` → `(string, string, error)`**
- **Propósito**: Carrega par de chaves SSH existente
- **Localização**: `~/.syntropy/keys/{nodeName}-node.key`
- **Funcionalidade**: Reutiliza chaves existentes para o mesmo nó
- **Retorna**: `(chavePrivada, chavePublica, erro)`

**`generateCertificates(nodeName, ownerKey)` → `(*Certificates, error)`**
- **Propósito**: Gera certificados TLS para o nó
- **Certificados**:
  - **CA**: RSA 4096 bits, válido por 10 anos
  - **Nó**: RSA 2048 bits, válido por 1 ano
- **Campos**:
  - Organização: "Syntropy Cooperative Grid"
  - País: "BR"
  - Localidade: "São Paulo"
  - CN: nome do nó
  - IPs: 127.0.0.1
  - DNS: nome do nó, localhost

**`saveCertificates(certs, workDir)` → `(string, string, string, string, error)`**
- **Propósito**: Salva certificados no diretório de trabalho
- **Estrutura**:
  ```
  ~/.syntropy/work/usb-{timestamp}/certs/
  ├── ca.key    (0600)
  ├── ca.crt    (0644)
  ├── node.key  (0600)
  └── node.crt  (0644)
  ```
- **Retorna**: Caminhos dos arquivos salvos

---

### ☁️ `cloudinit.go` (567 linhas)
**Configuração do cloud-init**

Gera e gerencia arquivos de configuração do cloud-init para inicialização automática dos nós.

#### Funções Principais:

**`generateCloudInitConfig(config, workDir, certs)` → `(*CloudInitConfig, error)`**
- **Propósito**: Gera configuração completa do cloud-init
- **Parâmetros**:
  - `config`: Configuração do nó
  - `workDir`: Diretório de trabalho
  - `certs`: Certificados TLS
- **Gera**:
  - ID único da instância
  - Configurações de rede padrão
  - Sufixo IP único
  - Caminhos de certificados

**`renderTemplate(templateStr, data)` → `(string, error)`**
- **Propósito**: Renderiza template Go com dados
- **Configuração**: `missingkey=error` para detectar variáveis não definidas

**`createCloudInitFiles(config, workDir, certPaths)` → `error`**
- **Propósito**: Cria arquivos de configuração do cloud-init
- **Arquivos Gerados**:
  - `user-data`: Configuração principal do cloud-init
  - `meta-data`: Metadados do nó
  - `network-config`: Configuração de rede

#### Template user-data:
- **Localização**: pt_BR.UTF-8, America/Sao_Paulo
- **Usuário**: `syntropy` com sudo sem senha
- **SSH**: Apenas chave pública, sem senha
- **Pacotes**: Docker, Kubernetes, WireGuard, monitoramento
- **Firewall**: UFW configurado com portas Syntropy
- **Syntropy Agent**: Download, configuração e service systemd
- **Logs**: Logrotate configurado

#### Template meta-data:
- **Instance ID**: Único baseado em timestamp
- **Metadados Syntropy**: Configuração completa do nó
- **Auditoria**: Configurações de log e retenção

#### Template network-config:
- **DHCP**: Configuração automática para en*, eth*, enp*
- **DNS**: Google DNS e Cloudflare
- **Bridge**: br0 para virtualização (172.20.0.x)
- **VLAN**: vlan100 para gerenciamento (192.168.100.x)
- **Rotas**: Redes Syntropy com tabelas de roteamento

**`copyScripts(workDir)` → `error`**
- **Propósito**: Copia scripts de instalação para diretório de trabalho
- **Scripts**:
  - `hardware-detection.sh`
  - `network-discovery.sh`
  - `syntropy-install.sh`
  - `cluster-join.sh`
- **Origem**: `infrastructure/cloud-init/scripts/`

---

### 🛠️ `utils.go` (246 linhas)
**Funções auxiliares e utilitários**

Contém funções utilitárias para formatação, cache e download.

#### Funções de Formatação:

**`formatSize(bytes)` → `string`**
- **Propósito**: Formata bytes em string legível
- **Unidades**: B, KB, MB, GB, TB
- **Exemplo**: `1073741824` → `"1.0 GB"`

**`outputTable(devices)` → `error`**
- **Propósito**: Exibe dispositivos em formato tabela
- **Adaptativo**: Layout diferente para Windows/WSL vs Linux
- **Colunas**: Dispositivo, Tamanho, Modelo, Serial, Plataforma

**`outputJSON(devices)` → `error`**
- **Propósito**: Exibe dispositivos em formato JSON
- **Configuração**: Indentação de 2 espaços

**`outputYAML(devices)` → `error`**
- **Propósito**: Exibe dispositivos em formato YAML
- **Formato**: Lista de dispositivos com todos os campos

#### Funções de Cache e Download:

**`downloadUbuntuISO(destPath)` → `error`**
- **Propósito**: Baixa ISO Ubuntu 24.04 LTS Server
- **URL**: `https://releases.ubuntu.com/24.04/ubuntu-24.04-live-server-amd64.iso`
- **Ferramenta**: `wget` com progresso e retry

**`manageISOCache(cacheDir)` → `(string, error)`**
- **Propósito**: Gerencia cache inteligente de ISOs Ubuntu
- **Estratégia**:
  1. Verifica ISOs existentes no cache
  2. Valida tamanho (> 500MB)
  3. Tenta baixar versões mais recentes primeiro
  4. Fallback para versões mais antigas
- **Versões Suportadas**:
  - Ubuntu 24.04 LTS Server
  - Ubuntu 22.04.5 LTS Server
  - Ubuntu 22.04.4 LTS Server
  - Ubuntu 20.04.6 LTS Server

**`checkURLExists(url)` → `bool`**
- **Propósito**: Verifica se URL existe via HEAD request
- **Método**: `curl -I --head --fail`

**`downloadWithProgress(url, destPath)` → `error`**
- **Propósito**: Baixa arquivo com indicador de progresso
- **Características**:
  - Arquivo temporário durante download
  - Renomeação atômica no final
  - Retry automático (3 tentativas)
  - Timeout de 30 segundos

---

## 🔄 Fluxo de Execução

### 1. **Listagem de Dispositivos**
```
NewUSBCommand() → newUSBListCommand() → listUSBDevices() → ListDevices() → 
[listDevicesLinux() | listDevicesWSL() | listDevicesWindows()] → 
outputTable()/outputJSON()/outputYAML()
```

### 2. **Criação de USB**
```
NewUSBCommand() → newUSBCreateCommand() → createUSB() → 
[validateDevice()] → [generateSSHKeyPair()] → [generateCertificates()] → 
[saveCertificates()] → [generateCloudInitConfig()] → [createCloudInitFiles()] → 
[copyScripts()] → [createUSBWithNoCloudLinux() | createUSBWithNoCloudWSL() | createUSBWithNoCloudWindows()]
```

### 3. **Formatação de USB**
```
NewUSBCommand() → newUSBFormatCommand() → formatUSB() → 
[validateDevice()] → [formatUSBLinux() | formatUSBWSL() | formatUSBWindows()]
```

---

## 🔒 Considerações de Segurança

### **Validação de Dispositivos**
- Verificação de existência do dispositivo
- Detecção de dispositivos do sistema
- Confirmação do usuário para operações destrutivas

### **Certificados TLS**
- Geração de CA própria (4096 bits)
- Certificados de nó com 2048 bits
- Validação de nomes e IPs

### **Chaves SSH - Sistema Centralizado** ✅
- **Algoritmo Principal**: ED25519 (mais seguro e rápido)
- **Fallback**: RSA 2048 bits (compatibilidade)
- **Armazenamento**: `~/.syntropy/keys/` (centralizado e persistente)
- **Nomenclatura**: `{nodeName}-node.key` (padronizada)
- **Fingerprints**: Gerados automaticamente para auditoria
- **Reutilização**: Chaves existentes são reutilizadas para o mesmo nó
- **Segurança**: Chave privada nunca enviada para o nó

### **Gerenciamento de Chaves**
- **KeyManager Centralizado**: Usa `infrastructure/KeyManager`
- **Persistência**: Chaves salvas automaticamente
- **Auditoria**: Fingerprints para rastreamento
- **Organização**: Estrutura padronizada por propósito

### **Permissões**
- Arquivos de chave privada: 0600
- Arquivos de certificado: 0644
- Scripts executáveis: 0755
- Diretório de chaves: 0700

---

## 🧪 Testabilidade

### **Pontos de Teste**
- Detecção de plataforma
- Validação de dispositivos
- Geração de certificados
- Renderização de templates
- Cache de ISOs

### **Mocking**
- Funções de sistema (`os.Stat`, `exec.Command`)
- Download de ISOs
- Validação de dispositivos

### **Cobertura**
- Testes unitários por arquivo
- Testes de integração por plataforma
- Testes de regressão para templates

---

## 📈 Performance

### **Otimizações**
- Cache de ISOs para evitar downloads repetidos
- Timeout em comandos longos (30 min para dd)
- Detecção incremental de dispositivos
- Fallbacks para métodos alternativos

### **Recursos**
- Uso eficiente de memória para certificados
- Stream de dados para downloads grandes
- Cleanup automático de arquivos temporários

---

## 🚀 Extensibilidade

### **Novas Plataformas**
- Implementar funções específicas em novo arquivo
- Adicionar detecção em `detectPlatform()`
- Registrar em `ListDevices()` e `formatUSB()`

### **Novos Formatos**
- Adicionar função `outputXXX()` em `utils.go`
- Registrar em `listUSBDevices()`

### **Novas Funcionalidades**
- Extender `Config` e `CloudInitConfig`
- Adicionar comandos em `commands.go`
- Implementar lógica específica por plataforma

---

## 📚 Referências

- [Cloud-init Documentation](https://cloudinit.readthedocs.io/)
- [Ubuntu Server Installation Guide](https://ubuntu.com/server/docs/installation)
- [Go Cobra CLI Framework](https://cobra.dev/)
- [PowerShell Disk Management](https://docs.microsoft.com/en-us/powershell/module/storage/)
- [Linux Block Devices](https://www.kernel.org/doc/Documentation/block/)
