# Melhoria: Seleção Interativa de Dispositivos USB

## Resumo

Esta melhoria implementa seleção interativa de dispositivos USB durante o comando `syntropy node create`, similar à funcionalidade já existente para seleção de ISO. O sistema agora:

1. **Detecta automaticamente** todos os dispositivos USB removíveis conectados
2. **Lista opções** quando múltiplos dispositivos estão disponíveis
3. **Permite seleção interativa** do dispositivo desejado
4. **Seleciona automaticamente** quando apenas um dispositivo está disponível

## Funcionalidades Implementadas

### 1. Detecção Automática de Dispositivos USB

- O sistema detecta automaticamente todos os dispositivos USB removíveis
- Filtra dispositivos do sistema para evitar sobrescrever dados importantes
- Exibe informações detalhadas sobre cada dispositivo

### 2. Seleção Interativa

#### Cenário 1: Um Dispositivo Disponível
```
🔌 Dispositivo USB encontrado:
   Caminho: /dev/sdb
   Capacidade: 8.00 GB
   Fabricante: SanDisk
   Modelo: Ultra USB 3.0
   Removível: true

✅ Usando dispositivo automaticamente
```

#### Cenário 2: Múltiplos Dispositivos Disponíveis
```
🔌 Múltiplos dispositivos USB encontrados:

  1. /dev/sdb
     Capacidade: 8.00 GB
     Fabricante: SanDisk
     Modelo: Ultra USB 3.0
     Removível: true

  2. /dev/sdc
     Capacidade: 16.00 GB
     Fabricante: Kingston
     Modelo: DataTraveler
     Removível: true

❓ Selecione um dispositivo (1-2): 2
✅ Dispositivo selecionado: /dev/sdc (16.00 GB)
```

### 3. Informações Exibidas

Para cada dispositivo, o sistema exibe:
- **Caminho**: Localização do dispositivo no sistema
- **Capacidade**: Tamanho do dispositivo em GB
- **Fabricante**: Marca do dispositivo (quando disponível)
- **Modelo**: Modelo específico do dispositivo (quando disponível)
- **Removível**: Confirmação de que o dispositivo é removível

## Implementação Técnica

### Arquivos Modificados

1. **`manager/interfaces/cli/node/src/create.go`**
   - Função `detectUSBDevice()` atualizada para usar seleção interativa
   - Nova função `selectUSBDeviceFromDetected()` implementada

2. **`manager/interfaces/cli/node/src/device_selection_test.go`**
   - Testes unitários para validar a funcionalidade

### Funções Principais

#### `detectUSBDevice(ctx context.Context) (string, error)`
- Detecta dispositivos USB removíveis
- Chama a função de seleção interativa
- Retorna o caminho do dispositivo selecionado

#### `selectUSBDeviceFromDetected(devices []types.USBDevice) (*types.USBDevice, error)`
- Gerencia a seleção interativa de dispositivos
- Exibe informações detalhadas sobre cada dispositivo
- Permite seleção manual quando múltiplos dispositivos estão disponíveis
- Seleciona automaticamente quando apenas um dispositivo está disponível

## Compatibilidade

- **Linux**: Usa `lsblk` para detectar dispositivos
- **macOS**: Usa `diskutil` para detectar dispositivos externos
- **Windows**: Usa PowerShell/WMIC para detectar dispositivos USB

## Benefícios

1. **Melhor Experiência do Usuário**: Interface clara e intuitiva
2. **Segurança**: Evita seleção acidental de dispositivos do sistema
3. **Flexibilidade**: Permite seleção manual quando necessário
4. **Automação**: Seleciona automaticamente quando apropriado
5. **Informação Detalhada**: Exibe informações úteis sobre cada dispositivo

## Exemplo de Uso

```bash
# Executar comando de criação de nó
syntropy node create

# O sistema detectará automaticamente dispositivos USB
# Se houver múltiplos dispositivos, apresentará opções para seleção
# Se houver apenas um dispositivo, será selecionado automaticamente
```

## Testes

A funcionalidade inclui testes unitários que validam:
- Seleção automática com um dispositivo
- Seleção manual com múltiplos dispositivos
- Tratamento de erro quando nenhum dispositivo está disponível
- Formatação correta das informações do dispositivo

Execute os testes com:
```bash
go test -v ./src -run TestDeviceSelection
```
