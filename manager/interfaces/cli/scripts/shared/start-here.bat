@echo off
REM Syntropy CLI Manager - Start Here Script
REM Script principal para iniciar o workflow no Windows

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
set "MAGENTA=[95m"
set "RESET=[0m"

echo.
echo %MAGENTA%╔══════════════════════════════════════════════════════════════╗%RESET%
echo %MAGENTA%║              SYNTROPY CLI MANAGER                           ║%RESET%
echo %MAGENTA%║                 Windows Workflow                            ║%RESET%
echo %MAGENTA%║                                                              ║%RESET%
echo %MAGENTA%║  Bem-vindo ao workflow de desenvolvimento do Syntropy CLI!   ║%RESET%
echo %MAGENTA%╚══════════════════════════════════════════════════════════════╝%RESET%
echo.

REM Verificar se Go está instalado
echo %BLUE%[INFO] Verificando Go...%RESET%
go version >nul 2>&1
if errorlevel 1 (
    echo %RED%[ERROR] Go não está instalado ou não está no PATH%RESET%
    echo %YELLOW%[INFO] Por favor, instale Go 1.22.5 ou superior%RESET%
    echo %YELLOW%[INFO] Download: https://golang.org/dl/%RESET%
    echo.
    echo %BLUE%[INFO] Ou execute o script de setup automático:%RESET%
    echo %BLUE%[INFO] .\scripts\setup-environment.ps1 install%RESET%
    echo.
    pause
    exit /b 1
)

for /f "tokens=3" %%i in ('go version') do set GO_VERSION=%%i
echo %GREEN%[SUCCESS] Go %GO_VERSION% encontrado%RESET%

REM Verificar se main.go existe
if not exist "%SCRIPT_DIR%main.go" (
    echo %RED%[ERROR] main.go não encontrado em %SCRIPT_DIR%%RESET%
    echo %YELLOW%[INFO] Certifique-se de estar no diretório correto%RESET%
    pause
    exit /b 1
)

echo %GREEN%[SUCCESS] Estrutura do projeto verificada%RESET%
echo.

REM Menu principal
:main_menu
echo %CYAN%=== WORKFLOW PRINCIPAL ===%RESET%
echo %BLUE%1. 🚀 Início Rápido (Recomendado para iniciantes)%RESET%
echo %BLUE%2. 🔧 Setup do Ambiente%RESET%
echo %BLUE%3. 🏗️  Build e Execução%RESET%
echo %BLUE%4. 🧪 Desenvolvimento Completo%RESET%
echo %BLUE%5. 🤖 Automação e CI/CD%RESET%
echo %BLUE%6. 📚 Exemplos de Uso%RESET%
echo %BLUE%7. 📖 Documentação%RESET%
echo %BLUE%8. 🛠️  Solução de Problemas%RESET%
echo %BLUE%0. ❌ Sair%RESET%
echo.
set /p CHOICE="Escolha uma opção (0-8): "

if "%CHOICE%"=="1" goto :quick_start
if "%CHOICE%"=="2" goto :setup_environment
if "%CHOICE%"=="3" goto :build_run
if "%CHOICE%"=="4" goto :development
if "%CHOICE%"=="5" goto :automation
if "%CHOICE%"=="6" goto :examples
if "%CHOICE%"=="7" goto :documentation
if "%CHOICE%"=="8" goto :troubleshooting
if "%CHOICE%"=="0" goto :exit
echo %RED%[ERROR] Opção inválida%RESET%
goto :main_menu

:quick_start
echo.
echo %CYAN%=== INÍCIO RÁPIDO ===%RESET%
echo %BLUE%[INFO] Executando script de início rápido...%RESET%
echo.
if exist "%SCRIPT_DIR%quick-start.bat" (
    call "%SCRIPT_DIR%quick-start.bat"
) else (
    echo %RED%[ERROR] quick-start.bat não encontrado%RESET%
)
echo.
pause
goto :main_menu

:setup_environment
echo.
echo %CYAN%=== SETUP DO AMBIENTE ===%RESET%
echo %BLUE%[INFO] Configurando ambiente de desenvolvimento...%RESET%
echo.
if exist "%SCRIPT_DIR%scripts\setup-environment.ps1" (
    powershell -ExecutionPolicy Bypass -File "%SCRIPT_DIR%scripts\setup-environment.ps1" check
    echo.
    set /p SETUP_CHOICE="Deseja instalar/configurar automaticamente? (S/N): "
    if /i "!SETUP_CHOICE!"=="S" (
        powershell -ExecutionPolicy Bypass -File "%SCRIPT_DIR%scripts\setup-environment.ps1" install
        powershell -ExecutionPolicy Bypass -File "%SCRIPT_DIR%scripts\setup-environment.ps1" configure
        powershell -ExecutionPolicy Bypass -File "%SCRIPT_DIR%scripts\setup-environment.ps1" verify
    )
) else (
    echo %RED%[ERROR] setup-environment.ps1 não encontrado%RESET%
)
echo.
pause
goto :main_menu

:build_run
echo.
echo %CYAN%=== BUILD E EXECUÇÃO ===%RESET%
echo %BLUE%[INFO] Compilando e executando aplicação...%RESET%
echo.
if exist "%SCRIPT_DIR%build-windows.ps1" (
    powershell -ExecutionPolicy Bypass -File "%SCRIPT_DIR%build-windows.ps1" build
    echo.
    set /p RUN_CHOICE="Deseja executar a aplicação? (S/N): "
    if /i "!RUN_CHOICE!"=="S" (
        echo %BLUE%[INFO] Executando aplicação...%RESET%
        powershell -ExecutionPolicy Bypass -File "%SCRIPT_DIR%build-windows.ps1" run
    )
) else (
    echo %RED%[ERROR] build-windows.ps1 não encontrado%RESET%
)
echo.
pause
goto :main_menu

:development
echo.
echo %CYAN%=== DESENVOLVIMENTO COMPLETO ===%RESET%
echo %BLUE%[INFO] Executando workflow de desenvolvimento...%RESET%
echo.
if exist "%SCRIPT_DIR%dev-workflow.ps1" (
    powershell -ExecutionPolicy Bypass -File "%SCRIPT_DIR%dev-workflow.ps1" dev
) else (
    echo %RED%[ERROR] dev-workflow.ps1 não encontrado%RESET%
)
echo.
pause
goto :main_menu

:automation
echo.
echo %CYAN%=== AUTOMAÇÃO E CI/CD ===%RESET%
echo %BLUE%[INFO] Executando workflow de automação...%RESET%
echo.
if exist "%SCRIPT_DIR%automation-workflow.ps1" (
    powershell -ExecutionPolicy Bypass -File "%SCRIPT_DIR%automation-workflow.ps1" full
) else (
    echo %RED%[ERROR] automation-workflow.ps1 não encontrado%RESET%
)
echo.
pause
goto :main_menu

:examples
echo.
echo %CYAN%=== EXEMPLOS DE USO ===%RESET%
echo %BLUE%[INFO] Executando exemplos da CLI...%RESET%
echo.
if exist "%SCRIPT_DIR%scripts\run-examples.bat" (
    call "%SCRIPT_DIR%scripts\run-examples.bat"
) else (
    echo %RED%[ERROR] run-examples.bat não encontrado%RESET%
)
echo.
pause
goto :main_menu

:documentation
echo.
echo %CYAN%=== DOCUMENTAÇÃO ===%RESET%
echo %BLUE%[INFO] Abrindo documentação...%RESET%
echo.
echo %BLUE%📚 Documentação disponível:%RESET%
echo %BLUE%  - README_WINDOWS.md (Início rápido)%RESET%
echo %BLUE%  - WINDOWS_WORKFLOW.md (Documentação completa)%RESET%
echo %BLUE%  - BUILD_AND_TEST.md (Build e testes)%RESET%
echo %BLUE%  - GUIDE.md (Guia completo)%RESET%
echo.
set /p DOC_CHOICE="Deseja abrir a documentação? (S/N): "
if /i "!DOC_CHOICE!"=="S" (
    if exist "%SCRIPT_DIR%README_WINDOWS.md" (
        start notepad "%SCRIPT_DIR%README_WINDOWS.md"
    )
    if exist "%SCRIPT_DIR%WINDOWS_WORKFLOW.md" (
        start notepad "%SCRIPT_DIR%WINDOWS_WORKFLOW.md"
    )
)
echo.
pause
goto :main_menu

:troubleshooting
echo.
echo %CYAN%=== SOLUÇÃO DE PROBLEMAS ===%RESET%
echo %BLUE%[INFO] Verificando problemas comuns...%RESET%
echo.

REM Verificar Go
echo %BLUE%[1/5] Verificando Go...%RESET%
go version >nul 2>&1
if errorlevel 1 (
    echo %RED%[ERROR] Go não encontrado%RESET%
    echo %YELLOW%[SOLUÇÃO] Instale Go: https://golang.org/dl/%RESET%
) else (
    echo %GREEN%[OK] Go encontrado%RESET%
)

REM Verificar Git
echo %BLUE%[2/5] Verificando Git...%RESET%
git --version >nul 2>&1
if errorlevel 1 (
    echo %RED%[ERROR] Git não encontrado%RESET%
    echo %YELLOW%[SOLUÇÃO] Instale Git: https://git-scm.com/download/win%RESET%
) else (
    echo %GREEN%[OK] Git encontrado%RESET%
)

REM Verificar main.go
echo %BLUE%[3/5] Verificando main.go...%RESET%
if exist "%SCRIPT_DIR%main.go" (
    echo %GREEN%[OK] main.go encontrado%RESET%
) else (
    echo %RED%[ERROR] main.go não encontrado%RESET%
    echo %YELLOW%[SOLUÇÃO] Certifique-se de estar no diretório correto%RESET%
)

REM Verificar binário
echo %BLUE%[4/5] Verificando binário...%RESET%
if exist "%BINARY_PATH%" (
    echo %GREEN%[OK] Binário encontrado%RESET%
) else (
    echo %YELLOW%[WARNING] Binário não encontrado%RESET%
    echo %YELLOW%[SOLUÇÃO] Execute build primeiro%RESET%
)

REM Verificar scripts
echo %BLUE%[5/5] Verificando scripts...%RESET%
set SCRIPT_COUNT=0
if exist "%SCRIPT_DIR%build-windows.ps1" set /a SCRIPT_COUNT+=1
if exist "%SCRIPT_DIR%dev-workflow.ps1" set /a SCRIPT_COUNT+=1
if exist "%SCRIPT_DIR%automation-workflow.ps1" set /a SCRIPT_COUNT+=1
if exist "%SCRIPT_DIR%quick-start.bat" set /a SCRIPT_COUNT+=1

if %SCRIPT_COUNT% geq 3 (
    echo %GREEN%[OK] Scripts encontrados (%SCRIPT_COUNT%/4)%RESET%
) else (
    echo %RED%[ERROR] Scripts não encontrados%RESET%
    echo %YELLOW%[SOLUÇÃO] Verifique se todos os arquivos foram criados%RESET%
)

echo.
echo %BLUE%💡 Dicas de solução de problemas:%RESET%
echo %BLUE%  - Verifique se Go está no PATH%RESET%
echo %BLUE%  - Execute como administrador se necessário%RESET%
echo %BLUE%  - Verifique logs em logs/ para detalhes%RESET%
echo %BLUE%  - Use PowerShell se batch não funcionar%RESET%
echo.
pause
goto :main_menu

:exit
echo.
echo %GREEN%[SUCCESS] Obrigado por usar o Syntropy CLI Manager!%RESET%
echo %BLUE%[INFO] Para mais informações, consulte a documentação.%RESET%
echo.
pause
exit /b 0
