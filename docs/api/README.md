# Syntropy Cooperative Grid - API Documentation

This directory contains comprehensive API documentation for all services in the Syntropy Cooperative Grid.

## API Categories

### [Cooperative Services](cooperative-services/)
- **Credit System API** - Manage credits and transactions
- **Node Discovery API** - Node registration and discovery
- **Resource Broker API** - Resource matching and allocation
- **Reputation System API** - Node reputation and trust scoring

### [Resource Management](resource-management/)
- **Resource Allocation API** - CPU, memory, storage allocation
- **Workload Management API** - Container and service lifecycle
- **Performance Monitoring API** - Resource usage and performance metrics
- **SLA Management API** - Service level agreement tracking

### [Blockchain](blockchain/)
- **Consensus API** - Blockchain consensus and validation
- **Smart Contracts API** - Contract deployment and interaction
- **Wallet API** - Token and transaction management
- **Governance API** - Voting and proposal management

## API Design Principles

1. **RESTful Design** - Standard HTTP methods and status codes
2. **OpenAPI Specification** - All APIs documented with OpenAPI 3.0
3. **Consistent Versioning** - Semantic versioning for all APIs
4. **Security First** - Authentication and authorization built-in
5. **Developer Friendly** - Comprehensive examples and SDKs

## Authentication

All APIs use a combination of:
- **JWT Tokens** for user authentication
- **API Keys** for service-to-service communication
- **mTLS** for secure service mesh communication

## Rate Limiting

- **User APIs**: 1000 requests per hour per user
- **Service APIs**: 10,000 requests per hour per service
- **Public APIs**: 100 requests per hour per IP

## Getting Started

1. Register for an API key
2. Review the OpenAPI specifications
3. Try the interactive API documentation
4. Use the SDKs for your preferred language
