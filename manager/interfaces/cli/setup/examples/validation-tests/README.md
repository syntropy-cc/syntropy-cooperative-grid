# Exemplos de Testes de Validação - Syntropy CLI

Este diretório contém ferramentas de teste e validação avançadas para o setup component do Syntropy CLI, incluindo testes de performance, validação de ambiente e cenários de teste end-to-end.

## Visão Geral

Os scripts de validação permitem:
- Testagem automática do ambiente pronto para setup
- Testes de performance do sistema com diferentes cargas
- Validação de infraestrutura em escala
- Cenários de teste para environments de produção
- Diagnósticos para troubleshoot de problemas

## Estrutura dos Arquivos

- `test-environment.sh` - Script principal de validação de ambiente
- `performance-test.sh` - Testes de performance e capacidade do sistema  
- `README.md` - Esta documentação completa

## Testes Disponíveis

### 1. Testes de Ambiente

**Sistema base:**
```bash
./test-environment.sh
```

**Validação específica por SO:**
```bash
./test-environment.sh --platform linux
./test-environment.sh --platform windows  
./test-environment.sh --platform macos
```

**Testes personalizáveis:**
```bash
./test-environment.sh --check-connectivity --check-disk --check-memory
./test-environment.sh --full-suite
./test-environment.sh --silent --output-file results.json
```

### 2. Testes de Performance

**Performance básico:**
```bash
./performance-test.sh
```

**Teste com carga:**
```bash
./performance-test.sh --concurrent-jobs 10
./performance-test.sh --duration 5m --max-load 75
```

**Benchmark empresarial:**
```bash
./performance-test.sh --enterprise-mode --stress-test --report
```

### 3. Validação de Configuração

**Validar arquivo de configuração:**
```bash
./validate-config.sh ../advanced-setup/custom-config.yaml
./validate-config.sh manager.yaml --strict
```

**Validar topologia de rede:**
```bash
./validate-network.sh ../advanced-setup/network-topology.yaml
```

## Tipos de Testes Implementados

### Testes de Ambiente
- ✅ Detecção de sistema operacional e versão
- ✅ Validação de recursos computacionais (CPU, RAM, disk)
- ✅ Conectividade de rede e proxy
- ✅ Permissões de sistema
- ✅ Verificação de dependências
- ✅ Testes de segurança básicos

### Testes de Performance  
- ✅ Latência de rede e API
- ✅ Throughput máximo do sistema
- ✅ Performance I/O (disk, memory)
- ✅ Disponibilidade durante alta carga
- ✅ Stress test para cenários extremos
- ✅ Benchmarking comparativo

### Testes de Integração
- ✅ Setup completo end-to-end
- ✅ Funcionalidades avançadas
- ✅ Estado posterior ao setup
- ✅ Integração com componentes internos
- ✅ Rollback e recovery
- ✅ Validação de configurações

### Testes de Segurança
- ✅ Validação de protocolos de criptografia
- ✅ Testes de firewall e ACLs
- ✅ Permissões de arquivos
- ✅ Certificados e chaves  
- ✅ Acessos delegados e RBAC

## Cenários de Execução

### Desenvolvimento Local
Para desenvolvimento e debug local:

```bash
# Validação rápida para desenvolvimento  
./test-environment.sh --quick --developer

# Teste completo incluindo performance
./test-environment.sh --full-suite
./performance-test.sh --dev-mode
```

### Ambiente de Staging
Para testes em ambiente de staging:

```bash
# Teste de staging completo
export TEST_ENVIRONMENT=staging
./test-environment.sh --environment tagging
./performance-test.sh --staging-load

# Validação integrada com infraestrutura externa
./integration-test.sh
```

### Ambiente de Produção
Para validação de pré-produção:

```bash
# Testes enterprise robustos
export TEST_ENVIRONMENT=production
./test-environment.sh --enterprise --non-destructive
./performance-test.sh --production-benchmark

# Teste de alta disponibilidade
./ha-simulation.sh
```

## Scripts Integrados

### Validação de Rede
- Diagnóstico de conectividade IPv4/IPv6
- Teste de proxy e túneis corporativo
- Validação de DNS e service discovery
- QoS e latency testing

### Validação de Segurança
- Auditoria mínima de segurança
- Validação de configurações SSL/TLS
- Teste de chaves e certificados
- Verificação de governança de acesso

### Benchmark
- Comparação de performance entre benchmarks
- Métricas históricas
- Relatórios automáticos para stakeholder
- Grafana dashboards de integração

## Relatórios e Métricas

### Formatos de Saída

**Console human-readable:**
```bash
./test-environment.sh --human-readable
```

**JSON para automação:**
```bash
./test-environment.sh --json > results.json
./performance-test.sh --json-report report.json
```

**Formato de métricas Prometheus:**
```bash
./performance-test.sh --prometheus-metrics
```

### Dados Coletados

Conjunto de métricas:

- Sistema: OS, arch, RAM, disk, CPU modelo
- Rede: Latency, throughput, DNS resolve
- Performance: Operations per second, memory footprint
- Configuração: Validações YAML, sintaxe config
- Segurança: Accessibilidade, permissões
- Integração: API endpoints, service discovery

## Integração com CI/CD

### GitHub Actions
```yaml
# .github/workflows/setup-validation.yml
name: Setup Validation
on: [push] 
jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Run Setup Validation 
        run: |
          cd examples/validation-tests
          ./test-environment.sh --ci-mode
```

### Jenkins Pipeline  
```groovy
pipeline {
  agent any
  stages {
    stage('Validate)' {
      steps {
        sh './examples/validation-tests/test-environment.sh --ci-mode'
        sh './examples/validation-tests/performance-test.sh --automated'
      }
    }
  }
}
```

### GitLab CI
```yaml
validate_setup:
  stage: test
  script:
    - cd examples/validation-tests
    - ./test-environment.sh --ci-mode > validation.json
    - ./performance-test.sh --quiet --exit-on-error
  artifacts:
    reports:
      performance: "${CI_PROJECT_DIR}/validation.json"
```

## Configuração e Customização

### Arquivos de Configuração

Crie arquivo `validation-config.yaml` para configurar validações:

```yaml
validation:
  timeout: 60
  scale:
    stress_test_level: high
    performance_ceiling: cpu_usage_80_percent
    
  environment:
    allowed_platforms: ['linux', 'windows']
    security_level: enterprise
    
  reporting:
    formats: ['json', 'html']
    output_directory: ./test-reports
```

### Variáveis de Ambiente

Configure através de variáveis específicas:

```bash
export VALIDATION_STRICT_MODE=true
export VALIDATION_TIMEOUT=300  
export VALIDATION_OUTPUT_FORMAT=json+yaml
export PERF_STRESS_LOAD_SECONDS=600
```

## Troubleshooting

### Problemas Comuns

**Permissions insufficient:**
```bash
# Run with sudo when required for resource checks
sudo ./test-environment.sh
```

**Network connectivity issues:**
```bash
# Diagnose network issues before running full test
./diagnose-network.sh
```

**Resources constraints:**
```bash
# Adjust test parameters for resource-limited environments
./performance-test.sh --low-resource --adaptive-load
```

### Debugging

**Ver detalhes da validação:**
```bash
./test-environment.sh --verbose --debug-log > debug.log
```

**Status e meterials checking:**
```bash
./test-environment.sh --quiet --check-only-resources
```

## Best Practices

1. **Execução em Pre-Builds:** Execute validação de ambiente antes de setup
2. **Testes Incrementais:** Use flag --skip-existing em re-testing
3. **Métricas de Monitoring:** Integre dados numa casa de observabilidade
4. **CI/CD Integration:** Sempre compare benchmarks consistency 
5. **Production Safety:** Nunca execute stress tests em produção

## Próximos Passos

1. **Material Support:** Examine `../../GUIDE.md` and `../../COMPILACAO_E_TESTE.md`
2. **Integration:** Adicione validações para `../../scripts/`
3. **Observability:** Integre com `Prometheus` e `Grafana`
4. **Custom Trim:** Customize para suas necessidades de infraestrutura
5. **Reporting:** Implement dashboards relevant para sua organização
