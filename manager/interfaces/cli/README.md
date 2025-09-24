# CLI Interface

Interface de linha de comando do Syntropy Manager, implementada em Go usando Cobra.

## Responsabilidades

- Fornecer comandos para gerenciamento de nós
- Implementar comandos para criação e gerenciamento de workloads
- Permitir monitoramento do estado da rede
- Configurar parâmetros do manager e rede
- Integrar com sistema cloud-init existente

## Comandos Principais

- **node**: Adicionar, listar, remover e configurar nós
- **workload**: Criar, gerenciar e monitorar workloads
- **state**: Visualizar estado da rede
- **config**: Configurar parâmetros do manager
- **init**: Inicializar configuração do manager
