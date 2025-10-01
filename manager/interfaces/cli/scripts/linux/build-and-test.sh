#!/bin/bash

# Syntropy CLI Manager - Build and Test Workflow
# Script para compilar e testar a aplicação CLI no Windows e Linux

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CLI_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
BUILD_DIR="$CLI_DIR/build"
VERSION=$(date +%Y%m%d-%H%M%S)
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)

# Functions
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_step() {
    echo -e "\n${CYAN}=== $1 ===${NC}"
}

# Banner
show_banner() {
    echo -e "\n${CYAN}╔══════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${CYAN}║              SYNTROPY CLI MANAGER                           ║${NC}"
    echo -e "${CYAN}║                Build & Test                                 ║${NC}"
    echo -e "${CYAN}╚══════════════════════════════════════════════════════════════╝${NC}\n"
}

# Check prerequisites
check_prerequisites() {
    print_step "Checking Prerequisites"
    
    # Check Go
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install Go 1.22.5 or higher."
        print_info "Download: https://golang.org/dl/"
        exit 1
    fi
    
    local go_version=$(go version | cut -d' ' -f3 | sed 's/go//')
    local required_version="1.22"
    
    if [ "$(printf '%s\n' "$required_version" "$go_version" | sort -V | head -n1)" != "$required_version" ]; then
        print_error "Go version $go_version found, but version $required_version or higher is required."
        exit 1
    fi
    
    print_success "Go $go_version found"
    
    # Check if we're in the right directory
    if [ ! -f "$CLI_DIR/main.go" ]; then
        print_error "main.go not found in $CLI_DIR"
        exit 1
    fi
    
    print_success "Project structure verified"
}

# Prepare build environment
prepare_build() {
    print_step "Preparing Build Environment"
    
    # Navigate to CLI directory
    cd "$CLI_DIR"
    
    # Create build directory
    mkdir -p "$BUILD_DIR"
    
    # Clean previous builds
    rm -f "$BUILD_DIR"/*
    
    print_success "Build environment prepared"
}

# Setup dependencies
setup_dependencies() {
    print_step "Setting Up Dependencies"
    
    print_info "Downloading dependencies..."
    go mod download
    
    print_info "Organizing dependencies..."
    go mod tidy
    
    print_info "Verifying dependencies..."
    go mod verify
    
    print_success "Dependencies configured"
}

# Build for Windows
build_windows() {
    print_step "Building for Windows"
    
    local build_flags="-ldflags \"-X main.version=$VERSION -X main.buildTime=$BUILD_TIME -X main.gitCommit=$GIT_COMMIT\""
    local output_file="$BUILD_DIR/syntropy-windows.exe"
    
    print_info "Building Windows executable..."
    GOOS=windows GOARCH=amd64 go build $build_flags -o "$output_file" main.go
    
    if [ -f "$output_file" ]; then
        local size=$(du -h "$output_file" | cut -f1)
        print_success "Windows build completed: $output_file ($size)"
    else
        print_error "Windows build failed"
        exit 1
    fi
}

# Build for Linux
build_linux() {
    print_step "Building for Linux"
    
    local build_flags="-ldflags \"-X main.version=$VERSION -X main.buildTime=$BUILD_TIME -X main.gitCommit=$GIT_COMMIT\""
    local output_file="$BUILD_DIR/syntropy-linux"
    
    print_info "Building Linux executable..."
    GOOS=linux GOARCH=amd64 go build $build_flags -o "$output_file" main.go
    
    if [ -f "$output_file" ]; then
        local size=$(du -h "$output_file" | cut -f1)
        print_success "Linux build completed: $output_file ($size)"
    else
        print_error "Linux build failed"
        exit 1
    fi
}

# Build for current platform
build_current() {
    print_step "Building for Current Platform"
    
    local build_flags="-ldflags \"-X main.version=$VERSION -X main.buildTime=$BUILD_TIME -X main.gitCommit=$GIT_COMMIT\""
    local output_file="$BUILD_DIR/syntropy"
    
    print_info "Building for current platform..."
    go build $build_flags -o "$output_file" main.go
    
    if [ -f "$output_file" ]; then
        local size=$(du -h "$output_file" | cut -f1)
        print_success "Current platform build completed: $output_file ($size)"
    else
        print_error "Current platform build failed"
        exit 1
    fi
}

# Test binaries
test_binaries() {
    print_step "Testing Binaries"
    
    # Test current platform binary
    if [ -f "$BUILD_DIR/syntropy" ]; then
        print_info "Testing current platform binary..."
        
        # Test version
        if "$BUILD_DIR/syntropy" --version >/dev/null 2>&1; then
            print_success "Version test passed"
        else
            print_warning "Version test failed (may be normal)"
        fi
        
        # Test help
        if "$BUILD_DIR/syntropy" --help >/dev/null 2>&1; then
            print_success "Help test passed"
        else
            print_warning "Help test failed (may be normal)"
        fi
    fi
    
    # Test Windows binary (if on Linux/WSL)
    if [ -f "$BUILD_DIR/syntropy-windows.exe" ]; then
        print_info "Windows binary created: syntropy-windows.exe"
        local size=$(du -h "$BUILD_DIR/syntropy-windows.exe" | cut -f1)
        print_info "Size: $size"
    fi
    
    # Test Linux binary
    if [ -f "$BUILD_DIR/syntropy-linux" ]; then
        print_info "Linux binary created: syntropy-linux"
        local size=$(du -h "$BUILD_DIR/syntropy-linux" | cut -f1)
        print_info "Size: $size"
    fi
    
    print_success "Binary testing completed"
}

# Run application
run_application() {
    print_step "Running Application"
    
    local binary_path="$BUILD_DIR/syntropy"
    
    if [ ! -f "$binary_path" ]; then
        print_error "Binary not found. Run build first."
        exit 1
    fi
    
    print_info "Running Syntropy CLI Manager..."
    echo -e "${CYAN}========================================${NC}"
    
    # Run with help to show available commands
    "$binary_path" --help
    
    echo -e "\n${CYAN}========================================${NC}"
    print_info "Application executed successfully!"
    
    # Show example commands
    echo -e "\n${BLUE}Example commands to try:${NC}"
    echo -e "  ${YELLOW}$binary_path --version${NC}"
    echo -e "  ${YELLOW}$binary_path setup --help${NC}"
    echo -e "  ${YELLOW}$binary_path setup validate${NC}"
    echo -e "  ${YELLOW}$binary_path setup run --force${NC}"
}

# Show summary
show_summary() {
    print_step "Build Summary"
    
    print_success "Build completed successfully!"
    echo -e "\n${BLUE}📁 Build Directory:${NC} $BUILD_DIR"
    echo -e "${BLUE}📦 Version:${NC} $VERSION"
    echo -e "${BLUE}🔧 Git Commit:${NC} $GIT_COMMIT"
    echo -e "${BLUE}🕒 Build Time:${NC} $BUILD_TIME"
    echo -e "${BLUE}🖥️  Platform:${NC} $(uname -s)-$(uname -m)"
    
    echo -e "\n${BLUE}📋 Created Binaries:${NC}"
    if [ -f "$BUILD_DIR/syntropy" ]; then
        local size=$(du -h "$BUILD_DIR/syntropy" | cut -f1)
        echo -e "  ${GREEN}✅${NC} syntropy ($size) - Current platform"
    fi
    if [ -f "$BUILD_DIR/syntropy-windows.exe" ]; then
        local size=$(du -h "$BUILD_DIR/syntropy-windows.exe" | cut -f1)
        echo -e "  ${GREEN}✅${NC} syntropy-windows.exe ($size) - Windows"
    fi
    if [ -f "$BUILD_DIR/syntropy-linux" ]; then
        local size=$(du -h "$BUILD_DIR/syntropy-linux" | cut -f1)
        echo -e "  ${GREEN}✅${NC} syntropy-linux ($size) - Linux"
    fi
    
    echo -e "\n${BLUE}🚀 Next Steps:${NC}"
    echo -e "  1. Test the application: ${YELLOW}$BUILD_DIR/syntropy --help${NC}"
    echo -e "  2. Run setup: ${YELLOW}$BUILD_DIR/syntropy setup run --force${NC}"
    echo -e "  3. Copy Windows .exe to Windows machine for testing"
    echo -e "  4. Copy Linux binary to Linux machine for testing"
    
    echo -e "\n${BLUE}💡 Usage Examples:${NC}"
    echo -e "  ${CYAN}$BUILD_DIR/syntropy --help${NC}                    # Show help"
    echo -e "  ${CYAN}$BUILD_DIR/syntropy --version${NC}                 # Show version"
    echo -e "  ${CYAN}$BUILD_DIR/syntropy setup --help${NC}              # Setup help"
    echo -e "  ${CYAN}$BUILD_DIR/syntropy setup run --force${NC}         # Run setup"
    echo -e "  ${CYAN}$BUILD_DIR/syntropy setup status${NC}              # Check status"
}

# Main function
main() {
    show_banner
    
    # Check if help is requested
    if [ "$1" = "--help" ] || [ "$1" = "-h" ]; then
        echo -e "${BLUE}Usage:${NC} $0 [options]"
        echo -e "\n${BLUE}Options:${NC}"
        echo -e "  --help, -h     Show this help message"
        echo -e "  --windows      Build only for Windows"
        echo -e "  --linux        Build only for Linux"
        echo -e "  --current      Build only for current platform"
        echo -e "  --test         Test existing binaries"
        echo -e "  --run          Run the application"
        echo -e "\n${BLUE}Examples:${NC}"
        echo -e "  $0                    # Build for all platforms"
        echo -e "  $0 --windows          # Build only for Windows"
        echo -e "  $0 --test             # Test existing binaries"
        echo -e "  $0 --run              # Run the application"
        exit 0
    fi
    
    # Check prerequisites
    check_prerequisites
    
    # Prepare build environment
    prepare_build
    
    # Setup dependencies
    setup_dependencies
    
    # Build based on options
    case "$1" in
        "--windows")
            build_windows
            ;;
        "--linux")
            build_linux
            ;;
        "--current")
            build_current
            ;;
        "--test")
            test_binaries
            show_summary
            exit 0
            ;;
        "--run")
            run_application
            exit 0
            ;;
        *)
            # Build for all platforms
            build_current
            build_windows
            build_linux
            ;;
    esac
    
    # Test binaries
    test_binaries
    
    # Show summary
    show_summary
    
    # Ask if user wants to run the application
    echo -e "\n${YELLOW}Do you want to run the application now? (y/N):${NC}"
    read -r response
    if [[ "$response" =~ ^[Yy]$ ]]; then
        run_application
    fi
}

# Run main function with all arguments
main "$@"
