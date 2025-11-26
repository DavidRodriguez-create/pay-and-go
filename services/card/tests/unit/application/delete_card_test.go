package application_test

import (
	"testing"
	"time"

	"github.com/DavidRodriguez-create/pay-and-go/services/card/application"
	"github.com/DavidRodriguez-create/pay-and-go/services/card/domain"
)

func TestDeleteCard(t *testing.T) {
	t.Run("Successful card deletion", func(t *testing.T) {
		cardRepo := NewMockCardRepository()

		// Setup: Create a card
		card, _ := domain.NewCard("card-123", "US-12345", "US", "acc-123", time.Now())
		cardRepo.Create(card)

		useCase := application.NewDeleteCard(cardRepo)

		req := &application.DeleteCardRequest{
			ID: "card-123",
		}

		err := useCase.Execute(req)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		// Verify card is marked as deleted
		deletedCard, _ := cardRepo.GetByID("card-123")
		if !deletedCard.Deleted {
			t.Error("Card should be marked as deleted")
		}
	})

	t.Run("Missing card ID", func(t *testing.T) {
		cardRepo := NewMockCardRepository()
		useCase := application.NewDeleteCard(cardRepo)

		req := &application.DeleteCardRequest{
			ID: "",
		}

		err := useCase.Execute(req)

		if err != domain.ErrCardIDRequired {
			t.Errorf("Expected error %v, got %v", domain.ErrCardIDRequired, err)
		}
	})

	t.Run("Card not found", func(t *testing.T) {
		cardRepo := NewMockCardRepository()
		useCase := application.NewDeleteCard(cardRepo)

		req := &application.DeleteCardRequest{
			ID: "nonexistent",
		}

		err := useCase.Execute(req)

		if err != domain.ErrCardNotFound {
			t.Errorf("Expected error %v, got %v", domain.ErrCardNotFound, err)
		}
	})

	t.Run("Delete already deleted card", func(t *testing.T) {
		cardRepo := NewMockCardRepository()

		// Setup: Create and delete a card
		card, _ := domain.NewCard("card-123", "US-12345", "US", "acc-123", time.Now())
		card.Delete()
		cardRepo.Create(card)

		useCase := application.NewDeleteCard(cardRepo)

		req := &application.DeleteCardRequest{
			ID: "card-123",
		}

		err := useCase.Execute(req)

		if err != domain.ErrCardAlreadyDeleted {
			t.Errorf("Expected error %v, got %v", domain.ErrCardAlreadyDeleted, err)
		}
	})

	t.Run("Repository delete error", func(t *testing.T) {
		cardRepo := NewMockCardRepository()
		cardRepo.deleteErr = domain.ErrCardNotFound

		// Setup: Create a card
		card, _ := domain.NewCard("card-123", "US-12345", "US", "acc-123", time.Now())
		cardRepo.Create(card)

		// Clear the error for GetByID to succeed
		cardRepo.deleteErr = nil
		useCase := application.NewDeleteCard(cardRepo)

		// Set error for Delete operation
		cardRepo.deleteErr = domain.ErrCardNotFound

		req := &application.DeleteCardRequest{
			ID: "card-123",
		}

		err := useCase.Execute(req)

		if err == nil {
			t.Error("Expected repository error, got nil")
		}
	})
}
