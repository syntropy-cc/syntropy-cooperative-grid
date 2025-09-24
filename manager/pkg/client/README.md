# Client

Cliente para comunicação com a rede Syntropy.

## Responsabilidades

- Fornecer interface para comunicação com Kubernetes
- Implementar cliente para Wireguard mesh
- Comunicar com sistema de créditos
- Monitorar estado da rede
- Gerenciar conexões e autenticação

## Componentes

- **KubernetesClient**: Cliente para cluster Kubernetes
- **WireguardClient**: Cliente para Wireguard mesh
- **CreditClient**: Cliente para sistema de créditos
- **DiscoveryClient**: Cliente para descoberta de nós
- **GridClient**: Cliente unificado para a rede
