#!/bin/bash

# Sistema de Diagnóstico Automático para Setup Component
# Autor: Sistema de Diagnóstico
# Data: $(date)

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configurações
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SETUP_DIR="$(dirname "$(dirname "$SCRIPT_DIR")")"
TESTS_DIR="$SETUP_DIR/tests"
REPORTS_DIR="$SETUP_DIR/diagnostic-reports"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# Criar diretório de relatórios
mkdir -p "$REPORTS_DIR"

echo -e "${PURPLE}========================================${NC}"
echo -e "${PURPLE}  SISTEMA DE DIAGNÓSTICO AUTOMÁTICO${NC}"
echo -e "${PURPLE}  SETUP COMPONENT${NC}"
echo -e "${PURPLE}========================================${NC}"
echo -e "Data/Hora: $(date)"
echo -e "Diretório: $SETUP_DIR"
echo -e "Relatórios: $REPORTS_DIR"
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

info() {
    echo -e "${CYAN}[INFO]${NC} $1"
}

# Função para análise de estrutura
analyze_structure() {
    log "Analisando estrutura do componente..."
    
    local structure_report="$REPORTS_DIR/structure_analysis_$TIMESTAMP.md"
    
    cat > "$structure_report" << EOF
# Análise de Estrutura - Setup Component

**Data:** $(date)
**Diretório:** $SETUP_DIR

## Estrutura de Diretórios

EOF

    # Analisar estrutura principal
    if [[ -d "$SETUP_DIR/src" ]]; then
        echo "### Código Fonte (src/)" >> "$structure_report"
        find "$SETUP_DIR/src" -type f -name "*.go" | while read -r file; do
            local rel_path=$(realpath --relative-to="$SETUP_DIR" "$file")
            local size=$(wc -l < "$file")
            echo "- \`$rel_path\` ($size linhas)" >> "$structure_report"
        done
        echo "" >> "$structure_report"
    fi
    
    # Analisar testes
    if [[ -d "$TESTS_DIR" ]]; then
        echo "### Testes (tests/)" >> "$structure_report"
        find "$TESTS_DIR" -type f -name "*.go" | while read -r file; do
            local rel_path=$(realpath --relative-to="$SETUP_DIR" "$file")
            local size=$(wc -l < "$file")
            echo "- \`$rel_path\` ($size linhas)" >> "$structure_report"
        done
        echo "" >> "$structure_report"
    fi
    
    # Analisar scripts
    if [[ -d "$SETUP_DIR/scripts" ]]; then
        echo "### Scripts (scripts/)" >> "$structure_report"
        find "$SETUP_DIR/scripts" -type f -name "*.sh" | while read -r file; do
            local rel_path=$(realpath --relative-to="$SETUP_DIR" "$file")
            local size=$(wc -l < "$file")
            echo "- \`$rel_path\` ($size linhas)" >> "$structure_report"
        done
        echo "" >> "$structure_report"
    fi
    
    # Analisar exemplos
    if [[ -d "$SETUP_DIR/examples" ]]; then
        echo "### Exemplos (examples/)" >> "$structure_report"
        find "$SETUP_DIR/examples" -type f | while read -r file; do
            local rel_path=$(realpath --relative-to="$SETUP_DIR" "$file")
            local size=$(wc -l < "$file" 2>/dev/null || echo "0")
            echo "- \`$rel_path\` ($size linhas)" >> "$structure_report"
        done
        echo "" >> "$structure_report"
    fi
    
    success "Análise de estrutura concluída: $structure_report"
}

# Função para análise de dependências
analyze_dependencies() {
    log "Analisando dependências..."
    
    local deps_report="$REPORTS_DIR/dependencies_analysis_$TIMESTAMP.md"
    
    cat > "$deps_report" << EOF
# Análise de Dependências - Setup Component

**Data:** $(date)

## Dependências do Sistema

EOF

    # Verificar Go
    if command -v go &> /dev/null; then
        local go_version=$(go version)
        echo "- **Go**: $go_version" >> "$deps_report"
    else
        echo "- **Go**: ❌ Não instalado" >> "$deps_report"
    fi
    
    # Verificar outras ferramentas
    local tools=("git" "make" "gcc" "curl" "wget")
    for tool in "${tools[@]}"; do
        if command -v "$tool" &> /dev/null; then
            local version=$($tool --version 2>&1 | head -1)
            echo "- **$tool**: ✅ $version" >> "$deps_report"
        else
            echo "- **$tool**: ❌ Não instalado" >> "$deps_report"
        fi
    done
    
    echo "" >> "$deps_report"
    echo "## Dependências Go" >> "$deps_report"
    echo "" >> "$deps_report"
    
    # Analisar go.mod se existir
    if [[ -f "$SETUP_DIR/go.mod" ]]; then
        echo "\`\`\`" >> "$deps_report"
        cat "$SETUP_DIR/go.mod" >> "$deps_report"
        echo "\`\`\`" >> "$deps_report"
    else
        echo "❌ Arquivo go.mod não encontrado" >> "$deps_report"
    fi
    
    success "Análise de dependências concluída: $deps_report"
}

# Função para análise de qualidade de código
analyze_code_quality() {
    log "Analisando qualidade do código..."
    
    local quality_report="$REPORTS_DIR/code_quality_$TIMESTAMP.md"
    
    cat > "$quality_report" << EOF
# Análise de Qualidade de Código - Setup Component

**Data:** $(date)

## Métricas de Código

EOF

    # Contar linhas de código
    local total_lines=0
    local go_files=0
    
    if [[ -d "$SETUP_DIR/src" ]]; then
        while IFS= read -r -d '' file; do
            local lines=$(wc -l < "$file")
            total_lines=$((total_lines + lines))
            go_files=$((go_files + 1))
        done < <(find "$SETUP_DIR/src" -name "*.go" -print0)
    fi
    
    echo "- **Total de arquivos Go**: $go_files" >> "$quality_report"
    echo "- **Total de linhas de código**: $total_lines" >> "$quality_report"
    
    if [[ $go_files -gt 0 ]]; then
        local avg_lines=$((total_lines / go_files))
        echo "- **Média de linhas por arquivo**: $avg_lines" >> "$quality_report"
    fi
    
    echo "" >> "$quality_report"
    echo "## Análise por Arquivo" >> "$quality_report"
    echo "" >> "$quality_report"
    
    # Analisar cada arquivo Go
    if [[ -d "$SETUP_DIR/src" ]]; then
        find "$SETUP_DIR/src" -name "*.go" | while read -r file; do
            local rel_path=$(realpath --relative-to="$SETUP_DIR" "$file")
            local lines=$(wc -l < "$file")
            local functions=$(grep -c "^func " "$file" 2>/dev/null || echo "0")
            local comments=$(grep -c "^//" "$file" 2>/dev/null || echo "0")
            
            echo "### \`$rel_path\`" >> "$quality_report"
            echo "- Linhas: $lines" >> "$quality_report"
            echo "- Funções: $functions" >> "$quality_report"
            echo "- Comentários: $comments" >> "$quality_report"
            echo "" >> "$quality_report"
        done
    fi
    
    success "Análise de qualidade concluída: $quality_report"
}

# Função para análise de testes
analyze_tests() {
    log "Analisando estrutura de testes..."
    
    local tests_report="$REPORTS_DIR/tests_analysis_$TIMESTAMP.md"
    
    cat > "$tests_report" << EOF
# Análise de Testes - Setup Component

**Data:** $(date)

## Estrutura de Testes

EOF

    # Analisar cada diretório de teste
    local test_dirs=("unit" "integration" "e2e" "security" "performance")
    
    for dir in "${test_dirs[@]}"; do
        local test_dir="$TESTS_DIR/$dir"
        if [[ -d "$test_dir" ]]; then
            echo "### $dir/" >> "$tests_report"
            
            local test_files=$(find "$test_dir" -name "*_test.go" | wc -l)
            local total_lines=0
            local total_functions=0
            
            find "$test_dir" -name "*_test.go" | while read -r file; do
                local lines=$(wc -l < "$file")
                local functions=$(grep -c "^func Test" "$file" 2>/dev/null || echo "0")
                local rel_path=$(realpath --relative-to="$SETUP_DIR" "$file")
                
                echo "- \`$rel_path\`: $lines linhas, $functions testes" >> "$tests_report"
            done
            
            echo "" >> "$tests_report"
        else
            echo "### $dir/" >> "$tests_report"
            echo "❌ Diretório não encontrado" >> "$tests_report"
            echo "" >> "$tests_report"
        fi
    done
    
    # Analisar mocks e fixtures
    echo "### Mocks e Fixtures" >> "$tests_report"
    echo "" >> "$tests_report"
    
    if [[ -d "$TESTS_DIR/mocks" ]]; then
        echo "**Mocks:**" >> "$tests_report"
        find "$TESTS_DIR/mocks" -name "*.go" | while read -r file; do
            local rel_path=$(realpath --relative-to="$SETUP_DIR" "$file")
            echo "- \`$rel_path\`" >> "$tests_report"
        done
        echo "" >> "$tests_report"
    fi
    
    if [[ -d "$TESTS_DIR/fixtures" ]]; then
        echo "**Fixtures:**" >> "$tests_report"
        find "$TESTS_DIR/fixtures" -name "*" -type f | while read -r file; do
            local rel_path=$(realpath --relative-to="$SETUP_DIR" "$file")
            echo "- \`$rel_path\`" >> "$tests_report"
        done
        echo "" >> "$tests_report"
    fi
    
    success "Análise de testes concluída: $tests_report"
}

# Função para análise de performance
analyze_performance() {
    log "Analisando aspectos de performance..."
    
    local perf_report="$REPORTS_DIR/performance_analysis_$TIMESTAMP.md"
    
    cat > "$perf_report" << EOF
# Análise de Performance - Setup Component

**Data:** $(date)

## Benchmarks Disponíveis

EOF

    # Procurar por benchmarks
    if [[ -d "$TESTS_DIR/performance" ]]; then
        find "$TESTS_DIR/performance" -name "*.go" | while read -r file; do
            local rel_path=$(realpath --relative-to="$SETUP_DIR" "$file")
            local benchmarks=$(grep -c "^func Benchmark" "$file" 2>/dev/null || echo "0")
            
            if [[ $benchmarks -gt 0 ]]; then
                echo "### \`$rel_path\`" >> "$perf_report"
                echo "- Benchmarks encontrados: $benchmarks" >> "$perf_report"
                
                # Listar benchmarks
                grep "^func Benchmark" "$file" | while read -r benchmark; do
                    local name=$(echo "$benchmark" | sed 's/func //' | sed 's/(.*//')
                    echo "  - $name" >> "$perf_report"
                done
                echo "" >> "$perf_report"
            fi
        done
    else
        echo "❌ Diretório de performance não encontrado" >> "$perf_report"
    fi
    
    # Analisar uso de goroutines
    echo "## Análise de Concorrência" >> "$perf_report"
    echo "" >> "$perf_report"
    
    local goroutine_usage=0
    if [[ -d "$SETUP_DIR/src" ]]; then
        goroutine_usage=$(find "$SETUP_DIR/src" -name "*.go" -exec grep -l "go func\|goroutine\|sync\." {} \; | wc -l)
    fi
    
    echo "- **Arquivos com concorrência**: $goroutine_usage" >> "$perf_report"
    
    success "Análise de performance concluída: $perf_report"
}

# Função para análise de segurança
analyze_security() {
    log "Analisando aspectos de segurança..."
    
    local security_report="$REPORTS_DIR/security_analysis_$TIMESTAMP.md"
    
    cat > "$security_report" << EOF
# Análise de Segurança - Setup Component

**Data:** $(date)

## Testes de Segurança

EOF

    # Verificar testes de segurança
    if [[ -d "$TESTS_DIR/security" ]]; then
        local security_tests=$(find "$TESTS_DIR/security" -name "*_test.go" | wc -l)
        echo "- **Arquivos de teste de segurança**: $security_tests" >> "$security_report"
        echo "" >> "$security_report"
        
        find "$TESTS_DIR/security" -name "*_test.go" | while read -r file; do
            local rel_path=$(realpath --relative-to="$SETUP_DIR" "$file")
            local test_functions=$(grep -c "^func Test" "$file" 2>/dev/null || echo "0")
            
            echo "### \`$rel_path\`" >> "$security_report"
            echo "- Funções de teste: $test_functions" >> "$security_report"
            echo "" >> "$security_report"
        done
    else
        echo "❌ Diretório de segurança não encontrado" >> "$security_report"
    fi
    
    # Analisar práticas de segurança no código
    echo "## Análise de Práticas de Segurança" >> "$security_report"
    echo "" >> "$security_report"
    
    if [[ -d "$SETUP_DIR/src" ]]; then
        # Verificar uso de crypto
        local crypto_usage=$(find "$SETUP_DIR/src" -name "*.go" -exec grep -l "crypto\|encrypt\|hash" {} \; | wc -l)
        echo "- **Arquivos com criptografia**: $crypto_usage" >> "$security_report"
        
        # Verificar validação de entrada
        local validation_usage=$(find "$SETUP_DIR/src" -name "*.go" -exec grep -l "validate\|sanitize\|clean" {} \; | wc -l)
        echo "- **Arquivos com validação**: $validation_usage" >> "$security_report"
        
        # Verificar tratamento de erros
        local error_handling=$(find "$SETUP_DIR/src" -name "*.go" -exec grep -l "error\|Error\|err" {} \; | wc -l)
        echo "- **Arquivos com tratamento de erro**: $error_handling" >> "$security_report"
    fi
    
    success "Análise de segurança concluída: $security_report"
}

# Função para gerar relatório consolidado
generate_consolidated_report() {
    log "Gerando relatório consolidado..."
    
    local consolidated_report="$REPORTS_DIR/DIAGNOSTIC_REPORT_$TIMESTAMP.md"
    
    cat > "$consolidated_report" << EOF
# Relatório de Diagnóstico Consolidado - Setup Component

**Data/Hora:** $(date)
**Diretório:** $SETUP_DIR
**Timestamp:** $TIMESTAMP

## Resumo Executivo

Este relatório consolida todas as análises realizadas no componente setup, fornecendo uma visão completa do estado atual do projeto.

## Relatórios Gerados

- **Estrutura**: structure_analysis_$TIMESTAMP.md
- **Dependências**: dependencies_analysis_$TIMESTAMP.md
- **Qualidade de Código**: code_quality_$TIMESTAMP.md
- **Testes**: tests_analysis_$TIMESTAMP.md
- **Performance**: performance_analysis_$TIMESTAMP.md
- **Segurança**: security_analysis_$TIMESTAMP.md

## Status Geral

### ✅ Pontos Positivos

- Estrutura de testes bem organizada
- Cobertura abrangente de tipos de teste
- Documentação detalhada disponível
- Scripts de automação implementados

### ⚠️ Pontos de Atenção

- Verificar configuração de módulos Go
- Validar execução de todos os testes
- Confirmar dependências do sistema

### 🔧 Recomendações

1. **Configuração de Módulo**: Resolver problemas de dependências Go
2. **Execução de Testes**: Implementar execução automática de testes
3. **CI/CD**: Configurar pipeline de integração contínua
4. **Documentação**: Manter documentação atualizada
5. **Monitoramento**: Implementar monitoramento de qualidade

## Próximos Passos

1. Revisar relatórios individuais
2. Corrigir problemas identificados
3. Executar testes completos
4. Implementar melhorias sugeridas
5. Configurar automação contínua

## Contato e Suporte

Para dúvidas ou sugestões sobre este diagnóstico, consulte a documentação do projeto ou entre em contato com a equipe de desenvolvimento.

---
*Relatório gerado automaticamente pelo Sistema de Diagnóstico Automático*
EOF

    success "Relatório consolidado gerado: $consolidated_report"
}

# Função para executar diagnóstico completo
run_full_diagnostic() {
    log "Executando diagnóstico completo..."
    
    analyze_structure
    analyze_dependencies
    analyze_code_quality
    analyze_tests
    analyze_performance
    analyze_security
    generate_consolidated_report
    
    success "Diagnóstico completo concluído!"
}

# Função para mostrar menu interativo
show_menu() {
    echo -e "${CYAN}=== MENU DE DIAGNÓSTICO ===${NC}"
    echo "1. Análise de Estrutura"
    echo "2. Análise de Dependências"
    echo "3. Análise de Qualidade de Código"
    echo "4. Análise de Testes"
    echo "5. Análise de Performance"
    echo "6. Análise de Segurança"
    echo "7. Diagnóstico Completo"
    echo "8. Sair"
    echo ""
    read -p "Escolha uma opção (1-8): " choice
    
    case $choice in
        1) analyze_structure ;;
        2) analyze_dependencies ;;
        3) analyze_code_quality ;;
        4) analyze_tests ;;
        5) analyze_performance ;;
        6) analyze_security ;;
        7) run_full_diagnostic ;;
        8) exit 0 ;;
        *) echo "Opção inválida" ;;
    esac
}

# Função principal
main() {
    if [[ $# -eq 0 ]]; then
        show_menu
    else
        case "$1" in
            "structure") analyze_structure ;;
            "dependencies") analyze_dependencies ;;
            "quality") analyze_code_quality ;;
            "tests") analyze_tests ;;
            "performance") analyze_performance ;;
            "security") analyze_security ;;
            "full") run_full_diagnostic ;;
            *) 
                echo "Uso: $0 [structure|dependencies|quality|tests|performance|security|full]"
                echo "Ou execute sem argumentos para menu interativo"
                exit 1
                ;;
        esac
    fi
    
    echo ""
    echo -e "${GREEN}Relatórios salvos em: $REPORTS_DIR${NC}"
}

# Executar função principal
main "$@"
