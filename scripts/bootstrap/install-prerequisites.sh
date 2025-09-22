#!/bin/bash

# Syntropy Cooperative Grid - Prerequisites Installation
# Version: 1.0.0

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}Installing prerequisites for Syntropy USB Creator...${NC}"

# Detect OS
if [ -f /etc/os-release ]; then
    . /etc/os-release
    OS=$ID
    VER=$VERSION_ID
else
    echo -e "${RED}Cannot detect OS${NC}"
    exit 1
fi

# Detect WSL
is_wsl() {
    if grep -qi microsoft /proc/version; then
        return 0  # True, is WSL
    else
        return 1  # False, not WSL
    fi
}

# Install WSL-specific requirements if needed
if is_wsl; then
    echo -e "${BLUE}Installing WSL-specific requirements...${NC}"
    
    # Check if PowerShell is accessible
    if ! command -v powershell.exe &> /dev/null; then
        echo -e "${RED}Error: PowerShell not found. Please ensure you're running in WSL2${NC}"
        exit 1
    fi
    
    # Additional WSL utilities
    case $OS in
        ubuntu|debian)
            sudo apt-get install -y \
                jq \
                bc \
                dos2unix \
                wslu
            ;;
        fedora|rhel|centos)
            sudo dnf install -y \
                jq \
                bc \
                dos2unix
            ;;
    esac
    
    echo -e "${GREEN}WSL-specific requirements installed${NC}"
fi

# Function to check command existence
check_command() {
    command -v "$1" >/dev/null 2>&1
}

# Function to install packages based on OS
install_packages() {
    local os=$1
    shift
    local packages=("$@")
    
    case $os in
        ubuntu|debian)
            echo -e "${BLUE}Installing packages for Ubuntu/Debian...${NC}"
            if ! check_command apt-get; then
                echo -e "${RED}apt-get not found. Please install manually.${NC}"
                exit 1
            fi
            sudo apt-get update
            sudo apt-get install -y "${packages[@]}"
            ;;
        
        fedora)
            echo -e "${BLUE}Installing packages for Fedora...${NC}"
            if ! check_command dnf; then
                echo -e "${RED}dnf not found. Please install manually.${NC}"
                exit 1
            fi
            sudo dnf install -y "${packages[@]}"
            ;;
        
        rhel|centos)
            echo -e "${BLUE}Installing packages for RHEL/CentOS...${NC}"
            if ! check_command yum; then
                echo -e "${RED}yum not found. Please install manually.${NC}"
                exit 1
            fi
            sudo yum install -y "${packages[@]}"
            ;;
        
        *)
            echo -e "${RED}Unsupported OS: $OS${NC}"
            echo "Please install the following packages manually:"
            printf '%s\n' "${packages[@]}"
            exit 1
            ;;
    esac
}

# Define common packages
COMMON_PACKAGES=(
    curl
    wget
    git
    jq
    python3
    python3-pip
    openssh-client
    openssl
    parted
    dosfstools
    gdisk
    nmap
    bc
)

# Define OS-specific packages
DEBIAN_PACKAGES=(
    "${COMMON_PACKAGES[@]}"
    syslinux
    syslinux-utils
    isolinux
    xorriso
    genisoimage
    cloud-init
    whois # for mkpasswd
)

REDHAT_PACKAGES=(
    "${COMMON_PACKAGES[@]}"
    syslinux
    xorriso
    genisoimage
    cloud-init
    python3-policycoreutils # for selinux
)

# Install packages based on OS
case $OS in
    ubuntu|debian)
        install_packages "$OS" "${DEBIAN_PACKAGES[@]}"
        ;;
    
    fedora|rhel|centos)
        install_packages "$OS" "${REDHAT_PACKAGES[@]}"
        ;;
    
    *)
        echo -e "${RED}Unsupported OS: $OS${NC}"
        exit 1
        ;;
esac

# Verify Python packages
echo -e "${BLUE}Installing required Python packages...${NC}"
pip3 install --user PyYAML requests cryptography

# Verify tool versions
echo -e "${BLUE}Verifying installed tools...${NC}"
for cmd in curl wget git python3 openssl parted mkfs.fat dd; do
    if check_command "$cmd"; then
        echo -e "${GREEN}✓ $cmd installed${NC}"
    else
        echo -e "${RED}✗ $cmd not found${NC}"
        MISSING_TOOLS=1
    fi
done

# Final status
if [ "${MISSING_TOOLS:-0}" -eq 1 ]; then
    echo -e "${RED}Some required tools are missing. Please install them manually.${NC}"
    exit 1
else
    echo -e "${GREEN}All prerequisites installed successfully!${NC}"
fi