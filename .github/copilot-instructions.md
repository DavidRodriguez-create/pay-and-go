# Copilot Instructions for pay-and-go

## 1. Project Overview

**Purpose:** Demo payment application to learn about Golang, showcasing microservices architecture with event-driven communication.

**Architecture:** Clean Architecture with strict layer separation:
- **Domain** (innermost): Pure business logic, NO external dependencies
- **Application**: Use cases, DTOs for serialization, orchestrates domain
- **Infrastructure**: External concerns (Kafka, repositories, HTTP clients)
- **Presentation**: HTTP handlers, routing, request/response formatting

**Tech Stack:**
- Go 1.23
- Apache Kafka (event streaming)
- Podman (containerization)
- In-memory storage (non-persistent)

---

## 2. Critical Architecture Rules

### 🚨 NEVER add JSON tags to Domain entities

Domain layer must remain pure - NO serialization concerns, NO external dependencies.

**❌ WRONG:**
```go
// domain/account.go
type Account struct {
    ID     string `json:"id"`  // ❌ JSON tags in domain
    Status string `json:"status"`
}
```

**✅ CORRECT:**
```go
// domain/account.go
type Account struct {
    ID     string  // ✅ Pure Go struct
    Status AccountStatus
}

// application/dtos.go
type AccountResponse struct {
    ID     string `json:"id"`  // ✅ JSON tags in DTOs
    Status string `json:"status"`
}
```

**Why:** Domain entities represent business concepts independent of serialization format. DTOs handle translation between domain and external systems.

---

## 3. Layer Responsibilities

### Domain Layer (`domain/`)
- **Contains:** Entities, value objects, interfaces (Repository, EventPublisher)
- **Rules:**
  - NO external dependencies (no JSON, HTTP, Kafka libraries)
  - NO implementation details (only interfaces)
  - Pure business logic and validation
  - Constructor functions return pointers and errors

**Example:**
```go
// domain/account.go
type Account struct {
    ID            string
    AccountNumber string
    Status        AccountStatus  // Custom type, not string
}

func NewAccount(id, accountNumber, beholderName, countryCode string) (*Account, error) {
    if id == "" || accountNumber == "" {
        return nil, errors.New("required fields missing")
    }
    return &Account{...}, nil
}

// domain/account_repository.go
type AccountRepository interface {
    Create(account *Account) error
    GetByID(id string) (*Account, error)
    Update(account *Account) error
}

// domain/event_publisher.go
type EventPublisher interface {
    PublishAccountCreated(accountID string, status string) error
    PublishAccountStatusChanged(accountID string, status string) error
}
```

### Application Layer (`application/`)
- **Contains:** Use cases, DTOs, service implementations, mappers
- **Rules:**
  - DTOs have JSON tags for API serialization
  - Coordinates domain entities via repositories
  - Publishes events via EventPublisher interface
  - Maps between domain entities and DTOs

**Example:**
```go
// application/dtos.go
type CreateAccountRequest struct {
    AccountNumber string `json:"account_number"`
    BeholderName  string `json:"beholder_name"`
    CountryCode   string `json:"country_code"`
}

type AccountResponse struct {
    ID            string `json:"id"`
    AccountNumber string `json:"account_number"`
    Status        string `json:"status"`
}

// application/service.go
type AccountServiceImpl struct {
    repository     domain.AccountRepository
    eventPublisher domain.EventPublisher
}

func (s *AccountServiceImpl) CreateAccount(req CreateAccountRequest) (*AccountResponse, error) {
    account, err := domain.NewAccount(req.ID, req.AccountNumber, req.BeholderName, req.CountryCode)
    if err != nil {
        return nil, err
    }
    
    if err := s.repository.Create(account); err != nil {
        return nil, err
    }
    
    s.eventPublisher.PublishAccountCreated(account.ID, string(account.Status))
    
    return MapAccountToResponse(account), nil
}
```

### Infrastructure Layer (`infrastructure/`)
- **Contains:** Repository implementations, Kafka producer/consumer, external clients
- **Rules:**
  - Implements domain interfaces
  - Handles external system communication
  - Converts domain errors to infrastructure errors

**Example:**
```go
// infrastructure/memory_account_repository.go
type MemoryAccountRepository struct {
    accounts map[string]*domain.Account
    mu       sync.RWMutex
}

func (r *MemoryAccountRepository) Create(account *domain.Account) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.accounts[account.ID] = account
    return nil
}

// infrastructure/kafka_producer.go
type KafkaProducer struct {
    writer *kafka.Writer
}

func (p *KafkaProducer) PublishAccountCreated(accountID string, status string) error {
    event := AccountEvent{Type: "account.created", AccountID: accountID, Status: status}
    value, _ := json.Marshal(event)
    return p.writer.WriteMessages(context.Background(), kafka.Message{Value: value})
}
```

### Presentation Layer (`presentation/`)
- **Contains:** HTTP controllers, routes, response presenters
- **Rules:**
  - Controllers: One per use case, delegate to application service
  - Routes: RESTful design with CORS middleware
  - Query parameters for resource IDs, request body for data

**Example:**
```go
// presentation/controllers/create_account_controller.go
type CreateAccountController struct {
    service application.AccountService
}

func (c *CreateAccountController) Handle(w http.ResponseWriter, r *http.Request) {
    var req application.CreateAccountRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    resp, err := c.service.CreateAccount(req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

// presentation/routes/routes.go
func SetupRoutes(ctrls *Controllers) *http.ServeMux {
    mux := http.NewServeMux()
    
    // Collection endpoint (plural) - GET /accounts
    mux.HandleFunc("/accounts", corsMiddleware(handleAccountList(ctrls)))
    
    // Single resource endpoint (singular)
    // POST /account - Create
    // GET /account?id=xxx - Read
    // PUT /account?id=xxx - Update
    // DELETE /account?id=xxx - Delete
    mux.HandleFunc("/account", corsMiddleware(handleAccount(ctrls)))
    
    return mux
}
```

---

## 4. API Design Conventions

### Endpoint Structure
- **Plural `/accounts`**: Collection operations (list all)
- **Singular `/account`**: Single resource operations (create, get, update, delete)
- **Query parameters**: Resource IDs (e.g., `?id=123`)
- **Request body**: Data for create/update operations

### Example Requests
```bash
# Create account (no ID in URL)
POST /account
Body: {"account_number": "ACC001", "beholder_name": "John", "country_code": "US"}

# Get account by ID
GET /account?id=123

# Update account (ID in query, data in body)
PUT /account?id=123
Body: {"status": "BLOCKED"}

# Delete account
DELETE /account?id=123

# List all accounts
GET /accounts

# Search by account number
GET /accounts/by-number?account_number=ACC001
```

### CORS Middleware Pattern
```go
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next(w, r)
    }
}
```

---

## 5. Event-Driven Architecture

### Kafka Event Flow
1. Account service publishes events to `account-events` topic
2. Card service consumes events and updates local account cache
3. **Sync delay**: 2-3 seconds between services (eventual consistency)

### Event Types
```go
type AccountEvent struct {
    Type      string `json:"type"`       // "account.created" | "account.status_changed"
    AccountID string `json:"account_id"`
    Status    string `json:"status"`     // "ACTIVE" | "BLOCKED" | "DELETED"
}
```

### When to Publish Events
- **account.created**: After successful account creation
- **account.status_changed**: After update (status change) or delete

**Example:**
```go
// application/delete_account.go
func (s *AccountServiceImpl) DeleteAccount(id string) error {
    if err := s.repository.Delete(id); err != nil {
        return err
    }
    
    // 🔔 CRITICAL: Always publish event after state change
    return s.eventPublisher.PublishAccountStatusChanged(id, "DELETED")
}
```

### Consumer Pattern (Card Service)
```go
// infrastructure/kafka_account_consumer.go
func (c *KafkaAccountConsumer) Start(ctx context.Context) error {
    for {
        msg, err := c.reader.ReadMessage(ctx)
        var event AccountEvent
        json.Unmarshal(msg.Value, &event)
        
        switch event.Type {
        case "account.created":
            c.accountRepo.SaveAccount(event.AccountID, event.Status)
        case "account.status_changed":
            if event.Status == "DELETED" {
                c.accountRepo.DeleteAccount(event.AccountID)
            } else {
                c.accountRepo.UpdateAccountStatus(event.AccountID, event.Status)
            }
        }
    }
}
```

---

## 6. Testing Strategy

### Table-Driven Tests (Standard Go Pattern)
```go
func TestNewAccount(t *testing.T) {
    tests := []struct {
        name          string
        id            string
        accountNumber string
        wantErr       bool
    }{
        {
            name:          "Valid account creation",
            id:            "123",
            accountNumber: "ACC001",
            wantErr:       false,
        },
        {
            name:          "Missing ID",
            id:            "",
            accountNumber: "ACC001",
            wantErr:       true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            account, err := domain.NewAccount(tt.id, tt.accountNumber, ...)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("NewAccount() error = %v, wantErr %v", err, tt.wantErr)
            }
            
            if !tt.wantErr && account.ID != tt.id {
                t.Errorf("NewAccount() ID = %v, want %v", account.ID, tt.id)
            }
        })
    }
}
```

### Test Organization
```
tests/
  unit/
    domain/        # Pure business logic tests (100% coverage goal)
    application/   # Use case tests with mocks
    infrastructure/
  integration/     # End-to-end tests with real dependencies
```

### Running Tests
```bash
# Run all tests
go test ./tests/... -v

# Run specific package
go test ./tests/unit/domain -v

# Coverage report
go test ./tests/... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## 7. Deployment & Operations

### Primary Command: `./manage-services.sh`
**Location:** Root of repository

```bash
# Start all services (builds images, runs containers, opens browser)
./manage-services.sh start

# Start without opening browser
./manage-services.sh start --no-browser

# Stop all services
./manage-services.sh stop

# Restart services (preserves Kafka data)
./manage-services.sh restart

# Check running containers
./manage-services.sh status
```

### Container Architecture
- **Zookeeper**: Port 2181 (Kafka coordination)
- **Kafka**: Port 9092 (event streaming)
- **Account Service**: Port 8080 (HTTP API)
- **Card Service**: Port 8081 (HTTP API)
- **UI Dashboard**: `ui.html` (static file, opens automatically)

### Important Notes
- **Podman only** - No Docker daemon required
- **In-memory storage** - Data lost on container restart
- **No docker-compose** - Uses individual `podman run` commands
- **Kafka startup**: 8-second wait before starting services (allows Kafka to initialize)
- **Sync delay**: 2-3 seconds after account creation before card operations

### Kubernetes Deployment
```bash
# Deploy to Kubernetes
./k8s/deploy.sh

# Remove from Kubernetes
./k8s/undeploy.sh
```

---

## 8. Common Tasks & Patterns

### Adding a New Use Case
1. Define domain entity/method in `domain/`
2. Create DTO in `application/dtos.go`
3. Implement use case in `application/<use_case>.go`
4. Create controller in `presentation/controllers/`
5. Add route in `presentation/routes/routes.go`
6. Write tests in `tests/unit/application/`

### Adding Event Publishing
```go
// After ANY state change in application layer:
account.Status = domain.StatusBlocked
s.repository.Update(account)

// ✅ Always publish event
s.eventPublisher.PublishAccountStatusChanged(account.ID, string(account.Status))
```

### Debugging Kafka Sync Issues
- Check Kafka logs: `podman logs kafka`
- Check consumer logs: `podman logs card-service`
- Verify event structure matches consumer expectations
- Remember 2-3 second sync delay (UI shows warning)

### Reading Query Parameters in Controller
```go
func (c *UpdateAccountController) Handle(w http.ResponseWriter, r *http.Request) {
    // Read ID from query parameter
    id := r.URL.Query().Get("id")
    if id == "" {
        http.Error(w, "ID is required", http.StatusBadRequest)
        return
    }
    
    // Read data from request body
    var req application.UpdateAccountRequest
    json.NewDecoder(r.Body).Decode(&req)
    req.ID = id  // Inject ID from query param
    
    c.service.UpdateAccount(req)
}
```

---

## 9. Code Style & Conventions

### Go Standards
- Use `gofmt` for formatting (automatically applied)
- Exported names start with capital letter (e.g., `Account`, not `account`)
- Unexported helpers start with lowercase (e.g., `validateID()`)
- Constructor pattern: `New<Type>(params) (*Type, error)`

### Error Handling
```go
// ✅ Return errors, don't panic
account, err := domain.NewAccount(...)
if err != nil {
    return nil, err
}

// ✅ Wrap errors with context
if err := s.repository.Create(account); err != nil {
    return nil, fmt.Errorf("failed to create account: %w", err)
}
```

### Logging
```go
// Use log package for infrastructure concerns
log.Printf("Published event: type=%s, account_id=%s", event.Type, event.AccountID)

// Domain/Application layers should NOT log (return errors instead)
```

### Dependency Injection
```go
// Services receive dependencies via constructor
func NewAccountService(
    repository domain.AccountRepository,
    eventPublisher domain.EventPublisher,
) *AccountServiceImpl {
    return &AccountServiceImpl{
        repository:     repository,
        eventPublisher: eventPublisher,
    }
}

// Wire dependencies in main.go
func main() {
    repo := infrastructure.NewMemoryAccountRepository()
    publisher := infrastructure.NewKafkaProducer(brokers, topic)
    service := application.NewAccountService(repo, publisher)
    
    ctrls := &routes.Controllers{
        CreateAccount: controllers.NewCreateAccountController(service),
        // ...
    }
    
    mux := routes.SetupRoutes(ctrls)
    http.ListenAndServe(":8080", mux)
}
```

---

## 10. Quick Reference

### File Structure Pattern
```
services/<service-name>/
  cmd/
    main.go              # Entry point, dependency wiring
  domain/
    <entity>.go          # ❌ NO JSON tags
    <entity>_repository.go  # Interfaces only
    event_publisher.go   # Interfaces only
  application/
    dtos.go              # ✅ JSON tags HERE
    service.go           # Service interface
    <use_case>.go        # Use case implementations
    mappers.go           # Domain ↔ DTO conversion
  infrastructure/
    memory_<entity>_repository.go
    kafka_producer.go
    kafka_consumer.go
  presentation/
    controllers/
      <use_case>_controller.go
    presenters/
      response_presenter.go
    routes/
      routes.go          # HTTP routing + CORS
  tests/
    unit/
      domain/
      application/
      infrastructure/
    integration/
```

### Environment Variables
```bash
# Account Service (.env)
PORT=8080
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC=account-events

# Card Service (.env)
PORT=8081
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC=account-events
KAFKA_GROUP_ID=card-service-group
```

### Useful Commands
```bash
# Start services
./manage-services.sh start

# Run tests
go test ./tests/... -v

# Check Kafka messages
podman exec -it kafka kafka-console-consumer --bootstrap-server localhost:9092 --topic account-events --from-beginning

# View logs
podman logs account-service
podman logs card-service

# Rebuild single service
podman build -f Dockerfile.account -t account-service:latest .
```

---

## 11. Troubleshooting

### "Account not found" in Card Service
**Cause:** Kafka sync delay (2-3 seconds)
**Solution:** Wait briefly after creating account, UI shows warning

### CORS Errors in Browser
**Cause:** Missing CORS middleware
**Solution:** Ensure `corsMiddleware` wraps all handlers in `routes.go`

### Kafka Connection Refused
**Cause:** Kafka not fully started
**Solution:** `manage-services.sh` includes 8-second wait, manual deploy needs `sleep 8`

### Tests Failing After Code Change
**Cause:** Likely domain/application boundary violation
**Solution:** Check for JSON tags in domain entities, external dependencies in domain layer

---

## 12. Learning Resources

- **Clean Architecture:** "Clean Architecture" by Robert C. Martin
- **Go Testing:** https://go.dev/doc/tutorial/add-a-test
- **Kafka Basics:** https://kafka.apache.org/intro
- **RESTful API Design:** https://restfulapi.net/

---

**Last Updated:** Generated from codebase analysis
**Maintainer:** Ask questions to improve this document!

