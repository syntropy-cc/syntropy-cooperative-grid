package cli

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// NodeInfo representa informações de um nó
type NodeInfo struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Created     string            `json:"created"`
	LastSeen    string            `json:"last_seen"`
	Status      string            `json:"status"`
	Network     NetworkInfo       `json:"network"`
	Hardware    HardwareInfo      `json:"hardware"`
	Software    SoftwareInfo      `json:"software"`
	Metadata    map[string]string `json:"metadata"`
}

type NetworkInfo struct {
	IPAddress    string `json:"ip_address"`
	SSHPort      int    `json:"ssh_port"`
	MACAddress   string `json:"mac_address"`
	Hostname     string `json:"hostname"`
	LastPing     string `json:"last_ping"`
	Latency      int    `json:"latency_ms"`
}

type HardwareInfo struct {
	CPU        string `json:"cpu"`
	Memory     string `json:"memory"`
	Storage    string `json:"storage"`
	Architecture string `json:"architecture"`
	Model      string `json:"model"`
}

type SoftwareInfo struct {
	OS         string `json:"os"`
	Kernel     string `json:"kernel"`
	Docker     string `json:"docker"`
	Python     string `json:"python"`
	LastUpdate string `json:"last_update"`
}

// NewManagerCommand cria o comando de gerenciamento
func NewManagerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "manager",
		Short: "Manage Syntropy nodes",
		Long: `Manage Syntropy Cooperative Grid nodes including discovery, connection,
status monitoring, and configuration management.

This command provides comprehensive node management capabilities:
- List and discover nodes on the network
- Connect to nodes via SSH
- Monitor node status and health
- Manage node configurations
- Backup and restore node data`,
	}

	// Adicionar subcomandos
	cmd.AddCommand(newManagerListCommand())
	cmd.AddCommand(newManagerConnectCommand())
	cmd.AddCommand(newManagerStatusCommand())
	cmd.AddCommand(newManagerDiscoverCommand())
	cmd.AddCommand(newManagerBackupCommand())
	cmd.AddCommand(newManagerRestoreCommand())
	cmd.AddCommand(newManagerHealthCommand())

	return cmd
}

// newManagerListCommand cria o comando de listagem
func newManagerListCommand() *cobra.Command {
	var (
		format string
		filter string
		sortBy string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all managed nodes",
		Long: `List all nodes managed by the Syntropy Cooperative Grid.

You can filter nodes by status, sort by different criteria, and output
in various formats (table, json, yaml).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listNodes(format, filter, sortBy)
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "table", "Output format (table, json, yaml)")
	cmd.Flags().StringVar(&filter, "filter", "", "Filter nodes by status (online, offline, unknown)")
	cmd.Flags().StringVar(&sortBy, "sort", "name", "Sort by field (name, created, last_seen, status)")

	return cmd
}

// newManagerConnectCommand cria o comando de conexão
func newManagerConnectCommand() *cobra.Command {
	var (
		interactive bool
		command     string
	)

	cmd := &cobra.Command{
		Use:   "connect <node-name>",
		Short: "Connect to a specific node",
		Long: `Connect to a specific node via SSH.

This command will:
1. Look up the node's IP address and SSH key
2. Establish SSH connection
3. Provide interactive shell access

You can also execute a single command remotely.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			nodeName := args[0]
			return connectToNode(nodeName, interactive, command)
		},
	}

	cmd.Flags().BoolVarP(&interactive, "interactive", "i", true, "Start interactive SSH session")
	cmd.Flags().StringVarP(&command, "command", "c", "", "Execute single command remotely")

	return cmd
}

// newManagerStatusCommand cria o comando de status
func newManagerStatusCommand() *cobra.Command {
	var (
		format string
		watch  bool
	)

	cmd := &cobra.Command{
		Use:   "status [node-name]",
		Short: "Show node status",
		Long: `Show detailed status information for nodes.

If no node name is provided, shows status for all nodes.
Use --watch to continuously monitor status changes.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			nodeName := ""
			if len(args) > 0 {
				nodeName = args[0]
			}
			return showNodeStatus(nodeName, format, watch)
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "table", "Output format (table, json, yaml)")
	cmd.Flags().BoolVarP(&watch, "watch", "w", false, "Watch for changes")

	return cmd
}

// newManagerDiscoverCommand cria o comando de descoberta
func newManagerDiscoverCommand() *cobra.Command {
	var (
		networks    []string
		port        int
		timeout     int
		parallel    int
		updateCache bool
	)

	cmd := &cobra.Command{
		Use:   "discover",
		Short: "Discover nodes on network",
		Long: `Discover Syntropy nodes on the local network.

This command will:
1. Scan specified networks for SSH-enabled devices
2. Attempt to identify Syntropy nodes
3. Update node metadata with discovered information
4. Cache results for faster future lookups`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return discoverNodes(networks, port, timeout, parallel, updateCache)
		},
	}

	cmd.Flags().StringSliceVarP(&networks, "networks", "n", []string{}, "Networks to scan (e.g., 192.168.1.0/24)")
	cmd.Flags().IntVarP(&port, "port", "p", 22, "SSH port to check")
	cmd.Flags().IntVarP(&timeout, "timeout", "t", 10, "Connection timeout in seconds")
	cmd.Flags().IntVar(&parallel, "parallel", 5, "Number of parallel scans")
	cmd.Flags().BoolVar(&updateCache, "update-cache", true, "Update discovery cache")

	return cmd
}

// newManagerBackupCommand cria o comando de backup
func newManagerBackupCommand() *cobra.Command {
	var (
		output   string
		compress bool
		include  []string
	)

	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Backup node configurations",
		Long: `Backup all node configurations and metadata.

This creates a complete backup of:
- Node configurations and metadata
- SSH keys and certificates
- Manager configuration
- Discovery cache`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return backupNodes(output, compress, include)
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file path (default: auto-generated)")
	cmd.Flags().BoolVarP(&compress, "compress", "c", true, "Compress backup")
	cmd.Flags().StringSliceVar(&include, "include", []string{"nodes", "keys", "config", "cache"}, "Components to include")

	return cmd
}

// newManagerRestoreCommand cria o comando de restore
func newManagerRestoreCommand() *cobra.Command {
	var (
		force bool
	)

	cmd := &cobra.Command{
		Use:   "restore <backup-file>",
		Short: "Restore node configurations",
		Long: `Restore node configurations from a backup file.

This will restore:
- Node configurations and metadata
- SSH keys and certificates
- Manager configuration
- Discovery cache`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			backupFile := args[0]
			return restoreNodes(backupFile, force)
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Force restore without confirmation")

	return cmd
}

// newManagerHealthCommand cria o comando de health check
func newManagerHealthCommand() *cobra.Command {
	var (
		format string
		watch  bool
	)

	cmd := &cobra.Command{
		Use:   "health",
		Short: "Check health of all nodes",
		Long: `Perform comprehensive health check on all managed nodes.

This will check:
- Network connectivity
- SSH accessibility
- Service status
- Resource usage
- System health`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return checkNodeHealth(format, watch)
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "table", "Output format (table, json, yaml)")
	cmd.Flags().BoolVarP(&watch, "watch", "w", false, "Watch for changes")

	return cmd
}

// Implementações dos comandos

func listNodes(format, filter, sortBy string) error {
	nodes, err := loadAllNodes()
	if err != nil {
		return fmt.Errorf("failed to load nodes: %w", err)
	}

	// Aplicar filtro
	if filter != "" {
		nodes = filterNodes(nodes, filter)
	}

	// Aplicar ordenação
	sortNodes(nodes, sortBy)

	// Output
	switch format {
	case "json":
		return outputNodesJSON(nodes)
	case "yaml":
		return outputNodesYAML(nodes)
	default:
		return outputNodesTable(nodes)
	}
}

func connectToNode(nodeName string, interactive bool, command string) error {
	node, err := loadNode(nodeName)
	if err != nil {
		return fmt.Errorf("node not found: %w", err)
	}

	if node.Network.IPAddress == "" {
		return fmt.Errorf("no IP address for node %s. Try: syntropy manager discover", nodeName)
	}

	keyFile := getNodeKeyFile(nodeName)
	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		return fmt.Errorf("SSH key not found for node %s: %s", nodeName, keyFile)
	}

	sshArgs := []string{
		"-i", keyFile,
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "LogLevel=ERROR",
	}

	if !interactive {
		sshArgs = append(sshArgs, "-o", "BatchMode=yes")
	}

	if command != "" {
		sshArgs = append(sshArgs, "admin@"+node.Network.IPAddress, command)
	} else {
		sshArgs = append(sshArgs, "admin@"+node.Network.IPAddress)
	}

	cmd := exec.Command("ssh", sshArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func showNodeStatus(nodeName, format string, watch bool) error {
	if nodeName != "" {
		// Status de um nó específico
		node, err := loadNode(nodeName)
		if err != nil {
			return fmt.Errorf("node not found: %w", err)
		}

		// Atualizar status em tempo real
		updateNodeStatus(&node)

		switch format {
		case "json":
			return outputNodeJSON(node)
		case "yaml":
			return outputNodeYAML(node)
		default:
			return outputNodeTable(node)
		}
	} else {
		// Status de todos os nós
		nodes, err := loadAllNodes()
		if err != nil {
			return fmt.Errorf("failed to load nodes: %w", err)
		}

		// Atualizar status de todos os nós
		for i := range nodes {
			updateNodeStatus(&nodes[i])
		}

		if watch {
			return watchNodeStatus(nodes, format)
		}

		switch format {
		case "json":
			return outputNodesJSON(nodes)
		case "yaml":
			return outputNodesYAML(nodes)
		default:
			return outputNodesTable(nodes)
		}
	}
}

func discoverNodes(networks []string, port, timeout, parallel int, updateCache bool) error {
	fmt.Println("🔍 Discovering Syntropy nodes on network...")

	// Usar redes padrão se não especificadas
	if len(networks) == 0 {
		networks = getDefaultNetworks()
	}

	fmt.Printf("Scanning networks: %v\n", networks)
	fmt.Printf("SSH port: %d, Timeout: %ds, Parallel: %d\n", port, timeout, parallel)

	discoveredNodes := []DiscoveredNode{}

	// Executar descoberta em paralelo
	for _, network := range networks {
		fmt.Printf("Scanning %s...\n", network)
		nodes := scanNetwork(network, port, timeout)
		discoveredNodes = append(discoveredNodes, nodes...)
	}

	fmt.Printf("Found %d potential nodes\n", len(discoveredNodes))

	// Identificar nós Syntropy
	syntropyNodes := identifySyntropyNodes(discoveredNodes)

	fmt.Printf("Identified %d Syntropy nodes\n", len(syntropyNodes))

	// Atualizar ou criar nós
	for _, node := range syntropyNodes {
		if err := updateOrCreateNode(node); err != nil {
			fmt.Printf("⚠️  Failed to update node %s: %v\n", node.IP, err)
		} else {
			fmt.Printf("✅ Updated node: %s (%s)\n", node.Hostname, node.IP)
		}
	}

	// Atualizar cache se solicitado
	if updateCache {
		updateDiscoveryCache(discoveredNodes)
	}

	return nil
}

func backupNodes(output string, compress bool, include []string) error {
	fmt.Println("💾 Creating backup of node configurations...")

	syntropyDir := getSyntropyDir()
	timestamp := time.Now().Format("20060102_150405")

	if output == "" {
		output = filepath.Join(syntropyDir, "backups", fmt.Sprintf("backup_%s", timestamp))
		if compress {
			output += ".tar.gz"
		}
	}

	// Criar diretório de backup temporário
	tempDir := filepath.Join(syntropyDir, "backups", fmt.Sprintf("temp_%s", timestamp))
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Copiar componentes incluídos
	for _, component := range include {
		src := filepath.Join(syntropyDir, component)
		dst := filepath.Join(tempDir, component)
		
		if _, err := os.Stat(src); err == nil {
			cmd := exec.Command("cp", "-r", src, dst)
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to copy %s: %w", component, err)
			}
		}
	}

	// Criar arquivo de informações do backup
	backupInfo := map[string]interface{}{
		"created":      time.Now().UTC().Format(time.RFC3339),
		"system":       getHostname(),
		"user":         os.Getenv("USER"),
		"components":   include,
		"backup_type":  "full",
	}

	infoData, _ := json.MarshalIndent(backupInfo, "", "  ")
	infoFile := filepath.Join(tempDir, "backup_info.json")
	os.WriteFile(infoFile, infoData, 0644)

	// Comprimir se solicitado
	if compress {
		cmd := exec.Command("tar", "-czf", output, "-C", filepath.Dir(tempDir), filepath.Base(tempDir))
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to compress backup: %w", err)
		}
	} else {
		cmd := exec.Command("cp", "-r", tempDir, output)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to copy backup: %w", err)
		}
	}

	fmt.Printf("✅ Backup created: %s\n", output)
	return nil
}

func restoreNodes(backupFile string, force bool) error {
	if !force {
		fmt.Printf("⚠️  This will restore from backup: %s\n", backupFile)
		fmt.Print("Are you sure? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Restore cancelled.")
			return nil
		}
	}

	fmt.Println("🔄 Restoring node configurations from backup...")

	syntropyDir := getSyntropyDir()
	tempDir := filepath.Join(syntropyDir, "backups", "restore_temp")

	// Extrair backup
	if strings.HasSuffix(backupFile, ".tar.gz") {
		cmd := exec.Command("tar", "-xzf", backupFile, "-C", filepath.Dir(tempDir))
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to extract backup: %w", err)
		}
	} else {
		cmd := exec.Command("cp", "-r", backupFile, tempDir)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to copy backup: %w", err)
		}
	}
	defer os.RemoveAll(tempDir)

	// Restaurar componentes
	components := []string{"nodes", "keys", "config", "cache"}
	for _, component := range components {
		src := filepath.Join(tempDir, component)
		dst := filepath.Join(syntropyDir, component)
		
		if _, err := os.Stat(src); err == nil {
			// Fazer backup do existente
			if _, err := os.Stat(dst); err == nil {
				backupDst := dst + ".backup." + time.Now().Format("20060102_150405")
				exec.Command("mv", dst, backupDst).Run()
			}
			
			cmd := exec.Command("cp", "-r", src, dst)
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to restore %s: %w", component, err)
			}
		}
	}

	fmt.Println("✅ Restore completed successfully")
	return nil
}

func checkNodeHealth(format string, watch bool) error {
	fmt.Println("❤️  Checking health of all nodes...")

	nodes, err := loadAllNodes()
	if err != nil {
		return fmt.Errorf("failed to load nodes: %w", err)
	}

	healthResults := []HealthResult{}

	for _, node := range nodes {
		health := checkSingleNodeHealth(node)
		healthResults = append(healthResults, health)
	}

	if watch {
		return watchHealthResults(healthResults, format)
	}

	switch format {
	case "json":
		return outputHealthJSON(healthResults)
	case "yaml":
		return outputHealthYAML(healthResults)
	default:
		return outputHealthTable(healthResults)
	}
}

// Estruturas e funções auxiliares

type DiscoveredNode struct {
	IP       string
	Hostname string
	Port     int
	SSH      bool
	Latency  int
}

type HealthResult struct {
	NodeName    string `json:"node_name"`
	Status      string `json:"status"`
	IPAddress   string `json:"ip_address"`
	SSH         bool   `json:"ssh_accessible"`
	Latency     int    `json:"latency_ms"`
	LastSeen    string `json:"last_seen"`
	Error       string `json:"error,omitempty"`
}

func loadAllNodes() ([]NodeInfo, error) {
	syntropyDir := getSyntropyDir()
	nodesDir := filepath.Join(syntropyDir, "nodes")
	
	files, err := filepath.Glob(filepath.Join(nodesDir, "*.json"))
	if err != nil {
		return nil, err
	}

	nodes := []NodeInfo{}
	for _, file := range files {
		node, err := loadNodeFromFile(file)
		if err != nil {
			continue // Ignorar arquivos corrompidos
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func loadNode(nodeName string) (NodeInfo, error) {
	syntropyDir := getSyntropyDir()
	nodeFile := filepath.Join(syntropyDir, "nodes", nodeName+".json")
	return loadNodeFromFile(nodeFile)
}

func loadNodeFromFile(filePath string) (NodeInfo, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return NodeInfo{}, err
	}

	var node NodeInfo
	if err := json.Unmarshal(data, &node); err != nil {
		return NodeInfo{}, err
	}

	return node, nil
}

func saveNode(node NodeInfo) error {
	syntropyDir := getSyntropyDir()
	nodesDir := filepath.Join(syntropyDir, "nodes")
	os.MkdirAll(nodesDir, 0755)

	nodeFile := filepath.Join(nodesDir, node.Name+".json")
	data, err := json.MarshalIndent(node, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(nodeFile, data, 0644)
}

func getNodeKeyFile(nodeName string) string {
	syntropyDir := getSyntropyDir()
	return filepath.Join(syntropyDir, "keys", nodeName+"_owner.key")
}

func filterNodes(nodes []NodeInfo, filter string) []NodeInfo {
	filtered := []NodeInfo{}
	for _, node := range nodes {
		if strings.ToLower(node.Status) == strings.ToLower(filter) {
			filtered = append(filtered, node)
		}
	}
	return filtered
}

func sortNodes(nodes []NodeInfo, sortBy string) {
	switch sortBy {
	case "created":
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].Created > nodes[j].Created
		})
	case "last_seen":
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].LastSeen > nodes[j].LastSeen
		})
	case "status":
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].Status < nodes[j].Status
		})
	default: // name
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].Name < nodes[j].Name
		})
	}
}

func updateNodeStatus(node *NodeInfo) {
	if node.Network.IPAddress == "" {
		node.Status = "unknown"
		return
	}

	// Testar conectividade SSH
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:22", node.Network.IPAddress), 5*time.Second)
	if err != nil {
		node.Status = "offline"
		node.LastSeen = time.Now().Format(time.RFC3339)
		return
	}
	conn.Close()

	node.Status = "online"
	node.LastSeen = time.Now().Format(time.RFC3339)
}

func getDefaultNetworks() []string {
	// Tentar detectar redes locais
	cmd := exec.Command("ip", "route")
	output, err := cmd.Output()
	if err != nil {
		return []string{"192.168.1.0/24", "10.0.0.0/24"}
	}

	networks := []string{}
	lines := strings.Split(string(output), "\n")
	
	for _, line := range lines {
		if strings.Contains(line, "192.168.") || strings.Contains(line, "10.") {
			parts := strings.Fields(line)
			if len(parts) > 0 && strings.Contains(parts[0], "/") {
				networks = append(networks, parts[0])
			}
		}
	}

	if len(networks) == 0 {
		return []string{"192.168.1.0/24", "10.0.0.0/24"}
	}

	return networks[:min(len(networks), 3)] // Limitar a 3 redes
}

func scanNetwork(network string, port, timeout int) []DiscoveredNode {
	// Implementação simplificada - em produção usar nmap
	nodes := []DiscoveredNode{}
	
	// Para demonstração, retornar alguns IPs fictícios
	// Em produção, usar nmap para scan real
	return nodes
}

func identifySyntropyNodes(discovered []DiscoveredNode) []DiscoveredNode {
	// Implementação simplificada - em produção verificar se é nó Syntropy
	return discovered
}

func updateOrCreateNode(discovered DiscoveredNode) error {
	// Implementação simplificada
	return nil
}

func updateDiscoveryCache(discovered []DiscoveredNode) {
	// Implementação simplificada
}

func checkSingleNodeHealth(node NodeInfo) HealthResult {
	result := HealthResult{
		NodeName:  node.Name,
		IPAddress: node.Network.IPAddress,
		LastSeen:  node.LastSeen,
	}

	if node.Network.IPAddress == "" {
		result.Status = "unknown"
		result.Error = "No IP address"
		return result
	}

	// Testar conectividade
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:22", node.Network.IPAddress), 5*time.Second)
	if err != nil {
		result.Status = "offline"
		result.SSH = false
		result.Error = err.Error()
		return result
	}
	conn.Close()

	result.Status = "online"
	result.SSH = true
	result.Latency = 10 // Placeholder

	return result
}

// Funções de output

func outputNodesTable(nodes []NodeInfo) error {
	fmt.Printf("%-20s %-15s %-10s %-20s %s\n", "NAME", "IP ADDRESS", "STATUS", "LAST SEEN", "DESCRIPTION")
	fmt.Println(strings.Repeat("-", 80))

	for _, node := range nodes {
		lastSeen := "never"
		if node.LastSeen != "" {
			lastSeen = node.LastSeen[:10] // Apenas data
		}
		fmt.Printf("%-20s %-15s %-10s %-20s %s\n",
			node.Name, node.Network.IPAddress, node.Status, lastSeen, node.Description)
	}

	return nil
}

func outputNodesJSON(nodes []NodeInfo) error {
	data, err := json.MarshalIndent(nodes, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func outputNodesYAML(nodes []NodeInfo) error {
	// Implementação simplificada
	fmt.Println("nodes:")
	for _, node := range nodes {
		fmt.Printf("- name: %s\n", node.Name)
		fmt.Printf("  ip_address: %s\n", node.Network.IPAddress)
		fmt.Printf("  status: %s\n", node.Status)
	}
	return nil
}

func outputNodeTable(node NodeInfo) error {
	fmt.Printf("Node: %s\n", node.Name)
	fmt.Printf("Description: %s\n", node.Description)
	fmt.Printf("Status: %s\n", node.Status)
	fmt.Printf("IP Address: %s\n", node.Network.IPAddress)
	fmt.Printf("Last Seen: %s\n", node.LastSeen)
	fmt.Printf("Created: %s\n", node.Created)
	return nil
}

func outputNodeJSON(node NodeInfo) error {
	data, err := json.MarshalIndent(node, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func outputNodeYAML(node NodeInfo) error {
	// Implementação simplificada
	fmt.Printf("name: %s\n", node.Name)
	fmt.Printf("description: %s\n", node.Description)
	fmt.Printf("status: %s\n", node.Status)
	return nil
}

func outputHealthTable(results []HealthResult) error {
	fmt.Printf("%-20s %-15s %-10s %-8s %-10s %s\n", "NODE", "IP ADDRESS", "STATUS", "SSH", "LATENCY", "ERROR")
	fmt.Println(strings.Repeat("-", 80))

	for _, result := range results {
		ssh := "❌"
		if result.SSH {
			ssh = "✅"
		}
		fmt.Printf("%-20s %-15s %-10s %-8s %-10d %s\n",
			result.NodeName, result.IPAddress, result.Status, ssh, result.Latency, result.Error)
	}

	return nil
}

func outputHealthJSON(results []HealthResult) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func outputHealthYAML(results []HealthResult) error {
	// Implementação simplificada
	fmt.Println("health_results:")
	for _, result := range results {
		fmt.Printf("- node_name: %s\n", result.NodeName)
		fmt.Printf("  status: %s\n", result.Status)
		fmt.Printf("  ssh_accessible: %t\n", result.SSH)
	}
	return nil
}

func watchNodeStatus(nodes []NodeInfo, format string) error {
	// Implementação simplificada - em produção usar ticker
	fmt.Println("Watching node status (press Ctrl+C to stop)...")
	return nil
}

func watchHealthResults(results []HealthResult, format string) error {
	// Implementação simplificada - em produção usar ticker
	fmt.Println("Watching health results (press Ctrl+C to stop)...")
	return nil
}

func getHostname() string {
	hostname, _ := os.Hostname()
	return hostname
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

