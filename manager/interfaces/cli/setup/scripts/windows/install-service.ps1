#Requires -Version 5.1

<#
.SYNOPSIS
    Syntropy CLI - Instalação como Serviço Windows
    Instala o Syntropy CLI como um serviço do sistema

.DESCRIPTION  
    Configura o Syntropy CLI para execução automática como
    serviço Windows incluindo dependências e policies

.PARAMETER ServiceName
    Nome do serviço (default: SyntropyService)

.PARAMETER AutoStart
    Habilitar auto-start (default: true)

.EXAMPLE
    .\install-service.ps1
    .\install-service.ps1 -ServiceName "SyntropyGrid" -AutoStart
#>

param(
    [string]$ServiceName = "SyntropyService",
    [switch]$AutoStart = $true,
    [switch]$Force = $false
)

Set-StrictMode -Version Latest

$ServiceUser = "LocalSystem"
$ServiceDescription = "Syntropy CLI Cooperative Grid Service"

# Colors  
$SUCCESS_COLOR = 'Green'
$ERROR_COLOR = 'Red'
$INFO_COLOR = 'Blue'

function Write-Log {
    param([string]$Level, [string]$Message)
    $timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
    $color = if ($Level -eq "ERROR") { $ERROR_COLOR } 
             elseif ($Level -eq "SUCCESS") { $SUCCESS_COLOR }
             else { $INFO_COLOR }
    
    Write-Host "[$Level] [$timestamp] $Message" -ForegroundColor $color
}

function Test-AdminPrivileges {
    $currentUser = [Security.Principal.WindowsIdentity]::GetCurrent()
    $principal = New-Object Security.Principal.WindowsPrincipal($currentUser)
    return $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

function Test-SyntropyCLI {
    try {
        $null = Get-Command syntropy -ErrorAction Stop
        return $true
    }
    catch {
        Write-Log "WARN" "Syntropy CLI not found in PATH"
        return $false
    }
}

function Install-WindowsService {
    Write-Log "INFO" "Iniciando instalação do serviço Windows"
    
    if (-not (Test-AdminPrivileges)) {
        Write-Log "ERROR" "Privilégios administrativos necessários - execute como Administrador"
        return $false
    }

    $serviceExists = Get-Service -Name $ServiceName -ErrorAction SilentlyContinue
    
    if ($serviceExists -and -not $Force) {
        Write-Log "WARN" "Serviço já existe: $ServiceName. Use -Force para substituir"
        return $false
    }
    
    # Clean existing service if it exists
    if ($serviceExists) {
        Write-Log "INFO" "Removendo serviço existente: $ServiceName"
        Stop-Service -Name $ServiceName -Force -ErrorAction SilentlyContinue
        Start-Sleep -Seconds 2
        Remove-Service -Name $ServiceName -ErrorAction SilentlyContinue
        Start-Sleep -Seconds 1
    }
    
    # Find Syntropy CLI path
    $syntropyPath = ""
    if (Test-SyntropyCLI) {
        $syntropyPath = (Get-Command syntropy).Source
    } else {
        # Look for syntropy in common installation paths
        $commonPaths = @(
            "$env:ProgramFiles\Syntropy",
            "$env:ProgramFiles(x86)\Syntropy",
            "$env:LOCALAPPDATA\Syntropy"
        )
        
        foreach ($path in $commonPaths) {
            $syntropyExe = Join-Path $path "syntropy.exe"
            if (Test-Path $syntropyExe) {
                $syntropyPath = $syntropyExe
                break
            }
        }
        
        if ([string]::IsNullOrEmpty($syntropyPath)) {
            Write-Log "ERROR" "Syntropy CLI not found. Please install it first."
            return $false
        }
    }
    
    $ServiceExecutable = "`"$syntropyPath`" run"
    $displayName = "Syntropy CLI Service"
    
    try {
        Write-Log "INFO" "Criando serviço: $ServiceName"
        New-Service -Name $ServiceName `
                   -BinaryPathName $ServiceExecutable `
                   -DisplayName $displayName `
                   -Description $ServiceDescription `
                   -StartupType ($AutoStart ? "Automatic" : "Manual") `
                   -ErrorAction Stop
                   
        Write-Log "SUCCESS" "Serviço '$ServiceName' criado com sucesso"
        return $true
    }
    catch {
        Write-Log "ERROR" "Erro criando serviço: $($_.Exception.Message)"
        return $false
    }
}

function Start-SyntropyService {
    param([string]$ServiceName)
    
    try {
        Set-Service -Name $ServiceName -StartupType Automatic -ErrorAction Stop
        Start-Service -Name $ServiceName -ErrorAction Stop
        Write-Log "SUCCESS" "Serviço '$ServiceName' iniciado com sucesso"
        return $true
    }
    catch {
        Write-Log "ERROR" "Falha ao iniciar serviço: $($_.Exception.Message)"
        return $false
    }
}

function Test-ServiceStatus {
    param([string]$ServiceName)
    
    try {
        $service = Get-Service -Name $ServiceName -ErrorAction Stop
        Write-Log "INFO" "Status do serviço: $($service.Status)"
        return $service
    }
    catch {
        Write-Log "WARN" "Não foi possível verificar status do serviço"
        return $null
    }
}

function Main {
    Write-Log "INFO" "Syntropy CLI Windows Service Installer"
    Write-Log "INFO" "Serviço: $ServiceName"
    
    if (Install-WindowsService) {
        if ($AutoStart) {
            Start-SyntropyService -ServiceName $ServiceName
        }
        
        Test-ServiceStatus -ServiceName $ServiceName
        
        Write-Log "SUCCESS" "Instalação do serviço concluída"
    } else {
        Write-Log "ERROR" "Falha na instalação do serviço"
        exit 1
    }
}

# Run main function
if ($MyInvocation.InvocationName -ne '.') {
    Main
}