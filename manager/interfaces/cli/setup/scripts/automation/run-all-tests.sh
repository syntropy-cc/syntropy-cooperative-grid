#!/bin/bash

# Script para executar todos os testes do componente setup
# Autor: Sistema de Diagnóstico Automático
# Data: $(date)

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configurações
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SETUP_DIR="$(dirname "$(dirname "$SCRIPT_DIR")")"
TESTS_DIR="$SETUP_DIR/tests"
REPORT_DIR="$SETUP_DIR/reports"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# Criar diretório de relatórios
mkdir -p "$REPORT_DIR"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  DIAGNÓSTICO COMPLETO - SETUP COMPONENT${NC}"
echo -e "${BLUE}========================================${NC}"
echo -e "Data/Hora: $(date)"
echo -e "Diretório: $SETUP_DIR"
echo -e "Relatórios: $REPORT_DIR"
echo ""

# Função para log
log() {
    echo -e "${BLUE}[$(date +'%H:%M:%S')]${NC} $1"
}

error() {
    echo -e "${RED}[ERRO]${NC} $1" >&2
}

success() {
    echo -e "${GREEN}[SUCESSO]${NC} $1"
}

warning() {
    echo -e "${YELLOW}[AVISO]${NC} $1"
}

# Função para verificar se Go está instalado
check_go() {
    log "Verificando instalação do Go..."
    if ! command -v go &> /dev/null; then
        error "Go não está instalado ou não está no PATH"
        return 1
    fi
    
    GO_VERSION=$(go version | cut -d' ' -f3)
    success "Go encontrado: $GO_VERSION"
    return 0
}

# Função para verificar estrutura de arquivos
check_structure() {
    log "Verificando estrutura de arquivos..."
    
    local missing_files=()
    
    # Verificar arquivos principais
    local required_files=(
        "$SETUP_DIR/src/setup.go"
        "$SETUP_DIR/src/validator.go"
        "$SETUP_DIR/src/configurator.go"
        "$SETUP_DIR/src/logger.go"
        "$TESTS_DIR/unit/setup_test.go"
        "$TESTS_DIR/integration/setup_integration_test.go"
        "$TESTS_DIR/e2e/setup_e2e_test.go"
        "$TESTS_DIR/security/security_test.go"
        "$TESTS_DIR/performance/load_test.go"
    )
    
    for file in "${required_files[@]}"; do
        if [[ ! -f "$file" ]]; then
            missing_files+=("$file")
        fi
    done
    
    if [[ ${#missing_files[@]} -eq 0 ]]; then
        success "Todos os arquivos principais encontrados"
    else
        warning "Arquivos ausentes:"
        for file in "${missing_files[@]}"; do
            echo "  - $file"
        done
    fi
    
    # Verificar diretórios de teste
    local test_dirs=("unit" "integration" "e2e" "security" "performance" "mocks" "fixtures" "helpers")
    local missing_dirs=()
    
    for dir in "${test_dirs[@]}"; do
        if [[ ! -d "$TESTS_DIR/$dir" ]]; then
            missing_dirs+=("$dir")
        fi
    done
    
    if [[ ${#missing_dirs[@]} -eq 0 ]]; then
        success "Todos os diretórios de teste encontrados"
    else
        warning "Diretórios de teste ausentes:"
        for dir in "${missing_dirs[@]}"; do
            echo "  - $TESTS_DIR/$dir"
        done
    fi
}

# Função para executar testes unitários
run_unit_tests() {
    log "Executando testes unitários..."
    
    local unit_report="$REPORT_DIR/unit_tests_$TIMESTAMP.txt"
    local unit_coverage="$REPORT_DIR/unit_coverage_$TIMESTAMP.out"
    
    cd "$SETUP_DIR"
    
    # Tentar executar testes unitários
    if go test -v -coverprofile="$unit_coverage" ./tests/unit/... > "$unit_report" 2>&1; then
        success "Testes unitários executados com sucesso"
        echo "  Relatório: $unit_report"
        echo "  Cobertura: $unit_coverage"
    else
        error "Falha na execução dos testes unitários"
        echo "  Verifique o relatório: $unit_report"
        return 1
    fi
}

# Função para executar testes de integração
run_integration_tests() {
    log "Executando testes de integração..."
    
    local integration_report="$REPORT_DIR/integration_tests_$TIMESTAMP.txt"
    local integration_coverage="$REPORT_DIR/integration_coverage_$TIMESTAMP.out"
    
    cd "$SETUP_DIR"
    
    if go test -v -coverprofile="$integration_coverage" ./tests/integration/... > "$integration_report" 2>&1; then
        success "Testes de integração executados com sucesso"
        echo "  Relatório: $integration_report"
        echo "  Cobertura: $integration_coverage"
    else
        error "Falha na execução dos testes de integração"
        echo "  Verifique o relatório: $integration_report"
        return 1
    fi
}

# Função para executar testes E2E
run_e2e_tests() {
    log "Executando testes end-to-end..."
    
    local e2e_report="$REPORT_DIR/e2e_tests_$TIMESTAMP.txt"
    
    cd "$SETUP_DIR"
    
    if go test -v ./tests/e2e/... > "$e2e_report" 2>&1; then
        success "Testes E2E executados com sucesso"
        echo "  Relatório: $e2e_report"
    else
        error "Falha na execução dos testes E2E"
        echo "  Verifique o relatório: $e2e_report"
        return 1
    fi
}

# Função para executar testes de segurança
run_security_tests() {
    log "Executando testes de segurança..."
    
    local security_report="$REPORT_DIR/security_tests_$TIMESTAMP.txt"
    
    cd "$SETUP_DIR"
    
    if go test -v ./tests/security/... > "$security_report" 2>&1; then
        success "Testes de segurança executados com sucesso"
        echo "  Relatório: $security_report"
    else
        error "Falha na execução dos testes de segurança"
        echo "  Verifique o relatório: $security_report"
        return 1
    fi
}

# Função para executar testes de performance
run_performance_tests() {
    log "Executando testes de performance..."
    
    local performance_report="$REPORT_DIR/performance_tests_$TIMESTAMP.txt"
    
    cd "$SETUP_DIR"
    
    if go test -v -bench=. ./tests/performance/... > "$performance_report" 2>&1; then
        success "Testes de performance executados com sucesso"
        echo "  Relatório: $performance_report"
    else
        error "Falha na execução dos testes de performance"
        echo "  Verifique o relatório: $performance_report"
        return 1
    fi
}

# Função para gerar relatório de cobertura
generate_coverage_report() {
    log "Gerando relatório de cobertura..."
    
    local coverage_report="$REPORT_DIR/coverage_report_$TIMESTAMP.html"
    local coverage_summary="$REPORT_DIR/coverage_summary_$TIMESTAMP.txt"
    
    cd "$SETUP_DIR"
    
    # Gerar cobertura completa
    if go test -coverprofile="$REPORT_DIR/coverage_$TIMESTAMP.out" ./... > /dev/null 2>&1; then
        # Gerar relatório HTML
        go tool cover -html="$REPORT_DIR/coverage_$TIMESTAMP.out" -o "$coverage_report"
        
        # Gerar resumo
        go tool cover -func="$REPORT_DIR/coverage_$TIMESTAMP.out" > "$coverage_summary"
        
        success "Relatório de cobertura gerado"
        echo "  HTML: $coverage_report"
        echo "  Resumo: $coverage_summary"
    else
        warning "Não foi possível gerar relatório de cobertura"
    fi
}

# Função para executar análise estática
run_static_analysis() {
    log "Executando análise estática..."
    
    local static_report="$REPORT_DIR/static_analysis_$TIMESTAMP.txt"
    
    cd "$SETUP_DIR"
    
    {
        echo "=== ANÁLISE ESTÁTICA - SETUP COMPONENT ==="
        echo "Data: $(date)"
        echo ""
        
        echo "=== GO VET ==="
        go vet ./... 2>&1 || echo "go vet encontrou problemas"
        echo ""
        
        echo "=== GO FMT ==="
        go fmt ./... 2>&1 || echo "go fmt encontrou problemas"
        echo ""
        
        echo "=== GO MOD TIDY ==="
        go mod tidy 2>&1 || echo "go mod tidy encontrou problemas"
        echo ""
        
        echo "=== VERIFICAÇÃO DE DEPENDÊNCIAS ==="
        go list -m all 2>&1 || echo "Erro ao listar dependências"
        
    } > "$static_report"
    
    success "Análise estática concluída"
    echo "  Relatório: $static_report"
}

# Função para gerar relatório final
generate_final_report() {
    log "Gerando relatório final..."
    
    local final_report="$REPORT_DIR/diagnostic_report_$TIMESTAMP.md"
    
    cat > "$final_report" << EOF
# Relatório de Diagnóstico - Setup Component

**Data/Hora:** $(date)
**Diretório:** $SETUP_DIR
**Timestamp:** $TIMESTAMP

## Resumo Executivo

Este relatório contém os resultados do diagnóstico completo do componente setup.

## Estrutura de Arquivos

### Arquivos Principais
EOF

    # Verificar arquivos principais
    local main_files=(
        "src/setup.go"
        "src/validator.go" 
        "src/configurator.go"
        "src/logger.go"
    )
    
    for file in "${main_files[@]}"; do
        if [[ -f "$SETUP_DIR/$file" ]]; then
            echo "- ✅ $file" >> "$final_report"
        else
            echo "- ❌ $file (AUSENTE)" >> "$final_report"
        fi
    done
    
    cat >> "$final_report" << EOF

### Diretórios de Teste
EOF

    local test_dirs=("unit" "integration" "e2e" "security" "performance" "mocks" "fixtures" "helpers")
    for dir in "${test_dirs[@]}"; do
        if [[ -d "$TESTS_DIR/$dir" ]]; then
            echo "- ✅ $dir/" >> "$final_report"
        else
            echo "- ❌ $dir/ (AUSENTE)" >> "$final_report"
        fi
    done
    
    cat >> "$final_report" << EOF

## Relatórios Gerados

- **Testes Unitários:** unit_tests_$TIMESTAMP.txt
- **Testes de Integração:** integration_tests_$TIMESTAMP.txt
- **Testes E2E:** e2e_tests_$TIMESTAMP.txt
- **Testes de Segurança:** security_tests_$TIMESTAMP.txt
- **Testes de Performance:** performance_tests_$TIMESTAMP.txt
- **Análise Estática:** static_analysis_$TIMESTAMP.txt
- **Cobertura de Código:** coverage_report_$TIMESTAMP.html

## Próximos Passos

1. Revisar relatórios individuais
2. Corrigir problemas identificados
3. Executar novamente o diagnóstico
4. Implementar melhorias sugeridas

---
*Relatório gerado automaticamente pelo sistema de diagnóstico*
EOF

    success "Relatório final gerado: $final_report"
}

# Função principal
main() {
    echo -e "${BLUE}Iniciando diagnóstico completo...${NC}"
    echo ""
    
    # Verificações iniciais
    if ! check_go; then
        error "Go não está disponível. Abortando."
        exit 1
    fi
    
    check_structure
    echo ""
    
    # Executar testes
    local test_results=()
    
    run_unit_tests && test_results+=("unit:SUCCESS") || test_results+=("unit:FAILED")
    echo ""
    
    run_integration_tests && test_results+=("integration:SUCCESS") || test_results+=("integration:FAILED")
    echo ""
    
    run_e2e_tests && test_results+=("e2e:SUCCESS") || test_results+=("e2e:FAILED")
    echo ""
    
    run_security_tests && test_results+=("security:SUCCESS") || test_results+=("security:FAILED")
    echo ""
    
    run_performance_tests && test_results+=("performance:SUCCESS") || test_results+=("performance:FAILED")
    echo ""
    
    # Análises adicionais
    run_static_analysis
    echo ""
    
    generate_coverage_report
    echo ""
    
    generate_final_report
    echo ""
    
    # Resumo final
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}  RESUMO DO DIAGNÓSTICO${NC}"
    echo -e "${BLUE}========================================${NC}"
    
    for result in "${test_results[@]}"; do
        local test_type=$(echo "$result" | cut -d':' -f1)
        local status=$(echo "$result" | cut -d':' -f2)
        
        if [[ "$status" == "SUCCESS" ]]; then
            echo -e "✅ $test_type: ${GREEN}SUCESSO${NC}"
        else
            echo -e "❌ $test_type: ${RED}FALHOU${NC}"
        fi
    done
    
    echo ""
    echo -e "Relatórios salvos em: ${BLUE}$REPORT_DIR${NC}"
    echo -e "Relatório principal: ${BLUE}diagnostic_report_$TIMESTAMP.md${NC}"
    echo ""
    
    # Contar sucessos
    local success_count=$(printf '%s\n' "${test_results[@]}" | grep -c "SUCCESS" || true)
    local total_count=${#test_results[@]}
    
    if [[ $success_count -eq $total_count ]]; then
        success "Todos os testes passaram! ($success_count/$total_count)"
        exit 0
    else
        warning "Alguns testes falharam ($success_count/$total_count)"
        exit 1
    fi
}

# Executar função principal
main "$@"
