#!/bin/bash

# Syntropy CLI - Script de Execução de Testes de Desenvolvimento
# Executa todos os testes para o setup component
#
# Author: Syntropy Team  
# Version: 1.0.0
# Date: 2025-01-27

set -euo pipefail

SCRIPT_NAME="run-tests.sh"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

log() {
    local level="$1"
    shift
    
    case "$level" in
        "ERROR") echo -e "${RED}[ERROR]${NC} $*" >&2 ;;
        "WARN")  echo -e "${YELLOW}[WARN]${NC}  $*" ;;
        "INFO")  echo -e "${BLUE}[INFO]${NC}  $*" ;;
        "SUCCESS") echo -e "${GREEN}[SUCCESS]${NC} $*" ;;
    esac
}

check_prerequisites() {
    log "INFO" "Verificando prerequisitos..."
    
    command -v go >/dev/null || {
        log "ERROR" "Go compiler não encontrado"
        return 1
    }
    
    log SUCCESS "✓ Go disponível"
}

run_unit_tests() {
    log "INFO" "Executando testes unitários..."
    
    cd ../../.. 2>/dev/null || cd .. 
    
    if go test -short -race -v ./tests/unit/... 2>&1 | grep -v "^#" ; then
        log SUCCESS "✓ Testes unitários passaram"
        return 0
    else
        log ERROR "✗ Alguns testes unitários falharam"
        return 1
    fi
}

run_integration_tests() {
    log "INFO" "Executando testes de integração..."
    
    if go test -run=Integration -v ./tests/integration/... ; then
        log SUCCESS "✓ Testes de integração passaram"
    else
        log ERROR "✗ Alguns testes de integração falharam"
        return 1
    fi
}

parse_arguments() {
    PARSE_UNIT=""
    PARSE_INTEGRATION=""
    
    while [[ $# -gt 0 ]]; do
        case "$1" in
            --unit)     PARSE_UNIT="true";     shift ;;
            --integration) PARSE_INTEGRATION="true"; shift ;;  
            --help)     
                cat <<EOF
Usage: $0 [options]
Options:
  --unit          Executa apenas testes unitários  
  --integration   Executa apenas testes integração
  --help          Exibir esta ajuda
EOF
                exit 0 ;;
        esac
    done
}

main() {
    log "INFO" "Syntropy CLI - Test Runners"
    
    parse_arguments "$@"
    check_prerequisites
    
    local exit_code=0
    
    if [[ "$PARSE_UNIT" == "true" ]] || [[ "$PARSE_INTEGRATION" == "" ]]; then
        run_unit_tests || exit_code=1
    fi
    
    if [[ "$PARSE_INTEGRATION" == "true" ]] || [[ "$PARSE_UNIT" == "" ]]; then
        run_integration_tests || exit_code=1
    fi
    
    [[ $exit_code -eq 0 ]] && {
        log SUCCESS "🎉 Todos os testes passaram!"
    } || {
        log ERROR "❌ Testes falharam"
    }
    
    exit $exit_code
}

main "$@"
