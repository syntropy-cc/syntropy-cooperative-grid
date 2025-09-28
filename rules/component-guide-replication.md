# 📖 Regras para Replicação de Guias de Componentes

> **Regras técnicas para LLMs criarem guias de componentes com estrutura padronizada e qualidade consistente**

## 📋 **Visão Geral**

Esta regra define o padrão para criação de guias de desenvolvimento (GUIDE.md) para componentes do projeto Syntropy Cooperative Grid, baseado na estrutura otimizada do guia de setup como referência de qualidade e implementação.

## 🎯 **Objetivo**

Criar guias consistentes, abrangentes e otimizados para LLMs que permitam implementação completa de componentes por agentes de construção de código, seguindo a mesma estrutura, qualidade e profundidade técnica do guia de setup.

## 🏗️ **Estrutura Hierárquica do Projeto**

### **Hierarquia de Organização:**
```
Projeto → Módulos → Sub-módulos → Componentes
```

### **Exemplo Prático:**
```
Syntropy Cooperative Grid (Projeto)
└── Syntropy Manager (Módulo)
    └── CLI (Sub-módulo)
        └── Setup (Componente)
```

### **Definição de Componente:**
- **Menor unidade entregável** do projeto
- **Resolve um problema específico** e bem definido
- **Funcionalidade independente** e testável
- **Interface clara** com outros componentes

## 📁 **Estrutura Padrão de Componente**

### **Estrutura de Diretórios Obrigatória:**
```
componente/
├── config/          # Configurações e templates
├── docs/            # Documentação adicional
├── examples/        # Exemplos de uso
├── scripts/         # Scripts de automação
├── src/             # Código fonte principal
└── test/            # Testes e validação
```

### **Diretório `src/` - Código Fonte Principal:**
- **Localização**: Pasta principal do código fonte
- **Organização**: Arquivos Go organizados por funcionalidade
- **Padrão**: Interfaces, implementações e utilitários

## 📐 **Estrutura Padrão Obrigatória do Guia**

### **1. Contexto e Objetivos (Obrigatório)**
```markdown
## Contexto e Objetivos

### Syntropy Cooperative Grid
[Breve explicação do projeto principal e seu propósito]

### Syntropy Manager
[Contexto do manager no ecossistema e suas responsabilidades]

### [Nome do Componente]
[Contextualização específica do componente, seu papel e objetivos]
```

### **2. Princípios de Implementação (Obrigatório)**
```markdown
## Princípios de Implementação

- **Simplicidade**: Arquitetura simples e direta, evitando over-engineering
- **Multiplataforma**: Suporte a Windows, Linux e macOS usando interfaces Go
- **Thread-Safe**: Operações atômicas e controle de concorrência
- **Segurança**: [Aspectos de segurança específicos do componente]
- **Observabilidade**: Logging estruturado e métricas
- **Testabilidade**: Componentes desacoplados e testáveis
- **Manutenibilidade**: Código limpo e bem documentado
```

### **3. Arquitetura Simplificada (Obrigatório)**
```markdown
## Arquitetura Simplificada

[Diagrama ASCII da arquitetura em 3 níveis]
┌─────────────────────────────────────────────────────────────┐
│ [Componente] (Orquestrador Principal)                      │
│ ─────────────────────────────────────────────────────────── │
│ • arquivo1.go    • arquivo2.go    • arquivo3.go            │
│ • manager.go     • validator.go   • configurator.go        │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│ Serviços Internos (Implementação por Interface)           │
│ ─────────────────────────────────────────────────────────── │
│ • Service1       • Service2       • Service3               │
│ • Provider1      • Provider2      • Provider3              │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│ Sistema de Estado Local (Persistência)                     │
│ ─────────────────────────────────────────────────────────── │
│ • config/        • data/          • cache/                 │
│ • logs/          • backups/       • temp/                  │
└─────────────────────────────────────────────────────────────┘
```

### **4. Explicação da Arquitetura por Níveis (Obrigatório)**
```markdown
## Explicação da Arquitetura por Níveis

### Nível 1: Componentes Principais (Orquestração)
[Explicação de cada arquivo principal com função, responsabilidades e interface]

### Nível 2: Serviços Internos (Implementação)
[Explicação de cada serviço com função, responsabilidades e implementações]

### Nível 3: Sistema de Estado Local (Persistência)
[Explicação de cada diretório de estado com função, conteúdo e segurança]

### Fluxo de Interação entre Níveis
[Diagrama de fluxo mostrando como os níveis interagem]

### Princípios de Design
[Princípios fundamentais que guiam a arquitetura]
```

### **5. Estrutura de Projeto Otimizada (Obrigatório)**
```markdown
## Estrutura de Projeto Otimizada

[Estrutura detalhada de diretórios com tamanhos de arquivo]
componente/
├── arquivo1.go                 # Descrição (X-Y linhas)
├── arquivo2.go                 # Descrição (X-Y linhas)
├── internal/
│   ├── types/
│   │   ├── tipos.go           # Estruturas de dados
│   │   └── interfaces.go      # Interfaces para desacoplamento
│   ├── services/
│   │   ├── service1/          # Serviço específico
│   │   └── service2/          # Outro serviço
│   └── utils/
│       ├── util1.go           # Utilitários específicos
│       └── util2.go           # Outros utilitários
├── config/
│   ├── templates/             # Templates de configuração
│   └── schemas/               # Schemas de validação
└── tests/
    ├── unit/                  # Testes unitários
    └── integration/           # Testes de integração
```

### **6. Interfaces e Contratos (Obrigatório)**
```markdown
## Interfaces e Contratos

### Interface Principal do [Componente]
[Interface principal com métodos e estruturas de dados]

### Interface de [Funcionalidade1]
[Interface específica com métodos e responsabilidades]

### Interface de [Funcionalidade2]
[Interface específica com métodos e responsabilidades]

[Continuar para todas as interfaces principais]
```

### **7. Implementação por Sistema Operacional (Quando Aplicável)**
```markdown
## Implementação por Sistema Operacional

### Interface de [Funcionalidade] por SO
[Interface com implementações específicas por SO]

### Implementações Específicas
[Implementações para Windows, Linux e macOS]
```

### **8. Sistema de [Funcionalidade] (Quando Aplicável)**
```markdown
## Sistema de [Funcionalidade]

### Interface de [Funcionalidade]
[Interface específica com métodos e estruturas]

### Implementação
[Exemplo de implementação com código Go]
```

### **9. Sistema de Erros Estruturado (Obrigatório)**
```markdown
## Sistema de Erros Estruturado

### Códigos de Erro e Contexto
[Estrutura de erro com códigos específicos]

### Códigos de Erro Específicos
[Constantes com códigos de erro do componente]
```

### **10. Comandos [Interface] (Quando Aplicável)**
```markdown
## Comandos [Interface]

### Estrutura de Comandos
[Comandos específicos do componente]

### Exemplo de Uso
[Exemplos práticos de uso]
```

### **11. Fluxo de Implementação (Obrigatório)**
```markdown
## Fluxo de Implementação

### 1. Implementação do [Componente Principal]
[Código de exemplo da implementação principal]

### 2. Implementação do [Subcomponente1]
[Código de exemplo da implementação]

### 3. Implementação do [Subcomponente2]
[Código de exemplo da implementação]

[Continuar para todos os componentes principais]
```

### **12. Sistema de [Funcionalidade Específica] (Quando Aplicável)**
```markdown
## Sistema de [Funcionalidade Específica]

### [Funcionalidade1]
[Código de exemplo da funcionalidade]

### [Funcionalidade2]
[Código de exemplo da funcionalidade]
```

### **13. Testes e Validação (Obrigatório)**
```markdown
## Testes e Validação

### Estrutura de Testes
[Estrutura de testes com exemplos]

### Exemplo de Teste
[Código de exemplo de teste]
```

### **14. Configurações e Templates (Quando Aplicável)**
```markdown
## Configurações e Templates

### Template de [Configuração]
[Template com variáveis]

### Schema de Validação
[Schema JSON para validação]
```

### **15. Considerações de Implementação (Obrigatório)**
```markdown
## Considerações de Implementação

### 1. [Aspecto1]
[Considerações específicas]

### 2. [Aspecto2]
[Considerações específicas]

### 3. [Aspecto3]
[Considerações específicas]

[Continuar para todos os aspectos importantes]
```

### **16. Conclusão (Obrigatório)**
```markdown
## Conclusão

Este guia fornece uma arquitetura [características] para o componente [nome], focando em:

1. **[Característica1]**: [Descrição]
2. **[Característica2]**: [Descrição]
3. **[Característica3]**: [Descrição]
4. **[Característica4]**: [Descrição]
5. **[Característica5]**: [Descrição]

A implementação deve seguir os padrões Go estabelecidos e garantir que o componente seja [objetivos].
```

## 🔧 **Regras Específicas para LLMs**

### **Ao Gerar Guias:**

#### **1. Estrutura Obrigatória**
- **SEMPRE** incluir todas as 16 seções padrão
- **SEMPRE** manter a ordem das seções
- **SEMPRE** usar emojis para identificação visual
- **SEMPRE** incluir diagramas ASCII quando aplicável

#### **2. Conteúdo Técnico**
- **SEMPRE** incluir interfaces Go completas
- **SEMPRE** especificar tamanho de arquivo (X-Y linhas)
- **SEMPRE** incluir código de exemplo funcional
- **SEMPRE** detalhar implementações específicas

#### **3. Qualidade e Consistência**
- **SEMPRE** usar linguagem técnica precisa
- **SEMPRE** incluir exemplos práticos de código
- **SEMPRE** manter consistência com guia de setup
- **SEMPRE** otimizar para implementação por LLM

#### **4. Especificidade do Componente**
- **ADAPTAR** conteúdo para o componente específico
- **MANTER** estrutura padrão
- **INCLUIR** funcionalidades específicas
- **DETALHAR** casos de uso específicos

### **Padrões de Escrita:**

#### **1. Linguagem**
- **Técnica e precisa**
- **Otimizada para implementação por LLM**
- **Consistente** com terminologia do projeto
- **Clara** e objetiva

#### **2. Formatação**
- **Markdown** bem estruturado
- **Código Go** em blocos apropriados
- **Interfaces** bem definidas
- **Seções** bem delimitadas

#### **3. Exemplos**
- **Código** executável e funcional
- **Interfaces** completas
- **Implementações** práticas
- **Testes** funcionais

## 📊 **Template de Referência**

### **Estrutura Base:**
```markdown
# [Nome do Componente] - Guia de Implementação para LLMs

## Contexto e Objetivos
### Syntropy Cooperative Grid
### Syntropy Manager
### [Nome do Componente]

## Princípios de Implementação
[Princípios específicos do componente]

## Arquitetura Simplificada
[Diagrama ASCII + explicação]

## Explicação da Arquitetura por Níveis
[Explicação detalhada dos 3 níveis]

## Estrutura de Projeto Otimizada
[Estrutura de diretórios com tamanhos]

## Interfaces e Contratos
[Interfaces Go completas]

## Implementação por Sistema Operacional
[Quando aplicável]

## Sistema de [Funcionalidade]
[Funcionalidades específicas]

## Sistema de Erros Estruturado
[Códigos de erro específicos]

## Comandos [Interface]
[Quando aplicável]

## Fluxo de Implementação
[Código de exemplo completo]

## Sistema de [Funcionalidade Específica]
[Funcionalidades adicionais]

## Testes e Validação
[Estrutura de testes]

## Configurações e Templates
[Quando aplicável]

## Considerações de Implementação
[Aspectos importantes]

## Conclusão
[Resumo e objetivos]
```

## 🎯 **Critérios de Qualidade**

### **Obrigatórios:**
1. **Estrutura completa** (16 seções)
2. **Interfaces Go** completas e funcionais
3. **Código de exemplo** executável
4. **Arquitetura em 3 níveis** bem explicada
5. **Tamanho de arquivo** especificado
6. **Sistema de erros** estruturado
7. **Considerações** de implementação

### **Desejáveis:**
1. **Diagramas ASCII** claros
2. **Código funcional** testado
3. **Casos de uso** realistas
4. **Troubleshooting** incluído
5. **Performance** considerada

## 🚨 **Regras Críticas**

### **NUNCA:**
1. **Omitir** seções obrigatórias
2. **Criar** estrutura diferente do padrão
3. **Esquecer** interfaces Go
4. **Ignorar** aspectos de implementação
5. **Usar** linguagem não técnica

### **SEMPRE:**
1. **Seguir** estrutura padrão
2. **Incluir** todas as seções
3. **Fornecer** código de exemplo
4. **Especificar** tamanho de arquivo
5. **Otimizar** para implementação por LLM

## 📚 **Recursos de Referência**

### **Guia de Referência:**
- `manager/interfaces/cli/setup/GUIDE.md` - Padrão de qualidade

### **Estrutura de Projeto:**
- `manager/` - Módulo principal
- `interfaces/` - Interfaces do sistema
- `core/` - Lógica de negócio

### **Padrões Técnicos:**
- Interfaces Go: `type Interface interface { ... }`
- Tamanho de arquivo: 200-500 linhas
- Nomenclatura: `componente_funcionalidade.go`

## 🎯 **Objetivo Final**

**Criar guias que permitam a LLMs implementarem componentes completos e funcionais, seguindo a mesma estrutura, qualidade e profundidade técnica do guia de setup, garantindo consistência e excelência em todo o projeto.**

---

**Esta regra garante que todos os guias de componentes gerados mantenham o mesmo padrão de qualidade, estrutura e capacidade de implementação do guia de referência do setup!** 🚀
