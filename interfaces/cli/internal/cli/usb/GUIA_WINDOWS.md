# 🪟 Guia Completo - USB Syntropy para Windows

## 📋 Visão Geral

Este guia fornece instruções detalhadas para usar o módulo USB Syntropy especificamente no Windows, incluindo todas as verificações de permissões, tratamento de erros e soluções para problemas comuns.

## 🚀 Início Rápido

### 1. **Verificar Pré-requisitos**
```powershell
# Executar como Administrador
syntropy usb-win debug
```

### 2. **Listar Dispositivos USB**
```powershell
syntropy usb-win list
```

### 3. **Criar USB Bootável**
```powershell
syntropy usb-win create --node-name "meu-no"
```

---

## 🔧 Pré-requisitos

### **Sistema Operacional**
- ✅ Windows 10 (versão 1903 ou superior)
- ✅ Windows 11 (todas as versões)
- ✅ Arquitetura: x64 (64-bit)

### **Privilégios Necessários**
- ✅ **Administrador**: Obrigatório para acesso a dispositivos
- ✅ **UAC**: Deve estar habilitado mas não bloqueando
- ✅ **PowerShell**: Versão 5.1 ou superior

### **Software Necessário**
- ✅ **WSL 2**: Windows Subsystem for Linux
- ✅ **Ubuntu WSL**: Distribuição Linux no WSL
- ✅ **PowerShell**: Com política de execução adequada

---

## 📥 Instalação e Configuração

### **1. Instalar WSL (se não estiver instalado)**
```powershell
# Executar no PowerShell como Administrador
wsl --install

# Ou instalar Ubuntu especificamente
wsl --install -d Ubuntu

# Reiniciar o computador se solicitado
```

### **2. Configurar Política de Execução**
```powershell
# Verificar política atual
Get-ExecutionPolicy

# Configurar política (se necessário)
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser

# Ou para apenas o processo atual
Set-ExecutionPolicy -ExecutionPolicy Bypass -Scope Process
```

### **3. Verificar WSL**
```powershell
# Verificar versão do WSL
wsl --version

# Listar distribuições instaladas
wsl --list --verbose

# Testar execução
wsl echo "WSL funcionando!"
```

---

## 🔍 Comandos Disponíveis

### **`syntropy usb-win list`**
Lista dispositivos USB com informações detalhadas.

```powershell
# Listar em formato tabela (padrão)
syntropy usb-win list

# Listar em formato JSON
syntropy usb-win list --format json

# Listar em formato YAML
syntropy usb-win list --format yaml
```

**Saída esperada:**
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
💾 Dispositivos USB Detectados (Windows Only)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
[1] PHYSICALDRIVE1
     Nome: USB Drive
     Tamanho: 8.0 GB
     Modelo: SanDisk USB 3.0
     Serial: 1234567890
     Status: Healthy
     Sistema: False | Boot: False | Offline: False
     Partições: 1
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### **`syntropy usb-win create`**
Cria USB bootável para nó Syntropy.

```powershell
# Criação básica com auto-detecção
syntropy usb-win create --node-name "node-01"

# Criação com dispositivo específico
syntropy usb-win create PHYSICALDRIVE1 --node-name "node-01"

# Criação com configurações completas
syntropy usb-win create PHYSICALDRIVE1 \
  --node-name "node-01" \
  --description "Nó principal da rede" \
  --coordinates "-23.5505,-46.6333" \
  --label "SYNTROPY-NODE" \
  --discovery-server "discovery.syntropy.local"
```

**Parâmetros disponíveis:**
- `--node-name` (obrigatório): Nome único do nó
- `--description`: Descrição do nó
- `--coordinates`: Coordenadas geográficas (lat,lon)
- `--owner-key`: Arquivo de chave de proprietário
- `--label`: Rótulo do sistema de arquivos (padrão: SYNTROPY)
- `--iso`: Caminho para ISO Ubuntu personalizada
- `--discovery-server`: Servidor de descoberta da rede
- `--created-by`: Usuário criador
- `--temp-dir`: Diretório temporário
- `--log-level`: Nível de log (debug, info, warn, error)

### **`syntropy usb-win format`**
Formata dispositivo USB.

```powershell
# Formatar com confirmação
syntropy usb-win format PHYSICALDRIVE1

# Formatar com rótulo personalizado
syntropy usb-win format PHYSICALDRIVE1 --label "MEU-USB"

# Formatar sem confirmação
syntropy usb-win format PHYSICALDRIVE1 --force
```

### **`syntropy usb-win debug`**
Executa diagnóstico completo do ambiente.

```powershell
syntropy usb-win debug
```

**Verificações incluídas:**
- Privilégios de Administrador
- WSL disponível e configurado
- Política de execução do PowerShell
- Ferramentas necessárias instaladas
- Dispositivos USB disponíveis
- Espaço em disco e permissões

---

## ⚠️ Tratamento de Erros

### **Erro: Privilégios de Administrador**
```
[NO_ADMIN_PRIVILEGES] Privilégios de administrador são necessários
```

**Solução:**
1. Feche o PowerShell atual
2. Clique com botão direito no PowerShell
3. Selecione "Executar como administrador"
4. Execute o comando novamente

### **Erro: WSL Não Disponível**
```
[WSL_NOT_AVAILABLE] WSL não está disponível ou configurado
```

**Solução:**
```powershell
# Instalar WSL
wsl --install

# Ou instalar Ubuntu especificamente
wsl --install -d Ubuntu

# Reiniciar se necessário
```

### **Erro: Política de Execução Restrita**
```
[EXECUTION_POLICY_RESTRICTED] Política de execução do PowerShell está restrita
```

**Solução:**
```powershell
# Configurar política para usuário atual
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser

# Ou para processo atual apenas
Set-ExecutionPolicy -ExecutionPolicy Bypass -Scope Process
```

### **Erro: Dispositivo Não Encontrado**
```
[DISK_NOT_FOUND] Dispositivo X não encontrado
```

**Soluções:**
1. Verifique se o USB está conectado
2. Aguarde alguns segundos para reconhecimento
3. Execute `syntropy usb-win list` para verificar
4. Tente conectar em outra porta USB

### **Erro: Dispositivo do Sistema**
```
[SYSTEM_DISK] Dispositivo é um disco do sistema
```

**Solução:**
- Use apenas dispositivos USB removíveis
- NÃO use o disco principal do Windows (geralmente C:)

### **Erro: WSL Montagem Falhou**
```
Falha ao montar disco no WSL
```

**Soluções:**
1. Verifique se WSL está funcionando: `wsl --status`
2. Reinicie o serviço WSL: `wsl --shutdown && wsl`
3. Verifique se não há outros programas usando o USB
4. Execute `syntropy usb-win debug` para diagnóstico

---

## 🔧 Solução de Problemas

### **1. USB Não Aparece na Lista**
```powershell
# Verificar dispositivos via PowerShell
Get-Disk | Where-Object {$_.BusType -eq 'USB'}

# Verificar via WMI
Get-WmiObject -Class Win32_DiskDrive | Where-Object {$_.InterfaceType -eq 'USB'}
```

### **2. WSL Não Detecta Dispositivo**
```powershell
# Verificar status do WSL
wsl --status

# Listar distribuições
wsl --list --verbose

# Testar comando simples
wsl ls /dev/sd*
```

### **3. Erro de Permissão**
```powershell
# Verificar se está executando como Administrador
[Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")

# Verificar política de execução
Get-ExecutionPolicy -List
```

### **4. ISO Não Baixa**
```powershell
# Verificar conectividade
Test-NetConnection releases.ubuntu.com -Port 443

# Verificar espaço em disco
Get-WmiObject -Class Win32_LogicalDisk | Select-Object DeviceID, @{Name="Size(GB)";Expression={[math]::Round($_.Size/1GB,2)}}, @{Name="FreeSpace(GB)";Expression={[math]::Round($_.FreeSpace/1GB,2)}}
```

---

## 📊 Monitoramento e Logs

### **Logs do Sistema**
Os logs são salvos em:
- **Temporário**: `%TEMP%\syntropy-usb\{timestamp}\`
- **Cache**: `%USERPROFILE%\.syntropy\cache\`
- **Configurações**: `%USERPROFILE%\.syntropy\config\`

### **Verificar Logs**
```powershell
# Listar arquivos de log recentes
Get-ChildItem -Path "$env:TEMP\syntropy-usb" -Recurse -Include "*.log" | Sort-Object LastWriteTime -Descending

# Ver último log
Get-Content -Path "$env:TEMP\syntropy-usb\*\*.log" -Tail 50
```

---

## 🚀 Fluxo de Criação de Nó

### **1. Preparação**
```powershell
# 1. Executar como Administrador
# 2. Verificar ambiente
syntropy usb-win debug

# 3. Listar dispositivos
syntropy usb-win list
```

### **2. Criação**
```powershell
# Criar USB para nó
syntropy usb-win create --node-name "node-01" --description "Nó principal"
```

### **3. Verificação**
```powershell
# Verificar se USB foi criado corretamente
# (O comando mostrará informações de sucesso)
```

### **4. Boot do Nó**
1. Conectar USB no hardware do nó
2. Configurar BIOS/UEFI para boot via USB
3. Iniciar o sistema
4. A configuração cloud-init será aplicada automaticamente

---

## 🔒 Considerações de Segurança

### **Certificados TLS**
- CA própria gerada (4096 bits)
- Certificados de nó (2048 bits)
- Validade: CA (10 anos), Nó (1 ano)

### **Chaves SSH**
- Algoritmo: ED25519 (padrão), RSA 2048 (fallback)
- Armazenamento: `%USERPROFILE%\.syntropy\keys\`
- Chave privada NÃO é enviada para o nó

### **Permissões**
- Chaves privadas: 0600
- Certificados: 0644
- Scripts: 0755

---

## 📞 Suporte

### **Diagnóstico Automático**
```powershell
syntropy usb-win debug
```

### **Informações do Sistema**
```powershell
# Informações do Windows
Get-ComputerInfo | Select-Object WindowsProductName, WindowsVersion, TotalPhysicalMemory

# Informações do WSL
wsl --version
wsl --list --verbose

# Informações do PowerShell
$PSVersionTable
```

### **Logs de Erro**
Se encontrar problemas:
1. Execute `syntropy usb-win debug`
2. Verifique os logs em `%TEMP%\syntropy-usb\`
3. Anote a mensagem de erro exata
4. Verifique se todos os pré-requisitos estão atendidos

---

## 📚 Referências

- [WSL Documentation](https://docs.microsoft.com/en-us/windows/wsl/)
- [PowerShell Execution Policies](https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_execution_policies)
- [Windows Disk Management](https://docs.microsoft.com/en-us/windows-server/storage/disk-management/overview-of-disk-management)
- [Cloud-init Documentation](https://cloudinit.readthedocs.io/)
- [Ubuntu Server Installation](https://ubuntu.com/server/docs/installation)

---

**💡 Dica:** Para a melhor experiência, sempre execute os comandos como Administrador e mantenha o WSL atualizado!
