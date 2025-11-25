# Account Service Tests

This directory contains all tests for the account service, organized by test type and layer.

## Directory Structure

```
tests/
├── integration/              # Integration/End-to-end tests
│   └── integration_test.go   # Full API lifecycle tests
├── unit/                     # Unit tests organized by layer
│   ├── application/          # Application layer (use cases) tests
│   │   ├── create_account_test.go
│   │   ├── delete_account_test.go
│   │   ├── update_account_test.go
│   │   └── view_account_test.go
│   ├── domain/              # Domain layer (entities) tests
│   │   └── account_test.go
│   └── infrastructure/      # Infrastructure layer (repository) tests
│       └── memory_account_repository_test.go
└── README.md                # This file
```

## Test Organization

### Unit Tests (`tests/unit/`)
Unit tests are organized by architectural layer, matching the clean architecture structure:

#### Domain Layer (`tests/unit/domain/`)
- Tests domain entities and their business logic
- Package: `domain_test`
- Coverage: Entity creation, validation, state management

#### Application Layer (`tests/unit/application/`)
- Tests use cases and application services
- Package: `application_test`
- Uses `MockAccountRepository` for isolation
- Coverage: All CRUD operations and business flows

#### Infrastructure Layer (`tests/unit/infrastructure/`)
- Tests infrastructure implementations
- Package: `infrastructure_test`
- Coverage: Repository operations, data persistence

### Integration Tests (`tests/integration/`)
- Tests complete API flows end-to-end
- Package: `tests`
- Coverage: Full HTTP request/response cycles, controller → service → repository

## Running Tests

### All Tests
```bash
go test ./tests/...
```

### Unit Tests Only
```bash
go test ./tests/unit/...
```

### Integration Tests Only
```bash
go test ./tests/integration/...
```

### Specific Layer
```bash
# Domain tests
go test ./tests/unit/domain/...

# Application tests
go test ./tests/unit/application/...

# Infrastructure tests
go test ./tests/unit/infrastructure/...
```

### With Verbose Output
```bash
go test ./tests/... -v
```

### With Coverage
```bash
go test ./tests/... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

### With Coverage HTML Report
```bash
go test ./tests/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## Test Naming Conventions

- Test files: `*_test.go`
- Test functions: `Test<FunctionName>`
- Test subtests: Use `t.Run()` with descriptive names
- Package naming: `<package>_test` (black-box testing)

## Mock Objects

The `MockAccountRepository` is defined in `tests/unit/application/create_account_test.go` and provides:
- Configurable behavior for each repository method
- Error injection for testing failure scenarios
- Zero dependencies for fast execution

## Test Coverage

Current coverage:
- **Domain Layer**: 100%
- **Application Layer**: 100%
- **Infrastructure Layer**: 97.7%
- **Integration Tests**: Full API coverage

## Best Practices

1. **Isolation**: Each test is independent and uses fresh dependencies
2. **Clarity**: Descriptive test names explaining what is being tested
3. **AAA Pattern**: Arrange-Act-Assert structure
4. **Table-Driven**: Multiple test cases in single test functions
5. **Error Paths**: Both success and failure scenarios tested
6. **Black-Box**: Tests use `*_test` packages to test public APIs

## Adding New Tests

When adding new functionality:

1. **Domain Changes**: Add tests to `tests/unit/domain/`
2. **Use Cases**: Add tests to `tests/unit/application/`
3. **Repository**: Add tests to `tests/unit/infrastructure/`
4. **API Endpoints**: Add tests to `tests/integration/`

Example:
```go
// tests/unit/application/new_feature_test.go
package application_test

import (
    "testing"
    "github.com/DavidRodriguez-create/pay-and-go/services/account/application"
)

func TestNewFeature(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        // test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

## Continuous Integration

These tests are designed to run in CI/CD pipelines:

```yaml
# Example CI configuration
test:
  script:
    - cd services/account
    - go test ./tests/... -v -race -coverprofile=coverage.out
    - go tool cover -func=coverage.out
  coverage: '/total.*?(\d+\.\d+)%/'
```

## Troubleshooting

### Tests Not Running
```bash
# Ensure you're in the correct directory
cd services/account

# Clean Go cache
go clean -testcache

# Run with verbose output to see errors
go test ./tests/... -v
```

### Import Issues
```bash
# Download dependencies
go mod download

# Tidy up modules
go mod tidy
```

### Coverage Not Generating
```bash
# Ensure coverage directory exists
mkdir -p coverage

# Generate coverage with explicit output
go test ./tests/... -coverprofile=coverage/coverage.out
```

---

**Last Updated**: November 25, 2025  
**Test Count**: 35+ test cases  
**Status**: All tests passing ✅
