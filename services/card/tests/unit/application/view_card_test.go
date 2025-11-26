package application_test

import (
	"testing"
	"time"

	"github.com/DavidRodriguez-create/pay-and-go/services/card/application"
	"github.com/DavidRodriguez-create/pay-and-go/services/card/domain"
)

func TestGetCardByID(t *testing.T) {
	t.Run("Successful retrieval", func(t *testing.T) {
		cardRepo := NewMockCardRepository()

		// Setup: Create a card
		card, _ := domain.NewCard("card-123", "US-12345", "US", "acc-123", time.Now())
		cardRepo.Create(card)

		useCase := application.NewViewCard(cardRepo)

		req := &application.GetCardRequest{
			ID: "card-123",
		}

		resp, err := useCase.GetByID(req)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if resp.ID != "card-123" {
			t.Errorf("Expected ID card-123, got %s", resp.ID)
		}
		if resp.CardNumber != "US-12345" {
			t.Errorf("Expected CardNumber US-12345, got %s", resp.CardNumber)
		}
	})

	t.Run("Missing ID", func(t *testing.T) {
		cardRepo := NewMockCardRepository()
		useCase := application.NewViewCard(cardRepo)

		req := &application.GetCardRequest{
			ID: "",
		}

		_, err := useCase.GetByID(req)

		if err != domain.ErrCardIDRequired {
			t.Errorf("Expected error %v, got %v", domain.ErrCardIDRequired, err)
		}
	})

	t.Run("Card not found", func(t *testing.T) {
		cardRepo := NewMockCardRepository()
		useCase := application.NewViewCard(cardRepo)

		req := &application.GetCardRequest{
			ID: "nonexistent",
		}

		_, err := useCase.GetByID(req)

		if err != domain.ErrCardNotFound {
			t.Errorf("Expected error %v, got %v", domain.ErrCardNotFound, err)
		}
	})
}

func TestGetCardByCardNumber(t *testing.T) {
	t.Run("Successful retrieval", func(t *testing.T) {
		cardRepo := NewMockCardRepository()

		// Setup: Create a card
		card, _ := domain.NewCard("card-123", "US-12345", "US", "acc-123", time.Now())
		cardRepo.Create(card)

		useCase := application.NewViewCard(cardRepo)

		req := &application.GetCardByNumberRequest{
			CardNumber: "US-12345",
		}

		resp, err := useCase.GetByCardNumber(req)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if resp.CardNumber != "US-12345" {
			t.Errorf("Expected CardNumber US-12345, got %s", resp.CardNumber)
		}
	})

	t.Run("Missing card number", func(t *testing.T) {
		cardRepo := NewMockCardRepository()
		useCase := application.NewViewCard(cardRepo)

		req := &application.GetCardByNumberRequest{
			CardNumber: "",
		}

		_, err := useCase.GetByCardNumber(req)

		if err != domain.ErrCardNumberRequired {
			t.Errorf("Expected error %v, got %v", domain.ErrCardNumberRequired, err)
		}
	})

	t.Run("Card not found", func(t *testing.T) {
		cardRepo := NewMockCardRepository()
		useCase := application.NewViewCard(cardRepo)

		req := &application.GetCardByNumberRequest{
			CardNumber: "nonexistent",
		}

		_, err := useCase.GetByCardNumber(req)

		if err != domain.ErrCardNotFound {
			t.Errorf("Expected error %v, got %v", domain.ErrCardNotFound, err)
		}
	})
}

func TestGetCardsByAccountID(t *testing.T) {
	t.Run("Successful retrieval with multiple cards", func(t *testing.T) {
		cardRepo := NewMockCardRepository()

		// Setup: Create multiple cards for same account
		card1, _ := domain.NewCard("card-1", "US-111", "US", "acc-123", time.Now())
		card2, _ := domain.NewCard("card-2", "US-222", "US", "acc-123", time.Now())
		card3, _ := domain.NewCard("card-3", "UK-333", "UK", "acc-456", time.Now())
		cardRepo.Create(card1)
		cardRepo.Create(card2)
		cardRepo.Create(card3)

		useCase := application.NewViewCard(cardRepo)

		req := &application.GetCardsByAccountRequest{
			AccountID: "acc-123",
		}

		resp, err := useCase.GetByAccountID(req)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if resp.Total != 2 {
			t.Errorf("Expected 2 cards, got %d", resp.Total)
		}
	})

	t.Run("Missing account ID", func(t *testing.T) {
		cardRepo := NewMockCardRepository()
		useCase := application.NewViewCard(cardRepo)

		req := &application.GetCardsByAccountRequest{
			AccountID: "",
		}

		_, err := useCase.GetByAccountID(req)

		if err != domain.ErrAccountIDRequired {
			t.Errorf("Expected error %v, got %v", domain.ErrAccountIDRequired, err)
		}
	})

	t.Run("No cards found", func(t *testing.T) {
		cardRepo := NewMockCardRepository()
		useCase := application.NewViewCard(cardRepo)

		req := &application.GetCardsByAccountRequest{
			AccountID: "acc-999",
		}

		resp, err := useCase.GetByAccountID(req)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if resp.Total != 0 {
			t.Errorf("Expected 0 cards, got %d", resp.Total)
		}
	})
}

func TestListCards(t *testing.T) {
	t.Run("Successful list with cards", func(t *testing.T) {
		cardRepo := NewMockCardRepository()

		// Setup: Create multiple cards
		card1, _ := domain.NewCard("card-1", "US-111", "US", "acc-123", time.Now())
		card2, _ := domain.NewCard("card-2", "US-222", "US", "acc-456", time.Now())
		cardRepo.Create(card1)
		cardRepo.Create(card2)

		useCase := application.NewListCards(cardRepo)

		resp, err := useCase.Execute()

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if resp.Total != 2 {
			t.Errorf("Expected 2 cards, got %d", resp.Total)
		}
	})

	t.Run("Empty list", func(t *testing.T) {
		cardRepo := NewMockCardRepository()
		useCase := application.NewListCards(cardRepo)

		resp, err := useCase.Execute()

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if resp.Total != 0 {
			t.Errorf("Expected 0 cards, got %d", resp.Total)
		}
	})

	t.Run("Repository error", func(t *testing.T) {
		cardRepo := NewMockCardRepository()
		cardRepo.listErr = domain.ErrCardNotFound

		useCase := application.NewListCards(cardRepo)

		_, err := useCase.Execute()

		if err == nil {
			t.Error("Expected repository error, got nil")
		}
	})
}
