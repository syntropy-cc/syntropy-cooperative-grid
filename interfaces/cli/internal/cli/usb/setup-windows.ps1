# Script de Configuração Automática - USB Syntropy para Windows
# Execute este script como Administrador para configurar o ambiente automaticamente

param(
    [string]$WSLDistro = "Ubuntu",
    [switch]$Force,
    [switch]$SkipWSL,
    [switch]$SkipPolicy
)

$ErrorActionPreference = "Stop"

Write-Host "🚀 Configuração Automática - USB Syntropy para Windows" -ForegroundColor Green
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host ""

# Função para verificar se está executando como Administrador
function Test-Administrator {
    $currentUser = [Security.Principal.WindowsIdentity]::GetCurrent()
    $principal = New-Object Security.Principal.WindowsPrincipal($currentUser)
    return $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

# Verificar privilégios de administrador
if (-not (Test-Administrator)) {
    Write-Host "❌ ERRO: Este script deve ser executado como Administrador!" -ForegroundColor Red
    Write-Host "💡 Solução: Clique com botão direito no PowerShell e selecione 'Executar como administrador'" -ForegroundColor Yellow
    exit 1
}

Write-Host "✅ Executando como Administrador" -ForegroundColor Green
Write-Host ""

# Função para verificar se comando existe
function Test-Command($command) {
    try {
        Get-Command $command -ErrorAction Stop | Out-Null
        return $true
    } catch {
        return $false
    }
}

# Função para instalar WSL
function Install-WSL {
    Write-Host "🐧 Instalando WSL..." -ForegroundColor Yellow
    
    try {
        # Verificar se WSL já está instalado
        if (Test-Command "wsl") {
            $wslVersion = wsl --version 2>$null
            if ($LASTEXITCODE -eq 0) {
                Write-Host "✅ WSL já está instalado" -ForegroundColor Green
                return $true
            }
        }
        
        # Instalar WSL
        Write-Host "📦 Instalando WSL..." -ForegroundColor Yellow
        wsl --install --no-launch
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "✅ WSL instalado com sucesso" -ForegroundColor Green
            Write-Host "⚠️  Reinicialização pode ser necessária" -ForegroundColor Yellow
            return $true
        } else {
            Write-Host "❌ Falha na instalação do WSL" -ForegroundColor Red
            return $false
        }
    } catch {
        Write-Host "❌ Erro na instalação do WSL: $($_.Exception.Message)" -ForegroundColor Red
        return $false
    }
}

# Função para configurar distribuição WSL
function Setup-WSLDistribution($distro) {
    Write-Host "🐧 Configurando distribuição WSL: $distro" -ForegroundColor Yellow
    
    try {
        # Verificar se distribuição já está instalada
        $installedDistros = wsl --list --quiet 2>$null
        if ($installedDistros -contains $distro) {
            Write-Host "✅ Distribuição $distro já está instalada" -ForegroundColor Green
            return $true
        }
        
        # Instalar distribuição
        Write-Host "📦 Instalando $distro..." -ForegroundColor Yellow
        wsl --install -d $distro
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "✅ $distro instalado com sucesso" -ForegroundColor Green
            return $true
        } else {
            Write-Host "❌ Falha na instalação de $distro" -ForegroundColor Red
            return $false
        }
    } catch {
        Write-Host "❌ Erro na configuração de $distro: $($_.Exception.Message)" -ForegroundColor Red
        return $false
    }
}

# Função para configurar política de execução
function Set-ExecutionPolicyForSyntropy {
    Write-Host "⚙️  Configurando política de execução do PowerShell..." -ForegroundColor Yellow
    
    try {
        $currentPolicy = Get-ExecutionPolicy -Scope CurrentUser
        Write-Host "📋 Política atual (CurrentUser): $currentPolicy" -ForegroundColor Cyan
        
        if ($currentPolicy -eq "Restricted") {
            Write-Host "🔧 Alterando política para RemoteSigned..." -ForegroundColor Yellow
            Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser -Force
            Write-Host "✅ Política configurada com sucesso" -ForegroundColor Green
        } else {
            Write-Host "✅ Política já está configurada adequadamente" -ForegroundColor Green
        }
        
        return $true
    } catch {
        Write-Host "❌ Erro na configuração da política: $($_.Exception.Message)" -ForegroundColor Red
        return $false
    }
}

# Função para instalar ferramentas necessárias
function Install-RequiredTools {
    Write-Host "🛠️  Verificando ferramentas necessárias..." -ForegroundColor Yellow
    
    $tools = @(
        @{Name="PowerShell"; Command="powershell.exe"; Required=$true},
        @{Name="WSL"; Command="wsl.exe"; Required=$true},
        @{Name="DiskPart"; Command="diskpart.exe"; Required=$true}
    )
    
    $allAvailable = $true
    
    foreach ($tool in $tools) {
        if (Test-Command $tool.Command) {
            Write-Host "✅ $($tool.Name): Disponível" -ForegroundColor Green
        } else {
            Write-Host "❌ $($tool.Name): Não encontrado" -ForegroundColor Red
            if ($tool.Required) {
                $allAvailable = $false
            }
        }
    }
    
    return $allAvailable
}

# Função para configurar diretórios
function Setup-Directories {
    Write-Host "📁 Configurando diretórios..." -ForegroundColor Yellow
    
    try {
        $directories = @(
            "$env:USERPROFILE\.syntropy",
            "$env:USERPROFILE\.syntropy\cache",
            "$env:USERPROFILE\.syntropy\cache\iso",
            "$env:USERPROFILE\.syntropy\work",
            "$env:USERPROFILE\.syntropy\keys",
            "$env:USERPROFILE\.syntropy\config",
            "$env:USERPROFILE\.syntropy\scripts",
            "$env:USERPROFILE\.syntropy\backups"
        )
        
        foreach ($dir in $directories) {
            if (-not (Test-Path $dir)) {
                New-Item -ItemType Directory -Path $dir -Force | Out-Null
                Write-Host "✅ Criado: $dir" -ForegroundColor Green
            } else {
                Write-Host "✅ Existe: $dir" -ForegroundColor Green
            }
        }
        
        return $true
    } catch {
        Write-Host "❌ Erro na criação de diretórios: $($_.Exception.Message)" -ForegroundColor Red
        return $false
    }
}

# Função para testar WSL
function Test-WSLFunctionality {
    Write-Host "🧪 Testando funcionalidade do WSL..." -ForegroundColor Yellow
    
    try {
        # Teste básico
        $testResult = wsl echo "WSL funcionando!" 2>$null
        if ($LASTEXITCODE -eq 0 -and $testResult -eq "WSL funcionando!") {
            Write-Host "✅ WSL básico funcionando" -ForegroundColor Green
        } else {
            Write-Host "❌ WSL básico não está funcionando" -ForegroundColor Red
            return $false
        }
        
        # Teste de comandos necessários
        $requiredCommands = @("dd", "sgdisk", "mkfs.vfat", "mount", "umount")
        $missingCommands = @()
        
        foreach ($cmd in $requiredCommands) {
            $cmdTest = wsl bash -c "command -v $cmd" 2>$null
            if ($LASTEXITCODE -eq 0) {
                Write-Host "✅ Comando $cmd disponível" -ForegroundColor Green
            } else {
                Write-Host "❌ Comando $cmd não disponível" -ForegroundColor Red
                $missingCommands += $cmd
            }
        }
        
        if ($missingCommands.Count -gt 0) {
            Write-Host "⚠️  Comandos ausentes: $($missingCommands -join ', ')" -ForegroundColor Yellow
            Write-Host "💡 Execute no WSL: sudo apt update && sudo apt install -y gdisk dosfstools" -ForegroundColor Cyan
            return $false
        }
        
        return $true
    } catch {
        Write-Host "❌ Erro no teste do WSL: $($_.Exception.Message)" -ForegroundColor Red
        return $false
    }
}

# Função para criar script de verificação
function Create-VerificationScript {
    Write-Host "📝 Criando script de verificação..." -ForegroundColor Yellow
    
    try {
        $scriptContent = @"
# Script de Verificação - USB Syntropy para Windows
# Execute este script para verificar se tudo está configurado corretamente

Write-Host "🔍 Verificação do Ambiente USB Syntropy" -ForegroundColor Green
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan

# Verificar privilégios
`$isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
if (`$isAdmin) {
    Write-Host "✅ Executando como Administrador" -ForegroundColor Green
} else {
    Write-Host "❌ NÃO está executando como Administrador" -ForegroundColor Red
}

# Verificar WSL
try {
    `$wslVersion = wsl --version 2>`$null
    if (`$LASTEXITCODE -eq 0) {
        Write-Host "✅ WSL disponível" -ForegroundColor Green
    } else {
        Write-Host "❌ WSL não disponível" -ForegroundColor Red
    }
} catch {
    Write-Host "❌ WSL não disponível" -ForegroundColor Red
}

# Verificar política de execução
`$policy = Get-ExecutionPolicy -Scope CurrentUser
if (`$policy -eq "Restricted") {
    Write-Host "❌ Política de execução restrita: `$policy" -ForegroundColor Red
} else {
    Write-Host "✅ Política de execução adequada: `$policy" -ForegroundColor Green
}

# Verificar ferramentas
`$tools = @("powershell.exe", "wsl.exe", "diskpart.exe")
foreach (`$tool in `$tools) {
    if (Get-Command `$tool -ErrorAction SilentlyContinue) {
        Write-Host "✅ `$tool disponível" -ForegroundColor Green
    } else {
        Write-Host "❌ `$tool não disponível" -ForegroundColor Red
    }
}

# Verificar dispositivos USB
try {
    `$usbDisks = Get-Disk | Where-Object {`$_.BusType -eq 'USB'}
    Write-Host "💾 Dispositivos USB encontrados: `$(`$usbDisks.Count)" -ForegroundColor Cyan
    foreach (`$disk in `$usbDisks) {
        Write-Host "   • PHYSICALDRIVE`$(`$disk.Number) - `$(`$disk.FriendlyName)" -ForegroundColor White
    }
} catch {
    Write-Host "❌ Erro ao listar dispositivos USB" -ForegroundColor Red
}

Write-Host ""
Write-Host "🎉 Verificação concluída!" -ForegroundColor Green
"@
        
        $scriptPath = "$env:USERPROFILE\.syntropy\scripts\verify-environment.ps1"
        $scriptContent | Out-File -FilePath $scriptPath -Encoding UTF8
        Write-Host "✅ Script de verificação criado: $scriptPath" -ForegroundColor Green
        
        return $true
    } catch {
        Write-Host "❌ Erro na criação do script: $($_.Exception.Message)" -ForegroundColor Red
        return $false
    }
}

# Execução principal
Write-Host "🔧 Iniciando configuração automática..." -ForegroundColor Yellow
Write-Host ""

$success = $true

# 1. Verificar ferramentas
if (-not (Install-RequiredTools)) {
    Write-Host "❌ Ferramentas necessárias não estão disponíveis" -ForegroundColor Red
    $success = $false
}

# 2. Configurar política de execução
if (-not $SkipPolicy) {
    if (-not (Set-ExecutionPolicyForSyntropy)) {
        Write-Host "❌ Falha na configuração da política de execução" -ForegroundColor Red
        $success = $false
    }
}

# 3. Instalar WSL
if (-not $SkipWSL) {
    if (-not (Install-WSL)) {
        Write-Host "❌ Falha na instalação do WSL" -ForegroundColor Red
        $success = $false
    } else {
        # Configurar distribuição
        if (-not (Setup-WSLDistribution -distro $WSLDistro)) {
            Write-Host "❌ Falha na configuração da distribuição WSL" -ForegroundColor Red
            $success = $false
        }
    }
}

# 4. Configurar diretórios
if (-not (Setup-Directories)) {
    Write-Host "❌ Falha na configuração de diretórios" -ForegroundColor Red
    $success = $false
}

# 5. Testar WSL (se instalado)
if (-not $SkipWSL) {
    if (-not (Test-WSLFunctionality)) {
        Write-Host "⚠️  WSL instalado mas pode precisar de configuração adicional" -ForegroundColor Yellow
        Write-Host "💡 Execute no WSL: sudo apt update && sudo apt install -y gdisk dosfstools" -ForegroundColor Cyan
    }
}

# 6. Criar script de verificação
if (-not (Create-VerificationScript)) {
    Write-Host "⚠️  Falha na criação do script de verificação" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan

if ($success) {
    Write-Host "🎉 Configuração concluída com sucesso!" -ForegroundColor Green
    Write-Host ""
    Write-Host "📋 Próximos passos:" -ForegroundColor Cyan
    Write-Host "1. Execute: syntropy usb-win debug" -ForegroundColor White
    Write-Host "2. Execute: syntropy usb-win list" -ForegroundColor White
    Write-Host "3. Execute: syntropy usb-win create --node-name 'meu-no'" -ForegroundColor White
    Write-Host ""
    Write-Host "💡 Para verificar a configuração:" -ForegroundColor Cyan
    Write-Host "   . $env:USERPROFILE\.syntropy\scripts\verify-environment.ps1" -ForegroundColor White
} else {
    Write-Host "❌ Configuração concluída com erros" -ForegroundColor Red
    Write-Host ""
    Write-Host "🔧 Soluções:" -ForegroundColor Cyan
    Write-Host "1. Verifique se você está executando como Administrador" -ForegroundColor White
    Write-Host "2. Verifique sua conexão com a internet" -ForegroundColor White
    Write-Host "3. Execute o script novamente com -Force se necessário" -ForegroundColor White
    Write-Host "4. Consulte o guia: GUIA_WINDOWS.md" -ForegroundColor White
}

Write-Host ""
Write-Host "📚 Documentação: GUIA_WINDOWS.md" -ForegroundColor Yellow
