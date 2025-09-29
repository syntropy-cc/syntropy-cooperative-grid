package types

import (
	"fmt"
	"time"
)

// SetupError representa um erro estruturado do setup
type SetupError struct {
	Code        string                 `json:"code"`
	Message     string                 `json:"message"`
	Context     map[string]interface{} `json:"context"`
	Suggestions []string               `json:"suggestions"`
	Timestamp   time.Time              `json:"timestamp"`
	Cause       error                  `json:"-"`
}

// Códigos de erro específicos
const (
	ErrOSNotSupported      = "SETUP_001"
	ErrInsufficientPerms   = "SETUP_002"
	ErrMissingDependency   = "SETUP_003"
	ErrInsufficientSpace   = "SETUP_004"
	ErrKeyGeneration       = "SETUP_005"
	ErrNetworkConnectivity = "SETUP_006"
	ErrConfigCorrupted     = "SETUP_007"
	ErrStateCorrupted      = "SETUP_008"
	ErrBackupFailed        = "SETUP_009"
	ErrRestoreFailed       = "SETUP_010"
	ErrValidationFailed    = "SETUP_011"
	ErrConfigGeneration    = "SETUP_012"
	ErrStructureCreation   = "SETUP_013"
	ErrServiceInstall      = "SETUP_014"
	ErrKeyStorage          = "SETUP_015"
	ErrKeyRotation         = "SETUP_016"
	ErrStateSave           = "SETUP_017"
	ErrStateLoad           = "SETUP_018"
	ErrIntegrityCheck      = "SETUP_019"
	ErrTemplateProcess     = "SETUP_020"
	ErrSchemaValidation    = "SETUP_021"
)

// Error implementa a interface error
func (e *SetupError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap retorna o erro causador
func (e *SetupError) Unwrap() error {
	return e.Cause
}

// NewSetupError cria um novo erro estruturado
func NewSetupError(code, message string, cause error) *SetupError {
	return &SetupError{
		Code:        code,
		Message:     message,
		Context:     make(map[string]interface{}),
		Suggestions: []string{},
		Timestamp:   time.Now(),
		Cause:       cause,
	}
}

// WithContext adiciona contexto ao erro
func (e *SetupError) WithContext(key string, value interface{}) *SetupError {
	e.Context[key] = value
	return e
}

// WithSuggestion adiciona uma sugestão ao erro
func (e *SetupError) WithSuggestion(suggestion string) *SetupError {
	e.Suggestions = append(e.Suggestions, suggestion)
	return e
}

// WithSuggestions adiciona múltiplas sugestões ao erro
func (e *SetupError) WithSuggestions(suggestions []string) *SetupError {
	e.Suggestions = append(e.Suggestions, suggestions...)
	return e
}

// Erros pré-definidos comuns

// ErrOSNotSupportedError erro para SO não suportado
func ErrOSNotSupportedError(os string) *SetupError {
	return NewSetupError(
		ErrOSNotSupported,
		fmt.Sprintf("Sistema operacional não suportado: %s", os),
		nil,
	).WithSuggestion("Use Windows 10+, Linux (Ubuntu/Debian/CentOS/Fedora) ou macOS 10.15+")
}

// ErrInsufficientPermissionsError erro para permissões insuficientes
func ErrInsufficientPermissionsError(operation string) *SetupError {
	return NewSetupError(
		ErrInsufficientPerms,
		fmt.Sprintf("Permissões insuficientes para: %s", operation),
		nil,
	).WithSuggestions([]string{
		"Execute como administrador/root",
		"Verifique as permissões do usuário",
		"Use sudo no Linux/macOS",
	})
}

// ErrMissingDependencyError erro para dependência ausente
func ErrMissingDependencyError(dependency string) *SetupError {
	return NewSetupError(
		ErrMissingDependency,
		fmt.Sprintf("Dependência ausente: %s", dependency),
		nil,
	).WithSuggestions([]string{
		"Instale a dependência necessária",
		"Verifique se está no PATH",
		"Use o gerenciador de pacotes do seu SO",
	})
}

// ErrInsufficientSpaceError erro para espaço insuficiente
func ErrInsufficientSpaceError(required, available float64) *SetupError {
	return NewSetupError(
		ErrInsufficientSpace,
		fmt.Sprintf("Espaço em disco insuficiente: %.2f GB necessários, %.2f GB disponíveis", required, available),
		nil,
	).WithSuggestions([]string{
		"Libere espaço em disco",
		"Use um diretório diferente",
		"Remova arquivos desnecessários",
	})
}

// ErrKeyGenerationError erro para falha na geração de chaves
func ErrKeyGenerationError(algorithm string, cause error) *SetupError {
	return NewSetupError(
		ErrKeyGeneration,
		fmt.Sprintf("Falha na geração de chaves %s", algorithm),
		cause,
	).WithSuggestions([]string{
		"Verifique a fonte de entropia",
		"Tente novamente",
		"Verifique permissões de arquivo",
	})
}

// ErrNetworkConnectivityError erro para problemas de conectividade
func ErrNetworkConnectivityError(cause error) *SetupError {
	return NewSetupError(
		ErrNetworkConnectivity,
		"Falha na conectividade de rede",
		cause,
	).WithSuggestions([]string{
		"Verifique sua conexão com a internet",
		"Verifique configurações de proxy",
		"Verifique firewall",
	})
}

// ErrConfigCorruptedError erro para configuração corrompida
func ErrConfigCorruptedError(path string, cause error) *SetupError {
	return NewSetupError(
		ErrConfigCorrupted,
		fmt.Sprintf("Configuração corrompida: %s", path),
		cause,
	).WithSuggestions([]string{
		"Execute 'syntropy setup repair'",
		"Restore de um backup",
		"Execute 'syntropy setup reset'",
	})
}

// ErrStateCorruptedError erro para estado corrompido
func ErrStateCorruptedError(cause error) *SetupError {
	return NewSetupError(
		ErrStateCorrupted,
		"Estado do setup corrompido",
		cause,
	).WithSuggestions([]string{
		"Execute 'syntropy setup repair'",
		"Restore de um backup",
		"Execute 'syntropy setup reset'",
	})
}

// ErrBackupFailedError erro para falha no backup
func ErrBackupFailedError(cause error) *SetupError {
	return NewSetupError(
		ErrBackupFailed,
		"Falha na criação do backup",
		cause,
	).WithSuggestions([]string{
		"Verifique espaço em disco",
		"Verifique permissões de escrita",
		"Tente um local diferente",
	})
}

// ErrRestoreFailedError erro para falha na restauração
func ErrRestoreFailedError(cause error) *SetupError {
	return NewSetupError(
		ErrRestoreFailed,
		"Falha na restauração do backup",
		cause,
	).WithSuggestions([]string{
		"Verifique se o arquivo de backup é válido",
		"Verifique permissões de leitura",
		"Tente um backup diferente",
	})
}

// ErrValidationFailedError erro para falha na validação
func ErrValidationFailedError(issues []ValidationIssue) *SetupError {
	context := map[string]interface{}{
		"issues_count": len(issues),
		"issues":       issues,
	}

	return NewSetupError(
		ErrValidationFailed,
		"Falha na validação do ambiente",
		nil,
	).WithContext("validation_issues", context).WithSuggestion("Execute 'syntropy setup --validate-only' para detalhes")
}

// ErrConfigGenerationError erro para falha na geração de configuração
func ErrConfigGenerationError(cause error) *SetupError {
	return NewSetupError(
		ErrConfigGeneration,
		"Falha na geração de configuração",
		cause,
	).WithSuggestions([]string{
		"Verifique permissões de escrita",
		"Verifique espaço em disco",
		"Verifique templates de configuração",
	})
}

// ErrStructureCreationError erro para falha na criação de estrutura
func ErrStructureCreationError(path string, cause error) *SetupError {
	return NewSetupError(
		ErrStructureCreation,
		fmt.Sprintf("Falha na criação da estrutura de diretórios: %s", path),
		cause,
	).WithSuggestions([]string{
		"Verifique permissões de escrita",
		"Verifique espaço em disco",
		"Verifique se o diretório pai existe",
	})
}

// ErrServiceInstallError erro para falha na instalação de serviço
func ErrServiceInstallError(service string, cause error) *SetupError {
	return NewSetupError(
		ErrServiceInstall,
		fmt.Sprintf("Falha na instalação do serviço: %s", service),
		cause,
	).WithSuggestions([]string{
		"Execute como administrador/root",
		"Verifique se o serviço já existe",
		"Verifique dependências do serviço",
	})
}

// ErrKeyStorageError erro para falha no armazenamento de chaves
func ErrKeyStorageError(keyID string, cause error) *SetupError {
	return NewSetupError(
		ErrKeyStorage,
		fmt.Sprintf("Falha no armazenamento da chave: %s", keyID),
		cause,
	).WithSuggestions([]string{
		"Verifique permissões de escrita",
		"Verifique espaço em disco",
		"Verifique integridade da chave",
	})
}

// ErrKeyRotationError erro para falha na rotação de chaves
func ErrKeyRotationError(keyID string, cause error) *SetupError {
	return NewSetupError(
		ErrKeyRotation,
		fmt.Sprintf("Falha na rotação da chave: %s", keyID),
		cause,
	).WithSuggestions([]string{
		"Verifique se a chave existe",
		"Verifique permissões",
		"Verifique backup da chave anterior",
	})
}

// ErrStateSaveError erro para falha no salvamento do estado
func ErrStateSaveError(cause error) *SetupError {
	return NewSetupError(
		ErrStateSave,
		"Falha no salvamento do estado",
		cause,
	).WithSuggestions([]string{
		"Verifique permissões de escrita",
		"Verifique espaço em disco",
		"Verifique integridade do estado",
	})
}

// ErrStateLoadError erro para falha no carregamento do estado
func ErrStateLoadError(cause error) *SetupError {
	return NewSetupError(
		ErrStateLoad,
		"Falha no carregamento do estado",
		cause,
	).WithSuggestions([]string{
		"Verifique se o arquivo de estado existe",
		"Verifique permissões de leitura",
		"Verifique integridade do arquivo",
	})
}

// ErrIntegrityCheckError erro para falha na verificação de integridade
func ErrIntegrityCheckError(component string, cause error) *SetupError {
	return NewSetupError(
		ErrIntegrityCheck,
		fmt.Sprintf("Falha na verificação de integridade: %s", component),
		cause,
	).WithSuggestions([]string{
		"Execute 'syntropy setup repair'",
		"Verifique arquivos corrompidos",
		"Restore de um backup",
	})
}

// ErrTemplateProcessError erro para falha no processamento de template
func ErrTemplateProcessError(template string, cause error) *SetupError {
	return NewSetupError(
		ErrTemplateProcess,
		fmt.Sprintf("Falha no processamento do template: %s", template),
		cause,
	).WithSuggestions([]string{
		"Verifique sintaxe do template",
		"Verifique variáveis disponíveis",
		"Verifique permissões de leitura",
	})
}

// ErrSchemaValidationError erro para falha na validação de schema
func ErrSchemaValidationError(schema string, cause error) *SetupError {
	return NewSetupError(
		ErrSchemaValidation,
		fmt.Sprintf("Falha na validação do schema: %s", schema),
		cause,
	).WithSuggestions([]string{
		"Verifique formato do arquivo",
		"Verifique campos obrigatórios",
		"Verifique tipos de dados",
	})
}
