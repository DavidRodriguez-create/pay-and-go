# pay-and-go
Demo payment application to learn about Golang

## Architecture

This project follows **Clean Architecture** principles with clear separation of concerns across multiple layers:

- **Domain**: Core business entities and repository interfaces (no external dependencies)
- **Application**: Use cases, DTOs, mappers, and service orchestration
- **Infrastructure**: Repository implementations (currently in-memory, TODO: database)
- **Presentation**: REST API controllers, presenters, and routing

## Services

### Account Service
Fully implemented account management microservice with CRUD operations:
- **Port**: 8081 (NodePort: 30081)
- **Endpoints**:
  - `POST /accounts` - Create account
  - `GET /accounts` - List all accounts
  - `GET /accounts?id={id}` - Get account by ID
  - `GET /accounts/by-number?account_number={number}` - Get account by number
  - `PUT /accounts?id={id}` - Update account
  - `DELETE /accounts?id={id}` - Delete account (soft delete)
  - `GET /health` - Health check

### Card Service
Prepared for deployment (minimal implementation):
- **Port**: 8082 (NodePort: 30082)

## Deployment with Kubernetes and Podman

### Prerequisites
- Podman installed
- kubectl installed
- Kubernetes cluster (or Podman Desktop with Kubernetes enabled)

### Build Container Images
```bash
./build-images.sh
```

This builds both `account-service:latest` and `card-service:latest` images.

### Deploy to Kubernetes
```bash
./deploy.sh
```

This deploys the account service to Kubernetes. The card service deployment is prepared but commented out until full implementation.

### Access Services
Once deployed:
- Account Service: `http://localhost:30081`
- Health check: `http://localhost:30081/health`

### Undeploy Services
```bash
./undeploy.sh
```

This removes all deployed services from the Kubernetes cluster.

## Development

### Project Structure
```
pay-and-go/
├── services/
│   └── account/
|       ├── cmd/              # 
│       ├── domain/           # Entities and interfaces
│       ├── application/      # Use cases, DTOs, mappers
│       ├── infrastructure/   # Repository implementations
│       └── presentation/     # Controllers, presenters, routes
├── k8s/                      # Kubernetes manifests
│   ├── account-service.yaml
│   └── card-service.yaml
├── Dockerfile.account
├── Dockerfile.card
├── build-images.sh
├── deploy.sh
└── undeploy.sh
```

### Clean Architecture Guidelines
- Domain entities MUST NOT have JSON tags (purely internal)
- DTOs handle JSON serialization in the Application layer
- Each use case has its own controller
- Repository interfaces in Domain, implementations in Infrastructure
