#!/bin/bash

# Syntropy Cooperative Grid - Exemplo de Deployment de Nó
# Este script demonstra como usar a infraestrutura como código para deploy de nós

set -e  # Parar em caso de erro

# Configurações
NODE_NAME="${NODE_NAME:-syntropy-node-01}"
COORDINATES="${COORDINATES:--23.5505,-46.6333}"
NODE_DESCRIPTION="${NODE_DESCRIPTION:-Syntropy Cooperative Grid Node}"
ENVIRONMENT="${ENVIRONMENT:-dev}"
NODE_COUNT="${NODE_COUNT:-1}"
INSTANCE_TYPE="${INSTANCE_TYPE:-t3.medium}"

echo "🚀 Iniciando deployment de nó Syntropy usando Infrastructure as Code"
echo "   Nome: $NODE_NAME"
echo "   Coordenadas: $COORDINATES"
echo "   Ambiente: $ENVIRONMENT"
echo "   Quantidade: $NODE_COUNT"
echo ""

# Função para verificar pré-requisitos
check_prerequisites() {
    echo "🔍 Verificando pré-requisitos..."
    
    # Verificar se as ferramentas estão instaladas
    for tool in terraform ansible; do
        if ! command -v $tool &> /dev/null; then
            echo "❌ $tool não está instalado. Instale antes de continuar."
            exit 1
        fi
    done
    
    echo "✅ Pré-requisitos verificados"
}

# Função para provisionar infraestrutura
provision_infrastructure() {
    echo "🌍 Provisionando infraestrutura com Terraform..."
    cd terraform/syntropy-nodes
    
    # Inicializar Terraform
    terraform init
    
    # Aplicar configuração
    terraform apply -auto-approve \
        -var="environment=$ENVIRONMENT" \
        -var="node_count=$NODE_COUNT" \
        -var="node_name_prefix=$NODE_NAME" \
        -var="instance_type=$INSTANCE_TYPE" \
        -var="coordinates=$COORDINATES" \
        -var="node_description=$NODE_DESCRIPTION"
    
    echo "✅ Infraestrutura provisionada"
    cd ../..
}

# Função para configurar nós
configure_nodes() {
    echo "🔧 Configurando nós com Ansible..."
    cd ansible
    
    # Aguardar um pouco para os nós ficarem prontos
    echo "⏳ Aguardando nós ficarem prontos..."
    sleep 30
    
    # Executar playbook
    ansible-playbook \
        -i ../terraform/syntropy-nodes/inventory/hosts.yml \
        -e "environment=$ENVIRONMENT" \
        -e "node_description=$NODE_DESCRIPTION" \
        playbooks/syntropy-node-setup.yml
    
    echo "✅ Nós configurados"
    cd ..
}

# Função para verificar deployment
verify_deployment() {
    echo "🔍 Verificando deployment..."
    
    # Verificar status dos nós via Ansible
    cd ansible
    ansible all -i ../terraform/syntropy-nodes/inventory/hosts.yml -m ping
    
    echo "✅ Deployment verificado"
    cd ..
}

# Função principal
main() {
    echo "🎯 Syntropy Cooperative Grid - Deployment de Nós"
    echo "================================================"
    echo ""
    
    check_prerequisites
    provision_infrastructure
    configure_nodes
    verify_deployment
    
    echo ""
    echo "🎉 Deployment concluído com sucesso!"
}

# Executar função principal
main "$@"