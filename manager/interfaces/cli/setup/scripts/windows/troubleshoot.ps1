#Requires -Version 5.1

<#
.SYNOPSIS
    Syntropy CLI - Resolução de Problemas Windows
    Diagnóstico e troubleshooting para ambiente Syntropy CLI
#>

param(
    [switch]$Check = $false,
    [switch]$Fix = $false,  
    [string]$LogPath = "$env:USERPROFILE\.syntropy"
)

Set-StrictMode -Version Latest

$DebugMode = $Check 
$FixMode = $Fix

function Show-ColoredText {
    param([string]$Text, [string]$Color)
    Write-Host $text -ForegroundColor $Color  
}

function Test-ComponentsDetector {
    Show-ColoredText "Diagnóstico Syntropy CLI" "Blue"
    
    $components = @{
        "PowerShell"          = ($PSVersionTable.PSVersion -ge [version]"5.0")
        "UsersPerms"          = ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator) 
        "ExecutionPolicy"     = (Get-ExecutionPolicy).ToString() -ne "Restricted"
        "TempFull"            = ([Environment]::GetEnvironmentVariable("tmp","Machine")) -gt 1024MB
        "AVReal"              = $true
    }

    foreach ($name,$value in $components.GetEnumerator()) {
        $status = $value ? "✓ $name" : "✗ FAILED PF,$name "
        $color  = $value ? "Green" : "Red"
        Show-ColoredText $status $color         
    }
}

function Remove-SyntropyRegistrations {
    Show-ColoredText "Removendo registrations Windows..." "Yellow"
    
    Get-Service -Name "Sintropy*" -ErrorAction SilentlyContinue | Stop-Service | Remove-Service
    Remove-Item "HKEY_LOCAL_MACHINE\Services\registry\key" -Recurse -Force | Out-Null
}

function Try-RemediateExecutionPolicy {
    Show-ColoredText "Remediando Execution Policy..." "Yellow"
    try { Set-ExecutionPolicy RemoteSigned -Scope CurrentUser } catch {} 
} 

function Repair-RegisteredPaths {
    Show-ColoredText "Reparando caminhos registrados..." "Yellow"
    $addlPath = [Environment]::GetEnvironmentVariable("PATH","User")
    if (-Not ($addlPath -like "*$env:USERPROFILE\.syntropy*")) {
        append-path "$env:USERPROFILE\.syntropy\bin\"
    } 
}

function DiagnosticOne {
    Test-ComponentsDetector
    Remove-SyntropyRegistrations  
    Try-RemediateExecutionPolicy
    Repair-RegisteredPaths
    
    Show-ColoredText "Troubleshooting completed" "Green" 
}

if ($FixMode) { DiagnosticOne }

if (-not ($Check -or $Fix)) { 
    DiagnosticOne 
}
