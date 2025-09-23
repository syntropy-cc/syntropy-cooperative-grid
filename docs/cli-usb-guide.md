# Syntropy CLI - Guia do Comando USB

## Visão Geral

O comando `syntropy usb` é a interface principal para criar USBs bootáveis personalizados para nós da Syntropy Cooperative Grid. Este guia detalha todas as funcionalidades, opções e casos de uso do comando.

## Índice

1. [Comandos Disponíveis](#comandos-disponíveis)
2. [Comando Create](#comando-create)
3. [Comando List](#comando-list)
4. [Comando Format](#comando-format)
5. [Exemplos Práticos](#exemplos-práticos)
6. [Troubleshooting](#troubleshooting)

---

## Comandos Disponíveis

### Estrutura Geral
```bash
syntropy usb <comando> [argumentos] [flags]
```

### Comandos Principais
- `create` - Cria USB com boot para um nó Syntropy
- `list` - Lista dispositivos USB disponíveis
- `format` - Formata um dispositivo USB

---

## Comando Create

### Sintaxe
```bash
syntropy usb create [device] [flags]
```

### Descrição
Cria um USB bootável contendo Ubuntu Server e configuração automática para um nó da Syntropy Cooperative Grid. O sistema gera automaticamente:

- Certificados TLS únicos
- Chaves SSH
- Configuração cloud-init personalizada
- Scripts de instalação
- ISO Ubuntu personalizada

### Argumentos
- `device` (opcional): Caminho para o dispositivo USB (ex: `/dev/sdb`, `PHYSICALDRIVE1`)

### Flags Obrigatórias
- `--node-name string`: Nome único do nó (obrigatório)

### Flags Opcionais
- `--description string`: Descrição do nó
- `--coordinates string`: Coordenadas geográficas (formato: lat,lon)
- `--owner-key string`: Arquivo de chave de proprietário existente
- `--auto-detect`: Detectar automaticamente dispositivo USB
- `--label string`: Rótulo do sistema de arquivos (padrão: "SYNTROPY")
- `--work-dir string`: Diretório de trabalho (padrão: `/tmp/syntropy-work`)
- `--cache-dir string`: Diretório de cache (padrão: `~/.syntropy/cache`)
- `--iso string`: Caminho para ISO Ubuntu (baixa automaticamente se não especificado)
- `--discovery-server string`: Servidor de descoberta da rede (padrão: "syntropy-discovery.local")
- `--created-by string`: Usuário que criou o nó (padrão: usuário atual)

### Exemplos

#### Criação Básica
```bash
# Criar USB com auto-detecção de dispositivo
syntropy usb create --auto-detect --node-name "server-01"
```

#### Criação com Descrição
```bash
# Criar USB com descrição detalhada
syntropy usb create --auto-detect --node-name "home-server" --description "Servidor doméstico principal"
```

#### Criação com Coordenadas
```bash
# Criar USB com localização geográfica
syntropy usb create --auto-detect --node-name "node-sp" --coordinates "-23.5505,-46.6333"
```

#### Criação Especificando Dispositivo
```bash
# Linux
syntropy usb create /dev/sdb --node-name "node-01"

# Windows/WSL
syntropy usb create PHYSICALDRIVE1 --node-name "node-01"
```

#### Criação com Servidor de Descoberta Personalizado
```bash
# Usar servidor de descoberta específico
syntropy usb create --auto-detect --node-name "node-02" --discovery-server "192.168.1.100"
```

#### Criação com ISO Personalizada
```bash
# Usar ISO Ubuntu específica
syntropy usb create --auto-detect --node-name "node-03" --iso "/path/to/ubuntu-24.04.iso"
```

#### Criação com Diretórios Personalizados
```bash
# Especificar diretórios de trabalho e cache
syntropy usb create --auto-detect --node-name "node-04" \
    --work-dir "/tmp/custom-work" \
    --cache-dir "/tmp/custom-cache"
```

### Saída Esperada
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

🔧 Gravando ISO no dispositivo USB...
✅ USB criado com sucesso!
```

### Estrutura de Arquivos Gerados

#### Diretório de Trabalho
```
/tmp/syntropy-work/
├── cloud-init/
│   ├── user-data
│   ├── meta-data
│   └── network-config
├── scripts/
│   ├── hardware-detection.sh
│   ├── network-discovery.sh
│   ├── syntropy-install.sh
│   └── cluster-join.sh
├── certs/
│   ├── ca.crt
│   ├── ca.key
│   ├── node.crt
│   └── node.key
├── ubuntu-base.iso
└── syntropy-node-01.iso
```

#### Estrutura do USB
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

---

## Comando List

### Sintaxe
```bash
syntropy usb list [flags]
```

### Descrição
Lista todos os dispositivos USB disponíveis no sistema, mostrando informações detalhadas sobre cada dispositivo.

### Flags
- `--format string`: Formato de saída (table, json, yaml) (padrão: "table")

### Exemplos

#### Listagem em Formato Tabela
```bash
syntropy usb list
```

**Saída**:
```
🔍 Dispositivos USB detectados:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
DISPOSITIVO    TAMANHO    MODELO                    FABRICANTE       REMOVÍVEL    PLATAFORMA
/dev/sdb       32.0 GB    SanDisk Ultra             SanDisk          Sim          linux
/dev/sdc       16.0 GB    Kingston DataTraveler     Kingston         Sim          linux
```

#### Listagem em Formato JSON
```bash
syntropy usb list --format json
```

**Saída**:
```json
[
  {
    "path": "/dev/sdb",
    "size": "32.0 GB",
    "size_gb": 32,
    "model": "SanDisk Ultra",
    "vendor": "SanDisk",
    "serial": "1234567890",
    "removable": true,
    "platform": "linux"
  },
  {
    "path": "/dev/sdc",
    "size": "16.0 GB",
    "size_gb": 16,
    "model": "Kingston DataTraveler",
    "vendor": "Kingston",
    "serial": "0987654321",
    "removable": true,
    "platform": "linux"
  }
]
```

#### Listagem em Formato YAML
```bash
syntropy usb list --format yaml
```

**Saída**:
```yaml
devices:
  - path: /dev/sdb
    size: 32.0 GB
    size_gb: 32
    model: SanDisk Ultra
    vendor: SanDisk
    serial: 1234567890
    removable: true
    platform: linux
  - path: /dev/sdc
    size: 16.0 GB
    size_gb: 16
    model: Kingston DataTraveler
    vendor: Kingston
    serial: 0987654321
    removable: true
    platform: linux
```

### Detecção de Dispositivos

#### Linux
O comando detecta dispositivos USB através do sistema de arquivos:
- Verifica `/sys/block/sd*` para dispositivos de bloco
- Confirma que são removíveis através de `/sys/block/*/removable`
- Obtém informações de modelo e fabricante

#### Windows/WSL
O comando usa PowerShell para detectar dispositivos:
- Executa `Get-Disk` para listar discos
- Filtra por tipo USB ou SCSI com tamanho apropriado
- Obtém informações detalhadas de cada dispositivo

---

## Comando Format

### Sintaxe
```bash
syntropy usb format <device> [flags]
```

### Descrição
Formata um dispositivo USB com sistema de arquivos FAT32. **⚠️ ATENÇÃO: Esta operação apagará TODOS os dados do dispositivo!**

### Argumentos
- `device` (obrigatório): Caminho para o dispositivo USB

### Flags
- `--label string`: Rótulo do sistema de arquivos (padrão: "SYNTROPY")
- `--force`: Não pedir confirmação

### Exemplos

#### Formatação Básica
```bash
# Linux
syntropy usb format /dev/sdb

# Windows/WSL
syntropy usb format PHYSICALDRIVE1
```

#### Formatação com Rótulo Personalizado
```bash
syntropy usb format /dev/sdb --label "MYUSB"
```

#### Formatação sem Confirmação
```bash
syntropy usb format /dev/sdb --force
```

### Saída Esperada
```
⚠️  ATENÇÃO: Esta operação apagará TODOS os dados em /dev/sdb!
Tem certeza que deseja continuar? (y/N): y

🔧 Formatando dispositivo /dev/sdb...
✅ Dispositivo /dev/sdb formatado com sucesso!
```

### Processo de Formatação

#### Linux
1. Desmonta partições existentes
2. Cria nova tabela de partições (MSDOS)
3. Cria partição primária FAT32
4. Formata partição com mkfs.vfat

#### Windows/WSL
1. Limpa disco com Clear-Disk
2. Cria nova partição com New-Partition
3. Formata volume com Format-Volume

---

## Exemplos Práticos

### Cenário 1: Primeiro Nó da Rede

#### Objetivo
Criar o primeiro nó que será o líder da rede Syntropy.

#### Comando
```bash
syntropy usb create --auto-detect \
    --node-name "leader-01" \
    --description "Primeiro nó da rede Syntropy" \
    --coordinates "-23.5505,-46.6333" \
    --discovery-server "syntropy-discovery.local"
```

#### Resultado
- Nó será configurado como líder
- Iniciará serviços de descoberta
- Outros nós se conectarão a este nó

### Cenário 2: Nó Worker em Casa

#### Objetivo
Criar um nó worker para uso doméstico.

#### Comando
```bash
syntropy usb create --auto-detect \
    --node-name "home-worker-01" \
    --description "Nó worker doméstico" \
    --coordinates "-23.5505,-46.6333" \
    --discovery-server "192.168.1.100"
```

#### Resultado
- Nó será configurado como worker
- Conectará ao líder em 192.168.1.100
- Contribuirá com recursos da rede

### Cenário 3: Nó com ISO Personalizada

#### Objetivo
Criar um nó usando uma ISO Ubuntu específica.

#### Comando
```bash
syntropy usb create --auto-detect \
    --node-name "custom-node-01" \
    --description "Nó com ISO personalizada" \
    --iso "/path/to/ubuntu-24.04-custom.iso"
```

#### Resultado
- Usará ISO personalizada em vez de baixar
- Configuração será aplicada sobre a ISO existente
- Útil para ambientes com restrições de rede

### Cenário 4: Múltiplos Nós

#### Objetivo
Criar vários nós para um cluster.

#### Script
```bash
#!/bin/bash

# Criar líder
syntropy usb create --auto-detect \
    --node-name "cluster-leader" \
    --description "Líder do cluster" \
    --coordinates "-23.5505,-46.6333"

# Criar workers
for i in {1..3}; do
    syntropy usb create --auto-detect \
        --node-name "cluster-worker-$i" \
        --description "Worker $i do cluster" \
        --coordinates "-23.5505,-46.6333" \
        --discovery-server "syntropy-discovery.local"
done
```

#### Resultado
- Um líder e três workers
- Todos configurados para a mesma localização
- Workers se conectarão ao líder

---

## Troubleshooting

### Problemas Comuns

#### 1. Dispositivo USB Não Detectado

**Sintomas**:
```
❌ Nenhum dispositivo USB encontrado.
```

**Soluções**:
1. Verificar se USB está conectado
2. Verificar se USB é reconhecido pelo sistema
3. Tentar outro USB
4. Verificar permissões (Linux: usar sudo)

**Comandos de Diagnóstico**:
```bash
# Linux
lsusb
fdisk -l
lsblk

# Windows
Get-Disk
Get-Volume
```

#### 2. Falha na Geração de Certificados

**Sintomas**:
```
🔐 Gerando certificados TLS...
❌ Erro ao gerar certificados: crypto/rsa: too few primes of given length to generate an RSA key
```

**Soluções**:
1. Verificar se sistema tem entropia suficiente
2. Instalar rng-tools (Linux)
3. Reiniciar serviço de entropia

**Comandos de Correção**:
```bash
# Linux
sudo apt install rng-tools
sudo systemctl start rng-tools
sudo systemctl enable rng-tools

# Verificar entropia
cat /proc/sys/kernel/random/entropy_avail
```

#### 3. Falha no Download da ISO

**Sintomas**:
```
📥 Baixando ISO Ubuntu base...
❌ Erro ao baixar ISO Ubuntu: wget: unable to resolve host address 'releases.ubuntu.com'
```

**Soluções**:
1. Verificar conectividade com internet
2. Verificar configuração DNS
3. Usar ISO local com --iso

**Comandos de Diagnóstico**:
```bash
# Verificar conectividade
ping 8.8.8.8
ping releases.ubuntu.com

# Verificar DNS
nslookup releases.ubuntu.com
dig releases.ubuntu.com

# Usar ISO local
syntropy usb create --auto-detect --node-name "node-01" --iso "/path/to/ubuntu.iso"
```

#### 4. Falha na Gravação do USB

**Sintomas**:
```
🔧 Gravando ISO no dispositivo USB...
❌ Erro ao gravar ISO: dd: /dev/sdb: Permission denied
```

**Soluções**:
1. Usar sudo (Linux)
2. Executar como administrador (Windows)
3. Verificar se dispositivo não está montado

**Comandos de Correção**:
```bash
# Linux
sudo umount /dev/sdb*
sudo syntropy usb create /dev/sdb --node-name "node-01"

# Windows
# Executar PowerShell como Administrador
```

#### 5. Falha na Criação da ISO

**Sintomas**:
```
💿 Criando ISO personalizada...
❌ Erro ao criar ISO personalizada: genisoimage: command not found
```

**Soluções**:
1. Instalar genisoimage (Linux)
2. Instalar cdrtools (Linux)
3. Usar alternativa como mkisofs

**Comandos de Correção**:
```bash
# Ubuntu/Debian
sudo apt install genisoimage

# CentOS/RHEL
sudo yum install genisoimage

# Alternativa
sudo apt install cdrtools
```

### Logs e Diagnóstico

#### Localização dos Logs
```bash
# Logs da CLI
~/.syntropy/logs/

# Logs do sistema
/var/log/syslog
/var/log/kern.log

# Logs do cloud-init
/var/log/cloud-init.log
```

#### Comandos de Diagnóstico
```bash
# Verificar status do USB
lsblk
fdisk -l /dev/sdb

# Verificar se ISO foi gravada
file /dev/sdb
hexdump -C /dev/sdb | head

# Verificar conectividade
ping 8.8.8.8
curl -I https://releases.ubuntu.com/

# Verificar recursos do sistema
free -h
df -h
```

### Recuperação

#### Limpar Diretório de Trabalho
```bash
# Limpar arquivos temporários
rm -rf /tmp/syntropy-work
rm -rf ~/.syntropy/cache/iso
```

#### Recriar USB
```bash
# Formatar USB
syntropy usb format /dev/sdb --force

# Recriar USB
syntropy usb create /dev/sdb --node-name "node-01"
```

#### Verificar Integridade
```bash
# Verificar se USB é bootável
file /dev/sdb

# Verificar se ISO está correta
file /tmp/syntropy-work/syntropy-node-01.iso

# Verificar se arquivos estão presentes
ls -la /tmp/syntropy-work/
```

---

## Conclusão

O comando `syntropy usb` oferece uma interface completa para criar USBs bootáveis personalizados para a rede Syntropy Cooperative Grid. Com recursos como:

- Detecção automática de dispositivos
- Geração automática de certificados e chaves
- Configuração personalizada por nó
- Suporte a múltiplas plataformas
- Sistema robusto de diagnóstico

Este comando é essencial para o funcionamento do MVP da rede Syntropy, permitindo que usuários criem nós de forma simples e segura.
