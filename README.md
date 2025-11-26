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
Fully implemented account management microservice with CRUD operations and event publishing:
- **Port**: 8081 (HTTP)
- **Status**: Deployed and tested
- **Test Coverage**: 100% (domain & application), 97.7% (infrastructure)
- **Event Publishing**: Publishes `account.created` and `account.status_changed` events to Kafka
- **Endpoints**:
  - `POST /account` - Create account (publishes event)
  - `GET /accounts` - List all accounts
  - `GET /account?id={id}` - Get account by ID
  - `GET /accounts/by-number?account_number={number}` - Get account by number
  - `PUT /account?id={id}` - Update account (publishes event on status change)
  - `DELETE /account?id={id}` - Delete account (publishes event)
  - `GET /health` - Health check

**Example Usage**:
```bash
# Create account
curl -X POST http://localhost:8081/account \
  -H "Content-Type: application/json" \
  -d '{"beholder_name":"John Doe","country_code":"US"}'

# List accounts
curl http://localhost:8081/accounts

# Health check
curl http://localhost:8081/health
```

### Card Service ✅
Fully implemented card management microservice with event-driven account synchronization:
- **Port**: 8082 (HTTP)
- **Status**: Deployed and tested
- **Test Coverage**: 109 tests passing (domain, application, infrastructure, integration)
- **Event Consumption**: Consumes `account.created` and `account.status_changed` events from Kafka
- **Endpoints**:
  - `POST /card` - Create card (requires account synced via Kafka)
  - `GET /cards` - List all cards
  - `GET /card?id={id}` - Get card by ID
  - `GET /cards/by-number?card_number={number}` - Get card by card number
  - `GET /cards/by-account?account_id={id}` - Get cards by account ID
  - `DELETE /card?id={id}` - Delete card (soft delete)
  - `GET /health` - Health check

**Example Usage**:
```bash
# Create card (account must exist and be synced via Kafka)
curl -X POST http://localhost:8082/card \
  -H "Content-Type: application/json" \
  -d '{"account_id":"<ACCOUNT_ID>","card_type":"DEBIT","country":"US"}'

# List all cards
curl http://localhost:8082/cards

# Health check
curl http://localhost:8082/health
```

## Deployment

### Prerequisites
- Podman installed
- Go 1.23+ (for local development)

### Quick Start (Recommended)

Deploy all services with a single command:

```bash
./deploy-and-test.sh
```

This will:
- Clean up any existing containers
- Build service images
- Start Zookeeper and Kafka
- Start Account and Card services
- Display sample test cases and verification commands

**Services Available**:
- Account Service: http://localhost:8081
- Card Service: http://localhost:8082
- Kafka Broker: localhost:9092
- Zookeeper: localhost:2181

### Manual Deployment

If you prefer manual control:

```bash
# Build images
./build-images.sh

# Clean up existing containers
podman rm -f account-service card-service kafka zookeeper

# Follow the manual steps in deploy-and-test.sh
```

### Stop Services

```bash
podman rm -f account-service card-service kafka zookeeper
```

### Deploy to Kubernetes
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

### Unit and Integration Tests

```bash
# Account Service
cd services/account
go test ./tests/... -v

# Card Service  
cd services/card
go test ./tests/... -v
```

### Test Coverage
```bash
cd services/account
go test ./tests/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

See service-specific test documentation:
- [Account Service Tests](services/account/tests/README.md)
- [Card Service Tests](services/card/tests/README.md)

### API Testing

After deploying with `./deploy-and-test.sh`, test the API endpoints:

#### 1. Health Checks
```bash
curl http://localhost:8081/health
curl http://localhost:8082/health
```

#### 2. Create Account (triggers Kafka event)
```bash
curl -X POST http://localhost:8081/account \
  -H "Content-Type: application/json" \
  -d '{"beholder_name":"John Doe","country_code":"US"}'
```

#### 3. List Accounts
```bash
curl http://localhost:8081/accounts
```

#### 4. Update Account Status (triggers event on status change)
```bash
curl -X PUT 'http://localhost:8081/account?id=<ACCOUNT_ID>' \
  -H "Content-Type: application/json" \
  -d '{"status":"SUSPENDED"}'
```

#### 5. Create Card (requires account synced via Kafka)
```bash
curl -X POST http://localhost:8082/card \
  -H "Content-Type: application/json" \
  -d '{"account_id":"<ACCOUNT_ID>","card_type":"DEBIT","country":"US"}'
```

#### 6. List Cards
```bash
curl http://localhost:8082/cards
```

### Verify Event-Driven Flow

```bash
# Check Kafka events in card service logs
podman logs card-service | grep "Received account event"

# Check Kafka connection in account service
podman logs account-service | grep -i kafka

# Monitor real-time logs
podman logs -f card-service
```

## Development

### Project Structure
```
pay-and-go/
├── services/
│   ├── account/                   # Account management service
│   │   ├── cmd/                   # Application entry point
│   │   ├── domain/                # Entities and interfaces
│   │   ├── application/           # Use cases, DTOs, mappers
│   │   ├── infrastructure/        # Repository & Kafka implementations
│   │   ├── presentation/          # Controllers, presenters, routes
│   │   ├── tests/                 # Test suite (unit + integration)
│   │   └── go.mod
│   └── card/                      # Card management service
│       ├── cmd/
│       ├── domain/
│       ├── application/
│       ├── infrastructure/
│       ├── presentation/
│       ├── tests/
│       └── go.mod
├── k8s/                           # Kubernetes manifests
│   ├── all-services.yaml          # Complete deployment
│   ├── kafka.yaml                 # Kafka & Zookeeper
│   ├── account-service.yaml       # Account service
│   └── card-service.yaml          # Card service
├── build-images.sh                # Build container images
├── deploy-and-test.sh             # Deploy all services (Podman)
├── deploy.sh                      # Deploy to Kubernetes
└── undeploy.sh                    # Remove Kubernetes deployment
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
- **Message Broker**: Apache Kafka (for event-driven communication)
- **Testing**: Go testing framework with table-driven tests
- **Architecture**: Clean Architecture with DDD principles and Event-Driven Architecture

## Event-Driven Communication

The services use Kafka for asynchronous event-driven communication:

- **Account Service** publishes events when accounts are created or their status changes
- **Card Service** consumes these events to maintain a local cache of account states
- **Benefits**: Loose coupling, eventual consistency, improved resilience

For detailed integration guide, see [INTEGRATION.md](INTEGRATION.md).
