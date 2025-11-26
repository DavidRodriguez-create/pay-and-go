package infrastructure_test

import (
	"testing"
	"time"

	"github.com/DavidRodriguez-create/pay-and-go/services/card/domain"
	"github.com/DavidRodriguez-create/pay-and-go/services/card/infrastructure"
)

func TestMemoryCardRepository_Create(t *testing.T) {
	t.Run("Successful card creation", func(t *testing.T) {
		repo := infrastructure.NewInMemoryCardRepository()

		card, _ := domain.NewCard("card-123", "US-12345", "US", "acc-123", time.Now())

		err := repo.Create(card)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		// Verify card was created
		found, err := repo.GetByID("card-123")
		if err != nil {
			t.Fatalf("Card not found after creation: %v", err)
		}
		if found.ID != "card-123" {
			t.Errorf("Expected ID card-123, got %s", found.ID)
		}
	})

	t.Run("Create nil card", func(t *testing.T) {
		repo := infrastructure.NewInMemoryCardRepository()

		err := repo.Create(nil)

		if err == nil {
			t.Error("Expected error when creating nil card, got nil")
		}
	})
}

func TestMemoryCardRepository_GetByID(t *testing.T) {
	t.Run("Successful retrieval", func(t *testing.T) {
		repo := infrastructure.NewInMemoryCardRepository()

		card, _ := domain.NewCard("card-123", "US-12345", "US", "acc-123", time.Now())
		repo.Create(card)

		found, err := repo.GetByID("card-123")

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if found.ID != "card-123" {
			t.Errorf("Expected ID card-123, got %s", found.ID)
		}
	})

	t.Run("Card not found", func(t *testing.T) {
		repo := infrastructure.NewInMemoryCardRepository()

		_, err := repo.GetByID("nonexistent")

		if err != domain.ErrCardNotFound {
			t.Errorf("Expected error %v, got %v", domain.ErrCardNotFound, err)
		}
	})

	t.Run("Empty ID", func(t *testing.T) {
		repo := infrastructure.NewInMemoryCardRepository()

		_, err := repo.GetByID("")

		if err == nil {
			t.Error("Expected error for empty ID, got nil")
		}
	})
}

func TestMemoryCardRepository_GetByCardNumber(t *testing.T) {
	t.Run("Successful retrieval", func(t *testing.T) {
		repo := infrastructure.NewInMemoryCardRepository()

		card, _ := domain.NewCard("card-123", "US-12345", "US", "acc-123", time.Now())
		repo.Create(card)

		found, err := repo.GetByCardNumber("US-12345")

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if found.CardNumber != "US-12345" {
			t.Errorf("Expected CardNumber US-12345, got %s", found.CardNumber)
		}
	})

	t.Run("Card not found", func(t *testing.T) {
		repo := infrastructure.NewInMemoryCardRepository()

		_, err := repo.GetByCardNumber("nonexistent")

		if err != domain.ErrCardNotFound {
			t.Errorf("Expected error %v, got %v", domain.ErrCardNotFound, err)
		}
	})

	t.Run("Empty card number", func(t *testing.T) {
		repo := infrastructure.NewInMemoryCardRepository()

		_, err := repo.GetByCardNumber("")

		if err == nil {
			t.Error("Expected error for empty card number, got nil")
		}
	})
}

func TestMemoryCardRepository_GetByAccountID(t *testing.T) {
	t.Run("Multiple cards for same account", func(t *testing.T) {
		repo := infrastructure.NewInMemoryCardRepository()

		card1, _ := domain.NewCard("card-1", "US-111", "US", "acc-123", time.Now())
		card2, _ := domain.NewCard("card-2", "US-222", "US", "acc-123", time.Now())
		card3, _ := domain.NewCard("card-3", "UK-333", "UK", "acc-456", time.Now())
		repo.Create(card1)
		repo.Create(card2)
		repo.Create(card3)

		cards, err := repo.GetByAccountID("acc-123")

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(cards) != 2 {
			t.Errorf("Expected 2 cards, got %d", len(cards))
		}
	})

	t.Run("No cards for account", func(t *testing.T) {
		repo := infrastructure.NewInMemoryCardRepository()

		cards, err := repo.GetByAccountID("acc-999")

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(cards) != 0 {
			t.Errorf("Expected 0 cards, got %d", len(cards))
		}
	})

	t.Run("Empty account ID", func(t *testing.T) {
		repo := infrastructure.NewInMemoryCardRepository()

		_, err := repo.GetByAccountID("")

		if err == nil {
			t.Error("Expected error for empty account ID, got nil")
		}
	})
}

func TestMemoryCardRepository_List(t *testing.T) {
	t.Run("List multiple cards", func(t *testing.T) {
		repo := infrastructure.NewInMemoryCardRepository()

		card1, _ := domain.NewCard("card-1", "US-111", "US", "acc-123", time.Now())
		card2, _ := domain.NewCard("card-2", "US-222", "US", "acc-456", time.Now())
		repo.Create(card1)
		repo.Create(card2)

		cards, err := repo.List()

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(cards) != 2 {
			t.Errorf("Expected 2 cards, got %d", len(cards))
		}
	})

	t.Run("Empty repository", func(t *testing.T) {
		repo := infrastructure.NewInMemoryCardRepository()

		cards, err := repo.List()

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(cards) != 0 {
			t.Errorf("Expected 0 cards, got %d", len(cards))
		}
	})
}

func TestMemoryCardRepository_Delete(t *testing.T) {
	t.Run("Successful soft delete", func(t *testing.T) {
		repo := infrastructure.NewInMemoryCardRepository()

		card, _ := domain.NewCard("card-123", "US-12345", "US", "acc-123", time.Now())
		repo.Create(card)

		err := repo.Delete("card-123")

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		// Verify card is marked as deleted
		deletedCard, _ := repo.GetByID("card-123")
		if !deletedCard.Deleted {
			t.Error("Card should be marked as deleted")
		}
	})

	t.Run("Delete nonexistent card", func(t *testing.T) {
		repo := infrastructure.NewInMemoryCardRepository()

		err := repo.Delete("nonexistent")

		if err != domain.ErrCardNotFound {
			t.Errorf("Expected error %v, got %v", domain.ErrCardNotFound, err)
		}
	})

	t.Run("Delete empty ID", func(t *testing.T) {
		repo := infrastructure.NewInMemoryCardRepository()

		err := repo.Delete("")

		if err == nil {
			t.Error("Expected error for empty ID, got nil")
		}
	})
}

func TestMemoryCardRepository_ConcurrentAccess(t *testing.T) {
	t.Run("Concurrent reads and writes", func(t *testing.T) {
		repo := infrastructure.NewInMemoryCardRepository()

		// Create initial card
		card, _ := domain.NewCard("card-123", "US-12345", "US", "acc-123", time.Now())
		repo.Create(card)

		// Test concurrent access
		done := make(chan bool, 3)

		// Concurrent reader 1
		go func() {
			_, err := repo.GetByID("card-123")
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
			newCard, _ := domain.NewCard("card-456", "US-67890", "US", "acc-456", time.Now())
			err := repo.Create(newCard)
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
}
