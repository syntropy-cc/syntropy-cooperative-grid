#!/bin/bash

# Script para executar todos os exemplos de teste
# Autor: Sistema de Exemplos

set -e

# Cores
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}=== EXECUTANDO TODOS OS EXEMPLOS ===${NC}"

# Lista de exemplos
examples=(
    "example1_basic_setup.go"
    "example2_validation.go"
    "example3_configuration.go"
    "example4_performance.go"
    "example5_security.go"
    "example6_integration.go"
)

# Executar cada exemplo
for example in "${examples[@]}"; do
    if [[ -f "$example" ]]; then
        echo -e "${BLUE}Executando $example...${NC}"
        if go run "$example"; then
            echo -e "${GREEN}✅ $example executado com sucesso${NC}"
        else
            echo -e "${RED}❌ $example falhou${NC}"
        fi
        echo ""
    else
        echo -e "${RED}❌ $example não encontrado${NC}"
    fi
done

echo -e "${BLUE}=== EXECUÇÃO CONCLUÍDA ===${NC}"
