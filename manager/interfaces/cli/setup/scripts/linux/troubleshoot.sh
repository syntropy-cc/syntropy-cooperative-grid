#!/bin/bash

# Syntropy CLI - Resolução de Problemas Linux
# Script para diagnóstico e troubleshooting de ambiente
#
# Author: Syntropy Team
# Version: 1.0.0
# Created: 2025-01-27

set -euo pipefail

SCRIPT_NAME="troubleshoot.sh"

# Global counters
exit_code=0
checkmode="false"
fixmode="false"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

log() {
    local level="$1"
    shift
    local msg="$*"
    case "$level" in
        "ERROR") echo -e "${RED}[ERROR]${NC} $msg" >&2; ((exit_code++)) ;;
        "WARN")  echo -e "${YELLOW}[WARN]${NC}  $msg" ;;
        "INFO")  echo -e "${BLUE}[INFO]${NC}  $msg" ;;
        "SUCCESS") echo -e "${GREEN}[SUCCESS]${NC} $msg" ;;
    esac
}

# System environment check
check_systemd_legacy() {
    log "INFO" "Checking systemd compatibility..."
    
    if command -v systemctl >/dev/null 2>&1; then
        log "SUCCESS" "✓ systemctl is available"
    else
        log "ERROR" "✗ systemctl NOT found (required)"
        return 1
    fi
    return 0
}

check_privileges() {
    log "INFO" "Checking user capabilities and user account..."
    
    if [[ ${EUID} -ne 0 ]]; then
        log "SUCCESS" "✓ Running as normal user (acceptable)"
        log "INFO"  "ℹ Note: You may need administrator rights for services"
    else
        log "WARN" "⚠ Running as root (may not be necessary)"
    fi
}

check_filesystem_permissions() {
    log "INFO" "Verifying filesystem writabilities..."
    
    local home_accessible=$(test -w "$HOME" && echo "OK" || echo "FAIL")
    
    if [[ "$home_accessible" == "OK" ]]; then
        log "SUCCESS" "✓ $HOME is writeable"
    else
        log "ERROR" "✗ $HOME NOT writable"
        return 1
    fi
}

probe_running_syntropy() {
    log "INFO" "Searching for active Syntropy/Syntropy-related processes..."
    
    local syntropy_processes=$(pgrep -fl syntropy 2>/dev/null || echo "")
    local service_name="${SERVICE_NAME:-syntropy}"
    local systemctl_status=$(systemctl show "$service_name" 2>&1 | grep ActiveState= || cat < <(printf '%s' 'NotFound value'))
    
    if [[ -n "$syntropy_processes" ]]; then
        log "SUCCESS" "✓ Found existing processes named syntropy"
    else
        log "INFO" "ℹ No running syntropy processes found"
    fi
    
    if [[ "$systemctl_status" == *'NotFound'* ]]; then
        log "INFO" "ℹ Service unit no longer exists (normal if removed)"
    fi
}

collect_syntropy_logs() {
    log "INFO" "Aggregating retrievable Syntropy logs..."
    
    # User space logs:
    if [[ -d "$HOME/.syntropy/logs" ]]; then
        echo "Found $HOME/.syntropy/logs/"
        ls "$HOME/.syntropy/logs/" 2>/dev/null || true
    else
        log "INFO" "No directory $HOME/.syntropy/logs found"
    fi
    
    # Try journalctl for user unit
    journalctl --user -u syntropy -n10 --no-pager 2>/dev/null || log "INFO" "No user service logs found"
}

self_repair_actions() {
    local action="$1"
    
    case "${action}" in
        "fix_permissions")
            if [[ -d "$HOME/.syntropy" ]]; then
                chmod -R 0700 "$HOME/.syntropy" 2>/dev/null || true
                log "SUCCESS" "Set directory mode 0700 on ~/.syntropy"
            fi
            ;;
        *)
            log "INFO" "No repair action for: $action"
            ;;
    esac
}

print_troubleshooting_summary() {
    echo "================================================================================"
    log "INFO" "Troubleshooting summary"
    log "INFO" "• Script: $SCRIPT_NAME version 1.0"
    log "INFO" "• Hostname: $(hostname)"
    log "INFO" "• OS: $(uname -s) $(uname -m)"
    
    if [[ $exit_code -eq 0 ]]; then
        log "SUCCESS" "✓ All probes reported no issues."
    else
        log "WARN" "Found $exit_code indicator(s) requiring attention."
    fi
}

# Command dispatch routine
invoke_main() {
    local args="$@"
    
    while [[ $# -gt 0 ]]; do
        case "$1" in
            --check)
                checkmode="true"
                shift
                ;;
            --fix)
                fixmode="true"
                shift
                ;;
            *)
                log "WARN" "Unknown option: $1"
                shift
                ;;
        esac
    done
    
    if [[ "$checkmode" == "true" ]]; then
        check_systemd_legacy
        check_privileges
        check_filesystem_permissions
        probe_running_syntropy
    fi
    
    if [[ "$fixmode" == "true" ]]; then
        log "INFO" "Executing self-repairs (case fix-mode)"
        self_repair_actions "fix_permissions"
    fi
    
    collect_syntropy_logs
    print_troubleshooting_summary
}

# Script main entrance
main() {
    invoke_main "$@"
    exit $exit_code
}

# Execute if script is run directly
if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
    main "$@"
fi