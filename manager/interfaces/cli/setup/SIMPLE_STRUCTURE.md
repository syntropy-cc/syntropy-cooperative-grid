# Estrutura Simplificada ~/.syntropy/

## Estrutura Final
```
~/.syntropy/
├── config/
│   └── manager.yaml           # Configuração principal
├── keys/
│   ├── owner.key              # Chave privada do administrador
│   └── owner.key.pub          # Chave pública do administrador
├── nodes/                     # Nós gerenciados
│   ├── lab-raspberry-01/      # Nome do nó como pasta
│   │   ├── metadata.yaml      # Metadados do nó
│   │   ├── config.yaml        # Configuração do nó
│   │   ├── status.json        # Status atual
│   │   ├── community.key      # Chave community do nó
│   │   └── community.key.pub  # Chave pública do nó
│   └── mini-pc-02/            # Outro nó
│       ├── metadata.yaml
│       ├── config.yaml
│       ├── status.json
│       ├── community.key
│       └── community.key.pub
├── logs/
│   ├── setup.log              # Logs do setup
│   ├── manager.log            # Logs do manager
│   ├── node-creation.log      # Logs de criação de nós
│   └── security.log           # Logs de segurança
├── cache/
│   └── iso/                   # Cache de imagens ISO
└── backups/                   # Backups automáticos
    ├── config/
    ├── keys/
    └── nodes/
```

## Detalhamento

### config/
- **manager.yaml**: Configuração única do quartel geral

### keys/
- **owner.key**: Chave privada única do administrador
- **owner.key.pub**: Chave pública do administrador

### nodes/
- **Pasta por nó**: Nome da pasta = nome do nó
- **metadata.yaml**: Informações do nó
- **config.yaml**: Configuração específica do nó
- **status.json**: Status atual em tempo real
- **community.key**: Chave privada do nó
- **community.key.pub**: Chave pública do nó

### logs/
- **setup.log**: Logs do processo de setup
- **manager.log**: Logs gerais do manager
- **node-creation.log**: Logs de criação de nós
- **security.log**: Logs de segurança

### cache/
- **iso/**: Cache de imagens ISO para criação de nós

### backups/
- **config/**: Backups de configuração
- **keys/**: Backups de chaves
- **nodes/**: Backups dos nós

## Vantagens

- **Simples**: Estrutura mínima e direta
- **Claro**: Nome da pasta = nome do nó
- **Seguro**: Chaves organizadas por nó
- **Eficiente**: Cache de ISOs
- **Confiável**: Backups automáticos
