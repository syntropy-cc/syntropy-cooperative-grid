# Syntropy CLI Manager - Environment Setup Script
# Script para configurar o ambiente de desenvolvimento no Windows

param(
    [Parameter(Position=0)]
    [ValidateSet("check", "install", "configure", "verify", "help")]
    [string]$Action = "check"
)

# Configurações
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$CLIDir = Split-Path -Parent $ScriptDir

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
    Write-Host "║                Environment Setup                            ║" -ForegroundColor Magenta
    Write-Host "╚══════════════════════════════════════════════════════════════╝" -ForegroundColor Magenta
    Write-Host ""
}

# Verificar pré-requisitos
function Test-Prerequisites {
    Write-Step "Verificando Pré-requisitos"
    
    $results = @{
        Go = $false
        Git = $false
        PowerShell = $false
        Chocolatey = $false
        Winget = $false
    }
    
    # Verificar Go
    try {
        $goVersion = (go version).Split(' ')[2].Substring(2)
        $requiredVersion = [Version]"1.22"
        $currentVersion = [Version]$goVersion
        
        if ($currentVersion -ge $requiredVersion) {
            Write-Success "Go $goVersion encontrado"
            $results.Go = $true
        } else {
            Write-Warning "Go $goVersion encontrado, mas versão $requiredVersion ou superior é necessária"
        }
    }
    catch {
        Write-Error "Go não está instalado"
    }
    
    # Verificar Git
    try {
        $gitVersion = git --version
        Write-Success "Git encontrado: $gitVersion"
        $results.Git = $true
    }
    catch {
        Write-Error "Git não está instalado"
    }
    
    # Verificar PowerShell
    $psVersion = $PSVersionTable.PSVersion
    if ($psVersion.Major -ge 5) {
        Write-Success "PowerShell $psVersion encontrado"
        $results.PowerShell = $true
    } else {
        Write-Error "PowerShell $psVersion encontrado, mas versão 5.1 ou superior é necessária"
    }
    
    # Verificar Chocolatey
    try {
        choco --version | Out-Null
        Write-Success "Chocolatey encontrado"
        $results.Chocolatey = $true
    }
    catch {
        Write-Info "Chocolatey não encontrado (opcional)"
    }
    
    # Verificar Winget
    try {
        winget --version | Out-Null
        Write-Success "Winget encontrado"
        $results.Winget = $true
    }
    catch {
        Write-Info "Winget não encontrado (opcional)"
    }
    
    return $results
}

# Instalar pré-requisitos
function Install-Prerequisites {
    Write-Step "Instalando Pré-requisitos"
    
    $results = Test-Prerequisites
    
    # Instalar Go se necessário
    if (-not $results.Go) {
        Write-Info "Instalando Go..."
        
        if ($results.Chocolatey) {
            choco install golang -y
        } elseif ($results.Winget) {
            winget install GoLang.Go
        } else {
            Write-Error "Chocolatey ou Winget necessário para instalar Go automaticamente"
            Write-Info "Instale manualmente: https://golang.org/dl/"
            return $false
        }
    }
    
    # Instalar Git se necessário
    if (-not $results.Git) {
        Write-Info "Instalando Git..."
        
        if ($results.Chocolatey) {
            choco install git -y
        } elseif ($results.Winget) {
            winget install Git.Git
        } else {
            Write-Error "Chocolatey ou Winget necessário para instalar Git automaticamente"
            Write-Info "Instale manualmente: https://git-scm.com/download/win"
            return $false
        }
    }
    
    # Instalar Chocolatey se necessário
    if (-not $results.Chocolatey) {
        Write-Info "Instalando Chocolatey..."
        try {
            Set-ExecutionPolicy Bypass -Scope Process -Force
            [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
            iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
            Write-Success "Chocolatey instalado"
        }
        catch {
            Write-Warning "Falha ao instalar Chocolatey: $($_.Exception.Message)"
        }
    }
    
    Write-Success "Instalação de pré-requisitos concluída"
    return $true
}

# Configurar ambiente
function Configure-Environment {
    Write-Step "Configurando Ambiente"
    
    # Navegar para o diretório CLI
    Set-Location $CLIDir
    
    # Verificar se estamos no diretório correto
    if (-not (Test-Path "main.go")) {
        Write-Error "main.go não encontrado. Execute este script do diretório CLI."
        return $false
    }
    
    # Configurar Go modules
    Write-Info "Configurando Go modules..."
    go mod download
    go mod tidy
    go mod verify
    
    # Criar diretórios necessários
    $directories = @("build", "logs", "temp", "dist", "scripts")
    foreach ($dir in $directories) {
        if (-not (Test-Path $dir)) {
            New-Item -ItemType Directory -Path $dir -Force | Out-Null
            Write-Info "Diretório criado: $dir"
        }
    }
    
    # Configurar Git (se necessário)
    try {
        git config --global user.name | Out-Null
        git config --global user.email | Out-Null
        Write-Success "Git configurado"
    }
    catch {
        Write-Warning "Git não configurado. Configure com:"
        Write-Info "git config --global user.name 'Seu Nome'"
        Write-Info "git config --global user.email 'seu@email.com'"
    }
    
    # Configurar Go environment
    Write-Info "Configurando Go environment..."
    $goEnv = @{
        GOOS = "windows"
        GOARCH = "amd64"
        CGO_ENABLED = "0"
    }
    
    foreach ($key in $goEnv.Keys) {
        [Environment]::SetEnvironmentVariable($key, $goEnv[$key], "User")
        Write-Info "Variável de ambiente definida: $key=$($goEnv[$key])"
    }
    
    Write-Success "Ambiente configurado"
    return $true
}

# Verificar configuração
function Verify-Configuration {
    Write-Step "Verificando Configuração"
    
    # Verificar Go
    try {
        $goVersion = go version
        Write-Success "Go: $goVersion"
    }
    catch {
        Write-Error "Go não encontrado"
        return $false
    }
    
    # Verificar Git
    try {
        $gitVersion = git --version
        Write-Success "Git: $gitVersion"
    }
    catch {
        Write-Error "Git não encontrado"
        return $false
    }
    
    # Verificar diretório CLI
    if (Test-Path (Join-Path $CLIDir "main.go")) {
        Write-Success "Diretório CLI: OK"
    } else {
        Write-Error "Diretório CLI: main.go não encontrado"
        return $false
    }
    
    # Verificar Go modules
    try {
        go mod verify
        Write-Success "Go modules: OK"
    }
    catch {
        Write-Warning "Go modules: Problemas detectados"
    }
    
    # Verificar variáveis de ambiente
    $envVars = @("GOOS", "GOARCH", "CGO_ENABLED")
    foreach ($var in $envVars) {
        $value = [Environment]::GetEnvironmentVariable($var, "User")
        if ($value) {
            Write-Success "Variável $var: $value"
        } else {
            Write-Warning "Variável $var: Não definida"
        }
    }
    
    # Testar compilação
    Write-Info "Testando compilação..."
    try {
        go build -o "temp\test-build.exe" main.go
        if (Test-Path "temp\test-build.exe") {
            Remove-Item "temp\test-build.exe" -Force
            Write-Success "Compilação de teste: OK"
        } else {
            Write-Error "Compilação de teste: Falhou"
            return $false
        }
    }
    catch {
        Write-Error "Compilação de teste: Falhou - $($_.Exception.Message)"
        return $false
    }
    
    Write-Success "Configuração verificada com sucesso"
    return $true
}

# Mostrar ajuda
function Show-Help {
    Write-Host "Uso: .\setup-environment.ps1 [ação]"
    Write-Host ""
    Write-Host "Ações:"
    Write-Host "  check      Verificar pré-requisitos (padrão)"
    Write-Host "  install    Instalar pré-requisitos"
    Write-Host "  configure  Configurar ambiente"
    Write-Host "  verify     Verificar configuração"
    Write-Host "  help       Mostrar esta ajuda"
    Write-Host ""
    Write-Host "Exemplos:"
    Write-Host "  .\setup-environment.ps1 check      # Verificar pré-requisitos"
    Write-Host "  .\setup-environment.ps1 install    # Instalar pré-requisitos"
    Write-Host "  .\setup-environment.ps1 configure  # Configurar ambiente"
    Write-Host "  .\setup-environment.ps1 verify     # Verificar configuração"
}

# Função principal
function Main {
    Show-Banner
    
    switch ($Action.ToLower()) {
        "check" {
            $results = Test-Prerequisites
            $allGood = $results.Go -and $results.Git -and $results.PowerShell
            
            if ($allGood) {
                Write-Success "Todos os pré-requisitos estão instalados!"
            } else {
                Write-Warning "Alguns pré-requisitos estão faltando."
                Write-Info "Execute '.\setup-environment.ps1 install' para instalar automaticamente."
            }
        }
        "install" {
            $success = Install-Prerequisites
            if ($success) {
                Write-Success "Instalação concluída!"
                Write-Info "Execute '.\setup-environment.ps1 configure' para configurar o ambiente."
            } else {
                Write-Error "Instalação falhou!"
            }
        }
        "configure" {
            $success = Configure-Environment
            if ($success) {
                Write-Success "Configuração concluída!"
                Write-Info "Execute '.\setup-environment.ps1 verify' para verificar a configuração."
            } else {
                Write-Error "Configuração falhou!"
            }
        }
        "verify" {
            $success = Verify-Configuration
            if ($success) {
                Write-Success "Ambiente configurado e pronto para uso!"
                Write-Info "Execute '.\build-windows.ps1 build' para compilar a aplicação."
            } else {
                Write-Error "Verificação falhou!"
            }
        }
        "help" {
            Show-Help
        }
        default {
            Write-Error "Ação desconhecida: $Action"
            Write-Host "Use '.\setup-environment.ps1 help' para opções disponíveis"
            exit 1
        }
    }
}

# Executar função principal
Main
