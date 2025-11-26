# Card Service

Event-driven microservice for payment card management with Kafka-based account synchronization.

## Architecture

Follows **Clean Architecture** with event-driven integration:

- **Domain**: Core entities (Card, AccountCache) and repository interfaces
- **Application**: Use cases, DTOs, and business logic
- **Infrastructure**: In-memory repositories and Kafka event consumer
- **Presentation**: REST API controllers, presenters, and routes

## Event-Driven Design

The card service maintains a **local account cache** synchronized via Kafka events from the account service:

### Account Event Flow

```
Account Service                    Card Service
     |                                  |
     |-- Creates Account -------------> |
     |   (publishes event)              |
     |                                  |-- Receives "account.created"
     |                                  |-- Upserts AccountCache
     |                                  |   (id, status=ACTIVE)
     |                                  |
     |-- Updates Account Status ------> |
     |   (publishes event)              |
     |                                  |-- Receives "account.status_changed"
     |                                  |-- Updates AccountCache
     |                                  |   (status=DELETED/BLOCKED)
```

### Why This Pattern?

1. **High Performance**: Card creation doesn't make HTTP calls to account service
2. **Resilience**: Service works even if account service is temporarily down
3. **Scalability**: Supports high card creation volumes
4. **Eventual Consistency**: Account status synced asynchronously via events

## Domain Model

### Card Entity

```go
Card {
    ID                string    // UUID
    CardNumber        string    // Generated: COUNTRY-UUID
    Country           string    // Country code
    AccountID         string    // Reference to account
    Deleted           bool      // Soft delete flag
    CreationTimestamp time.Time // Creation time
}
```

### AccountCache Entity

```go
AccountCache {
    ID     string        // Account ID
    Status AccountStatus // ACTIVE, BLOCKED, DELETED
}
```

## API Endpoints

### Card Operations

| Method | Endpoint | Description | Body |
|--------|----------|-------------|------|
| POST | `/card` | Create card | `{"country": "US", "account_id": "xxx"}` |
| GET | `/card?id=xxx` | Get card by ID | - |
| DELETE | `/card?id=xxx` | Delete card (soft) | - |
| GET | `/cards` | List all cards | - |
| GET | `/cards/by-number?card_number=xxx` | Get by card number | - |
| GET | `/cards/by-account?account_id=xxx` | Get by account ID | - |
| GET | `/health` | Health check | - |

### Business Rules

- ✅ Cards can only be created for **ACTIVE** accounts
- ❌ Cannot create cards for **DELETED** accounts
- ❌ Cannot create cards for **BLOCKED** accounts
- ❌ Cannot create cards for **non-existent** accounts
- ✅ Card deletion is **soft delete** (sets `deleted=true`)

## Kafka Integration

### Configuration

#### Using .env File (Recommended)

Create a `.env` file in the `services/card/` directory:

```env
# Server Configuration
PORT=8082

# Kafka Configuration
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC=account-events
KAFKA_GROUP_ID=card-service
```

The service automatically loads `.env` on startup. Use `.env.example` as a template.

#### Using Environment Variables

Alternatively, set environment variables directly:
- `PORT`: HTTP server port (default: `8082`)
- `KAFKA_BROKERS`: Comma-separated broker list (default: `localhost:9092`)
- `KAFKA_TOPIC`: Topic to consume (default: `account-events`)
- `KAFKA_GROUP_ID`: Consumer group ID (default: `card-service`)

Environment variables override `.env` file values.

### Event Schema

**Account Created Event:**
```json
{
  "type": "account.created",
  "account_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "ACTIVE"
}
```

**Account Status Changed Event:**
```json
{
  "type": "account.status_changed",
  "account_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "DELETED"
}
```

### Consumer Behavior

- **Consumer Group**: Enables horizontal scaling
- **Auto-commit**: Enabled for simplicity
- **Error Handling**: Logs errors but continues consuming
- **Graceful Shutdown**: Properly closes consumer on service termination

## Running the Service

### Prerequisites
- Go 1.23+
- Kafka cluster (optional, service starts without it)

### Setup

1. **Copy the example configuration:**
```bash
cd services/card
cp .env.example .env
```

2. **Edit `.env` with your settings** (optional, defaults work for local development)

### Local Development

```bash
cd services/card
go run cmd/main.go
```

The service automatically loads configuration from `.env` file.

### With Custom Configuration

You can override `.env` values with environment variables:

```bash
# Environment variables take precedence over .env
export PORT=9999
export KAFKA_BROKERS=kafka1:9092,kafka2:9092
go run cmd/main.go
```

### With Docker/Podman

```bash
podman build -f ../../Dockerfile.card -t card-service:latest .
podman run -d --name card-service \
  -p 8082:8082 \
  -e PORT=8082 \
  -e KAFKA_BROKERS=kafka:9092 \
  card-service:latest
```

## Testing

### Example Card Creation

```bash
# Create a card (requires active account in cache)
curl -X POST http://localhost:8082/card \
  -H "Content-Type: application/json" \
  -d '{
    "country": "US",
    "account_id": "550e8400-e29b-41d4-a716-446655440000"
  }'

# Response:
{
  "id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
  "card_number": "US-7c9e6679",
  "country": "US",
  "account_id": "550e8400-e29b-41d4-a716-446655440000",
  "deleted": false,
  "creation_timestamp": "2025-11-25T10:30:00Z"
}
```

### Example Card Retrieval

```bash
# Get card by ID
curl http://localhost:8082/card?id=7c9e6679-7425-40de-944b-e07fc1f90ae7

# Get cards by account
curl http://localhost:8082/cards/by-account?account_id=550e8400-e29b-41d4-a716-446655440000

# List all cards
curl http://localhost:8082/cards
```

### Example Card Deletion

```bash
# Soft delete a card
curl -X DELETE http://localhost:8082/card?id=7c9e6679-7425-40de-944b-e07fc1f90ae7

# Response:
{
  "message": "Card deleted successfully"
}
```

## Error Responses

| Error | Status Code | Scenario |
|-------|-------------|----------|
| `account not found` | 404 | Account not in cache |
| `cannot create card for deleted account` | 403 | Account is deleted |
| `cannot create card for inactive account` | 403 | Account is blocked |
| `card not found` | 404 | Card doesn't exist |
| `card is already deleted` | 409 | Attempting to delete twice |

## Future Enhancements

- [ ] **Database Integration**: Replace in-memory repos with PostgreSQL
- [ ] **Redis Cache**: Add Redis for distributed account cache
- [ ] **Event Publishing**: Publish card events to Kafka
- [ ] **Metrics**: Add Prometheus metrics
- [ ] **Tracing**: Add distributed tracing (Jaeger/Zipkin)
- [ ] **Circuit Breaker**: Handle Kafka unavailability gracefully
- [ ] **Admin API**: Endpoints to view account cache status
- [ ] **Batch Operations**: Bulk card creation/deletion

## Technologies

- **Language**: Go 1.23
- **Event Streaming**: Kafka (segmentio/kafka-go)
- **Architecture**: Clean Architecture + Event-Driven
- **Patterns**: Repository, Use Case, Dependency Injection
- **Storage**: In-memory (development), ready for PostgreSQL

## Integration with Account Service

To enable full event-driven flow, the **account service** needs to:

1. **Publish Events** to Kafka topic `account-events`
2. **Event Types**:
   - `account.created` - when account is created
   - `account.status_changed` - when status changes (DELETED, BLOCKED)
3. **Event Payload**: Include `account_id` and `status`

Example Kafka producer in account service:
```go
// After creating account
event := AccountEvent{
    Type:      "account.created",
    AccountID: account.ID,
    Status:    string(account.Status),
}
kafkaProducer.Publish("account-events", event)
```
