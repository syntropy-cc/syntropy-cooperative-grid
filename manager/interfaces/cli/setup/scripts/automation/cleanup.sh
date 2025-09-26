#!/bin/bash

# Syntropy CLI - Cleanup e Reset Completo
# Script para cleanup completo e reset do Syntropy CLI
#
# Author: Syntropy Team
# Version: 1.0.0  
# Created: 2025-01-27

set -euo pipefail

SCRIPT_NAME="cleanup.sh"

# Configuration paths
SYNTHROPY_HOME="$HOME/.syntropy"
SYNTHROPY_BACKUP_DIR="$HOME/.syntropy-backups"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

log() {
    local level="$1"
    shift
    echo -e "${!level}[$level]${NC} $*" 
}

confirm_action() {
    local message="$1"
    read -p "$message (y/N): " confirm
    [[ "$confirm" == "y" ]] || {
        log YELLOW "Operação cancelada"
        exit 0
    }
}

cleanup_configuration() {
    log BLUE "Limpando configurações..."
    
    [[ -d "$SYNTHROPY_HOME" ]] && {
        confirm_action "Remover diretório $SYNTHROPY_HOME?"
        rm -rf "$SYNTHROPY_HOME"
        log GREEN "✓ Configurações removidas"
    } || {
        log YELLOW "Diretório $SYNTHROPY_HOME não existe"
    }
}

cleanup_backups() {
    log BLUE "Limpando backups..."
    
    [[ -d "$SYNTHROPY_BACKUP_DIR" ]] && {
        confirm_action "Remover diretório $SYNTHROPY_BACKUP_DIR?"
        rm -rf "$SYNTHROPY_BACKUP_DIR"
        log GREEN "✓ Backups removidos"
    } || {
        log YELLOW "Diretório $SYNTHROPY_BACKUP_DIR não existe"
    }
}

stop_services() {
    log BLUE "Parando services..."
    
    if command -v systemctl >/dev/null 2>&1; then
        if systemctl is-active --quiet syntropy; then
            sudo systemctl stop syntropy 2>/dev/null || log YELLOW "Failed stopping systemd service"
        fi
    fi
    
    if command -v launchctl >/dev/null 2>&1; then
        launchctl unload ~/Library/LaunchAgents/com.syntropy.* 2>/dev/null || log YELLOW "Failed stopping launchd service"
    fi
    
    log GREEN "✓ Services stopped"
}

reset_all() {
    log BLUE "Executando reset completo..."
    
    cleanup_configuration
    cleanup_backups  
    stop_services
    
    log GREEN "Reset completo finalizado"
}

main() {
    log BLUE "Syntropy CLI - Cleanup e Reset"
    
    case "${1:-cleanup}" in
        "cleanup")
            cleanup_configuration
            ;;
        "reset")
            reset_all
            ;;
        "services")
            stop_services
            ;;
        *)
            echo "Usage: $0 [cleanup|reset|services]"
            exit 1
            ;;
    esac
    
    log GREEN "Cleanup finalizado com sucesso"
}

main "$@"
