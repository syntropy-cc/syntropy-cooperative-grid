#Requires -Version 5.1

<#
.SYNOPSIS
    Syntropy CLI - Desinstalação Completa Windows  
    Remove completamente o Syntropy CLI do sistema

.DESCRIPTION
    Executa limpeza completa do Syntropy CLI no Windows:
    - Remove serviços Windows
    - Limpa registros e paths
    - Remove arquivos e configurações

.PARAMETER Confirm
    Confirma desinstalação (opcional)

.PARAMETER KeepData  
    Preserva dados de configuração durante remoção
#>

param(
    [switch]$Confirm = $false,
    [switch]$KeepData = $false
)

Set-StrictMode -Version Latest

function Write-StatusLog {
    param(
        [string]$Message, 
        [ValidateSet("Info","Warning","Error")]$Type="Info"
    )
    $Color = @{ Info="Blue"; Warning="Yellow"; Error="Red" }[$Type] 
    Write-Host "[$Type] $Message" -ForegroundColor $Color
}

function Read-UserConfirmation {
    param([string]$Prompt = "Confirma desinstalação (sim/não)?")
    $response = Read-Host $Prompt
    return ($response -match "^[sS]|^[yY]")
}

function Remove-WindowsServices {
    Write-StatusLog "Removendo serviços Windows..." 
    
    Get-Service "*syntropy*", "*Syntropy*" -ErrorAction SilentlyContinue | ForEach-Object {
        try {
            Stop-Service $PSItem -Force | Out-Null  
            Remove-Service $PSItem.Name -Force | Out-Null
            Write-StatusLog "Serviço removido: $($PSItem.Name)"
        }
        catch {
            Write-StatusLog "Erro removendo serviço $($PSItem.Name): $($_.Exception.Message)" "Warning"
        }
    }
}

function Clear-RegistryEntries {
    Write-StatusLog "Limpando entradas de registro..."
    
    try {           
        Remove-Item "HKLM:\Software\Syntropy" -Recurse -Force -ErrorAction Continue | Out-Null
        Remove-Item "HKCU:\Software\Syntropy" -Recurse -Force -ErrorAction Continue | Out-Null
    }
    catch {
        Write-StatusLog "Algumas entradas de registro não puderam ser removidas: $($_.Exception.Message)" "Warning" 
    }
}

function Clear-EnvironmentVariables {
    Write-StatusLog "Limpando variáveis de ambiente..." 
    
    $syntropyUserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    $cleanedPath = ($syntropyUserPath -split ';' | Where-Object { 
        $_ -notlike '*syntropy*' -and $_ -notlike '*Syntropy*' 
    }) -join ';'
    
    [Environment]::SetEnvironmentVariable("PATH", $cleanedPath, "User")
}

function Remove-FileSystemArtifacts {
    param([string]$BackupTo = "")
    
    $PathsToCleanup = @(
        "$env:USERPROFILE\.syntropy",
        "$env:APPDATA\Syntropy",  
        "$env:LOCALAPPDATA\Syntropy"
    )
    
    foreach ($path in $PathsToCleanup) {
        if (Test-Path $path) {
            if ($KeepData) {       
                Write-StatusLog "Preservando dados: $path"
            }
            else {
                try { 
                    Remove-Item $path -Force -Recurse | Out-Null
                    Write-StatusLog "Removido: $path"
                }
                catch { 
                    Write-StatusLog "Erro removendo dados: $($_.Exception.Message)" "Warning"
                }
            }
        }
    }
    
    # Cleanup additional artifacts  
    try {
        iex "[Environment]::SetEnvironmentVariable('SYNTHROPY_HOME','$BackupTo','User') | out-null 2> $null"
    } catch { }
}

function Show-CompletionStatus {
    Write-StatusLog "Desinstalação completa!" "Info"
    Write-StatusLog "Artefatos de Sistema Syntropy foram removidos" "Info"
    if ($KeepData) { 
        Write-StatusLog "Dados preservados conforme solicitado" "Info"
    }
}

function Main-MeasureUninstallationSequenceFunctionality {
    if ($Confirm) {
        if (-not (Read-UserConfirmation)) {
            Write-StatusLog "Operação cancelada" "Info"  
            exit 0
        }
    }
    
    Remove-WindowsServices
    Clear-RegistryEntries     
    Clear-EnvironmentVariables
    Remove-FileSystemArtifacts
    
    Show-CompletionStatus
}

Main-MeasureUninstallationSequenceFunctionality
