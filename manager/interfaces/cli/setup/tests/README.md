# Testes do Componente Setup

Este diretório contém a suíte completa de testes para o componente `setup` do Syntropy CLI, implementada seguindo as regras definidas em `test-suite.rules.md`.

## Estrutura dos Testes

```
tests/
├── unit/                    # Testes unitários
│   ├── setup_test.go       # Testes principais do SetupManager
│   ├── setup_additional_test.go  # Testes adicionais do SetupManager
│   ├── validator_test.go   # Testes do Validator
│   ├── validator_auxiliary_test.go  # Testes auxiliares do Validator
│   ├── validator_windows_test.go    # Testes específicos do Windows
│   ├── validator_linux_test.go      # Testes específicos do Linux
│   ├── validator_darwin_test.go     # Testes específicos do macOS
│   ├── validator_generic_test.go    # Testes genéricos
│   ├── configurator_test.go         # Testes do Configurator
│   ├── key_manager_test.go          # Testes do KeyManager
│   ├── state_manager_test.go        # Testes do StateManager
│   ├── logger_test.go               # Testes do Logger
│   ├── types_test.go                # Testes dos tipos internos
│   ├── utils_test.go                # Testes dos utilitários
│   ├── services_test.go             # Testes dos serviços
│   └── benchmark_test.go            # Testes de benchmark
├── integration/             # Testes de integração
│   └── setup_integration_test.go
├── e2e/                     # Testes end-to-end
│   └── setup_e2e_test.go
├── performance/             # Testes de performance
│   └── setup_performance_test.go
├── security/                # Testes de segurança
│   └── setup_security_test.go
├── fixtures/                # Dados de teste
├── mocks/                   # Implementações mock
├── helpers/                 # Funções auxiliares
├── test_config.yaml         # Configuração dos testes
├── analysis.md              # Análise do código fonte
└── README.md                # Este arquivo
```

## Tipos de Testes

### 1. Testes Unitários (`unit/`)

Testam componentes individuais em isolamento, com 100% de cobertura de código.

**Características:**
- Testam cada função e método individualmente
- Usam mocks para dependências externas
- Verificam casos de sucesso e falha
- Incluem testes de benchmark para performance
- Cobertura obrigatória de 100%

**Arquivos principais:**
- `setup_test.go`: Testes do `SetupManager`
- `validator_test.go`: Testes do `Validator`
- `configurator_test.go`: Testes do `Configurator`
- `key_manager_test.go`: Testes do `KeyManager`
- `state_manager_test.go`: Testes do `StateManager`
- `logger_test.go`: Testes do `Logger`

### 2. Testes de Integração (`integration/`)

Testam a interação entre componentes do sistema.

**Características:**
- Testam fluxos completos entre componentes
- Verificam integração com sistemas externos
- Testam cenários de uso real
- Validam contratos entre interfaces

### 3. Testes End-to-End (`e2e/`)

Testam o fluxo completo do usuário, do início ao fim.

**Características:**
- Simulam uso real do componente
- Testam fluxos completos de setup
- Verificam integridade de dados
- Testam recuperação de erros

### 4. Testes de Performance (`performance/`)

Avaliam a performance e escalabilidade do sistema.

**Características:**
- Medem tempos de execução
- Testam operações concorrentes
- Verificam limites de recursos
- Validam SLAs de performance

### 5. Testes de Segurança (`security/`)

Verificam aspectos de segurança do sistema.

**Características:**
- Testam permissões de arquivos
- Verificam integridade de dados
- Validam geração segura de chaves
- Testam validações de segurança

## Executando os Testes

### Pré-requisitos

- Go 1.19 ou superior
- Acesso ao diretório `src/` (somente leitura)
- Permissões para criar diretórios temporários

### Comandos de Execução

```bash
# Executar todos os testes
go test ./tests/...

# Executar apenas testes unitários
go test -tags=!integration,!e2e,!performance,!security ./tests/unit/...

# Executar apenas testes de integração
go test -tags=integration ./tests/integration/...

# Executar apenas testes end-to-end
go test -tags=e2e ./tests/e2e/...

# Executar apenas testes de performance
go test -tags=performance ./tests/performance/...

# Executar apenas testes de segurança
go test -tags=security ./tests/security/...

# Executar com cobertura
go test -cover ./tests/...

# Executar com cobertura detalhada
go test -coverprofile=coverage.out ./tests/...
go tool cover -html=coverage.out
```

### Build Tags

Os testes usam build tags para separar diferentes tipos:

- `!integration,!e2e,!performance,!security`: Testes unitários
- `integration`: Testes de integração
- `e2e`: Testes end-to-end
- `performance`: Testes de performance
- `security`: Testes de segurança

## Configuração

### Arquivo de Configuração

O arquivo `test_config.yaml` contém configurações específicas para os testes:

- Configurações de ambiente
- Timeouts e limites
- Configurações de cobertura
- Configurações de mock
- Configurações de relatórios

### Variáveis de Ambiente

```bash
# Diretório temporário para testes
export TEST_TEMP_DIR="/tmp/syntropy-setup-tests"

# Nível de log para testes
export TEST_LOG_LEVEL="debug"

# Timeout para operações
export TEST_TIMEOUT="30s"
```

## Cobertura de Código

### Requisitos

- **100% de cobertura de linha**: Todas as linhas de código devem ser executadas
- **100% de cobertura de branch**: Todos os caminhos condicionais devem ser testados
- **100% de cobertura de path**: Todos os caminhos de execução devem ser cobertos

### Verificação de Cobertura

```bash
# Gerar relatório de cobertura
go test -coverprofile=coverage.out ./tests/...

# Visualizar cobertura
go tool cover -html=coverage.out

# Verificar cobertura mínima
go test -covermode=count -coverprofile=coverage.out ./tests/...
go tool cover -func=coverage.out | grep total
```

## Estrutura de Dados de Teste

### Fixtures

O diretório `fixtures/` contém dados de teste reutilizáveis:

- Configurações de exemplo
- Estados de setup válidos
- Dados de validação
- Templates de configuração

### Mocks

O diretório `mocks/` contém implementações mock:

- Mock do `FileService`
- Mock do `NetworkService`
- Mock do `SystemService`
- Mock de validadores específicos do OS

## Relatórios

### Relatórios de Teste

Os relatórios são gerados em formato JSON, XML ou HTML:

```bash
# Gerar relatório JSON
go test -json ./tests/... > test_report.json

# Gerar relatório XML
go test -v ./tests/... | go-junit-report > test_report.xml
```

### Relatórios de Cobertura

```bash
# Relatório HTML
go tool cover -html=coverage.out -o coverage_report.html

# Relatório XML
gocov convert coverage.out | gocov-xml > coverage_report.xml
```

## Troubleshooting

### Problemas Comuns

1. **Falha de permissões**: Verificar se o usuário tem permissões para criar diretórios temporários
2. **Timeout de testes**: Ajustar timeouts no arquivo de configuração
3. **Falha de cobertura**: Verificar se todos os caminhos de código estão sendo testados
4. **Falha de mock**: Verificar se as implementações mock estão corretas

### Logs de Debug

```bash
# Executar com logs detalhados
go test -v ./tests/...

# Executar com logs de debug
TEST_LOG_LEVEL=debug go test -v ./tests/...
```

## Contribuição

### Adicionando Novos Testes

1. Seguir a estrutura existente
2. Usar build tags apropriados
3. Garantir 100% de cobertura
4. Incluir testes de casos de erro
5. Adicionar documentação

### Padrões de Código

- Usar nomes descritivos para testes
- Incluir comentários explicativos
- Seguir convenções de Go
- Usar tabelas de teste quando apropriado
- Limpar recursos após testes

## Referências

- [Regras de Teste](test-suite.rules.md)
- [Análise do Código](analysis.md)
- [Configuração](test_config.yaml)
- [Documentação do Go Testing](https://golang.org/pkg/testing/)
- [Go Test Coverage](https://blog.golang.org/cover)
