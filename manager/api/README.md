# API

API central do Syntropy Manager. Contém todos os componentes relacionados à API REST, gRPC e WebSocket que servem como base para todas as interfaces do manager.

## Estrutura

- **handlers/**: Implementação dos handlers HTTP para endpoints REST
- **middleware/**: Middleware para autenticação, logging, CORS, rate limiting
- **routes/**: Definição e configuração das rotas da API

## Responsabilidades

- Expor endpoints para gerenciamento de nós, workloads e configurações
- Implementar autenticação e autorização
- Fornecer API unificada para todas as interfaces (CLI, Web, Mobile)
- Gerenciar comunicação com o core do manager
