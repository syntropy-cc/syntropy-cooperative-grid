# Exemplos de Uso - Syntropy Cooperative Grid USB Creator

Este documento contém exemplos práticos de como usar o Syntropy Cooperative Grid USB Creator.

## 🚀 Início Rápido

### 1. Compilação e Instalação

```bash
# Compilar o projeto
make build

# Verificar dependências do sistema
make check-deps

# Instalar no sistema
make install
```

### 2. Verificar Ambiente

```bash
# Verificar se está no WSL
make check-wsl

# Verificar dispositivos USB disponíveis
make check-usb

# Listar dispositivos USB
syntropy usb list
```

## 📱 Comandos Básicos

### Detecção de Dispositivos USB

```bash
# Listar todos os dispositivos USB
syntropy usb list

# Listar em formato JSON
syntropy usb list --format json

# Listar em formato YAML
syntropy usb list --format yaml
```

**Saída esperada:**
```
DISPOSITIVO   TAMANHO   MODELO                FABRICANTE       REMOVÍVEL   PLATAFORMA
--------------------------------------------------------------------------------
/dev/sdb      32G       SanDisk Ultra        SanDisk          Sim         linux
/dev/sdc      16G       Kingston DataTraveler Kingston        Sim         linux
```

### Criação de USB com Boot

```bash
# Criação básica com auto-detecção
syntropy usb create --auto-detect --node-name "servidor-01"

# Criação especificando dispositivo
syntropy usb create /dev/sdb --node-name "servidor-01"

# Criação com descrição e coordenadas
syntropy usb create /dev/sdb \
  --node-name "servidor-01" \
  --description "Servidor de produção principal" \
  --coordinates "-23.5505,-46.6333"
```

### Formatação de USB

```bash
# Formatação básica
syntropy usb format /dev/sdb

# Formatação com rótulo personalizado
syntropy usb format /dev/sdb --label "MEU_USB"

# Formatação sem confirmação (cuidado!)
syntropy usb format /dev/sdb --force
```

## 🏗️ Workflows Completos

### Workflow 1: Primeiro Nó da Grid

```bash
# 1. Verificar dispositivos disponíveis
syntropy usb list

# 2. Criar USB para nó principal
syntropy usb create --auto-detect \
  --node-name "grid-master" \
  --description "Nó principal da Syntropy Cooperative Grid" \
  --coordinates "-23.5505,-46.6333"

# 3. Verificar criação
ls -la ~/.syntropy/nodes/

# 4. Conectar ao nó (após instalação)
ssh -i ~/.syntropy/keys/grid-master_owner.key admin@<IP_DO_NÓ>
```

### Workflow 2: Múltiplos Nós

```bash
# 1. Criar nó principal
syntropy usb create /dev/sdb \
  --node-name "master-node" \
  --description "Nó mestre da grid"

# 2. Criar nós de borda
syntropy usb create /dev/sdc \
  --node-name "edge-node-01" \
  --description "Nó de borda 01"

syntropy usb create /dev/sdd \
  --node-name "edge-node-02" \
  --description "Nó de borda 02"

# 3. Listar todos os nós criados
syntropy node list
```

### Workflow 3: Deploy de Containers

```bash
# 1. Listar containers
syntropy container list

# 2. Deploy de nginx
syntropy container deploy nginx \
  --node master-node \
  --port "8080:80"

# 3. Deploy com escala
syntropy container deploy nginx \
  --node master-node \
  --port "8080:80" \
  --scale 3

# 4. Verificar logs
syntropy container logs nginx-01 --follow
```

## 🔧 Configurações Avançadas

### Usando Arquivo de Configuração

```bash
# Copiar configuração de exemplo
cp config.yaml.example ~/.syntropy/config.yaml

# Editar configurações
nano ~/.syntropy/config.yaml

# Usar configuração personalizada
syntropy usb create --config ~/.syntropy/config.yaml
```

### Configuração de Rede

```bash
# Habilitar service mesh
syntropy network mesh enable --encryption --monitoring

# Criar rota entre nós
syntropy network routes create \
  --source master-node \
  --dest edge-node-01 \
  --priority 1

# Verificar topologia
syntropy network topology --format graphviz
```

### Gerenciamento Cooperativo

```bash
# Ver saldo de créditos
syntropy cooperative credits balance --node master-node

# Transferir créditos
syntropy cooperative credits transfer \
  --from master-node \
  --to edge-node-01 \
  --amount 100

# Votar em proposta
syntropy cooperative governance vote \
  --proposal prop-001 \
  --vote yes
```

## 🐛 Solução de Problemas

### Problemas de Detecção USB

```bash
# Verificar dispositivos manualmente
lsblk

# Verificar se está no WSL
cat /proc/version | grep -i microsoft

# Reiniciar WSL (do PowerShell)
wsl --shutdown

# Verificar permissões
sudo fdisk -l /dev/sdb
```

### Problemas de Formatação

```bash
# Desmontar partições manualmente
sudo umount /dev/sdb*

# Limpar assinaturas
sudo wipefs -a /dev/sdb

# Verificar se dispositivo está em uso
sudo fuser /dev/sdb

# Matar processos usando o dispositivo
sudo fuser -k /dev/sdb
```

### Problemas de Conectividade

```bash
# Verificar status do nó
syntropy node status master-node

# Verificar saúde da rede
syntropy network health --detailed

# Testar conectividade SSH
ssh -i ~/.syntropy/keys/master-node_owner.key admin@<IP_DO_NÓ>
```

## 📊 Monitoramento

### Verificação de Status

```bash
# Status de todos os nós
syntropy node status

# Status de um nó específico
syntropy node status master-node --watch

# Status em formato JSON
syntropy node status master-node --format json
```

### Logs e Métricas

```bash
# Logs de container
syntropy container logs nginx-01 --follow --tail 100

# Logs do sistema (no nó)
ssh -i ~/.syntropy/keys/master-node_owner.key admin@<IP_DO_NÓ> \
  "journalctl -u syntropy-first-boot -f"

# Métricas do nó
ssh -i ~/.syntropy/keys/master-node_owner.key admin@<IP_DO_NÓ> \
  "curl http://localhost:9100/metrics"
```

## 🔒 Segurança

### Gerenciamento de Chaves

```bash
# Listar chaves SSH
ls -la ~/.syntropy/keys/

# Verificar permissões
ls -la ~/.syntropy/keys/*.key

# Backup de chaves
tar -czf ~/syntropy-keys-backup.tar.gz ~/.syntropy/keys/
```

### Firewall e Acesso

```bash
# Verificar regras de firewall
ssh -i ~/.syntropy/keys/master-node_owner.key admin@<IP_DO_NÓ> \
  "sudo ufw status"

# Configurar acesso restrito
ssh -i ~/.syntropy/keys/master-node_owner.key admin@<IP_DO_NÓ> \
  "sudo ufw allow from 192.168.1.0/24 to any port 22"
```

## 🚀 Automação

### Script de Deploy Automático

```bash
#!/bin/bash
# deploy-grid.sh

set -e

# Configurações
MASTER_NODE="master-node"
EDGE_NODES=("edge-node-01" "edge-node-02" "edge-node-03")
USB_DEVICES=("/dev/sdb" "/dev/sdc" "/dev/sdd")

# Criar nó mestre
echo "Criando nó mestre..."
syntropy usb create ${USB_DEVICES[0]} \
  --node-name "$MASTER_NODE" \
  --description "Nó mestre da grid"

# Criar nós de borda
for i in "${!EDGE_NODES[@]}"; do
  device_index=$((i + 1))
  echo "Criando ${EDGE_NODES[$i]}..."
  syntropy usb create ${USB_DEVICES[$device_index]} \
    --node-name "${EDGE_NODES[$i]}" \
    --description "Nó de borda $(($i + 1))"
done

echo "Grid criada com sucesso!"
syntropy node list
```

### Monitoramento Contínuo

```bash
#!/bin/bash
# monitor-grid.sh

while true; do
  echo "=== Status da Grid $(date) ==="
  
  # Status dos nós
  syntropy node status
  
  # Status da rede
  syntropy network health
  
  # Status dos containers
  syntropy container list
  
  echo "Aguardando 60 segundos..."
  sleep 60
done
```

## 📚 Recursos Adicionais

### Documentação

- [README Principal](../README.md)
- [Configuração](../config.yaml.example)
- [Makefile](../Makefile)

### Comandos de Ajuda

```bash
# Ajuda geral
syntropy --help

# Ajuda de comando específico
syntropy usb --help
syntropy usb create --help
syntropy node --help
syntropy container --help
syntropy network --help
syntropy cooperative --help
```

### Logs do Sistema

```bash
# Logs do CLI
tail -f ~/.syntropy/logs/syntropy.log

# Logs do sistema (Ubuntu)
ssh -i ~/.syntropy/keys/master-node_owner.key admin@<IP_DO_NÓ> \
  "journalctl -f"

# Logs de cloud-init
ssh -i ~/.syntropy/keys/master-node_owner.key admin@<IP_DO_NÓ> \
  "sudo tail -f /var/log/cloud-init-output.log"
```

---

**Nota**: Sempre teste em ambiente de desenvolvimento antes de usar em produção. Use dispositivos USB de teste para evitar perda de dados importantes.
