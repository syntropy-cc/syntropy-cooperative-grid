@echo off
REM Syntropy CLI Manager - Main Build Script for Windows
REM Script principal para compilar a aplicação CLI no Windows

setlocal enabledelayedexpansion

REM Configuration
set SCRIPT_DIR=%~dp0
set BUILD_SCRIPT=%SCRIPT_DIR%scripts\shared\build-and-test.bat

REM Colors
set "GREEN=[92m"
set "BLUE=[94m"
set "YELLOW=[93m"
set "RED=[91m"
set "CYAN=[96m"
set "RESET=[0m"

echo.
echo %CYAN%╔══════════════════════════════════════════════════════════════╗%RESET%
echo %CYAN%║              SYNTROPY CLI MANAGER                           ║%RESET%
echo %CYAN%║                Main Build Script                            ║%RESET%
echo %CYAN%╚══════════════════════════════════════════════════════════════╝%RESET%
echo.

REM Check if build script exists
if not exist "%BUILD_SCRIPT%" (
    echo %RED%[ERROR] Build script not found: %BUILD_SCRIPT%%RESET%
    pause
    exit /b 1
)

echo %BLUE%[INFO] Starting Syntropy CLI Manager build...%RESET%
echo %BLUE%[INFO] Using build script: %BUILD_SCRIPT%%RESET%

REM Execute the build script
call "%BUILD_SCRIPT%"

