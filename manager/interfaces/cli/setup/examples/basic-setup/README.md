# Exemplo de Setup Básico - Syntropy CLI

Este diretório contém exemplos práticos de configuração básica do setup component para o Syntropy CLI.

## Visão Geral

O setup básico demonstra como configurar o Syntropy CLI em diferentes sistemas operacionais usando configurações padrão e comandos simples.

## Estrutura dos Arquivos

- `setup-basic.sh` - Script de setup básico para Linux/macOS
- `setup-basic.ps1` - Script de setup básico para Windows
- `config-example.yaml` - Arquivo de configuração de exemplo
- `README.md` - Este arquivo com documentação completa

## Requisitos

### Linux/macOS
- Shell: Bash 4.0+
- Permissões: Usuário regular (algumas operações podem precisar de sudo)

### Windows
- PowerShell: Versão 5.1 ou superior
- Permissões: Administrador para instalação como serviço

## Configuração Básica

### 1. Usando o Script Automático

**Linux/macOS:**
```bash
chmod +x setup-basic.sh
./setup-basic.sh
```

**Windows:**
```powershell
# Execute com permissões de administrador
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
.\setup-basic.ps1
```

### 2. Setup Manual

```bash
# Linux/macOS
go run . syntropy setup

# Windows 
# (No PowerShell com permissões de administrador)
.\syntropy.exe setup
```

## Opções de Configuração

### Configuração Básica (config-example.yaml)
```yaml
manager:
  home_dir: "~/syntropy"
  log_level: "info"
  api_endpoint: "https://api.syntropy.io"
  directories:
    config: "config"
    logs: "logs"
    data: "data"
  default_paths:
    config: "config/manager.yaml"
    owner_key: "config/owner.key"

owner_key:
  type: "Ed25519"
  path: "config/owner.key"

environment:
  os: "{{.OS}}"
  architecture: "{{.Arch}}"
  home_dir: "{{.HomeDir}}"
```

## Exemplos de Comandos

### Setup Completo
```bash
# Setup básico sem forçar
syntropy setup

# Setup forçado (ignora validações)
syntropy setup --force

# Setup com instalação de serviço
syntropy setup --install-service
```

### Validação
```bash
# Validação apenas do ambiente
syntropy setup --validate-only

# Status atual do setup
syntropy setup status
```

### Configuração Avançada
```bash
# Setup com arquivo de configuração customizado
syntropy setup --config /custom/path/config.yaml

# Setup com diretório home customizado
syntropy setup --home-dir /custom/home/syntropy
```

## Verificação do Setup

Após a execução do setup, você pode verificar se tudo foi configurado corretamente:

```bash
# Verificar status
syntropy setup status

# Verificar configuração
cat ~/.syntropy/config/manager.yaml

# Listar arquivos criados
ls -la ~/.syntropy/
```

## Estrutura de Diretórios Criada

Após um setup bem-sucedido, a seguinte estrutura será criada:

```
~/.syntropy/                 # Diretório principal
├── config/                 # Configurações
│   ├── manager.yaml        # Configuração principal
│   └── owner.key          # Chave do proprietário
├── logs/                  # Logs do sistema
├── data/                  # Dados do usuário
└── services/              # Scripts de serviço (se aplicável)
    ├── install.sh         # Linux/macOS
    └── install.ps1        # Windows
```

## Troubleshooting Básico

### Problema: Permissões nesta plataforma
```bash
# Linux/macOS - Execute com sudo para instalação como serviço
sudo ./setup-basic.sh

# Ou use --user apenas para configuração local
./setup-basic.sh --user-only
```

### Problema: PowerShell execution policy no Windows
```powershell
# Habilitar execução de scripts (execute como administrador)
Set-ExecutionPolicy RemoteSigned -Force
```

### Problema: Espaço em disco insuficiente
- Verifique se há pelo menos 1GB de espaço livre
- Use `df -h` (Linux) ou `Get-WmiObject -Class Win32_LogicalDisk` (Windows) para verificar

## Próximos Passos

Após completar o setup básico, consulte:
- `../advanced-setup/` para configurações avançadas
- `../validation-tests/` para testes de sistema
- `../../README.md` para documentação completa

## Logs e Depuração

O setup mantém logs detalhados que podem ser consultados em caso de problemas:

```bash
# Linux/macOS
tail -f ~/.syntropy/logs/setup.log

# Windows
Get-Content "$env:USERPROFILE\.syntropy\logs\setup.log" -Tail 50 -Wait
```

Para debug avançado, use:
```bash
# Setup com logs detalhados
syntropy setup --log-level debug
```

## Suporte

Se encontrar problemas ou tiver dúvidas:
1. Consulte a documentação principal: `../../GUIDE.md`
2. Verifique logs em `~/.syntropy/logs/`
3. Execute testes de validação: `../validation-tests/test-environment.sh`
