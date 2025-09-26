#!/bin/bash

# Syntropy CLI - Build Script para Desenvolvimento
# Script para compilar o setup component
#
# Author: Syntropy Team
# Version: 1.0.0

set -euo pipefail

SCRIPT_NAME="build.sh"

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

build_setup() {
    log "INFO" "Building Syntropy CLI setup component..."
    
    # Navigate to project root
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../../.." && pwd)"
    cd "$PROJECT_ROOT"
    
    # Check if Go is available
    if ! command -v go >/dev/null 2>&1; then
        log "ERROR" "Go não está disponível"
        return 1
    fi
    
    # Build for multiple platforms
    local platforms=("linux/amd64" "linux/arm64" "windows/amd64" "darwin/amd64" "darwin/arm64")
    local build_dir="build"
    
    mkdir -p "$build_dir"
    
    for platform in "${platforms[@]}"; do
        IFS='/' read -r os arch <<< "$platform"
        output_name="syntropy-setup"
        
        if [ "$os" = "windows" ]; then
            output_name="${output_name}.exe"
        fi
        
        output_path="$build_dir/syntropy-setup-${os}-${arch}"
        if [ "$os" = "windows" ]; then
            output_path="$build_dir/syntropy-setup-${os}-${arch}.exe"
        fi
        
        log "INFO" "Building for $platform..."
        
        if GOOS="$os" GOARCH="$arch" go build -o "$output_path" ./...; then
            log "SUCCESS" "✓ Build concluído: $output_path"
        else
            log "ERROR" "✗ Build falhou para $platform"
        fi
    done
}

main() {
    log "INFO" "Syntropy CLI - Build Script"
    build_setup
}

main "$@"
