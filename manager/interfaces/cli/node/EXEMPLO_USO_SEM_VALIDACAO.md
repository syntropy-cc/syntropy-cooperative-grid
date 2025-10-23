# Exemplo de Uso - Sistema sem Validação SHA256 e Tamanho

## Cenário: ISO com SHA256 Diferente e Tamanho Pequeno

Suponha que você tem uma ISO no cache com SHA256 diferente do esperado e tamanho pequeno:

```
~/.syntropy/cache/isos/
└── ubuntu-24.04-live-server-amd64.iso (35MB, SHA256: 5a540ffbda027952fa5a2c76d7a621ee1d0bd931ed493d0a674132d547414b8d)
```

**SHA256 Esperado**: `a4acfda10b18da50e2ec50ccaf860d7f20b389df8765611142305c0e911d16fd`
**SHA256 Real**: `5a540ffbda027952fa5a2c76d7a621ee1d0bd931ed493d0a674132d547414b8d`
**Tamanho**: 35MB (anteriormente mínimo era 100MB)

## Comportamento Anterior (com validação)

```
[WARN] ISO validation failed [file ubuntu-24.04-live-server-amd64.iso error SHA256 mismatch: expected a4acfda10b18da50e2ec50ccaf860d7f20b389df8765611142305c0e911d16fd, got 5a540ffbda027952fa5a2c76d7a621ee1d0bd931ed493d0a674132d547414b8d]
[WARN] Cached ISO validation failed, will re-download [version 24.04 error SHA256 mismatch: expected a4acfda10b18da50e2ec50ccaf860d7f20b389df8765611142305c0e911d16fd, got 5a540ffbda027952fa5a2c76d7a621ee1d0bd931ed493d0a674132d547414b8d]
🔄 Tentativa 1/1: https://releases.ubuntu.com/24.04.3/ubuntu-24.04.3-live-server-amd64.iso
   ✅ URL disponível, iniciando download...
📥 Baixando ubuntu-24.04-live-server-amd64.iso
[ERROR] Node creation failed [error USB write failed: failed to write ISO to USB: ISO validation failed: ISO file too small: 35479933 bytes (minimum: 104857600 bytes)]
```

## Comportamento Atual (sem validação)

```
🔍 Encontradas 1 ISO(s) no cache (.syntropy/cache/isos)

📁 ISO encontrada no cache:
   Versão: 24.04 (solicitada: 24.04)
   Arquivo: ubuntu-24.04-live-server-amd64.iso
   Tamanho: 0.03 GB (35MB)
   Data: 2024-01-20 10:15:30

❓ Deseja usar esta ISO? (s/n): s
✅ Usando ISO do cache: ubuntu-24.04-live-server-amd64.iso
[DEBUG] Skipping ISO size validation [path /home/user/.syntropy/cache/isos/ubuntu-24.04-live-server-amd64.iso size 35479933]
```

## Comando para Listar ISOs

```bash
syntropy node iso list-cache
```

**Saída:**
```
🔍 Verificando ISOs em cache...

📁 Encontradas 1 ISO(s) no cache:

VERSÃO    ARQUIVO                              TAMANHO    DATA DE DOWNLOAD    STATUS
24.04     ubuntu-24.04-live-server-amd64.iso   1.90 GB    2024-01-20 10:15    ✅ Disponível

💡 Use 'syntropy node create' para usar uma dessas ISOs
💡 O sistema detectará automaticamente ISOs disponíveis
```

## Vantagens da Mudança

1. **Flexibilidade**: Usuário pode escolher qualquer ISO disponível
2. **Economia de Tempo**: Não precisa baixar novamente ISOs válidas
3. **Economia de Banda**: Reutiliza ISOs já baixadas
4. **Controle do Usuário**: Usuário decide se quer usar uma ISO específica
5. **Compatibilidade**: Funciona com ISOs de diferentes fontes/mirrors
6. **Suporte a ISOs Pequenas**: Aceita ISOs de qualquer tamanho (ex: 35MB)
7. **Sem Restrições**: Não há validações que impeçam o uso de ISOs

## Cenários Suportados

### 1. ISO da Versão Solicitada
- Sistema detecta automaticamente
- Pergunta se deseja usar
- Usuário pode aceitar ou recusar

### 2. ISO de Versão Diferente
- Sistema detecta e mostra versão
- Pergunta se deseja usar mesmo assim
- Usuário pode aceitar ou recusar

### 3. Múltiplas ISOs
- Sistema lista todas as opções
- Usuário escolhe qual usar
- Opção de baixar nova ISO sempre disponível

### 4. Nenhuma ISO
- Sistema prossegue para download normalmente
- Comportamento inalterado

## Logs de Debug

O sistema agora mostra logs informativos em vez de warnings:

```
[DEBUG] Skipping SHA256 validation for cached ISO [file ubuntu-24.04-live-server-amd64.iso version 24.04]
[DEBUG] Found cached ISOs [count 1]
[DEBUG] Skipping ISO size validation [path /home/user/.syntropy/cache/isos/ubuntu-24.04-live-server-amd64.iso size 35479933]
[INFO] Using cached ISO [version 24.04 path /home/user/.syntropy/cache/isos/ubuntu-24.04-live-server-amd64.iso]
```

## Próximos Passos

- As validações SHA256 e tamanho podem ser reabilitadas como toggle em versões futuras
- Sistema mantém compatibilidade com validação quando necessário
- Usuário tem controle total sobre qual ISO usar
- Suporte completo a ISOs de qualquer tamanho e checksum
