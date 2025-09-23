# Syntropy Cooperative Grid - Node Infrastructure
# Terraform configuration for managing Syntropy nodes

terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0"
    }
    google = {
      source  = "hashicorp/google"
      version = "~> 4.0"
    }
    local = {
      source  = "hashicorp/local"
      version = "~> 2.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }
}

# Variables
variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
  default     = "dev"
}

variable "node_count" {
  description = "Number of nodes to create"
  type        = number
  default     = 1
}

variable "node_name_prefix" {
  description = "Prefix for node names"
  type        = string
  default     = "syntropy-node"
}

variable "instance_type" {
  description = "Instance type for cloud providers"
  type        = string
  default     = "t3.medium"
}

variable "ssh_public_key" {
  description = "SSH public key for node access"
  type        = string
  default     = ""
}

variable "coordinates" {
  description = "Geographic coordinates (lat,lon)"
  type        = string
  default     = "0,0"
}

variable "node_description" {
  description = "Description for the nodes"
  type        = string
  default     = "Syntropy Cooperative Grid Node"
}

# Local values
locals {
  common_tags = {
    Environment   = var.environment
    Project       = "SyntropyCooperativeGrid"
    ManagedBy     = "Terraform"
    Platform      = "syntropy"
    Version       = "2.0.0"
  }
  
  node_names = [for i in range(var.node_count) : "${var.node_name_prefix}-${var.environment}-${i + 1}"]
  
  # Extract coordinates
  coordinates_split = split(",", var.coordinates)
  latitude  = tonumber(local.coordinates_split[0])
  longitude = tonumber(local.coordinates_split[1])
}

# Data sources
data "template_file" "user_data" {
  count    = var.node_count
  template = file("${path.module}/../cloud-init/user-data-template.yaml")
  vars = {
    NodeName              = local.node_names[count.index]
    Coordinates           = var.coordinates
    CreatedAt             = timestamp()
    AdminPasswordHash     = "$6$rounds=4096$syntropy$N8mVzFK0Y1OelT1SKEjg0jIXzKMzL3ZcOGcE5xR8nS6E8qSO5qFV6eJs1g7T6E0cC7w.kfNO3FqC3YhE9Gz19."
    KeyInstallationCommands = "echo 'Key installation will be handled by Ansible'"
    MetadataCreationCommands = "echo 'Metadata creation will be handled by Ansible'"
    TemplateCreationCommands = "echo 'Template creation will be handled by Ansible'"
    StartupServiceCommands = "echo 'Startup services will be handled by Ansible'"
  }
}

data "template_file" "meta_data" {
  count    = var.node_count
  template = file("${path.module}/../cloud-init/meta-data-template.yaml")
  vars = {
    NodeID             = random_id.node_ids[count.index].hex
    NodeName           = local.node_names[count.index]
    OwnerPublicKey     = var.ssh_public_key
    CommunityPublicKey = var.ssh_public_key
  }
}

# Random IDs for nodes
resource "random_id" "node_ids" {
  count       = var.node_count
  byte_length = 8
}

# AWS Provider Configuration
provider "aws" {
  region = "us-west-2"
  
  default_tags {
    tags = local.common_tags
  }
}

# AWS Security Group
resource "aws_security_group" "syntropy_nodes" {
  count       = var.node_count
  name_prefix = "${local.node_names[count.index]}-"
  description = "Security group for Syntropy node ${local.node_names[count.index]}"
  
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "SSH access"
  }
  
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "HTTP access"
  }
  
  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "HTTPS access"
  }
  
  ingress {
    from_port   = 9100
    to_port     = 9100
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "Prometheus Node Exporter"
  }
  
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
    description = "All outbound traffic"
  }
  
  tags = merge(local.common_tags, {
    Name = "${local.node_names[count.index]}-sg"
  })
}

# AWS EC2 Instances
resource "aws_instance" "syntropy_nodes" {
  count                  = var.node_count
  ami                    = "ami-0c02fb55956c7d316" # Ubuntu 22.04 LTS
  instance_type          = var.instance_type
  key_name               = aws_key_pair.syntropy_key[0].key_name
  vpc_security_group_ids = [aws_security_group.syntropy_nodes[count.index].id]
  
  user_data = base64encode(data.template_file.user_data[count.index].rendered)
  
  root_block_device {
    volume_size = 20
    volume_type = "gp3"
    encrypted   = true
    
    tags = merge(local.common_tags, {
      Name = "${local.node_names[count.index]}-root"
    })
  }
  
  tags = merge(local.common_tags, {
    Name        = local.node_names[count.index]
    Description = var.node_description
    Coordinates = var.coordinates
    NodeID      = random_id.node_ids[count.index].hex
  })
}

# AWS Key Pair
resource "aws_key_pair" "syntropy_key" {
  count      = var.ssh_public_key != "" ? 1 : 0
  key_name   = "syntropy-${var.environment}-key"
  public_key = var.ssh_public_key
  
  tags = local.common_tags
}

# Local file outputs for Ansible inventory
resource "local_file" "ansible_inventory" {
  content = templatefile("${path.module}/ansible-inventory.tpl", {
    nodes = aws_instance.syntropy_nodes
  })
  filename = "${path.module}/../ansible/inventory/hosts.yml"
}

# Local file outputs for node metadata
resource "local_file" "node_metadata" {
  count = var.node_count
  
  content = jsonencode({
    node_id      = random_id.node_ids[count.index].hex
    node_name    = local.node_names[count.index]
    instance_id  = aws_instance.syntropy_nodes[count.index].id
    public_ip    = aws_instance.syntropy_nodes[count.index].public_ip
    private_ip   = aws_instance.syntropy_nodes[count.index].private_ip
    coordinates  = var.coordinates
    description  = var.node_description
    created_at   = timestamp()
    environment  = var.environment
  })
  
  filename = "${path.module}/../output/${local.node_names[count.index]}-metadata.json"
}

# Outputs
output "node_instances" {
  description = "Information about created Syntropy nodes"
  value = {
    for i, instance in aws_instance.syntropy_nodes : local.node_names[i] => {
      instance_id = instance.id
      public_ip   = instance.public_ip
      private_ip  = instance.private_ip
      node_id     = random_id.node_ids[i].hex
      coordinates = var.coordinates
    }
  }
}

output "ansible_inventory_file" {
  description = "Path to the generated Ansible inventory file"
  value       = local_file.ansible_inventory.filename
}

output "node_metadata_files" {
  description = "Paths to the generated node metadata files"
  value       = local_file.node_metadata[*].filename
}
