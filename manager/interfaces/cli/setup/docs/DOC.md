# Setup Component

## What is This?

[MACRO_VIEW]
O componente Setup é responsável por configurar e inicializar o ambiente Syntropy CLI em diferentes sistemas operacionais, garantindo que todos os componentes necessários estejam prontos para funcionamento no ecossistema Syntropy Cooperative Grid.
[/MACRO_VIEW]

[MESO_VIEW]
Este componente funciona como a interface de configuração inicial do CLI Manager, interagindo com outros componentes do sistema para validar o ambiente, gerar configurações e estabelecer a estrutura necessária para operações subsequentes.
[/MESO_VIEW]

[MICRO_VIEW]
O componente Setup resolve o problema de inicialização e configuração automática do ambiente Syntropy, incluindo validação de dependências, geração de chaves criptográficas, criação de estruturas de diretórios e configuração de serviços do sistema.
[/MICRO_VIEW]

## Why Use This?

### Problems It Solves
- **Configuração Manual Complexa**: Elimina a necessidade de configuração manual detalhada do ambiente Syntropy
- **Incompatibilidade de Sistemas**: Resolve problemas de compatibilidade entre diferentes sistemas operacionais
- **Gerenciamento de Chaves**: Automatiza a geração e gerenciamento de chaves criptográficas necessárias
- **Validação de Ambiente**: Verifica automaticamente se o sistema atende aos requisitos mínimos

### Key Benefits
- **Configuração Automática**: Setup completo em uma única operação
- **Multiplataforma**: Suporte nativo para Windows, Linux e macOS
- **Backup Automático**: Cria backups automáticos de configurações existentes
- **Validação Inteligente**: Verifica ambiente, dependências e permissões
- **Logging Estruturado**: Registra todas as operações para auditoria e debugging

## Quick Start

### Prerequisites
- Go 1.21 ou superior
- Permissões de administrador (para instalação de serviços)
- 100MB de espaço em disco livre
- Conexão com internet (para validações)

### Installation

```bash
# Clonar o repositório
git clone https://github.com/syntropy-cc/syntropy-cooperative-grid.git
cd syntropy-cooperative-grid/manager/interfaces/cli/setup

# Compilar o componente
go build -o setup src/setup.go
```

### Basic Usage

```go
// Exemplo básico de uso
import "setup-component"

// Criar opções de setup
options := &types.SetupOptions{
    Force:        false,
    Verbose:      true,
    CustomSettings: map[string]string{
        "owner_name":  "João Silva",
        "owner_email": "joao@syntropy.network",
    },
}

// Executar setup
manager, err := setup.NewSetupManager()
if err != nil {
    log.Fatal(err)
}

err = manager.Setup(options)
if err != nil {
    log.Fatal("Setup falhou:", err)
}

fmt.Println("Setup concluído com sucesso!")
```

**Input**: Configurações básicas do usuário
**Output**: Ambiente Syntropy completamente configurado e pronto para uso

## Features

| Feature | Description | Status |
|---------|-------------|--------|
| Validação de Ambiente | Verifica SO, permissões e dependências | Stable |
| Geração de Chaves | Cria pares de chaves Ed25519 automaticamente | Stable |
| Configuração Automática | Gera arquivos de configuração necessários | Stable |
| Backup Inteligente | Preserva configurações existentes | Stable |
| Logging Estruturado | Registra todas as operações detalhadamente | Stable |
| Reset Seguro | Remove configurações com confirmação | Stable |
| Validação de Integridade | Verifica integridade de arquivos e chaves | Beta |
| Instalação de Serviços | Configura serviços do sistema (Windows/Linux) | Beta |

## Component Structure

```
setup/
├── docs/           # Documentação completa
├── src/            # Código fonte principal
│   ├── setup.go    # Interface principal
│   ├── logger.go   # Sistema de logging
│   ├── validator.go # Validação de ambiente
│   ├── configurator.go # Geração de configurações
│   ├── key_manager.go # Gerenciamento de chaves
│   └── state_manager.go # Gerenciamento de estado
├── internal/       # Tipos e utilitários internos
├── tests/          # Testes automatizados
├── examples/       # Exemplos de uso
└── config/         # Arquivos de configuração
```

## Next Steps

- [Explore the API](./API.md) - Instruções detalhadas de uso
- [Developer Guide](./DEV.md) - Entendendo a arquitetura interna
- [Testing Guide](./TEST.md) - Executando e escrevendo testes
- [Learning Path](./LEARN.md) - Mergulho profundo em conceitos e teoria
- [Examples](../examples/) - Mais cenários de uso

## Support

- Issue Tracker: [GitHub Issues](https://github.com/syntropy-cc/syntropy-cooperative-grid/issues)
- Discussion Forum: [Syntropy Community](https://community.syntropy.network)
- Contact: [Documentação Oficial](https://docs.syntropy.network)

## License

MIT License - Veja o arquivo LICENSE para detalhes








