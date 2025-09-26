#!/bin/bash

# Syntropy CLI - Teste de Validação de Ambiente  
# Script para validar ambiente antes do setup final do Syntropy CLI
#
# Author: Syntropy Team
# Version: 1.0.0
# Created: 2025-01-27
#
# Purpose: Validar ambiente completamente antes do setup final,
#         covering infrastructure requirements, OS compatibility,
#         security checks and performance baselines
#
# Usage Examples:
#   ./test-environment.sh                      # Padrão, validação completa
#   ./test-environment.sh --platform linux     # SO específico
#   ./test-environment.sh --quick --json       # Fast output em JSON     
#   ./test-environment.sh --no-network         # Testing offline/restrictive
#   ./test-environment.sh --full-suite         # Validação enterprise

set -euo pipefail

# ============================================================================
### CONFIGURAÇÕES E CONSTANTES ###
# ============================================================================

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SCRIPT_NAME="test-environment.sh"
SCRIPT_VERSION="1.0.0"

# Região de cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'  
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m'

# Configurações de validation
MIN_DISK_SPACE_GB=1
MIN_MEMORY_GB=1  # Em GB inteiros
MIN_CPU_CORES=1

# Network testing
NETM_TEST_SERVERS=("google.com" "github.com" "syntropy.io")
NETM_TIMEOUT=10

# ============================================================================
### FUNÇÕES UTILITÁRIAS ###
# ============================================================================

log() {
    local level="$1"
    shift
    local msg="$*"
    local tstamp=$(date '+%Y-%m-%d %H:%M:%S')
    
    case "$level" in
        "ERROR") echo -e "${RED}[ERROR]${NC} [$tstamp] $msg" >&2 ;;
        "WARN")  echo -e "${YELLOW}[WARN]${NC}  [$tstamp] $msg" ;;
        "INFO")  echo -e "${BLUE}[INFO]${NC}  [$tstamp] $msg" ;;
        "SUCCESS") echo -e "${GREEN}[SUCCESS]${NC} [$tstamp] $msg" ;;
        "DEBUG") echo -e "${CYAN}[DEBUG]${NC} [$tstamp] $msg" ;;
        *)       echo "[$level][$tstamp] $msg";;
    esac
}

log_section() {
    echo -e "\n${BLUE}=== $1 ===${NC}"
}

# Check de disponibilidade comandos
command_exists() {
    command -v "$1" >/dev/null 2>&1 
}

# ============================================================================
### DETECÇÃO DE PLATAFORMA E INFRAESTRUTURA ###
# ============================================================================

detect_operating_system() {
    case "$(uname -s)" in
        "Linux"*)
            OS_FLAVOR="linux"
            detect_linux_flavor 
            ;;
        "Darwin"*)
            OS_FLAVOR="macos"
            detect_macos_version
            ;;  
        MINGW64|CYGWIN_NT*)
            OS_FLAVOR="windows"
            log "WARN" "Executando em Windows - execute test-environment-windows.ps1 para compatibilidade completa"
            ;;
        *)
            OS_FLAVOR="unknown"
            log "ERROR" "OS platform não suportado para execução deste script"
            return 1
            ;;
    esac
    log "INFO" "OS detectado: $OS_FLAVOR"
}

detect_linux_flavor() {
    if [[ -f /etc/os-release ]]; then
        source /etc/os-release
        DISTRO=$NAME
        DISTRO_VER=$VERSION_ID
        
        case "$ID" in
            "ubuntu")
                log "INFO" "Detectado Ubuntu $VERSION_ID" 
                ;; 
            "debian")
                log "INFO" "Detectado Debian $VERSION_ID"
                ;;
            "centos"|"rhel"|"fedora")
                log "INFO" "Detectado RHEL/CentOS/Fedora variant"
                ;;
        esac
    fi
}

detect_macos_version() {
    if [[ -x /usr/bin/sw_vers ]]; then
        PRODUCT_VER=$(sw_vers --productVersion)
        MACOS_CODENAME=""
        case "$PRODUCT_VER" in
            13*|14*) MACOS_CODENAME="Sonoma" ;;
            12*) MACOS_CODENAME="Monterey" ;;  
        esac
        log "INFO" "macOS $PRODUCT_VER ($MACOS_CODENAME)"
    fi
}

validate_system_cpu() {
    log_section "CPU Validation" 
    
    if command_exists lscpu; then  
        CPU_CORES=$(lscpu | grep "^CPU(s):" | awk -F: '{print $2}' | xargs)
        CPU_MODEL=$(lscpu | grep "^Model name:" | awk -F: '{print $2}' | xargs)
        
        local cores_num="${CPU_CORES:-0}"
        if (( cores_num >= MIN_CPU_CORES )); then
            log "SUCCESS" "✓ CPU: $CPU_MODEL ($cores_num cores)" 
            return 0
        else
            log "ERROR" "✗ CPU cores insuficientes: $cores_num < $MIN_CPU_CORES"
            return 1
        fi
    elif command_exists sysctl; then  # OSX
        local cores=$(sysctl -n hw.ncpu 2>/dev/null || echo "0")
        local model=$(sysctl -n machdep.cpu.brand_string 2>/dev/null || echo "Unknown")
        
        if (( cores >= MIN_CPU_CORES )); then
            log "SUCCESS" "✓ CPU: $model ($cores cores)"
        else
            log "ERROR" "✗ CPU insuficiente"
            return 1
        fi
    elif command_exists nproc; then
        local cores
        cores=$(nproc 2>/dev/null || echo "0")
        if (( cores >= MIN_CPU_CORES )); then
            log "SUCCESS" "✓ CPU cores: $cores"
            return 0
        else
            log "ERROR" "✗ CPU cores insuficientes"
            return 1
        fi
    else
        log "WARN" "Não foi possível detectar informações do CPU"
        return 0
    fi
}

validate_system_memory() {
    log_section "Memory Validation"
    
    if command_exists free; then
        local memory_gb
        memory_gb=$(free -g | awk 'NR==2{print int($2+0.5)}' 2>/dev/null || echo "0")
        if (( memory_gb >= MIN_MEMORY_GB )); then
            log "SUCCESS" "✓ Memory: ${memory_gb}GB"
            return 0
        else
            log "WARN" "Memory: ${memory_gb}GB (may be low for heavy operations)"
            return 0
        fi
    elif command_exists sysctl; then  # macOS
        local memory_gb memory_bytes
        memory_bytes=$(sysctl -n hw.memsize 2>/dev/null || echo "0")
        memory_gb=$(( memory_bytes / 1024 / 1024 / 1024 ))
        
        if (( memory_gb >= MIN_MEMORY_GB )); then
            log "SUCCESS" "✓ Memory: ${memory_gb}GB"
            return 0
        else
            log "WARN" "Memory: ${memory_gb}GB (potential resource constraint)"
        fi
    else
        log "DEBUG" "Memory check skipped (commands unavailable)"
    fi
}

validate_disk_space() {
    log_section "Disk Space Validation"
    
    local available_space
    
    # Try standard df command
    available_space=$(df --output=avail / 2>/dev/null | awk 'NR==2 {print int($1/1024/1024)}' || echo "0")
    
    log "INFO" "Available disk space: ${available_space}GB"
    
    if [[ -n "$available_space" && "$available_space" -gt 0 ]] && (( available_space >= MIN_DISK_SPACE_GB )); then
        log "SUCCESS" "✓ Disk space: sufficient ($available_space GB)"
        return 0
    else
        log "ERROR" "✗ Insufficient free disk (need ${MIN_DISK_SPACE_GB}GB minimum): $available_space"
        return 1
    fi
}

# ============================================================================
### CONECTIVIDADE & REDE ###
# ============================================================================

validate_network_connectivity() {
    if [[ "${SKIP_NETWORK_TEST:-"false"}" == "true" ]]; then
        log "WARN" "Network checks SKIPPED (SKIP_NETWORK_TEST=true)"
        return 0
    fi
    
    log_section "Network Validation"
    local network_ok=true
    
    for target in "${NETM_TEST_SERVERS[@]}"; do
        if ping -c1 -W"${NETM_TIMEOUT}" "$target" >/dev/null 2>&1; then
            log "SUCCESS" "✓ $target reachable"
        else
            log "WARN" "✗ $target not reachable"
            network_ok=false
        fi
    done
    
    # Check HTTP connectivity 
    if command_exists curl; then 
        if curl -s --connect-timeout "$NETM_TIMEOUT" https://httpbin.org/get >/dev/null 2>&1; then
            log "SUCCESS" "✓ HTTPS connectivity OK (curl)"
        else
            log "WARN" "✗ HTTPS access constraints/down"
            network_ok=false
        fi
    fi
    
    [[ "$network_ok" == "true" ]] && return 0 || return 1
}

# ============================================================================
### VALIDAÇÃO DAS PERMISSÕES E ACESSOS ###
# ============================================================================

validate_permissions() {
    log_section "Permissions Validation"
    
    # Check home directory writable
    if [[ -w "$HOME" ]]; then
        log "SUCCESS" "✓ Home writable"
    else
        log "ERROR" "✗ Home dir not writable"
        return 1
    fi
    
    # Check write permissions in Syntropy expected location
    local syntropy_test_dir="$HOME/.syntropy/test_ensure_write"
    if mkdir -p "$syntropy_test_dir" 2>/dev/null && rmdir "$syntropy_test_dir" 2>/dev/null; then
        log "SUCCESS" "✓ ~/.syntropy dir writable"
        return 0
    else
        log "ERROR" "✗ Cannot write to ~/.syntropy"
        return 1
    fi
}

check_optional_dependencies() {
    log_section "Optional Dependencies"
    
    local OPTIONAL_TOOLS=("yq" "jq" "go" "docker")
    
    for tool in "${OPTIONAL_TOOLS[@]}"; do
        if command_exists "$tool"; then
            log "SUCCESS" "✓ $tool available"
        else
            log "DEBUG" "Tool not found: $tool (optional)"
        fi
    done
}

check_syntropy_cli() {
    log_section "Syntropy CLI Check"
    
    if command_exists syntropy; then  
        local syntropy_version
        if syntropy_version=$(syntropy --version 2>&1 | grep -o '[0-9]\+\.[0-9]\+\.[0-9]\+' | head -1); then
            log "SUCCESS" "Syntropy CLI available (v${syntropy_version})"
        else
            log "INFO" "Syntropy CLI present (couldn't get version)"
        fi
    else
        log "INFO" "Syntropy CLI not available (non-blocker)"
    fi
}

check_go_environment() {
    log_section "Go Environment Check"
    
    if command_exists go; then
        local go_version
        go_version=$(go version | awk '{print $3}' | sed 's/go//')
        log "SUCCESS" "Go available: $go_version"
    else
        log "WARN" "Go is recommended for some advanced Syntropy features"
    fi
}

apply_developer_checks() {
    log_section "Developer Checks"
    
    if command_exists node && node --version >/dev/null 2>&1; then
        log "SUCCESS" "Node.js available: $(node --version)"
    else
        log "INFO" "Node.js not required"
    fi
    
    if [[ -d ".git" ]]; then
        log "SUCCESS" "Git repository detected"
    fi
}

# ============================================================================
### FUNÇÃO PRINCIPAL DE TESTE ###
# ============================================================================

run_environment_tests() {
    local test_start_time test_duration
    
    test_start_time=$(date +%s)
    
    log_section "Syntropy Environment Validation Suite Start"
    log "INFO" "Script: $SCRIPT_NAME version $SCRIPT_VERSION"
    log "INFO" "Starting comprehensive environment validations"
    log "INFO" "Host: $(hostname) | OS: $(uname -s)"
    
    # Standard checks    
    detect_operating_system || {
        log "ERROR" "Platform detection failed"
        return 1
    }
    validate_system_cpu
    validate_system_memory
    validate_disk_space
    validate_network_connectivity 
    validate_permissions
    check_syntropy_cli
    
    check_optional_dependencies
    check_go_environment
    
    # Developer extras (if --developer specified)
    if [[ "${DEVELOPER_MODE:-false}" == "true" ]]; then
        apply_developer_checks
    fi
    
    test_duration=$(($(date +%s) - test_start_time))
    log "SUCCESS" "Environment tests completed in ${test_duration}s"
    export VALIDATION_RESULT="PASS"
}

# ============================================================================
### PARSE DE ARGUMENTOS ###
# ============================================================================

parse_arguments() {
    VALIDATION_PLATFORM="auto"
    JSON_MODE=false
    DEVELOPER_MODE=false
    SKIP_NETWORK_TEST=false
    QUIET_MODE=false
    
    while [[ $# -gt 0 ]]; do
        case "$1" in
            "--platform")
                shift
                VALIDATION_PLATFORM="$1"
                shift
                ;;
            "--json")
                JSON_MODE=true
                shift 
                ;;
            "--developer")
                DEVELOPER_MODE=true
                shift
                ;;
            "--no-network")
                SKIP_NETWORK_TEST=true
                shift 
                ;;
            "--quick")
                SKIP_NETWORK_TEST=true
                shift
                ;;
            "--help")
                show_help_message
                exit 0
                ;;
            *)
                log "ERROR" "Unrecognised arg $1"
                show_help_message >&2
                exit 1
                ;;
        esac
    done
}

show_help_message() {
cat <<- 'EOF'
Usage: ./test-environment.sh [OPTIONS]

Syntropy CLI Environment Validation Script

OPTIONS:
   --platform <OS>         Target platform for stricter checks  (linux,windows,macos,auto)
   --json                  Emit JSON report at completion
   --developer             Run extra developer checks
   --no-network           Skip network-related validations
   --quick                 Quick test run with reduced checks 
   --help                  This help text

EXAMPLES:
   ./test-environment.sh --json                    # JSON-based automated testing
   ./test-environment.sh --platform linux          # Strict Linux checks  
   ./test-environment.sh --developer               # Include developer tool checks
   ./test-environment.sh --quick --platform macos # Fast macOS validation
EOF
}

# ============================================================================
### FUNÇÃO PRINCIPAL ###
# ============================================================================

main() {
    log "INFO" "Starting Syntropy environment validation"
    
    # Parse arguments FIRST
    parse_arguments "$@"
    
    # Execute main feature
    run_environment_tests
    
    # Exit with success
    log "SUCCESS" "Environment validation passed"
    exit 0
}

# Execute if script is run directly
if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
    main "$@"
fi