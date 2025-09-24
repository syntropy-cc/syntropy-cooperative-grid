# Core

Lógica central do Syntropy Manager. Contém todos os componentes que implementam a funcionalidade principal do manager.

## Estrutura

- **state/**: Gerenciamento de estado desejado da rede
- **network/**: Controle e configuração da rede
- **workload/**: Gerenciamento de workloads e containers
- **reconciliation/**: Sistema de reconciliação entre estado desejado e atual

## Responsabilidades

- Manter estado desejado da rede
- Controlar operações na rede Kubernetes
- Gerenciar criação e execução de workloads
- Implementar sistema de reconciliação contínua
- Coordenar com componentes da rede existente
