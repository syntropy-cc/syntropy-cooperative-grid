# Syntropy CLI Manager - Windows Workflow

Este documento descreve o workflow completo para compilar, executar e gerenciar a aplicação CLI do Syntropy no Windows.

## 📋 Pré-requisitos

### Software Necessário
- **Go 1.22.5 ou superior** - [Download](https://golang.org/dl/)
- **Git** - [Download](https://git-scm.com/download/win)
- **PowerShell 5.1 ou superior** (incluído no Windows 10/11)

### Verificação dos Pré-requisitos
```powershell
# Verificar Go
go version

# Verificar Git
git --version

# Verificar PowerShell
$PSVersionTable.PSVersion
```

## 🚀 Scripts Disponíveis

### 1. `build-windows.ps1` - Script Principal de Build
Script PowerShell otimizado para compilação e execução básica.

**Uso:**
```powershell
.\build-windows.ps1 [ação] [argumentos]
```

**Ações disponíveis:**
- `build` - Compilar a aplicação (padrão)
- `run` - Executar a aplicação compilada
- `test` - Executar testes
- `install` - Instalar binário globalmente
- `uninstall` - Desinstalar binário
- `clean` - Limpar diretório de build
- `help` - Mostrar ajuda

**Exemplos:**
```powershell
# Compilar
.\build-windows.ps1 build

# Executar
.\build-windows.ps1 run

# Executar com argumentos
.\build-windows.ps1 run '--help'
.\build-windows.ps1 run 'setup run --force'

# Executar testes
.\build-windows.ps1 test

# Instalar globalmente
.\build-windows.ps1 install
```

### 2. `dev-workflow.ps1` - Workflow de Desenvolvimento
Script completo para desenvolvimento com verificações de qualidade.

**Uso:**
```powershell
.\dev-workflow.ps1 [ação] [argumentos]
```

**Ações disponíveis:**
- `setup` - Configurar ambiente de desenvolvimento
- `build` - Compilar a aplicação
- `test` - Executar testes
- `run` - Executar a aplicação
- `dev` - Modo desenvolvimento (build + test + quality)
- `install` - Instalar binário globalmente
- `uninstall` - Desinstalar binário
- `status` - Verificar status do sistema
- `clean` - Limpar diretórios
- `help` - Mostrar ajuda

**Exemplos:**
```powershell
# Setup inicial
.\dev-workflow.ps1 setup

# Modo desenvolvimento completo
.\dev-workflow.ps1 dev

# Verificar status
.\dev-workflow.ps1 status

# Executar com argumentos
.\dev-workflow.ps1 run 'setup status'
```

### 3. `automation-workflow.ps1` - Workflow de Automação
Script avançado para CI/CD e automação completa.

**Uso:**
```powershell
.\automation-workflow.ps1 [ação] [argumentos]
```

**Ações disponíveis:**
- `full` - Workflow completo (build + test + quality + verify)
- `build` - Apenas compilação
- `test` - Apenas testes
- `deploy` - Deploy e distribuição
- `ci` - Workflow de CI/CD completo
- `release` - Preparar release
- `help` - Mostrar ajuda

**Exemplos:**
```powershell
# Workflow completo
.\automation-workflow.ps1 full

# CI/CD completo
.\automation-workflow.ps1 ci

# Preparar release
.\automation-workflow.ps1 release
```

### 4. `quick-start.bat` - Script de Início Rápido
Script batch simples para setup rápido e execução interativa.

**Uso:**
```cmd
quick-start.bat
```

Este script oferece um menu interativo para:
- Compilar a aplicação
- Executar com diferentes argumentos
- Mostrar ajuda
- Executar setup
- Verificar status

### 5. `run-cli.bat` - Executor Simples
Script batch para executar a aplicação compilada.

**Uso:**
```cmd
run-cli.bat [argumentos]
```

**Exemplos:**
```cmd
# Executar sem argumentos
run-cli.bat

# Executar com argumentos
run-cli.bat --help
run-cli.bat setup run --force
```

## 📁 Estrutura de Diretórios

Após executar os scripts, a seguinte estrutura será criada:

```
cli/
├── build/                    # Binários compilados
│   └── syntropy.exe         # Executável principal
├── logs/                     # Logs de execução
│   ├── automation-*.log     # Logs de automação
│   ├── test-results-*.txt   # Resultados de testes
│   ├── quality-results-*.txt # Resultados de qualidade
│   └── binary-tests-*.txt   # Testes do binário
├── dist/                     # Artefatos de distribuição
│   ├── syntropy.exe         # Binário para distribuição
│   └── build-info.txt       # Informações do build
├── temp/                     # Arquivos temporários
└── scripts/                  # Scripts de workflow
```

## 🔧 Workflows Recomendados

### Para Desenvolvimento Diário
```powershell
# 1. Setup inicial (apenas uma vez)
.\dev-workflow.ps1 setup

# 2. Desenvolvimento
.\dev-workflow.ps1 dev

# 3. Executar aplicação
.\dev-workflow.ps1 run 'setup --help'
```

### Para Testes e Qualidade
```powershell
# Executar todos os testes e verificações
.\automation-workflow.ps1 full
```

### Para CI/CD
```powershell
# Workflow completo de CI/CD
.\automation-workflow.ps1 ci
```

### Para Iniciantes
```cmd
# Use o script de início rápido
quick-start.bat
```

## 📊 Monitoramento e Logs

### Logs de Automação
Os scripts de automação geram logs detalhados em `logs/`:
- `automation-*.log` - Log principal de execução
- `test-results-*.txt` - Resultados de testes
- `quality-results-*.txt` - Resultados de verificações de qualidade
- `binary-tests-*.txt` - Testes do binário
- `workflow-results-*.json` - Resultados do workflow em JSON

### Relatórios HTML
O script de automação gera relatórios HTML em `logs/ci-report-*.html` com:
- Resumo do build
- Status de cada etapa
- Links para logs detalhados
- Informações de versão e commit

## 🛠️ Comandos Úteis da CLI

Após compilar, você pode usar os seguintes comandos:

```powershell
# Mostrar ajuda
.\build\syntropy.exe --help

# Mostrar versão
.\build\syntropy.exe --version

# Ajuda do setup
.\build\syntropy.exe setup --help

# Executar setup
.\build\syntropy.exe setup run --force

# Verificar status
.\build\syntropy.exe setup status

# Validar ambiente
.\build\syntropy.exe setup validate

# Reset configuração
.\build\syntropy.exe setup reset --force
```

## 🔍 Solução de Problemas

### Erro: "Go não está instalado"
```powershell
# Verificar se Go está no PATH
go version

# Se não estiver, reinstalar Go e adicionar ao PATH
# Download: https://golang.org/dl/
```

### Erro: "main.go não encontrado"
```powershell
# Verificar se está no diretório correto
Get-Location
# Deve estar em: .../manager/interfaces/cli/

# Navegar para o diretório correto
cd "C:\caminho\para\syntropy-cooperative-grid\manager\interfaces\cli"
```

### Erro: "Falha na compilação"
```powershell
# Limpar e tentar novamente
.\build-windows.ps1 clean
.\build-windows.ps1 build

# Verificar dependências
go mod download
go mod tidy
```

### Erro: "Binário não encontrado"
```powershell
# Verificar se o build foi executado
Test-Path ".\build\syntropy.exe"

# Se não existir, executar build
.\build-windows.ps1 build
```

## 📈 Performance e Otimização

### Build Rápido
Para builds rápidos durante desenvolvimento:
```powershell
.\build-windows.ps1 build
```

### Build Completo
Para builds com todas as verificações:
```powershell
.\automation-workflow.ps1 full
```

### Limpeza
Para limpar arquivos temporários:
```powershell
.\build-windows.ps1 clean
# ou
.\dev-workflow.ps1 clean
```

## 🔐 Segurança

### Verificação de Integridade
Os scripts incluem verificações de:
- Integridade das dependências (`go mod verify`)
- Formatação de código (`go fmt`)
- Análise estática (`go vet`)
- Linting (`golangci-lint`)

### Logs Seguros
Os logs não contêm informações sensíveis e podem ser compartilhados para debugging.

## 📞 Suporte

### Logs para Suporte
Para obter suporte, forneça:
1. Logs de execução (`logs/automation-*.log`)
2. Resultados de testes (`logs/test-results-*.txt`)
3. Informações do sistema:
   ```powershell
   go version
   git --version
   $PSVersionTable.PSVersion
   ```

### Comandos de Diagnóstico
```powershell
# Verificar status completo
.\dev-workflow.ps1 status

# Executar workflow completo para diagnóstico
.\automation-workflow.ps1 full
```

## 🎯 Próximos Passos

1. **Primeira Execução**: Use `quick-start.bat` para setup inicial
2. **Desenvolvimento**: Use `dev-workflow.ps1 dev` para desenvolvimento
3. **Testes**: Use `automation-workflow.ps1 full` para testes completos
4. **Produção**: Use `automation-workflow.ps1 release` para releases

---

**Nota**: Este workflow foi projetado para ser executado no Windows com PowerShell. Para outros sistemas operacionais, use os scripts correspondentes (Makefile para Linux/macOS).
