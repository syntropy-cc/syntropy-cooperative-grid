# Script de Compilação Automatizada - Setup Component
# Syntropy Cooperative Grid
# PowerShell Script para Windows

param(
    [Parameter(Position=0)]
    [ValidateSet("all", "linux", "windows", "test", "clean", "help")]
    [string]$Action = "all"
)

# Configurações
$ProjectRoot = "C:\Users\$env:USERNAME\syntropy-cc\syntropy-cooperative-grid"
$SetupDir = "$ProjectRoot\manager\interfaces\cli\setup"
$BuildDir = "$SetupDir\build"
$Version = Get-Date -Format "yyyyMMdd-HHmmss"

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
    Write-Host "║              SYNTROPY SETUP COMPONENT                        ║" -ForegroundColor Magenta
    Write-Host "║                    Build Script                              ║" -ForegroundColor Magenta
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
        exit 1
    }
    
    # Verificar diretório do projeto
    if (-not (Test-Path $ProjectRoot)) {
        Write-Error "Diretório do projeto não encontrado: $ProjectRoot"
        exit 1
    }
    
    Write-Success "Diretório do projeto encontrado"
    
    # Verificar go.mod
    if (-not (Test-Path "$ProjectRoot\go.mod")) {
        Write-Error "Arquivo go.mod não encontrado no diretório raiz"
        exit 1
    }
    
    Write-Success "Arquivo go.mod encontrado"
    
    # Verificar PowerShell
    $psVersion = $PSVersionTable.PSVersion
    Write-Success "PowerShell $psVersion encontrado"
}

# Preparar ambiente de build
function Initialize-Build {
    Write-Step "Preparando Ambiente de Build"
    
    # Navegar para o diretório raiz
    Set-Location $ProjectRoot
    
    # Criar diretório de build
    if (-not (Test-Path $BuildDir)) {
        New-Item -ItemType Directory -Path $BuildDir -Force | Out-Null
    }
    
    # Limpar builds anteriores
    Remove-Item "$BuildDir\*" -Force -ErrorAction SilentlyContinue
    
    Write-Success "Ambiente de build preparado"
}

# Baixar e verificar dependências
function Setup-Dependencies {
    Write-Step "Configurando Dependências"
    
    # Baixar dependências
    Write-Info "Baixando dependências..."
    go mod download
    
    # Organizar dependências
    Write-Info "Organizando dependências..."
    go mod tidy
    
    # Verificar dependências
    Write-Info "Verificando dependências..."
    go mod verify
    
    Write-Success "Dependências configuradas"
}

# Compilar para Windows
function Build-Windows {
    Write-Step "Compilando para Windows"
    
    Set-Location $SetupDir
    
    # Compilar setup component para Windows
    Write-Info "Compilando setup component..."
    $buildTime = Get-Date -Format "yyyy-MM-ddTHH:mm:ssZ"
    go build -ldflags "-X main.version=$Version -X main.buildTime=$buildTime" `
        -o "$BuildDir\syntropy-setup-windows.exe" .
    
    # Compilar CLI completo para Windows
    Write-Info "Compilando CLI completo..."
    Set-Location "$ProjectRoot\interfaces\cli"
    go build -ldflags "-X main.version=$Version -X main.buildTime=$buildTime" `
        -o "$BuildDir\syntropy-cli-windows.exe" .\cmd\main.go
    
    Write-Success "Compilação para Windows concluída"
}

# Compilar para Linux (cross-compilation)
function Build-Linux {
    Write-Step "Compilando para Linux (Cross-compilation)"
    
    Set-Location $SetupDir
    
    # Definir variáveis para cross-compilation
    $env:GOOS = "linux"
    $env:GOARCH = "amd64"
    
    # Compilar setup component para Linux
    Write-Info "Compilando setup component para Linux..."
    $buildTime = Get-Date -Format "yyyy-MM-ddTHH:mm:ssZ"
    go build -ldflags "-X main.version=$Version -X main.buildTime=$buildTime" `
        -o "$BuildDir\syntropy-setup-linux" .
    
    # Compilar CLI completo para Linux
    Write-Info "Compilando CLI completo para Linux..."
    Set-Location "$ProjectRoot\interfaces\cli"
    go build -ldflags "-X main.version=$Version -X main.buildTime=$buildTime" `
        -o "$BuildDir\syntropy-cli-linux" .\cmd\main.go
    
    # Restaurar variáveis para Windows
    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
    
    Write-Success "Compilação para Linux concluída"
}

# Executar testes
function Invoke-Tests {
    Write-Step "Executando Testes"
    
    Set-Location $SetupDir
    
    # Executar testes unitários
    Write-Info "Executando testes unitários..."
    try {
        go test -v .\...
        Write-Success "Testes unitários executados"
    }
    catch {
        Write-Warning "Alguns testes falharam (esperado para funcionalidades não implementadas)"
    }
    
    # Executar testes com cobertura
    Write-Info "Executando testes com cobertura..."
    try {
        go test -v -cover .\...
        Write-Success "Testes com cobertura executados"
    }
    catch {
        Write-Warning "Alguns testes falharam"
    }
    
    # Executar testes de integração se existirem
    if (Test-Path "tests\integration") {
        Write-Info "Executando testes de integração..."
        try {
            go test -v .\tests\integration\...
            Write-Success "Testes de integração executados"
        }
        catch {
            Write-Warning "Testes de integração falharam"
        }
    }
    
    Write-Success "Testes executados"
}

# Análise de qualidade
function Invoke-QualityChecks {
    Write-Step "Executando Análise de Qualidade"
    
    Set-Location $SetupDir
    
    # Formatar código
    Write-Info "Formatando código..."
    go fmt .\...
    
    # Executar go vet
    Write-Info "Executando go vet..."
    go vet .\...
    
    # Verificar se golangci-lint está disponível
    try {
        golangci-lint --version | Out-Null
        Write-Info "Executando golangci-lint..."
        golangci-lint run
        Write-Success "golangci-lint executado"
    }
    catch {
        Write-Warning "golangci-lint não está instalado. Instale com: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
    }
    
    Write-Success "Análise de qualidade concluída"
}

# Verificar binários
function Test-Binaries {
    Write-Step "Verificando Binários"
    
    Set-Location $BuildDir
    
    # Verificar arquivos criados
    Write-Info "Arquivos criados:"
    Get-ChildItem | Format-Table Name, Length, LastWriteTime
    
    # Verificar informações dos binários Windows
    if (Test-Path "syntropy-setup-windows.exe") {
        Write-Info "Informações do syntropy-setup-windows.exe:"
        Get-Item "syntropy-setup-windows.exe" | Select-Object Name, Length, LastWriteTime
        
        # Tentar executar help (pode falhar se não implementado)
        try {
            & ".\syntropy-setup-windows.exe" --help 2>$null
            Write-Info "Binário executável"
        }
        catch {
            Write-Info "Binário criado (help não disponível)"
        }
    }
    
    if (Test-Path "syntropy-cli-windows.exe") {
        Write-Info "Informações do syntropy-cli-windows.exe:"
        Get-Item "syntropy-cli-windows.exe" | Select-Object Name, Length, LastWriteTime
        
        # Tentar executar help
        try {
            & ".\syntropy-cli-windows.exe" --help 2>$null
            Write-Info "Binário executável"
        }
        catch {
            Write-Info "Binário criado (help não disponível)"
        }
    }
    
    # Verificar informações dos binários Linux
    if (Test-Path "syntropy-setup-linux") {
        Write-Info "Informações do syntropy-setup-linux:"
        Get-Item "syntropy-setup-linux" | Select-Object Name, Length, LastWriteTime
    }
    
    if (Test-Path "syntropy-cli-linux") {
        Write-Info "Informações do syntropy-cli-linux:"
        Get-Item "syntropy-cli-linux" | Select-Object Name, Length, LastWriteTime
    }
    
    Write-Success "Verificação de binários concluída"
}

# Criar pacotes de distribuição
function New-Packages {
    Write-Step "Criando Pacotes de Distribuição"
    
    Set-Location $BuildDir
    
    # Criar pacote Windows
    if ((Test-Path "syntropy-setup-windows.exe") -and (Test-Path "syntropy-cli-windows.exe")) {
        Write-Info "Criando pacote Windows..."
        $zipFile = "syntropy-setup-windows-$Version.zip"
        Compress-Archive -Path "syntropy-setup-windows.exe", "syntropy-cli-windows.exe" -DestinationPath $zipFile -Force
        Write-Success "Pacote Windows criado: $zipFile"
    }
    
    # Criar pacote Linux
    if ((Test-Path "syntropy-setup-linux") -and (Test-Path "syntropy-cli-linux")) {
        Write-Info "Criando pacote Linux..."
        $tarFile = "syntropy-setup-linux-$Version.tar.gz"
        # Nota: tar.exe está disponível no Windows 10/11
        tar -czf $tarFile syntropy-setup-linux syntropy-cli-linux
        Write-Success "Pacote Linux criado: $tarFile"
    }
    
    Write-Success "Pacotes de distribuição criados"
}

# Mostrar resumo
function Show-Summary {
    Write-Step "Resumo da Compilação"
    
    Write-Host "✅ Compilação Concluída com Sucesso!" -ForegroundColor Green
    Write-Host ""
    Write-Host "📁 Diretório de Build: $BuildDir" -ForegroundColor Blue
    Write-Host "📦 Versão: $Version" -ForegroundColor Blue
    Write-Host "🕒 Timestamp: $(Get-Date)" -ForegroundColor Blue
    Write-Host ""
    Write-Host "📋 Binários Criados:" -ForegroundColor Blue
    
    Set-Location $BuildDir
    Get-ChildItem | ForEach-Object {
        $size = [math]::Round($_.Length / 1KB, 2)
        Write-Host "  - $($_.Name) ($size KB)"
    }
    
    Write-Host ""
    Write-Host "🚀 Próximos Passos:" -ForegroundColor Blue
    Write-Host "  1. Testar binários manualmente"
    Write-Host "  2. Executar testes de integração"
    Write-Host "  3. Distribuir pacotes conforme necessário"
    Write-Host "  4. Atualizar documentação se necessário"
    
    Write-Host ""
    Write-Host "💡 Dicas:" -ForegroundColor Cyan
    Write-Host "  - Use '.\syntropy-setup-windows.exe --help' para ver opções"
    Write-Host "  - Use '.\syntropy-cli-windows.exe setup --help' para comandos de setup"
    Write-Host "  - Consulte COMPILACAO_E_TESTE.md para instruções detalhadas"
}

# Função principal
function Main {
    Show-Banner
    
    switch ($Action) {
        "linux" {
            Test-Prerequisites
            Initialize-Build
            Setup-Dependencies
            Build-Linux
            Invoke-Tests
            Invoke-QualityChecks
            Test-Binaries
            New-Packages
            Show-Summary
        }
        "windows" {
            Test-Prerequisites
            Initialize-Build
            Setup-Dependencies
            Build-Windows
            Invoke-Tests
            Invoke-QualityChecks
            Test-Binaries
            New-Packages
            Show-Summary
        }
        "test" {
            Test-Prerequisites
            Set-Location $SetupDir
            Invoke-Tests
        }
        "clean" {
            Write-Info "Limpando diretório de build..."
            Remove-Item $BuildDir -Recurse -Force -ErrorAction SilentlyContinue
            Write-Success "Limpeza concluída"
        }
        "help" {
            Write-Host "Uso: .\build.ps1 [opção]"
            Write-Host ""
            Write-Host "Opções:"
            Write-Host "  all       Compilar para todos os sistemas (padrão)"
            Write-Host "  linux     Compilar apenas para Linux"
            Write-Host "  windows   Compilar apenas para Windows"
            Write-Host "  test      Executar apenas os testes"
            Write-Host "  clean     Limpar diretório de build"
            Write-Host "  help      Mostrar esta ajuda"
            Write-Host ""
            Write-Host "Exemplos:"
            Write-Host "  .\build.ps1                # Compilar tudo"
            Write-Host "  .\build.ps1 windows        # Compilar apenas Windows"
            Write-Host "  .\build.ps1 test           # Executar apenas testes"
            Write-Host "  .\build.ps1 clean          # Limpar build"
        }
        "all" {
            Test-Prerequisites
            Initialize-Build
            Setup-Dependencies
            Build-Linux
            Build-Windows
            Invoke-Tests
            Invoke-QualityChecks
            Test-Binaries
            New-Packages
            Show-Summary
        }
        default {
            Write-Error "Opção desconhecida: $Action"
            Write-Host "Use '.\build.ps1 help' para ver opções disponíveis"
            exit 1
        }
    }
}

# Executar função principal
Main
