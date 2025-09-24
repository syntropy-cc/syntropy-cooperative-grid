# TODO - Setup Component Implementation (Simplified)

## 📋 **Visão Geral**
Lista de tarefas para implementação do Setup Component simplificado - o quartel geral para criação e gestão de nós da rede Syntropy.

## 🎯 **Objetivo**
Implementar componente Setup com 2 subcomponentes (Validation + Configuration) com foco no Windows, seguindo arquitetura simplificada.

---

## 📁 **FASE 1: Estrutura Base (1 dia)**

### ✅ **1.1 Estrutura de Diretórios**
- [ ] Criar diretório `internal/` com subdiretórios:
  - [ ] `internal/types/` - Tipos e estruturas de dados
  - [ ] `internal/services/` - Serviços internos
- [ ] Criar diretório `config/` com subdiretórios:
  - [ ] `config/templates/` - Templates de configuração
- [ ] Criar diretório `tests/` com subdiretórios:
  - [ ] `tests/unit/` - Testes unitários
  - [ ] `tests/integration/` - Testes de integração

### ✅ **1.2 Orquestrador Principal**
- [ ] Implementar `setup.go` (300-500 linhas)
  - [ ] Função principal `Setup()`
  - [ ] Switch por sistema operacional
  - [ ] Orquestração de subcomponentes
  - [ ] Tratamento de erros centralizado

---

## 🔧 **FASE 2: Subcomponentes Windows (6 dias total)**

### ✅ **2.1 Validation Subcomponent (3 dias)**
- [ ] **Dia 1**: Estrutura básica
  - [ ] Implementar `validation_windows.go` (300-500 linhas)
  - [ ] Criar `internal/types/validation.go`
  - [ ] Criar `tests/unit/validation_test.go`
- [ ] **Dia 2**: Detecção completa
  - [ ] Detecção de SO (Windows, versão, arquitetura)
  - [ ] Verificação de permissões administrativas
  - [ ] Verificação de espaço em disco (mínimo 1GB)
  - [ ] Verificação de PowerShell (versão 5.1+)
  - [ ] Verificação de conectividade de rede
  - [ ] Testes de detecção
- [ ] **Dia 3**: Integração
  - [ ] Integração com API central
  - [ ] Sistema de logs estruturado
  - [ ] Testes de integração
  - [ ] Documentação

### ✅ **2.2 Configuration Subcomponent (3 dias)**
- [ ] **Dia 1**: Estrutura básica
  - [ ] Implementar `configuration_windows.go` (300-500 linhas)
  - [ ] Criar `internal/types/config.go`
  - [ ] Criar `tests/unit/configuration_test.go`
- [ ] **Dia 2**: Geração de configuração
  - [ ] Geração de `manager.yaml`
  - [ ] Criação de estrutura `~/.syntropy/`
  - [ ] Geração de owner key (Ed25519)
  - [ ] Validação de configuração
  - [ ] Testes de configuração
- [ ] **Dia 3**: Integração
  - [ ] Integração com API central
  - [ ] Sistema de logs estruturado
  - [ ] Testes de integração
  - [ ] Documentação

### ✅ **2.3 Setup Windows Orchestrator**
- [ ] Implementar `setup_windows.go` (300-500 linhas)
  - [ ] Orquestração específica para Windows
  - [ ] Sequência: Validation → Configuration
  - [ ] Sistema de rollback em caso de erro
  - [ ] Validação final do setup
  - [ ] Geração de relatório de setup

---

## 🏗️ **FASE 3: Serviços e Integração (2 dias)**

### ✅ **3.1 Tipos Internos**
- [ ] Implementar `internal/types/validation.go`
  - [ ] Estrutura `ValidationResult`
  - [ ] Estrutura `EnvironmentInfo`
  - [ ] Estrutura `SystemResources`
- [ ] Implementar `internal/types/config.go`
  - [ ] Estrutura `SetupConfig`
  - [ ] Estrutura `ManagerConfig`
  - [ ] Estrutura `OwnerKey`

### ✅ **3.2 Templates de Configuração**
- [ ] Criar `config/templates/manager.yaml`
  - [ ] Template de configuração principal
  - [ ] Configurações mínimas necessárias

### ✅ **3.3 Integração com API Central**
- [ ] Integrar com `manager/api/handlers/config/`
  - [ ] Reutilizar lógica de configuração
  - [ ] Compartilhar tipos de dados
- [ ] Integrar com `manager/api/services/validation/`
  - [ ] Reutilizar serviços de validação
  - [ ] Compartilhar lógica de validação

---

## 🧪 **FASE 4: Testes e Qualidade (2 dias)**

### ✅ **4.1 Testes Unitários**
- [ ] Implementar `tests/unit/validation_test.go`
  - [ ] Testes de detecção de ambiente
  - [ ] Testes de validação de recursos
  - [ ] Testes de permissões
- [ ] Implementar `tests/unit/configuration_test.go`
  - [ ] Testes de geração de configuração
  - [ ] Testes de validação de configuração
  - [ ] Testes de geração de chaves
- [ ] Implementar `tests/unit/setup_test.go`
  - [ ] Testes de orquestração
  - [ ] Testes de rollback
  - [ ] Testes de validação final

### ✅ **4.2 Testes de Integração**
- [ ] Implementar `tests/integration/setup_integration_test.go`
  - [ ] Teste completo de setup
  - [ ] Teste de integração com API
  - [ ] Teste de cenários de erro

---

## 🖥️ **FASE 5: Interface e Documentação (1 dia)**

### ✅ **5.1 Comandos CLI**
- [ ] Integrar com comando `syntropy setup`
  - [ ] Setup completo (valida + configura)
  - [ ] Validação apenas (`--validate-only`)
  - [ ] Forçar setup (`--force`)
  - [ ] Status do setup (`status`)
  - [ ] Reset completo (`reset`)

### ✅ **5.2 Documentação do Usuário**
- [ ] Criar `README.md`
  - [ ] Visão geral do Setup Component
  - [ ] Comandos disponíveis
  - [ ] Exemplos de uso
  - [ ] Troubleshooting
  - [ ] FAQ

---

## 📊 **Critérios de Sucesso**

### ✅ **Funcionalidade**
- [ ] Usuário pode executar `syntropy setup` com sucesso no Windows
- [ ] Ambiente é detectado e validado automaticamente
- [ ] Configuração é gerada e validada
- [ ] Owner key é gerada e armazenada com segurança
- [ ] Sistema funciona offline após setup

### ✅ **Qualidade**
- [ ] Cobertura de testes >= 80%
- [ ] Todos os testes passando
- [ ] Linting sem erros
- [ ] Documentação completa e atualizada

### ✅ **Integração**
- [ ] Integração funcional com API central
- [ ] Reutilização de componentes existentes
- [ ] Consistência com padrões do projeto
- [ ] Comandos CLI funcionais

---

## 🚀 **Ordem de Implementação Recomendada**

1. **Estrutura Base** → Criar diretórios e orquestrador principal (1 dia)
2. **Validation Subcomponent** → Implementar validação completa (3 dias)
3. **Configuration Subcomponent** → Implementar configuração (3 dias)
4. **Setup Windows** → Implementar orquestração específica (1 dia)
5. **Tipos e Serviços** → Implementar tipos internos e integração (2 dias)
6. **Testes** → Implementar testes unitários e de integração (2 dias)
7. **CLI e Documentação** → Integrar comandos e criar documentação (1 dia)

**Total**: 13 dias para implementação completa

---

## 📝 **Notas de Implementação**

- **Tamanho de arquivo**: Cada arquivo deve ter entre 300-500 linhas
- **Build tags**: Usar `//go:build windows`
- **Segurança**: Implementar criptografia Ed25519 para owner key
- **Logs**: Usar logging estruturado (logrus)
- **Erros**: Implementar tratamento de erros com contexto
- **Validação**: Validar todas as entradas e configurações
- **Performance**: Otimizar para operações I/O bound

---

## 🎯 **Comandos Simplificados**

```bash
# Setup completo
syntropy setup

# Só validar
syntropy setup --validate-only

# Forçar setup
syntropy setup --force

# Status
syntropy setup status

# Reset
syntropy setup reset

# Configuração
syntropy setup config generate
syntropy setup config validate
syntropy setup config backup
```

---

**Status**: 🚧 Em desenvolvimento  
**Prioridade**: 🔥 Alta (Fase 1 do projeto CLI)  
**Responsável**: Equipe de desenvolvimento  
**Prazo**: 13 dias para implementação completa