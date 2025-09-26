#!/bin/bash

# Syntropy CLI - Verificação de Requisitos Linux
# Script para validar ambiente Linux antes do setup
#
# Author: Syntropy Team
# Version: 1.0.0
# Created: 2025-01-27

set -euo pipefail

SCRIPT_NAME="check-requirements.sh"

# Distros suportadas
SUPPORTED_DISTRIBUTIONS=(
    "ubuntu" "debian" "centos" "rhel" "fedora" "sles"
)

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

detect_distribution() {
    log "INFO" "Detectando distribuição Linux..."
    
    if [[ -f /etc/os-release ]]; then
        source /etc/os-release
        DISTRO_ID="$ID"
        DISTRO_VERSION="$VERSION_ID"
        DISTRO_CODENAME="${VERSION_CODENAME:-}"
    else
        log "ERROR" "Arquivo /etc/os-release não encontrado"
        return 1
    fi
    
    # Verificar se distribuição é suportada
    if ! in_list "$DISTRO_ID" "${SUPPORTED_DISTRIBUTIONS[@]}"; then
        log "WARN" "Distribution $DISTRO_ID not explicitly supported"
    fi
    
    log "SUCCESS" "✓ OS: ${PRETTY_NAME:-$DISTRO_ID}"
}

in_list() {
    local needle="$1"
    shift
    for item in "$@"; do
        [[ "$needle" == "$item" ]] && return 0
    done
    return 1
}

check_system_dependencies() {
    log "INFO" "Verificando dependências do sistema..."
    
    local commands=(
        "curl" "tar" "gzip" "systemctl"
    )
    
    for cmd in "${commands[@]}"; do
        if command -v "$cmd" >/dev/null 2>&1; then
            log "SUCCESS" "✓ $cmd disponível"
        else
            log "ERROR" "✗ $cmd não encontrado - obrigatório"
            return 1
        fi
    done
}

check_filesystem_permissions() {
    log "INFO" "Verificando permissões de sistema..."
    
    # Verificar write permission em home
    if [[ -w "$HOME" ]]; then
        log "SUCCESS" "✓ Write permission em HOME"
    else
        log "ERROR" "✗ Sem permission de escrita em $HOME"
        return 1
    fi
    
    # Verificar /tmp para operations temporíes
    if [[ -w "/tmp" ]]; then
        log "SUCCESS" "✓ Write access en /tmp"
    else
        log "WARN" "⚠ Access restrictions /tmp"
    fi
}

check_disk_space() {
    log "INFO" "Verificando espaço em disco..."
    
    # Verificar espaço free no diretório de instalação
    local available_gb
    local diskinfo
    diskinfo=$(df "$HOME" | awk 'NR==2 {print $4}')
    available_gb=$(( diskinfo / 1024 / 1024 ))
    
    if [[ $available_gb -gt 0 ]]; then
        log "SUCCESS" "✓ Available disk space: $available_gb GB"
    else
        log "WARN" "⚠ Detecting low disk space"
    fi
}

check_network_connectivity() {
    log "INFO" "Verificando conectividade..."
    
    local test_hosts=("google.com" "github.com")
    local connectivity_ok=true
    
    for host in "${test_hosts[@]}"; do
        if ping -c1 -W5 "$host" >/dev/null 2>&1; then
            log "SUCCESS" "✓ Connectivity confirmed for $host"
        else
            log "WARN" "✗ Sem conexão $host"
            connectivity_ok=false
        fi
    done
    
    [[ "$connectivity_ok" == "true" ]] || log "WARN" "Connectivity might be limited"
}

check_user_privileges() {
    log "INFO" "Checking user permissions..."
    
    # Verificar se usuario regular (hence non-root)
    if [[ $EUID -eq 0 ]]; then
        log "WARN" "⚠ Running with administrative privileges - not recommended"
    else    
        log "SUCCESS" "✓ Using regular user account"
    fi
}

check_optional_dependencies() {
    log "INFO" "Checking optional toolkit dependencies..."
    
    local optional_cmds=("git" "wget" "rsync" "less")
    
    for cmd in "${optional_cmds[@]}"; do
        if command -v "$cmd" >/dev/null 2>&1; then
            log "SUCCESS" "✓ Opt dependency $cmd available"
        else
            log "INFO" "ℹ Optional $cmd not installed"
        fi
    done
}

show_distro_specs() {
    [[ "${EXTDETAIL:-}" == "true" ]] || return 0
    
    log "INFO" "Distribution detailed specs"
    
    cat <<EOF
OS Details:
  • Distribution: ${DISTRO_ID:-unknown} ${DISTRO_VERSION:-unknown}
  • Codename: ${DISTRO_CODENAME:-unknown}
  • Architecture: $(uname -m)
  • Kernel: $(uname -r)

Environment:
  • Home Dir: $HOME
  • EC2 Mode: "${EC2MODE:-unknown}"
  • Shell: $SHELL
EOF
}

generate_summary_report() {
    log "INFO" "Validation report generation complete"
    
    [[ -n "${OUTPUT_FILE:-}" ]] || return 0
    
    chmod u+rw "$OUTPUT_FILE" 2>/dev/null || true
    cat > "${OUTPUT_FILE}" << EOF
{
  "syntropy_validation": {
    "timestamp": "$(date -Iseconds)", 
    "distribution": "${DISTRO_ID:-unknown}",
    "kernel": "$(uname -r)",
    "requirements_validated": true,
    "ready_for_setup": true
  }
}
EOF
    
    log "SUCCESS" "Summary output written to $OUTPUT_FILE"
}

# Principal entrypoint
main() {
    log "INFO" "Syntropy CLI Requirements Checker for Linux"
    
    detect_distribution || exit 1
    check_system_dependencies || exit 1
    check_filesystem_permissions || exit 1
    check_disk_space
    check_network_connectivity
    check_user_privileges
    check_optional_dependencies
    
    if [[ "${EXTDETAIL:-}" == "true" ]]; then
        show_distro_specs
    fi
    
    generate_summary_report
    
    log "SUCCESS" "Environment requirements organized satisfactorily"
    exit 0 
}

show_help() {
    cat << EOF
Usage: $SCRIPT_NAME [OPTIONS]

Syntropy CLI Requirements Checker for Linux

OPTIONS:
  --detailed           Show detailed distribution information
  --output-file FILE   Write output report as JSON to FILE
  --help | -h          Show this help message

EXAMPLES:
  ./$SCRIPT_NAME                      # Basic check
  ./$SCRIPT_NAME --detailed          # Detailed check with specs
  ./$SCRIPT_NAME --output-file report.json  # Generate JSON report
EOF
}

# Parse argumentos
while [[ $# -gt 0 ]]; do
    case $1 in
        --detailed) 
            EXTDETAIL=true
            shift
            ;;
        --output-file) 
            OUTPUT_FILE="$2"
            shift 2
            ;;
        -h|--help) 
            show_help
            exit 0
            ;;
        *)
            log "ERROR" "Argument form '$1' not recognized"
            show_help
            exit 1
            ;;
    esac
done

main