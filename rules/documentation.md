# 📚 Regras para Documentação

> **Regras técnicas e sucintas para LLMs criarem documentação padronizada para o Syntropy Cooperative Grid**

## 📋 **Visão Geral**

Esta regra define o padrão para criação de documentação técnica no projeto Syntropy Cooperative Grid, garantindo consistência, completude e otimização para leitura tanto por humanos quanto por LLMs.

## 🎯 **Objetivo**

Estabelecer um padrão de documentação unificado que maximize a compreensão técnica, facilite a manutenção e promova a colaboração eficiente entre desenvolvedores e LLMs.

## 🏗️ **Estrutura de Documentação**

### **Hierarquia de Documentação**
```
docs/
├── README.md                    # Visão geral da documentação
├── architecture/                # Documentação de arquitetura
├── api/                         # Documentação de APIs
├── development/                 # Guias de desenvolvimento
├── setup/                       # Instruções de configuração
├── management-system/           # Documentação do sistema de gerenciamento
└── roadmap/                     # Planejamento futuro
```

### **Tipos de Documentação**
- **README.md**: Introdução e visão geral de componentes
- **GUIDE.md**: Guias detalhados de desenvolvimento
- **API.md**: Documentação de APIs
- **SETUP.md**: Instruções de configuração
- **TEST.md**: Documentação de testes

## 📄 **Padrões por Tipo de Documento**

### **README.md (Obrigatório)**
- **Tamanho**: 100-300 linhas
- **Audiência**: Usuários e desenvolvedores
- **Conteúdo**:
  - Título e descrição sucinta
  - Badges (status, versão, licença)
  - Visão geral do componente
  - Pré-requisitos
  - Instalação rápida
  - Uso básico com exemplos
  - Estrutura de diretórios
  - Links para documentação adicional

### **GUIDE.md (Obrigatório para Componentes Principais)**
- **Tamanho**: 300-500 linhas
- **Audiência**: Desenvolvedores
- **Conteúdo**:
  - Contexto e objetivos
  - Princípios fundamentais
  - Arquitetura do componente
  - Hierarquia de implementação
  - Tecnologias e padrões
  - Boas práticas de desenvolvimento
  - Considerações técnicas
  - Processo de desenvolvimento
  - Exemplos de uso

### **API.md (Obrigatório para Serviços)**
- **Tamanho**: 200-400 linhas
- **Audiência**: Desenvolvedores e integradores
- **Conteúdo**:
  - Visão geral da API
  - Autenticação e autorização
  - Endpoints com métodos HTTP
  - Parâmetros de requisição
  - Formatos de resposta
  - Códigos de status
  - Exemplos de requisição e resposta
  - Rate limiting e quotas
  - Versionamento

### **TEST.md (Obrigatório para Componentes Testáveis)**
- **Tamanho**: 100-300 linhas
- **Audiência**: Desenvolvedores e QA
- **Conteúdo**:
  - Visão geral dos testes
  - Pré-requisitos para execução
  - Estrutura dos testes
  - Comandos para execução
  - Interpretação de resultados
  - Mocks e fixtures
  - Cobertura de código
  - Troubleshooting comum

## 🔍 **Elementos Obrigatórios**

### **1. Cabeçalho Padronizado**
```markdown
# 📚 [Nome do Componente] - [Tipo de Documento]

> **[Descrição sucinta em uma linha]**

![Status](https://img.shields.io/badge/status-[status]-[cor])
![Versão](https://img.shields.io/badge/versão-[x.y.z]-blue)
```

### **2. Seções Obrigatórias**
- **Visão Geral**: Contexto e propósito
- **Pré-requisitos**: Dependências e configurações necessárias
- **Instalação/Configuração**: Passo a passo detalhado
- **Uso**: Exemplos práticos com código
- **Estrutura**: Organização de arquivos/diretórios
- **Referências**: Links para documentação relacionada

### **3. Formatação de Código**
- **Blocos de código**: Sempre com identificação de linguagem
```markdown
```go
// Exemplo de código Go
func Example() string {
    return "Exemplo"
}
```
```

### **4. Tabelas para Dados Estruturados**
```markdown
| Campo | Tipo | Descrição | Obrigatório |
|-------|------|-----------|-------------|
| id    | string | Identificador único | Sim |
| nome  | string | Nome do recurso | Sim |
| config | object | Configurações | Não |
```

## 🛠️ **Boas Práticas**

### **Para Desenvolvedores**
- **Atualizar em conjunto com código**: Documentação deve ser atualizada junto com mudanças no código
- **Usar exemplos reais**: Exemplos devem ser testados e funcionais
- **Incluir troubleshooting**: Antecipar problemas comuns e suas soluções
- **Versionar documentação**: Manter histórico de mudanças importantes

### **Para LLMs**
- **Seguir estrutura exata**: Respeitar todas as seções obrigatórias
- **Manter consistência**: Usar mesma terminologia em todo documento
- **Otimizar para parsing**: Usar formatação consistente para facilitar extração
- **Balancear detalhes**: Fornecer informação suficiente sem verbosidade excessiva
- **Incluir metadados**: Adicionar tags e categorias para facilitar busca

## 📊 **Métricas de Qualidade**

### **Completude**
- **100%**: Todas seções obrigatórias presentes
- **75%**: Maioria das seções obrigatórias presentes
- **50%**: Apenas metade das seções obrigatórias
- **25%**: Documentação mínima ou incompleta

### **Clareza**
- **Linguagem técnica precisa**
- **Frases curtas e diretas**
- **Terminologia consistente**
- **Hierarquia visual clara**

### **Acionabilidade**
- **Instruções passo-a-passo**
- **Comandos completos e copiáveis**
- **Verificações de sucesso após cada etapa**
- **Troubleshooting para erros comuns**

## 🔄 **Processo de Documentação**

### **1. Planejamento**
- Identificar audiência e objetivos
- Selecionar tipo de documento apropriado
- Listar tópicos principais a cobrir

### **2. Desenvolvimento**
- Criar estrutura com todas seções obrigatórias
- Desenvolver conteúdo técnico preciso
- Incluir exemplos práticos e testados

### **3. Revisão**
- Verificar completude (todas seções obrigatórias)
- Validar precisão técnica
- Confirmar clareza e acionabilidade

### **4. Publicação**
- Adicionar ao repositório no local correto
- Atualizar links em documentos relacionados
- Comunicar atualização à equipe

## 📝 **Exemplos**

### **Exemplo de README.md**
```markdown
# 📚 Node Manager - README

> **Componente para gerenciamento de nós na Syntropy Cooperative Grid**

![Status](https://img.shields.io/badge/status-stable-green)
![Versão](https://img.shields.io/badge/versão-1.2.0-blue)

## Visão Geral
O Node Manager permite detectar, configurar e monitorar nós físicos e virtuais na rede Syntropy.

## Pré-requisitos
- Go 1.18+
- Acesso root/admin
- Dispositivos USB (para nós físicos)

## Instalação Rápida
```bash
git clone https://github.com/syntropy-cc/syntropy-cooperative-grid.git
cd syntropy-cooperative-grid
./scripts/setup.sh node-manager
```

## Uso Básico
```bash
syntropy node create --usb /dev/sdb --name "node-01"
syntropy node list --format table
```

## Estrutura
```
node-manager/
├── cmd/                # Comandos CLI
├── internal/           # Lógica interna
├── api/                # API REST
└── tests/              # Testes unitários
```

## Documentação Adicional
- [Guia de Desenvolvimento](./GUIDE.md)
- [API Reference](./API.md)
- [Troubleshooting](./TROUBLESHOOTING.md)
```

## ⚠️ **Regras Críticas**

- **SEMPRE** incluir todas as seções obrigatórias
- **SEMPRE** usar formatação consistente (Markdown)
- **SEMPRE** incluir exemplos práticos e testados
- **SEMPRE** especificar versões de dependências
- **SEMPRE** incluir troubleshooting para problemas comuns
- **NUNCA** deixar links quebrados ou referências obsoletas
- **NUNCA** usar linguagem ambígua em instruções técnicas
- **NUNCA** omitir pré-requisitos importantes
- **NUNCA** incluir credenciais ou segredos, mesmo em exemplos