# Syntropy CLI Manager - Quick Start

## 🎯 Objetivo Simples

Compilar o Go, gerar binários .exe para Windows e equivalente para Linux, e testar a aplicação Syntropy.

## 🚀 Instalação Rápida

### Opção 1: Script Mais Simples (Recomendado)
```bash
# Execute este comando no diretório cli/
./install-syntropy.sh
```

### Opção 2: Script Completo
```bash
# Para Linux/WSL
./build-and-test.sh

# Para Windows
build-and-test.bat
```

## 📁 Resultado

Após executar, você terá:
```
build/
├── syntropy-windows.exe    # Para Windows
└── syntropy-linux          # Para Linux
```

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

## 📋 Pré-requisitos

- **Go 1.22.5+** - [Download](https://golang.org/dl/)
- **Git** (opcional)

## 🎯 Comandos Úteis

```bash
# Ajuda
./build/syntropy-linux --help

# Versão
./build/syntropy-linux --version

# Setup
./build/syntropy-linux setup run --force

# Status
./build/syntropy-linux setup status
```

## 🛠️ Solução de Problemas

### Go não encontrado
```bash
go version  # Verificar se Go está instalado
# Se não estiver: https://golang.org/dl/
```

### Erro de compilação
```bash
rm -rf build/  # Limpar
./install-syntropy.sh  # Tentar novamente
```

---

**Simples e direto!** 🚀

Execute `./install-syntropy.sh` e pronto!
