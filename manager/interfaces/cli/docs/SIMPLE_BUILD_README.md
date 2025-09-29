# Syntropy CLI Manager - Simple Build & Test

## 🎯 Objetivo

Compilar a aplicação CLI do Syntropy, gerar binários para Windows (.exe) e Linux, e testar a aplicação.

## 📋 Pré-requisitos

- **Go 1.22.5 ou superior** - [Download](https://golang.org/dl/)
- **Git** (opcional, para informações de commit)

## 🚀 Como Usar

### No Windows

```cmd
# Execute o script batch
build-and-test.bat
```

### No Linux/WSL

```bash
# Execute o script shell
./build-and-test.sh
```

## 📁 Arquivos Gerados

Após a execução, você terá:

```
build/
├── syntropy-windows.exe    # Executável para Windows
└── syntropy-linux          # Executável para Linux
```

## 🧪 Testando a Aplicação

### Windows
```cmd
# Mostrar ajuda
build\syntropy-windows.exe --help

# Mostrar versão
build\syntropy-windows.exe --version

# Executar setup
build\syntropy-windows.exe setup run --force

# Verificar status
build\syntropy-windows.exe setup status
```

### Linux
```bash
# Mostrar ajuda
./build/syntropy-linux --help

# Mostrar versão
./build/syntropy-linux --version

# Executar setup
./build/syntropy-linux setup run --force

# Verificar status
./build/syntropy-linux setup status
```

## 🔧 Opções Avançadas (Linux/WSL)

```bash
# Build apenas para Windows
./build-and-test.sh --windows

# Build apenas para Linux
./build-and-test.sh --linux

# Testar binários existentes
./build-and-test.sh --test

# Executar aplicação
./build-and-test.sh --run

# Mostrar ajuda
./build-and-test.sh --help
```

## 📊 O que o Script Faz

1. **Verifica pré-requisitos** - Go instalado e versão correta
2. **Prepara ambiente** - Cria diretório build e limpa builds anteriores
3. **Configura dependências** - Download, organização e verificação
4. **Compila para Windows** - Gera `syntropy-windows.exe`
5. **Compila para Linux** - Gera `syntropy-linux`
6. **Testa binários** - Verifica se executam corretamente
7. **Mostra resumo** - Informações dos binários criados
8. **Executa aplicação** - (opcional) Roda a aplicação para teste

## 🛠️ Solução de Problemas

### Erro: "Go não encontrado"
```bash
# Verificar se Go está instalado
go version

# Se não estiver, instalar Go
# Windows: https://golang.org/dl/
# Linux: sudo apt install golang-go
```

### Erro: "main.go não encontrado"
```bash
# Verificar se está no diretório correto
pwd
# Deve estar em: .../manager/interfaces/cli/

# Navegar para o diretório correto
cd /caminho/para/syntropy-cooperative-grid/manager/interfaces/cli
```

### Erro: "Falha na compilação"
```bash
# Limpar e tentar novamente
rm -rf build/
./build-and-test.sh

# Ou no Windows
del /q build\*
build-and-test.bat
```

## 📦 Distribuição

### Para Windows
- Copie `build/syntropy-windows.exe` para a máquina Windows
- Execute: `syntropy-windows.exe --help`

### Para Linux
- Copie `build/syntropy-linux` para a máquina Linux
- Torne executável: `chmod +x syntropy-linux`
- Execute: `./syntropy-linux --help`

## 🎯 Comandos Úteis da CLI

```bash
# Ajuda geral
./syntropy-linux --help

# Ajuda do setup
./syntropy-linux setup --help

# Validar ambiente (sem fazer mudanças)
./syntropy-linux setup validate

# Executar setup
./syntropy-linux setup run --force

# Verificar status
./syntropy-linux setup status

# Reset configuração
./syntropy-linux setup reset --force
```

## 📈 Informações do Build

O script inclui automaticamente:
- **Versão**: Timestamp da compilação
- **Git Commit**: Hash do commit atual (se disponível)
- **Build Time**: Data e hora da compilação
- **Plataforma**: Sistema operacional e arquitetura

## 🚀 Próximos Passos

1. **Execute** o script apropriado para seu sistema
2. **Teste** os binários gerados
3. **Copie** os binários para as máquinas de destino
4. **Execute** a aplicação e teste os comandos
5. **Configure** o ambiente com `setup run --force`

---

**Scripts criados:**
- `build-and-test.sh` - Para Linux/WSL
- `build-and-test.bat` - Para Windows

**Simples, direto e funcional!** 🎉
