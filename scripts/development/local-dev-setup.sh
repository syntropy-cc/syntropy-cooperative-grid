#!/bin/bash

# Local Development Environment Setup

set -e

echo "Setting up Syntropy Cooperative Grid development environment..."

# Check prerequisites
echo "Checking prerequisites..."
command -v docker >/dev/null 2>&1 || { echo "Docker is required but not installed."; exit 1; }
command -v docker-compose >/dev/null 2>&1 || { echo "Docker Compose is required but not installed."; exit 1; }

# Create necessary directories
mkdir -p tools/monitoring
mkdir -p configs/defaults

# Create basic Prometheus config
cat > tools/monitoring/prometheus.yml << 'PROMETHEUS_EOF'
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']
  
  - job_name: 'syntropy-services'
    static_configs:
      - targets: ['localhost:8080']  # Placeholder for future services
PROMETHEUS_EOF

# Create basic PostgreSQL init script
cat > configs/defaults/postgres-init.sql << 'SQL_EOF'
-- Syntropy Cooperative Grid Database Initialization

-- Create schemas
CREATE SCHEMA IF NOT EXISTS cooperative;
CREATE SCHEMA IF NOT EXISTS monitoring;
CREATE SCHEMA IF NOT EXISTS security;

-- Create basic tables (placeholder)
CREATE TABLE IF NOT EXISTS cooperative.nodes (
    id SERIAL PRIMARY KEY,
    node_id VARCHAR(64) UNIQUE NOT NULL,
    node_type VARCHAR(32) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE cooperative.nodes IS 'Registry of nodes in the cooperative grid';
SQL_EOF

echo "Development environment setup complete!"
echo "Run 'make dev-start' to start local services"
