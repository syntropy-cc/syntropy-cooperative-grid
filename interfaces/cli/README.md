# Syntropy Cooperative Grid CLI

Interface completa de linha de comando para gerenciar a Syntropy Cooperative Grid, incluindo setup do ambiente, gerenciamento de nós, criação de USB bootável, deploy de aplicações e operações de rede cooperativa.

## 🚀 Instalação

### Pré-requisitos

- Go 1.21 ou superior
- Sistema operacional: Linux, Windows ou WSL
- Permissões de administrador (para operações USB)

### Compilação

```bash
# Clonar o repositório
git clone https://github.com/syntropy-cc/syntropy-cooperative-grid.git
cd syntropy-cooperative-grid/interfaces/cli

# Instalar dependências
make deps

# Compilar
make build
```

### Instalação Global

```bash
# Instalar globalmente
make install

# Desinstalar
make uninstall
```

## ⚙️ Setup Inicial

### Configurar Ambiente de Gerenciamento

```bash
# Setup completo do ambiente de gerenciamento
syntropy setup

# Verificar status do ambiente
syntropy manager health
```

## 🖥️ Comandos de Gerenciamento

### Listar e Conectar a Nós

```bash
# Listar todos os nós gerenciados
syntropy manager list

# Conectar a um nó específico
syntropy manager connect node-01

# Status detalhado de um nó
syntropy manager status node-01

# Status de todos os nós
syntropy manager status
```

### Descoberta de Rede

```bash
# Descobrir nós na rede local
syntropy manager discover

# Descobrir em redes específicas
syntropy manager discover --networks "192.168.1.0/24,10.0.0.0/24"

# Descoberta com configurações customizadas
syntropy manager discover --port 2222 --timeout 15 --parallel 10
```

### Backup e Restore

```bash
# Backup completo das configurações
syntropy manager backup

# Backup com arquivo específico
syntropy manager backup --output /path/to/backup.tar.gz

# Restore de backup
syntropy manager restore backup_20240115_143022.tar.gz
```

## 📋 Comandos de Templates

### Gerenciar Templates de Aplicação

```bash
# Listar templates disponíveis
syntropy templates list

# Mostrar detalhes de um template
syntropy templates show fortran-computation

# Deploy de template para um nó
syntropy templates deploy jupyter-lab --node node-01

# Deploy com valores customizados
syntropy templates deploy python-datascience --node node-01 --set "memory=2Gi" --set "cpu=1000m"

# Deploy em modo dry-run
syntropy templates deploy nginx --node node-01 --dry-run

# Criar novo template
syntropy templates create my-app --category web --description "My custom application"
```

## 📱 Comandos USB

### Listar Dispositivos USB

```bash
# Listar em formato tabela
syntropy usb list

# Listar em formato JSON
syntropy usb list --format json

# Listar em formato YAML
syntropy usb list --format yaml
```

### Criar USB com Boot

```bash
# Criar USB com auto-detecção
syntropy usb create --auto-detect --node-name "node-01"

# Criar USB especificando dispositivo
syntropy usb create /dev/sdb --node-name "node-01"

# Criar USB com coordenadas específicas
syntropy usb create /dev/sdb --node-name "node-01" --coordinates "-23.5505,-46.6333"

# Criar USB usando chave de proprietário existente
syntropy usb create /dev/sdb --node-name "node-02" --owner-key ~/.syntropy/keys/main.key

# Criar USB com descrição personalizada
syntropy usb create /dev/sdb --node-name "node-03" --description "Nó de borda em São Paulo"
```

### Formatar USB

```bash
# Formatar USB
syntropy usb format /dev/sdb

# Formatar USB com rótulo personalizado
syntropy usb format /dev/sdb --label "MYUSB"

# Formatar USB sem confirmação
syntropy usb format /dev/sdb --force
```

## 🔧 Comandos de Nós

### Listar Nós

```bash
# Listar todos os nós
syntropy node list

# Listar nós em formato JSON
syntropy node list --format json

# Filtrar nós por status
syntropy node list --filter running
```

### Criar Nó

```bash
# Criar nó com USB
syntropy node create --usb /dev/sdb --name "node-01"

# Criar nó com auto-detecção
syntropy node create --auto-detect --name "node-02"

# Criar nó com descrição
syntropy node create --usb /dev/sdb --name "node-03" --description "Nó de produção"
```

### Gerenciar Nós

```bash
# Ver status do nó
syntropy node status node-01

# Atualizar configuração
syntropy node update node-01 --name "node-01-updated"

# Reiniciar nó
syntropy node restart node-01

# Deletar nó
syntropy node delete node-01
```

## 🐳 Comandos de Containers

### Listar Containers

```bash
# Listar containers
syntropy container list

# Listar containers por nó
syntropy container list --node node-01
```

### Deploy de Containers

```bash
# Deploy de container simples
syntropy container deploy nginx --node node-01

# Deploy com configurações personalizadas
syntropy container deploy nginx --node node-01 --port 8080:80 --env "ENV=production"

# Deploy com volume
syntropy container deploy postgres --node node-01 --volume "pgdata:/var/lib/postgresql/data"
```

### Gerenciar Containers

```bash
# Iniciar container
syntropy container start container-id

# Parar container
syntropy container stop container-id

# Reiniciar container
syntropy container restart container-id

# Ver logs
syntropy container logs container-id

# Remover container
syntropy container remove container-id
```

## 🌐 Comandos de Rede

### Gerenciar Mesh

```bash
# Habilitar mesh com criptografia
syntropy network mesh enable --encryption

# Desabilitar mesh
syntropy network mesh disable

# Ver status do mesh
syntropy network mesh status
```

### Gerenciar Rotas

```bash
# Listar rotas
syntropy network routes list

# Adicionar rota
syntropy network routes add --from node-01 --to node-02

# Remover rota
syntropy network routes remove route-id
```

### Topologia e Saúde

```bash
# Ver topologia da rede
syntropy network topology

# Verificar saúde da rede
syntropy network health

# Testar conectividade
syntropy network test --from node-01 --to node-02
```

## 🤝 Comandos Cooperativos

### Créditos

```bash
# Ver saldo de créditos
syntropy cooperative credits balance

# Ver histórico de transações
syntropy cooperative credits history

# Transferir créditos
syntropy cooperative credits transfer --to node-02 --amount 100
```

### Governança

```bash
# Listar propostas
syntropy cooperative governance proposals

# Votar em proposta
syntropy cooperative governance vote --proposal-id 123 --vote yes

# Ver resultados
syntropy cooperative governance results --proposal-id 123
```

### Reputação

```bash
# Ver reputação do nó
syntropy cooperative reputation show --node node-01

# Ver ranking de reputação
syntropy cooperative reputation ranking

# Ver histórico de reputação
syntropy cooperative reputation history --node node-01
```

## ⚙️ Configuração

### Arquivo de Configuração

Crie um arquivo `config.yaml` em `~/.syntropy/`:

```yaml
# Configurações globais
global:
  log_level: info
  cache_dir: ~/.syntropy/cache
  work_dir: /tmp/syntropy-work

# Configurações de rede
network:
  default_encryption: true
  mesh_port: 443
  api_port: 8080

# Configurações de USB
usb:
  default_label: "SYNTROPY"
  auto_confirm: false
  safety_checks: true

# Configurações de nós
nodes:
  default_ssh_port: 22
  default_user: syntropy
  key_size: 4096
```

### Variáveis de Ambiente

```bash
# Nível de log
export SYNTROPY_LOG_LEVEL=debug

# Diretório de cache
export SYNTROPY_CACHE_DIR=/tmp/syntropy-cache

# Diretório de trabalho
export SYNTROPY_WORK_DIR=/tmp/syntropy-work

# Chave de proprietário padrão
export SYNTROPY_OWNER_KEY=~/.syntropy/keys/main.key
```

## 🔒 Segurança

### Validações de Segurança

- ✅ Verificação de discos do sistema
- ✅ Confirmação múltipla para operações destrutivas
- ✅ Validação de tamanho e tipo de dispositivo
- ✅ Geração segura de chaves SSH (RSA 4096-bit)
- ✅ Configuração automática de firewall

### Permissões

```bash
# Adicionar usuário ao grupo disk (Linux)
sudo usermod -aG disk $USER

# Adicionar usuário ao grupo docker
sudo usermod -aG docker $USER

# Reiniciar sessão para aplicar grupos
newgrp disk
newgrp docker
```

## 🐛 Solução de Problemas

### Problemas Comuns

**Erro: "dispositivo não encontrado"**
```bash
# Verificar dispositivos disponíveis
syntropy usb list

# Verificar permissões
ls -la /dev/sd*
```

**Erro: "falha ao montar"**
```bash
# Verificar se dispositivo está em uso
sudo lsof /dev/sdb

# Desmontar manualmente
sudo umount /dev/sdb*
```

**Erro: "falha ao formatar"**
```bash
# Verificar se é disco do sistema
syntropy usb list --format json | grep -i system

# Usar formatação forçada
syntropy usb format /dev/sdb --force
```

### Logs

```bash
# Ver logs detalhados
syntropy --log-level debug usb create --auto-detect --node-name "node-01"

# Logs do sistema
journalctl -u syntropy -f
```

## 🤝 Contribuição

### Desenvolvimento

```bash
# Fork do repositório
git clone https://github.com/seu-usuario/syntropy-cooperative-grid.git

# Instalar dependências de desenvolvimento
make deps

# Executar testes
make test

# Verificar código
make lint
make vet
```

### Estrutura do Projeto

```
interfaces/cli/
├── cmd/main.go                 # Ponto de entrada
├── internal/cli/               # Comandos CLI
│   ├── root.go                 # Comando raiz
│   ├── setup.go               # Setup do ambiente
│   ├── manager.go             # Gerenciamento de nós
│   ├── templates.go           # Templates de aplicação
│   ├── usb/usb.go             # Comandos USB
│   ├── node.go                # Comandos de nós
│   ├── container.go           # Comandos de containers
│   ├── network.go             # Comandos de rede
│   └── cooperative.go         # Comandos cooperativos
├── go.mod                     # Dependências Go
├── config.yaml                # Configuração padrão
├── Makefile                   # Automação de build
└── README.md                  # Documentação
```

### Estrutura de Dados

Após o setup, a CLI cria a seguinte estrutura em `~/.syntropy/`:

```
~/.syntropy/
├── nodes/                     # Metadados dos nós
│   ├── node-01.json          # Configuração do nó
│   └── node-02.json
├── keys/                      # Chaves SSH
│   ├── node-01_owner.key     # Chave privada
│   ├── node-01_owner.key.pub # Chave pública
│   └── node-02_owner.key
├── config/                    # Configuração do gerenciador
│   ├── manager.json          # Configuração principal
│   ├── templates/            # Templates de aplicação
│   │   └── applications/
│   │       ├── fortran-computation.yaml
│   │       └── python-datascience.yaml
│   └── logs/                 # Logs do sistema
├── cache/                     # Cache de descoberta
├── scripts/                   # Scripts auxiliares
│   ├── discover-network.sh
│   ├── backup-all-nodes.sh
│   └── health-check-all.sh
└── backups/                   # Backups das configurações
    └── backup_20240115_143022.tar.gz
```

## 📄 Licença

Este projeto está licenciado sob a Licença MIT. Veja o arquivo `LICENSE` para mais detalhes.

## 🆘 Suporte

- 📧 Email: support@syntropy.coop
- 💬 Discord: https://discord.gg/syntropy
- 📖 Documentação: https://docs.syntropy.coop
- 🐛 Issues: https://github.com/syntropy-cc/syntropy-cooperative-grid/issues
