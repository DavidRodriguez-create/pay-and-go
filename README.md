# pay-and-go
Demo payment application to learn about Golang

## Architecture

This project follows **Clean Architecture** principles with clear separation of concerns across multiple layers:

- **Domain**: Core business entities and repository interfaces (no external dependencies)
- **Application**: Use cases, DTOs, mappers, and service orchestration
- **Infrastructure**: Repository implementations (currently in-memory, TODO: database)
- **Presentation**: REST API controllers, presenters, and routing

## Services

### Account Service ✅
Fully implemented account management microservice with CRUD operations:
- **Port**: 8081 (HTTP)
- **Status**: Deployed and tested
- **Test Coverage**: 100% (domain & application), 97.7% (infrastructure)
- **Endpoints**:
  - `POST /accounts` - Create account
  - `GET /accounts` - List all accounts
  - `GET /accounts?id={id}` - Get account by ID
  - `GET /accounts/by-number?account_number={number}` - Get account by number
  - `PUT /accounts?id={id}` - Update account
  - `DELETE /accounts?id={id}` - Delete account (soft delete)
  - `GET /health` - Health check

**Example Usage**:
```bash
# Create account
curl -X POST http://localhost:8081/accounts \
  -H "Content-Type: application/json" \
  -d '{"beholder_name":"John Doe","country_code":"US"}'

# List accounts
curl http://localhost:8081/accounts

# Health check
curl http://localhost:8081/health
```

### Card Service 🚧
Prepared for deployment (minimal implementation):
- **Port**: 8082 (HTTP)
- **Status**: In development

## Deployment

### Prerequisites
- Podman or Docker installed
- Go 1.23+ (for local development)

### Option 1: Run with Podman (Recommended)

#### Build Container Image
```bash
./build-images.sh
```

#### Run Container
```bash
# Stop and remove existing container (if any)
podman rm -f account-service

# Start account service
podman run -d --name account-service -p 8081:8081 -e PORT=8081 localhost/account-service:latest

# View logs
podman logs -f account-service

# Stop service
podman stop account-service
```

### Option 2: Run Locally (Development)
```bash
cd services/account
go run cmd/main.go
```

### Option 3: Deploy to Kubernetes
```bash
# Deploy services
./deploy.sh

# Access via NodePort
# Account Service: http://localhost:30081

# Undeploy
./undeploy.sh
```

**Note**: Kubernetes deployment requires `kubectl` and a running cluster.

## Testing

### Run All Tests
```bash
cd services/account
go test ./tests/... -v
```

### Test Organization
```
tests/
├── integration/              # End-to-end API tests
│   └── integration_test.go
└── unit/                     # Unit tests by layer
    ├── application/          # Use case tests (100% coverage)
    ├── domain/              # Entity tests (100% coverage)
    └── infrastructure/      # Repository tests (97.7% coverage)
```

### Coverage Report
```bash
cd services/account
go test ./tests/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

See [TEST_SUMMARY.md](services/account/TEST_SUMMARY.md) and [tests/README.md](services/account/tests/README.md) for detailed test documentation.

## Development

### Project Structure
```
pay-and-go/
├── services/
│   └── account/
│       ├── cmd/                    # Application entry point
│       ├── domain/                 # Entities and interfaces
│       ├── application/            # Use cases, DTOs, mappers
│       ├── infrastructure/         # Repository implementations
│       ├── presentation/           # Controllers, presenters, routes
│       │   ├── controllers/        # HTTP request handlers
│       │   ├── presenters/         # Response formatters
│       │   └── routes/            # Route configuration
│       ├── tests/                 # Test suite (unit + integration)
│       │   ├── unit/
│       │   │   ├── application/
│       │   │   ├── domain/
│       │   │   └── infrastructure/
│       │   └── integration/
│       ├── TEST_SUMMARY.md        # Test documentation
│       └── go.mod
├── k8s/                           # Kubernetes manifests
│   ├── account-service.yaml
│   └── card-service.yaml
├── Dockerfile.account
├── Dockerfile.card
├── build-images.sh
├── deploy.sh
└── undeploy.sh
```

### Clean Architecture Guidelines
- **Domain entities** MUST NOT have JSON tags (purely internal)
- **DTOs** handle JSON serialization in the Application layer
- **Each use case** has its own controller
- **Repository interfaces** in Domain, implementations in Infrastructure
- **Tests** are organized by layer with 100% coverage on core logic

### Adding New Features
1. Define entity in `domain/`
2. Create use case in `application/`
3. Implement repository in `infrastructure/`
4. Add controller in `presentation/controllers/`
5. Update routes in `presentation/routes/`
6. Write tests in `tests/unit/` and `tests/integration/`

## Technologies

- **Language**: Go 1.23
- **Containerization**: Podman/Docker
- **Orchestration**: Kubernetes (optional)
- **Testing**: Go testing framework with table-driven tests
- **Architecture**: Clean Architecture with DDD principles
