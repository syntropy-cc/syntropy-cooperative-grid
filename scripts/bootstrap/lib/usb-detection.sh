#!/bin/bash

# Check if running in WSL
is_wsl() {
    grep -qi microsoft /proc/version
}

# Syntropy Cooperative Grid - USB Detection and Safety
# Version: 2.0.0 - Enhanced USB Detection

# USB detection configuration
USB_MIN_SIZE_MB=1024  # Minimum USB size in MB
USB_MAX_SIZE_GB=128   # Maximum USB size in GB for safety

# Convert size string to MB
convert_size_to_mb() {
    local size="$1"
    local number="${size%[A-Za-z]*}"
    local unit="${size##*[0-9]}"
    
    case "$unit" in
        "K"|"KB") echo "$(echo "scale=0; $number/1024" | bc)";;
        "M"|"MB") echo "$number";;
        "G"|"GB") echo "$(echo "scale=0; $number*1024" | bc)";;
        "T"|"TB") echo "$(echo "scale=0; $number*1024*1024" | bc)";;
        *) echo "0";;
    esac
}

# Detectar dispositivos USB removíveis com validações aprimoradas
detect_usb_devices() {
    local detected=false
    
    # Buscar dispositivos removíveis usando lsblk com informações detalhadas
    lsblk -d -n -o NAME,SIZE,TYPE,RM,MODEL,SERIAL,VENDOR 2>/dev/null | \
    while read -r name size type rm model serial vendor; do
        if [ "$type" = "disk" ]; then
            local device="/dev/$name"
            local size_mb=$(convert_size_to_mb "$size")
            
            # No WSL, ser mais permissivo com dispositivos que podem ser USBs
            local is_potential_usb=false
            
            if is_wsl; then
                # No WSL, considerar dispositivos que:
                # 1. São marcados como removíveis OU
                # 2. Têm tamanho típico de USB (1GB-1TB) e não são sda (disco do sistema)
                if [ "$rm" = "1" ] || ([ "$name" != "sda" ] && [ "$size_mb" -ge 1024 ] && [ "$size_mb" -le 1048576 ]); then
                    is_potential_usb=true
                fi
            else
                # No Linux nativo, usar apenas dispositivos marcados como removíveis
                if [ "$rm" = "1" ]; then
                    is_potential_usb=true
                fi
            fi
            
            if [ "$is_potential_usb" = true ]; then
                # Validar tamanho mínimo e máximo
                if [ "$size_mb" -lt "$USB_MIN_SIZE_MB" ]; then
                    log WARN "Dispositivo $device ignorado: muito pequeno (mínimo ${USB_MIN_SIZE_MB}MB)"
                    continue
                fi
                
                if [ "$size_mb" -gt "$((USB_MAX_SIZE_GB*1024))" ]; then
                    log WARN "Dispositivo $device ignorado: muito grande (máximo ${USB_MAX_SIZE_GB}GB)"
                    continue
                fi
                
                # Verificar se o dispositivo existe e não é um disco do sistema
                if [ -b "$device" ] && ! is_system_disk "$device"; then
                    detected=true
                    # Format: device:size:model:vendor:serial
                    echo "$device:$size:${model:-Unknown}:${vendor:-Unknown}:${serial:-Unknown}"
                fi
            fi
        fi
    done
    
    if [ "$detected" = false ]; then
        log ERROR "Nenhum dispositivo USB válido encontrado"
        return 1
    fi
}

# Verificar se é um disco do sistema com validações adicionais
is_system_disk() {
    local device="$1"
    
    # Verificar se o dispositivo é válido
    if [ ! -b "$device" ]; then
        log ERROR "Dispositivo inválido: $device não existe"
        return 0
    fi
    
    # Verificar se o dispositivo é um disco do sistema via fstab
    if grep -q "^$device" /etc/fstab; then
        log WARN "Dispositivo $device encontrado em /etc/fstab - possível disco do sistema"
        return 0
    fi
    
    # Verificar pontos de montagem críticos
    local critical_mounts=("/" "/boot" "/boot/efi" "/usr" "/var" "/home" "/opt" "[SWAP]")
    
    lsblk -n -o NAME,MOUNTPOINT "$device" 2>/dev/null | \
    while read -r name mountpoint; do
        for mount in "${critical_mounts[@]}"; do
            if [ "$mountpoint" = "$mount" ]; then
                log WARN "Dispositivo $device montado em ponto crítico: $mountpoint"
                return 0
            fi
        done
    done
    
    # Verificação especial para WSL - ser mais permissivo com USBs
    if is_wsl; then
        log DEBUG "WSL detectado - aplicando validações específicas para WSL"
        
        # No WSL, verificar se o dispositivo é realmente removível
        local device_name=$(basename "$device")
        local removable_status=$(cat "/sys/block/$device_name/removable" 2>/dev/null || echo "0")
        
        if [ "$removable_status" = "1" ]; then
            log DEBUG "Dispositivo $device é removível no WSL - permitindo uso"
            return 1  # Não é disco do sistema
        fi
        
        # Verificar se é um dispositivo de swap (comum no WSL)
        local fstype=$(lsblk -n -o FSTYPE "$device" 2>/dev/null || echo "")
        if [ "$fstype" = "swap" ]; then
            log DEBUG "Dispositivo $device é um dispositivo de swap - não é USB"
            return 0  # É dispositivo do sistema (swap)
        fi
        
        # Verificar tamanho - USBs típicos são menores que 1TB
        local size_bytes=$(lsblk -b -d -n -o SIZE "$device" 2>/dev/null || echo "0")
        local size_gb=$((size_bytes / 1024 / 1024 / 1024))
        
        # No WSL, ser mais permissivo com dispositivos que não são sda (disco principal)
        # Mas excluir dispositivos muito pequenos (menos de 1GB) que podem ser partições do sistema
        if [ "$device_name" != "sda" ] && [ "$size_gb" -ge 1 ] && [ "$size_gb" -lt 1024 ]; then
            log DEBUG "Dispositivo $device tem tamanho típico de USB (${size_gb}GB) e não é sda - permitindo uso"
            return 1  # Não é disco do sistema
        fi
        
        # Se o dispositivo é muito pequeno (menos de 1GB), aplicar validação mais rigorosa
        if [ "$size_gb" -lt 1 ]; then
            log WARN "Dispositivo $device muito pequeno (${size_gb}GB) - pode ser partição do sistema"
            echo -e "${YELLOW}AVISO: Este dispositivo é muito pequeno (${size_gb}GB).${NC}"
            echo -e "Pode ser uma partição do sistema ou dispositivo de entrada."
            echo -e "Continuar irá APAGAR TODOS OS DADOS no dispositivo."
            read -p "Tem certeza que deseja continuar? (y/N): " -n 1 -r
            echo
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                return 0  # É dispositivo do sistema
            fi
            return 1  # Usuário confirmou, permitir uso
        fi
        
        # Se chegou até aqui no WSL, aplicar validação mais rigorosa
        log WARN "Dispositivo $device no WSL pode ser um disco do sistema"
        echo -e "${YELLOW}AVISO: Este dispositivo pode ser um disco do sistema no WSL.${NC}"
        echo -e "Continuar irá APAGAR TODOS OS DADOS no dispositivo."
        read -p "Tem certeza que deseja continuar? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            return 0
        fi
        return 1  # Usuário confirmou, permitir uso
    fi
    
    # Verificação padrão para Linux nativo
    if lsblk -n -o FSTYPE "$device" | grep -qE "ext[234]|btrfs|xfs|zfs"; then
        log WARN "Dispositivo $device contém sistema de arquivos não-removível"
        echo -e "${YELLOW}AVISO: Este dispositivo parece conter um sistema de arquivos não-removível.${NC}"
        echo -e "Continuar irá APAGAR TODOS OS DADOS no dispositivo."
        read -p "Tem certeza que deseja continuar? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            return 0
        fi
    fi
    
    return 1
}

# Validar USB criado
validate_usb() {
    local device="$1"
    local success=true
    
    log INFO "Iniciando validação do USB em $device..."
    
    # 1. Verificar se o dispositivo ainda existe
    if [ ! -b "$device" ]; then
        log ERROR "Dispositivo $device não encontrado"
        return 1
    fi
    
    # 2. Tentar montar o dispositivo
    local mount_point="/tmp/syntropy_validate_$$"
    mkdir -p "$mount_point"
    
    if ! mount "$device"1 "$mount_point" 2>/dev/null; then
        log ERROR "Falha ao montar partição de boot"
        rm -rf "$mount_point"
        return 1
    fi
    
    # 3. Verificar arquivos essenciais
    local required_files=(
        "syntropy/config.json"
        "syntropy/keys/node.key"
        "syntropy/keys/node.pub"
        "boot/grub/grub.cfg"
    )
    
    for file in "${required_files[@]}"; do
        if [ ! -f "$mount_point/$file" ]; then
            log ERROR "Arquivo essencial não encontrado: $file"
            success=false
        else
            log DEBUG "Arquivo validado: $file"
        fi
    done
    
    # 4. Verificar permissões dos arquivos sensíveis
    if [ -d "$mount_point/syntropy/keys" ]; then
        local key_perms=$(stat -c %a "$mount_point/syntropy/keys/node.key")
        if [ "$key_perms" != "600" ]; then
            log ERROR "Permissões incorretas em node.key: $key_perms (deveria ser 600)"
            success=false
        fi
    fi
    
    # Limpar
    umount "$mount_point"
    rm -rf "$mount_point"
    
    if [ "$success" = true ]; then
        log SUCCESS "Validação do USB concluída com sucesso"
        return 0
    else
        log ERROR "Validação do USB falhou"
        return 1
    fi
}

# Função principal de seleção de dispositivo USB
select_usb_device() {
    log INFO "Detectando dispositivos USB..."
    
    # Capturar lista de dispositivos
    local temp_devices=$(detect_usb_devices)
    local devices=()
    
    # Converter para array
    if [ -n "$temp_devices" ]; then
        while IFS= read -r line; do
            devices+=("$line")
        done <<< "$temp_devices"
    fi
    
    # Caso 1: Nenhum dispositivo encontrado
    if [ ${#devices[@]} -eq 0 ]; then
        log WARN "Nenhum dispositivo USB removível encontrado automaticamente"
        echo ""
        
        # Verificar se estamos no WSL
        if is_wsl; then
            echo -e "${YELLOW}⚠ PROBLEMA DETECTADO NO WSL:${NC}"
            echo "O USB pode não estar sendo detectado pelo WSL."
            echo ""
            echo "Possíveis causas:"
            echo "• USB não está conectado"
            echo "• USB está sendo usado pelo Windows"
            echo "• WSL precisa ser reiniciado"
            echo "• Problema de permissão do WSL"
            echo ""
            echo "Soluções recomendadas:"
            echo "1. Verifique se o USB está conectado e aparece no Windows Explorer"
            echo "2. Feche todas as janelas do Windows Explorer"
            echo "3. Ejecte o USB com segurança no Windows"
            echo "4. Reconecte o USB"
            echo "5. Execute no PowerShell: wsl --shutdown"
            echo "6. Reinicie o WSL e teste novamente"
            echo ""
        fi
        
        echo "Diagnóstico do sistema:"
        echo ""
        
        echo "Todos os dispositivos de disco:"
        lsblk -d -o NAME,SIZE,TYPE,RM,MODEL 2>/dev/null | grep disk
        echo ""
        
        echo "Status de removibilidade por dispositivo:"
        for dev in /dev/sd[a-z]; do
            if [ -b "$dev" ]; then
                name=$(basename "$dev")
                removable=$(cat "/sys/block/$name/removable" 2>/dev/null || echo "N/A")
                size=$(lsblk -d -n -o SIZE "$dev" 2>/dev/null || echo "N/A")
                echo "  $dev: removível=$removable, tamanho=$size"
            fi
        done
        echo ""
        
        # Procurar por dispositivos que podem ser USBs não detectados
        echo -e "${BLUE}Procurando por possíveis USBs não detectados:${NC}"
        found_potential_usb=false
        for dev in /dev/sd[a-z]; do
            if [ -b "$dev" ]; then
                size_bytes=$(lsblk -b -d -n -o SIZE "$dev" 2>/dev/null || echo "0")
                size_gb=$((size_bytes / 1024 / 1024 / 1024))
                name=$(basename "$dev")
                removable=$(cat "/sys/block/$name/removable" 2>/dev/null || echo "0")
                
                # Procurar por dispositivos de tamanho típico de USB (8GB-1TB) que não são sda
                if [ "$name" != "sda" ] && [ "$size_gb" -ge 8 ] && [ "$size_gb" -le 1024 ]; then
                    echo -e "  ${GREEN}Possível USB: $dev (${size_gb}GB) - removível=$removable${NC}"
                    found_potential_usb=true
                fi
            fi
        done
        
        if [ "$found_potential_usb" = false ]; then
            echo -e "  ${RED}Nenhum dispositivo com tamanho típico de USB encontrado${NC}"
        fi
        echo ""
        
        read -p "Especificar um dispositivo manualmente? (y/N): " -r
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            echo ""
            echo "Dispositivos disponíveis:"
            lsblk -d -o NAME,SIZE,TYPE,MODEL 2>/dev/null | grep disk
            echo ""
            read -p "Digite o caminho do dispositivo (ex: /dev/sdb): " manual_device
            
            if [ -b "$manual_device" ]; then
                echo ""
                echo -e "${YELLOW}⚠ ATENÇÃO: Seleção manual de dispositivo${NC}"
                echo "Você está selecionando manualmente: $manual_device"
                echo "Certifique-se de que este é o dispositivo correto!"
                echo ""
                read -p "Confirmar seleção de $manual_device? (y/N): " -r
                if [[ $REPLY =~ ^[Yy]$ ]]; then
                    echo "$manual_device"
                    return 0
                else
                    log INFO "Seleção manual cancelada"
                    return 1
                fi
            else
                log ERROR "Dispositivo inválido: $manual_device"
                return 1
            fi
        else
            log ERROR "Nenhum dispositivo USB selecionado"
            return 1
        fi
    fi
    
    # Caso 2: Um dispositivo encontrado - SELEÇÃO AUTOMÁTICA
    if [ ${#devices[@]} -eq 1 ]; then
        local device_info="${devices[0]}"
        local device_path=$(echo "$device_info" | cut -d: -f1)
        local device_size=$(echo "$device_info" | cut -d: -f2)
        local device_model=$(echo "$device_info" | cut -d: -f3)
        
        log INFO "Dispositivo USB detectado automaticamente: $device_path"
        echo ""
        echo "Dispositivo selecionado automaticamente:"
        echo "  Caminho: $device_path"
        echo "  Tamanho: $device_size"
        echo "  Modelo: $device_model"
        
        show_device_details "$device_path"
        
        # Seleção automática - apenas retorna o dispositivo
        echo "$device_path"
        return 0
    fi
    
    # Caso 3: Múltiplos dispositivos encontrados
    log INFO "Múltiplos dispositivos USB encontrados"
    echo ""
    echo -e "${YELLOW}⚠ ATENÇÃO: Múltiplos dispositivos USB detectados!${NC}"
    echo -e "${YELLOW}   Certifique-se de selecionar o dispositivo correto.${NC}"
    echo -e "${YELLOW}   A formatação apagará TODOS os dados do dispositivo selecionado.${NC}"
    echo ""
    
    echo -e "${CYAN}Dispositivos USB disponíveis:${NC}"
    
    local i=1
    for device_info in "${devices[@]}"; do
        local device_path=$(echo "$device_info" | cut -d: -f1)
        local device_size=$(echo "$device_info" | cut -d: -f2)
        local device_model=$(echo "$device_info" | cut -d: -f3)
        local device_vendor=$(echo "$device_info" | cut -d: -f4)
        
        # Destacar dispositivos que podem ser USBs de teclado/mouse
        local is_likely_input_device=false
        if [[ "$device_model" =~ -i[0-9]+$ ]] || [[ "$device_model" =~ [Kk]eyboard ]] || [[ "$device_model" =~ [Mm]ouse ]]; then
            is_likely_input_device=true
        fi
        
        if [ "$is_likely_input_device" = true ]; then
            echo -e "  ${YELLOW}$i) $device_path ($device_size) - $device_model (${device_vendor})${NC}"
            echo -e "     ${YELLOW}⚠ Possível dispositivo de entrada (teclado/mouse)${NC}"
        else
            echo -e "  ${GREEN}$i) $device_path ($device_size) - $device_model (${device_vendor})${NC}"
        fi
        ((i++))
    done
    
    echo ""
    echo -e "${BLUE}Dicas para seleção segura:${NC}"
    echo "• Verifique o tamanho do dispositivo (USBs de armazenamento são geralmente 8GB+)"
    echo "• Dispositivos muito pequenos (alguns MB) podem ser de teclado/mouse"
    echo "• Verifique o modelo e fabricante"
    echo "• Se em dúvida, desconecte outros USBs e execute novamente"
    echo ""
    
    while true; do
        read -p "Selecione o dispositivo (1-${#devices[@]}) ou 'q' para sair: " choice
        
        if [ "$choice" = "q" ] || [ "$choice" = "Q" ]; then
            log INFO "Seleção cancelada pelo usuário"
            return 1
        fi
        
        if [[ "$choice" =~ ^[0-9]+$ ]] && [ "$choice" -ge 1 ] && [ "$choice" -le ${#devices[@]} ]; then
            local selected_info="${devices[$((choice-1))]}"
            local selected_device=$(echo "$selected_info" | cut -d: -f1)
            local selected_size=$(echo "$selected_info" | cut -d: -f2)
            local selected_model=$(echo "$selected_info" | cut -d: -f3)
            
            echo ""
            echo -e "${CYAN}Dispositivo selecionado: $selected_device${NC}"
            echo "  Tamanho: $selected_size"
            echo "  Modelo: $selected_model"
            echo ""
            
            show_device_details "$selected_device"
            echo ""
            
            # Verificação adicional para dispositivos pequenos
            local size_bytes=$(lsblk -b -d -n -o SIZE "$selected_device" 2>/dev/null || echo "0")
            local size_gb=$((size_bytes / 1024 / 1024 / 1024))
            
            if [ "$size_gb" -lt 1 ]; then
                echo -e "${RED}⚠ ATENÇÃO: Este dispositivo é muito pequeno (${size_gb}GB)!${NC}"
                echo -e "${RED}   Pode ser um dispositivo de entrada (teclado/mouse).${NC}"
                echo -e "${RED}   Tem certeza que deseja continuar?${NC}"
                echo ""
                read -p "Continuar mesmo assim? (y/N): " -r
                if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                    echo "Seleção cancelada. Tente novamente."
                    echo ""
                    continue
                fi
            fi
            
            read -p "Confirmar seleção de $selected_device? (y/N): " -r
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                echo "$selected_device"
                return 0
            fi
            echo ""
        else
            echo "Escolha inválida. Digite um número entre 1-${#devices[@]} ou 'q'"
        fi
    done
}

# Mostrar detalhes do dispositivo com informações completas
show_device_details() {
    local device="$1"
    
    echo ""
    echo -e "${CYAN}═══════════════════════════════════════════════════════════════${NC}"
    echo -e "${CYAN}                    DETALHES DO DISPOSITIVO                    ${NC}"
    echo -e "${CYAN}═══════════════════════════════════════════════════════════════${NC}"
    
    # Obter informações básicas
    local size=$(lsblk -d -n -o SIZE "$device" 2>/dev/null || echo "Desconhecido")
    local model=$(lsblk -d -n -o MODEL "$device" 2>/dev/null || echo "Desconhecido")
    local vendor=$(lsblk -d -n -o VENDOR "$device" 2>/dev/null || echo "Desconhecido")
    local serial=$(lsblk -d -n -o SERIAL "$device" 2>/dev/null || echo "Desconhecido")
    
    # Calcular tamanho em GB
    local size_bytes=$(lsblk -b -d -n -o SIZE "$device" 2>/dev/null || echo "0")
    local size_gb=$((size_bytes / 1024 / 1024 / 1024))
    
    # Obter informações de removibilidade
    local device_name=$(basename "$device")
    local removable_status=$(cat "/sys/block/$device_name/removable" 2>/dev/null || echo "0")
    local removable_text="Não"
    if [ "$removable_status" = "1" ]; then
        removable_text="Sim"
    fi
    
    echo -e "${BLUE}Informações Básicas:${NC}"
    echo "  Caminho do Dispositivo: $device"
    echo "  Tamanho Total: $size (${size_gb}GB)"
    echo "  Modelo: $model"
    echo "  Fabricante: $vendor"
    echo "  Número de Série: $serial"
    echo "  Dispositivo Removível: $removable_text"
    
    # Verificar se é WSL e mostrar informações adicionais
    if is_wsl; then
        echo -e "${BLUE}Informações WSL:${NC}"
        echo "  Ambiente: Windows Subsystem for Linux"
        echo "  Acesso via: WSL → Windows → Hardware USB"
    fi
    
    # Avisos sobre tamanho
    echo -e "${BLUE}Validações de Tamanho:${NC}"
    if [ "$size_gb" -lt 8 ]; then
        echo "  ${YELLOW}⚠ AVISO: Dispositivo menor que 8GB (requisito mínimo)${NC}"
    else
        echo "  ${GREEN}✓ Tamanho adequado para instalação${NC}"
    fi
    
    if [ "$size_gb" -gt 1024 ]; then
        echo "  ${YELLOW}⚠ AVISO: Dispositivo muito grande (${size_gb}GB) - pode ser armazenamento interno${NC}"
    fi
    
    # Verificar partições existentes
    echo -e "${BLUE}Partições Existentes:${NC}"
    local partitions=$(lsblk -n -o NAME,SIZE,FSTYPE,MOUNTPOINT "$device" 2>/dev/null | grep -v "^$(basename "$device") ")
    if [ -n "$partitions" ]; then
        echo "$partitions" | sed 's/^/  /'
        echo ""
        echo -e "${RED}  ⚠ ATENÇÃO: DADOS EXISTENTES SERÃO PERMANENTEMENTE APAGADOS!${NC}"
    else
        echo "  Nenhuma (dispositivo limpo)"
    fi
    
    echo -e "${CYAN}═══════════════════════════════════════════════════════════════${NC}"
}

# Validação abrangente de segurança do USB
validate_usb_safety() {
    local device="$1"
    
    log INFO "Executando validação de segurança para $device..."
    
    # Verificação básica de existência
    if [ ! -b "$device" ]; then
        log ERROR "$device não é um dispositivo de bloco válido"
        return 1
    fi
    
    # Verificação de disco do sistema
    if is_system_disk "$device"; then
        log ERROR "ERRO: Dispositivo parece ser um disco do sistema!"
        echo ""
        echo "Este dispositivo contém arquivos críticos do sistema."
        echo "Usar este dispositivo danificaria seu sistema operacional."
        echo ""
        echo "Por favor, use um dispositivo USB dedicado para armazenamento."
        return 1
    fi
    
    # Verificação de tamanho
    local size_bytes=$(lsblk -b -d -n -o SIZE "$device" 2>/dev/null || echo "0")
    local size_gb=$((size_bytes / 1024 / 1024 / 1024))
    
    if [ "$size_gb" -lt 8 ]; then
        log WARN "Dispositivo muito pequeno (${size_gb}GB) - mínimo recomendado: 8GB"
        read -p "Continuar mesmo assim? (y/N): " -r
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log INFO "Operação cancelada devido ao tamanho insuficiente"
            return 1
        fi
    fi
    
    # Ajustar limite superior baseado no ambiente
    local max_size_gb=512
    if is_wsl; then
        max_size_gb=1024  # No WSL, ser mais permissivo
        log DEBUG "WSL detectado - usando limite superior de ${max_size_gb}GB"
    fi
    
    if [ "$size_gb" -gt "$max_size_gb" ]; then
        log WARN "Dispositivo muito grande (${size_gb}GB) - pode ser armazenamento interno"
        echo "Dispositivos USB típicos são menores que ${max_size_gb}GB."
        echo ""
        read -p "Tem certeza de que este é um dispositivo USB? (y/N): " -r
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log INFO "Validação cancelada por segurança"
            return 1
        fi
    fi
    
    # Verificação de dados existentes
    check_existing_data "$device"
    
    log SUCCESS "Validação de segurança concluída"
    return 0
}

# Verificar dados existentes
check_existing_data() {
    local device="$1"
    
    # Verificar partições
    local partitions=$(lsblk -ln -o NAME,FSTYPE "$device" 2>/dev/null | grep -v "^$(basename "$device")$")
    
    if [ -n "$partitions" ]; then
        echo ""
        echo "${YELLOW}AVISO: Dispositivo contém partições existentes:${NC}"
        echo "$partitions" | sed 's/^/  /'
        echo ""
        echo "${RED}TODOS OS DADOS NESTE DISPOSITIVO SERÃO PERMANENTEMENTE PERDIDOS!${NC}"
        echo ""
        read -p "Continuar e apagar todos os dados? (y/N): " -r
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log INFO "Operação cancelada para preservar dados"
            exit 0
        fi
    fi
}

# Confirmação final crítica antes da formatação
final_device_confirmation() {
    local device="$1"
    local node_name="$2"
    
    echo ""
    echo -e "${RED}╔════════════════════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${RED}║                           ⚠ CONFIRMAÇÃO CRÍTICA ⚠                            ║${NC}"
    echo -e "${RED}║                                                                                ║${NC}"
    echo -e "${RED}║  ESTE É O ÚLTIMO MOMENTO PARA CANCELAR A OPERAÇÃO!                            ║${NC}"
    echo -e "${RED}║                                                                                ║${NC}"
    echo -e "${RED}║  O DISPOSITIVO SELECIONADO SERÁ COMPLETAMENTE FORMATADO E                      ║${NC}"
    echo -e "${RED}║  TODOS OS DADOS EXISTENTES SERÃO PERMANENTEMENTE PERDIDOS!                     ║${NC}"
    echo -e "${RED}║                                                                                ║${NC}"
    echo -e "${RED}║  NÃO É POSSÍVEL REVERTER ESTA OPERAÇÃO!                                       ║${NC}"
    echo -e "${RED}╚════════════════════════════════════════════════════════════════════════════════╝${NC}"
    echo ""
    
    # Mostrar informações do dispositivo selecionado
    echo -e "${CYAN}DISPOSITIVO SELECIONADO:${NC}"
    echo "  Caminho: $device"
    
    # Obter informações detalhadas
    local size=$(lsblk -d -n -o SIZE "$device" 2>/dev/null || echo "Desconhecido")
    local model=$(lsblk -d -n -o MODEL "$device" 2>/dev/null || echo "Desconhecido")
    local vendor=$(lsblk -d -n -o VENDOR "$device" 2>/dev/null || echo "Desconhecido")
    
    echo "  Tamanho: $size"
    echo "  Modelo: $model"
    echo "  Fabricante: $vendor"
    echo ""
    
    # Mostrar informações do nó
    echo -e "${CYAN}CONFIGURAÇÃO DO NÓ:${NC}"
    echo "  Nome do Nó: $node_name"
    echo "  Sistema: Ubuntu Server + Syntropy Cooperative Grid"
    echo ""
    
    # Verificar se há outros dispositivos similares
    echo -e "${YELLOW}VERIFICAÇÃO DE SEGURANÇA:${NC}"
    echo "Dispositivos similares no sistema:"
    lsblk -d -o NAME,SIZE,MODEL 2>/dev/null | grep -E "disk|NAME" | head -10
    echo ""
    
    # Confirmação tripla
    echo -e "${RED}CONFIRMAÇÃO TRIPLA REQUERIDA:${NC}"
    echo ""
    
    # Primeira confirmação
    echo -e "${YELLOW}1. Confirme que este é o dispositivo correto:${NC}"
    read -p "   Digite o caminho completo do dispositivo ($device): " confirm_path
    
    if [ "$confirm_path" != "$device" ]; then
        log ERROR "Caminho do dispositivo não confere!"
        log ERROR "Esperado: $device"
        log ERROR "Digitado: $confirm_path"
        log ERROR "Operação cancelada por segurança."
        exit 1
    fi
    
    # Segunda confirmação
    echo -e "${YELLOW}2. Confirme que você tem certeza absoluta:${NC}"
    read -p "   Digite 'SIM, TENHO CERTEZA' (exatamente): " confirm_certainty
    
    if [ "$confirm_certainty" != "SIM, TENHO CERTEZA" ]; then
        log ERROR "Confirmação não reconhecida."
        log ERROR "Operação cancelada por segurança."
        exit 1
    fi
    
    # Terceira confirmação
    echo -e "${YELLOW}3. Confirmação final:${NC}"
    read -p "   Digite 'FORMATAR AGORA' para prosseguir: " confirm_format
    
    if [ "$confirm_format" != "FORMATAR AGORA" ]; then
        log ERROR "Confirmação final não reconhecida."
        log ERROR "Operação cancelada por segurança."
        exit 1
    fi
    
    echo ""
    echo -e "${GREEN}✓ Confirmação tripla aprovada!${NC}"
    echo -e "${GREEN}✓ Iniciando formatação do dispositivo $device${NC}"
    echo -e "${GREEN}✓ Criando nó Syntropy: $node_name${NC}"
    echo ""
    
    # Pequena pausa para o usuário ver a confirmação
    sleep 2
}

# Resumo dos dispositivos USB
show_usb_summary() {
    log INFO "Resumo da Detecção de Dispositivos USB:"
    
    local device_list=$(detect_usb_devices)
    local count=0
    
    if [ -n "$device_list" ]; then
        count=$(echo "$device_list" | wc -l)
        echo "  Dispositivos USB encontrados: $count"
        echo "  Dispositivos disponíveis:"
        while IFS=: read -r device size model; do
            echo "    $device - $size - $model"
        done <<< "$device_list"
    else
        echo "  Dispositivos USB encontrados: 0"
    fi
}