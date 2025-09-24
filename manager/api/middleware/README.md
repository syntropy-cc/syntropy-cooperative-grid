# Middleware

Middleware para processamento de requisições HTTP da API do Syntropy Manager.

## Responsabilidades

- **Autenticação**: Validação de tokens JWT
- **Autorização**: Verificação de permissões RBAC
- **Logging**: Registro estruturado de requisições
- **CORS**: Configuração de Cross-Origin Resource Sharing
- **Rate Limiting**: Controle de taxa de requisições
- **Recovery**: Tratamento de panics e erros não capturados

## Middleware Disponíveis

- **AuthMiddleware**: Autenticação JWT
- **LoggingMiddleware**: Logs estruturados
- **CORSMiddleware**: Configuração CORS
- **RateLimitMiddleware**: Limitação de taxa
- **RecoveryMiddleware**: Recuperação de panics
