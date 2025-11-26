# Event-Driven Architecture: Account + Card Services

## Overview

This demonstrates how the account and card services work together using Kafka for event-driven communication.

## Architecture Diagram

```
┌─────────────────────┐                    ┌──────────────────────┐
│  Account Service    │                    │   Card Service       │
│     (Port 8081)     │                    │    (Port 8082)       │
└──────────┬──────────┘                    └──────────┬───────────┘
           │                                           │
           │  1. POST /account                        │
           │     Create Account                       │
           │     ───────────────                      │
           │                                           │
           │  2. Publish Event                        │
           │     ─────────────────►                   │
           │                       │                  │
           │                       ▼                  │
           │              ┌─────────────────┐         │
           │              │  Kafka Broker   │         │
           │              │ account-events  │         │
           │              └────────┬────────┘         │
           │                       │                  │
           │                       │  3. Consume      │
           │                       └─────────────────►│
           │                                          │
           │                          4. Update       │
           │                             AccountCache │
           │                             in memory    │
           │                                          │
           │                       5. POST /card      │
           │                          ────────────────►│
           │                          Create Card     │
           │                          (validates      │
           │                           account cache) │
```

## Event Flow

### 1. Account Creation Flow

**Account Service:**
```bash
# Client creates account
POST /account
{
  "beholder_name": "John Doe",
  "country_code": "US"
}

# Account service:
# 1. Creates account in DB
# 2. Publishes Kafka event
```

**Kafka Event:**
```json
{
  "type": "account.created",
  "account_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "ACTIVE"
}
```

**Card Service:**
```
# Kafka consumer receives event
# Updates local AccountCache
# AccountCache[550e8400...] = {ID: "550e8400...", Status: "ACTIVE"}
```

### 2. Card Creation Flow

**Card Service:**
```bash
# Client creates card
POST /card
{
  "country": "US",
  "account_id": "550e8400-e29b-41d4-a716-446655440000"
}

# Card service:
# 1. Checks AccountCache (in-memory)
# 2. Validates account is ACTIVE
# 3. Creates card
# 4. Returns card details
```

### 3. Account Status Change Flow

**Account Service:**
```bash
# Client deletes account
DELETE /account?id=550e8400-e29b-41d4-a716-446655440000

# Account service:
# 1. Marks account as DELETED
# 2. Publishes Kafka event
```

**Kafka Event:**
```json
{
  "type": "account.status_changed",
  "account_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "DELETED"
}
```

**Card Service:**
```
# Kafka consumer receives event
# Updates local AccountCache
# AccountCache[550e8400...] = {ID: "550e8400...", Status: "DELETED"}

# Future card creation attempts will fail:
# Error: "cannot create card for deleted account"
```

## Running the Complete System

### Prerequisites

1. **Kafka Cluster** (using Docker)
```bash
# Run Kafka with Podman/Docker
podman run -d --name kafka \
  -p 9092:9092 \
  -e KAFKA_ENABLE_KRAFT=yes \
  -e KAFKA_CFG_PROCESS_ROLES=broker,controller \
  -e KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER \
  -e KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093 \
  -e KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT \
  -e KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092 \
  -e KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=1@localhost:9093 \
  -e KAFKA_CFG_NODE_ID=1 \
  bitnami/kafka:latest
```

2. **Create Kafka Topic**
```bash
podman exec -it kafka kafka-topics.sh \
  --create \
  --topic account-events \
  --bootstrap-server localhost:9092 \
  --partitions 3 \
  --replication-factor 1
```

### Start Services

**Terminal 1: Account Service**
```bash
cd services/account
export PORT=8081
export KAFKA_BROKERS=localhost:9092
export KAFKA_TOPIC=account-events
go run cmd/main.go
```

> 📝 **Note**: Account service now publishes events to Kafka when:
> - Creating accounts (`account.created` event)
> - Updating account status (`account.status_changed` event)

**Terminal 2: Card Service**
```bash
cd services/card
export PORT=8082
export KAFKA_BROKERS=localhost:9092
export KAFKA_TOPIC=account-events
export KAFKA_GROUP_ID=card-service
go run cmd/main.go
```

## Testing the Integration

### Complete Flow Test

```bash
# 1. Create an account
curl -X POST http://localhost:8081/account \
  -H "Content-Type: application/json" \
  -d '{
    "beholder_name": "John Doe",
    "country_code": "US"
  }'

# Response:
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "account_number": "ACC-12345678",
  "beholder_name": "John Doe",
  "country_code": "US",
  "status": "ACTIVE"
}

# 2. Wait a moment for Kafka event to be consumed

# 3. Create a card for the account
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

# 4. Get all cards for the account
curl http://localhost:8082/cards/by-account?account_id=550e8400-e29b-41d4-a716-446655440000

# 5. Delete the account
curl -X DELETE http://localhost:8081/account?id=550e8400-e29b-41d4-a716-446655440000

# 6. Try to create another card (should fail)
curl -X POST http://localhost:8082/card \
  -H "Content-Type: application/json" \
  -d '{
    "country": "US",
    "account_id": "550e8400-e29b-41d4-a716-446655440000"
  }'

# Response:
{
  "error": "cannot create card for deleted account"
}
```

## Required Account Service Changes

To enable this integration, the **account service** needs Kafka producer integration:

### 1. Add Kafka Dependency

```bash
cd services/account
go get github.com/segmentio/kafka-go
```

### 2. Create Kafka Producer (infrastructure/kafka_producer.go)

```go
package infrastructure

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
)

type AccountEventPublisher struct {
	writer *kafka.Writer
}

func NewAccountEventPublisher(brokers []string, topic string) *AccountEventPublisher {
	return &AccountEventPublisher{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *AccountEventPublisher) PublishCreated(accountID, status string) error {
	event := map[string]string{
		"type":       "account.created",
		"account_id": accountID,
		"status":     status,
	}
	return p.publish(event)
}

func (p *AccountEventPublisher) PublishStatusChanged(accountID, status string) error {
	event := map[string]string{
		"type":       "account.status_changed",
		"account_id": accountID,
		"status":     status,
	}
	return p.publish(event)
}

func (p *AccountEventPublisher) publish(event interface{}) error {
	data, _ := json.Marshal(event)
	return p.writer.WriteMessages(context.Background(), kafka.Message{Value: data})
}

func (p *AccountEventPublisher) Close() error {
	return p.writer.Close()
}
```

### 3. Integrate in Use Cases

**create_account.go:**
```go
// After creating account
if err := uc.accountRepo.Create(account); err != nil {
	return nil, err
}

// Publish event
uc.eventPublisher.PublishCreated(account.ID, string(account.Status))

return AccountToResponse(account), nil
```

**delete_account.go:**
```go
// After deleting account
if err := uc.accountRepo.Delete(req.ID); err != nil {
	return err
}

// Publish event
uc.eventPublisher.PublishStatusChanged(req.ID, "DELETED")

return nil
```

## Benefits of This Architecture

### 1. **Loose Coupling**
- Services don't call each other directly
- Can deploy/scale independently
- Can fail independently

### 2. **High Performance**
- Card creation is fast (no HTTP calls)
- Only local cache lookup
- Async event processing

### 3. **Resilience**
- Card service works if account service is down
- Events are durable in Kafka
- Can replay events if needed

### 4. **Scalability**
- Can run multiple card service instances
- Kafka consumer group distributes load
- Can handle high card creation volumes

### 5. **Eventual Consistency**
- Account status eventually consistent
- Acceptable for most use cases
- Can add sync endpoints if needed

## Monitoring

### Check Kafka Consumer Status

```bash
# List consumer groups
podman exec -it kafka kafka-consumer-groups.sh \
  --bootstrap-server localhost:9092 \
  --list

# Check consumer lag
podman exec -it kafka kafka-consumer-groups.sh \
  --bootstrap-server localhost:9092 \
  --group card-service \
  --describe
```

### Service Logs

**Card Service Logs:**
```
Kafka consumer started successfully
Received account event: type=account.created, account_id=550e8400..., status=ACTIVE
Updated account cache: account_id=550e8400..., status=ACTIVE
```

## Troubleshooting

### Card Creation Fails with "account not found"

**Cause:** Account event not yet consumed

**Solution:**
1. Check Kafka consumer is running
2. Verify Kafka broker is accessible
3. Check consumer logs for errors
4. Verify topic exists and has messages

### Events Not Being Consumed

**Cause:** Kafka connection issues

**Solution:**
1. Verify KAFKA_BROKERS environment variable
2. Check network connectivity to Kafka
3. Ensure topic exists: `kafka-topics.sh --list`
4. Check consumer group: `kafka-consumer-groups.sh --describe`

### High Consumer Lag

**Cause:** Card service processing too slowly

**Solution:**
1. Scale card service horizontally (multiple instances)
2. Increase Kafka partition count
3. Optimize event processing logic
4. Add caching/batching
