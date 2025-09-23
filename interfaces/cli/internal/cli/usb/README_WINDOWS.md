# 🪟 USB Syntropy - Versão Windows

## 📋 Visão Geral

Esta é a versão específica para Windows do módulo USB Syntropy, otimizada para criação de nós da Syntropy Cooperative Grid com validações robustas e tratamento de erros específicos do ambiente Windows.

## 🚀 Início Rápido

### **1. Configuração Automática**
```powershell
# Execute como Administrador
.\setup-windows.ps1
```

### **2. Verificar Ambiente**
```powershell
syntropy usb-win debug
```

### **3. Listar Dispositivos**
```powershell
syntropy usb-win list
```

### **4. Criar USB Bootável**
```powershell
syntropy usb-win create --node-name "meu-no"
```

---

## 🔧 Características da Versão Windows

### **✅ Validações Robustas**
- Privilégios de Administrador
- WSL disponível e configurado
- Política de execução do PowerShell
- Ferramentas necessárias instaladas
- Dispositivos USB válidos e seguros

### **✅ Tratamento de Erros Específicos**
- Erros de permissão (UAC, Administrador)
- Problemas de WSL (instalação, configuração)
- Dispositivos do sistema (proteção)
- Falhas de montagem WSL
- Problemas de rede (download ISO)

### **✅ Comandos Otimizados**
- `usb-win list`: Lista dispositivos com informações detalhadas
- `usb-win create`: Criação com validações completas
- `usb-win format`: Formatação segura com confirmação
- `usb-win debug`: Diagnóstico completo do ambiente

---

## 📁 Arquivos da Versão Windows

```
usb/
├── windows-only.go          # Implementação específica para Windows
├── windows-commands.go      # Comandos CLI para Windows
├── setup-windows.ps1       # Script de configuração automática
├── GUIA_WINDOWS.md         # Guia completo de uso
└── README_WINDOWS.md       # Este arquivo
```

---

## 🛠️ Pré-requisitos

### **Sistema Operacional**
- Windows 10 (versão 1903+)
- Windows 11 (todas as versões)
- Arquitetura x64 (64-bit)

### **Software Necessário**
- PowerShell 5.1+
- WSL 2 com Ubuntu
- Privilégios de Administrador

### **Ferramentas WSL**
```bash
# Execute no WSL após instalação
sudo apt update
sudo apt install -y gdisk dosfstools
```

---

## 📖 Comandos Disponíveis

### **`syntropy usb-win list`**
Lista dispositivos USB com validações Windows.

**Flags:**
- `--format`: Formato de saída (table, json, yaml)

**Exemplo:**
```powershell
syntropy usb-win list --format json
```

### **`syntropy usb-win create`**
Cria USB bootável com validações completas.

**Flags obrigatórias:**
- `--node-name`: Nome único do nó

**Flags opcionais:**
- `--description`: Descrição do nó
- `--coordinates`: Coordenadas geográficas
- `--owner-key`: Arquivo de chave de proprietário
- `--label`: Rótulo do sistema de arquivos
- `--iso`: Caminho para ISO Ubuntu
- `--discovery-server`: Servidor de descoberta
- `--created-by`: Usuário criador
- `--temp-dir`: Diretório temporário
- `--log-level`: Nível de log

**Exemplo:**
```powershell
syntropy usb-win create PHYSICALDRIVE1 \
  --node-name "node-01" \
  --description "Nó principal" \
  --coordinates "-23.5505,-46.6333"
```

### **`syntropy usb-win format`**
Formata dispositivo USB com validações.

**Flags:**
- `--label`: Rótulo do sistema de arquivos
- `--force`: Não pedir confirmação

**Exemplo:**
```powershell
syntropy usb-win format PHYSICALDRIVE1 --label "SYNTROPY"
```

### **`syntropy usb-win debug`**
Executa diagnóstico completo do ambiente.

**Exemplo:**
```powershell
syntropy usb-win debug
```

---

## 🔍 Tratamento de Erros

### **Tipos de Erro**
- `NO_ADMIN_PRIVILEGES`: Privilégios de Administrador necessários
- `WSL_NOT_AVAILABLE`: WSL não instalado ou configurado
- `EXECUTION_POLICY_RESTRICTED`: Política PowerShell restrita
- `DISK_NOT_FOUND`: Dispositivo USB não encontrado
- `SYSTEM_DISK`: Tentativa de usar disco do sistema
- `WSL_MOUNT_FAILED`: Falha na montagem WSL

### **Soluções Automáticas**
- Verificação de privilégios
- Instalação automática do WSL
- Configuração de política de execução
- Validação de dispositivos
- Diagnóstico de ambiente

---

## 🔧 Configuração

### **Configuração Automática**
```powershell
# Execute como Administrador
.\setup-windows.ps1
```

### **Configuração Manual**
1. **Instalar WSL:**
   ```powershell
   wsl --install -d Ubuntu
   ```

2. **Configurar PowerShell:**
   ```powershell
   Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
   ```

3. **Instalar ferramentas WSL:**
   ```bash
   sudo apt update
   sudo apt install -y gdisk dosfstools
   ```

---

## 📊 Estrutura de Dados

### **WindowsOnlyUSBDevice**
```go
type WindowsOnlyUSBDevice struct {
    DiskNumber     int    `json:"disk_number"`
    FriendlyName   string `json:"friendly_name"`
    Size           int64  `json:"size"`
    SizeFormatted  string `json:"size_formatted"`
    SerialNumber   string `json:"serial_number"`
    BusType        string `json:"bus_type"`
    Model          string `json:"model"`
    IsSystem       bool   `json:"is_system"`
    IsBoot         bool   `json:"is_boot"`
    IsOffline      bool   `json:"is_offline"`
    PartitionCount int    `json:"partition_count"`
    Status         string `json:"status"`
}
```

### **WindowsOnlyConfig**
```go
type WindowsOnlyConfig struct {
    NodeName        string `json:"node_name"`
    NodeDescription string `json:"node_description"`
    Coordinates     string `json:"coordinates"`
    OwnerKeyFile    string `json:"owner_key_file"`
    Label           string `json:"label"`
    ISOPath         string `json:"iso_path"`
    DiscoveryServer string `json:"discovery_server"`
    CreatedBy       string `json:"created_by"`
    ExecutionPolicy string `json:"execution_policy"`
    PowerShellPath  string `json:"powershell_path"`
    WSLDistro       string `json:"wsl_distro"`
    TempDir         string `json:"temp_dir"`
    LogLevel        string `json:"log_level"`
}
```

### **WindowsOnlyError**
```go
type WindowsOnlyError struct {
    Code        string `json:"code"`
    Message     string `json:"message"`
    Suggestion  string `json:"suggestion"`
    ErrorType   string `json:"error_type"`
    Recoverable bool   `json:"recoverable"`
}
```

---

## 🔒 Segurança

### **Validações de Segurança**
- Verificação de dispositivos do sistema
- Validação de permissões
- Proteção contra formatação acidental
- Verificação de integridade de dispositivos

### **Certificados e Chaves**
- CA própria (4096 bits)
- Certificados de nó (2048 bits)
- Chaves SSH ED25519
- Armazenamento seguro local

---

## 📈 Performance

### **Otimizações**
- Cache inteligente de ISOs
- Validações em paralelo
- Limpeza automática de recursos
- Timeouts configuráveis

### **Monitoramento**
- Logs detalhados
- Progresso em tempo real
- Diagnóstico automático
- Métricas de performance

---

## 🧪 Testes

### **Testes Automáticos**
```powershell
# Verificar ambiente
syntropy usb-win debug

# Testar listagem
syntropy usb-win list

# Verificar script de configuração
.\setup-windows.ps1 -Force
```

### **Testes Manuais**
1. Conectar dispositivo USB
2. Executar `syntropy usb-win list`
3. Criar USB com `syntropy usb-win create`
4. Verificar boot no hardware

---

## 📞 Suporte

### **Diagnóstico**
```powershell
# Diagnóstico completo
syntropy usb-win debug

# Verificação do ambiente
.\$env:USERPROFILE\.syntropy\scripts\verify-environment.ps1
```

### **Logs**
- **Temporário**: `%TEMP%\syntropy-usb\`
- **Cache**: `%USERPROFILE%\.syntropy\cache\`
- **Configurações**: `%USERPROFILE%\.syntropy\config\`

### **Documentação**
- [GUIA_WINDOWS.md](GUIA_WINDOWS.md) - Guia completo
- [SOLUCAO_ERROS.md](SOLUCAO_ERROS.md) - Soluções de problemas
- [README.md](README.md) - Documentação geral

---

## 🚀 Roadmap

### **Próximas Versões**
- [ ] Suporte a Windows Server
- [ ] Interface gráfica (GUI)
- [ ] Integração com Windows Update
- [ ] Suporte a múltiplas distribuições WSL
- [ ] Backup automático de configurações

### **Melhorias Planejadas**
- [ ] Validação de hardware
- [ ] Suporte a RAID
- [ ] Integração com Active Directory
- [ ] Monitoramento de performance
- [ ] Relatórios automatizados

---

## 📚 Referências

- [WSL Documentation](https://docs.microsoft.com/en-us/windows/wsl/)
- [PowerShell Documentation](https://docs.microsoft.com/en-us/powershell/)
- [Windows Disk Management](https://docs.microsoft.com/en-us/windows-server/storage/disk-management/)
- [Cloud-init Documentation](https://cloudinit.readthedocs.io/)
- [Ubuntu Server Guide](https://ubuntu.com/server/docs)

---

**💡 Dica:** Para a melhor experiência, sempre execute os comandos como Administrador e mantenha o WSL atualizado!
