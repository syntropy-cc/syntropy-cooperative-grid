#!/bin/bash

# Syntropy Cooperative Grid - Logging System
# Version: 2.0.0

# Logging configuration
LOG_LEVEL=${LOG_LEVEL:-"INFO"}
LOG_FILE=${LOG_FILE:-"/tmp/syntropy-usb-creator.log"}
ENABLE_FILE_LOGGING=${ENABLE_FILE_LOGGING:-false}

# Log levels (numeric for comparison)
declare -A LOG_LEVELS=(
    ["DEBUG"]=0
    ["INFO"]=1
    ["WARN"]=2
    ["ERROR"]=3
    ["SUCCESS"]=1
)

# Get current timestamp
get_timestamp() {
    date '+%Y-%m-%d %H:%M:%S'
}

# Check if log level should be shown
should_log() {
    local level="$1"
    local current_level_num=${LOG_LEVELS[$LOG_LEVEL]:-1}
    local message_level_num=${LOG_LEVELS[$level]:-1}
    
    [ $message_level_num -ge $current_level_num ]
}

# Main logging function
log() {
    local level="$1"
    shift
    local message="$*"
    local timestamp=$(get_timestamp)
    
    # Check if we should log this level
    if ! should_log "$level"; then
        return 0
    fi
    
    # Format message for console
    local console_message=""
    case "$level" in
        DEBUG)
            console_message="${GRAY}[DEBUG]${NC} $message"
            ;;
        INFO)
            console_message="${BLUE}[INFO]${NC} $message"
            ;;
        WARN)
            console_message="${YELLOW}[WARN]${NC} $message"
            ;;
        ERROR)
            console_message="${RED}[ERROR]${NC} $message"
            ;;
        SUCCESS)
            console_message="${GREEN}[SUCCESS]${NC} $message"
            ;;
        *)
            console_message="[$level] $message"
            ;;
    esac
    
    # Output to console
    echo -e "$console_message"
    
    # Output to file if enabled
    if [ "$ENABLE_FILE_LOGGING" = true ]; then
        echo "[$timestamp] [$level] $message" >> "$LOG_FILE"
    fi
}

# Convenience functions
log_debug() {
    log DEBUG "$@"
}

log_info() {
    log INFO "$@"
}

log_warn() {
    log WARN "$@"
}

log_error() {
    log ERROR "$@"
}

log_success() {
    log SUCCESS "$@"
}

# Progress indicator functions
show_spinner() {
    local pid=$1
    local message="$2"
    local delay=0.1
    local spinstr='|/-\'
    
    while [ "$(ps a | awk '{print $1}' | grep $pid)" ]; do
        local temp=${spinstr#?}
        printf " [%c] %s\r" "$spinstr" "$message"
        local spinstr=$temp${spinstr%"$temp"}
        sleep $delay
    done
    printf "    \r"
}

# Progress bar function
show_progress() {
    local current=$1
    local total=$2
    local message="$3"
    local width=50
    
    local percentage=$((current * 100 / total))
    local completed=$((current * width / total))
    local remaining=$((width - completed))
    
    printf "\r%s [" "$message"
    printf "%${completed}s" | tr ' ' '='
    printf "%${remaining}s" | tr ' ' '-'
    printf "] %d%%" "$percentage"
    
    if [ $current -eq $total ]; then
        echo ""
    fi
}

# Error handling
handle_error() {
    local exit_code=$?
    local line_number=$1
    local command="$2"
    
    if [ $exit_code -ne 0 ]; then
        log ERROR "Command failed with exit code $exit_code on line $line_number: $command"
        
        # Log additional context if available
        if [ -n "${FUNCNAME[1]:-}" ]; then
            log ERROR "Function: ${FUNCNAME[1]}"
        fi
        
        # Log to file for debugging
        if [ "$ENABLE_FILE_LOGGING" = true ]; then
            echo "--- Error Context ---" >> "$LOG_FILE"
            echo "Exit Code: $exit_code" >> "$LOG_FILE"
            echo "Line: $line_number" >> "$LOG_FILE"
            echo "Command: $command" >> "$LOG_FILE"
            echo "Function Stack: ${FUNCNAME[*]}" >> "$LOG_FILE"
            echo "-------------------" >> "$LOG_FILE"
        fi
    fi
    
    return $exit_code
}

# Set up error trapping
trap 'handle_error ${LINENO} "$BASH_COMMAND"' ERR

# Initialize logging
init_logging() {
    if [ "$ENABLE_FILE_LOGGING" = true ]; then
        # Create log file directory if it doesn't exist
        local log_dir=$(dirname "$LOG_FILE")
        mkdir -p "$log_dir"
        
        # Initialize log file
        echo "=== Syntropy USB Creator Log - $(get_timestamp) ===" > "$LOG_FILE"
        log INFO "Logging initialized - file: $LOG_FILE"
    fi
}