#!/bin/bash

# Syntropy CLI - Backup de Configuração
# Script para backup automatizado das configurações do setup
#
# Author: Syntropy Team
# Version: 1.0.0  
# Created: 2025-01-27
# Purpose: Backup automático de configurações e dados do Syntropy

set -euo pipefail

SCRIPT_NAME="backup-config.sh"
SCRIPT_VERSION="1.0.0"

# Configuration
BACKUP_DIR="${SYNTHROPY_BACKUP_DIR:-$HOME/.syntropy-backups}"
BACKUP_TIMESTAMP=$(date '+%Y%m%d_%H%M%S')
SOURCE_PATH="$HOME/.syntropy"

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

validate_source() {
    [[ -d "$SOURCE_PATH" ]] || {
        log YELLOW "Diretório de configuração não encontrado: $SOURCE_PATH"
        return 1
    }
    log GREEN "✓ Fonte encontrada: $SOURCE_PATH"
}

create_backup() {
    local backup_file="$BACKUP_DIR/syntropy-backup-$BACKUP_TIMESTAMP.tar.gz"
    
    log BLUE "Criando backup..."
    
    mkdir -p "$BACKUP_DIR"
    tar -czf "$backup_file" -C "$(dirname "$SOURCE_PATH")" "$(basename "$SOURCE_PATH")"
    
    log GREEN "✓ Backup criado: $backup_file"
}

main() {
    log BLUE "Syntropy CLI - Backup de Configuração"
    
    validate_source || {
        log RED "Erro: Configuração não encontrada"
        exit 1
    }
    
    create_backup
    log GREEN "Backup finalizado com sucesso"
}

main "$@"
