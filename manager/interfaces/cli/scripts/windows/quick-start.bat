@echo off
REM Syntropy CLI Manager - Quick Start Script
REM Script para setup rápido e execução da aplicação CLI

setlocal enabledelayedexpansion

REM Configurações
set SCRIPT_DIR=%~dp0
set BUILD_DIR=%SCRIPT_DIR%build
set BINARY_NAME=syntropy.exe
set BINARY_PATH=%BUILD_DIR%\%BINARY_NAME%

REM Cores
set "GREEN=[92m"
set "BLUE=[94m"
set "YELLOW=[93m"
set "RED=[91m"
set "CYAN=[96m"
set "RESET=[0m"

echo.
echo %CYAN%╔══════════════════════════════════════════════════════════════╗%RESET%
echo %CYAN%║              SYNTROPY CLI MANAGER                           ║%RESET%
echo %CYAN%║                 Quick Start                                 ║%RESET%
echo %CYAN%╚══════════════════════════════════════════════════════════════╝%RESET%
echo.

REM Verificar se Go está instalado
echo %BLUE%[INFO] Verificando Go...%RESET%
go version >nul 2>&1
if errorlevel 1 (
    echo %RED%[ERROR] Go não está instalado ou não está no PATH%RESET%
    echo %YELLOW%[INFO] Por favor, instale Go 1.22.5 ou superior%RESET%
    echo %YELLOW%[INFO] Download: https://golang.org/dl/%RESET%
    echo.
    pause
    exit /b 1
)

for /f "tokens=3" %%i in ('go version') do set GO_VERSION=%%i
echo %GREEN%[SUCCESS] Go %GO_VERSION% encontrado%RESET%

REM Verificar se main.go existe
if not exist "%SCRIPT_DIR%main.go" (
    echo %RED%[ERROR] main.go não encontrado em %SCRIPT_DIR%%RESET%
    pause
    exit /b 1
)

echo %GREEN%[SUCCESS] Estrutura do projeto verificada%RESET%

REM Verificar se já existe build
if exist "%BINARY_PATH%" (
    echo %YELLOW%[INFO] Binário já existe. Deseja recompilar? (S/N)%RESET%
    set /p REBUILD=
    if /i "!REBUILD!"=="S" (
        echo %BLUE%[INFO] Recompilando...%RESET%
        goto :build
    ) else (
        echo %BLUE%[INFO] Usando binário existente%RESET%
        goto :run
    )
) else (
    echo %BLUE%[INFO] Binário não encontrado. Compilando...%RESET%
    goto :build
)

:build
echo.
echo %CYAN%=== COMPILANDO APLICAÇÃO ===%RESET%

REM Criar diretório de build
if not exist "%BUILD_DIR%" mkdir "%BUILD_DIR%"

REM Baixar dependências
echo %BLUE%[INFO] Baixando dependências...%RESET%
cd /d "%SCRIPT_DIR%"
go mod download
if errorlevel 1 (
    echo %RED%[ERROR] Falha ao baixar dependências%RESET%
    pause
    exit /b 1
)

REM Organizar dependências
echo %BLUE%[INFO] Organizando dependências...%RESET%
go mod tidy
if errorlevel 1 (
    echo %RED%[ERROR] Falha ao organizar dependências%RESET%
    pause
    exit /b 1
)

REM Compilar
echo %BLUE%[INFO] Compilando aplicação...%RESET%
for /f %%i in ('powershell -command "Get-Date -Format 'yyyyMMdd-HHmmss'"') do set VERSION=%%i
for /f %%i in ('git rev-parse --short HEAD 2^>nul ^|^| echo unknown') do set GIT_COMMIT=%%i
for /f %%i in ('powershell -command "Get-Date -Format 'yyyy-MM-ddTHH:mm:ssZ'"') do set BUILD_TIME=%%i

set BUILD_FLAGS=-ldflags "-X main.version=%VERSION% -X main.buildTime=%BUILD_TIME% -X main.gitCommit=%GIT_COMMIT%"

go build %BUILD_FLAGS% -o "%BINARY_PATH%" main.go
if errorlevel 1 (
    echo %RED%[ERROR] Falha na compilação%RESET%
    pause
    exit /b 1
)

echo %GREEN%[SUCCESS] Compilação concluída: %BINARY_PATH%%RESET%

REM Verificar tamanho do binário
for %%i in ("%BINARY_PATH%") do set BINARY_SIZE=%%~zi
set /a BINARY_SIZE_KB=%BINARY_SIZE%/1024
echo %BLUE%[INFO] Tamanho: %BINARY_SIZE_KB% KB%RESET%

:run
echo.
echo %CYAN%=== EXECUTANDO APLICAÇÃO ===%RESET%

REM Executar testes básicos
echo %BLUE%[INFO] Executando testes básicos...%RESET%
"%BINARY_PATH%" --version >nul 2>&1
if errorlevel 1 (
    echo %YELLOW%[WARNING] Teste de versão falhou (pode ser normal)%RESET%
) else (
    echo %GREEN%[SUCCESS] Teste de versão passou%RESET%
)

"%BINARY_PATH%" --help >nul 2>&1
if errorlevel 1 (
    echo %YELLOW%[WARNING] Teste de ajuda falhou (pode ser normal)%RESET%
) else (
    echo %GREEN%[SUCCESS] Teste de ajuda passou%RESET%
)

echo.
echo %GREEN%[SUCCESS] Setup concluído com sucesso!%RESET%
echo.
echo %BLUE%📋 Informações do Build:%RESET%
echo %BLUE%   - Binário: %BINARY_PATH%%RESET%
echo %BLUE%   - Versão: %VERSION%%RESET%
echo %BLUE%   - Git Commit: %GIT_COMMIT%%RESET%
echo %BLUE%   - Tamanho: %BINARY_SIZE_KB% KB%RESET%
echo.

REM Menu de opções
:menu
echo %CYAN%=== OPÇÕES DISPONÍVEIS ===%RESET%
echo %BLUE%1. Executar CLI (sem argumentos)%RESET%
echo %BLUE%2. Executar CLI com argumentos%RESET%
echo %BLUE%3. Mostrar ajuda%RESET%
echo %BLUE%4. Executar setup%RESET%
echo %BLUE%5. Verificar status%RESET%
echo %BLUE%6. Sair%RESET%
echo.
set /p CHOICE="Escolha uma opção (1-6): "

if "%CHOICE%"=="1" goto :run_cli
if "%CHOICE%"=="2" goto :run_with_args
if "%CHOICE%"=="3" goto :show_help
if "%CHOICE%"=="4" goto :run_setup
if "%CHOICE%"=="5" goto :check_status
if "%CHOICE%"=="6" goto :exit
echo %RED%[ERROR] Opção inválida%RESET%
goto :menu

:run_cli
echo.
echo %BLUE%[INFO] Executando CLI...%RESET%
echo.
"%BINARY_PATH%"
echo.
pause
goto :menu

:run_with_args
echo.
set /p CLI_ARGS="Digite os argumentos para o CLI: "
echo.
echo %BLUE%[INFO] Executando CLI com argumentos: %CLI_ARGS%%RESET%
echo.
"%BINARY_PATH%" %CLI_ARGS%
echo.
pause
goto :menu

:show_help
echo.
echo %BLUE%[INFO] Mostrando ajuda...%RESET%
echo.
"%BINARY_PATH%" --help
echo.
pause
goto :menu

:run_setup
echo.
echo %BLUE%[INFO] Executando setup...%RESET%
echo.
"%BINARY_PATH%" setup run --force
echo.
pause
goto :menu

:check_status
echo.
echo %BLUE%[INFO] Verificando status...%RESET%
echo.
"%BINARY_PATH%" setup status
echo.
pause
goto :menu

:exit
echo.
echo %GREEN%[SUCCESS] Obrigado por usar o Syntropy CLI Manager!%RESET%
echo.
pause
exit /b 0
