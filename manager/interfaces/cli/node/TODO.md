# Node Component - Lista de Afazeres

## Progresso Geral
- [x] Estrutura de diretórios criada
- [x] Documentação completa (docs/)
- [x] Código fonte base (src/)
- [ ] Testes (tests/)
- [ ] Integração com CLI

## Etapa 1: Documentação Completa (docs/) - RIGOROSA

### Documentos Obrigatórios
- [ ] TODO.md (este arquivo)
- [ ] README.md - Visão geral rápida
- [ ] DOC.md - Quick start e overview (2 páginas)
- [ ] DEV.md - Arquitetura e implementação (8 páginas)
- [ ] API.md - Interface completa (12 páginas)
- [ ] TEST.md - Estratégias de teste (6 páginas)
- [ ] LEARN.md - Jornada educacional (20+ páginas)
- [ ] GUIDE.md - Guia de implementação para LLMs (16 seções)

### Progresso Documentação
- [x] TODO.md criado
- [x] README.md
- [x] DOC.md
- [x] DEV.md
- [x] API.md
- [x] TEST.md
- [x] LEARN.md
- [x] GUIDE.md

## Etapa 2: Estrutura Base (src/)

### Arquivos Base
- [x] src/node.go - Orquestrador principal (500 linhas)
- [x] src/types.go - Tipos públicos exportados
- [x] src/internal/types/interfaces.go - Interfaces internas
- [x] src/internal/helpers/helpers.go - Funções auxiliares
- [x] src/internal/constants/constants.go - Constantes

## Etapa 3: Token Integration

### Integração com Setup
- [x] src/token_integration.go - Integração com Setup TokenManager
  - [x] Importar TokenManager do Setup
  - [x] Obter Grid Token do Keyring
  - [x] Validar token antes de criar nó
  - [x] Compartilhar token com cloud-init

## Etapa 4: Auto Config Generator

### Geração Automática
- [x] src/auto_config_generator.go - Configurações automáticas
  - [x] NodeID sequencial (node-01, node-02...)
  - [x] Par de chaves SSH RSA 2048
  - [x] Certificado do nó Ed25519
  - [x] Detectar IP da Command Station
  - [x] Criar estrutura NodeConfig completa

## Etapa 5: Cloud-Init Generator ⚠️ CRÍTICO

### Cloud-Init Completo
- [x] src/cloud_init_generator.go - Geração de cloud-init
  - [x] user-data.yaml com instalação do Syntropy CLI
  - [x] Script de auto-registro completo
  - [x] Serviços systemd para heartbeat
  - [x] network-config.yaml (DHCP)
  - [x] meta-data.yaml (hostname, instance-id)
  - [x] Validar sintaxe YAML

## Etapa 6: USB Detection

### Detecção Multi-plataforma
- [x] src/usb_detector.go - Interface base
- [x] src/usb_detector_windows.go - Implementação Windows (WMIC/PowerShell)
- [x] src/usb_detector_linux.go - Implementação Linux (lsblk/fdisk)
- [x] Validação de segurança (não gravar em disco do sistema)

## Etapa 7: ISO Downloader

### Download e Cache
- [x] src/iso_downloader.go - Download de ISOs
  - [x] Ubuntu Server 24.04 LTS
  - [x] Verificação SHA256
  - [x] Cache em ~/.syntropy/cache/isos/
  - [x] Progress bar

## Etapa 8: USB Writer

### Gravação Multi-plataforma
- [x] src/usb_writer.go - Interface base
- [x] src/usb_writer_windows.go - Implementação Windows (dd/diskpart)
- [x] src/usb_writer_linux.go - Implementação Linux (dd/sync)
- [x] Injeção de cloud-init no ISO

## Etapa 9: Create Subcomponent

### Orquestração de Criação
- [x] src/create.go - Subcomponente Create
  - [ ] Coordenar todos os componentes de criação
  - [ ] Validar pré-requisitos
  - [ ] Gerar configurações
  - [ ] Criar USB bootável
  - [ ] Iniciar listener

## Etapa 10: Handshake Protocol ⚠️ CRÍTICO

### Handshake Seguro
- [x] src/handshake.go - Protocolo de handshake
  - [ ] Receber NodeAnnouncement do nó
  - [ ] Validar Grid Token contra Keyring
  - [ ] Validar NodeID e certificado
  - [ ] Enviar configurações para o nó
  - [ ] Confirmar registro

## Etapa 11: Listener

### Listener Automático
- [x] src/listener.go - Listener TCP
  - [ ] TCP na porta 51000
  - [ ] Um listener por nó criado
  - [ ] Timeout de 30 minutos
  - [ ] Processar handshake
  - [ ] Notificar sucesso/falha

## Etapa 12: Heartbeat

### Manutenção de Conexão
- [x] src/heartbeat.go - Heartbeat contínuo
  - [ ] Heartbeat a cada 30 segundos
  - [ ] Detectar nós inativos (3 falhas)
  - [ ] Reconexão automática
  - [ ] Coleta de métricas básicas

## Etapa 13: Node State

### Gerenciamento de Estado
- [ ] src/node_state.go - Estado de nós
  - [ ] Lista de nós pendentes (aguardando conexão)
  - [ ] Lista de nós ativos (conectados)
  - [ ] Lista de nós inativos (desconectados)
  - [ ] Transições de estado thread-safe
  - [ ] Persistência em ~/.syntropy/nodes/

## Etapa 14: Registration Subcomponent

### Orquestração de Registro
- [ ] src/registration.go - Subcomponente Registration
  - [ ] Coordenar listener, handshake e heartbeat
  - [ ] Gerenciar estado de nós
  - [ ] Processar eventos de registro
  - [ ] Logging estruturado

## Etapa 15: Main Orchestrator

### Orquestrador Principal
- [ ] src/node.go - API pública
  - [ ] Integrar Create e Registration
  - [ ] Comandos: Create, List, Status, Logs
  - [ ] Tratamento de erros unificado

## Etapa 16: CLI Integration

### Integração com CLI
- [x] Integrar com manager/interfaces/cli/main.go
  - [x] Comando `syntropy node create`
  - [x] Comando `syntropy node list`
  - [x] Comando `syntropy node status <id>`
  - [x] Comando `syntropy node logs <id>`

## Etapa 17: Suite de Testes Completa - RIGOROSA

### Análise Pré-Testes (OBRIGATÓRIO)
- [ ] Ler TODOS os arquivos em src/
- [ ] Catalogar TODOS os utilitários em src/internal/
- [ ] Documentar análise em tests/analysis.md
- [ ] Criar estratégia em tests/test-strategy.md

### Estrutura de Testes
- [ ] tests/unit/ - Testes unitários (1 arquivo por src/)
- [ ] tests/integration/ - Testes de integração
- [ ] tests/e2e/ - Testes end-to-end
- [ ] tests/performance/ - Testes de performance
- [ ] tests/security/ - Testes de segurança
- [ ] tests/fixtures/ - Dados de teste
- [ ] tests/mocks/ - Test doubles
- [ ] tests/helpers/ - Helpers de teste

### Cobertura Obrigatória
- [ ] 100% line coverage de TODOS os arquivos src/
- [ ] 100% branch coverage de TODOS os arquivos src/
- [ ] 100% function coverage de TODOS os arquivos src/
- [ ] Todos os subcomponentes têm testes dedicados
- [ ] Todos os códigos OS-específicos têm testes de plataforma

## Etapa 18: Documentação Final

### Finalização
- [ ] README.md - Guia do usuário
- [ ] EXAMPLES.md - Exemplos práticos
- [ ] Atualizar CLI README com comandos Node

## Critérios de Sucesso

- [ ] `syntropy node create` cria USB plug-and-play automaticamente
- [ ] Cloud-init instala Syntropy CLI no nó
- [ ] Nó executa handshake automático ao bootar
- [ ] Registro é aceito pela Command Station
- [ ] Nó fica PRONTO para workloads sem intervenção
- [ ] Heartbeat mantém conexão ativa
- [ ] Testes 100% cobertura (test-suite.rules.md)
- [ ] Documentação completa (documentation.rules.md)

## Observações Críticas

1. **Plug-and-Play**: Tudo 100% automático após `syntropy node create`
2. **Cloud-Init**: Arquivo mais crítico - instala tudo no nó
3. **Segurança 3 Camadas**: Grid Token + Node Certificate + SSH Keys
4. **Handshake**: Protocolo TCP/JSON na porta 51000
5. **CLI no Nó**: Baixado de GitHub releases durante cloud-init
6. **Systemd Services**: Auto-registro e heartbeat como services
7. **src/ vs tests/**: Código fonte em src/, testes em tests/ (imutável)
8. **Manager Separado**: Componente Manager é diferente, não entra em node/
9. **TODO.md**: Manter lista atualizada de progresso
10. **Test Suite**: Seguir test-suite.rules.md para 100% coverage
11. **Documentação**: Seguir documentation.rules.md rigorosamente

---

**Última atualização**: 2025-01-27
**Status**: Em progresso - Etapa 1: Documentação
