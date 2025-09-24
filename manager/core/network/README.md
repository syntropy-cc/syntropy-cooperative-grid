# Network Control

Controle e configuração da rede Syntropy.

## Responsabilidades

- Gerenciar nós do cluster Kubernetes
- Configurar Wireguard mesh network
- Implementar descoberta de nós
- Aplicar configurações de rede
- Monitorar conectividade e saúde da rede

## Componentes

- **NetworkController**: Controlador principal da rede
- **WireguardManager**: Gerenciamento do Wireguard
- **DiscoveryManager**: Descoberta e registro de nós
- **NetworkConfig**: Configurações de rede
- **HealthMonitor**: Monitoramento de saúde da rede
