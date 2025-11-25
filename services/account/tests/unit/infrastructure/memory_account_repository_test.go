package infrastructure_test

import (
"testing"
"time"

"github.com/DavidRodriguez-create/pay-and-go/services/account/domain"
"github.com/DavidRodriguez-create/pay-and-go/services/account/infrastructure"
)

func TestInMemoryAccountRepository(t *testing.T) {
repo := infrastructure.NewInMemoryAccountRepository()

t.Run("Create and retrieve account", func(t *testing.T) {
account, _ := domain.NewAccount("123", "ACC001", "John Doe", "US")

err := repo.Create(account)
if err != nil {
t.Fatalf("Failed to create account: %v", err)
}

retrieved, err := repo.GetByID("123")
if err != nil {
t.Fatalf("Failed to retrieve account: %v", err)
}

if retrieved.ID != account.ID {
t.Errorf("Expected ID %s, got %s", account.ID, retrieved.ID)
}
})

t.Run("Create duplicate account", func(t *testing.T) {
repo := infrastructure.NewInMemoryAccountRepository()
account, _ := domain.NewAccount("123", "ACC001", "John Doe", "US")

repo.Create(account)
err := repo.Create(account)

if err == nil {
t.Error("Expected error when creating duplicate account")
}
})

t.Run("Get by account number", func(t *testing.T) {
repo := infrastructure.NewInMemoryAccountRepository()
account, _ := domain.NewAccount("123", "ACC001", "John Doe", "US")
repo.Create(account)

retrieved, err := repo.GetByAccountNumber("ACC001")
if err != nil {
t.Fatalf("Failed to retrieve account: %v", err)
}

if retrieved.AccountNumber != "ACC001" {
t.Errorf("Expected account number ACC001, got %s", retrieved.AccountNumber)
}
})

t.Run("Get non-existent account", func(t *testing.T) {
repo := infrastructure.NewInMemoryAccountRepository()

_, err := repo.GetByID("nonexistent")
if err == nil {
t.Error("Expected error when getting non-existent account")
}
})

t.Run("Update account", func(t *testing.T) {
repo := infrastructure.NewInMemoryAccountRepository()
account, _ := domain.NewAccount("123", "ACC001", "John Doe", "US")
repo.Create(account)

account.BeholderName = "Jane Doe"
account.CountryCode = "UK"
account.UpdatedAt = time.Now()

err := repo.Update(account)
if err != nil {
t.Fatalf("Failed to update account: %v", err)
}

retrieved, _ := repo.GetByID("123")
if retrieved.BeholderName != "Jane Doe" {
t.Errorf("Expected beholder name Jane Doe, got %s", retrieved.BeholderName)
}
if retrieved.CountryCode != "UK" {
t.Errorf("Expected country code UK, got %s", retrieved.CountryCode)
}
})

t.Run("Update non-existent account", func(t *testing.T) {
repo := infrastructure.NewInMemoryAccountRepository()
account, _ := domain.NewAccount("123", "ACC001", "John Doe", "US")

err := repo.Update(account)
if err == nil {
t.Error("Expected error when updating non-existent account")
}
})

t.Run("Delete account", func(t *testing.T) {
repo := infrastructure.NewInMemoryAccountRepository()
account, _ := domain.NewAccount("123", "ACC001", "John Doe", "US")
repo.Create(account)

err := repo.Delete("123")
if err != nil {
t.Fatalf("Failed to delete account: %v", err)
}

retrieved, _ := repo.GetByID("123")
if retrieved.Status != domain.StatusDeleted {
t.Errorf("Expected status DELETED, got %s", retrieved.Status)
}
})

t.Run("Delete non-existent account", func(t *testing.T) {
repo := infrastructure.NewInMemoryAccountRepository()

err := repo.Delete("nonexistent")
if err == nil {
t.Error("Expected error when deleting non-existent account")
}
})

t.Run("List accounts", func(t *testing.T) {
repo := infrastructure.NewInMemoryAccountRepository()

for i := 0; i < 3; i++ {
account, _ := domain.NewAccount(
string(rune('1'+i)),
"ACC00"+string(rune('1'+i)),
"User "+string(rune('A'+i)),
"US",
)
repo.Create(account)
}

accounts, err := repo.List()
if err != nil {
t.Fatalf("Failed to list accounts: %v", err)
}

if len(accounts) != 3 {
t.Errorf("Expected 3 accounts, got %d", len(accounts))
}
})

t.Run("List empty repository", func(t *testing.T) {
repo := infrastructure.NewInMemoryAccountRepository()

accounts, err := repo.List()
if err != nil {
t.Fatalf("Failed to list accounts: %v", err)
}

if len(accounts) != 0 {
t.Errorf("Expected 0 accounts, got %d", len(accounts))
}
})
}
