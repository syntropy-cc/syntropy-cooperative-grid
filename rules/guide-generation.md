# 📖 Regras para Geração de Guias de Desenvolvimento

> **Regras técnicas para LLMs gerarem guias de desenvolvimento equivalentes ao padrão estabelecido**

## 📋 **Visão Geral**

Esta regra define o padrão para criação de guias de desenvolvimento (GUIDE.md) para componentes do projeto Syntropy Cooperative Grid, baseado no guia da CLI como referência de qualidade e estrutura.

## 🎯 **Objetivo**

Criar guias consistentes, abrangentes e otimizados para LLMs que sigam o mesmo padrão de qualidade, estrutura e profundidade técnica do guia da CLI.

## 📐 **Estrutura Padrão Obrigatória**

### **1. Contexto e Objetivos (Obrigatório)**
```markdown
## Contexto e Objetivos

### Syntropy Cooperative Grid
[Breve explicação do projeto principal]

### Syntropy Manager
[Contexto do manager no ecossistema]

### [Nome do Componente]
[Contextualização específica do componente]
```

### **2. Princípios Fundamentais (Obrigatório)**
- **Desenvolvimento Baseado em Componentes**: Sempre mencionar
- **Multiplataforma**: Tags `//go:build` quando aplicável
- **Integração com API**: Referência a `manager/api/`
- **Padrões específicos** do componente

### **3. Arquitetura do Componente (Obrigatório)**
- **Diagrama ASCII** da arquitetura
- **Estrutura de diretórios** detalhada
- **Integração com API central**
- **Integração com rede existente**

### **4. Hierarquia de Implementação (Obrigatório)**
```markdown
## Hierarquia de Implementação Baseada em Componentes

### Macro Etapa: [Nome do Componente]
### Meso Etapas (Componentes)
### Micro Etapas (Subcomponentes)
### Foco de Implementação
### Implementação por Sistema Operacional
```

### **5. Comandos por Componente (Quando Aplicável)**
- **Comandos específicos** do componente
- **Exemplos práticos** de uso
- **Parâmetros e flags** importantes

### **6. Sistema de Estado (Quando Aplicável)**
- **Estrutura de dados** local
- **Gerenciamento de estado**
- **Persistência** e backup

### **7. Tecnologias e Padrões (Obrigatório)**
- **Stack tecnológico** específico
- **Padrões de desenvolvimento**
- **Integração** com outros componentes
- **Infraestrutura** necessária

### **8. Boas Práticas de Desenvolvimento (Obrigatório)**
- **Princípios fundamentais** (SOLID, Clean Code, Design Patterns)
- **Arquitetura e estrutura**
- **Qualidade de código**
- **Performance e otimização**
- **Segurança e compliance**
- **DevOps e CI/CD**
- **Monitoring e observabilidade**
- **Code organization**
- **Code quality tools**
- **Documentation and knowledge management**

### **9. Considerações Técnicas (Obrigatório)**
- **Segurança** (detalhada com criptografia quantum-resistante)
- **Performance**
- **Usabilidade**
- **Extensibilidade**

### **10. Processo de Desenvolvimento (Obrigatório)**
- **Etapas por componente**
- **Prioridades** de implementação
- **Critérios de sucesso**

### **11. Primeira Meso Etapa (Obrigatório)**
- **Objetivo** específico
- **Entregáveis** claros
- **Critérios de sucesso**
- **Micro etapas detalhadas** com tamanho de arquivo (300-500 linhas)

### **12. Exemplos de Uso (Obrigatório)**
- **Exemplos práticos** por componente
- **Casos de uso** reais
- **Comandos** específicos

### **13. Padrões de Nomenclatura (Obrigatório)**
- **Estrutura de arquivos** por componente
- **Build tags** específicas
- **Exemplo de orquestração** em Go

### **14. Documentação por Componente (Obrigatório)**
- **GUIDE.md** (guia de implementação)
- **README.md** (documentação do usuário)

## 🔧 **Regras Específicas para LLMs**

### **Ao Gerar Guias:**

#### **1. Estrutura Obrigatória**
- **SEMPRE** incluir todas as 14 seções padrão
- **SEMPRE** manter a ordem das seções
- **SEMPRE** usar emojis para identificação visual
- **SEMPRE** incluir diagramas ASCII quando aplicável

#### **2. Conteúdo Técnico**
- **SEMPRE** referenciar `manager/api/` para integração
- **SEMPRE** incluir tags `//go:build` para multiplataforma
- **SEMPRE** especificar tamanho de arquivo (300-500 linhas)
- **SEMPRE** incluir seção de segurança detalhada

#### **3. Qualidade e Consistência**
- **SEMPRE** usar linguagem técnica precisa
- **SEMPRE** incluir exemplos práticos
- **SEMPRE** manter consistência com guia da CLI
- **SEMPRE** otimizar para leitura por LLM

#### **4. Especificidade do Componente**
- **ADAPTAR** conteúdo para o componente específico
- **MANTER** estrutura padrão
- **INCLUIR** funcionalidades específicas
- **DETALHAR** casos de uso específicos

### **Padrões de Escrita:**

#### **1. Linguagem**
- **Técnica e precisa**
- **Otimizada para LLMs**
- **Consistente** com terminologia do projeto
- **Clara** e objetiva

#### **2. Formatação**
- **Markdown** bem estruturado
- **Código** em blocos apropriados
- **Listas** organizadas
- **Seções** bem delimitadas

#### **3. Exemplos**
- **Comandos** funcionais
- **Código** executável
- **Casos de uso** reais
- **Configurações** práticas

## 📊 **Template de Referência**

### **Estrutura Base:**
```markdown
# [Nome do Componente] - Guia de Desenvolvimento

## Contexto e Objetivos
### Syntropy Cooperative Grid
### Syntropy Manager  
### [Nome do Componente]

## Princípios Fundamentais
- Desenvolvimento Baseado em Componentes
- Multiplataforma
- Integração com API
- [Específicos do componente]

## Arquitetura do [Componente]
[Diagrama ASCII + estrutura de diretórios]

## Integração com API Central
**FUNDAMENTAL**: [Detalhes específicos]

## Integração com Rede Existente
[Componentes utilizados]

## Hierarquia de Implementação Baseada em Componentes
[Macro/Meso/Micro etapas]

## Comandos por Componente
[Comandos específicos]

## Sistema de Estado Local
[Quando aplicável]

## Tecnologias e Padrões
[Stack específico + padrões]

## Boas Práticas de Desenvolvimento
[Seção completa obrigatória]

## Considerações Técnicas
[Segurança detalhada + outros aspectos]

## Processo de Desenvolvimento por Componentes
[Etapas específicas]

## Primeira Meso Etapa: [Componente Principal]
[Objetivo + entregáveis + micro etapas]

## Exemplos de Uso por Componente
[Exemplos práticos]

## Padrões de Nomenclatura de Arquivos
[Estrutura + build tags + exemplo Go]

## Documentação por Componente
[GUIDE.md + README.md]

---

**Objetivo**: [Resumo do objetivo do componente]
```

## 🎯 **Critérios de Qualidade**

### **Obrigatórios:**
1. **Estrutura completa** (14 seções)
2. **Integração com API** mencionada
3. **Segurança detalhada** incluída
4. **Exemplos práticos** fornecidos
5. **Tamanho de arquivo** especificado
6. **Build tags** incluídas
7. **Boas práticas** completas

### **Desejáveis:**
1. **Diagramas ASCII** claros
2. **Comandos funcionais** testados
3. **Casos de uso** realistas
4. **Troubleshooting** incluído
5. **Performance** considerada

## 📝 **Exemplos de Aplicação**

### **Para Web Interface:**
- Adaptar comandos para interface web
- Incluir componentes React/Next.js
- Adicionar autenticação e autorização
- Considerar responsividade

### **Para Mobile Interface:**
- Adaptar para Flutter
- Considerar limitações mobile
- Incluir notificações push
- Adicionar offline capabilities

### **Para Desktop Interface:**
- Adaptar para Electron
- Considerar tray icon
- Incluir notificações do sistema
- Adicionar auto-updater

## 🚨 **Regras Críticas**

### **NUNCA:**
1. **Omitir** seções obrigatórias
2. **Criar** estrutura diferente do padrão
3. **Esquecer** integração com API central
4. **Ignorar** aspectos de segurança
5. **Usar** linguagem não técnica

### **SEMPRE:**
1. **Seguir** estrutura padrão
2. **Incluir** todas as seções
3. **Referenciar** `manager/api/`
4. **Especificar** tamanho de arquivo
5. **Otimizar** para LLMs

## 📚 **Recursos de Referência**

### **Guia de Referência:**
- `manager/interfaces/cli/GUIDE.md` - Padrão de qualidade

### **Estrutura de Projeto:**
- `manager/api/` - API central
- `interfaces/` - Interfaces do sistema
- `core/` - Lógica de negócio

### **Padrões Técnicos:**
- Build tags: `//go:build windows|linux|darwin`
- Tamanho de arquivo: 300-500 linhas
- Nomenclatura: `componente_sistemaoperacional.go`

---

**Esta regra garante que todos os guias gerados mantenham o mesmo padrão de qualidade, estrutura e profundidade técnica do guia de referência da CLI!** 🚀
