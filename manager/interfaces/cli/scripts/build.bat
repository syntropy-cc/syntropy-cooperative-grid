@echo off
REM Syntropy CLI Manager - Universal Build Runner
REM Script que detecta a plataforma e executa o script apropriado

setlocal enabledelayedexpansion

REM Detectar se estamos no Windows
if "%OS%"=="Windows_NT" (
    REM Windows - usar PowerShell
    echo [INFO] Detected Windows platform, using PowerShell...
    powershell -ExecutionPolicy Bypass -File "%~dp0build-all.ps1" %*
) else (
    REM Linux/macOS - usar bash
    echo [INFO] Detected Unix-like platform, using bash...
    bash "%~dp0build-all.sh" %*
)

endlocal
