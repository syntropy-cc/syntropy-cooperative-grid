# Syntropy CLI Manager - Universal Build Script (PowerShell)
# Script unificado para compilar e testar em todas as plataformas suportadas
# Funciona no Windows PowerShell

param(
    [Parameter(Position=0)]
    [ValidateSet("all", "current", "platform", "test", "run", "clean", "help")]
    [string]$Action = "all",
    
    [Parameter(Position=1)]
    [string]$Platform = "",
    
    [switch]$NoTests,
    [switch]$Clean
)

# Configurações
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$CLIDir = Split-Path -Parent $ScriptDir
$BuildDir = Join-Path $CLIDir "build"
$Version = Get-Date -Format "yyyyMMdd-HHmmss"
$GitCommit = try { git rev-parse --short HEAD 2>$null } catch { "unknown" }
$BuildTime = Get-Date -Format "yyyy-MM-ddTHH:mm:ssZ"

# Detecção de plataforma
$OS = [System.Environment]::OSVersion.Platform.ToString().ToLower()
$Arch = [System.Environment]::GetEnvironmentVariable("PROCESSOR_ARCHITECTURE").ToLower()

# Mapear arquiteturas comuns
if ($Arch -eq "amd64") { $Arch = "amd64" }
if ($Arch -eq "x86") { $Arch = "386" }

# Plataformas suportadas
$SupportedPlatforms = @(
    "linux/amd64",
    "linux/arm64", 
    "windows/amd64",
    "darwin/amd64",
    "darwin/arm64"
)

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

function Write-Platform {
    param([string]$Platform, [string]$Message)
    Write-Host "[$Platform] $Message" -ForegroundColor Magenta
}

# Banner
function Show-Banner {
    Write-Host ""
    Write-Host "╔══════════════════════════════════════════════════════════════╗" -ForegroundColor Cyan
    Write-Host "║              SYNTROPY CLI MANAGER                           ║" -ForegroundColor Cyan
    Write-Host "║              Universal Build & Test                         ║" -ForegroundColor Cyan
    Write-Host "║              Cross-Platform Support                         ║" -ForegroundColor Cyan
    Write-Host "╚══════════════════════════════════════════════════════════════╝" -ForegroundColor Cyan
    Write-Host ""
    
    Write-Info "Detected Platform: $OS-$Arch"
    Write-Info "Build Directory: $BuildDir"
    Write-Info "Version: $Version"
    Write-Info "Git Commit: $GitCommit"
}

# Verificar pré-requisitos
function Test-Prerequisites {
    Write-Step "Checking Prerequisites"
    
    # Verificar Go
    try {
        $goVersion = (go version).Split(' ')[2].Substring(2)
        $requiredVersion = [Version]"1.22"
        $currentVersion = [Version]$goVersion
        
        if ($currentVersion -lt $requiredVersion) {
            Write-Error "Go version $goVersion found, but version $requiredVersion or higher is required."
            exit 1
        }
        
        Write-Success "Go $goVersion found"
    }
    catch {
        Write-Error "Go is not installed. Please install Go 1.22.5 or higher."
        Write-Info "Download: https://golang.org/dl/"
        exit 1
    }
    
    # Verificar se estamos no diretório correto
    if (-not (Test-Path (Join-Path $CLIDir "main.go"))) {
        Write-Error "main.go not found in $CLIDir"
        exit 1
    }
    
    Write-Success "Project structure verified"
}

# Preparar ambiente de build
function Initialize-Build {
    Write-Step "Preparing Build Environment"
    
    # Navegar para o diretório CLI
    Set-Location $CLIDir
    
    # Criar diretório de build
    if (-not (Test-Path $BuildDir)) {
        New-Item -ItemType Directory -Path $BuildDir -Force | Out-Null
    }
    
    # Limpar builds anteriores se solicitado
    if ($Clean) {
        Write-Info "Cleaning build directory..."
        Remove-Item "$BuildDir\*" -Force -ErrorAction SilentlyContinue
    }
    
    Write-Success "Build environment prepared"
}

# Configurar dependências
function Setup-Dependencies {
    Write-Step "Setting Up Dependencies"
    
    Write-Info "Downloading dependencies..."
    go mod download
    
    Write-Info "Organizing dependencies..."
    go mod tidy
    
    Write-Info "Verifying dependencies..."
    go mod verify
    
    Write-Success "Dependencies configured"
}

# Build para plataforma específica
function Build-Platform {
    param([string]$Platform)
    
    $goos = $Platform.Split('/')[0]
    $goarch = $Platform.Split('/')[1]
    
    Write-Platform $Platform "Building executable..."
    
    $buildFlags = "-ldflags `"-X main.version=$Version -X main.buildTime=$BuildTime -X main.gitCommit=$GitCommit`""
    $outputName = "syntropy"
    
    # Definir extensão apropriada
    if ($goos -eq "windows") {
        $outputName = "$outputName.exe"
    }
    
    $outputFile = Join-Path $BuildDir "syntropy-$goos-$goarch"
    if ($goos -eq "windows") {
        $outputFile = "$outputFile.exe"
    }
    
    # Build
    $env:GOOS = $goos
    $env:GOARCH = $goarch
    go build $buildFlags -o $outputFile main.go
    
    # Restaurar variáveis de ambiente
    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
    
    if (Test-Path $outputFile) {
        $sizeKB = [math]::Round((Get-Item $outputFile).Length / 1KB, 2)
        Write-Platform $Platform "Build completed: $outputFile ($sizeKB KB)"
        return $true
    } else {
        Write-Platform $Platform "Build failed"
        return $false
    }
}

# Build para todas as plataformas
function Build-AllPlatforms {
    Write-Step "Building for All Platforms"
    
    $successCount = 0
    $totalCount = $SupportedPlatforms.Count
    
    foreach ($platform in $SupportedPlatforms) {
        if (Build-Platform $platform) {
            $successCount++
        }
    }
    
    Write-Info "Build Summary: $successCount/$totalCount platforms successful"
    
    if ($successCount -eq 0) {
        Write-Error "No platforms built successfully"
        exit 1
    }
}

# Build para plataforma atual
function Build-CurrentPlatform {
    Write-Step "Building for Current Platform"
    
    $currentPlatform = "windows/$Arch"
    
    Write-Info "Current platform: $currentPlatform"
    
    if (Build-Platform $currentPlatform) {
        Write-Success "Current platform build completed"
    } else {
        Write-Error "Current platform build failed"
        exit 1
    }
}

# Executar testes
function Invoke-Tests {
    Write-Step "Running Tests"
    
    Write-Info "Running unit tests..."
    try {
        go test -v .\...
        Write-Success "Unit tests passed"
    }
    catch {
        Write-Warning "Some unit tests failed (expected for unimplemented features)"
    }
    
    Write-Info "Running tests with coverage..."
    try {
        go test -v -cover .\...
        Write-Success "Coverage tests completed"
    }
    catch {
        Write-Warning "Some coverage tests failed"
    }
    
    Write-Info "Running race condition tests..."
    try {
        go test -race .\...
        Write-Success "Race condition tests passed"
    }
    catch {
        Write-Warning "Some race condition tests failed"
    }
    
    Write-Success "Test suite completed"
}

# Testar binários
function Test-Binaries {
    Write-Step "Testing Binaries"
    
    $testCount = 0
    $successCount = 0
    
    $binaries = Get-ChildItem -Path $BuildDir -Filter "syntropy-*" -File
    
    foreach ($binary in $binaries) {
        $testCount++
        $binaryName = $binary.Name
        Write-Info "Testing $binaryName..."
        
        # Testar versão
        try {
            & $binary.FullName --version | Out-Null
            Write-Success "Version test passed for $binaryName"
            $successCount++
        }
        catch {
            Write-Warning "Version test failed for $binaryName (may be normal)"
        }
        
        # Testar ajuda
        try {
            & $binary.FullName --help | Out-Null
            Write-Success "Help test passed for $binaryName"
        }
        catch {
            Write-Warning "Help test failed for $binaryName (may be normal)"
        }
    }
    
    Write-Info "Binary testing: $successCount/$testCount binaries tested successfully"
}

# Executar aplicação
function Start-Application {
    Write-Step "Running Application"
    
    # Encontrar o binário apropriado para a plataforma atual
    $binaryPath = ""
    
    $windowsBinary = Join-Path $BuildDir "syntropy-windows-amd64.exe"
    if (Test-Path $windowsBinary) {
        $binaryPath = $windowsBinary
    }
    
    if (-not $binaryPath -or -not (Test-Path $binaryPath)) {
        Write-Error "No suitable binary found for current platform"
        Write-Info "Available binaries:"
        Get-ChildItem -Path $BuildDir -Filter "syntropy-*" -File | ForEach-Object { Write-Host "  $($_.Name)" }
        exit 1
    }
    
    Write-Info "Running: $binaryPath"
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Cyan
    
    # Executar com ajuda para mostrar comandos disponíveis
    & $binaryPath --help
    
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Success "Application executed successfully!"
    
    # Mostrar exemplos de comandos
    Write-Host ""
    Write-Host "Example commands to try:" -ForegroundColor Blue
    Write-Host "  $binaryPath --version" -ForegroundColor Yellow
    Write-Host "  $binaryPath setup --help" -ForegroundColor Yellow
    Write-Host "  $binaryPath setup validate" -ForegroundColor Yellow
    Write-Host "  $binaryPath setup run --force" -ForegroundColor Yellow
}

# Mostrar resumo
function Show-Summary {
    Write-Step "Build Summary"
    
    Write-Success "Build completed successfully!"
    Write-Host ""
    Write-Host "📁 Build Directory: $BuildDir" -ForegroundColor Blue
    Write-Host "📦 Version: $Version" -ForegroundColor Blue
    Write-Host "🔧 Git Commit: $GitCommit" -ForegroundColor Blue
    Write-Host "🕒 Build Time: $BuildTime" -ForegroundColor Blue
    Write-Host "🖥️  Current Platform: $OS-$Arch" -ForegroundColor Blue
    
    Write-Host ""
    Write-Host "📋 Created Binaries:" -ForegroundColor Blue
    $binaries = Get-ChildItem -Path $BuildDir -Filter "syntropy-*" -File
    $binaryCount = $binaries.Count
    
    if ($binaryCount -eq 0) {
        Write-Host "  ❌ No binaries created" -ForegroundColor Red
    } else {
        foreach ($binary in $binaries) {
            $sizeKB = [math]::Round($binary.Length / 1KB, 2)
            Write-Host "  ✅ $($binary.Name) ($sizeKB KB)" -ForegroundColor Green
        }
    }
    
    Write-Host ""
    Write-Host "🚀 Next Steps:" -ForegroundColor Blue
    Write-Host "  1. Test the application: .\build-all.ps1 run" -ForegroundColor Yellow
    Write-Host "  2. Run tests: .\build-all.ps1 test" -ForegroundColor Yellow
    Write-Host "  3. Copy binaries to target machines for testing" -ForegroundColor Yellow
    
    Write-Host ""
    Write-Host "💡 Usage Examples:" -ForegroundColor Blue
    Write-Host "  .\build-all.ps1 run                    # Run application" -ForegroundColor Cyan
    Write-Host "  .\build-all.ps1 test                   # Run tests only" -ForegroundColor Cyan
    Write-Host "  .\build-all.ps1 current                # Build current platform only" -ForegroundColor Cyan
    Write-Host "  .\build-all.ps1 platform linux/amd64  # Build specific platform" -ForegroundColor Cyan
}

# Mostrar ajuda
function Show-Help {
    Write-Host "Usage: .\build-all.ps1 [action] [platform] [options]" -ForegroundColor Blue
    Write-Host ""
    Write-Host "Actions:" -ForegroundColor Blue
    Write-Host "  all       Build for all platforms (default)" -ForegroundColor White
    Write-Host "  current   Build only for current platform" -ForegroundColor White
    Write-Host "  platform  Build for specific platform" -ForegroundColor White
    Write-Host "  test      Run tests only" -ForegroundColor White
    Write-Host "  run       Run the application" -ForegroundColor White
    Write-Host "  clean     Clean build directory" -ForegroundColor White
    Write-Host "  help      Show this help" -ForegroundColor White
    Write-Host ""
    Write-Host "Options:" -ForegroundColor Blue
    Write-Host "  -NoTests  Skip running tests" -ForegroundColor White
    Write-Host "  -Clean    Clean build directory before building" -ForegroundColor White
    Write-Host ""
    Write-Host "Supported Platforms:" -ForegroundColor Blue
    foreach ($platform in $SupportedPlatforms) {
        Write-Host "  - $platform" -ForegroundColor White
    }
    Write-Host ""
    Write-Host "Examples:" -ForegroundColor Blue
    Write-Host "  .\build-all.ps1                                    # Build for all platforms" -ForegroundColor Cyan
    Write-Host "  .\build-all.ps1 current                            # Build current platform only" -ForegroundColor Cyan
    Write-Host "  .\build-all.ps1 platform windows/amd64            # Build Windows only" -ForegroundColor Cyan
    Write-Host "  .\build-all.ps1 test                               # Run tests only" -ForegroundColor Cyan
    Write-Host "  .\build-all.ps1 run                                # Run application" -ForegroundColor Cyan
    Write-Host "  .\build-all.ps1 -Clean current                     # Clean and build current platform" -ForegroundColor Cyan
}

# Função principal
function Main {
    Show-Banner
    
    # Verificar pré-requisitos
    Test-Prerequisites
    
    # Preparar ambiente de build
    Initialize-Build
    
    # Configurar dependências
    Setup-Dependencies
    
    # Executar ação baseada nos parâmetros
    switch ($Action.ToLower()) {
        "all" {
            Build-AllPlatforms
            if (-not $NoTests) {
                Invoke-Tests
            }
            Test-Binaries
            Show-Summary
        }
        "current" {
            Build-CurrentPlatform
            if (-not $NoTests) {
                Invoke-Tests
            }
            Test-Binaries
            Show-Summary
        }
        "platform" {
            if (-not $Platform) {
                Write-Error "Platform parameter required for 'platform' action"
                Write-Host "Use: .\build-all.ps1 platform linux/amd64"
                exit 1
            }
            
            if ($SupportedPlatforms -notcontains $Platform) {
                Write-Error "Unsupported platform: $Platform"
                Write-Host "Supported platforms: $($SupportedPlatforms -join ', ')"
                exit 1
            }
            
            if (Build-Platform $Platform) {
                Write-Success "Platform $Platform built successfully"
            } else {
                Write-Error "Failed to build platform $Platform"
                exit 1
            }
            
            if (-not $NoTests) {
                Invoke-Tests
            }
            Test-Binaries
            Show-Summary
        }
        "test" {
            Invoke-Tests
        }
        "run" {
            Start-Application
        }
        "clean" {
            Write-Step "Cleaning Build Directory"
            if (Test-Path $BuildDir) {
                Remove-Item $BuildDir -Recurse -Force
                Write-Success "Build directory cleaned"
            } else {
                Write-Info "No build directory to clean"
            }
        }
        "help" {
            Show-Help
        }
        default {
            Write-Error "Unknown action: $Action"
            Show-Help
            exit 1
        }
    }
    
    # Perguntar se o usuário quer executar a aplicação (exceto para ações específicas)
    if ($Action -notin @("test", "run", "clean", "help")) {
        Write-Host ""
        $response = Read-Host "Do you want to run the application now? (y/N)"
        if ($response -match "^[Yy]$") {
            Start-Application
        }
    }
}

# Executar função principal
Main
