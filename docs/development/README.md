# 🛠️ Guia de Desenvolvimento - Syntropy Cooperative Grid

> **Documentação para Desenvolvedores**

## 📋 **Índice**

1. [Setup do Ambiente](#setup-do-ambiente)
2. [Estrutura do Projeto](#estrutura-do-projeto)
3. [Convenções de Código](#convenções-de-código)
4. [Workflow de Desenvolvimento](#workflow-de-desenvolvimento)
5. [Testes](#testes)
6. [Debugging](#debugging)
7. [Contribuição](#contribuição)

---

## 🚀 **Setup do Ambiente**

### **Pré-requisitos**

- **Go 1.21+**: [Download](https://golang.org/dl/)
- **Node.js 18+**: [Download](https://nodejs.org/)
- **Docker & Docker Compose**: [Download](https://www.docker.com/)
- **Git**: [Download](https://git-scm.com/)
- **Make**: Para usar o Makefile

### **Instalação**

```bash
# Clone o repositório
git clone https://github.com/syntropy-cc/cooperative-grid.git
cd cooperative-grid

# Instale dependências
make install

# Inicie o ambiente de desenvolvimento
make up

# Verifique se tudo está funcionando
make test
```

### **Configuração do IDE**

#### **VS Code**
```json
// .vscode/settings.json
{
  "go.toolsManagement.checkForUpdates": "local",
  "go.useLanguageServer": true,
  "go.lintTool": "golangci-lint",
  "go.formatTool": "goimports",
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.organizeImports": true
  }
}
```

#### **Extensões Recomendadas**
- Go (Google)
- REST Client
- Docker
- Kubernetes
- GitLens

---

## 📁 **Estrutura do Projeto**

```
syntropy-cooperative-grid/
├── cmd/                    # Entry points das aplicações
│   ├── cli/               # CLI interface
│   ├── api-server/        # API Gateway
│   ├── node-manager/      # Node management service
│   ├── container-manager/ # Container orchestration
│   └── network-manager/   # Network management
├── internal/              # Código privado da aplicação
│   ├── api/              # API handlers e middleware
│   ├── services/         # Business logic services
│   ├── platform/         # Platform-specific code
│   ├── storage/          # Data access layer
│   ├── config/           # Configuration management
│   └── utils/            # Shared utilities
├── pkg/                   # Pacotes públicos reutilizáveis
│   ├── models/           # Data models
│   ├── types/            # Type definitions
│   ├── errors/           # Error handling
│   └── constants/        # Application constants
├── web/                   # Interface web
├── mobile/                # App mobile
├── desktop/               # App desktop
├── api/                   # Definições de API
├── deployments/           # Configurações de deploy
├── docs/                  # Documentação
└── scripts/               # Scripts de build
```

### **Convenções de Nomenclatura**

#### **Go**
- **Pacotes**: lowercase, sem underscores (`usermanagement`)
- **Funções**: PascalCase para públicas, camelCase para privadas
- **Variáveis**: camelCase
- **Constantes**: PascalCase ou UPPER_CASE

#### **Arquivos**
- **Go**: snake_case (`user_service.go`)
- **Config**: kebab-case (`user-config.yaml`)
- **Docs**: kebab-case (`user-guide.md`)

---

## 📝 **Convenções de Código**

### **Go Style Guide**

#### **Estrutura de Arquivo**
```go
package usermanagement

import (
    "context"
    "fmt"
    
    "github.com/gin-gonic/gin"
    "github.com/syntropy-cc/cooperative-grid/pkg/models"
)

// Service handles user management operations
type Service struct {
    repo UserRepository
    log  *logrus.Logger
}

// NewService creates a new user service
func NewService(repo UserRepository, log *logrus.Logger) *Service {
    return &Service{
        repo: repo,
        log:  log,
    }
}

// CreateUser creates a new user
func (s *Service) CreateUser(ctx context.Context, req *CreateUserRequest) (*models.User, error) {
    // Implementation
}
```

#### **Error Handling**
```go
// Use custom errors
var (
    ErrUserNotFound = errors.New("user not found")
    ErrInvalidInput = errors.New("invalid input")
)

// Wrap errors with context
func (s *Service) GetUser(ctx context.Context, id string) (*models.User, error) {
    user, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user %s: %w", id, err)
    }
    return user, nil
}
```

#### **Logging**
```go
// Structured logging
s.log.WithFields(logrus.Fields{
    "user_id": userID,
    "action":  "create_user",
    "duration_ms": time.Since(start).Milliseconds(),
}).Info("User created successfully")
```

### **API Design**

#### **REST Endpoints**
```go
// Consistent URL structure
GET    /api/v1/users
POST   /api/v1/users
GET    /api/v1/users/{id}
PUT    /api/v1/users/{id}
DELETE /api/v1/users/{id}

// Nested resources
GET    /api/v1/users/{id}/nodes
POST   /api/v1/users/{id}/nodes
```

#### **Response Format**
```go
type APIResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   *APIError   `json:"error,omitempty"`
    Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
    Timestamp string `json:"timestamp"`
    Version   string `json:"version"`
    RequestID string `json:"request_id"`
}
```

---

## 🔄 **Workflow de Desenvolvimento**

### **Git Workflow**

#### **Branching Strategy**
```
main
├── develop
│   ├── feature/user-management
│   ├── feature/container-orchestration
│   └── hotfix/security-patch
└── release/v1.0.0
```

#### **Commit Messages**
```
feat: add user authentication system
fix: resolve USB detection issue on Windows
docs: update API documentation
test: add unit tests for user service
refactor: simplify node creation logic
```

### **Development Process**

#### **1. Feature Development**
```bash
# Create feature branch
git checkout -b feature/user-management

# Make changes
# ... code changes ...

# Commit changes
git add .
git commit -m "feat: add user management endpoints"

# Push branch
git push origin feature/user-management

# Create pull request
gh pr create --title "Add user management system" --body "Implements user CRUD operations"
```

#### **2. Code Review**
- **Reviewers**: Pelo menos 2 aprovações
- **Checks**: Todos os testes devem passar
- **Coverage**: Mínimo 80% de cobertura
- **Security**: Scan de segurança deve passar

#### **3. Merge Process**
```bash
# Squash and merge
git checkout main
git pull origin main
git merge --squash feature/user-management
git commit -m "feat: add user management system"
git push origin main
```

---

## 🧪 **Testes**

### **Estrutura de Testes**

#### **Unit Tests**
```go
// user_service_test.go
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
// integration_test.go
func TestUserAPI_Integration(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    // Setup test server
    server := setupTestServer(t, db)
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

# Run specific package tests
go test ./internal/services/user

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...

# Run integration tests
go test -tags=integration ./...
```

---

## 🐛 **Debugging**

### **Local Development**

#### **Debugging Go Services**
```bash
# Run with debugger
dlv debug cmd/api-server/main.go

# Run with profiling
go run -cpuprofile=cpu.prof cmd/api-server/main.go

# Analyze profile
go tool pprof cpu.prof
```

#### **Debugging Web Frontend**
```bash
# Run with debug mode
cd web/frontend
npm run dev:debug

# Use browser dev tools
# Chrome DevTools -> Sources -> Debugger
```

### **Logging**

#### **Structured Logging**
```go
// Development logging
log := logrus.New()
log.SetLevel(logrus.DebugLevel)
log.SetFormatter(&logrus.JSONFormatter{})

// Production logging
log := logrus.New()
log.SetLevel(logrus.InfoLevel)
log.SetFormatter(&logrus.JSONFormatter{
    TimestampFormat: time.RFC3339,
})
```

#### **Log Levels**
- **DEBUG**: Informações detalhadas para debugging
- **INFO**: Informações gerais sobre o funcionamento
- **WARN**: Avisos sobre situações anômalas
- **ERROR**: Erros que não impedem o funcionamento
- **FATAL**: Erros críticos que impedem o funcionamento

---

## 🤝 **Contribuição**

### **Como Contribuir**

#### **1. Fork e Clone**
```bash
# Fork no GitHub, depois clone
git clone https://github.com/SEU_USERNAME/cooperative-grid.git
cd cooperative-grid

# Adicione o upstream
git remote add upstream https://github.com/syntropy-cc/cooperative-grid.git
```

#### **2. Desenvolvimento**
```bash
# Crie uma branch para sua feature
git checkout -b feature/sua-feature

# Faça suas alterações
# ... código ...

# Teste suas alterações
make test
make lint

# Commit suas alterações
git add .
git commit -m "feat: sua feature description"
```

#### **3. Pull Request**
```bash
# Push sua branch
git push origin feature/sua-feature

# Crie um Pull Request
gh pr create --title "Sua Feature" --body "Descrição detalhada"
```

### **Guidelines de Contribuição**

#### **Code Quality**
- ✅ Código limpo e bem documentado
- ✅ Testes unitários e de integração
- ✅ Cobertura de testes mínima de 80%
- ✅ Linting sem erros
- ✅ Seguimento das convenções de código

#### **Documentation**
- ✅ README atualizado se necessário
- ✅ Comentários em código complexo
- ✅ Documentação de APIs
- ✅ Exemplos de uso

#### **Testing**
- ✅ Testes para casos de sucesso
- ✅ Testes para casos de erro
- ✅ Testes de edge cases
- ✅ Testes de performance se aplicável

---

## 📚 **Recursos Adicionais**

### **Documentação**
- [Go Documentation](https://golang.org/doc/)
- [Gin Framework](https://gin-gonic.com/docs/)
- [Cobra CLI](https://cobra.dev/)
- [GORM](https://gorm.io/docs/)

### **Ferramentas**
- [golangci-lint](https://golangci-lint.run/)
- [GoLand IDE](https://www.jetbrains.com/go/)
- [Postman](https://www.postman.com/) para testes de API
- [Docker Desktop](https://www.docker.com/products/docker-desktop)

### **Comunidade**
- [Discord](https://discord.gg/syntropy) - Chat da comunidade
- [GitHub Discussions](https://github.com/syntropy-cc/cooperative-grid/discussions)
- [Issues](https://github.com/syntropy-cc/cooperative-grid/issues)

---

**Bem-vindo à equipe de desenvolvimento do Syntropy Cooperative Grid! 🚀**


