# 📋 Regras para LLMs - Syntropy Cooperative Grid

> **Diretório de regras técnicas e sucintas para auxiliar LLMs a executar tarefas no projeto**

## 🎯 **Objetivo**

Este diretório contém regras otimizadas para LLMs trabalharem com o projeto Syntropy Cooperative Grid. As regras são:

- **Sucintas**: Informações essenciais sem verbosidade
- **Técnicas**: Focadas em implementação e uso prático
- **Otimizadas para LLMs**: Estruturadas para fácil parsing e compreensão
- **Específicas**: Direcionadas para componentes específicos do projeto

## 📁 **Estrutura**

```
rules/
├── README.md                    # Este arquivo - visão geral
├── management-system.md         # Regras para o Management System
├── guide-generation.md          # Regras para geração de guias de desenvolvimento
├── development.md              # Regras para desenvolvimento (futuro)
├── deployment.md               # Regras para deployment (futuro)
├── testing.md                  # Regras para testes (futuro)
└── security.md                 # Regras para segurança (futuro)
```

## 🔧 **Componentes Disponíveis**

### **Management System** (`management-system.md`)
Sistema unificado para gerenciar a Syntropy Cooperative Grid:

### **Guide Generation** (`guide-generation.md`)
Regras para geração de guias de desenvolvimento equivalentes ao padrão estabelecido:

**Funcionalidades principais:**
- Estrutura padrão obrigatória (14 seções)
- Integração com API central (`manager/api/`)
- Boas práticas de desenvolvimento
- Padrões de segurança e performance
- Exemplos práticos e casos de uso
- Otimização para leitura por LLMs

**Regras críticas:**
- SEMPRE incluir todas as seções obrigatórias
- SEMPRE referenciar integração com API central
- SEMPRE especificar tamanho de arquivo (300-500 linhas)
- SEMPRE incluir seção de segurança detalhada
- NUNCA omitir estrutura padrão

**Estrutura de Interfaces:**
- **CLI**: `interfaces/cli/` - Interface de linha de comando (Go + Cobra)
- **Web**: `interfaces/web/` - Dashboard web (React + Next.js + Go backend)
- **Mobile**: `interfaces/mobile/flutter/` - App mobile (Flutter)
- **Desktop**: `interfaces/desktop/electron/` - App desktop (Electron)

**Funcionalidades principais:**
- Gerenciamento de nós (detecção USB, configuração, monitoramento)
- Gerenciamento de containers (deploy, orquestração, escalabilidade)
- Gerenciamento de rede (service mesh, roteamento, conectividade)
- Gerenciamento cooperativo (créditos, governança, reputação)

**Regras críticas de estrutura:**
- NUNCA criar arquivos fora dos diretórios corretos
- Cada interface tem sua estrutura específica bem definida
- Core logic deve ir em `core/`, código compartilhado em `internal/`

## 📖 **Como Usar**

### **Para LLMs:**
1. **Leia as regras relevantes** antes de executar tarefas
2. **Siga os padrões estabelecidos** para comandos e configurações
3. **Use os exemplos fornecidos** como base para implementações
4. **Valide inputs** e implemente error handling adequado
5. **Documente mudanças** importantes
6. **RESPEITE a estrutura de diretórios** - NUNCA crie arquivos em locais incorretos
7. **Use os caminhos corretos** para cada tipo de interface

### **Para Desenvolvedores:**
1. **Consulte as regras** ao trabalhar com componentes específicos
2. **Mantenha as regras atualizadas** quando houver mudanças
3. **Adicione novas regras** para novos componentes
4. **Use as regras como referência** para documentação

## 🎯 **Padrões das Regras**

### **Estrutura Padrão:**
```markdown
# 🎯 Título do Componente

> **Descrição sucinta do componente**

## 📋 **Visão Geral**
- Objetivo e propósito
- Contexto no projeto

## 🏗️ **Arquitetura**
- Estrutura de diretórios do projeto
- Componentes principais
- Relacionamentos
- Stack tecnológico

## 🔧 **Funcionalidades**
- Lista de funcionalidades principais
- Casos de uso
- Comandos essenciais

## 📝 **Regras para LLMs**
- Estrutura de diretórios - regras críticas
- Instruções específicas
- Padrões a seguir
- Exemplos práticos
```

### **Características:**
- **Emojis**: Para identificação visual rápida
- **Código**: Exemplos práticos e comandos
- **Estrutura clara**: Seções bem definidas
- **Informações técnicas**: Detalhes de implementação
- **Casos de uso**: Exemplos reais de uso

## 🚀 **Roadmap de Regras**

### **Implementadas:**
- ✅ [Management System](management-system.md) - Sistema de gerenciamento completo
- ✅ [Guide Generation](guide-generation.md) - Regras para geração de guias

### **Planejadas:**
- ⏳ [Development](development.md) - Regras para desenvolvimento
- ⏳ [Deployment](deployment.md) - Regras para deployment
- ⏳ [Testing](testing.md) - Regras para testes
- ⏳ [Security](security.md) - Regras para segurança
- ⏳ [API](api.md) - Regras para APIs
- ⏳ [Database](database.md) - Regras para banco de dados
- ⏳ [Monitoring](monitoring.md) - Regras para monitoramento

## 📝 **Contribuindo**

### **Ao adicionar novas regras:**
1. **Siga a estrutura padrão** definida acima
2. **Mantenha o foco técnico** e sucinto
3. **Inclua exemplos práticos** sempre que possível
4. **Use emojis** para identificação visual
5. **Atualize este README** com a nova regra

### **Ao atualizar regras existentes:**
1. **Mantenha a compatibilidade** com versões anteriores
2. **Documente mudanças** importantes
3. **Atualize exemplos** se necessário
4. **Valide comandos** e configurações

## 🔍 **Exemplos de Uso**

### **Para LLMs trabalhando com Management System:**
```bash
# Sempre validar antes de executar
syntropy node list --format json | jq '.[] | select(.status == "unhealthy")'

# Usar templates para consistência
syntropy container deploy --template nginx --node node-01

# Monitorar operações
syntropy node status node-01 --watch --format table
```

### **Para desenvolvimento:**
```go
// Seguir padrões estabelecidos
func (s *NodeService) CreateNode(req *CreateNodeRequest) (*Node, error) {
    // Validação
    if err := validateNodeRequest(req); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }
    
    // Implementação
    // ...
}
```

## 📚 **Recursos Relacionados**

- [Documentação Principal](../docs/) - Documentação completa do projeto
- [Management System Docs](../docs/management-system/) - Documentação detalhada
- [Development Guide](../docs/development/) - Guia de desenvolvimento
- [API Reference](../docs/api/) - Referência das APIs

## 🎯 **Objetivos**

1. **Facilitar trabalho de LLMs** com informações técnicas precisas
2. **Padronizar abordagens** para tarefas comuns
3. **Reduzir tempo de onboarding** para novos componentes
4. **Manter consistência** entre diferentes partes do projeto
5. **Fornecer referência rápida** para comandos e configurações
6. **Prevenir criação de arquivos em locais incorretos**

## 📁 **Estrutura de Diretórios - Resumo**

### **Interfaces (`interfaces/`)**
- `interfaces/cli/` - CLI Interface (Go + Cobra)
- `interfaces/web/frontend/` - Web Frontend (React + Next.js)
- `interfaces/web/backend/` - Web Backend (Go)
- `interfaces/mobile/flutter/` - Mobile App (Flutter)
- `interfaces/desktop/electron/` - Desktop App (Electron)

### **Core e Serviços**
- `core/` - Management Core (lógica de negócio)
- `internal/` - Código interno compartilhado
- `services/` - Microserviços

### **Documentação e Regras**
- `docs/` - Documentação completa
- `rules/` - Regras para LLMs (este diretório)

### **Regra Crítica:**
**NUNCA crie arquivos fora dos diretórios corretos. Cada interface tem sua estrutura específica bem definida.**

---

**Este diretório é um recurso vivo que evolui com o projeto. Mantenha as regras atualizadas e relevantes!** 🚀
