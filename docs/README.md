# Documentação - Syntropy Cooperative Grid

## Visão Geral

Este diretório contém toda a documentação técnica do projeto Syntropy Cooperative Grid, incluindo arquitetura, implementações e guias de uso.

## Estrutura da Documentação

### 📋 Documentação Principal
- **[architecture/README.md](architecture/README.md)** - Arquitetura técnica completa do projeto
- **[implementation-summary.md](implementation-summary.md)** - Resumo das implementações realizadas

### 🔧 Documentação Técnica
- **[cloud-init-architecture.md](cloud-init-architecture.md)** - Arquitetura detalhada do sistema de cloud-init
- **[cli-usb-guide.md](cli-usb-guide.md)** - Guia completo do comando USB da CLI
- **[installation-scripts.md](installation-scripts.md)** - Documentação dos scripts de instalação

### 📁 Diretórios
- **architecture/** - Documentação de arquitetura do sistema
- **api/** - Documentação de APIs (futuro)
- **deployment/** - Guias de deployment (futuro)
- **tutorials/** - Tutoriais e exemplos (futuro)

## Implementações Realizadas

### ✅ Sistema de Cloud-Init Completo
- Templates cloud-init personalizados
- Scripts de instalação inteligentes
- Sistema de segurança automática
- Descoberta inteligente de rede
- Auditoria completa

### ✅ CLI Go para Gerenciamento de USBs
- Comando `syntropy usb create`
- Comando `syntropy usb list`
- Comando `syntropy usb format`
- Suporte a múltiplas plataformas

### ✅ Sistema de Segurança
- Geração automática de certificados TLS
- Geração automática de chaves SSH
- Firewall automático
- Sistema de auditoria

## Como Usar

### Criar um USB Bootável
```bash
# Criar USB com auto-detecção
syntropy usb create --auto-detect --node-name "node-01"

# Criar USB especificando dispositivo
syntropy usb create /dev/sdb --node-name "node-02" --description "Servidor principal"

# Criar USB com coordenadas
syntropy usb create --auto-detect --node-name "node-03" --coordinates "-23.5505,-46.6333"
```

### Listar Dispositivos USB
```bash
# Listagem em formato tabela
syntropy usb list

# Listagem em formato JSON
syntropy usb list --format json

# Listagem em formato YAML
syntropy usb list --format yaml
```

### Formatar USB
```bash
# Formatar USB
syntropy usb format /dev/sdb

# Formatar com rótulo personalizado
syntropy usb format /dev/sdb --label "MYUSB"

# Formatar sem confirmação
syntropy usb format /dev/sdb --force
```

## Arquitetura

### PC de Trabalho (Quartel General)
O PC de trabalho atua como centro de comando, gerando USBs personalizados com:
- Certificados TLS únicos
- Chaves SSH específicas
- Configuração cloud-init personalizada
- Scripts de instalação

### USB Bootável (DNA do Nó)
Cada USB contém:
- ISO Ubuntu personalizada
- Configuração cloud-init
- Certificados e chaves
- Scripts de instalação

### Hardware Virgem (Nó Syntropy)
O hardware inicia automaticamente e:
- Detecta tipo de hardware
- Configura rede
- Instala Syntropy Agent
- Conecta ao cluster

## Segurança

### Certificados TLS
- CA gerada automaticamente (RSA 4096 bits)
- Certificados de nó únicos (RSA 2048 bits)
- Validade: CA (10 anos), Nó (1 ano)

### Chaves SSH
- Par de chaves RSA 2048 bits
- Acesso apenas por chave pública
- Usuário: syntropy

### Firewall
- UFW configurado automaticamente
- Regras específicas para Syntropy
- Fail2ban para proteção

## Descoberta de Rede

### Métodos Implementados
1. **DNS**: Resolve `syntropy-discovery.local`
2. **Broadcast**: Envia na rede local
3. **Multicast**: Usa grupos específicos
4. **Configuração Manual**: Hosts pré-definidos

### Fallbacks Automáticos
- Se DNS falha → tenta broadcast
- Se broadcast falha → tenta multicast
- Se multicast falha → tenta manual
- Se tudo falha → vira líder (primeiro nó)

## Auditoria

### Logs Centralizados
- `~/.syntropy/nodes/*/audit.log`
- Logs de todas as operações
- Retenção configurável (90 dias)

### Rastreamento Completo
- Boot e inicialização
- Descoberta de rede
- Conexão ao cluster
- Operações de crédito
- Eventos de segurança

## Troubleshooting

### Problemas Comuns
- Dispositivo USB não detectado
- Falha na geração de certificados
- Falha no download da ISO
- Falha na gravação do USB

### Comandos de Diagnóstico
```bash
# Verificar status do USB
lsblk
fdisk -l /dev/sdb

# Verificar conectividade
ping 8.8.8.8
curl -I https://releases.ubuntu.com/

# Verificar recursos
free -h
df -h
```

## Contribuição

Para contribuir com a documentação:

1. Faça fork do repositório
2. Crie uma branch para sua contribuição
3. Faça as alterações necessárias
4. Teste as alterações
5. Submeta um pull request

## Licença

Este projeto está licenciado sob a licença MIT. Veja o arquivo LICENSE para mais detalhes.

## Suporte

Para suporte técnico:
- Abra uma issue no GitHub
- Consulte a documentação técnica
- Entre em contato com a equipe de desenvolvimento
