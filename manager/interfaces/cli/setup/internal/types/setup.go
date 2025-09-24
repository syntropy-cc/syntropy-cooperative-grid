package types

import (
	"errors"
	"time"
)

// SetupOptions define as opções para o processo de setup
type SetupOptions struct {
	Force          bool   // Forçar setup mesmo com validações falhas
	InstallService bool   // Instalar serviço do sistema
	ConfigPath     string // Caminho personalizado para o arquivo de configuração
	HomeDir        string // Diretório home personalizado
}

// SetupResult contém o resultado do processo de setup
type SetupResult struct {
	Success     bool      // Indica se o setup foi bem-sucedido
	StartTime   time.Time // Hora de início do setup
	EndTime     time.Time // Hora de término do setup
	ConfigPath  string    // Caminho do arquivo de configuração
	Environment string    // Ambiente (windows, linux, darwin)
	Options     SetupOptions // Opções utilizadas no setup
	Error       error     // Erro, se houver
}

// ErrNotImplemented é retornado quando uma funcionalidade não está implementada
var ErrNotImplemented = errors.New("funcionalidade não implementada para este sistema operacional")