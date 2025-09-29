@echo off
REM Syntropy CLI Manager - Windows Batch Runner
REM Script simples para executar a aplicação CLI no Windows

setlocal enabledelayedexpansion

REM Configurações
set SCRIPT_DIR=%~dp0
set BUILD_DIR=%SCRIPT_DIR%build
set BINARY_NAME=syntropy.exe
set BINARY_PATH=%BUILD_DIR%\%BINARY_NAME%

REM Cores (se suportado)
set "GREEN=[92m"
set "BLUE=[94m"
set "YELLOW=[93m"
set "RED=[91m"
set "RESET=[0m"

echo.
echo %BLUE%╔══════════════════════════════════════════════════════════════╗%RESET%
echo %BLUE%║              SYNTROPY CLI MANAGER                           ║%RESET%
echo %BLUE%║                 Windows Runner                               ║%RESET%
echo %BLUE%╚══════════════════════════════════════════════════════════════╝%RESET%
echo.

REM Verificar se o binário existe
if not exist "%BINARY_PATH%" (
    echo %RED%[ERROR] Binário não encontrado: %BINARY_PATH%%RESET%
    echo %YELLOW%[INFO] Execute 'build-windows.ps1 build' primeiro para compilar.%RESET%
    echo.
    pause
    exit /b 1
)

REM Verificar argumentos
if "%~1"=="" (
    echo %BLUE%[INFO] Executando CLI sem argumentos...%RESET%
    echo %BLUE%[INFO] Use 'run-cli.bat --help' para ver opções disponíveis%RESET%
    echo.
    "%BINARY_PATH%"
) else (
    echo %BLUE%[INFO] Executando CLI com argumentos: %*%RESET%
    echo.
    "%BINARY_PATH%" %*
)

echo.
echo %GREEN%[SUCCESS] Execução concluída%RESET%
pause
