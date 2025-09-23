# Syntropy Cooperative Grid - Infrastructure as Code

Este diretório contém toda a infraestrutura como código (IaC) para a Syntropy Cooperative Grid, organizando templates, configurações e automação para deployment e gerenciamento de nós.

## 📁 Estrutura do Diretório

```
infrastructure/
├── cloud-init/              # Templates de configuração cloud-init
│   ├── user-data-template.yaml
│   ├── meta-data-template.yaml
│   └── network-config-template.yaml
├── packer/                  # Templates Packer para construção de imagens
│   ├── syntropy-node-base.pkr.hcl
│   └── http/                # Arquivos para servidor HTTP do Packer
├── terraform/               # Configurações Terraform
│   └── syntropy-nodes/
│       ├── main.tf
│       └── ansible-inventory.tpl
├── ansible/                 # Playbooks e configurações Ansible
│   ├── playbooks/
│   │   └── syntropy-node-setup.yml
│   └── inventory/
└── output/                  # Arquivos gerados (gitignored)
```

## 🚀 Cloud-Init Templates

### user-data-template.yaml
Template principal para configuração automática do Ubuntu Server durante a instalação. Inclui:
- Configuração de rede automática
- Instalação de pacotes essenciais
- Configuração de usuários e SSH
- Setup do Docker e serviços
- Configuração de firewall e segurança
- Criação de estrutura de diretórios Syntropy

### meta-data-template.yaml
Metadados para identificação do nó, incluindo:
- Instance ID único
- Hostname do nó
- Chaves públicas SSH
- Configuração de rede

### network-config-template.yaml
Configuração avançada de rede com:
- Auto-detecção de interfaces
- Configuração DHCP
- Servidores DNS
- Configuração de loopback

## 🏗️ Packer Templates

### syntropy-node-base.pkr.hcl
Template para construção de imagens base do Ubuntu Server com:
- Configuração automática via cloud-init
- Instalação de pacotes essenciais
- Configuração de serviços
- Preparação para deployment

**Uso:**
```bash
cd infrastructure/packer
packer build syntropy-node-base.pkr.hcl
```

## 🌍 Terraform

### syntropy-nodes/main.tf
Configuração para provisionamento de nós em cloud providers:
- AWS EC2 instances
- Security groups
- Key pairs
- Outputs para Ansible

**Uso:**
```bash
cd infrastructure/terraform/syntropy-nodes
terraform init
terraform plan
terraform apply
```

## 🔧 Ansible

### syntropy-node-setup.yml
Playbook completo para configuração de nós:
- Instalação de pacotes
- Configuração de Docker
- Setup de SSH e segurança
- Configuração de firewall
- Setup de monitoramento
- Criação de templates Kubernetes
- Configuração de metadados

**Uso:**
```bash
cd infrastructure/ansible
ansible-playbook -i inventory/hosts.yml playbooks/syntropy-node-setup.yml
```

## 🔄 Workflow de Deployment

### 1. Preparação de Imagens Base (Packer)
```bash
# Construir imagem base
cd infrastructure/packer
packer build syntropy-node-base.pkr.hcl
```

### 2. Provisionamento de Infraestrutura (Terraform)
```bash
# Provisionar nós
cd infrastructure/terraform/syntropy-nodes
terraform init
terraform apply -var="node_count=3"
```

### 3. Configuração de Nós (Ansible)
```bash
# Configurar nós provisionados
cd infrastructure/ansible
ansible-playbook -i ../terraform/syntropy-nodes/inventory/hosts.yml playbooks/syntropy-node-setup.yml
```

## 🔐 Gerenciamento de Chaves

O sistema utiliza dois tipos de chaves SSH:

### Owner Keys (Chaves do Proprietário)
- **Propósito**: Acesso SSH e gerenciamento do nó
- **Algoritmo**: ED25519 (recomendado) ou RSA 4096-bit
- **Localização**: `/opt/syntropy/identity/owner/`
- **Uso**: Login SSH, administração

### Community Keys (Chaves da Comunidade)
- **Propósito**: Comunicação entre nós
- **Algoritmo**: ED25519
- **Localização**: `/opt/syntropy/identity/community/`
- **Uso**: Mesh networking, autenticação entre nós

## 📊 Monitoramento

### Prometheus Node Exporter
- **Porta**: 9100
- **Coletas**: CPU, memória, disco, rede, sistema
- **Configuração**: Automática via cloud-init e Ansible

### Health Checks
- Docker service status
- SSH service status
- Disk space monitoring
- Memory usage monitoring

## 🛡️ Segurança

### Firewall (UFW)
- **Política padrão**: Deny incoming, Allow outgoing
- **Portas abertas**: SSH (22), Prometheus (9100)
- **Configuração**: Automática via cloud-init

### Fail2ban
- **Proteção SSH**: 3 tentativas, ban de 1 hora
- **Logs monitorados**: `/var/log/auth.log`
- **Configuração**: Automática via Ansible

### SSH
- **Autenticação**: Apenas chaves (sem senha)
- **Usuário**: admin
- **Porta**: 22 (configurável)

## 🔧 Configuração Avançada

### Variáveis de Ambiente
```bash
export SYNTROPY_NODE_COUNT=5
export SYNTROPY_INSTANCE_TYPE=t3.medium
export SYNTROPY_COORDINATES="-23.5505,-46.6333"
export SYNTROPY_NODE_DESCRIPTION="Production Node"
```

### Customização de Templates
Todos os templates suportam variáveis personalizáveis:
- `{{.NodeName}}` - Nome do nó
- `{{.Coordinates}}` - Coordenadas geográficas
- `{{.NodeDescription}}` - Descrição do nó
- `{{.CreatedAt}}` - Timestamp de criação

## 📝 Logs e Debugging

### Logs do Sistema
```bash
# Logs do cloud-init
sudo journalctl -u cloud-init

# Logs do Ansible
ansible-playbook -v playbooks/syntropy-node-setup.yml

# Logs do Terraform
terraform apply -auto-approve
```

### Verificação de Status
```bash
# Verificar status dos serviços
systemctl status docker ssh prometheus-node-exporter

# Verificar configuração Syntropy
cat /opt/syntropy/config/platform.yaml

# Verificar metadados do nó
cat /opt/syntropy/metadata/node.json
```

## 🔄 Atualizações e Manutenção

### Atualização de Templates
1. Modificar templates em `infrastructure/`
2. Re-executar workflow de deployment
3. Aplicar configurações via Ansible

### Backup de Configurações
```bash
# Backup de configurações Syntropy
tar -czf syntropy-config-backup.tar.gz /opt/syntropy/

# Backup de chaves
cp -r /opt/syntropy/identity/ ~/syntropy-keys-backup/
```

## 🆘 Troubleshooting

### Problemas Comuns

**Cloud-init não executa:**
```bash
sudo cloud-init status
sudo cloud-init analyze show
```

**Ansible falha na conexão:**
```bash
ansible all -m ping -i inventory/hosts.yml
```

**Terraform state lock:**
```bash
terraform force-unlock <lock-id>
```

**Serviços não iniciam:**
```bash
sudo systemctl status <service>
sudo journalctl -u <service>
```

## 📚 Referências

- [Cloud-Init Documentation](https://cloudinit.readthedocs.io/)
- [Packer Documentation](https://www.packer.io/docs)
- [Terraform Documentation](https://www.terraform.io/docs)
- [Ansible Documentation](https://docs.ansible.com/)
- [Ubuntu Server Guide](https://ubuntu.com/server/docs)

## 🤝 Contribuição

Para contribuir com a infraestrutura:

1. Fork do repositório
2. Criar branch para feature
3. Modificar templates e configurações
4. Testar com ambiente de desenvolvimento
5. Submeter pull request

### Testes
```bash
# Testar templates cloud-init
cloud-init devel schema --config-file cloud-init/user-data-template.yaml

# Validar Terraform
terraform validate

# Testar Ansible
ansible-playbook --check playbooks/syntropy-node-setup.yml
```
