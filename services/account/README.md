# Account Service

RESTful microservice for managing bank accounts with event-driven architecture support.

## Features

- ✅ **Account Management**: Create, read, update, and delete accounts
- ✅ **Status Management**: Track account status (ACTIVE, BLOCKED, DELETED)
- ✅ **Event Publishing**: Publishes events to Kafka for account lifecycle changes
- ✅ **Clean Architecture**: Domain-driven design with clear separation of concerns
- ✅ **In-Memory Storage**: Fast development with in-memory repository

## Event-Driven Architecture

The account service publishes events to Kafka that other services can consume:

### Published Events

#### 1. Account Created Event
Published when a new account is created.

```json
{
  "type": "account.created",
  "account_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "ACTIVE"
}
```

#### 2. Account Status Changed Event
Published when an account's status is updated (e.g., blocked, deleted).

```json
{
  "type": "account.status_changed",
  "account_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "BLOCKED"
}
```

### Event Publishing Configuration

Events are published to Kafka when configured via environment variables:

```bash
# Enable event publishing
export KAFKA_BROKERS=localhost:9092
export KAFKA_TOPIC=account-events
```

If Kafka is not configured, the service runs normally but without event publishing (graceful degradation).

## API Endpoints

### Create Account
```bash
POST /account
Content-Type: application/json

{
  "beholder_name": "John Doe",
  "country_code": "US"
}
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "account_number": "ACC-12345678",
  "beholder_name": "John Doe",
  "country_code": "US",
  "status": "ACTIVE",
  "created_at": "2024-01-01T12:00:00Z"
}
```

### Get Account
```bash
GET /account?id=550e8400-e29b-41d4-a716-446655440000
```

### List All Accounts
```bash
GET /accounts
```

### Update Account
```bash
PUT /account?id=550e8400-e29b-41d4-a716-446655440000
Content-Type: application/json

{
  "beholder_name": "Jane Doe",
  "status": "BLOCKED"
}
```

**Triggers Event:** `account.status_changed` (if status was modified)

### Delete Account (Soft Delete)
```bash
DELETE /account?id=550e8400-e29b-41d4-a716-446655440000
```

**Triggers Event:** `account.status_changed` with status "DELETED"

### Health Check
```bash
GET /health
```

## Running the Service

### With Docker/Podman (Recommended)

```bash
# Build image
podman build -f Dockerfile.account -t account-service .

# Run with Kafka
podman run -p 8081:8081 \
  -e KAFKA_BROKERS=localhost:9092 \
  -e KAFKA_TOPIC=account-events \
  account-service

# Run without Kafka (standalone)
podman run -p 8081:8081 account-service
```

### Local Development

```bash
cd services/account

# Install dependencies
go mod download

# Run with .env file
go run cmd/main.go
```

### Using .env File

Create a `.env` file in the service root:

```env
# Server Configuration
PORT=8081

# Kafka Configuration (optional)
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC=account-events
```

## Configuration

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PORT` | HTTP server port | `8081` | No |
| `KAFKA_BROKERS` | Comma-separated Kafka broker addresses | - | No |
| `KAFKA_TOPIC` | Kafka topic for account events | - | No |

## Architecture

```
services/account/
├── cmd/
│   └── main.go                    # Application entry point
├── domain/
│   ├── account.go                 # Account entity
│   ├── account_repository.go     # Repository interface
│   └── event_publisher.go        # Event publisher interface
├── application/
│   ├── create_account.go         # Create account use case
│   ├── update_account.go         # Update account use case (publishes events)
│   ├── delete_account.go         # Delete account use case
│   ├── view_account.go           # View account use cases
│   └── service.go                # Service orchestration
├── infrastructure/
│   ├── memory_account_repository.go  # In-memory repository
│   └── kafka_producer.go             # Kafka event publisher
└── presentation/
    ├── controllers/              # HTTP handlers
    └── routes/                   # Route configuration
```

## Event Publishing Flow

1. **Account Creation**:
   - Account is created in repository
   - `account.created` event is published to Kafka
   - Event contains account ID and initial status (ACTIVE)

2. **Status Update**:
   - Account status is updated in repository
   - If status changed, `account.status_changed` event is published
   - Event contains account ID and new status

3. **Graceful Degradation**:
   - If Kafka is unavailable, events are logged but operations continue
   - Service remains functional even without event publishing
   - Errors are logged for monitoring

## Integration with Card Service

The card service consumes events from the account service to maintain a local cache of account states. This enables:

- **Fast validation**: No cross-service calls needed for card creation
- **Eventual consistency**: Account cache updated asynchronously
- **Decoupling**: Services communicate via events, not direct API calls

See [INTEGRATION.md](../../INTEGRATION.md) for complete integration guide.

## Testing

```bash
# Run all tests
go test ./... -v

# Run with coverage
go test ./... -cover

# Run specific test
go test ./tests/integration/... -v
```

## Development

### Adding New Events

1. Define event structure in `infrastructure/kafka_producer.go`
2. Add publish method to `EventPublisher` interface
3. Call publish method in appropriate use case
4. Update consuming services to handle new event type

### Best Practices

- Events are published **after** database operations succeed
- Event publishing failures don't rollback database operations
- Always log event publishing errors for monitoring
- Keep events small and focused (only essential data)

## Deployment

### Kubernetes

See `k8s/account-service.yaml` for Kubernetes deployment configuration.

### Environment Variables for Production

```yaml
env:
  - name: PORT
    value: "8081"
  - name: KAFKA_BROKERS
    value: "kafka-broker-1:9092,kafka-broker-2:9092,kafka-broker-3:9092"
  - name: KAFKA_TOPIC
    value: "account-events"
```

## Monitoring

Monitor these key metrics:
- Event publishing success/failure rate
- Event publishing latency
- Account creation/update rates
- Failed event publishing attempts (check logs)

## Troubleshooting

### Events Not Being Published

1. Check Kafka broker connectivity:
   ```bash
   telnet localhost 9092
   ```

2. Verify environment variables are set:
   ```bash
   echo $KAFKA_BROKERS
   echo $KAFKA_TOPIC
   ```

3. Check service logs for Kafka errors:
   ```bash
   # Look for "Failed to publish event" messages
   ```

### Service Starts But Events Missing

- Ensure the Kafka topic exists
- Verify broker addresses are correct
- Check network connectivity to Kafka cluster

## License

MIT
