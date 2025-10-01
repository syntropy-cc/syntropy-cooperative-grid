#!/bin/bash

# Syntropy CLI Manager - Install Script
# Script para instalar a aplicação CLI

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
INSTALL_SCRIPT="$SCRIPT_DIR/scripts/linux/install-syntropy.sh"

# Functions
print_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
print_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
print_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# Banner
echo -e "\n${CYAN}╔══════════════════════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║              SYNTROPY CLI MANAGER                           ║${NC}"
echo -e "${CYAN}║                Install Script                               ║${NC}"
echo -e "${CYAN}╚══════════════════════════════════════════════════════════════╝${NC}\n"

# Check if install script exists
if [ ! -f "$INSTALL_SCRIPT" ]; then
    print_error "Install script not found: $INSTALL_SCRIPT"
    exit 1
fi

print_info "Starting Syntropy CLI Manager installation..."
print_info "Using install script: $INSTALL_SCRIPT"

# Execute the install script
exec "$INSTALL_SCRIPT"

