#!/bin/bash

# Syntropy CLI - Linting Script
# Verifica qualidade do código
#
# Author: Syntropy Team
# Version: 1.0.0

set -euo pipefail

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

log() {
    local level="$1"; shift
    case "$level" in
        "ERROR") echo -e "${RED}[ERROR]${NC} $*" >&2 ;;
        "INFO")  echo -e "${BLUE}[INFO]${NC}  $*" ;;
        "SUCCESS") echo -e "${GREEN}[SUCCESS]${NC} $*" ;;
    esac
}

run_golint() {
    log "INFO" "Executando golint..."
    
    cd ../../..
    
    if command -v golint >/dev/null; then
        golint ./... && log "SUCCESS" "✓ Golint passou"
    else
        log "INFO" "golint não instalado - pulando"
    fi
}

run_govet() {
    log "INFO" "Executando go vet..."
    
    go vet ./... && log "SUCCESS" "✓ Go vet passou"
}

main() {
    log "INFO" "Syntropy CLI - Linting"
    
    run_govet
    run_golint
    
    log "SUCCESS" "Linting concluído"
}

main "$@"
