#!/bin/bash

# Syntropy CLI Manager - Universal Build Runner
# Script que detecta a plataforma e executa o script apropriado

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Detectar plataforma
OS=$(uname -s | tr '[:upper:]' '[:lower:]')

case "$OS" in
    "linux"|"darwin")
        echo "[INFO] Detected $OS platform, using bash..."
        bash "$SCRIPT_DIR/build-all.sh" "$@"
        ;;
    "mingw"*|"cygwin"*|"msys"*)
        echo "[INFO] Detected Windows environment, using PowerShell..."
        powershell -ExecutionPolicy Bypass -File "$SCRIPT_DIR/build-all.ps1" "$@"
        ;;
    *)
        echo "[WARNING] Unknown platform: $OS, trying bash..."
        bash "$SCRIPT_DIR/build-all.sh" "$@"
        ;;
esac
