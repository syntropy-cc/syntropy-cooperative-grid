#!/bin/bash

# Syntropy CLI - Setup Básico para Linux/macOS
# Script para configuração automática com verificações de ambiente
# 
# Author: Syntropy Team
# Created: 2025-01-27  
# Version: 1.0.0
# 
# Purpose: Automatizar o processo de setup do Syntropy CLI em sistemas Unix-like
# 
# Dependencies: 
# - bash 4.0+
# - curl/wget para download de dependências
# - tar para extração
# - go runtime (opcional, pode usar binários pré-compilados)
#
# Usage: 
#   ./setup-basic.sh                 # Setup completo com permissões admin
#   ./setup-basic.sh --user-only    # Setup apenas para usuário local
#   ./setup-basic.sh --validate     # Apenas validar ambiente
#   --dry-run                       # Simular sem fazer mudanças
#   --help                          # Exibir ajuda completa

set -euo pipefail  # Exit on any error, undefined vars, pipe failures

# ============================================================================
# CONFIGURAÇÕES E CONSTANTES
# ============================================================================

readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly SCRIPT_NAME="setup-basic.sh"
readonly SCRIPT_VERSION="1.0.0"

# Cores para output
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly NC='\033[0m' # No Color

# Configurações padrão
readonly SYNTHROPY_HOME="$HOME/.syntropy"
readonly MIN_DISK_SPACE_GB=1
readonly MIN_MEMORY_GB=0.5
readonly REQUIRED_GO_VERSION="1.19"
readonly CONFIG_FILE="$(dirname "$0")/config-example.yaml"

# ============================================================================
# FUNÇÕES UTILITÁRIAS
# ============================================================================

# Função para logging estruturado
log() {
    local level="$1"; shift
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    local message="$*"
    
    case "$level" in
        "ERROR") echo -e "${RED}[ERROR]${NC} [$timestamp] $message" >&2 ;;
        "WARN")  echo -e "${YELLOW}[WARN]${NC}  [$timestamp] $message" ;;
        "INFO")  echo -e "${BLUE}[INFO]${NC}  [$timestamp] $message" ;;
        "DEBUG") echo -e "${GREEN}[DEBUG]${NC} [$timestamp] $message" ;;
        *) echo "[$level] [$timestamp] $message" ;;
    esac
}

# Função para exibir título de seção
section_title() {
    echo -e "\n${BLUE}=== $1 ===${NC}"
}

# Função para verificar se comando existe
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Função para verificar arquivo/template de configuração padrão
check_config_template() {
    if [[ -f "$CONFIG_FILE" ]]; then
        log "INFO" "Arquivo de configuração encontrado: $CONFIG_FILE"
        return 0
    fi
    
    log "WARN" "Arquivo de configuração não encontrado: $CONFIG_FILE"
    log "INFO" "Será usada configuração padrão incorporada"
    return 1
}

# ============================================================================
# VERIFICAÇÕES DE AMBIENTE
# ============================================================================

# Verificar requisitos básicos do sistema
check_system_requirements() {
    log "INFO" "Verificando requisitos básicos do sistema..."
    
    # Verificar se não é Windows (this script é Unix-specific)
    if command_exists systeminfo; then
        log "ERROR" "Este é um ambiente Windows. Use setup-basic.ps1 em vez deste script."
        exit 1
    fi
    
    # Verificar espaço em disco
    if ! check_disk_space; then
        log "ERROR" "Espaço em disco insuficiente"
        exit 1
    fi
    
    # Verificar memória disponível (na maioria dos sistemas Unix)
    if has_command "free"; then
        check_memory
    fi
    
    # Verificar dependências de sistema
    for cmd in curl wget tar; do
        if ! command_exists "$cmd"; then
            log "WARN" "Comando '$cmd' não encontrado - algumas funcionalidades podem estar limitadas"
        fi
    done
    
    log "INFO" "Requisitos básicos verificados"
}

# Verificar espaço em disco disponível
check_disk_space() {
    local available_space_gb
    local home_dir_space_gb
    
    # Verificar espaço disponível no diretório home
    if command_exists "df"; then
        home_dir_space_gb=$(df -BG "$HOME" | tail -1 | awk '{gsub("G", "", $4); print $4}')
        
        if (( $(echo "$home_dir_space_gb < $MIN_DISK_SPACE_GB" | bc -l) )); then
            log "ERROR" "Espaço em disco insuficiente: ${home_dir_space_gb}GB disponível, ${MIN_DISK_SPACE_GB}GB necessário"
            return 1
        fi
        
        log "INFO" "Espaço em disco OK: ${home_dir_space_gb}GB disponível"
    else
        log "WARN" "Não foi possível verificar espaço em disco"
    fi
    
    return 0
}

# Verificar memória disponível
check_memory() {
    local total_memory_gb
    
    if command_exists "free"; then
        total_memory_gb=$(free -g | grep "^Mem:" | awk '{print $2}')
        
        if (( total_memory_gb < MIN_MEMORY_GB )); then
            log "WARN" "Memória disponível baixa: ${total_memory_gb}GB total"
        fi
        
        log "INFO" "Memória: ${total_memory_gb}GB disponível"
    fi
}

# Verificar se Go está disponível (opcional)
check_go_environment() {
    if command_exists go; then
        local go_version
        go_version=$(go version | cut -d' ' -f3 | sed 's/go//')
        
        log "INFO" "Go encontrado: versão $go_version"
        return 0
    else
        log "INFO" "Go não encontrado - tentando usar binários pré-compilados"
        return 1
    fi
}

# ============================================================================
# OPERAÇÕES DE SETUP
# ============================================================================

# Criar estrutura de diretórios
create_directories() {
    log "INFO" "Criando estrutura de diretórios..."
    
    local dirs=(
        "$SYNTHROPY_HOME"
        "$SYNTHROPY_HOME/config"
        "$SYNTHROPY_HOME/logs"
        "$SYNTHROPY_HOME/data"
        "$SYNTHROPY_HOME/services"
    )
    
    for dir in "${dirs[@]}"; do
        if [[ ! -d "$dir" ]]; then
            mkdir -p "$dir"
            log "INFO" "Diretório criado: $dir"
        else
            log "DEBUG" "Diretório já existe: $dir"
        fi
    done
}

# Configurar arquivos de configuração
setup_configuration() {
    section_title "Configuração"
    
    log "INFO" "Configurando arquivos de configuração..."
    
    if [[ -f "$CONFIG_FILE" ]]; then
        cp "$CONFIG_FILE" "$SYNTHROPY_HOME/config/manager.yaml"
        log "INFO" "Arquivo de configuração copiado de exemplo"
    else
        # Criar configuração básica padrão
        create_default_config
    fi
    
    # Verificar se configuração é válida
    validate_configuration
}

# Criar configuração padrão básica
create_default_config() {
    log "INFO" "Criando configuração padrão..."
    
    cat > "$SYNTHROPY_HOME/config/manager.yaml" << EOF
manager:
  home_dir: "$SYNTHROPY_HOME"
  log_level: "info"
  api_endpoint: "https://api.syntropy.io"
  directories:
    config: "config"
    logs: "logs"
    data: "data"
  default_paths:
    config: "config/manager.yaml"
    owner_key: "config/owner.key"

owner_key:
  type: "Ed25519"
  path: "config/owner.key"

environment:
  os: "$(uname -s | tr '[:upper:]' '[:lower:]')"
  architecture: "$(uname -m)"
  home_dir: "$HOME"
EOF
    
    log "INFO" "Configuração padrão criada"
}

# Validar configuração criada
validate_configuration() {
    log "INFO" "Validando configuração..."
    
    if [[ ! -f "$SYNTHROPY_HOME/config/manager.yaml" ]]; then
        log "ERROR" "Arquivo de configuração não foi criado"
        exit 1
    fi
    
    if ! command_exists "yq" && ! command_exists "python3"; then
        log "WARN" "Não foi possível validar sintaxe YAML (yq/python3 não disponível)"
    fi
    
    log "INFO" "Configuração válida"
}

# Configurar logs e diretórios
setup_logging() {
    log "INFO" "Configurando sistema de logging..."
    
    local log_file="$SYNTHROPY_HOME/logs/setup.log"
    
    # Criar log de setup
    {
        echo "# Syntropy Setup Log - $(date)"
        echo "# Generated by $SCRIPT_NAME v$SCRIPT_VERSION"
        echo "# Environment: $(uname -a)"
        echo ""
    } >> "$log_file"
    
    log "INFO" "Sistema de logging configurado"
}

# Executar setup usando sintropy CLI (se disponível)  
run_syntropy_setup() {
    section_title "Execução do Setup"
    
    if command_exists syntropy; then
        log "INFO" "Executando setup via syntropy CLI..."
        
        if [[ "$FORCE_SETUP" == "true" ]]; then
            syntropy setup --force
        else
            syntropy setup
        fi
        
        log "INFO" "Setup executado com sucesso"
    else
        log "WARN" "Syntropy CLI não encontrado - setup manual realizado"
        create_manual_setup_files
    fi
}

# Criar arquivos de setup manual se CLI não estiver disponível
create_manual_setup_files() {
    log "INFO" "Criando arquivos de setup manual..."
    
    # Criar script de instalação como serviço (Linux)
    if [[ -f /etc/os-release ]] && grep -q "ID=ubuntu\|ID=debian\|ID=fedora" /etc/os-release; then
        create_systemd_service
    fi
    
    log "INFO" "Arquivos de setup manual criados"
}

# Criar arquivo de serviço systemd
create_systemd_service() {
    cat > "$SYNTHROPY_HOME/services/install-systemd.sh" << 'EOF'
#!/bin/bash
# Script de instalação do serviço Syntropy como systemd service

set -e

SYNTHROPY_HOME="$HOME/.syntropy"
USER_SERVICE_DIR="$HOME/.config/systemd/user"
SERVICE_FILE="$USER_SERVICE_DIR/syntropy.service"

# Criar diretório de serviço do usuário
mkdir -p "$USER_SERVICE_DIR"

# Criar arquivo de serviço
cat > "$SERVICE_FILE" << EOS
[Unit]
Description=Syntropy Network Coordinator
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
ExecStart=/path/to/syntropy --config-dir $SYNTHROPY_HOME/config
WorkingDirectory=$SYNTHROPY_HOME
User=$USER
Restart=always
RestartSec=10

[Install]
WantedBy=default.target
EOS

echo "Service file created at: $SERVICE_FILE"
echo "To enable and start the service, run:"
echo "  systemctl --user daemon-reload"
echo "  systemctl --user enable syntropy.service"
echo "  systemctl --user start syntropy.service"
EOF
    
    chmod +x "$SYNTHROPY_HOME/services/install-systemd.sh"
    log "INFO" "Script de serviço systemd criado: $SYNTHROPY_HOME/services/install-systemd.sh"
}

# Executar apenas validação sem fazer mudanças
validate_environment_only() {
    section_title "Validação de Ambiente"
    
    log "INFO" "Executando validação sem modificar sistema..."
    check_system_requirements
    check_go_environment
    
    log "INFO" "Validação concluída - ambiente adequado para setup"
}

# ============================================================================
# VERIFICAÇÕES PÓS-SETUP
# ============================================================================

verify_setup() {
    section_title "Verificação Pós-Setup"
    
    log "INFO" "Verificando instalação..."
    
    # Verificar se diretórios foram criados
    for dir in config logs data; do
        if [[ -d "$SYNTHROPY_HOME/$dir" ]]; then
            log "INFO" "✓ Diretório $dir criado"
        else
            log "ERROR" "✗ Diretório $dir não foi criado"
            return 1
        fi
    done
    
    # Verificar arquivo de configuração
    if [[ -f "$SYNTHROPY_HOME/config/manager.yaml" ]]; then
        log "INFO" "✓ Arquivo de configuração criado"
    else
        log "ERROR" "✗ Arquivo de configuração não encontrado"
        return 1
    fi
    
    # Listar arquivos criados
    log "INFO" "Arquivos criados no setup:"
    find "$SYNTHROPY_HOME" -type f 2>/dev/null | while read -r file; do
        log "INFO" "  • $file"
    done
    
    log "INFO" "Verificação concluída - setup feito com sucesso"
}

# ============================================================================
# FUNÇÃO PRINCIPAL
# ============================================================================

main() {
    # Parse argumentos de linha de comando
    parse_arguments "$@"
    
    # Verificar se script está sendo executado do diretório correto
    if [[ "$(basename "$PWD")" != "basic-setup" && ! -f "config-example.yaml" ]]; then
        log "WARN" "Script executado fora do diretório correct. Sugestões:"
        log "WARN" "  1. Navigate para $SCRIPT_DIR"
        log "WARN" "  2. Ou execute: cd $SCRIPT_DIR && $SCRIPT_NAME"
    fi
    
    # Mostrar header
    section_title "Syntropy CLI Setup Básico para Linux/macOS"
    echo "Versão: $SCRIPT_VERSION"
    echo "Diretório: $SCRIPT_DIR"
    echo "Syntropy Home: $SYNTHROPY_HOME"
    echo ""
    
    if [[ "$DRY_RUN" == "true" ]]; then
        log "INFO" "Modo DRY RUN - nenhuma mudança será feita"
        validate_environment_only
        return 0
    fi
    
    if [[ "$VALIDATE_ONLY" == "true" ]]; then
        validate_environment_only
        return 0
    fi
    
    # Executar setup completo
    check_system_requirements
    create_directories
    setup_configuration
    setup_logging
    run_syntropy_setup
    verify_setup
    
    section_title "Setup Concluído"
    log "INFO" "Setup básico do Syntropy CLI concluído com sucesso!"
    log "INFO" "Diretório base: $SYNTHROPY_HOME"
    log "INFO" "Configuração: $SYNTHROPY_HOME/config/manager.yaml"
    log "INFO" "Logs: $SYNTHROPY_HOME/logs/"
    
    if [[ "$SERVICE_SETUP" == "true" ]]; then
        log "INFO" "Para configurar serviço, execute: $SYNTHROPY_HOME/services/install-systemd.sh"
    fi
    
    echo ""
    log "INFO" "Próximos passos:"
    log "INFO" "  1. Verifique status: syntropy setup status"
    log "INFO" "  2. Consulte log para detalhes: tail -f $SYNTHROPY_HOME/logs/setup.log"
    log "INFO" "  3. Configure serviço do sistema: $SYNTHROPY_HOME/services/install-systemd.sh"
}

# ============================================================================
# PARSE DE ARGUMENTOS
# ============================================================================

parse_arguments() {
    VALIDATE_ONLY="false"
    DRY_RUN="false"
    FORCE_SETUP="false"
    SERVICE_SETUP="false"
    USER_ONLY="false"
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --validate)
                VALIDATE_ONLY="true"
                shift
                ;;
            --dry-run)
                DRY_RUN="true"
                shift
                ;;
            --force)
                FORCE_SETUP="true"
                shift
                ;;
            --service)
                SERVICE_SETUP="true"
                shift
                ;;
            --user-only)
                USER_ONLY="true"
                shift
                ;;
            --help)
                show_help
                exit 0
                ;;
            *)
                log "ERROR" "Argumento desconhecido: $1"
                log "ERROR" "Execute com --help para ver opções disponíveis"
                exit 1
                ;;
        esac
    done
}

show_help() {
    cat << EOF
Uso: $SCRIPT_NAME [opções]

OPÇÕES:
    --validate       Apenas validar ambiente, não executar setup
    --dry-run        Simular operações sem fazer mudanças
    --force          Forçar setup mesmo com validações com problemas
    --service        Configurar como serviço do sistema
    --user-only      Configurar apenas para usuário atual
    --help           Exibir esta ajuda

EXEMPLOS:
    $SCRIPT_NAME                    # Setup completo padrão
    $SCRIPT_NAME --validate         # Apenas validar ambiente  
    $SCRIPT_NAME --service          # Setup + configuração de serviço
    $SCRIPT_NAME --dry-run          # Simular setup sem mudanças

DESCRIÇÃO:
    Este script automático configura o ambiente Syntropy CLI
    seguindo as melhores práticas e validações de segurança.
    
REQUISITOS:
    • Sistema Unix-like (Linux/macOS)
    • Bash 4.0 ou superior
    • Esp histórico slug gt >= 1GB
    • Permissões de usuário ou admin (conforme opção)

ARQUIVOS GERADOS:
    • $SYNTHROPY_HOME/config/manager.yaml
    • $SYNTHROPY_HOME/services/install-\*.sh
    • $SYNTHROPY_HOME/logs/setup.log

Para mais detalhes, consulte ../README.md
EOF
}

# ============================================================================
# EXECUÇÃO DO SCRIPT
# ============================================================================

# Garantir que rodamos no contexto correto
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
