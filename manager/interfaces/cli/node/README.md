# Node Component

## Visão Geral

O **Node Component** é o componente responsável pelo provisionamento automático e registro de nós físicos na Syntropy Cooperative Grid. Ele implementa um sistema plug-and-play completo que permite criar USBs bootáveis automaticamente e registrar nós na rede sem intervenção manual.

## Funcionalidades Principais

- 🔍 **Detecção automática de USBs** - Identifica dispositivos USB disponíveis
- 📥 **Download e cache de ISOs** - Ubuntu Server com verificação de integridade
- ⚙️ **Geração automática de configurações** - NodeID, Grid Token, SSH keys
- 💾 **Gravação de USBs bootáveis** - Multi-plataforma com cloud-init customizado
- 🔐 **Registro automático** - Handshake automático após instalação
- 🛡️ **Validações de segurança** - Prevenção de gravação em discos do sistema
- 🖥️ **Suporte multi-plataforma** - Windows/Linux/macOS
- 🔄 **Múltiplos nós simultâneos** - Gerenciamento de vários nós pendentes

## Quick Start

### Pré-requisitos
- Setup Component configurado (100% completo)
- Grid Token disponível via Keyring
- Dispositivo USB (mínimo 8GB)
- Ubuntu Server 24.04 LTS (baixado automaticamente)

### Comando Principal
```bash
# Criação automática de nó (comando único)
syntropy node create

# Comportamento automático:
# 1. Detecta USB automaticamente
# 2. Gera configurações automaticamente
# 3. Cria USB bootável
# 4. Inicia listener automático
# 5. Aguarda conexão do nó
```

### Comandos Disponíveis
```bash
# Listar nós
syntropy node list

# Status de nó específico
syntropy node status <node-id>

# Logs de nó específico
syntropy node logs <node-id>

# Especificar USB específico
syntropy node create --usb /dev/sdb
```

## Fluxo Plug-and-Play

```
1. Usuário: syntropy node create
   ↓
2. Sistema detecta USB automaticamente
   ↓
3. Sistema gera configurações (NodeID, chaves SSH, tokens)
   ↓
4. Sistema cria cloud-init com instalação automática do Syntropy CLI
   ↓
5. Sistema grava USB bootável
   ↓
6. Sistema inicia Listener na Command Station
   ↓
7. Usuário conecta USB no hardware virgem
   ↓
8. Hardware faz boot automático
   ↓
9. Nó executa cloud-init:
   - Instala Ubuntu Server
   - Instala Docker
   - Baixa Syntropy CLI
   - Executa handshake automático
   - Registra na Command Station
   ↓
10. Command Station aceita nó
    ↓
11. Nó fica PRONTO para workloads
```

## Estrutura do Componente

```
node/
├── docs/           # Documentação completa
├── src/            # Código fonte
├── tests/          # Testes (100% cobertura)
├── GUIDE.md        # Guia de implementação para LLMs
├── README.md       # Este arquivo
└── TODO.md         # Lista de afazeres
```

## Arquitetura de Segurança

### Sistema de Autenticação em 3 Camadas

1. **Grid Token** (Nível Global)
   - Origem: Setup Component via Keyring
   - Uso: Autenticação inicial do nó na grid

2. **Node Certificate** (Nível Nó)
   - Formato: Par de chaves Ed25519
   - Uso: Identificação única e permanente do nó

3. **SSH Keys** (Nível Comunicação)
   - Formato: RSA 2048 bits
   - Uso: Comunicação SSH entre Command Station e Nó

## Próximos Passos

- [Explorar a API](./docs/API.md) - Instruções detalhadas de uso
- [Guia do Desenvolvedor](./docs/DEV.md) - Entendendo os internos
- [Guia de Testes](./docs/TEST.md) - Executando e escrevendo testes
- [Jornada de Aprendizado](./docs/LEARN.md) - Mergulho profundo nos conceitos
- [Exemplos](../examples/) - Mais cenários de uso
- [Guia de Implementação](./GUIDE.md) - Para LLMs implementarem

## Suporte

- **Issue Tracker**: [Link para issues]
- **Discussion Forum**: [Link para discussões]
- **Contact**: [Método de contato]

## Licença

MIT - Veja arquivo LICENSE para detalhes

---

**Status**: 🚧 Em desenvolvimento
**Versão**: 0.1.0
**Última atualização**: 2025-01-27
