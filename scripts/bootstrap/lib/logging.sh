#!/bin/bash

# Syntropy Cooperative Grid - Enhanced Logging System
# Version: 2.1.0

# Default log configuration
SYNTROPY_LOG_DIR="${HOME}/.syntropy/logs"
LOG_LEVEL=${LOG_LEVEL:-"INFO"}
LOG_FILE=${LOG_FILE:-"${SYNTROPY_LOG_DIR}/syntropy-$(date +%Y%m%d).log"}
ENABLE_FILE_LOGGING=${ENABLE_FILE_LOGGING:-true}
MAX_LOG_SIZE_MB=10
MAX_LOG_FILES=7

# Ensure log directory exists
mkdir -p "${SYNTROPY_LOG_DIR}"

# Log levels with numeric values and colors
declare -A LOG_LEVELS=(
    ["TRACE"]=0
    ["DEBUG"]=1
    ["INFO"]=2
    ["WARN"]=3
    ["ERROR"]=4
    ["FATAL"]=5
    ["SUCCESS"]=2
)

declare -A LOG_COLORS=(
    ["TRACE"]="${GRAY}"
    ["DEBUG"]="${CYAN}"
    ["INFO"]="${BLUE}"
    ["WARN"]="${YELLOW}"
    ["ERROR"]="${RED}"
    ["FATAL"]="${RED}${BOLD}"
    ["SUCCESS"]="${GREEN}"
)

# Initialize logging
init_logging() {
    # Rotate logs if needed
    if [ -f "$LOG_FILE" ]; then
        local size_mb=$(du -m "$LOG_FILE" | cut -f1)
        if [ "$size_mb" -ge "$MAX_LOG_SIZE_MB" ]; then
            rotate_logs
        fi
    fi
    
    # Create log file if it doesn't exist
    if [ ! -f "$LOG_FILE" ]; then
        touch "$LOG_FILE"
        log INFO "Initialized new log file: $LOG_FILE"
    fi
}

# Rotate log files
rotate_logs() {
    for i in $(seq $((MAX_LOG_FILES-1)) -1 0); do
        [ -f "${LOG_FILE}.$i" ] && mv "${LOG_FILE}.$i" "${LOG_FILE}.$((i+1))"
    done
    [ -f "$LOG_FILE" ] && mv "$LOG_FILE" "${LOG_FILE}.0"
    touch "$LOG_FILE"
}

# Get current timestamp with microseconds
get_timestamp() {
    date '+%Y-%m-%d %H:%M:%S.%3N'
}

# Verificar se o nível de log deve ser mostrado
should_log() {
    local level="$1"
    local current_level_num=${LOG_LEVELS[$LOG_LEVEL]:-2}  # Default to INFO level
    local message_level_num=${LOG_LEVELS[$level]:-2}
    
    [ $message_level_num -ge $current_level_num ]
}

# Format log message
format_log_message() {
    local level="$1"
    local message="$2"
    local timestamp="$3"
    local source="${4:-unknown}"
    local line="${5:-0}"
    
    echo "[$timestamp] [$level] [$source:$line] $message"
}

# Main logging function with enhanced features
log() {
    local level="$1"
    shift
    local message="$*"
    local timestamp=$(get_timestamp)
    
    # Get caller information
    local source_info=""
    if [ "${BASH_SOURCE[1]:-}" ]; then
        source_info="$(basename "${BASH_SOURCE[1]}")"
        local caller_line="${BASH_LINENO[0]}"
    else
        source_info="main"
        local caller_line="0"
    fi
    
    # Check if we should log this level
    if ! should_log "$level"; then
        return 0
    fi
    
    # Format messages
    local log_message=$(format_log_message "$level" "$message" "$timestamp" "$source_info" "$caller_line")
    local console_message=""
    
    # Apply color formatting for console
    if [ -n "${LOG_COLORS[$level]:-}" ]; then
        console_message="${LOG_COLORS[$level]}[${level}]${NC} $message"
    else
        # Fallback colors if not defined in LOG_COLORS
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
    fi
    
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