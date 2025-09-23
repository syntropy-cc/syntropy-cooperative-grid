# Syntropy Cooperative Grid CLI - Exemplos de Uso

Este documento fornece exemplos práticos de como usar a CLI da Syntropy Cooperative Grid.

## 🚀 Quick Start

### 1. Setup Inicial

```bash
# Compilar a CLI
make build

# Setup do ambiente de gerenciamento
./bin/syntropy setup

# Verificar se o setup foi bem-sucedido
./bin/syntropy manager health
```

### 2. Criar Primeiro Nó

```bash
# Listar dispositivos USB disponíveis
./bin/syntropy usb list

# Criar USB bootável para o primeiro nó
./bin/syntropy usb create /dev/sdb --node-name "main-server" --coordinates "-23.5505,-46.6333"

# Aguardar instalação (~30 minutos) e depois descobrir o nó
./bin/syntropy manager discover
```

### 3. Gerenciar Nós

```bash
# Listar todos os nós
./bin/syntropy manager list

# Conectar ao nó
./bin/syntropy manager connect main-server

# Verificar status
./bin/syntropy manager status main-server
```

## 📋 Exemplos Detalhados

### Setup e Configuração

```bash
# Setup completo com reconfiguração
./bin/syntropy setup

# Verificar configuração do gerenciador
cat ~/.syntropy/config/manager.json

# Listar scripts auxiliares criados
ls ~/.syntropy/scripts/
```

### Gerenciamento de Nós

```bash
# Listar nós com filtros
./bin/syntropy manager list --filter online
./bin/syntropy manager list --sort last_seen

# Status detalhado em diferentes formatos
./bin/syntropy manager status --format json
./bin/syntropy manager status --format yaml

# Conectar com comando específico
./bin/syntropy manager connect main-server --command "docker ps"
./bin/syntropy manager connect main-server --command "df -h"
```

### Descoberta de Rede

```bash
# Descoberta básica
./bin/syntropy manager discover

# Descoberta em redes específicas
./bin/syntropy manager discover --networks "192.168.1.0/24,10.0.0.0/24"

# Descoberta com configurações customizadas
./bin/syntropy manager discover --port 2222 --timeout 15 --parallel 10

# Descoberta sem atualizar cache
./bin/syntropy manager discover --update-cache false
```

### Backup e Restore

```bash
# Backup completo
./bin/syntropy manager backup

# Backup com arquivo específico
./bin/syntropy manager backup --output /tmp/syntropy-backup.tar.gz

# Backup sem compressão
./bin/syntropy manager backup --compress false

# Backup de componentes específicos
./bin/syntropy manager backup --include "nodes,keys"

# Restore de backup
./bin/syntropy manager restore backup_20240115_143022.tar.gz

# Restore forçado (sem confirmação)
./bin/syntropy manager restore backup.tar.gz --force
```

### Health Check

```bash
# Health check de todos os nós
./bin/syntropy manager health

# Health check em formato JSON
./bin/syntropy manager health --format json

# Health check com monitoramento contínuo
./bin/syntropy manager health --watch
```

### Templates de Aplicação

```bash
# Listar templates disponíveis
./bin/syntropy templates list

# Listar por categoria
./bin/syntropy templates list --category scientific

# Mostrar detalhes de um template
./bin/syntropy templates show fortran-computation

# Deploy básico
./bin/syntropy templates deploy jupyter-lab --node main-server

# Deploy com valores customizados
./bin/syntropy templates deploy python-datascience --node main-server \
  --set "memory=2Gi" \
  --set "cpu=1000m" \
  --set "jupyter_token=mysecret123"

# Deploy em modo dry-run
./bin/syntropy templates deploy nginx --node main-server --dry-run

# Criar novo template
./bin/syntropy templates create my-web-app \
  --category web \
  --description "Minha aplicação web personalizada"
```

### Criação de USB

```bash
# Listar dispositivos USB
./bin/syntropy usb list --format json

# Criação básica
./bin/syntropy usb create /dev/sdb --node-name "edge-node-01"

# Criação com auto-detecção
./bin/syntropy usb create --auto-detect --node-name "edge-node-02"

# Criação com coordenadas
./bin/syntropy usb create /dev/sdb --node-name "sp-node" --coordinates "-23.5505,-46.6333"

# Criação com chave de proprietário existente
./bin/syntropy usb create /dev/sdb --node-name "secure-node" \
  --owner-key ~/.syntropy/keys/main.key

# Criação com configurações customizadas
./bin/syntropy usb create /dev/sdb --node-name "custom-node" \
  --description "Nó de produção em São Paulo" \
  --label "PROD-SP" \
  --work-dir /tmp/syntropy-work \
  --cache-dir ~/.syntropy/cache

# Formatação de USB
./bin/syntropy usb format /dev/sdb --label "SYNTROPY"

# Formatação forçada
./bin/syntropy usb format /dev/sdb --force
```

## 🔧 Workflows Completos

### Workflow 1: Deploy de Aplicação Científica

```bash
# 1. Setup inicial
./bin/syntropy setup

# 2. Criar nó de computação
./bin/syntropy usb create /dev/sdb --node-name "compute-node" --coordinates "-23.5505,-46.6333"

# 3. Aguardar instalação e descobrir
./bin/syntropy manager discover

# 4. Verificar se o nó está online
./bin/syntropy manager status compute-node

# 5. Deploy de aplicação Fortran
./bin/syntropy templates deploy fortran-computation --node compute-node

# 6. Monitorar execução
./bin/syntropy manager connect compute-node --command "docker logs fortran-simulation"
```

### Workflow 2: Setup de Laboratório de Data Science

```bash
# 1. Criar múltiplos nós
./bin/syntropy usb create /dev/sdb --node-name "jupyter-server"
./bin/syntropy usb create /dev/sdc --node-name "data-node-1"
./bin/syntropy usb create /dev/sdd --node-name "data-node-2"

# 2. Descobrir todos os nós
./bin/syntropy manager discover

# 3. Verificar saúde da rede
./bin/syntropy manager health

# 4. Deploy Jupyter Lab no servidor principal
./bin/syntropy templates deploy python-datascience --node jupyter-server \
  --set "memory=4Gi" \
  --set "cpu=2000m"

# 5. Deploy aplicações de dados nos nós secundários
./bin/syntropy templates deploy python-datascience --node data-node-1
./bin/syntropy templates deploy python-datascience --node data-node-2

# 6. Backup da configuração
./bin/syntropy manager backup --output lab-backup.tar.gz
```

### Workflow 3: Monitoramento e Manutenção

```bash
# 1. Health check regular
./bin/syntropy manager health --format json > health-report.json

# 2. Backup automático
./bin/syntropy manager backup

# 3. Atualizar descoberta de rede
./bin/syntropy manager discover --update-cache true

# 4. Verificar logs de nós específicos
./bin/syntropy manager connect main-server --command "journalctl -u docker --since '1 hour ago'"

# 5. Monitorar recursos
./bin/syntropy manager connect main-server --command "htop"
```

## 🎯 Casos de Uso Avançados

### Integração com Scripts

```bash
#!/bin/bash
# Script de monitoramento automático

# Health check e notificação
HEALTH_OUTPUT=$(./bin/syntropy manager health --format json)
OFFLINE_NODES=$(echo "$HEALTH_OUTPUT" | jq '.[] | select(.status == "offline") | .node_name')

if [ ! -z "$OFFLINE_NODES" ]; then
    echo "⚠️  Nós offline detectados: $OFFLINE_NODES"
    # Enviar notificação (email, Slack, etc.)
fi

# Backup automático se houver mudanças
./bin/syntropy manager backup
```

### Automação com Cron

```bash
# Adicionar ao crontab para execução automática
# Health check a cada 15 minutos
*/15 * * * * /path/to/syntropy manager health --format json >> /var/log/syntropy-health.log

# Backup diário às 2h da manhã
0 2 * * * /path/to/syntropy manager backup

# Descoberta de rede a cada hora
0 * * * * /path/to/syntropy manager discover --update-cache true
```

### Integração com CI/CD

```yaml
# .github/workflows/syntropy-deploy.yml
name: Deploy to Syntropy Grid

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Setup Syntropy CLI
        run: |
          make build
          ./bin/syntropy setup
      
      - name: Deploy Application
        run: |
          ./bin/syntropy templates deploy my-app --node production-node
      
      - name: Health Check
        run: |
          ./bin/syntropy manager health --format json
```

## 🐛 Troubleshooting

### Problemas Comuns

```bash
# Nó não aparece na descoberta
./bin/syntropy manager discover --networks "192.168.1.0/24" --timeout 30

# Erro de conexão SSH
./bin/syntropy manager connect node-01 --command "echo 'test'"

# Verificar chaves SSH
ls -la ~/.syntropy/keys/

# Verificar configuração
cat ~/.syntropy/config/manager.json

# Logs detalhados
./bin/syntropy manager health --format json | jq '.[] | select(.error)'
```

### Recuperação de Backup

```bash
# Listar backups disponíveis
ls -la ~/.syntropy/backups/

# Restore de backup específico
./bin/syntropy manager restore backup_20240115_143022.tar.gz

# Verificar restore
./bin/syntropy manager list
```

## 📚 Próximos Passos

1. **Explore templates**: Crie templates personalizados para suas aplicações
2. **Automatize**: Configure scripts de monitoramento e backup automático
3. **Escale**: Adicione mais nós à sua grid cooperativa
4. **Integre**: Conecte com ferramentas de CI/CD e monitoramento
5. **Contribua**: Desenvolva novos templates e funcionalidades

