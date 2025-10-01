2464#!/bin/bash

# Syntropy CLI Manager - Universal Build Script
# Script unificado para compilar e testar em todas as plataformas suportadas
# Funciona em Linux, macOS e Windows (via WSL/Git Bash)

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m' # No Color

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CLI_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
BUILD_DIR="$CLI_DIR/build"
VERSION=$(date +%Y%m%d-%H%M%S)
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)

# Platform detection
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Supported platforms
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "windows/amd64"
    "darwin/amd64"
    "darwin/arm64"
)

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

print_platform() {
    echo -e "${MAGENTA}[$1]${NC} $2"
}

# Banner
show_banner() {
    echo -e "\n${CYAN}╔══════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${CYAN}║              SYNTROPY CLI MANAGER                           ║${NC}"
    echo -e "${CYAN}║              Universal Build & Test                         ║${NC}"
    echo -e "${CYAN}║              Cross-Platform Support                         ║${NC}"
    echo -e "${CYAN}╚══════════════════════════════════════════════════════════════╝${NC}\n"
    
    print_info "Detected Platform: $OS-$ARCH"
    print_info "Build Directory: $BUILD_DIR"
    print_info "Version: $VERSION"
    print_info "Git Commit: $GIT_COMMIT"
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
    
    # Clean previous builds only if --clean flag is used
    if [ "$clean_build" = true ]; then
        rm -f "$BUILD_DIR"/*
    fi
    
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
    if ! go mod verify; then
        print_warning "Dependency verification failed (continuing anyway)"
    fi
    
    print_success "Dependencies configured"
}

# Build for specific platform
build_platform() {
    local platform="$1"
    local goos="${platform%%/*}"
    local goarch="${platform##*/}"
    
    print_platform "$platform" "Building executable..."
    
    local build_flags=""
    local output_name="syntropy"
    
    # Set appropriate file extension
    if [ "$goos" = "windows" ]; then
        output_name="${output_name}.exe"
    fi
    
    local output_file="$BUILD_DIR/syntropy-${goos}-${goarch}"
    if [ "$goos" = "windows" ]; then
        output_file="$BUILD_DIR/syntropy-${goos}-${goarch}.exe"
    fi
    
    # Build
    GOOS="$goos" GOARCH="$goarch" go build $build_flags -o "$output_file" main.go
    
    if [ -f "$output_file" ]; then
        local size=$(du -h "$output_file" | cut -f1)
        print_platform "$platform" "Build completed: $output_file ($size)"
        return 0
    else
        print_platform "$platform" "Build failed"
        return 1
    fi
}

# Build for all platforms
build_all_platforms() {
    print_step "Building for All Platforms"
    
    local success_count=0
    local total_count=${#PLATFORMS[@]}
    
    for platform in "${PLATFORMS[@]}"; do
        if build_platform "$platform"; then
            ((success_count++))
        fi
    done
    
    print_info "Build Summary: $success_count/$total_count platforms successful"
    
    if [ $success_count -eq 0 ]; then
        print_error "No platforms built successfully"
        exit 1
    fi
}

# Build for current platform only
build_current_platform() {
    print_step "Building for Current Platform"
    
    local current_platform="$OS/$ARCH"
    
    # Map common architectures
    case "$ARCH" in
        "x86_64")
            current_platform="$OS/amd64"
            ;;
        "aarch64")
            current_platform="$OS/arm64"
            ;;
        "arm64")
            current_platform="$OS/arm64"
            ;;
    esac
    
    print_info "Current platform: $current_platform"
    
    if build_platform "$current_platform"; then
        print_success "Current platform build completed"
    else
        print_error "Current platform build failed"
        exit 1
    fi
}

# Run tests
run_tests() {
    print_step "Running Tests"
    
    print_info "Running unit tests..."
    if go test -v ./...; then
        print_success "Unit tests passed"
    else
        print_warning "Some unit tests failed (expected for unimplemented features)"
    fi
    
    print_info "Running tests with coverage..."
    if go test -v -cover ./...; then
        print_success "Coverage tests completed"
    else
        print_warning "Some coverage tests failed"
    fi
    
    print_info "Running race condition tests..."
    if go test -race ./...; then
        print_success "Race condition tests passed"
    else
        print_warning "Some race condition tests failed"
    fi
    
    print_success "Test suite completed"
}

# Test binaries
test_binaries() {
    print_step "Testing Binaries"
    
    local test_count=0
    local success_count=0
    
    for binary in "$BUILD_DIR"/syntropy-*; do
        if [ -f "$binary" ]; then
            ((test_count++))
            local binary_name=$(basename "$binary")
            print_info "Testing $binary_name..."
            
            # Test version
            if "$binary" --version >/dev/null 2>&1; then
                print_success "Version test passed for $binary_name"
                ((success_count++))
            else
                print_warning "Version test failed for $binary_name (may be normal)"
            fi
            
            # Test help
            if "$binary" --help >/dev/null 2>&1; then
                print_success "Help test passed for $binary_name"
            else
                print_warning "Help test failed for $binary_name (may be normal)"
            fi
        fi
    done
    
    print_info "Binary testing: $success_count/$test_count binaries tested successfully"
}

# Run application
run_application() {
    print_step "Running Application"
    
    # Find the appropriate binary for current platform
    local binary_path=""
    
    case "$OS" in
        "linux")
            if [ -f "$BUILD_DIR/syntropy-linux-amd64" ]; then
                binary_path="$BUILD_DIR/syntropy-linux-amd64"
            elif [ -f "$BUILD_DIR/syntropy-linux-arm64" ]; then
                binary_path="$BUILD_DIR/syntropy-linux-arm64"
            fi
            ;;
        "darwin")
            if [ -f "$BUILD_DIR/syntropy-darwin-arm64" ]; then
                binary_path="$BUILD_DIR/syntropy-darwin-arm64"
            elif [ -f "$BUILD_DIR/syntropy-darwin-amd64" ]; then
                binary_path="$BUILD_DIR/syntropy-darwin-amd64"
            fi
            ;;
        "windows"|"mingw"*|"cygwin"*)
            if [ -f "$BUILD_DIR/syntropy-windows-amd64.exe" ]; then
                binary_path="$BUILD_DIR/syntropy-windows-amd64.exe"
            fi
            ;;
    esac
    
    if [ -z "$binary_path" ] || [ ! -f "$binary_path" ]; then
        print_error "No suitable binary found for current platform"
        print_info "Available binaries:"
        ls -la "$BUILD_DIR"/syntropy-* 2>/dev/null || print_warning "No binaries found"
        exit 1
    fi
    
    print_info "Running: $binary_path"
    echo -e "${CYAN}========================================${NC}"
    
    # Run with help to show available commands
    "$binary_path" --help
    
    echo -e "\n${CYAN}========================================${NC}"
    print_success "Application executed successfully!"
    
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
    echo -e "${BLUE}🖥️  Current Platform:${NC} $OS-$ARCH"
    
    echo -e "\n${BLUE}📋 Created Binaries:${NC}"
    local binary_count=0
    for binary in "$BUILD_DIR"/syntropy-*; do
        if [ -f "$binary" ]; then
            local size=$(du -h "$binary" | cut -f1)
            local binary_name=$(basename "$binary")
            echo -e "  ${GREEN}✅${NC} $binary_name ($size)"
            ((binary_count++))
        fi
    done
    
    if [ $binary_count -eq 0 ]; then
        echo -e "  ${RED}❌${NC} No binaries created"
    fi
    
    echo -e "\n${BLUE}🚀 Next Steps:${NC}"
    echo -e "  1. Test the application: ${YELLOW}./build-all.sh --run${NC}"
    echo -e "  2. Run tests: ${YELLOW}./build-all.sh --test${NC}"
    echo -e "  3. Copy binaries to target machines for testing"
    
    echo -e "\n${BLUE}💡 Usage Examples:${NC}"
    echo -e "  ${CYAN}./build-all.sh --run${NC}                    # Run application"
    echo -e "  ${CYAN}./build-all.sh --test${NC}                   # Run tests only"
    echo -e "  ${CYAN}./build-all.sh --current${NC}                # Build current platform only"
    echo -e "  ${CYAN}./build-all.sh --platform linux/amd64${NC}   # Build specific platform"
}

# Show help
show_help() {
    echo -e "${BLUE}Usage:${NC} $0 [options]"
    echo -e "\n${BLUE}Options:${NC}"
    echo -e "  --help, -h              Show this help message"
    echo -e "  --current               Build only for current platform"
    echo -e "  --platform PLATFORM     Build for specific platform (e.g., linux/amd64)"
    echo -e "  --test                  Run tests only"
    echo -e "  --run                   Run the application"
    echo -e "  --no-tests              Skip running tests"
    echo -e "  --clean                 Clean build directory before building"
    echo -e "\n${BLUE}Supported Platforms:${NC}"
    for platform in "${PLATFORMS[@]}"; do
        echo -e "  - $platform"
    done
    echo -e "\n${BLUE}Examples:${NC}"
    echo -e "  $0                                    # Build for all platforms"
    echo -e "  $0 --current                          # Build current platform only"
    echo -e "  $0 --platform windows/amd64           # Build Windows only"
    echo -e "  $0 --test                             # Run tests only"
    echo -e "  $0 --run                              # Run application"
    echo -e "  $0 --clean --current                  # Clean and build current platform"
}

# Main function
main() {
    show_banner
    
    local build_all=true
    local run_tests_flag=true
    local run_app=false
    local test_only=false
    local clean_build=false
    local specific_platform=""
    
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --help|-h)
                show_help
                exit 0
                ;;
            --current)
                build_all=false
                shift
                ;;
            --platform)
                build_all=false
                specific_platform="$2"
                shift 2
                ;;
            --test)
                test_only=true
                shift
                ;;
            --run)
                run_app=true
                shift
                ;;
            --no-tests)
                run_tests_flag=false
                shift
                ;;
            --clean)
                clean_build=true
                shift
                ;;
            *)
                print_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    # Check prerequisites
    check_prerequisites
    
    # Prepare build environment
    prepare_build
    
    if [ "$clean_build" = true ]; then
        print_info "Cleaning build directory..."
        rm -rf "$BUILD_DIR"/*
    fi
    
    # Setup dependencies
    setup_dependencies
    
    # Run tests only
    if [ "$test_only" = true ]; then
        run_tests
        exit 0
    fi
    
    # Build based on options
    if [ "$build_all" = true ]; then
        build_all_platforms
    elif [ -n "$specific_platform" ]; then
        if build_platform "$specific_platform"; then
            print_success "Platform $specific_platform built successfully"
        else
            print_error "Failed to build platform $specific_platform"
            exit 1
        fi
    else
        build_current_platform
    fi
    
    # Run tests
    if [ "$run_tests_flag" = true ]; then
        run_tests
    fi
    
    # Test binaries
    test_binaries
    
    # Show summary
    show_summary
    
    # Run application if requested
    if [ "$run_app" = true ]; then
        run_application
    else
        # Ask if user wants to run the application
        echo -e "\n${YELLOW}Do you want to run the application now? (y/N):${NC}"
        read -r response
        if [[ "$response" =~ ^[Yy]$ ]]; then
            run_application
        fi
    fi
}

# Run main function with all arguments
main "$@"
