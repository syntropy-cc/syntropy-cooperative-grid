#!/bin/bash

# Syntropy CLI - Validação Completa do Sistema  
# Script para validar sistema completo antes de setup do Syntropy CLI
#
# Author: Syntropy Team
# Version: 1.0.0
# Created: 2025-01-27
# Purpose: Executar validação abrangente do sistema
#
# Usage:
#   ./validate-system.sh
#   ./validate-system.sh --output report.json
#   ./validate-system.sh --exclude network

set -euo pipefail

SCRIPT_NAME="validate-system.sh"
SCRIPT_VERSION="1.0.0"

# Configuration defaults
OUTPUT_FILE=""
EXCLUDE_TESTS=()
CONTINUE_ON_ERRORS=false

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Global validation results
declare -A results=( )
total_tests=0
passed_tests=0

# Functions
log() {
    local level="$1"
    shift
    local msg="$*"
    local timestamp=$(date '+%Y-%m-%dT%H:%M:%S')
    
    case "$level" in
        "ERROR") echo -e "${RED}[ERROR]${NC} [$timestamp] $msg" >&2 ;;
        "WARN")  echo -e "${YELLOW}[WARN]${NC}  [$timestamp] $msg" ;;
        "INFO")  echo -e "${BLUE}[INFO]${NC}  [$timestamp] $msg" ;;
        "SUCCESS") echo -e "${GREEN}[SUCCESS]${NC} [$timestamp] $msg" ;;
    esac
}

# Check system requirements
validate_os_compatibility() {
    ((++total_tests))
    log "INFO" "Validando compatibilidade do sistema operacional..."
    
    local os=$(uname -s)
    local arch=$(uname -m)
    
    case "$os" in
        Linux)
            # Check distributions support
            [[ -f /etc/os-release ]] && source /etc/os-release
            results["os"]="linux-$ID_${VERSION_ID:-unknown}"
            ((++passed_tests))
            log "SUCCESS" "✓ Linux compatível detectado: ${PRETTY_NAME:-Unknown}"
            ;;
        Darwin|FreeBSD|OpenBSD)
            results["os"]="unix-$os"
            ((++passed_tests))
            log "SUCCESS" "✓ Unix-like OS compatível: $os"
            ;;
        *)
            log "ERROR" "✗ OS não suportado: $os"
            [[ "$CONTINUE_ON_ERRORS" == "false" ]] && exit 1 || true
            ;;
    esac
}

validate_system_resources() {
    ((++total_tests))
    log "INFO" "Validando recursos do sistema..."
    
    # Check available memory
    local memory_kb
    if command -v free >/dev/null 2>&1; then
        memory_kb=$(free -m | grep '^Mem:' | awk '{print $7}' || echo "0")
        [[ $memory_kb -gt 512 ]] || log "WARN" "⚠ Memória disponível baixa: ${memory_kb}MB"
    fi
    
    # Check disk space
    local disk_free
    disk_free=$(df / | awk 'NR==2 {print int($4/1024)}' || echo "0")
    [[ $disk_free -gt 1024 ]] || log "WARN" "⚠ Espaço em disco baixo: ${disk_free}MB"
    
    ((++passed_tests))
    log "SUCCESS" "✓ Recursos do sistema verificados"
    results["resources"]="memory:${memory_kb:-?}MB disk:${disk_free}MB"
}

validate_permissions() {
    ((++total_tests))
    log "INFO" "Validando permissões necessárias..."
    
    if [[ -w "$HOME" ]]; then
        ((++passed_tests))
        log "SUCCESS" "✓ Home directory writable"
        results["permissions"]="ok"
    else
        log "ERROR" "✗ Home directory não pode ser escrito"
        [[ "$CONTINUE_ON_ERRORS" == "false" ]] && exit 1
    fi
}

validate_network_connectivity() {
    # Skip já incluído na flag exclude
    [[ "network" == "${EXCLUDE_TESTS[*]}" && " ${EXCLUDE_TESTS[*]}" == *"network"* ]] && return
    
    ((++total_tests))
    log "INFO" "Validando conectividade de rede..."
    
    local test_servers=("google.com" "github.com")
    local connectivity_ok=true
    
    for server in "${test_servers[@]}"; do
        if ping -c1 -W5 "$server" >/dev/null 2>&1; then
            log "SUCCESS" "✓ Conectividade confirmada para $server"
        else
            log "WARN" "⚠ Sem resposta de $server"
            connectivity_ok=false
        fi
    done
    
    [[ "$connectivity_ok" == "true" ]] && {
        ((++passed_tests))
        results["network"]="ok"
        log "SUCCESS" "✓ Conectividade de rede funcional"
    } || log "WARN" "⚠ Conectividade limitada detectada"
}

validate_syntropy_cli() {
    ((++total_tests))
    log "INFO" "Verificando sintropy CLI..."
    
    if command -v syntropy >/dev/null 2>&1; then
        local version=$(syntropy --version 2>/dev/null | head -1 || echo "version_unavailable")
        ((++passed_tests))
        log "SUCCESS" "✓ Syntropy CLI encontrado: $version"
        results["cli"]="$version"
    else
        log "INFO" "ℹ Syntropy CLI não encontrado (setup irá instalar)"
        results["cli"]="not_installed"
    fi
}

validate_dependencies() {
    ((++total_tests))
    log "INFO" "Validando dependências opcionais..."
    
    local required=("curl" "tar")
    local optional=("git" "wget" "python3")
    local found_required=0
    
    for dep in "${required[@]}"; do
        if command -v "$dep" >/dev/null 2>&1; then
            log "SUCCESS" "✓ $dep disponível"
            ((++found_required))
        else
            log "ERROR" "✗ $dep não encontrado (obrigatório)"
            [[ "$CONTINUE_ON_ERRORS" == "false" ]] && exit 1
        fi
    done
    
    for dep in "${optional[@]}"; do
        if command -v "$dep" >/dev/null 2>&1; then
            log "INFO" "ℹ $dep disponível (opcional)"
        fi
    done
    
    ((++passed_tests))
    results["dependencies"]="${found_required}/${#required[@]}"
}

generate_report() {
    local duration=$((SECONDS))
    local success_rate
    [[ $total_tests -gt 0 ]] && success_rate=$((passed_tests * 100 / total_tests))
    
    log "INFO" "======= Validação completada em ${duration}s ======="
    log "INFO" "Testes executados: $total_tests"
    log "SUCCESS" "Testes passaram: $passed_tests"
    
    if [[ -n "$OUTPUT_FILE" ]]; then
        cat > "$OUTPUT_FILE" <<EOF
{
  "validation": {
    "script": "$SCRIPT_NAME",
    "version": "$SCRIPT_VERSION", 
    "timestamp": "$(date -ISO8601)",
    "summary": {
      "total_tests": $total_tests,
      "passed_tests": $passed_tests,
      "success_rate": "${success_rate:-0}%",
      "duration_sec": $duration
    },
    "results": {
$(for key in "${!results[@]}"; do
  echo "      \"$key\": \"${results[$key]}\","
done)
    }
  }
}
EOF
        log "INFO" "Relatório JSON salvo em: $OUTPUT_FILE"
    fi
    
    # Return result
    [[ $passed_tests -eq $total_tests ]] && return 0 || return 1
}

parse_arguments() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --output | -o)
                OUTPUT_FILE="$2"
                shift 2
                ;;
            --exclude)
                EXCLUDE_TESTS+=("$2")
                shift 2
                ;;
            --continue)
                CONTINUE_ON_ERRORS=true
                shift
                ;;
            --help)
                cat <<'EOF'
Usage: $0 [options]
Options:
  --output FILE    Salvar relatório JSON em FILE
  --exclude test   Excluir tipo de teste (ex: network)
  --continue       Continuar mesmo com alguns erros
  --help          Mostrar esta ajuda
EOF
                exit 0
                ;;
        esac
    done
}

# Main execution flow
run_validation_suite() {
    log "INFO" "Iniciando validação completa do sistema"
    echo "--------------------------------------------------"
    
    validate_os_compatibility
    validate_system_resources  
    validate_permissions
    validate_network_connectivity
    validate_syntropy_cli
    validate_dependencies
    
    generate_report
}

# Entry point
main() {
    log "INFO" "Validador do Sistema Syntropy CLI"
    log "INFO" "Script: $SCRIPT_NAME versão $SCRIPT_VERSION"
    
    parse_arguments "$@"
    run_validation_suite
}

# Execute main() when script is run directly
[[ "${BASH_SOURCE[0]}" == "$0" ]] && main "$@"
