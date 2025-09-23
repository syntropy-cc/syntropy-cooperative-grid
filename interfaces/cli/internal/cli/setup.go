package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// ManagerConfig representa a configuração do gerenciador
type ManagerConfig struct {
	Version     string                 `json:"version"`
	Created     string                 `json:"created"`
	ManagerID   string                 `json:"manager_id"`
	SystemInfo  SystemInfo             `json:"system_info"`
	Discovery   DiscoveryConfig        `json:"discovery"`
	Security    SecurityConfig         `json:"security"`
	Preferences PreferencesConfig      `json:"preferences"`
	Notifications NotificationsConfig  `json:"notifications"`
	Backup      BackupConfig           `json:"backup"`
}

type SystemInfo struct {
	Hostname     string `json:"hostname"`
	User         string `json:"user"`
	OS           string `json:"os"`
	Architecture string `json:"architecture"`
}

type DiscoveryConfig struct {
	Enabled            bool     `json:"enabled"`
	ScanNetworks       []string `json:"scan_networks"`
	DefaultSSHPort     int      `json:"default_ssh_port"`
	ConnectionTimeout  int      `json:"connection_timeout"`
	ParallelScans      int      `json:"parallel_scans"`
	CacheResults       bool     `json:"cache_results"`
	CacheTTLMinutes    int      `json:"cache_ttl_minutes"`
}

type SecurityConfig struct {
	KeyRotationDays      int  `json:"key_rotation_days"`
	RequireConfirmation  bool `json:"require_confirmation"`
	AuditLog            bool `json:"audit_log"`
	BackupKeys          bool `json:"backup_keys"`
	VerifyFingerprints  bool `json:"verify_fingerprints"`
}

type PreferencesConfig struct {
	DefaultEditor         string `json:"default_editor"`
	ShowCoordinates       bool   `json:"show_coordinates"`
	AutoUpdateMetadata    bool   `json:"auto_update_metadata"`
	ConcurrentConnections int    `json:"concurrent_connections"`
	VerboseOutput         bool   `json:"verbose_output"`
	ColorOutput           bool   `json:"color_output"`
}

type NotificationsConfig struct {
	Enabled       bool   `json:"enabled"`
	Email         string `json:"email"`
	WebhookURL    string `json:"webhook_url"`
	SlackChannel  string `json:"slack_channel"`
}

type BackupConfig struct {
	AutoBackup         bool `json:"auto_backup"`
	BackupFrequencyDays int `json:"backup_frequency_days"`
	MaxBackups         int `json:"max_backups"`
	CompressBackups    bool `json:"compress_backups"`
}

// NewSetupCommand cria o comando de setup
func NewSetupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Setup Syntropy Cooperative Grid management environment",
		Long: `Setup the complete Syntropy Cooperative Grid management infrastructure.

This command initializes the management environment by:
1. Creating directory structure
2. Installing dependencies
3. Creating manager configuration
4. Setting up helper scripts
5. Creating application templates
6. Setting up command aliases

The setup creates a complete management environment in ~/.syntropy/
that allows you to manage nodes, deploy applications, and monitor
the cooperative grid network.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSetup()
		},
	}

	return cmd
}

// runSetup executa o processo completo de setup
func runSetup() error {
	fmt.Println("🚀 Syntropy Cooperative Grid - Management Environment Setup")
	fmt.Println("Setting up complete node management infrastructure...")
	fmt.Println()

	// Verificar se já existe
	syntropyDir := getSyntropyDir()
	configDir := filepath.Join(syntropyDir, "config")
	managerConfig := filepath.Join(configDir, "manager.json")

	if _, err := os.Stat(managerConfig); err == nil {
		fmt.Println("⚠️  Syntropy management environment already exists.")
		fmt.Printf("Manager config: %s\n", managerConfig)
		fmt.Println()
		
		var response string
		fmt.Print("Reinitialize? This will preserve existing nodes but reset configuration (y/N): ")
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Setup cancelled. Use 'syntropy manager' to manage existing nodes.")
			return nil
		}
	}

	// Etapa 1: Criar estrutura de diretórios
	fmt.Println("📁 [1/6] Creating directory structure...")
	if err := createDirectoryStructure(); err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}

	// Etapa 2: Instalar dependências
	fmt.Println("📦 [2/6] Installing dependencies...")
	if err := installDependencies(); err != nil {
		return fmt.Errorf("failed to install dependencies: %w", err)
	}

	// Etapa 3: Criar configuração do gerenciador
	fmt.Println("⚙️  [3/6] Creating manager configuration...")
	if err := createManagerConfig(); err != nil {
		return fmt.Errorf("failed to create manager configuration: %w", err)
	}

	// Etapa 4: Criar scripts auxiliares
	fmt.Println("🔧 [4/6] Creating helper scripts and tools...")
	if err := createHelperScripts(); err != nil {
		return fmt.Errorf("failed to create helper scripts: %w", err)
	}

	// Etapa 5: Criar templates de aplicações
	fmt.Println("📋 [5/6] Creating application templates...")
	if err := createApplicationTemplates(); err != nil {
		return fmt.Errorf("failed to create application templates: %w", err)
	}

	// Etapa 6: Configurar aliases e PATH
	fmt.Println("🔗 [6/6] Setting up command aliases and PATH...")
	if err := setupCommandAliases(); err != nil {
		return fmt.Errorf("failed to setup command aliases: %w", err)
	}

	// Mostrar resumo final
	showSetupComplete()

	return nil
}

// createDirectoryStructure cria a estrutura de diretórios
func createDirectoryStructure() error {
	syntropyDir := getSyntropyDir()
	dirs := []string{
		filepath.Join(syntropyDir, "nodes"),
		filepath.Join(syntropyDir, "keys"),
		filepath.Join(syntropyDir, "config"),
		filepath.Join(syntropyDir, "cache"),
		filepath.Join(syntropyDir, "scripts"),
		filepath.Join(syntropyDir, "backups"),
		filepath.Join(syntropyDir, "config", "removed"),
		filepath.Join(syntropyDir, "config", "templates"),
		filepath.Join(syntropyDir, "config", "logs"),
		filepath.Join(syntropyDir, "config", "templates", "applications"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	fmt.Printf("Directory structure created:\n")
	fmt.Printf("├── %s\n", syntropyDir)
	fmt.Printf("│   ├── nodes/           # Node metadata and configurations\n")
	fmt.Printf("│   ├── keys/            # SSH keys for all managed nodes\n")
	fmt.Printf("│   ├── config/          # Manager configuration and templates\n")
	fmt.Printf("│   ├── cache/           # Temporary files and discovery cache\n")
	fmt.Printf("│   ├── scripts/         # Custom scripts and tools\n")
	fmt.Printf("│   └── backups/         # Node configuration backups\n")

	return nil
}

// installDependencies instala dependências necessárias
func installDependencies() error {
	missingTools := []string{}

	// Verificar ferramentas necessárias
	requiredTools := []string{"jq", "nmap", "python3", "ssh-keygen", "curl"}
	
	for _, tool := range requiredTools {
		if !commandExists(tool) {
			missingTools = append(missingTools, tool)
		}
	}

	if len(missingTools) > 0 {
		fmt.Printf("Installing missing dependencies: %v\n", missingTools)
		
		// Detectar gerenciador de pacotes e instalar
		if commandExists("apt-get") {
			cmd := exec.Command("sudo", "apt-get", "update")
			cmd.Run()
			
			args := append([]string{"apt-get", "install", "-y"}, missingTools...)
			cmd = exec.Command("sudo", args...)
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to install dependencies with apt-get: %w", err)
			}
		} else if commandExists("yum") {
			args := append([]string{"yum", "install", "-y"}, missingTools...)
			cmd := exec.Command("sudo", args...)
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to install dependencies with yum: %w", err)
			}
		} else if commandExists("brew") {
			cmd := exec.Command("brew", "install")
			cmd.Args = append(cmd.Args, missingTools...)
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to install dependencies with brew: %w", err)
			}
		} else {
			fmt.Printf("⚠️  Please install manually: %v\n", missingTools)
		}
	}

	return nil
}

// createManagerConfig cria a configuração do gerenciador
func createManagerConfig() error {
	syntropyDir := getSyntropyDir()
	configDir := filepath.Join(syntropyDir, "config")
	managerConfig := filepath.Join(configDir, "manager.json")

	// Gerar ID único do gerenciador
	managerID := generateManagerID()

	// Detectar redes locais para descoberta
	localNetworks := detectLocalNetworks()

	// Obter informações do sistema
	hostname, _ := os.Hostname()
	user := os.Getenv("USER")
	if user == "" {
		user = os.Getenv("USERNAME")
	}

	config := ManagerConfig{
		Version:   "1.0.0",
		Created:   time.Now().UTC().Format(time.RFC3339),
		ManagerID: managerID,
		SystemInfo: SystemInfo{
			Hostname:     hostname,
			User:         user,
			OS:           getOSInfo(),
			Architecture: getArchitecture(),
		},
		Discovery: DiscoveryConfig{
			Enabled:           true,
			ScanNetworks:      localNetworks,
			DefaultSSHPort:    22,
			ConnectionTimeout: 10,
			ParallelScans:     5,
			CacheResults:      true,
			CacheTTLMinutes:   30,
		},
		Security: SecurityConfig{
			KeyRotationDays:     90,
			RequireConfirmation: true,
			AuditLog:           true,
			BackupKeys:         true,
			VerifyFingerprints: true,
		},
		Preferences: PreferencesConfig{
			DefaultEditor:         "nano",
			ShowCoordinates:       true,
			AutoUpdateMetadata:    true,
			ConcurrentConnections: 3,
			VerboseOutput:         false,
			ColorOutput:           true,
		},
		Notifications: NotificationsConfig{
			Enabled:      false,
			Email:        "",
			WebhookURL:   "",
			SlackChannel: "",
		},
		Backup: BackupConfig{
			AutoBackup:         true,
			BackupFrequencyDays: 7,
			MaxBackups:         30,
			CompressBackups:    true,
		},
	}

	// Salvar configuração
	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(managerConfig, configData, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Printf("Manager configuration created with ID: %s\n", managerID)
	return nil
}

// createHelperScripts cria scripts auxiliares
func createHelperScripts() error {
	syntropyDir := getSyntropyDir()
	scriptsDir := filepath.Join(syntropyDir, "scripts")

	// Script de descoberta de rede
	discoverScript := `#!/bin/bash

# Quick network discovery for Syntropy nodes
echo "Discovering devices on local networks..."

# Get network interfaces and their networks
NETWORKS=$(ip route | grep -E "192\.168\.|10\.|172\." | grep "/" | awk '{print $1}' | sort -u)

for network in $NETWORKS; do
    echo "Scanning $network..."
    nmap -sn "$network" 2>/dev/null | grep "Nmap scan report" | awk '{print $5}' | while read ip; do
        if [[ $ip =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
            # Check if SSH is open
            if nc -z -w2 "$ip" 22 2>/dev/null; then
                echo "  $ip - SSH available"
            fi
        fi
    done
done`

	discoverPath := filepath.Join(scriptsDir, "discover-network.sh")
	if err := os.WriteFile(discoverPath, []byte(discoverScript), 0755); err != nil {
		return fmt.Errorf("failed to create discover script: %w", err)
	}

	// Script de backup
	backupScript := `#!/bin/bash

# Backup all node configurations
BACKUP_DIR="$HOME/.syntropy/backups/$(date +%Y%m%d_%H%M%S)"
mkdir -p "$BACKUP_DIR"

echo "Creating backup in: $BACKUP_DIR"

# Copy all node metadata
cp -r "$HOME/.syntropy/nodes" "$BACKUP_DIR/"
cp -r "$HOME/.syntropy/keys" "$BACKUP_DIR/"
cp -r "$HOME/.syntropy/config" "$BACKUP_DIR/"

# Create backup info
cat > "$BACKUP_DIR/backup_info.json" << EOF
{
  "created": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "system": "$(hostname)",
  "user": "$USER",
  "nodes_count": $(ls -1 "$HOME/.syntropy/nodes"/*.json 2>/dev/null | wc -l),
  "backup_type": "full"
}
EOF

# Compress backup
tar -czf "$BACKUP_DIR.tar.gz" -C "$(dirname "$BACKUP_DIR")" "$(basename "$BACKUP_DIR")"
rm -rf "$BACKUP_DIR"

echo "Backup created: $BACKUP_DIR.tar.gz"`

	backupPath := filepath.Join(scriptsDir, "backup-all-nodes.sh")
	if err := os.WriteFile(backupPath, []byte(backupScript), 0755); err != nil {
		return fmt.Errorf("failed to create backup script: %w", err)
	}

	// Script de health check
	healthScript := `#!/bin/bash

# Check health of all managed nodes
TOTAL=0
ONLINE=0
OFFLINE=0

echo "=== SYNTROPY NETWORK HEALTH CHECK ==="
echo ""

for node_file in "$HOME/.syntropy/nodes"/*.json; do
    if [ -f "$node_file" ]; then
        ((TOTAL++))
        
        node_name=$(basename "$node_file" .json)
        ip=$(jq -r '.network.ip_address // "unknown"' "$node_file")
        
        if [ "$ip" != "unknown" ] && [ "$ip" != "null" ]; then
            if nc -z -w3 "$ip" 22 2>/dev/null; then
                echo "✓ $node_name ($ip) - Online"
                ((ONLINE++))
            else
                echo "✗ $node_name ($ip) - Offline"
                ((OFFLINE++))
            fi
        else
            echo "? $node_name - No IP address"
            ((OFFLINE++))
        fi
    fi
done

echo ""
echo "Summary: $ONLINE/$TOTAL nodes online"`

	healthPath := filepath.Join(scriptsDir, "health-check-all.sh")
	if err := os.WriteFile(healthPath, []byte(healthScript), 0755); err != nil {
		return fmt.Errorf("failed to create health script: %w", err)
	}

	return nil
}

// createApplicationTemplates cria templates de aplicações
func createApplicationTemplates() error {
	syntropyDir := getSyntropyDir()
	templatesDir := filepath.Join(syntropyDir, "config", "templates", "applications")

	// Template Fortran
	fortranTemplate := `apiVersion: batch/v1
kind: Job
metadata:
  name: fortran-simulation
  labels:
    app: scientific-computing
    language: fortran
spec:
  template:
    spec:
      containers:
      - name: fortran-runner
        image: gcc:latest
        workingDir: /workspace
        command: ["/bin/bash"]
        args:
          - -c
          - |
            # Install Fortran compiler
            apt-get update && apt-get install -y gfortran
            
            # Compile and run Fortran program
            echo "program hello
              print *, 'Hello from Syntropy Cooperative Grid!'
              print *, 'Running Fortran simulation...'
            end program hello" > simulation.f90
            
            gfortran -o simulation simulation.f90
            ./simulation
            
            echo "Computation complete"
        resources:
          requests:
            cpu: "500m"
            memory: "512Mi"
          limits:
            cpu: "2000m"
            memory: "2Gi"
        volumeMounts:
        - name: workspace
          mountPath: /workspace
      volumes:
      - name: workspace
        emptyDir: {}
      restartPolicy: Never
  backoffLimit: 3`

	fortranPath := filepath.Join(templatesDir, "fortran-computation.yaml")
	if err := os.WriteFile(fortranPath, []byte(fortranTemplate), 0644); err != nil {
		return fmt.Errorf("failed to create fortran template: %w", err)
	}

	// Template Python Data Science
	pythonTemplate := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: jupyter-lab
  labels:
    app: data-science
    language: python
spec:
  replicas: 1
  selector:
    matchLabels:
      app: jupyter-lab
  template:
    metadata:
      labels:
        app: jupyter-lab
    spec:
      containers:
      - name: jupyter
        image: jupyter/datascience-notebook:latest
        ports:
        - containerPort: 8888
        env:
        - name: JUPYTER_ENABLE_LAB
          value: "yes"
        - name: JUPYTER_TOKEN
          value: "syntropy123"
        resources:
          requests:
            cpu: "200m"
            memory: "512Mi"
          limits:
            cpu: "2000m"
            memory: "4Gi"
        volumeMounts:
        - name: notebooks
          mountPath: /home/jovyan/work
      volumes:
      - name: notebooks
        emptyDir: {}
---
apiVersion: v1
kind: Service
metadata:
  name: jupyter-service
spec:
  selector:
    app: jupyter-lab
  ports:
  - port: 8888
    targetPort: 8888
  type: NodePort`

	pythonPath := filepath.Join(templatesDir, "python-datascience.yaml")
	if err := os.WriteFile(pythonPath, []byte(pythonTemplate), 0644); err != nil {
		return fmt.Errorf("failed to create python template: %w", err)
	}

	return nil
}

// setupCommandAliases configura aliases de comando
func setupCommandAliases() error {
	syntropyDir := getSyntropyDir()
	configDir := filepath.Join(syntropyDir, "config")
	bashrcPath := filepath.Join(configDir, "syntropy.bashrc")

	// Criar arquivo de aliases
	aliases := `# Syntropy Cooperative Grid - Management Aliases and Functions

# Core commands
alias syntropy-list='syntropy manager list'
alias syntropy-status='syntropy manager status-all'
alias syntropy-discover='syntropy manager discover'
alias syntropy-health='~/.syntropy/scripts/health-check-all.sh'
alias syntropy-backup='~/.syntropy/scripts/backup-all-nodes.sh'

# Quick node access
syntropy-connect() {
    if [ -z "$1" ]; then
        echo "Usage: syntropy-connect <node-name>"
        syntropy manager list
        return 1
    fi
    syntropy manager connect "$1"
}

syntropy-ssh() {
    if [ -z "$1" ]; then
        echo "Usage: syntropy-ssh <node-name>"
        return 1
    fi
    
    local node_file="$HOME/.syntropy/nodes/${1}.json"
    local key_file="$HOME/.syntropy/keys/${1}_owner.key"
    
    if [ ! -f "$node_file" ]; then
        echo "Node $1 not found"
        return 1
    fi
    
    local ip=$(jq -r '.network.ip_address' "$node_file" 2>/dev/null)
    if [ "$ip" = "null" ] || [ -z "$ip" ]; then
        echo "No IP address for node $1. Try: syntropy-discover"
        return 1
    fi
    
    ssh -i "$key_file" admin@"$ip"
}

# Node creation shortcut
syntropy-create() {
    local usb_device="$1"
    shift
    syntropy usb create "$usb_device" "$@"
}

# Show syntropy status
syntropy-info() {
    echo "=== SYNTROPY MANAGEMENT STATUS ==="
    echo "Managed nodes: $(ls -1 "$HOME/.syntropy/nodes"/*.json 2>/dev/null | wc -l)"
    echo "SSH keys: $(ls -1 "$HOME/.syntropy/keys"/*_owner.key 2>/dev/null | wc -l)"
    echo "Config dir: $HOME/.syntropy"
    echo ""
    echo "Recent nodes:"
    ls -t "$HOME/.syntropy/nodes"/*.json 2>/dev/null | head -3 | while read f; do
        local name=$(basename "$f" .json)
        local ip=$(jq -r '.network.ip_address // "unknown"' "$f")
        echo "  $name ($ip)"
    done
}

export PATH="$HOME/.syntropy/scripts:$PATH"`

	if err := os.WriteFile(bashrcPath, []byte(aliases), 0644); err != nil {
		return fmt.Errorf("failed to create aliases file: %w", err)
	}

	// Adicionar ao bashrc do usuário se não existir
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	userBashrc := filepath.Join(homeDir, ".bashrc")
	bashrcContent, err := os.ReadFile(userBashrc)
	if err != nil {
		// Se não conseguir ler, criar novo
		bashrcContent = []byte{}
	}

	if !strings.Contains(string(bashrcContent), "syntropy.bashrc") {
		addition := fmt.Sprintf("\n# Syntropy Cooperative Grid Management\nsource %s\n", bashrcPath)
		if err := os.WriteFile(userBashrc, append(bashrcContent, []byte(addition)...), 0644); err != nil {
			return fmt.Errorf("failed to update user bashrc: %w", err)
		}
		fmt.Println("Syntropy aliases added to ~/.bashrc")
	} else {
		fmt.Println("Syntropy aliases already configured in ~/.bashrc")
	}

	return nil
}

// showSetupComplete mostra o resumo final do setup
func showSetupComplete() {
	syntropyDir := getSyntropyDir()
	configDir := filepath.Join(syntropyDir, "config")
	managerConfig := filepath.Join(configDir, "manager.json")

	fmt.Println()
	fmt.Println("🎉 SETUP COMPLETE - READY TO USE!")
	fmt.Println()
	fmt.Println("═══ SYNTROPY MANAGEMENT ENVIRONMENT ═══")
	fmt.Printf("🏠 Base directory: %s\n", syntropyDir)
	fmt.Printf("📋 Configuration: %s\n", managerConfig)
	fmt.Println()
	fmt.Println("═══ AVAILABLE COMMANDS ═══")
	fmt.Println("Core Management:")
	fmt.Println("  syntropy manager list              # List all nodes")
	fmt.Println("  syntropy manager connect <node>    # Connect to specific node")
	fmt.Println("  syntropy manager status <node>     # Check node status")
	fmt.Println("  syntropy manager discover          # Discover nodes on network")
	fmt.Println()
	fmt.Println("Quick Aliases (after source ~/.bashrc):")
	fmt.Println("  syntropy-list                         # List all nodes")
	fmt.Println("  syntropy-connect <node>               # Connect to node")
	fmt.Println("  syntropy-ssh <node>                   # Direct SSH to node")
	fmt.Println("  syntropy-status                       # Status of all nodes")
	fmt.Println("  syntropy-discover                     # Discover network nodes")
	fmt.Println("  syntropy-health                       # Health check all nodes")
	fmt.Println("  syntropy-backup                       # Backup all configurations")
	fmt.Println("  syntropy-info                         # Show management summary")
	fmt.Println()
	fmt.Println("Node Creation:")
	fmt.Println("  syntropy usb create /dev/sdb --node-name my-server       # Create node")
	fmt.Println()
	fmt.Println("═══ NEXT STEPS ═══")
	fmt.Println("1. Reload shell: source ~/.bashrc")
	fmt.Println("2. Create your first node:")
	fmt.Println("   syntropy usb create /dev/sdb --node-name main-server")
	fmt.Println("3. Boot USB on target hardware")
	fmt.Println("4. Wait for installation (~30 minutes)")
	fmt.Println("5. Connect to node: syntropy manager connect main-server")
	fmt.Println()
	fmt.Println("Management environment ready! Start creating nodes with 'syntropy usb create'")
}

// Funções auxiliares

func getSyntropyDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".syntropy")
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func generateManagerID() string {
	cmd := exec.Command("openssl", "rand", "-hex", "8")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Sprintf("mgr-%d", time.Now().Unix())
	}
	return fmt.Sprintf("mgr-%s", strings.TrimSpace(string(output)))
}

func detectLocalNetworks() []string {
	cmd := exec.Command("ip", "route")
	output, err := cmd.Output()
	if err != nil {
		return []string{"192.168.1.0/24", "10.0.0.0/24", "172.16.0.0/16"}
	}

	lines := strings.Split(string(output), "\n")
	networks := []string{}
	
	for _, line := range lines {
		if strings.Contains(line, "192.168.") || strings.Contains(line, "10.") || strings.Contains(line, "172.") {
			parts := strings.Fields(line)
			if len(parts) > 0 && strings.Contains(parts[0], "/") {
				networks = append(networks, parts[0])
			}
		}
	}

	// Remover duplicatas e limitar a 5
	unique := make(map[string]bool)
	result := []string{}
	for _, net := range networks {
		if !unique[net] && len(result) < 5 {
			unique[net] = true
			result = append(result, net)
		}
	}

	// Adicionar redes padrão se não encontrou nenhuma
	if len(result) == 0 {
		result = []string{"192.168.1.0/24", "10.0.0.0/24", "172.16.0.0/16"}
	}

	return result
}

func getOSInfo() string {
	cmd := exec.Command("uname", "-s")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(output))
}

func getArchitecture() string {
	cmd := exec.Command("uname", "-m")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(output))
}

