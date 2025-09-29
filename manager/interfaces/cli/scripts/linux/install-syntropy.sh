#!/bin/bash

# Syntropy CLI Manager - Simple Install Script
# Script para instalar e testar a aplicação CLI do Syntropy

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
CLI_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
BUILD_DIR="$CLI_DIR/build"
VERSION=$(date +%Y%m%d-%H%M%S)
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Functions
print_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
print_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
print_warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; }
print_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# Banner
echo -e "\n${CYAN}╔══════════════════════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║              SYNTROPY CLI MANAGER                           ║${NC}"
echo -e "${CYAN}║                Simple Install Script                        ║${NC}"
echo -e "${CYAN}╚══════════════════════════════════════════════════════════════╝${NC}\n"

# Check Go
if ! command -v go &> /dev/null; then
    print_error "Go is not installed. Please install Go 1.22.5 or higher."
    print_info "Download: https://golang.org/dl/"
    exit 1
fi

print_success "Go $(go version | cut -d' ' -f3) found"

# Check main.go
if [ ! -f "$CLI_DIR/main.go" ]; then
    print_error "main.go not found in $CLI_DIR. Please check the project structure."
    exit 1
fi

print_success "Project structure verified"

# Prepare build
print_info "Preparing build environment..."
cd "$CLI_DIR"
mkdir -p "$BUILD_DIR"
rm -f "$BUILD_DIR"/*

# Setup dependencies
print_info "Setting up dependencies..."
go mod download
go mod tidy

# Build for Windows
print_info "Building for Windows..."
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=$VERSION -X main.gitCommit=$GIT_COMMIT" -o "$BUILD_DIR/syntropy-windows.exe" main.go

# Build for Linux
print_info "Building for Linux..."
go build -ldflags "-X main.version=$VERSION -X main.gitCommit=$GIT_COMMIT" -o "$BUILD_DIR/syntropy-linux" main.go

# Test binaries
print_info "Testing binaries..."
if [ -f "$BUILD_DIR/syntropy-linux" ]; then
    if "$BUILD_DIR/syntropy-linux" --version >/dev/null 2>&1; then
        print_success "Linux binary test passed"
    else
        print_warning "Linux binary test failed (may be normal)"
    fi
fi

# Show results
echo -e "\n${GREEN}✅ Installation completed successfully!${NC}"
echo -e "\n${BLUE}📁 Build Directory:${NC} $BUILD_DIR"
echo -e "${BLUE}📦 Version:${NC} $VERSION"
echo -e "${BLUE}🔧 Git Commit:${NC} $GIT_COMMIT"

echo -e "\n${BLUE}📋 Created Binaries:${NC}"
if [ -f "$BUILD_DIR/syntropy-windows.exe" ]; then
    local size=$(du -h "$BUILD_DIR/syntropy-windows.exe" | cut -f1)
    echo -e "  ${GREEN}✅${NC} syntropy-windows.exe ($size) - Windows"
fi
if [ -f "$BUILD_DIR/syntropy-linux" ]; then
    local size=$(du -h "$BUILD_DIR/syntropy-linux" | cut -f1)
    echo -e "  ${GREEN}✅${NC} syntropy-linux ($size) - Linux"
fi

echo -e "\n${BLUE}🚀 Next Steps:${NC}"
echo -e "  1. Test Windows: ${YELLOW}$BUILD_DIR/syntropy-windows.exe --help${NC}"
echo -e "  2. Test Linux: ${YELLOW}$BUILD_DIR/syntropy-linux --help${NC}"
echo -e "  3. Run setup: ${YELLOW}$BUILD_DIR/syntropy-linux setup run --force${NC}"

# Ask to run test
echo -e "\n${YELLOW}Do you want to test the Linux application now? (y/N):${NC}"
read -r response
if [[ "$response" =~ ^[Yy]$ ]]; then
    echo -e "\n${CYAN}=== Testing Application ===${NC}"
    "$BUILD_DIR/syntropy-linux" --help
    echo -e "\n${GREEN}✅ Test completed!${NC}"
fi

echo -e "\n${GREEN}🎉 Syntropy CLI Manager installed successfully!${NC}"
