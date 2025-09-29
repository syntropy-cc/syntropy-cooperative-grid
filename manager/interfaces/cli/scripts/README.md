# Syntropy CLI Manager - Scripts Directory

Esta pasta contém todos os scripts organizados por plataforma e funcionalidade.

## 📁 Estrutura

```
scripts/
├── windows/          # Scripts específicos para Windows
├── linux/            # Scripts específicos para Linux/WSL
├── shared/           # Scripts compartilhados entre plataformas
└── README.md         # Este arquivo
```

## 🪟 Windows Scripts (`windows/`)

### Scripts Principais
- **`build-windows.ps1`** - Script principal de build e execução
- **`dev-workflow.ps1`** - Workflow completo de desenvolvimento
- **`automation-workflow.ps1`** - Workflow de automação e CI/CD

### Scripts Auxiliares
- **`quick-start.bat`** - Setup rápido e execução interativa
- **`run-cli.bat`** - Executor simples da aplicação
- **`run-examples.bat`** - Execução de exemplos da CLI
- **`setup-environment.ps1`** - Configuração automática do ambiente

### Como Usar
```powershell
# Build básico
.\scripts\windows\build-windows.ps1 build

# Desenvolvimento completo
.\scripts\windows\dev-workflow.ps1 dev

# Automação completa
.\scripts\windows\automation-workflow.ps1 full
```

## 🐧 Linux Scripts (`linux/`)

### Scripts Principais
- **`install-syntropy.sh`** - Script mais simples para instalação
- **`build-and-test.sh`** - Script completo de build e teste

### Como Usar
```bash
# Instalação simples
./scripts/linux/install-syntropy.sh

# Build completo
./scripts/linux/build-and-test.sh
```

## 🔄 Shared Scripts (`shared/`)

### Scripts Compartilhados
- **`build-and-test.bat`** - Build e teste para Windows (compatível)
- **`start-here.bat`** - Script de entrada principal

### Como Usar
```cmd
# Build compartilhado
scripts\shared\build-and-test.bat

# Entrada principal
scripts\shared\start-here.bat
```

## 🚀 Scripts de Entrada (Raiz)

### Scripts Principais
- **`build.sh`** - Script principal para Linux/WSL
- **`build.bat`** - Script principal para Windows

### Como Usar
```bash
# Linux/WSL
./build.sh

# Windows
build.bat
```

## 📋 Funcionalidades por Script

### Build e Compilação
- ✅ Compilação para Windows (.exe)
- ✅ Compilação para Linux
- ✅ Configuração automática de dependências
- ✅ Verificação de integridade

### Testes e Qualidade
- ✅ Testes unitários
- ✅ Testes com cobertura
- ✅ Verificações de qualidade (go vet, golangci-lint)
- ✅ Testes de race condition

### Automação
- ✅ Workflow completo de CI/CD
- ✅ Geração de logs detalhados
- ✅ Relatórios HTML
- ✅ Artefatos de distribuição

### Execução
- ✅ Execução direta da aplicação
- ✅ Suporte a argumentos
- ✅ Menu interativo
- ✅ Exemplos de uso

## 🎯 Workflows Recomendados

### Para Iniciantes
```bash
# Linux/WSL
./build.sh

# Windows
build.bat
```

### Para Desenvolvimento
```bash
# Linux/WSL
./scripts/linux/install-syntropy.sh

# Windows
.\scripts\windows\dev-workflow.ps1 dev
```

### Para Testes Completos
```bash
# Linux/WSL
./scripts/linux/build-and-test.sh

# Windows
.\scripts\windows\automation-workflow.ps1 full
```

## 🛠️ Manutenção

### Adicionar Novos Scripts
1. Coloque na pasta apropriada (`windows/`, `linux/`, `shared/`)
2. Atualize este README
3. Teste em todas as plataformas suportadas

### Remover Scripts
1. Verifique se não há dependências
2. Atualize este README
3. Teste os scripts restantes

## 📞 Suporte

Para problemas com scripts:
1. Verifique se está na pasta correta
2. Execute com permissões apropriadas
3. Consulte os logs gerados
4. Verifique a documentação específica de cada script

---

**Scripts organizados e prontos para uso!** 🚀

