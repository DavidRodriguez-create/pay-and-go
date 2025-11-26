# Pay & Go

Event-driven payment microservices platform built with Go, demonstrating Clean Architecture and real-time event streaming with Kafka.

## 🚀 Quick Start

**Prerequisites**: 
- [Podman](https://podman.io/getting-started/installation)
- [Go 1.23+](https://go.dev/dl/) (optional, for local development)

### Start Everything

```bash
./manage-services.sh start
```

Then open **`ui.html`** in your browser to interact with the services.

### Stop Everything

```bash
./manage-services.sh stop
```

### Other Commands

```bash
./manage-services.sh status   # Check service health
./manage-services.sh restart  # Restart all services
```

**That's it!** The UI connects to:
- 💼 Account Service: http://localhost:8081
- 💳 Card Service: http://localhost:8082

## 🎮 Using the UI

1. **Open `ui.html`** in any web browser (just double-click the file)
2. **Create an Account** - Enter a name and country, click "Create Account"
3. **Create a Card** - Use the account ID (auto-filled) to create a card
4. **Test Operations** - List, suspend, delete accounts and cards
5. **Watch Kafka Events** - Account changes automatically sync to card service

The UI automatically checks service health and shows real-time status indicators.

## 🏗️ Architecture

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

## 🐳 Deployment

### Prerequisites
- **[Podman](https://podman.io/getting-started/installation)** - Container runtime
- **[Go 1.23+](https://go.dev/dl/)** - For local development (optional)

### Service Management (Recommended)

Use the unified management script for all operations:

```bash
# Start all services
./manage-services.sh start

# Stop all services
./manage-services.sh stop

# Restart all services
./manage-services.sh restart

# Check service status
./manage-services.sh status
```

**What `start` does**:
1. ✅ Cleans up existing containers
2. ✅ Builds account-service and card-service images
3. ✅ Starts Zookeeper and Kafka
4. ✅ Deploys both microservices
5. ✅ Shows service URLs and next steps

**Services Available**:
- 🌐 **UI**: Open `ui.html` in your browser
- 💼 Account Service: http://localhost:8081
- 💳 Card Service: http://localhost:8082
- 📨 Kafka Broker: localhost:9092
- 🔧 Zookeeper: localhost:2181

### Kubernetes Deployment (Optional)

```bash
cd k8s

# Deploy to cluster
./deploy.sh

# Access via NodePort
# Account Service: http://localhost:30081

# Remove deployment
./undeploy.sh
```

**Note**: Requires `kubectl` and a running Kubernetes cluster.

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

**Option 1: Use the Web UI** (Recommended)
- Open `ui.html` in your browser
- Visual interface with auto-fill and real-time responses

**Option 2: Command Line with curl**

After deploying, test the API endpoints:

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
├── ui.html                        # 🎮 Web UI for testing services
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
├── docker-compose.yml             # Service orchestration
├── podman/                        # Container build files
│   ├── Dockerfile.account         # Account service image
│   └── Dockerfile.card            # Card service image
├── k8s/                           # Kubernetes manifests
│   ├── all-services.yaml          # Complete deployment
│   ├── kafka.yaml                 # Kafka & Zookeeper
│   ├── account-service.yaml       # Account service
│   ├── card-service.yaml          # Card service
│   ├── deploy.sh                  # Deploy to Kubernetes
│   └── undeploy.sh                # Remove Kubernetes deployment
└── manage-services.sh             # Main deployment script (Podman)
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
