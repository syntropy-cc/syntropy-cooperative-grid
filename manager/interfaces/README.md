# Interfaces

Interfaces de usuário do Syntropy Manager. Cada interface consome a API central para fornecer diferentes formas de interação com o manager.

## Estrutura

- **cli/**: Interface de linha de comando (Go + Cobra)
- **desktop/**: Aplicação desktop (Electron + React)
- **web/**: Frontend web (React + TypeScript)
- **mobile/**: Aplicação mobile (Flutter)

## Responsabilidades

- Fornecer interfaces específicas para diferentes tipos de usuários
- Consumir a API central do manager
- Implementar UX/UI otimizada para cada plataforma
- Gerenciar estado local e cache quando necessário
- Implementar autenticação específica da plataforma
