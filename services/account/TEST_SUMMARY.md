# Account Service Test Suite

## Overview
Comprehensive unit and integration tests for the account service following clean architecture principles.

## Test Coverage Summary

### Overall Coverage: 97.7%+

#### By Layer:
- **Domain Layer**: 100% coverage
- **Application Layer**: 100% coverage  
- **Infrastructure Layer**: 97.7% coverage
- **Integration Tests**: Full API lifecycle coverage

## Test Structure

### 1. Domain Layer Tests (`domain/account_test.go`)

#### Tests Implemented:
- ✅ `TestNewAccount` - Account entity creation with validation
  - Valid account creation
  - Missing ID validation
  - Missing account number validation
  - Missing beholder name validation
  - Missing country code validation
  
- ✅ `TestAccountStatusMethods` - Account status checks
  - IsActive() method
  - IsDeleted() method
  - IsBlocked() method

**Coverage**: 100% of domain logic

---

### 2. Application Layer Tests

#### Create Account (`application/create_account_test.go`)
- ✅ Successful account creation
- ✅ Missing beholder name validation
- ✅ Missing country code validation
- ✅ Repository creation error handling

#### Update Account (`application/update_account_test.go`)
- ✅ Successful account update
- ✅ Update all fields
- ✅ Account not found error
- ✅ Update deleted account prevention
- ✅ Repository update error handling

#### Delete Account (`application/delete_account_test.go`)
- ✅ Successful account deletion (soft delete)
- ✅ Account not found error
- ✅ Delete already deleted account prevention
- ✅ Repository delete error handling

#### View Accounts (`application/view_account_test.go`)
- ✅ Get account by ID - successful retrieval
- ✅ Get account by ID - account not found
- ✅ Get account by account number - successful retrieval
- ✅ Get account by account number - not found
- ✅ List accounts - with accounts
- ✅ List accounts - empty list
- ✅ List accounts - repository error

**Coverage**: 100% of application use cases

---

### 3. Infrastructure Layer Tests (`infrastructure/memory_account_repository_test.go`)

#### Repository Operations:
- ✅ Create and retrieve account
- ✅ Create duplicate account error
- ✅ Get by account number
- ✅ Get non-existent account error
- ✅ Update account
- ✅ Update non-existent account error
- ✅ Delete account (soft delete)
- ✅ Delete non-existent account error
- ✅ List accounts
- ✅ List empty repository

**Coverage**: 97.7% of repository implementation

---

### 4. Integration Tests (`tests/integration_test.go`)

#### Complete Account Lifecycle Test:
1. ✅ Create an account via POST /accounts
2. ✅ Get account by ID via GET /accounts/?id={id}
3. ✅ Get account by account number via GET /accounts/by-number?account_number={number}
4. ✅ Update account via PUT /accounts/?id={id}
5. ✅ Verify update
6. ✅ List all accounts via GET /accounts
7. ✅ Delete account via DELETE /accounts/?id={id}
8. ✅ Verify soft delete status

#### Multiple Accounts Test:
- ✅ Create multiple accounts
- ✅ List all accounts and verify count

#### Health Endpoint Test:
- ✅ GET /health returns healthy status

#### Error Handling Tests:
- ✅ Invalid JSON payload
- ✅ Missing required fields
- ✅ Get non-existent account (404)
- ✅ Update non-existent account (400)
- ✅ Delete non-existent account (400)
- ✅ Method not allowed (405)

---

## Running Tests

### Run All Tests
```bash
cd services/account
go test ./... -v
```

### Run with Coverage
```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

### Run Specific Test Package
```bash
# Domain tests only
go test ./domain -v

# Application tests only
go test ./application -v

# Infrastructure tests only
go test ./infrastructure -v

# Integration tests only
go test ./tests -v
```

### Generate HTML Coverage Report
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

---

## Test Results

```
=== Test Summary ===
PASS: application   (100.0% coverage)
PASS: domain        (100.0% coverage)
PASS: infrastructure (97.7% coverage)
PASS: tests         (integration tests)

Total: 35+ test cases
Status: ALL PASSING ✅
```

---

## Mock Implementation

A `MockAccountRepository` is provided in `application/create_account_test.go` for testing application layer in isolation. This mock implements the `domain.AccountRepository` interface with configurable behavior.

### Mock Features:
- Configurable return values for each method
- Error injection for testing error paths
- No external dependencies
- Fast execution

---

## Test Best Practices Applied

1. **Arrange-Act-Assert Pattern**: All tests follow AAA structure
2. **Table-Driven Tests**: Multiple test cases in single test functions
3. **Isolation**: Each test is independent and uses fresh dependencies
4. **Clear Naming**: Descriptive test names explaining what is being tested
5. **Error Path Testing**: Both success and failure scenarios covered
6. **Edge Cases**: Boundary conditions and invalid inputs tested
7. **Integration Testing**: Full API lifecycle tested end-to-end

---

## Future Enhancements

### Potential Additions:
- [ ] Performance/benchmark tests
- [ ] Concurrent operation tests
- [ ] Database integration tests (when moving from in-memory to real DB)
- [ ] API contract tests
- [ ] Load testing
- [ ] Chaos testing for error resilience

---

## Continuous Integration

These tests are designed to run in CI/CD pipelines:

```yaml
# Example CI configuration
test:
  script:
    - cd services/account
    - go test ./... -v -race -coverprofile=coverage.out
    - go tool cover -func=coverage.out
  coverage: '/total.*?(\d+\.\d+)%/'
```

---

## Manual API Testing

The service is deployed and can be tested manually:

```bash
# Health check
curl http://localhost:8081/health

# Create account
curl -X POST http://localhost:8081/accounts \
  -H "Content-Type: application/json" \
  -d '{"beholder_name":"Test User","country_code":"US"}'

# List accounts
curl http://localhost:8081/accounts

# Get by ID
curl "http://localhost:8081/accounts/?id={id}"

# Update account
curl -X PUT "http://localhost:8081/accounts/?id={id}" \
  -H "Content-Type: application/json" \
  -d '{"id":"{id}","beholder_name":"Updated Name","country_code":"UK"}'

# Delete account
curl -X DELETE "http://localhost:8081/accounts/?id={id}"
```

---

## Test Maintenance

- Tests are co-located with implementation files for easy maintenance
- Mock implementations are shared across test files
- Test helpers reduce code duplication
- Clear documentation makes tests serve as usage examples

---

**Last Updated**: November 25, 2025
**Test Suite Version**: 1.0
**Go Version**: 1.23
