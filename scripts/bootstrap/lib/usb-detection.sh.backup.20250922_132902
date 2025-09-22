#!/bin/bash

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
        if [ "$type" = "disk" ] && [ "$rm" = "1" ]; then
            local device="/dev/$name"
            local size_mb=$(convert_size_to_mb "$size")
            
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
    
    # Verificar se tem partições com dados importantes
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
        log WARN "Nenhum dispositivo USB removível encontrado"
        echo ""
        echo "Diagnóstico do sistema:"
        echo ""
        
        echo "Todos os dispositivos de disco:"
        lsblk -d -o NAME,SIZE,TYPE,RM,MODEL 2>/dev/null | grep disk
        echo ""
        
        echo "Status de removibilidade por dispositivo:"
        for dev in /dev/sd[a-z]; do
            if [ -b "$dev" ]; then
                local name=$(basename "$dev")
                local removable=$(cat "/sys/block/$name/removable" 2>/dev/null || echo "N/A")
                local size=$(lsblk -d -n -o SIZE "$dev" 2>/dev/null || echo "N/A")
                echo "  $dev: removível=$removable, tamanho=$size"
            fi
        done
        echo ""
        
        read -p "Especificar um dispositivo manualmente? (y/N): " -r
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            echo ""
            echo "Dispositivos disponíveis:"
            lsblk -d -o NAME,SIZE,TYPE 2>/dev/null | grep disk
            echo ""
            read -p "Digite o caminho do dispositivo (ex: /dev/sdb): " manual_device
            
            if [ -b "$manual_device" ]; then
                echo "$manual_device"
                return 0
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
    echo "Dispositivos USB disponíveis:"
    
    local i=1
    for device_info in "${devices[@]}"; do
        local device_path=$(echo "$device_info" | cut -d: -f1)
        local device_size=$(echo "$device_info" | cut -d: -f2)
        local device_model=$(echo "$device_info" | cut -d: -f3)
        echo "  $i) $device_path ($device_size) - $device_model"
        ((i++))
    done
    
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
            
            show_device_details "$selected_device"
            echo ""
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

# Mostrar detalhes do dispositivo
show_device_details() {
    local device="$1"
    
    echo ""
    echo "Detalhes do dispositivo:"
    
    # Obter informações básicas
    local size=$(lsblk -d -n -o SIZE "$device" 2>/dev/null || echo "Desconhecido")
    local model=$(lsblk -d -n -o MODEL "$device" 2>/dev/null || echo "Desconhecido")
    local vendor=$(lsblk -d -n -o VENDOR "$device" 2>/dev/null || echo "Desconhecido")
    
    # Calcular tamanho em GB
    local size_bytes=$(lsblk -b -d -n -o SIZE "$device" 2>/dev/null || echo "0")
    local size_gb=$((size_bytes / 1024 / 1024 / 1024))
    
    echo "  Caminho: $device"
    echo "  Tamanho Total: $size (${size_gb}GB)"
    echo "  Modelo: $model"
    echo "  Fabricante: $vendor"
    
    # Avisos sobre tamanho
    if [ "$size_gb" -lt 8 ]; then
        echo "  ${YELLOW}AVISO: Dispositivo menor que 8GB (requisito mínimo)${NC}"
    fi
    
    if [ "$size_gb" -gt 512 ]; then
        echo "  ${YELLOW}AVISO: Dispositivo muito grande (${size_gb}GB) - pode ser armazenamento interno${NC}"
    fi
    
    # Verificar partições existentes
    echo "  Partições atuais:"
    local partitions=$(lsblk -n -o NAME,SIZE,FSTYPE,MOUNTPOINT "$device" 2>/dev/null | grep -v "^$(basename "$device") ")
    if [ -n "$partitions" ]; then
        echo "$partitions" | sed 's/^/    /'
        echo "  ${YELLOW}AVISO: Dados existentes serão permanentemente apagados!${NC}"
    else
        echo "    Nenhuma (dispositivo limpo)"
    fi
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
    
    if [ "$size_gb" -gt 512 ]; then
        log WARN "Dispositivo muito grande (${size_gb}GB) - pode ser armazenamento interno"
        echo "Dispositivos USB típicos são menores que 512GB."
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