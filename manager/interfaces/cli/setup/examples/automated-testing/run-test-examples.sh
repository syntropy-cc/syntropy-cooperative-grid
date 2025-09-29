#!/bin/bash

# Exemplos de Teste Automático para Setup Component
# Autor: Sistema de Exemplos
# Data: $(date)

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Configurações
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SETUP_DIR="$(dirname "$(dirname "$SCRIPT_DIR")")"
EXAMPLES_DIR="$SCRIPT_DIR"
RESULTS_DIR="$EXAMPLES_DIR/results"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# Criar diretório de resultados
mkdir -p "$RESULTS_DIR"

echo -e "${PURPLE}========================================${NC}"
echo -e "${PURPLE}  EXEMPLOS DE TESTE AUTOMÁTICO${NC}"
echo -e "${PURPLE}  SETUP COMPONENT${NC}"
echo -e "${PURPLE}========================================${NC}"
echo -e "Data/Hora: $(date)"
echo -e "Diretório: $SETUP_DIR"
echo -e "Resultados: $RESULTS_DIR"
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

# Exemplo 1: Teste de Setup Básico
example_basic_setup() {
    log "Executando Exemplo 1: Teste de Setup Básico"
    
    local example_log="$RESULTS_DIR/example1_basic_setup_$TIMESTAMP.log"
    
    cat > "$example_log" << EOF
=== EXEMPLO 1: TESTE DE SETUP BÁSICO ===
Data: $(date)
Descrição: Testa o fluxo básico de setup do componente

=== CENÁRIO ===
- Criar novo SetupManager
- Executar setup com opções padrão
- Verificar resultado

=== IMPLEMENTAÇÃO ===
package main

import (
    "fmt"
    "log"
    "github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src"
    "github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src/internal/types"
)

func main() {
    // Criar SetupManager
    manager, err := src.NewSetupManager()
    if err != nil {
        log.Fatalf("Erro ao criar SetupManager: %v", err)
    }
    
    // Configurar opções de setup
    options := &types.SetupOptions{
        Force: false,
        Verbose: true,
        ConfigPath: "./config",
    }
    
    // Executar setup
    result, err := manager.Setup(options)
    if err != nil {
        log.Fatalf("Erro no setup: %v", err)
    }
    
    // Verificar resultado
    if result.Success {
        fmt.Printf("Setup concluído com sucesso em %v\n", result.Duration)
        fmt.Printf("Configuração criada em: %s\n", result.ConfigPath)
    } else {
        log.Fatalf("Setup falhou: %s", result.Error)
    }
}

=== RESULTADO ESPERADO ===
- SetupManager criado com sucesso
- Setup executado sem erros
- Configuração criada no diretório especificado
- Tempo de execução registrado

=== COMO EXECUTAR ===
go run example1_basic_setup.go

EOF

    success "Exemplo 1 criado: $example_log"
}

# Exemplo 2: Teste de Validação
example_validation_test() {
    log "Executando Exemplo 2: Teste de Validação"
    
    local example_log="$RESULTS_DIR/example2_validation_$TIMESTAMP.log"
    
    cat > "$example_log" << EOF
=== EXEMPLO 2: TESTE DE VALIDAÇÃO ===
Data: $(date)
Descrição: Testa o sistema de validação do ambiente

=== CENÁRIO ===
- Criar Validator
- Executar validação completa
- Verificar problemas encontrados
- Aplicar correções automáticas

=== IMPLEMENTAÇÃO ===
package main

import (
    "fmt"
    "log"
    "github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src"
)

func main() {
    // Criar Validator
    validator, err := src.NewValidator()
    if err != nil {
        log.Fatalf("Erro ao criar Validator: %v", err)
    }
    
    // Executar validação completa
    result, err := validator.ValidateAll()
    if err != nil {
        log.Fatalf("Erro na validação: %v", err)
    }
    
    // Mostrar resultados
    fmt.Printf("Validação concluída:\n")
    fmt.Printf("- Ambiente: %s\n", getStatus(result.Environment))
    fmt.Printf("- Dependências: %s\n", getStatus(result.Dependencies))
    fmt.Printf("- Rede: %s\n", getStatus(result.Network))
    fmt.Printf("- Permissões: %s\n", getStatus(result.Permissions))
    
    // Aplicar correções se necessário
    if !result.Success {
        fmt.Println("Aplicando correções automáticas...")
        fixResult, err := validator.FixIssues(result.Issues)
        if err != nil {
            log.Fatalf("Erro ao aplicar correções: %v", err)
        }
        
        if fixResult.Success {
            fmt.Println("Correções aplicadas com sucesso!")
        } else {
            fmt.Printf("Algumas correções falharam: %s\n", fixResult.Error)
        }
    }
}

func getStatus(valid bool) string {
    if valid {
        return "✅ Válido"
    }
    return "❌ Inválido"
}

=== RESULTADO ESPERADO ===
- Validator criado com sucesso
- Validação executada para todos os aspectos
- Status de cada validação exibido
- Correções automáticas aplicadas se necessário

=== COMO EXECUTAR ===
go run example2_validation.go

EOF

    success "Exemplo 2 criado: $example_log"
}

# Exemplo 3: Teste de Configuração
example_configuration_test() {
    log "Executando Exemplo 3: Teste de Configuração"
    
    local example_log="$RESULTS_DIR/example3_configuration_$TIMESTAMP.log"
    
    cat > "$example_log" << EOF
=== EXEMPLO 3: TESTE DE CONFIGURAÇÃO ===
Data: $(date)
Descrição: Testa o sistema de configuração e gerenciamento de chaves

=== CENÁRIO ===
- Criar Configurator
- Gerar configuração personalizada
- Criar estrutura de diretórios
- Gerar e gerenciar chaves
- Validar configuração

=== IMPLEMENTAÇÃO ===
package main

import (
    "fmt"
    "log"
    "github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src"
    "github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src/internal/types"
)

func main() {
    // Criar Configurator
    configurator, err := src.NewConfigurator()
    if err != nil {
        log.Fatalf("Erro ao criar Configurator: %v", err)
    }
    
    // Configurar opções
    configOptions := &types.ConfigOptions{
        OutputPath: "./test-config",
        GenerateKeys: true,
        BackupExisting: true,
    }
    
    // Gerar configuração
    fmt.Println("Gerando configuração...")
    config, err := configurator.GenerateConfig(configOptions)
    if err != nil {
        log.Fatalf("Erro ao gerar configuração: %v", err)
    }
    
    // Criar estrutura
    fmt.Println("Criando estrutura de diretórios...")
    err = configurator.CreateStructure(config)
    if err != nil {
        log.Fatalf("Erro ao criar estrutura: %v", err)
    }
    
    // Gerar chaves
    fmt.Println("Gerando chaves...")
    keys, err := configurator.GenerateKeys()
    if err != nil {
        log.Fatalf("Erro ao gerar chaves: %v", err)
    }
    
    // Validar configuração
    fmt.Println("Validando configuração...")
    isValid, err := configurator.ValidateConfig(config)
    if err != nil {
        log.Fatalf("Erro na validação: %v", err)
    }
    
    if isValid {
        fmt.Println("✅ Configuração válida!")
        fmt.Printf("Configuração salva em: %s\n", config.Path)
        fmt.Printf("Chaves geradas: %d\n", len(keys))
    } else {
        fmt.Println("❌ Configuração inválida!")
    }
}

=== RESULTADO ESPERADO ===
- Configurator criado com sucesso
- Configuração gerada com opções personalizadas
- Estrutura de diretórios criada
- Chaves geradas e armazenadas
- Configuração validada com sucesso

=== COMO EXECUTAR ===
go run example3_configuration.go

EOF

    success "Exemplo 3 criado: $example_log"
}

# Exemplo 4: Teste de Performance
example_performance_test() {
    log "Executando Exemplo 4: Teste de Performance"
    
    local example_log="$RESULTS_DIR/example4_performance_$TIMESTAMP.log"
    
    cat > "$example_log" << EOF
=== EXEMPLO 4: TESTE DE PERFORMANCE ===
Data: $(date)
Descrição: Testa a performance do componente com diferentes cargas

=== CENÁRIO ===
- Executar setup múltiplas vezes
- Medir tempo de execução
- Testar com diferentes tamanhos de configuração
- Verificar uso de memória

=== IMPLEMENTAÇÃO ===
package main

import (
    "fmt"
    "log"
    "runtime"
    "time"
    "github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src"
    "github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src/internal/types"
)

func main() {
    // Criar SetupManager
    manager, err := src.NewSetupManager()
    if err != nil {
        log.Fatalf("Erro ao criar SetupManager: %v", err)
    }
    
    // Teste 1: Setup simples
    fmt.Println("=== TESTE 1: SETUP SIMPLES ===")
    testSimpleSetup(manager)
    
    // Teste 2: Setup com configuração grande
    fmt.Println("\n=== TESTE 2: CONFIGURAÇÃO GRANDE ===")
    testLargeConfig(manager)
    
    // Teste 3: Setup concorrente
    fmt.Println("\n=== TESTE 3: SETUP CONCORRENTE ===")
    testConcurrentSetup(manager)
    
    // Mostrar estatísticas de memória
    showMemoryStats()
}

func testSimpleSetup(manager *src.SetupManager) {
    start := time.Now()
    
    options := &types.SetupOptions{
        Force: false,
        Verbose: false,
    }
    
    result, err := manager.Setup(options)
    duration := time.Since(start)
    
    if err != nil {
        log.Printf("Erro no setup simples: %v", err)
        return
    }
    
    fmt.Printf("Setup simples: %v\n", duration)
    fmt.Printf("Sucesso: %t\n", result.Success)
}

func testLargeConfig(manager *src.SetupManager) {
    start := time.Now()
    
    options := &types.SetupOptions{
        Force: false,
        Verbose: false,
        ConfigPath: "./large-config",
    }
    
    result, err := manager.Setup(options)
    duration := time.Since(start)
    
    if err != nil {
        log.Printf("Erro no setup grande: %v", err)
        return
    }
    
    fmt.Printf("Setup grande: %v\n", duration)
    fmt.Printf("Sucesso: %t\n", result.Success)
}

func testConcurrentSetup(manager *src.SetupManager) {
    const numGoroutines = 5
    results := make(chan time.Duration, numGoroutines)
    
    start := time.Now()
    
    for i := 0; i < numGoroutines; i++ {
        go func(id int) {
            options := &types.SetupOptions{
                Force: false,
                Verbose: false,
                ConfigPath: fmt.Sprintf("./concurrent-config-%d", id),
            }
            
            _, err := manager.Setup(options)
            if err != nil {
                log.Printf("Erro no setup concorrente %d: %v", id, err)
            }
            
            results <- time.Since(start)
        }(i)
    }
    
    // Coletar resultados
    var totalDuration time.Duration
    for i := 0; i < numGoroutines; i++ {
        duration := <-results
        totalDuration += duration
    }
    
    avgDuration := totalDuration / numGoroutines
    fmt.Printf("Setup concorrente (média): %v\n", avgDuration)
}

func showMemoryStats() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("\n=== ESTATÍSTICAS DE MEMÓRIA ===\n")
    fmt.Printf("Memória alocada: %d KB\n", m.Alloc/1024)
    fmt.Printf("Total de alocações: %d\n", m.Mallocs)
    fmt.Printf("Goroutines: %d\n", runtime.NumGoroutine())
}

=== RESULTADO ESPERADO ===
- Tempos de execução medidos para diferentes cenários
- Performance comparada entre configurações simples e grandes
- Teste de concorrência executado
- Estatísticas de memória exibidas

=== COMO EXECUTAR ===
go run example4_performance.go

=== BENCHMARKS ===
Para executar benchmarks oficiais:
go test -bench=. ./tests/performance/...

EOF

    success "Exemplo 4 criado: $example_log"
}

# Exemplo 5: Teste de Segurança
example_security_test() {
    log "Executando Exemplo 5: Teste de Segurança"
    
    local example_log="$RESULTS_DIR/example5_security_$TIMESTAMP.log"
    
    cat > "$example_log" << EOF
=== EXEMPLO 5: TESTE DE SEGURANÇA ===
Data: $(date)
Descrição: Testa aspectos de segurança do componente

=== CENÁRIO ===
- Testar validação de entrada
- Verificar permissões de arquivo
- Testar proteção contra ataques comuns
- Validar gerenciamento de chaves

=== IMPLEMENTAÇÃO ===
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src"
    "github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src/internal/types"
)

func main() {
    fmt.Println("=== TESTES DE SEGURANÇA ===")
    
    // Teste 1: Validação de entrada
    testInputValidation()
    
    // Teste 2: Permissões de arquivo
    testFilePermissions()
    
    // Teste 3: Proteção contra directory traversal
    testDirectoryTraversal()
    
    // Teste 4: Gerenciamento seguro de chaves
    testSecureKeyManagement()
    
    fmt.Println("\n✅ Todos os testes de segurança concluídos!")
}

func testInputValidation() {
    fmt.Println("\n--- Teste: Validação de Entrada ---")
    
    // Testar entradas maliciosas
    maliciousInputs := []string{
        "../../../etc/passwd",
        "'; DROP TABLE users; --",
        "<script>alert('xss')</script>",
        "rm -rf /",
    }
    
    for _, input := range maliciousInputs {
        if isValidInput(input) {
            fmt.Printf("❌ Entrada maliciosa aceita: %s\n", input)
        } else {
            fmt.Printf("✅ Entrada maliciosa rejeitada: %s\n", input)
        }
    }
}

func testFilePermissions() {
    fmt.Println("\n--- Teste: Permissões de Arquivo ---")
    
    // Criar arquivo temporário
    tempFile := "test_permissions.tmp"
    file, err := os.Create(tempFile)
    if err != nil {
        log.Printf("Erro ao criar arquivo: %v", err)
        return
    }
    defer os.Remove(tempFile)
    defer file.Close()
    
    // Verificar permissões
    info, err := file.Stat()
    if err != nil {
        log.Printf("Erro ao obter info do arquivo: %v", err)
        return
    }
    
    mode := info.Mode()
    if mode&0o777 == 0o600 {
        fmt.Println("✅ Permissões corretas (600)")
    } else {
        fmt.Printf("❌ Permissões incorretas: %o\n", mode&0o777)
    }
}

func testDirectoryTraversal() {
    fmt.Println("\n--- Teste: Directory Traversal ---")
    
    // Testar caminhos maliciosos
    maliciousPaths := []string{
        "../../../etc/passwd",
        "..\\..\\..\\windows\\system32",
        "/etc/shadow",
        "C:\\Windows\\System32",
    }
    
    for _, path := range maliciousPaths {
        if isSecurePath(path) {
            fmt.Printf("✅ Caminho seguro: %s\n", path)
        } else {
            fmt.Printf("❌ Caminho inseguro: %s\n", path)
        }
    }
}

func testSecureKeyManagement() {
    fmt.Println("\n--- Teste: Gerenciamento de Chaves ---")
    
    // Criar KeyManager
    keyManager, err := src.NewKeyManager()
    if err != nil {
        log.Printf("Erro ao criar KeyManager: %v", err)
        return
    }
    
    // Gerar par de chaves
    keyPair, err := keyManager.GenerateKeyPair()
    if err != nil {
        log.Printf("Erro ao gerar chaves: %v", err)
        return
    }
    
    // Verificar integridade
    isValid, err := keyManager.VerifyKeyIntegrity(keyPair)
    if err != nil {
        log.Printf("Erro na verificação: %v", err)
        return
    }
    
    if isValid {
        fmt.Println("✅ Chaves geradas e validadas com sucesso")
    } else {
        fmt.Println("❌ Falha na validação das chaves")
    }
}

// Funções auxiliares
func isValidInput(input string) bool {
    // Implementar validação de entrada
    dangerousChars := []string{"../", "..\\", "<script", "DROP TABLE", "rm -rf"}
    
    for _, char := range dangerousChars {
        if contains(input, char) {
            return false
        }
    }
    
    return true
}

func isSecurePath(path string) bool {
    // Verificar se o caminho é seguro
    cleanPath := filepath.Clean(path)
    return !contains(cleanPath, "..") && !contains(cleanPath, "/etc/") && !contains(cleanPath, "\\Windows\\")
}

func contains(s, substr string) bool {
    return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsInMiddle(s, substr))))
}

func containsInMiddle(s, substr string) bool {
    for i := 0; i <= len(s)-len(substr); i++ {
        if s[i:i+len(substr)] == substr {
            return true
        }
    }
    return false
}

=== RESULTADO ESPERADO ===
- Entradas maliciosas rejeitadas
- Permissões de arquivo corretas
- Proteção contra directory traversal
- Chaves geradas e validadas com segurança

=== COMO EXECUTAR ===
go run example5_security.go

=== TESTES OFICIAIS ===
Para executar testes oficiais de segurança:
go test -v ./tests/security/...

EOF

    success "Exemplo 5 criado: $example_log"
}

# Exemplo 6: Teste de Integração Completa
example_integration_test() {
    log "Executando Exemplo 6: Teste de Integração Completa"
    
    local example_log="$RESULTS_DIR/example6_integration_$TIMESTAMP.log"
    
    cat > "$example_log" << EOF
=== EXEMPLO 6: TESTE DE INTEGRAÇÃO COMPLETA ===
Data: $(date)
Descrição: Testa o fluxo completo de integração do componente

=== CENÁRIO ===
- Executar fluxo completo de setup
- Testar todas as integrações
- Verificar estado final
- Executar cleanup

=== IMPLEMENTAÇÃO ===
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src"
    "github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src/internal/types"
)

func main() {
    fmt.Println("=== TESTE DE INTEGRAÇÃO COMPLETA ===")
    
    // Fase 1: Validação
    fmt.Println("\n--- FASE 1: VALIDAÇÃO ---")
    if !runValidationPhase() {
        log.Fatal("Fase de validação falhou")
    }
    
    // Fase 2: Setup
    fmt.Println("\n--- FASE 2: SETUP ---")
    if !runSetupPhase() {
        log.Fatal("Fase de setup falhou")
    }
    
    // Fase 3: Verificação
    fmt.Println("\n--- FASE 3: VERIFICAÇÃO ---")
    if !runVerificationPhase() {
        log.Fatal("Fase de verificação falhou")
    }
    
    // Fase 4: Teste de Operação
    fmt.Println("\n--- FASE 4: TESTE DE OPERAÇÃO ---")
    if !runOperationPhase() {
        log.Fatal("Fase de operação falhou")
    }
    
    // Fase 5: Cleanup
    fmt.Println("\n--- FASE 5: CLEANUP ---")
    if !runCleanupPhase() {
        log.Fatal("Fase de cleanup falhou")
    }
    
    fmt.Println("\n✅ Teste de integração completo concluído com sucesso!")
}

func runValidationPhase() bool {
    validator, err := src.NewValidator()
    if err != nil {
        log.Printf("Erro ao criar validator: %v", err)
        return false
    }
    
    result, err := validator.ValidateAll()
    if err != nil {
        log.Printf("Erro na validação: %v", err)
        return false
    }
    
    if !result.Success {
        fmt.Println("Aplicando correções...")
        fixResult, err := validator.FixIssues(result.Issues)
        if err != nil || !fixResult.Success {
            log.Printf("Erro ao aplicar correções: %v", err)
            return false
        }
    }
    
    fmt.Println("✅ Validação concluída")
    return true
}

func runSetupPhase() bool {
    manager, err := src.NewSetupManager()
    if err != nil {
        log.Printf("Erro ao criar manager: %v", err)
        return false
    }
    
    options := &types.SetupOptions{
        Force: false,
        Verbose: true,
        ConfigPath: "./integration-test-config",
    }
    
    result, err := manager.Setup(options)
    if err != nil {
        log.Printf("Erro no setup: %v", err)
        return false
    }
    
    if !result.Success {
        log.Printf("Setup falhou: %s", result.Error)
        return false
    }
    
    fmt.Printf("✅ Setup concluído em %v\n", result.Duration)
    return true
}

func runVerificationPhase() bool {
    manager, err := src.NewSetupManager()
    if err != nil {
        log.Printf("Erro ao criar manager: %v", err)
        return false
    }
    
    status, err := manager.Status()
    if err != nil {
        log.Printf("Erro ao verificar status: %v", err)
        return false
    }
    
    if !status.IsSetup {
        log.Println("Setup não está ativo")
        return false
    }
    
    fmt.Println("✅ Verificação de status concluída")
    return true
}

func runOperationPhase() bool {
    // Simular operações normais
    fmt.Println("Executando operações de teste...")
    
    // Aguardar um pouco para simular operação
    time.Sleep(2 * time.Second)
    
    fmt.Println("✅ Operações de teste concluídas")
    return true
}

func runCleanupPhase() bool {
    manager, err := src.NewSetupManager()
    if err != nil {
        log.Printf("Erro ao criar manager: %v", err)
        return false
    }
    
    result, err := manager.Reset()
    if err != nil {
        log.Printf("Erro no reset: %v", err)
        return false
    }
    
    if !result.Success {
        log.Printf("Reset falhou: %s", result.Error)
        return false
    }
    
    fmt.Println("✅ Cleanup concluído")
    return true
}

=== RESULTADO ESPERADO ===
- Todas as fases executadas com sucesso
- Validação, setup, verificação e cleanup funcionando
- Fluxo completo de integração testado
- Estado limpo após o teste

=== COMO EXECUTAR ===
go run example6_integration.go

=== TESTES OFICIAIS ===
Para executar testes oficiais de integração:
go test -v ./tests/integration/...

EOF

    success "Exemplo 6 criado: $example_log"
}

# Função para criar script de execução de todos os exemplos
create_run_all_examples() {
    log "Criando script para executar todos os exemplos"
    
    local run_all_script="$EXAMPLES_DIR/run-all-examples.sh"
    
    cat > "$run_all_script" << 'EOF'
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
EOF

    chmod +x "$run_all_script"
    success "Script de execução criado: $run_all_script"
}

# Função para criar README dos exemplos
create_examples_readme() {
    log "Criando README dos exemplos"
    
    local readme_file="$EXAMPLES_DIR/README.md"
    
    cat > "$readme_file" << EOF
# Exemplos de Teste Automático - Setup Component

Este diretório contém exemplos práticos de como testar o componente setup de forma automática.

## Exemplos Disponíveis

### 1. Teste de Setup Básico
- **Arquivo**: \`example1_basic_setup.go\`
- **Descrição**: Demonstra o fluxo básico de setup
- **Executar**: \`go run example1_basic_setup.go\`

### 2. Teste de Validação
- **Arquivo**: \`example2_validation.go\`
- **Descrição**: Testa o sistema de validação do ambiente
- **Executar**: \`go run example2_validation.go\`

### 3. Teste de Configuração
- **Arquivo**: \`example3_configuration.go\`
- **Descrição**: Testa geração de configuração e chaves
- **Executar**: \`go run example3_configuration.go\`

### 4. Teste de Performance
- **Arquivo**: \`example4_performance.go\`
- **Descrição**: Testa performance com diferentes cargas
- **Executar**: \`go run example4_performance.go\`

### 5. Teste de Segurança
- **Arquivo**: \`example5_security.go\`
- **Descrição**: Testa aspectos de segurança
- **Executar**: \`go run example5_security.go\`

### 6. Teste de Integração Completa
- **Arquivo**: \`example6_integration.go\`
- **Descrição**: Testa fluxo completo de integração
- **Executar**: \`go run example6_integration.go\`

## Executar Todos os Exemplos

Para executar todos os exemplos de uma vez:

\`\`\`bash
./run-all-examples.sh
\`\`\`

## Resultados

Os resultados dos exemplos são salvos no diretório \`results/\` com timestamp para identificação.

## Pré-requisitos

- Go 1.21 ou superior
- Dependências do projeto instaladas
- Permissões adequadas para criação de arquivos

## Personalização

Cada exemplo pode ser personalizado modificando as variáveis e opções conforme necessário para seu ambiente específico.

## Troubleshooting

Se algum exemplo falhar:

1. Verifique se todas as dependências estão instaladas
2. Confirme que o componente setup está compilado
3. Verifique permissões de arquivo
4. Consulte os logs de erro para detalhes específicos

---
*Exemplos gerados automaticamente pelo Sistema de Teste Automático*
EOF

    success "README dos exemplos criado: $readme_file"
}

# Função principal
main() {
    echo -e "${BLUE}Criando exemplos de teste automático...${NC}"
    echo ""
    
    # Criar todos os exemplos
    example_basic_setup
    example_validation_test
    example_configuration_test
    example_performance_test
    example_security_test
    example_integration_test
    
    # Criar scripts auxiliares
    create_run_all_examples
    create_examples_readme
    
    # Resumo final
    echo ""
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}  EXEMPLOS CRIADOS COM SUCESSO${NC}"
    echo -e "${GREEN}========================================${NC}"
    
    echo -e "Exemplos criados em: ${BLUE}$EXAMPLES_DIR${NC}"
    echo -e "Resultados salvos em: ${BLUE}$RESULTS_DIR${NC}"
    echo ""
    
    echo -e "${BLUE}Arquivos criados:${NC}"
    ls -la "$EXAMPLES_DIR"
    echo ""
    
    echo -e "${YELLOW}Para executar todos os exemplos:${NC}"
    echo -e "  ${BLUE}cd $EXAMPLES_DIR${NC}"
    echo -e "  ${BLUE}./run-all-examples.sh${NC}"
    echo ""
    
    success "Exemplos de teste automático criados com sucesso!"
}

# Executar função principal
main "$@"
