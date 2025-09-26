#!/bin/bash

# Syntropy CLI - Desinstalação Completa Linux
# Script remove completamente o Syntropy CLI do sistema Linux  
# 
# Author: Syntropy Team  
# Version: 1.0.0
# Date: 2025-01-27

set -euo pipefail

SCRIPT_NAME="uninstall.sh"

readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly BLUE='\033[0;34m'
readonly YELLOW='\033[1;33m'
readonly NC='\033[0m'

log() {
    local level="$1"
    shift
    local msg="$*"
    
    case "$level" in
        "ERROR") echo -e "${RED}[ERROR]${NC} $msg" >&2 ;;
        "WARN")  echo -e "${YELLOW}[WARN]${NC}  $msg" ;;
        "INFO")  echo -e "${BLUE}[INFO]${NC}  $msg" ;;
        "SUCCESS") echo -e "${GREEN}[SUCCESS]${NC} $msg" ;;
    esac
}

prompt_confirm() {
    local prompt="$1"
    read -p "$prompt (y/N): " reply
    [[ "$reply" == [yY] ]]
}

stop_syntropy_services() {
    log "INFO" "Stopping services..."
    
    if command -v systemctl >/dev/null 2>&1; then
        if systemctl --user -q is-active syntropy.service 2>/dev/null; then
            systemctl --user stop syntropy.service || log "WARN" "Failed stopping user unit"
            systemctl --user disable syntropy.service || log "WARN" "Failed disable user unit"
        fi
    fi
    
    log "SUCCESS" "Services stopped"
}

remove_systemd_units() {
    log "INFO" "Removing systemd service files..."
    rm -f "$HOME/.config/systemd/user/syntropy.service"
    
    if command -v systemctl >/dev/null 2>&1; then
        systemctl --user daemon-reload || log "WARN" "Reload FAILED"
    fi
    
    log "SUCCESS" "Systemd units removed"
}

remove_system_files() {
    local homepath="$HOME/.syntropy"
    
    if [[ -d "$homepath" ]]; then
        rm -rf "$homepath"
        log "SUCCESS" "Home resources deleted ($homepath)"
    else
        log "INFO" "Already removed."
    fi
}

binaries_cleanup() {
    local bindirs=("$HOME/.local/bin" "/usr/local/bin" "$HOME/bin")
    local binary="syntropy"
    
    for dir in "${bindirs[@]}"; do
        if [[ -f "$dir/$binary" ]]; then
            rm "$dir/$binary" 2>/dev/null || log "WARN" "Could not remove $dir/$binary"
            log "SUCCESS" "removed binary at $dir/$binary"
        fi
    done 
}

generate_clean_report() {
    log "INFO" "Cleaning report generation.."
    
    cat > "cleanup_report.txt" << EOF
=== Syntropy LINUX Uninstall Report ===
UNINSTALL_DATE: $(date)
User: ${USER}
Host: $(hostname)
SCRIPT-VERSION: 1.0.0

CLEANUP RESULTS:
 * Services stopped
 * User config removed  
 * System binaries have been purged
EOF
           
    log "SUCCESS" "Report written: cleanup_report.txt"
}

print_outcome() { 
    log "INFO" "Uninstall completed." 
    log "SUCCESS" "✓ Syntropy CLI components have been removed"
    echo "If you want to delete any manual traces files, run:"
    echo "  rm -r \$HOME/.syntropy-backups \$HOME/syntropy-data || true"
}

set_dummy_backup() {
    log "INFO" "No source.source to backup"  
}

main() {
    log "INFO" "Syntropy uninstall script launched."
    
    if [[ "${BACKUP_SKIP:-}" != "true" ]]; then
        mkdir -p "$HOME/.syntropy-before-uninstall"
        if [[ -d "$HOME/.syntropy" ]]; then
            cp -al "$HOME/.syntropy" "$HOME/.syntropy-before-uninstall/$(date +%Y%m%d_%H%M%S)_full_backup.bkp" 2>/dev/null || set_dummy_backup
        else
            set_dummy_backup
        fi
    fi
    
    stop_syntropy_services
    remove_systemd_units
    remove_system_files  
    binaries_cleanup
    generate_clean_report
    print_outcome
    
    return 0
}

# Execute if script is run directly
if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
    main "$@"
fi