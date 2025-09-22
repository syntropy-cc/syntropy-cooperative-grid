# 🤝 Guia de Contribuição - Syntropy Cooperative Grid

Obrigado por considerar contribuir para o Syntropy Cooperative Grid! Este documento fornece diretrizes e informações para contribuidores.

## 📋 **Índice**

1. [Código de Conduta](#código-de-conduta)
2. [Como Contribuir](#como-contribuir)
3. [Processo de Desenvolvimento](#processo-de-desenvolvimento)
4. [Padrões de Código](#padrões-de-código)
5. [Testes](#testes)
6. [Documentação](#documentação)
7. [Reportar Bugs](#reportar-bugs)
8. [Sugerir Melhorias](#sugerir-melhorias)

---

## 📜 **Código de Conduta**

### **Nossa Promessa**

No interesse de promover um ambiente aberto e acolhedor, nós, como contribuidores e mantenedores, nos comprometemos a fazer da participação em nosso projeto e em nossa comunidade uma experiência livre de assédio para todos, independentemente de idade, tamanho corporal, deficiência, etnia, características sexuais, identidade e expressão de gênero, nível de experiência, educação, status socioeconômico, nacionalidade, aparência pessoal, raça, religião ou identidade e orientação sexual.

### **Nossos Padrões**

Exemplos de comportamento que contribuem para criar um ambiente positivo incluem:

- ✅ Usar linguagem acolhedora e inclusiva
- ✅ Respeitar pontos de vista e experiências diferentes
- ✅ Aceitar graciosamente críticas construtivas
- ✅ Focar no que é melhor para a comunidade
- ✅ Mostrar empatia para com outros membros da comunidade

Exemplos de comportamento inaceitável incluem:

- ❌ Uso de linguagem ou imagens sexualizadas
- ❌ Comentários insultuosos ou depreciativos
- ❌ Assédio público ou privado
- ❌ Publicar informações privadas sem permissão
- ❌ Outras condutas inadequadas em um ambiente profissional

---

## 🚀 **Como Contribuir**

### **Tipos de Contribuição**

#### **🐛 Bug Reports**
- Reporte bugs através de [GitHub Issues](https://github.com/syntropy-cc/cooperative-grid/issues)
- Use o template de bug report
- Inclua informações detalhadas sobre o problema

#### **✨ Feature Requests**
- Sugira novas funcionalidades via [GitHub Discussions](https://github.com/syntropy-cc/cooperative-grid/discussions)
- Descreva o problema que a feature resolve
- Explique como a feature deve funcionar

#### **📝 Documentação**
- Melhore documentação existente
- Adicione exemplos de uso
- Traduza documentação para outros idiomas

#### **🧪 Testes**
- Adicione testes para funcionalidades existentes
- Melhore cobertura de testes
- Adicione testes de integração

#### **🔧 Código**
- Corrija bugs
- Implemente novas funcionalidades
- Refatore código existente
- Otimize performance

### **Processo de Contribuição**

#### **1. Fork e Setup**
```bash
# Fork o repositório no GitHub
# Clone seu fork
git clone https://github.com/SEU_USERNAME/cooperative-grid.git
cd cooperative-grid

# Adicione o upstream
git remote add upstream https://github.com/syntropy-cc/cooperative-grid.git

# Instale dependências
make install
```

#### **2. Criar Branch**
```bash
# Crie uma branch para sua contribuição
git checkout -b feature/sua-feature
# ou
git checkout -b fix/bug-description
# ou
git checkout -b docs/improve-readme
```

#### **3. Desenvolvimento**
```bash
# Faça suas alterações
# ... código ...

# Teste suas alterações
make test
make lint

# Commit suas alterações
git add .
git commit -m "feat: add user authentication system"
```

#### **4. Push e Pull Request**
```bash
# Push sua branch
git push origin feature/sua-feature

# Crie um Pull Request
gh pr create --title "Add user authentication system" --body "Descrição detalhada"
```

---

## 🔄 **Processo de Desenvolvimento**

### **Workflow de Branches**

```
main (production)
├── develop (integration)
│   ├── feature/user-auth
│   ├── feature/container-mgmt
│   └── hotfix/security-patch
└── release/v1.0.0
```

### **Convenções de Commit**

#### **Formato**
```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

#### **Tipos**
- **feat**: Nova funcionalidade
- **fix**: Correção de bug
- **docs**: Documentação
- **style**: Formatação, sem mudança de código
- **refactor**: Refatoração de código
- **test**: Adição ou correção de testes
- **chore**: Mudanças em build, dependências, etc.

#### **Exemplos**
```
feat(auth): add JWT authentication system
fix(usb): resolve USB detection issue on Windows
docs(api): update API documentation
test(user): add unit tests for user service
refactor(db): simplify database connection logic
```

### **Pull Request Process**

#### **Antes de Submeter**
- [ ] Código segue as convenções do projeto
- [ ] Testes passam localmente
- [ ] Linting não apresenta erros
- [ ] Documentação atualizada se necessário
- [ ] Commits seguem o formato correto

#### **Template de PR**
```markdown
## Descrição
Breve descrição das mudanças

## Tipo de Mudança
- [ ] Bug fix
- [ ] Nova funcionalidade
- [ ] Breaking change
- [ ] Documentação

## Checklist
- [ ] Testes adicionados/atualizados
- [ ] Documentação atualizada
- [ ] Código segue convenções
- [ ] Self-review realizado

## Screenshots (se aplicável)
Adicione screenshots para mudanças na UI

## Testes
Descreva os testes realizados
```

---

## 📝 **Padrões de Código**

### **Go**

#### **Style Guide**
- Siga o [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` para formatação
- Use `golint` para linting
- Documente funções públicas

#### **Estrutura de Arquivo**
```go
package example

import (
    "context"
    "fmt"
    
    "github.com/gin-gonic/gin"
    "github.com/syntropy-cc/cooperative-grid/pkg/models"
)

// Service handles example operations
type Service struct {
    repo Repository
    log  *logrus.Logger
}

// NewService creates a new service
func NewService(repo Repository, log *logrus.Logger) *Service {
    return &Service{
        repo: repo,
        log:  log,
    }
}
```

### **JavaScript/TypeScript**

#### **Style Guide**
- Use ESLint e Prettier
- Siga as convenções do Airbnb
- Use TypeScript para tipagem
- Documente funções complexas

#### **Estrutura de Componente**
```typescript
import React from 'react';
import { User } from '../types';

interface UserCardProps {
  user: User;
  onEdit: (user: User) => void;
}

export const UserCard: React.FC<UserCardProps> = ({ user, onEdit }) => {
  return (
    <div className="user-card">
      <h3>{user.name}</h3>
      <button onClick={() => onEdit(user)}>Edit</button>
    </div>
  );
};
```

---

## 🧪 **Testes**

### **Cobertura Mínima**
- **Go**: 80% de cobertura de código
- **JavaScript**: 80% de cobertura de código
- **Flutter**: 80% de cobertura de código

### **Tipos de Testes**

#### **Unit Tests**
```go
func TestUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name    string
        input   *CreateUserRequest
        want    *models.User
        wantErr bool
    }{
        {
            name: "valid user creation",
            input: &CreateUserRequest{
                Email: "test@example.com",
                Name:  "Test User",
            },
            want: &models.User{
                Email: "test@example.com",
                Name:  "Test User",
            },
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

#### **Integration Tests**
```go
func TestUserAPI_Integration(t *testing.T) {
    // Setup test environment
    server := setupTestServer(t)
    defer server.Close()
    
    // Test API endpoints
    t.Run("create user", func(t *testing.T) {
        resp := httptest.NewRecorder()
        req := createUserRequest(t, "test@example.com")
        
        server.ServeHTTP(resp, req)
        
        assert.Equal(t, http.StatusCreated, resp.Code)
    })
}
```

### **Running Tests**
```bash
# Run all tests
make test

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/services/user

# Run integration tests
go test -tags=integration ./...
```

---

## 📚 **Documentação**

### **Tipos de Documentação**

#### **Code Documentation**
```go
// CreateUser creates a new user in the system
// It validates the input and stores the user in the database
func (s *Service) CreateUser(ctx context.Context, req *CreateUserRequest) (*models.User, error) {
    // Implementation
}
```

#### **API Documentation**
```yaml
# OpenAPI specification
paths:
  /api/v1/users:
    post:
      summary: Create a new user
      description: Creates a new user in the system
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateUserRequest'
      responses:
        '201':
          description: User created successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/User'
```

#### **README Files**
- Documente como usar cada módulo
- Inclua exemplos de uso
- Mantenha atualizado

### **Documentation Standards**
- Use Markdown para documentação
- Inclua exemplos práticos
- Mantenha linguagem clara e concisa
- Atualize quando o código mudar

---

## 🐛 **Reportar Bugs**

### **Template de Bug Report**

```markdown
**Descrição do Bug**
Descrição clara e concisa do bug.

**Passos para Reproduzir**
1. Vá para '...'
2. Clique em '...'
3. Role até '...'
4. Veja o erro

**Comportamento Esperado**
Descrição clara do que deveria acontecer.

**Comportamento Atual**
Descrição do que está acontecendo.

**Screenshots**
Se aplicável, adicione screenshots.

**Ambiente:**
 - OS: [e.g. Windows 10, macOS 12, Ubuntu 20.04]
 - Versão: [e.g. v1.0.0]
 - Browser: [e.g. Chrome 91, Firefox 89]

**Informações Adicionais**
Qualquer outra informação relevante.
```

### **Critérios para Bug Reports**
- ✅ Descrição clara do problema
- ✅ Passos para reproduzir
- ✅ Comportamento esperado vs atual
- ✅ Informações do ambiente
- ✅ Screenshots se aplicável

---

## 💡 **Sugerir Melhorias**

### **Template de Feature Request**

```markdown
**Funcionalidade Sugerida**
Descrição clara da funcionalidade.

**Problema que Resolve**
Descrição do problema que esta funcionalidade resolve.

**Solução Proposta**
Descrição detalhada de como a funcionalidade deve funcionar.

**Alternativas Consideradas**
Outras soluções que foram consideradas.

**Contexto Adicional**
Qualquer outro contexto sobre a funcionalidade.
```

### **Critérios para Feature Requests**
- ✅ Descrição clara da funcionalidade
- ✅ Justificativa para a necessidade
- ✅ Proposta de implementação
- ✅ Consideração de alternativas
- ✅ Impacto na API (se aplicável)

---

## 🏷️ **Labels e Categorização**

### **Labels de Issues**
- **bug**: Algo não está funcionando
- **enhancement**: Nova funcionalidade ou melhoria
- **documentation**: Melhorias na documentação
- **good first issue**: Bom para novos contribuidores
- **help wanted**: Precisa de ajuda extra
- **priority:high**: Alta prioridade
- **priority:medium**: Média prioridade
- **priority:low**: Baixa prioridade

### **Labels de PRs**
- **ready for review**: Pronto para revisão
- **work in progress**: Em desenvolvimento
- **needs testing**: Precisa de testes
- **breaking change**: Mudança que quebra compatibilidade
- **security**: Relacionado a segurança

---

## 🎯 **Roadmap de Contribuição**

### **Para Novos Contribuidores**

#### **Primeiros Passos**
1. **Leia a documentação**: Comece com o README e docs de desenvolvimento
2. **Explore o código**: Entenda a estrutura do projeto
3. **Procure por "good first issue"**: Issues marcadas para iniciantes
4. **Participe da comunidade**: Discord, GitHub Discussions

#### **Issues Recomendadas**
- Correções de bugs simples
- Melhorias na documentação
- Adição de testes
- Pequenas funcionalidades

#### **Projetos Futuros**
- Implementação de novas interfaces
- Otimizações de performance
- Integrações com serviços externos
- Ferramentas de desenvolvimento

---

## 📞 **Suporte e Comunidade**

### **Canais de Comunicação**
- **GitHub Issues**: Para bugs e feature requests
- **GitHub Discussions**: Para discussões gerais
- **Discord**: Chat da comunidade
- **Email**: contato@syntropy.coop

### **Reuniões da Comunidade**
- **Weekly Sync**: Toda quarta-feira às 14:00 UTC
- **Monthly Review**: Primeira sexta-feira do mês
- **Sprint Planning**: A cada 2 semanas

---

## 🙏 **Reconhecimento**

### **Contribuidores**
Todos os contribuidores são reconhecidos no README do projeto e em releases.

### **Tipos de Contribuição**
- **Code**: Desenvolvimento de código
- **Documentation**: Melhoria da documentação
- **Testing**: Adição de testes
- **Community**: Ajuda na comunidade
- **Design**: Design de interfaces
- **Translation**: Tradução de documentação

---

**Obrigado por contribuir para o Syntropy Cooperative Grid! 🌐**

Juntos, estamos construindo o futuro da computação cooperativa.