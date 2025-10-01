# Syntropy CLI Manager

## 🚀 Quick Start

### Linux/WSL
```bash
# Instalação completa (recomendada)
./build.sh

# Instalação alternativa
./install.sh
```

### Windows
```cmd
# Build completo
build.bat

# Build simplificado
scripts\shared\build-and-test.bat
```

### ⚠️ Importante - Versões Disponíveis:
- **Versão Completa**: Funcionalidade total com `setup-component`
- **Versão Simplificada**: Versão demo para testes básicos

## 📁 Estrutura Organizada

```
cli/
├── build.sh                    # Script principal para Linux/WSL
├── build.bat                   # Script principal para Windows
├── install.sh                  # Instalação para Linux/WSL
├── scripts/                    # Scripts organizados por plataforma
│   ├── windows/               # Scripts específicos para Windows
│   │   ├── build-windows.ps1  # Build principal
│   │   ├── dev-workflow.ps1   # Desenvolvimento
│   │   ├── automation-workflow.ps1 # CI/CD
│   │   ├── quick-start.bat    # Início rápido
│   │   ├── run-cli.bat        # Executor
│   │   ├── run-examples.bat   # Exemplos
│   │   └── setup-environment.ps1 # Setup automático
│   ├── linux/                 # Scripts específicos para Linux
│   │   ├── install-syntropy.sh # Instalação
│   │   └── build-and-test.sh  # Build completo
│   ├── shared/                # Scripts compartilhados
│   │   ├── build-and-test.bat # Build para Windows
│   │   └── start-here.bat     # Entrada principal
│   └── README.md              # Documentação dos scripts
├── build/                     # Binários compilados
│   ├── syntropy-windows.exe   # Para Windows
│   └── syntropy-linux         # Para Linux
└── docs/                      # Documentação
    ├── QUICK_START.md         # Início rápido
    └── WINDOWS_WORKFLOW.md    # Workflow Windows
```

## 🎯 Scripts Principais

### Para Linux/WSL
- **`./install.sh`** - Instalação alternativa
- **`./build.sh`** - Build completo

### Para Windows
- **`build.bat`** - Build e teste

## 📋 Funcionalidades

- ✅ **Compilação** para Windows (.exe) e Linux
- ✅ **Testes** automáticos dos binários
- ✅ **Configuração** automática de dependências
- ✅ **Execução** da aplicação
- ✅ **Documentação** completa

## 🧪 Testar a Aplicação

### Windows
```cmd
build\syntropy-windows.exe --help
build\syntropy-windows.exe setup run --force
```

### Linux
```bash
./build/syntropy-linux --help
./build/syntropy-linux setup run --force
```

## 📚 Documentação

- **`QUICK_START.md`** - Início rápido
- **`WINDOWS_WORKFLOW.md`** - Workflow Windows
- **`scripts/README.md`** - Documentação dos scripts

## 🛠️ Pré-requisitos

- **Go 1.22.5+** - [Download](https://golang.org/dl/)
- **Git** (opcional)

## 🎉 Próximos Passos

1. **Execute** `./install.sh` (Linux) ou `build.bat` (Windows)
2. **Teste** os binários gerados
3. **Execute** a aplicação com `setup run --force`
4. **Consulte** a documentação para mais detalhes

---

**Scripts organizados e prontos para uso!** 🚀