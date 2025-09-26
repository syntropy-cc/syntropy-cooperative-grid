#!/bin/bash

# Syntropy CLI - Instalação como Serviço systemd
# Script para configurar Syntropy CLI como serviço systemd
#
# Author: Syntropy Team
# Version: 1.0.0
# Created: 2025-01-27

set -euo pipefail

SCRIPT_NAME="install-systemd.sh"

# Service settings
SERVICE_NAME="syntropy"
SERVICE_DESCRIPTION="Syntropy CLI Network Service"
SYNTHROPY_HOME="$HOME/.syntropy"

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

check_systemd() {
    log "INFO" "Verificando systemd..."
    
    command -v systemctl >/dev/null || {
        log "ERROR" "systemctl não encontrado - systemd necessário"
        return 1
    }
    
    log "SUCCESS" "✓ systemd disponível"
}

create_user_service() {
    log "INFO" "Criando serviço systemd para usuário..."
    
    mkdir -p "$HOME/.config/systemd/user"
    
    cat > "$HOME/.config/systemd/user/${SERVICE_NAME}.service" << EOF
[Unit]
Description=$SERVICE_DESCRIPTION
After=network.target

[Service]
Type=simple  
ExecStart=$SYNTHROPY_HOME/bin/syntropy run-daemon
WorkingDirectory=$SYNTHROPY_HOME
Restart=always
RestartSec=5

[Install]
WantedBy=default.target
EOF

    log "SUCCESS" "✓ Arquivo de serviço criado"
}

enable_service() {
    log "INFO" "Habilitando serviço..."
    
    systemctl --user daemon-reload
    systemctl --user enable "${SERVICE_NAME}.service"
    
    log "SUCCESS" "✓ Serviço habilitado"
}

show_commands() {
    log "INFO" "Para gerenciar o serviço:"
    echo "  Start:    systemctl --user start $SERVICE_NAME"
    echo "  Stop:     systemctl --user stop $SERVICE_NAME"  
    echo "  Status:   systemctl --user status $SERVICE_NAME"
    echo "  Logs:     journalctl --user -u $SERVICE_NAME"
}

parse_arguments() {
    local activate="false"
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --activate)
                activate="true"
                shift
                ;;
            --help)
                echo "Usage: $0 [--activate]"
                exit 0
                ;;
        esac
    done
}

main() {
    log "INFO" "Syntropy CLI Systemd Service Installation"
    
    check_systemd || exit 1
    create_user_service
    enable_service
    
    show_commands
    
    # Optionally auto-start service
    if [[ "${1:-}" == "--activate" ]]; then
        systemctl --user start "$SERVICE_NAME"
        log "SUCCESS" "✓ Serviço iniciado automaticamente"
    fi
    
    log "SUCCESS" "Instalação concluída"
}

main "$@"
