#!/bin/bash

# Syntropy Cooperative Grid - Color Definitions
# Version: 2.0.0

# ANSI Color codes for terminal output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
WHITE='\033[1;37m'
GRAY='\033[0;37m'
NC='\033[0m' # No Color

# Background colors
BG_RED='\033[0;41m'
BG_GREEN='\033[0;42m'
BG_YELLOW='\033[0;43m'
BG_BLUE='\033[0;44m'
BG_PURPLE='\033[0;45m'
BG_CYAN='\033[0;46m'

# Text formatting
BOLD='\033[1m'
DIM='\033[2m'
UNDERLINE='\033[4m'
ITALIC='\033[3m'
REVERSE='\033[7m'

# Function to check if terminal supports colors
supports_color() {
    if [ -t 1 ]; then
        local ncolors=$(tput colors 2>/dev/null || echo 0)
        [ "$ncolors" -ge 8 ]
    else
        false
    fi
}

# Disable colors if terminal doesn't support them
if ! supports_color; then
    RED=""
    GREEN=""
    YELLOW=""
    BLUE=""
    PURPLE=""
    CYAN=""
    WHITE=""
    GRAY=""
    NC=""
    BG_RED=""
    BG_GREEN=""
    BG_YELLOW=""
    BG_BLUE=""
    BG_PURPLE=""
    BG_CYAN=""
    BOLD=""
    DIM=""
    UNDERLINE=""
    ITALIC=""
    REVERSE=""
fi

# Color helper functions
print_colored() {
    local color="$1"
    shift
    echo -e "${color}$*${NC}"
}

print_red() {
    print_colored "$RED" "$@"
}

print_green() {
    print_colored "$GREEN" "$@"
}

print_yellow() {
    print_colored "$YELLOW" "$@"
}

print_blue() {
    print_colored "$BLUE" "$@"
}

print_purple() {
    print_colored "$PURPLE" "$@"
}

print_cyan() {
    print_colored "$CYAN" "$@"
}

# Status indication helpers
status_ok() {
    echo -e "[${GREEN}✓${NC}] $*"
}

status_error() {
    echo -e "[${RED}✗${NC}] $*"
}

status_warning() {
    echo -e "[${YELLOW}!${NC}] $*"
}

status_info() {
    echo -e "[${BLUE}i${NC}] $*"
}

# Progress indicators
show_loading() {
    local message="$1"
    local delay=${2:-0.1}
    local spinstr='|/-\'
    
    printf "%s " "$message"
    while true; do
        local temp=${spinstr#?}
        printf "[%c]" "$spinstr"
        local spinstr=$temp${spinstr%"$temp"}
        sleep $delay
        printf "\b\b\b"
    done
}

stop_loading() {
    printf "\b\b\b   \b\b\b"
}

# Box drawing helpers
print_box() {
    local title="$1"
    local content="$2"
    local color="${3:-$CYAN}"
    
    local title_len=${#title}
    local box_width=$((title_len + 4))
    
    # Top border
    echo -e "${color}╔$(printf '%*s' $((box_width-2)) | tr ' ' '═')╗${NC}"
    
    # Title
    echo -e "${color}║ ${BOLD}$title${NC}${color} ║${NC}"
    
    # Content separator
    echo -e "${color}╠$(printf '%*s' $((box_width-2)) | tr ' ' '═')╣${NC}"
    
    # Content lines
    while IFS= read -r line; do
        local line_len=${#line}
        local padding=$((box_width - line_len - 2))
        echo -e "${color}║${NC} $line$(printf '%*s' $padding)${color}║${NC}"
    done <<< "$content"
    
    # Bottom border
    echo -e "${color}╚$(printf '%*s' $((box_width-2)) | tr ' ' '═')╝${NC}"
}

# Banner helpers
print_banner() {
    local text="$1"
    local color="${2:-$PURPLE}"
    
    echo -e "${color}"
    echo "╔════════════════════════════════════════════════════════════════════════════╗"
    printf "║%*s║\n" 76 ""
    printf "║%*s%s%*s║\n" $(((76-${#text})/2)) "" "$text" $(((76-${#text}+1)/2)) ""
    printf "║%*s║\n" 76 ""
    echo "╚════════════════════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

# Table helpers
print_table_header() {
    local -a headers=("$@")
    local total_width=0
    local col_widths=()
    
    # Calculate column widths
    for header in "${headers[@]}"; do
        local width=$((${#header} + 4))
        col_widths+=($width)
        total_width=$((total_width + width))
    done
    
    # Print header separator
    printf "+"
    for width in "${col_widths[@]}"; do
        printf "%*s+" $((width-1)) "" | tr ' ' '-'
    done
    echo ""
    
    # Print headers
    printf "|"
    for i in "${!headers[@]}"; do
        local header="${headers[$i]}"
        local width="${col_widths[$i]}"
        printf " %-*s|" $((width-2)) "$header"
    done
    echo ""
    
    # Print separator
    printf "+"
    for width in "${col_widths[@]}"; do
        printf "%*s+" $((width-1)) "" | tr ' ' '-'
    done
    echo ""
}

print_table_row() {
    local -a cols=("$@")
    
    printf "|"
    for col in "${cols[@]}"; do
        printf " %-15s|" "$col"
    done
    echo ""
}

# Progress bar
print_progress_bar() {
    local current=$1
    local total=$2
    local message="$3"
    local width=${4:-50}
    local color="${5:-$GREEN}"
    
    local percentage=$((current * 100 / total))
    local filled=$((current * width / total))
    local empty=$((width - filled))
    
    printf "\r%s [" "$message"
    printf "%s" "${color}"
    printf "%*s" "$filled" | tr ' ' '█'
    printf "%s" "${NC}"
    printf "%*s" "$empty" | tr ' ' '░'
    printf "] %d%%" "$percentage"
    
    if [ "$current" -eq "$total" ]; then
        echo ""
    fi
}

# Section dividers
section_divider() {
    local title="$1"
    local color="${2:-$CYAN}"
    
    echo -e "\n${color}═══ $title ═══${NC}\n"
}

subsection_divider() {
    local title="$1"
    local color="${2:-$BLUE}"
    
    echo -e "\n${color}--- $title ---${NC}\n"
}

# Highlight important text
highlight() {
    local text="$1"
    local color="${2:-$YELLOW}"
    
    echo -e "${color}${BOLD}$text${NC}"
}

# Color test function
test_colors() {
    echo "Color Test - Syntropy Cooperative Grid"
    echo "======================================"
    echo ""
    
    echo "Basic Colors:"
    print_red "Red text"
    print_green "Green text"
    print_yellow "Yellow text"
    print_blue "Blue text"
    print_purple "Purple text"
    print_cyan "Cyan text"
    echo ""
    
    echo "Status Indicators:"
    status_ok "Success status"
    status_error "Error status"
    status_warning "Warning status"
    status_info "Information status"
    echo ""
    
    echo "Progress Bar:"
    for i in {0..10}; do
        print_progress_bar $i 10 "Loading" 30
        sleep 0.1
    done
    echo ""
    
    echo "Section Dividers:"
    section_divider "Main Section"
    subsection_divider "Sub Section"
    echo ""
    
    echo "Highlighted Text:"
    highlight "Important information"
    echo ""
    
    echo "Banner:"
    print_banner "SYNTROPY GRID"
    echo ""
}

# Export color variables for use in other scripts
export RED GREEN YELLOW BLUE PURPLE CYAN WHITE GRAY NC
export BG_RED BG_GREEN BG_YELLOW BG_BLUE BG_PURPLE BG_CYAN
export BOLD DIM UNDERLINE ITALIC REVERSE