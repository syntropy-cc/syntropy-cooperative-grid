# Syntropy Cooperative Grid - Base Node Image Builder
# Packer template for creating Ubuntu-based node images

packer {
  required_plugins {
    qemu = {
      version = "~> 1"
      source  = "github.com/hashicorp/qemu"
    }
    virtualbox = {
      version = "~> 1"
      source  = "github.com/hashicorp/virtualbox"
    }
  }
}

# Variables
variable "node_name" {
  type        = string
  description = "Name of the Syntropy node"
  default     = "syntropy-node-base"
}

variable "ubuntu_version" {
  type        = string
  description = "Ubuntu version to use"
  default     = "22.04"
}

variable "image_size" {
  type        = string
  description = "Size of the disk image"
  default     = "8G"
}

variable "memory" {
  type        = string
  description = "Amount of memory for the VM"
  default     = "2048"
}

variable "cpus" {
  type        = number
  description = "Number of CPUs for the VM"
  default     = 2
}

# Data sources
data "template_file" "user_data" {
  template = file("${path.root}/../cloud-init/user-data-base.yaml")
  vars = {
    node_name = var.node_name
  }
}

data "template_file" "meta_data" {
  template = file("${path.root}/../cloud-init/meta-data-base.yaml")
  vars = {
    node_name = var.node_name
  }
}

data "template_file" "network_config" {
  template = file("${path.root}/../cloud-init/network-config-base.yaml")
  vars = {
    node_name = var.node_name
  }
}

# QEMU Builder for Linux
source "qemu" "syntropy-node-linux" {
  # VM Configuration
  vm_name     = "${var.node_name}-linux.qcow2"
  memory      = var.memory
  cpus        = var.cpus
  disk_size   = var.image_size
  format      = "qcow2"
  
  # Boot configuration
  boot_wait   = "30s"
  boot_command = [
    "<esc><wait>",
    "linux /casper/vmlinuz quiet autoinstall ds=nocloud-net;s=http://{{ .HTTPIP }}:{{ .HTTPPort }}/",
    "<enter><wait>",
    "initrd /casper/initrd",
    "<enter><wait>",
    "boot",
    "<enter>"
  ]
  
  # ISO configuration
  iso_url      = "https://releases.ubuntu.com/${var.ubuntu_version}/ubuntu-${var.ubuntu_version}-live-server-amd64.iso"
  iso_checksum = "file:https://releases.ubuntu.com/${var.ubuntu_version}/SHA256SUMS"
  
  # Network configuration
  http_port_min = 8000
  http_port_max = 8000
  http_directory = "${path.root}/http"
  
  # Output configuration
  output_directory = "${path.root}/../output"
  shutdown_command = "sudo shutdown -P now"
  shutdown_timeout = "5m"
  
  # SSH configuration
  ssh_username = "admin"
  ssh_timeout  = "20m"
  
  # Headless mode
  headless = false
}

# VirtualBox Builder for cross-platform compatibility
source "virtualbox" "syntropy-node-vbox" {
  # VM Configuration
  vm_name           = "${var.node_name}-vbox"
  memory            = var.memory
  cpus              = var.cpus
  hard_drive_size   = tonumber(replace(var.image_size, "G", "")) * 1024
  
  # Boot configuration
  boot_wait         = "30s"
  boot_command = [
    "<esc><wait>",
    "linux /casper/vmlinuz quiet autoinstall ds=nocloud-net;s=http://{{ .HTTPIP }}:{{ .HTTPPort }}/",
    "<enter><wait>",
    "initrd /casper/initrd",
    "<enter><wait>",
    "boot",
    "<enter>"
  ]
  
  # ISO configuration
  iso_url      = "https://releases.ubuntu.com/${var.ubuntu_version}/ubuntu-${var.ubuntu_version}-live-server-amd64.iso"
  iso_checksum = "file:https://releases.ubuntu.com/${var.ubuntu_version}/SHA256SUMS"
  
  # Network configuration
  http_port_min = 8001
  http_port_max = 8001
  http_directory = "${path.root}/http"
  
  # Output configuration
  output_directory = "${path.root}/../output"
  format          = "ova"
  shutdown_command = "sudo shutdown -P now"
  shutdown_timeout = "5m"
  
  # SSH configuration
  ssh_username = "admin"
  ssh_timeout  = "20m"
  
  # VirtualBox specific
  vboxmanage = [
    ["modifyvm", "{{.Name}}", "--natpf1", "guestssh,tcp,,2222,,22"],
    ["modifyvm", "{{.Name}}", "--natpf1", "guesthttp,tcp,,8080,,80"]
  ]
}

# Build configurations
build {
  name = "syntropy-node-base"
  
  sources = [
    "source.qemu.syntropy-node-linux",
    "source.virtualbox.syntropy-node-vbox"
  ]
  
  # Provisioning steps
  provisioner "shell" {
    inline = [
      # Update system
      "sudo apt-get update -y",
      "sudo apt-get upgrade -y",
      
      # Install essential packages
      "sudo apt-get install -y curl wget git vim htop jq python3 python3-pip",
      "sudo apt-get install -y docker.io docker-compose",
      "sudo apt-get install -y fail2ban ufw prometheus-node-exporter",
      "sudo apt-get install -y openssh-server nmap ncdu tree tmux net-tools bc",
      "sudo apt-get install -y cloud-init cloud-guest-utils",
      
      # Configure Docker
      "sudo systemctl enable docker",
      "sudo systemctl start docker",
      "sudo usermod -aG docker admin",
      
      # Configure SSH
      "sudo systemctl enable ssh",
      "sudo systemctl start ssh",
      
      # Configure firewall
      "sudo ufw default deny incoming",
      "sudo ufw default allow outgoing",
      "sudo ufw allow ssh",
      "sudo ufw allow 9100/tcp",
      "sudo ufw --force enable",
      
      # Configure monitoring
      "sudo systemctl enable prometheus-node-exporter",
      "sudo systemctl start prometheus-node-exporter",
      
      # Create Syntropy directory structure
      "sudo mkdir -p /opt/syntropy/{identity/{owner,community},platform/{templates,scripts,data},metadata,logs,backups,config}",
      "sudo chown -R admin:admin /opt/syntropy/",
      
      # Create base configuration files
      "sudo tee /opt/syntropy/config/platform.yaml > /dev/null << 'EOF'",
      "platform:",
      "  name: syntropy_cooperative_grid",
      "  version: '2.0.0'",
      "  capabilities:",
      "    - container_orchestration",
      "    - resource_sharing",
      "    - cooperative_computing",
      "    - distributed_storage",
      "    - service_mesh",
      "    - universal_applications",
      "  services:",
      "    docker:",
      "      enabled: true",
      "    ssh:",
      "      enabled: true",
      "      port: 22",
      "    prometheus_exporter:",
      "      enabled: true",
      "      port: 9100",
      "    firewall:",
      "      enabled: true",
      "      policy: 'deny_incoming'",
      "EOF",
      
      # Set proper permissions
      "sudo chown admin:admin /opt/syntropy/config/platform.yaml",
      "sudo chmod 644 /opt/syntropy/config/platform.yaml",
      
      # Clean up
      "sudo apt-get autoremove -y",
      "sudo apt-get autoclean",
      "sudo rm -rf /tmp/*",
      "sudo rm -rf /var/tmp/*",
      
      # Create ready indicator
      "sudo touch /opt/syntropy/.ready",
      "sudo chown admin:admin /opt/syntropy/.ready",
      
      # Log completion
      "echo 'Syntropy base node image creation completed' | sudo logger -t syntropy"
    ]
  }
  
  # Post-processing
  post-processor "compress" {
    output = "${path.root}/../output/${var.node_name}-compressed.tar.gz"
  }
  
  post-processor "checksum" {
    checksum_types = ["sha256"]
    output = "${path.root}/../output/${var.node_name}.checksum"
  }
}
