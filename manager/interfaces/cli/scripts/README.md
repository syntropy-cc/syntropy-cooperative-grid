# Syntropy CLI Manager - Scripts Directory

Esta pasta contém scripts unificados para compilação e teste em todas as plataformas suportadas.

## 📁 Nova Estrutura Simplificada

```
scripts/
├── build-all.sh      # Script universal para Linux/macOS
├── build-all.ps1     # Script universal para Windows PowerShell
├── build.sh          # Runner universal para Linux/macOS
├── build.bat         # Runner universal para Windows
├── test-build.sh     # Script de teste rápido
├── windows/          # Scripts legados específicos para Windows
├── linux/            # Scripts legados específicos para Linux
├── shared/           # Scripts legados compartilhados
└── README.md         # Este arquivo
```

## 🚀 Scripts Principais (Recomendados)

### Scripts Universais
- **`build.sh`** / **`build.bat`** - Runners universais que detectam a plataforma automaticamente
- **`build-all.sh`** - Script principal para Linux/macOS com suporte completo
- **`build-all.ps1`** - Script principal para Windows PowerShell com suporte completo
- **`test-build.sh`** - Script de teste rápido para validar a instalação

## 🎯 Uso Rápido (Recomendado)

### Para Qualquer Plataforma
```bash
# Linux/macOS
./scripts/build.sh

# Windows
scripts\build.bat
```

### Opções Avançadas
```bash
# Linux/macOS
./scripts/build-all.sh --help                    # Ver todas as opções
./scripts/build-all.sh --current                 # Build apenas para plataforma atual
./scripts/build-all.sh --platform windows/amd64  # Build para plataforma específica
./scripts/build-all.sh --test                    # Executar apenas testes
./scripts/build-all.sh --run                     # Executar aplicação após build

# Windows PowerShell
.\scripts\build-all.ps1 help                     # Ver todas as opções
.\scripts\build-all.ps1 current                  # Build apenas para plataforma atual
.\scripts\build-all.ps1 platform linux/amd64     # Build para plataforma específica
.\scripts\build-all.ps1 test                     # Executar apenas testes
.\scripts\build-all.ps1 run                      # Executar aplicação após build
```

### Teste Rápido
```bash
# Validar instalação
./scripts/test-build.sh
```

## 🖥️ Plataformas Suportadas

### Build Cross-Platform
- **Linux**: amd64, arm64
- **Windows**: amd64
- **macOS**: amd64, arm64 (Apple Silicon)

### Funcionalidades
- ✅ Compilação para múltiplas plataformas simultaneamente
- ✅ Testes automatizados (unitários, cobertura, race conditions)
- ✅ Detecção automática de plataforma
- ✅ Scripts universais (funcionam em qualquer SO)
- ✅ Build otimizado com flags de versão e commit
- ✅ Validação de binários gerados

## 📚 Scripts Legados (Compatibilidade)

### Windows Scripts (`windows/`)
- **`build-windows.ps1`** - Script legado de build e execução
- **`dev-workflow.ps1`** - Workflow legado de desenvolvimento
- **`automation-workflow.ps1`** - Workflow legado de automação

### Linux Scripts (`linux/`)
- **`install-syntropy.sh`** - Script legado de instalação
- **`build-and-test.sh`** - Script legado de build e teste

### Shared Scripts (`shared/`)
- **`build-and-test.bat`** - Script legado compartilhado
- **`start-here.bat`** - Script legado de entrada

## 🎯 Workflows Recomendados

### Para Iniciantes (Mais Simples)
```bash
# Qualquer plataforma - detecta automaticamente
./scripts/build.sh          # Linux/macOS
scripts\build.bat           # Windows
```

### Para Desenvolvimento
```bash
# Build apenas para plataforma atual (mais rápido)
./scripts/build-all.sh --current
.\scripts\build-all.ps1 current
```

### Para Distribuição
```bash
# Build para todas as plataformas
./scripts/build-all.sh
.\scripts\build-all.ps1 all
```

### Para Testes
```bash
# Executar apenas testes
./scripts/build-all.sh --test
.\scripts\build-all.ps1 test
```

### Para Validação
```bash
# Teste rápido da instalação
./scripts/test-build.sh
```

## 🔧 Funcionalidades Avançadas

### Build Cross-Platform
- ✅ Compilação simultânea para 5 plataformas
- ✅ Detecção automática de arquitetura
- ✅ Flags de build otimizadas
- ✅ Informações de versão e commit Git

### Testes Automatizados
- ✅ Testes unitários com verbose
- ✅ Testes de cobertura
- ✅ Testes de race conditions
- ✅ Validação de binários gerados

### Qualidade e Confiabilidade
- ✅ Verificação de pré-requisitos
- ✅ Validação de dependências
- ✅ Logs coloridos e informativos
- ✅ Tratamento de erros robusto

### Usabilidade
- ✅ Interface unificada para todas as plataformas
- ✅ Detecção automática de SO
- ✅ Mensagens de ajuda detalhadas
- ✅ Exemplos de uso integrados

## 🛠️ Manutenção

### Adicionar Novas Plataformas
1. Edite `build-all.sh` e `build-all.ps1`
2. Adicione a nova plataforma ao array `PLATFORMS`
3. Teste em todas as plataformas suportadas
4. Atualize este README

### Modificar Scripts
1. **Scripts principais**: `build-all.sh`, `build-all.ps1`
2. **Runners**: `build.sh`, `build.bat`
3. **Teste**: `test-build.sh`
4. Sempre teste em múltiplas plataformas

### Scripts Legados
- Os scripts em `windows/`, `linux/`, `shared/` são mantidos para compatibilidade
- Novos desenvolvimentos devem usar os scripts universais
- Considere migrar funcionalidades úteis para os scripts principais

## 🐛 Solução de Problemas

### Problemas Comuns
1. **"Permission denied"**: Execute `chmod +x scripts/*.sh`
2. **"Go not found"**: Instale Go 1.22+ e adicione ao PATH
3. **"main.go not found"**: Execute do diretório correto
4. **Build fails**: Verifique dependências com `go mod tidy`

### Logs e Debug
- Todos os scripts geram logs coloridos e informativos
- Use `--help` para ver opções disponíveis
- Execute `./scripts/test-build.sh` para diagnóstico

### Suporte por Plataforma
- **Linux/macOS**: Use `build-all.sh` diretamente
- **Windows**: Use `build-all.ps1` ou `build.bat`
- **WSL**: Use scripts Linux (`build-all.sh`)

## 📞 Suporte

Para problemas:
1. Execute `./scripts/test-build.sh` para diagnóstico
2. Verifique os logs coloridos dos scripts
3. Consulte a documentação específica de cada script
4. Teste com `--help` para ver opções disponíveis

---

**Sistema de build unificado e otimizado!** 🚀

**Principais benefícios:**
- ✅ **Simplicidade**: Um comando para todas as plataformas
- ✅ **Eficiência**: Build simultâneo para múltiplas plataformas  
- ✅ **Confiabilidade**: Testes automatizados e validação
- ✅ **Manutenibilidade**: Código unificado e bem documentado

