@echo off
REM Syntropy CLI Manager - Examples Runner
REM Script para executar exemplos de uso da CLI

setlocal enabledelayedexpansion

REM Configurações
set SCRIPT_DIR=%~dp0
set BUILD_DIR=%SCRIPT_DIR%..\build
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
echo %CYAN%║                 Examples Runner                             ║%RESET%
echo %CYAN%╚══════════════════════════════════════════════════════════════╝%RESET%
echo.

REM Verificar se o binário existe
if not exist "%BINARY_PATH%" (
    echo %RED%[ERROR] Binário não encontrado: %BINARY_PATH%%RESET%
    echo %YELLOW%[INFO] Execute 'build-windows.ps1 build' primeiro para compilar.%RESET%
    echo.
    pause
    exit /b 1
)

echo %GREEN%[SUCCESS] Binário encontrado: %BINARY_PATH%%RESET%
echo.

REM Menu de exemplos
:menu
echo %CYAN%=== EXEMPLOS DISPONÍVEIS ===%RESET%
echo %BLUE%1. Mostrar ajuda geral%RESET%
echo %BLUE%2. Mostrar versão%RESET%
echo %BLUE%3. Ajuda do comando setup%RESET%
echo %BLUE%4. Validar ambiente (sem mudanças)%RESET%
echo %BLUE%5. Verificar status do setup%RESET%
echo %BLUE%6. Executar setup completo%RESET%
echo %BLUE%7. Executar setup com força%RESET%
echo %BLUE%8. Reset configuração%RESET%
echo %BLUE%9. Executar todos os exemplos%RESET%
echo %BLUE%0. Sair%RESET%
echo.
set /p CHOICE="Escolha um exemplo (0-9): "

if "%CHOICE%"=="1" goto :help_general
if "%CHOICE%"=="2" goto :version
if "%CHOICE%"=="3" goto :help_setup
if "%CHOICE%"=="4" goto :validate
if "%CHOICE%"=="5" goto :status
if "%CHOICE%"=="6" goto :setup_run
if "%CHOICE%"=="7" goto :setup_force
if "%CHOICE%"=="8" goto :reset
if "%CHOICE%"=="9" goto :run_all
if "%CHOICE%"=="0" goto :exit
echo %RED%[ERROR] Opção inválida%RESET%
goto :menu

:help_general
echo.
echo %BLUE%[INFO] Executando: %BINARY_NAME% --help%RESET%
echo.
"%BINARY_PATH%" --help
echo.
pause
goto :menu

:version
echo.
echo %BLUE%[INFO] Executando: %BINARY_NAME% --version%RESET%
echo.
"%BINARY_PATH%" --version
echo.
pause
goto :menu

:help_setup
echo.
echo %BLUE%[INFO] Executando: %BINARY_NAME% setup --help%RESET%
echo.
"%BINARY_PATH%" setup --help
echo.
pause
goto :menu

:validate
echo.
echo %BLUE%[INFO] Executando: %BINARY_NAME% setup validate%RESET%
echo.
"%BINARY_PATH%" setup validate
echo.
pause
goto :menu

:status
echo.
echo %BLUE%[INFO] Executando: %BINARY_NAME% setup status%RESET%
echo.
"%BINARY_PATH%" setup status
echo.
pause
goto :menu

:setup_run
echo.
echo %BLUE%[INFO] Executando: %BINARY_NAME% setup run%RESET%
echo.
"%BINARY_PATH%" setup run
echo.
pause
goto :menu

:setup_force
echo.
echo %BLUE%[INFO] Executando: %BINARY_NAME% setup run --force%RESET%
echo.
"%BINARY_PATH%" setup run --force
echo.
pause
goto :menu

:reset
echo.
echo %YELLOW%[WARNING] Este comando irá resetar toda a configuração!%RESET%
echo %YELLOW%[WARNING] Execute: %BINARY_NAME% setup reset --force%RESET%
echo.
set /p CONFIRM="Tem certeza? (S/N): "
if /i "!CONFIRM!"=="S" (
    echo %BLUE%[INFO] Executando: %BINARY_NAME% setup reset --force%RESET%
    echo.
    "%BINARY_PATH%" setup reset --force
) else (
    echo %BLUE%[INFO] Reset cancelado%RESET%
)
echo.
pause
goto :menu

:run_all
echo.
echo %CYAN%=== EXECUTANDO TODOS OS EXEMPLOS ===%RESET%
echo.

echo %BLUE%[1/8] Ajuda geral...%RESET%
"%BINARY_PATH%" --help
echo.
echo %BLUE%[2/8] Versão...%RESET%
"%BINARY_PATH%" --version
echo.
echo %BLUE%[3/8] Ajuda do setup...%RESET%
"%BINARY_PATH%" setup --help
echo.
echo %BLUE%[4/8] Validar ambiente...%RESET%
"%BINARY_PATH%" setup validate
echo.
echo %BLUE%[5/8] Status do setup...%RESET%
"%BINARY_PATH%" setup status
echo.
echo %BLUE%[6/8] Executar setup...%RESET%
"%BINARY_PATH%" setup run
echo.
echo %BLUE%[7/8] Executar setup com força...%RESET%
"%BINARY_PATH%" setup run --force
echo.
echo %BLUE%[8/8] Status final...%RESET%
"%BINARY_PATH%" setup status
echo.

echo %GREEN%[SUCCESS] Todos os exemplos executados!%RESET%
echo.
pause
goto :menu

:exit
echo.
echo %GREEN%[SUCCESS] Obrigado por usar o Syntropy CLI Manager!%RESET%
echo.
pause
exit /b 0
