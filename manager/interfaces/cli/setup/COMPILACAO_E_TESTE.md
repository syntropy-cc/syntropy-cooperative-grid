# Instruções de Compilação e Teste - Setup Component

Este documento fornece instruções passo-a-passo para compilar e testar o código do Setup Component do Syntropy CLI nos sistemas Linux e Windows.

## 📋 Visão Geral

O Setup Component é um componente Go que fornece funcionalidades de configuração para o Syntropy CLI em diferentes sistemas operacionais. Ele está localizado em `manager/interfaces/cli/setup/` e inclui:

- **Arquivo principal**: `setup.go` (orquestrador)
- **Implementações específicas**: `setup_linux.go`, `setup_windows.go`
- **Validação**: `validation_linux.go`, `validation_windows.go`
- **Configuração**: `configuration_linux.go`, `configuration_windows.go`
- **Testes**: `setup_test.go`, `tests/` (unitários e integração)

## 🔧 Pré-requisitos

### Dependências Comuns
- **Go 1.22.5+** (conforme `go.mod`)
- **Git** para controle de versão
- **Make** (opcional, mas recomendado)

### Dependências Específicas por SO

#### Linux
- `systemd` (para serviços)
- `systemctl` (para gerenciamento de serviços)
- Permissões de administrador (para instalação de serviços)

#### Windows
- **PowerShell 5.1+**
- Permissões de administrador
- **Windows Service Control Manager**

## 🏗️ Estrutura do Projeto

```
manager/interfaces/cli/setup/
├── setup.go                     # Orquestrador principal
├── setup_linux.go               # Implementação Linux
├── setup_windows.go             # Implementação Windows
├── validation_linux.go          # Validação Linux
├── validation_windows.go        # Validação Windows
├── configuration_linux.go       # Configuração Linux
├── configuration_windows.go     # Configuração Windows
├── internal/                    # Tipos e serviços internos
│   ├── types/                   # Estruturas de dados
│   ├── services/                # Serviços internos
│   └── utils/                   # Utilitários
├── tests/                       # Testes
│   ├── unit/                    # Testes unitários
│   ├── integration/             # Testes de integração
│   └── fixtures/                # Dados de teste
├── config/                      # Configurações e templates
└── COMPILACAO_E_TESTE.md        # Este documento
```

## 🐧 Compilação no Linux

### Passo 1: Preparar o Ambiente

```bash
# Navegar para o diretório do projeto
cd /home/jescott/syntropy-cc/syntropy-cooperative-grid

# Verificar versão do Go
go version
# Deve mostrar Go 1.22.5 ou superior

# Verificar se estamos no diretório correto
pwd
# Deve mostrar: /home/jescott/syntropy-cc/syntropy-cooperative-grid
```

### Passo 2: Instalar Dependências

```bash
# Baixar dependências do Go
go mod download

# Verificar e organizar dependências
go mod tidy

# Verificar se não há problemas de dependência
go mod verify
```

### Passo 3: Compilar o Setup Component

#### Opção A: Compilação Direta (Recomendada)

```bash
# Navegar para o diretório setup
cd manager/interfaces/cli/setup

# Compilar apenas o setup component
go build -o syntropy-setup-linux ./setup.go

# Verificar se o binário foi criado
ls -la syntropy-setup-linux
file syntropy-setup-linux
```

#### Opção B: Compilação com Build Tags

```bash
# Compilar especificamente para Linux
go build -tags linux -o syntropy-setup-linux .

# Compilar com informações de debug
go build -ldflags "-X main.version=dev -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o syntropy-setup-linux .
```

#### Opção C: Compilação do CLI Completo

```bash
# Voltar ao diretório raiz
cd /home/jescott/syntropy-cc/syntropy-cooperative-grid

# Compilar o CLI completo (inclui setup component)
go build -o syntropy-cli-linux ./interfaces/cli/cmd/main.go
```

### Passo 4: Usar o Makefile (Alternativo)

```bash
# Navegar para o diretório CLI
cd /home/jescott/syntropy-cc/syntropy-cooperative-grid/interfaces/cli

# Verificar o Makefile
cat Makefile

# Compilar usando Make
make build

# Executar testes
make test

# Limpar arquivos de build
make clean
```

## 🪟 Compilação no Windows

### Passo 1: Preparar o Ambiente

```powershell
# Abrir PowerShell como Administrador
# Navegar para o diretório do projeto
cd C:\Users\%USERNAME%\syntropy-cc\syntropy-cooperative-grid

# Verificar versão do Go
go version
# Deve mostrar Go 1.22.5 ou superior

# Verificar se estamos no diretório correto
Get-Location
# Deve mostrar o caminho do projeto
```

### Passo 2: Instalar Dependências

```powershell
# Baixar dependências do Go
go mod download

# Verificar e organizar dependências
go mod tidy

# Verificar se não há problemas de dependência
go mod verify
```

### Passo 3: Compilar o Setup Component

#### Opção A: Compilação Direta

```powershell
# Navegar para o diretório setup
cd manager\interfaces\cli\setup

# Compilar apenas o setup component
go build -o syntropy-setup-windows.exe ./setup.go

# Verificar se o executável foi criado
dir syntropy-setup-windows.exe
```

#### Opção B: Compilação com Build Tags

```powershell
# Compilar especificamente para Windows
go build -tags windows -o syntropy-setup-windows.exe .

# Compilar com informações de debug
$buildTime = Get-Date -Format "yyyy-MM-ddTHH:mm:ssZ"
go build -ldflags "-X main.version=dev -X main.buildTime=$buildTime" -o syntropy-setup-windows.exe .
```

#### Opção C: Compilação do CLI Completo

```powershell
# Voltar ao diretório raiz
cd C:\Users\%USERNAME%\syntropy-cc\syntropy-cooperative-grid

# Compilar o CLI completo
go build -o syntropy-cli-windows.exe .\interfaces\cli\cmd\main.go
```

### Passo 4: Compilação Cross-Platform (do Linux para Windows)

```bash
# No Linux, compilar para Windows
cd /home/jescott/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup

# Definir variáveis de ambiente para cross-compilation
export GOOS=windows
export GOARCH=amd64

# Compilar para Windows
go build -o syntropy-setup-windows.exe ./setup.go

# Restaurar variáveis para Linux
export GOOS=linux
export GOARCH=amd64
```

## 🧪 Testes

### Executar Testes Unitários

#### Linux
```bash
# Navegar para o diretório setup
cd /home/jescott/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup

# Executar todos os testes
go test -v ./...

# Executar testes com cobertura
go test -v -cover ./...

# Executar testes específicos
go test -v -run TestSetupFlow ./setup_test.go

# Executar testes com race detection
go test -v -race ./...
```

#### Windows
```powershell
# Navegar para o diretório setup
cd manager\interfaces\cli\setup

# Executar todos os testes
go test -v .\...

# Executar testes com cobertura
go test -v -cover .\...

# Executar testes específicos
go test -v -run TestSetupFlow .\setup_test.go

# Executar testes com race detection
go test -v -race .\...
```

### Executar Testes de Integração

```bash
# Navegar para o diretório de testes
cd /home/jescott/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/tests

# Executar testes de integração
go test -v ./integration/...

# Executar testes com dados de teste específicos
go test -v ./integration/ -testdata=./fixtures/
```

### Executar Testes com Diferentes Build Tags

```bash
# Testes apenas para Linux
go test -tags linux -v ./...

# Testes apenas para Windows (no Linux)
GOOS=windows go test -tags windows -v ./...

# Testes para ambos os sistemas
go test -tags "linux windows" -v ./...
```

## 🚀 Execução e Teste Manual

### Testar o Setup Component

#### Linux
```bash
# Navegar para o diretório com o binário
cd /home/jescott/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup

# Executar setup (se compilado como binário standalone)
./syntropy-setup-linux --help

# Testar funcionalidade de setup
./syntropy-setup-linux setup --validate-only

# Testar com diferentes opções
./syntropy-setup-linux setup --force --install-service
```

#### Windows
```powershell
# Navegar para o diretório com o executável
cd manager\interfaces\cli\setup

# Executar setup
.\syntropy-setup-windows.exe --help

# Testar funcionalidade de setup
.\syntropy-setup-windows.exe setup --validate-only

# Testar com diferentes opções
.\syntropy-setup-windows.exe setup --force --install-service
```

### Testar Integração com CLI Principal

#### Linux
```bash
# Navegar para o diretório CLI
cd /home/jescott/syntropy-cc/syntropy-cooperative-grid/interfaces/cli

# Compilar CLI completo
go build -o syntropy-cli ./cmd/main.go

# Testar comandos de setup
./syntropy-cli setup --help
./syntropy-cli setup status
./syntropy-cli setup --validate-only
```

#### Windows
```powershell
# Navegar para o diretório CLI
cd manager\interfaces\cli

# Compilar CLI completo
go build -o syntropy-cli.exe .\cmd\main.go

# Testar comandos de setup
.\syntropy-cli.exe setup --help
.\syntropy-cli.exe setup status
.\syntropy-cli.exe setup --validate-only
```

## 📊 Verificação de Qualidade

### Análise Estática do Código

```bash
# Instalar ferramentas de análise (se não estiver instalado)
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Formatar código
go fmt ./...
goimports -w .

# Executar linter
golangci-lint run

# Verificar com go vet
go vet ./...

# Verificar dependências vulneráveis
go list -json -deps ./... | nancy sleuth
```

### Verificação de Segurança

```bash
# Instalar ferramentas de segurança
go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

# Executar análise de segurança
gosec ./...

# Verificar dependências
go mod why <dependency-name>
```

## 🔍 Troubleshooting

### Problemas Comuns de Compilação

#### Erro: "package not found"
```bash
# Solução: Baixar dependências
go mod download
go mod tidy
```

#### Erro: "build constraints exclude all Go files"
```bash
# Solução: Verificar build tags
# Para Linux: go build -tags linux
# Para Windows: go build -tags windows
```

#### Erro: "permission denied" (Linux)
```bash
# Solução: Dar permissões de execução
chmod +x syntropy-setup-linux
```

#### Erro: "cannot find main package" (Windows)
```powershell
# Solução: Verificar se está no diretório correto
# O arquivo main.go deve estar em interfaces/cli/cmd/main.go
```

### Problemas de Teste

#### Testes falham com "not implemented"
```bash
# Isso é esperado em sistemas não-Windows
# Os stubs retornam ErrNotImplemented para funcionalidades não implementadas
```

#### Testes de integração falham
```bash
# Verificar se as dependências do sistema estão instaladas
# No Linux: systemd, systemctl
# No Windows: PowerShell, permissões de administrador
```

### Problemas de Execução

#### "command not found" no Linux
```bash
# Solução: Adicionar ao PATH ou usar caminho completo
export PATH=$PATH:/caminho/para/syntropy-setup-linux
# ou
./syntropy-setup-linux --help
```

#### "execution policy" no Windows
```powershell
# Solução: Alterar política de execução
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

## 📝 Exemplos de Uso

### Exemplo de Compilação Completa (Linux)

```bash
#!/bin/bash
# Script de compilação completa para Linux

set -e

echo "=== Compilação do Setup Component - Linux ==="

# Navegar para o diretório do projeto
cd /home/jescott/syntropy-cc/syntropy-cooperative-grid

# Verificar Go
echo "Verificando Go..."
go version

# Baixar dependências
echo "Baixando dependências..."
go mod download
go mod tidy

# Compilar setup component
echo "Compilando setup component..."
cd manager/interfaces/cli/setup
go build -ldflags "-X main.version=$(date +%Y%m%d-%H%M%S)" -o syntropy-setup-linux .

# Executar testes
echo "Executando testes..."
go test -v ./...

# Verificar binário
echo "Verificando binário..."
file syntropy-setup-linux
ls -la syntropy-setup-linux

echo "=== Compilação concluída com sucesso! ==="
```

### Exemplo de Compilação Completa (Windows)

```powershell
# Script de compilação completa para Windows

Write-Host "=== Compilação do Setup Component - Windows ===" -ForegroundColor Green

# Navegar para o diretório do projeto
Set-Location "C:\Users\$env:USERNAME\syntropy-cc\syntropy-cooperative-grid"

# Verificar Go
Write-Host "Verificando Go..."
go version

# Baixar dependências
Write-Host "Baixando dependências..."
go mod download
go mod tidy

# Compilar setup component
Write-Host "Compilando setup component..."
Set-Location "manager\interfaces\cli\setup"
$buildTime = Get-Date -Format "yyyyMMdd-HHmmss"
go build -ldflags "-X main.version=$buildTime" -o syntropy-setup-windows.exe .

# Executar testes
Write-Host "Executando testes..."
go test -v .\...

# Verificar executável
Write-Host "Verificando executável..."
Get-Item syntropy-setup-windows.exe

Write-Host "=== Compilação concluída com sucesso! ===" -ForegroundColor Green
```

## 📚 Recursos Adicionais

### Documentação Relacionada
- [README.md](./README.md) - Documentação do usuário
- [GUIDE.md](./GUIDE.md) - Guia de desenvolvimento
- [TODO.md](./TODO.md) - Lista de tarefas de implementação

### Comandos Úteis
```bash
# Limpar cache do Go
go clean -cache

# Verificar dependências desnecessárias
go mod why <package>

# Atualizar dependências
go get -u ./...

# Verificar vulnerabilidades
go list -json -deps ./... | nancy sleuth
```

### Links Úteis
- [Documentação oficial do Go](https://golang.org/doc/)
- [Go Build Constraints](https://pkg.go.dev/go/build#hdr-Build_Constraints)
- [Go Testing](https://golang.org/pkg/testing/)

---

## ✅ Checklist de Compilação

### Linux
- [ ] Go 1.22.5+ instalado
- [ ] Dependências baixadas (`go mod download`)
- [ ] Código compilado sem erros
- [ ] Testes executados com sucesso
- [ ] Binário criado e funcional
- [ ] Permissões de execução configuradas

### Windows
- [ ] Go 1.22.5+ instalado
- [ ] PowerShell 5.1+ disponível
- [ ] Dependências baixadas (`go mod download`)
- [ ] Código compilado sem erros
- [ ] Testes executados com sucesso
- [ ] Executável criado e funcional
- [ ] Política de execução configurada

### Ambos
- [ ] Linting executado sem erros
- [ ] Análise de segurança executada
- [ ] Documentação atualizada
- [ ] Testes de integração executados
- [ ] Funcionalidade testada manualmente

---

**Última atualização**: $(date)
**Versão**: 1.0
**Autor**: Equipe de Desenvolvimento Syntropy
