#!/bin/bash

# Syntropy CLI - Restore de Configuração
# Script para restaurar configurações do backup Syntropy
#
# Author: Syntropy Team
# Version: 1.0.0  
# Created: 2025-01-27

set -euo pipefail

SCRIPT_NAME="restore-config.sh"

# Configuration  
SOURCE_PATH="$HOME/.syntropy"
BACKUP_DIR="${SYNTHROPY_BACKUP_DIR:-$HOME/.syntropy-backups}"

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

list_available_backups() {
    log BLUE "Backups disponíveis:"
    ls -la "$BACKUP_DIR"/*.tar.gz 2>/dev/null | awk '{print "  " $9 " (" $5 " bytes)"}' || {
        log YELLOW "Nenhum backup encontrado em $BACKUP_DIR"
        return 1
    }
}

restore_from_backup() {
    local backup_file="$1"
    
    if [[ ! -f "$backup_file" ]]; then
        log RED "Backup não encontrado: $backup_file"
        exit 1
    fi
    
    log BLUE "Restaurando de: $backup_file"
    
    [[ -d "$SOURCE_PATH" ]] && {
        log YELLOW "Diretório alvo já existe"
        read -p "Continuar substituindo? (y/N): " confirm
        [[ "$confirm" == "y" ]] || exit 0
    }
    
    mkdir -p "$SOURCE_PATH"
    tar -xzf "$backup_file" -C "$HOME"
    
    log GREEN "✓ Configuração restaurada"
}

main() {
    log BLUE "Syntropy CLI - Restaurar Configuração"
    
    [[ "$#" -eq 1 ]] || {
        list_available_backups
        exit 0
    }
    
    restore_from_backup "$1"
    log GREEN "Restore finalizado com sucesso"
}

main "$@"
