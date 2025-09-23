package usb

// Este arquivo contém apenas a função principal NewUSBCommand
// Todas as outras funcionalidades foram movidas para arquivos específicos:
//
// - types.go: Estruturas de dados
// - commands.go: Comandos CLI e função createUSB
// - platform.go: Detecção de plataforma e funções comuns
// - linux.go: Implementações específicas do Linux
// - windows.go: Implementações específicas do Windows/WSL
// - certificates.go: Geração de certificados TLS e chaves SSH
// - cloudinit.go: Configuração do cloud-init
// - utils.go: Funções auxiliares e formatação
