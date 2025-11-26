package infrastructure_test

import (
	"testing"

	"github.com/DavidRodriguez-create/pay-and-go/services/card/domain"
	"github.com/DavidRodriguez-create/pay-and-go/services/card/infrastructure"
)

func TestMemoryAccountCacheRepository_Upsert(t *testing.T) {
	t.Run("Insert new account", func(t *testing.T) {
		repo := infrastructure.NewInMemoryAccountCacheRepository()

		account := domain.NewAccountCache("acc-123", "ACTIVE")

		err := repo.Upsert(account)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		// Verify account was created
		found, err := repo.GetByID("acc-123")
		if err != nil {
			t.Fatalf("Account not found after creation: %v", err)
		}
		if found.ID != "acc-123" {
			t.Errorf("Expected ID acc-123, got %s", found.ID)
		}
		if found.Status != "ACTIVE" {
			t.Errorf("Expected Status ACTIVE, got %s", found.Status)
		}
	})

	t.Run("Update existing account", func(t *testing.T) {
		repo := infrastructure.NewInMemoryAccountCacheRepository()

		// Create initial account
		account := domain.NewAccountCache("acc-123", "ACTIVE")
		repo.Upsert(account)

		// Update account status
		updatedAccount := domain.NewAccountCache("acc-123", "BLOCKED")
		err := repo.Upsert(updatedAccount)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		// Verify account was updated
		found, _ := repo.GetByID("acc-123")
		if found.Status != "BLOCKED" {
			t.Errorf("Expected Status BLOCKED, got %s", found.Status)
		}
	})

	t.Run("Upsert nil account", func(t *testing.T) {
		repo := infrastructure.NewInMemoryAccountCacheRepository()

		err := repo.Upsert(nil)

		if err == nil {
			t.Error("Expected error when upserting nil account, got nil")
		}
	})
}

func TestMemoryAccountCacheRepository_GetByID(t *testing.T) {
	t.Run("Successful retrieval", func(t *testing.T) {
		repo := infrastructure.NewInMemoryAccountCacheRepository()

		account := domain.NewAccountCache("acc-123", "ACTIVE")
		repo.Upsert(account)

		found, err := repo.GetByID("acc-123")

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if found.ID != "acc-123" {
			t.Errorf("Expected ID acc-123, got %s", found.ID)
		}
	})

	t.Run("Account not found", func(t *testing.T) {
		repo := infrastructure.NewInMemoryAccountCacheRepository()

		_, err := repo.GetByID("nonexistent")

		if err != domain.ErrAccountCacheNotFound {
			t.Errorf("Expected error %v, got %v", domain.ErrAccountCacheNotFound, err)
		}
	})

	t.Run("Empty ID", func(t *testing.T) {
		repo := infrastructure.NewInMemoryAccountCacheRepository()

		_, err := repo.GetByID("")

		if err == nil {
			t.Error("Expected error for empty ID, got nil")
		}
	})
}

func TestMemoryAccountCacheRepository_Delete(t *testing.T) {
	t.Run("Successful deletion", func(t *testing.T) {
		repo := infrastructure.NewInMemoryAccountCacheRepository()

		account := domain.NewAccountCache("acc-123", "ACTIVE")
		repo.Upsert(account)

		err := repo.Delete("acc-123")

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		// Verify account was deleted
		_, err = repo.GetByID("acc-123")
		if err != domain.ErrAccountCacheNotFound {
			t.Error("Account should not exist after deletion")
		}
	})

	t.Run("Delete nonexistent account", func(t *testing.T) {
		repo := infrastructure.NewInMemoryAccountCacheRepository()

		err := repo.Delete("nonexistent")

		if err != domain.ErrAccountCacheNotFound {
			t.Errorf("Expected error %v, got %v", domain.ErrAccountCacheNotFound, err)
		}
	})

	t.Run("Empty ID", func(t *testing.T) {
		repo := infrastructure.NewInMemoryAccountCacheRepository()

		err := repo.Delete("")

		if err == nil {
			t.Error("Expected error for empty ID, got nil")
		}
	})
}

func TestMemoryAccountCacheRepository_List(t *testing.T) {
	t.Run("List multiple accounts", func(t *testing.T) {
		repo := infrastructure.NewInMemoryAccountCacheRepository()

		account1 := domain.NewAccountCache("acc-123", "ACTIVE")
		account2 := domain.NewAccountCache("acc-456", "BLOCKED")
		repo.Upsert(account1)
		repo.Upsert(account2)

		accounts, err := repo.List()

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(accounts) != 2 {
			t.Errorf("Expected 2 accounts, got %d", len(accounts))
		}
	})

	t.Run("Empty repository", func(t *testing.T) {
		repo := infrastructure.NewInMemoryAccountCacheRepository()

		accounts, err := repo.List()

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(accounts) != 0 {
			t.Errorf("Expected 0 accounts, got %d", len(accounts))
		}
	})
}

func TestMemoryAccountCacheRepository_ConcurrentAccess(t *testing.T) {
	t.Run("Concurrent reads and writes", func(t *testing.T) {
		repo := infrastructure.NewInMemoryAccountCacheRepository()

		// Create initial account
		account := domain.NewAccountCache("acc-123", "ACTIVE")
		repo.Upsert(account)

		// Test concurrent access
		done := make(chan bool, 3)

		// Concurrent reader 1
		go func() {
			_, err := repo.GetByID("acc-123")
			if err != nil {
				t.Errorf("Concurrent read 1 failed: %v", err)
			}
			done <- true
		}()

		// Concurrent reader 2
		go func() {
			_, err := repo.List()
			if err != nil {
				t.Errorf("Concurrent read 2 failed: %v", err)
			}
			done <- true
		}()

		// Concurrent writer
		go func() {
			newAccount := domain.NewAccountCache("acc-456", "BLOCKED")
			err := repo.Upsert(newAccount)
			if err != nil {
				t.Errorf("Concurrent write failed: %v", err)
			}
			done <- true
		}()

		// Wait for all goroutines
		for i := 0; i < 3; i++ {
			<-done
		}
	})

	t.Run("Concurrent upserts on same account", func(t *testing.T) {
		repo := infrastructure.NewInMemoryAccountCacheRepository()

		// Test concurrent upserts on the same account
		done := make(chan bool, 2)

		go func() {
			account := domain.NewAccountCache("acc-123", "ACTIVE")
			repo.Upsert(account)
			done <- true
		}()

		go func() {
			account := domain.NewAccountCache("acc-123", "BLOCKED")
			repo.Upsert(account)
			done <- true
		}()

		// Wait for both goroutines
		<-done
		<-done

		// Verify account exists (final state may be either ACTIVE or BLOCKED)
		found, err := repo.GetByID("acc-123")
		if err != nil {
			t.Error("Account should exist after concurrent upserts")
		}
		if found.Status != "ACTIVE" && found.Status != "BLOCKED" {
			t.Errorf("Unexpected status: %s", found.Status)
		}
	})
}
