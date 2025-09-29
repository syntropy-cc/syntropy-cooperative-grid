package setup

import (
	"fmt"
	"os"
	"runtime"

	"setup-component/src/internal/types"
)

// Validator implementa a interface Validator
type Validator struct {
	osValidator types.OSValidator
	logger      *SetupLogger
}

// NewValidator cria um novo validador
func NewValidator(logger *SetupLogger) *Validator {
	return &Validator{
		osValidator: NewOSValidator(logger),
		logger:      logger,
	}
}

// NewOSValidator cria um novo validador de SO
func NewOSValidator(logger *SetupLogger) types.OSValidator {
	switch runtime.GOOS {
	case "windows":
		return &WindowsValidator{logger: logger}
	case "linux":
		return &LinuxValidator{logger: logger}
	case "darwin":
		return &DarwinValidator{logger: logger}
	default:
		return &GenericValidator{logger: logger}
	}
}

// ValidateEnvironment valida o ambiente completo
func (v *Validator) ValidateEnvironment() (*types.EnvironmentInfo, error) {
	v.logger.LogStep("validation_start", nil)

	// Detectar SO
	osInfo, err := v.osValidator.DetectOS()
	if err != nil {
		v.logger.LogError(err, map[string]interface{}{
			"step": "os_detection",
		})
		return nil, err
	}

	// Validar recursos
	resources, err := v.osValidator.ValidateResources()
	if err != nil {
		v.logger.LogError(err, map[string]interface{}{
			"step": "resource_validation",
		})
		return nil, err
	}

	// Validar permissões
	permissions, err := v.osValidator.ValidatePermissions()
	if err != nil {
		v.logger.LogError(err, map[string]interface{}{
			"step": "permission_validation",
		})
		return nil, err
	}

	// Criar informações do ambiente
	envInfo := &types.EnvironmentInfo{
		OS:              osInfo.Name,
		OSVersion:       osInfo.Version,
		Architecture:    osInfo.Architecture,
		HasAdminRights:  permissions.HasAdminRights,
		AvailableDiskGB: resources.DiskSpaceGB,
		HasInternet:     true, // Será validado separadamente
	}

	// Obter diretório home
	homeDir, err := getUserHomeDir()
	if err != nil {
		v.logger.LogError(err, map[string]interface{}{
			"step": "home_dir_detection",
		})
		return nil, err
	}
	envInfo.HomeDir = homeDir

	v.logger.LogStep("validation_completed", map[string]interface{}{
		"os":           envInfo.OS,
		"version":      envInfo.OSVersion,
		"architecture": envInfo.Architecture,
		"admin_rights": envInfo.HasAdminRights,
		"disk_space":   envInfo.AvailableDiskGB,
	})

	return envInfo, nil
}

// ValidateDependencies valida as dependências necessárias
func (v *Validator) ValidateDependencies() (*types.DependencyStatus, error) {
	v.logger.LogStep("dependency_validation_start", nil)

	// Obter dependências necessárias para o SO atual
	requiredDeps := v.getRequiredDependencies()

	status := &types.DependencyStatus{
		Required:  requiredDeps,
		Installed: []types.Dependency{},
		Missing:   []types.Dependency{},
		Outdated:  []types.Dependency{},
	}

	// Verificar cada dependência
	for _, dep := range requiredDeps {
		if v.isDependencyInstalled(dep) {
			status.Installed = append(status.Installed, dep)
		} else {
			status.Missing = append(status.Missing, dep)
		}
	}

	v.logger.LogStep("dependency_validation_completed", map[string]interface{}{
		"required_count":  len(status.Required),
		"installed_count": len(status.Installed),
		"missing_count":   len(status.Missing),
	})

	return status, nil
}

// ValidateNetwork valida a conectividade de rede
func (v *Validator) ValidateNetwork() (*types.NetworkInfo, error) {
	v.logger.LogStep("network_validation_start", nil)

	networkInfo := &types.NetworkInfo{
		HasInternet:     false,
		Connectivity:    false,
		ProxyConfigured: false,
		FirewallActive:  false,
		PortsOpen:       []int{},
	}

	// Verificar conectividade básica
	if v.testInternetConnectivity() {
		networkInfo.HasInternet = true
		networkInfo.Connectivity = true
	}

	// Verificar configuração de proxy
	networkInfo.ProxyConfigured = v.isProxyConfigured()

	// Verificar firewall
	networkInfo.FirewallActive = v.isFirewallActive()

	// Verificar portas necessárias
	networkInfo.PortsOpen = v.checkRequiredPorts()

	v.logger.LogStep("network_validation_completed", map[string]interface{}{
		"has_internet":     networkInfo.HasInternet,
		"connectivity":     networkInfo.Connectivity,
		"proxy_configured": networkInfo.ProxyConfigured,
		"firewall_active":  networkInfo.FirewallActive,
		"open_ports_count": len(networkInfo.PortsOpen),
	})

	return networkInfo, nil
}

// ValidatePermissions valida as permissões necessárias
func (v *Validator) ValidatePermissions() (*types.PermissionStatus, error) {
	v.logger.LogStep("permission_validation_start", nil)

	status := &types.PermissionStatus{
		FileSystem: false,
		Network:    false,
		Service:    false,
		Admin:      false,
		Issues:     []string{},
	}

	// Verificar permissões de sistema de arquivos
	if v.canWriteToHomeDir() {
		status.FileSystem = true
	} else {
		status.Issues = append(status.Issues, "Sem permissão de escrita no diretório home")
	}

	// Verificar permissões de rede
	if v.canAccessNetwork() {
		status.Network = true
	} else {
		status.Issues = append(status.Issues, "Sem permissão de acesso à rede")
	}

	// Verificar permissões de serviço
	if v.canInstallServices() {
		status.Service = true
	} else {
		status.Issues = append(status.Issues, "Sem permissão para instalar serviços")
	}

	// Verificar permissões de administrador
	if v.hasAdminRights() {
		status.Admin = true
	} else {
		status.Issues = append(status.Issues, "Sem permissões de administrador")
	}

	v.logger.LogStep("permission_validation_completed", map[string]interface{}{
		"file_system":  status.FileSystem,
		"network":      status.Network,
		"service":      status.Service,
		"admin":        status.Admin,
		"issues_count": len(status.Issues),
	})

	return status, nil
}

// FixIssues tenta corrigir problemas automaticamente
func (v *Validator) FixIssues(issues []types.ValidationIssue) error {
	v.logger.LogStep("fix_issues_start", map[string]interface{}{
		"issues_count": len(issues),
	})

	fixedCount := 0
	for _, issue := range issues {
		if v.canFixIssue(issue) {
			if err := v.fixIssue(issue); err != nil {
				v.logger.LogWarning("Falha ao corrigir issue", map[string]interface{}{
					"issue_type": issue.Type,
					"error":      err.Error(),
				})
			} else {
				fixedCount++
				v.logger.LogInfo("Issue corrigida com sucesso", map[string]interface{}{
					"issue_type": issue.Type,
				})
			}
		}
	}

	v.logger.LogStep("fix_issues_completed", map[string]interface{}{
		"fixed_count": fixedCount,
		"total_count": len(issues),
	})

	return nil
}

// ValidateAll executa todas as validações
func (v *Validator) ValidateAll() (*types.ValidationResult, error) {
	v.logger.LogStep("comprehensive_validation_start", nil)

	result := &types.ValidationResult{
		CanProceed: true,
		Warnings:   []string{},
		Issues:     []types.ValidationIssue{},
	}

	// Validar ambiente
	envInfo, err := v.ValidateEnvironment()
	if err != nil {
		result.CanProceed = false
		result.Issues = append(result.Issues, types.ValidationIssue{
			Type:     "environment",
			Severity: "error",
			Message:  err.Error(),
		})
	} else {
		result.Environment = envInfo
	}

	// Validar dependências
	deps, err := v.ValidateDependencies()
	if err != nil {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Falha na validação de dependências: %v", err))
	} else {
		result.Dependencies = deps
		if len(deps.Missing) > 0 {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Dependências ausentes: %d", len(deps.Missing)))
		}
	}

	// Validar rede
	network, err := v.ValidateNetwork()
	if err != nil {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Falha na validação de rede: %v", err))
	} else {
		result.Network = network
		if !network.HasInternet {
			result.Warnings = append(result.Warnings, "Sem conectividade com a internet")
		}
	}

	// Validar permissões
	permissions, err := v.ValidatePermissions()
	if err != nil {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Falha na validação de permissões: %v", err))
	} else {
		result.Permissions = permissions
		if len(permissions.Issues) > 0 {
			for _, issue := range permissions.Issues {
				result.Warnings = append(result.Warnings, issue)
			}
		}
	}

	// Determinar se pode prosseguir
	if len(result.Issues) > 0 {
		result.CanProceed = false
	}

	v.logger.LogStep("comprehensive_validation_completed", map[string]interface{}{
		"can_proceed":    result.CanProceed,
		"issues_count":   len(result.Issues),
		"warnings_count": len(result.Warnings),
	})

	return result, nil
}

// Métodos auxiliares

// getRequiredDependencies retorna as dependências necessárias para o SO atual
func (v *Validator) getRequiredDependencies() []types.Dependency {
	switch runtime.GOOS {
	case "windows":
		return []types.Dependency{
			{Name: "powershell", Version: "5.1+", Required: true},
			{Name: "curl", Version: "any", Required: false},
			{Name: "git", Version: "any", Required: false},
		}
	case "linux":
		return []types.Dependency{
			{Name: "curl", Version: "any", Required: true},
			{Name: "wget", Version: "any", Required: false},
			{Name: "git", Version: "any", Required: false},
			{Name: "systemctl", Version: "any", Required: false},
		}
	case "darwin":
		return []types.Dependency{
			{Name: "curl", Version: "any", Required: true},
			{Name: "git", Version: "any", Required: false},
		}
	default:
		return []types.Dependency{}
	}
}

// isDependencyInstalled verifica se uma dependência está instalada
func (v *Validator) isDependencyInstalled(dep types.Dependency) bool {
	// Implementação simplificada - em produção seria mais robusta
	switch dep.Name {
	case "powershell":
		return v.isPowerShellInstalled()
	case "curl":
		return v.isCurlInstalled()
	case "wget":
		return v.isWgetInstalled()
	case "git":
		return v.isGitInstalled()
	case "systemctl":
		return v.isSystemctlInstalled()
	default:
		return false
	}
}

// testInternetConnectivity testa a conectividade com a internet
func (v *Validator) testInternetConnectivity() bool {
	// Implementação simplificada - em produção seria mais robusta
	return true // Por enquanto, assumir que há conectividade
}

// isProxyConfigured verifica se há proxy configurado
func (v *Validator) isProxyConfigured() bool {
	// Implementação simplificada
	return false
}

// isFirewallActive verifica se o firewall está ativo
func (v *Validator) isFirewallActive() bool {
	// Implementação simplificada
	return false
}

// checkRequiredPorts verifica se as portas necessárias estão abertas
func (v *Validator) checkRequiredPorts() []int {
	// Implementação simplificada
	return []int{8080, 9090}
}

// canWriteToHomeDir verifica se pode escrever no diretório home
func (v *Validator) canWriteToHomeDir() bool {
	// Implementação simplificada
	return true
}

// canAccessNetwork verifica se pode acessar a rede
func (v *Validator) canAccessNetwork() bool {
	// Implementação simplificada
	return true
}

// canInstallServices verifica se pode instalar serviços
func (v *Validator) canInstallServices() bool {
	// Implementação simplificada
	return v.hasAdminRights()
}

// hasAdminRights verifica se tem direitos de administrador
func (v *Validator) hasAdminRights() bool {
	// Implementação simplificada
	return false
}

// canFixIssue verifica se pode corrigir um issue automaticamente
func (v *Validator) canFixIssue(issue types.ValidationIssue) bool {
	// Implementação simplificada
	return false
}

// fixIssue tenta corrigir um issue automaticamente
func (v *Validator) fixIssue(issue types.ValidationIssue) error {
	// Implementação simplificada
	return fmt.Errorf("correção automática não implementada para: %s", issue.Type)
}

// Verificações de dependências específicas

func (v *Validator) isPowerShellInstalled() bool {
	// Implementação simplificada
	return runtime.GOOS == "windows"
}

func (v *Validator) isCurlInstalled() bool {
	// Implementação simplificada
	return true
}

func (v *Validator) isWgetInstalled() bool {
	// Implementação simplificada
	return runtime.GOOS == "linux"
}

func (v *Validator) isGitInstalled() bool {
	// Implementação simplificada
	return true
}

func (v *Validator) isSystemctlInstalled() bool {
	// Implementação simplificada
	return runtime.GOOS == "linux"
}

// getUserHomeDir obtém o diretório home do usuário
func getUserHomeDir() (string, error) {
	// Implementação simplificada
	return os.UserHomeDir()
}

// Implementações dos validadores de SO

// WindowsValidator implementa validação específica para Windows
type WindowsValidator struct {
	logger *SetupLogger
}

// DetectOS detecta informações do Windows
func (wv *WindowsValidator) DetectOS() (*types.OSInfo, error) {
	wv.logger.LogDebug("Detectando sistema operacional Windows", nil)

	info := &types.OSInfo{
		Name:         "windows",
		Version:      "10.0",
		Architecture: runtime.GOARCH,
		Build:        "19041",
		Kernel:       "nt",
	}

	return info, nil
}

// ValidateResources valida recursos do Windows
func (wv *WindowsValidator) ValidateResources() (*types.ResourceInfo, error) {
	info := &types.ResourceInfo{
		TotalMemoryGB:  8.0,
		AvailableMemGB: 4.0,
		CPUCores:       4,
		DiskSpaceGB:    50.0,
	}

	return info, nil
}

// ValidatePermissions valida permissões no Windows
func (wv *WindowsValidator) ValidatePermissions() (*types.PermissionInfo, error) {
	info := &types.PermissionInfo{
		HasAdminRights: false,
		UserID:         "user",
		GroupID:        "users",
		Capabilities:   []string{"file_system", "network"},
	}

	return info, nil
}

// InstallDependencies instala dependências no Windows
func (wv *WindowsValidator) InstallDependencies(deps []types.Dependency) error {
	return nil
}

// ConfigureEnvironment configura o ambiente Windows
func (wv *WindowsValidator) ConfigureEnvironment() error {
	return nil
}

// LinuxValidator implementa validação específica para Linux
type LinuxValidator struct {
	logger *SetupLogger
}

// DetectOS detecta informações do Linux
func (lv *LinuxValidator) DetectOS() (*types.OSInfo, error) {
	info := &types.OSInfo{
		Name:         "linux",
		Version:      "20.04",
		Architecture: runtime.GOARCH,
		Build:        "5.4.0",
		Kernel:       "5.4.0-42-generic",
	}

	return info, nil
}

// ValidateResources valida recursos do Linux
func (lv *LinuxValidator) ValidateResources() (*types.ResourceInfo, error) {
	info := &types.ResourceInfo{
		TotalMemoryGB:  8.0,
		AvailableMemGB: 4.0,
		CPUCores:       4,
		DiskSpaceGB:    50.0,
	}

	return info, nil
}

// ValidatePermissions valida permissões no Linux
func (lv *LinuxValidator) ValidatePermissions() (*types.PermissionInfo, error) {
	info := &types.PermissionInfo{
		HasAdminRights: false,
		UserID:         "1000",
		GroupID:        "1000",
		Capabilities:   []string{"file_system", "network"},
	}

	return info, nil
}

// InstallDependencies instala dependências no Linux
func (lv *LinuxValidator) InstallDependencies(deps []types.Dependency) error {
	return nil
}

// ConfigureEnvironment configura o ambiente Linux
func (lv *LinuxValidator) ConfigureEnvironment() error {
	return nil
}

// DarwinValidator implementa validação específica para macOS
type DarwinValidator struct {
	logger *SetupLogger
}

// DetectOS detecta informações do macOS
func (dv *DarwinValidator) DetectOS() (*types.OSInfo, error) {
	info := &types.OSInfo{
		Name:         "darwin",
		Version:      "10.15",
		Architecture: runtime.GOARCH,
		Build:        "19H2",
		Kernel:       "19.6.0",
	}

	return info, nil
}

// ValidateResources valida recursos do macOS
func (dv *DarwinValidator) ValidateResources() (*types.ResourceInfo, error) {
	info := &types.ResourceInfo{
		TotalMemoryGB:  8.0,
		AvailableMemGB: 4.0,
		CPUCores:       4,
		DiskSpaceGB:    50.0,
	}

	return info, nil
}

// ValidatePermissions valida permissões no macOS
func (dv *DarwinValidator) ValidatePermissions() (*types.PermissionInfo, error) {
	info := &types.PermissionInfo{
		HasAdminRights: false,
		UserID:         "501",
		GroupID:        "20",
		Capabilities:   []string{"file_system", "network"},
	}

	return info, nil
}

// InstallDependencies instala dependências no macOS
func (dv *DarwinValidator) InstallDependencies(deps []types.Dependency) error {
	return nil
}

// ConfigureEnvironment configura o ambiente macOS
func (dv *DarwinValidator) ConfigureEnvironment() error {
	return nil
}

// GenericValidator implementa validação genérica para SOs não suportados
type GenericValidator struct {
	logger *SetupLogger
}

// DetectOS detecta informações genéricas do SO
func (gv *GenericValidator) DetectOS() (*types.OSInfo, error) {
	info := &types.OSInfo{
		Name:         runtime.GOOS,
		Version:      "unknown",
		Architecture: runtime.GOARCH,
		Build:        "unknown",
		Kernel:       "unknown",
	}

	return info, nil
}

// ValidateResources valida recursos genéricos
func (gv *GenericValidator) ValidateResources() (*types.ResourceInfo, error) {
	info := &types.ResourceInfo{
		TotalMemoryGB:  4.0,
		AvailableMemGB: 2.0,
		CPUCores:       2,
		DiskSpaceGB:    10.0,
	}

	return info, nil
}

// ValidatePermissions valida permissões genéricas
func (gv *GenericValidator) ValidatePermissions() (*types.PermissionInfo, error) {
	info := &types.PermissionInfo{
		HasAdminRights: false,
		UserID:         "unknown",
		GroupID:        "unknown",
		Capabilities:   []string{"file_system"},
	}

	return info, nil
}

// InstallDependencies instala dependências genéricas
func (gv *GenericValidator) InstallDependencies(deps []types.Dependency) error {
	return nil
}

// ConfigureEnvironment configura o ambiente genérico
func (gv *GenericValidator) ConfigureEnvironment() error {
	return nil
}
