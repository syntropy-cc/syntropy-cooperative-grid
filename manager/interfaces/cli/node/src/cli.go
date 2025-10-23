package node

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

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
	nodeCmd.AddCommand(cli.listenerCmd())
	nodeCmd.AddCommand(cli.isoCmd())
}

// createNodeCmd creates the node create command
func (cli *CLICommands) createNodeCmd() *cobra.Command {
	var (
		ubuntuVersion     string
		devicePath        string
		isoPath           string
		isoURL            string
		nodeID            string
		skipUSBDetection  bool
		skipISODownload   bool
		skipCloudInit     bool
		skipUSBWrite      bool
		forceOverwrite    bool
		autoStart         bool
		interactive       bool
		listDevices       bool
		skipISOValidation bool
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
				UbuntuVersion:     ubuntuVersion,
				DevicePath:        devicePath,
				ISOPath:           isoPath,
				ISOURL:            isoURL,
				NodeID:            nodeID,
				SkipUSBDetection:  skipUSBDetection,
				SkipISODownload:   skipISODownload,
				SkipCloudInit:     skipCloudInit,
				SkipUSBWrite:      skipUSBWrite,
				ForceOverwrite:    forceOverwrite,
				AutoStart:         autoStart,
				Interactive:       interactive,
				SkipISOValidation: skipISOValidation,
			})
		},
	}

	// Add flags
	cmd.Flags().StringVar(&ubuntuVersion, "ubuntu-version", "24.04", "Ubuntu Server version (e.g., 24.04, 22.04)")
	cmd.Flags().StringVar(&devicePath, "device", "", "USB device path (e.g., /dev/sdb on Linux, D: on Windows)")
	cmd.Flags().StringVar(&isoPath, "iso", "", "Path to local Ubuntu ISO file")
	cmd.Flags().StringVar(&isoURL, "iso-url", "", "Custom ISO download URL")
	cmd.Flags().StringVar(&nodeID, "node-id", "", "Pre-defined node ID (auto-generated if not specified)")
	cmd.Flags().BoolVar(&skipUSBDetection, "skip-usb-detection", false, "Skip USB device detection")
	cmd.Flags().BoolVar(&skipISODownload, "skip-iso-download", false, "Skip ISO download")
	cmd.Flags().BoolVar(&skipCloudInit, "skip-cloud-init", false, "Skip cloud-init generation")
	cmd.Flags().BoolVar(&skipUSBWrite, "skip-usb-write", false, "Skip USB writing")
	cmd.Flags().BoolVar(&forceOverwrite, "force", false, "Force overwrite USB device")
	cmd.Flags().BoolVar(&autoStart, "auto-start", true, "Auto-start listener after creation")
	cmd.Flags().BoolVar(&interactive, "interactive", false, "Interactive mode with prompts")
	cmd.Flags().BoolVar(&listDevices, "list-devices", false, "List available USB devices without creating node")
	cmd.Flags().BoolVar(&skipISOValidation, "skip-iso-validation", false, "Skip ISO SHA256 validation (temporary feature)")

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

// listenerCmd creates the listener command with subcommands
func (cli *CLICommands) listenerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listener",
		Short: "Manage node registration listener",
		Long: `Manage the TCP listener for node registration.

The listener accepts incoming connections from nodes that are booting up
and trying to register with the Command Station.`,
	}

	// Add subcommands
	cmd.AddCommand(cli.listenerStartCmd())
	cmd.AddCommand(cli.listenerStopCmd())
	cmd.AddCommand(cli.listenerStatusCmd())

	return cmd
}

// listenerStartCmd creates the listener start command
func (cli *CLICommands) listenerStartCmd() *cobra.Command {
	var (
		port int
	)

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the node registration listener",
		Long: `Start the TCP listener for node registration.

The listener will accept incoming connections from nodes and process
their registration handshakes. The command will block until Ctrl+C is pressed.

Examples:
  syntropy node listener start              # Start listener on default port 51000
  syntropy node listener start --port 51001 # Start listener on custom port`,
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

// listenerStopCmd creates the listener stop command
func (cli *CLICommands) listenerStopCmd() *cobra.Command {
	var (
		port int
	)

	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop the node registration listener",
		Long: `Stop the TCP listener for node registration.

This will close all active connections and stop accepting new node registrations.

Examples:
  syntropy node listener stop               # Stop listener on default port 51000
  syntropy node listener stop --port 51001  # Stop listener on custom port`,
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

// listenerStatusCmd creates the listener status command
func (cli *CLICommands) listenerStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show listener status",
		Long:  `Display the current status of the node registration listener.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cli.handleListenerStatus(cmd.Context())
		},
	}
	return cmd
}

// Command option types

// CreateNodeOptions represents options for creating a node
type CreateNodeOptions struct {
	UbuntuVersion     string
	DevicePath        string
	ISOPath           string
	ISOURL            string
	NodeID            string
	SkipUSBDetection  bool
	SkipISODownload   bool
	SkipCloudInit     bool
	SkipUSBWrite      bool
	ForceOverwrite    bool
	AutoStart         bool
	Interactive       bool
	SkipISOValidation bool
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
		UbuntuVersion:     options.UbuntuVersion,
		DevicePath:        options.DevicePath,
		ISOPath:           options.ISOPath,
		ISOURL:            options.ISOURL,
		SkipUSBDetection:  options.SkipUSBDetection,
		SkipISODownload:   options.SkipISODownload,
		SkipCloudInit:     options.SkipCloudInit,
		SkipUSBWrite:      options.SkipUSBWrite,
		ForceOverwrite:    options.ForceOverwrite,
		AutoStart:         options.AutoStart,
		SkipISOValidation: options.SkipISOValidation,
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

		fmt.Printf("\n📋 Next Steps:\n")
		fmt.Printf("1. Connect the USB device to your target hardware\n")
		fmt.Printf("2. Boot from USB (may require BIOS/UEFI configuration)\n")
		fmt.Printf("3. Wait for Ubuntu installation (~5-10 minutes)\n")
		fmt.Printf("4. Node will automatically register when ready\n\n")

		// Auto-start listener
		if !cli.nodeManager.IsListenerRunning() {
			fmt.Printf("🚀 Starting registration listener on port 51000...\n")
			if err := cli.nodeManager.StartRegistrationListener(); err != nil {
				return fmt.Errorf("failed to start listener: %w", err)
			}
			fmt.Printf("✅ Listener started successfully\n\n")
		} else {
			fmt.Printf("ℹ️  Registration listener is already running\n\n")
		}

		// Wait for registration
		fmt.Printf("⏳ Waiting for node '%s' to register (timeout: 60 minutes)...\n", result.NodeID)
		fmt.Printf("💡 Press Ctrl+C to stop waiting (listener will continue running)\n\n")

		waitCtx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
		defer cancel()

		if err := cli.nodeManager.WaitForNodeRegistration(waitCtx, result.NodeID, 60*time.Minute); err != nil {
			if errors.Is(err, context.Canceled) {
				fmt.Printf("\n⚠️  Wait canceled by user\n")
				fmt.Printf("💡 Node will still register when ready\n")
				fmt.Printf("💡 Use 'syntropy node list' to check status\n")
				fmt.Printf("💡 Use 'syntropy node listener stop' to stop the listener\n")
				return nil
			}
			if errors.Is(err, context.DeadlineExceeded) {
				fmt.Printf("\n⏰ Registration timeout reached\n")
				fmt.Printf("💡 Node might still be installing or booting\n")
				fmt.Printf("💡 Use 'syntropy node list' to check status later\n")
				fmt.Printf("💡 Listener is still running and accepting connections\n")
				return nil
			}
			return fmt.Errorf("wait failed: %w", err)
		}

		fmt.Printf("\n🎉 Node '%s' registered successfully!\n", result.NodeID)
		fmt.Printf("💡 Use 'syntropy node status %s' to view details\n", result.NodeID)
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

	// Initialize node manager
	if err := cli.nodeManager.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	// Check if already running
	if cli.nodeManager.IsListenerRunning() {
		fmt.Printf("ℹ️  Listener is already running on port %d\n", options.Port)
		return nil
	}

	// Start listener
	if err := cli.nodeManager.StartRegistrationListener(); err != nil {
		return fmt.Errorf("failed to start listener: %w", err)
	}

	fmt.Printf("✅ Listener started on port %d\n", options.Port)
	fmt.Printf("⏳ Waiting for nodes to register...\n")
	fmt.Printf("💡 Press Ctrl+C to stop\n\n")

	// Block until Ctrl+C
	<-ctx.Done()
	fmt.Printf("\n🛑 Stopping listener...\n")

	return cli.nodeManager.StopRegistrationListener()
}

// handleStopListener handles the stop listener command
func (cli *CLICommands) handleStopListener(ctx context.Context, options *StopListenerOptions) error {
	cli.logger.Info("Stopping listener", "port", options.Port)

	// Initialize node manager
	if err := cli.nodeManager.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	// Check if running
	if !cli.nodeManager.IsListenerRunning() {
		fmt.Printf("ℹ️  Listener is not running\n")
		return nil
	}

	// Stop listener
	if err := cli.nodeManager.StopRegistrationListener(); err != nil {
		return fmt.Errorf("failed to stop listener: %w", err)
	}

	fmt.Printf("✅ Listener stopped successfully\n")
	return nil
}

// handleListenerStatus handles the listener status command
func (cli *CLICommands) handleListenerStatus(ctx context.Context) error {
	if err := cli.nodeManager.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	running := cli.nodeManager.IsListenerRunning()

	if running {
		fmt.Printf("✅ Listener Status: RUNNING\n")
		fmt.Printf("   Port: 51000\n")
		fmt.Printf("   Ready to accept node registrations\n")
	} else {
		fmt.Printf("❌ Listener Status: STOPPED\n")
		fmt.Printf("💡 Use 'syntropy node listener start' to start\n")
	}

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

// isoCmd creates the ISO management command
func (cli *CLICommands) isoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "iso",
		Short: "Manage ISO downloads and configuration",
		Long:  `Manage Ubuntu Server ISO downloads, test URLs, and configure sources.`,
	}

	cmd.AddCommand(cli.isoTestURLsCmd())
	cmd.AddCommand(cli.isoListSourcesCmd())
	cmd.AddCommand(cli.isoConfigCmd())
	cmd.AddCommand(cli.isoTestSingleURLCmd())

	return cmd
}

// isoTestURLsCmd tests all configured ISO URLs
func (cli *CLICommands) isoTestURLsCmd() *cobra.Command {
	var version string

	cmd := &cobra.Command{
		Use:   "test-urls",
		Short: "Test all configured ISO download URLs",
		Long:  `Test accessibility of all configured ISO download URLs for a specific version.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cli.handleTestISOURLs(cmd.Context(), version)
		},
	}

	cmd.Flags().StringVar(&version, "version", "24.04", "Ubuntu version to test")
	return cmd
}

// isoListSourcesCmd lists all ISO sources
func (cli *CLICommands) isoListSourcesCmd() *cobra.Command {
	var version string

	cmd := &cobra.Command{
		Use:   "list-sources",
		Short: "List all configured ISO download sources",
		Long:  `Display all configured ISO download sources in priority order.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cli.handleListISOSources(cmd.Context(), version)
		},
	}

	cmd.Flags().StringVar(&version, "version", "24.04", "Ubuntu version")
	return cmd
}

// isoConfigCmd shows ISO configuration
func (cli *CLICommands) isoConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Show ISO download configuration",
		Long:  `Display current ISO download configuration from manager.yaml.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cli.handleShowISOConfig(cmd.Context())
		},
	}

	return cmd
}

// isoTestSingleURLCmd tests a specific URL
func (cli *CLICommands) isoTestSingleURLCmd() *cobra.Command {
	var url string

	cmd := &cobra.Command{
		Use:   "test-url",
		Short: "Test a specific ISO download URL",
		Long:  `Test accessibility and extract download URL from a specific URL.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cli.handleTestSingleURL(cmd.Context(), url)
		},
	}

	cmd.Flags().StringVar(&url, "url", "", "URL to test (required)")
	cmd.MarkFlagRequired("url")
	return cmd
}

// handleTestISOURLs tests all ISO URLs
func (cli *CLICommands) handleTestISOURLs(ctx context.Context, version string) error {
	fmt.Printf("🔍 Testando URLs para Ubuntu %s...\n\n", version)

	// Criar downloader temporário
	isoDownloader := NewISODownloader(cli.logger)
	urls, err := isoDownloader.buildDownloadURLs(version, "")
	if err != nil {
		return err
	}

	fmt.Printf("Total de URLs a testar: %d\n\n", len(urls))

	for i, url := range urls {
		fmt.Printf("%d. Testando: %s\n", i+1, url)
		if err := isoDownloader.validateURL(ctx, url); err != nil {
			fmt.Printf("   ❌ Não disponível: %v\n\n", err)
		} else {
			fmt.Printf("   ✅ Disponível\n\n")
		}
	}

	return nil
}

// handleListISOSources lists all ISO sources
func (cli *CLICommands) handleListISOSources(ctx context.Context, version string) error {
	isoDownloader := NewISODownloader(cli.logger)
	urls, err := isoDownloader.buildDownloadURLs(version, "")
	if err != nil {
		return err
	}

	fmt.Printf("📋 Fontes de ISO configuradas para Ubuntu %s\n\n", version)
	fmt.Printf("Ordem de prioridade:\n")
	for i, url := range urls {
		fmt.Printf("  %d. %s\n", i+1, url)
	}

	fmt.Printf("\n💡 URLs são tentadas nesta ordem até encontrar uma disponível\n")
	return nil
}

// handleTestSingleURL tests a specific URL
func (cli *CLICommands) handleTestSingleURL(ctx context.Context, url string) error {
	fmt.Printf("🔍 Testando URL específico: %s\n\n", url)

	// Criar downloader temporário
	isoDownloader := NewISODownloader(cli.logger)

	// Testar validação básica
	fmt.Printf("1. Testando acessibilidade básica...\n")
	if err := isoDownloader.validateURL(ctx, url); err != nil {
		fmt.Printf("   ❌ URL não acessível: %v\n", err)
	} else {
		fmt.Printf("   ✅ URL acessível\n")
	}

	// Se for uma página de download do Ubuntu, tentar extrair o link real
	if strings.Contains(url, "ubuntu.com/download") || strings.Contains(url, "thank-you") {
		fmt.Printf("\n2. Tentando extrair URL de download real...\n")
		realURL, err := isoDownloader.extractDownloadURL(ctx, url)
		if err != nil {
			fmt.Printf("   ❌ Falha ao extrair URL: %v\n", err)
		} else {
			fmt.Printf("   ✅ URL extraído: %s\n", realURL)

			// Testar o URL extraído
			fmt.Printf("\n3. Testando URL extraído...\n")
			if err := isoDownloader.validateURL(ctx, realURL); err != nil {
				fmt.Printf("   ❌ URL extraído não acessível: %v\n", err)
			} else {
				fmt.Printf("   ✅ URL extraído acessível\n")
			}
		}
	}

	return nil
}

// handleShowISOConfig shows ISO configuration
func (cli *CLICommands) handleShowISOConfig(ctx context.Context) error {
	isoDownloader := NewISODownloader(cli.logger)
	config, err := isoDownloader.loadISOConfig()
	if err != nil {
		return err
	}

	fmt.Printf("⚙️  Configuração de Download de ISOs\n\n")
	fmt.Printf("Arquivo: ~/.syntropy/config/manager.yaml\n\n")

	fmt.Printf("URLs Personalizadas:\n")
	if len(config.CustomURLs) == 0 {
		fmt.Printf("  (nenhuma configurada)\n")
	} else {
		for _, url := range config.CustomURLs {
			fmt.Printf("  - %s\n", url)
		}
	}

	fmt.Printf("\nMirrors Preferidos:\n")
	for _, mirror := range config.PreferredMirrors {
		fmt.Printf("  - %s\n", mirror)
	}

	fmt.Printf("\nFallback Automático: %v\n", config.EnableAutoFallback)
	fmt.Printf("Máximo de Tentativas: %d\n", config.MaxRetries)
	fmt.Printf("Timeout: %v\n", config.Timeout)

	return nil
}
