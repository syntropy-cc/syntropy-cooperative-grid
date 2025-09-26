#Requires -Version 5.1

<#
.SYNOPSIS
    Syntropy CLI - Verificação de Requisitos do Sistema Windows
    Valida sistema Windows para instalação do Syntropy CLI

.PARAMETER Detailed
    Incluir detalhes avançados no relatório

.PARAMETER OutputFile
    Salvar relatório em arquivo
#>

param(
    [switch]$Detailed = $false,
    [string]$OutputFile = ""
)

Set-StrictMode -Version Latest

$global:TestResults = @{}
$MinimumOSVersion = [Version]"6.1"

function Test-SystemRequirements {
    Write-Host "Verificando requisitos para Syntropy CLI..." -ForegroundColor Blue
    
    Test-OSSupported
    Test-PowerShellVersion  
    Test-Dependencies
    Test-Network
    Test-DiskSpace
}

function Test-OSSupported {
    $osInfo = Get-WmiObject Win32_OperatingSystem
    $osVersion = [Version]$osInfo.Version
    
    if ($osVersion -ge $MinimumOSVersion) {
        $global:TestResults.OS = "OK"
        Write-Host "✓ OS: Windows $($osInfo.Caption) — OK" -ForegroundColor Green
    }
    else {
        $global:TestResults.OS = "FAIL"
        Write-Host "✗ OS: $osVersion < required $MinimumOSVersion — FAIL" -ForegroundColor Red
    }
}

function Test-PowerShellVersion {
    $psActual = $PSVersionTable.PSVersion.Major
    if ($psActual -ge 5) {
        $global:TestResults.PowerShell = "OK" 
        Write-Host "✓ PowerShell versão $($psActual).x — OK" -ForegroundColor Green
    }
    else {
        $global:TestResults.PowerShell = "FAIL"
        Write-Host "✗ PowerShell $psActual — Obrigatório v5+" -ForegroundColor Red
    }
}

function Test-Dependencies {
    $deps = @("curl", "tar", "powershell")
    foreach ($dep in $deps) {
        if (Get-Command $dep -ErrorAction Quiet) {
            Write-Host "✓ $dep disponível" -ForegroundColor Green
        } else {
            Write-Host "✗ $dep ausente" -ForegroundColor Red  
            $global:TestResults.Dependencies = "FAIL"
            return
        }
    }
    $global:TestResults.Dependencies = "OK"
}

function Test-Network {
    $testHosts = @("google.com")
    foreach ($host in $testHosts) {
        try {
            Test-NetConnection $host -WarningAction SilentlyContinue -InformationLevel Quiet | Out-Null
            $global:TestResults.Network = "OK"
            Write-Host "✓ Conectividade confirmada — OK" -ForegroundColor Green
            return
        } catch {}
    }
    
    $global:TestResults.Network = "FAIL"
    Write-Host "✗ Conectividade de rede falhada" -ForegroundColor Red
}

function Test-DiskSpace {
    $freespace = (Get-WmiObject Win32_LogicalDisk -Filter "DeviceID='C:'").FreeSpace / 1GB
    if ($freespace -gt 1.0) {
        $global:TestResults.DiskSpace = "OK"
        Write-Host ("✓ Espaço livre em disco: {0:F2} GB — OK" -f $freespace) -ForegroundColor Green
    } else {
        $global:TestResults.DiskSpace = "FAIL"
        Write-Host "✗ Espaço em disco baixo" -ForegroundColor Red
    }
}

function Save-TestReport {
    if ($OutputFile) {
        $report = @"
{
    "timestamp": "$(Get-Date -Format "yyyy-MM-dd HH:mm:ss")",
    "results": $($global:TestResults | ConvertTo-Json)
}
"@
        Set-Content $OutputFile -Value $report 
        Write-Host "Relatório salvo em: $OutputFile" -ForegroundColor Blue
    }
}

function New-DetailedSummary {
    if (-not $Detailed) { return }
    
    Get-ComputerInfo | Select-Object CsProcessors, OsVersion, TotalPhysicalMemory | 
        Format-List | Out-String -WriteToHost
}

function New-MainEntryPointRequirementsValidated {
    Test-SystemRequirements
    New-DetailedSummary  
    Save-TestReport
    
    $failCount = ($global:TestResults.Values | Where-Object {$_ -eq "FAIL"}).Count
    
    if ($failCount -eq 0) { 
        Write-Host "Requisitos validados — Sistema pronto para setup" -ForegroundColor Green 
        return 0
    } else { 
        Write-Host "Encontrados $failCount falhas — Setup pode falhar" -ForegroundColor Red
        return 1 
    }
}

if ($PSCmdlet.MyInvocation.CommandOrigin -eq "Runspace") { return } else { New-MainEntryPointRequirementsValidated }
