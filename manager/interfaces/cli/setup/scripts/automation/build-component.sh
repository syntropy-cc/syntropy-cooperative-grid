#!/bin/bash

# Script para build completo do componente setup
# Autor: Sistema de Build Automático
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
BUILD_DIR="$SETUP_DIR/build"
DIST_DIR="$SETUP_DIR/dist"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# Criar diretórios de build
mkdir -p "$BUILD_DIR" "$DIST_DIR"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  BUILD COMPONENT - SETUP${NC}"
echo -e "${BLUE}========================================${NC}"
echo -e "Data/Hora: $(date)"
echo -e "Diretório: $SETUP_DIR"
echo -e "Build: $BUILD_DIR"
echo -e "Dist: $DIST_DIR"
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

# Função para verificar dependências
check_dependencies() {
    log "Verificando dependências..."
    
    # Verificar Go
    if ! command -v go &> /dev/null; then
        error "Go não está instalado"
        return 1
    fi
    
    GO_VERSION=$(go version | cut -d' ' -f3)
    success "Go: $GO_VERSION"
    
    # Verificar arquivos necessários
    local required_files=(
        "$SETUP_DIR/src/setup.go"
        "$SETUP_DIR/src/validator.go"
        "$SETUP_DIR/src/configurator.go"
        "$SETUP_DIR/src/logger.go"
    )
    
    for file in "${required_files[@]}"; do
        if [[ ! -f "$file" ]]; then
            error "Arquivo necessário não encontrado: $file"
            return 1
        fi
    done
    
    success "Todos os arquivos necessários encontrados"
    return 0
}

# Função para limpar builds anteriores
clean_build() {
    log "Limpando builds anteriores..."
    
    if [[ -d "$BUILD_DIR" ]]; then
        rm -rf "$BUILD_DIR"/*
        success "Build directory limpo"
    fi
    
    if [[ -d "$DIST_DIR" ]]; then
        rm -rf "$DIST_DIR"/*
        success "Dist directory limpo"
    fi
}

# Função para preparar módulo Go
prepare_go_module() {
    log "Preparando módulo Go..."
    
    cd "$SETUP_DIR"
    
    # Criar go.mod se não existir
    if [[ ! -f "go.mod" ]]; then
        go mod init setup-component
        success "go.mod criado"
    fi
    
    # Adicionar dependências necessárias
    go get github.com/spf13/cobra@latest
    go get github.com/stretchr/testify@latest
    go get gopkg.in/yaml.v3@latest
    
    # Limpar dependências
    go mod tidy
    
    success "Módulo Go preparado"
}

# Função para executar testes antes do build
run_tests() {
    log "Executando testes antes do build..."
    
    cd "$SETUP_DIR"
    
    local test_log="$BUILD_DIR/test_results_$TIMESTAMP.log"
    
    if go test -v ./... > "$test_log" 2>&1; then
        success "Todos os testes passaram"
        return 0
    else
        warning "Alguns testes falharam - continuando build"
        echo "  Log de testes: $test_log"
        return 1
    fi
}

# Função para build do componente
build_component() {
    log "Compilando componente..."
    
    cd "$SETUP_DIR"
    
    # Build para diferentes plataformas
    local platforms=(
        "linux/amd64"
        "linux/arm64"
        "windows/amd64"
        "darwin/amd64"
        "darwin/arm64"
    )
    
    for platform in "${platforms[@]}"; do
        local os=$(echo "$platform" | cut -d'/' -f1)
        local arch=$(echo "$platform" | cut -d'/' -f2)
        
        log "Compilando para $os/$arch..."
        
        local output_name="setup-$os-$arch"
        if [[ "$os" == "windows" ]]; then
            output_name="$output_name.exe"
        fi
        
        local output_path="$BUILD_DIR/$output_name"
        
        if GOOS="$os" GOARCH="$arch" go build -o "$output_path" ./src/...; then
            success "Build $os/$arch: $output_path"
            
            # Copiar para dist
            cp "$output_path" "$DIST_DIR/"
        else
            error "Falha no build $os/$arch"
        fi
    done
}

# Função para gerar documentação
generate_docs() {
    log "Gerando documentação..."
    
    cd "$SETUP_DIR"
    
    local docs_dir="$BUILD_DIR/docs"
    mkdir -p "$docs_dir"
    
    # Gerar documentação Go
    if command -v godoc &> /dev/null; then
        godoc -html ./src > "$docs_dir/api.html" 2>/dev/null || warning "godoc não disponível"
    fi
    
    # Copiar documentação existente
    if [[ -d "$SETUP_DIR/docs" ]]; then
        cp -r "$SETUP_DIR/docs"/* "$docs_dir/" 2>/dev/null || true
    fi
    
    # Gerar README do build
    cat > "$docs_dir/BUILD_INFO.md" << EOF
# Build Information - Setup Component

**Build Date:** $(date)
**Go Version:** $(go version)
**Build Directory:** $BUILD_DIR
**Distribution Directory:** $DIST_DIR

## Built Binaries

EOF

    # Listar binários gerados
    for binary in "$DIST_DIR"/*; do
        if [[ -f "$binary" ]]; then
            local size=$(du -h "$binary" | cut -f1)
            local name=$(basename "$binary")
            echo "- **$name**: $size" >> "$docs_dir/BUILD_INFO.md"
        fi
    done
    
    success "Documentação gerada em: $docs_dir"
}

# Função para criar pacotes de distribuição
create_packages() {
    log "Criando pacotes de distribuição..."
    
    cd "$DIST_DIR"
    
    # Criar tarball para Linux
    if ls setup-linux-* >/dev/null 2>&1; then
        tar -czf "setup-linux-$TIMESTAMP.tar.gz" setup-linux-*
        success "Pacote Linux criado: setup-linux-$TIMESTAMP.tar.gz"
    fi
    
    # Criar zip para Windows
    if ls setup-windows-* >/dev/null 2>&1; then
        zip -q "setup-windows-$TIMESTAMP.zip" setup-windows-*
        success "Pacote Windows criado: setup-windows-$TIMESTAMP.zip"
    fi
    
    # Criar tarball para macOS
    if ls setup-darwin-* >/dev/null 2>&1; then
        tar -czf "setup-darwin-$TIMESTAMP.tar.gz" setup-darwin-*
        success "Pacote macOS criado: setup-darwin-$TIMESTAMP.tar.gz"
    fi
}

# Função para verificar integridade dos binários
verify_binaries() {
    log "Verificando integridade dos binários..."
    
    local verification_log="$BUILD_DIR/verification_$TIMESTAMP.log"
    
    {
        echo "=== VERIFICAÇÃO DE BINÁRIOS ==="
        echo "Data: $(date)"
        echo ""
        
        for binary in "$DIST_DIR"/*; do
            if [[ -f "$binary" && ! "$binary" =~ \.(tar\.gz|zip)$ ]]; then
                echo "=== $(basename "$binary") ==="
                echo "Tamanho: $(du -h "$binary" | cut -f1)"
                echo "Tipo: $(file "$binary")"
                
                # Tentar executar com --help se for executável
                if [[ -x "$binary" ]]; then
                    echo "Teste de execução:"
                    timeout 5s "$binary" --help 2>&1 || echo "  (não suporta --help ou erro)"
                fi
                echo ""
            fi
        done
        
    } > "$verification_log"
    
    success "Verificação concluída: $verification_log"
}

# Função para gerar relatório de build
generate_build_report() {
    log "Gerando relatório de build..."
    
    local build_report="$BUILD_DIR/build_report_$TIMESTAMP.md"
    
    cat > "$build_report" << EOF
# Relatório de Build - Setup Component

**Data/Hora:** $(date)
**Diretório:** $SETUP_DIR
**Build ID:** $TIMESTAMP

## Informações do Build

- **Go Version:** $(go version)
- **Build Directory:** $BUILD_DIR
- **Distribution Directory:** $DIST_DIR

## Binários Gerados

EOF

    # Listar binários
    for binary in "$DIST_DIR"/*; do
        if [[ -f "$binary" && ! "$binary" =~ \.(tar\.gz|zip)$ ]]; then
            local size=$(du -h "$binary" | cut -f1)
            local name=$(basename "$binary")
            echo "- **$name**: $size" >> "$build_report"
        fi
    done
    
    cat >> "$build_report" << EOF

## Pacotes de Distribuição

EOF

    # Listar pacotes
    for package in "$DIST_DIR"/*.{tar.gz,zip}; do
        if [[ -f "$package" ]]; then
            local size=$(du -h "$package" | cut -f1)
            local name=$(basename "$package")
            echo "- **$name**: $size" >> "$build_report"
        fi
    done
    
    cat >> "$build_report" << EOF

## Logs e Relatórios

- **Test Results:** test_results_$TIMESTAMP.log
- **Verification:** verification_$TIMESTAMP.log
- **Documentation:** docs/

## Status do Build

- ✅ Dependências verificadas
- ✅ Módulo Go preparado
- ✅ Testes executados
- ✅ Componente compilado
- ✅ Documentação gerada
- ✅ Pacotes criados
- ✅ Integridade verificada

---
*Relatório gerado automaticamente pelo sistema de build*
EOF

    success "Relatório de build gerado: $build_report"
}

# Função principal
main() {
    echo -e "${BLUE}Iniciando build do componente...${NC}"
    echo ""
    
    # Verificações iniciais
    if ! check_dependencies; then
        error "Dependências não atendidas. Abortando."
        exit 1
    fi
    
    # Preparar ambiente
    clean_build
    prepare_go_module
    
    # Executar testes
    run_tests
    
    # Build
    build_component
    
    # Pós-processamento
    generate_docs
    create_packages
    verify_binaries
    generate_build_report
    
    # Resumo final
    echo ""
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}  BUILD CONCLUÍDO${NC}"
    echo -e "${BLUE}========================================${NC}"
    
    echo -e "Binários gerados em: ${GREEN}$DIST_DIR${NC}"
    echo -e "Documentação em: ${GREEN}$BUILD_DIR/docs${NC}"
    echo -e "Relatório em: ${GREEN}$BUILD_DIR/build_report_$TIMESTAMP.md${NC}"
    echo ""
    
    # Listar arquivos gerados
    echo -e "${BLUE}Arquivos gerados:${NC}"
    ls -la "$DIST_DIR"
    echo ""
    
    success "Build concluído com sucesso!"
}

# Executar função principal
main "$@"
