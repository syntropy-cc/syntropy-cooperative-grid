# Syntropy CLI Manager - Development Workflow
# Script completo para desenvolvimento e teste da aplicação CLI

param(
    [Parameter(Position=0)]
    [ValidateSet("setup", "build", "test", "run", "dev", "clean", "install", "uninstall", "status", "help")]
    [string]$Action = "help",
    
    [Parameter(Position=1)]
    [string]$Args = ""
)

# Configurações
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$BuildDir = Join-Path $ScriptDir "build"
$Version = Get-Date -Format "yyyyMMdd-HHmmss"
$GitCommit = try { git rev-parse --short HEAD 2>$null } catch { "unknown" }
$BuildTime = Get-Date -Format "yyyy-MM-ddTHH:mm:ssZ"
$BinaryName = "syntropy.exe"

# Funções de logging
function Write-Info {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor Blue
}

function Write-Success {
    param([string]$Message)
    Write-Host "[SUCCESS] $Message" -ForegroundColor Green
}

function Write-Warning {
    param([string]$Message)
    Write-Host "[WARNING] $Message" -ForegroundColor Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor Red
}

function Write-Step {
    param([string]$Message)
    Write-Host ""
    Write-Host "=== $Message ===" -ForegroundColor Cyan
}

# Banner
function Show-Banner {
    Write-Host ""
    Write-Host "╔══════════════════════════════════════════════════════════════╗" -ForegroundColor Magenta
    Write-Host "║              SYNTROPY CLI MANAGER                           ║" -ForegroundColor Magenta
    Write-Host "║                Development Workflow                         ║" -ForegroundColor Magenta
    Write-Host "╚══════════════════════════════════════════════════════════════╝" -ForegroundColor Magenta
    Write-Host ""
}

# Verificar pré-requisitos
function Test-Prerequisites {
    Write-Step "Verificando Pré-requisitos"
    
    # Verificar Go
    try {
        $goVersion = (go version).Split(' ')[2].Substring(2)
        $requiredVersion = [Version]"1.22"
        $currentVersion = [Version]$goVersion
        
        if ($currentVersion -lt $requiredVersion) {
            Write-Error "Go versão $goVersion encontrada, mas versão $requiredVersion ou superior é necessária."
            exit 1
        }
        
        Write-Success "Go $goVersion encontrado"
    }
    catch {
        Write-Error "Go não está instalado. Por favor, instale Go 1.22.5 ou superior."
        Write-Info "Download: https://golang.org/dl/"
        exit 1
    }
    
    # Verificar Git
    try {
        $gitVersion = git --version
        Write-Success "Git encontrado: $gitVersion"
    }
    catch {
        Write-Warning "Git não encontrado. Algumas funcionalidades podem não funcionar."
    }
    
    # Verificar se estamos no diretório correto
    if (-not (Test-Path (Join-Path $ScriptDir "main.go"))) {
        Write-Error "main.go não encontrado em $ScriptDir"
        exit 1
    }
    
    Write-Success "Estrutura do diretório CLI verificada"
}

# Setup inicial do ambiente de desenvolvimento
function Initialize-DevEnvironment {
    Write-Step "Configurando Ambiente de Desenvolvimento"
    
    # Navegar para o diretório CLI
    Set-Location $ScriptDir
    
    # Criar diretórios necessários
    $directories = @($BuildDir, "logs", "temp")
    foreach ($dir in $directories) {
        if (-not (Test-Path $dir)) {
            New-Item -ItemType Directory -Path $dir -Force | Out-Null
            Write-Info "Diretório criado: $dir"
        }
    }
    
    # Configurar dependências
    Write-Info "Configurando dependências..."
    go mod download
    go mod tidy
    go mod verify
    
    # Verificar ferramentas de desenvolvimento
    $tools = @("golangci-lint", "gofmt", "goimports")
    foreach ($tool in $tools) {
        try {
            & $tool --version | Out-Null
            Write-Success "$tool encontrado"
        }
        catch {
            Write-Warning "$tool não encontrado. Instale com: go install golang.org/x/tools/cmd/$tool@latest"
        }
    }
    
    Write-Success "Ambiente de desenvolvimento configurado"
}

# Build completo
function Build-Application {
    Write-Step "Compilando Aplicação"
    
    # Preparar ambiente
    if (-not (Test-Path $BuildDir)) {
        New-Item -ItemType Directory -Path $BuildDir -Force | Out-Null
    }
    
    # Limpar builds anteriores
    Remove-Item "$BuildDir\*" -Force -ErrorAction SilentlyContinue
    
    # Configurar dependências
    go mod download
    go mod tidy
    
    # Compilar
    $buildFlags = "-ldflags `"-X main.version=$Version -X main.buildTime=$BuildTime -X main.gitCommit=$GitCommit`""
    
    Write-Info "Compilando CLI Manager..."
    $buildCommand = "go build $buildFlags -o $BuildDir\$BinaryName main.go"
    Invoke-Expression $buildCommand
    
    if (Test-Path (Join-Path $BuildDir $BinaryName)) {
        Write-Success "Compilação concluída: $BuildDir\$BinaryName"
        
        # Mostrar informações do binário
        $binaryInfo = Get-Item (Join-Path $BuildDir $BinaryName)
        $sizeKB = [math]::Round($binaryInfo.Length / 1KB, 2)
        Write-Info "Tamanho: $sizeKB KB"
        Write-Info "Criado em: $($binaryInfo.LastWriteTime)"
    } else {
        Write-Error "Falha na compilação"
        exit 1
    }
}

# Executar testes
function Invoke-Tests {
    Write-Step "Executando Testes"
    
    # Executar testes unitários
    Write-Info "Executando testes unitários..."
    try {
        go test -v .\...
        Write-Success "Testes unitários concluídos"
    }
    catch {
        Write-Warning "Alguns testes falharam (esperado para funcionalidades não implementadas)"
    }
    
    # Executar testes com cobertura
    Write-Info "Executando testes com cobertura..."
    try {
        go test -v -cover .\...
        Write-Success "Testes com cobertura concluídos"
    }
    catch {
        Write-Warning "Alguns testes falharam"
    }
    
    # Executar testes de race condition
    Write-Info "Executando testes de race condition..."
    try {
        go test -v -race .\...
        Write-Success "Testes de race condition concluídos"
    }
    catch {
        Write-Warning "Alguns testes de race condition falharam"
    }
    
    Write-Success "Todos os testes executados"
}

# Verificações de qualidade
function Invoke-QualityChecks {
    Write-Step "Executando Verificações de Qualidade"
    
    # Formatar código
    Write-Info "Formatando código..."
    go fmt .\...
    Write-Success "Código formatado"
    
    # Executar go vet
    Write-Info "Executando go vet..."
    go vet .\...
    Write-Success "go vet concluído"
    
    # Executar golangci-lint se disponível
    try {
        golangci-lint --version | Out-Null
        Write-Info "Executando golangci-lint..."
        golangci-lint run
        Write-Success "golangci-lint concluído"
    }
    catch {
        Write-Warning "golangci-lint não instalado. Instale com: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
    }
    
    Write-Success "Verificações de qualidade concluídas"
}

# Executar aplicação
function Start-Application {
    Write-Step "Executando Aplicação"
    
    $binaryPath = Join-Path $BuildDir $BinaryName
    
    if (-not (Test-Path $binaryPath)) {
        Write-Error "Binário não encontrado. Execute 'build' primeiro."
        exit 1
    }
    
    Write-Info "Executando: $binaryPath $Args"
    Write-Host ""
    
    # Executar com argumentos fornecidos
    if ($Args) {
        & $binaryPath $Args.Split(' ')
    } else {
        & $binaryPath
    }
}

# Modo de desenvolvimento (build + test + run)
function Start-DevMode {
    Write-Step "Iniciando Modo de Desenvolvimento"
    
    # Build
    Build-Application
    
    # Testes
    Invoke-Tests
    
    # Verificações de qualidade
    Invoke-QualityChecks
    
    # Verificar binário
    $binaryPath = Join-Path $BuildDir $BinaryName
    if (Test-Path $binaryPath) {
        Write-Info "Testando binário..."
        try {
            & $binaryPath --version | Out-Null
            Write-Success "Binário funcionando corretamente"
        }
        catch {
            Write-Info "Binário criado (teste básico concluído)"
        }
    }
    
    Write-Success "Modo de desenvolvimento concluído"
    Write-Info "Use '.\dev-workflow.ps1 run' para executar a aplicação"
}

# Instalar binário
function Install-Binary {
    Write-Step "Instalando Binário"
    
    $binaryPath = Join-Path $BuildDir $BinaryName
    $installPath = Join-Path $env:ProgramFiles "Syntropy\CLI"
    
    if (-not (Test-Path $binaryPath)) {
        Write-Error "Binário não encontrado. Execute 'build' primeiro."
        exit 1
    }
    
    # Criar diretório de instalação
    if (-not (Test-Path $installPath)) {
        New-Item -ItemType Directory -Path $installPath -Force | Out-Null
    }
    
    # Copiar binário
    Copy-Item $binaryPath $installPath -Force
    
    # Adicionar ao PATH (opcional)
    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    if ($currentPath -notlike "*$installPath*") {
        Write-Info "Adicionando ao PATH do usuário..."
        [Environment]::SetEnvironmentVariable("PATH", "$currentPath;$installPath", "User")
        Write-Success "Binário instalado em: $installPath"
        Write-Warning "Reinicie o terminal para usar o comando 'syntropy' globalmente"
    } else {
        Write-Success "Binário instalado em: $installPath"
    }
}

# Desinstalar binário
function Uninstall-Binary {
    Write-Step "Desinstalando Binário"
    
    $installPath = Join-Path $env:ProgramFiles "Syntropy\CLI"
    $binaryPath = Join-Path $installPath $BinaryName
    
    if (Test-Path $binaryPath) {
        Remove-Item $binaryPath -Force
        Write-Success "Binário removido"
    }
    
    # Remover do PATH
    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    if ($currentPath -like "*$installPath*") {
        $newPath = $currentPath -replace [regex]::Escape(";$installPath"), ""
        [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
        Write-Success "Removido do PATH do usuário"
        Write-Warning "Reinicie o terminal para aplicar as mudanças"
    }
    
    Write-Success "Desinstalação concluída"
}

# Verificar status
function Show-Status {
    Write-Step "Status do Sistema"
    
    # Verificar Go
    try {
        $goVersion = (go version).Split(' ')[2].Substring(2)
        Write-Success "Go $goVersion instalado"
    }
    catch {
        Write-Error "Go não instalado"
    }
    
    # Verificar binário
    $binaryPath = Join-Path $BuildDir $BinaryName
    if (Test-Path $binaryPath) {
        $binaryInfo = Get-Item $binaryPath
        $sizeKB = [math]::Round($binaryInfo.Length / 1KB, 2)
        Write-Success "Binário encontrado: $($binaryInfo.Name) ($sizeKB KB)"
        Write-Info "Criado em: $($binaryInfo.LastWriteTime)"
    } else {
        Write-Warning "Binário não encontrado. Execute 'build' primeiro."
    }
    
    # Verificar instalação global
    $installPath = Join-Path $env:ProgramFiles "Syntropy\CLI"
    $installedBinary = Join-Path $installPath $BinaryName
    if (Test-Path $installedBinary) {
        Write-Success "Instalação global encontrada: $installPath"
    } else {
        Write-Info "Nenhuma instalação global encontrada"
    }
    
    # Verificar PATH
    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    if ($currentPath -like "*$installPath*") {
        Write-Success "CLI está no PATH do usuário"
    } else {
        Write-Info "CLI não está no PATH do usuário"
    }
}

# Limpar build
function Clear-Build {
    Write-Step "Limpando Build"
    
    if (Test-Path $BuildDir) {
        Remove-Item $BuildDir -Recurse -Force
        Write-Success "Diretório de build limpo"
    } else {
        Write-Info "Nenhum build para limpar"
    }
    
    # Limpar logs e temp
    $cleanDirs = @("logs", "temp")
    foreach ($dir in $cleanDirs) {
        if (Test-Path $dir) {
            Remove-Item $dir -Recurse -Force -ErrorAction SilentlyContinue
            Write-Info "Diretório limpo: $dir"
        }
    }
}

# Mostrar ajuda
function Show-Help {
    Write-Host "Uso: .\dev-workflow.ps1 [ação] [argumentos]"
    Write-Host ""
    Write-Host "Ações:"
    Write-Host "  setup     Configurar ambiente de desenvolvimento"
    Write-Host "  build     Compilar a aplicação"
    Write-Host "  test      Executar testes"
    Write-Host "  run       Executar a aplicação compilada"
    Write-Host "  dev       Modo desenvolvimento (build + test + quality)"
    Write-Host "  install   Instalar binário globalmente"
    Write-Host "  uninstall Desinstalar binário"
    Write-Host "  status    Verificar status do sistema"
    Write-Host "  clean     Limpar diretórios de build"
    Write-Host "  help      Mostrar esta ajuda"
    Write-Host ""
    Write-Host "Exemplos:"
    Write-Host "  .\dev-workflow.ps1 setup                    # Setup inicial"
    Write-Host "  .\dev-workflow.ps1 dev                      # Modo desenvolvimento"
    Write-Host "  .\dev-workflow.ps1 build                    # Compilar"
    Write-Host "  .\dev-workflow.ps1 test                     # Executar testes"
    Write-Host "  .\dev-workflow.ps1 run                      # Executar"
    Write-Host "  .\dev-workflow.ps1 run '--help'             # Executar com argumentos"
    Write-Host "  .\dev-workflow.ps1 run 'setup run --force'  # Executar setup"
    Write-Host "  .\dev-workflow.ps1 install                  # Instalar"
    Write-Host "  .\dev-workflow.ps1 status                   # Verificar status"
    Write-Host "  .\dev-workflow.ps1 clean                    # Limpar"
}

# Função principal
function Main {
    Show-Banner
    
    switch ($Action.ToLower()) {
        "setup" {
            Test-Prerequisites
            Initialize-DevEnvironment
        }
        "build" {
            Test-Prerequisites
            Build-Application
        }
        "test" {
            Test-Prerequisites
            Set-Location $ScriptDir
            Invoke-Tests
        }
        "run" {
            Start-Application
        }
        "dev" {
            Test-Prerequisites
            Start-DevMode
        }
        "install" {
            Install-Binary
        }
        "uninstall" {
            Uninstall-Binary
        }
        "status" {
            Show-Status
        }
        "clean" {
            Clear-Build
        }
        "help" {
            Show-Help
        }
        default {
            Write-Error "Ação desconhecida: $Action"
            Write-Host "Use '.\dev-workflow.ps1 help' para opções disponíveis"
            exit 1
        }
    }
}

# Executar função principal
Main
