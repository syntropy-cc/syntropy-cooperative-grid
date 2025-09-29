# Syntropy CLI Manager - Windows Quick Start

## 🚀 Início Rápido

### 1. Pré-requisitos
- **Go 1.22.5+** - [Download](https://golang.org/dl/)
- **Git** - [Download](https://git-scm.com/download/win)
- **PowerShell 5.1+** (incluído no Windows 10/11)

### 2. Execução Rápida
```cmd
# Opção 1: Script interativo (recomendado para iniciantes)
quick-start.bat

# Opção 2: Build e execução direta
.\build-windows.ps1 build
.\build-windows.ps1 run '--help'
```

### 3. Comandos Essenciais
```powershell
# Compilar
.\build-windows.ps1 build

# Executar
.\build-windows.ps1 run

# Executar com argumentos
.\build-windows.ps1 run 'setup run --force'

# Executar testes
.\build-windows.ps1 test

# Instalar globalmente
.\build-windows.ps1 install
```

## 📋 Scripts Disponíveis

| Script | Descrição | Uso |
|--------|-----------|-----|
| `quick-start.bat` | Setup interativo | `quick-start.bat` |
| `build-windows.ps1` | Build e execução básica | `.\build-windows.ps1 build` |
| `dev-workflow.ps1` | Desenvolvimento completo | `.\dev-workflow.ps1 dev` |
| `automation-workflow.ps1` | CI/CD e automação | `.\automation-workflow.ps1 full` |
| `run-cli.bat` | Executor simples | `run-cli.bat --help` |

## 🎯 Workflows Recomendados

### Para Iniciantes
```cmd
quick-start.bat
```

### Para Desenvolvimento
```powershell
.\dev-workflow.ps1 setup    # Setup inicial
.\dev-workflow.ps1 dev      # Desenvolvimento
.\dev-workflow.ps1 run      # Executar
```

### Para Testes Completos
```powershell
.\automation-workflow.ps1 full
```

## 📁 Estrutura Após Build

```
cli/
├── build/
│   └── syntropy.exe        # Executável principal
├── logs/                   # Logs de execução
├── dist/                   # Artefatos de distribuição
└── temp/                   # Arquivos temporários
```

## 🔧 Comandos da CLI

```powershell
# Ajuda geral
.\build\syntropy.exe --help

# Ajuda do setup
.\build\syntropy.exe setup --help

# Executar setup
.\build\syntropy.exe setup run --force

# Verificar status
.\build\syntropy.exe setup status

# Validar ambiente
.\build\syntropy.exe setup validate
```

## 🛠️ Solução de Problemas

### Go não encontrado
```powershell
go version  # Verificar se Go está instalado
# Se não estiver: https://golang.org/dl/
```

### Erro de compilação
```powershell
.\build-windows.ps1 clean
.\build-windows.ps1 build
```

### Binário não encontrado
```powershell
Test-Path ".\build\syntropy.exe"  # Verificar se existe
.\build-windows.ps1 build         # Compilar se necessário
```

## 📊 Logs e Monitoramento

- **Logs principais**: `logs/automation-*.log`
- **Resultados de testes**: `logs/test-results-*.txt`
- **Relatórios HTML**: `logs/ci-report-*.html`

## 📞 Suporte

Para suporte, forneça:
1. Logs de execução
2. Versão do Go: `go version`
3. Versão do PowerShell: `$PSVersionTable.PSVersion`

---

**Documentação completa**: [WINDOWS_WORKFLOW.md](WINDOWS_WORKFLOW.md)
