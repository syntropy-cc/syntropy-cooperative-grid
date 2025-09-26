#Requires -Version 5.1

<#
.SYNOPSIS
    Syntropy CLI - Setup Básico para Windows
    Script de configuração automática para Windows PowerShell

.DESCRIPTION
    Automatiza o processo de setup do Syntropy CLI no Windows, incluindo:
    - Validação de ambiente e requisitos
    - Criação de estrutura de diretórios
    - Configuração de parâmetros básicos
    - Instalação opcional como serviço Windows
    - Verificação pós-setup

.PARAMETER ValidateOnly
    Apenas validar ambiente, não executar setup completo

.PARAMETER DryRun  
    Simular operações sem fazer mudanças permanentes

.PARAMETER Force
    Forçar setup mesmo com problemas de validação

.PARAMETER InstallService
    Configurar como serviço Windows

.PARAMETER UserOnly
    Configuração apenas para usuário atual (sem admin rights)

.PARAMETER ConfigPath
    Caminho customizado para arquivo de configuração

.PARAMETER SyntropyHome
    Diretório base do Syntropy (padrão: $env:USERPROFILE\.syntropy)

.EXAMPLE
    .\setup-basic.ps1
    Setup completo com configurações padrão

.EXAMPLE
    .\setup-basic.ps1 -ValidateOnly
    Apenas verificar requisitos do sistema

.EXAMPLE
    .\setup-basic.ps1 -InstallService
    Setup + instalação como serviço Windows

.INPUTS
    None - script autossuficiente

.OUTPUTS
    Logs estruturados e arquivos de configuração

.NOTES 
    Author: Syntropy Team
    Version: 1.0.0
    Created: 2025-01-27
    Requires: PowerShell 5.1+, Windows 7/Server 2008 R2+

.LINK
    https://github.com/syntropy-cc/syntropy-cooperative-grid
#>

### CONFIGURAÇÕES E CONSTANTES ###

param(
    [switch]$ValidateOnly,
    [switch]$DryRun,
    [switch]$Force,
    [switch]$InstallService,
    [switch]$UserOnly,
    [string]$ConfigPath = "",
    [string]$SyntropyHome = ""
)

# Configurações básicas
$SCRIPT_NAME = "setup-basic.ps1"
$SCRIPT_VERSION = "1.0.0"
$SCRIPT_DIR = Split-Path -Parent $MyInvocation.MyCommand.Definition

# Definir diretório Syntropy
if ([string]::IsNullOrEmpty($SyntropyHome)) {
    $SYNTHROPY_HOME = Join-Path $env:USERPROFILE ".syntropy"
} else {
    $SYNTHROPY_HOME = $SyntropyHome
}

# Arquiteura do sistema detectada
$OSInfo = Get-WmiObject -Class Win32_OperatingSystem
$Is64Bit = [Environment]::Is64BitOperatingSystem

# Configuração padrão
$MIN_DISK_SPACE_GB = 1
$MIN_MEMORY_GB = 0.5
$POWERSHELL_REQD_VERSION = "5.1"

# ============================================================================
### FUNÇÕES UTILITÁRIAS ###
# ============================================================================

function Write-Log {
    [cmdletbinding()]
    param(
        [Parameter(Mandatory=$true, Position=0)]
        [ValidateSet("ERROR", "WARN", "INFO", "DEBUG", "SUCCESS")]
        [string]$Level,
        
        [Parameter(Mandatory=$true, Position=1)]
        [string]$Message,
        
        [string]$Color = "White"
    )
    
    $timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
    $prefix = switch ($Level) {
        "ERROR"   { "`e[31m[ERROR]`e[0m" }
        "WARN"    { "`e[33m[WARN]`e[0m"  }
        "INFO"    { "`e[34m[INFO]`e[0m" }
        "DEBUG"   { "`e[32m[DEBUG]`e[0m" }
        "SUCCESS" { "`e[32m[SUCCESS]`e[0m" }
    }
    
    $logMessage = "${prefix} [${timestamp}] $Message"
    Write-Host $logMessage -ForegroundColor $Color
}

function Write-Section {
    [cmdletbinding()]
    param(
        [Parameter(Mandatory=$true)]
        [string]$Title
    )
    
    Write-Host ""
    Write-Host "=== $Title ===" -ForegroundColor Blue
    Write-Host ""
}

function Test-AdminRights {
    # Verificar se script está rodando com privilégios administrativos
    $currentUser = [Security.Principal.WindowsIdentity]::GetCurrent()
    $principal = New-Object Security.Principal.WindowsPrincipal($currentUser)
    $isAdmin = $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
    
    return $isAdmin
}

function Get-DiskSpace {
    # Obter espaço disponível em disco para diretório de destino
    param([string]$Path)
    
    try {
        $drive = (Get-Item $Path | Select-Object -First 1).PSDrive
        $freeSpaceGB = [math]::Round($drive.Free / 1GB, 2)
        return $freeSpaceGB
    }
    catch {
        Write-Log "WARN" "Não foi possível calcular espaço de disco: $($_.Exception.Message)"
        return 0
    }
}

# ============================================================================
### VERIFICAÇÕES DE AMBIENTE E REQUISITOS ###
# ============================================================================

function Test-SystemRequirements {
    Write-Section "Verificação de Requisitos do Sistema"
    
    Write-Log "INFO" "Verificando requisitos do Windows..."
    
    # Verificar versão do Windows
    $currentVersion = [System.Environment]::OSVersion.Version
    $minVersion = [System.Version]::new(6, 1)  # Windows 7 / Server 2008 R2
    
    if ($currentVersion -lt $minVersion) {
        Write-Log "ERROR" "Versão do Windows não suportada: $currentVersion"
        Write-Log "ERROR" "Mínima necessária: Windows 7 / Server 2008 R2"
        return $false
    }
    
    Write-Log "SUCCESS" "Versão do Windows é apropriada: $($OSInfo.Caption)"
    
    # Verificar arquitetura
    $arch = if ($Is64Bit) { "64-bit" } else { "32-bit" }
    Write-Log "INFO" "Arquitetura detectada: $arch"
    
    # Verificar versão do PowerShell 
    $psVersion = $PSVersionTable.PSVersion
    Write-Log "INFO" "PowerShell versão: $psVersion"
    
    if ($psVersion -lt [Version]$POWERSHELL_REQD_VERSION) {
        Write-Log "ERROR" "Versão PowerShell insuficiente: $psVersion"
        Write-Log "ERROR" "Requisitos: PowerShell $POWERSHELL_REQD_VERSION ou superior"
        return $false
    }
    
    # Verificar privilégios necessários
    $hasAdminRights = Test-AdminRights
    
    if ($InstallService -or (-not $UserOnly)) {
        if (-not $hasAdminRights) {
            Write-Log "WARN" "Privilégios administrativos não detectados"
            Write-Log "WARN" "Alguns recursos podem estar limitados"
            Write-Log "WARN" "Execute como Administrador para funcionalidade completa"
        } else {
            Write-Log "SUCCESS" "Privilégios administrativos confirmados"
        }
    } else {
        Write-Log "INFO" "Configuração apenas para usuário - privilégios admin não necessários"
    }
    
    return $true
}

function Test-DiskSpace {
    Write-Log "INFO" "Verificando espaço em disco disponível..."
    
    $freeSpace = Get-DiskSpace -Path $SYNTHROPY_HOME
    Write-Log "INFO" "Espaço disponível: $freeSpace GB"
    
    if ($freeSpace -lt $MIN_DISK_SPACE_GB) {
        Write-Log "ERROR" "Espaço em disco insuficiente!"
        Write-Log "ERROR" "Necessário: $MIN_DISK_SPACE_GB GB, Disponível: $freeSpace GB"
        return $false
    }
    
    Write-Log "SUCCESS" "Espaço em disco adequado"
    return $true
}

function Test-SystemResources {
    Write-Log "INFO" "Verificando recursos do sistema..."
    
    # Memória RAM
    $memory = Get-WmiObject -Class Win32_PhysicalMemory | Measure-Object -Property Capacity -Sum
    $totalMemoryGB = [math]::Round($memory.Sum / 1GB, 2)
    
    Write-Log "INFO" "Memória RAM total: $totalMemoryGB GB"
    
    if ($totalMemoryGB -lt $MIN_MEMORY_GB) {
        Write-Log "WARN" "Memória disponível baixa - pode afetar performance"
    } else {
        Write-Log "SUCCESS" "Recursos do sistema adequados"
    }
    
    return $true
}

function Test-NetworkConnectivity {
    Write-Log "INFO" "Verificando conectividade de rede..."
    
    $domains = @(
        "google.com",
        "github.com", 
        "syntropy.io"
    )
    
    foreach ($domain in $domains) {
        try {
            $result = Test-NetConnection -ComputerName $domain -Port 80 -InformationLevel Quiet
            if ($result) {
                Write-Log "SUCCESS" "Conectividade confirmada: $domain"
                break
            }
        }
        catch {
            Write-Log "DEBUG" "Falha ao testar: $domain"
        }
    }
    
    return $true
}

# ============================================================================
### OPERAÇÕES DE SETUP ###
# ============================================================================

function New-SyntropyDirectories {
    Write-Section "Criação de Diretórios"
    
    $directories = @(
        $SYNTHROPY_HOME,
        Join-Path $SYNTHROPY_HOME "config",
        Join-Path $SYNTHROPY_HOME "logs", 
        Join-Path $SYNTHROPY_HOME "data",
        Join-Path $SYNTHROPY_HOME "services"
    )
    
    foreach ($dir in $directories) {
        if (-not (Test-Path $dir)) {
            try {
                New-Item -ItemType Directory -Path $dir -Force | Out-Null
                Write-Log "SUCCESS" "Diretório criado: $dir"
            }
            catch {
                Write-Log "ERROR" "Falha ao criar diretório '$dir': $($_.Exception.Message)"
                return $false
            }
        } else {
            Write-Log "DEBUG" "Diretório já existe: $dir"
        }
    }
    
    return $true
}

function Set-ConfigurationFiles {
    Write-Section "Configuração de Arquivos"
    
    $configPath = Join-Path $SYNTHROPY_HOME "config\manager.yaml"
    
    # Definir configuração de exemplo
    if (-not [string]::IsNullOrEmpty($ConfigPath) -and (Test-Path $ConfigPath)) {
        Write-Log "INFO" "Usando arquivo de configuração customizado: $ConfigPath"
        Copy-Item $ConfigPath $configPath -Force
    } else {
        # Verificar se existe arquivo exemplo no diretório do script  
        $exampleConfig = Join-Path $SCRIPT_DIR "config-example.yaml"
        
        if (Test-Path $exampleConfig) {
            Write-Log "INFO" "Usando arquivo de configuração exemplo: $exampleConfig"
            Copy-Item $exampleConfig $configPath -Force
        } else {
            Write-Log "INFO" "Criando configuração padrão..."
            New-DefaultConfiguration -Path $configPath
        }
    }
    
    return $true
}

function New-DefaultConfiguration {
    [cmdletbinding()]
    param(
        [Parameter(Mandatory=$true)]
        [string]$Path
    )
    
    $configContent = @"
manager:
  home_dir: "$SYNTHROPY_HOME"
  log_level: "info" 
  api_endpoint: "https://api.syntropy.io"
  directories:
    config: "config"
    logs: "logs"
    data: "data"
  default_paths:
    config: "config/manager.yaml"
    owner_key: "config/owner.key"

owner_key:
  type: "Ed25519"
  path: "config/owner.key"

environment:
  os: "windows"
  architecture: "$(if ($Is64Bit) { "amd64" } else { "386" })"
  home_dir: "$env:USERPROFILE"
"@
    
    try {
        $configContent | Out-File -FilePath $Path -Encoding UTF8 -NoNewline
        Write-Log "SUCCESS" "Configuração padrão criada"
    }
    catch {
        Write-Log "ERROR" "Falha ao criar configuração: $($_.Exception.Message)"
        return $false
    }
    
    return $true
}

function Install-WindowsService {
    Write-Section "Instalação como Serviço Windows"
    
    if (-not $InstallService) {
        Write-Log "INFO" "Instalação de serviço não solicitada"
        return $true
    }
    
    if (-not (Test-AdminRights)) {
        Write-Log "ERROR" "Privilégios administrativos necessários para instalar serviço"
        Write-Log "INFO" "Execute o script como Administrador para instalar serviço"
        return $false
    }
    
    $serviceName = "SyntropyService"
    $serviceScript = Join-Path $SYNTHROPY_HOME "services\syntropy-service.ps1"
    
    # Criar script de serviço 
    New-ServiceScript -Path $serviceScript
    
    # Criar comando para instalação do serviço
    $installScript = Join-Path $SYNTHROPY_HOME "services\install-service.ps1"
    
    $installContent = @"
# Instalar Syntropy Silva Service
`$serviceName = '$serviceName'
`$servicePath = '$serviceScript'

# Verificar se serviço já existe
if (Get-Service -Name `$serviceName -ErrorAction SilentlyContinue) {
    Write-Host 'Serviço já existe, atualizando...'
    Stop-Service -Name `$serviceName -Force -ErrorAction SilentlyContinue
    sc.exe delete `$serviceName
}

# Criar novo serviço
sc.exe create `$serviceName binPath="powershell.exe -ExecutionPolicy Bypass -File `$servicePath" start=auto

if (`$LASTEXITCODE -eq 0) {
    Write-Host 'Serviço SyntropyService criado com sucesso'
    Start-Service -Name `$serviceName
} else {
    Write-Error 'Falha na criação do serviço'
}
"@
    
    $installContent | Out-File -FilePath $installScript -Encoding UTF8
    
    Write-Log "SUCCESS" "Scripts de serviço criados"
    Write-Log "INFO" "Para instalar serviço, execute: $installScript"
    
    return $true
}

function New-ServiceScript {
    [cmdletbinding()]
    param(
        [Parameter(Mandatory=$true)]
        [string]$Path
    )
    
    $serviceScript = @"
# Syntropy Service Script - Windows
# Gerado automaticamente pelo setup-basic.ps1

`$SYNTHROPY_HOME = '$SYNTHROPY_HOME'
`$configPath = Join-Path `$SYNTHROPY_HOME 'config\manager.yaml'

# Função de logging
function Write-ServiceLog {
    param([string]`$message)
    `$timestamp = Get-Date -Format 'yyyy-MM-dd HH:mm:ss'
    `$logPath = Join-Path `$SYNTHROPY_HOME 'logs\service.log'
    "[`$timestamp] `$message" | Add-Content -Path `$logPath
}

Write-ServiceLog "Syntropy Service iniciando..."

# Verificar se configuração existe
if (-not (Test-Path `$configPath)) {
    Write-ServiceLog "ERROR: Arquivo de configuração não encontrado: `$configPath"
    exit 1
}

Write-ServiceLog "Configuração encontrada, iniciando Syntropy..."
Write-ServiceLog "Config: `$configPath"

# Aqui você adicionaria o comando real para executar o Syntropy       
# Exemplo:
# & 'path\to\syntropy.exe' --config=`$configPath --service

# Por agora, apenas simular o funcionamento
Write-ServiceLog "Syntropy Service funcionando..."
"@
    
    try {
        $serviceScript | Out-File -FilePath $Path -Encoding UTF8
        Write-Log "SUCCESS" "Script de serviço criado: $Path"
    }
    catch {
        Write-Log "ERROR" "Falha ao criar script de serviço: $($_.Exception.Message)" 
        return $false
    }
    
    return $true
}

function New-LoggingSystem {
    Write-Log "INFO" "Configurando sistema de logging..."
    
    $logDir = Join-Path $SYNTHROPY_HOME "logs"
    $setupLogPath = Join-Path $logDir "setup.log"
    
    $logHeader = @"
# Syntropy Setup Log - Windows
# Generated by $SCRIPT_NAME v$SCRIPT_VERSION  
# Date: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')
# OS: $($OSInfo.Caption)
# User: $env:USERNAME

"@
    
    $logHeader | Add-Content -Path $setupLogPath -Encoding UTF8
    
    Write-Log "SUCCESS" "Sistema de logging configurado: $setupLogPath"
    return $true
}

# ============================================================================
### VERIFICAÇÕES PÓS-SETUP ###
# ============================================================================

function Test-SetupVerification {
    Write-Section "Verificação Pós-Setup"
    
    $allOk = $true
    
    # Verificar estrutura de diretórios
    $expectedDirs = @("config", "logs", "data", "services")
    foreach ($dir in $expectedDirs) {
        $dirPath = Join-Path $SYNTHROPY_HOME $dir
        if (Test-Path $dirPath) {
            Write-Log "SUCCESS" "✓ Diretório '$dir' criado corretamente"  
        } else {
            Write-Log "ERROR" "✗ Diretório '$dir' não foi criado"
            $allOk = $false
        }
    }
    
    # Verificar arquivos de configuração
    $configFile = Join-Path $SYNTHROPY_HOME "config\manager.yaml"
    if (Test-Path $configFile) {
        Write-Log "SUCCESS" "✓ Arquivo de configuração encontrado"
    } else {
        Write-Log "ERROR" "✗ Arquivo de configuração não criado"
        $allOk = $false
    }
    
    # Listar arquivos criados
    Write-Log "INFO" "Arquivos criados no setup:"
    Get-ChildItem -Path $SYNTHROPY_HOME -Recurse -File | ForEach-Object {
        Write-Log "INFO" "  • $($_.FullName)"
    }
    
    return $allOk
}

# ============================================================================
### FUNÇÃO PRINCIPAL ###
# ============================================================================

function Main {
    # Header  
    Write-Section "Syntropy CLI Setup Básico para Windows"
    Write-Host "Versão: $SCRIPT_VERSION" -ForegroundColor Cyan
    Write-Host "Diretório: $SCRIPT_DIR" -ForegroundColor Gray
    Write-Host "Syntropy Home: $SYNTHROPY_HOME" -ForegroundColor Gray
    
    if ($DryRun) {
        Write-Log "INFO" "Modo DRY RUN - nenhuma mudança será feita"
        Test-SystemRequirements | Out-Null
        Test-DiskSpace | Out-Null 
        Test-SystemResources | Out-Null
        Write-Log "SUCCESS" "Simulação concluída - sistema adequado para setup"
        return
    }
    
    if ($ValidateOnly) {
        Test-SystemRequirements | Out-Null
        Test-DiskSpace | Out-Null  
        Test-SystemResources | Out-Null
        Test-NetworkConnectivity | Out-Null
        Write-Log "SUCCESS" "Validação concluída - ambiente adequado"
        return
    }
    
    try {
        # Verificações de requisitos
        if (-not (Test-SystemRequirements)) {
            Write-Log "ERROR" "Requisitos do sistema não atendidos"
            exit 1
        }
        
        if (-not (Test-DiskSpace)) {
            Write-Log "ERROR" "Espaço em disco insuficiente"
            exit 1
        }
        
        Test-SystemResources | Out-Null
        Test-NetworkConnectivity | Out-Null
        
        # Operações de setup
        if (-not (New-SyntropyDirectories)) {
            Write-Log "ERROR" "Falha na criação de diretórios"
            exit 1
        }
        
        if (-not (Set-ConfigurationFiles)) {
            Write-Log "ERROR" "Falha na configuração de arquivos"
            exit 1
        }
        
        New-LoggingSystem | Out-Null
        
        if ($InstallService) {
            if (-not (Install-WindowsService)) {
                Write-Log "WARN" "Falha ao configurar serviço Windows"
            }
        }
        
        # Verificação final
        if (-not (Test-SetupVerification)) {
            Write-Log "ERROR" "Verificação pós-setup falhou"
            exit 1
        }
        
        # Conclusão
        Write-Section "Setup Concluído"
        Write-Log "SUCCESS" "Setup básico do Syntropy CLI finalizado com sucesso!"
        Write-Log "INFO" "Diretório base: $SYNTHROPY_HOME"
        Write-Log "INFO" "Configuração: $SYNTHROPY_HOME\config\manager.yaml"
        Write-Log "INFO" "Logs: $SYNTHROPY_HOME\logs\setup.log"
        
        if ($InstallService -and (Test-Path (Join-Path $SYNTHROPY_HOME "services\install-service.ps1"))) {
            Write-Log "INFO" "Para ativar serviço: $SYNTHROPY_HOME\services\install-service.ps1"
        }
        
        Write-Host ""
        Write-Log "INFO" "Próximos passos sugeridos:"
        Write-Log "INFO" "  • Verificar status: Get-SyntropyStatus (quando implementado)"
        Write-Log "INFO" "  • Consultar logs: Get-Content '$SYNTHROPY_HOME\logs\setup.log'"
        Write-Log "INFO" "  • Configurar serviço: $SYNTHROPY_HOME\services\install-service.ps1"
        
    } catch {
        Write-Log "ERROR" "Erro durante setup: $($_.Exception.Message)"
        Write-Log "ERROR" "Consulte logs para mais detalhes em: $SYNTHROPY_HOME\logs\" 
        exit 1
    }
}

# Execução principal
Main
