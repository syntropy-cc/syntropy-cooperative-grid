#!/bin/bash

# Syntropy Cooperative Grid - Progress Bar and User Feedback
# Version: 1.0.0

# Configuração da barra de progresso
PROGRESS_WIDTH=50
PROGRESS_CHAR="="
PROGRESS_EMPTY=" "
PROGRESS_START="["
PROGRESS_END="]"

# Variáveis de estado
current_progress=0
current_message=""
start_time=0

# Inicializar barra de progresso
init_progress() {
    current_progress=0
    start_time=$(date +%s)
    echo -e "\n"  # Linha em branco antes da barra
}

# Atualizar barra de progresso
update_progress() {
    local progress=$1
    local message=$2
    local force=${3:-false}
    
    # Validar progresso
    if [ "$progress" -lt 0 ]; then progress=0; fi
    if [ "$progress" -gt 100 ]; then progress=100; fi
    
    # Se não houver mudança e não forçado, sair
    if [ "$progress" -eq "$current_progress" ] && [ "$message" = "$current_message" ] && [ "$force" = false ]; then
        return
    fi
    
    current_progress=$progress
    current_message=$message
    
    # Calcular número de caracteres preenchidos
    local filled=$(( progress * PROGRESS_WIDTH / 100 ))
    local empty=$(( PROGRESS_WIDTH - filled ))
    
    # Calcular tempo decorrido
    local current_time=$(date +%s)
    local elapsed=$((current_time - start_time))
    local elapsed_str=$(date -u -d "@$elapsed" +"%M:%S")
    
    # Construir barra
    printf "\r%-50s" "$message"
    printf "%s" "$PROGRESS_START"
    printf "%${filled}s" | tr " " "$PROGRESS_CHAR"
    printf "%${empty}s" | tr " " "$PROGRESS_EMPTY"
    printf "%s" "$PROGRESS_END"
    printf " %3d%% (%s)" "$progress" "$elapsed_str"
    
    # Se 100%, adicionar nova linha
    if [ "$progress" -eq 100 ]; then
        echo -e "\nConcluído em $elapsed_str\n"
    fi
}

# Mostrar mensagem de status
show_status() {
    local message=$1
    local type=${2:-info}
    
    case $type in
        info)
            echo -e "\n${BLUE}ℹ️  $message${NC}"
            ;;
        success)
            echo -e "\n${GREEN}✅ $message${NC}"
            ;;
        warning)
            echo -e "\n${YELLOW}⚠️  $message${NC}"
            ;;
        error)
            echo -e "\n${RED}❌ $message${NC}"
            ;;
    esac
}

# Solicitar confirmação do usuário
confirm_action() {
    local message=$1
    local default=${2:-n}
    
    local prompt
    if [ "$default" = "y" ]; then
        prompt="Y/n"
    else
        prompt="y/N"
    fi
    
    echo -e "\n${YELLOW}⚠️  $message${NC}"
    read -p "Confirmar? [$prompt]: " -n 1 -r
    echo
    
    if [ "$default" = "y" ]; then
        [[ ! $REPLY =~ ^[Nn]$ ]]
    else
        [[ $REPLY =~ ^[Yy]$ ]]
    fi
}

# Exemplo de uso:
# init_progress
# update_progress 0 "Iniciando criação do USB..."
# update_progress 20 "Preparando partições..."
# update_progress 40 "Copiando arquivos..."
# update_progress 60 "Configurando bootloader..."
# update_progress 80 "Gerando chaves..."
# update_progress 100 "Concluído"
#
# show_status "USB criado com sucesso!" success
# show_status "Atenção: faça backup dos dados" warning
# show_status "Erro ao montar dispositivo" error
#
# if confirm_action "Deseja formatar o dispositivo?"; then
#     echo "Formatando..."
# fi