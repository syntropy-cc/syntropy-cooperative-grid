# Setup Component - Testing Documentation

## Testing Philosophy

[MACRO_VIEW]
O testing deste componente contribui para os objetivos gerais de qualidade do projeto Syntropy Cooperative Grid, garantindo confiabilidade e estabilidade do processo crítico de setup em diferentes ambientes e plataformas.
[/MACRO_VIEW]

[MESO_VIEW]
A estratégia de testing integra-se com a infraestrutura de testes compartilhada do módulo CLI Manager, utilizando frameworks e utilitários comuns para manter consistência e reduzir duplicação de código de teste.
[/MESO_VIEW]

[MICRO_VIEW]
A abordagem de testing específica do componente foca em validação de comportamento, integridade de dados e tratamento de erros, com métricas de qualidade que garantem robustez em cenários de produção.
[/MICRO_VIEW]

## Testing Strategy

### Test Pyramid
```
        /\
       /E2E\      5%  - Fluxos críticos de usuário
      /------\
     /Integration\ 25% - Interações entre componentes  
    /------------\
   /     Unit     \ 70% - Funções individuais
  /----------------\
```

### Testing Dimensions

| Dimension | Coverage Goal | Current Coverage | Priority |
|-----------|--------------|------------------|----------|
| Functional | 90% | 85% | High |
| Performance | Key paths | 60% | Medium |
| Security | All inputs | 80% | High |
| Cross-platform | Windows/Linux/macOS | 70% | Medium |

## Test Environment Setup

### Prerequisites

#### System Requirements

| Component | Minimum Version | Recommended | Notes |
|-----------|-----------------|-------------|-------|
| Go Runtime | 1.21 | 1.22 | Required for testing |
| Testify | 1.8.4 | Latest | Testing framework |
| Temp Directory | 100MB | 500MB | For test files |
| Permissions | Write access | Admin rights | For integration tests |

#### Environment Variables

```bash
# Required
TEST_ENV=local              # Ambiente de teste
TEST_TIMEOUT=30s            # Timeout para testes
TEST_TEMP_DIR=/tmp/setup    # Diretório temporário

# Optional
TEST_VERBOSE=true           # Logging detalhado
TEST_PARALLEL=4             # Número de testes paralelos
TEST_COVERAGE=true          # Gerar relatório de cobertura
```

### Installation Steps

```bash
# Step 1: Install test dependencies
go mod download
go mod tidy

# Step 2: Setup test database (if needed)
# Not required for this component

# Step 3: Generate test fixtures
go generate ./tests/fixtures

# Step 4: Verify installation
go test -v ./tests/helpers
Expected output: All helper tests pass
```

## Test Organization

### Directory Structure

```
tests/
├── unit/              # Testes isolados de componentes
│   ├── core/          # Testes de lógica central
│   ├── validator/     # Testes de validação
│   ├── configurator/  # Testes de configuração
│   └── key_manager/   # Testes de gerenciamento de chaves
├── integration/       # Testes de interação entre componentes
│   ├── setup_flow/    # Fluxo completo de setup
│   ├── error_handling/# Tratamento de erros
│   └── state_management/ # Gerenciamento de estado
├── e2e/               # Testes end-to-end de jornadas de usuário
│   ├── scenarios/     # Cenários de usuário
│   └── cross_platform/ # Testes multiplataforma
├── performance/       # Benchmarks de performance
├── security/          # Testes de segurança
├── fixtures/          # Dados de teste
│   ├── valid/         # Casos de teste válidos
│   └── invalid/       # Casos de teste inválidos
├── mocks/             # Implementações mock
└── helpers/           # Utilitários de teste
```

## Running Tests

### All Tests

```bash
# Standard run
go test ./...

# With coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# In watch mode
go test -watch ./...
```

### Test Suites by Type

#### Unit Tests

```bash
# Run all unit tests
go test ./tests/unit/...

# Run specific unit test file
go test ./tests/unit/core/setup_manager_test.go

# Run with debugging
go test -v -run TestSetupManager ./tests/unit/core/
```

#### Integration Tests

```bash
# Prerequisites
export TEST_ENV=integration

# Run integration tests
go test ./tests/integration/...

# Cleanup after tests
go clean -testcache
```

#### End-to-End Tests

```bash
# Start test environment
export TEST_TEMP_DIR=$(mktemp -d)

# Run E2E tests
go test ./tests/e2e/...

# Run specific scenario
go test -run TestCompleteSetupFlow ./tests/e2e/scenarios/
```

#### Performance Tests

```bash
# Run performance benchmarks
go test -bench=. ./tests/performance/

# Generate performance report
go test -bench=. -benchmem ./tests/performance/ > benchmark.txt
```

### Platform-Specific Instructions

#### Windows

```powershell
# PowerShell specific commands
$env:TEST_ENV = "windows"
go test ./tests/...

# CMD specific adjustments
set TEST_TEMP_DIR=%TEMP%\setup-tests
go test ./tests/...

# Common issues and solutions
# Issue: Permission denied on temp files
# Solution: Run as Administrator or use user temp directory
```

#### macOS

```bash
# macOS specific setup
export TEST_ENV=darwin
go test ./tests/...

# Permission requirements
# Tests may require permissions for file system operations
sudo go test ./tests/integration/...
```

#### Linux

```bash
# Distribution-specific notes
export TEST_ENV=linux
go test ./tests/...

# Container-based testing
docker run --rm -v $(pwd):/app -w /app golang:1.21 go test ./tests/...
```

#### CI/CD Environment

```yaml
# Example pipeline configuration
test:
  stage: test
  script:
    - go test -v ./tests/unit/...
    - go test -v ./tests/integration/...
    - go test -v ./tests/e2e/...
  artifacts:
    reports:
      coverage_report:
        coverage_format: cobertura
        path: coverage.xml
```

## Test Data Management

### Fixtures

#### Structure
```
fixtures/
├── valid/           # Dados de teste válidos
│   ├── setup_options.json
│   ├── environment_info.json
│   └── key_pairs.json
├── invalid/         # Dados de teste inválidos
│   ├── corrupted_state.json
│   └── malformed_config.yaml
└── scenarios/       # Cenários de teste
    ├── basic_setup.json
    ├── setup_with_errors.json
    └── cross_platform.json
```

#### Loading Fixtures
```go
func loadFixture(t *testing.T, name string) []byte {
    fixturePath := filepath.Join("tests", "fixtures", name)
    data, err := os.ReadFile(fixturePath)
    require.NoError(t, err)
    return data
}
```

#### Generating Fixtures
```bash
# Generate new fixtures
go run tests/scripts/generate_fixtures.go

# Update existing fixtures
go run tests/scripts/update_fixtures.go --force
```

### Mocks and Stubs

#### Available Mocks

| Mock Name | Purpose | Configuration |
|-----------|---------|---------------|
| MockValidator | Simula validação de ambiente | Configurable success/failure |
| MockConfigurator | Simula geração de configuração | Configurable output paths |
| MockKeyManager | Simula gerenciamento de chaves | Configurable key types |
| MockStateManager | Simula persistência de estado | Configurable storage behavior |

#### Creating Custom Mocks
```go
type MockSetupLogger struct {
    logs []LogEntry
}

func (m *MockSetupLogger) LogStep(step string, data map[string]interface{}) {
    m.logs = append(m.logs, LogEntry{
        Step: step,
        Data: data,
        Time: time.Now(),
    })
}
```

#### Mock Best Practices
- Always reset mocks between tests
- Use realistic mock data that matches production
- Document mock limitations and assumptions
- Provide both success and failure scenarios

## Test Categories Explained

### Unit Tests
**Purpose**: Test individual components in isolation
**Characteristics**:
- No external dependencies
- Fast execution (<100ms per test)
- Deterministic results

**Example Structure**:
```go
func TestSetupManager_Setup(t *testing.T) {
    tests := []struct {
        name    string
        options *types.SetupOptions
        wantErr bool
        errType error
    }{
        {
            name: "valid setup options",
            options: &types.SetupOptions{
                Force: false,
                Verbose: true,
            },
            wantErr: false,
        },
        {
            name:    "nil options should fail",
            options: nil,
            wantErr: true,
            errType: ErrInvalidOptions,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            manager := createTestManager(t)
            
            // Act
            err := manager.Setup(tt.options)
            
            // Assert
            if tt.wantErr {
                assert.Error(t, err)
                assert.IsType(t, tt.errType, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### Integration Tests
**Purpose**: Test component interactions
**Characteristics**:
- May use real dependencies
- Slower than unit tests
- Test data flow between components

**Key Test Scenarios**:

| Scenario | Components Involved | Validation |
|----------|-------------------|------------|
| Complete Setup Flow | SetupManager + All Components | State persistence and validation |
| Error Propagation | SetupManager + Validator | Error handling consistency |
| State Recovery | StateManager + SetupManager | State integrity after failures |

### E2E Tests
**Purpose**: Validate complete user workflows
**Characteristics**:
- Test from user perspective
- Slowest test type
- Catch integration issues

**Critical User Journeys**:
1. **First Time Setup**: User runs setup for the first time → Complete configuration created
2. **Setup with Existing Config**: User runs setup with existing configuration → Backup created, new config applied
3. **Failed Setup Recovery**: Setup fails mid-process → Clean state maintained, retry possible
4. **Cross-Platform Setup**: Setup runs on different OS → Platform-specific configurations created

## Code Coverage

### Coverage Requirements

| Metric | Minimum | Target | Current |
|--------|---------|--------|---------|
| Line Coverage | 70% | 85% | 82% |
| Branch Coverage | 65% | 80% | 78% |
| Function Coverage | 75% | 90% | 88% |
| Statement Coverage | 70% | 85% | 83% |

### Generating Coverage Reports

```bash
# Generate HTML report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Generate JSON report for CI
go test -coverprofile=coverage.out -covermode=count ./...
go tool cover -func=coverage.out

# View coverage report
open coverage.html  # macOS
xdg-open coverage.html  # Linux
start coverage.html  # Windows
```

### Coverage Exclusions
```go
// Patterns excluded from coverage
// Generated files: *_mock.go, *_test.go
// Config files: *.yaml, *.json
// Test files: *_test.go
// Build constraints: +build ignore
```

## Continuous Integration

### Test Pipeline
```
1. Lint/Format Check
   └── Unit Tests
       └── Integration Tests
           └── E2E Tests
               └── Performance Tests
                   └── Coverage Report
```

### CI Configuration
```yaml
stages:
  - test-unit
  - test-integration
  - test-e2e
  - test-performance

test-unit:
  stage: test-unit
  script:
    - go test -v ./tests/unit/...
  coverage: '/coverage: \d+\.\d+%/'

test-integration:
  stage: test-integration
  script:
    - go test -v ./tests/integration/...
  dependencies:
    - test-unit

test-e2e:
  stage: test-e2e
  script:
    - go test -v ./tests/e2e/...
  dependencies:
    - test-integration
```

### Test Reports
**Location**: `tests/reports/`
**Format**: JUnit XML, HTML Coverage
**Retention**: 30 days

## Performance Testing

### Benchmarks

| Operation | Target | Current | Variance Allowed |
|-----------|--------|---------|------------------|
| Setup Complete | <30s | 25s | +/-10% |
| Validation | <5s | 3s | +/-20% |
| Key Generation | <2s | 1.5s | +/-15% |
| State Persistence | <1s | 0.8s | +/-25% |

### Load Testing
```bash
# Run load test
go test -bench=BenchmarkSetupLoad -benchtime=30s ./tests/performance/

# Parameters
Users: 10
Duration: 30 seconds
Ramp-up: 5 seconds
```

### Performance Regression Detection
**Threshold**: 20% performance degradation triggers alert
**Baseline**: Established from 10 successful runs in CI

## Security Testing

### Security Test Suite

| Test Type | Tool/Method | Frequency |
|-----------|-------------|-----------|
| Input Validation | Manual + Fuzzing | Every commit |
| Key Security | Cryptographic validation | Weekly |
| File Permissions | Permission audit | Daily |
| Path Traversal | Injection testing | Every commit |

### Security Test Scenarios
1. **Input validation testing**: Malformed configuration files
2. **Key security testing**: Validation of cryptographic key generation
3. **Permission testing**: File system access boundaries
4. **Path traversal testing**: Directory traversal attempts

## Debugging Tests

### Debug Strategies

| Problem Type | Debug Approach | Tools |
|--------------|----------------|-------|
| Flaky tests | Increase timeouts, add retries | Testify Eventually |
| Slow tests | Profile execution, optimize | Go pprof |
| False failures | Check environment, isolate tests | Test isolation |

### Interactive Debugging
```bash
# Run tests in debug mode
go test -v -race ./tests/...

# Attach debugger
dlv test ./tests/unit/core/setup_manager_test.go -- -test.run TestSetupManager
```

### Logging During Tests
```bash
# Enable verbose logging
export TEST_VERBOSE=true
export SYNTROPY_DEBUG=true

# Log locations
Test logs: tests/logs/
Application logs: .syntropy/logs/
```

## Test Maintenance

### Test Health Metrics

| Metric | Good | Warning | Critical |
|--------|------|---------|----------|
| Flaky Test Rate | <1% | 1-5% | >5% |
| Test Runtime | <5min | 5-15min | >15min |
| Test Failures | <2% | 2-5% | >5% |

### Test Review Checklist
- [ ] Test name clearly describes scenario
- [ ] Test is independent (no order dependencies)
- [ ] Test data is properly cleaned up
- [ ] Assertions are specific and meaningful
- [ ] Test covers edge cases
- [ ] Performance impact is acceptable

### Updating Tests
**When to update**:
- API changes in setup methods
- Business logic changes in validation
- Bug fixes (add regression test)
- Performance improvements

**Update process**:
1. Run existing tests to establish baseline
2. Make changes to component
3. Update affected tests
4. Verify all tests pass
5. Update documentation

## Troubleshooting

### Common Issues

| Issue | Platform | Symptoms | Solution |
|-------|----------|----------|----------|
| Permission denied | All | Test failures on file operations | Run with appropriate permissions |
| Temp directory full | All | Tests fail with disk space errors | Clean temp directory, increase space |
| Race conditions | All | Intermittent test failures | Add proper synchronization |
| Mock not working | All | Tests pass but behavior unexpected | Check mock configuration |

### FAQ

**Q: Tests pass locally but fail in CI**
A: Common causes: Different Go versions, missing environment variables, permission differences, timing issues in CI environment

**Q: How to run a single test?**
A: `go test -run TestSpecificTest ./path/to/test/file.go`

**Q: Tests are running slowly**
A: Check for: Unnecessary file I/O, missing parallel execution, large test data sets, inefficient mocks

## Test Documentation

### Writing Test Descriptions
```go
// GOOD: "should return error when setup options are nil"
// BAD: "test nil case"

// GOOD: "should complete setup in under 30 seconds"  
// BAD: "performance test"

// GOOD: "should create backup before overwriting existing configuration"
// BAD: "backup test"
```

### Test Comments
```go
func TestSetupManager_Setup(t *testing.T) {
    // GIVEN: A setup manager with valid configuration
    manager := createTestManager(t)
    
    // WHEN: Setup is called with valid options
    err := manager.Setup(validOptions)
    
    // THEN: Setup should complete without errors
    assert.NoError(t, err)
}
```

## Resources
- [Go Testing Best Practices](https://golang.org/pkg/testing/)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Test Report Dashboard](http://localhost:8080/test-reports)
- [Performance Baseline History](http://localhost:8080/performance)

## Learning Resources
Para princípios de desenvolvimento orientado a testes e filosofia de testing, veja [LEARN.md](./LEARN.md).





