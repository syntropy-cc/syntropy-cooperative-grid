#!/bin/bash

# Syntropy CLI - Format Script
# Formata código Go automaticamente
#
# Author: Syntropy Team
# Version: 1.0.0

set -euo pipefail

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

log() {
    local level="$1"; shift
    case "$level" in
        "INFO")  echo -e "${BLUE}[INFO]${NC}  $*" ;;
        "SUCCESS") echo -e "${GREEN}[SUCCESS]${NC} $*" ;;
    esac
}

format_code() {
    log "INFO" "Formatando código Go..."
    
    cd ../../..
    
    go fmt ./... && log "SUCCESS" "✓ Código formatado"
}

main() {
    log "INFO" "Syntropy CLI - Code Formatter"
    format_code
}

main "$@"
