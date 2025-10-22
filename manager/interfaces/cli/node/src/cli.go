package node

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"node-component/src/internal/constants"
	"node-component/src/internal/helpers"
	"node-component/src/internal/types"

	"github.com/spf13/cobra"
)

// CLICommands provides CLI commands for the node component
type CLICommands struct {
	nodeManager *NodeManager
	logger      types.Logger
}

// NewCLICommands creates a new CLI commands handler
func NewCLICommands(nodeManager *NodeManager, logger types.Logger) *CLICommands {
	return &CLICommands{
		nodeManager: nodeManager,
		logger:      logger,
	}
}

// RegisterCommands registers all node CLI commands
func (cli *CLICommands) RegisterCommands(nodeCmd *cobra.Command) {
	// Add subcommands to the existing node command
	nodeCmd.AddCommand(cli.createNodeCmd())
	nodeCmd.AddCommand(cli.listNodesCmd())
	nodeCmd.AddCommand(cli.nodeStatusCmd())
	nodeCmd.AddCommand(cli.nodeLogsCmd())
	nodeCmd.AddCommand(cli.removeNodeCmd())
	nodeCmd.AddCommand(cli.startListenerCmd())
	nodeCmd.AddCommand(cli.stopListenerCmd())
}

// createNodeCmd creates the node create command
func (cli *CLICommands) createNodeCmd() *cobra.Command {
	var (
		ubuntuVersion    string
		devicePath       string
		isoPath          string
		nodeID           string
		skipUSBDetection bool
		skipISODownload  bool
		skipCloudInit    bool
		skipUSBWrite     bool
		forceOverwrite   bool
		autoStart        bool
		interactive      bool
		listDevices      bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new Syntropy node",
		Long: `Create a new Syntropy node with plug-and-play USB bootable.

This command will:
1. Generate node configuration (NodeID, SSH keys, certificates)
2. Detect and prepare USB device
3. Download Ubuntu Server ISO (if needed)
4. Generate cloud-init configuration
5. Create bootable USB with auto-registration
6. Start listener for node registration

Examples:
  syntropy node create                           # Interactive mode
  syntropy node create --list-devices            # List available USB devices
  syntropy node create --ubuntu-version 24.04    # Specify Ubuntu version
  syntropy node create --device /dev/sdb         # Specify USB device
  syntropy node create --auto-start              # Auto-start listener`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// If --list-devices flag is set, only list devices without creating node
			if listDevices {
				return cli.handleListUSBDevices(cmd.Context())
			}

			return cli.handleCreateNode(cmd.Context(), &CreateNodeOptions{
				UbuntuVersion:    ubuntuVersion,
				DevicePath:       devicePath,
				ISOPath:          isoPath,
				NodeID:           nodeID,
				SkipUSBDetection: skipUSBDetection,
				SkipISODownload:  skipISODownload,
				SkipCloudInit:    skipCloudInit,
				SkipUSBWrite:     skipUSBWrite,
				ForceOverwrite:   forceOverwrite,
				AutoStart:        autoStart,
				Interactive:      interactive,
			})
		},
	}

	// Add flags
	cmd.Flags().StringVar(&ubuntuVersion, "ubuntu-version", "24.04", "Ubuntu Server version (e.g., 24.04, 22.04)")
	cmd.Flags().StringVar(&devicePath, "device", "", "USB device path (e.g., /dev/sdb on Linux, D: on Windows)")
	cmd.Flags().StringVar(&isoPath, "iso", "", "Path to local Ubuntu ISO file")
	cmd.Flags().StringVar(&nodeID, "node-id", "", "Pre-defined node ID (auto-generated if not specified)")
	cmd.Flags().BoolVar(&skipUSBDetection, "skip-usb-detection", false, "Skip USB device detection")
	cmd.Flags().BoolVar(&skipISODownload, "skip-iso-download", false, "Skip ISO download")
	cmd.Flags().BoolVar(&skipCloudInit, "skip-cloud-init", false, "Skip cloud-init generation")
	cmd.Flags().BoolVar(&skipUSBWrite, "skip-usb-write", false, "Skip USB writing")
	cmd.Flags().BoolVar(&forceOverwrite, "force", false, "Force overwrite USB device")
	cmd.Flags().BoolVar(&autoStart, "auto-start", true, "Auto-start listener after creation")
	cmd.Flags().BoolVar(&interactive, "interactive", false, "Interactive mode with prompts")
	cmd.Flags().BoolVar(&listDevices, "list-devices", false, "List available USB devices without creating node")

	return cmd
}

// listNodesCmd creates the node list command
func (cli *CLICommands) listNodesCmd() *cobra.Command {
	var (
		state   string
		verbose bool
		json    bool
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all nodes",
		Long: `List all nodes in the system.

Examples:
  syntropy node list                    # List all nodes
  syntropy node list --state active     # List only active nodes
  syntropy node list --verbose          # Show detailed information
  syntropy node list --json             # Output in JSON format`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cli.handleListNodes(cmd.Context(), &ListNodeOptions{
				State:   state,
				Verbose: verbose,
				JSON:    json,
			})
		},
	}

	// Add flags
	cmd.Flags().StringVar(&state, "state", "", "Filter by node state (pending, active, inactive)")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "Show detailed information")
	cmd.Flags().BoolVar(&json, "json", false, "Output in JSON format")

	return cmd
}

// nodeStatusCmd creates the node status command
func (cli *CLICommands) nodeStatusCmd() *cobra.Command {
	var (
		json bool
	)

	cmd := &cobra.Command{
		Use:   "status <node-id>",
		Short: "Show node status and information",
		Long: `Show detailed status and information for a specific node.

Examples:
  syntropy node status node-01          # Show status for node-01
  syntropy node status node-01 --json   # Output in JSON format`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cli.handleNodeStatus(cmd.Context(), args[0], &NodeStatusOptions{
				JSON: json,
			})
		},
	}

	// Add flags
	cmd.Flags().BoolVar(&json, "json", false, "Output in JSON format")

	return cmd
}

// nodeLogsCmd creates the node logs command
func (cli *CLICommands) nodeLogsCmd() *cobra.Command {
	var (
		lines  int
		follow bool
	)

	cmd := &cobra.Command{
		Use:   "logs <node-id>",
		Short: "Show node logs",
		Long: `Show logs for a specific node.

Examples:
  syntropy node logs node-01            # Show recent logs
  syntropy node logs node-01 --lines 100 # Show last 100 lines
  syntropy node logs node-01 --follow   # Follow logs in real-time`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cli.handleNodeLogs(cmd.Context(), args[0], &NodeLogsOptions{
				Lines:  lines,
				Follow: follow,
			})
		},
	}

	// Add flags
	cmd.Flags().IntVar(&lines, "lines", 50, "Number of lines to show")
	cmd.Flags().BoolVar(&follow, "follow", false, "Follow logs in real-time")

	return cmd
}

// removeNodeCmd creates the node remove command
func (cli *CLICommands) removeNodeCmd() *cobra.Command {
	var (
		force bool
	)

	cmd := &cobra.Command{
		Use:   "remove <node-id>",
		Short: "Remove a node from the system",
		Long: `Remove a node from the system.

This will:
1. Stop heartbeat monitoring
2. Remove node from active/pending/inactive lists
3. Clean up node files and configuration

Examples:
  syntropy node remove node-01          # Remove node-01
  syntropy node remove node-01 --force  # Force removal without confirmation`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cli.handleRemoveNode(cmd.Context(), args[0], &RemoveNodeOptions{
				Force: force,
			})
		},
	}

	// Add flags
	cmd.Flags().BoolVar(&force, "force", false, "Force removal without confirmation")

	return cmd
}

// startListenerCmd creates the start listener command
func (cli *CLICommands) startListenerCmd() *cobra.Command {
	var (
		port int
	)

	cmd := &cobra.Command{
		Use:   "start-listener",
		Short: "Start the node registration listener",
		Long: `Start the TCP listener for node registration.

Examples:
  syntropy node start-listener          # Start listener on default port 51000
  syntropy node start-listener --port 51001 # Start listener on custom port`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cli.handleStartListener(cmd.Context(), &StartListenerOptions{
				Port: port,
			})
		},
	}

	// Add flags
	cmd.Flags().IntVar(&port, "port", 51000, "Port to listen on")

	return cmd
}

// stopListenerCmd creates the stop listener command
func (cli *CLICommands) stopListenerCmd() *cobra.Command {
	var (
		port int
	)

	cmd := &cobra.Command{
		Use:   "stop-listener",
		Short: "Stop the node registration listener",
		Long: `Stop the TCP listener for node registration.

Examples:
  syntropy node stop-listener           # Stop listener on default port 51000
  syntropy node stop-listener --port 51001 # Stop listener on custom port`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cli.handleStopListener(cmd.Context(), &StopListenerOptions{
				Port: port,
			})
		},
	}

	// Add flags
	cmd.Flags().IntVar(&port, "port", 51000, "Port to stop listening on")

	return cmd
}

// Command option types

// CreateNodeOptions represents options for creating a node
type CreateNodeOptions struct {
	UbuntuVersion    string
	DevicePath       string
	ISOPath          string
	NodeID           string
	SkipUSBDetection bool
	SkipISODownload  bool
	SkipCloudInit    bool
	SkipUSBWrite     bool
	ForceOverwrite   bool
	AutoStart        bool
	Interactive      bool
}

// ListNodeOptions represents options for listing nodes
type ListNodeOptions struct {
	State   string
	Verbose bool
	JSON    bool
}

// NodeStatusOptions represents options for showing node status
type NodeStatusOptions struct {
	JSON bool
}

// NodeLogsOptions represents options for showing node logs
type NodeLogsOptions struct {
	Lines  int
	Follow bool
}

// RemoveNodeOptions represents options for removing a node
type RemoveNodeOptions struct {
	Force bool
}

// StartListenerOptions represents options for starting listener
type StartListenerOptions struct {
	Port int
}

// StopListenerOptions represents options for stopping listener
type StopListenerOptions struct {
	Port int
}

// Command handlers

// handleCreateNode handles the create node command
func (cli *CLICommands) handleCreateNode(ctx context.Context, options *CreateNodeOptions) error {
	cli.logger.Info("Creating new node", "options", options)

	// Initialize node manager if needed
	if err := cli.nodeManager.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize node manager: %w", err)
	}

	// Create node
	result, err := cli.nodeManager.CreateNode(&CreateOptions{
		UbuntuVersion:    options.UbuntuVersion,
		DevicePath:       options.DevicePath,
		SkipUSBDetection: options.SkipUSBDetection,
		SkipISODownload:  options.SkipISODownload,
		SkipCloudInit:    options.SkipCloudInit,
		SkipUSBWrite:     options.SkipUSBWrite,
		ForceOverwrite:   options.ForceOverwrite,
		AutoStart:        options.AutoStart,
	})
	if err != nil {
		return fmt.Errorf("failed to create node: %w", err)
	}

	// Display results
	if result.Success {
		fmt.Printf("✅ Node created successfully!\n\n")
		fmt.Printf("Node ID: %s\n", result.NodeID)
		fmt.Printf("Device: %s\n", result.DevicePath)
		fmt.Printf("ISO: %s\n", result.ISOPath)
		fmt.Printf("Duration: %v\n", result.Duration)
		fmt.Printf("Steps completed: %d\n", len(result.StepsCompleted))

		if len(result.StepsCompleted) > 0 {
			fmt.Printf("Completed steps:\n")
			for _, step := range result.StepsCompleted {
				fmt.Printf("  ✅ %s\n", step)
			}
		}

		fmt.Printf("\n🚀 Your node is ready! Connect the USB to your hardware and boot.\n")
		fmt.Printf("The node will automatically register with this Command Station.\n")
	} else {
		fmt.Printf("❌ Node creation failed: %s\n", result.ErrorMessage)
		if len(result.StepsFailed) > 0 {
			fmt.Printf("Failed steps:\n")
			for _, step := range result.StepsFailed {
				fmt.Printf("  ❌ %s\n", step)
			}
		}
	}

	return nil
}

// handleListNodes handles the list nodes command
func (cli *CLICommands) handleListNodes(ctx context.Context, options *ListNodeOptions) error {
	cli.logger.Info("Listing nodes", "options", options)

	// Get nodes from node manager
	nodeList, err := cli.nodeManager.ListNodes()
	if err != nil {
		return fmt.Errorf("failed to list nodes: %w", err)
	}

	// Filter nodes by state if specified
	var nodes []*NodeStatus
	if options.State == "" {
		// Combine all nodes
		nodes = append(nodes, nodeList.Active...)
		nodes = append(nodes, nodeList.Pending...)
		nodes = append(nodes, nodeList.Inactive...)
	} else {
		// Filter by specific state
		switch options.State {
		case "active":
			nodes = nodeList.Active
		case "pending":
			nodes = nodeList.Pending
		case "inactive":
			nodes = nodeList.Inactive
		}
	}

	if options.JSON {
		// Output in JSON format
		return cli.outputNodesJSON(nodes)
	}

	// Output in table format
	return cli.outputNodesTable(nodes, options.Verbose)
}

// handleNodeStatus handles the node status command
func (cli *CLICommands) handleNodeStatus(ctx context.Context, nodeID string, options *NodeStatusOptions) error {
	cli.logger.Info("Getting node status", "node_id", nodeID)

	// Get node status
	status, err := cli.nodeManager.GetNodeStatus(nodeID)
	if err != nil {
		return fmt.Errorf("failed to get node status: %w", err)
	}

	if options.JSON {
		// Output in JSON format
		return cli.outputNodeStatusJSON(status)
	}

	// Output in table format
	return cli.outputNodeStatusTable(status)
}

// handleNodeLogs handles the node logs command
func (cli *CLICommands) handleNodeLogs(ctx context.Context, nodeID string, options *NodeLogsOptions) error {
	cli.logger.Info("Getting node logs", "node_id", nodeID, "lines", options.Lines)

	// Get node logs
	logOptions := &LogOptions{
		Lines:  options.Lines,
		Follow: options.Follow,
	}
	logs, err := cli.nodeManager.GetNodeLogs(nodeID, logOptions)
	if err != nil {
		return fmt.Errorf("failed to get node logs: %w", err)
	}

	// Output logs
	for _, log := range logs.Logs {
		fmt.Printf("%s\n", log)
	}

	return nil
}

// handleRemoveNode handles the remove node command
func (cli *CLICommands) handleRemoveNode(ctx context.Context, nodeID string, options *RemoveNodeOptions) error {
	cli.logger.Info("Removing node", "node_id", nodeID)

	if !options.Force {
		// Ask for confirmation
		fmt.Printf("Are you sure you want to remove node '%s'? (yes/no): ", nodeID)
		var confirm string
		fmt.Scanln(&confirm)
		if strings.ToLower(confirm) != "yes" {
			fmt.Println("Operation cancelled.")
			return nil
		}
	}

	// Remove node
	if err := cli.nodeManager.DeleteNode(nodeID); err != nil {
		return fmt.Errorf("failed to remove node: %w", err)
	}

	fmt.Printf("✅ Node '%s' removed successfully.\n", nodeID)
	return nil
}

// handleStartListener handles the start listener command
func (cli *CLICommands) handleStartListener(ctx context.Context, options *StartListenerOptions) error {
	cli.logger.Info("Starting listener", "port", options.Port)

	// Start listener
	if err := cli.nodeManager.StartRegistrationListener(); err != nil {
		return fmt.Errorf("failed to start listener: %w", err)
	}

	fmt.Printf("✅ Listener started on port %d\n", options.Port)
	fmt.Printf("Waiting for nodes to register...\n")
	return nil
}

// handleStopListener handles the stop listener command
func (cli *CLICommands) handleStopListener(ctx context.Context, options *StopListenerOptions) error {
	cli.logger.Info("Stopping listener", "port", options.Port)

	// Stop listener
	if err := cli.nodeManager.StopRegistrationListener(); err != nil {
		return fmt.Errorf("failed to stop listener: %w", err)
	}

	fmt.Printf("✅ Listener stopped on port %d\n", options.Port)
	return nil
}

// handleListUSBDevices handles the list USB devices command
func (cli *CLICommands) handleListUSBDevices(ctx context.Context) error {
	cli.logger.Info("Listing USB devices")

	// Initialize node manager if needed
	if err := cli.nodeManager.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize node manager: %w", err)
	}

	// Get USB devices
	devices, err := cli.nodeManager.ListUSBDevices(ctx)
	if err != nil {
		return fmt.Errorf("failed to list USB devices: %w", err)
	}

	// Display devices
	return cli.outputUSBDevicesTable(devices)
}

// Output helpers

// outputNodesJSON outputs nodes in JSON format
func (cli *CLICommands) outputNodesJSON(nodes []*NodeStatus) error {
	// TODO: Implement JSON output
	fmt.Printf("JSON output not yet implemented\n")
	return nil
}

// outputNodesTable outputs nodes in table format
func (cli *CLICommands) outputNodesTable(nodes []*NodeStatus, verbose bool) error {
	if len(nodes) == 0 {
		fmt.Println("No nodes found.")
		return nil
	}

	// Create table writer
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	// Print header
	if verbose {
		fmt.Fprintf(w, "NODE ID\tSTATUS\tCREATED\tLAST HEARTBEAT\tIP ADDRESS\tHARDWARE\n")
	} else {
		fmt.Fprintf(w, "NODE ID\tSTATUS\tCREATED\tLAST HEARTBEAT\n")
	}

	// Print nodes
	for _, node := range nodes {
		hardwareInfo := "N/A"
		if node.Hardware != nil {
			hardwareInfo = fmt.Sprintf("%d cores, %dGB RAM", node.Hardware.CPUCores, node.Hardware.MemoryGB)
		}

		if verbose {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
				node.NodeID,
				node.Status,
				node.CreatedAt.Format("2006-01-02 15:04:05"),
				node.LastHeartbeat.Format("2006-01-02 15:04:05"),
				node.IPAddress,
				hardwareInfo,
			)
		} else {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				node.NodeID,
				node.Status,
				node.CreatedAt.Format("2006-01-02 15:04:05"),
				node.LastHeartbeat.Format("2006-01-02 15:04:05"),
			)
		}
	}

	return nil
}

// outputNodeStatusJSON outputs node status in JSON format
func (cli *CLICommands) outputNodeStatusJSON(status *NodeStatus) error {
	// TODO: Implement JSON output
	fmt.Printf("JSON output not yet implemented\n")
	return nil
}

// outputNodeStatusTable outputs node status in table format
func (cli *CLICommands) outputNodeStatusTable(status *NodeStatus) error {
	fmt.Printf("Node Status: %s\n\n", status.NodeID)
	fmt.Printf("Status: %s\n", status.Status)
	fmt.Printf("Created: %s\n", status.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Registered: %s\n", status.RegisteredAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Last Heartbeat: %s\n", status.LastHeartbeat.Format("2006-01-02 15:04:05"))
	fmt.Printf("IP Address: %s\n", status.IPAddress)
	fmt.Printf("Uptime: %v\n", status.Uptime)

	if status.Hardware != nil {
		fmt.Printf("\nHardware:\n")
		fmt.Printf("  CPU Cores: %d\n", status.Hardware.CPUCores)
		fmt.Printf("  Memory: %dGB\n", status.Hardware.MemoryGB)
		fmt.Printf("  Disk: %dGB\n", status.Hardware.DiskGB)
		fmt.Printf("  Hostname: %s\n", status.Hardware.Hostname)
	}

	return nil
}

// outputUSBDevicesTable outputs USB devices in table format
func (cli *CLICommands) outputUSBDevicesTable(devices []types.USBDevice) error {
	if len(devices) == 0 {
		fmt.Println("No USB devices found.")
		fmt.Println("Please connect a USB device and try again.")
		return nil
	}

	// Create table writer
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	// Print header
	fmt.Printf("Found %d USB device(s):\n\n", len(devices))
	fmt.Fprintf(w, "DEVICE\tCAPACITY\tVENDOR\tMODEL\tREMOVABLE\tSYSTEM\tSTATUS\n")

	// Print devices
	for _, device := range devices {
		capacity := helpers.FormatBytes(device.Capacity)
		removable := "No"
		if device.IsRemovable {
			removable = "Yes"
		}
		system := "No"
		if device.IsSystem {
			system = "Yes"
		}

		// Determine status
		status := "Ready"
		if device.IsSystem {
			status = "System Device (Cannot Use)"
		} else if !device.IsRemovable {
			status = "Not Removable"
		} else if device.Capacity < constants.DefaultMinUSBCapacity {
			status = fmt.Sprintf("Too Small (Min %s)", helpers.FormatBytes(constants.DefaultMinUSBCapacity))
		} else if device.Capacity > constants.DefaultMaxUSBCapacity {
			status = fmt.Sprintf("Too Large (Max %s)", helpers.FormatBytes(constants.DefaultMaxUSBCapacity))
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			device.Path,
			capacity,
			device.Vendor,
			device.Model,
			removable,
			system,
			status,
		)
	}

	fmt.Println("\nRecommendation: Use devices with status 'Ready' for node creation.")
	return nil
}
